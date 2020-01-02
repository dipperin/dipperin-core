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

package builder

import (
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/log"
	"github.com/dipperin/dipperin-core/core/chain"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/core/chain/cachedb"
	"github.com/dipperin/dipperin-core/core/chain/registerdb"
	"github.com/dipperin/dipperin-core/core/cs-chain"
	"github.com/dipperin/dipperin-core/core/cs-chain/chain-state"
	"github.com/dipperin/dipperin-core/core/cs-chain/chain-writer"
	"github.com/dipperin/dipperin-core/core/economy-model"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/tests"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

var testFee = economy_model.GetMinimumTxFee(50)

func TestMakeBftBlockBuilder(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	mts := NewMockSigner(controller)
	mps := NewMockPbftSigner(controller)
	mtpool := NewMockTxPool(controller)
	ccs, _, tb := getTestEnv(mtpool)
	builder := MakeBftBlockBuilder(ModelConfig{
		ChainReader:        ccs,
		TxPool:             mtpool,
		PriorityCalculator: model.DefaultPriorityCalculator,
		TxSigner:           mts,
		MsgSigner:          mps,
		ChainConfig:        *ccs.GetChainConfig(),
	})

	tmpMiner := tests.AccFactory.GenAccount()
	mps.EXPECT().PublicKey().Return(&tmpMiner.Pk.PublicKey).AnyTimes()
	seed, proof := crypto.Evaluate(tmpMiner.Pk, common.Hash{}.Bytes())
	mps.EXPECT().Evaluate(gomock.Any(), gomock.Any()).Return(seed, proof, nil).AnyTimes()
	mtpool.EXPECT().Pending().Return(map[common.Address][]model.AbstractTransaction{
		tb.From(): {tb.Build()},
	}, nil).AnyTimes()
	mts.EXPECT().GetSender(gomock.Any()).Return(tb.From(), nil).AnyTimes()
	mts.EXPECT().Equal(gomock.Any()).Return(true).AnyTimes()
	assert.NotNil(t, builder.BuildWaitPackBlock(common.Address{0x12}, 1e1, 1e1))

	tb.Nonce = 99
	mtpool.EXPECT().Pending().Return(map[common.Address][]model.AbstractTransaction{
		tb.From(): {tb.Build()},
	}, nil).AnyTimes()
	mtpool.EXPECT().RemoveTxs(gomock.Any()).AnyTimes()
	assert.NotNil(t, builder.BuildWaitPackBlock(common.Address{0x12}, 1e1, 1e1))

	builder.ChainReader = &fakeChainReader{}
	assert.Panics(t, func() {
		builder.BuildWaitPackBlock(common.Address{0x12}, 1e1, 1e1)
	})
}

func getTestEnv(p TxPool) (*cs_chain.CsChainService, *tests.GenesisEnv, *tests.TxBuilder) {
	conf := chain_config.GetChainConfig()
	conf.SlotSize = 3
	conf.VerifierNumber = 3
	f := chain_writer.NewChainWriterFactory()
	cs := chain_state.NewChainState(&chain_state.ChainStateConfig{
		DataDir:       "",
		WriterFactory: f,
		ChainConfig:   conf,
	})
	f.SetChain(cs)

	env := tests.NewGenesisEnv(cs.GetChainDB(), cs.GetStateStorage(), nil)

	txBuilder := &tests.TxBuilder{
		Nonce:  0,
		To:     common.HexToAddress(fmt.Sprintf("0x123a%v", 1)),
		Amount: big.NewInt(1),
		Pk:     env.DefaultVerifiers()[0].Pk,
		Fee:    testFee,
	}

	ccs := cs_chain.NewCsChainService(&cs_chain.CsChainServiceConfig{
		CacheDB: cachedb.NewCacheDB(ethdb.NewMemDatabase()),
		TxPool:  p,
	}, cs)

	return ccs, env, txBuilder
}

func TestBftBlockBuilder_BuildWaitPackBlock(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	mts := NewMockSigner(controller)
	mps := NewMockPbftSigner(controller)
	mtpool := NewMockTxPool(controller)
	ccs, _, tb := getTestEnv(mtpool)
	builder := MakeBftBlockBuilder(ModelConfig{
		ChainReader:        ccs,
		TxPool:             mtpool,
		PriorityCalculator: model.DefaultPriorityCalculator,
		TxSigner:           mts,
		MsgSigner:          mps,
		ChainConfig:        *ccs.GetChainConfig(),
	})

	tmpMiner := tests.AccFactory.GenAccount()
	mps.EXPECT().PublicKey().Return(&tmpMiner.Pk.PublicKey).AnyTimes()
	seed, proof := crypto.Evaluate(tmpMiner.Pk, common.Hash{}.Bytes())
	mps.EXPECT().Evaluate(gomock.Any(), gomock.Any()).Return(seed, proof, nil).AnyTimes()
	mtpool.EXPECT().Pending().Return(map[common.Address][]model.AbstractTransaction{
		tb.From(): {tb.Build()},
	}, nil)
	mts.EXPECT().GetSender(gomock.Any()).Return(tb.From(), nil).AnyTimes()
	mts.EXPECT().Equal(gomock.Any()).Return(true).AnyTimes()
	assert.NotNil(t, builder.BuildWaitPackBlock(common.Address{0x12}, 1e1, 1e1))
}

func TestBftBlockBuilder_GetDifficulty(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	mts := NewMockSigner(controller)
	mps := NewMockPbftSigner(controller)
	mtpool := NewMockTxPool(controller)
	ccs, _, tb := getTestEnv(mtpool)
	builder := MakeBftBlockBuilder(ModelConfig{
		ChainReader:        ccs,
		TxPool:             mtpool,
		PriorityCalculator: model.DefaultPriorityCalculator,
		TxSigner:           mts,
		MsgSigner:          mps,
		ChainConfig:        *ccs.GetChainConfig(),
	})

	tmpMiner := tests.AccFactory.GenAccount()
	mps.EXPECT().PublicKey().Return(&tmpMiner.Pk.PublicKey).AnyTimes()
	seed, proof := crypto.Evaluate(tmpMiner.Pk, common.Hash{}.Bytes())
	mps.EXPECT().Evaluate(gomock.Any(), gomock.Any()).Return(seed, proof, nil).AnyTimes()
	mtpool.EXPECT().Pending().Return(map[common.Address][]model.AbstractTransaction{
		tb.From(): {tb.Build()},
	}, nil)
	mts.EXPECT().GetSender(gomock.Any()).Return(tb.From(), nil).AnyTimes()
	mts.EXPECT().Equal(gomock.Any()).Return(true).AnyTimes()

	assert.NotNil(t, builder.BuildWaitPackBlock(common.Address{0x12}, 1e1, 1e1))
	assert.NotNil(t, builder.GetDifficulty())
}

func creatBlockWithGasLimitAndGasUsed(gasUsed, gasLimit, number uint64) *model.Block {
	header := &model.Header{
		GasUsed:  gasUsed,
		GasLimit: gasLimit,
		Number:   number,
	}
	return model.NewBlock(header, nil, nil)
}

func TestGasLimitAdjust(t *testing.T) {
	/*	controller := gomock.NewController(t)
		defer controller.Finish()*/
	var blockNumber uint64
	gasLimit := uint64(chain_config.BlockGasLimit)
	gasUsed := uint64(0)
	gasInfo := make(map[uint64]uint64, 11)
	for {
		Block := creatBlockWithGasLimitAndGasUsed(gasUsed, gasLimit, blockNumber)
		if blockNumber%360 == 0 {
			gasInfo[blockNumber] = Block.GasLimit()
		}

		blockNumber++
		gasLimit = CalcGasLimit(Block, 2e5, 2e10)
		//gasLimit = CalcGasLimit(Block,Block.GasLimit(),Block.GasLimit())
		if blockNumber%5 != 0 {
			gasUsed = uint64(0)
		} else {
			gasUsed = gasLimit
		}

		log.DLogger.Info("")
		log.DLogger.Info("********************next block***********************")
		log.DLogger.Info("********************next block***********************")
		log.DLogger.Info("")

		/*if blockNumber == 1 {
			break
		}*/
		//break
		/*		if gasLimit == 5000{
				log.DLogger.Info("the gasLimit is out of the MLimit")
				break
			}*/
		if blockNumber == 20000 {
			break
		}
	}
}

type fakeChainReader struct {
}

func (f *fakeChainReader) CurrentBlock() model.AbstractBlock {
	return nil
}

func (f *fakeChainReader) GetBlockByNumber(number uint64) model.AbstractBlock {
	panic("implement me")
}

func (f *fakeChainReader) GetVerifiers(round uint64) []common.Address {
	panic("implement me")
}

func (f *fakeChainReader) IsChangePoint(block model.AbstractBlock, isProcessPackageBlock bool) bool {
	panic("implement me")
}

func (f *fakeChainReader) GetLastChangePoint(block model.AbstractBlock) *uint64 {
	panic("implement me")
}

func (f *fakeChainReader) GetSlot(block model.AbstractBlock) *uint64 {
	panic("implement me")
}

func (f *fakeChainReader) GetSeenCommit(height uint64) []model.AbstractVerification {
	panic("implement me")
}

func (f *fakeChainReader) GetLatestNormalBlock() model.AbstractBlock {
	panic("implement me")
}

func (f *fakeChainReader) BlockProcessor(root common.Hash) (*chain.BlockProcessor, error) {
	panic("implement me")
}

func (f *fakeChainReader) BuildRegisterProcessor(preRoot common.Hash) (*registerdb.RegisterDB, error) {
	panic("implement me")
}

func (f *fakeChainReader) GetEconomyModel() economy_model.EconomyModel {
	panic("implement me")
}
