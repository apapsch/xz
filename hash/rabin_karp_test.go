package hash

import (
	"math/rand"
	"testing"
)

func TestRabinKarpSimple(t *testing.T) {
	p := []byte("abcde")
	r := NewRabinKarp(4)
	h2 := ComputeHashes(r, p)
	for i, h := range h2 {
		w := ComputeHashes(r, p[i:i+4])[0]
		t.Logf("%d h=%#016x w=%#016x", i, h, w)
		if h != w {
			t.Errorf("rolling hash %d: %#016x; want %#016x",
				i, h, w)
		}
	}
}

func makeBenchmarkBytes(n int) []byte {
	rnd := rand.New(rand.NewSource(42))
	p := make([]byte, n)
	for i := range p {
		p[i] = byte(rnd.Uint32())
	}
	return p
}

func BenchmarkRabinKarp(b *testing.B) {
	p := makeBenchmarkBytes(4096)
	r := NewRabinKarp(4)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ComputeHashes(r, p)
	}
}
