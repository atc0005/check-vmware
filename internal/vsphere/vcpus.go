// Copyright 2021 Adam Chalkley
//
// https://github.com/atc0005/check-vmware
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package vsphere

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/atc0005/go-nagios"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
)

// VirtualCPUsOneLineCheckSummary is used to generate a one-line Nagios
// service check results summary. This is the line most prominent in
// notifications.
func VirtualCPUsOneLineCheckSummary(
	stateLabel string,
	vCPUsAllocated int32,
	vCPUsMax int,
	evaluatedVMs []mo.VirtualMachine,
	rps []mo.ResourcePool,
) string {

	funcTimeStart := time.Now()

	defer func() {
		fmt.Fprintf(
			os.Stderr,
			"It took %v to execute VirtualCPUsOneLineCheckSummary func.\n",
			time.Since(funcTimeStart),
		)
	}()

	vCPUsPercentageUsed := float32(vCPUsAllocated) / float32(vCPUsMax) * 100

	switch {

	case vCPUsAllocated > int32(vCPUsMax):
		vCPUsOverage := vCPUsAllocated - int32(vCPUsMax)
		return fmt.Sprintf(
			"%s: %d vCPUs allocated (%.1f%%); %d more allocated than %d allowed"+
				" (evaluated %d VMs, %d Resource Pools)",
			stateLabel,
			vCPUsAllocated,
			vCPUsPercentageUsed,
			vCPUsOverage,
			vCPUsMax,
			len(evaluatedVMs),
			len(rps),
		)

	default:
		vCPUsRemaining := int32(vCPUsMax) - vCPUsAllocated
		return fmt.Sprintf(
			"%s: %d vCPUs allocated (%.1f%%); %d more remaining from %d allowed"+
				" (evaluated %d VMs, %d Resource Pools)",
			stateLabel,
			vCPUsAllocated,
			vCPUsPercentageUsed,
			vCPUsRemaining,
			vCPUsMax,
			len(evaluatedVMs),
			len(rps),
		)

	}
}

// VirtualCPUsReport generates a summary of vCPU usage along with various
// verbose details intended to aid in troubleshooting check results at a
// glance. This information is provided for use with the Long Service Output
// field commonly displayed on the detailed service check results display in
// the web UI or in the body of many notifications.
func VirtualCPUsReport(
	c *vim25.Client,
	vCPUsAllocated int32,
	vCPUsMax int,
	allVMs []mo.VirtualMachine,
	evaluatedVMs []mo.VirtualMachine,
	vmsToExclude []string,
	includeRPs []string,
	excludeRPs []string,
	rps []mo.ResourcePool,
) string {

	funcTimeStart := time.Now()

	defer func() {
		fmt.Fprintf(
			os.Stderr,
			"It took %v to execute VirtualCPUsReport func.\n",
			time.Since(funcTimeStart),
		)
	}()

	rpNames := make([]string, len(rps))
	for i := range rps {
		rpNames[i] = rps[i].Name
	}

	var vmsReport strings.Builder

	fmt.Fprintf(
		&vmsReport,
		"* vCPUs%s** Allocated: %d (%.1f%%)%s** Max Allowed: %d%s",
		nagios.CheckOutputEOL,
		vCPUsAllocated,
		float32(vCPUsAllocated)/float32(vCPUsMax)*100,
		nagios.CheckOutputEOL,
		vCPUsMax,
		nagios.CheckOutputEOL,
	)

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
