// Copyright 2021 Adam Chalkley
//
// https://github.com/atc0005/check-vmware
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package vsphere

// ParentResourcePool represents the hidden resource pool named Resources
// which is present on virtual machine hosts. This resource pool is a parent
// of all resource pools of the host. Including this pool in "eligible"
// resource pool lists throws off calculations (e.g., causes a VM to show up
// twice).
const ParentResourcePool string = "Resources"

const dcFailedToUseFailedToFallback string = "error: failed to use provided datacenter, failed to fallback to default datacenter"

const dcNotProvidedFailedToFallback string = "error: datacenter not provided, failed to fallback to default datacenter"

const crFailedToUseFailedToFallback string = "error: failed to use provided cluster to obtain compute resource, failed to fallback to default compute resource"

const crNotProvidedFailedToFallback string = "error: cluster not provided, failed to fallback to default compute resource"

// virtualHardwareVersionPrefix is used as a prefix for virtual hardware
// versions used by VirtualMachines. Examples include vmx-15, vmx-14 and so on.
const virtualHardwareVersionPrefix string = "vmx-"

// CustomAttributeValNotSet is used to indicate that a Custom Attribute value
// was not set on an object.
const CustomAttributeValNotSet string = "NotSet"

// Managed Object Reference types
const (
	MgObjRefTypeFolder          string = "Folder"
	MgObjRefTypeDatacenter      string = "Datacenter"
	MgObjRefTypeComputeResource string = "ComputeResource"
	MgObjRefTypeResourcePool    string = "ResourcePool"
	MgObjRefTypeHostSystem      string = "HostSystem"
)

// used with snapshots reports that provide Long Service Output
const (
	snapshotThresholdTypeAge   string = "age"
	snapshotThresholdTypeCount string = "count"
	snapshotThresholdTypeSize  string = "size"
)

// used with snapshots reports that provide Long Service Output
const (
	snapshotThresholdTypeAgeSuffix   string = "day"
	snapshotThresholdTypeCountSuffix string = "snapshots"
	snapshotThresholdTypeSizeSuffix  string = "GB"
)
