# Snapshot notes

## Notes

- `types.VirtualMachineFileLayoutExSnapshotLayout` has two fields that tie
  together a `VirtualMachine.Snapshot` (`*types.VirtualMachineSnapshotInfo`)
  with a `types.VirtualMachineSnapshotTree`

## Hierarchy of types

Levels above and some below skipped for simplicity and due to my ignorance.

- `VirtualMachine`
  - `VirtualMachineSnapshotInfo`
    - `VirtualMachineSnapshotTree`
  - `VirtualMachineFileLayoutEx`
    - `VirtualMachineFileLayoutExFileInfo`
    - `VirtualMachineFileLayoutExSnapshotLayout`

## Fields expanded

- `vm`
  - `Name` (**Used in output**)
  - `Snapshot` (`VirtualMachineSnapshotInfo`)
    - `RootSnapshotList` (`VirtualMachineSnapshotTree`)
      - `Name` (**Used in output**)
      - `Id` (`int32`)
      - `CreateTime` (**Used in output**)
      - `Snapshot` (`ManagedObjectReference`) (e.g., `snapshot-229099`)
        - `Value`
          - The specific instance of Managed Object this
            `ManagedObjectReference` refers to.
          - *Links to `LayoutEx.Snapshot.Key.Value`*
      - `ChildSnapshotList`
        - *can be multiple levels deep, or null*
  - `LayoutEx`
    - `File`
      - `type` (*filtered to `snapshotData`*)
      - `Key` (e.g., `40`)
        - *links to `LayoutEx.Snapshot.DataKey`*
      - `Size` (**Used in output; by itself & aggregate**)
    - `Snapshot`
      - `DataKey` (e.g., `40`)
        - *links to `LayoutEx.File.Key`*
      - `Key` (`ManagedObjectReference`)  (e.g., `snapshot-229099`)
        - `Value`
          - *links to `vm.Snapshot.Value`*

Superset Custom type:

```golang
// SnapshotSummary is intended to be a summary of the most commonly used
// snapshot details for a specific VirtualMachine snapshot.
type SnapshotSummary struct {
    Name        string
    Description string
    CreateTime  time.Time
    Size        int64
}

// SnapshotSet ties a collection of snapshot summary values to a specific
// VirtualMachine by way of a Managed Object Reference.
type SnapshotSet struct {
    VM        types.ManagedObjectReference
    Snapshots []SnapshotSummary // a collection for easy aggregation
}

// Size returns the size of all snapshots in the set.
func (ss SnapshotSet) Size() int64 {
    var sum int64
    for i := range ss.Snapshots {
        sum += ss.Snapshots[i].Size
    }

    return sum
}
```

Index type used to collect all snapshot details to build aggregate Superset:

```golang

// SnapshotsIndex is a map of ManagedObjectReference (e.g., `snapshot-229099`)
// to ... ?
type SnapshotsIndex map[string]
```

## Building an index

### Current index build notes

Question: Can we build a summary object in one loop, or do we *have* to do it
in two separate steps?

---

- for vm in vms
  - Create SnapshotSet
  - Populate VM field with VirtualMachine MOID
  - Process `vm.Snapshot.RootSnapshotList`

Recursive function:

- loop over `VirtualMachineSnapshotTree`
  - top level of snapshot tree
- create
- grab Name
- grab CreateTime
- grab Snapshot.Value
- check ChildSnapshotList
  - if nil, done

### earlier index build notes

Start with building a map of `Snapshot` `ManagedObjectReference` value
  may require value (from key/value pairs) to be zero value

Loop over map

  create new SnapshotSummary
    this will include a ManagedObjectReference for the snapshot

  for each map key, use that as search string
  Loop over `vm.LayoutEx.Snapshot`
  Filter `Key` to vm.Snapshot`ManagedObjectReference`
  Retrieve `DataKey`
  Loop over `vm.LayoutEx.File`
  Filter to snapshotData
  Filter `Key` to earlier `DataKey` value
  Retrieve Size
  Insert into index

## Scratch

vm.Snapshot entries
Parent snapshot [Name: "pre-upgrade-pre-mass-patch", MOR: snapshot-229098, ID: 27, Has Children: true]
ChildSnapshot [Name: "post-mass-patch-pre-upgrade-tasks", MOR: snapshot-229099, ID: 28, Has Children: true]

for _, snapshotTree := range vm.Snapshot.RootSnapshotList
    snapshotTree.Snapshot.Value == Managed Object Reference instance "ID" (e.g., snapshot-229098)

for _, childSnapshot := range snapshotTree.ChildSnapshotList
    childSnapshot.Snapshot.Value == Managed Object Reference instance "ID" (e.g., snapshot-229099)

loop over vm.LayoutEx.File slice, filter to snapshotData
    Key: 40,
    Size: 28944,
    Name: [DATASTORE-NAME-PLACEHOLDER] server1.example.com/server1.example.com-Snapshot27.vmsn,
    BackingObjectID:

loop over vm.LayoutEx.Snapshot
    Snapshot [DataKey: 40, Key: VirtualMachineSnapshot:snapshot-229098]

`vm.Snapshot` == `*types.VirtualMachineSnapshotInfo`

```golang
type VirtualMachine struct {
    ManagedEntity

    Capability           types.VirtualMachineCapability    `mo:"capability"`
    Config               *types.VirtualMachineConfigInfo   `mo:"config"`
    Layout               *types.VirtualMachineFileLayout   `mo:"layout"`
    LayoutEx             *types.VirtualMachineFileLayoutEx `mo:"layoutEx"`
    Storage              *types.VirtualMachineStorageInfo  `mo:"storage"`
    EnvironmentBrowser   types.ManagedObjectReference      `mo:"environmentBrowser"`
    ResourcePool         *types.ManagedObjectReference     `mo:"resourcePool"`
    ParentVApp           *types.ManagedObjectReference     `mo:"parentVApp"`
    ResourceConfig       *types.ResourceConfigSpec         `mo:"resourceConfig"`
    Runtime              types.VirtualMachineRuntimeInfo   `mo:"runtime"`
    Guest                *types.GuestInfo                  `mo:"guest"`
    Summary              types.VirtualMachineSummary       `mo:"summary"`
    Datastore            []types.ManagedObjectReference    `mo:"datastore"`
    Network              []types.ManagedObjectReference    `mo:"network"`
    Snapshot             *types.VirtualMachineSnapshotInfo `mo:"snapshot"`
    RootSnapshot         []types.ManagedObjectReference    `mo:"rootSnapshot"`
    GuestHeartbeatStatus types.ManagedEntityStatus         `mo:"guestHeartbeatStatus"`
}
```

```golang
type VirtualMachineSnapshotInfo struct {
    DynamicData

    CurrentSnapshot  *ManagedObjectReference      `xml:"currentSnapshot,omitempty"`
    RootSnapshotList []VirtualMachineSnapshotTree `xml:"rootSnapshotList"`
}
```

```golang
type VirtualMachineSnapshotTree struct {
    DynamicData

    Snapshot          ManagedObjectReference       `xml:"snapshot"`
    Vm                ManagedObjectReference       `xml:"vm"`
    Name              string                       `xml:"name"`
    Description       string                       `xml:"description"`
    Id                int32                        `xml:"id,omitempty"`
    CreateTime        time.Time                    `xml:"createTime"`
    State             VirtualMachinePowerState     `xml:"state"`
    Quiesced          bool                         `xml:"quiesced"`
    BackupManifest    string                       `xml:"backupManifest,omitempty"`
    ChildSnapshotList []VirtualMachineSnapshotTree `xml:"childSnapshotList,omitempty"`
    ReplaySupported   *bool                        `xml:"replaySupported"`
}

type VirtualMachineFileLayoutEx struct {
    DynamicData

    File      []VirtualMachineFileLayoutExFileInfo       `xml:"file,omitempty"`
    Disk      []VirtualMachineFileLayoutExDiskLayout     `xml:"disk,omitempty"`
    Snapshot  []VirtualMachineFileLayoutExSnapshotLayout `xml:"snapshot,omitempty"`
    Timestamp time.Time                                  `xml:"timestamp"`
}

type VirtualMachineFileLayoutExFileInfo struct {
    DynamicData

    Key             int32  `xml:"key"`
    Name            string `xml:"name"`
    Type            string `xml:"type"`
    Size            int64  `xml:"size"`
    UniqueSize      int64  `xml:"uniqueSize,omitempty"`
    BackingObjectId string `xml:"backingObjectId,omitempty"`
    Accessible      *bool  `xml:"accessible"`
}

type VirtualMachineFileLayoutExSnapshotLayout struct {
    DynamicData

    Key       ManagedObjectReference                 `xml:"key"`
    DataKey   int32                                  `xml:"dataKey"`
    MemoryKey int32                                  `xml:"memoryKey,omitempty"`
    Disk      []VirtualMachineFileLayoutExDiskLayout `xml:"disk,omitempty"`
}
```

## References

- <https://pkg.go.dev/github.com/vmware/govmomi@v0.24.0/vim25/mo#VirtualMachine>
  - <https://pkg.go.dev/github.com/vmware/govmomi@v0.24.0/vim25/types#VirtualMachineSnapshotInfo>
    - <https://pkg.go.dev/github.com/vmware/govmomi@v0.24.0/vim25/types#VirtualMachineSnapshotTree>
  - <https://pkg.go.dev/github.com/vmware/govmomi@v0.24.0/vim25/types#VirtualMachineFileLayoutEx>
    - <https://pkg.go.dev/github.com/vmware/govmomi@v0.24.0/vim25/types#VirtualMachineFileLayoutExFileInfo>
    - <https://pkg.go.dev/github.com/vmware/govmomi@v0.24.0/vim25/types#VirtualMachineFileLayoutExSnapshotLayout>
