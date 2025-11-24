// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

import (
	"math"
	"math/big"
)

// Generic implementations for basic operations (used as fallback)

// bigFloorGeneric returns the greatest integer value less than or equal to x
func bigFloorGeneric(x *BigFloat, prec uint) *BigFloat {
	if prec == 0 {
		prec = x.Prec()
	}

	// Handle special cases
	if x.IsInf() {
		return new(BigFloat).SetPrec(prec).Set(x)
	}

	// Convert to integer using Int() method (truncates toward zero)
	intVal := new(big.Int)
	x.Int(intVal)

	// Convert back to BigFloat
	result := new(BigFloat).SetPrec(prec).SetInt(intVal)

	// If x is negative and has fractional part, we need to subtract 1
	if x.Sign() < 0 {
		diff := new(BigFloat).SetPrec(prec).Sub(x, result)
		if diff.Sign() < 0 {
			one := NewBigFloat(1.0, prec)
			result.Sub(result, one)
		}
	}

	return result
}

// bigCeilGeneric returns the smallest integer value greater than or equal to x
func bigCeilGeneric(x *BigFloat, prec uint) *BigFloat {
	if prec == 0 {
		prec = x.Prec()
	}

	// Handle special cases
	if x.IsInf() {
		return new(BigFloat).SetPrec(prec).Set(x)
	}

	// Convert to integer using Int() method (truncates toward zero)
	intVal := new(big.Int)
	x.Int(intVal)

	// Convert back to BigFloat
	result := new(BigFloat).SetPrec(prec).SetInt(intVal)

	// If x is positive and has fractional part, we need to add 1
	if x.Sign() > 0 {
		diff := new(BigFloat).SetPrec(prec).Sub(x, result)
		if diff.Sign() > 0 {
			one := NewBigFloat(1.0, prec)
			result.Add(result, one)
		}
	}

	return result
}

// bigTruncGeneric returns the integer value of x truncated toward zero
func bigTruncGeneric(x *BigFloat, prec uint) *BigFloat {
	if prec == 0 {
		prec = x.Prec()
	}

	// Handle special cases
	if x.IsInf() {
		return new(BigFloat).SetPrec(prec).Set(x)
	}

	// Convert to integer using Int() method (truncates toward zero)
	intVal := new(big.Int)
	x.Int(intVal)

	// Convert back to BigFloat
	result := new(BigFloat).SetPrec(prec).SetInt(intVal)

	return result
}

// bigModGeneric returns x mod y (x - y*floor(x/y))
func bigModGeneric(x, y *BigFloat, prec uint) *BigFloat {
	if prec == 0 {
		prec = x.Prec()
	}

	// Handle special cases
	if y.Sign() == 0 {
		return NewBigFloat(0.0, prec)
	}

	if x.IsInf() || y.IsInf() {
		if x.IsInf() {
			return NewBigFloat(math.NaN(), prec)
		}
		return new(BigFloat).SetPrec(prec).Set(x)
	}

	// Compute x/y and take floor
	quo := new(BigFloat).SetPrec(prec+32).Quo(x, y)
	floorQuo := bigFloorGeneric(quo, prec)

	// Compute y * floor(x/y)
	product := new(BigFloat).SetPrec(prec).Mul(y, floorQuo)

	// Compute x - y*floor(x/y)
	result := new(BigFloat).SetPrec(prec).Sub(x, product)

	return result
}

// bigRemGeneric returns the remainder of x/y (IEEE 754 remainder)
func bigRemGeneric(x, y *BigFloat, prec uint) *BigFloat {
	if prec == 0 {
		prec = x.Prec()
	}

	// Handle special cases
	if y.Sign() == 0 {
		return NewBigFloat(0.0, prec)
	}

	if x.IsInf() || y.IsInf() {
		if x.IsInf() {
			return NewBigFloat(math.NaN(), prec)
		}
		return new(BigFloat).SetPrec(prec).Set(x)
	}

	// IEEE 754 remainder: x - y*round(x/y)
	quo := new(BigFloat).SetPrec(prec+32).Quo(x, y)

	// Round to nearest integer
	half := NewBigFloat(0.5, prec+32)
	var quoInt *BigFloat
	if quo.Sign() >= 0 {
		quoPlusHalf := new(BigFloat).SetPrec(prec+32).Add(quo, half)
		quoInt = bigFloorGeneric(quoPlusHalf, prec)
	} else {
		quoMinusHalf := new(BigFloat).SetPrec(prec+32).Sub(quo, half)
		quoInt = bigCeilGeneric(quoMinusHalf, prec)
	}

	// Compute y * round(x/y)
	product := new(BigFloat).SetPrec(prec).Mul(y, quoInt)
	// Compute x - y*round(x/y)
	result := new(BigFloat).SetPrec(prec).Sub(x, product)

	return result
}
