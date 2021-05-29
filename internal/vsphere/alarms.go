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
	"sort"
	"strings"
	"time"

	"github.com/atc0005/check-vmware/internal/textutils"
	"github.com/atc0005/go-nagios"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

// ErrAlarmNotExcludedFromEvaluation indicates that one or more alarms were
// detected and not excluded from evaluation.
var ErrAlarmNotExcludedFromEvaluation = errors.New("alarm detected and not excluded from evaluation")

// AlarmEntity is the affected resource associated with an alarm. For example,
// for a triggered "Datastore usage on disk" alarm, AlarmEntity represents the
// affected datastore.
type AlarmEntity struct {

	// Name is the name of the entity (e.g., HUSVM-DC1-vol6) associated with a
	// triggered alarm.
	Name string

	// MOID is the Managed Object Reference of the entity.
	MOID types.ManagedObjectReference

	// OverallStatus is the entity's top-level or overall status. vSphere
	// represents this status (aka, ManagedEntityStatus) as a color (gray,
	// green, red or yellow) with green indicating "OK" and red "CRITICAL".
	OverallStatus types.ManagedEntityStatus
}

// TriggeredAlarm represents the state of an alarm along with the affected
// resource associated with the alarm.
type TriggeredAlarm struct {

	// Entity is the affected resource associated with the triggered alarm.
	Entity AlarmEntity

	// AcknowledgedTime is the time when the triggered alarm was acknowledged
	// by an admin. Defaults to zero value if the triggered alarm has not been
	// acknowledged.
	AcknowledgedTime time.Time

	// Time is when the alarm was triggered.
	Time time.Time

	// Name is the name of the defined alarm.
	Name string

	// MOID is the Managed Object Reference to the defined alarm.
	MOID types.ManagedObjectReference

	// Key is the unique identifier for the triggered alarm or alarm "state"
	// (AlarmState), not the defined alarm (AlarmInfo).
	Key string

	// Description is the description of the defined alarm.
	Description string

	// Datacenter is the datacenter where the alarm was triggered.
	Datacenter string

	// AcknowledgedByUser is the user which acknowledged the alarm.
	AcknowledgedByUser string

	// OverallStatus is the alarm's top-level or overall status of the alarm.
	// vSphere represents this status (aka, ManagedEntityStatus) as a color
	// (gray, green, red or yellow) with green indicating "OK" and red
	// "CRITICAL".
	OverallStatus types.ManagedEntityStatus

	// Acknowledged indicates whether the triggered alarm has been
	// acknowledged by an admin user.
	Acknowledged bool
}

// TriggeredAlarms is a collection of alarms which have been triggered across
// one or more Datacenters.
type TriggeredAlarms []TriggeredAlarm

// IgnoredAlarms receives a collection of TriggeredAlarms that remained from
// earlier filtering and compares each entry against the current collection. A
// new collection is returned containing only the TriggeredAlarms not present
// in the current collection.
func (tas TriggeredAlarms) IgnoredAlarms(filteredAlarms TriggeredAlarms) TriggeredAlarms {

	// If the collections are of the same length, return an empty collection
	// to indicate that no TriggeredAlarms in the current collection were
	// ignored.
	if len(tas) == len(filteredAlarms) {
		return TriggeredAlarms{}
	}

	ignoredAlarms := make(TriggeredAlarms, 0, len(tas))

	filteredAlarmkeys := filteredAlarms.Keys()

	for i := range tas {
		if !textutils.InList(tas[i].Key, filteredAlarmkeys, false) {
			ignoredAlarms = append(ignoredAlarms, tas[i])
		}
	}

	sort.Slice(ignoredAlarms, func(i, j int) bool {
		return strings.ToLower(ignoredAlarms[i].Name) < strings.ToLower(ignoredAlarms[j].Name)
	})

	return ignoredAlarms

}

// FilterByKey returns the matching TriggeredAlarm for the provided unique
// identifier (key) for a TriggeredAlarm.
func (tas TriggeredAlarms) FilterByKey(key string) (TriggeredAlarm, error) {
	for i := range tas {
		if tas[i].Key == key {
			return tas[i], nil
		}
	}

	return TriggeredAlarm{}, fmt.Errorf(
		"provided key does not match TriggeredAlarm in this collection: %s",
		key,
	)
}

// CountPerDatacenter returns a map of Datacenter name to triggered alarms
// associated with the Datacenter name.
func (tas TriggeredAlarms) CountPerDatacenter() map[string]int {

	alarmCount := make(map[string]int)

	for i := range tas {
		alarmCount[tas[i].Datacenter]++
	}

	return alarmCount

}

// Keys returns a list of TriggeredAlarm keys or unique identifiers associated
// with each TriggeredAlarm in the collection.
func (tas TriggeredAlarms) Keys() []string {

	keys := make([]string, 0, len(tas))
	for i := range tas {
		keys = append(keys, tas[i].Key)
	}

	return keys

}

// Datacenters returns a list of Datacenter names associated with the
// collection of TriggeredAlarms.
func (tas TriggeredAlarms) Datacenters() []string {

	dcsIdx := make(map[string]struct{})
	dcs := make([]string, 0, len(dcsIdx))

	for i := range tas {
		dcsIdx[tas[i].Datacenter] = struct{}{}
	}

	for k := range dcsIdx {
		dcs = append(dcs, k)
	}

	return dcs

}

// HasCriticalState indicates whether the collection of TriggeredAlarms
// contains an alarm considered to be in a CRITICAL state. The caller is
// responsible for filtering the collection; processing of inclusion or
// exclusion lists should be performed prior to calling this method.
func (tas TriggeredAlarms) HasCriticalState() bool {

	var hasCriticalState bool

	for i := range tas {
		_, exitCode := EntityStatusToNagiosState(tas[i].OverallStatus)
		if exitCode == nagios.StateCRITICALExitCode {
			hasCriticalState = true
			break
		}
	}

	return hasCriticalState

}

// NumCriticalState indicates how many TriggeredAlarms in the collection are
// considered to be in a CRITICAL state. The caller is responsible for
// filtering the collection; processing of inclusion or exclusion lists should
// be performed prior to calling this method.
func (tas TriggeredAlarms) NumCriticalState() int {

	var numCriticalState int

	for i := range tas {
		_, exitCode := EntityStatusToNagiosState(tas[i].OverallStatus)
		if exitCode == nagios.StateCRITICALExitCode {
			numCriticalState++
		}
	}

	return numCriticalState

}

// HasWarningState indicates whether the collection of TriggeredAlarms
// contains an alarm considered to be in a WARNING state. The caller is
// responsible for filtering the collection; processing of inclusion or
// exclusion lists should be performed prior to calling this method.
func (tas TriggeredAlarms) HasWarningState() bool {

	var hasWarningState bool

	for i := range tas {
		_, exitCode := EntityStatusToNagiosState(tas[i].OverallStatus)
		if exitCode == nagios.StateWARNINGExitCode {
			hasWarningState = true
			break
		}
	}

	return hasWarningState

}

// NumWarningState indicates how many TriggeredAlarms in the collection are
// considered to be in a WARNING state. The caller is responsible for
// filtering the collection; processing of inclusion or exclusion lists should
// be performed prior to calling this method.
func (tas TriggeredAlarms) NumWarningState() int {

	var numWarningState int

	for i := range tas {
		_, exitCode := EntityStatusToNagiosState(tas[i].OverallStatus)
		if exitCode == nagios.StateWARNINGExitCode {
			numWarningState++
		}
	}

	return numWarningState

}

// HasUnknownState indicates whether the collection of TriggeredAlarms
// contains an alarm considered to be in an UNKNOWN state. The caller is
// responsible for filtering the collection; processing of inclusion or
// exclusion lists should be performed prior to calling this method.
func (tas TriggeredAlarms) HasUnknownState() bool {

	var hasUnknownState bool

	for i := range tas {
		_, exitCode := EntityStatusToNagiosState(tas[i].OverallStatus)
		if exitCode == nagios.StateUNKNOWNExitCode {
			hasUnknownState = true
			break
		}
	}

	return hasUnknownState

}

// NumUnknownState indicates how many TriggeredAlarms in the collection are
// considered to be in an UNKNOWN state. The caller is responsible for
// filtering the collection; processing of inclusion or exclusion lists should
// be performed prior to calling this method.
func (tas TriggeredAlarms) NumUnknownState() int {

	var numUnknownState int

	for i := range tas {
		_, exitCode := EntityStatusToNagiosState(tas[i].OverallStatus)
		if exitCode == nagios.StateUNKNOWNExitCode {
			numUnknownState++
		}
	}

	return numUnknownState

}

// IsOKState indicates whether all alarms in the collection of TriggeredAlarms
// are considered to be in an OK state. The caller is responsible for
// filtering the collection; processing of inclusion or exclusion lists should
// be performed prior to calling this method.
func (tas TriggeredAlarms) IsOKState() bool {

	switch {
	case tas.HasCriticalState():
		return false
	case tas.HasWarningState():
		return false
	case tas.HasUnknownState():
		return false
	default:
		return true
	}

}

// NumOKState indicates how many TriggeredAlarms in the collection are
// considered to be in an OK state. The caller is responsible for filtering
// the collection; processing of inclusion or exclusion lists should be
// performed prior to calling this method.
func (tas TriggeredAlarms) NumOKState() int {

	var numOKState int

	for i := range tas {
		_, exitCode := EntityStatusToNagiosState(tas[i].OverallStatus)
		if exitCode == nagios.StateOKExitCode {
			numOKState++
		}
	}

	return numOKState

}

// GetTriggeredAlarms accepts a list of Datacenters and a boolean value
// indicating whether only a subset of properties for datacenters and alarms
// should be returned. If requested, a subset of all available properties will
// be retrieved (faster) instead of recursively fetching all properties (about
// 2x as slow). Any TriggeredAlarms found are returned or an error if an empty
// list is provided or if there are issues retrieving properties for any
// TriggeredAlarms.
func GetTriggeredAlarms(ctx context.Context, c *govmomi.Client, datacenters []mo.Datacenter, propsSubset bool) (TriggeredAlarms, error) {
	//
	funcTimeStart := time.Now()

	// declare this early so that we can grab a pointer to it in order to
	// access the entries later
	var alarms TriggeredAlarms

	defer func(alarms *TriggeredAlarms, dcs []mo.Datacenter) {
		logger.Printf(
			"It took %v to execute GetTriggeredAlarms func (and retrieve %d Triggered Alarms from %d Datacenters).\n",
			time.Since(funcTimeStart),
			len(*alarms),
			len(dcs),
		)
	}(&alarms, datacenters)

	if datacenters == nil {
		return TriggeredAlarms{}, fmt.Errorf("empty datacenters list provided")
	}

	// Fetch all triggered AlarmState values for applicable datacenters.
	for _, dc := range datacenters {

		for _, alarmState := range dc.TriggeredAlarmState {

			var alarm mo.Alarm
			var alarmProps []string

			if propsSubset {
				alarmProps = getAlarmPropsSubset()
			}

			// Fetch Alarm definition associated with Triggered Alarm
			err := c.RetrieveOne(ctx, alarmState.Alarm, alarmProps, &alarm)
			if err != nil {
				return nil, err
			}

			// Fetch ManagedEntity associated with TriggeredAlarm
			var entity mo.ManagedEntity
			err = c.RetrieveOne(ctx, alarmState.Entity, nil, &entity)
			if err != nil {
				return nil, err
			}

			// Setup default time.Time value for alarm AcknowledgedTime if
			// the alarm hasn't yet been acknowledged.
			var acknowledgedTime time.Time
			if alarmState.AcknowledgedTime != nil {
				acknowledgedTime = *alarmState.AcknowledgedTime
			}

			var acknowledged bool
			if alarmState.Acknowledged != nil {
				acknowledged = *alarmState.Acknowledged
			}

			triggeredAlarm := TriggeredAlarm{
				Entity: AlarmEntity{
					Name:          entity.Name,
					MOID:          entity.Self,
					OverallStatus: entity.OverallStatus,
				},
				AcknowledgedTime:   acknowledgedTime,
				Time:               alarmState.Time,
				Name:               alarm.Info.Name,
				MOID:               alarm.Self,
				Key:                alarmState.Key,
				Description:        alarm.Info.Description,
				Datacenter:         dc.Name,
				OverallStatus:      alarmState.OverallStatus,
				AcknowledgedByUser: alarmState.AcknowledgedByUser,
				Acknowledged:       acknowledged,
			}

			alarms = append(alarms, triggeredAlarm)
		}
	}

	sort.Slice(alarms, func(i, j int) bool {
		return strings.ToLower(alarms[i].Entity.Name) < strings.ToLower(alarms[j].Entity.Name)
	})

	return alarms, nil
}

// EntityStatusToNagiosState converts a vSphere Managed Entity Status (e.g.,
// "red", "yellow") to a Nagios state label and exit code.
func EntityStatusToNagiosState(entityStatus types.ManagedEntityStatus) (string, int) {

	switch entityStatus {
	case types.ManagedEntityStatusGray:
		// Entity status is unknown, should be reviewed
		return nagios.StateUNKNOWNLabel, nagios.StateUNKNOWNExitCode

	case types.ManagedEntityStatusGreen:
		// Entity is OK
		return nagios.StateOKLabel, nagios.StateOKExitCode

	case types.ManagedEntityStatusYellow:
		// Entity monitoring thresholds have been crossed, should be reviewed
		return nagios.StateWARNINGLabel, nagios.StateWARNINGExitCode

	case types.ManagedEntityStatusRed:

		// Entity has a problem in need of remediation
		return nagios.StateCRITICALLabel, nagios.StateCRITICALExitCode

	default:
		// this shouldn't be reached, so assume the worst
		logger.Println("unknown entity status provided, assuming worst case")
		return nagios.StateCRITICALLabel, nagios.StateCRITICALExitCode
	}

}

// FilterTriggeredAlarmsByEntityType accepts a collection of TriggeredAlarms
// and slices of entity type values to include and exclude. These slices are
// used to determine what TriggeredAlarm values should be included in the
// returned collection. If the collection of provided TriggeredAlarms is
// empty, an empty collection is returned.
func FilterTriggeredAlarmsByEntityType(triggeredAlarms TriggeredAlarms, includeTypes []string, excludeTypes []string) TriggeredAlarms {

	// setup early so we can reference it from deferred stats output
	filteredTriggeredAlarms := make(TriggeredAlarms, 0, len(triggeredAlarms))

	funcTimeStart := time.Now()

	defer func(alarms TriggeredAlarms, filteredAlarms *TriggeredAlarms) {
		logger.Printf(
			"It took %v to execute FilterTriggeredAlarmsByEntityType func (for %d TriggeredAlarms, yielding %d TriggeredAlarms)\n",
			time.Since(funcTimeStart),
			len(alarms),
			len(*filteredAlarms),
		)
	}(triggeredAlarms, &filteredTriggeredAlarms)

	switch {
	// if the collection of TriggeredAlarm values is empty, return the empty
	// collection as-is.
	case len(triggeredAlarms) == 0:
		return triggeredAlarms

	// if we're not limiting the triggered alarm by entity type, return the
	// entire collection.
	case len(includeTypes) == 0 && len(excludeTypes) == 0:
		filteredTriggeredAlarms = triggeredAlarms
		return filteredTriggeredAlarms
	}

	switch {
	case len(includeTypes) > 0:
		logger.Println("Include list provided; keeping triggered alarms for any specified types, excluding others")

	case len(excludeTypes) > 0:
		logger.Println("Exclude list provided; ignoring triggered alarms for any specified types, keeping others")
	}

	for _, triggeredAlarm := range triggeredAlarms {
		switch {

		case len(includeTypes) > 0:
			if textutils.InList(triggeredAlarm.Entity.MOID.Type, includeTypes, true) {
				filteredTriggeredAlarms = append(filteredTriggeredAlarms, triggeredAlarm)
				continue
			}
			logger.Printf(
				"Alarm %s for %s of type %s filtered out",
				triggeredAlarm.Name,
				triggeredAlarm.Entity.Name,
				triggeredAlarm.Entity.MOID.Type,
			)

		case len(excludeTypes) > 0:
			if !textutils.InList(triggeredAlarm.Entity.MOID.Type, excludeTypes, true) {
				filteredTriggeredAlarms = append(filteredTriggeredAlarms, triggeredAlarm)
				continue
			}
			logger.Printf(
				"Alarm %s for %s of type %s filtered out",
				triggeredAlarm.Name,
				triggeredAlarm.Entity.Name,
				triggeredAlarm.Entity.MOID.Type,
			)
		}
	}

	sort.Slice(filteredTriggeredAlarms, func(i, j int) bool {
		return strings.ToLower(filteredTriggeredAlarms[i].Entity.Name) < strings.ToLower(filteredTriggeredAlarms[j].Entity.Name)
	})

	return filteredTriggeredAlarms

}

// FilterTriggeredAlarmsByAcknowledgedState accepts a collection of
// TriggeredAlarms and a boolean value to indicate whether previously
// acknowledged alarms should be included in the returned collection. If the
// collection of provided TriggeredAlarms is empty, an empty collection is
// returned.
func FilterTriggeredAlarmsByAcknowledgedState(triggeredAlarms TriggeredAlarms, includeAcknowledged bool) TriggeredAlarms {

	// setup early so we can reference it from deferred stats output
	filteredTriggeredAlarms := make(TriggeredAlarms, 0, len(triggeredAlarms))

	funcTimeStart := time.Now()

	defer func(vms TriggeredAlarms, filteredTriggeredAlarms *TriggeredAlarms) {
		logger.Printf(
			"It took %v to execute FilterTriggeredAlarmsByAcknowledgedState func (for %d TriggeredAlarms, yielding %d TriggeredAlarms)\n",
			time.Since(funcTimeStart),
			len(vms),
			len(*filteredTriggeredAlarms),
		)
	}(triggeredAlarms, &filteredTriggeredAlarms)

	if len(triggeredAlarms) == 0 {
		return triggeredAlarms
	}

	for _, alarm := range triggeredAlarms {
		switch {
		case !alarm.Acknowledged:
			filteredTriggeredAlarms = append(filteredTriggeredAlarms, alarm)

		case alarm.Acknowledged && includeAcknowledged:
			filteredTriggeredAlarms = append(filteredTriggeredAlarms, alarm)

		}
	}

	return filteredTriggeredAlarms

}

// AlarmsOneLineCheckSummary is used to generate a one-line Nagios service
// check results summary. This is the line most prominent in notifications.
func AlarmsOneLineCheckSummary(
	stateLabel string,
	allAlarms TriggeredAlarms,
	filteredAlarms TriggeredAlarms,
	includedAlarmEntityTypes []string,
	excludedAlarmEntityTypes []string,
) string {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute AlarmsOneLineCheckSummary func.\n",
			time.Since(funcTimeStart),
		)
	}()

	datacentersEvaluated := len(filteredAlarms.Datacenters())

	switch {
	case !filteredAlarms.IsOKState():
		return fmt.Sprintf(
			"%s: %d non-excluded Triggered Alarms detected (evaluated %d Datacenters, %d Triggered Alarms)",
			stateLabel,
			len(filteredAlarms),
			datacentersEvaluated,
			len(allAlarms),
		)

	default:
		return fmt.Sprintf(
			"%s: No non-excluded Triggered Alarms detected (evaluated %d Datacenters, %d Triggered Alarms)",
			stateLabel,
			datacentersEvaluated,
			len(allAlarms),
		)
	}
}

// AlarmsReport generates a summary of detected alarms along with various
// verbose details intended to aid in troubleshooting check results at a
// glance. This information is provided for use with the Long Service Output
// field commonly displayed on the detailed service check results display in
// the web UI or in the body of many notifications.
func AlarmsReport(
	c *vim25.Client,
	allAlarms TriggeredAlarms,
	filteredAlarms TriggeredAlarms,
	includedAlarmEntityTypes []string,
	excludedAlarmEntityTypes []string,
	evalAcknowledgedAlarms bool,
	specifiedDatacenters []string,
	datacentersEvaluated []string,
) string {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute AlarmsReport func.\n",
			time.Since(funcTimeStart),
		)
	}()

	// Build list of triggered alarms that have been filtered out
	ignoredAlarms := allAlarms.IgnoredAlarms(filteredAlarms)

	var report strings.Builder

	fmt.Fprintf(
		&report,
		"Non-excluded Triggered Alarms detected:%s%s",
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
	)

	switch {
	case len(filteredAlarms) == 0:
		fmt.Fprintf(&report, "* None%s", nagios.CheckOutputEOL)
	default:
		var alarmCtr int
		for i := range filteredAlarms {
			alarmCtr++
			fmt.Fprintf(
				&report,
				"* (%.2d) %s (type %s): %s%s",
				alarmCtr,
				filteredAlarms[i].Entity.Name,
				filteredAlarms[i].Entity.MOID.Type,
				filteredAlarms[i].Name,
				nagios.CheckOutputEOL,
			)
		}

		fmt.Fprintf(&report, "%s", nagios.CheckOutputEOL)

	}

	fmt.Fprintf(
		&report,
		"Excluded Triggered Alarms (as requested):%s%s",
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
	)

	switch {
	case len(ignoredAlarms) == 0:
		fmt.Fprintf(&report, "* None%s", nagios.CheckOutputEOL)
	default:
		var alarmCtr int
		for i := range ignoredAlarms {
			alarmCtr++
			fmt.Fprintf(
				&report,
				"* (%.2d) %s (type %s): %s%s",
				alarmCtr,
				ignoredAlarms[i].Entity.Name,
				ignoredAlarms[i].Entity.MOID.Type,
				ignoredAlarms[i].Name,
				nagios.CheckOutputEOL,
			)
		}
	}

	fmt.Fprintf(
		&report,
		"%s---%s%s",
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		&report,
		"* vSphere environment: %s%s",
		c.URL().String(),
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		&report,
		"* Triggered Alarms (evaluated: %d, ignored: %d, total: %d)%s",
		len(filteredAlarms),
		len(ignoredAlarms),
		len(allAlarms),
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		&report,
		"* Acknowledged Alarms evaluated: %t%s",
		evalAcknowledgedAlarms,
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		&report,
		"* Specified Triggered Alarm entity types to explicitly include (%d): [%v]%s",
		len(includedAlarmEntityTypes),
		strings.Join(includedAlarmEntityTypes, ", "),
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		&report,
		"* Specified Triggered Alarm entity types to explicitly exclude (%d): [%v]%s",
		len(excludedAlarmEntityTypes),
		strings.Join(excludedAlarmEntityTypes, ", "),
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		&report,
		"* Datacenters specified (%d): [%v]%s",
		len(specifiedDatacenters),
		strings.Join(specifiedDatacenters, ", "),
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		&report,
		"* Datacenters evaluated (%d): [%v]%s",
		len(datacentersEvaluated),
		strings.Join(datacentersEvaluated, ", "),
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		&report,
		"* Datacenters with Triggered Alarms (%d): [%v]%s",
		len(allAlarms.Datacenters()),
		strings.Join(allAlarms.Datacenters(), ", "),
		nagios.CheckOutputEOL,
	)

	return report.String()
}
