package htime

import (
	"testing"
	"time"
)

func TestAfterFuncSync(t *testing.T) {
	c := make(chan struct{})

	setOk := func() {
		close(c)
	}

	tt := AfterFuncSync(time.Millisecond, setOk)
	defer tt.Stop()

	<-c
}

func BenchmarkTimeNow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = time.Now()
	}
}

func BenchmarkNow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _, _ = Now()
	}
}

func BenchmarkTimeNowUnixNano(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = time.Now().UnixNano()
	}
}

func BenchmarkUnixNano(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = UnixNano()
	}
}

func BenchmarkDateAndClock(b *testing.B) {
	t := time.Now()

	for i := 0; i < b.N; i++ {
		_, _, _ = t.Date()
		_, _, _ = t.Clock()
	}
}

func BenchmarkDateClock(b *testing.B) {
	t := time.Now()

	for i := 0; i < b.N; i++ {
		_, _, _, _, _, _ = DateClock(t)
	}
}
