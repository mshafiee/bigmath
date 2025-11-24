// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

//go:build amd64

package bigmath

// AMD64 implementations for exponential and logarithmic functions
// Using optimized Go implementations instead of assembly wrappers
// to avoid GC stackmap issues

func bigExpAsm(x *BigFloat, prec uint) *BigFloat {
	return bigExpOptimized(x, prec)
}

func bigLogAsm(x *BigFloat, prec uint) *BigFloat {
	return bigLogOptimized(x, prec)
}
