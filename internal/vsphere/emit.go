// Copyright 2026 Adam Chalkley
//
// https://github.com/atc0005/check-vmware
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package vsphere

import (
	"fmt"
	"io"

	"github.com/atc0005/go-nagios"
)

// emitSeparator emits a given separator string to the given Writer. If not
// provided, this function skips emitting any output to the Writer.
func emitSeparator(w io.Writer, separatorText string) {
	if separatorText != "" {
		_, _ = fmt.Fprintf(
			w,
			"%s%s%s%s",
			nagios.CheckOutputEOL,
			separatorText,
			nagios.CheckOutputEOL,
			nagios.CheckOutputEOL,
		)
	}
}
