// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

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
	fmt.Println()

	// Example 8: Basic math utilities
	fmt.Println("=== Example 8: Basic Math Utilities ===")
	num := bigmath.NewBigFloat(3.7, 256)
	floorVal := bigmath.BigFloor(num, 256)
	ceilVal := bigmath.BigCeil(num, 256)
	truncVal := bigmath.BigTrunc(num, 256)

	floorValF, _ := floorVal.Float64()
	ceilValF, _ := ceilVal.Float64()
	truncValF, _ := truncVal.Float64()

	fmt.Printf("floor(3.7) = %g\n", floorValF)
	fmt.Printf("ceil(3.7) = %g\n", ceilValF)
	fmt.Printf("trunc(3.7) = %g\n", truncValF)
	fmt.Println()

	// Example 9: Root functions
	fmt.Println("=== Example 9: Root Functions ===")
	cubeVal := bigmath.NewBigFloat(8.0, 256)
	cbrtVal := bigmath.BigCbrt(cubeVal, 256)
	fourthRoot := bigmath.BigRoot(bigmath.NewBigFloat(4.0, 256), bigmath.NewBigFloat(16.0, 256), 256)

	cbrtValF, _ := cbrtVal.Float64()
	fourthRootF, _ := fourthRoot.Float64()

	fmt.Printf("cbrt(8) = %g\n", cbrtValF)
	fmt.Printf("16^(1/4) = %g\n", fourthRootF)
	fmt.Println()

	// Example 10: Advanced vector operations
	fmt.Println("=== Example 10: Advanced Vector Operations ===")
	vec1 := bigmath.NewBigVec3(1.0, 0.0, 0.0, 256)
	vec2 := bigmath.NewBigVec3(0.0, 1.0, 0.0, 256)

	cross := bigmath.BigVec3Cross(vec1, vec2, 256)
	normalized := bigmath.BigVec3Normalize(bigmath.NewBigVec3(3.0, 4.0, 0.0, 256), 256)
	angleVal := bigmath.BigVec3Angle(vec1, vec2, 256)

	crossF := cross.ToFloat64()
	normalizedF := normalized.ToFloat64()
	angleF, _ := angleVal.Float64()

	fmt.Printf("(1,0,0) × (0,1,0) = (%g, %g, %g)\n", crossF[0], crossF[1], crossF[2])
	fmt.Printf("normalize(3,4,0) = (%g, %g, %g)\n", normalizedF[0], normalizedF[1], normalizedF[2])
	fmt.Printf("angle between (1,0,0) and (0,1,0) = %g radians\n", angleF)
	fmt.Println()

	// Example 11: Advanced matrix operations
	fmt.Println("=== Example 11: Advanced Matrix Operations ===")
	mat := bigmath.NewIdentityMatrix(256)
	det := bigmath.BigMatDet(mat, 256)
	inv, err := bigmath.BigMatInverse(mat, 256)
	if err != nil {
		fmt.Printf("Error computing inverse: %v\n", err)
		return
	}

	detF, _ := det.Float64()
	invF, _ := inv.M[0][0].Float64()

	fmt.Printf("det(identity) = %g\n", detF)
	fmt.Printf("inv(identity)[0][0] = %g\n", invF)
	fmt.Println()

	// Example 12: Extended constants
	fmt.Println("=== Example 12: Extended Constants ===")
	phi := bigmath.BigPhi(256)
	sqrt2 := bigmath.BigSqrt2(256)
	ln10 := bigmath.BigLn10(256)

	phiF, _ := phi.Float64()
	sqrt2F, _ := sqrt2.Float64()
	ln10F, _ := ln10.Float64()

	fmt.Printf("φ (golden ratio) = %g\n", phiF)
	fmt.Printf("√2 = %g\n", sqrt2F)
	fmt.Printf("ln(10) = %g\n", ln10F)
	fmt.Println()

	// Example 13: Combinatorics
	fmt.Println("=== Example 13: Combinatorics ===")
	fact5 := bigmath.BigFactorial(5, 256)
	binom := bigmath.BigBinomial(10, 3, 256)

	fact5F, _ := fact5.Float64()
	binomF, _ := binom.Float64()

	fmt.Printf("5! = %g\n", fact5F)
	fmt.Printf("C(10,3) = %g\n", binomF)
	fmt.Println()

	// Example 14: Advanced logarithmic functions
	fmt.Println("=== Example 14: Advanced Logarithmic Functions ===")
	smallX := bigmath.NewBigFloat(0.001, 256)
	log1pVal := bigmath.BigLog1p(smallX, 256)
	logbVal := bigmath.BigLogb(bigmath.NewBigFloat(8.0, 256), bigmath.NewBigFloat(2.0, 256), 256)

	log1pValF, _ := log1pVal.Float64()
	logbValF, _ := logbVal.Float64()

	fmt.Printf("log1p(0.001) = %g\n", log1pValF)
	fmt.Printf("log₂(8) = %g\n", logbValF)
}
