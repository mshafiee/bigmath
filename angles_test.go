// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

import (
	"math"
	"testing"
)

// TestDegNormBig tests degree normalization to [0, 360)
func TestDegNormBig(t *testing.T) {
	tests := []struct {
		name     string
		deg      float64
		expected float64
	}{
		{"zero", 0.0, 0.0},
		{"positive", 45.0, 45.0},
		{"360", 360.0, 0.0},
		{"negative", -45.0, 315.0},
		{"large_positive", 720.0, 0.0},
		{"large_negative", -720.0, 0.0},
		{"between_cycles", 450.0, 90.0},
		{"negative_between_cycles", -90.0, 270.0},
		{"almost_360", 359.9, 359.9},
		{"just_over_360", 361.0, 1.0},
	}

	prec := uint(256)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deg := NewBigFloat(tt.deg, prec)
			result := DegNormBig(deg, prec)
			resultFloat, _ := result.Float64()

			if math.Abs(resultFloat-tt.expected) > 1e-10 {
				t.Errorf("DegNormBig(%v) = %v, want %v", tt.deg, resultFloat, tt.expected)
			}

			// Verify result is in [0, 360)
			if resultFloat < 0.0 || resultFloat >= 360.0 {
				t.Errorf("DegNormBig(%v) = %v, not in [0, 360)", tt.deg, resultFloat)
			}
		})
	}
}

// TestRadNormBig tests radian normalization to [-π, π]
func TestRadNormBig(t *testing.T) {
	tests := []struct {
		name      string
		rad       float64
		expected  float64
		tolerance float64
	}{
		{"zero", 0.0, 0.0, 1e-10},
		{"positive", math.Pi / 4, math.Pi / 4, 1e-10},
		{"negative", -math.Pi / 4, -math.Pi / 4, 1e-10},
		{"pi", math.Pi, math.Pi, 1e-10},
		{"negative_pi", -math.Pi, -math.Pi, 1e-10},
		{"just_over_pi", math.Pi + 0.1, -math.Pi + 0.1, 1e-9},
		{"just_under_neg_pi", -math.Pi - 0.1, math.Pi - 0.1, 1e-9},
		{"two_pi", 2 * math.Pi, 0.0, 1e-9},
		{"negative_two_pi", -2 * math.Pi, 0.0, 1e-9},
		{"three_pi", 3 * math.Pi, math.Pi, 1e-9}, // 3π normalizes to π (or -π, same point)
	}

	prec := uint(256)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rad := NewBigFloat(tt.rad, prec)
			result := RadNormBig(rad, prec)
			resultFloat, _ := result.Float64()

			if math.Abs(resultFloat-tt.expected) > tt.tolerance {
				t.Errorf("RadNormBig(%v) = %v, want %v", tt.rad, resultFloat, tt.expected)
			}

			// Verify result is in [-π, π]
			if resultFloat < -math.Pi-1e-9 || resultFloat > math.Pi+1e-9 {
				t.Errorf("RadNormBig(%v) = %v, not in [-π, π]", tt.rad, resultFloat)
			}
		})
	}
}

// TestRadNorm02PiBig tests radian normalization to [0, 2π)
func TestRadNorm02PiBig(t *testing.T) {
	tests := []struct {
		name      string
		rad       float64
		expected  float64
		tolerance float64
	}{
		{"zero", 0.0, 0.0, 1e-10},
		{"positive", math.Pi / 4, math.Pi / 4, 1e-10},
		{"pi", math.Pi, math.Pi, 1e-10},
		{"three_pi_over_two", 3 * math.Pi / 2, 3 * math.Pi / 2, 1e-10},
		{"two_pi", 2 * math.Pi, 0.0, 1e-9},
		{"negative", -math.Pi / 4, 2*math.Pi - math.Pi/4, 1e-9},
		{"negative_pi", -math.Pi, math.Pi, 1e-9},
		{"large_positive", 5 * math.Pi, math.Pi, 1e-9},
		{"large_negative", -5 * math.Pi, math.Pi, 1e-9},
		{"just_under_two_pi", 2*math.Pi - 0.1, 2*math.Pi - 0.1, 1e-10},
	}

	prec := uint(256)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rad := NewBigFloat(tt.rad, prec)
			result := RadNorm02PiBig(rad, prec)
			resultFloat, _ := result.Float64()

			if math.Abs(resultFloat-tt.expected) > tt.tolerance {
				t.Errorf("RadNorm02PiBig(%v) = %v, want %v", tt.rad, resultFloat, tt.expected)
			}

			// Verify result is in [0, 2π)
			if resultFloat < -1e-9 || resultFloat >= 2*math.Pi+1e-9 {
				t.Errorf("RadNorm02PiBig(%v) = %v, not in [0, 2π)", tt.rad, resultFloat)
			}
		})
	}
}

// TestAngleNormalizationConsistency tests that normalization is consistent
func TestAngleNormalizationConsistency(t *testing.T) {
	prec := uint(256)

	// Test that normalizing twice gives same result
	t.Run("idempotent_deg", func(t *testing.T) {
		deg := NewBigFloat(450.0, prec)
		result1 := DegNormBig(deg, prec)
		result2 := DegNormBig(result1, prec)

		r1, _ := result1.Float64()
		r2, _ := result2.Float64()

		if math.Abs(r1-r2) > 1e-10 {
			t.Errorf("DegNormBig not idempotent: %v != %v", r1, r2)
		}
	})

	t.Run("idempotent_rad", func(t *testing.T) {
		rad := NewBigFloat(3*math.Pi, prec)
		result1 := RadNormBig(rad, prec)
		result2 := RadNormBig(result1, prec)

		r1, _ := result1.Float64()
		r2, _ := result2.Float64()

		if math.Abs(r1-r2) > 1e-9 {
			t.Errorf("RadNormBig not idempotent: %v != %v", r1, r2)
		}
	})

	t.Run("idempotent_rad_02pi", func(t *testing.T) {
		rad := NewBigFloat(5*math.Pi, prec)
		result1 := RadNorm02PiBig(rad, prec)
		result2 := RadNorm02PiBig(result1, prec)

		r1, _ := result1.Float64()
		r2, _ := result2.Float64()

		if math.Abs(r1-r2) > 1e-9 {
			t.Errorf("RadNorm02PiBig not idempotent: %v != %v", r1, r2)
		}
	})
}

// TestAngleConversionEquivalence tests equivalence between degree and radian normalization
func TestAngleConversionEquivalence(t *testing.T) {
	prec := uint(256)

	tests := []struct {
		deg float64
	}{
		{45.0},
		{90.0},
		{180.0},
		{270.0},
		{360.0},
		{-45.0},
		{450.0},
		{-90.0},
	}

	for _, tt := range tests {
		t.Run("deg_"+string(rune(int(tt.deg))), func(t *testing.T) {
			// Convert to radians
			degBig := NewBigFloat(tt.deg, prec)
			radBig := new(BigFloat).SetPrec(prec).Mul(degBig, NewBigFloat(math.Pi/180.0, prec))

			// Normalize in degrees
			normDeg := DegNormBig(degBig, prec)
			normDegFloat, _ := normDeg.Float64()

			// Normalize in radians and convert back
			normRad := RadNorm02PiBig(radBig, prec)
			normRadFloat, _ := normRad.Float64()
			normDegFromRad := normRadFloat * 180.0 / math.Pi

			// Verify they're equivalent (within tolerance)
			if math.Abs(normDegFloat-normDegFromRad) > 1e-6 {
				t.Errorf("Normalization mismatch: deg=%v, rad=%v (from %v)",
					normDegFloat, normDegFromRad, tt.deg)
			}
		})
	}
}
