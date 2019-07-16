package exec

import (
	"github.com/dipperin/dipperin-core/third-party/log"
	"testing"
	"unsafe"
)

func assertNum(t *testing.T, a int, b int, args ...interface{}) {
	if a != b {
		t.Error(args, "result:", a, "except:", b)
	}
}

func TestNewMemPool(t *testing.T) {
	pool := NewMemPool(5, 4)
	log.Info("the pool size is:","poolSize",unsafe.Sizeof(pool))
	m1 := pool.Get(1)
	log.Info("the m1 len is:","len(m1)",len(m1)/1024)
/*	pool.Put(m1)
	assertNum(t, len(pool.memBlock[0].FreeMem), 4, "alloc error")

	mem := pool.Get(2 + DefaultMemoryPages)
	assertNum(t, len(pool.memBlock[1].FreeMem), 3, "alloc error")
	assertNum(t, len(pool.largeMem), 0, "alloc error")

	pool.Put(mem)
	assertNum(t, len(pool.memBlock[1].FreeMem), 4, "alloc error")

	pool.Get(2 + DefaultMemoryPages)
	pool.Get(2 + DefaultMemoryPages)
	pool.Get(2 + DefaultMemoryPages)
	pool.Get(2 + DefaultMemoryPages)
	pool.Get(2 + DefaultMemoryPages)

	assertNum(t, len(pool.memBlock[1].FreeMem), 0, "alloc error")

	pool.Get(12)
	assertNum(t, len(pool.largeMem), 0, "alloc error")

	lm := pool.Get(45)

	assertNum(t, len(pool.largeMem), 1, "alloc error")
	pool.Put(lm)*/
}
