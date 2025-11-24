// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

//go:build amd64

package bigmath

// AMD64 wrapper functions for special functions
// Now using optimized implementations with early convergence detection

func bigGammaAsm(x *BigFloat, prec uint) *BigFloat {
	return bigGammaOptimized(x, prec)
}

func bigErfAsm(x *BigFloat, prec uint) *BigFloat {
	return bigErfOptimized(x, prec)
}

func bigErfcAsm(x *BigFloat, prec uint) *BigFloat {
	return bigErfcOptimized(x, prec)
}

func bigBesselJAsm(n int, x *BigFloat, prec uint) *BigFloat {
	return bigBesselJOptimized(n, x, prec)
}

func bigBesselYAsm(n int, x *BigFloat, prec uint) *BigFloat {
	return bigBesselYOptimized(n, x, prec)
}
