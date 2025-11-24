// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

import (
	"math"
	"testing"
)

func TestBigFloor(t *testing.T) {
	prec := uint(256)

	tests := []struct {
		name     string
		input    float64
		expected float64
		useInf   bool
		infSign  bool
	}{
		{"positive", 3.7, 3.0, false, false},
		{"negative", -3.7, -4.0, false, false},
		{"integer", 5.0, 5.0, false, false},
		{"zero", 0.0, 0.0, false, false},
		{"small", 0.1, 0.0, false, false},
		{"near_integer_below", 3.9999999, 3.0, false, false},
		{"near_integer_above", 4.0000001, 4.0, false, false},
		{"negative_near_integer", -3.9999999, -4.0, false, false},
		{"exact_negative_integer", -5.0, -5.0, false, false},
		{"positive_infinity", 0, 0, true, false},
		{"negative_infinity", 0, 0, true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var x *BigFloat
			if tt.useInf {
				x = new(BigFloat).SetPrec(prec).SetInf(tt.infSign)
			} else {
				x = NewBigFloat(tt.input, prec)
			}
			result := BigFloor(x, prec)
			
			if tt.useInf {
				if !result.IsInf() {
					t.Errorf("BigFloor(Inf) should return Inf, got %v", result)
				}
				return
			}
			
			got, _ := result.Float64()
			if math.Abs(got-tt.expected) > 1e-10 {
				t.Errorf("BigFloor(%g) = %g, want %g", tt.input, got, tt.expected)
			}
			
			// Property: floor(x) <= x < floor(x) + 1
			xVal, _ := x.Float64()
			if got > xVal {
				t.Errorf("Property violated: floor(%g) = %g > %g", tt.input, got, xVal)
			}
			if got+1.0 <= xVal {
				t.Errorf("Property violated: floor(%g) + 1 = %g <= %g", tt.input, got+1.0, xVal)
			}
		})
	}
	
	// Test very large numbers
	t.Run("very_large_positive", func(t *testing.T) {
		x := NewBigFloat(1e100, prec)
		result := BigFloor(x, prec)
		got, _ := result.Float64()
		if got > 1e100 || got < 1e100-1 {
			t.Errorf("BigFloor(1e100) = %g, expected approximately 1e100", got)
		}
	})
	
	t.Run("very_large_negative", func(t *testing.T) {
		x := NewBigFloat(-1e100, prec)
		result := BigFloor(x, prec)
		got, _ := result.Float64()
		if got > -1e100 || got < -1e100-1 {
			t.Errorf("BigFloor(-1e100) = %g, expected approximately -1e100", got)
		}
	})
	
	// Test precision levels
	t.Run("precision_levels", func(t *testing.T) {
		testCases := []uint{64, 128, 256, 512}
		for _, p := range testCases {
			x := NewBigFloat(3.7, p)
			result := BigFloor(x, p)
			got, _ := result.Float64()
			if math.Abs(got-3.0) > 1e-10 {
				t.Errorf("BigFloor(3.7) at prec %d = %g, want 3.0", p, got)
			}
		}
	})
}

func TestBigCeil(t *testing.T) {
	prec := uint(256)

	tests := []struct {
		name     string
		input    float64
		expected float64
		useInf   bool
		infSign  bool
	}{
		{"positive", 3.2, 4.0, false, false},
		{"negative", -3.2, -3.0, false, false},
		{"integer", 5.0, 5.0, false, false},
		{"zero", 0.0, 0.0, false, false},
		{"near_integer_below", 3.9999999, 4.0, false, false},
		{"near_integer_above", 4.0000001, 5.0, false, false},
		{"negative_near_integer", -3.9999999, -3.0, false, false},
		{"exact_negative_integer", -5.0, -5.0, false, false},
		{"positive_infinity", 0, 0, true, false},
		{"negative_infinity", 0, 0, true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var x *BigFloat
			if tt.useInf {
				x = new(BigFloat).SetPrec(prec).SetInf(tt.infSign)
			} else {
				x = NewBigFloat(tt.input, prec)
			}
			result := BigCeil(x, prec)
			
			if tt.useInf {
				if !result.IsInf() {
					t.Errorf("BigCeil(Inf) should return Inf, got %v", result)
				}
				return
			}
			
			got, _ := result.Float64()
			if math.Abs(got-tt.expected) > 1e-10 {
				t.Errorf("BigCeil(%g) = %g, want %g", tt.input, got, tt.expected)
			}
			
			// Property: ceil(x) >= x > ceil(x) - 1
			xVal, _ := x.Float64()
			if got < xVal {
				t.Errorf("Property violated: ceil(%g) = %g < %g", tt.input, got, xVal)
			}
			if got-1.0 >= xVal {
				t.Errorf("Property violated: ceil(%g) - 1 = %g >= %g", tt.input, got-1.0, xVal)
			}
		})
	}
	
	// Test very large numbers
	t.Run("very_large_positive", func(t *testing.T) {
		x := NewBigFloat(1e100, prec)
		result := BigCeil(x, prec)
		got, _ := result.Float64()
		if got < 1e100 || got > 1e100+1 {
			t.Errorf("BigCeil(1e100) = %g, expected approximately 1e100", got)
		}
	})
	
	t.Run("very_large_negative", func(t *testing.T) {
		x := NewBigFloat(-1e100, prec)
		result := BigCeil(x, prec)
		got, _ := result.Float64()
		if got < -1e100-1 || got > -1e100 {
			t.Errorf("BigCeil(-1e100) = %g, expected approximately -1e100", got)
		}
	})
}

func TestBigTrunc(t *testing.T) {
	prec := uint(256)

	tests := []struct {
		name     string
		input    float64
		expected float64
		useInf   bool
		infSign  bool
	}{
		{"positive", 3.7, 3.0, false, false},
		{"negative", -3.7, -3.0, false, false},
		{"integer", 5.0, 5.0, false, false},
		{"zero", 0.0, 0.0, false, false},
		{"near_integer_below", 3.9999999, 3.0, false, false},
		{"near_integer_above", 4.0000001, 4.0, false, false},
		{"negative_near_integer", -3.9999999, -3.0, false, false},
		{"exact_negative_integer", -5.0, -5.0, false, false},
		{"positive_infinity", 0, 0, true, false},
		{"negative_infinity", 0, 0, true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var x *BigFloat
			if tt.useInf {
				x = new(BigFloat).SetPrec(prec).SetInf(tt.infSign)
			} else {
				x = NewBigFloat(tt.input, prec)
			}
			result := BigTrunc(x, prec)
			
			if tt.useInf {
				if !result.IsInf() {
					t.Errorf("BigTrunc(Inf) should return Inf, got %v", result)
				}
				return
			}
			
			got, _ := result.Float64()
			if math.Abs(got-tt.expected) > 1e-10 {
				t.Errorf("BigTrunc(%g) = %g, want %g", tt.input, got, tt.expected)
			}
			
			// Property: trunc(x) = floor(x) for x >= 0, trunc(x) = ceil(x) for x < 0
			xVal, _ := x.Float64()
			if xVal >= 0 {
				floor := BigFloor(x, prec)
				floorVal, _ := floor.Float64()
				if math.Abs(got-floorVal) > 1e-10 {
					t.Errorf("Property violated: trunc(%g) = %g != floor(%g) = %g", tt.input, got, tt.input, floorVal)
				}
			} else {
				ceil := BigCeil(x, prec)
				ceilVal, _ := ceil.Float64()
				if math.Abs(got-ceilVal) > 1e-10 {
					t.Errorf("Property violated: trunc(%g) = %g != ceil(%g) = %g", tt.input, got, tt.input, ceilVal)
				}
			}
		})
	}
	
	// Test very large numbers
	t.Run("very_large_positive", func(t *testing.T) {
		x := NewBigFloat(1e100, prec)
		result := BigTrunc(x, prec)
		got, _ := result.Float64()
		if got > 1e100 || got < 1e100-1 {
			t.Errorf("BigTrunc(1e100) = %g, expected approximately 1e100", got)
		}
	})
	
	t.Run("very_large_negative", func(t *testing.T) {
		x := NewBigFloat(-1e100, prec)
		result := BigTrunc(x, prec)
		got, _ := result.Float64()
		if got < -1e100-1 || got > -1e100 {
			t.Errorf("BigTrunc(-1e100) = %g, expected approximately -1e100", got)
		}
	})
}

func TestBigMod(t *testing.T) {
	prec := uint(256)

	tests := []struct {
		name     string
		x, y     float64
		expected float64
		checkNaN bool
	}{
		{"basic", 10.0, 3.0, 1.0, false},
		{"negative_x", -10.0, 3.0, 2.0, false},
		{"x_equals_y", 5.0, 5.0, 0.0, false},
		{"x_equals_neg_y", -5.0, 5.0, 0.0, false},
		{"x_less_than_y", 2.0, 5.0, 2.0, false},
		{"x_zero", 0.0, 5.0, 0.0, false},
		{"negative_x_positive_y", -7.0, 3.0, 2.0, false},
		{"positive_x_negative_y", 7.0, -3.0, -2.0, false},
		{"both_negative", -7.0, -3.0, -1.0, false}, // mod(-7, -3) = -1 (same sign as y)
		{"zero_y", 10.0, 0.0, 0.0, true}, // big.Float doesn't support NaN, returns 0
		{"very_large", 1e50, 1e30, 3.0735722628689113e29, false}, // 1e50 mod 1e30 is not 0
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := NewBigFloat(tt.x, prec)
			y := NewBigFloat(tt.y, prec)
			result := BigMod(x, y, prec)
			got, _ := result.Float64()
			if tt.checkNaN {
				// big.Float doesn't support NaN, so division by zero returns 0
				// This is a limitation we document
				if got != 0.0 {
					t.Errorf("BigMod(%g, %g) = %g, want 0 (big.Float limitation: no NaN support)", tt.x, tt.y, got)
				}
			} else if math.Abs(got-tt.expected) > 1e-10 {
				t.Errorf("BigMod(%g, %g) = %g, want %g", tt.x, tt.y, got, tt.expected)
			}
			
			// Property: mod(x, y) should have same sign as y (when y != 0)
			if tt.y != 0 && !tt.checkNaN {
				if (got > 0 && tt.y < 0) || (got < 0 && tt.y > 0) {
					t.Errorf("Property violated: mod(%g, %g) = %g should have same sign as y = %g", tt.x, tt.y, got, tt.y)
				}
			}
		})
	}
	
	// Test property: mod(x, y) = x - y*floor(x/y)
	t.Run("property_verification", func(t *testing.T) {
		testCases := [][]float64{
			{10.0, 3.0},
			{-10.0, 3.0},
			{7.0, 3.0},
			{-7.0, 3.0},
		}
		for _, tc := range testCases {
			x := NewBigFloat(tc[0], prec)
			y := NewBigFloat(tc[1], prec)
			mod := BigMod(x, y, prec)
			
			quo := new(BigFloat).SetPrec(prec + 32).Quo(x, y)
			floorQuo := BigFloor(quo, prec)
			expected := new(BigFloat).SetPrec(prec).Sub(x, new(BigFloat).SetPrec(prec).Mul(y, floorQuo))
			
			if mod.Cmp(expected) != 0 {
				modVal, _ := mod.Float64()
				expVal, _ := expected.Float64()
				t.Errorf("Property violated: mod(%g, %g) = %g, but x - y*floor(x/y) = %g", tc[0], tc[1], modVal, expVal)
			}
		}
	})
}

func TestBigRem(t *testing.T) {
	prec := uint(256)

	tests := []struct {
		name     string
		x, y     float64
		expected float64
	}{
		{"basic", 10.0, 3.0, 1.0},
		{"negative_x", -10.0, 3.0, -1.0},
		{"x_equals_y", 5.0, 5.0, 0.0},
		{"x_zero", 0.0, 5.0, 0.0},
		{"x_less_than_y", 2.0, 5.0, 2.0},
		{"negative_x_positive_y", -7.0, 3.0, -1.0},
		{"positive_x_negative_y", 7.0, -3.0, 1.0},
		{"both_negative", -7.0, -3.0, -1.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := NewBigFloat(tt.x, prec)
			y := NewBigFloat(tt.y, prec)
			result := BigRem(x, y, prec)
			got, _ := result.Float64()
			if math.Abs(got-tt.expected) > 1e-10 {
				t.Errorf("BigRem(%g, %g) = %g, want %g", tt.x, tt.y, got, tt.expected)
			}
			
			// Property: rem(x, y) should have same sign as x (when y != 0)
			if tt.y != 0 {
				if (got > 0 && tt.x < 0) || (got < 0 && tt.x > 0) {
					t.Errorf("Property violated: rem(%g, %g) = %g should have same sign as x = %g", tt.x, tt.y, got, tt.x)
				}
			}
		})
	}
	
	// Compare with standard library math.Remainder
	t.Run("compare_with_math_remainder", func(t *testing.T) {
		testCases := [][]float64{
			{10.0, 3.0},
			{-10.0, 3.0},
			{7.0, 3.0},
			{-7.0, 3.0},
		}
		for _, tc := range testCases {
			x := NewBigFloat(tc[0], prec)
			y := NewBigFloat(tc[1], prec)
			rem := BigRem(x, y, prec)
			remVal, _ := rem.Float64()
			
			expected := math.Remainder(tc[0], tc[1])
			if math.Abs(remVal-expected) > 1e-8 {
				t.Errorf("BigRem(%g, %g) = %g, math.Remainder = %g", tc[0], tc[1], remVal, expected)
			}
		}
	})
}
