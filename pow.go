// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

import (
	"math"
)

// BigPow computes x^y with specified precision
// Uses exp(y * ln(x)) for non-integer y
// If prec == ExtendedPrecision and x87 is available, uses hardware extended precision.
func BigPow(x, y *BigFloat, prec uint) *BigFloat {
	if CanUseExtendedPrecision(prec) {
		xVal := BigFloatToExtendedFloat(x)
		yVal := BigFloatToExtendedFloat(y)
		result := extendedPow(xVal, yVal)
		return ExtendedFloatToBigFloat(result, prec)
	}
	// Use dispatcher to select assembly or generic implementation
	return getDispatcher().BigPowImpl(x, y, prec)
}

// bigPowGeneric is the generic implementation (called by dispatcher)
// This is the actual implementation - do not call BigPow from here to avoid recursion
func bigPowGeneric(x, y *BigFloat, prec uint) *BigFloat {
	if prec == 0 {
		prec = x.Prec()
	}

	// Special cases

	// x^0 = 1
	zero := NewBigFloat(0.0, prec)
	if y.Cmp(zero) == 0 {
		return NewBigFloat(1.0, prec)
	}

	// 1^y = 1
	one := NewBigFloat(1.0, prec)
	if x.Cmp(one) == 0 {
		return NewBigFloat(1.0, prec)
	}

	// x^1 = x
	if y.Cmp(one) == 0 {
		return new(BigFloat).SetPrec(prec).Set(x)
	}

	// Check if y is integer
	if y.IsInt() {
		yInt, _ := y.Int64()
		// Use integer power if y fits in int64
		// But wait, IsInt() returns true for 1.0, 2.0 etc.
		// If y is very large integer, we might still want exp/log
		// but for small integers, repeated squaring is better.
		// Let's use generic integer pow if small enough.
		if yInt >= -1000000 && yInt <= 1000000 {
			return bigPowInteger(x, yInt, prec)
		}
	}

	// x < 0
	if x.Sign() < 0 {
		// If y is integer, we can compute.
		if y.IsInt() {
			yInt, _ := y.Int64()
			absX := new(BigFloat).SetPrec(prec).Abs(x)
			// Use dispatcher directly to avoid recursion
			res := getDispatcher().BigPowImpl(absX, y, prec)
			if yInt%2 != 0 {
				res.Neg(res)
			}
			return res
		}
		// Negative base, non-integer exponent -> NaN (complex)
		return new(BigFloat).SetPrec(prec).SetFloat64(math.NaN())
	}

	// x = 0
	if x.Sign() == 0 {
		if y.Sign() > 0 {
			return NewBigFloat(0.0, prec)
		}
		// 0^-y -> Inf
		return new(BigFloat).SetPrec(prec).SetInf(false)
	}

	// General case: x^y = exp(y * ln(x))
	workPrec := prec + 32

	// Use dispatcher directly to avoid recursion
	lnX := getDispatcher().BigLogImpl(x, workPrec)
	term := new(BigFloat).SetPrec(workPrec).Mul(y, lnX)

	res := getDispatcher().BigExpImpl(term, workPrec)

	return new(BigFloat).SetPrec(prec).Set(res)
}

// bigPowInteger computes x^n for integer n
func bigPowInteger(x *BigFloat, n int64, prec uint) *BigFloat {
	if n == 0 {
		return NewBigFloat(1.0, prec)
	}

	if n < 0 {
		// x^-n = 1/x^n
		res := bigPowInteger(x, -n, prec)
		return new(BigFloat).SetPrec(prec).Quo(NewBigFloat(1.0, prec), res)
	}

	// Binary exponentiation
	res := NewBigFloat(1.0, prec)
	base := new(BigFloat).SetPrec(prec).Set(x)

	for n > 0 {
		if n%2 == 1 {
			res.Mul(res, base)
		}
		base.Mul(base, base)
		n /= 2
	}

	return res
}
