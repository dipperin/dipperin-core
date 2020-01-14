package mem_manage

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
		//the basic memory is 4k in buddy
		tree := buildTree(pages * DefaultPageSize / BuddyMinimumSize)
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
