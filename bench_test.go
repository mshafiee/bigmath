// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

import (
	"math"
	"testing"
)

const benchPrec = 256

// BenchmarkBigSin benchmarks the BigSin function
func BenchmarkBigSin(b *testing.B) {
	x := NewBigFloat(math.Pi/4, benchPrec)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = BigSin(x, benchPrec)
	}
}

// BenchmarkBigCos benchmarks the BigCos function
func BenchmarkBigCos(b *testing.B) {
	x := NewBigFloat(math.Pi/4, benchPrec)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = BigCos(x, benchPrec)
	}
}

// BenchmarkBigExp benchmarks the BigExp function
// Note: Skipped due to assembly/Go interop issues with stack maps
// func BenchmarkBigExp(b *testing.B) {
// 	x := NewBigFloat(1.0, benchPrec)
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		_ = BigExp(x, benchPrec)
// 	}
// }

// BenchmarkBigLog benchmarks the BigLog function
// Note: Skipped due to assembly/Go interop issues with stack maps
// func BenchmarkBigLog(b *testing.B) {
// 	x := NewBigFloat(2.0, benchPrec)
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		_ = BigLog(x, benchPrec)
// 	}
// }

// BenchmarkBigPow benchmarks the BigPow function
func BenchmarkBigPow(b *testing.B) {
	x := NewBigFloat(2.0, benchPrec)
	y := NewBigFloat(3.0, benchPrec)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = BigPow(x, y, benchPrec)
	}
}

// BenchmarkBigSqrt benchmarks the BigSqrt function
func BenchmarkBigSqrt(b *testing.B) {
	x := NewBigFloat(2.0, benchPrec)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = BigSqrt(x, benchPrec)
	}
}

// BenchmarkBigVec3Add benchmarks vector addition
func BenchmarkBigVec3Add(b *testing.B) {
	v1 := NewBigVec3(1.0, 2.0, 3.0, benchPrec)
	v2 := NewBigVec3(4.0, 5.0, 6.0, benchPrec)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = BigVec3Add(v1, v2, benchPrec)
	}
}

// BenchmarkBigVec3Dot benchmarks vector dot product
func BenchmarkBigVec3Dot(b *testing.B) {
	v1 := NewBigVec3(1.0, 2.0, 3.0, benchPrec)
	v2 := NewBigVec3(4.0, 5.0, 6.0, benchPrec)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = BigVec3Dot(v1, v2, benchPrec)
	}
}

// BenchmarkBigVec3Magnitude benchmarks vector magnitude
func BenchmarkBigVec3Magnitude(b *testing.B) {
	v := NewBigVec3(3.0, 4.0, 5.0, benchPrec)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = BigVec3Magnitude(v, benchPrec)
	}
}

// BenchmarkBigFloatAdd benchmarks BigFloat addition
func BenchmarkBigFloatAdd(b *testing.B) {
	x := NewBigFloat(3.14159, benchPrec)
	y := NewBigFloat(2.71828, benchPrec)
	result := new(BigFloat).SetPrec(benchPrec)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result.Add(x, y)
	}
}

// BenchmarkBigFloatMul benchmarks BigFloat multiplication
func BenchmarkBigFloatMul(b *testing.B) {
	x := NewBigFloat(3.14159, benchPrec)
	y := NewBigFloat(2.71828, benchPrec)
	result := new(BigFloat).SetPrec(benchPrec)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result.Mul(x, y)
	}
}

// BenchmarkEvaluateChebyshevBig benchmarks Chebyshev polynomial evaluation
func BenchmarkEvaluateChebyshevBig(b *testing.B) {
	t := NewBigFloat(0.5, benchPrec)
	coeffs := make([]*BigFloat, 10)
	for i := range coeffs {
		coeffs[i] = NewBigFloat(float64(i+1)*0.1, benchPrec)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = EvaluateChebyshevBig(t, coeffs, len(coeffs), benchPrec)
	}
}

