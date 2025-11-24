// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

// BigVec3Cross computes the cross product of two 3D vectors: v1 × v2
// Result = (v1.Y*v2.Z - v1.Z*v2.Y, v1.Z*v2.X - v1.X*v2.Z, v1.X*v2.Y - v1.Y*v2.X)
func BigVec3Cross(v1, v2 *BigVec3, prec uint) *BigVec3 {
	return getDispatcher().BigVec3CrossImpl(v1, v2, prec)
}

// BigVec3Normalize normalizes a 3D vector to unit length
// Returns a unit vector in the same direction, or zero vector if input is zero
func BigVec3Normalize(v *BigVec3, prec uint) *BigVec3 {
	return getDispatcher().BigVec3NormalizeImpl(v, prec)
}

// BigVec3Angle computes the angle between two 3D vectors in radians
// Returns angle in range [0, π] using: angle = arccos((v1·v2) / (|v1|*|v2|))
func BigVec3Angle(v1, v2 *BigVec3, prec uint) *BigFloat {
	return getDispatcher().BigVec3AngleImpl(v1, v2, prec)
}

// BigVec3Project projects vector v1 onto vector v2
// Returns the projection: ((v1·v2) / |v2|^2) * v2
func BigVec3Project(v1, v2 *BigVec3, prec uint) *BigVec3 {
	return getDispatcher().BigVec3ProjectImpl(v1, v2, prec)
}
