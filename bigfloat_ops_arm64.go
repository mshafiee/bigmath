// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

//go:build arm64

package bigmath

// ARM64 fallback implementations that forward to the generic Go versions.

func bigfloatAddAsm(a, b *BigFloat, prec uint) *BigFloat {
	return bigfloatAddGeneric(a, b, prec)
}

func bigfloatSubAsm(a, b *BigFloat, prec uint) *BigFloat {
	return bigfloatSubGeneric(a, b, prec)
}

func bigfloatMulAsm(a, b *BigFloat, prec uint) *BigFloat {
	return bigfloatMulGeneric(a, b, prec)
}

func bigfloatDivAsm(a, b *BigFloat, prec uint) *BigFloat {
	return bigfloatDivGeneric(a, b, prec)
}

func bigfloatSqrtAsm(x *BigFloat, prec uint) *BigFloat {
	return bigfloatSqrtGeneric(x, prec)
}
