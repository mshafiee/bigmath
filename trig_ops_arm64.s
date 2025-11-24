#include "textflag.h"

// Trigonometric functions for ARM64

// func bigSinAsm(x *BigFloat, prec uint) *BigFloat
TEXT ·bigSinAsm(SB), NOSPLIT, $128-24
	MOVD	x+0(FP), R0
	MOVW	prec+8(FP), R1

	MOVD	R0, 0(RSP)
	MOVW	R1, 8(RSP)
	CALL	·bigSinGeneric(SB)
	MOVD	16(RSP), R0
	MOVD	R0, ret+16(FP)
	RET

// func bigCosAsm(x *BigFloat, prec uint) *BigFloat
TEXT ·bigCosAsm(SB), NOSPLIT, $128-24
	MOVD	x+0(FP), R0
	MOVW	prec+8(FP), R1

	MOVD	R0, 0(RSP)
	MOVW	R1, 8(RSP)
	CALL	·bigCosGeneric(SB)
	MOVD	16(RSP), R0
	MOVD	R0, ret+16(FP)
	RET

// func bigTanAsm(x *BigFloat, prec uint) *BigFloat
TEXT ·bigTanAsm(SB), NOSPLIT, $128-24
	MOVD	x+0(FP), R0
	MOVW	prec+8(FP), R1

	MOVD	R0, 0(RSP)
	MOVW	R1, 8(RSP)
	CALL	·bigTanGeneric(SB)
	MOVD	16(RSP), R0
	MOVD	R0, ret+16(FP)
	RET

// func bigAtanAsm(x *BigFloat, prec uint) *BigFloat
TEXT ·bigAtanAsm(SB), NOSPLIT, $128-24
	MOVD	x+0(FP), R0
	MOVW	prec+8(FP), R1

	MOVD	R0, 0(RSP)
	MOVW	R1, 8(RSP)
	CALL	·bigAtanGeneric(SB)
	MOVD	16(RSP), R0
	MOVD	R0, ret+16(FP)
	RET

// func bigAsinAsm(x *BigFloat, prec uint) *BigFloat
TEXT ·bigAsinAsm(SB), NOSPLIT, $128-24
	MOVD	x+0(FP), R0
	MOVW	prec+8(FP), R1

	MOVD	R0, 0(RSP)
	MOVW	R1, 8(RSP)
	CALL	·bigAsinGeneric(SB)
	MOVD	16(RSP), R0
	MOVD	R0, ret+16(FP)
	RET

// func bigAcosAsm(x *BigFloat, prec uint) *BigFloat
TEXT ·bigAcosAsm(SB), NOSPLIT, $128-24
	MOVD	x+0(FP), R0
	MOVW	prec+8(FP), R1

	MOVD	R0, 0(RSP)
	MOVW	R1, 8(RSP)
	CALL	·bigAcosGeneric(SB)
	MOVD	16(RSP), R0
	MOVD	R0, ret+16(FP)
	RET

// func bigAtan2Asm(y, x *BigFloat, prec uint) *BigFloat
TEXT ·bigAtan2Asm(SB), NOSPLIT, $128-32
	MOVD	y+0(FP), R0
	MOVD	x+8(FP), R1
	MOVW	prec+16(FP), R2

	MOVD	R0, 0(RSP)
	MOVD	R1, 8(RSP)
	MOVW	R2, 16(RSP)
	CALL	·bigAtan2Generic(SB)
	MOVD	24(RSP), R0
	MOVD	R0, ret+24(FP)
	RET

