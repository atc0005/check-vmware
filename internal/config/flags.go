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

		flag.Var(&c.IncludedResourcePools, "include-rp", includedResourcePoolsFlagHelp)
		flag.Var(&c.ExcludedResourcePools, "exclude-rp", excludedResourcePoolsFlagHelp)
		flag.Var(&c.IgnoredVMs, "ignore-vm", ignoreVMsFlagHelp)
		flag.BoolVar(&c.PoweredOff, "powered-off", defaultPoweredOff, poweredOffFlagHelp)

	case pluginType.SnapshotsAge:

		flag.Var(&c.IncludedResourcePools, "include-rp", includedResourcePoolsFlagHelp)
		flag.Var(&c.ExcludedResourcePools, "exclude-rp", excludedResourcePoolsFlagHelp)
		flag.Var(&c.IgnoredVMs, "ignore-vm", ignoreVMsFlagHelp)

		// NOTE: This plugin is hard-coded to evaluate powered off and powered
		// on VMs equally. I'm not sure whether ignoring powered off VMs by
		// default makes sense for this particular plugin.
		//
		// Please share your feedback on this GitHub issue if you feel differently:
		// https://github.com/atc0005/check-vmware/issues/79
		//
		// flag.BoolVar(&c.PoweredOff, "powered-off", defaultPoweredOff, poweredOffFlagHelp)

		flag.IntVar(&c.SnapshotsAgeWarning, "aw", defaultSnapshotsAgeWarning, snapshotsAgeWarningFlagHelp)
		flag.IntVar(&c.SnapshotsAgeWarning, "age-warning", defaultSnapshotsAgeWarning, snapshotsAgeWarningFlagHelp)

		flag.IntVar(&c.SnapshotsAgeCritical, "ac", defaultSnapshotsAgeCritical, snapshotsAgeCriticalFlagHelp)
		flag.IntVar(&c.SnapshotsAgeCritical, "age-critical", defaultSnapshotsAgeCritical, snapshotsAgeCriticalFlagHelp)

	case pluginType.SnapshotsCount:

		flag.Var(&c.IncludedResourcePools, "include-rp", includedResourcePoolsFlagHelp)
		flag.Var(&c.ExcludedResourcePools, "exclude-rp", excludedResourcePoolsFlagHelp)
		flag.Var(&c.IgnoredVMs, "ignore-vm", ignoreVMsFlagHelp)

		// NOTE: This plugin is hard-coded to evaluate powered off and powered
		// on VMs equally. I'm not sure whether ignoring powered off VMs by
		// default makes sense for this particular plugin.
		//
		// Please share your feedback on this GitHub issue if you feel differently:
		// https://github.com/atc0005/check-vmware/issues/79
		//
		// flag.BoolVar(&c.PoweredOff, "powered-off", defaultPoweredOff, poweredOffFlagHelp)

		flag.IntVar(&c.SnapshotsCountWarning, "cw", defaultSnapshotsCountWarning, snapshotsCountWarningFlagHelp)
		flag.IntVar(&c.SnapshotsCountWarning, "count-warning", defaultSnapshotsCountWarning, snapshotsCountWarningFlagHelp)

		flag.IntVar(&c.SnapshotsCountCritical, "cc", defaultSnapshotsCountCritical, snapshotsCountCriticalFlagHelp)
		flag.IntVar(&c.SnapshotsCountCritical, "count-critical", defaultSnapshotsCountCritical, snapshotsCountCriticalFlagHelp)

	case pluginType.SnapshotsSize:

		flag.Var(&c.IncludedResourcePools, "include-rp", includedResourcePoolsFlagHelp)
		flag.Var(&c.ExcludedResourcePools, "exclude-rp", excludedResourcePoolsFlagHelp)
		flag.Var(&c.IgnoredVMs, "ignore-vm", ignoreVMsFlagHelp)

		// NOTE: This plugin is hard-coded to evaluate powered off and powered
		// on VMs equally. I'm not sure whether ignoring powered off VMs by
		// default makes sense for this particular plugin.
		//
		// Please share your feedback on this GitHub issue if you feel differently:
		// https://github.com/atc0005/check-vmware/issues/79
		//
		// flag.BoolVar(&c.PoweredOff, "powered-off", defaultPoweredOff, poweredOffFlagHelp)

		flag.IntVar(&c.SnapshotsSizeWarning, "sw", defaultSnapshotsSizeWarning, snapshotsSizeWarningFlagHelp)
		flag.IntVar(&c.SnapshotsSizeWarning, "size-warning", defaultSnapshotsSizeWarning, snapshotsSizeWarningFlagHelp)

		flag.IntVar(&c.SnapshotsSizeCritical, "sc", defaultSnapshotsSizeCritical, snapshotsSizeCriticalFlagHelp)
		flag.IntVar(&c.SnapshotsSizeCritical, "size-critical", defaultSnapshotsSizeCritical, snapshotsSizeCriticalFlagHelp)

	case pluginType.VirtualMachinePowerCycleUptime:

		flag.Var(&c.IncludedResourcePools, "include-rp", includedResourcePoolsFlagHelp)
		flag.Var(&c.ExcludedResourcePools, "exclude-rp", excludedResourcePoolsFlagHelp)
		flag.Var(&c.IgnoredVMs, "ignore-vm", ignoreVMsFlagHelp)

		// NOTE: This plugin is hard-coded to evaluate powered off and powered
		// on VMs equally. I'm not sure whether ignoring powered off VMs by
		// default makes sense for this particular plugin.
		//
		// Please share your feedback on this GitHub issue if you feel differently:
		// https://github.com/atc0005/check-vmware/issues/79
		//
		// flag.BoolVar(&c.PoweredOff, "powered-off", defaultPoweredOff, poweredOffFlagHelp)

		flag.IntVar(&c.VMPowerCycleUptimeWarning, "uw", defaultVMPowerCycleUptimeWarning, vmPowerCycleUptimeWarningFlagHelp)
		flag.IntVar(&c.VMPowerCycleUptimeWarning, "uptime-warning", defaultVMPowerCycleUptimeWarning, vmPowerCycleUptimeWarningFlagHelp)

		flag.IntVar(&c.VMPowerCycleUptimeCritical, "uc", defaultVMPowerCycleUptimeCritical, vmPowerCycleUptimeCriticalFlagHelp)
		flag.IntVar(&c.VMPowerCycleUptimeCritical, "uptime-critical", defaultVMPowerCycleUptimeCritical, vmPowerCycleUptimeCriticalFlagHelp)

	case pluginType.DatastoresSize:

		flag.StringVar(&c.DatacenterName, "dc-name", defaultDatacenterName, datacenterNameFlagHelp)

		flag.StringVar(&c.DatastoreName, "ds-name", defaultDatastoreName, datastoreNameFlagHelp)

		flag.IntVar(&c.DatastoreUsageWarning, "ds-usage-warning", defaultDatastoreUsageWarning, datastoreUsageWarningFlagHelp)
		flag.IntVar(&c.DatastoreUsageWarning, "dsuw", defaultDatastoreUsageWarning, datastoreUsageWarningFlagHelp+" (shorthand)")

		flag.IntVar(&c.DatastoreUsageCritical, "ds-usage-critical", defaultDatastoreUsageCritical, datastoreUsageCriticalFlagHelp)
		flag.IntVar(&c.DatastoreUsageCritical, "dsuc", defaultDatastoreUsageCritical, datastoreUsageCriticalFlagHelp+" (shorthand)")

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

		flag.Var(&c.IncludedResourcePools, "include-rp", includedResourcePoolsFlagHelp)
		flag.Var(&c.ExcludedResourcePools, "exclude-rp", excludedResourcePoolsFlagHelp)

		flag.IntVar(&c.ResourcePoolsMemoryUseWarning, "memory-use-warning", defaultMemoryUseWarning, resourcePoolsMemoryUseWarningFlagHelp)
		flag.IntVar(&c.ResourcePoolsMemoryUseWarning, "mw", defaultMemoryUseWarning, resourcePoolsMemoryUseWarningFlagHelp+" (shorthand)")

		flag.IntVar(&c.ResourcePoolsMemoryUseCritical, "memory-use-critical", defaultMemoryUseCritical, resourcePoolsMemoryUseCriticalFlagHelp)
		flag.IntVar(&c.ResourcePoolsMemoryUseCritical, "mc", defaultMemoryUseCritical, resourcePoolsMemoryUseCriticalFlagHelp+" (shorthand)")

		flag.IntVar(&c.ResourcePoolsMemoryMaxAllowed, "memory-max-allowed", defaultResourcePoolsMemoryMaxAllowed, resourcePoolsMemoryMaxAllowedFlagHelp)
		flag.IntVar(&c.ResourcePoolsMemoryMaxAllowed, "mma", defaultResourcePoolsMemoryMaxAllowed, resourcePoolsMemoryMaxAllowedFlagHelp+" (shorthand)")

	case pluginType.VirtualCPUsAllocation:

		flag.Var(&c.IncludedResourcePools, "include-rp", includedResourcePoolsFlagHelp)
		flag.Var(&c.ExcludedResourcePools, "exclude-rp", excludedResourcePoolsFlagHelp)
		flag.Var(&c.IgnoredVMs, "ignore-vm", ignoreVMsFlagHelp)
		flag.BoolVar(&c.PoweredOff, "powered-off", defaultPoweredOff, poweredOffFlagHelp)

		flag.IntVar(&c.VCPUsAllocatedWarning, "vcpus-warning", defaultVCPUsAllocatedWarning, vCPUsAllocatedWarningFlagHelp)
		flag.IntVar(&c.VCPUsAllocatedWarning, "vw", defaultVCPUsAllocatedWarning, vCPUsAllocatedWarningFlagHelp+" (shorthand)")

		flag.IntVar(&c.VCPUsAllocatedCritical, "vcpus-critical", defaultVCPUsAllocatedCritical, vCPUsAllocatedCriticalFlagHelp)
		flag.IntVar(&c.VCPUsAllocatedCritical, "vc", defaultVCPUsAllocatedCritical, vCPUsAllocatedCriticalFlagHelp+" (shorthand)")

		flag.IntVar(&c.VCPUsMaxAllowed, "vcpus-max-allowed", defaultVCPUsMaxAllowed, vCPUsAllocatedMaxAllowedFlagHelp)
		flag.IntVar(&c.VCPUsMaxAllowed, "vcma", defaultVCPUsMaxAllowed, vCPUsAllocatedMaxAllowedFlagHelp+" (shorthand)")

	case pluginType.VirtualHardwareVersion:

		flag.Var(&c.IncludedResourcePools, "include-rp", includedResourcePoolsFlagHelp)
		flag.Var(&c.ExcludedResourcePools, "exclude-rp", excludedResourcePoolsFlagHelp)
		flag.Var(&c.IgnoredVMs, "ignore-vm", ignoreVMsFlagHelp)
		flag.BoolVar(&c.PoweredOff, "powered-off", defaultPoweredOff, poweredOffFlagHelp)

	case pluginType.Host2Datastores2VMs:

		flag.Var(&c.IncludedResourcePools, "include-rp", includedResourcePoolsFlagHelp)
		flag.Var(&c.ExcludedResourcePools, "exclude-rp", excludedResourcePoolsFlagHelp)
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

	flag.IntVar(&c.timeout, "t", defaultConnectTimeout, timeoutConnectFlagHelp)
	flag.IntVar(&c.timeout, "timeout", defaultConnectTimeout, timeoutConnectFlagHelp)

	flag.StringVar(&c.LoggingLevel, "ll", defaultLogLevel, logLevelFlagHelp)
	flag.StringVar(&c.LoggingLevel, "log-level", defaultLogLevel, logLevelFlagHelp)

	flag.BoolVar(&c.ShowVersion, "v", defaultDisplayVersionAndExit, versionFlagHelp)
	flag.BoolVar(&c.ShowVersion, "version", defaultDisplayVersionAndExit, versionFlagHelp)

	// Allow our function to override the default Help output
	flag.Usage = Usage

	// parse flag definitions from the argument list
	flag.Parse()

}
