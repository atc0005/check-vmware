// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/go-nagios
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package nagios

import (
	"errors"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

// Nagios plugin/service check states. These constants replicate the values
// from utils.sh which is normally found at one of these two locations,
// depending on which Linux distribution you're using:
//
// - /usr/lib/nagios/plugins/utils.sh
// - /usr/local/nagios/libexec/utils.sh
//
// See also http://nagios-plugins.org/doc/guidelines.html
const (
	StateOKExitCode        int = 0
	StateWARNINGExitCode   int = 1
	StateCRITICALExitCode  int = 2
	StateUNKNOWNExitCode   int = 3
	StateDEPENDENTExitCode int = 4
)

// Nagios plugin/service check state "labels". These constants are provided as
// an alternative to using literal state strings throughout client application
// code.
const (
	StateOKLabel        string = "OK"
	StateWARNINGLabel   string = "WARNING"
	StateCRITICALLabel  string = "CRITICAL"
	StateUNKNOWNLabel   string = "UNKNOWN"
	StateDEPENDENTLabel string = "DEPENDENT"
)

// CheckOutputEOL is the newline character(s) used with formatted service and
// host check output. Based on previous testing, Nagios treats LF newlines
// (without a leading space) within the `$LONGSERVICEOUTPUT$` macro as literal
// values instead of parsing them for display purposes.
//
// Using DOS EOL values with fmt.Fprintf() (or fmt.Fprintln()) gives expected
// formatting results in the Nagios Core web UI, but results in double
// newlines in Nagios XI output (see GH-109). Using a UNIX EOL with a single
// leading space appears to give the intended results for both Nagios Core and
// Nagios XI.
const CheckOutputEOL string = " \n"

// Default header text for various sections of the output if not overridden.
const (
	defaultThresholdsLabel   string = "THRESHOLDS"
	defaultErrorsLabel       string = "ERRORS"
	defaultDetailedInfoLabel string = "DETAILED INFO"
)

// Default performance data metrics emitted if not specified by client code.
const (
	defaultTimeMetricLabel             string = "time"
	defaultTimeMetricUnitOfMeasurement string = "ms"
)

// Sentinel error collection. Exported for potential use by client code to
// detect & handle specific error scenarios.
var (
	// ErrPanicDetected indicates that client code has an unhandled panic and
	// that this library detected it before it could cause the plugin to
	// abort. This error is included in the LongServiceOutput emitted by the
	// plugin.
	ErrPanicDetected = errors.New("plugin crash/panic detected")

	// ErrPerformanceDataMissingLabel indicates that client code did not
	// provide a PerformanceData value in the expected format; the label for
	// the label/value pair is missing.
	ErrPerformanceDataMissingLabel = errors.New("provided performance data missing required label")

	// ErrPerformanceDataMissingValue indicates that client code did not
	// provide a PerformanceData value in the expected format; the value for
	// the label/value pair is missing.
	ErrPerformanceDataMissingValue = errors.New("provided performance data missing required value")

	// ErrNoPerformanceDataProvided indicates that client code did not provide
	// the expected PerformanceData value(s).
	ErrNoPerformanceDataProvided = errors.New("no performance data provided")

	// ErrInvalidPerformanceDataFormat indicates that a given performance data
	// metric is not in a supported format.
	ErrInvalidPerformanceDataFormat = errors.New("invalid performance data format")

	// TODO: Should we use field-specific errors or is the more general
	// ErrInvalidPerformanceDataFormat "good enough" ? Wrapped versions of
	// that error will likely already indicate which field is a problem, but
	// that approach won't lend itself to automatic behavior regarding issues
	// with a specific field.
	//
	// ErrInvalidPerformanceDataLabelField = errors.New("invalid field Label in parsed performance data")
	// ErrInvalidPerformanceDataValueField = errors.New("invalid field Value in parsed performance data")
	// ErrInvalidPerformanceDataUoMField   = errors.New("invalid field UnitOfMeasurement in parsed performance data")
	// ErrInvalidPerformanceDataWarnField  = errors.New("invalid field Warn in parsed performance data")
	// ErrInvalidPerformanceDataCritField  = errors.New("invalid field Crit in parsed performance data")
	// ErrInvalidPerformanceDataMinField   = errors.New("invalid field Min in parsed performance data")
	// ErrInvalidPerformanceDataMaxField   = errors.New("invalid field Max in parsed performance data")
)

// ServiceState represents the status label and exit code for a service check.
type ServiceState struct {

	// Label maps directly to one of the supported Nagios state labels.
	Label string

	// ExitCode is the exit or exit status code associated with a Nagios
	// service check.
	ExitCode int
}

// ExitCallBackFunc represents a function that is called as a final step
// before application termination so that branding information can be emitted
// for inclusion in the notification. This helps identify which specific
// application (and its version) that is responsible for the notification.
type ExitCallBackFunc func() string

// Plugin represents the state of a monitoring plugin, including the most
// recent error and the final intended plugin state.
type Plugin struct {
	// outputSink is the user-specified or fallback target for plugin output.
	outputSink io.Writer

	// start tracks when the associated plugin begins executing. This value is
	// used to generate a default `time` performance data metric (which can be
	// overridden by client code).
	start time.Time

	// LastError is the last error encountered which should be reported as
	// part of ending the service check (e.g., "Failed to connect to XYZ to
	// check contents of Inbox").
	//
	// Deprecated: Use Errors field or AddError method instead.
	LastError error

	// Errors is a collection of one or more recorded errors to be displayed
	// in LongServiceOutput as a list when ending the service check.
	Errors []error

	// ExitStatusCode is the exit or exit status code provided to the Nagios
	// instance that calls this service check. These status codes indicate to
	// Nagios "state" the service is considered to be in. The most common
	// states are OK (0), WARNING (1) and CRITICAL (2).
	ExitStatusCode int

	// ServiceOutput is the first line of text output from the last service
	// check (i.e. "Ping OK").
	ServiceOutput string

	// LongServiceOutput is the full text output (aside from the first line)
	// from the last service check.
	LongServiceOutput string

	// perfData is the collection of zero or more PerformanceData values
	// generated by the plugin. Each entry in the collection is unique.
	perfData map[string]PerformanceData

	// WarningThreshold is the value used to determine when the service check
	// has crossed between an existing state into a WARNING state. This value
	// is used for display purposes.
	WarningThreshold string

	// CriticalThreshold is the value used to determine when the service check
	// has crossed between an existing state into a CRITICAL state. This value
	// is used for display purposes.
	CriticalThreshold string

	// thresholdLabel is an optional custom label used in place of the
	// standard text prior to a list of threshold values.
	thresholdsLabel string

	// errorsLabel is an optional custom label used in place of the standard
	// text prior to a list of recorded error values.
	errorsLabel string

	// detailedInfoLabel is an optional custom label used in place of the
	// standard text prior to emitting LongServiceOutput.
	detailedInfoLabel string

	// hideThresholdsSection indicates whether client code has opted to hide
	// the thresholds section, regardless of whether client code previously
	// specified values for display.
	hideThresholdsSection bool

	// hideErrorsSection indicates whether client code has opted to hide the
	// errors section, regardless of whether client code previously specified
	// values for display.
	hideErrorsSection bool

	// shouldSkipOSExit is intended to support tests where actually performing
	// the final os.Exit(x) call results in a panic (Go 1.16+). If set,
	// calling os.Exit(x) is skipped and a message is logged to os.Stderr
	// instead.
	shouldSkipOSExit bool

	// BrandingCallback is a function that is called before application
	// termination to emit branding details at the end of the notification.
	// See also ExitCallBackFunc.
	BrandingCallback ExitCallBackFunc
}

// NewPlugin constructs a new Plugin value in the same way that client code
// has been using this library. We also record a default time performance data
// metric. This default metric is ignored if supplied by client code.
func NewPlugin() *Plugin {
	es := Plugin{
		start:          time.Now(),
		LastError:      nil,
		ExitStatusCode: StateOKExitCode,
	}

	return &es
}

// ReturnCheckResults is intended to provide a reliable way to return a
// desired exit code from applications used as Nagios plugins. In most cases,
// this method should be registered as the first deferred function in client
// code. See remarks regarding "masking" or "swallowing" application panics.
//
// Since Nagios relies on plugin exit codes to determine success/failure of
// checks, the approach that is most often used with other languages is to use
// something like Using os.Exit() directly and force an early exit of the
// application with an explicit exit code. Using os.Exit() directly in Go does
// not run deferred functions. Go-based plugins that do not rely on deferring
// function calls may be able to use os.Exit(), but introducing new
// dependencies later could introduce problems if those dependencies rely on
// deferring functions.
//
// Before calling this method, client code should first set appropriate field
// values on the receiver. When called, this method will process them and exit
// with the desired exit code and status output.
//
// To repeat, if scheduled via defer, this method should be registered first;
// because this method calls os.Exit to set the intended plugin exit state, no
// other deferred functions will have an opportunity to run, so register this
// method first so that when deferred, it will be run last (FILO).
//
// Because this method is (or should be) deferred first within client code, it
// will run after all other deferred functions. It will also run before a
// panic in client code forces the application to exit. As already noted, this
// method calls os.Exit to set the plugin exit state. Because os.Exit forces
// the application to terminate immediately without running other deferred
// functions or processing panics, this "masks", "swallows" or "blocks" panics
// from client code from surfacing. This method checks for unhandled panics
// and if found, overrides exit state details from client code and surfaces
// details from the panic instead as a CRITICAL state.
func (p *Plugin) ReturnCheckResults() {

	var output strings.Builder

	// ##################################################################
	// Note: fmt.Println() (and fmt.Fprintln()) has the same issue as `\n`:
	// Nagios seems to interpret them literally instead of emitting an actual
	// newline. We work around that by using fmt.Fprintf() and fmt.Fprint()
	// for output that is intended for display within the Nagios web UI.
	// ##################################################################

	// Check for unhandled panic in client code. If present, override
	// Plugin and make clear that the client code/plugin crashed.
	if err := recover(); err != nil {

		p.AddError(fmt.Errorf("%w: %s", ErrPanicDetected, err))

		p.ServiceOutput = fmt.Sprintf(
			"%s: plugin crash detected. See details via web UI or run plugin manually via CLI.",
			StateCRITICALLabel,
		)

		// Gather stack trace associated with panic.
		stackTrace := debug.Stack()

		// Wrap stack trace details in an attempt to prevent these details
		// from being interpreted as formatting characters when passed through
		// web UI, text, email, Teams, etc. We use Markdown fenced code blocks
		// instead of `<pre>` start/end tags because Nagios strips out angle
		// brackets (due to default `illegal_macro_output_chars` settings).
		p.LongServiceOutput = fmt.Sprintf(
			"```%s%s%s%s%s%s```",
			CheckOutputEOL,
			err,
			CheckOutputEOL,
			CheckOutputEOL,
			stackTrace,
			CheckOutputEOL,
		)

		p.ExitStatusCode = StateCRITICALExitCode

	}

	p.handleServiceOutputSection(&output)

	p.handleErrorsSection(&output)

	p.handleThresholdsSection(&output)

	p.handleLongServiceOutput(&output)

	// If set, call user-provided branding function before emitting
	// performance data and exiting application.
	if p.BrandingCallback != nil {
		fmt.Fprintf(&output, "%s%s%s", CheckOutputEOL, p.BrandingCallback(), CheckOutputEOL)
	}

	p.handlePerformanceData(&output)

	// Emit all collected plugin output using user-specified or fallback
	// output target.
	p.emitOutput(output.String())

	// TODO: Should we offer an option to redirect the log message to stderr
	// to another error output sink?
	//
	// TODO: Perhaps just don't emit anything at all?
	switch {
	case p.shouldSkipOSExit:
		fmt.Fprintln(os.Stderr, "Skipping os.Exit call as requested.")
	default:
		os.Exit(p.ExitStatusCode)
	}
}

// AddPerfData adds provided performance data to the collection overwriting
// any previous performance data metrics using the same label.
//
// Validation is skipped if requested, otherwise an error is returned if
// validation fails. Validation failure results in no performance data being
// appended. Client code may wish to disable validation if performing this
// step directly.
func (p *Plugin) AddPerfData(skipValidate bool, perfData ...PerformanceData) error {

	if len(perfData) == 0 {
		return ErrNoPerformanceDataProvided
	}

	if !skipValidate {
		for i := range perfData {
			if err := perfData[i].Validate(); err != nil {
				return err
			}
		}
	}

	if p.perfData == nil {
		p.perfData = make(map[string]PerformanceData)
	}

	for _, pd := range perfData {
		p.perfData[strings.ToLower(pd.Label)] = pd
	}

	return nil
}

// AddError appends provided errors to the collection.
func (p *Plugin) AddError(err ...error) {
	p.Errors = append(p.Errors, err...)
}

// SetOutputTarget assigns a target for Nagios plugin output. By default
// output is emitted to os.Stdout.
func (p *Plugin) SetOutputTarget(w io.Writer) {
	// Guard against potential nil argument.
	if w == nil {
		p.outputSink = os.Stdout
	}

	p.outputSink = w
}

// SkipOSExit indicates that the os.Exit(x) step used to signal to Nagios what
// state plugin execution has completed in (e.g., OK, WARNING, ...) should be
// skipped. If skipped, a message is logged to os.Stderr in place of the
// os.Exit(x) call.
//
// Disabling the call to os.Exit is needed by tests to prevent panics in Go
// 1.16 and newer.
func (p *Plugin) SkipOSExit() {
	p.shouldSkipOSExit = true
}

// emitOutput writes final plugin output to the previously set output target.
// No further modifications to plugin output are performed.
func (p Plugin) emitOutput(pluginOutput string) {

	// Emit all collected output using user-specified output target. Fall back
	// to standard output if not set.
	if p.outputSink == nil {
		p.outputSink = os.Stdout
	}

	fmt.Fprint(p.outputSink, pluginOutput)
}

// tryAddDefaultTimeMetric inserts a default `time` performance data metric
// into the collection IF client code has not already specified such a value
// AND we have a non-zero start value to use.
func (p *Plugin) tryAddDefaultTimeMetric() {

	// We already have an existing time metric, skip replacing it.
	if _, hasTimeMetric := p.perfData[defaultTimeMetricLabel]; hasTimeMetric {
		return
	}

	// Our Plugin value was not generated from the constructor, so we do
	// not have an internal plugin start time that we can use to generate a
	// default time metric.
	if p.start.IsZero() {
		return
	}

	if p.perfData == nil {
		p.perfData = make(map[string]PerformanceData)
	}

	p.perfData[defaultTimeMetricLabel] = defaultTimeMetric(p.start)
}

// defaultTimeMetric is a helper function that wraps the logic used to provide
// a default performance data metric that tracks plugin execution time.
func defaultTimeMetric(start time.Time) PerformanceData {
	return PerformanceData{
		Label:             defaultTimeMetricLabel,
		Value:             fmt.Sprintf("%d", time.Since(start).Milliseconds()),
		UnitOfMeasurement: defaultTimeMetricUnitOfMeasurement,
	}
}

// SupportedStateLabels returns a list of valid plugin state labels.
func SupportedStateLabels() []string {
	return []string{
		StateOKLabel,
		StateWARNINGLabel,
		StateCRITICALLabel,
		StateUNKNOWNLabel,
		StateDEPENDENTLabel,
	}
}

// SupportedExitCodes returns a list of valid plugin exit codes.
func SupportedExitCodes() []int {
	return []int{
		StateOKExitCode,
		StateWARNINGExitCode,
		StateCRITICALExitCode,
		StateUNKNOWNExitCode,
		StateDEPENDENTExitCode,
	}
}

// SupportedServiceStates returns a collection of valid plugin service states.
func SupportedServiceStates() []ServiceState {
	return []ServiceState{
		{
			Label:    StateOKLabel,
			ExitCode: StateOKExitCode,
		},
		{
			Label:    StateWARNINGLabel,
			ExitCode: StateWARNINGExitCode,
		},
		{
			Label:    StateCRITICALLabel,
			ExitCode: StateCRITICALExitCode,
		},
		{
			Label:    StateUNKNOWNLabel,
			ExitCode: StateUNKNOWNExitCode,
		},
		{
			Label:    StateDEPENDENTLabel,
			ExitCode: StateDEPENDENTExitCode,
		},
	}
}

// StateLabelToExitCode returns the corresponding plugin exit code for the
// given plugin state label. If an invalid value is provided the
// StateUNKNOWNExitCode value is returned.
func StateLabelToExitCode(label string) int {
	switch strings.ToUpper(label) {
	case StateOKLabel:
		return StateOKExitCode
	case StateWARNINGLabel:
		return StateWARNINGExitCode
	case StateCRITICALLabel:
		return StateCRITICALExitCode
	case StateUNKNOWNLabel:
		return StateUNKNOWNExitCode
	case StateDEPENDENTLabel:
		return StateDEPENDENTExitCode
	default:
		return StateUNKNOWNExitCode
	}
}

// ExitCodeToStateLabel returns the corresponding plugin state label for the
// given plugin exit code. If an invalid value is provided the
// StateUNKNOWNLabel value is returned.
func ExitCodeToStateLabel(exitCode int) string {
	switch exitCode {
	case StateOKExitCode:
		return StateOKLabel
	case StateWARNINGExitCode:
		return StateWARNINGLabel
	case StateCRITICALExitCode:
		return StateCRITICALLabel
	case StateUNKNOWNExitCode:
		return StateUNKNOWNLabel
	case StateDEPENDENTExitCode:
		return StateDEPENDENTLabel
	default:
		return StateUNKNOWNLabel
	}
}
