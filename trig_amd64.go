// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

//go:build amd64

package bigmath

// AMD64 assembly implementations for trigonometric functions

// bigSinAMD64 is the basic AMD64 implementation (optimized)
func bigSinAMD64(x *BigFloat, prec uint) *BigFloat {
	return bigSinOptimized(x, prec)
}

// bigSinAVX2 is the AVX2-optimized implementation (optimized)
func bigSinAVX2(x *BigFloat, prec uint) *BigFloat {
	return bigSinOptimized(x, prec)
}

// bigCosAMD64 is the basic AMD64 implementation (optimized)
func bigCosAMD64(x *BigFloat, prec uint) *BigFloat {
	return bigCosOptimized(x, prec)
}

// bigCosAVX2 is the AVX2-optimized implementation (optimized)
func bigCosAVX2(x *BigFloat, prec uint) *BigFloat {
	return bigCosOptimized(x, prec)
}

//go:noescape
func bigSinAsm(x *BigFloat, prec uint) *BigFloat

//go:noescape
func bigCosAsm(x *BigFloat, prec uint) *BigFloat

