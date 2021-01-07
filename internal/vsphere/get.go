package vsphere

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/property"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

func getVirtualMachinePropsSubset() []string {
	// https://code.vmware.com/apis/1067/vsphere
	// https://vdc-download.vmware.com/vmwb-repository/dcr-public/a5f4000f-1ea8-48a9-9221-586adff3c557/7ff50256-2cf2-45ea-aacd-87d231ab1ac7/vim.VirtualMachine.html
	return []string{
		"summary",
		"datastore",
		"resourcePool",
		"config",
		"snapshot",
		"guest",
		"name",
		"network",
		"runtime", // Host system is listed here
		"customValue",
		"availableField",
	}
}
func getNetworkPropsSubset() []string {
	// https://code.vmware.com/apis/1067/vsphere
	// https://vdc-download.vmware.com/vmwb-repository/dcr-public/a5f4000f-1ea8-48a9-9221-586adff3c557/7ff50256-2cf2-45ea-aacd-87d231ab1ac7/vim.Network.html
	return []string{
		"summary", // properties of this network
		"name",    // name of this network
		"host",    // hosts attached to this network
		"vm",      // virtual machines using this network
	}
}
func getResourcePoolPropsSubset() []string {
	// https://code.vmware.com/apis/1067/vsphere
	// https://vdc-download.vmware.com/vmwb-repository/dcr-public/a5f4000f-1ea8-48a9-9221-586adff3c557/7ff50256-2cf2-45ea-aacd-87d231ab1ac7/vim.ResourcePool.html
	return []string{
		"summary",
		"resourcePool", // potential child resource pools
		"config",
		"name",
		"runtime",
	}
}
func getHostSystemPropsSubset() []string {
	// https://code.vmware.com/apis/1067/vsphere
	// https://vdc-download.vmware.com/vmwb-repository/dcr-public/a5f4000f-1ea8-48a9-9221-586adff3c557/7ff50256-2cf2-45ea-aacd-87d231ab1ac7/vim.HostSystem.html
	return []string{
		"summary",
		"vm",
		"name",
		"datastore",
		"customValue",
		"availableField",
	}
}
func getDatastorePropsSubset() []string {
	// https://code.vmware.com/apis/1067/vsphere
	// https://vdc-download.vmware.com/vmwb-repository/dcr-public/a5f4000f-1ea8-48a9-9221-586adff3c557/7ff50256-2cf2-45ea-aacd-87d231ab1ac7/vim.Datastore.html
	return []string{
		"summary",
		"vm",
		"host",
		"name",
		"customValue",
		"availableField",
	}
}

func getObjects(ctx context.Context, c *vim25.Client, dst interface{}, objRef types.ManagedObjectReference, propsSubset bool) error {

	funcTimeStart := time.Now()

	var objKind string

	// If the properties slice is nil, all properties are loaded.
	var props []string

	// this is set just before this deferred func executes due to deferred
	// length checks in type switch below
	var objCount int
	defer func(count *int, kind *string) {
		fmt.Fprintf(
			os.Stderr,
			"It took %v to execute getObjects func (and retrieve %d %s objects).\n",
			time.Since(funcTimeStart),
			*count,
			*kind,
		)
	}(&objCount, &objKind)

	// Create a view of caller-specified objects
	m := view.NewManager(c)

	switch u := dst.(type) {
	case *[]mo.Datastore:
		defer func() {
			objCount = len(*u)
		}()

		objKind = "Datastore"

		if propsSubset {
			props = getDatastorePropsSubset()
		}

	case *[]mo.HostSystem:
		defer func() {
			objCount = len(*u)
		}()
		objKind = "HostSystem"

		if propsSubset {
			props = getHostSystemPropsSubset()
		}

	case *[]mo.VirtualMachine:
		defer func() {
			objCount = len(*u)
		}()
		objKind = "VirtualMachine"

		if propsSubset {
			props = getVirtualMachinePropsSubset()
		}

	case *[]mo.Network:
		defer func() {
			objCount = len(*u)
		}()
		objKind = "Network"

		if propsSubset {
			props = getNetworkPropsSubset()
		}

	case *[]mo.ResourcePool:
		defer func() {
			objCount = len(*u)
		}()
		objKind = "ResourcePool"

		if propsSubset {
			props = getResourcePoolPropsSubset()
		}

	default:

		return fmt.Errorf("func getObjects: unknown type provided as destination")

	}

	v, err := m.CreateContainerView(
		ctx,
		objRef,
		[]string{objKind},
		true,
	)
	if err != nil {
		return err
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

	err = v.Retrieve(ctx, []string{objKind}, props, dst)
	if err != nil {
		return err
	}

	return nil

}

func getObjectByName(ctx context.Context, c *vim25.Client, dst interface{}, objName string, datacenter string, propsSubset bool) error {

	funcTimeStart := time.Now()

	var objKind string

	defer func(kind *string) {
		fmt.Fprintf(
			os.Stderr,
			"It took %v to execute getObjectByName func (and retrieve %s object).\n",
			time.Since(funcTimeStart),
			*kind,
		)
	}(&objKind)

	finder := find.NewFinder(c, true)

	switch {
	case datacenter == "":
		dc, findDCErr := finder.DefaultDatacenter(ctx)
		if findDCErr != nil {
			return fmt.Errorf("%s: %w", dcNotProvidedFailedToFallback, findDCErr)
		}
		finder.SetDatacenter(dc)

	default:
		dc, findDCErr := finder.DatacenterOrDefault(ctx, datacenter)
		if findDCErr != nil {
			return fmt.Errorf("%s: %w", failedToUseFailedToFallback, findDCErr)
		}
		finder.SetDatacenter(dc)
	}

	// If the properties slice is nil, all properties are loaded.
	var props []string

	pc := property.DefaultCollector(c)

	var objRef types.ManagedObjectReference

	switch dst.(type) {
	case *mo.Datastore:

		objKind = "Datastore"
		if propsSubset {
			props = getDatastorePropsSubset()
		}

		obj, err := finder.Datastore(ctx, objName)
		if err != nil {
			return err
		}
		objRef = obj.Reference()

	case *mo.HostSystem:

		objKind = "HostSystem"
		if propsSubset {
			props = getHostSystemPropsSubset()
		}

		obj, err := finder.HostSystem(ctx, objName)
		if err != nil {
			return err
		}
		objRef = obj.Reference()

	case *mo.VirtualMachine:

		objKind = "VirtualMachine"
		if propsSubset {
			props = getVirtualMachinePropsSubset()
		}

		obj, err := finder.VirtualMachine(ctx, objName)
		if err != nil {
			return err
		}

		objRef = obj.Reference()

	case *mo.Network:

		objKind = "Network"
		if propsSubset {
			props = getNetworkPropsSubset()
		}

		obj, err := finder.Network(ctx, objName)
		if err != nil {
			return err
		}

		objRef = obj.Reference()

	case *mo.ResourcePool:

		objKind = "ResourcePool"
		if propsSubset {
			props = getResourcePoolPropsSubset()
		}

		obj, err := finder.ResourcePool(ctx, objName)
		if err != nil {
			return err
		}

		objRef = obj.Reference()

	default:

		objKind = "unknown"

		return fmt.Errorf("func getObjectByName: unknown type provided as destination")

	}

	err := pc.RetrieveOne(
		ctx,
		objRef,
		props,
		&dst,
	)

	if err != nil {
		return err
	}

	return nil

}
