// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

//go:build arm64

package bigmath

// ARM64 wrappers for trigonometric functions
// Temporarily use generic implementations until ARM64 assembly calling convention is fixed
// The assembly files (trig_arm64.s) exist but have calling convention issues
// TODO: Fix ARM64 assembly calling convention - ensure proper stack setup for CALL instructions

func bigSinAsmARM64(x *BigFloat, prec uint) *BigFloat {
	return bigSinOptimized(x, prec)
}

func bigCosAsmARM64(x *BigFloat, prec uint) *BigFloat {
	return bigCosOptimized(x, prec)
}
