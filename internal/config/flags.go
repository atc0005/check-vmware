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

		flag.Var(&c.IncludedFolders, IncludeFolderIDFlagLong, vmIncludedFoldersFlagHelp)
		flag.Var(&c.ExcludedFolders, ExcludeFolderIDFlagLong, vmExcludedFoldersFlagHelp)

		flag.Var(&c.IncludedResourcePools, IncludeResourcePoolFlagLong, vmIncludedResourcePoolsFlagHelp)
		flag.Var(&c.ExcludedResourcePools, ExcludeResourcePoolFlagLong, vmExcludedResourcePoolsFlagHelp)
		flag.Var(&c.IgnoredVMs, IgnoreVMFlagLong, ignoreVMsFlagHelp)
		flag.BoolVar(&c.PoweredOff, IncludePoweredOffVMsFlagLong, defaultPoweredOff, poweredOffFlagHelp)

	case pluginType.SnapshotsAge:

		flag.Var(&c.IncludedFolders, IncludeFolderIDFlagLong, vmIncludedFoldersFlagHelp)
		flag.Var(&c.ExcludedFolders, ExcludeFolderIDFlagLong, vmExcludedFoldersFlagHelp)

		flag.Var(&c.IncludedResourcePools, IncludeResourcePoolFlagLong, vmIncludedResourcePoolsFlagHelp)
		flag.Var(&c.ExcludedResourcePools, ExcludeResourcePoolFlagLong, vmExcludedResourcePoolsFlagHelp)
		flag.Var(&c.IgnoredVMs, IgnoreVMFlagLong, ignoreVMsFlagHelp)

		// NOTE: This plugin is hard-coded to evaluate powered off and powered
		// on VMs equally. I'm not sure whether ignoring powered off VMs by
		// default makes sense for this particular plugin.
		//
		// Please share your feedback here if you feel differently:
		// https://github.com/atc0005/check-vmware/discussions/177
		//
		// flag.BoolVar(&c.PoweredOff, IncludePoweredOffVMsFlagLong, defaultPoweredOff, poweredOffFlagHelp)

		flag.IntVar(&c.SnapshotsAgeWarning, SnapshotAgeWarningFlagLong, defaultSnapshotsAgeWarning, snapshotsAgeWarningFlagHelp)
		flag.IntVar(&c.SnapshotsAgeWarning, SnapshotAgeWarningFlagShort, defaultSnapshotsAgeWarning, snapshotsAgeWarningFlagHelp+shorthandFlagSuffix)

		flag.IntVar(&c.SnapshotsAgeCritical, SnapshotAgeCriticalFlagLong, defaultSnapshotsAgeCritical, snapshotsAgeCriticalFlagHelp)
		flag.IntVar(&c.SnapshotsAgeCritical, SnapshotAgeCriticalFlagShort, defaultSnapshotsAgeCritical, snapshotsAgeCriticalFlagHelp+shorthandFlagSuffix)

	case pluginType.SnapshotsCount:

		flag.Var(&c.IncludedFolders, IncludeFolderIDFlagLong, vmIncludedFoldersFlagHelp)
		flag.Var(&c.ExcludedFolders, ExcludeFolderIDFlagLong, vmExcludedFoldersFlagHelp)

		flag.Var(&c.IncludedResourcePools, IncludeResourcePoolFlagLong, vmIncludedResourcePoolsFlagHelp)
		flag.Var(&c.ExcludedResourcePools, ExcludeResourcePoolFlagLong, vmExcludedResourcePoolsFlagHelp)
		flag.Var(&c.IgnoredVMs, IgnoreVMFlagLong, ignoreVMsFlagHelp)

		// NOTE: This plugin is hard-coded to evaluate powered off and powered
		// on VMs equally. I'm not sure whether ignoring powered off VMs by
		// default makes sense for this particular plugin.
		//
		// Please share your feedback here if you feel differently:
		// https://github.com/atc0005/check-vmware/discussions/177
		//
		// flag.BoolVar(&c.PoweredOff, IncludePoweredOffVMsFlagLong, defaultPoweredOff, poweredOffFlagHelp)

		flag.IntVar(&c.SnapshotsCountWarning, SnapshotCountWarningFlagLong, defaultSnapshotsCountWarning, snapshotsCountWarningFlagHelp)
		flag.IntVar(&c.SnapshotsCountWarning, SnapshotCountWarningFlagShort, defaultSnapshotsCountWarning, snapshotsCountWarningFlagHelp+shorthandFlagSuffix)

		flag.IntVar(&c.SnapshotsCountCritical, SnapshotCountCriticalFlagLong, defaultSnapshotsCountCritical, snapshotsCountCriticalFlagHelp)
		flag.IntVar(&c.SnapshotsCountCritical, SnapshotCountCriticalFlagShort, defaultSnapshotsCountCritical, snapshotsCountCriticalFlagHelp+shorthandFlagSuffix)

	case pluginType.SnapshotsSize:

		flag.Var(&c.IncludedFolders, IncludeFolderIDFlagLong, vmIncludedFoldersFlagHelp)
		flag.Var(&c.ExcludedFolders, ExcludeFolderIDFlagLong, vmExcludedFoldersFlagHelp)

		flag.Var(&c.IncludedResourcePools, IncludeResourcePoolFlagLong, vmIncludedResourcePoolsFlagHelp)
		flag.Var(&c.ExcludedResourcePools, ExcludeResourcePoolFlagLong, vmExcludedResourcePoolsFlagHelp)
		flag.Var(&c.IgnoredVMs, IgnoreVMFlagLong, ignoreVMsFlagHelp)

		// NOTE: This plugin is hard-coded to evaluate powered off and powered
		// on VMs equally. I'm not sure whether ignoring powered off VMs by
		// default makes sense for this particular plugin.
		//
		// Please share your feedback here if you feel differently:
		// https://github.com/atc0005/check-vmware/discussions/177
		//
		// flag.BoolVar(&c.PoweredOff, IncludePoweredOffVMsFlagLong, defaultPoweredOff, poweredOffFlagHelp)

		flag.IntVar(&c.SnapshotsSizeWarning, SnapshotSizeWarningFlagLong, defaultSnapshotsSizeWarning, snapshotsSizeWarningFlagHelp)
		flag.IntVar(&c.SnapshotsSizeWarning, SnapshotSizeWarningFlagShort, defaultSnapshotsSizeWarning, snapshotsSizeWarningFlagHelp+shorthandFlagSuffix)

		flag.IntVar(&c.SnapshotsSizeCritical, SnapshotSizeCriticalFlagLong, defaultSnapshotsSizeCritical, snapshotsSizeCriticalFlagHelp)
		flag.IntVar(&c.SnapshotsSizeCritical, SnapshotSizeCriticalFlagShort, defaultSnapshotsSizeCritical, snapshotsSizeCriticalFlagHelp+shorthandFlagSuffix)

	case pluginType.VirtualMachinePowerCycleUptime:

		flag.Var(&c.IncludedFolders, IncludeFolderIDFlagLong, vmIncludedFoldersFlagHelp)
		flag.Var(&c.ExcludedFolders, ExcludeFolderIDFlagLong, vmExcludedFoldersFlagHelp)

		flag.Var(&c.IncludedResourcePools, IncludeResourcePoolFlagLong, vmIncludedResourcePoolsFlagHelp)
		flag.Var(&c.ExcludedResourcePools, ExcludeResourcePoolFlagLong, vmExcludedResourcePoolsFlagHelp)
		flag.Var(&c.IgnoredVMs, IgnoreVMFlagLong, ignoreVMsFlagHelp)

		flag.IntVar(&c.VMPowerCycleUptimeWarning, PowerUptimeWarningFlagLong, defaultVMPowerCycleUptimeWarning, vmPowerCycleUptimeWarningFlagHelp)
		flag.IntVar(&c.VMPowerCycleUptimeWarning, PowerUptimeWarningFlagShort, defaultVMPowerCycleUptimeWarning, vmPowerCycleUptimeWarningFlagHelp+shorthandFlagSuffix)

		flag.IntVar(&c.VMPowerCycleUptimeCritical, PowerUptimeCriticalFlagLong, defaultVMPowerCycleUptimeCritical, vmPowerCycleUptimeCriticalFlagHelp)
		flag.IntVar(&c.VMPowerCycleUptimeCritical, PowerUptimeCriticalFlagShort, defaultVMPowerCycleUptimeCritical, vmPowerCycleUptimeCriticalFlagHelp+shorthandFlagSuffix)

	case pluginType.DiskConsolidation:

		flag.Var(&c.IncludedFolders, IncludeFolderIDFlagLong, vmIncludedFoldersFlagHelp)
		flag.Var(&c.ExcludedFolders, ExcludeFolderIDFlagLong, vmExcludedFoldersFlagHelp)

		flag.Var(&c.IncludedResourcePools, IncludeResourcePoolFlagLong, vmIncludedResourcePoolsFlagHelp)
		flag.Var(&c.ExcludedResourcePools, ExcludeResourcePoolFlagLong, vmExcludedResourcePoolsFlagHelp)
		flag.Var(&c.IgnoredVMs, IgnoreVMFlagLong, ignoreVMsFlagHelp)
		flag.BoolVar(&c.TriggerReloadStateData, TriggerReloadFlagLong, defaultTriggerReloadStateData, triggerReloadStateDataFlagHelp)

		// NOTE: This plugin is hard-coded to evaluate powered off and powered
		// on VMs equally. I'm not sure whether ignoring powered off VMs by
		// default makes sense for this particular plugin.
		//
		// Please share your feedback here if you feel differently:
		// https://github.com/atc0005/check-vmware/discussions/176
		//
		// Please expand on some use cases for ignoring powered off VMs by default.
		//
		// flag.BoolVar(&c.PoweredOff, IncludePoweredOffVMsFlagLong, defaultPoweredOff, poweredOffFlagHelp)

	case pluginType.InteractiveQuestion:

		flag.Var(&c.IncludedFolders, IncludeFolderIDFlagLong, vmIncludedFoldersFlagHelp)
		flag.Var(&c.ExcludedFolders, ExcludeFolderIDFlagLong, vmExcludedFoldersFlagHelp)

		flag.Var(&c.IncludedResourcePools, IncludeResourcePoolFlagLong, vmIncludedResourcePoolsFlagHelp)
		flag.Var(&c.ExcludedResourcePools, ExcludeResourcePoolFlagLong, vmExcludedResourcePoolsFlagHelp)
		flag.Var(&c.IgnoredVMs, IgnoreVMFlagLong, ignoreVMsFlagHelp)

	case pluginType.Alarms:

		flag.Var(&c.DatacenterNames, DatacenterNameFlagLong, datacenterNamesFlagHelp)
		flag.Var(&c.IncludedAlarmEntityTypes, AlarmIncludeEntityTypeFlagLong, includedAlarmEntityTypesFlagHelp)
		flag.Var(&c.ExcludedAlarmEntityTypes, AlarmExcludeEntityTypeFlagLong, excludedAlarmEntityTypesFlagHelp)

		flag.BoolVar(&c.EvaluateAcknowledgedAlarms, AlarmEvalAcknowledgedFlagLong, defaultEvaluateAcknowledgedAlarms, evaluateAcknowledgedTriggeredAlarmFlagHelp)

		flag.Var(&c.IncludedAlarmNames, AlarmIncludeNameFlagLong, includedAlarmNamesFlagHelp)
		flag.Var(&c.ExcludedAlarmNames, AlarmExcludeNameFlagLong, excludedAlarmNamesFlagHelp)

		flag.Var(&c.IncludedAlarmDescriptions, AlarmIncludeDescFlagLong, includedAlarmDescriptionsFlagHelp)
		flag.Var(&c.ExcludedAlarmDescriptions, AlarmExcludeDescFlagLong, excludedAlarmDescriptionsFlagHelp)

		flag.Var(&c.includedAlarmStatuses, AlarmIncludeStatusFlagLong, includedAlarmStatusesFlagHelp)
		flag.Var(&c.excludedAlarmStatuses, AlarmExcludeStatusFlagLong, excludedAlarmStatusesFlagHelp)

		flag.Var(&c.IncludedAlarmEntityNames, AlarmIncludeEntityNameFlagLong, includedAlarmEntityNamesFlagHelp)
		flag.Var(&c.ExcludedAlarmEntityNames, AlarmExcludeEntityNameFlagLong, excludedAlarmEntityNamesFlagHelp)

		flag.Var(&c.IncludedAlarmEntityResourcePools, AlarmIncludeEntityRPoolFlagLong, includedAlarmEntityResourcePoolsFlagHelp)
		flag.Var(&c.ExcludedAlarmEntityResourcePools, AlarmExcludeEntityRPoolFlagLong, excludedAlarmEntityResourcePoolsFlagHelp)

	case pluginType.DatastoresSpace:

		flag.StringVar(&c.DatacenterName, DatacenterNameFlagLong, defaultDatacenterName, datacenterNameFlagHelp)

		flag.StringVar(&c.DatastoreName, DatastoreNameFlagLong, defaultDatastoreName, datastoreNameFlagHelp)

		flag.IntVar(&c.DatastoreSpaceUsageWarning, DatastoreSpaceUsageWarningFlagLong, defaultDatastoreSpaceUsageWarning, datastoreSpaceUsageWarningFlagHelp)
		flag.IntVar(&c.DatastoreSpaceUsageWarning, DatastoreSpaceUsageWarningFlagShort, defaultDatastoreSpaceUsageWarning, datastoreSpaceUsageWarningFlagHelp+shorthandFlagSuffix)

		flag.IntVar(&c.DatastoreSpaceUsageCritical, DatastoreSpaceUsageCriticalFlagLong, defaultDatastoreSpaceUsageCritical, datastoreSpaceUsageCriticalFlagHelp)
		flag.IntVar(&c.DatastoreSpaceUsageCritical, DatastoreSpaceUsageCriticalFlagShort, defaultDatastoreSpaceUsageCritical, datastoreSpaceUsageCriticalFlagHelp+shorthandFlagSuffix)

	case pluginType.DatastoresPerformance:

		flag.StringVar(&c.DatacenterName, DatacenterNameFlagLong, defaultDatacenterName, datacenterNameFlagHelp)

		flag.StringVar(&c.DatastoreName, DatastoreNameFlagLong, defaultDatastoreName, datastoreNameFlagHelp)

		flag.BoolVar(&c.IgnoreMissingDatastorePerfMetrics, DatastorePerformanceIgnoreMissingMetricsFlagLong, defaultIgnoreMissingDatastoreMetrics, ignoreMissingDatastorePerfMetricsFlagHelp)
		flag.BoolVar(&c.IgnoreMissingDatastorePerfMetrics, DatastorePerformanceIgnoreMissingMetricsFlagShort, defaultIgnoreMissingDatastoreMetrics, ignoreMissingDatastorePerfMetricsFlagHelp+shorthandFlagSuffix)

		flag.BoolVar(&c.HideHistoricalDatastorePerfMetricSets, DatastorePerformanceHideHistoricalMetricSetsFlagLong, defaultHideHistoricalDatastorePerfMetricSets, hideHistoricalDatastorePerfMetricSetsFlagHelp)
		flag.BoolVar(&c.HideHistoricalDatastorePerfMetricSets, DatastorePerformanceHideHistoricalMetricSetsFlagShort, defaultHideHistoricalDatastorePerfMetricSets, hideHistoricalDatastorePerfMetricSetsFlagHelp+shorthandFlagSuffix)

		flag.Var(&c.datastoreReadLatencyWarning, DatastorePerformanceReadLatencyWarningFlagLong, datastoreReadLatencyWarningFlagHelp)
		flag.Var(&c.datastoreReadLatencyWarning, DatastorePerformanceReadLatencyWarningFlagShort, datastoreReadLatencyWarningFlagHelp+shorthandFlagSuffix)

		flag.Var(&c.datastoreReadLatencyCritical, DatastorePerformanceReadLatencyCriticalFlagLong, datastoreReadLatencyCriticalFlagHelp)
		flag.Var(&c.datastoreReadLatencyCritical, DatastorePerformanceReadLatencyCriticalFlagShort, datastoreReadLatencyCriticalFlagHelp+shorthandFlagSuffix)

		flag.Var(&c.datastoreWriteLatencyWarning, DatastorePerformanceWriteLatencyWarningFlagLong, datastoreWriteLatencyWarningFlagHelp)
		flag.Var(&c.datastoreWriteLatencyWarning, DatastorePerformanceWriteLatencyWarningFlagShort, datastoreWriteLatencyWarningFlagHelp+shorthandFlagSuffix)

		flag.Var(&c.datastoreWriteLatencyCritical, DatastorePerformanceWriteLatencyCriticalFlagLong, datastoreWriteLatencyCriticalFlagHelp)
		flag.Var(&c.datastoreWriteLatencyCritical, DatastorePerformanceWriteLatencyCriticalFlagShort, datastoreWriteLatencyCriticalFlagHelp+shorthandFlagSuffix)

		flag.Var(&c.datastoreVMLatencyWarning, DatastorePerformanceVMLatencyWarningFlagLong, datastoreVMLatencyWarningFlagHelp)
		flag.Var(&c.datastoreVMLatencyWarning, DatastorePerformanceVMLatencyWarningFlagShort, datastoreVMLatencyWarningFlagHelp+shorthandFlagSuffix)

		flag.Var(&c.datastoreVMLatencyCritical, DatastorePerformanceVMLatencyCriticalFlagLong, datastoreVMLatencyCriticalFlagHelp)
		flag.Var(&c.datastoreVMLatencyCritical, DatastorePerformanceVMLatencyCriticalFlagShort, datastoreVMLatencyCriticalFlagHelp+shorthandFlagSuffix)

		flag.Var(&c.datastorePerformancePercentileSet, DatastoreLatencyPercentileSetFlagLong, datastoreLatencyPercintileSetFlagHelp)
		flag.Var(&c.datastorePerformancePercentileSet, DatastoreLatencyPercentileSetFlagShort, datastoreLatencyPercintileSetFlagHelp+shorthandFlagSuffix)

	case pluginType.HostSystemMemory:

		flag.StringVar(&c.DatacenterName, DatacenterNameFlagLong, defaultDatacenterName, datacenterNameFlagHelp)

		flag.StringVar(&c.HostSystemName, HostNameFlagLong, defaultHostSystemName, hostSystemNameFlagHelp)

		flag.IntVar(&c.HostSystemMemoryUseWarning, HostMemoryUsageWarningFlagLong, defaultMemoryUseWarning, hostSystemMemoryUseWarningFlagHelp)
		flag.IntVar(&c.HostSystemMemoryUseWarning, HostMemoryUsageWarningFlagShort, defaultMemoryUseWarning, hostSystemMemoryUseWarningFlagHelp+shorthandFlagSuffix)

		flag.IntVar(&c.HostSystemMemoryUseCritical, HostMemoryUsageCriticalFlagLong, defaultMemoryUseCritical, hostSystemMemoryUseCriticalFlagHelp)
		flag.IntVar(&c.HostSystemMemoryUseCritical, HostMemoryUsageCriticalFlagShort, defaultMemoryUseCritical, hostSystemMemoryUseCriticalFlagHelp+shorthandFlagSuffix)

	case pluginType.HostSystemCPU:

		flag.StringVar(&c.DatacenterName, DatacenterNameFlagLong, defaultDatacenterName, datacenterNameFlagHelp)

		flag.StringVar(&c.HostSystemName, HostNameFlagLong, defaultHostSystemName, hostSystemNameFlagHelp)

		flag.IntVar(&c.HostSystemCPUUseWarning, HostCPUUsageWarningFlagLong, defaultCPUUseWarning, hostSystemCPUUseWarningFlagHelp)
		flag.IntVar(&c.HostSystemCPUUseWarning, HostCPUUsageWarningFlagShort, defaultCPUUseWarning, hostSystemCPUUseWarningFlagHelp+shorthandFlagSuffix)

		flag.IntVar(&c.HostSystemCPUUseCritical, HostCPUUsageCriticalFlagLong, defaultCPUUseCritical, hostSystemCPUUseCriticalFlagHelp)
		flag.IntVar(&c.HostSystemCPUUseCritical, HostCPUUsageCriticalFlagShort, defaultCPUUseCritical, hostSystemCPUUseCriticalFlagHelp+shorthandFlagSuffix)

	case pluginType.ResourcePoolsMemory:

		flag.Var(&c.IncludedResourcePools, IncludeResourcePoolFlagLong, vmIncludedResourcePoolsFlagHelp)
		flag.Var(&c.ExcludedResourcePools, ExcludeResourcePoolFlagLong, vmExcludedResourcePoolsFlagHelp)

		flag.IntVar(&c.ResourcePoolsMemoryUseWarning, RPMemoryUseWarningFlagLong, defaultMemoryUseWarning, resourcePoolsMemoryUseWarningFlagHelp)
		flag.IntVar(&c.ResourcePoolsMemoryUseWarning, RPMemoryUseWarningFlagShort, defaultMemoryUseWarning, resourcePoolsMemoryUseWarningFlagHelp+shorthandFlagSuffix)

		flag.IntVar(&c.ResourcePoolsMemoryUseCritical, RPMemoryUseCriticalFlagLong, defaultMemoryUseCritical, resourcePoolsMemoryUseCriticalFlagHelp)
		flag.IntVar(&c.ResourcePoolsMemoryUseCritical, RPMemoryUseCriticalFlagShort, defaultMemoryUseCritical, resourcePoolsMemoryUseCriticalFlagHelp+shorthandFlagSuffix)

		flag.IntVar(&c.ResourcePoolsMemoryMaxAllowed, RPMemoryMaxAllowedFlagLong, defaultResourcePoolsMemoryMaxAllowed, resourcePoolsMemoryMaxAllowedFlagHelp)
		flag.IntVar(&c.ResourcePoolsMemoryMaxAllowed, RPMemoryMaxAllowedFlagShort, defaultResourcePoolsMemoryMaxAllowed, resourcePoolsMemoryMaxAllowedFlagHelp+shorthandFlagSuffix)

	case pluginType.VirtualCPUsAllocation:

		flag.Var(&c.IncludedFolders, IncludeFolderIDFlagLong, vmIncludedFoldersFlagHelp)
		flag.Var(&c.ExcludedFolders, ExcludeFolderIDFlagLong, vmExcludedFoldersFlagHelp)

		flag.Var(&c.IncludedResourcePools, IncludeResourcePoolFlagLong, vmIncludedResourcePoolsFlagHelp)
		flag.Var(&c.ExcludedResourcePools, ExcludeResourcePoolFlagLong, vmExcludedResourcePoolsFlagHelp)
		flag.Var(&c.IgnoredVMs, IgnoreVMFlagLong, ignoreVMsFlagHelp)
		flag.BoolVar(&c.PoweredOff, IncludePoweredOffVMsFlagLong, defaultPoweredOff, poweredOffFlagHelp)

		flag.IntVar(&c.VCPUsAllocatedWarning, VirtualCPUsWarningFlagLong, defaultVCPUsAllocatedWarning, vCPUsAllocatedWarningFlagHelp)
		flag.IntVar(&c.VCPUsAllocatedWarning, VirtualCPUsWarningFlagShort, defaultVCPUsAllocatedWarning, vCPUsAllocatedWarningFlagHelp+shorthandFlagSuffix)

		flag.IntVar(&c.VCPUsAllocatedCritical, VirtualCPUsCriticalFlagLong, defaultVCPUsAllocatedCritical, vCPUsAllocatedCriticalFlagHelp)
		flag.IntVar(&c.VCPUsAllocatedCritical, VirtualCPUsCriticalFlagShort, defaultVCPUsAllocatedCritical, vCPUsAllocatedCriticalFlagHelp+shorthandFlagSuffix)

		flag.IntVar(&c.VCPUsMaxAllowed, VirtualCPUsMaxAllowedFlagLong, defaultVCPUsMaxAllowed, vCPUsAllocatedMaxAllowedFlagHelp)
		flag.IntVar(&c.VCPUsMaxAllowed, VirtualCPUsMaxAllowedFlagShort, defaultVCPUsMaxAllowed, vCPUsAllocatedMaxAllowedFlagHelp+shorthandFlagSuffix)

	case pluginType.VirtualHardwareVersion:

		flag.Var(&c.IncludedFolders, IncludeFolderIDFlagLong, vmIncludedFoldersFlagHelp)
		flag.Var(&c.ExcludedFolders, ExcludeFolderIDFlagLong, vmExcludedFoldersFlagHelp)

		flag.Var(&c.IncludedResourcePools, IncludeResourcePoolFlagLong, vmIncludedResourcePoolsFlagHelp)
		flag.Var(&c.ExcludedResourcePools, ExcludeResourcePoolFlagLong, vmExcludedResourcePoolsFlagHelp)
		flag.Var(&c.IgnoredVMs, IgnoreVMFlagLong, ignoreVMsFlagHelp)
		flag.BoolVar(&c.PoweredOff, IncludePoweredOffVMsFlagLong, defaultPoweredOff, poweredOffFlagHelp)

		flag.StringVar(&c.DatacenterName, DatacenterNameFlagLong, defaultDatacenterName, datacenterNameFlagHelp)
		flag.StringVar(&c.HostSystemName, HostNameFlagLong, defaultHostSystemName, hostSystemNameFlagHelp)
		flag.StringVar(&c.ClusterName, ClusterNameFlagLong, defaultClusterName, clusterNameFlagHelp)

		flag.IntVar(&c.VirtualHardwareOutdatedByWarning, OutdatedByWarningFlagLong, defaultVirtualHardwareOutdatedByWarning, virtualHardwareOutdatedByWarningFlagHelp)
		flag.IntVar(&c.VirtualHardwareOutdatedByWarning, OutdatedByWarningFlagShort, defaultVirtualHardwareOutdatedByWarning, virtualHardwareOutdatedByWarningFlagHelp+shorthandFlagSuffix)

		flag.IntVar(&c.VirtualHardwareOutdatedByCritical, OutdatedByCriticalFlagLong, defaultVirtualHardwareOutdatedByCritical, virtualHardwareOutdatedByCriticalFlagHelp)
		flag.IntVar(&c.VirtualHardwareOutdatedByCritical, OutdatedByCriticalFlagShort, defaultVirtualHardwareOutdatedByCritical, virtualHardwareOutdatedByCriticalFlagHelp+shorthandFlagSuffix)

		flag.IntVar(&c.VirtualHardwareMinimumVersion, MinimumVersionFlagLong, defaultVirtualHardwareMinimumVersion, virtualHardwareMinimumVersionFlagHelp)
		flag.IntVar(&c.VirtualHardwareMinimumVersion, MinimumVersionFlagShort, defaultVirtualHardwareMinimumVersion, virtualHardwareMinimumVersionFlagHelp+shorthandFlagSuffix)

		flag.BoolVar(&c.VirtualHardwareDefaultVersionIsMinimum, DefaultIsMinimumVersionFlagLong, defaultVirtualHardwareDefaultIsMinimum, virtualHardwareDefaultIsMinimumFlagHelp)
		flag.BoolVar(&c.VirtualHardwareDefaultVersionIsMinimum, DefaultIsMinimumVersionFlagShort, defaultVirtualHardwareDefaultIsMinimum, virtualHardwareDefaultIsMinimumFlagHelp+shorthandFlagSuffix)

	case pluginType.Host2Datastores2VMs:

		flag.Var(&c.IncludedFolders, IncludeFolderIDFlagLong, vmIncludedFoldersFlagHelp)
		flag.Var(&c.ExcludedFolders, ExcludeFolderIDFlagLong, vmExcludedFoldersFlagHelp)

		flag.Var(&c.IncludedResourcePools, IncludeResourcePoolFlagLong, vmIncludedResourcePoolsFlagHelp)
		flag.Var(&c.ExcludedResourcePools, ExcludeResourcePoolFlagLong, vmExcludedResourcePoolsFlagHelp)
		flag.Var(&c.IgnoredVMs, IgnoreVMFlagLong, ignoreVMsFlagHelp)
		flag.BoolVar(&c.PoweredOff, IncludePoweredOffVMsFlagLong, defaultPoweredOff, poweredOffFlagHelp)

		flag.Var(&c.IgnoredDatastores, IgnoreDatastoreFlagLong, ignoreDatastoreFlagHelp)

		flag.StringVar(&c.sharedCustomAttributeName, CustomAttributeNameFlagLong, defaultCustomAttributeName, sharedCustomAttributeNameFlagHelp)
		flag.StringVar(&c.sharedCustomAttributePrefixSeparator, CustomAttributePrefixSeparatorFlagLong, defaultCustomAttributePrefixSeparator, sharedCustomAttributePrefixSeparatorFlagHelp)

		flag.StringVar(&c.hostCustomAttributeName, HostCustomAttributeNameFlagLong, defaultCustomAttributeName, hostCustomAttributeNameFlagHelp)
		flag.StringVar(&c.hostCustomAttributePrefixSeparator, HostCustomAttributePrefixSeparatorFlagLong, defaultCustomAttributePrefixSeparator, hostCustomAttributePrefixSeparatorFlagHelp)

		flag.StringVar(&c.datastoreCustomAttributeName, DatastoreCustomAttributeNameFlagLong, defaultCustomAttributeName, datastoreCustomAttributeNameFlagHelp)
		flag.StringVar(&c.datastoreCustomAttributePrefixSeparator, DatastoreCustomAttributePrefixSeparatorFlagLong, defaultCustomAttributePrefixSeparator, datastoreCustomAttributePrefixSeparatorFlagHelp)

		flag.BoolVar(&c.IgnoreMissingCustomAttribute, CustomAttributeIgnoreMissingCAFlagLong, defaultIgnoreMissingCustomAttribute, ignoreMissingCustomAttributeFlagHelp)

	case pluginType.VirtualMachineLastBackupViaCA:

		// NOTE: This plugin is hard-coded to evaluate powered off and powered
		// on VMs equally. I'm not sure whether ignoring powered off VMs by
		// default makes sense for this particular plugin.
		//
		// Please share your feedback here if you feel differently:
		// https://github.com/atc0005/check-vmware/discussions
		//
		// Please expand on some use cases for ignoring powered off VMs by default.
		//
		// flag.BoolVar(&c.PoweredOff, IncludePoweredOffVMsFlagLong, defaultPoweredOff, poweredOffFlagHelp)

		flag.Var(&c.IncludedFolders, IncludeFolderIDFlagLong, vmIncludedFoldersFlagHelp)
		flag.Var(&c.ExcludedFolders, ExcludeFolderIDFlagLong, vmExcludedFoldersFlagHelp)

		flag.Var(&c.IncludedResourcePools, IncludeResourcePoolFlagLong, vmIncludedResourcePoolsFlagHelp)
		flag.Var(&c.ExcludedResourcePools, ExcludeResourcePoolFlagLong, vmExcludedResourcePoolsFlagHelp)
		flag.Var(&c.IgnoredVMs, IgnoreVMFlagLong, ignoreVMsFlagHelp)

		flag.StringVar(&c.VMBackupDateCustomAttribute, BackupDateCAFlagLong, defaultVMBackupDateCustomAttribute, vmBackupDateCustomAttributeFlagHelp)
		flag.StringVar(&c.VMBackupMetadataCustomAttribute, BackupMetadataCAFlagLong, defaultVMBackupMetadataCustomAttribute, vmBackupMetadataCustomAttributeFlagHelp)
		flag.StringVar(&c.VMBackupDateFormat, BackupDateFormatFlagLong, defaultVMBackupDateFormat, vmBackupDateFormatFlagHelp)
		flag.StringVar(&c.VMBackupDateTimezone, BackupDateTimezoneFlagLong, defaultVMBackupDateTimezone, vmBackupDateTimezoneFlagHelp)

		flag.IntVar(&c.VMBackupAgeWarning, BackupAgeWarningFlagLong, defaultVMBackupAgeWarning, vmBackupAgeWarningFlagHelp)
		flag.IntVar(&c.VMBackupAgeWarning, BackupAgeWarningFlagShort, defaultVMBackupAgeWarning, vmBackupAgeWarningFlagHelp+shorthandFlagSuffix)

		flag.IntVar(&c.VMBackupAgeCritical, BackupAgeCriticalFlagLong, defaultVMBackupAgeCritical, vmBackupAgeCriticalFlagHelp)
		flag.IntVar(&c.VMBackupAgeCritical, BackupAgeCriticalFlagShort, defaultVMBackupAgeCritical, vmBackupAgeCriticalFlagHelp+shorthandFlagSuffix)

	case pluginType.VirtualMachineList:

		// FIXME: Need to update README to include this flag.
		//
		// TODO: Consider moving this to the "common" flags section and allow
		// specifying it for ALL plugins. This will provide a useful way to
		// collect a total VMs count.
		// flag.Var(&c.DatacenterNames, DatacenterNameFlagLong, datacenterNamesFlagHelp)

		flag.BoolVar(&c.PoweredOff, IncludePoweredOffVMsFlagLong, defaultPoweredOff, poweredOffFlagHelp)

		flag.Var(&c.IncludedFolders, IncludeFolderIDFlagLong, vmIncludedFoldersFlagHelp)
		flag.Var(&c.ExcludedFolders, ExcludeFolderIDFlagLong, vmExcludedFoldersFlagHelp)

		flag.Var(&c.IncludedResourcePools, IncludeResourcePoolFlagLong, vmIncludedResourcePoolsFlagHelp)
		flag.Var(&c.ExcludedResourcePools, ExcludeResourcePoolFlagLong, vmExcludedResourcePoolsFlagHelp)
		flag.Var(&c.IgnoredVMs, IgnoreVMFlagLong, ignoreVMsFlagHelp)

	}

	// Shared flags for all plugin types

	flag.StringVar(&c.Username, UsernameFlagLong, defaultUsername, usernameFlagHelp)
	flag.StringVar(&c.Username, UsernameFlagShort, defaultUsername, usernameFlagHelp+shorthandFlagSuffix)
	flag.StringVar(&c.Password, PasswordFlagLong, defaultPassword, passwordFlagHelp)
	flag.StringVar(&c.Password, PasswordFlagShort, defaultPassword, passwordFlagHelp+shorthandFlagSuffix)

	flag.StringVar(&c.Domain, DomainFlagLong, defaultUserDomain, userDomainFlagHelp)

	flag.BoolVar(&c.TrustCert, TrustCertFlagLong, defaultTrustCert, trustCertFlagHelp)

	flag.BoolVar(&c.EmitBranding, BrandingFlag, defaultBranding, brandingFlagHelp)

	flag.StringVar(&c.Server, ServerFlagLong, defaultServer, serverFlagHelp)
	flag.StringVar(&c.Server, ServerFlagShort, defaultServer, serverFlagHelp+shorthandFlagSuffix)

	flag.IntVar(&c.Port, PortFlagLong, defaultPort, portFlagHelp)
	flag.IntVar(&c.Port, PortFlagShort, defaultPort, portFlagHelp+shorthandFlagSuffix)

	flag.IntVar(&c.timeout, TimeoutFlagLong, defaultPluginRuntimeTimeout, timeoutPluginRuntimeFlagHelp)
	flag.IntVar(&c.timeout, TimeoutFlagShort, defaultPluginRuntimeTimeout, timeoutPluginRuntimeFlagHelp+shorthandFlagSuffix)

	flag.StringVar(&c.LoggingLevel, LogLevelFlagLong, defaultLogLevel, logLevelFlagHelp)
	flag.StringVar(&c.LoggingLevel, LogLevelFlagShort, defaultLogLevel, logLevelFlagHelp+shorthandFlagSuffix)

	flag.BoolVar(&c.ShowVersion, VersionFlagLong, defaultDisplayVersionAndExit, versionFlagHelp)
	flag.BoolVar(&c.ShowVersion, VersionFlagShort, defaultDisplayVersionAndExit, versionFlagHelp+shorthandFlagSuffix)

	// Allow our function to override the default Help output
	flag.Usage = Usage

	// parse flag definitions from the argument list
	flag.Parse()

}
