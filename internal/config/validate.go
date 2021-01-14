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

	case pluginType.SnapshotsAge:

	// 	if c.AgeWarning < 0 {
	// 		return fmt.Errorf(
	// 			"invalid cert expiration WARNING threshold number: %d",
	// 			c.AgeWarning,
	// 		)
	// 	}
	//
	// 	if c.AgeCritical < 0 {
	// 		return fmt.Errorf(
	// 			"invalid cert expiration CRITICAL threshold number: %d",
	// 			c.AgeCritical,
	// 		)
	// 	}
	//
	// 	if c.AgeCritical > c.AgeWarning {
	// 		return fmt.Errorf(
	// 			"critical threshold set higher than warning threshold",
	// 		)
	// 	}

	case pluginType.SnapshotsSize:

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

		if c.DatastoreUsageCritical < c.DatastoreUsageWarning {
			return fmt.Errorf(
				"datastore critical threshold set lower than warning threshold",
			)
		}

	case pluginType.ResourcePoolsMemory:

	case pluginType.VirtualCPUsAllocation:

		if c.VCPUsMaxAllowed < 1 {
			return fmt.Errorf(
				"invalid maximum number of vCPUs allowed: %d",
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

		if c.VCPUsAllocatedCritical < c.VCPUsAllocatedWarning {
			return fmt.Errorf(
				"vCPUs allocation critical threshold set lower than warning threshold",
			)
		}

	case pluginType.Host2Datastores2VMs:

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

	}

	if c.Server == "" {
		return fmt.Errorf("server FQDN or IP Address not provided")
	}

	if c.Username == "" {
		return fmt.Errorf("username not provided")
	}

	if c.Password == "" {
		return fmt.Errorf("password not provided")
	}

	// only one of these options may be used
	if len(c.ExcludedResourcePools) > 0 && len(c.IncludedResourcePools) > 0 {
		return fmt.Errorf(
			"only one of %q or %q flags may be specified",
			"include-rp",
			"exclude-rp",
		)
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
