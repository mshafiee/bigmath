// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

// Generic implementations that assembly functions can call
// These serve as fallbacks and reference implementations

//nolint:unused // Used by assembly code in bigfloat_ops_amd64.s and bigfloat_ops_arm64.s
func bigfloatAddGeneric(a, b *BigFloat, prec uint) *BigFloat {
	if prec == 0 {
		prec = a.Prec()
	}
	result := new(BigFloat).SetPrec(prec)
	result.Add(a, b)
	return result
}

//nolint:unused // Used by assembly code in bigfloat_ops_amd64.s and bigfloat_ops_arm64.s
func bigfloatSubGeneric(a, b *BigFloat, prec uint) *BigFloat {
	if prec == 0 {
		prec = a.Prec()
	}
	result := new(BigFloat).SetPrec(prec)
	result.Sub(a, b)
	return result
}

//nolint:unused // Used by assembly code in bigfloat_ops_amd64.s and bigfloat_ops_arm64.s
func bigfloatMulGeneric(a, b *BigFloat, prec uint) *BigFloat {
	if prec == 0 {
		prec = a.Prec()
	}
	result := new(BigFloat).SetPrec(prec)
	result.Mul(a, b)
	return result
}

//nolint:unused // Used by assembly code in bigfloat_ops_amd64.s and bigfloat_ops_arm64.s
func bigfloatDivGeneric(a, b *BigFloat, prec uint) *BigFloat {
	if prec == 0 {
		prec = a.Prec()
	}
	result := new(BigFloat).SetPrec(prec)
	result.Quo(a, b)
	return result
}

//nolint:unused // Used by assembly code in bigfloat_ops_amd64.s and bigfloat_ops_arm64.s
func bigfloatSqrtGeneric(x *BigFloat, prec uint) *BigFloat {
	if prec == 0 {
		prec = x.Prec()
	}
	// Use existing BigSqrt implementation
	return BigSqrt(x, prec)
}
