// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

// BigPhi returns the golden ratio φ = (1 + √5) / 2 ≈ 1.6180339887498948482...
// with specified precision
func BigPhi(prec uint) *BigFloat {
	if prec == 0 {
		prec = DefaultPrecision
	}

	// φ = (1 + √5) / 2
	sqrt5 := BigSqrt(NewBigFloat(5.0, prec), prec)
	one := NewBigFloat(1.0, prec)
	numerator := new(BigFloat).SetPrec(prec).Add(one, sqrt5)
	two := NewBigFloat(2.0, prec)
	result := new(BigFloat).SetPrec(prec).Quo(numerator, two)

	return result
}

// BigSqrt2 returns √2 ≈ 1.4142135623730950488... with specified precision
func BigSqrt2(prec uint) *BigFloat {
	if prec == 0 {
		prec = DefaultPrecision
	}

	return BigSqrt(NewBigFloat(2.0, prec), prec)
}

// BigSqrt3 returns √3 ≈ 1.7320508075688772935... with specified precision
func BigSqrt3(prec uint) *BigFloat {
	if prec == 0 {
		prec = DefaultPrecision
	}

	return BigSqrt(NewBigFloat(3.0, prec), prec)
}

// BigLn10 returns ln(10) ≈ 2.3025850929940456840... with specified precision
func BigLn10(prec uint) *BigFloat {
	if prec == 0 {
		prec = DefaultPrecision
	}

	return BigLog(NewBigFloat(10.0, prec), prec)
}
