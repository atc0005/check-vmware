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
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

// DatastoreIDToNameIndex maps a Datastore's ID value to its name.
type DatastoreIDToNameIndex map[string]string

// GetDatastores accepts a context, a connected client and a boolean value
// indicating whether a subset of properties per Datastore are retrieved. A
// collection of Datastores with requested properties is returned. If
// requested, a subset of all available properties will be retrieved (faster)
// instead of recursively fetching all properties (about 2x as slow).
func GetDatastores(ctx context.Context, c *vim25.Client, propsSubset bool) ([]mo.Datastore, error) {

	funcTimeStart := time.Now()

	// declare this early so that we can grab a pointer to it in order to
	// access the entries later
	var dss []mo.Datastore

	defer func(dss *[]mo.Datastore) {
		fmt.Fprintf(
			os.Stderr,
			"It took %v to execute GetDatastores func (and retrieve %d Datastores).\n",
			time.Since(funcTimeStart),
			len(*dss),
		)
	}(&dss)

	err := getObjects(ctx, c, &dss, c.ServiceContent.RootFolder, propsSubset)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve Datastores: %w", err)
	}

	sort.Slice(dss, func(i, j int) bool {
		return strings.ToLower(dss[i].Name) < strings.ToLower(dss[j].Name)
	})

	return dss, nil
}

// GetDatastoreByName accepts the name of a network, the name of a datacenter
// and a boolean value indicating whether only a subset of properties for the
// Datastore should be returned. If requested, a subset of all available
// properties will be retrieved (faster) instead of recursively fetching all
// properties (about 2x as slow). If the datacenter name is an empty string
// then the default datacenter will be used.
func GetDatastoreByName(ctx context.Context, c *vim25.Client, dsName string, datacenter string, propsSubset bool) (mo.Datastore, error) {

	funcTimeStart := time.Now()

	defer func() {
		fmt.Fprintf(
			os.Stderr,
			"It took %v to execute GetDatastoreByName func.\n",
			time.Since(funcTimeStart),
		)
	}()

	var datastore mo.Datastore

	err := getObjectByName(ctx, c, &datastore, dsName, datacenter, propsSubset)
	if err != nil {
		return mo.Datastore{}, err
	}

	return datastore, nil

}

// FilterDatastoreByName accepts a collection of Datastores and a Datastore
// name to filter against. An error is returned if the list of Datastores is
// empty or if a match was not found.
func FilterDatastoreByName(dss []mo.Datastore, dsName string) (mo.Datastore, error) {

	funcTimeStart := time.Now()

	defer func() {
		fmt.Fprintf(
			os.Stderr,
			"It took %v to execute FilterDatastoreByName func.\n",
			time.Since(funcTimeStart),
		)
	}()

	if len(dss) == 0 {
		return mo.Datastore{}, fmt.Errorf("received empty list of datastores to filter by name")
	}

	for _, ds := range dss {
		if ds.Name == dsName {
			return ds, nil
		}
	}

	return mo.Datastore{}, fmt.Errorf(
		"error: failed to retrieve Datastore using provided name %q",
		dsName,
	)

}

// FilterDatastoreByID receives a collection of Datastores and a Datastore ID
// to filter against. An error is returned if the list of Datastores is empty
// or if a match was not found.
func FilterDatastoreByID(dss []mo.Datastore, dsID string) (mo.Datastore, error) {

	funcTimeStart := time.Now()

	defer func() {
		fmt.Fprintf(
			os.Stderr,
			"It took %v to execute FilterDatastoreByID func.\n",
			time.Since(funcTimeStart),
		)
	}()

	if len(dss) == 0 {
		return mo.Datastore{}, fmt.Errorf("received empty list of datastores to filter by ID")
	}

	for _, ds := range dss {
		// return match, if available
		// TODO: Refactor, use abstract type here
		// ds.GetManagedEntity().Reference().Value
		if ds.Summary.Datastore.Value == dsID {
			return ds, nil
		}
	}

	return mo.Datastore{}, fmt.Errorf(
		"error: failed to retrieve Datastore using provided id %q",
		dsID,
	)

}

// DatastoreIDsToNames returns a list of matching Datastore names for the
// provided list of Managed Object References for Datastores.
func DatastoreIDsToNames(dsRefs []types.ManagedObjectReference, dss []mo.Datastore) []string {

	dsNames := make([]string, 0, len(dsRefs))
	dsIDs := make([]string, 0, len(dsRefs))

	for _, dsRef := range dsRefs {
		dsIDs = append(dsIDs, dsRef.Value)
	}

	for _, ds := range dss {
		if textutils.InList(ds.Summary.Datastore.Value, dsIDs, true) {
			dsNames = append(dsNames, ds.Name)
		}
	}

	return dsNames

}
