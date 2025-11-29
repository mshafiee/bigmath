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

// Combined function: extract IEEE 754 components directly from bytes with endianness conversion
// This reduces function call overhead by combining two operations into one
// Phase 3 & 4: Optimized with conditional moves and better register usage
// func extractIEEE754FromBytes(bytes *[8]byte, bigEndian uint8) (sign uint64, exponent int64, mantissa uint64)
TEXT ·extractIEEE754FromBytes(SB), NOSPLIT, $0-40
	MOVQ bytes+0(FP), SI         // SI = bytes pointer
	MOVB bigEndian+8(FP), AL      // AL = bigEndian flag (0 or 1)

	// Load 8 bytes as uint64
	MOVQ (SI), AX                 // AX = bytes as uint64 (little-endian)

	// Phase 3: Use conditional move instead of branch to avoid misprediction
	MOVQ AX, BX                   // BX = copy of original
	BSWAPQ BX                     // BX = swapped version
	TESTB AL, AL                  // Test if bigEndian != 0
	CMOVQNE BX, AX                // If bigEndian != 0, use swapped version (CMOVQNE = conditional move if not equal)

	// Phase 4: Optimize register usage - extract all components in parallel, keep in registers
	// Extract sign (bit 63)
	MOVQ AX, BX                   // BX = bits
	SHRQ $63, BX                  // BX = sign (0 or 1)

	// Extract exponent (bits 52-62, 11 bits)
	MOVQ AX, CX                   // CX = bits
	SHRQ $52, CX                  // Shift right 52 bits
	ANDQ $0x7FF, CX               // Mask to get 11 bits (0x7FF = 2047)

	// Extract mantissa (bits 0-51, 52 bits)
	MOVQ AX, DX                   // DX = bits
	ANDQ $0xFFFFFFFFFFFFF, DX     // Mask to get lower 52 bits

	// Phase 4: Write all results in sequence for better cache locality
	MOVQ BX, sign+16(FP)          // Store sign
	MOVQ CX, exponent+24(FP)      // Store exponent as int64 (sign-extended)
	MOVQ DX, mantissa+32(FP)      // Store mantissa

	RET

// Phase 3: Construct float64 directly from IEEE 754 components
// This avoids Go-level bit manipulation and is faster
// func constructFloat64FromIEEE754(sign uint64, exponent int64, mantissa uint64) float64
TEXT ·constructFloat64FromIEEE754(SB), NOSPLIT, $0-24
	MOVQ sign+0(FP), AX           // AX = sign (0 or 1)
	MOVQ exponent+8(FP), BX       // BX = exponent (int64)
	MOVQ mantissa+16(FP), CX       // CX = mantissa (uint64)

	// Reconstruct IEEE 754 double from components
	// Format: sign (bit 63) | exponent (bits 52-62) | mantissa (bits 0-51)
	
	// Shift sign to bit 63
	SHLQ $63, AX                  // AX = sign << 63

	// Shift exponent to bits 52-62 and mask
	MOVQ BX, DX                   // DX = exponent
	SHLQ $52, DX                  // DX = exponent << 52
	ANDQ $0x7FF0000000000000, DX  // Mask to 11 bits (0x7FF = 2047)

	// Mantissa is already in correct position (bits 0-51)
	ANDQ $0xFFFFFFFFFFFFF, CX     // Mask mantissa to 52 bits (ensure clean)

	// Combine: sign | exponent | mantissa
	ORQ DX, AX                    // AX = sign | exponent
	ORQ CX, AX                    // AX = sign | exponent | mantissa

	// Return as float64 (same bit pattern, stored in XMM0 for float64)
	MOVQ AX, ret+24(FP)           // Store as uint64 (will be interpreted as float64)
	RET

