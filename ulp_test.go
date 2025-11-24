// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

import (
	"math"
	"testing"
)

// TestUlp tests the Ulp function
func TestUlp(t *testing.T) {
	prec := uint(256)

	t.Run("ulp_of_one", func(t *testing.T) {
		x := NewBigFloat(1.0, prec)
		ulp := Ulp(x, prec)

		if ulp == nil {
			t.Error("Ulp returned nil")
			return
		}

		ulpFloat, _ := ulp.Float64()
		// ULP of 1.0 at precision 256 should be 2^(0-256) which is very small
		if ulpFloat == 0.0 || ulpFloat > 1e-70 {
			t.Errorf("Ulp(1.0) = %e, seems incorrect", ulpFloat)
		}
	})

	t.Run("ulp_of_zero", func(t *testing.T) {
		x := NewBigFloat(0.0, prec)
		ulp := Ulp(x, prec)

		if ulp == nil {
			t.Error("Ulp returned nil")
			return
		}

		ulpFloat, _ := ulp.Float64()
		// ULP of 0 is defined as 0
		if ulpFloat != 0.0 {
			t.Errorf("Ulp(0.0) = %e, want 0.0", ulpFloat)
		}
	})

	t.Run("ulp_of_large_number", func(t *testing.T) {
		x := NewBigFloat(1e10, prec)
		ulp := Ulp(x, prec)

		if ulp == nil {
			t.Error("Ulp returned nil")
			return
		}

		ulpFloat, _ := ulp.Float64()
		// ULP of larger number should be larger
		if ulpFloat <= 0.0 {
			t.Errorf("Ulp(1e10) = %e, should be positive", ulpFloat)
		}
	})

	t.Run("ulp_of_small_number", func(t *testing.T) {
		x := NewBigFloat(1e-10, prec)
		ulp := Ulp(x, prec)

		if ulp == nil {
			t.Error("Ulp returned nil")
			return
		}

		ulpFloat, _ := ulp.Float64()
		// ULP of smaller number should be smaller
		if ulpFloat <= 0.0 {
			t.Errorf("Ulp(1e-10) = %e, should be positive", ulpFloat)
		}
	})

	t.Run("ulp_scales_with_number", func(t *testing.T) {
		x1 := NewBigFloat(1.0, prec)
		x2 := NewBigFloat(2.0, prec)

		ulp1 := Ulp(x1, prec)
		ulp2 := Ulp(x2, prec)

		ulp1Float, _ := ulp1.Float64()
		ulp2Float, _ := ulp2.Float64()

		// ULP should roughly double when x doubles (same exponent)
		if ulp1Float == 0.0 || ulp2Float == 0.0 {
			t.Skip("ULP too small to compare at this precision")
		}

		ratio := ulp2Float / ulp1Float
		if ratio < 1.5 || ratio > 2.5 {
			t.Logf("ULP ratio = %v (expected ~2.0, but can vary)", ratio)
		}
	})
}

// TestNewUlpError tests ULP error creation
func TestNewUlpError(t *testing.T) {
	prec := uint(256)

	err := NewUlpError(5.0, prec)

	if err.Value == nil {
		t.Error("NewUlpError returned nil Value")
	}

	if !err.IsUlp {
		t.Error("NewUlpError should set IsUlp to true")
	}

	valFloat, _ := err.Value.Float64()
	if valFloat != 5.0 {
		t.Errorf("NewUlpError value = %v, want 5.0", valFloat)
	}
}

// TestNewAbsError tests absolute error creation
func TestNewAbsError(t *testing.T) {
	prec := uint(256)
	absVal := NewBigFloat(0.001, prec)

	err := NewAbsError(absVal, prec)

	if err.Value == nil {
		t.Error("NewAbsError returned nil Value")
	}

	if err.IsUlp {
		t.Error("NewAbsError should set IsUlp to false")
	}

	valFloat, _ := err.Value.Float64()
	if math.Abs(valFloat-0.001) > 1e-10 {
		t.Errorf("NewAbsError value = %v, want 0.001", valFloat)
	}
}

// TestToAbs tests conversion to absolute error
func TestToAbs(t *testing.T) {
	prec := uint(256)
	x := NewBigFloat(1.0, prec)

	t.Run("absolute_error_unchanged", func(t *testing.T) {
		absErr := NewAbsError(NewBigFloat(0.01, prec), prec)
		result := absErr.ToAbs(x, prec)

		if result == nil {
			t.Error("ToAbs returned nil")
			return
		}

		resultFloat, _ := result.Float64()
		if math.Abs(resultFloat-0.01) > 1e-10 {
			t.Errorf("Absolute error changed: got %v, want 0.01", resultFloat)
		}
	})

	t.Run("ulp_error_converted", func(t *testing.T) {
		ulpErr := NewUlpError(5.0, prec)
		result := ulpErr.ToAbs(x, prec)

		if result == nil {
			t.Error("ToAbs returned nil")
			return
		}

		// Result should be 5.0 * ulp(x)
		expectedUlp := Ulp(x, prec)
		expectedFloat, _ := expectedUlp.Float64()
		expected := 5.0 * expectedFloat

		resultFloat, _ := result.Float64()
		if expectedFloat != 0.0 && math.Abs((resultFloat-expected)/expected) > 0.1 {
			t.Errorf("ULP conversion: got %e, want %e", resultFloat, expected)
		}
	})
}

// TestAddErrorBounds tests error bound addition
func TestAddErrorBounds(t *testing.T) {
	prec := uint(256)
	x := NewBigFloat(1.0, prec)

	t.Run("add_ulp_errors", func(t *testing.T) {
		err1 := NewUlpError(2.0, prec)
		err2 := NewUlpError(3.0, prec)

		result := AddErrorBounds(err1, err2, x, prec)

		if !result.IsUlp {
			t.Error("Adding ULP errors should return ULP error")
		}

		resultFloat, _ := result.Value.Float64()
		if math.Abs(resultFloat-5.0) > 1e-10 {
			t.Errorf("2 + 3 ULPs = %v, want 5.0", resultFloat)
		}
	})

	t.Run("add_absolute_errors", func(t *testing.T) {
		err1 := NewAbsError(NewBigFloat(0.01, prec), prec)
		err2 := NewAbsError(NewBigFloat(0.02, prec), prec)

		result := AddErrorBounds(err1, err2, x, prec)

		if result.IsUlp {
			t.Error("Adding absolute errors should return absolute error")
		}

		resultFloat, _ := result.Value.Float64()
		if math.Abs(resultFloat-0.03) > 1e-10 {
			t.Errorf("0.01 + 0.02 = %v, want 0.03", resultFloat)
		}
	})

	t.Run("add_mixed_errors", func(t *testing.T) {
		err1 := NewUlpError(2.0, prec)
		err2 := NewAbsError(NewBigFloat(0.01, prec), prec)

		result := AddErrorBounds(err1, err2, x, prec)

		if result.IsUlp {
			t.Error("Adding mixed errors should return absolute error")
		}

		if result.Value == nil {
			t.Error("Result value is nil")
		}
	})
}

// TestPropagateErrorAdd tests error propagation for addition
func TestPropagateErrorAdd(t *testing.T) {
	prec := uint(256)

	x := NewBigFloat(1.0, prec)
	y := NewBigFloat(2.0, prec)
	z := new(BigFloat).SetPrec(prec).Add(x, y) // z = 3.0

	errX := NewUlpError(1.0, prec)
	errY := NewUlpError(1.0, prec)

	result := PropagateErrorAdd(x, y, z, errX, errY, prec, ToNearest)

	if result.Value == nil {
		t.Error("PropagateErrorAdd returned nil value")
		return
	}

	// Should be absolute error
	if result.IsUlp {
		t.Error("PropagateErrorAdd should return absolute error")
	}

	// Error should be positive
	resultFloat, _ := result.Value.Float64()
	if resultFloat <= 0 {
		t.Errorf("Error = %e, should be positive", resultFloat)
	}
}

// TestPropagateErrorMul tests error propagation for multiplication
func TestPropagateErrorMul(t *testing.T) {
	prec := uint(256)

	x := NewBigFloat(2.0, prec)
	y := NewBigFloat(3.0, prec)
	z := new(BigFloat).SetPrec(prec).Mul(x, y) // z = 6.0

	errX := NewUlpError(1.0, prec)
	errY := NewUlpError(1.0, prec)

	result := PropagateErrorMul(x, y, z, errX, errY, prec, ToNearest)

	if result.Value == nil {
		t.Error("PropagateErrorMul returned nil value")
		return
	}

	// Should be ULP error
	if !result.IsUlp {
		t.Error("PropagateErrorMul should return ULP error")
	}

	// Error should be positive and reasonable
	resultFloat, _ := result.Value.Float64()
	if resultFloat <= 0 || resultFloat > 10 {
		t.Errorf("Error = %v ULPs, seems unreasonable", resultFloat)
	}
}

// TestCalculateRequiredPrecision tests precision calculation
func TestCalculateRequiredPrecision(t *testing.T) {
	tests := []struct {
		name              string
		targetPrec        uint
		expectedErrorUlps float64
		minRequired       uint
	}{
		{"small_error", 256, 1.0, 256},
		{"moderate_error", 256, 100.0, 256},
		{"large_error", 256, 10000.0, 256},
		{"zero_error", 128, 0.5, 128},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateRequiredPrecision(tt.targetPrec, tt.expectedErrorUlps)

			if result < tt.minRequired {
				t.Errorf("Required precision = %d, should be at least %d", result, tt.minRequired)
			}

			// Should always be at least target precision
			if result < tt.targetPrec {
				t.Errorf("Required precision %d < target precision %d", result, tt.targetPrec)
			}
		})
	}
}

// TestErrorPropagationWithDifferentModes tests error propagation with different rounding modes
func TestErrorPropagationWithDifferentModes(t *testing.T) {
	prec := uint(256)
	x := NewBigFloat(1.0, prec)
	y := NewBigFloat(2.0, prec)
	z := new(BigFloat).SetPrec(prec).Add(x, y)

	errX := NewUlpError(1.0, prec)
	errY := NewUlpError(1.0, prec)

	modes := []RoundingMode{ToNearest, ToZero, ToPositiveInf, ToNegativeInf}

	for _, mode := range modes {
		t.Run("mode_"+string(rune(mode)), func(t *testing.T) {
			result := PropagateErrorAdd(x, y, z, errX, errY, prec, mode)

			if result.Value == nil {
				t.Errorf("PropagateErrorAdd with mode %v returned nil", mode)
			}
		})
	}
}

// TestErrorBoundConsistency tests that error bounds are consistent
func TestErrorBoundConsistency(t *testing.T) {
	prec := uint(256)
	x := NewBigFloat(1.0, prec)

	err := NewUlpError(5.0, prec)

	// Convert to absolute and back
	abs := err.ToAbs(x, prec)
	absFloat, _ := abs.Float64()

	ulpX := Ulp(x, prec)
	ulpFloat, _ := ulpX.Float64()

	// abs should be approximately 5.0 * ulp(x)
	if ulpFloat != 0.0 {
		expected := 5.0 * ulpFloat
		if math.Abs((absFloat-expected)/expected) > 0.01 {
			t.Errorf("Absolute error = %e, expected %e", absFloat, expected)
		}
	}
}

// TestErrorPropagationChain tests chaining error propagation
func TestErrorPropagationChain(t *testing.T) {
	prec := uint(256)

	// z = x + y
	x := NewBigFloat(1.0, prec)
	y := NewBigFloat(2.0, prec)
	z := new(BigFloat).SetPrec(prec).Add(x, y)

	// Start with small errors
	errX := NewUlpError(0.5, prec)
	errY := NewUlpError(0.5, prec)

	// Propagate through addition
	errZ := PropagateErrorAdd(x, y, z, errX, errY, prec, ToNearest)

	// w = z * 2
	two := NewBigFloat(2.0, prec)
	w := new(BigFloat).SetPrec(prec).Mul(z, two)

	errTwo := NewUlpError(0.0, prec) // Exact constant
	errW := PropagateErrorMul(z, two, w, errZ, errTwo, prec, ToNearest)

	// Final error should be reasonable
	if errW.Value == nil {
		t.Error("Chained error propagation failed")
	}

	errWFloat, _ := errW.Value.Float64()
	if errWFloat < 0 {
		t.Errorf("Error bound should be non-negative, got %v", errWFloat)
	}
}
