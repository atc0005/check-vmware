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
	"github.com/vmware/govmomi/units"
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
	cfg, cfgErr := config.New(config.PluginType{ResourcePoolsMemory: true})
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
	nagiosExitState.CriticalThreshold = fmt.Sprintf(
		"%d%% usage of %d GB memory",
		cfg.ResourcePoolsMemoryUseCritical,
		cfg.ResourcePoolsMemoryMaxAllowed,
	)

	nagiosExitState.WarningThreshold = fmt.Sprintf(
		"%d%% usage of %d GB memory",
		cfg.ResourcePoolsMemoryUseWarning,
		cfg.ResourcePoolsMemoryMaxAllowed,
	)

	if cfg.EmitBranding {
		// If enabled, show application details at end of notification
		nagiosExitState.BrandingCallback = config.Branding("Notification generated by ")
	}

	// Explicitly ignore the default `Resources` resource pool so that we only
	// use the Resource Pools specified by the sysadmin.
	if err := cfg.ExcludedResourcePools.Set(vsphere.ParentResourcePool); err != nil {
		// We're using the standalone Err function from rs/zerolog/log as we
		// have not created our custom `log` zerolog.Logger instance yet.
		zlog.Err(cfgErr).Msg("Error excluding default Resources Pool from evaluation")
		nagiosExitState.ServiceOutput = fmt.Sprintf(
			"%s: Error excluding default Resources Pool from evaluation",
			nagios.StateCRITICALLabel,
		)
		nagiosExitState.LastError = err
		nagiosExitState.ExitStatusCode = nagios.StateCRITICALExitCode

		return
	}

	log := cfg.Log.With().
		Str("included_resource_pools", cfg.IncludedResourcePools.String()).
		Str("excluded_resource_pools", cfg.ExcludedResourcePools.String()).
		Int("max_memory_usage_allowed", cfg.ResourcePoolsMemoryMaxAllowed).
		Int("memory_usage_critical", cfg.ResourcePoolsMemoryUseCritical).
		Int("memory_usage_warning", cfg.ResourcePoolsMemoryUseWarning).
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

	var aggregateMemoryUsage int64
	for _, rp := range resourcePools {
		// Per vSphere API docs, `rp.Runtime.Memory.OverallUsage` was
		// deprecated in v6.5, so we use `hostMemoryUsage` instead.
		rpMemoryUsage := rp.Summary.GetResourcePoolSummary().QuickStats.HostMemoryUsage * units.MB
		aggregateMemoryUsage += rpMemoryUsage
		log.Debug().
			Str("resource_pool_name", rp.Name).
			Str("resource_pool_memory_usage", units.ByteSize(rpMemoryUsage).String()).
			Msg("")
	}

	clusterMemory, getMemErr := vsphere.GetHostSystemsTotalMemory(ctx, c.Client, false)
	if getMemErr != nil {
		log.Error().Err(getMemErr).Msg(
			"error retrieving hosts memory capacity",
		)

		nagiosExitState.LastError = getMemErr
		nagiosExitState.ServiceOutput = fmt.Sprintf(
			"%s: Error retrieving memory capacity of hosts from %q",
			nagios.StateCRITICALLabel,
			cfg.Server,
		)
		nagiosExitState.ExitStatusCode = nagios.StateCRITICALExitCode

		return
	}

	clusterMemoryInGB := clusterMemory / units.GB
	memoryPercentageUsedOfClusterCapacity := vsphere.MemoryUsedPercentage(
		aggregateMemoryUsage,
		int(clusterMemoryInGB),
	)

	log.Debug().
		Int64("cluster_memory_bytes", clusterMemory).
		Int64("cluster_memory_gb", clusterMemoryInGB).
		Str("cluster_memory_hr", units.ByteSize(clusterMemory).String()).
		Float64("percent_memory_used_from_cluster_raw", memoryPercentageUsedOfClusterCapacity).
		Str("percent_memory_used_from_cluster_hr", fmt.Sprintf("%0.2f", memoryPercentageUsedOfClusterCapacity)).
		Msg("")

	log.Debug().
		Int64("aggregate_memory_usage_raw", aggregateMemoryUsage).
		Str("aggregate_memory_usage_human_readable", units.ByteSize(aggregateMemoryUsage).String()).
		Msg("Finished evaluating Resource Pool memory usage")

	memoryPercentageUsedOfAllowed := vsphere.MemoryUsedPercentage(aggregateMemoryUsage, cfg.ResourcePoolsMemoryMaxAllowed)
	var memoryRemaining int64

	switch {
	case aggregateMemoryUsage > int64(cfg.ResourcePoolsMemoryMaxAllowed):
		memoryRemaining = 0
	default:
		memoryRemaining = int64(cfg.ResourcePoolsMemoryMaxAllowed) - aggregateMemoryUsage
	}

	log.Debug().
		Float64("memory_percent_used", memoryPercentageUsedOfAllowed).
		Int64("memory_remaining", memoryRemaining).
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

	switch {
	case memoryPercentageUsedOfAllowed > float64(cfg.ResourcePoolsMemoryUseCritical):

		log.Error().
			Float64("memory_percent_used", memoryPercentageUsedOfAllowed).
			Int64("memory_remaining", memoryRemaining).
			Msg("memory usage critical")

		nagiosExitState.LastError = vsphere.ErrResourcePoolMemoryUsageThresholdCrossed

		nagiosExitState.ServiceOutput = vsphere.RPMemoryUsageOneLineCheckSummary(
			nagios.StateCRITICALLabel,
			aggregateMemoryUsage,
			cfg.ResourcePoolsMemoryMaxAllowed,
			clusterMemoryInGB,
			resourcePools,
		)

		nagiosExitState.LongServiceOutput = vsphere.ResourcePoolsMemoryReport(
			c.Client,
			aggregateMemoryUsage,
			cfg.ResourcePoolsMemoryMaxAllowed,
			clusterMemoryInGB,
			cfg.IncludedResourcePools,
			cfg.ExcludedResourcePools,
			resourcePools,
			vms,
		)

		nagiosExitState.ExitStatusCode = nagios.StateCRITICALExitCode

		return

	case memoryPercentageUsedOfAllowed > float64(cfg.ResourcePoolsMemoryUseWarning):

		log.Error().
			Float64("memory_percent_used", memoryPercentageUsedOfAllowed).
			Int64("memory_remaining", memoryRemaining).
			Msg("memory usage warning")

		nagiosExitState.LastError = vsphere.ErrResourcePoolMemoryUsageThresholdCrossed

		nagiosExitState.ServiceOutput = vsphere.RPMemoryUsageOneLineCheckSummary(
			nagios.StateWARNINGLabel,
			aggregateMemoryUsage,
			cfg.ResourcePoolsMemoryMaxAllowed,
			clusterMemoryInGB,
			resourcePools,
		)

		nagiosExitState.LongServiceOutput = vsphere.ResourcePoolsMemoryReport(
			c.Client,
			aggregateMemoryUsage,
			cfg.ResourcePoolsMemoryMaxAllowed,
			clusterMemoryInGB,
			cfg.IncludedResourcePools,
			cfg.ExcludedResourcePools,
			resourcePools,
			vms,
		)

		nagiosExitState.ExitStatusCode = nagios.StateWARNINGExitCode

		return

	default:

		nagiosExitState.LastError = nil

		nagiosExitState.ServiceOutput = vsphere.RPMemoryUsageOneLineCheckSummary(
			nagios.StateOKLabel,
			aggregateMemoryUsage,
			cfg.ResourcePoolsMemoryMaxAllowed,
			clusterMemoryInGB,
			resourcePools,
		)

		nagiosExitState.LongServiceOutput = vsphere.ResourcePoolsMemoryReport(
			c.Client,
			aggregateMemoryUsage,
			cfg.ResourcePoolsMemoryMaxAllowed,
			clusterMemoryInGB,
			cfg.IncludedResourcePools,
			cfg.ExcludedResourcePools,
			resourcePools,
			vms,
		)

		nagiosExitState.ExitStatusCode = nagios.StateOKExitCode

		return

	}

}
