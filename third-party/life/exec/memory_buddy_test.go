package exec

import (
	"github.com/dipperin/dipperin-core/third-party/log"
	"testing"
)

func TestMemory_Malloc(t *testing.T) {

/*	testTree := buildTree(16*DefaultPageSize)
	log.Info("the testTree size is:","size",len(testTree)/1024*4/1024)

	return*/

	tree := buildTree(8)
	mem := Memory{
		Memory: make([]byte,8),
		Start:0,
		Size:8,
		tree:tree,
	}

	log.Info("the build tree is:","tree",mem.tree)

	offset0 := mem.Malloc(2)
	log.Info("the malloc data offset is:","offset",offset0)
	log.Info("the tree after malloc is:","tree",mem.tree)

	offset1 := mem.Malloc(2)
	log.Info("the malloc data offset is:","offset",offset1)
	log.Info("the tree after malloc is:","tree",mem.tree)

	offset2 := mem.Malloc(2)
	log.Info("the malloc data offset is:","offset",offset2)
	log.Info("the tree after malloc is:","tree",mem.tree)

	offset3 := mem.Malloc(1)
	log.Info("the malloc data offset is:","offset",offset3)
	log.Info("the tree after malloc is:","tree",mem.tree)

	offset4 := mem.Malloc(1)
	log.Info("the malloc data offset is:","offset",offset4)
	log.Info("the tree after malloc is:","tree",mem.tree)


	mem.Free(offset4)
	log.Info("the tree after free is:","tree",mem.tree)
	mem.Free(offset3)
	log.Info("the tree after free is:","tree",mem.tree)
	mem.Free(offset2)
	log.Info("the tree after free is:","tree",mem.tree)
	mem.Free(offset1)
	log.Info("the tree after free is:","tree",mem.tree)
	mem.Free(offset0)
	log.Info("the tree after free is:","tree",mem.tree)

}
