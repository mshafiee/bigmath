// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

//go:build !amd64 && !arm64

package bigmath

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"math/big"
)

// Cached constants for performance optimization
// These are initialized once and reused to avoid repeated allocations
var (
	two52ConstGeneric *big.Float // 2^52, cached for mantissa division
	oneConstGeneric   *big.Float // 1.0, cached for mantissa addition
	// Phase 5: Pre-computed common values for fast paths
	cachedOneGeneric  *big.Float // 1.0, cached for common value optimization
	cachedZeroGeneric *big.Float // 0.0, cached for common value optimization
)

func init() {
	// Initialize constants with maximum precision to support all use cases
	// They will be used with SetPrec before operations to match requested precision
	two52ConstGeneric = new(big.Float).SetUint64(1 << 52)
	oneConstGeneric = new(big.Float).SetUint64(1)
	// Pre-compute common values
	cachedOneGeneric = new(big.Float).SetUint64(1)
	cachedZeroGeneric = new(big.Float).SetUint64(0)
}

// readDoubleAsBigFloatImpl dispatches to the generic version
func readDoubleAsBigFloatImpl(r io.Reader, bigEndian bool, prec uint) (*BigFloat, error) {
	return readDoubleAsBigFloatGeneric(r, bigEndian, prec)
}

// readDoubleAsBigFloatGeneric is the generic (non-assembly) version of ReadDoubleAsBigFloat
func readDoubleAsBigFloatGeneric(r io.Reader, bigEndian bool, prec uint) (*BigFloat, error) {
	if prec == 0 {
		prec = DefaultPrecision
	}

	// Read 8 bytes
	doubleBytes := make([]byte, 8)
	if _, err := io.ReadFull(r, doubleBytes); err != nil {
		return nil, fmt.Errorf("failed to read 8 bytes: %w", err)
	}

	// Interpret as uint64 with correct endianness
	var bits uint64
	if bigEndian {
		bits = binary.BigEndian.Uint64(doubleBytes)
	} else {
		bits = binary.LittleEndian.Uint64(doubleBytes)
	}

	// Extract IEEE 754 components
	// Sign: bit 63
	sign := (bits >> 63) != 0

	// Exponent: bits 52-62 (11 bits)
	exponent := int((bits >> 52) & 0x7FF)

	// Mantissa: bits 0-51 (52 bits)
	mantissa := bits & 0xFFFFFFFFFFFFF

	// Handle special cases with optimized fast paths
	if exponent == 0 {
		if mantissa == 0 {
			// Zero (positive or negative) - Phase 5: use pre-computed zero
			result := new(big.Float).SetPrec(prec)
			result.Set(cachedZeroGeneric) // Use cached zero instead of SetInt64
			if sign {
				result.Neg(result)
			}
			return result, nil
		}
		// Denormalized number (subnormal)
		// For denormalized: value = (-1)^sign * 2^(-1022) * (mantissa / 2^52)
		// This is a very small number, handle as zero for now
		// TODO: Implement denormalized number handling if needed
		result := new(big.Float).SetPrec(prec)
		result.Set(cachedZeroGeneric) // Phase 5: use pre-computed zero
		return result, nil
	}

	if exponent == 0x7FF {
		// Infinity or NaN - optimized: single allocation, direct SetInf
		if mantissa == 0 {
			// Infinity - fast path: single operation
			result := new(big.Float).SetPrec(prec)
			result.SetInf(sign)
			return result, nil
		}
		// NaN - Phase 5: use pre-computed zero
		result := new(big.Float).SetPrec(prec)
		result.Set(cachedZeroGeneric) // big.Float doesn't have NaN, so we'll return zero
		// Caller should check for NaN if needed
		return result, nil
	}

	// Normalized number
	// Value = (-1)^sign * 2^(exponent - 1023) * (1 + mantissa / 2^52)

	// Phase 5: Fast path for common value 1.0 (exponent=1023, mantissa=0, sign=0)
	if exponent == 1023 && mantissa == 0 && !sign {
		result := new(big.Float).SetPrec(prec)
		result.Set(cachedOneGeneric)
		return result, nil
	}
	// Fast path for -1.0 (exponent=1023, mantissa=0, sign=1)
	if exponent == 1023 && mantissa == 0 && sign {
		result := new(big.Float).SetPrec(prec)
		result.Set(cachedOneGeneric)
		result.Neg(result)
		return result, nil
	}

	// Phase 1: Direct Float64 construction path - fastest for normalized numbers
	// For common exponent ranges, use float64 arithmetic and SetFloat64 (highly optimized)
	// Fall back to exact method for very large/small exponents to maintain precision
	expValue := exponent - 1023

	// Use float64 fast path for common exponent ranges (-1022 to 1023)
	// This avoids expensive BigFloat arithmetic operations (Quo, Add, MantExp, SetMantExp)
	if expValue >= -1022 && expValue <= 1023 {
		// Calculate float64 value directly: (-1)^sign * 2^expValue * (1 + mantissa / 2^52)
		mantissaFloat := float64(mantissa)/(1<<52) + 1.0
		value := mantissaFloat * math.Pow(2, float64(expValue))
		if sign {
			value = -value
		}

		// Use SetFloat64 which is highly optimized in Go's big.Float
		result := new(big.Float).SetPrec(prec)
		result.SetFloat64(value)
		return result, nil
	}

	// Fall back to exact method for very large/small exponents to maintain precision
	// Phase 1 & 7 optimization: Eliminate temporary allocations and reduce precision setting
	// Create result with precision first, then reuse for intermediate calculations
	result := new(big.Float).SetPrec(prec)

	// Step 1: Construct mantissa value (1 + mantissa / 2^52) using cached constants
	// Use result for mantissa construction to avoid extra SetPrec calls
	result.SetUint64(mantissa)

	// Reuse single temporary variable for both two52 and one operations
	// Phase 7: Set precision once on temp, reuse for both operations
	temp := new(big.Float).SetPrec(prec)
	temp.Set(two52ConstGeneric)
	result.Quo(result, temp)

	temp.Set(oneConstGeneric) // Reuse same variable instead of creating new one
	result.Add(result, temp)

	// Step 2: Calculate exponent and construct result directly
	// result is in range [1, 2), so MantExp will return mantExp = 1
	// Reuse result for mant extraction, then use it directly
	// This eliminates the need for a separate 'mant' variable
	mantExp := result.MantExp(result) // Extract mantissa to [0.5, 1), mantExp is 1

	// Construct result directly using SetMantExp
	// result now contains the mantissa in [0.5, 1) range
	result.SetMantExp(result, expValue+mantExp)

	// Apply sign
	if sign {
		result.Neg(result)
	}

	return result, nil
}

// readDoubleAsBigFloatImpl dispatches to the generic version
func readDoubleAsBigFloatImpl(r io.Reader, bigEndian bool, prec uint) (*BigFloat, error) {
	return readDoubleAsBigFloatGeneric(r, bigEndian, prec)
}
