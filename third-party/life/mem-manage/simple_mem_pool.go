package mem_manage

import (
	"sync"
)

type SimpleMemPool struct {
	sync.Mutex
	memory *[]byte
}

func (memPool *SimpleMemPool) Get(pages int) []byte {
	memPool.Lock()
	defer memPool.Unlock()
	if pages <= 0 {
		return nil
	}
	if memPool.memory != nil {
		memset(*memPool.memory)
		return *memPool.memory
	}

	//pages = fixSize(pages-DefaultMemoryPages) + DefaultMemoryPages
	memory := make([]byte, pages*DefaultPageSize)
	memset(memory)
	memPool.memory = &memory
	return memory
}

func (memPool *SimpleMemPool) Put(mem []byte) {
	memPool.Lock()
	defer memPool.Unlock()
	memPool.memory = nil
}
