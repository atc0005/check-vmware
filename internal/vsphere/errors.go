// Copyright 2021 Adam Chalkley
//
// https://github.com/atc0005/check-vmware
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package vsphere

import (
	"time"

	"github.com/atc0005/go-nagios"
)

// VMwareAdminAssistanceNeeded indicates that a known/detected problem can
// only be resolved with the assistance of the administrators of the VMware
// environment(s) monitored by plugins in this project. While this team may be
// the same ones to receive the notifications from the monitoring system using
// this project's plugin, that may not always be the case.
const VMwareAdminAssistanceNeeded = "assistance needed from vmware administrators to resolve issue (see plugin doc for details)"

// AnnotateError is a helper function used to add additional human-readable
// explanation for errors encountered during plugin execution. We first apply
// common advice for more general errors then apply advice specific to errors
// routinely encountered by this specific project.
func AnnotateError(plugin *nagios.Plugin) {
	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute AnnotateError func(errors evaluated: %d)",
			time.Since(funcTimeStart),
			len(plugin.Errors),
		)
	}()

	// If nothing to process, skip setup/processing steps.
	if len(plugin.Errors) == 0 {
		return
	}

	// Start off with the default advice collection.
	errorAdviceMap := nagios.DefaultErrorAnnotationMappings()

	// Add project-specific error feedback.
	errorAdviceMap[ErrDatastoreIormConfigurationStatisticsCollectionDisabled] = VMwareAdminAssistanceNeeded

	// Apply error advice annotations.
	plugin.AnnotateRecordedErrors(errorAdviceMap)
}
