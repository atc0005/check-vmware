/*

Nagios plugin used to monitor datastore performance.

PURPOSE

In addition to reporting current datastore performance, this plugin also
reports which VMs reside on the datastore.

The output for this plugin is designed to provide the one-line summary needed
by Nagios for quick identification of a problem while providing longer, more
detailed information for use in email and Teams notifications
(https://github.com/atc0005/send2teams).

This plugin uses an experimental API and is at greater risk of future breakage
than other plugins in this project (which use stable APIs).

If you use this plugin, please provide feedback by opening a new issue
(https://github.com/atc0005/check-vmware/issues/new) or commenting on the
original discussion thread
(https://github.com/atc0005/check-vmware/issues/505).

PROJECT HOME

See our GitHub repo (https://github.com/atc0005/check-vmware) for the latest
code, to file an issue or submit improvements for review and potential
inclusion into the project.

USAGE

See our main README for supported settings and examples.

*/
package main
