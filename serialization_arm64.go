// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

//go:build arm64

package bigmath

import (
	"fmt"
	"io"
	"math/big"
)

// Assembly function declarations
//
//go:noescape
func extractIEEE754ComponentsARM64(bits uint64) (sign uint64, exponent int64, mantissa uint64)

//go:noescape
func convertEndiannessBytesARM64(bytes *[8]byte, bigEndian uint8) uint64

// readDoubleAsBigFloatAsm is the assembly-optimized version of ReadDoubleAsBigFloat
func readDoubleAsBigFloatAsm(r io.Reader, bigEndian bool, prec uint) (*BigFloat, error) {
	if prec == 0 {
		prec = DefaultPrecision
	}

	// Read 8 bytes
	var doubleBytes [8]byte
	if _, err := io.ReadFull(r, doubleBytes[:]); err != nil {
		return nil, fmt.Errorf("failed to read 8 bytes: %w", err)
	}

	// Convert endianness using assembly
	var bigEndianFlag uint8
	if bigEndian {
		bigEndianFlag = 1
	}
	bits := convertEndiannessBytesARM64(&doubleBytes, bigEndianFlag)

	// Extract IEEE 754 components using assembly
	signUint, exponentInt, mantissaUint := extractIEEE754ComponentsARM64(bits)
	sign := signUint != 0
	exponent := int(exponentInt)
	mantissa := mantissaUint

	// Handle special cases
	if exponent == 0 {
		if mantissa == 0 {
			// Zero (positive or negative)
			result := new(big.Float).SetPrec(prec)
			if sign {
				result.Neg(result)
			}
			return result, nil
		}
		// Denormalized number (subnormal)
		// For denormalized: value = (-1)^sign * 2^(-1022) * (mantissa / 2^52)
		// This is a very small number, handle as zero for now
		// TODO: Implement denormalized number handling if needed
		result := new(big.Float).SetPrec(prec)
		return result, nil
	}

	if exponent == 0x7FF {
		// Infinity or NaN
		if mantissa == 0 {
			// Infinity
			result := new(big.Float).SetPrec(prec)
			result.SetInf(sign)
			return result, nil
		}
		// NaN
		result := new(big.Float).SetPrec(prec)
		// big.Float doesn't have NaN, so we'll return zero
		// Caller should check for NaN if needed
		return result, nil
	}

	// Normalized number
	// Value = (-1)^sign * 2^(exponent - 1023) * (1 + mantissa / 2^52)

	// Construct mantissa as BigFloat: 1 + mantissa / 2^52
	// This gives us the full 53-bit precision (1 implicit + 52 explicit)
	mantissaBig := new(big.Float).SetPrec(prec)
	mantissaBig.SetUint64(mantissa)

	// Divide by 2^52 to get fractional part
	two52 := new(big.Float).SetPrec(prec)
	two52.SetUint64(1 << 52) // 2^52
	mantissaBig.Quo(mantissaBig, two52)

	// Add 1 (implicit leading bit)
	one := new(big.Float).SetPrec(prec)
	one.SetUint64(1)
	mantissaBig.Add(mantissaBig, one)

	// Calculate exponent: 2^(exponent - 1023)
	expValue := exponent - 1023

	// Construct result: mantissa * 2^expValue
	// Use SetMantExp to set mantissa and exponent together
	// This is more efficient and handles large exponents correctly
	// mantissaBig is in range [1, 2), so we extract it and add expValue
	mant := new(big.Float).SetPrec(prec)
	mantExp := mantissaBig.MantExp(mant) // Extract mantissa to [0.5, 1), mantExp is 1

	// Set result = mant * 2^(expValue + mantExp)
	// mantExp is 1 because mantissaBig was in [1, 2), so we add it back
	result := new(big.Float).SetPrec(prec)
	result.SetMantExp(mant, expValue+mantExp)

	// Apply sign
	if sign {
		result.Neg(result)
	}

	return result, nil
}

// readDoubleAsBigFloatImpl dispatches to the assembly-optimized version
func readDoubleAsBigFloatImpl(r io.Reader, bigEndian bool, prec uint) (*BigFloat, error) {
	return readDoubleAsBigFloatAsm(r, bigEndian, prec)
}
