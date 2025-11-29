# Build Tags and Assembly Function Declarations

## Extended Precision (x87 FPU) Support

The library includes support for hardware extended precision (80-bit x87 FPU) on x86/x86-64 platforms:

- **Build tag**: `//go:build amd64 || 386` - Enables x87 extended precision support
- **Files**:
  - `extended_precision_decl.go` - Function declarations for x87 operations
  - `extended_precision_amd64.s` - x87 FPU assembly implementations
  - `extended_precision_fallback.go` - Fallback to standard math on other platforms

Extended precision mode is activated by setting `prec = ExtendedPrecision` (80) when calling functions.
On platforms without x87 support, operations automatically fall back to BigFloat implementations.

## Pattern for Assembly Function Declarations

When adding assembly function declarations, follow this pattern to ensure cross-platform compatibility:

### Rule: Assembly declarations MUST have either:
1. **Build tags limiting to platforms with assembly implementations** (e.g., `//go:build amd64`)
2. **Fallback implementations for generic platforms** (e.g., `//go:build !amd64 && !arm64`)

### Examples:

#### ✅ Correct Pattern 1: Build tag for specific platform
```go
//go:build amd64

package bigmath

//go:noescape
func myAsmFunction(x *BigFloat) *BigFloat
```
**File:** `my_asm_decl.go` (only included for amd64)

#### ✅ Correct Pattern 2: Fallback implementation for generic platforms
```go
//go:build !amd64 && !arm64

package bigmath

func myAsmFunction(x *BigFloat) *BigFloat {
    return myGenericFunction(x)
}
```
**File:** `my_fallback_generic.go` (only included for generic platforms)

#### ❌ Wrong Pattern: No build tags, no fallback
```go
package bigmath

//go:noescape
func myAsmFunction(x *BigFloat) *BigFloat  // Missing body - fails on generic platforms!
```

### Current Files Following This Pattern:

- `mpn.go` - `//go:build amd64 || arm64` (assembly implementations exist)
- `mpn_fallback_generic_ops.go` - `//go:build !amd64 && !arm64` (fallback implementations)
- `bigfloat_ops_decl.go` - `//go:build amd64` (assembly implementations exist)
- `bigfloat_ops_arm64.go` - `//go:build arm64` (fallback implementations)
- `exp_log_decl.go` - `//go:build !amd64 && !arm64` (fallback implementations)
- `rounding_asm_decl.go` - `//go:build amd64` (assembly implementations exist)
- `extended_precision_decl.go` - `//go:build amd64 || 386` (x87 FPU assembly implementations)
- `extended_precision_amd64.s` - `//go:build amd64 || 386` (x87 FPU assembly code)
- `extended_precision_fallback.go` - `//go:build !amd64 && !386` (fallback to standard math)

### Testing

Always test builds on multiple platforms:
```bash
# Test generic platform (e.g., s390x)
GOOS=linux GOARCH=s390x go build ./...

# Test amd64
GOOS=linux GOARCH=amd64 go build ./...

# Test arm64
GOOS=linux GOARCH=arm64 go build ./...
```

### Prevention

Before committing:
1. Check if new assembly declarations have proper build tags
2. Verify fallback implementations exist for generic platforms
3. Test build on at least one generic platform (s390x, ppc64le, etc.)

