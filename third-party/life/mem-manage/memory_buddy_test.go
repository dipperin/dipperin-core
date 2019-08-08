package mem_manage

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMemory_Malloc(t *testing.T) {
	size := 8
	tree := buildTree(size)
	mem := BuddyMemory{
		Memory: make([]byte, size*int(DefaultSlabSize)),
		Start:  1,
		Size:   size,
		Tree:   tree,
	}

	offsets := make([]int,0)
	for i:=0;i<size;i++{
		offset := mem.Malloc(1)
		offsets = append(offsets,offset)
		assert.Equal(t,i*int(DefaultSlabSize)+mem.Start,offset)
	}

	for i:=0;i<size;i++{
		err := mem.Free(offsets[size-1-i])
		assert.NoError(t,err)
	}

	for i:=0;i<size;i++{
		offset := mem.Malloc(1)
		assert.Equal(t,mem.Start,offset)
		err := mem.Free(offset)
		assert.NoError(t,err)
	}

}
