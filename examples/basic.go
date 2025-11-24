package main

import (
	"fmt"
	"math"

	"github.com/mshafiee/bigmath"
)

func main() {
	// Example 1: Basic arithmetic with arbitrary precision
	fmt.Println("=== Example 1: Basic Arithmetic ===")
	x := bigmath.NewBigFloat(3.141592653589793, 256)
	y := bigmath.NewBigFloat(2.718281828459045, 256)
	
	sum := new(bigmath.BigFloat).SetPrec(256).Add(x, y)
	product := new(bigmath.BigFloat).SetPrec(256).Mul(x, y)
	
	xVal, _ := x.Float64()
	yVal, _ := y.Float64()
	sumVal, _ := sum.Float64()
	productVal, _ := product.Float64()
	
	fmt.Printf("x = %g\n", xVal)
	fmt.Printf("y = %g\n", yVal)
	fmt.Printf("x + y = %g\n", sumVal)
	fmt.Printf("x * y = %g\n", productVal)
	fmt.Println()

	// Example 2: Trigonometric functions
	fmt.Println("=== Example 2: Trigonometric Functions ===")
	angle := bigmath.NewBigFloat(math.Pi/4, 256)
	sinVal := bigmath.BigSin(angle, 256)
	cosVal := bigmath.BigCos(angle, 256)
	tanVal := bigmath.BigTan(angle, 256)
	
	sinValF, _ := sinVal.Float64()
	cosValF, _ := cosVal.Float64()
	tanValF, _ := tanVal.Float64()
	
	fmt.Printf("sin(π/4) = %g\n", sinValF)
	fmt.Printf("cos(π/4) = %g\n", cosValF)
	fmt.Printf("tan(π/4) = %g\n", tanValF)
	fmt.Println()

	// Example 3: Constants
	fmt.Println("=== Example 3: Mathematical Constants ===")
	pi := bigmath.BigPI(256)
	e := bigmath.BigE(256)
	twoPi := bigmath.BigTwoPI(256)
	
	piVal, _ := pi.Float64()
	eVal, _ := e.Float64()
	twoPiVal, _ := twoPi.Float64()
	
	fmt.Printf("π = %g\n", piVal)
	fmt.Printf("e = %g\n", eVal)
	fmt.Printf("2π = %g\n", twoPiVal)
	fmt.Println()

	// Example 4: Vector operations
	fmt.Println("=== Example 4: Vector Operations ===")
	v1 := bigmath.NewBigVec3(1.0, 2.0, 3.0, 256)
	v2 := bigmath.NewBigVec3(4.0, 5.0, 6.0, 256)
	
	sumVec := bigmath.BigVec3Add(v1, v2, 256)
	dotProd := bigmath.BigVec3Dot(v1, v2, 256)
	magnitude := bigmath.BigVec3Magnitude(v1, 256)
	
	sumX, _ := sumVec.X.Float64()
	sumY, _ := sumVec.Y.Float64()
	sumZ, _ := sumVec.Z.Float64()
	dotProdVal, _ := dotProd.Float64()
	magVal, _ := magnitude.Float64()
	
	fmt.Printf("v1 = (1, 2, 3)\n")
	fmt.Printf("v2 = (4, 5, 6)\n")
	fmt.Printf("v1 + v2 = (%g, %g, %g)\n", sumX, sumY, sumZ)
	fmt.Printf("v1 · v2 = %g\n", dotProdVal)
	fmt.Printf("|v1| = %g\n", magVal)
	fmt.Println()

	// Example 5: Exponential and logarithmic functions
	fmt.Println("=== Example 5: Exponential and Logarithmic ===")
	val := bigmath.NewBigFloat(2.0, 256)
	expVal := bigmath.BigExp(val, 256)
	logVal := bigmath.BigLog(val, 256)
	powVal := bigmath.BigPow(val, bigmath.NewBigFloat(3.0, 256), 256)
	
	expValF, _ := expVal.Float64()
	logValF, _ := logVal.Float64()
	powValF, _ := powVal.Float64()
	
	fmt.Printf("exp(2) = %g\n", expValF)
	fmt.Printf("ln(2) = %g\n", logValF)
	fmt.Printf("2^3 = %g\n", powValF)
	fmt.Println()

	// Example 6: Square root
	fmt.Println("=== Example 6: Square Root ===")
	sqVal := bigmath.NewBigFloat(2.0, 256)
	sqrtVal := bigmath.BigSqrt(sqVal, 256)
	sqrtValF, _ := sqrtVal.Float64()
	fmt.Printf("√2 = %g\n", sqrtValF)
	fmt.Println()

	// Example 7: Rounding modes
	fmt.Println("=== Example 7: Rounding Modes ===")
	value := bigmath.NewBigFloat(1.5, 256)
	
	roundedNearest, _ := bigmath.Round(value, 64, bigmath.ToNearest)
	roundedZero, _ := bigmath.Round(value, 64, bigmath.ToZero)
	roundedPosInf, _ := bigmath.Round(value, 64, bigmath.ToPositiveInf)
	roundedNegInf, _ := bigmath.Round(value, 64, bigmath.ToNegativeInf)
	
	valueF, _ := value.Float64()
	roundedNearestF, _ := roundedNearest.Float64()
	roundedZeroF, _ := roundedZero.Float64()
	roundedPosInfF, _ := roundedPosInf.Float64()
	roundedNegInfF, _ := roundedNegInf.Float64()
	
	fmt.Printf("Original: %g\n", valueF)
	fmt.Printf("Round to nearest: %g\n", roundedNearestF)
	fmt.Printf("Round toward zero: %g\n", roundedZeroF)
	fmt.Printf("Round toward +∞: %g\n", roundedPosInfF)
	fmt.Printf("Round toward -∞: %g\n", roundedNegInfF)
}

