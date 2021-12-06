// Copyright 2021 Adam Chalkley
//
// https://github.com/atc0005/check-vmware
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package config

import "flag"

// handleFlagsConfig handles toggling the exposure of specific configuration
// flags to the user. This behavior is controlled via the specified plugin
// type as set by each cmd. Based on the plugin type, a smaller subset of
// flags specific to each type are exposed along with a set common to all
// plugin types.
func (c *Config) handleFlagsConfig(pluginType PluginType) {

	// Flags specific to one plugin type or the other
	switch {
	case pluginType.Tools:

		flag.Var(&c.IncludedResourcePools, "include-rp", vmIncludedResourcePoolsFlagHelp)
		flag.Var(&c.ExcludedResourcePools, "exclude-rp", vmExcludedResourcePoolsFlagHelp)
		flag.Var(&c.IgnoredVMs, "ignore-vm", ignoreVMsFlagHelp)
		flag.BoolVar(&c.PoweredOff, "powered-off", defaultPoweredOff, poweredOffFlagHelp)

	case pluginType.SnapshotsAge:

		flag.Var(&c.IncludedResourcePools, "include-rp", vmIncludedResourcePoolsFlagHelp)
		flag.Var(&c.ExcludedResourcePools, "exclude-rp", vmExcludedResourcePoolsFlagHelp)
		flag.Var(&c.IgnoredVMs, "ignore-vm", ignoreVMsFlagHelp)

		// NOTE: This plugin is hard-coded to evaluate powered off and powered
		// on VMs equally. I'm not sure whether ignoring powered off VMs by
		// default makes sense for this particular plugin.
		//
		// Please share your feedback here if you feel differently:
		// https://github.com/atc0005/check-vmware/discussions/177
		//
		// flag.BoolVar(&c.PoweredOff, "powered-off", defaultPoweredOff, poweredOffFlagHelp)

		flag.IntVar(&c.SnapshotsAgeWarning, "age-warning", defaultSnapshotsAgeWarning, snapshotsAgeWarningFlagHelp)
		flag.IntVar(&c.SnapshotsAgeWarning, "aw", defaultSnapshotsAgeWarning, snapshotsAgeWarningFlagHelp+" (shorthand)")

		flag.IntVar(&c.SnapshotsAgeCritical, "age-critical", defaultSnapshotsAgeCritical, snapshotsAgeCriticalFlagHelp)
		flag.IntVar(&c.SnapshotsAgeCritical, "ac", defaultSnapshotsAgeCritical, snapshotsAgeCriticalFlagHelp+" (shorthand)")

	case pluginType.SnapshotsCount:

		flag.Var(&c.IncludedResourcePools, "include-rp", vmIncludedResourcePoolsFlagHelp)
		flag.Var(&c.ExcludedResourcePools, "exclude-rp", vmExcludedResourcePoolsFlagHelp)
		flag.Var(&c.IgnoredVMs, "ignore-vm", ignoreVMsFlagHelp)

		// NOTE: This plugin is hard-coded to evaluate powered off and powered
		// on VMs equally. I'm not sure whether ignoring powered off VMs by
		// default makes sense for this particular plugin.
		//
		// Please share your feedback here if you feel differently:
		// https://github.com/atc0005/check-vmware/discussions/177
		//
		// flag.BoolVar(&c.PoweredOff, "powered-off", defaultPoweredOff, poweredOffFlagHelp)

		flag.IntVar(&c.SnapshotsCountWarning, "count-warning", defaultSnapshotsCountWarning, snapshotsCountWarningFlagHelp)
		flag.IntVar(&c.SnapshotsCountWarning, "cw", defaultSnapshotsCountWarning, snapshotsCountWarningFlagHelp+" (shorthand)")

		flag.IntVar(&c.SnapshotsCountCritical, "count-critical", defaultSnapshotsCountCritical, snapshotsCountCriticalFlagHelp)
		flag.IntVar(&c.SnapshotsCountCritical, "cc", defaultSnapshotsCountCritical, snapshotsCountCriticalFlagHelp+" (shorthand)")

	case pluginType.SnapshotsSize:

		flag.Var(&c.IncludedResourcePools, "include-rp", vmIncludedResourcePoolsFlagHelp)
		flag.Var(&c.ExcludedResourcePools, "exclude-rp", vmExcludedResourcePoolsFlagHelp)
		flag.Var(&c.IgnoredVMs, "ignore-vm", ignoreVMsFlagHelp)

		// NOTE: This plugin is hard-coded to evaluate powered off and powered
		// on VMs equally. I'm not sure whether ignoring powered off VMs by
		// default makes sense for this particular plugin.
		//
		// Please share your feedback here if you feel differently:
		// https://github.com/atc0005/check-vmware/discussions/177
		//
		// flag.BoolVar(&c.PoweredOff, "powered-off", defaultPoweredOff, poweredOffFlagHelp)

		flag.IntVar(&c.SnapshotsSizeWarning, "size-warning", defaultSnapshotsSizeWarning, snapshotsSizeWarningFlagHelp)
		flag.IntVar(&c.SnapshotsSizeWarning, "sw", defaultSnapshotsSizeWarning, snapshotsSizeWarningFlagHelp+" (shorthand)")

		flag.IntVar(&c.SnapshotsSizeCritical, "size-critical", defaultSnapshotsSizeCritical, snapshotsSizeCriticalFlagHelp)
		flag.IntVar(&c.SnapshotsSizeCritical, "sc", defaultSnapshotsSizeCritical, snapshotsSizeCriticalFlagHelp+" (shorthand)")

	case pluginType.VirtualMachinePowerCycleUptime:

		flag.Var(&c.IncludedResourcePools, "include-rp", vmIncludedResourcePoolsFlagHelp)
		flag.Var(&c.ExcludedResourcePools, "exclude-rp", vmExcludedResourcePoolsFlagHelp)
		flag.Var(&c.IgnoredVMs, "ignore-vm", ignoreVMsFlagHelp)

		flag.IntVar(&c.VMPowerCycleUptimeWarning, "uptime-warning", defaultVMPowerCycleUptimeWarning, vmPowerCycleUptimeWarningFlagHelp)
		flag.IntVar(&c.VMPowerCycleUptimeWarning, "uw", defaultVMPowerCycleUptimeWarning, vmPowerCycleUptimeWarningFlagHelp+" (shorthand)")

		flag.IntVar(&c.VMPowerCycleUptimeCritical, "uptime-critical", defaultVMPowerCycleUptimeCritical, vmPowerCycleUptimeCriticalFlagHelp)
		flag.IntVar(&c.VMPowerCycleUptimeCritical, "uc", defaultVMPowerCycleUptimeCritical, vmPowerCycleUptimeCriticalFlagHelp+" (shorthand)")

	case pluginType.DiskConsolidation:

		flag.Var(&c.IncludedResourcePools, "include-rp", vmIncludedResourcePoolsFlagHelp)
		flag.Var(&c.ExcludedResourcePools, "exclude-rp", vmExcludedResourcePoolsFlagHelp)
		flag.Var(&c.IgnoredVMs, "ignore-vm", ignoreVMsFlagHelp)
		flag.BoolVar(&c.TriggerReloadStateData, "trigger-reload", defaultTriggerReloadStateData, triggerReloadStateDataFlagHelp)

		// NOTE: This plugin is hard-coded to evaluate powered off and powered
		// on VMs equally. I'm not sure whether ignoring powered off VMs by
		// default makes sense for this particular plugin.
		//
		// Please share your feedback here if you feel differently:
		// https://github.com/atc0005/check-vmware/discussions/176
		//
		// Please expand on some use cases for ignoring powered off VMs by default.
		//
		// flag.BoolVar(&c.PoweredOff, "powered-off", defaultPoweredOff, poweredOffFlagHelp)

	case pluginType.InteractiveQuestion:

		flag.Var(&c.IncludedResourcePools, "include-rp", vmIncludedResourcePoolsFlagHelp)
		flag.Var(&c.ExcludedResourcePools, "exclude-rp", vmExcludedResourcePoolsFlagHelp)
		flag.Var(&c.IgnoredVMs, "ignore-vm", ignoreVMsFlagHelp)

	case pluginType.Alarms:

		flag.Var(&c.DatacenterNames, "dc-name", datacenterNamesFlagHelp)
		flag.Var(&c.IncludedAlarmEntityTypes, "include-entity-type", includedAlarmEntityTypesFlagHelp)
		flag.Var(&c.ExcludedAlarmEntityTypes, "exclude-entity-type", excludedAlarmEntityTypesFlagHelp)

		flag.BoolVar(&c.EvaluateAcknowledgedAlarms, "eval-acknowledged", defaultEvaluateAcknowledgedAlarms, evaluateAcknowledgedTriggeredAlarmFlagHelp)

		flag.Var(&c.IncludedAlarmNames, "include-name", includedAlarmNamesFlagHelp)
		flag.Var(&c.ExcludedAlarmNames, "exclude-name", excludedAlarmNamesFlagHelp)

		flag.Var(&c.IncludedAlarmDescriptions, "include-desc", includedAlarmDescriptionsFlagHelp)
		flag.Var(&c.ExcludedAlarmDescriptions, "exclude-desc", excludedAlarmDescriptionsFlagHelp)

		flag.Var(&c.includedAlarmStatuses, "include-status", includedAlarmStatusesFlagHelp)
		flag.Var(&c.excludedAlarmStatuses, "exclude-status", excludedAlarmStatusesFlagHelp)

		flag.Var(&c.IncludedAlarmEntityNames, "include-entity-name", includedAlarmEntityNamesFlagHelp)
		flag.Var(&c.ExcludedAlarmEntityNames, "exclude-entity-name", excludedAlarmEntityNamesFlagHelp)

		flag.Var(&c.IncludedAlarmEntityResourcePools, "include-entity-rp", includedAlarmEntityResourcePoolsFlagHelp)
		flag.Var(&c.ExcludedAlarmEntityResourcePools, "exclude-entity-rp", excludedAlarmEntityResourcePoolsFlagHelp)

	case pluginType.DatastoresSize:

		flag.StringVar(&c.DatacenterName, "dc-name", defaultDatacenterName, datacenterNameFlagHelp)

		flag.StringVar(&c.DatastoreName, "ds-name", defaultDatastoreName, datastoreNameFlagHelp)

		flag.IntVar(&c.DatastoreUsageWarning, "ds-usage-warning", defaultDatastoreUsageWarning, datastoreUsageWarningFlagHelp)
		flag.IntVar(&c.DatastoreUsageWarning, "dsuw", defaultDatastoreUsageWarning, datastoreUsageWarningFlagHelp+" (shorthand)")

		flag.IntVar(&c.DatastoreUsageCritical, "ds-usage-critical", defaultDatastoreUsageCritical, datastoreUsageCriticalFlagHelp)
		flag.IntVar(&c.DatastoreUsageCritical, "dsuc", defaultDatastoreUsageCritical, datastoreUsageCriticalFlagHelp+" (shorthand)")

	case pluginType.DatastoresPerformance:

		flag.StringVar(&c.DatacenterName, "dc-name", defaultDatacenterName, datacenterNameFlagHelp)

		flag.StringVar(&c.DatastoreName, "ds-name", defaultDatastoreName, datastoreNameFlagHelp)

		flag.BoolVar(&c.IgnoreMissingDatastorePerfMetrics, "ds-ignore-missing-metrics", defaultIgnoreMissingDatastoreMetrics, ignoreMissingDatastorePerfMetricsFlagHelp)
		flag.BoolVar(&c.IgnoreMissingDatastorePerfMetrics, "dsim", defaultIgnoreMissingDatastoreMetrics, ignoreMissingDatastorePerfMetricsFlagHelp+" (shorthand)")

		flag.BoolVar(&c.HideHistoricalDatastorePerfMetricSets, "ds-hide-historical-metric-sets", defaultHideHistoricalDatastorePerfMetricSets, hideHistoricalDatastorePerfMetricSetsFlagHelp)
		flag.BoolVar(&c.HideHistoricalDatastorePerfMetricSets, "dshhms", defaultHideHistoricalDatastorePerfMetricSets, hideHistoricalDatastorePerfMetricSetsFlagHelp+" (shorthand)")

		flag.Var(c.datastoreReadLatencyWarning, "ds-read-latency-warning", datastoreReadLatencyWarningFlagHelp)
		flag.Var(c.datastoreReadLatencyWarning, "dsrlw", datastoreReadLatencyWarningFlagHelp+" (shorthand)")

		flag.Var(c.datastoreReadLatencyCritical, "ds-read-latency-critical", datastoreReadLatencyCriticalFlagHelp)
		flag.Var(c.datastoreReadLatencyCritical, "dsrlc", datastoreReadLatencyCriticalFlagHelp+" (shorthand)")

		flag.Var(c.datastoreWriteLatencyWarning, "ds-write-latency-warning", datastoreWriteLatencyWarningFlagHelp)
		flag.Var(c.datastoreWriteLatencyWarning, "dswlw", datastoreWriteLatencyWarningFlagHelp+" (shorthand)")

		flag.Var(c.datastoreWriteLatencyCritical, "ds-write-latency-critical", datastoreWriteLatencyCriticalFlagHelp)
		flag.Var(c.datastoreWriteLatencyCritical, "dswlc", datastoreWriteLatencyCriticalFlagHelp+" (shorthand)")

		flag.Var(c.datastoreVMLatencyWarning, "ds-vm-latency-warning", datastoreVMLatencyWarningFlagHelp)
		flag.Var(c.datastoreVMLatencyWarning, "dsvmlw", datastoreVMLatencyWarningFlagHelp+" (shorthand)")

		flag.Var(c.datastoreVMLatencyCritical, "ds-vm-latency-critical", datastoreVMLatencyCriticalFlagHelp)
		flag.Var(c.datastoreVMLatencyCritical, "dsvmlc", datastoreVMLatencyCriticalFlagHelp+" (shorthand)")

		flag.Var(&c.datastorePerformancePercentileSet, "ds-latency-percentile-set", datastoreLatencyPercintileSetFlagHelp)
		flag.Var(&c.datastorePerformancePercentileSet, "dslps", datastoreLatencyPercintileSetFlagHelp+" (shorthand)")

	case pluginType.HostSystemMemory:

		flag.StringVar(&c.DatacenterName, "dc-name", defaultDatacenterName, datacenterNameFlagHelp)

		flag.StringVar(&c.HostSystemName, "host-name", defaultHostSystemName, hostSystemNameFlagHelp)

		flag.IntVar(&c.HostSystemMemoryUseWarning, "memory-usage-warning", defaultMemoryUseWarning, hostSystemMemoryUseWarningFlagHelp)
		flag.IntVar(&c.HostSystemMemoryUseWarning, "mw", defaultMemoryUseWarning, hostSystemMemoryUseWarningFlagHelp+" (shorthand)")

		flag.IntVar(&c.HostSystemMemoryUseCritical, "memory-usage-critical", defaultMemoryUseCritical, hostSystemMemoryUseCriticalFlagHelp)
		flag.IntVar(&c.HostSystemMemoryUseCritical, "mc", defaultMemoryUseCritical, hostSystemMemoryUseCriticalFlagHelp+" (shorthand)")

	case pluginType.HostSystemCPU:

		flag.StringVar(&c.DatacenterName, "dc-name", defaultDatacenterName, datacenterNameFlagHelp)

		flag.StringVar(&c.HostSystemName, "host-name", defaultHostSystemName, hostSystemNameFlagHelp)

		flag.IntVar(&c.HostSystemCPUUseWarning, "cpu-usage-warning", defaultCPUUseWarning, hostSystemCPUUseWarningFlagHelp)
		flag.IntVar(&c.HostSystemCPUUseWarning, "cw", defaultCPUUseWarning, hostSystemCPUUseWarningFlagHelp+" (shorthand)")

		flag.IntVar(&c.HostSystemCPUUseCritical, "cpu-usage-critical", defaultCPUUseCritical, hostSystemCPUUseCriticalFlagHelp)
		flag.IntVar(&c.HostSystemCPUUseCritical, "cc", defaultCPUUseCritical, hostSystemCPUUseCriticalFlagHelp+" (shorthand)")

	case pluginType.ResourcePoolsMemory:

		flag.Var(&c.IncludedResourcePools, "include-rp", vmIncludedResourcePoolsFlagHelp)
		flag.Var(&c.ExcludedResourcePools, "exclude-rp", vmExcludedResourcePoolsFlagHelp)

		flag.IntVar(&c.ResourcePoolsMemoryUseWarning, "memory-use-warning", defaultMemoryUseWarning, resourcePoolsMemoryUseWarningFlagHelp)
		flag.IntVar(&c.ResourcePoolsMemoryUseWarning, "mw", defaultMemoryUseWarning, resourcePoolsMemoryUseWarningFlagHelp+" (shorthand)")

		flag.IntVar(&c.ResourcePoolsMemoryUseCritical, "memory-use-critical", defaultMemoryUseCritical, resourcePoolsMemoryUseCriticalFlagHelp)
		flag.IntVar(&c.ResourcePoolsMemoryUseCritical, "mc", defaultMemoryUseCritical, resourcePoolsMemoryUseCriticalFlagHelp+" (shorthand)")

		flag.IntVar(&c.ResourcePoolsMemoryMaxAllowed, "memory-max-allowed", defaultResourcePoolsMemoryMaxAllowed, resourcePoolsMemoryMaxAllowedFlagHelp)
		flag.IntVar(&c.ResourcePoolsMemoryMaxAllowed, "mma", defaultResourcePoolsMemoryMaxAllowed, resourcePoolsMemoryMaxAllowedFlagHelp+" (shorthand)")

	case pluginType.VirtualCPUsAllocation:

		flag.Var(&c.IncludedResourcePools, "include-rp", vmIncludedResourcePoolsFlagHelp)
		flag.Var(&c.ExcludedResourcePools, "exclude-rp", vmExcludedResourcePoolsFlagHelp)
		flag.Var(&c.IgnoredVMs, "ignore-vm", ignoreVMsFlagHelp)
		flag.BoolVar(&c.PoweredOff, "powered-off", defaultPoweredOff, poweredOffFlagHelp)

		flag.IntVar(&c.VCPUsAllocatedWarning, "vcpus-warning", defaultVCPUsAllocatedWarning, vCPUsAllocatedWarningFlagHelp)
		flag.IntVar(&c.VCPUsAllocatedWarning, "vw", defaultVCPUsAllocatedWarning, vCPUsAllocatedWarningFlagHelp+" (shorthand)")

		flag.IntVar(&c.VCPUsAllocatedCritical, "vcpus-critical", defaultVCPUsAllocatedCritical, vCPUsAllocatedCriticalFlagHelp)
		flag.IntVar(&c.VCPUsAllocatedCritical, "vc", defaultVCPUsAllocatedCritical, vCPUsAllocatedCriticalFlagHelp+" (shorthand)")

		flag.IntVar(&c.VCPUsMaxAllowed, "vcpus-max-allowed", defaultVCPUsMaxAllowed, vCPUsAllocatedMaxAllowedFlagHelp)
		flag.IntVar(&c.VCPUsMaxAllowed, "vcma", defaultVCPUsMaxAllowed, vCPUsAllocatedMaxAllowedFlagHelp+" (shorthand)")

	case pluginType.VirtualHardwareVersion:

		flag.Var(&c.IncludedResourcePools, "include-rp", vmIncludedResourcePoolsFlagHelp)
		flag.Var(&c.ExcludedResourcePools, "exclude-rp", vmExcludedResourcePoolsFlagHelp)
		flag.Var(&c.IgnoredVMs, "ignore-vm", ignoreVMsFlagHelp)
		flag.BoolVar(&c.PoweredOff, "powered-off", defaultPoweredOff, poweredOffFlagHelp)

		flag.StringVar(&c.DatacenterName, "dc-name", defaultDatacenterName, datacenterNameFlagHelp)
		flag.StringVar(&c.HostSystemName, "host-name", defaultHostSystemName, hostSystemNameFlagHelp)
		flag.StringVar(&c.ClusterName, "cluster-name", defaultClusterName, clusterNameFlagHelp)

		flag.IntVar(&c.VirtualHardwareOutdatedByWarning, "outdated-by-warning", defaultVirtualHardwareOutdatedByWarning, virtualHardwareOutdatedByWarningFlagHelp)
		flag.IntVar(&c.VirtualHardwareOutdatedByWarning, "obw", defaultVirtualHardwareOutdatedByWarning, virtualHardwareOutdatedByWarningFlagHelp+" (shorthand)")

		flag.IntVar(&c.VirtualHardwareOutdatedByCritical, "outdated-by-critical", defaultVirtualHardwareOutdatedByCritical, virtualHardwareOutdatedByCriticalFlagHelp)
		flag.IntVar(&c.VirtualHardwareOutdatedByCritical, "obc", defaultVirtualHardwareOutdatedByCritical, virtualHardwareOutdatedByCriticalFlagHelp+" (shorthand)")

		flag.IntVar(&c.VirtualHardwareMinimumVersion, "minimum-version", defaultVirtualHardwareMinimumVersion, virtualHardwareMinimumVersionFlagHelp)
		flag.IntVar(&c.VirtualHardwareMinimumVersion, "mv", defaultVirtualHardwareMinimumVersion, virtualHardwareMinimumVersionFlagHelp+" (shorthand)")

		flag.BoolVar(&c.VirtualHardwareDefaultVersionIsMinimum, "default-is-min-version", defaultVirtualHardwareDefaultIsMinimum, virtualHardwareDefaultIsMinimumFlagHelp)
		flag.BoolVar(&c.VirtualHardwareDefaultVersionIsMinimum, "dimv", defaultVirtualHardwareDefaultIsMinimum, virtualHardwareDefaultIsMinimumFlagHelp+" (shorthand)")

	case pluginType.Host2Datastores2VMs:

		flag.Var(&c.IncludedResourcePools, "include-rp", vmIncludedResourcePoolsFlagHelp)
		flag.Var(&c.ExcludedResourcePools, "exclude-rp", vmExcludedResourcePoolsFlagHelp)
		flag.Var(&c.IgnoredVMs, "ignore-vm", ignoreVMsFlagHelp)
		flag.BoolVar(&c.PoweredOff, "powered-off", defaultPoweredOff, poweredOffFlagHelp)

		flag.Var(&c.IgnoredDatastores, "ignore-ds", ignoreDatastoreFlagHelp)

		flag.StringVar(&c.sharedCustomAttributeName, "ca-name", defaultCustomAttributeName, sharedCustomAttributeNameFlagHelp)
		flag.StringVar(&c.sharedCustomAttributePrefixSeparator, "ca-prefix-sep", defaultCustomAttributePrefixSeparator, sharedCustomAttributePrefixSeparatorFlagHelp)

		flag.StringVar(&c.hostCustomAttributeName, "host-ca-name", defaultCustomAttributeName, hostCustomAttributeNameFlagHelp)
		flag.StringVar(&c.hostCustomAttributePrefixSeparator, "host-ca-prefix-sep", defaultCustomAttributePrefixSeparator, hostCustomAttributePrefixSeparatorFlagHelp)

		flag.StringVar(&c.datastoreCustomAttributeName, "ds-ca-name", defaultCustomAttributeName, datastoreCustomAttributeNameFlagHelp)
		flag.StringVar(&c.datastoreCustomAttributePrefixSeparator, "ds-ca-prefix-sep", defaultCustomAttributePrefixSeparator, datastoreCustomAttributePrefixSeparatorFlagHelp)

		flag.BoolVar(&c.IgnoreMissingCustomAttribute, "ignore-missing-ca", defaultIgnoreMissingCustomAttribute, ignoreMissingCustomAttributeFlagHelp)

	}

	// Shared flags for all plugin types

	flag.StringVar(&c.Username, "username", defaultUsername, usernameFlagHelp)
	flag.StringVar(&c.Username, "u", defaultUsername, usernameFlagHelp+" (shorthand)")
	flag.StringVar(&c.Password, "password", defaultPassword, passwordFlagHelp)
	flag.StringVar(&c.Password, "pw", defaultPassword, passwordFlagHelp+" (shorthand)")

	// TODO: Is this actually needed?
	flag.StringVar(&c.Domain, "domain", defaultUserDomain, userDomainFlagHelp)

	flag.BoolVar(&c.TrustCert, "trust-cert", defaultTrustCert, trustCertFlagHelp)

	flag.BoolVar(&c.EmitBranding, "branding", defaultBranding, brandingFlagHelp)

	flag.StringVar(&c.Server, "s", defaultServer, serverFlagHelp+" (shorthand)")
	flag.StringVar(&c.Server, "server", defaultServer, serverFlagHelp)

	flag.IntVar(&c.Port, "p", defaultPort, portFlagHelp+" (shorthand)")
	flag.IntVar(&c.Port, "port", defaultPort, portFlagHelp)

	flag.IntVar(&c.timeout, "t", defaultPluginRuntimeTimeout, timeoutPluginRuntimeFlagHelp+" (shorthand)")
	flag.IntVar(&c.timeout, "timeout", defaultPluginRuntimeTimeout, timeoutPluginRuntimeFlagHelp)

	flag.StringVar(&c.LoggingLevel, "ll", defaultLogLevel, logLevelFlagHelp+" (shorthand)")
	flag.StringVar(&c.LoggingLevel, "log-level", defaultLogLevel, logLevelFlagHelp)

	flag.BoolVar(&c.ShowVersion, "v", defaultDisplayVersionAndExit, versionFlagHelp+" (shorthand)")
	flag.BoolVar(&c.ShowVersion, "version", defaultDisplayVersionAndExit, versionFlagHelp)

	// Allow our function to override the default Help output
	flag.Usage = Usage

	// parse flag definitions from the argument list
	flag.Parse()

}
