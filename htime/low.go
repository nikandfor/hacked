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

//go:linkname startTimer time.startTimer
//func startTimer(*runtimeTimer)

////go:linkname when time.when
//func when(time.Duration) int64

//go:linkname runtimeNano runtime.nanotime
func runtimeNano() int64

// when is a helper function for setting the 'when' field of a runtimeTimer.
// It returns what the time will be, in nanoseconds, Duration d in the future.
// If d is negative, it is ignored. If the returned value would be less than
// zero because of an overflow, MaxInt64 is returned.
func when(d time.Duration) int64 {
	if d <= 0 {
		return runtimeNano()
	}
	t := runtimeNano() + int64(d)
	if t < 0 {
		// N.B. runtimeNano() and d are always positive, so addition
		// (including overflow) will never result in t == 0.
		t = 1<<63 - 1 // math.MaxInt64
	}
	return t
}
