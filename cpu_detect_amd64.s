// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

#include "textflag.h"

// func cpuidAVX() bool
TEXT ·cpuidAVX(SB), NOSPLIT, $0-1
	// CPUID function 1
	MOVL	$1, AX
	CPUID
	// Check ECX bit 28 (AVX)
	SHRL	$28, CX
	ANDL	$1, CX
	MOVB	CX, ret+0(FP)
	RET

// func cpuidAVX2() bool
TEXT ·cpuidAVX2(SB), NOSPLIT, $0-1
	// Check max extended function
	MOVL	$0, AX
	CPUID
	CMPL	AX, $7
	JL		no_avx2
	
	// CPUID function 7, subfunction 0
	MOVL	$7, AX
	MOVL	$0, CX
	CPUID
	// Check EBX bit 5 (AVX2)
	SHRL	$5, BX
	ANDL	$1, BX
	MOVB	BX, ret+0(FP)
	RET

no_avx2:
	MOVB	$0, ret+0(FP)
	RET

// func cpuidAVX512() bool
TEXT ·cpuidAVX512(SB), NOSPLIT, $0-1
	// Check max extended function
	MOVL	$0, AX
	CPUID
	CMPL	AX, $7
	JL		no_avx512
	
	// CPUID function 7, subfunction 0
	MOVL	$7, AX
	MOVL	$0, CX
	CPUID
	// Check EBX bit 16 (AVX-512F)
	SHRL	$16, BX
	ANDL	$1, BX
	MOVB	BX, ret+0(FP)
	RET

no_avx512:
	MOVB	$0, ret+0(FP)
	RET

// func cpuidFMA() bool
TEXT ·cpuidFMA(SB), NOSPLIT, $0-1
	// CPUID function 1
	MOVL	$1, AX
	CPUID
	// Check ECX bit 12 (FMA)
	SHRL	$12, CX
	ANDL	$1, CX
	MOVB	CX, ret+0(FP)
	RET

// func cpuidBMI2() bool
TEXT ·cpuidBMI2(SB), NOSPLIT, $0-1
	// Check max extended function
	MOVL	$0, AX
	CPUID
	CMPL	AX, $7
	JL	no_bmi2
	
	// CPUID function 7, subfunction 0
	MOVL	$7, AX
	MOVL	$0, CX
	CPUID
	// Check EBX bit 8 (BMI2)
	SHRL	$8, BX
	ANDL	$1, BX
	MOVB	BX, ret+0(FP)
	RET

no_bmi2:
	MOVB	$0, ret+0(FP)
	RET

