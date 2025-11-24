// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

import "math/big"

// RoundingMode is an alias for big.RoundingMode
type RoundingMode = big.RoundingMode

// Rounding mode constants (aliases for big package constants)
const (
	ToNearest     = big.ToNearestEven // Round to nearest, ties to even
	ToNearestAway = big.ToNearestAway // Round to nearest, ties away from zero
	ToZero        = big.ToZero        // Round toward zero
	ToPositiveInf = big.ToPositiveInf // Round toward +∞
	ToNegativeInf = big.ToNegativeInf // Round toward -∞
	AwayFromZero  = big.AwayFromZero  // Round away from zero
)

// Rounding mode implementations
// These provide IEEE 754-2008 compliant rounding

//nolint:unused // May be used in dispatch or as fallback
func roundToNearestEvenGeneric(x *BigFloat, prec uint) *BigFloat {
	result := new(BigFloat).SetPrec(prec).Set(x)
	// big.Float uses round-to-nearest-even by default
	return result
}

//nolint:unused // May be used in dispatch or as fallback
func roundTowardZeroGeneric(x *BigFloat, prec uint) *BigFloat {
	result := new(BigFloat).SetPrec(prec).Set(x)
	result.SetMode(ToZero)
	return result
}

//nolint:unused // May be used in dispatch or as fallback
func roundTowardInfGeneric(x *BigFloat, prec uint) *BigFloat {
	result := new(BigFloat).SetPrec(prec).Set(x)
	result.SetMode(ToPositiveInf)
	return result
}

//nolint:unused // May be used in dispatch or as fallback
func roundTowardNegInfGeneric(x *BigFloat, prec uint) *BigFloat {
	result := new(BigFloat).SetPrec(prec).Set(x)
	result.SetMode(ToNegativeInf)
	return result
}

// Assembly function declarations

//go:noescape
//nolint:unused // Implemented in assembly (rounding_amd64.s)
func roundToNearestEvenAsm(x *BigFloat, prec uint) *BigFloat

//go:noescape
//nolint:unused // Implemented in assembly (rounding_amd64.s)
func roundTowardZeroAsm(x *BigFloat, prec uint) *BigFloat

//go:noescape
//nolint:unused // Implemented in assembly (rounding_amd64.s)
func roundTowardInfAsm(x *BigFloat, prec uint) *BigFloat

//go:noescape
//nolint:unused // Implemented in assembly (rounding_amd64.s)
func roundTowardNegInfAsm(x *BigFloat, prec uint) *BigFloat

// Round rounds x to prec bits using the specified rounding mode
// Returns the rounded value and a ternary value:
//
//	-1 if rounded down
//	 0 if exact
//	+1 if rounded up
func Round(x *BigFloat, prec uint, mode RoundingMode) (result *BigFloat, ternary int) {
	if prec == 0 {
		prec = x.Prec()
	}

	// Create result with specified precision
	result = new(BigFloat).SetPrec(prec)
	result.SetMode(mode)
	result.Set(x)

	// Determine ternary value by comparing original and rounded
	diff := new(BigFloat).SetPrec(prec).Sub(x, result)
	if diff.Sign() == 0 {
		return result, 0 // Exact
	} else if diff.Sign() > 0 {
		return result, -1 // x > result means rounded down
	} else {
		return result, 1 // x < result means rounded up
	}
}

// SqrtRounded computes sqrt(x) and rounds according to mode
func SqrtRounded(x *BigFloat, prec uint, mode RoundingMode) (result *BigFloat, ternary int) {
	if prec == 0 {
		prec = x.Prec()
	}

	// Compute sqrt with higher precision
	workPrec := prec + 32
	sqrt := BigSqrt(x, workPrec)

	// Round result
	return Round(sqrt, prec, mode)
}

// AddRounded adds two BigFloats and rounds according to mode
func AddRounded(a, b *BigFloat, prec uint, mode RoundingMode) (result *BigFloat, ternary int) {
	if prec == 0 {
		prec = a.Prec()
	}

	// Compute sum with higher precision
	workPrec := prec + 32
	sum := new(BigFloat).SetPrec(workPrec).Add(a, b)

	// Round result
	return Round(sum, prec, mode)
}

// QuoRounded divides two BigFloats and rounds according to mode
func QuoRounded(a, b *BigFloat, prec uint, mode RoundingMode) (result *BigFloat, ternary int) {
	if prec == 0 {
		prec = a.Prec()
	}

	// Compute quotient with higher precision
	workPrec := prec + 32
	quo := new(BigFloat).SetPrec(workPrec).Quo(a, b)

	// Round result
	return Round(quo, prec, mode)
}
