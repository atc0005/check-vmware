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

	"github.com/atc0005/check-vmware/internal/textutils"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

// GetFolders accepts a context, a connected client and a boolean value
// indicating whether a subset of properties per Folder are retrieved. If
// requested, a subset of all available properties will be retrieved (faster)
// instead of recursively fetching all properties (about 2x as slow) A
// collection of Folders with requested properties is returned or nil and an
// error, if one occurs.
func GetFolders(ctx context.Context, c *vim25.Client, propsSubset bool) ([]mo.Folder, error) {
	funcTimeStart := time.Now()

	// declare this early so that we can grab a pointer to it in order to
	// access the entries later
	var folders []mo.Folder

	defer func(folders *[]mo.Folder) {
		logger.Printf(
			"It took %v to execute GetFolders func (and retrieve %d Folders).\n",
			time.Since(funcTimeStart),
			len(*folders),
		)
	}(&folders)

	err := getObjects(ctx, c, &folders, c.ServiceContent.RootFolder, propsSubset, true)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve Folders: %w", err)
	}

	// FIXME: Should we sort retrieved folders? Does the original order matter?
	sort.Slice(folders, func(i, j int) bool {
		return strings.ToLower(folders[i].Name) < strings.ToLower(folders[j].Name)
	})

	return folders, nil
}

// FilterFoldersByID receives a collection of Folders and a Folder ID to
// filter against. An error is returned if the list of Folders is empty or if
// a match was not found. The matching Folder is returned along with the
// number of Folders that were excluded.
func FilterFoldersByID(folders []mo.Folder, folderID string) (mo.Folder, int, error) {
	funcTimeStart := time.Now()

	// If error condition, no Folders are excluded
	numExcluded := 0

	defer func() {
		logger.Printf(
			"It took %v to execute FilterFoldersByID func.\n",
			time.Since(funcTimeStart),
		)
	}()

	if len(folders) == 0 {
		return mo.Folder{},
			numExcluded,
			fmt.Errorf("received empty list of folders to filter by ID")
	}

	if folderID == "" {
		return mo.Folder{},
			numExcluded,
			fmt.Errorf("received empty folder ID to use as filter")
	}

	for _, folder := range folders {
		// return match, if available
		if folder.Self.Value == folderID {
			// we are excluding everything but the single ID value match
			numExcluded = len(folders) - 1
			return folder, numExcluded, nil
		}
	}

	return mo.Folder{}, numExcluded, fmt.Errorf(
		"error: failed to retrieve Folder using provided ID %q",
		folderID,
	)

}

// FolderNames receives a list of Folder values and returns a list of Folder
// Name values.
func FolderNames(foldersList []mo.Folder) []string {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute FolderNames func.\n",
			time.Since(funcTimeStart),
		)
	}()

	folderNames := make([]string, 0, len(foldersList))
	for _, folder := range foldersList {
		folderNames = append(folderNames, folder.Name)
	}

	return folderNames
}

// FolderManagedEntityVals receives a list of Folder values and returns a list
// of Folder ManagedEntity values.
func FolderManagedEntityVals(foldersList []mo.Folder) []mo.ManagedEntity {
	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute FolderManagedEntityVals func.\n",
			time.Since(funcTimeStart),
		)
	}()

	folderMEs := make([]mo.ManagedEntity, 0, len(foldersList))
	for _, folder := range foldersList {
		folderMEs = append(folderMEs, folder.ManagedEntity)
	}

	return folderMEs
}

// FolderNamesIDs receives a list of Folder values and returns a list of
// Folder Name/IP string values.
func FolderNamesIDs(foldersList []mo.Folder) []string {
	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute FolderNamesIDs func.\n",
			time.Since(funcTimeStart),
		)
	}()

	folderNameIDTmpl := "%s (ID: %s)"

	folderList := make([]string, 0, len(foldersList))
	for _, folder := range foldersList {
		folderList = append(
			folderList,
			fmt.Sprintf(
				folderNameIDTmpl,
				folder.Name,
				folder.Self.Value,
			),
		)
	}
	return folderList
}

// ValidateFolders is responsible for receiving two lists of Folder IDs,
// explicitly "included" (aka, "whitelisted") and explicitly "excluded" (aka,
// "blacklisted"). If any list entries are not found in the vSphere
// environment an error is returned listing which ones.
func ValidateFolders(ctx context.Context, c *vim25.Client, includeFolders []string, excludeFolders []string) error {

	funcTimeStart := time.Now()

	defer func(iFolders []string, eFolders []string) {
		logger.Printf(
			"It took %v to execute ValidateFolders func (and validate %d Folders).\n",
			time.Since(funcTimeStart),
			len(iFolders)+len(eFolders),
		)
	}(includeFolders, excludeFolders)

	m := view.NewManager(c)

	// Create a view of Folder objects
	v, createViewErr := m.CreateContainerView(
		ctx,
		c.ServiceContent.RootFolder,
		[]string{MgObjRefTypeFolder},
		true,
	)
	if createViewErr != nil {
		return fmt.Errorf("failed to create Folder view: %w", createViewErr)
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

	// Retrieve name property for all folders.
	props := []string{"name"}
	var folderSearchResults []mo.Folder
	retrieveErr := v.Retrieve(ctx, []string{MgObjRefTypeFolder}, props, &folderSearchResults)
	if retrieveErr != nil {
		return fmt.Errorf(
			"failed to retrieve Folder properties: %w",
			retrieveErr,
		)
	}

	foldersFound := make(map[string]mo.Folder, len(folderSearchResults))
	for _, folder := range folderSearchResults {
		foldersFound[folder.Self.Value] = folder
	}

	folderIDsFound := func() []string {
		ids := make([]string, 0, len(foldersFound))
		for id := range foldersFound {
			ids = append(ids, id)
		}
		return ids
	}()

	// If any specified folder names are not found, note that so we can
	// provide the full list of invalid names together as a convenience for
	// the user.
	var notFound []string
	switch {
	case len(includeFolders) > 0:
		for _, iFolderID := range includeFolders {
			if !textutils.InList(iFolderID, folderIDsFound, true) {
				notFound = append(
					notFound,
					fmt.Sprintf(
						"%s (ID: %s)",
						foldersFound[iFolderID].Name,
						iFolderID,
					),
				)
			}
		}

		if len(notFound) > 0 {
			return fmt.Errorf(
				"specified Folders (to include) not found: %v",
				notFound,
			)
		}

		// all specified folders were found
		return nil

	case len(excludeFolders) > 0:
		for _, eFolderID := range excludeFolders {
			if !textutils.InList(eFolderID, folderIDsFound, true) {
				notFound = append(
					notFound,
					fmt.Sprintf(
						"%s (ID: %s)",
						foldersFound[eFolderID].Name,
						eFolderID,
					),
				)
			}
		}

		if len(notFound) > 0 {
			return fmt.Errorf(
				"specified Folders (to exclude) not found: %v",
				notFound,
			)
		}

		// all specified folders were found
		return nil

	default:

		// no restrictions specified by user; all folders are "eligible" for
		// evaluation
		return nil
	}

}

// GetEligibleFolders receives a list of Folder IDs that should either be
// explicitly included or excluded along with a boolean value indicating
// whether only a subset of properties for the Folders should be returned. If
// requested, a subset of all available properties will be retrieved (faster)
// instead of recursively fetching all properties (about 2x as slow). The
// filtered list of Folders is returned, or an error if one occurs.
func GetEligibleFolders(ctx context.Context, c *vim25.Client, includeFolders []string, excludeFolders []string, propsSubset bool) ([]mo.Folder, error) {
	funcTimeStart := time.Now()

	// Declare slice early so that we can grab a pointer to it in order to
	// access the entries later. This holds the filtered list of folders that
	// will be returned to the caller.
	var eligibleFolders []mo.Folder

	defer func(folders *[]mo.Folder) {
		logger.Printf(
			"It took %v to execute GetEligibleFolders func (and retrieve %d Folders).\n",
			time.Since(funcTimeStart),
			len(*folders),
		)
	}(&eligibleFolders)

	// All available/accessible folders will be retrieved and stored here. We
	// will filter the results before returning a trimmed list to the caller.
	var folderSearchResults []mo.Folder

	err := getObjects(ctx, c, &folderSearchResults, c.ServiceContent.RootFolder, propsSubset, true)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve Folders: %w", err)
	}

	foldersFound := make(map[string]mo.Folder, len(folderSearchResults))
	for _, folder := range folderSearchResults {
		foldersFound[folder.Self.Value] = folder
	}

	listFolders := func(index map[string]mo.Folder) []string {
		folderNameIDs := make([]string, 0, len(index))
		for folderID, folder := range index {
			folderNameIDs = append(
				folderNameIDs,
				fmt.Sprintf("%s (ID: %s)", folder.Name, folderID),
			)
		}
		return folderNameIDs
	}

	logger.Printf(
		"Retrieved %d Folder objects: %v",
		len(folderSearchResults),
		strings.Join(listFolders(foldersFound), ", "),
	)

	for folderID, folder := range foldersFound {
		// Config validation asserts that only one of include/exclude folder
		// flags are specified.
		switch {

		// If specified, only include folders that have been intentionally
		// included (aka, "whitelisted").
		case len(includeFolders) > 0:
			if textutils.InList(folderID, includeFolders, true) {
				eligibleFolders = append(eligibleFolders, folder)
			}

		// If specified, don't include folders that have been intentionally
		// excluded (aka, "blacklisted").
		case len(excludeFolders) > 0:
			if !textutils.InList(folderID, excludeFolders, true) {
				eligibleFolders = append(eligibleFolders, folder)
			}

		// If we are not explicitly excluding or including folders, then we
		// can only assume (at this point in the filtering process) that we
		// are working with all folders. It is up to the caller to filter
		// further based on other criteria.
		default:
			eligibleFolders = append(eligibleFolders, folder)
		}
	}

	// TODO: SHOULD we sort based on folder name?
	//
	// When looking at debug logs it is probably more intuitive when seeing
	// the VM retrieval operations occurring in named order?
	sort.Slice(eligibleFolders, func(i, j int) bool {
		return strings.ToLower(eligibleFolders[i].Name) < strings.ToLower(eligibleFolders[j].Name)
	})

	return eligibleFolders, nil

}

// GetFoldersByIDs receives a list of Folder IDs that should resolved to
// Folder values along with a boolean value indicating whether only a subset
// of properties for the Folders should be returned. If requested, a subset of
// all available properties will be retrieved (faster) instead of recursively
// fetching all properties (about 2x as slow). The list of Folders is
// returned, or an error if one occurs.
func GetFoldersByIDs(ctx context.Context, c *vim25.Client, folderIDs []string, propsSubset bool) ([]mo.Folder, error) {
	funcTimeStart := time.Now()

	// Declare slice early so that we can grab a pointer to it in order to
	// access the entries later. This holds the filtered list of folders that
	// will be returned to the caller.
	var folders []mo.Folder

	defer func(folders *[]mo.Folder) {
		logger.Printf(
			"It took %v to execute GetFoldersByIDs func (and retrieve %d Folders).\n",
			time.Since(funcTimeStart),
			len(*folders),
		)
	}(&folders)

	if len(folderIDs) == 0 {
		return nil,
			fmt.Errorf("received empty list of folder IDs")
	}

	// All available/accessible folders will be retrieved and stored here. We
	// will filter the results before returning a trimmed list to the caller.
	var folderSearchResults []mo.Folder

	err := getObjects(ctx, c, &folderSearchResults, c.ServiceContent.RootFolder, propsSubset, true)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve Folders: %w", err)
	}

	foldersFound := make(map[string]mo.Folder, len(folderSearchResults))
	for _, folder := range folderSearchResults {
		foldersFound[folder.Self.Value] = folder
	}

	listFolders := func(index map[string]mo.Folder) []string {
		folderNameIDs := make([]string, 0, len(index))
		for folderID, folder := range index {
			folderNameIDs = append(
				folderNameIDs,
				fmt.Sprintf("%s (ID: %s)", folder.Name, folderID),
			)
		}
		return folderNameIDs
	}

	logger.Printf(
		"Retrieved %d Folder objects: %v",
		len(folderSearchResults),
		strings.Join(listFolders(foldersFound), ", "),
	)

	for folderID, folder := range foldersFound {
		if textutils.InList(folderID, folderIDs, true) {
			folders = append(folders, folder)
		}
	}

	// TODO: SHOULD we sort based on folder name?
	//
	// When looking at debug logs it is probably more intuitive when seeing
	// the VM retrieval operations occurring in named order?
	sort.Slice(folders, func(i, j int) bool {
		return strings.ToLower(folders[i].Name) < strings.ToLower(folders[j].Name)
	})

	return folders, nil
}

// getFoldersCountUsingContainerView accepts a context, a connected client, a
// container type ManagedObjectReference and a boolean value indicating
// whether the container type should be recursively searched for Folders. An
// error is returned if the provided ManagedObjectReference is not for a
// supported container type.
func getFoldersCountUsingContainerView(
	ctx context.Context,
	c *vim25.Client,
	containerRef types.ManagedObjectReference,
	recursive bool,
) (int, error) {

	funcTimeStart := time.Now()

	var allFolders []types.ObjectContent

	defer func(folders *[]types.ObjectContent, objRef types.ManagedObjectReference) {
		logger.Printf(
			"It took %v to execute getFoldersCountUsingContainerView func (and count %d Folders from %s).\n",
			time.Since(funcTimeStart),
			len(*folders),
			objRef.Type,
		)
	}(&allFolders, containerRef)

	// Create a view of caller-specified objects
	m := view.NewManager(c)

	logger.Printf("Container type is %s", containerRef.Type)

	// https://vdc-download.vmware.com/vmwb-repository/dcr-public/a5f4000f-1ea8-48a9-9221-586adff3c557/7ff50256-2cf2-45ea-aacd-87d231ab1ac7/vim.view.ContainerView.html
	switch containerRef.Type {

	// AFAIK only Folder container type can hold Folders?
	case MgObjRefTypeFolder:

	default:
		return 0, fmt.Errorf(
			"unsupported container type specified for ContainerView: %s",
			containerRef.Type,
		)
	}

	kind := []string{MgObjRefTypeFolder}

	// FIXME: Should this filter to a specific datacenter? See GH-219.
	v, createViewErr := m.CreateContainerView(
		ctx,
		containerRef,
		kind,
		recursive,
	)
	if createViewErr != nil {
		return 0, createViewErr
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

	// Perform as lightweight of a search as possible as we're only interested
	// in counting the total folders in a specified container.
	prop := []string{"overallStatus"}
	retrieveErr := v.Retrieve(ctx, kind, prop, &allFolders)
	if retrieveErr != nil {
		return 0, retrieveErr
	}

	return len(allFolders), nil
}

// validateFolders verifies that all explicitly specified Folders exist in the
// inventory.
func validateFolders(ctx context.Context, client *vim25.Client, filterOptions VMsFilterOptions) error {
	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute validateFolders func.\n",
			time.Since(funcTimeStart),
		)
	}()

	switch {
	case len(filterOptions.FoldersIncluded) > 0 || len(filterOptions.FoldersExcluded) > 0:
		logger.Println("Validating folders")

		validateErr := ValidateFolders(ctx, client, filterOptions.FoldersIncluded, filterOptions.FoldersExcluded)
		if validateErr != nil {
			logger.Printf(
				"%v: %v",
				ErrValidationOfIncludeExcludeFolderIDLists,
				validateErr,
			)

			return fmt.Errorf(
				"%v: %v",
				ErrValidationOfIncludeExcludeFolderIDLists,
				validateErr,
			)
		}
		logger.Println("Successfully validated folders")

		return nil
	default:
		logger.Println("Skipping folder validation; folder filtering not requested")
		return nil
	}
}

// GetNumTotalFolders returns the count of all Folders in the inventory.
func GetNumTotalFolders(ctx context.Context, client *vim25.Client) (int, error) {
	funcTimeStart := time.Now()

	var numAllFolders int

	defer func(allFolders *int) {
		logger.Printf(
			"It took %v to execute GetNumTotalFolders func (and count %d Folders).\n",
			time.Since(funcTimeStart),
			*allFolders,
		)
	}(&numAllFolders)

	var getFoldersErr error
	numAllFolders, getFoldersErr = getFoldersCountUsingContainerView(
		ctx,
		client,
		client.ServiceContent.RootFolder,
		true,
	)
	if getFoldersErr != nil {
		logger.Printf(
			"error retrieving list of all folders: %v",
			getFoldersErr,
		)

		return 0, fmt.Errorf(
			"error retrieving list of all folders: %w",
			getFoldersErr,
		)
	}
	logger.Printf(
		"Finished retrieving count of all folders: %d",
		numAllFolders,
	)

	return numAllFolders, nil
}
