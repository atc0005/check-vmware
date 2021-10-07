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
// the VMware Tools status for each one, providing an overall Nagios state
// label and exit code for the collection.
//
// NOTE: This function does *NOT* differentiate between VirtualMachines that
// are powered on and those that are powered off. If only powered on
// VirtualMachines should be evaluated, the caller should perform this
// filtering first before passing the collection to this function.
func GetVMToolsStatusSummary(vms []mo.VirtualMachine) nagios.ServiceState {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute GetVMToolsStatusSummary func.\n",
			time.Since(funcTimeStart),
		)
	}()

	var serviceState nagios.ServiceState

	for _, vm := range vms {

		serviceState = getVMwareToolsServiceState(vm)

		switch serviceState.ExitCode {

		// CRITICAL is as bad as it gets, so if we encounter this state go
		// ahead and consider the entire collection in this "overall" state.
		case nagios.StateCRITICALExitCode:
			return serviceState

		// UNKNOWN likely indicates an issue matching up the VMware Tools
		// status with known API values, so consider the entire collection in
		// this "overall" state.
		case nagios.StateUNKNOWNExitCode:
			return serviceState

		// For WARNING or OK states, we continue on to evaluating the next VM
		// retaining the service state we just received. If we don't find a VM
		// with a more severe service state we will return the last result as
		// the "overall" state.
		default:
			continue
		}

	}

	return serviceState
}

// FilterVMsWithToolsIssues filters the provided collection of VirtualMachines
// to just those with non-OK status, unless powered off VMs are also
// evaluated. In that case, ignore any powered off VirtualMachines with VMware
// Tools in a "not running" state which are otherwise current or unmanaged.
func FilterVMsWithToolsIssues(vms []mo.VirtualMachine, includePoweredOff bool) []mo.VirtualMachine {

	// setup early so we can reference it from deferred stats output
	vmsWithIssues := make([]mo.VirtualMachine, 0, len(vms))

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

		// If sysadmin did not opt to evaluate powered off VMs and VM is
		// powered off, skip evaluating VMware Tools state.
		if !includePoweredOff &&
			vm.Runtime.PowerState == types.VirtualMachinePowerStatePoweredOff {
			continue
		}

		if getVMwareToolsServiceState(vm).ExitCode != nagios.StateOKExitCode {
			vmsWithIssues = append(vmsWithIssues, vm)
		}

	}

	return vmsWithIssues

}

// VMToolsOneLineCheckSummary is used to generate a one-line Nagios service
// check results summary. This is the line most prominent in notifications.
func VMToolsOneLineCheckSummary(
	stateLabel string,
	evaluatedVMs []mo.VirtualMachine,
	vmsWithIssues []mo.VirtualMachine,
	rps []mo.ResourcePool,
) string {

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
				"* %02d) %s (%s, %s)%s",
				idx+1,
				vm.Name,
				string(vm.Runtime.PowerState),
				vm.Guest.ToolsVersionStatus2,
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
		"* Plugin User Agent: %s%s",
		c.Client.UserAgent,
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

// getVMwareToolsServiceState evaluates a VirtualMachine to determine an
// overall service state for the VM's VMware Tools status.
//
// References:
//
// https://developer.vmware.com/docs/vsphere-automation/latest/vcenter/data-structures/Vm/Tools/VersionStatus/
// https://vdc-repo.vmware.com/vmwb-repository/dcr-public/7989f521-fd57-4fff-9653-e6a5d5265089/1fd5908d-b8ce-49ca-887a-fefb3656e828/doc/vim.vm.GuestInfo.ToolsVersionStatus.html
// https://vdc-download.vmware.com/vmwb-repository/dcr-public/b50dcbbf-051d-4204-a3e7-e1b618c1e384/538cf2ec-b34f-4bae-a332-3820ef9e7773/vim.vm.GuestInfo.ToolsVersionStatus.html
func getVMwareToolsServiceState(vm mo.VirtualMachine) nagios.ServiceState {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute getVMwareToolsServiceState func.\n",
			time.Since(funcTimeStart),
		)
	}()

	switch {

	// If VM is powered off, don't evaluate whether the VMware Tools
	// status indicates "not running", focus only on the
	// VirtualMachineToolsVersionStatus values in the toolsVersionStatus2
	// API field.

	// VM is powered on, but Tools are not running.
	case vm.Runtime.PowerState == types.VirtualMachinePowerStatePoweredOn &&
		types.VirtualMachineToolsStatus(vm.Guest.ToolsRunningStatus) ==
			types.VirtualMachineToolsStatusToolsNotRunning:

		return nagios.ServiceState{
			Label:    nagios.StateCRITICALLabel,
			ExitCode: nagios.StateCRITICALExitCode,
		}

	// VMware Tools is not installed.
	case types.VirtualMachineToolsVersionStatus(vm.Guest.ToolsVersionStatus2) ==
		types.VirtualMachineToolsVersionStatusGuestToolsNotInstalled:

		return nagios.ServiceState{
			Label:    nagios.StateCRITICALLabel,
			ExitCode: nagios.StateCRITICALExitCode,
		}

	// VMware Tools is installed, and the version is current.
	case types.VirtualMachineToolsVersionStatus(vm.Guest.ToolsVersionStatus2) ==
		types.VirtualMachineToolsVersionStatusGuestToolsCurrent:

		return nagios.ServiceState{
			Label:    nagios.StateOKLabel,
			ExitCode: nagios.StateOKExitCode,
		}

	// VMware Tools is installed, but it is not managed by VMware. This
	// includes open-vm-tools or OSPs which should be managed by the guest
	// operating system.
	case types.VirtualMachineToolsVersionStatus(vm.Guest.ToolsVersionStatus2) ==
		types.VirtualMachineToolsVersionStatusGuestToolsUnmanaged:

		return nagios.ServiceState{
			Label:    nagios.StateOKLabel,
			ExitCode: nagios.StateOKExitCode,
		}

	// VMware Tools is installed, but the version is too old.
	case types.VirtualMachineToolsVersionStatus(vm.Guest.ToolsVersionStatus2) ==
		types.VirtualMachineToolsVersionStatusGuestToolsTooOld:

		return nagios.ServiceState{
			Label:    nagios.StateCRITICALLabel,
			ExitCode: nagios.StateCRITICALExitCode,
		}

	// VMware Tools is installed, supported, but a newer version is
	// available.
	case types.VirtualMachineToolsVersionStatus(vm.Guest.ToolsVersionStatus2) ==
		types.VirtualMachineToolsVersionStatusGuestToolsSupportedOld:

		return nagios.ServiceState{
			Label:    nagios.StateWARNINGLabel,
			ExitCode: nagios.StateWARNINGExitCode,
		}

	// VMware Tools is installed, but the version is not current.
	//
	// NOTE: This is a separate enum value in the API; it is not clear how
	// it differs from the guestToolsSupportedOld value, so we check for
	// it separately from that value.
	case types.VirtualMachineToolsVersionStatus(vm.Guest.ToolsVersionStatus2) ==
		types.VirtualMachineToolsVersionStatusGuestToolsNeedUpgrade:

		return nagios.ServiceState{
			Label:    nagios.StateWARNINGLabel,
			ExitCode: nagios.StateWARNINGExitCode,
		}

	// VMware Tools is installed, supported, and newer than the version
	// available on the host.
	case types.VirtualMachineToolsVersionStatus(vm.Guest.ToolsVersionStatus2) ==
		types.VirtualMachineToolsVersionStatusGuestToolsSupportedNew:

		return nagios.ServiceState{
			Label:    nagios.StateOKLabel,
			ExitCode: nagios.StateOKExitCode,
		}

	// VMware Tools is installed, and the version is known to be too new
	// to work correctly with this virtual machine.
	case types.VirtualMachineToolsVersionStatus(vm.Guest.ToolsVersionStatus2) ==
		types.VirtualMachineToolsVersionStatusGuestToolsTooNew:

		return nagios.ServiceState{
			Label:    nagios.StateCRITICALLabel,
			ExitCode: nagios.StateCRITICALExitCode,
		}

	// VMware Tools is installed, but the installed version is known to
	// have a grave bug and should be immediately upgraded.
	case types.VirtualMachineToolsVersionStatus(vm.Guest.ToolsVersionStatus2) ==
		types.VirtualMachineToolsVersionStatusGuestToolsBlacklisted:

		return nagios.ServiceState{
			Label:    nagios.StateCRITICALLabel,
			ExitCode: nagios.StateCRITICALExitCode,
		}

	// We should only reach this point if the vSphere API has been extended
	// and this library hasn't been updated to account for those changes.
	default:

		return nagios.ServiceState{
			Label:    nagios.StateUNKNOWNLabel,
			ExitCode: nagios.StateUNKNOWNExitCode,
		}

	}

}
