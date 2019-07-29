package exec

type TreePoolInterface interface {
	GetTree(pages int) tree
	PutTree(tree []int)
}
