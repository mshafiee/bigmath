// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

import (
	"math"
)

// Generic implementations for combinatorics functions (used as fallback)

// bigFactorialGeneric computes n! (factorial) using pure Go implementation
func bigFactorialGeneric(n int64, prec uint) *BigFloat {
	if prec == 0 {
		prec = DefaultPrecision
	}

	// Handle special cases
	if n < 0 {
		return NewBigFloat(math.NaN(), prec)
	}
	if n == 0 || n == 1 {
		return NewBigFloat(1.0, prec)
	}

	// For small n, compute directly for better performance
	if n <= 20 {
		result := NewBigFloat(1.0, prec)
		for i := int64(2); i <= n; i++ {
			result.Mul(result, NewBigFloat(float64(i), prec))
		}
		return result
	}

	// For large n, use Gamma function: n! = Γ(n+1)
	nPlusOne := NewBigFloat(float64(n+1), prec)
	return BigGamma(nPlusOne, prec)
}

// bigBinomialGeneric computes the binomial coefficient C(n, k) using pure Go implementation
func bigBinomialGeneric(n, k int64, prec uint) *BigFloat {
	if prec == 0 {
		prec = DefaultPrecision
	}

	// Handle special cases
	if k < 0 || k > n {
		return NewBigFloat(0.0, prec)
	}
	if k == 0 || k == n {
		return NewBigFloat(1.0, prec)
	}

	// Use symmetry: C(n, k) = C(n, n-k)
	if k > n-k {
		k = n - k
	}

	// For small values, compute directly
	if n <= 20 {
		return bigBinomialDirect(n, k, prec)
	}

	// For larger values, compute incrementally
	result := NewBigFloat(1.0, prec)

	for i := int64(0); i < k; i++ {
		numerator := NewBigFloat(float64(n-i), prec)
		denominator := NewBigFloat(float64(k-i), prec)
		term := new(BigFloat).SetPrec(prec).Quo(numerator, denominator)
		result.Mul(result, term)
	}

	return result
}
