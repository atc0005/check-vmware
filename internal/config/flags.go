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

	case pluginType.SnapshotsAge:

		flag.Var(&c.IncludedResourcePools, "include-rp", includedResourcePoolsFlagHelp)
		flag.Var(&c.ExcludedResourcePools, "exclude-rp", excludedResourcePoolsFlagHelp)
		flag.Var(&c.IgnoredVMs, "ignore-vm", ignoreVMsFlagHelp)

	// 	flag.IntVar(&c.AgeWarning, "w", defaultSnapshotAgeWarning, snapshotAgeWarningFlagHelp)
	// 	flag.IntVar(&c.AgeWarning, "age-warning", defaultSnapshotAgeWarning, snapshotAgeWarningFlagHelp)
	//
	// 	flag.IntVar(&c.AgeCritical, "c", defaultSnapshotAgeCritical, snapshotAgeCriticalFlagHelp)
	// 	flag.IntVar(&c.AgeCritical, "age-critical", defaultSnapshotAgeCritical, snapshotAgeCriticalFlagHelp)

	case pluginType.SnapshotsSize:

		flag.Var(&c.IncludedResourcePools, "include-rp", includedResourcePoolsFlagHelp)
		flag.Var(&c.ExcludedResourcePools, "exclude-rp", excludedResourcePoolsFlagHelp)
		flag.Var(&c.IgnoredVMs, "ignore-vm", ignoreVMsFlagHelp)

	// 	flag.IntVar(&c.SizeWarning, "w", defaultSnapshotSizeWarning, snapshotSizeWarningFlagHelp)
	// 	flag.IntVar(&c.SizeWarning, "size-warning", defaultSnapshotSizeWarning, snapshotSizeWarningFlagHelp)
	//
	// 	flag.IntVar(&c.SizeCritical, "c", defaultSnapshotSizeCritical, snapshotSizeCriticalFlagHelp)
	// 	flag.IntVar(&c.SizeCritical, "size-critical", defaultSnapshotSizeCritical, snapshotSizeCriticalFlagHelp)

	case pluginType.DatastoresSize:

	case pluginType.ResourcePoolsMemory:

		flag.Var(&c.IncludedResourcePools, "include-rp", includedResourcePoolsFlagHelp)
		flag.Var(&c.ExcludedResourcePools, "exclude-rp", excludedResourcePoolsFlagHelp)

	case pluginType.VirtualCPUsAllocation:

		flag.Var(&c.IncludedResourcePools, "include-rp", includedResourcePoolsFlagHelp)
		flag.Var(&c.ExcludedResourcePools, "exclude-rp", excludedResourcePoolsFlagHelp)
		flag.Var(&c.IgnoredVMs, "ignore-vm", ignoreVMsFlagHelp)

	}

	// Shared flags for all plugin types

	flag.StringVar(&c.Username, "username", defaultUsername, usernameFlagHelp)
	flag.StringVar(&c.Password, "password", defaultPassword, passwordFlagHelp)

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
