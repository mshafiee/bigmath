// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

//go:build !amd64 && !386

package bigmath

import "math"

// Fallback implementations for platforms without x87 FPU support.
// These functions use standard Go math operations and convert to/from BigFloat.

// ExtendedAdd: result = a + b
func extendedAdd(a, b float64) float64 {
	return a + b
}

// ExtendedSub: result = a - b
func extendedSub(a, b float64) float64 {
	return a - b
}

// ExtendedMul: result = a * b
func extendedMul(a, b float64) float64 {
	return a * b
}

// ExtendedDiv: result = a / b
func extendedDiv(a, b float64) float64 {
	return a / b
}

// ExtendedSin: result = sin(x)
func extendedSin(x float64) float64 {
	return math.Sin(x)
}

// ExtendedCos: result = cos(x)
func extendedCos(x float64) float64 {
	return math.Cos(x)
}

// ExtendedTan: result = tan(x)
func extendedTan(x float64) float64 {
	return math.Tan(x)
}

// ExtendedAtan: result = atan(x)
func extendedAtan(x float64) float64 {
	return math.Atan(x)
}

// ExtendedAtan2: result = atan2(y, x)
func extendedAtan2(y, x float64) float64 {
	return math.Atan2(y, x)
}

// ExtendedExp: result = exp(x)
func extendedExp(x float64) float64 {
	return math.Exp(x)
}

// ExtendedLog: result = log(x)
func extendedLog(x float64) float64 {
	return math.Log(x)
}

// ExtendedSqrt: result = sqrt(x)
func extendedSqrt(x float64) float64 {
	return math.Sqrt(x)
}

// ExtendedPow: result = pow(x, y) = x^y
func extendedPow(x, y float64) float64 {
	return math.Pow(x, y)
}

