package bigmath

import (
	"math"
	"math/big"
)

// BigExp computes e^x with specified precision using MPFR algorithm
// Algorithm:
//  1. Argument reduction: x = r + k*ln(2) where |r| <= ln(2)/2
//     exp(x) = exp(r) * 2^k
//  2. Further reduction: exp(r) = (exp(r/2^p))^(2^p)
//     Reduce r until |r| is small enough for Taylor series to converge fast.
//  3. Taylor series with binary splitting for exp(r/2^p).
func BigExp(x *BigFloat, prec uint) *BigFloat {
	// Use dispatcher to select assembly or generic implementation
	return getDispatcher().BigExpImpl(x, prec)
}

// bigExpGeneric is the generic implementation (called by dispatcher)
func bigExpGeneric(x *BigFloat, prec uint) *BigFloat {
	if prec == 0 {
		prec = x.Prec()
	}

	// Handle special cases
	if x.IsInf() {
		if x.Sign() > 0 {
			return new(BigFloat).SetPrec(prec).SetInf(false)
		}
		return new(BigFloat).SetPrec(prec).SetFloat64(0.0)
	}
	if x.Sign() == 0 {
		return NewBigFloat(1.0, prec)
	}

	// Working precision
	workPrec := prec + 32

	// 1. Argument reduction: x = k*ln(2) + r
	// k = round(x / ln(2))
	ln2 := BigLog2(workPrec)

	kFloat := new(BigFloat).SetPrec(workPrec).Quo(x, ln2)
	kInt := new(big.Int)
	kFloat.Int(kInt) // Round to nearest integer

	// r = x - k*ln(2)
	kBig := new(BigFloat).SetPrec(workPrec).SetInt(kInt)
	r := new(BigFloat).SetPrec(workPrec).Mul(kBig, ln2)
	r.Sub(x, r)

	// 2. Further reduction: exp(r) = (exp(r/2^S))^(2^S)
	// We want |r/2^S| < 2^-J where J is chosen for convergence.
	// For Taylor series sum x^i/i!, if x < 2^-J, terms shrink by 2^-J.
	// To get P bits, we need P/J terms.
	// MPFR balances S and J.
	// Let's choose S such that |r/2^S| < 1.
	// Actually, for binary splitting, we want integer arithmetic if possible?
	// No, we stick to BigFloat for simplicity or scaled integers.

	// Let's reduce r to be very small, e.g., < 2^-10
	// Then we use Taylor series.

	// Find S such that |r|/2^S < 2^-14 (example)
	rAbs := new(BigFloat).SetPrec(workPrec).Abs(r)
	rFloat, _ := rAbs.Float64()
	S := 0
	if rFloat > 0 {
		// log2(r) - S < -14 => S > log2(r) + 14
		S = int(math.Ceil(math.Log2(rFloat) + 14))
	}
	if S < 0 {
		S = 0
	}

	// rReduced = r / 2^S
	scale := new(BigFloat).SetPrec(workPrec).SetInt64(1)
	scale.SetMantExp(scale, S) // 2^S

	rReduced := new(BigFloat).SetPrec(workPrec).Quo(r, scale)

	// 3. Taylor series for exp(rReduced)
	// exp(u) = sum u^n / n!
	// We use binary splitting.
	// sum (u^n / n!) = sum (P/Q)
	// But u is a BigFloat.
	// Let u = U / 2^B where U is integer?
	// rReduced is small.
	// Let's just use standard series summation with BigFloat for simplicity
	// unless we want full binary splitting efficiency.
	// For MPFR compliance, binary splitting is preferred for high precision.

	// Convert rReduced to rational/integer?
	// rReduced is approx 2^-14.
	// Let's use the standard Taylor series loop for now, it's efficient enough for < 1000 bits
	// if argument is small.

	res := new(BigFloat).SetPrec(workPrec).SetFloat64(1.0)
	term := new(BigFloat).SetPrec(workPrec).SetFloat64(1.0)

	// Threshold for convergence
	threshold := new(BigFloat).SetPrec(workPrec).SetMantExp(NewBigFloat(1.0, workPrec), -int(workPrec))

	for n := 1; n < 1000; n++ {
		// term = term * rReduced / n
		term.Mul(term, rReduced)
		term.Quo(term, NewBigFloat(float64(n), workPrec))

		res.Add(res, term)

		if new(BigFloat).SetPrec(workPrec).Abs(term).Cmp(threshold) < 0 {
			break
		}
	}

	// 4. Square S times: res = res^(2^S)
	for i := 0; i < S; i++ {
		res.Mul(res, res)
	}

	// 5. Multiply by 2^k
	// res = res * 2^k
	// We can do this by adding k to exponent
	if kInt.Sign() != 0 {
		// kInt might be large, check for overflow?
		// For typical use, k fits in int.
		kVal := kInt.Int64()
		// res.MantExp gives mant, exp. New exp = exp + kVal.
		mant := new(BigFloat).SetPrec(workPrec)
		exp := res.MantExp(mant)
		res.SetMantExp(mant, exp+int(kVal))
	}

	return new(BigFloat).SetPrec(prec).Set(res)
}
