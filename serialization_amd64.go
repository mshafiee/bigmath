// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

//go:build amd64

package bigmath

import (
	"fmt"
	"io"
	"math"
	"math/big"
)

// Assembly function declarations
//
// These functions are kept for backward compatibility but are no longer used
// in favor of the combined extractIEEE754FromBytes function
//
//nolint:unused // Kept for API compatibility
//go:noescape
func extractIEEE754Components(bits uint64) (sign uint64, exponent int64, mantissa uint64)

//nolint:unused // Kept for API compatibility
//go:noescape
func convertEndiannessBytes(bytes *[8]byte, bigEndian uint8) uint64

// Combined function: extract IEEE 754 components directly from bytes with endianness conversion
// This reduces function call overhead by combining two operations into one
//
//go:noescape
func extractIEEE754FromBytes(bytes *[8]byte, bigEndian uint8) (sign uint64, exponent int64, mantissa uint64)

// Phase 3: Assembly function to construct float64 directly from IEEE 754 components
// This avoids Go-level bit manipulation overhead
//
//go:noescape
func constructFloat64FromIEEE754(sign uint64, exponent int64, mantissa uint64) float64

// Cached constants for performance optimization
// These are initialized once and reused to avoid repeated allocations
var (
	two52Const *big.Float // 2^52, cached for mantissa division
	oneConst   *big.Float // 1.0, cached for mantissa addition
	// Phase 5: Pre-computed common values for fast paths
	cachedOne  *big.Float // 1.0, cached for common value optimization
	cachedZero *big.Float // 0.0, cached for common value optimization
)

func init() {
	// Initialize constants with maximum precision to support all use cases
	// They will be used with SetPrec before operations to match requested precision
	two52Const = new(big.Float).SetUint64(1 << 52)
	oneConst = new(big.Float).SetUint64(1)
	// Pre-compute common values
	cachedOne = new(big.Float).SetUint64(1)
	cachedZero = new(big.Float).SetUint64(0)
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

	// Optimized: Use combined assembly function to do endianness conversion + extraction in one call
	// This reduces function call overhead from 2 calls to 1
	var bigEndianFlag uint8
	if bigEndian {
		bigEndianFlag = 1
	}
	signUint, exponentInt, mantissaUint := extractIEEE754FromBytes(&doubleBytes, bigEndianFlag)
	sign := signUint != 0
	exponent := int(exponentInt)
	mantissa := mantissaUint
	// Keep signUint, exponentInt, mantissaUint for Phase 3 assembly function

	// Handle special cases with optimized fast paths
	if exponent == 0 {
		if mantissa == 0 {
			// Zero (positive or negative) - Phase 5: use pre-computed zero
			result := new(big.Float).SetPrec(prec)
			result.Set(cachedZero) // Use cached zero instead of SetInt64
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
		result.Set(cachedZero) // Phase 5: use pre-computed zero
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
		result.Set(cachedZero) // big.Float doesn't have NaN, so we'll return zero
		// Caller should check for NaN if needed
		return result, nil
	}

	// Normalized number
	// Value = (-1)^sign * 2^(exponent - 1023) * (1 + mantissa / 2^52)

	// Phase 5: Fast path for common value 1.0 (exponent=1023, mantissa=0, sign=0)
	if exponent == 1023 && mantissa == 0 && !sign {
		result := new(big.Float).SetPrec(prec)
		result.Set(cachedOne)
		return result, nil
	}
	// Fast path for -1.0 (exponent=1023, mantissa=0, sign=1)
	if exponent == 1023 && mantissa == 0 && sign {
		result := new(big.Float).SetPrec(prec)
		result.Set(cachedOne)
		result.Neg(result)
		return result, nil
	}

	// Phase 1 & 3: Direct Float64 construction path - fastest for normalized numbers
	// For common exponent ranges, use float64 arithmetic and SetFloat64 (highly optimized)
	// Fall back to exact method for very large/small exponents to maintain precision
	expValue := exponent - 1023

	// Use float64 fast path for common exponent ranges (-1022 to 1023)
	// This avoids expensive BigFloat arithmetic operations (Quo, Add, MantExp, SetMantExp)
	if expValue >= -1022 && expValue <= 1023 {
		// Phase 3: Use assembly to construct float64 directly from components
		// This is faster than Go-level bit manipulation
		var signUint uint64
		if sign {
			signUint = 1
		}
		floatValue := constructFloat64FromIEEE754(signUint, exponentInt, mantissaUint)

		// Use SetFloat64 which is highly optimized in Go's big.Float
		result := new(big.Float).SetPrec(prec)
		result.SetFloat64(floatValue)
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
	temp.Set(two52Const)
	result.Quo(result, temp)

	temp.Set(oneConst) // Reuse same variable instead of creating new one
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

// readDoubleAsBigFloatImpl dispatches to the assembly-optimized version
func readDoubleAsBigFloatImpl(r io.Reader, bigEndian bool, prec uint) (*BigFloat, error) {
	return readDoubleAsBigFloatAsm(r, bigEndian, prec)
}
