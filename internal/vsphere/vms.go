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
	"io"
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

// ErrVirtualMachineMissingBackupDate indicates that a Virtual Machine is
// missing an expected backup date (or value).
var ErrVirtualMachineMissingBackupDate = errors.New(
	"virtual machine missing backup date",
)

// ErrVirtualMachineBackupDateOld indicates that a Virtual Machine has a
// backup date specified, but the backup is older than the user specified
// threshold.
var ErrVirtualMachineBackupDateOld = errors.New(
	"virtual machine backup date exceeds specified threshold",
)

// VMWithCA wraps the vSphere VirtualMachine managed object type with a
// specific Custom Attribute name/value pair.
type VMWithCA struct {
	mo.VirtualMachine

	// CustomAttribute represents the name/value Custom Attribute pair as
	// specified by the user. This is often used to initially filter a
	// collection to a specific attribute (e.g., VMs with current backup
	// information) and then review further to determine whether the VM passes
	// specific criteria (e.g., backup within an expected window of time).
	CustomAttribute CustomAttribute
}

// VMWithCAs wraps the vSphere VirtualMachine managed object type with a
// Custom Attribute name/value index (map) of all Custom Attributes associated
// with the managed object.
type VMWithCAs struct {
	mo.VirtualMachine

	// Custom Attribute name/value index (map) of all Custom Attributes
	// associated with the managed object.
	CustomAttributes CustomAttributes
}

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

// VMWithBackup is a VirtualMachine with backup date details.
type VMWithBackup struct {
	// mo.VirtualMachine

	// VMWithCAs is embedded to provide access to the original
	// mo.VirtualMachine value and also access to the complete index of Custom
	// Attributes for the VM.
	VMWithCAs

	// BackupDateCA is the Custom Attribute value indicating when the last
	// backup occurred for this VirtualMachine.
	// TODO: Perhaps implement as a method?
	// TODO: Perhaps have this as a string field, use a method to retrieve the value from the included map?
	// BackupDateCA CustomAttribute

	// BackupDateCAName is the name (not the value) of the Custom Attribute
	// which indicates when the last backup occurred for this VirtualMachine.
	BackupDateCAName string

	// BackupMetadataCA is the Custom Attribute value providing additional
	// context for the last backup for this VirtualMachine.
	// TODO: Perhaps implement as a method?
	// TODO: Perhaps have this as a string field, use a method to retrieve the value from the included map?
	// BackupMetadataCA CustomAttribute

	// BackupMetadataCAName is the name (not the value) of the Custom
	// Attribute which provides additional context for the last backup for
	// this VirtualMachine.
	BackupMetadataCAName string

	// BackupDate is the date/time of the last backup for this VirtualMachine.
	// This value is set to the user-specified time zone or location (or the
	// default location if not specified).
	//
	// TODO: Perhaps implement as a method, record the Location instead?
	BackupDate time.Time

	// WarningAgeInDaysThreshold is the age in days when a recorded backup
	// date is considered stale and a WARNING threshold reached (but not yet
	// crossed).
	WarningAgeInDaysThreshold int

	// CriticalAgeInDaysThreshold is the age in days when a recorded backup
	// date is considered stale and a WARNING threshold reached (but not yet
	// crossed).
	CriticalAgeInDaysThreshold int
}

// VMsWithBackup is a collection of VirtualMachines which track backup date
// details.
type VMsWithBackup []VMWithBackup

// IsWarningState indicates whether the WARNING threshold has been crossed or
// if the Virtual Machine is missing the expected Custom Attribute for
// tracking last backup date or if the Custom Attribute has an empty value.
func (vmwb VMWithBackup) IsWarningState() bool {

	// Look for the requested Custom Attributes
	backupDateCA, hasBackupDateCA := vmwb.CustomAttributes[vmwb.BackupDateCAName]

	// NOTE: This is an optional Custom Attribute, so we don't require it here.
	// _, hasBackupMetadataCA := vmwb.CustomAttributes[vmwb.BackupMetadataCAName]

	if !hasBackupDateCA {
		logger.Printf(
			"Custom Attribute %q missing from %s",
			vmwb.BackupDateCAName,
			vmwb.Name,
		)

		return true
	}

	if strings.TrimSpace(backupDateCA) == "" {
		logger.Printf(
			"Custom Attribute %q is blank for %s",
			vmwb.BackupDateCAName,
			vmwb.Name,
		)
		return true
	}

	// TODO: Should we apply any special behavior for zero value backup date?
	// Should we treat this differently than if the VM is missing Custom
	// Attributes?
	// if vmwb.BackupDate.IsZero() {
	// 	return true
	// }

	if ExceedsAge(vmwb.BackupDate, vmwb.WarningAgeInDaysThreshold) &&
		!ExceedsAge(vmwb.BackupDate, vmwb.CriticalAgeInDaysThreshold) {
		return true
	}

	// TODO: What other criteria would indicate WARNING state?

	return false

}

// IsCriticalState indicates whether the CRITICAL threshold has been crossed.
//
// A CRITICAL state is NOT returned if a Virtual Machine is missing the
// expected Custom Attribute for tracking last backup date. Instead, the
// caller is expected to also validate the WARNING state which IS expected to
// handle that scenario.
func (vmwb VMWithBackup) IsCriticalState() bool {

	// TODO: Should we apply any special behavior for zero value backup date?
	// Should we treat this differently than if the VM is missing Custom
	// Attributes?
	// if vmwb.BackupDate.IsZero() {
	// 	return true
	// }

	// if ExceedsAge(vmwb.BackupDate, vmwb.CriticalAgeInDaysThreshold) {
	// 	return true
	// }

	return ExceedsAge(vmwb.BackupDate, vmwb.CriticalAgeInDaysThreshold)

	// TODO: What other criteria would indicate CRITICAL state?

	// return false

}

// IsOKState indicates whether the WARNING or CRITICAL thresholds have been
// crossed.
func (vmwb VMWithBackup) IsOKState() bool {
	if vmwb.IsCriticalState() || vmwb.IsWarningState() {
		return false
	}

	return true

}

// HasBackup indicates whether a Virtual Machine has a user-specified Custom
// Attribute used to track last backup date with a non-empty value. This
// method does not validate the format of the Custom Attribute value, only
// that the requested value exists. This method does not consider whether the
// optional metadata Custom Attribute is present.
func (vmwb VMWithBackup) HasBackup() bool {
	backupDateVal, exists := vmwb.CustomAttributes[vmwb.BackupDateCAName]
	if exists && strings.TrimSpace(backupDateVal) != "" {
		return true
	}

	return false
}

// HasOldBackup indicates whether a Virtual Machine (with a recorded backup)
// has a backup date which exceeds a user-specified age threshold. If a backup
// is not present false is returned. For best results, the caller should first
// filter by the HasBackup() method for the most reliable result.
func (vmwb VMWithBackup) HasOldBackup() bool {
	if !vmwb.HasBackup() {
		return false
	}

	return ExceedsAge(vmwb.BackupDate, vmwb.WarningAgeInDaysThreshold)
}

// FormattedBackupAge returns the formatted age of a Virtual Machine's backup
// date.
func (vmwb VMWithBackup) FormattedBackupAge() string {
	return FormattedTimeSinceEvent(vmwb.BackupDate)
}

// BackupDaysAgo returns the age of a Virtual Machine's backup date in days.
func (vmwb VMWithBackup) BackupDaysAgo() int {
	return DaysAgo(vmwb.BackupDate)
}

// IsWarningState indicates whether any Virtual Machines in the collection
// have crossed the WARNING threshold.
func (vmswb VMsWithBackup) IsWarningState() bool {
	for _, vm := range vmswb {
		if vm.IsWarningState() {
			return true
		}
	}

	return false
}

// IsCriticalState indicates whether any Virtual Machines in the collection
// have crossed the CRITICAL threshold.
func (vmswb VMsWithBackup) IsCriticalState() bool {
	for _, vm := range vmswb {
		if vm.IsCriticalState() {
			return true
		}
	}

	return false
}

// IsOKState indicates whether all Virtual Machines in the collection have not
// crossed either of the WARNING or CRITICAL thresholds.
func (vmswb VMsWithBackup) IsOKState() bool {
	if vmswb.IsCriticalState() || vmswb.IsWarningState() {
		return false
	}

	return true
}

// AllHasBackup indicates whether all Virtual Machines in the collection have
// a last backup date.
//
// This method does not validate the format of the last backup date value,
// only that the requested value exists. This method does not consider whether
// the optional metadata Custom Attribute is present.
func (vmswb VMsWithBackup) AllHasBackup() bool {
	for _, vm := range vmswb {
		if !vm.HasBackup() {
			return false
		}
	}

	return true
}

// HasOldBackup indicates whether ANY of the Virtual Machines in the
// collection have a backup date which exceeds a user-specified age threshold.
// Because Virtual Machines without a recorded backup date are ignored, the
// caller should also use the AllHasBackup() method to first ensure that all
// Virtual Machines in the collection have a recorded backup.
func (vmswb VMsWithBackup) HasOldBackup() bool {
	for _, vm := range vmswb {
		if vm.HasOldBackup() {
			return true
		}
	}

	return false
}

// NumBackups returns the number of VirtualMachines in the collection which
// have a recorded backup via user specified Custom Attribute. This method
// does not validate the format of the Custom Attribute value, only that the
// requested value exists. This method does not consider whether the optional
// metadata Custom Attribute is present.
func (vmswb VMsWithBackup) NumBackups() int {
	var num int
	for _, vm := range vmswb {
		if vm.HasBackup() {
			num++
		}
	}

	return num
}

// NumOldBackups returns the number of VirtualMachines in the collection which
// have a recorded backup via user specified Custom Attribute older than
// specified thresholds.
func (vmswb VMsWithBackup) NumOldBackups() int {
	var num int
	for _, vm := range vmswb {
		if vm.HasBackup() && vm.HasOldBackup() {
			num++
		}
	}

	return num
}

// NumWithoutBackups returns the number of VirtualMachines in the collection
// which do not have a recorded backup via user specified Custom Attribute.
func (vmswb VMsWithBackup) NumWithoutBackups() int {
	var num int
	for _, vm := range vmswb {
		if !vm.HasBackup() {
			num++
		}
	}

	return num
}

// VMWithOldestBackup returns a pointer to a Virtual Machine from the
// collection with the oldest backup date. Any Virtual Machine without a
// recorded backup date is ignored. Nil is returned if the collection is empty
// or if there are no backups to evaluate.
//
// TODO: Is there a better approach to this? Just let the caller range over
// the collection to determine this?
func (vmswb VMsWithBackup) VMWithOldestBackup() *VMWithBackup {

	// Explicitly handle potential empty collection
	if len(vmswb) == 0 {
		return nil
	}

	// NOTE: Using pointers instead of value to allow us to determine when
	// this variable becomes set to the VM with the oldest backup date.
	var vmWithOldestBackup *VMWithBackup

	// Check the collection for a Virtual Machine with an older backup date
	// and if found, record it instead.
	for i := range vmswb {
		if vmswb[i].HasBackup() {
			if vmWithOldestBackup == nil {
				vmWithOldestBackup = &vmswb[i]
			}

			if vmswb[i].BackupDate.Before(vmWithOldestBackup.BackupDate) {
				vmWithOldestBackup = &vmswb[i]
			}
		}
	}

	return vmWithOldestBackup
}

// VMWithYoungestBackup returns a pointer to a Virtual Machine from the
// collection with the most recent backup date. Any Virtual Machine without a
// recorded backup date is ignored. Nil is returned if the collection is empty
// or if there are no backups to evaluate.
//
// TODO: Is there a better approach to this? Just let the caller range over
// the collection to determine this?
func (vmswb VMsWithBackup) VMWithYoungestBackup() *VMWithBackup {

	// Explicitly handle potential empty collection
	if len(vmswb) == 0 {
		return nil
	}

	// NOTE: Using pointers instead of value to allow us to determine when
	// this variable becomes set to the VM with the youngest backup date.
	var vmWithYoungestBackup *VMWithBackup

	// Check the collection for a Virtual Machine with an older backup date
	// and if found, record it instead.
	for i := range vmswb {
		if vmswb[i].HasBackup() {
			if vmWithYoungestBackup == nil {
				vmWithYoungestBackup = &vmswb[i]
			}

			if vmWithYoungestBackup.BackupDate.Before(vmswb[i].BackupDate) {
				vmWithYoungestBackup = &vmswb[i]
			}
		}
	}

	return vmWithYoungestBackup
}

// VMNames returns a list of sorted VirtualMachine names which have exceeded
// specified power cycle uptime thresholds. VirtualMachines which have yet to
// exceed specified thresholds are not listed.
func (vpcs VirtualMachinePowerCycleUptimeStatus) VMNames() string {

	funcTimeStart := time.Now()

	vmNames := make([]string, 0, len(vpcs.VMsCritical)+len(vpcs.VMsWarning))

	defer func(names *[]string) {
		logger.Printf(
			"It took %v to execute VMNames func (and retrieve %d Datastores).\n",
			time.Since(funcTimeStart),
			len(*names),
		)
	}(&vmNames)

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

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute TopTenOK func.\n",
			time.Since(funcTimeStart),
		)
	}()

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

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute BottomTenOK func.\n",
			time.Since(funcTimeStart),
		)
	}()

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

		logger.Printf(
			"Retrieving VirtualMachines from object %q of type %q",
			obj.Name,
			obj.Self.Type,
		)

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

// GetVMsWithCA receives a collection of VirtualMachines, a Custom Attribute
// name to filter VirtualMachines by and a boolean flag indicating whether
// VirtualMachines missing a Custom Attribute should be ignored. A collection
// of VMWithCA is returned along with an error (if applicable).
func GetVMsWithCA(vms []mo.VirtualMachine, vmCustomAttributeName string, ignoreMissingCA bool) ([]VMWithCA, error) {

	funcTimeStart := time.Now()

	vmsWithCA := make([]VMWithCA, 0, len(vms))

	defer func(vms *[]VMWithCA) {
		logger.Printf(
			"It took %v to execute GetVMsWithCA func (and retrieve %d VMWithCA).\n",
			time.Since(funcTimeStart),
			len(*vms),
		)
	}(&vmsWithCA)

	for _, vm := range vms {
		ca, err := GetObjectCustomAttribute(vm.ManagedEntity, vmCustomAttributeName, ignoreMissingCA)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to retrieve custom attribute for VM %s: %w",
				vm.Name,
				err,
			)
		}
		vmsWithCA = append(vmsWithCA, VMWithCA{
			VirtualMachine:  vm,
			CustomAttribute: ca,
		})
	}

	return vmsWithCA, nil

}

// GetVMsWithCAs receives a collection of VirtualMachines and returns a
// collection of sorted VirtualMachines. Each VirtualMachine includes an index
// of Custom Attributes for that VirtualMachine. If no Custom Attributes are
// available, the index is empty. An error is returned if one occurs.
func GetVMsWithCAs(vms []mo.VirtualMachine) ([]VMWithCAs, error) {

	funcTimeStart := time.Now()

	vmsWithAllCAs := make([]VMWithCAs, 0, len(vms))

	defer func(vms *[]VMWithCAs) {
		logger.Printf(
			"It took %v to execute GetVMsWithAllCAs func (and retrieve %d VMWithAllCAs).\n",
			time.Since(funcTimeStart),
			len(*vms),
		)
	}(&vmsWithAllCAs)

	for _, vm := range vms {
		customAttributes, err := GetObjectCustomAttributes(vm.ManagedEntity)
		switch {

		// Custom Attributes are not set for this object, though they exist in
		// the vSphere inventory as an attribute that could be set if the
		// vsphere admin wishes to do so.
		case errors.Is(err, ErrCustomAttributeNotSet):

			logger.Printf("Custom Attributes for virtual machine %q missing",
				vm.Name,
			)

			logger.Printf(
				"Adding VM %s to collection with empty Custom Attributes map",
				vm.Name,
			)
			vmsWithAllCAs = append(vmsWithAllCAs, VMWithCAs{
				VirtualMachine:   vm,
				CustomAttributes: make(CustomAttributes),
			})

		// Custom attributes are set, but some other error occurred
		case err != nil:
			return nil, fmt.Errorf(
				"failed to retrieve custom attributes for %s: %w",
				vm.Name,
				err,
			)

		// Custom attributes are set and successfully retrieved
		default:
			vmsWithAllCAs = append(vmsWithAllCAs, VMWithCAs{
				VirtualMachine:   vm,
				CustomAttributes: customAttributes,
			})
		}

	}

	sort.Slice(vmsWithAllCAs, func(i, j int) bool {
		return strings.ToLower(vmsWithAllCAs[i].Name) < strings.ToLower(vmsWithAllCAs[j].Name)
	})

	return vmsWithAllCAs, nil
}

// GetVMsWithBackup receives a collection of VirtualMachines, a user-specified
// time zone (i.e., "location"), a Custom Attribute name for the last backup
// (required), a Custom Attribute name for the last backup's metadata
// (optional), thresholds for when the backup should be considered in a
// CRITICAL or WARNING state and whether missing Custom Attributes should be
// ignored.
//
// An error is returned if the given empty collection of VirtualMachines is
// empty or the user specified time zone is not recognized.
func GetVMsWithBackup(
	vms []mo.VirtualMachine,
	backupTimezone string,
	lastBackupCA string,
	backupMetadataCA string,
	backupDateFormat string,
	criticalAgeThreshold int,
	warningAgeThreshold int,
	ignoreMissingCAs bool,
) (VMsWithBackup, error) {

	funcTimeStart := time.Now()

	vmsWithBackup := make(VMsWithBackup, 0, len(vms))

	defer func(vms *VMsWithBackup) {
		logger.Printf(
			"It took %v to execute GetVMsWithBackup func (and retrieve %d VMWithBackup).\n",
			time.Since(funcTimeStart),
			len(*vms),
		)
	}(&vmsWithBackup)

	if len(vms) == 0 {
		return nil, fmt.Errorf(
			"received empty collection of virtual machines to evaluate for backup details",
		)
	}

	vmsWithCAs, err := GetVMsWithCAs(vms)
	if err != nil {
		// TODO: Anything additional needed here?
		return nil, err
	}

	// TODO: Any valid reason for an empty list of VMWithCAs?
	// if len(vmsWithCAs) == 0 {
	// 	return nil, fmt.Errorf(
	// 		"failed to retrieve list of virtual machines with custom attributes",
	// 	)
	// }

	for _, vm := range vmsWithCAs {

		// Rely on a map's zero value behavior when an element is not present.
		backupDateCAVal := vm.CustomAttributes[lastBackupCA]

		var backupDateParsed time.Time
		if backupDateCAVal != "" {
			location, err := time.LoadLocation(backupTimezone)
			switch {
			case err != nil:
				return nil, fmt.Errorf(
					"error loading location data using user specified time zone of %q: %w",
					backupTimezone,
					err,
				)

			default:
				// We were able to retrieve the location, so use it when
				// attempting to parse the recorded backup date.
				var err error
				backupDateParsed, err = time.ParseInLocation(backupDateFormat, backupDateCAVal, location)
				if err != nil {
					return nil, fmt.Errorf(
						"error evaluating backup date for virtual machine %q: %w",
						vm.Name,
						err,
					)
				}
			}
		}

		vmsWithBackup = append(
			vmsWithBackup,
			VMWithBackup{
				VMWithCAs:                  vm,
				BackupDateCAName:           lastBackupCA,
				BackupMetadataCAName:       backupMetadataCA,
				BackupDate:                 backupDateParsed,
				WarningAgeInDaysThreshold:  warningAgeThreshold,
				CriticalAgeInDaysThreshold: criticalAgeThreshold,
			},
		)
	}

	return vmsWithBackup, nil

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

// FilterVMsByCustomAttributeNames searches a given collection of VMWithCAs
// using provided Custom Attribute names and returns a collection of VMWithCAs
// which have values for those Custom Attributes. Validation of Custom
// Attribute values is not performed. If the caller does not wish to require
// that a specific Custom Attribute be present, the caller should not specify
// the value in the given list of Custom Attribute names.
//
// If the collection of provided VirtualMachines is empty, an empty collection
// is returned. The collection is returned along with the number of
// VirtualMachines that were excluded due to missing Custom Attribute names.
func FilterVMsByCustomAttributeNames(vmsWithCAs []VMWithCAs, caNames []string) ([]VMWithCAs, int) {

	// setup early so we can reference it from deferred stats output
	matchedVMs := make([]VMWithCAs, 0, len(vmsWithCAs))

	funcTimeStart := time.Now()

	defer func(vms []VMWithCAs, filteredVMs *[]VMWithCAs) {
		logger.Printf(
			"It took %v to execute FilterVMsByCustomAttributeNames func (for %d VMs, yielding %d VMs).\n",
			time.Since(funcTimeStart),
			len(vms),
			len(*filteredVMs),
		)
	}(vmsWithCAs, &matchedVMs)

	for _, vm := range vmsWithCAs {

		var missingAttribute bool

		for _, caName := range caNames {
			if _, ok := vm.CustomAttributes[caName]; !ok {
				// VM does not have set value for specified Custom Attribute
				missingAttribute = true
				continue
			}
		}

		if !missingAttribute {
			matchedVMs = append(matchedVMs, vm)
		}
	}

	numVMsExcluded := len(vmsWithCAs) - len(matchedVMs)

	// Sort collection
	sort.Slice(matchedVMs, func(i, j int) bool {
		return matchedVMs[i].Name < matchedVMs[j].Name
	})

	return matchedVMs, numVMsExcluded

}

// FilterVMsByCustomAttributeStatus accepts a collection of VirtualMachines
// and evaluates whether Custom Attributes are recorded. If the collection of
// provided VirtualMachines is empty, an empty collection is returned. The
// collection is returned along with the number of VirtualMachines that were
// excluded due to missing Custom Attributes.
func FilterVMsByCustomAttributeStatus(vms []VMWithCAs) ([]VMWithCAs, int) {

	// setup early so we can reference it from deferred stats output
	vmsWithCAs := make([]VMWithCAs, 0, len(vms))

	funcTimeStart := time.Now()

	defer func(vms []VMWithCAs, filteredVMs *[]VMWithCAs) {
		logger.Printf(
			"It took %v to execute FilterVMsByCustomAttributeStatus func (for %d VMs, yielding %d VMs).\n",
			time.Since(funcTimeStart),
			len(vms),
			len(*filteredVMs),
		)
	}(vms, &vmsWithCAs)

	for _, vm := range vms {
		if len(vm.CustomAttributes) > 0 {
			vmsWithCAs = append(vmsWithCAs, vm)
		}
	}

	numVMsExcluded := len(vms) - len(vmsWithCAs)

	return vmsWithCAs, numVMsExcluded

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

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute ExcludeVMsByName func.\n",
			time.Since(funcTimeStart),
		)
	}()

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

// FilterVMsByDiskConsolidationState accepts a collection of VirtualMachines
// and evaluates whether their ConsolidationNeeded flag is set. This function
// assumes that the caller has already initiated a "reload" of each
// VirtualMachine in order to retrieve the most current status of its
// ConsolidationNeeded field. If the collection of provided VirtualMachines is
// empty, an empty collection is returned. The collection is returned along
// with the number of VirtualMachines that were excluded.
func FilterVMsByDiskConsolidationState(vms []mo.VirtualMachine) ([]mo.VirtualMachine, int) {

	// setup early so we can reference it from deferred stats output
	var vmsNeedingConsolidation []mo.VirtualMachine

	funcTimeStart := time.Now()

	defer func(vms []mo.VirtualMachine, filteredVMs *[]mo.VirtualMachine) {
		logger.Printf(
			"It took %v to execute FilterVMsByDiskConsolidationState func (for %d VMs, yielding %d VMs).\n",
			time.Since(funcTimeStart),
			len(vms),
			len(*filteredVMs),
		)
	}(vms, &vmsNeedingConsolidation)

	for _, vm := range vms {
		if vm.Runtime.ConsolidationNeeded != nil && *vm.Runtime.ConsolidationNeeded {
			vmsNeedingConsolidation = append(vmsNeedingConsolidation, vm)
		}
	}

	numExcluded := len(vms) - len(vmsNeedingConsolidation)

	return vmsNeedingConsolidation, numExcluded

}

// FilterVMsByInteractiveQuestionStatus accepts a collection of
// VirtualMachines and evaluates whether their Question flag is set. If the
// collection of provided VirtualMachines is empty, an empty collection is
// returned. The collection is returned along with the number of
// VirtualMachines that were excluded.
func FilterVMsByInteractiveQuestionStatus(vms []mo.VirtualMachine) ([]mo.VirtualMachine, int) {

	// setup early so we can reference it from deferred stats output
	var vmsWaitingOnInput []mo.VirtualMachine

	funcTimeStart := time.Now()

	defer func(vms []mo.VirtualMachine, filteredVMs *[]mo.VirtualMachine) {
		logger.Printf(
			"It took %v to execute FilterVMsByInteractiveQuestionStatus func (for %d VMs, yielding %d VMs).\n",
			time.Since(funcTimeStart),
			len(vms),
			len(*filteredVMs),
		)
	}(vms, &vmsWaitingOnInput)

	for _, vm := range vms {
		if vm.Summary.Runtime.Question != nil {
			vmsWaitingOnInput = append(vmsWaitingOnInput, vm)
		}
	}

	numExcluded := len(vms) - len(vmsWaitingOnInput)

	return vmsWaitingOnInput, numExcluded

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

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute VMNames func.\n",
			time.Since(funcTimeStart),
		)
	}()

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
			case vm.Summary.Runtime.Question != nil &&
				vm.Summary.Runtime.Question.Text != "":
				question = vm.Summary.Runtime.Question.Text
			default:
				question = "unknown"
			}

			possibleAnswers := make([]string, 0, len(vm.Summary.Runtime.Question.Choice.ChoiceInfo))
			for _, e := range vm.Summary.Runtime.Question.Choice.ChoiceInfo {
				ed := e.(*types.ElementDescription)
				// possibleAnswers = append(possibleAnswers, fmt.Sprintf(
				// 	"'%s > %s'",
				// 	ed.Key, ed.Description.Label,
				// ))
				possibleAnswers = append(possibleAnswers, fmt.Sprintf(
					"'%s'",
					ed.Description.Label,
				))
			}

			fmt.Fprintf(
				&report,
				"* %s (%q [%s])%s",
				vm.Name,
				question,
				strings.Join(possibleAnswers, ", "),
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

// VMBackupViaCAOneLineCheckSummary is used to generate a one-line
// Nagios service check results summary. This is the line most prominent in
// notifications.
func VMBackupViaCAOneLineCheckSummary(
	stateLabel string,
	allVMs []mo.VirtualMachine,
	evaluatedVMs []mo.VirtualMachine,
	vmsWithBackups VMsWithBackup,
	rps []mo.ResourcePool,
) string {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute VMBackupViaCAOneLineCheckSummary func.\n",
			time.Since(funcTimeStart),
		)
	}()

	numMissingBackups := vmsWithBackups.NumWithoutBackups()
	numWithBackups := vmsWithBackups.NumBackups()
	numWithOldBackups := vmsWithBackups.NumOldBackups()
	numCurrentBackups := numWithBackups - numWithOldBackups
	if numCurrentBackups < 0 {
		numCurrentBackups = 0
	}

	switch {
	case numWithOldBackups > 0:
		numCurrentBackups := numWithBackups - numWithOldBackups
		if numCurrentBackups < 0 {
			numCurrentBackups = 0
		}
		return fmt.Sprintf(
			"%s: %d VMs with old backups detected (%d current, evaluated %d VMs & %d Resource Pools)",
			stateLabel,
			numWithOldBackups,
			numCurrentBackups,
			len(evaluatedVMs),
			len(rps),
		)

	case numMissingBackups > 0:
		return fmt.Sprintf(
			"%s: %d VMs missing backups detected (%d present & %d current, evaluated %d VMs & %d Resource Pools)",
			stateLabel,
			numMissingBackups,
			numWithBackups,
			numCurrentBackups,
			len(evaluatedVMs),
			len(rps),
		)

	default:

		return fmt.Sprintf(
			"%s: No VMs with old or missing backups detected (%d present, evaluated %d VMs & %d Resource Pools)",
			stateLabel,
			numCurrentBackups,
			len(evaluatedVMs),
			len(rps),
		)
	}

}

// VMBackupViaCAReport generates a summary of VMs & their backup status along
// with various verbose details intended to aid in troubleshooting check
// results at a glance. This information is provided for use with the Long
// Service Output field commonly displayed on the detailed service check
// results display in the web UI or in the body of many notifications.
func VMBackupViaCAReport(
	c *vim25.Client,
	allVMs []mo.VirtualMachine,
	evaluatedVMs []mo.VirtualMachine,
	vmsWithBackup VMsWithBackup,
	vmsToExclude []string,
	includeRPs []string,
	excludeRPs []string,
	rps []mo.ResourcePool,
) string {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute VMBackupViaCAReport func.\n",
			time.Since(funcTimeStart),
		)
	}()

	rpNames := make([]string, len(rps))
	for i := range rps {
		rpNames[i] = rps[i].Name
	}

	var report strings.Builder

	// Somewhat arbitrary number intended to limit the number of VMs emitted
	// in order to keep output manageable.
	vmPrintLimit := 50

	printVM := func(w io.Writer, vm VMWithBackup) {

		backupDateCAVal := vm.CustomAttributes[vm.BackupDateCAName]

		fmt.Fprintf(
			w,
			"* %s%s",
			vm.Name,
			nagios.CheckOutputEOL,
		)

		fmt.Fprintf(
			w,
			"\t** %s: %t%s",
			"Old Backup",
			vm.HasOldBackup(),
			nagios.CheckOutputEOL,
		)

		fmt.Fprintf(
			w,
			"\t** %s: %s%s",
			"Backup age (formatted)",
			vm.FormattedBackupAge(),
			nagios.CheckOutputEOL,
		)

		fmt.Fprintf(
			w,
			"\t** %s: %d%s",
			"Backup age (in days)",
			vm.BackupDaysAgo(),
			nagios.CheckOutputEOL,
		)

		fmt.Fprintf(
			w,
			"\t** %s (raw value): %q%s",
			vm.BackupDateCAName,
			backupDateCAVal,
			nagios.CheckOutputEOL,
		)
		if vm.BackupMetadataCAName != "" {
			fmt.Fprintf(
				w,
				"\t** %s: %q%s",
				vm.BackupMetadataCAName,
				vm.CustomAttributes[vm.BackupMetadataCAName],
				nagios.CheckOutputEOL,
			)
		}

	}

	fmt.Fprintf(
		&report,
		"VMs without backups:%s%s",
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
	)
	switch {
	case vmsWithBackup.NumWithoutBackups() == 0:
		fmt.Fprintf(&report, "* None%s", nagios.CheckOutputEOL)

	case vmsWithBackup.NumWithoutBackups() > vmPrintLimit:
		fmt.Fprintf(
			&report,
			"* %d VMs without backups; output limit of %d reached, omitting list of VMs%s",
			vmsWithBackup.NumWithoutBackups(),
			vmPrintLimit,
			nagios.CheckOutputEOL,
		)
	default:
		for _, vm := range vmsWithBackup {
			if !vm.HasBackup() {
				printVM(&report, vm)
			}
		}
	}
	fmt.Fprint(&report, nagios.CheckOutputEOL)

	fmt.Fprintf(
		&report,
		"VMs with old backups: %s%s",
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
	)
	switch {
	case vmsWithBackup.NumOldBackups() == 0:
		fmt.Fprintf(&report, "* None%s", nagios.CheckOutputEOL)

	case vmsWithBackup.NumOldBackups() > vmPrintLimit:
		fmt.Fprintf(
			&report,
			"* %d VMs with old backups; output limit of %d reached, omitting list of VMs%s",
			vmsWithBackup.NumOldBackups(),
			vmPrintLimit,
			nagios.CheckOutputEOL,
		)

	default:
		for _, vm := range vmsWithBackup {
			if vm.HasBackup() && vm.HasOldBackup() {
				printVM(&report, vm)
			}
		}
	}
	fmt.Fprint(&report, nagios.CheckOutputEOL)

	fmt.Fprintf(
		&report,
		"Virtual Machines Backup Summary: %s%s",
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		&report,
		"* Missing Backups: %d%s",
		vmsWithBackup.NumWithoutBackups(),
		nagios.CheckOutputEOL,
	)
	fmt.Fprintf(
		&report,
		"* With Backups: %d%s",
		vmsWithBackup.NumBackups(),
		nagios.CheckOutputEOL,
	)
	fmt.Fprintf(
		&report,
		"* Old Backups: %d%s",
		vmsWithBackup.NumOldBackups(),
		nagios.CheckOutputEOL,
	)

	vmWithOldestBackup := vmsWithBackup.VMWithOldestBackup()
	if vmWithOldestBackup != nil {
		fmt.Fprintf(
			&report,
			"* Oldest backup: %s (%s)%s",
			vmWithOldestBackup.Name,
			vmWithOldestBackup.FormattedBackupAge(),
			nagios.CheckOutputEOL,
		)
	}

	vmWithYoungestBackup := vmsWithBackup.VMWithYoungestBackup()
	if vmWithYoungestBackup != nil {
		fmt.Fprintf(
			&report,
			"* Most recent backup: %s (%s)%s",
			vmWithYoungestBackup.Name,
			vmWithYoungestBackup.FormattedBackupAge(),
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
		"* VMs (total: %d, evaluated: %d)%s",
		len(allVMs),
		len(evaluatedVMs),
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		&report,
		"* Powered off VMs evaluated: %t%s",
		// NOTE: This plugin is used to detect the backup status for Virtual
		// Machines, regardless of power state; we report here that both
		// powered on and powered off VMs are evaluated for simplicity.
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
