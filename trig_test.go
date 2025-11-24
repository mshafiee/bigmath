// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

import (
	"math"
	"testing"
)

// TestBigCos tests the BigCos function
func TestBigCos(t *testing.T) {
	tests := []struct {
		name      string
		x         float64
		prec      uint
		expected  float64
		tolerance float64
	}{
		{"zero", 0.0, 256, 1.0, 1e-10},
		{"pi/2", math.Pi / 2, 256, 0.0, 1e-10},
		{"pi", math.Pi, 256, -1.0, 1e-10},
		{"3pi/2", 3 * math.Pi / 2, 256, 0.0, 1e-10},
		{"2pi", 2 * math.Pi, 256, 1.0, 1e-10},
		{"pi/4", math.Pi / 4, 256, math.Sqrt(2) / 2, 1e-10},
		{"pi/6", math.Pi / 6, 256, math.Sqrt(3) / 2, 1e-10},
		{"pi/3", math.Pi / 3, 256, 0.5, 1e-10},
		{"negative", -math.Pi / 4, 256, math.Sqrt(2) / 2, 1e-10},
		{"small", 0.001, 256, math.Cos(0.001), 1e-10},
		{"large", 10.0, 256, math.Cos(10.0), 1e-10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := NewBigFloat(tt.x, tt.prec)
			result := BigCos(x, tt.prec)
			resultFloat, _ := result.Float64()

			if math.Abs(resultFloat-tt.expected) > tt.tolerance {
				t.Errorf("BigCos(%v) = %v, want %v", tt.x, resultFloat, tt.expected)
			}
		})
	}
}

// TestBigTan tests the BigTan function
func TestBigTan(t *testing.T) {
	tests := []struct {
		name      string
		x         float64
		prec      uint
		expected  float64
		tolerance float64
	}{
		{"zero", 0.0, 256, 0.0, 1e-10},
		{"pi/4", math.Pi / 4, 256, 1.0, 1e-10},
		{"pi/6", math.Pi / 6, 256, 1.0 / math.Sqrt(3), 1e-10},
		{"pi/3", math.Pi / 3, 256, math.Sqrt(3), 1e-10},
		{"negative", -math.Pi / 4, 256, -1.0, 1e-10},
		{"small", 0.001, 256, math.Tan(0.001), 1e-10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := NewBigFloat(tt.x, tt.prec)
			result := BigTan(x, tt.prec)
			resultFloat, _ := result.Float64()

			if math.Abs(resultFloat-tt.expected) > tt.tolerance {
				t.Errorf("BigTan(%v) = %v, want %v", tt.x, resultFloat, tt.expected)
			}
		})
	}
}

// TestBigAtan tests the BigAtan function
func TestBigAtan(t *testing.T) {
	tests := []struct {
		name      string
		x         float64
		prec      uint
		expected  float64
		tolerance float64
	}{
		{"zero", 0.0, 256, 0.0, 1e-10},
		{"one", 1.0, 256, math.Pi / 4, 1e-10},
		{"negative_one", -1.0, 256, -math.Pi / 4, 1e-10},
		{"small", 0.1, 256, math.Atan(0.1), 1e-10},
		{"large", 10.0, 256, math.Atan(10.0), 1e-10},
		{"sqrt3", math.Sqrt(3), 256, math.Pi / 3, 1e-10},
		{"1/sqrt3", 1.0 / math.Sqrt(3), 256, math.Pi / 6, 1e-10},
		{"negative", -0.5, 256, math.Atan(-0.5), 1e-10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := NewBigFloat(tt.x, tt.prec)
			result := BigAtan(x, tt.prec)
			resultFloat, _ := result.Float64()

			if math.Abs(resultFloat-tt.expected) > tt.tolerance {
				t.Errorf("BigAtan(%v) = %v, want %v", tt.x, resultFloat, tt.expected)
			}
		})
	}
}

// TestBigAtan2 tests the BigAtan2 function
func TestBigAtan2(t *testing.T) {
	tests := []struct {
		name      string
		y         float64
		x         float64
		prec      uint
		expected  float64
		tolerance float64
	}{
		{"quadrant_1", 1.0, 1.0, 256, math.Pi / 4, 1e-10},
		{"quadrant_2", 1.0, -1.0, 256, 3 * math.Pi / 4, 1e-10},
		{"quadrant_3", -1.0, -1.0, 256, -3 * math.Pi / 4, 1e-10},
		{"quadrant_4", -1.0, 1.0, 256, -math.Pi / 4, 1e-10},
		{"positive_x_axis", 0.0, 1.0, 256, 0.0, 1e-10},
		{"negative_x_axis", 0.0, -1.0, 256, math.Pi, 1e-10},
		{"positive_y_axis", 1.0, 0.0, 256, math.Pi / 2, 1e-10},
		{"negative_y_axis", -1.0, 0.0, 256, -math.Pi / 2, 1e-10},
		{"both_zero", 0.0, 0.0, 256, 0.0, 1e-10},
		{"general", 3.0, 4.0, 256, math.Atan2(3.0, 4.0), 1e-10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			y := NewBigFloat(tt.y, tt.prec)
			x := NewBigFloat(tt.x, tt.prec)
			result := BigAtan2(y, x, tt.prec)
			resultFloat, _ := result.Float64()

			if math.Abs(resultFloat-tt.expected) > tt.tolerance {
				t.Errorf("BigAtan2(%v, %v) = %v, want %v", tt.y, tt.x, resultFloat, tt.expected)
			}
		})
	}
}

// TestBigAsin tests the BigAsin function
func TestBigAsin(t *testing.T) {
	tests := []struct {
		name      string
		x         float64
		prec      uint
		expected  float64
		tolerance float64
	}{
		{"zero", 0.0, 256, 0.0, 1e-10},
		{"one", 1.0, 256, math.Pi / 2, 1e-10},
		{"negative_one", -1.0, 256, -math.Pi / 2, 1e-10},
		{"half", 0.5, 256, math.Pi / 6, 1e-10},
		{"sqrt2/2", math.Sqrt(2) / 2, 256, math.Pi / 4, 1e-10},
		{"sqrt3/2", math.Sqrt(3) / 2, 256, math.Pi / 3, 1e-10},
		{"negative_half", -0.5, 256, -math.Pi / 6, 1e-10},
		{"small", 0.1, 256, math.Asin(0.1), 1e-10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := NewBigFloat(tt.x, tt.prec)
			result := BigAsin(x, tt.prec)
			resultFloat, _ := result.Float64()

			if math.Abs(resultFloat-tt.expected) > tt.tolerance {
				t.Errorf("BigAsin(%v) = %v, want %v", tt.x, resultFloat, tt.expected)
			}
		})
	}
}

// TestTrigRounded tests all rounded trigonometric functions
func TestTrigRounded(t *testing.T) {
	x := NewBigFloat(math.Pi/4, 256)

	t.Run("BigSinRounded", func(t *testing.T) {
		result, ternary := BigSinRounded(x, 64, ToNearest)
		if result == nil {
			t.Error("BigSinRounded returned nil")
		}
		_ = ternary // ternary is the rounding direction
	})

	t.Run("BigCosRounded", func(t *testing.T) {
		result, ternary := BigCosRounded(x, 64, ToNearest)
		if result == nil {
			t.Error("BigCosRounded returned nil")
		}
		_ = ternary
	})

	t.Run("BigTanRounded", func(t *testing.T) {
		result, ternary := BigTanRounded(x, 64, ToNearest)
		if result == nil {
			t.Error("BigTanRounded returned nil")
		}
		_ = ternary
	})

	t.Run("BigAtanRounded", func(t *testing.T) {
		result, ternary := BigAtanRounded(x, 64, ToNearest)
		if result == nil {
			t.Error("BigAtanRounded returned nil")
		}
		_ = ternary
	})

	t.Run("BigAtan2Rounded", func(t *testing.T) {
		y := NewBigFloat(1.0, 256)
		result, ternary := BigAtan2Rounded(y, x, 64, ToNearest)
		if result == nil {
			t.Error("BigAtan2Rounded returned nil")
		}
		_ = ternary
	})

	t.Run("BigAsinRounded", func(t *testing.T) {
		x := NewBigFloat(0.5, 256)
		result, ternary := BigAsinRounded(x, 64, ToNearest)
		if result == nil {
			t.Error("BigAsinRounded returned nil")
		}
		_ = ternary
	})

	t.Run("BigAcosRounded", func(t *testing.T) {
		x := NewBigFloat(0.5, 256)
		result, ternary := BigAcosRounded(x, 64, ToNearest)
		if result == nil {
			t.Error("BigAcosRounded returned nil")
		}
		_ = ternary
	})
}

// TestTrigIdentities tests mathematical identities
func TestTrigIdentities(t *testing.T) {
	prec := uint(256)
	x := NewBigFloat(0.7, prec)

	t.Run("sin^2 + cos^2 = 1", func(t *testing.T) {
		sinX := BigSin(x, prec)
		cosX := BigCos(x, prec)

		sin2 := new(BigFloat).SetPrec(prec).Mul(sinX, sinX)
		cos2 := new(BigFloat).SetPrec(prec).Mul(cosX, cosX)
		sum := new(BigFloat).SetPrec(prec).Add(sin2, cos2)

		one := NewBigFloat(1.0, prec)
		diff := new(BigFloat).SetPrec(prec).Sub(sum, one)
		diffFloat, _ := diff.Float64()

		if math.Abs(diffFloat) > 1e-10 {
			t.Errorf("sin^2 + cos^2 = %v, want 1.0", sum)
		}
	})

	t.Run("tan = sin/cos", func(t *testing.T) {
		sinX := BigSin(x, prec)
		cosX := BigCos(x, prec)
		tanX := BigTan(x, prec)

		computed := new(BigFloat).SetPrec(prec).Quo(sinX, cosX)
		diff := new(BigFloat).SetPrec(prec).Sub(tanX, computed)
		diffFloat, _ := diff.Float64()

		if math.Abs(diffFloat) > 1e-10 {
			t.Errorf("tan != sin/cos, difference: %v", diffFloat)
		}
	})

	t.Run("atan(tan(x)) = x", func(t *testing.T) {
		tanX := BigTan(x, prec)
		atanTanX := BigAtan(tanX, prec)

		diff := new(BigFloat).SetPrec(prec).Sub(atanTanX, x)
		diffFloat, _ := diff.Float64()

		if math.Abs(diffFloat) > 1e-10 {
			t.Errorf("atan(tan(x)) = %v, want %v", atanTanX, x)
		}
	})

	t.Run("asin(sin(x)) = x", func(t *testing.T) {
		xSmall := NewBigFloat(0.5, prec) // Use value in [-π/2, π/2]
		sinX := BigSin(xSmall, prec)
		asinSinX := BigAsin(sinX, prec)

		diff := new(BigFloat).SetPrec(prec).Sub(asinSinX, xSmall)
		diffFloat, _ := diff.Float64()

		if math.Abs(diffFloat) > 1e-10 {
			t.Errorf("asin(sin(x)) = %v, want %v", asinSinX, xSmall)
		}
	})
}

// TestTrigPrecisionLevels tests trigonometric functions at different precision levels
func TestTrigPrecisionLevels(t *testing.T) {
	x := 0.7
	precisions := []uint{64, 128, 256, 512}

	for _, prec := range precisions {
		t.Run("precision_"+string(rune(prec)), func(t *testing.T) {
			xBig := NewBigFloat(x, prec)

			// Test that functions don't panic and return reasonable values
			sinX := BigSin(xBig, prec)
			if sinX == nil {
				t.Error("BigSin returned nil")
			}

			cosX := BigCos(xBig, prec)
			if cosX == nil {
				t.Error("BigCos returned nil")
			}

			tanX := BigTan(xBig, prec)
			if tanX == nil {
				t.Error("BigTan returned nil")
			}
		})
	}
}
