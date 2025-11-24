// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

//go:build amd64

package bigmath

// AMD64 wrapper functions (currently calling generic implementations)
// Future: Replace with true assembly optimizations

func bigFloorAsm(x *BigFloat, prec uint) *BigFloat {
	return bigFloorOptimized(x, prec)
}

func bigCeilAsm(x *BigFloat, prec uint) *BigFloat {
	return bigCeilOptimized(x, prec)
}

func bigTruncAsm(x *BigFloat, prec uint) *BigFloat {
	return bigTruncOptimized(x, prec)
}

func bigModAsm(x, y *BigFloat, prec uint) *BigFloat {
	return bigModOptimized(x, y, prec)
}

func bigRemAsm(x, y *BigFloat, prec uint) *BigFloat {
	return bigRemOptimized(x, y, prec)
}
