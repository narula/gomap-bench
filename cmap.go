package cmap

type Key [16]byte
type Value interface{}
type KV struct {
	K Key
	V Value
}

type Map struct {
	store []KV
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
