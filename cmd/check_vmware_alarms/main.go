// Copyright 2021 Adam Chalkley
//
// https://github.com/atc0005/check-vmware
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/atc0005/go-nagios"

	"github.com/atc0005/check-vmware/internal/config"
	"github.com/atc0005/check-vmware/internal/vsphere"

	zlog "github.com/rs/zerolog/log"
)

func main() {

	// Set initial "state" as valid, adjust as we go.
	var nagiosExitState = nagios.ExitState{
		LastError:      nil,
		ExitStatusCode: nagios.StateOKExitCode,
	}

	// defer this from the start so it is the last deferred function to run
	defer nagiosExitState.ReturnCheckResults()

	// Disable library debug logging output by default
	// vsphere.EnableLogging()
	vsphere.DisableLogging()

	// Setup configuration by parsing user-provided flags. Note plugin type so
	// that only applicable CLI flags are exposed and any plugin-specific
	// settings are applied.
	cfg, cfgErr := config.New(config.PluginType{Alarms: true})
	switch {
	case errors.Is(cfgErr, config.ErrVersionRequested):
		fmt.Println(config.Version())

		return

	case cfgErr != nil:
		// We're using the standalone Err function from rs/zerolog/log as we
		// do not have a working configuration.
		zlog.Err(cfgErr).Msg("Error initializing application")
		nagiosExitState.ServiceOutput = fmt.Sprintf(
			"%s: Error initializing application",
			nagios.StateCRITICALLabel,
		)
		nagiosExitState.LastError = cfgErr
		nagiosExitState.ExitStatusCode = nagios.StateCRITICALExitCode

		return
	}

	// Enable library-level logging if debug logging level is enabled app-wide
	if cfg.LoggingLevel == config.LogLevelDebug {
		vsphere.EnableLogging()
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout())
	defer cancel()

	// Record thresholds for use as Nagios "Long Service Output" content. This
	// content is shown in the detailed web UI and in notifications generated
	// by Nagios.
	// https://vdc-download.vmware.com/vmwb-repository/dcr-public/a5f4000f-1ea8-48a9-9221-586adff3c557/7ff50256-2cf2-45ea-aacd-87d231ab1ac7/vim.ManagedEntity.html#overallStatus
	nagiosExitState.CriticalThreshold = "One or more non-excluded alarms with a red status"
	nagiosExitState.WarningThreshold = "One or more non-excluded alarms with a yellow status"

	if cfg.EmitBranding {
		// If enabled, show application details at end of notification
		nagiosExitState.BrandingCallback = config.Branding("Notification generated by ")
	}

	log := cfg.Log.With().
		Str("included_resource_pools", cfg.IncludedResourcePools.String()).
		Str("excluded_resource_pools", cfg.ExcludedResourcePools.String()).
		Str("ignored_vms", cfg.IgnoredVMs.String()).
		Logger()

	log.Debug().Msg("Logging into vSphere environment")
	c, loginErr := vsphere.Login(
		ctx, cfg.Server, cfg.Port, cfg.TrustCert,
		cfg.Username, cfg.Domain, cfg.Password,
	)
	if loginErr != nil {
		log.Error().Err(loginErr).Msgf("error logging into %s", cfg.Server)

		nagiosExitState.LastError = loginErr
		nagiosExitState.ServiceOutput = fmt.Sprintf(
			"%s: Error logging into %q",
			nagios.StateCRITICALLabel,
			cfg.Server,
		)
		nagiosExitState.ExitStatusCode = nagios.StateCRITICALExitCode

		return
	}
	log.Debug().Msg("Successfully logged into vSphere environment")

	defer func() {
		if err := c.Logout(ctx); err != nil {
			log.Error().
				Err(err).
				Msg("failed to logout")
		}
	}()

	// At this point we're logged in, ready to process alarms.

	allAlarms, fetchAlarmsErr := vsphere.GetTriggeredAlarms(
		ctx,
		c,
		cfg.DatacenterName,
		true,
	)

	if fetchAlarmsErr != nil {
		log.Error().Err(fetchAlarmsErr).Msg("error retrieving alarms")

		nagiosExitState.LastError = fetchAlarmsErr
		nagiosExitState.ServiceOutput = fmt.Sprintf(
			"%s: Error retrieving alarms",
			nagios.StateCRITICALLabel,
		)
		nagiosExitState.ExitStatusCode = nagios.StateCRITICALExitCode

		return
	}

	log.Debug().Int("total_alarms", len(allAlarms)).Msg("alarms found")

	log.Debug().Msg("Filtering triggered alarms by entity type")
	filteredAlarms := vsphere.FilterTriggeredAlarmsByEntityType(
		allAlarms,
		cfg.IncludedAlarmEntityTypes,
		cfg.ExcludedAlarmEntityTypes,
	)

	log.Debug().Msg("Filtering triggered alarms by acknowledged state")
	filteredAlarms = vsphere.FilterTriggeredAlarmsByAcknowledgedState(
		filteredAlarms,
		cfg.EvaluateAcknowledgedAlarms,
	)

	log.Debug().
		Int("remaining_alarms", len(filteredAlarms)).
		Msg("alarms remaining after filtering")

	switch {

	case len(filteredAlarms) > 0:

		log.Error().
			Int("total_alarms", len(allAlarms)).
			Int("filtered_alarms", len(filteredAlarms)).
			Int("excluded_alarms", len(allAlarms)-len(filteredAlarms)).
			Msg("Non-excluded alarms detected")

		// Same error no matter whether CRITICAL, WARNING or UNKNOWN state.
		nagiosExitState.LastError = vsphere.ErrAlarmNotExcludedFromEvaluation

		// Set state label and exit code based on most severe
		// ManagedEntityStatus found in the TriggeredAlarms collection.
		var stateLabel string
		switch {
		case filteredAlarms.HasCriticalState():
			stateLabel = nagios.StateCRITICALLabel
			nagiosExitState.ExitStatusCode = nagios.StateCRITICALExitCode

		case filteredAlarms.HasWarningState():
			stateLabel = nagios.StateWARNINGLabel
			nagiosExitState.ExitStatusCode = nagios.StateWARNINGExitCode

		case filteredAlarms.HasUnknownState():
			stateLabel = nagios.StateUNKNOWNLabel
			nagiosExitState.ExitStatusCode = nagios.StateUNKNOWNExitCode
		}

		nagiosExitState.ServiceOutput = vsphere.AlarmsOneLineCheckSummary(
			stateLabel,
			allAlarms,
			filteredAlarms,
			cfg.IncludedAlarmEntityTypes,
			cfg.ExcludedAlarmEntityTypes,
		)

		nagiosExitState.LongServiceOutput = vsphere.AlarmsReport(
			c.Client,
			allAlarms,
			filteredAlarms,
			cfg.IncludedAlarmEntityTypes,
			cfg.ExcludedAlarmEntityTypes,
			cfg.EvaluateAcknowledgedAlarms,
		)

		return

	default:

		// success path

		log.Info().Msg("No non-excluded alarms detected")

		nagiosExitState.LastError = nil

		nagiosExitState.ServiceOutput = vsphere.AlarmsOneLineCheckSummary(
			nagios.StateOKLabel,
			allAlarms,
			filteredAlarms,
			cfg.IncludedAlarmEntityTypes,
			cfg.ExcludedAlarmEntityTypes,
		)

		nagiosExitState.LongServiceOutput = vsphere.AlarmsReport(
			c.Client,
			allAlarms,
			filteredAlarms,
			cfg.IncludedAlarmEntityTypes,
			cfg.ExcludedAlarmEntityTypes,
			cfg.EvaluateAcknowledgedAlarms,
		)

		nagiosExitState.ExitStatusCode = nagios.StateOKExitCode

		return

	}

}