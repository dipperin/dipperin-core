package mem_manage

import (
	"github.com/dipperin/dipperin-core/third-party/log"
	"math"
	"testing"
)

var (
	poolSize  = 5
	cacheSize = 4
)

func TestBuildTree(t *testing.T) {
	//buildTree
	size := 16
	tree := buildTree(size)
	log.Info("the Tree len is:", "len(Tree)", len(tree))
	log.Info("the Tree is:", "Tree", tree)
}

func TestNewTreePool(t *testing.T) {
	treePool := NewTreePool(poolSize, cacheSize)
	for k, v := range treePool.pool {
		expect := int(math.Pow(2, float64(k)) * DefaultPageSize)
		if get := v.Get().(tree)[0]; get != expect {
			t.Fatalf("new pool error ,expect Tree[0]=%d,get=%d", expect, get)
		}
	}
	for k, v := range treePool.trees {
		expect := int(math.Pow(2, float64(k)) * DefaultPageSize)
		for _, tree := range v {
			if get := tree[0]; get != expect {
				t.Fatalf("new pool error ,expect Tree[0]=%d,get=%d", expect, get)
			}
		}
	}
}
