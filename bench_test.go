package cmap

import "testing"

func BenchmarkGoMapSimple(b *testing.B) {
	m := make(map[Key]Value)
	var buf [16]byte
	for i := 0; i < b.N; i++ {
		m[reuse_key(uint64(i), &buf)] = i
	}
}

func BenchmarkGoMapRead(b *testing.B) {
	b.StopTimer()
	m := make(map[Key]Value)
	var buf [16]byte
	for i := 0; i < 1000000; i++ {
		m[reuse_key(uint64(i), &buf)] = i
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		x, _ := m[reuse_key(uint64(i), &buf)]
		_ = x
	}
}
