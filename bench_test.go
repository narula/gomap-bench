package cmap

import (
	"runtime"
	"sync"
	"testing"
)

func BenchmarkGoMapWriteEmpty(b *testing.B) {
	m := make(map[Key]Value)
	var buf [16]byte
	for i := 0; i < b.N; i++ {
		m[reuse_key(uint64(i), &buf)] = i
	}
}

func BenchmarkGoMapWriteFixed(b *testing.B) {
	b.StopTimer()
	m := make(map[Key]Value)
	var buf [16]byte
	for i := 0; i < 1000000; i++ {
		m[reuse_key(uint64(i), &buf)] = i
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		m[reuse_key(uint64(i%1000000), &buf)] = i
	}
}

func BenchmarkGoMapReadFixed(b *testing.B) {
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

func BenchmarkGoMapConcurrent2(b *testing.B) {
	benchmarkGoMapReadConcurrent(b, 2)
}

func BenchmarkGoMapConcurrent4(b *testing.B) {
	benchmarkGoMapReadConcurrent(b, 4)
}

func BenchmarkGoMapConcurrent8(b *testing.B) {
	benchmarkGoMapReadConcurrent(b, 8)
}

func BenchmarkGoMapConcurrent12(b *testing.B) {
	benchmarkGoMapReadConcurrent(b, 12)
}

func BenchmarkGoMapWriteConcurrent1(b *testing.B) {
	benchmarkGoMapWriteConcurrent(b, 1)
}

func BenchmarkGoMapWriteConcurrent2(b *testing.B) {
	benchmarkGoMapWriteConcurrent(b, 2)
}

func BenchmarkGoMapWriteConcurrent4(b *testing.B) {
	benchmarkGoMapWriteConcurrent(b, 4)
}

func BenchmarkGoMapWriteConcurrent8(b *testing.B) {
	benchmarkGoMapWriteConcurrent(b, 8)
}

func BenchmarkGoMapWriteConcurrent12(b *testing.B) {
	benchmarkGoMapWriteConcurrent(b, 12)
}

func benchmarkGoMapReadConcurrent(b *testing.B, nprocs int) {
	b.StopTimer()
	runtime.GOMAXPROCS(nprocs)
	m := make(map[Key]Value)
	var rw sync.RWMutex
	var buf [16]byte
	for i := 0; i < 1000000; i++ {
		m[reuse_key(uint64(i), &buf)] = i
	}
	b.StartTimer()
	var wg sync.WaitGroup
	for j := 0; j < nprocs; j++ {
		wg.Add(1)
		go func() {
			for i := 0; i < b.N; i++ {
				rw.RLock()
				x, _ := m[reuse_key(uint64(i), &buf)]
				_ = x
				rw.RUnlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func benchmarkGoMapWriteConcurrent(b *testing.B, nprocs int) {
	b.StopTimer()
	runtime.GOMAXPROCS(nprocs)
	m := make(map[Key]Value)
	var rw sync.RWMutex
	var buf [16]byte
	for i := 0; i < 1000000; i++ {
		m[reuse_key(uint64(i), &buf)] = i
	}
	b.StartTimer()
	var wg sync.WaitGroup
	for j := 0; j < nprocs; j++ {
		wg.Add(1)
		go func() {
			for i := 0; i < b.N; i++ {
				rw.Lock()
				m[reuse_key(uint64(i%1000000), &buf)] = i
				rw.Unlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
