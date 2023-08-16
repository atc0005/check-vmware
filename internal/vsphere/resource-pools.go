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

	"github.com/atc0005/check-vmware/internal/textutils"

	"github.com/atc0005/go-nagios"

	"github.com/vmware/govmomi/units"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
)

// ResourcePoolsAggregateStats is a collection of aggregated statistics for
// one or more Resource Pools.
type ResourcePoolsAggregateStats struct {
	// MemoryUsageInBytes is the consumed host memory in bytes for one or more
	// specified Resource Pools.
	MemoryUsageInBytes int64

	// BalloonedMemoryInBytes is the size of the balloon driver in bytes
	// across all virtual machines in one or more specified Resource Pools.
	// The host will inflate the balloon driver to reclaim physical memory
	// from a virtual machine. This is a sign that there is memory pressure on
	// the host.
	BalloonedMemoryInBytes int64

	// SwappedMemoryInBytes is the the portion of memory in bytes that is granted
	// to virtual machines from the host's swap space. This is a sign that
	// there is memory pressure on the host.
	SwappedMemoryInBytes int64
}

// ErrResourcePoolMemoryUsageThresholdCrossed indicates that specified
// resource pools have exceeded a given threshold
var ErrResourcePoolMemoryUsageThresholdCrossed = errors.New("memory usage exceeds specified threshold")

// ErrResourcePoolStatisticUnavailable indicates that one or more statistics
// are missing from specified Resource Pools. This is usually due to
// retrieving an insufficient subset of properties from a vSphere View.
var ErrResourcePoolStatisticUnavailable = errors.New("resource pool missing expected statistic")

// BalloonedMemoryHR returns the size of the balloon driver across all virtual
// machines in one or more specified Resource Pools as a human readable
// string.
func (rps ResourcePoolsAggregateStats) BalloonedMemoryHR() string {
	return units.ByteSize(rps.BalloonedMemoryInBytes).String()
}

// SwappedMemoryHR returns the portion of memory granted to all virtual
// machines from the host's swap space across all virtual machines in one or
// more specified Resource Pools as a human readable string.
func (rps ResourcePoolsAggregateStats) SwappedMemoryHR() string {
	return units.ByteSize(rps.SwappedMemoryInBytes).String()
}

// MemoryUsageHR returns the consumed host memory for one or more specified
// Resource Pools as a human readable string.
func (rps ResourcePoolsAggregateStats) MemoryUsageHR() string {
	return units.ByteSize(rps.MemoryUsageInBytes).String()
}

// ValidateRPs is responsible for receiving two lists of resource pools,
// explicitly "included" (aka, "whitelisted") and explicitly "excluded" (aka,
// "blacklisted"). If any list entries are not found in the vSphere
// environment an error is returned listing which ones.
func ValidateRPs(ctx context.Context, c *vim25.Client, includeRPs []string, excludeRPs []string) error {

	funcTimeStart := time.Now()

	defer func(irps []string, erps []string) {
		logger.Printf(
			"It took %v to execute ValidateRPs func (and validate %d Resource Pools).\n",
			time.Since(funcTimeStart),
			len(irps)+len(erps),
		)
	}(includeRPs, excludeRPs)

	m := view.NewManager(c)

	// Create a view of Resource Pool objects
	v, createViewErr := m.CreateContainerView(
		ctx,
		c.ServiceContent.RootFolder,
		[]string{MgObjRefTypeResourcePool},
		true,
	)
	if createViewErr != nil {
		return fmt.Errorf("failed to create ResourcePool view: %w", createViewErr)
	}

	defer func() {
		// Per vSphere Web Services SDK Programming Guide - VMware vSphere 7.0
		// Update 1:
		//
		// A best practice when using views is to call the DestroyView()
		// method when a view is no longer needed. This practice frees memory
		// on the server.
		if err := v.Destroy(ctx); err != nil {
			logger.Printf("Error occurred while destroying view: %s", err)
		}
	}()

	// Retrieve name property for all resource pools.
	props := []string{"name"}
	var rpsSearchResults []mo.ResourcePool
	retrieveErr := v.Retrieve(ctx, []string{MgObjRefTypeResourcePool}, props, &rpsSearchResults)
	if retrieveErr != nil {
		return fmt.Errorf(
			"failed to retrieve ResourcePool properties: %w",
			retrieveErr,
		)
	}

	// We're only interested in working with resource pool names
	poolNamesFound := make([]string, 0, len(rpsSearchResults))
	for _, rp := range rpsSearchResults {
		poolNamesFound = append(poolNamesFound, rp.Name)
	}

	// If any specified resource pool names are not found, note that so we can
	// provide the full list of invalid pool names together as a convenience
	// for the user.
	var notFound []string
	switch {
	case len(includeRPs) > 0:
		for _, iRP := range includeRPs {
			if !textutils.InList(iRP, poolNamesFound, true) {
				notFound = append(notFound, iRP)
			}
		}

		if len(notFound) > 0 {
			return fmt.Errorf(
				"specified Resource Pools (to include) not found: %v",
				notFound,
			)
		}

		// all listed resource pools were found
		return nil

	case len(excludeRPs) > 0:
		for _, eRP := range excludeRPs {
			if !textutils.InList(eRP, poolNamesFound, true) {
				notFound = append(notFound, eRP)
			}
		}

		if len(notFound) > 0 {
			return fmt.Errorf(
				"specified Resource Pools (to exclude) not found: %v",
				notFound,
			)
		}

		// all listed resource pools were found
		return nil

	default:

		// no restrictions specified by user; all resource pools are
		// "eligible" for evaluation
		return nil
	}

}

// GetEligibleRPs receives a list of Resource Pool names that should either be
// explicitly included or excluded along with a boolean value indicating
// whether only a subset of properties for the Resource Pools should be
// returned. If requested, a subset of all available properties will be
// retrieved (faster) instead of recursively fetching all properties (about 2x
// as slow). The filtered list of Resource Pools is returned, or an error if
// one occurs.
func GetEligibleRPs(ctx context.Context, c *vim25.Client, includeRPs []string, excludeRPs []string, propsSubset bool) ([]mo.ResourcePool, error) {

	funcTimeStart := time.Now()

	// Declare slice early so that we can grab a pointer to it in order to
	// access the entries later. This holds the filtered list of resource
	// pools that will be returned to the caller.
	var rps []mo.ResourcePool

	defer func(rps *[]mo.ResourcePool) {
		logger.Printf(
			"It took %v to execute GetEligibleRPs func (and retrieve %d Resource Pools).\n",
			time.Since(funcTimeStart),
			len(*rps),
		)
	}(&rps)

	// All available/accessible resource pools will be retrieved and stored
	// here. We will filter the results before returning a trimmed list to the
	// caller.
	var rpsSearchResults []mo.ResourcePool

	err := getObjects(ctx, c, &rpsSearchResults, c.ServiceContent.RootFolder, propsSubset, true)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve ResourcePools: %w", err)
	}

	rpNames := make([]string, 0, len(rpsSearchResults))
	for _, rp := range rpsSearchResults {
		rpNames = append(rpNames, rp.Name)
	}

	logger.Printf(
		"Retrieved %d ResourcePool objects: %v",
		len(rpsSearchResults),
		strings.Join(rpNames, ", "),
	)

	for _, rp := range rpsSearchResults {

		// config validation asserts that only one of include/exclude resource
		// pools flags are specified
		switch {

		// if specified, only include resource pools that have been
		// intentionally included (aka, "whitelisted")
		case len(includeRPs) > 0:
			if textutils.InList(rp.Name, includeRPs, true) {
				rps = append(rps, rp)
			}

		// if specified, don't include resource pools that have been
		// intentionally excluded (aka, "blacklisted")
		case len(excludeRPs) > 0:
			if !textutils.InList(rp.Name, excludeRPs, true) {
				rps = append(rps, rp)
			}

		// if we are not explicitly excluding or including pools, then we are
		// working with all pools
		default:
			rps = append(rps, rp)
		}

	}

	sort.Slice(rps, func(i, j int) bool {
		return strings.ToLower(rps[i].Name) < strings.ToLower(rps[j].Name)
	})

	return rps, nil

}

// GetRPByName accepts the name of a Resource Pool, the name of a datacenter
// and a boolean value indicating whether only a subset of properties for the
// Network should be returned. If requested, a subset of all available
// properties will be retrieved (faster) instead of recursively fetching all
// properties (about 2x as slow). If the datacenter name is an empty string
// then the default datacenter will be used.
func GetRPByName(ctx context.Context, c *vim25.Client, rpName string, datacenter string, propsSubset bool) (mo.ResourcePool, error) {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute GetRPByName func.\n",
			time.Since(funcTimeStart),
		)
	}()

	var rPool mo.ResourcePool
	err := getObjectByName(ctx, c, &rPool, rpName, datacenter, propsSubset)

	if err != nil {
		return mo.ResourcePool{}, err
	}

	return rPool, nil

}

// ResourcePoolStats receives a collection of ResourcePool values and returns
// a collection of aggregate statistics (e.g., memory usage, ballooned memory,
// swapped memory, etc.). An error is returned if required properties are
// missing for one or more of the ResourcePool values and an initial attempt
// to populate the properties fails.
func ResourcePoolStats(ctx context.Context, client *vim25.Client, resourcePools []mo.ResourcePool) (ResourcePoolsAggregateStats, error) {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute ResourcePoolStats func (and process %d resource pools).\n",
			time.Since(funcTimeStart),
			len(resourcePools),
		)
	}()

	var aggregateBalloonedMemoryInBytes int64
	var aggregateSwappedMemoryInBytes int64
	var aggregateMemoryUsageInBytes int64

	for _, rp := range resourcePools {

		rpSummary := rp.Summary.GetResourcePoolSummary()

		// If required ResourcePool summary and quickStats properties are not
		// populated, trigger a state reload in an attempt to populate them.
		// Return an error if the attempt fails.
		if rpSummary == nil || rpSummary.QuickStats == nil {

			logger.Printf(
				"Required statistics unavailable for ResourcePool %q; "+
					"triggering state reload",
				rp.Name,
			)

			if err := TriggerEntityStateReload(ctx, client, rp.ManagedEntity); err != nil {
				return ResourcePoolsAggregateStats{}, fmt.Errorf(
					"failed to reload state for resource pool %q: %w",
					rp.Name,
					ErrEntityStateReloadUnsuccessful,
				)
			}

			logger.Printf(
				"State reload successfully triggered for ResourcePool %q",
				rp.Name,
			)

			logger.Print("Rechecking statistics availability")
		}

		// If summary and quickStats properties are *still* unpopulated,
		// return an error.
		//
		// TODO: Annotate error with additional details to help sysadmin
		// enable required settings in vSphere to ensure that the needed
		// properties are populated.
		switch {
		case rpSummary == nil:
			return ResourcePoolsAggregateStats{}, fmt.Errorf(
				"failed to retrieve summary property for resource pool %q: %w",
				rp.Name,
				ErrResourcePoolStatisticUnavailable,
			)
		case rpSummary.QuickStats == nil:
			return ResourcePoolsAggregateStats{}, fmt.Errorf(
				"failed to retrieve quickstats property for resource pool %q: %w",
				rp.Name,
				ErrResourcePoolStatisticUnavailable,
			)
		default:
			logger.Printf("Statistics available for ResourcePool %q", rp.Name)
		}

		// Per vSphere API docs, `rp.Runtime.Memory.OverallUsage` was
		// deprecated in v6.5, so we use `hostMemoryUsage` instead.
		//
		// The `hostMemoryUsage` property tracks consumed host memory in MB.
		// This includes the overhead memory of a virtual machine. We multiply
		// by units.MB in order to get the number of bytes in order to match
		// the same unit of measurement used by the `host.Hardware.MemorySize`
		// property.
		rpMemoryUsage := rpSummary.QuickStats.HostMemoryUsage * units.MB
		aggregateMemoryUsageInBytes += rpMemoryUsage

		// The size of the balloon driver in a virtual machine, in MB. The
		// host will inflate the balloon driver to reclaim physical memory
		// from a virtual machine. This is a sign that there is memory
		// pressure on the host. We multiply by units.MB in order to get the
		// number of bytes.
		rpBalloonedMemory := rpSummary.QuickStats.BalloonedMemory * units.MB
		aggregateBalloonedMemoryInBytes += rpBalloonedMemory

		rpSwappedMemory := rpSummary.QuickStats.SwappedMemory * units.MB
		aggregateSwappedMemoryInBytes += rpSwappedMemory

		logger.Printf(
			"resource pool %q (memory usage: %s, ballooned memory: %s, swapped memory: %s)",
			rp.Name,
			units.ByteSize(rpMemoryUsage).String(),
			units.ByteSize(rpBalloonedMemory).String(),
			units.ByteSize(rpSwappedMemory).String(),
		)
	}

	stats := ResourcePoolsAggregateStats{
		MemoryUsageInBytes:     aggregateMemoryUsageInBytes,
		BalloonedMemoryInBytes: aggregateBalloonedMemoryInBytes,
		SwappedMemoryInBytes:   aggregateSwappedMemoryInBytes,
	}

	return stats, nil

}

// MemoryUsedPercentage is a helper function used to calculate the current
// memory usage as a percentage of the specified maximum memory allowed to be
// used.
func MemoryUsedPercentage(
	aggregateMemoryUsageInBytes int64,
	maxMemoryUsageInBytes int64,
) float64 {
	memoryPercentageUsedOfAllowed := (float64(aggregateMemoryUsageInBytes) / float64(maxMemoryUsageInBytes)) * 100

	return memoryPercentageUsedOfAllowed
}

// RPMemoryUsageOneLineCheckSummary is used to generate a one-line Nagios
// service check results summary. This is the line most prominent in
// notifications.
func RPMemoryUsageOneLineCheckSummary(
	stateLabel string,
	aggregateMemoryUsageInBytes int64,
	maxMemoryUsageInBytes int64,
	clusterMemoryInBytes int64,
	rps []mo.ResourcePool,
) string {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute RPMemoryUsageOneLineCheckSummary func.\n",
			time.Since(funcTimeStart),
		)
	}()

	memoryPercentageUsedOfAllowed := MemoryUsedPercentage(aggregateMemoryUsageInBytes, maxMemoryUsageInBytes)
	memoryPercentageUsedOfClusterCapacity := MemoryUsedPercentage(
		aggregateMemoryUsageInBytes,
		clusterMemoryInBytes,
	)

	switch {

	case aggregateMemoryUsageInBytes > maxMemoryUsageInBytes:
		return fmt.Sprintf(
			"%s: %s (%.1f%%) memory used of %s allowed, "+
				"%.2f%% of %s total capacity (evaluated %d Resource Pools)",
			stateLabel,
			units.ByteSize(aggregateMemoryUsageInBytes),
			memoryPercentageUsedOfAllowed,
			units.ByteSize(maxMemoryUsageInBytes),
			memoryPercentageUsedOfClusterCapacity,
			units.ByteSize(clusterMemoryInBytes),
			len(rps),
		)

	default:
		memoryRemaining := maxMemoryUsageInBytes - aggregateMemoryUsageInBytes
		return fmt.Sprintf(
			"%s: %s memory used (%0.1f%%), %.2f%% of %s total capacity; "+
				"%s (%0.1f%%) of %s remaining "+
				"(evaluated %d Resource Pools)",
			stateLabel,
			units.ByteSize(aggregateMemoryUsageInBytes),
			memoryPercentageUsedOfAllowed,
			memoryPercentageUsedOfClusterCapacity,
			units.ByteSize(clusterMemoryInBytes),
			units.ByteSize(memoryRemaining),
			float64(100)-memoryPercentageUsedOfAllowed,
			units.ByteSize(maxMemoryUsageInBytes),
			len(rps),
		)

	}
}

// ResourcePoolsMemoryReport generates a summary of memory usage associated
// with specified Resource Pools along with various verbose details intended
// to aid in troubleshooting check results at a glance. This information is
// provided for use with the Long Service Output field commonly displayed on
// the detailed service check results display in the web UI or in the body of
// many notifications.
func ResourcePoolsMemoryReport(
	c *vim25.Client,
	maxMemoryUsageInBytes int64,
	clusterMemoryInBytes int64,
	includeRPs []string,
	excludeRPs []string,
	rps []mo.ResourcePool,
	rpsVMs []mo.VirtualMachine,
) string {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute ResourcePoolsMemoryReport func.\n",
			time.Since(funcTimeStart),
		)
	}()

	var report strings.Builder

	rpIDtoNameIdx := make(map[string]string)

	fmt.Fprintf(
		&report,
		"Memory usage by Resource Pool:%s%s",
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
	)
	for _, rp := range rps {

		// gather MOID to Name mappings for later lookup
		rpIDtoNameIdx[rp.Self.Value] = rp.Name

		rpSummary := rp.Summary.GetResourcePoolSummary()
		switch {
		case rpSummary == nil:
			fmt.Fprintf(
				&report,
				"* %s [Pool: (unavailable), Cluster: (unavailable)]%s",
				rp.Name,
				nagios.CheckOutputEOL,
			)

		default:
			rpMemoryUsage := rpSummary.QuickStats.HostMemoryUsage * units.MB
			rpMemoryPercentageUsed := MemoryUsedPercentage(rpMemoryUsage, maxMemoryUsageInBytes)
			memoryPercentageUsedOfClusterCapacity := MemoryUsedPercentage(
				rpMemoryUsage,
				clusterMemoryInBytes,
			)
			fmt.Fprintf(
				&report,
				"* %s [Pool: (%s, %0.1f%%), Cluster: (%.2f%%)]%s",
				rp.Name,
				units.ByteSize(rpMemoryUsage),
				rpMemoryPercentageUsed,
				memoryPercentageUsedOfClusterCapacity,
				nagios.CheckOutputEOL,
			)
		}
	}

	poweredVMs, numVMsPoweredOff := FilterVMsByPowerState(rpsVMs, false)
	numVMsPoweredOn := len(poweredVMs)

	fmt.Fprintf(
		&report,
		"%sTen VMS consuming most memory:%s%s",
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
	)

	switch {
	case numVMsPoweredOn == 0:
		fmt.Fprintf(
			&report,
			"* None (visible); %d powered off%s",
			numVMsPoweredOff,
			nagios.CheckOutputEOL,
		)

	default:

		sort.Slice(poweredVMs, func(i, j int) bool {
			return poweredVMs[i].Summary.QuickStats.HostMemoryUsage > poweredVMs[j].Summary.QuickStats.HostMemoryUsage
		})

		// grab up to the first 10 VMs, presorted by most memory usage
		sampleSize := len(poweredVMs)
		if sampleSize > 10 {
			sampleSize = 10
		}

		for _, vm := range poweredVMs[:sampleSize] {
			hostMemUsedBytes := int64(vm.Summary.QuickStats.HostMemoryUsage) * units.MB
			rpName := rpIDtoNameIdx[vm.ResourcePool.Value]

			fmt.Fprintf(
				&report,
				"* %s [Mem: %s, Pool: %s]%s",
				vm.Name,
				units.ByteSize(hostMemUsedBytes),
				rpName,
				nagios.CheckOutputEOL,
			)
		}

	}

	fmt.Fprintf(
		&report,
		"%sTen VMs most recently powered on:%s%s",
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
	)

	switch {
	case len(poweredVMs) == 0:
		fmt.Fprintf(
			&report,
			"* None (visible); %d powered off%s",
			numVMsPoweredOff,
			nagios.CheckOutputEOL,
		)

	default:

		sort.Slice(poweredVMs, func(i, j int) bool {
			return poweredVMs[i].Summary.QuickStats.UptimeSeconds < poweredVMs[j].Summary.QuickStats.UptimeSeconds
		})

		// grab up to the first 10 VMs, presorted by least uptime
		sampleSize := len(poweredVMs)
		if sampleSize > 10 {
			sampleSize = 10
		}

		for _, vm := range poweredVMs[:sampleSize] {
			hostMemUsedBytes := int64(vm.Summary.QuickStats.HostMemoryUsage) * units.MB
			uptime := time.Duration(vm.Summary.QuickStats.UptimeSeconds) * time.Second
			uptimeDays := uptime.Hours() / 24
			rpName := rpIDtoNameIdx[vm.ResourcePool.Value]

			fmt.Fprintf(
				&report,
				"* %s: [Uptime: %.2f days, Mem: %s, Pool: %s]%s",
				vm.Name,
				uptimeDays,
				units.ByteSize(hostMemUsedBytes),
				rpName,
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

	rpNames := make([]string, len(rps))
	for i := range rps {
		rpNames[i] = rps[i].Name
	}

	fmt.Fprintf(
		&report,
		"* Resource Pools evaluated (%d): [%v]%s",
		len(rpNames),
		strings.Join(rpNames, ", "),
		nagios.CheckOutputEOL,
	)

	return report.String()
}
