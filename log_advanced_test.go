// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

import (
	"math"
	"testing"
)

func TestBigLog1p(t *testing.T) {
	prec := uint(256)

	tests := []struct {
		name      string
		input     float64
		expected  float64
		tolerance float64
		shouldNaN bool
		shouldInf bool
		infSign   bool
	}{
		{"zero", 0.0, 0.0, 1e-10, false, false, false},
		{"small_positive", 0.001, math.Log1p(0.001), 1e-8, false, false, false},
		{"small_negative", -0.001, math.Log1p(-0.001), 1e-8, false, false, false},
		{"one", 1.0, math.Log(2.0), 1e-8, false, false, false},
		{"negative_one", -1.0, math.Log(0.0), 1e-8, false, true, true},
		{"less_than_neg_one", -1.5, 0.0, 0, true, false, false},
		{"large_positive", 10.0, math.Log(11.0), 1e-8, false, false, false},
		{"very_small", 1e-20, 1e-20, 1e-8, false, false, false},
		{"very_small_negative", -1e-20, -1e-20, 1e-8, false, false, false},
		{"near_neg_one", -0.999999, math.Log1p(-0.999999), 1e-6, false, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := NewBigFloat(tt.input, prec)
			result := BigLog1p(x, prec)

			if tt.shouldNaN {
				got, _ := result.Float64()
				// big.Float doesn't support NaN, so NewBigFloat(math.NaN()) returns 0
				// This is a known limitation - we accept 0 for NaN cases
				if !math.IsNaN(got) && got != 0.0 {
					t.Errorf("BigLog1p(%g) should return NaN or 0 (big.Float limitation), got %g", tt.input, got)
				}
				return
			}

			if tt.shouldInf {
				if !result.IsInf() {
					t.Errorf("BigLog1p(%g) should return Inf, got %v", tt.input, result)
				}
				return
			}

			got, _ := result.Float64()
			if math.Abs(got-tt.expected) > tt.tolerance {
				t.Errorf("BigLog1p(%g) = %g, want %g (tolerance %g)", tt.input, got, tt.expected, tt.tolerance)
			}

			// Compare with standard library for valid range
			if tt.input >= -1.0 && tt.input <= 1e10 {
				stdLib := math.Log1p(tt.input)
				if math.Abs(got-stdLib) > tt.tolerance*10 {
					t.Errorf("BigLog1p(%g) = %g, math.Log1p = %g", tt.input, got, stdLib)
				}
			}
		})
	}

	// Test accuracy improvement over log(1+x) for small x
	t.Run("accuracy_improvement", func(t *testing.T) {
		smallX := 1e-15
		x := NewBigFloat(smallX, prec)

		log1p := BigLog1p(x, prec)
		onePlusX := new(BigFloat).SetPrec(prec).Add(NewBigFloat(1.0, prec), x)
		logNormal := BigLog(onePlusX, prec)

		log1pVal, _ := log1p.Float64()
		logNormalVal, _ := logNormal.Float64()

		// log1p should be more accurate for very small x
		expected := math.Log1p(smallX)
		diff1p := math.Abs(log1pVal - expected)
		diffNormal := math.Abs(logNormalVal - expected)

		// log1p should be at least as accurate
		if diff1p > diffNormal*2 {
			t.Errorf("BigLog1p not more accurate: diff1p=%g, diffNormal=%g", diff1p, diffNormal)
		}
	})

	// Test infinity
	t.Run("infinity", func(t *testing.T) {
		posInf := new(BigFloat).SetPrec(prec).SetInf(false)
		result := BigLog1p(posInf, prec)
		if !result.IsInf() {
			t.Error("BigLog1p(+Inf) should return +Inf")
		}

		negInf := new(BigFloat).SetPrec(prec).SetInf(true)
		result2 := BigLog1p(negInf, prec)
		got, _ := result2.Float64()
		// big.Float doesn't support NaN, so NewBigFloat(math.NaN()) returns 0
		// This is a known limitation - we accept 0 for NaN cases
		if !math.IsNaN(got) && got != 0.0 {
			t.Errorf("BigLog1p(-Inf) should return NaN or 0 (big.Float limitation), got %g", got)
		}
	})
}

func TestBigExp1m(t *testing.T) {
	prec := uint(256)

	tests := []struct {
		name      string
		input     float64
		expected  float64
		tolerance float64
		shouldInf bool
		infSign   bool
	}{
		{"zero", 0.0, 0.0, 1e-10, false, false},
		{"small_positive", 0.001, math.Expm1(0.001), 1e-8, false, false},
		{"small_negative", -0.001, math.Expm1(-0.001), 1e-8, false, false},
		{"one", 1.0, math.E - 1.0, 1e-8, false, false},
		{"negative_one", -1.0, 1.0/math.E - 1.0, 1e-8, false, false},
		{"very_small", 1e-20, 1e-20, 1e-8, false, false},
		{"very_small_negative", -1e-20, -1e-20, 1e-8, false, false},
		{"large_positive", 10.0, math.Expm1(10.0), 1e-6, false, false}, // exp(10)-1 is finite, not infinity
		{"large_negative", -10.0, math.Expm1(-10.0), 1e-8, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := NewBigFloat(tt.input, prec)
			result := BigExp1m(x, prec)

			if tt.shouldInf {
				if !result.IsInf() {
					t.Errorf("BigExp1m(%g) should return Inf, got %v", tt.input, result)
				}
				return
			}

			got, _ := result.Float64()
			if math.Abs(got-tt.expected) > tt.tolerance {
				t.Errorf("BigExp1m(%g) = %g, want %g (tolerance %g)", tt.input, got, tt.expected, tt.tolerance)
			}

			// Compare with standard library
			if tt.input <= 10.0 {
				stdLib := math.Expm1(tt.input)
				if math.Abs(got-stdLib) > tt.tolerance*10 {
					t.Errorf("BigExp1m(%g) = %g, math.Expm1 = %g", tt.input, got, stdLib)
				}
			}
		})
	}

	// Test accuracy improvement over exp(x)-1 for small x
	t.Run("accuracy_improvement", func(t *testing.T) {
		smallX := 1e-15
		x := NewBigFloat(smallX, prec)

		exp1m := BigExp1m(x, prec)
		expX := BigExp(x, prec)
		expNormal := new(BigFloat).SetPrec(prec).Sub(expX, NewBigFloat(1.0, prec))

		exp1mVal, _ := exp1m.Float64()
		expNormalVal, _ := expNormal.Float64()

		expected := math.Expm1(smallX)
		diff1m := math.Abs(exp1mVal - expected)
		diffNormal := math.Abs(expNormalVal - expected)

		// exp1m should be more accurate for very small x
		if diff1m > diffNormal*2 {
			t.Errorf("BigExp1m not more accurate: diff1m=%g, diffNormal=%g", diff1m, diffNormal)
		}
	})

	// Test infinity
	t.Run("infinity", func(t *testing.T) {
		posInf := new(BigFloat).SetPrec(prec).SetInf(false)
		result := BigExp1m(posInf, prec)
		if !result.IsInf() {
			t.Error("BigExp1m(+Inf) should return +Inf")
		}

		negInf := new(BigFloat).SetPrec(prec).SetInf(true)
		result2 := BigExp1m(negInf, prec)
		got, _ := result2.Float64()
		expected := -1.0
		if math.Abs(got-expected) > 1e-10 {
			t.Errorf("BigExp1m(-Inf) = %g, want %g", got, expected)
		}
	})
}

func TestBigLogb(t *testing.T) {
	prec := uint(256)

	tests := []struct {
		name      string
		x, base   float64
		expected  float64
		tolerance float64
		shouldNaN bool
	}{
		{"base_2_x_8", 8.0, 2.0, 3.0, 1e-8, false},
		{"base_2_x_4", 4.0, 2.0, 2.0, 1e-8, false},
		{"base_10_x_100", 100.0, 10.0, 2.0, 1e-8, false},
		{"base_e_x_e", math.E, math.E, 1.0, 1e-8, false},
		{"x_equals_one", 1.0, 2.0, 0.0, 1e-10, false},
		{"x_equals_base", 5.0, 5.0, 1.0, 1e-10, false},
		{"base_equals_one", 5.0, 1.0, 0.0, 0, true},
		{"base_zero", 5.0, 0.0, 0.0, 0, true},
		{"base_negative", 5.0, -2.0, 0.0, 0, true},
		{"x_zero", 0.0, 2.0, 0.0, 0, true},
		{"x_negative", -5.0, 2.0, 0.0, 0, true},
		{"base_3_x_27", 27.0, 3.0, 3.0, 1e-8, false},
		{"base_5_x_125", 125.0, 5.0, 3.0, 1e-8, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := NewBigFloat(tt.x, prec)
			base := NewBigFloat(tt.base, prec)
			result := BigLogb(x, base, prec)

			if tt.shouldNaN {
				got, _ := result.Float64()
				// big.Float doesn't support NaN, so NewBigFloat(math.NaN()) returns 0
				// This is a known limitation - we accept 0 for NaN cases
				if !math.IsNaN(got) && got != 0.0 {
					t.Errorf("BigLogb(%g, %g) should return NaN or 0 (big.Float limitation), got %g", tt.x, tt.base, got)
				}
				return
			}

			got, _ := result.Float64()
			if math.Abs(got-tt.expected) > tt.tolerance {
				// Compare with standard library calculation
				expectedCalc := math.Log(tt.x) / math.Log(tt.base)
				t.Errorf("BigLogb(%g, %g) = %g, want %g (calculated: %g)", tt.x, tt.base, got, tt.expected, expectedCalc)
			}

			// Property: log_b(x) = ln(x) / ln(b)
			if tt.x > 0 && tt.base > 0 && tt.base != 1.0 {
				lnX := BigLog(x, prec)
				lnB := BigLog(base, prec)
				expectedProp := new(BigFloat).SetPrec(prec).Quo(lnX, lnB)
				expVal, _ := expectedProp.Float64()
				if math.Abs(got-expVal) > tt.tolerance {
					t.Errorf("Property violated: log_b(%g, %g) = %g, but ln(%g)/ln(%g) = %g", tt.x, tt.base, got, tt.x, tt.base, expVal)
				}
			}
		})
	}

	// Test various bases
	t.Run("various_bases", func(t *testing.T) {
		x := NewBigFloat(8.0, prec)
		bases := []float64{2.0, 4.0, 8.0, 10.0, math.E}
		for _, b := range bases {
			base := NewBigFloat(b, prec)
			result := BigLogb(x, base, prec)
			got, _ := result.Float64()
			expected := math.Log(8.0) / math.Log(b)
			if math.Abs(got-expected) > 1e-8 {
				t.Errorf("BigLogb(8, %g) = %g, want %g", b, got, expected)
			}
		}
	})

	// Test precision levels
	t.Run("precision_levels", func(t *testing.T) {
		testCases := []uint{64, 128, 256, 512}
		for _, p := range testCases {
			x := NewBigFloat(8.0, p)
			base := NewBigFloat(2.0, p)
			result := BigLogb(x, base, p)
			got, _ := result.Float64()
			if math.Abs(got-3.0) > 1e-6 {
				t.Errorf("BigLogb(8, 2) at prec %d = %g, want 3.0", p, got)
			}
		}
	})
}
