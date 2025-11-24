// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

//go:build !amd64 && !arm64

package bigmath

//go:noescape
func bigExpAsm(x *BigFloat, prec uint) *BigFloat

//go:noescape
func bigLogAsm(x *BigFloat, prec uint) *BigFloat
