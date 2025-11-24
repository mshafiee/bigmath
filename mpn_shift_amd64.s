#include "textflag.h"

// Multi-precision integer shift operations
// mpn_lshift left shifts n-limb number by count bits
// Returns bits shifted out
//
// func mpnLShift(dst, src *uint64, n int, count uint) uint64
TEXT ·mpnLShift(SB), NOSPLIT, $0-32
	MOVQ	dst+0(FP), DI      // DI = destination pointer
	MOVQ	src+8(FP), SI      // SI = source pointer
	MOVQ	n+16(FP), CX       // CX = number of limbs
	MOVQ	count+24(FP), DX   // DX = shift count

	// Handle zero case
	TESTQ	CX, CX
	JZ	lshift_done_zero

	// Handle zero shift
	TESTQ	DX, DX
	JZ	lshift_copy

	// Compute complement shift for SHLD
	MOVQ	$64, AX
	SUBQ	DX, AX  // AX = 64 - count

	// Process from high to low (right to left)
	MOVQ	CX, BX
	DECQ	BX      // BX = n - 1

	// Load first limb
	MOVQ	(SI)(BX*8), R8  // R8 = src[n-1]
	MOVQ	R8, R9
	SHLQ	DX, R9   // R9 = high bits shifted out
	MOVQ	R9, ret+32(FP)  // Save return value

	// Shift first limb
	MOVQ	R8, (DI)(BX*8)

	DECQ	BX
	JL	lshift_done

lshift_loop:
	// Load current limb
	MOVQ	(SI)(BX*8), R8
	// Load next limb for SHLD
	MOVQ	8(SI)(BX*8), R9

	// SHLD shifts R8 left by DX, filling from R9
	// We need to do: (R8 << DX) | (R9 >> (64-DX))
	MOVQ	R8, R10
	SHLQ	DX, R10
	MOVQ	R9, R11
	SHRQ	AX, R11
	ORQ	R11, R10

	MOVQ	R10, (DI)(BX*8)

	DECQ	BX
	JGE	lshift_loop

lshift_done:
	MOVQ	ret+32(FP), AX
	RET

lshift_copy:
	// Zero shift - just copy
	MOVQ	CX, AX
	SHLQ	$3, AX  // AX = n * 8
	MOVQ	SI, BX
	MOVQ	DI, DX
	REP; MOVSQ
	XORQ	AX, AX
	MOVQ	AX, ret+32(FP)
	RET

lshift_done_zero:
	XORQ	AX, AX
	MOVQ	AX, ret+32(FP)
	RET

// mpn_rshift right shifts n-limb number by count bits
// Returns bits shifted out
//
// func mpnRShift(dst, src *uint64, n int, count uint) uint64
TEXT ·mpnRShift(SB), NOSPLIT, $0-32
	MOVQ	dst+0(FP), DI      // DI = destination pointer
	MOVQ	src+8(FP), SI      // SI = source pointer
	MOVQ	n+16(FP), CX       // CX = number of limbs
	MOVQ	count+24(FP), DX   // DX = shift count

	// Handle zero case
	TESTQ	CX, CX
	JZ	rshift_done_zero

	// Handle zero shift
	TESTQ	DX, DX
	JZ	rshift_copy

	// Compute complement shift for SHRD
	MOVQ	$64, AX
	SUBQ	DX, AX  // AX = 64 - count

	// Process from low to high (left to right)
	XORQ	BX, BX  // BX = index

	// Load first limb
	MOVQ	(SI)(BX*8), R8  // R8 = src[0]
	MOVQ	R8, R9
	SHRQ	DX, R9   // R9 = low bits shifted out
	MOVQ	R9, ret+32(FP)  // Save return value

	// Shift first limb
	MOVQ	R8, (DI)(BX*8)

	INCQ	BX
	CMPQ	BX, CX
	JGE	rshift_done

rshift_loop:
	// Load current limb
	MOVQ	(SI)(BX*8), R8
	// Load previous limb for SHRD
	MOVQ	-8(SI)(BX*8), R9

	// SHRD shifts R9 right by DX, filling from R8
	// We need to do: (R9 >> DX) | (R8 << (64-DX))
	MOVQ	R9, R10
	SHRQ	DX, R10
	MOVQ	R8, R11
	SHLQ	AX, R11
	ORQ	R11, R10

	MOVQ	R10, (DI)(BX*8)

	INCQ	BX
	CMPQ	BX, CX
	JL	rshift_loop

rshift_done:
	MOVQ	ret+32(FP), AX
	RET

rshift_copy:
	// Zero shift - just copy
	MOVQ	CX, AX
	SHLQ	$3, AX  // AX = n * 8
	MOVQ	SI, BX
	MOVQ	DI, DX
	REP; MOVSQ
	XORQ	AX, AX
	MOVQ	AX, ret+32(FP)
	RET

rshift_done_zero:
	XORQ	AX, AX
	MOVQ	AX, ret+32(FP)
	RET

