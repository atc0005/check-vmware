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
	"github.com/vmware/govmomi/property"
	"github.com/vmware/govmomi/view"
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

// ErrValidationOfIncludeExcludeRPLists indicates that a validation attempt of
// the given ResourcePool include or exclude lists failed.
var ErrValidationOfIncludeExcludeRPLists = errors.New("validation failed for include/exclude resource pool lists")

// ErrValidationOfIncludeExcludeFolderIDLists indicates that a validation
// attempt of the given Folder include or exclude lists failed.
var ErrValidationOfIncludeExcludeFolderIDLists = errors.New("validation failed for include/exclude folder ID lists")

// ErrVirtualMachineConfigurationIsNil indicates that the configuration for a
// virtual machine is unset, which may occur if the property is not requested
// from the vSphere API or if the service account executing the plugin has
// insufficient privileges.
var ErrVirtualMachineConfigurationIsNil = errors.New("virtual machine configuration is nil")

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

	// BackupDateCAName is the name (not the value) of the Custom Attribute
	// which indicates when the last backup occurred for this VirtualMachine.
	BackupDateCAName string

	// BackupMetadataCAName is the name (not the value) of the Custom
	// Attribute which provides additional context for the last backup for
	// this VirtualMachine.
	BackupMetadataCAName string

	// BackupDate is the date/time of the last backup for this VirtualMachine.
	// If a backup date is recorded for a VM, then the time zone (aka,
	// "location") for the parsed date/time value is set to the user-specified
	// time zone. If a time zone is not specified, the default location is
	// used for a recorded backup date.
	BackupDate *time.Time

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

// VMsFilterOptions is the set of options used to filter a given collection of
// VirtualMachines.
type VMsFilterOptions struct {
	ResourcePoolsIncluded       []string
	ResourcePoolsExcluded       []string
	FoldersIncluded             []string
	FoldersExcluded             []string
	VirtualMachineNamesExcluded []string
	IncludePoweredOff           bool
}

// vmsRPFilterResults is the results of performing resource pool filtering
// operations using specified filter settings.
type vmsRPFilterResults struct {
	RPs []mo.ResourcePool
	VMs []mo.VirtualMachine
}

// vmsFolderFilterResults is the results of performing folder filtering
// operations on a given VirtualMachines collection.
type vmsFolderFilterResults struct {
	VMs                    []mo.VirtualMachine
	NumVMsExcludedByFolder int
	NumFoldersEvaluated    int
}

// VMsFilterResults is the results of performing filtering operations on a
// given VirtualMachines collection.
type VMsFilterResults struct {
	// numVMsAll is the count of all vms in the inventory which are not
	// templates.
	numVMsAll int

	// numVMsExcludedByName is the count of vms excluded or "filtered out" via
	// name filtering.
	numVMsExcludedByName int

	// numVMsExcludedByResourcePool is the count of vms excluded or "filtered
	// out" via resource pool filtering.
	numVMsExcludedByResourcePool int

	// numVMsExcludedByFolder is the count of vms excluded or "filtered out"
	// via folder filtering.
	numVMsExcludedByFolder int

	// numVMsExcludedByPowerState is the count of vms excluded or "filtered
	// out" via power state filtering.
	numVMsExcludedByPowerState int

	// numFoldersAll is the count of all folders in the inventory.
	numFoldersAll int

	// NumFoldersEvaluated is the number of folders remaining after they have
	// been explicitly included or excluded.
	numFoldersEvaluated int

	// NumFoldersExcluded is the count of folders explicitly specified via CLI
	// flag for inclusion.
	numFoldersIncluded int

	// NumFoldersExcluded is the count of folders explicitly specified via CLI
	// flag for exclusion.
	numFoldersExcluded int

	// numResourcePoolsAll is the count of all resource pools in the
	// inventory.
	numResourcePoolsAll int

	// numResourcePoolsEvaluated is the number of resource pools remaining
	// after they have been explicitly included or excluded.
	numResourcePoolsEvaluated int

	// numResourcePoolsIncluded is the count of resource pools explicitly
	// specified via CLI flag for inclusion.
	numResourcePoolsIncluded int

	// numResourcePoolsExcluded is the count of resource pools explicitly
	// specified via CLI flag for exclusion.
	numResourcePoolsExcluded int

	// rps is the collection of resource pools remaining after they have been
	// explicitly included or excluded.
	rpsAfterAllFiltering []mo.ResourcePool

	vmsAfterRPFiltering         []mo.VirtualMachine
	vmsAfterFolderFiltering     []mo.VirtualMachine
	vmsAfterVMNameFiltering     []mo.VirtualMachine
	vmsAfterPowerStateFiltering []mo.VirtualMachine
	vmsAfterAllFiltering        []mo.VirtualMachine
}

// IsWarningState indicates whether the WARNING threshold has been crossed or
// if the Virtual Machine is missing a backup.
func (vmwb VMWithBackup) IsWarningState() bool {
	if !vmwb.HasBackup() {
		return true
	}

	if ExceedsAge(*vmwb.BackupDate, vmwb.WarningAgeInDaysThreshold) &&
		!ExceedsAge(*vmwb.BackupDate, vmwb.CriticalAgeInDaysThreshold) {
		return true
	}

	return false
}

// IsCriticalState indicates whether the CRITICAL threshold has been crossed.
//
// A CRITICAL state is NOT if a Virtual Machine is missing a backup. Instead,
// the caller is expected to also validate the WARNING state which IS expected
// to handle that scenario.
func (vmwb VMWithBackup) IsCriticalState() bool {
	switch {
	case !vmwb.HasBackup():
		// Reminder: The IsWarningState() method is responsible for handling
		// this scenario.
		return false

	default:
		return ExceedsAge(*vmwb.BackupDate, vmwb.CriticalAgeInDaysThreshold)
	}
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
// Attribute used to track last backup date with a non-empty value and that
// the BackupDate field is non-nil (assumed to be set when retrieving the last
// backup date for a VM).
//
// This method does not validate the format of the Custom Attribute value,
// only that the requested value exists. This method does not consider whether
// the optional metadata Custom Attribute is present.
func (vmwb VMWithBackup) HasBackup() bool {
	backupDateVal, backupDateValExists := vmwb.CustomAttributes[vmwb.BackupDateCAName]

	// NOTE: This is an optional Custom Attribute, so we don't require it here.
	// _, hasBackupMetadataCA := vmwb.CustomAttributes[vmwb.BackupMetadataCAName]

	switch {
	case !backupDateValExists:
		logger.Printf(
			"Custom Attribute %q missing from %s",
			vmwb.BackupDateCAName,
			vmwb.Name,
		)
		return false

	case strings.TrimSpace(backupDateVal) == "":

		logger.Printf(
			"Custom Attribute %q is blank for %s",
			vmwb.BackupDateCAName,
			vmwb.Name,
		)
		return false

	case vmwb.BackupDate == nil:
		logger.Printf(
			"No backup date is recorded for %s",
			vmwb.Name,
		)
		return false

	default:
		logger.Printf(
			"No problems with recorded backup date detected for %s; "+
				"assuming valid backup date",
			vmwb.Name,
		)
		return true
	}
}

// HasOldBackup indicates whether a Virtual Machine (with a recorded backup)
// has a backup date which exceeds a user-specified age threshold. If a backup
// is not present false is returned. For best results, the caller should first
// filter by the HasBackup() method for the most reliable result.
func (vmwb VMWithBackup) HasOldBackup() bool {
	switch {
	case !vmwb.HasBackup():
		return false

	default:
		return ExceedsAge(*vmwb.BackupDate, vmwb.WarningAgeInDaysThreshold)
	}
}

// FormattedBackupAge returns the formatted age of a Virtual Machine's backup
// date. If a backup is not recorded, this is indicated instead of a formatted
// age.
func (vmwb VMWithBackup) FormattedBackupAge() string {
	switch {
	case vmwb.BackupDate == nil:
		return "Backup unavailable"
	default:
		return FormattedTimeSinceEvent(*vmwb.BackupDate)
	}
}

// BackupDaysAgo returns the age of a Virtual Machine's backup date in days.
// If a backup date is not available, 0 is returned.
func (vmwb VMWithBackup) BackupDaysAgo() int {
	switch {
	case vmwb.BackupDate == nil:
		return 0
	default:
		return DaysAgo(*vmwb.BackupDate)
	}
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

			if vmswb[i].BackupDate.Before(*vmWithOldestBackup.BackupDate) {
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

			if vmWithYoungestBackup.BackupDate.Before(*vmswb[i].BackupDate) {
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

// CountVMsPowerStateOn returns the count of VMs from the provided collection
// that are powered on.
func CountVMsPowerStateOn(vms []mo.VirtualMachine) int {
	var count int
	for _, vm := range vms {
		if vm.Runtime.PowerState == types.VirtualMachinePowerStatePoweredOn {
			count++
		}
	}

	return count
}

// CountVMsPowerStateSuspended returns the count of VMs from the provided
// collection that are suspended.
func CountVMsPowerStateSuspended(vms []mo.VirtualMachine) int {
	var count int
	for _, vm := range vms {
		if vm.Runtime.PowerState == types.VirtualMachinePowerStateSuspended {
			count++
		}
	}

	return count
}

// CountVMsPowerStateOff returns the count of VMs from the provided collection
// that are fully powered off. This count does not include VMs that are
// suspended.
func CountVMsPowerStateOff(vms []mo.VirtualMachine) int {
	var count int
	for _, vm := range vms {
		if vm.Runtime.PowerState == types.VirtualMachinePowerStatePoweredOff {
			count++
		}
	}

	return count
}

// CountVMsPoweredOff returns the count of VMs from the provided collection
// that are fully powered off and those which are suspended.
func CountVMsPoweredOff(vms []mo.VirtualMachine) int {
	var count int
	for _, vm := range vms {
		if vm.Runtime.PowerState == types.VirtualMachinePowerStatePoweredOff ||
			vm.Runtime.PowerState == types.VirtualMachinePowerStateSuspended {
			count++
		}
	}

	return count
}

// CountVMsPowerStates returns the count of VMs from the provided collection
// in each power state.
//
// The order of returned values:
//
//  1. Powered On
//  2. Suspended
//  3. Powered Off
func CountVMsPowerStates(vms []mo.VirtualMachine) (int, int, int) {
	var countPowerStateOn int
	var countPowerStateSuspended int
	var countPowerStateOff int

	for _, vm := range vms {
		switch {
		case vm.Runtime.PowerState != types.VirtualMachinePowerStatePoweredOn:
			countPowerStateOn++

		case vm.Runtime.PowerState != types.VirtualMachinePowerStateSuspended:
			countPowerStateSuspended++

		case vm.Runtime.PowerState != types.VirtualMachinePowerStatePoweredOff:
			countPowerStateOff++
		}
	}

	return countPowerStateOn, countPowerStateSuspended, countPowerStateOff
}

// FilterVMs is used as a high-level abstraction to handle the filtering of
// VirtualMachines in the inventory given a set of filtering options. The
// results of the filtering steps are returned as an aggregate or an error if
// one occurs.
//
//	Filter order:
//
//	1. Resource Pools
//	2. Folder
//	3. VirtualMachine Name
//	4. VirtualMachine Power State
//
// Separate filtering functions are provided for a more fine-tuned, manual
// approach to filtering VirtualMachines.
func FilterVMs(ctx context.Context, client *vim25.Client, filterOptions VMsFilterOptions) (VMsFilterResults, error) {
	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute FilterVMs func.\n",
			time.Since(funcTimeStart),
		)
	}()

	if err := validateRPs(ctx, client, filterOptions); err != nil {
		return VMsFilterResults{}, err
	}

	if err := validateFolders(ctx, client, filterOptions); err != nil {
		return VMsFilterResults{}, err
	}

	numAllRPs, rpsCountErr := GetNumTotalRPs(ctx, client)
	if rpsCountErr != nil {
		return VMsFilterResults{}, rpsCountErr
	}

	numFolders, foldersCountErr := GetNumTotalFolders(ctx, client)
	if foldersCountErr != nil {
		return VMsFilterResults{}, foldersCountErr
	}

	numNonTemplateVMs, vmsCountErr := GetNumTotalVMs(ctx, client)
	if vmsCountErr != nil {
		return VMsFilterResults{}, vmsCountErr
	}

	logger.Println("Filtering VMs by resource pool")
	vmsRPResults, rpsFilterErr := filterVMsByRP(ctx, client, filterOptions)
	if rpsFilterErr != nil {
		return VMsFilterResults{}, rpsFilterErr
	}

	logger.Println("Filtering VMs by folder")
	vmsFolderResults, folderFilterErr := filterVMsByFolder(
		ctx, client, vmsRPResults.VMs, filterOptions,
	)
	if folderFilterErr != nil {
		return VMsFilterResults{}, folderFilterErr
	}

	logger.Println("Filtering VMs by name")
	vmsAfterNameFiltering, numVMsExcludedByName := ExcludeVMsByName(vmsFolderResults.VMs, filterOptions.VirtualMachineNamesExcluded)
	logger.Printf(
		"VMs after name filtering: (filteredByName: %v, excludedByName: %d)",
		strings.Join(VMNames(vmsAfterNameFiltering), ", "),
		numVMsExcludedByName,
	)

	logger.Println("Filtering VMs by specified power state")
	vmsAfterPowerStateFiltering, numVMsExcludedByPowerState := FilterVMsByPowerState(vmsAfterNameFiltering, filterOptions.IncludePoweredOff)
	logger.Printf(
		"VMs after power state filtering: (filteredByPowerState: %v, excludedByPowerState: %d)",
		strings.Join(VMNames(vmsAfterPowerStateFiltering), ", "),
		numVMsExcludedByPowerState,
	)

	return VMsFilterResults{
		numVMsAll:                    numNonTemplateVMs,
		numVMsExcludedByResourcePool: numNonTemplateVMs - len(vmsRPResults.VMs),
		numVMsExcludedByFolder:       vmsFolderResults.NumVMsExcludedByFolder,
		numVMsExcludedByName:         numVMsExcludedByName,
		numVMsExcludedByPowerState:   numVMsExcludedByPowerState,

		numFoldersAll:       numFolders,
		numFoldersEvaluated: vmsFolderResults.NumFoldersEvaluated,
		numFoldersIncluded:  len(filterOptions.FoldersIncluded),
		numFoldersExcluded:  len(filterOptions.FoldersExcluded),

		numResourcePoolsAll:       numAllRPs,
		numResourcePoolsEvaluated: len(vmsRPResults.RPs),
		numResourcePoolsIncluded:  len(filterOptions.FoldersIncluded),
		numResourcePoolsExcluded:  len(filterOptions.FoldersExcluded),

		vmsAfterRPFiltering:         vmsRPResults.VMs,
		vmsAfterFolderFiltering:     vmsFolderResults.VMs,
		vmsAfterVMNameFiltering:     vmsAfterNameFiltering,
		vmsAfterPowerStateFiltering: vmsAfterPowerStateFiltering,
		vmsAfterAllFiltering:        vmsAfterPowerStateFiltering,
		rpsAfterAllFiltering:        vmsRPResults.RPs,
	}, nil

}

// GetNumTotalVMs provides the total number of non-template VirtualMachines
// in the inventory.
func GetNumTotalVMs(ctx context.Context, client *vim25.Client) (int, error) {
	funcTimeStart := time.Now()

	var numAllVMs int

	defer func(allVMs *int) {
		logger.Printf(
			"It took %v to execute GetNumTotalVMs func (and count %d VMs).\n",
			time.Since(funcTimeStart),
			*allVMs,
		)
	}(&numAllVMs)

	var getVMsErr error
	numAllVMs, getVMsErr = getVMsCountUsingRootFolderContainerView(
		ctx,
		client,
		true,
	)

	if getVMsErr != nil {
		logger.Printf(
			"error retrieving count of all virtual machines and templates: %v",
			getVMsErr,
		)

		return 0, fmt.Errorf(
			"error retrieving count of all virtual machines and templates: %w",
			getVMsErr,
		)
	}

	return numAllVMs, nil
}

// GetAllVMs provides every VirtualMachine in the inventory using the
// RootFolder as the starting point. In contrast to retrieving VirtualMachine
// values from ResourcePools, this function also returns template
// VirtualMachines.
func GetAllVMs(ctx context.Context, client *vim25.Client) ([]mo.VirtualMachine, error) {
	funcTimeStart := time.Now()

	var allVMs []mo.VirtualMachine

	defer func(vms *[]mo.VirtualMachine) {
		logger.Printf(
			"It took %v to execute GetAllVMs func (and retrieve %d VMs).\n",
			time.Since(funcTimeStart),
			len(*vms),
		)
	}(&allVMs)

	err := getObjects(
		ctx,
		client,
		&allVMs,
		client.ServiceContent.RootFolder,
		true,
		true,
	)
	if err != nil {
		logger.Printf(
			"error retrieving all virtual machines: %v",
			err,
		)

		return nil, fmt.Errorf(
			"error retrieving all virtual machines: %w",
			err,
		)
	}

	return allVMs, nil
}

// filterVMsByRP uses the given filtering options to obtain VirtualMachines
// from eligible resource pools.
func filterVMsByRP(
	ctx context.Context,
	client *vim25.Client,
	filterOptions VMsFilterOptions,
) (vmsRPFilterResults, error) {
	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute filterVMsByRP func.\n",
			time.Since(funcTimeStart),
		)
	}()

	logger.Println("Retrieving eligible resource pools")
	filteredResourcePools, getRPsErr := GetEligibleRPs(
		ctx,
		client,
		filterOptions.ResourcePoolsIncluded,
		filterOptions.ResourcePoolsExcluded,
		true,
	)
	if getRPsErr != nil {
		logger.Printf(
			"Error retrieving list of resource pools: %v",
			getRPsErr,
		)

		return vmsRPFilterResults{}, fmt.Errorf(
			"failed to retrieve list of resource pools: %w",
			getRPsErr,
		)
	}
	logger.Println("Finished retrieving eligible resource pools")

	rpNames := make([]string, 0, len(filteredResourcePools))
	for _, rp := range filteredResourcePools {
		rpNames = append(rpNames, rp.Name)
	}

	logger.Printf("Resource Pools: %v", strings.Join(rpNames, ", "))

	logger.Println("Retrieving VMs from eligible resource pools")
	rpEntityVals := make([]mo.ManagedEntity, 0, len(filteredResourcePools))
	for i := range filteredResourcePools {
		rpEntityVals = append(rpEntityVals, filteredResourcePools[i].ManagedEntity)
	}

	vmsFromRPs, getVMsErr := GetVMsFromContainer(ctx, client, true, rpEntityVals...)
	if getVMsErr != nil {
		logger.Printf(
			"Error retrieving list of VMs from resource pools list: %v",
			getVMsErr,
		)

		return vmsRPFilterResults{}, fmt.Errorf(
			"failed to retrieve VMs from resource pools: %w",
			getVMsErr,
		)
	}
	logger.Printf(
		"Finished retrieving %d vms from %d resource pools",
		len(vmsFromRPs),
		len(filteredResourcePools),
	)

	logger.Printf("VMs to evaluate: %v", strings.Join(VMNames(vmsFromRPs), ", "))

	return vmsRPFilterResults{
		RPs: filteredResourcePools,
		VMs: vmsFromRPs,
	}, nil
}

func filterVMsByFolder(
	ctx context.Context,
	client *vim25.Client,
	vms []mo.VirtualMachine,
	filterOptions VMsFilterOptions,
) (vmsFolderFilterResults, error) {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute filterVMsByFolder func.\n",
			time.Since(funcTimeStart),
		)
	}()

	switch {
	case len(filterOptions.FoldersIncluded) > 0 || len(filterOptions.FoldersExcluded) > 0:
		// If the exclude list is specified, only grab VMs from the excluded
		// folders and use with ExcludeVMsByVMs.
		//
		// If the include list is specified, grab VMs from the included
		// folders and use with FilterVMsByVMs.
		//
		// We use these "sift" variables to reflect which of the two
		// approaches we're using.
		var (
			siftList            []string
			numFoldersEvaluated int
			siftListDesc        string
			keepMatchedVMs      bool
		)

		switch {
		case len(filterOptions.FoldersIncluded) > 0:
			siftList = filterOptions.FoldersIncluded
			siftListDesc = "included folders list"
			numFoldersEvaluated = len(filterOptions.FoldersIncluded)
			keepMatchedVMs = true
		case len(filterOptions.FoldersExcluded) > 0:
			siftList = filterOptions.FoldersExcluded
			siftListDesc = "excluded folders list"
			numFoldersEvaluated = len(filterOptions.FoldersExcluded)
			keepMatchedVMs = false
		}

		logger.Println("Resolving folder IDs to folder values")
		folders, retrieveErr := GetFoldersByIDs(ctx, client, siftList, true)
		if retrieveErr != nil {
			logger.Printf(
				"Error retrieving %s: %v",
				siftListDesc,
				retrieveErr,
			)
			return vmsFolderFilterResults{}, fmt.Errorf(
				"failed to retrieve %s: %w",
				siftListDesc,
				retrieveErr,
			)
		}

		folderEntityVals := FolderManagedEntityVals(folders)

		logger.Printf("Retrieving vms from %d folders", len(folderEntityVals))
		vmsFromFolders, getVMsErr := GetVMsFromContainer(ctx, client, true, folderEntityVals...)
		if getVMsErr != nil {
			logger.Printf(
				"Error retrieving list of VMs from %s: %v",
				siftListDesc,
				getVMsErr,
			)

			return vmsFolderFilterResults{}, fmt.Errorf(
				"failed to retrieve list of VMs from %s: %w",
				siftListDesc,
				getVMsErr,
			)
		}
		logger.Printf(
			"Finished retrieving %d vms from %s",
			len(vmsFromFolders),
			siftListDesc,
		)

		logger.Printf(
			"Filtering %d given VMs against %d VMs retrieved from %d folders",
			len(vms),
			len(vmsFromFolders),
			len(folderEntityVals),
		)
		filteredVMs, numVMsExcludedByFolder := SiftVMsByVMs(vms, vmsFromFolders, keepMatchedVMs)

		logger.Printf(
			"VMs after folder filtering: %v (kept: %d, excluded: %d)",
			strings.Join(VMNames(filteredVMs), ", "),
			len(filteredVMs),
			numVMsExcludedByFolder,
		)

		return vmsFolderFilterResults{
			VMs:                    filteredVMs,
			NumFoldersEvaluated:    numFoldersEvaluated,
			NumVMsExcludedByFolder: numVMsExcludedByFolder,
		}, nil

	default:
		logger.Println("Skipping filter by folder; folder filtering not requested")

		// Return the original collection untouched.
		return vmsFolderFilterResults{
			VMs:                    vms,
			NumFoldersEvaluated:    0,
			NumVMsExcludedByFolder: 0,
		}, nil
	}
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

	err := getObjects(ctx, c, &vms, c.ServiceContent.RootFolder, propsSubset, true)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve VirtualMachines: %w", err)
	}

	sort.Slice(vms, func(i, j int) bool {
		return strings.ToLower(vms[i].Name) < strings.ToLower(vms[j].Name)
	})

	return vms, nil
}

// getVMsCountUsingParentPoolsAndContainerView returns only VM counts, not VM
// template counts.
//
// This function retrieves all VirtualMachines using "parent" resource pools
// named 'Resources' (ParentResourcePool). One of these pools are present for
// each standalone host (ComputeResource) in a datacenter.
//
// For example, if we have one standalone host and one cluster in a datacenter
// like so:
//
// Datacenter (ExampleDC)
// - Hosts & Clusters (hostFolder property)
//   - HostSystem (192.168.2.200)
//   - Cluster (XYZ-Hosted)
//
// There are two 'Resources' resource pools.
//
// We count VMs within each 'Resources' resource pool. This approach ignores
// VM templates as templates are not rooted to a resource pool.
// func getVMsCountUsingParentPoolsAndContainerView(
// 	ctx context.Context,
// 	c *vim25.Client,
// 	recursive bool,
// ) (int, error) {
// 	funcTimeStart := time.Now()
//
// 	var allVMs []types.ObjectContent
//
// 	defer func(vms *[]types.ObjectContent) {
// 		logger.Printf(
// 			"It took %v to execute getVMsCountUsingParentPoolsAndContainerView func (and count %d VMs).\n",
// 			time.Since(funcTimeStart),
// 			len(allVMs),
// 		)
// 	}(&allVMs)
//
// 	// Create a view of caller-specified objects
// 	m := view.NewManager(c)
//
// 	// FIXME: Should this filter to a specific datacenter? See GH-219.
// 	rpView, createViewErr := m.CreateContainerView(
// 		ctx,
// 		c.ServiceContent.RootFolder,
// 		[]string{MgObjRefTypeResourcePool},
// 		recursive,
// 	)
// 	if createViewErr != nil {
// 		return 0, createViewErr
// 	}
//
// 	defer func() {
// 		// Per vSphere Web Services SDK Programming Guide - VMware vSphere 7.0
// 		// Update 1:
// 		//
// 		// A best practice when using views is to call the DestroyView()
// 		// method when a view is no longer needed. This practice frees memory
// 		// on the server.
// 		if err := rpView.Destroy(ctx); err != nil {
// 			logger.Printf("Error occurred while destroying datacenter view: %s", err)
// 		}
// 	}()
//
// 	var rps []mo.ResourcePool
// 	retrieveErr := rpView.Retrieve(ctx, []string{MgObjRefTypeResourcePool}, []string{"name"}, &rps)
// 	if retrieveErr != nil {
// 		return 0, retrieveErr
// 	}
//
// 	for _, rp := range rps {
// 		if rp.Name == ParentResourcePool {
// 			parentPoolID := rp.Reference()
// 			logger.Printf("Pool ID %s for Pool %s", parentPoolID, rp.Name)
// 			vmView, createViewErr := m.CreateContainerView(
// 				ctx,
// 				parentPoolID,
// 				[]string{MgObjRefTypeVirtualMachine},
// 				recursive,
// 			)
// 			if createViewErr != nil {
// 				return 0, createViewErr
// 			}
// 			defer func() {
// 				if err := vmView.Destroy(ctx); err != nil {
// 					logger.Printf("Error occurred while destroying virtual machine view: %s", err)
// 				}
// 			}()
//
// 			var poolVMs []types.ObjectContent
//
// 			kind := []string{MgObjRefTypeVirtualMachine}
// 			prop := []string{"overallStatus"}
// 			retrieveErr = vmView.Retrieve(ctx, kind, prop, &poolVMs)
// 			if retrieveErr != nil {
// 				return 0, fmt.Errorf(
// 					"failed to retrieve VMs list from parent resource pool %s: %w",
// 					ParentResourcePool,
// 					retrieveErr,
// 				)
// 			}
//
// 			allVMs = append(allVMs, poolVMs...)
// 		}
// 	}
//
// 	// Remove any duplicate entries which could occur if we process nested
// 	// resource pools. Because we are processing only "parent" resource pools
// 	// this should not happen, but we guard against any unexpected edge cases
// 	// just to be sure.
// 	logger.Println("Deduplicating VMs")
// 	numVMsBeforeDeduping := len(allVMs)
// 	allVMs = dedupeObjects(allVMs)
// 	numVMsAfterDeduping := len(allVMs)
//
// 	logger.Printf("Before deduping VMs: %d", numVMsBeforeDeduping)
// 	logger.Printf("After deduping VMs: %d", numVMsAfterDeduping)
//
// 	return len(allVMs), nil
// }

// getVMsCountUsingRootFolderContainerView returns only VM counts, not VM
// template counts.
func getVMsCountUsingRootFolderContainerView(
	ctx context.Context,
	c *vim25.Client,
	recursive bool,
) (int, error) {

	funcTimeStart := time.Now()

	var allVMs []types.ObjectContent

	defer func(vms *[]types.ObjectContent) {
		logger.Printf(
			"It took %v to execute getVMsCountUsingRootFolderContainerView func (and count %d VMs).\n",
			time.Since(funcTimeStart),
			len(*vms),
		)
	}(&allVMs)

	// Create a view of caller-specified objects
	m := view.NewManager(c)

	kind := []string{MgObjRefTypeVirtualMachine}
	v, createViewErr := m.CreateContainerView(
		ctx,
		c.ServiceContent.RootFolder,
		kind,
		recursive,
	)
	if createViewErr != nil {
		return 0, createViewErr
	}
	defer func() {
		if err := v.Destroy(ctx); err != nil {
			logger.Printf("Error occurred while destroying virtual machine view: %s", err)
		}
	}()

	filter := property.Match{"config.template": false}
	prop := []string{"overallStatus"}

	retrieveErr := v.RetrieveWithFilter(ctx, kind, prop, &allVMs, filter)
	if retrieveErr != nil {
		return 0, fmt.Errorf(
			"failed to retrieve VMs list: %w",
			retrieveErr,
		)
	}

	return len(allVMs), nil
}

// GetVMsFromContainer receives one or many ManagedEntity values for Folder,
// Datacenter, ComputeResource, ResourcePool, VirtualApp or HostSystem types
// and returns a list of VirtualMachine object references. Deduplication of
// VirtualMachines is applied in order to properly handle nested resource
// pools.
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
			"It took %v to execute GetVMsFromContainer func (and retrieve %d VMs).\n",
			time.Since(funcTimeStart),
			len(*vms),
		)
	}(&allVMs)

	for _, obj := range objs {

		var vmsFromContainer []mo.VirtualMachine

		// Perform a recursive retrieval by default.
		recursiveRetrieval := true
		if obj.Name == ParentResourcePool &&
			obj.Self.Type == MgObjRefTypeResourcePool {

			// If we're retrieving Virtual Machines from the parent Resource
			// Pool perform a shallow retrieval instead so that we do not pull
			// in any Virtual Machines that should be ignored due to a
			// Resource Pool exclusion.
			recursiveRetrieval = false
		}

		logger.Printf(
			"Retrieving VirtualMachines (recursively: %t) from object %q of type %q ",
			recursiveRetrieval,
			obj.Name,
			obj.Self.Type,
		)

		err := getObjects(ctx, c, &vmsFromContainer, obj.Reference(), propsSubset, recursiveRetrieval)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to retrieve VirtualMachines from object: %s: %w",
				obj.Name,
				err,
			)
		}

		allVMs = append(allVMs, vmsFromContainer...)

	}

	// Remove any duplicate entries which could occur if we process nested
	// resource pools.
	logger.Println("Deduplicating VMs")
	numVMsBeforeDeduping := len(allVMs)
	allVMs = dedupeVMs(allVMs)
	numVMsAfterDeduping := len(allVMs)

	logger.Printf("Before deduping VMs: %d", numVMsBeforeDeduping)
	logger.Printf("After deduping VMs: %d", numVMsAfterDeduping)

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
	err := getObjects(ctx, c, &allVMs, c.ServiceContent.RootFolder, propsSubset, true)
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
				"failed to retrieve custom attribute for %s %s: %w",
				vm.ManagedEntity.Self.Type,
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
			"It took %v to execute GetVMsWithCAs func (and retrieve %d VMWithAllCAs).\n",
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

			logger.Printf("Custom attributes for virtual machine %q missing",
				vm.Name,
			)

			logger.Printf(
				"Adding VM %s to collection with empty custom attributes map",
				vm.Name,
			)
			vmsWithAllCAs = append(vmsWithAllCAs, VMWithCAs{
				VirtualMachine:   vm,
				CustomAttributes: make(CustomAttributes),
			})

		// Custom attributes are set, but some other error occurred
		case err != nil:
			return nil, fmt.Errorf(
				"failed to retrieve custom attribute for %s %s: %w",
				vm.ManagedEntity.Self.Type,
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

		vmWithBackup := VMWithBackup{
			VMWithCAs:                  vm,
			BackupDateCAName:           lastBackupCA,
			BackupMetadataCAName:       backupMetadataCA,
			WarningAgeInDaysThreshold:  warningAgeThreshold,
			CriticalAgeInDaysThreshold: criticalAgeThreshold,
		}

		// If we managed to parse the backup date, record it for the VM,
		// otherwise leave the pointer field nil as a signal that no backup
		// date is available for the VM.
		if !backupDateParsed.IsZero() {
			vmWithBackup.BackupDate = &backupDateParsed
		}

		vmsWithBackup = append(vmsWithBackup, vmWithBackup)

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

// SiftVMsByVMs accepts an original collection of VMs and a collection to
// match against. If specified, the collection of matches are returned,
// otherwise VMs not matched are returned. An error is returned if one occurs.
func SiftVMsByVMs(vmsToExamine []mo.VirtualMachine, vmsToMatch []mo.VirtualMachine, keepMatches bool) ([]mo.VirtualMachine, int) {
	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute SiftVMsByVMs func.\n",
			time.Since(funcTimeStart),
		)
	}()

	if keepMatches {
		return FilterVMsByVMs(vmsToExamine, vmsToMatch)
	}
	return ExcludeVMsByVMs(vmsToExamine, vmsToMatch)
}

// FilterVMsByVMs receives a collection of VirtualMachines to examine and
// another collection of VirtualMachines to filter against. VirtualMachines
// from the first collection present in the second collection are returned. If
// the collection to review or the collection to filter against is empty then
// an empty collection is returned. The number of excluded VirtualMachines (if
// any) is also returned.
func FilterVMsByVMs(vmsToExamine []mo.VirtualMachine, vmsOkToKeep []mo.VirtualMachine) ([]mo.VirtualMachine, int) {
	funcTimeStart := time.Now()

	var numExcluded int
	vmsToKeep := make([]mo.VirtualMachine, 0, len(vmsToExamine))

	defer func() {
		logger.Printf(
			"It took %v to execute FilterVMsByVMs func.\n",
			time.Since(funcTimeStart),
		)
	}()

	// We require populated source and filter collections.
	switch {
	case len(vmsToExamine) == 0:
		// if nothing to examine, nothing to exclude
		return []mo.VirtualMachine{}, 0

	case len(vmsOkToKeep) == 0:
		// if nothing to keep, then everything to review is to be excluded
		return []mo.VirtualMachine{}, len(vmsToExamine)
	}

	vmIDsToKeep := make(map[string]struct{}, len(vmsOkToKeep))
	for _, vm := range vmsOkToKeep {
		vmIDsToKeep[vm.Self.Value] = struct{}{}
	}

	for _, vmToReview := range vmsToExamine {
		if _, ok := vmIDsToKeep[vmToReview.Self.Value]; ok {
			vmsToKeep = append(vmsToKeep, vmToReview)
			continue
		}
		numExcluded++
	}

	return vmsToKeep, numExcluded
}

// ExcludeVMsByVMs receives a collection of VirtualMachines to examine and
// another collection of VirtualMachines to match against. VirtualMachines
// from the first collection NOT present in the second collection are
// returned. If the collection to review is empty then an empty collection is
// returned. If the collection to match against is empty then the first
// collection is returned unmodified.
//
// The number of excluded VirtualMachines (if any) is also returned.
func ExcludeVMsByVMs(vmsToExamine []mo.VirtualMachine, vmsToExclude []mo.VirtualMachine) ([]mo.VirtualMachine, int) {
	funcTimeStart := time.Now()

	var numExcluded int
	vmsToKeep := make([]mo.VirtualMachine, 0, len(vmsToExamine))

	defer func() {
		logger.Printf(
			"It took %v to execute ExcludeVMsByVMs func.\n",
			time.Since(funcTimeStart),
		)
	}()

	// We require populated source and filter collections.
	switch {
	case len(vmsToExamine) == 0:
		// if nothing to examine, nothing to exclude
		return []mo.VirtualMachine{}, 0

	case len(vmsToExclude) == 0:
		// if nothing to exclude, then the original collection is returned
		return vmsToExamine, 0
	}

	vmIDsToExclude := make(map[string]struct{}, len(vmsToExclude))
	for _, vm := range vmsToExclude {
		vmIDsToExclude[vm.Self.Value] = struct{}{}
	}

	for _, vmToReview := range vmsToExamine {
		if _, ok := vmIDsToExclude[vmToReview.Self.Value]; ok {
			numExcluded++
			continue
		}
		vmsToKeep = append(vmsToKeep, vmToReview)
	}

	return vmsToKeep, numExcluded
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
// VirtualMachines to just those with WARNING or CRITICAL values. The
// collection is returned along with the number of VirtualMachines that were
// excluded.
func FilterVMsByPowerCycleUptime(vms []mo.VirtualMachine, warningThreshold int) ([]mo.VirtualMachine, int) {

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

// dedupeObjects receives a collection of ObjectContent data objects and
// returns unique values from the collection.
//
// NOTE: The ObjectContent data object type contains the contents retrieved
// for a single managed object.
//
// https://vdc-download.vmware.com/vmwb-repository/dcr-public/fa5d1ee7-fad5-4ebf-b150-bdcef1d38d35/a5e46da1-9b96-4f0c-a1d0-7b8f3ebfd4f5/doc/vmodl.query.PropertyCollector.ObjectContent.html
// func dedupeObjects(objects []types.ObjectContent) []types.ObjectContent {
//
// 	funcTimeStart := time.Now()
//
// 	defer func(objs *[]types.ObjectContent) {
// 		logger.Printf(
// 			"It took %v to execute dedupeObjects func (evaluated %d ObjectContent values).\n",
// 			time.Since(funcTimeStart),
// 			len(*objs),
// 		)
// 	}(&objects)
//
// 	seen := make(map[string]struct{}, len(objects))
// 	j := 0
// 	for _, item := range objects {
// 		if _, ok := seen[item.Obj.Value]; ok {
// 			continue
// 		}
// 		seen[item.Obj.Value] = struct{}{}
// 		objects[j] = item
// 		j++
// 	}
//
// 	return objects[:j]
// }

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

// getVMHostID retrieves the VM host MOID value for a specified VM. If the
// host MOID value is not available an error is returned. The host MOID value
// may be unavailable if the service account executing the plugin does not
// have sufficient permissions.
func getVMHostID(vm mo.VirtualMachine) (string, error) {
	switch {

	case vm.Runtime.Host == nil:
		return "", fmt.Errorf(
			"error retrieving associated Host MOID for VM %s: %w",
			vm.Name,
			ErrManagedObjectIDIsNil,
		)

	case vm.Runtime.Host.Value == "":
		return "", fmt.Errorf(
			"error retrieving associated Host MOID for VM %s: %w",
			vm.Name,
			ErrManagedObjectIDIsEmpty,
		)

	default:

		// Safe to reference now that we have guarded against potential
		// nil Host field pointer and empty MOID.

		return vm.Runtime.Host.Value, nil

	}
}

// NumVMsAll is the count of all VirtualMachines in the inventory.
func (vfr VMsFilterResults) NumVMsAll() int {
	return vfr.numVMsAll
}

// NumVMsExcluded is the count of all VirtualMachines excluded by filtering
// operations.
func (vfr VMsFilterResults) NumVMsExcluded() int {
	return vfr.numVMsAll - len(vfr.vmsAfterPowerStateFiltering)
}

// NumVMsAfterFiltering is the count of all VirtualMachines after filtering
// was applied.
func (vfr VMsFilterResults) NumVMsAfterFiltering() int {
	return len(vfr.vmsAfterAllFiltering)
}

// NumVMsExcludedByName is the count of all VirtualMachines excluded by name
// filtering.
func (vfr VMsFilterResults) NumVMsExcludedByName() int {
	return vfr.numVMsExcludedByName
}

// NumVMsExcludedByPowerState is the count of all VirtualMachines excluded by
// power state filtering.
func (vfr VMsFilterResults) NumVMsExcludedByPowerState() int {
	return vfr.numVMsExcludedByPowerState
}

// NumVMsExcludedByResourcePool is the count of all VirtualMachines excluded
// by resource pool filtering.
func (vfr VMsFilterResults) NumVMsExcludedByResourcePool() int {
	return vfr.numVMsExcludedByResourcePool
}

// NumVMsExcludedByFolder is the count of all VirtualMachines excluded by
// folder filtering.
func (vfr VMsFilterResults) NumVMsExcludedByFolder() int {
	return vfr.numVMsExcludedByFolder
}

// NumFoldersAll is the count of all Folders in the inventory.
func (vfr VMsFilterResults) NumFoldersAll() int {
	return vfr.numFoldersAll
}

// NumFoldersIncluded is the count of all ResourcePools excluded by filtering
// operations.
func (vfr VMsFilterResults) NumFoldersIncluded() int {
	return vfr.numFoldersIncluded
}

// NumFoldersExcluded is the count of all ResourcePools excluded by filtering
// operations.
func (vfr VMsFilterResults) NumFoldersExcluded() int {
	return vfr.numFoldersExcluded
}

// NumFoldersAfterFiltering is the count of all ResourcePools after filtering was
// applied.
func (vfr VMsFilterResults) NumFoldersAfterFiltering() int {
	return vfr.numFoldersEvaluated
}

// NumRPsAll is the count of all ResourcePools in the inventory.
func (vfr VMsFilterResults) NumRPsAll() int {
	return vfr.numResourcePoolsAll
}

// NumRPsIncluded is the count of all ResourcePools excluded by filtering
// operations.
func (vfr VMsFilterResults) NumRPsIncluded() int {
	return vfr.numResourcePoolsIncluded
}

// NumRPsExcluded is the count of all ResourcePools excluded by filtering
// operations.
func (vfr VMsFilterResults) NumRPsExcluded() int {
	return vfr.numResourcePoolsExcluded
}

// NumRPsAfterFiltering is the count of all ResourcePools after filtering was
// applied.
func (vfr VMsFilterResults) NumRPsAfterFiltering() int {
	return vfr.numResourcePoolsEvaluated
}

// VMsAfterFiltering is the collection of VirtualMachines after all filtering
// steps were applied.
func (vfr VMsFilterResults) VMsAfterFiltering() []mo.VirtualMachine {
	return vfr.vmsAfterAllFiltering
}

// VMsAfterResourcePoolFiltering is the collection of VirtualMachines after
// resource pool filtering was applied.
func (vfr VMsFilterResults) VMsAfterResourcePoolFiltering() []mo.VirtualMachine {
	return vfr.vmsAfterRPFiltering
}

// VMsAfterFolderFiltering is the collection of VirtualMachines after
// folder filtering was applied.
func (vfr VMsFilterResults) VMsAfterFolderFiltering() []mo.VirtualMachine {
	return vfr.vmsAfterFolderFiltering
}

// VMsBeforeFolderFiltering is the collection of VirtualMachines before
// folder filtering was applied.
func (vfr VMsFilterResults) VMsBeforeFolderFiltering() []mo.VirtualMachine {
	return vfr.vmsAfterRPFiltering
}

// VMsAfterVMNameFiltering is the collection of VirtualMachines after
// VirtualMachine name filtering was applied.
func (vfr VMsFilterResults) VMsAfterVMNameFiltering() []mo.VirtualMachine {
	return vfr.vmsAfterVMNameFiltering
}

// VMsBeforeVMNameFiltering is the collection of VirtualMachines before
// VirtualMachine name filtering was applied.
func (vfr VMsFilterResults) VMsBeforeVMNameFiltering() []mo.VirtualMachine {
	return vfr.vmsAfterFolderFiltering
}

// VMsAfterPowerStateFiltering is the collection of VirtualMachines after
// VirtualMachine power state filtering was applied.
func (vfr VMsFilterResults) VMsAfterPowerStateFiltering() []mo.VirtualMachine {
	return vfr.vmsAfterPowerStateFiltering
}

// VMsBeforePowerStateFiltering is the collection of VirtualMachines before
// VirtualMachine power state filtering was applied.
func (vfr VMsFilterResults) VMsBeforePowerStateFiltering() []mo.VirtualMachine {
	return vfr.vmsAfterVMNameFiltering
}

// RPsAfterFiltering is the collection of ResourcePools after all filtering
// steps were applied.
func (vfr VMsFilterResults) RPsAfterFiltering() []mo.ResourcePool {
	return vfr.rpsAfterAllFiltering
}

// VMNamesAfterFiltering is the collection of names for VirtualMachines
// remaining after all filtering steps were applied.
func (vfr VMsFilterResults) VMNamesAfterFiltering() []string {
	vmNames := make([]string, 0, vfr.NumVMsAfterFiltering())
	for _, vm := range vfr.VMsAfterFiltering() {
		vmNames = append(vmNames, vm.Name)
	}
	return vmNames
}

// RPNamesAfterFiltering is the collection of names for ResourcePools
// remaining after all filtering steps were applied.
func (vfr VMsFilterResults) RPNamesAfterFiltering() []string {
	rpNames := make([]string, 0, vfr.NumRPsAfterFiltering())
	for _, rp := range vfr.RPsAfterFiltering() {
		rpNames = append(rpNames, rp.Name)
	}
	return rpNames
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
	vmsFilterResults VMsFilterResults,
	uptimeSummary VirtualMachinePowerCycleUptimeStatus,

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
			vmsFilterResults.NumVMsAfterFiltering(),
			vmsFilterResults.NumRPsAfterFiltering(),
		)

	case len(uptimeSummary.VMsWarning) > 0:
		return fmt.Sprintf(
			"%s: %d VMs with power cycle uptime exceeding %d days detected (evaluated %d VMs, %d Resource Pools)",
			stateLabel,
			len(uptimeSummary.VMsWarning),
			uptimeSummary.WarningThreshold,
			vmsFilterResults.NumVMsAfterFiltering(),
			vmsFilterResults.NumRPsAfterFiltering(),
		)

	default:

		return fmt.Sprintf(
			"%s: No VMs with power cycle uptime exceeding %d days detected (evaluated %d VMs, %d Resource Pools)",
			stateLabel,
			uptimeSummary.WarningThreshold,
			vmsFilterResults.NumVMsAfterFiltering(),
			vmsFilterResults.NumRPsAfterFiltering(),
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
	vmsFilterOptions VMsFilterOptions,
	vmsFilterResults VMsFilterResults,
	uptimeSummary VirtualMachinePowerCycleUptimeStatus,

) string {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute VMPowerCycleUptimeReport func.\n",
			time.Since(funcTimeStart),
		)
	}()

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

	vmFilterResultsReportTrailer(
		&report,
		c,
		vmsFilterOptions,
		vmsFilterResults,
		true,
	)

	return report.String()
}

// VMDiskConsolidationOneLineCheckSummary is used to generate a one-line Nagios
// service check results summary. This is the line most prominent in
// notifications.
func VMDiskConsolidationOneLineCheckSummary(
	stateLabel string,
	vmsFilterResults VMsFilterResults,
	vmsNeedingConsolidation []mo.VirtualMachine,
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
			vmsFilterResults.NumVMsAfterFiltering(),
			vmsFilterResults.NumRPsAfterFiltering(),
		)

	default:

		return fmt.Sprintf(
			"%s: No VMs requiring disk consolidation detected (evaluated %d VMs, %d Resource Pools)",
			stateLabel,
			vmsFilterResults.NumVMsAfterFiltering(),
			vmsFilterResults.NumRPsAfterFiltering(),
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
	vmsFilterOptions VMsFilterOptions,
	vmsFilterResults VMsFilterResults,
	vmsNeedingConsolidation []mo.VirtualMachine,
) string {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute VMDiskConsolidationReport func.\n",
			time.Since(funcTimeStart),
		)
	}()

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
	vmFilterResultsReportTrailer(
		&report,
		c,
		vmsFilterOptions,
		vmsFilterResults,
		true,
	)

	return report.String()
}

// VMInteractiveQuestionOneLineCheckSummary is used to generate a one-line
// Nagios service check results summary. This is the line most prominent in
// notifications.
func VMInteractiveQuestionOneLineCheckSummary(
	stateLabel string,
	vmsFilterResults VMsFilterResults,
	vmsNeedingResponse []mo.VirtualMachine,
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
			vmsFilterResults.NumVMsAfterFiltering(),
			vmsFilterResults.NumRPsAfterFiltering(),
		)

	default:

		return fmt.Sprintf(
			"%s: No VMs requiring interactive response detected (evaluated %d VMs, %d Resource Pools)",
			stateLabel,
			vmsFilterResults.NumVMsAfterFiltering(),
			vmsFilterResults.NumRPsAfterFiltering(),
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
	vmsFilterOptions VMsFilterOptions,
	vmsFilterResults VMsFilterResults,
	vmsNeedingResponse []mo.VirtualMachine,
) string {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute VMInteractiveQuestionReport func.\n",
			time.Since(funcTimeStart),
		)
	}()

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

	vmFilterResultsReportTrailer(
		&report,
		c,
		vmsFilterOptions,
		vmsFilterResults,
		true,
	)

	return report.String()
}

// VMBackupViaCAOneLineCheckSummary is used to generate a one-line
// Nagios service check results summary. This is the line most prominent in
// notifications.
func VMBackupViaCAOneLineCheckSummary(
	stateLabel string,
	vmsFilterResults VMsFilterResults,
	vmsWithBackups VMsWithBackup,

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
		return fmt.Sprintf(
			"%s: %d VMs with old backups detected (%d current; evaluated %d VMs, %d Resource Pools)",
			stateLabel,
			numWithOldBackups,
			numCurrentBackups,
			vmsFilterResults.NumVMsAfterFiltering(),
			vmsFilterResults.NumRPsAfterFiltering(),
		)

	case numMissingBackups > 0:
		return fmt.Sprintf(
			"%s: %d VMs missing backups detected (%d present, %d current; evaluated %d VMs, %d Resource Pools)",
			stateLabel,
			numMissingBackups,
			numWithBackups,
			numCurrentBackups,
			vmsFilterResults.NumVMsAfterFiltering(),
			vmsFilterResults.NumRPsAfterFiltering(),
		)

	default:

		return fmt.Sprintf(
			"%s: No VMs with old or missing backups detected (%d present; evaluated %d VMs, %d Resource Pools)",
			stateLabel,
			numCurrentBackups,
			vmsFilterResults.NumVMsAfterFiltering(),
			vmsFilterResults.NumRPsAfterFiltering(),
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
	vmsFilterOptions VMsFilterOptions,
	vmsFilterResults VMsFilterResults,
	vmsWithBackup VMsWithBackup,
) string {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute VMBackupViaCAReport func.\n",
			time.Since(funcTimeStart),
		)
	}()

	var report strings.Builder

	// Somewhat arbitrary number intended to limit the number of VMs emitted
	// in order to keep output manageable.
	vmPrintLimit := 50

	printVM := func(w io.Writer, vm VMWithBackup) {

		backupDateCAVal := vm.CustomAttributes[vm.BackupDateCAName]
		backupDateMetadataVal := vm.CustomAttributes[vm.BackupMetadataCAName]

		fmt.Fprintf(
			w,
			"* %s%s",
			vm.Name,
			nagios.CheckOutputEOL,
		)

		if vm.HasBackup() {
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
		}

		if backupDateCAVal != "" {
			fmt.Fprintf(
				w,
				"\t** %s (raw value): %q%s",
				vm.BackupDateCAName,
				backupDateCAVal,
				nagios.CheckOutputEOL,
			)
		}

		if backupDateMetadataVal != "" {
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

	vmFilterResultsReportTrailer(
		&report,
		c,
		vmsFilterOptions,
		vmsFilterResults,
		true,
	)

	return report.String()
}

// VMListOneLineCheckSummary is used to generate a one-line Nagios service
// check results summary. This is the line most prominent in notifications.
func VMListOneLineCheckSummary(stateLabel string, vmsFilterResults VMsFilterResults) string {
	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute VMListOneLineCheckSummary func.\n",
			time.Since(funcTimeStart),
		)
	}()

	switch {
	case len(vmsFilterResults.VMsAfterFiltering()) > 0:
		return fmt.Sprintf(
			"%s: %d VMs remaining after filtering (evaluated %d of %d VMs, %d of %d Resource Pools)",
			stateLabel,
			vmsFilterResults.NumVMsAfterFiltering(),
			vmsFilterResults.NumVMsAfterFiltering(),
			vmsFilterResults.NumVMsAll(),
			vmsFilterResults.NumRPsAfterFiltering(),
			vmsFilterResults.NumRPsAll(),
		)

	default:
		return fmt.Sprintf(
			"%s: No VMs remaining after filtering (evaluated %d of %d VMs, %d of %d Resource Pools)",
			stateLabel,
			vmsFilterResults.NumVMsAfterFiltering(),
			vmsFilterResults.NumVMsAll(),
			vmsFilterResults.NumRPsAfterFiltering(),
			vmsFilterResults.NumRPsAll(),
		)
	}
}

// VMListReport generates a summary of VMs before filtering and after along
// with various verbose details intended to aid in troubleshooting check
// results at a glance. This information is provided for use with the Long
// Service Output field commonly displayed on the detailed service check
// results display in the web UI or in the body of many notifications.
func VMListReport(
	c *vim25.Client,
	vmsFilterOptions VMsFilterOptions,
	vmsFilterResults VMsFilterResults,
) string {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute VMListReport func.\n",
			time.Since(funcTimeStart),
		)
	}()

	var report strings.Builder

	fmt.Fprintf(
		&report,
		"Summary of inventory before before any filtering was applied:%s%s",
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		&report,
		"* %d Virtual Machines%s",
		vmsFilterResults.NumVMsAll(),
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		&report,
		"* %d Resource Pools%s",
		vmsFilterResults.NumRPsAll(),
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		&report,
		"* %d Folders%s",
		vmsFilterResults.NumFoldersAll(),
		nagios.CheckOutputEOL,
	)

	vmListReportFilteringBeforeAfterResults(
		&report,
		vmsFilterResults.NumVMsAll(),
		"",
		nil,
		"resource pool",
		nil,
		vmsFilterResults.VMsAfterResourcePoolFiltering,
	)

	fmt.Fprint(&report, nagios.CheckOutputEOL)

	vmListReportFilteringBeforeAfterResults(
		&report,
		vmsFilterResults.NumVMsAll(),
		"resource pool",
		vmsFilterResults.VMsAfterResourcePoolFiltering,
		"folder",
		vmsFilterResults.VMsBeforeFolderFiltering,
		vmsFilterResults.VMsAfterFolderFiltering,
	)

	fmt.Fprint(&report, nagios.CheckOutputEOL)

	vmListReportFilteringBeforeAfterResults(
		&report,
		vmsFilterResults.NumVMsAll(),
		"folder",
		vmsFilterResults.VMsAfterFolderFiltering,
		"VM name",
		vmsFilterResults.VMsBeforeVMNameFiltering,
		vmsFilterResults.VMsAfterVMNameFiltering,
	)

	fmt.Fprint(&report, nagios.CheckOutputEOL)

	vmListReportFilteringBeforeAfterResults(
		&report,
		vmsFilterResults.NumVMsAll(),
		"VM name",
		vmsFilterResults.VMsAfterVMNameFiltering,
		"VM power state",
		vmsFilterResults.VMsBeforePowerStateFiltering,
		vmsFilterResults.VMsAfterPowerStateFiltering,
	)

	fmt.Fprint(&report, nagios.CheckOutputEOL)

	vmListReportAfterAllFiltering(&report, vmsFilterResults)

	fmt.Fprint(&report, nagios.CheckOutputEOL)

	vmFilterResultsReportTrailer(
		&report,
		c,
		vmsFilterOptions,
		vmsFilterResults,
		true,
	)

	return report.String()
}

func vmListReportAfterAllFiltering(w io.Writer, vmsFilterResults VMsFilterResults) {
	fmt.Fprintf(
		w,
		"%s(%d of %d) VMs after all filtering was applied:%s%s",
		nagios.CheckOutputEOL,
		len(vmsFilterResults.VMsAfterFiltering()),
		vmsFilterResults.NumVMsAll(),
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
	)

	switch {
	case len(vmsFilterResults.VMsAfterFiltering()) == vmsFilterResults.NumVMsAll():
		fmt.Fprintf(
			w,
			"No filtering applied: %d VMs remain.%s",
			vmsFilterResults.NumVMsAll(),
			nagios.CheckOutputEOL,
		)

	case len(vmsFilterResults.VMsAfterFiltering()) == len(vmsFilterResults.VMsAfterResourcePoolFiltering()):
		fmt.Fprintf(
			w,
			"* Same list as after resource pool filtering.%s",
			nagios.CheckOutputEOL,
		)

	case len(vmsFilterResults.VMsAfterFiltering()) == len(vmsFilterResults.VMsAfterFolderFiltering()):
		fmt.Fprintf(
			w,
			"* Same list as after folder filtering.%s",
			nagios.CheckOutputEOL,
		)

	case len(vmsFilterResults.VMsAfterFiltering()) == len(vmsFilterResults.VMsAfterVMNameFiltering()):
		fmt.Fprintf(
			w,
			"* Same list as after VM name filtering.%s",
			nagios.CheckOutputEOL,
		)

	case len(vmsFilterResults.VMsAfterFiltering()) == len(vmsFilterResults.VMsAfterPowerStateFiltering()):
		fmt.Fprintf(
			w,
			"* Same list as after VM power state filtering.%s",
			nagios.CheckOutputEOL,
		)
	default:
		for _, vmName := range vmsFilterResults.VMNamesAfterFiltering() {
			fmt.Fprintf(
				w,
				"* %s%s",
				vmName,
				nagios.CheckOutputEOL,
			)
		}
	}
}

func vmFilterResultsReportTrailer(
	w io.Writer,
	c *vim25.Client,
	vmsFilterOptions VMsFilterOptions,
	vmsFilterResults VMsFilterResults,
	emitSeparator bool,
) {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute vmFilterResultsReportTrailer func.\n",
			time.Since(funcTimeStart),
		)
	}()

	if emitSeparator {
		fmt.Fprintf(
			w,
			"%s---%s%s",
			nagios.CheckOutputEOL,
			nagios.CheckOutputEOL,
			nagios.CheckOutputEOL,
		)
	}

	fmt.Fprintf(
		w,
		"* vSphere environment: %s%s",
		c.URL().String(),
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		w,
		"* Plugin User Agent: %s%s",
		c.Client.UserAgent,
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		w,
		"* VMs evaluated: %d of %d%s",
		len(vmsFilterResults.VMsAfterFiltering()),
		vmsFilterResults.NumVMsAll(),
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		w,
		"* Powered off VMs evaluated: %t%s",
		vmsFilterOptions.IncludePoweredOff,
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		w,
		"* Specified VMs to exclude (%d): [%v]%s",
		len(vmsFilterOptions.VirtualMachineNamesExcluded),
		strings.Join(vmsFilterOptions.VirtualMachineNamesExcluded, ", "),
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		w,
		"* Specified Folders to explicitly include (%d): [%v]%s",
		len(vmsFilterOptions.FoldersIncluded),
		strings.Join(vmsFilterOptions.FoldersIncluded, ", "),
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		w,
		"* Specified Folders to explicitly exclude (%d): [%v]%s",
		len(vmsFilterOptions.FoldersExcluded),
		strings.Join(vmsFilterOptions.FoldersExcluded, ", "),
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		w,
		"* Folders evaluated: %d of %d%s",
		vmsFilterResults.NumFoldersAfterFiltering(),
		vmsFilterResults.NumFoldersAll(),
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		w,
		"* Specified Resource Pools to explicitly include (%d): [%v]%s",
		len(vmsFilterOptions.ResourcePoolsIncluded),
		strings.Join(vmsFilterOptions.ResourcePoolsIncluded, ", "),
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		w,
		"* Specified Resource Pools to explicitly exclude (%d): [%v]%s",
		len(vmsFilterOptions.ResourcePoolsExcluded),
		strings.Join(vmsFilterOptions.ResourcePoolsExcluded, ", "),
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		w,
		"* Resource Pools evaluated (%d of %d): [%v]%s",
		vmsFilterResults.NumRPsAfterFiltering(),
		vmsFilterResults.NumRPsAll(),
		strings.Join(vmsFilterResults.RPNamesAfterFiltering(), ", "),
		nagios.CheckOutputEOL,
	)
}

// vmListReportFilteringBeforeAfterResults is a helper function used by the
// VMListReport function to generate a summary of filtering results.
//
// This summary is generating using provided functions to provide a collection
// of VirtualMachines before a filtering step against a collection of
// VirtualMachines after a filtering step has completed.
//
// If nil is provided in place of previousAfterFilterFunc the assumption will
// be made that no VirtualMachine collection is available for the previous
// filtering step. This covers cases where you are dealing with the first
// filtering step for VirtualMachines.
//
// If nil is provided in place of currentBeforeFilterFunc the assumption will
// be made that no VirtualMachine collection is available for before the
// current filtering step began. This covers cases where you are dealing with
// the first filtering step for VirtualMachines.
//
// The currentAfterFilterFunc argument is required.
func vmListReportFilteringBeforeAfterResults(
	w io.Writer,
	numAllVMs int,
	previousAfterFilterDesc string,
	previousAfterFilterFunc func() []mo.VirtualMachine,
	currentFilterDesc string,
	currentBeforeFilterFunc func() []mo.VirtualMachine,
	currentAfterFilterFunc func() []mo.VirtualMachine,
) {

	fmt.Fprintf(
		w,
		"%s(%d of %d) VMs before %s filtering was applied:%s%s",
		nagios.CheckOutputEOL,
		func() int {
			if currentBeforeFilterFunc == nil {
				return numAllVMs
			}
			return len(currentBeforeFilterFunc())
		}(),
		numAllVMs,
		currentFilterDesc,
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
	)

	switch {
	case previousAfterFilterFunc == nil:
		fmt.Fprintf(
			w,
			"* No filtering applied yet; skipping listing of all VMs.%s",
			nagios.CheckOutputEOL,
		)

	case len(currentBeforeFilterFunc()) == len(previousAfterFilterFunc()):
		fmt.Fprintf(
			w,
			"* Same list as after %s filtering.%s",
			previousAfterFilterDesc,
			nagios.CheckOutputEOL,
		)

	default:
		for _, vm := range currentBeforeFilterFunc() {
			fmt.Fprintf(
				w,
				"* %s%s",
				vm.Name,
				nagios.CheckOutputEOL,
			)
		}
	}

	fmt.Fprint(w, nagios.CheckOutputEOL)

	fmt.Fprintf(
		w,
		"%s(%d of %d) VMs after %s filtering was applied:%s%s",
		nagios.CheckOutputEOL,
		len(currentAfterFilterFunc()),
		numAllVMs,
		currentFilterDesc,
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
	)

	switch {
	case currentAfterFilterFunc != nil && currentBeforeFilterFunc == nil:
		for _, vm := range currentAfterFilterFunc() {
			fmt.Fprintf(
				w,
				"* %s%s",
				vm.Name,
				nagios.CheckOutputEOL,
			)
		}

	case len(currentAfterFilterFunc()) == len(currentBeforeFilterFunc()):
		fmt.Fprintf(
			w,
			"* Same list as before %s filtering.%s",
			currentFilterDesc,
			nagios.CheckOutputEOL,
		)

	default:
		for _, vm := range currentAfterFilterFunc() {
			fmt.Fprintf(
				w,
				"* %s%s",
				vm.Name,
				nagios.CheckOutputEOL,
			)
		}
	}
}

// VMFilterResultsPerfData provides performance data metrics for the results
// of performing filtering operations on a given VirtualMachines collection.
func VMFilterResultsPerfData(vmsFilterResults VMsFilterResults) []nagios.PerformanceData {
	return []nagios.PerformanceData{
		// The `time` (runtime) metric is appended at plugin exit, so do not
		// duplicate it here.
		{
			// This metric represents all non-template VirtualMachines in the
			// inventory.
			Label: "vms",
			Value: fmt.Sprintf("%d", vmsFilterResults.NumVMsAll()),
		},
		{
			// Alias to vms metric.
			Label: "vms_all",
			Value: fmt.Sprintf("%d", vmsFilterResults.NumVMsAll()),
		},
		{
			// Alias to vms_after_filtering performance metric.
			Label: "vms_evaluated",
			Value: fmt.Sprintf("%d", vmsFilterResults.NumVMsAfterFiltering()),
		},
		{
			// Alias to vms_evaluated performance metric.
			Label: "vms_after_filtering",
			Value: fmt.Sprintf("%d", vmsFilterResults.NumVMsAfterFiltering()),
		},
		{
			// We pull this metric from the collection remaining after
			// Resource Pool filtering so that we are looking at VMs which
			// have yet to be (potentially) filtered out based on power state.
			Label: "vms_powered_off",
			Value: fmt.Sprintf(
				"%d",
				CountVMsPowerStateOff(vmsFilterResults.VMsAfterResourcePoolFiltering()),
			),
		},
		{
			// We pull this metric from the collection remaining after
			// Resource Pool filtering to be consistent with how the
			// vms_powered_off metric is calculated.
			Label: "vms_powered_on",
			Value: fmt.Sprintf(
				"%d",
				CountVMsPowerStateOn(vmsFilterResults.VMsAfterResourcePoolFiltering()),
			),
		},
		{
			Label: "vms_excluded_by_name",
			Value: fmt.Sprintf("%d", vmsFilterResults.NumVMsExcludedByName()),
		},
		{
			Label: "vms_excluded_by_folder",
			Value: fmt.Sprintf("%d", vmsFilterResults.NumVMsExcludedByFolder()),
		},
		{
			Label: "vms_excluded_by_resource_pool",
			Value: fmt.Sprintf("%d", vmsFilterResults.NumVMsExcludedByResourcePool()),
		},
		{
			Label: "vms_excluded_by_power_state",
			Value: fmt.Sprintf("%d", vmsFilterResults.NumVMsExcludedByPowerState()),
		},
		{
			Label: "folders_all",
			Value: fmt.Sprintf("%d", vmsFilterResults.NumFoldersAll()),
		},
		{
			Label: "folders_excluded",
			Value: fmt.Sprintf("%d", vmsFilterResults.NumFoldersExcluded()),
		},
		{
			Label: "folders_included",
			Value: fmt.Sprintf("%d", vmsFilterResults.NumFoldersIncluded()),
		},
		{
			Label: "folders_evaluated",
			Value: fmt.Sprintf("%d", vmsFilterResults.NumFoldersAfterFiltering()),
		},
		{
			Label: "resource_pools_all",
			Value: fmt.Sprintf("%d", vmsFilterResults.NumRPsAll()),
		},
		{
			Label: "resource_pools_excluded",
			Value: fmt.Sprintf("%d", vmsFilterResults.NumRPsExcluded()),
		},
		{
			Label: "resource_pools_included",
			Value: fmt.Sprintf("%d", vmsFilterResults.NumRPsIncluded()),
		},
		{
			Label: "resource_pools_evaluated",
			Value: fmt.Sprintf("%d", vmsFilterResults.NumRPsAfterFiltering()),
		},
	}
}
