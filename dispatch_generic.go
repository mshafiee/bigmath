// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

//go:build !amd64 && !arm64

package bigmath

// initDispatcherImpl sets up generic (pure-Go) function pointers for non-AMD64/ARM64 platforms
func initDispatcherImpl(d *Dispatcher) {
	// Use generic pure-Go implementations as fallback
	d.BigVec3AddImpl = bigVec3AddGeneric
	d.BigVec3SubImpl = bigVec3SubGeneric
	d.BigVec3MulImpl = bigVec3MulGeneric
	d.BigVec3DotImpl = bigVec3DotGeneric
	d.BigMatMulImpl = bigMatMulGeneric
	// BigVec6 operations
	d.BigVec6AddImpl = bigVec6AddGeneric
	d.BigVec6SubImpl = bigVec6SubGeneric
	d.BigVec6NegateImpl = bigVec6NegateGeneric
	d.BigVec6MagnitudeImpl = bigVec6MagnitudeGeneric
	d.EvaluateChebyshevBigImpl = evaluateChebyshevBigGeneric
	d.EvaluateChebyshevDerivativeBigImpl = evaluateChebyshevDerivativeBigGeneric
	d.BigSinImpl = bigSinOptimized
	d.BigCosImpl = bigCosOptimized
	d.BigTanImpl = bigTanGeneric // tan = sin/cos, already optimized
	d.BigAtanImpl = bigAtanOptimized
	d.BigAsinImpl = bigAsinOptimized
	d.BigAcosImpl = bigAcosOptimized
	d.BigAtan2Impl = bigAtan2Optimized
	d.BigExpImpl = bigExpGeneric
	d.BigLogImpl = bigLogGeneric
	d.BigPowImpl = bigPowGeneric
	d.BigSinhImpl = bigSinhGeneric
	d.BigCoshImpl = bigCoshGeneric
	d.BigTanhImpl = bigTanhGeneric
	d.BigAsinhImpl = bigAsinhGeneric
	d.BigAcoshImpl = bigAcoshGeneric
	d.BigAtanhImpl = bigAtanhGeneric
}

