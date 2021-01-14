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
	"strings"

	"github.com/rs/zerolog"
)

// Updated via Makefile builds. Setting placeholder value here so that
// something resembling a version string will be provided for non-Makefile
// builds.
var version string = "x.y.z"

// ErrVersionRequested indicates that the user requested application version
// information.
var ErrVersionRequested = errors.New("version information requested")

// PluginType represents the type of plugin that is being
// configured/initialized. Not all plugin types will use the same features and
// as a result will not accept the same flags. Unless noted otherwise, each of
// the plugin types are incompatible with each other, though some flags are
// common to all types.
type PluginType struct {
	Tools                  bool
	SnapshotsAge           bool
	SnapshotsSize          bool
	DatastoresSize         bool
	ResourcePoolsMemory    bool
	VirtualCPUsAllocation  bool
	VirtualHardwareVersion bool
	Host2Datastores2VMs    bool

	// TODO:
	// - Hosts (memory, CPU usage)
	// - vCenter/server time (NTP)

}

// multiValueStringFlag is a custom type that satisfies the flag.Value
// interface in order to accept multiple string values for some of our flags.
type multiValueStringFlag []string

// String returns a comma separated string consisting of all slice elements.
func (mvs *multiValueStringFlag) String() string {

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
	// inventory. Not applicable to standline ESXi hosts.
	DatacenterName string

	// Port is the TCP port used by the certifcate-enabled service.
	Port int

	// AgeWarning is the number of days remaining before certificate
	// expiration when this application will flag the NotAfter certificate
	// field as a WARNING state.
	// AgeWarning int

	// AgeCritical is the number of days remaining before certificate
	// expiration when this application will flag the NotAfter certificate
	// field as a CRITICAL state.
	// AgeCritical int

	// timeout is the number of seconds allowed before the connection attempt
	// to a standalone ESXi host or vCenter instance is abandoned and an error
	// returned.
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

	// DatastoreUsageWarning specifies the percentage of a datastore's storage
	// usage (as a whole number) when a WARNING threshold is reached.
	DatastoreUsageWarning int

	// DatastoreUsageCritical specifies the percentage of a datastore's storage
	// usage (as a whole number) when a CRITICAL threshold is reached.
	DatastoreUsageCritical int

	// IgnoreMissingCustomAttribute indicates whether a host or datastore
	// missing the specified Custom Attribute should be ignored.
	IgnoreMissingCustomAttribute bool

	// PoweredOff indicates whether powered off VMs are evaluated in addition
	// to powered on VMs.
	PoweredOff bool

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

	// IncludedResourcePools lists resource pools that are explicitly
	// monitored.
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

	// Log is an embedded zerolog Logger initialized via config.New().
	Log zerolog.Logger
}

// Usage is a custom override for the default Help text provided by the flag
// package. Here we prepend some additional metadata to the existing output.
var Usage = func() {
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

// New is a factory function that produces a new Config object based on user
// provided flag and config file values. It is responsible for validating
// user-provided values and initializing the logging settings used by this
// application.
func New(pluginType PluginType) (*Config, error) {
	var config Config

	config.handleFlagsConfig(pluginType)

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

	return &config, nil

}
