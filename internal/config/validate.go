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
				"critical threshold set lower than warning threshold",
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
