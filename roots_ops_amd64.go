// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

//go:build amd64

package bigmath

// AMD64 wrapper functions for root functions
// Currently calling generic implementations
// Future: Replace with optimized assembly with unrolled iterations

func bigCbrtAsm(x *BigFloat, prec uint) *BigFloat {
	return bigCbrtOptimized(x, prec)
}

func bigRootAsm(n, x *BigFloat, prec uint) *BigFloat {
	return bigRootOptimized(n, x, prec)
}
