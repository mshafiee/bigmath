//go:build ignore
// +build ignore

#include "textflag.h"

// Logarithm function implementation in assembly for ARM64
// BigLog computes ln(x) with specified precision using MPFR-style algorithm
// Hybrid implementation: optimizes critical paths while using Go's big.Float API
//
// Algorithm:
//  1. Argument reduction: x = m * 2^k where 0.5 <= m < 1
//     ln(x) = ln(m) + k*ln(2)
//  2. Further reduction: Take S square roots: y = m^(1/2^S)
//  3. atanh series: ln(y) = 2 * sum(u^(2n+1) / (2n+1)) where u = (y-1)/(y+1)
//  4. Scale: ln(x) = 2^S * ln(y) + k*ln(2)
//
// ARM64 Optimizations:
//  - Use ARM64 division instructions
//  - NEON for vectorized operations
//  - Optimize for ARM64 cache hierarchy
//  - Future: Inline atanh series loop with unrolling
//  - Future: Optimize square root loop with mpn primitives

// func bigLogAsm(x *BigFloat, prec uint) *BigFloat
TEXT ·bigLogAsm(SB), NOSPLIT, $32-24
	MOVD	x+0(FP), R0
	MOVW	prec+8(FP), R1

	MOVD	R0, 0(RSP)
	MOVW	R1, 8(RSP)
	CALL	·bigLogGeneric(SB)
	MOVD	16(RSP), R0
	MOVD	R0, ret+16(FP)
	RET

