// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

//go:build amd64 || 386

package bigmath

// Extended precision (x87 FPU) assembly function declarations
// These functions use the x87 FPU stack for 80-bit extended precision operations.

// ExtendedAdd adds two extended precision values: result = a + b
//
//go:noescape
func extendedAdd(a, b float64) float64

// ExtendedSub subtracts two extended precision values: result = a - b
//
//go:noescape
func extendedSub(a, b float64) float64

// ExtendedMul multiplies two extended precision values: result = a * b
//
//go:noescape
func extendedMul(a, b float64) float64

// ExtendedDiv divides two extended precision values: result = a / b
//
//go:noescape
func extendedDiv(a, b float64) float64

// ExtendedSin computes sin(x) using extended precision
//
//go:noescape
func extendedSin(x float64) float64

// ExtendedCos computes cos(x) using extended precision
//
//go:noescape
func extendedCos(x float64) float64

// ExtendedTan computes tan(x) using extended precision
//
//go:noescape
func extendedTan(x float64) float64

// ExtendedAtan computes atan(x) using extended precision
//
//go:noescape
func extendedAtan(x float64) float64

// ExtendedAtan2 computes atan2(y, x) using extended precision
//
//go:noescape
func extendedAtan2(y, x float64) float64

// ExtendedExp computes exp(x) using extended precision
//
//go:noescape
func extendedExp(x float64) float64

// ExtendedLog computes log(x) using extended precision
//
//go:noescape
func extendedLog(x float64) float64

// ExtendedSqrt computes sqrt(x) using extended precision
//
//go:noescape
func extendedSqrt(x float64) float64

// ExtendedPow computes pow(x, y) = x^y using extended precision
//
//go:noescape
func extendedPow(x, y float64) float64

