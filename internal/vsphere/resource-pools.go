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

	"github.com/atc0005/check-vmware/internal/textutils"

	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
)

// ValidateRPs is responsible for receiving two lists of resource pools,
// explicitly "included" (aka, "whitelisted") and explicitly "excluded" (aka,
// "blacklisted"). If any list entries are not found in the vSphere
// environment an error is returned listing which ones.
func ValidateRPs(ctx context.Context, c *vim25.Client, includeRPs []string, excludeRPs []string) error {

	funcTimeStart := time.Now()

	defer func(irps []string, erps []string) {
		fmt.Fprintf(
			os.Stderr,
			"It took %v to execute ValidateRPs func (and validate %d Resource Pools).\n",
			time.Since(funcTimeStart),
			len(irps)+len(erps),
		)
	}(includeRPs, excludeRPs)

	m := view.NewManager(c)

	// Create a view of Resource Pool objects
	v, createViewErr := m.CreateContainerView(
		ctx,
		c.ServiceContent.RootFolder,
		[]string{
			"ResourcePool",
		},
		true,
	)
	if createViewErr != nil {
		return fmt.Errorf("failed to create ResourcePool view: %w", createViewErr)
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

	// Retrieve name property for all resource pools.
	props := []string{"name"}
	var rpsSearchResults []mo.ResourcePool
	retrieveErr := v.Retrieve(ctx, []string{"ResourcePool"}, props, &rpsSearchResults)
	if retrieveErr != nil {
		return fmt.Errorf(
			"failed to retrieve ResourcePool properties: %w",
			retrieveErr,
		)
	}

	// We're only interested in working with resource pool names
	poolNamesFound := make([]string, 0, len(rpsSearchResults))
	for _, rp := range rpsSearchResults {
		poolNamesFound = append(poolNamesFound, rp.Name)
	}

	// If any specified resource pool names are not found, note that so we can
	// provide the full list of invalid pool names together as a convenience
	// for the user.
	var notFound []string
	switch {
	case len(includeRPs) > 0:
		for _, iRP := range includeRPs {
			if !textutils.InList(iRP, poolNamesFound, true) {
				notFound = append(notFound, iRP)
			}
		}

		if len(notFound) > 0 {
			return fmt.Errorf(
				"specified Resource Pools (to include) not found: %v",
				notFound,
			)
		}

		// all listed resource pools were found
		return nil

	case len(excludeRPs) > 0:
		for _, eRP := range excludeRPs {
			if !textutils.InList(eRP, poolNamesFound, true) {
				notFound = append(notFound, eRP)
			}
		}

		if len(notFound) > 0 {
			return fmt.Errorf(
				"specified Resource Pools (to exclude) not found: %v",
				notFound,
			)
		}

		// all listed resource pools were found
		return nil

	default:

		// no restrictions specified by user; all resource pools are
		// "eligible" for evaluation
		return nil
	}

}

// GetEligibleRPs receives a list of Resource Pool names that should either be
// explicitly included or excluded along with a boolean value indicating
// whether only a subset of properties for the Resource Pools should be
// returned. If requested, a subset of all available properties will be
// retrieved (faster) instead of recursively fetching all properties (about 2x
// as slow). The filtered list of Resource Pools is returned, or an error if
// one occurs.
func GetEligibleRPs(ctx context.Context, c *vim25.Client, includeRPs []string, excludeRPs []string, propsSubset bool) ([]mo.ResourcePool, error) {

	funcTimeStart := time.Now()

	// declare this early so that we can grab a pointer to it in order to
	// access the entries later
	var rps []mo.ResourcePool

	defer func(rps *[]mo.ResourcePool) {
		fmt.Fprintf(
			os.Stderr,
			"It took %v to execute GetEligibleRPs func (and retrieve %d Resource Pools).\n",
			time.Since(funcTimeStart),
			len(*rps),
		)
	}(&rps)

	// By default, all resource pools will be retrieved. We will filter and
	// return a trimmed list.
	var rpsSearchResults []mo.ResourcePool

	err := getObjects(ctx, c, &rpsSearchResults, c.ServiceContent.RootFolder, propsSubset)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve ResourcePools: %w", err)
	}

	for _, rp := range rpsSearchResults {

		// Virtual machine hosts have a hidden resource pool named Resources,
		// which is a parent of all resource pools of the host. Including this
		// pool can throw off our calculations, so we ignore it *unless* the
		// list of provided Resource Pools to explicitly include is empty.
		// Because this is a hidden pool and non-obvious, we try to avoid
		// requiring the sysadmin to specify it explicitly.
		if strings.EqualFold(rp.Name, ParentResourcePool) {
			if len(includeRPs) != 0 {
				continue
			}
		}

		// config validation asserts that only one of include/exclude resource
		// pools flags are specified
		switch {

		// if specified, only include resource pools that have been
		// intentionally included (aka, "whitelisted")
		case len(includeRPs) > 0:
			if textutils.InList(rp.Name, includeRPs, true) {
				rps = append(rps, rp)
			}

		// if specified, don't include resource pools that have been
		// intentionally excluded (aka, "blacklisted")
		case len(excludeRPs) > 0:
			if !textutils.InList(rp.Name, excludeRPs, true) {
				rps = append(rps, rp)
			}

		// if we are not explicitly excluding or including pools, then we are
		// working with all pools
		default:
			rps = append(rps, rp)
		}

	}

	sort.Slice(rps, func(i, j int) bool {
		return strings.ToLower(rps[i].Name) < strings.ToLower(rps[j].Name)
	})

	return rps, nil

}

// GetRPByName accepts the name of a Resource Pool, the name of a datacenter
// and a boolean value indicating whether only a subset of properties for the
// Network should be returned. If requested, a subset of all available
// properties will be retrieved (faster) instead of recursively fetching all
// properties (about 2x as slow). If the datacenter name is an empty string
// then the default datacenter will be used.
func GetRPByName(ctx context.Context, c *vim25.Client, rpName string, datacenter string, propsSubset bool) (mo.ResourcePool, error) {

	funcTimeStart := time.Now()

	defer func() {
		fmt.Fprintf(
			os.Stderr,
			"It took %v to execute GetRPByName func.\n",
			time.Since(funcTimeStart),
		)
	}()

	var rPool mo.ResourcePool
	err := getObjectByName(ctx, c, &rPool, rpName, datacenter, propsSubset)

	if err != nil {
		return mo.ResourcePool{}, err
	}

	return rPool, nil

}
