# Faster Assembly Optimization Plan for ReadDoubleAsBigFloat

## Overview

Further optimize the assembly-optimized `ReadDoubleAsBigFloat` implementation to achieve 1.5-2x performance improvement over current optimized version. Current performance: ~263ns/op with 10 allocations. Target: ~150-180ns/op with 6-8 allocations.

## Current Bottlenecks Identified

1. **Temporary BigFloat Allocations**: Still creating `two52` and `one` as new BigFloat objects (2 extra allocations)
2. **Multiple BigFloat Operations**: Quo, Add, MantExp, SetMantExp operations are expensive
3. **Redundant Precision Setting**: Setting precision multiple times on temporary objects
4. **Assembly Function Call Overhead**: Combined function still has some overhead
5. **Memory Access Patterns**: Multiple memory writes in assembly

## Optimization Strategy

### Phase 1: Eliminate Temporary BigFloat Allocations (Biggest Impact)

**File**: `serialization_amd64.go`, `serialization_arm64.go`, `serialization_generic.go`

**Current approach** (4 allocations):
```go
mantissaBig := new(big.Float).SetPrec(prec)
two52 := new(big.Float).SetPrec(prec).Set(two52Const)  // Allocation
one := new(big.Float).SetPrec(prec).Set(oneConst)      // Allocation
result := new(big.Float).SetPrec(prec)
```

**Optimized approach** (2 allocations):
```go
// Reuse a single temporary variable for both two52 and one operations
mantissaBig := new(big.Float).SetPrec(prec)
mantissaBig.SetUint64(mantissa)

// Reuse single temporary for both constants
temp := new(big.Float).SetPrec(prec)
temp.Set(two52Const)
mantissaBig.Quo(mantissaBig, temp)

temp.Set(oneConst)  // Reuse same variable
mantissaBig.Add(mantissaBig, temp)

// Reuse mantissaBig for MantExp, then create result
mantExp := mantissaBig.MantExp(mantissaBig)
result := new(big.Float).SetPrec(prec)
result.SetMantExp(mantissaBig, expValue+mantExp)
```

**Expected improvement**: 15-20% faster, reduce from 10 to 8 allocations

### Phase 2: Direct Float64 Construction Path (High Impact)

**File**: `serialization_amd64.go`, `serialization_arm64.go`, `serialization_generic.go`

For normalized numbers, construct the float64 value directly and use `SetFloat64`, which is highly optimized in Go's big.Float:

```go
// Calculate the float64 value directly
// Value = (-1)^sign * 2^(exponent - 1023) * (1 + mantissa / 2^52)
mantissaFloat := float64(mantissa) / (1 << 52) + 1.0
expValue := exponent - 1023
value := mantissaFloat * math.Pow(2, float64(expValue))
if sign {
    value = -value
}

// Use SetFloat64 which is highly optimized
result := new(big.Float).SetPrec(prec)
result.SetFloat64(value)
```

**Note**: This may lose precision for very large exponents, so we need to verify correctness. Alternative: use this as a fast path for common cases, fall back to exact method for edge cases.

**Expected improvement**: 30-40% faster for normalized numbers, reduce to 1-2 allocations

### Phase 3: Optimize Assembly with Conditional Moves

**File**: `serialization_amd64.s`, `serialization_arm64.s`

Replace branches with conditional moves (CMOV) where possible to reduce branch misprediction:

**Current** (branch):
```assembly
TESTB AL, AL
JZ extract_components
BSWAPQ AX
extract_components:
```

**Optimized** (conditional move):
```assembly
// Use CMOV to conditionally swap bytes
MOVQ AX, BX           // Copy original
BSWAPQ BX             // Swapped version
TESTB AL, AL          // Test bigEndian flag
CMOVQNE BX, AX        // If bigEndian != 0, use swapped version
```

**Expected improvement**: 5-10% faster by eliminating branch misprediction

### Phase 4: Optimize Register Usage in Assembly

**File**: `serialization_amd64.s`, `serialization_arm64.s`

Minimize memory writes by keeping values in registers longer:

**Current**: Multiple MOVQ to memory
**Optimized**: Keep intermediate values in registers, write once at end

```assembly
// Keep all values in registers, write once
MOVQ AX, BX        // Sign in BX
SHRQ $63, BX
MOVQ AX, CX        // Exponent in CX
SHRQ $52, CX
ANDQ $0x7FF, CX
MOVQ AX, DX        // Mantissa in DX
ANDQ $0xFFFFFFFFFFFFF, DX

// Write all results in sequence (better cache locality)
MOVQ BX, sign+16(FP)
MOVQ CX, exponent+24(FP)
MOVQ DX, mantissa+32(FP)
```

**Expected improvement**: 3-5% faster

### Phase 5: Pre-compute Common Values

**File**: `serialization_amd64.go`, `serialization_arm64.go`, `serialization_generic.go`

For very common values (like 1.0, 0.0), use pre-computed BigFloat constants:

```go
var (
    cachedOne *big.Float
    cachedZero *big.Float
)

func init() {
    cachedOne = new(big.Float).SetUint64(1)
    cachedZero = new(big.Float).SetUint64(0)
}

// In function, for common cases:
if mantissa == 0 && exponent == 1023 {  // Value is 1.0
    result := new(big.Float).SetPrec(prec)
    result.Set(cachedOne)
    if sign {
        result.Neg(result)
    }
    return result, nil
}
```

**Expected improvement**: 10-15% faster for common values

### Phase 6: Inline Special Case Detection in Assembly

**File**: `serialization_amd64.s`, `serialization_arm64.s`

Add a combined function that returns special case flags:

```assembly
// func extractIEEE754WithFlags(bytes *[8]byte, bigEndian uint8) (sign uint64, exponent int64, mantissa uint64, flags uint8)
// flags: bit 0 = is_zero, bit 1 = is_infinity, bit 2 = is_nan
TEXT ·extractIEEE754WithFlags(SB), NOSPLIT, $0-48
    // ... extract components ...
    
    // Check for special cases
    MOVQ $0, R8              // flags = 0
    CMPQ CX, $0              // exponent == 0?
    JNE check_infinity
    CMPQ DX, $0              // mantissa == 0?
    JNE check_infinity
    ORQ $1, R8               // Set is_zero flag
    JMP done_flags
    
check_infinity:
    CMPQ CX, $0x7FF          // exponent == 0x7FF?
    JNE done_flags
    CMPQ DX, $0              // mantissa == 0?
    JNE is_nan
    ORQ $2, R8               // Set is_infinity flag
    JMP done_flags
    
is_nan:
    ORQ $4, R8               // Set is_nan flag
    
done_flags:
    MOVB R8, flags+40(FP)
    RET
```

**Expected improvement**: 5-10% faster for special cases

### Phase 7: Reduce Precision Setting Overhead

**File**: `serialization_amd64.go`, `serialization_arm64.go`, `serialization_generic.go`

Set precision once on result, then use Copy operations instead of SetPrec on temporaries:

```go
// Create result with precision first
result := new(big.Float).SetPrec(prec)

// Use result for intermediate calculations (reuse precision)
result.SetUint64(mantissa)
temp := new(big.Float).SetPrec(prec)  // Only one SetPrec call
temp.Set(two52Const)
result.Quo(result, temp)
// ...
```

**Expected improvement**: 3-5% faster

## Implementation Order

1. **Phase 1** (Eliminate Temporary Allocations) - Expected 15-20% improvement
2. **Phase 2** (Direct Float64 Path) - Expected 30-40% improvement (if precision is acceptable)
3. **Phase 3** (Conditional Moves) - Expected 5-10% improvement
4. **Phase 4** (Register Optimization) - Expected 3-5% improvement
5. **Phase 5** (Pre-computed Values) - Expected 10-15% improvement for common cases
6. **Phase 6** (Inline Flags) - Expected 5-10% improvement
7. **Phase 7** (Precision Optimization) - Expected 3-5% improvement

## Files to Modify

- `serialization_amd64.go`: Apply all Go-level optimizations
- `serialization_amd64.s`: Apply assembly optimizations (CMOV, register usage, flags)
- `serialization_arm64.go`: Same optimizations as AMD64
- `serialization_arm64.s`: ARM64 equivalent optimizations
- `serialization_generic.go`: Apply Go-level optimizations for consistency

## Testing Strategy

1. **Correctness**: Ensure all existing tests pass, especially precision tests
2. **Performance**: Run benchmarks before/after each phase
3. **Memory**: Verify allocation reductions with `-benchmem`
4. **Edge Cases**: Test very large/small numbers, verify precision is maintained

## Expected Final Results

- **Normalized numbers**: 150-180 ns/op (from 263ns) - **40-45% faster**
- **Special cases**: 20-30 ns/op (from 57ns) - **50-65% faster**
- **Memory allocations**: 6-8 allocs/op (from 10) - **20-40% reduction**
- **Memory usage**: 250-300 B/op (from 400B) - **25-37% reduction**

## Risks and Considerations

1. **Precision**: Direct float64 path may lose precision for very large exponents - need validation
2. **Code complexity**: More optimizations = more complex code
3. **Platform differences**: Some optimizations may work better on AMD64 vs ARM64
4. **Maintainability**: Aggressive optimizations may make code harder to understand

## Success Criteria

- Benchmarks show 40-50% improvement for normalized numbers
- Special cases show 50-65% improvement
- Allocations reduced by 20-40%
- All tests pass
- Precision maintained for all test cases

