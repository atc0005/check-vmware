/*

Nagios plugin used to monitor virtual hardware versions.

PURPOSE

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

CAVEATS

As of this writing, I am unaware of a way to query the current vSphere
environment for the latest available hardware version. As a workaround for
that lack of knowledge, this plugin applies an automatic baseline of "highest
version discovered" across evaluated VMs. Any VMs with a hardware version not
at that highest version are flagged as problematic. Please file an issue or
open a discussion in this project's repo if you're aware of a way to directly
query the desired value from the current vSphere environment.

Instead of trying to determine how far behind each VM is from the newest
version, this plugin assumes that any deviation is a WARNING level issue.
See GH-33 for future potential changes to this behavior.

*/
package main
