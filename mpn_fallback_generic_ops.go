// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

//go:build !amd64 && !arm64

package bigmath

import (
	"math/bits"
	"unsafe"
)

// Generic fallback implementations for mpn operations on platforms without assembly support
// These functions use unsafe pointer arithmetic to match the assembly function signatures

// mpnAddN adds two n-limb numbers: dst = src1 + src2
// Returns carry (0 or 1)
func mpnAddN(dst, src1, src2 *uint64, n int) uint64 {
	dstSlice := unsafe.Slice(dst, n)
	src1Slice := unsafe.Slice(src1, n)
	src2Slice := unsafe.Slice(src2, n)
	var carry uint64
	for i := 0; i < n; i++ {
		sum := src1Slice[i] + src2Slice[i] + carry
		dstSlice[i] = sum
		if sum < src1Slice[i] || (carry != 0 && sum == src1Slice[i]) {
			carry = 1
		} else {
			carry = 0
		}
	}
	return carry
}

// mpnSubN subtracts two n-limb numbers: dst = src1 - src2
// Returns borrow (0 or 1)
func mpnSubN(dst, src1, src2 *uint64, n int) uint64 {
	dstSlice := unsafe.Slice(dst, n)
	src1Slice := unsafe.Slice(src1, n)
	src2Slice := unsafe.Slice(src2, n)
	var borrow uint64
	for i := 0; i < n; i++ {
		diff := src1Slice[i] - src2Slice[i] - borrow
		dstSlice[i] = diff
		if diff > src1Slice[i] || (borrow != 0 && diff == src1Slice[i]) {
			borrow = 1
		} else {
			borrow = 0
		}
	}
	return borrow
}

// mpnMul1 multiplies n-limb number by single limb: dst = src * multiplier
// Returns high limb of result
func mpnMul1(dst, src *uint64, n int, multiplier uint64) uint64 {
	dstSlice := unsafe.Slice(dst, n)
	srcSlice := unsafe.Slice(src, n)
	var hi uint64
	for i := 0; i < n; i++ {
		lo, carry := bits.Mul64(srcSlice[i], multiplier)
		lo += hi
		if lo < hi {
			carry++
		}
		dstSlice[i] = lo
		hi = carry
	}
	return hi
}

// mpnAddMul1 multiplies and accumulates: dst += src * multiplier
// Returns carry (high limb)
func mpnAddMul1(dst, src *uint64, n int, multiplier uint64) uint64 {
	dstSlice := unsafe.Slice(dst, n)
	srcSlice := unsafe.Slice(src, n)
	var hi uint64
	for i := 0; i < n; i++ {
		lo, carry := bits.Mul64(srcSlice[i], multiplier)
		lo += hi
		if lo < hi {
			carry++
		}
		lo += dstSlice[i]
		if lo < dstSlice[i] {
			carry++
		}
		dstSlice[i] = lo
		hi = carry
	}
	return hi
}

// mpnLShift left shifts n-limb number by count bits
// Returns bits shifted out
func mpnLShift(dst, src *uint64, n int, count uint) uint64 {
	if n == 0 || count == 0 {
		return 0
	}
	dstSlice := unsafe.Slice(dst, n)
	srcSlice := unsafe.Slice(src, n)
	if count >= 64 {
		// Shift by multiples of 64
		shiftWords := int(count / 64)
		if shiftWords >= n {
			// All bits shifted out
			for i := 0; i < n; i++ {
				dstSlice[i] = 0
			}
			return 0
		}
		// Shift remaining bits
		count = count % 64
		for i := 0; i < n-shiftWords; i++ {
			dstSlice[i] = srcSlice[i+shiftWords]
		}
		for i := n - shiftWords; i < n; i++ {
			dstSlice[i] = 0
		}
		if count == 0 {
			return 0
		}
		srcSlice = dstSlice
	}
	var carry uint64
	for i := 0; i < n; i++ {
		val := srcSlice[i]
		dstSlice[i] = (val << count) | carry
		carry = val >> (64 - count)
	}
	return carry
}

// mpnRShift right shifts n-limb number by count bits
// Returns bits shifted out
func mpnRShift(dst, src *uint64, n int, count uint) uint64 {
	if n == 0 || count == 0 {
		return 0
	}
	dstSlice := unsafe.Slice(dst, n)
	srcSlice := unsafe.Slice(src, n)
	if count >= 64 {
		// Shift by multiples of 64
		shiftWords := int(count / 64)
		if shiftWords >= n {
			// All bits shifted out
			return srcSlice[n-1] >> (count % 64)
		}
		// Shift remaining bits
		count = count % 64
		var carry uint64
		if shiftWords > 0 {
			carry = srcSlice[shiftWords-1] >> (64 - count)
		}
		for i := 0; i < n-shiftWords; i++ {
			val := srcSlice[i+shiftWords]
			dstSlice[i] = (val >> count) | carry
			carry = val << (64 - count)
		}
		return carry
	}
	var carry uint64
	for i := n - 1; i >= 0; i-- {
		val := srcSlice[i]
		dstSlice[i] = (val >> count) | carry
		carry = val << (64 - count)
	}
	return carry >> (64 - count)
}

