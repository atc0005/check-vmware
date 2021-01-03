// Copyright 2021 Adam Chalkley
//
// https://github.com/atc0005/check-vmware
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package vsphere

import (
	"context"
	"fmt"
	"net/url"

	"github.com/vmware/govmomi/session/cache"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/soap"
)

// Login receives credentials and related settings used to handle creating a
// new client and logging into a specified vSphere environment. The
// initialized and logged-in client is returned for further use.
func Login(
	ctx context.Context,
	server string,
	port int,
	trustCert bool,
	username string,
	domain string,
	password string,
) (*vim25.Client, error) {

	// TODO: Do we really need to support user domains?

	vCenterURL := fmt.Sprintf("https://%s:%d/sdk", server, port)

	// TODO: soap.ParseURL automatically adds missing scheme and path. It may
	// be worth using that as a fallback if there are issues logging in?
	u, parseErr := soap.ParseURL(vCenterURL)
	if parseErr != nil {
		return nil, parseErr
	}

	u.User = url.UserPassword(username, password)

	// Use session cache to help avoid "leaking sessions"; Session.Login will
	// only create a new authenticated session if the cached session does not
	// exist or is invalid.
	s := &cache.Session{
		URL:      u,
		Insecure: trustCert,
	}

	c := new(vim25.Client)
	loginErr := s.Login(ctx, c, nil)
	if loginErr != nil {
		return nil, loginErr
	}

	return c, nil

}
