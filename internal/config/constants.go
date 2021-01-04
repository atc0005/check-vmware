// Copyright 2021 Adam Chalkley
//
// https://github.com/atc0005/check-vmware
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package config

const myAppName string = "check-cert"
const myAppURL string = "https://github.com/atc0005/check-vmware"

const (
	versionFlagHelp                  string = "Whether to display application version and then immediately exit application."
	logLevelFlagHelp                 string = "Sets log level to one of disabled, panic, fatal, error, warn, info, debug or trace."
	serverFlagHelp                   string = "The fully-qualified domain name or IP Address of the remote ESXi host or vCenter instance."
	trustCertFlagHelp                string = "Whether the certificate should be trusted as-is without validation. WARNING: TLS is susceptible to man-in-the-middle attacks if enabling this option."
	portFlagHelp                     string = "TCP port of the remote ESXi host or vCenter instance. This is usually 443 (HTTPS)."
	timeoutConnectFlagHelp           string = "Timeout value in seconds allowed before a plugin execution attempt is abandoned and an error returned."
	brandingFlagHelp                 string = "Toggles emission of branding details with plugin status details. This output is disabled by default."
	usernameFlagHelp                 string = "Username with permission to access specified ESXi host or vCenter instance."
	passwordFlagHelp                 string = "Password used to login to ESXi host or vCenter instance."
	userDomainFlagHelp               string = "(Optional) domain for user account used to login to ESXi host or vCenter instance."
	includedResourcePoolsFlagHelp    string = "Specifies a comma-separated list of Resource Pools that should be exclusively used when evaluating VMs. This option is incompatible with specifying a list of Resource Pools to ignore or exclude from evaluation."
	excludedResourcePoolsFlagHelp    string = "Specifies a comma-separated list of Resource Pools that should be ignored when evaluating VMs. This option is incompatible with specifying a list of Resource Pools to include for evaluation."
	ignoreVMsFlagHelp                string = "Specifies a comma-separated list of VM names that should be ignored or excluded from evaluation."
	poweredOffFlagHelp               string = "Toggles evaluation of powered off VMs in addition to powered on VMs. Evaluation of powered off VMs is disabled by default."
	vCPUsAllocatedMaxAllowedFlagHelp string = "Specifies the maximum amount of virtual CPUs (as a whole number) that we are allowed to allocate in the target VMware environment."
	vCPUsAllocatedCriticalFlagHelp   string = "Specifies the percentage of vCPUs allocation (as a whole number) when a CRITICAL threshold is reached."
	vCPUsAllocatedWarningFlagHelp    string = "Specifies the percentage of vCPUs allocation (as a whole number) when a WARNING threshold is reached."
)

// Default flag settings if not overridden by user input
const (
	defaultLogLevel               string = "info"
	defaultServer                 string = ""
	defaultTrustCert              bool   = false
	defaultUsername               string = ""
	defaultPassword               string = ""
	defaultUserDomain             string = ""
	defaultPort                   int    = 443
	defaultBranding               bool   = false
	defaultDisplayVersionAndExit  bool   = false
	defaultPoweredOff             bool   = false
	defaultVCPUsAllocatedCritical int    = 100
	defaultVCPUsAllocatedWarning  int    = 95

	// Intentionally set low to trigger validation failure if not specified by
	// the end user.
	defaultVCPUsMaxAllowedAllowed int = 0

	// Default timeout (in seconds) used when connecting to a remote server
	defaultConnectTimeout int = 10
)

// Plugin types provided by this project.
const (
	PluginTypeTools                 string = "vmware-tools"
	PluginTypeSnapshotsAge          string = "snapshots-age"
	PluginTypeSnapshotsSize         string = "snapshots-size"
	PluginTypeDatastoresSize        string = "datastore-size"
	PluginTypeResourcePoolsMemory   string = "resource-pools-memory"
	PluginTypeVirtualCPUsAllocation string = "virtual-cpus-allocation"
)
