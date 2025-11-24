//go:build arm64

package bigmath

// ARM64 assembly implementations for Chebyshev polynomial evaluation

// evaluateChebyshevBigARM64 is the ARM64/NEON implementation
// Using generic directly - optimized version has allocation overhead
func evaluateChebyshevBigARM64(t *BigFloat, c []*BigFloat, neval int, prec uint) *BigFloat {
	return evaluateChebyshevBigGeneric(t, c, neval, prec)
}

// evaluateChebyshevDerivativeBigARM64 is the ARM64/NEON implementation (optimized)
func evaluateChebyshevDerivativeBigARM64(t *BigFloat, c []*BigFloat, neval int, prec uint) *BigFloat {
	return evaluateChebyshevDerivativeBigOptimized(t, c, neval, prec)
}

//go:noescape
func evaluateChebyshevBigAsmARM64(t *BigFloat, c []*BigFloat, neval int, prec uint) *BigFloat

//go:noescape
func evaluateChebyshevDerivativeBigAsmARM64(t *BigFloat, c []*BigFloat, neval int, prec uint) *BigFloat

