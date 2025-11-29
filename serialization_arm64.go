// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

//go:build arm64

package bigmath

import (
	"encoding/binary"
	"fmt"
	"io"
	"math/big"
)

// Assembly function declarations
//
// These functions are kept for backward compatibility but are no longer used
// in favor of the combined extractIEEE754FromBytesARM64 function
//
//nolint:unused // Kept for API compatibility
//go:noescape
func extractIEEE754ComponentsARM64(bits uint64) (sign uint64, exponent int64, mantissa uint64)

//nolint:unused // Kept for API compatibility
//go:noescape
func convertEndiannessBytesARM64(bytes *[8]byte, bigEndian uint8) uint64

// Combined function: extract IEEE 754 components directly from bytes with endianness conversion
// This reduces function call overhead by combining two operations into one
//
//nolint:unused // Kept for API compatibility but now using binary package for endianness
//go:noescape
func extractIEEE754FromBytesARM64(bytes *[8]byte, bigEndian uint8) (sign uint64, exponent int64, mantissa uint64)

// Phase 3: Assembly function to construct float64 directly from IEEE 754 components
// This avoids Go-level bit manipulation overhead
//
//go:noescape
func constructFloat64FromIEEE754ARM64(sign uint64, exponent int64, mantissa uint64) float64

// Cached constants for performance optimization
// These are initialized once and reused to avoid repeated allocations
var (
	two52ConstARM64 *big.Float // 2^52, cached for mantissa division
	oneConstARM64   *big.Float // 1.0, cached for mantissa addition
	// Phase 5: Pre-computed common values for fast paths
	cachedOneARM64  *big.Float // 1.0, cached for common value optimization
	cachedZeroARM64 *big.Float // 0.0, cached for common value optimization
)

func init() {
	// Initialize constants with maximum precision to support all use cases
	// They will be used with SetPrec before operations to match requested precision
	two52ConstARM64 = new(big.Float).SetUint64(1 << 52)
	oneConstARM64 = new(big.Float).SetUint64(1)
	// Pre-compute common values
	cachedOneARM64 = new(big.Float).SetUint64(1)
	cachedZeroARM64 = new(big.Float).SetUint64(0)
}

// handleZeroOrDenormalizedARM64 handles zero and denormalized numbers
func handleZeroOrDenormalizedARM64(sign bool, prec uint) *BigFloat {
	result := new(big.Float).SetPrec(prec)
	result.Set(cachedZeroARM64)
	if sign {
		result.Neg(result)
	}
	return result
}

// handleInfinityOrNaNARM64 handles infinity and NaN cases
func handleInfinityOrNaNARM64(sign bool, mantissa uint64, prec uint) *BigFloat {
	if mantissa == 0 {
		// Infinity - fast path: single operation
		result := new(big.Float).SetPrec(prec)
		result.SetInf(sign)
		return result
	}
	// NaN - use pre-computed zero
	result := new(big.Float).SetPrec(prec)
	result.Set(cachedZeroARM64) // big.Float doesn't have NaN, so we'll return zero
	return result
}

// handleCommonValuesARM64 handles fast paths for common values like 1.0 and -1.0
func handleCommonValuesARM64(exponent int, mantissa uint64, sign bool, prec uint) *BigFloat {
	if exponent == 1023 && mantissa == 0 {
		result := new(big.Float).SetPrec(prec)
		result.Set(cachedOneARM64)
		if sign {
			result.Neg(result)
		}
		return result
	}
	return nil
}

// handleNormalizedFastPathARM64 handles normalized numbers using fast float64 path
//
//nolint:unparam // sign parameter kept for consistency with other handlers
func handleNormalizedFastPathARM64(sign bool, signUint uint64, exponentInt int64, mantissaUint uint64, expValue int, prec uint) *BigFloat {
	if expValue >= -1022 && expValue <= 1023 {
		// Phase 3: Use assembly to construct float64 directly from components
		floatValue := constructFloat64FromIEEE754ARM64(signUint, exponentInt, mantissaUint)
		result := new(big.Float).SetPrec(prec)
		result.SetFloat64(floatValue)
		return result
	}
	return nil
}

// handleNormalizedExactARM64 handles normalized numbers using exact BigFloat arithmetic
func handleNormalizedExactARM64(sign bool, mantissa uint64, expValue int, prec uint) *BigFloat {
	result := new(big.Float).SetPrec(prec)
	result.SetUint64(mantissa)

	// Reuse single temporary variable for both two52 and one operations
	temp := new(big.Float).SetPrec(prec)
	temp.Set(two52ConstARM64)
	result.Quo(result, temp)

	temp.Set(oneConstARM64)
	result.Add(result, temp)

	// Extract mantissa and construct result
	mantExp := result.MantExp(result)
	result.SetMantExp(result, expValue+mantExp)

	if sign {
		result.Neg(result)
	}

	return result
}

// readDoubleAsBigFloatAsm is the assembly-optimized version of ReadDoubleAsBigFloat
func readDoubleAsBigFloatAsm(r io.Reader, bigEndian bool, prec uint) (*BigFloat, error) {
	if prec == 0 {
		prec = DefaultPrecision
	}

	// Read 8 bytes
	var doubleBytes [8]byte
	if _, err := io.ReadFull(r, doubleBytes[:]); err != nil {
		return nil, fmt.Errorf("failed to read 8 bytes: %w", err)
	}

	// Extract IEEE 754 components with correct endianness handling
	// Use binary package to correctly interpret bytes based on endianness
	// This matches the generic implementation and fixes the endianness bug
	var bits uint64
	if bigEndian {
		bits = binary.BigEndian.Uint64(doubleBytes[:])
	} else {
		bits = binary.LittleEndian.Uint64(doubleBytes[:])
	}

	// Extract components from bits
	sign := (bits >> 63) != 0
	exponent := int((bits >> 52) & 0x7FF)
	mantissa := bits & 0xFFFFFFFFFFFFF

	// Convert to types needed for assembly functions
	var signUint uint64
	if sign {
		signUint = 1
	}
	exponentInt := int64(exponent)
	mantissaUint := mantissa

	// Handle special cases
	if exponent == 0 {
		return handleZeroOrDenormalizedARM64(sign, prec), nil
	}

	if exponent == 0x7FF {
		return handleInfinityOrNaNARM64(sign, mantissa, prec), nil
	}

	// Handle common values
	if result := handleCommonValuesARM64(exponent, mantissa, sign, prec); result != nil {
		return result, nil
	}

	// Handle normalized numbers
	expValue := exponent - 1023

	// Try fast path first
	if result := handleNormalizedFastPathARM64(sign, signUint, exponentInt, mantissaUint, expValue, prec); result != nil {
		return result, nil
	}

	// Fall back to exact method
	return handleNormalizedExactARM64(sign, mantissa, expValue, prec), nil
}

// readDoubleAsBigFloatImpl dispatches to the assembly-optimized version
func readDoubleAsBigFloatImpl(r io.Reader, bigEndian bool, prec uint) (*BigFloat, error) {
	return readDoubleAsBigFloatAsm(r, bigEndian, prec)
}
