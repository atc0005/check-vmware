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
	"strings"

	"github.com/atc0005/go-nagios"

	"github.com/atc0005/check-vmware/internal/config"
	"github.com/atc0005/check-vmware/internal/vsphere"

	zlog "github.com/rs/zerolog/log"
)

//go:generate go-winres make --product-version=git-tag --file-version=git-tag

func main() {

	plugin := nagios.NewPlugin()

	// defer this from the start so it is the last deferred function to run
	defer plugin.ReturnCheckResults()

	// Annotate all errors (if any) with remediation advice just before ending
	// plugin execution.
	defer vsphere.AnnotateError(plugin)

	// Setup configuration by parsing user-provided flags. Note plugin type so
	// that only applicable CLI flags are exposed and any plugin-specific
	// settings are applied.
	cfg, cfgErr := config.New(config.PluginType{VirtualMachineLastBackupViaCA: true})
	switch {
	case errors.Is(cfgErr, config.ErrVersionRequested):
		fmt.Println(config.Version())

		return

	case cfgErr != nil:
		// We're using the standalone Err function from rs/zerolog/log as we
		// do not have a working configuration.
		zlog.Err(cfgErr).Msg("Error initializing application")
		plugin.ServiceOutput = fmt.Sprintf(
			"%s: Error initializing application",
			nagios.StateUNKNOWNLabel,
		)
		plugin.AddError(cfgErr)
		plugin.ExitStatusCode = nagios.StateUNKNOWNExitCode

		return
	}

	// Enable library-level logging if debug or greater logging level is
	// enabled app-wide.
	handleLibraryLogging()

	// Set context deadline equal to user-specified timeout value for plugin
	// runtime/execution.
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout())
	defer cancel()

	// Record thresholds for use as Nagios "Long Service Output" content. This
	// content is shown in the detailed web UI and in notifications generated
	// by Nagios.
	plugin.CriticalThreshold = fmt.Sprintf(
		"non-excluded VM with: %s\t"+strings.Join(
			[]string{
				"backup date exceeding specified CRITICAL threshold",
			},
			nagios.CheckOutputEOL+"\t",
		),
		nagios.CheckOutputEOL,
	)
	plugin.WarningThreshold = fmt.Sprintf(
		"non-excluded VM with: %s\t"+strings.Join(
			[]string{
				"backup date exceeding specified WARNING threshold, but not CRITICAL threshold",
				"backup date missing",
				"backup date does not match default/user-specified format",
			},
			nagios.CheckOutputEOL+"\t",
		),
		nagios.CheckOutputEOL,
	)

	if cfg.EmitBranding {
		// If enabled, show application details at end of notification
		plugin.BrandingCallback = config.Branding("Notification generated by ")
	}

	log := cfg.Log.With().
		Str("included_resource_pools", cfg.IncludedResourcePools.String()).
		Str("excluded_resource_pools", cfg.ExcludedResourcePools.String()).
		Str("ignored_vms", cfg.IgnoredVMs.String()).
		Int("backup_age_critical", cfg.VMBackupAgeCritical).
		Int("backup_age_warning", cfg.VMBackupAgeWarning).
		Logger()

	log.Debug().Msg("Logging into vSphere environment")
	c, loginErr := vsphere.Login(
		ctx, cfg.Server, cfg.Port, cfg.TrustCert,
		cfg.Username, cfg.Domain, cfg.Password,
		cfg.UserAgent(),
	)
	if loginErr != nil {
		log.Error().Err(loginErr).Msgf("error logging into %s", cfg.Server)

		plugin.AddError(loginErr)
		plugin.ServiceOutput = fmt.Sprintf(
			"%s: Error logging into %q",
			nagios.StateCRITICALLabel,
			cfg.Server,
		)
		plugin.ExitStatusCode = nagios.StateCRITICALExitCode

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

	log.Debug().Msg("Performing initial filtering of vms")
	vmsFilterOptions := vsphere.VMsFilterOptions{
		ResourcePoolsIncluded:       cfg.IncludedResourcePools,
		ResourcePoolsExcluded:       cfg.ExcludedResourcePools,
		FoldersIncluded:             cfg.IncludedFolders,
		FoldersExcluded:             cfg.ExcludedFolders,
		VirtualMachineNamesExcluded: cfg.IgnoredVMs,

		// NOTE: This plugin is hard-coded to evaluate powered off and powered
		// on VMs equally. I'm not sure whether ignoring powered off VMs by
		// default makes sense for this particular plugin.
		//
		// Please share your feedback here if you feel differently:
		// https://github.com/atc0005/check-vmware/discussions
		//
		// Please expand on some use cases for ignoring powered off VMs by
		// default.
		// IncludePoweredOff:           cfg.PoweredOff,
		IncludePoweredOff: true,
	}
	vmsFilterResults, vmsFilterErr := vsphere.FilterVMs(
		ctx,
		c.Client,
		vmsFilterOptions,
	)
	if vmsFilterErr != nil {
		log.Error().Err(vmsFilterErr).Msg(
			"error filtering VMs",
		)

		plugin.AddError(vmsFilterErr)
		plugin.ServiceOutput = fmt.Sprintf(
			"%s: Error filtering VMs",
			nagios.StateCRITICALLabel,
		)
		plugin.ExitStatusCode = nagios.StateCRITICALExitCode

		return
	}
	log.Debug().Msg("Finished initial filtering of vms")

	// Here we diverge from most other plugins in this project

	vmsWithBackup, vmsLookupErr := vsphere.GetVMsWithBackup(
		vmsFilterResults.VMsAfterFiltering(),
		cfg.VMBackupDateTimezone,
		cfg.VMBackupDateCustomAttribute,
		cfg.VMBackupMetadataCustomAttribute,
		cfg.VMBackupDateFormat,
		cfg.VMBackupAgeCritical,
		cfg.VMBackupAgeWarning,
	)
	if vmsLookupErr != nil {

		log.Error().Err(vmsLookupErr).
			Msg("error retrieving virtual machines with requested backup custom attributes")

		plugin.AddError(vmsLookupErr)
		plugin.ServiceOutput = fmt.Sprintf(
			"%s: Error retrieving virtual machines with requested backup custom attributes",
			nagios.StateCRITICALLabel,
		)
		plugin.ExitStatusCode = nagios.StateCRITICALExitCode

		return

	}

	log.Debug().Msg("Compiling Performance Data details")

	pd := append(
		vsphere.VMFilterResultsPerfData(vmsFilterResults),
		[]nagios.PerformanceData{
			// The `time` (runtime) metric is appended at plugin exit, so do not
			// duplicate it here.
			{
				Label: "vms_with_backup_dates",
				Value: fmt.Sprintf("%d", vmsWithBackup.NumBackups()),
			},
			{
				Label: "vms_without_backup_dates",
				Value: fmt.Sprintf("%d", vmsWithBackup.NumWithoutBackups()),
			},
		}...,
	)

	if err := plugin.AddPerfData(false, pd...); err != nil {
		log.Error().
			Err(err).
			Msg("failed to add performance data")

		// Surface the error in plugin output.
		plugin.AddError(err)

		plugin.ExitStatusCode = nagios.StateUNKNOWNExitCode
		plugin.ServiceOutput = fmt.Sprintf(
			"%s: Failed to process performance data metrics",
			nagios.StateUNKNOWNLabel,
		)

		return
	}

	// Update logger with new performance data related fields
	log = log.With().
		Int("resource_pools_evaluated", vmsFilterResults.NumRPsAfterFiltering()).
		Int("vms_total", vmsFilterResults.NumVMsAll()).
		Int("vms_after_filtering", vmsFilterResults.NumVMsAfterFiltering()).
		Int("vms_excluded_by_name", vmsFilterResults.NumVMsExcludedByName()).
		Int("vms_excluded_by_power_state", vmsFilterResults.NumVMsExcludedByPowerState()).
		Int("vms_with_backup_dates", vmsWithBackup.NumBackups()).
		Int("vms_without_backup_dates", vmsWithBackup.NumWithoutBackups()).
		Logger()

	switch {
	case vmsWithBackup.IsCriticalState() || vmsWithBackup.IsWarningState():

		plugin.AddError(func() error {
			switch {

			// Something prevented a regularly scheduled backup from
			// running/completing.
			//
			// We consider this error to be of a higher priority, so we check
			// for it first before we look for missing backups.
			case vmsWithBackup.HasOldBackup():
				return vsphere.ErrVirtualMachineBackupDateOld

			// One or more of the non-excluded VMs does not have a backup
			// associated with it (for whatever reason).
			case !vmsWithBackup.AllHasBackup():
				return vsphere.ErrVirtualMachineMissingBackupDate

			default:
				return errors.New("unknown error state; please report this")

			}
		}())

		stateLabel := nagios.StateCRITICALLabel
		stateExitCode := nagios.StateCRITICALExitCode
		if vmsWithBackup.IsWarningState() {
			stateLabel = nagios.StateWARNINGLabel
			stateExitCode = nagios.StateWARNINGExitCode
		}

		plugin.ServiceOutput = vsphere.VMBackupViaCAOneLineCheckSummary(
			stateLabel,
			vmsFilterResults,
			vmsWithBackup,
		)

		plugin.LongServiceOutput = vsphere.VMBackupViaCAReport(
			c.Client,
			vmsFilterOptions,
			vmsFilterResults,
			vmsWithBackup,
		)

		plugin.ExitStatusCode = stateExitCode

	default:

		// success if we made it here

		log.Debug().Msg("No non-excluded VMs with old or missing backups detected")

		plugin.ServiceOutput = vsphere.VMBackupViaCAOneLineCheckSummary(
			nagios.StateOKLabel,
			vmsFilterResults,
			vmsWithBackup,
		)

		plugin.LongServiceOutput = vsphere.VMBackupViaCAReport(
			c.Client,
			vmsFilterOptions,
			vmsFilterResults,
			vmsWithBackup,
		)

		plugin.ExitStatusCode = nagios.StateOKExitCode

	}
}
