package low

import (
	"errors"
	"io"
)

const Spaces = "                                                                                                                                "

type (
	// Buf is a buffer that implements a few Writer/Reader interfaces.
	// Mostly used for tests.
	Buf []byte
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

func (w *Buf) ReadFrom(r io.Reader) (int64, error) {
	if wt, ok := r.(io.WriterTo); ok {
		return wt.WriteTo(w)
	}

	st := len(*w)
	end := 0

	for {
		*w = append((*w)[:end], make([]byte, 512)...)

		n, err := r.Read((*w)[end:cap(*w)])
		end += n
		*w = (*w)[:end]

		if errors.Is(err, io.EOF) {
			return int64(end - st), nil
		}
		if err != nil {
			return int64(end - st), err
		}
	}
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

func (w *Buf) Reset()       { *w = (*w)[:0] }
func (w Buf) Len() int      { return len(w) }
func (w Buf) LenF() float64 { return float64(w.Len()) }
func (w Buf) Bytes() []byte { return w }
