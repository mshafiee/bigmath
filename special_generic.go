// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

import (
	"math"
)

// Generic implementations for special functions (used as fallback)

// bigGammaGeneric computes the Gamma function using pure Go implementation
func bigGammaGeneric(x *BigFloat, prec uint) *BigFloat {
	// Call back to the helpers in special.go
	// The helpers are: bigGammaPositive, etc.
	if prec == 0 {
		prec = x.Prec()
	}

	// Handle special cases
	if x.Sign() == 0 {
		return new(BigFloat).SetPrec(prec).SetInf(false)
	}
	if x.IsInf() {
		if x.Sign() > 0 {
			return new(BigFloat).SetPrec(prec).SetInf(false)
		}
		return NewBigFloat(math.NaN(), prec)
	}

	workPrec := prec + 32

	// Check if x is negative
	if x.Sign() < 0 {
		// Use reflection formula: Γ(x) = π / (Γ(1-x) * sin(π*x))
		one := NewBigFloat(1.0, workPrec)
		oneMinusX := new(BigFloat).SetPrec(workPrec).Sub(one, x)
		gamma1MinusX := bigGammaPositive(oneMinusX, workPrec)

		piX := new(BigFloat).SetPrec(workPrec).Mul(BigPI(workPrec), x)
		sinPiX := BigSin(piX, workPrec)

		denom := new(BigFloat).SetPrec(workPrec).Mul(gamma1MinusX, sinPiX)
		result := new(BigFloat).SetPrec(workPrec).Quo(BigPI(workPrec), denom)

		return new(BigFloat).SetPrec(prec).Set(result)
	}

	return new(BigFloat).SetPrec(prec).Set(bigGammaPositive(x, workPrec))
}

// bigErfGeneric computes the error function using pure Go implementation
func bigErfGeneric(x *BigFloat, prec uint) *BigFloat {
	if prec == 0 {
		prec = x.Prec()
	}

	// Handle special cases
	if x.Sign() == 0 {
		return NewBigFloat(0.0, prec)
	}
	if x.IsInf() {
		if x.Sign() > 0 {
			return NewBigFloat(1.0, prec)
		}
		return NewBigFloat(-1.0, prec)
	}

	workPrec := prec + 32

	// For small |x|, use series expansion
	xAbs := BigAbs(x, workPrec)
	smallThreshold := NewBigFloat(0.8, workPrec)
	moderateThreshold := NewBigFloat(2.0, workPrec)

	if xAbs.Cmp(smallThreshold) < 0 {
		return bigErfSeries(x, workPrec, prec)
	}

	if xAbs.Cmp(moderateThreshold) < 0 {
		if x.Sign() > 0 {
			one := NewBigFloat(1.0, workPrec)
			erfcX := bigErfcImproved(x, workPrec, workPrec)
			result := new(BigFloat).SetPrec(workPrec).Sub(one, erfcX)
			return new(BigFloat).SetPrec(prec).Set(result)
		} else {
			negOne := NewBigFloat(-1.0, workPrec)
			negX := new(BigFloat).SetPrec(workPrec).Neg(x)
			erfcNegX := bigErfcImproved(negX, workPrec, workPrec)
			result := new(BigFloat).SetPrec(workPrec).Add(negOne, erfcNegX)
			return new(BigFloat).SetPrec(prec).Set(result)
		}
	}

	// For large |x|, use erfc
	if x.Sign() > 0 {
		one := NewBigFloat(1.0, workPrec)
		erfcX := bigErfcGeneric(x, workPrec)
		result := new(BigFloat).SetPrec(workPrec).Sub(one, erfcX)
		return new(BigFloat).SetPrec(prec).Set(result)
	} else {
		negOne := NewBigFloat(-1.0, workPrec)
		negX := new(BigFloat).SetPrec(workPrec).Neg(x)
		erfcNegX := bigErfcGeneric(negX, workPrec)
		result := new(BigFloat).SetPrec(workPrec).Add(negOne, erfcNegX)
		return new(BigFloat).SetPrec(prec).Set(result)
	}
}

// bigErfcGeneric computes the complementary error function using pure Go implementation
func bigErfcGeneric(x *BigFloat, prec uint) *BigFloat {
	if prec == 0 {
		prec = x.Prec()
	}

	// Handle special cases
	if x.Sign() == 0 {
		return NewBigFloat(1.0, prec)
	}
	if x.IsInf() {
		if x.Sign() > 0 {
			return NewBigFloat(0.0, prec)
		}
		return NewBigFloat(2.0, prec)
	}

	workPrec := prec + 32

	// For negative x, use erfc(-x) = 2 - erfc(x)
	if x.Sign() < 0 {
		negX := new(BigFloat).SetPrec(workPrec).Neg(x)
		erfcNegX := bigErfcGeneric(negX, workPrec)
		two := NewBigFloat(2.0, workPrec)
		result := new(BigFloat).SetPrec(workPrec).Sub(two, erfcNegX)
		return new(BigFloat).SetPrec(prec).Set(result)
	}

	// For small x, compute as 1 - erf(x)
	xAbs := BigAbs(x, workPrec)
	smallThreshold := NewBigFloat(0.8, workPrec)

	if xAbs.Cmp(smallThreshold) < 0 {
		one := NewBigFloat(1.0, workPrec)
		erfX := bigErfGeneric(x, workPrec)
		result := new(BigFloat).SetPrec(workPrec).Sub(one, erfX)
		return new(BigFloat).SetPrec(prec).Set(result)
	}

	// For moderate and large x, use improved erfc computation
	return bigErfcImproved(x, workPrec, prec)
}

// bigBesselJGeneric computes the Bessel function of the first kind using pure Go implementation
func bigBesselJGeneric(n int, x *BigFloat, prec uint) *BigFloat {
	if prec == 0 {
		prec = x.Prec()
	}

	// Handle special cases
	if x.Sign() == 0 {
		if n == 0 {
			return NewBigFloat(1.0, prec)
		}
		return NewBigFloat(0.0, prec)
	}
	if x.IsInf() {
		return NewBigFloat(0.0, prec)
	}

	workPrec := prec + 32

	// Handle negative n
	if n < 0 {
		result := bigBesselJGeneric(-n, x, workPrec)
		if (-n)%2 != 0 {
			result.Neg(result)
		}
		return new(BigFloat).SetPrec(prec).Set(result)
	}

	// For large x, use asymptotic expansion
	xAbs := BigAbs(x, workPrec)
	largeXThreshold := NewBigFloat(10.0, workPrec)

	if xAbs.Cmp(largeXThreshold) > 0 {
		sqrtTwoOverPiX := BigSqrt(new(BigFloat).SetPrec(workPrec).Quo(NewBigFloat(2.0, workPrec),
			new(BigFloat).SetPrec(workPrec).Mul(BigPI(workPrec), x)), workPrec)

		nPiOver2 := new(BigFloat).SetPrec(workPrec).Mul(NewBigFloat(float64(n), workPrec), BigHalfPI(workPrec))
		phase := new(BigFloat).SetPrec(workPrec).Sub(x, nPiOver2)
		phase.Sub(phase, BigHalfPI(workPrec))

		cosPhase := BigCos(phase, workPrec)
		result := new(BigFloat).SetPrec(workPrec).Mul(sqrtTwoOverPiX, cosPhase)

		if n > 0 || xAbs.Cmp(NewBigFloat(20.0, workPrec)) < 0 {
			return bigBesselJSeries(n, x, workPrec, prec)
		}

		correction := new(BigFloat).SetPrec(workPrec).Quo(
			NewBigFloat(float64(4*n*n-1), workPrec),
			new(BigFloat).SetPrec(workPrec).Mul(NewBigFloat(8.0, workPrec), x),
		)
		sinPhase := BigSin(phase, workPrec)
		corrTerm := new(BigFloat).SetPrec(workPrec).Mul(correction, sinPhase)
		result.Sub(result, corrTerm)

		return new(BigFloat).SetPrec(prec).Set(result)
	}

	// For small to moderate x, use series expansion
	return bigBesselJSeries(n, x, workPrec, prec)
}

// bigBesselYGeneric computes the Bessel function of the second kind using pure Go implementation
func bigBesselYGeneric(n int, x *BigFloat, prec uint) *BigFloat {
	if prec == 0 {
		prec = x.Prec()
	}

	// Handle special cases
	if x.Sign() <= 0 {
		return NewBigFloat(math.NaN(), prec)
	}
	if x.IsInf() {
		return NewBigFloat(0.0, prec)
	}

	workPrec := prec + 32

	if n == 0 {
		return bigBesselY0(x, workPrec, prec)
	}
	if n == 1 {
		return bigBesselY1(x, workPrec, prec)
	}

	// Use recurrence: Y_{n+1}(x) = (2n/x) * Y_n(x) - Y_{n-1}(x)
	y0 := bigBesselY0(x, workPrec, workPrec)
	y1 := bigBesselY1(x, workPrec, workPrec)

	ynMinus1 := y0
	yn := y1

	for i := 1; i < n; i++ {
		twoI := NewBigFloat(float64(2*i), workPrec)
		coeff := new(BigFloat).SetPrec(workPrec).Quo(twoI, x)
		term := new(BigFloat).SetPrec(workPrec).Mul(coeff, yn)
		ynPlus1 := new(BigFloat).SetPrec(workPrec).Sub(term, ynMinus1)

		ynMinus1 = yn
		yn = ynPlus1
	}

	return new(BigFloat).SetPrec(prec).Set(yn)
}
