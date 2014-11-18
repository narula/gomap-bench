package gomap

import (
	"fmt"
	"gotomic"
	"testing"
	"unsafe"
)

func TestMap(t *testing.T) {
	h := gotomic.NewHash()
	x := 1
	h.Put(gotomic.MakeKey(1), unsafe.Pointer(&x))
	v, ok := h.Get(gotomic.MakeKey(1))
	if !ok {
		t.Fatalf("Could not get key\n")
	}
	va := (*int)(v)
	if *va != 1 {
		t.Fatalf("Wrong va %v\n", *va)
	}
	fmt.Println("PASSED")
}
