package htime

import (
	"time"
	_ "unsafe"
)

//go:linkname Now time.now
func Now() (sec int64, nsec int32, mono int64)

func UnixNano() (t int64) {
	t, nsec, _ := Now()

	return t*1e9 + int64(nsec)
}

func Monotonic() (c int64) {
	_, _, c = Now()

	return
}

func MonotonicOf(t time.Time) int64 {
	return mono(&t)
}

// DateClock is faster version of t.Date(); t.Clock().
func DateClock(t time.Time) (year, month, day, hour, min, sec int) { //nolint:gocritic
	u := timeAbs(t)
	year, month, day, _ = absDate(u, true)
	hour, min, sec = absClock(u)
	return
}

//go:linkname timeAbs time.Time.abs
func timeAbs(time.Time) uint64

//go:linkname absClock time.absClock
func absClock(uint64) (hour, min, sec int)

//go:linkname absDate time.absDate
func absDate(uint64, bool) (year, month, day, yday int)

//go:noescape
//go:linkname mono time.(*Time).mono
func mono(*time.Time) int64
