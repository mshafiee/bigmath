// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

// BigFloor returns the greatest integer value less than or equal to x
// Uses rounding toward negative infinity
func BigFloor(x *BigFloat, prec uint) *BigFloat {
	return getDispatcher().BigFloorImpl(x, prec)
}

// BigCeil returns the smallest integer value greater than or equal to x
// Uses rounding toward positive infinity
func BigCeil(x *BigFloat, prec uint) *BigFloat {
	return getDispatcher().BigCeilImpl(x, prec)
}

// BigTrunc returns the integer value of x truncated toward zero
// Uses rounding toward zero
func BigTrunc(x *BigFloat, prec uint) *BigFloat {
	return getDispatcher().BigTruncImpl(x, prec)
}

// BigMod returns x mod y (x - y*floor(x/y))
// The result has the same sign as y
func BigMod(x, y *BigFloat, prec uint) *BigFloat {
	return getDispatcher().BigModImpl(x, y, prec)
}

// BigRem returns the remainder of x/y
// The result has the same sign as x
// This is IEEE 754 remainder operation: x - y*round(x/y)
func BigRem(x, y *BigFloat, prec uint) *BigFloat {
	return getDispatcher().BigRemImpl(x, y, prec)
}
