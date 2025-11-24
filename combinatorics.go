// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

import (
	"math"
)

// BigFactorial computes n! (factorial) using Gamma function
// n! = Γ(n+1)
func BigFactorial(n int64, prec uint) *BigFloat {
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

// BigBinomial computes the binomial coefficient C(n, k) = n! / (k! * (n-k)!)
// Uses optimized computation to avoid computing large factorials
func BigBinomial(n, k int64, prec uint) *BigFloat {
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

	// For larger values, use Gamma function: C(n, k) = Γ(n+1) / (Γ(k+1) * Γ(n-k+1))
	// But we can optimize by computing incrementally to avoid large intermediate values
	// C(n, k) = (n * (n-1) * ... * (n-k+1)) / (k * (k-1) * ... * 1)
	result := NewBigFloat(1.0, prec)

	for i := int64(0); i < k; i++ {
		numerator := NewBigFloat(float64(n-i), prec)
		denominator := NewBigFloat(float64(k-i), prec)
		term := new(BigFloat).SetPrec(prec).Quo(numerator, denominator)
		result.Mul(result, term)
	}

	return result
}

// bigBinomialDirect computes binomial coefficient directly for small values
func bigBinomialDirect(n, k int64, prec uint) *BigFloat {
	// C(n, k) = n! / (k! * (n-k)!)
	nFact := BigFactorial(n, prec)
	kFact := BigFactorial(k, prec)
	nMinusKFact := BigFactorial(n-k, prec)

	denom := new(BigFloat).SetPrec(prec).Mul(kFact, nMinusKFact)
	result := new(BigFloat).SetPrec(prec).Quo(nFact, denom)

	return result
}
