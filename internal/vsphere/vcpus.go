// Copyright 2021 Adam Chalkley
//
// https://github.com/atc0005/check-vmware
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package vsphere

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/atc0005/go-nagios"
	"github.com/vmware/govmomi/vim25"
)

// ErrVCPUsUsageThresholdCrossed indicates that specified
// vCPUs allocation has exceeded a given threshold
var ErrVCPUsUsageThresholdCrossed = errors.New("vCPUS allocation exceeds specified threshold")

// VirtualCPUsOneLineCheckSummary is used to generate a one-line Nagios
// service check results summary. This is the line most prominent in
// notifications.
func VirtualCPUsOneLineCheckSummary(
	stateLabel string,
	vmsFilterResults VMsFilterResults,
	vCPUsAllocated int32,
	vCPUsMax int,
) string {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
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
			vmsFilterResults.NumVMsAfterFiltering(),
			vmsFilterResults.NumRPsAfterFiltering(),
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
			vmsFilterResults.NumVMsAfterFiltering(),
			vmsFilterResults.NumRPsAfterFiltering(),
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
	vmsFilterOptions VMsFilterOptions,
	vmsFilterResults VMsFilterResults,
	vCPUsAllocated int32,
	vCPUsMax int,
) string {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute VirtualCPUsReport func.\n",
			time.Since(funcTimeStart),
		)
	}()

	var vmsReport strings.Builder

	// This is shown regardless of whether the plugin is considered to be in a
	// non-OK state.
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
		"%sTop 10 vCPU consumers:%s%s",
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
	)

	evaluatedVMs := vmsFilterResults.VMsAfterFiltering()
	sort.Slice(evaluatedVMs, func(i, j int) bool {
		return evaluatedVMs[i].Summary.Config.NumCpu > evaluatedVMs[j].Summary.Config.NumCpu
	})

	// grab up to the first 10 VMs, presorted by large vCPU consumption
	sampleSize := len(evaluatedVMs)
	if sampleSize > 10 {
		sampleSize = 10
	}
	topTen := evaluatedVMs[:sampleSize]

	switch {
	case len(topTen) == 0:
		fmt.Fprintf(&vmsReport, "* None %s", nagios.CheckOutputEOL)
	default:
		for _, vm := range topTen {
			fmt.Fprintf(
				&vmsReport,
				"* %s (%d vCPUs)%s",
				vm.Name,
				vm.Summary.Config.NumCpu,
				nagios.CheckOutputEOL,
			)
		}
	}

	fmt.Fprintf(
		&vmsReport,
		"%sTen most recently started VMs:%s%s",
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
	)

	// Regardless of earlier decision whether to exclude powered off VMs from
	// vCPU consumption calculations, we explicitly exclude here in order to
	// limit evaluation of "most recently booted" to powered on VMs only.
	poweredOnVMs, _ := FilterVMsByPowerState(evaluatedVMs, false)

	// sort before we sample the VMs so that we only get the ones with lowest
	// power cycle uptime; require that the VM be powered on in order to sort
	// in the intended order.
	sort.Slice(poweredOnVMs, func(i, j int) bool {
		return poweredOnVMs[i].Summary.QuickStats.UptimeSeconds < poweredOnVMs[j].Summary.QuickStats.UptimeSeconds

	})

	// Grab a sampling of the powered on VMs which were most recently booted.
	sampleSize = len(poweredOnVMs)
	if sampleSize > 10 {
		sampleSize = 10
	}
	bottomTen := poweredOnVMs[:sampleSize]

	switch {
	case len(bottomTen) == 0:
		fmt.Fprintf(&vmsReport, "* None %s", nagios.CheckOutputEOL)
	default:
		for _, vm := range bottomTen {
			uptime := time.Duration(vm.Summary.QuickStats.UptimeSeconds) * time.Second
			uptimeDays := uptime.Hours() / 24

			fmt.Fprintf(
				&vmsReport,
				"* %s: (%.2f days)%s",
				vm.Name,
				uptimeDays,
				nagios.CheckOutputEOL,
			)
		}
	}

	vmFilterResultsReportTrailer(
		&vmsReport,
		c,
		vmsFilterOptions,
		vmsFilterResults,
		true,
	)

	return vmsReport.String()
}
