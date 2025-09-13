package low

import (
	"fmt"
	"io"
	"unicode/utf8"
)

type (
	BufReader struct {
		Buf
		R int
	}
)

func (r *BufReader) Read(p []byte) (n int, err error) {
	n = copy(p, r.Buf[r.R:])
	r.R += n

	if r.R == len(r.Buf) {
		err = io.EOF
	}

	return
}

func (r *BufReader) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(r.Buf[r.R:])
	r.R += n

	return int64(n), err
}

func (r *BufReader) ReadByte() (byte, error) {
	if r.R >= len(r.Buf) {
		return 0, io.EOF
	}

	r.R++

	return r.Buf[r.R-1], nil
}

func (r *BufReader) UnreadByte() error {
	if r.R == 0 {
		return fmt.Errorf("unread before the start")
	}

	r.R--
	return nil
}

func (r *BufReader) ReadRune() (rune, int, error) {
	if r.R >= len(r.Buf) {
		return 0, 0, io.EOF
	}
	if !utf8.FullRune(r.Buf[r.R:]) {
		return 0, 0, io.ErrUnexpectedEOF
	}

	rr, size := utf8.DecodeRune(r.Buf[r.R:])
	r.R += size

	return rr, size, nil
}

func (r *BufReader) Reset()       { r.R = 0 }
func (r BufReader) Len() int      { return r.Buf.Len() - r.R }
func (r BufReader) LenF() float64 { return float64(r.Len()) }
func (r BufReader) Bytes() []byte { return r.Buf[r.R:] }
