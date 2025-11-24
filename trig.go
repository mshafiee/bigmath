package bigmath

// BigSin computes sin(x) using Taylor series with arbitrary precision
// sin(x) = x - x³/3! + x⁵/5! - x⁷/7! + x⁹/9! - ...
func BigSin(x *BigFloat, prec uint) *BigFloat {
	return getDispatcher().BigSinImpl(x, prec)
}

// BigSinRounded computes sin(x) and rounds the result according to the mode
func BigSinRounded(x *BigFloat, prec uint, mode RoundingMode) (*BigFloat, int) {
	// Compute with higher precision then round
	workPrec := prec + 32
	res := BigSin(x, workPrec)
	return Round(res, prec, mode)
}

// BigCos computes cos(x) using Taylor series with arbitrary precision
// cos(x) = 1 - x²/2! + x⁴/4! - x⁶/6! + x⁸/8! - ...
func BigCos(x *BigFloat, prec uint) *BigFloat {
	return getDispatcher().BigCosImpl(x, prec)
}

// BigCosRounded computes cos(x) and rounds the result according to the mode
func BigCosRounded(x *BigFloat, prec uint, mode RoundingMode) (*BigFloat, int) {
	workPrec := prec + 32
	res := BigCos(x, workPrec)
	return Round(res, prec, mode)
}

// BigTan computes tan(x) = sin(x) / cos(x)
func BigTan(x *BigFloat, prec uint) *BigFloat {
	return getDispatcher().BigTanImpl(x, prec)
}

// BigTanRounded computes tan(x) and rounds the result according to the mode
func BigTanRounded(x *BigFloat, prec uint, mode RoundingMode) (*BigFloat, int) {
	workPrec := prec + 32
	res := BigTan(x, workPrec)
	return Round(res, prec, mode)
}

// BigAtan computes arctan(x) using Taylor series
// atan(x) = x - x³/3 + x⁵/5 - x⁷/7 + ... for |x| ≤ 1
// For |x| > 1, use atan(x) = π/2 - atan(1/x)
func BigAtan(x *BigFloat, prec uint) *BigFloat {
	return getDispatcher().BigAtanImpl(x, prec)
}

// BigAtanRounded computes atan(x) and rounds the result according to the mode
func BigAtanRounded(x *BigFloat, prec uint, mode RoundingMode) (*BigFloat, int) {
	workPrec := prec + 32
	res := BigAtan(x, workPrec)
	return Round(res, prec, mode)
}

// BigAtan2 computes atan2(y, x) with arbitrary precision
// Returns the angle in radians between the positive x-axis and the point (x, y)
func BigAtan2(y, x *BigFloat, prec uint) *BigFloat {
	return getDispatcher().BigAtan2Impl(y, x, prec)
}

// BigAtan2Rounded computes atan2(y, x) and rounds the result according to the mode
func BigAtan2Rounded(y, x *BigFloat, prec uint, mode RoundingMode) (*BigFloat, int) {
	workPrec := prec + 32
	res := BigAtan2(y, x, workPrec)
	return Round(res, prec, mode)
}

// BigAsin computes arcsin(x) using the relation: asin(x) = atan(x / sqrt(1 - x²))
func BigAsin(x *BigFloat, prec uint) *BigFloat {
	return getDispatcher().BigAsinImpl(x, prec)
}

// BigAsinRounded computes asin(x) and rounds the result according to the mode
func BigAsinRounded(x *BigFloat, prec uint, mode RoundingMode) (*BigFloat, int) {
	workPrec := prec + 32
	res := BigAsin(x, workPrec)
	return Round(res, prec, mode)
}

// BigAcos computes arccos(x) using the relation: acos(x) = π/2 - asin(x)
func BigAcos(x *BigFloat, prec uint) *BigFloat {
	return getDispatcher().BigAcosImpl(x, prec)
}

// BigAcosRounded computes acos(x) and rounds the result according to the mode
func BigAcosRounded(x *BigFloat, prec uint, mode RoundingMode) (*BigFloat, int) {
	workPrec := prec + 32
	res := BigAcos(x, workPrec)
	return Round(res, prec, mode)
}
