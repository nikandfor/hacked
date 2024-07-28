package hfmt

import (
	"bytes"
	"fmt"
	"testing"
)

type (
	testformatter struct {
		flags [64]bool

		wid int
		prc int

		widok bool
		prcok bool

		verb rune

		//	buf  [100]byte
		//	bufi int
	}
)

func TestSubPrintArg(t *testing.T) {
	var b bytes.Buffer

	var f testformatter

	fmt.Fprintf(&b, "%+012.6q", &f)

	if (testformatter{
		flags: flags("+0"),
		wid:   12,
		widok: true,
		prc:   6,
		prcok: true,
		verb:  'q',
	}) != f {
		t.Errorf("not expected")
	}

	//

	f = testformatter{}
	f2 := testformatter{
		flags: flags("-# "),
		prc:   3,
		prcok: true,
	}

	PrintArg(&f2, &f, 'v')

	if (testformatter{
		flags: flags(" -#"),
		prc:   3,
		prcok: true,
		verb:  'v',
	}) != f {
		t.Errorf("not expected")
	}
}

/*
func BenchmarkPringArg(b *testing.B) {
	b.ReportAllocs()

	f := testformatter{}

	for i := 0; i < b.N; i++ {
		pp := newPrinter()
		printArg(pp, &f, 'v')
		ppFree(pp)
	}
}
*/

func BenchmarkPringArgFallback(b *testing.B) {
	b.ReportAllocs()

	f := testformatter{}
	f2 := testformatter{
		flags: flags("-# "),
		prc:   3,
		prcok: true,
	}

	for i := 0; i < b.N; i++ {
		PrintArg(&f2, &f, 'v')
	}
}

func BenchmarkFmtFprintf(b *testing.B) {
	b.ReportAllocs()

	var buf bytes.Buffer

	for i := 0; i < b.N; i++ {
		buf.Reset()

		fmt.Fprintf(&buf, "message %v %v %v", 1, "string", 3.12)
	}
}

func BenchmarkFmtAppendf(b *testing.B) {
	b.ReportAllocs()

	var buf []byte

	for i := 0; i < b.N; i++ {
		buf = Appendf(buf[:0], "message %v %v %v", 1, "string", 3.12)
	}
}

func (f *testformatter) Format(s fmt.State, verb rune) {
	f.flags = [64]bool{}

	for _, q := range "-+# 0" {
		if s.Flag(int(q)) {
			f.flags[q] = true
		}
	}

	f.wid, f.widok = s.Width()
	f.prc, f.prcok = s.Precision()

	f.verb = verb
}

func (f *testformatter) Flag(c int) bool {
	return f.flags[c]
}

func (f *testformatter) Width() (int, bool) {
	return f.wid, f.widok
}

func (f *testformatter) Precision() (int, bool) {
	return f.prc, f.prcok
}

func (f *testformatter) Write(p []byte) (int, error) {
	//	f.bufi += copy(f.buf[f.bufi:], p)

	return len(p), nil
}

func flags(f string) (r [64]bool) {
	for _, q := range f {
		r[q] = true
	}

	return
}
