// Copyright 2021 Adam Chalkley
//
// https://github.com/atc0005/check-vmware
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package vsphere

// https://godoc.org/github.com/vmware/govmomi/property#example-Collector-Retrieve
// https://github.com/vmware/govmomi/issues/2167
//     pc := property.DefaultCollector(c)
//
//     obj, err := find.NewFinder(c).HostSystem(ctx, "DC0_H0")
//     if err != nil {
//         return err
//     }
//
//     var host mo.HostSystem
//     err = pc.RetrieveOne(ctx, obj.Reference(), []string{"vm"}, &host)
//     if err != nil {
//         return err
//     }
//
//     var vms []mo.VirtualMachine
//     err = pc.Retrieve(ctx, host.Vm, []string{"name"}, &vms)
//     fmt.Printf("host has %d vms:", len(vms))
//     for i := range vms {
//         fmt.Print(" ", vms[i].Name)
//     }
