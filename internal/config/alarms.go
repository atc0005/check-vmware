// Copyright 2021 Adam Chalkley
//
// https://github.com/atc0005/check-vmware
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package config

import (
	"fmt"
	"strings"
)

// getTriggeredAlarmStatuses is a helper function that returns a map of
// supported status keywords to vSphere ManagedEntityStatus values. This is
// used to provide keyword validation and support conversion of supported flag
// keywords to valid ManagedEntityStatus values.
func getTriggeredAlarmStatuses() map[string]string {
	return map[string]string{
		AlarmStatusRed:      AlarmStatusRed,
		AlarmStatusCritical: AlarmStatusRed,
		AlarmStatusYellow:   AlarmStatusYellow,
		AlarmStatusWarning:  AlarmStatusYellow,
		AlarmStatusGreen:    AlarmStatusGreen, // here for completeness; alarms are not triggered for this status
		AlarmStatusOk:       AlarmStatusGreen, // here for completeness; alarms are not triggered for this status
		AlarmStatusGray:     AlarmStatusGray,
		AlarmStatusUnknown:  AlarmStatusGray,
	}
}

// setAlarmStatuses evaluates user-provided triggered alarm status keywords
// and assigns a list of valid/equivalent (and de-duplicated)
// ManagedEntityStatus keywords to exported fields for later use. This method
// should be called *after* config validation has been performed.
func (c *Config) setAlarmStatuses() error {

	if len(c.includedAlarmStatuses) == 0 && len(c.excludedAlarmStatuses) == 0 {
		return nil
	}

	alarmStatuses := getTriggeredAlarmStatuses()

	includedStatusesIdx := make(map[string]struct{})
	for _, keyword := range c.includedAlarmStatuses {
		requestedkeyword := strings.ToLower(keyword)
		if _, ok := alarmStatuses[requestedkeyword]; !ok {
			return fmt.Errorf("invalid triggered alarm status for inclusion: %q", keyword)
		}
		includedStatusesIdx[requestedkeyword] = struct{}{}
	}
	for keyword := range includedStatusesIdx {
		c.IncludedAlarmStatuses = append(c.IncludedAlarmStatuses, keyword)
	}

	excludedStatusesIdx := make(map[string]struct{})
	for _, keyword := range c.excludedAlarmStatuses {
		requestedkeyword := strings.ToLower(keyword)
		if _, ok := alarmStatuses[requestedkeyword]; !ok {
			return fmt.Errorf("invalid triggered alarm status for exclusion: %q", keyword)
		}
		excludedStatusesIdx[requestedkeyword] = struct{}{}
	}
	for keyword := range excludedStatusesIdx {
		c.ExcludedAlarmStatuses = append(c.ExcludedAlarmStatuses, keyword)
	}

	return nil

}
