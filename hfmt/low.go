package hfmt

import (
	"fmt"
	"io"
	"reflect"
	"strconv"
	"unsafe"
)

// nolint
type (
	iface struct {
		typ, word unsafe.Pointer
	}

	// pp is used to store a printer's state and is reused with sync.Pool to avoid allocations.
	pp struct {
		buf []byte

		// arg holds the current item, as an interface{}.
		arg interface{}

		// value is used instead of arg for reflect values.
		value reflect.Value

		// fmt is used to format basic items such as integers or strings.
		fmt fmtt

		// reordered records whether the format string used argument reordering.
		reordered bool
		// goodArgNum records whether the most recent reordering directive was valid.
		goodArgNum bool
		// panicking is set by catchPanic to avoid infinite panic, recover, panic, ... recursion.
		panicking bool
		// erroring is set when printing an error string to guard against calling handleMethods.
		erroring bool
		// wrapErrs is set when the format string may contain a %w verb.
		wrapErrs bool
		// wrappedErr records the target of the %w verb.
		wrappedErr error
	}

	// A fmt is the raw formatter used by Printf etc.
	// It prints into a buffer that must be set up separately.
	fmtt struct {
		buf *[]byte

		fmtFlags

		wid  int // width
		prec int // precision

		// intbuf is large enough to store %b of an int64 with a sign and
		// avoids padding at the end of the struct on 32 bit architectures.
		intbuf [68]byte
	}

	// flags placed in a separate struct for easy clearing.
	fmtFlags struct {
		widPresent  bool
		precPresent bool
		minus       bool
		plus        bool
		sharp       bool
		space       bool
		zero        bool

		// For the formats %+v %#v, we set the plusV/sharpV flags
		// and clear the plus/sharp flags since %+v and %#v are in effect
		// different, flagless formats set at the top level.
		plusV  bool
		sharpV bool
	}

	formatter struct{}
)

var ppType unsafe.Pointer

func init() {
	fmt.Fprintf(io.Discard, "%v", formatter{})
}

// Appendf is similar to fmt.Fprintf but a little bit hacked.
//
// There is no sync.Pool.Get and Put. There is no copying buffer to io.Writer or conversion to string. There is no io.Writer interface dereference.
// All that gives advantage about 30-50 ns per call. Yes, I know :).
func Appendf(b []byte, format string, a ...interface{}) []byte {
	return fmt.Appendf(b, format, a...)
}

// Appendln is similar to fmt.Sprintln but faster. See doc for Appendf for more details.
func Appendln(b []byte, a ...interface{}) []byte {
	return fmt.Appendln(b, a...)
}

// Append is similar to fmt.Sprint but faster. See doc for Appendf for more details.
func Append(b []byte, a ...interface{}) []byte {
	return fmt.Append(b, a...)
}

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

	fmt.Fprintf(s, bytesToString(buf[:i]), arg)
}

//go:linkname doPrintf fmt.(*pp).doPrintf
//go:noescape
func doPrintf(p *pp, format string, a []interface{})

//go:linkname doPrintln fmt.(*pp).doPrintln
//go:noescape
func doPrintln(p *pp, a []interface{})

//go:linkname doPrint fmt.(*pp).doPrint
//go:noescape
func doPrint(p *pp, a []interface{})

//go:linkname newPrinter fmt.newPrinter
//go:noescape
func newPrinter() unsafe.Pointer

//go:linkname ppFree fmt.(*pp).free
//go:noescape
func ppFree(unsafe.Pointer)

//go:linkname printArg fmt.(*pp).printArg
//go:noescape
func printArg(p unsafe.Pointer, arg interface{}, verb rune)

func (formatter) Format(s fmt.State, c rune) {
	i := *(*iface)(unsafe.Pointer(&s))

	ppType = i.typ
}

// noescape hides a pointer from escape analysis.  noescape is
// the identity function but escape analysis doesn't think the
// output depends on the input.  noescape is inlined and currently
// compiles down to zero instructions.
// USE CAREFULLY!
//
//go:nosplit
func noescape(p unsafe.Pointer) unsafe.Pointer {
	x := uintptr(p)
	return unsafe.Pointer(x ^ 0) //nolint:staticcheck
}

func bytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
