<!-- omit in toc -->
# [check-vmware][repo-url] | `check_vmware_host_memory` plugin

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

Nagios plugin used to monitor ESXi host memory.

In addition to reporting current host memory usage, this plugin also reports
which VMs are on the host (running or not), how much memory each VM is using
as a fixed value and as a percentage of the host's total memory.

Thresholds for `CRITICAL` and `WARNING` memory usage have usable defaults, but
max memory usage is required before this plugin can be used. See the
[configuration options](#configuration-options) section for details.

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

**NOTE**: These metrics are based on the visibility of the service account
used to login to the target VMware environment. If the service account cannot
see a resource, it cannot evaluate the resource.

| Metric             | Unit of Measurement | Description                                   |
| ------------------ | ------------------- | --------------------------------------------- |
| `time`             | milliseconds        | plugin runtime                                |
| `vms`              |                     | all (visible) virtual machines on the host    |
| `vms_powered_on`   |                     | virtual machines powered on                   |
| `vms_powered_off`  |                     | virtual machines powered off                  |
| `memory_usage`     | percentage          | cpu usage                                     |
| `memory_total`     | Hz                  | the total amount of CPU capacity for the host |
| `memory_used`      | Hz                  | the consumed host memory                      |
| `memory_remaining` | Hz                  | the remaining host memory                     |

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

| Nagios State | Description                                                                    |
| ------------ | ------------------------------------------------------------------------------ |
| `OK`         | Ideal state, memory usage for the specified ESXi host system is within bounds. |
| `WARNING`    | Memory usage crossed user-specified threshold for this state.                  |
| `CRITICAL`   | Memory usage crossed user-specified threshold for this state.                  |

### Command-line arguments

- Use the `-h` or `--help` flag to display current usage information.
- Flags marked as **`required`** must be set via CLI flag.
- Flags *not* marked as required are for settings where a useful default is
  already defined, but may be overridden if desired.

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
| `domain`                      | No       |         | No     | *valid user domain*                                                     | (Optional) domain for user account used to login to ESXi host or vCenter instance. This is needed for user accounts residing in a non-default domain (e.g., SSO specific domain).                      |
| `trust-cert`                  | No       | `false` | No     | `true`, `false`                                                         | Whether the certificate should be trusted as-is without validation. WARNING: TLS is susceptible to man-in-the-middle attacks if enabling this option.                                                  |
| `dc-name`                     | No       |         | No     | *valid vSphere datacenter name*                                         | Specifies the name of a vSphere Datacenter. If not specified, applicable plugins will attempt to use the default datacenter found in the vSphere environment. Not applicable to standalone ESXi hosts. |
| `host-name`                   | **Yes**  |         | No     | *valid ESXi host name*                                                  | ESXi host/server name as it is found within the vSphere inventory.                                                                                                                                     |
| `mc`, `memory-usage-critical` | No       | `95`    | No     | *percentage as positive whole number*                                   | Specifies the percentage of memory use (as a whole number) when a CRITICAL threshold is reached.                                                                                                       |
| `mw`, `memory-usage-warning`  | No       | `80`    | No     | *percentage as positive whole number*                                   | Specifies the percentage of memory use (as a whole number) when a WARNING threshold is reached.                                                                                                        |

### Configuration file

Not currently supported. This feature may be added later if there is
sufficient interest.

## Contrib

See the [main project README](../../README.md) for details.

## Examples

### CLI invocation

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
  - logging output is intended to be seen when invoking the plugin directly
    via CLI (often for troubleshooting)
    - see the [Output section](../../README.md#output) of the main README for
      potential conflicts with some monitoring systems

### Command definition

```shell
# /etc/nagios-plugins/config/vmware-host-memory.cfg

# Look at a specific host and explicitly provide custom WARNING and CRITICAL
# threshold values.
define command{
    command_name    check_vmware_host_memory
    command_line    $USER1$/check_vmware_host_memory --server '$HOSTNAME$' --domain '$ARG1$' --username '$ARG2$' --password '$ARG3$' --memory-usage-warning '$ARG4$' --memory-usage-critical '$ARG5$' --host-name '$ARG6$' --trust-cert  --log-level info
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
