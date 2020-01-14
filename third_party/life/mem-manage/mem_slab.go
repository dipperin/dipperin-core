package mem_manage

import (
	"errors"
	"fmt"
	"github.com/dipperin/dipperin-core/common/log"
	"github.com/willf/bitset"
	"go.uber.org/zap"
	"math"
	"sort"
)

type SlabStatus uint8

const (
	SlabFull SlabStatus = iota
	SlabNotEmpty
	SlabEmpty
)

// chunk size: 2~2^10(2byte~2048byte)
const SlabClassNumber = 11
const DefaultSlabSize = uint(4 * 1024)
const StartChunkSize = uint(2)
const GrowthFactor = float64(2)

var ErrOffsetNotInSlab = errors.New("the offset isn't in the slab pool")

//go:generate mockgen -destination=./slabNeedInterface_mock_test.go -package=mem_manage github.com/dipperin/dipperin-core/third-party/life/mem-manage SlabNeedInterface
type SlabNeedInterface interface {
	Malloc(size int) int
	Free(offset int) error
}

type Slab struct {
	SlabAddr  uint // slab offset addr
	SlabSize  uint
	ChunkSize uint
	//chunkNumbers int

	//mark the chunk status
	BitSet *bitset.BitSet
}

type SlabClass struct {
	slabs           []*Slab
	defaultSlabSize uint
	chunkSize       uint

	emptyChunkNumber uint
	memorySource     SlabNeedInterface
}

type SlabMemory struct {
	memorySource   SlabNeedInterface
	growthFactor   float64
	startChunkSize uint
	maxChunkSize   uint

	slabClasses     []*SlabClass // slabClasses's chunkSizes grow by growthFactor.
	defaultSlabSize uint
}

func NewSlab(slabAddr, slabSize, chunkSize uint) (*Slab, error) {
	if slabSize < chunkSize {
		return nil, errors.New("defaultSlabSize < chunkSize")
	}

	slab := &Slab{
		SlabAddr:  slabAddr,
		SlabSize:  slabSize,
		ChunkSize: chunkSize,
	}

	chunkNumber := slabSize / chunkSize
	slab.BitSet = bitset.New(chunkNumber)

	//set bit when initializing
	for i := 0; i < int(chunkNumber); i++ {
		slab.BitSet.Set(uint(i))
	}

	return slab, nil
}

func (s *Slab) emptyChunkIndex(number uint) []uint {
	emptyChunks := make([]uint, number)
	_, emptyChunks = s.BitSet.NextSetMany(0, emptyChunks)
	return emptyChunks
}

func (s *Slab) MallocChunks(number uint) (offset []uint) {
	if s.Status() == SlabFull {
		return []uint{}
	}

	chunks := s.emptyChunkIndex(number)
	chunkOffset := make([]uint, number)
	for i, chunkIndex := range chunks {
		chunkOffset[i] = s.SlabAddr + chunkIndex*s.ChunkSize
		s.BitSet.Clear(chunkIndex)
	}

	return chunkOffset
}

func (s *Slab) Status() SlabStatus {
	setCount := s.BitSet.Count()
	if setCount == s.BitSet.Len() {
		return SlabEmpty
	} else if setCount == 0 {
		return SlabFull
	} else {
		return SlabNotEmpty
	}
}

func (s *Slab) FreeChunks(offset []uint) {
	for _, chunkAddr := range offset {
		index := (chunkAddr - s.SlabAddr) / s.ChunkSize
		if s.BitSet.Test(index) {
			panic(fmt.Errorf("the chunk is already free. index:%v", index))
		}
		s.BitSet.Set(index)
	}
}

func (s *Slab) ChunkNumber() uint {
	return s.SlabSize / s.ChunkSize
}

func (s *Slab) EmptyChunkNumber() uint {
	return s.BitSet.Count()
}

func NewSlabClass(slabSize, chunkSize uint, memorySource SlabNeedInterface) *SlabClass {
	return &SlabClass{
		slabs:            make([]*Slab, 0),
		emptyChunkNumber: 0,
		defaultSlabSize:  slabSize,
		chunkSize:        chunkSize,
		memorySource:     memorySource,
	}
}

func (c *SlabClass) SlabNumber() uint {
	return uint(len(c.slabs))
}

func (c *SlabClass) AddSlab(slabSize uint) error {
	slabAddr := c.memorySource.Malloc(int(slabSize))
	slab, err := NewSlab(uint(slabAddr), slabSize, c.chunkSize)
	if err != nil {
		return err
	}

	c.slabs = append(c.slabs, slab)
	c.emptyChunkNumber += slab.ChunkNumber()
	return nil
}

func (c *SlabClass) MallocChunks(number uint, slabSize uint) (offset []uint, err error) {
	//log.DLogger.Debug("malloc chunks from slab class","number",number,"chunkSize",c.chunkSize,"slabNumber",c.SlabNumber(),"emptyChunkNumber",c.emptyChunkNumber)
	if slabSize == 0 {
		slabSize = c.defaultSlabSize
	}

	//add enough slab number for malloc
	tmpNumber0 := number
	for {
		if tmpNumber0 > c.emptyChunkNumber {
			err := c.AddSlab(slabSize)
			if err != nil {
				return []uint{}, err
			}
		} else {
			break
		}
	}

	tmpNumber := number
	for _, slab := range c.slabs {
		if tmpNumber == 0 {
			break
		}
		if tmpNumber < slab.EmptyChunkNumber() {
			offset = append(offset, slab.MallocChunks(tmpNumber)...)
			break
		} else {
			tmpNumber -= slab.EmptyChunkNumber()
			offset = append(offset, slab.MallocChunks(slab.EmptyChunkNumber())...)
		}
	}

	c.emptyChunkNumber -= number
	//log.DLogger.Debug("the malloc addr is:","addr",offset)
	return offset, nil
}

func (c *SlabClass) findAddrPos(offset []uint, f func(slabIndex int, addr uint) error) (error, []uint) {
	freeLen := 0
	invalidOffset := make([]uint, 0)
	for _, addr := range offset {
		findFlag := false
		for slabIndex, slab := range c.slabs {
			if addr >= slab.SlabAddr && addr < (slab.SlabAddr+slab.ChunkNumber()*slab.ChunkSize) {
				err := f(slabIndex, addr)
				if err != nil {
					return err, []uint{}
				}
				freeLen++
				findFlag = true
				break
			}
		}

		if !findFlag {
			invalidOffset = append(invalidOffset, addr)
		}

	}

	if 0 == len(invalidOffset) {
		//all chunk were free in offset slice
		return nil, []uint{}
	} else {
		//there is a offset not found in the slab
		return ErrOffsetNotInSlab, invalidOffset
	}
}

func (c *SlabClass) FreeChunks(offset []uint) (error, []uint) {
	//log.DLogger.Debug("the free offset is:","offsets",offset)
	//log.DLogger.Debug("free chunks from slab class","chunkSize",c.chunkSize,"slabNumber",c.SlabNumber(),"emptyChunkNumber",c.emptyChunkNumber)
	//free chunk from slab
	f := func(slabIndex int, addr uint) error {
		//l.Info("the slabClass empty chunk number is:","number",c.emptyChunkNumber)
		slab := c.slabs[slabIndex]
		slab.FreeChunks([]uint{addr})
		c.emptyChunkNumber += 1
		return nil
	}

	err, invalidOffset := c.findAddrPos(offset, f)
	if err != nil {
		return err, invalidOffset
	}

	//free and delete empty slab from slab class
	j := 0
	for _, slab := range c.slabs {
		if slab.Status() == SlabEmpty {
			err := c.memorySource.Free(int(slab.SlabAddr))
			if err != nil {
				return err, []uint{}
			}
			c.emptyChunkNumber -= slab.ChunkNumber()
		} else {
			c.slabs[j] = slab
			j++
		}
	}
	c.slabs = c.slabs[:j]
	return nil, []uint{}
}

func NewSlabMemory(growthFactor float64, startChunkSize, slabSize uint, memoryOp SlabNeedInterface) *SlabMemory {
	slabMemory := &SlabMemory{
		memorySource:    memoryOp,
		growthFactor:    growthFactor,
		startChunkSize:  startChunkSize,
		defaultSlabSize: slabSize,
		slabClasses:     make([]*SlabClass, SlabClassNumber),
	}

	chunkSize := slabMemory.startChunkSize
	for i := 0; i < SlabClassNumber; i++ {
		if i == SlabClassNumber-1 {
			slabMemory.maxChunkSize = chunkSize
		}
		slabMemory.slabClasses[i] = NewSlabClass(slabMemory.defaultSlabSize, chunkSize, slabMemory.memorySource)
		chunkSize = uint(math.Ceil(float64(chunkSize) * slabMemory.growthFactor))
	}
	return slabMemory
}

func (m *SlabMemory) MaxChunkSize() int {
	return int(m.maxChunkSize)
}

func (m *SlabMemory) SlabNumber() uint {
	var totalSlabNumber uint
	for _, slabClass := range m.slabClasses {
		totalSlabNumber += slabClass.SlabNumber()
	}

	return totalSlabNumber
}

func (m *SlabMemory) Malloc(size int, slabSize uint) (int, error) {
	log.DLogger.Debug("[**malloc from slab**]", zap.Int("size", size), zap.Int("len(slabClass)", len(m.slabClasses)), zap.Int("maxChunkSize", m.MaxChunkSize()))
	if size <= 0 {
		panic(fmt.Errorf("wrong Size=%d", size))
	}

	i := sort.Search(len(m.slabClasses),
		func(i int) bool { return uint(size) <= m.slabClasses[i].chunkSize })
	if i == SlabClassNumber {
		return 0, errors.New("the size isn't in slab pool")
	}

	slabClass := m.slabClasses[i]
	addr, err := slabClass.MallocChunks(1, slabSize)
	if err != nil {
		return 0, err
	}

	return int(addr[0]), nil
}

func (m *SlabMemory) Free(offset int) (err error) {
	log.DLogger.Debug("[**free from slab**]", zap.Int("offset", offset))
	for _, slabClass := range m.slabClasses {
		err, _ = slabClass.FreeChunks([]uint{uint(offset)})
		if err == nil {
			return nil
		}
	}

	return err
}
