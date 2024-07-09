package c_test

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"
	"unsafe"
	"zmem/c"
)

func IsLittleEndian() bool {
	var i int32 = 0x01020304

	// 获取一个字节大小的指针
	u := unsafe.Pointer(&i)
	pb := (*byte)(u)

	// 将指针转换为int类型
	b := *pb

	// 判断大小端
	// 0x04 (03 02 01) 小端
	// 0x01 (02 03 04) 大端
	return (b == 0x04)
}

func IntoBytes(n uint32) []byte {
	x := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})

	var order binary.ByteOrder
	if IsLittleEndian() {
		order = binary.LittleEndian
	} else {
		order = binary.BigEndian
	}
	binary.Write(bytesBuffer, order, x)

	return bytesBuffer.Bytes()
}

func TestMemoryC(t *testing.T) {
	data := c.Malloc(4)
	fmt.Printf("\tdata: %+v, %T\n", data, data)
	myData := (*uint32)(data)
	*myData = 7
	fmt.Printf("\tmyData: %+v, %T\n", *myData, *myData)

	var a uint32 = 111
	c.Memcpy(data, IntoBytes(a), 4)
	fmt.Printf("\tmyData: %+v, %T\n", *myData, *myData)

	c.Free(data)
}
