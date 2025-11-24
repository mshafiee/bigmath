// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

//go:build amd64 || arm64

package bigmath

import "math"

// Optimized root functions with reduced allocations

// bigCbrtOptimized implements optimized cube root
// Optimization: Reuse temporaries, reduce allocation in Newton-Raphson iteration
func bigCbrtOptimized(x *BigFloat, prec uint) *BigFloat {
	if prec == 0 {
		prec = x.Prec()
	}

	if x.Sign() == 0 {
		return NewBigFloat(0.0, prec)
	}

	if x.IsInf() {
		result := NewBigFloat(0.0, prec)
		if x.Sign() > 0 {
			result.SetInf(false)
		} else {
			result.SetInf(true)
		}
		return result
	}

	negative := x.Sign() < 0
	var xWork *BigFloat
	if negative {
		xWork = new(BigFloat).SetPrec(prec + 64).Neg(x)
	} else {
		xWork = new(BigFloat).SetPrec(prec + 64).Set(x)
	}

	result := bigCbrtPositiveOptimized(xWork, prec+64)
	if negative {
		result.Neg(result)
	}

	return new(BigFloat).SetPrec(prec).Set(result)
}

// bigCbrtPositiveOptimized implements optimized Newton-Raphson for positive values
// Formula: x_new = (2*x + a/x²)/3
func bigCbrtPositiveOptimized(a *BigFloat, prec uint) *BigFloat {
	// Get initial estimate
	aFloat, _ := a.Float64()
	x := NewBigFloat(1.0, prec)
	if aFloat > 0 && !a.IsInf() {
		// Better initial guess reduces iterations
		guess := math.Cbrt(aFloat)
		x.SetFloat64(guess)
	}

	// Preallocate temporaries to reduce allocations in loop
	temp1 := new(BigFloat).SetPrec(prec)
	temp2 := new(BigFloat).SetPrec(prec)
	x2 := new(BigFloat).SetPrec(prec)
	two := NewBigFloat(2.0, prec)
	three := NewBigFloat(3.0, prec)

	threshold := new(BigFloat).SetPrec(prec).SetMantExp(
		NewBigFloat(1.0, prec),
		-int(prec+10),
	)

	// Newton-Raphson with optimized allocation
	maxIter := 100
	for i := 0; i < maxIter; i++ {
		// x² = x * x
		x2.Mul(x, x)

		// temp1 = a / x²
		temp1.Quo(a, x2)

		// Compute temp2 = 2*x + temp1
		temp2.Mul(two, x)
		temp2.Add(temp2, temp1)

		// Compute x_new = temp2 / 3
		xNew := new(BigFloat).SetPrec(prec).Quo(temp2, three)

		// Check convergence
		temp1.Sub(xNew, x)
		temp1.Abs(temp1)
		if temp1.Cmp(threshold) < 0 {
			return xNew
		}
		x = xNew
	}

	return x
}

// bigRootOptimized implements optimized nth root
// Note: n is a BigFloat to match the signature, but should represent an integer
func bigRootOptimized(n, x *BigFloat, prec uint) *BigFloat {
	// Use the generic implementation for now
	// Optimization would require converting the algorithm to work with BigFloat n
	return bigRootGeneric(n, x, prec)
}
