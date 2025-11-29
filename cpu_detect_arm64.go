// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

//go:build arm64

package bigmath

// detectAMD64Features is not applicable on ARM64
func detectAMD64Features(features *CPUFeatures) {
	// No AMD64 features on ARM64
	features.HasX87 = false
}

// detectARM64Features detects ARM64-specific CPU features
func detectARM64Features(features *CPUFeatures) {
	// NEON is standard on ARMv8
	features.HasNEON = true

	// SVE detection requires reading ID_AA64PFR0_EL1 register
	// For now, we'll conservatively set to false
	// In a production system, this would check hwcap or similar
	features.HasSVE = false
}
