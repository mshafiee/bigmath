# bigmath

[![Go Reference](https://pkg.go.dev/badge/github.com/mshafiee/bigmath.svg)](https://pkg.go.dev/github.com/mshafiee/bigmath)
[![License: BSD-3-Clause](https://img.shields.io/badge/License-BSD--3--Clause-blue.svg)](https://opensource.org/licenses/BSD-3-Clause)
[![Go Report Card](https://goreportcard.com/badge/github.com/mshafiee/bigmath)](https://goreportcard.com/report/github.com/mshafiee/bigmath)

A high-performance arbitrary-precision mathematics library for Go, optimized with assembly implementations for AMD64 and ARM64 architectures. Provides MPFR-compatible algorithms for scientific computing, astronomical calculations, and high-precision numerical analysis.

## Features

- ⚡ **High Performance**: Assembly-optimized implementations achieving 20-40% performance improvements
- 🔢 **Arbitrary Precision**: Default 256-bit precision (77 decimal digits), configurable to any precision
- 🎯 **Comprehensive Math**: Trigonometric, hyperbolic, exponential, logarithmic, and power functions
- 📐 **Vector & Matrix Operations**: 3D/6D vectors and 3x3 matrices with arbitrary precision
- 🧮 **Chebyshev Polynomials**: Optimized evaluation for astronomical ephemeris calculations
- 🔍 **CPU Feature Detection**: Automatically uses optimal code paths (BMI2, AVX2, NEON)
- ✅ **MPFR-Compatible**: Algorithms designed to match MPFR behavior for scientific computing

## Installation

```bash
go get github.com/mshafiee/bigmath
```

## Why bigmath?

While Go's standard `math/big` package provides arbitrary-precision arithmetic, `bigmath` adds:

- **Assembly optimizations** for critical operations (20-40% faster)
- **Complete transcendental functions** (trig, hyperbolic, exp, log) with MPFR-compatible algorithms
- **Vector and matrix operations** optimized for scientific computing
- **Chebyshev polynomial evaluation** for astronomical calculations
- **CPU-specific optimizations** that automatically adapt to your hardware

Perfect for applications requiring both high precision and high performance, such as astronomical ephemeris calculations, financial modeling, and scientific simulations.

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/mshafiee/bigmath"
)

func main() {
    // Create a high-precision number
    x := bigmath.NewBigFloat(3.141592653589793, 256)
    
    // Compute sin(x) with arbitrary precision
    sinX := bigmath.BigSin(x, 256)
    
    // Convert back to float64
    result, _ := sinX.Float64()
    fmt.Printf("sin(π) ≈ %g\n", result)
    
    // Vector operations
    v1 := bigmath.NewBigVec3(1.0, 2.0, 3.0, 256)
    v2 := bigmath.NewBigVec3(4.0, 5.0, 6.0, 256)
    sum := bigmath.BigVec3Add(v1, v2, 256)
    
    fmt.Printf("Vector sum: (%g, %g, %g)\n", 
        sum.X.Float64(), sum.Y.Float64(), sum.Z.Float64())
}
```

## Core Types

### BigFloat
An alias for `big.Float` providing arbitrary-precision floating-point arithmetic.

```go
// Create from float64
x := bigmath.NewBigFloat(3.14, 256)

// Create from string
x, err := bigmath.NewBigFloatFromString("3.141592653589793238462643383279", 256)

// Use constants
pi := bigmath.BigPI(256)
twoPi := bigmath.BigTwoPI(256)
halfPi := bigmath.BigHalfPI(256)
```

### BigVec3
A 3D vector with arbitrary-precision components.

```go
v := bigmath.NewBigVec3(1.0, 2.0, 3.0, 256)

// Operations
sum := bigmath.BigVec3Add(v1, v2, 256)
diff := bigmath.BigVec3Sub(v1, v2, 256)
scaled := bigmath.BigVec3Mul(v, scalar, 256)
dot := bigmath.BigVec3Dot(v1, v2, 256)
mag := bigmath.BigVec3Magnitude(v, 256)
```

### BigVec6
A 6D vector representing position and velocity (X, Y, Z, VX, VY, VZ).

```go
v := bigmath.NewBigVec6(1.0, 2.0, 3.0, 0.1, 0.2, 0.3, 256)

// Operations
sum := bigmath.BigVec6Add(v1, v2, 256)
diff := bigmath.BigVec6Sub(v1, v2, 256)
neg := bigmath.BigVec6Negate(v, 256)
mag := bigmath.BigVec6Magnitude(v, 256)
```

### BigMatrix3x3
A 3x3 matrix with arbitrary-precision elements.

```go
// Create identity matrix
m := bigmath.NewIdentityMatrix(256)

// Matrix-vector multiplication
result := bigmath.BigMatMul(m, v, 256)
```

## Mathematical Functions

### Trigonometric Functions

```go
x := bigmath.NewBigFloat(math.Pi/4, 256)

sinX := bigmath.BigSin(x, 256)
cosX := bigmath.BigCos(x, 256)
tanX := bigmath.BigTan(x, 256)

asinX := bigmath.BigAsin(x, 256)
acosX := bigmath.BigAcos(x, 256)
atanX := bigmath.BigAtan(x, 256)
atan2XY := bigmath.BigAtan2(y, x, 256)
```

### Hyperbolic Functions

```go
x := bigmath.NewBigFloat(1.0, 256)

sinhX := bigmath.BigSinh(x, 256)
coshX := bigmath.BigCosh(x, 256)
tanhX := bigmath.BigTanh(x, 256)

asinhX := bigmath.BigAsinh(x, 256)
acoshX := bigmath.BigAcosh(x, 256)
atanhX := bigmath.BigAtanh(x, 256)
```

### Exponential and Logarithmic Functions

```go
x := bigmath.NewBigFloat(2.0, 256)

// Exponential
expX := bigmath.BigExp(x, 256)

// Natural logarithm
logX := bigmath.BigLog(x, 256)

// Base-2 logarithm
log2X := bigmath.BigLog2(256)

// Power: x^y
powXY := bigmath.BigPow(x, y, 256)
```

### Square Root

```go
x := bigmath.NewBigFloat(2.0, 256)
sqrtX := bigmath.BigSqrt(x, 256)
```

## Rounding Modes

The library supports MPFR-compatible rounding modes:

```go
x := bigmath.NewBigFloat(1.5, 256)

// Round to nearest (ties to even)
rounded, _ := bigmath.Round(x, 64, bigmath.RoundNearestEven)

// Round toward zero
rounded, _ := bigmath.Round(x, 64, bigmath.RoundTowardZero)

// Round toward +infinity
rounded, _ := bigmath.Round(x, 64, bigmath.RoundTowardPosInf)

// Round toward -infinity
rounded, _ := bigmath.Round(x, 64, bigmath.RoundTowardNegInf)

// Round away from zero
rounded, _ := bigmath.Round(x, 64, bigmath.RoundAwayFromZero)
```

Functions with `Rounded` suffix support rounding:

```go
sinX, _ := bigmath.BigSinRounded(x, 256, bigmath.RoundNearestEven)
```

## Chebyshev Polynomial Evaluation

For evaluating Chebyshev polynomial series (useful for astronomical calculations):

```go
// Coefficients for Chebyshev series
coeffs := []*bigmath.BigFloat{
    bigmath.NewBigFloat(1.0, 256),
    bigmath.NewBigFloat(0.5, 256),
    bigmath.NewBigFloat(0.25, 256),
}

t := bigmath.NewBigFloat(0.5, 256)
result := bigmath.EvaluateChebyshevBig(t, coeffs, len(coeffs), 256)
derivative := bigmath.EvaluateChebyshevDerivativeBig(t, coeffs, len(coeffs), 256)
```

## Performance Optimizations

The library includes extensive assembly optimizations:

- **AMD64**: Uses MULX, ADCX/ADOX (BMI2), AVX2 when available
- **ARM64**: Uses NEON, LDP/STP for efficient memory access
- **CPU Detection**: Automatically selects optimal code paths
- **Loop Unrolling**: Processes 4 limbs at a time in multi-precision operations
- **Fused Multiply-Add**: Optimized FMA operations where supported

Performance improvements:
- Multi-precision operations: 20-30% faster
- Chebyshev evaluation: 30-40% faster
- Trigonometric functions: 30-40% faster
- Exponential/logarithmic: 25-35% faster

## Architecture Support

- **AMD64/x86_64**: Full assembly support with BMI2 optimizations
- **ARM64**: Full assembly support with NEON optimizations
- **Generic**: Pure Go fallback implementations for other architectures

## Precision

Default precision is 256 bits (77 decimal digits), which is sufficient for most astronomical and scientific calculations. You can specify any precision:

```go
// Low precision (64 bits)
x := bigmath.NewBigFloat(3.14, 64)

// High precision (1024 bits)
x := bigmath.NewBigFloat(3.14, 1024)
```

## Use Cases

- **Astronomical Calculations**: Ephemeris computation, orbital mechanics
- **Scientific Computing**: High-precision numerical analysis
- **Financial Calculations**: Currency and interest rate computations requiring exact precision
- **Cryptography**: Key generation and verification requiring exact arithmetic
- **Computer Algebra Systems**: Symbolic computation and exact arithmetic

## Examples

See the [examples directory](examples/) for more detailed usage examples.

## API Documentation

See [DOCS.md](DOCS.md) for complete API documentation.

## Requirements

- Go 1.21 or later
- AMD64 or ARM64 architecture (fallback to pure Go on other architectures)

## Performance

Benchmarked improvements over standard implementations:
- Multi-precision operations: **20-30% faster**
- Chebyshev evaluation: **30-40% faster**
- Trigonometric functions: **30-40% faster**
- Exponential/logarithmic: **25-35% faster**

## License

This project is licensed under the BSD 3-Clause License - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

## References

- [MPFR Library](https://www.mpfr.org/) - Multi-Precision Floating-Point Reliable Library
- [Go big.Float](https://pkg.go.dev/math/big#Float) - Standard library arbitrary-precision floating-point
- [Awesome Go](https://github.com/avelino/awesome-go) - Curated list of Go frameworks and libraries

## Related Projects

- [gmp](https://github.com/ncw/gmp) - Go wrapper for GMP
- [decimal](https://github.com/shopspring/decimal) - Arbitrary-precision decimal arithmetic
- [apd](https://github.com/cockroachdb/apd) - Arbitrary-precision decimal library

