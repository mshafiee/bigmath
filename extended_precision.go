// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

import (
	"math"
)

// ExtendedFloat represents an 80-bit extended precision floating-point value
// stored in the x87 FPU. This type is used internally for hardware extended
// precision operations on x86/x86-64 platforms.
//
// Note: ExtendedFloat values are not directly accessible from Go code.
// They exist only during x87 FPU operations and must be converted to/from
// BigFloat for use in Go code.
type ExtendedFloat struct {
	// value stores the 80-bit extended precision value as a float64
	// for conversion purposes. The actual x87 operations work directly
	// with the FPU stack.
	value float64
}

// BigFloatToExtendedFloat converts a BigFloat to extended precision format.
// This function prepares the value for x87 FPU operations.
// On platforms without x87 support, this returns the float64 representation.
func BigFloatToExtendedFloat(x *BigFloat) float64 {
	if x == nil {
		return 0.0
	}
	val, _ := x.Float64()
	return val
}

// ExtendedFloatToBigFloat converts an extended precision value to BigFloat.
// The value is stored from the x87 FPU stack and converted to BigFloat
// with the specified precision.
func ExtendedFloatToBigFloat(val float64, prec uint) *BigFloat {
	if prec == 0 {
		prec = DefaultPrecision
	}
	result := new(BigFloat).SetPrec(prec)
	
	// Handle special values
	if math.IsNaN(val) {
		// big.Float doesn't support NaN, return 0
		return result
	}
	if math.IsInf(val, 0) {
		result.SetInf(math.IsInf(val, -1))
		return result
	}
	
	return result.SetFloat64(val)
}

// IsExtendedPrecisionMode checks if the given precision value indicates
// extended precision mode should be used.
func IsExtendedPrecisionMode(prec uint) bool {
	return prec == ExtendedPrecision
}

// CanUseExtendedPrecision checks if extended precision mode can be used
// on the current platform with the given precision.
func CanUseExtendedPrecision(prec uint) bool {
	if !IsExtendedPrecisionMode(prec) {
		return false
	}
	features := GetCPUFeatures()
	return features.HasX87
}

