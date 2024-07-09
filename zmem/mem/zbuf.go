package mem

import "fmt"

// Application interface to memory buffer
type ZBuf struct {
	b *Buf
}

func (zb *ZBuf) Clear() {
	MemPool().Revert(zb.b)
	zb.b = nil
}

func (zb *ZBuf) Pop(len int) {
	if zb.b == nil || len > zb.b.length {
		return
	}

	zb.b.Pop(len)

	if zb.b.Length() == 0 {
		MemPool().Revert(zb.b)
		zb.b = nil
	}
}

func (zb *ZBuf) Data() []byte {
	if zb.b == nil {
		return nil
	}

	return zb.b.GetBytes()
}

func (zb *ZBuf) Adjust() {
	if zb.b == nil {
		zb.b.Adjust()
	}
}

func (zb *ZBuf) Read(src []byte) (err error) {
	if zb.b == nil {
		zb.b, err = MemPool().Alloc(len(src))
		if err != nil {
			fmt.Printf("pool alloc error: %v\r\n", err)
		}
	} else {
		if zb.b.Head() != 0 {
			return nil
		}

		if zb.b.Cap-zb.b.Length() < len(src) {
			newBuf, err := MemPool().Alloc(len(src) + zb.b.Length())
			if err != nil {
				fmt.Printf("pool alloc error: %v\r\n", err)
				return nil
			}

			newBuf.Copy(zb.b)
			MemPool().Revert(zb.b)
			zb.b = newBuf
		}
	}

	zb.b.SetBytes(src)
	return nil
}
