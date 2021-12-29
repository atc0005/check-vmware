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
	"strings"
	"time"

	"github.com/vmware/govmomi"
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
	userAgent string,
) (*govmomi.Client, error) {

	funcTimeStart := time.Now()

	defer func() {
		logger.Printf(
			"It took %v to execute Login func.\n",
			time.Since(funcTimeStart),
		)
	}()

	vCenterURL := fmt.Sprintf("https://%s:%d/sdk", server, port)

	// TODO: soap.ParseURL automatically adds missing scheme and path. It may
	// be worth using that as a fallback if there are issues logging in?
	u, parseErr := url.Parse(vCenterURL)
	if parseErr != nil {
		return nil, parseErr
	}

	if domain != "" {
		username = strings.Join([]string{username, domain}, "@")
	}

	c, err := govmomi.NewClient(ctx, u, trustCert)
	if err != nil {
		return nil, err
	}

	// Override default user agent
	c.Client.UserAgent = userAgent

	// provide credentials *after* we create the client so that the desired
	// User Agent value can be set before logging in.
	u.User = url.UserPassword(username, password)

	// Login, supplying our custom user agent in place of the default
	authErr := c.Login(ctx, u.User)
	if authErr != nil {
		return nil, authErr
	}

	return c, nil

}
