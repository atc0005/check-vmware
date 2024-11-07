// Copyright 2023 Codeweavers Ltd
// Copyright 2023 Adam Chalkley
//
// https://github.com/atc0005/go-nagios
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.
//
// Portions of the code in this file inspired by or generated with the help of
// ChatGPT and Google Gemini.

package nagios

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Range represents the thresholds that the user can pass in for warning and
// critical, this format is based on the [Nagios Plugin Dev Guidelines:
// Threshold and Ranges] definition.
//
// [Nagios Plugin Dev Guidelines: Threshold and Ranges]: https://nagios-plugins.org/doc/guidelines.html#THRESHOLDFORMAT
type Range struct {
	StartInfinity bool
	EndInfinity   bool
	AlertOn       string
	Start         float64
	End           float64
}

// CheckRange returns true if an alert should be raised for a given
// performance data Value, otherwise false.
func (r Range) CheckRange(value string) bool {
	valueAsAFloat, _ := strconv.ParseFloat(value, 64)
	isOutsideRange := r.checkOutsideRange(valueAsAFloat)
	if r.AlertOn == "INSIDE" {
		return !isOutsideRange
	}
	return isOutsideRange
}

// checkOutsideRange returns in the inverse of CheckRange. It is used to
// handle the inverting logic of "inside" vs "outside" ranges.
//
// See the [Nagios Plugin Dev Guidelines: Threshold and Ranges] definition for
// additional details.
//
// [Nagios Plugin Dev Guidelines: Threshold and Ranges]: https://nagios-plugins.org/doc/guidelines.html#THRESHOLDFORMAT
func (r Range) checkOutsideRange(valueAsAFloat float64) bool {
	// This alternative implementation was provided by Google Gemini (model
	// 'Gemini 1.5 Flash').
	//
	// Explanation of the logic:
	//
	// Infinite Bounds:
	//
	// If the start is infinite, the value is outside the range only if it's
	// greater than the end. If the end is infinite, the value is outside the
	// range only if it's less than the start.
	//
	// Finite Bounds:
	//
	// The value is outside the range if it's either less than the start or
	// greater than the end.

	// Handle infinite bounds first
	if r.StartInfinity {
		return valueAsAFloat > r.End
	} else if r.EndInfinity {
		return valueAsAFloat < r.Start
	}

	// Handle finite bounds
	return valueAsAFloat < r.Start || valueAsAFloat > r.End
}

// checkOutsideRange is provided by ChatGPT (model 'GPT-4o') as a
// simplification of the original checkOutsideRange function.
// func (r Range) checkOutsideRange(value float64) bool {
// 	// Explanation of the Simplification:
// 	//
// 	// Each case is now focused only on the conditions that make the value
// 	// outside the range. The final default case covers the fully infinite
// 	// range, simplifying the logic to just return false since no bounds
// 	// restrict the range.
//
// 	switch {
// 	case !r.StartInfinity && value < r.Start:
// 		return true
// 	case !r.EndInfinity && value > r.End:
// 		return true
// 	case r.StartInfinity && r.EndInfinity:
// 		return false
// 	default:
// 		return false
// 	}
// }

// ParseRangeString static method to construct a Range object from the string
// representation based on the [Nagios Plugin Dev Guidelines: Threshold and
// Ranges] definition.
//
// [Nagios Plugin Dev Guidelines: Threshold and Ranges]: https://nagios-plugins.org/doc/guidelines.html#THRESHOLDFORMAT
func ParseRangeString(input string) *Range {
	// Initialize range with default values
	r := Range{
		Start:         0,
		End:           0,
		StartInfinity: false,
		EndInfinity:   false,
		AlertOn:       "OUTSIDE",
	}

	// Define regular expressions
	digitOrInfinity := regexp.MustCompile(`[\d~]`)
	optionalInvertAndRange := regexp.MustCompile(`^@?([-+]?[\d.]+(?:e[-+]?[\d.]+)?|~)?(:([-+]?[\d.]+(?:e[-+]?[\d.]+)?)?)?$`)
	firstHalfOfRange := regexp.MustCompile(`^([-+]?[\d.]+(?:e[-+]?[\d.]+)?)?:`)
	endOfRange := regexp.MustCompile(`^[-+]?[\d.]+(?:e[-+]?[\d.]+)?$`)

	// Validate input format
	if !(digitOrInfinity.MatchString(input) && optionalInvertAndRange.MatchString(input)) {
		return nil
	}

	switch {
	// Parse alert inversion (starts with @)
	case strings.HasPrefix(input, "@"):
		r.AlertOn = "INSIDE"
		input = input[1:]

	// Parse start infinity (~ symbol at start)
	case strings.HasPrefix(input, "~"):
		r.StartInfinity = true
		input = input[1:]
	}

	// Parse start of range (e.g., "10:")
	if rangeComponents := firstHalfOfRange.FindStringSubmatch(input); rangeComponents != nil {
		if rangeComponents[1] != "" {
			r.Start, _ = strconv.ParseFloat(rangeComponents[1], 64)
			r.StartInfinity = false
		}
		r.EndInfinity = true
		input = strings.TrimPrefix(input, rangeComponents[0])
	}

	// Parse end of range (e.g., "10" or "x:10")
	if endOfRangeComponents := endOfRange.FindStringSubmatch(input); endOfRangeComponents != nil {
		r.End, _ = strconv.ParseFloat(endOfRangeComponents[0], 64)
		r.EndInfinity = false
	}

	// Ensure valid range boundaries
	if r.StartInfinity || r.EndInfinity || r.Start <= r.End {
		return &r
	}

	return nil
}

// EvaluateThreshold causes the performance data to be checked against the
// Warn and Crit thresholds provided by client code and sets the
// ExitStatusCode of the plugin as appropriate.
func (p *Plugin) EvaluateThreshold(perfData ...PerformanceData) error {
	for i := range perfData {
		// Evaluate critical threshold
		if inCritical, err := evaluateThreshold(perfData[i].Crit, perfData[i].Value); err != nil {
			p.ExitStatusCode = StateUNKNOWNExitCode
			return err
		} else if inCritical {
			p.ExitStatusCode = StateCRITICALExitCode
			return nil
		}

		// Evaluate warning threshold
		if inWarning, err := evaluateThreshold(perfData[i].Warn, perfData[i].Value); err != nil {
			p.ExitStatusCode = StateUNKNOWNExitCode
			return err
		} else if inWarning {
			p.ExitStatusCode = StateWARNINGExitCode
			return nil
		}
	}

	return nil
}

// evaluateThreshold is a helper function used to handle both parsing and
// range-checking, taking rangeStr (the threshold string), value, and
// exitCode. If the parsing fails, it returns an error to simplify error
// handling within the caller.
func evaluateThreshold(rangeStr, value string) (bool, error) {
	if rangeStr == "" {
		return false, nil // Skip empty thresholds
	}
	thresholdObj := ParseRangeString(rangeStr)
	if thresholdObj == nil {
		return false, fmt.Errorf("failed to parse range string %s: %w", rangeStr, ErrInvalidRangeThreshold)
	}
	return thresholdObj.CheckRange(value), nil
}
