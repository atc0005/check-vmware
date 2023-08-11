<!-- omit in toc -->
# [check-vmware][repo-url] | `check_vmware_power_uptime` plugin

- [Main project README](../../README.md)
- [Documentation index](../README.md)

<!-- omit in toc -->
## Table of Contents

- [Overview](#overview)
- [Output](#output)
- [Performance Data](#performance-data)
  - [Background](#background)
  - [Supported metrics](#supported-metrics)
- [Optional evaluation](#optional-evaluation)
- [Installation](#installation)
- [Configuration options](#configuration-options)
  - [Threshold calculations](#threshold-calculations)
  - [Command-line arguments](#command-line-arguments)
  - [Configuration file](#configuration-file)
- [Contrib](#contrib)
- [Examples](#examples)
  - [CLI invocation](#cli-invocation)
  - [Command definition](#command-definition)
- [License](#license)
- [References](#references)

## Overview

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
1. Evaluate virtual machines for power uptime issues

For example, the count of virtual machines powered on is obtained based on VMs
remaining after resource pool filtering is complete at the time of applying
power state filtering.

**NOTE**: These metrics are based on the visibility of the service account
used to login to the target VMware environment. If the service account can't
see a resource, it cannot evaluate the resource.

| Metric                           | Alias of              | Unit of Measurement | Description                                                                              |
| -------------------------------- | --------------------- | ------------------- | ---------------------------------------------------------------------------------------- |
| `time`                           |                       | milliseconds        | plugin runtime                                                                           |
| `vms`                            | `vms_all`             |                     | all (visible) virtual machines in the inventory                                          |
| `vms_all`                        | `vms`                 |                     | all (visible) virtual machines in the inventory                                          |
| `vms_evaluated`                  | `vms_after_filtering` |                     | virtual machines after filtering, evaluated for plugin-specific threshold violations     |
| `vms_after_filtering`            | `vms_evaluated`       |                     | virtual machines after filtering, evaluated for plugin-specific threshold violations     |
| `vms_powered_on`                 |                       |                     | virtual machines powered on                                                              |
| `vms_powered_off`                |                       |                     | virtual machines powered off                                                             |
| `vms_excluded_by_name`           |                       |                     | virtual machines excluded based on fixed name values                                     |
| `vms_excluded_by_folder`         |                       |                     | virtual machines excluded based on folder IDs                                            |
| `vms_excluded_by_power_state`    |                       |                     | virtual machines excluded based on power state (powered off VMs are excluded by default) |
| `vms_excluded_by_resource_pool`  |                       |                     | virtual machines excluded based on resource pool name                                    |
| `folders_all`                    |                       |                     | all folders in the inventory                                                             |
| `folders_excluded`               |                       |                     | folders excluded by request                                                              |
| `folders_included`               |                       |                     | folders included by request (all non-listed folders excluded)                            |
| `folders_evaluated`              |                       |                     | folders remaining after inclusion/exclusion filtering logic is applied                   |
| `resource_pools_all`             |                       |                     | all resource pools in the inventory                                                      |
| `resource_pools_excluded`        |                       |                     | resource pools excluded by request                                                       |
| `resource_pools_included`        |                       |                     | resource pools included by request (all non-listed resource pools excluded)              |
| `resource_pools_evaluated`       |                       |                     | resource pools remaining after inclusion/exclusion filtering logic is applied            |
| `vms_with_critical_power_uptime` |                       |                     | virtual machines with a power uptime that has exceeded the given CRITICAL threshold      |
| `vms_with_warning_power_uptime`  |                       |                     | virtual machines with a power uptime that has exceeded the given WARNING threshold       |

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

| Nagios State | Description                                                            |
| ------------ | ---------------------------------------------------------------------- |
| `OK`         | Ideal state, VM power cycle uptime is within bounds.                   |
| `WARNING`    | VM power cycle uptime crossed user-specified threshold for this state. |
| `CRITICAL`   | VM power cycle uptime crossed user-specified threshold for this state. |

### Command-line arguments

- Use the `-h` or `--help` flag to display current usage information.
- Flags marked as **`required`** must be set via CLI flag.
- Flags *not* marked as required are for settings where a useful default is
  already defined, but may be overridden if desired.

| Flag                    | Required | Default | Repeat | Possible                                                                | Description                                                                                                                                                                                                                                                                                                                          |
| ----------------------- | -------- | ------- | ------ | ----------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `branding`              | No       | `false` | No     | `branding`                                                              | Toggles emission of branding details with plugin status details. This output is disabled by default.                                                                                                                                                                                                                                 |
| `h`, `help`             | No       | `false` | No     | `h`, `help`                                                             | Show Help text along with the list of supported flags.                                                                                                                                                                                                                                                                               |
| `v`, `version`          | No       | `false` | No     | `v`, `version`                                                          | Whether to display application version and then immediately exit application.                                                                                                                                                                                                                                                        |
| `ll`, `log-level`       | No       | `info`  | No     | `disabled`, `panic`, `fatal`, `error`, `warn`, `info`, `debug`, `trace` | Log message priority filter. Log messages with a lower level are ignored. Log messages are sent to `stderr` by default. See [Output](#output) for more information.                                                                                                                                                                  |
| `p`, `port`             | No       | `443`   | No     | *positive whole number between 1-65535, inclusive*                      | TCP port of the remote ESXi host or vCenter instance. This is usually 443 (HTTPS).                                                                                                                                                                                                                                                   |
| `t`, `timeout`          | No       | `10`    | No     | *positive whole number of seconds*                                      | Timeout value in seconds allowed before a plugin execution attempt is abandoned and an error returned.                                                                                                                                                                                                                               |
| `s`, `server`           | **Yes**  |         | No     | *fully-qualified domain name or IP Address*                             | The fully-qualified domain name or IP Address of the remote ESXi host or vCenter instance.                                                                                                                                                                                                                                           |
| `u`, `username`         | **Yes**  |         | No     | *valid username*                                                        | Username with permission to access specified ESXi host or vCenter instance.                                                                                                                                                                                                                                                          |
| `pw`, `password`        | **Yes**  |         | No     | *valid password*                                                        | Password used to login to ESXi host or vCenter instance.                                                                                                                                                                                                                                                                             |
| `domain`                | No       |         | No     | *valid user domain*                                                     | (Optional) domain for user account used to login to ESXi host or vCenter instance. This is needed for user accounts residing in a non-default domain (e.g., SSO specific domain).                                                                                                                                                    |
| `trust-cert`            | No       | `false` | No     | `true`, `false`                                                         | Whether the certificate should be trusted as-is without validation. WARNING: TLS is susceptible to man-in-the-middle attacks if enabling this option.                                                                                                                                                                                |
| `include-rp`            | No       |         | No     | *comma-separated list of resource pool names*                           | Specifies a comma-separated list of Resource Pool names that should be exclusively used when evaluating VMs. Specifying this option will also exclude any VMs from evaluation that are *outside* of a Resource Pool. This option is incompatible with specifying a list of Resource Pool names to ignore or exclude from evaluation. |
| `exclude-rp`            | No       |         | No     | *comma-separated list of resource pool names*                           | Specifies a comma-separated list of Resource Pool names that should be ignored when evaluating VMs. This option is incompatible with specifying a list of Resource Pool names to include for evaluation.                                                                                                                             |
| `include-folder-id`     | No       |         | No     | *comma-separated list of folder ID values*                              | Specifies a comma-separated list of Folder Managed Object ID (MOID) values (e.g., group-v34) that should be exclusively used when evaluating VMs. This option is incompatible with specifying a list of Folder IDs to ignore or exclude from evaluation.                                                                             |
| `exclude-folder-id`     | No       |         | No     | *comma-separated list of folder ID values*                              | Specifies a comma-separated list of Folder Managed Object ID (MOID) values (e.g., group-v34) that should be ignored when evaluating VMs. This option is incompatible with specifying a list of Folder Managed Object ID (MOID) values to include for evaluation.                                                                     |
| `ignore-vm`             | No       |         | No     | *comma-separated list of (vSphere) virtual machine names*               | Specifies a comma-separated list of VM names that should be ignored or excluded from evaluation.                                                                                                                                                                                                                                     |
| `uc`, `uptime-critical` | No       | `90`    | No     | *days as positive whole number*                                         | Specifies the power cycle (off/on) uptime in days per VM when a CRITICAL threshold is reached.                                                                                                                                                                                                                                       |
| `uw`, `uptime-warning`  | No       | `60`    | No     | *days as positive whole number*                                         | Specifies the power cycle (off/on) uptime in days per VM when a WARNING threshold is reached.                                                                                                                                                                                                                                        |

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

### CLI invocation

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
  - logging output is intended to be seen when invoking the plugin directly
    via CLI (often for troubleshooting)
    - see the [Output section](../../README.md#output) of the main README for
      potential conflicts with some monitoring systems

### Command definition

```shell
# /etc/nagios-plugins/config/vmware-vm-power-uptime.cfg

# Look at all pools, all VMs, do not evaluate any VMs that are powered off.
# This variation of the command is most useful for environments where all VMs
# are monitored equally.
define command{
    command_name    check_vmware_vm_power_uptime
    command_line    $USER1$/check_vmware_vm_power_uptime --server '$HOSTNAME$' --domain '$ARG1$' --username '$ARG2$' --password '$ARG3$' --uptime-warning '$ARG4$' --uptime-critical '$ARG5$' --trust-cert  --log-level info
    }
```

## License

See the [main project README](../../README.md) for details.

## References

- [Main project README](../../README.md)
- [Documentation index](../README.md)
- [Project repo][repo-url]

<!-- Footnotes here  -->

[repo-url]: <https://github.com/atc0005/check-vmware>  "This project's GitHub repo"

<!-- []: PLACEHOLDER "DESCRIPTION_HERE" -->
