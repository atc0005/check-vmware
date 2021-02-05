// Copyright 2021 Adam Chalkley
//
// https://github.com/atc0005/check-vmware
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package vsphere

import "fmt"

// CPUSpeed represents the speed of a CPU
type CPUSpeed float64

// CPU speed values
const (

	// https://stackoverflow.com/questions/34124294/writing-powers-of-10-as-constants-compactly
	// untyped float
	Hz  = 1
	KHz = 1e3
	MHz = 1e6
	GHz = 1e9
	THz = 1e12
	PHz = 1e15
	EHz = 1e18
	ZHz = 1e21
	YHz = 1e24

	// untyped int
	// Hz  = 1
	// KHz = Hz * 1000
	// MHz = KHz * 1000
	// GHz = MHz * 1000
	// THz = GHz * 1000
	// PHz = THz * 1000
	// EHz = PHz * 1000
	// ZHz = EHz * 1000
	// YHz = ZHz * 1000
)

func (cpu CPUSpeed) String() string {
	switch {
	case cpu >= YHz:
		return fmt.Sprintf("%.1f EHz", float64(cpu)/YHz)
	case cpu >= ZHz:
		return fmt.Sprintf("%.1f EHz", float64(cpu)/ZHz)
	case cpu >= EHz:
		return fmt.Sprintf("%.1f EHz", float64(cpu)/EHz)
	case cpu >= PHz:
		return fmt.Sprintf("%.1f PHz", float64(cpu)/PHz)
	case cpu >= THz:
		return fmt.Sprintf("%.1f THz", float64(cpu)/THz)
	case cpu >= GHz:
		return fmt.Sprintf("%.1f GHz", float64(cpu)/GHz)
	case cpu >= MHz:
		return fmt.Sprintf("%.1f MHz", float64(cpu)/MHz)
	case cpu >= KHz:
		return fmt.Sprintf("%.1f KHz", float64(cpu)/KHz)
	}
	return fmt.Sprintf("%d Hz", int(cpu))
}
