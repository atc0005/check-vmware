<!-- omit in toc -->
# check-vmware

Go-based tooling to monitor VMware environments; **NOT** affiliated with
or endorsed by VMware, Inc.

[![Latest Release](https://img.shields.io/github/release/atc0005/check-vmware.svg?style=flat-square)](https://github.com/atc0005/check-vmware/releases/latest)
[![Go Reference](https://pkg.go.dev/badge/github.com/atc0005/check-vmware.svg)](https://pkg.go.dev/github.com/atc0005/check-vmware)
[![go.mod Go version](https://img.shields.io/github/go-mod/go-version/atc0005/check-vmware)](https://github.com/atc0005/check-vmware)
[![Validate Codebase](https://github.com/atc0005/check-vmware/workflows/Validate%20Codebase/badge.svg)](https://github.com/atc0005/check-vmware/actions?query=workflow%3A%22Validate+Codebase%22)
[![Validate Docs](https://github.com/atc0005/check-vmware/workflows/Validate%20Docs/badge.svg)](https://github.com/atc0005/check-vmware/actions?query=workflow%3A%22Validate+Docs%22)
[![Lint and Build using Makefile](https://github.com/atc0005/check-vmware/workflows/Lint%20and%20Build%20using%20Makefile/badge.svg)](https://github.com/atc0005/check-vmware/actions?query=workflow%3A%22Lint+and+Build+using+Makefile%22)
[![Quick Validation](https://github.com/atc0005/check-vmware/workflows/Quick%20Validation/badge.svg)](https://github.com/atc0005/check-vmware/actions?query=workflow%3A%22Quick+Validation%22)

<!-- omit in toc -->
## Table of Contents

- [Project home](#project-home)
- [Overview](#overview)
  - [Plugin index](#plugin-index)
  - [Output](#output)
  - [Performance Data](#performance-data)
  - [Optional evaluation](#optional-evaluation)
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
- [Contrib](#contrib)
- [Examples](#examples)
- [License](#license)
- [References](#references)

## Project home

See [our GitHub repo][repo-url] for the latest code, to file an issue or
submit improvements for review and potential inclusion into the project.

Just to be 100% clear: this project is not affiliated with or endorsed by
VMware, Inc.

## Overview

This repo contains various tools and plugins used to monitor/validate VMware
environments. Documentation for the project as a whole is available in this
README file while details specific to a plugin are recorded in a dedicated
file. See the [plugin index](#plugin-index) for a quick reference.

### Plugin index

| Plugin or Tool Name                                                                        | Description                                                                         |
| ------------------------------------------------------------------------------------------ | ----------------------------------------------------------------------------------- |
| [`check_vmware_tools`](docs/plugins/check_vmware_tools.md)                                 | Nagios plugin used to monitor VMware Tools installations.                           |
| [`check_vmware_vcpus`](docs/plugins/check_vmware_vcpus.md)                                 | Nagios plugin used to monitor allocation of virtual CPUs (vCPUs).                   |
| [`check_vmware_vhw`](docs/plugins/check_vmware_vhw.md)                                     | Nagios plugin used to monitor virtual hardware versions.                            |
| [`check_vmware_hs2ds2vms`](docs/plugins/check_vmware_hs2ds2vms.md)                         | Nagios plugin used to monitor host/datastore/vm pairings.                           |
| [`check_vmware_datastore_space`](docs/plugins/check_vmware_datastore_space.md)             | Nagios plugin used to monitor datastore usage.                                      |
| [`check_vmware_datastore_performance`](docs/plugins/check_vmware_datastore_performance.md) | Nagios plugin used to monitor datastore performance.                                |
| [`check_vmware_snapshots_age`](docs/plugins/check_vmware_snapshots_age.md)                 | Nagios plugin used to monitor the age of Virtual Machine snapshots.                 |
| [`check_vmware_snapshots_count`](docs/plugins/check_vmware_snapshots_count.md)             | Nagios plugin used to monitor the count of Virtual Machine snapshots.               |
| [`check_vmware_snapshots_size`](docs/plugins/check_vmware_snapshots_size.md)               | Nagios plugin used to monitor the **cumulative** size of Virtual Machine snapshots. |
| [`check_vmware_rps_memory`](docs/plugins/check_vmware_rps_memory.md)                       | Nagios plugin used to monitor memory usage across Resource Pools.                   |
| [`check_vmware_host_memory`](docs/plugins/check_vmware_host_memory.md)                     | Nagios plugin used to monitor memory usage for a specific ESXi host system.         |
| [`check_vmware_host_cpu`](docs/plugins/check_vmware_host_cpu.md)                           | Nagios plugin used to monitor CPU usage for a specific ESXi host system.            |
| [`check_vmware_vm_power_uptime`](docs/plugins/check_vmware_vm_power_uptime.md)             | Nagios plugin used to monitor VM power cycle uptime.                                |
| [`check_vmware_disk_consolidation`](docs/plugins/check_vmware_disk_consolidation.md)       | Nagios plugin used to monitor VM disk consolidation status.                         |
| [`check_vmware_question`](docs/plugins/check_vmware_question.md)                           | Nagios plugin used to monitor VM interactive question status.                       |
| [`check_vmware_alarms`](docs/plugins/check_vmware_alarms.md)                               | Nagios plugin used to monitor for Triggered Alarms in one or more datacenters.      |

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

Initial support has been added for emitting Performance Data / Metrics, but
refinement suggestions are welcome.

Consult the list of available metrics for each plugin for details. See the
[plugin index](#plugin-index) for a quick reference of available plugins.

### Optional evaluation

Some plugins provide optional support to limit evaluation of VMs to specific
Resource Pools (explicitly including or excluding) and power states (on or
off). Other plugins support similar filtering options (e.g., `Acknowledged`
state of Triggered Alarms). See the [configuration
options](#configuration-options), [examples](#examples) and
[contrib](#contrib) sections for more information.

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
  - Datastore performance
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

- Go
  - see this project's `go.mod` file for *preferred* version
  - see upstream `vmware/govmomi` library for required version
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
     - `go build -mod=vendor ./cmd/check_vmware_datastore_space/`
     - `go build -mod=vendor ./cmd/check_vmware_datastore_performance/`
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
     - look in `/tmp/check-vmware/release_assets/check_vmware_datastore_space/`
     - look in `/tmp/check-vmware/release_assets/check_vmware_datastore_performance/`
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

See the [plugin index](#plugin-index) for a quick reference to each plugin's
documentation.

## Contrib

Example Nagios configuration files are provided in an effort to illustrate
usage of plugins provided by this project. See the [Contrib
README](contrib/README.md) and [directory contents](./contrib/) for details.

## Examples

See the [plugin index](#plugin-index) for a quick reference to each plugin's
documentation.

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
  - <https://github.com/atc0005/check-illiad>
  - <https://github.com/atc0005/check-mail>
  - <https://github.com/atc0005/check-path>
  - <https://github.com/atc0005/check-statuspage>
  - <https://github.com/atc0005/check-whois>
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

[vsphere-query-datastore-performance-summary-method]: <https://vdc-download.vmware.com/vmwb-repository/dcr-public/bf660c0a-f060-46e8-a94d-4b5e6ffc77ad/208bc706-e281-49b6-a0ce-b402ec19ef82/SDK/vsphere-ws/docs/ReferenceGuide/vim.StorageResourceManager.html#queryDatastorePerformanceSummary>

[vsphere-storage-performance-summary-data-object]: <https://vdc-download.vmware.com/vmwb-repository/dcr-public/b50dcbbf-051d-4204-a3e7-e1b618c1e384/538cf2ec-b34f-4bae-a332-3820ef9e7773/vim.StorageResourceManager.StoragePerformanceSummary.html>

[vsphere-storage-io-resource-management-data-object]: <https://vdc-download.vmware.com/vmwb-repository/dcr-public/b50dcbbf-051d-4204-a3e7-e1b618c1e384/538cf2ec-b34f-4bae-a332-3820ef9e7773/vim.StorageResourceManager.IORMConfigInfo.html>

<!-- []: PLACEHOLDER "DESCRIPTION_HERE" -->
