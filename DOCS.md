# bigmath API Documentation

Complete API reference for the bigmath package.

## Table of Contents

- [Constants](#constants)
- [Types](#types)
- [BigFloat Operations](#bigfloat-operations)
- [Basic Math Utilities](#basic-math-utilities)
- [Vector Operations](#vector-operations)
- [Advanced Vector Operations](#advanced-vector-operations)
- [Matrix Operations](#matrix-operations)
- [Advanced Matrix Operations](#advanced-matrix-operations)
- [Trigonometric Functions](#trigonometric-functions)
- [Hyperbolic Functions](#hyperbolic-functions)
- [Exponential and Logarithmic Functions](#exponential-and-logarithmic-functions)
- [Advanced Logarithmic Functions](#advanced-logarithmic-functions)
- [Power Functions](#power-functions)
- [Root Functions](#root-functions)
- [Special Functions](#special-functions)
- [Combinatorics](#combinatorics)
- [Rounding Functions](#rounding-functions)
- [Angle Normalization](#angle-normalization)
- [Chebyshev Polynomial Evaluation](#chebyshev-polynomial-evaluation)
- [Mathematical Constants](#mathematical-constants)
- [Extended Constants](#extended-constants)
- [Serialization](#serialization)
- [Error Handling](#error-handling)
- [CPU Feature Detection](#cpu-feature-detection)
- [Extended Precision Mode](#extended-precision-mode)

## Constants

### DefaultPrecision

```go
const DefaultPrecision = 256
```

Default precision in bits (77 decimal digits). Used when `prec` parameter is 0.

### ExtendedPrecision

```go
const ExtendedPrecision = 80
```

Extended precision constant (80 bits) that enables hardware extended precision mode using the x87 FPU.
When `prec == ExtendedPrecision` and x87 is available (x86/x86-64 platforms), operations use the hardware
80-bit extended precision format for faster intermediate calculations. On other platforms or when x87 is
unavailable, operations automatically fall back to BigFloat implementations.

**Platform Support**: Only available on x86/x86-64 (amd64 or 386) platforms with x87 FPU support.

**Usage**:
```go
// Use extended precision for sin calculation
x := bigmath.NewBigFloat(math.Pi/4, bigmath.ExtendedPrecision)
result := bigmath.BigSin(x, bigmath.ExtendedPrecision)
```

**Supported Operations**: Trigonometric (sin, cos, tan, atan, atan2), exponential/logarithmic (exp, log),
power (pow), and square root (sqrt) functions support extended precision mode.

### Rounding Modes

```go
const (
    ToNearest     RoundingMode = big.ToNearestEven  // Round to nearest, ties to even
    ToNearestAway RoundingMode = big.ToNearestAway  // Round to nearest, ties away from zero
    ToZero        RoundingMode = big.ToZero         // Round toward zero
    ToPositiveInf RoundingMode = big.ToPositiveInf   // Round toward +∞
    ToNegativeInf RoundingMode = big.ToNegativeInf  // Round toward -∞
    AwayFromZero  RoundingMode = big.AwayFromZero   // Round away from zero
)
```

## Types

### BigFloat

```go
type BigFloat = big.Float
```

An alias for `big.Float` providing arbitrary-precision floating-point arithmetic.

### BigVec3

```go
type BigVec3 struct {
    X, Y, Z *BigFloat
}
```

A 3D vector with arbitrary-precision components.

### BigVec6

```go
type BigVec6 struct {
    X, Y, Z    *BigFloat // Position
    VX, VY, VZ *BigFloat // Velocity
}
```

A 6D vector representing position and velocity.

### BigMatrix3x3

```go
type BigMatrix3x3 struct {
    M [3][3]*BigFloat
}
```

A 3x3 matrix with arbitrary-precision elements.

### RoundingMode

```go
type RoundingMode = big.RoundingMode
```

An alias for `big.RoundingMode` specifying rounding behavior.

### SegmentInfoBig

```go
type SegmentInfoBig struct {
    Body         int
    SegmentStart *BigFloat
    SegmentEnd   *BigFloat
    SegmentSize  *BigFloat
    NumCoeffs    int
    ElemEpoch    *BigFloat
    Qrot         *BigFloat
    DQrot        *BigFloat
    Prot         *BigFloat
    DProt        *BigFloat
    Peri         *BigFloat
    DPeri        *BigFloat
    RefEllipse   []*BigFloat
    Flags        int
}
```

Holds segment information for Chebyshev polynomial evaluation (used in astronomical calculations).

### CPUFeatures

```go
type CPUFeatures struct {
    // AMD64 features
    HasBMI2  bool
    HasAVX2  bool
    HasFMA   bool
    
    // ARM64 features
    HasNEON  bool
}
```

CPU feature detection results.

## BigFloat Operations

### NewBigFloat

```go
func NewBigFloat(f float64, prec uint) *BigFloat
```

Creates a new `BigFloat` from a `float64` with specified precision. If `prec` is 0, uses `DefaultPrecision`.

**Example:**
```go
x := bigmath.NewBigFloat(3.14159, 256)
```

### NewBigFloatFromString

```go
func NewBigFloatFromString(s string, prec uint) (*BigFloat, error)
```

Creates a new `BigFloat` from a string representation. Returns an error if the string is invalid.

**Example:**
```go
x, err := bigmath.NewBigFloatFromString("3.141592653589793238462643383279", 256)
```

### BigAbs

```go
func BigAbs(x *BigFloat, prec uint) *BigFloat
```

Returns the absolute value of `x`.

### BigMax

```go
func BigMax(a, b *BigFloat, prec uint) *BigFloat
```

Returns the maximum of `a` and `b`.

### BigMin

```go
func BigMin(a, b *BigFloat, prec uint) *BigFloat
```

Returns the minimum of `a` and `b`.

### BigSqrt

```go
func BigSqrt(x *BigFloat, prec uint) *BigFloat
```

Computes the square root of `x` using Newton-Raphson method. Returns NaN for negative inputs.

### BigSqrtRounded

```go
func BigSqrtRounded(x *BigFloat, prec uint, mode RoundingMode) (*BigFloat, int)
```

Computes the square root with specified rounding mode. Returns the result and rounding direction.

### BigFloatFMA

```go
func BigFloatFMA(a, b, c *BigFloat, prec uint) *BigFloat
```

Fused multiply-add: computes `a * b + c` with a single rounding operation.

## Basic Math Utilities

### BigFloor

```go
func BigFloor(x *BigFloat, prec uint) *BigFloat
```

Returns the greatest integer value less than or equal to x (rounds toward negative infinity).

### BigCeil

```go
func BigCeil(x *BigFloat, prec uint) *BigFloat
```

Returns the smallest integer value greater than or equal to x (rounds toward positive infinity).

### BigTrunc

```go
func BigTrunc(x *BigFloat, prec uint) *BigFloat
```

Returns the integer value of x truncated toward zero.

### BigMod

```go
func BigMod(x, y *BigFloat, prec uint) *BigFloat
```

Returns x mod y (x - y*floor(x/y)). The result has the same sign as y.

### BigRem

```go
func BigRem(x, y *BigFloat, prec uint) *BigFloat
```

Returns the remainder of x/y (IEEE 754 remainder operation). The result has the same sign as x.

## Vector Operations

### NewBigVec3

```go
func NewBigVec3(x, y, z float64, prec uint) *BigVec3
```

Creates a new 3D vector from float64 values.

### NewBigVec6

```go
func NewBigVec6(x, y, z, vx, vy, vz float64, prec uint) *BigVec6
```

Creates a new 6D vector (position and velocity) from float64 values.

### BigVec3Add

```go
func BigVec3Add(v1, v2 *BigVec3, prec uint) *BigVec3
```

Adds two 3D vectors: `result = v1 + v2`.

### BigVec3Sub

```go
func BigVec3Sub(v1, v2 *BigVec3, prec uint) *BigVec3
```

Subtracts two 3D vectors: `result = v1 - v2`.

### BigVec3Mul

```go
func BigVec3Mul(v *BigVec3, scalar *BigFloat, prec uint) *BigVec3
```

Multiplies a 3D vector by a scalar: `result = v * scalar`.

### BigVec3Dot

```go
func BigVec3Dot(v1, v2 *BigVec3, prec uint) *BigFloat
```

Computes the dot product of two 3D vectors: `v1 · v2`.

### BigVec3Magnitude

```go
func BigVec3Magnitude(v *BigVec3, prec uint) *BigFloat
```

Computes the magnitude (length) of a 3D vector: `√(x² + y² + z²)`.

### BigVec6Add

```go
func BigVec6Add(v1, v2 *BigVec6, prec uint) *BigVec6
```

Adds two 6D vectors.

### BigVec6Sub

```go
func BigVec6Sub(v1, v2 *BigVec6, prec uint) *BigVec6
```

Subtracts two 6D vectors.

### BigVec6Negate

```go
func BigVec6Negate(v *BigVec6, prec uint) *BigVec6
```

Negates all components of a 6D vector.

### BigVec6Magnitude

```go
func BigVec6Magnitude(v *BigVec6, prec uint) *BigFloat
```

Computes the magnitude of the position component of a 6D vector.

### ApplyRotationMatrixToBigVec6

```go
func ApplyRotationMatrixToBigVec6(m *BigMatrix3x3, v *BigVec6, prec uint) *BigVec6
```

Applies a rotation matrix to both position and velocity components of a 6D vector.

### Copy Methods

```go
func (v *BigVec3) Copy() *BigVec3
func (v *BigVec6) Copy() *BigVec6
```

Creates a deep copy of the vector.

### ToFloat64 Methods

```go
func (v *BigVec3) ToFloat64() [3]float64
func (v *BigVec6) ToFloat64() [6]float64
```

Converts the vector to a float64 array.

## Advanced Vector Operations

### BigVec3Cross

```go
func BigVec3Cross(v1, v2 *BigVec3, prec uint) *BigVec3
```

Computes the cross product of two 3D vectors: `v1 × v2`.

### BigVec3Normalize

```go
func BigVec3Normalize(v *BigVec3, prec uint) *BigVec3
```

Normalizes a 3D vector to unit length. Returns a unit vector in the same direction, or zero vector if input is zero.

### BigVec3Angle

```go
func BigVec3Angle(v1, v2 *BigVec3, prec uint) *BigFloat
```

Computes the angle between two 3D vectors in radians using: `angle = arccos((v1·v2) / (|v1|*|v2|))`. Returns angle in range [0, π].

### BigVec3Project

```go
func BigVec3Project(v1, v2 *BigVec3, prec uint) *BigVec3
```

Projects vector v1 onto vector v2: `((v1·v2) / |v2|²) * v2`.

## Matrix Operations

### NewIdentityMatrix

```go
func NewIdentityMatrix(prec uint) *BigMatrix3x3
```

Creates a 3x3 identity matrix.

### BigMatMul

```go
func BigMatMul(m *BigMatrix3x3, v *BigVec3, prec uint) *BigVec3
```

Multiplies a 3x3 matrix by a 3D vector: `result = M * v`.

### CreateRotationMatrix

```go
func CreateRotationMatrix(angles [3]*BigFloat, prec uint) *BigMatrix3x3
```

Creates a rotation matrix from three Euler angles (used for precession and coordinate transformations).

## Advanced Matrix Operations

### BigMatTranspose

```go
func BigMatTranspose(m *BigMatrix3x3, prec uint) *BigMatrix3x3
```

Returns the transpose of a 3x3 matrix.

### BigMatMulMat

```go
func BigMatMulMat(m1, m2 *BigMatrix3x3, prec uint) *BigMatrix3x3
```

Multiplies two 3x3 matrices: `result = m1 * m2`.

### BigMatDet

```go
func BigMatDet(m *BigMatrix3x3, prec uint) *BigFloat
```

Computes the determinant of a 3x3 matrix using cofactor expansion.

### BigMatInverse

```go
func BigMatInverse(m *BigMatrix3x3, prec uint) (*BigMatrix3x3, error)
```

Computes the inverse of a 3x3 matrix using adjugate/determinant. Returns error if matrix is singular (determinant is zero).

## Trigonometric Functions

### BigSin

```go
func BigSin(x *BigFloat, prec uint) *BigFloat
```

Computes sin(x) using Taylor series expansion.

### BigCos

```go
func BigCos(x *BigFloat, prec uint) *BigFloat
```

Computes cos(x) using Taylor series expansion.

### BigTan

```go
func BigTan(x *BigFloat, prec uint) *BigFloat
```

Computes tan(x) = sin(x) / cos(x).

### BigAtan

```go
func BigAtan(x *BigFloat, prec uint) *BigFloat
```

Computes arctan(x) using Taylor series. For |x| > 1, uses atan(x) = π/2 - atan(1/x).

### BigAtan2

```go
func BigAtan2(y, x *BigFloat, prec uint) *BigFloat
```

Computes atan2(y, x), returning the angle in radians between the positive x-axis and the point (x, y).

### BigAsin

```go
func BigAsin(x *BigFloat, prec uint) *BigFloat
```

Computes arcsin(x) using the relation: asin(x) = atan(x / sqrt(1 - x²)).

### BigAcos

```go
func BigAcos(x *BigFloat, prec uint) *BigFloat
```

Computes arccos(x) using the relation: acos(x) = π/2 - asin(x).

### Rounded Variants

All trigonometric functions have `Rounded` variants that accept a rounding mode:

```go
func BigSinRounded(x *BigFloat, prec uint, mode RoundingMode) (*BigFloat, int)
func BigCosRounded(x *BigFloat, prec uint, mode RoundingMode) (*BigFloat, int)
func BigTanRounded(x *BigFloat, prec uint, mode RoundingMode) (*BigFloat, int)
func BigAtanRounded(x *BigFloat, prec uint, mode RoundingMode) (*BigFloat, int)
func BigAtan2Rounded(y, x *BigFloat, prec uint, mode RoundingMode) (*BigFloat, int)
func BigAsinRounded(x *BigFloat, prec uint, mode RoundingMode) (*BigFloat, int)
func BigAcosRounded(x *BigFloat, prec uint, mode RoundingMode) (*BigFloat, int)
```

## Hyperbolic Functions

### BigSinh

```go
func BigSinh(x *BigFloat, prec uint) *BigFloat
```

Computes sinh(x) = (e^x - e^-x) / 2.

### BigCosh

```go
func BigCosh(x *BigFloat, prec uint) *BigFloat
```

Computes cosh(x) = (e^x + e^-x) / 2.

### BigTanh

```go
func BigTanh(x *BigFloat, prec uint) *BigFloat
```

Computes tanh(x) = sinh(x) / cosh(x).

### BigAsinh

```go
func BigAsinh(x *BigFloat, prec uint) *BigFloat
```

Computes asinh(x) = ln(x + sqrt(x² + 1)).

### BigAcosh

```go
func BigAcosh(x *BigFloat, prec uint) *BigFloat
```

Computes acosh(x) = ln(x + sqrt(x² - 1)) for x ≥ 1.

### BigAtanh

```go
func BigAtanh(x *BigFloat, prec uint) *BigFloat
```

Computes atanh(x) = (1/2) * ln((1+x)/(1-x)) for |x| < 1.

## Exponential and Logarithmic Functions

### BigExp

```go
func BigExp(x *BigFloat, prec uint) *BigFloat
```

Computes e^x using argument reduction and Taylor series expansion.

**Algorithm:**
1. Argument reduction: x = r + k*ln(2) where |r| ≤ ln(2)/2
2. Further reduction: exp(r) = (exp(r/2^p))^(2^p)
3. Taylor series with binary splitting for exp(r/2^p)

### BigLog

```go
func BigLog(x *BigFloat, prec uint) *BigFloat
```

Computes the natural logarithm ln(x) using argument reduction and series expansion.

### BigLog10

```go
func BigLog10(x *BigFloat, prec uint) *BigFloat
```

Computes the base-10 logarithm log₁₀(x) = ln(x) / ln(10).

### BigLog2

```go
func BigLog2(prec uint) *BigFloat
```

Returns ln(2) with specified precision.

## Advanced Logarithmic Functions

### BigLog1p

```go
func BigLog1p(x *BigFloat, prec uint) *BigFloat
```

Computes log(1+x) accurately for values near zero. Uses series expansion to avoid precision loss when x is very small.

### BigExp1m

```go
func BigExp1m(x *BigFloat, prec uint) *BigFloat
```

Computes exp(x)-1 accurately for values near zero. Uses series expansion to avoid precision loss when x is very small.

### BigLogb

```go
func BigLogb(x, base *BigFloat, prec uint) *BigFloat
```

Computes the logarithm of x with base b: `log_b(x) = ln(x) / ln(b)`.

## Power Functions

### BigPow

```go
func BigPow(x, y *BigFloat, prec uint) *BigFloat
```

Computes x^y. Uses exp(y * ln(x)) for non-integer y. For integer exponents, uses binary exponentiation.

**Special cases:**
- x^0 = 1
- 1^y = 1
- x^1 = x
- For negative base with non-integer exponent, returns NaN

## Root Functions

### BigSqrt

```go
func BigSqrt(x *BigFloat, prec uint) *BigFloat
```

Computes √x using Newton-Raphson method. Returns NaN for negative inputs.

**Note:** See [Root Functions](#root-functions) for additional root functions like cube root and nth root.

## Rounding Functions

### Round

```go
func Round(x *BigFloat, prec uint, mode RoundingMode) (*BigFloat, int)
```

Rounds `x` to `prec` bits using the specified rounding mode. Returns the rounded value and rounding direction (-1, 0, or +1).

### SqrtRounded

```go
func SqrtRounded(x *BigFloat, prec uint, mode RoundingMode) (*BigFloat, int)
```

Computes square root with specified rounding mode.

### AddRounded

```go
func AddRounded(a, b *BigFloat, prec uint, mode RoundingMode) (*BigFloat, int)
```

Adds two numbers with specified rounding mode.

### QuoRounded

```go
func QuoRounded(a, b *BigFloat, prec uint, mode RoundingMode) (*BigFloat, int)
```

Divides two numbers with specified rounding mode.

## Angle Normalization

### DegNormBig

```go
func DegNormBig(deg *BigFloat, prec uint) *BigFloat
```

Normalizes an angle in degrees to the range [0, 360).

### RadNormBig

```go
func RadNormBig(rad *BigFloat, prec uint) *BigFloat
```

Normalizes an angle in radians to the range [-π, π].

### RadNorm02PiBig

```go
func RadNorm02PiBig(rad *BigFloat, prec uint) *BigFloat
```

Normalizes an angle in radians to the range [0, 2π).

## Chebyshev Polynomial Evaluation

### EvaluateChebyshevBig

```go
func EvaluateChebyshevBig(t *BigFloat, c []*BigFloat, neval int, prec uint) *BigFloat
```

Evaluates a Chebyshev polynomial series at point `t` with coefficients `c`. `neval` is the number of coefficients to evaluate.

**Algorithm:** Uses Clenshaw's recurrence relation for numerical stability.

### EvaluateChebyshevDerivativeBig

```go
func EvaluateChebyshevDerivativeBig(t *BigFloat, c []*BigFloat, neval int, prec uint) *BigFloat
```

Evaluates the derivative of a Chebyshev polynomial series at point `t`.

### EvaluateSegmentBig

```go
func EvaluateSegmentBig(tjd *BigFloat, coeffs []*BigFloat, segStart, segEnd *BigFloat, neval int, prec uint) *BigVec6
```

Evaluates a Chebyshev segment (used in astronomical ephemeris calculations). Returns a 6D vector (position and velocity).

### RotateCoeffsToJ2000Big

```go
func RotateCoeffsToJ2000Big(coeffs []*BigFloat, segInfo *SegmentInfoBig, isMoon bool, prec uint) ([]*BigFloat, int)
```

Rotates Chebyshev coefficients from orbital plane to equatorial J2000 coordinates using arbitrary precision.

### ConvertToBigFloatCoeffs

```go
func ConvertToBigFloatCoeffs(coeffsFloat64 []float64, prec uint) []*BigFloat
```

Converts float64 coefficients to BigFloat coefficients.

## Mathematical Constants

### BigPI

```go
func BigPI(prec uint) *BigFloat
```

Returns π (pi) with specified precision. Computed using Chudnovsky algorithm.

### BigTwoPI

```go
func BigTwoPI(prec uint) *BigFloat
```

Returns 2π with specified precision.

### BigHalfPI

```go
func BigHalfPI(prec uint) *BigFloat
```

Returns π/2 with specified precision.

### BigE

```go
func BigE(prec uint) *BigFloat
```

Returns Euler's number e with specified precision.

### BigEulerGamma

```go
func BigEulerGamma(prec uint) *BigFloat
```

Returns Euler-Mascheroni constant γ with specified precision.

### BigCatalan

```go
func BigCatalan(prec uint) *BigFloat
```

Returns Catalan's constant with specified precision.

### BigJ2000

```go
func BigJ2000(prec uint) *BigFloat
```

Returns the Julian Day number for J2000.0 epoch (2451545.0).

### BigJulianCentury

```go
func BigJulianCentury(prec uint) *BigFloat
```

Returns the number of days in a Julian century (36525.0).

### BigJulianMillennium

```go
func BigJulianMillennium(prec uint) *BigFloat
```

Returns the number of days in a Julian millennium (365250.0).

### BigLightSpeedAUperDay

```go
func BigLightSpeedAUperDay(prec uint) *BigFloat
```

Returns the speed of light in AU per day (173.14463267424034).

## Extended Constants

### BigPhi

```go
func BigPhi(prec uint) *BigFloat
```

Returns the golden ratio φ = (1 + √5) / 2 ≈ 1.6180339887498948482... with specified precision.

### BigSqrt2

```go
func BigSqrt2(prec uint) *BigFloat
```

Returns √2 ≈ 1.4142135623730950488... with specified precision.

### BigSqrt3

```go
func BigSqrt3(prec uint) *BigFloat
```

Returns √3 ≈ 1.7320508075688772935... with specified precision.

### BigLn10

```go
func BigLn10(prec uint) *BigFloat
```

Returns ln(10) ≈ 2.3025850929940456840... with specified precision.

## Serialization

### BigFloatMarshalJSON

```go
func BigFloatMarshalJSON(x *BigFloat) ([]byte, error)
```

Marshals a BigFloat to JSON using string representation for precision.

### BigFloatUnmarshalJSON

```go
func BigFloatUnmarshalJSON(data []byte, prec uint) (*BigFloat, error)
```

Unmarshals a BigFloat from JSON.

### BigVec3 JSON Methods

`BigVec3` implements `json.Marshaler` and `json.Unmarshaler` interfaces. Vectors are serialized as arrays of strings.

### BigVec6 JSON Methods

`BigVec6` implements `json.Marshaler` and `json.Unmarshaler` interfaces. Vectors are serialized as arrays of strings.

### BigMatrix3x3 JSON Methods

`BigMatrix3x3` implements `json.Marshaler` and `json.Unmarshaler` interfaces. Matrices are serialized as 3x3 arrays of strings.

## Error Handling

### Ulp

```go
func Ulp(x *BigFloat, prec uint) *BigFloat
```

Returns the Unit in the Last Place (ULP) for `x` at the specified precision.

### ErrorBound

```go
type ErrorBound interface {
    GetUlps() float64
    GetAbsError(x *BigFloat, prec uint) *BigFloat
}
```

Interface for representing error bounds.

### NewUlpError

```go
func NewUlpError(ulps float64, prec uint) ErrorBound
```

Creates an error bound specified in ULPs.

### NewAbsError

```go
func NewAbsError(absVal *BigFloat, prec uint) ErrorBound
```

Creates an error bound specified as an absolute value.

### AddErrorBounds

```go
func AddErrorBounds(e1, e2 ErrorBound, x *BigFloat, prec uint) ErrorBound
```

Adds two error bounds.

### PropagateErrorAdd

```go
func PropagateErrorAdd(x, y, z *BigFloat, errX, errY ErrorBound, prec uint, mode RoundingMode) ErrorBound
```

Propagates error through addition operation.

### PropagateErrorMul

```go
func PropagateErrorMul(x, y, z *BigFloat, errX, errY ErrorBound, prec uint, mode RoundingMode) ErrorBound
```

Propagates error through multiplication operation.

### CalculateRequiredPrecision

```go
func CalculateRequiredPrecision(targetPrec uint, expectedErrorUlps float64) uint
```

Calculates the required precision to achieve a target precision with expected error in ULPs.

## CPU Feature Detection

### GetCPUFeatures

```go
func GetCPUFeatures() CPUFeatures
```

Returns detected CPU features. Used internally by the dispatcher to select optimal code paths.

**AMD64 Features:**
- `HasBMI2`: Bit Manipulation Instructions 2 (enables dual carry chains)
- `HasAVX2`: Advanced Vector Extensions 2
- `HasFMA`: Fused Multiply-Add

**ARM64 Features:**
- `HasNEON`: NEON SIMD instructions (always available on ARMv8)
- `HasX87`: x87 FPU support (80-bit extended precision, x86/x86-64 only)

## Extended Precision Mode

The library supports hardware extended precision (80-bit x87 FPU) mode for faster intermediate calculations
on x86/x86-64 platforms. This mode uses the x87 FPU's native 80-bit extended precision format instead of
arbitrary-precision BigFloat calculations.

### Activation

Extended precision mode is activated by setting `prec = ExtendedPrecision` (80) when calling functions:

```go
// Use extended precision for trigonometric calculations
x := bigmath.NewBigFloat(math.Pi/4, bigmath.ExtendedPrecision)
sinX := bigmath.BigSin(x, bigmath.ExtendedPrecision)
cosX := bigmath.BigCos(x, bigmath.ExtendedPrecision)
```

### Platform Support

Extended precision mode is only available on:
- x86-64 (amd64) platforms
- x86 (386) platforms
- Systems with x87 FPU support

On other platforms or when x87 is unavailable, operations automatically fall back to BigFloat implementations.

### Supported Operations

The following operations support extended precision mode:
- **Trigonometric**: `BigSin`, `BigCos`, `BigTan`, `BigAtan`, `BigAtan2`
- **Exponential/Logarithmic**: `BigExp`, `BigLog`
- **Power**: `BigPow`
- **Root**: `BigSqrt`

### Checking Availability

Use `CanUseExtendedPrecision(prec)` to check if extended precision can be used:

```go
if bigmath.CanUseExtendedPrecision(bigmath.ExtendedPrecision) {
    // Extended precision is available
    result := bigmath.BigSin(x, bigmath.ExtendedPrecision)
} else {
    // Fall back to BigFloat
    result := bigmath.BigSin(x, 256)
}
```

### Performance Considerations

Extended precision mode provides:
- **Faster execution**: Hardware x87 FPU operations are typically faster than software BigFloat calculations
- **80-bit precision**: Provides approximately 19 decimal digits of precision (64-bit mantissa)
- **Automatic fallback**: Seamlessly falls back to BigFloat when extended precision is unavailable

**Use cases:**
- Intermediate calculations where 80-bit precision is sufficient
- Performance-critical code paths on x86/x86-64 platforms
- When exact arbitrary precision is not required

**Limitations:**
- Fixed 80-bit precision (not arbitrary)
- Platform-specific (x86/x86-64 only)
- May have different rounding behavior than BigFloat

## Performance Notes

- All functions automatically select optimized assembly implementations when available
- CPU feature detection happens once at package initialization
- Precision parameter `prec` of 0 uses `DefaultPrecision` (256 bits)
- Functions with `Rounded` suffix compute with higher precision internally, then round to the requested precision
- Assembly optimizations provide 20-40% performance improvements over pure Go implementations

## Thread Safety

- The package is thread-safe for concurrent use
- CPU feature detection is performed once and cached
- Function dispatcher is initialized once using `sync.Once`

