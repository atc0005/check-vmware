// Copyright 2021 Adam Chalkley
//
// https://github.com/atc0005/check-vmware
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package config

import (
	"fmt"
	"strings"
)

// validate verifies all Config struct fields have been provided acceptable
// values.
func (c Config) validate(pluginType PluginType) error {

	// Flags specific to one plugin type or the other
	switch {
	case pluginType.Tools:

		// only one of these options may be used
		if len(c.ExcludedResourcePools) > 0 && len(c.IncludedResourcePools) > 0 {
			return fmt.Errorf(
				"only one of %q or %q flags may be specified",
				"include-rp",
				"exclude-rp",
			)
		}

	case pluginType.SnapshotsAge:

		// only one of these options may be used
		if len(c.ExcludedResourcePools) > 0 && len(c.IncludedResourcePools) > 0 {
			return fmt.Errorf(
				"only one of %q or %q flags may be specified",
				"include-rp",
				"exclude-rp",
			)
		}

		if c.SnapshotsAgeWarning < 0 {
			return fmt.Errorf(
				"invalid snapshot age WARNING threshold number: %d",
				c.SnapshotsAgeWarning,
			)
		}

		if c.SnapshotsAgeCritical < 0 {
			return fmt.Errorf(
				"invalid snapshot age CRITICAL threshold number: %d",
				c.SnapshotsAgeCritical,
			)
		}

		if c.SnapshotsAgeCritical <= c.SnapshotsAgeWarning {
			return fmt.Errorf(
				"critical threshold set lower than or equal to warning threshold",
			)
		}

	case pluginType.SnapshotsCount:

		// only one of these options may be used
		if len(c.ExcludedResourcePools) > 0 && len(c.IncludedResourcePools) > 0 {
			return fmt.Errorf(
				"only one of %q or %q flags may be specified",
				"include-rp",
				"exclude-rp",
			)
		}

		if c.SnapshotsCountWarning < 0 {
			return fmt.Errorf(
				"invalid snapshot count WARNING threshold number: %d",
				c.SnapshotsCountWarning,
			)
		}

		if c.SnapshotsCountCritical < 0 {
			return fmt.Errorf(
				"invalid snapshot count CRITICAL threshold number: %d",
				c.SnapshotsCountCritical,
			)
		}

		if c.SnapshotsCountCritical <= c.SnapshotsCountWarning {
			return fmt.Errorf(
				"critical threshold set lower than or equal to warning threshold",
			)
		}

	case pluginType.SnapshotsSize:

		// only one of these options may be used
		if len(c.ExcludedResourcePools) > 0 && len(c.IncludedResourcePools) > 0 {
			return fmt.Errorf(
				"only one of %q or %q flags may be specified",
				"include-rp",
				"exclude-rp",
			)
		}

		if c.SnapshotsSizeWarning < 0 {
			return fmt.Errorf(
				"invalid snapshot size WARNING threshold number: %d",
				c.SnapshotsSizeWarning,
			)
		}

		if c.SnapshotsSizeCritical < 0 {
			return fmt.Errorf(
				"invalid snapshot size CRITICAL threshold number: %d",
				c.SnapshotsSizeCritical,
			)
		}

		if c.SnapshotsSizeCritical <= c.SnapshotsSizeWarning {
			return fmt.Errorf(
				"critical threshold set lower than or equal to warning threshold",
			)
		}

	case pluginType.VirtualMachinePowerCycleUptime:

		// only one of these options may be used
		if len(c.ExcludedResourcePools) > 0 && len(c.IncludedResourcePools) > 0 {
			return fmt.Errorf(
				"only one of %q or %q flags may be specified",
				"include-rp",
				"exclude-rp",
			)
		}

		if c.VMPowerCycleUptimeWarning < 0 {
			return fmt.Errorf(
				"invalid VM power cycle uptime WARNING threshold number: %d",
				c.VMPowerCycleUptimeWarning,
			)
		}

		if c.VMPowerCycleUptimeCritical < 0 {
			return fmt.Errorf(
				"invalid VM power cycle uptime CRITICAL threshold number: %d",
				c.VMPowerCycleUptimeCritical,
			)
		}

		if c.VMPowerCycleUptimeCritical <= c.VMPowerCycleUptimeWarning {
			return fmt.Errorf(
				"critical threshold set lower than or equal to warning threshold",
			)
		}

	case pluginType.DatastoresSize:

		if c.DatastoreName == "" {
			return fmt.Errorf("datastore name not provided")
		}

		if c.DatastoreUsageCritical < 1 {
			return fmt.Errorf(
				"invalid datastore usage (percentage as whole number) CRITICAL threshold number: %d",
				c.DatastoreUsageCritical,
			)
		}

		if c.DatastoreUsageWarning < 1 {
			return fmt.Errorf(
				"invalid datastore usage (percentage as whole number) WARNING threshold number: %d",
				c.DatastoreUsageWarning,
			)
		}

		if c.DatastoreUsageCritical <= c.DatastoreUsageWarning {
			return fmt.Errorf(
				"datastore critical threshold set lower than or equal to warning threshold",
			)
		}

	case pluginType.DatastoresPerformance:

		if c.DatastoreName == "" {
			return fmt.Errorf("datastore name not provided")
		}

		switch {

		// Datastore Performance percentile set flags were not provided.
		// Evaluate individual latency threshold flag values.
		case len(c.datastorePerformancePercentileSet) == 0:

			latencyPerfThresholds := c.DatastorePerfThresholds()

			tt := []struct {
				latencyCritical    float64
				latencyWarning     float64
				latencyDescription string
			}{
				{
					latencyCritical:    latencyPerfThresholds.ReadLatencyCritical,
					latencyWarning:     latencyPerfThresholds.ReadLatencyWarning,
					latencyDescription: "read",
				},
				{
					latencyCritical:    latencyPerfThresholds.WriteLatencyCritical,
					latencyWarning:     latencyPerfThresholds.WriteLatencyWarning,
					latencyDescription: "write",
				},
				{
					latencyCritical:    latencyPerfThresholds.VMLatencyCritical,
					latencyWarning:     latencyPerfThresholds.VMLatencyWarning,
					latencyDescription: "VM",
				},
			}

			for _, threshold := range tt {

				if threshold.latencyCritical < 1 {
					return fmt.Errorf(
						"invalid datastore latency (percentage as whole number) CRITICAL threshold number: %f",
						threshold.latencyCritical,
					)
				}

				if threshold.latencyWarning < 1 {
					return fmt.Errorf(
						"invalid datastore latency (percentage as whole number) WARNING threshold number: %f",
						threshold.latencyWarning,
					)
				}

				if threshold.latencyCritical <= threshold.latencyWarning {
					return fmt.Errorf(
						"datastore latency critical threshold set lower than or equal to warning threshold",
					)
				}

			}

		// Datastore performance percentile set was specified. Individual
		// latency flags are not permitted.
		case len(c.datastorePerformancePercentileSet) > 0:

			latencyThresholdFlags := []dsPerfLatencyMetricFlag{
				c.datastoreReadLatencyWarning,
				c.datastoreReadLatencyCritical,
				c.datastoreWriteLatencyWarning,
				c.datastoreWriteLatencyCritical,
				c.datastoreVMLatencyWarning,
				c.datastoreVMLatencyCritical,
			}

			for i := range latencyThresholdFlags {
				if latencyThresholdFlags[i].isSet {
					return fmt.Errorf(
						"invalid combination of flags; percentile set flag is incompatible with individual latency threshold flags",
					)
				}
			}

			// If percentile sets were provided, evaluate each percentile from
			// the set against known supported percentiles.
			supportedPercentiles := getSupportedDatastorePerfPercentiles()
			isSupportedPercentile := func(specified int, supported []int) bool {
				for i := range supported {
					if specified == supported[i] {
						return true
					}
				}

				return false
			}

			for specifiedPercentile := range c.datastorePerformancePercentileSet {
				if !isSupportedPercentile(specifiedPercentile, supportedPercentiles) {
					return fmt.Errorf(
						"invalid percentile specified; got percentile %v, expected one of %v",
						specifiedPercentile,
						supportedPercentiles,
					)
				}
			}

		}

	case pluginType.HostSystemMemory:

		if c.HostSystemName == "" {
			return fmt.Errorf("host name not provided")
		}

		if c.HostSystemMemoryUseCritical < 1 {
			return fmt.Errorf(
				"invalid host memory usage (percentage as whole number) CRITICAL threshold number: %d",
				c.HostSystemMemoryUseCritical,
			)
		}

		if c.HostSystemMemoryUseWarning < 1 {
			return fmt.Errorf(
				"invalid host memory usage (percentage as whole number) WARNING threshold number: %d",
				c.HostSystemMemoryUseWarning,
			)
		}

		if c.HostSystemMemoryUseCritical <= c.HostSystemMemoryUseWarning {
			return fmt.Errorf(
				"critical threshold set lower than or equal to warning threshold",
			)
		}

	case pluginType.HostSystemCPU:

		if c.HostSystemName == "" {
			return fmt.Errorf("host name not provided")
		}

		if c.HostSystemCPUUseCritical < 1 {
			return fmt.Errorf(
				"invalid host CPU usage (percentage as whole number) CRITICAL threshold number: %d",
				c.HostSystemCPUUseCritical,
			)
		}

		if c.HostSystemCPUUseWarning < 1 {
			return fmt.Errorf(
				"invalid host CPU usage (percentage as whole number) WARNING threshold number: %d",
				c.HostSystemCPUUseWarning,
			)
		}

		if c.HostSystemCPUUseCritical <= c.HostSystemCPUUseWarning {
			return fmt.Errorf(
				"critical threshold set lower than or equal to warning threshold",
			)
		}

	case pluginType.ResourcePoolsMemory:

		// only one of these options may be used
		if len(c.ExcludedResourcePools) > 0 && len(c.IncludedResourcePools) > 0 {
			return fmt.Errorf(
				"only one of %q or %q flags may be specified",
				"include-rp",
				"exclude-rp",
			)
		}

		if c.ResourcePoolsMemoryMaxAllowed < 1 {
			return fmt.Errorf(
				"invalid value specified for maximum memory usage allowed: %d",
				c.ResourcePoolsMemoryMaxAllowed,
			)
		}

		if c.ResourcePoolsMemoryUseCritical < 1 {
			return fmt.Errorf(
				"invalid memory usage CRITICAL threshold number: %d",
				c.ResourcePoolsMemoryUseCritical,
			)
		}

		if c.ResourcePoolsMemoryUseWarning < 1 {
			return fmt.Errorf(
				"invalid memory usage WARNING threshold number: %d",
				c.ResourcePoolsMemoryUseWarning,
			)
		}

		if c.ResourcePoolsMemoryUseCritical <= c.ResourcePoolsMemoryUseWarning {
			return fmt.Errorf(
				"memory usage critical threshold set lower than or equal to warning threshold",
			)
		}

	case pluginType.VirtualCPUsAllocation:

		// only one of these options may be used
		if len(c.ExcludedResourcePools) > 0 && len(c.IncludedResourcePools) > 0 {
			return fmt.Errorf(
				"only one of %q or %q flags may be specified",
				"include-rp",
				"exclude-rp",
			)
		}

		if c.VCPUsMaxAllowed < 1 {
			return fmt.Errorf(
				"invalid value specified for maximum number of vCPUs allowed: %d",
				c.VCPUsMaxAllowed,
			)
		}

		if c.VCPUsAllocatedCritical < 1 {
			return fmt.Errorf(
				"invalid vCPUs allocation CRITICAL threshold number: %d",
				c.VCPUsAllocatedCritical,
			)
		}

		if c.VCPUsAllocatedWarning < 1 {
			return fmt.Errorf(
				"invalid vCPUs allocation WARNING threshold number: %d",
				c.VCPUsAllocatedWarning,
			)
		}

		if c.VCPUsAllocatedCritical <= c.VCPUsAllocatedWarning {
			return fmt.Errorf(
				"vCPUs allocation critical threshold set lower than or equal to warning threshold",
			)
		}

	case pluginType.Host2Datastores2VMs:

		// only one of these options may be used
		if len(c.ExcludedResourcePools) > 0 && len(c.IncludedResourcePools) > 0 {
			return fmt.Errorf(
				"only one of %q or %q flags may be specified",
				"include-rp",
				"exclude-rp",
			)
		}

		// Validate that *only one* of shared Custom Attribute name is
		// provided or both datastore and host Custom Attribute names are
		// provided.
		switch {

		// no Custom Attribute provided
		case c.sharedCustomAttributeName == "" &&
			(c.datastoreCustomAttributeName == "" && c.hostCustomAttributeName == ""):

			return fmt.Errorf(
				"one of shared or resource-specific Custom Attribute name must be specified",
			)

		// shared Custom Attribute and one of resource-specific Custom
		// Attribute provided
		case c.sharedCustomAttributeName != "" &&
			(c.datastoreCustomAttributeName != "" || c.hostCustomAttributeName != ""):

			return fmt.Errorf(
				"only one of shared or resource-specific Custom Attribute name may be specified",
			)

		// shared Custom Attribute not provided and either of datastore or
		// host Custom Attribute not provided
		case c.sharedCustomAttributeName == "" &&
			c.datastoreCustomAttributeName == "" && c.hostCustomAttributeName != "":

			return fmt.Errorf(
				"datastore Custom Attribute name must be specified if providing Custom Attribute name for hosts",
			)

		case c.sharedCustomAttributeName == "" &&
			c.datastoreCustomAttributeName != "" && c.hostCustomAttributeName == "":

			return fmt.Errorf(
				"host Custom Attribute name must be specified if providing Custom Attribute name for datastores",
			)

		}

		// Validate that shared Custom Attribute separator is provided, both
		// datastore and host Custom Attribute separators are provided (and
		// not shared), or no Custom Attribute separator is provided.
		switch {

		// no Custom Attribute prefix separator provided
		case c.sharedCustomAttributePrefixSeparator == "" &&
			(c.datastoreCustomAttributePrefixSeparator == "" && c.hostCustomAttributePrefixSeparator == ""):

			// this is a valid scenario and indicates that literal Custom
			// Attribute value matching is performed.

		// shared Custom Attribute prefix separator and one of
		// resource-specific Custom Attribute prefix separators provided
		case c.sharedCustomAttributePrefixSeparator != "" &&
			(c.datastoreCustomAttributePrefixSeparator != "" || c.hostCustomAttributePrefixSeparator != ""):

			return fmt.Errorf(
				"error: Custom Attribute prefix separators may only be specified as a shared value, or for both datastore and hosts",
			)

		case c.sharedCustomAttributePrefixSeparator == "" &&
			c.datastoreCustomAttributePrefixSeparator == "" && c.hostCustomAttributePrefixSeparator != "":

			return fmt.Errorf(
				"datastore Custom Attribute prefix must be specified if providing prefix for hosts",
			)

		case c.sharedCustomAttributePrefixSeparator == "" &&
			c.datastoreCustomAttributePrefixSeparator != "" && c.hostCustomAttributePrefixSeparator == "":

			return fmt.Errorf(
				"host Custom Attribute prefix must be specified if providing prefix for datastores",
			)

		}

	case pluginType.VirtualHardwareVersion:

		// only one of these options may be used
		if len(c.ExcludedResourcePools) > 0 && len(c.IncludedResourcePools) > 0 {
			return fmt.Errorf(
				"only one of %q or %q flags may be specified",
				"include-rp",
				"exclude-rp",
			)
		}

		// optional flag; if not default value, assert known requirements
		if c.ClusterName != defaultClusterName {
			if len(c.ClusterName) > MaxClusterNameChars {
				return fmt.Errorf(
					"invalid cluster name specified; max supported length is %d, received %d",
					MaxClusterNameChars,
					len(c.ClusterName),
				)
			}
		}

		// both are optional flags, but only one at a time is supported
		if c.ClusterName != defaultClusterName && c.HostSystemName != defaultHostSystemName {
			return fmt.Errorf(
				"only one of cluster or host name supported",
			)
		}

		// assert that only one type of behavior is used for plugin
		switch {

		// homogeneous version checks
		case c.VirtualHardwareMinimumVersion == defaultVirtualHardwareMinimumVersion &&
			c.VirtualHardwareOutdatedByCritical == defaultVirtualHardwareOutdatedByCritical &&
			c.VirtualHardwareOutdatedByWarning == defaultVirtualHardwareOutdatedByWarning &&
			!c.VirtualHardwareDefaultVersionIsMinimum:

		// host/cluster default is minimum version check
		case c.VirtualHardwareMinimumVersion == defaultVirtualHardwareMinimumVersion &&
			c.VirtualHardwareOutdatedByCritical == defaultVirtualHardwareOutdatedByCritical &&
			c.VirtualHardwareOutdatedByWarning == defaultVirtualHardwareOutdatedByWarning &&
			c.VirtualHardwareDefaultVersionIsMinimum:

		// minimum version check
		case c.VirtualHardwareMinimumVersion != defaultVirtualHardwareMinimumVersion &&
			c.VirtualHardwareOutdatedByCritical == defaultVirtualHardwareOutdatedByCritical &&
			c.VirtualHardwareOutdatedByWarning == defaultVirtualHardwareOutdatedByWarning &&
			!c.VirtualHardwareDefaultVersionIsMinimum:

			// ESX 2.x, GSX Server 3.x, Workstation 4.x & 5.x, ...
			// https://kb.vmware.com/s/article/1003746
			if c.VirtualHardwareMinimumVersion < 3 {
				return fmt.Errorf("invalid value specified for minimum virtual hardware version")
			}

		// outdated by version thresholds check; apply further validation
		case c.VirtualHardwareMinimumVersion == defaultVirtualHardwareMinimumVersion &&
			(c.VirtualHardwareOutdatedByCritical != defaultVirtualHardwareOutdatedByCritical ||
				c.VirtualHardwareOutdatedByWarning != defaultVirtualHardwareOutdatedByWarning) &&
			!c.VirtualHardwareDefaultVersionIsMinimum:

			switch {
			// user did not specify a value, do not apply further validation
			// checks for this field
			case c.VirtualHardwareOutdatedByCritical == defaultVirtualHardwareOutdatedByCritical:

			case c.VirtualHardwareOutdatedByCritical < 1:
				return fmt.Errorf("invalid value specified for outdated by critical threshold")
			}

			switch {
			// user did not specify a value, do not apply further validation
			// checks for this field
			case c.VirtualHardwareOutdatedByWarning == defaultVirtualHardwareOutdatedByWarning:

			case c.VirtualHardwareOutdatedByWarning < 1:
				return fmt.Errorf("invalid value specified for outdated by warning threshold")
			}

			switch {
			case c.VirtualHardwareOutdatedByWarning == defaultVirtualHardwareOutdatedByWarning &&
				c.VirtualHardwareOutdatedByCritical != defaultVirtualHardwareOutdatedByCritical:

				return fmt.Errorf(
					"outdated by critical threshold specified, but not warning threshold; both critical and warning thresholds must be set if using outdated-by plugin mode",
				)

			case c.VirtualHardwareOutdatedByWarning != defaultVirtualHardwareOutdatedByWarning &&
				c.VirtualHardwareOutdatedByCritical == defaultVirtualHardwareOutdatedByCritical:

				return fmt.Errorf(
					"outdated by warning threshold specified, but not critical threshold; both critical and warning thresholds must be set if using outdated-by plugin mode",
				)
			}

			if c.VirtualHardwareOutdatedByCritical <= c.VirtualHardwareOutdatedByWarning {
				return fmt.Errorf(
					"outdated by critical threshold set lower than or equal to warning threshold",
				)
			}

		default:

			return fmt.Errorf("unsupported plugin mode requested")

		}

	case pluginType.Alarms:

		// only one of these options may be used
		if len(c.IncludedAlarmEntityTypes) > 0 && len(c.ExcludedAlarmEntityTypes) > 0 {
			return fmt.Errorf(
				"only one of %q or %q flags may be specified",
				"include-entity-type",
				"exclude-entity-type",
			)
		}

		// only one of these options may be used
		if len(c.IncludedAlarmEntityNames) > 0 && len(c.ExcludedAlarmEntityNames) > 0 {
			return fmt.Errorf(
				"only one of %q or %q flags may be specified",
				"include-entity-name",
				"exclude-entity-name",
			)
		}

		// only one of these options may be used
		if len(c.ExcludedAlarmEntityResourcePools) > 0 && len(c.IncludedAlarmEntityResourcePools) > 0 {
			return fmt.Errorf(
				"only one of %q or %q flags may be specified",
				"include-entity-rp",
				"exclude-entity-rp",
			)
		}

		// only one of these options may be used
		if len(c.IncludedAlarmNames) > 0 && len(c.ExcludedAlarmNames) > 0 {
			return fmt.Errorf(
				"only one of %q or %q flags may be specified",
				"include-name",
				"exclude-name",
			)
		}

		// only one of these options may be used
		if len(c.IncludedAlarmDescriptions) > 0 && len(c.ExcludedAlarmDescriptions) > 0 {
			return fmt.Errorf(
				"only one of %q or %q flags may be specified",
				"include-desc",
				"exclude-desc",
			)
		}

		if len(c.IncludedAlarmStatuses) > 0 || len(c.ExcludedAlarmStatuses) > 0 {

			// only one of these options may be used
			if len(c.IncludedAlarmStatuses) > 0 && len(c.ExcludedAlarmStatuses) > 0 {
				return fmt.Errorf(
					"only one of %q or %q flags may be specified",
					"include-status",
					"exclude-status",
				)
			}

			alarmStatuses := getTriggeredAlarmStatuses()

			// If there are explicit inclusions, validate keyword and then
			// reject "green" or "ok" status for Triggered Alarms since alarms
			// (afaik) don't trigger for that state.
			if len(c.includedAlarmStatuses) > 0 {
				for _, keyword := range c.includedAlarmStatuses {
					requestedkeyword := strings.ToLower(keyword)
					_, ok := alarmStatuses[requestedkeyword]

					if !ok || requestedkeyword == AlarmStatusOk ||
						requestedkeyword == AlarmStatusGreen {
						return fmt.Errorf(
							"invalid triggered alarm status for inclusion: %q",
							keyword,
						)
					}
				}
			}

			// If there are explicit exclusions, validate keyword and then
			// reject "green" or "ok" status for Triggered Alarms since alarms
			// (afaik) don't trigger for that state.
			if len(c.excludedAlarmStatuses) > 0 {
				for _, keyword := range c.excludedAlarmStatuses {
					requestedkeyword := strings.ToLower(keyword)
					_, ok := alarmStatuses[requestedkeyword]

					if !ok || requestedkeyword == AlarmStatusOk ||
						requestedkeyword == AlarmStatusGreen {
						return fmt.Errorf(
							"invalid triggered alarm status for exclusion: %q",
							keyword,
						)
					}
				}
			}

		}
	}

	// shared validation checks

	if c.Server == "" {
		return fmt.Errorf("server FQDN or IP Address not provided")
	}

	if c.Username == "" {
		return fmt.Errorf("username not provided")
	}

	if c.Password == "" {
		return fmt.Errorf("password not provided")
	}

	if c.Port < 0 {
		return fmt.Errorf("invalid TCP port number %d", c.Port)
	}

	if c.Timeout() < 1 {
		return fmt.Errorf("invalid timeout value %d provided", c.Timeout())
	}

	requestedLoggingLevel := strings.ToLower(c.LoggingLevel)
	if _, ok := loggingLevels[requestedLoggingLevel]; !ok {
		return fmt.Errorf("invalid logging level %q", c.LoggingLevel)
	}

	// Optimist
	return nil

}
