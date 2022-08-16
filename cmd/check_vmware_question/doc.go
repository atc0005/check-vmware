/*
Nagios plugin used to monitor whether a Virtual Machine is blocked from
execution due to one or more Virtual Machines requiring an interactive
response.

# PURPOSE

This plugin monitors the `question` property of evaluated Virtual Machines.
The status of this property indicates whether an interactive question is
blocking the virtual machine's execution. While a Virtual Machine is in this
state it is not available for normal use.

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
