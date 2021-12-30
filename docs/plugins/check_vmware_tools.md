<!-- omit in toc -->
# [check-vmware][repo-url] | `check_vmware_tools` plugin

- [Main project README](../../README.md)
- [Documentation index](../README.md)

<!-- omit in toc -->
## Table of Contents

- [Overview](#overview)
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

Nagios plugin used to monitor VMware Tools installations. See the
[configuration options](#configuration-options) section for details regarding
how the various Tools states are evaluated.

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
- `vms_excluded_by_power_state`
- `vms_with_tools_issues`
- `vms_without_tools_issues`
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

### Command-line arguments

- Use the `-h` or `--help` flag to display current usage information.
- Flags marked as **`required`** must be set via CLI flag.
- Flags *not* marked as required are for settings where a useful default is
  already defined, but may be overridden if desired.

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
| `domain`          | No       |         | No     | *valid user domain*                                                     | (Optional) domain for user account used to login to ESXi host or vCenter instance. This is needed for user accounts residing in a non-default domain (e.g., SSO specific domain).                                                                                                                                          |
| `trust-cert`      | No       | `false` | No     | `true`, `false`                                                         | Whether the certificate should be trusted as-is without validation. WARNING: TLS is susceptible to man-in-the-middle attacks if enabling this option.                                                                                                                                                                      |
| `include-rp`      | No       |         | No     | *comma-separated list of resource pool names*                           | Specifies a comma-separated list of Resource Pools that should be exclusively used when evaluating VMs. Specifying this option will also exclude any VMs from evaluation that are *outside* of a Resource Pool. This option is incompatible with specifying a list of Resource Pools to ignore or exclude from evaluation. |
| `exclude-rp`      | No       |         | No     | *comma-separated list of resource pool names*                           | Specifies a comma-separated list of Resource Pools that should be ignored when evaluating VMs. This option is incompatible with specifying a list of Resource Pools to include for evaluation.                                                                                                                             |
| `ignore-vm`       | No       |         | No     | *comma-separated list of (vSphere) virtual machine names*               | Specifies a comma-separated list of VM names that should be ignored or excluded from evaluation.                                                                                                                                                                                                                           |
| `powered-off`     | No       | `false` | No     | `true`, `false`                                                         | Toggles evaluation of powered off VMs in addition to powered on VMs. Evaluation of powered off VMs is disabled by default.                                                                                                                                                                                                 |

### Configuration file

Not currently supported. This feature may be added later if there is
sufficient interest.

## Contrib

See the [main project README](../../README.md) for details.

## Examples

### CLI invocation

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
  - logging output is intended to be seen when invoking the plugin directly
    via CLI (often for troubleshooting)
    - see the [Output section](../../README.md#output) of the main README for
      potential conflicts with some monitoring systems

### Command definition

```shell
# /etc/nagios-plugins/config/vmware-tools.cfg

# Look at all pools, all VMs, do not evaluate any VMs that are powered off.
# This variation of the command is most useful for environments where all VMs
# are monitored equally.
define command{
    command_name    check_vmware_tools
    command_line    $USER1$/check_vmware_tools --server '$HOSTNAME$' --domain '$ARG1$' --username '$ARG2$' --password '$ARG3$'  --trust-cert  --log-level info
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

[vsphere-guestinfo-data-object]: <https://vdc-download.vmware.com/vmwb-repository/dcr-public/b50dcbbf-051d-4204-a3e7-e1b618c1e384/538cf2ec-b34f-4bae-a332-3820ef9e7773/vim.vm.GuestInfo.html>

<!-- []: PLACEHOLDER "DESCRIPTION_HERE" -->
