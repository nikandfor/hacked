package hfmt

import (
	"fmt"
	"strconv"
	"unsafe"
)

// nolint
type (
	iface struct {
		typ, word unsafe.Pointer
	}

	formatter struct{}
)

func PrintArg(s fmt.State, arg interface{}, verb rune) {
	var buf [64]byte

	i := 0

	buf[i] = '%'
	i++

	for _, f := range "-+# 0" {
		if s.Flag(int(f)) {
			buf[i] = byte(f)
			i++
		}
	}

	if w, ok := s.Width(); ok {
		q := strconv.AppendInt(buf[:i], int64(w), 10)
		i = len(q)
	}

	if p, ok := s.Precision(); ok {
		buf[i] = '.'
		i++

		q := strconv.AppendInt(buf[:i], int64(p), 10)
		i = len(q)
	}

	buf[i] = byte(verb)
	i++

	_, _ = fmt.Fprintf(s, bytesToString(buf[:i]), arg)
}

//go:linkname printArg fmt.(*pp).printArg
//go:noescape
func printArg(p unsafe.Pointer, arg interface{}, verb rune)

// noescape hides a pointer from escape analysis.  noescape is
// the identity function but escape analysis doesn't think the
// output depends on the input.  noescape is inlined and currently
// compiles down to zero instructions.
// USE CAREFULLY!
//
//go:nosplit
func noescape(p unsafe.Pointer) unsafe.Pointer {
	x := uintptr(p)
	return unsafe.Pointer(x ^ 0) //nolint:staticcheck,govet
}

func bytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
