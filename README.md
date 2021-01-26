<!-- omit in toc -->

# check-vmware

Go-based tooling to monitor VMware environments; **NOT** affiliated with
or endorsed by VMware, Inc.

[![Latest Release](https://img.shields.io/github/release/atc0005/check-vmware.svg?style=flat-square)](https://github.com/atc0005/check-vmware/releases/latest)
[![Go Reference](https://pkg.go.dev/badge/github.com/atc0005/check-vmware.svg)](https://pkg.go.dev/github.com/atc0005/check-vmware)
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
    - [`check_vmware_datastore`](#check_vmware_datastore)
    - [`check_vmware_snapshots_age`](#check_vmware_snapshots_age)
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
      - [`check_vmware_datastore`](#check_vmware_datastore-1)
      - [`check_vmware_snapshots_age`](#check_vmware_snapshots_age-1)
    - [Command-line arguments](#command-line-arguments)
      - [`check_vmware_tools`](#check_vmware_tools-2)
      - [`check_vmware_vcpus`](#check_vmware_vcpus-2)
      - [`check_vmware_vhw`](#check_vmware_vhw-2)
      - [`check_vmware_hs2ds2vms`](#check_vmware_hs2ds2vms-2)
      - [`check_vmware_datastore`](#check_vmware_datastore-2)
      - [`check_vmware_snapshots_age`](#check_vmware_snapshots_age-2)
    - [Configuration file](#configuration-file)
  - [Contrib](#contrib)
  - [Examples](#examples)
    - [`check_vmware_tools` Nagios plugin](#check_vmware_tools-nagios-plugin)
      - [CLI invocation](#cli-invocation)
      - [Command definition](#command-definition)
    - [`check_vmware_vcpus` Nagios plugin](#check_vmware_vcpus-nagios-plugin)
      - [CLI invocation](#cli-invocation-1)
      - [Command definition](#command-definition-1)
    - [`check_vmware_vhw` Nagios plugin](#check_vmware_vhw-nagios-plugin)
      - [CLI invocation](#cli-invocation-2)
      - [Command definition](#command-definition-2)
    - [`check_vmware_hs2ds2vms` Nagios plugin](#check_vmware_hs2ds2vms-nagios-plugin)
      - [CLI invocation](#cli-invocation-3)
      - [Command definition](#command-definition-3)
    - [`check_vmware_datastore` Nagios plugin](#check_vmware_datastore-nagios-plugin)
      - [CLI invocation](#cli-invocation-4)
      - [Command definition](#command-definition-4)
    - [`check_vmware_snapshots_age` Nagios plugin](#check_vmware_snapshots_age-nagios-plugin)
      - [CLI invocation](#cli-invocation-5)
      - [Command definition](#command-definition-5)
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

| Tool Name                    | Status | Description                                                         |
| ---------------------------- | ------ | ------------------------------------------------------------------- |
| `check_vmware_tools`         | Alpha  | Nagios plugin used to monitor VMware Tools installations.           |
| `check_vmware_vcpus`         | Alpha  | Nagios plugin used to monitor allocation of virtual CPUs (vCPUs).   |
| `check_vmware_vhw`           | Alpha  | Nagios plugin used to monitor virtual hardware versions.            |
| `check_vmware_hs2ds2vms`     | Alpha  | Nagios plugin used to monitor host/datastore/vm pairings.           |
| `check_vmware_datastore`     | Alpha  | Nagios plugin used to monitor datastore usage.                      |
| `check_vmware_snapshots_age` | Alpha  | Nagios plugin used to monitor the age of Virtual Machine snapshots. |

The output for these plugins is designed to provide the one-line summary
needed by Nagios for quick identification of a problem while providing longer,
more detailed information for use in email and Teams notifications
([atc0005/send2teams](https://github.com/atc0005/send2teams)).

Some plugins provide optional support to limit evaluation of VMs to specific
Resource Pools (explicitly including or excluding) and power states (on or
off). See the [configuration options](#configuration-options),
[examples](#examples) and [contrib](#contrib) sections for more information.

### `check_vmware_tools`

Nagios plugin used to monitor VMware Tools installations. See the
[configuration options](#configuration-options) section for details regarding
how the various Tools states are evaluated.

### `check_vmware_vcpus`

Nagios plugin used to monitor allocation of virtual CPUs (vCPUs).

Thresholds for `CRITICAL` and `WARNING` vCPUs allocation have usable defaults,
but Max vCPUs allocation is required before this plugin can be used. See the
[configuration options](#configuration-options) section for details.

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

### `check_vmware_datastore`

Nagios plugin used to monitor datastore usage.

In addition to reporting current datastore usage, this plugin also reports
which VMs reside on the datastore along with their percentage of the total
datastore space used.

### `check_vmware_snapshots_age`

Nagios plugin used to monitor the age of Virtual Machine snapshots.

The current design of this plugin is to evaluate *all* Virtual Machines,
whether powered off or powered on. If you have a use case for evaluating
*only* powered on VMs by default, please file a GitHub issue in this project
providing some details for your use-case. In our environment, I have yet to
see a need to *only* evaluate powered on VMs for old snapshots. For cases
where the snapshots needed to be ignored, we added the VM to the ignore list.
We then relied on datastore usage monitoring to let us know when space was
becoming an issue.

Thresholds for `CRITICAL` and `WARNING` age values have usable defaults, but
may require adjustment for your environment. See the [configuration
options](#configuration-options) section for details.

## Features

- Multiple plugins for monitoring VMware vSphere environments (standalone ESXi
  hosts or vCenter instances) for select (or all) Resource Pools.
  - VMware Tools
  - Virtual CPU allocations
  - Virtual hardware versions
  - Host/Datastore/Virtual Machine pairings (using provided Custom Attribute)
  - Datastore usage
  - Snapshots age

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
     - `go build -mod=vendor ./cmd/check_vmware_datastore/`
     - `go build -mod=vendor ./cmd/check_vmware_snapshots_age/`
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
     - look in `/tmp/check-vmware/release_assets/check_vmware_datastore/`
     - look in `/tmp/check-vmware/release_assets/check_vmware_snapshots_age/`
   - if using `go build`
     - look in `/tmp/check-vmware/`
1. Review [configuration options](#configuration-options),
   [`examples`](#examples) and [`contrib`](#contrib) sections usage details.

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

| Nagios State | Description                                |
| ------------ | ------------------------------------------ |
| `OK`         | Ideal state, homogenous hardware versions. |
| `WARNING`    | Non-homogenous hardware versions.          |
| `CRITICAL`   | Not used by this plugin.                   |

#### `check_vmware_hs2ds2vms`

| Nagios State | Description                                                                  |
| ------------ | ---------------------------------------------------------------------------- |
| `OK`         | Ideal state, no mismatched Host/Datastore/Virtual machine pairings detected. |
| `WARNING`    | Not used by this plugin.                                                     |
| `CRITICAL`   | Any errors encountered or Hosts/Datastores/VM mismatches.                    |

#### `check_vmware_datastore`

| Nagios State | Description                                                      |
| ------------ | ---------------------------------------------------------------- |
| `OK`         | Ideal state, Datastore usage within bounds.                      |
| `WARNING`    | Datastore usage crossed user-specified threshold for this state. |
| `CRITICAL`   | Datastore usage crossed user-specified threshold for this state. |

#### `check_vmware_snapshots_age`

| Nagios State | Description                                                    |
| ------------ | -------------------------------------------------------------- |
| `OK`         | Ideal state, snapshots age within bounds.                      |
| `WARNING`    | Snapshots age crossed user-specified threshold for this state. |
| `CRITICAL`   | Snapshots age crossed user-specified threshold for this state. |

### Command-line arguments

- Use the `-h` or `--help` flag to display current usage information.
- Flags marked as **`required`** must be set via CLI flag.
- Flags *not* marked as required are for settings where a useful default is
  already defined, but may be overridden if desired.

#### `check_vmware_tools`

| Flag              | Required | Default | Repeat | Possible                                                                | Description                                                                                                                                                                                                                                                                                                                |
| ----------------- | -------- | ------- | ------ | ----------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `branding`        | No       | `false` | No     | `branding`                                                              | Toggles emission of branding details with plugin status details. This output is disabled by default.                                                                                                                                                                                                                       |
| `h`, `help`       | No       | `false` | No     | `h`, `help`                                                             | Show Help text along with the list of supported flags.                                                                                                                                                                                                                                                                     |
| `v`, `version`    | No       | `false` | No     | `v`, `version`                                                          | Whether to display application version and then immediately exit application.                                                                                                                                                                                                                                              |
| `ll`, `log-level` | No       | `info`  | No     | `disabled`, `panic`, `fatal`, `error`, `warn`, `info`, `debug`, `trace` | Log message priority filter. Log messages with a lower level are ignored.                                                                                                                                                                                                                                                  |
| `p`, `port`       | No       | `443`   | No     | *positive whole number between 1-65535, inclusive*                      | TCP port of the remote ESXi host or vCenter instance. This is usually 443 (HTTPS).                                                                                                                                                                                                                                         |
| `t`, `timeout`    | No       | `10`    | No     | *positive whole number of seconds*                                      | Timeout value in seconds allowed before a plugin execution attempt is abandoned and an error returned.                                                                                                                                                                                                                     |
| `s`, `server`     | **Yes**  |         | No     | *fully-qualified domain name or IP Address*                             | The fully-qualified domain name or IP Address of the remote ESXi host or vCenter instance.                                                                                                                                                                                                                                 |
| `u`, `username`   | **Yes**  |         | No     | *valid username*                                                        | Username with permission to access specified ESXi host or vCenter instance.                                                                                                                                                                                                                                                |
| `pw`, `password`  | **Yes**  |         | No     | *valid password*                                                        | Password used to login to ESXi host or vCenter instance.                                                                                                                                                                                                                                                                   |
| `domain`          | No       |         | No     | *valid user domain*                                                     | (Optional) domain for user account used to login to ESXi host or vCenter instance.                                                                                                                                                                                                                                         |
| `trust-cert`      | No       | `false` | No     | `true`, `false`                                                         | Whether the certificate should be trusted as-is without validation. WARNING: TLS is susceptible to man-in-the-middle attacks if enabling this option.                                                                                                                                                                      |
| `include-rp`      | No       |         | No     | *comma-separated list of resource pool names*                           | Specifies a comma-separated list of Resource Pools that should be exclusively used when evaluating VMs. Specifying this option will also exclude any VMs from evaluation that are *outside* of a Resource Pool. This option is incompatible with specifying a list of Resource Pools to ignore or exclude from evaluation. |
| `exclude-rp`      | No       |         | No     | *comma-separated list of resource pool names*                           | Specifies a comma-separated list of Resource Pools that should be ignored when evaluating VMs. This option is incompatible with specifying a list of Resource Pools to include for evaluation.                                                                                                                             |
| `ignore-vm`       | No       |         | No     | *comma-separated list of (vSphere) virtual machine names*               | Specifies a comma-separated list of VM names that should be ignored or excluded from evaluation.                                                                                                                                                                                                                           |
| `powered-off`     | No       | `false` | No     | `true`, `false`                                                         | Toggles evaluation of powered off VMs in addition to powered on VMs. Evaluation of powered off VMs is disabled by default.                                                                                                                                                                                                 |

#### `check_vmware_vcpus`

| Flag                        | Required | Default | Repeat | Possible                                                                | Description                                                                                                                                                                                                                                                                                                                |
| --------------------------- | -------- | ------- | ------ | ----------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `branding`                  | No       | `false` | No     | `branding`                                                              | Toggles emission of branding details with plugin status details. This output is disabled by default.                                                                                                                                                                                                                       |
| `h`, `help`                 | No       | `false` | No     | `h`, `help`                                                             | Show Help text along with the list of supported flags.                                                                                                                                                                                                                                                                     |
| `v`, `version`              | No       | `false` | No     | `v`, `version`                                                          | Whether to display application version and then immediately exit application.                                                                                                                                                                                                                                              |
| `ll`, `log-level`           | No       | `info`  | No     | `disabled`, `panic`, `fatal`, `error`, `warn`, `info`, `debug`, `trace` | Log message priority filter. Log messages with a lower level are ignored.                                                                                                                                                                                                                                                  |
| `p`, `port`                 | No       | `443`   | No     | *positive whole number between 1-65535, inclusive*                      | TCP port of the remote ESXi host or vCenter instance. This is usually 443 (HTTPS).                                                                                                                                                                                                                                         |
| `t`, `timeout`              | No       | `10`    | No     | *positive whole number of seconds*                                      | Timeout value in seconds allowed before a plugin execution attempt is abandoned and an error returned.                                                                                                                                                                                                                     |
| `s`, `server`               | **Yes**  |         | No     | *fully-qualified domain name or IP Address*                             | The fully-qualified domain name or IP Address of the remote ESXi host or vCenter instance.                                                                                                                                                                                                                                 |
| `u`, `username`             | **Yes**  |         | No     | *valid username*                                                        | Username with permission to access specified ESXi host or vCenter instance.                                                                                                                                                                                                                                                |
| `pw`, `password`            | **Yes**  |         | No     | *valid password*                                                        | Password used to login to ESXi host or vCenter instance.                                                                                                                                                                                                                                                                   |
| `domain`                    | No       |         | No     | *valid user domain*                                                     | (Optional) domain for user account used to login to ESXi host or vCenter instance.                                                                                                                                                                                                                                         |
| `trust-cert`                | No       | `false` | No     | `true`, `false`                                                         | Whether the certificate should be trusted as-is without validation. WARNING: TLS is susceptible to man-in-the-middle attacks if enabling this option.                                                                                                                                                                      |
| `include-rp`                | No       |         | No     | *comma-separated list of resource pool names*                           | Specifies a comma-separated list of Resource Pools that should be exclusively used when evaluating VMs. Specifying this option will also exclude any VMs from evaluation that are *outside* of a Resource Pool. This option is incompatible with specifying a list of Resource Pools to ignore or exclude from evaluation. |
| `exclude-rp`                | No       |         | No     | *comma-separated list of resource pool names*                           | Specifies a comma-separated list of Resource Pools that should be ignored when evaluating VMs. This option is incompatible with specifying a list of Resource Pools to include for evaluation.                                                                                                                             |
| `ignore-vm`                 | No       |         | No     | *comma-separated list of (vSphere) virtual machine names*               | Specifies a comma-separated list of VM names that should be ignored or excluded from evaluation.                                                                                                                                                                                                                           |
| `powered-off`               | No       | `false` | No     | `true`, `false`                                                         | Toggles evaluation of powered off VMs in addition to powered on VMs. Evaluation of powered off VMs is disabled by default.                                                                                                                                                                                                 |
| `vcma`, `vcpus-max-allowed` | **Yes**  | `0`     | No     | *positive whole number of vCPUs*                                        | Specifies the maximum amount of virtual CPUs (as a whole number) that we are allowed to allocate in the target VMware environment.                                                                                                                                                                                         |
| `vc`, `vcpus-critical`      | No       | `100`   | No     | *percentage as positive whole number*                                   | Specifies the percentage of vCPUs allocation (as a whole number) when a CRITICAL threshold is reached.                                                                                                                                                                                                                     |
| `vw`, `vcpus-warning`       | No       | `95`    | No     | *percentage as positive whole number*                                   | Specifies the percentage of vCPUs allocation (as a whole number) when a WARNING threshold is reached.                                                                                                                                                                                                                      |

#### `check_vmware_vhw`

| Flag              | Required | Default | Repeat | Possible                                                                | Description                                                                                                                                                                                                                                                                                                                |
| ----------------- | -------- | ------- | ------ | ----------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `branding`        | No       | `false` | No     | `branding`                                                              | Toggles emission of branding details with plugin status details. This output is disabled by default.                                                                                                                                                                                                                       |
| `h`, `help`       | No       | `false` | No     | `h`, `help`                                                             | Show Help text along with the list of supported flags.                                                                                                                                                                                                                                                                     |
| `v`, `version`    | No       | `false` | No     | `v`, `version`                                                          | Whether to display application version and then immediately exit application.                                                                                                                                                                                                                                              |
| `ll`, `log-level` | No       | `info`  | No     | `disabled`, `panic`, `fatal`, `error`, `warn`, `info`, `debug`, `trace` | Log message priority filter. Log messages with a lower level are ignored.                                                                                                                                                                                                                                                  |
| `p`, `port`       | No       | `443`   | No     | *positive whole number between 1-65535, inclusive*                      | TCP port of the remote ESXi host or vCenter instance. This is usually 443 (HTTPS).                                                                                                                                                                                                                                         |
| `t`, `timeout`    | No       | `10`    | No     | *positive whole number of seconds*                                      | Timeout value in seconds allowed before a plugin execution attempt is abandoned and an error returned.                                                                                                                                                                                                                     |
| `s`, `server`     | **Yes**  |         | No     | *fully-qualified domain name or IP Address*                             | The fully-qualified domain name or IP Address of the remote ESXi host or vCenter instance.                                                                                                                                                                                                                                 |
| `u`, `username`   | **Yes**  |         | No     | *valid username*                                                        | Username with permission to access specified ESXi host or vCenter instance.                                                                                                                                                                                                                                                |
| `pw`, `password`  | **Yes**  |         | No     | *valid password*                                                        | Password used to login to ESXi host or vCenter instance.                                                                                                                                                                                                                                                                   |
| `domain`          | No       |         | No     | *valid user domain*                                                     | (Optional) domain for user account used to login to ESXi host or vCenter instance.                                                                                                                                                                                                                                         |
| `trust-cert`      | No       | `false` | No     | `true`, `false`                                                         | Whether the certificate should be trusted as-is without validation. WARNING: TLS is susceptible to man-in-the-middle attacks if enabling this option.                                                                                                                                                                      |
| `include-rp`      | No       |         | No     | *comma-separated list of resource pool names*                           | Specifies a comma-separated list of Resource Pools that should be exclusively used when evaluating VMs. Specifying this option will also exclude any VMs from evaluation that are *outside* of a Resource Pool. This option is incompatible with specifying a list of Resource Pools to ignore or exclude from evaluation. |
| `exclude-rp`      | No       |         | No     | *comma-separated list of resource pool names*                           | Specifies a comma-separated list of Resource Pools that should be ignored when evaluating VMs. This option is incompatible with specifying a list of Resource Pools to include for evaluation.                                                                                                                             |
| `ignore-vm`       | No       |         | No     | *comma-separated list of (vSphere) virtual machine names*               | Specifies a comma-separated list of VM names that should be ignored or excluded from evaluation.                                                                                                                                                                                                                           |
| `powered-off`     | No       | `false` | No     | `true`, `false`                                                         | Toggles evaluation of powered off VMs in addition to powered on VMs. Evaluation of powered off VMs is disabled by default.                                                                                                                                                                                                 |

#### `check_vmware_hs2ds2vms`

| Flag                 | Required  | Default | Repeat | Possible                                                                | Description                                                                                                                                                                                                                                                                                                                |
| -------------------- | --------- | ------- | ------ | ----------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `branding`           | No        | `false` | No     | `branding`                                                              | Toggles emission of branding details with plugin status details. This output is disabled by default.                                                                                                                                                                                                                       |
| `h`, `help`          | No        | `false` | No     | `h`, `help`                                                             | Show Help text along with the list of supported flags.                                                                                                                                                                                                                                                                     |
| `v`, `version`       | No        | `false` | No     | `v`, `version`                                                          | Whether to display application version and then immediately exit application.                                                                                                                                                                                                                                              |
| `ll`, `log-level`    | No        | `info`  | No     | `disabled`, `panic`, `fatal`, `error`, `warn`, `info`, `debug`, `trace` | Log message priority filter. Log messages with a lower level are ignored.                                                                                                                                                                                                                                                  |
| `p`, `port`          | No        | `443`   | No     | *positive whole number between 1-65535, inclusive*                      | TCP port of the remote ESXi host or vCenter instance. This is usually 443 (HTTPS).                                                                                                                                                                                                                                         |
| `t`, `timeout`       | No        | `10`    | No     | *positive whole number of seconds*                                      | Timeout value in seconds allowed before a plugin execution attempt is abandoned and an error returned.                                                                                                                                                                                                                     |
| `s`, `server`        | **Yes**   |         | No     | *fully-qualified domain name or IP Address*                             | The fully-qualified domain name or IP Address of the remote ESXi host or vCenter instance.                                                                                                                                                                                                                                 |
| `u`, `username`      | **Yes**   |         | No     | *valid username*                                                        | Username with permission to access specified ESXi host or vCenter instance.                                                                                                                                                                                                                                                |
| `pw`, `password`     | **Yes**   |         | No     | *valid password*                                                        | Password used to login to ESXi host or vCenter instance.                                                                                                                                                                                                                                                                   |
| `domain`             | No        |         | No     | *valid user domain*                                                     | (Optional) domain for user account used to login to ESXi host or vCenter instance.                                                                                                                                                                                                                                         |
| `trust-cert`         | No        | `false` | No     | `true`, `false`                                                         | Whether the certificate should be trusted as-is without validation. WARNING: TLS is susceptible to man-in-the-middle attacks if enabling this option.                                                                                                                                                                      |
| `include-rp`         | No        |         | No     | *comma-separated list of resource pool names*                           | Specifies a comma-separated list of Resource Pools that should be exclusively used when evaluating VMs. Specifying this option will also exclude any VMs from evaluation that are *outside* of a Resource Pool. This option is incompatible with specifying a list of Resource Pools to ignore or exclude from evaluation. |
| `exclude-rp`         | No        |         | No     | *comma-separated list of resource pool names*                           | Specifies a comma-separated list of Resource Pools that should be ignored when evaluating VMs. This option is incompatible with specifying a list of Resource Pools to include for evaluation.                                                                                                                             |
| `ignore-vm`          | No        |         | No     | *comma-separated list of (vSphere) virtual machine names*               | Specifies a comma-separated list of VM names that should be ignored or excluded from evaluation.                                                                                                                                                                                                                           |
| `ignore-ds`          | No        |         | No     | *comma-separated list of (vSphere) datastore names*                     | Specifies a comma-separated list of Datastore names that should be ignored or excluded from evaluation.                                                                                                                                                                                                                    |
| `powered-off`        | No        | `false` | No     | `true`, `false`                                                         | Toggles evaluation of powered off VMs in addition to powered on VMs. Evaluation of powered off VMs is disabled by default.                                                                                                                                                                                                 |
| `ca-name`            | **Maybe** |         | No     | *valid Custom Attribute name*                                           | Custom Attribute name for host ESXi systems and datastores. Optional if specifying resource-specific custom attribute names.                                                                                                                                                                                               |
| `ca-prefix-sep`      | **Maybe** |         | No     | *valid Custom Attribute prefix separator character*                     | Custom Attribute prefix separator for host ESXi systems and datastores. Skip if using Custom Attribute values as-is for comparison, otherwise optional if specifying resource-specific custom attribute prefix separator, or using the default separator.                                                                  |
| `ignore-missing-ca`  | No        | `false` | No     | `true`, `false`                                                         | Toggles how missing specified Custom Attributes will be handled. By default, ESXi hosts and datastores missing the Custom Attribute are treated as an error condition.                                                                                                                                                     |
| `host-ca-name`       | **Maybe** |         | No     | *valid Custom Attribute name*                                           | Custom Attribute name specific to host ESXi systems. Optional if specifying shared custom attribute flag.                                                                                                                                                                                                                  |
| `host-ca-prefix-sep` | **Maybe** |         | No     | *valid Custom Attribute prefix separator character*                     | Custom Attribute prefix separator specific to host ESXi systems. Skip if using Custom Attribute values as-is for comparison, otherwise optional if specifying shared custom attribute prefix separator, or using the default separator.                                                                                    |
| `ds-ca-name`         | **Maybe** |         | No     | *valid Custom Attribute name*                                           | Custom Attribute name specific to datastores. Optional if specifying shared custom attribute flag.                                                                                                                                                                                                                         |
| `ds-ca-prefix-sep`   | **Maybe** |         | No     | *valid Custom Attribute prefix separator character*                     | Custom Attribute prefix separator specific to datastores. Skip if using Custom Attribute values as-is for comparison, otherwise optional if specifying shared custom attribute prefix separator, or using the default separator.                                                                                           |

#### `check_vmware_datastore`

| Flag                        | Required | Default | Repeat | Possible                                                                | Description                                                                                                                                                                                            |
| --------------------------- | -------- | ------- | ------ | ----------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `branding`                  | No       | `false` | No     | `branding`                                                              | Toggles emission of branding details with plugin status details. This output is disabled by default.                                                                                                   |
| `h`, `help`                 | No       | `false` | No     | `h`, `help`                                                             | Show Help text along with the list of supported flags.                                                                                                                                                 |
| `v`, `version`              | No       | `false` | No     | `v`, `version`                                                          | Whether to display application version and then immediately exit application.                                                                                                                          |
| `ll`, `log-level`           | No       | `info`  | No     | `disabled`, `panic`, `fatal`, `error`, `warn`, `info`, `debug`, `trace` | Log message priority filter. Log messages with a lower level are ignored.                                                                                                                              |
| `p`, `port`                 | No       | `443`   | No     | *positive whole number between 1-65535, inclusive*                      | TCP port of the remote ESXi host or vCenter instance. This is usually 443 (HTTPS).                                                                                                                     |
| `t`, `timeout`              | No       | `10`    | No     | *positive whole number of seconds*                                      | Timeout value in seconds allowed before a plugin execution attempt is abandoned and an error returned.                                                                                                 |
| `s`, `server`               | **Yes**  |         | No     | *fully-qualified domain name or IP Address*                             | The fully-qualified domain name or IP Address of the remote ESXi host or vCenter instance.                                                                                                             |
| `u`, `username`             | **Yes**  |         | No     | *valid username*                                                        | Username with permission to access specified ESXi host or vCenter instance.                                                                                                                            |
| `pw`, `password`            | **Yes**  |         | No     | *valid password*                                                        | Password used to login to ESXi host or vCenter instance.                                                                                                                                               |
| `domain`                    | No       |         | No     | *valid user domain*                                                     | (Optional) domain for user account used to login to ESXi host or vCenter instance.                                                                                                                     |
| `trust-cert`                | No       | `false` | No     | `true`, `false`                                                         | Whether the certificate should be trusted as-is without validation. WARNING: TLS is susceptible to man-in-the-middle attacks if enabling this option.                                                  |
| `dc-name`                   | No       |         | No     | *valid vSphere datacenter name*                                         | Specifies the name of a vSphere Datacenter. If not specified, applicable plugins will attempt to use the default datacenter found in the vSphere environment. Not applicable to standalone ESXi hosts. |
| `ds-name`                   | **Yes**  |         | No     | *valid datastore name*                                                  | Datastore name as it is found within the vSphere inventory.                                                                                                                                            |
| `dsuc`, `ds-usage-critical` | No       | `95`    | No     | *percentage as positive whole number*                                   | Specifies the percentage of a datastore's storage usage (as a whole number) when a `CRITICAL` threshold is reached.                                                                                    |
| `dsuw`, `ds-usage-warning`  | No       | `90`    | No     | *percentage as positive whole number*                                   | Specifies the percentage of a datastore's storage usage (as a whole number) when a `WARNING` threshold is reached.                                                                                     |

#### `check_vmware_snapshots_age`

| Flag                 | Required | Default | Repeat | Possible                                                                | Description                                                                                                                                                                                                                                                                                                                |
| -------------------- | -------- | ------- | ------ | ----------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `branding`           | No       | `false` | No     | `branding`                                                              | Toggles emission of branding details with plugin status details. This output is disabled by default.                                                                                                                                                                                                                       |
| `h`, `help`          | No       | `false` | No     | `h`, `help`                                                             | Show Help text along with the list of supported flags.                                                                                                                                                                                                                                                                     |
| `v`, `version`       | No       | `false` | No     | `v`, `version`                                                          | Whether to display application version and then immediately exit application.                                                                                                                                                                                                                                              |
| `ll`, `log-level`    | No       | `info`  | No     | `disabled`, `panic`, `fatal`, `error`, `warn`, `info`, `debug`, `trace` | Log message priority filter. Log messages with a lower level are ignored.                                                                                                                                                                                                                                                  |
| `p`, `port`          | No       | `443`   | No     | *positive whole number between 1-65535, inclusive*                      | TCP port of the remote ESXi host or vCenter instance. This is usually 443 (HTTPS).                                                                                                                                                                                                                                         |
| `t`, `timeout`       | No       | `10`    | No     | *positive whole number of seconds*                                      | Timeout value in seconds allowed before a plugin execution attempt is abandoned and an error returned.                                                                                                                                                                                                                     |
| `s`, `server`        | **Yes**  |         | No     | *fully-qualified domain name or IP Address*                             | The fully-qualified domain name or IP Address of the remote ESXi host or vCenter instance.                                                                                                                                                                                                                                 |
| `u`, `username`      | **Yes**  |         | No     | *valid username*                                                        | Username with permission to access specified ESXi host or vCenter instance.                                                                                                                                                                                                                                                |
| `pw`, `password`     | **Yes**  |         | No     | *valid password*                                                        | Password used to login to ESXi host or vCenter instance.                                                                                                                                                                                                                                                                   |
| `domain`             | No       |         | No     | *valid user domain*                                                     | (Optional) domain for user account used to login to ESXi host or vCenter instance.                                                                                                                                                                                                                                         |
| `trust-cert`         | No       | `false` | No     | `true`, `false`                                                         | Whether the certificate should be trusted as-is without validation. WARNING: TLS is susceptible to man-in-the-middle attacks if enabling this option.                                                                                                                                                                      |
| `include-rp`         | No       |         | No     | *comma-separated list of resource pool names*                           | Specifies a comma-separated list of Resource Pools that should be exclusively used when evaluating VMs. Specifying this option will also exclude any VMs from evaluation that are *outside* of a Resource Pool. This option is incompatible with specifying a list of Resource Pools to ignore or exclude from evaluation. |
| `exclude-rp`         | No       |         | No     | *comma-separated list of resource pool names*                           | Specifies a comma-separated list of Resource Pools that should be ignored when evaluating VMs. This option is incompatible with specifying a list of Resource Pools to include for evaluation.                                                                                                                             |
| `ignore-vm`          | No       |         | No     | *comma-separated list of (vSphere) virtual machine names*               | Specifies a comma-separated list of VM names that should be ignored or excluded from evaluation.                                                                                                                                                                                                                           |
| `ac`, `age-critical` | No       | `2`     | No     | *age in days as positive whole number*                                  | Specifies the age of a snapshot in days when a CRITICAL threshold is reached.                                                                                                                                                                                                                                              |
| `aw`, `age-warning`  | No       | `1`     | No     | *age in days as positive whole number*                                  | Specifies the age of a snapshot in days when a WARNING threshold is reached.                                                                                                                                                                                                                                               |

### Configuration file

Not currently supported. This feature may be added later if there is
sufficient interest.

## Contrib

Example Nagios configuration files are provided in an effort to illustrate
usage of plugins provided by this project. See the [Contrib
README](contrib/README.md) and [directory contents](./contrib/) for details.

## Examples

While entries in this section attempt to provide a brief overview of usage, it
is recommended that you review the provided command definitions and other
Nagios configuration files within the [`contrib`](#contrib) directory for more
complete examples.

See the [configuration options](#configuration-options) section for all
command-line settings supported by this plugin along with descriptions of
each.

### `check_vmware_tools` Nagios plugin

#### CLI invocation

```ShellSession
/usr/lib/nagios/plugins/check_vmware_tools --username SERVICE_ACCOUNT_NAME --password "SERVICE_ACCOUNT_PASSWORD" --server vc1.example.com --exclude-rp "Desktops" --ignore-vm "test1.example.com,redmine.example.com,TESTING-AC,RHEL7-TEST" --trust-cert --log-level info
```

See the [configuration options](#configuration-options) section for all
command-line settings supported by this plugin along with descriptions of
each. See the [contrib](#contrib) section for information regarding example
command definitions and Nagios configuration files.

Of note:

- The resource pool named `Desktops` is excluded from evaluation.
  - this results in *all other* resource pools visible to the specified user
    account being used for evaluation
  - this also results in *all* VMs *outside* of a Resource Pool visible to the
    specified user account being used for evaluation
- Multiple Virtual machines (vSphere inventory name, not OS hostname), are
  ignored, regardless of which Resource Pool they are part of.
  - `test1.example.com`
  - `redmine.example.com`
  - `TESTING-AC`
  - `RHEL7-TEST`
- Certificate warnings are ignored.
  - not best practice, but many vCenter instances use self-signed certs per
    various freely available guides
- Logging is enabled at the `info` level.
  - this output is sent to `stderr` by default, which Nagios ignores
  - this output is only seen (at least as of Nagios v3.x) when invoking the
    plugin directly via CLI (often for troubleshooting)

#### Command definition

```shell
# /etc/nagios-plugins/config/vmware-tools.cfg

# Look at all pools, all VMs, do not evaluate any VMs that are powered off.
# This variation of the command is most useful for environments where all VMs
# are monitored equally.
define command{
    command_name    check_vmware_tools
    command_line    /usr/lib/nagios/plugins/check_vmware_tools --server '$HOSTNAME$' --domain '$ARG1$' --username '$ARG2$' --password '$ARG3$'  --trust-cert  --log-level info
    }
```

See the [configuration options](#configuration-options) section for all
command-line settings supported by this plugin along with descriptions of
each. See the [contrib](#contrib) section for information regarding example
command definitions and Nagios configuration files.

### `check_vmware_vcpus` Nagios plugin

#### CLI invocation

```ShellSession
/usr/lib/nagios/plugins/check_vmware_tools --username SERVICE_ACCOUNT_NAME --password "SERVICE_ACCOUNT_PASSWORD" --server vc1.example.com --exclude-rp "Desktops" --ignore-vm "test1.example.com,redmine.example.com,TESTING-AC,RHEL7-TEST" --vcpus-warning 97 --vcpus-critical 100  --vcpus-max-allowed 160 --trust-cert --log-level info
```

See the [configuration options](#configuration-options) section for all
command-line settings supported by this plugin along with descriptions of
each. See the [contrib](#contrib) section for information regarding example
command definitions and Nagios configuration files.

Of note:

- The resource pool named `Desktops` is excluded from evaluation.
  - this results in *all other* resource pools visible to the specified user
    account being used for evaluation
  - this also results in *all* VMs *outside* of a Resource Pool visible to the
    specified user account being used for evaluation
- Multiple Virtual machines (vSphere inventory name, not OS hostname), are
  ignored, regardless of which Resource Pool they are part of.
  - `test1.example.com`
  - `redmine.example.com`
  - `TESTING-AC`
  - `RHEL7-TEST`
- Certificate warnings are ignored.
  - not best practice, but many vCenter instances use self-signed certs per
    various freely available guides
- Logging is enabled at the `info` level.
  - this output is sent to `stderr` by default, which Nagios ignores
  - this output is only seen (at least as of Nagios v3.x) when invoking the
    plugin directly via CLI (often for troubleshooting)

#### Command definition

```shell
# /etc/nagios-plugins/config/vmware-vcpus.cfg

# Look at all pools, all VMs, do not evaluate any VMs that are powered off.
# This variation of the command is most useful for environments where all VMs
# are monitored equally.
define command{
    command_name    check_vmware_vcpus
    command_line    /usr/lib/nagios/plugins/check_vmware_vcpus --server '$HOSTNAME$' --domain '$ARG1$' --username '$ARG2$' --password '$ARG3$' --vcpus-warning '$ARG4$' --vcpus-critical '$ARG5$' --vcpus-max-allowed '$ARG6$' --trust-cert --log-level info
    }
```

See the [configuration options](#configuration-options) section for all
command-line settings supported by this plugin along with descriptions of
each. See the [contrib](#contrib) section for information regarding example
command definitions and Nagios configuration files.

### `check_vmware_vhw` Nagios plugin

#### CLI invocation

```ShellSession
/usr/lib/nagios/plugins/check_vmware_vhw --username SERVICE_ACCOUNT_NAME --password "SERVICE_ACCOUNT_PASSWORD" --server vc1.example.com --exclude-rp "Desktops" --ignore-vm "test1.example.com,redmine.example.com,TESTING-AC,RHEL7-TEST" --trust-cert --log-level info
```

See the [configuration options](#configuration-options) section for all
command-line settings supported by this plugin along with descriptions of
each. See the [contrib](#contrib) section for information regarding example
command definitions and Nagios configuration files.

Of note:

- The resource pool named `Desktops` is excluded from evaluation.
  - this results in *all other* resource pools visible to the specified user
    account being used for evaluation
  - this also results in *all* VMs *outside* of a Resource Pool visible to the
    specified user account being used for evaluation
- Multiple Virtual machines (vSphere inventory name, not OS hostname), are
  ignored, regardless of which Resource Pool they are part of.
  - `test1.example.com`
  - `redmine.example.com`
  - `TESTING-AC`
  - `RHEL7-TEST`
- Certificate warnings are ignored.
  - not best practice, but many vCenter instances use self-signed certs per
    various freely available guides
- Logging is enabled at the `info` level.
  - this output is sent to `stderr` by default, which Nagios ignores
  - this output is only seen (at least as of Nagios v3.x) when invoking the
    plugin directly via CLI (often for troubleshooting)

#### Command definition

```shell
# /etc/nagios-plugins/config/vmware-virtual-hardware.cfg

# Look at all pools, all VMs, do not evaluate any VMs that are powered off.
# This variation of the command is most useful for environments where all VMs
# are monitored equally.
define command{
    command_name    check_vmware_vhw
    command_line    /usr/lib/nagios/plugins/check_vmware_vhw --server '$HOSTNAME$' --domain '$ARG1$' --username '$ARG2$' --password '$ARG3$' --trust-cert --log-level info
    }
```

See the [configuration options](#configuration-options) section for all
command-line settings supported by this plugin along with descriptions of
each. See the [contrib](#contrib) section for information regarding example
command definitions and Nagios configuration files.

### `check_vmware_hs2ds2vms` Nagios plugin

#### CLI invocation

```ShellSession
/usr/lib/nagios/plugins/check_vmware_vhw --username SERVICE_ACCOUNT_NAME --password "SERVICE_ACCOUNT_PASSWORD" --server vc1.example.com --exclude-rp "Desktops" --ignore-vm "test1.example.com,redmine.example.com,TESTING-AC,RHEL7-TEST" --ca-name "Location" --ca-prefix-sep "-" --trust-cert --log-level info
```

See the [configuration options](#configuration-options) section for all
command-line settings supported by this plugin along with descriptions of
each. See the [contrib](#contrib) section for information regarding example
command definitions and Nagios configuration files.

Of note:

- The resource pool named `Desktops` is excluded from evaluation.
  - this results in *all other* resource pools visible to the specified user
    account being used for evaluation
  - this also results in *all* VMs *outside* of a Resource Pool visible to the
    specified user account being used for evaluation
- Multiple Virtual machines (vSphere inventory name, not OS hostname), are
  ignored, regardless of which Resource Pool they are part of.
  - `test1.example.com`
  - `redmine.example.com`
  - `TESTING-AC`
  - `RHEL7-TEST`
- The Custom Attribute named `Location` is used to dynamically build pairs of
  Hosts and Datastores. Any Host or Datastore missing that Custom Attribute is
  reported as an error condition unless the appropriate CLI flag is provided.
  See the [Configuration options](#configuration-options) section for the flag
  name and further details.
- The Custom Attribute prefix separator `-` is provided in order to "split"
  the value found for the Custom Attribute named `Location` into pairs. The
  second value is thrown away, leaving the first to be used as the `Location`
  value for comparison. VMs running on a host with one value have their
  datastores checked for the same value. If a mismatch is found, this is
  assumed to be a `CRITICAL` level event and reported as such.
- Certificate warnings are ignored.
  - not best practice, but many vCenter instances use self-signed certs per
    various freely available guides
- Logging is enabled at the `info` level.
  - this output is sent to `stderr` by default, which Nagios ignores
  - this output is only seen (at least as of Nagios v3.x) when invoking the
    plugin directly via CLI (often for troubleshooting)

#### Command definition

```shell
# /etc/nagios-plugins/config/vmware-host-datastore-vms-pairings.cfg

# Look at all pools, all VMs, do not evaluate any VMs that are powered off.
# Use the same Custom Attribute for hosts and datastores. Use the same Custom
# Attribute prefix separator for hosts and datastores.
#
# This variation of the command is most useful for environments where all VMs
# are monitored equally.
define command{
    command_name   check_vmware_hs2ds2vms
    command_line   /usr/lib/nagios/plugins/check_vmware_hs2ds2vms --server '$HOSTNAME$' --domain '$ARG1$' --username '$ARG2$' --password '$ARG3$' --ca-name '$ARG4$' --ca-prefix-sep '$ARG5$' --trust-cert --log-level info
    }
```

See the [configuration options](#configuration-options) section for all
command-line settings supported by this plugin along with descriptions of
each. See the [contrib](#contrib) section for information regarding example
command definitions and Nagios configuration files.

### `check_vmware_datastore` Nagios plugin

#### CLI invocation

```ShellSession
/usr/lib/nagios/plugins/check_vmware_datastore --username SERVICE_ACCOUNT_NAME --password "SERVICE_ACCOUNT_PASSWORD" --server vc1.example.com --ds-name "HUSVM-DC1-vol6" --ds-usage-warning 95 --ds-usage-critical 97 --trust-cert --log-level info
```

See the [configuration options](#configuration-options) section for all
command-line settings supported by this plugin along with descriptions of
each. See the [contrib](#contrib) section for information regarding example
command definitions and Nagios configuration files.

Of note:

- Certificate warnings are ignored.
  - not best practice, but many vCenter instances use self-signed certs per
    various freely available guides
- Logging is enabled at the `info` level.
  - this output is sent to `stderr` by default, which Nagios ignores
  - this output is only seen (at least as of Nagios v3.x) when invoking the
    plugin directly via CLI (often for troubleshooting)

#### Command definition

```shell
# /etc/nagios-plugins/config/vmware-datastores.cfg

# Look at specific datastore and explicitly provide custom WARNING and
# CRITICAL threshold values.
define command{
    command_name    check_vmware_datastore
    command_line    /usr/lib/nagios/plugins/check_vmware_tools --server '$HOSTNAME$' --domain '$ARG1$' --username '$ARG2$' --password '$ARG3$' --ds-usage-warning '$ARG4$' --ds-usage-critical '$ARG5$' --ds-name '$ARG6$' --trust-cert  --log-level info
    }

```

See the [configuration options](#configuration-options) section for all
command-line settings supported by this plugin along with descriptions of
each. See the [contrib](#contrib) section for information regarding example
command definitions and Nagios configuration files.

### `check_vmware_snapshots_age` Nagios plugin

#### CLI invocation

```ShellSession
/usr/lib/nagios/plugins/check_vmware_snapshots_age --username SERVICE_ACCOUNT_NAME --password "SERVICE_ACCOUNT_PASSWORD" --server vc1.example.com --age-warning 1 --age-critical 2 --trust-cert --log-level info
```

See the [configuration options](#configuration-options) section for all
command-line settings supported by this plugin along with descriptions of
each. See the [contrib](#contrib) section for information regarding example
command definitions and Nagios configuration files.

Of note:

- No Resource Pools are explicitly *included* or *excluded*
  - this results in *all* Resource Pools visible to the specified user account
    being used for evaluation
  - this also results in *all* VMs *outside* of a Resource Pool visible to the
    specified user account being used for evaluation
- Certificate warnings are ignored.
  - not best practice, but many vCenter instances use self-signed certs per
    various freely available guides
- Logging is enabled at the `info` level.
  - this output is sent to `stderr` by default, which Nagios ignores
  - this output is only seen (at least as of Nagios v3.x) when invoking the
    plugin directly via CLI (often for troubleshooting)

#### Command definition

```shell
# /etc/nagios-plugins/config/vmware-snapshots-age.cfg

# Look at all pools, all VMs, do not evaluate any VMs that are powered off.
# This variation of the command is most useful for environments where all VMs
# are monitored equally.
define command{
    command_name    check_vmware_snapshots_age
    command_line    /usr/lib/nagios/plugins/check_vmware_snapshots_age --server '$HOSTNAME$' --domain '$ARG1$' --username '$ARG2$' --password '$ARG3$' --age-warning '$ARG4$' --age-critical '$ARG5$' --trust-cert --log-level info
    }
```

See the [configuration options](#configuration-options) section for all
command-line settings supported by this plugin along with descriptions of
each. See the [contrib](#contrib) section for information regarding example
command definitions and Nagios configuration files.

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

- Related projects
  - <https://github.com/atc0005/send2teams>
  - <https://github.com/atc0005/check-cert>
  - <https://github.com/atc0005/check-mail>
  - <https://github.com/atc0005/check-path>
  - <https://github.com/atc0005/nagios-debug>
  - <https://github.com/atc0005/go-nagios>

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
