// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

#include "textflag.h"

// Power function for ARM64

// func bigPowAsm(x, y *BigFloat, prec uint) *BigFloat
TEXT ·bigPowAsm(SB), $128-32
	MOVD	x+0(FP), R0
	MOVD	y+8(FP), R1
	MOVW	prec+16(FP), R2

	MOVD	R0, 0(RSP)
	MOVD	R1, 8(RSP)
	MOVW	R2, 16(RSP)
	CALL	·bigPowGeneric(SB)
	MOVD	24(RSP), R0
	MOVD	R0, ret+24(FP)
	RET

