// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

#include "textflag.h"

// Multi-precision integer multiplication (optimized with MULX and loop unrolling)
// mpn_mul_1 multiplies n-limb number by single limb: dst = src * multiplier
// Returns high limb of result
//
// func mpnMul1(dst, src *uint64, n int, multiplier uint64) uint64
TEXT ·mpnMul1(SB), NOSPLIT, $0-32
	MOVQ	dst+0(FP), DI      // DI = destination pointer
	MOVQ	src+8(FP), SI      // SI = source pointer
	MOVQ	n+16(FP), CX       // CX = number of limbs
	MOVQ	multiplier+24(FP), DX // DX = multiplier (MULX uses RDX as implicit source)

	// Handle zero case
	TESTQ	CX, CX
	JZ	mul1_done_zero

	// Clear high part accumulator
	XORQ	R8, R8  // R8 = high part accumulator (carry from previous multiplication)

	// Fast paths for small n (most common cases: 1-4, 8 limbs)
	CMPQ	CX, $1
	JE	mul1_n1
	CMPQ	CX, $2
	JE	mul1_n2
	CMPQ	CX, $3
	JE	mul1_n3
	CMPQ	CX, $4
	JE	mul1_n4
	CMPQ	CX, $8
	JE	mul1_n8

	// Check if we have at least 4 limbs for loop
	JL	mul1_remainder

	// Main loop: process 4 limbs at a time
	// MULX instruction (BMI2): multiplies RDX * src, stores low in dst1, high in dst2
	// Advantage: Doesn't modify flags, allows better pipelining
mul1_loop_unrolled:
	// Load 4 limbs from source
	MOVQ	0(SI), R9
	MOVQ	8(SI), R10
	MOVQ	16(SI), R11
	MOVQ	24(SI), R12

	// Multiply limb 0: MULXQ src, low, high
	MULXQ	R9, R9, R13      // R13:R9 = RDX * src[0]
	ADDQ	R8, R9           // Add previous high part
	MOVQ	R9, 0(DI)        // Store low result
	ADCQ	$0, R13          // Add carry from ADDQ (if overflow)
	MOVQ	R13, R8          // Save high for next iteration

	// Multiply limb 1
	MULXQ	R10, R10, R13    // R13:R10 = RDX * src[1]
	ADDQ	R8, R10           // Add previous high part
	MOVQ	R10, 8(DI)
	ADCQ	$0, R13
	MOVQ	R13, R8

	// Multiply limb 2
	MULXQ	R11, R11, R13    // R13:R11 = RDX * src[2]
	ADDQ	R8, R11
	MOVQ	R11, 16(DI)
	ADCQ	$0, R13
	MOVQ	R13, R8

	// Multiply limb 3
	MULXQ	R12, R12, R13    // R13:R12 = RDX * src[3]
	ADDQ	R8, R12
	MOVQ	R12, 24(DI)
	ADCQ	$0, R13
	MOVQ	R13, R8

	// Advance pointers by 4 limbs (32 bytes)
	ADDQ	$32, SI
	ADDQ	$32, DI

	// Decrement counter and loop if >= 4 remain
	SUBQ	$4, CX
	CMPQ	CX, $4
	JGE	mul1_loop_unrolled

mul1_remainder:
	// Process remaining 0-3 limbs one at a time
	TESTQ	CX, CX
	JZ	mul1_done

mul1_remainder_loop:
	MOVQ	0(SI), R9
	MULXQ	R9, R9, R13       // R13:R9 = RDX * src[i]
	ADDQ	R8, R9            // Add previous high part
	MOVQ	R9, 0(DI)         // Store result
	ADCQ	$0, R13           // Add carry from ADDQ
	MOVQ	R13, R8           // Save high for next iteration

	ADDQ	$8, SI
	ADDQ	$8, DI

	DECQ	CX
	JNZ	mul1_remainder_loop

mul1_done:
	MOVQ	R8, ret+32(FP)   // Return final high limb
	RET

mul1_done_zero:
	XORQ	AX, AX
	MOVQ	AX, ret+32(FP)
	RET

// Fast paths for small n multiplication (straight-line code, no loops)
mul1_n1:
	MOVQ	0(SI), R9
	MULXQ	R9, R9, R8       // R8:R9 = RDX * src[0]
	MOVQ	R9, 0(DI)
	MOVQ	R8, ret+32(FP)
	RET

mul1_n2:
	MOVQ	0(SI), R9
	MOVQ	8(SI), R10
	MULXQ	R9, R9, R8       // R8:R9 = RDX * src[0]
	MOVQ	R9, 0(DI)
	MULXQ	R10, R10, R11    // R11:R10 = RDX * src[1]
	ADDQ	R8, R10
	MOVQ	R10, 8(DI)
	ADCQ	$0, R11
	MOVQ	R11, ret+32(FP)
	RET

mul1_n3:
	MOVQ	0(SI), R9
	MOVQ	8(SI), R10
	MOVQ	16(SI), R11
	MULXQ	R9, R9, R8       // R8:R9 = RDX * src[0]
	MOVQ	R9, 0(DI)
	MULXQ	R10, R10, R12    // R12:R10 = RDX * src[1]
	ADDQ	R8, R10
	MOVQ	R10, 8(DI)
	ADCQ	$0, R12
	MULXQ	R11, R11, R8     // R8:R11 = RDX * src[2]
	ADDQ	R12, R11
	MOVQ	R11, 16(DI)
	ADCQ	$0, R8
	MOVQ	R8, ret+32(FP)
	RET

mul1_n4:
	MOVQ	0(SI), R9
	MOVQ	8(SI), R10
	MOVQ	16(SI), R11
	MOVQ	24(SI), R12
	MULXQ	R9, R9, R8       // R8:R9 = RDX * src[0]
	MOVQ	R9, 0(DI)
	MULXQ	R10, R10, R13    // R13:R10 = RDX * src[1]
	ADDQ	R8, R10
	MOVQ	R10, 8(DI)
	ADCQ	$0, R13
	MULXQ	R11, R11, R8     // R8:R11 = RDX * src[2]
	ADDQ	R13, R11
	MOVQ	R11, 16(DI)
	ADCQ	$0, R8
	MULXQ	R12, R12, R13    // R13:R12 = RDX * src[3]
	ADDQ	R8, R12
	MOVQ	R12, 24(DI)
	ADCQ	$0, R13
	MOVQ	R13, ret+32(FP)
	RET

mul1_n8:
	// Process first 4 limbs
	MOVQ	0(SI), R9
	MOVQ	8(SI), R10
	MOVQ	16(SI), R11
	MOVQ	24(SI), R12
	MULXQ	R9, R9, R8       // R8:R9 = RDX * src[0]
	MOVQ	R9, 0(DI)
	MULXQ	R10, R10, R13    // R13:R10 = RDX * src[1]
	ADDQ	R8, R10
	MOVQ	R10, 8(DI)
	ADCQ	$0, R13
	MULXQ	R11, R11, R8     // R8:R11 = RDX * src[2]
	ADDQ	R13, R11
	MOVQ	R11, 16(DI)
	ADCQ	$0, R8
	MULXQ	R12, R12, R13    // R13:R12 = RDX * src[3]
	ADDQ	R8, R12
	MOVQ	R12, 24(DI)
	ADCQ	$0, R13
	// Process next 4 limbs
	MOVQ	32(SI), R9
	MOVQ	40(SI), R10
	MOVQ	48(SI), R11
	MOVQ	56(SI), R12
	MULXQ	R9, R9, R8       // R8:R9 = RDX * src[4]
	ADDQ	R13, R9
	MOVQ	R9, 32(DI)
	ADCQ	$0, R8
	MULXQ	R10, R10, R13    // R13:R10 = RDX * src[5]
	ADDQ	R8, R10
	MOVQ	R10, 40(DI)
	ADCQ	$0, R13
	MULXQ	R11, R11, R8     // R8:R11 = RDX * src[6]
	ADDQ	R13, R11
	MOVQ	R11, 48(DI)
	ADCQ	$0, R8
	MULXQ	R12, R12, R13    // R13:R12 = RDX * src[7]
	ADDQ	R8, R12
	MOVQ	R12, 56(DI)
	ADCQ	$0, R13
	MOVQ	R13, ret+32(FP)
	RET

// mpn_addmul_1 multiplies and accumulates: dst += src * multiplier
// Returns carry (high limb)
//
// func mpnAddMul1(dst, src *uint64, n int, multiplier uint64) uint64
TEXT ·mpnAddMul1(SB), NOSPLIT, $0-32
	MOVQ	dst+0(FP), DI      // DI = destination pointer
	MOVQ	src+8(FP), SI      // SI = source pointer
	MOVQ	n+16(FP), CX       // CX = number of limbs
	MOVQ	multiplier+24(FP), DX // DX = multiplier (MULX uses RDX)

	// Handle zero case
	TESTQ	CX, CX
	JZ	addmul1_done_zero

	// Clear high part accumulator
	XORQ	R8, R8  // R8 = high part accumulator

	// Check if we have at least 4 limbs
	CMPQ	CX, $4
	JL	addmul1_remainder

	// Main loop: process 4 limbs at a time
addmul1_loop_unrolled:
	// Load 4 source limbs
	MOVQ	0(SI), R9
	MOVQ	8(SI), R10
	MOVQ	16(SI), R11
	MOVQ	24(SI), R12

	// Load 4 destination limbs
	MOVQ	0(DI), R13
	MOVQ	8(DI), R14
	MOVQ	16(DI), R15
	MOVQ	24(DI), AX  // Use AX as temp

	// Multiply-accumulate limb 0: dst[0] += src[0] * multiplier + high
	MULXQ	R9, R9, BX         // BX:R9 = RDX * src[0]
	ADDQ	R8, R9             // Add previous high part
	ADCQ	$0, BX              // Add carry from ADDQ
	ADDQ	R13, R9             // Add destination value
	ADCQ	$0, BX              // Add carry from ADDQ
	MOVQ	R9, 0(DI)           // Store result
	MOVQ	BX, R8              // Save high for next iteration

	// Multiply-accumulate limb 1
	MULXQ	R10, R10, BX       // BX:R10 = RDX * src[1]
	ADDQ	R8, R10
	ADCQ	$0, BX
	ADDQ	R14, R10
	ADCQ	$0, BX
	MOVQ	R10, 8(DI)
	MOVQ	BX, R8

	// Multiply-accumulate limb 2
	MULXQ	R11, R11, BX       // BX:R11 = RDX * src[2]
	ADDQ	R8, R11
	ADCQ	$0, BX
	ADDQ	R15, R11
	ADCQ	$0, BX
	MOVQ	R11, 16(DI)
	MOVQ	BX, R8

	// Multiply-accumulate limb 3
	MULXQ	R12, R12, BX       // BX:R12 = RDX * src[3]
	ADDQ	R8, R12
	ADCQ	$0, BX
	ADDQ	AX, R12             // Add destination value (was in AX)
	ADCQ	$0, BX
	MOVQ	R12, 24(DI)
	MOVQ	BX, R8

	// Advance pointers by 4 limbs (32 bytes)
	ADDQ	$32, SI
	ADDQ	$32, DI

	// Decrement counter and loop if >= 4 remain
	SUBQ	$4, CX
	CMPQ	CX, $4
	JGE	addmul1_loop_unrolled

addmul1_remainder:
	// Process remaining 0-3 limbs one at a time
	TESTQ	CX, CX
	JZ	addmul1_done

addmul1_remainder_loop:
	MOVQ	0(SI), R9           // Load source limb
	MOVQ	0(DI), R10          // Load destination limb
	MULXQ	R9, R9, BX          // BX:R9 = RDX * src[i]
	ADDQ	R8, R9               // Add previous high part
	ADCQ	$0, BX                // Add carry
	ADDQ	R10, R9              // Add destination value
	ADCQ	$0, BX                // Add carry
	MOVQ	R9, 0(DI)            // Store result
	MOVQ	BX, R8               // Save high for next iteration

	ADDQ	$8, SI
	ADDQ	$8, DI

	DECQ	CX
	JNZ	addmul1_remainder_loop

addmul1_done:
	MOVQ	R8, ret+32(FP)      // Return final high limb
	RET

addmul1_done_zero:
	XORQ	AX, AX
	MOVQ	AX, ret+32(FP)
	RET

