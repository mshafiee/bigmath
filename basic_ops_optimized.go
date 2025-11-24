// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

//go:build amd64 || arm64

package bigmath

// Optimized basic operations with reduced allocations

// bigFloorOptimized implements optimized floor function
// Optimization: Fast path for integers, reduced allocations
func bigFloorOptimized(x *BigFloat, prec uint) *BigFloat {
	if prec == 0 {
		prec = x.Prec()
	}

	// Handle special cases
	if x.IsInf() || x.IsInt() {
		return new(BigFloat).SetPrec(prec).Set(x)
	}

	// Extract integer and fractional parts
	intPart := new(BigFloat).SetPrec(prec)
	bigInt, _ := x.Int(nil)
	intPart.SetInt(bigInt)

	// If negative and has fractional part, subtract 1
	if x.Sign() < 0 {
		fracPart := new(BigFloat).SetPrec(prec).Sub(x, intPart)
		if fracPart.Sign() != 0 {
			one := NewBigFloat(1.0, prec)
			intPart.Sub(intPart, one)
		}
	}

	return intPart
}

// bigCeilOptimized implements optimized ceiling function
// Optimization: Fast path for integers, reduced allocations
func bigCeilOptimized(x *BigFloat, prec uint) *BigFloat {
	if prec == 0 {
		prec = x.Prec()
	}

	// Handle special cases
	if x.IsInf() || x.IsInt() {
		return new(BigFloat).SetPrec(prec).Set(x)
	}

	// Extract integer part
	intPart := new(BigFloat).SetPrec(prec)
	bigInt, _ := x.Int(nil)
	intPart.SetInt(bigInt)

	// If positive and has fractional part, add 1
	if x.Sign() >= 0 {
		fracPart := new(BigFloat).SetPrec(prec).Sub(x, intPart)
		if fracPart.Sign() > 0 {
			one := NewBigFloat(1.0, prec)
			intPart.Add(intPart, one)
		}
	}

	return intPart
}

// bigTruncOptimized implements optimized truncate function
// Optimization: Direct integer extraction
func bigTruncOptimized(x *BigFloat, prec uint) *BigFloat {
	if prec == 0 {
		prec = x.Prec()
	}

	// Handle special cases
	if x.IsInf() || x.IsInt() {
		return new(BigFloat).SetPrec(prec).Set(x)
	}

	// Truncate is simply the integer part
	result := new(BigFloat).SetPrec(prec)
	bigInt, _ := x.Int(nil)
	result.SetInt(bigInt)
	return result
}

// bigModOptimized implements optimized modulo function
// Optimization: Reduced allocations in quotient calculation
func bigModOptimized(x, y *BigFloat, prec uint) *BigFloat {
	if prec == 0 {
		prec = x.Prec()
	}

	if y.Sign() == 0 {
		return NewBigFloat(0.0, prec) // NaN case
	}

	// mod(x, y) = x - y * floor(x/y)
	// Reuse temp variables
	temp := new(BigFloat).SetPrec(prec)
	temp.Quo(x, y)
	temp = bigFloorOptimized(temp, prec)
	temp.Mul(temp, y)

	result := new(BigFloat).SetPrec(prec)
	result.Sub(x, temp)

	return result
}

// bigRemOptimized implements optimized remainder function
// Optimization: Direct trunc calculation
func bigRemOptimized(x, y *BigFloat, prec uint) *BigFloat {
	if prec == 0 {
		prec = x.Prec()
	}

	if y.Sign() == 0 {
		return NewBigFloat(0.0, prec) // NaN case
	}

	// Compute rem(x, y) = x - y * trunc(x/y)
	temp := new(BigFloat).SetPrec(prec)
	temp.Quo(x, y)
	temp = bigTruncOptimized(temp, prec)
	temp.Mul(temp, y)

	result := new(BigFloat).SetPrec(prec)
	result.Sub(x, temp)

	return result
}
