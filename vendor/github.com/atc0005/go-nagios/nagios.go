// Copyright 2020 Adam Chalkley
//
// https://github.com/atc0005/go-nagios
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

// Package nagios is a small collection of common types and package-level
// variables intended for use with various plugins to reduce code duplication.
package nagios

import (
	"fmt"
	"os"
	"runtime"
)

// Nagios plugin/service check states. These constants replicate the values
// from utils.sh which is normally found at one of these two locations,
// depending on which Linux distribution you're using:
//
//     /usr/lib/nagios/plugins/utils.sh
//     /usr/local/nagios/libexec/utils.sh
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
	StateUNKNOWNLabel   string = "UKNOWN"
	StateDEPENDENTLabel string = "DEPENDENT"
)

// CheckOutputEOL is the newline used with formatted service and host check
// output. Based on previous testing, Nagios treats LF newlines within the
// `$LONGSERVICEOUTPUT$` macro as literal values instead of parsing them for
// display purposes. Using DOS EOL values with `fmt.Printf()` gives the
// results that we're looking for with that output and (presumably) host check
// output as well.
const CheckOutputEOL string = "\r\n"

// ExitCallBackFunc represents a function that is called as a final step
// before application termination so that branding information can be emitted
// for inclusion in the notification. This helps identify which specific
// application (and its version) that is responsible for the notification.
type ExitCallBackFunc func() string

// ExitState represents the last known execution state of the
// application, including the most recent error and the final intended plugin
// state.
type ExitState struct {

	// LastError is the last error encountered which should be reported as
	// part of ending the service check (e.g., "Failed to connect to XYZ to
	// check contents of Inbox").
	LastError error

	// LastError is the exit or exit status code provided to the Nagios
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

	// WarningThreshold is the value used to determine when the service check
	// has crossed between an existing state into a WARNING state. This value
	// is used for display purposes.
	WarningThreshold string

	// CriticalThreshold is the value used to determine when the service check
	// has crossed between an existing state into a CRITICAL state. This value
	// is used for display purposes.
	CriticalThreshold string

	// BrandingCallback is a function that is called before application
	// termination to emit branding details at the end of the notification.
	// See also ExitCallBackFunc.
	BrandingCallback ExitCallBackFunc
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
func (es *ExitState) ReturnCheckResults() {

	// Check for unhandled panic in client code. If present, override
	// ExitState and make clear that the client code/plugin crashed.
	if err := recover(); err != nil {

		es.LastError = fmt.Errorf("plugin crash/panic detected")

		es.ServiceOutput = fmt.Sprintf(
			"%s: plugin crash detected. See details via web UI or run plugin manually via CLI.",
			StateCRITICALLabel,
		)

		// Gather stack trace associated with panic for display.
		// NOTE: runtime.Stack *requires* that we preallocate the slice and
		// with zero values.
		stackTrace := make([]byte, 512)
		stackTraceSize := runtime.Stack(stackTrace, false)

		// Using literal `<pre>` tags in an effort to wrap the stack trace
		// details so that they are not interpreted as formatting characters
		// when passed through web UI, text, email, Teams, etc.
		es.LongServiceOutput = fmt.Sprintf(
			"<pre>%s%s%s</pre>",
			CheckOutputEOL,
			stackTrace[0:stackTraceSize],
			CheckOutputEOL,
		)

		es.ExitStatusCode = StateCRITICALExitCode

	}

	// ##################################################################
	// Note: fmt.Println() has the same issue as `\n`: Nagios seems to
	// interpret them literally instead of emitting an actual newline.
	// We work around that by using fmt.Printf() for output that is
	// intended for display within the Nagios web UI.
	// ##################################################################

	// One-line output used as the summary or short explanation for the
	// specific Nagios state that we are returning.
	fmt.Printf(es.ServiceOutput)

	if es.LongServiceOutput != "" || es.LastError != nil {

		fmt.Printf("%s%s**ERRORS**%s", CheckOutputEOL, CheckOutputEOL, CheckOutputEOL)

		// If an error occurred or if there are additional details to share ...

		if es.LastError != nil {
			fmt.Printf("%s* %v%s", CheckOutputEOL, es.LastError, CheckOutputEOL)
		} else {
			fmt.Printf("%s* None%s", CheckOutputEOL, CheckOutputEOL)
		}

		if es.LongServiceOutput != "" {

			fmt.Printf("%s**THRESHOLDS**%s", CheckOutputEOL, CheckOutputEOL)

			if es.CriticalThreshold != "" || es.WarningThreshold != "" {

				fmt.Print(CheckOutputEOL)

				if es.CriticalThreshold != "" {
					fmt.Printf(
						"* %s: %v%s",
						StateCRITICALLabel,
						es.CriticalThreshold,
						CheckOutputEOL,
					)
				}

				if es.WarningThreshold != "" {
					fmt.Printf(
						"* %s: %v%s",
						StateWARNINGLabel,
						es.WarningThreshold,
						CheckOutputEOL,
					)
				}
			} else {
				fmt.Printf("%s* Not specified%s", CheckOutputEOL, CheckOutputEOL)
			}

			fmt.Printf("%s**DETAILED INFO**%s", CheckOutputEOL, CheckOutputEOL)

			// Note: fmt.Println() has the same issue as `\n`: Nagios seems to
			// interpret them literally instead of emitting an actual newline.
			// We work around that by using fmt.Printf() for output that is
			// intended for display within the Nagios web UI.
			fmt.Printf(
				"%s%v%s",
				CheckOutputEOL,
				es.LongServiceOutput,
				CheckOutputEOL,
			)
		}

	}

	// If set, call user-provided branding function just before exiting
	// application
	if es.BrandingCallback != nil {
		fmt.Printf("%s%s%s", CheckOutputEOL, es.BrandingCallback(), CheckOutputEOL)
	}

	os.Exit(es.ExitStatusCode)
}
