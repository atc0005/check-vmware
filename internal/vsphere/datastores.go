// Copyright 2021 Adam Chalkley
//
// https://github.com/atc0005/check-vmware
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package vsphere

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/atc0005/check-vmware/internal/textutils"
	"github.com/atc0005/go-nagios"
	"github.com/vmware/govmomi/units"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

// ErrDatastoreUsageThresholdCrossed indicates that specified
// resource pools have exceeded a given threshold
var ErrDatastoreUsageThresholdCrossed = errors.New("datastore usage exceeds specified threshold")

// DatastoreIDToNameIndex maps a Datastore's ID value to its name.
type DatastoreIDToNameIndex map[string]string

// DatastoreVMs provides an overview of all (visible) VirtualMachines residing
// on a specific Datastore.
type DatastoreVMs []DatastoreVM

// DatastoreVM is a summary of details for a VirtualMachine found on a
// specific datastore.
type DatastoreVM struct {

	// Name is the display name of the VirtualMachine.
	Name string

	// VMSize is the human readable or formatted size of the VirtualMachine.
	VMSize string

	// DatastoreUsage is the human readable or formatted percentage of the
	// Datastore space consumed by this VirtualMachine.
	DatastoreUsage string

	// PowerState tracks the current power state for a VirtualMachine.
	PowerState types.VirtualMachinePowerState

	// DatastoreMOID is the MOID or MoRef ID for the Datastore where this
	// VirtualMachine resides.
	DatastoreMOID types.ManagedObjectReference
}

// DatastoreUsageSummary tracks usage details for a specific Datastore.
type DatastoreUsageSummary struct {
	Datastore               mo.Datastore
	StorageRemainingPercent float64
	StorageUsedPercent      float64
	StorageTotal            int64
	StorageUsed             int64
	StorageRemaining        int64
	CriticalThreshold       int
	WarningThreshold        int
	VMs                     DatastoreVMs
}

// DatastoreVMsSummary evaluates provided Datastore and collection of
// VirtualMachines and provides a basic human readable / formatted summary of
// VirtualMachine details.
func DatastoreVMsSummary(ds mo.Datastore, vms []mo.VirtualMachine) DatastoreVMs {

	datastoreVMs := make(DatastoreVMs, 0, len(vms))

	for _, vm := range vms {

		var vmStorageUsed int64
		for _, usage := range vm.Storage.PerDatastoreUsage {
			if usage.Datastore == ds.Reference() {
				vmStorageUsed += usage.Committed + usage.Uncommitted
			}
		}

		vmPercentOfDSUsed := float64(vmStorageUsed) / float64(ds.Summary.Capacity) * 100
		dsVM := DatastoreVM{
			Name:           vm.Name,
			VMSize:         units.ByteSize(vmStorageUsed).String(),
			DatastoreUsage: fmt.Sprintf("%2.2f%%", vmPercentOfDSUsed),
			PowerState:     vm.Runtime.PowerState,
		}

		datastoreVMs = append(datastoreVMs, dsVM)

	}

	return datastoreVMs

}

// NewDatastoreUsageSummary receives a Datastore and generates summary
// information used to determine if usage levels have crossed user-specified
// thresholds.
// func NewDatastoreUsageSummary(ds mo.Datastore, dsVMs []mo.VirtualMachine, criticalThreshold int, warningThreshold int) DatastoreUsageSummary {
func NewDatastoreUsageSummary(
	ctx context.Context,
	c *vim25.Client,
	ds mo.Datastore,
	criticalThreshold int,
	warningThreshold int,
) (DatastoreUsageSummary, error) {

	storageRemainingPercentage := float64(ds.Summary.FreeSpace) / float64(ds.Summary.Capacity) * 100
	storageUsedPercentage := 100 - storageRemainingPercentage
	storageRemaining := ds.Summary.FreeSpace
	storageTotal := ds.Summary.Capacity
	storageUsed := storageTotal - storageRemaining

	dsVMs, err := GetVMsFromDatastore(ctx, c, ds, true)
	if err != nil {
		return DatastoreUsageSummary{}, err
	}

	dsUsage := DatastoreUsageSummary{
		Datastore:               ds,
		VMs:                     DatastoreVMsSummary(ds, dsVMs),
		StorageRemainingPercent: storageRemainingPercentage,
		StorageUsedPercent:      storageUsedPercentage,
		StorageTotal:            storageTotal,
		StorageUsed:             storageUsed,
		StorageRemaining:        storageRemaining,
		CriticalThreshold:       criticalThreshold,
		WarningThreshold:        warningThreshold,
	}

	return dsUsage, nil

}

// IsWarningState indicates whether Datastore usage has crossed the WARNING
// level threshold.
func (dus DatastoreUsageSummary) IsWarningState() bool {
	return dus.StorageUsedPercent < float64(dus.CriticalThreshold) &&
		dus.StorageUsedPercent > float64(dus.WarningThreshold)
}

// IsCriticalState indicates whether Datastore usage has crossed the CRITICAL
// level threshold.
func (dus DatastoreUsageSummary) IsCriticalState() bool {
	return dus.StorageUsedPercent > float64(dus.CriticalThreshold)
}

// NumVMsPoweredOn indicates how many VirtualMachines on a specific Datastore
// are powered on.
func (dsVMs DatastoreVMs) NumVMsPoweredOn() int {

	var numOn int
	for _, vm := range dsVMs {
		if vm.PowerState == types.VirtualMachinePowerStatePoweredOn {
			numOn++
		}
	}

	return numOn
}

// NumVMsPoweredOff indicates how many VirtualMachines on a specific Datastore
// are powered off OR suspended.
func (dsVMs DatastoreVMs) NumVMsPoweredOff() int {
	return len(dsVMs) - dsVMs.NumVMsPoweredOn()
}

// GetDatastores accepts a context, a connected client and a boolean value
// indicating whether a subset of properties per Datastore are retrieved. A
// collection of Datastores with requested properties is returned. If
// requested, a subset of all available properties will be retrieved (faster)
// instead of recursively fetching all properties (about 2x as slow).
func GetDatastores(ctx context.Context, c *vim25.Client, propsSubset bool) ([]mo.Datastore, error) {

	funcTimeStart := time.Now()

	// declare this early so that we can grab a pointer to it in order to
	// access the entries later
	var dss []mo.Datastore

	defer func(dss *[]mo.Datastore) {
		logger.Printf(
			"It took %v to execute GetDatastores func (and retrieve %d Datastores).\n",
			time.Since(funcTimeStart),
			len(*dss),
		)
	}(&dss)

	err := getObjects(ctx, c, &dss, c.ServiceContent.RootFolder, propsSubset)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve Datastores: %w", err)
	}

	sort.Slice(dss, func(i, j int) bool {
		return strings.ToLower(dss[i].Name) < strings.ToLower(dss[j].Name)
	})

	return dss, nil
}

// GetDatastoreByName accepts the name of a datastore, the name of a
// datacenter and a boolean value indicating whether only a subset of
// properties for the Datastore should be returned. If requested, a subset of
// all available properties will be retrieved (faster) instead of recursively
// fetching all properties (about 2x as slow). If the datacenter name is an
// empty string then the default datacenter will be used.
func GetDatastoreByName(ctx context.Context, c *vim25.Client, dsName string, datacenter string, propsSubset bool) (mo.Datastore, error) {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute GetDatastoreByName func.\n",
			time.Since(funcTimeStart),
		)
	}()

	var datastore mo.Datastore

	err := getObjectByName(ctx, c, &datastore, dsName, datacenter, propsSubset)
	if err != nil {
		return mo.Datastore{}, err
	}

	return datastore, nil

}

// FilterDatastoresByName accepts a collection of Datastores and a Datastore
// name to filter against. An error is returned if the list of Datastores is
// empty or if a match was not found. The matching Datastore is returned along
// with the number of Datastores that were excluded.
func FilterDatastoresByName(dss []mo.Datastore, dsName string) (mo.Datastore, int, error) {

	funcTimeStart := time.Now()

	// If error condition, no exclusions are made
	numExcluded := 0

	defer func() {
		logger.Printf(
			"It took %v to execute FilterDatastoresByName func.\n",
			time.Since(funcTimeStart),
		)
	}()

	if len(dss) == 0 {
		return mo.Datastore{}, numExcluded, fmt.Errorf("received empty list of datastores to filter by name")
	}

	for _, ds := range dss {
		if ds.Name == dsName {
			// we are excluding everything but the single name value match
			numExcluded = len(dss) - 1
			return ds, numExcluded, nil
		}
	}

	return mo.Datastore{}, numExcluded, fmt.Errorf(
		"error: failed to retrieve Datastore using provided name %q",
		dsName,
	)

}

// FilterDatastoresByID receives a collection of Datastores and a Datastore ID
// to filter against. An error is returned if the list of Datastores is empty
// or if a match was not found. The matching Datastore is returned along with
// the number of Datastores that were excluded.
func FilterDatastoresByID(dss []mo.Datastore, dsID string) (mo.Datastore, int, error) {

	funcTimeStart := time.Now()

	// If error condition, no exclusions are made
	numExcluded := 0

	defer func() {
		logger.Printf(
			"It took %v to execute FilterDatastoresByID func.\n",
			time.Since(funcTimeStart),
		)
	}()

	if len(dss) == 0 {
		return mo.Datastore{}, numExcluded, fmt.Errorf("received empty list of datastores to filter by ID")
	}

	for _, ds := range dss {
		// return match, if available
		// TODO: Refactor, use abstract type here
		// ds.GetManagedEntity().Reference().Value
		if ds.Summary.Datastore.Value == dsID {
			// we are excluding everything but the single name value match
			numExcluded = len(dss) - 1
			return ds, numExcluded, nil
		}
	}

	return mo.Datastore{}, numExcluded, fmt.Errorf(
		"error: failed to retrieve Datastore using provided id %q",
		dsID,
	)

}

// DatastoreIDsToNames returns a list of matching Datastore names for the
// provided list of Managed Object References for Datastores.
func DatastoreIDsToNames(dsRefs []types.ManagedObjectReference, dss []mo.Datastore) []string {

	dsNames := make([]string, 0, len(dsRefs))
	dsIDs := make([]string, 0, len(dsRefs))

	for _, dsRef := range dsRefs {
		dsIDs = append(dsIDs, dsRef.Value)
	}

	for _, ds := range dss {
		if textutils.InList(ds.Summary.Datastore.Value, dsIDs, true) {
			dsNames = append(dsNames, ds.Name)
		}
	}

	return dsNames

}

// DatastoreUsageOneLineCheckSummary is used to generate a one-line Nagios
// service check results summary. This is the line most prominent in
// notifications.
func DatastoreUsageOneLineCheckSummary(
	stateLabel string,
	dsUsageSummary DatastoreUsageSummary,
) string {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute DatastoreUsageOneLineCheckSummary func.\n",
			time.Since(funcTimeStart),
		)
	}()

	return fmt.Sprintf(
		"%s: Datastore %s usage (%d VMs) is %.2f%% of %s with %s remaining [WARNING: %d%% , CRITICAL: %d%%]",
		stateLabel,
		dsUsageSummary.Datastore.Name,
		len(dsUsageSummary.VMs),
		dsUsageSummary.StorageUsedPercent,
		units.ByteSize(dsUsageSummary.StorageTotal),
		units.ByteSize(dsUsageSummary.StorageRemaining),
		dsUsageSummary.WarningThreshold,
		dsUsageSummary.CriticalThreshold,
	)

}

// DatastoreUsageReport generates a summary of Datastore usage along with
// various verbose details intended to aid in troubleshooting check results at
// a glance. This information is provided for use with the Long Service Output
// field commonly displayed on the detailed service check results display in
// the web UI or in the body of many notifications.
func DatastoreUsageReport(
	c *vim25.Client,
	dsUsageSummary DatastoreUsageSummary,
) string {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute DatastoreUsageReport func.\n",
			time.Since(funcTimeStart),
		)
	}()

	var report strings.Builder

	fmt.Fprintf(
		&report,
		"Datastore Summary:%s%s"+
			"* Name: %s%s"+
			"* Used: %v (%.2f%%)%s"+
			"* Remaining: %v (%.2f%%)%s"+
			"* VMs: %v %s%s",
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
		dsUsageSummary.Datastore.Name,
		nagios.CheckOutputEOL,
		units.ByteSize(dsUsageSummary.StorageUsed),
		dsUsageSummary.StorageUsedPercent,
		nagios.CheckOutputEOL,
		units.ByteSize(dsUsageSummary.StorageRemaining),
		dsUsageSummary.StorageRemainingPercent,
		nagios.CheckOutputEOL,
		len(dsUsageSummary.VMs),
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
	)

	printVMSummary := func(powerState types.VirtualMachinePowerState) {

		// Skip efforts to list VM summary details if there is nothing to show.
		if len(dsUsageSummary.VMs) == 0 {
			return
		}

		var powerStateVMs int
		switch powerState {
		case types.VirtualMachinePowerStatePoweredOn:
			powerStateVMs = dsUsageSummary.VMs.NumVMsPoweredOn()
		default:
			powerStateVMs = dsUsageSummary.VMs.NumVMsPoweredOff()
		}

		fmt.Fprintf(
			&report,
			"%d %s VMs on datastore:%s%s",
			powerStateVMs,
			powerState,
			nagios.CheckOutputEOL,
			nagios.CheckOutputEOL,
		)

		for _, vm := range dsUsageSummary.VMs {
			if vm.PowerState == powerState {
				fmt.Fprintf(
					&report,
					"* %s [Size: %s, Datastore Usage: %s]%s",
					vm.Name,
					vm.VMSize,
					vm.DatastoreUsage,
					nagios.CheckOutputEOL,
				)
			}
		}

		fmt.Fprintf(&report, nagios.CheckOutputEOL)
	}

	printVMSummary(types.VirtualMachinePowerStatePoweredOn)

	printVMSummary(types.VirtualMachinePowerStatePoweredOff)

	fmt.Fprintf(
		&report,
		"%s---%s%s",
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		&report,
		"* vSphere environment: %s%s",
		c.URL().String(),
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		&report,
		"* Plugin User Agent: %s%s",
		c.Client.UserAgent,
		nagios.CheckOutputEOL,
	)

	return report.String()
}
