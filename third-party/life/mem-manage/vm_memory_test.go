package mem_manage

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var testInitialLimit = 1
var testVmMemory = NewVmMemory(testInitialLimit)

func TestNewVmMemory(t *testing.T) {
	assert.Equal(t, DynamicMemoryPages*DefaultPageSize/DefaultSlabSize, uint(testVmMemory.Size))
	assert.Equal(t, 2048, testVmMemory.Slab.MaxChunkSize())
}

func TestVmMemory_MallocAndFreeFromBuddy(t *testing.T) {
	//malloc
	offsets := make([]int, 0)
	for i := 0; i < DynamicMemoryPages; i++ {
		offset := testVmMemory.Malloc(int(DefaultSlabSize))
		assert.Equal(t, i*int(DefaultSlabSize)+testVmMemory.BuddyMemory.Start, offset)
		offsets = append(offsets, offset)
	}

	//free
	for i := 0; i < DynamicMemoryPages; i++ {
		err := testVmMemory.Free(offsets[i])
		assert.NoError(t, err)
	}

	//malloc and free
	for i := 0; i < DynamicMemoryPages; i++ {
		offset := testVmMemory.Malloc(int(DefaultSlabSize))
		assert.Equal(t, testVmMemory.Start, offset)
		err := testVmMemory.Free(offset)
		assert.NoError(t, err)
	}
}

func TestVmMemory_MallocAndFreeFromSlab(t *testing.T) {

	//malloc
	offsets := make([]int, 0)
	for i := 0; i < SlabClassNumber; i++ {
		chunkSize := StartChunkSize << uint(i)
		chunkNumber := DefaultSlabSize / chunkSize
		for j := 0; j < int(chunkNumber); j++ {
			offset := testVmMemory.Malloc(int(chunkSize))
			assert.Equal(t, i*BuddyMinimumSize+j*int(chunkSize)+testVmMemory.BuddyMemory.Start, offset)
			offsets = append(offsets, offset)
		}
	}

	//free
	for _, offset := range offsets {
		err := testVmMemory.Free(offset)
		assert.NoError(t, err)
	}

	//malloc and free
	for i := 0; i < SlabClassNumber; i++ {
		chunkSize := StartChunkSize << uint(i)
		offset := testVmMemory.Malloc(int(chunkSize))
		assert.Equal(t, testVmMemory.Start, offset)
		err := testVmMemory.Free(offset)
		assert.NoError(t, err)
	}
}
