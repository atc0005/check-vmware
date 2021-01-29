// Copyright 2021 Adam Chalkley
//
// https://github.com/atc0005/check-vmware
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package vsphere

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/atc0005/go-nagios"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

// GetVMToolsStatusSummary accepts a collection of VirtualMachines and checks
// the ToolsStatus for each one providing an overall Nagios state label and
// exit code for the collection.
func GetVMToolsStatusSummary(vms []mo.VirtualMachine) (string, int) {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute GetVMToolsStatusSummary func.\n",
			time.Since(funcTimeStart),
		)
	}()

	var nagiosExitStateLabel string
	var nagiosExitStateCode int

Loop:
	for _, vm := range vms {

		// check specific tools issue to determine final Nagios state
		switch vm.Guest.ToolsStatus {

		case types.VirtualMachineToolsStatusToolsOk:
			continue

		case types.VirtualMachineToolsStatusToolsOld:

			// Not severe enough to immediately break as other more severe
			// issues may be present. Set state and allow the state to "carry"
			// at function exit.
			nagiosExitStateLabel = nagios.StateWARNINGLabel
			nagiosExitStateCode = nagios.StateWARNINGExitCode

		case types.VirtualMachineToolsStatusToolsNotRunning:
			nagiosExitStateLabel = nagios.StateCRITICALLabel
			nagiosExitStateCode = nagios.StateCRITICALExitCode
			break Loop

		case types.VirtualMachineToolsStatusToolsNotInstalled:
			nagiosExitStateLabel = nagios.StateCRITICALLabel
			nagiosExitStateCode = nagios.StateCRITICALExitCode
			break Loop

		// This should not be reached
		default:
			nagiosExitStateLabel = nagios.StateUNKNOWNLabel
			nagiosExitStateCode = nagios.StateUNKNOWNExitCode
			break Loop
		}

	}

	return nagiosExitStateLabel, nagiosExitStateCode

}

// FilterVMsWithToolsIssues filters the provided collection of VirtualMachines
// to just those with non-OK status.
func FilterVMsWithToolsIssues(vms []mo.VirtualMachine) []mo.VirtualMachine {

	// setup early so we can reference it from deferred stats output
	var vmsWithIssues []mo.VirtualMachine

	funcTimeStart := time.Now()

	defer func(vms []mo.VirtualMachine, filteredVMs *[]mo.VirtualMachine) {
		logger.Printf(
			"It took %v to execute FilterVMsWithToolsIssues func (for %d VMs, yielding %d VMs).\n",
			time.Since(funcTimeStart),
			len(vms),
			len(*filteredVMs),
		)
	}(vms, &vmsWithIssues)

	for _, vm := range vms {
		if vm.Guest.ToolsStatus != types.VirtualMachineToolsStatusToolsOk {
			vmsWithIssues = append(vmsWithIssues, vm)
		}
	}

	return vmsWithIssues

}

// VMToolsOneLineCheckSummary is used to generate a one-line Nagios service
// check results summary. This is the line most prominent in notifications.
func VMToolsOneLineCheckSummary(stateLabel string, vmsWithIssues []mo.VirtualMachine, evaluatedVMs []mo.VirtualMachine, rps []mo.ResourcePool) string {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute VMToolsOneLineCheckSummary func.\n",
			time.Since(funcTimeStart),
		)
	}()

	switch {
	case len(vmsWithIssues) > 0:
		return fmt.Sprintf(
			"%s: %d VMs with VMware Tools issues detected (evaluated %d VMs, %d Resource Pools)",
			stateLabel,
			len(vmsWithIssues),
			len(evaluatedVMs),
			len(rps),
		)

	default:

		return fmt.Sprintf(
			"%s: No VMware Tools issues detected (evaluated %d VMs, %d Resource Pools)",
			stateLabel,
			len(evaluatedVMs),
			len(rps),
		)

	}
}

// VMToolsReport generates a comprehensive summary including any active issues
// along with various verbose details intended to aid in troubleshooting check
// results at a glance. This information is provided for use with the Long
// Service Output field commonly displayed on the detailed service check
// results display in the web UI or in the body of many notifications.
func VMToolsReport(
	c *vim25.Client,
	allVMs []mo.VirtualMachine,
	evaluatedVMs []mo.VirtualMachine,
	vmsWithIssues []mo.VirtualMachine,
	vmsToExclude []string,
	evalPoweredOffVMs bool,
	includeRPs []string,
	excludeRPs []string,
	rps []mo.ResourcePool,
) string {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute VMToolsReport func.\n",
			time.Since(funcTimeStart),
		)
	}()

	rpNames := make([]string, len(rps))
	for i := range rps {
		rpNames[i] = rps[i].Name
	}

	var vmsReport strings.Builder

	switch {

	case len(vmsWithIssues) > 0:

		sort.Slice(vmsWithIssues, func(i, j int) bool {
			return strings.ToLower(vmsWithIssues[i].Name) < strings.ToLower(vmsWithIssues[j].Name)
		})

		for idx, vm := range vmsWithIssues {
			fmt.Fprintf(
				&vmsReport,
				"* %02d) %s (%s)%s",
				idx+1,
				vm.Name,
				string(vm.Guest.ToolsStatus),
				nagios.CheckOutputEOL,
			)
		}

	default:
		fmt.Fprintf(
			&vmsReport,
			"* No VMware Tools issues detected.%s",
			nagios.CheckOutputEOL,
		)
	}

	fmt.Fprintf(
		&vmsReport,
		"%s---%s%s",
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		&vmsReport,
		"* vSphere environment: %s%s",
		c.URL().String(),
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		&vmsReport,
		"* VMs (evaluated: %d, total: %d)%s",
		len(evaluatedVMs),
		len(allVMs),
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		&vmsReport,
		"* Powered off VMs evaluated: %t%s",
		evalPoweredOffVMs,
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		&vmsReport,
		"* Specified VMs to exclude (%d): [%v]%s",
		len(vmsToExclude),
		strings.Join(vmsToExclude, ", "),
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		&vmsReport,
		"* Specified Resource Pools to explicitly include (%d): [%v]%s",
		len(includeRPs),
		strings.Join(includeRPs, ", "),
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		&vmsReport,
		"* Specified Resource Pools to explicitly exclude (%d): [%v]%s",
		len(excludeRPs),
		strings.Join(excludeRPs, ", "),
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		&vmsReport,
		"* Resource Pools evaluated (%d): [%v]%s",
		len(rpNames),
		strings.Join(rpNames, ", "),
		nagios.CheckOutputEOL,
	)

	return vmsReport.String()
}
