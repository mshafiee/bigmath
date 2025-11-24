// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

#include "textflag.h"

// Hyperbolic functions in assembly

// func bigSinhAsm(x *BigFloat, prec uint) *BigFloat
TEXT ·bigSinhAsm(SB), NOSPLIT, $128-24
	MOVQ	x+0(FP), AX
	MOVL	prec+8(FP), BX

	MOVQ	AX, 0(SP)
	MOVL	BX, 8(SP)
	CALL	·bigSinhGeneric(SB)
	MOVQ	16(SP), AX
	MOVQ	AX, ret+16(FP)
	RET

// func bigCoshAsm(x *BigFloat, prec uint) *BigFloat
TEXT ·bigCoshAsm(SB), NOSPLIT, $128-24
	MOVQ	x+0(FP), AX
	MOVL	prec+8(FP), BX

	MOVQ	AX, 0(SP)
	MOVL	BX, 8(SP)
	CALL	·bigCoshGeneric(SB)
	MOVQ	16(SP), AX
	MOVQ	AX, ret+16(FP)
	RET

// func bigTanhAsm(x *BigFloat, prec uint) *BigFloat
TEXT ·bigTanhAsm(SB), NOSPLIT, $128-24
	MOVQ	x+0(FP), AX
	MOVL	prec+8(FP), BX

	MOVQ	AX, 0(SP)
	MOVL	BX, 8(SP)
	CALL	·bigTanhGeneric(SB)
	MOVQ	16(SP), AX
	MOVQ	AX, ret+16(FP)
	RET

// func bigAsinhAsm(x *BigFloat, prec uint) *BigFloat
TEXT ·bigAsinhAsm(SB), NOSPLIT, $128-24
	MOVQ	x+0(FP), AX
	MOVL	prec+8(FP), BX

	MOVQ	AX, 0(SP)
	MOVL	BX, 8(SP)
	CALL	·bigAsinhGeneric(SB)
	MOVQ	16(SP), AX
	MOVQ	AX, ret+16(FP)
	RET

// func bigAcoshAsm(x *BigFloat, prec uint) *BigFloat
TEXT ·bigAcoshAsm(SB), NOSPLIT, $128-24
	MOVQ	x+0(FP), AX
	MOVL	prec+8(FP), BX

	MOVQ	AX, 0(SP)
	MOVL	BX, 8(SP)
	CALL	·bigAcoshGeneric(SB)
	MOVQ	16(SP), AX
	MOVQ	AX, ret+16(FP)
	RET

// func bigAtanhAsm(x *BigFloat, prec uint) *BigFloat
TEXT ·bigAtanhAsm(SB), NOSPLIT, $128-24
	MOVQ	x+0(FP), AX
	MOVL	prec+8(FP), BX

	MOVQ	AX, 0(SP)
	MOVL	BX, 8(SP)
	CALL	·bigAtanhGeneric(SB)
	MOVQ	16(SP), AX
	MOVQ	AX, ret+16(FP)
	RET

