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
	vmIncludedResourcePoolsFlagHelp                 string = "Specifies a comma-separated list of Resource Pools that should be exclusively used when evaluating VMs. Specifying this option will also exclude any VMs from evaluation that are *outside* of a Resource Pool. This option is incompatible with specifying a list of Resource Pools to ignore or exclude from evaluation."
	vmExcludedResourcePoolsFlagHelp                 string = "Specifies a comma-separated list of Resource Pools that should be ignored when evaluating VMs. This option is incompatible with specifying a list of Resource Pools to include for evaluation."
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
	datacenterNamesFlagHelp                         string = "Specifies the name of one or more vSphere Datacenters. If not specified, applicable plugins will attempt to evaluate all visible datacenters found in the vSphere environment. Not applicable to standalone ESXi hosts."
	clusterNameFlagHelp                             string = "Specifies the name of a vSphere Cluster. If not specified, applicable plugins will attempt to use the default cluster found in the vSphere environment. Not applicable to standalone ESXi hosts."
	snapshotsAgeCriticalFlagHelp                    string = "Specifies the age of a snapshot in days when a CRITICAL threshold is reached."
	snapshotsAgeWarningFlagHelp                     string = "Specifies the age of a snapshot in days when a WARNING threshold is reached."
	snapshotsCountCriticalFlagHelp                  string = "Specifies the number of snapshots per VM when a CRITICAL threshold is reached."
	snapshotsCountWarningFlagHelp                   string = "Specifies the number of snapshots per VM when a WARNING threshold is reached."
	snapshotsSizeCriticalFlagHelp                   string = "Specifies the cumulative size in GB of all snapshots for a Virtual Machine when a CRITICAL threshold is reached."
	snapshotsSizeWarningFlagHelp                    string = "Specifies the cumulative size in GB of all snapshots for a Virtual Machine when a WARNING threshold is reached."
	resourcePoolsMemoryMaxAllowedFlagHelp           string = "Specifies the maximum amount of memory that we are allowed to consume in GB (as a whole number) in the target VMware environment across all specified Resource Pools. VMs that are running outside of resource pools are not considered in these calculations."
	resourcePoolsMemoryUseCriticalFlagHelp          string = "Specifies the percentage of memory use (as a whole number) across all specified Resource Pools when a CRITICAL threshold is reached."
	resourcePoolsMemoryUseWarningFlagHelp           string = "Specifies the percentage of memory use (as a whole number) across all specified Resource Pools when a WARNING threshold is reached."
	hostSystemMemoryUseCriticalFlagHelp             string = "Specifies the percentage of memory use (as a whole number) when a CRITICAL threshold is reached."
	hostSystemMemoryUseWarningFlagHelp              string = "Specifies the percentage of memory use (as a whole number) when a WARNING threshold is reached."
	hostSystemNameFlagHelp                          string = "ESXi host/server name as it is found within the vSphere inventory."
	hostSystemCPUUseCriticalFlagHelp                string = "Specifies the percentage of CPU use (as a whole number) when a CRITICAL threshold is reached."
	hostSystemCPUUseWarningFlagHelp                 string = "Specifies the percentage of CPU use (as a whole number) when a WARNING threshold is reached."
	vmPowerCycleUptimeCriticalFlagHelp              string = "Specifies the power cycle (off/on) uptime in days per VM when a CRITICAL threshold is reached."
	vmPowerCycleUptimeWarningFlagHelp               string = "Specifies the power cycle (off/on) uptime in days per VM when a WARNING threshold is reached."
	virtualHardwareOutdatedByCriticalFlagHelp       string = "If provided, this value is the CRITICAL threshold for outdated virtual hardware versions. If the current virtual hardware version for a VM is found to be more than this many versions older than the latest version a CRITICAL state is triggered. Required if specifying the WARNING threshold for outdated virtual hardware versions."
	virtualHardwareOutdatedByWarningFlagHelp        string = "If provided, this value is the WARNING threshold for outdated virtual hardware versions. If the current virtual hardware version for a VM is found to be more than this many versions older than the latest version a WARNING state is triggered. Required if specifying the CRITICAL threshold for outdated virtual hardware versions."
	virtualHardwareMinimumVersionFlagHelp           string = "If provided, this value is the minimum virtual hardware version accepted for each Virtual Machine. Any Virtual Machine not meeting this minimum value is considered to be in a CRITICAL state. Per KB 1003746, version 3 appears to be the oldest version supported."
	virtualHardwareDefaultIsMinimumFlagHelp         string = "If specified, the host or cluster default virtual hardware version is the minimum hardware version allowed. Any Virtual Machine not meeting this minimum value is considered to be in a WARNING state."
	includedAlarmEntityTypesFlagHelp                string = "If specified, triggered alarms will only be evaluated if the associated entity type (e.g., Datastore) matches one of the provided values."
	excludedAlarmEntityTypesFlagHelp                string = "If specified, triggered alarms will only be evaluated if the associated entity type (e.g., Datastore) does NOT match one of the provided values."
	includedAlarmEntityNamesFlagHelp                string = "If specified, triggered alarms will only be evaluated if the associated entity name (e.g., \"node1.example.com\") matches one of the provided values."
	excludedAlarmEntityNamesFlagHelp                string = "If specified, triggered alarms will only be evaluated if the associated entity name (e.g., \"node1.example.com\") does NOT match one of the provided values."
	evaluateAcknowledgedTriggeredAlarmFlagHelp      string = "Toggles evaluation of acknowledged triggered alarms in addition to unacknowledged triggered alarms. Evaluation of acknowledged alarms is disabled by default."
	includedAlarmNamesFlagHelp                      string = "If specified, triggered alarms will only be evaluated if the alarm name (e.g., \"Datastore usage on disk\") case-insensitively matches one of the provided substring values (e.g., \"datastore\" or \"datastore usage\") and is not explicitly excluded by another filter in the pipeline; while multiple explicit inclusions are allowed, explicit exclusions have precedence over explicit inclusions and will exclude the triggered alarm from further evaluation."
	excludedAlarmNamesFlagHelp                      string = "If specified, triggered alarms will only be evaluated if the alarm name (e.g., \"Datastore usage on disk\") DOES NOT case-insensitively match one of the provided substring values (e.g., \"datastore\" or \"datastore usage\") and is not explicitly excluded by another filter in the pipeline; while multiple explicit inclusions are allowed, explicit exclusions have precedence over explicit inclusions and will exclude the triggered alarm from further evaluation."
	includedAlarmDescriptionsFlagHelp               string = "If specified, triggered alarms will only be evaluated if the alarm description (e.g., \"Default alarm to monitor datastore disk usage\") case-insensitively matches one of the provided substring values (e.g., \"datastore disk\" or \"monitor datastore\") and is not explicitly excluded by another filter in the pipeline; while multiple explicit inclusions are allowed, explicit exclusions have precedence over explicit inclusions and will exclude the triggered alarm from further evaluation."
	excludedAlarmDescriptionsFlagHelp               string = "If specified, triggered alarms will only be evaluated if the alarm description (e.g., \"Default alarm to monitor datastore disk usage\") DOES NOT case-insensitively match one of the provided substring values (e.g., \"datastore disk\" or \"monitor datastore\") and is not explicitly excluded by another filter in the pipeline; while multiple explicit inclusions are allowed, explicit exclusions have precedence over explicit inclusions and will exclude the triggered alarm from further evaluation."
	includedAlarmStatusesFlagHelp                   string = "If specified, triggered alarms will only be evaluated if the alarm status (e.g., \"yellow\") case-insensitively matches one of the provided keywords (e.g., \"yellow\") and is not explicitly excluded by another filter in the pipeline; while multiple explicit inclusions are allowed, explicit exclusions have precedence over explicit inclusions and will exclude the triggered alarm from further evaluation."
	excludedAlarmStatusesFlagHelp                   string = "If specified, triggered alarms will only be evaluated if the alarm status (e.g., \"yellow\") DOES NOT case-insensitively match one of the provided keywords (e.g., \"yellow\") and is not explicitly excluded by another filter in the pipeline; while multiple explicit inclusions are allowed, explicit exclusions have precedence over explicit inclusions and will exclude the triggered alarm from further evaluation."
	includedAlarmEntityResourcePoolsFlagHelp        string = "If specified, triggered alarms will only be evaluated if the associated entity is part of one of the provided Resource Pools (case-insensitive match on the name) and is not explicitly excluded by another filter in the pipeline; while multiple explicit inclusions are allowed, explicit exclusions have precedence over explicit inclusions and will exclude the triggered alarm from further evaluation."
	excludedAlarmEntityResourcePoolsFlagHelp        string = "If specified, triggered alarms will only be evaluated if the associated entity is NOT part of one of the provided Resource Pools (case-insensitive match on the name) and is not explicitly excluded by another filter in the pipeline; while multiple explicit inclusions are allowed, explicit exclusions have precedence over explicit inclusions and will exclude the triggered alarm from further evaluation."
)

// Default flag settings if not overridden by user input
const (
	defaultLogLevel                     string = "info"
	defaultServer                       string = ""
	defaultTrustCert                    bool   = false
	defaultUsername                     string = ""
	defaultPassword                     string = ""
	defaultUserDomain                   string = ""
	defaultClusterName                  string = ""
	defaultPort                         int    = 443
	defaultBranding                     bool   = false
	defaultDisplayVersionAndExit        bool   = false
	defaultPoweredOff                   bool   = false
	defaultEvaluateAcknowledgedAlarms   bool   = false
	defaultVCPUsAllocatedCritical       int    = 100
	defaultVCPUsAllocatedWarning        int    = 95
	defaultIgnoreMissingCustomAttribute bool   = false
	defaultDatastoreName                string = ""
	defaultDatastoreUsageCritical       int    = 95
	defaultDatastoreUsageWarning        int    = 90
	defaultDatacenterName               string = ""
	defaultSnapshotsAgeCritical         int    = 2
	defaultSnapshotsAgeWarning          int    = 1
	defaultSnapshotsCountCritical       int    = 25 // max is 32
	defaultSnapshotsCountWarning        int    = 4  // recommended cap is 3-4
	defaultSnapshotsSizeCritical        int    = 40 // size in GB
	defaultSnapshotsSizeWarning         int    = 20 // size in GB
	defaultHostSystemName               string = ""
	defaultVMPowerCycleUptimeCritical   int    = 90
	defaultVMPowerCycleUptimeWarning    int    = 60

	// The default values are intentionally invalid to help determine whether
	// the user has supplied values for the flags.
	defaultVirtualHardwareMinimumVersion     int = -1
	defaultVirtualHardwareOutdatedByWarning  int = -1
	defaultVirtualHardwareOutdatedByCritical int = -1

	// Whether the default host or cluster hardware version is the minimum
	// version allowed
	defaultVirtualHardwareDefaultIsMinimum bool = false

	// default memory usage values for Resource Pools and ESXi Host systems
	defaultMemoryUseCritical int = 95
	defaultMemoryUseWarning  int = 80

	// HostSystem CPU usage thresholds
	defaultCPUUseCritical int = 95
	defaultCPUUseWarning  int = 80

	// Intentionally set low to trigger validation failure if not specified by
	// the end user.
	defaultVCPUsMaxAllowed               int = 0
	defaultResourcePoolsMemoryMaxAllowed int = 0

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
	PluginTypeTools                          string = "vmware-tools"
	PluginTypeSnapshotsAge                   string = "snapshots-age"
	PluginTypeSnapshotsCount                 string = "snapshots-count"
	PluginTypeSnapshotsSize                  string = "snapshots-size"
	PluginTypeDatastoresSize                 string = "datastore-size"
	PluginTypeResourcePoolsMemory            string = "resource-pools-memory"
	PluginTypeVirtualCPUsAllocation          string = "virtual-cpus-allocation"
	PluginTypeVirtualHardwareVersion         string = "virtual-hardware-version"
	PluginTypeHostDatastoreVMsPairings       string = "host-to-ds-to-vms"
	PluginTypeHostSystemMemory               string = "host-system-memory"
	PluginTypeHostSystemCPU                  string = "host-system-cpu"
	PluginTypeVirtualMachinePowerCycleUptime string = "vm-power-uptime"
	PluginTypeDiskConsolidation              string = "disk-consolidation"
	PluginTypeInteractiveQuestion            string = "interactive-question"
	PluginTypeAlarms                         string = "alarms"
)

// Known limits
// https://trainingrevolution.wordpress.com/2018/07/22/vmware-vsphere-6-7-character-limits-for-objects/
const (
	MaxClusterNameChars int = 80
)

// ThresholdNotUsed indicates that a plugin is not using a specific threshold.
// This is visible in locations where Long Service Output text is displayed.
const ThresholdNotUsed string = "Not used."

const (

	// LogLevelDisabled maps to zerolog.Disabled logging level
	LogLevelDisabled string = "disabled"

	// LogLevelPanic maps to zerolog.PanicLevel logging level
	LogLevelPanic string = "panic"

	// LogLevelFatal maps to zerolog.FatalLevel logging level
	LogLevelFatal string = "fatal"

	// LogLevelError maps to zerolog.ErrorLevel logging level
	LogLevelError string = "error"

	// LogLevelWarn maps to zerolog.WarnLevel logging level
	LogLevelWarn string = "warn"

	// LogLevelInfo maps to zerolog.InfoLevel logging level
	LogLevelInfo string = "info"

	// LogLevelDebug maps to zerolog.DebugLevel logging level
	LogLevelDebug string = "debug"

	// LogLevelTrace maps to zerolog.TraceLevel logging level
	LogLevelTrace string = "trace"
)

// Valid Triggered Alarm status keywords. Provided by sysadmin, maps to
// ManagedEntityStatus values.
const (

	// native vSphere keywords
	AlarmStatusRed    string = "red"
	AlarmStatusYellow string = "yellow"
	AlarmStatusGreen  string = "green"
	AlarmStatusGray   string = "gray"

	// Nagios keywords, though these values are displayed within the web UI
	AlarmStatusCritical string = "critical"
	AlarmStatusWarning  string = "warning"
	AlarmStatusOk       string = "ok"
	AlarmStatusUnknown  string = "unknown"
)
