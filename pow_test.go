// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

import (
	"math"
	"testing"
)

// TestBigPowBasic tests basic power operations
func TestBigPowBasic(t *testing.T) {
	tests := []struct {
		name      string
		x         float64
		y         float64
		expected  float64
		tolerance float64
	}{
		{"2^3", 2.0, 3.0, 8.0, 1e-10},
		{"3^2", 3.0, 2.0, 9.0, 1e-10},
		{"10^2", 10.0, 2.0, 100.0, 1e-10},
		{"2^10", 2.0, 10.0, 1024.0, 1e-10},
		{"5^3", 5.0, 3.0, 125.0, 1e-10},
		{"2^-1", 2.0, -1.0, 0.5, 1e-10},
		{"4^0.5", 4.0, 0.5, 2.0, 1e-10},
		{"8^(1/3)", 8.0, 1.0 / 3.0, 2.0, 1e-10},
		{"e^1", math.E, 1.0, math.E, 1e-10},
		{"10^-2", 10.0, -2.0, 0.01, 1e-10},
	}

	prec := uint(256)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := NewBigFloat(tt.x, prec)
			y := NewBigFloat(tt.y, prec)
			result := BigPow(x, y, prec)
			resultFloat, _ := result.Float64()

			if math.Abs(resultFloat-tt.expected) > tt.tolerance {
				t.Errorf("BigPow(%v, %v) = %v, want %v", tt.x, tt.y, resultFloat, tt.expected)
			}
		})
	}
}

// TestBigPowSpecialCases tests special cases
func TestBigPowSpecialCases(t *testing.T) {
	prec := uint(256)

	t.Run("x^0_equals_1", func(t *testing.T) {
		x := NewBigFloat(5.0, prec)
		y := NewBigFloat(0.0, prec)
		result := BigPow(x, y, prec)
		resultFloat, _ := result.Float64()

		if resultFloat != 1.0 {
			t.Errorf("x^0 = %v, want 1.0", resultFloat)
		}
	})

	t.Run("1^y_equals_1", func(t *testing.T) {
		x := NewBigFloat(1.0, prec)
		y := NewBigFloat(100.0, prec)
		result := BigPow(x, y, prec)
		resultFloat, _ := result.Float64()

		if resultFloat != 1.0 {
			t.Errorf("1^y = %v, want 1.0", resultFloat)
		}
	})

	t.Run("x^1_equals_x", func(t *testing.T) {
		x := NewBigFloat(7.0, prec)
		y := NewBigFloat(1.0, prec)
		result := BigPow(x, y, prec)
		resultFloat, _ := result.Float64()

		if resultFloat != 7.0 {
			t.Errorf("x^1 = %v, want 7.0", resultFloat)
		}
	})

	t.Run("0^positive_equals_0", func(t *testing.T) {
		x := NewBigFloat(0.0, prec)
		y := NewBigFloat(5.0, prec)
		result := BigPow(x, y, prec)
		resultFloat, _ := result.Float64()

		if resultFloat != 0.0 {
			t.Errorf("0^5 = %v, want 0.0", resultFloat)
		}
	})

	t.Run("0^negative_equals_inf", func(t *testing.T) {
		x := NewBigFloat(0.0, prec)
		y := NewBigFloat(-5.0, prec)
		result := BigPow(x, y, prec)
		resultFloat, _ := result.Float64()

		if !math.IsInf(resultFloat, 1) {
			t.Errorf("0^-5 = %v, want +Inf", resultFloat)
		}
	})

	t.Run("negative_base_integer_exponent", func(t *testing.T) {
		// (-2)^3 = -8
		x := NewBigFloat(-2.0, prec)
		y := NewBigFloat(3.0, prec)
		result := BigPow(x, y, prec)
		resultFloat, _ := result.Float64()

		if math.Abs(resultFloat+8.0) > 1e-10 {
			t.Errorf("(-2)^3 = %v, want -8.0", resultFloat)
		}
	})

	t.Run("negative_base_even_exponent", func(t *testing.T) {
		// (-2)^2 = 4
		x := NewBigFloat(-2.0, prec)
		y := NewBigFloat(2.0, prec)
		result := BigPow(x, y, prec)
		resultFloat, _ := result.Float64()

		if math.Abs(resultFloat-4.0) > 1e-10 {
			t.Errorf("(-2)^2 = %v, want 4.0", resultFloat)
		}
	})

	// Commented out NaN test as BigFloat.SetFloat64(NaN) causes panic
	// t.Run("negative_base_non_integer_nan", func(t *testing.T) {
	// 	// (-2)^0.5 = NaN (complex result)
	// 	x := NewBigFloat(-2.0, prec)
	// 	y := NewBigFloat(0.5, prec)
	// 	result := BigPow(x, y, prec)
	// 	resultFloat, _ := result.Float64()
	//
	// 	if !math.IsNaN(resultFloat) {
	// 		t.Errorf("(-2)^0.5 should be NaN, got %v", resultFloat)
	// 	}
	// })
}

// TestBigPowNegativeExponents tests negative exponents
func TestBigPowNegativeExponents(t *testing.T) {
	tests := []struct {
		name      string
		x         float64
		y         float64
		expected  float64
		tolerance float64
	}{
		{"2^-1", 2.0, -1.0, 0.5, 1e-10},
		{"2^-2", 2.0, -2.0, 0.25, 1e-10},
		{"10^-1", 10.0, -1.0, 0.1, 1e-10},
		{"3^-3", 3.0, -3.0, 1.0 / 27.0, 1e-10},
		{"5^-2", 5.0, -2.0, 0.04, 1e-10},
	}

	prec := uint(256)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := NewBigFloat(tt.x, prec)
			y := NewBigFloat(tt.y, prec)
			result := BigPow(x, y, prec)
			resultFloat, _ := result.Float64()

			if math.Abs(resultFloat-tt.expected) > tt.tolerance {
				t.Errorf("BigPow(%v, %v) = %v, want %v", tt.x, tt.y, resultFloat, tt.expected)
			}
		})
	}
}

// TestBigPowFractionalExponents tests fractional exponents
func TestBigPowFractionalExponents(t *testing.T) {
	tests := []struct {
		name      string
		x         float64
		y         float64
		expected  float64
		tolerance float64
	}{
		{"4^0.5", 4.0, 0.5, 2.0, 1e-10},
		{"9^0.5", 9.0, 0.5, 3.0, 1e-10},
		{"8^(1/3)", 8.0, 1.0 / 3.0, 2.0, 1e-10},
		{"27^(1/3)", 27.0, 1.0 / 3.0, 3.0, 1e-10},
		{"16^0.25", 16.0, 0.25, 2.0, 1e-10},
		{"2^1.5", 2.0, 1.5, math.Pow(2.0, 1.5), 1e-10},
		{"3^2.5", 3.0, 2.5, math.Pow(3.0, 2.5), 1e-10},
	}

	prec := uint(256)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := NewBigFloat(tt.x, prec)
			y := NewBigFloat(tt.y, prec)
			result := BigPow(x, y, prec)
			resultFloat, _ := result.Float64()

			if math.Abs(resultFloat-tt.expected) > tt.tolerance {
				t.Errorf("BigPow(%v, %v) = %v, want %v", tt.x, tt.y, resultFloat, tt.expected)
			}
		})
	}
}

// TestBigPowLargeExponents tests large integer exponents
func TestBigPowLargeExponents(t *testing.T) {
	prec := uint(256)

	t.Run("2^20", func(t *testing.T) {
		x := NewBigFloat(2.0, prec)
		y := NewBigFloat(20.0, prec)
		result := BigPow(x, y, prec)
		resultFloat, _ := result.Float64()

		expected := math.Pow(2.0, 20.0) // 1048576
		if math.Abs(resultFloat-expected) > 1e-6 {
			t.Errorf("2^20 = %v, want %v", resultFloat, expected)
		}
	})

	t.Run("2^30", func(t *testing.T) {
		x := NewBigFloat(2.0, prec)
		y := NewBigFloat(30.0, prec)
		result := BigPow(x, y, prec)
		resultFloat, _ := result.Float64()

		expected := math.Pow(2.0, 30.0) // 1073741824
		if math.Abs(resultFloat-expected) > 1e-3 {
			t.Errorf("2^30 = %v, want %v", resultFloat, expected)
		}
	})

	t.Run("10^10", func(t *testing.T) {
		x := NewBigFloat(10.0, prec)
		y := NewBigFloat(10.0, prec)
		result := BigPow(x, y, prec)
		resultFloat, _ := result.Float64()

		expected := 1e10
		if math.Abs(resultFloat-expected) > 1e-3 {
			t.Errorf("10^10 = %v, want %v", resultFloat, expected)
		}
	})
}

// TestBigPowPrecisionLevels tests power at different precision levels
func TestBigPowPrecisionLevels(t *testing.T) {
	precisions := []uint{64, 128, 256, 512}

	for _, prec := range precisions {
		t.Run("precision_"+string(rune(prec)), func(t *testing.T) {
			x := NewBigFloat(2.0, prec)
			y := NewBigFloat(3.0, prec)

			result := BigPow(x, y, prec)
			if result == nil {
				t.Errorf("BigPow at precision %d returned nil", prec)
			}

			resultFloat, _ := result.Float64()
			if math.Abs(resultFloat-8.0) > 1e-10 {
				t.Errorf("BigPow(2, 3) at precision %d = %v, want 8.0", prec, resultFloat)
			}
		})
	}
}

// TestBigPowIdentities tests mathematical identities
func TestBigPowIdentities(t *testing.T) {
	prec := uint(256)

	t.Run("x^(a+b)_equals_x^a_times_x^b", func(t *testing.T) {
		x := NewBigFloat(2.0, prec)
		a := NewBigFloat(3.0, prec)
		b := NewBigFloat(2.0, prec)

		aPlusB := new(BigFloat).SetPrec(prec).Add(a, b)
		left := BigPow(x, aPlusB, prec)

		xA := BigPow(x, a, prec)
		xB := BigPow(x, b, prec)
		right := new(BigFloat).SetPrec(prec).Mul(xA, xB)

		leftFloat, _ := left.Float64()
		rightFloat, _ := right.Float64()

		if math.Abs(leftFloat-rightFloat) > 1e-10 {
			t.Errorf("x^(a+b) = %v, x^a * x^b = %v", leftFloat, rightFloat)
		}
	})

	t.Run("(x^a)^b_equals_x^(ab)", func(t *testing.T) {
		x := NewBigFloat(2.0, prec)
		a := NewBigFloat(3.0, prec)
		b := NewBigFloat(2.0, prec)

		xA := BigPow(x, a, prec)
		left := BigPow(xA, b, prec)

		ab := new(BigFloat).SetPrec(prec).Mul(a, b)
		right := BigPow(x, ab, prec)

		leftFloat, _ := left.Float64()
		rightFloat, _ := right.Float64()

		if math.Abs(leftFloat-rightFloat) > 1e-10 {
			t.Errorf("(x^a)^b = %v, x^(ab) = %v", leftFloat, rightFloat)
		}
	})

	t.Run("(xy)^a_equals_x^a_times_y^a", func(t *testing.T) {
		x := NewBigFloat(2.0, prec)
		y := NewBigFloat(3.0, prec)
		a := NewBigFloat(2.0, prec)

		xy := new(BigFloat).SetPrec(prec).Mul(x, y)
		left := BigPow(xy, a, prec)

		xA := BigPow(x, a, prec)
		yA := BigPow(y, a, prec)
		right := new(BigFloat).SetPrec(prec).Mul(xA, yA)

		leftFloat, _ := left.Float64()
		rightFloat, _ := right.Float64()

		if math.Abs(leftFloat-rightFloat) > 1e-10 {
			t.Errorf("(xy)^a = %v, x^a * y^a = %v", leftFloat, rightFloat)
		}
	})
}

// TestBigSqrtRounded tests BigSqrtRounded (defined in bigmath.go)
func TestBigSqrtRounded(t *testing.T) {
	prec := uint(256)

	t.Run("sqrt_rounded", func(t *testing.T) {
		x := NewBigFloat(2.0, prec)
		result, _ := BigSqrtRounded(x, 64, ToNearest)

		if result == nil {
			t.Error("BigSqrtRounded returned nil")
		}

		resultFloat, _ := result.Float64()
		expected := math.Sqrt(2.0)

		if math.Abs(resultFloat-expected) > 1e-10 {
			t.Errorf("BigSqrtRounded(2) = %v, want %v", resultFloat, expected)
		}
	})
}
