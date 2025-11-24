// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

//go:build ignore
// +build ignore

#include "textflag.h"

// BigFloat arithmetic operations for ARM64
// These work with Go's big.Float type through its API

// func bigfloatAddAsm(a, b *BigFloat, prec uint) *BigFloat
TEXT ·bigfloatAddAsm(SB), $0-32
	MOVD	a+0(FP), R0
	MOVD	b+8(FP), R1
	MOVW	prec+16(FP), R2
	CALL	·bigfloatAddGeneric(SB)
	MOVD	R0, ret+24(FP)
	RET

// func bigfloatSubAsm(a, b *BigFloat, prec uint) *BigFloat
TEXT ·bigfloatSubAsm(SB), $0-32
	MOVD	a+0(FP), R0
	MOVD	b+8(FP), R1
	MOVW	prec+16(FP), R2
	CALL	·bigfloatSubGeneric(SB)
	MOVD	R0, ret+24(FP)
	RET

// func bigfloatMulAsm(a, b *BigFloat, prec uint) *BigFloat
TEXT ·bigfloatMulAsm(SB), $0-32
	MOVD	a+0(FP), R0
	MOVD	b+8(FP), R1
	MOVW	prec+16(FP), R2
	CALL	·bigfloatMulGeneric(SB)
	MOVD	R0, ret+24(FP)
	RET

// func bigfloatDivAsm(a, b *BigFloat, prec uint) *BigFloat
TEXT ·bigfloatDivAsm(SB), $0-32
	MOVD	a+0(FP), R0
	MOVD	b+8(FP), R1
	MOVW	prec+16(FP), R2
	CALL	·bigfloatDivGeneric(SB)
	MOVD	R0, ret+24(FP)
	RET

// func bigfloatSqrtAsm(x *BigFloat, prec uint) *BigFloat
TEXT ·bigfloatSqrtAsm(SB), $0-24
	MOVD	x+0(FP), R0
	MOVW	prec+8(FP), R1
	CALL	·bigfloatSqrtGeneric(SB)
	MOVD	R0, ret+16(FP)
	RET

