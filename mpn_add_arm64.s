// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

#include "textflag.h"

// Multi-precision integer addition for ARM64 (optimized with loop unrolling)
// mpn_add_n adds two n-limb numbers: dst = src1 + src2
// Returns carry (0 or 1)
//
// func mpnAddN(dst, src1, src2 *uint64, n int) uint64
TEXT ·mpnAddN(SB), NOSPLIT, $0-32
	MOVD	dst+0(FP), R0   // R0 = destination pointer
	MOVD	src1+8(FP), R1  // R1 = source 1 pointer
	MOVD	src2+16(FP), R2 // R2 = source 2 pointer
	MOVD	n+24(FP), R3    // R3 = number of limbs

	// Handle zero case
	CBZ	R3, done_zero

	// Clear carry by comparing register with itself (sets flags such that carry = 0)
	CMP	R0, R0

	// Fast paths for small n (most common cases: 1-4 limbs)
	CMP	$1, R3
	BEQ	n1
	CMP	$2, R3
	BEQ	n2
	CMP	$3, R3
	BEQ	n3
	CMP	$4, R3
	BEQ	n4

	// Check if we have at least 4 limbs for loop
	BLT	remainder

	// Main loop: process 4 limbs at a time
loop_unrolled:
	// Load 4 limbs from src1 using load-pair instructions
	LDP	0(R1), (R5, R6)  // R5 = src1[0], R6 = src1[1]
	LDP	16(R1), (R7, R8) // R7 = src1[2], R8 = src1[3]

	// Load 4 limbs from src2
	LDP	0(R2), (R9, R10)  // R9 = src2[0], R10 = src2[1]
	LDP	16(R2), (R11, R12) // R11 = src2[2], R12 = src2[3]

	// Add with carry propagation: ADDS sets flags, ADCS uses carry
	ADDS	R9, R5, R5    // R5 = R5 + R9, set flags
	ADCS	R10, R6, R6   // R6 = R6 + R10 + carry, propagate carry
	ADCS	R11, R7, R7   // R7 = R7 + R11 + carry
	ADCS	R12, R8, R8   // R8 = R8 + R12 + carry

	// Store results using store-pair instructions
	STP	(R5, R6), 0(R0)  // dst[0] = R5, dst[1] = R6
	STP	(R7, R8), 16(R0) // dst[2] = R7, dst[3] = R8

	// Advance pointers by 4 limbs (32 bytes)
	ADD	$32, R1, R1
	ADD	$32, R2, R2
	ADD	$32, R0, R0

	// Decrement counter and loop if >= 4 remain
	SUBS	$4, R3, R3
	CMP	$4, R3
	BGE	loop_unrolled

remainder:
	// Process remaining 0-3 limbs
	CBZ	R3, done

remainder_loop:
	// Process 2 limbs at a time if possible
	CMP	$2, R3
	BLT	remainder_single

	LDP	0(R1), (R5, R6)
	LDP	0(R2), (R9, R10)

	ADDS	R9, R5, R5
	ADCS	R10, R6, R6

	STP	(R5, R6), 0(R0)

	ADD	$16, R1, R1
	ADD	$16, R2, R2
	ADD	$16, R0, R0

	SUBS	$2, R3, R3
	BGT	remainder_loop

remainder_single:
	// Handle last limb if odd
	CBZ	R3, done
	MOVD	0(R1), R5
	MOVD	0(R2), R9
	ADDS	R9, R5, R5
	MOVD	R5, 0(R0)

done:
	// Return carry flag as 0 or 1
	CSET	HS, R4         // R4 = 1 if carry set (HS = higher or same), else 0
	MOVD	R4, ret+32(FP)
	RET

done_zero:
	MOVD	$0, R4
	MOVD	R4, ret+32(FP)
	RET

// Fast paths for small n (straight-line code, no loops)
n1:
	MOVD	0(R1), R5
	MOVD	0(R2), R9
	ADDS	R9, R5, R5
	MOVD	R5, 0(R0)
	CSET	HS, R4
	MOVD	R4, ret+32(FP)
	RET

n2:
	LDP	0(R1), (R5, R6)
	LDP	0(R2), (R9, R10)
	ADDS	R9, R5, R5
	ADCS	R10, R6, R6
	STP	(R5, R6), 0(R0)
	CSET	HS, R4
	MOVD	R4, ret+32(FP)
	RET

n3:
	LDP	0(R1), (R5, R6)
	MOVD	16(R1), R7
	LDP	0(R2), (R9, R10)
	MOVD	16(R2), R11
	ADDS	R9, R5, R5
	ADCS	R10, R6, R6
	ADCS	R11, R7, R7
	STP	(R5, R6), 0(R0)
	MOVD	R7, 16(R0)
	CSET	HS, R4
	MOVD	R4, ret+32(FP)
	RET

n4:
	LDP	0(R1), (R5, R6)
	LDP	16(R1), (R7, R8)
	LDP	0(R2), (R9, R10)
	LDP	16(R2), (R11, R12)
	ADDS	R9, R5, R5
	ADCS	R10, R6, R6
	ADCS	R11, R7, R7
	ADCS	R12, R8, R8
	STP	(R5, R6), 0(R0)
	STP	(R7, R8), 16(R0)
	CSET	HS, R4
	MOVD	R4, ret+32(FP)
	RET

// mpn_sub_n subtracts two n-limb numbers: dst = src1 - src2
// Returns borrow (0 or 1)
//
// func mpnSubN(dst, src1, src2 *uint64, n int) uint64
TEXT ·mpnSubN(SB), NOSPLIT, $0-32
	MOVD	dst+0(FP), R0   // R0 = destination pointer
	MOVD	src1+8(FP), R1  // R1 = source 1 pointer
	MOVD	src2+16(FP), R2 // R2 = source 2 pointer
	MOVD	n+24(FP), R3    // R3 = number of limbs

	// Handle zero case
	CBZ	R3, sub_done_zero

	// Clear borrow flag (set carry = 1 means no initial borrow)
	CMP	R0, R0  // This sets flags such that carry = 0 (no borrow initially)

	// Fast paths for small n (most common cases: 1-4 limbs)
	CMP	$1, R3
	BEQ	sub_n1
	CMP	$2, R3
	BEQ	sub_n2
	CMP	$3, R3
	BEQ	sub_n3
	CMP	$4, R3
	BEQ	sub_n4

	// Check if we have at least 4 limbs for loop
	BLT	sub_remainder

	// Main loop: process 4 limbs at a time
sub_loop_unrolled:
	// Load 4 limbs from src1 using load-pair instructions
	LDP	0(R1), (R5, R6)  // R5 = src1[0], R6 = src1[1]
	LDP	16(R1), (R7, R8) // R7 = src1[2], R8 = src1[3]

	// Load 4 limbs from src2
	LDP	0(R2), (R9, R10)  // R9 = src2[0], R10 = src2[1]
	LDP	16(R2), (R11, R12) // R11 = src2[2], R12 = src2[3]

	// Subtract with borrow propagation: SUBS sets flags, SBCS uses borrow
	SUBS	R9, R5, R5    // R5 = R5 - R9, set flags
	SBCS	R10, R6, R6   // R6 = R6 - R10 - borrow, propagate borrow
	SBCS	R11, R7, R7   // R7 = R7 - R11 - borrow
	SBCS	R12, R8, R8   // R8 = R8 - R12 - borrow

	// Store results using store-pair instructions
	STP	(R5, R6), 0(R0)  // dst[0] = R5, dst[1] = R6
	STP	(R7, R8), 16(R0) // dst[2] = R7, dst[3] = R8

	// Advance pointers by 4 limbs (32 bytes)
	ADD	$32, R1, R1
	ADD	$32, R2, R2
	ADD	$32, R0, R0

	// Decrement counter and loop if >= 4 remain
	SUBS	$4, R3, R3
	CMP	$4, R3
	BGE	sub_loop_unrolled

sub_remainder:
	// Process remaining 0-3 limbs
	CBZ	R3, sub_done

sub_remainder_loop:
	// Process 2 limbs at a time if possible
	CMP	$2, R3
	BLT	sub_remainder_single

	LDP	0(R1), (R5, R6)
	LDP	0(R2), (R9, R10)

	SUBS	R9, R5, R5
	SBCS	R10, R6, R6

	STP	(R5, R6), 0(R0)

	ADD	$16, R1, R1
	ADD	$16, R2, R2
	ADD	$16, R0, R0

	SUBS	$2, R3, R3
	BGT	sub_remainder_loop

sub_remainder_single:
	// Handle last limb if odd
	CBZ	R3, sub_done
	MOVD	0(R1), R5
	MOVD	0(R2), R9
	SUBS	R9, R5, R5
	MOVD	R5, 0(R0)

sub_done:
	// Return borrow flag as 0 or 1 (carry clear = borrow)
	CSET	LO, R4         // R4 = 1 if borrow (LO = lower = borrow), else 0
	MOVD	R4, ret+32(FP)
	RET

sub_done_zero:
	MOVD	$0, R4
	MOVD	R4, ret+32(FP)
	RET

// Fast paths for small n subtraction (straight-line code, no loops)
sub_n1:
	MOVD	0(R1), R5
	MOVD	0(R2), R9
	SUBS	R9, R5, R5
	MOVD	R5, 0(R0)
	CSET	LO, R4
	MOVD	R4, ret+32(FP)
	RET

sub_n2:
	LDP	0(R1), (R5, R6)
	LDP	0(R2), (R9, R10)
	SUBS	R9, R5, R5
	SBCS	R10, R6, R6
	STP	(R5, R6), 0(R0)
	CSET	LO, R4
	MOVD	R4, ret+32(FP)
	RET

sub_n3:
	LDP	0(R1), (R5, R6)
	MOVD	16(R1), R7
	LDP	0(R2), (R9, R10)
	MOVD	16(R2), R11
	SUBS	R9, R5, R5
	SBCS	R10, R6, R6
	SBCS	R11, R7, R7
	STP	(R5, R6), 0(R0)
	MOVD	R7, 16(R0)
	CSET	LO, R4
	MOVD	R4, ret+32(FP)
	RET

sub_n4:
	LDP	0(R1), (R5, R6)
	LDP	16(R1), (R7, R8)
	LDP	0(R2), (R9, R10)
	LDP	16(R2), (R11, R12)
	SUBS	R9, R5, R5
	SBCS	R10, R6, R6
	SBCS	R11, R7, R7
	SBCS	R12, R8, R8
	STP	(R5, R6), 0(R0)
	STP	(R7, R8), 16(R0)
	CSET	LO, R4
	MOVD	R4, ret+32(FP)
	RET

