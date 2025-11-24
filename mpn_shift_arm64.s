#include "textflag.h"

// Multi-precision integer shift operations for ARM64
// mpn_lshift left shifts n-limb number by count bits
// Returns bits shifted out
//
// func mpnLShift(dst, src *uint64, n int, count uint) uint64
TEXT ·mpnLShift(SB), NOSPLIT, $0-32
	MOVD	dst+0(FP), R0      // R0 = destination pointer
	MOVD	src+8(FP), R1      // R1 = source pointer
	MOVD	n+16(FP), R2       // R2 = number of limbs
	MOVD	count+24(FP), R3   // R3 = shift count

	// Handle zero case
	CBZ	R2, lshift_done_zero

	// Handle zero shift
	CBZ	R3, lshift_copy

	// Compute complement shift
	MOVD	$64, R4
	SUB	R3, R4, R4  // R4 = 64 - count

	// Process from high to low (right to left)
	SUB	$1, R2, R5  // R5 = n - 1

	// Load first limb
	MOVD	R5, R6
	LSL	$3, R6      // R6 = (n-1) * 8
	MOVD	(R1)(R6), R7  // R7 = src[n-1]
	MOVD	R7, R8
	LSL	R3, R8, R8   // R8 = high bits shifted out
	MOVD	R8, ret+32(FP)  // Save return value

	// Shift first limb
	MOVD	R7, (R0)(R6)

	SUBS	$1, R5, R5
	BLT	lshift_done

lshift_loop:
	// Load current limb
	MOVD	R5, R6
	LSL	$3, R6
	MOVD	(R1)(R6), R7
	// Load next limb
	ADD	$8, R6, R8
	MOVD	(R1)(R8), R9

	// Shift: (R7 << R3) | (R9 >> (64-R3))
	LSL	R3, R7, R10
	LSR	R4, R9, R11
	ORR	R11, R10, R10

	MOVD	R10, (R0)(R6)

	SUBS	$1, R5, R5
	BGE	lshift_loop

lshift_done:
	MOVD	ret+32(FP), R0
	RET

lshift_copy:
	// Zero shift - just copy
	MOVD	R2, R3
	LSL	$3, R3      // R3 = n * 8
	MOVD	R1, R4
	MOVD	R0, R5
	// Simple copy loop
	CBZ	R3, lshift_copy_done
lshift_copy_loop:
	MOVD	(R4), R6
	MOVD	R6, (R5)
	ADD	$8, R4, R4
	ADD	$8, R5, R5
	SUBS	$8, R3, R3
	BNE	lshift_copy_loop
lshift_copy_done:
	MOVD	$0, R0
	MOVD	R0, ret+32(FP)
	RET

lshift_done_zero:
	MOVD	$0, R0
	MOVD	R0, ret+32(FP)
	RET

// mpn_rshift right shifts n-limb number by count bits
// Returns bits shifted out
//
// func mpnRShift(dst, src *uint64, n int, count uint) uint64
TEXT ·mpnRShift(SB), NOSPLIT, $0-32
	MOVD	dst+0(FP), R0      // R0 = destination pointer
	MOVD	src+8(FP), R1      // R1 = source pointer
	MOVD	n+16(FP), R2       // R2 = number of limbs
	MOVD	count+24(FP), R3   // R3 = shift count

	// Handle zero case
	CBZ	R2, rshift_done_zero

	// Handle zero shift
	CBZ	R3, rshift_copy

	// Compute complement shift
	MOVD	$64, R4
	SUB	R3, R4, R4  // R4 = 64 - count

	// Process from low to high (left to right)
	MOVD	$0, R5  // R5 = index

	// Load first limb
	MOVD	(R1), R7  // R7 = src[0]
	MOVD	R7, R8
	LSR	R3, R8, R8   // R8 = low bits shifted out
	MOVD	R8, ret+32(FP)  // Save return value

	// Shift first limb
	MOVD	R7, (R0)

	ADD	$1, R5, R5
	CMP	R5, R2
	BGE	rshift_done

rshift_loop:
	// Load current limb
	MOVD	R5, R6
	LSL	$3, R6
	MOVD	(R1)(R6), R7
	// Load previous limb
	SUB	$8, R6, R8
	MOVD	(R1)(R8), R9

	// Shift: (R9 >> R3) | (R7 << (64-R3))
	LSR	R3, R9, R10
	LSL	R4, R7, R11
	ORR	R11, R10, R10

	MOVD	R5, R6
	LSL	$3, R6
	MOVD	R10, (R0)(R6)

	ADD	$1, R5, R5
	CMP	R5, R2
	BLT	rshift_loop

rshift_done:
	MOVD	ret+32(FP), R0
	RET

rshift_copy:
	// Zero shift - just copy
	MOVD	R2, R3
	LSL	$3, R3      // R3 = n * 8
	MOVD	R1, R4
	MOVD	R0, R5
	// Simple copy loop
	CBZ	R3, rshift_copy_done
rshift_copy_loop:
	MOVD	(R4), R6
	MOVD	R6, (R5)
	ADD	$8, R4, R4
	ADD	$8, R5, R5
	SUBS	$8, R3, R3
	BNE	rshift_copy_loop
rshift_copy_done:
	MOVD	$0, R0
	MOVD	R0, ret+32(FP)
	RET

rshift_done_zero:
	MOVD	$0, R0
	MOVD	R0, ret+32(FP)
	RET

