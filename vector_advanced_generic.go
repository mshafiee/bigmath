// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

import (
	"math"
)

// Generic implementations for advanced vector operations (used as fallback)

// bigVec3CrossGeneric computes the cross product using pure Go implementation
func bigVec3CrossGeneric(v1, v2 *BigVec3, prec uint) *BigVec3 {
	if prec == 0 {
		prec = v1.X.Prec()
	}

	return &BigVec3{
		X: new(BigFloat).SetPrec(prec).Sub(
			new(BigFloat).SetPrec(prec).Mul(v1.Y, v2.Z),
			new(BigFloat).SetPrec(prec).Mul(v1.Z, v2.Y),
		),
		Y: new(BigFloat).SetPrec(prec).Sub(
			new(BigFloat).SetPrec(prec).Mul(v1.Z, v2.X),
			new(BigFloat).SetPrec(prec).Mul(v1.X, v2.Z),
		),
		Z: new(BigFloat).SetPrec(prec).Sub(
			new(BigFloat).SetPrec(prec).Mul(v1.X, v2.Y),
			new(BigFloat).SetPrec(prec).Mul(v1.Y, v2.X),
		),
	}
}

// bigVec3NormalizeGeneric normalizes a 3D vector using pure Go implementation
func bigVec3NormalizeGeneric(v *BigVec3, prec uint) *BigVec3 {
	if prec == 0 {
		prec = v.X.Prec()
	}

	magnitude := BigVec3Magnitude(v, prec)

	// Check if magnitude is zero
	zero := NewBigFloat(0.0, prec)
	if magnitude.Cmp(zero) == 0 {
		// Return zero vector
		return &BigVec3{
			X: NewBigFloat(0.0, prec),
			Y: NewBigFloat(0.0, prec),
			Z: NewBigFloat(0.0, prec),
		}
	}

	// Normalize each component
	return &BigVec3{
		X: new(BigFloat).SetPrec(prec).Quo(v.X, magnitude),
		Y: new(BigFloat).SetPrec(prec).Quo(v.Y, magnitude),
		Z: new(BigFloat).SetPrec(prec).Quo(v.Z, magnitude),
	}
}

// bigVec3AngleGeneric computes the angle between two 3D vectors using pure Go implementation
func bigVec3AngleGeneric(v1, v2 *BigVec3, prec uint) *BigFloat {
	if prec == 0 {
		prec = v1.X.Prec()
	}

	// Compute dot product
	dot := BigVec3Dot(v1, v2, prec)

	// Compute magnitudes
	mag1 := BigVec3Magnitude(v1, prec)
	mag2 := BigVec3Magnitude(v2, prec)

	// Check for zero vectors
	zero := NewBigFloat(0.0, prec)
	if mag1.Cmp(zero) == 0 || mag2.Cmp(zero) == 0 {
		return NewBigFloat(math.NaN(), prec)
	}

	// Compute cos(angle) = (v1·v2) / (|v1|*|v2|)
	magProduct := new(BigFloat).SetPrec(prec).Mul(mag1, mag2)
	cosAngle := new(BigFloat).SetPrec(prec).Quo(dot, magProduct)

	// Clamp cosAngle to [-1, 1] to avoid numerical errors
	one := NewBigFloat(1.0, prec)
	negOne := NewBigFloat(-1.0, prec)
	if cosAngle.Cmp(one) > 0 {
		cosAngle.Set(one)
	}
	if cosAngle.Cmp(negOne) < 0 {
		cosAngle.Set(negOne)
	}

	// Compute angle = arccos(cosAngle)
	return BigAcos(cosAngle, prec)
}

// bigVec3ProjectGeneric projects vector v1 onto vector v2 using pure Go implementation
func bigVec3ProjectGeneric(v1, v2 *BigVec3, prec uint) *BigVec3 {
	if prec == 0 {
		prec = v1.X.Prec()
	}

	// Compute dot product
	dot := BigVec3Dot(v1, v2, prec)

	// Compute |v2|^2
	mag2Sq := BigVec3Dot(v2, v2, prec)

	// Check if v2 is zero vector
	zero := NewBigFloat(0.0, prec)
	if mag2Sq.Cmp(zero) == 0 {
		// Return zero vector
		return &BigVec3{
			X: NewBigFloat(0.0, prec),
			Y: NewBigFloat(0.0, prec),
			Z: NewBigFloat(0.0, prec),
		}
	}

	// Compute scalar: (v1·v2) / |v2|^2
	scalar := new(BigFloat).SetPrec(prec).Quo(dot, mag2Sq)

	// Multiply v2 by scalar
	return BigVec3Mul(v2, scalar, prec)
}
