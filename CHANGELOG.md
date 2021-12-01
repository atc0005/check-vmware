# Changelog

## Overview

All notable changes to this project will be documented in this file.

The format is based on [Keep a
Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to
[Semantic Versioning](https://semver.org/spec/v2.0.0.html).

Please [open an issue](https://github.com/atc0005/check-vmware/issues) for any
deviations that you spot; I'm still learning!.

## Types of changes

The following types of changes will be recorded in this file:

- `Added` for new features.
- `Changed` for changes in existing functionality.
- `Deprecated` for soon-to-be removed features.
- `Removed` for now removed features.
- `Fixed` for any bug fixes.
- `Security` in case of vulnerabilities.

## [Unreleased]

- placeholder

## [v0.28.0] - 2021-XX-XX

### Overview

- placeholder
- Breaking changes
  - the `check_vmware_datastore` plugin has been renamed
- built using Go 1.16.10
  - Statically linked
  - Linux (x86, x64)

### Breaking

- (GH-510) `check_vmware_datastore` plugin
  - this plugin has been renamed to `check_vmware_datastore_space`

### Added

- placeholder

### Changed

- placeholder

### Fixed

- placeholder

## [v0.27.0] - 2021-12-01

### Overview

- New plugin
- Bugfixes
- Dependency updates
- Deprecated `check_vmware_datastore` plugin name
  - will be renamed in the `v0.28.0` release
- built using Go 1.16.10
  - Statically linked
  - Linux (x86, x64)

### Added

- (GH-505) New plugin: `check_vmware_datastore_performance`

### Changed

- Dependencies
  - `Go`
    - drop Go version in `go.mod` from `1.15` to `1.14`
      - attempt to reflect actual base Go version required by dependencies
  - `github/codeql-action`
    - `v1.0.21` to `v1.0.22`
  - `vmware/govmomi`
    - `v0.27.1` to `v0.27.2`
  - `actions/setup-node`
    - `v2.4.1` to `v2.5.0`

- (GH-511) Add `go-mod/go-version` shield to README
- (GH-517) Replace fully-qualified path to plugins/binaries in command
  definitions with `$USER1$` macro reference

### Deprecated

- (GH-510, GH-530) `check_vmware_datastore` plugin
  - this plugin will be renamed in the `v0.28.0` release to
    `check_vmware_datastore_space`
  - documentation, `Makefile` and other changes will be applied in the
    `v0.28.0` release to accommodate the rename of this plugin

### Fixed

- (GH-514) Add missing deferred timing log messages from `internal/vsphere`
  package functions
- (GH-515) Log message for `check_vmware_datastore` incorrect for `WARNING`
  state
- (GH-516) Boolean state checks for `DatastoreUsageSummary` trigger on
  threshold matches in addition to exceeding thresholds
- (GH-521) Boolean state checks for `HostSystemCPUSummary` trigger on
  threshold matches in addition to exceeding thresholds
- (GH-522) Boolean state checks for `HostSystemMemorySummary` trigger on
  threshold matches in addition to exceeding thresholds
- (GH-526) golangci-lint | WARN [runner] The linter 'golint' is deprecated;
  replaced by `revive`
- (GH-527) internal/config/config.go:23:13: var-declaration: should omit type
  string from declaration of var version; it will be inferred from the
  right-hand side (revive)

## [v0.26.0] - 2021-11-10

### Overview

- Expand support for Nagios Performance Data
- Bugfixes
- Dependency updates
- built using Go 1.16.10
  - Statically linked
  - Linux (x86, x64)

### Added

- Add additional Nagios Performance Data metrics
  - (GH-357) `check_vmware_vm_power_uptime` plugin
    - `time`
    - `vms`
    - `vms_with_critical_power_uptime`
    - `vms_with_warning_power_uptime`
    - `resource_pools_excluded`
    - `resource_pools_included`
    - `resource_pools_evaluated`
  - (GH-355) `check_vmware_vcpus` plugin
    - `time`
    - `vms`
    - `vms_excluded_by_name`
    - `vms_excluded_by_power_state`
    - `vcpus_usage`
    - `vcpus_used`
    - `vcpus_remaining`
    - `resource_pools_excluded`
    - `resource_pools_included`
    - `resource_pools_evaluated`
  - (GH-356) `check_vmware_vhw` plugin
    - `time`
    - `vms`
    - `vms_excluded_by_name`
    - `vms_excluded_by_power_state`
    - `hardware_versions_unique`
    - `hardware_versions_newest`
    - `hardware_versions_default`
    - `hardware_versions_oldest`
    - `resource_pools_excluded`
    - `resource_pools_included`
    - `resource_pools_evaluated`
  - (GH-348) `check_vmware_hs2ds2vms` plugin
    - `time`
    - `vms`
    - `vms_excluded_by_name`
    - `vms_excluded_by_power_state`
    - `pairing_issues`
    - `datastores`
    - `hosts`
    - `resource_pools_excluded`
    - `resource_pools_included`
    - `resource_pools_evaluated`
  - (GH-344) `check_vmware_alarms` plugin
    - `time`
    - `datacenters`
    - `triggered_alarms`
    - `triggered_alarms_excluded`
    - `triggered_alarms_included`
    - `triggered_alarms_critical`
    - `triggered_alarms_warning`
    - `triggered_alarms_unknown`
    - `triggered_alarms_ok`
  - (GH-350) `check_vmware_rps_memory` plugin
    - `time`
    - `vms`
    - `memory_usage`
    - `memory_used`
    - `memory_remaining`
    - `resource_pools_excluded`
    - `resource_pools_included`
    - `resource_pools_evaluated`

### Changed

- Dependencies
  - `Go`
    - `1.16.9` to `1.16.10`
  - `github/codeql-action`
    - `v1.0.21` to `v1.0.22`
  - `rs/zerolog`
    - `v1.25.0` to `v1.26.0`
  - `actions/checkout`
    - `v2.3.5` to `v2.4.0`

### Fixed

- (GH-484) `vsphere.DefaultHardwareVersion()`, `vsphere.NewHardwareVersion()`
  fail to set count of Virtual Machines at this version
- (GH-483) `check_vmware_vhw` | Clarify that only one monitoring mode (at a
  time) is supported
- (GH-489) `check_vmware_hs2ds2vms` plugin does not explicitly close session
- (GH-491) Makefile | `go get` executable installation deprecated
- (GH-494) Early exit logic for `vsphere.GetEligibleRPs()` is incomplete
- (GH-496) Incorrect calculation of memory remaining in evaluated resource
  pools for `check_vmware_rps_memory` plugin
- (GH-500) Use bytes as baseline value for RPs memory plugin
- (GH-333) Update documentation to reflect support for Performance Data output

## [v0.25.1] - 2021-11-02

### Overview

- Bugfixes
- built using Go 1.16.9
  - Statically linked
  - Linux (x86, x64)

### Changed

- Dependencies
  - `atc0005/go-nagios`
    - `v0.8.0` to `v0.8.1`

### Fixed

- (GH-473) Runtime error after deploying v0.25.0

## [v0.25.0] - 2021-11-01

### Overview

- Expand support for Nagios Performance Data
- Refactored `check_vmware_snapshots_*` plugins
- Append additional context to errors related to plugin runtime timeouts
- Bugfixes
- Dependency updates
- built using Go 1.16.9
  - Statically linked
  - Linux (x86, x64)

### Added

- Add additional Nagios Performance Data metrics
  - (GH-446, GH-459) `check_vmware_snapshots_age` plugin
    - `time`
    - `vms`
    - `vms_with_critical_snapshots`
    - `vms_with_warning_snapshots`
    - `snapshots`
    - `critical_snapshots`
    - `warning_snapshots`
    - `resource_pools_excluded`
    - `resource_pools_included`
    - `resource_pools_evaluated`
  - (GH-353) `check_vmware_snapshots_size` plugin
    - `time`
    - `vms`
    - `vms_with_critical_snapshots`
    - `vms_with_warning_snapshots`
    - `snapshots`
    - `critical_snapshots`
    - `warning_snapshots`
    - `resource_pools_excluded`
    - `resource_pools_included`
    - `resource_pools_evaluated`
  - (GH-352) `check_vmware_snapshots_count` plugin
    - `time`
    - `vms`
    - `vms_with_critical_snapshots`
    - `vms_with_warning_snapshots`
    - `snapshots`
    - `critical_snapshots`
    - `warning_snapshots`
    - `resource_pools_excluded`
    - `resource_pools_included`
    - `resource_pools_evaluated`

### Changed

- Dependencies
  - `github/codeql-action`
    - `v1.0.19` to `v1.0.21`

- (GH-463) Annotate `context deadline exceeded` errors with additional human
  readable details

### Fixed

- (GH-448) Tweak logging for datastore plugin
- (GH-451) Update `check_vmware_snapshots_*` plugins to track crossing of
  thresholds, use methods to determine exactly one of CRITICAL or WARNING
  state
- (GH-455) Snapshots list headers for the `check_vmware_snapshots_count`
  plugin do not have the right context
- (GH-461) Provide time (runtime) performance data metric for all exit points
  from plugins
- (GH-462) `Error occurred while destroying view` message emitted to stdout
- (GH-468) Threshold values for `memory_usage` metric for
  `check_vmware_host_memory` plugin are always zero

## [v0.24.0] - 2021-10-25

### Overview

- Expand support for Nagios Performance Data
- Tweaks to `check_vmware_question` plugin
- Bugfixes
- built using Go 1.16.9
  - Statically linked
  - Linux (x86, x64)

### Added

- Add additional Nagios Performance Data metrics
  - (GH-435) `check_vmware_host_cpu` plugin
    - `cpu_used`
  - (GH-435) `check_vmware_host_memory` plugin
    - `memory_used`
  - (GH-435) `check_vmware_datastore` plugin
    - `datastore_storage_used`

### Changed

- (GH-422) Extend `check_vmware_question` plugin to list pending question(s)
  in `LongServiceOutput`

### Fixed

- (GH-436) Incorrect UoM used by `memory_total` metric for
  `check_vmware_host_memory` plugin
- (GH-439) Potential nil pointer dereference in
  `vsphere.VMInteractiveQuestionReport()` func
- (GH-442) Incorrect error value returned by `check_vmware_question plugin`
  when interactive response needed

## [v0.23.0] - 2021-10-22

### Overview

- Expand support for Nagios Performance Data
- Refactoring/cleanup
- Bugfixes
- Dependency updates
- built using Go 1.16.9
  - Statically linked
  - Linux (x86, x64)

### Added

- Add additional Nagios Performance Data metrics
  - (GH-346) `check_vmware_host_cpu` plugin
    - `time`
    - `cpu_usage`
    - `cpu_total`
    - `cpu_remaining`
    - `vms`
    - `vms_powered_off`
    - `vms_powered_on`
  - (GH-347) `check_vmware_host_memory` plugin
    - `time`
    - `memory_usage`
    - `memory_total`
    - `memory_remaining`
    - `vms`
    - `vms_powered_off`
    - `vms_powered_on`
  - (GH-349) `check_vmware_question` plugin
    - `time`
    - `vms`
    - `vms_excluded_by_name`
    - `vms_requiring_input`
    - `vms_not_requiring_input`
    - `resource_pools_excluded`
    - `resource_pools_included`
    - `resource_pools_evaluated`
  - (GH-345) `check_vmware_disk_consolidation` plugin
    - `time`
    - `vms`
    - `vms_excluded_by_name`
    - `vms_with_consolidation_need`
    - `vms_without_consolidation_need`
    - `resource_pools_excluded`
    - `resource_pools_included`
    - `resource_pools_evaluated`
- `check_vmware_disk_consolidation` plugin
  - (GH-419) Add `--trigger-reload` flag to (optionally) force state reload
    for Virtual Machines before checking their disk consolidation status
    - NOTE: Omitting the flag retains previous plugin behavior of only
      checking the current state as recorded by vSphere (faster, but with
      potentially stale data)

### Changed

- Dependencies
  - `vmware/govmomi`
    - `v0.27.0` to `v0.27.1`

### Fixed

- (GH-278) Issue with `check_vmware_disk_consolidation` plugin
- (GH-293) Review `check_vmware_disk_consolidation` results
- (GH-424) Incorrect CPU speed UoM provided by `(vsphere.CPUSpeed).String()`
- (GH-425) Potential nil pointer dereference in
  `vsphere.NewHostSystemCPUUsageSummary()`
- (GH-430) Potential nil pointer dereference in
  `vsphere.NewHostSystemMemoryUsageSummary()`

## [v0.22.0] - 2021-10-19

### Overview

- Refactoring/cleanup
- Bugfixes
- Dependency updates
- built using Go 1.16.9
  - Statically linked
  - Linux (x86, x64)

### Added

- (GH-408) Add `quick` Makefile recipe

### Changed

- Dependencies
  - `actions/checkout`
    - `v2.3.4` to `v2.3.5`

- (GH-394) Refactor check_vmware_hs2ds2vms plugin

### Fixed

- (GH-410) Dependabot did not create a PR for actions/checkout 2.3.5
- (GH-413) Update GHAW doc comment headers to list correct repo
- (GH-414) Timeout handling is unclear: connection timeout or plugin
  runtime/execution timeout?

## [v0.21.1] - 2021-10-16

### Overview

- Refactoring/cleanup
- Bugfixes
- Dependency updates
- built using Go 1.16.9
  - Statically linked
  - Linux (x86, x64)

### Added

- (GH-389) Add CodeQL analysis GitHub Actions workflow

### Changed

- (GH-406) Cover Performance Data metrics work in README
- (GH-403) Note type of container for retrieved VMs via debug log message

- Dependencies
  - `vmware/govmomi`
    - `v0.26.1` to `v0.27.0`

### Fixed

- (GH-391) README CLI example for `check_vmware_hs2ds2vms` lists
  `check_vmware_vhw` plugin instead
- (GH-393) `check_vmware_hs2ds2vms` plugin emits empty string in place of ESXi
  host for mismatched VM
- (GH-397) Unable to retrieve VirtualMachine objects from VirtualApp
  "containers"
- (GH-400) When called from `vsphere.GetVMsFromContainer()`, the
  `vsphere.getObjects()` function reports incorrect number of objects
  retrieved
- (GH-399) `vsphere.GetEligibleRPs()` should exit early if the hidden/parent
  `Resources` ResourcePool is evaluated
  - light testing indicates that this can reduce plugin runtime for applicable
    plugins as much as 2x and is particularly notable when evaluating 200+ VMs

## [v0.21.0] - 2021-10-13

### Overview

- Expand support for Nagios Performance Data
- Refactoring/cleanup
- Bugfixes
- Dependency updates
- built using Go 1.16.9
  - Statically linked
  - Linux (x86, x64)

### Added

- Add additional Nagios Performance Data metrics
  - `check_vmware_datastore` plugin
    - `time`
    - `vms`
    - `vms_powered_on`
    - `vms_powered_off`
  - `check_vmware_tools` plugin
    - `time`
    - `vms`
    - `vms_excluded_by_name`
    - `vms_excluded_by_power_state`
    - `vms_with_tools_issues`
    - `vms_without_tools_issues`
    - `resource_pools_excluded`
    - `resource_pools_included`
    - `resource_pools_evaluated`
- Add project security policy

### Changed

- Dependencies
  - `Go`
    - `1.16.8` to `1.16.9`

- `check_vmware_datastore` plugin
  - (GH-361) List number of VMs on datastore along with their power state
  - (GH-362) Refactor datastore usage summary generation
- `check_vmware_tools` plugin
  - refactor: Update Tools summary func to provide ServiceState
  - (GH-366, GH-370) Consider additional VMware Tools status fields when
    determining overall plugin state (`toolsStatus` is deprecated, see README
    updates for details)
  - (GH-372) Extend structured logging fields to include Performance Data
    metrics
- `check_vmware_host_cpu` plugin
  - (GH-377) Log VMs found running on host
- `check_vmware_host_memory` plugin
  - (GH-378) Log VMs found running on host
- `internal/vsphere` package
  - (GH-375) Evaluate suspended VMs as powered off VMs
  - (GH-373) Update `FilterXByY()` and `ExcludeXByY()` functions to return the
    number of excluded items
  - (GH-380) Update filter func names that return a single value to still use
    plural source naming

### Removed

- Linux x86 assets no longer provided for new releases
  - still available via local builds (e.g., `make linux` or `make linux-x86`)

### Fixed

- (GH-342) Download links generated incorrectly
- (GH-340) Missing footnote reference for v0.20.0 release
- (GH-358) Potential nil pointer dereference in `vsphere.FilterVMByID()`
- (GH-359) Potential nil pointer dereference in `vsphere.dedupeVMs()`
- (GH-365) Powered off VMs without running Tools are listed as problem VMs
  when `--powered-off` flag is specified
- (GH-374) The `vphere.FilterVMsByPowerState()` function should evaluate
  suspended VMs

## [v0.20.0] - 2021-10-01

### Overview

- Add (limited) initial support for Nagios Performance Data output
- Add list of download links
- Build tweaks
- DEPRECATE Linux x86 assets
- Dependency updates
- built using Go 1.16.8
  - Statically linked
  - Linux (x86, x64)

### Added

- Add (limited) initial support for Nagios Performance Data output
  - initial support provided for the `check_vmware_datastore` plugin
  - updates to other plugins are planned for future releases

- List of download links is now provided with releases

### Changed

- Dependencies
  - `atc0005/go-nagios`
    - `v0.7.0` to `v0.8.0`
  - `actions/setup-node`
    - `v2.4.0` to `v2.4.1`

- `Makefile`
  - text file containing download links is generated by build recipes
  - new `release-build` recipe intended to act as a wrapper for all tasks
    required to generate assets for public release
  - existing `windows` and `linux` recipes split into platform-arch specific
    recipes
    - `windows` and `linux` recipes retained as wrappers for generating x86
      and x64 assets for each platform
  - `clean` recipe
    - dependency removed from `windows` and `linux` recipes and
      platform-arch specific recipes
    - now used only by `all` and `release-build` recipes
    - empty plugin-specific subdirectories are now removed

### Deprecated

- Linux 32-bit (x86) binaries have been discontinued for now
  - if there is interest, they can be provided in future releases
  - see GH-192 for details (and to leave feedback)

## [v0.19.1] - 2021-09-20

### Overview

- Dependency updates
- built using Go 1.16.8
  - Statically linked
  - Linux (x86, x64)

### Changed

- Dependencies
  - `Go`
    - `1.16.7` to `1.16.8`
  - `rs/zerolog`
    - `v1.23.0` to `v1.25.0`
  - `vmware/govmomi`
    - `v0.26.0` to `v0.26.1`
  - `atc0005/go-nagios`
    - `v0.6.0` to `v0.7.0`

- Documentation
  - Update README coverage for `stderr` output

## [v0.19.0] - 2021-08-11

### Overview

- Dependency updates
- built using Go 1.16.7
  - Statically linked
  - Linux (x86, x64)

### Changed

- List associated datastore for each snapshot (in place of MOID value
  previously used)

### Fixed

- Add handling of potential nil pointer for VM's snapshot property (of type
  `VirtualMachineSnapshotInfo`)

- README
  - Fix alarm plugin thresholds

## [v0.18.1] - 2021-08-06

### Overview

- Dependency updates
- built using Go 1.16.7
  - Statically linked
  - Linux (x86, x64)

### Changed

- Dependencies
  - `Go`
    - `1.16.6` to `1.16.7`
  - `actions/setup-node`
    - updated from `v2.3.0` to `v2.4.0`

## [v0.18.0] - 2021-07-25

### Overview

- Plugin improvements
- Dependency update
- built using Go 1.16.6
  - Statically linked
  - Linux (x86, x64)

### Changed

- `check_vmware_vpus` plugin
  - list 10 most vCPU consuming VMs in extended output
  - list 10 most recently booted VMs in extended output

- Dependencies
  - `actions/setup-node`
    - updated from `v2.2.0` to `v2.3.0`

## [v0.17.5] - 2021-07-15

### Overview

- Dependency update
- built using Go 1.16.6
  - Statically linked
  - Linux (x86, x64)

### Changed

- Dependencies
  - `atc0005/go-nagios`
    - `v0.6.0` to `v0.6.1`
  - `Go`
    - canary file updated from `1.16.5` to `1.16.6`

### Fixed

- CHANGELOG
  - Add missing Go dependency update

## [v0.17.4] - 2021-07-13

### Overview

- Output tweak
- Dependency updates
- built using Go 1.16.6
  - Statically linked
  - Linux (x86, x64)

### Added

- Add "canary" Dockerfile to track stable Go releases, serve as a reminder to
  generate fresh binaries

### Changed

- `check_vmware_disk_consolidation` plugin
  - list the power state for VMs in need of disk consolidation

- Dependencies
  - `Go`
    - `1.16.5` to `1.16.6`
  - `actions/setup-node`
    - updated from `v2.1.5` to `v2.2.0`
    - update `node-version` value to always use latest LTS version instead of
      hard-coded version

## [v0.17.3] - 2021-06-27

### Overview

- Output tweaks
- Minor cleanup & bug fixes
- built using Go 1.16.5
  - Statically linked
  - Linux (x86, x64)

### Added

- contrib files
  - Update send2teams command definition to use new `target-url` flag from
    [send2teams v0.6.0
    release](https://github.com/atc0005/send2teams/releases/tag/v0.6.0)
- plugin output
  - Set & report custom User Agent as shared value for all plugins

### Changed

- Triggered Alarms | Extend missing resource pool error message

### Fixed

- contrib files
  - check_vmware_alarms | vc1 host service checks contrib file has several
    misconfigured checks
  - check_vmware_alarms | command definitions contrib file has several
    duplicated doc comments

## [v0.17.2] - 2021-06-17

### Overview

- Plugin output tweak
- built using Go 1.16.5
  - Statically linked
  - Linux (x86, x64)

### Changed

- Expand triggered alarms listing in report to help indicate *why* an alarm
  was excluded

## [v0.17.1] - 2021-06-16

### Overview

- Bug fixes
- built using Go 1.16.5
  - Statically linked
  - Linux (x86, x64)

### Fixed

- Missing date for prior version header
- logger includes fields not used by this plugin
- plugin_type structured logger field empty

## [v0.17.0] - 2021-06-16

### Overview

- Improved `check_vmware_alarms` plugin
- Bug fixes
- Dependency updates
- built using Go 1.16.5
  - Statically linked
  - Linux (x86, x64)

### Added

- Add contrib files coverage for alarms plugin
  - initial coverage of v0.16.0 `check_vmware_alarms` plugin functionality
  - expanded coverage of new filtering modes
- Add support for specifying multiple datacenters
- Add support for ignoring (or explicitly including) triggered alarms by alarm
  description or name substrings
- Add support for filtering triggered alarms by alarm description substrings
- Add support for filtering triggered alarms by alarm name substrings
- Add support for filtering triggered alarms by specific status keywords
- Add support for filtering triggered alarms by associated entity name
- Add support for filtering triggered alarms by associated entity resource
  pool
- Add test coverage for filtering behavior

### Changed

- Dependencies
  - `Go`
    - `1.16.4` to `1.16.5`
  - `rs/zerolog`
    - `v1.22.0` to `v1.23.0`
  - `vmware/govmomi`
    - `v0.25.0` to `v0.26.0`

- Linting
  - Re-enable (deprecated) `maligned` linter
  - Disable `fieldalignment` settings
    - until the Go team offers more control over the types of checks provided
      by the `fieldalignment` linter or `golangci-lint` does so.

- Datacenters handling
  - reverse the decision to fallback to the default datacenter if the
    specified list is not found and instead consider any missing requested
    datacenter names as a `CRITICAL` error

### Fixed

- Invalid logic when filtering by acknowledged state
- README and command definition incorrectly uses a flag in the "show
  everything" example
- Virtual Hardware check documentation outdated

## [v0.16.0] - 2021-05-27

### Overview

- New plugin
- Bug fix
- built using Go 1.16.4
  - Statically linked
  - Linux (x86, x64)

### Added

- New plugin: `check_vmware_alarms`

### Fixed

- check_vmware_question plugin incorrectly labels VMs as requiring an
  interactive response as needing disk consolidation

## [v0.15.3] - 2021-05-21

### Overview

- Bug fixes
- Dependency updates
- built using Go 1.16.4
  - Statically linked
  - Linux (x86, x64)

### Changed

- Dependencies
  - `Go`
    - `1.16.3` to `1.16.4`
  - `rs/zerolog`
    - `v1.21.0` to `v1.22.0`
  - `vmware/govmomi`
    - `v0.24.1` to `v0.25.0`

### Fixed

- Login error message does not interpolate server value as intended
- Missing section header in CHANGELOG
- Stray newline breaks intended README Table of Contents ordering
- Doc comments missing from `getObjects` func

## [v0.15.2] - 2021-04-02

### Overview

- Bug fixes
- built using Go 1.16.3
  - Statically linked
  - Linux (x86, x64)

### Changed

- Dependencies
  - `Go`
    - `1.16.2` to `1.16.3`

### Fixed

- linting
  - fieldalignment: struct with X pointer bytes could be Y (govet)
  - `golangci/golangci-lint`
    - replace deprecated `maligned` linter with `govet: fieldalignment`
    - replace deprecated `scopelint` linter with `exportloopref`

## [v0.15.1] - 2021-03-30

### Overview

- Bug fixes
- built using Go 1.16.2
  - Statically linked
  - Linux (x86, x64)

### Fixed

- CHANGELOG
  - Fix `Deprecated` section header level for v0.15.0 release entry of this
    CHANGELOG

- `check_vmware_question` plugin
  - invalid threshold text/label

## [v0.15.0] - 2021-03-30

### Overview

- New plugin
- Bug fixes
- built using Go 1.16.2
  - Statically linked
  - Linux (x86, x64)

### Added

- New plugin: `check_vmware_question`

### Deprecated

- 32-bit binaries are *likely* going to be dropped from future releases
  - if there is interest, they will continue to be provided
  - see GH-192 for details (and to leave feedback)

### Fixed

- CHANGELOG
  - Fix `Deprecated` section header level for v0.12.0 release entry of this
    CHANGELOG

- `contrib`
  - multiple mistakes in `vc1.example.com.cfg` host config file
    - missing service group
    - referenced wrong service group
    - service description in consistent with existing entries
  - `vmware-disk-consolidation` command definition doc comments regarding
    powered on/off VM state evaluation incorrect

- internal constants do not reflect `check_vmware_disk_consolidation` plugin
  rename

- README
  - `vmware-disk-consolidation` command definition doc comments regarding
    powered on/off VM state evaluation incorrect

## [v0.14.0] - 2021-03-29

### Overview

- New plugin
- Misc tweaks
- Bug fixes
- built using Go 1.16.2
  - Statically linked
  - Linux (x86, x64)

### Added

- New plugin: `check_vmware_disk_consolidation`

### Changed

- README
  - remove `Status` column from plugin table
  - replace GH issue links with GH discussion links
    - GH-160 to GH-178
    - GH-79 to GH-176

- `check_vmware_vhw` plugin
  - thresholds output text adjusted

- Refactoring
  - very minor tweaks, much, much more to do

- Dependencies
  - `rs/zerolog`
    - `v1.20.0` to `v1.21.0`
  - `vmware/govmomi`
    - `v0.24.0` to `v0.24.1`

### Fixed

- Doc comments
  - typo
  - copy/paste/forget

- `contrib`
  - multiple mistakes in `vc1.example.com.cfg` host config file

- `check_vmware_datastore` plugin
  - remove `included_resource_pools` and `excluded_resource_pools` fields from
    log messages as they do not apply to this plugin

- README
  - invalid `host-name` flag used in `check_vmware_vm_power_update` examples

## [v0.13.1] - 2021-03-17

### Overview

- Bug fixes
- built using Go 1.16.2
  - Statically linked
  - Linux (x86, x64)

### Changed

- Dependencies
  - `Go`
    - `1.15.8` to `1.16.2`
  - `actions/setup-node`
    - `v2.1.4` to `v2.1.5`

### Fixed

- `check_vmware_vm_power_uptime` plugin
  - incorrect exit code returned for `WARNING` state

## [v0.13.0] - 2021-02-18

### Overview

- Plugin improvements
- Bug fixes
- built using Go 1.15.8
- *shelve* building Windows binaries until feedback is provided

### Breaking

Precompiled Windows binaries have been discontinued. If this affects you,
please provide feedback to
[GH-160](https://github.com/atc0005/check-vmware/issues/160).

See the build instructions in the project README file for other options.

### Added

- New monitoring mode for `check_vmware_vhw` plugin
  - `Default is minimum required version`
  - Asserts that all evaluated Virtual Machines meet or exceed the default
    version specified by either a specified host or a specified cluster

### Changed

- Makefile
  - add `clean` recipe for single OS builds (`linux`, `windows`) to match what
    is already done for the `all` recipe

- README
  - Update coverage of installation steps
  - Note shelving of `Windows` builds until feedback is received indicating
    that others find them useful
  - Add "other OSes" section inviting requests for precompiled binaries for
    other supported OSes

### Fixed

- Missing plugin "type" in log messages due to missing constant/check in
  logging setup
- Invalid preallocaction of vhw slices used by `Oldest` and `Newest`
  `HardwareVersionsIndex` methods
- Add `Deprecated` section to v0.12.0 release entry of this CHANGELOG
  - this matches the current Release entry for this project's GH repo (which
    was previously labeled as Breaking)

## [v0.12.0] - 2021-02-16

### Overview

- Plugin improvements
- built using Go 1.15.8

### Deprecated

- Windows binaries are *likely* not going to be included in the next release
  - if there is interest, they will continue to be provided
  - see GH-160 for details (and to leave feedback)

### Changed

- `check_vmware_vm_power_uptime` plugin
  - list 10 most recently booted VMs in extended output
- `check_vmware_rps_memory` plugin
  - list 10 most recently booted VMs in extended output
  - list 10 most memory consuming VMs in extended output

### Fixed

- Typo in plugin name in prior CHANGELOG version entry

## [v0.11.0] - 2021-02-12

### Overview

- Bug fixes
- Misc tweaks
- built using Go 1.15.8

### Changed

- `check_vmware_vhw` plugin
  - two new additional monitoring modes added
    - minimum required version check
    - outdated-by or threshold range check
- CI build timeout adjusted from `20` to `40` minutes
- `check_vmware_vm_power_uptime` plugin
  - Extend check_vmware_vm_power_uptime to list 5-10 highest uptime VMs when
    state is OK

### Fixed

- CI-driven Makefile builds timing out after v0.10.0 release
- "shorthand" suffix missing for plugin-specific help output for short flag
  options
- Validation checks for CRITICAL/WARNING threshold checks does not (negative)
  assert `CRITICAL <= WARNING`

## [v0.10.0] - 2021-02-09

### Overview

- New plugin
- Misc tweaks
- built using Go 1.15.8

### Added

- New plugin: `check_vmware_vm_power_uptime`

### Changed

- minor refactor of summary and report functions

## [v0.9.0] - 2021-02-07

### Overview

- New plugin
- built using Go 1.15.8

### Added

- New plugin: `check_vmware_host_cpu`

### Changed

- Dependencies
  - built using Go 1.15.8
    - Statically linked
    - Windows (x86, x64)
    - Linux (x86, x64)

## [v0.8.0] - 2021-02-04

### Overview

- New plugin
- Bug fixes
- Minor tweaks
- built using Go 1.15.7

### Added

- New plugin: `check_vmware_host_memory`

### Changed

- Remove doc comments regarding hidden resource pool from
  `vc1.example.com.cfg` Nagios host configuration file
  - this topic is already covered by the main README file
  - some details have changed slightly since the remarks were written

## [v0.7.0] - 2021-02-04

### Overview

- New plugin
- Bug fixes
- Misc adjustments to output
- built using Go 1.15.7

### Added

- New plugin: `check_vmware_snapshots_count`

### Changed

- `check_vmware_rp_memory` plugin
  - one-line summary tweaks
    - list usage percentage of total capacity as aggregate value for all
      specified resource pools
    - attempt to make message more concise
    - drop explicit "overage" detail, rely on implied overage in total memory
      usage percentage
  - extended output
    - list usage percentage of total capacity per-Resource Pool
- Review and update threshold listings in extended output

### Fixed

- Fix invalid `ExceedsAge` logic
- `check_vmware_snapshots_age`: Misreported VMs, snapshots count

## [v0.6.1] - 2021-02-02

### Overview

- Bug fix
- built using Go 1.15.7

### Fixed

- Snapshots age evaluation counts misreported (swapped VMs, snapshots)

## [v0.6.0] - 2021-02-01

### Overview

- New plugin
- Bug fixes
- built using Go 1.15.7

### Added

- New plugin: `check_vmware_rps_memory`

### Changed

- `check_vmware_vcpus`
  - Adjust errors listing to use fixed error for crossing vCPUs allocation
    thresholds vs dynamically generated details (GH-104, GH-108)
    - the dynamic version mostly duplicated what the one-line summary was
      already conveying

### Fixed

- Fix variable name with stutter
- `check_vmware_datastore` command definition template uses wrong plugin name

## [v0.5.1] - 2021-01-29

### Overview

- Bug fixes
- built using Go 1.15.7

### Changed

- `internal/vsphere` package logging output disabled by default, exposed via
  `debug` logging level (user configurable)

### Fixed

- snapshots size plugin properly detects `WARNING` cumulative size state, but
  unhelpfully notes 0 (individual) snapshots exceeding size

## [v0.5.0] - 2021-01-27

### Overview

- New plugin
- Bug fixes
- built using Go 1.15.7

### Added

- New plugin: `check_vmware_snapshots_size`

### Changed

- Makefile: indent output per plugin build step

### Fixed

- check_vmware_tools plugin does not clearly define what thresholds are used
  for service check logic
- GoDoc coverage missing for project plugins

## [v0.4.3] - 2021-01-26

### Overview

- Bug fixes
- built using Go 1.15.7

### Changed

- The default `Resources` Resource Pool is now evaluated *unless* a Resource
  Pool is explicitly *included* via the `--include-rp` flag
  - previously VirtualMachine objects *outside* of a Resource Pool were
    unintentionally ignored
  - affected multiple plugins (see Fixed section)
  - credit: bug report from @HisArchness via Twitter

### Fixed

- check_vmware_tools: Long Service Output listing omits affected VMs when only one affected
- VMs outside of Resource Pools excluded from evaluation
  - multiple plugins affected
    - `check_vmware_tools`
    - `check_vmware_hs2ds2vms`
    - `check_vmware_snapshots_age`
    - `check_vmware_vhw`
    - `check_vmware_vcpus`

## [v0.4.2] - 2021-01-25

### Overview

- Bug fixes
- built using Go 1.15.7

### Changed

- Swap out GoDoc badge for pkg.go.dev badge
- Dependencies
  - built using Go 1.15.7
    - Statically linked
    - Windows (x86, x64)
    - Linux (x86, x64)

### Fixed

- check_vmware_snapshots_age plugin: duplicated structured logging field
- "Snapshots *not yet* exceeding age thresholds" list not populated
- Replace "please submit issue" request text with link to recently created
  issue for feedback collection
- Misc typos, copy/paste/modify mistakes

## [v0.4.1] - 2021-01-19

### Overview

- Bug fixes
- built using Go 1.15.6
  - Statically linked
  - Windows (x86, x64)
  - Linux (x86, x64)

### Fixed

- check_vmware_snapshots_age plugin: incomplete logic for young snapshots
  switch case
- check_vmware_snapshots_age plugin: wrong state label for OK check results

## [v0.4.0] - 2021-01-19

### Overview

- New plugin
- Bug fixes
- built using Go 1.15.6
  - Statically linked
  - Windows (x86, x64)
  - Linux (x86, x64)

### Added

- New plugin: `check_vmware_snapshots_age`

### Fixed

- check_vmware_datastore | Datastore-specific storage usage for VMs appears to
  be incorrect
- check_vmware_datastore | Datastore-specific storage usage for VMs is rounded
  without sufficient precision
- check_vmware_datastore | Angle brackets for pre tags (in VM listing)
  shown in CLI, missing from Nagios generated notifications

## [v0.3.0] - 2021-01-14

### Overview

- New plugin
- CI build timeout tweak
- Bug fixes
- built using Go 1.15.6
  - Statically linked
  - Windows (x86, x64)
  - Linux (x86, x64)

### Added

- New plugin: `check_vmware_datastore`
  - used to monitor space usage on datastores (one per service check)

### Changed

- GitHub Actions Workflow: `Build codebase using Makefile`
  - build timeout adjusted from 10 minutes to 20 minutes

### Fixed

- `vsphere.getObjectByName`
  - `RetrieveOne` PropertyCollector method called with an empty interface
    causing `unexpected type` panic
- `sphere.getObjects`
  - accepts unsupported `types.ManagedObjectReference` for use with
    `CreateContainerView`

## [v0.2.1] - 2021-01-14

### Overview

- documentation updates
- dependency update

- built using Go 1.15.6
  - Statically linked
  - Windows (x86, x64)
  - Linux (x86, x64)

### Added

- Add "contrib" files
  - demonstrate plugin usage

### Changed

- dependencies
  - `atc0005/go-nagios`
    - `v0.5.3` to `v0.6.0`

### Fixed

- Update documentation for v0.2.x release
  - examples
  - contrib overview
  - minor fixes
- Fix typo in project breadcrumb URL

## [v0.2.0] - 2021-01-12

### Added

- New plugin: `check_vmware_hs2ds2vms`
  - used to assert that VMs are running and housed on intended hosts and
    datastores

### Changed

- Allow specifying lists of values via CLI flag with or without quotes

### Fixed

- Expose via Long Service Output whether powered off VMs are evaluated
- Wire-up user domain support

## [v0.1.1] - 2021-01-11

### Fixed

- Incorrect project name in version output
- Plugins require write permission on home directory in order to cache login
  sessions

## [v0.1.0] - 2021-01-06

Initial release!

This release provides early versions of several Nagios plugins used to monitor
VMware vSphere environments (with more hopefully on the way soon).

### Added

- Nagios plugin for monitoring VMware Tools for select (or all) Resource
  Pools.

- Nagios plugin for monitoring virtual CPU allocations for select (or all)
  Resource Pools.

- Nagios plugin for monitoring virtual hardware versions for select (or all)
  Resource Pools.

[Unreleased]: https://github.com/atc0005/check-vmware/compare/v0.27.0...HEAD
[v0.28.0]: https://github.com/atc0005/check-vmware/releases/tag/v0.28.0
[v0.27.0]: https://github.com/atc0005/check-vmware/releases/tag/v0.27.0
[v0.26.0]: https://github.com/atc0005/check-vmware/releases/tag/v0.26.0
[v0.25.1]: https://github.com/atc0005/check-vmware/releases/tag/v0.25.1
[v0.25.0]: https://github.com/atc0005/check-vmware/releases/tag/v0.25.0
[v0.24.0]: https://github.com/atc0005/check-vmware/releases/tag/v0.24.0
[v0.23.0]: https://github.com/atc0005/check-vmware/releases/tag/v0.23.0
[v0.22.0]: https://github.com/atc0005/check-vmware/releases/tag/v0.22.0
[v0.21.1]: https://github.com/atc0005/check-vmware/releases/tag/v0.21.1
[v0.21.0]: https://github.com/atc0005/check-vmware/releases/tag/v0.21.0
[v0.20.0]: https://github.com/atc0005/check-vmware/releases/tag/v0.20.0
[v0.19.1]: https://github.com/atc0005/check-vmware/releases/tag/v0.19.1
[v0.19.0]: https://github.com/atc0005/check-vmware/releases/tag/v0.19.0
[v0.18.1]: https://github.com/atc0005/check-vmware/releases/tag/v0.18.1
[v0.18.0]: https://github.com/atc0005/check-vmware/releases/tag/v0.18.0
[v0.17.5]: https://github.com/atc0005/check-vmware/releases/tag/v0.17.5
[v0.17.4]: https://github.com/atc0005/check-vmware/releases/tag/v0.17.4
[v0.17.3]: https://github.com/atc0005/check-vmware/releases/tag/v0.17.3
[v0.17.2]: https://github.com/atc0005/check-vmware/releases/tag/v0.17.2
[v0.17.1]: https://github.com/atc0005/check-vmware/releases/tag/v0.17.1
[v0.17.0]: https://github.com/atc0005/check-vmware/releases/tag/v0.17.0
[v0.16.0]: https://github.com/atc0005/check-vmware/releases/tag/v0.16.0
[v0.15.3]: https://github.com/atc0005/check-vmware/releases/tag/v0.15.3
[v0.15.2]: https://github.com/atc0005/check-vmware/releases/tag/v0.15.2
[v0.15.1]: https://github.com/atc0005/check-vmware/releases/tag/v0.15.1
[v0.15.0]: https://github.com/atc0005/check-vmware/releases/tag/v0.15.0
[v0.14.0]: https://github.com/atc0005/check-vmware/releases/tag/v0.14.0
[v0.13.1]: https://github.com/atc0005/check-vmware/releases/tag/v0.13.1
[v0.13.0]: https://github.com/atc0005/check-vmware/releases/tag/v0.13.0
[v0.12.0]: https://github.com/atc0005/check-vmware/releases/tag/v0.12.0
[v0.11.0]: https://github.com/atc0005/check-vmware/releases/tag/v0.11.0
[v0.10.0]: https://github.com/atc0005/check-vmware/releases/tag/v0.10.0
[v0.9.0]: https://github.com/atc0005/check-vmware/releases/tag/v0.9.0
[v0.8.0]: https://github.com/atc0005/check-vmware/releases/tag/v0.8.0
[v0.7.0]: https://github.com/atc0005/check-vmware/releases/tag/v0.7.0
[v0.6.1]: https://github.com/atc0005/check-vmware/releases/tag/v0.6.1
[v0.6.0]: https://github.com/atc0005/check-vmware/releases/tag/v0.6.0
[v0.5.1]: https://github.com/atc0005/check-vmware/releases/tag/v0.5.1
[v0.5.0]: https://github.com/atc0005/check-vmware/releases/tag/v0.5.0
[v0.4.3]: https://github.com/atc0005/check-vmware/releases/tag/v0.4.3
[v0.4.2]: https://github.com/atc0005/check-vmware/releases/tag/v0.4.2
[v0.4.1]: https://github.com/atc0005/check-vmware/releases/tag/v0.4.1
[v0.4.0]: https://github.com/atc0005/check-vmware/releases/tag/v0.4.0
[v0.3.0]: https://github.com/atc0005/check-vmware/releases/tag/v0.3.0
[v0.2.1]: https://github.com/atc0005/check-vmware/releases/tag/v0.2.1
[v0.2.0]: https://github.com/atc0005/check-vmware/releases/tag/v0.2.0
[v0.1.1]: https://github.com/atc0005/check-vmware/releases/tag/v0.1.1
[v0.1.0]: https://github.com/atc0005/check-vmware/releases/tag/v0.1.0
