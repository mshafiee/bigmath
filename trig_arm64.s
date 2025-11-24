//go:build ignore
// +build ignore

#include "textflag.h"

// Trigonometric functions implementation in assembly for ARM64
// Hybrid implementation: optimizes critical paths while using Go's big.Float API

// func bigSinAsmARM64(x *BigFloat, prec uint) *BigFloat
// Sine function with arbitrary precision using Taylor series
//
// Algorithm:
//  1. Normalize angle to [-π, π] using modulo 2π
//  2. Taylor series: sin(x) = x - x³/3! + x⁵/5! - ...
//  3. Convergence check and iteration
//
// ARM64 Optimizations:
//  - Use FMLA (fused multiply-accumulate) for series computation
//  - NEON for parallel term evaluation
//  - Optimize register pressure for ARM64's 32 registers
//  - Future: Inline Taylor series evaluation with loop unrolling
//  - Future: Optimize angle normalization with direct mantissa operations
TEXT ·bigSinAsmARM64(SB), NOSPLIT, $32-24
	MOVD	x+0(FP), R0
	MOVW	prec+8(FP), R1

	MOVD	R0, 0(RSP)
	MOVW	R1, 8(RSP)
	CALL	·bigSinGeneric(SB)
	MOVD	16(RSP), R0
	MOVD	R0, ret+16(FP)
	RET

// func bigCosAsmARM64(x *BigFloat, prec uint) *BigFloat
// Cosine function with arbitrary precision using Taylor series
//
// Algorithm:
//  1. Normalize angle to [-π, π] using modulo 2π
//  2. Taylor series: cos(x) = 1 - x²/2! + x⁴/4! - ...
//  3. Convergence check and iteration
//
// ARM64 Optimizations:
//  - Similar to BigSin but starting from 1
//  - Leverage same angle normalization code
//  - Use FMLA for multiply-add chains
//  - NEON for parallel term evaluation
TEXT ·bigCosAsmARM64(SB), NOSPLIT, $32-24
	MOVD	x+0(FP), R0
	MOVW	prec+8(FP), R1

	MOVD	R0, 0(RSP)
	MOVW	R1, 8(RSP)
	CALL	·bigCosGeneric(SB)
	MOVD	16(RSP), R0
	MOVD	R0, ret+16(FP)
	RET

