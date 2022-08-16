/*
Go-based tooling to monitor VMware environments; **NOT** affiliated with
or endorsed by VMware, Inc.

# PROJECT HOME

See our GitHub repo (https://github.com/atc0005/check-vmware) for the latest
code, to file an issue or submit improvements for review and potential
inclusion into the project.

# PURPOSE

# Monitor VMware environments

# FEATURES

Nagios plugins for monitoring VMware vSphere environments (standalone ESXi
hosts or vCenter instances).

  - VMware Tools

  - Virtual CPU allocations

  - Virtual hardware versions: homogenous, outdated-by threshold range, minimum required and default is minimum required checks

  - Host/Datastore/Virtual Machine pairings (using provided Custom Attribute)

  - Datastore usage

  - Datastore performance

  - Snapshots age

  - Snapshots count

  - Snapshots size

  - Resource Pools: Memory usage

  - Host Memory usage

  - Host CPU usage

  - Virtual Machine (power cycle) uptime

  - Virtual Machine disk consolidation status (with optional forced refresh of Virtual Machine state data)

  - Virtual Machine interactive question status

  - Triggered Alarms in one or more datacenters

  - Last Backup date for VMs (via specified custom attribute)

# USAGE

See our main README for supported settings and examples.
*/
package main
