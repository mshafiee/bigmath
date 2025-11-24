// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

import (
	"encoding/json"
	"testing"
)

func TestBigVec3JSON(t *testing.T) {
	prec := uint(256)

	v := NewBigVec3(1.0, 2.0, 3.0, prec)

	data, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var v2 BigVec3
	if err := json.Unmarshal(data, &v2); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	orig := v.ToFloat64()
	unmarshaled := v2.ToFloat64()

	for i := 0; i < 3; i++ {
		if orig[i] != unmarshaled[i] {
			t.Errorf("Component[%d] = %g, want %g", i, unmarshaled[i], orig[i])
		}
	}
}

func TestBigVec6JSON(t *testing.T) {
	prec := uint(256)

	v := NewBigVec6(1.0, 2.0, 3.0, 0.1, 0.2, 0.3, prec)

	data, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var v2 BigVec6
	if err := json.Unmarshal(data, &v2); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	orig := v.ToFloat64()
	unmarshaled := v2.ToFloat64()

	for i := 0; i < 6; i++ {
		if orig[i] != unmarshaled[i] {
			t.Errorf("Component[%d] = %g, want %g", i, unmarshaled[i], orig[i])
		}
	}
}

func TestBigMatrix3x3JSON(t *testing.T) {
	prec := uint(256)

	m := NewIdentityMatrix(prec)

	data, err := json.Marshal(m)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var m2 BigMatrix3x3
	if err := json.Unmarshal(data, &m2); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			expected := 0.0
			if i == j {
				expected = 1.0
			}
			got, _ := m2.M[i][j].Float64()
			if got != expected {
				t.Errorf("M[%d][%d] = %g, want %g", i, j, got, expected)
			}
		}
	}
}
