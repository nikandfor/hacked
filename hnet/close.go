package hnet

import (
	"fmt"
	"io"
)

// CloserOnErr closes resource on exit if error happened.
//
//	defer CloserOnErr(c, &err)
func CloserOnErr(c io.Closer, errp *error) {
	if *errp == nil {
		return
	}

	_ = c.Close()
}

// Closer is for closing resource on exit and handling an error.
//
//	defer Closer(c, &err, "close conn")
func Closer(c io.Closer, errp *error, msg string) {
	CloserFunc(c.Close, errp, msg)
}

// CloserFunc is for closing resource on exit and handling an error.
//
//	defer CloserFunc(s.Finish, &err, "finish something")
func CloserFunc(c func() error, errp *error, msg string) {
	err := c()
	if err != nil && *errp == nil {
		*errp = fmt.Errorf("%v: %w", msg, err)
	}
}

// CloserWriter closes writer and handles an error.
// It calls CloseWriter if c has it and handles an error.
// net.TCPConn has CloseWriter for example.
//
//	defer CloserWriter(c, &err, "close writer")
func CloserWriter(c any, errp *error, msg string) {
	cw, ok := c.(interface {
		CloseWrite() error
	})
	if !ok {
		return
	}

	CloserFunc(cw.CloseWrite, errp, msg)
}

// CloseWriter check if CloseWriter exists and calls it.
// net.TCPConn has that method for example.
func CloseWriter(c any) error {
	cw, ok := c.(interface {
		CloseWrite() error
	})
	if !ok {
		return nil
	}

	return cw.CloseWrite()
}
