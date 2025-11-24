// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

import (
	"math"
	"testing"
)

// TestBigVec3Add tests BigVec3 addition
func TestBigVec3Add(t *testing.T) {
	prec := uint(256)
	v1 := NewBigVec3(1.0, 2.0, 3.0, prec)
	v2 := NewBigVec3(4.0, 5.0, 6.0, prec)

	result := BigVec3Add(v1, v2, prec)

	expected := [3]float64{5.0, 7.0, 9.0}
	actual := result.ToFloat64()

	for i := 0; i < 3; i++ {
		if math.Abs(actual[i]-expected[i]) > 1e-10 {
			t.Errorf("BigVec3Add component %d = %v, want %v", i, actual[i], expected[i])
		}
	}
}

// TestBigVec3Sub tests BigVec3 subtraction
func TestBigVec3Sub(t *testing.T) {
	prec := uint(256)
	v1 := NewBigVec3(4.0, 5.0, 6.0, prec)
	v2 := NewBigVec3(1.0, 2.0, 3.0, prec)

	result := BigVec3Sub(v1, v2, prec)

	expected := [3]float64{3.0, 3.0, 3.0}
	actual := result.ToFloat64()

	for i := 0; i < 3; i++ {
		if math.Abs(actual[i]-expected[i]) > 1e-10 {
			t.Errorf("BigVec3Sub component %d = %v, want %v", i, actual[i], expected[i])
		}
	}
}

// TestBigVec3Mul tests BigVec3 scalar multiplication
func TestBigVec3Mul(t *testing.T) {
	prec := uint(256)
	v := NewBigVec3(1.0, 2.0, 3.0, prec)
	scalar := NewBigFloat(2.0, prec)

	result := BigVec3Mul(v, scalar, prec)

	expected := [3]float64{2.0, 4.0, 6.0}
	actual := result.ToFloat64()

	for i := 0; i < 3; i++ {
		if math.Abs(actual[i]-expected[i]) > 1e-10 {
			t.Errorf("BigVec3Mul component %d = %v, want %v", i, actual[i], expected[i])
		}
	}
}

// TestBigVec3Copy tests BigVec3 copy functionality
func TestBigVec3Copy(t *testing.T) {
	prec := uint(256)
	v1 := NewBigVec3(1.0, 2.0, 3.0, prec)
	v2 := v1.Copy()

	// Verify values are equal
	v1Arr := v1.ToFloat64()
	v2Arr := v2.ToFloat64()

	for i := 0; i < 3; i++ {
		if v1Arr[i] != v2Arr[i] {
			t.Errorf("Copy component %d differs: %v != %v", i, v1Arr[i], v2Arr[i])
		}
	}

	// Modify v2 and verify v1 is unchanged
	v2.X = NewBigFloat(999.0, prec)
	v1X, _ := v1.X.Float64()
	if v1X != 1.0 {
		t.Error("Modifying copy affected original")
	}
}

// TestBigVec6Add tests BigVec6 addition
func TestBigVec6Add(t *testing.T) {
	prec := uint(256)
	v1 := NewBigVec6(1.0, 2.0, 3.0, 0.1, 0.2, 0.3, prec)
	v2 := NewBigVec6(4.0, 5.0, 6.0, 0.4, 0.5, 0.6, prec)

	result := BigVec6Add(v1, v2, prec)

	expected := [6]float64{5.0, 7.0, 9.0, 0.5, 0.7, 0.9}
	actual := result.ToFloat64()

	for i := 0; i < 6; i++ {
		if math.Abs(actual[i]-expected[i]) > 1e-10 {
			t.Errorf("BigVec6Add component %d = %v, want %v", i, actual[i], expected[i])
		}
	}
}

// TestBigVec6Sub tests BigVec6 subtraction
func TestBigVec6Sub(t *testing.T) {
	prec := uint(256)
	v1 := NewBigVec6(4.0, 5.0, 6.0, 0.4, 0.5, 0.6, prec)
	v2 := NewBigVec6(1.0, 2.0, 3.0, 0.1, 0.2, 0.3, prec)

	result := BigVec6Sub(v1, v2, prec)

	expected := [6]float64{3.0, 3.0, 3.0, 0.3, 0.3, 0.3}
	actual := result.ToFloat64()

	for i := 0; i < 6; i++ {
		if math.Abs(actual[i]-expected[i]) > 1e-10 {
			t.Errorf("BigVec6Sub component %d = %v, want %v", i, actual[i], expected[i])
		}
	}
}

// TestBigVec6Negate tests BigVec6 negation
func TestBigVec6Negate(t *testing.T) {
	prec := uint(256)
	v := NewBigVec6(1.0, -2.0, 3.0, -0.1, 0.2, -0.3, prec)

	result := BigVec6Negate(v, prec)

	expected := [6]float64{-1.0, 2.0, -3.0, 0.1, -0.2, 0.3}
	actual := result.ToFloat64()

	for i := 0; i < 6; i++ {
		if math.Abs(actual[i]-expected[i]) > 1e-10 {
			t.Errorf("BigVec6Negate component %d = %v, want %v", i, actual[i], expected[i])
		}
	}
}

// TestBigVec6Magnitude tests BigVec6 magnitude calculation
func TestBigVec6Magnitude(t *testing.T) {
	prec := uint(256)
	v := NewBigVec6(3.0, 4.0, 0.0, 0.0, 0.0, 0.0, prec)

	result := BigVec6Magnitude(v, prec)
	resultFloat, _ := result.Float64()

	expected := 5.0 // sqrt(3^2 + 4^2)
	if math.Abs(resultFloat-expected) > 1e-10 {
		t.Errorf("BigVec6Magnitude = %v, want %v", resultFloat, expected)
	}
}

// TestBigVec6Copy tests BigVec6 copy functionality
func TestBigVec6Copy(t *testing.T) {
	prec := uint(256)
	v1 := NewBigVec6(1.0, 2.0, 3.0, 0.1, 0.2, 0.3, prec)
	v2 := v1.Copy()

	// Verify values are equal
	v1Arr := v1.ToFloat64()
	v2Arr := v2.ToFloat64()

	for i := 0; i < 6; i++ {
		if v1Arr[i] != v2Arr[i] {
			t.Errorf("Copy component %d differs: %v != %v", i, v1Arr[i], v2Arr[i])
		}
	}

	// Modify v2 and verify v1 is unchanged
	v2.X = NewBigFloat(999.0, prec)
	v1X, _ := v1.X.Float64()
	if v1X != 1.0 {
		t.Error("Modifying copy affected original")
	}
}

// TestBigMatMul tests matrix-vector multiplication
func TestBigMatMul(t *testing.T) {
	prec := uint(256)

	// Identity matrix
	m := NewIdentityMatrix(prec)
	v := NewBigVec3(1.0, 2.0, 3.0, prec)

	result := BigMatMul(m, v, prec)
	expected := v.ToFloat64()
	actual := result.ToFloat64()

	for i := 0; i < 3; i++ {
		if math.Abs(actual[i]-expected[i]) > 1e-10 {
			t.Errorf("BigMatMul with identity: component %d = %v, want %v", i, actual[i], expected[i])
		}
	}

	// Scale matrix (2x diagonal)
	scale := NewBigFloat(2.0, prec)
	zero := NewBigFloat(0.0, prec)
	m2 := &BigMatrix3x3{
		M: [3][3]*BigFloat{
			{new(BigFloat).Set(scale), new(BigFloat).Set(zero), new(BigFloat).Set(zero)},
			{new(BigFloat).Set(zero), new(BigFloat).Set(scale), new(BigFloat).Set(zero)},
			{new(BigFloat).Set(zero), new(BigFloat).Set(zero), new(BigFloat).Set(scale)},
		},
	}

	result2 := BigMatMul(m2, v, prec)
	expected2 := [3]float64{2.0, 4.0, 6.0}
	actual2 := result2.ToFloat64()

	for i := 0; i < 3; i++ {
		if math.Abs(actual2[i]-expected2[i]) > 1e-10 {
			t.Errorf("BigMatMul with scale: component %d = %v, want %v", i, actual2[i], expected2[i])
		}
	}
}

// TestNewIdentityMatrix tests identity matrix creation
func TestNewIdentityMatrix(t *testing.T) {
	prec := uint(256)
	m := NewIdentityMatrix(prec)

	// Check diagonal is 1
	for i := 0; i < 3; i++ {
		val, _ := m.M[i][i].Float64()
		if math.Abs(val-1.0) > 1e-10 {
			t.Errorf("Identity matrix diagonal[%d] = %v, want 1.0", i, val)
		}
	}

	// Check off-diagonal is 0
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if i != j {
				val, _ := m.M[i][j].Float64()
				if math.Abs(val) > 1e-10 {
					t.Errorf("Identity matrix[%d][%d] = %v, want 0.0", i, j, val)
				}
			}
		}
	}

	// Test with default precision
	m2 := NewIdentityMatrix(0)
	if m2 == nil {
		t.Error("NewIdentityMatrix(0) returned nil")
	}
}

// TestApplyRotationMatrixToBigVec6 tests rotation matrix application
func TestApplyRotationMatrixToBigVec6(t *testing.T) {
	prec := uint(256)

	// Identity rotation should leave vector unchanged
	m := NewIdentityMatrix(prec)
	v := NewBigVec6(1.0, 2.0, 3.0, 0.1, 0.2, 0.3, prec)

	result := ApplyRotationMatrixToBigVec6(m, v, prec)
	expected := v.ToFloat64()
	actual := result.ToFloat64()

	for i := 0; i < 6; i++ {
		if math.Abs(actual[i]-expected[i]) > 1e-10 {
			t.Errorf("ApplyRotationMatrixToBigVec6 component %d = %v, want %v", i, actual[i], expected[i])
		}
	}
}

// TestCreateRotationMatrix tests rotation matrix creation
func TestCreateRotationMatrix(t *testing.T) {
	prec := uint(256)

	// Zero rotation should give identity-like behavior
	angles := [3]*BigFloat{
		NewBigFloat(0.0, prec),
		NewBigFloat(0.0, prec),
		NewBigFloat(0.0, prec),
	}

	m := CreateRotationMatrix(angles, prec)
	if m == nil {
		t.Error("CreateRotationMatrix returned nil")
	}

	// Test with π/2 rotation
	angles2 := [3]*BigFloat{
		NewBigFloat(math.Pi/2, prec),
		NewBigFloat(0.0, prec),
		NewBigFloat(0.0, prec),
	}

	m2 := CreateRotationMatrix(angles2, prec)
	if m2 == nil {
		t.Error("CreateRotationMatrix with π/2 returned nil")
	}

	// Verify it's a valid rotation matrix (orthogonal)
	// For now just check it's not nil and has proper structure
}

// TestBigMaxMin tests BigMax and BigMin functions
func TestBigMaxMin(t *testing.T) {
	prec := uint(256)
	a := NewBigFloat(3.5, prec)
	b := NewBigFloat(2.1, prec)

	t.Run("BigMax", func(t *testing.T) {
		result := BigMax(a, b, prec)
		resultFloat, _ := result.Float64()
		if resultFloat != 3.5 {
			t.Errorf("BigMax(3.5, 2.1) = %v, want 3.5", resultFloat)
		}

		result2 := BigMax(b, a, prec)
		result2Float, _ := result2.Float64()
		if result2Float != 3.5 {
			t.Errorf("BigMax(2.1, 3.5) = %v, want 3.5", result2Float)
		}

		// Test with equal values
		c := NewBigFloat(3.5, prec)
		result3 := BigMax(a, c, prec)
		result3Float, _ := result3.Float64()
		if result3Float != 3.5 {
			t.Errorf("BigMax(3.5, 3.5) = %v, want 3.5", result3Float)
		}

		// Test with default precision
		result4 := BigMax(a, b, 0)
		if result4 == nil {
			t.Error("BigMax with prec=0 returned nil")
		}
	})

	t.Run("BigMin", func(t *testing.T) {
		result := BigMin(a, b, prec)
		resultFloat, _ := result.Float64()
		if resultFloat != 2.1 {
			t.Errorf("BigMin(3.5, 2.1) = %v, want 2.1", resultFloat)
		}

		result2 := BigMin(b, a, prec)
		result2Float, _ := result2.Float64()
		if result2Float != 2.1 {
			t.Errorf("BigMin(2.1, 3.5) = %v, want 2.1", result2Float)
		}

		// Test with equal values
		c := NewBigFloat(2.1, prec)
		result3 := BigMin(b, c, prec)
		result3Float, _ := result3.Float64()
		if result3Float != 2.1 {
			t.Errorf("BigMin(2.1, 2.1) = %v, want 2.1", result3Float)
		}

		// Test with default precision
		result4 := BigMin(a, b, 0)
		if result4 == nil {
			t.Error("BigMin with prec=0 returned nil")
		}
	})
}

// TestBigFloatFMA tests Fused Multiply-Add
func TestBigFloatFMA(t *testing.T) {
	prec := uint(256)
	a := NewBigFloat(2.0, prec)
	b := NewBigFloat(3.0, prec)
	c := NewBigFloat(4.0, prec)

	// FMA: a*b + c = 2*3 + 4 = 10
	result := BigFloatFMA(a, b, c, prec)
	resultFloat, _ := result.Float64()

	expected := 10.0
	if math.Abs(resultFloat-expected) > 1e-10 {
		t.Errorf("BigFloatFMA(2, 3, 4) = %v, want %v", resultFloat, expected)
	}

	// Test with default precision
	result2 := BigFloatFMA(a, b, c, 0)
	if result2 == nil {
		t.Error("BigFloatFMA with prec=0 returned nil")
	}
}

// TestBigFloatDotProduct tests dot product calculation
func TestBigFloatDotProduct(t *testing.T) {
	prec := uint(256)

	v1 := []*BigFloat{
		NewBigFloat(1.0, prec),
		NewBigFloat(2.0, prec),
		NewBigFloat(3.0, prec),
	}
	v2 := []*BigFloat{
		NewBigFloat(4.0, prec),
		NewBigFloat(5.0, prec),
		NewBigFloat(6.0, prec),
	}

	// 1*4 + 2*5 + 3*6 = 4 + 10 + 18 = 32
	result := BigFloatDotProduct(v1, v2, prec)
	resultFloat, _ := result.Float64()

	expected := 32.0
	if math.Abs(resultFloat-expected) > 1e-10 {
		t.Errorf("BigFloatDotProduct = %v, want %v", resultFloat, expected)
	}

	// Test with empty vectors
	empty1 := []*BigFloat{}
	empty2 := []*BigFloat{}
	result2 := BigFloatDotProduct(empty1, empty2, prec)
	result2Float, _ := result2.Float64()
	if result2Float != 0.0 {
		t.Errorf("BigFloatDotProduct of empty vectors = %v, want 0.0", result2Float)
	}

	// Test with default precision
	result3 := BigFloatDotProduct(v1, v2, 0)
	if result3 == nil {
		t.Error("BigFloatDotProduct with prec=0 returned nil")
	}
}
