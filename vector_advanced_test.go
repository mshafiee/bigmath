// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

import (
	"math"
	"testing"
)

func TestBigVec3Cross(t *testing.T) {
	prec := uint(256)

	tests := []struct {
		name     string
		v1, v2   [3]float64
		expected [3]float64
		tolerance float64
	}{
		{"unit_x_unit_y", [3]float64{1.0, 0.0, 0.0}, [3]float64{0.0, 1.0, 0.0}, [3]float64{0.0, 0.0, 1.0}, 1e-10},
		{"unit_y_unit_z", [3]float64{0.0, 1.0, 0.0}, [3]float64{0.0, 0.0, 1.0}, [3]float64{1.0, 0.0, 0.0}, 1e-10},
		{"unit_z_unit_x", [3]float64{0.0, 0.0, 1.0}, [3]float64{1.0, 0.0, 0.0}, [3]float64{0.0, 1.0, 0.0}, 1e-10},
		{"parallel_vectors", [3]float64{1.0, 2.0, 3.0}, [3]float64{2.0, 4.0, 6.0}, [3]float64{0.0, 0.0, 0.0}, 1e-10},
		{"zero_vector", [3]float64{0.0, 0.0, 0.0}, [3]float64{1.0, 2.0, 3.0}, [3]float64{0.0, 0.0, 0.0}, 1e-10},
		{"both_zero", [3]float64{0.0, 0.0, 0.0}, [3]float64{0.0, 0.0, 0.0}, [3]float64{0.0, 0.0, 0.0}, 1e-10},
		{"general_case", [3]float64{1.0, 2.0, 3.0}, [3]float64{4.0, 5.0, 6.0}, [3]float64{-3.0, 6.0, -3.0}, 1e-10},
		{"negative_components", [3]float64{-1.0, -2.0, -3.0}, [3]float64{1.0, 2.0, 3.0}, [3]float64{0.0, 0.0, 0.0}, 1e-10},
		{"very_large", [3]float64{1e10, 2e10, 3e10}, [3]float64{4e10, 5e10, 6e10}, [3]float64{-3e20, 6e20, -3e20}, 1e5},
		{"very_small", [3]float64{1e-10, 2e-10, 3e-10}, [3]float64{4e-10, 5e-10, 6e-10}, [3]float64{-3e-20, 6e-20, -3e-20}, 1e-30},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v1 := NewBigVec3(tt.v1[0], tt.v1[1], tt.v1[2], prec)
			v2 := NewBigVec3(tt.v2[0], tt.v2[1], tt.v2[2], prec)
			result := BigVec3Cross(v1, v2, prec)
			got := result.ToFloat64()

			for i := 0; i < 3; i++ {
				if math.Abs(got[i]-tt.expected[i]) > tt.tolerance {
					t.Errorf("BigVec3Cross component[%d] = %g, want %g (tolerance %g)", i, got[i], tt.expected[i], tt.tolerance)
				}
			}
			
			// Property: cross(v1, v2) = -cross(v2, v1)
			result2 := BigVec3Cross(v2, v1, prec)
			got2 := result2.ToFloat64()
			for i := 0; i < 3; i++ {
				if math.Abs(got[i]+got2[i]) > tt.tolerance {
					t.Errorf("Property violated: cross(v1,v2)[%d] = %g, but -cross(v2,v1)[%d] = %g", i, got[i], i, -got2[i])
				}
			}
			
			// Property: cross(v1, v2) is perpendicular to both v1 and v2
			if tt.name != "zero_vector" && tt.name != "both_zero" && tt.name != "parallel_vectors" {
				dot1 := BigVec3Dot(result, v1, prec)
				dot2 := BigVec3Dot(result, v2, prec)
				dot1Val, _ := dot1.Float64()
				dot2Val, _ := dot2.Float64()
				if math.Abs(dot1Val) > tt.tolerance*10 {
					t.Errorf("Property violated: cross(v1,v2) not perpendicular to v1, dot = %g", dot1Val)
				}
				if math.Abs(dot2Val) > tt.tolerance*10 {
					t.Errorf("Property violated: cross(v1,v2) not perpendicular to v2, dot = %g", dot2Val)
				}
			}
		})
	}
	
	// Test precision levels
	t.Run("precision_levels", func(t *testing.T) {
		testCases := []uint{64, 128, 256, 512}
		for _, p := range testCases {
			v1 := NewBigVec3(1.0, 0.0, 0.0, p)
			v2 := NewBigVec3(0.0, 1.0, 0.0, p)
			result := BigVec3Cross(v1, v2, p)
			got := result.ToFloat64()
			expected := [3]float64{0.0, 0.0, 1.0}
			for i := 0; i < 3; i++ {
				if math.Abs(got[i]-expected[i]) > 1e-6 {
					t.Errorf("BigVec3Cross at prec %d component[%d] = %g, want %g", p, i, got[i], expected[i])
				}
			}
		}
	})
}

func TestBigVec3Normalize(t *testing.T) {
	prec := uint(256)

	tests := []struct {
		name      string
		v         [3]float64
		shouldNaN bool
		tolerance float64
	}{
		{"unit_vector", [3]float64{1.0, 0.0, 0.0}, false, 1e-10},
		{"pythagorean_triple", [3]float64{3.0, 4.0, 0.0}, false, 1e-10},
		{"general_case", [3]float64{1.0, 2.0, 3.0}, false, 1e-10},
		{"negative_components", [3]float64{-1.0, -2.0, -3.0}, false, 1e-10},
		{"very_large", [3]float64{1e10, 2e10, 3e10}, false, 1e-8},
		{"very_small", [3]float64{1e-10, 2e-10, 3e-10}, false, 1e-8},
		{"zero_vector", [3]float64{0.0, 0.0, 0.0}, true, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewBigVec3(tt.v[0], tt.v[1], tt.v[2], prec)
			normalized := BigVec3Normalize(v, prec)
			
			if tt.shouldNaN {
				// Zero vector normalization should return zero vector or NaN
				mag := BigVec3Magnitude(normalized, prec)
				magVal, _ := mag.Float64()
				if magVal != 0.0 && !math.IsNaN(magVal) {
					t.Errorf("BigVec3Normalize(zero) should return zero or NaN, got magnitude %g", magVal)
				}
				return
			}
			
			magnitude := BigVec3Magnitude(normalized, prec)
			mag, _ := magnitude.Float64()

			if math.Abs(mag-1.0) > tt.tolerance {
				t.Errorf("BigVec3Normalize magnitude = %g, want 1.0 (tolerance %g)", mag, tt.tolerance)
			}
			
			// Property: normalized vector should be in same direction
			if tt.v[0] != 0 || tt.v[1] != 0 || tt.v[2] != 0 {
				dot := BigVec3Dot(v, normalized, prec)
				dotVal, _ := dot.Float64()
				origMag := BigVec3Magnitude(v, prec)
				origMagVal, _ := origMag.Float64()
				expectedDot := origMagVal
				if math.Abs(dotVal-expectedDot) > tt.tolerance*origMagVal {
					t.Errorf("Property violated: dot(v, normalized) = %g, want %g", dotVal, expectedDot)
				}
			}
		})
	}
	
	// Test precision levels
	t.Run("precision_levels", func(t *testing.T) {
		testCases := []uint{64, 128, 256, 512}
		for _, p := range testCases {
			v := NewBigVec3(3.0, 4.0, 0.0, p)
			normalized := BigVec3Normalize(v, p)
			magnitude := BigVec3Magnitude(normalized, p)
			mag, _ := magnitude.Float64()
			if math.Abs(mag-1.0) > 1e-6 {
				t.Errorf("BigVec3Normalize at prec %d magnitude = %g, want 1.0", p, mag)
			}
		}
	})
}

func TestBigVec3Angle(t *testing.T) {
	prec := uint(256)

	tests := []struct {
		name      string
		v1, v2    [3]float64
		expected  float64
		tolerance float64
		shouldNaN bool
	}{
		{"perpendicular", [3]float64{1.0, 0.0, 0.0}, [3]float64{0.0, 1.0, 0.0}, math.Pi / 2.0, 1e-8, false},
		{"parallel", [3]float64{1.0, 0.0, 0.0}, [3]float64{2.0, 0.0, 0.0}, 0.0, 1e-8, false},
		{"anti_parallel", [3]float64{1.0, 0.0, 0.0}, [3]float64{-1.0, 0.0, 0.0}, math.Pi, 1e-8, false},
		{"forty_five_degrees", [3]float64{1.0, 1.0, 0.0}, [3]float64{1.0, 0.0, 0.0}, math.Pi / 4.0, 1e-8, false},
		{"zero_vector", [3]float64{0.0, 0.0, 0.0}, [3]float64{1.0, 0.0, 0.0}, 0.0, 0, true},
		{"both_zero", [3]float64{0.0, 0.0, 0.0}, [3]float64{0.0, 0.0, 0.0}, 0.0, 0, true},
		{"general_case", [3]float64{1.0, 2.0, 3.0}, [3]float64{4.0, 5.0, 6.0}, 0.2257261285527342, 1e-6, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v1 := NewBigVec3(tt.v1[0], tt.v1[1], tt.v1[2], prec)
			v2 := NewBigVec3(tt.v2[0], tt.v2[1], tt.v2[2], prec)
			angle := BigVec3Angle(v1, v2, prec)
			angleVal, _ := angle.Float64()
			
			if tt.shouldNaN {
				// big.Float doesn't support NaN, so NewBigFloat(math.NaN()) returns 0
				// This is a known limitation - we accept 0 for NaN cases
				if !math.IsNaN(angleVal) && angleVal != 0.0 {
					t.Errorf("BigVec3Angle should return NaN or 0 (big.Float limitation) for zero vector, got %g", angleVal)
				}
				return
			}
			
			if math.Abs(angleVal-tt.expected) > tt.tolerance {
				t.Errorf("BigVec3Angle = %g, want %g (tolerance %g)", angleVal, tt.expected, tt.tolerance)
			}
			
			// Property: angle(v1, v2) = angle(v2, v1)
			angle2 := BigVec3Angle(v2, v1, prec)
			angle2Val, _ := angle2.Float64()
			if math.Abs(angleVal-angle2Val) > tt.tolerance {
				t.Errorf("Property violated: angle(v1,v2) = %g != angle(v2,v1) = %g", angleVal, angle2Val)
			}
			
			// Property: 0 <= angle <= pi
			if angleVal < 0 || angleVal > math.Pi+tt.tolerance {
				t.Errorf("Property violated: angle = %g should be in [0, pi]", angleVal)
			}
		})
	}
	
	// Test precision levels
	t.Run("precision_levels", func(t *testing.T) {
		testCases := []uint{64, 128, 256, 512}
		for _, p := range testCases {
			v1 := NewBigVec3(1.0, 0.0, 0.0, p)
			v2 := NewBigVec3(0.0, 1.0, 0.0, p)
			angle := BigVec3Angle(v1, v2, p)
			angleVal, _ := angle.Float64()
			expected := math.Pi / 2.0
			if math.Abs(angleVal-expected) > 1e-6 {
				t.Errorf("BigVec3Angle at prec %d = %g, want %g", p, angleVal, expected)
			}
		}
	})
}

func TestBigVec3Project(t *testing.T) {
	prec := uint(256)

	tests := []struct {
		name      string
		v1, v2    [3]float64
		expected  [3]float64
		tolerance float64
		shouldNaN bool
	}{
		{"onto_x_axis", [3]float64{3.0, 4.0, 0.0}, [3]float64{1.0, 0.0, 0.0}, [3]float64{3.0, 0.0, 0.0}, 1e-10, false},
		{"parallel_vectors", [3]float64{2.0, 4.0, 6.0}, [3]float64{1.0, 2.0, 3.0}, [3]float64{2.0, 4.0, 6.0}, 1e-10, false},
		{"perpendicular_vectors", [3]float64{0.0, 1.0, 0.0}, [3]float64{1.0, 0.0, 0.0}, [3]float64{0.0, 0.0, 0.0}, 1e-10, false},
		{"general_case", [3]float64{1.0, 2.0, 3.0}, [3]float64{1.0, 1.0, 0.0}, [3]float64{1.5, 1.5, 0.0}, 1e-8, false},
		{"zero_v1", [3]float64{0.0, 0.0, 0.0}, [3]float64{1.0, 0.0, 0.0}, [3]float64{0.0, 0.0, 0.0}, 1e-10, false},
		{"zero_v2", [3]float64{1.0, 2.0, 3.0}, [3]float64{0.0, 0.0, 0.0}, [3]float64{0.0, 0.0, 0.0}, 0, true},
		{"very_large", [3]float64{1e10, 2e10, 3e10}, [3]float64{1.0, 0.0, 0.0}, [3]float64{1e10, 0.0, 0.0}, 1e5, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v1 := NewBigVec3(tt.v1[0], tt.v1[1], tt.v1[2], prec)
			v2 := NewBigVec3(tt.v2[0], tt.v2[1], tt.v2[2], prec)
			projection := BigVec3Project(v1, v2, prec)
			proj := projection.ToFloat64()
			
			if tt.shouldNaN {
				// Projection onto zero vector should return zero or NaN
				mag := math.Sqrt(proj[0]*proj[0] + proj[1]*proj[1] + proj[2]*proj[2])
				if mag != 0.0 && !math.IsNaN(mag) {
					t.Errorf("BigVec3Project onto zero vector should return zero, got %v", proj)
				}
				return
			}

			for i := 0; i < 3; i++ {
				if math.Abs(proj[i]-tt.expected[i]) > tt.tolerance {
					t.Errorf("BigVec3Project component[%d] = %g, want %g (tolerance %g)", i, proj[i], tt.expected[i], tt.tolerance)
				}
			}
			
			// Property: projection should be parallel to v2
			if tt.v2[0] != 0 || tt.v2[1] != 0 || tt.v2[2] != 0 {
				projVec := NewBigVec3(proj[0], proj[1], proj[2], prec)
				cross := BigVec3Cross(projVec, v2, prec)
				crossMag := BigVec3Magnitude(cross, prec)
				crossMagVal, _ := crossMag.Float64()
				if crossMagVal > tt.tolerance*10 {
					t.Errorf("Property violated: projection not parallel to v2, cross magnitude = %g", crossMagVal)
				}
			}
		})
	}
	
	// Test precision levels
	t.Run("precision_levels", func(t *testing.T) {
		testCases := []uint{64, 128, 256, 512}
		for _, p := range testCases {
			v1 := NewBigVec3(3.0, 4.0, 0.0, p)
			v2 := NewBigVec3(1.0, 0.0, 0.0, p)
			projection := BigVec3Project(v1, v2, p)
			proj := projection.ToFloat64()
			expected := [3]float64{3.0, 0.0, 0.0}
			for i := 0; i < 3; i++ {
				if math.Abs(proj[i]-expected[i]) > 1e-6 {
					t.Errorf("BigVec3Project at prec %d component[%d] = %g, want %g", p, i, proj[i], expected[i])
				}
			}
		}
	})
}
