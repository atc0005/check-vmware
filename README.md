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

- [Project home](#project-home)
- [Overview](#overview)
  - [Output](#output)
  - [Performance Data](#performance-data)
  - [Optional evaluation](#optional-evaluation)
  - [`check_vmware_tools`](#check_vmware_tools)
  - [`check_vmware_vcpus`](#check_vmware_vcpus)
  - [`check_vmware_vhw`](#check_vmware_vhw)
    - [Homogeneous version check](#homogeneous-version-check)
    - [Outdated-by or threshold range check](#outdated-by-or-threshold-range-check)
    - [Minimum required version check](#minimum-required-version-check)
    - [Default is minimum required version check](#default-is-minimum-required-version-check)
  - [`check_vmware_hs2ds2vms`](#check_vmware_hs2ds2vms)
  - [`check_vmware_datastore`](#check_vmware_datastore)
  - [`check_vmware_snapshots_age`](#check_vmware_snapshots_age)
  - [`check_vmware_snapshots_count`](#check_vmware_snapshots_count)
  - [`check_vmware_snapshots_size`](#check_vmware_snapshots_size)
  - [`check_vmware_rps_memory`](#check_vmware_rps_memory)
  - [`check_vmware_host_memory`](#check_vmware_host_memory)
  - [`check_vmware_host_cpu`](#check_vmware_host_cpu)
  - [`check_vmware_vm_power_uptime`](#check_vmware_vm_power_uptime)
  - [`check_vmware_disk_consolidation`](#check_vmware_disk_consolidation)
  - [`check_vmware_question`](#check_vmware_question)
  - [`check_vmware_alarms`](#check_vmware_alarms)
- [Features](#features)
- [Changelog](#changelog)
- [Requirements](#requirements)
  - [Building source code](#building-source-code)
  - [Running](#running)
- [Installation](#installation)
  - [From source](#from-source)
  - [Using provided binaries](#using-provided-binaries)
    - [Linux](#linux)
    - [Windows](#windows)
    - [Other operating systems](#other-operating-systems)
- [Configuration options](#configuration-options)
  - [Threshold calculations](#threshold-calculations)
    - [`check_vmware_tools`](#check_vmware_tools-1)
    - [`check_vmware_vcpus`](#check_vmware_vcpus-1)
    - [`check_vmware_vhw`](#check_vmware_vhw-1)
      - [Homogeneous version check](#homogeneous-version-check-1)
      - [Outdated-by or threshold range check](#outdated-by-or-threshold-range-check-1)
      - [Minimum required version check](#minimum-required-version-check-1)
      - [Default is minimum required version check](#default-is-minimum-required-version-check-1)
    - [`check_vmware_hs2ds2vms`](#check_vmware_hs2ds2vms-1)
    - [`check_vmware_datastore`](#check_vmware_datastore-1)
    - [`check_vmware_snapshots_age`](#check_vmware_snapshots_age-1)
    - [`check_vmware_snapshots_count`](#check_vmware_snapshots_count-1)
    - [`check_vmware_snapshots_size`](#check_vmware_snapshots_size-1)
    - [`check_vmware_rps_memory`](#check_vmware_rps_memory-1)
    - [`check_vmware_host_memory`](#check_vmware_host_memory-1)
    - [`check_vmware_host_cpu`](#check_vmware_host_cpu-1)
    - [`check_vmware_vm_power_uptime`](#check_vmware_vm_power_uptime-1)
    - [`check_vmware_disk_consolidation`](#check_vmware_disk_consolidation-1)
    - [`check_vmware_question`](#check_vmware_question-1)
    - [`check_vmware_alarms`](#check_vmware_alarms-1)
  - [Command-line arguments](#command-line-arguments)
    - [`check_vmware_tools`](#check_vmware_tools-2)
    - [`check_vmware_vcpus`](#check_vmware_vcpus-2)
    - [`check_vmware_vhw`](#check_vmware_vhw-2)
    - [`check_vmware_hs2ds2vms`](#check_vmware_hs2ds2vms-2)
    - [`check_vmware_datastore`](#check_vmware_datastore-2)
    - [`check_vmware_snapshots_age`](#check_vmware_snapshots_age-2)
    - [`check_vmware_snapshots_count`](#check_vmware_snapshots_count-2)
    - [`check_vmware_snapshots_size`](#check_vmware_snapshots_size-2)
    - [`check_vmware_rps_memory`](#check_vmware_rps_memory-2)
    - [`check_vmware_host_memory`](#check_vmware_host_memory-2)
    - [`check_vmware_host_cpu`](#check_vmware_host_cpu-2)
    - [`check_vmware_vm_power_uptime`](#check_vmware_vm_power_uptime-2)
    - [`check_vmware_disk_consolidation`](#check_vmware_disk_consolidation-2)
    - [`check_vmware_question`](#check_vmware_question-2)
    - [`check_vmware_alarms`](#check_vmware_alarms-2)
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
    - [Homogeneous version check](#homogeneous-version-check-2)
      - [CLI invocation](#cli-invocation-2)
      - [Command definition](#command-definition-2)
    - [Outdated-by or threshold range check](#outdated-by-or-threshold-range-check-2)
      - [CLI invocation](#cli-invocation-3)
      - [Command definition](#command-definition-3)
    - [Minimum required version check](#minimum-required-version-check-2)
      - [CLI invocation](#cli-invocation-4)
      - [Command definition](#command-definition-4)
    - [Default is minimum required version check](#default-is-minimum-required-version-check-2)
      - [CLI invocation](#cli-invocation-5)
      - [Command definition](#command-definition-5)
  - [`check_vmware_hs2ds2vms` Nagios plugin](#check_vmware_hs2ds2vms-nagios-plugin)
    - [CLI invocation](#cli-invocation-6)
    - [Command definition](#command-definition-6)
  - [`check_vmware_datastore` Nagios plugin](#check_vmware_datastore-nagios-plugin)
    - [CLI invocation](#cli-invocation-7)
    - [Command definition](#command-definition-7)
  - [`check_vmware_snapshots_age` Nagios plugin](#check_vmware_snapshots_age-nagios-plugin)
    - [CLI invocation](#cli-invocation-8)
    - [Command definition](#command-definition-8)
  - [`check_vmware_snapshots_count` Nagios plugin](#check_vmware_snapshots_count-nagios-plugin)
    - [CLI invocation](#cli-invocation-9)
    - [Command definition](#command-definition-9)
  - [`check_vmware_snapshots_size` Nagios plugin](#check_vmware_snapshots_size-nagios-plugin)
    - [CLI invocation](#cli-invocation-10)
    - [Command definition](#command-definition-10)
  - [`check_vmware_rps_memory` Nagios plugin](#check_vmware_rps_memory-nagios-plugin)
    - [CLI invocation](#cli-invocation-11)
    - [Command definition](#command-definition-11)
  - [`check_vmware_host_memory` Nagios plugin](#check_vmware_host_memory-nagios-plugin)
    - [CLI invocation](#cli-invocation-12)
    - [Command definition](#command-definition-12)
  - [`check_vmware_host_cpu` Nagios plugin](#check_vmware_host_cpu-nagios-plugin)
    - [CLI invocation](#cli-invocation-13)
    - [Command definition](#command-definition-13)
  - [`check_vmware_vm_power_uptime` Nagios plugin](#check_vmware_vm_power_uptime-nagios-plugin)
    - [CLI invocation](#cli-invocation-14)
    - [Command definition](#command-definition-14)
  - [`check_vmware_disk_consolidation` Nagios plugin](#check_vmware_disk_consolidation-nagios-plugin)
    - [CLI invocation](#cli-invocation-15)
    - [Command definition](#command-definition-15)
  - [`check_vmware_question` Nagios plugin](#check_vmware_question-nagios-plugin)
    - [CLI invocation](#cli-invocation-16)
    - [Command definition](#command-definition-16)
  - [`check_vmware_alarms` Nagios plugin](#check_vmware_alarms-nagios-plugin)
    - [CLI invocation](#cli-invocation-17)
    - [Command definition](#command-definition-17)
- [License](#license)
- [References](#references)

## Project home

See [our GitHub repo][repo-url] for the latest code, to file an issue or
submit improvements for review and potential inclusion into the project.

Just to be 100% clear: this project is not affiliated with or endorsed by
VMware, Inc.

## Overview

This repo contains various tools and plugins used to monitor/validate VMware
environments.

| Plugin or Tool Name               | Description                                                                         |
| --------------------------------- | ----------------------------------------------------------------------------------- |
| `check_vmware_tools`              | Nagios plugin used to monitor VMware Tools installations.                           |
| `check_vmware_vcpus`              | Nagios plugin used to monitor allocation of virtual CPUs (vCPUs).                   |
| `check_vmware_vhw`                | Nagios plugin used to monitor virtual hardware versions.                            |
| `check_vmware_hs2ds2vms`          | Nagios plugin used to monitor host/datastore/vm pairings.                           |
| `check_vmware_datastore`          | Nagios plugin used to monitor datastore usage.                                      |
| `check_vmware_snapshots_age`      | Nagios plugin used to monitor the age of Virtual Machine snapshots.                 |
| `check_vmware_snapshots_count`    | Nagios plugin used to monitor the count of Virtual Machine snapshots.               |
| `check_vmware_snapshots_size`     | Nagios plugin used to monitor the **cumulative** size of Virtual Machine snapshots. |
| `check_vmware_rps_memory`         | Nagios plugin used to monitor memory usage across Resource Pools.                   |
| `check_vmware_host_memory`        | Nagios plugin used to monitor memory usage for a specific ESXi host system.         |
| `check_vmware_host_cpu`           | Nagios plugin used to monitor CPU usage for a specific ESXi host system.            |
| `check_vmware_vm_power_uptime`    | Nagios plugin used to monitor VM power cycle uptime.                                |
| `check_vmware_disk_consolidation` | Nagios plugin used to monitor VM disk consolidation status.                         |
| `check_vmware_question`           | Nagios plugin used to monitor VM interactive question status.                       |
| `check_vmware_alarms`             | Nagios plugin used to monitor for Triggered Alarms in one or more datacenters.      |

### Output

The output for these plugins is designed to provide the one-line summary
needed by Nagios for quick identification of a problem while providing longer,
more detailed information for display within the web UI, use in email and
Teams notifications
([atc0005/send2teams](https://github.com/atc0005/send2teams)).

By default, output intended for processing by Nagios is sent to `stdout` and
output intended for troubleshooting by the sysadmin is sent to `stderr`.

For some monitoring systems or addons (e.g., Icinga Web 2), the `stderr`
output is mixed in with the `stdout` output (GH-314) in the web UI for the
service check. This may add visual noise when viewing the service check
output. For those cases, you may wish to explicitly disable the output to
`stderr` via the `--log-level "disabled"` CLI flag.

If this impacts you, please [provide feedback
here](https://github.com/atc0005/check-vmware/discussions/323). Future
releases of this project may modify plugins to not emit to `stderr` by default
or the example command definitions may be updated to specify the `--log-level
"disabled"` CLI flag.

### Performance Data

Adding support for Performance Data / Metrics to plugins in this project is an
ongoing effort.

Consult the table below (and the GH issues listed) for the metrics implemented
thus far, the applicable GH issue for pending work and the [Add Performance
Data / Metrics support](https://github.com/atc0005/check-vmware/projects/1)
project board for a quick overview of work.

Please provide feedback on the applicable issue(s) if you have any, good or
bad.

| Plugin                            | Emitted Performance Data / Metrics                                                                                                                                                                          |
| --------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `check_vmware_tools`              | `time`, `vms`, `vms_excluded_by_name`, `vms_excluded_by_power_state`, `vms_with_tools_issues`, `vms_without_tools_issues`, `resource_pools_excluded`, `resource_pools_included`, `resource_pools_evaluated` |
| `check_vmware_vcpus`              | TBD. See GH-355 for details.                                                                                                                                                                                |
| `check_vmware_vhw`                | TBD. See GH-356 for details.                                                                                                                                                                                |
| `check_vmware_hs2ds2vms`          | TBD. See GH-348 for details.                                                                                                                                                                                |
| `check_vmware_datastore`          | `time`, `datastore_usage`, `datastore_storage_remaining`, `vms`, `vms_powered_on`, `vms_powered_off`                                                                                                        |
| `check_vmware_snapshots_age`      | TBD. See GH-351 for details.                                                                                                                                                                                |
| `check_vmware_snapshots_count`    | TBD. See GH-352 for details.                                                                                                                                                                                |
| `check_vmware_snapshots_size`     | TBD. See GH-353 for details.                                                                                                                                                                                |
| `check_vmware_rps_memory`         | TBD. See GH-350 for details.                                                                                                                                                                                |
| `check_vmware_host_memory`        | TBD. See GH-347 for details.                                                                                                                                                                                |
| `check_vmware_host_cpu`           | TBD. See GH-346 for details.                                                                                                                                                                                |
| `check_vmware_vm_power_uptime`    | TBD. See GH-357 for details.                                                                                                                                                                                |
| `check_vmware_disk_consolidation` | TBD. See GH-345 for details.                                                                                                                                                                                |
| `check_vmware_question`           | TBD. See GH-349 for details.                                                                                                                                                                                |
| `check_vmware_alarms`             | TBD. See GH-344 for details.                                                                                                                                                                                |

### Optional evaluation

Some plugins provide optional support to limit evaluation of VMs to specific
Resource Pools (explicitly including or excluding) and power states (on or
off). Other plugins support similar filtering options (e.g., `Acknowledged`
state of Triggered Alarms). See the [configuration
options](#configuration-options), [examples](#examples) and
[contrib](#contrib) sections for more information.

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

This plugin supports four monitoring modes:

1. Homogeneous version check
1. Outdated-by or threshold range check
1. Minimum required version check
1. Default is minimum required version check

#### Homogeneous version check

This monitoring mode applies an automatic baseline of "highest version
discovered" across evaluated VMs. Any VMs with a hardware version not at that
highest version are flagged as problematic.

Instead of trying to determine how far behind each VM is from the newest
version, this monitoring mode assumes that any deviation is a `WARNING` state.

#### Outdated-by or threshold range check

This mode applies the standard WARNING and CRITICAL level threshold checks to
determine the current plugin state. Any VM with virtual hardware older than
the specified thresholds triggers the associated state. This mode is useful
for catching VMs with outdated hardware outside of an acceptable range.

The highest version used as a baseline for comparison is provided using the
same logic as provided by the "homogeneous" version check: latest visible
hardware version.

#### Minimum required version check

This mode requires that all hardware versions match or exceed the specified
minimum hardware version. This monitoring mode assumes that any deviation is
considered a `CRITICAL` state.

#### Default is minimum required version check

This mode requires that all hardware versions match or exceed the host or
cluster default hardware version. This monitoring mode assumes that any
deviation is considered a `WARNING` state.

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
*only* powered on VMs by default, please [share it
here](https://github.com/atc0005/check-vmware/discussions/177) providing some
details for your use-case. In our environment, I have yet to see a need to
*only* evaluate powered on VMs for old snapshots. For cases where the
snapshots needed to be ignored, we added the VM to the ignore list. We then
relied on datastore usage monitoring to let us know when space was becoming an
issue.

Thresholds for `CRITICAL` and `WARNING` age values have usable defaults, but
may require adjustment for your environment. See the [configuration
options](#configuration-options) section for details.

### `check_vmware_snapshots_count`

Nagios plugin used to monitor the number of snapshots per Virtual Machine.

Monitor the number of snapshots for each Virtual Machine. VMware recommends
using no more than 3 or 4 snapshots per Virtual Machine and only for a limited
duration. A maximum of 32 snapshots per Virtual Machine are supported. See
<https://kb.vmware.com/s/article/1025279> for more information.

The current design of this plugin is to evaluate *all* Virtual Machines,
whether powered off or powered on. If you have a use case for evaluating
*only* powered on VMs by default, please [share it
here](https://github.com/atc0005/check-vmware/discussions/177) providing some
details for your use-case. In our environment, I have yet to see a need to
*only* evaluate powered on VMs for old snapshots. For cases where the
snapshots needed to be ignored, we added the VM to the ignore list. We then
relied on datastore usage monitoring to let us know when space was becoming an
issue.

Thresholds for `CRITICAL` and `WARNING` count values have usable defaults, but
may require adjustment for your environment. See the [configuration
options](#configuration-options) section for details.

### `check_vmware_snapshots_size`

Nagios plugin used to monitor the **cumulative** size of snapshots for each
Virtual Machine.

While individual snapshots are listed, it is the cumulative size for a Virtual
Machine crossing a given size threshold that determines the overall check
result.

The current design of this plugin is to evaluate *all* Virtual Machines,
whether powered off or powered on. If you have a use case for evaluating
*only* powered on VMs by default, please [share it
here](https://github.com/atc0005/check-vmware/discussions/177) providing some
details for your use-case. In our environment, I have yet to see a need to
*only* evaluate powered on VMs for old snapshots. For cases where the
snapshots needed to be ignored, we added the VM to the ignore list. We then
relied on datastore usage monitoring to let us know when space was becoming an
issue.

Thresholds for `CRITICAL` and `WARNING` age values have usable defaults, but
may require adjustment for your environment. See the [configuration
options](#configuration-options) section for details.

### `check_vmware_rps_memory`

Nagios plugin used to monitor memory usage across Resource Pools.

In addition to reporting memory usage for each Resource Pool, this plugin also
reports the ten most recently booted VMs along with their memory usage. This
is intended to help spot which VM is responsible for a state change alert.

Thresholds for `CRITICAL` and `WARNING` memory usage have usable defaults, but
max memory usage is required before this plugin can be used. See the
[configuration options](#configuration-options) section for details.

### `check_vmware_host_memory`

Nagios plugin used to monitor ESXi host memory.

In addition to reporting current host memory usage, this plugin also reports
which VMs are on the host (running or not), how much memory each VM is using
as a fixed value and as a percentage of the host's total memory.

Thresholds for `CRITICAL` and `WARNING` memory usage have usable defaults, but
max memory usage is required before this plugin can be used. See the
[configuration options](#configuration-options) section for details.

### `check_vmware_host_cpu`

Nagios plugin used to monitor ESXi host CPU usage.

In addition to reporting current host CPU usage, this plugin also reports
which VMs are on the host (running or not), how much CPU each VM is using
as a fixed value and as a percentage of the host's total CPU capacity.

Thresholds for `CRITICAL` and `WARNING` CPU usage have usable defaults, but
may require adjustment for your environment. See the [configuration
options](#configuration-options) section for details.

### `check_vmware_vm_power_uptime`

Nagios plugin used to monitor Virtual Machine (power cycle) uptime.

This is essentially the time since the VM was last powered off and then back
on (e.g., for a snapshot).

In addition to reporting current power cycle uptime, this plugin also reports:

- which VMs have crossed thresholds (if any) and the uptime for each
- which VMs have yet to cross thresholds (only if there are not any which
  have) and the uptime for each
- the ten most recently booted VMs

Thresholds for `CRITICAL` and `WARNING` CPU usage have usable defaults, but
may require adjustment for your environment. See the [configuration
options](#configuration-options) section for details.

### `check_vmware_disk_consolidation`

Nagios plugin used to monitor Virtual Machine disk consolidation status.

The status of this property indicates whether one or more disks for a Virtual
Machine require consolidation. This can happen when a snapshot is deleted, but
its associated disk is not committed back to the base disk. This situation can
cause backup failures and performance issues.

By default, this plugin does not trigger a state reload for each Virtual
Machine that it evaluates, instead evaluating the disk consolidation status as
currently reflected in the vSphere environment. The state data appears to only
be updated during vMotion and Fault Tolerant related methods, when a VM is
first added to inventory or when manually reloaded via the vSphere web UI. If
not refreshed by one of these tasks or a custom job configured on the cluster
the consolidation status may be stale.

You can work around this potentially stale state by specifying a
`--trigger-reload` flag for this plugin. This flag enables a state reload for
each evaluated Virtual Machine. This reload will refresh state data for the
Virtual Machine to ensure that the disk consolidation status reflects the
actual state of the VM. This option does not come without a cost however.

Due to the time required for each reload operation to complete, this plugin
can require a much longer timeout value than other plugins which only evaluate
(and not refresh) existing state data for vSphere objects. You should
configure the `--timeout` value for this plugin accordingly and also configure
the timeout settings in your monitoring system (e.g., `service_check_timeout`
within `nagios.cfg` for Nagios) to permit longer plugin execution times.

Instead of enabling this flag, you may wish to schedule a job on the cluster
or an "admin box" that handles the reload/refresh of each Virtual Machine.
This will be significantly faster than evaluating the state of each VM every
time the associated service check executes and depending on the frequency of
the job should be "fresh enough" to allow this plugin to accurately detect
disk consolidation needs.

### `check_vmware_question`

Nagios plugin used to monitor whether a Virtual Machine is blocked from
execution due to one or more Virtual Machines requiring an interactive
response.

This plugin monitors the `question` property of evaluated Virtual Machines.
The status of this property indicates whether an interactive question is
blocking the virtual machine's execution. While a Virtual Machine is in this
state it is not available for normal use.

### `check_vmware_alarms`

Nagios plugin used to monitor for Triggered Alarms in one or more datacenters.

- Explicit exclusions take priority over either implicit or explicit
  inclusions.
- All filtering is currently applied in batches/bulk.

It helps to think of the process working this way for each filter in the
"pipeline":

1. Explicit inclusions are applied, marking matching triggered alarms as
   explicitly included and non-matches as *implicitly* excluded
1. Explicit exclusions are applied, marking matching triggered alarms as
   explicitly excluded, permanently "dropping" the triggered alarm from
   further evaluation
1. After all filters have finished processing, any triggered alarms marked as
   excluded (implicit or explicit) are removed from final evaluation (i.e.,
   ignored and not reported as a problem).

Filtering is available for explicitly *including* or *excluding* based on:

- `Acknowledged` status
- [Managed Entity type][vsphere-managed-object-reference] (e.g.,
  `Datastore`, `VirtualMachine`) associated with the Triggered Alarm
- Inventory object `name` (e.g., `node1.example.com`, `vc1.example.com`)
  associated with the Triggered Alarm
- Alarm `Name` field substring match
- Alarm `Description` field substring match
- Triggered Alarm `Status` (e.g., `red`, `yellow`, `gray`)
- `Resource Pool` for the [Managed Entity
  type][vsphere-managed-object-reference] (e.g., `ResourcePool`,
  `VirtualMachine`) associated with the Triggered Alarm

## Features

- Multiple plugins for monitoring VMware vSphere environments (standalone ESXi
  hosts or vCenter instances) for select (or all) Resource Pools.
  - VMware Tools
  - Virtual CPU allocations
  - Virtual hardware versions (multiple modes)
    - homogeneous version check
    - outdated-by or threshold range check
    - minimum required version check
    - default is minimum required version check
  - Host/Datastore/Virtual Machine pairings (using provided Custom Attribute)
  - Datastore usage
  - Snapshots age
  - Snapshots count
  - Snapshots size
  - Resource Pools: Memory usage
  - Host Memory usage
  - Host CPU usage
  - Virtual Machine (power cycle) uptime
  - Virtual Machine disk consolidation status
    - with optional forced refresh of Virtual Machine state data
  - Virtual Machine interactive question status
  - Triggered Alarms in one or more datacenters

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
  - dependent on current upstream `vmware/govmomi` library
- GCC
  - if building with custom options (as the provided `Makefile` does)
- `make`
  - if using the provided `Makefile`

### Running

- Windows 10
- Ubuntu Linux 16.04+

## Installation

### From source

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
     - `go build -mod=vendor ./cmd/check_vmware_snapshots_count/`
     - `go build -mod=vendor ./cmd/check_vmware_snapshots_size/`
     - `go build -mod=vendor ./cmd/check_vmware_rps_memory/`
     - `go build -mod=vendor ./cmd/check_vmware_host_memory/`
     - `go build -mod=vendor ./cmd/check_vmware_host_cpu/`
     - `go build -mod=vendor ./cmd/check_vmware_vm_power_uptime/`
     - `go build -mod=vendor ./cmd/check_vmware_disk_consolidation/`
     - `go build -mod=vendor ./cmd/check_vmware_question/`
     - `go build -mod=vendor ./cmd/check_vmware_alarms/`
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
     - look in `/tmp/check-vmware/release_assets/check_vmware_snapshots_count/`
     - look in `/tmp/check-vmware/release_assets/check_vmware_snapshots_size/`
     - look in `/tmp/check-vmware/release_assets/check_vmware_rps_memory/`
     - look in `/tmp/check-vmware/release_assets/check_vmware_host_memory/`
     - look in `/tmp/check-vmware/release_assets/check_vmware_host_cpu/`
     - look in `/tmp/check-vmware/release_assets/check_vmware_vm_power_uptime/`
     - look in `/tmp/check-vmware/release_assets/check_vmware_disk_consolidation/`
     - look in `/tmp/check-vmware/release_assets/check_vmware_question/`
     - look in `/tmp/check-vmware/release_assets/check_vmware_alarms/`
   - if using `go build`
     - look in `/tmp/check-vmware/`
1. Review [configuration options](#configuration-options),
   [`examples`](#examples) and [`contrib`](#contrib) sections usage details.

### Using provided binaries

#### Linux

1. Download plugins from the [Latest
   release](https://github.com/atc0005/check-vmware/releases/latest) that you
   are interested in
1. Review [configuration options](#configuration-options),
   [`examples`](#examples) and [`contrib`](#contrib) sections usage details.

#### Windows

Note: As of the `v0.13.0` release, precompiled Windows binaries are no longer
provided. This change was made primarily due to the lengthy build times
required and the perception that most users of this project would not benefit
from having them. If you *do* use Windows binaries or would like to (e.g., on
a Windows system within a restricted environment that has access to your
vSphere cluster or hosts), please provide feedback on
[GH-178](https://github.com/atc0005/check-vmware/discussions/178).

#### Other operating systems

As of the `v0.13.0` release, only Linux precompiled binaries are provided. If
you would benefit from precompiled binaries for other platforms, please let us
know by opening a new issue or responding to an existing issue with an
up-vote. See <https://golang.org/doc/install/source> for a list of supported
architectures and operating systems.

## Configuration options

### Threshold calculations

#### `check_vmware_tools`

This plugin evaluates two fields from the [GuestInfo Data
Object][vsphere-guestinfo-data-object] vSphere API:

- `toolsRunningStatus`
- `toolsVersionStatus2`

The overall state of the service check is determined based on these fields,
the power state of an evaluated VM and whether the `powered-off` flag has been
specified:

- If it has not, then powered off VMs are ignored.
- If it has, then powered off VMs are evaluated for combinations of field
  values that appear to be relevant.
  - For example, this plugin does not consider VMware Tools with a "not
    running" status to be a problem if the Virtual Machine is powered off.

To simplify this table, most entries assume that the `powered-off` flag has
been specified. If it is not specified in your Nagios instance, then powered
off Virtual Machines will be ignored; the details of this table will not apply
to those Virtual Machines.

| VM Power State | `powered-off` flag | API Field Name        | API Field Value          | Nagios State | Description                                                                                                              |
| -------------- | ------------------ | --------------------- | ------------------------ | ------------ | ------------------------------------------------------------------------------------------------------------------------ |
| Powered Off    | Yes                | `toolsRunningStatus`  | `guestToolsNotRunning`   | `OK`         | Virtual Machine is not running, so VMware Tools is not expected to run either.                                           |
| Powered On     | N/A                | `toolsRunningStatus`  | `guestToolsNotRunning`   | `CRITICAL`   | VMware Tools (or `open-vm-tools`) not currently running. It likely crashed or was terminated due to low memory scenario. |
| N/A            | Yes                | `toolsVersionStatus2` | `guestToolsNotInstalled` | `CRITICAL`   | VMware Tools is not installed.                                                                                           |
| N/A            | Yes                | `toolsVersionStatus2` | `guestToolsCurrent`      | `OK`         | Ideal state, no problems with VMware Tools (or `open-vm-tools`) detected.                                                |
| N/A            | Yes                | `toolsVersionStatus2` | `guestToolsUnmanaged`    | `OK`         | *Assumed* to be an `OK` state; VMware Tools is installed, but it is not managed by VMware (e.g., `open-vm-tools`).       |
| N/A            | Yes                | `toolsVersionStatus2` | `guestToolsTooOld`       | `CRITICAL`   | VMware Tools is installed, but the version is too old.                                                                   |
| N/A            | Yes                | `toolsVersionStatus2` | `guestToolsSupportedOld` | `WARNING`    | VMware Tools is installed, supported, but a newer version is available.                                                  |
| N/A            | Yes                | `toolsVersionStatus2` | `guestToolsNeedUpgrade`  | `WARNING`    | VMware Tools is installed, but the version is not current. Assumed to be roughly equivalent to `guestToolsSupportedOld`. |
| N/A            | Yes                | `toolsVersionStatus2` | `guestToolsSupportedNew` | `OK`         | VMware Tools is installed, supported, and newer than the version available on the host.                                  |
| N/A            | Yes                | `toolsVersionStatus2` | `guestToolsTooNew`       | `CRITICAL`   | VMware Tools is installed, and the version is known to be too new to work correctly with this virtual machine.           |
| N/A            | Yes                | `toolsRunningStatus2` | `guestToolsBlacklisted`  | `CRITICAL`   | VMware Tools is installed, but the installed version is known to have a grave bug and should be immediately upgraded.    |
| N/A            | Yes                | `toolsVersionStatus2` | Unknown to this plugin   | `UNKNOWN`    | This field in the vSphere API has been extended and this library hasn't been updated to account for those changes.       |

#### `check_vmware_vcpus`

| Nagios State | Description                                                       |
| ------------ | ----------------------------------------------------------------- |
| `OK`         | Ideal state, vCPU allocations within bounds.                      |
| `WARNING`    | vCPU allocations crossed user-specified threshold for this state. |
| `CRITICAL`   | vCPU allocations crossed user-specified threshold for this state. |

#### `check_vmware_vhw`

This plugin supports multiple modes. Each mode applies slightly different
logic for determining plugin state.

##### Homogeneous version check

| Nagios State | Description                                |
| ------------ | ------------------------------------------ |
| `OK`         | Ideal state, homogenous hardware versions. |
| `WARNING`    | Non-homogenous hardware versions.          |
| `CRITICAL`   | Not used by this monitoring mode.          |

##### Outdated-by or threshold range check

| Nagios State | Description                                                        |
| ------------ | ------------------------------------------------------------------ |
| `OK`         | Ideal state, hardware versions within tolerance.                   |
| `WARNING`    | Hardware versions crossed user-specified threshold for this state. |
| `CRITICAL`   | Hardware versions crossed user-specified threshold for this state. |

##### Minimum required version check

| Nagios State | Description                                                       |
| ------------ | ----------------------------------------------------------------- |
| `OK`         | Ideal state, hardware versions within tolerance.                  |
| `WARNING`    | Not used by this monitoring mode.                                 |
| `CRITICAL`   | Hardware versions older than the minimum specified value present. |

##### Default is minimum required version check

| Nagios State | Description                                                             |
| ------------ | ----------------------------------------------------------------------- |
| `OK`         | Ideal state, hardware versions within tolerance.                        |
| `WARNING`    | Hardware versions older than the host or cluster default value present. |
| `CRITICAL`   | Not used by this monitoring mode.                                       |

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

#### `check_vmware_snapshots_count`

| Nagios State | Description                                                             |
| ------------ | ----------------------------------------------------------------------- |
| `OK`         | Ideal state, snapshots count per VM within bounds.                      |
| `WARNING`    | Snapshots count per VM crossed user-specified threshold for this state. |
| `CRITICAL`   | Snapshots count per VM crossed user-specified threshold for this state. |

#### `check_vmware_snapshots_size`

| Nagios State | Description                                                                         |
| ------------ | ----------------------------------------------------------------------------------- |
| `OK`         | Ideal state, snapshots size within bounds.                                          |
| `WARNING`    | Cumulative snapshots size for a VM crossed user-specified threshold for this state. |
| `CRITICAL`   | Cumulative snapshots size for a VM crossed user-specified threshold for this state. |

#### `check_vmware_rps_memory`

| Nagios State | Description                                                     |
| ------------ | --------------------------------------------------------------- |
| `OK`         | Ideal state, memory usage across Resources Pools within bounds. |
| `WARNING`    | Memory usage crossed user-specified threshold for this state.   |
| `CRITICAL`   | Memory usage crossed user-specified threshold for this state.   |

#### `check_vmware_host_memory`

| Nagios State | Description                                                                    |
| ------------ | ------------------------------------------------------------------------------ |
| `OK`         | Ideal state, memory usage for the specified ESXi host system is within bounds. |
| `WARNING`    | Memory usage crossed user-specified threshold for this state.                  |
| `CRITICAL`   | Memory usage crossed user-specified threshold for this state.                  |

#### `check_vmware_host_cpu`

| Nagios State | Description                                                                 |
| ------------ | --------------------------------------------------------------------------- |
| `OK`         | Ideal state, CPU usage for the specified ESXi host system is within bounds. |
| `WARNING`    | CPU usage crossed user-specified threshold for this state.                  |
| `CRITICAL`   | CPU usage crossed user-specified threshold for this state.                  |

#### `check_vmware_vm_power_uptime`

| Nagios State | Description                                                            |
| ------------ | ---------------------------------------------------------------------- |
| `OK`         | Ideal state, VM power cycle uptime is within bounds.                   |
| `WARNING`    | VM power cycle uptime crossed user-specified threshold for this state. |
| `CRITICAL`   | VM power cycle uptime crossed user-specified threshold for this state. |

#### `check_vmware_disk_consolidation`

| Nagios State | Description                                    |
| ------------ | ---------------------------------------------- |
| `OK`         | Ideal state, VM disk consolidation not needed. |
| `WARNING`    | Not used by this plugin.                       |
| `CRITICAL`   | Disk consolidation needed for one or more VMs. |

#### `check_vmware_question`

| Nagios State | Description                                              |
| ------------ | -------------------------------------------------------- |
| `OK`         | Ideal state, no VMs require an interactive response.     |
| `WARNING`    | Not used by this plugin.                                 |
| `CRITICAL`   | An interactive response is required for one or more VMs. |

#### `check_vmware_alarms`

| Nagios State | Description                                             |
| ------------ | ------------------------------------------------------- |
| `OK`         | Ideal state, no non-excluded Triggered Alarms detected. |
| `WARNING`    | One or more non-excluded alarms with a yellow status.   |
| `CRITICAL`   | One or more non-excluded alarms with a red status.      |

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
| `ll`, `log-level` | No       | `info`  | No     | `disabled`, `panic`, `fatal`, `error`, `warn`, `info`, `debug`, `trace` | Log message priority filter. Log messages with a lower level are ignored. Log messages are sent to `stderr` by default. See [Output](#output) for more information.                                                                                                                                                        |
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
| `ll`, `log-level`           | No       | `info`  | No     | `disabled`, `panic`, `fatal`, `error`, `warn`, `info`, `debug`, `trace` | Log message priority filter. Log messages with a lower level are ignored. Log messages are sent to `stderr` by default. See [Output](#output) for more information.                                                                                                                                                        |
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

This plugin supports multiple monitoring modes. Each mode has options which
are incompatible with the others. As of this writing these monitoring modes
are *not* implemented as subcommands, though this may change in the future
based on feedback.

| Flag                             | Required  | Default | Repeat | Possible                                                                | Description                                                                                                                                                                                                                                                                                                                                                                                   |
| -------------------------------- | --------- | ------- | ------ | ----------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `branding`                       | No        | `false` | No     | `branding`                                                              | Toggles emission of branding details with plugin status details. This output is disabled by default.                                                                                                                                                                                                                                                                                          |
| `h`, `help`                      | No        | `false` | No     | `h`, `help`                                                             | Show Help text along with the list of supported flags.                                                                                                                                                                                                                                                                                                                                        |
| `v`, `version`                   | No        | `false` | No     | `v`, `version`                                                          | Whether to display application version and then immediately exit application.                                                                                                                                                                                                                                                                                                                 |
| `ll`, `log-level`                | No        | `info`  | No     | `disabled`, `panic`, `fatal`, `error`, `warn`, `info`, `debug`, `trace` | Log message priority filter. Log messages with a lower level are ignored. Log messages are sent to `stderr` by default. See [Output](#output) for more information.                                                                                                                                                                                                                           |
| `p`, `port`                      | No        | `443`   | No     | *positive whole number between 1-65535, inclusive*                      | TCP port of the remote ESXi host or vCenter instance. This is usually 443 (HTTPS).                                                                                                                                                                                                                                                                                                            |
| `t`, `timeout`                   | No        | `10`    | No     | *positive whole number of seconds*                                      | Timeout value in seconds allowed before a plugin execution attempt is abandoned and an error returned.                                                                                                                                                                                                                                                                                        |
| `s`, `server`                    | **Yes**   |         | No     | *fully-qualified domain name or IP Address*                             | The fully-qualified domain name or IP Address of the remote ESXi host or vCenter instance.                                                                                                                                                                                                                                                                                                    |
| `u`, `username`                  | **Yes**   |         | No     | *valid username*                                                        | Username with permission to access specified ESXi host or vCenter instance.                                                                                                                                                                                                                                                                                                                   |
| `pw`, `password`                 | **Yes**   |         | No     | *valid password*                                                        | Password used to login to ESXi host or vCenter instance.                                                                                                                                                                                                                                                                                                                                      |
| `domain`                         | No        |         | No     | *valid user domain*                                                     | (Optional) domain for user account used to login to ESXi host or vCenter instance.                                                                                                                                                                                                                                                                                                            |
| `trust-cert`                     | No        | `false` | No     | `true`, `false`                                                         | Whether the certificate should be trusted as-is without validation. WARNING: TLS is susceptible to man-in-the-middle attacks if enabling this option.                                                                                                                                                                                                                                         |
| `dc-name`                        | No        |         | No     | *valid vSphere datacenter name*                                         | Specifies the name of a vSphere Datacenter. If not specified, applicable plugins will attempt to use the default datacenter found in the vSphere environment. Not applicable to standalone ESXi hosts.                                                                                                                                                                                        |
| `host-name`                      | No        |         | No     | *valid ESXi host name*                                                  | ESXi host/server name as it is found within the vSphere inventory.                                                                                                                                                                                                                                                                                                                            |
| `cluster-name`                   | No        |         | No     | *valid vSphere cluster name*                                            | Specifies the name of a vSphere Cluster. If not specified, applicable plugins will attempt to use the default cluster found in the vSphere environment. Not applicable to standalone ESXi hosts.                                                                                                                                                                                              |
| `include-rp`                     | No        |         | No     | *comma-separated list of resource pool names*                           | Specifies a comma-separated list of Resource Pools that should be exclusively used when evaluating VMs. Specifying this option will also exclude any VMs from evaluation that are *outside* of a Resource Pool. This option is incompatible with specifying a list of Resource Pools to ignore or exclude from evaluation.                                                                    |
| `exclude-rp`                     | No        |         | No     | *comma-separated list of resource pool names*                           | Specifies a comma-separated list of Resource Pools that should be ignored when evaluating VMs. This option is incompatible with specifying a list of Resource Pools to include for evaluation.                                                                                                                                                                                                |
| `ignore-vm`                      | No        |         | No     | *comma-separated list of (vSphere) virtual machine names*               | Specifies a comma-separated list of VM names that should be ignored or excluded from evaluation.                                                                                                                                                                                                                                                                                              |
| `powered-off`                    | No        | `false` | No     | `true`, `false`                                                         | Toggles evaluation of powered off VMs in addition to powered on VMs. Evaluation of powered off VMs is disabled by default.                                                                                                                                                                                                                                                                    |
| `obw`, `outdated-by-warning`     | **Maybe** |         | No     | *positive whole number 1 or greater*                                    | If provided, this value is the WARNING threshold for outdated virtual hardware versions. If the current virtual hardware version for a VM is found to be more than this many versions older than the latest version a WARNING state is triggered. Required if specifying the CRITICAL threshold for outdated virtual hardware versions, incompatible with the minimum required version flag.  |
| `obw`, `outdated-by-critical`    | **Maybe** |         | No     | *positive whole number 1 or greater*                                    | If provided, this value is the CRITICAL threshold for outdated virtual hardware versions. If the current virtual hardware version for a VM is found to be more than this many versions older than the latest version a CRITICAL state is triggered. Required if specifying the WARNING threshold for outdated virtual hardware versions, incompatible with the minimum required version flag. |
| `mv`, `minimum-version`          | **Maybe** |         | No     | *positive whole number greater than 3*                                  | If provided, this value is the minimum virtual hardware version accepted for each Virtual Machine. Any Virtual Machine not meeting this minimum value is considered to be in a CRITICAL state. Per [KB 1003746](https://kb.vmware.com/s/article/1003746), version 3 appears to be the oldest version supported. Incompatible with the CRITICAL and WARNING threshold flags.                   |
| `dimv`, `default-is-min-version` | **Maybe** |         | No     | *positive whole number greater than 3*                                  | If provided, this value is the minimum virtual hardware version accepted for each Virtual Machine. Any Virtual Machine not meeting this minimum value is considered to be in a CRITICAL state. Per [KB 1003746](https://kb.vmware.com/s/article/1003746), version 3 appears to be the oldest version supported. Incompatible with the CRITICAL and WARNING threshold flags.                   |

#### `check_vmware_hs2ds2vms`

| Flag                 | Required  | Default | Repeat | Possible                                                                | Description                                                                                                                                                                                                                                                                                                                |
| -------------------- | --------- | ------- | ------ | ----------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `branding`           | No        | `false` | No     | `branding`                                                              | Toggles emission of branding details with plugin status details. This output is disabled by default.                                                                                                                                                                                                                       |
| `h`, `help`          | No        | `false` | No     | `h`, `help`                                                             | Show Help text along with the list of supported flags.                                                                                                                                                                                                                                                                     |
| `v`, `version`       | No        | `false` | No     | `v`, `version`                                                          | Whether to display application version and then immediately exit application.                                                                                                                                                                                                                                              |
| `ll`, `log-level`    | No        | `info`  | No     | `disabled`, `panic`, `fatal`, `error`, `warn`, `info`, `debug`, `trace` | Log message priority filter. Log messages with a lower level are ignored. Log messages are sent to `stderr` by default. See [Output](#output) for more information.                                                                                                                                                        |
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
| `ll`, `log-level`           | No       | `info`  | No     | `disabled`, `panic`, `fatal`, `error`, `warn`, `info`, `debug`, `trace` | Log message priority filter. Log messages with a lower level are ignored. Log messages are sent to `stderr` by default. See [Output](#output) for more information.                                    |
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
| `ll`, `log-level`    | No       | `info`  | No     | `disabled`, `panic`, `fatal`, `error`, `warn`, `info`, `debug`, `trace` | Log message priority filter. Log messages with a lower level are ignored. Log messages are sent to `stderr` by default. See [Output](#output) for more information.                                                                                                                                                        |
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

#### `check_vmware_snapshots_count`

| Flag                   | Required | Default | Repeat | Possible                                                                | Description                                                                                                                                                                                                                                                                                                                |
| ---------------------- | -------- | ------- | ------ | ----------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `branding`             | No       | `false` | No     | `branding`                                                              | Toggles emission of branding details with plugin status details. This output is disabled by default.                                                                                                                                                                                                                       |
| `h`, `help`            | No       | `false` | No     | `h`, `help`                                                             | Show Help text along with the list of supported flags.                                                                                                                                                                                                                                                                     |
| `v`, `version`         | No       | `false` | No     | `v`, `version`                                                          | Whether to display application version and then immediately exit application.                                                                                                                                                                                                                                              |
| `ll`, `log-level`      | No       | `info`  | No     | `disabled`, `panic`, `fatal`, `error`, `warn`, `info`, `debug`, `trace` | Log message priority filter. Log messages with a lower level are ignored. Log messages are sent to `stderr` by default. See [Output](#output) for more information.                                                                                                                                                        |
| `p`, `port`            | No       | `443`   | No     | *positive whole number between 1-65535, inclusive*                      | TCP port of the remote ESXi host or vCenter instance. This is usually 443 (HTTPS).                                                                                                                                                                                                                                         |
| `t`, `timeout`         | No       | `10`    | No     | *positive whole number of seconds*                                      | Timeout value in seconds allowed before a plugin execution attempt is abandoned and an error returned.                                                                                                                                                                                                                     |
| `s`, `server`          | **Yes**  |         | No     | *fully-qualified domain name or IP Address*                             | The fully-qualified domain name or IP Address of the remote ESXi host or vCenter instance.                                                                                                                                                                                                                                 |
| `u`, `username`        | **Yes**  |         | No     | *valid username*                                                        | Username with permission to access specified ESXi host or vCenter instance.                                                                                                                                                                                                                                                |
| `pw`, `password`       | **Yes**  |         | No     | *valid password*                                                        | Password used to login to ESXi host or vCenter instance.                                                                                                                                                                                                                                                                   |
| `domain`               | No       |         | No     | *valid user domain*                                                     | (Optional) domain for user account used to login to ESXi host or vCenter instance.                                                                                                                                                                                                                                         |
| `trust-cert`           | No       | `false` | No     | `true`, `false`                                                         | Whether the certificate should be trusted as-is without validation. WARNING: TLS is susceptible to man-in-the-middle attacks if enabling this option.                                                                                                                                                                      |
| `include-rp`           | No       |         | No     | *comma-separated list of resource pool names*                           | Specifies a comma-separated list of Resource Pools that should be exclusively used when evaluating VMs. Specifying this option will also exclude any VMs from evaluation that are *outside* of a Resource Pool. This option is incompatible with specifying a list of Resource Pools to ignore or exclude from evaluation. |
| `exclude-rp`           | No       |         | No     | *comma-separated list of resource pool names*                           | Specifies a comma-separated list of Resource Pools that should be ignored when evaluating VMs. This option is incompatible with specifying a list of Resource Pools to include for evaluation.                                                                                                                             |
| `ignore-vm`            | No       |         | No     | *comma-separated list of (vSphere) virtual machine names*               | Specifies a comma-separated list of VM names that should be ignored or excluded from evaluation.                                                                                                                                                                                                                           |
| `cc`, `count-critical` | No       | `4`     | No     | *count as positive whole number*                                        | Specifies the number of snapshots per VM when a CRITICAL threshold is reached.                                                                                                                                                                                                                                             |
| `cw`, `count-warning`  | No       | `25`    | No     | *count as positive whole number*                                        | Specifies the number of snapshots per VM when a WARNING threshold is reached.                                                                                                                                                                                                                                              |

#### `check_vmware_snapshots_size`

| Flag                  | Required | Default | Repeat | Possible                                                                | Description                                                                                                                                                                                                                                                                                                                |
| --------------------- | -------- | ------- | ------ | ----------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `branding`            | No       | `false` | No     | `branding`                                                              | Toggles emission of branding details with plugin status details. This output is disabled by default.                                                                                                                                                                                                                       |
| `h`, `help`           | No       | `false` | No     | `h`, `help`                                                             | Show Help text along with the list of supported flags.                                                                                                                                                                                                                                                                     |
| `v`, `version`        | No       | `false` | No     | `v`, `version`                                                          | Whether to display application version and then immediately exit application.                                                                                                                                                                                                                                              |
| `ll`, `log-level`     | No       | `info`  | No     | `disabled`, `panic`, `fatal`, `error`, `warn`, `info`, `debug`, `trace` | Log message priority filter. Log messages with a lower level are ignored. Log messages are sent to `stderr` by default. See [Output](#output) for more information.                                                                                                                                                        |
| `p`, `port`           | No       | `443`   | No     | *positive whole number between 1-65535, inclusive*                      | TCP port of the remote ESXi host or vCenter instance. This is usually 443 (HTTPS).                                                                                                                                                                                                                                         |
| `t`, `timeout`        | No       | `10`    | No     | *positive whole number of seconds*                                      | Timeout value in seconds allowed before a plugin execution attempt is abandoned and an error returned.                                                                                                                                                                                                                     |
| `s`, `server`         | **Yes**  |         | No     | *fully-qualified domain name or IP Address*                             | The fully-qualified domain name or IP Address of the remote ESXi host or vCenter instance.                                                                                                                                                                                                                                 |
| `u`, `username`       | **Yes**  |         | No     | *valid username*                                                        | Username with permission to access specified ESXi host or vCenter instance.                                                                                                                                                                                                                                                |
| `pw`, `password`      | **Yes**  |         | No     | *valid password*                                                        | Password used to login to ESXi host or vCenter instance.                                                                                                                                                                                                                                                                   |
| `domain`              | No       |         | No     | *valid user domain*                                                     | (Optional) domain for user account used to login to ESXi host or vCenter instance.                                                                                                                                                                                                                                         |
| `trust-cert`          | No       | `false` | No     | `true`, `false`                                                         | Whether the certificate should be trusted as-is without validation. WARNING: TLS is susceptible to man-in-the-middle attacks if enabling this option.                                                                                                                                                                      |
| `include-rp`          | No       |         | No     | *comma-separated list of resource pool names*                           | Specifies a comma-separated list of Resource Pools that should be exclusively used when evaluating VMs. Specifying this option will also exclude any VMs from evaluation that are *outside* of a Resource Pool. This option is incompatible with specifying a list of Resource Pools to ignore or exclude from evaluation. |
| `exclude-rp`          | No       |         | No     | *comma-separated list of resource pool names*                           | Specifies a comma-separated list of Resource Pools that should be ignored when evaluating VMs. This option is incompatible with specifying a list of Resource Pools to include for evaluation.                                                                                                                             |
| `ignore-vm`           | No       |         | No     | *comma-separated list of (vSphere) virtual machine names*               | Specifies a comma-separated list of VM names that should be ignored or excluded from evaluation.                                                                                                                                                                                                                           |
| `sc`, `size-critical` | No       | `40`    | No     | *size in GB as positive whole number*                                   | Specifies the cumulative size in GB of all snapshots for a Virtual Machine when a CRITICAL threshold is reached.                                                                                                                                                                                                           |
| `sw`, `size-warning`  | No       | `20`    | No     | *size in GB as positive whole number*                                   | Specifies the cumulative size in GB of all snapshots for a Virtual Machine when a WARNING threshold is reached.                                                                                                                                                                                                            |

#### `check_vmware_rps_memory`

| Flag                        | Required | Default | Repeat | Possible                                                                | Description                                                                                                                                                                                                                                                                                                                |
| --------------------------- | -------- | ------- | ------ | ----------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `branding`                  | No       | `false` | No     | `branding`                                                              | Toggles emission of branding details with plugin status details. This output is disabled by default.                                                                                                                                                                                                                       |
| `h`, `help`                 | No       | `false` | No     | `h`, `help`                                                             | Show Help text along with the list of supported flags.                                                                                                                                                                                                                                                                     |
| `v`, `version`              | No       | `false` | No     | `v`, `version`                                                          | Whether to display application version and then immediately exit application.                                                                                                                                                                                                                                              |
| `ll`, `log-level`           | No       | `info`  | No     | `disabled`, `panic`, `fatal`, `error`, `warn`, `info`, `debug`, `trace` | Log message priority filter. Log messages with a lower level are ignored. Log messages are sent to `stderr` by default. See [Output](#output) for more information.                                                                                                                                                        |
| `p`, `port`                 | No       | `443`   | No     | *positive whole number between 1-65535, inclusive*                      | TCP port of the remote ESXi host or vCenter instance. This is usually 443 (HTTPS).                                                                                                                                                                                                                                         |
| `t`, `timeout`              | No       | `10`    | No     | *positive whole number of seconds*                                      | Timeout value in seconds allowed before a plugin execution attempt is abandoned and an error returned.                                                                                                                                                                                                                     |
| `s`, `server`               | **Yes**  |         | No     | *fully-qualified domain name or IP Address*                             | The fully-qualified domain name or IP Address of the remote ESXi host or vCenter instance.                                                                                                                                                                                                                                 |
| `u`, `username`             | **Yes**  |         | No     | *valid username*                                                        | Username with permission to access specified ESXi host or vCenter instance.                                                                                                                                                                                                                                                |
| `pw`, `password`            | **Yes**  |         | No     | *valid password*                                                        | Password used to login to ESXi host or vCenter instance.                                                                                                                                                                                                                                                                   |
| `domain`                    | No       |         | No     | *valid user domain*                                                     | (Optional) domain for user account used to login to ESXi host or vCenter instance.                                                                                                                                                                                                                                         |
| `trust-cert`                | No       | `false` | No     | `true`, `false`                                                         | Whether the certificate should be trusted as-is without validation. WARNING: TLS is susceptible to man-in-the-middle attacks if enabling this option.                                                                                                                                                                      |
| `include-rp`                | No       |         | No     | *comma-separated list of resource pool names*                           | Specifies a comma-separated list of Resource Pools that should be exclusively used when evaluating VMs. Specifying this option will also exclude any VMs from evaluation that are *outside* of a Resource Pool. This option is incompatible with specifying a list of Resource Pools to ignore or exclude from evaluation. |
| `exclude-rp`                | No       |         | No     | *comma-separated list of resource pool names*                           | Specifies a comma-separated list of Resource Pools that should be ignored when evaluating VMs. This option is incompatible with specifying a list of Resource Pools to include for evaluation.                                                                                                                             |
| `mma`, `memory-max-allowed` | **Yes**  | `0`     | No     | *positive whole number of vCPUs*                                        | Specifies the maximum amount of memory that we are allowed to consume in GB (as a whole number) in the target VMware environment across all specified Resource Pools. VMs that are running outside of resource pools are not considered in these calculations.                                                             |
| `mc`, `memory-use-critical` | No       | `95`    | No     | *percentage as positive whole number*                                   | Specifies the percentage of memory use (as a whole number) across all specified Resource Pools when a CRITICAL threshold is reached.                                                                                                                                                                                       |
| `mw`, `memory-use-warning`  | No       | `100`   | No     | *percentage as positive whole number*                                   | Specifies the percentage of memory use (as a whole number) across all specified Resource Pools when a WARNING threshold is reached.                                                                                                                                                                                        |

#### `check_vmware_host_memory`

| Flag                          | Required | Default | Repeat | Possible                                                                | Description                                                                                                                                                                                            |
| ----------------------------- | -------- | ------- | ------ | ----------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `branding`                    | No       | `false` | No     | `branding`                                                              | Toggles emission of branding details with plugin status details. This output is disabled by default.                                                                                                   |
| `h`, `help`                   | No       | `false` | No     | `h`, `help`                                                             | Show Help text along with the list of supported flags.                                                                                                                                                 |
| `v`, `version`                | No       | `false` | No     | `v`, `version`                                                          | Whether to display application version and then immediately exit application.                                                                                                                          |
| `ll`, `log-level`             | No       | `info`  | No     | `disabled`, `panic`, `fatal`, `error`, `warn`, `info`, `debug`, `trace` | Log message priority filter. Log messages with a lower level are ignored. Log messages are sent to `stderr` by default. See [Output](#output) for more information.                                    |
| `p`, `port`                   | No       | `443`   | No     | *positive whole number between 1-65535, inclusive*                      | TCP port of the remote ESXi host or vCenter instance. This is usually 443 (HTTPS).                                                                                                                     |
| `t`, `timeout`                | No       | `10`    | No     | *positive whole number of seconds*                                      | Timeout value in seconds allowed before a plugin execution attempt is abandoned and an error returned.                                                                                                 |
| `s`, `server`                 | **Yes**  |         | No     | *fully-qualified domain name or IP Address*                             | The fully-qualified domain name or IP Address of the remote ESXi host or vCenter instance.                                                                                                             |
| `u`, `username`               | **Yes**  |         | No     | *valid username*                                                        | Username with permission to access specified ESXi host or vCenter instance.                                                                                                                            |
| `pw`, `password`              | **Yes**  |         | No     | *valid password*                                                        | Password used to login to ESXi host or vCenter instance.                                                                                                                                               |
| `domain`                      | No       |         | No     | *valid user domain*                                                     | (Optional) domain for user account used to login to ESXi host or vCenter instance.                                                                                                                     |
| `trust-cert`                  | No       | `false` | No     | `true`, `false`                                                         | Whether the certificate should be trusted as-is without validation. WARNING: TLS is susceptible to man-in-the-middle attacks if enabling this option.                                                  |
| `dc-name`                     | No       |         | No     | *valid vSphere datacenter name*                                         | Specifies the name of a vSphere Datacenter. If not specified, applicable plugins will attempt to use the default datacenter found in the vSphere environment. Not applicable to standalone ESXi hosts. |
| `host-name`                   | **Yes**  |         | No     | *valid ESXi host name*                                                  | ESXi host/server name as it is found within the vSphere inventory.                                                                                                                                     |
| `mc`, `memory-usage-critical` | No       | `95`    | No     | *percentage as positive whole number*                                   | Specifies the percentage of memory use (as a whole number) when a CRITICAL threshold is reached.                                                                                                       |
| `mw`, `memory-usage-warning`  | No       | `80`    | No     | *percentage as positive whole number*                                   | Specifies the percentage of memory use (as a whole number) when a WARNING threshold is reached.                                                                                                        |

#### `check_vmware_host_cpu`

| Flag                       | Required | Default | Repeat | Possible                                                                | Description                                                                                                                                                                                            |
| -------------------------- | -------- | ------- | ------ | ----------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `branding`                 | No       | `false` | No     | `branding`                                                              | Toggles emission of branding details with plugin status details. This output is disabled by default.                                                                                                   |
| `h`, `help`                | No       | `false` | No     | `h`, `help`                                                             | Show Help text along with the list of supported flags.                                                                                                                                                 |
| `v`, `version`             | No       | `false` | No     | `v`, `version`                                                          | Whether to display application version and then immediately exit application.                                                                                                                          |
| `ll`, `log-level`          | No       | `info`  | No     | `disabled`, `panic`, `fatal`, `error`, `warn`, `info`, `debug`, `trace` | Log message priority filter. Log messages with a lower level are ignored. Log messages are sent to `stderr` by default. See [Output](#output) for more information.                                    |
| `p`, `port`                | No       | `443`   | No     | *positive whole number between 1-65535, inclusive*                      | TCP port of the remote ESXi host or vCenter instance. This is usually 443 (HTTPS).                                                                                                                     |
| `t`, `timeout`             | No       | `10`    | No     | *positive whole number of seconds*                                      | Timeout value in seconds allowed before a plugin execution attempt is abandoned and an error returned.                                                                                                 |
| `s`, `server`              | **Yes**  |         | No     | *fully-qualified domain name or IP Address*                             | The fully-qualified domain name or IP Address of the remote ESXi host or vCenter instance.                                                                                                             |
| `u`, `username`            | **Yes**  |         | No     | *valid username*                                                        | Username with permission to access specified ESXi host or vCenter instance.                                                                                                                            |
| `pw`, `password`           | **Yes**  |         | No     | *valid password*                                                        | Password used to login to ESXi host or vCenter instance.                                                                                                                                               |
| `domain`                   | No       |         | No     | *valid user domain*                                                     | (Optional) domain for user account used to login to ESXi host or vCenter instance.                                                                                                                     |
| `trust-cert`               | No       | `false` | No     | `true`, `false`                                                         | Whether the certificate should be trusted as-is without validation. WARNING: TLS is susceptible to man-in-the-middle attacks if enabling this option.                                                  |
| `dc-name`                  | No       |         | No     | *valid vSphere datacenter name*                                         | Specifies the name of a vSphere Datacenter. If not specified, applicable plugins will attempt to use the default datacenter found in the vSphere environment. Not applicable to standalone ESXi hosts. |
| `host-name`                | **Yes**  |         | No     | *valid ESXi host name*                                                  | ESXi host/server name as it is found within the vSphere inventory.                                                                                                                                     |
| `cc`, `cpu-usage-critical` | No       | `95`    | No     | *percentage as positive whole number*                                   | Specifies the percentage of CPU use (as a whole number) when a CRITICAL threshold is reached.                                                                                                          |
| `cw`, `cpu-usage-warning`  | No       | `80`    | No     | *percentage as positive whole number*                                   | Specifies the percentage of CPU use (as a whole number) when a WARNING threshold is reached.                                                                                                           |

#### `check_vmware_vm_power_uptime`

| Flag                    | Required | Default | Repeat | Possible                                                                | Description                                                                                                                                                                                                                                                                                                                |
| ----------------------- | -------- | ------- | ------ | ----------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `branding`              | No       | `false` | No     | `branding`                                                              | Toggles emission of branding details with plugin status details. This output is disabled by default.                                                                                                                                                                                                                       |
| `h`, `help`             | No       | `false` | No     | `h`, `help`                                                             | Show Help text along with the list of supported flags.                                                                                                                                                                                                                                                                     |
| `v`, `version`          | No       | `false` | No     | `v`, `version`                                                          | Whether to display application version and then immediately exit application.                                                                                                                                                                                                                                              |
| `ll`, `log-level`       | No       | `info`  | No     | `disabled`, `panic`, `fatal`, `error`, `warn`, `info`, `debug`, `trace` | Log message priority filter. Log messages with a lower level are ignored. Log messages are sent to `stderr` by default. See [Output](#output) for more information.                                                                                                                                                        |
| `p`, `port`             | No       | `443`   | No     | *positive whole number between 1-65535, inclusive*                      | TCP port of the remote ESXi host or vCenter instance. This is usually 443 (HTTPS).                                                                                                                                                                                                                                         |
| `t`, `timeout`          | No       | `10`    | No     | *positive whole number of seconds*                                      | Timeout value in seconds allowed before a plugin execution attempt is abandoned and an error returned.                                                                                                                                                                                                                     |
| `s`, `server`           | **Yes**  |         | No     | *fully-qualified domain name or IP Address*                             | The fully-qualified domain name or IP Address of the remote ESXi host or vCenter instance.                                                                                                                                                                                                                                 |
| `u`, `username`         | **Yes**  |         | No     | *valid username*                                                        | Username with permission to access specified ESXi host or vCenter instance.                                                                                                                                                                                                                                                |
| `pw`, `password`        | **Yes**  |         | No     | *valid password*                                                        | Password used to login to ESXi host or vCenter instance.                                                                                                                                                                                                                                                                   |
| `domain`                | No       |         | No     | *valid user domain*                                                     | (Optional) domain for user account used to login to ESXi host or vCenter instance.                                                                                                                                                                                                                                         |
| `trust-cert`            | No       | `false` | No     | `true`, `false`                                                         | Whether the certificate should be trusted as-is without validation. WARNING: TLS is susceptible to man-in-the-middle attacks if enabling this option.                                                                                                                                                                      |
| `include-rp`            | No       |         | No     | *comma-separated list of resource pool names*                           | Specifies a comma-separated list of Resource Pools that should be exclusively used when evaluating VMs. Specifying this option will also exclude any VMs from evaluation that are *outside* of a Resource Pool. This option is incompatible with specifying a list of Resource Pools to ignore or exclude from evaluation. |
| `exclude-rp`            | No       |         | No     | *comma-separated list of resource pool names*                           | Specifies a comma-separated list of Resource Pools that should be ignored when evaluating VMs. This option is incompatible with specifying a list of Resource Pools to include for evaluation.                                                                                                                             |
| `ignore-vm`             | No       |         | No     | *comma-separated list of (vSphere) virtual machine names*               | Specifies a comma-separated list of VM names that should be ignored or excluded from evaluation.                                                                                                                                                                                                                           |
| `uc`, `uptime-critical` | No       | `90`    | No     | *days as positive whole number*                                         | Specifies the power cycle (off/on) uptime in days per VM when a CRITICAL threshold is reached.                                                                                                                                                                                                                             |
| `uw`, `uptime-warning`  | No       | `60`    | No     | *days as positive whole number*                                         | Specifies the power cycle (off/on) uptime in days per VM when a WARNING threshold is reached.                                                                                                                                                                                                                              |

#### `check_vmware_disk_consolidation`

| Flag              | Required | Default | Repeat | Possible                                                                | Description                                                                                                                                                                                                                                                                                                                |
| ----------------- | -------- | ------- | ------ | ----------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `branding`        | No       | `false` | No     | `branding`                                                              | Toggles emission of branding details with plugin status details. This output is disabled by default.                                                                                                                                                                                                                       |
| `h`, `help`       | No       | `false` | No     | `h`, `help`                                                             | Show Help text along with the list of supported flags.                                                                                                                                                                                                                                                                     |
| `v`, `version`    | No       | `false` | No     | `v`, `version`                                                          | Whether to display application version and then immediately exit application.                                                                                                                                                                                                                                              |
| `ll`, `log-level` | No       | `info`  | No     | `disabled`, `panic`, `fatal`, `error`, `warn`, `info`, `debug`, `trace` | Log message priority filter. Log messages with a lower level are ignored. Log messages are sent to `stderr` by default. See [Output](#output) for more information.                                                                                                                                                        |
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
| `trigger-reload`  | No       | `false` | No     | `true`, `false`                                                         | Trigger a reload operation for each VM evaluated. This option ensures that the most current state data is evaluated, but increases plugin runtime. If using this, you should also adjust the `--timeout` value and potentially your monitor system's service check timeout setting.                                        |

#### `check_vmware_question`

| Flag              | Required | Default | Repeat | Possible                                                                | Description                                                                                                                                                                                                                                                                                                                |
| ----------------- | -------- | ------- | ------ | ----------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `branding`        | No       | `false` | No     | `branding`                                                              | Toggles emission of branding details with plugin status details. This output is disabled by default.                                                                                                                                                                                                                       |
| `h`, `help`       | No       | `false` | No     | `h`, `help`                                                             | Show Help text along with the list of supported flags.                                                                                                                                                                                                                                                                     |
| `v`, `version`    | No       | `false` | No     | `v`, `version`                                                          | Whether to display application version and then immediately exit application.                                                                                                                                                                                                                                              |
| `ll`, `log-level` | No       | `info`  | No     | `disabled`, `panic`, `fatal`, `error`, `warn`, `info`, `debug`, `trace` | Log message priority filter. Log messages with a lower level are ignored. Log messages are sent to `stderr` by default. See [Output](#output) for more information.                                                                                                                                                        |
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

#### `check_vmware_alarms`

| Flag                  | Required | Default | Repeat | Possible                                                                                                                                                                       | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 |
| --------------------- | -------- | ------- | ------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `branding`            | No       | `false` | No     | `branding`                                                                                                                                                                     | Toggles emission of branding details with plugin status details. This output is disabled by default.                                                                                                                                                                                                                                                                                                                                                                                                        |
| `h`, `help`           | No       | `false` | No     | `h`, `help`                                                                                                                                                                    | Show Help text along with the list of supported flags.                                                                                                                                                                                                                                                                                                                                                                                                                                                      |
| `v`, `version`        | No       | `false` | No     | `v`, `version`                                                                                                                                                                 | Whether to display application version and then immediately exit application.                                                                                                                                                                                                                                                                                                                                                                                                                               |
| `ll`, `log-level`     | No       | `info`  | No     | `disabled`, `panic`, `fatal`, `error`, `warn`, `info`, `debug`, `trace`                                                                                                        | Log message priority filter. Log messages with a lower level are ignored. Log messages are sent to `stderr` by default. See [Output](#output) for more information.                                                                                                                                                                                                                                                                                                                                         |
| `p`, `port`           | No       | `443`   | No     | *positive whole number between 1-65535, inclusive*                                                                                                                             | TCP port of the remote ESXi host or vCenter instance. This is usually 443 (HTTPS).                                                                                                                                                                                                                                                                                                                                                                                                                          |
| `t`, `timeout`        | No       | `10`    | No     | *positive whole number of seconds*                                                                                                                                             | Timeout value in seconds allowed before a plugin execution attempt is abandoned and an error returned.                                                                                                                                                                                                                                                                                                                                                                                                      |
| `s`, `server`         | **Yes**  |         | No     | *fully-qualified domain name or IP Address*                                                                                                                                    | The fully-qualified domain name or IP Address of the remote ESXi host or vCenter instance.                                                                                                                                                                                                                                                                                                                                                                                                                  |
| `u`, `username`       | **Yes**  |         | No     | *valid username*                                                                                                                                                               | Username with permission to access specified ESXi host or vCenter instance.                                                                                                                                                                                                                                                                                                                                                                                                                                 |
| `pw`, `password`      | **Yes**  |         | No     | *valid password*                                                                                                                                                               | Password used to login to ESXi host or vCenter instance.                                                                                                                                                                                                                                                                                                                                                                                                                                                    |
| `domain`              | No       |         | No     | *valid user domain*                                                                                                                                                            | (Optional) domain for user account used to login to ESXi host or vCenter instance.                                                                                                                                                                                                                                                                                                                                                                                                                          |
| `trust-cert`          | No       | `false` | No     | `true`, `false`                                                                                                                                                                | Whether the certificate should be trusted as-is without validation. WARNING: TLS is susceptible to man-in-the-middle attacks if enabling this option.                                                                                                                                                                                                                                                                                                                                                       |
| `dc-name`             | No       |         | No     | *comma-separated list of valid vSphere datacenter names*                                                                                                                       | Specifies the name of one or more vSphere Datacenters. If not specified, applicable plugins will attempt to evaluate all visible datacenters found in the vSphere environment. Not applicable to standalone ESXi hosts.                                                                                                                                                                                                                                                                                     |
| `include-entity-type` | No       |         | No     | [*comma-separated list of valid managed object type keywords*][vsphere-managed-object-reference]                                                                               | If specified, triggered alarms will only be evaluated if the associated entity type (e.g., `Datastore`) matches one of the specified values; while multiple explicit inclusions are allowed, explicit exclusions have precedence over explicit inclusions and will exclude the triggered alarm from further evaluation.                                                                                                                                                                                     |
| `exclude-entity-type` | No       |         | No     | [*comma-separated list of valid managed object type keywords*][vsphere-managed-object-reference]                                                                               | If specified, triggered alarms will only be evaluated if the associated entity type (e.g., `Datastore`) does NOT match one of the specified values; while multiple explicit inclusions are allowed, explicit exclusions have precedence over explicit inclusions and will exclude the triggered alarm from further evaluation.                                                                                                                                                                              |
| `include-entity-name` | No       |         | No     | *comma-separated list of vSphere inventory object names*                                                                                                                       | If specified, triggered alarms will only be evaluated if the associated entity name (e.g., `node1.example.com`) matches one of the specified values; while multiple explicit inclusions are allowed, explicit exclusions have precedence over explicit inclusions and will exclude the triggered alarm from further evaluation.                                                                                                                                                                             |
| `exclude-entity-name` | No       |         | No     | *comma-separated list of vSphere inventory object names*                                                                                                                       | If specified, triggered alarms will only be evaluated if the associated entity name (e.g., `node1.example.com`) does NOT match one of the specified values; while multiple explicit inclusions are allowed, explicit exclusions have precedence over explicit inclusions and will exclude the triggered alarm from further evaluation.                                                                                                                                                                      |
| `include-entity-rp`   | No       |         | No     | *comma-separated list of resource pool names*                                                                                                                                  | If specified, triggered alarms will only be evaluated if the associated entity is part of one of the specified Resource Pools (case-insensitive match on the name) and is not explicitly excluded by another filter in the pipeline; while multiple explicit inclusions are allowed, explicit exclusions have precedence over explicit inclusions and will exclude the triggered alarm from further evaluation.                                                                                             |
| `exclude-entity-rp`   | No       |         | No     | *comma-separated list of resource pool names*                                                                                                                                  | If specified, triggered alarms will only be evaluated if the associated entity is NOT part of one of the specified Resource Pools (case-insensitive match on the name) and is not explicitly excluded by another filter in the pipeline; while multiple explicit inclusions are allowed, explicit exclusions have precedence over explicit inclusions and will exclude the triggered alarm from further evaluation.                                                                                         |
| `eval-acknowledged`   | No       | `false` | No     | `true`, `false`                                                                                                                                                                | Toggles evaluation of acknowledged triggered alarms in addition to unacknowledged triggered alarms. Evaluation of acknowledged alarms is disabled by default.                                                                                                                                                                                                                                                                                                                                               |
| `include-name`        | No       |         | No     | *valid custom or* [*default alarm names*][vsphere-default-alarms]                                                                                                              | If specified, triggered alarms will only be evaluated if the alarm name (e.g., `Datastore usage on disk`) case-insensitively matches one of the specified substring values (e.g., `datastore` or `datastore usage`) and is not explicitly excluded by another filter in the pipeline; while multiple explicit inclusions are allowed, explicit exclusions have precedence over explicit inclusions and will exclude the triggered alarm from further evaluation.                                            |
| `exclude-name`        | No       |         | No     | *valid custom or* [*default alarm names*][vsphere-default-alarms]                                                                                                              | If specified, triggered alarms will only be evaluated if the alarm name (e.g., `Datastore usage on disk`) DOES NOT case-insensitively match one of the specified substring values (e.g., `datastore` or `datastore usage`) and is not explicitly excluded by another filter in the pipeline; while multiple explicit inclusions are allowed, explicit exclusions have precedence over explicit inclusions and will exclude the triggered alarm from further evaluation.                                     |
| `include-desc`        | No       |         | No     | *valid custom or* [*default alarm descriptions*][vsphere-default-alarms]                                                                                                       | If specified, triggered alarms will only be evaluated if the alarm description (e.g., `Default alarm to monitor datastore disk usage`) case-insensitively matches one of the specified substring values (e.g., `datastore disk` or `monitor datastore`) and is not explicitly excluded by another filter in the pipeline; while multiple explicit inclusions are allowed, explicit exclusions have precedence over explicit inclusions and will exclude the triggered alarm from further evaluation.        |
| `exclude-desc`        | No       |         | No     | *valid custom or* [*default alarm descriptions*][vsphere-default-alarms]                                                                                                       | If specified, triggered alarms will only be evaluated if the alarm description (e.g., `Default alarm to monitor datastore disk usage`) DOES NOT case-insensitively match one of the specified substring values (e.g., `datastore disk` or `monitor datastore`) and is not explicitly excluded by another filter in the pipeline; while multiple explicit inclusions are allowed, explicit exclusions have precedence over explicit inclusions and will exclude the triggered alarm from further evaluation. |
| `include-status`      | No       |         | No     | *valid* [*managed entity status*][vsphere-manged-entity-status] (excluding `green`) or [Nagios state][nagios-state-types] (excluding `OK`) (`WARNING`, `CRITICAL` , `UNKNOwN`) | If specified, triggered alarms will only be evaluated if the alarm status (e.g., `yellow`) case-insensitively matches one of the specified keywords (e.g., `yellow` or `warning`) and is not explicitly excluded by another filter in the pipeline; while multiple explicit inclusions are allowed, explicit exclusions have precedence over explicit inclusions and will exclude the triggered alarm from further evaluation.                                                                              |
| `exclude-status`      | No       |         | No     | *valid* [*managed entity status*][vsphere-manged-entity-status]                                                                                                                | If specified, triggered alarms will only be evaluated if the alarm status (e.g., `yellow`) DOES NOT case-insensitively match one of the specified keywords (e.g., `yellow` or `warning`) and is not explicitly excluded by another filter in the pipeline; while multiple explicit inclusions are allowed, explicit exclusions have precedence over explicit inclusions and will exclude the triggered alarm from further evaluation.                                                                       |

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
- Service Check results output is sent to `stdout`
- Logging output is enabled at the `info` level.
  - logging output is sent to `stderr` by default
  - logging output is intended to be seen when invoking the
    plugin directly via CLI (often for troubleshooting)
    - see [Output](#output) for potential conflicts with some monitoring
      systems

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
- Service Check results output is sent to `stdout`
- Logging output is enabled at the `info` level.
  - logging output is sent to `stderr` by default
  - logging output is intended to be seen when invoking the
    plugin directly via CLI (often for troubleshooting)
    - see [Output](#output) for potential conflicts with some monitoring
      systems

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

This plugin supports four monitoring modes. Each is incompatible with the
other, so an example is provided for each mode. See the [overview](#overview)
section for further information.

#### Homogeneous version check

##### CLI invocation

```ShellSession
/usr/lib/nagios/plugins/check_vmware_vhw --username SERVICE_ACCOUNT_NAME --password "SERVICE_ACCOUNT_PASSWORD" --server vc1.example.com --exclude-rp "Desktops" --ignore-vm "test1.example.com,redmine.example.com,TESTING-AC,RHEL7-TEST" --trust-cert --log-level info
```

See the [configuration options](#configuration-options) section for all
command-line settings supported by this plugin along with descriptions of
each. See the [contrib](#contrib) section for information regarding example
command definitions and Nagios configuration files.

Of note:

- This monitoring mode asserts that all hardware versions match.
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
- Service Check results output is sent to `stdout`
- Logging output is enabled at the `info` level.
  - logging output is sent to `stderr` by default
  - logging output is intended to be seen when invoking the
    plugin directly via CLI (often for troubleshooting)
    - see [Output](#output) for potential conflicts with some monitoring
      systems

##### Command definition

```shell
# /etc/nagios-plugins/config/vmware-virtual-hardware.cfg

# Look at all pools, all VMs, do not evaluate any VMs that are powered off.
# This variation of the command is most useful for environments where all VMs
# are monitored equally.
define command{
    command_name    check_vmware_vhw_homogeneous
    command_line    /usr/lib/nagios/plugins/check_vmware_vhw --server '$HOSTNAME$' --domain '$ARG1$' --username '$ARG2$' --password '$ARG3$' --trust-cert --log-level info
    }
```

See the [configuration options](#configuration-options) section for all
command-line settings supported by this plugin along with descriptions of
each. See the [contrib](#contrib) section for information regarding example
command definitions and Nagios configuration files.

#### Outdated-by or threshold range check

##### CLI invocation

```ShellSession
/usr/lib/nagios/plugins/check_vmware_vhw --username SERVICE_ACCOUNT_NAME --password "SERVICE_ACCOUNT_PASSWORD" --server vc1.example.com --exclude-rp "Desktops" --ignore-vm "test1.example.com,redmine.example.com,TESTING-AC,RHEL7-TEST" --outdated-by-warning 1 --outdated-by-critical 5 --trust-cert --log-level info
```

See the [configuration options](#configuration-options) section for all
command-line settings supported by this plugin along with descriptions of
each. See the [contrib](#contrib) section for information regarding example
command definitions and Nagios configuration files.

Of note:

- Assuming that the latest hardware version is `15`, this monitoring mode
  permits hardware versions as old as `14` without `WARNING` state and as old
  as `10` without `CRITICAL` state change.
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
- Service Check results output is sent to `stdout`
- Logging output is enabled at the `info` level.
  - logging output is sent to `stderr` by default
  - logging output is intended to be seen when invoking the
    plugin directly via CLI (often for troubleshooting)
    - see [Output](#output) for potential conflicts with some monitoring
      systems

##### Command definition

```shell
# /etc/nagios-plugins/config/vmware-virtual-hardware.cfg

# Look at all pools, all VMs, do not evaluate any VMs that are powered off.
# This variation of the command is most useful for environments where all VMs
# are monitored equally.
define command{
    command_name    check_vmware_vhw_thresholds
    command_line    /usr/lib/nagios/plugins/check_vmware_vhw --server '$HOSTNAME$' --domain '$ARG1$' --username '$ARG2$' --password '$ARG3$' --outdated-by-warning '$ARG4$' --outdated-by-critical '$ARG5$' --trust-cert --log-level info
    }

```

See the [configuration options](#configuration-options) section for all
command-line settings supported by this plugin along with descriptions of
each. See the [contrib](#contrib) section for information regarding example
command definitions and Nagios configuration files.

#### Minimum required version check

##### CLI invocation

```ShellSession
/usr/lib/nagios/plugins/check_vmware_vhw --username SERVICE_ACCOUNT_NAME --password "SERVICE_ACCOUNT_PASSWORD" --server vc1.example.com --exclude-rp "Desktops" --ignore-vm "test1.example.com,redmine.example.com,TESTING-AC,RHEL7-TEST" --minimum-version 15 --trust-cert --log-level info
```

See the [configuration options](#configuration-options) section for all
command-line settings supported by this plugin along with descriptions of
each. See the [contrib](#contrib) section for information regarding example
command definitions and Nagios configuration files.

Of note:

- The minimum hardware version `15` is required, while newer versions are
  permitted, older versions will trigger a plugin state change.
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
- Service Check results output is sent to `stdout`
- Logging output is enabled at the `info` level.
  - logging output is sent to `stderr` by default
  - logging output is intended to be seen when invoking the
    plugin directly via CLI (often for troubleshooting)
    - see [Output](#output) for potential conflicts with some monitoring
      systems

##### Command definition

```shell
# /etc/nagios-plugins/config/vmware-virtual-hardware.cfg

# Look at all pools, all VMs, do not evaluate any VMs that are powered off.
# This variation of the command is most useful for environments where all VMs
# are monitored equally.
define command{
    command_name    check_vmware_vhw_minreq
    command_line    /usr/lib/nagios/plugins/check_vmware_vhw --server '$HOSTNAME$' --domain '$ARG1$' --username '$ARG2$' --password '$ARG3$' --minimum-version '$ARG4$' --trust-cert --log-level info
    }

```

See the [configuration options](#configuration-options) section for all
command-line settings supported by this plugin along with descriptions of
each. See the [contrib](#contrib) section for information regarding example
command definitions and Nagios configuration files.

#### Default is minimum required version check

##### CLI invocation

```ShellSession
/usr/lib/nagios/plugins/check_vmware_vhw --username SERVICE_ACCOUNT_NAME --password "SERVICE_ACCOUNT_PASSWORD" --server vc1.example.com --exclude-rp "Desktops" --ignore-vm "test1.example.com,redmine.example.com,TESTING-AC,RHEL7-TEST" --default-is-min-version --trust-cert --log-level info
```

See the [configuration options](#configuration-options) section for all
command-line settings supported by this plugin along with descriptions of
each. See the [contrib](#contrib) section for information regarding example
command definitions and Nagios configuration files.

Of note:

- The default host or cluster hardware version is required
  - while newer versions are permitted, older versions will trigger a plugin
    state change.
- Neither a host name nor a cluster name is provided
  - the plugin will attempt to use the default `ComputeResource` in order to
    determine the default hardware version
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
- Service Check results output is sent to `stdout`
- Logging output is enabled at the `info` level.
  - logging output is sent to `stderr` by default
  - logging output is intended to be seen when invoking the
    plugin directly via CLI (often for troubleshooting)
    - see [Output](#output) for potential conflicts with some monitoring
      systems

##### Command definition

```shell
# /etc/nagios-plugins/config/vmware-virtual-hardware.cfg

# Look at all pools, all VMs, do not evaluate any VMs that are powered off.
# This variation of the command is most useful for environments where all VMs
# are monitored equally.
define command{
    command_name    check_vmware_vhw_defreq
    command_line    /usr/lib/nagios/plugins/check_vmware_vhw --server '$HOSTNAME$' --domain '$ARG1$' --username '$ARG2$' --password '$ARG3$' --cluster-name '$ARG4$' --default-is-minimum-version --trust-cert --log-level info
    }

```

See the [configuration options](#configuration-options) section for all
command-line settings supported by this plugin along with descriptions of
each. See the [contrib](#contrib) section for information regarding example
command definitions and Nagios configuration files.

### `check_vmware_hs2ds2vms` Nagios plugin

#### CLI invocation

```ShellSession
/usr/lib/nagios/plugins/check_vmware_hs2ds2vms --username SERVICE_ACCOUNT_NAME --password "SERVICE_ACCOUNT_PASSWORD" --server vc1.example.com --exclude-rp "Desktops" --ignore-vm "test1.example.com,redmine.example.com,TESTING-AC,RHEL7-TEST" --ca-name "Location" --ca-prefix-sep "-" --trust-cert --log-level info
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
- Service Check results output is sent to `stdout`
- Logging output is enabled at the `info` level.
  - logging output is sent to `stderr` by default
  - logging output is intended to be seen when invoking the
    plugin directly via CLI (often for troubleshooting)
    - see [Output](#output) for potential conflicts with some monitoring
      systems

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
- Service Check results output is sent to `stdout`
- Logging output is enabled at the `info` level.
  - logging output is sent to `stderr` by default
  - logging output is intended to be seen when invoking the
    plugin directly via CLI (often for troubleshooting)
    - see [Output](#output) for potential conflicts with some monitoring
      systems

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
- Service Check results output is sent to `stdout`
- Logging output is enabled at the `info` level.
  - logging output is sent to `stderr` by default
  - logging output is intended to be seen when invoking the
    plugin directly via CLI (often for troubleshooting)
    - see [Output](#output) for potential conflicts with some monitoring
      systems

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

### `check_vmware_snapshots_count` Nagios plugin

#### CLI invocation

```ShellSession
/usr/lib/nagios/plugins/check_vmware_snapshots_count --username SERVICE_ACCOUNT_NAME --password "SERVICE_ACCOUNT_PASSWORD" --server vc1.example.com --count-warning 4 --count-critical 25 --trust-cert --log-level info
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
- Service Check results output is sent to `stdout`
- Logging output is enabled at the `info` level.
  - logging output is sent to `stderr` by default
  - logging output is intended to be seen when invoking the
    plugin directly via CLI (often for troubleshooting)
    - see [Output](#output) for potential conflicts with some monitoring
      systems

#### Command definition

```shell
# /etc/nagios-plugins/config/vmware-snapshots-count.cfg

# Look at all pools, all VMs, do not evaluate any VMs that are powered off.
# This variation of the command is most useful for environments where all VMs
# are monitored equally.
define command{
    command_name    check_vmware_snapshots_count
    command_line    /usr/lib/nagios/plugins/check_vmware_snapshots_count --server '$HOSTNAME$' --domain '$ARG1$' --username '$ARG2$' --password '$ARG3$' --count-warning '$ARG4$' --count-critical '$ARG5$' --trust-cert --log-level info
    }
```

### `check_vmware_snapshots_size` Nagios plugin

#### CLI invocation

```ShellSession
/usr/lib/nagios/plugins/check_vmware_snapshots_size --username SERVICE_ACCOUNT_NAME --password "SERVICE_ACCOUNT_PASSWORD" --server vc1.example.com --size-warning 20 --size-critical 40 --trust-cert --log-level info
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
- Service Check results output is sent to `stdout`
- Logging output is enabled at the `info` level.
  - logging output is sent to `stderr` by default
  - logging output is intended to be seen when invoking the
    plugin directly via CLI (often for troubleshooting)
    - see [Output](#output) for potential conflicts with some monitoring
      systems

#### Command definition

```shell
# /etc/nagios-plugins/config/vmware-snapshots-size.cfg

# Look at all pools, all VMs, do not evaluate any VMs that are powered off.
# This variation of the command is most useful for environments where all VMs
# are monitored equally.
define command{
    command_name    check_vmware_snapshots_size
    command_line    /usr/lib/nagios/plugins/check_vmware_snapshots_size --server '$HOSTNAME$' --domain '$ARG1$' --username '$ARG2$' --password '$ARG3$' --size-warning '$ARG4$' --size-critical '$ARG5$' --trust-cert --log-level info
    }
```

See the [configuration options](#configuration-options) section for all
command-line settings supported by this plugin along with descriptions of
each. See the [contrib](#contrib) section for information regarding example
command definitions and Nagios configuration files.

### `check_vmware_rps_memory` Nagios plugin

#### CLI invocation

```ShellSession
/usr/lib/nagios/plugins/check_vmware_rps_memory --username SERVICE_ACCOUNT_NAME --password "SERVICE_ACCOUNT_PASSWORD" --server vc1.example.com --exclude-rp "Desktops" --memory-use-warning 80 --memory-use-critical 95  --memory-max-allowed 320 --trust-cert --log-level info
```

See the [configuration options](#configuration-options) section for all
command-line settings supported by this plugin along with descriptions of
each. See the [contrib](#contrib) section for information regarding example
command definitions and Nagios configuration files.

Of note:

- The resource pool named `Desktops` is excluded from evaluation.
  - this results in *all other* resource pools visible to the specified user
    account being used for evaluation
  - VMs *outside* of a Resource Pool (visible to the specified user account
    or not) do not contribute to memory usage calculations
- Certificate warnings are ignored.
  - not best practice, but many vCenter instances use self-signed certs per
    various freely available guides
- Service Check results output is sent to `stdout`
- Logging output is enabled at the `info` level.
  - logging output is sent to `stderr` by default
  - logging output is intended to be seen when invoking the
    plugin directly via CLI (often for troubleshooting)
    - see [Output](#output) for potential conflicts with some monitoring
      systems

#### Command definition

NOTE: This is the inverse of the command-line example for this plugin; only
specified Resource Pools are evaluated.

```shell
# /etc/nagios-plugins/config/vmware-resource-pools.cfg

# This variation of the command does not allow exclusions
define command{
    command_name    check_vmware_resource_pools_include_pools
    command_line    /usr/lib/nagios/plugins/check_vmware_rps_memory --server '$HOSTNAME$' --domain '$ARG1$' --username '$ARG2$' --password '$ARG3$' --memory-use-warning '$ARG4$' --memory-use-critical '$ARG5$' --memory-max-allowed '$ARG6$' --include-rp '$ARG7$' --trust-cert  --log-level info
    }
```

See the [configuration options](#configuration-options) section for all
command-line settings supported by this plugin along with descriptions of
each. See the [contrib](#contrib) section for information regarding example
command definitions and Nagios configuration files.

### `check_vmware_host_memory` Nagios plugin

#### CLI invocation

```ShellSession
/usr/lib/nagios/plugins/check_vmware_host_memory --username SERVICE_ACCOUNT_NAME --password "SERVICE_ACCOUNT_PASSWORD" --server vc1.example.com --host-name "esx1.example.com" --memory-usage-warning 80 --memory-usage-critical 95 --trust-cert --log-level info
```

See the [configuration options](#configuration-options) section for all
command-line settings supported by this plugin along with descriptions of
each. See the [contrib](#contrib) section for information regarding example
command definitions and Nagios configuration files.

Of note:

- The host name is specified (via `host-name` flag) using the exact value
  shown in the vSphere inventory (e.g., `esx1.example.com`)
- Certificate warnings are ignored.
  - not best practice, but many vCenter instances use self-signed certs per
    various freely available guides
- Service Check results output is sent to `stdout`
- Logging output is enabled at the `info` level.
  - logging output is sent to `stderr` by default
  - logging output is intended to be seen when invoking the
    plugin directly via CLI (often for troubleshooting)
    - see [Output](#output) for potential conflicts with some monitoring
      systems

#### Command definition

```shell
# /etc/nagios-plugins/config/vmware-host-memory.cfg

# Look at a specific host and explicitly provide custom WARNING and CRITICAL
# threshold values.
define command{
    command_name    check_vmware_host_memory
    command_line    /usr/lib/nagios/plugins/check_vmware_host_memory --server '$HOSTNAME$' --domain '$ARG1$' --username '$ARG2$' --password '$ARG3$' --memory-usage-warning '$ARG4$' --memory-usage-critical '$ARG5$' --host-name '$ARG6$' --trust-cert  --log-level info
    }
```

### `check_vmware_host_cpu` Nagios plugin

#### CLI invocation

```ShellSession
/usr/lib/nagios/plugins/check_vmware_host_cpu --username SERVICE_ACCOUNT_NAME --password "SERVICE_ACCOUNT_PASSWORD" --server vc1.example.com --host-name "esx1.example.com" --cpu-usage-warning 80 --cpu-usage-critical 95 --trust-cert --log-level info
```

See the [configuration options](#configuration-options) section for all
command-line settings supported by this plugin along with descriptions of
each. See the [contrib](#contrib) section for information regarding example
command definitions and Nagios configuration files.

Of note:

- The host name is specified (via `host-name` flag) using the exact value
  shown in the vSphere inventory (e.g., `esx1.example.com`)
- Certificate warnings are ignored.
  - not best practice, but many vCenter instances use self-signed certs per
    various freely available guides
- Service Check results output is sent to `stdout`
- Logging output is enabled at the `info` level.
  - logging output is sent to `stderr` by default
  - logging output is intended to be seen when invoking the
    plugin directly via CLI (often for troubleshooting)
    - see [Output](#output) for potential conflicts with some monitoring
      systems

#### Command definition

```shell
# /etc/nagios-plugins/config/vmware-host-cpu.cfg

# Look at a specific host and explicitly provide custom WARNING and CRITICAL
# threshold values.
define command{
    command_name    check_vmware_host_cpu
    command_line    /usr/lib/nagios/plugins/check_vmware_host_cpu --server '$HOSTNAME$' --domain '$ARG1$' --username '$ARG2$' --password '$ARG3$' --cpu-usage-warning '$ARG4$' --cpu-usage-critical '$ARG5$' --host-name '$ARG6$' --trust-cert  --log-level info
    }
```

### `check_vmware_vm_power_uptime` Nagios plugin

#### CLI invocation

```ShellSession
/usr/lib/nagios/plugins/check_vmware_vm_power_uptime --username SERVICE_ACCOUNT_NAME --password "SERVICE_ACCOUNT_PASSWORD" --server vc1.example.com --uptime-warning 60 --uptime-critical 90 --trust-cert --log-level info
```

See the [configuration options](#configuration-options) section for all
command-line settings supported by this plugin along with descriptions of
each. See the [contrib](#contrib) section for information regarding example
command definitions and Nagios configuration files.

Of note:

- Certificate warnings are ignored.
  - not best practice, but many vCenter instances use self-signed certs per
    various freely available guides
- Service Check results output is sent to `stdout`
- Logging output is enabled at the `info` level.
  - logging output is sent to `stderr` by default
  - logging output is intended to be seen when invoking the
    plugin directly via CLI (often for troubleshooting)
    - see [Output](#output) for potential conflicts with some monitoring
      systems

#### Command definition

```shell
# /etc/nagios-plugins/config/vmware-vm-power-uptime.cfg

# Look at all pools, all VMs, do not evaluate any VMs that are powered off.
# This variation of the command is most useful for environments where all VMs
# are monitored equally.
define command{
    command_name    check_vmware_vm_power_uptime
    command_line    /usr/lib/nagios/plugins/check_vmware_vm_power_uptime --server '$HOSTNAME$' --domain '$ARG1$' --username '$ARG2$' --password '$ARG3$' --uptime-warning '$ARG4$' --uptime-critical '$ARG5$' --trust-cert  --log-level info
    }
```

### `check_vmware_disk_consolidation` Nagios plugin

#### CLI invocation

```ShellSession
/usr/lib/nagios/plugins/check_vmware_disk_consolidation --username SERVICE_ACCOUNT_NAME --password "SERVICE_ACCOUNT_PASSWORD" --server vc1.example.com  --trust-cert --log-level info --trigger-reload --timeout 110
```

See the [configuration options](#configuration-options) section for all
command-line settings supported by this plugin along with descriptions of
each. See the [contrib](#contrib) section for information regarding example
command definitions and Nagios configuration files.

Of note:

- Certificate warnings are ignored.
  - not best practice, but many vCenter instances use self-signed certs per
    various freely available guides
- Service Check results output is sent to `stdout`
- Logging output is enabled at the `info` level.
  - logging output is sent to `stderr` by default
  - logging output is intended to be seen when invoking the
    plugin directly via CLI (often for troubleshooting)
    - see [Output](#output) for potential conflicts with some monitoring
      systems
- A forced state data reload/refresh is triggered for each evaluated Virtual
  Machine
  - NOTE: this operation is expensive, omit to rely on existing (potentially
    stale) state data
  - see [`check_vmware_disk_consolidation`](#check_vmware_disk_consolidation)
    section for additional details, including potential alternatives to the
    use of this flag
- A custom timeout value is specified
  - NOTE: in order for this timeout value to be respected, you may need to
    adjust the service check timeout value in your monitoring system (e.g.,
    `service_check_timeout` value in your `nagios.cfg` file)

#### Command definition

```shell
# /etc/nagios-plugins/config/vmware-disk-consolidation.cfg

# Look at all pools, all VMs.  Use existing (potentially stale) state data for
# evaluation of disk consolidation status instead of triggering (potentially
# expensive) reload/refresh of state data.
#
# This variation of the command is most useful for environments where all VMs
# are monitored equally and no filtering based on pool membership or VM name
# is needed.
define command{
    command_name    check_vmware_disk_consolidation
    command_line    /usr/lib/nagios/plugins/check_vmware_disk_consolidation --server '$HOSTNAME$' --domain '$ARG1$' --username '$ARG2$' --password '$ARG3$'  --trust-cert --log-level info
    }

# Look at all pools, all VMs, trigger potentially expensive reload operation
# on each evaluated VM.
#
# This variation of the command is most useful for environments where all VMs
# are monitored equally and where the time required to reload/refresh data
# data for each VM is acceptable.
#
# The tradeoff in having current state data comes at the cost of increased
# execution time. If this proves too expensive for your environment, you may
# wish to schedule a job on the cluster to handle refreshing state data.
define command{
    command_name    check_vmware_disk_consolidation_trigger_reload
    command_line    /usr/lib/nagios/plugins/check_vmware_disk_consolidation --server '$HOSTNAME$' --domain '$ARG1$' --username '$ARG2$' --password '$ARG3$'  --trust-cert --log-level info --trigger-reload --timeout 110
    }

```

### `check_vmware_question` Nagios plugin

#### CLI invocation

```ShellSession
/usr/lib/nagios/plugins/check_vmware_question --username SERVICE_ACCOUNT_NAME --password "SERVICE_ACCOUNT_PASSWORD" --server vc1.example.com  --trust-cert --log-level info
```

See the [configuration options](#configuration-options) section for all
command-line settings supported by this plugin along with descriptions of
each. See the [contrib](#contrib) section for information regarding example
command definitions and Nagios configuration files.

Of note:

- Certificate warnings are ignored.
  - not best practice, but many vCenter instances use self-signed certs per
    various freely available guides
- Service Check results output is sent to `stdout`
- Logging output is enabled at the `info` level.
  - logging output is sent to `stderr` by default
  - logging output is intended to be seen when invoking the
    plugin directly via CLI (often for troubleshooting)
    - see [Output](#output) for potential conflicts with some monitoring
      systems

#### Command definition

```shell
# /etc/nagios-plugins/config/vmware-interactive-question.cfg

# Look at all pools, all VMs. This variation of the command is most useful for
# environments where all VMs are monitored equally.
define command{
    command_name    check_vmware_question
    command_line    /usr/lib/nagios/plugins/check_vmware_question --server '$HOSTNAME$' --domain '$ARG1$' --username '$ARG2$' --password '$ARG3$'  --trust-cert --log-level info
    }
```

### `check_vmware_alarms` Nagios plugin

#### CLI invocation

```ShellSession
/usr/lib/nagios/plugins/check_vmware_alarms --username SERVICE_ACCOUNT_NAME --password "SERVICE_ACCOUNT_PASSWORD" --server vc1.example.com  --trust-cert --log-level info
```

See the [configuration options](#configuration-options) section for all
command-line settings supported by this plugin along with descriptions of
each. See the [contrib](#contrib) section for information regarding example
command definitions and Nagios configuration files.

Of note:

- Triggered alarms are evaluated for all detected datacenters
  - due to lack of specified datacenter name (or names)
- Triggered alarms are not filtered based on associated [managed
  object][vsphere-managed-object-reference] (aka, `managed entity`) type
  - due to lack of explicit exclusions or inclusions
- Triggered alarms are not filtered based on associated [managed
  object][vsphere-managed-object-reference] (aka, `managed entity`) name
  - due to lack of explicit exclusions or inclusions
- Triggered alarms are not filtered based on associated [managed
  object][vsphere-managed-object-reference] (aka, `managed entity`) resource
  pool
  - due to lack of explicit exclusions or inclusions
- Triggered alarms that were previously acknowledged are ignored
- Triggered alarms are *not* filtered based on defined Alarm name
  - due to lack of explicit exclusions or inclusions
- Triggered alarms are *not* filtered based on defined Alarm description
  - due to lack of explicit exclusions or inclusions
- Triggered alarms are *not* filtered based on Triggered Alarm status
  - due to lack of explicit exclusions or inclusions
- Certificate warnings are ignored.
  - not best practice, but many vCenter instances use self-signed certs per
    various freely available guides
- Service Check results output is sent to `stdout`
- Logging output is enabled at the `info` level.
  - logging output is sent to `stderr` by default
  - logging output is intended to be seen when invoking the
    plugin directly via CLI (often for troubleshooting)
    - see [Output](#output) for potential conflicts with some monitoring
      systems

#### Command definition

```shell
# /etc/nagios-plugins/config/vmware-alarms.cfg

# Look at triggered alarms across all detected datacenters, do not evaluate
# any triggered alarms which have been previously acknowledged.
define command{
    command_name    check_vmware_alarms
    command_line    /usr/lib/nagios/plugins/check_vmware_alarms --server '$HOSTNAME$' --domain '$ARG1$' --username '$ARG2$' --password '$ARG3$' --trust-cert --log-level info
    }

# Look at triggered alarms within specified datacenters. Do not evaluate any
# triggered alarms which have been previously acknowledged.
define command{
    command_name    check_vmware_alarms_specific_dc
    command_line    /usr/lib/nagios/plugins/check_vmware_alarms --server '$HOSTNAME$' --domain '$ARG1$' --username '$ARG2$' --password '$ARG3$' --dc-name '$ARG4$' --trust-cert --log-level info
    }
```

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
  - VMware Tools | VersionStatus field
    - [vCenter Data Structures](https://developer.vmware.com/docs/vsphere-automation/latest/vcenter/data-structures/Vm/Tools/VersionStatus/)
    - [vCenter Storage Monitoring Service API Reference](https://vdc-repo.vmware.com/vmwb-repository/dcr-public/7989f521-fd57-4fff-9653-e6a5d5265089/1fd5908d-b8ce-49ca-887a-fefb3656e828/doc/vim.vm.GuestInfo.ToolsVersionStatus.html)

- Logging
  - <https://github.com/rs/zerolog>

- Nagios
  - <https://github.com/atc0005/go-nagios>
  - <https://nagios-plugins.org/doc/guidelines.html>
  - <https://www.monitoring-plugins.org/doc/guidelines.html>
  - <https://icinga.com/docs/icinga-2/latest/doc/05-service-monitoring/>

<!-- Footnotes here  -->

[repo-url]: <https://github.com/atc0005/check-vmware>  "This project's GitHub repo"

[go-docs-download]: <https://golang.org/dl>  "Download Go"

[go-docs-install]: <https://golang.org/doc/install>  "Install Go"

[vsphere-managed-object-reference]: <https://vdc-download.vmware.com/vmwb-repository/dcr-public/a5f4000f-1ea8-48a9-9221-586adff3c557/7ff50256-2cf2-45ea-aacd-87d231ab1ac7/vmodl.ManagedObjectReference.html> "Data Object - ManagedObjectReference(vmodl.ManagedObjectReference)"

[vsphere-manged-entity-status]: <https://vdc-repo.vmware.com/vmwb-repository/dcr-public/91f5f971-bf1d-4904-9942-37c6109da8a3/b79fa83f-dc4e-491d-9785-dc9d91aa0c67/doc/vim.ManagedEntity.Status.html>

[vsphere-default-alarms]: <https://docs.vmware.com/en/VMware-vSphere/7.0/com.vmware.vsphere.monitoring.doc/GUID-82933270-1D72-4CF3-A1AF-E5A1343F62DE.html>

[nagios-state-types]: <https://assets.nagios.com/downloads/nagioscore/docs/nagioscore/3/en/statetypes.html>

[vsphere-guestinfo-data-object]: <https://vdc-download.vmware.com/vmwb-repository/dcr-public/b50dcbbf-051d-4204-a3e7-e1b618c1e384/538cf2ec-b34f-4bae-a332-3820ef9e7773/vim.vm.GuestInfo.html>

<!-- []: PLACEHOLDER "DESCRIPTION_HERE" -->
