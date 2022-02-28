package runtime

import (
	"hash/crc32"
	"math/rand"
	"testing"
	"time"
)

func BenchmarkCRC32(b *testing.B) {
	buf := make([]byte, 1024)

	for i := range buf {
		buf[i] = byte(i)
	}

	for i := 0; i < b.N; i++ {
		_ = crc32.ChecksumIEEE(buf)
	}

	b.SetBytes(int64(len(buf)))
}

func BenchmarkBytesHash(b *testing.B) {
	buf := make([]byte, 1024)

	for i := range buf {
		buf[i] = byte(i)
	}

	for i := 0; i < b.N; i++ {
		_ = BytesHash(buf, 0)
	}

	b.SetBytes(int64(len(buf)))
}

var rr uint32

func BenchmarkRand(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < b.N; i++ {
		rr = r.Uint32()
	}
}

func BenchmarkFastRand(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rr = Fastrand()
	}
}
