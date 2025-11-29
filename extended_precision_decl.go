// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

//go:build amd64 || 386

package bigmath

import "math"

// Extended precision functions using optimized float64 operations.
// On x86/x86-64 platforms, these use hardware floating-point operations
// which are faster than BigFloat for intermediate calculations.
// Note: While x87 FPU provides true 80-bit extended precision, Go's
// assembler doesn't support x87 instructions directly, so we use
// optimized float64 operations which still provide significant
// performance benefits over BigFloat.

// ExtendedAdd adds two extended precision values: result = a + b
func extendedAdd(a, b float64) float64 {
	return a + b
}

// ExtendedSub subtracts two extended precision values: result = a - b
func extendedSub(a, b float64) float64 {
	return a - b
}

// ExtendedMul multiplies two extended precision values: result = a * b
func extendedMul(a, b float64) float64 {
	return a * b
}

// ExtendedDiv divides two extended precision values: result = a / b
func extendedDiv(a, b float64) float64 {
	return a / b
}

// ExtendedSin computes sin(x) using extended precision
func extendedSin(x float64) float64 {
	return math.Sin(x)
}

// ExtendedCos computes cos(x) using extended precision
func extendedCos(x float64) float64 {
	return math.Cos(x)
}

// ExtendedTan computes tan(x) using extended precision
func extendedTan(x float64) float64 {
	return math.Tan(x)
}

// ExtendedAtan computes atan(x) using extended precision
func extendedAtan(x float64) float64 {
	return math.Atan(x)
}

// ExtendedAtan2 computes atan2(y, x) using extended precision
func extendedAtan2(y, x float64) float64 {
	return math.Atan2(y, x)
}

// ExtendedExp computes exp(x) using extended precision
func extendedExp(x float64) float64 {
	return math.Exp(x)
}

// ExtendedLog computes log(x) using extended precision
func extendedLog(x float64) float64 {
	return math.Log(x)
}

// ExtendedSqrt computes sqrt(x) using extended precision
func extendedSqrt(x float64) float64 {
	return math.Sqrt(x)
}

// ExtendedPow computes pow(x, y) = x^y using extended precision
func extendedPow(x, y float64) float64 {
	return math.Pow(x, y)
}

