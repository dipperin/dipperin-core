package components

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/dipperin/dipperin-core/common"
)

func TestNewBlockPool(t *testing.T) {
	rsp := NewBlockPool(0, nil)
	t.Log("www:", rsp)
}

func TestBlockPool_SetNodeConfig(t *testing.T) {
	NewBlockPool(0, nil).SetNodeConfig(nil)
}

func TestBlockPool_SetPoolEventNotifier(t *testing.T) {
	NewBlockPool(0, nil).SetPoolEventNotifier(nil)

}

func TestBlockPool_Start(t *testing.T) {
	NewBlockPool(0, nil).Start()
}

func TestBlockPool_Stop(t *testing.T) {
	NewBlockPool(0, nil).Stop()
}

func TestBlockPool_IsEmpty(t *testing.T) {
	sign := NewBlockPool(0, nil).IsEmpty()
	assert.Equal(t, sign, true)
}

func TestBlockPool_IsRunning(t *testing.T) {
	sign := NewBlockPool(0, nil).IsRunning()
	assert.Equal(t, sign, false)
}

func TestBlockPool_RemoveBlock(t *testing.T) {
	NewBlockPool(0, nil).RemoveBlock(common.Hash{})
}

func TestBlockPool_NewHeight(t *testing.T) {
	NewBlockPool(0, nil).NewHeight(0)
}

func TestBlockPool_AddBlock(t *testing.T) {
	NewBlockPool(0, nil).SetNodeConfig(nil)
}

func TestBlockPool_GetBlockByHash(t *testing.T) {
	//pool := NewBlockPool(0,nil)
	//rsp := pool.GetBlockByHash(common.Hash{})
	//resultC := make(chan model.AbstractBlock)
	//defer close(resultC)
	//assert.Equal(t, rsp, resultC)
}

func TestBlockPool_GetProposalBlock(t *testing.T) {
	//pool := NewBlockPool(0,nil)
	//rsp := pool.GetProposalBlock()
	//resultC := make(chan model.AbstractBlock)
	//defer close(resultC)
	//assert.Equal(t, rsp, resultC)
}
