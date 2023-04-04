package hnet

import (
	"context"
	"io"
	"net"
	"net/netip"
	"time"
)

type (
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

func Accept(ctx context.Context, l net.Listener) (net.Conn, error) {
	d, ok := l.(interface {
		SetDeadline(time.Time) error
	})

	if !ok {
		return l.Accept()
	}

	stopc := make(chan struct{})
	defer close(stopc)

	go func() {
		select {
		case <-ctx.Done():
		case <-stopc:
			return
		}

		_ = d.SetDeadline(time.Unix(1, 0))
	}()

	c, err := l.Accept()
	if c != nil || !isTimeout(err) {
		return c, err
	}

	select {
	case <-ctx.Done():
		err = ctx.Err()
	default:
	}

	return nil, err
}

func Read(ctx context.Context, r io.Reader, p []byte) (int, error) {
	d, ok := r.(interface {
		SetReadDeadline(time.Time) error
	})

	if !ok {
		return r.Read(p)
	}

	stopc := make(chan struct{})
	defer close(stopc)

	go func() {
		select {
		case <-ctx.Done():
		case <-stopc:
			return
		}

		_ = d.SetReadDeadline(time.Unix(1, 0))
	}()

	n, err := r.Read(p)

	if isTimeout(err) {
		select {
		case <-ctx.Done():
			err = ctx.Err()
		default:
		}
	}

	return n, err
}

func ReadFrom(ctx context.Context, r ReaderFrom, p []byte) (int, net.Addr, error) {
	d, ok := r.(interface {
		SetReadDeadline(time.Time) error
	})

	if !ok {
		return r.ReadFrom(p)
	}

	stopc := make(chan struct{})
	defer close(stopc)

	go func() {
		select {
		case <-ctx.Done():
		case <-stopc:
			return
		}

		_ = d.SetReadDeadline(time.Unix(1, 0))
	}()

	n, addr, err := r.ReadFrom(p)

	if isTimeout(err) {
		select {
		case <-ctx.Done():
			err = ctx.Err()
		default:
		}
	}

	return n, addr, err
}

func ReadFromUDP(ctx context.Context, r ReaderFromUDP, p []byte) (int, *net.UDPAddr, error) {
	d, ok := r.(interface {
		SetReadDeadline(time.Time) error
	})

	if !ok {
		return r.ReadFromUDP(p)
	}

	stopc := make(chan struct{})
	defer close(stopc)

	go func() {
		select {
		case <-ctx.Done():
		case <-stopc:
			return
		}

		_ = d.SetReadDeadline(time.Unix(1, 0))
	}()

	n, addr, err := r.ReadFromUDP(p)

	if isTimeout(err) {
		select {
		case <-ctx.Done():
			err = ctx.Err()
		default:
		}
	}

	return n, addr, err
}

func ReadFromUDPAddrPort(ctx context.Context, r ReaderFromUDPAddrPort, p []byte) (int, netip.AddrPort, error) {
	d, ok := r.(interface {
		SetReadDeadline(time.Time) error
	})

	if !ok {
		return r.ReadFromUDPAddrPort(p)
	}

	stopc := make(chan struct{})
	defer close(stopc)

	go func() {
		select {
		case <-ctx.Done():
		case <-stopc:
			return
		}

		_ = d.SetReadDeadline(time.Unix(1, 0))
	}()

	n, addr, err := r.ReadFromUDPAddrPort(p)

	if isTimeout(err) {
		select {
		case <-ctx.Done():
			err = ctx.Err()
		default:
		}
	}

	return n, addr, err
}

func ReadMsgUDP(ctx context.Context, r ReaderMsgUDP, p, oob []byte) (n, oobn, flags int, addr *net.UDPAddr, err error) {
	d, ok := r.(interface {
		SetReadDeadline(time.Time) error
	})

	if !ok {
		return r.ReadMsgUDP(p, oob)
	}

	stopc := make(chan struct{})
	defer close(stopc)

	go func() {
		select {
		case <-ctx.Done():
		case <-stopc:
			return
		}

		_ = d.SetReadDeadline(time.Unix(1, 0))
	}()

	n, oobn, flags, addr, err = r.ReadMsgUDP(p, oob)

	if isTimeout(err) {
		select {
		case <-ctx.Done():
			err = ctx.Err()
		default:
		}
	}

	return
}

func ReadMsgUDPAddrPort(ctx context.Context, r ReaderMsgUDPAddrPort, p, oob []byte) (n, oobn, flags int, addr netip.AddrPort, err error) {
	d, ok := r.(interface {
		SetReadDeadline(time.Time) error
	})

	if !ok {
		return r.ReadMsgUDPAddrPort(p, oob)
	}

	stopc := make(chan struct{})
	defer close(stopc)

	go func() {
		select {
		case <-ctx.Done():
		case <-stopc:
			return
		}

		_ = d.SetReadDeadline(time.Unix(1, 0))
	}()

	n, oobn, flags, addr, err = r.ReadMsgUDPAddrPort(p, oob)

	if isTimeout(err) {
		select {
		case <-ctx.Done():
			err = ctx.Err()
		default:
		}
	}

	return
}

func isTimeout(err error) bool {
	to, ok := err.(interface{ Timeout() bool })

	return ok && to.Timeout()
}
