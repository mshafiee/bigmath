// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

// BigMatTranspose returns the transpose of a 3x3 matrix
func BigMatTranspose(m *BigMatrix3x3, prec uint) *BigMatrix3x3 {
	return getDispatcher().BigMatTransposeImpl(m, prec)
}

// BigMatMulMat multiplies two 3x3 matrices: result = m1 * m2
func BigMatMulMat(m1, m2 *BigMatrix3x3, prec uint) *BigMatrix3x3 {
	return getDispatcher().BigMatMulMatImpl(m1, m2, prec)
}

// BigMatDet computes the determinant of a 3x3 matrix
// Uses the formula: det = a(ei-fh) - b(di-fg) + c(dh-eg)
// where matrix is:
//
//	[a b c]
//	[d e f]
//	[g h i]
func BigMatDet(m *BigMatrix3x3, prec uint) *BigFloat {
	return getDispatcher().BigMatDetImpl(m, prec)
}

// BigMatInverse computes the inverse of a 3x3 matrix
// Returns error if matrix is singular (determinant is zero)
// Uses adjugate matrix: M^-1 = (1/det(M)) * adj(M)
func BigMatInverse(m *BigMatrix3x3, prec uint) (*BigMatrix3x3, error) {
	return getDispatcher().BigMatInverseImpl(m, prec)
}
