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

	// OverallStatus is the entity's top-level or overall status. vSphere
	// represents this status (aka, ManagedEntityStatus) as a color (gray,
	// green, red or yellow) with green indicating "OK" and red "CRITICAL".
	OverallStatus types.ManagedEntityStatus

	// MOID is the Managed Object Reference of the entity.
	MOID types.ManagedObjectReference

	// ResourcePools are the names of the Resource Pool that the
	// TriggeredAlarm entity is part of. This applies to VirtualMachine and
	// ResourcePool types. VirtualMachine types have one entry and
	// ResourcePool types have two (self & parent).
	ResourcePools []string
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

	// ExcludeReason gives a brief explanation of why a TriggeredAlarm is
	// excluded.
	ExcludeReason string

	// OverallStatus is the alarm's top-level or overall status of the alarm.
	// vSphere represents this status (aka, ManagedEntityStatus) as a color
	// (gray, green, red or yellow) with green indicating "OK" and red
	// "CRITICAL".
	OverallStatus types.ManagedEntityStatus

	// Acknowledged indicates whether the triggered alarm has been
	// acknowledged by an admin user.
	Acknowledged bool

	// Exclude indicates whether the TriggeredAlarm has been excluded from
	// final evaluation. During processing multiple filters are applied. We
	// track exclusion state through the filtering pipeline so that any
	// explicit inclusions chosen by the sysadmin will have the opportunity to
	// reset this state and have the TriggeredAlarm considered for evaluation.
	Exclude bool

	// ExplicitlyIncluded indicates whether the TriggeredAlarm has been marked
	// for explicit inclusion by a step in the filtering pipeline. A
	// TriggeredAlarm marked in this way is not "dropped" by later explicit
	// inclusion filtering steps in the pipeline.
	ExplicitlyIncluded bool

	// ExplicitlyExcluded indicates whether the TriggeredAlarm has been marked
	// for explicit exclusion by a step in the filtering pipeline.
	ExplicitlyExcluded bool
}

// TriggeredAlarms is a collection of alarms which have been triggered across
// one or more Datacenters.
type TriggeredAlarms []TriggeredAlarm

// TriggeredAlarmFilters is a collection of the options specified by the user
// for filtering detected TriggeredAlarms. This is most often used for
// providing summary information in logging or user-facing output.
type TriggeredAlarmFilters struct {
	IncludedAlarmEntityTypes         []string
	ExcludedAlarmEntityTypes         []string
	IncludedAlarmEntityNames         []string
	ExcludedAlarmEntityNames         []string
	IncludedAlarmEntityResourcePools []string
	ExcludedAlarmEntityResourcePools []string
	IncludedAlarmNames               []string
	ExcludedAlarmNames               []string
	IncludedAlarmDescriptions        []string
	ExcludedAlarmDescriptions        []string
	IncludedAlarmStatuses            []string
	ExcludedAlarmStatuses            []string
	EvaluateAcknowledgedAlarms       bool
}

// NumExcluded returns the number of TriggeredAlarms that have been implicitly
// or explicitly excluded.
func (tas TriggeredAlarms) NumExcluded() int {
	var num int
	for i := range tas {
		if tas[i].Excluded() {
			num++
		}
	}

	return num
}

// NumExcludedFinal returns the number of TriggeredAlarms that have been
// explicitly excluded from further evaluation.
func (tas TriggeredAlarms) NumExcludedFinal() int {
	var num int
	for i := range tas {
		if tas[i].ExcludedFinal() {
			num++
		}
	}

	return num
}

// FilterByKey returns the matching TriggeredAlarm for the provided unique
// identifier (key) for a TriggeredAlarm.
func (tas TriggeredAlarms) FilterByKey(key string) (TriggeredAlarm, error) {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute FilterByKey func.\n",
			time.Since(funcTimeStart),
		)
	}()

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
// with each TriggeredAlarm in the collection. If specified, keys are also
// returned for acknowledged triggered alarms. Keys are returned in ascending
// order.
func (tas TriggeredAlarms) Keys(evalAcknowledged bool, evalExcluded bool) []string {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute Keys func.\n",
			time.Since(funcTimeStart),
		)
	}()

	keys := make([]string, 0, len(tas))
	for i := range tas {
		switch {
		case tas[i].Acknowledged && !evalAcknowledged:
			continue
		case tas[i].Exclude && !evalExcluded:
			continue
		default:
			keys = append(keys, tas[i].Key)
		}
	}

	sort.Slice(keys, func(i, j int) bool {
		return strings.ToLower(keys[i]) < strings.ToLower(keys[j])
	})

	return keys

}

// KeysExcluded returns a list of TriggeredAlarm keys or unique identifiers
// associated with each TriggeredAlarm in the collection that has been
// excluded. Keys are returned in ascending order.
func (tas TriggeredAlarms) KeysExcluded() []string {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute KeysExcluded func.\n",
			time.Since(funcTimeStart),
		)
	}()

	keysExcl := make([]string, 0, len(tas))
	for i := range tas {
		if tas[i].Exclude {
			keysExcl = append(keysExcl, tas[i].Key)
		}
	}

	sort.Slice(keysExcl, func(i, j int) bool {
		return strings.ToLower(keysExcl[i]) < strings.ToLower(keysExcl[j])
	})

	return keysExcl

}

// Datacenters returns a list of Datacenter names associated with the
// collection of TriggeredAlarms.
func (tas TriggeredAlarms) Datacenters() []string {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute Datacenters func.\n",
			time.Since(funcTimeStart),
		)
	}()

	dcsIdx := make(map[string]struct{})
	dcs := make([]string, 0, len(tas))

	for i := range tas {
		dcsIdx[tas[i].Datacenter] = struct{}{}
	}

	for k := range dcsIdx {
		dcs = append(dcs, k)
	}

	sort.Slice(dcs, func(i, j int) bool {
		return strings.ToLower(dcs[i]) < strings.ToLower(dcs[j])
	})

	return dcs

}

// ResourcePools returns a list of ResourcePool names associated with the
// collection of TriggeredAlarms.
func (tas TriggeredAlarms) ResourcePools() []string {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute ResourcePools func.\n",
			time.Since(funcTimeStart),
		)
	}()

	rpsIdx := make(map[string]struct{})
	rps := make([]string, 0, len(tas))

	for i := range tas {
		for j := range tas[i].Entity.ResourcePools {
			rpsIdx[tas[i].Entity.ResourcePools[j]] = struct{}{}
		}

	}

	for k := range rpsIdx {
		rps = append(rps, k)
	}

	sort.Slice(rps, func(i, j int) bool {
		return strings.ToLower(rps[i]) < strings.ToLower(rps[j])
	})

	return rps

}

// HasCriticalState indicates whether the collection of TriggeredAlarms
// contains an alarm considered to be in a CRITICAL state. A boolean value is
// accepted which indicates whether TriggeredAlarm values marked for exclusion
// (during filtering) should also be considered. The caller is responsible for
// filtering the collection; processing of inclusion or exclusion lists should
// be performed prior to calling this method.
func (tas TriggeredAlarms) HasCriticalState(evalExcluded bool) bool {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute HasCriticalState func.\n",
			time.Since(funcTimeStart),
		)
	}()

	if len(tas) == 0 {
		return false
	}

	var hasCriticalState bool

	for i := range tas {
		if hasCriticalState {
			// NOTE: We are interested in whether *any* alarm is in CRITICAL
			// state. If we found a single match that is sufficient for our
			// purposes; we do not need to look any further.
			break
		}
		switch {
		case tas[i].Exclude && !evalExcluded:
			continue
		default:
			_, exitCode := EntityStatusToNagiosState(tas[i].OverallStatus)
			if exitCode == nagios.StateCRITICALExitCode {
				hasCriticalState = true
			}
		}
	}

	return hasCriticalState

}

// NumCriticalState indicates how many TriggeredAlarms in the collection are
// considered to be in a CRITICAL state. A boolean value is accepted which
// indicates whether all TriggeredAlarm values are evaluated or only those not
// marked for exclusion. The caller is responsible for filtering the
// collection; processing of inclusion or exclusion lists should be performed
// prior to calling this method.
func (tas TriggeredAlarms) NumCriticalState(evalExcluded bool) int {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute NumCriticalState func.\n",
			time.Since(funcTimeStart),
		)
	}()

	if len(tas) == 0 {
		return 0
	}

	var numCriticalState int

	for i := range tas {
		switch {
		case tas[i].Exclude && !evalExcluded:
			continue
		default:
			_, exitCode := EntityStatusToNagiosState(tas[i].OverallStatus)
			if exitCode == nagios.StateCRITICALExitCode {
				numCriticalState++
			}
		}
	}

	return numCriticalState

}

// HasWarningState indicates whether the collection of TriggeredAlarms
// contains an alarm considered to be in a WARNING state. A boolean value is
// accepted which indicates whether TriggeredAlarm values marked for exclusion
// (during filter) should also be considered. The caller is responsible for
// filtering the collection; processing of inclusion or exclusion lists should
// be performed prior to calling this method.
func (tas TriggeredAlarms) HasWarningState(evalExcluded bool) bool {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute HasWarningState func.\n",
			time.Since(funcTimeStart),
		)
	}()

	if len(tas) == 0 {
		return false
	}

	var hasWarningState bool

	for i := range tas {
		if hasWarningState {
			// NOTE: We are interested in whether *any* alarm is in WARNING
			// state. If we found a single match that is sufficient for our
			// purposes; we do not need to look any further.
			break
		}
		switch {
		case tas[i].Exclude && !evalExcluded:
			continue
		default:
			_, exitCode := EntityStatusToNagiosState(tas[i].OverallStatus)
			if exitCode == nagios.StateWARNINGExitCode {
				hasWarningState = true
			}
		}
	}

	return hasWarningState

}

// NumWarningState indicates how many TriggeredAlarms in the collection are
// considered to be in a WARNING state. A boolean value is accepted which
// indicates whether TriggeredAlarm values marked for exclusion (during
// filtering) should also be considered. The caller is responsible for
// filtering the collection; processing of inclusion or exclusion lists should
// be performed prior to calling this method.
func (tas TriggeredAlarms) NumWarningState(evalExcluded bool) int {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute NumWarningState func.\n",
			time.Since(funcTimeStart),
		)
	}()

	if len(tas) == 0 {
		return 0
	}

	var numWarningState int

	for i := range tas {
		switch {
		case tas[i].Exclude && !evalExcluded:
			continue
		default:
			_, exitCode := EntityStatusToNagiosState(tas[i].OverallStatus)
			if exitCode == nagios.StateWARNINGExitCode {
				numWarningState++
			}
		}

	}

	return numWarningState

}

// HasUnknownState indicates whether the collection of TriggeredAlarms
// contains an alarm considered to be in an UNKNOWN state. A boolean value is
// accepted which indicates whether TriggeredAlarm values marked for exclusion
// (during filtering) should also be considered. The caller is responsible for
// filtering the collection; processing of inclusion or exclusion lists should
// be performed prior to calling this method.
func (tas TriggeredAlarms) HasUnknownState(evalExcluded bool) bool {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute HasUnknownState func.\n",
			time.Since(funcTimeStart),
		)
	}()

	if len(tas) == 0 {
		return false
	}

	var hasUnknownState bool

	for i := range tas {
		if hasUnknownState {
			// NOTE: We are interested in whether *any* alarm is in UNKNOWN
			// state. If we found a single match that is sufficient for our
			// purposes; we do not need to look any further.
			break
		}
		switch {
		case tas[i].Exclude && !evalExcluded:
			continue
		default:
			_, exitCode := EntityStatusToNagiosState(tas[i].OverallStatus)
			if exitCode == nagios.StateUNKNOWNExitCode {
				hasUnknownState = true
			}
		}
	}

	return hasUnknownState

}

// NumUnknownState indicates how many TriggeredAlarms in the collection are
// considered to be in an UNKNOWN state. A boolean value is accepted which
// indicates whether TriggeredAlarm values marked for exclusion (during
// filtering) should also be considered. The caller is responsible for
// filtering the collection; processing of inclusion or exclusion lists should
// be performed prior to calling this method.
func (tas TriggeredAlarms) NumUnknownState(evalExcluded bool) int {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute NumUnknownState func.\n",
			time.Since(funcTimeStart),
		)
	}()

	if len(tas) == 0 {
		return 0
	}

	var numUnknownState int

	for i := range tas {
		switch {
		case tas[i].Exclude && !evalExcluded:
			continue
		default:
			_, exitCode := EntityStatusToNagiosState(tas[i].OverallStatus)
			if exitCode == nagios.StateUNKNOWNExitCode {
				numUnknownState++
			}
		}
	}

	return numUnknownState

}

// IsOKState indicates whether all alarms in the collection of TriggeredAlarms
// are considered to be in an OK state. A boolean value is accepted to control
// whether all TriggeredAlarm values are evaluated or only those not marked
// for exclusion. The caller is responsible for filtering the collection;
// processing of inclusion or exclusion lists should be performed prior to
// calling this method.
func (tas TriggeredAlarms) IsOKState(evalExcluded bool) bool {

	switch {
	case tas.HasCriticalState(evalExcluded):
		return false
	case tas.HasWarningState(evalExcluded):
		return false
	case tas.HasUnknownState(evalExcluded):
		return false
	default:
		return true
	}

}

// NumOKState indicates how many TriggeredAlarms in the collection are
// considered to be in an OK state. A boolean value is accepted which
// indicates whether TriggeredAlarm values marked for exclusion (during
// filtering) should also be considered. The caller is responsible for filtering
// the collection; processing of inclusion or exclusion lists should be
// performed prior to calling this method.
func (tas TriggeredAlarms) NumOKState(evalExcluded bool) int {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute NumOKState func.\n",
			time.Since(funcTimeStart),
		)
	}()

	var numOKState int

	for i := range tas {
		switch {
		case tas[i].Exclude && !evalExcluded:
			continue
		default:
			_, exitCode := EntityStatusToNagiosState(tas[i].OverallStatus)
			if exitCode == nagios.StateOKExitCode {
				numOKState++
			}
		}

	}

	return numOKState

}

// Excluded indicates whether a TriggeredAlarm has been excluded implicitly
// (for now) or explicitly (permanently) from further evaluation.
func (ta TriggeredAlarm) Excluded() bool {

	if ta.ExplicitlyExcluded || ta.Exclude {
		return true
	}

	return false

}

// ExcludedFinal indicates whether a TriggeredAlarm has been permanently
// excluded from further evaluation.
func (ta TriggeredAlarm) ExcludedFinal() bool {

	return ta.ExplicitlyExcluded

}

// logExcluded is a helper method for logging when a TriggeredAlarm has been
// marked for exclusion, mostly for debugging purposes.
func (ta TriggeredAlarm) logExcluded(explicit bool) {
	logTriggeredAlarmMarked(ta, false, explicit)
}

// logIncluded is a helper method for logging when a TriggeredAlarm has been
// marked for inclusion, mostly for debugging purposes.
func (ta TriggeredAlarm) logIncluded(explicit bool) {
	logTriggeredAlarmMarked(ta, true, explicit)
}

// logTriggeredAlarmMarked is a helper function for logging when a
// TriggeredAlarm has been marked for inclusion or exclusion, mostly for
// debugging purposes.
func logTriggeredAlarmMarked(triggeredAlarm TriggeredAlarm, keep bool, explicit bool) {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute logTriggeredAlarmMarked func.\n",
			time.Since(funcTimeStart),
		)
	}()

	markType := "implicitly"
	if explicit {
		markType = "explicitly"
	}

	// create comma-separated list of resource pools for entity if provided,
	// otherwise produce a NOOP
	var rpsList string
	if len(triggeredAlarm.Entity.ResourcePools) > 0 {
		rpsList = fmt.Sprintf(
			" from pools [%q]",
			strings.Join(triggeredAlarm.Entity.ResourcePools, ", "),
		)
	}

	switch {
	case keep:
		logger.Printf(
			"Alarm (%s) for entity name %q of type %q%s with alarm name %q %s marked for inclusion",
			triggeredAlarm.OverallStatus,
			triggeredAlarm.Entity.Name,
			triggeredAlarm.Entity.MOID.Type,
			rpsList,
			triggeredAlarm.Name,
			markType,
		)

	default:
		logger.Printf(
			"Alarm (%s) for entity name %q of type %q%s with alarm name %q %s marked for exclusion",
			triggeredAlarm.OverallStatus,
			triggeredAlarm.Entity.Name,
			triggeredAlarm.Entity.MOID.Type,
			rpsList,
			triggeredAlarm.Name,
			markType,
		)

	}

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

	if len(datacenters) == 0 {
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

			resourcePoolNames := make([]string, 0, 2)
			switch {
			case entity.Self.Type == MgObjRefTypeResourcePool ||
				entity.Self.Type == MgObjRefTypeVirtualApp ||
				entity.Self.Type == MgObjRefTypeVirtualMachine:

				rps, err := getResourcePools(ctx, c, entity.Self, propsSubset)
				if err != nil {
					return nil, err
				}
				for _, rp := range rps {
					resourcePoolNames = append(resourcePoolNames, rp.Name)
				}

			default:
				// As far as I know, no other types can be "part" of a
				// Resource Pool. Set to empty slice to help prevent nil
				// pointer dereferencing attempts later. This matches the
				// pattern used for other fields in our custom TriggeredAlarm
				// type.
				resourcePoolNames = []string{}
			}

			triggeredAlarm := TriggeredAlarm{
				Entity: AlarmEntity{
					Name:          entity.Name,
					MOID:          entity.Self,
					ResourcePools: resourcePoolNames,
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

	// Sorting is only needed at initialization as we retain all
	// TriggeredAlarm values during later filtering/processing pipeline
	// phases.
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

// getSubstringFilterKeywords is a helper function that returns a map of all
// valid keywords used by the TriggeredAlarms.filterByString method.
// func getSubstringFilterKeywords() map[string]struct{} {
// 	return map[string]struct{}{
// 		AlarmDescription: struct{}{},
// 		AlarmName:        struct{}{},
// 		EntityName:       struct{}{},
// 	}
// }

// Filter explicitly includes or excludes TriggeredAlarms based on specified
// filter settings.
func (tas *TriggeredAlarms) Filter(filters TriggeredAlarmFilters) {

	logger.Println("Filtering triggered alarms by acknowledged state")
	tas.FilterByAcknowledgedState(filters.EvaluateAcknowledgedAlarms)

	logger.Println("Filtering triggered alarms by entity type")
	tas.filterByEntityType(filters.IncludedAlarmEntityTypes, filters.ExcludedAlarmEntityTypes)

	logger.Println("Filtering triggered alarms by name")
	tas.filterBySubstring(alarmName, filters.IncludedAlarmNames, filters.ExcludedAlarmNames)

	logger.Println("Filtering triggered alarms by description")
	tas.filterBySubstring(alarmDescription, filters.IncludedAlarmDescriptions, filters.ExcludedAlarmDescriptions)

	logger.Println("Filtering triggered alarms by status")
	tas.filterByStatus(filters.IncludedAlarmStatuses, filters.ExcludedAlarmStatuses)

	logger.Println("Filtering triggered alarms by entity name")
	tas.filterBySubstring(entityName, filters.IncludedAlarmEntityNames, filters.ExcludedAlarmEntityNames)

	logger.Println("Filtering triggered alarms by entity resource pool")
	tas.filterByEntityResourcePool(filters.IncludedAlarmEntityResourcePools, filters.ExcludedAlarmEntityResourcePools)

}

// FilterByIncludedEntityType accepts a slice of entity type keywords to use
// in comparison against the entity type associated with a TriggeredAlarm. For
// any matches, the TriggeredAlarm is marked as explicitly included. This will
// prevent later filtering from implicitly excluding the TriggeredAlarm, but
// will not stop explicit exclusions from "dropping" the TriggeredAlarm from
// further evaluation in the filtering pipeline.
func (tas *TriggeredAlarms) FilterByIncludedEntityType(includeTypes []string) {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute FilterByIncludedEntityType func for %d types\n",
			time.Since(funcTimeStart),
			len(includeTypes),
		)
	}()

	tas.filterByEntityType(includeTypes, []string{})

}

// FilterByExcludedEntityType accepts a slice of entity type keywords to use
// in comparison against the entity type associated with a TriggeredAlarm. For
// any matches, the TriggeredAlarm is marked as explicitly excluded. This will
// result in "dropping" the TriggeredAlarm from further evaluation in the
// filtering pipeline.
func (tas *TriggeredAlarms) FilterByExcludedEntityType(excludeTypes []string) {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute FilterByExcludedEntityType func for %d types\n",
			time.Since(funcTimeStart),
			len(excludeTypes),
		)
	}()

	tas.filterByEntityType([]string{}, excludeTypes)

}

// filterByEntityType uses slices of entity type values to explicitly mark
// TriggeredAlarm values for inclusion or exclusion in the final evaluation.
// Flag evaluation logic prevents sysadmins from providing both an inclusion
// and exclusion list.
func (tas *TriggeredAlarms) filterByEntityType(include []string, exclude []string) {

	funcTimeStart := time.Now()

	// Collect number of non-excluded TriggeredAlarms at the start of this
	// filtering process. We'll collect this number again after filtering has
	// been applied in order to show the results of this filter.
	nonExcludedStart := len(*tas) - tas.NumExcluded()

	defer func(start *int) {
		logger.Printf(
			"It took %v to execute filterByEntityType func (for %d non-excluded TriggeredAlarms, yielding %d non-excluded TriggeredAlarms)\n",
			time.Since(funcTimeStart),
			*start,
			len(*tas)-tas.NumExcluded(),
		)
	}(&nonExcludedStart)

	switch {
	// if the collection of TriggeredAlarms is empty, skip filtering attempts.
	case len(*tas) == 0:
		logger.Println("Triggered Alarms list is empty, aborting")
		return

	// if we're not limiting TriggeredAlarms by entity type, skip filtering
	// attempts.
	case len(include) == 0 && len(exclude) == 0:
		logger.Println("Triggered Alarms entity type inclusion and exclusion lists are empty, aborting")
		return
	}

	switch {
	case len(include) > 0:
		logger.Printf(
			"Include list provided; explicitly marking TriggeredAlarms for inclusion for %d specified types",
			len(include),
		)

	case len(exclude) > 0:
		logger.Printf(
			"Exclude list provided; explicitly marking TriggeredAlarms for exclusion for %d specified types",
			len(exclude),
		)
	}

	for i := range *tas {

		switch {

		case len(include) > 0:

			switch {

			// If the Entity Type of the TriggeredAlarm matches one of the
			// provided type keywords mark TriggeredAlarm as explicitly
			// included.
			case textutils.InList((*tas)[i].Entity.MOID.Type, include, true):

				// Don't explicitly *include* the TriggeredAlarm if the
				// TriggeredAlarm has already been explicitly *excluded*.
				if !(*tas)[i].ExplicitlyExcluded {
					(*tas)[i].Exclude = false
					(*tas)[i].ExplicitlyIncluded = true
					(*tas)[i].logIncluded(true)
				}

			// if not explicitly included by another filter in the pipeline,
			// implicitly mark as excluded
			default:
				if !(*tas)[i].ExplicitlyIncluded {
					(*tas)[i].Exclude = true
					(*tas)[i].ExcludeReason = alarmExcludeReasonEntityType
					(*tas)[i].logExcluded(false)
				}

			}

		case len(exclude) > 0:

			// explicitly excluded
			//
			// no implicit inclusions are applied for non-matching alarm types
			// as that could unintentionally flip the results from earlier
			// filtering stages.
			if textutils.InList((*tas)[i].Entity.MOID.Type, exclude, true) {
				(*tas)[i].Exclude = true
				(*tas)[i].ExcludeReason = alarmExcludeReasonEntityType
				(*tas)[i].ExplicitlyExcluded = true
				// (*tas)[i].ExplicitlyIncluded = false
				(*tas)[i].logExcluded(true)
			}

		}
	}

}

// FilterByIncludedEntityResourcePool accepts a slice of Resource Pool names
// to use in comparison against the Resource Pool for an entity associated
// with a TriggeredAlarm. For any matches, the TriggeredAlarm is marked as
// explicitly included. This will prevent later filtering from implicitly
// excluding the TriggeredAlarm, but will not stop explicit exclusions from
// "dropping" the TriggeredAlarm from further evaluation in the filtering
// pipeline.
func (tas *TriggeredAlarms) FilterByIncludedEntityResourcePool(include []string) {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute FilterByIncludedEntityResourcePool func for %d Resource Pools\n",
			time.Since(funcTimeStart),
			len(include),
		)
	}()

	tas.filterByEntityResourcePool(include, []string{})

}

// FilterByExcludedEntityResourcePool accepts a slice of Resource Pool names
// to use in comparison against the Resource Pool for an entity associated
// with a TriggeredAlarm. For any matches, the TriggeredAlarm is marked as
// explicitly excluded. This will result in "dropping" the TriggeredAlarm from
// further evaluation in the filtering pipeline.
func (tas *TriggeredAlarms) FilterByExcludedEntityResourcePool(exclude []string) {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute FilterByExcludedEntityResourcePool func for %d Resource Pools\n",
			time.Since(funcTimeStart),
			len(exclude),
		)
	}()

	tas.filterByEntityResourcePool([]string{}, exclude)

}

// filterByEntityResourcePool accepts slices of Resource Pool names to use in
// comparison against the Resource Pool for an entity associated with a
// TriggeredAlarm. This is done to explicitly mark TriggeredAlarm values for
// inclusion or exclusion in the final evaluation. Flag evaluation logic
// prevents sysadmins from providing both an inclusion and exclusion list.
func (tas *TriggeredAlarms) filterByEntityResourcePool(include []string, exclude []string) {

	funcTimeStart := time.Now()

	// Collect number of non-excluded TriggeredAlarms at the start of this
	// filtering process. We'll collect this number again after filtering has
	// been applied in order to show the results of this filter.
	nonExcludedStart := len(*tas) - tas.NumExcluded()

	defer func(start *int) {
		logger.Printf(
			"It took %v to execute filterByEntityResourcePool func (for %d non-excluded TriggeredAlarms, yielding %d non-excluded TriggeredAlarms)\n",
			time.Since(funcTimeStart),
			*start,
			len(*tas)-tas.NumExcluded(),
		)
	}(&nonExcludedStart)

	switch {
	// if the collection of TriggeredAlarms is empty, skip filtering attempts.
	case len(*tas) == 0:
		logger.Println("Triggered Alarms list is empty, aborting")
		return

	// if we're not limiting TriggeredAlarms by entity type, skip filtering
	// attempts.
	case len(include) == 0 && len(exclude) == 0:
		logger.Println("Triggered Alarms entity Resource Pool inclusion and exclusion lists are empty, aborting")
		return
	}

	switch {
	case len(include) > 0:
		logger.Printf(
			"Include list provided; explicitly marking TriggeredAlarms for inclusion for %d specified Resource Pools",
			len(include),
		)

	case len(exclude) > 0:
		logger.Printf(
			"Exclude list provided; explicitly marking TriggeredAlarms for exclusion for %d specified Resource Pools",
			len(exclude),
		)
	}

	for i := range *tas {

		switch {

		case len(include) > 0:

			// at this point we have a list of Resource Pool names (via
			// include list) to compare against a list of actual Resource Pool
			// names associated with the Triggered Alarms.

			switch {

			// if we have resource pools for the entity, examine to see if any
			// are in our explicit inclusion list.
			case len((*tas)[i].Entity.ResourcePools) > 0:
				for j := range (*tas)[i].Entity.ResourcePools {

					switch {

					// If the Resource Pool names associated with the
					// TriggeredAlarm matches one of the provided Resource
					// Pool names to compare against mark the TriggeredAlarm
					// as explicitly included.
					case textutils.InList((*tas)[i].Entity.ResourcePools[j], include, true):

						// Don't explicitly *include* the TriggeredAlarm if
						// the TriggeredAlarm has already been explicitly
						// *excluded*.
						if !(*tas)[i].ExplicitlyExcluded {
							(*tas)[i].Exclude = false
							(*tas)[i].ExplicitlyIncluded = true
							(*tas)[i].logIncluded(true)
						}

					// if not explicitly included by another filter in the
					// pipeline, implicitly mark as excluded
					default:
						if !(*tas)[i].ExplicitlyIncluded {
							(*tas)[i].Exclude = true
							(*tas)[i].ExcludeReason = alarmExcludeReasonEntityResourcePool
							(*tas)[i].logExcluded(false)
						}

					}
				}

			// there are no resource pools for the entity, so therefore no
			// match is possible
			default:
				// if not explicitly included by another filter in the
				// pipeline, implicitly mark as excluded
				if !(*tas)[i].ExplicitlyIncluded {
					(*tas)[i].Exclude = true
					(*tas)[i].ExcludeReason = alarmExcludeReasonEntityResourcePool
					(*tas)[i].logExcluded(false)
				}
			}

		case len(exclude) > 0:

			// If we have resource pools for the entity, examine to see if any
			// are in our explicit exclusion list. If there are no resource
			// pools for the entity, then we won't exclude anything based on a
			// resource pool match (none to match against).
			for j := range (*tas)[i].Entity.ResourcePools {

				// explicitly excluded
				//
				// no implicit inclusions are applied for non-matching alarm
				// types as that could unintentionally flip the results from
				// earlier filtering stages.
				if textutils.InList((*tas)[i].Entity.ResourcePools[j], exclude, true) {
					(*tas)[i].Exclude = true
					(*tas)[i].ExcludeReason = alarmExcludeReasonEntityResourcePool
					(*tas)[i].ExplicitlyExcluded = true
					// (*tas)[i].ExplicitlyIncluded = false
					(*tas)[i].logExcluded(true)
				}

			}

		}
	}

}

// FilterByAcknowledgedState accepts a boolean value to indicate whether
// previously acknowledged alarms should be included in the final evaluation.
//
// If false, previously acknowledged TriggeredAlarms are marked as explicitly
// excluded and will be "dropped" from further evaluation in the filtering
// pipeline. If true, no changes are made. Further evaluation in the filtering
// pipeline can still mark the TriggeredAlarm as excluded.
func (tas *TriggeredAlarms) FilterByAcknowledgedState(includeAcknowledged bool) {

	// Collect number of non-excluded TriggeredAlarms at the start of this
	// filtering process. We'll collect this number again after filtering has
	// been applied in order to show the results of this filter.
	nonExcludedStart := len(*tas) - tas.NumExcluded()

	funcTimeStart := time.Now()

	defer func(start *int) {
		logger.Printf(
			"It took %v to execute FilterByAcknowledgedState func (for %d non-excluded TriggeredAlarms, yielding %d non-excluded TriggeredAlarms)\n",
			time.Since(funcTimeStart),
			*start,
			len(*tas)-tas.NumExcluded(),
		)
	}(&nonExcludedStart)

	// if the collection of TriggeredAlarms is empty, skip filtering attempts.
	if len(*tas) == 0 {
		return
	}

	for i := range *tas {

		// Mark TriggeredAlarm as explicitly excluded if sysadmin did not opt
		// to evaluate previously acknowledged alarms
		if (*tas)[i].Acknowledged && !includeAcknowledged {
			(*tas)[i].Exclude = true
			(*tas)[i].ExcludeReason = alarmExcludeReasonAlarmAcknowledged
			(*tas)[i].ExplicitlyExcluded = true
			// (*tas)[i].ExplicitlyIncluded = false
			(*tas)[i].logExcluded(true)
		}

	}

}

// FilterByIncludedNameSubstring accepts a slice of substrings to use in
// comparisons against TriggeredAlarm names. For any matches, the
// TriggeredAlarm is marked as explicitly included. This will prevent later
// filtering from implicitly excluding the TriggeredAlarm, but will not stop
// explicit exclusions from "dropping" the TriggeredAlarm from further
// evaluation in the filtering pipeline.
func (tas *TriggeredAlarms) FilterByIncludedNameSubstring(include []string) {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute FilterByIncludedNameSubstring func for %d substrings\n",
			time.Since(funcTimeStart),
			len(include),
		)
	}()

	tas.filterBySubstring(alarmName, include, []string{})

}

// FilterByExcludedNameSubstring accepts a slice of substrings to use in
// comparisons against TriggeredAlarm names in order to explicitly mark
// TriggeredAlarms for exclusion in the final evaluation. Flag evaluation
// logic prevents sysadmins from providing both an inclusion and exclusion
// list.
func (tas *TriggeredAlarms) FilterByExcludedNameSubstring(exclude []string) {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute FilterByExcludedNameSubstring func for %d substrings\n",
			time.Since(funcTimeStart),
			len(exclude),
		)
	}()

	tas.filterBySubstring(alarmName, []string{}, exclude)

}

// FilterByIncludedDescriptionSubstring accepts a slice of substrings to use
// in comparisons against TriggeredAlarm descriptions.
//
// For any matches, the TriggeredAlarm is marked as explicitly included. This
// will prevent later filtering from implicitly excluding the TriggeredAlarm,
// but will not stop explicit exclusions from "dropping" the TriggeredAlarm
// from further evaluation in the filtering pipeline.
func (tas *TriggeredAlarms) FilterByIncludedDescriptionSubstring(include []string) {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute FilterByIncludedDescriptionSubstring func for %d substrings\n",
			time.Since(funcTimeStart),
			len(include),
		)
	}()

	tas.filterBySubstring(alarmDescription, include, []string{})

}

// FilterByExcludedDescriptionSubstring accepts a slice of substrings to use
// in comparisons against TriggeredAlarm descriptions in order to explicitly
// mark TriggeredAlarms for exclusion in the final evaluation. Flag evaluation
// logic prevents sysadmins from providing both an inclusion and exclusion
// list.
func (tas *TriggeredAlarms) FilterByExcludedDescriptionSubstring(exclude []string) {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute FilterByExcludedDescriptionSubstring func for %d substrings\n",
			time.Since(funcTimeStart),
			len(exclude),
		)
	}()

	tas.filterBySubstring(alarmDescription, []string{}, exclude)

}

// filterBySubstring accepts a field keyword and slices of substrings to use
// in comparisons against TriggeredAlarm fields in order to explicitly mark
// TriggeredAlarms for inclusion or exclusion in the final evaluation. The
// provided field keyword indicates which field the comparison should be
// against. If an invlaid field keyword is supplied the field comparison will
// default to using the alarm name.
//
// Flag evaluation logic prevents sysadmins from providing both an inclusion
// and exclusion list.
func (tas *TriggeredAlarms) filterBySubstring(fieldKeyword string, include []string, exclude []string) {

	funcTimeStart := time.Now()

	// Collect number of non-excluded TriggeredAlarms at the start of this
	// filtering process. We'll collect this number again after filtering has
	// been applied in order to show the results of this filter.
	nonExcludedStart := len(*tas) - tas.NumExcluded()

	defer func(start *int, keyword string) {
		logger.Printf(
			"It took %v to execute filterBySubstring func (for %d non-excluded TriggeredAlarms, using keyword %s, yielding %d non-excluded TriggeredAlarms)\n",
			time.Since(funcTimeStart),
			*start,
			keyword,
			len(*tas)-tas.NumExcluded(),
		)
	}(&nonExcludedStart, fieldKeyword)

	switch {
	// if the collection of TriggeredAlarms is empty, skip filtering attempts.
	case len(*tas) == 0:
		logger.Println("Triggered Alarms list is empty, aborting")
		return

	// if we're not limiting TriggeredAlarms by entity type, skip filtering
	// attempts.
	case len(include) == 0 && len(exclude) == 0:
		logger.Printf(
			"Triggered Alarms substring (%s) inclusion and exclusion lists are empty, aborting",
			fieldKeyword,
		)
		return
	}

	switch {
	case len(include) > 0:
		logger.Printf(
			"Include list provided; explicitly marking TriggeredAlarms for inclusion which match any of %d specified substrings",
			len(include),
		)

	case len(exclude) > 0:
		logger.Printf(
			"Exclude list provided; explicitly marking TriggeredAlarms for exclusion which match any of %d specified substrings",
			len(exclude),
		)
	}

	// validKeywords := getSubstringFilterKeywords()
	// if _, ok := validKeywords[fieldKeyword]; !ok {
	// 	logger.Printf("")
	// }

	logger.Printf("substring field keyword %q specified", fieldKeyword)
	for i := range *tas {

		var substrField string
		var excludeReason string
		switch fieldKeyword {
		case alarmDescription:
			substrField = (*tas)[i].Description
			excludeReason = alarmExcludeReasonAlarmDescription
		case alarmName:
			substrField = (*tas)[i].Name
			excludeReason = alarmExcludeReasonAlarmName
		case entityName:
			substrField = (*tas)[i].Entity.Name
			excludeReason = alarmExcludeReasonEntityName
		default:
			logger.Printf(
				"substring field %q not recognized, defaulting to alarm name",
				fieldKeyword,
			)
			substrField = (*tas)[i].Name
			excludeReason = alarmExcludeReasonAlarmName
		}

		switch {

		case len(include) > 0:

			for _, substr := range include {

				switch {

				// Attempt literal, case-insensitive match first then attempt
				// substring, case-insensitive match.
				case strings.EqualFold(substrField, substr) ||
					strings.Contains(substrField, substr):

					// Don't explicitly *include* the TriggeredAlarm if the
					// TriggeredAlarm has already been explicitly *excluded*.
					if !(*tas)[i].ExplicitlyExcluded {
						(*tas)[i].Exclude = false
						(*tas)[i].ExplicitlyIncluded = true
						(*tas)[i].logIncluded(true)
					}

				// If not explicitly included by another filter in the
				// pipeline, implicitly mark as excluded.
				default:
					if !(*tas)[i].ExplicitlyIncluded {
						(*tas)[i].Exclude = true
						(*tas)[i].ExcludeReason = excludeReason
						(*tas)[i].logExcluded(false)
					}
				}

			}

		case len(exclude) > 0:

			for _, substr := range exclude {

				// explicitly excluded
				//
				// no implicit inclusions are applied for non-matching alarm
				// types as that could unintentionally flip the results from
				// earlier filtering stages.
				if strings.EqualFold(substrField, substr) ||
					strings.Contains(substrField, substr) {
					(*tas)[i].Exclude = true
					(*tas)[i].ExcludeReason = excludeReason
					(*tas)[i].ExplicitlyExcluded = true
					// (*tas)[i].ExplicitlyIncluded = false
					(*tas)[i].logExcluded(true)
				}

			}

		}
	}
}

// FilterByIncludedStatus accepts a slice of ManagedEntityStatus keywords to
// use in comparisons against TriggeredAlarm statuses. For any matches, the
// TriggeredAlarm is marked as explicitly included. This will prevent later
// filtering from implicitly excluding the TriggeredAlarm, but will not stop
// explicit exclusions from "dropping" the TriggeredAlarm from further
// evaluation in the filtering pipeline.
func (tas *TriggeredAlarms) FilterByIncludedStatus(include []string) {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute FilterByIncludedStatus func for %d keywords\n",
			time.Since(funcTimeStart),
			len(include),
		)
	}()

	tas.filterByStatus(include, []string{})

}

// FilterByExcludedStatus accepts a slice of ManagedEntityStatus keywords to
// use in comparisons against TriggeredAlarm statuses in order to explicitly
// mark TriggeredAlarms for exclusion in the final evaluation. Flag evaluation
// logic prevents sysadmins from providing both an inclusion and exclusion
// list.
func (tas *TriggeredAlarms) FilterByExcludedStatus(exclude []string) {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute FilterByExcludedStatus func for %d keywords\n",
			time.Since(funcTimeStart),
			len(exclude),
		)
	}()

	tas.filterByStatus([]string{}, exclude)

}

// filterByStatus accepts slices of ManagedEntityStatus keywords to use in
// comparisons against TriggeredAlarm statuses n order to explicitly mark
// TriggeredAlarms for inclusion or exclusion in the final evaluation.
//
// Flag evaluation logic prevents sysadmins from providing both an inclusion
// and exclusion list.
func (tas *TriggeredAlarms) filterByStatus(include []string, exclude []string) {

	funcTimeStart := time.Now()

	// Collect number of non-excluded TriggeredAlarms at the start of this
	// filtering process. We'll collect this number again after filtering has
	// been applied in order to show the results of this filter.
	nonExcludedStart := len(*tas) - tas.NumExcluded()

	defer func(start *int) {
		logger.Printf(
			"It took %v to execute filterByStatus func (for %d non-excluded TriggeredAlarms, yielding %d non-excluded TriggeredAlarms)\n",
			time.Since(funcTimeStart),
			*start,
			len(*tas)-tas.NumExcluded(),
		)
	}(&nonExcludedStart)

	switch {
	// if the collection of TriggeredAlarms is empty, skip filtering attempts.
	case len(*tas) == 0:
		logger.Println("Triggered Alarms list is empty, aborting")
		return

	// if we're not limiting TriggeredAlarms by entity type, skip filtering
	// attempts.
	case len(include) == 0 && len(exclude) == 0:
		logger.Println("Triggered Alarms status inclusion and exclusion lists are empty, aborting")
		return
	}

	switch {
	case len(include) > 0:
		logger.Printf(
			"Include list provided; explicitly marking TriggeredAlarms for inclusion which match any of %d specified status keywords",
			len(include),
		)

	case len(exclude) > 0:
		logger.Printf(
			"Exclude list provided; explicitly marking TriggeredAlarms for exclusion which match any of %d specified status keywords",
			len(exclude),
		)
	}

	for i := range *tas {

		switch {

		case len(include) > 0:

			for _, keyword := range include {

				// logger.Printf(
				// 	"(incl) OverallStatus: %q, Keyword: %q",
				// 	string((*tas)[i].OverallStatus),
				// 	keyword,
				// )

				switch {

				case strings.EqualFold(string((*tas)[i].OverallStatus), keyword):

					// logger.Printf("SUCCESSFUL MATCH on keyword: %s\n", keyword)

					// Don't explicitly *include* the TriggeredAlarm if the
					// TriggeredAlarm has already been explicitly *excluded*.
					if !(*tas)[i].ExplicitlyExcluded {
						(*tas)[i].Exclude = false
						(*tas)[i].ExplicitlyIncluded = true
						(*tas)[i].logIncluded(true)
					}

				// If not explicitly included by another filter in the
				// pipeline, implicitly mark as excluded.
				default:

					// logger.Printf("FAILED MATCH on keyword: %s", keyword)

					if !(*tas)[i].ExplicitlyIncluded {
						(*tas)[i].Exclude = true
						(*tas)[i].ExcludeReason = alarmExcludeReasonAlarmStatus
						(*tas)[i].logExcluded(false)
					}
				}

			}

		case len(exclude) > 0:

			for _, keyword := range exclude {

				// logger.Printf(
				// 	"(excl) OverallStatus: %q, Keyword: %q",
				// 	string((*tas)[i].OverallStatus),
				// 	keyword,
				// )

				// explicitly excluded
				//
				// no implicit inclusions are applied for non-matching alarm
				// types as that could unintentionally flip the results from
				// earlier filtering stages.
				if strings.EqualFold(string((*tas)[i].OverallStatus), keyword) {
					(*tas)[i].Exclude = true
					(*tas)[i].ExcludeReason = alarmExcludeReasonAlarmStatus
					(*tas)[i].ExplicitlyExcluded = true
					// (*tas)[i].ExplicitlyIncluded = false
					(*tas)[i].logExcluded(true)
				}

			}

		}
	}
}

// AlarmsOneLineCheckSummary is used to generate a one-line Nagios service
// check results summary. This is the line most prominent in notifications.
func AlarmsOneLineCheckSummary(
	stateLabel string,
	triggeredAlarms TriggeredAlarms,
	datacentersEvaluated []string,
) string {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute AlarmsOneLineCheckSummary func.\n",
			time.Since(funcTimeStart),
		)
	}()

	switch {
	case !triggeredAlarms.IsOKState(false):
		return fmt.Sprintf(
			"%s: %d non-excluded Triggered Alarms detected (evaluated %d Datacenters, %d Triggered Alarms)",
			stateLabel,
			len(triggeredAlarms)-triggeredAlarms.NumExcluded(),
			len(datacentersEvaluated),
			len(triggeredAlarms),
		)

	default:
		return fmt.Sprintf(
			"%s: No non-excluded Triggered Alarms detected (evaluated %d Datacenters, %d Triggered Alarms)",
			stateLabel,
			len(datacentersEvaluated),
			len(triggeredAlarms),
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
	triggeredAlarms TriggeredAlarms,
	triggeredAlarmFilters TriggeredAlarmFilters,
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

	var report strings.Builder

	fmt.Fprintf(
		&report,
		"Non-excluded Triggered Alarms detected:%s%s",
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
	)

	numTriggeredAlarmsToReport := len(triggeredAlarms) - triggeredAlarms.NumExcluded()

	switch {
	case numTriggeredAlarmsToReport == 0:
		fmt.Fprintf(
			&report,
			"* None%s%s",
			nagios.CheckOutputEOL,
			nagios.CheckOutputEOL,
		)
	default:
		var alarmCtr int
		for i := range triggeredAlarms {
			// only look at non-excluded alarms
			if !triggeredAlarms[i].Exclude {
				alarmCtr++
				fmt.Fprintf(
					&report,
					"* (%.2d) %s (type %s): %s%s",
					alarmCtr,
					triggeredAlarms[i].Entity.Name,
					triggeredAlarms[i].Entity.MOID.Type,
					triggeredAlarms[i].Name,
					nagios.CheckOutputEOL,
				)
			}
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
	case triggeredAlarms.NumExcluded() == 0:
		fmt.Fprintf(
			&report,
			"* None%s%s",
			nagios.CheckOutputEOL,
			nagios.CheckOutputEOL,
		)
	default:
		var alarmCtr int
		for i := range triggeredAlarms {
			// only look at excluded alarms
			if triggeredAlarms[i].Exclude {
				alarmCtr++
				fmt.Fprintf(
					&report,
					"* (%.2d) %s (type: %q, alarm name: %q, exclude reason: %q)%s",
					alarmCtr,
					triggeredAlarms[i].Entity.Name,
					triggeredAlarms[i].Entity.MOID.Type,
					triggeredAlarms[i].Name,
					triggeredAlarms[i].ExcludeReason,
					nagios.CheckOutputEOL,
				)
			}
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
		"%s**NOTE: Explicit exclusions have precedence over inclusions**%s%s",
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
		"* Plugin User Agent: %s%s",
		c.Client.UserAgent,
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		&report,
		"* Triggered Alarms (evaluated: %d, ignored: %d, total: %d)%s",
		numTriggeredAlarmsToReport,
		triggeredAlarms.NumExcluded(),
		len(triggeredAlarms),
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		&report,
		"* Acknowledged Triggered Alarms evaluated: %t%s",
		triggeredAlarmFilters.EvaluateAcknowledgedAlarms,
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		&report,
		"* Triggered Alarms to explicitly include%s",
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		&report,
		"** entity types (%d): [%v]%s",
		len(triggeredAlarmFilters.IncludedAlarmEntityTypes),
		strings.Join(triggeredAlarmFilters.IncludedAlarmEntityTypes, ", "),
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		&report,
		"** entity names (%d): [%v]%s",
		len(triggeredAlarmFilters.IncludedAlarmEntityNames),
		strings.Join(triggeredAlarmFilters.IncludedAlarmEntityNames, ", "),
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		&report,
		"** entity resource pools (%d): [%v]%s",
		len(triggeredAlarmFilters.IncludedAlarmEntityResourcePools),
		strings.Join(triggeredAlarmFilters.IncludedAlarmEntityResourcePools, ", "),
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		&report,
		"** names (%d): [%v]%s",
		len(triggeredAlarmFilters.IncludedAlarmNames),
		strings.Join(triggeredAlarmFilters.IncludedAlarmNames, ", "),
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		&report,
		"** descriptions (%d): [%v]%s",
		len(triggeredAlarmFilters.IncludedAlarmDescriptions),
		strings.Join(triggeredAlarmFilters.IncludedAlarmDescriptions, ", "),
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		&report,
		"** statuses (%d): [%v]%s",
		len(triggeredAlarmFilters.IncludedAlarmStatuses),
		strings.Join(triggeredAlarmFilters.IncludedAlarmStatuses, ", "),
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		&report,
		"* Triggered Alarms to explicitly exclude%s",
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		&report,
		"** entity types (%d): [%v]%s",
		len(triggeredAlarmFilters.ExcludedAlarmEntityTypes),
		strings.Join(triggeredAlarmFilters.ExcludedAlarmEntityTypes, ", "),
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		&report,
		"** entity names (%d): [%v]%s",
		len(triggeredAlarmFilters.ExcludedAlarmEntityNames),
		strings.Join(triggeredAlarmFilters.ExcludedAlarmEntityNames, ", "),
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		&report,
		"** entity resource pools (%d): [%v]%s",
		len(triggeredAlarmFilters.ExcludedAlarmEntityResourcePools),
		strings.Join(triggeredAlarmFilters.ExcludedAlarmEntityResourcePools, ", "),
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		&report,
		"** names (%d): [%v]%s",
		len(triggeredAlarmFilters.ExcludedAlarmNames),
		strings.Join(triggeredAlarmFilters.ExcludedAlarmNames, ", "),
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		&report,
		"** descriptions (%d): [%v]%s",
		len(triggeredAlarmFilters.ExcludedAlarmDescriptions),
		strings.Join(triggeredAlarmFilters.ExcludedAlarmDescriptions, ", "),
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		&report,
		"** statuses (%d): [%v]%s",
		len(triggeredAlarmFilters.ExcludedAlarmStatuses),
		strings.Join(triggeredAlarmFilters.ExcludedAlarmStatuses, ", "),
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
		len(triggeredAlarms.Datacenters()),
		strings.Join(triggeredAlarms.Datacenters(), ", "),
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		&report,
		"* Resource Pools with Triggered Alarms (%d): [%v]%s",
		len(triggeredAlarms.ResourcePools()),
		strings.Join(triggeredAlarms.ResourcePools(), ", "),
		nagios.CheckOutputEOL,
	)

	return report.String()
}
