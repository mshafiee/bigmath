// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

import (
	"math"
	"testing"
)

func TestBigCbrt(t *testing.T) {
	prec := uint(256)

	tests := []struct {
		name      string
		input     float64
		expected  float64
		tolerance float64
	}{
		{"positive", 8.0, 2.0, 1e-10},
		{"one", 1.0, 1.0, 1e-10},
		{"zero", 0.0, 0.0, 1e-10},
		{"negative", -8.0, -2.0, 1e-10},
		{"perfect_cube_27", 27.0, 3.0, 1e-10},
		{"perfect_cube_64", 64.0, 4.0, 1e-10},
		{"perfect_cube_neg_27", -27.0, -3.0, 1e-10},
		{"non_perfect_cube", 10.0, math.Cbrt(10.0), 1e-8},
		{"very_small", 1e-30, math.Cbrt(1e-30), 1e-8},
		{"very_large", 1e30, math.Cbrt(1e30), 1e-8},
		{"fraction", 0.125, 0.5, 1e-10},
		{"negative_fraction", -0.125, -0.5, 1e-10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := NewBigFloat(tt.input, prec)
			result := BigCbrt(x, prec)
			got, _ := result.Float64()
			if math.Abs(got-tt.expected) > tt.tolerance {
				t.Errorf("BigCbrt(%g) = %g, want %g (tolerance %g)", tt.input, got, tt.expected, tt.tolerance)
			}

			// Property: cbrt(x)^3 ≈ x
			cube := new(BigFloat).SetPrec(prec).Mul(result, result)
			cube.Mul(cube, result)
			cubeVal, _ := cube.Float64()
			if math.Abs(cubeVal-tt.input) > tt.tolerance*10 {
				t.Errorf("Property violated: cbrt(%g)^3 = %g, want %g", tt.input, cubeVal, tt.input)
			}
		})
	}

	// Test precision levels
	t.Run("precision_levels", func(t *testing.T) {
		testCases := []uint{64, 128, 256, 512}
		for _, p := range testCases {
			x := NewBigFloat(8.0, p)
			result := BigCbrt(x, p)
			got, _ := result.Float64()
			if math.Abs(got-2.0) > 1e-8 {
				t.Errorf("BigCbrt(8) at prec %d = %g, want 2.0", p, got)
			}
		}
	})

	// Compare with standard library
	t.Run("compare_with_math_cbrt", func(t *testing.T) {
		testCases := []float64{1.0, 8.0, 27.0, 64.0, 10.0, 0.5, 0.001}
		for _, tc := range testCases {
			x := NewBigFloat(tc, prec)
			result := BigCbrt(x, prec)
			got, _ := result.Float64()
			expected := math.Cbrt(tc)
			if math.Abs(got-expected) > 1e-8 {
				t.Errorf("BigCbrt(%g) = %g, math.Cbrt = %g", tc, got, expected)
			}
		}
	})
}

func TestBigRoot(t *testing.T) {
	prec := uint(256)

	tests := []struct {
		name       string
		n, x       float64
		expected   float64
		tolerance  float64
		shouldFail bool
	}{
		{"square_root", 2.0, 4.0, 2.0, 1e-8, false},
		{"fourth_root", 4.0, 16.0, 2.0, 1e-8, false},
		{"one", 5.0, 1.0, 1.0, 1e-10, false},
		{"n_equals_one", 1.0, 5.0, 5.0, 1e-10, false},
		{"x_equals_zero", 3.0, 0.0, 0.0, 1e-10, false},
		{"n_equals_zero", 0.0, 4.0, 0.0, 0, true},
		{"n_negative", -2.0, 4.0, 0.0, 0, true},
		{"negative_x_even_n", 2.0, -4.0, 0.0, 0, true},
		{"negative_x_odd_n", 3.0, -8.0, -2.0, 1e-8, false},
		{"eighth_root", 8.0, 256.0, 2.0, 1e-8, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NewBigFloat(tt.n, prec)
			x := NewBigFloat(tt.x, prec)
			result := BigRoot(n, x, prec)
			got, _ := result.Float64()

			if tt.shouldFail {
				// big.Float doesn't support NaN, so NewBigFloat(math.NaN()) returns 0
				// This is a known limitation - we accept 0 for NaN cases
				if !math.IsNaN(got) && got != 0.0 {
					t.Errorf("BigRoot(%g, %g) should return NaN or 0 (big.Float limitation), got %g", tt.n, tt.x, got)
				}
				return
			}

			if math.Abs(got-tt.expected) > tt.tolerance {
				expected := math.Pow(tt.x, 1.0/tt.n)
				t.Errorf("BigRoot(%g, %g) = %g, want %g (expected from math.Pow: %g)", tt.n, tt.x, got, tt.expected, expected)
			}

			// Property: root(n, x)^n ≈ x
			if tt.n > 0 && tt.x > 0 {
				power := BigPow(result, n, prec)
				powerVal, _ := power.Float64()
				if math.Abs(powerVal-tt.x) > tt.tolerance*10 {
					t.Errorf("Property violated: root(%g, %g)^%g = %g, want %g", tt.n, tt.x, tt.n, powerVal, tt.x)
				}
			}
		})
	}

	// Test very large n
	t.Run("very_large_n", func(t *testing.T) {
		n := NewBigFloat(100.0, prec)
		x := NewBigFloat(2.0, prec)
		result := BigRoot(n, x, prec)
		got, _ := result.Float64()
		expected := math.Pow(2.0, 1.0/100.0)
		if math.Abs(got-expected) > 1e-6 {
			t.Errorf("BigRoot(100, 2) = %g, want %g", got, expected)
		}
	})

	// Test precision levels
	t.Run("precision_levels", func(t *testing.T) {
		testCases := []uint{64, 128, 256, 512}
		for _, p := range testCases {
			n := NewBigFloat(2.0, p)
			x := NewBigFloat(4.0, p)
			result := BigRoot(n, x, p)
			got, _ := result.Float64()
			if math.Abs(got-2.0) > 1e-6 {
				t.Errorf("BigRoot(2, 4) at prec %d = %g, want 2.0", p, got)
			}
		}
	})
}
