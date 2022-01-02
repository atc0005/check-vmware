// Copyright 2022 Adam Chalkley
//
// https://github.com/atc0005/check-vmware
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package vsphere

import (
	"fmt"
	"math"
	"strings"
	"time"
)

// ExceedsAge indicates whether a given event date is older than the specified
// number of days.
func ExceedsAge(event time.Time, days int) bool {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute ExceedsAge func.\n",
			time.Since(funcTimeStart),
		)
	}()

	// Flip specified number of days negative so that we can wind back that
	// many days from the current date.
	daysBack := -(days)
	ageThreshold := time.Now().AddDate(0, 0, daysBack)

	switch {
	case event.Before(ageThreshold):
		return true
	case event.Equal(ageThreshold):
		return false
	case event.After(ageThreshold):
		return false

	// TODO: Is there any other state than Before, Equal and After?
	// TODO: Perhaps remove 'After' and use this instead?
	default:
		return false
	}

}

// DaysAgo accepts an event date and returns the number of days ago that it
// occurred. Only full days are counted. If the event has not occurred yet, 0
// is returned.
func DaysAgo(event time.Time) int {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute DaysAgo func.\n",
			time.Since(funcTimeStart),
		)
	}()

	timeSince := time.Since(event).Hours()

	// Event has not occurred yet.
	if timeSince < 0 {
		timeSince = 0
	}

	// Toss remainder so that we only get the whole number of days
	daysSince := math.Trunc(timeSince / 24)

	return int(daysSince)

}

// FormattedTimeSinceEvent receives a Time value and converts it to a string
// representing the largest useful whole units of time in days and hours. For
// example, if an event occurred 1 year, 2 days and 3 hours ago, this function
// will return the string '367d 3h ago', but if the event was 3 hours ago then
// '3h ago' will be returned. If the event occurs in the future, a suffix of
// 'until' is used in the formatted string.
func FormattedTimeSinceEvent(event time.Time) string {

	timeSince := time.Since(event).Hours()

	logger.Printf("Given event time in string format: %s", event.String())
	logger.Printf("Hours since event: %v", timeSince)

	var futureEvent bool
	var formattedTime string
	var daysSinceStr string
	var hoursSinceStr string

	// Flip sign back to positive, note that the event occurs in the future
	// for later use.
	if timeSince < 0 {
		futureEvent = true
		timeSince *= -1
	}

	// Toss remainder so that we only get the whole number of days
	daysSince := math.Trunc(timeSince / 24)

	if daysSince > 0 {
		daysSinceStr = fmt.Sprintf("%dd", int64(daysSince))
	}

	// Multiply the whole number of days by 24 to get the hours value, then
	// subtract from the original number of hours until event to get the
	// number of hours leftover from the days calculation.
	hoursSince := math.Trunc(timeSince - (daysSince * 24))
	hoursSinceStr = fmt.Sprintf("%dh", int64(hoursSince))

	// Handle any leading space if the event has or will occur in less than a
	// day's time.
	formattedTime = strings.TrimSpace(daysSinceStr + " " + hoursSinceStr)

	switch {
	case futureEvent:
		formattedTime += " " + "until"
	case !futureEvent:
		formattedTime += " " + "ago"
	}

	return formattedTime

}
