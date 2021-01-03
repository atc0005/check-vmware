// Copyright 2021 Adam Chalkley
//
// https://github.com/atc0005/check-vmware
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package vsphere

// Virtual machine hosts have a hidden resource pool named Resources, which is
// a parent of all resource pools of the host. Including this pool in
// "eligible" resource pool lists throws off calculations (e.g., causes a VM
// to show up twice).
const ParentResourcePool string = "Resources"
