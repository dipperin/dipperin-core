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


package minemaster

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/bloom"
	"github.com/dipperin/dipperin-core/core/model"
	"testing"
	"math/big"
	"github.com/stretchr/testify/assert"
	"sync/atomic"
)

func TestDefaultPerformance_GetPerformance(t *testing.T) {
	/*m := newDefaultWorkManager(&fakeContext{})
	m.submitBlock(common.Address{1,}, &model.Block{})

	assert.EqualValues(t, 1, m.performance[common.Address{1,}].getPerformance())

	m.submitBlock(common.Address{1,}, &model.Block{})

	assert.EqualValues(t, 2, m.performance[common.Address{1,}].getPerformance())*/
}

type fakeCalculableBlock struct{}

func (fakeCalculableBlock) SetReceiptHash(receiptHash common.Hash) {
	panic("implement me")
}

func (fakeCalculableBlock) GetReceiptHash() common.Hash {
	panic("implement me")
}

func (fakeCalculableBlock) IsSpecial() bool {
	panic("implement me")
}

func (fakeCalculableBlock) GetRegisterRoot() common.Hash {
	panic("implement me")
}

func (fakeCalculableBlock) SetRegisterRoot(root common.Hash) {
	panic("implement me")
}

func (fakeCalculableBlock) Body() model.AbstractBody {
	panic("implement me")
}

func (fakeCalculableBlock) RefreshHashCache() common.Hash {
	panic("implement me")
}

func (fakeCalculableBlock) Header() model.AbstractHeader {
	panic("implement me")
}

func (fakeCalculableBlock) GetBlockTxsBloom() *iblt.Bloom {
	panic("implement me")
}

func (fakeCalculableBlock) SetVerifications(vs []model.AbstractVerification) {
	panic("implement me")
}

func (fakeCalculableBlock) VersIterator(func(int, model.AbstractVerification, model.AbstractBlock) error) (error) {
	panic("implement me")
}

func (fakeCalculableBlock) GetVerifications() ([]model.AbstractVerification) {
	panic("implement me")
}

func (fakeCalculableBlock) GetInterlinks() model.InterLink {
	panic("implement me")
}

func (fakeCalculableBlock) SetInterLinkRoot(root common.Hash) {
	panic("implement me")
}

func (fakeCalculableBlock) GetInterLinkRoot() (root common.Hash) {
	panic("implement me")
}

func (fakeCalculableBlock) SetInterLinks(inter model.InterLink) {
	panic("implement me")
}

func (fakeCalculableBlock) GetEiBloomBlockData(reqEstimator *iblt.HybridEstimator) *model.BloomBlockData {
	panic("implement me")
}

func (fakeCalculableBlock) GetAbsTransactions() []model.AbstractTransaction {
	panic("implement me")
}

func (fakeCalculableBlock) TxCount() int { panic("implement me") }

func (fakeCalculableBlock) GetTransactions() []*model.Transaction {
	panic("implement me")
}

func (fakeCalculableBlock) GetBloom() iblt.Bloom {
	panic("implement me")
}

func (fakeCalculableBlock) VerificationRoot() common.Hash {
	panic("implement me")
}

func (fakeCalculableBlock) HeaderRoot() common.Hash {

	panic("implement me")
}

func (fakeCalculableBlock) GetBloomBlockData() *model.BloomBlockData {
	panic("implement me")
}

func (fakeCalculableBlock) GetBodyDataByBloom() *model.BloomBlockData {
	panic("implement me")
}

func (fakeCalculableBlock) Version() uint64 {
	panic("implement me")
}

func (fakeCalculableBlock) Number() uint64 {
	panic("implement me")
}

func (fakeCalculableBlock) Difficulty() common.Difficulty {
	panic("implement me")
}

func (fakeCalculableBlock) PreHash() common.Hash {
	panic("implement me")
}

func (fakeCalculableBlock) Hash() common.Hash {
	panic("implement me")
}

func (fakeCalculableBlock) EncodeRlpToBytes() ([]byte, error) {
	panic("implement me")
}

func (fakeCalculableBlock) TxIterator(func(index int, tx model.AbstractTransaction) ( error)) (error){
	panic("implement me")
}

func (fakeCalculableBlock) TxRoot() common.Hash {
	panic("implement me")
}

func (fakeCalculableBlock) Timestamp() *big.Int {
	panic("implement me")
}

func (fakeCalculableBlock) Nonce() common.BlockNonce {
	panic("implement me")
}

func (fakeCalculableBlock) GetBody() model.AbstractBody {
	panic("implement me")
}

func (fakeCalculableBlock) GetHeader() model.AbstractHeader {
	panic("implement me")
}

func (fakeCalculableBlock) StateRoot() common.Hash {
	panic("implement me")
}

func (fakeCalculableBlock) SetStateRoot(root common.Hash) {
	panic("implement me")
}

func (fakeCalculableBlock) FormatForRpc() interface{} {
	panic("implement me")
}

func (fakeCalculableBlock) SetNonce(nonce common.BlockNonce) {
	panic("implement me")
}

func (fakeCalculableBlock) CoinBaseAddress() common.Address {
	panic("implement me")
}

func (fakeCalculableBlock) CoinBase() *big.Int {
	return big.NewInt(6e9)
}

func (fakeCalculableBlock) GetTransactionFees() *big.Int {
	return big.NewInt(15e9)
}

func (fakeCalculableBlock) Seed() common.Hash{
	return common.Hash{}
}

type fakeBlockBroadcaster struct {

}

func (fakeBlockBroadcaster) BroadcastMinedBlock(block model.AbstractBlock) {
	return
}

func fakeMineConfig() MineConfig {
	addr := common.HexToAddress("0x123")
	av := &atomic.Value{}
	av.Store(addr)
	return MineConfig{
		CoinbaseAddress: av,
		BlockBroadcaster: &fakeBlockBroadcaster{},
	}
}

func TestWorkerManager_GetReward(t *testing.T) {
	manager := newDefaultWorkManager(fakeMineConfig())
	block := fakeCalculableBlock{}

	worker1 := common.HexToAddress("0x123")

	worker2 := common.HexToAddress("0x123223")

	assert.EqualValues(t, 0, manager.getPerformance(worker1))
	assert.EqualValues(t, 0, manager.getPerformance(worker2))

	// submits block, the input parameter should be verified
	// before the this function called. So we assume the input
	// block is valid, here we use empty block instead.
	manager.submitBlock(worker1, &model.Block{})
	manager.submitBlock(worker2, &model.Block{})
	manager.submitBlock(worker2, &model.Block{})

	assert.EqualValues(t, 1, manager.getPerformance(worker1))
	assert.EqualValues(t, 2, manager.getPerformance(worker2))

	assert.EqualValues(t, big.NewInt(0), manager.getReward(worker1))
	assert.EqualValues(t, big.NewInt(0), manager.getReward(worker2))

	// receives from previous mined block, assumes the block was
	// mined by the same master, which is verified before this
	// function is invoked.
	manager.onNewBlock(block)

	assert.EqualValues(t, 1, manager.getPerformance(worker1))
	assert.EqualValues(t, 2, manager.getPerformance(worker2))

	// we have only one worker, all the rewards go to that worker
	assert.EqualValues(t, big.NewInt(7e9), manager.getReward(worker1))
	assert.EqualValues(t, big.NewInt(14e9), manager.getReward(worker2))

	// receives previous mined block again
	manager.onNewBlock(block)

	assert.EqualValues(t, big.NewInt(14e9), manager.getReward(worker1))
	assert.EqualValues(t, big.NewInt(28e9), manager.getReward(worker2))
}

func TestWorkerManager_WithdrawReward(t *testing.T) {
	manager := newDefaultWorkManager(fakeMineConfig())
	rcBlock := fakeCalculableBlock{}
	minedBlock := &model.Block{}

	worker1 := common.HexToAddress("0x123")
	worker2 := common.HexToAddress("0x123223")

	manager.submitBlock(worker1, minedBlock)
	manager.submitBlock(worker2, minedBlock)
	manager.submitBlock(worker2, minedBlock)

	manager.onNewBlock(rcBlock)
	manager.onNewBlock(rcBlock)

	assert.EqualValues(t, big.NewInt(14e9), manager.getReward(worker1))
	assert.EqualValues(t, big.NewInt(28e9), manager.getReward(worker2))

	manager.clearReward(worker1)
	manager.clearPerformance(worker1)
	assert.EqualValues(t, big.NewInt(0), manager.getReward(worker1))
	assert.EqualValues(t, 0, manager.getPerformance(worker1))

	manager.subtractReward(worker2, big.NewInt(14e9))
	assert.EqualValues(t, big.NewInt(14e9), manager.getReward(worker2))

	manager.subtractPerformance(worker2, 1)
	assert.EqualValues(t, 1, manager.getPerformance(worker2))
}

type fakeContext struct {
	coinbase *atomic.Value
}

func TestChangeCoinbase(t *testing.T) {
	fc := &fakeContext{ coinbase: &atomic.Value{} }
	conf := MineConfig{ CoinbaseAddress: fc.coinbase }
	m, _ := MakeMineMaster(conf)
	fc.coinbase.Store(common.HexToAddress("0x123"))
	assert.Equal(t, common.HexToAddress("0x123"), m.CurrentCoinbaseAddress())
}
