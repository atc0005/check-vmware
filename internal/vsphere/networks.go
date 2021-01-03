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
	"os"
	"sort"
	"strings"
	"time"

	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/property"
	"github.com/vmware/govmomi/view"
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
		fmt.Fprintf(
			os.Stderr,
			"It took %v to execute GetNetworks func (and retrieve %d Networks).\n",
			time.Since(funcTimeStart),
			len(*nets),
		)
	}(&nets)

	// Create a view of Network objects
	m := view.NewManager(c)

	v, err := m.CreateContainerView(
		ctx,
		c.ServiceContent.RootFolder,
		[]string{
			"Network", // managed object types we are exposing via this view
		},
		true,
	)
	if err != nil {
		return nil, err
	}

	defer func() {
		// Per vSphere Web Services SDK Programming Guide - VMware vSphere 7.0
		// Update 1:
		//
		// A best practice when using views is to call the DestroyView()
		// method when a view is no longer needed. This practice frees memory
		// on the server.
		if err := v.Destroy(ctx); err != nil {
			fmt.Println("Error occurred while destroying view")
		}
	}()

	var props []string
	if propsSubset {
		// https://code.vmware.com/apis/1067/vsphere
		// https://vdc-download.vmware.com/vmwb-repository/dcr-public/a5f4000f-1ea8-48a9-9221-586adff3c557/7ff50256-2cf2-45ea-aacd-87d231ab1ac7/vim.Network.html
		props = []string{
			"summary", // properties of this network
			"name",    // name of this network
			"host",    // hosts attached to this network
			"vm",      // virtual machines using this network
		}
	}

	err = v.Retrieve(
		ctx,
		[]string{
			"Network", // managed object type we are retrieving from the view
		},
		props,
		&nets,
	)
	if err != nil {
		return nil, err
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
		fmt.Fprintf(
			os.Stderr,
			"It took %v to execute GetNetworkByName func.\n",
			time.Since(funcTimeStart),
		)
	}()

	finder := find.NewFinder(c, true)

	var dc *object.Datacenter
	var findDCErr error
	var errMsg string
	switch {
	case datacenter == "":
		dc, findDCErr = finder.DefaultDatacenter(ctx)
		errMsg = "error: datacenter not provided, failed to fallback to default datacenter"
	default:
		dc, findDCErr = finder.DatacenterOrDefault(ctx, datacenter)
		errMsg = "error: failed to use provided datacenter, failed to fallback to default datacenter"
	}

	if findDCErr != nil {
		return mo.Network{}, fmt.Errorf("%s: %w", errMsg, findDCErr)
	}
	finder.SetDatacenter(dc)

	netObj, err := finder.Network(ctx, netName)
	if err != nil {
		return mo.Network{}, err
	}

	pc := property.DefaultCollector(c)

	var props []string
	if propsSubset {
		// https://code.vmware.com/apis/1067/vsphere
		// https://vdc-download.vmware.com/vmwb-repository/dcr-public/a5f4000f-1ea8-48a9-9221-586adff3c557/7ff50256-2cf2-45ea-aacd-87d231ab1ac7/vim.Network.html
		props = []string{
			"summary", // properties of this network
			"name",    // name of this network
			"host",    // hosts attached to this network
			"vm",      // virtual machines using this network
		}
	}

	var network mo.Network
	err = pc.RetrieveOne(
		ctx,
		netObj.Reference(),
		props,
		&network,
	)

	if err != nil {
		return mo.Network{}, err
	}

	return network, nil

}

// FilterNetworkByName accepts a collection of Networks and a Network name to
// filter against. An error is returned if the list of Networks is empty or if
// a match was not found.
func FilterNetworkByName(nets []mo.Network, netName string) (mo.Network, error) {

	funcTimeStart := time.Now()

	defer func() {
		fmt.Fprintf(
			os.Stderr,
			"It took %v to execute FilterNetworkByName func.\n",
			time.Since(funcTimeStart),
		)
	}()

	if len(nets) == 0 {
		return mo.Network{}, fmt.Errorf("received empty list of networks to filter by name")
	}

	for _, net := range nets {
		if net.Name == netName {
			return net, nil
		}
	}

	return mo.Network{}, fmt.Errorf(
		"error: failed to retrieve Network using provided name %q",
		netName,
	)

}

// FilterNetworkByID receives a collection of Networks and a Network ID to
// filter against. An error is returned if the list of Networks is empty or if
// a match was not found.
func FilterNetworkByID(nets []mo.Network, netID string) (mo.Network, error) {

	funcTimeStart := time.Now()

	defer func() {
		fmt.Fprintf(
			os.Stderr,
			"It took %v to execute FilterNetworkByID func.\n",
			time.Since(funcTimeStart),
		)
	}()

	if len(nets) == 0 {
		return mo.Network{}, fmt.Errorf("received empty list of networks to filter by ID")
	}

	for _, net := range nets {
		// return match, if available
		if net.Self.Value == netID {
			return net, nil
		}
	}

	return mo.Network{}, fmt.Errorf(
		"error: failed to retrieve Network using provided id %q",
		netID,
	)
}
