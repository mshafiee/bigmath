// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

// DegNormBig normalizes an angle in degrees to [0, 360) using BigFloat
func DegNormBig(deg *BigFloat, prec uint) *BigFloat {
	result := NewBigFloat(0, prec).Set(deg)
	deg360 := NewBigFloat(360.0, prec)
	zero := NewBigFloat(0, prec)
	
	// Reduce to [0, 360)
	for result.Cmp(zero) < 0 {
		result.Add(result, deg360)
	}
	for result.Cmp(deg360) >= 0 {
		result.Sub(result, deg360)
	}
	
	return result
}

// RadNormBig normalizes an angle in radians to [-π, π] using BigFloat
func RadNormBig(rad *BigFloat, prec uint) *BigFloat {
	result := NewBigFloat(0, prec).Set(rad)
	
	// π
	pi := NewBigFloat(3.14159265358979323846264338327950288419716939937510582097494459, prec)
	twoPi := NewBigFloat(0, prec).Mul(pi, NewBigFloat(2, prec))
	negPi := NewBigFloat(0, prec).Neg(pi)
	
	// Reduce to [-π, π]
	for result.Cmp(negPi) < 0 {
		result.Add(result, twoPi)
	}
	for result.Cmp(pi) > 0 {
		result.Sub(result, twoPi)
	}
	
	return result
}

// RadNorm02PiBig normalizes an angle in radians to [0, 2π) using BigFloat
func RadNorm02PiBig(rad *BigFloat, prec uint) *BigFloat {
	result := NewBigFloat(0, prec).Set(rad)
	
	// 2π
	pi := NewBigFloat(3.14159265358979323846264338327950288419716939937510582097494459, prec)
	twoPi := NewBigFloat(0, prec).Mul(pi, NewBigFloat(2, prec))
	zero := NewBigFloat(0, prec)
	
	// Reduce to [0, 2π)
	for result.Cmp(zero) < 0 {
		result.Add(result, twoPi)
	}
	for result.Cmp(twoPi) >= 0 {
		result.Sub(result, twoPi)
	}
	
	return result
}

