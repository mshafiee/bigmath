//go:build !arm64

package bigmath

//go:noescape
func bigExpAsm(x *BigFloat, prec uint) *BigFloat

//go:noescape
func bigLogAsm(x *BigFloat, prec uint) *BigFloat
