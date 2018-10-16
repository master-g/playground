package bench

import (
	"strings"
	"testing"
	"time"
)

func BenchmarkFib(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Fib(20)
	}
}

const m1 = 0x5555555555555555
const m2 = 0x3333333333333333
const m4 = 0x0f0f0f0f0f0f0f0f
const h01 = 0x0101010101010101

func popcnt(x uint64) uint64 {
	x -= (x >> 1) & m1
	x = (x & m2) + ((x >> 2) & m2)
	x = (x + (x >> 4)) & m4
	return (x * h01) >> 56
}

func BenchmarkPopcnt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		popcnt(uint64(i))
	}
}

func SCRaw() string {
	s := "1"
	s += " " + "wow"
	s += time.Now().String()
	return s
}

func SCBuilder(b *strings.Builder) string {
	b.Reset()
	b.WriteString("1")
	b.WriteString(" ")
	b.WriteString(time.Now().String())

	return b.String()
}

var gStr string

func BenchmarkSCRaw(b *testing.B) {
	var s string
	for i := 0; i < b.N; i++ {
		s = SCRaw()
	}
	gStr = s
}

func BenchmarkSCBuilder(b *testing.B) {
	var s string
	builder := &strings.Builder{}
	for i := 0; i < b.N; i++ {
		s = SCBuilder(builder)
	}
	gStr = s
}
