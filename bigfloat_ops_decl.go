// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

//go:build amd64

package bigmath

// Assembly function declarations
// These will be implemented in architecture-specific .s files (AMD64).

//go:noescape
//nolint:unused // Declared for assembly implementation
func bigfloatAddAsm(a, b *BigFloat, prec uint) *BigFloat

//go:noescape
//nolint:unused // Declared for assembly implementation
func bigfloatSubAsm(a, b *BigFloat, prec uint) *BigFloat

//go:noescape
//nolint:unused // Declared for assembly implementation
func bigfloatMulAsm(a, b *BigFloat, prec uint) *BigFloat

//go:noescape
//nolint:unused // Declared for assembly implementation
func bigfloatDivAsm(a, b *BigFloat, prec uint) *BigFloat

//go:noescape
//nolint:unused // Declared for assembly implementation
func bigfloatSqrtAsm(x *BigFloat, prec uint) *BigFloat
