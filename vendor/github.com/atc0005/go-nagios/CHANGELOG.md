# Changelog

## Overview

All notable changes to this project will be documented in this file.

The format is based on [Keep a
Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to
[Semantic Versioning](https://semver.org/spec/v2.0.0.html).

Please [open an issue](https://github.com/atc0005/go-nagios/issues) for any
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

## [v0.11.0] - 2022-12-14

### Added

- (GH-102) Add constructor to initialize plugin start time, provide `time`
  Performance Data metric "automatically"

## [v0.10.2] - 2022-10-19

### Fixed

- (GH-164) Sort performance data metrics for consistent output

## [v0.10.1] - 2022-10-09

### Changed

- Dependencies
  - `google/go-cmp`
    - `v0.5.8` to `v0.5.9`
- (GH-103) Emit Performance Data on same line as `ServiceOutput` if
  `LongServiceOutput` is empty
- (GH-146) Move examples from README to Example tests
- (GH-154) Refactor GitHub Actions workflows to import logic
- (GH-156) Update README to include go.mod badge
- (GH-157) Prevent duplication of Performance Data metrics

## [v0.10.0] - 2022-09-18

### Added

- (GH-144) Add test to prevent further Service Output interpolation
  regressions
- (GH-145) Add support for setting preferred output target
- (GH-147) Add option to skip os.Exit call at end of plugin execution
- (GH-149) Add initial plugin output test

### Changed

- Dependencies
  - `github/codeql-action`
    - `v2.1.22` to `v2.1.23`

## [v0.9.2] - 2022-09-15

### Changed

- (GH-137) Update project to Go 1.19
- (GH-138) Update Makefile and GitHub Actions Workflows

### Fixed

- (GH-139) `ExitState.ServiceOutput` is (incorrectly) interpreted as
  containing format specifier

## [v0.9.1] - 2022-06-20

### Fixed

- (GH-132) `LongServiceOutput` header is emitted when
  `ExitState.LongServiceOutput` field is empty
- (GH-134) Update lintinstall Makefile recipe
- (GH-135) Fix Go 1.19beta1 gofmt linter warning

## [v0.9.0] - 2022-06-16

### Added

- (GH-107) Add support for overriding section headers/labels
- (GH-117) Add options to explicitly "hide" or omit Thresholds and Errors
  section, automatically hide Detailed Info section if those sections are
  unused
- (GH-120) Allow collecting multiple errors

### Changed

- Dependencies
  - `actions/checkout`
    - `v2.4.0` to `v3`
  - `actions/setup-node`
    - `v2.5.1` to `v3`

- (GH-118) Automatically omit Thresholds and Errors sections if unused
- (GH-112) Expand linting GitHub Actions Workflow to include `oldstable`,
  `unstable` container images
- (GH-113) Switch Docker image source from Docker Hub to GitHub Container
  Registry (GHCR)
- (GH-130) Refresh README, go.doc files for v0.9.0 release

### Fixed

- (GH-121) `err113` linting issues: `do not define dynamic errors, use wrapped
  static errors instead`
- (GH-127) cognitive complexity 32 of func `(*ExitState).ReturnCheckResults`
  is high (> 10) (gocognit)

## [v0.8.2] - 2022-01-01

### Changed

- Dependencies
  - `actions/checkout`
    - `v2.3.5` to `v2.4.0`
  - `actions/setup-node`
    - `v2.4.1` to `v2.5.1`

### Fixed

- (GH-104) Swap nil PerformanceData check for length check
- (GH-109) Nagios XI interprets current `nagios.CheckOutputEOL` value as two
  newlines

## [v0.8.1] - 2021-11-02

### Changed

- Dependencies
  - `actions/checkout`
    - `v2.3.4` to `v2.3.5`

### Fixed

- (GH-98) Performance Data emitted when `ServiceOutput` and
  `LongServiceOutput` empty
- (GH-96) CHANGELOG | Fix typo and clarify perfdata support

## [v0.8.0] - 2021-09-30

### Added

- Add initial support for collecting, formatting & emitting performance data

### Changed

- Dependencies
  - `actions/setup-node`
    - `v2.4.0` to `v2.4.1`

### Fixed

- Fix typo in CHANGELOG release tags list
- Remove duplicate package doc comments

## [v0.7.0] - 2021-09-07

### Added

- (GH-80) Add `ServiceState` type

### Changed

- Dependencies
  - `actions/setup-node`
    - (`v2.2.0` to `v2.4.0`

### Fixed

- (GH-76) Fix typo in doc comment for `ExitStatusCode` field
- (GH-82) Typo in `UNKNOWN` state label text

## [v0.6.1] - 2021-07-15

### Changed

- Panic error message
  - (GH-69) stack trace wrapped with Markdown code fences instead of `pre`
    tags
  - (GH-68) provide complete stack trace instead of snippet
- Documentation
  - Replace GoDoc badge with pkg.go.dev badge
- Dependencies
  - `actions/setup-node`
    - (`v2.1.4` to `v2.2.0`
    - update `node-version` value to always use latest LTS version instead of
      hard-coded version

## [v0.6.0] - 2021-01-14

### Added

- Expose panic error message
  - emit as first line item in pre block above stack trace
  - emit as part of error entry

## [v0.5.3] - 2021-01-05

### Changed

- dependencies
  - `actions/setup-node`
    - `v2.1.2` to `v2.1.4`

### Fixed

- Preformatted `ServiceOutput` string subjected to another (failed) formatting
  operation
- Breadcrumb URL formatting

## [v0.5.2] - 2020-11-08

### Changed

- dependencies
  - `actions/setup-node`
    - `v2.1.1` to `v2.1.2`
  - `actions/checkout`
    - `v2.3.2` to `v2.3.4`

### Fixed

- (*nagios.ExitState).ReturnCheckResults() unintentionally masks or "swallows"
  panics

## [v0.5.1] - 2020-09-22

### Changed

- `ExitState` receiver type is now a pointer for `ReturnCheckResults()` method

### Fixed

- Documentation
  - Update doc comments and README examples to (hopefully) better explain
    method usage

## [v0.5.0] - 2020-09-20

### Changed

- Add explicit state labels to threshold list items
  - `CRITICAL:` and a single space
  - `WARNING:` and a single space
  - **BREAKING**: this will require updates to client code to accommodate this
    change
- Update whitespace/EOL handling within `Long Service Output` or `DETAILED
  INFO` section
  - **BREAKING**: this will require updates to client code to accommodate this
    change

### Fixed

- Don't assume that state thresholds will be provided
- `YYYY-MM-DD` changelog version entries

## [v0.4.0] - 2020-08-31

### Added

- Add initial "framework workflow"
  - `ExitState` type with `ReturnCheckResults` method
    - used to process and return all applicable check results to Nagios for
      further processing/display
    - supports "branding" callback function to display application name,
      version, or other information as a "trailer" for check results provided
      to Nagios
      - this could be useful for identifying what version of a plugin
        determined the service or host state to be an issue
  - README
    - extend examples to reflect new type/method

### Changed

- GoDoc coverage
  - simple example retained, reader referred to README for further examples

### Fixed

- GitHub Actions Workflow shallow build depth
- `YYYY-MM-DD` changelog version entries

## [v0.3.1] - 2020-08-23

### Added

- Docker-based GitHub Actions Workflows
  - Replace native GitHub Actions with containers created and managed through
    the `atc0005/go-ci` project.

  - New, primary workflow
    - with parallel linting, testing and building tasks
    - with three Go environments
      - "old stable"
      - "stable"
      - "unstable"
    - Makefile is *not* used in this workflow
    - staticcheck linting using latest stable version provided by the
      `atc0005/go-ci` containers

  - Separate Makefile-based linting and building workflow
    - intended to help ensure that local Makefile-based builds that are
      referenced in project README files continue to work as advertised until
      a better local tool can be discovered/explored further
    - use `golang:latest` container to allow for Makefile-based linting
      tooling installation testing since the `atc0005/go-ci` project provides
      containers with those tools already pre-installed
      - linting tasks use container-provided `golangci-lint` config file
        *except* for the Makefile-driven linting task which continues to use
        the repo-provided copy of the `golangci-lint` configuration file

  - Add Quick Validation workflow
    - run on every push, everything else on pull request updates
    - linting via `golangci-lint` only
    - testing
    - no builds

- Add new README badges for additional CI workflows
  - each badge also links to the associated workflow results

### Changed

- Disable `golangci-lint` default exclusions

- dependencies
  - `go.mod` Go version
    - updated from `1.13` to `1.14`
  - `actions/setup-go`
    - updated from `v2.1.0` to `v2.1.2`
      - since replaced with Docker containers
  - `actions/setup-node`
    - updated from `v2.1.0` to `v2.1.1`
  - `actions/checkout`
    - updated from `v2.3.1` to `v2.3.2`

- README
  - Link badges to applicable GitHub Actions workflows results
  - Add Table of contents

- Linting
  - Local
    - `Makefile`
      - install latest stable `golangci-lint` binary instead of using a fixed
          version
  - CI
    - remove repo-provided copy of `golangci-lint` config file at start of
      linting task in order to force use of Docker container-provided config
      file

## [v0.3.0] - 2020-07-05

### Added

- Add State "labels" to provide an alternative to using literal state strings

- Add GitHub Actions Workflow, Makefile for builds
  - Lint codebase
  - Build codebase

- Enable Dependabot updates
  - GitHub Actions
  - Go Modules

### Changed

- BREAKING: Rename existing exit code constants to explicitly note that they
  are exit codes
  - the thinking was that since we have text "labels" for state, it would be
    good to help clarify the difference between the new constants and the
    existing exit code constants

- Minor tweaks to README to reference changes, wording

- Update dependencies
  - `actions/checkout`
    - `v1` to `v2.3.1`
  - `actions/setup-go`
    - `v2.0.3` to `v2.1.0`
  - `actions/setup-node`
    - `v1` to `v2.1.0`

## [v0.2.0] - 2020-02-02

### Added

- Add Nagios `State` constants

### Removed

- Nagios `State` map

## [v0.1.0] - 2020-01-20

### Added

Initial package state

- Nagios state map

[Unreleased]: https://github.com/atc0005/go-nagios/compare/v0.11.0...HEAD
[v0.11.0]: https://github.com/atc0005/go-nagios/releases/tag/v0.11.0
[v0.10.2]: https://github.com/atc0005/go-nagios/releases/tag/v0.10.2
[v0.10.1]: https://github.com/atc0005/go-nagios/releases/tag/v0.10.1
[v0.10.0]: https://github.com/atc0005/go-nagios/releases/tag/v0.10.0
[v0.9.2]: https://github.com/atc0005/go-nagios/releases/tag/v0.9.2
[v0.9.1]: https://github.com/atc0005/go-nagios/releases/tag/v0.9.1
[v0.9.0]: https://github.com/atc0005/go-nagios/releases/tag/v0.9.0
[v0.8.2]: https://github.com/atc0005/go-nagios/releases/tag/v0.8.2
[v0.8.1]: https://github.com/atc0005/go-nagios/releases/tag/v0.8.1
[v0.8.0]: https://github.com/atc0005/go-nagios/releases/tag/v0.8.0
[v0.7.0]: https://github.com/atc0005/go-nagios/releases/tag/v0.7.0
[v0.6.1]: https://github.com/atc0005/go-nagios/releases/tag/v0.6.1
[v0.6.0]: https://github.com/atc0005/go-nagios/releases/tag/v0.6.0
[v0.5.3]: https://github.com/atc0005/go-nagios/releases/tag/v0.5.3
[v0.5.2]: https://github.com/atc0005/go-nagios/releases/tag/v0.5.2
[v0.5.1]: https://github.com/atc0005/go-nagios/releases/tag/v0.5.1
[v0.5.0]: https://github.com/atc0005/go-nagios/releases/tag/v0.5.0
[v0.4.0]: https://github.com/atc0005/go-nagios/releases/tag/v0.4.0
[v0.3.1]: https://github.com/atc0005/go-nagios/releases/tag/v0.3.1
[v0.3.0]: https://github.com/atc0005/go-nagios/releases/tag/v0.3.0
[v0.2.0]: https://github.com/atc0005/go-nagios/releases/tag/v0.2.0
[v0.1.0]: https://github.com/atc0005/go-nagios/releases/tag/v0.1.0
