// Copyright 2021 Adam Chalkley
//
// https://github.com/atc0005/check-vmware
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package config

import "time"

// Timeout converts the user-specified connection timeout value in
// seconds to an appropriate time duration value for use with setting
// initial connection attempt timeout value.
func (c Config) Timeout() time.Duration {
	return time.Duration(c.timeout) * time.Second
}
