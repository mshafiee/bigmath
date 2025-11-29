// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

import (
	"math"
)

// BigLog computes ln(x) with specified precision using MPFR-style algorithm
// Algorithm:
//  1. Argument reduction: x = m * 2^k
//     ln(x) = ln(m) + k*ln(2)
//     where 0.5 <= m < 1 (or similar range)
//  2. Further reduction using square roots:
//     ln(m) = P * ln(m^(1/P))
//     Let y = m^(1/P) be close to 1.
//  3. Taylor series for ln(y) using atanh series:
//     ln(y) = 2 * atanh((y-1)/(y+1))
//     = 2 * sum_{n=0} (u^(2n+1) / (2n+1)) where u = (y-1)/(y+1)
//
// If prec == ExtendedPrecision and x87 is available, uses hardware extended precision.
func BigLog(x *BigFloat, prec uint) *BigFloat {
	if CanUseExtendedPrecision(prec) {
		val := BigFloatToExtendedFloat(x)
		result := extendedLog(val)
		return ExtendedFloatToBigFloat(result, prec)
	}
	// Use dispatcher to select assembly or generic implementation
	return getDispatcher().BigLogImpl(x, prec)
}

// bigLogGeneric is the generic implementation (called by dispatcher)
//
//nolint:unused // Used in dispatch system and called from assembly
func bigLogGeneric(x *BigFloat, prec uint) *BigFloat {
	if prec == 0 {
		prec = x.Prec()
	}

	// Handle special cases
	if x.Sign() <= 0 {
		// NaN for negative or zero
		return new(BigFloat).SetPrec(prec).SetFloat64(math.NaN())
	}
	if x.IsInf() {
		return new(BigFloat).SetPrec(prec).SetInf(false)
	}

	// Check for 1.0
	one := NewBigFloat(1.0, prec)
	if x.Cmp(one) == 0 {
		return NewBigFloat(0.0, prec)
	}

	workPrec := prec + 32

	// 1. Argument reduction: x = m * 2^k
	mant := new(BigFloat).SetPrec(workPrec)
	exp := x.MantExp(mant) // mant in [0.5, 1)

	// 2. Further reduction: Take square roots until close to 1
	// We want u = (y-1)/(y+1) to be small, e.g., < 2^-10
	// y = mant^(1/2^S)
	// If mant = 0.5, ln(0.5) = -0.69
	// After S sqrts, ln(y) = -0.69 / 2^S
	// y = exp(-0.69/2^S) approx 1 - 0.69/2^S
	// u approx 0.35/2^S
	// To get u < 2^-14, we need S approx 14.

	S := 14 // Sufficient for fast convergence

	y := new(BigFloat).SetPrec(workPrec).Set(mant)
	// Perform S square roots
	// We can use BigSqrt, but we need to be careful about precision loss?
	// No, precision loss is minimal if we keep extra bits.
	for i := 0; i < S; i++ {
		y = BigSqrt(y, workPrec)
	}

	// 3. Series expansion
	// u = (y-1)/(y+1)
	num := new(BigFloat).SetPrec(workPrec).Sub(y, one)
	den := new(BigFloat).SetPrec(workPrec).Add(y, one)
	u := new(BigFloat).SetPrec(workPrec).Quo(num, den)

	u2 := new(BigFloat).SetPrec(workPrec).Mul(u, u)

	res := new(BigFloat).SetPrec(workPrec).Set(u) // First term n=0
	term := new(BigFloat).SetPrec(workPrec).Set(u)

	threshold := new(BigFloat).SetPrec(workPrec).SetMantExp(NewBigFloat(1.0, workPrec), -int(workPrec))

	for n := 1; n < 1000; n++ {
		// term *= u^2
		term.Mul(term, u2)

		// Add term / (2n+1)
		denom := NewBigFloat(float64(2*n+1), workPrec)
		val := new(BigFloat).SetPrec(workPrec).Quo(term, denom)

		res.Add(res, val)

		if new(BigFloat).SetPrec(workPrec).Abs(val).Cmp(threshold) < 0 {
			break
		}
	}

	// Multiply by 2
	res.Mul(res, NewBigFloat(2.0, workPrec))

	// Multiply by 2^S (undo square roots)
	// res = res * 2^S
	scale := new(BigFloat).SetPrec(workPrec).SetInt64(1)
	scale.SetMantExp(scale, S)
	res.Mul(res, scale)

	// Add k*ln(2)
	ln2 := BigLog2(workPrec)
	kLn2 := new(BigFloat).SetPrec(workPrec).Mul(NewBigFloat(float64(exp), workPrec), ln2)

	res.Add(res, kLn2)

	return new(BigFloat).SetPrec(prec).Set(res)
}

// BigLog10 computes log10(x) = ln(x) / ln(10)
func BigLog10(x *BigFloat, prec uint) *BigFloat {
	lnX := BigLog(x, prec)
	ln10 := BigLog(NewBigFloat(10.0, prec), prec)
	return new(BigFloat).SetPrec(prec).Quo(lnX, ln10)
}
