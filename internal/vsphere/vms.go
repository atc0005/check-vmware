// Copyright 2021 Adam Chalkley
//
// https://github.com/atc0005/check-vmware
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package vsphere

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/atc0005/check-vmware/internal/textutils"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

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

// GetVMsFromRPs receives a list of ResourcePool object references and returns
// a list of VirtualMachine object references. The propsSubset boolean value
// indicates whether a subset of properties per VirtualMachine are retrieved.
// If requested, a subset of all available properties will be retrieved
// (faster) instead of recursively fetching all properties (about 2x as slow)
// A collection of VirtualMachines with requested properties is returned or
// nil and an error, if one occurs.
func GetVMsFromRPs(ctx context.Context, c *vim25.Client, rps []mo.ResourcePool, propsSubset bool) ([]mo.VirtualMachine, error) {

	funcTimeStart := time.Now()

	// declare this early so that we can grab a pointer to it in order to
	// access the entries later
	var vms []mo.VirtualMachine

	defer func(vms *[]mo.VirtualMachine) {
		logger.Printf(
			"It took %v to execute GetVMsFromRPs func (and retrieve %d VMs).\n",
			time.Since(funcTimeStart),
			len(*vms),
		)
	}(&vms)

	for _, rp := range rps {

		err := getObjects(ctx, c, &vms, rp.Reference(), propsSubset)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to retrieve VirtualMachines from Resource Pool %s: %w",
				rp.Name,
				err,
			)
		}
	}

	// remove any potential duplicate entries which could occur if we are
	// evaluating the (default, hidden) 'Resources' Resource Pool
	vms = dedupeVMs(vms)

	sort.Slice(vms, func(i, j int) bool {
		return strings.ToLower(vms[i].Name) < strings.ToLower(vms[j].Name)
	})

	return vms, nil

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
		vm, err := FilterVMByID(allVMs, ds.Vm[i].Value)
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

// FilterVMByName accepts a collection of VirtualMachines and a VirtualMachine
// name to filter against. An error is returned if the list of VirtualMachines
// is empty or if a match was not found.
func FilterVMByName(vms []mo.VirtualMachine, vmName string) (mo.VirtualMachine, error) {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute FilterVMByName func.\n",
			time.Since(funcTimeStart),
		)
	}()

	if len(vms) == 0 {
		return mo.VirtualMachine{}, fmt.Errorf("received empty list of virtual machines to filter by name")
	}

	for _, vm := range vms {
		if vm.Name == vmName {
			return vm, nil
		}
	}

	return mo.VirtualMachine{}, fmt.Errorf(
		"error: failed to retrieve VirtualMachine using provided name %q",
		vmName,
	)

}

// FilterVMByID receives a collection of VirtualMachines and a VirtualMachine
// ID to filter against. An error is returned if the list of VirtualMachines
// is empty or if a match was not found.
func FilterVMByID(vms []mo.VirtualMachine, vmID string) (mo.VirtualMachine, error) {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute FilterVMByID func.\n",
			time.Since(funcTimeStart),
		)
	}()

	if len(vms) == 0 {
		return mo.VirtualMachine{}, fmt.Errorf("received empty list of virtual machines to filter by ID")
	}

	for _, vm := range vms {
		// return match, if available
		if vm.Summary.Vm.Value == vmID {
			return vm, nil
		}
	}

	return mo.VirtualMachine{}, fmt.Errorf(
		"error: failed to retrieve VirtualMachine using provided ID %q",
		vmID,
	)

}

// ExcludeVMsByName receives a collection of VirtualMachines and a list of VMs
// that should be ignored. A new collection minus ignored VirtualMachines is
// returned. If the collection of VirtualMachine is empty, an empty collection
// is returned. If the list of ignored VirtualMachines is empty, the same
// items from the received collection of VirtualMachines is returned. If the
// list of ignored VirtualMachines is greater than the list of received
// VirtualMachines, then only matching VirtualMachines will be excluded and
// any others silently skipped.
func ExcludeVMsByName(allVMs []mo.VirtualMachine, ignoreList []string) []mo.VirtualMachine {

	if len(allVMs) == 0 || len(ignoreList) == 0 {
		return allVMs
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

	return vmsToKeep

}

// FilterVMsByPowerState accepts a collection of VirtualMachines and a boolean
// value to indicate whether powered off VMs should be included in the
// returned collection. If the collection of provided VirtualMachines is
// empty, an empty collection is returned.
func FilterVMsByPowerState(vms []mo.VirtualMachine, includePoweredOff bool) []mo.VirtualMachine {

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
		return vms
	}

	for _, vm := range vms {
		switch {
		// case includePoweredOff && vm.Guest.ToolsStatus != types.VirtualMachineToolsStatusToolsOk:
		// 	vmsWithIssues = append(vmsWithIssues, vm)

		case vm.Runtime.PowerState == types.VirtualMachinePowerStatePoweredOn:
			filteredVMs = append(filteredVMs, vm)

		case includePoweredOff &&
			vm.Runtime.PowerState == types.VirtualMachinePowerStatePoweredOff:
			filteredVMs = append(filteredVMs, vm)

		}
	}

	return filteredVMs

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
		if _, ok := seen[vm.Summary.Vm.Value]; ok {
			continue
		}
		seen[vm.Summary.Vm.Value] = struct{}{}
		vmsList[j] = vm
		j++
	}

	return vmsList[:j]
}
