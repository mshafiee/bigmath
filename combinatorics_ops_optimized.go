// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

//go:build amd64 || arm64

package bigmath

// Optimized combinatorics functions

// bigFactorialOptimized implements optimized factorial
// Optimization: Fast path for small values with lookup table
func bigFactorialOptimized(n int64, prec uint) *BigFloat {
	// Fast path for small values - precomputed lookup
	if n >= 0 && n <= 20 {
		factorialTable := []float64{
			1, 1, 2, 6, 24, 120, 720, 5040, 40320, 362880,
			3628800, 39916800, 479001600, 6227020800, 87178291200,
			1307674368000, 20922789888000, 355687428096000,
			6402373705728000, 121645100408832000, 2432902008176640000,
		}
		return NewBigFloat(factorialTable[n], prec)
	}

	// For larger values, use original implementation
	return bigFactorialGeneric(n, prec)
}

// bigBinomialOptimized implements optimized binomial coefficient
// Optimization: Early returns, optimized computation order
func bigBinomialOptimized(n, k int64, prec uint) *BigFloat {
	// Fast paths
	if k > n || k < 0 || n < 0 {
		return NewBigFloat(0.0, prec)
	}
	if k == 0 || k == n {
		return NewBigFloat(1.0, prec)
	}
	if k == 1 || k == n-1 {
		return NewBigFloat(float64(n), prec)
	}

	// Use symmetry: C(n,k) = C(n, n-k)
	if k > n-k {
		k = n - k
	}

	// For small k, use direct computation
	if k <= 5 {
		return bigBinomialDirectOptimized(n, k, prec)
	}

	// For larger values, use original implementation
	return bigBinomialGeneric(n, k, prec)
}

// bigBinomialDirectOptimized computes C(n,k) directly for small k
// Optimization: Accumulate multiplications, reduce divisions
func bigBinomialDirectOptimized(n, k int64, prec uint) *BigFloat {
	result := NewBigFloat(1.0, prec)
	temp := new(BigFloat).SetPrec(prec)

	// Compute n!/(n-k)! / k!
	// = (n * (n-1) * ... * (n-k+1)) / (k * (k-1) * ... * 1)
	for i := int64(0); i < k; i++ {
		// Multiply by (n-i)
		temp.SetFloat64(float64(n - i))
		result.Mul(result, temp)

		// Divide by (i+1)
		temp.SetFloat64(float64(i + 1))
		result.Quo(result, temp)
	}

	return result
}
