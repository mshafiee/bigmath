#include "textflag.h"

// Multi-precision fused multiply-add for ARM64
// mpnFMA computes: dst = multiplier * src + addend (all n-limb numbers)
// Returns carry (high limb)
//
// mpnFMA computes fused multiply-add: dst = multiplier * src + addend
// This is useful for patterns like: result = scale * a + b
// Returns carry (high limb)
// func mpnFMA(dst, src, addend *uint64, n int, multiplier uint64) uint64
TEXT ·mpnFMA(SB), NOSPLIT, $0-40
	MOVD	dst+0(FP), R0         // R0 = destination pointer
	MOVD	src+8(FP), R1         // R1 = source pointer
	MOVD	addend+16(FP), R2     // R2 = addend pointer
	MOVD	n+24(FP), R3          // R3 = number of limbs
	MOVD	multiplier+32(FP), R4 // R4 = multiplier

	// Handle zero case
	CBZ	R3, fma_done_zero

	// Clear high part accumulator
	MOVD	$0, R5  // R5 = high part accumulator

	// Check if we have at least 4 limbs
	CMP	$4, R3
	BLT	fma_remainder

	// Main loop: process 4 limbs at a time
fma_loop_unrolled:
	// Load 4 source limbs using load-pair
	LDP	0(R1), (R6, R7)   // R6 = src[0], R7 = src[1]
	LDP	16(R1), (R8, R9)  // R8 = src[2], R9 = src[3]

	// Load 4 addend limbs
	LDP	0(R2), (R10, R11) // R10 = addend[0], R11 = addend[1]
	LDP	16(R2), (R12, R13) // R12 = addend[2], R13 = addend[3]

	// Fused multiply-add limb 0: dst[0] = multiplier * src[0] + addend[0] + high
	MUL	R6, R4, R14       // R14 = low 64 bits of R6 * R4
	UMULH	R6, R4, R15       // R15 = high 64 bits of R6 * R4
	ADDS	R5, R14, R14      // Add previous high part, set flags
	ADC	$0, R15, R15      // Add carry to high part
	ADDS	R10, R14, R14     // Add addend value, set flags
	ADC	$0, R15, R15      // Add carry to high part
	MOVD	R14, 0(R0)        // Store result
	MOVD	R15, R5           // Save high for next iteration

	// Fused multiply-add limb 1
	MUL	R7, R4, R14
	UMULH	R7, R4, R15
	ADDS	R5, R14, R14
	ADC	$0, R15, R15
	ADDS	R11, R14, R14
	ADC	$0, R15, R15
	MOVD	R14, 8(R0)
	MOVD	R15, R5

	// Fused multiply-add limb 2
	MUL	R8, R4, R14
	UMULH	R8, R4, R15
	ADDS	R5, R14, R14
	ADC	$0, R15, R15
	ADDS	R12, R14, R14
	ADC	$0, R15, R15
	MOVD	R14, 16(R0)
	MOVD	R15, R5

	// Fused multiply-add limb 3
	MUL	R9, R4, R14
	UMULH	R9, R4, R15
	ADDS	R5, R14, R14
	ADC	$0, R15, R15
	ADDS	R13, R14, R14
	ADC	$0, R15, R15
	MOVD	R14, 24(R0)
	MOVD	R15, R5

	// Advance pointers by 4 limbs (32 bytes)
	ADD	$32, R1, R1
	ADD	$32, R2, R2
	ADD	$32, R0, R0

	// Decrement counter and loop if >= 4 remain
	SUBS	$4, R3, R3
	CMP	$4, R3
	BGE	fma_loop_unrolled

fma_remainder:
	// Process remaining 0-3 limbs one at a time
	CBZ	R3, fma_done

fma_remainder_loop:
	MOVD	0(R1), R6           // Load source limb
	MOVD	0(R2), R7           // Load addend limb
	MUL	R6, R4, R8          // R8 = low 64 bits
	UMULH	R6, R4, R9          // R9 = high 64 bits
	ADDS	R5, R8, R8          // Add previous high part
	ADC	$0, R9, R9          // Add carry
	ADDS	R7, R8, R8          // Add addend value
	ADC	$0, R9, R9          // Add carry
	MOVD	R8, 0(R0)           // Store result
	MOVD	R9, R5              // Save high for next iteration

	ADD	$8, R1, R1
	ADD	$8, R2, R2
	ADD	$8, R0, R0

	SUBS	$1, R3, R3
	BNE	fma_remainder_loop

fma_done:
	MOVD	R5, ret+40(FP)
	RET

fma_done_zero:
	MOVD	$0, R5
	MOVD	R5, ret+40(FP)
	RET

