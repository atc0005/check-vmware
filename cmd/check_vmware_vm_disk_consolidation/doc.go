/*

Nagios plugin used to monitor Virtual Machine disk consolidation status.

PURPOSE

This plugin monitors the `consolidationNeeded` property of evaluated Virtual
Machines. The status of this property indicates whether one or more disks for
a Virtual Machine require consolidation. This can happen when a snapshot is
deleted, but its associated disk is not committed back to the base disk. This
situation can cause backup failures and performance issues.

The current design of this plugin is to evaluate *all* Virtual Machines,
whether powered off or powered on. If you have a use case for evaluating
*only* powered on VMs by default, please post it to
https://github.com/atc0005/check-vmware/discussions/176 providing some details
for your use-case.

The output for this plugin is designed to provide the one-line summary needed
by Nagios for quick identification of a problem while providing longer, more
detailed information for use in email and Teams notifications
(https://github.com/atc0005/send2teams).

PROJECT HOME

See our GitHub repo (https://github.com/atc0005/check-vmware) for the latest
code, to file an issue or submit improvements for review and potential
inclusion into the project.

USAGE

See our main README for supported settings and examples.

*/
package main
