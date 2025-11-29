// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

// BigFloatMarshalJSON marshals a BigFloat to JSON
// Uses string representation for precision
func BigFloatMarshalJSON(x *BigFloat) ([]byte, error) {
	if x == nil {
		return []byte("null"), nil
	}
	return json.Marshal(x.Text('g', -1))
}

// BigFloatUnmarshalJSON unmarshals a BigFloat from JSON
func BigFloatUnmarshalJSON(data []byte, prec uint) (*BigFloat, error) {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}

	if prec == 0 {
		prec = DefaultPrecision
	}

	return NewBigFloatFromString(s, prec)
}

// MarshalJSON implements json.Marshaler for BigVec3
func (v *BigVec3) MarshalJSON() ([]byte, error) {
	if v == nil {
		return []byte("null"), nil
	}

	return json.Marshal([3]string{
		v.X.Text('g', -1),
		v.Y.Text('g', -1),
		v.Z.Text('g', -1),
	})
}

// UnmarshalJSON implements json.Unmarshaler for BigVec3
func (v *BigVec3) UnmarshalJSON(data []byte) error {
	if v == nil {
		return errors.New("cannot unmarshal into nil BigVec3")
	}

	var arr [3]string
	if err := json.Unmarshal(data, &arr); err != nil {
		return err
	}

	var prec uint = DefaultPrecision
	if v.X != nil {
		prec = v.X.Prec()
	}
	if prec == 0 {
		prec = DefaultPrecision
	}

	x, err := NewBigFloatFromString(arr[0], prec)
	if err != nil {
		return fmt.Errorf("invalid X component: %w", err)
	}

	y, err := NewBigFloatFromString(arr[1], prec)
	if err != nil {
		return fmt.Errorf("invalid Y component: %w", err)
	}

	z, err := NewBigFloatFromString(arr[2], prec)
	if err != nil {
		return fmt.Errorf("invalid Z component: %w", err)
	}

	v.X = x
	v.Y = y
	v.Z = z

	return nil
}

// MarshalJSON implements json.Marshaler for BigVec6
func (v *BigVec6) MarshalJSON() ([]byte, error) {
	if v == nil {
		return []byte("null"), nil
	}

	return json.Marshal([6]string{
		v.X.Text('g', -1),
		v.Y.Text('g', -1),
		v.Z.Text('g', -1),
		v.VX.Text('g', -1),
		v.VY.Text('g', -1),
		v.VZ.Text('g', -1),
	})
}

// UnmarshalJSON implements json.Unmarshaler for BigVec6
func (v *BigVec6) UnmarshalJSON(data []byte) error {
	if v == nil {
		return errors.New("cannot unmarshal into nil BigVec6")
	}

	var arr [6]string
	if err := json.Unmarshal(data, &arr); err != nil {
		return err
	}

	var prec uint = DefaultPrecision
	if v.X != nil {
		prec = v.X.Prec()
	}
	if prec == 0 {
		prec = DefaultPrecision
	}

	x, err := NewBigFloatFromString(arr[0], prec)
	if err != nil {
		return fmt.Errorf("invalid X component: %w", err)
	}

	y, err := NewBigFloatFromString(arr[1], prec)
	if err != nil {
		return fmt.Errorf("invalid Y component: %w", err)
	}

	z, err := NewBigFloatFromString(arr[2], prec)
	if err != nil {
		return fmt.Errorf("invalid Z component: %w", err)
	}

	vx, err := NewBigFloatFromString(arr[3], prec)
	if err != nil {
		return fmt.Errorf("invalid VX component: %w", err)
	}

	vy, err := NewBigFloatFromString(arr[4], prec)
	if err != nil {
		return fmt.Errorf("invalid VY component: %w", err)
	}

	vz, err := NewBigFloatFromString(arr[5], prec)
	if err != nil {
		return fmt.Errorf("invalid VZ component: %w", err)
	}

	v.X = x
	v.Y = y
	v.Z = z
	v.VX = vx
	v.VY = vy
	v.VZ = vz

	return nil
}

// MarshalJSON implements json.Marshaler for BigMatrix3x3
func (m *BigMatrix3x3) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	}

	matrix := [3][3]string{}
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			matrix[i][j] = m.M[i][j].Text('g', -1)
		}
	}

	return json.Marshal(matrix)
}

// UnmarshalJSON implements json.Unmarshaler for BigMatrix3x3
func (m *BigMatrix3x3) UnmarshalJSON(data []byte) error {
	if m == nil {
		return errors.New("cannot unmarshal into nil BigMatrix3x3")
	}

	var matrix [3][3]string
	if err := json.Unmarshal(data, &matrix); err != nil {
		return err
	}

	var prec uint = DefaultPrecision
	if m.M[0][0] != nil {
		prec = m.M[0][0].Prec()
	}
	if prec == 0 {
		prec = DefaultPrecision
	}

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			val, err := NewBigFloatFromString(matrix[i][j], prec)
			if err != nil {
				return fmt.Errorf("invalid element [%d][%d]: %w", i, j, err)
			}
			if m.M[i][j] == nil {
				m.M[i][j] = val
			} else {
				m.M[i][j].Set(val)
			}
		}
	}

	return nil
}

// ReadDoubleAsBigFloat reads 8 bytes from the reader and converts them directly to BigFloat
// without going through float64. This preserves the full 53-bit precision of IEEE 754 doubles.
//
// The bytes are interpreted as an IEEE 754 double precision floating-point number:
// - 1 sign bit
// - 11 exponent bits
// - 52 mantissa bits (with implicit leading 1)
//
// This function extracts the sign, exponent, and mantissa, then constructs a BigFloat
// directly from these components, avoiding the float64 intermediate conversion.
//
// Parameters:
//   - r: io.Reader to read 8 bytes from
//   - bigEndian: true for big-endian byte order, false for little-endian
//   - prec: BigFloat precision in bits (0 uses DefaultPrecision)
//
// Returns:
//   - *BigFloat: The converted value with full 53-bit precision
//   - error: Any error encountered during reading or conversion
//
// Example:
//
//	reader := bytes.NewReader(doubleBytes)
//	value, err := bigmath.ReadDoubleAsBigFloat(reader, false, 256)
//	if err != nil {
//	    return err
//	}
func ReadDoubleAsBigFloat(r io.Reader, bigEndian bool, prec uint) (*BigFloat, error) {
	// Use platform-specific optimized implementation when available
	// Falls back to generic implementation on unsupported platforms
	return readDoubleAsBigFloatImpl(r, bigEndian, prec)
}
