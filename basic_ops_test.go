// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

import (
	"math"
	"testing"
)

func TestBigFloor(t *testing.T) {
	prec := uint(256)

	tests := []struct {
		name     string
		input    float64
		expected float64
	}{
		{"positive", 3.7, 3.0},
		{"negative", -3.7, -4.0},
		{"integer", 5.0, 5.0},
		{"zero", 0.0, 0.0},
		{"small", 0.1, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := NewBigFloat(tt.input, prec)
			result := BigFloor(x, prec)
			got, _ := result.Float64()
			if math.Abs(got-tt.expected) > 1e-10 {
				t.Errorf("BigFloor(%g) = %g, want %g", tt.input, got, tt.expected)
			}
		})
	}
}

func TestBigCeil(t *testing.T) {
	prec := uint(256)

	tests := []struct {
		name     string
		input    float64
		expected float64
	}{
		{"positive", 3.2, 4.0},
		{"negative", -3.2, -3.0},
		{"integer", 5.0, 5.0},
		{"zero", 0.0, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := NewBigFloat(tt.input, prec)
			result := BigCeil(x, prec)
			got, _ := result.Float64()
			if math.Abs(got-tt.expected) > 1e-10 {
				t.Errorf("BigCeil(%g) = %g, want %g", tt.input, got, tt.expected)
			}
		})
	}
}

func TestBigTrunc(t *testing.T) {
	prec := uint(256)

	tests := []struct {
		name     string
		input    float64
		expected float64
	}{
		{"positive", 3.7, 3.0},
		{"negative", -3.7, -3.0},
		{"integer", 5.0, 5.0},
		{"zero", 0.0, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := NewBigFloat(tt.input, prec)
			result := BigTrunc(x, prec)
			got, _ := result.Float64()
			if math.Abs(got-tt.expected) > 1e-10 {
				t.Errorf("BigTrunc(%g) = %g, want %g", tt.input, got, tt.expected)
			}
		})
	}
}

func TestBigMod(t *testing.T) {
	prec := uint(256)

	tests := []struct {
		name     string
		x, y     float64
		expected float64
		checkNaN bool
	}{
		{"basic", 10.0, 3.0, 1.0, false},
		{"negative", -10.0, 3.0, 2.0, false},
		{"zero_y", 10.0, 0.0, 0.0, true}, // big.Float doesn't support NaN, returns 0
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := NewBigFloat(tt.x, prec)
			y := NewBigFloat(tt.y, prec)
			result := BigMod(x, y, prec)
			got, _ := result.Float64()
			if tt.checkNaN {
				// big.Float doesn't support NaN, so division by zero returns 0
				// This is a limitation we document
				if got != 0.0 {
					t.Errorf("BigMod(%g, %g) = %g, want 0 (big.Float limitation: no NaN support)", tt.x, tt.y, got)
				}
			} else if math.Abs(got-tt.expected) > 1e-10 {
				t.Errorf("BigMod(%g, %g) = %g, want %g", tt.x, tt.y, got, tt.expected)
			}
		})
	}
}

func TestBigRem(t *testing.T) {
	prec := uint(256)

	tests := []struct {
		name     string
		x, y     float64
		expected float64
	}{
		{"basic", 10.0, 3.0, 1.0},
		{"negative_x", -10.0, 3.0, -1.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := NewBigFloat(tt.x, prec)
			y := NewBigFloat(tt.y, prec)
			result := BigRem(x, y, prec)
			got, _ := result.Float64()
			if math.Abs(got-tt.expected) > 1e-10 {
				t.Errorf("BigRem(%g, %g) = %g, want %g", tt.x, tt.y, got, tt.expected)
			}
		})
	}
}
