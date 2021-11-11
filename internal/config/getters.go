// Copyright 2021 Adam Chalkley
//
// https://github.com/atc0005/check-vmware
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package config

import (
	"fmt"
	"time"
)

// Timeout converts the user-specified plugin runtime/execution timeout value
// in seconds to an appropriate time duration value for use with setting
// context deadline value.
func (c Config) Timeout() time.Duration {
	return time.Duration(c.timeout) * time.Second
}

// add getters to indicate whether user has specified a shared custom
// attribute or whether separate host and datastore attributes are used.

// UsingSharedCA indicates whether the user opted to use a shared Custom
// Attribute for linking Hosts with specific Datastores, or whether they opted
// to instead specify Custom Attributes for Hosts and Datastores separately.
// This method relies heavily on config validation to enforce only one flag
// set for specifying the required Custom Attribute.
func (c Config) UsingSharedCA() bool {
	return c.sharedCustomAttributeName != ""
}

// UsingCAPrefixes indicates whether the user opted to use a Custom Attribute
// prefix in place of a literal value for linking Hosts with specific
// Datastores.
func (c Config) UsingCAPrefixes() bool {
	return c.sharedCustomAttributePrefixSeparator != "" ||
		(c.datastoreCustomAttributePrefixSeparator != "" && c.hostCustomAttributePrefixSeparator != "")
}

// DatastoreCAName returns the user-provided Custom Attribute name specific to
// datastores or the shared Custom Attribute name used for both datastores and
// hosts.
func (c Config) DatastoreCAName() string {
	if c.datastoreCustomAttributeName != "" {
		return c.datastoreCustomAttributeName
	}

	return c.sharedCustomAttributeName
}

// HostCAName returns the user-provided Custom Attribute name specific to
// hosts or the shared Custom Attribute name used for both hosts and
// datastores.
func (c Config) HostCAName() string {
	if c.hostCustomAttributeName != "" {
		return c.hostCustomAttributeName
	}

	return c.sharedCustomAttributeName
}

// DatastoreCASep returns the user-provided Custom Attribute prefix separator
// specific to datastores, the shared separator provided for both datastores
// and hosts or the default separator if not specified.
func (c Config) DatastoreCASep() string {
	switch {
	case c.datastoreCustomAttributePrefixSeparator != "":
		return c.datastoreCustomAttributePrefixSeparator
	case c.sharedCustomAttributePrefixSeparator != "":
		return c.sharedCustomAttributePrefixSeparator
	default:
		return defaultCustomAttributePrefixSeparator
	}
}

// HostCASep returns the user-provided Custom Attribute prefix separator
// specific to datastores, the shared separator provided for both datastores
// and hosts or the default separator if not specified.
func (c Config) HostCASep() string {
	switch {
	case c.hostCustomAttributePrefixSeparator != "":
		return c.hostCustomAttributePrefixSeparator
	case c.sharedCustomAttributePrefixSeparator != "":
		return c.sharedCustomAttributePrefixSeparator
	default:
		return defaultCustomAttributePrefixSeparator
	}
}

// VirtualHardwareApplyMinVersionCheck indicates whether all virtual machines
// are required to have the specified minimum hardware version or greater.
// This is only used if the other behaviors were not requested.
func (c Config) VirtualHardwareApplyMinVersionCheck() bool {

	return c.VirtualHardwareMinimumVersion != defaultVirtualHardwareMinimumVersion &&
		c.VirtualHardwareOutdatedByCritical == defaultVirtualHardwareOutdatedByCritical &&
		c.VirtualHardwareOutdatedByWarning == defaultVirtualHardwareOutdatedByWarning &&
		!c.VirtualHardwareDefaultVersionIsMinimum

}

// VirtualHardwareApplyDefaultIsMinVersionCheck indicates whether all virtual
// machines are required to have the host or cluster default hardware version
// or greater. This is only used if the other behaviors were not requested.
func (c Config) VirtualHardwareApplyDefaultIsMinVersionCheck() bool {

	return c.VirtualHardwareMinimumVersion == defaultVirtualHardwareMinimumVersion &&
		c.VirtualHardwareOutdatedByCritical == defaultVirtualHardwareOutdatedByCritical &&
		c.VirtualHardwareOutdatedByWarning == defaultVirtualHardwareOutdatedByWarning &&
		c.VirtualHardwareDefaultVersionIsMinimum

}

// VirtualHardwareApplyOutdatedByVersionCheck indicates whether the outdated
// by CRITICAL and WARNING threshold checks are applied to assert that virtual
// hardware versions are within the specified thresholds. This is only used if
// the other behaviors were not requested.
func (c Config) VirtualHardwareApplyOutdatedByVersionCheck() bool {

	return c.VirtualHardwareMinimumVersion == defaultVirtualHardwareMinimumVersion &&
		(c.VirtualHardwareOutdatedByCritical != defaultVirtualHardwareOutdatedByCritical ||
			c.VirtualHardwareOutdatedByWarning != defaultVirtualHardwareOutdatedByWarning) &&
		!c.VirtualHardwareDefaultVersionIsMinimum

}

// VirtualHardwareApplyHomogeneousVersionCheck indicates whether the default
// behavior of asserting that all virtual hardware versions are the same is
// applied. This is only used if the other behaviors were not requested.
func (c Config) VirtualHardwareApplyHomogeneousVersionCheck() bool {

	return c.VirtualHardwareMinimumVersion == defaultVirtualHardwareMinimumVersion &&
		c.VirtualHardwareOutdatedByCritical == defaultVirtualHardwareOutdatedByCritical &&
		c.VirtualHardwareOutdatedByWarning == defaultVirtualHardwareOutdatedByWarning &&
		!c.VirtualHardwareDefaultVersionIsMinimum

}

// UserAgent returns a string usable as-is as a custom user agent for plugins
// provided by this project.
func (c Config) UserAgent() string {

	// Default User Agent: (Go-http-client/1.1)
	// https://datatracker.ietf.org/doc/html/draft-ietf-httpbis-p2-semantics-22#section-5.5.3
	return fmt.Sprintf(
		"%s/%s",
		c.App.Name,
		c.App.Version,
	)

}

// DatastorePerfThresholds returns Datastore Performance Summary latency
// thresholds for the default percentile. If defined by the user, those values
// are returned. If the user did not specify individual threshold values,
// default values are returned.
func (c Config) DatastorePerfThresholds() DSPerformanceSummaryThresholds {

	readLatencyWarning := defaultDatastoreReadLatencyWarning
	if c.datastoreReadLatencyWarning.isSet {
		readLatencyWarning = c.datastoreReadLatencyWarning.value
	}

	readLatencyCritical := defaultDatastoreReadLatencyCritical
	if c.datastoreReadLatencyCritical.isSet {
		readLatencyCritical = c.datastoreReadLatencyCritical.value
	}

	writeLatencyWarning := defaultDatastoreWriteLatencyWarning
	if c.datastoreWriteLatencyWarning.isSet {
		writeLatencyWarning = c.datastoreWriteLatencyWarning.value
	}

	writeLatencyCritical := defaultDatastoreWriteLatencyCritical
	if c.datastoreWriteLatencyCritical.isSet {
		writeLatencyCritical = c.datastoreWriteLatencyCritical.value
	}

	vmLatencyWarning := defaultDatastoreVMLatencyWarning
	if c.datastoreVMLatencyWarning.isSet {
		vmLatencyWarning = c.datastoreVMLatencyWarning.value
	}

	vmLatencyCritical := defaultDatastoreVMLatencyCritical
	if c.datastoreVMLatencyCritical.isSet {
		vmLatencyCritical = c.datastoreVMLatencyCritical.value
	}

	return DSPerformanceSummaryThresholds{
		ReadLatencyWarning:   readLatencyWarning,
		ReadLatencyCritical:  readLatencyCritical,
		WriteLatencyWarning:  writeLatencyWarning,
		WriteLatencyCritical: writeLatencyCritical,
		VMLatencyWarning:     vmLatencyWarning,
		VMLatencyCritical:    vmLatencyCritical,
	}

}

// DatastorePerfPercentileSet returns a mapping of Datastore Performance
// Summary percentile to latency thresholds. This mapping is used to evaluate
// Datastore Performance metrics to determine overall plugin state.
func (c Config) DatastorePerfPercentileSet() MultiValueDSPerfPercentileSetFlag {

	switch {

	// A percentile set wasn't explicitly defined, so use the
	// explicit/specified individual latency metric flag values (or default as
	// set by the flags package) to build a default percentile set.
	case len(c.datastorePerformancePercentileSet) == 0:

		return MultiValueDSPerfPercentileSetFlag{
			defaultDatastorePerfSumPercentile: c.DatastorePerfThresholds(),
		}

	default:

		// At least one percentile set *was* defined, so return the set for
		// evaluation. Ignore any thresholds set by individual latency metric
		// flags.
		return c.datastorePerformancePercentileSet

	}

}
