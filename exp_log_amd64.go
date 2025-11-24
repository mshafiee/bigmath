// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

//go:build amd64

package bigmath

// AMD64 assembly implementations for exponential and logarithmic functions
// These functions are declared in exp_asm_amd64.s and log_asm_amd64.s
// Currently using optimized Go implementations

//go:noescape
func bigExpAsm(x *BigFloat, prec uint) *BigFloat

//go:noescape
func bigLogAsm(x *BigFloat, prec uint) *BigFloat

// Wrapper functions that use optimized implementations
//nolint:unused // May be used in dispatch or called from assembly
func bigExpAsmWrapper(x *BigFloat, prec uint) *BigFloat {
	return bigExpOptimized(x, prec)
}

//nolint:unused // May be used in dispatch or called from assembly
func bigLogAsmWrapper(x *BigFloat, prec uint) *BigFloat {
	return bigLogOptimized(x, prec)
}
