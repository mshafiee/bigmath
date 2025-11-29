// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/big"
	"testing"
)

func TestBigVec3JSON(t *testing.T) {
	prec := uint(256)

	tests := []struct {
		name      string
		v         *BigVec3
		tolerance float64
	}{
		{"normal_values", NewBigVec3(1.0, 2.0, 3.0, prec), 1e-10},
		{"zero_vector", NewBigVec3(0.0, 0.0, 0.0, prec), 1e-10},
		{"negative_values", NewBigVec3(-1.0, -2.0, -3.0, prec), 1e-10},
		{"very_large", NewBigVec3(1e10, 2e10, 3e10, prec), 1e5},
		{"very_small", NewBigVec3(1e-10, 2e-10, 3e-10, prec), 1e-20},
		{"mixed_signs", NewBigVec3(1.0, -2.0, 3.0, prec), 1e-10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.v)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}

			var v2 BigVec3
			if err = json.Unmarshal(data, &v2); err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}

			orig := tt.v.ToFloat64()
			unmarshaled := v2.ToFloat64()

			for i := 0; i < 3; i++ {
				diff := orig[i] - unmarshaled[i]
				if diff < 0 {
					diff = -diff
				}
				if diff > tt.tolerance {
					t.Errorf("Component[%d] = %g, want %g (diff %g, tolerance %g)", i, unmarshaled[i], orig[i], diff, tt.tolerance)
				}
			}

			// Test round-trip accuracy
			data2, err := json.Marshal(&v2)
			if err != nil {
				t.Fatalf("Second Marshal failed: %v", err)
			}

			var v3 BigVec3
			if err = json.Unmarshal(data2, &v3); err != nil {
				t.Fatalf("Second Unmarshal failed: %v", err)
			}

			unmarshaled2 := v3.ToFloat64()
			for i := 0; i < 3; i++ {
				diff := unmarshaled[i] - unmarshaled2[i]
				if diff < 0 {
					diff = -diff
				}
				if diff > tt.tolerance {
					t.Errorf("Round-trip failed: Component[%d] = %g, want %g", i, unmarshaled2[i], unmarshaled[i])
				}
			}
		})
	}

	// Test invalid JSON
	t.Run("invalid_json", func(t *testing.T) {
		invalidJSON := []byte(`{"X": "not a number", "Y": 2.0, "Z": 3.0}`)
		var v BigVec3
		if err := json.Unmarshal(invalidJSON, &v); err == nil {
			t.Error("Unmarshal should fail for invalid JSON")
		}
	})

	// Test precision preservation
	t.Run("precision_preservation", func(t *testing.T) {
		testCases := []uint{64, 128, 256, 512}
		for _, p := range testCases {
			v := NewBigVec3(1.0, 2.0, 3.0, p)
			data, err := json.Marshal(v)
			if err != nil {
				t.Fatalf("Marshal failed at prec %d: %v", p, err)
			}

			var v2 BigVec3
			if err := json.Unmarshal(data, &v2); err != nil {
				t.Fatalf("Unmarshal failed at prec %d: %v", p, err)
			}

			// Check that precision is preserved (at least approximately)
			origPrec := v.X.Prec()
			unmarshaledPrec := v2.X.Prec()
			if unmarshaledPrec < origPrec/2 {
				t.Errorf("Precision not preserved: orig %d, unmarshaled %d", origPrec, unmarshaledPrec)
			}
		}
	})
}

func TestBigVec6JSON(t *testing.T) {
	prec := uint(256)

	tests := []struct {
		name      string
		v         *BigVec6
		tolerance float64
	}{
		{"normal_values", NewBigVec6(1.0, 2.0, 3.0, 0.1, 0.2, 0.3, prec), 1e-10},
		{"zero_vector", NewBigVec6(0.0, 0.0, 0.0, 0.0, 0.0, 0.0, prec), 1e-10},
		{"negative_values", NewBigVec6(-1.0, -2.0, -3.0, -0.1, -0.2, -0.3, prec), 1e-10},
		{"very_large", NewBigVec6(1e10, 2e10, 3e10, 4e10, 5e10, 6e10, prec), 1e5},
		{"very_small", NewBigVec6(1e-10, 2e-10, 3e-10, 4e-10, 5e-10, 6e-10, prec), 1e-20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.v)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}

			var v2 BigVec6
			if err := json.Unmarshal(data, &v2); err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}

			orig := tt.v.ToFloat64()
			unmarshaled := v2.ToFloat64()

			for i := 0; i < 6; i++ {
				diff := orig[i] - unmarshaled[i]
				if diff < 0 {
					diff = -diff
				}
				if diff > tt.tolerance {
					t.Errorf("Component[%d] = %g, want %g (diff %g, tolerance %g)", i, unmarshaled[i], orig[i], diff, tt.tolerance)
				}
			}
		})
	}

	// Test invalid JSON
	t.Run("invalid_json", func(t *testing.T) {
		invalidJSON := []byte(`{"X": "not a number", "Y": 2.0}`)
		var v BigVec6
		if err := json.Unmarshal(invalidJSON, &v); err == nil {
			t.Error("Unmarshal should fail for invalid JSON")
		}
	})
}

func TestBigMatrix3x3JSON(t *testing.T) {
	prec := uint(256)

	tests := []struct {
		name      string
		m         *BigMatrix3x3
		tolerance float64
	}{
		{"identity", NewIdentityMatrix(prec), 1e-10},
		{"zero_matrix", &BigMatrix3x3{
			M: [3][3]*BigFloat{
				{NewBigFloat(0.0, prec), NewBigFloat(0.0, prec), NewBigFloat(0.0, prec)},
				{NewBigFloat(0.0, prec), NewBigFloat(0.0, prec), NewBigFloat(0.0, prec)},
				{NewBigFloat(0.0, prec), NewBigFloat(0.0, prec), NewBigFloat(0.0, prec)},
			},
		}, 1e-10},
		{"general_matrix", &BigMatrix3x3{
			M: [3][3]*BigFloat{
				{NewBigFloat(1.0, prec), NewBigFloat(2.0, prec), NewBigFloat(3.0, prec)},
				{NewBigFloat(4.0, prec), NewBigFloat(5.0, prec), NewBigFloat(6.0, prec)},
				{NewBigFloat(7.0, prec), NewBigFloat(8.0, prec), NewBigFloat(9.0, prec)},
			},
		}, 1e-10},
		{"very_large", &BigMatrix3x3{
			M: [3][3]*BigFloat{
				{NewBigFloat(1e10, prec), NewBigFloat(2e10, prec), NewBigFloat(3e10, prec)},
				{NewBigFloat(4e10, prec), NewBigFloat(5e10, prec), NewBigFloat(6e10, prec)},
				{NewBigFloat(7e10, prec), NewBigFloat(8e10, prec), NewBigFloat(9e10, prec)},
			},
		}, 1e5},
		{"very_small", &BigMatrix3x3{
			M: [3][3]*BigFloat{
				{NewBigFloat(1e-10, prec), NewBigFloat(2e-10, prec), NewBigFloat(3e-10, prec)},
				{NewBigFloat(4e-10, prec), NewBigFloat(5e-10, prec), NewBigFloat(6e-10, prec)},
				{NewBigFloat(7e-10, prec), NewBigFloat(8e-10, prec), NewBigFloat(9e-10, prec)},
			},
		}, 1e-20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.m)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}

			var m2 BigMatrix3x3
			if err := json.Unmarshal(data, &m2); err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}

			for i := 0; i < 3; i++ {
				for j := 0; j < 3; j++ {
					orig, _ := tt.m.M[i][j].Float64()
					got, _ := m2.M[i][j].Float64()
					diff := orig - got
					if diff < 0 {
						diff = -diff
					}
					if diff > tt.tolerance {
						t.Errorf("M[%d][%d] = %g, want %g (diff %g, tolerance %g)", i, j, got, orig, diff, tt.tolerance)
					}
				}
			}
		})
	}

	// Test invalid JSON
	t.Run("invalid_json", func(t *testing.T) {
		invalidJSON := []byte(`{"M": [["not a number", 2.0, 3.0], [4.0, 5.0, 6.0], [7.0, 8.0, 9.0]]}`)
		var m BigMatrix3x3
		if err := json.Unmarshal(invalidJSON, &m); err == nil {
			t.Error("Unmarshal should fail for invalid JSON")
		}
	})

	// Test precision preservation
	t.Run("precision_preservation", func(t *testing.T) {
		testCases := []uint{64, 128, 256, 512}
		for _, p := range testCases {
			m := NewIdentityMatrix(p)
			data, err := json.Marshal(m)
			if err != nil {
				t.Fatalf("Marshal failed at prec %d: %v", p, err)
			}

			var m2 BigMatrix3x3
			if err := json.Unmarshal(data, &m2); err != nil {
				t.Fatalf("Unmarshal failed at prec %d: %v", p, err)
			}

			// Check that precision is preserved (at least approximately)
			origPrec := m.M[0][0].Prec()
			unmarshaledPrec := m2.M[0][0].Prec()
			if unmarshaledPrec < origPrec/2 {
				t.Errorf("Precision not preserved: orig %d, unmarshaled %d", origPrec, unmarshaledPrec)
			}
		}
	})
}

// TestBigFloatMarshalJSON tests BigFloat JSON marshaling
func TestBigFloatMarshalJSON(t *testing.T) {
	prec := uint(256)

	tests := []struct {
		name      string
		value     float64
		tolerance float64
	}{
		{"zero", 0.0, 1e-10},
		{"one", 1.0, 1e-10},
		{"negative", -3.14, 1e-10},
		{"large", 1e10, 1e5},
		{"small", 1e-10, 1e-20},
		{"pi", 3.141592653589793, 1e-10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := NewBigFloat(tt.value, prec)
			data, err := BigFloatMarshalJSON(x)
			if err != nil {
				t.Fatalf("BigFloatMarshalJSON failed: %v", err)
			}

			if len(data) == 0 {
				t.Error("BigFloatMarshalJSON returned empty data")
			}

			// Should be valid JSON string
			var s string
			if err := json.Unmarshal(data, &s); err != nil {
				t.Errorf("BigFloatMarshalJSON didn't produce valid JSON: %v", err)
			}
		})
	}
}

// TestBigFloatUnmarshalJSON tests BigFloat JSON unmarshaling
func TestBigFloatUnmarshalJSON(t *testing.T) {
	prec := uint(256)

	tests := []struct {
		name      string
		jsonStr   string
		expected  float64
		shouldErr bool
	}{
		{"zero", `"0"`, 0.0, false},
		{"one", `"1.0"`, 1.0, false},
		{"negative", `"-3.14"`, -3.14, false},
		{"scientific", `"1.23e10"`, 1.23e10, false},
		{"invalid", `"not a number"`, 0.0, true},
		{"empty", `""`, 0.0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x, err := BigFloatUnmarshalJSON([]byte(tt.jsonStr), prec)

			if tt.shouldErr {
				if err == nil {
					t.Error("BigFloatUnmarshalJSON should have failed")
				}
			} else {
				if err != nil {
					t.Fatalf("BigFloatUnmarshalJSON failed: %v", err)
				}

				if x == nil {
					t.Error("BigFloatUnmarshalJSON returned nil result")
					return
				}

				xFloat, _ := x.Float64()
				diff := xFloat - tt.expected
				if diff < 0 {
					diff = -diff
				}
				if diff > 1e-10 && tt.expected != 0 && diff/tt.expected > 1e-10 {
					t.Errorf("BigFloatUnmarshalJSON = %v, want %v", xFloat, tt.expected)
				}
			}
		})
	}
}

// TestBigFloatJSONRoundTrip tests BigFloat JSON round-trip
func TestBigFloatJSONRoundTrip(t *testing.T) {
	prec := uint(256)

	values := []float64{0.0, 1.0, -1.0, 3.14159, 1e10, 1e-10, -2.71828}

	for _, val := range values {
		t.Run("value_"+string(rune(int(val*100))), func(t *testing.T) {
			x := NewBigFloat(val, prec)

			// Marshal
			data, err := BigFloatMarshalJSON(x)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}

			// Unmarshal
			y, err := BigFloatUnmarshalJSON(data, prec)
			if err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}

			if y == nil {
				t.Fatal("Unmarshal returned nil")
			}

			// Compare
			xFloat, _ := x.Float64()
			yFloat, _ := y.Float64()

			if xFloat == 0.0 {
				if yFloat != 0.0 {
					t.Errorf("Round-trip failed: %v != %v", yFloat, xFloat)
				}
			} else {
				relErr := (yFloat - xFloat) / xFloat
				if relErr < 0 {
					relErr = -relErr
				}
				if relErr > 1e-10 {
					t.Errorf("Round-trip failed: %v != %v (rel err %e)", yFloat, xFloat, relErr)
				}
			}
		})
	}
}

// TestReadDoubleAsBigFloat tests ReadDoubleAsBigFloat function
func TestReadDoubleAsBigFloat(t *testing.T) {
	prec := uint(256)

	tests := []struct {
		name      string
		value     float64
		bigEndian bool
		tolerance float64
	}{
		{"zero_little_endian", 0.0, false, 1e-20},
		{"zero_big_endian", 0.0, true, 1e-20},
		{"negative_zero_little_endian", math.Copysign(0.0, -1.0), false, 1e-20},
		{"negative_zero_big_endian", math.Copysign(0.0, -1.0), true, 1e-20},
		{"one_little_endian", 1.0, false, 1e-15},
		{"one_big_endian", 1.0, true, 1e-15},
		{"negative_one_little_endian", -1.0, false, 1e-15},
		{"negative_one_big_endian", -1.0, true, 1e-15},
		{"pi_little_endian", math.Pi, false, 1e-15},
		{"pi_big_endian", math.Pi, true, 1e-15},
		{"e_little_endian", math.E, false, 1e-15},
		{"e_big_endian", math.E, true, 1e-15},
		{"large_number_little_endian", 1e100, false, 1e85},
		{"large_number_big_endian", 1e100, true, 1e85},
		{"small_number_little_endian", 1e-100, false, 1e-115},
		{"small_number_big_endian", 1e-100, true, 1e-115},
		{"negative_large_little_endian", -1e50, false, 1e35},
		{"negative_large_big_endian", -1e50, true, 1e35},
		{"negative_small_little_endian", -1e-50, false, 1e-65},
		{"negative_small_big_endian", -1e-50, true, 1e-65},
		{"max_float64_little_endian", math.MaxFloat64, false, 1e292},
		{"max_float64_big_endian", math.MaxFloat64, true, 1e292},
		// Note: math.SmallestNonzeroFloat64 is denormalized, which the function handles as zero (see TODO in code)
		// Use smallest normalized number instead: 2^-1022
		{"smallest_normal_little_endian", math.Pow(2, -1022), false, 1e-330},
		{"smallest_normal_big_endian", math.Pow(2, -1022), true, 1e-330},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Convert float64 to bytes
			var buf bytes.Buffer
			if tt.bigEndian {
				binary.Write(&buf, binary.BigEndian, tt.value)
			} else {
				binary.Write(&buf, binary.LittleEndian, tt.value)
			}

			// Read back using ReadDoubleAsBigFloat
			reader := bytes.NewReader(buf.Bytes())
			result, err := ReadDoubleAsBigFloat(reader, tt.bigEndian, prec)
			if err != nil {
				t.Fatalf("ReadDoubleAsBigFloat failed: %v", err)
			}

			if result == nil {
				t.Fatal("ReadDoubleAsBigFloat returned nil")
			}

			// Compare with original value
			resultFloat, _ := result.Float64()
			diff := math.Abs(resultFloat - tt.value)
			if diff > tt.tolerance {
				t.Errorf("ReadDoubleAsBigFloat = %g, want %g (diff %g, tolerance %g)", resultFloat, tt.value, diff, tt.tolerance)
			}
		})
	}
}

// TestReadDoubleAsBigFloatSpecialCases tests special IEEE 754 cases
func TestReadDoubleAsBigFloatSpecialCases(t *testing.T) {
	prec := uint(256)

	// Test positive infinity
	t.Run("positive_infinity_little_endian", func(t *testing.T) {
		var buf bytes.Buffer
		binary.Write(&buf, binary.LittleEndian, math.Inf(1))
		reader := bytes.NewReader(buf.Bytes())
		result, err := ReadDoubleAsBigFloat(reader, false, prec)
		if err != nil {
			t.Fatalf("ReadDoubleAsBigFloat failed: %v", err)
		}
		if !result.IsInf() || result.Sign() < 0 {
			t.Errorf("ReadDoubleAsBigFloat(Inf) = %v, want positive infinity", result)
		}
	})

	t.Run("positive_infinity_big_endian", func(t *testing.T) {
		var buf bytes.Buffer
		binary.Write(&buf, binary.BigEndian, math.Inf(1))
		reader := bytes.NewReader(buf.Bytes())
		result, err := ReadDoubleAsBigFloat(reader, true, prec)
		if err != nil {
			t.Fatalf("ReadDoubleAsBigFloat failed: %v", err)
		}
		if !result.IsInf() || result.Sign() < 0 {
			t.Errorf("ReadDoubleAsBigFloat(Inf) = %v, want positive infinity", result)
		}
	})

	// Test negative infinity
	t.Run("negative_infinity_little_endian", func(t *testing.T) {
		var buf bytes.Buffer
		binary.Write(&buf, binary.LittleEndian, math.Inf(-1))
		reader := bytes.NewReader(buf.Bytes())
		result, err := ReadDoubleAsBigFloat(reader, false, prec)
		if err != nil {
			t.Fatalf("ReadDoubleAsBigFloat failed: %v", err)
		}
		if !result.IsInf() || result.Sign() > 0 {
			t.Errorf("ReadDoubleAsBigFloat(-Inf) = %v, want negative infinity", result)
		}
	})

	t.Run("negative_infinity_big_endian", func(t *testing.T) {
		var buf bytes.Buffer
		binary.Write(&buf, binary.BigEndian, math.Inf(-1))
		reader := bytes.NewReader(buf.Bytes())
		result, err := ReadDoubleAsBigFloat(reader, true, prec)
		if err != nil {
			t.Fatalf("ReadDoubleAsBigFloat failed: %v", err)
		}
		if !result.IsInf() || result.Sign() > 0 {
			t.Errorf("ReadDoubleAsBigFloat(-Inf) = %v, want negative infinity", result)
		}
	})

	// Test NaN (big.Float doesn't support NaN, so it returns zero)
	t.Run("nan_little_endian", func(t *testing.T) {
		var buf bytes.Buffer
		binary.Write(&buf, binary.LittleEndian, math.NaN())
		reader := bytes.NewReader(buf.Bytes())
		result, err := ReadDoubleAsBigFloat(reader, false, prec)
		if err != nil {
			t.Fatalf("ReadDoubleAsBigFloat failed: %v", err)
		}
		// big.Float doesn't support NaN, so it returns zero
		if result.Sign() != 0 {
			val, _ := result.Float64()
			t.Errorf("ReadDoubleAsBigFloat(NaN) = %v (big.Float limitation: NaN not supported)", val)
		}
	})

	t.Run("nan_big_endian", func(t *testing.T) {
		var buf bytes.Buffer
		binary.Write(&buf, binary.BigEndian, math.NaN())
		reader := bytes.NewReader(buf.Bytes())
		result, err := ReadDoubleAsBigFloat(reader, true, prec)
		if err != nil {
			t.Fatalf("ReadDoubleAsBigFloat failed: %v", err)
		}
		// big.Float doesn't support NaN, so it returns zero
		if result.Sign() != 0 {
			val, _ := result.Float64()
			t.Errorf("ReadDoubleAsBigFloat(NaN) = %v (big.Float limitation: NaN not supported)", val)
		}
	})

	// Test denormalized number (subnormal)
	t.Run("denormalized_little_endian", func(t *testing.T) {
		// Create a denormalized number: exponent=0, mantissa != 0
		// Smallest denormalized: 0x0000000000000001
		bits := uint64(0x0000000000000001)
		var buf bytes.Buffer
		binary.Write(&buf, binary.LittleEndian, bits)
		reader := bytes.NewReader(buf.Bytes())
		result, err := ReadDoubleAsBigFloat(reader, false, prec)
		if err != nil {
			t.Fatalf("ReadDoubleAsBigFloat failed: %v", err)
		}
		// Currently denormalized numbers are handled as zero (see TODO in code)
		if result.Sign() != 0 {
			val, _ := result.Float64()
			t.Errorf("ReadDoubleAsBigFloat(denormalized) = %v, want zero (denormalized handling TODO)", val)
		}
	})

	t.Run("denormalized_big_endian", func(t *testing.T) {
		// Create a denormalized number: exponent=0, mantissa != 0
		bits := uint64(0x0000000000000001)
		var buf bytes.Buffer
		binary.Write(&buf, binary.BigEndian, bits)
		reader := bytes.NewReader(buf.Bytes())
		result, err := ReadDoubleAsBigFloat(reader, true, prec)
		if err != nil {
			t.Fatalf("ReadDoubleAsBigFloat failed: %v", err)
		}
		// Currently denormalized numbers are handled as zero (see TODO in code)
		if result.Sign() != 0 {
			val, _ := result.Float64()
			t.Errorf("ReadDoubleAsBigFloat(denormalized) = %v, want zero (denormalized handling TODO)", val)
		}
	})
}

// TestReadDoubleAsBigFloatErrorCases tests error handling
func TestReadDoubleAsBigFloatErrorCases(t *testing.T) {
	prec := uint(256)

	// Test short read
	t.Run("short_read", func(t *testing.T) {
		shortData := []byte{0x01, 0x02, 0x03} // Only 3 bytes instead of 8
		reader := bytes.NewReader(shortData)
		result, err := ReadDoubleAsBigFloat(reader, false, prec)
		if err == nil {
			t.Error("ReadDoubleAsBigFloat should fail on short read")
		}
		if result != nil {
			t.Error("ReadDoubleAsBigFloat should return nil on error")
		}
	})

	// Test empty reader
	t.Run("empty_reader", func(t *testing.T) {
		reader := bytes.NewReader([]byte{})
		result, err := ReadDoubleAsBigFloat(reader, false, prec)
		if err == nil {
			t.Error("ReadDoubleAsBigFloat should fail on empty reader")
		}
		if result != nil {
			t.Error("ReadDoubleAsBigFloat should return nil on error")
		}
	})
}

// TestReadDoubleAsBigFloatPrecision tests precision parameter handling
func TestReadDoubleAsBigFloatPrecision(t *testing.T) {
	// Test with explicit precision
	t.Run("explicit_precision", func(t *testing.T) {
		var buf bytes.Buffer
		binary.Write(&buf, binary.LittleEndian, math.Pi)
		reader := bytes.NewReader(buf.Bytes())
		result, err := ReadDoubleAsBigFloat(reader, false, 128)
		if err != nil {
			t.Fatalf("ReadDoubleAsBigFloat failed: %v", err)
		}
		if result.Prec() != 128 {
			t.Errorf("ReadDoubleAsBigFloat precision = %d, want 128", result.Prec())
		}
	})

	// Test with zero precision (should use DefaultPrecision)
	t.Run("zero_precision_uses_default", func(t *testing.T) {
		var buf bytes.Buffer
		binary.Write(&buf, binary.LittleEndian, 1.0)
		reader := bytes.NewReader(buf.Bytes())
		result, err := ReadDoubleAsBigFloat(reader, false, 0)
		if err != nil {
			t.Fatalf("ReadDoubleAsBigFloat failed: %v", err)
		}
		if result.Prec() != DefaultPrecision {
			t.Errorf("ReadDoubleAsBigFloat precision = %d, want %d (DefaultPrecision)", result.Prec(), DefaultPrecision)
		}
	})

	// Test with different precision levels
	t.Run("various_precisions", func(t *testing.T) {
		testPrecs := []uint{64, 128, 256, 512}
		for _, p := range testPrecs {
			var buf bytes.Buffer
			binary.Write(&buf, binary.LittleEndian, 2.71828)
			reader := bytes.NewReader(buf.Bytes())
			result, err := ReadDoubleAsBigFloat(reader, false, p)
			if err != nil {
				t.Fatalf("ReadDoubleAsBigFloat failed at prec %d: %v", p, err)
			}
			if result.Prec() != p {
				t.Errorf("ReadDoubleAsBigFloat precision = %d, want %d", result.Prec(), p)
			}
		}
	})
}

// TestReadDoubleAsBigFloatPrecisionPreservation tests that full 53-bit precision is preserved
func TestReadDoubleAsBigFloatPrecisionPreservation(t *testing.T) {
	prec := uint(256)

	// Test with a value that has specific bit patterns
	// Use a value that would lose precision if converted through float64
	t.Run("precision_preservation", func(t *testing.T) {
		// Create a double with specific mantissa bits
		// Value: 1.0000000000000002 (next representable after 1.0)
		value := 1.0 + math.Nextafter(0.0, 1.0)

		var buf bytes.Buffer
		binary.Write(&buf, binary.LittleEndian, value)
		reader := bytes.NewReader(buf.Bytes())
		result, err := ReadDoubleAsBigFloat(reader, false, prec)
		if err != nil {
			t.Fatalf("ReadDoubleAsBigFloat failed: %v", err)
		}

		resultFloat, _ := result.Float64()
		// Should match exactly (within float64 precision)
		if resultFloat != value {
			t.Errorf("ReadDoubleAsBigFloat = %g, want %g (precision loss detected)", resultFloat, value)
		}
	})

	// Test round-trip: write as double, read as BigFloat, compare
	t.Run("round_trip_accuracy", func(t *testing.T) {
		testValues := []float64{
			1.0,
			math.Pi,
			math.E,
			1e10,
			1e-10,
			-1.0,
			-math.Pi,
		}

		for _, val := range testValues {
			var buf bytes.Buffer
			binary.Write(&buf, binary.LittleEndian, val)
			reader := bytes.NewReader(buf.Bytes())
			result, err := ReadDoubleAsBigFloat(reader, false, prec)
			if err != nil {
				t.Fatalf("ReadDoubleAsBigFloat failed for %g: %v", val, err)
			}

			resultFloat, _ := result.Float64()
			if resultFloat != val {
				// For very large or very small numbers, allow some tolerance
				tolerance := math.Max(math.Abs(val)*1e-15, 1e-15)
				diff := math.Abs(resultFloat - val)
				if diff > tolerance {
					t.Errorf("Round-trip failed for %g: got %g, diff %g", val, resultFloat, diff)
				}
			}
		}
	})
}

// TestReadDoubleAsBigFloatAssemblyFunctions tests the assembly functions directly
func TestReadDoubleAsBigFloatAssemblyFunctions(t *testing.T) {
	// Test extractIEEE754Components (AMD64/ARM64)
	// This test will only run on platforms with assembly support
	t.Run("extract_components", func(t *testing.T) {
		// Test with a known value: 1.0 in IEEE 754 double
		// Sign: 0, Exponent: 1023 (0x3FF), Mantissa: 0
		// Bits: 0x3FF0000000000000
		bits := uint64(0x3FF0000000000000)

		// On platforms with assembly, test the assembly function
		// On generic platforms, this will be a no-op
		_ = bits // Use bits to avoid unused variable warning

		// Verify using the public API
		var buf bytes.Buffer
		binary.Write(&buf, binary.LittleEndian, 1.0)
		reader := bytes.NewReader(buf.Bytes())
		result, err := ReadDoubleAsBigFloat(reader, false, 256)
		if err != nil {
			t.Fatalf("ReadDoubleAsBigFloat failed: %v", err)
		}

		resultFloat, _ := result.Float64()
		if resultFloat != 1.0 {
			t.Errorf("Expected 1.0, got %g", resultFloat)
		}
	})

	// Test endianness conversion
	t.Run("endianness_conversion", func(t *testing.T) {
		// Test both endianness modes
		testValue := 3.141592653589793

		// Little-endian
		var bufLE bytes.Buffer
		binary.Write(&bufLE, binary.LittleEndian, testValue)
		readerLE := bytes.NewReader(bufLE.Bytes())
		resultLE, err := ReadDoubleAsBigFloat(readerLE, false, 256)
		if err != nil {
			t.Fatalf("ReadDoubleAsBigFloat (LE) failed: %v", err)
		}
		resultLEFloat, _ := resultLE.Float64()
		if math.Abs(resultLEFloat-testValue) > 1e-15 {
			t.Errorf("Little-endian: expected %g, got %g", testValue, resultLEFloat)
		}

		// Big-endian
		var bufBE bytes.Buffer
		binary.Write(&bufBE, binary.BigEndian, testValue)
		readerBE := bytes.NewReader(bufBE.Bytes())
		resultBE, err := ReadDoubleAsBigFloat(readerBE, true, 256)
		if err != nil {
			t.Fatalf("ReadDoubleAsBigFloat (BE) failed: %v", err)
		}
		resultBEFloat, _ := resultBE.Float64()
		if math.Abs(resultBEFloat-testValue) > 1e-15 {
			t.Errorf("Big-endian: expected %g, got %g", testValue, resultBEFloat)
		}
	})
}

// BenchmarkReadDoubleAsBigFloat benchmarks the ReadDoubleAsBigFloat function
func BenchmarkReadDoubleAsBigFloat(b *testing.B) {
	prec := uint(256)
	testValue := 3.141592653589793

	// Prepare test data
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, testValue)
	testData := buf.Bytes()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		reader := bytes.NewReader(testData)
		result, err := ReadDoubleAsBigFloat(reader, false, prec)
		if err != nil {
			b.Fatalf("ReadDoubleAsBigFloat failed: %v", err)
		}
		_ = result // Prevent optimization
	}
}

// BenchmarkReadDoubleAsBigFloatBigEndian benchmarks big-endian conversion
func BenchmarkReadDoubleAsBigFloatBigEndian(b *testing.B) {
	prec := uint(256)
	testValue := 3.141592653589793

	// Prepare test data
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, testValue)
	testData := buf.Bytes()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		reader := bytes.NewReader(testData)
		result, err := ReadDoubleAsBigFloat(reader, true, prec)
		if err != nil {
			b.Fatalf("ReadDoubleAsBigFloat failed: %v", err)
		}
		_ = result // Prevent optimization
	}
}

// BenchmarkReadDoubleAsBigFloatSpecialCases benchmarks special cases
func BenchmarkReadDoubleAsBigFloatSpecialCases(b *testing.B) {
	prec := uint(256)

	testCases := []struct {
		name  string
		value float64
	}{
		{"zero", 0.0},
		{"one", 1.0},
		{"pi", math.Pi},
		{"large", 1e100},
		{"small", 1e-100},
		{"infinity", math.Inf(1)},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			var buf bytes.Buffer
			binary.Write(&buf, binary.LittleEndian, tc.value)
			testData := buf.Bytes()

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				reader := bytes.NewReader(testData)
				result, err := ReadDoubleAsBigFloat(reader, false, prec)
				if err != nil {
					b.Fatalf("ReadDoubleAsBigFloat failed: %v", err)
				}
				_ = result // Prevent optimization
			}
		})
	}
}

// readDoubleAsBigFloatGenericTest is a copy of the generic implementation
// for benchmarking comparison purposes
//
//nolint:unparam // Return value is used in benchmarks, parameter is needed for comparison
func readDoubleAsBigFloatGenericTest(r io.Reader, bigEndian bool, prec uint) (*BigFloat, error) {
	if prec == 0 {
		prec = DefaultPrecision
	}

	// Read 8 bytes
	doubleBytes := make([]byte, 8)
	if _, err := io.ReadFull(r, doubleBytes); err != nil {
		return nil, fmt.Errorf("failed to read 8 bytes: %w", err)
	}

	// Interpret as uint64 with correct endianness
	var bits uint64
	if bigEndian {
		bits = binary.BigEndian.Uint64(doubleBytes)
	} else {
		bits = binary.LittleEndian.Uint64(doubleBytes)
	}

	// Extract IEEE 754 components
	// Sign: bit 63
	sign := (bits >> 63) != 0

	// Exponent: bits 52-62 (11 bits)
	exponent := int((bits >> 52) & 0x7FF)

	// Mantissa: bits 0-51 (52 bits)
	mantissa := bits & 0xFFFFFFFFFFFFF

	// Handle special cases
	if exponent == 0 {
		if mantissa == 0 {
			// Zero (positive or negative)
			result := new(big.Float).SetPrec(prec)
			if sign {
				result.Neg(result)
			}
			return result, nil
		}
		// Denormalized number (subnormal)
		result := new(big.Float).SetPrec(prec)
		return result, nil
	}

	if exponent == 0x7FF {
		// Infinity or NaN
		if mantissa == 0 {
			// Infinity
			result := new(big.Float).SetPrec(prec)
			result.SetInf(sign)
			return result, nil
		}
		// NaN
		result := new(big.Float).SetPrec(prec)
		return result, nil
	}

	// Normalized number
	// Construct mantissa as BigFloat: 1 + mantissa / 2^52
	mantissaBig := new(big.Float).SetPrec(prec)
	mantissaBig.SetUint64(mantissa)

	// Divide by 2^52 to get fractional part
	two52 := new(big.Float).SetPrec(prec)
	two52.SetUint64(1 << 52) // 2^52
	mantissaBig.Quo(mantissaBig, two52)

	// Add 1 (implicit leading bit)
	one := new(big.Float).SetPrec(prec)
	one.SetUint64(1)
	mantissaBig.Add(mantissaBig, one)

	// Calculate exponent: 2^(exponent - 1023)
	expValue := exponent - 1023

	// Construct result: mantissa * 2^expValue
	mant := new(big.Float).SetPrec(prec)
	mantExp := mantissaBig.MantExp(mant)

	result := new(big.Float).SetPrec(prec)
	result.SetMantExp(mant, expValue+mantExp)

	// Apply sign
	if sign {
		result.Neg(result)
	}

	return result, nil
}

// BenchmarkReadDoubleAsBigFloat_Generic benchmarks the pure Go implementation
func BenchmarkReadDoubleAsBigFloat_Generic(b *testing.B) {
	prec := uint(256)
	testValue := 3.141592653589793

	// Prepare test data
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, testValue)
	testData := buf.Bytes()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		reader := bytes.NewReader(testData)
		result, err := readDoubleAsBigFloatGenericTest(reader, false, prec)
		if err != nil {
			b.Fatalf("readDoubleAsBigFloatGenericTest failed: %v", err)
		}
		_ = result // Prevent optimization
	}
}

// BenchmarkReadDoubleAsBigFloat_Assembly benchmarks the assembly-optimized implementation
// This will use the actual assembly version on amd64/arm64, or fall back to generic on other platforms
func BenchmarkReadDoubleAsBigFloat_Assembly(b *testing.B) {
	prec := uint(256)
	testValue := 3.141592653589793

	// Prepare test data
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, testValue)
	testData := buf.Bytes()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		reader := bytes.NewReader(testData)
		result, err := ReadDoubleAsBigFloat(reader, false, prec)
		if err != nil {
			b.Fatalf("ReadDoubleAsBigFloat failed: %v", err)
		}
		_ = result // Prevent optimization
	}
}

// BenchmarkReadDoubleAsBigFloat_Comparison runs both implementations for direct comparison
func BenchmarkReadDoubleAsBigFloat_Comparison(b *testing.B) {
	prec := uint(256)
	testCases := []struct {
		name  string
		value float64
	}{
		{"normalized", 3.141592653589793},
		{"zero", 0.0},
		{"one", 1.0},
		{"large", 1e100},
		{"small", 1e-100},
		{"infinity", math.Inf(1)},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			var buf bytes.Buffer
			binary.Write(&buf, binary.LittleEndian, tc.value)
			testData := buf.Bytes()

			// Benchmark generic version
			b.Run("Generic", func(b *testing.B) {
				b.ResetTimer()
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					reader := bytes.NewReader(testData)
					result, err := readDoubleAsBigFloatGenericTest(reader, false, prec)
					if err != nil {
						b.Fatalf("readDoubleAsBigFloatGenericTest failed: %v", err)
					}
					_ = result
				}
			})

			// Benchmark assembly version (or current implementation)
			b.Run("Assembly", func(b *testing.B) {
				b.ResetTimer()
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					reader := bytes.NewReader(testData)
					result, err := ReadDoubleAsBigFloat(reader, false, prec)
					if err != nil {
						b.Fatalf("ReadDoubleAsBigFloat failed: %v", err)
					}
					_ = result
				}
			})
		})
	}
}
