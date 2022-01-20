// Copyright 2021 Adam Chalkley
//
// https://github.com/atc0005/check-vmware
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package vsphere

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/atc0005/check-vmware/internal/textutils"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
)

// ErrDatacenterNotFound indicates that one or more specified datacenters were
// not located.
var ErrDatacenterNotFound = errors.New("specified Datacenters not found")

// ValidateDCs receives a list of Datacenter names and compares against all
// visible Datacenter objects within the vSphere environment. If any are not
// found an error is returned listing which ones. If an empty list of
// Datacenter names is provided validation is considered successful.
func ValidateDCs(ctx context.Context, c *vim25.Client, datacenters []string) error {

	funcTimeStart := time.Now()

	defer func(datacenters []string) {
		logger.Printf(
			"It took %v to execute ValidateDCs func (and validate %d Datacenters).\n",
			time.Since(funcTimeStart),
			len(datacenters),
		)
	}(datacenters)

	// If the requested list to validate is empty, declare successful
	// validation.
	if len(datacenters) == 0 {
		return nil
	}

	m := view.NewManager(c)

	// Create a view of Datacenter objects
	v, createViewErr := m.CreateContainerView(
		ctx,
		c.ServiceContent.RootFolder,
		[]string{
			"Datacenter",
		},
		true,
	)
	if createViewErr != nil {
		return fmt.Errorf("failed to create Datacenter view: %w", createViewErr)
	}

	defer func() {
		// Per vSphere Web Services SDK Programming Guide - VMware vSphere 7.0
		// Update 1:
		//
		// A best practice when using views is to call the DestroyView()
		// method when a view is no longer needed. This practice frees memory
		// on the server.
		if err := v.Destroy(ctx); err != nil {
			logger.Printf("Error occurred while destroying view: %s", err)
		}
	}()

	// Retrieve only the name property for all Datacenters.
	props := []string{"name"}
	var dcsSearchResults []mo.Datacenter
	err := v.Retrieve(ctx, []string{"Datacenter"}, props, &dcsSearchResults)
	if err != nil {
		return fmt.Errorf(
			"failed to retrieve Datacenter properties: %w",
			err,
		)
	}

	// Gather Datacenter names
	dcNamesFound := make([]string, 0, len(dcsSearchResults))
	for _, dc := range dcsSearchResults {
		dcNamesFound = append(dcNamesFound, dc.Name)
	}

	// If any specified Datacenter names are not found, note that so we can
	// provide the full list of invalid names together as a convenience for
	// the user.
	var notFound []string

	for _, dc := range datacenters {
		if !textutils.InList(dc, dcNamesFound, true) {
			notFound = append(notFound, dc)
		}
	}

	if len(notFound) > 0 {
		return fmt.Errorf(
			"%w: %v",
			ErrDatacenterNotFound,
			notFound,
		)
	}

	// all specified Datacenters were found
	return nil

}

// GetDatacenters receives a list of Datacenter names along with a boolean
// value indicating whether only a subset of properties for the Datacenters
// should be returned. If requested, a subset of all available properties will
// be retrieved (faster) instead of recursively fetching all properties (about
// 2x as slow). If an empty list of Datacenter names is provided then all
// visible Datacenters will be retrieved. The list of Datacenters found is
// returned, or an error if one occurs.
func GetDatacenters(ctx context.Context, c *vim25.Client, dcNames []string, propsSubset bool) ([]mo.Datacenter, error) {

	funcTimeStart := time.Now()

	// declare this early so that we can grab a pointer to it in order to
	// access the entries later
	var filteredDCs []mo.Datacenter

	defer func(dcs *[]mo.Datacenter) {
		logger.Printf(
			"It took %v to execute GetDatacenters func (and retrieve %d Datacenters).\n",
			time.Since(funcTimeStart),
			len(*dcs),
		)
	}(&filteredDCs)

	switch {
	case dcNames == nil:
		logger.Println("empty datacenters list provided")
	default:
		logger.Println("one or more datacenters specified")
	}

	// Fetch all visible datacenters.
	var allDCs []mo.Datacenter

	err := getObjects(ctx, c, &allDCs, c.ServiceContent.RootFolder, propsSubset, true)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve all Datacenters: %w", err)
	}

	switch {
	// If a specific list of Datacenter names was not provided, return all
	// Datacenters that are visible.
	case len(dcNames) == 0:
		filteredDCs = allDCs
		return filteredDCs, nil
	default:
		for _, dc := range allDCs {
			if textutils.InList(dc.Name, dcNames, true) {
				filteredDCs = append(filteredDCs, dc)
			}
		}
		return filteredDCs, nil
	}

}
