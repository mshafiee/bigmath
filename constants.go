// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

// Constants exported from assembly data sections
// These are high-precision precomputed values
// Note: bigPI, bigTwoPI, bigHalfPI are already declared in bigmath.go as *BigFloat
// These assembly constants can be used for low-level operations if needed

// LoadConstants loads constants from assembly data sections
// This is called during package initialization
//
//nolint:unused // May be called from assembly or init functions
func loadConstants() {
	// Constants are loaded from assembly data sections
	// The actual loading happens at link time
}
