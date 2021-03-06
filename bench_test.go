package gomap

import (
	"gotomic"
	"runtime"
	"sync"
	"testing"
	"unsafe"

	lgotomic "github.com/zond/gotomic"
)

// GOMAP READS

func BenchmarkGoMapReadOneThreadFixed(b *testing.B) {
	m := make(map[Key1]Value)
	keys := PreallocGotomicKeys(NUMKEYS)
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
	m := make(map[Key1]Value)
	keys := PreallocGotomicKeys(NUMKEYS)
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
	m := make(map[Key1]Value)
	keys := PreallocGotomicKeys(NUMKEYS)
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
	m := make(map[Key1]int)
	keys := PreallocGotomicKeys(NUMKEYS)
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
	m := make(map[Key1]Value)
	keys := PreallocGotomicKeys(NUMKEYS)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m[keys[i]] = i
	}
}

func BenchmarkGoMapWriteOneThreadFixed(b *testing.B) {
	m := make(map[Key1]Value)
	keys := PreallocGotomicKeys(NUMKEYS)
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
	m := make(map[Key1]Value)
	keys := PreallocGotomicKeys(NUMKEYS)
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

// ZGOTOMIC

func BenchmarkZgotomicReadOneThreadFixed(b *testing.B) {
	h := lgotomic.NewHash()
	keys := PreallocGotomicKeys(NUMKEYS)
	for i := 0; i < NUMKEYS; i++ {
		h.Put(keys[i], i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x, _ := h.Get(keys[i&WRAPPER])
		_ = x
	}
}

func BenchmarkZgotomicReadConcurrent(b *testing.B) {
	nprocs := runtime.GOMAXPROCS(0)
	h := lgotomic.NewHash()
	keys := PreallocGotomicKeys(NUMKEYS)
	hcs := make([]uint32, NUMKEYS)
	for i := 0; i < NUMKEYS; i++ {
		h.Put(keys[i], i)
		hcs[i] = keys[i].HashCode()
	}
	var wg sync.WaitGroup
	b.ResetTimer()
	for j := 0; j < nprocs; j++ {
		wg.Add(1)
		go func() {
			for i := 0; i < (b.N / nprocs); i++ {
				it := i & WRAPPER
				x, ok := h.GetHC(hcs[it], keys[it])
				if !ok {
					b.Fatalf("Could not get %v\n", keys[it])
				}
				_ = x
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkZgotomicWriteOneThreadEmpty(b *testing.B) {
	h := lgotomic.NewHash()
	keys := PreallocGotomicKeys(NUMKEYS)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Put(keys[i], i)
	}
}

func BenchmarkZgotomicWriteOneThreadFixed(b *testing.B) {
	h := lgotomic.NewHash()
	keys := PreallocGotomicKeys(NUMKEYS)
	for i := 0; i < NUMKEYS; i++ {
		h.Put(keys[i], i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.Put(keys[i&WRAPPER], i)
	}
}

func BenchmarkZgotomicWriteConcurrent(b *testing.B) {
	nprocs := runtime.GOMAXPROCS(0)
	keys := PreallocGotomicKeys(NUMKEYS)
	h := lgotomic.NewHash()
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

// MY GOTOMIC (GOTOMIC)

func BenchmarkGotomicReadOneThreadFixed(b *testing.B) {
	h := gotomic.NewHash()
	keys := PreallocLocalKeys(NUMKEYS)
	hcs := make([]uint32, NUMKEYS)
	ld := gotomic.InitLocalData()
	for i := 0; i < NUMKEYS; i++ {
		x := i
		h.Put(keys[i], unsafe.Pointer(&x))
		hcs[i] = keys[i].HashCode()
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		it := i & WRAPPER
		x, _ := h.GetHC(hcs[it], keys[it], ld)
		_ = x
	}
}

func BenchmarkGotomicReadConcurrent(b *testing.B) {
	nprocs := runtime.GOMAXPROCS(0)
	h := gotomic.NewHash()
	keys := PreallocLocalKeys(NUMKEYS)
	hcs := make([]uint32, NUMKEYS)
	for i := 0; i < NUMKEYS; i++ {
		x := i
		h.Put(keys[i], unsafe.Pointer(&x))
		hcs[i] = keys[i].HashCode()
	}
	var wg sync.WaitGroup
	b.ResetTimer()
	for j := 0; j < nprocs; j++ {
		wg.Add(1)
		go func() {
			ld := gotomic.InitLocalData()
			for i := 0; i < (b.N / nprocs); i++ {
				it := i & WRAPPER
				x, ok := h.GetHC(hcs[it], keys[it], ld)
				if !ok {
					b.Fatalf("Could not get %v\n", keys[it])
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
	keys := PreallocLocalKeys(NUMKEYS)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x := i
		h.Put(keys[i], unsafe.Pointer(&x))
	}
}

func BenchmarkGotomicWriteOneThreadFixed(b *testing.B) {
	h := gotomic.NewHash()
	keys := PreallocLocalKeys(NUMKEYS)
	for i := 0; i < NUMKEYS; i++ {
		x := i
		h.Put(keys[i], unsafe.Pointer(&x))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x := i
		h.Put(keys[i&WRAPPER], unsafe.Pointer(&x))
	}
}

func BenchmarkGotomicWriteConcurrent(b *testing.B) {
	nprocs := runtime.GOMAXPROCS(0)
	keys := PreallocLocalKeys(NUMKEYS)
	h := gotomic.NewHash()
	for i := 0; i < NUMKEYS; i++ {
		x := i
		h.Put(keys[i], unsafe.Pointer(&x))
	}
	var wg sync.WaitGroup
	b.ResetTimer()
	for j := 0; j < nprocs; j++ {
		wg.Add(1)
		go func() {
			for i := 0; i < (b.N / nprocs); i++ {
				x := i
				h.Put(keys[i], unsafe.Pointer(&x))
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
