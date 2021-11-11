// Copyright 2021 Adam Chalkley
//
// https://github.com/atc0005/check-vmware
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package config

// getSupportedDatastorePerfPercentiles is a helper function that returns a
// list of supported Datastore Performance Summary percentiles. This is used
// to provide validation of user-specified percentiles for Datastore
// Performance Percentile sets.
//
// FIXME: These percentiles are tightly bound to observed results from testing
// against vSphere 6.7.0 instances. This collection of supported percentiles
// should be validated against published API docs if at all possible.
func getSupportedDatastorePerfPercentiles() []int {

	return []int{
		90,
		80,
		70,
		60,
		50,
	}
}
