package hnet

import (
	"context"
	"io"
	"net"
	"net/netip"
	"sync"
	"time"
)

type (
	StoppableListener struct {
		context.Context
		net.Listener
	}

	StoppableConn struct {
		context.Context
		net.Conn
	}

	ReaderFrom interface {
		ReadFrom(p []byte) (int, net.Addr, error)
	}

	ReaderFromUDP interface {
		ReadFromUDP(p []byte) (int, *net.UDPAddr, error)
	}

	ReaderFromUDPAddrPort interface {
		ReadFromUDPAddrPort(p []byte) (int, netip.AddrPort, error)
	}

	ReaderMsgUDP interface {
		ReadMsgUDP(p, oob []byte) (n, oobn, flags int, addr *net.UDPAddr, err error)
	}

	ReaderMsgUDPAddrPort interface {
		ReadMsgUDPAddrPort(p, oob []byte) (n, oobn, flags int, addr netip.AddrPort, err error)
	}
)

// Accept wraps l.Accept but aborts the operation on context cancelation.
func Accept(ctx context.Context, l net.Listener) (net.Conn, error) {
	d, ok := l.(interface {
		SetDeadline(time.Time) error
	})
	if !ok {
		return l.Accept()
	}

	defer Stopper(ctx, d.SetDeadline)()

	c, err := l.Accept()
	if c != nil {
		return c, err
	}

	err = FixError(ctx, err)

	return nil, err
}

// Read wraps r.Read but aborts the operation on context cancelation.
func Read(ctx context.Context, r io.Reader, p []byte) (int, error) {
	d, ok := r.(interface {
		SetReadDeadline(time.Time) error
	})
	if !ok {
		return r.Read(p)
	}

	defer Stopper(ctx, d.SetReadDeadline)()

	n, err := r.Read(p)
	err = FixError(ctx, err)

	return n, err
}

// ReadFrom wraps r.ReadFrom but aborts the operation on context cancelation.
func ReadFrom(ctx context.Context, r ReaderFrom, p []byte) (int, net.Addr, error) {
	d, ok := r.(interface {
		SetReadDeadline(time.Time) error
	})
	if !ok {
		return r.ReadFrom(p)
	}

	defer Stopper(ctx, d.SetReadDeadline)()

	n, addr, err := r.ReadFrom(p)
	err = FixError(ctx, err)

	return n, addr, err
}

// ReadFromUDP wraps r.ReadFromUDP but aborts the operation on context cancelation.
func ReadFromUDP(ctx context.Context, r ReaderFromUDP, p []byte) (int, *net.UDPAddr, error) {
	d, ok := r.(interface {
		SetReadDeadline(time.Time) error
	})
	if !ok {
		return r.ReadFromUDP(p)
	}

	defer Stopper(ctx, d.SetReadDeadline)()

	n, addr, err := r.ReadFromUDP(p)
	err = FixError(ctx, err)

	return n, addr, err
}

// ReadFromUDPAddrPort wraps r.ReadFromUDPAddrPort but aborts the operation on context cancelation.
func ReadFromUDPAddrPort(ctx context.Context, r ReaderFromUDPAddrPort, p []byte) (int, netip.AddrPort, error) {
	d, ok := r.(interface {
		SetReadDeadline(time.Time) error
	})
	if !ok {
		return r.ReadFromUDPAddrPort(p)
	}

	defer Stopper(ctx, d.SetReadDeadline)()

	n, addr, err := r.ReadFromUDPAddrPort(p)
	err = FixError(ctx, err)

	return n, addr, err
}

// ReadMsgUDP wraps r.ReadMsgUDP but aborts the operation on context cancelation.
func ReadMsgUDP(ctx context.Context, r ReaderMsgUDP, p, oob []byte) (n, oobn, flags int, addr *net.UDPAddr, err error) {
	d, ok := r.(interface {
		SetReadDeadline(time.Time) error
	})
	if !ok {
		return r.ReadMsgUDP(p, oob)
	}

	defer Stopper(ctx, d.SetReadDeadline)()

	n, oobn, flags, addr, err = r.ReadMsgUDP(p, oob)
	err = FixError(ctx, err)

	return
}

// ReadMsgUDPAddrPort wraps r.ReadMsgUDPAddrPort but aborts the operation on context cancelation.
func ReadMsgUDPAddrPort(ctx context.Context, r ReaderMsgUDPAddrPort, p, oob []byte) (n, oobn, flags int, addr netip.AddrPort, err error) {
	d, ok := r.(interface {
		SetReadDeadline(time.Time) error
	})
	if !ok {
		return r.ReadMsgUDPAddrPort(p, oob)
	}

	defer Stopper(ctx, d.SetReadDeadline)()

	n, oobn, flags, addr, err = r.ReadMsgUDPAddrPort(p, oob)
	err = FixError(ctx, err)

	return
}

// NewStoppableListener wraps listener to have the same interface,
// but returned listener aborts on context cancellation.
func NewStoppableListener(ctx context.Context, l net.Listener) net.Listener {
	return StoppableListener{
		Context:  ctx,
		Listener: l,
	}
}

func (l StoppableListener) Accept() (net.Conn, error) {
	return Accept(l.Context, l.Listener)
}

// NewStoppableConn wraps connection to have the same interface,
// but returned connection aborts on context cancellation.
func NewStoppableConn(ctx context.Context, c net.Conn) net.Conn {
	return StoppableConn{
		Context: ctx,
		Conn:    c,
	}
}

func (c StoppableConn) Read(p []byte) (n int, err error) {
	defer Stopper(c.Context, c.Conn.SetReadDeadline)()

	n, err = c.Conn.Read(p)
	err = FixError(c.Context, err)

	return
}

func (c StoppableConn) Write(p []byte) (n int, err error) {
	defer Stopper(c.Context, c.Conn.SetWriteDeadline)()

	n, err = c.Conn.Write(p)
	err = FixError(c.Context, err)

	return
}

// Stopper is a helper function which calls dead if context is canceled.
// Returned function must be called with defer to release goroutine used internally.
func Stopper(ctx context.Context, dead func(time.Time) error) func() {
	donec := make(chan struct{})

	var mu sync.Mutex
	var killed bool

	go func() {
		select {
		case <-ctx.Done():
		case <-donec:
			return
		}

		defer mu.Unlock()
		mu.Lock()

		select {
		case <-donec:
			return
		default:
		}

		_ = dead(time.Unix(1, 0))
		killed = true
	}()

	return func() {
		close(donec)

		defer mu.Unlock()
		mu.Lock()

		if killed {
			_ = dead(time.Time{})
		}
	}
}

func isTimeout(err error) bool {
	to, ok := err.(interface{ Timeout() bool })

	return ok && to.Timeout()
}

// FixError replaces internal error caused by operation abortion
// and replaces it with one returned by ctx.Err if context was canceled.
// Otherwise it returns the error unchanged.
func FixError(ctx context.Context, err error) error {
	if isTimeout(err) {
		select {
		case <-ctx.Done():
			err = ctx.Err()
		default:
		}
	}

	return err
}
