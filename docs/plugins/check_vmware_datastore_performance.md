<!-- omit in toc -->
# [check-vmware][repo-url] | `check_vmware_datastore_performance` plugin

- [Main project README](../../README.md)
- [Documentation index](../README.md)

<!-- omit in toc -->
## Table of Contents

- [Overview](#overview)
  - [Requirements](#requirements)
  - [How datastore performance metrics are evaluated](#how-datastore-performance-metrics-are-evaluated)
  - [Performance Data metrics](#performance-data-metrics)
  - [Stability of this plugin](#stability-of-this-plugin)
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

Nagios plugin used to monitor datastore performance.

In addition to reporting current [datastore performance
details][vsphere-storage-performance-summary-data-object], this plugin also
reports which VMs reside on the datastore along with their percentage of the
total datastore space used. This is intended to help pinpoint potential causes
of high latency at a glance.

### Requirements

This plugin requires that the `Statistics Collection` setting (part of
`Storage I/O Control`) for a monitored datastore be enabled. If it is not,
this plugin is unable to evaluate performance for a specified datastore. This
plugin attempts to detect and report this condition so that vSphere
administrators can assist with enabling this feature.

To help with locating datastores in need of adjustment, the following PowerCLI
snippet may be used:

```powershell
$credential = Get-Credential -Message "Enter your credentials (DOMAIN\ID)"
$server = Connect-VIServer -Server vc1.example.com -Credential $credential

Get-View -ViewType Datastore |
    Where-Object {$_.IormConfiguration.StatsCollectionEnabled -eq $false} |
    Select -Property Name, @{Label="StatsCollectionEnabled"; Expression={$_.IormConfiguration.StatsCollectionEnabled}} |
    Sort-Object -Property Name

Disconnect-VIServer $server
```

Available settings For `Storage I/O Control`:

- `Disabled`
- `Statistics enabled but Storage I/O disabled`
- `Statistics and Storage I/O enabled`

### How datastore performance metrics are evaluated

Performance metrics are provided by vSphere in aggregated quantiles over a
period of time (intervals). Aggregated metrics correspond with a specific
percentile. As of this plugin's initial development, vSphere provides metrics
associated with these percentiles:

- `90`
- `80`
- `70`
- `60`
- `50`

If not otherwise specified, percentile `90` is used to evaluate datastore
performance metrics. While the vSphere API provides metrics in multiple
intervals (one active, up to seven historical), only the active interval is
used for evaluating current datastore performance.

There is a brief window between when the current interval ends and the new
active interval begins that no metrics are available for the active interval.
Testing shows that this is approximately 30 minutes. The current plugin design
is to omit performance data latency metrics if no metrics are available. This
is done in an attempt to prevent skewing historical data already collected.

This plugin accepts flags to:

- specify individual latency metric thresholds (e.g., read latency CRITICAL,
  read latency WARNING, write latency ...)
- specify percentile *sets*
  - multiple sets supported, each composed of a percentile and pairs of
    CRITICAL and WARNING threshold values

If you specify a percentile set, the plugin will not accept individual latency
threshold flags. The reverse is also true, specifying one or more latency
threshold flags is incompatible with specifying one or more percentile sets.

By specifying multiple percentile sets, you are indicating that crossing the
thresholds of any one set is enough to trigger a state change.

### Performance Data metrics

This plugin emits Nagios performance data metrics for each percentile in the
active interval that is not completely of value `0`. Any percentile with all
`0` metrics are omitted from the performance data metrics collected & emitted
by the plugin.

Please provide feedback by [opening a new
issue](https://github.com/atc0005/check-vmware/issues/new) or commenting on
the original discussion thread [here
(GH-316)](https://github.com/atc0005/check-vmware/discussions/316) if you find
that this decision causes problems with gathering metrics.

### Stability of this plugin

**NOTE**: This plugin uses the [`QueryDatastorePerformanceSummary()` method
provided by the `StorageResourceManager` Managed
Object][vsphere-query-datastore-performance-summary-method]. While available
since vSphere API 5.1, this API is marked as experimental (and subject to
change/removal):

> This is an experimental interface that is not intended for use in production
> code.

In addition to using the experimental `QueryDatastorePerformanceSummary()`
API, this plugin uses the deprecated `statsCollectionEnabled` property from
the [`StorageIORMInfo` Data
Object][vsphere-storage-io-resource-management-data-object] to determine
whether `Statistics Collection` is enabled for a datastore. Using the
prescribed `enabled` property for [that Data
Object][vsphere-storage-io-resource-management-data-object] to determine
`Statistics Collection` does not work.

If you use this plugin, please provide feedback by [opening a new discussion
thread](https://github.com/atc0005/check-vmware/discussions/new) or commenting
on the original discussion thread [here
(GH-316)](https://github.com/atc0005/check-vmware/discussions/316).

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
- `p*_read_latency`
- `p*_write_latency`
- `p*_vm_latency`
- `p*_read_iops`
- `p*_write_iops`
- `vms`
- `vms_powered_off`
- `vms_powered_on`

`*` is a placeholder for `90`, `80`, `70`, `60` & `50` percentiles.

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

**TODO**: Research & note why metric sets might contain all values of `0`.

| Nagios State | Description                                                                                            |
| ------------ | ------------------------------------------------------------------------------------------------------ |
| `OK`         | Ideal state, Datastore performance within bounds for the active interval for the chosen percentile(s). |
| `UNKNOWN`    | Datastore performance metric sets are all value `0` or metrics collection for a datastore is disabled. |
| `WARNING`    | Datastore performance crossed user-specified latency thresholds for this state.                        |
| `CRITICAL`   | Datastore performance crossed user-specified latency thresholds for this state.                        |

### Command-line arguments

- Use the `-h` or `--help` flag to display current usage information.
- Flags marked as **`required`** must be set via CLI flag.
- Flags *not* marked as required are for settings where a useful default is
  already defined, but may be overridden if desired.

| Flag                                       | Required | Default                | Repeat | Possible                                                                | Description                                                                                                                                                                                            |
| ------------------------------------------ | -------- | ---------------------- | ------ | ----------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `branding`                                 | No       | `false`                | No     | `branding`                                                              | Toggles emission of branding details with plugin status details. This output is disabled by default.                                                                                                   |
| `h`, `help`                                | No       | `false`                | No     | `h`, `help`                                                             | Show Help text along with the list of supported flags.                                                                                                                                                 |
| `v`, `version`                             | No       | `false`                | No     | `v`, `version`                                                          | Whether to display application version and then immediately exit application.                                                                                                                          |
| `ll`, `log-level`                          | No       | `info`                 | No     | `disabled`, `panic`, `fatal`, `error`, `warn`, `info`, `debug`, `trace` | Log message priority filter. Log messages with a lower level are ignored. Log messages are sent to `stderr` by default. See [Output](#output) for more information.                                    |
| `p`, `port`                                | No       | `443`                  | No     | *positive whole number between 1-65535, inclusive*                      | TCP port of the remote ESXi host or vCenter instance. This is usually 443 (HTTPS).                                                                                                                     |
| `t`, `timeout`                             | No       | `10`                   | No     | *positive whole number of seconds*                                      | Timeout value in seconds allowed before a plugin execution attempt is abandoned and an error returned.                                                                                                 |
| `s`, `server`                              | **Yes**  |                        | No     | *fully-qualified domain name or IP Address*                             | The fully-qualified domain name or IP Address of the remote ESXi host or vCenter instance.                                                                                                             |
| `u`, `username`                            | **Yes**  |                        | No     | *valid username*                                                        | Username with permission to access specified ESXi host or vCenter instance.                                                                                                                            |
| `pw`, `password`                           | **Yes**  |                        | No     | *valid password*                                                        | Password used to login to ESXi host or vCenter instance.                                                                                                                                               |
| `domain`                                   | No       |                        | No     | *valid user domain*                                                     | (Optional) domain for user account used to login to ESXi host or vCenter instance. This is needed for user accounts residing in a non-default domain (e.g., SSO specific domain).                      |
| `trust-cert`                               | No       | `false`                | No     | `true`, `false`                                                         | Whether the certificate should be trusted as-is without validation. WARNING: TLS is susceptible to man-in-the-middle attacks if enabling this option.                                                  |
| `dc-name`                                  | No       |                        | No     | *valid vSphere datacenter name*                                         | Specifies the name of a vSphere Datacenter. If not specified, applicable plugins will attempt to use the default datacenter found in the vSphere environment. Not applicable to standalone ESXi hosts. |
| `ds-name`                                  | **Yes**  |                        | No     | *valid datastore name*                                                  | Datastore name as it is found within the vSphere inventory.                                                                                                                                            |
| `dsim`, `ds-ignore-missing-metrics`        | No       | `false`                | No     | `true`, `false`                                                         | Toggles how missing Datastore Performance metrics will be handled.This is believed to occur when a datastore is newly created and metrics have not yet been collected.                                 |
| `dshhms`, `ds-hide-historical-metric-sets` | No       | `false`                | No     | `true`, `false`                                                         | Toggles display of historical Datastore Performance metrics at plugin completion. By default historical metrics are listed.                                                                            |
| `dsrlc`, `ds-read-latency-critical`        | No       | `15`                   | No     | *positive whole number or float*                                        | Specifies the read latency of a datastore's storage (in ms) when a `CRITICAL` threshold is reached. The default percentile is used (`90`).                                                             |
| `dsrlw`, `ds-read-latency-warning`         | No       | `30`                   | No     | *positive whole number or float*                                        | Specifies the read latency of a datastore's storage (in ms) when a `WARNING` threshold is reached. The default percentile is used (`90`).                                                              |
| `dswlc`, `ds-write-latency-critical`       | No       | `15`                   | No     | *positive whole number or float*                                        | Specifies the write latency of a datastore's storage (in ms) when a `CRITICAL` threshold is reached. The default percentile is used (`90`).                                                            |
| `dswlw`, `ds-write-latency-warning`        | No       | `30`                   | No     | *positive whole number or float*                                        | Specifies the write latency of a datastore's storage (in ms) when a `WARNING` threshold is reached. The default percentile is used (`90`).                                                             |
| `dsvmlc`, `ds-vm-latency-critical`         | No       | `15`                   | No     | *positive whole number or float*                                        | Specifies the latency (in ms) as observed by VMs using the datastore when a `CRITICAL` threshold is reached. The default percentile is used (`90`).                                                    |
| `dsvmlw`, `ds-vm-latency-warning`          | No       | `30`                   | No     | *positive whole number or float*                                        | Specifies the latency (in ms) as observed by VMs using the datastore when a `WARNING` threshold is reached. The default percentile is used (`90`).                                                     |
| `dslps`, `ds-latency-percentile-set`       | No       | `90,15,30,15,30,15,30` | Yes    | *complete percentile set* in `P,RLW,RLC,WLW,WLC,VMLW,VMLC` format       | Specifies the performance percentile set used for threshold calculations. Incompatible with individual latency threshold flags. All comma-separated field values are required for each set.            |

### Configuration file

Not currently supported. This feature may be added later if there is
sufficient interest.

## Contrib

See the [main project README](../../README.md) for details.

## Examples

### CLI invocation

```ShellSession
/usr/lib/nagios/plugins/check_vmware_datastore_performance --server vc1.example.com --username SERVICE_ACCOUNT_NAME --password "SERVICE_ACCOUNT_PASSWORD" --ds-latency-percentile-set '90,15,30,15,30,15,30' --ds-name "HUSVM-DC1-vol6" --trust-cert  --log-level info
```

See the [configuration options](#configuration-options) section for all
command-line settings supported by this plugin along with descriptions of
each. See the [contrib](#contrib) section for information regarding example
command definitions and Nagios configuration files.

Of note:

- We use a datastore performance percentile set instead of individual latency
  flags
  - `90`th percentile
  - read latency `WARNING` threshold of `15 ms`
  - read latency `CRITICAL` threshold of `30 ms`
  - write latency `WARNING` threshold of `15 ms`
  - write latency `CRITICAL` threshold of `30 ms`
  - vm latency `WARNING` threshold of `15 ms`
  - vm latency `CRITICAL` threshold of `30 ms`
- Due to plugin design, only the active interval is evaluated for threshold
  violations
  - historical interval metrics are reported via `LongServiceOutput` *unless*
    the flag to skip emitting those metrics is specified
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
# /etc/nagios-plugins/config/vmware-datastores-performance.cfg

# Look at specific datastore and explicitly provide custom WARNING and
# CRITICAL latency threshold values via individual flags.
define command{
    command_name    check_vmware_datastore_performance_via_individual_flags
    command_line    $USER1$/check_vmware_datastore_performance --server '$HOSTNAME$' --domain '$ARG1$' --username '$ARG2$' --password '$ARG3$' --ds-read-latency-warning '$ARG4$' --ds-read-latency-critical '$ARG5$' --ds-write-latency-warning '$ARG6$' --ds-write-latency-critical '$ARG7$' --ds-vm-latency-warning '$ARG8$' --ds-vm-latency-critical '$ARG9$' --ds-name '$ARG10$' --trust-cert  --log-level info
    }

# Look at specific datastore and explicitly provide custom WARNING and
# CRITICAL latency threshold values for a single percentile via a percentile
# flag set.
define command{
    command_name    check_vmware_datastore_performance_via_1percentile_set
    command_line    $USER1$/check_vmware_datastore_performance --server '$HOSTNAME$' --domain '$ARG1$' --username '$ARG2$' --password '$ARG3$' --ds-latency-percentile-set '$ARG4$' --ds-name '$ARG5$' --trust-cert  --log-level info
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

[vsphere-query-datastore-performance-summary-method]: <https://vdc-download.vmware.com/vmwb-repository/dcr-public/bf660c0a-f060-46e8-a94d-4b5e6ffc77ad/208bc706-e281-49b6-a0ce-b402ec19ef82/SDK/vsphere-ws/docs/ReferenceGuide/vim.StorageResourceManager.html#queryDatastorePerformanceSummary>

[vsphere-storage-performance-summary-data-object]: <https://vdc-download.vmware.com/vmwb-repository/dcr-public/b50dcbbf-051d-4204-a3e7-e1b618c1e384/538cf2ec-b34f-4bae-a332-3820ef9e7773/vim.StorageResourceManager.StoragePerformanceSummary.html>

[vsphere-storage-io-resource-management-data-object]: <https://vdc-download.vmware.com/vmwb-repository/dcr-public/b50dcbbf-051d-4204-a3e7-e1b618c1e384/538cf2ec-b34f-4bae-a332-3820ef9e7773/vim.StorageResourceManager.IORMConfigInfo.html>

<!-- []: PLACEHOLDER "DESCRIPTION_HERE" -->
