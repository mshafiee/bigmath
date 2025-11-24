// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

//go:build amd64

package bigmath

// Assembly function declarations for AMD64
// These are implemented in rounding_amd64.s

//go:noescape
//nolint:unused // Implemented in assembly (rounding_amd64.s)
func roundToNearestEvenAsm(x *BigFloat, prec uint) *BigFloat

//go:noescape
//nolint:unused // Implemented in assembly (rounding_amd64.s)
func roundTowardZeroAsm(x *BigFloat, prec uint) *BigFloat

//go:noescape
//nolint:unused // Implemented in assembly (rounding_amd64.s)
func roundTowardInfAsm(x *BigFloat, prec uint) *BigFloat

//go:noescape
//nolint:unused // Implemented in assembly (rounding_amd64.s)
func roundTowardNegInfAsm(x *BigFloat, prec uint) *BigFloat
