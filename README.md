<!-- omit in toc -->

# check-vmware

Go-based tooling to monitor VMware environments; **NOT** affiliated with
or endorsed by VMware, Inc.

[![Latest Release](https://img.shields.io/github/release/atc0005/check-vmware.svg?style=flat-square)](https://github.com/atc0005/check-vmware/releases/latest)
[![GoDoc](https://godoc.org/github.com/atc0005/check-vmware?status.svg)](https://godoc.org/github.com/atc0005/check-vmware)
[![Validate Codebase](https://github.com/atc0005/check-vmware/workflows/Validate%20Codebase/badge.svg)](https://github.com/atc0005/check-vmware/actions?query=workflow%3A%22Validate+Codebase%22)
[![Validate Docs](https://github.com/atc0005/check-vmware/workflows/Validate%20Docs/badge.svg)](https://github.com/atc0005/check-vmware/actions?query=workflow%3A%22Validate+Docs%22)
[![Lint and Build using Makefile](https://github.com/atc0005/check-vmware/workflows/Lint%20and%20Build%20using%20Makefile/badge.svg)](https://github.com/atc0005/check-vmware/actions?query=workflow%3A%22Lint+and+Build+using+Makefile%22)
[![Quick Validation](https://github.com/atc0005/check-vmware/workflows/Quick%20Validation/badge.svg)](https://github.com/atc0005/check-vmware/actions?query=workflow%3A%22Quick+Validation%22)

<!-- omit in toc -->
## Table of Contents

- [check-vmware](#check-vmware)
  - [Project home](#project-home)
  - [Overview](#overview)
    - [`check_vmware_tools`](#check_vmware_tools)
    - [`check_vmware_vcpus`](#check_vmware_vcpus)
    - [`check_vmware_vhw`](#check_vmware_vhw)
    - [`check_vmware_hs2ds2vms`](#check_vmware_hs2ds2vms)
  - [Features](#features)
  - [Changelog](#changelog)
  - [Requirements](#requirements)
    - [Building source code](#building-source-code)
    - [Running](#running)
  - [Installation](#installation)
  - [Configuration options](#configuration-options)
    - [Threshold calculations](#threshold-calculations)
      - [`check_vmware_tools`](#check_vmware_tools-1)
      - [`check_vmware_vcpus`](#check_vmware_vcpus-1)
      - [`check_vmware_vhw`](#check_vmware_vhw-1)
      - [`check_vmware_hs2ds2vms`](#check_vmware_hs2ds2vms-1)
    - [Command-line arguments](#command-line-arguments)
      - [Shared](#shared)
      - [`check_vmware_tools`](#check_vmware_tools-2)
      - [`check_vmware_vcpus`](#check_vmware_vcpus-2)
      - [`check_vmware_vhw`](#check_vmware_vhw-2)
      - [`check_vmware_hs2ds2vms`](#check_vmware_hs2ds2vms-2)
    - [Configuration file](#configuration-file)
  - [Examples](#examples)
    - [`check_vmware_tools` Nagios plugin](#check_vmware_tools-nagios-plugin)
      - [OK results](#ok-results)
      - [WARNING results](#warning-results)
      - [CRITICAL results](#critical-results)
    - [`check_vmware_vcpus` Nagios plugin](#check_vmware_vcpus-nagios-plugin)
      - [OK results](#ok-results-1)
      - [WARNING results](#warning-results-1)
      - [CRITICAL results](#critical-results-1)
    - [`check_vmware_vhw` Nagios plugin](#check_vmware_vhw-nagios-plugin)
      - [OK results](#ok-results-2)
      - [WARNING results](#warning-results-2)
      - [CRITICAL results](#critical-results-2)
    - [`check_vmware_hs2ds2vms` Nagios plugin](#check_vmware_hs2ds2vms-nagios-plugin)
      - [OK results](#ok-results-3)
      - [WARNING results](#warning-results-3)
      - [CRITICAL results](#critical-results-3)
  - [License](#license)
  - [References](#references)

## Project home

See [our GitHub repo](https://github.com/atc0005/check-vmware) for the latest
code, to file an issue or submit improvements for review and potential
inclusion into the project.

Just to be 100% clear: this project is not affiliated with or endorsed by
VMware, Inc.

## Overview

This repo contains various tools used to monitor/validate VMware environments.

| Tool Name                | Status | Description                                                       |
| ------------------------ | ------ | ----------------------------------------------------------------- |
| `check_vmware_tools`     | Alpha  | Nagios plugin used to monitor VMware Tools installations.         |
| `check_vmware_vcpus`     | Alpha  | Nagios plugin used to monitor allocation of virtual CPUs (vCPUs). |
| `check_vmware_vhw`       | Alpha  | Nagios plugin used to monitor virtual hardware versions.          |
| `check_vmware_hs2ds2vms` | Alpha  | Nagios plugin used to monitor host/datastore/vm pairings.         |

The output for these plugins is designed to provide the one-line summary
needed by Nagios for quick identification of a problem while providing longer,
more detailed information for use in email and Teams notifications
([atc0005/send2teams](https://github.com/atc0005/send2teams)).

Some plugins provide optional support to limit evaluation of VMs to specific
Resource Pools (explicitly including or excluding) and power states (on or
off). See the [Configuration options](#configuration-options) section for more
information.

### `check_vmware_tools`

Nagios plugin used to monitor VMware Tools installations. See the
[Configuration options](#configuration-options) section for details regarding
how the various Tools states are evaluated.

### `check_vmware_vcpus`

Nagios plugin used to monitor allocation of virtual CPUs (vCPUs).

Thresholds for `CRITICAL` and `WARNING` vCPUs allocation have usable defaults,
but Max vCPUs allocation is required before this plugin can be used.

### `check_vmware_vhw`

Nagios plugin used to monitor virtual hardware versions.

As of this writing, I am unaware of a way to query the current vSphere
environment for the latest available hardware version. As a workaround for
that lack of knowledge, this plugin applies an automatic baseline of "highest
version discovered" across evaluated VMs. Any VMs with a hardware version not
at that highest version are flagged as problematic. Please file an issue or
open a discussion in this project's repo if you're aware of a way to directly
query the desired value from the current vSphere environment.

Instead of trying to determine how far behind each VM is from the newest
version, this plugin assumes that any deviation is a `WARNING` level issue.
See GH-33 for future potential changes to this behavior.

### `check_vmware_hs2ds2vms`

Nagios plugin used to monitor host/datastore/vm pairings.

This is a functional plugin responsible for verifying that each VM is housed
on a datastore (best) intended for the host associated with the VM.

By default, the evaluation is limited to powered on VMs, but this can be
toggled to also include powered off VMs.

The association between datastores and hosts is determined by a user-provided
Custom Attribute. Flags for this plugin allow specifying separate Custom
Attribute names for hosts and datastores along with optional separate prefixes
for the provided Custom Attributes.

This allows for example, hosts to use a `Location` Custom Attribute that
shares a datacenter name with datastores using the same `Location` Custom
Attribute. If not specifying a prefix separator, the plugin assumes that a
literal, case-insensitive match of the `Location` field is required. If a
prefix separator is provided, then the separator is used to retrieve the
common prefix for the `Location` Custom Attribute for both hosts and
datastores.

This is intended to work around hosts that may include both the datacenter
name and rack location details in their Custom Attribute (e.g., `Location`).

This plugin optionally allows ignoring a list of datastores, and both hosts
and datastores that are missing the specified Custom Attribute.

In addition to specifying separate Custom Attribute names (required) and
prefix separators (optional), the plugin also accepts a single Custom
Attribute used by both hosts and datastores and an optional prefix separator,
also used by both hosts and datastores.

If specifying a shared Custom Attribute or prefix, per-resource Custom
Attribute flags are rejected (error condition).

Known issues:

- While testing indicates that this plugin is functional, significant
  refactoring is still needed

- As with other plugins, documentation is incredibly thin

## Features

- Multiple plugins for monitoring VMware vSphere environments (standalone ESXi
  hosts or vCenter instances) for select (or all) Resource Pools.
  - VMware Tools
  - Virtual CPU allocations
  - Virtual hardware versions
  - Host/Datastore/Virtual Machine pairings (using provided Custom Attribute)

- Optional, leveled logging using `rs/zerolog` package
  - JSON-format output (to `stderr`)
  - choice of `disabled`, `panic`, `fatal`, `error`, `warn`, `info` (the
    default), `debug` or `trace`.

- Optional, user-specified timeout value for plugin execution.

## Changelog

See the [`CHANGELOG.md`](CHANGELOG.md) file for the changes associated with
each release of this application. Changes that have been merged to `master`,
but not yet an official release may also be noted in the file under the
`Unreleased` section. A helpful link to the Git commit history since the last
official release is also provided for further review.

## Requirements

The following is a loose guideline. Other combinations of Go and operating
systems for building and running tools from this repo may work, but have not
been tested.

### Building source code

- Go 1.14+
- GCC
  - if building with custom options (as the provided `Makefile` does)
- `make`
  - if using the provided `Makefile`

### Running

- Windows 7, Server 2008R2 or later
  - per official [Go install notes][go-docs-install]
- Windows 10 Version 1909
  - tested
- Ubuntu Linux 16.04, 18.04

## Installation

1. [Download][go-docs-download] Go
1. [Install][go-docs-install] Go
   - NOTE: Pay special attention to the remarks about `$HOME/.profile`
1. Clone the repo
   1. `cd /tmp`
   1. `git clone https://github.com/atc0005/check-vmware`
   1. `cd check-vmware`
1. Install dependencies (optional)
   - for Ubuntu Linux
     - `sudo apt-get install make gcc`
   - for CentOS Linux
     - `sudo yum install make gcc`
   - for Windows
     - Emulated environments (*easier*)
       - Skip all of this and build using the default `go build` command in
         Windows (see below for use of the `-mod=vendor` flag)
       - build using Windows Subsystem for Linux Ubuntu environment and just
         copy out the Windows binaries from that environment
       - If already running a Docker environment, use a container with the Go
         tool-chain already installed
       - If already familiar with LXD, create a container and follow the
         installation steps given previously to install required dependencies
     - Native tooling (*harder*)
       - see the StackOverflow Question `32127524` link in the
         [References](references.md) section for potential options for
         installing `make` on Windows
       - see the mingw-w64 project homepage link in the
         [References](references.md) section for options for installing `gcc`
         and related packages on Windows
1. Build binaries
   - for the current operating system, explicitly using bundled dependencies
         in top-level `vendor` folder
     - `go build -mod=vendor ./cmd/check_vmware_tools/`
     - `go build -mod=vendor ./cmd/check_vmware_vcpus/`
     - `go build -mod=vendor ./cmd/check_vmware_vhw/`
     - `go build -mod=vendor ./cmd/check_vmware_hs2ds2vms/`
   - for all supported platforms (where `make` is installed)
      - `make all`
   - for use on Windows
      - `make windows`
   - for use on Linux
     - `make linux`
1. Copy the newly compiled binary from the applicable `/tmp` subdirectory path
   (based on the clone instructions in this section) below and deploy where
   needed.
   - if using `Makefile`
     - look in `/tmp/check-vmware/release_assets/check_vmware_tools/`
     - look in `/tmp/check-vmware/release_assets/check_vmware_vcpus/`
     - look in `/tmp/check-vmware/release_assets/check_vmware_vhw/`
     - look in `/tmp/check-vmware/release_assets/check_vmware_hs2ds2vms/`
   - if using `go build`
     - look in `/tmp/check-vmware/`

## Configuration options

### Threshold calculations

#### `check_vmware_tools`

| Tools Status        | Nagios State | Description                                                                                                              |
| ------------------- | ------------ | ------------------------------------------------------------------------------------------------------------------------ |
| `toolsOk`           | `OK`         | Ideal state, no problems with VMware Tools (or `open-vm-tools`) detected.                                                |
| `toolsOld`          | `WARNING`    | Outdated VMware Tools installation. The host ESXi system was likely recently updated.                                    |
| `toolsNotRunning`   | `CRITICAL`   | VMware Tools (or `open-vm-tools`) not currently running. It likely crashed or was terminated due to low memory scenario. |
| `toolsNotInstalled` | `CRITICAL`   | Fresh virtual environment, or VMware Tools removed as part of an upgrade of an existing installation.                    |

#### `check_vmware_vcpus`

| Nagios State | Description                                                       |
| ------------ | ----------------------------------------------------------------- |
| `OK`         | Ideal state, vCPU allocations within bounds.                      |
| `WARNING`    | vCPU allocations crossed user-specified threshold for this state. |
| `CRITICAL`   | vCPU allocations crossed user-specified threshold for this state. |

#### `check_vmware_vhw`

| Nagios State | Description                                                                  |
| ------------ | ---------------------------------------------------------------------------- |
| `OK`         | Ideal state, no mismatched Host/Datastore/Virtual machine pairings detected. |
| `WARNING`    | Non-homogenous hardware versions.                                            |
| `CRITICAL`   | Not used by this plugin.                                                     |

#### `check_vmware_hs2ds2vms`

| Nagios State | Description                                                                  |
| ------------ | ---------------------------------------------------------------------------- |
| `OK`         | Ideal state, no mismatched Host/Datastore/Virtual machine pairings detected. |
| `WARNING`    | Not used by this plugin.                                                     |
| `CRITICAL`   | Any errors encountered or Hosts/Datastores/VM mismatches.                    |

### Command-line arguments

- Use the `-h` or `--help` flag to display current usage information.
- Flags marked as **`required`** must be set via CLI flag.
- Flags *not* marked as required are for settings where a useful default is
  already defined, but may be overridden if desired.

#### Shared

TODO: Remove this section and instead duplicate the full table for each plugin.

#### `check_vmware_tools`

| Flag              | Required | Default | Repeat | Possible                                                                | Description                                                                                                                                                                                                        |
| ----------------- | -------- | ------- | ------ | ----------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `branding`        | No       | `false` | No     | `branding`                                                              | Toggles emission of branding details with plugin status details. This output is disabled by default.                                                                                                               |
| `h`, `help`       | No       | `false` | No     | `h`, `help`                                                             | Show Help text along with the list of supported flags.                                                                                                                                                             |
| `v`, `version`    | No       | `false` | No     | `v`, `version`                                                          | Whether to display application version and then immediately exit application.                                                                                                                                      |
| `ll`, `log-level` | No       | `info`  | No     | `disabled`, `panic`, `fatal`, `error`, `warn`, `info`, `debug`, `trace` | Log message priority filter. Log messages with a lower level are ignored.                                                                                                                                          |
| `p`, `port`       | No       | `443`   | No     | *positive whole number between 1-65535, inclusive*                      | TCP port of the remote ESXi host or vCenter instance. This is usually 443 (HTTPS).                                                                                                                                 |
| `t`, `timeout`    | No       | `10`    | No     | *positive whole number of seconds*                                      | Timeout value in seconds allowed before a plugin execution attempt is abandoned and an error returned.                                                                                                             |
| `s`, `server`     | **Yes**  |         | No     | *fully-qualified domain name or IP Address*                             | The fully-qualified domain name or IP Address of the remote ESXi host or vCenter instance.                                                                                                                         |
| `u`, `username`   | **Yes**  |         | No     | *valid username*                                                        | Username with permission to access specified ESXi host or vCenter instance.                                                                                                                                        |
| `pw`, `password`  | **Yes**  |         | No     | *valid password*                                                        | Password used to login to ESXi host or vCenter instance.                                                                                                                                                           |
| `domain`          | No       |         | No     | *valid password*                                                        | (Optional) domain for user account used to login to ESXi host or vCenter instance.                                                                                                                                 |
| `trust-cert`      | No       | `false` | No     | `true`, `false`                                                         | Whether the certificate should be trusted as-is without validation. WARNING: TLS is susceptible to man-in-the-middle attacks if enabling this option.                                                              |
| `include-rp`      | No       |         | No     | *comma-separated list of resource pool names*                           | Specifies a comma-separated list of Resource Pools that should be exclusively used when evaluating VMs. This option is incompatible with specifying a list of Resource Pools to ignore or exclude from evaluation. |
| `exclude-rp`      | No       |         | No     | *comma-separated list of resource pool names*                           | Specifies a comma-separated list of Resource Pools that should be ignored when evaluating VMs. This option is incompatible with specifying a list of Resource Pools to include for evaluation.                     |
| `ignore-vm`       | No       |         | No     | *comma-separated list of (vSphere) virtual machine names*               | Specifies a comma-separated list of VM names that should be ignored or excluded from evaluation.                                                                                                                   |
| `powered-off`     | No       | `false` | No     | `true`, `false`                                                         | Toggles evaluation of powered off VMs in addition to powered on VMs. Evaluation of powered off VMs is disabled by default.                                                                                         |

#### `check_vmware_vcpus`

| Flag                        | Required | Default | Repeat | Possible                                                                | Description                                                                                                                                                                                                        |
| --------------------------- | -------- | ------- | ------ | ----------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `branding`                  | No       | `false` | No     | `branding`                                                              | Toggles emission of branding details with plugin status details. This output is disabled by default.                                                                                                               |
| `h`, `help`                 | No       | `false` | No     | `h`, `help`                                                             | Show Help text along with the list of supported flags.                                                                                                                                                             |
| `v`, `version`              | No       | `false` | No     | `v`, `version`                                                          | Whether to display application version and then immediately exit application.                                                                                                                                      |
| `ll`, `log-level`           | No       | `info`  | No     | `disabled`, `panic`, `fatal`, `error`, `warn`, `info`, `debug`, `trace` | Log message priority filter. Log messages with a lower level are ignored.                                                                                                                                          |
| `p`, `port`                 | No       | `443`   | No     | *positive whole number between 1-65535, inclusive*                      | TCP port of the remote ESXi host or vCenter instance. This is usually 443 (HTTPS).                                                                                                                                 |
| `t`, `timeout`              | No       | `10`    | No     | *positive whole number of seconds*                                      | Timeout value in seconds allowed before a plugin execution attempt is abandoned and an error returned.                                                                                                             |
| `s`, `server`               | **Yes**  |         | No     | *fully-qualified domain name or IP Address*                             | The fully-qualified domain name or IP Address of the remote ESXi host or vCenter instance.                                                                                                                         |
| `u`, `username`             | **Yes**  |         | No     | *valid username*                                                        | Username with permission to access specified ESXi host or vCenter instance.                                                                                                                                        |
| `pw`, `password`            | **Yes**  |         | No     | *valid password*                                                        | Password used to login to ESXi host or vCenter instance.                                                                                                                                                           |
| `domain`                    | No       |         | No     | *valid password*                                                        | (Optional) domain for user account used to login to ESXi host or vCenter instance.                                                                                                                                 |
| `trust-cert`                | No       | `false` | No     | `true`, `false`                                                         | Whether the certificate should be trusted as-is without validation. WARNING: TLS is susceptible to man-in-the-middle attacks if enabling this option.                                                              |
| `include-rp`                | No       |         | No     | *comma-separated list of resource pool names*                           | Specifies a comma-separated list of Resource Pools that should be exclusively used when evaluating VMs. This option is incompatible with specifying a list of Resource Pools to ignore or exclude from evaluation. |
| `exclude-rp`                | No       |         | No     | *comma-separated list of resource pool names*                           | Specifies a comma-separated list of Resource Pools that should be ignored when evaluating VMs. This option is incompatible with specifying a list of Resource Pools to include for evaluation.                     |
| `ignore-vm`                 | No       |         | No     | *comma-separated list of (vSphere) virtual machine names*               | Specifies a comma-separated list of VM names that should be ignored or excluded from evaluation.                                                                                                                   |
| `powered-off`               | No       | `false` | No     | `true`, `false`                                                         | Toggles evaluation of powered off VMs in addition to powered on VMs. Evaluation of powered off VMs is disabled by default.                                                                                         |
| `vcma`, `vcpus-max-allowed` | **Yes**  | `0`     | No     | *positive whole number of vCPUs*                                        | Specifies the maximum amount of virtual CPUs (as a whole number) that we are allowed to allocate in the target VMware environment.                                                                                 |
| `vc`, `vcpus-critical`      | No       | `100`   | No     | *percentage as positive whole number*                                   | Specifies the percentage of vCPUs allocation (as a whole number) when a CRITICAL threshold is reached.                                                                                                             |
| `vw`, `vcpus-warning`       | No       | `95`    | No     | *percentage as positive whole number*                                   | Specifies the percentage of vCPUs allocation (as a whole number) when a WARNING threshold is reached.                                                                                                              |

#### `check_vmware_vhw`

| Flag              | Required | Default | Repeat | Possible                                                                | Description                                                                                                                                                                                                        |
| ----------------- | -------- | ------- | ------ | ----------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `branding`        | No       | `false` | No     | `branding`                                                              | Toggles emission of branding details with plugin status details. This output is disabled by default.                                                                                                               |
| `h`, `help`       | No       | `false` | No     | `h`, `help`                                                             | Show Help text along with the list of supported flags.                                                                                                                                                             |
| `v`, `version`    | No       | `false` | No     | `v`, `version`                                                          | Whether to display application version and then immediately exit application.                                                                                                                                      |
| `ll`, `log-level` | No       | `info`  | No     | `disabled`, `panic`, `fatal`, `error`, `warn`, `info`, `debug`, `trace` | Log message priority filter. Log messages with a lower level are ignored.                                                                                                                                          |
| `p`, `port`       | No       | `443`   | No     | *positive whole number between 1-65535, inclusive*                      | TCP port of the remote ESXi host or vCenter instance. This is usually 443 (HTTPS).                                                                                                                                 |
| `t`, `timeout`    | No       | `10`    | No     | *positive whole number of seconds*                                      | Timeout value in seconds allowed before a plugin execution attempt is abandoned and an error returned.                                                                                                             |
| `s`, `server`     | **Yes**  |         | No     | *fully-qualified domain name or IP Address*                             | The fully-qualified domain name or IP Address of the remote ESXi host or vCenter instance.                                                                                                                         |
| `u`, `username`   | **Yes**  |         | No     | *valid username*                                                        | Username with permission to access specified ESXi host or vCenter instance.                                                                                                                                        |
| `pw`, `password`  | **Yes**  |         | No     | *valid password*                                                        | Password used to login to ESXi host or vCenter instance.                                                                                                                                                           |
| `domain`          | No       |         | No     | *valid password*                                                        | (Optional) domain for user account used to login to ESXi host or vCenter instance.                                                                                                                                 |
| `trust-cert`      | No       | `false` | No     | `true`, `false`                                                         | Whether the certificate should be trusted as-is without validation. WARNING: TLS is susceptible to man-in-the-middle attacks if enabling this option.                                                              |
| `include-rp`      | No       |         | No     | *comma-separated list of resource pool names*                           | Specifies a comma-separated list of Resource Pools that should be exclusively used when evaluating VMs. This option is incompatible with specifying a list of Resource Pools to ignore or exclude from evaluation. |
| `exclude-rp`      | No       |         | No     | *comma-separated list of resource pool names*                           | Specifies a comma-separated list of Resource Pools that should be ignored when evaluating VMs. This option is incompatible with specifying a list of Resource Pools to include for evaluation.                     |
| `ignore-vm`       | No       |         | No     | *comma-separated list of (vSphere) virtual machine names*               | Specifies a comma-separated list of VM names that should be ignored or excluded from evaluation.                                                                                                                   |
| `powered-off`     | No       | `false` | No     | `true`, `false`                                                         | Toggles evaluation of powered off VMs in addition to powered on VMs. Evaluation of powered off VMs is disabled by default.                                                                                         |

#### `check_vmware_hs2ds2vms`

| Flag                 | Required  | Default | Repeat | Possible                                                                | Description                                                                                                                                                                                                                                               |
| -------------------- | --------- | ------- | ------ | ----------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `branding`           | No        | `false` | No     | `branding`                                                              | Toggles emission of branding details with plugin status details. This output is disabled by default.                                                                                                                                                      |
| `h`, `help`          | No        | `false` | No     | `h`, `help`                                                             | Show Help text along with the list of supported flags.                                                                                                                                                                                                    |
| `v`, `version`       | No        | `false` | No     | `v`, `version`                                                          | Whether to display application version and then immediately exit application.                                                                                                                                                                             |
| `ll`, `log-level`    | No        | `info`  | No     | `disabled`, `panic`, `fatal`, `error`, `warn`, `info`, `debug`, `trace` | Log message priority filter. Log messages with a lower level are ignored.                                                                                                                                                                                 |
| `p`, `port`          | No        | `443`   | No     | *positive whole number between 1-65535, inclusive*                      | TCP port of the remote ESXi host or vCenter instance. This is usually 443 (HTTPS).                                                                                                                                                                        |
| `t`, `timeout`       | No        | `10`    | No     | *positive whole number of seconds*                                      | Timeout value in seconds allowed before a plugin execution attempt is abandoned and an error returned.                                                                                                                                                    |
| `s`, `server`        | **Yes**   |         | No     | *fully-qualified domain name or IP Address*                             | The fully-qualified domain name or IP Address of the remote ESXi host or vCenter instance.                                                                                                                                                                |
| `u`, `username`      | **Yes**   |         | No     | *valid username*                                                        | Username with permission to access specified ESXi host or vCenter instance.                                                                                                                                                                               |
| `pw`, `password`     | **Yes**   |         | No     | *valid password*                                                        | Password used to login to ESXi host or vCenter instance.                                                                                                                                                                                                  |
| `domain`             | No        |         | No     | *valid password*                                                        | (Optional) domain for user account used to login to ESXi host or vCenter instance.                                                                                                                                                                        |
| `trust-cert`         | No        | `false` | No     | `true`, `false`                                                         | Whether the certificate should be trusted as-is without validation. WARNING: TLS is susceptible to man-in-the-middle attacks if enabling this option.                                                                                                     |
| `include-rp`         | No        |         | No     | *comma-separated list of resource pool names*                           | Specifies a comma-separated list of Resource Pools that should be exclusively used when evaluating VMs. This option is incompatible with specifying a list of Resource Pools to ignore or exclude from evaluation.                                        |
| `exclude-rp`         | No        |         | No     | *comma-separated list of resource pool names*                           | Specifies a comma-separated list of Resource Pools that should be ignored when evaluating VMs. This option is incompatible with specifying a list of Resource Pools to include for evaluation.                                                            |
| `ignore-vm`          | No        |         | No     | *comma-separated list of (vSphere) virtual machine names*               | Specifies a comma-separated list of VM names that should be ignored or excluded from evaluation.                                                                                                                                                          |
| `ignore-ds`          | No        |         | No     | *comma-separated list of (vSphere) datastore names*                     | Specifies a comma-separated list of Datastore names that should be ignored or excluded from evaluation.                                                                                                                                                   |
| `powered-off`        | No        | `false` | No     | `true`, `false`                                                         | Toggles evaluation of powered off VMs in addition to powered on VMs. Evaluation of powered off VMs is disabled by default.                                                                                                                                |
| `ca-name`            | **Maybe** |         | No     | *valid Custom Attribute name*                                           | Custom Attribute name for host ESXi systems and datastores. Optional if specifying resource-specific custom attribute names.                                                                                                                              |
| `ca-prefix-sep`      | **Maybe** |         | No     | *valid Custom Attribute prefix separator character*                     | Custom Attribute prefix separator for host ESXi systems and datastores. Skip if using Custom Attribute values as-is for comparison, otherwise optional if specifying resource-specific custom attribute prefix separator, or using the default separator. |
| `ignore-missing-ca`  | No        | `false` | No     | `true`, `false`                                                         | Toggles how missing specified Custom Attributes will be handled. By default, ESXi hosts and datastores missing the Custom Attribute are treated as an error condition.                                                                                    |
| `host-ca-name`       | **Maybe** |         | No     | *valid Custom Attribute name*                                           | Custom Attribute name specific to host ESXi systems. Optional if specifying shared custom attribute flag.                                                                                                                                                 |
| `host-ca-prefix-sep` | **Maybe** |         | No     | *valid Custom Attribute prefix separator character*                     | Custom Attribute prefix separator specific to host ESXi systems. Skip if using Custom Attribute values as-is for comparison, otherwise optional if specifying shared custom attribute prefix separator, or using the default separator.                   |
| `ds-ca-name`         | **Maybe** |         | No     | *valid Custom Attribute name*                                           | Custom Attribute name specific to datastores. Optional if specifying shared custom attribute flag.                                                                                                                                                        |
| `ds-ca-prefix-sep`   | **Maybe** |         | No     | *valid Custom Attribute prefix separator character*                     | Custom Attribute prefix separator specific to datastores. Skip if using Custom Attribute values as-is for comparison, otherwise optional if specifying shared custom attribute prefix separator, or using the default separator.                          |

### Configuration file

Not currently supported. This feature may be added later if there is
sufficient interest.

## Examples

### `check_vmware_tools` Nagios plugin

#### OK results

TODO

#### WARNING results

TODO

#### CRITICAL results

TODO

### `check_vmware_vcpus` Nagios plugin

#### OK results

TODO

#### WARNING results

TODO

#### CRITICAL results

TODO

### `check_vmware_vhw` Nagios plugin

#### OK results

TODO

#### WARNING results

TODO

#### CRITICAL results

TODO

### `check_vmware_hs2ds2vms` Nagios plugin

#### OK results

TODO

#### WARNING results

TODO

#### CRITICAL results

TODO

## License

From the [LICENSE](LICENSE) file:

```license
MIT License

Copyright (c) 2021 Adam Chalkley

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

## References

- vSphere
  - [Go library for the VMware vSphere API](https://github.com/vmware/govmomi)
  - [vSphere Web Services API](https://code.vmware.com/apis/1067/vsphere)

- Logging
  - <https://github.com/rs/zerolog>

- Nagios
  - <https://github.com/atc0005/go-nagios>
  - <https://nagios-plugins.org/doc/guidelines.html>

<!-- Footnotes here  -->

[repo-url]: <https://github.com/atc0005/check-vmware>  "This project's GitHub repo"

[go-docs-download]: <https://golang.org/dl>  "Download Go"

[go-docs-install]: <https://golang.org/doc/install>  "Install Go"

<!-- []: PLACEHOLDER "DESCRIPTION_HERE" -->
