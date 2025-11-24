// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

// Multi-precision integer operations
// These are low-level building blocks for arbitrary precision arithmetic

// mpnAddN adds two n-limb numbers: dst = src1 + src2
// Returns carry (0 or 1)
func mpnAddN(dst, src1, src2 *uint64, n int) uint64

// mpnSubN subtracts two n-limb numbers: dst = src1 - src2
// Returns borrow (0 or 1)
func mpnSubN(dst, src1, src2 *uint64, n int) uint64

// mpnMul1 multiplies n-limb number by single limb: dst = src * multiplier
// Returns high limb of result
func mpnMul1(dst, src *uint64, n int, multiplier uint64) uint64

// mpnAddMul1 multiplies and accumulates: dst += src * multiplier
// Returns carry (high limb)
func mpnAddMul1(dst, src *uint64, n int, multiplier uint64) uint64

// mpnLShift left shifts n-limb number by count bits
// Returns bits shifted out
func mpnLShift(dst, src *uint64, n int, count uint) uint64

// mpnRShift right shifts n-limb number by count bits
// Returns bits shifted out
func mpnRShift(dst, src *uint64, n int, count uint) uint64

// mpnAddNDualCarry and mpnFMA are declared in their respective implementation files
// to avoid conflicts between assembly and Go implementations.
// See: mpn_add_amd64.s, mpn_fma_amd64.s, mpn_fma_arm64.s, mpn_fallback_arm64.go, mpn_fallback_generic.go

