// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

import (
	"math"
)

// BigLog1p computes log(1+x) accurately for values near zero
// Uses series expansion: log(1+x) = x - x^2/2 + x^3/3 - x^4/4 + ...
// This avoids precision loss when x is very small
func BigLog1p(x *BigFloat, prec uint) *BigFloat {
	if prec == 0 {
		prec = x.Prec()
	}

	// Handle special cases
	if x.Sign() == 0 {
		return NewBigFloat(0.0, prec)
	}

	// Check if x < -1
	negOne := NewBigFloat(-1.0, prec)
	if x.Cmp(negOne) < 0 {
		// x < -1, log(1+x) is undefined
		return NewBigFloat(math.NaN(), prec)
	}

	// Check if x == -1
	if x.Cmp(negOne) == 0 {
		// log(0) = -Inf
		return new(BigFloat).SetPrec(prec).SetInf(true)
	}

	if x.IsInf() {
		if x.Sign() > 0 {
			return new(BigFloat).SetPrec(prec).SetInf(false)
		}
		return NewBigFloat(math.NaN(), prec)
	}

	workPrec := prec + 32

	// For small |x|, use series expansion directly
	// For larger |x|, use log(1+x) = log((1+x)) computed normally
	xAbs := BigAbs(x, workPrec)
	threshold := NewBigFloat(0.1, workPrec)

	if xAbs.Cmp(threshold) < 0 {
		// Use series: log(1+x) = sum_{n=1} (-1)^(n+1) * x^n / n
		// = x - x^2/2 + x^3/3 - x^4/4 + ...
		result := new(BigFloat).SetPrec(workPrec).Set(x)
		term := new(BigFloat).SetPrec(workPrec).Set(x)
		xPower := new(BigFloat).SetPrec(workPrec).Set(x)

		threshold := new(BigFloat).SetPrec(workPrec).SetMantExp(NewBigFloat(1.0, workPrec), -int(workPrec))

		for n := 2; n < 1000; n++ {
			// Compute next term: (-1)^(n+1) * x^n / n
			xPower.Mul(xPower, x)
			term.Set(xPower)
			denom := NewBigFloat(float64(n), workPrec)
			term.Quo(term, denom)

			if n%2 == 0 {
				// Even n: subtract
				result.Sub(result, term)
			} else {
				// Odd n: add
				result.Add(result, term)
			}

			if new(BigFloat).SetPrec(workPrec).Abs(term).Cmp(threshold) < 0 {
				break
			}
		}

		return new(BigFloat).SetPrec(prec).Set(result)
	}

	// For larger |x|, compute normally: log(1+x)
	one := NewBigFloat(1.0, workPrec)
	onePlusX := new(BigFloat).SetPrec(workPrec).Add(one, x)
	return new(BigFloat).SetPrec(prec).Set(BigLog(onePlusX, workPrec))
}

// BigExp1m computes exp(x)-1 accurately for values near zero
// Uses series expansion: exp(x)-1 = x + x^2/2! + x^3/3! + ...
// This avoids precision loss when x is very small
func BigExp1m(x *BigFloat, prec uint) *BigFloat {
	if prec == 0 {
		prec = x.Prec()
	}

	// Handle special cases
	if x.Sign() == 0 {
		return NewBigFloat(0.0, prec)
	}
	if x.IsInf() {
		if x.Sign() > 0 {
			return new(BigFloat).SetPrec(prec).SetInf(false)
		}
		return NewBigFloat(-1.0, prec)
	}

	workPrec := prec + 32

	// For small |x|, use series expansion directly
	// For larger |x|, compute exp(x) - 1 normally
	xAbs := BigAbs(x, workPrec)
	threshold := NewBigFloat(0.1, workPrec)

	if xAbs.Cmp(threshold) < 0 {
		// Use series: exp(x)-1 = sum_{n=1} x^n / n!
		// = x + x^2/2! + x^3/3! + ...
		result := new(BigFloat).SetPrec(workPrec).Set(x)
		term := new(BigFloat).SetPrec(workPrec).Set(x)
		xPower := new(BigFloat).SetPrec(workPrec).Set(x)
		factorial := NewBigFloat(1.0, workPrec)

		threshold := new(BigFloat).SetPrec(workPrec).SetMantExp(NewBigFloat(1.0, workPrec), -int(workPrec))

		for n := 2; n < 1000; n++ {
			// term = x^n / n!
			xPower.Mul(xPower, x)
			factorial.Mul(factorial, NewBigFloat(float64(n), workPrec))
			term.Quo(xPower, factorial)

			result.Add(result, term)

			if new(BigFloat).SetPrec(workPrec).Abs(term).Cmp(threshold) < 0 {
				break
			}
		}

		return new(BigFloat).SetPrec(prec).Set(result)
	}

	// For larger |x|, compute normally: exp(x) - 1
	expX := BigExp(x, workPrec)
	one := NewBigFloat(1.0, workPrec)
	result := new(BigFloat).SetPrec(workPrec).Sub(expX, one)
	return new(BigFloat).SetPrec(prec).Set(result)
}

// BigLogb computes the logarithm of x with base b: log_b(x) = ln(x) / ln(b)
func BigLogb(x, base *BigFloat, prec uint) *BigFloat {
	if prec == 0 {
		prec = x.Prec()
	}

	// Handle special cases
	if base.Sign() <= 0 || base.Cmp(NewBigFloat(1.0, prec)) == 0 {
		return NewBigFloat(math.NaN(), prec)
	}
	if x.Sign() <= 0 {
		return NewBigFloat(math.NaN(), prec)
	}

	workPrec := prec + 32

	// Compute logarithm with arbitrary base using change of base formula
	lnX := BigLog(x, workPrec)
	lnB := BigLog(base, workPrec)

	result := new(BigFloat).SetPrec(workPrec).Quo(lnX, lnB)
	return new(BigFloat).SetPrec(prec).Set(result)
}
