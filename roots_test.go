// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

import (
	"math"
	"testing"
)

func TestBigCbrt(t *testing.T) {
	prec := uint(256)

	tests := []struct {
		name     string
		input    float64
		expected float64
	}{
		{"positive", 8.0, 2.0},
		{"one", 1.0, 1.0},
		{"zero", 0.0, 0.0},
		{"negative", -8.0, -2.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := NewBigFloat(tt.input, prec)
			result := BigCbrt(x, prec)
			got, _ := result.Float64()
			if math.Abs(got-tt.expected) > 1e-10 {
				t.Errorf("BigCbrt(%g) = %g, want %g", tt.input, got, tt.expected)
			}
		})
	}
}

func TestBigRoot(t *testing.T) {
	prec := uint(256)

	tests := []struct {
		name     string
		n, x     float64
		expected float64
	}{
		{"square_root", 2.0, 4.0, 2.0},
		{"fourth_root", 4.0, 16.0, 2.0},
		{"one", 5.0, 1.0, 1.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NewBigFloat(tt.n, prec)
			x := NewBigFloat(tt.x, prec)
			result := BigRoot(n, x, prec)
			got, _ := result.Float64()
			expected := math.Pow(tt.x, 1.0/tt.n)
			if math.Abs(got-expected) > 1e-8 {
				t.Errorf("BigRoot(%g, %g) = %g, want %g", tt.n, tt.x, got, expected)
			}
		})
	}
}
