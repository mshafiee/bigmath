// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

import (
	"testing"
)

// Comparison benchmarks: Generic (Pure Go) vs Optimized implementations

func BenchmarkBigFactorial_Generic(b *testing.B) {
	x := NewBigFloat(10.0, 256)
	for i := 0; i < b.N; i++ {
		bigFactorialGeneric(10, 256)
		_ = x
	}
}

func BenchmarkBigFactorial_Optimized(b *testing.B) {
	x := NewBigFloat(10.0, 256)
	for i := 0; i < b.N; i++ {
		bigFactorialOptimized(10, 256)
		_ = x
	}
}

func BenchmarkBigFloor_Generic(b *testing.B) {
	x := NewBigFloat(3.7, 256)
	for i := 0; i < b.N; i++ {
		bigFloorGeneric(x, 256)
	}
}

func BenchmarkBigFloor_Optimized(b *testing.B) {
	x := NewBigFloat(3.7, 256)
	for i := 0; i < b.N; i++ {
		bigFloorOptimized(x, 256)
	}
}

func BenchmarkBigCeil_Generic(b *testing.B) {
	x := NewBigFloat(3.2, 256)
	for i := 0; i < b.N; i++ {
		bigCeilGeneric(x, 256)
	}
}

func BenchmarkBigCeil_Optimized(b *testing.B) {
	x := NewBigFloat(3.2, 256)
	for i := 0; i < b.N; i++ {
		bigCeilOptimized(x, 256)
	}
}

func BenchmarkBigTrunc_Generic(b *testing.B) {
	x := NewBigFloat(3.9, 256)
	for i := 0; i < b.N; i++ {
		bigTruncGeneric(x, 256)
	}
}

func BenchmarkBigTrunc_Optimized(b *testing.B) {
	x := NewBigFloat(3.9, 256)
	for i := 0; i < b.N; i++ {
		bigTruncOptimized(x, 256)
	}
}

func BenchmarkBigCbrt_Generic(b *testing.B) {
	x := NewBigFloat(27.0, 256)
	for i := 0; i < b.N; i++ {
		bigCbrtGeneric(x, 256)
	}
}

func BenchmarkBigCbrt_Optimized(b *testing.B) {
	x := NewBigFloat(27.0, 256)
	for i := 0; i < b.N; i++ {
		bigCbrtOptimized(x, 256)
	}
}

func BenchmarkBigMod_Generic(b *testing.B) {
	x := NewBigFloat(10.5, 256)
	y := NewBigFloat(3.0, 256)
	for i := 0; i < b.N; i++ {
		bigModGeneric(x, y, 256)
	}
}

func BenchmarkBigMod_Optimized(b *testing.B) {
	x := NewBigFloat(10.5, 256)
	y := NewBigFloat(3.0, 256)
	for i := 0; i < b.N; i++ {
		bigModOptimized(x, y, 256)
	}
}

func BenchmarkBigRem_Generic(b *testing.B) {
	x := NewBigFloat(10.5, 256)
	y := NewBigFloat(3.0, 256)
	for i := 0; i < b.N; i++ {
		bigRemGeneric(x, y, 256)
	}
}

func BenchmarkBigRem_Optimized(b *testing.B) {
	x := NewBigFloat(10.5, 256)
	y := NewBigFloat(3.0, 256)
	for i := 0; i < b.N; i++ {
		bigRemOptimized(x, y, 256)
	}
}

func BenchmarkBigBinomial_Generic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bigBinomialGeneric(20, 5, 256)
	}
}

func BenchmarkBigBinomial_Optimized(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bigBinomialOptimized(20, 5, 256)
	}
}

func BenchmarkBigVec3Cross_Generic(b *testing.B) {
	v1 := &BigVec3{
		X: NewBigFloat(1.0, 256),
		Y: NewBigFloat(2.0, 256),
		Z: NewBigFloat(3.0, 256),
	}
	v2 := &BigVec3{
		X: NewBigFloat(4.0, 256),
		Y: NewBigFloat(5.0, 256),
		Z: NewBigFloat(6.0, 256),
	}
	for i := 0; i < b.N; i++ {
		bigVec3CrossGeneric(v1, v2, 256)
	}
}

func BenchmarkBigVec3Cross_Optimized(b *testing.B) {
	v1 := &BigVec3{
		X: NewBigFloat(1.0, 256),
		Y: NewBigFloat(2.0, 256),
		Z: NewBigFloat(3.0, 256),
	}
	v2 := &BigVec3{
		X: NewBigFloat(4.0, 256),
		Y: NewBigFloat(5.0, 256),
		Z: NewBigFloat(6.0, 256),
	}
	for i := 0; i < b.N; i++ {
		bigVec3CrossOptimized(v1, v2, 256)
	}
}

func BenchmarkBigVec3Normalize_Generic(b *testing.B) {
	v := &BigVec3{
		X: NewBigFloat(1.0, 256),
		Y: NewBigFloat(2.0, 256),
		Z: NewBigFloat(3.0, 256),
	}
	for i := 0; i < b.N; i++ {
		bigVec3NormalizeGeneric(v, 256)
	}
}

func BenchmarkBigVec3Normalize_Optimized(b *testing.B) {
	v := &BigVec3{
		X: NewBigFloat(1.0, 256),
		Y: NewBigFloat(2.0, 256),
		Z: NewBigFloat(3.0, 256),
	}
	for i := 0; i < b.N; i++ {
		bigVec3NormalizeOptimized(v, 256)
	}
}

func BenchmarkBigVec3Angle_Generic(b *testing.B) {
	v1 := &BigVec3{
		X: NewBigFloat(1.0, 256),
		Y: NewBigFloat(0.0, 256),
		Z: NewBigFloat(0.0, 256),
	}
	v2 := &BigVec3{
		X: NewBigFloat(0.0, 256),
		Y: NewBigFloat(1.0, 256),
		Z: NewBigFloat(0.0, 256),
	}
	for i := 0; i < b.N; i++ {
		bigVec3AngleGeneric(v1, v2, 256)
	}
}

func BenchmarkBigVec3Angle_Optimized(b *testing.B) {
	v1 := &BigVec3{
		X: NewBigFloat(1.0, 256),
		Y: NewBigFloat(0.0, 256),
		Z: NewBigFloat(0.0, 256),
	}
	v2 := &BigVec3{
		X: NewBigFloat(0.0, 256),
		Y: NewBigFloat(1.0, 256),
		Z: NewBigFloat(0.0, 256),
	}
	for i := 0; i < b.N; i++ {
		bigVec3AngleOptimized(v1, v2, 256)
	}
}

func BenchmarkBigVec3Project_Generic(b *testing.B) {
	v1 := &BigVec3{
		X: NewBigFloat(1.0, 256),
		Y: NewBigFloat(2.0, 256),
		Z: NewBigFloat(3.0, 256),
	}
	v2 := &BigVec3{
		X: NewBigFloat(1.0, 256),
		Y: NewBigFloat(0.0, 256),
		Z: NewBigFloat(0.0, 256),
	}
	for i := 0; i < b.N; i++ {
		bigVec3ProjectGeneric(v1, v2, 256)
	}
}

func BenchmarkBigVec3Project_Optimized(b *testing.B) {
	v1 := &BigVec3{
		X: NewBigFloat(1.0, 256),
		Y: NewBigFloat(2.0, 256),
		Z: NewBigFloat(3.0, 256),
	}
	v2 := &BigVec3{
		X: NewBigFloat(1.0, 256),
		Y: NewBigFloat(0.0, 256),
		Z: NewBigFloat(0.0, 256),
	}
	for i := 0; i < b.N; i++ {
		bigVec3ProjectOptimized(v1, v2, 256)
	}
}

func BenchmarkBigMatTranspose_Generic(b *testing.B) {
	m := &BigMatrix3x3{
		M: [3][3]*BigFloat{
			{NewBigFloat(1.0, 256), NewBigFloat(2.0, 256), NewBigFloat(3.0, 256)},
			{NewBigFloat(4.0, 256), NewBigFloat(5.0, 256), NewBigFloat(6.0, 256)},
			{NewBigFloat(7.0, 256), NewBigFloat(8.0, 256), NewBigFloat(9.0, 256)},
		},
	}
	for i := 0; i < b.N; i++ {
		bigMatTransposeGeneric(m, 256)
	}
}

func BenchmarkBigMatTranspose_Optimized(b *testing.B) {
	m := &BigMatrix3x3{
		M: [3][3]*BigFloat{
			{NewBigFloat(1.0, 256), NewBigFloat(2.0, 256), NewBigFloat(3.0, 256)},
			{NewBigFloat(4.0, 256), NewBigFloat(5.0, 256), NewBigFloat(6.0, 256)},
			{NewBigFloat(7.0, 256), NewBigFloat(8.0, 256), NewBigFloat(9.0, 256)},
		},
	}
	for i := 0; i < b.N; i++ {
		bigMatTransposeOptimized(m, 256)
	}
}

func BenchmarkBigMatMulMat_Generic(b *testing.B) {
	m1 := &BigMatrix3x3{
		M: [3][3]*BigFloat{
			{NewBigFloat(1.0, 256), NewBigFloat(2.0, 256), NewBigFloat(3.0, 256)},
			{NewBigFloat(4.0, 256), NewBigFloat(5.0, 256), NewBigFloat(6.0, 256)},
			{NewBigFloat(7.0, 256), NewBigFloat(8.0, 256), NewBigFloat(9.0, 256)},
		},
	}
	m2 := &BigMatrix3x3{
		M: [3][3]*BigFloat{
			{NewBigFloat(1.0, 256), NewBigFloat(0.0, 256), NewBigFloat(0.0, 256)},
			{NewBigFloat(0.0, 256), NewBigFloat(1.0, 256), NewBigFloat(0.0, 256)},
			{NewBigFloat(0.0, 256), NewBigFloat(0.0, 256), NewBigFloat(1.0, 256)},
		},
	}
	for i := 0; i < b.N; i++ {
		bigMatMulMatGeneric(m1, m2, 256)
	}
}

func BenchmarkBigMatMulMat_Optimized(b *testing.B) {
	m1 := &BigMatrix3x3{
		M: [3][3]*BigFloat{
			{NewBigFloat(1.0, 256), NewBigFloat(2.0, 256), NewBigFloat(3.0, 256)},
			{NewBigFloat(4.0, 256), NewBigFloat(5.0, 256), NewBigFloat(6.0, 256)},
			{NewBigFloat(7.0, 256), NewBigFloat(8.0, 256), NewBigFloat(9.0, 256)},
		},
	}
	m2 := &BigMatrix3x3{
		M: [3][3]*BigFloat{
			{NewBigFloat(1.0, 256), NewBigFloat(0.0, 256), NewBigFloat(0.0, 256)},
			{NewBigFloat(0.0, 256), NewBigFloat(1.0, 256), NewBigFloat(0.0, 256)},
			{NewBigFloat(0.0, 256), NewBigFloat(0.0, 256), NewBigFloat(1.0, 256)},
		},
	}
	for i := 0; i < b.N; i++ {
		bigMatMulMatOptimized(m1, m2, 256)
	}
}

func BenchmarkBigMatDet_Generic(b *testing.B) {
	m := &BigMatrix3x3{
		M: [3][3]*BigFloat{
			{NewBigFloat(1.0, 256), NewBigFloat(2.0, 256), NewBigFloat(3.0, 256)},
			{NewBigFloat(4.0, 256), NewBigFloat(5.0, 256), NewBigFloat(6.0, 256)},
			{NewBigFloat(7.0, 256), NewBigFloat(8.0, 256), NewBigFloat(9.0, 256)},
		},
	}
	for i := 0; i < b.N; i++ {
		bigMatDetGeneric(m, 256)
	}
}

func BenchmarkBigMatDet_Optimized(b *testing.B) {
	m := &BigMatrix3x3{
		M: [3][3]*BigFloat{
			{NewBigFloat(1.0, 256), NewBigFloat(2.0, 256), NewBigFloat(3.0, 256)},
			{NewBigFloat(4.0, 256), NewBigFloat(5.0, 256), NewBigFloat(6.0, 256)},
			{NewBigFloat(7.0, 256), NewBigFloat(8.0, 256), NewBigFloat(9.0, 256)},
		},
	}
	for i := 0; i < b.N; i++ {
		bigMatDetOptimized(m, 256)
	}
}

func BenchmarkBigErf_Generic(b *testing.B) {
	x := NewBigFloat(0.5, 256)
	for i := 0; i < b.N; i++ {
		bigErfGeneric(x, 256)
	}
}

func BenchmarkBigErf_Optimized(b *testing.B) {
	x := NewBigFloat(0.5, 256)
	for i := 0; i < b.N; i++ {
		bigErfOptimized(x, 256)
	}
}

func BenchmarkBigErfc_Generic(b *testing.B) {
	x := NewBigFloat(0.5, 256)
	for i := 0; i < b.N; i++ {
		bigErfcGeneric(x, 256)
	}
}

func BenchmarkBigErfc_Optimized(b *testing.B) {
	x := NewBigFloat(0.5, 256)
	for i := 0; i < b.N; i++ {
		bigErfcOptimized(x, 256)
	}
}

func BenchmarkBigGamma_Generic(b *testing.B) {
	x := NewBigFloat(5.0, 256)
	for i := 0; i < b.N; i++ {
		bigGammaGeneric(x, 256)
	}
}

func BenchmarkBigGamma_Optimized(b *testing.B) {
	x := NewBigFloat(5.0, 256)
	for i := 0; i < b.N; i++ {
		bigGammaOptimized(x, 256)
	}
}
