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

package middleware

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/bloom"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
	"time"
)

func TestValidateBlockNumber(t *testing.T) {
	_, _, _, passChain := getTxTestEnv(t)
	assert.Error(t, ValidateBlockNumber(&BlockContext{
		Block: &fakeBlock{num: 10, ts: big.NewInt(time.Now().UnixNano())},
		Chain: passChain,
	})())
	assert.Error(t, ValidateBlockNumber(&BlockContext{
		Block: &fakeBlock{num: 10, ts: big.NewInt(time.Now().Add(time.Second*maxTimeFutureBlocks + time.Second).UnixNano())},
		Chain: passChain,
	})())
}

func TestValidateBlockHash(t *testing.T) {
	_, _, _, passChain := getTxTestEnv(t)
	assert.NoError(t, ValidateBlockHash(&BlockContext{
		Block: &fakeBlock{preHash: common.Hash{}},
		Chain: passChain,
	})())

	assert.Error(t, ValidateBlockHash(&BlockContext{
		Block: &fakeBlock{preHash: common.Hash{0x12}},
		Chain: passChain,
	})())

	passChain.block = nil
	assert.Error(t, ValidateBlockHash(&BlockContext{
		Block: &fakeBlock{preHash: common.Hash{}},
		Chain: passChain,
	})())
}

func TestValidBlockSize(t *testing.T) {
	assert.Error(t, ValidateBlockSize(&BlockContext{
		Block: &fakeWrongBlock{X: 1},
		Chain: &fakeChainInterface{},
	})())

	assert.Error(t, ValidateBlockSize(&BlockContext{
		Block: &fakeBlock{ExtraData: make([]byte, chain_config.MaxBlockSize+1)},
		Chain: &fakeChainInterface{},
	})())

	assert.NoError(t, ValidateBlockSize(&BlockContext{
		Block: &fakeBlock{},
		Chain: &fakeChainInterface{},
	})())
}

func TestValidateBlockDifficulty(t *testing.T) {
	nt := time.Now()
	assert.Error(t, ValidateBlockDifficulty(&BlockContext{
		Block: &fakeBlock{ts: big.NewInt(nt.UnixNano())},
		Chain: &fakeChainInterface{block: &fakeBlock{ts: big.NewInt(nt.Add(-time.Second).UnixNano())}},
	})())

	assert.NoError(t, ValidateBlockDifficulty(&BlockContext{
		Block: &fakeBlock{isSpecial: true},
		Chain: &fakeChainInterface{block: &fakeBlock{ts: big.NewInt(nt.Add(-time.Second).UnixNano())}},
	})())

	assert.NoError(t, ValidateBlockDifficulty(&BlockContext{
		Block: &fakeBlock{ts: big.NewInt(nt.UnixNano()), diff: common.HexToDiff("0x1f3fffff"), num: 2},
		Chain: &fakeChainInterface{block: &fakeBlock{ts: big.NewInt(nt.Add(-time.Second).UnixNano()), diff: common.HexToDiff("0x1f3fffff"), num: 1}},
	})())

	assert.Error(t, ValidateBlockDifficulty(&BlockContext{
		Block: &fakeBlock{ts: big.NewInt(nt.UnixNano()), diff: common.HexToDiff("0x1f3fffff"), num: 2, preHash: common.Hash{0x21}},
		Chain: &fakeChainInterface{block: &fakeBlock{ts: big.NewInt(nt.Add(-time.Second).UnixNano()), diff: common.HexToDiff("0x1f3fffff"), num: 1}},
	})())
}

func TestValidateBlockCoinBase(t *testing.T) {
	assert.NoError(t, ValidateBlockCoinBase(&BlockContext{
		Block: &fakeBlock{},
		Chain: &fakeChainInterface{},
	})())

	assert.Error(t, ValidateBlockCoinBase(&BlockContext{
		Block: &fakeBlock{isSpecial: true, cb: common.Address{0x12}},
		Chain: &fakeChainInterface{},
	})())
}

func TestValidateSeed(t *testing.T) {
	assert.Error(t, ValidateSeed(&BlockContext{
		Block: &fakeBlock{},
		Chain: &fakeChainInterface{block: &fakeBlock{}},
	})())

	a := NewAccount()
	aPk := crypto.FromECDSAPub(&a.Pk.PublicKey)
	seed, proof := crypto.Evaluate(a.Pk, common.Hash{}.Bytes())
	assert.Error(t, ValidateSeed(&BlockContext{
		Block: &fakeBlock{
			seed:  seed,
			proof: proof,
			mPk:   aPk,
		},
		Chain: &fakeChainInterface{block: &fakeBlock{seed: common.Hash{0x12}}},
	})())

	assert.Error(t, ValidateSeed(&BlockContext{
		Block: &fakeBlock{
			seed:  seed,
			proof: proof,
			mPk:   aPk,
		},
		Chain: &fakeChainInterface{block: &fakeBlock{seed: common.Hash{}}},
	})())
}

func TestValidateBlockVersion(t *testing.T) {
	assert.NoError(t, ValidateBlockVersion(&BlockContext{
		Block: &fakeBlock{},
		Chain: &fakeChainInterface{block: &fakeBlock{}},
	})())

	assert.Error(t, ValidateBlockVersion(&BlockContext{
		Block: &fakeBlock{version: 999},
		Chain: &fakeChainInterface{block: &fakeBlock{}},
	})())

	tn := time.Now()
	assert.NoError(t, ValidateBlockTime(&BlockContext{
		Block: &fakeBlock{ts: big.NewInt(tn.UnixNano())},
		Chain: &fakeChainInterface{block: &fakeBlock{ts: big.NewInt(tn.Add(-time.Second).UnixNano())}},
	})())

	assert.Error(t, ValidateBlockTime(&BlockContext{
		Block: &fakeBlock{ts: big.NewInt(tn.Add(chain_config.GetChainConfig().BlockTimeRestriction + time.Second*30).UnixNano())},
		Chain: &fakeChainInterface{block: &fakeBlock{}},
	})())
}

type fakeWrongBlock struct {
	X int
}

func (f *fakeWrongBlock) SetReceiptHash(receiptHash common.Hash) {
	panic("implement me")
}

func (f *fakeWrongBlock) GetReceiptHash() common.Hash {
	panic("implement me")
}

func (f *fakeWrongBlock) Body() model.AbstractBody {
	panic("implement me")
}

func (f *fakeWrongBlock) CoinBase() *big.Int {
	panic("implement me")
}

func (f *fakeWrongBlock) CoinBaseAddress() common.Address {
	panic("implement me")
}

func (f *fakeWrongBlock) Difficulty() common.Difficulty {
	panic("implement me")
}

func (f *fakeWrongBlock) EncodeRlpToBytes() ([]byte, error) {
	panic("implement me")
}

func (f *fakeWrongBlock) FormatForRpc() interface{} {
	panic("implement me")
}

func (f *fakeWrongBlock) GetAbsTransactions() []model.AbstractTransaction {
	panic("implement me")
}

func (f *fakeWrongBlock) GetBlockTxsBloom() *iblt.Bloom {
	panic("implement me")
}

func (f *fakeWrongBlock) GetBloom() iblt.Bloom {
	panic("implement me")
}

func (f *fakeWrongBlock) GetEiBloomBlockData(reqEstimator *iblt.HybridEstimator) *model.BloomBlockData {
	panic("implement me")
}

func (f *fakeWrongBlock) GetInterLinkRoot() (root common.Hash) {
	panic("implement me")
}

func (f *fakeWrongBlock) GetInterlinks() model.InterLink {
	panic("implement me")
}

func (f *fakeWrongBlock) GetRegisterRoot() common.Hash {
	panic("implement me")
}

func (f *fakeWrongBlock) GetTransactionFees() *big.Int {
	panic("implement me")
}

func (f *fakeWrongBlock) GetTransactions() []*model.Transaction {
	panic("implement me")
}

func (f *fakeWrongBlock) GetVerifications() []model.AbstractVerification {
	panic("implement me")
}

func (f *fakeWrongBlock) Hash() common.Hash {
	panic("implement me")
}

func (f *fakeWrongBlock) Header() model.AbstractHeader {
	panic("implement me")
}

func (f *fakeWrongBlock) IsSpecial() bool {
	panic("implement me")
}

func (f *fakeWrongBlock) Nonce() common.BlockNonce {
	panic("implement me")
}

func (f *fakeWrongBlock) Number() uint64 {
	panic("implement me")
}

func (f *fakeWrongBlock) PreHash() common.Hash {
	panic("implement me")
}

func (f *fakeWrongBlock) RefreshHashCache() common.Hash {
	panic("implement me")
}

func (f *fakeWrongBlock) Seed() common.Hash {
	panic("implement me")
}

func (f *fakeWrongBlock) SetInterLinkRoot(root common.Hash) {
	panic("implement me")
}

func (f *fakeWrongBlock) SetInterLinks(inter model.InterLink) {
	panic("implement me")
}

func (f *fakeWrongBlock) SetNonce(nonce common.BlockNonce) {
	panic("implement me")
}

func (f *fakeWrongBlock) SetRegisterRoot(root common.Hash) {
	panic("implement me")
}

func (f *fakeWrongBlock) SetStateRoot(root common.Hash) {
	panic("implement me")
}

func (f *fakeWrongBlock) SetVerifications(vs []model.AbstractVerification) {
	panic("implement me")
}

func (f *fakeWrongBlock) StateRoot() common.Hash {
	panic("implement me")
}

func (f *fakeWrongBlock) Timestamp() *big.Int {
	panic("implement me")
}

func (f *fakeWrongBlock) TxCount() int {
	panic("implement me")
}

func (f *fakeWrongBlock) TxIterator(cb func(int, model.AbstractTransaction) error) error {
	panic("implement me")
}

func (f *fakeWrongBlock) TxRoot() common.Hash {
	panic("implement me")
}

func (f *fakeWrongBlock) VerificationRoot() common.Hash {
	panic("implement me")
}

func (f *fakeWrongBlock) VersIterator(func(int, model.AbstractVerification, model.AbstractBlock) error) error {
	panic("implement me")
}

func (f *fakeWrongBlock) Version() uint64 {
	panic("implement me")
}
