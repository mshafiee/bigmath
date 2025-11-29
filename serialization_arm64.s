// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

#include "textflag.h"

// Extract IEEE 754 double precision components from 64-bit value
// func extractIEEE754ComponentsARM64(bits uint64) (sign uint64, exponent int64, mantissa uint64)
TEXT ·extractIEEE754ComponentsARM64(SB), NOSPLIT, $0-32
	MOVD bits+0(FP), R0         // R0 = bits (IEEE 754 double)

	// Extract sign (bit 63)
	MOVD R0, R1
	LSR	$63, R1, R1            // R1 = sign (0 or 1)
	MOVD R1, sign+8(FP)

	// Extract exponent (bits 52-62, 11 bits)
	MOVD R0, R2
	LSR	$52, R2, R2             // Shift right 52 bits
	MOVD $0x7FF, R3              // Load mask (0x7FF = 2047)
	AND	R3, R2, R2              // Mask to get 11 bits
	MOVD R2, exponent+16(FP)     // Store as int64

	// Extract mantissa (bits 0-51, 52 bits)
	MOVD R0, R4
	MOVD $0xFFFFFFFFFFFFF, R5   // Load mask for lower 52 bits
	AND	R5, R4, R4              // Mask to get lower 52 bits
	MOVD R4, mantissa+24(FP)

	RET

// Convert 8 bytes to uint64 with endianness conversion
// func convertEndiannessBytesARM64(bytes *[8]byte, bigEndian uint8) uint64
TEXT ·convertEndiannessBytesARM64(SB), NOSPLIT, $0-24
	MOVD bytes+0(FP), R0         // R0 = bytes pointer
	MOVBU bigEndian+8(FP), R1    // R1 = bigEndian flag (0 or 1)

	// Load 8 bytes as uint64
	MOVD (R0), R2                // R2 = bytes as uint64 (native endian)

	// Check if big-endian conversion needed
	CMP	$0, R1                   // Compare bigEndian with 0
	BEQ	little_endian            // If zero, skip swap

	// Big-endian: reverse bytes using REV
	REV	R2, R2                   // Reverse byte order
	B	done

little_endian:
	// Little-endian: value already correct (ARM64 native is little-endian)
	// No operation needed

done:
	MOVD R2, ret+16(FP)           // Return the converted value
	RET

