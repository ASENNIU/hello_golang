package mem

import (
	"errors"
	"fmt"
	"sync"
)

const (
	m4K   = 4 * 1024
	m16K  = 16 * 1024
	m64K  = 64 * 1024
	m256K = 256 * 1024
	m1M   = 1024 * 1024
	m4M   = 4 * 1024 * 1024
	m8M   = 8 * 1024 * 1024
)

const (
	// extra memory limit: 5GB
	EXTRA_MEM_LIMIT = 5 * 1024 * 1024 // kb
)

type Pool map[int]*Buf

type BufPool struct {
	Pool     Pool
	PoolLock sync.RWMutex

	TotalMem uint64 // total bytes(kb) allocated
}

var bufPoolInstance *BufPool
var once sync.Once

func MemPool() *BufPool {
	once.Do(func() {
		bufPoolInstance = new(BufPool)
		bufPoolInstance.Pool = make(map[int]*Buf)
		bufPoolInstance.TotalMem = 0
		bufPoolInstance.initPool()
	})
	return bufPoolInstance
}

// initial the pool
func (bp *BufPool) initPool() {
	bp.makeBufList(m4K, 5000)  // 4K with 5000, totally 20MB
	bp.makeBufList(m16K, 1000) // 16K with 1000, totally 16MB
	bp.makeBufList(m64K, 500)  // 64K with 500, totally 32MB
	bp.makeBufList(m256K, 200) // 256K with 200, totally 51.2MB
	bp.makeBufList(m1M, 50)    // 1MB with 50, totally 51.2MB
	bp.makeBufList(m4M, 20)    // 4MB with 10, totally 80MB
	bp.makeBufList(m8M, 10)    // 8MB with 10, totally 80MB

}

// make a list of buf
func (bp *BufPool) makeBufList(cap, num int) {
	bp.Pool[cap] = NewBuf(cap)

	var prev *Buf
	prev = bp.Pool[cap]

	for i := 1; i < num; i++ {
		prev.Next = NewBuf(cap)
		prev = prev.Next
	}

	bp.TotalMem += (uint64(cap) / 1024) * uint64(num)
}

func (bp *BufPool) Alloc(N int) (*Buf, error) {
	var index int
	if N <= m4K {
		index = m4K
	} else if N <= m16K {
		index = m16K
	} else if N <= m64K {
		index = m64K
	} else if N <= m256K {
		index = m256K
	} else if N <= m1M {
		index = m1M
	} else if N <= m4M {
		index = m4M
	} else if N <= m8M {
		index = m8M
	} else {
		return nil, errors.New("alloc size too large")
	}

	bp.PoolLock.Lock()
	// if the group in pool is empty, try to allocate a new buf
	if bp.Pool[index] == nil {
		if (bp.TotalMem + uint64(index/1024)) >= uint64(EXTRA_MEM_LIMIT) {
			errStr := fmt.Sprintf("Already use %dKB, no more memory to allocate %dKB!", bp.TotalMem, index/1024)
			return nil, errors.New(errStr)
		}

		newBuf := NewBuf(index)
		bp.TotalMem += uint64(index / 1024)
		bp.PoolLock.Unlock()
		fmt.Printf("Allocate %dKB\r\n", newBuf.Cap/1024)

		return newBuf, nil
	}

	// if the group in pool is not empty, pop a buf from the pool
	targetBuf := bp.Pool[index]
	bp.Pool[index] = targetBuf.Next
	bp.PoolLock.Unlock()
	targetBuf.Next = nil
	fmt.Printf("Allocate %dKB\r\n", targetBuf.Cap/1024)
	return targetBuf, nil
}

func (bp *BufPool) Revert(buf *Buf) error {
	index := buf.Cap
	buf.Clear()

	bp.PoolLock.Lock()
	if _, ok := bp.Pool[index]; !ok {
		ereStr := fmt.Sprintf("Revert %dKB buf, but the Index is not in pool!", buf.Cap/1024)
		return errors.New(ereStr)
	}

	buf.Next = bp.Pool[index]
	bp.Pool[index] = buf
	bp.TotalMem += uint64(index / 1024)
	bp.PoolLock.Unlock()
	fmt.Printf("Revert %dKB mem\r\n", buf.Cap/1024)
	return nil
}
