/*
Nagios plugin used to monitor datastore usage.

# PURPOSE

In addition to reporting current datastore usage, this plugin also reports
which VMs reside on the datastore along with their percentage of the total
datastore space used.

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
