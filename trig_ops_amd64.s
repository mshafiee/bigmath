// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

#include "textflag.h"

// Additional trigonometric functions in assembly

// func bigTanAsm(x *BigFloat, prec uint) *BigFloat
TEXT ·bigTanAsm(SB), $128-24
	MOVQ	x+0(FP), AX
	MOVL	prec+8(FP), BX

	MOVQ	AX, 0(SP)
	MOVL	BX, 8(SP)
	CALL	·bigTanGeneric(SB)
	MOVQ	16(SP), AX
	MOVQ	AX, ret+16(FP)
	RET

// func bigAtanAsm(x *BigFloat, prec uint) *BigFloat
TEXT ·bigAtanAsm(SB), $128-24
	MOVQ	x+0(FP), AX
	MOVL	prec+8(FP), BX

	MOVQ	AX, 0(SP)
	MOVL	BX, 8(SP)
	CALL	·bigAtanGeneric(SB)
	MOVQ	16(SP), AX
	MOVQ	AX, ret+16(FP)
	RET

// func bigAsinAsm(x *BigFloat, prec uint) *BigFloat
TEXT ·bigAsinAsm(SB), $128-24
	MOVQ	x+0(FP), AX
	MOVL	prec+8(FP), BX

	MOVQ	AX, 0(SP)
	MOVL	BX, 8(SP)
	CALL	·bigAsinGeneric(SB)
	MOVQ	16(SP), AX
	MOVQ	AX, ret+16(FP)
	RET

// func bigAcosAsm(x *BigFloat, prec uint) *BigFloat
TEXT ·bigAcosAsm(SB), $128-24
	MOVQ	x+0(FP), AX
	MOVL	prec+8(FP), BX

	MOVQ	AX, 0(SP)
	MOVL	BX, 8(SP)
	CALL	·bigAcosGeneric(SB)
	MOVQ	16(SP), AX
	MOVQ	AX, ret+16(FP)
	RET

// func bigAtan2Asm(y, x *BigFloat, prec uint) *BigFloat
TEXT ·bigAtan2Asm(SB), $128-32
	MOVQ	y+0(FP), AX
	MOVQ	x+8(FP), BX
	MOVL	prec+16(FP), CX

	MOVQ	AX, 0(SP)
	MOVQ	BX, 8(SP)
	MOVL	CX, 16(SP)
	CALL	·bigAtan2Generic(SB)
	MOVQ	24(SP), AX
	MOVQ	AX, ret+24(FP)
	RET

