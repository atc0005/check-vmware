<!-- omit in toc -->
# [check-vmware][repo-url] | `check_vmware_disk_consolidation` plugin

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

- `time`
- `vms`
- `vms_excluded_by_name`
- `vms_with_consolidation_need`
- `vms_without_consolidation_need`
- `resource_pools_excluded`
- `resource_pools_included`
- `resource_pools_evaluated`

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

| Nagios State | Description                                    |
| ------------ | ---------------------------------------------- |
| `OK`         | Ideal state, VM disk consolidation not needed. |
| `WARNING`    | Not used by this plugin.                       |
| `CRITICAL`   | Disk consolidation needed for one or more VMs. |

### Command-line arguments

- Use the `-h` or `--help` flag to display current usage information.
- Flags marked as **`required`** must be set via CLI flag.
- Flags *not* marked as required are for settings where a useful default is
  already defined, but may be overridden if desired.

| Flag              | Required | Default | Repeat | Possible                                                                | Description                                                                                                                                                                                                                                                                                                                          |
| ----------------- | -------- | ------- | ------ | ----------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `branding`        | No       | `false` | No     | `branding`                                                              | Toggles emission of branding details with plugin status details. This output is disabled by default.                                                                                                                                                                                                                                 |
| `h`, `help`       | No       | `false` | No     | `h`, `help`                                                             | Show Help text along with the list of supported flags.                                                                                                                                                                                                                                                                               |
| `v`, `version`    | No       | `false` | No     | `v`, `version`                                                          | Whether to display application version and then immediately exit application.                                                                                                                                                                                                                                                        |
| `ll`, `log-level` | No       | `info`  | No     | `disabled`, `panic`, `fatal`, `error`, `warn`, `info`, `debug`, `trace` | Log message priority filter. Log messages with a lower level are ignored. Log messages are sent to `stderr` by default. See [Output](#output) for more information.                                                                                                                                                                  |
| `p`, `port`       | No       | `443`   | No     | *positive whole number between 1-65535, inclusive*                      | TCP port of the remote ESXi host or vCenter instance. This is usually 443 (HTTPS).                                                                                                                                                                                                                                                   |
| `t`, `timeout`    | No       | `10`    | No     | *positive whole number of seconds*                                      | Timeout value in seconds allowed before a plugin execution attempt is abandoned and an error returned.                                                                                                                                                                                                                               |
| `s`, `server`     | **Yes**  |         | No     | *fully-qualified domain name or IP Address*                             | The fully-qualified domain name or IP Address of the remote ESXi host or vCenter instance.                                                                                                                                                                                                                                           |
| `u`, `username`   | **Yes**  |         | No     | *valid username*                                                        | Username with permission to access specified ESXi host or vCenter instance.                                                                                                                                                                                                                                                          |
| `pw`, `password`  | **Yes**  |         | No     | *valid password*                                                        | Password used to login to ESXi host or vCenter instance.                                                                                                                                                                                                                                                                             |
| `domain`          | No       |         | No     | *valid user domain*                                                     | (Optional) domain for user account used to login to ESXi host or vCenter instance. This is needed for user accounts residing in a non-default domain (e.g., SSO specific domain).                                                                                                                                                    |
| `trust-cert`      | No       | `false` | No     | `true`, `false`                                                         | Whether the certificate should be trusted as-is without validation. WARNING: TLS is susceptible to man-in-the-middle attacks if enabling this option.                                                                                                                                                                                |
| `include-rp`      | No       |         | No     | *comma-separated list of resource pool names*                           | Specifies a comma-separated list of Resource Pool names that should be exclusively used when evaluating VMs. Specifying this option will also exclude any VMs from evaluation that are *outside* of a Resource Pool. This option is incompatible with specifying a list of Resource Pool names to ignore or exclude from evaluation. |
| `exclude-rp`      | No       |         | No     | *comma-separated list of resource pool names*                           | Specifies a comma-separated list of Resource Pool names that should be ignored when evaluating VMs. This option is incompatible with specifying a list of Resource Pool names to include for evaluation.                                                                                                                             |
| `ignore-vm`       | No       |         | No     | *comma-separated list of (vSphere) virtual machine names*               | Specifies a comma-separated list of VM names that should be ignored or excluded from evaluation.                                                                                                                                                                                                                                     |
| `trigger-reload`  | No       | `false` | No     | `true`, `false`                                                         | Trigger a reload operation for each VM evaluated. This option ensures that the most current state data is evaluated, but increases plugin runtime. If using this, you should also adjust the `--timeout` value and potentially your monitor system's service check timeout setting.                                                  |

### Configuration file

Not currently supported. This feature may be added later if there is
sufficient interest.

## Contrib

See the [main project README](../../README.md) for details.

## Examples

### CLI invocation

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
  - logging output is intended to be seen when invoking the plugin directly
    via CLI (often for troubleshooting)
    - see the [Output section](../../README.md#output) of the main README for
      potential conflicts with some monitoring systems
- A forced state data reload/refresh is triggered for each evaluated Virtual
  Machine
  - NOTE: this operation is expensive, omit to rely on existing (potentially
    stale) state data
  - see the [`overview`](#overview) section for additional details, including
    potential alternatives to the use of this flag
- A custom timeout value is specified
  - NOTE: in order for this timeout value to be respected, you may need to
    adjust the service check timeout value in your monitoring system (e.g.,
    `service_check_timeout` value in your `nagios.cfg` file)

### Command definition

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
    command_line    $USER1$/check_vmware_disk_consolidation --server '$HOSTNAME$' --domain '$ARG1$' --username '$ARG2$' --password '$ARG3$'  --trust-cert --log-level info
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
    command_line    $USER1$/check_vmware_disk_consolidation --server '$HOSTNAME$' --domain '$ARG1$' --username '$ARG2$' --password '$ARG3$'  --trust-cert --log-level info --trigger-reload --timeout 110
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
