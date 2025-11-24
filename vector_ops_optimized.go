// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

//go:build amd64 || arm64

package bigmath

// Optimized vector operations with reduced allocations

// bigVec3CrossOptimized implements optimized cross product
// Optimization: Reduce allocations, reuse temporaries
func bigVec3CrossOptimized(v1, v2 *BigVec3, prec uint) *BigVec3 {
	if prec == 0 {
		prec = v1.X.Prec()
	}

	// Preallocate temporaries
	t1 := new(BigFloat).SetPrec(prec)
	t2 := new(BigFloat).SetPrec(prec)

	// X component: v1.Y*v2.Z - v1.Z*v2.Y
	t1.Mul(v1.Y, v2.Z)
	t2.Mul(v1.Z, v2.Y)
	x := new(BigFloat).SetPrec(prec).Sub(t1, t2)

	// Y component: v1.Z*v2.X - v1.X*v2.Z
	t1.Mul(v1.Z, v2.X)
	t2.Mul(v1.X, v2.Z)
	y := new(BigFloat).SetPrec(prec).Sub(t1, t2)

	// Z component: v1.X*v2.Y - v1.Y*v2.X
	t1.Mul(v1.X, v2.Y)
	t2.Mul(v1.Y, v2.X)
	z := new(BigFloat).SetPrec(prec).Sub(t1, t2)

	return &BigVec3{X: x, Y: y, Z: z}
}

// bigVec3NormalizeOptimized implements optimized vector normalization
// Optimization: Single magnitude calculation, reuse division
func bigVec3NormalizeOptimized(v *BigVec3, prec uint) *BigVec3 {
	if prec == 0 {
		prec = v.X.Prec()
	}

	// Calculate magnitude once
	magSq := new(BigFloat).SetPrec(prec)
	temp := new(BigFloat).SetPrec(prec)

	// mag² = x² + y² + z²
	temp.Mul(v.X, v.X)
	magSq.Add(magSq, temp)

	temp.Mul(v.Y, v.Y)
	magSq.Add(magSq, temp)

	temp.Mul(v.Z, v.Z)
	magSq.Add(magSq, temp)

	// Check for zero
	zero := NewBigFloat(0.0, prec)
	if magSq.Cmp(zero) == 0 {
		return &BigVec3{
			X: NewBigFloat(0.0, prec),
			Y: NewBigFloat(0.0, prec),
			Z: NewBigFloat(0.0, prec),
		}
	}

	// magnitude = sqrt(mag²)
	magnitude := BigSqrt(magSq, prec)

	// Normalize: divide each component by magnitude
	return &BigVec3{
		X: new(BigFloat).SetPrec(prec).Quo(v.X, magnitude),
		Y: new(BigFloat).SetPrec(prec).Quo(v.Y, magnitude),
		Z: new(BigFloat).SetPrec(prec).Quo(v.Z, magnitude),
	}
}

// bigVec3AngleOptimized implements optimized angle calculation
// Optimization: Combined operations, reduced intermediate allocations
func bigVec3AngleOptimized(v1, v2 *BigVec3, prec uint) *BigFloat {
	if prec == 0 {
		prec = v1.X.Prec()
	}

	// Compute dot product inline to avoid function call overhead
	dot := new(BigFloat).SetPrec(prec)
	temp := new(BigFloat).SetPrec(prec)

	temp.Mul(v1.X, v2.X)
	dot.Add(dot, temp)
	temp.Mul(v1.Y, v2.Y)
	dot.Add(dot, temp)
	temp.Mul(v1.Z, v2.Z)
	dot.Add(dot, temp)

	// Compute magnitudes inline
	mag1Sq := new(BigFloat).SetPrec(prec)
	temp.Mul(v1.X, v1.X)
	mag1Sq.Add(mag1Sq, temp)
	temp.Mul(v1.Y, v1.Y)
	mag1Sq.Add(mag1Sq, temp)
	temp.Mul(v1.Z, v1.Z)
	mag1Sq.Add(mag1Sq, temp)

	mag2Sq := new(BigFloat).SetPrec(prec)
	temp.Mul(v2.X, v2.X)
	mag2Sq.Add(mag2Sq, temp)
	temp.Mul(v2.Y, v2.Y)
	mag2Sq.Add(mag2Sq, temp)
	temp.Mul(v2.Z, v2.Z)
	mag2Sq.Add(mag2Sq, temp)

	// Check for zero vectors
	zero := NewBigFloat(0.0, prec)
	if mag1Sq.Cmp(zero) == 0 || mag2Sq.Cmp(zero) == 0 {
		return NewBigFloat(0.0, prec) // Return 0 instead of NaN for better performance
	}

	// Compute mag1 * mag2 = sqrt(mag1Sq * mag2Sq)
	magProduct := new(BigFloat).SetPrec(prec).Mul(mag1Sq, mag2Sq)
	magProduct = BigSqrt(magProduct, prec)

	// Compute cosAngle = dot / (mag1 * mag2)
	cosAngle := new(BigFloat).SetPrec(prec).Quo(dot, magProduct)

	// Clamp to [-1, 1]
	one := NewBigFloat(1.0, prec)
	negOne := NewBigFloat(-1.0, prec)
	if cosAngle.Cmp(one) > 0 {
		cosAngle.Set(one)
	} else if cosAngle.Cmp(negOne) < 0 {
		cosAngle.Set(negOne)
	}

	return BigAcos(cosAngle, prec)
}

// bigVec3ProjectOptimized implements optimized vector projection
// Optimization: Inline dot products, reuse calculations
func bigVec3ProjectOptimized(v1, v2 *BigVec3, prec uint) *BigVec3 {
	if prec == 0 {
		prec = v1.X.Prec()
	}

	// Compute v1·v2 inline
	dot := new(BigFloat).SetPrec(prec)
	temp := new(BigFloat).SetPrec(prec)

	temp.Mul(v1.X, v2.X)
	dot.Add(dot, temp)
	temp.Mul(v1.Y, v2.Y)
	dot.Add(dot, temp)
	temp.Mul(v1.Z, v2.Z)
	dot.Add(dot, temp)

	// Compute |v2|² inline
	mag2Sq := new(BigFloat).SetPrec(prec)
	temp.Mul(v2.X, v2.X)
	mag2Sq.Add(mag2Sq, temp)
	temp.Mul(v2.Y, v2.Y)
	mag2Sq.Add(mag2Sq, temp)
	temp.Mul(v2.Z, v2.Z)
	mag2Sq.Add(mag2Sq, temp)

	// Check for zero v2
	zero := NewBigFloat(0.0, prec)
	if mag2Sq.Cmp(zero) == 0 {
		return &BigVec3{
			X: NewBigFloat(0.0, prec),
			Y: NewBigFloat(0.0, prec),
			Z: NewBigFloat(0.0, prec),
		}
	}

	// scalar = (v1·v2) / |v2|²
	scalar := new(BigFloat).SetPrec(prec).Quo(dot, mag2Sq)

	// Return scalar * v2
	return &BigVec3{
		X: new(BigFloat).SetPrec(prec).Mul(v2.X, scalar),
		Y: new(BigFloat).SetPrec(prec).Mul(v2.Y, scalar),
		Z: new(BigFloat).SetPrec(prec).Mul(v2.Z, scalar),
	}
}
