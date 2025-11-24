// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

import (
	"testing"
)

func TestBigFactorial(t *testing.T) {
	prec := uint(256)

	tests := []struct {
		name     string
		n        int64
		expected int64
	}{
		{"zero", 0, 1},
		{"one", 1, 1},
		{"five", 5, 120},
		{"ten", 10, 3628800},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BigFactorial(tt.n, prec)
			got, _ := result.Int64()
			if got != tt.expected {
				t.Errorf("BigFactorial(%d) = %d, want %d", tt.n, got, tt.expected)
			}
		})
	}
}

func TestBigBinomial(t *testing.T) {
	prec := uint(256)

	tests := []struct {
		name     string
		n, k     int64
		expected int64
	}{
		{"C(5,2)", 5, 2, 10},
		{"C(10,3)", 10, 3, 120},
		{"C(6,0)", 6, 0, 1},
		{"C(6,6)", 6, 6, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BigBinomial(tt.n, tt.k, prec)
			got, _ := result.Int64()
			if got != tt.expected {
				t.Errorf("BigBinomial(%d, %d) = %d, want %d", tt.n, tt.k, got, tt.expected)
			}
		})
	}
}
