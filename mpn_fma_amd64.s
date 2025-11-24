#include "textflag.h"

// Multi-precision fused multiply-add
// mpnFMA computes: dst = multiplier * src + addend (all n-limb numbers)
// This is useful for patterns like: result = scale * a + b
// Returns carry (high limb)
//
// mpnFMA computes fused multiply-add: dst = multiplier * src + addend
// This is useful for patterns like: result = scale * a + b
// Returns carry (high limb)
// func mpnFMA(dst, src, addend *uint64, n int, multiplier uint64) uint64
TEXT ·mpnFMA(SB), NOSPLIT, $0-40
	MOVQ	dst+0(FP), DI         // DI = destination pointer
	MOVQ	src+8(FP), SI         // SI = source pointer
	MOVQ	addend+16(FP), R9     // R9 = addend pointer (save DX for multiplier)
	MOVQ	n+24(FP), CX          // CX = number of limbs
	MOVQ	multiplier+32(FP), DX // DX = multiplier (MULX uses RDX as implicit source)

	// Handle zero case
	TESTQ	CX, CX
	JZ	fma_done_zero

	// Clear high part accumulator
	XORQ	R8, R8  // R8 = high part accumulator

	// Check if we have at least 4 limbs
	CMPQ	CX, $4
	JL	fma_remainder

	// Main loop: process 4 limbs at a time
	// Compute: dst[i] = multiplier * src[i] + addend[i] + high
fma_loop_unrolled:
	// Load 4 source limbs
	MOVQ	0(SI), R10
	MOVQ	8(SI), R11
	MOVQ	16(SI), R12
	MOVQ	24(SI), R13

	// Load 4 addend limbs (using R9 as addend pointer)
	MOVQ	0(R9), R14
	MOVQ	8(R9), R15
	MOVQ	16(R9), AX
	MOVQ	24(R9), BX

	// Fused multiply-add limb 0: dst[0] = multiplier * src[0] + addend[0] + high
	MULXQ	R10, R10, BP         // BP:R10 = RDX * src[0]
	ADDQ	R8, R10              // Add previous high part
	ADCQ	$0, BP               // Add carry from ADDQ
	ADDQ	R14, R10             // Add addend value
	ADCQ	$0, BP               // Add carry from ADDQ
	MOVQ	R10, 0(DI)           // Store result
	MOVQ	BP, R8               // Save high for next iteration

	// Fused multiply-add limb 1
	MULXQ	R11, R11, BP         // BP:R11 = RDX * src[1]
	ADDQ	R8, R11
	ADCQ	$0, BP
	ADDQ	R15, R11
	ADCQ	$0, BP
	MOVQ	R11, 8(DI)
	MOVQ	BP, R8

	// Fused multiply-add limb 2
	MULXQ	R12, R12, BP         // BP:R12 = RDX * src[2]
	ADDQ	R8, R12
	ADCQ	$0, BP
	ADDQ	AX, R12              // Add addend value
	ADCQ	$0, BP
	MOVQ	R12, 16(DI)
	MOVQ	BP, R8

	// Fused multiply-add limb 3
	MULXQ	R13, R13, BP         // BP:R13 = RDX * src[3]
	ADDQ	R8, R13
	ADCQ	$0, BP
	ADDQ	BX, R13              // Add addend value
	ADCQ	$0, BP
	MOVQ	R13, 24(DI)
	MOVQ	BP, R8

	// Advance pointers by 4 limbs (32 bytes)
	ADDQ	$32, SI
	ADDQ	$32, R9               // Advance addend pointer
	ADDQ	$32, DI

	// Decrement counter and loop if >= 4 remain
	SUBQ	$4, CX
	CMPQ	CX, $4
	JGE	fma_loop_unrolled

fma_remainder:
	// Process remaining 0-3 limbs one at a time
	TESTQ	CX, CX
	JZ	fma_done

fma_remainder_loop:
	MOVQ	0(SI), R10           // Load source limb
	MOVQ	0(R9), R11           // Load addend limb
	MULXQ	R10, R10, BP         // BP:R10 = RDX * src[i]
	ADDQ	R8, R10              // Add previous high part
	ADCQ	$0, BP               // Add carry
	ADDQ	R11, R10             // Add addend value
	ADCQ	$0, BP               // Add carry
	MOVQ	R10, 0(DI)           // Store result
	MOVQ	BP, R8               // Save high for next iteration

	ADDQ	$8, SI
	ADDQ	$8, R9                // Advance addend pointer
	ADDQ	$8, DI

	DECQ	CX
	JNZ	fma_remainder_loop

fma_done:
	MOVQ	R8, ret+40(FP)       // Return final high limb
	RET

fma_done_zero:
	XORQ	AX, AX
	MOVQ	AX, ret+40(FP)
	RET
