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

## [v0.36.16] - 2024-11-11

### Changed

#### Dependency Updates

- (GH-1300) Build Image: Bump atc0005/go-ci from go-ci-oldstable-build-v0.21.14 to go-ci-oldstable-build-v0.21.15 in /dependabot/docker/builds
- (GH-1294) Go Dependency: Bump github.com/atc0005/go-nagios from 0.16.2 to 0.17.0
- (GH-1303) Go Dependency: Bump github.com/atc0005/go-nagios from 0.17.0 to 0.17.1
- (GH-1273) Go Dependency: Bump github.com/vmware/govmomi from 0.44.1 to 0.45.0
- (GH-1276) Go Dependency: Bump github.com/vmware/govmomi from 0.45.0 to 0.45.1
- (GH-1295) Go Dependency: Bump github.com/vmware/govmomi from 0.45.1 to 0.46.0
- (GH-1296) Go Dependency: Bump golang.org/x/sys from 0.26.0 to 0.27.0
- (GH-1291) Go Runtime: Bump golang from 1.22.8 to 1.22.9 in /dependabot/docker/go

## [v0.36.15] - 2024-10-17

### Changed

#### Dependency Updates

- (GH-1254) Build Image: Bump atc0005/go-ci from go-ci-oldstable-build-v0.21.13 to go-ci-oldstable-build-v0.21.14 in /dependabot/docker/builds
- (GH-1264) Go Dependency: Bump github.com/atc0005/go-nagios from 0.16.1 to 0.16.2
- (GH-1257) Go Dependency: Bump github.com/vmware/govmomi from 0.43.0 to 0.44.1
- (GH-1251) Go Dependency: Bump golang.org/x/sys from 0.25.0 to 0.26.0
- (GH-1250) Go Runtime: Bump golang from 1.22.7 to 1.22.8 in /dependabot/docker/go
- (GH-1260) Update Go version to 1.22

## [v0.36.14] - 2024-09-25

### Changed

#### Dependency Updates

- (GH-1241) Build Image: Bump atc0005/go-ci from go-ci-oldstable-build-v0.21.11 to go-ci-oldstable-build-v0.21.12 in /dependabot/docker/builds
- (GH-1246) Build Image: Bump atc0005/go-ci from go-ci-oldstable-build-v0.21.12 to go-ci-oldstable-build-v0.21.13 in /dependabot/docker/builds
- (GH-1233) Build Image: Bump atc0005/go-ci from go-ci-oldstable-build-v0.21.9 to go-ci-oldstable-build-v0.21.11 in /dependabot/docker/builds
- (GH-1243) Go Dependency: Bump github.com/vmware/govmomi from 0.42.0 to 0.43.0
- (GH-1238) Go Dependency: Bump golang.org/x/sys from 0.24.0 to 0.25.0
- (GH-1239) Go Runtime: Bump golang from 1.22.6 to 1.22.7 in /dependabot/docker/go

### Fixed

- (GH-1236) Fix gosec G115 linting errors

## [v0.36.13] - 2024-08-21

### Changed

#### Dependency Updates

- (GH-1220) Build Image: Bump atc0005/go-ci from go-ci-oldstable-build-v0.21.8 to go-ci-oldstable-build-v0.21.9 in /dependabot/docker/builds
- (GH-1214) Go Dependency: Bump github.com/vmware/govmomi from 0.39.0 to 0.40.0
- (GH-1224) Go Dependency: Bump github.com/vmware/govmomi from 0.40.0 to 0.42.0
- (GH-1223) Go Runtime: Bump golang from 1.21.13 to 1.22.6 in /dependabot/docker/go
- (GH-1222) Update project to Go 1.22 series

## [v0.36.12] - 2024-08-13

### Changed

#### Dependency Updates

- (GH-1190) Build Image: Bump atc0005/go-ci from go-ci-oldstable-build-v0.21.4 to go-ci-oldstable-build-v0.21.5 in /dependabot/docker/builds
- (GH-1192) Build Image: Bump atc0005/go-ci from go-ci-oldstable-build-v0.21.5 to go-ci-oldstable-build-v0.21.6 in /dependabot/docker/builds
- (GH-1201) Build Image: Bump atc0005/go-ci from go-ci-oldstable-build-v0.21.6 to go-ci-oldstable-build-v0.21.7 in /dependabot/docker/builds
- (GH-1208) Build Image: Bump atc0005/go-ci from go-ci-oldstable-build-v0.21.7 to go-ci-oldstable-build-v0.21.8 in /dependabot/docker/builds
- (GH-1195) Go Dependency: Bump github.com/vmware/govmomi from 0.38.0 to 0.39.0
- (GH-1205) Go Dependency: Bump golang.org/x/sys from 0.22.0 to 0.23.0
- (GH-1211) Go Dependency: Bump golang.org/x/sys from 0.23.0 to 0.24.0
- (GH-1207) Go Runtime: Bump golang from 1.21.12 to 1.21.13 in /dependabot/docker/go
- (GH-1198) Update Go version to 1.21

#### Other

- (GH-1202) Push `REPO_VERSION` var into containers for builds

## [v0.36.11] - 2024-07-10

### Changed

#### Dependency Updates

- (GH-1169) Build Image: Bump atc0005/go-ci from go-ci-oldstable-build-v0.20.7 to go-ci-oldstable-build-v0.20.8 in /dependabot/docker/builds
- (GH-1171) Build Image: Bump atc0005/go-ci from go-ci-oldstable-build-v0.20.8 to go-ci-oldstable-build-v0.21.0 in /dependabot/docker/builds
- (GH-1176) Build Image: Bump atc0005/go-ci from go-ci-oldstable-build-v0.21.0 to go-ci-oldstable-build-v0.21.2 in /dependabot/docker/builds
- (GH-1177) Build Image: Bump atc0005/go-ci from go-ci-oldstable-build-v0.21.2 to go-ci-oldstable-build-v0.21.3 in /dependabot/docker/builds
- (GH-1182) Build Image: Bump atc0005/go-ci from go-ci-oldstable-build-v0.21.3 to go-ci-oldstable-build-v0.21.4 in /dependabot/docker/builds
- (GH-1173) Go Dependency: Bump github.com/vmware/govmomi from 0.37.3 to 0.38.0
- (GH-1183) Go Dependency: Bump golang.org/x/sys from 0.21.0 to 0.22.0
- (GH-1179) Go Runtime: Bump golang from 1.21.11 to 1.21.12 in /dependabot/docker/go

## [v0.36.10] - 2024-06-06

### Changed

#### Dependency Updates

- (GH-1144) Build Image: Bump atc0005/go-ci from go-ci-oldstable-build-v0.20.4 to go-ci-oldstable-build-v0.20.5 in /dependabot/docker/builds
- (GH-1150) Build Image: Bump atc0005/go-ci from go-ci-oldstable-build-v0.20.5 to go-ci-oldstable-build-v0.20.6 in /dependabot/docker/builds
- (GH-1164) Build Image: Bump atc0005/go-ci from go-ci-oldstable-build-v0.20.6 to go-ci-oldstable-build-v0.20.7 in /dependabot/docker/builds
- (GH-1146) Go Dependency: Bump github.com/rs/zerolog from 1.32.0 to 1.33.0
- (GH-1157) Go Dependency: Bump github.com/vmware/govmomi from 0.37.2 to 0.37.3
- (GH-1158) Go Dependency: Bump golang.org/x/sys from 0.20.0 to 0.21.0
- (GH-1162) Go Runtime: Bump golang from 1.21.10 to 1.21.11 in /dependabot/docker/go

### Fixed

- (GH-1151) Remove inactive maligned linter
- (GH-1152) Fix errcheck linting errors

## [v0.36.9] - 2024-05-11

### Changed

#### Dependency Updates

- (GH-1130) Build Image: Bump atc0005/go-ci from go-ci-oldstable-build-v0.20.1 to go-ci-oldstable-build-v0.20.2 in /dependabot/docker/builds
- (GH-1135) Build Image: Bump atc0005/go-ci from go-ci-oldstable-build-v0.20.2 to go-ci-oldstable-build-v0.20.3 in /dependabot/docker/builds
- (GH-1137) Build Image: Bump atc0005/go-ci from go-ci-oldstable-build-v0.20.3 to go-ci-oldstable-build-v0.20.4 in /dependabot/docker/builds
- (GH-1124) Go Dependency: Bump github.com/vmware/govmomi from 0.36.3 to 0.37.0
- (GH-1125) Go Dependency: Bump github.com/vmware/govmomi from 0.37.0 to 0.37.1
- (GH-1139) Go Dependency: Bump github.com/vmware/govmomi from 0.37.1 to 0.37.2
- (GH-1128) Go Dependency: Bump golang.org/x/sys from 0.19.0 to 0.20.0
- (GH-1132) Go Runtime: Bump golang from 1.21.9 to 1.21.10 in /dependabot/docker/go

## [v0.36.8] - 2024-04-08

### Changed

#### Dependency Updates

- (GH-1100) Build Image: Bump atc0005/go-ci from go-ci-oldstable-build-v0.15.4 to go-ci-oldstable-build-v0.16.0 in /dependabot/docker/builds
- (GH-1102) Build Image: Bump atc0005/go-ci from go-ci-oldstable-build-v0.16.0 to go-ci-oldstable-build-v0.16.1 in /dependabot/docker/builds
- (GH-1105) Build Image: Bump atc0005/go-ci from go-ci-oldstable-build-v0.16.1 to go-ci-oldstable-build-v0.19.0 in /dependabot/docker/builds
- (GH-1109) Build Image: Bump atc0005/go-ci from go-ci-oldstable-build-v0.19.0 to go-ci-oldstable-build-v0.20.0 in /dependabot/docker/builds
- (GH-1119) Build Image: Bump atc0005/go-ci from go-ci-oldstable-build-v0.20.0 to go-ci-oldstable-build-v0.20.1 in /dependabot/docker/builds
- (GH-1095) Go Dependency: Bump github.com/vmware/govmomi from 0.35.0 to 0.36.0
- (GH-1098) Go Dependency: Bump github.com/vmware/govmomi from 0.36.0 to 0.36.1
- (GH-1106) Go Dependency: Bump github.com/vmware/govmomi from 0.36.1 to 0.36.2
- (GH-1111) Go Dependency: Bump github.com/vmware/govmomi from 0.36.2 to 0.36.3
- (GH-1117) Go Dependency: Bump golang.org/x/sys from 0.18.0 to 0.19.0
- (GH-1114) Go Runtime: Bump golang from 1.21.8 to 1.21.9 in /dependabot/docker/go

## [v0.36.7] - 2024-03-07

### Changed

#### Dependency Updates

- (GH-1092) Add todo/release label to "Go Runtime" PRs
- (GH-1083) Build Image: Bump atc0005/go-ci from go-ci-oldstable-build-v0.15.2 to go-ci-oldstable-build-v0.15.3 in /dependabot/docker/builds
- (GH-1089) Build Image: Bump atc0005/go-ci from go-ci-oldstable-build-v0.15.3 to go-ci-oldstable-build-v0.15.4 in /dependabot/docker/builds
- (GH-1078) canary: bump golang from 1.21.6 to 1.21.7 in /dependabot/docker/go
- (GH-1075) docker: bump atc0005/go-ci from go-ci-oldstable-build-v0.15.0 to go-ci-oldstable-build-v0.15.2 in /dependabot/docker/builds
- (GH-1086) Go Dependency: Bump golang.org/x/sys from 0.17.0 to 0.18.0
- (GH-1087) Go Runtime: Bump golang from 1.21.7 to 1.21.8 in /dependabot/docker/go
- (GH-1080) Update Dependabot PR prefixes (redux)
- (GH-1079) Update Dependabot PR prefixes
- (GH-1077) Update project to Go 1.21 series

## [v0.36.6] - 2024-02-15

### Changed

#### Dependency Updates

- (GH-1056) canary: bump golang from 1.20.13 to 1.20.14 in /dependabot/docker/go
- (GH-1037) docker: bump atc0005/go-ci from go-ci-oldstable-build-v0.14.3 to go-ci-oldstable-build-v0.14.4 in /dependabot/docker/builds
- (GH-1045) docker: bump atc0005/go-ci from go-ci-oldstable-build-v0.14.4 to go-ci-oldstable-build-v0.14.5 in /dependabot/docker/builds
- (GH-1048) docker: bump atc0005/go-ci from go-ci-oldstable-build-v0.14.5 to go-ci-oldstable-build-v0.14.6 in /dependabot/docker/builds
- (GH-1053) docker: bump atc0005/go-ci from go-ci-oldstable-build-v0.14.6 to go-ci-oldstable-build-v0.14.8 in /dependabot/docker/builds
- (GH-1059) docker: bump atc0005/go-ci from go-ci-oldstable-build-v0.14.8 to go-ci-oldstable-build-v0.14.9 in /dependabot/docker/builds
- (GH-1062) docker: bump atc0005/go-ci from go-ci-oldstable-build-v0.14.9 to go-ci-oldstable-build-v0.15.0 in /dependabot/docker/builds
- (GH-1043) go.mod: bump github.com/atc0005/go-nagios from 0.16.0 to 0.16.1
- (GH-1049) go.mod: bump github.com/rs/zerolog from 1.31.0 to 1.32.0
- (GH-1065) go.mod: bump github.com/vmware/govmomi from 0.34.2 to 0.35.0
- (GH-1057) go.mod: bump golang.org/x/sys from 0.16.0 to 0.17.0

### Fixed

- (GH-1067) Replace property.Filter with property.Match
- (GH-1071) Fix `unused-parameter` revive linting errors

## [v0.36.5] - 2024-01-19

### Changed

#### Dependency Updates

- (GH-1032) canary: bump golang from 1.20.12 to 1.20.13 in /dependabot/docker/go
- (GH-1034) docker: bump atc0005/go-ci from go-ci-oldstable-build-v0.14.2 to go-ci-oldstable-build-v0.14.3 in /dependabot/docker/builds
- (GH-1023) ghaw: bump github/codeql-action from 2 to 3
- (GH-1021) go.mod: bump github.com/vmware/govmomi from 0.33.1 to 0.34.0
- (GH-1025) go.mod: bump github.com/vmware/govmomi from 0.34.0 to 0.34.1
- (GH-1029) go.mod: bump github.com/vmware/govmomi from 0.34.1 to 0.34.2
- (GH-1027) go.mod: bump golang.org/x/sys from 0.15.0 to 0.16.0

## [v0.36.4] - 2023-12-08

### Changed

#### Dependency Updates

- (GH-1012) canary: bump golang from 1.20.11 to 1.20.12 in /dependabot/docker/go
- (GH-1015) docker: bump atc0005/go-ci from go-ci-oldstable-build-v0.14.1 to go-ci-oldstable-build-v0.14.2 in /dependabot/docker/builds
- (GH-1011) go.mod: bump golang.org/x/sys from 0.14.0 to 0.15.0

## [v0.36.3] - 2023-11-15

### Changed

#### Dependency Updates

- (GH-1002) canary: bump golang from 1.20.10 to 1.20.11 in /dependabot/docker/go
- (GH-1004) docker: bump atc0005/go-ci from go-ci-oldstable-build-v0.14.0 to go-ci-oldstable-build-v0.14.1 in /dependabot/docker/builds
- (GH-1000) go.mod: bump golang.org/x/sys from 0.13.0 to 0.14.0

## [v0.36.2] - 2023-11-02

### Changed

#### Dependency Updates

- (GH-980) canary: bump golang from 1.20.8 to 1.20.10 in /dependabot/docker/go
- (GH-977) docker: bump atc0005/go-ci from go-ci-oldstable-build-v0.13.10 to go-ci-oldstable-build-v0.13.11 in /dependabot/docker/builds
- (GH-981) docker: bump atc0005/go-ci from go-ci-oldstable-build-v0.13.11 to go-ci-oldstable-build-v0.13.12 in /dependabot/docker/builds
- (GH-989) docker: bump atc0005/go-ci from go-ci-oldstable-build-v0.13.12 to go-ci-oldstable-build-v0.14.0 in /dependabot/docker/builds
- (GH-976) go.mod: bump github.com/google/go-cmp from 0.5.9 to 0.6.0
- (GH-985) go.mod: bump github.com/mattn/go-isatty from 0.0.19 to 0.0.20
- (GH-965) go.mod: bump github.com/rs/zerolog from 1.30.0 to 1.31.0
- (GH-967) go.mod: bump github.com/vmware/govmomi from 0.30.7 to 0.32.0
- (GH-987) go.mod: bump github.com/vmware/govmomi from 0.32.0 to 0.33.0
- (GH-993) go.mod: bump github.com/vmware/govmomi from 0.33.0 to 0.33.1
- (GH-969) go.mod: bump golang.org/x/sys from 0.12.0 to 0.13.0

### Fixed

- (GH-996) Fix goconst linting errors

## [v0.36.1] - 2023-10-06

### Changed

#### Dependency Updates

- (GH-947) canary: bump golang from 1.20.7 to 1.20.8 in /dependabot/docker/go
- (GH-949) docker: bump atc0005/go-ci from go-ci-oldstable-build-v0.13.7 to go-ci-oldstable-build-v0.13.8 in /dependabot/docker/builds
- (GH-956) docker: bump atc0005/go-ci from go-ci-oldstable-build-v0.13.8 to go-ci-oldstable-build-v0.13.9 in /dependabot/docker/builds
- (GH-959) docker: bump atc0005/go-ci from go-ci-oldstable-build-v0.13.9 to go-ci-oldstable-build-v0.13.10 in /dependabot/docker/builds
- (GH-945) ghaw: bump actions/checkout from 3 to 4
- (GH-943) go.mod: bump golang.org/x/sys from 0.11.0 to 0.12.0

## [v0.36.0] - 2023-08-31

### Added

- (GH-927) Update datastore plugins to indicate which VMs within datastore are
  templates

### Changed

- (GH-926) Refactor objects tally logic used to provide plugin "trailer"
  summary details
- (GH-932) Update `vsphere.VMwareAdminAssistanceNeeded` annotation to point
  sysadmins to plugin doc
- (GH-939) Update datastore performance plugin doc

## [v0.35.1] - 2023-08-25

### Changed

- Dependencies
  - `atc0005/go-ci`
    - `go-ci-oldstable-build-v0.13.4` to `go-ci-oldstable-build-v0.13.7`

### Fixed

- (GH-916) `gosec` `G601: Implicit memory aliasing in for loop` linting errors

## [v0.35.0] - 2023-08-17

### Added

- (GH-809) Add support for excluding/ignoring VMs by `Folder` Managed Object
  ID (MOID)

### Changed

- Dependencies
  - `Go`
    - `1.19.12` to `1.20.7`
  - `atc0005/go-ci`
    - `go-ci-oldstable-build-v0.13.2` to `go-ci-oldstable-build-v0.13.4`
- (GH-862) Replace hard-coded flag names with constants
- (GH-901) Update project to Go 1.20 series

### Fixed

- (GH-658) Potential nil pointer dereference in
  `vsphere.ResourcePoolsMemoryReport()`
- (GH-852) Evaluate consistency of terminology regarding VM collections
- (GH-897) README: Fix short flag for virtual hardware plugin

## [v0.34.1] - 2023-08-10

### Changed

- Dependencies
  - `Go`
    - `1.19.11` to `1.19.12`
  - `vmware/govmomi`
    - `v0.30.6` to `v0.30.7`
  - `atc0005/go-ci`
    - `go-ci-oldstable-build-v0.13.1` to `go-ci-oldstable-build-v0.13.2`
  - `rs/zerolog`
    `v1.29.1` to `v1.30.0`
  - `golang.org/x/sys`
    - `v0.10.0` to `v0.11.0`
- (GH-868) Update help text for include/exclude RP flags to emphasize
  filtering by name
- (GH-869) Add VM power state tally helper funcs
- (GH-890) Replace fixed MO type strings with constants

### Fixed

- (GH-887) Add missing resource pools flags validation
- (GH-888) Fix func names in deferred debug log messages
- (GH-889) Use correct mo type in error message

## [v0.34.0] - 2023-07-28

### Added

- (GH-839) Add initial automated release notes config
- (GH-841) Add initial automated release build workflow
- (GH-851) Create plugin to list Virtual Machines in order to test
  include/exclude options

### Changed

- Dependencies
  - `vmware/govmomi`
    - `v0.30.5` to `v0.30.6`
  - `atc0005/go-ci`
    - `go-ci-oldstable-build-v0.11.3` to `go-ci-oldstable-build-v0.13.1`
- (GH-844) Update Dependabot config to monitor both branches
- (GH-853) Update .gitignore patterns

### Fixed

- (GH-842) Update CodeQL GHAW timeout

## [v0.33.3] - 2023-07-13

### Overview

- RPM package improvements
- Change exit state for several scenarios
- Bug fixes
- Dependency updates
- built using Go 1.19.11
  - Statically linked
  - Linux x64

### Changed

- Dependencies
  - `Go`
    - `1.19.10` to `1.19.11`
  - `vmware/govmomi`
    - `v0.30.4` to `v0.30.5`
  - `atc0005/go-nagios`
    - `v0.15.0` to `v0.16.0`
  - `atc0005/go-ci`
    - `go-ci-oldstable-build-v0.11.0` to `go-ci-oldstable-build-v0.11.3`
  - `golang.org/x/sys`
    - `v0.9.0` to `v0.10.0`
- (GH-829) Update RPM postinstall scripts to use restorecon
- (GH-825) Update error annotation implementation

### Fixed

- (GH-828) Use UNKNOWN state for invalid command-line args
- (GH-832) Use UNKNOWN state for perfdata add failure

## [v0.33.2] - 2023-06-21

### Overview

- Bug fixes
- Dependency updates
- built using Go 1.19.10
  - Statically linked
  - Linux x64

### Changed

- Dependencies
  - `atc0005/go-ci`
    - `go-ci-oldstable-build-v0.10.6` to `go-ci-oldstable-build-v0.11.0`
  - `golang.org/x/sys`
    - `v0.8.0` to `v0.9.0`

### Fixed

- (GH-818) Restore local CodeQL workflow
- (GH-820) Fix helper function closure collection evaluation

## [v0.33.1] - 2023-06-09

### Overview

- Bug fixes
- Dependency updates
- built using Go 1.19.10
  - Statically linked
  - Linux x64

### Changed

- Dependencies
  - `Go`
    - `1.19.9` to `1.19.10`
  - `atc0005/go-nagios`
    - `v0.14.0` to `v0.15.0`
  - `atc0005/go-ci`
    - `go-ci-oldstable-build-v0.10.5` to `go-ci-oldstable-build-v0.10.6`
  - `mattn/go-isatty`
    - `v0.0.18` to `v0.0.19`
- (GH-814) Update vuln analysis GHAW to remove on.push hook

### Fixed

- (GH-811) Disable depguard linter

## [v0.33.0] - 2023-05-12

### Overview

- Bug fixes
- Dependency updates
- built using Go 1.19.9
  - Statically linked
  - Linux x64

### Changed

- Dependencies
  - `Go`
    - `1.19.7` to `1.19.9`
  - `atc0005/go-ci`
    - `go-ci-oldstable-build-v0.10.3` to `go-ci-oldstable-build-v0.10.5`
  - `rs/zerolog`
    - `v1.29.0` to `v1.29.1`
  - `golang.org/x/sys`
    - `v0.6.0` to `v0.8.0`
- (GH-796) Add .dockerignore file for use during image builds

### Fixed

- (GH-805) Fix revive linter errors

## [v0.32.0] - 2023-03-31

### Overview

- Build improvements
- Bug fixes
- Dependency updates
- built using Go 1.19.7
  - Statically linked
  - Linux x64

### Added

- (GH-793) Add rootless container builds via Docker/Podman

### Changed

- Dependencies
  - `Go`
    - `1.19.6` to `1.19.7`
  - `vmware/govmomi`
    - `v0.30.2` to `v0.30.4`
  - `mattn/go-isatty`
    - `v0.0.17` to `v0.0.18`
  - `atc0005/go-ci`
    - `go-ci-oldstable-build-v0.9.0` to `go-ci-oldstable-build-v0.10.3`
- (GH-781) Update .gitignore to exclude Windows syso files
- (GH-787) Update vuln analysis GHAW to use on.push hook

### Fixed

- (GH-780) Fix '*-all-links.txt' generation
- (GH-782) RPM installation output is "saw-toothed" and "noisy"

## [v0.31.0] - 2023-03-05

### Overview

- Add support for generating packages
- Generated binary changes
  - filename patterns
  - compression
  - executable metadata
- Build improvements
- Dependency updates
- built using Go 1.19.6
  - Statically linked
  - Linux x64

### Added

- (GH-771) Generate RPM/DEB packages using nFPM

### Changed

- Dependencies
  - `golang.org/x/sys`
    - `v0.5.0` to `v0.6.0`
- Build process
  - (GH-770) Switch to semantic versioning (semver) compatible versioning
    pattern
  - (GH-772) Add version metadata to Windows executables
  - (GH-773) Makefile: Compress binaries and use fixed filenames
  - (GH-774) Makefile: Refresh recipes to add "standard" set, new
    package-related options
  - (GH-775) Build dev/stable releases using go-ci Docker image

## [v0.30.8] - 2023-03-02

### Overview

- Bug fixes
- Dependency updates
- GitHub Actions Workflows updates
- built using Go 1.19.6
  - Statically linked
  - Linux x64

### Changed

- Dependencies
  - `Go`
    - `1.19.4` to `1.19.6`
  - `atc0005/go-nagios`
    - `v0.10.2` to `v0.14.0`
  - `rs/zerolog`
    - `v1.28.0` to `v1.29.0`
  - `vmware/govmomi`
    - `v0.29.0` to `v0.30.2`
  - `mattn/go-isatty`
    - `v0.0.16` to `v0.0.17`
  - `golang.org/x/sys`
    - `v0.3.0` to `v0.5.0`
- (GH-750) Drop plugin runtime tracking, update library usage
  - `time` metric is provided "automatically" via library
- (GH-755) Add Go Module Validation, Dependency Updates jobs
- (GH-763) Drop `Push Validation` workflow
- (GH-764) Rework workflow scheduling
- (GH-766) Remove `Push Validation` workflow status badge

### Fixed

- (GH-749) library logging is not enabled at Trace level

## [v0.30.7] - 2022-12-07

### Overview

- Bug fixes
- Dependency updates
- GitHub Actions Workflows updates
- built using Go 1.19.4
  - Statically linked
  - Linux x64

### Changed

- Dependencies
  - `Go`
    - `1.19.1` to `1.19.4`
  - `atc0005/go-nagios`
    - `v0.9.2` to `v0.10.2`
  - `github.com/mattn/go-colorable`
    - `v0.1.12` to `v0.1.13`
  - `github.com/mattn/go-isatty`
    - `v0.0.14` to `v0.0.16`
  - `golang.org/x/sys`
    - `v0.0.0-20210927094055-39ccf1dd6fa6` to `v0.3.0`
- (GH-734) Refactor GitHub Actions workflows to import logic

### Fixed

- (GH-739) Move perfdata debug message to correct location
- (GH-740) Fix Makefile Go module base path detection

## [v0.30.6] - 2022-09-15

### Overview

- Bug fixes
- Dependency updates
- built using Go 1.19.1
  - Statically linked
  - Linux x64

### Changed

- Dependencies
  - `atc0005/go-nagios`
    - `v0.9.1` to `v0.9.2`

### Fixed

- (GH-724) Timeout occurred for `Build codebase using Makefile all recipe`
  GHAW job
- (GH-725) check_vmware_datastore_space plugin Service Ouput contains invalid
  `MISSING` text

## [v0.30.5] - 2022-09-14

### Overview

- Bug fixes
- Dependency updates
- GitHub Actions Workflows updates
- built using Go 1.19.1
  - Statically linked
  - Linux x64

### Added

- (GH-720) Add Vulnerability Analysis GitHub Actions Workflow

### Changed

- Dependencies
  - `Go`
    - `1.17.13` to `1.19.1`
  - `rs/zerolog`
    - `v1.27.0` to `v1.28.0`
  - `google/go-cmp`
    - `v0.5.8` to `v0.5.9`
  - `github/codeql-action`
    - `v2.1.18` to `v2.1.22`

- (GH-710) Update project to Go 1.19
- (GH-712) Update Makefile and GitHub Actions Workflows
- (GH-719) Update check_vmware_hs2ds2vms plugin to list all hosts and
  datastores missing Custom Attributes

### Fixed

- (GH-711) Adjust doc comments formatting for doc.go
- (GH-714) vsphere.GetObjectCustomAttribute func (incorrectly) uses fixed
  object type in log/error messages

## [v0.30.4] - 2022-08-16

### Overview

- Dependency updates
- built using Go 1.17.13
  - Statically linked
  - Linux x64

### Changed

- Dependencies
  - `Go`
    - `1.17.12` to `1.17.13`
  - `github/codeql-action`
    - `v2.1.15` to `v2.1.18`

- (GH-702) Emit markdownlint CLI version after installation

### Fixed

- (GH-700) Markdown linting failures: invalid link fragments, unused
  link/image references
- (GH-703) Apply linting fixes for Go 1.19 release

## [v0.30.3] - 2022-07-13

### Overview

- Dependency updates
- built using Go 1.17.12
  - Statically linked
  - Linux x64

### Changed

- Dependencies
  - `Go`
    - `1.17.10` to `1.17.12`
  - `vmware/govmomi`
    - `v0.28.0` to `v0.29.0`
  - `rs/zerolog`
    - `v1.26.1` to `v1.27.0`
  - `atc0005/go-nagios`
    - `v0.8.2` to `v0.9.1`
  - `github/codeql-action`
    - `v2.1.10` to `v2.1.15`

### Fixed

- (GH-677) Tweak doc comments for check_vmware_alarms plugin
- (GH-679) Ampersands replaced by Nagios in plugin output
- (GH-689) Fix misc typo in doc comment

## [v0.30.2] - 2022-05-11

### Overview

- Dependency updates
- built using Go 1.17.10
  - Statically linked
  - Linux x64

### Changed

- Dependencies
  - `Go`
    - `1.17.9` to `1.17.10`
  - `github/codeql-action`
    - `v2.1.9` to `v2.1.10`

## [v0.30.1] - 2022-04-29

### Overview

- Bugfixes
- Dependency updates
- built using Go 1.17.9
  - Statically linked
  - Linux x64

### Changed

- Dependencies
  - `vmware/govmomi`
    - `v0.27.4` to `v0.28.0`
  - `github/codeql-action`
    - `v2.1.8` to `v2.1.9`

### Fixed

- (GH-668) Total memory calculation error in the `check_vmware_rps_memory`
  plugin

## [v0.30.0] - 2022-04-27

### Overview

- Expand support for Nagios Performance Data
- Bugfixes
- Dependency updates
- built using Go 1.17.9
  - Statically linked
  - Linux x64

### Added

- (GH-652) Extend `check_vmware_rps_memory` plugin to emit swapped memory
  perfdata
- (GH-644) Extend `check_vmware_rps_memory` plugin to emit ballooned memory
  perfdata

### Changed

- Dependencies
  - `Go`
    - `1.17.8` to `1.17.9`
  - `google/go-cmp`
    - `v0.5.7` to `v0.5.8`
  - `github/codeql-action`
    - `v1.1.5` to `v2.1.8`

- (GH-648) Update send2teams command definition to reflect changes from
  send2teams v0.9.0 release
- (GH-650) Update send2teams command definition to use condensed overview
  format

### Fixed

- (GH-653) Fix duplication in variable name
- (GH-659) Update doc for check_vmware_rps_memory plugin
- (GH-660) Expand logging, attempt state reload for Resource Pool stats
- (GH-663) Document that RPS plugin requires vCenter instance

## [v0.29.2] - 2022-03-16

### Overview

- Report output tweaks
- Bugfixes
- Dependency updates
- built using Go 1.17.8
  - Statically linked
  - Linux x64

### Changed

- Dependencies
  - `Go`
    - `1.17.7` to `1.17.8`
  - `github/codeql-action`
    - `v1.0.32` to `v1.1.5`
  - `actions/checkout`
    - `v2.4.0` to `v3`
  - `actions/setup-node`
    - `v2.5.1` to `v3`

- (GH-635) Expose custom attribute used as key for mismatched
  host/datastore/vm pairings in report output for `check_vmware_hs2ds2vms`
  plugin

- (GH-640) Add `FilterDatastoresByIDs()` func

### Fixed

- (GH-636) Incorrect vSphere object type mentioned in error message

## [v0.29.1] - 2022-02-11

### Overview

- Bugfixes
- Dependency updates
- built using Go 1.17.7
  - Statically linked
  - Linux x64

### Changed

- Dependencies
  - `Go`
    - `1.17.6` to `1.17.7`
  - `github/codeql-action`
    - `v1.0.28` to `v1.0.32`
  - `vmware/govmomi`
    - `v0.27.2` to `v0.27.4`

- (GH-611) Enable additional govet linter analyzers
- (GH-613) Various linting issues detected by nilness, shadow `govet`
  analyzers
- (GH-617) Expand linting GitHub Actions Workflow to include `oldstable`,
  `unstable` container images
- (GH-618) Switch Docker image source from Docker Hub to GitHub Container
  Registry (GHCR)

### Fixed

- CHANGELOG
  - v0.29.0 release incorrectly noted latest codeql-action release

## [v0.29.0] - 2022-01-23

### Overview

- New plugin
- Bugfixes
- Dependency updates
- built using Go 1.17.6
  - Statically linked
  - Linux x64

### Added

- (GH-506) New plugin: `check_vmware_vm_backup_via_ca`

### Changed

- Dependencies
  - `Go`
    - `1.17.5` to `1.17.6`
  - `github/codeql-action`
    - `v1.0.26` to `v1.0.28`
  - `google/go-cmp`
    - `v0.5.6` to `v0.5.7`

- (GH-580) Add project name to generated download links file

### Fixed

- (GH-581) Review & update timeout-minutes setting (if needed) for all GitHub
  Actions Workflows
- (GH-600) Bug in resource pool exclusion logic
- (GH-603) Incorrect power state evaluation noted in logger field for
  `check_vmware_vm_power_uptime` plugin

## [v0.28.1] - 2022-01-01

### Overview

- Dependency updates
- built using Go 1.17.5
  - Statically linked
  - Linux x64

### Changed

- Dependencies
  - `atc0005/go-nagios`
    - `v0.8.1` to `v0.8.2`
  - `actions/setup-node`
    - `v2.5.0` to `v2.5.1`

### Fixed

- (GH-575) CHANGELOG | Correct builds list for release entries

## [v0.28.0] - 2021-12-30

### Overview

- Dependency updates
- Breaking changes
  - the `check_vmware_datastore` plugin has been renamed
  - performance data metrics have been renamed
- built using Go 1.17.5
  - Statically linked
  - Linux x64

### Changed

- Dependencies
  - `Go`
    - (GH-563) Update go.mod file, canary Dockerfile to reflect current
      dependencies
    - `1.16.12` to `1.17.5`

- **Breaking**
  - (GH-510) `check_vmware_datastore` plugin
    - renamed from `check_vmware_datastore` to `check_vmware_datastore_space`
  - (GH-480) `*_used` and `*_remaining` metrics
    - renamed metrics
      - `datastore_usage` to `datastore_space_usage`
      - `datastore_storage_used` to `datastore_space_used`
      - `datastore_storage_remaining` to `datastore_space_remaining`
    - updated plugins
      - `check_vmware_datastore_space`
      - `check_vmware_host_cpu`
      - `check_vmware_host_memory`
      - `check_vmware_rps_memory`
      - `check_vmware_vcpus`

- (GH-533) Update datastore usage/space plugin to evaluate whether datastore
  is accessible

- (GH-539) Split README into separate files

### Fixed

- (GH-549) Affirm purpose of user domain flag
- (GH-558) The `vphere.GetTriggeredAlarms()` method fails to consider a
  `VirtualApp` entity as a `ResourcePool`

## [v0.27.1] - 2021-12-28

### Overview

- Bugfixes
- Dependency updates
- Deprecated `check_vmware_datastore` plugin name
  - will be renamed in the `v0.28.0` release
- built using Go 1.16.12
  - Statically linked
  - Linux x64

### Changed

- Dependencies
  - `Go`
    - `1.16.10` to `1.16.12`
  - `github/codeql-action`
    - `v1.0.24` to `v1.0.26`
  - `rs/zerolog`
    - `v1.26.0` to `v1.26.1`

- (GH-555) Help output generated by `-h`, `--help` flag is sent to `stderr`,
  probably should go to `stdout` instead

### Deprecated

- (GH-510, GH-530) `check_vmware_datastore` plugin
  - this plugin will be renamed in the `v0.28.0` release to
    `check_vmware_datastore_space`
  - documentation, `Makefile` and other changes will be applied in the
    `v0.28.0` release to accommodate the rename of this plugin

### Fixed

- CHANGELOG
  - `github/codeql-action` version in v0.27.0 release
- (GH-548) Add missing shorthand indicators to short flags
- (GH-551) `config.MultiValueDSPerfLatencyMetricFlag` type incorrectly named

## [v0.27.0] - 2021-12-01

### Overview

- New plugin
- Bugfixes
- Dependency updates
- Deprecated `check_vmware_datastore` plugin name
  - will be renamed in the `v0.28.0` release
- built using Go 1.16.10
  - Statically linked
  - Linux x64

### Added

- (GH-505) New plugin: `check_vmware_datastore_performance`

### Changed

- Dependencies
  - `Go`
    - drop Go version in `go.mod` from `1.15` to `1.14`
      - attempt to reflect actual base Go version required by dependencies
  - `github/codeql-action`
    - `v1.0.22` to `v1.0.24`
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
  - Linux x64

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
  - Linux x64

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
  - Linux x64

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
  - Linux x64

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
  - Linux x64

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
  - Linux x64

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
  - Linux x64

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

[Unreleased]: https://github.com/atc0005/check-vmware/compare/v0.36.16...HEAD
[v0.36.16]: https://github.com/atc0005/check-vmware/releases/tag/v0.36.16
[v0.36.15]: https://github.com/atc0005/check-vmware/releases/tag/v0.36.15
[v0.36.14]: https://github.com/atc0005/check-vmware/releases/tag/v0.36.14
[v0.36.13]: https://github.com/atc0005/check-vmware/releases/tag/v0.36.13
[v0.36.12]: https://github.com/atc0005/check-vmware/releases/tag/v0.36.12
[v0.36.11]: https://github.com/atc0005/check-vmware/releases/tag/v0.36.11
[v0.36.10]: https://github.com/atc0005/check-vmware/releases/tag/v0.36.10
[v0.36.9]: https://github.com/atc0005/check-vmware/releases/tag/v0.36.9
[v0.36.8]: https://github.com/atc0005/check-vmware/releases/tag/v0.36.8
[v0.36.7]: https://github.com/atc0005/check-vmware/releases/tag/v0.36.7
[v0.36.6]: https://github.com/atc0005/check-vmware/releases/tag/v0.36.6
[v0.36.5]: https://github.com/atc0005/check-vmware/releases/tag/v0.36.5
[v0.36.4]: https://github.com/atc0005/check-vmware/releases/tag/v0.36.4
[v0.36.3]: https://github.com/atc0005/check-vmware/releases/tag/v0.36.3
[v0.36.2]: https://github.com/atc0005/check-vmware/releases/tag/v0.36.2
[v0.36.1]: https://github.com/atc0005/check-vmware/releases/tag/v0.36.1
[v0.36.0]: https://github.com/atc0005/check-vmware/releases/tag/v0.36.0
[v0.35.1]: https://github.com/atc0005/check-vmware/releases/tag/v0.35.1
[v0.35.0]: https://github.com/atc0005/check-vmware/releases/tag/v0.35.0
[v0.34.1]: https://github.com/atc0005/check-vmware/releases/tag/v0.34.1
[v0.34.0]: https://github.com/atc0005/check-vmware/releases/tag/v0.34.0
[v0.33.3]: https://github.com/atc0005/check-vmware/releases/tag/v0.33.3
[v0.33.2]: https://github.com/atc0005/check-vmware/releases/tag/v0.33.2
[v0.33.1]: https://github.com/atc0005/check-vmware/releases/tag/v0.33.1
[v0.33.0]: https://github.com/atc0005/check-vmware/releases/tag/v0.33.0
[v0.32.0]: https://github.com/atc0005/check-vmware/releases/tag/v0.32.0
[v0.31.0]: https://github.com/atc0005/check-vmware/releases/tag/v0.31.0
[v0.30.8]: https://github.com/atc0005/check-vmware/releases/tag/v0.30.8
[v0.30.7]: https://github.com/atc0005/check-vmware/releases/tag/v0.30.7
[v0.30.6]: https://github.com/atc0005/check-vmware/releases/tag/v0.30.6
[v0.30.5]: https://github.com/atc0005/check-vmware/releases/tag/v0.30.5
[v0.30.4]: https://github.com/atc0005/check-vmware/releases/tag/v0.30.4
[v0.30.3]: https://github.com/atc0005/check-vmware/releases/tag/v0.30.3
[v0.30.2]: https://github.com/atc0005/check-vmware/releases/tag/v0.30.2
[v0.30.1]: https://github.com/atc0005/check-vmware/releases/tag/v0.30.1
[v0.30.0]: https://github.com/atc0005/check-vmware/releases/tag/v0.30.0
[v0.29.2]: https://github.com/atc0005/check-vmware/releases/tag/v0.29.2
[v0.29.1]: https://github.com/atc0005/check-vmware/releases/tag/v0.29.1
[v0.29.0]: https://github.com/atc0005/check-vmware/releases/tag/v0.29.0
[v0.28.1]: https://github.com/atc0005/check-vmware/releases/tag/v0.28.1
[v0.28.0]: https://github.com/atc0005/check-vmware/releases/tag/v0.28.0
[v0.27.1]: https://github.com/atc0005/check-vmware/releases/tag/v0.27.1
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
