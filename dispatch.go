// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

import (
	"sync"
)

// Function pointer types for dispatched functions
type (
	// Vector operations
	bigVec3AddFunc func(v1, v2 *BigVec3, prec uint) *BigVec3
	bigVec3SubFunc func(v1, v2 *BigVec3, prec uint) *BigVec3
	bigVec3MulFunc func(v *BigVec3, scalar *BigFloat, prec uint) *BigVec3
	bigVec3DotFunc func(v1, v2 *BigVec3, prec uint) *BigFloat

	// BigVec6 operations
	bigVec6AddFunc       func(v1, v2 *BigVec6, prec uint) *BigVec6
	bigVec6SubFunc       func(v1, v2 *BigVec6, prec uint) *BigVec6
	bigVec6NegateFunc    func(v *BigVec6, prec uint) *BigVec6
	bigVec6MagnitudeFunc func(v *BigVec6, prec uint) *BigFloat

	// Matrix operations
	bigMatMulFunc func(m *BigMatrix3x3, v *BigVec3, prec uint) *BigVec3

	// Chebyshev evaluation
	evaluateChebyshevBigFunc           func(t *BigFloat, c []*BigFloat, neval int, prec uint) *BigFloat
	evaluateChebyshevDerivativeBigFunc func(t *BigFloat, c []*BigFloat, neval int, prec uint) *BigFloat

	// Trigonometric functions
	bigSinFunc   func(x *BigFloat, prec uint) *BigFloat
	bigCosFunc   func(x *BigFloat, prec uint) *BigFloat
	bigTanFunc   func(x *BigFloat, prec uint) *BigFloat
	bigAtanFunc  func(x *BigFloat, prec uint) *BigFloat
	bigAsinFunc  func(x *BigFloat, prec uint) *BigFloat
	bigAcosFunc  func(x *BigFloat, prec uint) *BigFloat
	bigAtan2Func func(y, x *BigFloat, prec uint) *BigFloat

	// Exponential and logarithmic functions
	bigExpFunc func(x *BigFloat, prec uint) *BigFloat
	bigLogFunc func(x *BigFloat, prec uint) *BigFloat

	// Power function
	bigPowFunc func(x, y *BigFloat, prec uint) *BigFloat

	// Hyperbolic functions
	bigSinhFunc  func(x *BigFloat, prec uint) *BigFloat
	bigCoshFunc  func(x *BigFloat, prec uint) *BigFloat
	bigTanhFunc  func(x *BigFloat, prec uint) *BigFloat
	bigAsinhFunc func(x *BigFloat, prec uint) *BigFloat
	bigAcoshFunc func(x *BigFloat, prec uint) *BigFloat
	bigAtanhFunc func(x *BigFloat, prec uint) *BigFloat

	// Special functions
	bigGammaFunc   func(x *BigFloat, prec uint) *BigFloat
	bigErfFunc     func(x *BigFloat, prec uint) *BigFloat
	bigErfcFunc    func(x *BigFloat, prec uint) *BigFloat
	bigBesselJFunc func(n int, x *BigFloat, prec uint) *BigFloat
	bigBesselYFunc func(n int, x *BigFloat, prec uint) *BigFloat

	// Root functions
	bigCbrtFunc func(x *BigFloat, prec uint) *BigFloat
	bigRootFunc func(n, x *BigFloat, prec uint) *BigFloat

	// Basic operations
	bigFloorFunc func(x *BigFloat, prec uint) *BigFloat
	bigCeilFunc  func(x *BigFloat, prec uint) *BigFloat
	bigTruncFunc func(x *BigFloat, prec uint) *BigFloat
	bigModFunc   func(x, y *BigFloat, prec uint) *BigFloat
	bigRemFunc   func(x, y *BigFloat, prec uint) *BigFloat

	// Combinatorics
	bigFactorialFunc func(n int64, prec uint) *BigFloat
	bigBinomialFunc  func(n, k int64, prec uint) *BigFloat

	// Advanced vector operations
	bigVec3CrossFunc     func(v1, v2 *BigVec3, prec uint) *BigVec3
	bigVec3NormalizeFunc func(v *BigVec3, prec uint) *BigVec3
	bigVec3AngleFunc     func(v1, v2 *BigVec3, prec uint) *BigFloat
	bigVec3ProjectFunc   func(v1, v2 *BigVec3, prec uint) *BigVec3

	// Advanced matrix operations
	bigMatTransposeFunc func(m *BigMatrix3x3, prec uint) *BigMatrix3x3
	bigMatMulMatFunc    func(m1, m2 *BigMatrix3x3, prec uint) *BigMatrix3x3
	bigMatDetFunc       func(m *BigMatrix3x3, prec uint) *BigFloat
	bigMatInverseFunc   func(m *BigMatrix3x3, prec uint) (*BigMatrix3x3, error)
)

// Dispatcher holds function pointers selected at runtime
type Dispatcher struct {
	// Vector operations
	BigVec3AddImpl bigVec3AddFunc
	BigVec3SubImpl bigVec3SubFunc
	BigVec3MulImpl bigVec3MulFunc
	BigVec3DotImpl bigVec3DotFunc

	// BigVec6 operations
	BigVec6AddImpl       bigVec6AddFunc
	BigVec6SubImpl       bigVec6SubFunc
	BigVec6NegateImpl    bigVec6NegateFunc
	BigVec6MagnitudeImpl bigVec6MagnitudeFunc

	// Matrix operations
	BigMatMulImpl bigMatMulFunc

	// Chebyshev evaluation
	EvaluateChebyshevBigImpl           evaluateChebyshevBigFunc
	EvaluateChebyshevDerivativeBigImpl evaluateChebyshevDerivativeBigFunc

	// Trigonometric functions
	BigSinImpl   bigSinFunc
	BigCosImpl   bigCosFunc
	BigTanImpl   bigTanFunc
	BigAtanImpl  bigAtanFunc
	BigAsinImpl  bigAsinFunc
	BigAcosImpl  bigAcosFunc
	BigAtan2Impl bigAtan2Func

	// Exponential and logarithmic functions
	BigExpImpl bigExpFunc
	BigLogImpl bigLogFunc

	// Power function
	BigPowImpl bigPowFunc

	// Hyperbolic functions
	BigSinhImpl  bigSinhFunc
	BigCoshImpl  bigCoshFunc
	BigTanhImpl  bigTanhFunc
	BigAsinhImpl bigAsinhFunc
	BigAcoshImpl bigAcoshFunc
	BigAtanhImpl bigAtanhFunc

	// Special functions
	BigGammaImpl   bigGammaFunc
	BigErfImpl     bigErfFunc
	BigErfcImpl    bigErfcFunc
	BigBesselJImpl bigBesselJFunc
	BigBesselYImpl bigBesselYFunc

	// Root functions
	BigCbrtImpl bigCbrtFunc
	BigRootImpl bigRootFunc

	// Basic operations
	BigFloorImpl bigFloorFunc
	BigCeilImpl  bigCeilFunc
	BigTruncImpl bigTruncFunc
	BigModImpl   bigModFunc
	BigRemImpl   bigRemFunc

	// Combinatorics
	BigFactorialImpl bigFactorialFunc
	BigBinomialImpl  bigBinomialFunc

	// Advanced vector operations
	BigVec3CrossImpl     bigVec3CrossFunc
	BigVec3NormalizeImpl bigVec3NormalizeFunc
	BigVec3AngleImpl     bigVec3AngleFunc
	BigVec3ProjectImpl   bigVec3ProjectFunc

	// Advanced matrix operations
	BigMatTransposeImpl bigMatTransposeFunc
	BigMatMulMatImpl    bigMatMulMatFunc
	BigMatDetImpl       bigMatDetFunc
	BigMatInverseImpl   bigMatInverseFunc

	// CPU features used
	Features CPUFeatures
}

var (
	dispatcher     *Dispatcher
	dispatcherOnce sync.Once
)

// initDispatcher initializes the function dispatcher based on CPU capabilities
// The actual implementation selection is done in architecture-specific files
func initDispatcher() *Dispatcher {
	d := &Dispatcher{}
	d.Features = GetCPUFeatures()

	// Call architecture-specific initialization
	initDispatcherImpl(d)

	return d
}

// getDispatcher returns the initialized dispatcher (singleton)
func getDispatcher() *Dispatcher {
	dispatcherOnce.Do(func() {
		dispatcher = initDispatcher()
	})
	return dispatcher
}
