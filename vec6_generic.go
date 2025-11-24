// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

// Generic pure-Go implementations for BigVec6 operations
// These serve as reference implementations and fallbacks

// bigVec6AddGeneric adds two BigVec6 vectors (pure-Go)
func bigVec6AddGeneric(v1, v2 *BigVec6, prec uint) *BigVec6 {
	if prec == 0 {
		prec = DefaultPrecision
	}
	
	result := &BigVec6{
		X:  NewBigFloat(0, prec),
		Y:  NewBigFloat(0, prec),
		Z:  NewBigFloat(0, prec),
		VX: NewBigFloat(0, prec),
		VY: NewBigFloat(0, prec),
		VZ: NewBigFloat(0, prec),
	}
	
	result.X.Add(v1.X, v2.X)
	result.Y.Add(v1.Y, v2.Y)
	result.Z.Add(v1.Z, v2.Z)
	result.VX.Add(v1.VX, v2.VX)
	result.VY.Add(v1.VY, v2.VY)
	result.VZ.Add(v1.VZ, v2.VZ)
	
	return result
}

// bigVec6SubGeneric subtracts two BigVec6 vectors (pure-Go)
func bigVec6SubGeneric(v1, v2 *BigVec6, prec uint) *BigVec6 {
	if prec == 0 {
		prec = DefaultPrecision
	}
	
	result := &BigVec6{
		X:  NewBigFloat(0, prec),
		Y:  NewBigFloat(0, prec),
		Z:  NewBigFloat(0, prec),
		VX: NewBigFloat(0, prec),
		VY: NewBigFloat(0, prec),
		VZ: NewBigFloat(0, prec),
	}
	
	result.X.Sub(v1.X, v2.X)
	result.Y.Sub(v1.Y, v2.Y)
	result.Z.Sub(v1.Z, v2.Z)
	result.VX.Sub(v1.VX, v2.VX)
	result.VY.Sub(v1.VY, v2.VY)
	result.VZ.Sub(v1.VZ, v2.VZ)
	
	return result
}

// bigVec6NegateGeneric negates all components of a BigVec6 (pure-Go)
func bigVec6NegateGeneric(v *BigVec6, prec uint) *BigVec6 {
	if prec == 0 {
		prec = DefaultPrecision
	}
	
	result := &BigVec6{
		X:  NewBigFloat(0, prec),
		Y:  NewBigFloat(0, prec),
		Z:  NewBigFloat(0, prec),
		VX: NewBigFloat(0, prec),
		VY: NewBigFloat(0, prec),
		VZ: NewBigFloat(0, prec),
	}
	
	result.X.Neg(v.X)
	result.Y.Neg(v.Y)
	result.Z.Neg(v.Z)
	result.VX.Neg(v.VX)
	result.VY.Neg(v.VY)
	result.VZ.Neg(v.VZ)
	
	return result
}

// bigVec6MagnitudeGeneric computes the magnitude of the position component (pure-Go)
func bigVec6MagnitudeGeneric(v *BigVec6, prec uint) *BigFloat {
	if prec == 0 {
		prec = DefaultPrecision
	}
	
	// |r| = sqrt(x² + y² + z²)
	pos := &BigVec3{X: v.X, Y: v.Y, Z: v.Z}
	return BigVec3Magnitude(pos, prec)
}

