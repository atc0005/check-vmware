// Copyright 2021 Adam Chalkley
//
// https://github.com/atc0005/check-vmware
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package vsphere

import (
	"errors"
	"strings"

	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

// ErrConvertBaseCustomFieldValue is returned when a conversion error occurs
// or (type assertion failure) for a provided BaseCustomFieldValue.
//
// TODO: Document how/when this might occur.
var ErrConvertBaseCustomFieldValue = errors.New("failed to convert base custom field value to obtain key/value pair")

// ErrCustomAttributeKeyNotFound is returned when a Custom Attribute key (and
// thus the desired value) cannot be found for a vSphere object. This can
// occur when a Custom Attribute exists within an inventory, but is not
// applied to a specific managed object.
var ErrCustomAttributeKeyNotFound = errors.New("failed to find a matching custom attribute key/value pair for provided custom attribute key")

// ErrCustomAttributeNotSet is similar to ErrCustomAttributeKeyNotFound, but
// is returned when there are no Custom Attributes set for an associated
// managed object type.
var ErrCustomAttributeNotSet = errors.New("custom attributes not set")

// ErrAvailableFieldValueNotFound is returned when a specified Custom
// Attribute name cannot be located within an inventory. This might be the
// case if a specified Custom Attribute name has a typo.
var ErrAvailableFieldValueNotFound = errors.New("failed to find a matching available field name")

// ErrAvailableFieldValueNotDefined is returned when no Custom Attributes are
// defined within an inventory for an associated managed object type.
var ErrAvailableFieldValueNotDefined = errors.New("no custom attributes defined within vSphere environment for this type")

// CustomAttrKeyToValue receives the key of a Custom Attribute key/value pair
// a slice of BaseCustomFieldValue which represents the Custom Attribute
// values for vSphere objects (e.g., HostSystem, Datastore). An error is
// returned if conversion fails or if the specified Custom Attribute key is
// not found.
func CustomAttrKeyToValue(caKey int32, customValue []types.BaseCustomFieldValue) (string, error) {

	switch {

	// this vSphere object has at least one custom attribute set
	case len(customValue) > 0:

		for _, bcfv := range customValue {

			// TODO: Document scenarios where this might fail.
			cfsv, ok := bcfv.(*types.CustomFieldStringValue)
			if !ok {
				return "", ErrConvertBaseCustomFieldValue
			}

			// fmt.Printf("Object details: %#v\n", bcfv)
			if caKey == cfsv.Key {
				return cfsv.Value, nil
			}
		}

		return "", ErrCustomAttributeKeyNotFound

	// this vSphere object does not have any custom attributes set
	default:

		return "", ErrCustomAttributeNotSet
	}

}

// CustomAttrNameToKey receives a user-provided Custom Attribute name and a
// slice of values which contain key/value fields. If a match for the
// user-provided Custom Attribute name is found, the matching key is returned.
// This key can be used to search for a match in the Custom Values associated
// with vSphere objects (e.g., Hosts, Datastores). An error is returned if a
// match is not found.
func CustomAttrNameToKey(caName string, availableField []types.CustomFieldDef) (int32, error) {

	switch {
	// this vSphere object has at least one custom attribute defined for its
	// type (though not necessarily set)
	case len(availableField) > 0:

		for _, af := range availableField {

			// Allow user to provide Custom Attribute name in any mix of case
			if strings.EqualFold(af.Name, caName) {
				return af.Key, nil
			}
		}

		return 0, ErrAvailableFieldValueNotFound

	// this vSphere object has no custom attributes available (defined within
	// vSphere environment) for its type

	default:

		return 0, ErrAvailableFieldValueNotDefined

	}

}

// GetObjectCAVal receives the name of a Custom Attribute and a ManagedEntity,
// an abstract base type for all managed objects present in the inventory tree
// and returns the value for the specified Custom Attribute. An error is
// returned if the value could not be retrieved indicating the cause of the
// failure.
func GetObjectCAVal(caName string, obj mo.ManagedEntity) (string, error) {

	caKey, keyLookupErr := CustomAttrNameToKey(caName, obj.AvailableField)
	if keyLookupErr != nil {
		return "", keyLookupErr
	}

	caValue, valLookupErr := CustomAttrKeyToValue(caKey, obj.CustomValue)
	if valLookupErr != nil {
		return "", valLookupErr

	}

	return caValue, nil

}
