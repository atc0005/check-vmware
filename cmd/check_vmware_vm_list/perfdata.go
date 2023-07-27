// Copyright 2021 Adam Chalkley
//
// https://github.com/atc0005/check-vmware
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package main

import (
	"fmt"

	"github.com/atc0005/go-nagios"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

// FilteredItems is a tally of relevant details used during compilation of
// performance metrics.
type FilteredItems struct {
	NumVMsRemainingAfterFiltering int
	NumVMsExcludedByName          int
	NumVMsExcludedByPowerState    int

	// NumResourcePoolsEvaluated is the number of resource pools remaining
	// after they have been explicitly included or excluded.
	NumResourcePoolsEvaluated int
	NumResourcePoolsIncluded  int
	NumResourcePoolsExcluded  int
}

// getPerfData gathers performance data metrics that we wish to report.
func getPerfData(
	allVMs []mo.VirtualMachine,
	filteredItems FilteredItems,
) []nagios.PerformanceData {
	return []nagios.PerformanceData{
		// The `time` (runtime) metric is appended at plugin exit, so do not
		// duplicate it here.
		{
			Label: "vms",
			Value: fmt.Sprintf("%d", len(allVMs)),
		},
		{
			Label: "vms_after_filtering",
			Value: fmt.Sprintf("%d", filteredItems.NumVMsRemainingAfterFiltering),
		},
		{
			Label: "vms_powered_on",
			Value: fmt.Sprintf(
				"%d",
				func() int {
					var count int
					for _, vm := range allVMs {
						if vm.Runtime.PowerState == types.VirtualMachinePowerStatePoweredOn {
							count++
						}
					}
					return count
				}(),
			),
		},
		{
			Label: "vms_powered_off",
			Value: fmt.Sprintf(
				"%d",
				func() int {
					var count int
					for _, vm := range allVMs {
						if vm.Runtime.PowerState == types.VirtualMachinePowerStatePoweredOff ||
							vm.Runtime.PowerState == types.VirtualMachinePowerStateSuspended {
							count++
						}
					}
					return count
				}(),
			),
		},
		{
			Label: "vms_excluded_by_name",
			Value: fmt.Sprintf("%d", filteredItems.NumVMsExcludedByName),
		},
		{
			Label: "vms_excluded_by_power_state",
			Value: fmt.Sprintf("%d", filteredItems.NumVMsExcludedByPowerState),
		},
		{
			Label: "resource_pools_excluded",
			Value: fmt.Sprintf("%d", filteredItems.NumResourcePoolsExcluded),
		},
		{
			Label: "resource_pools_included",
			Value: fmt.Sprintf("%d", filteredItems.NumResourcePoolsIncluded),
		},
		{
			Label: "resource_pools_evaluated",
			Value: fmt.Sprintf("%d", filteredItems.NumResourcePoolsEvaluated),
		},
	}

}
