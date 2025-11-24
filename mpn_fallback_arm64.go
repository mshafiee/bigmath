// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

//go:build arm64

package bigmath

// mpnAddNDualCarry uses dual carry chains (ADCX/ADOX) for parallel carry propagation
// Requires BMI2 support (Intel Broadwell+, AMD Excavator+)
// Returns carry (0 or 1)
// ARM64 fallback: uses regular mpnAddN since ARM64 doesn't have equivalent dual carry instructions
func mpnAddNDualCarry(dst, src1, src2 *uint64, n int) uint64 {
	return mpnAddN(dst, src1, src2, n)
}

// mpnFMA computes fused multiply-add: dst = multiplier * src + addend
// This is useful for patterns like: result = scale * a + b
// Returns carry (high limb)
// Implemented in assembly: mpn_fma_arm64.s
func mpnFMA(dst, src, addend *uint64, n int, multiplier uint64) uint64

