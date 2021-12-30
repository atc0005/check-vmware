# Mappings

| Field                                           | Maps to                             | Description                             | Example                 |
| ----------------------------------------------- | ----------------------------------- | --------------------------------------- | ----------------------- |
| `vm.LayoutEx.Disk[].Chain[].FileKey`            |                                     | diskDescriptor, diskExtent pairs        | [3 4], [11 12], [24 25] |
| `vm.LayoutEx.Disk[].Key`                        | `vm.LayoutEx.Snapshot[].Disk[].Key` |                                         | 2000                    |
| `vm.LayoutEx.Snapshot[].Disk[].Key`             | `vm.LayoutEx.Disk[].Key`            |                                         | 2000                    |
| `vm.LayoutEx.Snapshot[].Disk[].Chain[].FileKey` |                                     | diskDescriptor, diskExtent pairs        | [3 4], [11 12], [24 25] |
| `vm.LayoutEx.Snapshot[].Key.Value`              |                                     | Managed Object Reference                | snapshot-163887         |
| `vm.LayoutEx.Snapshot[].DataKey`                | `vm.LayoutEx.File[].Key`            | individual file keys (logs, vmdks, etc) | 13, 14, 15, 16, 26, 29  |
| `vm.LayoutEx.File[].Key`                        | `vm.LayoutEx.Snapshot[].DataKey`    | individual file keys (logs, vmdks, etc) | 13, 14, 15, 16, 26, 29  |
| `vm.LayoutEx.File[].Size`                       |                                     | file size in bytes (int64)              | 8.4 GB                  |
