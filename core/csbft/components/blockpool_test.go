package components

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/dipperin/dipperin-core/common"
	"time"
	"github.com/dipperin/dipperin-core/core/model"
	"errors"
)

func TestNewBlockPool(t *testing.T) {
	rsp := NewBlockPool(0, nil)
	t.Log("www:", rsp)
	assert.NotEmpty(t, rsp)
}

func TestBlockPool_SetNodeConfig(t *testing.T) {
	rsp := NewBlockPool(0, nil)
	assert.NotEmpty(t, rsp)
	rsp.SetNodeConfig(nil)
	assert.Equal(t, nil, rsp.Blockpoolconfig)
}

func TestBlockPool_SetPoolEventNotifier(t *testing.T) {
	rsp := NewBlockPool(0, nil)
	assert.NotEmpty(t, rsp)
	rsp.SetPoolEventNotifier(nil)
	assert.Equal(t, nil, rsp.poolEventNotifier)
}

func TestBlockPool_Start(t *testing.T) {
	rsp := NewBlockPool(0, nil)
	assert.NotEmpty(t, rsp)
	assert.NoError(t, rsp.Start())
}

func TestBlockPool_Stop(t *testing.T) {
	rsp := NewBlockPool(0, nil)
	assert.NotEmpty(t, rsp)
	if assert.NoError(t, rsp.Start()) {
		time.Sleep(500 * time.Millisecond)
		rsp.Stop()
	}
}

func TestBlockPool_IsEmpty(t *testing.T) {
	sign := NewBlockPool(0, nil).IsEmpty()
	assert.Equal(t, sign, true)
}

func TestBlockPool_loop(t *testing.T) {
	rsp := NewBlockPool(0, nil)
	assert.NotEmpty(t, rsp)
	go rsp.loop()
	go func(rsp *BlockPool) {
		rsp.stopChan = make(chan struct{})
		rsp.stopChan <- struct{}{}
	}(rsp)
}

func TestBlockPool_doRemoveBlock(t *testing.T) {
	rsp := NewBlockPool(0, nil)
	assert.NotEmpty(t, rsp)
	var hashTmp = `0xd50866a60b4f7e4123400e0563efb987dc800d1a72af5cc1ae9ee68760bb18889`
	rsp.doRemoveBlock(common.HexToHash(hashTmp))
	assert.Equal(t, 0, len(rsp.blocks))
}

func TestBlockPool_IsRunning(t *testing.T) {
	sign := NewBlockPool(0, nil).IsRunning()
	assert.Equal(t, sign, false)
}

func TestBlockPool_RemoveBlock(t *testing.T) {
	rsp := NewBlockPool(0, nil)
	assert.NotEmpty(t, rsp)
	var hashTmp = `0xd50866a60b4f7e4123400e0563efb987dc800d1a72af5cc1ae9ee68760bb18889`
	rsp.RemoveBlock(common.HexToHash(hashTmp))
	assert.Equal(t, 0, len(rsp.rmBlockChan))
}

func TestBlockPool_NewHeight(t *testing.T) {
	rsp := NewBlockPool(0, nil)
	assert.NotEmpty(t, rsp)
	rsp.NewHeight(0)
	assert.Equal(t, 0, len(rsp.newHeightChan))
}

func TestBlockPool_doNewHeight(t *testing.T) {
	rsp := NewBlockPool(0, nil)
	assert.NotEmpty(t, rsp)
	rsp.doNewHeight(0)
	assert.Equal(t, 0, len(rsp.blocks))
	assert.Equal(t, uint64(0), rsp.height)
}

func TestBlockPool_AddBlock(t *testing.T) {
	rsp := NewBlockPool(0, nil)
	assert.NotEmpty(t, rsp)
	assert.Error(t, errors.New("block pool not running"), rsp.AddBlock(nil))
}

func TestBlockPool_GetBlockByHash(t *testing.T) {
	pool := NewBlockPool(0, nil)
	assert.NotEmpty(t, pool)
	var rsp model.AbstractBlock
	go func(block *model.AbstractBlock) {
		var hashTmp = `0xd50866a60b4f7e4123400e0563efb987dc800d1a72af5cc1ae9ee68760bb18889`
		rspT := pool.GetBlockByHash(common.HexToHash(hashTmp))
		block = &rspT
	}(&rsp)
	time.Sleep(500 * time.Millisecond)
	assert.Empty(t, rsp)
}

func TestBlockPool_doAddBlock(t *testing.T) {
	rsp := NewBlockPool(0, nil)
	assert.NotEmpty(t, rsp)
	resultChan := make(chan error)
	var block model.AbstractBlock
	b := newBlockWithResultErr{block: block, resultChan: resultChan}
	assert.Panics(t, func() {
		rsp.doAddBlock(b)
	})
}

func TestBlockPool_GetProposalBlock(t *testing.T) {
	pool := NewBlockPool(0, nil)
	assert.NotEmpty(t, pool)
	var rsp model.AbstractBlock
	go func(block *model.AbstractBlock) {
		rspT := pool.GetProposalBlock()
		block = &rspT
	}(&rsp)
	time.Sleep(500 * time.Millisecond)
	assert.Empty(t, rsp)
}

func TestBlockPool_getBlock(t *testing.T) {
	rsp := NewBlockPool(0, nil)
	assert.NotEmpty(t, rsp)
	resultC := make(chan model.AbstractBlock)
	var hashTmp = `0xd50866a60b4f7e4123400e0563efb987dc800d1a72af5cc1ae9ee68760bb18889`
	h := common.HexToHash(hashTmp)
	assert.NotPanics(t, func() {
		rsp.getBlock(
			&blockPoolGetter{
				blockHash:  h,
				resultChan: resultC,
			},
		)
	})
}

func TestBlockPool_doGetBlock(t *testing.T) {
	rsp := NewBlockPool(0, nil)
	assert.NotEmpty(t, rsp)
	resultC := make(chan model.AbstractBlock)
	var hashTmp = `0xd50866a60b4f7e4123400e0563efb987dc800d1a72af5cc1ae9ee68760bb18889`
	h := common.HexToHash(hashTmp)
	var getter = blockPoolGetter{
		blockHash:  h,
		resultChan: resultC,
	}
	go func(g *blockPoolGetter) {
		rsp.doGetBlock(
			&blockPoolGetter{
				blockHash:  h,
				resultChan: resultC,
			},
		)
	}(&getter)
	assert.Equal(t, 0, len(getter.resultChan))
}
