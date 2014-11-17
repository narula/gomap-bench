package main

import (
	"ddtxn/prof"
	"flag"
	"fmt"
	"gomap"
	"gotomic"
	"log"
	"runtime"
	"sync"
	"time"

	lgotomic "github.com/zond/gotomic"
)

var nprocs = flag.Int("nprocs", 2, "GOMAXPROCS default 2")
var clientGoRoutines = flag.Int("ngo", 0, "Number of goroutines/workers generating client requests.")
var mapType = flag.Int("map", 1, "Map type; 0=go's map, 1=my gotomic, 2=github gotomic")

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(*nprocs)

	if *clientGoRoutines == 0 {
		*clientGoRoutines = *nprocs
	}

	var wg sync.WaitGroup
	nitr := 2000000
	var p *prof.Profile
	var start time.Time

	switch *mapType {
	case 0:
		h := make(map[gomap.Key1]gomap.Value)
		keys := gomap.PreallocGotomicKeys(gomap.NUMKEYS)
		for i := 0; i < gomap.NUMKEYS; i++ {
			h[keys[i]] = i
		}
		p = prof.StartProfile()
		start = time.Now()

		for i := 0; i < *clientGoRoutines; i++ {
			wg.Add(1)
			go func(n int) {
				for j := 0; j < nitr; j++ {
					x, ok := h[keys[(j+n)&gomap.WRAPPER]]
					if !ok {
						log.Fatalf("Could not get %v\n", keys[(j+n)&gomap.WRAPPER])
					}
					_ = x
				}
				wg.Done()
			}(i)
		}
	case 1:
		h := gotomic.NewHash()
		keys := gomap.PreallocLocalKeys(gomap.NUMKEYS)
		for i := 0; i < gomap.NUMKEYS; i++ {
			h.Put(keys[i], i)
		}
		p = prof.StartProfile()
		start = time.Now()

		for i := 0; i < *clientGoRoutines; i++ {
			wg.Add(1)
			go func(n int) {
				for j := 0; j < nitr; j++ {
					x, ok := h.Get(keys[(j+n)&gomap.WRAPPER])
					if !ok {
						log.Fatalf("Could not get %v\n", keys[(j+n)&gomap.WRAPPER])
					}
					_ = x
				}
				wg.Done()
			}(i)
		}
	case 2:
		h := lgotomic.NewHash()
		keys := gomap.PreallocGotomicKeys(gomap.NUMKEYS)
		for i := 0; i < gomap.NUMKEYS; i++ {
			h.Put(keys[i], i)
		}
		p = prof.StartProfile()
		start = time.Now()

		for i := 0; i < *clientGoRoutines; i++ {
			wg.Add(1)
			go func(n int) {
				for j := 0; j < nitr; j++ {
					x, ok := h.Get(keys[(j+n)&gomap.WRAPPER])
					if !ok {
						log.Fatalf("Could not get %v\n", keys[(j+n)&gomap.WRAPPER])
					}
					_ = x
				}
				wg.Done()
			}(i)
		}
	}

	wg.Wait()
	end := time.Since(start)
	p.Stop()
	fmt.Printf("ns/txn: %v\n", end.Nanoseconds()/int64(nitr*(*clientGoRoutines)))
}
