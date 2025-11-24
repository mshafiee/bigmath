// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

import (
	"testing"
)

func TestBigFactorial(t *testing.T) {
	prec := uint(256)

	tests := []struct {
		name     string
		n        int64
		expected int64
		shouldErr bool
	}{
		{"zero", 0, 1, false},
		{"one", 1, 1, false},
		{"two", 2, 2, false},
		{"three", 3, 6, false},
		{"four", 4, 24, false},
		{"five", 5, 120, false},
		{"six", 6, 720, false},
		{"seven", 7, 5040, false},
		{"eight", 8, 40320, false},
		{"nine", 9, 362880, false},
		{"ten", 10, 3628800, false},
		{"fifteen", 15, 1307674368000, false},
		{"twenty", 20, 2432902008176640000, false},
		{"negative", -1, 0, true},
		{"negative_large", -10, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BigFactorial(tt.n, prec)
			
			if tt.shouldErr {
				// Negative factorial should return error or zero
				got, _ := result.Int64()
				if got != 0 {
					t.Errorf("BigFactorial(%d) should return 0 for negative, got %d", tt.n, got)
				}
				return
			}
			
			got, _ := result.Int64()
			if got != tt.expected {
				t.Errorf("BigFactorial(%d) = %d, want %d", tt.n, got, tt.expected)
			}
		})
	}
	
	// Property: factorial(n) = gamma(n+1)
	t.Run("factorial_equals_gamma", func(t *testing.T) {
		for n := int64(0); n <= 20; n++ {
			factorial := BigFactorial(n, prec)
			nPlusOne := NewBigFloat(float64(n+1), prec)
			gamma := BigGamma(nPlusOne, prec)
			
			// Use tolerance comparison due to floating point precision
			factVal, _ := factorial.Float64()
			gammaVal, _ := gamma.Float64()
			diff := factVal - gammaVal
			if diff < 0 {
				diff = -diff
			}
			// Relative tolerance: allow small relative error
			relTol := factVal * 1e-10
			if diff > relTol && diff > 1e-6 {
				t.Errorf("Property violated: factorial(%d) = %g != gamma(%d+1) = %g (diff %g)", n, factVal, n, gammaVal, diff)
			}
		}
	})
	
	// Property: factorial(n) = n * factorial(n-1)
	t.Run("recurrence_property", func(t *testing.T) {
		for n := int64(1); n <= 20; n++ {
			factorialN := BigFactorial(n, prec)
			factorialNMinusOne := BigFactorial(n-1, prec)
			nBig := NewBigFloat(float64(n), prec)
			expected := new(BigFloat).SetPrec(prec).Mul(nBig, factorialNMinusOne)
			
			if factorialN.Cmp(expected) != 0 {
				factNVal, _ := factorialN.Float64()
				expectedVal, _ := expected.Float64()
				t.Errorf("Property violated: factorial(%d) = %g != %d * factorial(%d) = %g", n, factNVal, n, n-1, expectedVal)
			}
		}
	})
	
	// Test precision levels
	t.Run("precision_levels", func(t *testing.T) {
		testCases := []uint{64, 128, 256, 512}
		for _, p := range testCases {
			result := BigFactorial(10, p)
			got, _ := result.Int64()
			if got != 3628800 {
				t.Errorf("BigFactorial(10) at prec %d = %d, want 3628800", p, got)
			}
		}
	})
}

func TestBigBinomial(t *testing.T) {
	prec := uint(256)

	tests := []struct {
		name     string
		n, k     int64
		expected int64
		shouldErr bool
	}{
		{"C(5,2)", 5, 2, 10, false},
		{"C(10,3)", 10, 3, 120, false},
		{"C(6,0)", 6, 0, 1, false},
		{"C(6,6)", 6, 6, 1, false},
		{"C(7,3)", 7, 3, 35, false},
		{"C(8,4)", 8, 4, 70, false},
		{"C(12,5)", 12, 5, 792, false},
		{"C(20,10)", 20, 10, 184756, false},
		{"C(15,0)", 15, 0, 1, false},
		{"C(15,15)", 15, 15, 1, false},
		{"C(15,1)", 15, 1, 15, false},
		{"C(15,14)", 15, 14, 15, false},
		{"k_equals_n_minus_1", 10, 9, 10, false},
		{"k_greater_than_n", 5, 10, 0, true},
		{"negative_n", -5, 2, 0, true},
		{"negative_k", 5, -2, 0, true},
		{"both_negative", -5, -2, 0, true},
		{"n_zero", 0, 0, 1, false},
		{"n_zero_k_one", 0, 1, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BigBinomial(tt.n, tt.k, prec)
			
			if tt.shouldErr {
				// Invalid input should return 0
				got, _ := result.Int64()
				if got != 0 {
					t.Errorf("BigBinomial(%d, %d) should return 0 for invalid input, got %d", tt.n, tt.k, got)
				}
				return
			}
			
			got, _ := result.Int64()
			if got != tt.expected {
				t.Errorf("BigBinomial(%d, %d) = %d, want %d", tt.n, tt.k, got, tt.expected)
			}
		})
	}
	
	// Property: C(n, k) = C(n, n-k)
	t.Run("symmetry_property", func(t *testing.T) {
		testCases := [][]int64{
			{10, 3},
			{15, 5},
			{20, 7},
		}
		for _, tc := range testCases {
			n, k := tc[0], tc[1]
			binomial1 := BigBinomial(n, k, prec)
			binomial2 := BigBinomial(n, n-k, prec)
			
			if binomial1.Cmp(binomial2) != 0 {
				val1, _ := binomial1.Int64()
				val2, _ := binomial2.Int64()
				t.Errorf("Property violated: C(%d, %d) = %d != C(%d, %d) = %d", n, k, val1, n, n-k, val2)
			}
		}
	})
	
	// Property: C(n, k) = C(n-1, k-1) + C(n-1, k) (Pascal's triangle)
	t.Run("pascal_triangle_property", func(t *testing.T) {
		testCases := [][]int64{
			{10, 5},
			{15, 7},
			{20, 10},
		}
		for _, tc := range testCases {
			n, k := tc[0], tc[1]
			if k > 0 && k < n {
				binomial := BigBinomial(n, k, prec)
				term1 := BigBinomial(n-1, k-1, prec)
				term2 := BigBinomial(n-1, k, prec)
				expected := new(BigFloat).SetPrec(prec).Add(term1, term2)
				
				if binomial.Cmp(expected) != 0 {
					binVal, _ := binomial.Int64()
					term1Val, _ := term1.Int64()
					term2Val, _ := term2.Int64()
					expectedVal, _ := expected.Int64()
					t.Errorf("Property violated: C(%d, %d) = %d != C(%d, %d) + C(%d, %d) = %d + %d = %d", n, k, binVal, n-1, k-1, n-1, k, term1Val, term2Val, expectedVal)
				}
			}
		}
	})
	
	// Property: sum of binomial coefficients = 2^n
	t.Run("sum_property", func(t *testing.T) {
		for n := int64(1); n <= 10; n++ {
			sum := NewBigFloat(0.0, prec)
			for k := int64(0); k <= n; k++ {
				binomial := BigBinomial(n, k, prec)
				sum.Add(sum, binomial)
			}
			
			two := NewBigFloat(2.0, prec)
			expected := BigPow(two, NewBigFloat(float64(n), prec), prec)
			
			if sum.Cmp(expected) != 0 {
				sumVal, _ := sum.Float64()
				expectedVal, _ := expected.Float64()
				t.Errorf("Property violated: sum(C(%d, k)) = %g != 2^%d = %g", n, sumVal, n, expectedVal)
			}
		}
	})
	
	// Test precision levels
	t.Run("precision_levels", func(t *testing.T) {
		testCases := []uint{64, 128, 256, 512}
		for _, p := range testCases {
			result := BigBinomial(10, 3, p)
			got, _ := result.Int64()
			if got != 120 {
				t.Errorf("BigBinomial(10, 3) at prec %d = %d, want 120", p, got)
			}
		}
	})
}
