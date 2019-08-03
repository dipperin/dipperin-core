package mem_manage

import (
	"github.com/dipperin/dipperin-core/third-party/log"
	"testing"
)

func TestMemory_Malloc(t *testing.T) {
	tree := buildTree(8)
	mem := BuddyMemory{
		Memory: make([]byte, 8),
		Start:  0,
		Size:   8,
		Tree:   tree,
	}

	log.Info("the build Tree is:", "Tree", mem.Tree)

	offset0 := mem.Malloc(2)
	log.Info("the malloc data offset is:", "offset", offset0)
	log.Info("the Tree after malloc is:", "Tree", mem.Tree)

	offset1 := mem.Malloc(2)
	log.Info("the malloc data offset is:", "offset", offset1)
	log.Info("the Tree after malloc is:", "Tree", mem.Tree)

	offset2 := mem.Malloc(2)
	log.Info("the malloc data offset is:", "offset", offset2)
	log.Info("the Tree after malloc is:", "Tree", mem.Tree)

	offset3 := mem.Malloc(1)
	log.Info("the malloc data offset is:", "offset", offset3)
	log.Info("the Tree after malloc is:", "Tree", mem.Tree)

	offset4 := mem.Malloc(1)
	log.Info("the malloc data offset is:", "offset", offset4)
	log.Info("the Tree after malloc is:", "Tree", mem.Tree)

	mem.Free(offset4)
	log.Info("the Tree after free is:", "Tree", mem.Tree)
	mem.Free(offset3)
	log.Info("the Tree after free is:", "Tree", mem.Tree)
	mem.Free(offset2)
	log.Info("the Tree after free is:", "Tree", mem.Tree)
	mem.Free(offset1)
	log.Info("the Tree after free is:", "Tree", mem.Tree)
	mem.Free(offset0)
	log.Info("the Tree after free is:", "Tree", mem.Tree)

}
