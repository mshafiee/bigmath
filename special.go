// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

import (
	"math"
)

// BigGamma computes the Gamma function Γ(x)
// Uses Lanczos approximation for x > 0
// For negative x, uses reflection formula: Γ(x) = π / (Γ(1-x) * sin(π*x))
func BigGamma(x *BigFloat, prec uint) *BigFloat {
	if prec == 0 {
		prec = x.Prec()
	}

	// Handle special cases
	if x.Sign() == 0 {
		return new(BigFloat).SetPrec(prec).SetInf(false)
	}
	if x.IsInf() {
		if x.Sign() > 0 {
			return new(BigFloat).SetPrec(prec).SetInf(false)
		}
		return NewBigFloat(math.NaN(), prec)
	}

	workPrec := prec + 32

	// Check if x is negative
	if x.Sign() < 0 {
		// Use reflection formula: Γ(x) = π / (Γ(1-x) * sin(π*x))
		one := NewBigFloat(1.0, workPrec)
		oneMinusX := new(BigFloat).SetPrec(workPrec).Sub(one, x)
		gamma1MinusX := bigGammaPositive(oneMinusX, workPrec)

		piX := new(BigFloat).SetPrec(workPrec).Mul(BigPI(workPrec), x)
		sinPiX := BigSin(piX, workPrec)

		denom := new(BigFloat).SetPrec(workPrec).Mul(gamma1MinusX, sinPiX)
		result := new(BigFloat).SetPrec(workPrec).Quo(BigPI(workPrec), denom)

		return new(BigFloat).SetPrec(prec).Set(result)
	}

	return new(BigFloat).SetPrec(prec).Set(bigGammaPositive(x, workPrec))
}

// bigGammaPositive computes Gamma for positive x using Lanczos approximation
// Lanczos approximation: Γ(z) ≈ sqrt(2π) * (z+g-0.5)^(z-0.5) * exp(-(z+g-0.5)) * A(z)
// where A(z) is a series approximation
func bigGammaPositive(x *BigFloat, prec uint) *BigFloat {
	// Lanczos coefficients (g=7, n=9 terms)
	// These are precomputed constants for the Lanczos approximation
	g := 7.0
	coeffs := []float64{
		0.99999999999980993,
		676.5203681218851,
		-1259.1392167224028,
		771.32342877765313,
		-176.61502916214059,
		12.507343278686905,
		-0.13857109526572012,
		9.9843695780195716e-6,
		1.5056327351493116e-7,
	}

	// Check if x is small (< 0.5), use reflection and recursion
	half := NewBigFloat(0.5, prec)
	if x.Cmp(half) < 0 {
		// Use reflection formula: Γ(x) = π / (Γ(1-x) * sin(π*x))
		one := NewBigFloat(1.0, prec)
		oneMinusX := new(BigFloat).SetPrec(prec).Sub(one, x)
		gamma1MinusX := bigGammaPositive(oneMinusX, prec)

		piX := new(BigFloat).SetPrec(prec).Mul(BigPI(prec), x)
		sinPiX := BigSin(piX, prec)

		denom := new(BigFloat).SetPrec(prec).Mul(gamma1MinusX, sinPiX)
		result := new(BigFloat).SetPrec(prec).Quo(BigPI(prec), denom)

		return result
	}

	// Reduce x by 1 to bring into range [0.5, 1.5]
	z := new(BigFloat).SetPrec(prec).Set(x)
	one := NewBigFloat(1.0, prec)
	reduction := 0

	for z.Cmp(NewBigFloat(1.5, prec)) >= 0 {
		z.Sub(z, one)
		reduction++
	}

	// Now z is in [0.5, 1.5), compute Γ(z) using Lanczos
	// A(z) = sum_{k=0}^{n-1} coeffs[k] / (z + k)
	zPlusG := new(BigFloat).SetPrec(prec).Add(z, NewBigFloat(g, prec))
	zPlusGMinusHalf := new(BigFloat).SetPrec(prec).Sub(zPlusG, NewBigFloat(0.5, prec))

	// Compute A(z)
	aZ := NewBigFloat(coeffs[0], prec)
	for k := 1; k < len(coeffs); k++ {
		zPlusK := new(BigFloat).SetPrec(prec).Add(z, NewBigFloat(float64(k-1), prec))
		term := new(BigFloat).SetPrec(prec).Quo(NewBigFloat(coeffs[k], prec), zPlusK)
		aZ.Add(aZ, term)
	}

	// Compute (z+g-0.5)^(z-0.5)
	zMinusHalf := new(BigFloat).SetPrec(prec).Sub(z, NewBigFloat(0.5, prec))
	power := BigPow(zPlusGMinusHalf, zMinusHalf, prec)

	// Compute exp(-(z+g-0.5))
	negZPlusGMinusHalf := new(BigFloat).SetPrec(prec).Neg(zPlusGMinusHalf)
	expTerm := BigExp(negZPlusGMinusHalf, prec)

	// Compute sqrt(2π)
	twoPi := BigTwoPI(prec)
	sqrtTwoPi := BigSqrt(twoPi, prec)

	// Combine: sqrt(2π) * power * expTerm * aZ
	result := new(BigFloat).SetPrec(prec).Mul(sqrtTwoPi, power)
	result.Mul(result, expTerm)
	result.Mul(result, aZ)

	// Multiply by (z+1)*(z+2)*...*(z+reduction) to account for reduction
	for i := 0; i < reduction; i++ {
		zPlusI := new(BigFloat).SetPrec(prec).Add(z, NewBigFloat(float64(i), prec))
		result.Mul(result, zPlusI)
	}

	return result
}

// BigErf computes the error function erf(x) = (2/√π) * ∫[0 to x] exp(-t²) dt
// Uses series expansion for small |x|, asymptotic expansion for large |x|
func BigErf(x *BigFloat, prec uint) *BigFloat {
	if prec == 0 {
		prec = x.Prec()
	}

	// Handle special cases
	if x.Sign() == 0 {
		return NewBigFloat(0.0, prec)
	}
	if x.IsInf() {
		if x.Sign() > 0 {
			return NewBigFloat(1.0, prec)
		}
		return NewBigFloat(-1.0, prec)
	}

	workPrec := prec + 32

	// For small |x|, use series expansion
	// erf(x) = (2/√π) * sum_{n=0} (-1)^n * x^(2n+1) / (n! * (2n+1))
	xAbs := BigAbs(x, workPrec)
	smallThreshold := NewBigFloat(0.8, workPrec)
	moderateThreshold := NewBigFloat(2.0, workPrec)

	if xAbs.Cmp(smallThreshold) < 0 {
		// Very small x: use direct series expansion
		return bigErfSeries(x, workPrec, prec)
	}

	if xAbs.Cmp(moderateThreshold) < 0 {
		// Moderate x: use improved erfc for better accuracy
		// Compute as 1 - erfc(x) using improved erfc
		if x.Sign() > 0 {
			one := NewBigFloat(1.0, workPrec)
			erfcX := bigErfcImproved(x, workPrec, workPrec)
			result := new(BigFloat).SetPrec(workPrec).Sub(one, erfcX)
			return new(BigFloat).SetPrec(prec).Set(result)
		} else {
			negOne := NewBigFloat(-1.0, workPrec)
			negX := new(BigFloat).SetPrec(workPrec).Neg(x)
			erfcNegX := bigErfcImproved(negX, workPrec, workPrec)
			result := new(BigFloat).SetPrec(workPrec).Add(negOne, erfcNegX)
			return new(BigFloat).SetPrec(prec).Set(result)
		}
	}

	// For large |x|, use erfc via complementary error function
	// erf(x) ≈ 1 - erfc(x) for x > 0
	// erf(x) ≈ -1 + erfc(-x) for x < 0
	if x.Sign() > 0 {
		one := NewBigFloat(1.0, workPrec)
		erfcX := BigErfc(x, workPrec)
		result := new(BigFloat).SetPrec(workPrec).Sub(one, erfcX)
		return new(BigFloat).SetPrec(prec).Set(result)
	} else {
		negOne := NewBigFloat(-1.0, workPrec)
		negX := new(BigFloat).SetPrec(workPrec).Neg(x)
		erfcNegX := BigErfc(negX, workPrec)
		result := new(BigFloat).SetPrec(workPrec).Add(negOne, erfcNegX)
		return new(BigFloat).SetPrec(prec).Set(result)
	}
}

// bigErfSeries computes erf(x) using series expansion for small |x|
func bigErfSeries(x *BigFloat, workPrec, targetPrec uint) *BigFloat {
	// Series expansion: erf(x) = (2/√π) * sum_{n=0} (-1)^n * x^(2n+1) / (n! * (2n+1))
	twoOverSqrtPi := new(BigFloat).SetPrec(workPrec).Quo(NewBigFloat(2.0, workPrec), BigSqrt(BigPI(workPrec), workPrec))

	result := new(BigFloat).SetPrec(workPrec).Set(x)
	term := new(BigFloat).SetPrec(workPrec).Set(x)
	x2 := new(BigFloat).SetPrec(workPrec).Mul(x, x)
	factorial := NewBigFloat(1.0, workPrec)

	// More lenient convergence threshold for better accuracy
	convThreshold := new(BigFloat).SetPrec(workPrec).SetMantExp(NewBigFloat(1.0, workPrec), -int(workPrec+15))
	resultAbs := new(BigFloat).SetPrec(workPrec).Abs(result)

	for n := 1; n < 3000; n++ {
		// term = (-1)^n * x^(2n+1) / (n! * (2n+1))
		term.Mul(term, x2)
		factorial.Mul(factorial, NewBigFloat(float64(n), workPrec))
		denom := new(BigFloat).SetPrec(workPrec).Mul(factorial, NewBigFloat(float64(2*n+1), workPrec))
		term.Quo(term, denom)

		if n%2 == 1 {
			// Odd n: subtract
			result.Sub(result, term)
		} else {
			// Even n: add
			result.Add(result, term)
		}

		// Check convergence relative to result magnitude
		termAbs := new(BigFloat).SetPrec(workPrec).Abs(term)
		if termAbs.Cmp(convThreshold) < 0 {
			// Also check relative convergence
			if resultAbs.Sign() > 0 {
				relErr := new(BigFloat).SetPrec(workPrec).Quo(termAbs, resultAbs)
				if relErr.Cmp(convThreshold) < 0 {
					break
				}
			} else {
				break
			}
		}
		resultAbs.Set(result)
		resultAbs.Abs(resultAbs)
	}

	result.Mul(result, twoOverSqrtPi)
	return new(BigFloat).SetPrec(targetPrec).Set(result)
}

// bigErfcImproved computes erfc(x) with improved accuracy for moderate x
func bigErfcImproved(x *BigFloat, workPrec, targetPrec uint) *BigFloat {
	// For moderate x (0.5 <= x < 4.0), use improved asymptotic expansion
	// with more terms and better convergence

	x2 := new(BigFloat).SetPrec(workPrec).Mul(x, x)
	expNegX2 := BigExp(new(BigFloat).SetPrec(workPrec).Neg(x2), workPrec)

	sqrtPi := BigSqrt(BigPI(workPrec), workPrec)
	xSqrtPi := new(BigFloat).SetPrec(workPrec).Mul(x, sqrtPi)
	base := new(BigFloat).SetPrec(workPrec).Quo(expNegX2, xSqrtPi)

	// Asymptotic series with more terms and better convergence
	series := NewBigFloat(1.0, workPrec)
	x2Inv := new(BigFloat).SetPrec(workPrec).Quo(NewBigFloat(1.0, workPrec), x2)

	// First term: -1/(2x²)
	term := new(BigFloat).SetPrec(workPrec).Set(x2Inv)
	term.Mul(term, NewBigFloat(-0.5, workPrec))

	// More strict convergence threshold for better accuracy
	convThreshold := new(BigFloat).SetPrec(workPrec).SetMantExp(NewBigFloat(1.0, workPrec), -int(workPrec+15))

	for n := 1; n < 300; n++ {
		series.Add(series, term)

		// Next term: a_{n+1} = a_n * (2n+1)/(2n+2) * (-1) / x²
		coeff := new(BigFloat).SetPrec(workPrec).Quo(
			NewBigFloat(float64(2*n+1), workPrec),
			NewBigFloat(float64(2*n+2), workPrec),
		)
		term.Mul(term, coeff)
		term.Mul(term, x2Inv)
		term.Neg(term) // Alternating sign

		termAbs := new(BigFloat).SetPrec(workPrec).Abs(term)
		if termAbs.Cmp(convThreshold) < 0 {
			// Also check relative convergence
			seriesAbs := new(BigFloat).SetPrec(workPrec).Abs(series)
			if seriesAbs.Sign() > 0 {
				relErr := new(BigFloat).SetPrec(workPrec).Quo(termAbs, seriesAbs)
				if relErr.Cmp(convThreshold) < 0 {
					break
				}
			} else {
				break
			}
		}
	}

	result := new(BigFloat).SetPrec(workPrec).Mul(base, series)
	return new(BigFloat).SetPrec(targetPrec).Set(result)
}

// BigErfc computes the complementary error function erfc(x) = 1 - erf(x)
// Uses continued fraction for large |x|
func BigErfc(x *BigFloat, prec uint) *BigFloat {
	if prec == 0 {
		prec = x.Prec()
	}

	// Handle special cases
	if x.Sign() == 0 {
		return NewBigFloat(1.0, prec)
	}
	if x.IsInf() {
		if x.Sign() > 0 {
			return NewBigFloat(0.0, prec)
		}
		return NewBigFloat(2.0, prec)
	}

	workPrec := prec + 32

	// For negative x, use erfc(-x) = 2 - erfc(x)
	if x.Sign() < 0 {
		negX := new(BigFloat).SetPrec(workPrec).Neg(x)
		erfcNegX := BigErfc(negX, workPrec)
		two := NewBigFloat(2.0, workPrec)
		result := new(BigFloat).SetPrec(workPrec).Sub(two, erfcNegX)
		return new(BigFloat).SetPrec(prec).Set(result)
	}

	// For small x, compute as 1 - erf(x)
	xAbs := BigAbs(x, workPrec)
	smallThreshold := NewBigFloat(0.8, workPrec)

	if xAbs.Cmp(smallThreshold) < 0 {
		// Very small x: use 1 - erf(x) for better accuracy
		one := NewBigFloat(1.0, workPrec)
		erfX := BigErf(x, workPrec)
		result := new(BigFloat).SetPrec(workPrec).Sub(one, erfX)
		return new(BigFloat).SetPrec(prec).Set(result)
	}

	// For moderate and large x, use improved erfc computation
	return bigErfcImproved(x, workPrec, prec)
}

// BigBesselJ computes the Bessel function of the first kind J_n(x)
// Uses series expansion: J_n(x) = sum_{k=0} (-1)^k * (x/2)^(n+2k) / (k! * (n+k)!)
func BigBesselJ(n int, x *BigFloat, prec uint) *BigFloat {
	if prec == 0 {
		prec = x.Prec()
	}

	// Handle special cases
	if x.Sign() == 0 {
		if n == 0 {
			return NewBigFloat(1.0, prec)
		}
		return NewBigFloat(0.0, prec)
	}
	if x.IsInf() {
		return NewBigFloat(0.0, prec) // Bessel functions decay at infinity
	}

	workPrec := prec + 32

	// Handle negative n using: J_{-n}(x) = (-1)^n * J_n(x)
	if n < 0 {
		result := BigBesselJ(-n, x, workPrec)
		if (-n)%2 != 0 {
			result.Neg(result)
		}
		return new(BigFloat).SetPrec(prec).Set(result)
	}

	// For large x, use asymptotic expansion for better accuracy
	xAbs := BigAbs(x, workPrec)
	largeXThreshold := NewBigFloat(10.0, workPrec)

	if xAbs.Cmp(largeXThreshold) > 0 {
		// Use asymptotic expansion for large x
		// J_n(x) ≈ sqrt(2/(πx)) * cos(x - nπ/2 - π/4) * (1 + O(1/x))
		// For better accuracy, use more terms of asymptotic expansion
		sqrtTwoOverPiX := BigSqrt(new(BigFloat).SetPrec(workPrec).Quo(NewBigFloat(2.0, workPrec),
			new(BigFloat).SetPrec(workPrec).Mul(BigPI(workPrec), x)), workPrec)

		// Phase: x - nπ/2 - π/4
		nPiOver2 := new(BigFloat).SetPrec(workPrec).Mul(NewBigFloat(float64(n), workPrec), BigHalfPI(workPrec))
		phase := new(BigFloat).SetPrec(workPrec).Sub(x, nPiOver2)
		phase.Sub(phase, BigHalfPI(workPrec))

		cosPhase := BigCos(phase, workPrec)
		result := new(BigFloat).SetPrec(workPrec).Mul(sqrtTwoOverPiX, cosPhase)

		// Add correction terms for better accuracy (first order correction)
		// J_n(x) ≈ sqrt(2/(πx)) * [cos(phase) * (1 - (4n²-1)/(8x)) - sin(phase) * (4n²-1)/(8x)]
		if n > 0 || xAbs.Cmp(NewBigFloat(20.0, workPrec)) < 0 {
			// For moderate x or small n, use series expansion instead
			return bigBesselJSeries(n, x, workPrec, prec)
		}

		correction := new(BigFloat).SetPrec(workPrec).Quo(
			NewBigFloat(float64(4*n*n-1), workPrec),
			new(BigFloat).SetPrec(workPrec).Mul(NewBigFloat(8.0, workPrec), x),
		)
		sinPhase := BigSin(phase, workPrec)
		corrTerm := new(BigFloat).SetPrec(workPrec).Mul(correction, sinPhase)
		result.Sub(result, corrTerm)

		return new(BigFloat).SetPrec(prec).Set(result)
	}

	// For small to moderate x, use series expansion
	return bigBesselJSeries(n, x, workPrec, prec)
}

// bigBesselJSeries computes J_n(x) using series expansion
func bigBesselJSeries(n int, x *BigFloat, workPrec, targetPrec uint) *BigFloat {
	// Series expansion: J_n(x) = sum_{k=0} (-1)^k * (x/2)^(n+2k) / (k! * (n+k)!)
	xHalf := new(BigFloat).SetPrec(workPrec).Quo(x, NewBigFloat(2.0, workPrec))
	xHalfPower := BigPow(xHalf, NewBigFloat(float64(n), workPrec), workPrec)

	// Compute n! for denominator
	nFactorial := NewBigFloat(1.0, workPrec)
	for i := 1; i <= n; i++ {
		nFactorial.Mul(nFactorial, NewBigFloat(float64(i), workPrec))
	}

	result := new(BigFloat).SetPrec(workPrec).Quo(xHalfPower, nFactorial)
	term := new(BigFloat).SetPrec(workPrec).Set(result)
	xHalf2 := new(BigFloat).SetPrec(workPrec).Mul(xHalf, xHalf)

	kFactorial := NewBigFloat(1.0, workPrec)
	nPlusKFactorial := new(BigFloat).SetPrec(workPrec).Set(nFactorial)

	convThreshold := new(BigFloat).SetPrec(workPrec).SetMantExp(NewBigFloat(1.0, workPrec), -int(workPrec+10))
	resultAbs := new(BigFloat).SetPrec(workPrec).Abs(result)

	for k := 1; k < 2000; k++ {
		// term = (-1)^k * (x/2)^(n+2k) / (k! * (n+k)!)
		term.Mul(term, xHalf2)
		kFactorial.Mul(kFactorial, NewBigFloat(float64(k), workPrec))
		nPlusKFactorial.Mul(nPlusKFactorial, NewBigFloat(float64(n+k), workPrec))
		denom := new(BigFloat).SetPrec(workPrec).Mul(kFactorial, nPlusKFactorial)
		term.Quo(term, denom)

		if k%2 == 1 {
			// Odd k: subtract
			result.Sub(result, term)
		} else {
			// Even k: add
			result.Add(result, term)
		}

		// Check convergence
		termAbs := new(BigFloat).SetPrec(workPrec).Abs(term)
		if termAbs.Cmp(convThreshold) < 0 {
			// Also check relative convergence
			resultAbs.Set(result)
			resultAbs.Abs(resultAbs)
			if resultAbs.Sign() > 0 {
				relErr := new(BigFloat).SetPrec(workPrec).Quo(termAbs, resultAbs)
				if relErr.Cmp(convThreshold) < 0 {
					break
				}
			} else {
				break
			}
		}
		resultAbs.Set(result)
		resultAbs.Abs(resultAbs)
	}

	return new(BigFloat).SetPrec(targetPrec).Set(result)
}

// BigBesselY computes the Bessel function of the second kind Y_n(x)
// Uses formula: Y_n(x) = (J_n(x)*cos(nπ) - J_{-n}(x)) / sin(nπ)
// For integer n, simplifies to: Y_n(x) = (J_n(x)*cos(nπ) - (-1)^n*J_n(x)) / sin(nπ)
// Actually, for integer n: Y_n(x) = limit_{ν→n} (J_ν(x)*cos(νπ) - J_{-ν}(x)) / sin(νπ)
// For computational purposes, we use series expansion or asymptotic forms
func BigBesselY(n int, x *BigFloat, prec uint) *BigFloat {
	if prec == 0 {
		prec = x.Prec()
	}

	// Handle special cases
	if x.Sign() <= 0 {
		return NewBigFloat(math.NaN(), prec)
	}
	if x.IsInf() {
		return NewBigFloat(0.0, prec)
	}

	workPrec := prec + 32

	// For small x, use series expansion
	// For large x, use asymptotic expansion
	// For simplicity, we'll use a combination approach

	// Y_n(x) can be computed using recurrence relations or direct series
	// For integer n, we can use:
	// Y_n(x) = (2/π) * (ln(x/2) + γ) * J_n(x) - (1/π) * sum_{k=0}^{n-1} ((n-k-1)!/k!) * (x/2)^(2k-n)
	//         - (1/π) * sum_{k=0} (-1)^k * (ψ(k+1) + ψ(n+k+1)) * (x/2)^(n+2k) / (k! * (n+k)!)
	// where γ is Euler's constant and ψ is the digamma function

	// Simplified approach: use recurrence relations starting from Y_0 and Y_1
	if n == 0 {
		return bigBesselY0(x, workPrec, prec)
	}
	if n == 1 {
		return bigBesselY1(x, workPrec, prec)
	}

	// Use recurrence: Y_{n+1}(x) = (2n/x) * Y_n(x) - Y_{n-1}(x)
	y0 := bigBesselY0(x, workPrec, workPrec)
	y1 := bigBesselY1(x, workPrec, workPrec)

	ynMinus1 := y0
	yn := y1

	for i := 1; i < n; i++ {
		// Use recurrence relation: Y_{i+1} = (2i/x) * Y_i - Y_{i-1}
		twoI := NewBigFloat(float64(2*i), workPrec)
		coeff := new(BigFloat).SetPrec(workPrec).Quo(twoI, x)
		term := new(BigFloat).SetPrec(workPrec).Mul(coeff, yn)
		ynPlus1 := new(BigFloat).SetPrec(workPrec).Sub(term, ynMinus1)

		ynMinus1 = yn
		yn = ynPlus1
	}

	return new(BigFloat).SetPrec(prec).Set(yn)
}

// bigBesselY0 computes Y_0(x) using series expansion
func bigBesselY0(x *BigFloat, workPrec, targetPrec uint) *BigFloat {
	// Y_0(x) = (2/π) * (ln(x/2) + γ) * J_0(x) - (2/π) * sum_{k=0} (-1)^k * H_k * (x/2)^(2k) / (k!)^2
	// where H_k is the k-th harmonic number (approximation of digamma)
	twoOverPi := new(BigFloat).SetPrec(workPrec).Quo(NewBigFloat(2.0, workPrec), BigPI(workPrec))
	j0 := BigBesselJ(0, x, workPrec)

	xHalf := new(BigFloat).SetPrec(workPrec).Quo(x, NewBigFloat(2.0, workPrec))
	lnXHalf := BigLog(xHalf, workPrec)
	gamma := BigEulerGamma(workPrec)
	lnTerm := new(BigFloat).SetPrec(workPrec).Add(lnXHalf, gamma)

	firstTerm := new(BigFloat).SetPrec(workPrec).Mul(twoOverPi, lnTerm)
	firstTerm.Mul(firstTerm, j0)

	// Series term: sum_{k=0} (-1)^k * H_k * (x/2)^(2k) / (k!)^2
	// H_k ≈ ln(k) + γ + 1/(2k) - 1/(12k²) + ... for large k
	// For small k, compute H_k exactly
	xHalf2 := new(BigFloat).SetPrec(workPrec).Mul(xHalf, xHalf)
	series := NewBigFloat(0.0, workPrec)
	harmonic := NewBigFloat(0.0, workPrec) // H_0 = 0
	kFactorial := NewBigFloat(1.0, workPrec)
	term := NewBigFloat(1.0, workPrec) // (x/2)^(2k) / (k!)^2, starts at k=0

	convThreshold := new(BigFloat).SetPrec(workPrec).SetMantExp(NewBigFloat(1.0, workPrec), -int(workPrec+10))

	for k := 0; k < 1000; k++ {
		if k > 0 {
			harmonic.Add(harmonic, new(BigFloat).SetPrec(workPrec).Quo(NewBigFloat(1.0, workPrec), NewBigFloat(float64(k), workPrec)))
			term.Mul(term, xHalf2)
			kFactorial.Mul(kFactorial, NewBigFloat(float64(k), workPrec))
			kFactorial2 := new(BigFloat).SetPrec(workPrec).Mul(kFactorial, kFactorial)
			term.Quo(term, kFactorial2)
		}

		seriesTerm := new(BigFloat).SetPrec(workPrec).Mul(harmonic, term)
		if k%2 == 1 {
			seriesTerm.Neg(seriesTerm)
		}
		series.Add(series, seriesTerm)

		termAbs := new(BigFloat).SetPrec(workPrec).Abs(seriesTerm)
		if termAbs.Cmp(convThreshold) < 0 {
			seriesAbs := new(BigFloat).SetPrec(workPrec).Abs(series)
			if seriesAbs.Sign() > 0 {
				relErr := new(BigFloat).SetPrec(workPrec).Quo(termAbs, seriesAbs)
				if relErr.Cmp(convThreshold) < 0 {
					break
				}
			} else {
				break
			}
		}
	}

	series.Mul(series, twoOverPi)
	result := new(BigFloat).SetPrec(workPrec).Sub(firstTerm, series)

	return new(BigFloat).SetPrec(targetPrec).Set(result)
}

// bigBesselY1 computes Y_1(x) using series expansion
func bigBesselY1(x *BigFloat, workPrec, targetPrec uint) *BigFloat {
	// Y_1(x) = (2/π) * (ln(x/2) + γ) * J_1(x) - (2/π) * (1/x) - (2/π) * sum_{k=0} (-1)^k * (H_k + H_{k+1}) * (x/2)^(2k+1) / (k! * (k+1)!)
	twoOverPi := new(BigFloat).SetPrec(workPrec).Quo(NewBigFloat(2.0, workPrec), BigPI(workPrec))
	j1 := BigBesselJ(1, x, workPrec)

	xHalf := new(BigFloat).SetPrec(workPrec).Quo(x, NewBigFloat(2.0, workPrec))
	lnXHalf := BigLog(xHalf, workPrec)
	gamma := BigEulerGamma(workPrec)
	lnTerm := new(BigFloat).SetPrec(workPrec).Add(lnXHalf, gamma)

	firstTerm := new(BigFloat).SetPrec(workPrec).Mul(twoOverPi, lnTerm)
	firstTerm.Mul(firstTerm, j1)

	// Add -2/(π*x) term
	negTwoOverPiX := new(BigFloat).SetPrec(workPrec).Quo(NewBigFloat(-2.0, workPrec), new(BigFloat).SetPrec(workPrec).Mul(BigPI(workPrec), x))
	result := new(BigFloat).SetPrec(workPrec).Add(firstTerm, negTwoOverPiX)

	// Series term: sum_{k=0} (-1)^k * (H_k + H_{k+1}) * (x/2)^(2k+1) / (k! * (k+1)!)
	xHalf2 := new(BigFloat).SetPrec(workPrec).Mul(xHalf, xHalf)
	series := NewBigFloat(0.0, workPrec)
	harmonic := NewBigFloat(0.0, workPrec) // H_0 = 0
	kFactorial := NewBigFloat(1.0, workPrec)
	kPlus1Factorial := NewBigFloat(1.0, workPrec)      // 1! = 1
	term := new(BigFloat).SetPrec(workPrec).Set(xHalf) // (x/2)^(2k+1) / (k! * (k+1)!), starts at k=0

	convThreshold := new(BigFloat).SetPrec(workPrec).SetMantExp(NewBigFloat(1.0, workPrec), -int(workPrec+10))

	for k := 0; k < 1000; k++ {
		if k > 0 {
			harmonic.Add(harmonic, new(BigFloat).SetPrec(workPrec).Quo(NewBigFloat(1.0, workPrec), NewBigFloat(float64(k), workPrec)))
			term.Mul(term, xHalf2)
			kFactorial.Mul(kFactorial, NewBigFloat(float64(k), workPrec))
			kPlus1Factorial.Mul(kPlus1Factorial, NewBigFloat(float64(k+1), workPrec))
			term.Quo(term, new(BigFloat).SetPrec(workPrec).Mul(kFactorial, kPlus1Factorial))
		}

		harmonicKPlus1 := new(BigFloat).SetPrec(workPrec).Add(harmonic,
			new(BigFloat).SetPrec(workPrec).Quo(NewBigFloat(1.0, workPrec), NewBigFloat(float64(k+1), workPrec)))
		harmonicSum := new(BigFloat).SetPrec(workPrec).Add(harmonic, harmonicKPlus1)

		seriesTerm := new(BigFloat).SetPrec(workPrec).Mul(harmonicSum, term)
		if k%2 == 1 {
			seriesTerm.Neg(seriesTerm)
		}
		series.Add(series, seriesTerm)

		termAbs := new(BigFloat).SetPrec(workPrec).Abs(seriesTerm)
		if termAbs.Cmp(convThreshold) < 0 {
			seriesAbs := new(BigFloat).SetPrec(workPrec).Abs(series)
			if seriesAbs.Sign() > 0 {
				relErr := new(BigFloat).SetPrec(workPrec).Quo(termAbs, seriesAbs)
				if relErr.Cmp(convThreshold) < 0 {
					break
				}
			} else {
				break
			}
		}
	}

	series.Mul(series, twoOverPi)
	result.Sub(result, series)

	return new(BigFloat).SetPrec(targetPrec).Set(result)
}
