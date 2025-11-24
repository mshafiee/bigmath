// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

//go:build amd64

package bigmath

// Forward declarations for AMD64 assembly implementations
// These functions are implemented in assembly files

// mpnAddNDualCarry uses dual carry chains (ADCX/ADOX) for parallel carry propagation
// Requires BMI2 support (Intel Broadwell+, AMD Excavator+)
// Returns carry (0 or 1)
//
//nolint:unused // Implemented in assembly (mpn_add_amd64.s)
func mpnAddNDualCarry(dst, src1, src2 *uint64, n int) uint64

// mpnFMA computes fused multiply-add: dst = multiplier * src + addend
// This is useful for patterns like: result = scale * a + b
// Returns carry (high limb)
//
//nolint:unused // Implemented in assembly (mpn_fma_amd64.s)
func mpnFMA(dst, src, addend *uint64, n int, multiplier uint64) uint64
