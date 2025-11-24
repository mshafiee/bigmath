// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

// Generic pure-Go implementations of vector and matrix operations
// These serve as reference implementations and fallbacks

// bigVec3AddGeneric adds two BigVec3 vectors: result = v1 + v2 (pure-Go)
func bigVec3AddGeneric(v1, v2 *BigVec3, prec uint) *BigVec3 {
	if prec == 0 {
		prec = v1.X.Prec()
	}
	return &BigVec3{
		X: new(BigFloat).SetPrec(prec).Add(v1.X, v2.X),
		Y: new(BigFloat).SetPrec(prec).Add(v1.Y, v2.Y),
		Z: new(BigFloat).SetPrec(prec).Add(v1.Z, v2.Z),
	}
}

// bigVec3SubGeneric subtracts two BigVec3 vectors: result = v1 - v2 (pure-Go)
func bigVec3SubGeneric(v1, v2 *BigVec3, prec uint) *BigVec3 {
	if prec == 0 {
		prec = v1.X.Prec()
	}
	return &BigVec3{
		X: new(BigFloat).SetPrec(prec).Sub(v1.X, v2.X),
		Y: new(BigFloat).SetPrec(prec).Sub(v1.Y, v2.Y),
		Z: new(BigFloat).SetPrec(prec).Sub(v1.Z, v2.Z),
	}
}

// bigVec3MulGeneric multiplies a BigVec3 by a scalar: result = v * scalar (pure-Go)
func bigVec3MulGeneric(v *BigVec3, scalar *BigFloat, prec uint) *BigVec3 {
	if prec == 0 {
		prec = v.X.Prec()
	}
	return &BigVec3{
		X: new(BigFloat).SetPrec(prec).Mul(v.X, scalar),
		Y: new(BigFloat).SetPrec(prec).Mul(v.Y, scalar),
		Z: new(BigFloat).SetPrec(prec).Mul(v.Z, scalar),
	}
}

// bigVec3DotGeneric computes the dot product of two BigVec3 vectors (pure-Go)
func bigVec3DotGeneric(v1, v2 *BigVec3, prec uint) *BigFloat {
	if prec == 0 {
		prec = v1.X.Prec()
	}

	// x1*x2 + y1*y2 + z1*z2
	result := new(BigFloat).SetPrec(prec)
	temp := new(BigFloat).SetPrec(prec)

	// x1 * x2
	result.Mul(v1.X, v2.X)

	// + y1 * y2
	temp.Mul(v1.Y, v2.Y)
	result.Add(result, temp)

	// + z1 * z2
	temp.Mul(v1.Z, v2.Z)
	result.Add(result, temp)

	return result
}

// bigMatMulGeneric multiplies a matrix by a vector: result = M * v (pure-Go)
func bigMatMulGeneric(m *BigMatrix3x3, v *BigVec3, prec uint) *BigVec3 {
	if prec == 0 {
		prec = v.X.Prec()
	}

	result := &BigVec3{
		X: new(BigFloat).SetPrec(prec),
		Y: new(BigFloat).SetPrec(prec),
		Z: new(BigFloat).SetPrec(prec),
	}

	temp := new(BigFloat).SetPrec(prec)

	// X = m[0][0]*v.X + m[0][1]*v.Y + m[0][2]*v.Z
	result.X.Mul(m.M[0][0], v.X)
	temp.Mul(m.M[0][1], v.Y)
	result.X.Add(result.X, temp)
	temp.Mul(m.M[0][2], v.Z)
	result.X.Add(result.X, temp)

	// Y = m[1][0]*v.X + m[1][1]*v.Y + m[1][2]*v.Z
	result.Y.Mul(m.M[1][0], v.X)
	temp.Mul(m.M[1][1], v.Y)
	result.Y.Add(result.Y, temp)
	temp.Mul(m.M[1][2], v.Z)
	result.Y.Add(result.Y, temp)

	// Z = m[2][0]*v.X + m[2][1]*v.Y + m[2][2]*v.Z
	result.Z.Mul(m.M[2][0], v.X)
	temp.Mul(m.M[2][1], v.Y)
	result.Z.Add(result.Z, temp)
	temp.Mul(m.M[2][2], v.Z)
	result.Z.Add(result.Z, temp)

	return result
}
