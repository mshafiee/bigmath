// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

import (
	"runtime"
	"sync"
)

// CPUFeatures holds detected CPU capabilities
type CPUFeatures struct {
	// AMD64 features
	HasAVX    bool
	HasAVX2   bool
	HasAVX512 bool
	HasFMA    bool
	HasBMI2   bool // Bit Manipulation Instructions 2 (ADCX, ADOX, MULX)
	HasX87    bool // x87 FPU support (80-bit extended precision)

	// ARM64 features
	HasNEON bool
	HasSVE  bool

	// Architecture
	IsAMD64 bool
	IsARM64 bool
}

var (
	cpuFeatures     CPUFeatures
	cpuFeaturesOnce sync.Once
)

// detectCPUFeatures performs runtime CPU feature detection
func detectCPUFeatures() CPUFeatures {
	var features CPUFeatures

	arch := runtime.GOARCH
	features.IsAMD64 = arch == "amd64"
	features.IsARM64 = arch == "arm64"

	if features.IsAMD64 {
		detectAMD64Features(&features)
	} else if features.IsARM64 {
		detectARM64Features(&features)
	}

	return features
}

// GetCPUFeatures returns the detected CPU features (cached)
func GetCPUFeatures() CPUFeatures {
	cpuFeaturesOnce.Do(func() {
		cpuFeatures = detectCPUFeatures()
	})
	return cpuFeatures
}

// detectAMD64Features and detectARM64Features are implemented in
// architecture-specific files: cpu_detect_amd64.go, cpu_detect_arm64.go, cpu_detect_generic.go
