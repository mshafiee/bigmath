// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

import (
	"math"
	"testing"
)

func TestBigMatTranspose(t *testing.T) {
	prec := uint(256)

	tests := []struct {
		name string
		m    *BigMatrix3x3
	}{
		{"general_matrix", &BigMatrix3x3{
			M: [3][3]*BigFloat{
				{NewBigFloat(1.0, prec), NewBigFloat(2.0, prec), NewBigFloat(3.0, prec)},
				{NewBigFloat(4.0, prec), NewBigFloat(5.0, prec), NewBigFloat(6.0, prec)},
				{NewBigFloat(7.0, prec), NewBigFloat(8.0, prec), NewBigFloat(9.0, prec)},
			},
		}},
		{"identity_matrix", NewIdentityMatrix(prec)},
		{"zero_matrix", &BigMatrix3x3{
			M: [3][3]*BigFloat{
				{NewBigFloat(0.0, prec), NewBigFloat(0.0, prec), NewBigFloat(0.0, prec)},
				{NewBigFloat(0.0, prec), NewBigFloat(0.0, prec), NewBigFloat(0.0, prec)},
				{NewBigFloat(0.0, prec), NewBigFloat(0.0, prec), NewBigFloat(0.0, prec)},
			},
		}},
		{"symmetric_matrix", &BigMatrix3x3{
			M: [3][3]*BigFloat{
				{NewBigFloat(1.0, prec), NewBigFloat(2.0, prec), NewBigFloat(3.0, prec)},
				{NewBigFloat(2.0, prec), NewBigFloat(4.0, prec), NewBigFloat(5.0, prec)},
				{NewBigFloat(3.0, prec), NewBigFloat(5.0, prec), NewBigFloat(6.0, prec)},
			},
		}},
		{"very_large_values", &BigMatrix3x3{
			M: [3][3]*BigFloat{
				{NewBigFloat(1e10, prec), NewBigFloat(2e10, prec), NewBigFloat(3e10, prec)},
				{NewBigFloat(4e10, prec), NewBigFloat(5e10, prec), NewBigFloat(6e10, prec)},
				{NewBigFloat(7e10, prec), NewBigFloat(8e10, prec), NewBigFloat(9e10, prec)},
			},
		}},
		{"very_small_values", &BigMatrix3x3{
			M: [3][3]*BigFloat{
				{NewBigFloat(1e-10, prec), NewBigFloat(2e-10, prec), NewBigFloat(3e-10, prec)},
				{NewBigFloat(4e-10, prec), NewBigFloat(5e-10, prec), NewBigFloat(6e-10, prec)},
				{NewBigFloat(7e-10, prec), NewBigFloat(8e-10, prec), NewBigFloat(9e-10, prec)},
			},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transposed := BigMatTranspose(tt.m, prec)
			
			// Property: transpose(transpose(M)) = M
			transposedTwice := BigMatTranspose(transposed, prec)
			for i := 0; i < 3; i++ {
				for j := 0; j < 3; j++ {
					if transposedTwice.M[i][j].Cmp(tt.m.M[i][j]) != 0 {
						t.Errorf("Property violated: transpose(transpose(M))[%d][%d] != M[%d][%d]", i, j, i, j)
					}
				}
			}
			
			// Property: transpose(M)[i][j] = M[j][i]
			for i := 0; i < 3; i++ {
				for j := 0; j < 3; j++ {
					if transposed.M[i][j].Cmp(tt.m.M[j][i]) != 0 {
						t.Errorf("Transpose failed: transpose(M)[%d][%d] != M[%d][%d]", i, j, j, i)
					}
				}
			}
		})
	}
	
	// Test precision levels
	t.Run("precision_levels", func(t *testing.T) {
		testCases := []uint{64, 128, 256, 512}
		for _, p := range testCases {
			m := &BigMatrix3x3{
				M: [3][3]*BigFloat{
					{NewBigFloat(1.0, p), NewBigFloat(2.0, p), NewBigFloat(3.0, p)},
					{NewBigFloat(4.0, p), NewBigFloat(5.0, p), NewBigFloat(6.0, p)},
					{NewBigFloat(7.0, p), NewBigFloat(8.0, p), NewBigFloat(9.0, p)},
				},
			}
			transposed := BigMatTranspose(m, p)
			if transposed.M[0][1].Cmp(m.M[1][0]) != 0 {
				t.Errorf("Transpose failed at prec %d", p)
			}
		}
	})
}

func TestBigMatMulMat(t *testing.T) {
	prec := uint(256)

	tests := []struct {
		name     string
		m1, m2   *BigMatrix3x3
		expected *BigMatrix3x3
		tolerance float64
	}{
		{"identity_identity", NewIdentityMatrix(prec), NewIdentityMatrix(prec), NewIdentityMatrix(prec), 1e-10},
		{"identity_general", NewIdentityMatrix(prec), &BigMatrix3x3{
			M: [3][3]*BigFloat{
				{NewBigFloat(1.0, prec), NewBigFloat(2.0, prec), NewBigFloat(3.0, prec)},
				{NewBigFloat(4.0, prec), NewBigFloat(5.0, prec), NewBigFloat(6.0, prec)},
				{NewBigFloat(7.0, prec), NewBigFloat(8.0, prec), NewBigFloat(9.0, prec)},
			},
		}, &BigMatrix3x3{
			M: [3][3]*BigFloat{
				{NewBigFloat(1.0, prec), NewBigFloat(2.0, prec), NewBigFloat(3.0, prec)},
				{NewBigFloat(4.0, prec), NewBigFloat(5.0, prec), NewBigFloat(6.0, prec)},
				{NewBigFloat(7.0, prec), NewBigFloat(8.0, prec), NewBigFloat(9.0, prec)},
			},
		}, 1e-10},
		{"zero_matrix", &BigMatrix3x3{
			M: [3][3]*BigFloat{
				{NewBigFloat(0.0, prec), NewBigFloat(0.0, prec), NewBigFloat(0.0, prec)},
				{NewBigFloat(0.0, prec), NewBigFloat(0.0, prec), NewBigFloat(0.0, prec)},
				{NewBigFloat(0.0, prec), NewBigFloat(0.0, prec), NewBigFloat(0.0, prec)},
			},
		}, NewIdentityMatrix(prec), &BigMatrix3x3{
			M: [3][3]*BigFloat{
				{NewBigFloat(0.0, prec), NewBigFloat(0.0, prec), NewBigFloat(0.0, prec)},
				{NewBigFloat(0.0, prec), NewBigFloat(0.0, prec), NewBigFloat(0.0, prec)},
				{NewBigFloat(0.0, prec), NewBigFloat(0.0, prec), NewBigFloat(0.0, prec)},
			},
		}, 1e-10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BigMatMulMat(tt.m1, tt.m2, prec)
			
			for i := 0; i < 3; i++ {
				for j := 0; j < 3; j++ {
					got, _ := result.M[i][j].Float64()
					expVal, _ := tt.expected.M[i][j].Float64()
					if math.Abs(got-expVal) > tt.tolerance {
						t.Errorf("BigMatMulMat[%d][%d] = %g, want %g (tolerance %g)", i, j, got, expVal, tt.tolerance)
					}
				}
			}
		})
	}
	
	// Test general matrix multiplication
	t.Run("general_multiplication", func(t *testing.T) {
		m1 := &BigMatrix3x3{
			M: [3][3]*BigFloat{
				{NewBigFloat(1.0, prec), NewBigFloat(2.0, prec), NewBigFloat(3.0, prec)},
				{NewBigFloat(4.0, prec), NewBigFloat(5.0, prec), NewBigFloat(6.0, prec)},
				{NewBigFloat(7.0, prec), NewBigFloat(8.0, prec), NewBigFloat(9.0, prec)},
			},
		}
		m2 := &BigMatrix3x3{
			M: [3][3]*BigFloat{
				{NewBigFloat(9.0, prec), NewBigFloat(8.0, prec), NewBigFloat(7.0, prec)},
				{NewBigFloat(6.0, prec), NewBigFloat(5.0, prec), NewBigFloat(4.0, prec)},
				{NewBigFloat(3.0, prec), NewBigFloat(2.0, prec), NewBigFloat(1.0, prec)},
			},
		}
		result := BigMatMulMat(m1, m2, prec)
		
		// Expected result: [30, 24, 18; 84, 69, 54; 138, 114, 90]
		expected := [3][3]float64{
			{30.0, 24.0, 18.0},
			{84.0, 69.0, 54.0},
			{138.0, 114.0, 90.0},
		}
		
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				got, _ := result.M[i][j].Float64()
				if math.Abs(got-expected[i][j]) > 1e-8 {
					t.Errorf("BigMatMulMat[%d][%d] = %g, want %g", i, j, got, expected[i][j])
				}
			}
		}
	})
	
	// Test associativity: (A*B)*C = A*(B*C)
	t.Run("associativity", func(t *testing.T) {
		m1 := &BigMatrix3x3{
			M: [3][3]*BigFloat{
				{NewBigFloat(1.0, prec), NewBigFloat(2.0, prec), NewBigFloat(3.0, prec)},
				{NewBigFloat(4.0, prec), NewBigFloat(5.0, prec), NewBigFloat(6.0, prec)},
				{NewBigFloat(7.0, prec), NewBigFloat(8.0, prec), NewBigFloat(9.0, prec)},
			},
		}
		m2 := NewIdentityMatrix(prec)
		m3 := &BigMatrix3x3{
			M: [3][3]*BigFloat{
				{NewBigFloat(2.0, prec), NewBigFloat(0.0, prec), NewBigFloat(0.0, prec)},
				{NewBigFloat(0.0, prec), NewBigFloat(2.0, prec), NewBigFloat(0.0, prec)},
				{NewBigFloat(0.0, prec), NewBigFloat(0.0, prec), NewBigFloat(2.0, prec)},
			},
		}
		
		left := BigMatMulMat(BigMatMulMat(m1, m2, prec), m3, prec)
		right := BigMatMulMat(m1, BigMatMulMat(m2, m3, prec), prec)
		
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				if left.M[i][j].Cmp(right.M[i][j]) != 0 {
					leftVal, _ := left.M[i][j].Float64()
					rightVal, _ := right.M[i][j].Float64()
					t.Errorf("Associativity violated: (A*B)*C[%d][%d] = %g != A*(B*C)[%d][%d] = %g", i, j, leftVal, i, j, rightVal)
				}
			}
		}
	})
	
	// Test precision levels
	t.Run("precision_levels", func(t *testing.T) {
		testCases := []uint{64, 128, 256, 512}
		for _, p := range testCases {
			m1 := NewIdentityMatrix(p)
			m2 := NewIdentityMatrix(p)
			result := BigMatMulMat(m1, m2, p)
			for i := 0; i < 3; i++ {
				for j := 0; j < 3; j++ {
					expected := 0.0
					if i == j {
						expected = 1.0
					}
					got, _ := result.M[i][j].Float64()
					if math.Abs(got-expected) > 1e-6 {
						t.Errorf("BigMatMulMat at prec %d[%d][%d] = %g, want %g", p, i, j, got, expected)
					}
				}
			}
		}
	})
}

func TestBigMatDet(t *testing.T) {
	prec := uint(256)

	tests := []struct {
		name      string
		m         *BigMatrix3x3
		expected  float64
		tolerance float64
	}{
		{"identity", NewIdentityMatrix(prec), 1.0, 1e-10},
		{"zero_matrix", &BigMatrix3x3{
			M: [3][3]*BigFloat{
				{NewBigFloat(0.0, prec), NewBigFloat(0.0, prec), NewBigFloat(0.0, prec)},
				{NewBigFloat(0.0, prec), NewBigFloat(0.0, prec), NewBigFloat(0.0, prec)},
				{NewBigFloat(0.0, prec), NewBigFloat(0.0, prec), NewBigFloat(0.0, prec)},
			},
		}, 0.0, 1e-10},
		{"singular_matrix", &BigMatrix3x3{
			M: [3][3]*BigFloat{
				{NewBigFloat(1.0, prec), NewBigFloat(2.0, prec), NewBigFloat(3.0, prec)},
				{NewBigFloat(2.0, prec), NewBigFloat(4.0, prec), NewBigFloat(6.0, prec)},
				{NewBigFloat(3.0, prec), NewBigFloat(6.0, prec), NewBigFloat(9.0, prec)},
			},
		}, 0.0, 1e-8},
		{"diagonal_matrix", &BigMatrix3x3{
			M: [3][3]*BigFloat{
				{NewBigFloat(2.0, prec), NewBigFloat(0.0, prec), NewBigFloat(0.0, prec)},
				{NewBigFloat(0.0, prec), NewBigFloat(3.0, prec), NewBigFloat(0.0, prec)},
				{NewBigFloat(0.0, prec), NewBigFloat(0.0, prec), NewBigFloat(4.0, prec)},
			},
		}, 24.0, 1e-10},
		{"general_matrix", &BigMatrix3x3{
			M: [3][3]*BigFloat{
				{NewBigFloat(1.0, prec), NewBigFloat(2.0, prec), NewBigFloat(3.0, prec)},
				{NewBigFloat(4.0, prec), NewBigFloat(5.0, prec), NewBigFloat(6.0, prec)},
				{NewBigFloat(7.0, prec), NewBigFloat(8.0, prec), NewBigFloat(9.0, prec)},
			},
		}, 0.0, 1e-8}, // This matrix is singular (rows are linearly dependent)
		{"non_singular", &BigMatrix3x3{
			M: [3][3]*BigFloat{
				{NewBigFloat(1.0, prec), NewBigFloat(0.0, prec), NewBigFloat(0.0, prec)},
				{NewBigFloat(0.0, prec), NewBigFloat(1.0, prec), NewBigFloat(0.0, prec)},
				{NewBigFloat(1.0, prec), NewBigFloat(1.0, prec), NewBigFloat(1.0, prec)},
			},
		}, 1.0, 1e-10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			det := BigMatDet(tt.m, prec)
			detVal, _ := det.Float64()
			if math.Abs(detVal-tt.expected) > tt.tolerance {
				t.Errorf("BigMatDet = %g, want %g (tolerance %g)", detVal, tt.expected, tt.tolerance)
			}
		})
	}
	
	// Property: det(A*B) = det(A) * det(B)
	t.Run("multiplicative_property", func(t *testing.T) {
		m1 := &BigMatrix3x3{
			M: [3][3]*BigFloat{
				{NewBigFloat(1.0, prec), NewBigFloat(2.0, prec), NewBigFloat(0.0, prec)},
				{NewBigFloat(0.0, prec), NewBigFloat(1.0, prec), NewBigFloat(0.0, prec)},
				{NewBigFloat(0.0, prec), NewBigFloat(0.0, prec), NewBigFloat(1.0, prec)},
			},
		}
		m2 := &BigMatrix3x3{
			M: [3][3]*BigFloat{
				{NewBigFloat(2.0, prec), NewBigFloat(0.0, prec), NewBigFloat(0.0, prec)},
				{NewBigFloat(0.0, prec), NewBigFloat(3.0, prec), NewBigFloat(0.0, prec)},
				{NewBigFloat(0.0, prec), NewBigFloat(0.0, prec), NewBigFloat(4.0, prec)},
			},
		}
		
		det1 := BigMatDet(m1, prec)
		det2 := BigMatDet(m2, prec)
		product := BigMatMulMat(m1, m2, prec)
		detProduct := BigMatDet(product, prec)
		
		expected := new(BigFloat).SetPrec(prec).Mul(det1, det2)
		if detProduct.Cmp(expected) != 0 {
			det1Val, _ := det1.Float64()
			det2Val, _ := det2.Float64()
			detProductVal, _ := detProduct.Float64()
			expectedVal, _ := expected.Float64()
			t.Errorf("Property violated: det(A*B) = %g != det(A)*det(B) = %g * %g = %g", detProductVal, det1Val, det2Val, expectedVal)
		}
	})
	
	// Property: det(transpose(A)) = det(A)
	t.Run("transpose_property", func(t *testing.T) {
		m := &BigMatrix3x3{
			M: [3][3]*BigFloat{
				{NewBigFloat(1.0, prec), NewBigFloat(2.0, prec), NewBigFloat(3.0, prec)},
				{NewBigFloat(4.0, prec), NewBigFloat(5.0, prec), NewBigFloat(6.0, prec)},
				{NewBigFloat(7.0, prec), NewBigFloat(8.0, prec), NewBigFloat(10.0, prec)},
			},
		}
		det := BigMatDet(m, prec)
		transposed := BigMatTranspose(m, prec)
		detTransposed := BigMatDet(transposed, prec)
		
		if det.Cmp(detTransposed) != 0 {
			detVal, _ := det.Float64()
			detTransposedVal, _ := detTransposed.Float64()
			t.Errorf("Property violated: det(transpose(A)) = %g != det(A) = %g", detTransposedVal, detVal)
		}
	})
	
	// Test precision levels
	t.Run("precision_levels", func(t *testing.T) {
		testCases := []uint{64, 128, 256, 512}
		for _, p := range testCases {
			identity := NewIdentityMatrix(p)
			det := BigMatDet(identity, p)
			detVal, _ := det.Float64()
			if math.Abs(detVal-1.0) > 1e-6 {
				t.Errorf("BigMatDet(identity) at prec %d = %g, want 1.0", p, detVal)
			}
		}
	})
}

func TestBigMatInverse(t *testing.T) {
	prec := uint(256)

	tests := []struct {
		name      string
		m         *BigMatrix3x3
		shouldErr bool
		tolerance float64
	}{
		{"identity", NewIdentityMatrix(prec), false, 1e-10},
		{"diagonal", &BigMatrix3x3{
			M: [3][3]*BigFloat{
				{NewBigFloat(2.0, prec), NewBigFloat(0.0, prec), NewBigFloat(0.0, prec)},
				{NewBigFloat(0.0, prec), NewBigFloat(3.0, prec), NewBigFloat(0.0, prec)},
				{NewBigFloat(0.0, prec), NewBigFloat(0.0, prec), NewBigFloat(4.0, prec)},
			},
		}, false, 1e-10},
		{"singular_matrix", &BigMatrix3x3{
			M: [3][3]*BigFloat{
				{NewBigFloat(1.0, prec), NewBigFloat(2.0, prec), NewBigFloat(3.0, prec)},
				{NewBigFloat(2.0, prec), NewBigFloat(4.0, prec), NewBigFloat(6.0, prec)},
				{NewBigFloat(3.0, prec), NewBigFloat(6.0, prec), NewBigFloat(9.0, prec)},
			},
		}, true, 0},
		{"zero_matrix", &BigMatrix3x3{
			M: [3][3]*BigFloat{
				{NewBigFloat(0.0, prec), NewBigFloat(0.0, prec), NewBigFloat(0.0, prec)},
				{NewBigFloat(0.0, prec), NewBigFloat(0.0, prec), NewBigFloat(0.0, prec)},
				{NewBigFloat(0.0, prec), NewBigFloat(0.0, prec), NewBigFloat(0.0, prec)},
			},
		}, true, 0},
		{"non_singular", &BigMatrix3x3{
			M: [3][3]*BigFloat{
				{NewBigFloat(1.0, prec), NewBigFloat(0.0, prec), NewBigFloat(0.0, prec)},
				{NewBigFloat(0.0, prec), NewBigFloat(1.0, prec), NewBigFloat(0.0, prec)},
				{NewBigFloat(1.0, prec), NewBigFloat(1.0, prec), NewBigFloat(1.0, prec)},
			},
		}, false, 1e-10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inv, err := BigMatInverse(tt.m, prec)
			
			if tt.shouldErr {
				if err == nil {
					t.Errorf("BigMatInverse should have failed for singular matrix")
				}
				return
			}
			
			if err != nil {
				t.Fatalf("BigMatInverse failed: %v", err)
			}
			
			// Property: M * M^-1 = I
			product := BigMatMulMat(tt.m, inv, prec)
			identity := NewIdentityMatrix(prec)
			for i := 0; i < 3; i++ {
				for j := 0; j < 3; j++ {
					got, _ := product.M[i][j].Float64()
					expected, _ := identity.M[i][j].Float64()
					if math.Abs(got-expected) > tt.tolerance {
						t.Errorf("Property violated: M * M^-1[%d][%d] = %g, want %g", i, j, got, expected)
					}
				}
			}
			
			// Property: M^-1 * M = I
			product2 := BigMatMulMat(inv, tt.m, prec)
			for i := 0; i < 3; i++ {
				for j := 0; j < 3; j++ {
					got, _ := product2.M[i][j].Float64()
					expected, _ := identity.M[i][j].Float64()
					if math.Abs(got-expected) > tt.tolerance {
						t.Errorf("Property violated: M^-1 * M[%d][%d] = %g, want %g", i, j, got, expected)
					}
				}
			}
		})
	}
	
	// Test precision levels
	t.Run("precision_levels", func(t *testing.T) {
		testCases := []uint{64, 128, 256, 512}
		for _, p := range testCases {
			identity := NewIdentityMatrix(p)
			inv, err := BigMatInverse(identity, p)
			if err != nil {
				t.Fatalf("BigMatInverse failed at prec %d: %v", p, err)
			}
			for i := 0; i < 3; i++ {
				for j := 0; j < 3; j++ {
					expected := 0.0
					if i == j {
						expected = 1.0
					}
					got, _ := inv.M[i][j].Float64()
					if math.Abs(got-expected) > 1e-6 {
						t.Errorf("BigMatInverse at prec %d[%d][%d] = %g, want %g", p, i, j, got, expected)
					}
				}
			}
		}
	})
}
