package low

import "io"

const Spaces = "                                                                                                                                "

type (
	Buf []byte

	BufReader struct {
		Buf
		R int
	}
)

func (w *Buf) Write(p []byte) (int, error) {
	*w = append(*w, p...)

	return len(p), nil
}

func (w *Buf) WriteAt(p []byte, off int64) (int, error) {
	if int(off)+len(p) <= len(*w) {
		return copy((*w)[off:], p), nil
	}

	for cap(*w) < int(off) {
		*w = append((*w)[:cap(*w)], 0, 0, 0, 0, 0, 0, 0, 0)
	}

	*w = (*w)[:off]
	*w = append(*w, p...)

	return len(p), nil
}

func (w *Buf) NewLine() {
	l := len(*w)
	if l == 0 || (*w)[l-1] != '\n' {
		*w = append(*w, '\n')
	}
}

func (w Buf) ReadAt(p []byte, off int64) (int, error) {
	if off >= int64(len(w)) {
		return 0, io.EOF
	}

	n := copy(p, w[off:])
	if n < len(p) {
		return n, io.EOF
	}

	return n, nil
}

func (w *Buf) Reset()        { *w = (*w)[:0] }
func (w *Buf) Len() int      { return len(*w) }
func (w *Buf) LenF() float64 { return float64(w.Len()) }
func (w *Buf) Bytes() []byte { return *w }

func (r *BufReader) Read(p []byte) (n int, err error) {
	n = copy(p, r.Buf[r.R:])
	r.R += n

	if r.R == len(r.Buf) {
		err = io.EOF
	}

	return
}

func (r *BufReader) Reset()        { r.R = 0 }
func (r *BufReader) Len() int      { return r.Buf.Len() - r.R }
func (r *BufReader) LenF() float64 { return float64(r.Len()) }
func (r *BufReader) Bytes() []byte { return r.Buf[r.R:] }
