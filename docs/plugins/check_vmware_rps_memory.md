<!-- omit in toc -->
# [check-vmware][repo-url] | `check_vmware_rps_memory` plugin

- [Main project README](../../README.md)
- [Documentation index](../README.md)

<!-- omit in toc -->
## Table of Contents

- [Overview](#overview)
- [Output](#output)
- [Limitations](#limitations)
- [`Resources` Resource Pool](#resources-resource-pool)
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

Nagios plugin used to monitor memory usage across Resource Pools.

If specific Resource Pools are not specified by the sysadmin for inclusion or
exclusion all visible Resource Pools will be evaluated.

In addition to reporting memory usage for each Resource Pool, this plugin also
reports the ten most recently booted VMs along with their memory usage. This
is intended to help spot which VM is responsible for a state change alert.

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

## Limitations

**NOTE**: This plugin is not compatible with standalone ESXi servers. It must
be used with a vCenter server instance.

See [GH-643](https://github.com/atc0005/check-vmware/discussions/643) and
[GH-657](https://github.com/atc0005/check-vmware/issues/657) for prior
discussion/troubleshooting. Please file a new GH issue if you have any
additional information that would assist with resolving the issue.

## `Resources` Resource Pool

**NOTE**: There is a parent or root Resource Pool named `Resources`. This
plugin is hard-coded to exclude this Resource Pool from evaluation. Since
other Resource Pools are descended from the `Resources` Resource Pool,
evaluating this resource pool directly would throw off calculations. The
`Resources` Resource Pool is listed in the plugin output as excluded,
regardless of whether the sysadmin opts to exclude any Resource Pools.

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
- `memory_usage`
- `memory_used`
- `memory_remaining`
- `memory_ballooned`
- `memory_swapped`
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

| Nagios State | Description                                                     |
| ------------ | --------------------------------------------------------------- |
| `OK`         | Ideal state, memory usage across Resources Pools within bounds. |
| `WARNING`    | Memory usage crossed user-specified threshold for this state.   |
| `CRITICAL`   | Memory usage crossed user-specified threshold for this state.   |

### Command-line arguments

- Use the `-h` or `--help` flag to display current usage information.
- Flags marked as **`required`** must be set via CLI flag.
- Flags *not* marked as required are for settings where a useful default is
  already defined, but may be overridden if desired.

| Flag                        | Required | Default | Repeat | Possible                                                                | Description                                                                                                                                                                                                                                                                                                                |
| --------------------------- | -------- | ------- | ------ | ----------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `branding`                  | No       | `false` | No     | `branding`                                                              | Toggles emission of branding details with plugin status details. This output is disabled by default.                                                                                                                                                                                                                       |
| `h`, `help`                 | No       | `false` | No     | `h`, `help`                                                             | Show Help text along with the list of supported flags.                                                                                                                                                                                                                                                                     |
| `v`, `version`              | No       | `false` | No     | `v`, `version`                                                          | Whether to display application version and then immediately exit application.                                                                                                                                                                                                                                              |
| `ll`, `log-level`           | No       | `info`  | No     | `disabled`, `panic`, `fatal`, `error`, `warn`, `info`, `debug`, `trace` | Log message priority filter. Log messages with a lower level are ignored. Log messages are sent to `stderr` by default. See [Output](#output) for more information.                                                                                                                                                        |
| `p`, `port`                 | No       | `443`   | No     | *positive whole number between 1-65535, inclusive*                      | TCP port of the remote vCenter instance. This is usually 443 (HTTPS).                                                                                                                                                                                                                                                      |
| `t`, `timeout`              | No       | `10`    | No     | *positive whole number of seconds*                                      | Timeout value in seconds allowed before a plugin execution attempt is abandoned and an error returned.                                                                                                                                                                                                                     |
| `s`, `server`               | **Yes**  |         | No     | *fully-qualified domain name or IP Address*                             | The fully-qualified domain name or IP Address of the remote vCenter instance.                                                                                                                                                                                                                                              |
| `u`, `username`             | **Yes**  |         | No     | *valid username*                                                        | Username with permission to access the specified vCenter instance.                                                                                                                                                                                                                                                         |
| `pw`, `password`            | **Yes**  |         | No     | *valid password*                                                        | Password used to login to the vCenter instance.                                                                                                                                                                                                                                                                            |
| `domain`                    | No       |         | No     | *valid user domain*                                                     | (Optional) domain for the user account used to login to the vCenter instance. This is needed for user accounts residing in a non-default domain (e.g., SSO specific domain).                                                                                                                                               |
| `trust-cert`                | No       | `false` | No     | `true`, `false`                                                         | Whether the certificate should be trusted as-is without validation. WARNING: TLS is susceptible to man-in-the-middle attacks if enabling this option.                                                                                                                                                                      |
| `include-rp`                | No       |         | No     | *comma-separated list of resource pool names*                           | Specifies a comma-separated list of Resource Pools that should be exclusively used when evaluating VMs. Specifying this option will also exclude any VMs from evaluation that are *outside* of a Resource Pool. This option is incompatible with specifying a list of Resource Pools to ignore or exclude from evaluation. |
| `exclude-rp`                | No       |         | No     | *comma-separated list of resource pool names*                           | Specifies a comma-separated list of Resource Pools that should be ignored when evaluating VMs. This option is incompatible with specifying a list of Resource Pools to include for evaluation.                                                                                                                             |
| `mma`, `memory-max-allowed` | **Yes**  | `0`     | No     | *positive whole number in GB*                                           | Specifies the maximum amount of memory that we are allowed to consume in GB (as a whole number) in the target VMware environment across all specified Resource Pools. VMs that are running outside of resource pools are not considered in these calculations.                                                             |
| `mc`, `memory-use-critical` | No       | `95`    | No     | *percentage as positive whole number*                                   | Specifies the percentage of memory use (as a whole number) across all specified Resource Pools when a CRITICAL threshold is reached.                                                                                                                                                                                       |
| `mw`, `memory-use-warning`  | No       | `100`   | No     | *percentage as positive whole number*                                   | Specifies the percentage of memory use (as a whole number) across all specified Resource Pools when a WARNING threshold is reached.                                                                                                                                                                                        |

### Configuration file

Not currently supported. This feature may be added later if there is
sufficient interest.

## Contrib

See the [main project README](../../README.md) for details.

## Examples

### CLI invocation

```ShellSession
/usr/lib/nagios/plugins/check_vmware_rps_memory --username SERVICE_ACCOUNT_NAME --password "SERVICE_ACCOUNT_PASSWORD" --server vc1.example.com --exclude-rp "Desktops" --memory-use-warning 80 --memory-use-critical 95  --memory-max-allowed 320 --trust-cert --log-level info
```

See the [configuration options](#configuration-options) section for all
command-line settings supported by this plugin along with descriptions of
each. See the [contrib](#contrib) section for information regarding example
command definitions and Nagios configuration files.

Of note:

- The default/parent/root resource pool named `Resources` is excluded from
  evaluation
  - this behavior is hard-coded into the plugin
  - since other resource pools are descended from this one, evaluating this
    resource pool directly would skew memory usage calculations
- The resource pool named `Desktops` was specified by the sysadmin to be
  excluded from evaluation
  - this results in *all other* resource pools visible to the specified user
    account being used for evaluation
  - VMs *outside* of a Resource Pool (visible to the specified user account or
    not) do not contribute to memory usage calculations
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

NOTE: This is the inverse of the command-line example for this plugin; only
specified Resource Pools are evaluated.

```shell
# /etc/nagios-plugins/config/vmware-resource-pools.cfg

# This variation of the command does not allow exclusions
define command{
    command_name    check_vmware_resource_pools_include_pools
    command_line    $USER1$/check_vmware_rps_memory --server '$HOSTNAME$' --domain '$ARG1$' --username '$ARG2$' --password '$ARG3$' --memory-use-warning '$ARG4$' --memory-use-critical '$ARG5$' --memory-max-allowed '$ARG6$' --include-rp '$ARG7$' --trust-cert  --log-level info
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
