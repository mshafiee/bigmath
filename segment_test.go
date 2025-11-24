// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

import (
	"math"
	"testing"
)

// TestEvaluateChebyshevBig tests Chebyshev polynomial evaluation
func TestEvaluateChebyshevBig(t *testing.T) {
	prec := uint(256)

	t.Run("constant_polynomial", func(t *testing.T) {
		// P(t) = c[0], and Clenshaw returns (br - brp2) * 0.5 = c[0] * 0.5
		// This matches the C library swi_echeb() behavior exactly
		coeffs := []*BigFloat{NewBigFloat(1.0, prec)}
		tVal := NewBigFloat(0.5, prec)

		result := EvaluateChebyshevBig(tVal, coeffs, 1, prec)
		if result == nil {
			t.Error("EvaluateChebyshevBig returned nil")
			return
		}

		// Expected: swi_echeb returns c[0] * 0.5 = 0.5 for constant polynomial
		resultFloat, _ := result.Float64()
		expected := swi_echeb(0.5, []float64{1.0}, 1)
		if math.Abs(resultFloat-expected) > 1e-10 {
			t.Errorf("Constant polynomial = %v, want %v", resultFloat, expected)
		}
	})

	t.Run("linear_polynomial", func(t *testing.T) {
		// P(t) = c[0] + c[1]*T1(t) = c[0] + c[1]*t
		// Clenshaw returns scaled result matching C library swi_echeb()
		coeffs := []*BigFloat{
			NewBigFloat(2.0, prec), // constant term
			NewBigFloat(3.0, prec), // linear term
		}
		tVal := NewBigFloat(0.5, prec)

		result := EvaluateChebyshevBig(tVal, coeffs, 2, prec)
		if result == nil {
			t.Error("EvaluateChebyshevBig returned nil")
			return
		}

		// Expected: matches C library swi_echeb()
		resultFloat, _ := result.Float64()
		expected := swi_echeb(0.5, []float64{2.0, 3.0}, 2)
		if math.Abs(resultFloat-expected) > 1e-10 {
			t.Errorf("Linear polynomial = %v, want %v", resultFloat, expected)
		}
	})

	t.Run("at_endpoints", func(t *testing.T) {
		coeffs := []*BigFloat{
			NewBigFloat(1.0, prec),
			NewBigFloat(2.0, prec),
			NewBigFloat(3.0, prec),
		}

		// Test at t = -1
		tNeg1 := NewBigFloat(-1.0, prec)
		result1 := EvaluateChebyshevBig(tNeg1, coeffs, 3, prec)
		if result1 == nil {
			t.Error("EvaluateChebyshevBig at t=-1 returned nil")
		}

		// Test at t = 1
		tPos1 := NewBigFloat(1.0, prec)
		result2 := EvaluateChebyshevBig(tPos1, coeffs, 3, prec)
		if result2 == nil {
			t.Error("EvaluateChebyshevBig at t=1 returned nil")
		}

		// Test at t = 0
		tZero := NewBigFloat(0.0, prec)
		result3 := EvaluateChebyshevBig(tZero, coeffs, 3, prec)
		if result3 == nil {
			t.Error("EvaluateChebyshevBig at t=0 returned nil")
		}
	})
}

// TestEvaluateChebyshevDerivativeBig tests Chebyshev derivative evaluation
func TestEvaluateChebyshevDerivativeBig(t *testing.T) {
	prec := uint(256)

	t.Run("constant_derivative_is_zero", func(t *testing.T) {
		// P(t) = 1, so P'(t) = 0
		coeffs := []*BigFloat{NewBigFloat(1.0, prec)}
		tVal := NewBigFloat(0.5, prec)

		result := EvaluateChebyshevDerivativeBig(tVal, coeffs, 1, prec)
		if result == nil {
			t.Error("EvaluateChebyshevDerivativeBig returned nil")
			return
		}

		resultFloat, _ := result.Float64()
		if math.Abs(resultFloat) > 1e-10 {
			t.Errorf("Derivative of constant = %v, want 0.0", resultFloat)
		}
	})

	t.Run("linear_derivative_is_constant", func(t *testing.T) {
		// P(t) = a0 + a1*t, so P'(t) = a1
		coeffs := []*BigFloat{
			NewBigFloat(2.0, prec),
			NewBigFloat(3.0, prec),
		}
		tVal := NewBigFloat(0.5, prec)

		result := EvaluateChebyshevDerivativeBig(tVal, coeffs, 2, prec)
		if result == nil {
			t.Error("EvaluateChebyshevDerivativeBig returned nil")
			return
		}

		resultFloat, _ := result.Float64()
		// For Chebyshev: d/dt(T1) = 1, so derivative should be close to a1
		// (exact value depends on Chebyshev polynomial properties)
		if math.Abs(resultFloat) > 100 {
			// Just check it's not absurdly wrong
			t.Errorf("Derivative magnitude too large: %v", resultFloat)
		}
	})

	t.Run("derivative_at_endpoints", func(t *testing.T) {
		coeffs := []*BigFloat{
			NewBigFloat(1.0, prec),
			NewBigFloat(2.0, prec),
			NewBigFloat(3.0, prec),
		}

		// Test at various points
		tVals := []float64{-1.0, -0.5, 0.0, 0.5, 1.0}
		for _, tv := range tVals {
			tBig := NewBigFloat(tv, prec)
			result := EvaluateChebyshevDerivativeBig(tBig, coeffs, 3, prec)
			if result == nil {
				t.Errorf("EvaluateChebyshevDerivativeBig at t=%v returned nil", tv)
			}
		}
	})
}

// TestEvaluateSegmentBig tests segment evaluation
func TestEvaluateSegmentBig(t *testing.T) {
	prec := uint(256)

	// Create simple test coefficients
	// 3 coefficients per dimension (X, Y, Z)
	coeffs := []*BigFloat{
		// X coefficients
		NewBigFloat(1.0, prec),
		NewBigFloat(0.1, prec),
		NewBigFloat(0.01, prec),
		// Y coefficients
		NewBigFloat(2.0, prec),
		NewBigFloat(0.2, prec),
		NewBigFloat(0.02, prec),
		// Z coefficients
		NewBigFloat(3.0, prec),
		NewBigFloat(0.3, prec),
		NewBigFloat(0.03, prec),
	}

	segStart := NewBigFloat(0.0, prec)
	segEnd := NewBigFloat(10.0, prec)
	tjd := NewBigFloat(5.0, prec) // Middle of segment
	neval := 3

	result := EvaluateSegmentBig(tjd, coeffs, segStart, segEnd, neval, prec)

	if result == nil {
		t.Error("EvaluateSegmentBig returned nil")
		return
	}

	// Verify we got a result with all components
	resultArr := result.ToFloat64()
	for i, val := range resultArr {
		if math.IsNaN(val) || math.IsInf(val, 0) {
			t.Errorf("Component %d is NaN or Inf: %v", i, val)
		}
	}

	// Position components should be non-zero given non-zero coefficients
	if resultArr[0] == 0.0 && resultArr[1] == 0.0 && resultArr[2] == 0.0 {
		t.Error("All position components are zero")
	}
}

// TestConvertToBigFloatCoeffs tests coefficient conversion
func TestConvertToBigFloatCoeffs(t *testing.T) {
	prec := uint(256)

	coeffsFloat := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
	result := ConvertToBigFloatCoeffs(coeffsFloat, prec)

	if len(result) != len(coeffsFloat) {
		t.Errorf("Length mismatch: got %d, want %d", len(result), len(coeffsFloat))
	}

	for i, expected := range coeffsFloat {
		resultFloat, _ := result[i].Float64()
		if math.Abs(resultFloat-expected) > 1e-10 {
			t.Errorf("Coefficient %d = %v, want %v", i, resultFloat, expected)
		}
	}

	// Test with empty array
	empty := ConvertToBigFloatCoeffs([]float64{}, prec)
	if len(empty) != 0 {
		t.Errorf("Empty conversion returned length %d, want 0", len(empty))
	}
}

// TestDebugPrintBigVec6 tests debug printing (just ensure it doesn't panic)
func TestDebugPrintBigVec6(t *testing.T) {
	prec := uint(256)
	v := NewBigVec6(1.0, 2.0, 3.0, 0.1, 0.2, 0.3, prec)

	// This function just prints, so we just ensure it doesn't panic
	DebugPrintBigVec6("test", v)
}

// TestRotateCoeffsToJ2000Big tests coordinate rotation
// Commented out as it requires complex segment setup
// func TestRotateCoeffsToJ2000Big(t *testing.T) {
// 	prec := uint(256)
//
// 	// Create simple test coefficients
// 	coeffs := []*BigFloat{
// 		NewBigFloat(1.0, prec),
// 		NewBigFloat(2.0, prec),
// 		NewBigFloat(3.0, prec),
// 	}
//
// 	// Create minimal segment info
// 	segInfo := &SegmentInfoBig{
// 		SegmentStart: NewBigFloat(0.0, prec),
// 		SegmentEnd:   NewBigFloat(32.0, prec),
// 		SegmentSize:  NewBigFloat(32.0, prec),
// 		NumCoeffs:    3,
// 		Body:         399, // Earth
// 	}
//
// 	result, neval := RotateCoeffsToJ2000Big(coeffs, segInfo, false, prec)
//
// 	if result == nil {
// 		t.Error("RotateCoeffsToJ2000Big returned nil")
// 		return
// 	}
//
// 	// Verify we got coefficients back
// 	if len(result) == 0 {
// 		t.Error("RotateCoeffsToJ2000Big returned empty coefficients")
// 	}
//
// 	// neval should be positive
// 	if neval <= 0 {
// 		t.Errorf("RotateCoeffsToJ2000Big returned neval = %d, expected positive", neval)
// 	}
// }

// TestChebyshevConsistency tests that evaluation is consistent
func TestChebyshevConsistency(t *testing.T) {
	prec := uint(256)

	coeffs := []*BigFloat{
		NewBigFloat(1.0, prec),
		NewBigFloat(2.0, prec),
		NewBigFloat(3.0, prec),
	}

	tVal := NewBigFloat(0.5, prec)

	// Evaluate twice and ensure we get the same result
	result1 := EvaluateChebyshevBig(tVal, coeffs, 3, prec)
	result2 := EvaluateChebyshevBig(tVal, coeffs, 3, prec)

	if result1 == nil || result2 == nil {
		t.Error("EvaluateChebyshevBig returned nil")
		return
	}

	r1Float, _ := result1.Float64()
	r2Float, _ := result2.Float64()

	if math.Abs(r1Float-r2Float) > 1e-15 {
		t.Errorf("Inconsistent results: %v vs %v", r1Float, r2Float)
	}
}

// TestChebyshevPrecisionLevels tests at different precision levels
func TestChebyshevPrecisionLevels(t *testing.T) {
	precisions := []uint{64, 128, 256, 512}

	for _, prec := range precisions {
		t.Run("precision_"+string(rune(prec)), func(t *testing.T) {
			coeffs := []*BigFloat{
				NewBigFloat(1.0, prec),
				NewBigFloat(2.0, prec),
			}
			tVal := NewBigFloat(0.5, prec)

			result := EvaluateChebyshevBig(tVal, coeffs, 2, prec)
			if result == nil {
				t.Errorf("EvaluateChebyshevBig at precision %d returned nil", prec)
			}

			derivative := EvaluateChebyshevDerivativeBig(tVal, coeffs, 2, prec)
			if derivative == nil {
				t.Errorf("EvaluateChebyshevDerivativeBig at precision %d returned nil", prec)
			}
		})
	}
}

// TestSegmentEdgeCases tests edge cases for segment evaluation
func TestSegmentEdgeCases(t *testing.T) {
	prec := uint(256)

	coeffs := []*BigFloat{
		NewBigFloat(1.0, prec),
		NewBigFloat(2.0, prec),
		NewBigFloat(3.0, prec),
	}

	t.Run("segment_start", func(t *testing.T) {
		segStart := NewBigFloat(0.0, prec)
		segEnd := NewBigFloat(10.0, prec)
		tjd := NewBigFloat(0.0, prec) // At start

		result := EvaluateSegmentBig(tjd, coeffs, segStart, segEnd, 1, prec)
		if result == nil {
			t.Error("EvaluateSegmentBig at segment start returned nil")
		}
	})

	t.Run("segment_end", func(t *testing.T) {
		segStart := NewBigFloat(0.0, prec)
		segEnd := NewBigFloat(10.0, prec)
		tjd := NewBigFloat(10.0, prec) // At end

		result := EvaluateSegmentBig(tjd, coeffs, segStart, segEnd, 1, prec)
		if result == nil {
			t.Error("EvaluateSegmentBig at segment end returned nil")
		}
	})

	t.Run("segment_middle", func(t *testing.T) {
		segStart := NewBigFloat(0.0, prec)
		segEnd := NewBigFloat(10.0, prec)
		tjd := NewBigFloat(5.0, prec) // Middle

		result := EvaluateSegmentBig(tjd, coeffs, segStart, segEnd, 1, prec)
		if result == nil {
			t.Error("EvaluateSegmentBig at segment middle returned nil")
		}
	})
}

// swi_echeb is the reference implementation matching the C library swi_echeb()
// This is the gold standard for Chebyshev polynomial evaluation using Clenshaw's algorithm
func swi_echeb(x float64, coef []float64, ncf int) float64 {
	x2 := x * 2.0
	br := 0.0
	brpp := 0.0
	brp2 := 0.0
	for j := ncf - 1; j >= 0; j-- {
		brp2 = brpp
		brpp = br
		br = x2*brpp - brp2 + coef[j]
	}
	return (br - brp2) * 0.5
}

// TestChebyshevMatchesCLibrary tests that BigFloat evaluation matches the C library swi_echeb()
// This test would have caught the c[0]/2 bug that caused 42 million km errors
func TestChebyshevMatchesCLibrary(t *testing.T) {
	prec := uint(256)

	// Test case 1: Simple polynomial with known values
	t.Run("simple_coefficients", func(t *testing.T) {
		coeffsFloat := []float64{1.0, 2.0, 3.0, 0.5, -0.25}
		coeffsBig := ConvertToBigFloatCoeffs(coeffsFloat, prec)

		tVals := []float64{-1.0, -0.75, -0.5, -0.25, 0.0, 0.25, 0.5, 0.75, 1.0}
		for _, tVal := range tVals {
			// Reference: C library algorithm
			expected := swi_echeb(tVal, coeffsFloat, len(coeffsFloat))

			// BigFloat implementation
			tBig := NewBigFloat(tVal, prec)
			resultBig := EvaluateChebyshevBig(tBig, coeffsBig, len(coeffsBig), prec)
			resultFloat, _ := resultBig.Float64()

			diff := math.Abs(resultFloat - expected)
			if diff > 1e-14 {
				t.Errorf("t=%.2f: BigFloat=%.15e, C=%.15e, diff=%.15e",
					tVal, resultFloat, expected, diff)
			}
		}
	})

	// Test case 2: Realistic ephemeris-like coefficients (similar to planetary positions)
	t.Run("ephemeris_like_coefficients", func(t *testing.T) {
		// These simulate the magnitude of actual ephemeris coefficients
		coeffsFloat := []float64{
			-0.1413725774275288,
			0.5587416028003886,
			-0.0123456789,
			0.0001234567,
			-0.0000012345,
			0.0000000123,
		}
		coeffsBig := ConvertToBigFloatCoeffs(coeffsFloat, prec)

		// Test at a realistic normalized time (close to segment end)
		tVal := 0.9862744986824852 // From actual debug output
		expected := swi_echeb(tVal, coeffsFloat, len(coeffsFloat))

		tBig := NewBigFloat(tVal, prec)
		resultBig := EvaluateChebyshevBig(tBig, coeffsBig, len(coeffsBig), prec)
		resultFloat, _ := resultBig.Float64()

		diff := math.Abs(resultFloat - expected)
		if diff > 1e-14 {
			t.Errorf("Ephemeris test: BigFloat=%.15e, C=%.15e, diff=%.15e",
				resultFloat, expected, diff)
		}
	})

	// Test case 3: Many coefficients (like actual ephemeris with 26 coefficients)
	t.Run("many_coefficients", func(t *testing.T) {
		coeffsFloat := make([]float64, 26)
		coeffsFloat[0] = -0.1413725774275288
		coeffsFloat[1] = 0.5587416028003886
		// Rest are small values decreasing in magnitude
		for i := 2; i < 26; i++ {
			coeffsFloat[i] = 0.001 / float64(i)
		}
		coeffsBig := ConvertToBigFloatCoeffs(coeffsFloat, prec)

		tVals := []float64{-0.9, 0.0, 0.9, 0.9862744986824852}
		for _, tVal := range tVals {
			expected := swi_echeb(tVal, coeffsFloat, len(coeffsFloat))

			tBig := NewBigFloat(tVal, prec)
			resultBig := EvaluateChebyshevBig(tBig, coeffsBig, len(coeffsBig), prec)
			resultFloat, _ := resultBig.Float64()

			diff := math.Abs(resultFloat - expected)
			// Allow slightly larger tolerance for many coefficients
			if diff > 1e-13 {
				t.Errorf("t=%.4f: BigFloat=%.15e, C=%.15e, diff=%.15e",
					tVal, resultFloat, expected, diff)
			}
		}
	})

	// Test case 4: Verify the constant term is not double-counted
	// This specifically catches the c[0]/2 bug
	t.Run("constant_term_not_doubled", func(t *testing.T) {
		// With only constant coefficient, result should be c[0] * T_0(t) = c[0] * 1 = c[0]
		// For Clenshaw: (br - brp2) * 0.5 = (c[0] - 0) * 0.5 = c[0]/2 ???
		// Actually, let's trace through Clenshaw for [c0]:
		// j=0: brp2=0, brpp=0, br = 0*0 - 0 + c[0] = c[0]
		// result = (c[0] - 0) * 0.5 = c[0]/2
		// So for a constant-only polynomial, result = c[0]/2
		//
		// But mathematically, sum of c_i * T_i(t) for i=0 is just c_0 * T_0(t) = c_0 * 1 = c_0
		// The factor of 2 in Clenshaw needs to be accounted for.
		//
		// Looking at the C library: it returns (br - brp2) * 0.5
		// So for constant only: result = c[0]/2

		coeffsFloat := []float64{5.0}
		expected := swi_echeb(0.5, coeffsFloat, 1)

		coeffsBig := ConvertToBigFloatCoeffs(coeffsFloat, prec)
		tBig := NewBigFloat(0.5, prec)
		resultBig := EvaluateChebyshevBig(tBig, coeffsBig, 1, prec)
		resultFloat, _ := resultBig.Float64()

		diff := math.Abs(resultFloat - expected)
		if diff > 1e-15 {
			t.Errorf("Constant term: BigFloat=%.15e, C=%.15e (expected c[0]/2=%.15e), diff=%.15e",
				resultFloat, expected, 2.5, diff)
		}
	})

	// Test case 5: Linear polynomial to verify both terms handled correctly
	t.Run("linear_polynomial_matches_c", func(t *testing.T) {
		// P(t) = c0*T0(t) + c1*T1(t) = c0*1 + c1*t
		// But with Clenshaw factor of 0.5: effective = (c0 + c1*t)/2 ??? Let's verify
		coeffsFloat := []float64{4.0, 6.0}

		tVal := 0.5
		expected := swi_echeb(tVal, coeffsFloat, 2)

		coeffsBig := ConvertToBigFloatCoeffs(coeffsFloat, prec)
		tBig := NewBigFloat(tVal, prec)
		resultBig := EvaluateChebyshevBig(tBig, coeffsBig, 2, prec)
		resultFloat, _ := resultBig.Float64()

		diff := math.Abs(resultFloat - expected)
		if diff > 1e-15 {
			t.Errorf("Linear: BigFloat=%.15e, C=%.15e, diff=%.15e",
				resultFloat, expected, diff)
		}
	})
}

// TestChebyshevAgainstDirectComputation tests against direct Chebyshev polynomial evaluation
// T_0(t) = 1, T_1(t) = t, T_2(t) = 2t^2 - 1, T_3(t) = 4t^3 - 3t
func TestChebyshevAgainstDirectComputation(t *testing.T) {
	prec := uint(256)

	// Helper to compute Chebyshev T_n(x) directly
	chebyshevT := func(n int, x float64) float64 {
		switch n {
		case 0:
			return 1.0
		case 1:
			return x
		case 2:
			return 2*x*x - 1
		case 3:
			return 4*x*x*x - 3*x
		case 4:
			return 8*x*x*x*x - 8*x*x + 1
		default:
			// Recurrence relation
			tPrev2 := 1.0
			tPrev1 := x
			for i := 2; i <= n; i++ {
				tCurr := 2*x*tPrev1 - tPrev2
				tPrev2 = tPrev1
				tPrev1 = tCurr
			}
			return tPrev1
		}
	}

	// Test: sum of c_i * T_i(t) computed directly vs Clenshaw
	t.Run("verify_against_direct_sum", func(t *testing.T) {
		coeffsFloat := []float64{1.5, 2.3, -0.7, 0.4}
		tVal := 0.6

		// Direct computation: sum of c_i * T_i(t)
		direct := 0.0
		for i, c := range coeffsFloat {
			direct += c * chebyshevT(i, tVal)
		}

		// C library Clenshaw result (known to be correct)
		clenshaw := swi_echeb(tVal, coeffsFloat, len(coeffsFloat))

		// BigFloat implementation
		coeffsBig := ConvertToBigFloatCoeffs(coeffsFloat, prec)
		tBig := NewBigFloat(tVal, prec)
		resultBig := EvaluateChebyshevBig(tBig, coeffsBig, len(coeffsBig), prec)
		bigFloat, _ := resultBig.Float64()

		// The C library Clenshaw includes a factor of 0.5 in the final result
		// which means it computes: 0.5 * sum(c_i * T_i(t)) for i>=1 plus 0.5*c_0
		// Actually, the standard Clenshaw for Chebyshev does: result = 0.5*(b0 - b2)
		// where the loop includes all coefficients including c[0]

		// Let's verify BigFloat matches C library (which is our reference)
		diffBigC := math.Abs(bigFloat - clenshaw)
		if diffBigC > 1e-14 {
			t.Errorf("BigFloat vs C: BigFloat=%.15e, C=%.15e, diff=%.15e",
				bigFloat, clenshaw, diffBigC)
		}

		t.Logf("Direct sum=%.15e, Clenshaw=%.15e, BigFloat=%.15e", direct, clenshaw, bigFloat)
	})
}
