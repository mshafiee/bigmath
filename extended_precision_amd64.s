// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

#include "textflag.h"

// Extended precision (80-bit x87 FPU) assembly implementations
// These functions use the x87 FPU stack for hardware extended precision operations.

// ExtendedAdd: result = a + b
TEXT ·extendedAdd(SB), NOSPLIT, $0
	// Load arguments from stack (float64 in XMM registers)
	MOVSD a+0(FP), X0
	MOVSD b+8(FP), X1
	// Convert to x87: store to memory then load to x87 stack
	MOVSD X0, -8(SP)
	FLDL -8(SP)      // Load a to ST(0) as extended precision
	MOVSD X1, -8(SP)
	FADDL -8(SP)     // Add b to ST(0), result in ST(0)
	// Store result back (extended precision to double)
	FSTPL -8(SP)     // Store extended precision result to memory
	MOVSD -8(SP), X0
	MOVSD X0, ret+16(FP)
	RET

// ExtendedSub: result = a - b
TEXT ·extendedSub(SB), NOSPLIT, $0
	MOVSD a+0(FP), X0
	MOVSD b+8(FP), X1
	MOVSD X0, -8(SP)
	FLDL -8(SP)      // Load a to ST(0)
	MOVSD X1, -8(SP)
	FSUBL -8(SP)     // Subtract b from ST(0)
	FSTPL -8(SP)
	MOVSD -8(SP), X0
	MOVSD X0, ret+16(FP)
	RET

// ExtendedMul: result = a * b
TEXT ·extendedMul(SB), NOSPLIT, $0
	MOVSD a+0(FP), X0
	MOVSD b+8(FP), X1
	MOVSD X0, -8(SP)
	FLDL -8(SP)      // Load a to ST(0)
	MOVSD X1, -8(SP)
	FMULL -8(SP)     // Multiply ST(0) by b
	FSTPL -8(SP)
	MOVSD -8(SP), X0
	MOVSD X0, ret+16(FP)
	RET

// ExtendedDiv: result = a / b
TEXT ·extendedDiv(SB), NOSPLIT, $0
	MOVSD a+0(FP), X0
	MOVSD b+8(FP), X1
	MOVSD X0, -8(SP)
	FLDL -8(SP)      // Load a to ST(0)
	MOVSD X1, -8(SP)
	FDIVL -8(SP)     // Divide ST(0) by b
	FSTPL -8(SP)
	MOVSD -8(SP), X0
	MOVSD X0, ret+16(FP)
	RET

// ExtendedSin: result = sin(x)
TEXT ·extendedSin(SB), NOSPLIT, $0
	MOVSD x+0(FP), X0
	MOVSD X0, -8(SP)
	FLDL -8(SP)      // Load x to ST(0)
	FSIN             // sin(ST(0)) -> ST(0)
	FSTPL -8(SP)
	MOVSD -8(SP), X0
	MOVSD X0, ret+8(FP)
	RET

// ExtendedCos: result = cos(x)
TEXT ·extendedCos(SB), NOSPLIT, $0
	MOVSD x+0(FP), X0
	MOVSD X0, -8(SP)
	FLDL -8(SP)      // Load x to ST(0)
	FCOS             // cos(ST(0)) -> ST(0)
	FSTPL -8(SP)
	MOVSD -8(SP), X0
	MOVSD X0, ret+8(FP)
	RET

// ExtendedTan: result = tan(x)
TEXT ·extendedTan(SB), NOSPLIT, $0
	MOVSD x+0(FP), X0
	MOVSD X0, -8(SP)
	FLDL -8(SP)      // Load x to ST(0)
	FSIN             // sin(x) -> ST(0)
	MOVSD X0, -8(SP)
	FLDL -8(SP)      // Reload x to ST(0), sin(x) -> ST(1)
	FCOS             // cos(x) -> ST(0), sin(x) -> ST(1)
	FDIVP            // ST(1)/ST(0) -> ST(0) (tan = sin/cos)
	FSTPL -8(SP)
	MOVSD -8(SP), X0
	MOVSD X0, ret+8(FP)
	RET

// ExtendedAtan: result = atan(x)
TEXT ·extendedAtan(SB), NOSPLIT, $0
	MOVSD x+0(FP), X0
	// For atan(x), we compute atan2(x, 1.0)
	MOVSD X0, -8(SP)
	FLDL -8(SP)      // Load x to ST(0)
	MOVQ $0x3FF0000000000000, AX  // 1.0 as double
	MOVQ AX, -8(SP)
	FLDL -8(SP)      // Load 1.0 to ST(0), x is now ST(1)
	FXCH             // Swap: ST(0)=x, ST(1)=1.0
	FPATAN           // atan2(x, 1.0) = atan(x) -> ST(0)
	FSTPL -8(SP)
	MOVSD -8(SP), X0
	MOVSD X0, ret+8(FP)
	RET

// ExtendedAtan2: result = atan2(y, x)
TEXT ·extendedAtan2(SB), NOSPLIT, $0
	MOVSD y+0(FP), X0
	MOVSD x+8(FP), X1
	MOVSD X0, -8(SP)
	FLDL -8(SP)      // Load y to ST(0)
	MOVSD X1, -8(SP)
	FLDL -8(SP)      // Load x to ST(0), y is now ST(1)
	FXCH             // Swap: ST(0)=y, ST(1)=x
	FPATAN           // atan2(y, x) -> ST(0)
	FSTPL -8(SP)
	MOVSD -8(SP), X0
	MOVSD X0, ret+16(FP)
	RET

// ExtendedExp: result = exp(x)
// Note: x87 doesn't have direct exp, so we use f2xm1 and scale
// exp(x) = 2^(x * log2(e)) = 2^(x / ln(2))
TEXT ·extendedExp(SB), NOSPLIT, $0
	MOVSD x+0(FP), X0
	MOVSD X0, -8(SP)
	FLDL -8(SP)      // Load x to ST(0)
	// exp(x) = 2^(x / ln(2))
	MOVQ $0x3FF71547652B82FE, AX  // 1/ln(2) as double
	MOVQ AX, -8(SP)
	FMULL -8(SP)     // x / ln(2) -> ST(0)
	// Now compute 2^ST(0)
	FLD ST(0)         // Duplicate
	FRNDINT           // Integer part -> ST(0)
	FSUB ST(1), ST(0) // Fractional part -> ST(0), original -> ST(1)
	FXCH              // Swap: fractional -> ST(1), original -> ST(0)
	FSTP ST(0)        // Pop original, keep fractional
	F2XM1             // 2^fractional - 1 -> ST(0)
	FLD1
	FADD              // 2^fractional -> ST(0)
	FXCH              // Swap with integer part
	FLD1
	FSCALE            // 2^integer -> ST(0)
	FSTP ST(1)        // Pop integer part
	FMUL              // 2^fractional * 2^integer = 2^(x/ln(2)) = exp(x)
	FSTPL -8(SP)
	MOVSD -8(SP), X0
	MOVSD X0, ret+8(FP)
	RET

// ExtendedLog: result = log(x) (natural logarithm)
TEXT ·extendedLog(SB), NOSPLIT, $0
	MOVSD x+0(FP), X0
	MOVSD X0, -8(SP)
	FLDL -8(SP)      // Load x to ST(0)
	FLD1             // Load 1.0 to ST(0), x -> ST(1)
	FXCH             // Swap: ST(0)=x, ST(1)=1.0
	FYL2X            // log2(x) -> ST(0) (pops both, result in ST(0))
	// ln(x) = log2(x) * ln(2)
	MOVQ $0x3FE62E42FEFA39EF, AX  // ln(2) as double
	MOVQ AX, -8(SP)
	FMULL -8(SP)     // Multiply by ln(2) to get ln(x)
	FSTPL -8(SP)
	MOVSD -8(SP), X0
	MOVSD X0, ret+8(FP)
	RET

// ExtendedSqrt: result = sqrt(x)
TEXT ·extendedSqrt(SB), NOSPLIT, $0
	MOVSD x+0(FP), X0
	MOVSD X0, -8(SP)
	FLDL -8(SP)      // Load x to ST(0)
	FSQRT            // sqrt(ST(0)) -> ST(0)
	FSTPL -8(SP)
	MOVSD -8(SP), X0
	MOVSD X0, ret+8(FP)
	RET

// ExtendedPow: result = pow(x, y) = x^y
TEXT ·extendedPow(SB), NOSPLIT, $0
	MOVSD x+0(FP), X0
	MOVSD y+8(FP), X1
	MOVSD X0, -8(SP)
	FLDL -8(SP)      // Load x to ST(0)
	MOVSD X1, -8(SP)
	FLDL -8(SP)      // Load y to ST(0), x is now ST(1)
	FXCH             // Swap: ST(0)=x, ST(1)=y
	// x^y = 2^(y * log2(x))
	FLD1             // Load 1.0 to ST(0), x -> ST(1), y -> ST(2)
	FXCH             // ST(0)=1.0, ST(1)=x, ST(2)=y
	FYL2X            // y * log2(x) -> ST(0) (pops x and y)
	// Now compute 2^ST(0)
	FLD ST(0)        // Duplicate y*log2(x)
	FRNDINT          // Integer part -> ST(0)
	FSUB ST(1), ST(0) // Fractional part -> ST(0), original -> ST(1)
	FXCH             // Swap: fractional -> ST(1), original -> ST(0)
	FSTP ST(0)       // Pop original, keep fractional
	F2XM1            // 2^fractional - 1 -> ST(0)
	FLD1
	FADD             // 2^fractional -> ST(0)
	FXCH             // Swap with integer part
	FLD1
	FSCALE           // 2^integer -> ST(0)
	FSTP ST(1)       // Pop integer part
	FMUL             // 2^fractional * 2^integer = 2^(y*log2(x)) = x^y
	FSTPL -8(SP)
	MOVSD -8(SP), X0
	MOVSD X0, ret+16(FP)
	RET

