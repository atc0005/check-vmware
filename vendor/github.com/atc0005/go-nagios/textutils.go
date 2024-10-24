// Copyright 2023 Adam Chalkley
//
// https://github.com/atc0005/go-nagios
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package nagios

import "strings"

// inList is a helper function to emulate Python's `if "x" in list:`
// functionality. The caller can optionally ignore case of compared items.
func inList(needle string, haystack []string, ignoreCase bool) bool {
	for _, item := range haystack {

		if ignoreCase {
			if strings.EqualFold(item, needle) {
				return true
			}
		}

		if item == needle {
			return true
		}
	}

	return false
}

// removeEntry is a helper function to allow removing an entry or "line" from
// input which matches a given substring. The specified delimiter is used to
// perform the initial line splitting for entry removal and then to rejoin the
// elements into the original input string (minus the intended entry to
// remove).
func removeEntry(input string, substr string, delimiter string) string {
	if len(input) == 0 || len(substr) == 0 || len(delimiter) == 0 {
		return input
	}

	// https://stackoverflow.com/a/57213476/903870
	removeAtIndex := func(s []string, index int) []string {
		// ret := make([]string, 0)
		ret := make([]string, 0, len(s)-1)
		ret = append(ret, s[:index]...)
		return append(ret, s[index+1:]...)
	}

	lines := strings.Split(input, delimiter)
	var idxToRemove int
	for idx, line := range lines {
		if strings.Contains(line, substr) {
			idxToRemove = idx
		}
	}

	return strings.Join(removeAtIndex(lines, idxToRemove), delimiter)
}
