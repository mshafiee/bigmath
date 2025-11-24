// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

//go:build !amd64 && !arm64

package bigmath

// Generic fallback implementations for platforms without assembly support

//nolint:unused // Declared for consistency with amd64/arm64 versions
func bigExpAsm(x *BigFloat, prec uint) *BigFloat {
	return bigExpGeneric(x, prec)
}

//nolint:unused // Declared for consistency with amd64/arm64 versions
func bigLogAsm(x *BigFloat, prec uint) *BigFloat {
	return bigLogGeneric(x, prec)
}
