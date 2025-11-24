// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

// Optimized inverse trigonometric functions with reduced allocations

// atanWorkspace holds pre-allocated buffers for arctangent calculations
type atanWorkspace struct {
	result      *BigFloat
	term        *BigFloat
	xSquared    *BigFloat
	threshold   *BigFloat
	numerator   *BigFloat
	denominator *BigFloat
	temp        *BigFloat
	prec        uint
}

// getAtanWorkspace returns a workspace with pre-allocated buffers
func getAtanWorkspace(prec uint) *atanWorkspace {
	if prec == 0 {
		prec = DefaultPrecision
	}
	return &atanWorkspace{
		result:      NewBigFloat(0.0, prec),
		term:        NewBigFloat(0.0, prec),
		xSquared:    NewBigFloat(0.0, prec),
		threshold:   NewBigFloat(0.0, prec),
		numerator:   NewBigFloat(0.0, prec),
		denominator: NewBigFloat(0.0, prec),
		temp:        NewBigFloat(0.0, prec),
		prec:        prec,
	}
}

// bigAtanOptimized computes arctan(x) with optimized allocation pattern
func bigAtanOptimized(x *BigFloat, prec uint) *BigFloat {
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
		result := bigAtanOptimized(negX, prec)
		return result.Neg(result)
	}

	// Step 2: Inversion - if x > 1, use atan(x) = π/2 - atan(1/x)
	if x.Cmp(one) > 0 {
		invX := new(BigFloat).SetPrec(prec).Quo(one, x)
		atanInvX := bigAtanOptimized(invX, prec)
		halfPi := BigHalfPI(prec)
		return new(BigFloat).SetPrec(prec).Sub(halfPi, atanInvX)
	}

	// Special case: atan(1) = π/4 (exact)
	if x.Cmp(one) == 0 {
		return new(BigFloat).SetPrec(prec).Quo(BigPI(prec), NewBigFloat(4.0, prec))
	}

	// Step 3: Argument Reduction using the halving formula
	threshold := NewBigFloat(0.5, prec)
	reductionCount := 0
	xReduced := new(BigFloat).SetPrec(prec).Set(x)

	// Pre-allocate buffers for reduction loop
	xSquared := NewBigFloat(0.0, prec)
	onePlusXSquared := NewBigFloat(0.0, prec)
	denominator := NewBigFloat(0.0, prec)

	for xReduced.Cmp(threshold) > 0 && reductionCount < 10 {
		xSquared.Mul(xReduced, xReduced)
		onePlusXSquared.Add(one, xSquared)
		sqrtTerm := BigSqrt(onePlusXSquared, prec)
		denominator.Add(one, sqrtTerm)
		xReduced.Quo(xReduced, denominator)
		reductionCount++
	}

	// Step 4: Compute atan(xReduced) using optimized Taylor series
	result := bigAtanTaylorSeriesOptimized(xReduced, prec)

	// Step 5: Multiply by 2^reductionCount
	if reductionCount > 0 {
		powerOfTwo := float64(int64(1) << uint(reductionCount))
		multiplier := NewBigFloat(powerOfTwo, prec)
		result.Mul(result, multiplier)
	}

	return result
}

// bigAtanTaylorSeriesOptimized computes atan for small |x| using optimized Taylor series
func bigAtanTaylorSeriesOptimized(x *BigFloat, prec uint) *BigFloat {
	ws := getAtanWorkspace(prec)

	// First term is x
	ws.term.Set(x)
	ws.result.Set(ws.term)

	// Pre-compute x² once
	ws.xSquared.Mul(x, x)

	// Convergence threshold
	ws.threshold.SetMantExp(NewBigFloat(1.0, prec), -int(prec))

	// Optimized Taylor series loop
	// atan(x) = x - x³/3 + x⁵/5 - ...
	// term = term * (-x²) * (2n-1) / (2n+1)
	for n := 1; n < 500; n++ {
		// Update term: term = term * (-x²)
		ws.term.Mul(ws.term, ws.xSquared)
		ws.term.Neg(ws.term)

		// Multiply by (2n-1) / (2n+1)
		ws.numerator.SetInt64(int64(2*n - 1))
		ws.denominator.SetInt64(int64(2*n + 1))
		ws.term.Mul(ws.term, ws.numerator)
		ws.term.Quo(ws.term, ws.denominator)

		ws.result.Add(ws.result, ws.term)

		// Check convergence using pre-allocated buffer
		ws.temp.Abs(ws.term)
		if ws.temp.Cmp(ws.threshold) < 0 {
			break
		}
	}

	return ws.result
}

// bigAtan2Optimized computes atan2(y, x) with optimized allocation pattern
func bigAtan2Optimized(y, x *BigFloat, prec uint) *BigFloat {
	if prec == 0 {
		prec = y.Prec()
	}

	zero := NewBigFloat(0.0, prec)

	// Handle special cases
	if y.Sign() == 0 {
		if x.Sign() >= 0 {
			return zero
		}
		return BigPI(prec)
	}

	if x.Sign() == 0 {
		if y.Sign() > 0 {
			return BigHalfPI(prec)
		}
		halfPi := BigHalfPI(prec)
		return halfPi.Neg(halfPi)
	}

	// Compute atan(y/x)
	ratio := new(BigFloat).SetPrec(prec).Quo(y, x)
	atan := bigAtanOptimized(ratio, prec)

	// Adjust quadrant based on signs of x and y
	if x.Sign() < 0 {
		if y.Sign() >= 0 {
			atan.Add(atan, BigPI(prec))
		} else {
			atan.Sub(atan, BigPI(prec))
		}
	}

	return atan
}

// bigAsinOptimized computes arcsin(x) using the relation: asin(x) = atan(x / sqrt(1 - x²))
func bigAsinOptimized(x *BigFloat, prec uint) *BigFloat {
	if prec == 0 {
		prec = x.Prec()
	}

	one := NewBigFloat(1.0, prec)
	xSquared := new(BigFloat).SetPrec(prec).Mul(x, x)
	oneMinusXSquared := new(BigFloat).SetPrec(prec).Sub(one, xSquared)
	sqrtTerm := BigSqrt(oneMinusXSquared, prec)
	ratio := new(BigFloat).SetPrec(prec).Quo(x, sqrtTerm)
	return bigAtanOptimized(ratio, prec)
}

// bigAcosOptimized computes arccos(x) using the relation: acos(x) = π/2 - asin(x)
func bigAcosOptimized(x *BigFloat, prec uint) *BigFloat {
	if prec == 0 {
		prec = x.Prec()
	}

	halfPi := BigHalfPI(prec)
	asin := bigAsinOptimized(x, prec)
	return new(BigFloat).SetPrec(prec).Sub(halfPi, asin)
}
