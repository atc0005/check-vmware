package vsphere

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/atc0005/go-nagios"
	"github.com/vmware/govmomi/units"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

// ErrSnapshotAgeThresholdCrossed indicates that a snapshot is older than a
// specified age threshold
var ErrSnapshotAgeThresholdCrossed = errors.New("snapshot exceeds specified age threshold")

// ErrSnapshotSizeThresholdCrossed indicates that a snapshot is larger than a
// specified size threshold
var ErrSnapshotSizeThresholdCrossed = errors.New("snapshot exceeds specified size threshold")

// ExceedsAge indicates whether a given snapshot creation date is older than
// the specified number of days.
func ExceedsAge(snapshotCreated time.Time, days int) bool {

	now := time.Now()

	// Flip user specified number of days negative so that we can wind
	// back that many days from the file modification time. This gives
	// us our threshold to compare file modification times against.
	daysBack := -(days)
	ageThreshold := now.AddDate(0, 0, daysBack)

	switch {
	case snapshotCreated.Before(ageThreshold):
		return true
	case snapshotCreated.Equal(ageThreshold):
		return false
	case snapshotCreated.After(ageThreshold):
		return false

	// TODO: Is there any other state than Before, Equal and After?
	// TODO: Perhaps remove 'After' and use this instead?
	default:
		return false
	}

}

// FilterVMsWithSnapshots filters the provided collection of VirtualMachines
// to just those with snapshots. Later steps are responsible for validating
// whether those snapshots place the VMs into non-OK states.
func FilterVMsWithSnapshots(vms []mo.VirtualMachine) []mo.VirtualMachine {

	// setup early so we can reference it from deferred stats output
	var vmsWithSnapshots []mo.VirtualMachine

	funcTimeStart := time.Now()

	defer func(vms []mo.VirtualMachine, filteredVMs *[]mo.VirtualMachine) {
		fmt.Fprintf(
			os.Stderr,
			"It took %v to execute FilterVMsWithSnapshots func (for %d VMs, yielding %d VMs).\n",
			time.Since(funcTimeStart),
			len(vms),
			len(*filteredVMs),
		)
	}(vms, &vmsWithSnapshots)

	for _, vm := range vms {

		if vm.Snapshot != nil && vm.Snapshot.RootSnapshotList != nil {
			vmsWithSnapshots = append(vmsWithSnapshots, vm)
		}
	}

	return vmsWithSnapshots

}

// SnapshotSummary is intended to be a summary of the most commonly used
// snapshot details for a specific VirtualMachine snapshot.
type SnapshotSummary struct {

	// Name of the snapshot in human readable format.
	Name string

	// MOID is the Managed Object Reference value for the snapshot.
	MOID string

	// Description of the snapshot in human readable format.
	Description string

	// createTime is when the snapshot was created.
	createTime time.Time

	// Size is the size of the snapshot.
	Size int64

	// ID is the unique identifier that distinguishes this snapshot from other
	// snapshots of the virtual machine.
	ID int32

	// ageWarningState indicates whether this snapshot is considered in a
	// WARNING state based on crossing snapshot age threshold.
	ageWarningState bool

	// ageCriticalState indicates whether this snapshot is considered in a
	// CRITICAL state based on crossing snapshot age threshold.
	ageCriticalState bool

	// sizeWarningState indicates Whether this snapshot is considered in a
	// WARNING state based on crossing snapshot size threshold.
	sizeWarningState bool

	// sizeCriticalState indicates whether this snapshot is considered in a
	// CRITICAL state based on crossing snapshot size threshold.
	sizeCriticalState bool

	VMName string
}

// SnapshotSummarySet ties a collection of snapshot summary values to a
// specific VirtualMachine by way of a VirtualMachine Managed Object
// Reference.
type SnapshotSummarySet struct {
	VM        types.ManagedObjectReference
	Snapshots []SnapshotSummary
}

// SnapshotSummarySets is a collection of SnapshotSummarySet types for bulk
// operations. Most often this is used when determining the overall state of
// all sets in the collection.
type SnapshotSummarySets []SnapshotSummarySet

// Size returns the size of all snapshots in the set.
// TODO: See atc0005/check-vmware#4,vmware/govmomi#2243
func (sss SnapshotSummarySet) Size() int64 {
	var sum int64
	for i := range sss.Snapshots {
		sum += sss.Snapshots[i].Size
	}

	return sum
}

// SizeHR returns the human readable size of all snapshots in the set.
// TODO: See atc0005/check-vmware#4,vmware/govmomi#2243
func (sss SnapshotSummarySet) SizeHR() string {
	return units.ByteSize(sss.Size()).String()
}

// ExceedsAge indicates how many snapshots in the set are older than the
// specified number of days.
func (sss SnapshotSummarySet) ExceedsAge(days int) int {

	var numExceeded int
	for _, snap := range sss.Snapshots {
		if snap.IsAgeExceeded(days) {
			numExceeded++
		}
	}

	return numExceeded
}

// ExceedsAge indicates how many snapshots in any of the sets are older
// than the specified number of days.
func (sss SnapshotSummarySets) ExceedsAge(days int) int {

	var numExceeded int
	for _, set := range sss {
		numExceeded += set.ExceedsAge(days)
	}

	return numExceeded
}

// SizeHR returns the human readable size of the snapshot.
// TODO: See atc0005/check-vmware#4,vmware/govmomi#2243
func (ss SnapshotSummary) SizeHR() string {
	return units.ByteSize(ss.Size).String()
}

// AgeDays returns the age of a snapshot in days.
func (ss SnapshotSummary) AgeDays() float64 {

	now := time.Now()

	return now.Sub(ss.createTime).Hours() / 24

}

// Age returns the age of a snapshot in formatted days.
func (ss SnapshotSummary) Age() string {

	return fmt.Sprintf(
		"%.2f days",
		ss.AgeDays(),
	)

}

// IsAgeExceeded indicates whether the snapshot is older than the specified
// number of days.
func (ss SnapshotSummary) IsAgeExceeded(days int) bool {
	return ExceedsAge(ss.createTime, days)
}

// IsWarningState indicates whether the snapshot has exceeded age or size
// WARNING thresholds.
func (ss SnapshotSummary) IsWarningState() bool {
	return ss.ageWarningState || ss.sizeWarningState
}

// IsCriticalState indicates whether the snapshot has exceeded age or size
// CRITICAL thresholds.
func (ss SnapshotSummary) IsCriticalState() bool {
	return ss.ageCriticalState || ss.sizeCriticalState
}

// IsAgeWarningState indicates whether the snapshot has exceeded the age
// WARNING threshold.
func (ss SnapshotSummary) IsAgeWarningState() bool {
	return ss.ageWarningState
}

// IsAgeCriticalState indicates whether the snapshot has exceeded the age
// CRITICAL threshold.
func (ss SnapshotSummary) IsAgeCriticalState() bool {
	return ss.ageCriticalState
}

// IsSizeWarningState indicates whether the snapshot has exceeded the size
// WARNING threshold.
func (ss SnapshotSummary) IsSizeWarningState() bool {
	return ss.sizeWarningState
}

// IsSizeCriticalState indicates whether the snapshot has exceeded the size
// CRITICAL threshold.
func (ss SnapshotSummary) IsSizeCriticalState() bool {
	return ss.sizeCriticalState
}

// IsWarningState indicates whether the snapshot set has exceeded age or size
// WARNING thresholds.
func (sss SnapshotSummarySet) IsWarningState() bool {
	for i := range sss.Snapshots {
		if sss.Snapshots[i].IsWarningState() {
			return true
		}
	}

	return false
}

// IsCriticalState indicates whether the snapshot set has exceeded age or size
// CRITICAL thresholds.
func (sss SnapshotSummarySet) IsCriticalState() bool {
	for i := range sss.Snapshots {
		if sss.Snapshots[i].IsCriticalState() {
			return true
		}
	}

	return false
}

// IsAgeWarningState indicates whether the snapshot set has exceeded the age
// WARNING threshold.
func (sss SnapshotSummarySet) IsAgeWarningState() bool {
	for i := range sss.Snapshots {
		if sss.Snapshots[i].IsAgeWarningState() {
			return true
		}
	}

	return false
}

// IsAgeCriticalState indicates whether the snapshot set has exceeded the age
// CRITICAL threshold.
func (sss SnapshotSummarySet) IsAgeCriticalState() bool {
	for i := range sss.Snapshots {
		if sss.Snapshots[i].IsAgeCriticalState() {
			return true
		}
	}

	return false
}

// IsSizeWarningState indicates whether the snapshot set has exceeded the
// size WARNING threshold.
func (sss SnapshotSummarySet) IsSizeWarningState() bool {
	for i := range sss.Snapshots {
		if sss.Snapshots[i].IsSizeWarningState() {
			return true
		}
	}

	return false
}

// IsSizeCriticalState indicates whether the snapshot set has exceeded the
// size CRITICAL threshold.
func (sss SnapshotSummarySet) IsSizeCriticalState() bool {
	for i := range sss.Snapshots {
		if sss.Snapshots[i].IsSizeCriticalState() {
			return true
		}
	}

	return false
}

// IsWarningState indicates whether the snapshot sets have exceeded age or
// size WARNING thresholds.
func (sss SnapshotSummarySets) IsWarningState() bool {
	for i := range sss {
		if sss[i].IsWarningState() {
			return true
		}
	}

	return false
}

// IsCriticalState indicates whether the snapshot sets have exceeded age or
// size CRITICAL thresholds.
func (sss SnapshotSummarySets) IsCriticalState() bool {
	for i := range sss {
		if sss[i].IsCriticalState() {
			return true
		}
	}

	return false
}

// IsAgeWarningState indicates whether the snapshot sets have exceeded the age
// WARNING threshold.
func (sss SnapshotSummarySets) IsAgeWarningState() bool {
	for i := range sss {
		if sss[i].IsAgeWarningState() {
			return true
		}
	}

	return false
}

// IsAgeCriticalState indicates whether the snapshot sets have exceeded the
// age CRITICAL threshold.
func (sss SnapshotSummarySets) IsAgeCriticalState() bool {
	for i := range sss {
		if sss[i].IsAgeCriticalState() {
			return true
		}
	}

	return false
}

// IsSizeWarningState indicates whether the snapshot sets have exceeded the
// size WARNING threshold.
func (sss SnapshotSummarySets) IsSizeWarningState() bool {
	for i := range sss {
		if sss[i].IsSizeWarningState() {
			return true
		}
	}

	return false
}

// IsSizeCriticalState indicates whether the snapshot sets have exceeded the
// size CRITICAL threshold.
func (sss SnapshotSummarySets) IsSizeCriticalState() bool {
	for i := range sss {
		if sss[i].IsSizeCriticalState() {
			return true
		}
	}

	return false
}

// SnapshotsIndex is a mapping of Snapshot ManagedObjectReference to a tree of
// snapshot details. This type is intended to help with producing a superset
// type combining a summary of snapshot metadata with the original
// VirtualMachine object.
//
// Deprecated ?
type SnapshotsIndex map[string]types.VirtualMachineSnapshotTree

// NewSnapshotSummarySet returns a set of SnapshotSummary values for snapshots
// associated with a specified VirtualMachine.
func NewSnapshotSummarySet(
	vm mo.VirtualMachine,
	snapshotsAgeCritical int,
	snapshotsAgeWarning int,
	snapshotsSizeCritical int,
	snapshotsSizeWarning int,
) SnapshotSummarySet {

	funcTimeStart := time.Now()

	var snapshots []SnapshotSummary

	defer func(ss *[]SnapshotSummary) {
		fmt.Fprintf(
			os.Stderr,
			"It took %v to execute NewSnapshotSummarySet func "+
				"(and retrieve %d snapshot summaries).\n",
			time.Since(funcTimeStart),
			len(*ss),
		)
	}(&snapshots)

	var crawlFunc func([]types.VirtualMachineSnapshotTree)

	crawlFunc = func(snapTree []types.VirtualMachineSnapshotTree) {

		if len(snapTree) == 0 {
			return
		}

		for _, snap := range snapTree {

			fmt.Fprintf(
				os.Stderr,
				"Processing Snapshot ID %s\n",
				snap.Snapshot.Value,
			)

			snapshots = append(snapshots, SnapshotSummary{
				Name:             snap.Name,
				VMName:           vm.Name,
				ID:               snap.Id,
				MOID:             snap.Snapshot.Value,
				Description:      snap.Description,
				createTime:       snap.CreateTime,
				ageWarningState:  ExceedsAge(snap.CreateTime, snapshotsAgeWarning),
				ageCriticalState: ExceedsAge(snap.CreateTime, snapshotsAgeCritical),

				// See atc0005/check-vmware#4,vmware/govmomi#2243
				//
				// NOTE:
				// Probably cleaner to implement as a separate helper function
				// Size:        fileLayout.Size,
				// ageWarningSize: ,
				// ageCriticalSize: ,
			})

			if snap.ChildSnapshotList != nil {
				crawlFunc(snap.ChildSnapshotList)
			}

		}
	}

	crawlFunc(vm.Snapshot.RootSnapshotList)

	return SnapshotSummarySet{
		VM:        vm.Self,
		Snapshots: snapshots,
	}

}

// SnapshotsAgeOneLineCheckSummary is used to generate a one-line Nagios
// service check results summary. This is the line most prominent in
// notifications.
func SnapshotsAgeOneLineCheckSummary(
	stateLabel string,
	snapshotSets SnapshotSummarySets,
	snapshotsAgeCritical int,
	snapshotsAgeWarning int,
	evaluatedVMs []mo.VirtualMachine,
	rps []mo.ResourcePool,
) string {

	funcTimeStart := time.Now()

	defer func() {
		fmt.Fprintf(
			os.Stderr,
			"It took %v to execute SnapshotsAgeOneLineCheckSummary func.\n",
			time.Since(funcTimeStart),
		)
	}()

	switch {

	case snapshotSets.IsAgeCriticalState():

		return fmt.Sprintf(
			"%s: %d snapshots older than %d days detected (evaluated %d VMs, %d Resource Pools)",
			stateLabel,
			snapshotSets.ExceedsAge(snapshotsAgeCritical),
			snapshotsAgeCritical,
			len(evaluatedVMs),
			len(rps),
		)

	case snapshotSets.IsAgeWarningState():

		return fmt.Sprintf(
			"%s: %d snapshots older than %d days detected (evaluated %d VMs, %d Resource Pools)",
			stateLabel,
			snapshotSets.ExceedsAge(snapshotsAgeWarning),
			snapshotsAgeWarning,
			len(evaluatedVMs),
			len(rps),
		)

	default:

		return fmt.Sprintf(
			"%s: No snapshots older than %d days detected (evaluated %d VMs, %d Resource Pools)",
			stateLabel,
			snapshotsAgeWarning,
			len(evaluatedVMs),
			len(rps),
		)

	}
}

// SnapshotsAgeReport generates a summary of snapshot details along with
// various verbose details intended to aid in troubleshooting check results at
// a glance. This information is provided for use with the Long Service Output
// field commonly displayed on the detailed service check results display in
// the web UI or in the body of many notifications.
func SnapshotsAgeReport(
	c *vim25.Client,
	snapshotSummarySets SnapshotSummarySets,
	snapshotsAgeCritical int,
	snapshotsAgeWarning int,
	allVMs []mo.VirtualMachine,
	evaluatedVMs []mo.VirtualMachine,
	vmsWithIssues []mo.VirtualMachine,
	vmsToExclude []string,
	evalPoweredOffVMs bool,
	includeRPs []string,
	excludeRPs []string,
	rps []mo.ResourcePool,
) string {

	funcTimeStart := time.Now()

	defer func() {
		fmt.Fprintf(
			os.Stderr,
			"It took %v to execute SnapshotsAgeReport func.\n",
			time.Since(funcTimeStart),
		)
	}()

	rpNames := make([]string, len(rps))
	for i := range rps {
		rpNames[i] = rps[i].Name
	}

	var report strings.Builder

	fmt.Fprintf(
		&report,
		"Snapshots exceeding CRITICAL (%dd) or WARNING (%dd) age thresholds:%s%s",
		snapshotsAgeCritical,
		snapshotsAgeWarning,
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
	)

	switch {

	case snapshotSummarySets.IsAgeCriticalState(), snapshotSummarySets.IsAgeWarningState():
		for _, snapSet := range snapshotSummarySets {
			for _, snap := range snapSet.Snapshots {
				if snap.IsAgeCriticalState() || snap.IsAgeWarningState() {
					fmt.Fprintf(
						&report,
						// See atc0005/check-vmware#4,vmware/govmomi#2243
						// "* %q [Age: %v, SnapSize: %v, Combined SnapSize: %v, Name: %q, SnapID: %v]\n"
						"* %q [Age: %v, Name: %q, SnapID: %v]\n",
						snap.VMName,
						snap.Age(),
						snap.Name,
						snap.MOID,
					)
				}
			}
		}

	default:
		fmt.Fprintln(&report, "* None detected")
	}

	fmt.Fprintf(
		&report,
		"%sSnapshots *not yet* exceeding age thresholds:%s%s",
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
		nagios.CheckOutputEOL,
	)

	switch {
	case !snapshotSummarySets.IsAgeCriticalState() && !snapshotSummarySets.IsAgeWarningState():
		for _, snapSet := range snapshotSummarySets {
			for _, snap := range snapSet.Snapshots {
				if !snap.IsAgeCriticalState() && !snap.IsAgeWarningState() {
					fmt.Fprintf(
						&report,
						// See atc0005/check-vmware#4,vmware/govmomi#2243
						// "* %q [Age: %v, SnapSize: %v, Combined SnapSize: %v, Name: %q, SnapID: %v]\n"
						"* %q [Age: %v, SnapName: %q, SnapID: %v, ]\n",
						snap.VMName,
						snap.Age(),
						snap.Name,
						snap.MOID,
					)
				}
			}
		}

	default:
		fmt.Fprintln(&report, "* None detected")
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
		"* VMs (evaluated: %d, total: %d)%s",
		len(evaluatedVMs),
		len(allVMs),
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		&report,
		"* Powered off VMs evaluated: %t%s",
		// NOTE: This plugin is hard-coded to evaluate powered off and powered
		// on VMs equally. I'm not sure whether ignoring powered off VMs by
		// default makes sense for this particular plugin.
		//
		// Please submit a GitHub issue if you feel differently and expand on
		// some use cases for ignoring powered off VMs by default.
		true,
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		&report,
		"* Specified VMs to exclude (%d): [%v]%s",
		len(vmsToExclude),
		strings.Join(vmsToExclude, ", "),
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		&report,
		"* Specified Resource Pools to explicitly include (%d): [%v]%s",
		len(includeRPs),
		strings.Join(includeRPs, ", "),
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		&report,
		"* Specified Resource Pools to explicitly exclude (%d): [%v]%s",
		len(excludeRPs),
		strings.Join(excludeRPs, ", "),
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		&report,
		"* Resource Pools evaluated (%d): [%v]%s",
		len(rpNames),
		strings.Join(rpNames, ", "),
		nagios.CheckOutputEOL,
	)

	return report.String()
}
