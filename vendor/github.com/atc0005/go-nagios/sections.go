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
func (p Plugin) handleServiceOutputSection(w io.Writer) {
	if p.LongServiceOutput == "" {
		// If Long Service Output was not specified, explicitly trim any
		// formatted trailing spacing so that performance data output will be
		// emitted immediately following the Service Output on the same line.

		// NOTE: We explicitly include a space character in the cut set just
		// on the off chance that a future update to the CheckOutputEOL
		// constant removes the explicitly leading whitespace character.
		cutSet := fmt.Sprintf(" \t%s", CheckOutputEOL)
		p.ServiceOutput = strings.TrimRight(p.ServiceOutput, cutSet)
	}

	// Aside from (potentially) trimming trailing whitespace, we apply no
	// formatting changes to this content, simply emit it as-is. This helps
	// avoid potential issues with literal characters being interpreted as
	// formatting verbs.
	fmt.Fprint(w, p.ServiceOutput)
}

// handleErrorsSection is a wrapper around the logic used to handle/process
// the Errors section header and listing.
func (p Plugin) handleErrorsSection(w io.Writer) {

	// If one or more errors were recorded and client code has not opted to
	// hide the section ...
	if !p.isErrorsHidden() {

		fmt.Fprintf(w,
			"%s%s**%s**%s%s",
			CheckOutputEOL,
			CheckOutputEOL,
			p.getErrorsLabelText(),
			CheckOutputEOL,
			CheckOutputEOL,
		)

		if p.LastError != nil {
			fmt.Fprintf(w, "* %v%s", p.LastError, CheckOutputEOL)
		}

		// Process any non-nil errors in the collection.
		for _, err := range p.Errors {
			if err != nil {
				fmt.Fprintf(w, "* %v%s", err, CheckOutputEOL)
			}
		}

	}

}

// handleThresholdsSection is a wrapper around the logic used to
// handle/process the Thresholds section header and listing.
func (p Plugin) handleThresholdsSection(w io.Writer) {

	// We skip emitting the thresholds section if there isn't any
	// LongServiceOutput to process.
	if p.LongServiceOutput != "" {

		// If one or more threshold values were recorded and client code has
		// not opted to hide the section ...
		if !p.isThresholdsSectionHidden() {

			fmt.Fprintf(w,
				"%s**%s**%s%s",
				CheckOutputEOL,
				p.getThresholdsLabelText(),
				CheckOutputEOL,
				CheckOutputEOL,
			)

			if p.CriticalThreshold != "" {
				fmt.Fprintf(w,
					"* %s: %v%s",
					StateCRITICALLabel,
					p.CriticalThreshold,
					CheckOutputEOL,
				)
			}

			if p.WarningThreshold != "" {
				fmt.Fprintf(w,
					"* %s: %v%s",
					StateWARNINGLabel,
					p.WarningThreshold,
					CheckOutputEOL,
				)
			}
		}
	}

}

// handleLongServiceOutput is a wrapper around the logic used to
// handle/process the LongServiceOutput content.
func (p Plugin) handleLongServiceOutput(w io.Writer) {

	// Early exit if there is no content to emit.
	if p.LongServiceOutput == "" {
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
	case !p.isThresholdsSectionHidden() || !p.isErrorsHidden():
		fmt.Fprintf(w,
			"%s**%s**%s",
			CheckOutputEOL,
			p.getDetailedInfoLabelText(),
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
		p.LongServiceOutput,
		CheckOutputEOL,
	)
}

// handlePerformanceData is a wrapper around the logic used to
// handle/process plugin Performance Data.
func (p *Plugin) handlePerformanceData(w io.Writer) {

	// We require that a one-line summary is set by client code before
	// emitting performance data metrics.
	if strings.TrimSpace(p.ServiceOutput) == "" {
		return
	}

	// If the value is available, use it, otherwise this is a NOOP.
	p.tryAddDefaultTimeMetric()

	// If no metrics have been collected by this point we have nothing further
	// to do.
	if len(p.perfData) == 0 {
		return
	}

	// Performance data metrics are appended to plugin output. These
	// metrics are provided as a single line, leading with a pipe
	// character, a space and one or more metrics each separated from
	// another by a single space.
	fmt.Fprint(w, " |")

	// Sort performance data values prior to emitting them so that the
	// output is consistent across plugin execution.
	perfData := p.getSortedPerfData()

	for _, pd := range perfData {
		fmt.Fprint(w, pd.String())
	}

	// Add final trailing newline to satisfy Nagios plugin output format.
	fmt.Fprint(w, CheckOutputEOL)

}

// isThresholdsSectionHidden indicates whether the Thresholds section should
// be omitted from output.
func (p Plugin) isThresholdsSectionHidden() bool {
	if p.hideThresholdsSection || (p.WarningThreshold == "" && p.CriticalThreshold == "") {
		return true
	}
	return false
}

// isErrorsHidden indicates whether the Thresholds section should be omitted
// from output.
func (p Plugin) isErrorsHidden() bool {
	if p.hideErrorsSection || (len(p.Errors) == 0 && p.LastError == nil) {
		return true
	}
	return false
}

// getThresholdsLabelText retrieves the custom thresholds label text if set,
// otherwise returns the default value.
func (p Plugin) getThresholdsLabelText() string {
	switch {
	case p.thresholdsLabel != "":
		return p.thresholdsLabel
	default:
		return defaultThresholdsLabel
	}
}

// getErrorsLabelText retrieves the custom errors label text if set, otherwise
// returns the default value.
func (p Plugin) getErrorsLabelText() string {
	switch {
	case p.errorsLabel != "":
		return p.errorsLabel
	default:
		return defaultErrorsLabel
	}
}

// getErrorsLabelText retrieves the custom detailed info label text if set,
// otherwise returns the default value.
func (p Plugin) getDetailedInfoLabelText() string {
	switch {
	case p.detailedInfoLabel != "":
		return p.detailedInfoLabel
	default:
		return defaultDetailedInfoLabel
	}
}

// SetThresholdsLabel overrides the default thresholds label text.
func (p *Plugin) SetThresholdsLabel(newLabel string) {
	p.thresholdsLabel = newLabel
}

// SetErrorsLabel overrides the default errors label text.
func (p *Plugin) SetErrorsLabel(newLabel string) {
	p.errorsLabel = newLabel
}

// SetDetailedInfoLabel overrides the default detailed info label text.
func (p *Plugin) SetDetailedInfoLabel(newLabel string) {
	p.detailedInfoLabel = newLabel
}

// HideThresholdsSection indicates that client code has opted to hide the
// thresholds section, regardless of whether values were previously provided
// for display.
func (p *Plugin) HideThresholdsSection() {
	p.hideThresholdsSection = true
}

// HideErrorsSection indicates that client code has opted to hide the errors
// section, regardless of whether values were previously provided for display.
func (p *Plugin) HideErrorsSection() {
	p.hideErrorsSection = true
}

// getSortedPerfData returns a sorted copy of the performance data metrics.
func (p Plugin) getSortedPerfData() []PerformanceData {
	keys := make([]string, 0, len(p.perfData))
	perfData := make([]PerformanceData, 0, len(p.perfData))

	for k := range p.perfData {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		pd := p.perfData[key]
		perfData = append(perfData, pd)
	}

	return perfData
}
