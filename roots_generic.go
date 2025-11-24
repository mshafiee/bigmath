// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

import (
	"math"
)

// Generic implementations for root functions (used as fallback)

// bigCbrtGeneric computes the cube root using pure Go implementation
func bigCbrtGeneric(x *BigFloat, prec uint) *BigFloat {
	if prec == 0 {
		prec = x.Prec()
	}

	// Handle special cases
	if x.Sign() == 0 {
		return NewBigFloat(0.0, prec)
	}
	if x.Sign() < 0 {
		// For negative numbers, compute cube root of absolute value and negate
		absX := BigAbs(x, prec)
		result := bigCbrtPositive(absX, prec)
		result.Neg(result)
		return result
	}

	return bigCbrtPositive(x, prec)
}

// bigRootGeneric computes the nth root using pure Go implementation
func bigRootGeneric(n, x *BigFloat, prec uint) *BigFloat {
	if prec == 0 {
		prec = x.Prec()
	}

	// Handle special cases
	if n.Sign() <= 0 {
		return NewBigFloat(math.NaN(), prec)
	}

	if x.Sign() == 0 {
		return NewBigFloat(0.0, prec)
	}

	if x.Sign() < 0 {
		// For even roots of negative numbers, return NaN
		// For odd roots, we could handle it
		if n.IsInt() {
			nInt, _ := n.Int64()
			if nInt%2 != 0 {
				// Odd root of negative number
				absX := BigAbs(x, prec)
				result := bigRootPositive(n, absX, prec)
				result.Neg(result)
				return result
			}
		}
		return NewBigFloat(math.NaN(), prec)
	}

	return bigRootPositive(n, x, prec)
}
