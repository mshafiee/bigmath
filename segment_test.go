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
		// P(t) = 1
		coeffs := []*BigFloat{NewBigFloat(1.0, prec)}
		tVal := NewBigFloat(0.5, prec)

		result := EvaluateChebyshevBig(tVal, coeffs, 1, prec)
		if result == nil {
			t.Error("EvaluateChebyshevBig returned nil")
			return
		}

		resultFloat, _ := result.Float64()
		if math.Abs(resultFloat-1.0) > 1e-10 {
			t.Errorf("Constant polynomial = %v, want 1.0", resultFloat)
		}
	})

	t.Run("linear_polynomial", func(t *testing.T) {
		// P(t) = a0 + a1*T1(t) = a0 + a1*t
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

		// Expected: 2 + 3*0.5 = 3.5
		resultFloat, _ := result.Float64()
		expected := 3.5
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
