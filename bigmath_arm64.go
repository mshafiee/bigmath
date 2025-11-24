//go:build arm64

package bigmath

// ARM64 assembly implementations and stubs for vector/matrix operations

// bigVec3AddARM64 is the ARM64/NEON implementation
func bigVec3AddARM64(v1, v2 *BigVec3, prec uint) *BigVec3 {
	return bigVec3AddAsmARM64(v1, v2, prec)
}

// bigVec3SubARM64 is the ARM64/NEON implementation
func bigVec3SubARM64(v1, v2 *BigVec3, prec uint) *BigVec3 {
	return bigVec3SubAsmARM64(v1, v2, prec)
}

// bigVec3MulARM64 is the ARM64/NEON implementation
func bigVec3MulARM64(v *BigVec3, scalar *BigFloat, prec uint) *BigVec3 {
	return bigVec3MulAsmARM64(v, scalar, prec)
}

// bigVec3DotARM64 is the ARM64/NEON implementation
func bigVec3DotARM64(v1, v2 *BigVec3, prec uint) *BigFloat {
	return bigVec3DotAsmARM64(v1, v2, prec)
}

// Assembly function declarations
// These will be implemented in bigmath_arm64.s

//go:noescape
func bigVec3AddAsmARM64(v1, v2 *BigVec3, prec uint) *BigVec3

//go:noescape
func bigVec3SubAsmARM64(v1, v2 *BigVec3, prec uint) *BigVec3

//go:noescape
func bigVec3MulAsmARM64(v *BigVec3, scalar *BigFloat, prec uint) *BigVec3

//go:noescape
func bigVec3DotAsmARM64(v1, v2 *BigVec3, prec uint) *BigFloat

// Matrix operations

// bigMatMulARM64 is the ARM64/NEON implementation
func bigMatMulARM64(m *BigMatrix3x3, v *BigVec3, prec uint) *BigVec3 {
	return bigMatMulAsmARM64(m, v, prec)
}

//go:noescape
func bigMatMulAsmARM64(m *BigMatrix3x3, v *BigVec3, prec uint) *BigVec3

