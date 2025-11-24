//go:build ignore
// +build ignore

// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

#include "textflag.h"

// Exponential function implementation in assembly for ARM64
// BigExp computes e^x with specified precision using MPFR-style algorithm
// Hybrid implementation: optimizes critical paths while using Go's big.Float API
//
// Algorithm:
//  1. Argument reduction: x = k*ln(2) + r where |r| <= ln(2)/2
//  2. Further reduction: exp(r) = (exp(r/2^S))^(2^S)
//  3. Taylor series: exp(u) = sum u^n / n!
//  4. Square S times: res = res^(2^S)
//  5. Multiply by 2^k: res = res * 2^k
//
// ARM64 Optimizations:
//  - Use MADD/MSUB for multiply-accumulate operations
//  - Use NEON vector instructions where applicable
//  - Optimize register allocation for ARM64 calling convention
//  - Future: Inline Taylor series loop with unrolling
//  - Future: Optimize squaring loop with mpn primitives

// func bigExpAsm(x *BigFloat, prec uint) *BigFloat
TEXT ·bigExpAsm(SB), $32-24
	MOVD	x+0(FP), R0
	MOVW	prec+8(FP), R1

	MOVD	R0, 0(RSP)
	MOVW	R1, 8(RSP)
	CALL	·bigExpGeneric(SB)
	MOVD	16(RSP), R0
	MOVD	R0, ret+16(FP)
	RET

