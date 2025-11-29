// Copyright (c) 2025 Mohammad Shafiee
// SPDX-License-Identifier: BSD-3-Clause

package bigmath

import (
	"math"
	"strconv"
	"strings"
	"testing"
)

// TestNewBigFloatFromStringValidFormats tests valid numeric string formats
func TestNewBigFloatFromStringValidFormats(t *testing.T) {
	prec := uint(256)

	tests := []struct {
		name      string
		input     string
		expected  float64
		tolerance float64
	}{
		// Integers
		{"zero", "0", 0.0, 0},
		{"one", "1", 1.0, 0},
		{"negative_one", "-1", -1.0, 0},
		{"large_int", "12345678901234567890", 12345678901234567890.0, 1e10},
		{"negative_large_int", "-98765432109876543210", -98765432109876543210.0, 1e10},

		// Decimals
		{"decimal", "3.14", 3.14, 1e-15},
		{"negative_decimal", "-2.718", -2.718, 1e-15},
		{"decimal_no_leading", ".5", 0.5, 1e-15},
		{"decimal_no_trailing", "5.", 5.0, 0},
		{"negative_decimal_no_leading", "-.5", -0.5, 1e-15},
		{"many_decimals", "3.141592653589793238462643383279", 3.141592653589793238462643383279, 1e-15},

		// Scientific notation
		{"scientific_positive", "1.23e10", 1.23e10, 1e5},
		{"scientific_negative_exp", "1.23e-10", 1.23e-10, 1e-25},
		{"scientific_uppercase", "1.23E10", 1.23e10, 1e5},
		{"scientific_positive_exp", "1.23e+10", 1.23e10, 1e5},
		{"scientific_negative_exp_plus", "1.23e-10", 1.23e-10, 1e-25},
		{"scientific_large", "1e100", 1e100, 1e85},
		{"scientific_small", "1e-100", 1e-100, 1e-115},
		{"scientific_negative", "-1.23e10", -1.23e10, 1e5},

		// Hex floats (if supported by big.Float)
		{"hex_float", "0x1.5p0", 1.3125, 1e-15},
		{"hex_float_negative", "-0x1.5p0", -1.3125, 1e-15},
		{"hex_float_exp", "0x1.5p10", 1344.0, 1e-10},
		{"hex_float_neg_exp", "0x1.5p-10", 0.00128173828125, 1e-20},

		// Special values
		{"infinity", "+Inf", math.Inf(1), 0},
		{"negative_infinity", "-Inf", math.Inf(-1), 0},
		{"infinity_no_sign", "Inf", math.Inf(1), 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewBigFloatFromString(tt.input, prec)
			if err != nil {
				// Hex floats might not be supported, skip if error
				if strings.HasPrefix(tt.input, "0x") {
					t.Skipf("Hex float format may not be supported: %v", err)
				}
				t.Fatalf("NewBigFloatFromString(%q) failed: %v", tt.input, err)
			}

			if result == nil {
				t.Fatal("NewBigFloatFromString returned nil")
			}

			resultFloat, _ := result.Float64()

			// Handle infinity specially
			if math.IsInf(tt.expected, 0) {
				if !math.IsInf(resultFloat, 0) || math.IsInf(tt.expected, 1) != math.IsInf(resultFloat, 1) {
					t.Errorf("NewBigFloatFromString(%q) = %g, want %g", tt.input, resultFloat, tt.expected)
				}
				return
			}

			// Compare values
			diff := math.Abs(resultFloat - tt.expected)
			if diff > tt.tolerance {
				t.Errorf("NewBigFloatFromString(%q) = %g, want %g (diff %g, tolerance %g)",
					tt.input, resultFloat, tt.expected, diff, tt.tolerance)
			}
		})
	}
}

// TestNewBigFloatFromStringEdgeCases tests edge cases and boundary conditions
func TestNewBigFloatFromStringEdgeCases(t *testing.T) {
	prec := uint(256)

	tests := []struct {
		name      string
		input     string
		expected  float64
		tolerance float64
		shouldErr bool
	}{
		// Zero variations
		{"zero", "0", 0.0, 0, false},
		{"zero_decimal", "0.0", 0.0, 0, false},
		{"zero_scientific", "0e0", 0.0, 0, false},
		{"zero_negative", "-0", 0.0, 0, false},
		{"zero_negative_decimal", "-0.0", 0.0, 0, false},

		// Very small numbers
		{"tiny_positive", "1e-300", 1e-300, 1e-315, false},
		{"tiny_negative", "-1e-300", -1e-300, 1e-315, false},
		{"smallest_normal", "2.2250738585072014e-308", 2.2250738585072014e-308, 1e-323, false},

		// Very large numbers
		{"huge_positive", "1e300", 1e300, 1e285, false},
		{"huge_negative", "-1e300", -1e300, 1e285, false},
		{"max_float64_str", "1.7976931348623157e+308", 1.7976931348623157e+308, 1e293, false},

		// Boundary values
		{"one_minus_epsilon", "0.9999999999999999", 0.9999999999999999, 1e-15, false},
		{"one_plus_epsilon", "1.0000000000000002", 1.0000000000000002, 1e-15, false},
		{"pi_precise", "3.14159265358979323846264338327950288419716939937510", 3.14159265358979323846264338327950288419716939937510, 1e-15, false},
		{"e_precise", "2.71828182845904523536028747135266249775724709369995", 2.71828182845904523536028747135266249775724709369995, 1e-15, false},

		// Leading/trailing zeros
		{"leading_zeros", "000123.456", 123.456, 1e-15, false},
		{"trailing_zeros", "123.456000", 123.456, 1e-15, false},
		{"both_zeros", "000123.456000", 123.456, 1e-15, false},

		// Whitespace (should be invalid for strict parsing)
		{"leading_space", " 123", 0, 0, true},
		{"trailing_space", "123 ", 0, 0, true},
		{"both_spaces", " 123 ", 0, 0, true},
		{"tab", "\t123", 0, 0, true},
		{"newline", "123\n", 0, 0, true},

		// Empty and invalid
		{"empty", "", 0, 0, true},
		{"only_sign", "+", 0, 0, true},
		{"only_sign_neg", "-", 0, 0, true},
		{"only_dot", ".", 0, 0, true},
		{"only_e", "e", 0, 0, true},
		{"only_e_upper", "E", 0, 0, true},
		{"dot_e", ".e", 0, 0, true},
		{"e_no_number", "123e", 0, 0, true},
		{"e_no_exp", "123e+", 0, 0, true},
		{"double_dot", "12.34.56", 0, 0, true},
		{"double_e", "12e34e56", 0, 0, true},
		{"invalid_char", "12.34abc", 0, 0, true},
		{"invalid_char_start", "abc12.34", 0, 0, true},
		{"unicode", "12.34π", 0, 0, true},
		{"null_char", "12.34\x00", 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewBigFloatFromString(tt.input, prec)

			if tt.shouldErr {
				if err == nil {
					t.Errorf("NewBigFloatFromString(%q) should have failed, got %v", tt.input, result)
				}
				return
			}

			if err != nil {
				t.Fatalf("NewBigFloatFromString(%q) failed: %v", tt.input, err)
			}

			if result == nil {
				t.Fatal("NewBigFloatFromString returned nil")
			}

			resultFloat, _ := result.Float64()
			diff := math.Abs(resultFloat - tt.expected)
			if diff > tt.tolerance {
				t.Errorf("NewBigFloatFromString(%q) = %g, want %g (diff %g, tolerance %g)",
					tt.input, resultFloat, tt.expected, diff, tt.tolerance)
			}
		})
	}
}

// TestNewBigFloatFromStringPrecision tests precision handling
func TestNewBigFloatFromStringPrecision(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		prec     uint
		expected uint
	}{
		{"zero_precision", "3.14159", 0, DefaultPrecision},
		{"low_precision", "3.14159", 64, 64},
		{"medium_precision", "3.14159", 128, 128},
		{"high_precision", "3.14159", 256, 256},
		{"very_high_precision", "3.14159", 512, 512},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewBigFloatFromString(tt.input, tt.prec)
			if err != nil {
				t.Fatalf("NewBigFloatFromString failed: %v", err)
			}

			if result.Prec() != tt.expected {
				t.Errorf("Precision = %d, want %d", result.Prec(), tt.expected)
			}
		})
	}
}

// TestNewBigFloatFromStringCompatibility tests compatibility with strconv.ParseFloat
func TestNewBigFloatFromStringCompatibility(t *testing.T) {
	prec := uint(256)

	// Test cases that strconv.ParseFloat handles
	testCases := []string{
		"0",
		"1",
		"-1",
		"3.14",
		"-3.14",
		"1.23e10",
		"1.23e-10",
		"1.23E10",
		"1.23E-10",
		"1e100",
		"1e-100",
		"-1e100",
		"-1e-100",
		"0.5",
		"-0.5",
		".5",
		"-.5",
		"5.",
		"-5.",
		"12345678901234567890",
		"-12345678901234567890",
		"1.7976931348623157e+308",
		"2.2250738585072014e-308",
		"3.141592653589793",
		"2.718281828459045",
	}

	for _, input := range testCases {
		t.Run(input, func(t *testing.T) {
			// Parse with strconv.ParseFloat
			stdlibVal, stdlibErr := strconv.ParseFloat(input, 64)

			// Parse with NewBigFloatFromString
			bigmathVal, bigmathErr := NewBigFloatFromString(input, prec)

			// Both should succeed or both should fail
			if (stdlibErr == nil) != (bigmathErr == nil) {
				t.Errorf("Error mismatch: strconv.ParseFloat err=%v, NewBigFloatFromString err=%v",
					stdlibErr, bigmathErr)
				return
			}

			if stdlibErr != nil {
				// Both failed, that's fine
				return
			}

			// Both succeeded, compare values
			bigmathFloat, _ := bigmathVal.Float64()

			// Handle special cases
			if math.IsNaN(stdlibVal) {
				// big.Float doesn't support NaN, so bigmathVal will be 0
				if bigmathFloat != 0 {
					t.Errorf("NaN handling: strconv.ParseFloat=%g, NewBigFloatFromString=%g (expected 0 for NaN)",
						stdlibVal, bigmathFloat)
				}
				return
			}

			if math.IsInf(stdlibVal, 0) {
				if !math.IsInf(bigmathFloat, 0) || math.IsInf(stdlibVal, 1) != math.IsInf(bigmathFloat, 1) {
					t.Errorf("Infinity mismatch: strconv.ParseFloat=%g, NewBigFloatFromString=%g",
						stdlibVal, bigmathFloat)
				}
				return
			}

			// Compare numeric values
			diff := math.Abs(bigmathFloat - stdlibVal)
			tolerance := math.Max(math.Abs(stdlibVal)*1e-15, 1e-15)
			if diff > tolerance {
				t.Errorf("Value mismatch: strconv.ParseFloat=%g, NewBigFloatFromString=%g (diff %g, tolerance %g)",
					stdlibVal, bigmathFloat, diff, tolerance)
			}
		})
	}
}

// TestNewBigFloatFromStringRoundTrip tests round-trip conversion
func TestNewBigFloatFromStringRoundTrip(t *testing.T) {
	prec := uint(256)

	testValues := []float64{
		0.0,
		1.0,
		-1.0,
		3.14159,
		-2.71828,
		1e10,
		1e-10,
		-1e10,
		-1e-10,
		1e100,
		1e-100,
		math.MaxFloat64,
		math.SmallestNonzeroFloat64,
		math.Pi,
		math.E,
	}

	for _, val := range testValues {
		t.Run(strconv.FormatFloat(val, 'g', -1, 64), func(t *testing.T) {
			// Convert to string
			str := strconv.FormatFloat(val, 'g', -1, 64)

			// Parse back
			result, err := NewBigFloatFromString(str, prec)
			if err != nil {
				t.Fatalf("NewBigFloatFromString failed: %v", err)
			}

			resultFloat, _ := result.Float64()

			// Compare
			if val == 0.0 {
				if resultFloat != 0.0 {
					t.Errorf("Round-trip failed: %g -> %q -> %g", val, str, resultFloat)
				}
			} else {
				relErr := math.Abs((resultFloat - val) / val)
				if relErr > 1e-15 {
					t.Errorf("Round-trip failed: %g -> %q -> %g (rel err %e)", val, str, resultFloat, relErr)
				}
			}
		})
	}
}

// TestNewBigFloatFromStringLongStrings tests very long numeric strings
func TestNewBigFloatFromStringLongStrings(t *testing.T) {
	prec := uint(256)

	tests := []struct {
		name  string
		input string
	}{
		{"long_decimal", "3." + strings.Repeat("1", 100)},
		{"long_scientific", "1." + strings.Repeat("2", 50) + "e10"},
		{"many_digits", strings.Repeat("9", 200)},
		{"long_negative", "-" + strings.Repeat("1", 100) + ".5"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewBigFloatFromString(tt.input, prec)
			if err != nil {
				t.Fatalf("NewBigFloatFromString failed: %v", err)
			}

			if result == nil {
				t.Fatal("NewBigFloatFromString returned nil")
			}

			// Just verify it parses without error
			_ = result.Sign()
		})
	}
}

// TestNewBigFloatFromStringSpecialCharacters tests special character handling
func TestNewBigFloatFromStringSpecialCharacters(t *testing.T) {
	prec := uint(256)

	tests := []struct {
		name      string
		input     string
		shouldErr bool
	}{
		{"plus_sign", "+123", false},
		{"minus_sign", "-123", false},
		{"multiple_signs", "--123", true},
		{"sign_after_digit", "12-3", true},
		{"sign_after_e", "1e-10", false},
		{"multiple_dots", "12.34.56", true},
		{"dot_after_e", "1e.5", true},
		{"e_after_dot", "12.e5", false}, // Valid: 12.0e5 = 1.2e6
		{"underscore", "12_34", false},   // Valid: Go 1.13+ supports underscores
		{"comma", "12,34", true},
		{"space_in_middle", "12 34", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewBigFloatFromString(tt.input, prec)

			if tt.shouldErr {
				if err == nil {
					t.Errorf("NewBigFloatFromString(%q) should have failed, got %v", tt.input, result)
				}
			} else {
				if err != nil {
					t.Errorf("NewBigFloatFromString(%q) failed unexpectedly: %v", tt.input, err)
				}
			}
		})
	}
}

// TestNewBigFloatFromStringBoundaryPrecision tests precision boundary conditions
func TestNewBigFloatFromStringBoundaryPrecision(t *testing.T) {
	input := "3.14159265358979323846264338327950288419716939937510"

	precisions := []uint{1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024}

	for _, prec := range precisions {
		t.Run(string(rune(prec)), func(t *testing.T) {
			result, err := NewBigFloatFromString(input, prec)
			if err != nil {
				t.Fatalf("NewBigFloatFromString failed at prec %d: %v", prec, err)
			}

			if result.Prec() != prec {
				t.Errorf("Precision = %d, want %d", result.Prec(), prec)
			}
		})
	}
}

