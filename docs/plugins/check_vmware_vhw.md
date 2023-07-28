<!-- omit in toc -->
# [check-vmware][repo-url] | `check_vmware_vhw` plugin

- [Main project README](../../README.md)
- [Documentation index](../README.md)

<!-- omit in toc -->
## Table of Contents

- [Overview](#overview)
  - [Homogeneous version check](#homogeneous-version-check)
  - [Outdated-by or threshold range check](#outdated-by-or-threshold-range-check)
  - [Minimum required version check](#minimum-required-version-check)
  - [Default is minimum required version check](#default-is-minimum-required-version-check)
- [Output](#output)
- [Performance Data](#performance-data)
  - [Background](#background)
  - [Supported metrics](#supported-metrics)
- [Optional evaluation](#optional-evaluation)
- [Installation](#installation)
- [Configuration options](#configuration-options)
  - [Threshold calculations](#threshold-calculations)
    - [Homogeneous version check](#homogeneous-version-check-1)
    - [Outdated-by or threshold range check](#outdated-by-or-threshold-range-check-1)
    - [Minimum required version check](#minimum-required-version-check-1)
    - [Default is minimum required version check](#default-is-minimum-required-version-check-1)
  - [Command-line arguments](#command-line-arguments)
  - [Configuration file](#configuration-file)
- [Contrib](#contrib)
- [Examples](#examples)
  - [Homogeneous version check](#homogeneous-version-check-2)
    - [CLI invocation](#cli-invocation)
    - [Command definition](#command-definition)
  - [Outdated-by or threshold range check](#outdated-by-or-threshold-range-check-2)
    - [CLI invocation](#cli-invocation-1)
    - [Command definition](#command-definition-1)
  - [Minimum required version check](#minimum-required-version-check-2)
    - [CLI invocation](#cli-invocation-2)
    - [Command definition](#command-definition-2)
  - [Default is minimum required version check](#default-is-minimum-required-version-check-2)
    - [CLI invocation](#cli-invocation-3)
    - [Command definition](#command-definition-3)
- [License](#license)
- [References](#references)

## Overview

Nagios plugin used to monitor virtual hardware versions.

This plugin supports four independent monitoring modes; only one mode can be
used at a time.

1. Homogeneous version check
1. Outdated-by or threshold range check
1. Minimum required version check
1. Default is minimum required version check

### Homogeneous version check

This monitoring mode applies an automatic baseline of "highest version
discovered" across evaluated VMs. Any VMs with a hardware version not at that
highest version are flagged as problematic.

Instead of trying to determine how far behind each VM is from the newest
version, this monitoring mode assumes that any deviation is a `WARNING` state.

### Outdated-by or threshold range check

This mode applies the standard WARNING and CRITICAL level threshold checks to
determine the current plugin state. Any VM with virtual hardware older than
the specified thresholds triggers the associated state. This mode is useful
for catching VMs with outdated hardware outside of an acceptable range.

The highest version used as a baseline for comparison is provided using the
same logic as provided by the "homogeneous" version check: latest visible
hardware version.

### Minimum required version check

This mode requires that all hardware versions match or exceed the specified
minimum hardware version. This monitoring mode assumes that any deviation is
considered a `CRITICAL` state.

### Default is minimum required version check

This mode requires that all hardware versions match or exceed the host or
cluster default hardware version. This monitoring mode assumes that any
deviation is considered a `WARNING` state.

## Output

The output for these plugins is designed to provide the one-line summary
needed by Nagios for quick identification of a problem while providing longer,
more detailed information for display within the web UI, use in email and
Teams notifications
([atc0005/send2teams](https://github.com/atc0005/send2teams)).

See the [main project README](../../README.md) for details.

## Performance Data

### Background

Initial support has been added for emitting Performance Data / Metrics, but
refinement suggestions are welcome.

Consult the list below for the metrics implemented thus far, [the original
discussion thread](https://github.com/atc0005/check-vmware/discussions/315)
and the [Add Performance Data / Metrics
support](https://github.com/atc0005/check-vmware/projects/1) project board for
an index of the initial implementation work.

Please add to an existing
[Discussion](https://github.com/atc0005/check-vmware/discussions) thread or
[open a new one](https://github.com/atc0005/check-vmware/discussions/new) with
any feedback that you may have. Thanks in advance!

### Supported metrics

Metrics below are obtained in this order:

1. Obtain count of all resource pools
1. Obtain count of all folders
1. Obtain count of all virtual machines
1. Filter virtual machines
   1. by resource pools
   1. by folders
   1. by name
   1. by power state
1. Evaluate virtual machine virtual hardware versions

For example, the count of virtual machines powered on is obtained based on VMs
remaining after resource pool filtering is complete at the time of applying
power state filtering.

**NOTE**: These metrics are based on the visibility of the service account
used to login to the target VMware environment. If the service account cannot
see a resource, it cannot evaluate the resource.

| Metric                          | Alias of              | Unit of Measurement | Description                                                                              |
| ------------------------------- | --------------------- | ------------------- | ---------------------------------------------------------------------------------------- |
| `time`                          |                       | milliseconds        | plugin runtime                                                                           |
| `vms`                           | `vms_all`             |                     | all (visible) virtual machines in the inventory                                          |
| `vms_all`                       | `vms`                 |                     | all (visible) virtual machines in the inventory                                          |
| `vms_evaluated`                 | `vms_after_filtering` |                     | virtual machines after filtering, evaluated for plugin-specific threshold violations     |
| `vms_after_filtering`           | `vms_evaluated`       |                     | virtual machines after filtering, evaluated for plugin-specific threshold violations     |
| `vms_powered_on`                |                       |                     | virtual machines powered on                                                              |
| `vms_powered_off`               |                       |                     | virtual machines powered off                                                             |
| `vms_excluded_by_name`          |                       |                     | virtual machines excluded based on fixed name values                                     |
| `vms_excluded_by_folder`        |                       |                     | virtual machines excluded based on folder IDs                                            |
| `vms_excluded_by_power_state`   |                       |                     | virtual machines excluded based on power state (powered off VMs are excluded by default) |
| `vms_excluded_by_resource_pool` |                       |                     | virtual machines excluded based on resource pool name                                    |
| `folders_all`                   |                       |                     | all folders in the inventory                                                             |
| `folders_excluded`              |                       |                     | folders excluded by request                                                              |
| `folders_included`              |                       |                     | folders included by request (all non-listed folders excluded)                            |
| `folders_evaluated`             |                       |                     | folders remaining after inclusion/exclusion filtering logic is applied                   |
| `resource_pools_all`            |                       |                     | all resource pools in the inventory                                                      |
| `resource_pools_excluded`       |                       |                     | resource pools excluded by request                                                       |
| `resource_pools_included`       |                       |                     | resource pools included by request (all non-listed resource pools excluded)              |
| `resource_pools_evaluated`      |                       |                     | resource pools remaining after inclusion/exclusion filtering logic is applied            |
| `hardware_versions_unique`      |                       |                     | virtual machines with unique virtual machine hardware versions                           |
| `hardware_versions_newest`      |                       |                     | virtual machines with the newest virtual machine hardware version                        |
| `hardware_versions_default`     |                       |                     | virtual machines with the default cluster hardware version                               |
| `hardware_versions_oldest`      |                       |                     | virtual machines with the oldest hardware version                                        |

## Optional evaluation

Some plugins provide optional support to limit evaluation of VMs to specific
Resource Pools (explicitly including or excluding) and power states (on or
off). Other plugins support similar filtering options (e.g., `Acknowledged`
state of Triggered Alarms). See the [configuration
options](#configuration-options), [examples](#examples) and
[contrib](#contrib) sections for more information.

## Installation

See the [main project README](../../README.md) for details.

## Configuration options

### Threshold calculations

This plugin supports multiple (independent) modes. Each mode applies slightly
different logic for determining plugin state. Only one mode can be used at a
time.

#### Homogeneous version check

| Nagios State | Description                                |
| ------------ | ------------------------------------------ |
| `OK`         | Ideal state, homogenous hardware versions. |
| `WARNING`    | Non-homogenous hardware versions.          |
| `CRITICAL`   | Not used by this monitoring mode.          |

#### Outdated-by or threshold range check

| Nagios State | Description                                                        |
| ------------ | ------------------------------------------------------------------ |
| `OK`         | Ideal state, hardware versions within tolerance.                   |
| `WARNING`    | Hardware versions crossed user-specified threshold for this state. |
| `CRITICAL`   | Hardware versions crossed user-specified threshold for this state. |

#### Minimum required version check

| Nagios State | Description                                                       |
| ------------ | ----------------------------------------------------------------- |
| `OK`         | Ideal state, hardware versions within tolerance.                  |
| `WARNING`    | Not used by this monitoring mode.                                 |
| `CRITICAL`   | Hardware versions older than the minimum specified value present. |

#### Default is minimum required version check

| Nagios State | Description                                                             |
| ------------ | ----------------------------------------------------------------------- |
| `OK`         | Ideal state, hardware versions within tolerance.                        |
| `WARNING`    | Hardware versions older than the host or cluster default value present. |
| `CRITICAL`   | Not used by this monitoring mode.                                       |

### Command-line arguments

This plugin supports multiple (independent) monitoring modes. Only one mode
can be used at a time and each mode has options which are incompatible with
the others.

As of this writing, these monitoring modes are *not* implemented as
subcommands, though this may change in the future based on feedback. See the
[examples](#examples) for this plugin for more information.

- Use the `-h` or `--help` flag to display current usage information.
- Flags marked as **`required`** must be set via CLI flag.
- Flags *not* marked as required are for settings where a useful default is
  already defined, but may be overridden if desired.

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
| `domain`                         | No        |         | No     | *valid user domain*                                                     | (Optional) domain for user account used to login to ESXi host or vCenter instance. This is needed for user accounts residing in a non-default domain (e.g., SSO specific domain).                                                                                                                                                                                                             |
| `trust-cert`                     | No        | `false` | No     | `true`, `false`                                                         | Whether the certificate should be trusted as-is without validation. WARNING: TLS is susceptible to man-in-the-middle attacks if enabling this option.                                                                                                                                                                                                                                         |
| `dc-name`                        | No        |         | No     | *valid vSphere datacenter name*                                         | Specifies the name of a vSphere Datacenter. If not specified, applicable plugins will attempt to use the default datacenter found in the vSphere environment. Not applicable to standalone ESXi hosts.                                                                                                                                                                                        |
| `host-name`                      | No        |         | No     | *valid ESXi host name*                                                  | ESXi host/server name as it is found within the vSphere inventory.                                                                                                                                                                                                                                                                                                                            |
| `cluster-name`                   | No        |         | No     | *valid vSphere cluster name*                                            | Specifies the name of a vSphere Cluster. If not specified, applicable plugins will attempt to use the default cluster found in the vSphere environment. Not applicable to standalone ESXi hosts.                                                                                                                                                                                              |
| `include-rp`                     | No        |         | No     | *comma-separated list of resource pool names*                           | Specifies a comma-separated list of Resource Pool names that should be exclusively used when evaluating VMs. Specifying this option will also exclude any VMs from evaluation that are *outside* of a Resource Pool. This option is incompatible with specifying a list of Resource Pool names to ignore or exclude from evaluation.                                                          |
| `exclude-rp`                     | No        |         | No     | *comma-separated list of resource pool names*                           | Specifies a comma-separated list of Resource Pool names that should be ignored when evaluating VMs. This option is incompatible with specifying a list of Resource Pool names to include for evaluation.                                                                                                                                                                                      |
| `include-folder-id`              | No        |         | No     | *comma-separated list of folder ID values*                              | Specifies a comma-separated list of Folder Managed Object ID (MOID) values (e.g., group-v34) that should be exclusively used when evaluating VMs. This option is incompatible with specifying a list of Folder IDs to ignore or exclude from evaluation.                                                                                                                                      |
| `exclude-folder-id`              | No        |         | No     | *comma-separated list of folder ID values*                              | Specifies a comma-separated list of Folder Managed Object ID (MOID) values (e.g., group-v34) that should be ignored when evaluating VMs. This option is incompatible with specifying a list of Folder Managed Object ID (MOID) values to include for evaluation.                                                                                                                              |
| `ignore-vm`                      | No        |         | No     | *comma-separated list of (vSphere) virtual machine names*               | Specifies a comma-separated list of VM names that should be ignored or excluded from evaluation.                                                                                                                                                                                                                                                                                              |
| `powered-off`                    | No        | `false` | No     | `true`, `false`                                                         | Toggles evaluation of powered off VMs in addition to powered on VMs. Evaluation of powered off VMs is disabled by default.                                                                                                                                                                                                                                                                    |
| `obw`, `outdated-by-warning`     | **Maybe** |         | No     | *positive whole number 1 or greater*                                    | If provided, this value is the WARNING threshold for outdated virtual hardware versions. If the current virtual hardware version for a VM is found to be more than this many versions older than the latest version a WARNING state is triggered. Required if specifying the CRITICAL threshold for outdated virtual hardware versions, incompatible with the minimum required version flag.  |
| `obc`, `outdated-by-critical`    | **Maybe** |         | No     | *positive whole number 1 or greater*                                    | If provided, this value is the CRITICAL threshold for outdated virtual hardware versions. If the current virtual hardware version for a VM is found to be more than this many versions older than the latest version a CRITICAL state is triggered. Required if specifying the WARNING threshold for outdated virtual hardware versions, incompatible with the minimum required version flag. |
| `mv`, `minimum-version`          | **Maybe** |         | No     | *positive whole number greater than 3*                                  | If provided, this value is the minimum virtual hardware version accepted for each Virtual Machine. Any Virtual Machine not meeting this minimum value is considered to be in a CRITICAL state. Per [KB 1003746](https://kb.vmware.com/s/article/1003746), version 3 appears to be the oldest version supported. Incompatible with the CRITICAL and WARNING threshold flags.                   |
| `dimv`, `default-is-min-version` | **Maybe** |         | No     | *positive whole number greater than 3*                                  | If provided, this value is the minimum virtual hardware version accepted for each Virtual Machine. Any Virtual Machine not meeting this minimum value is considered to be in a CRITICAL state. Per [KB 1003746](https://kb.vmware.com/s/article/1003746), version 3 appears to be the oldest version supported. Incompatible with the CRITICAL and WARNING threshold flags.                   |

### Configuration file

Not currently supported. This feature may be added later if there is
sufficient interest.

## Contrib

See the [main project README](../../README.md) for details.

## Examples

While entries in this section attempt to provide a brief overview of usage, it
is recommended that you review the provided command definitions and other
Nagios configuration files within the [`contrib`](#contrib) directory for more
complete examples.

See the [configuration options](#configuration-options) section for all
command-line settings supported by this plugin along with descriptions of
each.

This plugin supports four independent monitoring modes; only one mode can be
used at a time. Due to this, it is beneficial to have separate command
definitions for each. See the examples for each mode below or the
[overview](#overview) section for further information.

### Homogeneous version check

#### CLI invocation

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
  - logging output is intended to be seen when invoking the plugin directly
    via CLI (often for troubleshooting)
    - see the [Output section](../../README.md#output) of the main README for
      potential conflicts with some monitoring systems

#### Command definition

```shell
# /etc/nagios-plugins/config/vmware-virtual-hardware.cfg

# Look at all pools, all VMs, do not evaluate any VMs that are powered off.
# This variation of the command is most useful for environments where all VMs
# are monitored equally.
define command{
    command_name    check_vmware_vhw_homogeneous
    command_line    $USER1$/check_vmware_vhw --server '$HOSTNAME$' --domain '$ARG1$' --username '$ARG2$' --password '$ARG3$' --trust-cert --log-level info
    }
```

See the [configuration options](#configuration-options) section for all
command-line settings supported by this plugin along with descriptions of
each. See the [contrib](#contrib) section for information regarding example
command definitions and Nagios configuration files.

### Outdated-by or threshold range check

#### CLI invocation

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
  - logging output is intended to be seen when invoking the plugin directly
    via CLI (often for troubleshooting)
    - see the [Output section](../../README.md#output) of the main README for
      potential conflicts with some monitoring systems

#### Command definition

```shell
# /etc/nagios-plugins/config/vmware-virtual-hardware.cfg

# Look at all pools, all VMs, do not evaluate any VMs that are powered off.
# This variation of the command is most useful for environments where all VMs
# are monitored equally.
define command{
    command_name    check_vmware_vhw_thresholds
    command_line    $USER1$/check_vmware_vhw --server '$HOSTNAME$' --domain '$ARG1$' --username '$ARG2$' --password '$ARG3$' --outdated-by-warning '$ARG4$' --outdated-by-critical '$ARG5$' --trust-cert --log-level info
    }

```

See the [configuration options](#configuration-options) section for all
command-line settings supported by this plugin along with descriptions of
each. See the [contrib](#contrib) section for information regarding example
command definitions and Nagios configuration files.

### Minimum required version check

#### CLI invocation

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
  - logging output is intended to be seen when invoking the plugin directly
    via CLI (often for troubleshooting)
    - see the [Output section](../../README.md#output) of the main README for
      potential conflicts with some monitoring systems

#### Command definition

```shell
# /etc/nagios-plugins/config/vmware-virtual-hardware.cfg

# Look at all pools, all VMs, do not evaluate any VMs that are powered off.
# This variation of the command is most useful for environments where all VMs
# are monitored equally.
define command{
    command_name    check_vmware_vhw_minreq
    command_line    $USER1$/check_vmware_vhw --server '$HOSTNAME$' --domain '$ARG1$' --username '$ARG2$' --password '$ARG3$' --minimum-version '$ARG4$' --trust-cert --log-level info
    }

```

See the [configuration options](#configuration-options) section for all
command-line settings supported by this plugin along with descriptions of
each. See the [contrib](#contrib) section for information regarding example
command definitions and Nagios configuration files.

### Default is minimum required version check

#### CLI invocation

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
  - logging output is intended to be seen when invoking the plugin directly
    via CLI (often for troubleshooting)
    - see the [Output section](../../README.md#output) of the main README for
      potential conflicts with some monitoring systems

#### Command definition

```shell
# /etc/nagios-plugins/config/vmware-virtual-hardware.cfg

# Look at all pools, all VMs, do not evaluate any VMs that are powered off.
# This variation of the command is most useful for environments where all VMs
# are monitored equally.
define command{
    command_name    check_vmware_vhw_defreq
    command_line    $USER1$/check_vmware_vhw --server '$HOSTNAME$' --domain '$ARG1$' --username '$ARG2$' --password '$ARG3$' --cluster-name '$ARG4$' --default-is-minimum-version --trust-cert --log-level info
    }

```

See the [configuration options](#configuration-options) section for all
command-line settings supported by this plugin along with descriptions of
each. See the [contrib](#contrib) section for information regarding example
command definitions and Nagios configuration files.

## License

See the [main project README](../../README.md) for details.

## References

- [Main project README](../../README.md)
- [Documentation index](../README.md)
- [Project repo][repo-url]

<!-- Footnotes here  -->

[repo-url]: <https://github.com/atc0005/check-vmware>  "This project's GitHub repo"

<!-- []: PLACEHOLDER "DESCRIPTION_HERE" -->
