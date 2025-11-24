// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

// Generic pure-Go implementations of Chebyshev polynomial evaluation
// These serve as reference implementations and fallbacks

// computeChebyshevResult computes the final result from Clenshaw algorithm values
// This matches the C library swi_echeb() exactly: return (br - brp2) * 0.5
// Note: The Clenshaw algorithm already includes c[0] during the iteration when i=0,
// so we do NOT add c[0]/2 separately (that would double-count the constant term).
func computeChebyshevResult(b0, b2 *BigFloat, prec uint) *BigFloat {
	two := NewBigFloat(2.0, prec)
	result := new(BigFloat).SetPrec(prec)
	result.Sub(b0, b2)
	result.Quo(result, two)
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

	// Clenshaw's algorithm for Chebyshev polynomial evaluation
	// This matches swi_echeb() from swephlib.c exactly
	// T_n(t) = 2*t*T_{n-1}(t) - T_{n-2}(t)

	// Variables match C code: br=b0, brpp=b1, brp2=b2
	b0 := NewBigFloat(0.0, prec) // br
	b1 := NewBigFloat(0.0, prec) // brpp
	b2 := NewBigFloat(0.0, prec) // brp2

	two := NewBigFloat(2.0, prec)
	twoT := new(BigFloat).SetPrec(prec).Mul(two, t) // x2 = t * 2

	for i := neval - 1; i >= 0; i-- {
		// Match C code: brp2 = brpp; brpp = br; br = x2*brpp - brp2 + coef[j]
		b2.Set(b1) // brp2 = brpp
		b1.Set(b0) // brpp = br

		// b0 = 2*t*b1 - b2 + c[i] (matches: br = x2*brpp - brp2 + coef[j])
		b0.Mul(twoT, b1)
		b0.Sub(b0, b2)
		b0.Add(b0, c[i])
	}

	// Result = (br - brp2) * 0.5 (matches C code exactly)
	// The Clenshaw algorithm already incorporates c[0] during iteration i=0
	return computeChebyshevResult(b0, b2, prec)
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
