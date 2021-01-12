<!-- omit in toc -->
# Custom Attributes design

<!-- omit in toc -->
## Table of contents

- [About](#about)
- [Notes](#notes)

## About

This is related to the work for GH-8. As part of using Custom Attributes for
determining whether VMs are paired with proper Hosts/Datastores, we first have
to reliably read Custom Attributes from both Hosts and Datastore vSphere
objects and apparently, the inventory at large.

## Notes

For Hosts: Loop over each one, building a slice of host (name, id)

AvailableField: look for name (exact match)
AvailableField: get key

CustomValue: Get all entries
CustomValue: Type assert `cf.(*types.CustomFieldStringValue).Key` for each entry
CustomValue: Compare key from `AvailableField`
CustomValue: if match, retrieve Value
return value

For Datastores: Loop over each one, building a slice of datastore entries

Repeat
...
return datastores

---

Get-DatastoreCustomAttributes
Get-HostSystemCustomAttributes

Attempt direct comparison
Attempt prefix comparison

CLI: Require attribute name
CLI: Accept attribute prefix
CLI: Accept attribute prefix separator

If attribute prefix not set, attempt 1:1 match
If attribute prefix is provided, change match type to prefix only, skip 1:1 match

use default "-" prefix separator if not explicitly provided.

Flags:

- `host-custom-attribute-name`
- `host-custom-attribute-value`
- `host-custom-attribute-value-prefix`
- `host-custom-attribute-value-prefix-separator`

- `datastore-custom-attribute-name`
- `datastore-custom-attribute-value`
- `datastore-custom-attribute-value-prefix`
- `datastore-custom-attribute-value-prefix-separator`

Use validation to:

- reject specifying both `*-custom-attribute-value` and
  `*-custom-attribute-value-prefix`
- require one of `*-custom-attribute-value` and
  `*-custom-attribute-value-prefix`
- (unsure) require `*-prefix` if `*-prefix-separator` is supplied

---

We need a map. For all VMs retrieved, we need to determine datastore and host
and build a list of valid pairings.

Perhaps map of host to list of datastores? The assumption is that a VM object
can quickly reveal which host it is attached to, so from a plugin perspective
we would loop over all VMs and issue an "is in" style check, or a method call
on the map type to determine the same.

Host to datastores mapping makes a lot of sense. The Custom Attribute would be
used initially to create the mapping type, then from there would not be needed
for the bulk of the plugin code.

For some additional validation while building the map (probably a good idea?),
I could use this:
<https://pkg.go.dev/github.com/vmware/govmomi@v0.24.0/object#Datastore.AttachedHosts>

---

The assumption thus far is that the plugin could be used in X service checks,
where X is a datacenter. For example, `DCa` and `DCb`. For the service check
for `DCa`, we would specify `Location` as the Custom Attribute and a prefix of
`DCa`. This would be a mapping of hosts and datastores based on that provided
prefix. The same goes if `DCa` was a literal value to match 1:1.

This would be repeated for another datacenter, `DCb` in our example. This is
how we've previously handled monitoring datastores, hosts and other "single
focus" service checks.

Or, we can specify a Custom Attribute, and a separator and let the plugin
assume that when splitting on the separator that the first value (e.g.,
element 0, `DCa` from `DCa-Rack5-Bay10`) becomes our required prefix for
host/datastore pairings.

Perhaps this could be a "stretch" goal of sorts?

With this design, the following flags could be used:

- custom-attribute-name
- custom-attribute-prefix-separator

This design would assume that the lack of a separator is equivalent to element
0 from a "split" operation on the custom attribute value.

Based on the earlier design, this latter design could work as an alternative.
Not sure if making it a separate plugin would be beneficial.

---

A further enhancement would be to add a flag to enforce both the presence of
the provided custom attribute and a value for hosts and datastores. If either
are found to be missing the Custom Attribute, an error could be thrown putting
the service check into a `CRITICAL` state.

Maybe this could be a default state, with a flag such as
`--ignore-missing-custom-attribute`?

Or, a non-default state and a flag such as `--require-custom-attribute`.
