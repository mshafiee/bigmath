// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

#include "textflag.h"

// Hyperbolic functions for ARM64

// func bigSinhAsm(x *BigFloat, prec uint) *BigFloat
TEXT ·bigSinhAsm(SB), $128-24
	MOVD	x+0(FP), R0
	MOVW	prec+8(FP), R1

	MOVD	R0, 0(RSP)
	MOVW	R1, 8(RSP)
	CALL	·bigSinhGeneric(SB)
	MOVD	16(RSP), R0
	MOVD	R0, ret+16(FP)
	RET

// func bigCoshAsm(x *BigFloat, prec uint) *BigFloat
TEXT ·bigCoshAsm(SB), $128-24
	MOVD	x+0(FP), R0
	MOVW	prec+8(FP), R1

	MOVD	R0, 0(RSP)
	MOVW	R1, 8(RSP)
	CALL	·bigCoshGeneric(SB)
	MOVD	16(RSP), R0
	MOVD	R0, ret+16(FP)
	RET

// func bigTanhAsm(x *BigFloat, prec uint) *BigFloat
TEXT ·bigTanhAsm(SB), $128-24
	MOVD	x+0(FP), R0
	MOVW	prec+8(FP), R1

	MOVD	R0, 0(RSP)
	MOVW	R1, 8(RSP)
	CALL	·bigTanhGeneric(SB)
	MOVD	16(RSP), R0
	MOVD	R0, ret+16(FP)
	RET

// func bigAsinhAsm(x *BigFloat, prec uint) *BigFloat
TEXT ·bigAsinhAsm(SB), $128-24
	MOVD	x+0(FP), R0
	MOVW	prec+8(FP), R1

	MOVD	R0, 0(RSP)
	MOVW	R1, 8(RSP)
	CALL	·bigAsinhGeneric(SB)
	MOVD	16(RSP), R0
	MOVD	R0, ret+16(FP)
	RET

// func bigAcoshAsm(x *BigFloat, prec uint) *BigFloat
TEXT ·bigAcoshAsm(SB), $128-24
	MOVD	x+0(FP), R0
	MOVW	prec+8(FP), R1

	MOVD	R0, 0(RSP)
	MOVW	R1, 8(RSP)
	CALL	·bigAcoshGeneric(SB)
	MOVD	16(RSP), R0
	MOVD	R0, ret+16(FP)
	RET

// func bigAtanhAsm(x *BigFloat, prec uint) *BigFloat
TEXT ·bigAtanhAsm(SB), $128-24
	MOVD	x+0(FP), R0
	MOVW	prec+8(FP), R1

	MOVD	R0, 0(RSP)
	MOVW	R1, 8(RSP)
	CALL	·bigAtanhGeneric(SB)
	MOVD	16(RSP), R0
	MOVD	R0, ret+16(FP)
	RET

