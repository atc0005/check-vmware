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
	"github.com/vmware/govmomi/vim25/mo"

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
	cfg, cfgErr := config.New(config.PluginType{VirtualHardwareVersion: true})
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

	if cfg.EmitBranding {
		// If enabled, show application details at end of notification
		nagiosExitState.BrandingCallback = config.Branding("Notification generated by ")
	}

	log := cfg.Log.With().
		Str("included_resource_pools", cfg.IncludedResourcePools.String()).
		Str("excluded_resource_pools", cfg.ExcludedResourcePools.String()).
		Str("ignored_vms", cfg.IgnoredVMs.String()).
		Bool("eval_powered_off", cfg.PoweredOff).
		Logger()

	log.Debug().Msg("Logging into vSphere environment")
	c, loginErr := vsphere.Login(
		ctx, cfg.Server, cfg.Port, cfg.TrustCert,
		cfg.Username, cfg.Domain, cfg.Password,
		cfg.UserAgent(),
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

	// At this point we're logged in, ready to retrieve a list of VMs. If
	// specified, we should limit VMs based on include/exclude lists. First,
	// we'll make sure that all specified resource pools actually exist in the
	// vSphere environment.

	log.Debug().Msg("Validating resource pools")
	validateErr := vsphere.ValidateRPs(ctx, c.Client, cfg.IncludedResourcePools, cfg.ExcludedResourcePools)
	if validateErr != nil {
		log.Error().Err(validateErr).Msg("error validating include/exclude lists")

		nagiosExitState.LastError = validateErr
		nagiosExitState.ServiceOutput = fmt.Sprintf(
			"%s: Error validating include/exclude lists",
			nagios.StateCRITICALLabel,
		)
		nagiosExitState.ExitStatusCode = nagios.StateCRITICALExitCode

		return
	}

	log.Debug().Msg("Retrieving eligible resource pools")
	resourcePools, getRPsErr := vsphere.GetEligibleRPs(
		ctx,
		c.Client,
		cfg.IncludedResourcePools,
		cfg.ExcludedResourcePools,
		true,
	)
	if getRPsErr != nil {
		log.Error().Err(getRPsErr).Msg(
			"error retrieving list of resource pools",
		)

		nagiosExitState.LastError = getRPsErr
		nagiosExitState.ServiceOutput = fmt.Sprintf(
			"%s: Error retrieving list of resource pools from %q",
			nagios.StateCRITICALLabel,
			cfg.Server,
		)
		nagiosExitState.ExitStatusCode = nagios.StateCRITICALExitCode

		return
	}

	rpNames := make([]string, 0, len(resourcePools))
	for _, rp := range resourcePools {
		rpNames = append(rpNames, rp.Name)
	}

	log.Debug().
		Str("resource_pools", strings.Join(rpNames, ", ")).
		Msg("")

	log.Debug().Msg("Retrieving vms from eligible resource pools")
	rpEntityVals := make([]mo.ManagedEntity, 0, len(resourcePools))
	for i := range resourcePools {
		rpEntityVals = append(rpEntityVals, resourcePools[i].ManagedEntity)
	}
	vms, getVMsErr := vsphere.GetVMsFromContainer(ctx, c.Client, true, rpEntityVals...)
	if getVMsErr != nil {
		log.Error().Err(getVMsErr).Msg(
			"error retrieving list of VMs from resource pools list",
		)

		nagiosExitState.LastError = getVMsErr
		nagiosExitState.ServiceOutput = fmt.Sprintf(
			"%s: Error retrieving list of VMs from resource pools list",
			nagios.StateCRITICALLabel,
		)
		nagiosExitState.ExitStatusCode = nagios.StateCRITICALExitCode

		return
	}

	log.Debug().Msg("Drop any VMs we've been asked to exclude from checks")
	filteredVMs := vsphere.ExcludeVMsByName(vms, cfg.IgnoredVMs)

	log.Debug().Msg("Filter VMs to specified power state")
	filteredVMs = vsphere.FilterVMsByPowerState(filteredVMs, cfg.PoweredOff)

	log.Debug().
		Str("virtual_machines", strings.Join(vsphere.VMNames(filteredVMs), ", ")).
		Msg("Filtered VMs")

		// here we diverge from other plugins

	defaultHardwareVersion, getDefVerErr := vsphere.DefaultHardwareVersion(
		ctx,
		c.Client,
		cfg.HostSystemName,
		cfg.ClusterName,
		cfg.DatacenterName,
	)
	if getDefVerErr != nil {
		log.Error().Err(getDefVerErr).Msg(
			"error retrieving default hardware version",
		)

		nagiosExitState.LastError = getDefVerErr
		nagiosExitState.ServiceOutput = fmt.Sprintf(
			"%s: Error retrieving default hardware version",
			nagios.StateCRITICALLabel,
		)
		nagiosExitState.ExitStatusCode = nagios.StateCRITICALExitCode

		return
	}

	log.Debug().
		Int("default_hardware_version", defaultHardwareVersion.VersionNumber()).
		Msg("")

	hardwareVersionsIdx := make(vsphere.HardwareVersionsIndex)
	for _, vm := range filteredVMs {
		log.Debug().
			Str("vm_name", vm.Name).
			Str("hardware_version", vm.Config.Version).
			Msg("")

		// record the hardware version and count of that version
		hardwareVersionsIdx[vm.Config.Version]++
	}

	if cfg.VirtualHardwareApplyHomogeneousVersionCheck() {

		// Record thresholds for use as Nagios "Long Service Output" content. This
		// content is shown in the detailed web UI and in notifications generated
		// by Nagios.
		nagiosExitState.CriticalThreshold = config.ThresholdNotUsed
		nagiosExitState.WarningThreshold = "Non-homogenous hardware versions."

		switch {

		// There are at least two hardware versions present instead of a
		// uniform version across all VirtualMachines.
		case hardwareVersionsIdx.Count() > 1:

			log.Error().
				Int("vms_filtered", len(filteredVMs)).
				Int("unique_hardware_versions", hardwareVersionsIdx.Count()).
				Str("newest_hardware", hardwareVersionsIdx.Newest().String()).
				Str("outdated_hardware_list", strings.Join(
					hardwareVersionsIdx.Outdated().VersionNames(), ", ")).
				Msg("Virtual Hardware versions inconsistency detected")

			nagiosExitState.LastError = vsphere.ErrVirtualHardwareOutdatedVersionsFound

			nagiosExitState.ServiceOutput = vsphere.VirtualHardwareOneLineCheckSummary(
				nagios.StateWARNINGLabel,
				hardwareVersionsIdx,
				hardwareVersionsIdx.Newest().VersionNumber(),
				filteredVMs,
				resourcePools,
			)

			nagiosExitState.LongServiceOutput = vsphere.VirtualHardwareReport(
				c.Client,
				hardwareVersionsIdx,
				hardwareVersionsIdx.Newest().VersionNumber(),
				defaultHardwareVersion,
				vms,
				filteredVMs,
				cfg.IgnoredVMs,
				cfg.PoweredOff,
				cfg.IncludedResourcePools,
				cfg.ExcludedResourcePools,
				resourcePools,
			)

			nagiosExitState.ExitStatusCode = nagios.StateWARNINGExitCode

			return

		default:

			// same hardware version

			nagiosExitState.LastError = nil

			nagiosExitState.ServiceOutput = vsphere.VirtualHardwareOneLineCheckSummary(
				nagios.StateOKLabel,
				hardwareVersionsIdx,
				hardwareVersionsIdx.Newest().VersionNumber(),
				filteredVMs,
				resourcePools,
			)

			nagiosExitState.LongServiceOutput = vsphere.VirtualHardwareReport(
				c.Client,
				hardwareVersionsIdx,
				hardwareVersionsIdx.Newest().VersionNumber(),
				defaultHardwareVersion,
				vms,
				filteredVMs,
				cfg.IgnoredVMs,
				cfg.PoweredOff,
				cfg.IncludedResourcePools,
				cfg.ExcludedResourcePools,
				resourcePools,
			)

			nagiosExitState.ExitStatusCode = nagios.StateOKExitCode

			return

		}

	}

	if cfg.VirtualHardwareApplyMinVersionCheck() {

		// Record thresholds for use as Nagios "Long Service Output" content. This
		// content is shown in the detailed web UI and in notifications generated
		// by Nagios.
		nagiosExitState.CriticalThreshold = fmt.Sprintf(
			"Hardware versions older than the minimum (%d) present.",
			cfg.VirtualHardwareMinimumVersion,
		)
		nagiosExitState.WarningThreshold = config.ThresholdNotUsed

		hardwareVersions := hardwareVersionsIdx.Versions()

		switch {
		case !hardwareVersions.MeetsMinVersion(cfg.VirtualHardwareMinimumVersion):

			log.Error().
				Int("vms_filtered", len(filteredVMs)).
				Int("unique_hardware_versions", hardwareVersionsIdx.Count()).
				Str("newest_hardware", hardwareVersionsIdx.Newest().String()).
				Str("outdated_hardware_list", strings.Join(
					hardwareVersionsIdx.Outdated().VersionNames(), ", ")).
				Msg("Virtual Hardware versions older than the specified minimum version detected")

			nagiosExitState.LastError = vsphere.ErrVirtualHardwareOutdatedVersionsFound

			nagiosExitState.ServiceOutput = vsphere.VirtualHardwareOneLineCheckSummary(
				nagios.StateCRITICALLabel,
				hardwareVersionsIdx,
				cfg.VirtualHardwareMinimumVersion,
				filteredVMs,
				resourcePools,
			)

			nagiosExitState.LongServiceOutput = vsphere.VirtualHardwareReport(
				c.Client,
				hardwareVersionsIdx,
				cfg.VirtualHardwareMinimumVersion,
				defaultHardwareVersion,
				vms,
				filteredVMs,
				cfg.IgnoredVMs,
				cfg.PoweredOff,
				cfg.IncludedResourcePools,
				cfg.ExcludedResourcePools,
				resourcePools,
			)

			nagiosExitState.ExitStatusCode = nagios.StateCRITICALExitCode

			return

		default:
			nagiosExitState.LastError = nil

			nagiosExitState.ServiceOutput = vsphere.VirtualHardwareOneLineCheckSummary(
				nagios.StateOKLabel,
				hardwareVersionsIdx,
				cfg.VirtualHardwareMinimumVersion,
				filteredVMs,
				resourcePools,
			)

			nagiosExitState.LongServiceOutput = vsphere.VirtualHardwareReport(
				c.Client,
				hardwareVersionsIdx,
				cfg.VirtualHardwareMinimumVersion,
				defaultHardwareVersion,
				vms,
				filteredVMs,
				cfg.IgnoredVMs,
				cfg.PoweredOff,
				cfg.IncludedResourcePools,
				cfg.ExcludedResourcePools,
				resourcePools,
			)

			nagiosExitState.ExitStatusCode = nagios.StateOKExitCode

			return
		}

	}

	if cfg.VirtualHardwareApplyDefaultIsMinVersionCheck() {

		// Record thresholds for use as Nagios "Long Service Output" content. This
		// content is shown in the detailed web UI and in notifications generated
		// by Nagios.
		nagiosExitState.CriticalThreshold = config.ThresholdNotUsed

		nagiosExitState.WarningThreshold = fmt.Sprintf(
			"Hardware versions older than the default host or cluster (%d) present.",
			defaultHardwareVersion.VersionNumber(),
		)

		hardwareVersions := hardwareVersionsIdx.Versions()

		switch {
		case !hardwareVersions.MeetsMinVersion(defaultHardwareVersion.VersionNumber()):

			log.Error().
				Int("vms_filtered", len(filteredVMs)).
				Int("unique_hardware_versions", hardwareVersionsIdx.Count()).
				Str("newest_hardware", hardwareVersionsIdx.Newest().String()).
				Int("default_hardware_version", defaultHardwareVersion.VersionNumber()).
				Str("outdated_hardware_list", strings.Join(
					hardwareVersionsIdx.Outdated().VersionNames(), ", ")).
				Msg("Virtual Hardware versions older than the host or cluster default version detected")

			nagiosExitState.LastError = vsphere.ErrVirtualHardwareOutdatedVersionsFound

			nagiosExitState.ServiceOutput = vsphere.VirtualHardwareOneLineCheckSummary(
				nagios.StateWARNINGLabel,
				hardwareVersionsIdx,
				defaultHardwareVersion.VersionNumber(),
				filteredVMs,
				resourcePools,
			)

			nagiosExitState.LongServiceOutput = vsphere.VirtualHardwareReport(
				c.Client,
				hardwareVersionsIdx,
				defaultHardwareVersion.VersionNumber(),
				defaultHardwareVersion,
				vms,
				filteredVMs,
				cfg.IgnoredVMs,
				cfg.PoweredOff,
				cfg.IncludedResourcePools,
				cfg.ExcludedResourcePools,
				resourcePools,
			)

			nagiosExitState.ExitStatusCode = nagios.StateWARNINGExitCode

			return

		default:
			nagiosExitState.LastError = nil

			nagiosExitState.ServiceOutput = vsphere.VirtualHardwareOneLineCheckSummary(
				nagios.StateOKLabel,
				hardwareVersionsIdx,
				defaultHardwareVersion.VersionNumber(),
				filteredVMs,
				resourcePools,
			)

			nagiosExitState.LongServiceOutput = vsphere.VirtualHardwareReport(
				c.Client,
				hardwareVersionsIdx,
				defaultHardwareVersion.VersionNumber(),
				defaultHardwareVersion,
				vms,
				filteredVMs,
				cfg.IgnoredVMs,
				cfg.PoweredOff,
				cfg.IncludedResourcePools,
				cfg.ExcludedResourcePools,
				resourcePools,
			)

			nagiosExitState.ExitStatusCode = nagios.StateOKExitCode

			return
		}

	}

	if cfg.VirtualHardwareApplyOutdatedByVersionCheck() {

		// Record thresholds for use as Nagios "Long Service Output" content. This
		// content is shown in the detailed web UI and in notifications generated
		// by Nagios.
		nagiosExitState.CriticalThreshold = fmt.Sprintf(
			"Hardware versions outdated by more than %d versions present.",
			cfg.VirtualHardwareOutdatedByCritical,
		)
		nagiosExitState.WarningThreshold = fmt.Sprintf(
			"Hardware versions outdated by more than %d versions present.",
			cfg.VirtualHardwareOutdatedByWarning,
		)

		hardwareVersions := hardwareVersionsIdx.Versions()
		latestHWVerNum := hardwareVersionsIdx.Newest().VersionNumber()
		criticalThresholdVerNum := latestHWVerNum - cfg.VirtualHardwareOutdatedByCritical
		warningThresholdVerNum := latestHWVerNum - cfg.VirtualHardwareOutdatedByWarning

		switch {
		case !hardwareVersions.MeetsMinVersion(criticalThresholdVerNum):

			log.Error().
				Int("vms_filtered", len(filteredVMs)).
				Int("unique_hardware_versions", hardwareVersionsIdx.Count()).
				Str("newest_hardware", hardwareVersionsIdx.Newest().String()).
				Str("outdated_hardware_list", strings.Join(
					hardwareVersionsIdx.Outdated().VersionNames(), ", ")).
				Msg("Virtual Hardware versions older than the specified minimum version detected")

			nagiosExitState.LastError = vsphere.ErrVirtualHardwareOutdatedVersionsFound

			nagiosExitState.ServiceOutput = vsphere.VirtualHardwareOneLineCheckSummary(
				nagios.StateCRITICALLabel,
				hardwareVersionsIdx,
				criticalThresholdVerNum,
				filteredVMs,
				resourcePools,
			)

			nagiosExitState.LongServiceOutput = vsphere.VirtualHardwareReport(
				c.Client,
				hardwareVersionsIdx,
				criticalThresholdVerNum,
				defaultHardwareVersion,
				vms,
				filteredVMs,
				cfg.IgnoredVMs,
				cfg.PoweredOff,
				cfg.IncludedResourcePools,
				cfg.ExcludedResourcePools,
				resourcePools,
			)

			nagiosExitState.ExitStatusCode = nagios.StateCRITICALExitCode

			return

		case !hardwareVersions.MeetsMinVersion(warningThresholdVerNum):

			log.Error().
				Int("vms_filtered", len(filteredVMs)).
				Int("unique_hardware_versions", hardwareVersionsIdx.Count()).
				Str("newest_hardware", hardwareVersionsIdx.Newest().String()).
				Str("outdated_hardware_list", strings.Join(
					hardwareVersionsIdx.Outdated().VersionNames(), ", ")).
				Msg("Virtual Hardware versions older than the specified minimum version detected")

			nagiosExitState.LastError = vsphere.ErrVirtualHardwareOutdatedVersionsFound

			nagiosExitState.ServiceOutput = vsphere.VirtualHardwareOneLineCheckSummary(
				nagios.StateWARNINGLabel,
				hardwareVersionsIdx,
				warningThresholdVerNum,
				filteredVMs,
				resourcePools,
			)

			nagiosExitState.LongServiceOutput = vsphere.VirtualHardwareReport(
				c.Client,
				hardwareVersionsIdx,
				warningThresholdVerNum,
				defaultHardwareVersion,
				vms,
				filteredVMs,
				cfg.IgnoredVMs,
				cfg.PoweredOff,
				cfg.IncludedResourcePools,
				cfg.ExcludedResourcePools,
				resourcePools,
			)

			nagiosExitState.ExitStatusCode = nagios.StateWARNINGExitCode

			return

		default:
			nagiosExitState.LastError = nil

			nagiosExitState.ServiceOutput = vsphere.VirtualHardwareOneLineCheckSummary(
				nagios.StateOKLabel,
				hardwareVersionsIdx,
				warningThresholdVerNum,
				filteredVMs,
				resourcePools,
			)

			nagiosExitState.LongServiceOutput = vsphere.VirtualHardwareReport(
				c.Client,
				hardwareVersionsIdx,
				warningThresholdVerNum,
				defaultHardwareVersion,
				vms,
				filteredVMs,
				cfg.IgnoredVMs,
				cfg.PoweredOff,
				cfg.IncludedResourcePools,
				cfg.ExcludedResourcePools,
				resourcePools,
			)

			nagiosExitState.ExitStatusCode = nagios.StateOKExitCode

			return
		}

	}

}
