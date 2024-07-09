package mem

/*
#include<stdlib.h>
#include<string.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
	"zmem/c"
)

type Buf struct {
	Next   *Buf           // next buffer
	Cap    int            // capacity of current buffer
	length int            // length of bytes in current buffer
	head   int            // index of first byte not used
	data   unsafe.Pointer // address of data
}

func NewBuf(size int) *Buf {
	return &Buf{
		Cap:    size,
		length: 0,
		head:   0,
		data:   c.Malloc(size),
	}
}

// Fill the buffer with bytes
func (b *Buf) SetBytes(src []byte) {
	c.Memcpy(unsafe.Pointer(uintptr(b.data)+uintptr(b.head)), src, len(src))
	b.length += len(src)
}

// Get the bytes in the buffer
func (b *Buf) GetBytes() []byte {
	data := C.GoBytes(unsafe.Pointer(uintptr(b.data)+uintptr(b.head)), C.int(b.length))
	return data
}

// Copy the bytes from another buffer
func (b *Buf) Copy(other *Buf) {
	c.Memcpy(b.data, other.GetBytes(), other.length)
	b.length = other.length
	b.head = 0
}

// Pop the bytes from the buffer
func (b *Buf) Pop(len int) {
	if b.data == nil {
		fmt.Printf("buf is nil\r\n")
		return
	}

	if b.length < len {
		fmt.Printf("buf length: %d < pop len: %d\r\n", b.length, len)
	}

	b.head += len
	b.length -= len
}

// Adjust the buffer, move the data to the start of the buffer
func (b *Buf) Adjust() {
	if b.head != 0 {
		if b.length != 0 {
			c.Memmove(b.data, unsafe.Pointer(uintptr(b.data)+uintptr(b.head)), b.length)
		}
		b.head = 0
	}
}

// Clear the buffer, reset the length and head
func (b *Buf) Clear() {
	b.length = 0
	b.head = 0
}

func (b *Buf) Head() int {
	return b.head
}

func (b *Buf) Length() int {
	return b.length
}
