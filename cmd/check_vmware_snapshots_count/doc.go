/*

Nagios plugin used to monitor the number of snapshots per Virtual Machine.

PURPOSE

Monitor the number of snapshots for each Virtual Machine. VMware recommends
using no more than 3 or 4 snapshots per Virtual Machine and only for a limited
duration. A maximum of 32 snapshots per Virtual Machine are supported. See
https://kb.vmware.com/s/article/1025279 for more information.

The current design of this plugin is to evaluate *all* Virtual Machines,
whether powered off or powered on. If you have a use case for evaluating
*only* powered on VMs by default, please add a comment to
https://github.com/atc0005/check-vmware/issues/79 providing some details for
your use-case.

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
