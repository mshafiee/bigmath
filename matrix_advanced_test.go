// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

import (
	"math"
	"testing"
)

func TestBigMatTranspose(t *testing.T) {
	prec := uint(256)

	m := &BigMatrix3x3{
		M: [3][3]*BigFloat{
			{NewBigFloat(1.0, prec), NewBigFloat(2.0, prec), NewBigFloat(3.0, prec)},
			{NewBigFloat(4.0, prec), NewBigFloat(5.0, prec), NewBigFloat(6.0, prec)},
			{NewBigFloat(7.0, prec), NewBigFloat(8.0, prec), NewBigFloat(9.0, prec)},
		},
	}

	transposed := BigMatTranspose(m, prec)

	// Check that transpose worked
	if transposed.M[0][1].Cmp(m.M[1][0]) != 0 {
		t.Error("Transpose failed: M[0][1] != original M[1][0]")
	}
}

func TestBigMatMulMat(t *testing.T) {
	prec := uint(256)

	m1 := NewIdentityMatrix(prec)
	m2 := NewIdentityMatrix(prec)

	result := BigMatMulMat(m1, m2, prec)

	// Identity * Identity = Identity
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			expected := 0.0
			if i == j {
				expected = 1.0
			}
			got, _ := result.M[i][j].Float64()
			if math.Abs(got-expected) > 1e-10 {
				t.Errorf("BigMatMulMat[%d][%d] = %g, want %g", i, j, got, expected)
			}
		}
	}
}

func TestBigMatDet(t *testing.T) {
	prec := uint(256)

	// Identity matrix has determinant 1
	identity := NewIdentityMatrix(prec)
	det := BigMatDet(identity, prec)
	detVal, _ := det.Float64()

	if math.Abs(detVal-1.0) > 1e-10 {
		t.Errorf("BigMatDet(identity) = %g, want 1.0", detVal)
	}
}

func TestBigMatInverse(t *testing.T) {
	prec := uint(256)

	identity := NewIdentityMatrix(prec)
	inv, err := BigMatInverse(identity, prec)

	if err != nil {
		t.Fatalf("BigMatInverse failed: %v", err)
	}

	// Inverse of identity is identity
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			expected := 0.0
			if i == j {
				expected = 1.0
			}
			got, _ := inv.M[i][j].Float64()
			if math.Abs(got-expected) > 1e-10 {
				t.Errorf("BigMatInverse[%d][%d] = %g, want %g", i, j, got, expected)
			}
		}
	}
}
