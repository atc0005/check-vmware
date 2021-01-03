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
	"os"
	"sort"
	"strings"
	"time"

	"github.com/atc0005/check-vmware/internal/textutils"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
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
		fmt.Fprintf(
			os.Stderr,
			"It took %v to execute GetVMs func (and retrieve %d VMs).\n",
			time.Since(funcTimeStart),
			len(*vms),
		)
	}(&vms)

	m := view.NewManager(c)

	v, err := m.CreateContainerView(
		ctx,
		c.ServiceContent.RootFolder,

		// Q: Is this a selection of what appears "in" the view we are creating?
		// A: Based on testing, it appears so. Much like a database view, you
		// can tie together multiple types into a single view that can be
		// selectively queried.
		[]string{
			"VirtualMachine",
			//
			// Q: What difference does it make to specify additional Managed
			// Object Types here?
			// A: It exposes additional types from this view.
			//
			// "Network",
			// "ResourcePool",
		},
		true,
	)
	if err != nil {
		return nil, err
	}

	defer func() {
		// Per vSphere Web Services SDK Programming Guide - VMware vSphere 7.0
		// Update 1:
		//
		// A best practice when using views is to call the DestroyView()
		// method when a view is no longer needed. This practice frees memory
		// on the server.
		if err := v.Destroy(ctx); err != nil {
			fmt.Println("Error occurred while destroying view")
		}
	}()

	// If the properties slice is nil, all properties are loaded.
	var props []string
	if propsSubset {
		// https://code.vmware.com/apis/1067/vsphere
		// https://vdc-download.vmware.com/vmwb-repository/dcr-public/a5f4000f-1ea8-48a9-9221-586adff3c557/7ff50256-2cf2-45ea-aacd-87d231ab1ac7/vim.VirtualMachine.html
		props = []string{
			"summary",
			"datastore",
			"resourcePool",
			"config",
			"snapshot",
			"guest",
			"name",
			"network",
			"runtime",
		}
	}

	err = v.Retrieve(
		ctx,

		// Q: Is this meant to indicate just one "kind" from the view?
		//
		// I: This type has to match the destination slice type setup
		// previously. If `mo.VirtualMachine` is the slice type, then
		// `VirtualMachine` is the required, single slice entry value here.
		// Adding other "kind" values here results in this library attempting
		// to retrieve and assign types (e.g.,
		// `mo.DistributedVirtualPortgroup`) to the slice of
		// `mo.VirtualMachine` causing a panic.
		//
		// Alternatively, it appears you can use a slice of empty interface,
		// then use a type switch later to figure out what you're working
		// with.
		[]string{
			"VirtualMachine",
			// "Network",
		},
		// https://code.vmware.com/apis/1067/vsphere
		// https://vdc-download.vmware.com/vmwb-repository/dcr-public/a5f4000f-1ea8-48a9-9221-586adff3c557/7ff50256-2cf2-45ea-aacd-87d231ab1ac7/vim.VirtualMachine.html
		props,
		&vms,
	)
	if err != nil {
		return nil, err
	}

	sort.Slice(vms, func(i, j int) bool {
		return strings.ToLower(vms[i].Name) < strings.ToLower(vms[j].Name)
	})

	return vms, nil

}

// GetVMsFromRP receives a list of ResourcePool object references and returns
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
	// vms := make([]mo.VirtualMachine, 0, 100)
	var vms []mo.VirtualMachine

	defer func(vms *[]mo.VirtualMachine) {
		fmt.Fprintf(
			os.Stderr,
			"It took %v to execute GetVMsFromRPs func (and retrieve %d VMs).\n",
			time.Since(funcTimeStart),
			len(*vms),
		)
	}(&vms)

	m := view.NewManager(c)

	for _, rp := range rps {

		// attempt to limit view to VMs within the resource pool
		v, err := m.CreateContainerView(
			ctx,
			rp.Reference(),
			[]string{
				"VirtualMachine",
			},
			true,
		)
		if err != nil {
			return nil, err
		}

		defer func() {
			// Per vSphere Web Services SDK Programming Guide - VMware vSphere 7.0
			// Update 1:
			//
			// A best practice when using views is to call the DestroyView()
			// method when a view is no longer needed. This practice frees memory
			// on the server.
			if err := v.Destroy(ctx); err != nil {
				fmt.Println("Error occurred while destroying view")
			}
		}()

		// If the properties slice is nil, all properties are loaded.
		var props []string
		if propsSubset {
			// https://code.vmware.com/apis/1067/vsphere
			// https://vdc-download.vmware.com/vmwb-repository/dcr-public/a5f4000f-1ea8-48a9-9221-586adff3c557/7ff50256-2cf2-45ea-aacd-87d231ab1ac7/vim.VirtualMachine.html
			props = []string{
				"summary",
				"datastore",
				"resourcePool",
				"config",
				"snapshot",
				"guest",
				"name",
				"network",
				"runtime",
			}
		}

		err = v.Retrieve(ctx, []string{"VirtualMachine"}, props, &vms)
		if err != nil {
			return nil, err
		}
	}

	sort.Slice(vms, func(i, j int) bool {
		return strings.ToLower(vms[i].Name) < strings.ToLower(vms[j].Name)
	})

	return vms, nil

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
		fmt.Fprintf(
			os.Stderr,
			"It took %v to execute GetVMByName func.\n",
			time.Since(funcTimeStart),
		)
	}()

	finder := find.NewFinder(c, true)

	var dc *object.Datacenter
	var findDCErr error
	var errMsg string
	switch {
	case datacenter == "":
		dc, findDCErr = finder.DefaultDatacenter(ctx)
		errMsg = "error: datacenter not provided, failed to fallback to default datacenter"
	default:
		dc, findDCErr = finder.DatacenterOrDefault(ctx, datacenter)
		errMsg = "error: failed to use provided datacenter, failed to fallback to default datacenter"
	}

	if findDCErr != nil {
		return mo.VirtualMachine{}, fmt.Errorf("%s: %w", errMsg, findDCErr)
	}
	finder.SetDatacenter(dc)

	vmo, err := finder.VirtualMachine(ctx, vmName)
	if err != nil {
		return mo.VirtualMachine{}, err
	}

	// If the properties slice is nil, all properties are loaded.
	var props []string
	if propsSubset {
		// https://code.vmware.com/apis/1067/vsphere
		// https://vdc-download.vmware.com/vmwb-repository/dcr-public/a5f4000f-1ea8-48a9-9221-586adff3c557/7ff50256-2cf2-45ea-aacd-87d231ab1ac7/vim.VirtualMachine.html
		props = []string{
			"summary",
			"datastore",
			"resourcePool",
			"config",
			"snapshot",
			"guest",
			"name",
			"network",
		}
	}

	var vm mo.VirtualMachine
	err = vmo.Common.Properties(
		ctx,
		vmo.Reference(),
		props,
		&vm,
	)

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
		fmt.Fprintf(
			os.Stderr,
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
		fmt.Fprintf(
			os.Stderr,
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
