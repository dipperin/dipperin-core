package mem_manage

import (
	"fmt"
	"github.com/dipperin/dipperin-core/common/log"
	"go.uber.org/zap"
	"sync"
)

//memory pool abstract interface
type MemPoolInterface interface {
	Get(pages int) []byte
	Put(mem []byte)
}

//tree pool abstract interface
type TreePoolInterface interface {
	GetTree(pages int) tree
	PutTree(tree []int)
}

type VmMemory struct {
	MemPoolInterface
	TreePoolInterface
	*BuddyMemory
	Slab *SlabMemory
	lock sync.RWMutex
}

func NewVmMemory(initialLimit int) *VmMemory {
	memory := &VmMemory{
		MemPoolInterface:  &SimpleMemPool{},
		TreePoolInterface: &SimpleTreePool{},
	}

	memory.BuddyMemory = &BuddyMemory{
		Memory: memory.Get(initialLimit + DynamicMemoryPages),
		Start:  initialLimit * DefaultPageSize,
		Tree:   memory.GetTree(DynamicMemoryPages),
	}
	memory.BuddyMemory.Size = (len(memory.BuddyMemory.Tree) + 1) / 2

	memory.Slab = NewSlabMemory(GrowthFactor, StartChunkSize, DefaultSlabSize, memory.BuddyMemory)
	return memory
}

func (m *VmMemory) Malloc(size int) int {
	m.lock.Lock()
	defer m.lock.Unlock()
	log.DLogger.Debug("[~~VmMemory malloc~~]", zap.Int("size", size))
	//size > 2048 malloc from buddy
	log.DLogger.Debug("the slab max chunkSize is:", zap.Int("chunkSize", m.Slab.MaxChunkSize()))
	if size > m.Slab.MaxChunkSize() {
		return m.BuddyMemory.Malloc(size)
	}

	//size <= 2048 malloc from slab
	offset, err := m.Slab.Malloc(size, 0)
	if err != nil {
		panic(fmt.Errorf("VmMemory malloc error : %v", size))
	}
	log.DLogger.Debug("[~~VmMemory malloc from slab ~~]", zap.Int("addr", offset))
	//clear malloc memory from slab
	clear(offset, offset+size, m.Memory)
	return offset
}

func (m *VmMemory) Free(offset int) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	log.DLogger.Debug("[~~VmMemory free~~]", zap.Int("offset", offset))
	//free from SlabMemory
	err := m.Slab.Free(offset)
	if err == nil {
		return nil
	}

	//free from buddy if offset isn't in slab
	if err != nil && err == ErrOffsetNotInSlab {
		return m.BuddyMemory.Free(offset)
	}

	return err
}
