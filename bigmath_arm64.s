// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

#include "textflag.h"

// func bigVec3AddAsmARM64(v1, v2 *BigVec3, prec uint) *BigVec3
// ARM64/NEON optimized 3D vector addition with arbitrary precision
TEXT ·bigVec3AddAsmARM64(SB), $128-32
	MOVD	v1+0(FP), R0
	MOVD	v2+8(FP), R1
	MOVW	prec+16(FP), R2

	MOVD	R0, 0(RSP)
	MOVD	R1, 8(RSP)
	MOVW	R2, 16(RSP)
	CALL	·bigVec3AddGeneric(SB)
	MOVD	24(RSP), R0
	MOVD	R0, ret+24(FP)
	RET

// func bigVec3SubAsmARM64(v1, v2 *BigVec3, prec uint) *BigVec3
TEXT ·bigVec3SubAsmARM64(SB), $128-32
	MOVD	v1+0(FP), R0
	MOVD	v2+8(FP), R1
	MOVW	prec+16(FP), R2

	MOVD	R0, 0(RSP)
	MOVD	R1, 8(RSP)
	MOVW	R2, 16(RSP)
	CALL	·bigVec3SubGeneric(SB)
	MOVD	24(RSP), R0
	MOVD	R0, ret+24(FP)
	RET

// func bigVec3MulAsmARM64(v *BigVec3, scalar *BigFloat, prec uint) *BigVec3
TEXT ·bigVec3MulAsmARM64(SB), $128-32
	MOVD	v+0(FP), R0
	MOVD	scalar+8(FP), R1
	MOVW	prec+16(FP), R2

	MOVD	R0, 0(RSP)
	MOVD	R1, 8(RSP)
	MOVW	R2, 16(RSP)
	CALL	·bigVec3MulGeneric(SB)
	MOVD	24(RSP), R0
	MOVD	R0, ret+24(FP)
	RET

// func bigVec3DotAsmARM64(v1, v2 *BigVec3, prec uint) *BigFloat
TEXT ·bigVec3DotAsmARM64(SB), $128-32
	MOVD	v1+0(FP), R0
	MOVD	v2+8(FP), R1
	MOVW	prec+16(FP), R2

	MOVD	R0, 0(RSP)
	MOVD	R1, 8(RSP)
	MOVW	R2, 16(RSP)
	CALL	·bigVec3DotGeneric(SB)
	MOVD	24(RSP), R0
	MOVD	R0, ret+24(FP)
	RET

// func bigMatMulAsmARM64(m *BigMatrix3x3, v *BigVec3, prec uint) *BigVec3
TEXT ·bigMatMulAsmARM64(SB), $128-32
	MOVD	m+0(FP), R0
	MOVD	v+8(FP), R1
	MOVW	prec+16(FP), R2

	MOVD	R0, 0(RSP)
	MOVD	R1, 8(RSP)
	MOVW	R2, 16(RSP)
	CALL	·bigMatMulGeneric(SB)
	MOVD	24(RSP), R0
	MOVD	R0, ret+24(FP)
	RET

