// Copyright 2021 Adam Chalkley
//
// https://github.com/atc0005/check-vmware
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package vsphere

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/atc0005/check-vmware/internal/textutils"
	"github.com/atc0005/go-nagios"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
)

// ErrHostDatastorePairingFailed is returned when compiling host and datastore
// pairings using provided the Custom Attribute fails. This is usually due to
// a lack of a match between Custom Attribute values used on hosts and
// datastores.
//
// For example, this may occur if prefix matching is not enabled and the host
// Custom Attribute uses a Location attribute which contains a separator
// between a datacenter and the hosts rack position, whereas a datastore
// contains only the datacenter.
var ErrHostDatastorePairingFailed = errors.New("failed to compile host/datastore pairings")

// ErrDatastoreIDToNameLookupFailed is returned when a search of the host to
// datastore index fails to yield a name for a specified ID value. This is
// expected to be an unusual scenario.
// var ErrDatastoreIDToNameLookupFailed = errors.New("failed to find a matching Datastore name for provided ID")

// ErrHostDatastoreIdxIDToNameLookupFailed is returned when a search of the host to
// datastore index fails to yield a name for a specified ID value. This can
// occur if the datastore for a VM is in a user-specified ignored datastores
// list.
type ErrHostDatastoreIdxIDToNameLookupFailed struct {
	Err         error
	DatastoreID string
}

func (dsIDFail ErrHostDatastoreIdxIDToNameLookupFailed) Error() string {
	return fmt.Sprintf(
		"id: %v; %v",
		dsIDFail.DatastoreID,
		dsIDFail.Err,
	)
}

// func (dsIDFail ErrDatastoreIDToNameLookupFailed) ID() string {
// 	return dsIDFail.DatastoreID
// }

// ErrVMDatastoreNotInVMHostPairedList is returned when one or more datastores
// for a VirtualMachine are not in the list of datastores paired with the
// VirtualMachine's current host.
var ErrVMDatastoreNotInVMHostPairedList = errors.New("host/datastore/vm mismatch")

// ErrManagedObjectIDIsNil indicates that a managed object ID is unset, which
// may occur if the property is not requested from the vSphere API or if the
// service account executing the plugin has insufficient privileges.
var ErrManagedObjectIDIsNil = errors.New("managed object ID is nil")

// ErrManagedObjectIDIsEmpty indicates that a managed object ID is empty.
var ErrManagedObjectIDIsEmpty = errors.New("managed object ID is empty")

// DatastoreWithCA wraps the vSphere Datastore managed object type with a
// specific Custom Attribute name/value pair. This Custom Attribute is
// intended to link this datastore to a specific ESXi host.
type DatastoreWithCA struct {
	mo.Datastore

	// CustomAttribute represents the name/value pair used to relate specific
	// hosts with specific datastores. Most often this takes the form of a
	// "Location" or "Datacenter" field to indicate which Datastore a
	// VirtualMachine should reside on when running on a specific HostSystem.
	CustomAttribute CustomAttribute
}

// HostWithCA wraps the vSphere HostSystem managed object type with a specific
// Custom Attribute name/value pair. This Custom Attribute is intended to link
// this host to one or more datastores.
type HostWithCA struct {
	mo.HostSystem

	// CustomAttribute represents the name/value pair used to relate specific
	// hosts with specific datastores. Most often this takes the form of a
	// "Location" or "Datacenter" field to indicate which Datastore a
	// VirtualMachine should reside on when running on a specific HostSystem.
	CustomAttribute CustomAttribute
}

// HostDatastoresPairing collects Host and Datastores pairings based on shared
// Custom Attribute name and value (literal) or prefix (if user-specified).
// This is intended to "pair" hosts and datastores (using a specific Custom
// Attribute) within an environment that are known to work well together.
type HostDatastoresPairing struct {
	Host       HostWithCA
	Datastores []DatastoreWithCA
}

// HostToDatastoreIndex indexes HostDatastorePairings based on host id values.
type HostToDatastoreIndex map[string]HostDatastoresPairing

// VMToMismatchedPairing indexes VirtualMachine name to a pairing of host and
// datastores. This index reflects mismatched Datastores based on the current
// host for the VirtualMachine.
type VMToMismatchedPairing map[string]HostDatastoresPairing

// GetVMDatastorePairingIssues receives a list of VirtualMachines, a
// HostToDatastoreIndex to evaluate each VirtualMachine by, a list of all
// datastores and a list of datastore names which should be ignored. A
// VMHostDatastoresPairing index is returned noting improperly paired
// VirtualMachines (if any) and an error (if applicable).
//
// NOTE: The full list of datastores is used instead of the list of datastores
// with specified pairing custom attribute. This is because some datastores
// may not have the specified custom attribute and the user may have opted to
// ignore any datastores or hosts which do not have it.
func GetVMDatastorePairingIssues(vms []mo.VirtualMachine, h2dIdx HostToDatastoreIndex, dss []mo.Datastore, ignoredDatastoreNames []string) (VMToMismatchedPairing, error) {

	funcTimeStart := time.Now()

	vmDatastoresPairingIssues := make(VMToMismatchedPairing)

	defer func(issues *VMToMismatchedPairing) {
		logger.Printf(
			"It took %v to execute GetVMDatastorePairingIssues func (and retrieve VMToMismatchedPairing idx [hosts: %d, datastores: %d]).\n",
			time.Since(funcTimeStart),
			issues.NumHosts(),
			issues.NumDatastores(),
		)
	}(&vmDatastoresPairingIssues)

	for _, vm := range vms {

		// Assert that we can retrieve the required hostMOID for the VM.
		hostMOID, lookupErr := getVMHostID(vm)
		if lookupErr != nil {
			return nil, lookupErr
		}

		vmHostDatastoresPairing, ok := h2dIdx[hostMOID]
		if !ok {
			// FAILURE due to host id lookup; this should not occur since we (as
			// of GH-393) create stub entries for hosts without matching
			// datastores (via specified custom attribute).
			errMsg := "error retrieving host/datastores pairing using Host MOID " + hostMOID

			logger.Print(errMsg)

			return nil, errors.New(errMsg)
		}

		mismatchedDatastores, lookupErr := h2dIdx.ValidateVirtualMachinePairings(
			vm,
			dss,
			ignoredDatastoreNames,
		)

		if lookupErr != nil {
			errMsg := fmt.Sprintf(
				"error occurred while validating VM/Host/Datastore match for vm %s and host %s",
				vm.Name,
				vmHostDatastoresPairing.Host.Name,
			)

			logger.Print(errMsg)

			return nil, errors.New(errMsg)
		}

		switch {
		case len(mismatchedDatastores) > 0:

			logger.Printf(
				"VM/Host/Datastore validation failed for VM %q on host %q",
				vm.Name,
				vmHostDatastoresPairing.Host.Name,
			)

			// Retrieve the complete list of datastore names for the VM,
			// regardless of any potential mismatches between host/datastore.
			vmDatastoreNames := DatastoreIDsToNames(vm.Datastore, dss)
			logger.Printf(
				"All datastores for VM %q: %q",
				vm.Name,
				strings.Join(vmDatastoreNames, ", "),
			)

			mismatchedDatastoreNames := make([]string, 0, len(mismatchedDatastores))
			for _, ds := range mismatchedDatastores {
				mismatchedDatastoreNames = append(mismatchedDatastoreNames, ds.Name)
			}

			logger.Printf(
				"VM/Datastores mismatched: %q",
				strings.Join(mismatchedDatastoreNames, ", "),
			)

			// index mismatched datastore names by VirtualMachine name, also
			// for later review
			vmDatastoresPairingIssues[vm.Name] = HostDatastoresPairing{
				Host:       vmHostDatastoresPairing.Host,
				Datastores: mismatchedDatastores,
			}

		default:

			logger.Printf(
				"All datastores for VM %q matched to host %q",
				vm.Name,
				vmHostDatastoresPairing.Host.Name,
			)

		}
	}

	return vmDatastoresPairingIssues, nil

}

// GetHostsWithCA receives a collection of Hosts, a Custom Attribute name to
// filter Hosts by and a boolean flag indicating whether Hosts missing a
// Custom Attribute should be ignored. A collection of HostWithCA is returned
// along with an error (if applicable).
func GetHostsWithCA(allHosts []mo.HostSystem, hostCustomAttributeName string, ignoreMissingCA bool) ([]HostWithCA, error) {

	funcTimeStart := time.Now()

	hostsWithCAs := make([]HostWithCA, 0, len(allHosts))

	defer func(hosts *[]HostWithCA) {
		logger.Printf(
			"It took %v to execute GetHostsWithCA func (and retrieve %d HostWithCAs).\n",
			time.Since(funcTimeStart),
			len(*hosts),
		)
	}(&hostsWithCAs)

	hostsMissingCAs := make([]string, 0, len(allHosts))
	for _, host := range allHosts {
		ca, err := GetObjectCustomAttribute(host.ManagedEntity, hostCustomAttributeName, ignoreMissingCA)
		switch {
		case errors.Is(err, ErrCustomAttributeNotSet):
			logger.Printf(
				"custom attributes missing for %s %q",
				host.ManagedEntity.Self.Type,
				host.Name,
			)

			hostsMissingCAs = append(hostsMissingCAs, host.Name)

		case err != nil:
			logger.Printf(
				"failed to retrieve custom attribute for %s %q: %s",
				host.ManagedEntity.Self.Type,
				host.Name,
				err,
			)

			// Unknown error occurred. Don't batch these retrieval errors,
			// report them immediately.
			return nil, fmt.Errorf(
				"failed to retrieve custom attribute for %s %q: %w",
				host.ManagedEntity.Self.Type,
				host.Name,
				err,
			)

		default:
			logger.Printf(
				"successfully retrieved custom attribute for %s %q",
				host.ManagedEntity.Self.Type,
				host.Name,
			)

			hostsWithCAs = append(hostsWithCAs, HostWithCA{
				HostSystem:      host,
				CustomAttribute: ca,
			})

		}

	}

	switch {
	case len(hostsMissingCAs) > 0:
		return nil, fmt.Errorf(
			"failed to retrieve custom attribute %q from hosts: [%v]: %w",
			hostCustomAttributeName,
			strings.Join(hostsMissingCAs, ", "),
			ErrCustomAttributeNotSet,
		)

	default:
		return hostsWithCAs, nil
	}

}

// GetDatastoresWithCA receives a collection of Datastores, a list of
// datastore names that should be ignored or excluded from evaluation, a
// Custom Attribute name to filter Datastores by and a boolean flag indicating
// whether datastores missing a Custom Attribute should be ignored. A
// collection of DatastoreWithCA is returned along with an error (if
// applicable).
func GetDatastoresWithCA(allDS []mo.Datastore, ignoredDatastoreNames []string, dsCustomAttributeName string, ignoreMissingCA bool) ([]DatastoreWithCA, error) {

	funcTimeStart := time.Now()

	dsNames := make([]string, 0, len(allDS))
	for _, ds := range allDS {
		dsNames = append(dsNames, ds.Name)
	}

	datastoresWithCA := make([]DatastoreWithCA, 0, len(allDS))

	defer func(dss *[]DatastoreWithCA) {
		logger.Printf(
			"It took %v to execute GetDatastoresWithCA func (and retrieve %d DatastoreWithCAs).\n",
			time.Since(funcTimeStart),
			len(*dss),
		)
	}(&datastoresWithCA)

	// validate the list of ignored datastores
	if len(ignoredDatastoreNames) > 0 {
		for _, ignDSName := range ignoredDatastoreNames {
			if !textutils.InList(ignDSName, dsNames, true) {

				validateIgnoredDSErr := errors.New(
					"error validating list of ignored datastores",
				)
				validateIgnoredDSErrMsg := fmt.Sprintf(
					"datastore %s could not be ignored as requested; "+
						"could not locate datastore within vSphere inventory",
					ignDSName,
				)

				logger.Printf(
					"%v: %v",
					validateIgnoredDSErr,
					validateIgnoredDSErrMsg,
				)

				return nil, fmt.Errorf(
					"%v: %v",
					validateIgnoredDSErr,
					validateIgnoredDSErrMsg,
				)

			}
		}
	}

	dsMissingCAs := make([]string, 0, len(allDS))
	for _, ds := range allDS {

		// if user opted to ignore the Datastore, skip attempts to retrieve
		// Custom Attribute for it.
		if textutils.InList(ds.Name, ignoredDatastoreNames, true) {
			continue
		}

		ca, err := GetObjectCustomAttribute(ds.ManagedEntity, dsCustomAttributeName, ignoreMissingCA)
		switch {
		case errors.Is(err, ErrCustomAttributeNotSet):
			logger.Printf(
				"custom attributes missing for %s %q",
				ds.ManagedEntity.Self.Type,
				ds.Name,
			)

			dsMissingCAs = append(dsMissingCAs, ds.Name)

		case err != nil:
			logger.Printf(
				"failed to retrieve custom attribute for %s %q: %s",
				ds.ManagedEntity.Self.Type,
				ds.Name,
				err,
			)

			// Unknown error occurred. Don't batch these retrieval errors,
			// report them immediately.
			return nil, fmt.Errorf(
				"failed to retrieve custom attribute for %s %q: %w",
				ds.ManagedEntity.Self.Type,
				ds.Name,
				err,
			)

		default:
			logger.Printf(
				"successfully retrieved custom attribute for %s %q",
				ds.ManagedEntity.Self.Type,
				ds.Name,
			)

			datastoresWithCA = append(datastoresWithCA, DatastoreWithCA{
				Datastore:       ds,
				CustomAttribute: ca,
			})

		}

	}

	switch {
	case len(dsMissingCAs) > 0:
		return nil, fmt.Errorf(
			"failed to retrieve custom attribute %q from datastores: [%v]: %w",
			dsCustomAttributeName,
			strings.Join(dsMissingCAs, ", "),
			ErrCustomAttributeNotSet,
		)

	default:
		return datastoresWithCA, nil
	}

}

// NewHostToDatastoreIndex receives a collection of hosts and datastores
// wrapped with user-specified Custom Attributes, prefix separators and a
// boolean flag indicating whether prefix matching will be used.
//
// The index is created using each ESXi host's MOID as the key and a
// HostDatastoresPairing type as the value. This effectively provides a
// mapping between a host and all datastores with matching specified custom
// attribute. If no datastores are found with a matching specified custom
// attribute, an empty list of datastores is recorded for the host to indicate
// this.
//
// The resulting HostToDatastoreIndex is returned if no errors occur,
// otherwise nil and the error.
func NewHostToDatastoreIndex(
	hosts []HostWithCA,
	datastores []DatastoreWithCA,
	usingPrefixes bool,
	hostCASep string,
	datastoreCASep string,
) (HostToDatastoreIndex, error) {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute NewHostToDatastoreIndex func.\n",
			time.Since(funcTimeStart),
		)
	}()

	h2dIdx := make(HostToDatastoreIndex)

	// Review incoming hosts slice for problem entries.
	for i, host := range hosts {
		logger.Printf(
			"index: %d, hostname: %v, hostID: %v",
			i,
			host.Name,
			host.Self.Value,
		)
	}

	// Review datastore slice for problem entries.
	for i, datastore := range datastores {
		logger.Printf(
			"index: %d, datastore name: %v, datastore ID: %v",
			i,
			datastore.Name,
			datastore.Self.Value,
		)
	}

	// Ensure that we were given some useful values to work with, otherwise
	// abort early.
	switch {
	case len(hosts) == 0:
		return nil, errors.New("empty hosts list provided; at least one host is required")
	case len(datastores) == 0:
		return nil, errors.New("empty datastores list provided; at least one datastore is required")
	case usingPrefixes && hostCASep == "":
		return nil, errors.New("missing host custom attribute prefix; prefix is required if using attribute prefix matching")
	case usingPrefixes && datastoreCASep == "":
		return nil, errors.New("missing datastore custom attribute prefix; prefix is required if using attribute prefix matching")
	}

	for _, host := range hosts {

		hostID := host.Self.Value

		for _, datastore := range datastores {

			// FIXME: Generating the hostCAValPrefix for each datastore is
			// inefficient, but not a serious problem. Review this with
			// future refactoring work.
			switch {
			case usingPrefixes:
				hostCAValPrefix := strings.SplitN(
					host.CustomAttribute.Value,
					hostCASep,
					2,
				)[0]

				datastoreCAValPrefix := strings.SplitN(
					datastore.CustomAttribute.Value,
					datastoreCASep,
					2,
				)[0]

				if strings.EqualFold(datastoreCAValPrefix, hostCAValPrefix) {
					h2dIdx[hostID] = HostDatastoresPairing{
						Host:       host,
						Datastores: append(h2dIdx[hostID].Datastores, datastore),
					}

					logger.Printf(
						"successful match using prefix equal fold between datastore %q (%q) and host %q (%q)",
						datastore.Name,
						datastoreCAValPrefix,
						host.Name,
						hostCAValPrefix,
					)

					continue
				}

				logger.Printf(
					"failed match using prefix equal fold between datastore %q (%q) and host %q (%q)",
					datastore.Name,
					datastoreCAValPrefix,
					host.Name,
					hostCAValPrefix,
				)

			default:
				// not using prefixes, so literal values
				if strings.EqualFold(datastore.CustomAttribute.Value, host.CustomAttribute.Value) {
					h2dIdx[hostID] = HostDatastoresPairing{
						Host:       host,
						Datastores: append(h2dIdx[hostID].Datastores, datastore),
					}

					logger.Printf(
						"successful match using literal equal fold between datastore %q (%q) and host %q (%q)",
						datastore.Name,
						datastore.CustomAttribute.Value,
						host.Name,
						host.CustomAttribute.Value,
					)

					continue
				}

				logger.Printf(
					"failed match using literal equal fold between datastore %q (%q) and host %q (%q)",
					datastore.Name,
					datastore.CustomAttribute.Value,
					host.Name,
					host.CustomAttribute.Value,
				)
			}
		}

		// If we did not find any matching datastores for the host (via
		// specified custom attribute) note as much by adding a
		// HostDatastoresPairing entry to the index for the host with an empty
		// datastores list.
		if _, ok := h2dIdx[hostID]; !ok {
			logger.Printf(
				"Adding zero value entry for host [name: %s, id: %s]",
				host.Name,
				hostID,
			)

			h2dIdx[hostID] = HostDatastoresPairing{
				Host:       host,
				Datastores: []DatastoreWithCA{},
			}
		}

	}

	// Baseline validation check. Unless an invalid custom attribute value was
	// provided we should have at last one pairing for the provided hosts and
	// datastores.
	if len(h2dIdx) == 0 {
		return nil, ErrHostDatastorePairingFailed
	}

	return h2dIdx, nil

}

// DatastoreNames returns a list of all Datastore names in the index.
func (hdi HostToDatastoreIndex) DatastoreNames() []string {

	funcTimeStart := time.Now()

	var dsNames []string

	defer func(dss *[]string) {
		logger.Printf(
			"It took %v to execute DatastoreNames func (and retrieve %d Datastore names).\n",
			time.Since(funcTimeStart),
			len(*dss),
		)
	}(&dsNames)

	for hostID := range hdi {
		for _, ds := range hdi[hostID].Datastores {
			dsNames = append(dsNames, ds.Name)
		}
	}

	sort.Slice(dsNames, func(i, j int) bool {
		return strings.ToLower(dsNames[i]) < strings.ToLower(dsNames[j])
	})

	return dsNames

}

// DatastoreIDToNameIndex returns an index of all Datastore IDs to names in the index.
func (hdi HostToDatastoreIndex) DatastoreIDToNameIndex() DatastoreIDToNameIndex {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute DatastoreIDToNameIndex func.\n",
			time.Since(funcTimeStart),
		)
	}()

	dsIdx := make(DatastoreIDToNameIndex)
	for hostID := range hdi {
		for _, ds := range hdi[hostID].Datastores {
			dsIdx[ds.Self.Value] = ds.Name
		}
	}

	return dsIdx

}

// IsDatastoreIDInIndex indicates whether a provided Datastore ID is in the
// index.
func (hdi HostToDatastoreIndex) IsDatastoreIDInIndex(dsID string) bool {

	for hostID := range hdi {
		for _, ds := range hdi[hostID].Datastores {
			if strings.EqualFold(dsID, ds.Self.Value) {
				return true
			}
		}
	}

	return false

}

// DatastoreIDToName returns the name associated with a Datastore ID. An error
// is returned if the name could not be retrieved from the index.
func (hdi HostToDatastoreIndex) DatastoreIDToName(dsID string) (string, error) {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute DatastoreIDToName func.\n",
			time.Since(funcTimeStart),
		)
	}()

	for hostID := range hdi {
		for _, ds := range hdi[hostID].Datastores {
			if ds.Self.Value == dsID {
				return ds.Name, nil
			}
		}
	}

	return "", &ErrHostDatastoreIdxIDToNameLookupFailed{
		DatastoreID: dsID,
		Err:         errors.New("datastore ID not found"),
	}
}

// DatastoreWithCAByID returns the DatastoreWithCA value associated with a
// Datastore ID. An error is returned if the DatastoreWithCA value was not
// found in the index.
func (hdi HostToDatastoreIndex) DatastoreWithCAByID(dsID string) (DatastoreWithCA, error) {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute DatastoreWithCAByID func.\n",
			time.Since(funcTimeStart),
		)
	}()

	for hostID := range hdi {
		for _, ds := range hdi[hostID].Datastores {
			if ds.Self.Value == dsID {
				return ds, nil
			}
		}
	}

	return DatastoreWithCA{}, &ErrHostDatastoreIdxIDToNameLookupFailed{
		DatastoreID: dsID,
		Err:         errors.New("datastore ID not found"),
	}
}

// ValidateVirtualMachinePairings receives a VirtualMachine host ID, a
// collection of all Datastores, the VirtualMachine, and an optional list of
// Datastore names to ignore. An index of mismatched virtual machine to host
// and datastore pairings (if applicable) is returned or an error if one
// occurs.
func (hdi HostToDatastoreIndex) ValidateVirtualMachinePairings(
	vm mo.VirtualMachine,
	allDatastores []mo.Datastore,
	dsNamesToIgnore []string,
) ([]DatastoreWithCA, error) {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute ValidateVirtualMachinePairings func.\n",
			time.Since(funcTimeStart),
		)
	}()

	// Assert that we can retrieve the required host ID for the VM.
	vmHostID, lookupErr := getVMHostID(vm)
	if lookupErr != nil {
		return nil, lookupErr
	}

	vmDatastoreIDs := make([]string, 0, len(vm.Datastore))
	for _, vmDSRef := range vm.Datastore {
		vmDatastoreIDs = append(vmDatastoreIDs, vmDSRef.Value)
	}

	hostDatastoreIDs := make([]string, 0, len(hdi[vmHostID].Datastores))
	for _, hostPairedDS := range hdi[vmHostID].Datastores {
		hostDatastoreIDs = append(hostDatastoreIDs, hostPairedDS.Self.Value)
	}

	// Assert that every datastore ID associated with the VM is within the
	// list of datastores associated with the current host for the VM.
	vmDatastoreIDMismatches := make([]string, 0, 10)
	for _, vmDatastoreID := range vmDatastoreIDs {
		if !textutils.InList(vmDatastoreID, hostDatastoreIDs, true) {
			vmDatastoreIDMismatches = append(vmDatastoreIDMismatches, vmDatastoreID)
		}
	}

	// There are no mismatches; all VMs are running on a host properly paired
	// with the VM's datastores. No further evaluation is necessary.
	if len(vmDatastoreIDMismatches) == 0 {
		return nil, nil
	}

	// For every datastore associated with the VM which isn't paired against
	// the currently running host, add it to the mismatch index.
	datastoreMismatches := make([]DatastoreWithCA, 0, 5)
	for _, vmDatastoreID := range vmDatastoreIDMismatches {

		// Lookup VM datastore in the index via ID.
		ds, lookupErr := hdi.DatastoreWithCAByID(vmDatastoreID)
		if lookupErr != nil {

			// The first lookup attempt could have failed if the sole
			// datastore for the VM is in the ignored list; double-check
			// that possibility before reporting the lookup failure.
			var dsIDLookupErr *ErrHostDatastoreIdxIDToNameLookupFailed
			if errors.As(lookupErr, &dsIDLookupErr) {
				logger.Printf(
					"Initial lookup failed for datastore ID %s\n",
					vmDatastoreID,
				)

				datastore, _, filterErr := FilterDatastoresByID(allDatastores, vmDatastoreID)
				if filterErr != nil {
					// This is our second lookup attempt using the
					// datastore id. The first failure is because the
					// datastore isn't in our host-to-datastore index,
					// this second because it could not be located in the
					// full datastores list from the vSphere inventory.
					return nil, fmt.Errorf(
						"second lookup attempt unsuccessful for datastore"+
							" ID %q for VM %q; failed to locate datastore ID in "+
							"index or full datastores list: %w",
						vmDatastoreID,
						vm.Name,
						filterErr,
					)
				}

				logger.Printf(
					"Second lookup successful; name: %q id: %q",
					datastore.Name,
					vmDatastoreID,
				)

				// Resolved datastore name is in the ignore list, skip it.
				if textutils.InList(datastore.Name, dsNamesToIgnore, true) {
					continue
				}
			}

			// Lookup failure occurred for some other reason.
			return nil, lookupErr
		}

		switch {
		case textutils.InList(ds.Name, dsNamesToIgnore, true):
			// if datastore name is in the ignore list, don't report
			// the mismatch, move on and check the next datastore
			continue
		default:
			// mismatched pairing; a VM datastore is not in the list
			// of valid datastores for its current host and is not in
			// the ignore list
			datastoreMismatches = append(datastoreMismatches, ds)
		}

	}

	// return any mismatches found, note that no lookup errors occurred
	return datastoreMismatches, nil

}

// NumDatastores returns the total number of datastores in the index for all
// recorded hosts.
func (vmtmp VMToMismatchedPairing) NumDatastores() int {
	var numDatastores int
	for k := range vmtmp {
		numDatastores += len(vmtmp[k].Datastores)
	}
	return numDatastores
}

// NumHosts returns the total number of records hosts in the index.
func (vmtmp VMToMismatchedPairing) NumHosts() int {
	return len(vmtmp)
}

// H2D2VMsOneLineCheckSummary is used to generate a one-line Nagios service
// check results summary. This is the line most prominent in notifications.
func H2D2VMsOneLineCheckSummary(
	stateLabel string,
	vmsFilterResults VMsFilterResults,
	vmDatastoresPairingIssues VMToMismatchedPairing,
) string {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute H2D2VMsOneLineCheckSummary func.\n",
			time.Since(funcTimeStart),
		)
	}()

	switch {
	case len(vmDatastoresPairingIssues) > 0:
		return fmt.Sprintf(
			"%s: %d mismatched Host/Datastore/VM pairings detected (evaluated %d VMs, %d Resource Pools)",
			stateLabel,
			len(vmDatastoresPairingIssues),
			vmsFilterResults.NumVMsAfterFiltering(),
			vmsFilterResults.NumRPsAfterFiltering(),
		)

	default:

		return fmt.Sprintf(
			"%s: No mismatched Host/Datastore/VM pairings detected (evaluated %d VMs, %d Resource Pools)",
			stateLabel,
			vmsFilterResults.NumVMsAfterFiltering(),
			vmsFilterResults.NumRPsAfterFiltering(),
		)

	}
}

// H2D2VMsReport generates a summary of host/datastore/vms pairings along with
// additional details intended to aid in troubleshooting check results at a
// glance. This information is provided for use with the Long Service Output
// field commonly displayed on the detailed service check results display in
// the web UI or in the body of many notifications.
func H2D2VMsReport(
	c *vim25.Client,
	h2dIdx HostToDatastoreIndex,
	vmsFilterOptions VMsFilterOptions,
	vmsFilterResults VMsFilterResults,
	vmDatastoresPairingIssues VMToMismatchedPairing,
	ignoreMissingCA bool,
	ignoredDatastores []string,
	datastoreCAPrefixSeparator string,
	hostCAPrefixSeparator string,
	datastoreCAName string,
	hostCAName string,
) string {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute HostToDatastoresToVMsReport func.\n",
			time.Since(funcTimeStart),
		)
	}()

	// Build lists of the objects that are missing requested Custom Attribute
	var datastoresMissingCA []string
	var hostsMissingCA []string

	for _, hostDSPairing := range h2dIdx {
		if hostDSPairing.Host.CustomAttribute.Value == CustomAttributeValNotSet {
			hostsMissingCA = append(hostsMissingCA, hostDSPairing.Host.Name)
		}

		for _, ds := range hostDSPairing.Datastores {
			if ds.CustomAttribute.Value == CustomAttributeValNotSet {
				datastoresMissingCA = append(datastoresMissingCA, ds.Name)
			}
		}
	}

	var report strings.Builder

	// if we have more than one hardware version in the index, we have at
	// least one outdated version to report

	switch {

	case len(vmDatastoresPairingIssues) > 0:

		_, _ = fmt.Fprintf(
			&report,
			"Mismatched Hosts / Datastores / Virtual Machines:%s%s",
			nagios.CheckOutputEOL,
			nagios.CheckOutputEOL,
		)

		// build list of VM names for ordered index access
		vmNames := make([]string, 0, len(vmDatastoresPairingIssues))
		for vmName := range vmDatastoresPairingIssues {
			vmNames = append(vmNames, vmName)
		}

		// order the list of VM names for ordered index access
		sort.Slice(vmNames, func(i, j int) bool {
			return strings.ToLower(vmNames[i]) < strings.ToLower(vmNames[j])
		})

		// sort datastore names also (for the associated VM)
		for key := range vmDatastoresPairingIssues {
			// prevent "using the variable on range scope in function literal"
			// linting error
			vmName := key
			sort.Slice(vmDatastoresPairingIssues[key].Datastores, func(i, j int) bool {
				return strings.ToLower(
					vmDatastoresPairingIssues[vmName].Datastores[i].Name,
				) < strings.ToLower(
					vmDatastoresPairingIssues[vmName].Datastores[j].Name,
				)
			})
		}

		for _, vmName := range vmNames {
			var dsNamesWithCA strings.Builder
			for i, ds := range vmDatastoresPairingIssues[vmName].Datastores {
				_, _ = fmt.Fprintf(&dsNamesWithCA, "%q (%s)", ds.Name, ds.CustomAttribute.Value)
				if i != len(vmDatastoresPairingIssues[vmName].Datastores)-1 {
					dsNamesWithCA.WriteString(", ")
				}
			}

			_, _ = fmt.Fprintf(
				&report,
				"* %s: [Host: %q (%s), Datastores: %s]%s",
				vmName,
				vmDatastoresPairingIssues[vmName].Host.Name,
				vmDatastoresPairingIssues[vmName].Host.CustomAttribute.Value,
				dsNamesWithCA.String(),
				nagios.CheckOutputEOL,
			)
		}

		_, _ = fmt.Fprint(&report, nagios.CheckOutputEOL)

	default:

		// homogenous

		_, _ = fmt.Fprintf(
			&report,
			"No mismatched Host/Datastore/VM pairings detected.%s%s",
			nagios.CheckOutputEOL,
			nagios.CheckOutputEOL,
		)

	}

	_, _ = fmt.Fprintf(
		&report,
		"%s---%s%s",
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
	)

	switch {
	case ignoreMissingCA:
		_, _ = fmt.Fprintf(
			&report,
			"* As requested, Hosts and Datastores with missing Custom Attribute are ignored [Host: %q, Datastore: %q]%s",
			hostCAName,
			datastoreCAName,
			nagios.CheckOutputEOL,
		)

	default:
		_, _ = fmt.Fprintf(
			&report,
			"* As requested, Hosts and Datastores with missing Custom Attribute is a fatal condition [Host: %q, Datastore: %q]%s",
			hostCAName,
			datastoreCAName,
			nagios.CheckOutputEOL,
		)
	}

	switch {
	case len(hostsMissingCA) > 0:

		_, _ = fmt.Fprintf(
			&report,
			"Hosts missing Custom Attribute %q: %s%s",
			hostCAName,
			nagios.CheckOutputEOL,
			nagios.CheckOutputEOL,
		)

		for _, hostName := range hostsMissingCA {
			_, _ = fmt.Fprintf(
				&report,
				"* %s%s",
				hostName,
				nagios.CheckOutputEOL,
			)
		}

	case len(datastoresMissingCA) > 0:

		_, _ = fmt.Fprintf(
			&report,
			"Datastores missing Custom Attribute %q: %s%s",
			hostCAName,
			nagios.CheckOutputEOL,
			nagios.CheckOutputEOL,
		)

		for _, dsName := range datastoresMissingCA {
			_, _ = fmt.Fprintf(
				&report,
				"* %s%s",
				dsName,
				nagios.CheckOutputEOL,
			)
		}
	default:
		_, _ = fmt.Fprintf(
			&report,
			"* No Hosts or Datastores are missing specified Custom Attribute%s",
			nagios.CheckOutputEOL,
		)

	}

	if hostCAPrefixSeparator != "" || datastoreCAPrefixSeparator != "" {
		_, _ = fmt.Fprintf(
			&report,
			"* Custom Attribute Prefix Separator: [Host: %q, Datastore: %q]%s",
			hostCAPrefixSeparator,
			datastoreCAPrefixSeparator,
			nagios.CheckOutputEOL,
		)
	}

	vmFilterResultsReportTrailer(
		&report,
		c,
		vmsFilterOptions,
		vmsFilterResults,
		false,
	)

	_, _ = fmt.Fprintf(
		&report,
		"* Specified Datastores to exclude (%d): [%v]%s",
		len(ignoredDatastores),
		strings.Join(ignoredDatastores, ", "),
		nagios.CheckOutputEOL,
	)

	return report.String()
}
