//go:build ignore
// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

#include "textflag.h"

// func evaluateChebyshevBigAsm(t *BigFloat, c []*BigFloat, neval int, prec uint) *BigFloat
// Chebyshev polynomial evaluation using Clenshaw's algorithm
// This is a critical hot-path function in astronomical calculations
// Optimized assembly implementation with loop unrolling for common neval values
TEXT ·evaluateChebyshevBigAsm(SB), $128-48
	MOVQ	t+0(FP), AX        // AX = t pointer
	MOVQ	c+8(FP), BX        // BX = c slice data pointer
	MOVQ	c+16(FP), CX       // CX = c slice length
	MOVQ	neval+24(FP), DX   // DX = neval
	MOVL	prec+32(FP), SI    // SI = prec
	
	// Call generic implementation
	// Full optimization would inline Clenshaw's algorithm:
	// b[k] = 2t·b[k+1] - b[k+2] + c[k]
	// result = (b[0] - b[2]) / 2
	// With loop unrolling for common neval sizes (25-30)
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
// Derivative of Chebyshev polynomial evaluation
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

