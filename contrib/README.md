<!-- omit in toc -->
# check-vmware | Contrib

[HOME: Main project README](../README.md)

<!-- omit in toc -->
## Table of contents

- [About](#about)
- [Config files](#config-files)
  - [Paths](#paths)
  - [Overview](#overview)
- [References](#references)

## About

This a collection of content pulled from a live production Nagios version 3
"console" hosted on an Ubuntu system. The plan is to update this content (if
needed) in the future to reflect usage with Nagios v4. Please file an issue if
you spot any incompatibilities.

Due to how this content was pulled (and aggressively trimmed in places), this
content does not reflect a full production system. An attempt was made to
provide enough content to give an idea how the plugins in this project can be
used with an actual system, but tries not to pull in *too* much unrelated
configuration detail. There are a few exceptions, primarily the main
`nagios.cfg` config file; this full file (with substitutions in rare places)
is included to illustrate some non-standard directory paths used by the Nagios
version 3 "console".

## Config files

### Paths

All paths listed here are relative to the root of this project, but intended
to map clearly to an Ubuntu system. For example,
`contrib/nagios/etc/nagios-plugins/config/vmware-tools.cfg` refers both to the
file within this project repo, but also to the fully-qualified
`/etc/nagios-plugins/config/vmware-tools.cfg` path on the production Nagios
console these example files were pulled from. You may need to adjust the
deployment path to match your specific environment.

Here are the current config files and structure as illustrated by the `tree`
command:

```ShellSession
$ tree contrib | sed 's/ / /g'
contrib
├── README.md
└── nagios
    └── etc
        ├── nagios-plugins
        │   └── config
        │       ├── send2teams.cfg
        │       ├── vmware-alarms.cfg
        │       ├── vmware-datastores-performance.cfg
        │       ├── vmware-datastores.cfg
        │       ├── vmware-disk-consolidation.cfg
        │       ├── vmware-host-cpu.cfg
        │       ├── vmware-host-datastore-vms-pairings.cfg
        │       ├── vmware-host-memory.cfg
        │       ├── vmware-interactive-question.cfg
        │       ├── vmware-resource-pools.cfg
        │       ├── vmware-snapshots-age.cfg
        │       ├── vmware-snapshots-count.cfg
        │       ├── vmware-snapshots-size.cfg
        │       ├── vmware-tools.cfg
        │       ├── vmware-vcpus.cfg
        │       ├── vmware-virtual-hardware.cfg
        │       └── vmware-vm-power-uptime.cfg
        └── nagios3
            ├── commands.cfg
            ├── conf
            │   ├── contacts
            │   │   ├── helpdesk.cfg
            │   │   └── msteams.cfg
            │   ├── groups
            │   │   ├── contact-groups.cfg
            │   │   ├── host-groups.cfg
            │   │   └── service-groups.cfg
            │   ├── hosts
            │   │   ├── servers
            │   │   │   └── vc1.example.com.cfg
            │   │   └── templates
            │   │       ├── generic-host.cfg
            │   │       ├── generic-linux-box.cfg
            │   │       └── generic-production-host.cfg
            │   └── services
            │       ├── service_host_group_pairings.cfg
            │       └── templates
            │           ├── generic-helpdesk-service.cfg
            │           ├── generic-service.cfg
            │           └── vmware-vsphere-service.cfg
            ├── nagios.cfg
            └── resource.cfg

13 directories, 34 files
```

### Overview

This is a brief overview of the paths/content provided. Please file an issue
if this explanation is unclear, though an expanded description of Nagios
configuration settings is outside the scope of this project. This content is
insufficient to stand up a new Nagios console, but should be useful when
adding additional service checks which make use of plugins from this project.

| Repo path                                                                  | Purpose                                                                                                                                      |
| -------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------- |
| `contrib/nagios/etc/nagios-plugins/config/*.cfg`                           | Supplementary command definition files for Nagios plugins. Commands for plugins in this project are defined here.                            |
|                                                                            |                                                                                                                                              |
| `contrib/nagios/etc/nagios3/commands.cfg`                                  | Primary command definition file for Nagios plugins.                                                                                          |
| `contrib/nagios/etc/nagios3/nagios.cfg`                                    | Primary Nagios configuration file.                                                                                                           |
| `contrib/nagios/etc/nagios3/resource.cfg`                                  | Resource configuration file. This holds `$USERx$` macro definitions referenced in service check and command definitions (e.g., Webhook URL). |
|                                                                            |                                                                                                                                              |
| `contrib/nagios/etc/nagios3/conf/contacts/*.cfg`                           | Contact entry definition files.                                                                                                              |
| `contrib/nagios/etc/nagios3/conf/groups/*.cfg`                             | Contact, Host and Service group definition files.                                                                                            |
| `contrib/nagios/etc/nagios3/conf/hosts/servers/vc1.example.com.cfg`        | Host and Service check definitions for VMware vCenter / vSphere environment. Review alongside plugin command definitions.                    |
| `contrib/nagios/etc/nagios3/conf/hosts/templates/*.cfg`                    | Host templates. Some are used by the VMware vCenter Host definition.                                                                         |
| `contrib/nagios/etc/nagios3/conf/services/service_host_group_pairings.cfg` | Custom shared Service check definitions. This is mostly a placeholder file to satisfy references from other config files.                    |
| `contrib/nagios/etc/nagios3/conf/templates/*.cfg`                          | Service check templates used by the service checks defined in the `vc1.example.com.cfg` file.                                                |

## References

Unordered list of references found within the provided config files. Listed
here for quick/easy access.

- <https://github.com/atc0005/send2teams>
- <https://assets.nagios.com/downloads/nagioscore/docs/nagioscore/3/en/customobjectvars.html>
- <https://www.monitoring-plugins.org/doc/man/check_http.html>
- <https://assets.nagios.com/downloads/nagioscore/docs/nagioscore/3/en/objectdefinitions.html>
- <http://kaotickreation.com/2011/01/30/nagios-check_http/>
- <http://www.jonwitts.co.uk/archives/196>
- <http://linux.101hacks.com/unix/check-http/>
