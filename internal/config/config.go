// Copyright 2021 Adam Chalkley
//
// https://github.com/atc0005/check-vmware
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
)

// Updated via Makefile builds. Setting placeholder value here so that
// something resembling a version string will be provided for non-Makefile
// builds.
var version = "x.y.z"

// ErrVersionRequested indicates that the user requested application version
// information.
var ErrVersionRequested = errors.New("version information requested")

// PluginType represents the type of plugin that is being
// configured/initialized. Not all plugin types will use the same features and
// as a result will not accept the same flags. Unless noted otherwise, each of
// the plugin types are incompatible with each other, though some flags are
// common to all types. See also the PluginType* constants.
type PluginType struct {
	Tools                          bool
	SnapshotsAge                   bool
	SnapshotsCount                 bool
	SnapshotsSize                  bool
	DatastoresSpace                bool
	DatastoresPerformance          bool
	ResourcePoolsMemory            bool
	VirtualCPUsAllocation          bool
	VirtualHardwareVersion         bool
	Host2Datastores2VMs            bool
	HostSystemMemory               bool
	HostSystemCPU                  bool
	VirtualMachinePowerCycleUptime bool
	DiskConsolidation              bool
	InteractiveQuestion            bool
	Alarms                         bool

	// TODO:
	// - vCenter/server time (NTP)

}

// AppInfo identifies common details about the plugins provided by this
// project.
type AppInfo struct {

	// Name specifies the public name shared by all plugins in this project.
	Name string

	// Version specifies the public version shared by all plugins in this
	// project.
	Version string

	// URL specifies the public repo URL shared by all plugins in this
	// project.
	URL string

	// Plugin indicates which plugin provided by this project is currently
	// executing.
	Plugin string
}

// multiValueStringFlag is a custom type that satisfies the flag.Value
// interface in order to accept multiple string values for some of our flags.
type multiValueStringFlag []string

// String returns a comma separated string consisting of all slice elements.
func (mvs *multiValueStringFlag) String() string {

	// The String() method is called by the flag.isZeroValue function in order
	// to determine whether the output string represents the zero value for a
	// flag. This occurs even if the flag is not specified by the user.

	// From the `flag` package docs:
	// "The flag package may call the String method with a zero-valued
	// receiver, such as a nil pointer."
	if mvs == nil {
		return ""
	}

	return strings.Join(*mvs, ", ")
}

// Set is called once by the flag package, in command line order, for each
// flag present.
func (mvs *multiValueStringFlag) Set(value string) error {

	// split comma-separated string into multiple values, toss leading and
	// trailing whitespace
	items := strings.Split(value, ",")
	for index, item := range items {
		items[index] = strings.TrimSpace(item)
		items[index] = strings.ReplaceAll(items[index], "'", "")
		items[index] = strings.ReplaceAll(items[index], "\"", "")
	}

	// add them to the collection
	*mvs = append(*mvs, items...)

	return nil
}

// DSPerformanceSummaryThresholds represents the thresholds used to evaluate
// Datastore Performance Summary values.
type DSPerformanceSummaryThresholds struct {
	// ReadLatencyWarning is the read latency in ms when a WARNING threshold
	// is reached.
	ReadLatencyWarning float64

	// ReadLatencyCritical is the read latency in ms when a CRITICAL threshold
	// is reached.
	ReadLatencyCritical float64

	// WriteLatencyWarning is the write latency in ms when a WARNING threshold
	// is reached.
	WriteLatencyWarning float64

	// WriteLatencyCritical is the write latency in ms when a CRITICAL
	// threshold is reached.
	WriteLatencyCritical float64

	// VMLatencyWarning is the latency in ms as observed by VMs using the
	// datastore when a WARNING threshold is reached.
	VMLatencyWarning float64

	// VMLatencyCritical is the latency in ms as observed by VMs using the
	// datastore when a CRITICAL threshold is reached.
	VMLatencyCritical float64
}

// dsPerfLatencyMetricFlag is a custom type that satisfies the flag.Value
// interface. This type is used to accept Datastore Performance Summary
// latency metric values. This flag type is incompatible with the flag type
// used to specify percentile sets.
type dsPerfLatencyMetricFlag struct {

	// value is the user-specified value
	value float64

	// isSet identifies whether a value was provided by the user
	isSet bool
}

// String satisfies the flag.Value interface method set requirements.
func (dspl *dsPerfLatencyMetricFlag) String() string {

	// The String() method is called by the flag.isZeroValue function in order
	// to determine whether the output string represents the zero value for a
	// flag. This occurs even if the flag is not specified by the user.

	if dspl == nil {
		return ""
	}

	return fmt.Sprintf(
		"value: %v, isSet: %v",
		dspl.value,
		dspl.isSet,
	)

}

// Set satisfies the flag.Value interface method set requirements.
func (dspl *dsPerfLatencyMetricFlag) Set(value string) error {

	// fmt.Println("dsPerfLatencyMetricFlag Set() called")

	var strConvErr error

	value = strings.TrimSpace(value)
	value = strings.ReplaceAll(value, "'", "")
	value = strings.ReplaceAll(value, "\"", "")

	var parsedVal float64
	parsedVal, strConvErr = strconv.ParseFloat(value, 64)
	if strConvErr != nil {
		return fmt.Errorf(
			"error processing flag; failed to convert %q: %v",
			value,
			strConvErr,
		)
	}

	dspl.value = parsedVal
	dspl.isSet = true

	return nil

}

// MultiValueDSPerfPercentileSetFlag is a custom type that satisfies the
// flag.Value interface. This type is used to accept Datastore Performance
// Summary percentile "sets". These sets define thresholds used to check
// Datastore Performance latency metrics to determine overall plugin state.
type MultiValueDSPerfPercentileSetFlag map[int]DSPerformanceSummaryThresholds

// String returns a comma separated string consisting of all map entries.
func (mvdsperf *MultiValueDSPerfPercentileSetFlag) String() string {

	// The String() method is called by the flag.isZeroValue function in order
	// to determine whether the output string represents the zero value for a
	// flag. This occurs even if the flag is not specified by the user.

	// From the `flag` package docs:
	// "The flag package may call the String method with a zero-valued
	// receiver, such as a nil pointer."
	if mvdsperf == nil {
		return "empty percentile set"
	}

	percentiles := make([]int, 0, len(*mvdsperf))
	for key := range *mvdsperf {
		percentiles = append(percentiles, key)
	}

	sort.Slice(percentiles, func(i, j int) bool {
		return percentiles[i] < percentiles[j]
	})

	var output strings.Builder

	for _, p := range percentiles {
		fmt.Fprintf(&output,
			"{Percentile: %v, ThresholdVals: %+v}, ",
			p,
			(*mvdsperf)[p],
		)
	}

	outputString := strings.TrimSuffix(output.String(), ", ")

	return outputString

}

// thresholdValues receives a string indicating either WARNING or CRITICAL
// state and returns a comma separated string consisting of all specified
// metric percentiles and the the associated WARNING or CRITICAL threshold
// values.
func (mvdsperf MultiValueDSPerfPercentileSetFlag) thresholdValues(state string) string {

	// From the `flag` package docs:
	// "The flag package may call the String method with a zero-valued
	// receiver, such as a nil pointer."
	if mvdsperf == nil {
		return "empty percentile set"
	}

	percentiles := make([]int, 0, len(mvdsperf))
	for key := range mvdsperf {
		percentiles = append(percentiles, key)
	}

	sort.Slice(percentiles, func(i, j int) bool {
		return percentiles[i] < percentiles[j]
	})

	var output strings.Builder

	var readLatency float64
	var writeLatency float64
	var vmLatency float64

	for _, p := range percentiles {

		switch {
		case strings.ToUpper(state) == StateCRITICALLabel:
			readLatency = mvdsperf[p].ReadLatencyCritical
			writeLatency = mvdsperf[p].WriteLatencyCritical
			vmLatency = mvdsperf[p].VMLatencyCritical

			// fmt.Printf(
			// 	"CRITICAL | readLatency: %v, writeLatency: %v, vmLatency: %v\n",
			// 	readLatency,
			// 	writeLatency,
			// 	vmLatency,
			// )

		case strings.ToUpper(state) == StateWARNINGLabel:
			readLatency = mvdsperf[p].ReadLatencyWarning
			writeLatency = mvdsperf[p].WriteLatencyWarning
			vmLatency = mvdsperf[p].VMLatencyWarning

			// fmt.Printf(
			// 	"WARNING | readLatency: %v, writeLatency: %v, vmLatency: %v\n",
			// 	readLatency,
			// 	writeLatency,
			// 	vmLatency,
			// )
		}

		fmt.Fprintf(&output,
			"{ Percentile: %v, ReadLatency: %+v, WriteLatency: %v, VMLatency: %v }, ",
			p,
			readLatency,
			writeLatency,
			vmLatency,
		)
	}

	outputString := strings.TrimSuffix(output.String(), ", ")

	return outputString

}

// CriticalThresholdValues returns a comma separated string consisting of all
// specified metric percentiles and the the associated CRITICAL threshold
// values.
func (mvdsperf MultiValueDSPerfPercentileSetFlag) CriticalThresholdValues() string {
	return mvdsperf.thresholdValues(StateCRITICALLabel)
}

// WarningThresholdValues returns a comma separated string consisting of all
// specified metric percentiles and the the associated WARNING threshold
// values.
func (mvdsperf MultiValueDSPerfPercentileSetFlag) WarningThresholdValues() string {
	return mvdsperf.thresholdValues(StateWARNINGLabel)
}

// Set is called once by the flag package, in command line order, for each
// flag present.
func (mvdsperf *MultiValueDSPerfPercentileSetFlag) Set(value string) error {

	// We require the same number of values as we have fields in the struct
	// plus one more to serve as the map index (perenctile).
	const expectedValues int = 7

	// Split comma-separated string into multiple values, toss whitespace,
	// then convert value in string format to integer.
	items := strings.Split(value, ",")

	if len(items) != expectedValues {
		return fmt.Errorf(
			"error processing flag; string %q provides %d values, expected %d values",
			value,
			len(items),
			expectedValues,
		)
	}

	// fmt.Println("items", items)

	percentileSet := make([]float64, len(items))
	var strConvErr error
	for i := range items {
		items[i] = strings.TrimSpace(items[i])
		items[i] = strings.ReplaceAll(items[i], "'", "")
		items[i] = strings.ReplaceAll(items[i], "\"", "")

		percentileSet[i], strConvErr = strconv.ParseFloat(strings.TrimSpace(items[i]), 64)
		if strConvErr != nil {
			return fmt.Errorf(
				"error processing flag; failed to convert %q: %v",
				items[i],
				strConvErr,
			)
		}

	}

	// We now have the latency values (along with the percentile) stored as
	// float64 values. The first element is the percentile which is an int.
	percentile := int(percentileSet[0])

	// fmt.Printf("mvdsperf before assignment to map: %+v (nil: %t)\n", mvdsperf, mvdsperf == nil)

	// The rest of the latency values have already been converted to the
	// necessary type, so we assign directly.
	(*mvdsperf)[percentile] = DSPerformanceSummaryThresholds{
		ReadLatencyWarning:   percentileSet[1],
		ReadLatencyCritical:  percentileSet[2],
		WriteLatencyWarning:  percentileSet[3],
		WriteLatencyCritical: percentileSet[4],
		VMLatencyWarning:     percentileSet[5],
		VMLatencyCritical:    percentileSet[6],
	}

	// 	fmt.Printf("mvdsperf[percentile]: %+v\n", mvdsperf[percentile])
	//
	// 	fmt.Printf("mvdsperf after assignment to map: %+v (nil: %t)\n", mvdsperf, mvdsperf == nil)

	return nil

}

// Config represents the application configuration as specified via
// command-line flags.
type Config struct {

	// Server is the fully-qualified domain name of the system running a
	// certificate-enabled service.
	Server string

	// Username is the user account used to login to the ESXi host or vCenter
	// instance.
	Username string

	// Password is associated with the account used to login to the ESXi host
	// or vCenter instance.
	Password string

	// Domain is the domain for the user account used to login to the ESXi
	// host or vCenter instance.
	Domain string

	// ClusterName is the name of the vSphere cluster where monitored objects
	// reside.
	ClusterName string

	// LoggingLevel is the supported logging level for this application.
	LoggingLevel string

	// hostCustomAttributeName is a Custom Attribute name specific to hosts.
	// If specified, the user must also specify the Custom Attribute name
	// specific to datastores.
	hostCustomAttributeName string

	// hostCustomAttributePrefixSeparator is a prefix separator for Custom
	// Attribute values specific to hosts. If specified, this separator is
	// used to split the value for the specified Custom Attribute. The first
	// element from the split value is used as the prefix when comparing
	// Custom Attribute values. Also if specified, the user must also specify
	// the Custom Attribute prefix separator specific to datastores.
	hostCustomAttributePrefixSeparator string

	// datastoreCustomAttributeName is a Custom Attribute name specific to
	// datastores. If specified, the user must also specify the Custom
	// Attribute name specific to hosts.
	datastoreCustomAttributeName string

	// datastoreCustomAttributePrefixSeparator is a prefix separator for
	// Custom Attribute values specific to datastores. If specified, this
	// separator is used to split the value for the specified Custom
	// Attribute. The first element from the split value is used as the prefix
	// when comparing Custom Attribute values. Also if specified, the user
	// must also specify the Custom Attribute prefix separator specific to
	// hosts.
	datastoreCustomAttributePrefixSeparator string

	// sharedCustomAttributeName is a Custom Attribute name shared by both
	// hosts and datastores. If specified, the user must not specify the
	// Custom Attribute name specific to hosts or datastores.
	sharedCustomAttributeName string

	// sharedCustomAttributePrefixSeparator is a prefix separator for Custom
	// Attribute values shared by both hosts and datastores. If specified,
	// this separator is used to split the value for the specified Custom
	// Attribute. The first element from the split value is used as the prefix
	// when comparing Custom Attribute values. If specified, the user must not
	// specify the Custom Attribute prefix separator specific to hosts or
	// datastores.
	sharedCustomAttributePrefixSeparator string

	// DatastoreName is the name of the datastore as it is found within the
	// vSphere inventory of the specified ESXi host or vCenter instance.
	DatastoreName string

	// DatacenterName is the name of a Datacenter in the associated vSphere
	// inventory. This field is used by plugins which support monitoring only
	// a single Datacenter. Not applicable to standalone ESXi hosts.
	DatacenterName string

	// DatacenterNames is the name of one or more Datacenters in the
	// associated vSphere inventory. This field is used by plugins which
	// support monitoring multiple Datacenters. Not applicable to standalone
	// ESXi hosts.
	DatacenterNames multiValueStringFlag

	// HostSystemName is the name of an ESXi host/server in the associated
	// vSphere inventory.
	HostSystemName string

	// IncludedResourcePools lists resource pools that are explicitly
	// monitored. Specifying list values automatically excludes VirtualMachine
	// objects outside a Resource Pool.
	IncludedResourcePools multiValueStringFlag

	// ExcludedResourcePools lists resource pools that are explicitly ignored
	// or excluded from being monitored.
	ExcludedResourcePools multiValueStringFlag

	// IgnoredVM is a list of VMs that are explicitly ignored or excluded
	// from being monitored.
	IgnoredVMs multiValueStringFlag

	// IgnoredDatastores is a list of datastore names for Datastores that are
	// allowed to be associated with a VirtualMachine that are not associated
	// with its current host.
	IgnoredDatastores multiValueStringFlag

	// IncludedAlarmEntityTypes is a list of entity types for Alarms that will
	// be explicitly included for evaluation. Unless included by later
	// filtering logic, unmatched Triggered Alarms will be excluded from final
	// evaluation. Explicitly included Triggered Alarms are still subject to
	// permanent exclusion if an explicit exclusion match is made.
	IncludedAlarmEntityTypes multiValueStringFlag

	// ExcludedAlarmEntityTypes is a list of entity types for Alarms that will
	// be explicitly excluded from further evaluation by other stages in the
	// filtering pipeline. Explicit exclusions have precedence over explicit
	// inclusions.
	ExcludedAlarmEntityTypes multiValueStringFlag

	// IncludedAlarmEntityNames is a list of entity names for Alarms that will
	// be explicitly included for evaluation. Unless included by later
	// filtering logic, unmatched Triggered Alarms will be excluded from final
	// evaluation. Explicitly included Triggered Alarms are still subject to
	// permanent exclusion if an explicit exclusion match is made.
	IncludedAlarmEntityNames multiValueStringFlag

	// ExcludedAlarmEntityTypes is a list of entity names for Alarms that will
	// be explicitly excluded from further evaluation by other stages in the
	// filtering pipeline. Explicit exclusions have precedence over explicit
	// inclusions.
	ExcludedAlarmEntityNames multiValueStringFlag

	// IncludedAlarmEntityResourcePools is a list of resource pools that are
	// compared against the name of a resource pool for an entity associated
	// with one or more Triggered Alarms. Any Triggered Alarm with an
	// associated entity that is part of one of these resource pools is
	// explicitly included for evaluation.
	//
	// Unless included by later filtering logic, unmatched Triggered Alarms
	// will be excluded from final evaluation. Explicitly included Triggered
	// Alarms are still subject to permanent exclusion if an explicit
	// exclusion match is made.
	IncludedAlarmEntityResourcePools multiValueStringFlag

	// ExcludedAlarmEntityTypes is a list of resource pools that are compared
	// against the name of a resource pool for an entity associated with one
	// or more Triggered Alarms. Any Triggered Alarm with an associated that
	// is NOT part of one of these resource pools will be explicitly excluded
	// from further evaluation by other stages in the filtering pipeline.
	// Explicit exclusions have precedence over explicit inclusions.
	ExcludedAlarmEntityResourcePools multiValueStringFlag

	// IncludedAlarmNames is a list of names for defined Alarms that will be
	// explicitly included for evaluation. Unless included by later filtering
	// logic, unmatched Triggered Alarms will be excluded from final
	// evaluation. Explicitly included Triggered Alarms are still subject to
	// permanent exclusion if an explicit exclusion match is made.
	IncludedAlarmNames multiValueStringFlag

	// ExcludedAlarmNames is a list of names for defined Alarms that will be
	// explicitly excluded from further evaluation by other stages in the
	// filtering pipeline. Explicit exclusions have precedence over explicit
	// inclusions.
	ExcludedAlarmNames multiValueStringFlag

	// IncludedAlarmDescriptions is a list of descriptions for defined Alarms
	// that will be explicitly included for evaluation. Unless included by
	// later filtering logic, unmatched Triggered Alarms will be excluded from
	// final evaluation. Explicitly included Triggered Alarms are still
	// subject to permanent exclusion if an explicit exclusion match is made.
	IncludedAlarmDescriptions multiValueStringFlag

	// ExcludedAlarmDescriptions is a list of descriptions for defined Alarms
	// that will be explicitly excluded from further evaluation by other
	// stages in the filtering pipeline. Explicit exclusions have precedence
	// over explicit inclusions.
	ExcludedAlarmDescriptions multiValueStringFlag

	// includedAlarmStatuses is a list of user-specified status keywords for
	// Triggered Alarms that should be explicitly included. This list will be
	// validated and then converted (where needed) into ManagedEntityStatus
	// keywords. See the exported field of the same name for more information.
	includedAlarmStatuses multiValueStringFlag

	// excludedAlarmStatuses is a list of user-specified status keywords for
	// Triggered Alarms that should be explicitly excluded. This list will be
	// validated and then converted (where needed) into ManagedEntityStatus
	// keywords. See the exported field of the same name for more information.
	excludedAlarmStatuses multiValueStringFlag

	// IncludedAlarmNames is a list of statuses for Triggered Alarms that will
	// be explicitly included for evaluation. Unless included by later
	// filtering logic, unmatched Triggered Alarms will be excluded from final
	// evaluation. Explicitly included Triggered Alarms are still subject to
	// permanent exclusion if an explicit exclusion match is made.
	IncludedAlarmStatuses multiValueStringFlag

	// ExcludedAlarmStatuses is a list of statuses for Triggered Alarms that
	// will be explicitly excluded from further evaluation by other stages in
	// the filtering pipeline. Explicit exclusions have precedence over
	// explicit inclusions.
	ExcludedAlarmStatuses multiValueStringFlag

	// App represents common details about the plugins provided by this
	// project.
	App AppInfo

	// Log is an embedded zerolog Logger initialized via config.New().
	Log zerolog.Logger

	// HostSystemMemoryUseWarning specifies the percentage of memory use (as a
	// whole number) for the specified ESXi host when a WARNING threshold is
	// reached.
	HostSystemMemoryUseWarning int

	// HostSystemMemoryUseCritical specifies the percentage of memory use (as
	// a whole number) for the specified ESXi host when a CRITICAL threshold
	// is reached.
	HostSystemMemoryUseCritical int

	// HostSystemCPUUseWarning specifies the percentage of CPU use (as a whole
	// number) for the specified ESXi host when a WARNING threshold is
	// reached.
	HostSystemCPUUseWarning int

	// HostSystemCPUUseCritical specifies the percentage of CPU use (as a
	// whole number) for the specified ESXi host when a CRITICAL threshold is
	// reached.
	HostSystemCPUUseCritical int

	// Port is the TCP port used by the certifcate-enabled service.
	Port int

	// timeout is the value in seconds allowed before a plugin execution
	// attempt is abandoned and an error returned.
	timeout int

	// VCPUsAllocatedWarning specifies the percentage of vCPUs allocation (as
	// a whole number) when a WARNING threshold is reached.
	VCPUsAllocatedWarning int

	// VCPUsAllocatedCritical specifies the percentage of vCPUs allocation (as
	// a whole number) when a CRITICAL threshold is reached.
	VCPUsAllocatedCritical int

	// VCPUsMaxAllowed specifies the maximum amount of virtual CPUs (as a
	// whole number) that we are allowed to allocate in the target VMware
	// environment.
	VCPUsMaxAllowed int

	// ResourcePoolsMemoryUseWarning specifies the percentage of memory use
	// (as a whole number) across all specified Resource Pools when a WARNING
	// threshold is reached.
	ResourcePoolsMemoryUseWarning int

	// ResourcePoolsMemoryUseCritical specifies the percentage of memory use
	// (as a whole number) across all specified Resource Pools when a CRITICAL
	// threshold is reached.
	ResourcePoolsMemoryUseCritical int

	// ResourcePoolsMemoryMaxAllowed specifies the maximum amount of memory
	// that we are allowed to consume in GB (as a whole number) in the target
	// VMware environment across all specified Resource Pools. VMs that are
	// running outside of resource pools are not considered in these
	// calculations.
	ResourcePoolsMemoryMaxAllowed int

	// DatastoreSpaceUsageWarning specifies the percentage of a datastore's
	// storage usage (as a whole number) when a WARNING threshold is reached.
	DatastoreSpaceUsageWarning int

	// DatastoreSpaceUsageCritical specifies the percentage of a datastore's
	// storage usage (as a whole number) when a CRITICAL threshold is reached.
	DatastoreSpaceUsageCritical int

	// datastoreReadLatencyWarning specifies the read latency of a datastore's
	// storage (in ms) when a WARNING threshold is reached.
	datastoreReadLatencyWarning dsPerfLatencyMetricFlag

	// datastoreReadLatencyWarning specifies the read latency of a datastore's
	// storage (in ms) when a CRITICAL threshold is reached.
	datastoreReadLatencyCritical dsPerfLatencyMetricFlag

	// datastoreWriteLatencyWarning specifies the write latency of a
	// datastore's storage (in ms) when a WARNING threshold is reached.
	datastoreWriteLatencyWarning dsPerfLatencyMetricFlag

	// datastoreWriteLatencyCritical specifies the write latency of a
	// datastore's storage (in ms) when a CRITICAL threshold is reached.
	datastoreWriteLatencyCritical dsPerfLatencyMetricFlag

	// datastoreVMLatencyWarning specifies the latency of a datastore's
	// storage (in ms) as observed by VMs using the datastore when a WARNING
	// threshold is reached.
	datastoreVMLatencyWarning dsPerfLatencyMetricFlag

	// datastoreVMLatencyWarning specifies the latency of a datastore's
	// storage (in ms) as observed by VMs using the datastore when a CRITICAL
	// threshold is reached.
	datastoreVMLatencyCritical dsPerfLatencyMetricFlag

	// datastorePerformancePercentileSet specifies the set of
	// DatastorePerformanceSummary latency thresholds associated with a
	// specific percentile.
	datastorePerformancePercentileSet MultiValueDSPerfPercentileSetFlag

	// SnapshotsSizeCritical specifies the cumulative size in GB of all
	// snapshots for a VM when a WARNING threshold is reached.
	SnapshotsSizeWarning int

	// SnapshotsSizeCritical specifies the cumulative size in GB of all
	// snapshots for a VM when a CRITICAL threshold is reached.
	SnapshotsSizeCritical int

	// SnapshotsAgeWarning specifies the age of a snapshot in days when a
	// WARNING threshold is reached.
	SnapshotsAgeWarning int

	// SnapshotsAgeCritical specifies the age of a snapshot in days when a
	// CRITICAL threshold is reached.
	SnapshotsAgeCritical int

	// SnapshotsCountWarning specifies the number of snapshots per VM when a
	// WARNING threshold is reached.
	SnapshotsCountWarning int

	// SnapshotsCountCritical specifies the number of snapshots per VM when a
	// CRITICAL threshold is reached.
	SnapshotsCountCritical int

	// VMPowerCycleUptimeWarning specifies the power cycle (off/on) uptime in
	// days per VM when a WARNING threshold is reached.
	VMPowerCycleUptimeWarning int

	// VMPowerCycleUptimeCritical specifies the power cycle (off/on) uptime in
	// days per VM when a CRITICAL threshold is reached.
	VMPowerCycleUptimeCritical int

	// VirtualHardwareMinimumVersion is the minimum virtual hardware version
	// accepted for each Virtual Machine. Any Virtual Machine not meeting this
	// minimum value is considered to be in a CRITICAL state. Per KB 1003746,
	// version 3 appears to be the oldest version supported.
	VirtualHardwareMinimumVersion int

	// VirtualHardwareOutdatedByWarning specifies the WARNING threshold for
	// outdated virtual hardware versions. If the current virtual hardware
	// version for a VM is found to be more than this many versions older than
	// the latest version a WARNING state is triggered.
	VirtualHardwareOutdatedByWarning int

	// VirtualHardwareOutdatedByCritical specifies the CRITICAL threshold for
	// outdated virtual hardware versions. If the current virtual hardware
	// version for a VM is found to be more than this many versions older than
	// the latest version a CRITICAL state is triggered.
	VirtualHardwareOutdatedByCritical int

	// VirtualHardwareDefaultVersionIsMinimum indicates whether the host or
	// cluster default hardware version is the minimum allowed.
	VirtualHardwareDefaultVersionIsMinimum bool

	// IgnoreMissingCustomAttribute indicates whether a host or datastore
	// missing the specified Custom Attribute should be ignored.
	IgnoreMissingCustomAttribute bool

	// IgnoreMissingDatastorePerfMetrics indicates whether the lack of
	// available metrics for a specific datastore should be ignored. This is
	// not intended to handle scenarios where metrics collection is disabled
	// entirely, but for new datastores where metrics have not yet been
	// collected for the active interval.
	IgnoreMissingDatastorePerfMetrics bool

	// HideHistoricalDatastorePerfMetricSets indicates whether metrics for a
	// specific datastore should be excluded from the performance summary
	// report emitted at plugin completion.
	HideHistoricalDatastorePerfMetricSets bool

	// PoweredOff indicates whether powered off VMs are evaluated in addition
	// to powered on VMs.
	PoweredOff bool

	// EvaluateAcknowledgedAlarms indicates whether acknowledged triggered
	// alarms are evaluated in addition to unacknowledged ones.
	EvaluateAcknowledgedAlarms bool

	// TriggerReloadStateData indicates whether the state data for evaluated
	// objects (e.g., VirtualMachines) will be reloaded/refreshed prior to
	// evaluation of specific properties.
	TriggerReloadStateData bool

	// Whether the certificate should be trusted as-is without validation.
	TrustCert bool

	// EmitBranding controls whether "generated by" text is included at the
	// bottom of application output. This output is included in the Nagios
	// dashboard and notifications. This output may not mix well with branding
	// output from other tools such as atc0005/send2teams which also insert
	// their own branding output.
	EmitBranding bool

	// ShowVersion is a flag indicating whether the user opted to display only
	// the version string and then immediately exit the application.
	ShowVersion bool
}

// Usage is a custom override for the default Help text provided by the flag
// package. Here we prepend some additional metadata to the existing output.
var Usage = func() {

	// Override default of stderr as destination for help output. This allows
	// Nagios XI and similar monitoring systems to call plugins with the
	// `--help` flag and have it display within the Admin web UI.
	flag.CommandLine.SetOutput(os.Stdout)

	fmt.Fprintln(flag.CommandLine.Output(), "\n"+Version()+"\n")
	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}

// Version emits application name, version and repo location.
func Version() string {
	return fmt.Sprintf("%s %s (%s)", myAppName, version, myAppURL)
}

// Branding accepts a message and returns a function that concatenates that
// message with version information. This function is intended to be called as
// a final step before application exit after any other output has already
// been emitted.
func Branding(msg string) func() string {
	return func() string {
		return strings.Join([]string{msg, Version()}, "")
	}
}

// pluginTypeLabel is used as a lookup to return the plugin type label
// associated with the active/specified PluginType.
func pluginTypeLabel(pluginType PluginType) string {

	var label string

	switch {
	case pluginType.SnapshotsAge:
		label = PluginTypeSnapshotsAge

	case pluginType.SnapshotsCount:
		label = PluginTypeSnapshotsCount

	case pluginType.SnapshotsSize:
		label = PluginTypeSnapshotsSize

	case pluginType.DatastoresSpace:
		label = PluginTypeDatastoresSpace

	case pluginType.DatastoresPerformance:
		label = PluginTypeDatastoresPerformance

	case pluginType.ResourcePoolsMemory:
		label = PluginTypeResourcePoolsMemory

	case pluginType.VirtualCPUsAllocation:
		label = PluginTypeVirtualCPUsAllocation

	case pluginType.VirtualHardwareVersion:
		label = PluginTypeVirtualHardwareVersion

	case pluginType.Host2Datastores2VMs:
		label = PluginTypeHostDatastoreVMsPairings

	case pluginType.HostSystemMemory:
		label = PluginTypeHostSystemMemory

	case pluginType.HostSystemCPU:
		label = PluginTypeHostSystemCPU

	case pluginType.VirtualMachinePowerCycleUptime:
		label = PluginTypeVirtualMachinePowerCycleUptime

	case pluginType.DiskConsolidation:
		label = PluginTypeDiskConsolidation

	case pluginType.InteractiveQuestion:
		label = PluginTypeInteractiveQuestion

	case pluginType.Alarms:
		label = PluginTypeAlarms

	case pluginType.Tools:
		label = PluginTypeTools

	default:
		label = "ERROR: Please report this; I evidently forgot to expand the PluginType collection"

	}

	return label

}

// New is a factory function that produces a new Config object based on user
// provided flag and config file values. It is responsible for validating
// user-provided values and initializing the logging settings used by this
// application.
func New(pluginType PluginType) (*Config, error) {
	var config Config

	// Ensure we're working with an initialized map
	if config.datastorePerformancePercentileSet == nil {
		config.datastorePerformancePercentileSet = make(MultiValueDSPerfPercentileSetFlag)
	}

	config.handleFlagsConfig(pluginType)

	config.App = AppInfo{
		Name:    myAppName,
		Version: version,
		URL:     myAppURL,
		Plugin:  pluginTypeLabel(pluginType),
	}

	if config.ShowVersion {
		return nil, ErrVersionRequested
	}

	if err := config.validate(pluginType); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	// initialize logging just as soon as validation is complete
	if err := config.setupLogging(pluginType); err != nil {
		return nil, fmt.Errorf(
			"failed to set logging configuration: %w",
			err,
		)
	}

	// initialize exported TriggeredAlarm status inclusion and exclusion lists
	// based on user-provided keywords after validation is complete
	if err := config.setAlarmStatuses(); err != nil {
		return nil, fmt.Errorf(
			"failed to evaluate provided triggered alarm status keywords: %w",
			err,
		)
	}

	return &config, nil

}
