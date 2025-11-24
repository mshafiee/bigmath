// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

//go:build arm64

package bigmath

// ARM64 assembly implementations and stubs for vector/matrix operations

//nolint:unused // May be used in future dispatch implementations
func bigVec3AddARM64(v1, v2 *BigVec3, prec uint) *BigVec3 {
	return bigVec3AddAsmARM64(v1, v2, prec)
}

//nolint:unused // May be used in future dispatch implementations
func bigVec3SubARM64(v1, v2 *BigVec3, prec uint) *BigVec3 {
	return bigVec3SubAsmARM64(v1, v2, prec)
}

//nolint:unused // May be used in future dispatch implementations
func bigVec3MulARM64(v *BigVec3, scalar *BigFloat, prec uint) *BigVec3 {
	return bigVec3MulAsmARM64(v, scalar, prec)
}

//nolint:unused // May be used in future dispatch implementations
func bigVec3DotARM64(v1, v2 *BigVec3, prec uint) *BigFloat {
	return bigVec3DotAsmARM64(v1, v2, prec)
}

// Assembly wrappers replaced with Go implementations to avoid GC stackmap issues
// The real implementation is in the generic functions

//nolint:all // Used by ARM64 dispatch
func bigVec3AddAsmARM64(v1, v2 *BigVec3, prec uint) *BigVec3 {
	return bigVec3AddGeneric(v1, v2, prec)
}

//nolint:all // Used by ARM64 dispatch
func bigVec3SubAsmARM64(v1, v2 *BigVec3, prec uint) *BigVec3 {
	return bigVec3SubGeneric(v1, v2, prec)
}

//nolint:all // Used by ARM64 dispatch
func bigVec3MulAsmARM64(v *BigVec3, scalar *BigFloat, prec uint) *BigVec3 {
	return bigVec3MulGeneric(v, scalar, prec)
}

//nolint:all // Used by ARM64 dispatch
func bigVec3DotAsmARM64(v1, v2 *BigVec3, prec uint) *BigFloat {
	return bigVec3DotGeneric(v1, v2, prec)
}

// Matrix operations

//nolint:all // Used by ARM64 dispatch
func bigMatMulARM64(m *BigMatrix3x3, v *BigVec3, prec uint) *BigVec3 {
	return bigMatMulAsmARM64(m, v, prec)
}

//nolint:all // Used by ARM64 dispatch
func bigMatMulAsmARM64(m *BigMatrix3x3, v *BigVec3, prec uint) *BigVec3 {
	return bigMatMulGeneric(m, v, prec)
}
