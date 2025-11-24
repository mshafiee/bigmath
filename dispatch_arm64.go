// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

//go:build arm64

package bigmath

// initDispatcherImpl sets up ARM64-specific function pointers
func initDispatcherImpl(d *Dispatcher) {
	// ARM64 assembly implementations available - use them when stable
	// For now, use assembly for basic functions and generic for complex ones
	d.BigVec3AddImpl = bigVec3AddGeneric
	d.BigVec3SubImpl = bigVec3SubGeneric // Already using generic directly
	d.BigVec3MulImpl = bigVec3MulGeneric
	d.BigVec3DotImpl = bigVec3DotGeneric
	d.BigMatMulImpl = bigMatMulGeneric
	// BigVec6 operations use generic implementations
	d.BigVec6AddImpl = bigVec6AddGeneric
	d.BigVec6SubImpl = bigVec6SubGeneric
	d.BigVec6NegateImpl = bigVec6NegateGeneric
	d.BigVec6MagnitudeImpl = bigVec6MagnitudeGeneric
	d.EvaluateChebyshevBigImpl = evaluateChebyshevBigARM64
	d.EvaluateChebyshevDerivativeBigImpl = evaluateChebyshevDerivativeBigARM64
	// Use assembly implementations for sin, cos, exp, and log
	// Note: ARM64 assembly implementations exist but have calling convention issues
	// Temporarily using generic until ARM64 assembly calling convention is fixed
	// TODO: Fix ARM64 assembly function calling convention - see exp_asm_arm64.s, log_asm_arm64.s, trig_arm64.s
	d.BigSinImpl = bigSinGeneric
	d.BigCosImpl = bigCosGeneric
	d.BigExpImpl = bigExpGeneric
	d.BigLogImpl = bigLogGeneric
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

	// Special functions
	d.BigGammaImpl = bigGammaAsm
	d.BigErfImpl = bigErfAsm
	d.BigErfcImpl = bigErfcAsm
	d.BigBesselJImpl = bigBesselJAsm
	d.BigBesselYImpl = bigBesselYAsm

	// Root functions
	d.BigCbrtImpl = bigCbrtAsm
	d.BigRootImpl = bigRootAsm

	// Basic operations
	d.BigFloorImpl = bigFloorAsm
	d.BigCeilImpl = bigCeilAsm
	d.BigTruncImpl = bigTruncAsm
	d.BigModImpl = bigModAsm
	d.BigRemImpl = bigRemAsm

	// Combinatorics
	d.BigFactorialImpl = bigFactorialAsm
	d.BigBinomialImpl = bigBinomialAsm

	// Advanced vector operations
	d.BigVec3CrossImpl = bigVec3CrossAsm
	d.BigVec3NormalizeImpl = bigVec3NormalizeAsm
	d.BigVec3AngleImpl = bigVec3AngleAsm
	d.BigVec3ProjectImpl = bigVec3ProjectAsm

	// Advanced matrix operations
	d.BigMatTransposeImpl = bigMatTransposeAsm
	d.BigMatMulMatImpl = bigMatMulMatAsm
	d.BigMatDetImpl = bigMatDetAsm
	d.BigMatInverseImpl = bigMatInverseGeneric // No asm for error-returning function yet
}
