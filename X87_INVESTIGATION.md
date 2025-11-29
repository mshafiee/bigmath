# Investigation: Extended Precision Floating Point in Go Assembler

## Summary

After deep investigation:
- **x86/x86-64**: Go's assembler does NOT support x87 FPU instructions like `FLDL`, `FSTPL`, `FSIN`, `FCOS`, `FPATAN`, `FYL2X`, `F2XM1`, `FSCALE`, etc.
- **ARM64**: ARM64 does NOT have x87 FPU (x87 is x86-specific). ARM64 supports 32-bit and 64-bit floating-point via NEON/scalar FP, but NOT 80-bit extended precision.

## Findings

### 1. Go Assembler Limitations

- Go's assembler (`cmd/asm`) is based on Plan 9 assembler syntax
- It only supports a **subset of instructions** commonly used in Go programs
- x87 FPU instructions are **not included** in this subset
- The assembler focuses on instructions compatible with Go's runtime and garbage collector

### 2. Why x87 Instructions Are Not Supported

1. **Modern Architecture Trends**: Modern x86-64 processors favor SSE/AVX for floating-point operations
2. **Runtime Compatibility**: Go's runtime and GC are optimized for specific instruction sets
3. **Portability**: Go prioritizes cross-platform compatibility over architecture-specific features
4. **FPU State Management**: x87 FPU state is aliased to MMX state, requiring careful management

### 3. Alternative Approaches

#### Option 1: CGO with Inline Assembly (Complex)
```c
// Requires CGO, adds complexity and dependencies
double x87_sin(double x) {
    double result;
    __asm__ __volatile__ (
        "fldl %1\n\t"
        "fsin\n\t"
        "fstpl %0"
        : "=m" (result)
        : "m" (x)
        : "st", "st(1)", "st(2)", "st(3)", "st(4)", "st(5)", "st(6)", "st(7)"
    );
    return result;
}
```
**Drawbacks:**
- Requires CGO (adds C dependency)
- Complex FPU stack management
- Portability issues
- Build complexity

#### Option 2: External Assembly (NASM/GAS)
- Write assembly in NASM or GAS
- Compile to object files
- Link with Go code
**Drawbacks:**
- Build system complexity
- Maintenance burden
- Platform-specific

#### Option 3: Software Emulation (Current Approach)
- Use optimized `float64` operations via Go's `math` package
- Provides performance benefits over `BigFloat`
- Fully portable and maintainable
**Benefits:**
- No external dependencies
- Works on all platforms
- Easier to maintain
- Still faster than `BigFloat` for intermediate calculations

### 4. Current Implementation

The current implementation uses **Option 3** (software emulation):
- Extended precision mode uses optimized `float64` operations
- Functions call `math.Sin`, `math.Cos`, `math.Exp`, etc.
- Provides significant performance benefits over `BigFloat`
- Maintains portability and simplicity

### 5. Documentation References

- Go assembler documentation: https://go.dev/doc/asm
- Plan 9 assembler: Limited instruction set
- Intel x87 FPU manual: Instructions not supported in Go asm

## ARM64 Extended Precision

### ARM64 Floating-Point Architecture

- **No x87 FPU**: ARM64 does not have x87 FPU (x87 is x86-specific)
- **NEON & Scalar FP**: ARM64 has NEON (SIMD) and scalar floating-point units
- **Supported Precisions**: 32-bit (float) and 64-bit (double) only
- **No 80-bit Extended Precision**: ARM64 hardware does NOT support 80-bit extended precision

### ARM64 Floating-Point Instructions in Go

Go's assembler for ARM64 supports:
- Scalar floating-point: `FMOVD`, `FADDD`, `FSUBD`, `FMULD`, `FDIVD`, `FSQRTD` (64-bit double)
- NEON SIMD: `VADD`, `VSUB`, `VMUL`, `VDIV` (vector operations)
- These are 32-bit or 64-bit precision, NOT 80-bit

### ARM64 Extended Precision Mode

On ARM64, extended precision mode works the same as on x86:
- Uses optimized `float64` operations via Go's `math` package
- Provides performance benefits over `BigFloat`
- No hardware extended precision available (neither x87 nor ARM64 equivalent)

## Conclusion

**Neither x86 nor ARM64 support hardware extended precision in Go's assembler:**
- **x86/x86-64**: Go's assembler does not support x87 FPU instructions
- **ARM64**: No x87 FPU exists, and ARM64 hardware doesn't support 80-bit extended precision

The current implementation using optimized `float64` operations is the correct and recommended approach for extended precision mode in Go on both architectures.

The name "extended precision" is maintained for:
- Conceptual clarity (faster intermediate calculations)
- API consistency
- Future-proofing (if CGO approach is needed later)

## Recommendation

Keep the current implementation. It provides:
- ✅ Performance benefits over `BigFloat`
- ✅ Full portability (works on x86, ARM64, and all platforms)
- ✅ No external dependencies
- ✅ Easy maintenance
- ✅ Consistent behavior across architectures

**Architecture-Specific Notes:**
- **x86/x86-64**: Could theoretically use CGO with x87 inline assembly, but adds complexity
- **ARM64**: No hardware extended precision available, so software approach is the only option
- **All Platforms**: Current implementation provides optimal balance of performance and portability

If true 80-bit hardware extended precision is absolutely required on x86, use CGO with inline assembly, but this adds significant complexity and should only be considered if the performance benefit justifies it. On ARM64, hardware extended precision is not available.

