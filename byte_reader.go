package tsgo

import (
	"sync"
)

type ByteReader struct {
	sync.Mutex

	b []byte
	i uint
}

// Index is used to check the buffer pointers current position.
func (r *ByteReader) Position() uint {
	r.Lock()
	defer r.Unlock()

	return r.i
}

// Jump moves the buffer pointer to a given position.
func (r *ByteReader) Jump(i uint) {
	r.Lock()
	defer r.Unlock()

	r.i = i
}

// Inc is used to jump the buffer forward n-bytes forward.
func (r *ByteReader) Inc(i uint) {
	r.Lock()
	defer r.Unlock()

	r.i += i
}

// Dec is used to jump the buffer forward n-bytes backward.
func (r *ByteReader) Dec(i uint) {
	r.Lock()
	defer r.Unlock()

	r.i -= i
}

// ReadByte is used to read the current byte at the buffer pointer and move
// forward by one byte.
func (r *ByteReader) ReadByte() byte {
	r.Lock()
	defer r.Unlock()

	b := r.b[r.i]
	r.i += 1
	return b
}

// ReadBytes is used to read the next N bytes, starting at the buffer pointer,
// and move forward by that many bytes.
func (r *ByteReader) ReadBytes(i uint) []byte {
	r.Lock()
	defer r.Unlock()

	b := r.b[r.i : r.i+i]
	r.i += i
	return b
}
