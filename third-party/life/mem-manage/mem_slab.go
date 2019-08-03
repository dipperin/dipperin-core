package mem_manage

import (
	"github.com/willf/bitset"
)

type NewSlab struct {
	slabAddr int // slab offset addr
	slabSize int
	chunkSize int
	//chunkNumbers int

	//mark the chunk status
	bitSet *bitset.BitSet
}

type NewSlabClass struct {
	slabs []*slab
	chunkSize int     // Each slab is sliced into fixed-sized chunks.
	//chunkFree Loc     // Chunks are tracked in a free-list per slabClass.

	numChunks     int64
	numChunksFree int64
}

type SlabMemory struct {
	memorySource MemManagementInterface
	growthFactor float64
	slabClasses  []NewSlabClass // slabClasses's chunkSizes grow by growthFactor.
	slabSize     int
}


func NewSlabMemory(growthFactor float64, startChunkSize ,slabSize int, memoryOp MemManagementInterface) *SlabMemory{
	slab:= &SlabMemory{
		memorySource:memoryOp,
		growthFactor:growthFactor,
		slabSize:slabSize,
	}
	return slab
}




func (m *SlabMemory) Malloc(size int) int {

	panic("implement me")
}

func (m *SlabMemory) Free(offset int) error {
	panic("implement me")
}
