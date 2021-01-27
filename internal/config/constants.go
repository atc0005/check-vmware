// Copyright 2021 Adam Chalkley
//
// https://github.com/atc0005/check-vmware
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package config

const myAppName string = "check-vmware"
const myAppURL string = "https://github.com/atc0005/" + myAppName

const (
	versionFlagHelp                                 string = "Whether to display application version and then immediately exit application."
	logLevelFlagHelp                                string = "Sets log level to one of disabled, panic, fatal, error, warn, info, debug or trace."
	serverFlagHelp                                  string = "The fully-qualified domain name or IP Address of the remote ESXi host or vCenter instance."
	trustCertFlagHelp                               string = "Whether the certificate should be trusted as-is without validation. WARNING: TLS is susceptible to man-in-the-middle attacks if enabling this option."
	portFlagHelp                                    string = "TCP port of the remote ESXi host or vCenter instance. This is usually 443 (HTTPS)."
	timeoutConnectFlagHelp                          string = "Timeout value in seconds allowed before a plugin execution attempt is abandoned and an error returned."
	brandingFlagHelp                                string = "Toggles emission of branding details with plugin status details. This output is disabled by default."
	usernameFlagHelp                                string = "Username with permission to access specified ESXi host or vCenter instance."
	passwordFlagHelp                                string = "Password used to login to ESXi host or vCenter instance."
	userDomainFlagHelp                              string = "(Optional) domain for user account used to login to ESXi host or vCenter instance."
	includedResourcePoolsFlagHelp                   string = "Specifies a comma-separated list of Resource Pools that should be exclusively used when evaluating VMs. Specifying this option will also exclude any VMs from evaluation that are *outside* of a Resource Pool. This option is incompatible with specifying a list of Resource Pools to ignore or exclude from evaluation."
	excludedResourcePoolsFlagHelp                   string = "Specifies a comma-separated list of Resource Pools that should be ignored when evaluating VMs. This option is incompatible with specifying a list of Resource Pools to include for evaluation."
	ignoreVMsFlagHelp                               string = "Specifies a comma-separated list of VM names that should be ignored or excluded from evaluation."
	poweredOffFlagHelp                              string = "Toggles evaluation of powered off VMs in addition to powered on VMs. Evaluation of powered off VMs is disabled by default."
	vCPUsAllocatedMaxAllowedFlagHelp                string = "Specifies the maximum amount of virtual CPUs (as a whole number) that we are allowed to allocate in the target VMware environment."
	vCPUsAllocatedCriticalFlagHelp                  string = "Specifies the percentage of vCPUs allocation (as a whole number) when a CRITICAL threshold is reached."
	vCPUsAllocatedWarningFlagHelp                   string = "Specifies the percentage of vCPUs allocation (as a whole number) when a WARNING threshold is reached."
	hostCustomAttributeNameFlagHelp                 string = "Custom Attribute name specific to host ESXi systems. Optional if specifying shared custom attribute flag."
	hostCustomAttributePrefixSeparatorFlagHelp      string = "Custom Attribute prefix separator specific to host ESXi systems. Skip if using Custom Attribute values as-is for comparison, otherwise optional if specifying shared custom attribute prefix separator, or using the default separator."
	datastoreCustomAttributeNameFlagHelp            string = "Custom Attribute name specific to datastores. Optional if specifying shared custom attribute flag."
	datastoreCustomAttributePrefixSeparatorFlagHelp string = "Custom Attribute prefix separator specific to datastores. Skip if using Custom Attribute values as-is for comparison, otherwise optional if specifying shared custom attribute prefix separator, or using the default separator."
	sharedCustomAttributeNameFlagHelp               string = "Custom Attribute name for host ESXi systems and datastores. Optional if specifying resource-specific custom attribute names."
	sharedCustomAttributePrefixSeparatorFlagHelp    string = "Custom Attribute prefix separator for host ESXi systems and datastores. Skip if using Custom Attribute values as-is for comparison, otherwise optional if specifying resource-specific custom attribute prefix separator, or using the default separator."
	ignoreMissingCustomAttributeFlagHelp            string = "Toggles how missing specified Custom Attributes will be handled. By default, ESXi hosts and datastores missing the Custom Attribute are treated as an error condition."
	ignoreDatastoreFlagHelp                         string = "Specifies a comma-separated list of Datastore names that should be ignored or excluded from evaluation."
	datastoreNameFlagHelp                           string = "Datastore name as it is found within the vSphere inventory."
	datastoreUsageCriticalFlagHelp                  string = "Specifies the percentage of a datastore's storage usage (as a whole number) when a CRITICAL threshold is reached."
	datastoreUsageWarningFlagHelp                   string = "Specifies the percentage of a datastore's storage usage (as a whole number) when a WARNING threshold is reached."
	datacenterNameFlagHelp                          string = "Specifies the name of a vSphere Datacenter. If not specified, applicable plugins will attempt to use the default datacenter found in the vSphere environment. Not applicable to standalone ESXi hosts."
	snapshotsAgeCriticalFlagHelp                    string = "Specifies the age of a snapshot in days when a CRITICAL threshold is reached."
	snapshotsAgeWarningFlagHelp                     string = "Specifies the age of a snapshot in days when a WARNING threshold is reached."
	snapshotsSizeCriticalFlagHelp                   string = "Specifies the cumulative size in GB of all snapshots for a Virtual Machine when a CRITICAL threshold is reached."
	snapshotsSizeWarningFlagHelp                    string = "Specifies the cumulative size in GB of all snapshots for a Virtual Machine when a WARNING threshold is reached."
)

// Default flag settings if not overridden by user input
const (
	defaultLogLevel                     string = "info"
	defaultServer                       string = ""
	defaultTrustCert                    bool   = false
	defaultUsername                     string = ""
	defaultPassword                     string = ""
	defaultUserDomain                   string = ""
	defaultPort                         int    = 443
	defaultBranding                     bool   = false
	defaultDisplayVersionAndExit        bool   = false
	defaultPoweredOff                   bool   = false
	defaultVCPUsAllocatedCritical       int    = 100
	defaultVCPUsAllocatedWarning        int    = 95
	defaultIgnoreMissingCustomAttribute bool   = false
	defaultDatastoreName                string = ""
	defaultDatastoreUsageCritical       int    = 95
	defaultDatastoreUsageWarning        int    = 90
	defaultDatacenterName               string = ""
	defaultSnapshotsAgeCritical         int    = 2
	defaultSnapshotsAgeWarning          int    = 1
	defaultSnapshotsSizeCritical        int    = 40 // size in GB
	defaultSnapshotsSizeWarning         int    = 20 // size in GB

	// Intentionally set low to trigger validation failure if not specified by
	// the end user.
	defaultVCPUsMaxAllowedAllowed int = 0

	// Default timeout (in seconds) used when connecting to a remote server
	defaultConnectTimeout int = 10

	defaultCustomAttributeName string = ""

	// Default separator for Custom Attribute values.
	//
	// When specified, this separator can be used to split Custom Attribute
	// values in order to get at the prefix for comparison with Custom
	// Attributes used for other vSphere object types.
	//
	// For example, ESXi hosts might have a "Location" field that specifies
	// the datacenter and rack details (with a separator between them),
	// whereas a datastore might have only the datacenter as its Location
	// field value.
	//
	// defaultCustomAttributePrefixSeparator string = "-"
	//
	// By not specifying a default separator value, when an attempt is made to
	// split on the separator the full string will be returned as-is. This
	// forces the user to provide an actual prefix separator to enable prefix
	// splitting.
	defaultCustomAttributePrefixSeparator string = ""
)

// Plugin types provided by this project.
const (
	PluginTypeTools                    string = "vmware-tools"
	PluginTypeSnapshotsAge             string = "snapshots-age"
	PluginTypeSnapshotsSize            string = "snapshots-size"
	PluginTypeDatastoresSize           string = "datastore-size"
	PluginTypeResourcePoolsMemory      string = "resource-pools-memory"
	PluginTypeVirtualCPUsAllocation    string = "virtual-cpus-allocation"
	PluginTypeHostDatastoreVMsPairings string = "host-to-ds-to-vms"
)
