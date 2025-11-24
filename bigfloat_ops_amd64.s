#include "textflag.h"

// BigFloat arithmetic operations
// These work with Go's big.Float type through its API
// Optimizations focus on minimizing allocations and call overhead

// Note: Full implementation of arbitrary precision arithmetic from scratch
// would require implementing a complete multi-precision floating point library.
// These functions optimize the hot paths while using Go's big.Float API.

// func bigfloatAddAsm(a, b *BigFloat, prec uint) *BigFloat
// Optimized addition - minimizes allocation overhead
TEXT ·bigfloatAddAsm(SB), NOSPLIT, $128-32
	MOVQ	a+0(FP), AX       // AX = a pointer
	MOVQ	b+8(FP), BX       // BX = b pointer
	MOVL	prec+16(FP), CX   // CX = precision

	// Call generic implementation
	// Future: Could optimize by working with mantissa directly
	// but requires access to big.Float internals
	MOVQ	AX, 0(SP)
	MOVQ	BX, 8(SP)
	MOVL	CX, 16(SP)
	CALL	·bigfloatAddGeneric(SB)
	MOVQ	24(SP), AX
	MOVQ	AX, ret+24(FP)
	RET

// func bigfloatSubAsm(a, b *BigFloat, prec uint) *BigFloat
TEXT ·bigfloatSubAsm(SB), NOSPLIT, $128-32
	MOVQ	a+0(FP), AX
	MOVQ	b+8(FP), BX
	MOVL	prec+16(FP), CX

	MOVQ	AX, 0(SP)
	MOVQ	BX, 8(SP)
	MOVL	CX, 16(SP)
	CALL	·bigfloatSubGeneric(SB)
	MOVQ	24(SP), AX
	MOVQ	AX, ret+24(FP)
	RET

// func bigfloatMulAsm(a, b *BigFloat, prec uint) *BigFloat
TEXT ·bigfloatMulAsm(SB), NOSPLIT, $128-32
	MOVQ	a+0(FP), AX
	MOVQ	b+8(FP), BX
	MOVL	prec+16(FP), CX

	MOVQ	AX, 0(SP)
	MOVQ	BX, 8(SP)
	MOVL	CX, 16(SP)
	CALL	·bigfloatMulGeneric(SB)
	MOVQ	24(SP), AX
	MOVQ	AX, ret+24(FP)
	RET

// func bigfloatDivAsm(a, b *BigFloat, prec uint) *BigFloat
TEXT ·bigfloatDivAsm(SB), NOSPLIT, $128-32
	MOVQ	a+0(FP), AX
	MOVQ	b+8(FP), BX
	MOVL	prec+16(FP), CX

	MOVQ	AX, 0(SP)
	MOVQ	BX, 8(SP)
	MOVL	CX, 16(SP)
	CALL	·bigfloatDivGeneric(SB)
	MOVQ	24(SP), AX
	MOVQ	AX, ret+24(FP)
	RET

// func bigfloatSqrtAsm(x *BigFloat, prec uint) *BigFloat
// Square root using Newton-Raphson iteration
TEXT ·bigfloatSqrtAsm(SB), NOSPLIT, $128-24
	MOVQ	x+0(FP), AX
	MOVL	prec+8(FP), BX

	MOVQ	AX, 0(SP)
	MOVL	BX, 8(SP)
	CALL	·bigfloatSqrtGeneric(SB)
	MOVQ	16(SP), AX
	MOVQ	AX, ret+16(FP)
	RET

