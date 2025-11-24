//go:build !amd64 && !arm64

package bigmath

// Generic fallback implementations for architectures without assembly support

// mpnAddNDualCarry uses dual carry chains (ADCX/ADOX) for parallel carry propagation
// Requires BMI2 support (Intel Broadwell+, AMD Excavator+)
// Returns carry (0 or 1)
// Generic fallback: uses regular mpnAddN
func mpnAddNDualCarry(dst, src1, src2 *uint64, n int) uint64 {
	return mpnAddN(dst, src1, src2, n)
}

// mpnFMA computes fused multiply-add: dst = multiplier * src + addend
// This is useful for patterns like: result = scale * a + b
// Returns carry (high limb)
// Generic fallback implementation
func mpnFMA(dst, src, addend *uint64, n int, multiplier uint64) uint64 {
	tmp := make([]uint64, n+1)
	high := mpnMul1(&tmp[0], src, n, multiplier)
	carry := mpnAddN(dst, &tmp[0], addend, n)
	if carry != 0 {
		if high != ^uint64(0) {
			high++
		} else {
			return 1
		}
	}
	return high
}

