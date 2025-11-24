// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

// Optimized Chebyshev polynomial evaluation with reduced allocations
// This version pre-allocates workspace and uses pointer rotation instead of copying

// chebyshevWorkspace holds pre-allocated buffers for Chebyshev evaluation
type chebyshevWorkspace struct {
	b0, b1, b2 *BigFloat
	twoT       *BigFloat
	two        *BigFloat
	prec       uint
}

// getChebyshevWorkspace returns a workspace with pre-allocated buffers
func getChebyshevWorkspace(prec uint) *chebyshevWorkspace {
	if prec == 0 {
		prec = DefaultPrecision
	}
	return &chebyshevWorkspace{
		b0:   NewBigFloat(0.0, prec),
		b1:   NewBigFloat(0.0, prec),
		b2:   NewBigFloat(0.0, prec),
		twoT: NewBigFloat(0.0, prec),
		two:  NewBigFloat(2.0, prec),
		prec: prec,
	}
}

// evaluateChebyshevBigOptimized evaluates Chebyshev polynomial with optimized allocation pattern
func evaluateChebyshevBigOptimized(t *BigFloat, c []*BigFloat, neval int, prec uint) *BigFloat {
	if prec == 0 {
		prec = DefaultPrecision
	}

	if neval <= 0 || len(c) == 0 {
		return NewBigFloat(0.0, prec)
	}

	// Pre-allocate workspace (reused across iterations)
	ws := getChebyshevWorkspace(prec)

	// Initialize b0, b1, b2 to zero (reuse existing buffers)
	ws.b0.SetFloat64(0.0)
	ws.b1.SetFloat64(0.0)
	ws.b2.SetFloat64(0.0)

	// Compute 2*t once (optimized: multiply by 2 is just exponent increment, but big.Float doesn't expose this easily)
	ws.twoT.Set(t)
	ws.twoT.Mul(ws.twoT, ws.two)

	// Main loop with optimized operations (reusing pre-allocated buffers)
	for i := neval - 1; i >= 0; i-- {
		// Rotate: b2 = b1, b1 = b0
		// Using Set() on pre-allocated buffers is faster than NewBigFloat()
		ws.b2.Set(ws.b1)
		ws.b1.Set(ws.b0)

		// Fused operation: b0 = 2*t*b1 - b2 + c[i]
		// All operations reuse ws.b0 buffer, reducing allocations
		ws.b0.Mul(ws.twoT, ws.b1)
		ws.b0.Sub(ws.b0, ws.b2)
		ws.b0.Add(ws.b0, c[i])
	}

	// Result = (b0 - b2) / 2 + c[0]/2
	result := new(BigFloat).SetPrec(prec)
	if len(c) > 0 {
		c0Half := new(BigFloat).SetPrec(prec).Quo(c[0], ws.two)
		result.Sub(ws.b0, ws.b2)
		result.Quo(result, ws.two)
		result.Add(result, c0Half)
	} else {
		result.Sub(ws.b0, ws.b2)
		result.Quo(result, ws.two)
	}

	return result
}

// evaluateChebyshevDerivativeBigOptimized evaluates derivative with optimized allocation pattern
func evaluateChebyshevDerivativeBigOptimized(t *BigFloat, c []*BigFloat, neval int, prec uint) *BigFloat {
	if prec == 0 {
		prec = DefaultPrecision
	}

	if neval <= 0 || len(c) == 0 {
		return NewBigFloat(0.0, prec)
	}

	ws := getChebyshevWorkspace(prec)
	ws.b0.SetFloat64(0.0)
	ws.b1.SetFloat64(0.0)
	ws.b2.SetFloat64(0.0)

	ws.twoT.Set(t)
	ws.twoT.Mul(ws.twoT, ws.two)

	// Pre-allocate weighted coefficient buffer
	weighted := NewBigFloat(0.0, prec)

	for i := neval - 1; i >= 1; i-- {
		ws.b2.Set(ws.b1)
		ws.b1.Set(ws.b0)

		// Weight by coefficient index for derivative
		weighted.SetFloat64(float64(i))
		weighted.Mul(c[i], weighted)

		// b0 = 2*t*b1 - b2 + i*c[i]
		ws.b0.Mul(ws.twoT, ws.b1)
		ws.b0.Sub(ws.b0, ws.b2)
		ws.b0.Add(ws.b0, weighted)
	}

	return ws.b0
}

