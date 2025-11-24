// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

import (
	"math"
	"testing"
)

// TestBigLog2 tests the BigLog2 constant
func TestBigLog2(t *testing.T) {
	prec := uint(256)
	result := BigLog2(prec)
	resultFloat, _ := result.Float64()

	expected := math.Ln2
	if math.Abs(resultFloat-expected) > 1e-10 {
		t.Errorf("BigLog2() = %v, want %v", resultFloat, expected)
	}

	// Test with default precision
	result2 := BigLog2(0)
	if result2 == nil {
		t.Error("BigLog2(0) returned nil")
	}
}

// TestBigJ2000 tests the BigJ2000 constant
func TestBigJ2000(t *testing.T) {
	prec := uint(256)
	result := BigJ2000(prec)
	resultFloat, _ := result.Float64()

	expected := 2451545.0
	if math.Abs(resultFloat-expected) > 1e-10 {
		t.Errorf("BigJ2000() = %v, want %v", resultFloat, expected)
	}

	// Test with default precision
	result2 := BigJ2000(0)
	if result2 == nil {
		t.Error("BigJ2000(0) returned nil")
	}
}

// TestBigLightSpeedAUperDay tests the light speed constant
func TestBigLightSpeedAUperDay(t *testing.T) {
	prec := uint(256)
	result := BigLightSpeedAUperDay(prec)
	resultFloat, _ := result.Float64()

	expected := 173.1446327205363
	if math.Abs(resultFloat-expected) > 1e-10 {
		t.Errorf("BigLightSpeedAUperDay() = %v, want %v", resultFloat, expected)
	}

	// Test with default precision
	result2 := BigLightSpeedAUperDay(0)
	if result2 == nil {
		t.Error("BigLightSpeedAUperDay(0) returned nil")
	}
}

// TestBigJulianCentury tests the Julian century constant
func TestBigJulianCentury(t *testing.T) {
	prec := uint(256)
	result := BigJulianCentury(prec)
	resultFloat, _ := result.Float64()

	expected := 36525.0
	if math.Abs(resultFloat-expected) > 1e-10 {
		t.Errorf("BigJulianCentury() = %v, want %v", resultFloat, expected)
	}

	// Test with default precision
	result2 := BigJulianCentury(0)
	if result2 == nil {
		t.Error("BigJulianCentury(0) returned nil")
	}
}

// TestBigJulianMillennium tests the Julian millennium constant
func TestBigJulianMillennium(t *testing.T) {
	prec := uint(256)
	result := BigJulianMillennium(prec)
	resultFloat, _ := result.Float64()

	expected := 365250.0
	if math.Abs(resultFloat-expected) > 1e-10 {
		t.Errorf("BigJulianMillennium() = %v, want %v", resultFloat, expected)
	}

	// Test with default precision
	result2 := BigJulianMillennium(0)
	if result2 == nil {
		t.Error("BigJulianMillennium(0) returned nil")
	}
}

// TestBigE tests Euler's number constant
func TestBigE(t *testing.T) {
	prec := uint(256)
	result := BigE(prec)
	resultFloat, _ := result.Float64()

	expected := math.E
	if math.Abs(resultFloat-expected) > 1e-10 {
		t.Errorf("BigE() = %v, want %v", resultFloat, expected)
	}

	// Test with default precision
	result2 := BigE(0)
	if result2 == nil {
		t.Error("BigE(0) returned nil")
	}
}

// TestBigEulerGamma tests Euler's constant
func TestBigEulerGamma(t *testing.T) {
	prec := uint(256)
	result := BigEulerGamma(prec)
	resultFloat, _ := result.Float64()

	expected := 0.5772156649015329
	if math.Abs(resultFloat-expected) > 1e-10 {
		t.Errorf("BigEulerGamma() = %v, want ~%v", resultFloat, expected)
	}

	// Test with default precision
	result2 := BigEulerGamma(0)
	if result2 == nil {
		t.Error("BigEulerGamma(0) returned nil")
	}
}

// TestBigCatalan tests Catalan's constant
func TestBigCatalan(t *testing.T) {
	prec := uint(256)
	result := BigCatalan(prec)
	resultFloat, _ := result.Float64()

	expected := 0.915965594177219
	if math.Abs(resultFloat-expected) > 1e-10 {
		t.Errorf("BigCatalan() = %v, want ~%v", resultFloat, expected)
	}

	// Test with default precision
	result2 := BigCatalan(0)
	if result2 == nil {
		t.Error("BigCatalan(0) returned nil")
	}
}

// TestBigLog10 tests the BigLog10 function
func TestBigLog10(t *testing.T) {
	tests := []struct {
		name      string
		x         float64
		expected  float64
		tolerance float64
	}{
		{"one", 1.0, 0.0, 1e-10},
		{"ten", 10.0, 1.0, 1e-10},
		{"hundred", 100.0, 2.0, 1e-10},
		{"thousand", 1000.0, 3.0, 1e-10},
		{"two", 2.0, math.Log10(2.0), 1e-10},
		{"half", 0.5, math.Log10(0.5), 1e-10},
		{"e", math.E, math.Log10(math.E), 1e-10},
	}

	prec := uint(256)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := NewBigFloat(tt.x, prec)
			result := BigLog10(x, prec)
			resultFloat, _ := result.Float64()

			if math.Abs(resultFloat-tt.expected) > tt.tolerance {
				t.Errorf("BigLog10(%v) = %v, want %v", tt.x, resultFloat, tt.expected)
			}
		})
	}
}

// TestConstantsPrecision tests that constants maintain precision
func TestConstantsPrecision(t *testing.T) {
	precisions := []uint{64, 128, 256, 512, 1024}

	for _, prec := range precisions {
		t.Run("precision_"+string(rune(prec)), func(t *testing.T) {
			// Test that all constants work at different precision levels
			constants := []struct {
				name string
				fn   func(uint) *BigFloat
			}{
				{"BigLog2", BigLog2},
				{"BigJ2000", BigJ2000},
				{"BigLightSpeedAUperDay", BigLightSpeedAUperDay},
				{"BigJulianCentury", BigJulianCentury},
				{"BigJulianMillennium", BigJulianMillennium},
				{"BigE", BigE},
				{"BigEulerGamma", BigEulerGamma},
				{"BigCatalan", BigCatalan},
			}

			for _, c := range constants {
				result := c.fn(prec)
				if result == nil {
					t.Errorf("%s(%d) returned nil", c.name, prec)
				}
				if result.Prec() != prec {
					t.Errorf("%s(%d) has precision %d, want %d", c.name, prec, result.Prec(), prec)
				}
			}
		})
	}
}

// TestConstantsRelationships tests mathematical relationships between constants
func TestConstantsRelationships(t *testing.T) {
	prec := uint(256)

	t.Run("e^ln(2)_approx_2", func(t *testing.T) {
		ln2 := BigLog2(prec)
		e_ln2 := BigExp(ln2, prec)
		e_ln2_float, _ := e_ln2.Float64()

		if math.Abs(e_ln2_float-2.0) > 1e-10 {
			t.Errorf("e^ln(2) = %v, want 2.0", e_ln2_float)
		}
	})

	t.Run("ln(e)_equals_1", func(t *testing.T) {
		e := BigE(prec)
		lnE := BigLog(e, prec)
		lnEFloat, _ := lnE.Float64()

		if math.Abs(lnEFloat-1.0) > 1e-10 {
			t.Errorf("ln(e) = %v, want 1.0", lnEFloat)
		}
	})

	t.Run("julian_millennium_10x_century", func(t *testing.T) {
		century := BigJulianCentury(prec)
		millennium := BigJulianMillennium(prec)

		centuryFloat, _ := century.Float64()
		millenniumFloat, _ := millennium.Float64()

		expected := centuryFloat * 10.0
		if math.Abs(millenniumFloat-expected) > 1e-10 {
			t.Errorf("Millennium = %v, want %v (10x century)", millenniumFloat, expected)
		}
	})
}
