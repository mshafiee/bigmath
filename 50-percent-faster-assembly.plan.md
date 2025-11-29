# Plan: Make Assembly 50% Faster

## Current Performance Baseline

- **Normalized numbers**: ~306.6 ns/op, 328 B/op, 8 allocs/op
- **Target**: ~150-180 ns/op (50% improvement)
- **One case**: ~83.85 ns/op (already fast, but can improve further)

## Bottleneck Analysis

Current implementation still has these bottlenecks:

1. **BigFloat Arithmetic Operations** (60-70% of time):
   - `Quo` (division) operation
   - `Add` operation  
   - `MantExp` operation
   - `SetMantExp` operation
   - These are expensive arbitrary-precision operations

2. **Allocations** (15-20% of time):
   - Still creating 2 BigFloat objects (result + temp)
   - Memory allocation overhead

3. **Go Wrapper Overhead** (10-15% of time):
   - Function call overhead
   - Interface conversions
   - Error handling

## Optimization Strategy

### Phase 1: Direct Float64 Construction Path (Biggest Impact - 30-40% improvement)

**File**: `serialization_amd64.go`, `serialization_arm64.go`, `serialization_generic.go`

For normalized numbers, construct the float64 value directly and use `SetFloat64`, which is highly optimized in Go's big.Float. This avoids all the BigFloat arithmetic operations.

**Implementation**:
```go
// For normalized numbers, calculate float64 directly
// Value = (-1)^sign * 2^(exponent - 1023) * (1 + mantissa / 2^52)
mantissaFloat := float64(mantissa) / (1 << 52) + 1.0
expValue := exponent - 1023

// Calculate value using float64 arithmetic
value := mantissaFloat * math.Pow(2, float64(expValue))
if sign {
    value = -value
}

// Use SetFloat64 which is highly optimized
result := new(big.Float).SetPrec(prec)
result.SetFloat64(value)
```

**Precision consideration**: For very large exponents (>1022 or <-1022), we need to fall back to exact method. Use this as fast path for common cases.

**Expected improvement**: 30-40% faster (from 306ns to ~180-210ns)

### Phase 2: Eliminate All Temporary Allocations (15-20% improvement)

**File**: `serialization_amd64.go`, `serialization_arm64.go`, `serialization_generic.go`

Use object pooling or reuse the result object for all intermediate calculations:

**Current**: Creates `result` and `temp` (2 allocations)
**Optimized**: Reuse `result` for all operations, eliminate `temp`:

```go
result := new(big.Float).SetPrec(prec)
result.SetUint64(mantissa)

// Reuse result for division - need to save value first
// Actually, we can use Quo with result as both arguments in some cases
// Or use a different approach: calculate in float64, then convert

// Better approach: Use float64 path (Phase 1) which eliminates temp entirely
```

**Expected improvement**: 15-20% faster, reduce to 1 allocation

### Phase 3: Assembly-Level Float64 Construction (10-15% improvement)

**File**: `serialization_amd64.s`, `serialization_arm64.s`

Create an assembly function that constructs the float64 value directly from IEEE 754 components:

```assembly
// func constructFloat64FromIEEE754(sign uint64, exponent int64, mantissa uint64) float64
TEXT ·constructFloat64FromIEEE754(SB), NOSPLIT, $0-24
    // Reconstruct IEEE 754 double from components
    // sign: bit 63
    // exponent: bits 52-62
    // mantissa: bits 0-51
    
    MOVQ sign+0(FP), AX
    SHLQ $63, AX                  // Shift sign to bit 63
    
    MOVQ exponent+8(FP), BX
    SHLQ $52, BX                  // Shift exponent to bits 52-62
    ANDQ $0x7FF0000000000000, BX  // Mask exponent bits
    
    MOVQ mantissa+16(FP), CX
    ANDQ $0xFFFFFFFFFFFFF, CX     // Mask mantissa bits
    
    // Combine: sign | exponent | mantissa
    ORQ BX, AX                    // AX = sign | exponent
    ORQ CX, AX                    // AX = sign | exponent | mantissa
    
    // Return as float64 (same bit pattern)
    MOVQ AX, ret+24(FP)
    RET
```

Then in Go code:
```go
// Use assembly to construct float64 directly
floatValue := constructFloat64FromIEEE754(signUint, exponentInt, mantissaUint)
result := new(big.Float).SetPrec(prec)
result.SetFloat64(floatValue)
```

**Expected improvement**: 10-15% faster by avoiding Go-level bit manipulation

### Phase 4: Batch Operations with SIMD (Optional, Advanced - 20-30% for batch)

**File**: `serialization_amd64.s`, `serialization_arm64.s`

For batch processing (reading multiple doubles), use SIMD instructions to process 2-4 values at once:

```assembly
// Process 2 doubles at once using SSE/AVX
// func extractIEEE754FromBytesBatch2(bytes *[16]byte, bigEndian uint8) (signs, exponents, mantissas [2]uint64)
```

**Note**: This only helps if caller processes multiple values. May not apply to single-value reads.

### Phase 5: Inline More Operations (5-10% improvement)

**File**: `serialization_amd64.go`, `serialization_arm64.go`

Use compiler hints and reduce function call overhead:

```go
// Mark function as inline candidate
//go:inline
func fastPathNormalized(mantissa uint64, exponent int, sign bool, prec uint) *big.Float {
    // Inline the fast path logic
}
```

### Phase 6: Optimize Special Cases Further (5-10% improvement)

**File**: `serialization_amd64.go`, `serialization_arm64.go`

Add more fast paths for common values:
- Powers of 2 (exponent patterns)
- Small integers (0-100)
- Common constants (Pi, e, etc.)

### Phase 7: Reduce Precision Setting (3-5% improvement)

**File**: `serialization_amd64.go`, `serialization_arm64.go`

Cache precision-specific BigFloat objects or use a pool:

```go
var precisionPool = sync.Map{} // map[uint]*big.Float

func getPrecisionTemplate(prec uint) *big.Float {
    if cached, ok := precisionPool.Load(prec); ok {
        return cached.(*big.Float)
    }
    // Create and cache
}
```

## Implementation Priority

1. **Phase 1** (Direct Float64 Path) - **30-40% improvement** - Highest priority
2. **Phase 2** (Eliminate Temp Allocations) - **15-20% improvement** - Works with Phase 1
3. **Phase 3** (Assembly Float64 Construction) - **10-15% improvement** - Medium priority
4. **Phase 5** (Inline Operations) - **5-10% improvement** - Low priority
5. **Phase 6** (More Fast Paths) - **5-10% improvement** - Low priority
6. **Phase 7** (Precision Pooling) - **3-5% improvement** - Optional

## Expected Combined Results

With Phases 1-3 implemented:
- **Normalized numbers**: ~150-180 ns/op (50% improvement) ✅
- **Allocations**: 1-2 allocs/op (from 8)
- **Memory**: ~200-250 B/op (from 328B)

## Risks and Considerations

1. **Precision Loss**: Float64 path may lose precision for very large exponents
   - **Solution**: Use float64 for common range, exact method for edge cases
   - **Validation**: Test with very large/small numbers

2. **Platform Differences**: Some optimizations may work better on AMD64 vs ARM64
   - **Solution**: Platform-specific implementations

3. **Code Complexity**: More optimizations = more complex code
   - **Solution**: Clear comments, maintain tests

## Testing Strategy

1. **Correctness**: Ensure all existing tests pass
2. **Precision**: Test with very large exponents (>1022, <-1022)
3. **Performance**: Benchmark before/after each phase
4. **Edge Cases**: Test denormalized, infinity, NaN, zero

## Success Criteria

- Benchmarks show 50% improvement (from ~306ns to ~150-180ns)
- Allocations reduced to 1-2 (from 8)
- Memory usage reduced by 30-40%
- All tests pass
- Precision maintained for all test cases

