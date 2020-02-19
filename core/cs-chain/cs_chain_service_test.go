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

package cs_chain

import (
	"errors"
	"fmt"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/core/chain/chaindb"
	"github.com/dipperin/dipperin-core/core/economy-model"
	"github.com/dipperin/dipperin-core/tests"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/golang/mock/gomock"
	"math/big"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/dipperin/dipperin-core/cmd/utils"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/g-error"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/core/cs-chain/chain-state"
	"github.com/dipperin/dipperin-core/core/cs-chain/chain-writer"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/core/tx-pool"
	"github.com/stretchr/testify/assert"
)

type fakeCacheDB struct{}

func (c *fakeCacheDB) DeleteSeenCommits(blockHeight uint64, blockHash common.Hash) error {
	panic("implement me")
}

func (c *fakeCacheDB) GetSeenCommits(blockHeight uint64, blockHash common.Hash) (result []model.AbstractVerification, err error) {
	return nil, nil
}

func (c *fakeCacheDB) SaveSeenCommits(blockHeight uint64, blockHash common.Hash, commits []model.AbstractVerification) error {
	return nil
}

type fakeTxPool struct{}

func (t *fakeTxPool) AddRemotes(txs []model.AbstractTransaction) []error {
	return []error{errors.New("add remotes failed")}
}

func (t *fakeTxPool) Reset(oldHead, newHead *model.Header) {
	return
}

func CsChainServiceBuilder() *CsChainService {
	f := chain_writer.NewChainWriterFactory()
	conf := chain_config.GetChainConfig()
	conf.SlotSize = 3
	conf.VerifierNumber = 3

	homeDir := util.HomeDir()
	dataDir := filepath.FromSlash(homeDir + "/tmp/cs_chain_service_test")
	csConfig := &chain_state.ChainStateConfig{
		DataDir:       dataDir,
		WriterFactory: f,
		ChainConfig:   conf,
	}
	defer os.RemoveAll(dataDir)

	utils.SetupGenesis(dataDir, conf)

	return NewCsChainService(&CsChainServiceConfig{
		CacheDB: &fakeCacheDB{},
		TxPool:  &fakeTxPool{},
	}, chain_state.NewChainState(csConfig))
}

func TestNewCsChainService(t *testing.T) {
	f := chain_writer.NewChainWriterFactory()
	conf := chain_config.GetChainConfig()
	conf.SlotSize = 3
	conf.VerifierNumber = 3

	homeDir := util.HomeDir()
	dataDir := filepath.FromSlash(homeDir + "/tmp/cs_chain_service_test")
	csConfig := &chain_state.ChainStateConfig{
		DataDir:       "",
		WriterFactory: f,
		ChainConfig:   conf,
	}
	defer os.RemoveAll(dataDir)

	assert.Panics(t, func() {
		tmpState := chain_state.NewChainState(csConfig)
		defer tmpState.GetDB().Close()
		NewCsChainService(&CsChainServiceConfig{
			CacheDB: &fakeCacheDB{},
			TxPool:  &fakeTxPool{},
		}, tmpState)
	})

	utils.SetupGenesis(dataDir, conf)
	csConfig.DataDir = dataDir
	assert.NotNil(t, NewCsChainService(&CsChainServiceConfig{
		CacheDB: &fakeCacheDB{},
		TxPool:  &fakeTxPool{},
	}, chain_state.NewChainState(csConfig)))
}

func TestCsChainService_Stop(t *testing.T) {
	if runtime.GOOS != "linux" {
		return
	}
	ccs := CsChainServiceBuilder()
	ccs.Stop()
}

func TestCsChainService_CurrentBalance(t *testing.T) {
	if runtime.GOOS != "linux" {
		return
	}
	assert.NoError(t, os.Setenv("boots_env", "test"))
	ccs := CsChainServiceBuilder()
	assert.Equal(t, ccs.CurrentBalance(common.HexToAddress("0x1234")), (*big.Int)(nil))

	_, err := ccs.GetTransactionNonce(common.Address{})
	assert.Error(t, err)

	assert.NotEmpty(t, economy_model.DIPProportion.DeveloperProportion)
	for addr := range economy_model.DIPProportion.DeveloperProportion {
		assert.True(t, ccs.CurrentBalance(common.HexToAddress(addr)).Cmp(big.NewInt(0)) > 0)
		_, err = ccs.GetTransactionNonce(common.HexToAddress(addr))
		assert.NoError(t, err)
	}

	ccs = &CsChainService{
		CacheChainState: &CacheChainState{ChainState: &chain_state.ChainState{
			ChainDB: chaindb.NewChainDB(ethdb.NewMemDatabase(), model.MakeDefaultBlockDecoder()),
		}},
	}
	assert.Nil(t, ccs.CurrentBalance(common.Address{}))
	_, err = ccs.GetTransactionNonce(common.Address{})
	assert.Error(t, err)
}

func TestCsChainService_GetSeenCommit(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	mCDB := NewMockCacheDB(controller)

	ccs := &CsChainService{
		CsChainServiceConfig: &CsChainServiceConfig{CacheDB: mCDB},
	}

	mCDB.EXPECT().GetSeenCommits(gomock.Eq(uint64(1)), gomock.Eq(common.Hash{})).Return([]model.AbstractVerification{model.VoteMsg{}}, nil)
	mCDB.EXPECT().GetSeenCommits(gomock.Eq(uint64(2)), gomock.Eq(common.Hash{})).Return(nil, errors.New("failed"))
	assert.Nil(t, ccs.GetSeenCommit(0))
	assert.Len(t, ccs.GetSeenCommit(1), 1)
	assert.Len(t, ccs.GetSeenCommit(2), 0)
}

func TestCsChainService_SaveBlock(t *testing.T) {
	cMock := &fakeCacheDB{}
	ccs, gEnv, txB, bB := getTestChainEnv(cMock)
	bB.SetMinerPk(gEnv.DefaultBootNodeVerifiers()[0].Pk)
	bB.Txs = []*model.Transaction{txB.Build()}

	b1 := bB.Build()
	b1Special := bB.BuildSpecialBlock()

	v := gEnv.VoteBlock(1, 1, b1)
	assert.Equal(t, g_error.ErrBlockVotesNotEnough, ccs.SaveBlock(b1, v))
	v = gEnv.VoteBlock(3, 1, b1)
	assert.NoError(t, ccs.SaveBlock(b1, v))

	bB.PreBlock = b1
	bB.Vers = v
	b2 := bB.Build()
	assert.NoError(t, ccs.SaveBlock(b2, gEnv.VoteBlock(3, 1, b2)))
	assert.NoError(t, ccs.SaveBlock(b1Special, gEnv.VoteSpecialBlock(b1Special)))
	assert.Equal(t, b1Special.Number(), ccs.CurrentBlock().Number())
	assert.Equal(t, b1Special.Number(), ccs.GetBlockByNumber(1).Number())
	assert.Nil(t, ccs.GetBlockByNumber(2))

	numLowBlockToReturnErr = 0
	block := ccs.GetBlockByNumber(0)
	assert.Equal(t, g_error.ErrAlreadyHaveThisBlock, ccs.checkBftBlock(block, nil))
}

func TestCsChainService_checkGenesis(t *testing.T) {
	ccs := &CsChainService{CacheChainState: &CacheChainState{ChainState: &chain_state.ChainState{}}}
	assert.Panics(t, func() {
		ccs.checkGenesis()
	})
	ccs.ChainState = chain_state.NewChainState(&chain_state.ChainStateConfig{
		ChainConfig:   chain_config.GetChainConfig(),
		DataDir:       "",
		WriterFactory: chain_writer.NewChainWriterFactory(),
	})
	ccs.checkGenesis()

	controller := gomock.NewController(t)
	defer controller.Finish()
	sm := NewMockStateStorage(controller)
	sm.EXPECT().OpenTrie(gomock.Any()).Return(nil, errors.New("failed"))
	sm.EXPECT().OpenTrie(gomock.Any()).Return(nil, nil)
	sm.EXPECT().OpenTrie(gomock.Any()).Return(nil, errors.New("failed"))
	sm.EXPECT().DiskDB().Return(ccs.GetDB()).AnyTimes()
	ccs.ChainState.StateStorage = sm
	assert.Panics(t, func() {
		ccs.checkGenesis()
	})
	assert.Panics(t, func() {
		ccs.checkGenesis()
	})
}

func TestCsChainService_initService(t *testing.T) {
	GenesisSetUp = true
	defer func() { GenesisSetUp = false }()
	ccs := NewCsChainService(&CsChainServiceConfig{}, chain_state.NewChainState(&chain_state.ChainStateConfig{
		ChainConfig:   chain_config.GetChainConfig(),
		DataDir:       "",
		WriterFactory: chain_writer.NewChainWriterFactory(),
	}))
	assert.NoError(t, ccs.InitService())
}

func TestCsChainService_handleFutureBlock(t *testing.T) {
	cMock := &fakeCacheDB{}
	ccs, gEnv, _, bB := getTestChainEnv(cMock)
	//fmt.Println(chain.VerifierAddress)
	block := bB.Build()

	ccs.FutureBlocks.Add(common.Hash{}, &futureBlock{block: block, seenCommits: gEnv.VoteBlock(1, 1, block)})
	ccs.FutureBlocks.Add(block.Hash(), &futureBlock{block: block, seenCommits: gEnv.VoteBlock(1, 1, block)})
	ccs.handleFutureBlock()
}

func getTestChainEnv(db CacheDB) (*CsChainService, *tests.GenesisEnv, *tests.TxBuilder, *tests.BlockBuilder) {
	model.IgnoreDifficultyValidation = true
	f := chain_writer.NewChainWriterFactory()

	cConf := chain_config.GetChainConfig()
	cConf.SlotSize = 3
	cConf.VerifierNumber = 3

	chainState := chain_state.NewChainState(&chain_state.ChainStateConfig{
		DataDir:       "",
		WriterFactory: f,
		ChainConfig:   cConf,
	})

	attackEnv := tests.NewGenesisEnv(chainState.GetChainDB(), chainState.GetStateStorage(), nil)

	txB := &tests.TxBuilder{
		Nonce:  0,
		To:     common.HexToAddress(fmt.Sprintf("0x123a%v", 1)),
		Amount: big.NewInt(1),
		Pk:     attackEnv.DefaultVerifiers()[0].Pk,
		Fee:    testFee,
	}
	bb := &tests.BlockBuilder{
		ChainState: chainState,
		PreBlock:   chainState.CurrentBlock(),
		MinerPk:    attackEnv.Miner().Pk,
	}

	txPool := tx_pool.NewTxPool(tx_pool.DefaultTxPoolConfig, *chain_config.GetChainConfig(), chainState)
	ccs := NewCsChainService(&CsChainServiceConfig{
		CacheDB: db,
		TxPool:  txPool,
	}, chainState)
	f.SetChain(ccs.CacheChainState)

	return ccs, attackEnv, txB, bb
}
