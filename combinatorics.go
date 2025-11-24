// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

// BigFactorial computes n! (factorial) using Gamma function
// n! = Γ(n+1)
func BigFactorial(n int64, prec uint) *BigFloat {
	return getDispatcher().BigFactorialImpl(n, prec)
}

// BigBinomial computes the binomial coefficient C(n, k) = n! / (k! * (n-k)!)
// Uses optimized computation to avoid computing large factorials
func BigBinomial(n, k int64, prec uint) *BigFloat {
	return getDispatcher().BigBinomialImpl(n, k, prec)
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
