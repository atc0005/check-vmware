// Copyright 2023 Codeweavers Ltd
// Copyright 2023 Adam Chalkley
//
// https://github.com/atc0005/go-nagios
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package nagios

import (
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
	switch {
	case !r.EndInfinity && !r.StartInfinity:
		if r.Start <= valueAsAFloat && valueAsAFloat <= r.End {
			return false
		}
		return true

	case !r.StartInfinity && r.EndInfinity:
		if valueAsAFloat >= r.Start {
			return false
		}
		return true

	case r.StartInfinity && !r.EndInfinity:
		if valueAsAFloat <= r.End {
			return false
		}
		return true

	default:
		return false
	}
}

// ParseRangeString static method to construct a Range object from the string
// representation based on the [Nagios Plugin Dev Guidelines: Threshold and
// Ranges] definition.
//
// [Nagios Plugin Dev Guidelines: Threshold and Ranges]: https://nagios-plugins.org/doc/guidelines.html#THRESHOLDFORMAT
func ParseRangeString(input string) *Range {
	r := Range{}

	digitOrInfinity := regexp.MustCompile(`[\d~]`)
	optionalInvertAndRange := regexp.MustCompile(`^\@?((?:[-+]?[\d\.]+)(?:e(?:[-+]?[\d\.]+))?|~)?(:((?:[-+]?[\d\.]+)(?:e(?:[-+]?[\d\.]+))?)?)?$`)
	firstHalfOfRange := regexp.MustCompile(`^((?:[-+]?[\d\.]+)(?:e(?:[-+]?[\d\.]+))?)?:`)
	endOfRange := regexp.MustCompile(`^(?:[-+]?[\d\.]+)(?:e(?:[-+]?[\d\.]+))?$`)

	r.Start = 0
	r.StartInfinity = false
	r.End = 0
	r.EndInfinity = false
	r.AlertOn = "OUTSIDE"

	valid := true

	// If regex does not match ...
	if !(digitOrInfinity.MatchString(input) && optionalInvertAndRange.MatchString(input)) {
		return nil
	}

	// Invert the range.
	//
	// i.e. @10:20 means ≥ 10 and ≤ 20 (inside the range of {10 .. 20}
	// inclusive)
	if strings.HasPrefix(input, "@") {
		r.AlertOn = "INSIDE"
		input = input[1:]
	}

	// ~ represents infinity
	if strings.HasPrefix(input, "~") {
		r.StartInfinity = true
		input = input[1:]
	}

	// 10:
	rangeComponents := firstHalfOfRange.FindAllStringSubmatch(input, -1)
	if rangeComponents != nil {
		if rangeComponents[0][1] != "" {
			r.Start, _ = strconv.ParseFloat(rangeComponents[0][1], 64)
			r.StartInfinity = false
		}

		r.EndInfinity = true
		input = strings.TrimPrefix(input, rangeComponents[0][0])
		valid = true
	}

	// x:10 or 10
	endOfRangeComponents := endOfRange.FindAllStringSubmatch(input, -1)
	if endOfRangeComponents != nil {

		r.End, _ = strconv.ParseFloat(endOfRangeComponents[0][0], 64)
		r.EndInfinity = false
		valid = true
	}

	if valid && (r.StartInfinity || r.EndInfinity || r.Start <= r.End) {
		return &r
	}

	return nil
}

// EvaluateThreshold causes the performance data to be checked against the
// Warn and Crit thresholds provided by client code and sets the
// ExitStatusCode of the plugin as appropriate.
func (p *Plugin) EvaluateThreshold(perfData ...PerformanceData) error {
	for i := range perfData {

		if perfData[i].Crit != "" {

			CriticalThresholdObject := ParseRangeString(perfData[i].Crit)

			if CriticalThresholdObject.CheckRange(perfData[i].Value) {
				p.ExitStatusCode = StateCRITICALExitCode
				return nil
			}
		}

		if perfData[i].Warn != "" {
			warningThresholdObject := ParseRangeString(perfData[i].Warn)

			if warningThresholdObject.CheckRange(perfData[i].Value) {
				p.ExitStatusCode = StateWARNINGExitCode
				return nil
			}
		}
	}

	return nil
}
