// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

import (
	"math"
	"testing"
)

func TestExtendedPrecisionConstant(t *testing.T) {
	if ExtendedPrecision != 80 {
		t.Errorf("Expected ExtendedPrecision to be 80, got %d", ExtendedPrecision)
	}
}

func TestIsExtendedPrecisionMode(t *testing.T) {
	tests := []struct {
		prec   uint
		result bool
	}{
		{ExtendedPrecision, true},
		{80, true},
		{256, false},
		{64, false},
		{0, false},
	}

	for _, tt := range tests {
		result := IsExtendedPrecisionMode(tt.prec)
		if result != tt.result {
			t.Errorf("IsExtendedPrecisionMode(%d) = %v, want %v", tt.prec, result, tt.result)
		}
	}
}

func TestCanUseExtendedPrecision(t *testing.T) {
	features := GetCPUFeatures()
	expected := features.HasX87 && IsExtendedPrecisionMode(ExtendedPrecision)

	result := CanUseExtendedPrecision(ExtendedPrecision)
	if result != expected {
		t.Errorf("CanUseExtendedPrecision(%d) = %v, want %v (HasX87=%v)", ExtendedPrecision, result, expected, features.HasX87)
	}

	// Should return false for non-extended precision
	result = CanUseExtendedPrecision(256)
	if result {
		t.Errorf("CanUseExtendedPrecision(256) = true, want false")
	}
}

func TestBigFloatToExtendedFloat(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected float64
	}{
		{"zero", 0.0, 0.0},
		{"one", 1.0, 1.0},
		{"pi", math.Pi, math.Pi},
		{"negative", -1.0, -1.0},
		{"small", 1e-10, 1e-10},
		{"large", 1e10, 1e10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bf := NewBigFloat(tt.input, 256)
			result := BigFloatToExtendedFloat(bf)
			if math.Abs(result-tt.expected) > 1e-15 {
				t.Errorf("BigFloatToExtendedFloat(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestExtendedFloatToBigFloat(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		prec     uint
		expected float64
	}{
		{"zero", 0.0, 256, 0.0},
		{"one", 1.0, 256, 1.0},
		{"pi", math.Pi, 256, math.Pi},
		{"negative", -1.0, 256, -1.0},
		{"extended_precision", math.Pi, ExtendedPrecision, math.Pi},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtendedFloatToBigFloat(tt.input, tt.prec)
			resultFloat, _ := result.Float64()
			if math.Abs(resultFloat-tt.expected) > 1e-15 {
				t.Errorf("ExtendedFloatToBigFloat(%v, %d) = %v, want %v", tt.input, tt.prec, resultFloat, tt.expected)
			}
		})
	}
}

func TestExtendedPrecisionOperations(t *testing.T) {
	if !CanUseExtendedPrecision(ExtendedPrecision) {
		t.Skip("Extended precision not available on this platform")
	}

	tests := []struct {
		name     string
		op       func(*BigFloat, uint) *BigFloat
		input    float64
		expected float64
		tol      float64
	}{
		{"sin(0)", BigSin, 0.0, 0.0, 1e-15},
		{"sin(π/2)", BigSin, math.Pi / 2, 1.0, 1e-15},
		{"sin(π)", BigSin, math.Pi, 0.0, 1e-15},
		{"cos(0)", BigCos, 0.0, 1.0, 1e-15},
		{"cos(π/2)", BigCos, math.Pi / 2, 0.0, 1e-15},
		{"cos(π)", BigCos, math.Pi, -1.0, 1e-15},
		{"tan(0)", BigTan, 0.0, 0.0, 1e-15},
		{"tan(π/4)", BigTan, math.Pi / 4, 1.0, 1e-15},
		{"atan(0)", BigAtan, 0.0, 0.0, 1e-15},
		{"atan(1)", BigAtan, 1.0, math.Pi / 4, 1e-15},
		{"exp(0)", BigExp, 0.0, 1.0, 1e-15},
		{"exp(1)", BigExp, 1.0, math.E, 1e-15},
		{"log(1)", BigLog, 1.0, 0.0, 1e-15},
		{"log(e)", BigLog, math.E, 1.0, 1e-15},
		{"sqrt(1)", BigSqrt, 1.0, 1.0, 1e-15},
		{"sqrt(4)", BigSqrt, 4.0, 2.0, 1e-15},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := NewBigFloat(tt.input, ExtendedPrecision)
			result := tt.op(x, ExtendedPrecision)
			resultFloat, _ := result.Float64()
			diff := math.Abs(resultFloat - tt.expected)
			if diff > tt.tol {
				t.Errorf("%s(%v) = %v, want %v (diff: %v)", tt.name, tt.input, resultFloat, tt.expected, diff)
			}
		})
	}
}

func TestExtendedPrecisionAtan2(t *testing.T) {
	if !CanUseExtendedPrecision(ExtendedPrecision) {
		t.Skip("Extended precision not available on this platform")
	}

	tests := []struct {
		name     string
		y, x     float64
		expected float64
		tol      float64
	}{
		{"atan2(0, 1)", 0.0, 1.0, 0.0, 1e-15},
		{"atan2(1, 1)", 1.0, 1.0, math.Pi / 4, 1e-15},
		{"atan2(1, 0)", 1.0, 0.0, math.Pi / 2, 1e-15},
		{"atan2(-1, 0)", -1.0, 0.0, -math.Pi / 2, 1e-15},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			y := NewBigFloat(tt.y, ExtendedPrecision)
			x := NewBigFloat(tt.x, ExtendedPrecision)
			result := BigAtan2(y, x, ExtendedPrecision)
			resultFloat, _ := result.Float64()
			diff := math.Abs(resultFloat - tt.expected)
			if diff > tt.tol {
				t.Errorf("%s = %v, want %v (diff: %v)", tt.name, resultFloat, tt.expected, diff)
			}
		})
	}
}

func TestExtendedPrecisionPow(t *testing.T) {
	if !CanUseExtendedPrecision(ExtendedPrecision) {
		t.Skip("Extended precision not available on this platform")
	}

	tests := []struct {
		name     string
		x, y     float64
		expected float64
		tol      float64
	}{
		{"2^0", 2.0, 0.0, 1.0, 1e-15},
		{"2^1", 2.0, 1.0, 2.0, 1e-15},
		{"2^2", 2.0, 2.0, 4.0, 1e-15},
		{"2^0.5", 2.0, 0.5, math.Sqrt(2.0), 1e-15},
		{"e^1", math.E, 1.0, math.E, 1e-15},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := NewBigFloat(tt.x, ExtendedPrecision)
			y := NewBigFloat(tt.y, ExtendedPrecision)
			result := BigPow(x, y, ExtendedPrecision)
			resultFloat, _ := result.Float64()
			diff := math.Abs(resultFloat - tt.expected)
			if diff > tt.tol {
				t.Errorf("%s = %v, want %v (diff: %v)", tt.name, resultFloat, tt.expected, diff)
			}
		})
	}
}

func TestExtendedPrecisionFallback(t *testing.T) {
	// Test that operations fall back to BigFloat when extended precision is not requested
	x := NewBigFloat(math.Pi/4, 256)
	result := BigSin(x, 256)
	resultFloat, _ := result.Float64()
	expected := math.Sin(math.Pi / 4)
	diff := math.Abs(resultFloat - expected)
	if diff > 1e-10 {
		t.Errorf("BigSin with prec=256 should use BigFloat, got %v, want %v (diff: %v)", resultFloat, expected, diff)
	}
}

