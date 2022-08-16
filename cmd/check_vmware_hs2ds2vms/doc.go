/*
Nagios plugin used to monitor host/datastore/vm pairings.

# PURPOSE

This is a functional plugin responsible for verifying that each VM is housed
on a datastore (best) intended for the host associated with the VM.

By default, the evaluation is limited to powered on VMs, but this can be
toggled to also include powered off VMs.

The association between datastores and hosts is determined by a user-provided
Custom Attribute. Flags for this plugin allow specifying separate Custom
Attribute names for hosts and datastores along with optional separate prefixes
for the provided Custom Attributes.

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
