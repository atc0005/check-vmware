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

	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
)

// GetNetworks accepts a context, a connected client and a boolean value
// indicating whether a subset of properties per Network are retrieved. If
// requested, a subset of all available properties will be retrieved (faster)
// instead of recursively fetching all properties (about 2x as slow). A
// collection of Networks with requested properties is returned or nil and an
// error, if one occurs.
func GetNetworks(ctx context.Context, c *vim25.Client, propsSubset bool) ([]mo.Network, error) {

	funcTimeStart := time.Now()

	// declare this early so that we can grab a pointer to it in order to
	// access the entries later
	var nets []mo.Network

	defer func(nets *[]mo.Network) {
		logger.Printf(
			"It took %v to execute GetNetworks func (and retrieve %d Networks).\n",
			time.Since(funcTimeStart),
			len(*nets),
		)
	}(&nets)

	err := getObjects(ctx, c, &nets, c.ServiceContent.RootFolder, propsSubset)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve host systems: %w", err)
	}

	sort.Slice(nets, func(i, j int) bool {
		return strings.ToLower(nets[i].Name) < strings.ToLower(nets[j].Name)
	})

	return nets, nil
}

// GetNetworkByName accepts the name of a network, the name of a datacenter
// and a boolean value indicating whether only a subset of properties for the
// Network should be returned. If requested, a subset of all available
// properties will be retrieved (faster) instead of recursively fetching all
// properties (about 2x as slow). If the datacenter name is an empty string
// then the default datacenter will be used.
func GetNetworkByName(ctx context.Context, c *vim25.Client, netName string, datacenter string, propsSubset bool) (mo.Network, error) {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute GetNetworkByName func.\n",
			time.Since(funcTimeStart),
		)
	}()

	var network mo.Network
	err := getObjectByName(ctx, c, &network, netName, datacenter, propsSubset)

	if err != nil {
		return mo.Network{}, err
	}

	return network, nil

}

// FilterNetworkByName accepts a collection of Networks and a Network name to
// filter against. An error is returned if the list of Networks is empty or if
// a match was not found. The matching Network is returned along with the
// number of Networks that were excluded.
func FilterNetworkByName(nets []mo.Network, netName string) (mo.Network, int, error) {

	funcTimeStart := time.Now()

	// If error condition, no exclusions are made
	numExcluded := 0

	defer func() {
		logger.Printf(
			"It took %v to execute FilterNetworkByName func.\n",
			time.Since(funcTimeStart),
		)
	}()

	if len(nets) == 0 {
		return mo.Network{}, numExcluded, fmt.Errorf("received empty list of networks to filter by name")
	}

	for _, net := range nets {
		if net.Name == netName {
			// we are excluding everything but the single name value match
			numExcluded = len(nets) - 1
			return net, numExcluded, nil
		}
	}

	return mo.Network{}, numExcluded, fmt.Errorf(
		"error: failed to retrieve Network using provided name %q",
		netName,
	)

}

// FilterNetworkByID receives a collection of Networks and a Network ID to
// filter against. An error is returned if the list of Networks is empty or if
// a match was not found. The matching Network is returned along with the
// number of Networks that were excluded.
func FilterNetworkByID(nets []mo.Network, netID string) (mo.Network, int, error) {

	funcTimeStart := time.Now()

	// If error condition, no exclusions are made
	numExcluded := 0

	defer func() {
		logger.Printf(
			"It took %v to execute FilterNetworkByID func.\n",
			time.Since(funcTimeStart),
		)
	}()

	if len(nets) == 0 {
		return mo.Network{}, numExcluded, fmt.Errorf("received empty list of networks to filter by ID")
	}

	for _, net := range nets {
		// return match, if available
		if net.Self.Value == netID {
			// we are excluding everything but the single ID value match
			numExcluded = len(nets) - 1
			return net, numExcluded, nil
		}
	}

	return mo.Network{}, numExcluded, fmt.Errorf(
		"error: failed to retrieve Network using provided id %q",
		netID,
	)
}
