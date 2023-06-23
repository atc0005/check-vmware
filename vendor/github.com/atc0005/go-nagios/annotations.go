// Copyright 2023 Adam Chalkley
//
// https://github.com/atc0005/go-nagios
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package nagios

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"syscall"
)

// runtimeTimeoutReachedAdvice offers advice to the sysadmin for remediating
// plugin timeouts if occurrence is routine.
const runtimeTimeoutReachedAdvice string = "plugin runtime exceeded specified timeout value; consider increasing value if this is routinely encountered"

// connectionResetByPeerAdvice offers advice to the sysadmin for remediating
// "read: connection reset by peer" errors. The cause of this error can vary
// but is often associated access control mechanisms such as firewalls,
// misconfigured IIS sites (e.g., missing certificate binding), overloaded
// services at maximum connection capacity or otherwise misconfigured services.
//
// See also https://stackoverflow.com/questions/1434451
const connectionResetByPeerAdvice string = "consider checking firewall, certificate/port bindings or maximum supported connections for remote service"

// connectionRefusedAdvice offers advice to the sysadmin to check the
// specified port and remote service state. The "connect: connection refused"
// error is often encountered when an application associated with the
// certificate being checked is stopped (e.g., troubleshooting purposes,
// replacing a certificate or the service has crashed) or if the wrong
// port was specified for a service.
const connectionRefusedAdvice string = "consider double-checking specified port and remote service state (i.e., make sure service is actually running on given port)"

// ErrorAnnotationMappings is a collection of error to string values. Each
// error is linked to specific advice for how to remediate the issue. The
// advice is appended to the list of errors (if any) which occurred during
// plugin execution. If the sysadmin opted to hide the errors section then no
// error output (and no advice) will be displayed.
type ErrorAnnotationMappings map[error]string

func isNilErrCollection(collection []error) bool {
	if len(collection) != 0 {
		for _, err := range collection {
			if err != nil {
				return false
			}
		}
	}
	return true
}

// DefaultErrorAnnotationMappings provides a default collection of error to
// string values which associate suggested advice with known/common error
// conditions.
//
// This collection is intended to serve as a starting point for plugin authors
// to further extend or override as needed.
func DefaultErrorAnnotationMappings() ErrorAnnotationMappings {
	return ErrorAnnotationMappings{
		context.DeadlineExceeded: runtimeTimeoutReachedAdvice,
		syscall.ECONNRESET:       connectionResetByPeerAdvice,
		syscall.ECONNREFUSED:     connectionRefusedAdvice,
	}
}

// AnnotateError adds additional human-readable explanation for errors
// encountered during plugin execution. This additional context is intended to
// help sysadmins remediate common issues detected by plugins.
//
// This function receives an optional map of error type to advice and one or
// more errors. If the map argument is nil or empty a default advice map is
// used. If an empty error collection or a collection of nil error values are
// provided for evaluation then nil is returned.
//
// Each error is evaluated for a match in its chain within the given advice
// map. If no advice map is given then the default advice map is used.
//
// If an error match is found then the advice is appended to the error via
// error wrapping. This process is repeated for each applicable error in the
// given error chain. The error is unmodified if no advice is found or if the
// error is already annotated with advice for the specific error.
//
// This updated error collection is returned to the caller.
//
// The original error collection is returned unmodified if no annotations were
// deemed necessary.
func AnnotateError(errorAdviceMap ErrorAnnotationMappings, errs ...error) []error {
	switch {

	// Process errors as long as the collection is not empty or not composed
	// entirely of nil values.
	case !isNilErrCollection(errs):

		// Create copy of provided errors collection so that we can safely
		// modify the copy without risk of touching the original.
		annotatedErrors := make([]error, 0, len(errs))
		// annotatedErrors := make([]error, len(errs))
		// copy(annotatedErrors, errs)

		// If the caller did not provide a custom error advice map (or
		// provided an empty map) then fallback to the default collection.
		if len(errorAdviceMap) == 0 {
			errorAdviceMap = DefaultErrorAnnotationMappings()
		}

		for _, givenErr := range errs {
			if givenErr == nil {
				// It is possible that the caller provided an error collection
				// with some entries set to nil. Skip over each nil entry.
				continue
			}

			annotatedErrors = append(annotatedErrors, annotateError(givenErr, errorAdviceMap))
		}

		return annotatedErrors

	// No errors were provided for evaluation.
	default:
		return nil
	}

}

// AnnotateRecordedErrors adds additional human-readable explanation for
// errors encountered during plugin execution. This additional context is
// intended to help sysadmins remediate common issues detected by plugins.
//
// Each error already recorded in the collection is evaluated for a match in
// the provided error advice map. If the existing error collection is empty
// then no annotations are performed.
//
// If an error match is found *AND* the error is not already annotated with
// specified advice then the advice is appended to the error via error
// wrapping. This process is repeated for each applicable error in the
// given error chain.
//
// If no advice is found then the error is unmodified. The existing error
// collection is replaced with this (potentially) updated collection of error
// chains.
//
// NOTE: Deduplication of errors already recorded in the collection is *not*
// performed. The caller is responsible for ensuring that a given error is not
// already recorded in the collection.
func (p *Plugin) AnnotateRecordedErrors(errorAdviceMap ErrorAnnotationMappings) {
	p.Errors = AnnotateError(errorAdviceMap, p.Errors...)
}

// AddAnnotatedError adds additional human-readable explanation for errors
// encountered during plugin execution. This additional context is intended to
// help sysadmins remediate common issues detected by plugins.
//
// Each entry in the error chain for a given error is evaluated for a match in
// the provided error advice map. If an empty error collection or a collection
// of nil error values are provided for evaluation then nil is returned.
//
// If a match is found and the error is not already annotated with specified
// advice then the advice is appended to the error (via wrapping). If no
// advice is found then the original error is unmodified.
//
// Existing errors in the collection are unmodified. Given errors are appended
// (annotated or not) to the existing error collection.
//
// NOTE: Deduplication of errors is *not* performed. The caller is responsible
// for ensuring that a given error (annotated or not) is not already recorded
// in the collection.
//
// Another method is provided for callers which wish to skip insertion of an
// error if it is already present in the collection.
func (p *Plugin) AddAnnotatedError(errorAdviceMap ErrorAnnotationMappings, errs ...error) {
	annotatedErrors := AnnotateError(errorAdviceMap, errs...)
	p.AddError(annotatedErrors...)
}

// AddUniqueAnnotatedError adds additional human-readable explanation for
// errors encountered during plugin execution. This additional context is
// intended to help sysadmins remediate common issues detected by plugins.
//
// Each entry in the error chain for a given error is evaluated for a match in
// the provided error advice map. If an empty error collection or a collection
// of nil error values are provided for evaluation then nil is returned.
//
// If a match is found and the error is not already annotated with specified
// advice then the advice is appended to the error (via wrapping). If no
// advice is found then the original error is unmodified.
//
// Existing errors in the collection are unmodified.
//
// Annotated errors are appended to the collection *unless* they are
// determined to already be present. This evaluation is performed using
// case-insensitive string comparison.
func (p *Plugin) AddUniqueAnnotatedError(errorAdviceMap ErrorAnnotationMappings, errs ...error) {
	annotatedErrors := AnnotateError(errorAdviceMap, errs...)
	p.AddUniqueError(annotatedErrors...)
}

// annotateError is a helper function used to process a given error using a
// provided advice map. Annotations continue for each applicable error in the
// given error chain; a given error may be annotated multiple times, once for
// each match in the provided advice map.
func annotateError(err error, errorAdviceMap ErrorAnnotationMappings) error {
	annotatedErr := err

	for knownErr, advice := range errorAdviceMap {
		if errors.Is(annotatedErr, knownErr) {
			if annotationAlreadyPresent(err, advice) {
				continue
			}

			// Only apply annotation if not already present.
			annotatedErr = fmt.Errorf(
				"%w: %s",
				err,
				advice,
			)
		}
	}

	return annotatedErr
}

// annotationAlreadyPresent represents the logic used to determine whether a
// given error message string already contains specific advice. This is used
// to prevent appending the same advice multiple times.
func annotationAlreadyPresent(err error, advice string) bool {
	return strings.Contains(
		strings.ToLower(err.Error()),
		strings.ToLower(advice),
	)
}
