package cmap

import (
	"runtime"
	"sync"
	"testing"

	"github.com/zond/gotomic"
)

const (
	NUMKEYS = 2 << 20
	WRAPPER = 2<<20 - 1
)

func preallocKeys(n int) []Key {
	keys := make([]Key, n)
	for i := 0; i < n; i++ {
		keys[i] = key(uint64(i))
	}
	return keys
}

// GOMAP READS

func BenchmarkGoMapReadOneThreadFixed(b *testing.B) {
	m := make(map[Key]Value)
	keys := preallocKeys(NUMKEYS)
	for i := 0; i < NUMKEYS; i++ {
		m[keys[i]] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x, _ := m[keys[i&WRAPPER]]
		_ = x
	}
}

func BenchmarkGoMapReadConcurrentNoLock(b *testing.B) {
	nprocs := runtime.GOMAXPROCS(0)
	m := make(map[Key]Value)
	keys := preallocKeys(NUMKEYS)
	for i := 0; i < NUMKEYS; i++ {
		m[keys[i]] = i
	}
	var wg sync.WaitGroup
	b.ResetTimer()
	for j := 0; j < nprocs; j++ {
		wg.Add(1)
		go func() {
			for i := 0; i < (b.N / nprocs); i++ {
				y := i & WRAPPER
				if y >= NUMKEYS {
					wg.Done()
					b.Fatalf("Somehow generated a %b; NUMKEYS: %b, WRAPPER: %v IDEAL: %v; i: %v\n", y, NUMKEYS, WRAPPER, ((2 << 20) - 1), i)
				}
				x, ok := m[keys[y]]
				if !ok {
					wg.Done()
					b.Fatalf("Could not get %v\n", keys[i&WRAPPER])
				}
				_ = x
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkGoMapReadConcurrentLocked(b *testing.B) {
	nprocs := runtime.GOMAXPROCS(0)
	m := make(map[Key]Value)
	keys := preallocKeys(NUMKEYS)
	var rw sync.RWMutex
	for i := 0; i < NUMKEYS; i++ {
		m[keys[i]] = i
	}
	var wg sync.WaitGroup
	b.ResetTimer()
	for j := 0; j < nprocs; j++ {
		wg.Add(1)
		go func() {
			for i := 0; i < (b.N / nprocs); i++ {
				rw.RLock()
				x, ok := m[keys[i&WRAPPER]]
				_ = x
				rw.RUnlock()
				if !ok {
					b.Fatalf("Could not get %v\n", keys[i&WRAPPER])
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkGoMapReadTypedLocked(b *testing.B) {
	nprocs := runtime.GOMAXPROCS(0)
	m := make(map[Key]int)
	keys := preallocKeys(NUMKEYS)
	var rw sync.RWMutex
	for i := 0; i < NUMKEYS; i++ {
		m[keys[i]] = i
	}
	var wg sync.WaitGroup
	b.ResetTimer()
	for j := 0; j < nprocs; j++ {
		wg.Add(1)
		go func() {
			for i := 0; i < (b.N / nprocs); i++ {
				rw.RLock()
				x, ok := m[keys[i&WRAPPER]]
				_ = x
				rw.RUnlock()
				if !ok {
					b.Fatalf("Could not get %v\n", keys[i&WRAPPER])
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkGoMapWriteOneThreadEmpty(b *testing.B) {
	m := make(map[Key]Value)
	keys := preallocKeys(NUMKEYS)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m[keys[i]] = i
	}
}

func BenchmarkGoMapWriteOneThreadFixed(b *testing.B) {
	m := make(map[Key]Value)
	keys := preallocKeys(NUMKEYS)
	for i := 0; i < NUMKEYS; i++ {
		m[keys[i]] = i
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m[keys[i&WRAPPER]] = i
	}
}

func BenchmarkGoMapWriteConcurrent(b *testing.B) {
	nprocs := runtime.GOMAXPROCS(0)
	m := make(map[Key]Value)
	keys := preallocKeys(NUMKEYS)
	var rw sync.RWMutex
	for i := 0; i < NUMKEYS; i++ {
		m[keys[i]] = i
	}
	var wg sync.WaitGroup
	b.ResetTimer()
	for j := 0; j < nprocs; j++ {
		wg.Add(1)
		go func() {
			for i := 0; i < (b.N / nprocs); i++ {
				rw.Lock()
				m[keys[i&WRAPPER]] = i
				rw.Unlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

// GOTOMIC

func BenchmarkGotomicReadOneThreadFixed(b *testing.B) {
	h := gotomic.NewHash()
	keys := preallocKeys(NUMKEYS)
	for i := 0; i < NUMKEYS; i++ {
		h.Put(keys[i], i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x, _ := h.Get(keys[i&WRAPPER])
		_ = x
	}
}

func BenchmarkGotomicReadConcurrent(b *testing.B) {
	nprocs := runtime.GOMAXPROCS(0)
	h := gotomic.NewHash()
	keys := preallocKeys(NUMKEYS)
	for i := 0; i < NUMKEYS; i++ {
		h.Put(keys[i], i)
	}
	var wg sync.WaitGroup
	b.ResetTimer()
	for j := 0; j < nprocs; j++ {
		wg.Add(1)
		go func() {
			for i := 0; i < (b.N / nprocs); i++ {
				x, ok := h.Get(keys[i&WRAPPER])
				if !ok {
					b.Fatalf("Could not get %v\n", keys[i&WRAPPER])
				}
				_ = x
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkGotomicWriteOneThreadEmpty(b *testing.B) {
	h := gotomic.NewHash()
	keys := preallocKeys(NUMKEYS)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Put(keys[i], i)
	}
}

func BenchmarkGotomicWriteOneThreadFixed(b *testing.B) {
	h := gotomic.NewHash()
	keys := preallocKeys(NUMKEYS)
	for i := 0; i < NUMKEYS; i++ {
		h.Put(keys[i], i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Put(keys[i&WRAPPER], i)
	}
}

func BenchmarkGotomicWriteConcurrent(b *testing.B) {
	nprocs := runtime.GOMAXPROCS(0)
	keys := preallocKeys(NUMKEYS)
	h := gotomic.NewHash()
	for i := 0; i < NUMKEYS; i++ {
		h.Put(keys[i], i)
	}
	var wg sync.WaitGroup
	b.ResetTimer()
	for j := 0; j < nprocs; j++ {
		wg.Add(1)
		go func() {
			for i := 0; i < (b.N / nprocs); i++ {
				h.Put(keys[i&WRAPPER], i)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
