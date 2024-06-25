package hnet

import (
	"context"
	"io"
	"net"
	"net/netip"
	"time"
)

type (
	WriterTo interface {
		WriteTo(p []byte, addr net.Addr) (int, error)
	}

	WriterToUDP interface {
		WriteToUDP(p []byte, addr *net.UDPAddr) (int, error)
	}

	WriterToUDPAddrPort interface {
		WriteToUDPAddrPort(p []byte, addr netip.AddrPort) (int, error)
	}

	WriterMsgUDP interface {
		WriteMsgUDP(p, oob []byte, addr *net.UDPAddr) (int, int, error)
	}

	WriterMsgUDPAddrPort interface {
		WriteMsgUDPAddrPort(p, oob []byte, addr netip.AddrPort) (int, int, error)
	}
)

func Write(ctx context.Context, w io.Writer, p []byte) (int, error) {
	d, ok := w.(interface {
		SetWriteDeadline(time.Time) error
	})
	if !ok {
		return w.Write(p)
	}

	defer Stopper(ctx, d.SetWriteDeadline)()

	n, err := w.Write(p)
	err = FixError(ctx, err)

	return n, err
}

func WriteTo(ctx context.Context, w WriterTo, p []byte, addr net.Addr) (int, error) {
	d, ok := w.(interface {
		SetWriteDeadline(time.Time) error
	})
	if !ok {
		return w.WriteTo(p, addr)
	}

	defer Stopper(ctx, d.SetWriteDeadline)()

	n, err := w.WriteTo(p, addr)
	err = FixError(ctx, err)

	return n, err
}

func WriteToUDP(ctx context.Context, w WriterToUDP, p []byte, addr *net.UDPAddr) (int, error) {
	d, ok := w.(interface {
		SetWriteDeadline(time.Time) error
	})
	if !ok {
		return w.WriteToUDP(p, addr)
	}

	defer Stopper(ctx, d.SetWriteDeadline)()

	n, err := w.WriteToUDP(p, addr)
	err = FixError(ctx, err)

	return n, err
}

func WriteToUDPAddrPort(ctx context.Context, w WriterToUDPAddrPort, p []byte, addr netip.AddrPort) (int, error) {
	d, ok := w.(interface {
		SetWriteDeadline(time.Time) error
	})
	if !ok {
		return w.WriteToUDPAddrPort(p, addr)
	}

	defer Stopper(ctx, d.SetWriteDeadline)()

	n, err := w.WriteToUDPAddrPort(p, addr)
	err = FixError(ctx, err)

	return n, err
}

func WriteMsgUDP(ctx context.Context, w WriterMsgUDP, p, oob []byte, addr *net.UDPAddr) (int, int, error) {
	d, ok := w.(interface {
		SetWriteDeadline(time.Time) error
	})
	if !ok {
		return w.WriteMsgUDP(p, oob, addr)
	}

	defer Stopper(ctx, d.SetWriteDeadline)()

	n, oobn, err := w.WriteMsgUDP(p, oob, addr)
	err = FixError(ctx, err)

	return n, oobn, err
}

func WriteMsgUDPAddrPort(ctx context.Context, w WriterMsgUDPAddrPort, p, oob []byte, addr netip.AddrPort) (int, int, error) {
	d, ok := w.(interface {
		SetWriteDeadline(time.Time) error
	})
	if !ok {
		return w.WriteMsgUDPAddrPort(p, oob, addr)
	}

	defer Stopper(ctx, d.SetWriteDeadline)()

	n, oobn, err := w.WriteMsgUDPAddrPort(p, oob, addr)
	err = FixError(ctx, err)

	return n, oobn, err
}
