package mem_manage

type TreePoolInterface interface {
	GetTree(pages int) tree
	PutTree(tree []int)
}
