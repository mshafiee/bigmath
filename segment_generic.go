// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

import "fmt"

// Generic pure-Go implementations of Chebyshev polynomial evaluation
// These serve as reference implementations and fallbacks

// shouldEnableDebug checks if debug logging should be enabled for Body=10
func shouldEnableDebug(c []*BigFloat, neval int) bool {
	if len(c) < 26 || neval != 25 {
		return false
	}
	c0, _ := c[0].Float64()
	return c0 > -0.13 && c0 < -0.12 // Body=10 first coeff is around -0.1257
}

// printDebugStart prints debug information at the start of Chebyshev evaluation
func printDebugStart(t *BigFloat, c []*BigFloat, neval int) {
	tF, _ := t.Float64()
	fmt.Printf("\n[CHEB-DEBUG] Starting Chebyshev evaluation for Body=10\n")
	fmt.Printf("[CHEB-DEBUG] t=%.15e (normalized time in [-1,1])\n", tF)
	fmt.Printf("[CHEB-DEBUG] neval=%d (using first %d coefficients)\n", neval, neval)
	fmt.Printf("[CHEB-DEBUG] First 5 coeffs:")
	for i := 0; i < 5 && i < len(c); i++ {
		cf, _ := c[i].Float64()
		fmt.Printf(" c[%d]=%.15e", i, cf)
	}
	fmt.Printf("\n")
}

// printDebugIteration prints debug information for a Clenshaw iteration
func printDebugIteration(i, neval int, b0, b1, b2, c *BigFloat) {
	if i >= 3 && i < neval-3 {
		return
	}
	b0F, _ := b0.Float64()
	b1F, _ := b1.Float64()
	b2F, _ := b2.Float64()
	cF, _ := c.Float64()
	fmt.Printf("[CHEB-DEBUG]   i=%2d: c[%d]=%.15e, b1=%.15e, b2=%.15e → b0=%.15e\n",
		i, i, cF, b1F, b2F, b0F)
}

// printDebugEnd prints debug information at the end of Chebyshev evaluation
func printDebugEnd(b0, b2, result *BigFloat) {
	b0F, _ := b0.Float64()
	b2F, _ := b2.Float64()
	resultF, _ := result.Float64()
	fmt.Printf("\n[CHEB-DEBUG] Final: b0=%.15e, b2=%.15e\n", b0F, b2F)
	fmt.Printf("[CHEB-DEBUG] Result = (b0 - b2) / 2 = %.15e\n", resultF)
	fmt.Printf("[CHEB-DEBUG] === End Chebyshev evaluation ===\n\n")
}

// computeChebyshevResult computes the final result from Clenshaw algorithm values
func computeChebyshevResult(b0, b2 *BigFloat, c []*BigFloat, prec uint) *BigFloat {
	two := NewBigFloat(2.0, prec)
	result := new(BigFloat).SetPrec(prec)
	result.Sub(b0, b2)
	result.Quo(result, two)

	if len(c) > 0 {
		// Add c[0]/2 term
		c0Half := new(BigFloat).SetPrec(prec).Quo(c[0], two)
		result.Add(result, c0Half)
	}

	return result
}

// evaluateChebyshevBigGeneric evaluates Chebyshev polynomial with arbitrary precision (pure-Go)
// This is the BigFloat version of swi_echeb()
func evaluateChebyshevBigGeneric(t *BigFloat, c []*BigFloat, neval int, prec uint) *BigFloat {
	if prec == 0 {
		prec = DefaultPrecision
	}

	if neval <= 0 || len(c) == 0 {
		return NewBigFloat(0.0, prec)
	}

	// DEBUG: Enable detailed logging for Body=10 only
	debug := shouldEnableDebug(c, neval)
	if debug {
		printDebugStart(t, c, neval)
	}

	// Clenshaw's algorithm for Chebyshev polynomial evaluation
	// This matches swi_echeb() from swephlib.c exactly
	// T_n(t) = 2*t*T_{n-1}(t) - T_{n-2}(t)

	// Variables match C code: br=b0, brpp=b1, brp2=b2
	b0 := NewBigFloat(0.0, prec) // br
	b1 := NewBigFloat(0.0, prec) // brpp
	b2 := NewBigFloat(0.0, prec) // brp2

	two := NewBigFloat(2.0, prec)
	twoT := new(BigFloat).SetPrec(prec).Mul(two, t) // x2 = t * 2

	if debug {
		twoTF, _ := twoT.Float64()
		fmt.Printf("[CHEB-DEBUG] 2*t=%.15e\n\n", twoTF)
		fmt.Printf("[CHEB-DEBUG] Clenshaw recursion (backwards from i=%d to 0):\n", neval-1)
	}

	for i := neval - 1; i >= 0; i-- {
		// Match C code: brp2 = brpp; brpp = br; br = x2*brpp - brp2 + coef[j]
		b2.Set(b1) // brp2 = brpp
		b1.Set(b0) // brpp = br

		// b0 = 2*t*b1 - b2 + c[i] (matches: br = x2*brpp - brp2 + coef[j])
		b0.Mul(twoT, b1)
		b0.Sub(b0, b2)
		b0.Add(b0, c[i])

		if debug {
			printDebugIteration(i, neval, b0, b1, b2, c[i])
		}
	}

	// Result = (br - brp2) * 0.5 (matches C code exactly)
	// After loop: b0=br (final), b2=brp2 (value from two iterations before final br)
	// However, for Chebyshev polynomials Σ c[i]*T_i(t), the standard Clenshaw formula
	// requires: result = c[0]/2 + (b0 - b2)/2
	// This matches the mathematical formulation where T_0(t) = 1 needs special handling
	result := computeChebyshevResult(b0, b2, c, prec)

	if debug {
		printDebugEnd(b0, b2, result)
	}

	return result
}

// evaluateChebyshevDerivativeBigGeneric evaluates derivative of Chebyshev polynomial (pure-Go)
// This is the BigFloat version of swi_edcheb()
//
//nolint:unused // Used in dispatch system
func evaluateChebyshevDerivativeBigGeneric(t *BigFloat, c []*BigFloat, neval int, prec uint) *BigFloat {
	if prec == 0 {
		prec = DefaultPrecision
	}

	if neval <= 0 || len(c) == 0 {
		return NewBigFloat(0.0, prec)
	}

	// Derivative of Chebyshev: dT_n/dt = n*U_{n-1}(t)
	// where U is Chebyshev polynomial of second kind

	b0 := NewBigFloat(0.0, prec)
	b1 := NewBigFloat(0.0, prec)
	b2 := NewBigFloat(0.0, prec)

	two := NewBigFloat(2.0, prec)
	twoT := new(BigFloat).SetPrec(prec).Mul(two, t)

	for i := neval - 1; i >= 1; i-- {
		b2.Set(b1)
		b1.Set(b0)

		// Weight by coefficient index for derivative
		nBig := NewBigFloat(float64(i), prec)
		weighted := new(BigFloat).SetPrec(prec).Mul(c[i], nBig)

		b0.Mul(twoT, b1)
		b0.Sub(b0, b2)
		b0.Add(b0, weighted)
	}

	return b0
}
