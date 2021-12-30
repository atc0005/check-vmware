<!-- omit in toc -->
# [check-vmware][repo-url] | `check_vmware_alarms` plugin

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

Nagios plugin used to monitor for Triggered Alarms in one or more datacenters.

- Explicit exclusions take priority over either implicit or explicit
  inclusions.
- All filtering is currently applied in batches/bulk.

It helps to think of the process working this way for each filter in the
"pipeline":

1. Explicit inclusions are applied, marking matching triggered alarms as
   explicitly included and non-matches as *implicitly* excluded
1. Explicit exclusions are applied, marking matching triggered alarms as
   explicitly excluded, permanently "dropping" the triggered alarm from
   further evaluation
1. After all filters have finished processing, any triggered alarms marked as
   excluded (implicit or explicit) are removed from final evaluation (i.e.,
   ignored and not reported as a problem).

Filtering is available for explicitly *including* or *excluding* based on:

- `Acknowledged` status
- [Managed Entity type][vsphere-managed-object-reference] (e.g.,
  `Datastore`, `VirtualMachine`) associated with the Triggered Alarm
- Inventory object `name` (e.g., `node1.example.com`, `vc1.example.com`)
  associated with the Triggered Alarm
- Alarm `Name` field substring match
- Alarm `Description` field substring match
- Triggered Alarm `Status` (e.g., `red`, `yellow`, `gray`)
- `Resource Pool` for the [Managed Entity
  type][vsphere-managed-object-reference] (e.g., `ResourcePool`,
  `VirtualMachine`) associated with the Triggered Alarm

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
- `datacenters`
- `triggered_alarms`
- `triggered_alarms_included`
- `triggered_alarms_excluded`
- `triggered_alarms_critical`
- `triggered_alarms_warning`
- `triggered_alarms_unknown`
- `triggered_alarms_ok`

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

| Nagios State | Description                                             |
| ------------ | ------------------------------------------------------- |
| `OK`         | Ideal state, no non-excluded Triggered Alarms detected. |
| `WARNING`    | One or more non-excluded alarms with a yellow status.   |
| `CRITICAL`   | One or more non-excluded alarms with a red status.      |

### Command-line arguments

- Use the `-h` or `--help` flag to display current usage information.
- Flags marked as **`required`** must be set via CLI flag.
- Flags *not* marked as required are for settings where a useful default is
  already defined, but may be overridden if desired.

| Flag                  | Required | Default | Repeat | Possible                                                                                                                                                                       | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 |
| --------------------- | -------- | ------- | ------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `branding`            | No       | `false` | No     | `branding`                                                                                                                                                                     | Toggles emission of branding details with plugin status details. This output is disabled by default.                                                                                                                                                                                                                                                                                                                                                                                                        |
| `h`, `help`           | No       | `false` | No     | `h`, `help`                                                                                                                                                                    | Show Help text along with the list of supported flags.                                                                                                                                                                                                                                                                                                                                                                                                                                                      |
| `v`, `version`        | No       | `false` | No     | `v`, `version`                                                                                                                                                                 | Whether to display application version and then immediately exit application.                                                                                                                                                                                                                                                                                                                                                                                                                               |
| `ll`, `log-level`     | No       | `info`  | No     | `disabled`, `panic`, `fatal`, `error`, `warn`, `info`, `debug`, `trace`                                                                                                        | Log message priority filter. Log messages with a lower level are ignored. Log messages are sent to `stderr` by default. See [Output](#output) for more information.                                                                                                                                                                                                                                                                                                                                         |
| `p`, `port`           | No       | `443`   | No     | *positive whole number between 1-65535, inclusive*                                                                                                                             | TCP port of the remote ESXi host or vCenter instance. This is usually 443 (HTTPS).                                                                                                                                                                                                                                                                                                                                                                                                                          |
| `t`, `timeout`        | No       | `10`    | No     | *positive whole number of seconds*                                                                                                                                             | Timeout value in seconds allowed before a plugin execution attempt is abandoned and an error returned.                                                                                                                                                                                                                                                                                                                                                                                                      |
| `s`, `server`         | **Yes**  |         | No     | *fully-qualified domain name or IP Address*                                                                                                                                    | The fully-qualified domain name or IP Address of the remote ESXi host or vCenter instance.                                                                                                                                                                                                                                                                                                                                                                                                                  |
| `u`, `username`       | **Yes**  |         | No     | *valid username*                                                                                                                                                               | Username with permission to access specified ESXi host or vCenter instance.                                                                                                                                                                                                                                                                                                                                                                                                                                 |
| `pw`, `password`      | **Yes**  |         | No     | *valid password*                                                                                                                                                               | Password used to login to ESXi host or vCenter instance.                                                                                                                                                                                                                                                                                                                                                                                                                                                    |
| `domain`              | No       |         | No     | *valid user domain*                                                                                                                                                            | (Optional) domain for user account used to login to ESXi host or vCenter instance. This is needed for user accounts residing in a non-default domain (e.g., SSO specific domain).                                                                                                                                                                                                                                                                                                                           |
| `trust-cert`          | No       | `false` | No     | `true`, `false`                                                                                                                                                                | Whether the certificate should be trusted as-is without validation. WARNING: TLS is susceptible to man-in-the-middle attacks if enabling this option.                                                                                                                                                                                                                                                                                                                                                       |
| `dc-name`             | No       |         | No     | *comma-separated list of valid vSphere datacenter names*                                                                                                                       | Specifies the name of one or more vSphere Datacenters. If not specified, applicable plugins will attempt to evaluate all visible datacenters found in the vSphere environment. Not applicable to standalone ESXi hosts.                                                                                                                                                                                                                                                                                     |
| `include-entity-type` | No       |         | No     | [*comma-separated list of valid managed object type keywords*][vsphere-managed-object-reference]                                                                               | If specified, triggered alarms will only be evaluated if the associated entity type (e.g., `Datastore`) matches one of the specified values; while multiple explicit inclusions are allowed, explicit exclusions have precedence over explicit inclusions and will exclude the triggered alarm from further evaluation.                                                                                                                                                                                     |
| `exclude-entity-type` | No       |         | No     | [*comma-separated list of valid managed object type keywords*][vsphere-managed-object-reference]                                                                               | If specified, triggered alarms will only be evaluated if the associated entity type (e.g., `Datastore`) does NOT match one of the specified values; while multiple explicit inclusions are allowed, explicit exclusions have precedence over explicit inclusions and will exclude the triggered alarm from further evaluation.                                                                                                                                                                              |
| `include-entity-name` | No       |         | No     | *comma-separated list of vSphere inventory object names*                                                                                                                       | If specified, triggered alarms will only be evaluated if the associated entity name (e.g., `node1.example.com`) matches one of the specified values; while multiple explicit inclusions are allowed, explicit exclusions have precedence over explicit inclusions and will exclude the triggered alarm from further evaluation.                                                                                                                                                                             |
| `exclude-entity-name` | No       |         | No     | *comma-separated list of vSphere inventory object names*                                                                                                                       | If specified, triggered alarms will only be evaluated if the associated entity name (e.g., `node1.example.com`) does NOT match one of the specified values; while multiple explicit inclusions are allowed, explicit exclusions have precedence over explicit inclusions and will exclude the triggered alarm from further evaluation.                                                                                                                                                                      |
| `include-entity-rp`   | No       |         | No     | *comma-separated list of resource pool names*                                                                                                                                  | If specified, triggered alarms will only be evaluated if the associated entity is part of one of the specified Resource Pools (case-insensitive match on the name) and is not explicitly excluded by another filter in the pipeline; while multiple explicit inclusions are allowed, explicit exclusions have precedence over explicit inclusions and will exclude the triggered alarm from further evaluation.                                                                                             |
| `exclude-entity-rp`   | No       |         | No     | *comma-separated list of resource pool names*                                                                                                                                  | If specified, triggered alarms will only be evaluated if the associated entity is NOT part of one of the specified Resource Pools (case-insensitive match on the name) and is not explicitly excluded by another filter in the pipeline; while multiple explicit inclusions are allowed, explicit exclusions have precedence over explicit inclusions and will exclude the triggered alarm from further evaluation.                                                                                         |
| `eval-acknowledged`   | No       | `false` | No     | `true`, `false`                                                                                                                                                                | Toggles evaluation of acknowledged triggered alarms in addition to unacknowledged triggered alarms. Evaluation of acknowledged alarms is disabled by default.                                                                                                                                                                                                                                                                                                                                               |
| `include-name`        | No       |         | No     | *valid custom or* [*default alarm names*][vsphere-default-alarms]                                                                                                              | If specified, triggered alarms will only be evaluated if the alarm name (e.g., `Datastore usage on disk`) case-insensitively matches one of the specified substring values (e.g., `datastore` or `datastore usage`) and is not explicitly excluded by another filter in the pipeline; while multiple explicit inclusions are allowed, explicit exclusions have precedence over explicit inclusions and will exclude the triggered alarm from further evaluation.                                            |
| `exclude-name`        | No       |         | No     | *valid custom or* [*default alarm names*][vsphere-default-alarms]                                                                                                              | If specified, triggered alarms will only be evaluated if the alarm name (e.g., `Datastore usage on disk`) DOES NOT case-insensitively match one of the specified substring values (e.g., `datastore` or `datastore usage`) and is not explicitly excluded by another filter in the pipeline; while multiple explicit inclusions are allowed, explicit exclusions have precedence over explicit inclusions and will exclude the triggered alarm from further evaluation.                                     |
| `include-desc`        | No       |         | No     | *valid custom or* [*default alarm descriptions*][vsphere-default-alarms]                                                                                                       | If specified, triggered alarms will only be evaluated if the alarm description (e.g., `Default alarm to monitor datastore disk usage`) case-insensitively matches one of the specified substring values (e.g., `datastore disk` or `monitor datastore`) and is not explicitly excluded by another filter in the pipeline; while multiple explicit inclusions are allowed, explicit exclusions have precedence over explicit inclusions and will exclude the triggered alarm from further evaluation.        |
| `exclude-desc`        | No       |         | No     | *valid custom or* [*default alarm descriptions*][vsphere-default-alarms]                                                                                                       | If specified, triggered alarms will only be evaluated if the alarm description (e.g., `Default alarm to monitor datastore disk usage`) DOES NOT case-insensitively match one of the specified substring values (e.g., `datastore disk` or `monitor datastore`) and is not explicitly excluded by another filter in the pipeline; while multiple explicit inclusions are allowed, explicit exclusions have precedence over explicit inclusions and will exclude the triggered alarm from further evaluation. |
| `include-status`      | No       |         | No     | *valid* [*managed entity status*][vsphere-manged-entity-status] (excluding `green`) or [Nagios state][nagios-state-types] (excluding `OK`) (`WARNING`, `CRITICAL` , `UNKNOwN`) | If specified, triggered alarms will only be evaluated if the alarm status (e.g., `yellow`) case-insensitively matches one of the specified keywords (e.g., `yellow` or `warning`) and is not explicitly excluded by another filter in the pipeline; while multiple explicit inclusions are allowed, explicit exclusions have precedence over explicit inclusions and will exclude the triggered alarm from further evaluation.                                                                              |
| `exclude-status`      | No       |         | No     | *valid* [*managed entity status*][vsphere-manged-entity-status]                                                                                                                | If specified, triggered alarms will only be evaluated if the alarm status (e.g., `yellow`) DOES NOT case-insensitively match one of the specified keywords (e.g., `yellow` or `warning`) and is not explicitly excluded by another filter in the pipeline; while multiple explicit inclusions are allowed, explicit exclusions have precedence over explicit inclusions and will exclude the triggered alarm from further evaluation.                                                                       |

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
/usr/lib/nagios/plugins/check_vmware_alarms --username SERVICE_ACCOUNT_NAME --password "SERVICE_ACCOUNT_PASSWORD" --server vc1.example.com  --trust-cert --log-level info
```

See the [configuration options](#configuration-options) section for all
command-line settings supported by this plugin along with descriptions of
each. See the [contrib](#contrib) section for information regarding example
command definitions and Nagios configuration files.

Of note:

- Triggered alarms are evaluated for all detected datacenters
  - due to lack of specified datacenter name (or names)
- Triggered alarms are not filtered based on associated [managed
  object][vsphere-managed-object-reference] (aka, `managed entity`) type
  - due to lack of explicit exclusions or inclusions
- Triggered alarms are not filtered based on associated [managed
  object][vsphere-managed-object-reference] (aka, `managed entity`) name
  - due to lack of explicit exclusions or inclusions
- Triggered alarms are not filtered based on associated [managed
  object][vsphere-managed-object-reference] (aka, `managed entity`) resource
  pool
  - due to lack of explicit exclusions or inclusions
- Triggered alarms that were previously acknowledged are ignored
- Triggered alarms are *not* filtered based on defined Alarm name
  - due to lack of explicit exclusions or inclusions
- Triggered alarms are *not* filtered based on defined Alarm description
  - due to lack of explicit exclusions or inclusions
- Triggered alarms are *not* filtered based on Triggered Alarm status
  - due to lack of explicit exclusions or inclusions
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
# /etc/nagios-plugins/config/vmware-alarms.cfg

# Look at triggered alarms across all detected datacenters, do not evaluate
# any triggered alarms which have been previously acknowledged.
define command{
    command_name    check_vmware_alarms
    command_line    $USER1$/check_vmware_alarms --server '$HOSTNAME$' --domain '$ARG1$' --username '$ARG2$' --password '$ARG3$' --trust-cert --log-level info
    }

# Look at triggered alarms within specified datacenters. Do not evaluate any
# triggered alarms which have been previously acknowledged.
define command{
    command_name    check_vmware_alarms_specific_dc
    command_line    $USER1$/check_vmware_alarms --server '$HOSTNAME$' --domain '$ARG1$' --username '$ARG2$' --password '$ARG3$' --dc-name '$ARG4$' --trust-cert --log-level info
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

[vsphere-managed-object-reference]: <https://vdc-download.vmware.com/vmwb-repository/dcr-public/a5f4000f-1ea8-48a9-9221-586adff3c557/7ff50256-2cf2-45ea-aacd-87d231ab1ac7/vmodl.ManagedObjectReference.html> "Data Object - ManagedObjectReference(vmodl.ManagedObjectReference)"

[vsphere-manged-entity-status]: <https://vdc-repo.vmware.com/vmwb-repository/dcr-public/91f5f971-bf1d-4904-9942-37c6109da8a3/b79fa83f-dc4e-491d-9785-dc9d91aa0c67/doc/vim.ManagedEntity.Status.html>

[vsphere-default-alarms]: <https://docs.vmware.com/en/VMware-vSphere/7.0/com.vmware.vsphere.monitoring.doc/GUID-82933270-1D72-4CF3-A1AF-E5A1343F62DE.html>

[nagios-state-types]: <https://assets.nagios.com/downloads/nagioscore/docs/nagioscore/3/en/statetypes.html>

<!-- []: PLACEHOLDER "DESCRIPTION_HERE" -->
