package hunsafe

import "unsafe"

func NoEscapeBuffer(b []byte) []byte {
	return *(*[]byte)(NoEscape(unsafe.Pointer(&b)))
}

func NoEscapeSlice[T any](b []T) []T {
	return *(*[]T)(NoEscape(unsafe.Pointer(&b)))
}

// noescape hides a pointer from escape analysis.  noescape is
// the identity function but escape analysis doesn't think the
// output depends on the input.  noescape is inlined and currently
// compiles down to zero instructions.
// USE CAREFULLY!
//
//go:nosplit
func NoEscape(p unsafe.Pointer) unsafe.Pointer {
	x := uintptr(p)
	return unsafe.Pointer(x ^ 0) //nolint:staticcheck,govet
}
