// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

// Extended precision operation wrappers that check for extended precision mode
// and route to x87 implementations when available, otherwise fall back to BigFloat.
// These functions are used in the dispatch system to provide extended precision
// support when prec == ExtendedPrecision.

// bigSinWithExtended computes sin(x), using extended precision if available
//
//nolint:unused // Used in dispatch system
func bigSinWithExtended(x *BigFloat, prec uint, fallback bigSinFunc) *BigFloat {
	if CanUseExtendedPrecision(prec) {
		val := BigFloatToExtendedFloat(x)
		result := extendedSin(val)
		return ExtendedFloatToBigFloat(result, prec)
	}
	return fallback(x, prec)
}

// bigCosWithExtended computes cos(x), using extended precision if available
//
//nolint:unused // Used in dispatch system
func bigCosWithExtended(x *BigFloat, prec uint, fallback bigCosFunc) *BigFloat {
	if CanUseExtendedPrecision(prec) {
		val := BigFloatToExtendedFloat(x)
		result := extendedCos(val)
		return ExtendedFloatToBigFloat(result, prec)
	}
	return fallback(x, prec)
}

// bigTanWithExtended computes tan(x), using extended precision if available
//
//nolint:unused // Used in dispatch system
func bigTanWithExtended(x *BigFloat, prec uint, fallback bigTanFunc) *BigFloat {
	if CanUseExtendedPrecision(prec) {
		val := BigFloatToExtendedFloat(x)
		result := extendedTan(val)
		return ExtendedFloatToBigFloat(result, prec)
	}
	return fallback(x, prec)
}

// bigAtanWithExtended computes atan(x), using extended precision if available
//
//nolint:unused // Used in dispatch system
func bigAtanWithExtended(x *BigFloat, prec uint, fallback bigAtanFunc) *BigFloat {
	if CanUseExtendedPrecision(prec) {
		val := BigFloatToExtendedFloat(x)
		result := extendedAtan(val)
		return ExtendedFloatToBigFloat(result, prec)
	}
	return fallback(x, prec)
}

// bigAtan2WithExtended computes atan2(y, x), using extended precision if available
//
//nolint:unused // Used in dispatch system
func bigAtan2WithExtended(y, x *BigFloat, prec uint, fallback bigAtan2Func) *BigFloat {
	if CanUseExtendedPrecision(prec) {
		yVal := BigFloatToExtendedFloat(y)
		xVal := BigFloatToExtendedFloat(x)
		result := extendedAtan2(yVal, xVal)
		return ExtendedFloatToBigFloat(result, prec)
	}
	return fallback(y, x, prec)
}

// bigExpWithExtended computes exp(x), using extended precision if available
//
//nolint:unused // Used in dispatch system
func bigExpWithExtended(x *BigFloat, prec uint, fallback bigExpFunc) *BigFloat {
	if CanUseExtendedPrecision(prec) {
		val := BigFloatToExtendedFloat(x)
		result := extendedExp(val)
		return ExtendedFloatToBigFloat(result, prec)
	}
	return fallback(x, prec)
}

// bigLogWithExtended computes log(x), using extended precision if available
//
//nolint:unused // Used in dispatch system
func bigLogWithExtended(x *BigFloat, prec uint, fallback bigLogFunc) *BigFloat {
	if CanUseExtendedPrecision(prec) {
		val := BigFloatToExtendedFloat(x)
		result := extendedLog(val)
		return ExtendedFloatToBigFloat(result, prec)
	}
	return fallback(x, prec)
}

// bigPowWithExtended computes pow(x, y), using extended precision if available
//
//nolint:unused // Used in dispatch system
func bigPowWithExtended(x, y *BigFloat, prec uint, fallback bigPowFunc) *BigFloat {
	if CanUseExtendedPrecision(prec) {
		xVal := BigFloatToExtendedFloat(x)
		yVal := BigFloatToExtendedFloat(y)
		result := extendedPow(xVal, yVal)
		return ExtendedFloatToBigFloat(result, prec)
	}
	return fallback(x, y, prec)
}
