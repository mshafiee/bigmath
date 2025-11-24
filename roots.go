// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

import (
	"math"
)

// BigCbrt computes the cube root of x using Newton-Raphson method
// Returns NaN for negative inputs (or could be extended to handle negative)
func BigCbrt(x *BigFloat, prec uint) *BigFloat {
	if prec == 0 {
		prec = x.Prec()
	}

	// Handle special cases
	if x.Sign() == 0 {
		return NewBigFloat(0.0, prec)
	}
	if x.Sign() < 0 {
		// For negative numbers, compute cube root of absolute value and negate
		absX := BigAbs(x, prec)
		result := bigCbrtPositive(absX, prec)
		result.Neg(result)
		return result
	}

	return bigCbrtPositive(x, prec)
}

// bigCbrtPositive computes cube root for positive numbers
func bigCbrtPositive(x *BigFloat, prec uint) *BigFloat {
	// Initial guess using float64 cbrt
	xFloat, _ := x.Float64()
	guess := NewBigFloat(math.Cbrt(xFloat), prec)

	// Newton-Raphson: x_{n+1} = (2*x_n + a/(x_n^2)) / 3
	// For cube root: f(x) = x^3 - a, f'(x) = 3*x^2
	// x_{n+1} = x_n - f(x_n)/f'(x_n) = x_n - (x_n^3 - a)/(3*x_n^2)
	// = (2*x_n + a/(x_n^2)) / 3
	three := NewBigFloat(3.0, prec)
	two := NewBigFloat(2.0, prec)
	temp := new(BigFloat).SetPrec(prec)
	temp2 := new(BigFloat).SetPrec(prec)
	diff := new(BigFloat).SetPrec(prec)
	threshold := new(BigFloat).SetPrec(prec).SetFloat64(1e-77) // Convergence threshold

	for i := 0; i < 100; i++ { // Max 100 iterations
		// Compute guess squared
		temp2.Mul(guess, guess)
		// Compute x divided by guess squared
		temp.Quo(x, temp2)

		// Compute 2*guess + x/(guess^2)
		temp2.Mul(two, guess)
		temp.Add(temp2, temp)

		// Update guess using Newton-Raphson formula: (2*guess + x/(guess^2)) / 3
		temp.Quo(temp, three)

		// Check convergence: |guess_new - guess|
		diff.Sub(temp, guess)
		diff = BigAbs(diff, prec)

		guess.Set(temp)

		if diff.Cmp(threshold) < 0 {
			break
		}
	}

	return guess
}

// BigRoot computes the nth root of x: x^(1/n)
// Uses Newton-Raphson method
// n must be positive
func BigRoot(n, x *BigFloat, prec uint) *BigFloat {
	if prec == 0 {
		prec = x.Prec()
	}

	// Handle special cases
	if n.Sign() <= 0 {
		return NewBigFloat(math.NaN(), prec)
	}

	if x.Sign() == 0 {
		return NewBigFloat(0.0, prec)
	}

	if x.Sign() < 0 {
		// For even roots of negative numbers, return NaN
		// For odd roots, we could handle it, but for simplicity return NaN
		// Check if n is integer and odd
		if n.IsInt() {
			nInt, _ := n.Int64()
			if nInt%2 != 0 {
				// Odd root of negative number
				absX := BigAbs(x, prec)
				result := bigRootPositive(n, absX, prec)
				result.Neg(result)
				return result
			}
		}
		return NewBigFloat(math.NaN(), prec)
	}

	return bigRootPositive(n, x, prec)
}

// bigRootPositive computes nth root for positive numbers
func bigRootPositive(n, x *BigFloat, prec uint) *BigFloat {
	// Initial guess: use x^(1/n) approximated with float64
	xFloat, _ := x.Float64()
	nFloat, _ := n.Float64()
	guess := NewBigFloat(math.Pow(xFloat, 1.0/nFloat), prec)

	// Newton-Raphson for nth root: x_{n+1} = ((n-1)*x_n + a/(x_n^(n-1))) / n
	// For f(x) = x^n - a, f'(x) = n*x^(n-1)
	// x_{n+1} = x_n - (x_n^n - a)/(n*x_n^(n-1))
	// = ((n-1)*x_n + a/(x_n^(n-1))) / n
	one := NewBigFloat(1.0, prec)
	nMinusOne := new(BigFloat).SetPrec(prec).Sub(n, one)

	temp := new(BigFloat).SetPrec(prec)
	temp2 := new(BigFloat).SetPrec(prec)
	diff := new(BigFloat).SetPrec(prec)
	threshold := new(BigFloat).SetPrec(prec).SetFloat64(1e-77)

	for i := 0; i < 100; i++ { // Max 100 iterations
		// Compute guess raised to power (n-1)
		temp2.Set(guess)
		if nMinusOne.Cmp(one) == 0 {
			// n-1 = 1, so guess^(n-1) = guess
		} else {
			// Compute guess^(n-1) using BigPow
			temp2 = BigPow(guess, nMinusOne, prec)
		}

		// Compute x divided by guess^(n-1)
		temp.Quo(x, temp2)

		// Compute (n-1)*guess + x/(guess^(n-1))
		temp2.Mul(nMinusOne, guess)
		temp.Add(temp2, temp)

		// Update guess using Newton-Raphson formula: ((n-1)*guess + x/(guess^(n-1))) / n
		temp.Quo(temp, n)

		// Check convergence: |guess_new - guess|
		diff.Sub(temp, guess)
		diff = BigAbs(diff, prec)

		guess.Set(temp)

		if diff.Cmp(threshold) < 0 {
			break
		}
	}

	return guess
}
