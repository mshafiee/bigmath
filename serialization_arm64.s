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

// Combined function: extract IEEE 754 components directly from bytes with endianness conversion
// This reduces function call overhead by combining two operations into one
// Phase 3 & 4: Optimized with conditional execution and better register usage
// func extractIEEE754FromBytesARM64(bytes *[8]byte, bigEndian uint8) (sign uint64, exponent int64, mantissa uint64)
TEXT ·extractIEEE754FromBytesARM64(SB), NOSPLIT, $0-40
	MOVD bytes+0(FP), R0          // R0 = bytes pointer
	MOVBU bigEndian+8(FP), R1     // R1 = bigEndian flag (0 or 1)

	// Load 8 bytes as uint64
	MOVD (R0), R2                 // R2 = bytes as uint64 (native endian)

	// Phase 3: ARM64 doesn't have CMOV, but we can use conditional execution
	// Create swapped version in parallel
	MOVD R2, R8                   // R8 = copy of original
	REV  R8, R8                   // R8 = swapped version
	CMP  $0, R1                   // Compare bigEndian with 0
	CSEL EQ, R2, R8, R2           // If bigEndian == 0, use R2, else use R8 (CSEL = conditional select)

	// Phase 4: Optimize register usage - extract all components in parallel, keep in registers
	// Extract sign (bit 63)
	MOVD R2, R3                   // R3 = bits
	LSR  $63, R3, R3              // R3 = sign (0 or 1)

	// Extract exponent (bits 52-62, 11 bits)
	MOVD R2, R4                   // R4 = bits
	LSR  $52, R4, R4              // Shift right 52 bits
	MOVD $0x7FF, R5               // Load mask (0x7FF = 2047)
	AND  R5, R4, R4               // Mask to get 11 bits

	// Extract mantissa (bits 0-51, 52 bits)
	MOVD R2, R6                   // R6 = bits
	MOVD $0xFFFFFFFFFFFFF, R7     // Load mask for lower 52 bits
	AND  R7, R6, R6               // Mask to get lower 52 bits

	// Phase 4: Write all results in sequence for better cache locality
	MOVD R3, sign+16(FP)          // Store sign
	MOVD R4, exponent+24(FP)      // Store exponent as int64
	MOVD R6, mantissa+32(FP)      // Store mantissa

	RET

// Phase 3: Construct float64 directly from IEEE 754 components
// This avoids Go-level bit manipulation and is faster
// func constructFloat64FromIEEE754ARM64(sign uint64, exponent int64, mantissa uint64) float64
TEXT ·constructFloat64FromIEEE754ARM64(SB), NOSPLIT, $0-24
	MOVD sign+0(FP), R0           // R0 = sign (0 or 1)
	MOVD exponent+8(FP), R1       // R1 = exponent (int64)
	MOVD mantissa+16(FP), R2       // R2 = mantissa (uint64)

	// Reconstruct IEEE 754 double from components
	// Format: sign (bit 63) | exponent (bits 52-62) | mantissa (bits 0-51)
	
	// Shift sign to bit 63
	LSL $63, R0, R0               // R0 = sign << 63

	// Shift exponent to bits 52-62 and mask
	MOVD R1, R3                   // R3 = exponent
	LSL $52, R3, R3               // R3 = exponent << 52
	MOVD $0x7FF0000000000000, R4  // Load mask for exponent bits
	AND R4, R3, R3                // Mask to 11 bits

	// Mantissa is already in correct position (bits 0-51)
	MOVD $0xFFFFFFFFFFFFF, R5     // Load mask for mantissa
	AND R5, R2, R2                // Mask mantissa to 52 bits (ensure clean)

	// Combine: sign | exponent | mantissa
	ORR R3, R0, R0                // R0 = sign | exponent
	ORR R2, R0, R0                // R0 = sign | exponent | mantissa

	// Return as float64 (same bit pattern)
	MOVD R0, ret+24(FP)           // Store as uint64 (will be interpreted as float64)
	RET

