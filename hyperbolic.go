package bigmath

import "math"

// BigSinh computes sinh(x) = (e^x - e^-x) / 2
func BigSinh(x *BigFloat, prec uint) *BigFloat {
	// Use dispatcher to select assembly or generic implementation
	return getDispatcher().BigSinhImpl(x, prec)
}

// bigSinhGeneric is the generic implementation (called by dispatcher)
func bigSinhGeneric(x *BigFloat, prec uint) *BigFloat {
	if prec == 0 {
		prec = x.Prec()
	}

	// Special cases
	if x.Sign() == 0 {
		return NewBigFloat(0.0, prec)
	}

	workPrec := prec + 32

	// Calculate exp(x) - use dispatcher directly to avoid recursion
	expX := getDispatcher().BigExpImpl(x, workPrec)

	// Calculate exp(-x) = 1/exp(x)
	one := NewBigFloat(1.0, workPrec)
	expNegX := new(BigFloat).SetPrec(workPrec).Quo(one, expX)

	// (exp(x) - exp(-x)) / 2
	res := new(BigFloat).SetPrec(workPrec).Sub(expX, expNegX)
	res.Quo(res, NewBigFloat(2.0, workPrec))

	return new(BigFloat).SetPrec(prec).Set(res)
}

// BigCosh computes cosh(x) = (e^x + e^-x) / 2
func BigCosh(x *BigFloat, prec uint) *BigFloat {
	// Use dispatcher to select assembly or generic implementation
	return getDispatcher().BigCoshImpl(x, prec)
}

// bigCoshGeneric is the generic implementation (called by dispatcher)
func bigCoshGeneric(x *BigFloat, prec uint) *BigFloat {
	if prec == 0 {
		prec = x.Prec()
	}

	// Special cases
	if x.Sign() == 0 {
		return NewBigFloat(1.0, prec)
	}

	workPrec := prec + 32

	// Use dispatcher directly to avoid recursion
	expX := getDispatcher().BigExpImpl(x, workPrec)
	one := NewBigFloat(1.0, workPrec)
	expNegX := new(BigFloat).SetPrec(workPrec).Quo(one, expX)

	res := new(BigFloat).SetPrec(workPrec).Add(expX, expNegX)
	res.Quo(res, NewBigFloat(2.0, workPrec))

	return new(BigFloat).SetPrec(prec).Set(res)
}

// BigTanh computes tanh(x) = (e^2x - 1) / (e^2x + 1)
func BigTanh(x *BigFloat, prec uint) *BigFloat {
	// Use dispatcher to select assembly or generic implementation
	return getDispatcher().BigTanhImpl(x, prec)
}

// bigTanhGeneric is the generic implementation (called by dispatcher)
func bigTanhGeneric(x *BigFloat, prec uint) *BigFloat {
	if prec == 0 {
		prec = x.Prec()
	}

	if x.Sign() == 0 {
		return NewBigFloat(0.0, prec)
	}

	workPrec := prec + 32

	// e^2x - use dispatcher directly to avoid recursion
	twoX := new(BigFloat).SetPrec(workPrec).Mul(x, NewBigFloat(2.0, workPrec))
	exp2X := getDispatcher().BigExpImpl(twoX, workPrec)

	one := NewBigFloat(1.0, workPrec)

	// num = e^2x - 1
	num := new(BigFloat).SetPrec(workPrec).Sub(exp2X, one)

	// den = e^2x + 1
	den := new(BigFloat).SetPrec(workPrec).Add(exp2X, one)

	res := new(BigFloat).SetPrec(workPrec).Quo(num, den)

	return new(BigFloat).SetPrec(prec).Set(res)
}

// BigAsinh computes asinh(x) = ln(x + sqrt(x^2 + 1))
func BigAsinh(x *BigFloat, prec uint) *BigFloat {
	// Use dispatcher to select assembly or generic implementation
	return getDispatcher().BigAsinhImpl(x, prec)
}

// bigAsinhGeneric is the generic implementation (called by dispatcher)
func bigAsinhGeneric(x *BigFloat, prec uint) *BigFloat {
	if prec == 0 {
		prec = x.Prec()
	}

	workPrec := prec + 32

	x2 := new(BigFloat).SetPrec(workPrec).Mul(x, x)
	one := NewBigFloat(1.0, workPrec)
	x2Plus1 := new(BigFloat).SetPrec(workPrec).Add(x2, one)
	sqrt := BigSqrt(x2Plus1, workPrec)

	arg := new(BigFloat).SetPrec(workPrec).Add(x, sqrt)
	// Use dispatcher directly to avoid recursion
	res := getDispatcher().BigLogImpl(arg, workPrec)

	return new(BigFloat).SetPrec(prec).Set(res)
}

// BigAcosh computes acosh(x) = ln(x + sqrt(x^2 - 1))
func BigAcosh(x *BigFloat, prec uint) *BigFloat {
	// Use dispatcher to select assembly or generic implementation
	return getDispatcher().BigAcoshImpl(x, prec)
}

// bigAcoshGeneric is the generic implementation (called by dispatcher)
func bigAcoshGeneric(x *BigFloat, prec uint) *BigFloat {
	if prec == 0 {
		prec = x.Prec()
	}

	// Domain check: x >= 1
	one := NewBigFloat(1.0, prec)
	if x.Cmp(one) < 0 {
		return new(BigFloat).SetPrec(prec).SetFloat64(math.NaN()) // NaN
	}

	workPrec := prec + 32

	x2 := new(BigFloat).SetPrec(workPrec).Mul(x, x)
	oneW := NewBigFloat(1.0, workPrec)
	x2Minus1 := new(BigFloat).SetPrec(workPrec).Sub(x2, oneW)
	sqrt := BigSqrt(x2Minus1, workPrec)

	arg := new(BigFloat).SetPrec(workPrec).Add(x, sqrt)
	// Use dispatcher directly to avoid recursion
	res := getDispatcher().BigLogImpl(arg, workPrec)

	return new(BigFloat).SetPrec(prec).Set(res)
}

// BigAtanh computes atanh(x) = 0.5 * ln((1+x)/(1-x))
func BigAtanh(x *BigFloat, prec uint) *BigFloat {
	// Use dispatcher to select assembly or generic implementation
	return getDispatcher().BigAtanhImpl(x, prec)
}

// bigAtanhGeneric is the generic implementation (called by dispatcher)
func bigAtanhGeneric(x *BigFloat, prec uint) *BigFloat {
	if prec == 0 {
		prec = x.Prec()
	}

	// Domain check: |x| < 1
	one := NewBigFloat(1.0, prec)
	absX := new(BigFloat).SetPrec(prec).Abs(x)
	if absX.Cmp(one) >= 0 {
		return new(BigFloat).SetPrec(prec).SetFloat64(math.NaN())
	}

	workPrec := prec + 32

	oneW := NewBigFloat(1.0, workPrec)
	num := new(BigFloat).SetPrec(workPrec).Add(oneW, x)
	den := new(BigFloat).SetPrec(workPrec).Sub(oneW, x)

	ratio := new(BigFloat).SetPrec(workPrec).Quo(num, den)
	ln := BigLog(ratio, workPrec)

	res := new(BigFloat).SetPrec(workPrec).Mul(ln, NewBigFloat(0.5, workPrec))

	return new(BigFloat).SetPrec(prec).Set(res)
}
