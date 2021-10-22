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
	"sort"
	"strings"
	"time"

	"github.com/atc0005/go-nagios"
	"github.com/vmware/govmomi/units"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

// ErrHostSystemMemoryUsageThresholdCrossed indicates that specified host
// memory usage has exceeded a given threshold
var ErrHostSystemMemoryUsageThresholdCrossed = errors.New("host memory usage exceeds specified threshold")

// ErrHostSystemCPUUsageThresholdCrossed indicates that specified host CPU
// usage has exceeded a given threshold
var ErrHostSystemCPUUsageThresholdCrossed = errors.New("host CPU usage exceeds specified threshold")

// ErrHostSystemHardwarePropertiesUnavailable indicates that specified host
// hardware properties are unavailable. This is likely due to permission
// issues for the service account or a shallow host properties retrieval
// request (coding error).
var ErrHostSystemHardwarePropertiesUnavailable = errors.New("host hardware properties unavailable")

// HostSystemMemorySummary tracks memory usage details for a specific
// HostSystem.
type HostSystemMemorySummary struct {
	HostSystem             mo.HostSystem
	MemoryUsedPercent      float64
	MemoryRemainingPercent float64

	// MemoryUsed is the amount of memory used by the host in bytes.
	MemoryUsed int64

	// MemoryUsed is the amount of memory remaining to the host in bytes.
	MemoryRemaining int64

	// MemoryTotal is the total amount of memory for the host in bytes.
	MemoryTotal       int64
	CriticalThreshold int
	WarningThreshold  int
}

// HostSystemCPUSummary tracks CPU usage details for a specific HostSystem.
type HostSystemCPUSummary struct {
	HostSystem          mo.HostSystem
	CPUUsedPercent      float64
	CPURemainingPercent float64

	// CPUUsed is the amount of CPU used by the host in Hz.
	CPUUsed float64

	// CPURemaining is the amount of CPU capacity remaining to the host in Hz.
	CPURemaining float64

	// CPUTotal is the total amount of CPU capacity for the host in Hz.
	CPUTotal          float64
	CriticalThreshold int
	WarningThreshold  int
}

// NewHostSystemMemoryUsageSummary receives a HostSystem and generates summary
// information used to determine if usage levels have crossed user-specified
// thresholds.
func NewHostSystemMemoryUsageSummary(hs mo.HostSystem, criticalThreshold int, warningThreshold int) HostSystemMemorySummary {

	// total memory in bytes
	memoryTotal := hs.Hardware.MemorySize

	// memory used in bytes
	memoryUsed := int64(hs.Summary.QuickStats.OverallMemoryUsage) * units.MB

	// memory remaining in bytes
	memoryRemaining := memoryTotal - memoryUsed

	memoryRemainingPercentage := float64(memoryRemaining) / float64(memoryTotal) * 100
	memoryUsedPercentage := 100 - memoryRemainingPercentage

	hsUsage := HostSystemMemorySummary{
		HostSystem:             hs,
		MemoryUsedPercent:      memoryUsedPercentage,
		MemoryRemainingPercent: memoryRemainingPercentage,
		MemoryUsed:             memoryUsed,
		MemoryRemaining:        memoryRemaining,
		MemoryTotal:            memoryTotal,
		CriticalThreshold:      criticalThreshold,
		WarningThreshold:       warningThreshold,
	}

	return hsUsage

}

// NewHostSystemCPUUsageSummary receives a HostSystem and generates summary
// information used to determine if usage levels have crossed user-specified
// thresholds.
func NewHostSystemCPUUsageSummary(hs mo.HostSystem, criticalThreshold int, warningThreshold int) (HostSystemCPUSummary, error) {

	if hs.Summary.Hardware == nil {
		return HostSystemCPUSummary{}, fmt.Errorf(
			"error creating HostSystemCPUSummary: %w",
			ErrHostSystemHardwarePropertiesUnavailable,
		)
	}

	numCPUCores := hs.Summary.Hardware.NumCpuCores

	// base value in MHz, convert to Hz
	cpuUsage := float64(hs.Summary.QuickStats.OverallCpuUsage) * MHz
	cpuSpeedPerCore := float64(hs.Summary.Hardware.CpuMhz) * MHz

	// capacity in Hz
	cpuTotalCapacity := (float64(numCPUCores) * cpuSpeedPerCore)
	cpuRemainingCapacity := cpuTotalCapacity - cpuUsage

	cpuUsagePercent := cpuUsage / cpuTotalCapacity * 100
	cpuCapacityRemainingPercent := 100 - cpuUsagePercent

	hsUsage := HostSystemCPUSummary{
		HostSystem:          hs,
		CPUUsedPercent:      cpuUsagePercent,
		CPURemainingPercent: cpuCapacityRemainingPercent,
		CPUUsed:             cpuUsage,
		CPURemaining:        cpuRemainingCapacity,
		CPUTotal:            cpuTotalCapacity,
		CriticalThreshold:   criticalThreshold,
		WarningThreshold:    warningThreshold,
	}

	return hsUsage, nil

}

// IsWarningState indicates whether HostSystem memory usage has crossed the
// WARNING level threshold.
func (hss HostSystemMemorySummary) IsWarningState() bool {
	return hss.MemoryUsedPercent < float64(hss.CriticalThreshold) &&
		hss.MemoryUsedPercent >= float64(hss.WarningThreshold)
}

// IsCriticalState indicates whether HostSystem memory usage has crossed the
// CRITICAL level threshold.
func (hss HostSystemMemorySummary) IsCriticalState() bool {
	return hss.MemoryUsedPercent >= float64(hss.CriticalThreshold)
}

// IsWarningState indicates whether HostSystem CPU usage has crossed the
// WARNING level threshold.
func (hss HostSystemCPUSummary) IsWarningState() bool {
	return hss.CPUUsedPercent < float64(hss.CriticalThreshold) &&
		hss.CPUUsedPercent >= float64(hss.WarningThreshold)
}

// IsCriticalState indicates whether HostSystem CPU usage has crossed the
// CRITICAL level threshold.
func (hss HostSystemCPUSummary) IsCriticalState() bool {
	return hss.CPUUsedPercent >= float64(hss.CriticalThreshold)
}

// GetHostSystems accepts a context, a connected client and a boolean value
// indicating whether a subset of properties per HostSystem are retrieved. A
// collection of HostSystems with requested properties is returned. If
// requested, a subset of all available properties will be retrieved (faster)
// instead of recursively fetching all properties (about 2x as slow).
func GetHostSystems(ctx context.Context, c *vim25.Client, propsSubset bool) ([]mo.HostSystem, error) {

	funcTimeStart := time.Now()

	// declare this early so that we can grab a pointer to it in order to
	// access the entries later
	var hss []mo.HostSystem

	defer func(hss *[]mo.HostSystem) {
		logger.Printf(
			"It took %v to execute GetHostSystems func (and retrieve %d HostSystems).\n",
			time.Since(funcTimeStart),
			len(*hss),
		)
	}(&hss)

	err := getObjects(ctx, c, &hss, c.ServiceContent.RootFolder, propsSubset)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve HostSystems: %w", err)
	}

	sort.Slice(hss, func(i, j int) bool {
		return strings.ToLower(hss[i].Name) < strings.ToLower(hss[j].Name)
	})

	return hss, nil
}

// GetHostSystemByName accepts the name of a HostSystem, the name of a
// datacenter and a boolean value indicating whether only a subset of
// properties for the HostSystem should be returned. If requested, a subset of
// all available properties will be retrieved (faster) instead of recursively
// fetching all properties (about 2x as slow). If the datacenter name is an
// empty string then the default datacenter will be used.
func GetHostSystemByName(ctx context.Context, c *vim25.Client, hsName string, datacenter string, propsSubset bool) (mo.HostSystem, error) {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute GetHostSystemByName func.\n",
			time.Since(funcTimeStart),
		)
	}()

	var hostSystem mo.HostSystem
	err := getObjectByName(ctx, c, &hostSystem, hsName, datacenter, propsSubset)

	if err != nil {
		return mo.HostSystem{}, err
	}

	return hostSystem, nil

}

// FilterHostSystemsByName accepts a collection of HostSystems and a
// HostSystem name to filter against. An error is returned if the list of
// HostSystems is empty or if a match was not found. The matching HostSystem
// is returned along with the number of HostSystems that were excluded.
func FilterHostSystemsByName(hss []mo.HostSystem, hsName string) (mo.HostSystem, int, error) {

	funcTimeStart := time.Now()

	// If error condition, no exclusions are made
	numExcluded := 0

	defer func() {
		logger.Printf(
			"It took %v to execute FilterHostSystemsByName func.\n",
			time.Since(funcTimeStart),
		)
	}()

	if len(hss) == 0 {
		return mo.HostSystem{}, numExcluded, fmt.Errorf("received empty list of HostSystems to filter by name")
	}

	for _, hs := range hss {
		if hs.Name == hsName {
			// we are excluding everything but the single name value match
			numExcluded = len(hss) - 1
			return hs, numExcluded, nil
		}
	}

	return mo.HostSystem{}, numExcluded, fmt.Errorf(
		"error: failed to retrieve HostSystem using provided name %q",
		hsName,
	)

}

// FilterHostSystemsByID receives a collection of HostSystems and a HostSystem
// ID to filter against. An error is returned if the list of HostSystems is
// empty or if a match was not found. The matching HostSystem is returned
// along with the number of HostSystems that were excluded.
func FilterHostSystemsByID(hss []mo.HostSystem, hsID string) (mo.HostSystem, int, error) {

	funcTimeStart := time.Now()

	// If error condition, no exclusions are made
	numExcluded := 0

	defer func() {
		logger.Printf(
			"It took %v to execute FilterHostSystemsByID func.\n",
			time.Since(funcTimeStart),
		)
	}()

	if len(hss) == 0 {
		return mo.HostSystem{}, numExcluded, fmt.Errorf("received empty list of HostSystems to filter by ID")
	}

	for _, hs := range hss {
		// return match, if available

		if hs.Summary.Host.Value == hsID {
			// we are excluding everything but the single ID value match
			numExcluded = len(hss) - 1
			return hs, numExcluded, nil
		}
	}

	return mo.HostSystem{}, numExcluded, fmt.Errorf(
		"error: failed to retrieve HostSystem using provided id %q",
		hsID,
	)

}

// GetHostSystemsTotalMemory returns the total memory capacity for all
// HostSystems. Unless requested, offline or otherwise unavailable hosts are
// included for evaluation based on the assumption that offline hosts are
// offline for only a brief time and should still be considered part of
// overall cluster capacity.
func GetHostSystemsTotalMemory(ctx context.Context, c *vim25.Client, excludeOffline bool) (int64, error) {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute GetHostSystemTotalMemory func.\n",
			time.Since(funcTimeStart),
		)
	}()

	clusterHosts, err := GetHostSystems(ctx, c, true)
	if err != nil {
		return 0, fmt.Errorf(
			"failed to gather total memory capacity for host systems: %w",
			err,
		)
	}

	var clusterMemory int64
	for _, host := range clusterHosts {

		// Evaluate offline systems by default, unless requested otherwise.
		if excludeOffline {

			logger.Printf("Checking host %s availability ... \n", host.Name)

			switch {

			case host.Runtime.PowerState == types.HostSystemPowerStatePoweredOn &&
				host.Runtime.ConnectionState == types.HostSystemConnectionStateConnected:
				// desired state, no other limiting factors detected

			case host.Runtime.InMaintenanceMode:
				logger.Printf("Host %s is in maintenance mode, skipping evaluation ...\n", host.Name)
				continue

			case host.Runtime.InQuarantineMode != nil && *host.Runtime.InQuarantineMode:
				logger.Printf("Host %s is in quarantine mode, skipping evaluation ...\n", host.Name)
				continue

			case host.Runtime.PowerState == types.HostSystemPowerStatePoweredOff:
				logger.Printf("Host %s is powered off, skipping evaluation ...\n", host.Name)
				continue

			case host.Runtime.PowerState == types.HostSystemPowerStateStandBy:
				logger.Printf("Host %s is in standby, skipping evaluation ...\n", host.Name)
				continue

			case host.Runtime.ConnectionState == types.HostSystemConnectionStateDisconnected:
				logger.Printf("Host %s is disconnected, skipping evaluation ...\n", host.Name)
				continue

			case host.Runtime.ConnectionState == types.HostSystemConnectionStateNotResponding:
				logger.Printf("Host %s is not responding, skipping evaluation ...\n", host.Name)
				continue

			default:
				logger.Printf("Host %s is in an UNKNOWN state, skipping evaluation ...\n", host.Name)
				continue

			}
		}

		logger.Printf(
			"Host %s has %s memory capacity.\n",
			host.Name,
			units.ByteSize(host.Hardware.MemorySize),
		)

		clusterMemory += host.Hardware.MemorySize
	}

	return clusterMemory, nil

}

// HostSystemMemoryUsageOneLineCheckSummary is used to generate a one-line
// Nagios service check results summary. This is the line most prominent in
// notifications.
func HostSystemMemoryUsageOneLineCheckSummary(
	stateLabel string,
	hsVMs []mo.VirtualMachine,
	hsUsageSummary HostSystemMemorySummary,
) string {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute HostSystemMemoryUsageOneLineCheckSummary func.\n",
			time.Since(funcTimeStart),
		)
	}()

	// drop any powered off/suspended VMs from our list
	hsVMs, _ = FilterVMsByPowerState(hsVMs, false)

	var vmsMemUsedBytes int64 // int64 used to prevent int32 overflow
	for _, vm := range hsVMs {

		// vm.Summary.QuickStats.HostMemoryUsage == memory usage in MB
		vmsMemUsedBytes += int64(vm.Summary.QuickStats.HostMemoryUsage) * units.MB
	}

	vmsMemUsedPercentOfHost := (float64(vmsMemUsedBytes) / float64(hsUsageSummary.MemoryTotal)) * 100

	summaryTemplate := "%s: Host %s using %s (%.2f%%) of %s with %s (%.2f%%) remaining (%d visible VMs using %s (%.2f%%) memory)"
	// summaryTemplate := "%s: Host %s memory usage is %.2f%% (%s) of %s with %s remaining (%d visible VMs using %s (%.2f%%) memory)"
	// summaryTemplate := "%s: Host %s memory usage is %.2f%% of %s with %s remaining [WARNING: %d%% , CRITICAL: %d%%]"

	return fmt.Sprintf(
		summaryTemplate,
		stateLabel,
		hsUsageSummary.HostSystem.Name,
		units.ByteSize(hsUsageSummary.MemoryUsed),
		hsUsageSummary.MemoryUsedPercent,
		units.ByteSize(hsUsageSummary.MemoryTotal),
		units.ByteSize(hsUsageSummary.MemoryRemaining),
		hsUsageSummary.MemoryRemainingPercent,
		len(hsVMs),
		units.ByteSize(vmsMemUsedBytes),
		vmsMemUsedPercentOfHost,
	)

}

// HostSystemMemoryUsageReport generates a summary of HostSystem memory usage
// along with various verbose details intended to aid in troubleshooting check
// results at a glance. This information is provided for use with the Long
// Service Output field commonly displayed on the detailed service check
// results display in the web UI or in the body of many notifications.
func HostSystemMemoryUsageReport(
	c *vim25.Client,
	hsVMs []mo.VirtualMachine,
	hsUsageSummary HostSystemMemorySummary,
) string {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute HostSystemMemoryUsageReport func.\n",
			time.Since(funcTimeStart),
		)
	}()

	var report strings.Builder

	var vmsMemUsedBytes int64 // int64 used to prevent int32 overflow
	var vmsPoweredOn int
	var vmsPoweredOff int
	for _, vm := range hsVMs {

		// vm.Summary.QuickStats.HostMemoryUsage == memory usage in MB
		vmsMemUsedBytes += int64(vm.Summary.QuickStats.HostMemoryUsage) * units.MB

		switch {
		case vm.Runtime.PowerState == types.VirtualMachinePowerStatePoweredOn:
			vmsPoweredOn++
		default:
			vmsPoweredOff++
		}

	}

	vmsMemUsedPercentOfHost := (float64(vmsMemUsedBytes) / float64(hsUsageSummary.MemoryTotal)) * 100

	fmt.Fprintf(
		&report,
		"Host Summary:%s%s"+
			"* Name: %s%s"+
			"* Memory%s"+
			"** Used by all VMs: %s (%.2f%%)%s"+
			"** Used by visible VMs: %s (%.2f%%)%s"+
			"** Remaining: %s (%.2f%%)%s"+
			"* VMs%s"+
			"** Visible: %d%s"+
			"** Running: %d%s"+
			"** Off: %d%s",
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
		hsUsageSummary.HostSystem.Name,
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
		units.ByteSize(hsUsageSummary.MemoryUsed),
		hsUsageSummary.MemoryUsedPercent,
		nagios.CheckOutputEOL,
		units.ByteSize(vmsMemUsedBytes),
		vmsMemUsedPercentOfHost,
		nagios.CheckOutputEOL,
		units.ByteSize(hsUsageSummary.MemoryRemaining),
		hsUsageSummary.MemoryRemainingPercent,
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
		len(hsVMs),
		nagios.CheckOutputEOL,
		vmsPoweredOn,
		nagios.CheckOutputEOL,
		vmsPoweredOff,
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		&report,
		"%sVMs on host consuming memory (descending order):%s%s",
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
	)

	sort.Slice(hsVMs, func(i, j int) bool {
		return hsVMs[i].Summary.QuickStats.HostMemoryUsage > hsVMs[j].Summary.QuickStats.HostMemoryUsage
	})

	for _, vm := range hsVMs {
		if vm.Runtime.PowerState == types.VirtualMachinePowerStatePoweredOn {
			hostMemUsedBytes := int64(vm.Summary.QuickStats.HostMemoryUsage) * units.MB
			vmPercentOfHostMemUsed := float64(hostMemUsedBytes) / float64(hsUsageSummary.MemoryTotal) * 100
			fmt.Fprintf(
				&report,
				"* %s (Memory: %v, Host Memory Usage: %2.2f%%)%s",
				vm.Name,
				units.ByteSize(hostMemUsedBytes),
				vmPercentOfHostMemUsed,
				nagios.CheckOutputEOL,
			)
		}
	}

	if vmsPoweredOn == 0 {
		fmt.Fprintf(
			&report,
			"* None (visible)%s",
			nagios.CheckOutputEOL,
		)
	}

	fmt.Fprintf(
		&report,
		"%sVMs on host not consuming memory:%s%s",
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
	)

	sort.Slice(hsVMs, func(i, j int) bool {
		return strings.ToLower(hsVMs[i].Name) < strings.ToLower(hsVMs[j].Name)
	})

	for _, vm := range hsVMs {
		if vm.Runtime.PowerState != types.VirtualMachinePowerStatePoweredOn {
			fmt.Fprintf(
				&report,
				"* %s%s",
				vm.Name,
				nagios.CheckOutputEOL,
			)
		}
	}

	if vmsPoweredOff == 0 {
		fmt.Fprintf(
			&report,
			"* None (visible)%s",
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

	return report.String()
}

// HostSystemCPUUsageOneLineCheckSummary is used to generate a one-line
// Nagios service check results summary. This is the line most prominent in
// notifications.
func HostSystemCPUUsageOneLineCheckSummary(
	stateLabel string,
	hsVMs []mo.VirtualMachine,
	hsUsageSummary HostSystemCPUSummary,
) string {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute HostSystemCPUUsageOneLineCheckSummary func.\n",
			time.Since(funcTimeStart),
		)
	}()

	// drop any powered off/suspended VMs from our list
	hsVMs, _ = FilterVMsByPowerState(hsVMs, false)

	var vmsCPUUsage int64
	for _, vm := range hsVMs {
		// usage in MHz, convert to Hz
		vmsCPUUsage += int64(vm.Summary.QuickStats.OverallCpuUsage) * MHz
	}

	vmsMemUsedPercentOfHost := (float64(vmsCPUUsage) / hsUsageSummary.CPUTotal) * 100

	// summaryTemplate := "%s: Host %s CPU usage is %s (%.2f%%) of %s with %s (%.2f%%) remaining (%d visible VMs using %s (%.2f%%) memory)"
	summaryTemplate := "%s: Host %s using %s (%.2f%%) of %s with %s (%.2f%%) remaining CPU capacity (%d visible VMs using %s (%.2f%%) CPU)"

	return fmt.Sprintf(
		summaryTemplate,
		stateLabel,
		hsUsageSummary.HostSystem.Name,
		CPUSpeed(hsUsageSummary.CPUUsed),
		hsUsageSummary.CPUUsedPercent,
		CPUSpeed(hsUsageSummary.CPUTotal),
		CPUSpeed(hsUsageSummary.CPURemaining),
		hsUsageSummary.CPURemainingPercent,
		len(hsVMs),
		CPUSpeed(vmsCPUUsage),
		vmsMemUsedPercentOfHost,
	)

}

// HostSystemCPUUsageReport generates a summary of HostSystem CPU usage along
// with various verbose details intended to aid in troubleshooting check
// results at a glance. This information is provided for use with the Long
// Service Output field commonly displayed on the detailed service check
// results display in the web UI or in the body of many notifications.
func HostSystemCPUUsageReport(
	c *vim25.Client,
	hsVMs []mo.VirtualMachine,
	hsUsageSummary HostSystemCPUSummary,
) string {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute HostSystemCPUUsageReport func.\n",
			time.Since(funcTimeStart),
		)
	}()

	var report strings.Builder

	var vmsCPUUsage int64
	var vmsPoweredOn int
	var vmsPoweredOff int
	for _, vm := range hsVMs {

		// usage in MHz, convert to Hz
		vmCPUUsage := int64(vm.Summary.QuickStats.OverallCpuUsage) * MHz
		vmsCPUUsage += vmCPUUsage
		logger.Printf("VM %s used %s \n", vm.Name, CPUSpeed(vmCPUUsage))

		switch {
		case vm.Runtime.PowerState == types.VirtualMachinePowerStatePoweredOn:
			vmsPoweredOn++
		default:
			vmsPoweredOff++
		}

	}

	vmsCPUUsedPercentOfHost := (float64(vmsCPUUsage) / hsUsageSummary.CPUTotal) * 100

	fmt.Fprintf(
		&report,
		"Host Summary:%s%s"+
			"* Name: %s%s"+
			"* CPU%s"+
			"** Used by all VMs: %s (%.2f%%)%s"+
			"** Used by visible VMs: %s (%.2f%%)%s"+
			"** Remaining: %s (%.2f%%)%s"+
			"* VMs%s"+
			"** Visible: %d%s"+
			"** Running: %d%s"+
			"** Off: %d%s",
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
		hsUsageSummary.HostSystem.Name,
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
		CPUSpeed(hsUsageSummary.CPUUsed),
		hsUsageSummary.CPUUsedPercent,
		nagios.CheckOutputEOL,
		CPUSpeed(vmsCPUUsage),
		vmsCPUUsedPercentOfHost,
		nagios.CheckOutputEOL,
		CPUSpeed(hsUsageSummary.CPURemaining),
		hsUsageSummary.CPURemainingPercent,
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
		len(hsVMs),
		nagios.CheckOutputEOL,
		vmsPoweredOn,
		nagios.CheckOutputEOL,
		vmsPoweredOff,
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		&report,
		"%sVMs on host consuming CPU (descending order):%s%s",
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
	)

	sort.Slice(hsVMs, func(i, j int) bool {
		return hsVMs[i].Summary.QuickStats.OverallCpuUsage > hsVMs[j].Summary.QuickStats.OverallCpuUsage
	})

	for _, vm := range hsVMs {
		if vm.Runtime.PowerState == types.VirtualMachinePowerStatePoweredOn {
			hostCPUUsed := int64(vm.Summary.QuickStats.OverallCpuUsage) * MHz
			vmPercentOfHostCPUUsed := (float64(hostCPUUsed) / hsUsageSummary.CPUTotal) * 100
			fmt.Fprintf(
				&report,
				"* %s (CPU: %s, Host CPU Usage: %2.2f%%)%s",
				vm.Name,
				CPUSpeed(hostCPUUsed),
				vmPercentOfHostCPUUsed,
				nagios.CheckOutputEOL,
			)
		}
	}

	if vmsPoweredOn == 0 {
		fmt.Fprintf(
			&report,
			"* None (visible)%s",
			nagios.CheckOutputEOL,
		)
	}

	fmt.Fprintf(
		&report,
		"%sVMs on host not consuming CPU:%s%s",
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
	)

	sort.Slice(hsVMs, func(i, j int) bool {
		return strings.ToLower(hsVMs[i].Name) < strings.ToLower(hsVMs[j].Name)
	})

	for _, vm := range hsVMs {
		if vm.Runtime.PowerState != types.VirtualMachinePowerStatePoweredOn {
			fmt.Fprintf(
				&report,
				"* %s%s",
				vm.Name,
				nagios.CheckOutputEOL,
			)
		}
	}

	if vmsPoweredOff == 0 {
		fmt.Fprintf(
			&report,
			"* None (visible)%s",
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

	return report.String()
}
