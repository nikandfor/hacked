package htime

import (
	"time"
	"unsafe"
	_ "unsafe"
)

type (
	timer struct {
		C <-chan time.Time
		r runtimeTimer
	}

	// Interface to timers implemented in package runtime.
	// Must be in sync with ../runtime/time.go:/^type timer
	runtimeTimer struct {
		pp       uintptr
		when     int64
		period   int64
		f        func(interface{}, uintptr) // NOTE: must not be closure
		arg      interface{}
		seq      uintptr
		nextwhen int64
		status   uint32
	}
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

// DateClock is faster version of t.Date(); t.Clock().
func DateClock(t time.Time) (year, month, day, hour, min, sec int) { //nolint:gocritic
	u := timeAbs(t)
	year, month, day, _ = absDate(u, true)
	hour, min, sec = absClock(u)
	return
}

func NewTimerFunc(d time.Duration, f func()) *time.Timer {
	t := &timer{
		r: runtimeTimer{
			when: when(d),
			f:    doSync,
			arg:  f,
		},
	}

	startTimer(&t.r)

	return (*time.Timer)(unsafe.Pointer(t))
}

//go:linkname AfterFuncSync NewTimerFunc

func AfterFuncSync(d time.Duration, f func()) *time.Timer

func doSync(arg interface{}, seq uintptr) {
	arg.(func())()
}

//go:linkname timeAbs time.Time.abs
func timeAbs(time.Time) uint64

//go:linkname absClock time.absClock
func absClock(uint64) (hour, min, sec int)

//go:linkname absDate time.absDate
func absDate(uint64, bool) (year, month, day, yday int)

//go:linkname startTimer time.startTimer
func startTimer(*runtimeTimer)

//go:linkname when time.when
func when(time.Duration) int64
