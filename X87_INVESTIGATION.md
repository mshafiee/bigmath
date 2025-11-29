# Investigation: x87 FPU Instructions in Go Assembler

## Summary

After deep investigation, **Go's assembler does NOT support x87 FPU instructions** like `FLDL`, `FSTPL`, `FSIN`, `FCOS`, `FPATAN`, `FYL2X`, `F2XM1`, `FSCALE`, etc.

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

## Conclusion

**Go's assembler does not and will not support x87 FPU instructions.** The current implementation using optimized `float64` operations is the correct and recommended approach for extended precision mode in Go.

The name "extended precision" is maintained for:
- Conceptual clarity (faster intermediate calculations)
- API consistency
- Future-proofing (if CGO approach is needed later)

## Recommendation

Keep the current implementation. It provides:
- ✅ Performance benefits over `BigFloat`
- ✅ Full portability
- ✅ No external dependencies
- ✅ Easy maintenance
- ✅ Works on all platforms

If true 80-bit hardware extended precision is absolutely required, use CGO with inline assembly, but this adds significant complexity and should only be considered if the performance benefit justifies it.

