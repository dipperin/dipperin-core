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

package cs_chain_test

import (
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/core/cs-chain"
	"github.com/dipperin/dipperin-core/core/cs-chain/chain-state"
	"github.com/dipperin/dipperin-core/core/cs-chain/chain-writer"
	"github.com/dipperin/dipperin-core/core/cs-chain/chain-writer/middleware"
	"github.com/dipperin/dipperin-core/core/economy-model"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/tests"
	"github.com/dipperin/dipperin-core/tests/g-testData"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

var (
	testVerBootAccounts  []tests.Account
	testVerifierAccounts []tests.Account
	alicePriv            = "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232031"
	aliceAddr            = common.HexToAddress("0x00005586B883Ec6dd4f8c26063E18eb4Bd228e59c3E9")
)

var testFee = economy_model.GetMinimumTxFee(20001)

func init() {
	log.Info("change ver boot node address for test")
	var err error
	testVerBootAccounts, err = tests.ChangeVerBootNodeAddress()
	if err != nil {
		panic("change verifier boot node address error for test")
	}
	testVerifierAccounts, err = tests.ChangeVerifierAddress(nil)
	if err != nil {
		panic("change verifier address error for test")
	}
}

func CreateCsChain() *cs_chain.CacheChainState {
	f := chain_writer.NewChainWriterFactory()
	chainState, _ := cs_chain.NewCacheChainState(chain_state.NewChainState(&chain_state.ChainStateConfig{
		DataDir:       "",
		WriterFactory: f,
		ChainConfig:   chain_config.GetChainConfig(),
	}))
	f.SetChain(chainState)

	tests.NewGenesisEnv(chainState.GetChainDB(), chainState.GetStateStorage(), nil)
	return chainState
}

func CreateFullBlock(chain *cs_chain.CacheChainState, txs []*model.Transaction, votes []model.AbstractVerification) model.AbstractBlock {
	key1, _ := crypto.HexToECDSA(alicePriv)
	bb := &tests.BlockBuilder{
		ChainState: chain,
		PreBlock:   chain.CurrentBlock(),
		Txs:        txs,
		Vers:       votes,
		MinerPk:    key1,
	}
	return bb.Build()
}

func CreateSpecialBlock(chain *cs_chain.CacheChainState, votes []model.AbstractVerification) model.AbstractBlock {
	bb := &tests.BlockBuilder{
		ChainState: chain,
		PreBlock:   chain.CurrentBlock(),
		Txs:        nil,
		Vers:       votes,
		MinerPk:    testVerBootAccounts[0].Pk,
	}
	return bb.BuildSpecialBlock()
}

func CreateVerifiersVotes(block model.AbstractBlock, votesNum int) (votes []model.AbstractVerification) {
	for i := 0; i < votesNum; i++ {
		voteA := model.NewVoteMsg(block.Number(), uint64(0), block.Hash(), model.VoteMessage)
		sign, _ := crypto.Sign(voteA.Hash().Bytes(), testVerifierAccounts[i].Pk)
		voteA.Witness.Address = testVerifierAccounts[i].Address()
		voteA.Witness.Sign = sign
		votes = append(votes, voteA)
	}
	return
}

func CreateVerBootVote(block model.AbstractBlock) model.AbstractVerification {
	voteA := model.NewVoteMsg(block.Number(), uint64(0), block.Hash(), model.VerBootNodeVoteMessage)
	sign, _ := crypto.Sign(voteA.Hash().Bytes(), testVerBootAccounts[0].Pk)
	voteA.Witness.Address = testVerBootAccounts[0].Address()
	voteA.Witness.Sign = sign
	return voteA
}

func CreateSignedTx(nonce uint64, to common.Address, amount *big.Int, account tests.Account) model.AbstractTransaction {
	fs1 := model.NewMercurySigner(big.NewInt(1))
	tx := model.NewTransaction(nonce, to, amount, g_testData.TestGasPrice, g_testData.TestGasLimit, []byte{})
	signedTx, _ := tx.SignTx(account.Pk, fs1)
	return signedTx
}

func CreateSignedRegisterTx(nonce uint64, amount *big.Int, account tests.Account) model.AbstractTransaction {
	fs1 := model.NewMercurySigner(big.NewInt(1))
	tx := model.NewRegisterTransaction(nonce, amount, g_testData.TestGasPrice, g_testData.TestGasLimit)
	signedTx, _ := tx.SignTx(account.Pk, fs1)
	return signedTx
}

func TestCreateCsChain(t *testing.T) {
	chain := CreateCsChain()
	block := chain.CurrentBlock()
	assert.Equal(t, block.Number(), uint64(0))
}

func TestCreateFullBlock(t *testing.T) {
	chain := CreateCsChain()
	block := CreateFullBlock(chain, nil, nil)
	assert.Equal(t, block.Number(), uint64(1))
	assert.Equal(t, block.IsSpecial(), false)
	assert.Equal(t, block.CoinBaseAddress(), aliceAddr)
}

func TestCreateSpecialBlock(t *testing.T) {
	chain := CreateCsChain()
	block := CreateSpecialBlock(chain, nil)
	assert.Equal(t, block.Number(), uint64(1))
	assert.Equal(t, block.IsSpecial(), true)
	assert.Equal(t, block.CoinBaseAddress(), testVerBootAccounts[0].Address())
}

func TestCreateVerifiersVotes(t *testing.T) {
	config := chain_config.GetChainConfig()
	chain := CreateCsChain()
	block := CreateFullBlock(chain, nil, nil)
	votes := CreateVerifiersVotes(block, config.VerifierNumber*2/3+1)

	// create block processor
	context := middleware.NewBftBlockContext(block, votes, chain)
	err := context.Process(middleware.ValidateVotes(context))
	assert.NoError(t, err)
}

func TestCreateVerBootVotes(t *testing.T) {
	chain := CreateCsChain()
	block := CreateSpecialBlock(chain, nil)
	vote := CreateVerBootVote(block)

	// create block processor
	context := middleware.NewBftBlockContext(block, []model.AbstractVerification{vote}, chain)
	err := context.Process(middleware.ValidateVotes(context))
	assert.NoError(t, err)
}

func TestCreateSignedTx(t *testing.T) {
	tx := CreateSignedTx(0, aliceAddr, big.NewInt(10000), testVerifierAccounts[0])
	chain := CreateCsChain()
	block := CreateFullBlock(chain, []*model.Transaction{tx.(*model.Transaction)}, nil)

	// create block processor
	context := middleware.NewBftBlockContext(block, nil, chain)
	err := context.Process(middleware.ValidateBlockTxs(&context.BlockContext))
	assert.NoError(t, err)
}

func insertBlockToChain(t *testing.T, chain *cs_chain.CacheChainState, num int) {
	curNum := int(chain.CurrentBlock().Number())
	config := chain_config.GetChainConfig()
	for i := curNum; i < curNum+num; i++ {
		curBlock := chain.CurrentBlock()
		var block model.AbstractBlock
		if curBlock.Number() == 0 {
			block = CreateFullBlock(chain, nil, nil)
		} else {

			// votes for curBlock on chain
			var curBlockVotes []model.AbstractVerification
			if curBlock.IsSpecial() {
				vote := CreateVerBootVote(curBlock)
				curBlockVotes = append(curBlockVotes, vote)
			} else {
				curBlockVotes = CreateVerifiersVotes(curBlock, config.VerifierNumber*2/3+1)
			}

			block = CreateFullBlock(chain, nil, curBlockVotes)
		}

		// votes for build block
		votes := CreateVerifiersVotes(block, config.VerifierNumber*2/3+1)
		err := chain.SaveBftBlock(block, votes)
		assert.NoError(t, err)
		assert.Equal(t, uint64(i+1), chain.CurrentBlock().Number())
		assert.Equal(t, false, chain.CurrentBlock().IsSpecial())
	}
}

func insertSpecialBlockToChain(t *testing.T, chain *cs_chain.CacheChainState, num int) {
	curNum := int(chain.CurrentBlock().Number())
	config := chain_config.GetChainConfig()
	for i := curNum; i < curNum+num; i++ {
		// curBlock on chain
		curBlock := chain.CurrentBlock()
		var block model.AbstractBlock
		if curBlock.Number() == 0 {
			block = CreateSpecialBlock(chain, nil)
		} else {
			// votes for curBlock on chain
			var curBlockVotes []model.AbstractVerification
			if curBlock.IsSpecial() {
				vote := CreateVerBootVote(curBlock)
				curBlockVotes = append(curBlockVotes, vote)
			} else {
				curBlockVotes = CreateVerifiersVotes(curBlock, config.VerifierNumber*2/3+1)
			}

			block = CreateSpecialBlock(chain, curBlockVotes)
		}

		// votes for build block
		vote := CreateVerBootVote(block)
		err := chain.SaveBftBlock(block, []model.AbstractVerification{vote})
		assert.NoError(t, err)
		assert.Equal(t, uint64(i+1), chain.CurrentBlock().Number())
		assert.Equal(t, true, chain.CurrentBlock().IsSpecial())
	}
}

func TestBftSaveBlock(t *testing.T) {
	//todo test not pass, waiting for debugging
	//t.Skip("")
	model.IgnoreDifficultyValidation = true
	chain := CreateCsChain()
	insertBlockToChain(t, chain, 1)
	insertSpecialBlockToChain(t, chain, 3)
	insertBlockToChain(t, chain, 1)
}

func TestCacheChainState_Rollback(t *testing.T) {
	chain := CreateCsChain()
	insertBlockToChain(t, chain, 3)

	// create special block
	config := chain_config.GetChainConfig()
	curBlockVote := CreateVerifiersVotes(chain.CurrentBlock(), config.VerifierNumber*2/3+1)
	specialBlock := CreateSpecialBlock(chain, curBlockVote)

	// make a fork
	insertBlockToChain(t, chain, 1)
	nBlock := chain.CurrentBlock()
	assert.Equal(t, nBlock.Number(), uint64(4))
	assert.Equal(t, nBlock.IsSpecial(), false)

	// reverse chain
	seenCommit := CreateVerBootVote(specialBlock)
	err := chain.SaveBftBlock(specialBlock, []model.AbstractVerification{seenCommit})
	assert.NoError(t, err)
	assert.Equal(t, chain.CurrentBlock().Hash(), specialBlock.Hash())
	assert.Equal(t, chain.CurrentBlock().IsSpecial(), true)
	assert.Equal(t, chain.GetBlockByNumber(4).Hash(), specialBlock.Hash())
	assert.Equal(t, chain.GetBlockByNumber(4).IsSpecial(), true)
	assert.Equal(t, chain.HasBlock(nBlock.Hash(), uint64(4)), true)
}

func TestCacheChainState_GetSlot_GetLastChangePoint(t *testing.T) {
	t.Skip()
	model.IgnoreDifficultyValidation = true
	chain := CreateCsChain()
	insertBlockToChain(t, chain, 10)

	num := chain.GetLastChangePoint(chain.CurrentBlock())
	slot := chain.GetSlot(chain.CurrentBlock())
	assert.Equal(t, num, uint64(9))
	assert.Equal(t, slot, uint64(1))
}

func TestCacheChainState_GetNumBySlot(t *testing.T) {
	t.Skip()
	model.IgnoreDifficultyValidation = true
	chain := CreateCsChain()
	insertBlockToChain(t, chain, 20)

	num := chain.GetNumBySlot(0)
	assert.Equal(t, num, uint64(9))

	num = chain.GetNumBySlot(1)
	assert.Equal(t, num, uint64(19))

	num = chain.GetNumBySlot(2)
	assert.Equal(t, num, nil)
}

func TestCacheChainState_CurrentBlock(t *testing.T) {
	chain := CreateCsChain()
	fmt.Println(chain.CurrentBlock())
}
