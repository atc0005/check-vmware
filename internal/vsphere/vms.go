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
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

// ErrVirtualMachinePowerCycleUptimeThresholdCrossed indicates that specified
// Virtual Machine power cycle thresholds have been exceeded.
var ErrVirtualMachinePowerCycleUptimeThresholdCrossed = errors.New("power cycle uptime exceeds specified threshold")

// ErrVirtualMachineDiskConsolidationNeeded indicates that disk consolidation
// is needed for one or more Virtual Machines.
var ErrVirtualMachineDiskConsolidationNeeded = errors.New("disk consolidation needed")

// ErrVirtualMachineInteractiveResponseNeeded indicates that an interactive
// response is needed for one or more Virtual Machines.
var ErrVirtualMachineInteractiveResponseNeeded = errors.New("interactive response needed")

// VirtualMachinePowerCycleUptimeStatus tracks VirtualMachines with power
// cycle uptimes that exceed specified thresholds while retaining a list of
// the VirtualMachines that have yet to exceed thresholds.
type VirtualMachinePowerCycleUptimeStatus struct {
	VMsCritical       []mo.VirtualMachine
	VMsWarning        []mo.VirtualMachine
	VMsOK             []mo.VirtualMachine
	WarningThreshold  int
	CriticalThreshold int
}

// VMNames returns a list of sorted VirtualMachine names which have exceeded
// specified power cycle uptime thresholds. VirtualMachines which have yet to
// exceed specified thresholds are not listed.
func (vpcs VirtualMachinePowerCycleUptimeStatus) VMNames() string {
	vmNames := make([]string, 0, len(vpcs.VMsCritical)+len(vpcs.VMsWarning))

	for _, vm := range vpcs.VMsWarning {
		vmNames = append(vmNames, vm.Name)
	}
	for _, vm := range vpcs.VMsCritical {
		vmNames = append(vmNames, vm.Name)
	}

	sort.Slice(vmNames, func(i, j int) bool {
		return strings.ToLower(vmNames[i]) < strings.ToLower(vmNames[j])
	})

	return strings.Join(vmNames, ", ")
}

// TopTenOK is a helper method that returns at most ten VMs with the highest
// power cycle uptime values that have yet to exceed specified thresholds.
func (vpcs VirtualMachinePowerCycleUptimeStatus) TopTenOK() []mo.VirtualMachine {

	// sort before we sample the VMs so that we only get the ones with highest
	// power cycle uptime
	sort.Slice(vpcs.VMsOK, func(i, j int) bool {
		return vpcs.VMsOK[i].Summary.QuickStats.UptimeSeconds > vpcs.VMsOK[j].Summary.QuickStats.UptimeSeconds
	})

	sampleSize := len(vpcs.VMsOK)
	switch {
	case sampleSize > 10:
		sampleSize = 10
	case sampleSize == 0:
		return []mo.VirtualMachine{}
	}

	topTen := make([]mo.VirtualMachine, 0, sampleSize)
	topTen = append(topTen, vpcs.VMsOK[:sampleSize]...)

	return topTen

}

// BottomTenOK is a helper method that returns at most ten VMs with the lowest
// power cycle uptime values that have yet to exceed specified thresholds.
// Only powered on VMs are considered.
func (vpcs VirtualMachinePowerCycleUptimeStatus) BottomTenOK() []mo.VirtualMachine {

	poweredOnVMs, _ := FilterVMsByPowerState(vpcs.VMsOK, false)

	// sort before we sample the VMs so that we only get the ones with lowest
	// power cycle uptime; require that the VM be powered on in order to sort
	// in the intended order.
	sort.Slice(poweredOnVMs, func(i, j int) bool {
		return poweredOnVMs[i].Summary.QuickStats.UptimeSeconds < poweredOnVMs[j].Summary.QuickStats.UptimeSeconds

	})

	sampleSize := len(poweredOnVMs)
	switch {
	case sampleSize > 10:
		sampleSize = 10
	case sampleSize == 0:
		return []mo.VirtualMachine{}
	}

	bottomTen := make([]mo.VirtualMachine, 0, sampleSize)
	bottomTen = append(bottomTen, poweredOnVMs[:sampleSize]...)

	return bottomTen

}

// GetVMs accepts a context, a connected client and a boolean value indicating
// whether a subset of properties per VirtualMachine are retrieved. If
// requested, a subset of all available properties will be retrieved (faster)
// instead of recursively fetching all properties (about 2x as slow) A
// collection of VirtualMachines with requested properties is returned or nil
// and an error, if one occurs.
func GetVMs(ctx context.Context, c *vim25.Client, propsSubset bool) ([]mo.VirtualMachine, error) {

	funcTimeStart := time.Now()

	// declare this early so that we can grab a pointer to it in order to
	// access the entries later
	// vms := make([]mo.VirtualMachine, 0, 100)
	var vms []mo.VirtualMachine

	defer func(vms *[]mo.VirtualMachine) {
		logger.Printf(
			"It took %v to execute GetVMs func (and retrieve %d VirtualMachines).\n",
			time.Since(funcTimeStart),
			len(*vms),
		)
	}(&vms)

	err := getObjects(ctx, c, &vms, c.ServiceContent.RootFolder, propsSubset)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve VirtualMachines: %w", err)
	}

	sort.Slice(vms, func(i, j int) bool {
		return strings.ToLower(vms[i].Name) < strings.ToLower(vms[j].Name)
	})

	return vms, nil
}

// GetVMsFromContainer receives one or many ManagedEntity values for Folder,
// Datacenter, ComputeResource, ResourcePool, VirtualApp or HostSystem types
// and returns a list of VirtualMachine object references.
//
// The propsSubset boolean value indicates whether a subset of properties per
// VirtualMachine are retrieved. If requested, a subset of all available
// properties will be retrieved (faster) instead of recursively fetching all
// properties (about 2x as slow). A collection of VirtualMachines with
// requested properties is returned or nil and an error, if one occurs.
func GetVMsFromContainer(ctx context.Context, c *vim25.Client, propsSubset bool, objs ...mo.ManagedEntity) ([]mo.VirtualMachine, error) {

	funcTimeStart := time.Now()

	// declare this early so that we can grab a pointer to it in order to
	// access the entries later
	var allVMs []mo.VirtualMachine

	defer func(vms *[]mo.VirtualMachine) {
		logger.Printf(
			"It took %v to execute GetVMsFromContainers func (and retrieve %d VMs).\n",
			time.Since(funcTimeStart),
			len(*vms),
		)
	}(&allVMs)

	for _, obj := range objs {

		var vmsFromContainer []mo.VirtualMachine

		err := getObjects(ctx, c, &vmsFromContainer, obj.Reference(), propsSubset)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to retrieve VirtualMachines from object: %s: %w",
				obj.Name,
				err,
			)
		}

		allVMs = append(allVMs, vmsFromContainer...)

	}

	// remove any potential duplicate entries which could occur if we are
	// evaluating the (default, hidden) 'Resources' Resource Pool
	allVMs = dedupeVMs(allVMs)

	sort.Slice(allVMs, func(i, j int) bool {
		return strings.ToLower(allVMs[i].Name) < strings.ToLower(allVMs[j].Name)
	})

	return allVMs, nil

}

// GetVMsFromDatastore receives a Datastore object reference and returns a
// list of VirtualMachine object references. The propsSubset boolean value
// indicates whether a subset of properties per VirtualMachine are retrieved.
// If requested, a subset of all available properties will be retrieved
// (faster) instead of recursively fetching all properties (about 2x as slow)
// A collection of VirtualMachines with requested properties is returned or
// nil and an error, if one occurs.
func GetVMsFromDatastore(ctx context.Context, c *vim25.Client, ds mo.Datastore, propsSubset bool) ([]mo.VirtualMachine, error) {

	funcTimeStart := time.Now()

	// declare this early so that we can grab a pointer to it in order to
	// access the entries later
	dsVMs := make([]mo.VirtualMachine, len(ds.Vm))

	defer func(vms *[]mo.VirtualMachine) {
		logger.Printf(
			"It took %v to execute GetVMsFromDatastore func (and retrieve %d VMs).\n",
			time.Since(funcTimeStart),
			len(*vms),
		)
	}(&dsVMs)

	var allVMs []mo.VirtualMachine
	err := getObjects(ctx, c, &allVMs, c.ServiceContent.RootFolder, propsSubset)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to retrieve VirtualMachines from Datastore %s: %w",
			ds.Name,
			err,
		)
	}

	for i := range ds.Vm {
		vm, _, err := FilterVMsByID(allVMs, ds.Vm[i].Value)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to retrieve VM for VM ID %s: %w",
				ds.Vm[i].Value,
				err,
			)
		}

		dsVMs[i] = vm
	}

	sort.Slice(dsVMs, func(i, j int) bool {
		return strings.ToLower(dsVMs[i].Name) < strings.ToLower(dsVMs[j].Name)
	})

	return dsVMs, nil

}

// GetVMByName accepts the name of a VirtualMachine, the name of a datacenter
// and a boolean value indicating whether only a subset of properties for the
// VirtualMachine should be returned. If requested, a subset of all available
// properties will be retrieved (faster) instead of recursively fetching all
// properties (about 2x as slow). If the datacenter name is an empty string
// then the default datacenter will be used.
func GetVMByName(ctx context.Context, c *vim25.Client, vmName string, datacenter string, propsSubset bool) (mo.VirtualMachine, error) {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute GetVMByName func.\n",
			time.Since(funcTimeStart),
		)
	}()

	var vm mo.VirtualMachine
	err := getObjectByName(ctx, c, &vm, vmName, datacenter, propsSubset)

	if err != nil {
		return mo.VirtualMachine{}, err
	}

	return vm, nil

}

// FilterVMsByName accepts a collection of VirtualMachines and a
// VirtualMachine name to filter against. An error is returned if the list of
// VirtualMachines is empty or if a match was not found. The matching
// VirtualMachine is returned along with the number of VirtualMachines that
// were excluded.
func FilterVMsByName(vms []mo.VirtualMachine, vmName string) (mo.VirtualMachine, int, error) {

	funcTimeStart := time.Now()

	// If error condition, no VMs are excluded
	numExcluded := 0

	defer func() {
		logger.Printf(
			"It took %v to execute FilterVMsByName func.\n",
			time.Since(funcTimeStart),
		)
	}()

	if len(vms) == 0 {
		return mo.VirtualMachine{}, numExcluded, fmt.Errorf("received empty list of virtual machines to filter by name")
	}

	for _, vm := range vms {
		if vm.Name == vmName {
			// we are excluding everything but the single name value match
			numExcluded = len(vms) - 1
			return vm, numExcluded, nil
		}
	}

	return mo.VirtualMachine{}, numExcluded, fmt.Errorf(
		"error: failed to retrieve VirtualMachine using provided name %q",
		vmName,
	)

}

// FilterVMsByID receives a collection of VirtualMachines and a VirtualMachine
// ID to filter against. An error is returned if the list of VirtualMachines
// is empty or if a match was not found. The matching VirtualMachine is
// returned along with the number of VirtualMachines that were excluded.
func FilterVMsByID(vms []mo.VirtualMachine, vmID string) (mo.VirtualMachine, int, error) {

	funcTimeStart := time.Now()

	// If error condition, no VMs are excluded
	numExcluded := 0

	defer func() {
		logger.Printf(
			"It took %v to execute FilterVMsByID func.\n",
			time.Since(funcTimeStart),
		)
	}()

	if len(vms) == 0 {
		return mo.VirtualMachine{},
			numExcluded,
			fmt.Errorf("received empty list of virtual machines to filter by ID")
	}

	for _, vm := range vms {
		// return match, if available
		if vm.Self.Value == vmID {
			// we are excluding everything but the single ID value match
			numExcluded = len(vms) - 1
			return vm, numExcluded, nil
		}
	}

	return mo.VirtualMachine{}, numExcluded, fmt.Errorf(
		"error: failed to retrieve VirtualMachine using provided ID %q",
		vmID,
	)

}

// ExcludeVMsByName receives a collection of VirtualMachines and a list of VMs
// that should be ignored. A new collection minus ignored VirtualMachines is
// returned along with the number of VMs that were excluded.
//
// If the collection of VirtualMachine is empty, an empty collection is
// returned. If the list of ignored VirtualMachines is empty, the same items
// from the received collection of VirtualMachines is returned. If the list of
// ignored VirtualMachines is greater than the list of received
// VirtualMachines, then only matching VirtualMachines will be excluded and
// any others silently skipped.
func ExcludeVMsByName(allVMs []mo.VirtualMachine, ignoreList []string) ([]mo.VirtualMachine, int) {

	if len(allVMs) == 0 || len(ignoreList) == 0 {
		return allVMs, 0
	}

	vmsToKeep := make([]mo.VirtualMachine, 0, len(allVMs))

	for _, vm := range allVMs {
		if textutils.InList(vm.Name, ignoreList, true) {
			continue
		}
		vmsToKeep = append(vmsToKeep, vm)
	}

	sort.Slice(vmsToKeep, func(i, j int) bool {
		return strings.ToLower(vmsToKeep[i].Name) < strings.ToLower(vmsToKeep[j].Name)
	})

	numExcluded := len(allVMs) - len(vmsToKeep)

	return vmsToKeep, numExcluded

}

// FilterVMsByPowerState accepts a collection of VirtualMachines and a boolean
// value to indicate whether powered off VMs should be included in the
// returned collection. If the collection of provided VirtualMachines is
// empty, an empty collection is returned. The collection is returned along
// with the number of VirtualMachines that were excluded.
func FilterVMsByPowerState(vms []mo.VirtualMachine, includePoweredOff bool) ([]mo.VirtualMachine, int) {

	// setup early so we can reference it from deferred stats output
	filteredVMs := make([]mo.VirtualMachine, 0, len(vms))

	funcTimeStart := time.Now()

	defer func(vms []mo.VirtualMachine, filteredVMs *[]mo.VirtualMachine) {
		logger.Printf(
			"It took %v to execute FilterVMsByPowerState func (for %d VMs, yielding %d VMs)\n",
			time.Since(funcTimeStart),
			len(vms),
			len(*filteredVMs),
		)
	}(vms, &filteredVMs)

	if len(vms) == 0 {
		return vms, 0
	}

	for _, vm := range vms {
		switch {
		case vm.Runtime.PowerState == types.VirtualMachinePowerStatePoweredOn:
			filteredVMs = append(filteredVMs, vm)

		case includePoweredOff &&
			vm.Runtime.PowerState == types.VirtualMachinePowerStatePoweredOff:
			filteredVMs = append(filteredVMs, vm)

		// Consider suspended VMs to be "powered off"
		case includePoweredOff &&
			vm.Runtime.PowerState == types.VirtualMachinePowerStateSuspended:
			filteredVMs = append(filteredVMs, vm)

		}
	}

	numExcluded := len(vms) - len(filteredVMs)

	return filteredVMs, numExcluded

}

// FilterVMsByPowerCycleUptime filters the provided collection of
// VirtualMachines to just those with WARNING or CRITICAL values based on
// provided thresholds. The collection is returned along with the number of
// VirtualMachines that were excluded.
func FilterVMsByPowerCycleUptime(vms []mo.VirtualMachine, warningThreshold int, criticalThreshold int) ([]mo.VirtualMachine, int) {

	// setup early so we can reference it from deferred stats output
	var vmsWithIssues []mo.VirtualMachine

	funcTimeStart := time.Now()

	defer func(vms []mo.VirtualMachine, filteredVMs *[]mo.VirtualMachine) {
		logger.Printf(
			"It took %v to execute FilterVMsByPowerCycleUptime func (for %d VMs, yielding %d VMs).\n",
			time.Since(funcTimeStart),
			len(vms),
			len(*filteredVMs),
		)
	}(vms, &vmsWithIssues)

	for _, vm := range vms {
		uptime := time.Duration(vm.Summary.QuickStats.UptimeSeconds) * time.Second
		uptimeDays := uptime.Hours() / 24

		// compare against the WARNING threshold as that will net VMs with
		// CRITICAL state as well.
		if uptimeDays > float64(warningThreshold) {
			vmsWithIssues = append(vmsWithIssues, vm)
		}
	}

	numExcluded := len(vms) - len(vmsWithIssues)

	return vmsWithIssues, numExcluded

}

// dedupeVMs receives a list of VirtualMachine values potentially containing
// one or more duplicate values and returns a new list of unique
// VirtualMachine values.
//
// Credit:
// https://www.reddit.com/r/golang/comments/5ia523/idiomatic_way_to_remove_duplicates_in_a_slice/db6qa2e
func dedupeVMs(vmsList []mo.VirtualMachine) []mo.VirtualMachine {

	funcTimeStart := time.Now()

	defer func(vms *[]mo.VirtualMachine) {
		logger.Printf(
			"It took %v to execute dedupeVMs func (evaluated %d VMs).\n",
			time.Since(funcTimeStart),
			len(*vms),
		)
	}(&vmsList)

	seen := make(map[string]struct{}, len(vmsList))
	j := 0
	for _, vm := range vmsList {
		if _, ok := seen[vm.Self.Value]; ok {
			continue
		}
		seen[vm.Self.Value] = struct{}{}
		vmsList[j] = vm
		j++
	}

	return vmsList[:j]
}

// VMNames receives a list of VirtualMachine values and returns a new list of
// VirtualMachine Name values.
func VMNames(vmsList []mo.VirtualMachine) []string {

	vmNames := make([]string, 0, len(vmsList))
	for i := range vmsList {
		vmNames = append(vmNames, vmsList[i].Name)
	}

	return vmNames

}

// GetVMPowerCycleUptimeStatusSummary accepts a list of VirtualMachines and
// threshold values and generates a collection of VirtualMachines that exceeds
// given thresholds along with those given thresholds.
func GetVMPowerCycleUptimeStatusSummary(
	vms []mo.VirtualMachine,
	warningThreshold int,
	criticalThreshold int,
) VirtualMachinePowerCycleUptimeStatus {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute GetVMPowerCycleUptimeStatusSummary func.\n",
			time.Since(funcTimeStart),
		)
	}()

	var vmsCritical []mo.VirtualMachine
	var vmsWarning []mo.VirtualMachine
	var vmsOK []mo.VirtualMachine

	for _, vm := range vms {

		uptime := time.Duration(vm.Summary.QuickStats.UptimeSeconds) * time.Second
		uptimeDays := uptime.Hours() / 24

		switch {
		case uptimeDays > float64(criticalThreshold):
			vmsCritical = append(vmsCritical, vm)

		case uptimeDays > float64(warningThreshold):
			vmsWarning = append(vmsWarning, vm)

		default:
			vmsOK = append(vmsOK, vm)

		}

	}

	return VirtualMachinePowerCycleUptimeStatus{
		VMsCritical:       vmsCritical,
		VMsWarning:        vmsWarning,
		VMsOK:             vmsOK,
		WarningThreshold:  warningThreshold,
		CriticalThreshold: criticalThreshold,
	}

}

// VMPowerCycleUptimeOneLineCheckSummary is used to generate a one-line Nagios
// service check results summary. This is the line most prominent in
// notifications.
func VMPowerCycleUptimeOneLineCheckSummary(
	stateLabel string,
	evaluatedVMs []mo.VirtualMachine,
	uptimeSummary VirtualMachinePowerCycleUptimeStatus,
	rps []mo.ResourcePool,
) string {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute VMPowerCycleUptimeOneLineCheckSummary func.\n",
			time.Since(funcTimeStart),
		)
	}()

	switch {
	case len(uptimeSummary.VMsCritical) > 0:
		return fmt.Sprintf(
			"%s: %d VMs with power cycle uptime exceeding %d days detected (evaluated %d VMs, %d Resource Pools)",
			stateLabel,
			len(uptimeSummary.VMsCritical),
			uptimeSummary.CriticalThreshold,
			len(evaluatedVMs),
			len(rps),
		)

	case len(uptimeSummary.VMsWarning) > 0:
		return fmt.Sprintf(
			"%s: %d VMs with power cycle uptime exceeding %d days detected (evaluated %d VMs, %d Resource Pools)",
			stateLabel,
			len(uptimeSummary.VMsWarning),
			uptimeSummary.WarningThreshold,
			len(evaluatedVMs),
			len(rps),
		)

	default:

		return fmt.Sprintf(
			"%s: No VMs with power cycle uptime exceeding %d days detected (evaluated %d VMs, %d Resource Pools)",
			stateLabel,
			uptimeSummary.WarningThreshold,
			len(evaluatedVMs),
			len(rps),
		)
	}
}

// VMPowerCycleUptimeReport generates a summary of VMs which exceed power
// cycle uptime thresholds along with various verbose details intended to aid
// in troubleshooting check results at a glance. This information is provided
// for use with the Long Service Output field commonly displayed on the
// detailed service check results display in the web UI or in the body of many
// notifications.
func VMPowerCycleUptimeReport(
	c *vim25.Client,
	allVMs []mo.VirtualMachine,
	evaluatedVMs []mo.VirtualMachine,
	uptimeSummary VirtualMachinePowerCycleUptimeStatus,
	vmsToExclude []string,
	evalPoweredOffVMs bool,
	includeRPs []string,
	excludeRPs []string,
	rps []mo.ResourcePool,
) string {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute VMPowerCycleUptimeReport func.\n",
			time.Since(funcTimeStart),
		)
	}()

	rpNames := make([]string, len(rps))
	for i := range rps {
		rpNames[i] = rps[i].Name
	}

	var report strings.Builder

	fmt.Fprintf(
		&report,
		"VMs with high power cycle uptime:%s%s",
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
	)

	switch {
	case len(uptimeSummary.VMsCritical) > 0 || len(uptimeSummary.VMsWarning) > 0:

		vmsWithHighUptime := make(
			[]mo.VirtualMachine,
			0,
			len(uptimeSummary.VMsCritical)+len(uptimeSummary.VMsWarning),
		)

		vmsWithHighUptime = append(vmsWithHighUptime, uptimeSummary.VMsWarning...)
		vmsWithHighUptime = append(vmsWithHighUptime, uptimeSummary.VMsCritical...)

		sort.Slice(vmsWithHighUptime, func(i, j int) bool {
			return vmsWithHighUptime[i].Summary.QuickStats.UptimeSeconds > vmsWithHighUptime[j].Summary.QuickStats.UptimeSeconds
		})

		for _, vm := range vmsWithHighUptime {

			uptime := time.Duration(vm.Summary.QuickStats.UptimeSeconds) * time.Second
			uptimeDays := uptime.Hours() / 24

			fmt.Fprintf(
				&report,
				"* %s: %.2f days%s",
				vm.Name,
				uptimeDays,
				nagios.CheckOutputEOL,
			)
		}
	default:

		fmt.Fprintf(&report, "* None %s", nagios.CheckOutputEOL)

		fmt.Fprintf(
			&report,
			"%sTop 10 VMs, not yet exceeding power cycle uptime thresholds:%s%s",
			nagios.CheckOutputEOL,
			nagios.CheckOutputEOL,
			nagios.CheckOutputEOL,
		)

		topTen := uptimeSummary.TopTenOK()
		switch {
		case len(topTen) == 0:
			fmt.Fprintf(&report, "* None %s", nagios.CheckOutputEOL)
		default:
			for _, vm := range topTen {
				uptime := time.Duration(vm.Summary.QuickStats.UptimeSeconds) * time.Second
				uptimeDays := uptime.Hours() / 24

				fmt.Fprintf(
					&report,
					"* %s: %.2f days%s",
					vm.Name,
					uptimeDays,
					nagios.CheckOutputEOL,
				)
			}
		}

	}

	fmt.Fprintf(
		&report,
		"%sTen most recently started VMs:%s%s",
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
	)

	bottomTen := uptimeSummary.BottomTenOK()
	switch {
	case len(bottomTen) == 0:
		fmt.Fprintf(&report, "* None %s", nagios.CheckOutputEOL)
	default:
		for _, vm := range bottomTen {
			uptime := time.Duration(vm.Summary.QuickStats.UptimeSeconds) * time.Second
			uptimeDays := uptime.Hours() / 24

			fmt.Fprintf(
				&report,
				"* %s: %.2f days%s",
				vm.Name,
				uptimeDays,
				nagios.CheckOutputEOL,
			)
		}
	}

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

// VMDiskConsolidationOneLineCheckSummary is used to generate a one-line Nagios
// service check results summary. This is the line most prominent in
// notifications.
func VMDiskConsolidationOneLineCheckSummary(
	stateLabel string,
	evaluatedVMs []mo.VirtualMachine,
	vmsNeedingConsolidation []mo.VirtualMachine,
	rps []mo.ResourcePool,
) string {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute VMDiskConsolidationOneLineCheckSummary func.\n",
			time.Since(funcTimeStart),
		)
	}()

	switch {
	case len(vmsNeedingConsolidation) > 0:
		return fmt.Sprintf(
			"%s: %d VMs requiring disk consolidation detected (evaluated %d VMs, %d Resource Pools)",
			stateLabel,
			len(vmsNeedingConsolidation),
			len(evaluatedVMs),
			len(rps),
		)

	default:

		return fmt.Sprintf(
			"%s: No VMs requiring disk consolidation detected (evaluated %d VMs, %d Resource Pools)",
			stateLabel,
			len(evaluatedVMs),
			len(rps),
		)
	}
}

// VMDiskConsolidationReport generates a summary of VMs which require disk
// consolidation along with various verbose details intended to aid in
// troubleshooting check results at a glance. This information is provided for
// use with the Long Service Output field commonly displayed on the detailed
// service check results display in the web UI or in the body of many
// notifications.
func VMDiskConsolidationReport(
	c *vim25.Client,
	allVMs []mo.VirtualMachine,
	evaluatedVMs []mo.VirtualMachine,
	vmsNeedingConsolidation []mo.VirtualMachine,
	vmsToExclude []string,
	evalPoweredOffVMs bool,
	includeRPs []string,
	excludeRPs []string,
	rps []mo.ResourcePool,
) string {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute VMDiskConsolidationReport func.\n",
			time.Since(funcTimeStart),
		)
	}()

	rpNames := make([]string, len(rps))
	for i := range rps {
		rpNames[i] = rps[i].Name
	}

	var report strings.Builder

	fmt.Fprintf(
		&report,
		"VMs requiring disk consolidation:%s%s",
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
	)

	switch {
	case len(vmsNeedingConsolidation) > 0:

		sort.Slice(vmsNeedingConsolidation, func(i, j int) bool {
			return vmsNeedingConsolidation[i].Name < vmsNeedingConsolidation[j].Name
		})

		for _, vm := range vmsNeedingConsolidation {
			fmt.Fprintf(
				&report,
				"* %s (%s)%s",
				vm.Name,
				vm.Runtime.PowerState,
				nagios.CheckOutputEOL,
			)
		}

	default:

		fmt.Fprintf(&report, "* None %s", nagios.CheckOutputEOL)

	}

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
		// NOTE: This plugin is hard-coded to evaluate powered off and powered
		// on VMs equally. I'm not sure whether ignoring powered off VMs by
		// default makes sense for this particular plugin.
		//
		// Please share your feedback here if you feel differently:
		// https://github.com/atc0005/check-vmware/discussions/176
		//
		// Please expand on some use cases for ignoring powered off VMs by default.
		true,
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

// VMInteractiveQuestionOneLineCheckSummary is used to generate a one-line
// Nagios service check results summary. This is the line most prominent in
// notifications.
func VMInteractiveQuestionOneLineCheckSummary(
	stateLabel string,
	evaluatedVMs []mo.VirtualMachine,
	vmsNeedingResponse []mo.VirtualMachine,
	rps []mo.ResourcePool,
) string {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute VMInteractiveQuestionOneLineCheckSummary func.\n",
			time.Since(funcTimeStart),
		)
	}()

	switch {
	case len(vmsNeedingResponse) > 0:
		return fmt.Sprintf(
			"%s: %d VMs requiring interactive response detected (evaluated %d VMs, %d Resource Pools)",
			stateLabel,
			len(vmsNeedingResponse),
			len(evaluatedVMs),
			len(rps),
		)

	default:

		return fmt.Sprintf(
			"%s: No VMs requiring interactive response detected (evaluated %d VMs, %d Resource Pools)",
			stateLabel,
			len(evaluatedVMs),
			len(rps),
		)
	}
}

// VMInteractiveQuestionReport generates a summary of VMs which require an
// interactive response along with various verbose details intended to aid in
// troubleshooting check results at a glance. This information is provided for
// use with the Long Service Output field commonly displayed on the detailed
// service check results display in the web UI or in the body of many
// notifications.
func VMInteractiveQuestionReport(
	c *vim25.Client,
	allVMs []mo.VirtualMachine,
	evaluatedVMs []mo.VirtualMachine,
	vmsNeedingResponse []mo.VirtualMachine,
	vmsToExclude []string,
	evalPoweredOffVMs bool,
	includeRPs []string,
	excludeRPs []string,
	rps []mo.ResourcePool,
) string {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute VMInteractiveQuestionReport func.\n",
			time.Since(funcTimeStart),
		)
	}()

	rpNames := make([]string, len(rps))
	for i := range rps {
		rpNames[i] = rps[i].Name
	}

	var report strings.Builder

	fmt.Fprintf(
		&report,
		"VMs requiring interactive response:%s%s",
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
	)

	switch {
	case len(vmsNeedingResponse) > 0:

		sort.Slice(vmsNeedingResponse, func(i, j int) bool {
			return vmsNeedingResponse[i].Name < vmsNeedingResponse[j].Name
		})

		for _, vm := range vmsNeedingResponse {

			var question string
			switch {
			case vm.Summary.Runtime.Question.Text != "":
				question = vm.Summary.Runtime.Question.Text
			default:
				question = "unknown"
			}

			fmt.Fprintf(
				&report,
				"* %s (%q)%s",
				vm.Name,
				question,
				nagios.CheckOutputEOL,
			)
		}

	default:

		fmt.Fprintf(&report, "* None %s", nagios.CheckOutputEOL)

	}

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
		// NOTE: This plugin is used to detect Virtual Machines which are
		// blocked from execution due to an interactive question. At this
		// stage you could argue that they are neither "on" nor "off", but
		// instead are in an in-between state, though it is likely that
		// vSphere would considered them to be in an "off" state,
		// transitioning to an "on" state. Either way, we report here that
		// both powered on and powered off VMs are evaluated for simplicity.
		true,
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
