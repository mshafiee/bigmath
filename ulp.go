// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

import (
	"math"
)

// Ulp computes the Unit in the Last Place for a BigFloat x.
// For x = m * 2^e with precision p, ulp(x) = 2^(e-p).
func Ulp(x *BigFloat, prec uint) *BigFloat {
	if x.Sign() == 0 {
		// For zero, ULP is the smallest representable number > 0
		// which is 2^(MinExp - prec) roughly, but practically 0 for error bounds
		return NewBigFloat(0.0, prec)
	}

	// Get exponent
	exp := x.MantExp(nil)

	ulpExp := exp - int(prec)

	res := new(BigFloat).SetPrec(prec)
	res.SetMantExp(new(BigFloat).SetFloat64(1.0), ulpExp)

	return res
}

// ErrorBound represents an error bound in ULPs or absolute value
type ErrorBound struct {
	Value *BigFloat
	IsUlp bool // If true, Value is in ULPs; if false, Value is absolute
}

// NewUlpError creates a new error bound in ULPs
func NewUlpError(ulps float64, prec uint) ErrorBound {
	return ErrorBound{
		Value: NewBigFloat(ulps, prec),
		IsUlp: true,
	}
}

// NewAbsError creates a new absolute error bound
func NewAbsError(absVal *BigFloat, prec uint) ErrorBound {
	return ErrorBound{
		Value: new(BigFloat).SetPrec(prec).Set(absVal),
		IsUlp: false,
	}
}

// ToAbs converts the error bound to an absolute value given the result x
func (e ErrorBound) ToAbs(x *BigFloat, prec uint) *BigFloat {
	if !e.IsUlp {
		return new(BigFloat).SetPrec(prec).Set(e.Value)
	}

	// Convert ULPs to absolute error: error = ulps * ulp(x)
	ulp := Ulp(x, prec)
	res := new(BigFloat).SetPrec(prec)
	res.Mul(e.Value, ulp)
	return res
}

// AddErrorBounds adds two error bounds (assuming they are independent)
// If both are ULPs, we can add them directly (approximation).
// If mixed, we convert to absolute.
func AddErrorBounds(e1, e2 ErrorBound, x *BigFloat, prec uint) ErrorBound {
	if e1.IsUlp && e2.IsUlp {
		sum := new(BigFloat).SetPrec(prec).Add(e1.Value, e2.Value)
		return ErrorBound{Value: sum, IsUlp: true}
	}

	abs1 := e1.ToAbs(x, prec)
	abs2 := e2.ToAbs(x, prec)
	sum := new(BigFloat).SetPrec(prec).Add(abs1, abs2)

	return ErrorBound{Value: sum, IsUlp: false}
}

// PropagateErrorAdd propagates error for addition z = x + y
// Error(z) <= Error(x) + Error(y) + RoundingError
func PropagateErrorAdd(x, y, z *BigFloat, errX, errY ErrorBound, prec uint, mode RoundingMode) ErrorBound {
	// Absolute errors add up
	absErrX := errX.ToAbs(x, prec)
	absErrY := errY.ToAbs(y, prec)

	totalAbsErr := new(BigFloat).SetPrec(prec).Add(absErrX, absErrY)

	// Add rounding error
	// Rounding error <= 0.5 ulp(z) for Nearest, 1 ulp(z) for Directed
	roundingErrUlps := 0.5
	if mode != ToNearest {
		roundingErrUlps = 1.0
	}

	roundingErr := NewUlpError(roundingErrUlps, prec).ToAbs(z, prec)
	totalAbsErr.Add(totalAbsErr, roundingErr)

	return NewAbsError(totalAbsErr, prec)
}

// PropagateErrorMul propagates error for multiplication z = x * y
// Relative errors add up roughly: RelErr(z) <= RelErr(x) + RelErr(y) + RoundingRelErr
// Since ULP error is similar to relative error (ulp(x)/x ≈ 2^-p), we can sum ULP errors.
func PropagateErrorMul(x, y, z *BigFloat, errX, errY ErrorBound, prec uint, mode RoundingMode) ErrorBound {
	// Convert everything to ULPs relative to the result z
	// This is an approximation valid for high precision

	// If errX is in ULPs of x, then AbsErrX = errX * ulp(x)
	// RelErrX = AbsErrX / |x| = errX * ulp(x) / |x| ≈ errX * 2^-p

	// RelErrZ = RelErrX + RelErrY + RoundingRelErr
	// ErrZ_ulps ≈ ErrX_ulps + ErrY_ulps + RoundingErr_ulps

	var ulpsX, ulpsY *BigFloat

	if errX.IsUlp {
		ulpsX = errX.Value
	} else {
		// Convert abs error to ulps: err / ulp(x)
		ulpX := Ulp(x, prec)
		ulpsX = new(BigFloat).SetPrec(prec).Quo(errX.Value, ulpX)
	}

	if errY.IsUlp {
		ulpsY = errY.Value
	} else {
		ulpY := Ulp(y, prec)
		ulpsY = new(BigFloat).SetPrec(prec).Quo(errY.Value, ulpY)
	}

	totalUlps := new(BigFloat).SetPrec(prec).Add(ulpsX, ulpsY)

	// Add rounding error
	roundingErrUlps := 0.5
	if mode != ToNearest {
		roundingErrUlps = 1.0
	}
	totalUlps.Add(totalUlps, NewBigFloat(roundingErrUlps, prec))

	return ErrorBound{Value: totalUlps, IsUlp: true}
}

// CalculateRequiredPrecision estimates the working precision needed to achieve
// target precision with given error bounds.
// Rule of thumb: WorkingPrec = TargetPrec + log2(AccumulatedErrorUlps) + GuardBits
func CalculateRequiredPrecision(targetPrec uint, expectedErrorUlps float64) uint {
	if expectedErrorUlps <= 1.0 {
		return targetPrec + 2 // Minimal guard bits
	}

	bits := math.Ceil(math.Log2(expectedErrorUlps))
	return targetPrec + uint(bits) + 5 // 5 extra bits for safety
}
