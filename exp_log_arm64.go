// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

//go:build arm64

package bigmath

// ARM64 wrappers for exp and log
// Temporarily use generic implementations until ARM64 assembly calling convention is fixed
// The assembly files (exp_asm_arm64.s, log_asm_arm64.s) exist but have calling convention issues
// TODO: Fix ARM64 assembly calling convention - ensure proper stack setup for CALL instructions

func bigExpAsm(x *BigFloat, prec uint) *BigFloat {
	return bigExpOptimized(x, prec)
}

func bigLogAsm(x *BigFloat, prec uint) *BigFloat {
	return bigLogOptimized(x, prec)
}
