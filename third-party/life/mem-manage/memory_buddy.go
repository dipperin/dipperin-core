package mem_manage

import (
	"fmt"
)

type BuddyMemory struct {
	Memory []byte
	Start  int //start position for malloc
	Size   int //memory size for malloc
	Tree   []int
}

const BuddyMinimumSize = 4*1024

func (m *BuddyMemory) Malloc(size int) int {
	l.Debug("[malloc from buddy memory]","size",size,"start",m.Start,"treeLen",len(m.Tree))
	//minimum memory is 4K in buddy
	if size <= 0 {
		panic(fmt.Errorf("wrong Size=%d", size))
	} else {
		size = (size + BuddyMinimumSize-1)/BuddyMinimumSize
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
	offset = offset * BuddyMinimumSize

	//Upward modify the size of the parent node affected by the size
	for index > 0 {
		index = parent(index)
		m.Tree[index] = max(m.Tree[left(index)], m.Tree[right(index)])
	}
	//Clear the memory data corresponding to the node
	clear(offset+m.Start, offset+m.Start+nodeSize, m.Memory)
	l.Debug("[the buddy memory malloc offset is:]","addr",offset + m.Start)
	return offset + m.Start
}

func (m *BuddyMemory) Free(offset int) error {
	l.Debug("[the buddy memory free offset is:]","addr",offset)
	//todo: 有的wasm会多一个free(0),待在cdt端解决
	if offset == 0 {
		l.Debug("free offset = 0...")
		return nil
	}

	if offset < m.Start{
		panic(fmt.Errorf("free offset is small than momory start,offset:%v start:%v",offset,m.Start))
	}
	offset = offset - m.Start
	if  offset%BuddyMinimumSize !=0{
		panic(fmt.Errorf("offset is invalid error offset=%d", offset))
	}

	offset = offset/BuddyMinimumSize
	if offset < 0 || offset >= m.Size{
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
