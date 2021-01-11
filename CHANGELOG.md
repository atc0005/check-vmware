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

[Unreleased]: https://github.com/atc0005/check-vmware/compare/v0.1.1...HEAD
[v0.1.1]: https://github.com/atc0005/check-vmware/releases/tag/v0.1.1
[v0.1.0]: https://github.com/atc0005/check-vmware/releases/tag/v0.1.0
