// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

import (
	"errors"
)

// BigMatTranspose returns the transpose of a 3x3 matrix
func BigMatTranspose(m *BigMatrix3x3, prec uint) *BigMatrix3x3 {
	if prec == 0 {
		prec = m.M[0][0].Prec()
	}

	return &BigMatrix3x3{
		M: [3][3]*BigFloat{
			{new(BigFloat).SetPrec(prec).Set(m.M[0][0]), new(BigFloat).SetPrec(prec).Set(m.M[1][0]), new(BigFloat).SetPrec(prec).Set(m.M[2][0])},
			{new(BigFloat).SetPrec(prec).Set(m.M[0][1]), new(BigFloat).SetPrec(prec).Set(m.M[1][1]), new(BigFloat).SetPrec(prec).Set(m.M[2][1])},
			{new(BigFloat).SetPrec(prec).Set(m.M[0][2]), new(BigFloat).SetPrec(prec).Set(m.M[1][2]), new(BigFloat).SetPrec(prec).Set(m.M[2][2])},
		},
	}
}

// BigMatMulMat multiplies two 3x3 matrices: result = m1 * m2
func BigMatMulMat(m1, m2 *BigMatrix3x3, prec uint) *BigMatrix3x3 {
	if prec == 0 {
		prec = m1.M[0][0].Prec()
	}

	result := &BigMatrix3x3{
		M: [3][3]*BigFloat{},
	}

	// Compute each element: result[i][j] = sum_k(m1[i][k] * m2[k][j])
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			sum := NewBigFloat(0.0, prec)
			for k := 0; k < 3; k++ {
				product := new(BigFloat).SetPrec(prec).Mul(m1.M[i][k], m2.M[k][j])
				sum.Add(sum, product)
			}
			result.M[i][j] = sum
		}
	}

	return result
}

// BigMatDet computes the determinant of a 3x3 matrix
// Uses the formula: det = a(ei-fh) - b(di-fg) + c(dh-eg)
// where matrix is:
//
//	[a b c]
//	[d e f]
//	[g h i]
func BigMatDet(m *BigMatrix3x3, prec uint) *BigFloat {
	if prec == 0 {
		prec = m.M[0][0].Prec()
	}

	// Extract elements for clarity
	a, b, c := m.M[0][0], m.M[0][1], m.M[0][2]
	d, e, f := m.M[1][0], m.M[1][1], m.M[1][2]
	g, h, i := m.M[2][0], m.M[2][1], m.M[2][2]

	// Compute sub-determinants
	// ei - fh
	ei := new(BigFloat).SetPrec(prec).Mul(e, i)
	fh := new(BigFloat).SetPrec(prec).Mul(f, h)
	term1 := new(BigFloat).SetPrec(prec).Sub(ei, fh)

	// di - fg
	di := new(BigFloat).SetPrec(prec).Mul(d, i)
	fg := new(BigFloat).SetPrec(prec).Mul(f, g)
	term2 := new(BigFloat).SetPrec(prec).Sub(di, fg)

	// dh - eg
	dh := new(BigFloat).SetPrec(prec).Mul(d, h)
	eg := new(BigFloat).SetPrec(prec).Mul(e, g)
	term3 := new(BigFloat).SetPrec(prec).Sub(dh, eg)

	// Compute determinant: a*term1 - b*term2 + c*term3
	aTerm1 := new(BigFloat).SetPrec(prec).Mul(a, term1)
	bTerm2 := new(BigFloat).SetPrec(prec).Mul(b, term2)
	cTerm3 := new(BigFloat).SetPrec(prec).Mul(c, term3)

	result := new(BigFloat).SetPrec(prec).Sub(aTerm1, bTerm2)
	result.Add(result, cTerm3)

	return result
}

// BigMatInverse computes the inverse of a 3x3 matrix
// Returns error if matrix is singular (determinant is zero)
// Uses adjugate matrix: M^-1 = (1/det(M)) * adj(M)
func BigMatInverse(m *BigMatrix3x3, prec uint) (*BigMatrix3x3, error) {
	if prec == 0 {
		prec = m.M[0][0].Prec()
	}

	// Compute determinant
	det := BigMatDet(m, prec)

	// Check if matrix is singular
	zero := NewBigFloat(0.0, prec)
	if det.Cmp(zero) == 0 {
		return nil, errors.New("matrix is singular (determinant is zero)")
	}

	// Compute adjugate matrix (transpose of cofactor matrix)
	// For 3x3, cofactors are:
	// C[0][0] = +(e*i - f*h), C[0][1] = -(d*i - f*g), C[0][2] = +(d*h - e*g)
	// C[1][0] = -(b*i - c*h), C[1][1] = +(a*i - c*g), C[1][2] = -(a*h - b*g)
	// C[2][0] = +(b*f - c*e), C[2][1] = -(a*f - c*d), C[2][2] = +(a*e - b*d)

	a, b, c := m.M[0][0], m.M[0][1], m.M[0][2]
	d, e, f := m.M[1][0], m.M[1][1], m.M[1][2]
	g, h, i := m.M[2][0], m.M[2][1], m.M[2][2]

	// Compute cofactors
	cofactor := [3][3]*BigFloat{}

	// Row 0
	cofactor[0][0] = new(BigFloat).SetPrec(prec).Sub(new(BigFloat).SetPrec(prec).Mul(e, i), new(BigFloat).SetPrec(prec).Mul(f, h))
	cofactor[0][1] = new(BigFloat).SetPrec(prec).Neg(new(BigFloat).SetPrec(prec).Sub(new(BigFloat).SetPrec(prec).Mul(d, i), new(BigFloat).SetPrec(prec).Mul(f, g)))
	cofactor[0][2] = new(BigFloat).SetPrec(prec).Sub(new(BigFloat).SetPrec(prec).Mul(d, h), new(BigFloat).SetPrec(prec).Mul(e, g))

	// Row 1
	cofactor[1][0] = new(BigFloat).SetPrec(prec).Neg(new(BigFloat).SetPrec(prec).Sub(new(BigFloat).SetPrec(prec).Mul(b, i), new(BigFloat).SetPrec(prec).Mul(c, h)))
	cofactor[1][1] = new(BigFloat).SetPrec(prec).Sub(new(BigFloat).SetPrec(prec).Mul(a, i), new(BigFloat).SetPrec(prec).Mul(c, g))
	cofactor[1][2] = new(BigFloat).SetPrec(prec).Neg(new(BigFloat).SetPrec(prec).Sub(new(BigFloat).SetPrec(prec).Mul(a, h), new(BigFloat).SetPrec(prec).Mul(b, g)))

	// Row 2
	cofactor[2][0] = new(BigFloat).SetPrec(prec).Sub(new(BigFloat).SetPrec(prec).Mul(b, f), new(BigFloat).SetPrec(prec).Mul(c, e))
	cofactor[2][1] = new(BigFloat).SetPrec(prec).Neg(new(BigFloat).SetPrec(prec).Sub(new(BigFloat).SetPrec(prec).Mul(a, f), new(BigFloat).SetPrec(prec).Mul(c, d)))
	cofactor[2][2] = new(BigFloat).SetPrec(prec).Sub(new(BigFloat).SetPrec(prec).Mul(a, e), new(BigFloat).SetPrec(prec).Mul(b, d))

	// Adjugate is transpose of cofactor
	adjugate := &BigMatrix3x3{
		M: [3][3]*BigFloat{
			{cofactor[0][0], cofactor[1][0], cofactor[2][0]},
			{cofactor[0][1], cofactor[1][1], cofactor[2][1]},
			{cofactor[0][2], cofactor[1][2], cofactor[2][2]},
		},
	}

	// Multiply adjugate by 1/det
	one := NewBigFloat(1.0, prec)
	invDet := new(BigFloat).SetPrec(prec).Quo(one, det)

	result := &BigMatrix3x3{
		M: [3][3]*BigFloat{},
	}

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			result.M[i][j] = new(BigFloat).SetPrec(prec).Mul(adjugate.M[i][j], invDet)
		}
	}

	return result, nil
}
