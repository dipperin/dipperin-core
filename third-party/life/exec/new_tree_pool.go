package exec

import (
	"sync"
)

type SimpleTreePool struct {
	sync.Mutex
	tree *tree
}

func (treePool *SimpleTreePool) GetTree(pages int) tree {
	treePool.Lock()
	defer treePool.Unlock()

	if treePool.tree == nil {
		pages = fixSize(pages)
		tree := buildTree(pages * DefaultPageSize)
		treePool.tree = &tree
		return *treePool.tree
	}

	return *treePool.tree
}

func (treePool *SimpleTreePool) PutTree(tree []int) {
	treePool.Lock()
	defer treePool.Unlock()
	treePool.tree = nil
}
