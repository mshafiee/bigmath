// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

import (
	"math"
)

// Generic pure-Go implementations of trigonometric functions
// These serve as reference implementations and fallbacks

// bigSinGeneric computes sin(x) using Taylor series with arbitrary precision (pure-Go)
func bigSinGeneric(x *BigFloat, prec uint) *BigFloat {
	if prec == 0 {
		prec = x.Prec()
	}

	// Normalize x to [-π, π] using high-precision PI
	x = normalizeAngle(x, prec)

	// Taylor series computation
	// sin(x) = x - x^3/3! + ...
	// We use a slightly higher working precision for intermediate steps
	workPrec := prec + 16

	result := new(BigFloat).SetPrec(workPrec)
	term := new(BigFloat).SetPrec(workPrec).Set(x) // First term is x
	result.Set(term)

	xSquared := new(BigFloat).SetPrec(workPrec)
	xSquared.Mul(x, x)

	// Convergence threshold
	threshold := new(BigFloat).SetPrec(workPrec).SetMantExp(NewBigFloat(1.0, workPrec), -int(prec))

	for n := 1; n < 200; n++ {
		// term = term * (-x²) / ((2n)(2n+1))
		denominator1 := new(BigFloat).SetPrec(workPrec).SetInt64(int64(2 * n))
		denominator2 := new(BigFloat).SetPrec(workPrec).SetInt64(int64(2*n + 1))

		term.Mul(term, xSquared)
		term.Neg(term)
		term.Quo(term, denominator1)
		term.Quo(term, denominator2)

		result.Add(result, term)

		// Check convergence
		if BigAbs(term, workPrec).Cmp(threshold) < 0 {
			break
		}
	}

	return new(BigFloat).SetPrec(prec).Set(result)
}

// bigCosGeneric computes cos(x) using Taylor series with arbitrary precision (pure-Go)
func bigCosGeneric(x *BigFloat, prec uint) *BigFloat {
	if prec == 0 {
		prec = x.Prec()
	}

	// Normalize x to [-π, π]
	x = normalizeAngle(x, prec)

	workPrec := prec + 16

	// Taylor series computation
	result := NewBigFloat(1.0, workPrec) // First term is 1
	term := NewBigFloat(1.0, workPrec)

	xSquared := new(BigFloat).SetPrec(workPrec)
	xSquared.Mul(x, x)

	threshold := new(BigFloat).SetPrec(workPrec).SetMantExp(NewBigFloat(1.0, workPrec), -int(prec))

	for n := 1; n < 200; n++ {
		// term = term * (-x²) / ((2n-1)(2n))
		denominator1 := new(BigFloat).SetPrec(workPrec).SetInt64(int64(2*n - 1))
		denominator2 := new(BigFloat).SetPrec(workPrec).SetInt64(int64(2 * n))

		term.Mul(term, xSquared)
		term.Neg(term)
		term.Quo(term, denominator1)
		term.Quo(term, denominator2)

		result.Add(result, term)

		// Check convergence
		if BigAbs(term, workPrec).Cmp(threshold) < 0 {
			break
		}
	}

	return new(BigFloat).SetPrec(prec).Set(result)
}

// bigTanGeneric computes tan(x) = sin(x) / cos(x) (pure-Go)
func bigTanGeneric(x *BigFloat, prec uint) *BigFloat {
	if prec == 0 {
		prec = x.Prec()
	}

	sinX := BigSin(x, prec)
	cosX := BigCos(x, prec)

	result := new(BigFloat).SetPrec(prec)
	result.Quo(sinX, cosX)

	return result
}

//nolint:unused // Used in dispatch system
func bigAtanGeneric(x *BigFloat, prec uint) *BigFloat {
	if prec == 0 {
		prec = x.Prec()
	}

	// Handle zero
	zero := NewBigFloat(0.0, prec)
	if x.Cmp(zero) == 0 {
		return zero
	}

	one := NewBigFloat(1.0, prec)

	// Step 1: Symmetry - if x < 0, return -atan(-x)
	if x.Sign() < 0 {
		negX := new(BigFloat).SetPrec(prec).Neg(x)
		result := bigAtanGeneric(negX, prec)
		return result.Neg(result)
	}

	// Step 2: Inversion - if x > 1, use atan(x) = π/2 - atan(1/x)
	if x.Cmp(one) > 0 {
		invX := new(BigFloat).SetPrec(prec).Quo(one, x)
		atanInvX := bigAtanGeneric(invX, prec)
		halfPi := BigHalfPI(prec)
		return new(BigFloat).SetPrec(prec).Sub(halfPi, atanInvX)
	}

	// Special case: atan(1) = π/4 (exact)
	if x.Cmp(one) == 0 {
		return new(BigFloat).SetPrec(prec).Quo(BigPI(prec), NewBigFloat(4.0, prec))
	}

	// Step 3: Argument Reduction using the halving formula
	// atan(x) = 2*atan(x/(1+sqrt(1+x²)))
	threshold := NewBigFloat(0.5, prec)
	reductionCount := 0
	xReduced := new(BigFloat).SetPrec(prec).Set(x)

	for xReduced.Cmp(threshold) > 0 && reductionCount < 10 {
		xSquared := new(BigFloat).SetPrec(prec).Mul(xReduced, xReduced)
		onePlusXSquared := new(BigFloat).SetPrec(prec).Add(one, xSquared)
		sqrtTerm := BigSqrt(onePlusXSquared, prec)
		denominator := new(BigFloat).SetPrec(prec).Add(one, sqrtTerm)
		xReduced = new(BigFloat).SetPrec(prec).Quo(xReduced, denominator)
		reductionCount++
	}

	// Step 4: Compute atan(xReduced) using Taylor series
	result := bigAtanTaylorSeries(xReduced, prec)

	// Step 5: Multiply by 2^reductionCount
	if reductionCount > 0 {
		powerOfTwo := float64(int64(1) << uint(reductionCount))
		multiplier := NewBigFloat(powerOfTwo, prec)
		result.Mul(result, multiplier)
	}

	return result
}

//nolint:unused // Used internally by bigAtanGeneric
func bigAtanTaylorSeries(x *BigFloat, prec uint) *BigFloat {
	result := new(BigFloat).SetPrec(prec)
	term := new(BigFloat).SetPrec(prec).Set(x)
	result.Set(term)

	xSquared := new(BigFloat).SetPrec(prec).Mul(x, x)

	threshold := new(BigFloat).SetPrec(prec).SetMantExp(NewBigFloat(1.0, prec), -int(prec))

	for n := 1; n < 500; n++ {
		term.Mul(term, xSquared)
		term.Neg(term)

		numerator := new(BigFloat).SetPrec(prec).SetInt64(int64(2*n - 1))
		denominator := new(BigFloat).SetPrec(prec).SetInt64(int64(2*n + 1))

		term.Mul(term, numerator)
		term.Quo(term, denominator)

		result.Add(result, term)

		if BigAbs(term, prec).Cmp(threshold) < 0 {
			break
		}
	}

	return result
}

//nolint:unused // Used in dispatch system
func bigAtan2Generic(y, x *BigFloat, prec uint) *BigFloat {
	if prec == 0 {
		prec = y.Prec()
	}

	zero := NewBigFloat(0.0, prec)
	pi := BigPI(prec)
	halfPi := BigHalfPI(prec)

	if x.Cmp(zero) == 0 && y.Cmp(zero) == 0 {
		return NewBigFloat(0.0, prec)
	}

	if x.Cmp(zero) == 0 {
		if y.Sign() > 0 {
			return new(BigFloat).SetPrec(prec).Set(halfPi)
		}
		result := new(BigFloat).SetPrec(prec).Set(halfPi)
		return result.Neg(result)
	}

	if y.Cmp(zero) == 0 {
		if x.Sign() > 0 {
			return NewBigFloat(0.0, prec)
		}
		return new(BigFloat).SetPrec(prec).Set(pi)
	}

	ratio := new(BigFloat).SetPrec(prec).Quo(y, x)
	atan := BigAtan(ratio, prec)

	if x.Sign() > 0 {
		return atan
	}

	if y.Sign() >= 0 {
		return new(BigFloat).SetPrec(prec).Add(atan, pi)
	}

	return new(BigFloat).SetPrec(prec).Sub(atan, pi)
}

//nolint:unused // Used in dispatch system
func bigAsinGeneric(x *BigFloat, prec uint) *BigFloat {
	if prec == 0 {
		prec = x.Prec()
	}

	one := NewBigFloat(1.0, prec)
	absX := BigAbs(x, prec)

	if absX.Cmp(one) > 0 {
		return NewBigFloat(math.NaN(), prec)
	}

	if x.Sign() == 0 {
		return NewBigFloat(0.0, prec)
	}

	if absX.Cmp(one) == 0 {
		halfPi := BigHalfPI(prec)
		if x.Sign() < 0 {
			halfPi.Neg(halfPi)
		}
		return halfPi
	}

	// Check if x² is very close to 0.5 (for pi/4 case)
	xSquared := new(BigFloat).SetPrec(prec).Mul(x, x)
	half := NewBigFloat(0.5, prec)
	diff := new(BigFloat).SetPrec(prec).Sub(xSquared, half)
	absDiff := BigAbs(diff, prec)
	tolerance := new(BigFloat).SetPrec(prec).SetFloat64(1e-70)

	if absDiff.Cmp(tolerance) < 0 {
		quarterPi := new(BigFloat).SetPrec(prec).Quo(BigPI(prec), NewBigFloat(4.0, prec))
		if x.Sign() < 0 {
			quarterPi.Neg(quarterPi)
		}
		return quarterPi
	}

	// asin(x) = atan(x / sqrt(1 - x²))
	// Reuse xSquared
	oneMinusXSquared := new(BigFloat).SetPrec(prec)
	oneMinusXSquared.Sub(one, xSquared)

	sqrtTerm := BigSqrt(oneMinusXSquared, prec)

	ratio := new(BigFloat).SetPrec(prec)
	ratio.Quo(x, sqrtTerm)

	return BigAtan(ratio, prec)
}

//nolint:unused // Used in dispatch system
func bigAcosGeneric(x *BigFloat, prec uint) *BigFloat {
	if prec == 0 {
		prec = x.Prec()
	}

	asinX := BigAsin(x, prec)
	halfPi := BigHalfPI(prec)

	result := new(BigFloat).SetPrec(prec)
	result.Sub(halfPi, asinX)

	return result
}

// normalizeAngle normalizes an angle to the range [-π, π]
func normalizeAngle(x *BigFloat, prec uint) *BigFloat {
	if prec == 0 {
		prec = x.Prec()
	}

	twoPi := BigTwoPI(prec)
	pi := BigPI(prec)
	negPi := new(BigFloat).SetPrec(prec).Neg(pi)

	// Quick check if already in range
	if x.Cmp(negPi) >= 0 && x.Cmp(pi) <= 0 {
		return new(BigFloat).SetPrec(prec).Set(x)
	}

	// Reduce x modulo 2π using extended precision for intermediate calculation
	// This minimizes rounding errors in the subtraction
	extPrec := prec + 64

	result := new(BigFloat).SetPrec(extPrec).Set(x)
	twoPiExt := new(BigFloat).SetPrec(extPrec).Set(twoPi)
	piExt := new(BigFloat).SetPrec(extPrec).Set(pi)

	// Number of full rotations
	nRotations := new(BigFloat).SetPrec(extPrec)
	nRotations.Quo(result, twoPiExt)

	// Get integer part (floor for positive, ceiling-1 for negative)
	nRotationsInt, _ := nRotations.Int(nil)
	nRotationsBig := new(BigFloat).SetPrec(extPrec).SetInt(nRotationsInt)

	// Subtract full rotations
	temp := new(BigFloat).SetPrec(extPrec)
	temp.Mul(nRotationsBig, twoPiExt)
	result.Sub(result, temp)

	// Normalize to [-π, π]
	negPiExt := new(BigFloat).SetPrec(extPrec).Neg(piExt)
	if result.Cmp(piExt) > 0 {
		result.Sub(result, twoPiExt)
	} else if result.Cmp(negPiExt) < 0 {
		result.Add(result, twoPiExt)
	}

	// Round back to target precision
	return new(BigFloat).SetPrec(prec).Set(result)
}
