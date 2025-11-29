// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

#include "textflag.h"

// Extract IEEE 754 double precision components from 64-bit value
// func extractIEEE754Components(bits uint64) (sign uint64, exponent int64, mantissa uint64)
TEXT ·extractIEEE754Components(SB), NOSPLIT, $0-32
	MOVQ bits+0(FP), AX        // AX = bits (IEEE 754 double)

	// Extract sign (bit 63)
	MOVQ AX, BX
	SHRQ $63, BX                // BX = sign (0 or 1)
	MOVQ BX, sign+8(FP)

	// Extract exponent (bits 52-62, 11 bits)
	MOVQ AX, CX
	SHRQ $52, CX                // Shift right 52 bits
	ANDQ $0x7FF, CX             // Mask to get 11 bits (0x7FF = 2047)
	MOVQ CX, exponent+16(FP)    // Store as int64 (sign-extended)

	// Extract mantissa (bits 0-51, 52 bits)
	MOVQ AX, DX
	ANDQ $0xFFFFFFFFFFFFF, DX   // Mask to get lower 52 bits
	MOVQ DX, mantissa+24(FP)

	RET

// Convert 8 bytes to uint64 with endianness conversion
// func convertEndiannessBytes(bytes *[8]byte, bigEndian uint8) uint64
TEXT ·convertEndiannessBytes(SB), NOSPLIT, $0-24
	MOVQ bytes+0(FP), SI         // SI = bytes pointer
	MOVB bigEndian+8(FP), AL     // AL = bigEndian flag (0 or 1)

	// Load 8 bytes as uint64
	MOVQ (SI), AX                // AX = bytes as uint64 (little-endian)

	// Check if big-endian conversion needed
	TESTB AL, AL                 // Test if bigEndian != 0
	JZ little_endian             // If zero, skip swap

	// Big-endian: swap bytes using BSWAP
	BSWAPQ AX                    // Byte swap the 64-bit value
	JMP done

little_endian:
	// Little-endian: value already correct (native x86-64 is little-endian)
	// No operation needed

done:
	MOVQ AX, ret+16(FP)          // Return the converted value
	RET

