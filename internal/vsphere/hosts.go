// Copyright 2021 Adam Chalkley
//
// https://github.com/atc0005/check-vmware
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package vsphere

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/vmware/govmomi/units"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

// GetHostSystems accepts a context, a connected client and a boolean value
// indicating whether a subset of properties per HostSystem are retrieved. A
// collection of HostSystems with requested properties is returned. If
// requested, a subset of all available properties will be retrieved (faster)
// instead of recursively fetching all properties (about 2x as slow).
func GetHostSystems(ctx context.Context, c *vim25.Client, propsSubset bool) ([]mo.HostSystem, error) {

	funcTimeStart := time.Now()

	// declare this early so that we can grab a pointer to it in order to
	// access the entries later
	var hss []mo.HostSystem

	defer func(hss *[]mo.HostSystem) {
		logger.Printf(
			"It took %v to execute GetHostSystems func (and retrieve %d HostSystems).\n",
			time.Since(funcTimeStart),
			len(*hss),
		)
	}(&hss)

	err := getObjects(ctx, c, &hss, c.ServiceContent.RootFolder, propsSubset)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve HostSystems: %w", err)
	}

	sort.Slice(hss, func(i, j int) bool {
		return strings.ToLower(hss[i].Name) < strings.ToLower(hss[j].Name)
	})

	return hss, nil
}

// GetHostSystemByName accepts the name of a HostSystem, the name of a
// datacenter and a boolean value indicating whether only a subset of
// properties for the HostSystem should be returned. If requested, a subset of
// all available properties will be retrieved (faster) instead of recursively
// fetching all properties (about 2x as slow). If the datacenter name is an
// empty string then the default datacenter will be used.
func GetHostSystemByName(ctx context.Context, c *vim25.Client, hsName string, datacenter string, propsSubset bool) (mo.HostSystem, error) {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute GetHostSystemByName func.\n",
			time.Since(funcTimeStart),
		)
	}()

	var hostSystem mo.HostSystem
	err := getObjectByName(ctx, c, &hostSystem, hsName, datacenter, propsSubset)

	if err != nil {
		return mo.HostSystem{}, err
	}

	return hostSystem, nil

}

// FilterHostSystemByName accepts a collection of HostSystems and a HostSystem
// name to filter against. An error is returned if the list of HostSystems is
// empty or if a match was not found.
func FilterHostSystemByName(hss []mo.HostSystem, hsName string) (mo.HostSystem, error) {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute FilterHostSystemByName func.\n",
			time.Since(funcTimeStart),
		)
	}()

	if len(hss) == 0 {
		return mo.HostSystem{}, fmt.Errorf("received empty list of HostSystems to filter by name")
	}

	for _, hs := range hss {
		if hs.Name == hsName {
			return hs, nil
		}
	}

	return mo.HostSystem{}, fmt.Errorf(
		"error: failed to retrieve HostSystem using provided name %q",
		hsName,
	)

}

// FilterHostSystemByID receives a collection of HostSystems and a HostSystem ID
// to filter against. An error is returned if the list of HostSystems is empty
// or if a match was not found.
func FilterHostSystemByID(hss []mo.HostSystem, hsID string) (mo.HostSystem, error) {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute FilterHostSystemByID func.\n",
			time.Since(funcTimeStart),
		)
	}()

	if len(hss) == 0 {
		return mo.HostSystem{}, fmt.Errorf("received empty list of HostSystems to filter by ID")
	}

	for _, hs := range hss {
		// return match, if available

		if hs.Summary.Host.Value == hsID {
			return hs, nil
		}
	}

	return mo.HostSystem{}, fmt.Errorf(
		"error: failed to retrieve HostSystem using provided id %q",
		hsID,
	)

}

// GetHostSystemsTotalMemory returns the total memory capacity for all
// HostSystems. Unless requested, offline or otherwise unavailable hosts are
// included for evaluation based on the assumption that offline hosts are
// offline for only a brief time and should still be considered part of
// overall cluster capacity.
func GetHostSystemsTotalMemory(ctx context.Context, c *vim25.Client, excludeOffline bool) (int64, error) {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute GetHostSystemTotalMemory func.\n",
			time.Since(funcTimeStart),
		)
	}()

	clusterHosts, err := GetHostSystems(ctx, c, true)
	if err != nil {
		return 0, fmt.Errorf(
			"failed to gather total memory capacity for host systems: %w",
			err,
		)
	}

	var clusterMemory int64
	for _, host := range clusterHosts {

		// Evaluate offline systems by default, unless requested otherwise.
		if excludeOffline {

			logger.Printf("Checking host %s availability ... \n", host.Name)

			switch {

			case host.Runtime.PowerState == types.HostSystemPowerStatePoweredOn &&
				host.Runtime.ConnectionState == types.HostSystemConnectionStateConnected:
				// desired state, no other limiting factors detected

			case host.Runtime.InMaintenanceMode:
				logger.Printf("Host %s is in maintenance mode, skipping evaluation ...\n", host.Name)
				continue

			case host.Runtime.InQuarantineMode != nil && *host.Runtime.InQuarantineMode:
				logger.Printf("Host %s is in quarantine mode, skipping evaluation ...\n", host.Name)
				continue

			case host.Runtime.PowerState == types.HostSystemPowerStatePoweredOff:
				logger.Printf("Host %s is powered off, skipping evaluation ...\n", host.Name)
				continue

			case host.Runtime.PowerState == types.HostSystemPowerStateStandBy:
				logger.Printf("Host %s is in standby, skipping evaluation ...\n", host.Name)
				continue

			case host.Runtime.ConnectionState == types.HostSystemConnectionStateDisconnected:
				logger.Printf("Host %s is disconnected, skipping evaluation ...\n", host.Name)
				continue

			case host.Runtime.ConnectionState == types.HostSystemConnectionStateNotResponding:
				logger.Printf("Host %s is not responding, skipping evaluation ...\n", host.Name)
				continue

			default:
				logger.Printf("Host %s is in an UNKNOWN state, skipping evaluation ...\n", host.Name)
				continue

			}
		}

		logger.Printf(
			"Host %s has %s memory capacity.\n",
			host.Name,
			units.ByteSize(host.Hardware.MemorySize),
		)

		clusterMemory += host.Hardware.MemorySize
	}

	return clusterMemory, nil

}
