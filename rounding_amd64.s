// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

#include "textflag.h"

// Rounding mode implementations for AMD64
// These implement IEEE 754-2008 rounding modes

// func roundToNearestEvenAsm(x *BigFloat, prec uint) *BigFloat
// Round to nearest, ties to even (default IEEE 754 mode)
TEXT ·roundToNearestEvenAsm(SB), NOSPLIT, $128-24
	MOVQ	x+0(FP), AX
	MOVL	prec+8(FP), BX

	MOVQ	AX, 0(SP)
	MOVL	BX, 8(SP)
	CALL	·roundToNearestEvenGeneric(SB)
	MOVQ	16(SP), AX
	MOVQ	AX, ret+16(FP)
	RET

// func roundTowardZeroAsm(x *BigFloat, prec uint) *BigFloat
// Round toward zero (truncate)
TEXT ·roundTowardZeroAsm(SB), NOSPLIT, $128-24
	MOVQ	x+0(FP), AX
	MOVL	prec+8(FP), BX

	MOVQ	AX, 0(SP)
	MOVL	BX, 8(SP)
	CALL	·roundTowardZeroGeneric(SB)
	MOVQ	16(SP), AX
	MOVQ	AX, ret+16(FP)
	RET

// func roundTowardInfAsm(x *BigFloat, prec uint) *BigFloat
// Round toward +∞ (ceiling)
TEXT ·roundTowardInfAsm(SB), NOSPLIT, $128-24
	MOVQ	x+0(FP), AX
	MOVL	prec+8(FP), BX

	MOVQ	AX, 0(SP)
	MOVL	BX, 8(SP)
	CALL	·roundTowardInfGeneric(SB)
	MOVQ	16(SP), AX
	MOVQ	AX, ret+16(FP)
	RET

// func roundTowardNegInfAsm(x *BigFloat, prec uint) *BigFloat
// Round toward -∞ (floor)
TEXT ·roundTowardNegInfAsm(SB), NOSPLIT, $128-24
	MOVQ	x+0(FP), AX
	MOVL	prec+8(FP), BX

	MOVQ	AX, 0(SP)
	MOVL	BX, 8(SP)
	CALL	·roundTowardNegInfGeneric(SB)
	MOVQ	16(SP), AX
	MOVQ	AX, ret+16(FP)
	RET

