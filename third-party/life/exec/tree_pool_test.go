package exec

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
	size := 8
	tree := buildTree(size)
	log.Info("the tree len is:", "len(tree)", len(tree))
	log.Info("the tree is:", "tree", tree)
}

func TestNewTreePool(t *testing.T) {
	treePool := NewTreePool(poolSize, cacheSize)
	for k, v := range treePool.pool {
		expect := int(math.Pow(2, float64(k)) * DefaultPageSize)
		if get := v.Get().(tree)[0]; get != expect {
			t.Fatalf("new pool error ,expect tree[0]=%d,get=%d", expect, get)
		}
	}
	for k, v := range treePool.trees {
		expect := int(math.Pow(2, float64(k)) * DefaultPageSize)
		for _, tree := range v {
			if get := tree[0]; get != expect {
				t.Fatalf("new pool error ,expect tree[0]=%d,get=%d", expect, get)
			}
		}
	}
}
