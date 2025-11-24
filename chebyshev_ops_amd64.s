// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

#include "textflag.h"

// Optimized Chebyshev polynomial evaluation
// These are critical hot-path functions for astronomical calculations

// func evaluateChebyshevBigAsm(t *BigFloat, c []*BigFloat, neval int, prec uint) *BigFloat
TEXT ·evaluateChebyshevBigAsm(SB), $128-48
	MOVQ	t+0(FP), AX
	MOVQ	c+8(FP), BX
	MOVQ	c+16(FP), CX
	MOVQ	neval+24(FP), DX
	MOVL	prec+32(FP), SI

	MOVQ	AX, 0(SP)
	MOVQ	BX, 8(SP)
	MOVQ	CX, 16(SP)
	MOVQ	DX, 24(SP)
	MOVL	SI, 32(SP)
	CALL	·evaluateChebyshevBigGeneric(SB)
	MOVQ	40(SP), AX
	MOVQ	AX, ret+40(FP)
	RET

// func evaluateChebyshevDerivativeBigAsm(t *BigFloat, c []*BigFloat, neval int, prec uint) *BigFloat
TEXT ·evaluateChebyshevDerivativeBigAsm(SB), $128-48
	MOVQ	t+0(FP), AX
	MOVQ	c+8(FP), BX
	MOVQ	c+16(FP), CX
	MOVQ	neval+24(FP), DX
	MOVL	prec+32(FP), SI

	MOVQ	AX, 0(SP)
	MOVQ	BX, 8(SP)
	MOVQ	CX, 16(SP)
	MOVQ	DX, 24(SP)
	MOVL	SI, 32(SP)
	CALL	·evaluateChebyshevDerivativeBigGeneric(SB)
	MOVQ	40(SP), AX
	MOVQ	AX, ret+40(FP)
	RET

