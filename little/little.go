package little

import (
	"fmt"
	"io"
)

func CloserAny(a any, errp *error, msg string, args ...any) {
	c, ok := a.(io.Closer)
	if !ok {
		return
	}

	Closer(c, errp, msg, args...)
}

func Closer(c io.Closer, errp *error, msg string, args ...any) {
	CloserFunc(c.Close, errp, msg, args...)
}

func CloserFunc(c func() error, errp *error, msg string, args ...any) {
	err := c()
	if *errp == nil && err != nil {
		if len(args) != 0 {
			msg = fmt.Sprintf(msg, args...)
		}

		*errp = fmt.Errorf("%v: %w", msg, err)
	}
}

func CloseOnErr(c io.Closer, errp *error) {
	if *errp == nil {
		return
	}

	_ = c.Close()
}

func CloseFuncOnErr(f func() error, errp *error) {
	if *errp == nil {
		return
	}

	_ = f()
}

func CSel[T any](c bool, t, e T) T {
	if c {
		return t
	} else {
		return e
	}
}
