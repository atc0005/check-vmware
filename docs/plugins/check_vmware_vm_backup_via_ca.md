<!-- omit in toc -->
# [check-vmware][repo-url] | `check_vmware_vm_backup_via_ca` plugin

- [Main project README](../../README.md)
- [Documentation index](../README.md)

<!-- omit in toc -->
## Table of Contents

- [Overview](#overview)
- [Output](#output)
- [Performance Data](#performance-data)
- [Optional evaluation](#optional-evaluation)
- [Installation](#installation)
- [Configuration options](#configuration-options)
  - [Threshold calculations](#threshold-calculations)
  - [Command-line arguments](#command-line-arguments)
  - [Configuration file](#configuration-file)
  - [Backup Date format](#backup-date-format)
- [Contrib](#contrib)
- [Examples](#examples)
  - [CLI invocation](#cli-invocation)
    - [One-line](#one-line)
    - [A more readable equivalent](#a-more-readable-equivalent)
    - [Explanation](#explanation)
  - [Command definitions](#command-definitions)
- [License](#license)
- [References](#references)

## Overview

Nagios plugin used to monitor the last backup date for virtual machines.

## Output

The output for these plugins is designed to provide the one-line summary
needed by Nagios for quick identification of a problem while providing longer,
more detailed information for display within the web UI, use in email and
Teams notifications
([atc0005/send2teams](https://github.com/atc0005/send2teams)).

See the [main project README](../../README.md) for details.

## Performance Data

These performance data metrics are currently supported:

- `time`
- `vms`
- `vms_excluded_by_name`
- `vms_evaluated`
- `vms_with_backup_dates`
- `vms_without_backup_dates`
- `resource_pools_excluded`
- `resource_pools_included`
- `resource_pools_evaluated`

## Optional evaluation

Virtual machines can be explicitly *included* by one or more resource pools
and explicitly *excluded* by one or more resource pools and (full) virtual
machine names.

See the [configuration options](#configuration-options), [examples](#examples)
and [contrib](#contrib) sections for more information.

## Installation

See the [main project README](../../README.md) for details.

## Configuration options

### Threshold calculations

| Nagios State | Description                                                                                  |
| ------------ | -------------------------------------------------------------------------------------------- |
| `OK`         | Ideal state, all non-excluded VMs have a backup and it is current.                           |
| `UNKNOWN`    | Not currently used by this plugin.                                                           |
| `WARNING`    | Virtual machine backup date exceeds specified WARNING threshold, but not CRITICAL threshold. |
| `WARNING`    | Virtual machine backup is missing.                                                           |
| `WARNING`    | Backup date does not match default/user-specified format.                                    |
| `CRITICAL`   | Virtual machine backup date exceeds specified CRITICAL threshold.                            |

### Command-line arguments

- Use the `-h` or `--help` flag to display current usage information.
- Flags marked as **`required`** must be set via CLI flag.
- Flags *not* marked as required are for settings where a useful default is
  already defined, but may be overridden if desired.

| Flag                         | Required | Default               | Repeat | Possible                                                                | Description                                                                                                                                                                                                                                                                                                                |
| ---------------------------- | -------- | --------------------- | ------ | ----------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `branding`                   | No       | `false`               | No     | `branding`                                                              | Toggles emission of branding details with plugin status details. This output is disabled by default.                                                                                                                                                                                                                       |
| `h`, `help`                  | No       | `false`               | No     | `h`, `help`                                                             | Show Help text along with the list of supported flags.                                                                                                                                                                                                                                                                     |
| `v`, `version`               | No       | `false`               | No     | `v`, `version`                                                          | Whether to display application version and then immediately exit application.                                                                                                                                                                                                                                              |
| `ll`, `log-level`            | No       | `info`                | No     | `disabled`, `panic`, `fatal`, `error`, `warn`, `info`, `debug`, `trace` | Log message priority filter. Log messages with a lower level are ignored. Log messages are sent to `stderr` by default. See [Output](#output) for more information.                                                                                                                                                        |
| `p`, `port`                  | No       | `443`                 | No     | *positive whole number between 1-65535, inclusive*                      | TCP port of the remote ESXi host or vCenter instance. This is usually 443 (HTTPS).                                                                                                                                                                                                                                         |
| `t`, `timeout`               | No       | `10`                  | No     | *positive whole number of seconds*                                      | Timeout value in seconds allowed before a plugin execution attempt is abandoned and an error returned.                                                                                                                                                                                                                     |
| `s`, `server`                | **Yes**  |                       | No     | *fully-qualified domain name or IP Address*                             | The fully-qualified domain name or IP Address of the remote ESXi host or vCenter instance.                                                                                                                                                                                                                                 |
| `u`, `username`              | **Yes**  |                       | No     | *valid username*                                                        | Username with permission to access specified ESXi host or vCenter instance.                                                                                                                                                                                                                                                |
| `pw`, `password`             | **Yes**  |                       | No     | *valid password*                                                        | Password used to login to ESXi host or vCenter instance.                                                                                                                                                                                                                                                                   |
| `domain`                     | No       |                       | No     | *valid user domain*                                                     | (Optional) domain for user account used to login to ESXi host or vCenter instance. This is needed for user accounts residing in a non-default domain (e.g., SSO specific domain).                                                                                                                                          |
| `trust-cert`                 | No       | `false`               | No     | `true`, `false`                                                         | Whether the certificate should be trusted as-is without validation. WARNING: TLS is susceptible to man-in-the-middle attacks if enabling this option.                                                                                                                                                                      |
| `include-rp`                 | No       |                       | No     | *comma-separated list of resource pool names*                           | Specifies a comma-separated list of Resource Pools that should be exclusively used when evaluating VMs. Specifying this option will also exclude any VMs from evaluation that are *outside* of a Resource Pool. This option is incompatible with specifying a list of Resource Pools to ignore or exclude from evaluation. |
| `exclude-rp`                 | No       |                       | No     | *comma-separated list of resource pool names*                           | Specifies a comma-separated list of Resource Pools that should be ignored when evaluating VMs. This option is incompatible with specifying a list of Resource Pools to include for evaluation.                                                                                                                             |
| `ignore-vm`                  | No       |                       | No     | *comma-separated list of (vSphere) virtual machine names*               | Specifies a comma-separated list of VM names that should be ignored or excluded from evaluation.                                                                                                                                                                                                                           |
| `backup-date-ca`             | No       | `Last Backup`         | No     | *valid custom attribute name*                                           | Specifies the name of the custom attribute used by virtual machine backup software to record when the last backup occurred.                                                                                                                                                                                                |
| `backup-metadata-ca`         | No       |                       | No     | *valid custom attribute name*                                           | Specifies the (optional) name of the custom attribute used by virtual machine backup software to record metadata / details for the last backup. If provided, this value is used in log messages and the final report.                                                                                                      |
| `backup-date-format`         | No       | `01/02/2006 15:04:05` | No     | *[supported layout string][official-time-pkg-docs]*                     | Specifies the format of the date recorded when the last backup occurred. See the [official docs][official-time-pkg-docs], [references](#references) and the [examples](#examples) section for more information.                                                                                                            |
| `backup-date-timezone`       | No       | `Local`               | No     | *[valid time zone database entry][tz-database]*                         | Specifies the time zone for the specified custom attribute used by virtual machine backup software to record when the last backup occurred. Requires tz database format (e.g., `Europe/Amsterdam`, `America/New_York`, `Europe/Paris`). See also [tz-database] for examples.                                               |
| `bac`, `backup-age-critical` | No       | `2`                   | No     | *positive whole number of days*                                         | Specifies the number of days since the last backup for a VM when a `CRITICAL` threshold is reached.                                                                                                                                                                                                                        |
| `baw`, `backup-age-warning`  | No       | `1`                   | No     | *positive whole number of days*                                         | Specifies the number of days since the last backup for a VM when a `WARNING` threshold is reached.                                                                                                                                                                                                                         |

### Configuration file

Not currently supported. This feature may be added later if there is
sufficient interest.

### Backup Date format

Instead of using the classic [`strftime` format codes][strftime-codes] from
the C programming language, the Go `time` package uses human readable date
"layout" strings to parse input strings as valid dates and times.

From the [official documentation][official-time-pkg-docs]:

> The reference time used in these layouts is the specific time stamp:
>
> `01/02 03:04:05PM '06 -0700`
>
> (`January 2, 15:04:05, 2006`, in time zone seven hours west of `GMT`). That
> value is recorded as the constant named `Layout`, listed below. As a Unix
> time, this is `1136239445`. Since `MST` is `GMT-0700`, the reference would be
> printed by the Unix `date` command as:
>
> `Mon Jan 2 15:04:05 MST 2006`
>
> It is a regrettable historic error that the date uses the American
> convention of putting the numerical month before the day.

The following table is intended to provide a quick reference for common date
formats and equivalent format strings for use with the `--backup-date-format`
flag. If not specified, the default value is used.

If backup dates for your Virtual Machines are recorded in a format on the
left, use the string across from it in the right column as an argument for the
`--backup-date-format` flag.

| Date format             | Format string           |
| ----------------------- | ----------------------- |
| `01/17/2022 20:14:12`   | `01/02/2006 15:04:05`   |
| `2021-11-09 9:07:21 PM` | `2006-01-02 3:04:05 PM` |

## Contrib

See the [main project README](../../README.md) for details.

## Examples

### CLI invocation

#### One-line

```shell
/usr/lib/nagios/plugins/check_vmware_vm_backup_via_ca --username "SERVICE_ACCOUNT_NAME" --password "SERVICE_ACCOUNT_PASSWORD" --server "vc1.example.com" --exclude-rp "Desktops" --ignore-vm "test1.example.com,redmine.example.com,TESTING-AC,RHEL7-TEST" --trust-cert --log-level info --backup-date-timezone "Europe/Amsterdam" --backup-date-ca "Last Backup" --backup-metadata-ca "Backup Status" --backup-date-format "01/02/2006 15:04:05"
```

#### A more readable equivalent

```shell
/usr/lib/nagios/plugins/check_vmware_vm_backup_via_ca \
    --username "SERVICE_ACCOUNT_NAME" \
    --password "SERVICE_ACCOUNT_PASSWORD" \
    --server "vc1.example.com" \
    --exclude-rp "Desktops" \
    --ignore-vm "test1.example.com,redmine.example.com,TESTING-AC,RHEL7-TEST" \
    --port "443" \
    --trust-cert \
    --log-level info \
    --backup-date-timezone "Europe/Amsterdam" \
    --backup-date-ca "Last Backup" \
    --backup-metadata-ca "Backup Status" \
    --backup-date-format "01/02/2006 15:04:05"
```

The examples above attempt to showcase the majority of the supported flags,
but not all are required (where default values are sufficient for your
environment).

See the [configuration options](#configuration-options) section for all
command-line settings supported by this plugin along with descriptions of
each. See the [contrib](#contrib) section for information regarding example
command definitions and Nagios configuration files.

#### Explanation

- We specify required connections settings
  - username
  - password
  - server
- We specify settings for including/excluding VMs from evaluation
  - exclude all VMs from the `Desktop` resource pool
  - explicitly ignore (full) names of VMs
    - `test1.example.com`
    - `redmine.example.com`
    - `TESTING-AC`
    - `RHEL7-TEST`
- We specify settings specific to backups
  - the time zone instead of defaulting to the local time zone
    - this affects parsing of the date/time recorded for a virtual machine's
      last backup date
  - the custom attribute name used by backup software to record when the last
    (presumably successful) backup occurred.
  - the (optional) custom attribute name used by backup software to record
    metadata for the last (presumably successful) backup.
  - the format or date "layout" used by the backup software to record when the
    last (presumably successful) backup occurred.
    - see the [official documentation][official-time-pkg-docs] and other
      third-party resources noted in the [References](#references) section and
      the table of [common backup date formats](#backup-date-format) for
      additional information
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

### Command definitions

```shell
# /etc/nagios-plugins/config/vmware-vm-backup-via-ca.cfg

# Look at all resource pools, all virtual machines. Use default values for
# time zone, backup date format, custom attribute name for last backup and
# thresholds. This variation of the command is most useful for environments
# where all VMs are monitored equally and where the default plugin values are
# sufficient.
define command{
    command_name    check_vmware_vm_backup_via_ca
    command_line    $USER1$/check_vmware_vm_backup_via_ca --server '$HOSTNAME$' --domain '$ARG1$' --username '$ARG2$' --password '$ARG3$' --trust-cert --log-level info
    }

# Look at specific pools, exclude other pools. Define all flags.
define command{
    command_name    check_vmware_vm_backup_via_ca_include_pools_specify_all
    command_line    $USER1$/check_vmware_vm_backup_via_ca --server '$HOSTNAME$' --domain '$ARG1$' --username '$ARG2$' --password '$ARG3$' --include-rp '$ARG4$' --backup-date-warning '$ARG5$' --backup-date-critical '$ARG6$' --backup-date-timezone '$ARG7$' --backup-date-format '$ARG8$' --backup-date-ca '$ARG9$' --backup-metadata-ca '$ARG10$' --trust-cert --log-level info
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

[strftime-codes]: <https://docs.python.org/3/library/datetime.html#strftime-and-strptime-format-codes> "C standard format codes"

[tz-database]: <https://en.wikipedia.org/wiki/Tz_database> "Time zone database"

[official-time-pkg-docs]: <https://pkg.go.dev/time#pkg-constants> "Go time package"

- Go time package / formatting & parsing
  - <https://pkg.go.dev/time#pkg-constants>
  - <https://pkg.go.dev/time#example-Time.Format>
  - <https://stackoverflow.com/questions/42217308/go-time-format-how-to-understand-meaning-of-2006-01-02-layout>
  - <https://yourbasic.org/golang/format-parse-string-time-date-example/>
  - <https://www.golangprograms.com/get-current-date-and-time-in-various-format-in-golang.html>
  - <https://gobyexample.com/time-formatting-parsing>

- Time zone database
  - <https://en.wikipedia.org/wiki/Tz_database>

<!-- []: PLACEHOLDER "DESCRIPTION_HERE" -->
