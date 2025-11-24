#include "textflag.h"

// Chebyshev polynomial evaluation for ARM64

// func evaluateChebyshevBigAsm(t *BigFloat, c []*BigFloat, neval int, prec uint) *BigFloat
TEXT ·evaluateChebyshevBigAsm(SB), NOSPLIT, $128-48
	MOVD	t+0(FP), R0
	MOVD	c+8(FP), R1
	MOVD	c+16(FP), R2
	MOVD	neval+24(FP), R3
	MOVW	prec+32(FP), R4

	MOVD	R0, 0(RSP)
	MOVD	R1, 8(RSP)
	MOVD	R2, 16(RSP)
	MOVD	R3, 24(RSP)
	MOVW	R4, 32(RSP)
	CALL	·evaluateChebyshevBigGeneric(SB)
	MOVD	40(RSP), R0
	MOVD	R0, ret+40(FP)
	RET

// func evaluateChebyshevDerivativeBigAsm(t *BigFloat, c []*BigFloat, neval int, prec uint) *BigFloat
TEXT ·evaluateChebyshevDerivativeBigAsm(SB), NOSPLIT, $128-48
	MOVD	t+0(FP), R0
	MOVD	c+8(FP), R1
	MOVD	c+16(FP), R2
	MOVD	neval+24(FP), R3
	MOVW	prec+32(FP), R4

	MOVD	R0, 0(RSP)
	MOVD	R1, 8(RSP)
	MOVD	R2, 16(RSP)
	MOVD	R3, 24(RSP)
	MOVW	R4, 32(RSP)
	CALL	·evaluateChebyshevDerivativeBigGeneric(SB)
	MOVD	40(RSP), R0
	MOVD	R0, ret+40(FP)
	RET

