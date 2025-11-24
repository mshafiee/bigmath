// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

// Optimized trigonometric functions with reduced allocations and improved performance

// trigWorkspace holds pre-allocated buffers for trigonometric calculations
type trigWorkspace struct {
	result    *BigFloat
	term      *BigFloat
	xSquared  *BigFloat
	threshold *BigFloat
	denom1    *BigFloat
	temp      *BigFloat
	prec      uint
}

// getTrigWorkspace returns a workspace with pre-allocated buffers
func getTrigWorkspace(prec uint) *trigWorkspace {
	if prec == 0 {
		prec = DefaultPrecision
	}
	workPrec := prec + 16
	return &trigWorkspace{
		result:    NewBigFloat(0.0, workPrec),
		term:      NewBigFloat(0.0, workPrec),
		xSquared:  NewBigFloat(0.0, workPrec),
		threshold: NewBigFloat(0.0, workPrec),
		denom1:    NewBigFloat(0.0, workPrec),
		temp:      NewBigFloat(0.0, workPrec),
		prec:      prec,
	}
}

// bigSinOptimized computes sin(x) using optimized Taylor series
func bigSinOptimized(x *BigFloat, prec uint) *BigFloat {
	if prec == 0 {
		prec = x.Prec()
	}

	// Normalize x to [-π, π]
	x = normalizeAngle(x, prec)

	workPrec := prec + 16
	ws := getTrigWorkspace(prec)

	// First term is x
	ws.term.Set(x)
	ws.result.Set(ws.term)

	// Pre-compute x² once (reused in loop)
	ws.xSquared.Mul(x, x)

	// Convergence threshold
	ws.threshold.SetMantExp(NewBigFloat(1.0, workPrec), -int(prec))

	// Optimized Taylor series loop
	// sin(x) = x - x³/3! + x⁵/5! - ...
	// term = term * (-x²) / ((2n)(2n+1))
	for n := 1; n < 200; n++ {
		// Pre-compute denominator as single value: (2n)(2n+1)
		// This reduces one division operation per iteration
		denom := int64(2*n) * int64(2*n+1)
		ws.denom1.SetInt64(denom)

		// Update term: term = term * (-x²) / ((2n)(2n+1))
		ws.term.Mul(ws.term, ws.xSquared)
		ws.term.Neg(ws.term)
		ws.term.Quo(ws.term, ws.denom1) // Single division instead of two

		ws.result.Add(ws.result, ws.term)

		// Check convergence using pre-allocated temp buffer
		ws.temp.Abs(ws.term)
		if ws.temp.Cmp(ws.threshold) < 0 {
			break
		}
	}

	return new(BigFloat).SetPrec(prec).Set(ws.result)
}

// bigCosOptimized computes cos(x) using optimized Taylor series
func bigCosOptimized(x *BigFloat, prec uint) *BigFloat {
	if prec == 0 {
		prec = x.Prec()
	}

	// Normalize x to [-π, π]
	x = normalizeAngle(x, prec)

	workPrec := prec + 16
	ws := getTrigWorkspace(prec)

	// First term is 1
	ws.result.SetFloat64(1.0)
	ws.term.SetFloat64(1.0)

	// Pre-compute x² once
	ws.xSquared.Mul(x, x)

	// Convergence threshold
	ws.threshold.SetMantExp(NewBigFloat(1.0, workPrec), -int(prec))

	// Optimized Taylor series loop
	// cos(x) = 1 - x²/2! + x⁴/4! - ...
	// term = term * (-x²) / ((2n-1)(2n))
	for n := 1; n < 200; n++ {
		// Pre-compute denominator as single value: (2n-1)(2n)
		// This reduces one division operation per iteration
		denom := int64(2*n-1) * int64(2*n)
		ws.denom1.SetInt64(denom)

		// Update term: term = term * (-x²) / ((2n-1)(2n))
		ws.term.Mul(ws.term, ws.xSquared)
		ws.term.Neg(ws.term)
		ws.term.Quo(ws.term, ws.denom1) // Single division instead of two

		ws.result.Add(ws.result, ws.term)

		// Check convergence
		ws.temp.Abs(ws.term)
		if ws.temp.Cmp(ws.threshold) < 0 {
			break
		}
	}

	return new(BigFloat).SetPrec(prec).Set(ws.result)
}

// bigSinCosOptimized computes both sin(x) and cos(x) together
// This is more efficient when both are needed because we can share x² computation
func bigSinCosOptimized(x *BigFloat, prec uint) (sin, cos *BigFloat) {
	if prec == 0 {
		prec = x.Prec()
	}

	// Normalize x to [-π, π]
	x = normalizeAngle(x, prec)

	// Compute both using optimized versions
	// They share the x² computation internally, but we can optimize further
	// by computing x² once and passing it, but for now this is simpler
	sin = bigSinOptimized(x, prec)
	cos = bigCosOptimized(x, prec)

	return sin, cos
}
