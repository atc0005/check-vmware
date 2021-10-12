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
	"strconv"
	"strings"
	"time"

	"github.com/atc0005/go-nagios"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/property"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

// ErrVirtualHardwareOutdatedVersionsFound indicates that hardware versions
// older than the minimum have been found.
var ErrVirtualHardwareOutdatedVersionsFound = errors.New("outdated hardware versions found")

// HardwareVersionsIndex is a map of hardware version to number of VMs present
// with that hardware version. This index serves as just that, an index.
// Accessor methods are provided to obtain HardwareVersion and
// HardwareVersions types which provide most of the useful methods for working
// with hardware version entries.
type HardwareVersionsIndex map[string]int

// HardwareVersion represents the virtual hardware version of a VirtualMachine.
type HardwareVersion struct {
	// value is the original value as provided by the
	// (types.VirtualMachineConfigInfo).Version field
	value string

	// count is the number of VirtualMachines with this hardware version
	count int

	// highest indicates whether this version is the greatest version found.
	// The default value properly indicates that a HardwareVersion is not the
	// greatest version found while applicable accessor methods expose this
	// value. Accessor methods on a HardwareVersionsIndex handle setting this
	// appropriately when constructing a collection of HardwareVersion or
	// explicitly returning the highest version.
	highest bool
}

// HardwareVersions represents a collection of HardwareVersion.
type HardwareVersions []HardwareVersion

// NewHardwareVersion creates a new HardwareVersion value using a provided
// string with "vmx-" prefix (e.g., vmx-15).
func NewHardwareVersion(verStr string) HardwareVersion {
	return HardwareVersion{
		value: verStr,
	}
}

// DefaultHardwareVersion accepts optional host, cluster and datacenter names
// and returns the default hardware version. If not specified, an attempt will
// be made to use the default Datacenter and default ComputeResource (obtained
// using cluster name). If a host name is supplied, it will be used to obtain
// the default hardware version. If a host name and a cluster name are
// provided, an error will be returned.
//
// The default version may not be the very latest version supported in the
// cluster (e.g., v14 is the default, but v15 is the latest supported).
func DefaultHardwareVersion(ctx context.Context, c *vim25.Client, hostName string, clusterName string, datacenterName string) (HardwareVersion, error) {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute DefaultHardwareVersionfunc.\n",
			time.Since(funcTimeStart),
		)
	}()

	if hostName != "" && clusterName != "" {
		return HardwareVersion{}, fmt.Errorf(
			"func DefaultHardwareVersion: only one of cluster or host name supported",
		)
	}

	finder := find.NewFinder(c, true)

	switch {
	case datacenterName == "":
		dc, findDCErr := finder.DefaultDatacenter(ctx)
		if findDCErr != nil {
			return HardwareVersion{},
				fmt.Errorf("%s: %w", dcNotProvidedFailedToFallback, findDCErr)
		}
		finder.SetDatacenter(dc)

	default:
		dc, findDCErr := finder.DatacenterOrDefault(ctx, datacenterName)
		if findDCErr != nil {
			return HardwareVersion{},
				fmt.Errorf("%s: %w", dcFailedToUseFailedToFallback, findDCErr)
		}
		finder.SetDatacenter(dc)
	}

	var computeResourceRef types.ManagedObjectReference
	switch {
	case clusterName == "":
		cr, findCRErr := finder.DefaultComputeResource(ctx)
		if findCRErr != nil {
			return HardwareVersion{},
				fmt.Errorf("%s: %w", crNotProvidedFailedToFallback, findCRErr)
		}
		computeResourceRef = cr.Reference()

	default:
		cr, findCRErr := finder.ComputeResourceOrDefault(ctx, clusterName)
		if findCRErr != nil {
			return HardwareVersion{},
				fmt.Errorf("%s: %w", crFailedToUseFailedToFallback, findCRErr)
		}
		computeResourceRef = cr.Reference()
	}

	if hostName != "" {
		hostSystem, err := GetHostSystemByName(
			ctx, c, hostName, datacenterName, true,
		)
		if err != nil {
			return HardwareVersion{}, fmt.Errorf(
				"failed to obtain default hardware version for host %s: %w",
				hostName,
				err,
			)
		}

		computeResourceRef = *hostSystem.Parent

	}

	var content []types.ObjectContent

	envBrowserErr := property.DefaultCollector(c).RetrieveOne(
		ctx,
		computeResourceRef,
		[]string{
			"environmentBrowser",
		},
		&content,
	)
	if envBrowserErr != nil {
		return HardwareVersion{}, fmt.Errorf(
			"%s: %w",
			"error creating environment browser",
			envBrowserErr,
		)
	}

	req := types.QueryConfigOptionEx{
		This: content[0].PropSet[0].Val.(types.ManagedObjectReference),
	}

	if req.Spec == nil {
		req.Spec = new(types.EnvironmentBrowserConfigOptionQuerySpec)
	}

	opt, optErr := methods.QueryConfigOptionEx(ctx, c, &req)
	if optErr != nil {
		return HardwareVersion{}, fmt.Errorf(
			"%s: %w",
			"error creating option",
			optErr,
		)

	}

	return NewHardwareVersion(opt.Returnval.Version), nil

}

// Versions returns a collection of all HardwareVersion entries from the index.
func (hvi HardwareVersionsIndex) Versions() HardwareVersions {

	newest := hvi.Newest().value
	versions := make([]HardwareVersion, 0, len(hvi))
	for hwv, count := range hvi {
		var isNewest bool
		if hwv == newest {
			isNewest = true
		}
		versions = append(versions, HardwareVersion{
			value:   hwv,
			count:   count,
			highest: isNewest,
		})
	}

	sort.Slice(versions, func(i, j int) bool {
		return strings.ToLower(versions[i].value) > strings.ToLower(versions[j].value)
	})

	return versions
}

// Outdated returns a collection of all older HardwareVersion.
func (hvi HardwareVersionsIndex) Outdated() HardwareVersions {

	newest := hvi.Newest().value
	var outliers []HardwareVersion
	for hwv, count := range hvi {
		if hwv != newest {
			outliers = append(outliers, HardwareVersion{
				value: hwv,
				count: count,
			})
		}
	}

	sort.Slice(outliers, func(i, j int) bool {
		return outliers[i].count > outliers[j].count
	})

	return outliers
}

// Newest returns the highest hardware version stored in the index. This value
// is returned as a HardwareVersion type, providing both the original vmx-123
// formatted string in addition to the actual version number.
func (hvi HardwareVersionsIndex) Newest() HardwareVersion {

	keys := make([]string, 0, len(hvi))
	for k := range hvi {
		keys = append(keys, k)
	}

	// highest version to the front to avoid potential negative slice indexing
	sort.Slice(keys, func(i, j int) bool { return keys[i] > keys[j] })

	highestVersion := keys[0]
	highestVersionCount := hvi[highestVersion]

	return HardwareVersion{
		value:   highestVersion,
		count:   highestVersionCount,
		highest: true,
	}
}

// Oldest returns the highest hardware version stored in the index. This value
// is returned as a HardwareVersion type, providing both the original vmx-123
// formatted string in addition to the actual version number.
func (hvi HardwareVersionsIndex) Oldest() HardwareVersion {

	keys := make([]string, 0, len(hvi))
	for k := range hvi {
		keys = append(keys, k)
	}

	// lowest version to the front to avoid potential negative slice indexing
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

	lowestVersion := keys[0]
	lowestVersionCount := hvi[lowestVersion]

	return HardwareVersion{
		value: lowestVersion,
		count: lowestVersionCount,

		// handled by default value
		// highest: false,
	}
}

// Count returns the number of hardware versions stored in the index.
func (hvi HardwareVersionsIndex) Count() int {
	return len(hvi)
}

// String is a Stringer implementation to return the original formatted string.
func (hv HardwareVersion) String() string {
	return hv.value
}

// Count returns the number of VirtualMachines with this specific virtual
// hardware version.
func (hv HardwareVersion) Count() int {
	return hv.count
}

// IsHighest indicates whether this HardwareVersion is the highest version in
// our inventory.
func (hv HardwareVersion) IsHighest() bool {
	return hv.highest
}

// VersionNumber returns the numeric version number of a VirtualMachine or -1
// if there was an issue converting the prefixed string value to a usable
// number.
func (hv HardwareVersion) VersionNumber() int {

	numStr := strings.Replace(hv.value, virtualHardwareVersionPrefix, "", 1)
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return -1
	}

	return num
}

// Sum provides the total count of all HardwareVersion entries.
func (hvs HardwareVersions) Sum() int {
	var sum int
	for i := range hvs {
		sum += hvs[i].count
	}

	return sum
}

// VersionNames returns a list of all hardware versions in their original
// string format.
func (hvs HardwareVersions) VersionNames() []string {

	names := make([]string, 0, len(hvs))
	for _, hwv := range hvs {
		names = append(names, hwv.value)
	}

	sort.Slice(names, func(i, j int) bool {
		return strings.ToLower(names[i]) < strings.ToLower(names[j])
	})

	return names
}

// VersionNumbers returns a list of all hardware versions in numerical format.
// -1 is returned for each hardware version if there was an issue converting
// the prefixed string value to a usable number.
func (hvs HardwareVersions) VersionNumbers() []int {

	versionNums := make([]int, 0, len(hvs))
	for _, hwv := range hvs {
		versionNums = append(versionNums, hwv.VersionNumber())
	}

	sort.Slice(versionNums, func(i, j int) bool {
		return versionNums[i] < versionNums[j]
	})

	return versionNums
}

// MeetsMinVersion accepts the minimum hardware version for all VMs and
// indicates whether all hardware versions meet or exceed the minimum.
func (hvs HardwareVersions) MeetsMinVersion(minVer int) bool {

	hvs.VersionNumbers()
	for _, num := range hvs.VersionNumbers() {
		if num < minVer {
			return false
		}
	}

	return true

}

// FilterVMsWithOldHardware filters the provided collection of VirtualMachines
// to just those with older hardware versions. The collection is returned
// along with the number of VirtualMachines that were excluded.
func FilterVMsWithOldHardware(vms []mo.VirtualMachine, hwIndex HardwareVersionsIndex) ([]mo.VirtualMachine, int) {

	var vmsWithOldHardware []mo.VirtualMachine
	for _, vm := range vms {
		if vm.Config.Version != hwIndex.Newest().String() {
			vmsWithOldHardware = append(vmsWithOldHardware, vm)
		}
	}

	sort.Slice(vmsWithOldHardware, func(i, j int) bool {
		return strings.ToLower(vmsWithOldHardware[i].Name) > strings.ToLower(vmsWithOldHardware[j].Name)
	})

	numExcluded := len(vms) - len(vmsWithOldHardware)

	return vmsWithOldHardware, numExcluded

}

// VirtualHardwareOneLineCheckSummary is used to generate a one-line Nagios
// service check results summary. This is the line most prominent in
// notifications.
func VirtualHardwareOneLineCheckSummary(
	stateLabel string,
	hwvIndex HardwareVersionsIndex,
	minHardwareVersion int,
	evaluatedVMs []mo.VirtualMachine,
	rps []mo.ResourcePool,
) string {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute VirtualHardwareOneLineCheckSummary func.\n",
			time.Since(funcTimeStart),
		)
	}()

	var outdatedVMs int
	minHardwareVersionString := fmt.Sprintf(
		"%s%d",
		virtualHardwareVersionPrefix,
		minHardwareVersion,
	)
	for _, vm := range evaluatedVMs {
		if vm.Config.Version == minHardwareVersionString {
			continue
		}

		hwVersion := NewHardwareVersion(vm.Config.Version)
		hwVerNum := hwVersion.VersionNumber()
		if hwVerNum < minHardwareVersion {
			outdatedVMs++
		}
	}

	switch {
	case outdatedVMs > 0:
		return fmt.Sprintf(
			"%s: %d VMs with hardware version older than %d (evaluated %d VMs, %d Resource Pools)",
			stateLabel,
			outdatedVMs,
			minHardwareVersion,
			len(evaluatedVMs),
			len(rps),
		)

	default:

		return fmt.Sprintf(
			"%s: No hardware versions older than %d detected (evaluated %d VMs, %d Resource Pools)",
			stateLabel,
			minHardwareVersion,
			len(evaluatedVMs),
			len(rps),
		)

	}
}

// VirtualHardwareReport generates a summary of virtual hardware details
// intended to aid in troubleshooting check results at a glance. This
// information is provided for use with the Long Service Output field commonly
// displayed on the detailed service check results display in the web UI or in
// the body of many notifications.
func VirtualHardwareReport(
	c *vim25.Client,
	hwvIndex HardwareVersionsIndex,
	minHardwareVersion int,
	defaultHardwareVersion HardwareVersion,
	allVMs []mo.VirtualMachine,
	evaluatedVMs []mo.VirtualMachine,
	vmsToExclude []string,
	evalPoweredOffVMs bool,
	includeRPs []string,
	excludeRPs []string,
	rps []mo.ResourcePool,
) string {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute VirtualHardwareReport func.\n",
			time.Since(funcTimeStart),
		)
	}()

	rpNames := make([]string, len(rps))
	for i := range rps {
		rpNames[i] = rps[i].Name
	}

	var report strings.Builder

	hardwareVersions := hwvIndex.Versions()
	hardwareVersions.MeetsMinVersion(minHardwareVersion)

	switch {

	// if we have more than one hardware version in the index, we have at
	// least one outdated version to report
	case hwvIndex.Count() > 1:

		fmt.Fprintf(
			&report,
			"Virtual Hardware Summary%s%s",
			nagios.CheckOutputEOL,
			nagios.CheckOutputEOL,
		)

		for _, hwv := range hwvIndex.Versions() {
			if !hwv.IsHighest() {
				fmt.Fprintf(
					&report,
					"version: %s, count: %d (outdated)\n",
					hwv.String(),
					hwv.Count(),
				)
				continue
			}
			fmt.Fprintf(
				&report,
				"version: %s, count: %d\n",
				hwv.String(),
				hwv.Count(),
			)
		}

	default:

		// homogenous

		fmt.Fprintf(
			&report,
			"All evaluated VMs are at hardware version %d.%s",
			hwvIndex.Newest().VersionNumber(),
			nagios.CheckOutputEOL,
		)

	}

	if !hardwareVersions.MeetsMinVersion(minHardwareVersion) {

		minHardwareVersionString := fmt.Sprintf(
			"%s%d",
			virtualHardwareVersionPrefix,
			minHardwareVersion,
		)

		fmt.Fprintf(
			&report,
			"%sVirtual Machines in need of upgrade:%s%s",
			nagios.CheckOutputEOL,
			nagios.CheckOutputEOL,
			nagios.CheckOutputEOL,
		)

		sort.Slice(evaluatedVMs, func(i, j int) bool {
			return evaluatedVMs[i].Config.Version < evaluatedVMs[j].Config.Version
		})

		for _, vm := range evaluatedVMs {
			if vm.Config.Version == minHardwareVersionString {
				continue
			}

			hwVersion := NewHardwareVersion(vm.Config.Version)
			hwVerNum := hwVersion.VersionNumber()
			if hwVerNum < minHardwareVersion {
				fmt.Fprintf(
					&report,
					"* %s (%s)%s",
					vm.Name,
					vm.Config.Version,
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
		"* Default Virtual Hardware Version: %d (%s) %s",
		defaultHardwareVersion.VersionNumber(),
		defaultHardwareVersion.String(),
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		&report,
		"* Newest Virtual Hardware Version: %d (%s) %s",
		hwvIndex.Newest().VersionNumber(),
		hwvIndex.Newest().String(),
		nagios.CheckOutputEOL,
	)

	fmt.Fprintf(
		&report,
		"* Oldest Virtual Hardware Version: %d (%s) %s",
		hwvIndex.Oldest().VersionNumber(),
		hwvIndex.Oldest().String(),
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
		evalPoweredOffVMs,
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
