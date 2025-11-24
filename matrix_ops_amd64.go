// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

//go:build amd64

package bigmath

// AMD64 wrapper functions for advanced matrix operations
// Now using optimized implementations with reduced allocations

func bigMatTransposeAsm(m *BigMatrix3x3, prec uint) *BigMatrix3x3 {
	return bigMatTransposeOptimized(m, prec)
}

func bigMatMulMatAsm(m1, m2 *BigMatrix3x3, prec uint) *BigMatrix3x3 {
	return bigMatMulMatOptimized(m1, m2, prec)
}

func bigMatDetAsm(m *BigMatrix3x3, prec uint) *BigFloat {
	return bigMatDetOptimized(m, prec)
}
