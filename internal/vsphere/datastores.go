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
	"io"
	"sort"
	"strings"
	"time"

	"github.com/atc0005/check-vmware/internal/textutils"
	"github.com/atc0005/go-nagios"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/units"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

// ErrDatastoreInaccessible indicates that a specified datastore is marked as
// inaccessible.
var ErrDatastoreInaccessible = errors.New("datastore is inaccessible")

// ErrDatastoreUsageThresholdCrossed indicates that a specified
// datastore has exceeded a given threshold.
var ErrDatastoreUsageThresholdCrossed = errors.New("datastore usage exceeds specified threshold")

// ErrDatastoreLatencyThresholdCrossed indicates that a specified datastore
// has exceeded a given latency threshold.
var ErrDatastoreLatencyThresholdCrossed = errors.New("datastore latency exceeds specified threshold")

// ErrDatastoreLatencyAllMetricSetsZero indicates that all performance
// metric metric sets for a specified datastore are of value 0.
var ErrDatastoreLatencyAllMetricSetsZero = errors.New("datastore latency metric sets are all value zero")

// ErrDatastoreIormConfigurationPropertyUnavailable indicates that the
// IORMConfigInfo property is not available for evaluation. Without this
// property, plugins in this project are unable to reliably determine whether
// datastore performance statistics are being gathered.
//
// https://vdc-download.vmware.com/vmwb-repository/dcr-public/b50dcbbf-051d-4204-a3e7-e1b618c1e384/538cf2ec-b34f-4bae-a332-3820ef9e7773/vim.StorageResourceManager.IORMConfigInfo.html
var ErrDatastoreIormConfigurationPropertyUnavailable = errors.New(
	"datastore storage I/O resource management configuration property unavailable",
)

// ErrDatastoreStatsCollectionPropertyUnavailable indicates that the
// statsCollectionEnabled property is not available for evaluation. Without
// access to this property, plugins in this project are unable to reliably
// determine whether datastore performance statistics are being gathered.
//
// https://vdc-download.vmware.com/vmwb-repository/dcr-public/b50dcbbf-051d-4204-a3e7-e1b618c1e384/538cf2ec-b34f-4bae-a332-3820ef9e7773/vim.StorageResourceManager.IORMConfigInfo.html
var ErrDatastoreStatsCollectionPropertyUnavailable = errors.New(
	"datastore storage I/O statistics collection property unavailable",
)

// ErrDatastoreIormConfigurationStatisticsCollectionDisabled indicates that
// I/O (Iops, Latency) statistics collection is disabled. The administrators
// of the vSphere environment must enable the Statistics Collection option for
// each datastore that is monitored by plugins in this project.
//
// https://vdc-download.vmware.com/vmwb-repository/dcr-public/b50dcbbf-051d-4204-a3e7-e1b618c1e384/538cf2ec-b34f-4bae-a332-3820ef9e7773/vim.StorageResourceManager.IORMConfigInfo.html
var ErrDatastoreIormConfigurationStatisticsCollectionDisabled = errors.New(
	"datastore storage I/O statistics collection disabled",
)

// ErrDatastorePerformanceMetricsMissing indicates that no datastore
// performance metrics are available.
//
// This is believed to occur when a datastore is newly created and metrics
// have not yet been collected. For a long-lived datastore this is a
// problematic scenario, but for a new datastore it is not unexpected.
var ErrDatastorePerformanceMetricsMissing = errors.New("datastore performance metrics results are unavailable")

// ErrDatastorePerformancePercentileUnavailable is returned when a specific
// performance percentile is requested, but is unavailable in the results
// returned from Datastore Performance Summary query.
var ErrDatastorePerformancePercentileUnavailable = errors.New(
	"datastore performance percentile unavailable",
)

// DatastoreIDToNameIndex maps a Datastore's ID value to its name.
type DatastoreIDToNameIndex map[string]string

// DatastoreVMs provides an overview of all (visible) VirtualMachines residing
// on a specific Datastore.
type DatastoreVMs []DatastoreVM

// DatastoreVM is a summary of details for a VirtualMachine found on a
// specific datastore.
type DatastoreVM struct {

	// Name is the display name of the VirtualMachine.
	Name string

	// VMSize is the human readable or formatted size of the VirtualMachine.
	VMSize string

	// DatastoreUsage is the human readable or formatted percentage of the
	// Datastore space consumed by this VirtualMachine.
	DatastoreUsage string

	// PowerState tracks the current power state for a VirtualMachine.
	PowerState types.VirtualMachinePowerState

	// DatastoreMOID is the MOID or MoRef ID for the Datastore where this
	// VirtualMachine resides.
	DatastoreMOID types.ManagedObjectReference
}

// DatastoreUsageSummary tracks usage details for a specific Datastore.
type DatastoreUsageSummary struct {
	Datastore               mo.Datastore
	StorageRemainingPercent float64
	StorageUsedPercent      float64
	StorageTotal            int64
	StorageUsed             int64
	StorageRemaining        int64
	CriticalThreshold       int
	WarningThreshold        int
	VMs                     DatastoreVMs
}

// DatastorePerformanceThresholds is a collection of threshold values used to
// determine the state of latency metrics for a specified Datastore.
type DatastorePerformanceThresholds struct {

	// ReadLatencyWarning is the read latency in ms when a WARNING threshold
	// is reached.
	ReadLatencyWarning float64

	// ReadLatencyCritical is the read latency in ms when a CRITICAL threshold
	// is reached.
	ReadLatencyCritical float64

	// WriteLatencyWarning is the write latency in ms when a WARNING threshold
	// is reached.
	WriteLatencyWarning float64

	// WriteLatencyCritical is the write latency in ms when a CRITICAL
	// threshold is reached.
	WriteLatencyCritical float64

	// VMLatencyWarning is the latency in ms as observed by VMs using the
	// datastore when a WARNING threshold is reached.
	VMLatencyWarning float64

	// VMLatencyCritical is the latency in ms as observed by VMs using the
	// datastore when a CRITICAL threshold is reached.
	VMLatencyCritical float64
}

// DatastorePerformanceThresholdsIndex is an index of Datastore Performance
// metrics percentile to DatastorePerformanceThresholds.
type DatastorePerformanceThresholdsIndex map[int]DatastorePerformanceThresholds

// DatastorePerformanceSummaryIndex is an index of Datastore performance
// metric percentile to DatastorePerformanceSummary values. At any time there
// is an active or "live" interval for aggregated metrics.
type DatastorePerformanceSummaryIndex struct {

	// Entries is a map of percentile to Datastore Performance Summary
	// metrics.
	Entries map[int]DatastorePerformanceSummary

	// Active indicates whether the associated performance metrics are for the
	// active or "live" Interval or "window" of time.
	Active bool
}

// DatastorePerformanceSummaryIntervals is a collection of
// DatastorePerformanceSummaryIndex values. Each element of this collection
// represents an interval or window of time where metrics were collected. The
// first element contains metrics for the active interval where metrics are
// being actively aggregated.
type DatastorePerformanceSummaryIntervals []DatastorePerformanceSummaryIndex

// DatastorePerformanceSet is set of performance metrics for a specific
// Datastore. A collection of indexes is used to group metrics based on
// intervals or windows of time (e.g., X days worth of metrics, further
// grouped by percentile). Each group of metrics, known as a
// DatastorePerformanceSummary contains the thresholds necessary to evaluate
// metrics and determine overall plugin state.
//
// FIXME: Not satisfied with the name. Revisit this.
//
// DatastorePerformance, DatastorePerformanceSummaryCollection,
// DatastorePerformanceSet?
type DatastorePerformanceSet struct {

	// Datastore is the Managed Object associated with the collected Datastore
	// performance metrics.
	Datastore mo.Datastore

	// VMs provides a summary of details for all VirtualMachines found on the
	// specified datastore.
	VMs DatastoreVMs

	// Intervals is a collection of percentile to Datastore performance metric
	// indexes. Each element of this collection represents an interval or
	// window of time where metrics were collected. The first element contains
	// metrics for the active interval which are being actively aggregated.
	Intervals DatastorePerformanceSummaryIntervals
}

// DatastorePerformanceSummary tracks performance metrics for a specific
// Datastore.
//
// https://vdc-download.vmware.com/vmwb-repository/dcr-public/bf660c0a-f060-46e8-a94d-4b5e6ffc77ad/208bc706-e281-49b6-a0ce-b402ec19ef82/SDK/vsphere-ws/docs/ReferenceGuide/vim.StorageResourceManager.html#queryDatastorePerformanceSummary
// https://vdc-download.vmware.com/vmwb-repository/dcr-public/b50dcbbf-051d-4204-a3e7-e1b618c1e384/538cf2ec-b34f-4bae-a332-3820ef9e7773/vim.StorageResourceManager.StoragePerformanceSummary.html
type DatastorePerformanceSummary struct {

	// DatastoreID is the Managed Object Reference (MOID or MoRef ID)
	// associated with the collected Datastore performance metrics.
	DatastoreMOID types.ManagedObjectReference

	// ReadIops is the aggregated datastore Read I/O rate (reads/second).
	ReadIops float64

	// ReadLatency is the aggregated datastore latency in milliseconds for
	// read operations.
	ReadLatency float64

	// VMLatency is the aggregated datastore latency as observed by
	// VirtualMachines using the datastore. The reported latency is in
	// milliseconds.
	VMLatency float64

	// WriteIops is the aggregated datastore Write I/O rate (writes/second).
	WriteIops float64

	// WriteLatency is the aggregated datastore latency in milliseconds for
	// write operations.
	WriteLatency float64

	// Percentile is the metric percentile specification.
	//
	// A percentile is a value between 1 and 100. Each metric value in this
	// type corresponds with the percentile value in this field. For example,
	// if the value of percentile is P, and the value of the ReadLatency is L,
	// then P% of all the read IOs performed during observation Interval is
	// less than L milliseconds.
	Percentile int32

	// Interval is the time period over which statistics are aggregated. The
	// reported time unit is in seconds.
	//
	// NOTE: By observation *only*, this appears to represent up to (roughly)
	// a full 24 hours before a new metrics collection entry is created. For
	// days where metrics have not been collected this is 0, for days where
	// stats have been collected values such as these have been noted:
	//
	// 86020, 85980, 86060, 83860
	Interval int32

	// thresholds is a collection of specified Datastore performance metric
	// thresholds. If not specified, this value is nil.
	thresholds *DatastorePerformanceThresholds
}

// printVMSummary is a helper function used by Datastore report functions to
// generate summary information for a collection of Virtual Machines present
// on a specific datastore.
func printVMSummary(w io.Writer, vms DatastoreVMs, powerState types.VirtualMachinePowerState) {

	// Skip efforts to list VM summary details if there is nothing to show.
	if len(vms) == 0 {
		return
	}

	var powerStateVMs int
	switch powerState {
	case types.VirtualMachinePowerStatePoweredOn:
		powerStateVMs = vms.NumVMsPoweredOn()
	default:
		powerStateVMs = vms.NumVMsPoweredOff()
	}

	sectionHeader := fmt.Sprintf(
		"%d %s VMs on datastore:%s%s",
		powerStateVMs,
		powerState,
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
	)

	if powerStateVMs == 0 {
		sectionHeader = strings.ReplaceAll(sectionHeader, ":", "")
	}

	fmt.Fprint(w, sectionHeader)

	for _, vm := range vms {
		if vm.PowerState == powerState {
			fmt.Fprintf(
				w,
				"* %s [Size: %s, Datastore Usage: %s]%s",
				vm.Name,
				vm.VMSize,
				vm.DatastoreUsage,
				nagios.CheckOutputEOL,
			)
		}
	}

	fmt.Fprintf(w, nagios.CheckOutputEOL)

}

// ValidateDatastoreAccessibility evaluates a given Datastore's accessibility
// and returns a list of reasons why and an error if the datastore is
// inaccessible. If the Datastore is accessible, nil is returned for both
// values.
func ValidateDatastoreAccessibility(ds mo.Datastore) ([]string, error) {

	reasons := []string{"unknown"}
	if !ds.Summary.Accessible {
		for _, hostMount := range ds.Host {
			if hostMount.MountInfo.Accessible != nil &&
				!*hostMount.MountInfo.Accessible {

				reasons = append(reasons, hostMount.MountInfo.InaccessibleReason)
			}
		}

		logger.Printf(
			"datastore %q is inaccessible due to: [%v]\n",
			ds.Name,
			strings.Join(reasons, ", "),
		)

		return reasons, ErrDatastoreInaccessible

	}

	return nil, nil

}

// ValidateDatastoreStatsCollectionStatus returns nil indicating that
// statistics collection for a Datastore is enabled or an error which provides
// more information.
func ValidateDatastoreStatsCollectionStatus(datastore mo.Datastore) error {

	switch {

	// This field is required to determine stats collection status.
	case datastore.IormConfiguration == nil:

		return ErrDatastoreIormConfigurationPropertyUnavailable

	case (datastore.IormConfiguration).StatsCollectionEnabled == nil:

		return ErrDatastoreStatsCollectionPropertyUnavailable

	case !*(datastore.IormConfiguration).StatsCollectionEnabled:

		return ErrDatastoreIormConfigurationStatisticsCollectionDisabled

	default:
		return nil
	}

}

// IsWarningState indicates whether a Datastore Performance Summary metric
// has crossed the WARNING level threshold.
func (dps DatastorePerformanceSummary) IsWarningState() bool {

	switch {

	// Only nil if the thresholds have not been defined, which indicates that
	// the percentile associated with these metrics was not requested.
	case dps.thresholds == nil:
		return false

	case dps.ReadLatency > dps.thresholds.ReadLatencyWarning &&
		dps.ReadLatency < dps.thresholds.ReadLatencyCritical:
		return true

	case dps.WriteLatency > dps.thresholds.WriteLatencyWarning &&
		dps.WriteLatency < dps.thresholds.WriteLatencyCritical:
		return true

	case dps.VMLatency > dps.thresholds.VMLatencyWarning &&
		dps.VMLatency < dps.thresholds.VMLatencyCritical:
		return true

	default:
		return false
	}

}

// IsCriticalState indicates whether a Datastore Performance Summary metric
// has crossed the CRITICAL level threshold.
func (dps DatastorePerformanceSummary) IsCriticalState() bool {

	switch {

	// Only nil if the thresholds have not been defined, which indicates that
	// the percentile associated with these metrics was not requested.
	case dps.thresholds == nil:
		return false

	case dps.ReadLatency > dps.thresholds.ReadLatencyCritical:
		return true

	case dps.WriteLatency > dps.thresholds.WriteLatencyCritical:
		return true

	case dps.VMLatency > dps.thresholds.VMLatencyCritical:
		return true

	default:
		return false
	}

}

// IsZero indicates whether Datastore Performance Summary metrics are all
// value 0. This is a common occurrence after a new interval begins. For
// approximately 30 minutes no metrics are available until (presumably)
// sufficient time has elapsed to reliably generate aggregates of performance
// data. This can also occur if performance metrics collection is disabled for
// a Datastore.
func (dps DatastorePerformanceSummary) IsZero() bool {

	if dps.ReadLatency > 0 ||
		dps.WriteLatency > 0 ||
		dps.VMLatency > 0 ||
		dps.ReadIops > 0 ||
		dps.WriteIops > 0 {

		logger.Printf(
			"DatastorePerformanceSummary: one or more metrics for percentile %d are set",
			dps.Percentile,
		)

		return false
	}

	logger.Printf(
		"DatastorePerformanceSummary: all metrics for percentile %d are value 0",
		dps.Percentile,
	)

	return true
}

// MetricsAboveThreshold returns a list of Datastore performance metrics which
// have exceeded specified thresholds.
func (dps DatastorePerformanceSummary) MetricsAboveThreshold() []string {

	// Read Latency, Write Latency, VM Latency
	//
	// TODO: Extend this if we evaluate IOPS values in the future. See
	// constants.go for commented constants.
	totalMetrics := 3
	exceeded := make([]string, 0, totalMetrics)

	// Only nil if the thresholds have not been defined, which indicates that
	// the percentile associated with these metrics was not requested. Nothing
	// else to evaluate, skip any further checks.
	if dps.thresholds == nil {

		return exceeded
	}

	if dps.ReadLatency > dps.thresholds.ReadLatencyCritical ||
		dps.ReadLatency > dps.thresholds.ReadLatencyWarning {
		exceeded = append(exceeded, readLatency)
	}

	if dps.WriteLatency > dps.thresholds.WriteLatencyCritical ||
		dps.WriteLatency > dps.thresholds.WriteLatencyWarning {
		exceeded = append(exceeded, writeLatency)
	}

	if dps.VMLatency > dps.thresholds.VMLatencyCritical ||
		dps.VMLatency > dps.thresholds.VMLatencyWarning {
		exceeded = append(exceeded, vmLatency)
	}

	return exceeded

}

// DatastoreVMsSummary evaluates provided Datastore and collection of
// VirtualMachines and provides a basic human readable / formatted summary of
// VirtualMachine details.
func DatastoreVMsSummary(ds mo.Datastore, vms []mo.VirtualMachine) DatastoreVMs {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute DatastoreVMsSummary func.\n",
			time.Since(funcTimeStart),
		)
	}()

	datastoreVMs := make(DatastoreVMs, 0, len(vms))

	for _, vm := range vms {

		var vmStorageUsed int64
		for _, usage := range vm.Storage.PerDatastoreUsage {
			if usage.Datastore == ds.Reference() {
				vmStorageUsed += usage.Committed + usage.Uncommitted
			}
		}

		vmPercentOfDSUsed := float64(vmStorageUsed) / float64(ds.Summary.Capacity) * 100
		dsVM := DatastoreVM{
			Name:           vm.Name,
			VMSize:         units.ByteSize(vmStorageUsed).String(),
			DatastoreUsage: fmt.Sprintf("%2.2f%%", vmPercentOfDSUsed),
			PowerState:     vm.Runtime.PowerState,
		}

		datastoreVMs = append(datastoreVMs, dsVM)

	}

	return datastoreVMs

}

// NewDatastorePerformanceSet receives a Datastore and a specified thresholds
// index and retrieves performance summary information used to determine if
// storage latency levels have crossed user-specified thresholds.
func NewDatastorePerformanceSet(
	ctx context.Context,
	c *vim25.Client,
	ds mo.Datastore,
	thresholdsIndex DatastorePerformanceThresholdsIndex,
) (DatastorePerformanceSet, error) {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute NewDatastorePerformanceSet func.\n",
			time.Since(funcTimeStart),
		)
	}()

	// Determine whether stats collection is enabled.
	statsCollectionStatusErr := ValidateDatastoreStatsCollectionStatus(ds)
	if statsCollectionStatusErr != nil {
		return DatastorePerformanceSet{}, statsCollectionStatusErr
	}

	srm := object.NewStorageResourceManager(c)
	dsObj := object.NewDatastore(c, ds.Reference())
	results, err := srm.QueryDatastorePerformanceSummary(ctx, dsObj)
	if err != nil {
		errMsg := fmt.Sprintf(
			"error retrieving performance summary details for datastore %s",
			ds.Name,
		)

		logger.Print(errMsg)

		return DatastorePerformanceSet{}, errors.New(errMsg)
	}

	// Return a sentinel error for empty or "unavailable" metrics collection
	// so that caller can optionally handle this specific scenario differently
	// based on whether the specified datastore is newly created.
	if len(results) == 0 {
		errMsg := fmt.Sprintf(
			"missing performance summary details for datastore %s",
			ds.Name,
		)

		logger.Print(errMsg)

		return DatastorePerformanceSet{}, ErrDatastorePerformanceMetricsMissing
	}

	dsVMs, err := GetVMsFromDatastore(ctx, c, ds, true)
	if err != nil {
		return DatastorePerformanceSet{}, fmt.Errorf(
			"error retrieving VMs for Datastore Performance Set: %w", err,
		)
	}

	perfSummaryIndexes := make([]DatastorePerformanceSummaryIndex, 0, len(results))

	for resultsIdx, result := range results {

		// A separate index is recorded for each distinct interval or window
		// of time that metrics are recorded.
		perfSummaryEntries := make(map[int]DatastorePerformanceSummary)

		// Metrics are aggregated based on percentile. We use the position of
		// the percentile value in its index to index into the metrics which
		// correspond to that percentile.
		for pIdx, percentile := range result.Percentile {

			// if specified by the user, record thresholds
			var thresholds *DatastorePerformanceThresholds
			if val, ok := thresholdsIndex[int(percentile)]; ok {
				thresholds = &val
			}

			perfSummaryEntries[int(percentile)] = DatastorePerformanceSummary{
				DatastoreMOID: ds.Self,
				ReadIops:      result.DatastoreReadIops[pIdx],
				ReadLatency:   result.DatastoreReadLatency[pIdx],
				WriteIops:     result.DatastoreWriteIops[pIdx],
				WriteLatency:  result.DatastoreWriteLatency[pIdx],
				VMLatency:     result.DatastoreVmLatency[pIdx],
				Percentile:    percentile,
				Interval:      result.Interval,
				thresholds:    thresholds,
			}
		}

		// Assert that all requested percentiles are represented in the
		// metrics returned by the QueryDatastorePerformanceSummary() method
		// call.
		inPercentilesList := func(percentile int, percentiles []int) bool {
			for i := range percentiles {
				if percentile == percentiles[i] {
					return true
				}
			}

			return false
		}
		percentilesAvailable := make([]int, 0, len(perfSummaryEntries))
		for key := range perfSummaryEntries {
			percentilesAvailable = append(percentilesAvailable, key)
		}
		for percentile := range thresholdsIndex {
			if !inPercentilesList(percentile, percentilesAvailable) {
				return DatastorePerformanceSet{}, ErrDatastorePerformancePercentileUnavailable
			}
		}

		perfSummaryIndexes = append(perfSummaryIndexes, DatastorePerformanceSummaryIndex{
			Entries: perfSummaryEntries,

			// Based on observation, the first element from the
			// `QueryDatastorePerformanceSummary()` query is for the active
			// interval or window of time. Other entries are for prior
			// windows, and appear to be placed at the end of the slice as
			// those windows of time close (based on element 0 holding active
			// details and element 7 holding the prior interval's metrics).
			Active: resultsIdx == 0,
		})
	}

	dsPerfSet := DatastorePerformanceSet{
		Datastore: ds,
		VMs:       DatastoreVMsSummary(ds, dsVMs),
		Intervals: perfSummaryIndexes,
	}

	return dsPerfSet, nil

}

// Percentiles returns a sorted list of all Datastore Performance Summary
// percentiles in the index.
func (dpsi DatastorePerformanceSummaryIndex) Percentiles() []int {

	// use sorted keys for consistency in percentile list order
	percentiles := make([]int, 0, len(dpsi.Entries))
	for percentile := range dpsi.Entries {
		percentiles = append(percentiles, percentile)
	}
	sort.Slice(percentiles, func(i, j int) bool {
		return percentiles[i] < percentiles[j]
	})

	return percentiles

}

// IsWarningState indicates whether a Datastore Performance Summary metric in
// the index has crossed the WARNING level threshold.
func (dpsi DatastorePerformanceSummaryIndex) IsWarningState() bool {

	for percentile, perfSummary := range dpsi.Entries {

		switch {
		case perfSummary.thresholds != nil:
			logger.Printf("thresholds provided for percentile %d", percentile)
			logger.Printf("evaluating performance metrics for percentile %d", percentile)

			if perfSummary.IsWarningState() {
				return true
			}

		default:
			logger.Printf("thresholds not provided for percentile %d", percentile)
		}
	}

	return false

}

// IsCriticalState indicates whether a Datastore Performance Summary metric in
// the index has crossed the WARNING level threshold.
func (dpsi DatastorePerformanceSummaryIndex) IsCriticalState() bool {

	for percentile, perfSummary := range dpsi.Entries {

		switch {
		case perfSummary.thresholds != nil:
			logger.Printf("thresholds provided for percentile %d", percentile)
			logger.Printf("evaluating performance metrics for percentile %d", percentile)

			if perfSummary.IsCriticalState() {
				return true
			}

		default:
			logger.Printf("thresholds not provided for percentile %d", percentile)
		}
	}

	return false

}

// IsZero indicates whether all Datastore Performance Summary metrics for all
// percentiles in the index are value 0.
func (dpsi DatastorePerformanceSummaryIndex) IsZero() bool {

	for percentile, perfSummary := range dpsi.Entries {
		if !perfSummary.IsZero() {
			logger.Printf("DatastorePerformanceSummaryIndex: one or more metrics for percentile %d are set", percentile)
			return false
		}
	}

	logger.Print("DatastorePerformanceSummaryIndex: all metrics for all percentiles are value 0")

	return true

}

// IsWarningState indicates whether a Datastore Performance Summary metric in
// the set has crossed the WARNING level threshold.
func (dps DatastorePerformanceSet) IsWarningState() bool {

	// The first element contains metrics for the active interval which are
	// being actively aggregated. The other intervals are historical values
	// useful for review, but not for determining current plugin state.
	//
	// TODO: Audit potential for this to cause an index out of range panic
	activeIndex := dps.Intervals[0]

	return activeIndex.IsWarningState()

}

// IsCriticalState indicates whether a Datastore Performance Summary metric in
// the set has crossed the CRITICAL level threshold.
func (dps DatastorePerformanceSet) IsCriticalState() bool {

	// The first element contains metrics for the active interval which are
	// being actively aggregated. The other intervals are historical values
	// useful for review, but not for determining current plugin state.
	//
	// TODO: Audit potential for this to cause an index out of range panic
	activeIndex := dps.Intervals[0]

	return activeIndex.IsCriticalState()

}

// IsUnknownState indicates whether a DatastorePerformanceSet is in an UNKNOWN
// state.
func (dps DatastorePerformanceSet) IsUnknownState() bool {
	return dps.UnknownState() != nil
}

// UnknownState provides the associated error for a DatastorePerformanceSet's
// UNKNOWN state.
func (dps DatastorePerformanceSet) UnknownState() error {

	statsCollectionStatus := ValidateDatastoreStatsCollectionStatus(dps.Datastore)

	switch {

	// One "unknowable" state is when the Statistics Collection setting for a
	// datastore is definitively disabled. Without collection of statistics,
	// this plugin cannot make a determination whether datastore performance
	// is within specified bounds.
	//
	// When disabled, the
	// StorageResourceManager.QueryDatastorePerformanceSummary() method
	// returns metric sets all of value 0 (specifically 0.000000).
	case errors.Is(statsCollectionStatus, ErrDatastoreIormConfigurationStatisticsCollectionDisabled):
		return statsCollectionStatus

	// If the Statistics Collection setting for a datastore is enabled and we
	// have running VMs on the datastore, we should have Datastore performance
	// statistics to evaluate. If we instead only have metrics of value 0,
	// then we have another "unknowable" state.
	case dps.IsZero():
		if dps.VMs.NumVMsPoweredOn() > 0 {
			return ErrDatastoreLatencyAllMetricSetsZero
		}

		// If there are no running VMs, not having metrics for the datastore
		// isn't considered a problem state.
		return nil

	default:
		return nil
	}

}

// IsZero indicates whether all Datastore Performance Summary metrics for all
// percentiles in all indexes in the set are value 0.
func (dps DatastorePerformanceSet) IsZero() bool {

	for interval := range dps.Intervals {
		if !dps.Intervals[interval].IsZero() {
			logger.Printf("DatastorePerformanceSet: one or more metrics for interval %d are set", interval)
			return false
		}
	}

	logger.Print("DatastorePerformanceSet: all metrics for all percentiles in all intervals of the set are value 0")

	return true

}

// ActivePerfSummaryIndex returns the active DatastorePerformanceSummaryIndex
// from the set or an error if it is not available.
func (dps DatastorePerformanceSet) ActivePerfSummaryIndex() (DatastorePerformanceSummaryIndex, error) {

	for i := range dps.Intervals {
		if dps.Intervals[i].Active {
			return dps.Intervals[i], nil
		}
	}

	return DatastorePerformanceSummaryIndex{}, fmt.Errorf(
		"unknown error; active performance summary index not found",
	)

}

// ActiveIntervalMetrics returns the associated DatastorePerformanceSummary
// for a specified percentile from the active interval. An error is returned
// if an invalid percentile is specified.
func (dps DatastorePerformanceSet) ActiveIntervalMetrics(percentile int) (DatastorePerformanceSummary, error) {

	var perSummaryIdx DatastorePerformanceSummaryIndex
	for i := range dps.Intervals {
		if dps.Intervals[i].Active {
			logger.Print("Found active interval:", i)
			perSummaryIdx = dps.Intervals[i]
		}
	}

	var summary DatastorePerformanceSummary
	var ok bool
	summary, ok = perSummaryIdx.Entries[percentile]
	if !ok {
		return DatastorePerformanceSummary{}, fmt.Errorf(
			"invalid percentile specified: %d",
			percentile,
		)
	}

	return summary, nil

}

// NewDatastoreUsageSummary receives a Datastore and generates summary
// information used to determine if usage levels have crossed user-specified
// thresholds.
// func NewDatastoreUsageSummary(ds mo.Datastore, dsVMs []mo.VirtualMachine, criticalThreshold int, warningThreshold int) DatastoreUsageSummary {
func NewDatastoreUsageSummary(
	ctx context.Context,
	c *vim25.Client,
	ds mo.Datastore,
	criticalThreshold int,
	warningThreshold int,
) (DatastoreUsageSummary, error) {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute NewDatastoreUsageSummary func.\n",
			time.Since(funcTimeStart),
		)
	}()

	storageRemainingPercentage := float64(ds.Summary.FreeSpace) / float64(ds.Summary.Capacity) * 100
	storageUsedPercentage := 100 - storageRemainingPercentage
	storageRemaining := ds.Summary.FreeSpace
	storageTotal := ds.Summary.Capacity
	storageUsed := storageTotal - storageRemaining

	dsVMs, err := GetVMsFromDatastore(ctx, c, ds, true)
	if err != nil {
		return DatastoreUsageSummary{}, err
	}

	dsUsage := DatastoreUsageSummary{
		Datastore:               ds,
		VMs:                     DatastoreVMsSummary(ds, dsVMs),
		StorageRemainingPercent: storageRemainingPercentage,
		StorageUsedPercent:      storageUsedPercentage,
		StorageTotal:            storageTotal,
		StorageUsed:             storageUsed,
		StorageRemaining:        storageRemaining,
		CriticalThreshold:       criticalThreshold,
		WarningThreshold:        warningThreshold,
	}

	return dsUsage, nil

}

// IsWarningState indicates whether Datastore usage has crossed the WARNING
// level threshold.
func (dus DatastoreUsageSummary) IsWarningState() bool {
	return dus.StorageUsedPercent < float64(dus.CriticalThreshold) &&
		dus.StorageUsedPercent > float64(dus.WarningThreshold)
}

// IsCriticalState indicates whether Datastore usage has crossed the CRITICAL
// level threshold.
func (dus DatastoreUsageSummary) IsCriticalState() bool {
	return dus.StorageUsedPercent > float64(dus.CriticalThreshold)
}

// NumVMsPoweredOn indicates how many VirtualMachines on a specific Datastore
// are powered on.
func (dsVMs DatastoreVMs) NumVMsPoweredOn() int {

	var numOn int
	for _, vm := range dsVMs {
		if vm.PowerState == types.VirtualMachinePowerStatePoweredOn {
			numOn++
		}
	}

	return numOn
}

// NumVMsPoweredOff indicates how many VirtualMachines on a specific Datastore
// are powered off OR suspended.
func (dsVMs DatastoreVMs) NumVMsPoweredOff() int {
	return len(dsVMs) - dsVMs.NumVMsPoweredOn()
}

// GetDatastores accepts a context, a connected client and a boolean value
// indicating whether a subset of properties per Datastore are retrieved. A
// collection of Datastores with requested properties is returned. If
// requested, a subset of all available properties will be retrieved (faster)
// instead of recursively fetching all properties (about 2x as slow).
func GetDatastores(ctx context.Context, c *vim25.Client, propsSubset bool) ([]mo.Datastore, error) {

	funcTimeStart := time.Now()

	// declare this early so that we can grab a pointer to it in order to
	// access the entries later
	var dss []mo.Datastore

	defer func(dss *[]mo.Datastore) {
		logger.Printf(
			"It took %v to execute GetDatastores func (and retrieve %d Datastores).\n",
			time.Since(funcTimeStart),
			len(*dss),
		)
	}(&dss)

	err := getObjects(ctx, c, &dss, c.ServiceContent.RootFolder, propsSubset)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve Datastores: %w", err)
	}

	sort.Slice(dss, func(i, j int) bool {
		return strings.ToLower(dss[i].Name) < strings.ToLower(dss[j].Name)
	})

	return dss, nil
}

// GetDatastoreByName accepts the name of a datastore, the name of a
// datacenter and a boolean value indicating whether only a subset of
// properties for the Datastore should be returned. If requested, a subset of
// all available properties will be retrieved (faster) instead of recursively
// fetching all properties (about 2x as slow). If the datacenter name is an
// empty string then the default datacenter will be used.
func GetDatastoreByName(ctx context.Context, c *vim25.Client, dsName string, datacenter string, propsSubset bool) (mo.Datastore, error) {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute GetDatastoreByName func.\n",
			time.Since(funcTimeStart),
		)
	}()

	var datastore mo.Datastore

	err := getObjectByName(ctx, c, &datastore, dsName, datacenter, propsSubset)
	if err != nil {
		return mo.Datastore{}, err
	}

	return datastore, nil

}

// FilterDatastoresByName accepts a collection of Datastores and a Datastore
// name to filter against. An error is returned if the list of Datastores is
// empty or if a match was not found. The matching Datastore is returned along
// with the number of Datastores that were excluded.
func FilterDatastoresByName(dss []mo.Datastore, dsName string) (mo.Datastore, int, error) {

	funcTimeStart := time.Now()

	// If error condition, no exclusions are made
	numExcluded := 0

	defer func() {
		logger.Printf(
			"It took %v to execute FilterDatastoresByName func.\n",
			time.Since(funcTimeStart),
		)
	}()

	if len(dss) == 0 {
		return mo.Datastore{}, numExcluded, fmt.Errorf("received empty list of datastores to filter by name")
	}

	for _, ds := range dss {
		if ds.Name == dsName {
			// we are excluding everything but the single name value match
			numExcluded = len(dss) - 1
			return ds, numExcluded, nil
		}
	}

	return mo.Datastore{}, numExcluded, fmt.Errorf(
		"error: failed to retrieve Datastore using provided name %q",
		dsName,
	)

}

// FilterDatastoresByID receives a collection of Datastores and a Datastore ID
// to filter against. An error is returned if the list of Datastores is empty
// or if a match was not found. The matching Datastore is returned along with
// the number of Datastores that were excluded.
func FilterDatastoresByID(dss []mo.Datastore, dsID string) (mo.Datastore, int, error) {

	funcTimeStart := time.Now()

	// If error condition, no exclusions are made
	numExcluded := 0

	defer func() {
		logger.Printf(
			"It took %v to execute FilterDatastoresByID func.\n",
			time.Since(funcTimeStart),
		)
	}()

	if len(dss) == 0 {
		return mo.Datastore{}, numExcluded, fmt.Errorf("received empty list of datastores to filter by ID")
	}

	for _, ds := range dss {
		// return match, if available
		// TODO: Refactor, use abstract type here
		// ds.GetManagedEntity().Reference().Value
		if ds.Summary.Datastore.Value == dsID {
			// we are excluding everything but the single name value match
			numExcluded = len(dss) - 1
			return ds, numExcluded, nil
		}
	}

	return mo.Datastore{}, numExcluded, fmt.Errorf(
		"error: failed to retrieve Datastore using provided id %q",
		dsID,
	)

}

// DatastoreIDsToNames returns a list of matching Datastore names for the
// provided list of Managed Object References for Datastores.
func DatastoreIDsToNames(dsRefs []types.ManagedObjectReference, dss []mo.Datastore) []string {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute DatastoreIDsToNames func.\n",
			time.Since(funcTimeStart),
		)
	}()

	dsNames := make([]string, 0, len(dsRefs))
	dsIDs := make([]string, 0, len(dsRefs))

	for _, dsRef := range dsRefs {
		dsIDs = append(dsIDs, dsRef.Value)
	}

	for _, ds := range dss {
		if textutils.InList(ds.Summary.Datastore.Value, dsIDs, true) {
			dsNames = append(dsNames, ds.Name)
		}
	}

	return dsNames

}

// DatastoreUsageOneLineCheckSummary is used to generate a one-line Nagios
// service check results summary. This is the line most prominent in
// notifications.
func DatastoreUsageOneLineCheckSummary(
	stateLabel string,
	dsUsageSummary DatastoreUsageSummary,
) string {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute DatastoreUsageOneLineCheckSummary func.\n",
			time.Since(funcTimeStart),
		)
	}()

	return fmt.Sprintf(
		"%s: Datastore %s usage (%d VMs) is %.2f%% of %s with %s remaining [WARNING: %d%% , CRITICAL: %d%%]",
		stateLabel,
		dsUsageSummary.Datastore.Name,
		len(dsUsageSummary.VMs),
		dsUsageSummary.StorageUsedPercent,
		units.ByteSize(dsUsageSummary.StorageTotal),
		units.ByteSize(dsUsageSummary.StorageRemaining),
		dsUsageSummary.WarningThreshold,
		dsUsageSummary.CriticalThreshold,
	)

}

// DatastoreUsageReport generates a summary of Datastore usage along with
// various verbose details intended to aid in troubleshooting check results at
// a glance. This information is provided for use with the Long Service Output
// field commonly displayed on the detailed service check results display in
// the web UI or in the body of many notifications.
func DatastoreUsageReport(
	c *vim25.Client,
	dsUsageSummary DatastoreUsageSummary,
) string {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute DatastoreUsageReport func.\n",
			time.Since(funcTimeStart),
		)
	}()

	var report strings.Builder

	fmt.Fprintf(
		&report,
		"Datastore Summary:%s%s"+
			"* Name: %s%s"+
			"* Used: %v (%.2f%%)%s"+
			"* Remaining: %v (%.2f%%)%s"+
			"* VMs: %v %s%s",
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
		dsUsageSummary.Datastore.Name,
		nagios.CheckOutputEOL,
		units.ByteSize(dsUsageSummary.StorageUsed),
		dsUsageSummary.StorageUsedPercent,
		nagios.CheckOutputEOL,
		units.ByteSize(dsUsageSummary.StorageRemaining),
		dsUsageSummary.StorageRemainingPercent,
		nagios.CheckOutputEOL,
		len(dsUsageSummary.VMs),
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
	)

	printVMSummary(&report, dsUsageSummary.VMs, types.VirtualMachinePowerStatePoweredOn)

	printVMSummary(&report, dsUsageSummary.VMs, types.VirtualMachinePowerStatePoweredOff)

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
		"* Plugin User Agent: %s%s",
		c.Client.UserAgent,
		nagios.CheckOutputEOL,
	)

	return report.String()
}

// DatastorePerformanceOneLineCheckSummary is used to generate a one-line Nagios
// service check results summary. This is the line most prominent in
// notifications.
func DatastorePerformanceOneLineCheckSummary(
	stateLabel string,
	dsPerfSet DatastorePerformanceSet,
) string {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute DatastorePerformanceOneLineCheckSummary func.\n",
			time.Since(funcTimeStart),
		)
	}()

	switch {

	case dsPerfSet.IsUnknownState():

		return fmt.Sprintf(
			"%s: Datastore %s (%d VMs) performance metrics are unavailable",
			stateLabel,
			dsPerfSet.Datastore.Name,
			len(dsPerfSet.VMs),
		)

	case dsPerfSet.IsWarningState() || dsPerfSet.IsCriticalState():

		var metricsExceededThresholds []string

		for _, perSummaryIndex := range dsPerfSet.Intervals {

			if !perSummaryIndex.Active {
				continue
			}

			percentiles := perSummaryIndex.Percentiles()
			for i := range percentiles {
				percentile := percentiles[i]
				summary := perSummaryIndex.Entries[percentile]

				if summary.IsCriticalState() || summary.IsWarningState() {
					metricsExceededThresholds = append(metricsExceededThresholds, summary.MetricsAboveThreshold()...)
				}
			}

		}

		metricsExceededThresholds = textutils.DedupeList(metricsExceededThresholds)
		sort.Strings(metricsExceededThresholds)

		return fmt.Sprintf(
			"%s: Datastore %s (%d VMs) exceeds specified performance thresholds: [%v]",
			stateLabel,
			dsPerfSet.Datastore.Name,
			len(dsPerfSet.VMs),
			strings.Join(metricsExceededThresholds, ", "),
		)

	default:

		reason := func() string {

			var reason string
			switch {
			case dsPerfSet.VMs.NumVMsPoweredOn() == 0:
				reason = " (no running VMs)"
			default:
				reason = " (thresholds not exceeded)"
			}

			return reason
		}

		return fmt.Sprintf(
			"%s: Datastore %s (%d VMs) meets specified performance thresholds%s",
			stateLabel,
			dsPerfSet.Datastore.Name,
			len(dsPerfSet.VMs),
			reason(),
		)

	}

}

// DatastorePerformanceReport generates a summary of Datastore usage along
// with various verbose details intended to aid in troubleshooting check
// results at a glance. This information is provided for use with the Long
// Service Output field commonly displayed on the detailed service check
// results display in the web UI or in the body of many notifications.
func DatastorePerformanceReport(
	c *vim25.Client,
	dsPerfSet DatastorePerformanceSet,
	hideHistoricalMetricSets bool,
) string {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute DatastorePerformanceReport func.\n",
			time.Since(funcTimeStart),
		)
	}()

	var report strings.Builder

	// TODO: Is there a useful header we can include here?
	//
	// fmt.Fprintf(
	// 	&report,
	// 	"Performance Summary for datastore %q (%d VMs):%s%s",
	// 	dsPerfSet.Datastore.Name,
	// 	len(dsPerfSet.VMs),
	// 	nagios.CheckOutputEOL,
	// 	nagios.CheckOutputEOL,
	// )

	// List metrics which exceed threshold.
	if dsPerfSet.IsWarningState() || dsPerfSet.IsCriticalState() {

		fmt.Fprintf(
			&report,
			"Metrics for datastore %q which exceed thresholds:%s",
			dsPerfSet.Datastore.Name,
			nagios.CheckOutputEOL,
		)

		for result, perSummaryIndex := range dsPerfSet.Intervals {

			if !perSummaryIndex.Active {
				continue
			}

			fmt.Fprintf(
				&report,
				"%sResult %v (active): %s",
				nagios.CheckOutputEOL,
				result+1,
				nagios.CheckOutputEOL,
			)

			// use sorted keys for consistency in percentile list order
			percentiles := perSummaryIndex.Percentiles()
			for i := range percentiles {

				percentile := percentiles[i]
				summary := perSummaryIndex.Entries[percentile]

				// Skip emitting any metrics which don't exceed the thresholds.
				if summary.IsCriticalState() || summary.IsWarningState() {
					fmt.Fprintf(
						&report,
						"  * { Percentile: %d, RLatency: %.2f, WLatency: %.2f, VMLatency: %.2f, RIops: %.2f, WIops: %.2f, Interval: %d }%s",
						percentile,
						summary.ReadLatency,
						summary.WriteLatency,
						summary.VMLatency,
						summary.ReadIops,
						summary.WriteIops,
						summary.Interval,
						nagios.CheckOutputEOL,
					)
				}
			}
		}

		fmt.Fprintf(&report, nagios.CheckOutputEOL)

	}

	var metricCollectionsHeaderTemplate string
	switch {
	case hideHistoricalMetricSets:
		metricCollectionsHeaderTemplate =
			"Active Performance Metrics for datastore %q (%d VMs):%s"

	default:
		metricCollectionsHeaderTemplate =
			"Full Collection of Performance Metrics for datastore %q (%d VMs):%s"
	}

	fmt.Fprintf(
		&report,
		metricCollectionsHeaderTemplate,
		dsPerfSet.Datastore.Name,
		len(dsPerfSet.VMs),
		nagios.CheckOutputEOL,
	)

	for result, perSummaryIndex := range dsPerfSet.Intervals {

		// Skip emission of historical metric sets if requested.
		if hideHistoricalMetricSets && !perSummaryIndex.Active {
			continue
		}

		activeIndicator := "historical"
		if perSummaryIndex.Active {
			activeIndicator = "active"
		}

		fmt.Fprintf(
			&report,
			"%sResult %v (%s): %s",
			nagios.CheckOutputEOL,
			result+1,
			activeIndicator,
			nagios.CheckOutputEOL,
		)

		// use sorted keys for consistency in percentile list order
		percentiles := perSummaryIndex.Percentiles()
		for i := range percentiles {

			percentile := percentiles[i]
			summary := perSummaryIndex.Entries[percentile]

			fmt.Fprintf(
				&report,
				// "\t* { Percentile: %d, Read Latency: %.2f, Write Latency: %.2f, VM Latency: %.2f, Read Iops: %.2f, Write Iops: %.2f, Interval: %d%s",
				"  * { Percentile: %d, RLatency: %.2f, WLatency: %.2f, VMLatency: %.2f, RIops: %.2f, WIops: %.2f, Interval: %d }%s",
				percentile,
				summary.ReadLatency,
				summary.WriteLatency,
				summary.VMLatency,
				summary.ReadIops,
				summary.WriteIops,
				summary.Interval,
				nagios.CheckOutputEOL,
			)
		}
	}

	fmt.Fprintf(&report, nagios.CheckOutputEOL)

	printVMSummary(&report, dsPerfSet.VMs, types.VirtualMachinePowerStatePoweredOn)

	printVMSummary(&report, dsPerfSet.VMs, types.VirtualMachinePowerStatePoweredOff)

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
		"* Plugin User Agent: %s%s",
		c.Client.UserAgent,
		nagios.CheckOutputEOL,
	)

	return report.String()
}
