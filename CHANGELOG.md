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

[Unreleased]: https://github.com/atc0005/check-vmware/compare/v0.15.3...HEAD
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
