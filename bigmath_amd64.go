// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

//go:build amd64

package bigmath

// AMD64 assembly implementations and stubs for vector/matrix operations

// bigVec3AddAMD64 is the basic AMD64 implementation (no AVX)
func bigVec3AddAMD64(v1, v2 *BigVec3, prec uint) *BigVec3 {
	// For now, delegate to generic until assembly is implemented
	return bigVec3AddGeneric(v1, v2, prec)
}

// bigVec3AddAVX2 is the AVX2-optimized implementation
func bigVec3AddAVX2(v1, v2 *BigVec3, prec uint) *BigVec3 {
	// Actual implementation will be in assembly
	return bigVec3AddAsm(v1, v2, prec)
}

//nolint:unused // May be used in dispatch system
func bigVec3SubAMD64(v1, v2 *BigVec3, prec uint) *BigVec3 {
	return bigVec3SubGeneric(v1, v2, prec)
}

//nolint:unused // May be used in dispatch system
func bigVec3SubAVX2(v1, v2 *BigVec3, prec uint) *BigVec3 {
	return bigVec3SubAsm(v1, v2, prec)
}

// bigVec3MulAMD64 is the basic AMD64 implementation
func bigVec3MulAMD64(v *BigVec3, scalar *BigFloat, prec uint) *BigVec3 {
	return bigVec3MulGeneric(v, scalar, prec)
}

// bigVec3MulAVX2 is the AVX2-optimized implementation
func bigVec3MulAVX2(v *BigVec3, scalar *BigFloat, prec uint) *BigVec3 {
	return bigVec3MulAsm(v, scalar, prec)
}

// bigVec3DotAMD64 is the basic AMD64 implementation
func bigVec3DotAMD64(v1, v2 *BigVec3, prec uint) *BigFloat {
	return bigVec3DotGeneric(v1, v2, prec)
}

// bigVec3DotAVX2 is the AVX2-optimized implementation
func bigVec3DotAVX2(v1, v2 *BigVec3, prec uint) *BigFloat {
	return bigVec3DotAsm(v1, v2, prec)
}

// Assembly function declarations
// These will be implemented in bigmath_amd64.s

//go:noescape
func bigVec3AddAsm(v1, v2 *BigVec3, prec uint) *BigVec3

//go:noescape
//nolint:unused // Declared for assembly implementation
func bigVec3SubAsm(v1, v2 *BigVec3, prec uint) *BigVec3

//go:noescape
func bigVec3MulAsm(v *BigVec3, scalar *BigFloat, prec uint) *BigVec3

//go:noescape
func bigVec3DotAsm(v1, v2 *BigVec3, prec uint) *BigFloat

// Matrix operations

// bigMatMulAMD64 is the basic AMD64 implementation
func bigMatMulAMD64(m *BigMatrix3x3, v *BigVec3, prec uint) *BigVec3 {
	return bigMatMulGeneric(m, v, prec)
}

// bigMatMulAVX2 is the AVX2-optimized implementation
func bigMatMulAVX2(m *BigMatrix3x3, v *BigVec3, prec uint) *BigVec3 {
	return bigMatMulAsm(m, v, prec)
}

//go:noescape
func bigMatMulAsm(m *BigMatrix3x3, v *BigVec3, prec uint) *BigVec3
