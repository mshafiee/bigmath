#include "textflag.h"

// Multi-precision integer multiplication for ARM64 (optimized with loop unrolling)
// mpn_mul_1 multiplies n-limb number by single limb: dst = src * multiplier
// Returns high limb of result
//
// func mpnMul1(dst, src *uint64, n int, multiplier uint64) uint64
TEXT ·mpnMul1(SB), NOSPLIT, $0-32
	MOVD	dst+0(FP), R0      // R0 = destination pointer
	MOVD	src+8(FP), R1      // R1 = source pointer
	MOVD	n+16(FP), R2       // R2 = number of limbs
	MOVD	multiplier+24(FP), R3 // R3 = multiplier

	// Handle zero case
	CBZ	R2, mul1_done_zero

	// Clear high part accumulator
	MOVD	$0, R4  // R4 = high part accumulator

	// Fast paths for small n (most common cases: 1-4, 8 limbs)
	CMP	$1, R2
	BEQ	mul1_n1
	CMP	$2, R2
	BEQ	mul1_n2
	CMP	$3, R2
	BEQ	mul1_n3
	CMP	$4, R2
	BEQ	mul1_n4
	CMP	$8, R2
	BEQ	mul1_n8

	// Check if we have at least 4 limbs for loop
	BLT	mul1_remainder

	// Main loop: process 4 limbs at a time
mul1_loop_unrolled:
	// Load 4 source limbs using load-pair
	LDP	0(R1), (R5, R6)   // R5 = src[0], R6 = src[1]
	LDP	16(R1), (R7, R8)  // R7 = src[2], R8 = src[3]

	// Multiply limb 0: MUL gives low 64 bits, UMULH gives high 64 bits
	MUL	R5, R3, R9        // R9 = low 64 bits of R5 * R3
	UMULH	R5, R3, R10       // R10 = high 64 bits of R5 * R3
	ADDS	R4, R9, R9        // Add previous high part, set flags
	MOVD	R9, 0(R0)         // Store low result
	ADC	$0, R10, R10      // Add carry (0 or 1) to high part
	MOVD	R10, R4           // Save high for next iteration

	// Multiply limb 1
	MUL	R6, R3, R9
	UMULH	R6, R3, R10
	ADDS	R4, R9, R9
	MOVD	R9, 8(R0)
	ADC	$0, R10, R10
	MOVD	R10, R4

	// Multiply limb 2
	MUL	R7, R3, R9
	UMULH	R7, R3, R10
	ADDS	R4, R9, R9
	MOVD	R9, 16(R0)
	ADC	$0, R10, R10
	MOVD	R10, R4

	// Multiply limb 3
	MUL	R8, R3, R9
	UMULH	R8, R3, R10
	ADDS	R4, R9, R9
	MOVD	R9, 24(R0)
	ADC	$0, R10, R10
	MOVD	R10, R4

	// Advance pointers by 4 limbs (32 bytes)
	ADD	$32, R1, R1
	ADD	$32, R0, R0

	// Decrement counter and loop if >= 4 remain
	SUBS	$4, R2, R2
	CMP	$4, R2
	BGE	mul1_loop_unrolled

mul1_remainder:
	// Process remaining 0-3 limbs one at a time
	CBZ	R2, mul1_done

mul1_remainder_loop:
	MOVD	0(R1), R5
	MUL	R5, R3, R9        // R9 = low 64 bits
	UMULH	R5, R3, R10       // R10 = high 64 bits
	ADDS	R4, R9, R9        // Add previous high part
	MOVD	R9, 0(R0)         // Store result
	ADC	$0, R10, R10      // Add carry to high part
	MOVD	R10, R4           // Save high for next iteration

	ADD	$8, R1, R1
	ADD	$8, R0, R0

	SUBS	$1, R2, R2
	BNE	mul1_remainder_loop

mul1_done:
	MOVD	R4, ret+32(FP)
	RET

mul1_done_zero:
	MOVD	$0, R4
	MOVD	R4, ret+32(FP)
	RET

// Fast paths for small n multiplication (straight-line code, no loops)
mul1_n1:
	MOVD	0(R1), R5
	MUL	R5, R3, R9        // R9 = low 64 bits
	UMULH	R5, R3, R4        // R4 = high 64 bits
	MOVD	R9, 0(R0)
	MOVD	R4, ret+32(FP)
	RET

mul1_n2:
	LDP	0(R1), (R5, R6)
	MUL	R5, R3, R9
	UMULH	R5, R3, R4
	MOVD	R9, 0(R0)
	MUL	R6, R3, R9
	UMULH	R6, R3, R10
	ADDS	R4, R9, R9
	MOVD	R9, 8(R0)
	ADC	$0, R10, R4
	MOVD	R4, ret+32(FP)
	RET

mul1_n3:
	LDP	0(R1), (R5, R6)
	MOVD	16(R1), R7
	MUL	R5, R3, R9
	UMULH	R5, R3, R4
	MOVD	R9, 0(R0)
	MUL	R6, R3, R9
	UMULH	R6, R3, R10
	ADDS	R4, R9, R9
	MOVD	R9, 8(R0)
	ADC	$0, R10, R4
	MUL	R7, R3, R9
	UMULH	R7, R3, R10
	ADDS	R4, R9, R9
	MOVD	R9, 16(R0)
	ADC	$0, R10, R4
	MOVD	R4, ret+32(FP)
	RET

mul1_n4:
	LDP	0(R1), (R5, R6)
	LDP	16(R1), (R7, R8)
	MUL	R5, R3, R9
	UMULH	R5, R3, R4
	MOVD	R9, 0(R0)
	MUL	R6, R3, R9
	UMULH	R6, R3, R10
	ADDS	R4, R9, R9
	MOVD	R9, 8(R0)
	ADC	$0, R10, R4
	MUL	R7, R3, R9
	UMULH	R7, R3, R10
	ADDS	R4, R9, R9
	MOVD	R9, 16(R0)
	ADC	$0, R10, R4
	MUL	R8, R3, R9
	UMULH	R8, R3, R10
	ADDS	R4, R9, R9
	MOVD	R9, 24(R0)
	ADC	$0, R10, R4
	MOVD	R4, ret+32(FP)
	RET

mul1_n8:
	// Process first 4 limbs
	LDP	0(R1), (R5, R6)
	LDP	16(R1), (R7, R8)
	MUL	R5, R3, R9
	UMULH	R5, R3, R4
	MOVD	R9, 0(R0)
	MUL	R6, R3, R9
	UMULH	R6, R3, R10
	ADDS	R4, R9, R9
	MOVD	R9, 8(R0)
	ADC	$0, R10, R4
	MUL	R7, R3, R9
	UMULH	R7, R3, R10
	ADDS	R4, R9, R9
	MOVD	R9, 16(R0)
	ADC	$0, R10, R4
	MUL	R8, R3, R9
	UMULH	R8, R3, R10
	ADDS	R4, R9, R9
	MOVD	R9, 24(R0)
	ADC	$0, R10, R4
	// Process next 4 limbs
	LDP	32(R1), (R5, R6)
	LDP	48(R1), (R7, R8)
	MUL	R5, R3, R9
	UMULH	R5, R3, R10
	ADDS	R4, R9, R9
	MOVD	R9, 32(R0)
	ADC	$0, R10, R4
	MUL	R6, R3, R9
	UMULH	R6, R3, R10
	ADDS	R4, R9, R9
	MOVD	R9, 40(R0)
	ADC	$0, R10, R4
	MUL	R7, R3, R9
	UMULH	R7, R3, R10
	ADDS	R4, R9, R9
	MOVD	R9, 48(R0)
	ADC	$0, R10, R4
	MUL	R8, R3, R9
	UMULH	R8, R3, R10
	ADDS	R4, R9, R9
	MOVD	R9, 56(R0)
	ADC	$0, R10, R4
	MOVD	R4, ret+32(FP)
	RET

// mpn_addmul_1 multiplies and accumulates: dst += src * multiplier
// Returns carry (high limb)
//
// func mpnAddMul1(dst, src *uint64, n int, multiplier uint64) uint64
TEXT ·mpnAddMul1(SB), NOSPLIT, $0-32
	MOVD	dst+0(FP), R0      // R0 = destination pointer
	MOVD	src+8(FP), R1      // R1 = source pointer
	MOVD	n+16(FP), R2       // R2 = number of limbs
	MOVD	multiplier+24(FP), R3 // R3 = multiplier

	// Handle zero case
	CBZ	R2, addmul1_done_zero

	// Clear high part accumulator
	MOVD	$0, R4  // R4 = high part accumulator

	// Check if we have at least 4 limbs
	CMP	$4, R2
	BLT	addmul1_remainder

	// Main loop: process 4 limbs at a time
addmul1_loop_unrolled:
	// Load 4 source limbs using load-pair
	LDP	0(R1), (R5, R6)   // R5 = src[0], R6 = src[1]
	LDP	16(R1), (R7, R8)  // R7 = src[2], R8 = src[3]

	// Load 4 destination limbs
	LDP	0(R0), (R9, R10)  // R9 = dst[0], R10 = dst[1]
	LDP	16(R0), (R11, R12) // R11 = dst[2], R12 = dst[3]

	// Multiply-accumulate limb 0: dst[0] += src[0] * multiplier + high
	MUL	R5, R3, R13       // R13 = low 64 bits of R5 * R3
	UMULH	R5, R3, R14       // R14 = high 64 bits of R5 * R3
	ADDS	R4, R13, R13      // Add previous high part, set flags
	ADC	$0, R14, R14      // Add carry to high part
	ADDS	R9, R13, R13      // Add destination value, set flags
	ADC	$0, R14, R14      // Add carry to high part
	MOVD	R13, 0(R0)        // Store result
	MOVD	R14, R4           // Save high for next iteration

	// Multiply-accumulate limb 1
	MUL	R6, R3, R13
	UMULH	R6, R3, R14
	ADDS	R4, R13, R13
	ADC	$0, R14, R14
	ADDS	R10, R13, R13
	ADC	$0, R14, R14
	MOVD	R13, 8(R0)
	MOVD	R14, R4

	// Multiply-accumulate limb 2
	MUL	R7, R3, R13
	UMULH	R7, R3, R14
	ADDS	R4, R13, R13
	ADC	$0, R14, R14
	ADDS	R11, R13, R13
	ADC	$0, R14, R14
	MOVD	R13, 16(R0)
	MOVD	R14, R4

	// Multiply-accumulate limb 3
	MUL	R8, R3, R13
	UMULH	R8, R3, R14
	ADDS	R4, R13, R13
	ADC	$0, R14, R14
	ADDS	R12, R13, R13
	ADC	$0, R14, R14
	MOVD	R13, 24(R0)
	MOVD	R14, R4

	// Advance pointers by 4 limbs (32 bytes)
	ADD	$32, R1, R1
	ADD	$32, R0, R0

	// Decrement counter and loop if >= 4 remain
	SUBS	$4, R2, R2
	CMP	$4, R2
	BGE	addmul1_loop_unrolled

addmul1_remainder:
	// Process remaining 0-3 limbs one at a time
	CBZ	R2, addmul1_done

addmul1_remainder_loop:
	MOVD	0(R1), R5         // Load source limb
	MOVD	0(R0), R9          // Load destination limb
	MUL	R5, R3, R13        // R13 = low 64 bits
	UMULH	R5, R3, R14        // R14 = high 64 bits
	ADDS	R4, R13, R13       // Add previous high part
	ADC	$0, R14, R14       // Add carry
	ADDS	R9, R13, R13       // Add destination value
	ADC	$0, R14, R14       // Add carry
	MOVD	R13, 0(R0)         // Store result
	MOVD	R14, R4            // Save high for next iteration

	ADD	$8, R1, R1
	ADD	$8, R0, R0

	SUBS	$1, R2, R2
	BNE	addmul1_remainder_loop

addmul1_done:
	MOVD	R4, ret+32(FP)
	RET

addmul1_done_zero:
	MOVD	$0, R4
	MOVD	R4, ret+32(FP)
	RET

