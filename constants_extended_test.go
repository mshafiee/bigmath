// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

import (
	"math"
	"testing"
)

func TestBigPhi(t *testing.T) {
	// Known high-precision value: φ = (1 + √5) / 2 ≈ 1.6180339887498948482...
	knownHighPrec := 1.6180339887498948482

	tests := []struct {
		name      string
		prec      uint
		tolerance float64
	}{
		{"low_precision", 64, 1e-8},
		{"medium_precision", 128, 1e-10},
		{"high_precision", 256, 1e-12},
		{"very_high_precision", 512, 1e-14},
		{"ultra_high_precision", 1024, 1e-16},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			phi := BigPhi(tt.prec)
			phiVal, _ := phi.Float64()

			expected := (1.0 + math.Sqrt(5.0)) / 2.0
			if math.Abs(phiVal-expected) > tt.tolerance {
				t.Errorf("BigPhi at prec %d = %g, want %g (tolerance %g)", tt.prec, phiVal, expected, tt.tolerance)
			}
			
			// Compare with known high-precision value
			if math.Abs(phiVal-knownHighPrec) > tt.tolerance {
				t.Errorf("BigPhi at prec %d = %g, known high-prec = %g (tolerance %g)", tt.prec, phiVal, knownHighPrec, tt.tolerance)
			}
			
			// Property: φ² = φ + 1
			phi2 := new(BigFloat).SetPrec(tt.prec).Mul(phi, phi)
			phiPlusOne := new(BigFloat).SetPrec(tt.prec).Add(phi, NewBigFloat(1.0, tt.prec))
			phi2Val, _ := phi2.Float64()
			phiPlusOneVal, _ := phiPlusOne.Float64()
			diff := phi2Val - phiPlusOneVal
			if diff < 0 {
				diff = -diff
			}
			if diff > tt.tolerance {
				t.Errorf("Property violated: φ² = %g != φ + 1 = %g (diff %g, tolerance %g)", phi2Val, phiPlusOneVal, diff, tt.tolerance)
			}
		})
	}
}

func TestBigSqrt2(t *testing.T) {
	// Known high-precision value: √2 ≈ 1.4142135623730950488...
	knownHighPrec := 1.4142135623730950488

	tests := []struct {
		name      string
		prec      uint
		tolerance float64
	}{
		{"low_precision", 64, 1e-8},
		{"medium_precision", 128, 1e-10},
		{"high_precision", 256, 1e-12},
		{"very_high_precision", 512, 1e-14},
		{"ultra_high_precision", 1024, 1e-16},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqrt2 := BigSqrt2(tt.prec)
			sqrt2Val, _ := sqrt2.Float64()

			expected := math.Sqrt(2.0)
			if math.Abs(sqrt2Val-expected) > tt.tolerance {
				t.Errorf("BigSqrt2 at prec %d = %g, want %g (tolerance %g)", tt.prec, sqrt2Val, expected, tt.tolerance)
			}
			
			// Compare with known high-precision value
			if math.Abs(sqrt2Val-knownHighPrec) > tt.tolerance {
				t.Errorf("BigSqrt2 at prec %d = %g, known high-prec = %g (tolerance %g)", tt.prec, sqrt2Val, knownHighPrec, tt.tolerance)
			}
			
			// Property: (√2)² = 2
			sqrt2Squared := new(BigFloat).SetPrec(tt.prec).Mul(sqrt2, sqrt2)
			two := NewBigFloat(2.0, tt.prec)
			squaredVal, _ := sqrt2Squared.Float64()
			twoVal, _ := two.Float64()
			diff := squaredVal - twoVal
			if diff < 0 {
				diff = -diff
			}
			if diff > tt.tolerance {
				t.Errorf("Property violated: (√2)² = %g != 2 (diff %g, tolerance %g)", squaredVal, diff, tt.tolerance)
			}
		})
	}
}

func TestBigSqrt3(t *testing.T) {
	// Known high-precision value: √3 ≈ 1.7320508075688772935...
	knownHighPrec := 1.7320508075688772935

	tests := []struct {
		name      string
		prec      uint
		tolerance float64
	}{
		{"low_precision", 64, 1e-8},
		{"medium_precision", 128, 1e-10},
		{"high_precision", 256, 1e-12},
		{"very_high_precision", 512, 1e-14},
		{"ultra_high_precision", 1024, 1e-16},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sqrt3 := BigSqrt3(tt.prec)
			sqrt3Val, _ := sqrt3.Float64()

			expected := math.Sqrt(3.0)
			if math.Abs(sqrt3Val-expected) > tt.tolerance {
				t.Errorf("BigSqrt3 at prec %d = %g, want %g (tolerance %g)", tt.prec, sqrt3Val, expected, tt.tolerance)
			}
			
			// Compare with known high-precision value
			if math.Abs(sqrt3Val-knownHighPrec) > tt.tolerance {
				t.Errorf("BigSqrt3 at prec %d = %g, known high-prec = %g (tolerance %g)", tt.prec, sqrt3Val, knownHighPrec, tt.tolerance)
			}
			
			// Property: (√3)² = 3
			sqrt3Squared := new(BigFloat).SetPrec(tt.prec).Mul(sqrt3, sqrt3)
			three := NewBigFloat(3.0, tt.prec)
			squaredVal, _ := sqrt3Squared.Float64()
			threeVal, _ := three.Float64()
			diff := squaredVal - threeVal
			if diff < 0 {
				diff = -diff
			}
			if diff > tt.tolerance {
				t.Errorf("Property violated: (√3)² = %g != 3 (diff %g, tolerance %g)", squaredVal, diff, tt.tolerance)
			}
		})
	}
}

func TestBigLn10(t *testing.T) {
	// Known high-precision value: ln(10) ≈ 2.3025850929940456840...
	knownHighPrec := 2.3025850929940456840

	tests := []struct {
		name      string
		prec      uint
		tolerance float64
	}{
		{"low_precision", 64, 1e-8},
		{"medium_precision", 128, 1e-10},
		{"high_precision", 256, 1e-12},
		{"very_high_precision", 512, 1e-14},
		{"ultra_high_precision", 1024, 1e-16},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ln10 := BigLn10(tt.prec)
			ln10Val, _ := ln10.Float64()

			expected := math.Log(10.0)
			if math.Abs(ln10Val-expected) > tt.tolerance {
				t.Errorf("BigLn10 at prec %d = %g, want %g (tolerance %g)", tt.prec, ln10Val, expected, tt.tolerance)
			}
			
			// Compare with known high-precision value
			if math.Abs(ln10Val-knownHighPrec) > tt.tolerance {
				t.Errorf("BigLn10 at prec %d = %g, known high-prec = %g (tolerance %g)", tt.prec, ln10Val, knownHighPrec, tt.tolerance)
			}
			
			// Property: exp(ln(10)) = 10
			expLn10 := BigExp(ln10, tt.prec)
			ten := NewBigFloat(10.0, tt.prec)
			expVal, _ := expLn10.Float64()
			tenVal, _ := ten.Float64()
			diff := expVal - tenVal
			if diff < 0 {
				diff = -diff
			}
			if diff > tt.tolerance {
				t.Errorf("Property violated: exp(ln(10)) = %g != 10 (diff %g, tolerance %g)", expVal, diff, tt.tolerance)
			}
		})
	}
}
