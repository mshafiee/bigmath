// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

//go:build amd64

package bigmath

// detectAMD64Features detects AMD64-specific CPU features
func detectAMD64Features(features *CPUFeatures) {
	// CPUID is always available on AMD64
	features.HasAVX = cpuidAVX()
	features.HasAVX2 = cpuidAVX2()
	features.HasAVX512 = cpuidAVX512()
	features.HasFMA = cpuidFMA()
	features.HasBMI2 = cpuidBMI2()
	// x87 FPU is always available on x86-64 (AMD64)
	features.HasX87 = true
}

// detectARM64Features is not applicable on AMD64
func detectARM64Features(features *CPUFeatures) {
	// No ARM64 features on AMD64
}

// cpuidAVX checks for AVX support using CPUID instruction
// CPUID EAX=1: ECX bit 28 indicates AVX support
//
//go:noescape
func cpuidAVX() bool

// cpuidAVX2 checks for AVX2 support using CPUID instruction
// CPUID EAX=7, ECX=0: EBX bit 5 indicates AVX2 support
//
//go:noescape
func cpuidAVX2() bool

// cpuidAVX512 checks for AVX-512 support using CPUID instruction
// CPUID EAX=7, ECX=0: EBX bit 16 indicates AVX-512F support
//
//go:noescape
func cpuidAVX512() bool

// cpuidFMA checks for FMA support using CPUID instruction
// CPUID EAX=1: ECX bit 12 indicates FMA support
//
//go:noescape
func cpuidFMA() bool

// cpuidBMI2 checks for BMI2 support using CPUID instruction
// CPUID EAX=7, ECX=0: EBX bit 8 indicates BMI2 support
//
//go:noescape
func cpuidBMI2() bool
