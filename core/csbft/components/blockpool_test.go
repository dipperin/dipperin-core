// Copyright 2019, Keychain Foundation Ltd.
// This file is part of the dipperin-core library.
//
// The dipperin-core library is free software: you can redistribute
// it and/or modify it under the terms of the GNU Lesser General Public License
// as published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// The dipperin-core library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package components

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/bloom"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
	"time"
)

func TestBlockPool_Start(t *testing.T) {
	blockPool := NewBlockPool(1, FakePoolEventNotifier{})
	b1 := &FakeBlock{
		Height:     uint64(1),
		HeaderHash: common.HexToHash("0x232"),
		Headers:    nil,
	}

	blockPool.Start()
	blockPool.AddBlock(b1)
	err := blockPool.Start()
	assert.Equal(t, err.Error(), "block pool already started")
	assert.Equal(t, blockPool.IsEmpty(), false)
}

type BlockpoolconfigImpl struct {
}

func (BlockpoolconfigImpl) CurrentBlock() model.AbstractBlock {
	return &FakeBlock{}
}

func TestBlockPool_Start2(t *testing.T) {
	blockpool := NewBlockPool(1, FakePoolEventNotifier{})
	blockpool.SetNodeConfig(BlockpoolconfigImpl{})

	b1 := &FakeBlock{
		Height:     uint64(1),
		HeaderHash: common.HexToHash("0x232"),
		Headers:    nil,
	}
	b2 := &FakeBlock{
		Height:     uint64(1),
		HeaderHash: common.HexToHash("0x234"),
		Headers:    nil,
	}

	blockpool.Start()
	blockpool.AddBlock(b1)
	blockpool.AddBlock(b2)
	assert.Equal(t, blockpool.IsEmpty(), false)
}

func TestBlockPool_Stop(t *testing.T) {
	blockPool := NewBlockPool(1, FakePoolEventNotifier{})
	b1 := &FakeBlock{
		Height:     uint64(1),
		HeaderHash: common.HexToHash("0x232"),
		Headers:    nil,
	}
	//blockPool.Start()
	err := blockPool.AddBlock(b1)
	assert.Equal(t, err.Error(), "block pool not running")
	assert.Equal(t, blockPool.IsEmpty(), true)

	blockPool.SetPoolEventNotifier(FakePoolEventNotifier{})
	time.Sleep(time.Microsecond * 1)
	blockPool.Stop()
	blockPool.Stop()
	assert.Equal(t, blockPool.IsEmpty(), true)
}

func TestBlockPool_GetBlockByHash(t *testing.T) {
	blockPool := NewBlockPool(1, FakePoolEventNotifier{})
	b1 := &FakeBlock{
		Height:     uint64(1),
		HeaderHash: common.HexToHash("0x232"),
		Headers:    nil,
	}
	blockPool.Start()
	blockPool.AddBlock(b1)
	time.Sleep(time.Microsecond * 1)
	assert.Equal(t, blockPool.IsEmpty(), false)

	block := blockPool.GetBlockByHash(b1.Hash())
	assert.Equal(t, b1, block)
	blockPool.RemoveBlock(b1.Hash())
	time.Sleep(time.Microsecond * 10)
	assert.Equal(t, blockPool.IsEmpty(), true)
}

func TestBlockPool_GetProposalBlock(t *testing.T) {
	blockPool := NewBlockPool(1, FakePoolEventNotifier{})
	b1 := &FakeBlock{
		Height:     uint64(1),
		HeaderHash: common.HexToHash("0x232"),
		Headers:    nil,
	}
	blockPool.Start()
	blockPool.AddBlock(b1)
	time.Sleep(time.Microsecond * 1)
	assert.Equal(t, blockPool.IsEmpty(), false)

	block := blockPool.GetProposalBlock()
	assert.Equal(t, b1, block)
	blockPool.RemoveBlock(block.Hash())
	time.Sleep(time.Microsecond * 1)
	assert.Equal(t, blockPool.IsEmpty(), true)
}

func TestBlockPool_NewHeight(t *testing.T) {
	blockPool := NewBlockPool(1, FakePoolEventNotifier{})
	blockPool.Start()

	blockPool.NewHeight(3)
	time.Sleep(time.Microsecond * 100)
	assert.Equal(t, blockPool.height, uint64(3))

	blockPool.NewHeight(1)
	assert.Equal(t, blockPool.height, uint64(3))
}

func TestBlockPool_AddBlock(t *testing.T) {
	blockPool := NewBlockPool(1, FakePoolEventNotifier{})
	blockPool.Start()

	b1 := &FakeBlock{
		Height:     uint64(3),
		HeaderHash: common.HexToHash("0x232"),
		Headers:    nil,
	}
	blockPool.AddBlock(b1)
	assert.Equal(t, 0, len(blockPool.blocks))

}

type FakePoolEventNotifier struct{}

func (fp FakePoolEventNotifier) BlockPoolNotEmpty() {
	return
}

type FakeBlock struct {
	Height     uint64
	HeaderHash common.Hash
	Headers    model.AbstractHeader
}

func (fb *FakeBlock) SetReceiptHash(receiptHash common.Hash) {
	panic("implement me")
}

func (fb *FakeBlock) GetReceiptHash() common.Hash {
	panic("implement me")
}

func (fb *FakeBlock) IsSpecial() bool {
	panic("implement me")
}

func (fb *FakeBlock) GetRegisterRoot() common.Hash {
	panic("implement me")
}

func (fb *FakeBlock) SetRegisterRoot(root common.Hash) {
	panic("implement me")
}

func (fb *FakeBlock) Body() model.AbstractBody {
	panic("implement me")
}

func (fb *FakeBlock) GetBlockTxsBloom() *iblt.Bloom {
	panic("implement me")
}

func (fb *FakeBlock) Version() uint64 {
	panic("implement me")
}

func (fb *FakeBlock) Number() uint64 {
	return fb.Height
}

func (fb *FakeBlock) Difficulty() common.Difficulty {
	panic("implement me")
}

func (fb *FakeBlock) PreHash() common.Hash {
	panic("implement me")
}

func (fb *FakeBlock) Seed() common.Hash {
	panic("implement me")
}

func (fb *FakeBlock) RefreshHashCache() common.Hash {
	panic("implement me")
}

func (fb *FakeBlock) Hash() common.Hash {
	return common.RlpHashKeccak256(fb)
}

func (fb *FakeBlock) EncodeRlpToBytes() ([]byte, error) {
	panic("implement me")
}

func (fb *FakeBlock) TxIterator(cb func(int, model.AbstractTransaction) error) error {
	panic("implement me")
}

func (fb *FakeBlock) TxRoot() common.Hash {
	panic("implement me")
}

func (fb *FakeBlock) Timestamp() *big.Int {
	panic("implement me")
}

func (fb *FakeBlock) Nonce() common.BlockNonce {
	panic("implement me")
}

func (fb *FakeBlock) GetBody() model.AbstractBody {
	panic("implement me")
}

func (fb *FakeBlock) GetHeader() model.AbstractHeader {
	panic("implement me")
}

func (fb *FakeBlock) StateRoot() common.Hash {
	panic("implement me")
}

func (fb *FakeBlock) SetStateRoot(root common.Hash) {
	panic("implement me")
}

func (fb *FakeBlock) FormatForRpc() interface{} {
	panic("implement me")
}

func (fb *FakeBlock) SetNonce(nonce common.BlockNonce) {
	panic("implement me")
}

func (fb *FakeBlock) CoinBaseAddress() common.Address {
	panic("implement me")
}

func (fb *FakeBlock) GetInterlinks() model.InterLink {
	panic("implement me")
}

func (fb *FakeBlock) SetInterLinkRoot(root common.Hash) {
	panic("implement me")
}

func (fb *FakeBlock) GetInterLinkRoot() (root common.Hash) {
	panic("implement me")
}

func (fb *FakeBlock) SetInterLinks(inter model.InterLink) {
	panic("implement me")
}

func (fb *FakeBlock) Header() model.AbstractHeader {
	return fb.Headers
}

func (fb *FakeBlock) GetEiBloomBlockData(reqEstimator *iblt.HybridEstimator) *model.BloomBlockData {
	panic("implement me")
}

func (fb *FakeBlock) SetVerifications(vs []model.AbstractVerification) {
	panic("implement me")
}

func (fb *FakeBlock) VersIterator(func(int, model.AbstractVerification, model.AbstractBlock) error) error {
	panic("implement me")
}

func (fb *FakeBlock) GetVerifications() []model.AbstractVerification {
	panic("implement me")
}

func (fb *FakeBlock) GetTransactionFees() *big.Int {
	panic("implement me")
}

func (fb *FakeBlock) CoinBase() *big.Int {
	panic("implement me")
}

func (fb *FakeBlock) GetBloomBlockData() *model.BloomBlockData {
	panic("implement me")
}

func (fb *FakeBlock) GetTransactions() []*model.Transaction {
	panic("implement me")
}

func (fb *FakeBlock) GetAbsTransactions() []model.AbstractTransaction {
	panic("implement me")
}

func (fb *FakeBlock) GetBloom() iblt.Bloom {
	panic("implement me")
}

func (fb *FakeBlock) VerificationRoot() common.Hash {
	panic("implement me")
}

func (fb *FakeBlock) HeaderRoot() common.Hash {
	panic("implement me")
}

func (fb *FakeBlock) TxCount() int {
	panic("implement me")
}
