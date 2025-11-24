// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

//go:build amd64

package bigmath

// AMD64 wrapper functions for advanced vector operations
// Now using optimized implementations with reduced allocations

func bigVec3CrossAsm(v1, v2 *BigVec3, prec uint) *BigVec3 {
	return bigVec3CrossOptimized(v1, v2, prec)
}

func bigVec3NormalizeAsm(v *BigVec3, prec uint) *BigVec3 {
	return bigVec3NormalizeOptimized(v, prec)
}

func bigVec3AngleAsm(v1, v2 *BigVec3, prec uint) *BigFloat {
	return bigVec3AngleOptimized(v1, v2, prec)
}

func bigVec3ProjectAsm(v1, v2 *BigVec3, prec uint) *BigVec3 {
	return bigVec3ProjectOptimized(v1, v2, prec)
}
