<!-- omit in toc -->
# [check-vmware][repo-url] | `check_vmware_snapshots_age` plugin

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
- `vms_with_critical_snapshots`
- `vms_with_warning_snapshots`
- `snapshots`
- `critical_snapshots`
- `warning_snapshots`
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

| Flag                 | Required | Default | Repeat | Possible                                                                | Description                                                                                                                                                                                                                                                                                                                          |
| -------------------- | -------- | ------- | ------ | ----------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `branding`           | No       | `false` | No     | `branding`                                                              | Toggles emission of branding details with plugin status details. This output is disabled by default.                                                                                                                                                                                                                                 |
| `h`, `help`          | No       | `false` | No     | `h`, `help`                                                             | Show Help text along with the list of supported flags.                                                                                                                                                                                                                                                                               |
| `v`, `version`       | No       | `false` | No     | `v`, `version`                                                          | Whether to display application version and then immediately exit application.                                                                                                                                                                                                                                                        |
| `ll`, `log-level`    | No       | `info`  | No     | `disabled`, `panic`, `fatal`, `error`, `warn`, `info`, `debug`, `trace` | Log message priority filter. Log messages with a lower level are ignored. Log messages are sent to `stderr` by default. See [Output](#output) for more information.                                                                                                                                                                  |
| `p`, `port`          | No       | `443`   | No     | *positive whole number between 1-65535, inclusive*                      | TCP port of the remote ESXi host or vCenter instance. This is usually 443 (HTTPS).                                                                                                                                                                                                                                                   |
| `t`, `timeout`       | No       | `10`    | No     | *positive whole number of seconds*                                      | Timeout value in seconds allowed before a plugin execution attempt is abandoned and an error returned.                                                                                                                                                                                                                               |
| `s`, `server`        | **Yes**  |         | No     | *fully-qualified domain name or IP Address*                             | The fully-qualified domain name or IP Address of the remote ESXi host or vCenter instance.                                                                                                                                                                                                                                           |
| `u`, `username`      | **Yes**  |         | No     | *valid username*                                                        | Username with permission to access specified ESXi host or vCenter instance.                                                                                                                                                                                                                                                          |
| `pw`, `password`     | **Yes**  |         | No     | *valid password*                                                        | Password used to login to ESXi host or vCenter instance.                                                                                                                                                                                                                                                                             |
| `domain`             | No       |         | No     | *valid user domain*                                                     | (Optional) domain for user account used to login to ESXi host or vCenter instance. This is needed for user accounts residing in a non-default domain (e.g., SSO specific domain).                                                                                                                                                    |
| `trust-cert`         | No       | `false` | No     | `true`, `false`                                                         | Whether the certificate should be trusted as-is without validation. WARNING: TLS is susceptible to man-in-the-middle attacks if enabling this option.                                                                                                                                                                                |
| `include-rp`         | No       |         | No     | *comma-separated list of resource pool names*                           | Specifies a comma-separated list of Resource Pool names that should be exclusively used when evaluating VMs. Specifying this option will also exclude any VMs from evaluation that are *outside* of a Resource Pool. This option is incompatible with specifying a list of Resource Pool names to ignore or exclude from evaluation. |
| `exclude-rp`         | No       |         | No     | *comma-separated list of resource pool names*                           | Specifies a comma-separated list of Resource Pool names that should be ignored when evaluating VMs. This option is incompatible with specifying a list of Resource Pool names to include for evaluation.                                                                                                                             |
| `ignore-vm`          | No       |         | No     | *comma-separated list of (vSphere) virtual machine names*               | Specifies a comma-separated list of VM names that should be ignored or excluded from evaluation.                                                                                                                                                                                                                                     |
| `ac`, `age-critical` | No       | `2`     | No     | *age in days as positive whole number*                                  | Specifies the age of a snapshot in days when a CRITICAL threshold is reached.                                                                                                                                                                                                                                                        |
| `aw`, `age-warning`  | No       | `1`     | No     | *age in days as positive whole number*                                  | Specifies the age of a snapshot in days when a WARNING threshold is reached.                                                                                                                                                                                                                                                         |

### Configuration file

Not currently supported. This feature may be added later if there is
sufficient interest.

## Contrib

See the [main project README](../../README.md) for details.

## Examples

### CLI invocation

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
  - logging output is intended to be seen when invoking the plugin directly
    via CLI (often for troubleshooting)
    - see the [Output section](../../README.md#output) of the main README for
      potential conflicts with some monitoring systems

### Command definition

```shell
# /etc/nagios-plugins/config/vmware-snapshots-age.cfg

# Look at all pools, all VMs, do not evaluate any VMs that are powered off.
# This variation of the command is most useful for environments where all VMs
# are monitored equally.
define command{
    command_name    check_vmware_snapshots_age
    command_line    $USER1$/check_vmware_snapshots_age --server '$HOSTNAME$' --domain '$ARG1$' --username '$ARG2$' --password '$ARG3$' --age-warning '$ARG4$' --age-critical '$ARG5$' --trust-cert --log-level info
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
