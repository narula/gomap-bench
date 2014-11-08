package cmap

import (
	"runtime"
	"sync"
	"testing"

	"github.com/zond/gotomic"
)

const (
	NUMKEYS = 2 << 20
	WRAPPER = 2<<21 - 1
)

// GOMAP READS

func BenchmarkGoMapReadFixed(b *testing.B) {
	b.StopTimer()
	m := make(map[Key]Value)
	var buf [16]byte
	for i := 0; i < NUMKEYS; i++ {
		m[reuse_key(uint64(i), &buf)] = i
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		x, _ := m[reuse_key(uint64(i&WRAPPER), &buf)]
		_ = x
	}
}

func BenchmarkGoMapReadConcurrentNoLock(b *testing.B) {
	nprocs := runtime.GOMAXPROCS(0)
	b.StopTimer()
	runtime.GOMAXPROCS(nprocs)
	m := make(map[Key]Value)
	var buf [16]byte
	for i := 0; i < NUMKEYS; i++ {
		m[reuse_key(uint64(i), &buf)] = i
	}
	b.StartTimer()
	var wg sync.WaitGroup
	for j := 0; j < nprocs; j++ {
		wg.Add(1)
		go func() {
			for i := 0; i < b.N; i++ {
				x, _ := m[reuse_key(uint64(i&WRAPPER), &buf)]
				_ = x
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkGoMapReadConcurrentLocked(b *testing.B) {
	nprocs := runtime.GOMAXPROCS(0)
	b.StopTimer()
	runtime.GOMAXPROCS(nprocs)
	m := make(map[Key]Value)
	var rw sync.RWMutex
	var buf [16]byte
	for i := 0; i < NUMKEYS; i++ {
		m[reuse_key(uint64(i), &buf)] = i
	}
	b.StartTimer()
	var wg sync.WaitGroup
	for j := 0; j < nprocs; j++ {
		wg.Add(1)
		go func() {
			for i := 0; i < b.N; i++ {
				rw.RLock()
				x, _ := m[reuse_key(uint64(i&WRAPPER), &buf)]
				_ = x
				rw.RUnlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

// GOTOMIC READS

func benchmarkGotomicReadFixed(b *testing.B) {
	b.StopTimer()
	h := gotomic.NewHash()
	var buf [16]byte
	for i := 0; i < NUMKEYS; i++ {
		h.Put(reuse_key(uint64(i), &buf), i)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		x, _ := h.Get(reuse_key(uint64(i&WRAPPER), &buf))
		_ = x
	}
}

func BenchmarkGotomicReadConcurrent(b *testing.B) {
	nprocs := runtime.GOMAXPROCS(0)
	b.StopTimer()
	runtime.GOMAXPROCS(nprocs)
	h := gotomic.NewHash()
	var buf [16]byte
	for i := 0; i < NUMKEYS; i++ {
		h.Put(reuse_key(uint64(i), &buf), i)
	}
	b.StartTimer()
	var wg sync.WaitGroup
	for j := 0; j < nprocs; j++ {
		wg.Add(1)
		go func() {
			for i := 0; i < b.N; i++ {
				x, _ := h.Get(reuse_key(uint64(i&WRAPPER), &buf))
				_ = x
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

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
	for i := 0; i < NUMKEYS; i++ {
		m[reuse_key(uint64(i), &buf)] = i
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		m[reuse_key(uint64(i&WRAPPER), &buf)] = i
	}
}

func BenchmarkGoMapWriteConcurrent(b *testing.B) {
	nprocs := runtime.GOMAXPROCS(0)
	b.StopTimer()
	runtime.GOMAXPROCS(nprocs)
	m := make(map[Key]Value)
	var rw sync.RWMutex
	var buf [16]byte
	for i := 0; i < NUMKEYS; i++ {
		m[reuse_key(uint64(i), &buf)] = i
	}
	b.StartTimer()
	var wg sync.WaitGroup
	for j := 0; j < nprocs; j++ {
		wg.Add(1)
		go func() {
			for i := 0; i < b.N; i++ {
				rw.Lock()
				m[reuse_key(uint64(i&WRAPPER), &buf)] = i
				rw.Unlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkGotomicWriteEmpty(b *testing.B) {
	h := gotomic.NewHash()
	var buf [16]byte
	for i := 0; i < b.N; i++ {
		h.Put(reuse_key(uint64(i), &buf), i)
	}
}

func benchmarkGotomicWriteFixed(b *testing.B) {
	b.StopTimer()
	h := gotomic.NewHash()
	var buf [16]byte
	for i := 0; i < NUMKEYS; i++ {
		h.Put(reuse_key(uint64(i), &buf), i)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		h.Put(reuse_key(uint64(i&WRAPPER), &buf), i)
	}
}

func BenchmarkGotomicWriteConcurrent(b *testing.B) {
	nprocs := runtime.GOMAXPROCS(0)
	b.StopTimer()
	runtime.GOMAXPROCS(nprocs)
	h := gotomic.NewHash()
	var buf [16]byte
	for i := 0; i < NUMKEYS; i++ {
		h.Put(reuse_key(uint64(i), &buf), i)
	}
	b.StartTimer()
	var wg sync.WaitGroup
	for j := 0; j < nprocs; j++ {
		wg.Add(1)
		go func() {
			for i := 0; i < b.N; i++ {
				h.Put(reuse_key(uint64(i&WRAPPER), &buf), i)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
