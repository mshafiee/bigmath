// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

//go:build amd64 || arm64

package bigmath

// Optimized matrix operations with reduced allocations and inlined operations

// bigMatTransposeOptimized implements an optimized matrix transpose
// Optimization: Reuse BigFloat objects, reduce allocations
func bigMatTransposeOptimized(m *BigMatrix3x3, prec uint) *BigMatrix3x3 {
	if prec == 0 {
		prec = m.M[0][0].Prec()
	}

	// Preallocate result matrix
	result := &BigMatrix3x3{M: [3][3]*BigFloat{}}

	// Transpose with direct assignment (cache-friendly access pattern)
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			result.M[i][j] = new(BigFloat).SetPrec(prec).Set(m.M[j][i])
		}
	}

	return result
}

// bigMatMulMatOptimized implements optimized 3x3 matrix multiplication
// Optimizations:
// - Reuse temporary BigFloat objects to reduce allocations
// - Unroll inner loop for better CPU pipelining
// - Use SetPrec once per element
func bigMatMulMatOptimized(m1, m2 *BigMatrix3x3, prec uint) *BigMatrix3x3 {
	if prec == 0 {
		prec = m1.M[0][0].Prec()
	}

	result := &BigMatrix3x3{M: [3][3]*BigFloat{}}

	// Preallocate temporary BigFloats for computation
	temp := new(BigFloat).SetPrec(prec)

	// Compute each element with unrolled multiplication
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			// Initialize sum
			sum := new(BigFloat).SetPrec(prec)

			// Unrolled loop: k=0
			temp.Mul(m1.M[i][0], m2.M[0][j])
			sum.Add(sum, temp)

			// k=1
			temp.Mul(m1.M[i][1], m2.M[1][j])
			sum.Add(sum, temp)

			// k=2
			temp.Mul(m1.M[i][2], m2.M[2][j])
			sum.Add(sum, temp)

			result.M[i][j] = sum
		}
	}

	return result
}

// bigMatDetOptimized implements optimized determinant calculation
// Optimizations:
// - Minimize allocations by reusing BigFloat objects
// - Compute sub-expressions in optimal order
// - Use FMA-like patterns where possible
func bigMatDetOptimized(m *BigMatrix3x3, prec uint) *BigFloat {
	if prec == 0 {
		prec = m.M[0][0].Prec()
	}

	// Extract elements (no allocation, just pointers)
	a, b, c := m.M[0][0], m.M[0][1], m.M[0][2]
	d, e, f := m.M[1][0], m.M[1][1], m.M[1][2]
	g, h, i := m.M[2][0], m.M[2][1], m.M[2][2]

	// Preallocate temporaries
	t1 := new(BigFloat).SetPrec(prec)
	t2 := new(BigFloat).SetPrec(prec)
	result := new(BigFloat).SetPrec(prec)

	// Compute: a(ei-fh) - b(di-fg) + c(dh-eg)
	// Unrolled and optimized computation order

	// First term: a * (e*i - f*h)
	t1.Mul(e, i)      // ei
	t2.Mul(f, h)      // fh
	t1.Sub(t1, t2)    // ei - fh
	result.Mul(a, t1) // a * (ei - fh)

	// Second term: -b * (d*i - f*g)
	t1.Mul(d, i)           // di
	t2.Mul(f, g)           // fg
	t1.Sub(t1, t2)         // di - fg
	t1.Mul(b, t1)          // b * (di - fg)
	result.Sub(result, t1) // subtract

	// Third term: c * (d*h - e*g)
	t1.Mul(d, h)           // dh
	t2.Mul(e, g)           // eg
	t1.Sub(t1, t2)         // dh - eg
	t1.Mul(c, t1)          // c * (dh - eg)
	result.Add(result, t1) // add

	return result
}
