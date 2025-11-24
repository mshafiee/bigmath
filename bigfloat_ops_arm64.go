// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

//go:build arm64

package bigmath

// ARM64 fallback implementations that forward to the generic Go versions.

//nolint:unused // Declared in bigfloat_ops_decl.go, may be called from assembly
func bigfloatAddAsm(a, b *BigFloat, prec uint) *BigFloat {
	return bigfloatAddGeneric(a, b, prec)
}

//nolint:unused // Declared in bigfloat_ops_decl.go, may be called from assembly
func bigfloatSubAsm(a, b *BigFloat, prec uint) *BigFloat {
	return bigfloatSubGeneric(a, b, prec)
}

//nolint:unused // Declared in bigfloat_ops_decl.go, may be called from assembly
func bigfloatMulAsm(a, b *BigFloat, prec uint) *BigFloat {
	return bigfloatMulGeneric(a, b, prec)
}

//nolint:unused // Declared in bigfloat_ops_decl.go, may be called from assembly
func bigfloatDivAsm(a, b *BigFloat, prec uint) *BigFloat {
	return bigfloatDivGeneric(a, b, prec)
}

//nolint:unused // Declared in bigfloat_ops_decl.go, may be called from assembly
func bigfloatSqrtAsm(x *BigFloat, prec uint) *BigFloat {
	return bigfloatSqrtGeneric(x, prec)
}
