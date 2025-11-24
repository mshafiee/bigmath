//go:build !arm64

package bigmath

// Assembly function declarations
// These will be implemented in architecture-specific .s files (e.g., AMD64).

//go:noescape
func bigfloatAddAsm(a, b *BigFloat, prec uint) *BigFloat

//go:noescape
func bigfloatSubAsm(a, b *BigFloat, prec uint) *BigFloat

//go:noescape
func bigfloatMulAsm(a, b *BigFloat, prec uint) *BigFloat

//go:noescape
func bigfloatDivAsm(a, b *BigFloat, prec uint) *BigFloat

//go:noescape
func bigfloatSqrtAsm(x *BigFloat, prec uint) *BigFloat
