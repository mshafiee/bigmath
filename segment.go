// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

import (
	"fmt"
	"math"
)

// SegmentInfoBig holds segment information with arbitrary precision
type SegmentInfoBig struct {
	SegmentStart *BigFloat
	SegmentEnd   *BigFloat
	SegmentSize  *BigFloat
	ElemEpoch    *BigFloat
	Qrot         *BigFloat
	DQrot        *BigFloat
	Prot         *BigFloat
	DProt        *BigFloat
	Peri         *BigFloat
	DPeri        *BigFloat
	RefEllipse   []*BigFloat
	Body         int
	NumCoeffs    int
	Flags        int
}

// reduceModulo2Pi reduces a BigFloat value modulo 2π
func reduceModulo2Pi(val *BigFloat, prec uint) *BigFloat {
	twoPi := BigTwoPI(prec)
	nRotations := new(BigFloat).SetPrec(prec).Quo(val, twoPi)
	nInt, _ := nRotations.Int(nil)
	nBig := new(BigFloat).SetPrec(prec).SetInt(nInt)
	temp := new(BigFloat).SetPrec(prec).Mul(nBig, twoPi)
	result := new(BigFloat).SetPrec(prec).Sub(val, temp)
	return result
}

// calculateOrbitalParametersMoon calculates qav and pav for Moon
func calculateOrbitalParametersMoon(tdiff *BigFloat, segInfo *SegmentInfoBig, prec uint) (qav, pav *BigFloat) {
	dn := new(BigFloat).SetPrec(prec)
	dn.Mul(tdiff, segInfo.DProt)
	dn.Add(segInfo.Prot, dn)
	dn = reduceModulo2Pi(dn, prec)

	qrotPlusTdiff := new(BigFloat).SetPrec(prec)
	qrotPlusTdiff.Mul(tdiff, segInfo.DQrot)
	qrotPlusTdiff.Add(segInfo.Qrot, qrotPlusTdiff)

	cosdn := BigCos(dn, prec)
	sindn := BigSin(dn, prec)

	qav = new(BigFloat).SetPrec(prec).Mul(qrotPlusTdiff, cosdn)
	pav = new(BigFloat).SetPrec(prec).Mul(qrotPlusTdiff, sindn)
	return qav, pav
}

// calculateOrbitalParametersPlanet calculates qav and pav for planets
func calculateOrbitalParametersPlanet(tdiff *BigFloat, segInfo *SegmentInfoBig, prec uint) (qav, pav *BigFloat) {
	qav = new(BigFloat).SetPrec(prec)
	qav.Mul(tdiff, segInfo.DQrot)
	qav.Add(segInfo.Qrot, qav)

	pav = new(BigFloat).SetPrec(prec)
	pav.Mul(tdiff, segInfo.DProt)
	pav.Add(segInfo.Prot, pav)

	qavF, _ := qav.Float64()
	pavF, _ := pav.Float64()
	tdiffF, _ := tdiff.Float64()
	fmt.Printf("[BIGFLOAT-QP] Body=%d qav=%.15e pav=%.15e tdiff=%.15e\n", segInfo.Body, qavF, pavF, tdiffF)
	return qav, pav
}

// addReferenceEllipse adds reference ellipse to coefficients if flag is set
func addReferenceEllipse(x [][3]*BigFloat, segInfo *SegmentInfoBig, tdiff *BigFloat, numCoeffs int, prec uint) {
	const SegFlagEllipse = 0x2
	if (segInfo.Flags & SegFlagEllipse) == 0 {
		fmt.Printf("[BIGFLOAT-ELLIPSE] Body=%d does NOT have SegFlagEllipse (Flags=0x%x)\n",
			segInfo.Body, segInfo.Flags)
		return
	}

	omtildF, _ := segInfo.Peri.Float64()
	fmt.Printf("[BIGFLOAT-ELLIPSE] Body=%d has SegFlagEllipse set, Peri=%.15e\n", segInfo.Body, omtildF)
	fmt.Printf("[BIGFLOAT-ELLIPSE] RefEllipse length=%d, need=%d\n", len(segInfo.RefEllipse), 2*numCoeffs)

	omtild := new(BigFloat).SetPrec(prec)
	omtild.Mul(tdiff, segInfo.DPeri)
	omtild.Add(segInfo.Peri, omtild)
	omtild = reduceModulo2Pi(omtild, prec)

	com := BigCos(omtild, prec)
	som := BigSin(omtild, prec)

	if len(segInfo.RefEllipse) < 2*numCoeffs {
		fmt.Printf("[BIGFLOAT-ELLIPSE] RefEllipse length=%d insufficient (need %d)\n",
			len(segInfo.RefEllipse), 2*numCoeffs)
		return
	}

	refepxF, _ := segInfo.RefEllipse[0].Float64()
	refepyF, _ := segInfo.RefEllipse[numCoeffs].Float64()
	comF, _ := com.Float64()
	somF, _ := som.Float64()
	fmt.Printf("[BIGFLOAT-ELLIPSE] Applying reference ellipse, refepx[0]=%.15e refepy[0]=%.15e\n", refepxF, refepyF)
	fmt.Printf("[BIGFLOAT-ELLIPSE] com=%.15e som=%.15e\n", comF, somF)

	x0BeforeF, _ := x[0][0].Float64()
	for i := 0; i < numCoeffs; i++ {
		refepx := segInfo.RefEllipse[i]
		refepy := segInfo.RefEllipse[i+numCoeffs]

		temp1 := new(BigFloat).SetPrec(prec).Mul(com, refepx)
		temp2 := new(BigFloat).SetPrec(prec).Mul(som, refepy)
		x[i][0].Add(x[i][0], temp1)
		x[i][0].Sub(x[i][0], temp2)

		temp1 = new(BigFloat).SetPrec(prec).Mul(com, refepy)
		temp2 = new(BigFloat).SetPrec(prec).Mul(som, refepx)
		x[i][1].Add(x[i][1], temp1)
		x[i][1].Add(x[i][1], temp2)
	}
	x0AfterF, _ := x[0][0].Float64()
	fmt.Printf("[BIGFLOAT-ELLIPSE] x[0][0]: before=%.15e after=%.15e (delta=%.15e)\n",
		x0BeforeF, x0AfterF, x0AfterF-x0BeforeF)

	fmt.Printf("[BIGFLOAT-ELLIPSE] First 10 X coeffs after ref ellipse:")
	for i := 0; i < 10 && i < numCoeffs; i++ {
		coeffF, _ := x[i][0].Float64()
		fmt.Printf(" %.15e", coeffF)
	}
	fmt.Printf("\n")
}

// constructRotationMatrix constructs the rotation matrix basis vectors uix, uiy, uiz
func constructRotationMatrix(qav, pav *BigFloat, segInfo *SegmentInfoBig, prec uint) (uix, uiy, uiz [3]*BigFloat) {
	one := NewBigFloat(1.0, prec)
	qavSq := new(BigFloat).SetPrec(prec).Mul(qav, qav)
	pavSq := new(BigFloat).SetPrec(prec).Mul(pav, pav)
	denom := new(BigFloat).SetPrec(prec).Add(one, qavSq)
	denom.Add(denom, pavSq)
	cosih2 := new(BigFloat).SetPrec(prec)
	cosih2.Quo(one, denom)

	cosih2F, _ := cosih2.Float64()
	fmt.Printf("[BIGFLOAT-COSIH2] Body=%d cosih2=%.15e\n", segInfo.Body, cosih2F)

	two := NewBigFloat(2.0, prec)
	temp := new(BigFloat).SetPrec(prec)

	// uix = [(1 + qav² - pav²) * cosih2, 2*qav*pav*cosih2, -2*pav*cosih2]
	uix = [3]*BigFloat{
		new(BigFloat).SetPrec(prec),
		new(BigFloat).SetPrec(prec),
		new(BigFloat).SetPrec(prec),
	}
	temp.Add(one, qavSq)
	temp.Sub(temp, pavSq)
	uix[0].Mul(temp, cosih2)

	temp.Mul(two, qav)
	temp.Mul(temp, pav)
	temp.Mul(temp, cosih2)
	uix[1].Set(temp)

	temp.Mul(two, pav)
	temp.Mul(temp, cosih2)
	uix[2].Neg(temp)

	// uiy = [2*qav*pav*cosih2, (1 - qav² + pav²) * cosih2, 2*qav*cosih2]
	uiy = [3]*BigFloat{
		new(BigFloat).SetPrec(prec),
		new(BigFloat).SetPrec(prec),
		new(BigFloat).SetPrec(prec),
	}
	temp.Mul(two, qav)
	temp.Mul(temp, pav)
	temp.Mul(temp, cosih2)
	uiy[0].Set(temp)

	temp.Sub(one, qavSq)
	temp.Add(temp, pavSq)
	uiy[1].Mul(temp, cosih2)

	temp.Mul(two, qav)
	temp.Mul(temp, cosih2)
	uiy[2].Set(temp)

	// uiz = [2*pav*cosih2, -2*qav*cosih2, (1 - qav² - pav²) * cosih2]
	uiz = [3]*BigFloat{
		new(BigFloat).SetPrec(prec),
		new(BigFloat).SetPrec(prec),
		new(BigFloat).SetPrec(prec),
	}
	temp.Mul(two, pav)
	temp.Mul(temp, cosih2)
	uiz[0].Set(temp)

	temp.Mul(two, qav)
	temp.Mul(temp, cosih2)
	uiz[1].Neg(temp)

	temp.Sub(one, qavSq)
	temp.Sub(temp, pavSq)
	uiz[2].Mul(temp, cosih2)

	uix0F, _ := uix[0].Float64()
	uix1F, _ := uix[1].Float64()
	uix2F, _ := uix[2].Float64()
	fmt.Printf("[BIGFLOAT-UIX] Body=%d uix=[%.15e, %.15e, %.15e]\n", segInfo.Body, uix0F, uix1F, uix2F)

	return uix, uiy, uiz
}

// rotateCoefficients rotates coefficients using the rotation matrix and tracks neval
func rotateCoefficients(x [][3]*BigFloat, uix, uiy, uiz [3]*BigFloat, segInfo *SegmentInfoBig, isMoon bool, numCoeffs int, prec uint) int {
	var threshold *BigFloat
	if segInfo.Body == 10 {
		threshold = NewBigFloat(3.67e-9, prec)
	} else {
		threshold = NewBigFloat(1e-14, prec)
	}

	neval := 0
	for i := 0; i < numCoeffs; i++ {
		xrot := new(BigFloat).SetPrec(prec)
		temp1 := new(BigFloat).SetPrec(prec).Mul(x[i][0], uix[0])
		temp2 := new(BigFloat).SetPrec(prec).Mul(x[i][1], uiy[0])
		temp3 := new(BigFloat).SetPrec(prec).Mul(x[i][2], uiz[0])
		xrot.Add(temp1, temp2)
		xrot.Add(xrot, temp3)

		yrot := new(BigFloat).SetPrec(prec)
		temp1 = new(BigFloat).SetPrec(prec).Mul(x[i][0], uix[1])
		temp2 = new(BigFloat).SetPrec(prec).Mul(x[i][1], uiy[1])
		temp3 = new(BigFloat).SetPrec(prec).Mul(x[i][2], uiz[1])
		yrot.Add(temp1, temp2)
		yrot.Add(yrot, temp3)

		zrot := new(BigFloat).SetPrec(prec)
		temp1 = new(BigFloat).SetPrec(prec).Mul(x[i][0], uix[2])
		temp2 = new(BigFloat).SetPrec(prec).Mul(x[i][1], uiy[2])
		temp3 = new(BigFloat).SetPrec(prec).Mul(x[i][2], uiz[2])
		zrot.Add(temp1, temp2)
		zrot.Add(zrot, temp3)

		xrotF, _ := xrot.Float64()
		yrotF, _ := yrot.Float64()
		zrotF, _ := zrot.Float64()
		magnitudeF := math.Abs(xrotF) + math.Abs(yrotF) + math.Abs(zrotF)
		thresholdF, _ := threshold.Float64()

		if magnitudeF >= thresholdF {
			neval = i
		}

		if isMoon {
			seps2000 := NewBigFloat(0.39777715572793088, prec)
			ceps2000 := NewBigFloat(0.91748206215761929, prec)

			yNew := new(BigFloat).SetPrec(prec)
			yNew.Mul(ceps2000, yrot)
			temp := new(BigFloat).SetPrec(prec).Mul(seps2000, zrot)
			yNew.Sub(yNew, temp)

			zNew := new(BigFloat).SetPrec(prec)
			zNew.Mul(seps2000, yrot)
			temp = new(BigFloat).SetPrec(prec).Mul(ceps2000, zrot)
			zNew.Add(zNew, temp)

			x[i][0] = xrot
			x[i][1] = yNew
			x[i][2] = zNew
		} else {
			x[i][0] = xrot
			x[i][1] = yrot
			x[i][2] = zrot
		}
	}
	return neval
}

// RotateCoeffsToJ2000Big rotates Chebyshev coefficients from orbital plane to equatorial J2000
// using arbitrary precision to eliminate all rounding errors.
// This is the BigFloat version of RotateCoeffsToJ2000() from segment_reader.go
func RotateCoeffsToJ2000Big(coeffs []*BigFloat, segInfo *SegmentInfoBig, isMoon bool, prec uint) (result []*BigFloat, neval int) {
	if prec == 0 {
		prec = DefaultPrecision
	}

	numCoeffs := segInfo.NumCoeffs

	// Time at middle of segment
	segStart := segInfo.SegmentStart
	segEnd := segInfo.SegmentEnd
	t := new(BigFloat).SetPrec(prec)
	t.Add(segStart, segEnd)
	t.Quo(t, NewBigFloat(2.0, prec))

	tdiff := new(BigFloat).SetPrec(prec)
	tdiff.Sub(t, segInfo.ElemEpoch)
	tdiff.Quo(tdiff, BigJulianMillennium(prec))

	if segInfo.Body == 9 { // Pluto
		tF, _ := t.Float64()
		telemF, _ := segInfo.ElemEpoch.Float64()
		tdiffF, _ := tdiff.Float64()
		fmt.Printf("[DEBUG-PLUTO] t=%.6f telem=%.6f tdiff=%.9e\n", tF, telemF, tdiffF)
	}

	// Calculate orbital parameters
	var qav, pav *BigFloat
	if isMoon {
		qav, pav = calculateOrbitalParametersMoon(tdiff, segInfo, prec)
	} else {
		qav, pav = calculateOrbitalParametersPlanet(tdiff, segInfo, prec)
	}

	// Copy coefficients to working array (3 sets of numCoeffs: X, Y, Z)
	x := make([][3]*BigFloat, numCoeffs)
	for i := 0; i < numCoeffs; i++ {
		x[i][0] = new(BigFloat).SetPrec(prec).Set(coeffs[i])             // X coeffs
		x[i][1] = new(BigFloat).SetPrec(prec).Set(coeffs[i+numCoeffs])   // Y coeffs
		x[i][2] = new(BigFloat).SetPrec(prec).Set(coeffs[i+2*numCoeffs]) // Z coeffs
	}

	// Add reference ellipse if flag is set
	addReferenceEllipse(x, segInfo, tdiff, numCoeffs, prec)

	// Construct rotation matrix basis vectors
	uix, uiy, uiz := constructRotationMatrix(qav, pav, segInfo, prec)

	// Rotate coefficients to actual orientation in space
	neval = rotateCoefficients(x, uix, uiy, uiz, segInfo, isMoon, numCoeffs, prec)

	// Write rotated coefficients back
	result = make([]*BigFloat, 3*numCoeffs)
	for i := 0; i < numCoeffs; i++ {
		result[i] = x[i][0]
		result[i+numCoeffs] = x[i][1]
		result[i+2*numCoeffs] = x[i][2]
	}

	// Debug: Print first 10 coefficients after rotation
	fmt.Printf("[BIGFLOAT-ROTATED] Body=%d First 10 X coeffs after rotation:", segInfo.Body)
	for i := 0; i < 10 && i < numCoeffs; i++ {
		coeffF, _ := result[i].Float64()
		fmt.Printf(" %.15e", coeffF)
	}
	fmt.Printf("\n")
	neval++
	fmt.Printf("[BIGFLOAT-ROTATED] Body=%d neval=%d\n", segInfo.Body, neval)

	return
}

// EvaluateChebyshevBig evaluates Chebyshev polynomial with arbitrary precision
// This is the BigFloat version of swi_echeb()
func EvaluateChebyshevBig(t *BigFloat, c []*BigFloat, neval int, prec uint) *BigFloat {
	return getDispatcher().EvaluateChebyshevBigImpl(t, c, neval, prec)
}

// EvaluateChebyshevDerivativeBig evaluates derivative of Chebyshev polynomial
// This is the BigFloat version of swi_edcheb()
func EvaluateChebyshevDerivativeBig(t *BigFloat, c []*BigFloat, neval int, prec uint) *BigFloat {
	return getDispatcher().EvaluateChebyshevDerivativeBigImpl(t, c, neval, prec)
}

// EvaluateSegmentBig evaluates segment coefficients to get position and velocity
func EvaluateSegmentBig(tjd *BigFloat, coeffs []*BigFloat, segStart, segEnd *BigFloat, neval int, prec uint) *BigVec6 {
	if prec == 0 {
		prec = DefaultPrecision
	}

	numCoeffs := len(coeffs) / 3

	// Normalize time to [-1, 1] range
	// t = 2 * (tjd - segStart) / (segEnd - segStart) - 1
	segSize := new(BigFloat).SetPrec(prec).Sub(segEnd, segStart)
	tOffset := new(BigFloat).SetPrec(prec).Sub(tjd, segStart)
	t := new(BigFloat).SetPrec(prec)
	t.Quo(tOffset, segSize)
	t.Mul(t, NewBigFloat(2.0, prec))
	t.Sub(t, NewBigFloat(1.0, prec))

	// DEBUG: Log segment evaluation details
	tjdF, _ := tjd.Float64()
	segStartF, _ := segStart.Float64()
	segEndF, _ := segEnd.Float64()
	segSizeF, _ := segSize.Float64()
	tF, _ := t.Float64()
	fmt.Printf("[EVAL-SEG] tjd=%.6f segStart=%.6f segEnd=%.6f segSize=%.6f\n", tjdF, segStartF, segEndF, segSizeF)
	fmt.Printf("[EVAL-SEG] t_normalized=%.15e (should be in [-1,1])\n", tF)

	// Evaluate position for X, Y, Z
	xCoeffs := coeffs[:numCoeffs]
	yCoeffs := coeffs[numCoeffs : 2*numCoeffs]
	zCoeffs := coeffs[2*numCoeffs:]

	x := EvaluateChebyshevBig(t, xCoeffs, neval, prec)
	y := EvaluateChebyshevBig(t, yCoeffs, neval, prec)
	z := EvaluateChebyshevBig(t, zCoeffs, neval, prec)

	// Evaluate velocity (derivative) for VX, VY, VZ
	// Velocity needs to be scaled by 2/segSize
	vx := EvaluateChebyshevDerivativeBig(t, xCoeffs, neval, prec)
	vy := EvaluateChebyshevDerivativeBig(t, yCoeffs, neval, prec)
	vz := EvaluateChebyshevDerivativeBig(t, zCoeffs, neval, prec)

	// Scale velocity: v = dpos/dt * (2/segSize)
	two := NewBigFloat(2.0, prec)
	velocityScale := new(BigFloat).SetPrec(prec).Quo(two, segSize)

	// DEBUG
	ssFloat, _ := segSize.Float64()
	vsFloat, _ := velocityScale.Float64()
	vxFloat, _ := vx.Float64()
	fmt.Printf("[SEGMENT] segSize=%.6f scale=%.6e raw_vx=%.6e\n", ssFloat, vsFloat, vxFloat)

	vx.Mul(vx, velocityScale)
	vy.Mul(vy, velocityScale)
	vz.Mul(vz, velocityScale)

	return &BigVec6{
		X:  x,
		Y:  y,
		Z:  z,
		VX: vx,
		VY: vy,
		VZ: vz,
	}
}

// ConvertToBigFloatCoeffs converts float64 coefficients to BigFloat
func ConvertToBigFloatCoeffs(coeffsFloat64 []float64, prec uint) []*BigFloat {
	result := make([]*BigFloat, len(coeffsFloat64))
	for i, v := range coeffsFloat64 {
		result[i] = NewBigFloat(v, prec)
	}
	return result
}

// DebugPrintBigVec6 prints a BigVec6 for debugging
func DebugPrintBigVec6(label string, v *BigVec6) {
	x, _ := v.X.Float64()
	y, _ := v.Y.Float64()
	z, _ := v.Z.Float64()
	vx, _ := v.VX.Float64()
	vy, _ := v.VY.Float64()
	vz, _ := v.VZ.Float64()

	fmt.Printf("[%s] pos=[%.15e, %.15e, %.15e] vel=[%.15e, %.15e, %.15e]\n",
		label, x, y, z, vx, vy, vz)
}
