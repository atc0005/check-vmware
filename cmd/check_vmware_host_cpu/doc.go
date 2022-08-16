/*
Nagios plugin used to monitor ESXi host CPU usage.

# PURPOSE

In addition to reporting current host CPU usage, this plugin also reports
which VMs are on the host (running or not), how much CPU each VM is using as a
fixed value and as a percentage of the host's total CPU capacity.

The output for this plugin is designed to provide the one-line summary needed
by Nagios for quick identification of a problem while providing longer, more
detailed information for use in email and Teams notifications
(https://github.com/atc0005/send2teams).

# PROJECT HOME

See our GitHub repo (https://github.com/atc0005/check-vmware) for the latest
code, to file an issue or submit improvements for review and potential
inclusion into the project.

# USAGE

See our main README for supported settings and examples.
*/
package main
