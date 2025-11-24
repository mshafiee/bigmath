// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

// BigVec6Add adds two BigVec6 vectors
func BigVec6Add(v1, v2 *BigVec6, prec uint) *BigVec6 {
	return getDispatcher().BigVec6AddImpl(v1, v2, prec)
}

// BigVec6Sub subtracts two BigVec6 vectors
func BigVec6Sub(v1, v2 *BigVec6, prec uint) *BigVec6 {
	return getDispatcher().BigVec6SubImpl(v1, v2, prec)
}

// BigVec6Negate negates all components of a BigVec6
func BigVec6Negate(v *BigVec6, prec uint) *BigVec6 {
	return getDispatcher().BigVec6NegateImpl(v, prec)
}

// BigVec6Magnitude computes the magnitude of the position component
func BigVec6Magnitude(v *BigVec6, prec uint) *BigFloat {
	return getDispatcher().BigVec6MagnitudeImpl(v, prec)
}

// ApplyRotationMatrixToBigVec6 applies a rotation matrix to position and velocity
func ApplyRotationMatrixToBigVec6(m *BigMatrix3x3, v *BigVec6, prec uint) *BigVec6 {
	if prec == 0 {
		prec = DefaultPrecision
	}

	// Rotate position
	pos := &BigVec3{X: v.X, Y: v.Y, Z: v.Z}
	rotPos := BigMatMul(m, pos, prec)

	// Rotate velocity
	vel := &BigVec3{X: v.VX, Y: v.VY, Z: v.VZ}
	rotVel := BigMatMul(m, vel, prec)

	return &BigVec6{
		X:  rotPos.X,
		Y:  rotPos.Y,
		Z:  rotPos.Z,
		VX: rotVel.X,
		VY: rotVel.Y,
		VZ: rotVel.Z,
	}
}

// CreateRotationMatrix creates a rotation matrix for given angles
// This is used for precession and coordinate transformations
func CreateRotationMatrix(angles [3]*BigFloat, prec uint) *BigMatrix3x3 {
	if prec == 0 {
		prec = DefaultPrecision
	}

	// Simple rotation around Z axis (first angle only for now)
	// For more complex rotations, this would be extended
	cosA := BigCos(angles[0], prec)
	sinA := BigSin(angles[0], prec)

	// Combined rotation matrix (Z-axis rotation)
	zero := NewBigFloat(0, prec)
	one := NewBigFloat(1, prec)
	negSinA := NewBigFloat(0, prec).Neg(sinA)

	return &BigMatrix3x3{
		M: [3][3]*BigFloat{
			{NewBigFloat(0, prec).Set(cosA), NewBigFloat(0, prec).Set(negSinA), NewBigFloat(0, prec).Set(zero)},
			{NewBigFloat(0, prec).Set(sinA), NewBigFloat(0, prec).Set(cosA), NewBigFloat(0, prec).Set(zero)},
			{NewBigFloat(0, prec).Set(zero), NewBigFloat(0, prec).Set(zero), NewBigFloat(0, prec).Set(one)},
		},
	}
}

// BigFloatFMA computes a*b + c using Fused Multiply-Add for higher precision
// This reduces rounding errors compared to separate multiply and add operations
// For BigFloat, we simulate FMA by using extended precision internally
func BigFloatFMA(a, b, c *BigFloat, prec uint) *BigFloat {
	if prec == 0 {
		prec = DefaultPrecision
	}

	// Use extended precision for intermediate calculation to simulate FMA behavior
	// This matches how hardware FMA works - single rounding at the end
	workPrec := prec + 64

	// Compute a*b with extended precision
	product := new(BigFloat).SetPrec(workPrec)
	product.Mul(a, b)

	// Add c with extended precision
	result := new(BigFloat).SetPrec(workPrec)
	result.Add(product, c)

	// Round to target precision (single rounding, like hardware FMA)
	return new(BigFloat).SetPrec(prec).Set(result)
}

// BigFloatDotProduct computes the dot product of two vectors using FMA for higher precision
// v1 and v2 must have the same length
// Uses FMA chain: FMA(v1[0], v2[0], FMA(v1[1], v2[1], FMA(v1[2], v2[2], 0)))
func BigFloatDotProduct(v1, v2 []*BigFloat, prec uint) *BigFloat {
	if prec == 0 {
		if len(v1) > 0 {
			prec = v1[0].Prec()
		} else {
			prec = DefaultPrecision
		}
	}

	if len(v1) != len(v2) {
		panic("BigFloatDotProduct: vectors must have same length")
	}

	if len(v1) == 0 {
		return NewBigFloat(0, prec)
	}

	// Use FMA chain for better numerical stability
	// Start with last term, then work backwards
	result := NewBigFloat(0, prec)

	for i := len(v1) - 1; i >= 0; i-- {
		result = BigFloatFMA(v1[i], v2[i], result, prec)
	}

	return result
}

// BigLog2 returns ln(2) with specified precision
func BigLog2(prec uint) *BigFloat {
	if prec == 0 {
		prec = DefaultPrecision
	}

	// Use high-precision string constant for ln(2)
	// 0.69314718055994530941723212145817656807550013436025525412068000949339362196969471
	log2Str := "0.69314718055994530941723212145817656807550013436025525412068000949339362196969471"
	result, _ := NewBigFloatFromString(log2Str, prec)
	return result
}

// BigJ2000 returns the Julian day for J2000.0 epoch (2451545.0) with specified precision
func BigJ2000(prec uint) *BigFloat {
	if prec == 0 {
		prec = DefaultPrecision
	}
	return NewBigFloat(2451545.0, prec)
}

// BigLightSpeedAUperDay returns the speed of light in AU/day (173.1446327205363) with specified precision
// This is calculated as AUNIT/CLIGHT/86400.0
// AUNIT = 1.49597870700e+11 m, CLIGHT = 2.99792458e+8 m/s
func BigLightSpeedAUperDay(prec uint) *BigFloat {
	if prec == 0 {
		prec = DefaultPrecision
	}
	return NewBigFloat(173.1446327205363, prec)
}

// BigJulianCentury returns 36525.0 (days in a Julian century) with specified precision
func BigJulianCentury(prec uint) *BigFloat {
	if prec == 0 {
		prec = DefaultPrecision
	}
	return NewBigFloat(36525.0, prec)
}

// BigJulianMillennium returns 365250.0 (days in a Julian millennium) with specified precision
func BigJulianMillennium(prec uint) *BigFloat {
	if prec == 0 {
		prec = DefaultPrecision
	}
	return NewBigFloat(365250.0, prec)
}

// BigE returns e (Euler's number) with specified precision
func BigE(prec uint) *BigFloat {
	if prec == 0 {
		prec = DefaultPrecision
	}
	// Use high-precision string constant for e
	// 2.71828182845904523536028747135266249775724709369995957496696762772407663035354759
	eStr := "2.71828182845904523536028747135266249775724709369995957496696762772407663035354759"
	result, _ := NewBigFloatFromString(eStr, prec)
	return result
}

// BigEulerGamma returns Euler's constant γ ≈ 0.57721... with specified precision
// This is a placeholder - full implementation would use a series or continued fraction
func BigEulerGamma(prec uint) *BigFloat {
	if prec == 0 {
		prec = DefaultPrecision
	}
	// Placeholder: return approximate value
	// Full implementation would compute using series: γ = lim(n→∞) (H_n - ln(n))
	gammaStr := "0.57721566490153286060651209008240243104215933593992359880576723488486772677766467"
	result, _ := NewBigFloatFromString(gammaStr, prec)
	return result
}

// BigCatalan returns Catalan's constant G ≈ 0.91596... with specified precision
// This is a placeholder - full implementation would use a series
func BigCatalan(prec uint) *BigFloat {
	if prec == 0 {
		prec = DefaultPrecision
	}
	// Placeholder: return approximate value
	// Full implementation would compute using series
	catalanStr := "0.91596559417721901505460351493238411077414937428167213426649811962176301977625476"
	result, _ := NewBigFloatFromString(catalanStr, prec)
	return result
}
