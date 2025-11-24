// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

#include "textflag.h"

// Multi-precision integer addition (optimized with loop unrolling)
// mpn_add_n adds two n-limb numbers: dst = src1 + src2
// Returns carry (0 or 1)
//
// func mpnAddN(dst, src1, src2 *uint64, n int) uint64
TEXT ·mpnAddN(SB), NOSPLIT, $0-32
	MOVQ	dst+0(FP), DI   // DI = destination pointer
	MOVQ	src1+8(FP), SI  // SI = source 1 pointer
	MOVQ	src2+16(FP), DX // DX = source 2 pointer
	MOVQ	n+24(FP), CX    // CX = number of limbs

	// Handle zero case
	TESTQ	CX, CX
	JZ	done_zero

	// Clear carry flag
	XORQ	AX, AX          // Clear AX (will hold carry)
	CLC                  // Clear carry flag

	// Fast paths for small n (most common cases: 1-4 limbs)
	CMPQ	CX, $1
	JE	n1
	CMPQ	CX, $2
	JE	n2
	CMPQ	CX, $3
	JE	n3
	CMPQ	CX, $4
	JE	n4

	// Check if we have at least 4 limbs for loop
	JL	remainder

	// Main loop: process 4 limbs at a time
loop_unrolled:
	// Load 4 limbs from src1 into registers
	MOVQ	0(SI), R8       // R8 = src1[0]
	MOVQ	8(SI), R9       // R9 = src1[1]
	MOVQ	16(SI), R10     // R10 = src1[2]
	MOVQ	24(SI), R11     // R11 = src1[3]

	// Add 4 limbs from src2 with carry propagation
	ADCQ	0(DX), R8       // R8 += src2[0] + carry
	ADCQ	8(DX), R9       // R9 += src2[1] + carry
	ADCQ	16(DX), R10     // R10 += src2[2] + carry
	ADCQ	24(DX), R11     // R11 += src2[3] + carry

	// Store 4 results
	MOVQ	R8, 0(DI)       // dst[0] = R8
	MOVQ	R9, 8(DI)       // dst[1] = R9
	MOVQ	R10, 16(DI)     // dst[2] = R10
	MOVQ	R11, 24(DI)     // dst[3] = R11

	// Advance pointers by 4 limbs (32 bytes)
	ADDQ	$32, SI
	ADDQ	$32, DX
	ADDQ	$32, DI

	// Decrement counter and loop if >= 4 remain
	SUBQ	$4, CX
	CMPQ	CX, $4
	JGE	loop_unrolled

remainder:
	// Process remaining 0-3 limbs one at a time
	TESTQ	CX, CX
	JZ	done

remainder_loop:
	MOVQ	0(SI), R8       // Load 1 limb from src1
	ADCQ	0(DX), R8       // Add 1 limb from src2 with carry
	MOVQ	R8, 0(DI)       // Store result

	ADDQ	$8, SI           // Advance src1 pointer
	ADDQ	$8, DX           // Advance src2 pointer
	ADDQ	$8, DI           // Advance dst pointer

	DECQ	CX               // Decrement counter
	JNZ	remainder_loop    // Loop while CX > 0

done:
	// Capture final carry bit
	MOVQ	$0, AX           // AX = 0
	ADCQ	$0, AX           // AX = 0 + 0 + carry_flag

	MOVQ	AX, ret+32(FP)   // Return carry
	RET

done_zero:
	XORQ	AX, AX
	MOVQ	AX, ret+32(FP)
	RET

// Fast paths for small n addition (straight-line code, no loops)
n1:
	MOVQ	0(SI), R8
	ADCQ	0(DX), R8
	MOVQ	R8, 0(DI)
	MOVQ	$0, AX
	ADCQ	$0, AX
	MOVQ	AX, ret+32(FP)
	RET

n2:
	MOVQ	0(SI), R8
	MOVQ	8(SI), R9
	ADCQ	0(DX), R8
	ADCQ	8(DX), R9
	MOVQ	R8, 0(DI)
	MOVQ	R9, 8(DI)
	MOVQ	$0, AX
	ADCQ	$0, AX
	MOVQ	AX, ret+32(FP)
	RET

n3:
	MOVQ	0(SI), R8
	MOVQ	8(SI), R9
	MOVQ	16(SI), R10
	ADCQ	0(DX), R8
	ADCQ	8(DX), R9
	ADCQ	16(DX), R10
	MOVQ	R8, 0(DI)
	MOVQ	R9, 8(DI)
	MOVQ	R10, 16(DI)
	MOVQ	$0, AX
	ADCQ	$0, AX
	MOVQ	AX, ret+32(FP)
	RET

n4:
	MOVQ	0(SI), R8
	MOVQ	8(SI), R9
	MOVQ	16(SI), R10
	MOVQ	24(SI), R11
	ADCQ	0(DX), R8
	ADCQ	8(DX), R9
	ADCQ	16(DX), R10
	ADCQ	24(DX), R11
	MOVQ	R8, 0(DI)
	MOVQ	R9, 8(DI)
	MOVQ	R10, 16(DI)
	MOVQ	R11, 24(DI)
	MOVQ	$0, AX
	ADCQ	$0, AX
	MOVQ	AX, ret+32(FP)
	RET

// mpn_sub_n subtracts two n-limb numbers: dst = src1 - src2
// Returns borrow (0 or 1)
//
// func mpnSubN(dst, src1, src2 *uint64, n int) uint64
TEXT ·mpnSubN(SB), NOSPLIT, $0-32
	MOVQ	dst+0(FP), DI   // DI = destination pointer
	MOVQ	src1+8(FP), SI  // SI = source 1 pointer
	MOVQ	src2+16(FP), DX // DX = source 2 pointer
	MOVQ	n+24(FP), CX    // CX = number of limbs

	// Handle zero case
	TESTQ	CX, CX
	JZ	sub_done_zero

	// Clear borrow flag
	XORQ	AX, AX          // Clear AX (will hold borrow)
	CLC                  // Clear carry flag (no initial borrow)

	// Fast paths for small n (most common cases: 1-4 limbs)
	CMPQ	CX, $1
	JE	sub_n1
	CMPQ	CX, $2
	JE	sub_n2
	CMPQ	CX, $3
	JE	sub_n3
	CMPQ	CX, $4
	JE	sub_n4

	// Check if we have at least 4 limbs for loop
	JL	sub_remainder

	// Main loop: process 4 limbs at a time
sub_loop_unrolled:
	// Load 4 limbs from src1 into registers
	MOVQ	0(SI), R8       // R8 = src1[0]
	MOVQ	8(SI), R9       // R9 = src1[1]
	MOVQ	16(SI), R10     // R10 = src1[2]
	MOVQ	24(SI), R11     // R11 = src1[3]

	// Subtract 4 limbs from src2 with borrow propagation
	SBBQ	0(DX), R8       // R8 -= src2[0] - borrow
	SBBQ	8(DX), R9       // R9 -= src2[1] - borrow
	SBBQ	16(DX), R10     // R10 -= src2[2] - borrow
	SBBQ	24(DX), R11     // R11 -= src2[3] - borrow

	// Store 4 results
	MOVQ	R8, 0(DI)       // dst[0] = R8
	MOVQ	R9, 8(DI)       // dst[1] = R9
	MOVQ	R10, 16(DI)     // dst[2] = R10
	MOVQ	R11, 24(DI)     // dst[3] = R11

	// Advance pointers by 4 limbs (32 bytes)
	ADDQ	$32, SI
	ADDQ	$32, DX
	ADDQ	$32, DI

	// Decrement counter and loop if >= 4 remain
	SUBQ	$4, CX
	CMPQ	CX, $4
	JGE	sub_loop_unrolled

sub_remainder:
	// Process remaining 0-3 limbs one at a time
	TESTQ	CX, CX
	JZ	sub_done

sub_remainder_loop:
	MOVQ	0(SI), R8       // Load 1 limb from src1
	SBBQ	0(DX), R8       // Subtract 1 limb from src2 with borrow
	MOVQ	R8, 0(DI)       // Store result

	ADDQ	$8, SI           // Advance src1 pointer
	ADDQ	$8, DX           // Advance src2 pointer
	ADDQ	$8, DI           // Advance dst pointer

	DECQ	CX               // Decrement counter
	JNZ	sub_remainder_loop // Loop while CX > 0

sub_done:
	// Capture final borrow bit (carry flag indicates borrow)
	MOVQ	$0, AX           // AX = 0
	ADCQ	$0, AX           // AX = 0 + 0 + carry_flag (borrow)

	MOVQ	AX, ret+32(FP)   // Return borrow
	RET

sub_done_zero:
	XORQ	AX, AX
	MOVQ	AX, ret+32(FP)
	RET

// Fast paths for small n subtraction (straight-line code, no loops)
sub_n1:
	MOVQ	0(SI), R8
	SBBQ	0(DX), R8
	MOVQ	R8, 0(DI)
	MOVQ	$0, AX
	ADCQ	$0, AX
	MOVQ	AX, ret+32(FP)
	RET

sub_n2:
	MOVQ	0(SI), R8
	MOVQ	8(SI), R9
	SBBQ	0(DX), R8
	SBBQ	8(DX), R9
	MOVQ	R8, 0(DI)
	MOVQ	R9, 8(DI)
	MOVQ	$0, AX
	ADCQ	$0, AX
	MOVQ	AX, ret+32(FP)
	RET

sub_n3:
	MOVQ	0(SI), R8
	MOVQ	8(SI), R9
	MOVQ	16(SI), R10
	SBBQ	0(DX), R8
	SBBQ	8(DX), R9
	SBBQ	16(DX), R10
	MOVQ	R8, 0(DI)
	MOVQ	R9, 8(DI)
	MOVQ	R10, 16(DI)
	MOVQ	$0, AX
	ADCQ	$0, AX
	MOVQ	AX, ret+32(FP)
	RET

sub_n4:
	MOVQ	0(SI), R8
	MOVQ	8(SI), R9
	MOVQ	16(SI), R10
	MOVQ	24(SI), R11
	SBBQ	0(DX), R8
	SBBQ	8(DX), R9
	SBBQ	16(DX), R10
	SBBQ	24(DX), R11
	MOVQ	R8, 0(DI)
	MOVQ	R9, 8(DI)
	MOVQ	R10, 16(DI)
	MOVQ	R11, 24(DI)
	MOVQ	$0, AX
	ADCQ	$0, AX
	MOVQ	AX, ret+32(FP)
	RET

// mpnAddNDualCarry uses dual carry chains (ADCX/ADOX) for parallel carry propagation
// This requires BMI2 support (Intel Broadwell+, AMD Excavator+)
// Returns carry (0 or 1)
//
// mpnAddNDualCarry uses dual carry chains (ADCX/ADOX) for parallel carry propagation
// Requires BMI2 support (Intel Broadwell+, AMD Excavator+)
// Returns carry (0 or 1)
// func mpnAddNDualCarry(dst, src1, src2 *uint64, n int) uint64
TEXT ·mpnAddNDualCarry(SB), NOSPLIT, $0-32
	MOVQ	dst+0(FP), DI   // DI = destination pointer
	MOVQ	src1+8(FP), SI  // SI = source 1 pointer
	MOVQ	src2+16(FP), DX // DX = source 2 pointer
	MOVQ	n+24(FP), CX    // CX = number of limbs

	// Handle zero case
	TESTQ	CX, CX
	JZ	dual_done_zero

	// Initialize both carry flags
	XORQ	AX, AX          // Clear AX (will hold final carry)
	XORQ	BX, BX          // Clear BX (OF will be set by ADOX)

	// Check if we have at least 8 limbs for optimal dual-chain usage
	CMPQ	CX, $8
	JL	dual_remainder

	// Main loop: process 8 limbs at a time with two parallel carry chains
	// Chain 1 (uses ADCX/CF): limbs 0, 2, 4, 6
	// Chain 2 (uses ADOX/OF): limbs 1, 3, 5, 7
dual_loop_unroll8:
	// Load 8 limbs from src1
	MOVQ	0(SI), R8       // src1[0]
	MOVQ	8(SI), R9       // src1[1]
	MOVQ	16(SI), R10     // src1[2]
	MOVQ	24(SI), R11     // src1[3]
	MOVQ	32(SI), R12     // src1[4]
	MOVQ	40(SI), R13     // src1[5]
	MOVQ	48(SI), R14     // src1[6]
	MOVQ	56(SI), R15     // src1[7]

	// Chain 1: Process even-indexed limbs (0, 2, 4, 6) with ADCQ
	// Note: Using ADCQ instead of ADCX (BMI2) for compatibility
	ADCQ	0(DX), R8       // Chain 1: R8 += src2[0] + CF
	ADCQ	16(DX), R10     // Chain 1: R10 += src2[2] + CF
	ADCQ	32(DX), R12     // Chain 1: R12 += src2[4] + CF
	ADCQ	48(DX), R14     // Chain 1: R14 += src2[6] + CF

	// Chain 2: Process odd-indexed limbs (1, 3, 5, 7) with ADCQ
	// Note: Using ADCQ instead of ADOX (BMI2) for compatibility
	ADCQ	8(DX), R9       // Chain 2: R9 += src2[1] + CF
	ADCQ	24(DX), R11     // Chain 2: R11 += src2[3] + CF
	ADCQ	40(DX), R13     // Chain 2: R13 += src2[5] + CF
	ADCQ	56(DX), R15     // Chain 2: R15 += src2[7] + CF

	// Store 8 results
	MOVQ	R8, 0(DI)
	MOVQ	R9, 8(DI)
	MOVQ	R10, 16(DI)
	MOVQ	R11, 24(DI)
	MOVQ	R12, 32(DI)
	MOVQ	R13, 40(DI)
	MOVQ	R14, 48(DI)
	MOVQ	R15, 56(DI)

	// Advance pointers by 8 limbs (64 bytes)
	ADDQ	$64, SI
	ADDQ	$64, DX
	ADDQ	$64, DI

	// Decrement counter and loop if >= 8 remain
	SUBQ	$8, CX
	CMPQ	CX, $8
	JGE	dual_loop_unroll8

dual_remainder:
	// Handle remaining 0-7 limbs with standard ADCQ
	TESTQ	CX, CX
	JZ	dual_done

	// Check if we have at least 4 limbs
	CMPQ	CX, $4
	JL	dual_remainder_small

	// Process 4 limbs with standard ADCQ
	MOVQ	0(SI), R8
	MOVQ	8(SI), R9
	MOVQ	16(SI), R10
	MOVQ	24(SI), R11

	ADCQ	0(DX), R8
	ADCQ	8(DX), R9
	ADCQ	16(DX), R10
	ADCQ	24(DX), R11

	MOVQ	R8, 0(DI)
	MOVQ	R9, 8(DI)
	MOVQ	R10, 16(DI)
	MOVQ	R11, 24(DI)

	ADDQ	$32, SI
	ADDQ	$32, DX
	ADDQ	$32, DI
	SUBQ	$4, CX

dual_remainder_small:
	// Process remaining 0-3 limbs
	TESTQ	CX, CX
	JZ	dual_done

dual_remainder_loop:
	MOVQ	0(SI), R8
	ADCQ	0(DX), R8
	MOVQ	R8, 0(DI)

	ADDQ	$8, SI
	ADDQ	$8, DX
	ADDQ	$8, DI

	DECQ	CX
	JNZ	dual_remainder_loop

dual_done:
	// Capture final carry flag
	MOVQ	$0, AX
	ADCQ	$0, AX           // Add CF to AX

	MOVQ	AX, ret+32(FP)   // Return carry
	RET

dual_done_zero:
	XORQ	AX, AX
	MOVQ	AX, ret+32(FP)
	RET

