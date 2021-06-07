// Copyright 2021 Adam Chalkley
//
// https://github.com/atc0005/check-vmware
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package main

import (
	"testing"
	"time"

	"github.com/atc0005/check-vmware/internal/config"
	"github.com/atc0005/check-vmware/internal/vsphere"
	"github.com/vmware/govmomi/vim25/types"
)

func getTestTriggeredAlarms() vsphere.TriggeredAlarms {

	return vsphere.TriggeredAlarms{

		// previously acknowledged (5 hours ago), triggered (24 hours ago)
		// yellow or WARNING datastore usage
		vsphere.TriggeredAlarm{
			Entity: vsphere.AlarmEntity{
				Name:          "RES-DC1-S6200-vol12",
				MOID:          types.ManagedObjectReference{Type: "Datastore", Value: "datastore-50120"},
				OverallStatus: types.ManagedEntityStatus("red"),
			},
			AcknowledgedTime:   time.Now().Add(-5 * time.Hour),
			Time:               time.Now().AddDate(0, 0, -1),
			Name:               "Datastore usage on disk",
			MOID:               types.ManagedObjectReference{Type: "Alarm", Value: "alarm-8"},
			Key:                "alarm-8.datastore-50120",
			Description:        "Default alarm to monitor datastore disk usage",
			Datacenter:         "Example",
			OverallStatus:      types.ManagedEntityStatus("yellow"),
			AcknowledgedByUser: "Ash",
			Acknowledged:       true,
		},

		// yellow or WARNING datastore usage
		vsphere.TriggeredAlarm{
			Entity: vsphere.AlarmEntity{
				Name:          "RES-DC1-S6200-vol11",
				MOID:          types.ManagedObjectReference{Type: "Datastore", Value: "datastore-50119"},
				OverallStatus: types.ManagedEntityStatus("yellow"),
			},
			AcknowledgedTime:   time.Time{},
			Time:               time.Now(),
			Name:               "Datastore usage on disk",
			MOID:               types.ManagedObjectReference{Type: "Alarm", Value: "alarm-8"},
			Key:                "alarm-8.datastore-50119",
			Description:        "Default alarm to monitor datastore disk usage",
			Datacenter:         "Example",
			OverallStatus:      types.ManagedEntityStatus("yellow"),
			AcknowledgedByUser: "",
			Acknowledged:       false,
		},

		// red or CRITICAL datastore usage
		vsphere.TriggeredAlarm{
			Entity: vsphere.AlarmEntity{
				Name:          "HUSVM-DC1-DigColl-vol8",
				MOID:          types.ManagedObjectReference{Type: "Datastore", Value: "datastore-141490"},
				OverallStatus: types.ManagedEntityStatus("red"),
			},
			AcknowledgedTime:   time.Time{},
			Time:               time.Now(),
			Name:               "Datastore usage on disk",
			MOID:               types.ManagedObjectReference{Type: "Alarm", Value: "alarm-8"},
			Key:                "alarm-8.datastore-141490",
			Description:        "Default alarm to monitor datastore disk usage",
			Datacenter:         "Example",
			OverallStatus:      types.ManagedEntityStatus("red"),
			AcknowledgedByUser: "",
			Acknowledged:       false,
		},

		// virtual machine CPU usage
		vsphere.TriggeredAlarm{
			Entity: vsphere.AlarmEntity{
				Name:          "node1.example.com",
				MOID:          types.ManagedObjectReference{Type: "VirtualMachine", Value: "vm-197"},
				OverallStatus: types.ManagedEntityStatus("red"),
			},
			AcknowledgedTime:   time.Time{},
			Time:               time.Now(),
			Name:               "Virtual machine CPU usage",
			MOID:               types.ManagedObjectReference{Type: "Alarm", Value: "alarm-6"},
			Key:                "alarm-6.vm-197",
			Description:        "Default alarm to monitor virtual machine CPU usage",
			Datacenter:         "Example",
			OverallStatus:      types.ManagedEntityStatus("red"),
			AcknowledgedByUser: "",
			Acknowledged:       false,
		},

		// virtual machine memory usage
		vsphere.TriggeredAlarm{
			Entity: vsphere.AlarmEntity{
				Name:          "node1.example.com",
				MOID:          types.ManagedObjectReference{Type: "VirtualMachine", Value: "vm-197"},
				OverallStatus: types.ManagedEntityStatus("red"),
			},
			AcknowledgedTime:   time.Time{},
			Time:               time.Now(),
			Name:               "Virtual machine memory usage",
			MOID:               types.ManagedObjectReference{Type: "Alarm", Value: "alarm-7"},
			Key:                "alarm-7.vm-197",
			Description:        "Default alarm to monitor virtual machine memory usage",
			Datacenter:         "Example",
			OverallStatus:      types.ManagedEntityStatus("red"),
			AcknowledgedByUser: "",
			Acknowledged:       false,
		},
	}
}

func TestFilters(t *testing.T) {

	if testing.Verbose() {
		t.Log("Enabling vsphere package logging output")
		vsphere.EnableLogging()
	}

	// setup table tests
	tests := []struct {
		testName                                  string
		cfg                                       config.Config
		pretendNoAlarms                           bool
		wantedNumTotalTriggeredAlarms             int
		wantedNumNonExcludedAlarmsBeforeFiltering int
		wantedNumNonExcludedAlarmsAfterFiltering  int
	}{
		{
			testName: "Include VirtualMachine, Exclude VM CPU usage",
			cfg: config.Config{
				Server:                     "vc1.example.com",
				Username:                   "vc1-read-only-service-account",
				Password:                   "placeholder",
				Domain:                     "example",
				LoggingLevel:               "info",
				DatacenterNames:            []string{"Example"},
				TrustCert:                  true,
				IncludedAlarmEntityTypes:   []string{"VirtualMachine"},
				ExcludedAlarmNames:         []string{"Virtual machine CPU usage"},
				EvaluateAcknowledgedAlarms: false,
			},
			wantedNumTotalTriggeredAlarms:             5,
			wantedNumNonExcludedAlarmsBeforeFiltering: 0,
			wantedNumNonExcludedAlarmsAfterFiltering:  1, // VirtuaslMachine memory usage
		},
		{
			testName: "Include VirtualMachine, Exclude VM CPU and memory usage",
			cfg: config.Config{
				Server:                     "vc1.example.com",
				Username:                   "vc1-read-only-service-account",
				Password:                   "placeholder",
				Domain:                     "example",
				LoggingLevel:               "info",
				DatacenterNames:            []string{"Example"},
				TrustCert:                  true,
				IncludedAlarmEntityTypes:   []string{"VirtualMachine"},
				ExcludedAlarmNames:         []string{"Virtual machine CPU usage", "memory usage"},
				EvaluateAcknowledgedAlarms: false,
			},
			wantedNumTotalTriggeredAlarms:             5,
			wantedNumNonExcludedAlarmsBeforeFiltering: 0,
			wantedNumNonExcludedAlarmsAfterFiltering:  0,
		},
		{
			testName: "Pretend no alarms",
			cfg: config.Config{
				Server:                     "vc1.example.com",
				Username:                   "vc1-read-only-service-account",
				Password:                   "placeholder",
				Domain:                     "example",
				LoggingLevel:               "info",
				DatacenterNames:            []string{"Example"},
				TrustCert:                  true,
				IncludedAlarmEntityTypes:   []string{"VirtualMachine"},
				ExcludedAlarmNames:         []string{"Virtual machine CPU usage"},
				EvaluateAcknowledgedAlarms: false,
			},
			pretendNoAlarms:                           true,
			wantedNumTotalTriggeredAlarms:             0,
			wantedNumNonExcludedAlarmsBeforeFiltering: 0,
			wantedNumNonExcludedAlarmsAfterFiltering:  0,
		},
		{
			testName: "Exclude datastore usage, tacos on sale",
			cfg: config.Config{
				Server:                     "vc1.example.com",
				Username:                   "vc1-read-only-service-account",
				Password:                   "placeholder",
				Domain:                     "example",
				LoggingLevel:               "info",
				DatacenterNames:            []string{"Example"},
				TrustCert:                  true,
				IncludedAlarmEntityTypes:   []string{},
				ExcludedAlarmNames:         []string{"datastore usage on disk", "tacos on sale"},
				EvaluateAcknowledgedAlarms: false,
			},
			wantedNumTotalTriggeredAlarms:             5,
			wantedNumNonExcludedAlarmsBeforeFiltering: 0,
			wantedNumNonExcludedAlarmsAfterFiltering:  2, // implicit VirtualMachine matches
		},
		{
			testName: "Include Tacos on sale",
			cfg: config.Config{
				Server:                     "vc1.example.com",
				Username:                   "vc1-read-only-service-account",
				Password:                   "placeholder",
				Domain:                     "example",
				LoggingLevel:               "info",
				DatacenterNames:            []string{"Example"},
				TrustCert:                  true,
				IncludedAlarmEntityTypes:   []string{},
				IncludedAlarmDescriptions:  []string{"tacos on sale"},
				ExcludedAlarmNames:         []string{},
				EvaluateAcknowledgedAlarms: false,
			},
			wantedNumTotalTriggeredAlarms:             5,
			wantedNumNonExcludedAlarmsBeforeFiltering: 0,
			wantedNumNonExcludedAlarmsAfterFiltering:  0,
		},
		{
			testName: "Include Tacos on sale, evaluate acknowledged alarms",
			cfg: config.Config{
				Server:                     "vc1.example.com",
				Username:                   "vc1-read-only-service-account",
				Password:                   "placeholder",
				Domain:                     "example",
				LoggingLevel:               "info",
				DatacenterNames:            []string{"Example"},
				TrustCert:                  true,
				IncludedAlarmEntityTypes:   []string{},
				IncludedAlarmDescriptions:  []string{"tacos on sale"},
				ExcludedAlarmNames:         []string{},
				EvaluateAcknowledgedAlarms: true,
			},
			wantedNumTotalTriggeredAlarms:             5,
			wantedNumNonExcludedAlarmsBeforeFiltering: 0,
			wantedNumNonExcludedAlarmsAfterFiltering:  0,
		},
		{
			testName: "Include datastore usage",
			cfg: config.Config{
				Server:                     "vc1.example.com",
				Username:                   "vc1-read-only-service-account",
				Password:                   "placeholder",
				Domain:                     "example",
				LoggingLevel:               "info",
				DatacenterNames:            []string{"Example"},
				TrustCert:                  true,
				IncludedAlarmEntityTypes:   []string{},
				IncludedAlarmNames:         []string{"datastore usage on disk"},
				ExcludedAlarmNames:         []string{},
				EvaluateAcknowledgedAlarms: false,
			},
			wantedNumTotalTriggeredAlarms:             5,
			wantedNumNonExcludedAlarmsBeforeFiltering: 0,
			wantedNumNonExcludedAlarmsAfterFiltering:  2,
		},
		{
			testName: "Include datastore usage, eval previously acknowledged",
			cfg: config.Config{
				Server:                     "vc1.example.com",
				Username:                   "vc1-read-only-service-account",
				Password:                   "placeholder",
				Domain:                     "example",
				LoggingLevel:               "info",
				DatacenterNames:            []string{"Example"},
				TrustCert:                  true,
				IncludedAlarmEntityTypes:   []string{},
				IncludedAlarmNames:         []string{"datastore usage on disk"},
				ExcludedAlarmNames:         []string{},
				EvaluateAcknowledgedAlarms: true,
			},
			wantedNumTotalTriggeredAlarms:             5,
			wantedNumNonExcludedAlarmsBeforeFiltering: 0,
			wantedNumNonExcludedAlarmsAfterFiltering:  3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {

			// initialize a fresh copy for every table test entry
			var triggeredAlarms vsphere.TriggeredAlarms
			if !tt.pretendNoAlarms {
				triggeredAlarms = getTestTriggeredAlarms()
			}

			// Pre-filtering checks
			numTriggeredAlarms := len(triggeredAlarms)
			switch {
			case numTriggeredAlarms != tt.wantedNumTotalTriggeredAlarms:
				t.Errorf(
					"want %d total triggered alarms before filtering; got %d",
					tt.wantedNumTotalTriggeredAlarms,
					numTriggeredAlarms,
				)
			default:
				t.Logf(
					"Got expected number (%d) of total triggered alarms before filtering",
					tt.wantedNumTotalTriggeredAlarms,
				)
			}

			numNonExcludedAlarmsBeforeFiltering := triggeredAlarms.NumExcluded()
			switch {
			case numNonExcludedAlarmsBeforeFiltering != tt.wantedNumNonExcludedAlarmsBeforeFiltering:
				t.Errorf(
					"want %d triggered alarms before filtering; got %d",
					tt.wantedNumNonExcludedAlarmsBeforeFiltering,
					numNonExcludedAlarmsBeforeFiltering,
				)
			default:
				t.Logf(
					"Got expected number (%d) of non-excluded alarms before filtering",
					tt.wantedNumNonExcludedAlarmsBeforeFiltering,
				)
			}

			switch {

			case len(triggeredAlarms) > 0:

				triggeredAlarmFilters := vsphere.TriggeredAlarmFilters{
					IncludedAlarmEntityTypes:   tt.cfg.IncludedAlarmEntityTypes,
					ExcludedAlarmEntityTypes:   tt.cfg.ExcludedAlarmEntityTypes,
					IncludedAlarmNames:         tt.cfg.IncludedAlarmNames,
					ExcludedAlarmNames:         tt.cfg.ExcludedAlarmNames,
					IncludedAlarmDescriptions:  tt.cfg.IncludedAlarmDescriptions,
					ExcludedAlarmDescriptions:  tt.cfg.ExcludedAlarmDescriptions,
					EvaluateAcknowledgedAlarms: tt.cfg.EvaluateAcknowledgedAlarms,
				}

				triggeredAlarms.Filter(triggeredAlarmFilters)

				numTriggeredAlarmsToIgnore := triggeredAlarms.NumExcluded()
				numTriggeredAlarmsToReport := len(triggeredAlarms) - numTriggeredAlarmsToIgnore

				t.Logf("%d Triggered Alarms to ignore", numTriggeredAlarmsToIgnore)
				t.Logf("%d Triggered Alarms to report", numTriggeredAlarmsToReport)

				// Post-filtering checks
				numNonExcludedAlarmsAfterFiltering := len(triggeredAlarms) - triggeredAlarms.NumExcluded()
				switch {
				case numNonExcludedAlarmsAfterFiltering != tt.wantedNumNonExcludedAlarmsAfterFiltering:
					t.Errorf(
						"want %d triggered alarms after filtering; got %d",
						tt.wantedNumNonExcludedAlarmsAfterFiltering,
						numNonExcludedAlarmsAfterFiltering,
					)
				default:
					t.Logf(
						"Got expected number (%d) of non-excluded alarms after filtering",
						tt.wantedNumNonExcludedAlarmsAfterFiltering,
					)
				}

				switch {
				case triggeredAlarms.HasCriticalState(false):
					t.Log("TriggeredAlarms have CRITICAL state")

				case triggeredAlarms.HasWarningState(false):
					t.Log("TriggeredAlarms have WARNING state")

				case triggeredAlarms.HasUnknownState(false):
					t.Log("TriggeredAlarms have UNKNOWN state")
				}

			default:

				t.Log("No non-excluded alarms detected")

			}
		})
	}

}
