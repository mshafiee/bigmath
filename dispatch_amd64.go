//go:build amd64

package bigmath

// initDispatcherImpl sets up AMD64-specific function pointers
func initDispatcherImpl(d *Dispatcher) {
	// AMD64 assembly implementations available
	if d.Features.HasAVX2 {
		// Use AVX2-optimized implementations when available
		d.BigVec3AddImpl = bigVec3AddAVX2
		d.BigVec3SubImpl = bigVec3SubGeneric // Use generic directly - assembly stub adds overhead
		d.BigVec3MulImpl = bigVec3MulAVX2
		d.BigVec3DotImpl = bigVec3DotAVX2
		d.BigMatMulImpl = bigMatMulAVX2
		// BigVec6 operations use generic implementations (no assembly benefit for arbitrary-precision)
		d.BigVec6AddImpl = bigVec6AddGeneric
		d.BigVec6SubImpl = bigVec6SubGeneric
		d.BigVec6NegateImpl = bigVec6NegateGeneric
		d.BigVec6MagnitudeImpl = bigVec6MagnitudeGeneric
		d.EvaluateChebyshevBigImpl = evaluateChebyshevBigAVX2
		d.EvaluateChebyshevDerivativeBigImpl = evaluateChebyshevDerivativeBigAVX2
		d.BigSinImpl = bigSinAVX2
		d.BigCosImpl = bigCosAVX2
		// Use assembly implementations for exp and log
		d.BigExpImpl = bigExpAsm
		d.BigLogImpl = bigLogAsm
		// Use optimized implementations
		d.BigTanImpl = bigTanGeneric // tan = sin/cos, already optimized
		d.BigAtanImpl = bigAtanOptimized
		d.BigAsinImpl = bigAsinOptimized
		d.BigAcosImpl = bigAcosOptimized
		d.BigAtan2Impl = bigAtan2Optimized
		d.BigPowImpl = bigPowGeneric
		d.BigSinhImpl = bigSinhGeneric
		d.BigCoshImpl = bigCoshGeneric
		d.BigTanhImpl = bigTanhGeneric
		d.BigAsinhImpl = bigAsinhGeneric
		d.BigAcoshImpl = bigAcoshGeneric
		d.BigAtanhImpl = bigAtanhGeneric
	} else {
		// Fallback to standard AMD64 assembly
		d.BigVec3AddImpl = bigVec3AddAMD64
		d.BigVec3SubImpl = bigVec3SubGeneric // Use generic directly - assembly stub adds overhead
		d.BigVec3MulImpl = bigVec3MulAMD64
		d.BigVec3DotImpl = bigVec3DotAMD64
		d.BigMatMulImpl = bigMatMulAMD64
		// BigVec6 operations use generic implementations (no assembly benefit for arbitrary-precision)
		d.BigVec6AddImpl = bigVec6AddGeneric
		d.BigVec6SubImpl = bigVec6SubGeneric
		d.BigVec6NegateImpl = bigVec6NegateGeneric
		d.BigVec6MagnitudeImpl = bigVec6MagnitudeGeneric
		d.EvaluateChebyshevBigImpl = evaluateChebyshevBigAMD64
		d.EvaluateChebyshevDerivativeBigImpl = evaluateChebyshevDerivativeBigAMD64
		d.BigSinImpl = bigSinAMD64
		d.BigCosImpl = bigCosAMD64
		// Use assembly implementations for exp and log
		d.BigExpImpl = bigExpAsm
		d.BigLogImpl = bigLogAsm
		// Use optimized implementations
		d.BigTanImpl = bigTanGeneric // tan = sin/cos, already optimized
		d.BigAtanImpl = bigAtanOptimized
		d.BigAsinImpl = bigAsinOptimized
		d.BigAcosImpl = bigAcosOptimized
		d.BigAtan2Impl = bigAtan2Optimized
		d.BigPowImpl = bigPowGeneric
		d.BigSinhImpl = bigSinhGeneric
		d.BigCoshImpl = bigCoshGeneric
		d.BigTanhImpl = bigTanhGeneric
		d.BigAsinhImpl = bigAsinhGeneric
		d.BigAcoshImpl = bigAcoshGeneric
		d.BigAtanhImpl = bigAtanhGeneric
	}
}

