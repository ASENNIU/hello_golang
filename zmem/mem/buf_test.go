package mem_test

import (
	"fmt"
	"testing"
	"zmem/mem"
)

func TestBufPoolSetGet(t *testing.T) {
	pool := mem.MemPool()

	buffer, err := pool.Alloc(1)
	if err != nil {
		fmt.Println("pool alloc error", err)
		return
	}

	buffer.SetBytes([]byte("hello lio"))
	fmt.Printf("get bytes: %+v, toString = %s\n", buffer.GetBytes(), string(buffer.GetBytes()))

	buffer.Pop(4)
	fmt.Printf("get bytes: %+v, toString = %s\n", buffer.GetBytes(), string(buffer.GetBytes()))
}

func TestBufferPoolCopy(t *testing.T) {
	pool := mem.MemPool()

	buffer, err := pool.Alloc(1)
	if err != nil {
		fmt.Println("pool alloc error", err)
		return
	}

	buffer.SetBytes([]byte("hello lio"))
	fmt.Printf("get bytes: %+v, toString = %s\n", buffer.GetBytes(), string(buffer.GetBytes()))

	buffer2, err := pool.Alloc(1)
	if err != nil {
		fmt.Println("pool alloc error", err)
		return
	}

	buffer2.Copy(buffer)
	fmt.Printf("get bytes: %+v, toString = %s\n", buffer2.GetBytes(), string(buffer2.GetBytes()))
}

func TestBufferPoolAdjust(t *testing.T) {
	pool := mem.MemPool()

	buffer, err := pool.Alloc(4096)
	if err != nil {
		fmt.Println("pool alloc error", err)
		return
	}
	buffer.SetBytes([]byte("hello liovale"))
	fmt.Printf("GetBytes: %+v, tostring = %s, Head = %d, Length = %d\r\n", buffer.GetBytes(), string(buffer.GetBytes()), buffer.Head(), buffer.Length())

	buffer.Pop(4)
	fmt.Printf("GetBytes: %+v, tostring = %s, Head = %d, Length = %d\r\n", buffer.GetBytes(), string(buffer.GetBytes()), buffer.Head(), buffer.Length())

	buffer.Adjust()
	fmt.Printf("GetBytes: %+v, tostring = %s, Head = %d, Length = %d\r\n", buffer.GetBytes(), string(buffer.GetBytes()), buffer.Head(), buffer.Length())

}
