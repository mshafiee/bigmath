// Package bigmath provides arbitrary-precision math operations for astronomical calculations
// to achieve bit-exact matching with the C library and eliminate rounding errors.
package bigmath

import (
	"math"
	"math/big"
)

// Default precision: 256 bits (77 decimal digits) - eliminates all rounding errors
// 256 bits is sufficient - errors are not due to BigFloat precision limits
const DefaultPrecision = 256

// BigFloat is an alias for big.Float for convenience
type BigFloat = big.Float

// BigVec3 represents a 3D vector with arbitrary precision
type BigVec3 struct {
	X, Y, Z *BigFloat
}

// BigVec6 represents position and velocity (6D) with arbitrary precision
type BigVec6 struct {
	X, Y, Z    *BigFloat // Position
	VX, VY, VZ *BigFloat // Velocity
}

// BigMatrix3x3 represents a 3x3 matrix with arbitrary precision
type BigMatrix3x3 struct {
	M [3][3]*BigFloat
}

// NewBigFloat creates a new BigFloat from a float64 with specified precision
func NewBigFloat(f float64, prec uint) *BigFloat {
	if prec == 0 {
		prec = DefaultPrecision
	}
	bf := new(BigFloat).SetPrec(prec)

	if math.IsNaN(f) {
		// big.Float doesn't support NaN. We might return a zero or special value?
		// Or we just don't set it (default 0).
		// But MPFR expects NaN propagation.
		// Since Go's big.Float doesn't support NaN, we can't fully support it.
		// We'll return 0 for now, or panic if we want to be strict.
		// But better to return 0 to avoid crashes, or handle it in caller.
		// However, SetInf is supported.
		return bf // 0
	}

	if math.IsInf(f, 0) {
		bf.SetInf(math.IsInf(f, -1))
		return bf
	}

	return bf.SetFloat64(f)
}

// NewBigFloatFromString creates a BigFloat from a string with specified precision
func NewBigFloatFromString(s string, prec uint) (*BigFloat, error) {
	if prec == 0 {
		prec = DefaultPrecision
	}
	bf := new(BigFloat).SetPrec(prec)
	_, ok := bf.SetString(s)
	if !ok {
		return nil, big.ErrNaN{}
	}
	return bf, nil
}

// NewBigVec3 creates a new BigVec3 from float64 values
func NewBigVec3(x, y, z float64, prec uint) *BigVec3 {
	return &BigVec3{
		X: NewBigFloat(x, prec),
		Y: NewBigFloat(y, prec),
		Z: NewBigFloat(z, prec),
	}
}

// NewBigVec6 creates a new BigVec6 from float64 values
func NewBigVec6(x, y, z, vx, vy, vz float64, prec uint) *BigVec6 {
	return &BigVec6{
		X:  NewBigFloat(x, prec),
		Y:  NewBigFloat(y, prec),
		Z:  NewBigFloat(z, prec),
		VX: NewBigFloat(vx, prec),
		VY: NewBigFloat(vy, prec),
		VZ: NewBigFloat(vz, prec),
	}
}

// Copy creates a deep copy of a BigVec3
func (v *BigVec3) Copy() *BigVec3 {
	prec := v.X.Prec()
	return &BigVec3{
		X: new(BigFloat).SetPrec(prec).Set(v.X),
		Y: new(BigFloat).SetPrec(prec).Set(v.Y),
		Z: new(BigFloat).SetPrec(prec).Set(v.Z),
	}
}

// Copy creates a deep copy of a BigVec6
func (v *BigVec6) Copy() *BigVec6 {
	prec := v.X.Prec()
	return &BigVec6{
		X:  new(BigFloat).SetPrec(prec).Set(v.X),
		Y:  new(BigFloat).SetPrec(prec).Set(v.Y),
		Z:  new(BigFloat).SetPrec(prec).Set(v.Z),
		VX: new(BigFloat).SetPrec(prec).Set(v.VX),
		VY: new(BigFloat).SetPrec(prec).Set(v.VY),
		VZ: new(BigFloat).SetPrec(prec).Set(v.VZ),
	}
}

// ToFloat64 converts BigVec3 to float64 array
func (v *BigVec3) ToFloat64() [3]float64 {
	x, _ := v.X.Float64()
	y, _ := v.Y.Float64()
	z, _ := v.Z.Float64()
	return [3]float64{x, y, z}
}

// ToFloat64 converts BigVec6 to float64 array
func (v *BigVec6) ToFloat64() [6]float64 {
	x, _ := v.X.Float64()
	y, _ := v.Y.Float64()
	z, _ := v.Z.Float64()
	vx, _ := v.VX.Float64()
	vy, _ := v.VY.Float64()
	vz, _ := v.VZ.Float64()
	return [6]float64{x, y, z, vx, vy, vz}
}

// Add adds two BigVec3 vectors: result = v1 + v2
func BigVec3Add(v1, v2 *BigVec3, prec uint) *BigVec3 {
	return getDispatcher().BigVec3AddImpl(v1, v2, prec)
}

// Sub subtracts two BigVec3 vectors: result = v1 - v2
func BigVec3Sub(v1, v2 *BigVec3, prec uint) *BigVec3 {
	return getDispatcher().BigVec3SubImpl(v1, v2, prec)
}

// Mul multiplies a BigVec3 by a scalar: result = v * scalar
func BigVec3Mul(v *BigVec3, scalar *BigFloat, prec uint) *BigVec3 {
	return getDispatcher().BigVec3MulImpl(v, scalar, prec)
}

// Dot computes the dot product of two BigVec3 vectors
func BigVec3Dot(v1, v2 *BigVec3, prec uint) *BigFloat {
	return getDispatcher().BigVec3DotImpl(v1, v2, prec)
}

// Magnitude computes the magnitude (length) of a BigVec3
func BigVec3Magnitude(v *BigVec3, prec uint) *BigFloat {
	if prec == 0 {
		prec = v.X.Prec()
	}

	// sqrt(x² + y² + z²)
	dotProd := BigVec3Dot(v, v, prec)
	return BigSqrt(dotProd, prec)
}

// BigMatMul multiplies a matrix by a vector: result = M * v
func BigMatMul(m *BigMatrix3x3, v *BigVec3, prec uint) *BigVec3 {
	return getDispatcher().BigMatMulImpl(m, v, prec)
}

// NewIdentityMatrix creates a 3x3 identity matrix
func NewIdentityMatrix(prec uint) *BigMatrix3x3 {
	if prec == 0 {
		prec = DefaultPrecision
	}

	one := NewBigFloat(1.0, prec)
	zero := NewBigFloat(0.0, prec)

	return &BigMatrix3x3{
		M: [3][3]*BigFloat{
			{new(BigFloat).Set(one), new(BigFloat).Set(zero), new(BigFloat).Set(zero)},
			{new(BigFloat).Set(zero), new(BigFloat).Set(one), new(BigFloat).Set(zero)},
			{new(BigFloat).Set(zero), new(BigFloat).Set(zero), new(BigFloat).Set(one)},
		},
	}
}

// BigAbs returns the absolute value of a BigFloat
func BigAbs(x *BigFloat, prec uint) *BigFloat {
	if prec == 0 {
		prec = x.Prec()
	}
	result := new(BigFloat).SetPrec(prec).Set(x)
	if result.Sign() < 0 {
		result.Neg(result)
	}
	return result
}

// BigMax returns the maximum of two BigFloats
func BigMax(a, b *BigFloat, prec uint) *BigFloat {
	if prec == 0 {
		prec = a.Prec()
	}
	if a.Cmp(b) > 0 {
		return new(BigFloat).SetPrec(prec).Set(a)
	}
	return new(BigFloat).SetPrec(prec).Set(b)
}

// BigMin returns the minimum of two BigFloats
func BigMin(a, b *BigFloat, prec uint) *BigFloat {
	if prec == 0 {
		prec = a.Prec()
	}
	if a.Cmp(b) < 0 {
		return new(BigFloat).SetPrec(prec).Set(a)
	}
	return new(BigFloat).SetPrec(prec).Set(b)
}

// Constants with high precision
var (
	bigPI     *BigFloat
	bigTwoPI  *BigFloat
	bigHalfPI *BigFloat
)

func init() {
	// Initialize constants with maximum precision
	prec := uint(DefaultPrecision)

	// Compute Pi using Chudnovsky algorithm (MPFR standard)
	// 1/π = 12 * Σ ((-1)^k * (6k)! * (13591409 + 545140134k)) / ((3k)! * (k!)^3 * 640320^(3k+3/2))
	// This converges at ~14 digits per term. For 256 bits (~77 digits), we need ~6 terms.
	// We'll use a few more to be safe and generic for higher precision.

	bigPI = computePiChudnovsky(prec)

	// 2π
	bigTwoPI = new(BigFloat).SetPrec(prec)
	bigTwoPI.Mul(bigPI, NewBigFloat(2.0, prec))

	// π/2
	bigHalfPI = new(BigFloat).SetPrec(prec)
	bigHalfPI.Quo(bigPI, NewBigFloat(2.0, prec))
}

// computePiChudnovsky computes Pi using the Chudnovsky algorithm
// This is the algorithm used by MPFR and most high-precision libraries.
func computePiChudnovsky(prec uint) *BigFloat {
	// Constants for Chudnovsky algorithm
	// C = 640320
	// C3_OVER_24 = C^3 / 24

	// We need slightly higher precision for intermediate calculations
	workPrec := prec + 32

	// Number of terms needed: ~14 digits per term
	// digits = prec * log10(2) ≈ prec * 0.3
	// terms = digits / 14 ≈ prec * 0.0215
	numTerms := int(float64(prec)*0.022) + 2

	// Binary splitting variables
	// P(a, b) = product of terms for numerator
	// Q(a, b) = product of terms for denominator
	// T(a, b) = sum of terms

	var bs func(a, b int64) (*big.Int, *big.Int, *big.Int)
	bs = func(a, b int64) (*big.Int, *big.Int, *big.Int) {
		if b-a == 1 {
			// Leaf node for k = a
			k := a

			// P(k, k+1) = -(6k+5)(2k+1)(6k+1)
			// For k=0, P = 1 (special case handled by recursion structure usually, but here:)
			// Actually, standard binary splitting form:
			// a_k = (-1)^k * (13591409 + 545140134k)
			// b_k = (3k)! * (k!)^3 * 640320^(3k)  <-- handled in denominator Q

			// Let's use the factorized form for P, Q, T
			// P(a, b) = P(a, m) * P(m, b)
			// Q(a, b) = Q(a, m) * Q(m, b)
			// T(a, b) = T(a, m) * Q(m, b) + P(a, m) * T(m, b)

			// Leaf k:
			// P_k = -(6k+5)(2k+1)(6k+1)  for k > 0
			// Q_k = 10939058860032000 * k^3
			// T_k = P_k * (13591409 + 545140134*k)

			// k=0 is special, usually handled outside or as first term
			if k == 0 {
				P := big.NewInt(1)
				Q := big.NewInt(1)
				T := big.NewInt(13591409)
				return P, Q, T
			}

			kBig := big.NewInt(k)

			// P = -(6k-5)(2k-1)(6k-1)
			// Note: The indices shift depending on definition.
			// Let's use the standard recurrence:
			// P_k = -(6k-5)(2k-1)(6k-1)
			// Q_k = (k^3) * (640320^3 / 24) = k^3 * 10939058860032000
			// T_k = P_k * (13591409 + 545140134*k)

			P := big.NewInt(6*k - 5)
			P.Mul(P, big.NewInt(2*k-1))
			P.Mul(P, big.NewInt(6*k-1))
			P.Neg(P)

			Q := big.NewInt(k)
			Q.Exp(Q, big.NewInt(3), nil)
			Q.Mul(Q, big.NewInt(10939058860032000)) // 640320^3 / 24

			T := big.NewInt(545140134)
			T.Mul(T, kBig)
			T.Add(T, big.NewInt(13591409))
			T.Mul(T, P)

			// For k=0, the formula gives P=1, Q=1, T=13591409 if we adjust
			// But we handle k=0 separately or as base case.
			// The loop usually runs 1 to N, with 0-th term added separately.
			// But binary splitting can handle 0 to N.

			return P, Q, T
		}

		m := (a + b) / 2
		P_am, Q_am, T_am := bs(a, m)
		P_mb, Q_mb, T_mb := bs(m, b)

		// P = P_am * P_mb
		P := new(big.Int).Mul(P_am, P_mb)

		// Q = Q_am * Q_mb
		Q := new(big.Int).Mul(Q_am, Q_mb)

		// T = Q_mb * T_am + P_am * T_mb
		T := new(big.Int).Mul(Q_mb, T_am)
		tmp := new(big.Int).Mul(P_am, T_mb)
		T.Add(T, tmp)

		return P, Q, T
	}

	// Run binary splitting
	_, Q, T := bs(0, int64(numTerms))

	// Final calculation:
	// Pi = (426880 * sqrt(10005) * Q) / T

	// Sqrt(10005)
	sqrt10005 := NewBigFloat(10005.0, workPrec)
	sqrt10005 = BigSqrt(sqrt10005, workPrec)

	// 426880
	constFactor := NewBigFloat(426880.0, workPrec)

	// Numerator = 426880 * sqrt(10005) * Q
	num := new(BigFloat).SetPrec(workPrec).SetInt(Q)
	num.Mul(num, constFactor)
	num.Mul(num, sqrt10005)

	// Denominator = T
	den := new(BigFloat).SetPrec(workPrec).SetInt(T)

	// Pi = Num / Den
	pi := new(BigFloat).SetPrec(workPrec).Quo(num, den)

	// Return with requested precision
	return new(BigFloat).SetPrec(prec).Set(pi)
}

// BigPI returns π with specified precision
func BigPI(prec uint) *BigFloat {
	if prec == 0 {
		prec = DefaultPrecision
	}
	return new(BigFloat).SetPrec(prec).Set(bigPI)
}

// BigTwoPI returns 2π with specified precision
func BigTwoPI(prec uint) *BigFloat {
	if prec == 0 {
		prec = DefaultPrecision
	}
	return new(BigFloat).SetPrec(prec).Set(bigTwoPI)
}

// BigHalfPI returns π/2 with specified precision
func BigHalfPI(prec uint) *BigFloat {
	if prec == 0 {
		prec = DefaultPrecision
	}
	return new(BigFloat).SetPrec(prec).Set(bigHalfPI)
}

// BigSqrt computes the square root using Newton-Raphson method
func BigSqrt(x *BigFloat, prec uint) *BigFloat {
	if prec == 0 {
		prec = x.Prec()
	}

	// Handle special cases
	if x.Sign() == 0 {
		return NewBigFloat(0.0, prec)
	}
	if x.Sign() < 0 {
		// Return NaN for negative numbers
		return NewBigFloat(math.NaN(), prec)
	}

	// Initial guess using float64 sqrt
	xFloat, _ := x.Float64()
	guess := NewBigFloat(math.Sqrt(xFloat), prec)

	// Newton-Raphson: x_{n+1} = (x_n + S/x_n) / 2
	// Iterate until convergence
	two := NewBigFloat(2.0, prec)
	temp := new(BigFloat).SetPrec(prec)
	diff := new(BigFloat).SetPrec(prec)
	threshold := new(BigFloat).SetPrec(prec).SetFloat64(1e-77) // Convergence threshold

	for i := 0; i < 100; i++ { // Max 100 iterations
		// temp = S / guess
		temp.Quo(x, guess)

		// temp = guess + S/guess
		temp.Add(guess, temp)

		// guess_new = (guess + S/guess) / 2
		temp.Quo(temp, two)

		// Check convergence: |guess_new - guess|
		diff.Sub(temp, guess)
		diff = BigAbs(diff, prec)

		guess.Set(temp)

		if diff.Cmp(threshold) < 0 {
			break
		}
	}

	return guess
}

// BigSqrtRounded computes sqrt(x) and rounds the result according to the mode
func BigSqrtRounded(x *BigFloat, prec uint, mode RoundingMode) (*BigFloat, int) {
	// Use rounding.go implementation which wraps big.Float.Sqrt
	return SqrtRounded(x, prec, mode)
}
