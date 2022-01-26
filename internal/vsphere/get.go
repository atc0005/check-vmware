package vsphere

import (
	"context"
	"fmt"
	"time"

	"github.com/vmware/govmomi"
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
		// "rootSnapshot", // TODO: need for this?
		"storage",
		"guest",
		"layoutEx",
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
func getVirtualAppPropsSubset() []string {
	// https://code.vmware.com/apis/1067/vsphere
	// https://vdc-download.vmware.com/vmwb-repository/dcr-public/a5f4000f-1ea8-48a9-9221-586adff3c557/7ff50256-2cf2-45ea-aacd-87d231ab1ac7/vim.VirtualApp.html

	// All of the properties that we need from a VirtualApp are inherited from
	// the enclosing ResourcePool, so we just use the same properties list
	// here also.
	return getResourcePoolPropsSubset()
}
func getHostSystemPropsSubset() []string {
	// https://code.vmware.com/apis/1067/vsphere
	// https://vdc-download.vmware.com/vmwb-repository/dcr-public/a5f4000f-1ea8-48a9-9221-586adff3c557/7ff50256-2cf2-45ea-aacd-87d231ab1ac7/vim.HostSystem.html
	return []string{
		"hardware", // memory capacity
		"runtime",  // connection, power state details
		"summary",
		"vm",
		"name",
		"datastore",
		"customValue",
		"availableField",
		"parent", // used to obtain ComputeResource
	}
}
func getDatastorePropsSubset() []string {
	// https://code.vmware.com/apis/1067/vsphere
	// https://vdc-download.vmware.com/vmwb-repository/dcr-public/a5f4000f-1ea8-48a9-9221-586adff3c557/7ff50256-2cf2-45ea-aacd-87d231ab1ac7/vim.Datastore.html
	return []string{
		"summary",
		"vm",
		"host",
		"iormConfiguration", // unreliable if DatastoreSummary.Accessible != true; used to determine whether stats are being collected
		"name",
		"customValue",
		"availableField",
	}
}
func getDatacenterPropsSubset() []string {
	// https://code.vmware.com/apis/1067/vsphere
	// https://vdc-download.vmware.com/vmwb-repository/dcr-public/a5f4000f-1ea8-48a9-9221-586adff3c557/7ff50256-2cf2-45ea-aacd-87d231ab1ac7/vim.Datacenter.html
	return []string{
		"name",
		"overallStatus",
		"triggeredAlarmState",
	}
}
func getAlarmPropsSubset() []string {
	// https://code.vmware.com/apis/1067/vsphere
	// https://vdc-download.vmware.com/vmwb-repository/dcr-public/a5f4000f-1ea8-48a9-9221-586adff3c557/7ff50256-2cf2-45ea-aacd-87d231ab1ac7/vim.alarm.Alarm.html
	return []string{
		"info",
	}
}

// getObjects retrieves one or more objects, filtered by the provided
// container type ManagedObjectReference. An error is returned if the provided
// ManagedObjectReference is not for a supported container type.
func getObjects(
	ctx context.Context,
	c *vim25.Client,
	dst interface{},
	objRef types.ManagedObjectReference,
	propsSubset bool,
	recursive bool,
) error {

	funcTimeStart := time.Now()

	var objKind string

	// If the properties slice is nil, all properties are loaded.
	var props []string

	// this is set just before this deferred func executes due to deferred
	// length checks in type switch below
	var objCount int
	defer func(count *int, kind *string) {
		logger.Printf(
			"It took %v to execute getObjects func (and retrieve %d %s objects from %s).\n",
			time.Since(funcTimeStart),
			*count,
			*kind,
			objRef.Type,
		)
	}(&objCount, &objKind)

	// Create a view of caller-specified objects
	m := view.NewManager(c)

	logger.Printf("Requested objRef type is %s", objRef.Type)

	// https://vdc-download.vmware.com/vmwb-repository/dcr-public/a5f4000f-1ea8-48a9-9221-586adff3c557/7ff50256-2cf2-45ea-aacd-87d231ab1ac7/vim.view.ContainerView.html
	switch objRef.Type {
	case MgObjRefTypeFolder:
	case MgObjRefTypeDatacenter:
	case MgObjRefTypeComputeResource:
	case MgObjRefTypeResourcePool:
	case MgObjRefTypeHostSystem:

	// A VirtualApp is not documented as a supported managed object type to
	// use for a container view, but testing has shown that it works for our
	// purposes.
	case MgObjRefTypeVirtualApp:

	default:
		return fmt.Errorf(
			"unsupported container type specified for ContainerView: %s",
			objRef.Type,
		)
	}

	switch u := dst.(type) {
	case *[]mo.Datacenter:
		defer func() {
			objCount = len(*u)
		}()

		objKind = "Datacenter"

		if propsSubset {
			props = getDatacenterPropsSubset()
		}

	case *[]mo.Alarm:
		defer func() {
			objCount = len(*u)
		}()

		objKind = "Alarm"

		if propsSubset {
			props = getAlarmPropsSubset()
		}

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

	case *[]mo.VirtualApp:
		defer func() {
			objCount = len(*u)
		}()
		objKind = "VirtualApp"

		if propsSubset {
			props = getVirtualAppPropsSubset()
		}

	default:

		return fmt.Errorf("func getObjects: unknown type provided as destination")

	}

	// FIXME: Should this filter to a specific datacenter? See GH-219.
	v, err := m.CreateContainerView(
		ctx,
		objRef,
		[]string{objKind},
		recursive,
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
		//
		// nolint:govet // err intentionally scoped; shadowing not a concern.
		if err := v.Destroy(ctx); err != nil {
			logger.Printf("Error occurred while destroying view: %s", err)
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
		logger.Printf(
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
			return fmt.Errorf("%s: %w", dcFailedToUseFailedToFallback, findDCErr)
		}
		finder.SetDatacenter(dc)
	}

	// If the properties slice is nil, all properties are loaded.
	var props []string

	pc := property.DefaultCollector(c)

	switch u := dst.(type) {
	case *mo.Datastore:

		objKind = "Datastore"
		if propsSubset {
			props = getDatastorePropsSubset()
		}

		obj, err := finder.Datastore(ctx, objName)
		if err != nil {
			return err
		}

		err = pc.RetrieveOne(
			ctx,
			obj.Reference(),
			props,
			u,
		)

		if err != nil {
			return err
		}

	case *mo.HostSystem:

		objKind = "HostSystem"
		if propsSubset {
			props = getHostSystemPropsSubset()
		}

		obj, err := finder.HostSystem(ctx, objName)
		if err != nil {
			return err
		}

		err = pc.RetrieveOne(
			ctx,
			obj.Reference(),
			props,
			u,
		)

		if err != nil {
			return err
		}

	case *mo.VirtualMachine:

		objKind = "VirtualMachine"
		if propsSubset {
			props = getVirtualMachinePropsSubset()
		}

		obj, err := finder.VirtualMachine(ctx, objName)
		if err != nil {
			return err
		}

		err = pc.RetrieveOne(
			ctx,
			obj.Reference(),
			props,
			u,
		)

		if err != nil {
			return err
		}

	case *mo.Network:

		objKind = "Network"
		if propsSubset {
			props = getNetworkPropsSubset()
		}

		obj, err := finder.Network(ctx, objName)
		if err != nil {
			return err
		}

		err = pc.RetrieveOne(
			ctx,
			obj.Reference(),
			props,
			u,
		)

		if err != nil {
			return err
		}

	case *mo.ResourcePool:

		objKind = "ResourcePool"
		if propsSubset {
			props = getResourcePoolPropsSubset()
		}

		obj, err := finder.ResourcePool(ctx, objName)
		if err != nil {
			return err
		}

		err = pc.RetrieveOne(
			ctx,
			obj.Reference(),
			props,
			u,
		)

		if err != nil {
			return err
		}

	default:

		objKind = "unknown"

		return fmt.Errorf("func getObjectByName: unknown type provided as destination")

	}

	return nil

}

// getResourcePools retrieves Resource Pools for VirtualMachine or
// ResourcePool types via the provided ManagedObjectReference (moRef) and
// Datacenter name. If the moRef is for a VirtualMachine, one ResourcePool is
// returned. If the moRef is for a ResourcePool then two are returned: the
// matching ResourcePool and the parent ResourcePool. If the moRef is for
// another ManagedEntity type an empty collection is returned. If specified, a
// subset of all properties are returned for discovered ResourcePools.
func getResourcePools(ctx context.Context, c *govmomi.Client, moRef types.ManagedObjectReference, propsSubset bool) ([]mo.ResourcePool, error) {

	funcTimeStart := time.Now()

	// Up to 2 can be returned, so set the initial size to match. Declare this
	// early so that we can grab a pointer to it in order to access the
	// entries later
	resourcePools := make([]mo.ResourcePool, 0, 2)

	defer func(rp *[]mo.ResourcePool) {
		logger.Printf(
			"It took %v to execute getResourcePools func (and retrieve %d ResourcePools).\n",
			time.Since(funcTimeStart),
			len(*rp),
		)
	}(&resourcePools)

	switch {
	case moRef.Type == MgObjRefTypeResourcePool:

		// if the associated entity is a Resource Pool, then its name
		// and the name of its parent (another Resource Pool) should
		// be considered.

		var rp mo.ResourcePool
		var rpParent mo.ResourcePool
		var rpProps []string

		if propsSubset {
			rpProps = getResourcePoolPropsSubset()
		}

		// Fetch Resource Pool directly associated with the Triggered
		// Alarm entity
		err := c.RetrieveOne(ctx, moRef, rpProps, &rp)
		if err != nil {
			return nil, err
		}

		// Add *this* Resource Pool
		resourcePools = append(resourcePools, rp)

		// Fetch the parent Resource Pool for the Resource Pool
		// associated with the Triggered Alarm.
		err = c.RetrieveOne(ctx, rp.Self, rpProps, &rpParent)
		if err != nil {
			return nil, err
		}

		// Add the parent Resource Pool
		resourcePools = append(resourcePools, rpParent)

		return resourcePools, nil

	case moRef.Type == MgObjRefTypeVirtualMachine:
		var vm mo.VirtualMachine
		var rp mo.ResourcePool
		var vmProps []string
		var rpProps []string

		if propsSubset {
			vmProps = getVirtualMachinePropsSubset()
			rpProps = getResourcePoolPropsSubset()
		}

		// Fetch VirtualMachine associated with Triggered Alarm
		// entity.
		err := c.RetrieveOne(ctx, moRef, vmProps, &vm)
		if err != nil {
			return nil, err
		}

		// guard against missing resource pool (nil pointer dereferencing)
		if vm.ResourcePool == nil {
			return nil, fmt.Errorf(
				"resource pool MOID not set for %q (%q); "+
					"permissions to view associated resource pool may not be available",
				vm.Name,
				vm.Self,
			)
		}

		// Fetch Resource Pool for VirtualMachine associated with
		// Triggered Alarm entity.
		err = c.RetrieveOne(ctx, *vm.ResourcePool, rpProps, &rp)
		if err != nil {
			return nil, err
		}

		// Add the VirtualMachine's ResourcePool.
		resourcePools = append(resourcePools, rp)

		return resourcePools, nil

	default:
		// As far as I know, no other types can be "part" of a
		// Resource Pool. Return an empty collection.
		return []mo.ResourcePool{}, nil
	}

}
