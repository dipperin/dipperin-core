package mem_manage

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"math"
	"sort"
	"testing"
)

func TestNewSlab(t *testing.T) {
	testSlab, err := NewSlab(0, DefaultSlabSize, StartChunkSize)
	assert.NoError(t, err)

	assert.Equal(t, DefaultSlabSize/StartChunkSize, testSlab.ChunkNumber())
	assert.Equal(t, DefaultSlabSize/StartChunkSize, testSlab.EmptyChunkNumber())
	assert.Equal(t, true, testSlab.BitSet.All())
	assert.Equal(t, SlabEmpty, testSlab.Status())
}

func TestSlab_MallocChunks(t *testing.T) {
	testSlab, err := NewSlab(0, DefaultSlabSize, StartChunkSize)
	assert.NoError(t, err)

	for i := 0; i < int(testSlab.ChunkNumber()); i++ {
		offset := testSlab.MallocChunks(1)
		assert.Equal(t, uint(i*2), offset[0])
	}

	assert.Equal(t, SlabFull, testSlab.Status())
}

func TestSlab_FreeChunks(t *testing.T) {
	testSlab, err := NewSlab(0, DefaultSlabSize, StartChunkSize)
	assert.NoError(t, err)

	offsets := testSlab.MallocChunks(50)
	assert.Equal(t, 50, len(offsets))
	assert.Equal(t, SlabNotEmpty, testSlab.Status())

	for i := 0; i < len(offsets); i++ {
		assert.Equal(t, uint(i*2), offsets[i])
	}

	testSlab.FreeChunks([]uint{offsets[0]})
	assert.Equal(t, uint(len(offsets)-1), testSlab.BitSet.Len()-testSlab.BitSet.Count())

	//free all
	testSlab.FreeChunks(offsets[1:])
	assert.Equal(t, SlabEmpty, testSlab.Status())

	offsets = testSlab.MallocChunks(0)
	assert.Equal(t, []uint{}, offsets)

	offsets = testSlab.MallocChunks(1)
	assert.Equal(t, uint(0), offsets[0])

	testSlab.BitSet.Set(offsets[0] / StartChunkSize)
	assert.Panics(t, func() {
		testSlab.FreeChunks(offsets)
	})
}

func newTestSlabClass(t *testing.T) (*SlabClass, *MockSlabNeedInterface) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMemory := NewMockSlabNeedInterface(ctrl)
	slabClass := NewSlabClass(DefaultSlabSize, StartChunkSize, mockMemory)
	assert.Equal(t, uint(0), slabClass.SlabNumber())
	assert.Equal(t, uint(0), slabClass.emptyChunkNumber)

	return slabClass, mockMemory
}

func TestNewSlabClass(t *testing.T) {
	newTestSlabClass(t)
}

func TestSlabClass_MallocChunks(t *testing.T) {
	testSlabClass, mockMemory := newTestSlabClass(t)

	mallocAddr0 := uint(0)
	mallocAddr1 := DefaultSlabSize
	mockMemory.EXPECT().Malloc(gomock.Any()).Return(int(mallocAddr0))
	mockMemory.EXPECT().Malloc(gomock.Any()).Return(int(mallocAddr1))
	mockMemory.EXPECT().Free(gomock.Any()).Return(nil).AnyTimes()

	//malloc 10 chunks from slabClass
	mallocChunkNumber0 := uint(10)
	offsets, err := testSlabClass.MallocChunks(mallocChunkNumber0, uint(0))
	assert.NoError(t, err)
	assert.Equal(t, DefaultSlabSize/StartChunkSize-mallocChunkNumber0, testSlabClass.emptyChunkNumber)
	assert.Equal(t, uint(1), testSlabClass.SlabNumber())

	for i := 0; i < int(mallocChunkNumber0); i++ {
		assert.Equal(t, StartChunkSize*uint(i), offsets[i])
	}
	slab := testSlabClass.slabs[0]
	assert.Equal(t, mallocAddr0, slab.SlabAddr)

	//malloc 2048 chunks from slabClass
	mallocChunkNumber1 := DefaultSlabSize / StartChunkSize
	offsets, err = testSlabClass.MallocChunks(mallocChunkNumber1, 0)
	assert.NoError(t, err)
	assert.Equal(t, uint(2), testSlabClass.SlabNumber())
	assert.Equal(t, mallocChunkNumber1-mallocChunkNumber0, testSlabClass.emptyChunkNumber)
	for i := 0; i < int(mallocChunkNumber1); i++ {
		assert.Equal(t, (uint(i)+mallocChunkNumber0)*StartChunkSize, offsets[i])
	}

	slab0 := testSlabClass.slabs[0]
	slab1 := testSlabClass.slabs[1]
	assert.Equal(t, SlabFull, slab0.Status())
	assert.Equal(t, mallocAddr0, slab0.SlabAddr)
	assert.Equal(t, mallocChunkNumber1-mallocChunkNumber0, slab1.EmptyChunkNumber())
	assert.Equal(t, mallocAddr1, slab1.SlabAddr)
}

func mallocForTest(t *testing.T, mallocChunkNumber uint) (*SlabClass, []uint) {
	testSlabClass, mockMemory := newTestSlabClass(t)

	mallocAddr0 := uint(0)
	mallocNumber := 4
	for i := 0; i < mallocNumber; i++ {
		mockMemory.EXPECT().Malloc(gomock.Any()).Return(int(mallocAddr0 + uint(i)*DefaultSlabSize))
		mockMemory.EXPECT().Free(int(mallocAddr0 + uint(i)*DefaultSlabSize)).Return(nil)
	}

	offsets, err := testSlabClass.MallocChunks(mallocChunkNumber, uint(0))
	assert.NoError(t, err)
	return testSlabClass, offsets
}

func TestSlabClass_findAddrPos(t *testing.T) {
	mallocChunkNumber0 := uint(2058)
	testSlabClass, offsets := mallocForTest(t, mallocChunkNumber0)

	f := func(slabIndex int, addr uint) error {
		result := assert.Equal(t, addr/testSlabClass.defaultSlabSize, uint(slabIndex))
		if !result {
			return errors.New("assert error")
		}
		return nil
	}

	testSlabClass.findAddrPos(offsets, f)
}

func TestSlabClass_FreeChunks(t *testing.T) {
	defaultChunkNumber := DefaultSlabSize / StartChunkSize

	mallocChunkNumber0 := uint(2058)
	testSlabClass, offsets := mallocForTest(t, mallocChunkNumber0)
	initialEmptyNumber := defaultChunkNumber*2 - mallocChunkNumber0
	assert.Equal(t, initialEmptyNumber, testSlabClass.emptyChunkNumber)

	//free chunk one bye one
	for i, addr := range offsets {
		err, _ := testSlabClass.FreeChunks([]uint{addr})
		assert.NoError(t, err)
		if uint(i+1) < defaultChunkNumber {
			//free chunk in the first slab
			assert.Equal(t, initialEmptyNumber+uint(i+1), testSlabClass.emptyChunkNumber)
		} else if uint(i+1) >= defaultChunkNumber && uint(i+1) < mallocChunkNumber0 {
			//free chunk in the second slab
			assert.Equal(t, uint(1), testSlabClass.SlabNumber())
			expectEmptyNumber := initialEmptyNumber + uint(i+1) - (defaultChunkNumber)
			assert.Equal(t, expectEmptyNumber, testSlabClass.emptyChunkNumber)
		} else if uint(i+1) == mallocChunkNumber0 {
			//all chunks were freed
			assert.Equal(t, uint(0), testSlabClass.SlabNumber())
			assert.Equal(t, uint(0), testSlabClass.emptyChunkNumber)
		}
	}

	//free all chunks on one time
	offsets, err := testSlabClass.MallocChunks(mallocChunkNumber0, uint(0))
	assert.NoError(t, err)
	assert.Equal(t, initialEmptyNumber, testSlabClass.emptyChunkNumber)
	err, invalidOffsets := testSlabClass.FreeChunks(offsets)
	assert.Equal(t, 0, len(invalidOffsets))
	assert.NoError(t, err)
	assert.Equal(t, uint(0), testSlabClass.emptyChunkNumber)
	assert.Equal(t, uint(0), testSlabClass.SlabNumber())
}

func TestSlabClass_AddSlab(t *testing.T) {
	testSlabClass, mockMemory := newTestSlabClass(t)
	mallocAddr0 := uint(0)
	mockMemory.EXPECT().Malloc(gomock.Any()).Return(int(mallocAddr0))
	mockMemory.EXPECT().Free(int(mallocAddr0)).Return(nil)

	addSlabSize := DefaultSlabSize * 2
	err := testSlabClass.AddSlab(addSlabSize)
	assert.NoError(t, err)
	assert.Equal(t, addSlabSize/StartChunkSize, testSlabClass.emptyChunkNumber)

	offsets, err := testSlabClass.MallocChunks(addSlabSize/StartChunkSize, 0)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), testSlabClass.SlabNumber())

	err, invalidOffsets := testSlabClass.FreeChunks(offsets)
	assert.Equal(t, 0, len(invalidOffsets))
	assert.NoError(t, err)
	assert.Equal(t, uint(0), testSlabClass.emptyChunkNumber)
}

func newTestSlabMemory(t *testing.T) (*SlabMemory, *MockSlabNeedInterface) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMemory := NewMockSlabNeedInterface(ctrl)
	return NewSlabMemory(GrowthFactor, StartChunkSize, DefaultSlabSize, mockMemory), mockMemory
}

func TestNewSlabMemory(t *testing.T) {
	testSlabMemory, _ := newTestSlabMemory(t)
	assert.Equal(t, StartChunkSize<<(SlabClassNumber-1), uint(testSlabMemory.MaxChunkSize()))
	for i, slabClass := range testSlabMemory.slabClasses {
		assert.Equal(t, StartChunkSize<<uint(i), slabClass.chunkSize)
		assert.Equal(t, uint(0), slabClass.emptyChunkNumber)
		assert.Equal(t, DefaultSlabSize, slabClass.defaultSlabSize)
		assert.Equal(t, uint(0), slabClass.SlabNumber())
	}
}

type TestCmpData struct {
	offsets         []uint
	totalSlabNumber uint
}

// generate a memory distribution test data to compare with the real data
func generateTestCmpData(growthFactor float64, startChunkSize, slabSize uint, slabClassNumber, mallocAddr int) TestCmpData {
	result := TestCmpData{
		offsets: make([]uint, 0),
	}

	chunkSize := make([]uint, slabClassNumber)
	chunkSize[0] = startChunkSize
	for i := 1; i < slabClassNumber; i++ {
		chunkSize[i] = uint(math.Ceil(float64(chunkSize[i-1]) * growthFactor))
	}
	mallocMaxSize := int(chunkSize[slabClassNumber-1])

	totalSlabNum := uint(0)
	addr := uint(mallocAddr)
	var tmpAddAddr uint
	for j := 1; j <= mallocMaxSize; j++ {
		index := sort.Search(len(chunkSize), func(i int) bool {
			return uint(j) <= chunkSize[i]
		})

		if j%int(chunkSize[index]) == 0 {
			var startSize, endSize, needMemSize uint
			if index == 0 {
				startSize = 1
			} else {
				startSize = chunkSize[index-1] + 1
			}
			endSize = chunkSize[index]
			needMemSize = (endSize - startSize + 1) * endSize
			needSlabNumber := (needMemSize + slabSize - 1) / slabSize
			totalSlabNum += needSlabNumber

			//slab中剩余地址要增加
			emptyMemory := needSlabNumber*slabSize - needMemSize + endSize
			tmpAddAddr = emptyMemory
		}

		if j > int(chunkSize[0]) && j%int(chunkSize[index-1]) == 1 {
			//slab的第一个地址需要考虑上一个slab存在空余地址的情况
			addr += tmpAddAddr
		} else if j != 1 {
			//当不用处理空余地址时,直接在上一个地址基础上加上chunkSize即可,
			addr += chunkSize[index]
		}

		result.offsets = append(result.offsets, addr)
	}
	result.totalSlabNumber = totalSlabNum
	return result
}

func TestGenerateTestCmpData(t *testing.T) {
	cmpData := generateTestCmpData(GrowthFactor, StartChunkSize, DefaultSlabSize, 4, 0)

	offset := []uint{
		0, 2, 4096, 4100, 8192, 8200, 8208, 8216, 12288, 12304, 12320, 12336, 12352, 12368, 12384, 12400,
	}
	assert.Equal(t, uint(4), cmpData.totalSlabNumber)
	assert.Equal(t, offset, cmpData.offsets)
}

func TestSlabMemory_MallocAndFree(t *testing.T) {
	testSlabMemory, mockMemory := newTestSlabMemory(t)
	//get test compared data
	cmpData := generateTestCmpData(GrowthFactor, StartChunkSize, DefaultSlabSize, SlabClassNumber, 0)

	mallocAddr0 := uint(0)
	mallocNumber := int(cmpData.totalSlabNumber)
	for i := 0; i < mallocNumber; i++ {
		mockMemory.EXPECT().Malloc(gomock.Any()).Return(int(mallocAddr0 + uint(i)*DefaultSlabSize))
		mockMemory.EXPECT().Free(int(mallocAddr0 + uint(i)*DefaultSlabSize)).Return(nil)
	}

	assert.Panics(t, func() {
		testSlabMemory.Malloc(0, 0)
	})

	//test malloc from 1 to 2048
	for i := 1; i <= testSlabMemory.MaxChunkSize(); i++ {
		offset, err := testSlabMemory.Malloc(i, 0)
		assert.NoError(t, err)
		assert.Equal(t, cmpData.offsets[i-1], uint(offset))
	}
	assert.Equal(t, cmpData.totalSlabNumber, testSlabMemory.SlabNumber())

	// test free
	for _, offset := range cmpData.offsets {
		err := testSlabMemory.Free(int(offset))
		assert.NoError(t, err)
	}
	assert.Equal(t, uint(0), testSlabMemory.SlabNumber())

	mockMemory.EXPECT().Malloc(gomock.Any()).Return(int(mallocAddr0))
	mockMemory.EXPECT().Free(int(mallocAddr0)).Return(nil)
	//test malloc 2048 * 1byte
	for i := 0; i < 2048; i++ {
		offset, err := testSlabMemory.Malloc(2, 0)
		assert.NoError(t, err)
		assert.Equal(t, 2*i, offset)
	}
	assert.Equal(t, uint(1), testSlabMemory.SlabNumber())

	//free
	for i := 0; i < 2048; i++ {
		err := testSlabMemory.Free(i * 2)
		assert.NoError(t, err)
	}
	assert.Equal(t, uint(0), testSlabMemory.SlabNumber())
}
