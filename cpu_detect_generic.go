// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

//go:build !amd64 && !arm64

package bigmath

// detectAMD64Features is not applicable for non-AMD64 platforms
func detectAMD64Features(features *CPUFeatures) {
	// No AMD64 features on non-AMD64 platforms
	features.HasX87 = false
}

// detectARM64Features is not applicable for non-ARM64 platforms
func detectARM64Features(features *CPUFeatures) {
	// No ARM64 features on non-ARM64 platforms
}
