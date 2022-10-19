// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/go-nagios
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package nagios

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// handleServiceOutputSection is a wrapper around the logic used to process
// the Service Output or "one-line summary" content.
func (es ExitState) handleServiceOutputSection(w io.Writer) {
	if es.LongServiceOutput == "" {
		// If Long Service Output was not specified, explicitly trim any
		// formatted trailing spacing so that performance data output will be
		// emitted immediately following the Service Output on the same line.

		// NOTE: We explicitly include a space character in the cut set just
		// on the off chance that a future update to the CheckOutputEOL
		// constant removes the explicitly leading whitespace character.
		cutSet := fmt.Sprintf(" \t%s", CheckOutputEOL)
		es.ServiceOutput = strings.TrimRight(es.ServiceOutput, cutSet)
	}

	// Aside from (potentially) trimming trailing whitespace, we apply no
	// formatting changes to this content, simply emit it as-is. This helps
	// avoid potential issues with literal characters being interpreted as
	// formatting verbs.
	fmt.Fprint(w, es.ServiceOutput)
}

// handleErrorsSection is a wrapper around the logic used to handle/process the
// Errors section header and listing.
func (es ExitState) handleErrorsSection(w io.Writer) {

	// If one or more errors were recorded and client code has not opted to
	// hide the section ...
	if !es.isErrorsHidden() {

		fmt.Fprintf(w,
			"%s%s**%s**%s%s",
			CheckOutputEOL,
			CheckOutputEOL,
			es.getErrorsLabelText(),
			CheckOutputEOL,
			CheckOutputEOL,
		)

		if es.LastError != nil {
			fmt.Fprintf(w, "* %v%s", es.LastError, CheckOutputEOL)
		}

		// Process any non-nil errors in the collection.
		for _, err := range es.Errors {
			if err != nil {
				fmt.Fprintf(w, "* %v%s", err, CheckOutputEOL)
			}
		}

	}

}

// handleThresholdsSection is a wrapper around the logic used to handle/process the
// Thresholds section header and listing.
func (es ExitState) handleThresholdsSection(w io.Writer) {

	// We skip emitting the thresholds section if there isn't any
	// LongServiceOutput to process.
	if es.LongServiceOutput != "" {

		// If one or more threshold values were recorded and client code has
		// not opted to hide the section ...
		if !es.isThresholdsSectionHidden() {

			fmt.Fprintf(w,
				"%s**%s**%s%s",
				CheckOutputEOL,
				es.getThresholdsLabelText(),
				CheckOutputEOL,
				CheckOutputEOL,
			)

			if es.CriticalThreshold != "" {
				fmt.Fprintf(w,
					"* %s: %v%s",
					StateCRITICALLabel,
					es.CriticalThreshold,
					CheckOutputEOL,
				)
			}

			if es.WarningThreshold != "" {
				fmt.Fprintf(w,
					"* %s: %v%s",
					StateWARNINGLabel,
					es.WarningThreshold,
					CheckOutputEOL,
				)
			}
		}
	}

}

// handleLongServiceOutput is a wrapper around the logic used to
// handle/process the LongServiceOutput content.
func (es ExitState) handleLongServiceOutput(w io.Writer) {

	// Early exit if there is no content to emit.
	if es.LongServiceOutput == "" {
		return
	}

	// Hide section header/label if threshold and error values were not
	// specified by client code or if client code opted to explicitly hide
	// those sections; there is no need to use a header to separate the
	// LongServiceOutput from those sections if they are not displayed.
	//
	// If we hide the section header, we still provide some padding to
	// prevent the LongServiceOutput from running up against the
	// ServiceOutput content.
	switch {
	case !es.isThresholdsSectionHidden() || !es.isErrorsHidden():
		fmt.Fprintf(w,
			"%s**%s**%s",
			CheckOutputEOL,
			es.getDetailedInfoLabelText(),
			CheckOutputEOL,
		)
	default:
		fmt.Fprint(w, CheckOutputEOL)
	}

	// Note: fmt.Println() (and fmt.Fprintln()) has the same issue as `\n`:
	// Nagios seems to interpret them literally instead of emitting an actual
	// newline. We work around that by using fmt.Fprintf() for output that is
	// intended for display within the Nagios web UI.
	fmt.Fprintf(w,
		"%s%v%s",
		CheckOutputEOL,
		es.LongServiceOutput,
		CheckOutputEOL,
	)
}

// handlePerformanceData is a wrapper around the logic used to
// handle/process plugin Performance Data.
func (es ExitState) handlePerformanceData(w io.Writer) {
	// Generate formatted performance data if provided. Only emit if a
	// one-line summary is set by client code.
	if len(es.perfData) != 0 && es.ServiceOutput != "" {

		// Performance data metrics are appended to plugin output. These
		// metrics are provided as a single line, leading with a pipe
		// character, a space and one or more metrics each separated from
		// another by a single space.
		fmt.Fprint(w, " |")

		// Sort performance data values prior to emitting them so that the
		// output is consistent across plugin execution.
		perfData := es.getSortedPerfData()

		for _, pd := range perfData {
			fmt.Fprintf(w,
				// The expected format of a performance data metric:
				//
				// 'label'=value[UOM];[warn];[crit];[min];[max]
				//
				// References:
				//
				// https://nagios-plugins.org/doc/guidelines.html
				// https://assets.nagios.com/downloads/nagioscore/docs/nagioscore/3/en/perfdata.html
				// https://assets.nagios.com/downloads/nagioscore/docs/nagioscore/3/en/pluginapi.html
				// https://www.monitoring-plugins.org/doc/guidelines.html
				// https://icinga.com/docs/icinga-2/latest/doc/05-service-monitoring/#performance-data-metrics
				" '%s'=%s%s;%s;%s;%s;%s",
				pd.Label,
				pd.Value,
				pd.UnitOfMeasurement,
				pd.Warn,
				pd.Crit,
				pd.Min,
				pd.Max,
			)
		}

		// Add final trailing newline to satisfy Nagios plugin output format.
		fmt.Fprint(w, CheckOutputEOL)

	}
}

// isThresholdsSectionHidden indicates whether the Thresholds section should
// be omitted from output.
func (es ExitState) isThresholdsSectionHidden() bool {
	if es.hideThresholdsSection || (es.WarningThreshold == "" && es.CriticalThreshold == "") {
		return true
	}
	return false
}

// isErrorsHidden indicates whether the Thresholds section should be omitted
// from output.
func (es ExitState) isErrorsHidden() bool {
	if es.hideErrorsSection || (len(es.Errors) == 0 && es.LastError == nil) {
		return true
	}
	return false
}

// getThresholdsLabelText retrieves the custom thresholds label text if set,
// otherwise returns the default value.
func (es ExitState) getThresholdsLabelText() string {
	switch {
	case es.thresholdsLabel != "":
		return es.thresholdsLabel
	default:
		return defaultThresholdsLabel
	}
}

// getErrorsLabelText retrieves the custom errors label text if set, otherwise
// returns the default value.
func (es ExitState) getErrorsLabelText() string {
	switch {
	case es.errorsLabel != "":
		return es.errorsLabel
	default:
		return defaultErrorsLabel
	}
}

// getErrorsLabelText retrieves the custom detailed info label text if set,
// otherwise returns the default value.
func (es ExitState) getDetailedInfoLabelText() string {
	switch {
	case es.detailedInfoLabel != "":
		return es.detailedInfoLabel
	default:
		return defaultDetailedInfoLabel
	}
}

// SetThresholdsLabel overrides the default thresholds label text.
func (es *ExitState) SetThresholdsLabel(newLabel string) {
	es.thresholdsLabel = newLabel
}

// SetErrorsLabel overrides the default errors label text.
func (es *ExitState) SetErrorsLabel(newLabel string) {
	es.errorsLabel = newLabel
}

// SetDetailedInfoLabel overrides the default detailed info label text.
func (es *ExitState) SetDetailedInfoLabel(newLabel string) {
	es.detailedInfoLabel = newLabel
}

// HideThresholdsSection indicates that client code has opted to hide the
// thresholds section, regardless of whether values were previously provided
// for display.
func (es *ExitState) HideThresholdsSection() {
	es.hideThresholdsSection = true
}

// HideErrorsSection indicates that client code has opted to hide the errors
// section, regardless of whether values were previously provided for display.
func (es *ExitState) HideErrorsSection() {
	es.hideErrorsSection = true
}

// getSortedPerfData returns a sorted copy of the performance data metrics.
func (es ExitState) getSortedPerfData() []PerformanceData {
	keys := make([]string, 0, len(es.perfData))
	perfData := make([]PerformanceData, 0, len(es.perfData))

	for k := range es.perfData {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		pd := es.perfData[key]
		perfData = append(perfData, pd)
	}

	return perfData
}
