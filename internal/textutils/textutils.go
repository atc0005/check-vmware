// Copyright 2021 Adam Chalkley
//
// https://github.com/atc0005/check-vmware
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package textutils

import (
	"fmt"
	"strconv"
	"strings"
)

// InList is a helper function to emulate Python's `if "x"
// in list:` functionality. The caller can optionally ignore case of compared
// items.
func InList(needle string, haystack []string, ignoreCase bool) bool {
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

// DedupeList returns a copy of a provided string slice with all duplicate
// entries removed.
// FIXME: Is there already a standard library version of this functionality?
func DedupeList(list []string) []string {

	// preallocate at least as much space as the original
	newList := make([]string, 0, len(list))
	uniqueItems := make(map[string]struct{})

	// build a map of unique list entries
	for _, item := range list {
		uniqueItems[item] = struct{}{}
	}

	// generate a new, deduped list from the map
	for key := range uniqueItems {
		newList = append(newList, key)
	}
	return newList
}

// IntSliceToStringSlice converts a slice of integers to a slice of string.
func IntSliceToStringSlice(ix []int) []string {
	sx := make([]string, len(ix))
	for i, v := range ix {
		sx[i] = strconv.Itoa(v)
	}
	return sx
}

// LowerCaseStringSlice is a helper function to convert all provided string
// slice elements to lowercase.
//
// FIXME: There is likely a better way to do this already; replace with that
// better way.
func LowerCaseStringSlice(xs []string) []string {
	lxs := make([]string, 0, len(xs))
	for idx := range xs {
		lxs = append(lxs, strings.ToLower(xs[idx]))
	}

	return lxs
}

// PrintHeader printers a section header to help separate otherwise
// potentially dense blocks of text.
func PrintHeader(headerText string) {
	headerBorderStr := strings.Repeat("=", len(headerText))
	fmt.Printf(
		"\n\n%s\n%s\n%s\n",
		headerBorderStr,
		headerText,
		headerBorderStr,
	)
}

// InsertDelimiter inserts a delimiter into the provided string every pos
// characters. If the length of the provided string is less than pos + 1
// characters the original string is returned unmodified as we are unable to
// insert delimiter between blocks of characters of specified (pos) size.
func InsertDelimiter(s string, delimiter string, pos int) string {

	if len(s) < pos+1 {
		return s
	}

	// convert string to rune slice in order to use unicode package functions
	// (which expect to work with runes).
	r := []rune(s)

	// to track position in string
	var ctr int

	var delimitedStr string
	for i, v := range r {
		c := string(v)
		ctr++

		// add delimiter when we have reached the specified position in the
		// string, provided that we've not reached the end of the string.
		if (ctr == pos) && (i+1 != len(r)) {
			delimitedStr += c + delimiter
			ctr = 0
			continue
		}
		delimitedStr += c
	}

	return delimitedStr
}

// BytesToDelimitedHexStr converts a byte slice to a delimited hex string.
func BytesToDelimitedHexStr(bx []byte, delimiter string) string {
	hexStr := make([]string, 0, len(bx))
	for _, v := range bx {
		hexStr = append(hexStr, fmt.Sprintf("%X", v))
	}
	return strings.Join(hexStr, delimiter)
}
