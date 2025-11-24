// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

import (
	"math"
	"testing"
)

func TestBigPhi(t *testing.T) {
	prec := uint(256)

	phi := BigPhi(prec)
	phiVal, _ := phi.Float64()

	expected := (1.0 + math.Sqrt(5.0)) / 2.0
	if math.Abs(phiVal-expected) > 1e-10 {
		t.Errorf("BigPhi = %g, want %g", phiVal, expected)
	}
}

func TestBigSqrt2(t *testing.T) {
	prec := uint(256)

	sqrt2 := BigSqrt2(prec)
	sqrt2Val, _ := sqrt2.Float64()

	expected := math.Sqrt(2.0)
	if math.Abs(sqrt2Val-expected) > 1e-10 {
		t.Errorf("BigSqrt2 = %g, want %g", sqrt2Val, expected)
	}
}

func TestBigSqrt3(t *testing.T) {
	prec := uint(256)

	sqrt3 := BigSqrt3(prec)
	sqrt3Val, _ := sqrt3.Float64()

	expected := math.Sqrt(3.0)
	if math.Abs(sqrt3Val-expected) > 1e-10 {
		t.Errorf("BigSqrt3 = %g, want %g", sqrt3Val, expected)
	}
}

func TestBigLn10(t *testing.T) {
	prec := uint(256)

	ln10 := BigLn10(prec)
	ln10Val, _ := ln10.Float64()

	expected := math.Log(10.0)
	if math.Abs(ln10Val-expected) > 1e-8 {
		t.Errorf("BigLn10 = %g, want %g", ln10Val, expected)
	}
}
