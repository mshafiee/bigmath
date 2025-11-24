// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

import (
	"encoding/json"
	"testing"
)

func TestBigVec3JSON(t *testing.T) {
	prec := uint(256)

	tests := []struct {
		name      string
		v         *BigVec3
		tolerance float64
	}{
		{"normal_values", NewBigVec3(1.0, 2.0, 3.0, prec), 1e-10},
		{"zero_vector", NewBigVec3(0.0, 0.0, 0.0, prec), 1e-10},
		{"negative_values", NewBigVec3(-1.0, -2.0, -3.0, prec), 1e-10},
		{"very_large", NewBigVec3(1e10, 2e10, 3e10, prec), 1e5},
		{"very_small", NewBigVec3(1e-10, 2e-10, 3e-10, prec), 1e-20},
		{"mixed_signs", NewBigVec3(1.0, -2.0, 3.0, prec), 1e-10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.v)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}

			var v2 BigVec3
			if err := json.Unmarshal(data, &v2); err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}

			orig := tt.v.ToFloat64()
			unmarshaled := v2.ToFloat64()

			for i := 0; i < 3; i++ {
				diff := orig[i] - unmarshaled[i]
				if diff < 0 {
					diff = -diff
				}
				if diff > tt.tolerance {
					t.Errorf("Component[%d] = %g, want %g (diff %g, tolerance %g)", i, unmarshaled[i], orig[i], diff, tt.tolerance)
				}
			}

			// Test round-trip accuracy
			data2, err := json.Marshal(&v2)
			if err != nil {
				t.Fatalf("Second Marshal failed: %v", err)
			}

			var v3 BigVec3
			if err := json.Unmarshal(data2, &v3); err != nil {
				t.Fatalf("Second Unmarshal failed: %v", err)
			}

			unmarshaled2 := v3.ToFloat64()
			for i := 0; i < 3; i++ {
				diff := unmarshaled[i] - unmarshaled2[i]
				if diff < 0 {
					diff = -diff
				}
				if diff > tt.tolerance {
					t.Errorf("Round-trip failed: Component[%d] = %g, want %g", i, unmarshaled2[i], unmarshaled[i])
				}
			}
		})
	}

	// Test invalid JSON
	t.Run("invalid_json", func(t *testing.T) {
		invalidJSON := []byte(`{"X": "not a number", "Y": 2.0, "Z": 3.0}`)
		var v BigVec3
		if err := json.Unmarshal(invalidJSON, &v); err == nil {
			t.Error("Unmarshal should fail for invalid JSON")
		}
	})

	// Test precision preservation
	t.Run("precision_preservation", func(t *testing.T) {
		testCases := []uint{64, 128, 256, 512}
		for _, p := range testCases {
			v := NewBigVec3(1.0, 2.0, 3.0, p)
			data, err := json.Marshal(v)
			if err != nil {
				t.Fatalf("Marshal failed at prec %d: %v", p, err)
			}

			var v2 BigVec3
			if err := json.Unmarshal(data, &v2); err != nil {
				t.Fatalf("Unmarshal failed at prec %d: %v", p, err)
			}

			// Check that precision is preserved (at least approximately)
			origPrec := v.X.Prec()
			unmarshaledPrec := v2.X.Prec()
			if unmarshaledPrec < origPrec/2 {
				t.Errorf("Precision not preserved: orig %d, unmarshaled %d", origPrec, unmarshaledPrec)
			}
		}
	})
}

func TestBigVec6JSON(t *testing.T) {
	prec := uint(256)

	tests := []struct {
		name      string
		v         *BigVec6
		tolerance float64
	}{
		{"normal_values", NewBigVec6(1.0, 2.0, 3.0, 0.1, 0.2, 0.3, prec), 1e-10},
		{"zero_vector", NewBigVec6(0.0, 0.0, 0.0, 0.0, 0.0, 0.0, prec), 1e-10},
		{"negative_values", NewBigVec6(-1.0, -2.0, -3.0, -0.1, -0.2, -0.3, prec), 1e-10},
		{"very_large", NewBigVec6(1e10, 2e10, 3e10, 4e10, 5e10, 6e10, prec), 1e5},
		{"very_small", NewBigVec6(1e-10, 2e-10, 3e-10, 4e-10, 5e-10, 6e-10, prec), 1e-20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.v)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}

			var v2 BigVec6
			if err := json.Unmarshal(data, &v2); err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}

			orig := tt.v.ToFloat64()
			unmarshaled := v2.ToFloat64()

			for i := 0; i < 6; i++ {
				diff := orig[i] - unmarshaled[i]
				if diff < 0 {
					diff = -diff
				}
				if diff > tt.tolerance {
					t.Errorf("Component[%d] = %g, want %g (diff %g, tolerance %g)", i, unmarshaled[i], orig[i], diff, tt.tolerance)
				}
			}
		})
	}

	// Test invalid JSON
	t.Run("invalid_json", func(t *testing.T) {
		invalidJSON := []byte(`{"X": "not a number", "Y": 2.0}`)
		var v BigVec6
		if err := json.Unmarshal(invalidJSON, &v); err == nil {
			t.Error("Unmarshal should fail for invalid JSON")
		}
	})
}

func TestBigMatrix3x3JSON(t *testing.T) {
	prec := uint(256)

	tests := []struct {
		name      string
		m         *BigMatrix3x3
		tolerance float64
	}{
		{"identity", NewIdentityMatrix(prec), 1e-10},
		{"zero_matrix", &BigMatrix3x3{
			M: [3][3]*BigFloat{
				{NewBigFloat(0.0, prec), NewBigFloat(0.0, prec), NewBigFloat(0.0, prec)},
				{NewBigFloat(0.0, prec), NewBigFloat(0.0, prec), NewBigFloat(0.0, prec)},
				{NewBigFloat(0.0, prec), NewBigFloat(0.0, prec), NewBigFloat(0.0, prec)},
			},
		}, 1e-10},
		{"general_matrix", &BigMatrix3x3{
			M: [3][3]*BigFloat{
				{NewBigFloat(1.0, prec), NewBigFloat(2.0, prec), NewBigFloat(3.0, prec)},
				{NewBigFloat(4.0, prec), NewBigFloat(5.0, prec), NewBigFloat(6.0, prec)},
				{NewBigFloat(7.0, prec), NewBigFloat(8.0, prec), NewBigFloat(9.0, prec)},
			},
		}, 1e-10},
		{"very_large", &BigMatrix3x3{
			M: [3][3]*BigFloat{
				{NewBigFloat(1e10, prec), NewBigFloat(2e10, prec), NewBigFloat(3e10, prec)},
				{NewBigFloat(4e10, prec), NewBigFloat(5e10, prec), NewBigFloat(6e10, prec)},
				{NewBigFloat(7e10, prec), NewBigFloat(8e10, prec), NewBigFloat(9e10, prec)},
			},
		}, 1e5},
		{"very_small", &BigMatrix3x3{
			M: [3][3]*BigFloat{
				{NewBigFloat(1e-10, prec), NewBigFloat(2e-10, prec), NewBigFloat(3e-10, prec)},
				{NewBigFloat(4e-10, prec), NewBigFloat(5e-10, prec), NewBigFloat(6e-10, prec)},
				{NewBigFloat(7e-10, prec), NewBigFloat(8e-10, prec), NewBigFloat(9e-10, prec)},
			},
		}, 1e-20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.m)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}

			var m2 BigMatrix3x3
			if err := json.Unmarshal(data, &m2); err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}

			for i := 0; i < 3; i++ {
				for j := 0; j < 3; j++ {
					orig, _ := tt.m.M[i][j].Float64()
					got, _ := m2.M[i][j].Float64()
					diff := orig - got
					if diff < 0 {
						diff = -diff
					}
					if diff > tt.tolerance {
						t.Errorf("M[%d][%d] = %g, want %g (diff %g, tolerance %g)", i, j, got, orig, diff, tt.tolerance)
					}
				}
			}
		})
	}

	// Test invalid JSON
	t.Run("invalid_json", func(t *testing.T) {
		invalidJSON := []byte(`{"M": [["not a number", 2.0, 3.0], [4.0, 5.0, 6.0], [7.0, 8.0, 9.0]]}`)
		var m BigMatrix3x3
		if err := json.Unmarshal(invalidJSON, &m); err == nil {
			t.Error("Unmarshal should fail for invalid JSON")
		}
	})

	// Test precision preservation
	t.Run("precision_preservation", func(t *testing.T) {
		testCases := []uint{64, 128, 256, 512}
		for _, p := range testCases {
			m := NewIdentityMatrix(p)
			data, err := json.Marshal(m)
			if err != nil {
				t.Fatalf("Marshal failed at prec %d: %v", p, err)
			}

			var m2 BigMatrix3x3
			if err := json.Unmarshal(data, &m2); err != nil {
				t.Fatalf("Unmarshal failed at prec %d: %v", p, err)
			}

			// Check that precision is preserved (at least approximately)
			origPrec := m.M[0][0].Prec()
			unmarshaledPrec := m2.M[0][0].Prec()
			if unmarshaledPrec < origPrec/2 {
				t.Errorf("Precision not preserved: orig %d, unmarshaled %d", origPrec, unmarshaledPrec)
			}
		}
	})
}
