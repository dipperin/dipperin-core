package mem_manage

import (
	"fmt"
	"github.com/dipperin/dipperin-core/third-party/log"
)

type BuddyMemory struct {
	Memory []byte
	Start  int //start position for malloc
	Size   int //memory size for malloc
	Tree   []int
}

func (m *BuddyMemory) Malloc(size int) int {
	if size <= 0 {
		panic(fmt.Errorf("wrong Size=%d", size))
	} else {
		size = fixSize(size)
	}
	if size > m.Tree[0] {
		panic(fmt.Errorf("malloc Size=%d exceed avalable memory Size", size))
	}

	/*
		find the suitable nodeSize
	*/
	index := 0
	nodeSize := 0
	for nodeSize = m.Size; nodeSize != size; nodeSize /= 2 {
		if m.Tree[left(index)] >= size {
			index = left(index)
		} else {
			index = right(index)
		}
	}
	m.Tree[index] = 0
	//Calculate the address corresponding to the node
	offset := (index+1)*nodeSize - m.Size

	//Upward modify the size of the parent node affected by the size
	for index > 0 {
		index = parent(index)
		m.Tree[index] = max(m.Tree[left(index)], m.Tree[right(index)])
	}
	//Clear the memory data corresponding to the node
	clear(offset+m.Start, offset+m.Start+nodeSize, m.Memory)
	return offset + m.Start
}

func (m *BuddyMemory) Free(offset int) error {
	if offset == 0 {
		log.Debug("free offset = 0...")
		return nil
	}
	offset = offset - m.Start
	if offset < 0 || offset >= m.Size {
		panic(fmt.Errorf("error offset=%d", offset))
	}

	//Lowermost node
	nodeSize := 1
	//Offset corresponds to the node index
	index := offset + m.Size - 1
	//From the last node, go up and find the node with size 0, that is, the size and position of the original allocation block.
	for ; m.Tree[index] != 0; index = parent(index) {
		nodeSize *= 2
		if index == 0 {
			return nil
		}
	}

	//Recovery node
	m.Tree[index] = nodeSize

	//Traverse up the nodes that are affected by the recovery
	var leftNode int
	var rightNode int
	for index = parent(index); index >= 0; index = parent(index) {
		nodeSize *= 2
		leftNode = m.Tree[left(index)]
		rightNode = m.Tree[right(index)]
		if leftNode+rightNode == nodeSize {
			m.Tree[index] = nodeSize
		} else {
			m.Tree[index] = max(leftNode, rightNode)
		}
	}

	return nil
}

func clear(start, end int, mem []byte) {
	for i := start; i < end; i++ {
		mem[i] = 0
	}
}

/**
Calculate the index of the current node to calculate the index of the left leaf node
*/
func left(index int) int {
	return index*2 + 1
}

/**
Calculate the index of the current node and calculate the index of the right leaf node
*/
func right(index int) int {
	return index*2 + 2
}

/**
Calculate the index of the current node to calculate the index of the left leaf node
*/
func parent(index int) int {
	return ((index)+1)/2 - 1
}

func max(a, b int) int {
	if a >= b {
		return a
	} else {
		return b
	}
}

/**
Determine if it is the power of 2
*/
func isPowOf2(n int) bool {
	if n <= 0 {
		return false
	}
	return n&(n-1) == 0
}

/*
Get the minimum power of 2 greater than size
*/
func fixSize(size int) int {

	result := 1
	for result < size {
		result = result << 1
	}
	return result
}
