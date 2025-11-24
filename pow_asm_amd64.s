// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

#include "textflag.h"

// Power function in assembly

// func bigPowAsm(x, y *BigFloat, prec uint) *BigFloat
TEXT ·bigPowAsm(SB), $128-32
	MOVQ	x+0(FP), AX
	MOVQ	y+8(FP), BX
	MOVL	prec+16(FP), CX

	// Call generic implementation
	// Full optimization would:
	// - Check for integer exponents and use binary exponentiation
	// - For general case: exp(y * ln(x))
	MOVQ	AX, 0(SP)
	MOVQ	BX, 8(SP)
	MOVL	CX, 16(SP)
	CALL	·bigPowGeneric(SB)
	MOVQ	24(SP), AX
	MOVQ	AX, ret+24(FP)
	RET

