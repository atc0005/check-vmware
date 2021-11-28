// Copyright 2021 Adam Chalkley
//
// https://github.com/atc0005/check-vmware
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package vsphere

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// ErrRuntimeTimeoutReached indicates that plugin runtime exceeded specified
// timeout value.
var ErrRuntimeTimeoutReached = errors.New("plugin runtime exceeded specified timeout value")

// ErrVMwareAdminAssistanceNeeded indicates that a known/detected problem can
// only be resolved with the assistance of the administrators of the VMware
// environment(s) monitored by plugins in this project. While this team may be
// the same ones to receive the notifications from the monitoring system using
// this project's plugin, that may not always be the case.
var ErrVMwareAdminAssistanceNeeded = errors.New(
	"assistance needed from vmware administrators to resolve issue",
)

// AnnotateError is a helper function used to add additional human-readable
// explanation for errors commonly emitted by dependencies. This function
// receives an error, evaluates whether it contains specific errors in its
// chain and then (potentially) appends additional details for later use. This
// updated error chain is returned to the caller, preserving the original
// wrapped error. The original error is returned unmodified if no annotations
// were deemed necessary.
func AnnotateError(err error) error {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute AnnotateError func.\n",
			time.Since(funcTimeStart),
		)
	}()

	switch {

	case errors.Is(err, context.DeadlineExceeded):
		return fmt.Errorf("%w: %s", err, ErrRuntimeTimeoutReached)

	case errors.Is(err, ErrDatastoreIormConfigurationStatisticsCollectionDisabled):
		return fmt.Errorf("%w: %s", err, ErrVMwareAdminAssistanceNeeded)

	default:

		// Return error unmodified if additional decoration isn't defined for the
		// error type.
		return err

	}

}
