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
	"github.com/vmware/govmomi/vim25/types"
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

// GetNumTotalRPs returns the count of all Resource Pools in the inventory.
func GetNumTotalRPs(ctx context.Context, client *vim25.Client) (int, error) {
	funcTimeStart := time.Now()

	var numAllRPs int

	defer func(allRPs *int) {
		logger.Printf(
			"It took %v to execute GetNumTotalRPs func (and count %d RPs).\n",
			time.Since(funcTimeStart),
			*allRPs,
		)
	}(&numAllRPs)

	var getRPsErr error
	numAllRPs, getRPsErr = getRPsCountUsingContainerView(
		ctx,
		client,
		client.ServiceContent.RootFolder,
		true,
	)
	if getRPsErr != nil {
		logger.Printf(
			"error retrieving list of all resource pools: %v",
			getRPsErr,
		)

		return 0, fmt.Errorf(
			"error retrieving list of all resource pools: %w",
			getRPsErr,
		)
	}
	logger.Printf(
		"Finished retrieving count of all resource pools: %d",
		numAllRPs,
	)

	return numAllRPs, nil
}

// validateRPs verifies that all explicitly specified ResourcePools exist in
// the inventory.
func validateRPs(ctx context.Context, client *vim25.Client, filterOptions VMsFilterOptions) error {
	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute validateRPs func.\n",
			time.Since(funcTimeStart),
		)
	}()

	switch {
	case len(filterOptions.FoldersIncluded) > 0 || len(filterOptions.FoldersExcluded) > 0:
		logger.Println("Validating resource pools")

		validateErr := ValidateRPs(ctx, client, filterOptions.ResourcePoolsIncluded, filterOptions.ResourcePoolsExcluded)
		if validateErr != nil {
			logger.Printf(
				"%v: %v",
				ErrValidationOfIncludeExcludeRPLists,
				validateErr,
			)

			return fmt.Errorf(
				"%v: %v",
				ErrValidationOfIncludeExcludeRPLists,
				validateErr,
			)
		}
		logger.Println("Successfully validated resource pools")

		return nil
	default:
		logger.Println("Skipping resource pool validation; resource pool filtering not requested")
		return nil
	}
}

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

// getRPsCountUsingContainerView accepts a context, a connected client, a
// container type ManagedObjectReference and a boolean value indicating
// whether the container type should be recursively searched for
// ResourcePools. An error is returned if the provided ManagedObjectReference
// is not for a supported container type.
func getRPsCountUsingContainerView(
	ctx context.Context,
	c *vim25.Client,
	containerRef types.ManagedObjectReference,
	recursive bool,
) (int, error) {

	funcTimeStart := time.Now()

	var allRPs []types.ObjectContent

	defer func(rps *[]types.ObjectContent, objRef types.ManagedObjectReference) {
		logger.Printf(
			"It took %v to execute getRPsCountUsingContainerView func (and count %d RPs from %s).\n",
			time.Since(funcTimeStart),
			len(*rps),
			objRef.Type,
		)
	}(&allRPs, containerRef)

	// Create a view of caller-specified objects
	m := view.NewManager(c)

	logger.Printf("Container type is %s", containerRef.Type)

	// https://vdc-download.vmware.com/vmwb-repository/dcr-public/a5f4000f-1ea8-48a9-9221-586adff3c557/7ff50256-2cf2-45ea-aacd-87d231ab1ac7/vim.view.ContainerView.html
	switch containerRef.Type {
	case MgObjRefTypeResourcePool:
	case MgObjRefTypeFolder:

	default:
		return 0, fmt.Errorf(
			"unsupported container type specified for ContainerView: %s",
			containerRef.Type,
		)
	}

	kind := []string{MgObjRefTypeResourcePool}

	// FIXME: Should this filter to a specific datacenter? See GH-219.
	v, createViewErr := m.CreateContainerView(
		ctx,
		containerRef,
		kind,
		recursive,
	)
	if createViewErr != nil {
		return 0, createViewErr
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

	// Perform as lightweight of a search as possible as we're only interested
	// in counting the total resource pools in a specified container.
	prop := []string{"overallStatus"}
	retrieveErr := v.Retrieve(ctx, kind, prop, &allRPs)
	if retrieveErr != nil {
		return 0, retrieveErr
	}

	return len(allRPs), nil
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
	vmsFilterResults VMsFilterResults,
	aggregateMemoryUsageInBytes int64,
	maxMemoryUsageInBytes int64,
	clusterMemoryInBytes int64,
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
			vmsFilterResults.NumRPsAfterFiltering(),
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
			vmsFilterResults.NumRPsAfterFiltering(),
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
	vmsFilterOptions VMsFilterOptions,
	vmsFilterResults VMsFilterResults,
	maxMemoryUsageInBytes int64,
	clusterMemoryInBytes int64,
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

	_, _ = fmt.Fprintf(
		&report,
		"Memory usage by Resource Pool:%s%s",
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
	)
	for _, rp := range vmsFilterResults.RPsAfterFiltering() {

		// gather MOID to Name mappings for later lookup
		rpIDtoNameIdx[rp.Self.Value] = rp.Name

		rpSummary := rp.Summary.GetResourcePoolSummary()
		switch {
		case rpSummary == nil:
			_, _ = fmt.Fprintf(
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
			_, _ = fmt.Fprintf(
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

	vms := vmsFilterResults.VMsAfterFiltering()

	// TODO: We already have these values in vmsFilterResults, provided that
	// we assume that the caller has already filtered out powered off VMs.
	//
	// Since it is possible (however unlikely) that this function will be
	// called by another plugin, it might be worth performing this separate
	// filtering step just to be sure.
	poweredVMs, numVMsPoweredOff := FilterVMsByPowerState(vms, false)
	numVMsPoweredOn := len(poweredVMs)

	_, _ = fmt.Fprintf(
		&report,
		"%sTen VMS consuming most memory:%s%s",
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
	)

	switch {
	case numVMsPoweredOn == 0:
		_, _ = fmt.Fprintf(
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

			_, _ = fmt.Fprintf(
				&report,
				"* %s [Mem: %s, Pool: %s]%s",
				vm.Name,
				units.ByteSize(hostMemUsedBytes),
				rpName,
				nagios.CheckOutputEOL,
			)
		}

	}

	_, _ = fmt.Fprintf(
		&report,
		"%sTen VMs most recently powered on:%s%s",
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
	)

	switch {
	case len(poweredVMs) == 0:
		_, _ = fmt.Fprintf(
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

			_, _ = fmt.Fprintf(
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

	_, _ = fmt.Fprintf(
		&report,
		"%s---%s%s",
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
	)

	_, _ = fmt.Fprintf(
		&report,
		"* vSphere environment: %s%s",
		c.URL().String(),
		nagios.CheckOutputEOL,
	)

	_, _ = fmt.Fprintf(
		&report,
		"* Plugin User Agent: %s%s",
		c.Client.UserAgent,
		nagios.CheckOutputEOL,
	)

	_, _ = fmt.Fprintf(
		&report,
		"* Specified Resource Pools to explicitly include (%d): [%v]%s",
		len(vmsFilterOptions.ResourcePoolsIncluded),
		strings.Join(vmsFilterOptions.ResourcePoolsIncluded, ", "),
		nagios.CheckOutputEOL,
	)

	_, _ = fmt.Fprintf(
		&report,
		"* Specified Resource Pools to explicitly exclude (%d): [%v]%s",
		len(vmsFilterOptions.ResourcePoolsExcluded),
		strings.Join(vmsFilterOptions.ResourcePoolsExcluded, ", "),
		nagios.CheckOutputEOL,
	)

	_, _ = fmt.Fprintf(
		&report,
		"* Resource Pools evaluated (%d of %d): [%v]%s",
		vmsFilterResults.NumRPsAfterFiltering(),
		vmsFilterResults.NumRPsAll(),
		strings.Join(vmsFilterResults.RPNamesAfterFiltering(), ", "),
		nagios.CheckOutputEOL,
	)

	return report.String()
}
