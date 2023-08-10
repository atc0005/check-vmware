// Copyright 2021 Adam Chalkley
//
// https://github.com/atc0005/check-vmware
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package config

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// DSPerformanceSummaryThresholds represents the thresholds used to evaluate
// Datastore Performance Summary values.
type DSPerformanceSummaryThresholds struct {
	// ReadLatencyWarning is the read latency in ms when a WARNING threshold
	// is reached.
	ReadLatencyWarning float64

	// ReadLatencyCritical is the read latency in ms when a CRITICAL threshold
	// is reached.
	ReadLatencyCritical float64

	// WriteLatencyWarning is the write latency in ms when a WARNING threshold
	// is reached.
	WriteLatencyWarning float64

	// WriteLatencyCritical is the write latency in ms when a CRITICAL
	// threshold is reached.
	WriteLatencyCritical float64

	// VMLatencyWarning is the latency in ms as observed by VMs using the
	// datastore when a WARNING threshold is reached.
	VMLatencyWarning float64

	// VMLatencyCritical is the latency in ms as observed by VMs using the
	// datastore when a CRITICAL threshold is reached.
	VMLatencyCritical float64
}

// dsPerfLatencyMetricFlag is a custom type that satisfies the flag.Value
// interface. This type is used to accept Datastore Performance Summary
// latency metric values. This flag type is incompatible with the flag type
// used to specify percentile sets.
type dsPerfLatencyMetricFlag struct {

	// value is the user-specified value
	value float64

	// isSet identifies whether a value was provided by the user
	isSet bool
}

// MultiValueDSPerfPercentileSetFlag is a custom type that satisfies the
// flag.Value interface. This type is used to accept Datastore Performance
// Summary percentile "sets". These sets define thresholds used to check
// Datastore Performance latency metrics to determine overall plugin state.
type MultiValueDSPerfPercentileSetFlag map[int]DSPerformanceSummaryThresholds

// String satisfies the flag.Value interface method set requirements.
func (dspl *dsPerfLatencyMetricFlag) String() string {

	// The String() method is called by the flag.isZeroValue function in order
	// to determine whether the output string represents the zero value for a
	// flag. This occurs even if the flag is not specified by the user.

	if dspl == nil {
		return ""
	}

	return fmt.Sprintf(
		"value: %v, isSet: %v",
		dspl.value,
		dspl.isSet,
	)

}

// Set satisfies the flag.Value interface method set requirements.
func (dspl *dsPerfLatencyMetricFlag) Set(value string) error {

	// fmt.Println("dsPerfLatencyMetricFlag Set() called")

	var strConvErr error

	value = strings.TrimSpace(value)
	value = strings.ReplaceAll(value, "'", "")
	value = strings.ReplaceAll(value, "\"", "")

	var parsedVal float64
	parsedVal, strConvErr = strconv.ParseFloat(value, 64)
	if strConvErr != nil {
		return fmt.Errorf(
			"error processing flag; failed to convert %q: %v",
			value,
			strConvErr,
		)
	}

	dspl.value = parsedVal
	dspl.isSet = true

	return nil

}

// String returns a comma separated string consisting of all map entries.
func (mvdsperf *MultiValueDSPerfPercentileSetFlag) String() string {

	// The String() method is called by the flag.isZeroValue function in order
	// to determine whether the output string represents the zero value for a
	// flag. This occurs even if the flag is not specified by the user.

	// From the `flag` package docs:
	// "The flag package may call the String method with a zero-valued
	// receiver, such as a nil pointer."
	if mvdsperf == nil {
		return "empty percentile set"
	}

	percentiles := make([]int, 0, len(*mvdsperf))
	for key := range *mvdsperf {
		percentiles = append(percentiles, key)
	}

	sort.Slice(percentiles, func(i, j int) bool {
		return percentiles[i] < percentiles[j]
	})

	var output strings.Builder

	for _, p := range percentiles {
		fmt.Fprintf(&output,
			"{Percentile: %v, ThresholdVals: %+v}, ",
			p,
			(*mvdsperf)[p],
		)
	}

	outputString := strings.TrimSuffix(output.String(), ", ")

	return outputString

}

// thresholdValues receives a string indicating either WARNING or CRITICAL
// state and returns a comma separated string consisting of all specified
// metric percentiles and the the associated WARNING or CRITICAL threshold
// values.
func (mvdsperf MultiValueDSPerfPercentileSetFlag) thresholdValues(state string) string {

	// From the `flag` package docs:
	// "The flag package may call the String method with a zero-valued
	// receiver, such as a nil pointer."
	if mvdsperf == nil {
		return "empty percentile set"
	}

	percentiles := make([]int, 0, len(mvdsperf))
	for key := range mvdsperf {
		percentiles = append(percentiles, key)
	}

	sort.Slice(percentiles, func(i, j int) bool {
		return percentiles[i] < percentiles[j]
	})

	var output strings.Builder

	var readLatency float64
	var writeLatency float64
	var vmLatency float64

	for _, p := range percentiles {

		switch {
		case strings.ToUpper(state) == StateCRITICALLabel:
			readLatency = mvdsperf[p].ReadLatencyCritical
			writeLatency = mvdsperf[p].WriteLatencyCritical
			vmLatency = mvdsperf[p].VMLatencyCritical

			// fmt.Printf(
			// 	"CRITICAL | readLatency: %v, writeLatency: %v, vmLatency: %v\n",
			// 	readLatency,
			// 	writeLatency,
			// 	vmLatency,
			// )

		case strings.ToUpper(state) == StateWARNINGLabel:
			readLatency = mvdsperf[p].ReadLatencyWarning
			writeLatency = mvdsperf[p].WriteLatencyWarning
			vmLatency = mvdsperf[p].VMLatencyWarning

			// fmt.Printf(
			// 	"WARNING | readLatency: %v, writeLatency: %v, vmLatency: %v\n",
			// 	readLatency,
			// 	writeLatency,
			// 	vmLatency,
			// )
		}

		fmt.Fprintf(&output,
			"{ Percentile: %v, ReadLatency: %+v, WriteLatency: %v, VMLatency: %v }, ",
			p,
			readLatency,
			writeLatency,
			vmLatency,
		)
	}

	outputString := strings.TrimSuffix(output.String(), ", ")

	return outputString

}

// CriticalThresholdValues returns a comma separated string consisting of all
// specified metric percentiles and the the associated CRITICAL threshold
// values.
func (mvdsperf MultiValueDSPerfPercentileSetFlag) CriticalThresholdValues() string {
	return mvdsperf.thresholdValues(StateCRITICALLabel)
}

// WarningThresholdValues returns a comma separated string consisting of all
// specified metric percentiles and the the associated WARNING threshold
// values.
func (mvdsperf MultiValueDSPerfPercentileSetFlag) WarningThresholdValues() string {
	return mvdsperf.thresholdValues(StateWARNINGLabel)
}

// Set is called once by the flag package, in command line order, for each
// flag present.
func (mvdsperf *MultiValueDSPerfPercentileSetFlag) Set(value string) error {

	// We require the same number of values as we have fields in the struct
	// plus one more to serve as the map index (percentile).
	const expectedValues int = 7

	// Split comma-separated string into multiple values, toss whitespace,
	// then convert value in string format to integer.
	items := strings.Split(value, ",")

	if len(items) != expectedValues {
		return fmt.Errorf(
			"error processing flag; string %q provides %d values, expected %d values",
			value,
			len(items),
			expectedValues,
		)
	}

	// fmt.Println("items", items)

	percentileSet := make([]float64, len(items))
	var strConvErr error
	for i := range items {
		items[i] = strings.TrimSpace(items[i])
		items[i] = strings.ReplaceAll(items[i], "'", "")
		items[i] = strings.ReplaceAll(items[i], "\"", "")

		percentileSet[i], strConvErr = strconv.ParseFloat(strings.TrimSpace(items[i]), 64)
		if strConvErr != nil {
			return fmt.Errorf(
				"error processing flag; failed to convert %q: %v",
				items[i],
				strConvErr,
			)
		}

	}

	// We now have the latency values (along with the percentile) stored as
	// float64 values. The first element is the percentile which is an int.
	percentile := int(percentileSet[0])

	// fmt.Printf("mvdsperf before assignment to map: %+v (nil: %t)\n", mvdsperf, mvdsperf == nil)

	// The rest of the latency values have already been converted to the
	// necessary type, so we assign directly.
	(*mvdsperf)[percentile] = DSPerformanceSummaryThresholds{
		ReadLatencyWarning:   percentileSet[1],
		ReadLatencyCritical:  percentileSet[2],
		WriteLatencyWarning:  percentileSet[3],
		WriteLatencyCritical: percentileSet[4],
		VMLatencyWarning:     percentileSet[5],
		VMLatencyCritical:    percentileSet[6],
	}

	// 	fmt.Printf("mvdsperf[percentile]: %+v\n", mvdsperf[percentile])
	//
	// 	fmt.Printf("mvdsperf after assignment to map: %+v (nil: %t)\n", mvdsperf, mvdsperf == nil)

	return nil

}
