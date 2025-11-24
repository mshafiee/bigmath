// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

#include "textflag.h"

// func bigVec3AddAsm(v1, v2 *BigVec3, prec uint) *BigVec3
// Optimized assembly for 3D vector addition with arbitrary precision
// Currently delegates to generic implementation but structured for future optimization
TEXT ·bigVec3AddAsm(SB), NOSPLIT, $128-32
	// Load arguments
	MOVQ	v1+0(FP), AX      // AX = v1
	MOVQ	v2+8(FP), BX      // BX = v2
	MOVL	prec+16(FP), CX   // CX = prec
	
	// Call generic implementation
	// Future optimization: inline the three Add operations with
	// reduced allocation overhead and better register usage
	MOVQ	AX, 0(SP)
	MOVQ	BX, 8(SP)
	MOVL	CX, 16(SP)
	CALL	·bigVec3AddGeneric(SB)
	MOVQ	24(SP), AX
	MOVQ	AX, ret+24(FP)
	RET

// func bigVec3SubAsm(v1, v2 *BigVec3, prec uint) *BigVec3
TEXT ·bigVec3SubAsm(SB), NOSPLIT, $128-32
	MOVQ	v1+0(FP), AX
	MOVQ	v2+8(FP), BX
	MOVL	prec+16(FP), CX
	
	MOVQ	AX, 0(SP)
	MOVQ	BX, 8(SP)
	MOVL	CX, 16(SP)
	CALL	·bigVec3SubGeneric(SB)
	MOVQ	24(SP), AX
	MOVQ	AX, ret+24(FP)
	RET

// func bigVec3MulAsm(v *BigVec3, scalar *BigFloat, prec uint) *BigVec3
TEXT ·bigVec3MulAsm(SB), NOSPLIT, $128-32
	MOVQ	v+0(FP), AX
	MOVQ	scalar+8(FP), BX
	MOVL	prec+16(FP), CX
	
	MOVQ	AX, 0(SP)
	MOVQ	BX, 8(SP)
	MOVL	CX, 16(SP)
	CALL	·bigVec3MulGeneric(SB)
	MOVQ	24(SP), AX
	MOVQ	AX, ret+24(FP)
	RET

// func bigVec3DotAsm(v1, v2 *BigVec3, prec uint) *BigFloat
TEXT ·bigVec3DotAsm(SB), NOSPLIT, $128-24
	MOVQ	v1+0(FP), AX
	MOVQ	v2+8(FP), BX
	MOVL	prec+16(FP), CX
	
	MOVQ	AX, 0(SP)
	MOVQ	BX, 8(SP)
	MOVL	CX, 16(SP)
	CALL	·bigVec3DotGeneric(SB)
	MOVQ	24(SP), AX
	MOVQ	AX, ret+24(FP)
	RET

// func bigMatMulAsm(m *BigMatrix3x3, v *BigVec3, prec uint) *BigVec3
// Matrix-vector multiplication optimized for AMD64
TEXT ·bigMatMulAsm(SB), NOSPLIT, $128-32
	MOVQ	m+0(FP), AX       // AX = matrix pointer
	MOVQ	v+8(FP), BX       // BX = vector pointer
	MOVL	prec+16(FP), CX   // CX = precision
	
	// Call generic implementation
	// In a full optimization, we would unroll the multiplication loops
	MOVQ	AX, 0(SP)
	MOVQ	BX, 8(SP)
	MOVL	CX, 16(SP)
	CALL	·bigMatMulGeneric(SB)
	MOVQ	24(SP), AX
	MOVQ	AX, ret+24(FP)
	RET

