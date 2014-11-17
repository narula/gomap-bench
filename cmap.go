package gomap

import (
	"hash/crc32"

	"gotomic"

	lgotomic "github.com/zond/gotomic"
)

type Key1 [16]byte
type Key2 [16]byte
type Value interface{}

func (k Key1) HashCode() uint32 {
	return crc32.ChecksumIEEE(k[:])

}

func (k Key1) Equals(t lgotomic.Thing) bool {
	if sk, ok := t.(Key1); ok {
		return sk == k
	}
	return false
}

func MakeKey1(x uint64) Key1 {
	var i uint64
	var b [16]byte
	for i = 0; i < 8; i++ {
		b[i] = byte((x >> (i * 8)))
	}
	return Key1(b)
}

func (k Key2) HashCode() uint32 {
	return crc32.ChecksumIEEE(k[:])

}

func (k Key2) Equals(t gotomic.Thing) bool {
	if sk, ok := t.(Key2); ok {
		return sk == k
	}
	return false
}

func MakeKey2(x uint64) Key2 {
	var i uint64
	var b [16]byte
	for i = 0; i < 8; i++ {
		b[i] = byte((x >> (i * 8)))
	}
	return Key2(b)
}

const (
	NUMKEYS = 2 << 20
	WRAPPER = 2<<20 - 1
)

func PreallocGotomicKeys(n int) []Key1 {
	keys := make([]Key1, n)
	for i := 0; i < n; i++ {
		keys[i] = MakeKey1(uint64(i))
	}
	return keys
}

func PreallocLocalKeys(n int) []Key2 {
	keys := make([]Key2, n)
	for i := 0; i < n; i++ {
		keys[i] = MakeKey2(uint64(i))
	}
	return keys
}
