// Copyright 2021 Adam Chalkley
//
// https://github.com/atc0005/check-vmware
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package vsphere

import (
	"context"
	"time"

	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

// TriggerEntityStateReload accepts a context, a client and a collection of
// ManagedEntity values whose state should be reloaded. This function is used
// when we need to ensure that we are working with the very latest state data
// for a vSphere object.
func TriggerEntityStateReload(ctx context.Context, c *vim25.Client, entities []mo.ManagedEntity) error {

	// https://vdc-download.vmware.com/vmwb-repository/dcr-public/b50dcbbf-051d-4204-a3e7-e1b618c1e384/538cf2ec-b34f-4bae-a332-3820ef9e7773/vim.ManagedEntity.html#reload
	// https://pkg.go.dev/github.com/vmware/govmomi@v0.27.0/vim25/methods#Reload

	funcTimeStart := time.Now()

	defer func(entities []mo.ManagedEntity) {
		logger.Printf(
			"It took %v to execute TriggerEntityStateReload func (for %d entities).\n",
			time.Since(funcTimeStart),
			len(entities),
		)
	}(entities)

	for _, entity := range entities {

		req := types.Reload{
			This: entity.Self,
		}

		logger.Printf(
			"Triggering reload for entity %s of type %s with id %s",
			entity.Name,
			entity.Self.Type,
			entity.Self.Value,
		)

		reloadTimeStart := time.Now()
		_, err := methods.Reload(ctx, c, &req)
		if err != nil {
			return err
		}

		logger.Printf(
			"It took %v to trigger state reload for entity %s",
			time.Since(reloadTimeStart),
			entity.Name,
		)

	}

	return nil

}
