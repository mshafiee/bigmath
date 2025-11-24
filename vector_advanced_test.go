// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

import (
	"math"
	"testing"
)

func TestBigVec3Cross(t *testing.T) {
	prec := uint(256)

	v1 := NewBigVec3(1.0, 0.0, 0.0, prec)
	v2 := NewBigVec3(0.0, 1.0, 0.0, prec)

	result := BigVec3Cross(v1, v2, prec)

	expected := [3]float64{0.0, 0.0, 1.0}
	got := result.ToFloat64()

	for i := 0; i < 3; i++ {
		if math.Abs(got[i]-expected[i]) > 1e-10 {
			t.Errorf("BigVec3Cross component[%d] = %g, want %g", i, got[i], expected[i])
		}
	}
}

func TestBigVec3Normalize(t *testing.T) {
	prec := uint(256)

	v := NewBigVec3(3.0, 4.0, 0.0, prec)
	normalized := BigVec3Normalize(v, prec)

	magnitude := BigVec3Magnitude(normalized, prec)
	mag, _ := magnitude.Float64()

	if math.Abs(mag-1.0) > 1e-10 {
		t.Errorf("BigVec3Normalize magnitude = %g, want 1.0", mag)
	}
}

func TestBigVec3Angle(t *testing.T) {
	prec := uint(256)

	v1 := NewBigVec3(1.0, 0.0, 0.0, prec)
	v2 := NewBigVec3(0.0, 1.0, 0.0, prec)

	angle := BigVec3Angle(v1, v2, prec)
	angleVal, _ := angle.Float64()

	expected := math.Pi / 2.0
	if math.Abs(angleVal-expected) > 1e-8 {
		t.Errorf("BigVec3Angle = %g, want %g", angleVal, expected)
	}
}

func TestBigVec3Project(t *testing.T) {
	prec := uint(256)

	v1 := NewBigVec3(3.0, 4.0, 0.0, prec)
	v2 := NewBigVec3(1.0, 0.0, 0.0, prec)

	projection := BigVec3Project(v1, v2, prec)
	proj := projection.ToFloat64()

	expected := [3]float64{3.0, 0.0, 0.0}
	for i := 0; i < 3; i++ {
		if math.Abs(proj[i]-expected[i]) > 1e-10 {
			t.Errorf("BigVec3Project component[%d] = %g, want %g", i, proj[i], expected[i])
		}
	}
}
