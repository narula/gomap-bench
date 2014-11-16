package cmap

import (
	"hash/crc32"

	"gotomic"
)

type Key [16]byte
type Value interface{}
type KV struct {
	K Key
	V Value
}

func (k Key) HashCode() uint32 {
	return crc32.ChecksumIEEE(k[:])

}

func (k Key) Equals(t gotomic.Thing) bool {
	if sk, ok := t.(Key); ok {
		return sk == k
	}
	return false
}

type Map struct {
	store []KV
	n     int
}

func (m *Map) Get(k Key) (Value, bool) {
	return nil, false
}

func (m *Map) Put(k Key, v Value) error {
	return nil
}

func (m *Map) Delete(k Key) error {
	return nil
}

func reuse_key(x uint64, b *[16]byte) Key {
	var i uint64
	for i = 0; i < 8; i++ {
		(*b)[i] = byte((x >> (i * 8)))
	}
	return Key(*b)
}

func key(x uint64) Key {
	var i uint64
	var b [16]byte
	for i = 0; i < 8; i++ {
		b[i] = byte((x >> (i * 8)))
	}
	return Key(b)
}

func convert_back(k Key) uint64 {
	b := [16]byte(k)
	var x uint64
	var i uint64
	for i = 0; i < 8; i++ {
		v := uint32(b[i])
		x = x + uint64(v<<(i*8))
	}
	return x
}
