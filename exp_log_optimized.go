// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

import (
	"math"
	"math/big"
)

// Optimized exponential and logarithmic functions with reduced allocations

// expWorkspace holds pre-allocated buffers for exponential calculations
type expWorkspace struct {
	result    *BigFloat
	term      *BigFloat
	rReduced  *BigFloat
	threshold *BigFloat
	scale     *BigFloat
	temp      *BigFloat
	ln2       *BigFloat
	kBig      *BigFloat
	prec      uint
}

// getExpWorkspace returns a workspace with pre-allocated buffers
func getExpWorkspace(prec uint) *expWorkspace {
	if prec == 0 {
		prec = DefaultPrecision
	}
	workPrec := prec + 32
	return &expWorkspace{
		result:    NewBigFloat(0.0, workPrec),
		term:      NewBigFloat(0.0, workPrec),
		rReduced:  NewBigFloat(0.0, workPrec),
		threshold: NewBigFloat(0.0, workPrec),
		scale:     NewBigFloat(0.0, workPrec),
		temp:      NewBigFloat(0.0, workPrec),
		ln2:       BigLog2(workPrec),
		kBig:      NewBigFloat(0.0, workPrec),
		prec:      prec,
	}
}

// bigExpOptimized computes e^x with optimized allocation pattern
func bigExpOptimized(x *BigFloat, prec uint) *BigFloat {
	if prec == 0 {
		prec = x.Prec()
	}

	// Handle special cases
	if x.IsInf() {
		if x.Sign() > 0 {
			return new(BigFloat).SetPrec(prec).SetInf(false)
		}
		return new(BigFloat).SetPrec(prec).SetFloat64(0.0)
	}
	if x.Sign() == 0 {
		return NewBigFloat(1.0, prec)
	}

	workPrec := prec + 32
	ws := getExpWorkspace(prec)

	// 1. Argument reduction: x = k*ln(2) + r
	kFloat := new(BigFloat).SetPrec(workPrec).Quo(x, ws.ln2)
	kInt := new(big.Int)
	kFloat.Int(kInt) // Round to nearest integer

	// r = x - k*ln(2)
	ws.kBig.SetInt(kInt)
	r := new(BigFloat).SetPrec(workPrec).Mul(ws.kBig, ws.ln2)
	r.Sub(x, r)

	// 2. Further reduction: exp(r) = (exp(r/2^S))^(2^S)
	rAbs := new(BigFloat).SetPrec(workPrec).Abs(r)
	rFloat, _ := rAbs.Float64()
	S := 0
	if rFloat > 0 {
		S = int(math.Ceil(math.Log2(rFloat) + 14))
	}
	if S < 0 {
		S = 0
	}

	// rReduced = r / 2^S
	ws.scale.SetInt64(1)
	ws.scale.SetMantExp(ws.scale, S) // 2^S
	ws.rReduced.Quo(r, ws.scale)

	// 3. Taylor series for exp(rReduced) with optimized loop
	ws.result.SetFloat64(1.0)
	ws.term.SetFloat64(1.0)

	ws.threshold.SetMantExp(NewBigFloat(1.0, workPrec), -int(workPrec))

	for n := 1; n < 1000; n++ {
		// term = term * rReduced / n
		ws.term.Mul(ws.term, ws.rReduced)
		ws.temp.SetFloat64(float64(n))
		ws.term.Quo(ws.term, ws.temp)

		ws.result.Add(ws.result, ws.term)

		// Check convergence using pre-allocated buffer
		ws.temp.Abs(ws.term)
		if ws.temp.Cmp(ws.threshold) < 0 {
			break
		}
	}

	// 4. Square S times: res = res^(2^S)
	for i := 0; i < S; i++ {
		ws.result.Mul(ws.result, ws.result)
	}

	// 5. Multiply by 2^k
	if kInt.Sign() != 0 {
		kVal := kInt.Int64()
		mant := new(BigFloat).SetPrec(workPrec)
		exp := ws.result.MantExp(mant)
		ws.result.SetMantExp(mant, exp+int(kVal))
	}

	return new(BigFloat).SetPrec(prec).Set(ws.result)
}

// logWorkspace holds pre-allocated buffers for logarithmic calculations
type logWorkspace struct {
	result    *BigFloat
	term      *BigFloat
	xReduced  *BigFloat
	threshold *BigFloat
	temp      *BigFloat
	prec      uint
}

// getLogWorkspace returns a workspace with pre-allocated buffers
func getLogWorkspace(prec uint) *logWorkspace {
	if prec == 0 {
		prec = DefaultPrecision
	}
	workPrec := prec + 32
	return &logWorkspace{
		result:    NewBigFloat(0.0, workPrec),
		term:      NewBigFloat(0.0, workPrec),
		xReduced:  NewBigFloat(0.0, workPrec),
		threshold: NewBigFloat(0.0, workPrec),
		temp:      NewBigFloat(0.0, workPrec),
		prec:      prec,
	}
}

// bigLogOptimized computes ln(x) with optimized allocation pattern
func bigLogOptimized(x *BigFloat, prec uint) *BigFloat {
	if prec == 0 {
		prec = x.Prec()
	}

	// Handle special cases
	if x.Sign() == 0 {
		return new(BigFloat).SetPrec(prec).SetInf(true) // -infinity
	}
	if x.Sign() < 0 {
		return new(BigFloat).SetPrec(prec).SetFloat64(math.NaN())
	}
	if x.IsInf() {
		return new(BigFloat).SetPrec(prec).SetInf(false) // +infinity
	}

	workPrec := prec + 32
	ws := getLogWorkspace(prec)

	// Range reduction to [1, 2)
	// Extract exponent and normalize mantissa
	mant := new(BigFloat).SetPrec(workPrec)
	exp := x.MantExp(mant)

	// x = mant * 2^exp, so ln(x) = ln(mant) + exp*ln(2)
	// Normalize mant to [1, 2)
	if mant.Cmp(NewBigFloat(1.0, workPrec)) < 0 {
		mant.Mul(mant, NewBigFloat(2.0, workPrec))
		exp--
	}

	// Now mant is in [1, 2), compute ln(mant) using Taylor series
	// ln(1+y) = y - y²/2 + y³/3 - ... where y = mant - 1
	ws.xReduced.Sub(mant, NewBigFloat(1.0, workPrec)) // y = mant - 1

	// Taylor series: ln(1+y) = sum (-1)^(n+1) * y^n / n
	ws.result.Set(ws.xReduced) // First term is y
	ws.term.Set(ws.xReduced)

	ws.threshold.SetMantExp(NewBigFloat(1.0, workPrec), -int(prec))

	for n := 2; n < 1000; n++ {
		// term = term * y
		ws.term.Mul(ws.term, ws.xReduced)

		// Add term with alternating sign
		ws.temp.Set(ws.term)
		ws.temp.Quo(ws.temp, NewBigFloat(float64(n), workPrec))
		if n%2 == 0 {
			ws.temp.Neg(ws.temp)
		}
		ws.result.Add(ws.result, ws.temp)

		// Check convergence
		ws.temp.Abs(ws.temp)
		if ws.temp.Cmp(ws.threshold) < 0 {
			break
		}
	}

	// Add exp*ln(2)
	if exp != 0 {
		ln2 := BigLog2(workPrec)
		expBig := NewBigFloat(float64(exp), workPrec)
		expTerm := new(BigFloat).SetPrec(workPrec).Mul(expBig, ln2)
		ws.result.Add(ws.result, expTerm)
	}

	return new(BigFloat).SetPrec(prec).Set(ws.result)
}

