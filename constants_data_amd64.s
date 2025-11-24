// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

#include "textflag.h"

// High-precision constants stored in read-only data section
// These are precomputed values for common mathematical constants
// NOTE: These are placeholder values. Actual constants are computed in Go init().
// Commented out to avoid linking issues - constants are computed at runtime.

/*
// π (Pi) - 256 bits precision (8 limbs of 64 bits each)
// Computed using Chudnovsky algorithm
DATA ·bigPI+0(SB)/8, $0x243f6a8885a308d3  // π limb 0
DATA ·bigPI+8(SB)/8, $0x13198a2e03707344  // π limb 1
DATA ·bigPI+16(SB)/8, $0xa4093822299f31d0  // π limb 2
DATA ·bigPI+24(SB)/8, $0x082efa98ec4e6c89  // π limb 3
DATA ·bigPI+32(SB)/8, $0x452821e638d01377  // π limb 4
DATA ·bigPI+40(SB)/8, $0xbe5466cf34e90c6c  // π limb 5
DATA ·bigPI+48(SB)/8, $0xc0ac29b7c97c50dd  // π limb 6
DATA ·bigPI+56(SB)/8, $0x3f84d5b5b5470917  // π limb 7
GLOBL ·bigPI(SB), RODATA, $64

// ln(2) - natural logarithm of 2
DATA ·bigLn2+0(SB)/8, $0xb17217f7d1cf79ab  // ln(2) limb 0
DATA ·bigLn2+8(SB)/8, $0xc9e3b39803f2f6af  // ln(2) limb 1
DATA ·bigLn2+16(SB)/8, $0x40f343267298b62d  // ln(2) limb 2
DATA ·bigLn2+24(SB)/8, $0x8c5ec1a3a03fbdff  // ln(2) limb 3
DATA ·bigLn2+32(SB)/8, $0x39e9c4a2d4f91b9d  // ln(2) limb 4
DATA ·bigLn2+40(SB)/8, $0x5e2d58d8b3bdf817  // ln(2) limb 5
DATA ·bigLn2+48(SB)/8, $0x9b5cb8f40692823d  // ln(2) limb 6
DATA ·bigLn2+56(SB)/8, $0x3fef324e7738925e  // ln(2) limb 7
GLOBL ·bigLn2(SB), RODATA, $64

// ln(10) - natural logarithm of 10
DATA ·bigLn10+0(SB)/8, $0x935d8dddaaa8ac16  // ln(10) limb 0
DATA ·bigLn10+8(SB)/8, $0x7e37be2022c09a98  // ln(10) limb 1
DATA ·bigLn10+16(SB)/8, $0x4c2eb6872a1f258b  // ln(10) limb 2
DATA ·bigLn10+24(SB)/8, $0x76e381ac4beaadf4  // ln(10) limb 3
DATA ·bigLn10+32(SB)/8, $0x9e8bc9b24b8e9b0b  // ln(10) limb 4
DATA ·bigLn10+40(SB)/8, $0x7e9474bf8eb5cdc0  // ln(10) limb 5
DATA ·bigLn10+48(SB)/8, $0x3e7abc9e88b57d2f  // ln(10) limb 6
DATA ·bigLn10+56(SB)/8, $0x40026bb1bbb55516  // ln(10) limb 7
GLOBL ·bigLn10(SB), RODATA, $64

// e (Euler's number) - base of natural logarithm
DATA ·bigE+0(SB)/8, $0xadf85458a2bb4a9a  // e limb 0
DATA ·bigE+8(SB)/8, $0xafb562e6ab36a1e7  // e limb 1
DATA ·bigE+16(SB)/8, $0x8b8f8bb8b8f8b8b8  // e limb 2
DATA ·bigE+24(SB)/8, $0x8b8f8bb8b8f8b8b8  // e limb 3
DATA ·bigE+32(SB)/8, $0x8b8f8bb8b8f8b8b8  // e limb 4
DATA ·bigE+40(SB)/8, $0x8b8f8bb8b8f8b8b8  // e limb 5
DATA ·bigE+48(SB)/8, $0x8b8f8bb8b8f8b8b8  // e limb 6
DATA ·bigE+56(SB)/8, $0x4002b7e151628aed  // e limb 7
GLOBL ·bigE(SB), RODATA, $64

// √2 (square root of 2)
DATA ·bigSqrt2+0(SB)/8, $0xb504f333f9de6484  // √2 limb 0
DATA ·bigSqrt2+8(SB)/8, $0x597ed6a310dd0c51  // √2 limb 1
DATA ·bigSqrt2+16(SB)/8, $0x8c4f1b9e6b169e57  // √2 limb 2
DATA ·bigSqrt2+24(SB)/8, $0x8c4f1b9e6b169e57  // √2 limb 3
DATA ·bigSqrt2+32(SB)/8, $0x8c4f1b9e6b169e57  // √2 limb 4
DATA ·bigSqrt2+40(SB)/8, $0x8c4f1b9e6b169e57  // √2 limb 5
DATA ·bigSqrt2+48(SB)/8, $0x8c4f1b9e6b169e57  // √2 limb 6
DATA ·bigSqrt2+56(SB)/8, $0x3ff6a09e667f3bcd  // √2 limb 7
GLOBL ·bigSqrt2(SB), RODATA, $64
*/

// Note: The actual values above are placeholders
// In a production implementation, these would be computed using
// high-precision algorithms (Chudnovsky for π, series for others)
// and stored with full 256-bit precision

