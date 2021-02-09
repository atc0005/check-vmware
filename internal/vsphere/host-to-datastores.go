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
	"os"
	"sort"
	"strings"
	"time"

	"github.com/atc0005/check-vmware/internal/textutils"
	"github.com/atc0005/go-nagios"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
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
	DatastoreID string
	Err         error
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

// PairingCustomAttribute represents the key/value Custom Attribute pair used
// to relate specific hosts with specific datastores. Most often this takes
// the form of a "Location" or "Datacenter" field to indicate which Datastore
// a VirtualMachine should reside on when running on a specific HostSystem.
type PairingCustomAttribute struct {
	Name  string
	Value string
}

// DatastoreWithCA wraps the vSphere Datastore managed object type with a
// Custom Attribute key/value pair. This Custom Attribute is intended to link
// this datastore to a specific ESXi host.
type DatastoreWithCA struct {
	mo.Datastore
	CustomAttribute PairingCustomAttribute
}

// HostWithCA wraps the vSphere HostSystem managed object type with a Custom
// Attribute key/value pair. This Custom Attribute is intended to link this
// host to one or more datastores.
type HostWithCA struct {
	mo.HostSystem
	CustomAttribute PairingCustomAttribute
}

// HostDatastoresPairing collects Host and Datastores pairings based on shared
// Custom Attribute name and value (literal) or prefix (if user-specified).
// This is intended to "pair" hosts and datastores within an environment that
// are known to work well together.
type HostDatastoresPairing struct {
	Host       HostWithCA
	Datastores []DatastoreWithCA
}

// VMHostDatastoresPairing collects HostSystem, VirtualMachine and Datastore
// name pairings.
type VMHostDatastoresPairing struct {
	HostName       string
	DatastoreNames []string
	// VMName         string
}

// HostToDatastoreIndex indexes HostDatastorePairings based on host id values.
type HostToDatastoreIndex map[string]HostDatastoresPairing

// VMToMismatchedDatastoreNames indexes VirtualMachine name to
// VMHostDatastoresPairing type. This index reflects mismatched Datastore
// names based on the current host for the VirtualMachine.
type VMToMismatchedDatastoreNames map[string]VMHostDatastoresPairing

// NewHostToDatastoreIndex receives a collection of hosts and datastores
// wrapped with user-specified Custom Attributes, prefix separators and a
// boolean flag indicating whether prefix matching will be used. The resulting
// HostToDatastoreIndex is returned if no errors occur, otherwise nil and the
// error.
func NewHostToDatastoreIndex(
	hosts []HostWithCA,
	datastores []DatastoreWithCA,
	usingPrefixes bool,
	hostCASep string,
	datastoreCASep string,
) (HostToDatastoreIndex, error) {

	h2dIdx := make(HostToDatastoreIndex)

	for _, host := range hosts {

		hostID := host.Summary.Host.Value

		for _, datastore := range datastores {

			if usingPrefixes {
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
				}
			}

			// not using prefixes, so literal values
			if strings.EqualFold(datastore.CustomAttribute.Value, host.CustomAttribute.Value) {
				h2dIdx[hostID] = HostDatastoresPairing{
					Host:       host,
					Datastores: append(h2dIdx[hostID].Datastores, datastore),
				}
			}
		}
	}

	if len(h2dIdx) == 0 {
		return nil, ErrHostDatastorePairingFailed
	}

	return h2dIdx, nil

}

// DatastoreNames returns a list of all Datastore names in the index.
func (hdi HostToDatastoreIndex) DatastoreNames() []string {

	var dsNames []string
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

	dsIdx := make(DatastoreIDToNameIndex)
	for hostID := range hdi {
		for _, ds := range hdi[hostID].Datastores {
			dsIdx[ds.Summary.Datastore.Value] = ds.Name
		}
	}

	return dsIdx

}

// IsDatastoreIDInIndex indicates whether a provided Datastore ID is in the
// index.
func (hdi HostToDatastoreIndex) IsDatastoreIDInIndex(dsID string) bool {

	for hostID := range hdi {
		for _, ds := range hdi[hostID].Datastores {
			if strings.EqualFold(dsID, ds.Summary.Datastore.Value) {
				return true
			}
		}
	}

	return false

}

// DatastoreIDToName returns the name associated with a Datastore ID. An error
// is returned if the name could not be retrieved from the index.
func (hdi HostToDatastoreIndex) DatastoreIDToName(dsID string) (string, error) {

	for hostID := range hdi {
		for _, ds := range hdi[hostID].Datastores {
			if ds.Summary.Datastore.Value == dsID {
				return ds.Name, nil
			}
		}
	}

	return "", &ErrHostDatastoreIdxIDToNameLookupFailed{
		DatastoreID: dsID,
		Err:         errors.New("datastore ID not found"),
	}
}

// ValidateVirtualMachinePairings receives a VirtualMachine ID, a collection
// of Datastore IDs associated with the VM and an optional list of Datastore
// names to ignore. A list of mismatched datastores is returned along with any
// errors that may occur.
func (hdi HostToDatastoreIndex) ValidateVirtualMachinePairings(
	vmHostID string,
	allDatastores []mo.Datastore,
	vmDatastoreRefs []types.ManagedObjectReference,
	dsNamesToIgnore []string,
) ([]string, error) {

	// fmt.Println("All datastores length:", len(allDatastores))
	// fmt.Println("vmDatastoreRefs length:", len(vmDatastoreRefs))
	// fmt.Println("dsNamesToIgnore length:", len(dsNamesToIgnore))
	// fmt.Println("vmHostID:", vmHostID)

	// defer func() {
	// 	if err := recover(); err != nil {
	// 		fmt.Println(err)
	// 		panic(err)
	// 	}
	// }()

	// assert that every datastore for the VM is in the list of datastores for
	// the host

	vmDatastoreIDs := make([]string, 0, len(vmDatastoreRefs))
	hostDatastoreIDs := make([]string, 0, len(hdi[vmHostID].Datastores))

	for _, vmDSRef := range vmDatastoreRefs {
		vmDatastoreIDs = append(vmDatastoreIDs, vmDSRef.Value)
	}

	for _, hostPairedDS := range hdi[vmHostID].Datastores {
		hostDatastoreIDs = append(hostDatastoreIDs, hostPairedDS.Summary.Datastore.Value)
	}

	var datastoreMismatches []string

	// assert that every datastore ID associated with the VM is within the
	// list of datastores associated with the current host for the VM.
	for _, vmDatastoreID := range vmDatastoreIDs {

		if !textutils.InList(vmDatastoreID, hostDatastoreIDs, true) {

			// lookup errors abort the validation process, unless ...
			dsName, lookupErr := hdi.DatastoreIDToName(vmDatastoreID)
			if lookupErr != nil {

				// lookup could have failed if the sole datastore for the
				// VM is in the ignored list; double-check that
				// possibility before reporting the lookup failure.
				var dsIDLookupErr *ErrHostDatastoreIdxIDToNameLookupFailed
				if errors.As(lookupErr, &dsIDLookupErr) {

					// TODO: Switch this off after sufficient testing has
					// been completed. For now, explicitly send to stderr
					// to keep Nagios from routing it to notifications or
					// the web UI
					fmt.Fprintf(
						os.Stderr,
						"Initial lookup failed for %s\n",
						vmDatastoreID,
					)

					dsID := dsIDLookupErr.DatastoreID
					datastore, filterErr := FilterDatastoreByID(allDatastores, dsID)
					if filterErr != nil {

						// This is our second attempt to lookup the
						// datastore using the datastore id. The first
						// failure is because the datastore isn't in our
						// host-to-datastore index, this second because it
						// could not be located in the full datastores
						// list from the vSphere inventory.
						return datastoreMismatches, fmt.Errorf(
							"second lookup attempt unsuccessful; "+
								"failed to locate datastore ID in "+
								"index or full datastores list: %w",
							filterErr,
						)
					}

					if textutils.InList(datastore.Name, dsNamesToIgnore, true) {

						// TODO: Switch this off after sufficient testing
						// has been completed. For now, explicitly send to
						// stderr to keep Nagios from routing it to
						// notifications or the web UI
						fmt.Fprintf(
							os.Stderr,
							"Second lookup successful; name: %s id: %s\n",
							vmDatastoreID,
							datastore.Name,
						)
						continue
					}

				}

				// Lookup failure occurred for some other reason.
				return datastoreMismatches, lookupErr

			}

			switch {
			case textutils.InList(dsName, dsNamesToIgnore, true):
				// if datastore name is in the ignore list, don't report
				// the mismatch, move on and check the next datastore
				continue
			default:
				// mismatched pairing; a VM datastore is not in the list
				// of valid datastores for its current host and is not in
				// the ignore list
				datastoreMismatches = append(datastoreMismatches, dsName)
			}
		}
	}

	// return any mismatches found, note that no lookup errors occurred
	return datastoreMismatches, nil

}

// H2D2VMsOneLineCheckSummary is used to generate a one-line Nagios service
// check results summary. This is the line most prominent in notifications.
func H2D2VMsOneLineCheckSummary(
	stateLabel string,
	evaluatedVMs []mo.VirtualMachine,
	vmDatastoresPairingIssues VMToMismatchedDatastoreNames,
	rps []mo.ResourcePool,
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
			len(evaluatedVMs),
			len(rps),
		)

	default:

		return fmt.Sprintf(
			"%s: No mismatched Host/Datastore/VM pairings detected (evaluated %d VMs, %d Resource Pools)",
			stateLabel,
			len(evaluatedVMs),
			len(rps),
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
	allVMs []mo.VirtualMachine,
	evaluatedVMs []mo.VirtualMachine,
	vmDatastoresPairingIssues VMToMismatchedDatastoreNames,
	vmsToExclude []string,
	evalPoweredOffVMs bool,
	includeRPs []string,
	excludeRPs []string,
	rps []mo.ResourcePool,
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

	rpNames := make([]string, len(rps))
	for i := range rps {
		rpNames[i] = rps[i].Name
	}

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

		fmt.Fprintf(
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
		// for key := range vmDatastoresPairingIssues {
		// 	// prevent "using the variable on range scope in function literal"
		// 	// linting error
		// 	key := key
		// 	sort.Slice(vmDatastoresPairingIssues[key], func(i, j int) bool {
		// 		return strings.ToLower(
		// 			vmDatastoresPairingIssues[key][i],
		// 		) < strings.ToLower(
		// 			vmDatastoresPairingIssues[key][j],
		// 		)
		// 	})
		// }
		// TODO: Is the sorted order acceptable?
		for key := range vmDatastoresPairingIssues {
			sort.Strings(vmDatastoresPairingIssues[key].DatastoreNames)
		}

		for _, vmName := range vmNames {
			fmt.Fprintf(
				&report,
				"* %s: [%s, %s]%s",
				vmName,
				vmDatastoresPairingIssues[vmName].HostName,
				strings.Join(vmDatastoresPairingIssues[vmName].DatastoreNames, ", "),
				nagios.CheckOutputEOL,
			)
		}

		fmt.Fprint(&report, nagios.CheckOutputEOL)

	default:

		// homogenous

		fmt.Fprintf(
			&report,
			"No mismatched Host/Datastore/VM pairings detected.%s%s",
			nagios.CheckOutputEOL,
			nagios.CheckOutputEOL,
		)

	}

	fmt.Fprintf(
		&report,
		"%s---%s%s",
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
	)

	switch {
	case ignoreMissingCA:
		fmt.Fprintf(
			&report,
			"* As requested, Hosts & Datastores with missing Custom Attribute are ignored [Host: %q, Datastore: %q]%s",
			hostCAName,
			datastoreCAName,
			nagios.CheckOutputEOL,
		)

	default:
		fmt.Fprintf(
			&report,
			"* As requested, Hosts & Datastores with missing Custom Attribute is a fatal condition [Host: %q, Datastore: %q]%s",
			hostCAName,
			datastoreCAName,
			nagios.CheckOutputEOL,
		)
	}

	switch {
	case len(hostsMissingCA) > 0:

		fmt.Fprintf(
			&report,
			"Hosts missing Custom Attribute %q: %s%s",
			hostCAName,
			nagios.CheckOutputEOL,
			nagios.CheckOutputEOL,
		)

		for _, hostName := range hostsMissingCA {
			fmt.Fprintf(
				&report,
				"* %s%s",
				hostName,
				nagios.CheckOutputEOL,
			)
		}

	case len(datastoresMissingCA) > 0:

		fmt.Fprintf(
			&report,
			"Datastores missing Custom Attribute %q: %s%s",
			hostCAName,
			nagios.CheckOutputEOL,
			nagios.CheckOutputEOL,
		)

		for _, dsName := range datastoresMissingCA {
			fmt.Fprintf(
				&report,
				"* %s%s",
				dsName,
				nagios.CheckOutputEOL,
			)
		}
	default:
		fmt.Fprintf(
			&report,
			"* No Hosts or Datastores are missing specified Custom Attribute%s",
			nagios.CheckOutputEOL,
		)

	}

	if hostCAPrefixSeparator != "" || datastoreCAPrefixSeparator != "" {
		fmt.Fprintf(
			&report,
			"* Custom Attribute Prefix Separator: [Host: %q, Datastore: %q]%s",
			hostCAPrefixSeparator,
			datastoreCAPrefixSeparator,
			nagios.CheckOutputEOL,
		)
	}

	fmt.Fprintf(
		&report,
		"* vSphere environment: %s%s",
		c.URL().String(),
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		&report,
		"* VMs (evaluated: %d, total: %d)%s",
		len(evaluatedVMs),
		len(allVMs),
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		&report,
		"* Powered off VMs evaluated: %t%s",
		evalPoweredOffVMs,
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		&report,
		"* Specified VMs to exclude (%d): [%v]%s",
		len(vmsToExclude),
		strings.Join(vmsToExclude, ", "),
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		&report,
		"* Specified Datastores to exclude (%d): [%v]%s",
		len(ignoredDatastores),
		strings.Join(ignoredDatastores, ", "),
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		&report,
		"* Specified Resource Pools to explicitly include (%d): [%v]%s",
		len(includeRPs),
		strings.Join(includeRPs, ", "),
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		&report,
		"* Specified Resource Pools to explicitly exclude (%d): [%v]%s",
		len(excludeRPs),
		strings.Join(excludeRPs, ", "),
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		&report,
		"* Resource Pools evaluated (%d): [%v]%s",
		len(rpNames),
		strings.Join(rpNames, ", "),
		nagios.CheckOutputEOL,
	)

	return report.String()
}
