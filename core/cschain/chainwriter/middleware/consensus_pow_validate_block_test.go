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
	"github.com/dipperin/dipperin-core/common/gerror"
	"github.com/dipperin/dipperin-core/core/chainconfig"
	"github.com/dipperin/dipperin-core/third_party/crypto"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
	"time"
)

func TestValidateBlockNumber(t *testing.T) {
	_, _, _, passChain := getTxTestEnv(t)
	rollBackNum := passChain.GetChainConfig().RollBackNum
	assert.Equal(t, ValidateBlockNumber(&BlockContext{
		Block: &fakeBlock{num: testBlockNum - rollBackNum},
		Chain: passChain,
	})(), gerror.ErrBlockHeightTooLow)
	
	assert.Equal(t, ValidateBlockNumber(&BlockContext{
		Block: &fakeBlock{num: testBlockNum - rollBackNum + 1},
		Chain: passChain,
	})(), gerror.ErrNormalBlockHeightTooLow)
	
	assert.Equal(t, ValidateBlockNumber(&BlockContext{
		Block: &fakeBlock{num: testBlockNum},
		Chain: passChain,
	})(), gerror.ErrNormalBlockHeightTooLow)
	
	assert.NoError(t, ValidateBlockNumber(&BlockContext{
		Block: &fakeBlock{num: testBlockNum + 1},
		Chain: passChain,
	})())
	
	assert.Equal(t, ValidateBlockNumber(&BlockContext{
		Block: &fakeBlock{num: testBlockNum + 2, ts: big.NewInt(time.Now().UnixNano())},
		Chain: passChain,
	})(), gerror.ErrFutureBlock)
	
	assert.Equal(t, ValidateBlockNumber(&BlockContext{
		Block: &fakeBlock{num: testBlockNum + 10, ts: big.NewInt(time.Now().Add(time.Second*maxTimeFutureBlocks + time.Second).UnixNano())},
		Chain: passChain,
	})(), gerror.ErrFutureBlockTooFarAway)
	
	// special block
	assert.Equal(t, ValidateBlockNumber(&BlockContext{
		Block: &fakeBlock{num: testBlockNum - rollBackNum, isSpecial: true},
		Chain: passChain,
	})(), gerror.ErrBlockHeightTooLow)
	
	assert.NoError(t, ValidateBlockNumber(&BlockContext{
		Block: &fakeBlock{num: testBlockNum - rollBackNum + 1, isSpecial: true, ts: big.NewInt(time.Now().UnixNano())},
		Chain: passChain,
	})())
	
	assert.NoError(t, ValidateBlockNumber(&BlockContext{
		Block: &fakeBlock{num: testBlockNum, isSpecial: true, ts: big.NewInt(time.Now().UnixNano())},
		Chain: passChain,
	})())
	
	assert.Equal(t, ValidateBlockNumber(&BlockContext{
		Block: &fakeBlock{num: testBlockNum + 2, isSpecial: true, ts: big.NewInt(time.Now().UnixNano())},
		Chain: passChain,
	})(), gerror.ErrFutureBlock)
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
		Block: &fakeBlock{ts: big.NewInt(tn.Add(chainconfig.GetChainConfig().BlockTimeRestriction + time.Second*30).UnixNano())},
		Chain: &fakeChainInterface{block: &fakeBlock{}},
	})())
}

func TestValidateGasLimit(t *testing.T) {
	gasLimit := chainconfig.BlockGasLimit
	
	nextGasLimit := gasLimit + gasLimit/1024 - 1
	assert.NoError(t, ValidateGasLimit(&BlockContext{
		Block: &fakeBlock{GasLimit: uint64(nextGasLimit)},
		Chain: &fakeChainInterface{block: &fakeBlock{GasLimit: uint64(gasLimit)}},
	})())
	
	nextGasLimit = gasLimit + gasLimit/1024
	assert.Error(t, ValidateGasLimit(&BlockContext{
		Block: &fakeBlock{GasLimit: uint64(nextGasLimit)},
		Chain: &fakeChainInterface{block: &fakeBlock{GasLimit: uint64(gasLimit)}},
	})())
	
	nextGasLimit = gasLimit - gasLimit/1024 + 1
	assert.NoError(t, ValidateGasLimit(&BlockContext{
		Block: &fakeBlock{GasLimit: uint64(nextGasLimit)},
		Chain: &fakeChainInterface{block: &fakeBlock{GasLimit: uint64(gasLimit)}},
	})())
	
	nextGasLimit = gasLimit - gasLimit/1024
	assert.Error(t, ValidateGasLimit(&BlockContext{
		Block: &fakeBlock{GasLimit: uint64(nextGasLimit)},
		Chain: &fakeChainInterface{block: &fakeBlock{GasLimit: uint64(gasLimit)}},
	})())
}