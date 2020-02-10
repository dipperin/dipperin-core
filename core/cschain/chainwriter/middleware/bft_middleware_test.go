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
	iblt "github.com/dipperin/dipperin-core/core/bloom"
	chain_true "github.com/dipperin/dipperin-core/core/chain"
	"github.com/dipperin/dipperin-core/core/chain/stateprocessor"
	"github.com/dipperin/dipperin-core/core/chainconfig"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third_party/crypto"
	cs_crypto "github.com/dipperin/dipperin-core/third_party/crypto/cs-crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"log"
	"math/big"
	"testing"
	"time"
)

func TestBftMiddleware(t *testing.T) {
	bc := NewBftBlockContext(nil, nil, nil)
	//assert.Equal(t, true, len(bc.middlewares) == 0)
	assert.Len(t, bc.middlewares, 0)
	bc.Use(CheckBlock(&bc.BlockContext))
	bc.Use(ValidateBlockNumber(&bc.BlockContext))
	bc.Use(UpdateStateRoot(&bc.BlockContext))
	bc.Use(UpdateBlockVerifier(&bc.BlockContext))
	bc.Use(InsertBlock(&bc.BlockContext))
	//assert.Equal(t, true, len(bc.middlewares) == 5)
	assert.Len(t, bc.middlewares, 5)
	err := bc.Process()
	assert.Error(t, err)
}

var minDiff = common.HexToDiff("0x20ffffff")

func TestBftMiddleware2(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	chain := NewMockChainInterface(ctl)
	validator := NewBftBlockValidator(chain)
	
	skCur, err := crypto.GenerateKey()
	pkCur := skCur.PublicKey
	coinbaseCur := cs_crypto.PubkeyToAddress(pkCur)
	timeStamp := time.Now().Nanosecond()
	headerCur := model.Header{
		Version:     1,
		Number:      9,
		PreHash:     common.Hash{},
		Seed:        common.Hash{},
		Diff:        minDiff,
		TimeStamp:   new(big.Int).SetInt64(int64(timeStamp)),
		CoinBase:    coinbaseCur,
		Nonce:       common.BlockNonceFromInt(432423),
		Bloom:       iblt.NewBloom(model.DefaultBlockBloomConfig),
		GasLimit:    model.DefaultGasLimit,
		Proof:       []byte{},
		MinerPubKey: crypto.FromECDSAPub(&pkCur),
	}
	
	blockCur := model.NewBlock(&headerCur, nil, nil)
	
	verifies := []common.Address{
		common.Address{},
	}
	
	chain.EXPECT().CurrentBlock().Return(blockCur)
	chain.EXPECT().CurrentHeader().Return(blockCur.Header())
	chain.EXPECT().GetBlockByNumber(gomock.Any()).Return(blockCur).AnyTimes()
	chain.EXPECT().GetLatestNormalBlock().Return(blockCur).AnyTimes()
	chain.EXPECT().GetChainConfig().Return(&chainconfig.ChainConfig{
		Version:              1,
		BlockTimeRestriction: blockCacheLimit,
	}).AnyTimes()
	slot := uint64(0)
	chain.EXPECT().GetSlot(blockCur).Return(&slot)
	chain.EXPECT().GetVerifiers(slot).Return(verifies)
	processor, _ := chain_true.NewBlockProcessor(nil, blockCur.StateRoot(), stateprocessor.NewStateStorageWithCache(ethdb.NewMemDatabase()))
	chain.EXPECT().BlockProcessor(blockCur.StateRoot()).Return(processor, nil).AnyTimes()
	chain.EXPECT().GetEconomyModel().Return(nil)
	chain.EXPECT().GetEconomyModel().Return(nil)
	
	sk, err := crypto.GenerateKey()
	pk := sk.PublicKey
	coinbase := cs_crypto.PubkeyToAddress(pk)
	seed, proof := crypto.Evaluate(sk, headerCur.Seed.Bytes())
	header := model.Header{
		Version:     1,
		Number:      10,
		PreHash:     blockCur.Hash(),
		Seed:        seed,
		Diff:        minDiff,
		TimeStamp:   new(big.Int).SetInt64(int64(timeStamp)),
		CoinBase:    coinbase,
		Nonce:       common.BlockNonceFromInt(432423),
		Bloom:       iblt.NewBloom(model.DefaultBlockBloomConfig),
		GasLimit:    model.DefaultGasLimit,
		Proof:       proof,
		MinerPubKey: crypto.FromECDSAPub(&pk),
	}
	
	ver1 := NewMockAbstractVerification(ctl)
	ver1.EXPECT().GetAddress().Return(common.Address{}).AnyTimes()
	ver1.EXPECT().Valid().Return(nil)
	ver1.EXPECT().GetBlockHash().Return(blockCur.Hash().Hex())
	blockVerifier := model.Verifications{
		ver1,
	}
	block := model.NewBlock(&header, nil, blockVerifier)
	log.Println("the block number is:", block.Number())
	chain.EXPECT().StateAtByBlockNumber(block.Number()-1).Return(&stateprocessor.AccountStateDB{}, nil)
	chain.EXPECT().CurrentState().Return(&stateprocessor.AccountStateDB{}, nil)
	err = validator.FullValid(block)
	
	assert.Equal(t, gerror.ErrContractNotExist, err)
	
}
