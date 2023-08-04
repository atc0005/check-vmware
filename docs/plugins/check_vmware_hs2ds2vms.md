<!-- omit in toc -->
# [check-vmware][repo-url] | `check_vmware_hs2ds2vms` plugin

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

Nagios plugin used to monitor host/datastore/vm pairings.

This is a functional plugin responsible for verifying that each VM is housed
on a datastore (best) intended for the host associated with the VM.

By default, the evaluation is limited to powered on VMs, but this can be
toggled to also include powered off VMs.

The association between datastores and hosts is determined by a user-provided
Custom Attribute. Flags for this plugin allow specifying separate Custom
Attribute names for hosts and datastores along with optional separate prefixes
for the provided Custom Attributes.

This allows for example, hosts to use a `Location` Custom Attribute that
shares a datacenter name with datastores using the same `Location` Custom
Attribute. If not specifying a prefix separator, the plugin assumes that a
literal, case-insensitive match of the `Location` field is required. If a
prefix separator is provided, then the separator is used to retrieve the
common prefix for the `Location` Custom Attribute for both hosts and
datastores.

This is intended to work around hosts that may include both the datacenter
name and rack location details in their Custom Attribute (e.g., `Location`).

This plugin optionally allows ignoring a list of datastores, and both hosts
and datastores that are missing the specified Custom Attribute.

In addition to specifying separate Custom Attribute names (required) and
prefix separators (optional), the plugin also accepts a single Custom
Attribute used by both hosts and datastores and an optional prefix separator,
also used by both hosts and datastores.

If specifying a shared Custom Attribute or prefix, per-resource Custom
Attribute flags are rejected (error condition).

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
- `vms_excluded_by_power_state`
- `pairing_issues`
- `datastores`
- `hosts`
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

| Nagios State | Description                                                                  |
| ------------ | ---------------------------------------------------------------------------- |
| `OK`         | Ideal state, no mismatched Host/Datastore/Virtual machine pairings detected. |
| `WARNING`    | Not used by this plugin.                                                     |
| `CRITICAL`   | Any errors encountered or Hosts/Datastores/VM mismatches.                    |

### Command-line arguments

- Use the `-h` or `--help` flag to display current usage information.
- Flags marked as **`required`** must be set via CLI flag.
- Flags *not* marked as required are for settings where a useful default is
  already defined, but may be overridden if desired.

| Flag                 | Required  | Default | Repeat | Possible                                                                | Description                                                                                                                                                                                                                                                                                                                          |
| -------------------- | --------- | ------- | ------ | ----------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `branding`           | No        | `false` | No     | `branding`                                                              | Toggles emission of branding details with plugin status details. This output is disabled by default.                                                                                                                                                                                                                                 |
| `h`, `help`          | No        | `false` | No     | `h`, `help`                                                             | Show Help text along with the list of supported flags.                                                                                                                                                                                                                                                                               |
| `v`, `version`       | No        | `false` | No     | `v`, `version`                                                          | Whether to display application version and then immediately exit application.                                                                                                                                                                                                                                                        |
| `ll`, `log-level`    | No        | `info`  | No     | `disabled`, `panic`, `fatal`, `error`, `warn`, `info`, `debug`, `trace` | Log message priority filter. Log messages with a lower level are ignored. Log messages are sent to `stderr` by default. See [Output](#output) for more information.                                                                                                                                                                  |
| `p`, `port`          | No        | `443`   | No     | *positive whole number between 1-65535, inclusive*                      | TCP port of the remote ESXi host or vCenter instance. This is usually 443 (HTTPS).                                                                                                                                                                                                                                                   |
| `t`, `timeout`       | No        | `10`    | No     | *positive whole number of seconds*                                      | Timeout value in seconds allowed before a plugin execution attempt is abandoned and an error returned.                                                                                                                                                                                                                               |
| `s`, `server`        | **Yes**   |         | No     | *fully-qualified domain name or IP Address*                             | The fully-qualified domain name or IP Address of the remote ESXi host or vCenter instance.                                                                                                                                                                                                                                           |
| `u`, `username`      | **Yes**   |         | No     | *valid username*                                                        | Username with permission to access specified ESXi host or vCenter instance.                                                                                                                                                                                                                                                          |
| `pw`, `password`     | **Yes**   |         | No     | *valid password*                                                        | Password used to login to ESXi host or vCenter instance.                                                                                                                                                                                                                                                                             |
| `domain`             | No        |         | No     | *valid user domain*                                                     | (Optional) domain for user account used to login to ESXi host or vCenter instance. This is needed for user accounts residing in a non-default domain (e.g., SSO specific domain).                                                                                                                                                    |
| `trust-cert`         | No        | `false` | No     | `true`, `false`                                                         | Whether the certificate should be trusted as-is without validation. WARNING: TLS is susceptible to man-in-the-middle attacks if enabling this option.                                                                                                                                                                                |
| `include-rp`         | No        |         | No     | *comma-separated list of resource pool names*                           | Specifies a comma-separated list of Resource Pool names that should be exclusively used when evaluating VMs. Specifying this option will also exclude any VMs from evaluation that are *outside* of a Resource Pool. This option is incompatible with specifying a list of Resource Pool names to ignore or exclude from evaluation. |
| `exclude-rp`         | No        |         | No     | *comma-separated list of resource pool names*                           | Specifies a comma-separated list of Resource Pool names that should be ignored when evaluating VMs. This option is incompatible with specifying a list of Resource Pool names to include for evaluation.                                                                                                                             |
| `ignore-vm`          | No        |         | No     | *comma-separated list of (vSphere) virtual machine names*               | Specifies a comma-separated list of VM names that should be ignored or excluded from evaluation.                                                                                                                                                                                                                                     |
| `ignore-ds`          | No        |         | No     | *comma-separated list of (vSphere) datastore names*                     | Specifies a comma-separated list of Datastore names that should be ignored or excluded from evaluation.                                                                                                                                                                                                                              |
| `powered-off`        | No        | `false` | No     | `true`, `false`                                                         | Toggles evaluation of powered off VMs in addition to powered on VMs. Evaluation of powered off VMs is disabled by default.                                                                                                                                                                                                           |
| `ca-name`            | **Maybe** |         | No     | *valid Custom Attribute name*                                           | Custom Attribute name for host ESXi systems and datastores. Optional if specifying resource-specific custom attribute names.                                                                                                                                                                                                         |
| `ca-prefix-sep`      | **Maybe** |         | No     | *valid Custom Attribute prefix separator character*                     | Custom Attribute prefix separator for host ESXi systems and datastores. Skip if using Custom Attribute values as-is for comparison, otherwise optional if specifying resource-specific custom attribute prefix separator, or using the default separator.                                                                            |
| `ignore-missing-ca`  | No        | `false` | No     | `true`, `false`                                                         | Toggles how missing specified Custom Attributes will be handled. By default, ESXi hosts and datastores missing the Custom Attribute are treated as an error condition.                                                                                                                                                               |
| `host-ca-name`       | **Maybe** |         | No     | *valid Custom Attribute name*                                           | Custom Attribute name specific to host ESXi systems. Optional if specifying shared custom attribute flag.                                                                                                                                                                                                                            |
| `host-ca-prefix-sep` | **Maybe** |         | No     | *valid Custom Attribute prefix separator character*                     | Custom Attribute prefix separator specific to host ESXi systems. Skip if using Custom Attribute values as-is for comparison, otherwise optional if specifying shared custom attribute prefix separator, or using the default separator.                                                                                              |
| `ds-ca-name`         | **Maybe** |         | No     | *valid Custom Attribute name*                                           | Custom Attribute name specific to datastores. Optional if specifying shared custom attribute flag.                                                                                                                                                                                                                                   |
| `ds-ca-prefix-sep`   | **Maybe** |         | No     | *valid Custom Attribute prefix separator character*                     | Custom Attribute prefix separator specific to datastores. Skip if using Custom Attribute values as-is for comparison, otherwise optional if specifying shared custom attribute prefix separator, or using the default separator.                                                                                                     |

### Configuration file

Not currently supported. This feature may be added later if there is
sufficient interest.

## Contrib

See the [main project README](../../README.md) for details.

## Examples

### CLI invocation

```ShellSession
/usr/lib/nagios/plugins/check_vmware_hs2ds2vms --username SERVICE_ACCOUNT_NAME --password "SERVICE_ACCOUNT_PASSWORD" --server vc1.example.com --exclude-rp "Desktops" --ignore-vm "test1.example.com,redmine.example.com,TESTING-AC,RHEL7-TEST" --ca-name "Location" --ca-prefix-sep "-" --trust-cert --log-level info
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
- The Custom Attribute named `Location` is used to dynamically build pairs of
  Hosts and Datastores. Any Host or Datastore missing that Custom Attribute is
  reported as an error condition unless the appropriate CLI flag is provided.
  See the [Configuration options](#configuration-options) section for the flag
  name and further details.
- The Custom Attribute prefix separator `-` is provided in order to "split"
  the value found for the Custom Attribute named `Location` into pairs. The
  second value is thrown away, leaving the first to be used as the `Location`
  value for comparison. VMs running on a host with one value have their
  datastores checked for the same value. If a mismatch is found, this is
  assumed to be a `CRITICAL` level event and reported as such.
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
# /etc/nagios-plugins/config/vmware-host-datastore-vms-pairings.cfg

# Look at all pools, all VMs, do not evaluate any VMs that are powered off.
# Use the same Custom Attribute for hosts and datastores. Use the same Custom
# Attribute prefix separator for hosts and datastores.
#
# This variation of the command is most useful for environments where all VMs
# are monitored equally.
define command{
    command_name   check_vmware_hs2ds2vms
    command_line   $USER1$/check_vmware_hs2ds2vms --server '$HOSTNAME$' --domain '$ARG1$' --username '$ARG2$' --password '$ARG3$' --ca-name '$ARG4$' --ca-prefix-sep '$ARG5$' --trust-cert --log-level info
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
