// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

//go:build amd64

package bigmath

// AMD64 wrapper functions for combinatorics
// Currently calling generic implementations
// Future: Replace with optimized assembly with unrolled multiplication chains

func bigFactorialAsm(n int64, prec uint) *BigFloat {
	return bigFactorialOptimized(n, prec)
}

func bigBinomialAsm(n, k int64, prec uint) *BigFloat {
	return bigBinomialOptimized(n, k, prec)
}
