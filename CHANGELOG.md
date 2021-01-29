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

[Unreleased]: https://github.com/atc0005/check-vmware/compare/v0.5.1...HEAD
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
