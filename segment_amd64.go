// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

//go:build amd64

package bigmath

// AMD64 assembly implementations for Chebyshev polynomial evaluation

// evaluateChebyshevBigAMD64 is the basic AMD64 implementation
// Using generic directly - optimized version has allocation overhead
func evaluateChebyshevBigAMD64(t *BigFloat, c []*BigFloat, neval int, prec uint) *BigFloat {
	return evaluateChebyshevBigGeneric(t, c, neval, prec)
}

// evaluateChebyshevBigAVX2 is the AVX2-optimized implementation
// Using generic directly - optimized version has allocation overhead
func evaluateChebyshevBigAVX2(t *BigFloat, c []*BigFloat, neval int, prec uint) *BigFloat {
	return evaluateChebyshevBigGeneric(t, c, neval, prec)
}

// evaluateChebyshevDerivativeBigAMD64 is the basic AMD64 implementation (optimized)
func evaluateChebyshevDerivativeBigAMD64(t *BigFloat, c []*BigFloat, neval int, prec uint) *BigFloat {
	return evaluateChebyshevDerivativeBigOptimized(t, c, neval, prec)
}

// evaluateChebyshevDerivativeBigAVX2 is the AVX2-optimized implementation (optimized)
func evaluateChebyshevDerivativeBigAVX2(t *BigFloat, c []*BigFloat, neval int, prec uint) *BigFloat {
	return evaluateChebyshevDerivativeBigOptimized(t, c, neval, prec)
}

//go:noescape
//nolint:unused // Implemented in assembly (segment_amd64.s)
func evaluateChebyshevBigAsm(t *BigFloat, c []*BigFloat, neval int, prec uint) *BigFloat

//go:noescape
//nolint:unused // Implemented in assembly (segment_amd64.s)
func evaluateChebyshevDerivativeBigAsm(t *BigFloat, c []*BigFloat, neval int, prec uint) *BigFloat
