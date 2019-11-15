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

package chain_state

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/economy-model"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/tests"
	"github.com/dipperin/dipperin-core/tests/g-testData"
	"github.com/stretchr/testify/assert"
	"gopkg.in/check.v1"
	"math/big"
)

var testFee = economy_model.GetMinimumTxFee(20001)

func (suite *chainWriterSuite) TestChainState_BuildRegisterProcessor(t *check.C) {
	config := suite.chainState.ChainConfig
	cv := suite.chainState.GetCurrVerifiers()
	assert.Len(t, cv, config.VerifierNumber)

	acc1 := tests.AccFactory.GetAccount(0)
	// send coins to this address so that it can initiate a registration transaction
	suite.txBuilder.To = acc1.Address()
	suite.txBuilder.Amount = big.NewInt(0).Add(economy_model.MiniPledgeValue, testFee)
	suite.blockBuilder.Txs = []*model.Transaction{suite.txBuilder.Build()}
	b1 := suite.blockBuilder.Build()
	assert.Equal(t, 1, b1.TxCount())
	b1Votes := suite.env.VoteBlock(config.VerifierNumber, 1, b1)
	err := suite.chainState.SaveBlock(b1)
	assert.NoError(t, err)
	// verify the correctness of the balance after the transaction
	b1S, err := suite.chainState.CurrentState()
	assert.NoError(t, err)
	acc1B, err := b1S.GetBalance(acc1.Address())
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(0).Add(economy_model.MiniPledgeValue, testFee), acc1B)

	suite.txBuilder.To = common.HexToAddress(common.AddressStake)
	suite.txBuilder.Pk = acc1.Pk
	suite.txBuilder.Amount = economy_model.MiniPledgeValue
	suite.blockBuilder.Txs = []*model.Transaction{suite.txBuilder.Build()}
	suite.blockBuilder.Vers = b1Votes
	suite.blockBuilder.PreBlock = suite.chainState.CurrentBlock()

	b2 := suite.blockBuilder.Build()
	assert.Equal(t, 1, b2.TxCount())
	assert.NotNil(t, b2)

	err = suite.chainState.SaveBftBlock(b2, suite.env.VoteBlock(config.VerifierNumber, 1, b2))
	t.Check(err, check.IsNil)

	curB := suite.chainState.CurrentBlock()
	assert.Equal(t, uint64(2), curB.Number())
	rDB, err := suite.chainState.BuildRegisterProcessor(curB.GetRegisterRoot())
	assert.NoError(t, err)
	rAddrs := rDB.GetRegisterData()
	assert.Equal(t, acc1.Address(), rAddrs[0])
}

func (suite *chainWriterSuite) TestChainState_CurrentSeed(t *check.C) {
	s, sn := suite.chainState.CurrentSeed()
	assert.Equal(t, common.Hash{}, s)
	assert.Equal(t, uint64(0), sn)
}

func (suite *chainWriterSuite) TestChainState_IsChangePoint(t *check.C) {

	model.IgnoreDifficultyValidation = true

	block := suite.chainState.CurrentBlock()
	assert.False(t, suite.chainState.IsChangePoint(block, false))

	// insert one slot blocks
	config := suite.chainState.ChainConfig
	suite.InsertBlock(t, int(config.SlotSize))
	block = suite.chainState.GetBlockByNumber(config.SlotSize - 1)
	assert.True(t, suite.chainState.IsChangePoint(block, false))
	assert.True(t, suite.chainState.IsChangePoint(block, true))

	// insert special block
	suite.InsertSpecialBlock(t, int(config.SlotSize))
	block = suite.chainState.CurrentBlock()
	assert.True(t, suite.chainState.IsChangePoint(block, false))
}

func (suite *chainWriterSuite) TestChainState_GetLastChangePoint(t *check.C) {
	n := suite.chainState.GetLastChangePoint(suite.chainState.CurrentBlock())
	assert.Equal(t, uint64(0), *n)

	// insert one slot blocks
	config := suite.chainState.ChainConfig
	suite.InsertBlock(t, int(config.SlotSize))

	curBlock := suite.chainState.CurrentBlock()
	point := suite.chainState.GetLastChangePoint(curBlock)
	assert.Equal(t, config.SlotSize-1, *point)
}

func (suite *chainWriterSuite) TestChainState_GetSlotByNum(t *check.C) {
	n := suite.chainState.GetSlotByNum(suite.chainState.CurrentBlock().Number())
	assert.Equal(t, uint64(0), *n)

	// insert one slot blocks
	config := suite.chainState.ChainConfig
	suite.InsertBlock(t, int(config.SlotSize))

	curBlock := suite.chainState.CurrentBlock()
	n = suite.chainState.GetSlotByNum(curBlock.Number())
	assert.Equal(t, uint64(1), *n)
	//fmt.Println("=========", config.SlotSize, curBlock.Number())
	n = suite.chainState.GetSlotByNum(curBlock.Number() + 5)
	assert.Nil(t, n)
}

func (suite *chainWriterSuite) TestChainState_GetCurrVerifiers(t *check.C) {
	block := suite.chainState.CurrentBlock()
	suite.chainState.ChainDB.DeleteBlock(block.Hash(), 0)

	n := suite.chainState.GetCurrVerifiers()
	assert.Nil(t, n)
}

func (suite *chainWriterSuite) TestChainState_GetVerifiers(t *check.C) {
	slot := suite.chainState.GetSlot(suite.chainState.CurrentBlock())
	m := suite.chainState.GetVerifiers(*slot)
	assert.NotEqual(t, nil, m)
}

func (suite *chainWriterSuite) TestChainState_GetNextVerifiers(t *check.C) {
	n := suite.chainState.GetNextVerifiers()
	config := suite.chainState.ChainConfig
	assert.Len(t, n, config.VerifierNumber)

	block := suite.chainState.CurrentBlock()
	suite.chainState.ChainDB.DeleteBlock(block.Hash(), 0)

	n = suite.chainState.GetNextVerifiers()
	assert.Nil(t, n)
}

func (suite *chainWriterSuite) TestChainState_NumBeforeLastBySlot(t *check.C) {
	n := suite.chainState.NumBeforeLastBySlot(0)
	assert.Equal(t, uint64(0), *n)

	n = suite.chainState.NumBeforeLastBySlot(3)
	assert.Nil(t, n)
}

func (suite *chainWriterSuite) TestChainState_GetNumBySlot(t *check.C) {
	n := suite.chainState.GetNumBySlot(3)
	assert.Nil(t, n)

	config := suite.chainState.ChainConfig
	suite.InsertBlock(t, int(config.SlotSize-2))
	n = suite.chainState.GetNumBySlot(0)
	assert.Nil(t, n)

}

func (suite *chainWriterSuite) TestChainState_CalVerifiers(t *check.C) {
	// insert block
	config := suite.chainState.ChainConfig
	ver, _ := tests.ChangeVerifierAddress(nil)
	tx := createRegisterTX(0, economy_model.MiniPledgeValue, ver[1])
	//sender,_ := tx.Sender(nil)
	//log.Info("the register tx sender address is:","addr",sender.Hex())
	suite.blockBuilder.PreBlock = suite.chainState.CurrentBlock()
	suite.blockBuilder.Txs = []*model.Transaction{tx}
	block := suite.blockBuilder.Build()
	seenCommit := suite.env.VoteBlock(config.VerifierNumber, 1, block)
	err := suite.chainState.SaveBftBlock(block, seenCommit)
	assert.NoError(t, err)

	verifiers := suite.chainState.CalVerifiers(block)
	assert.Equal(t, ver[1].Address().Hex(), verifiers[0].Hex())
	assert.Len(t, verifiers, config.VerifierNumber)

	var txs []*model.Transaction
	for i := 0; i < config.VerifierNumber; i++ {
		tx = createRegisterTX(0, economy_model.MiniPledgeValue, ver[i])
		txs = append(txs, tx)
	}
	suite.blockBuilder.Txs = txs
	suite.blockBuilder.SetPreBlock(block)
	suite.blockBuilder.SetVerifications(seenCommit)
	block = suite.blockBuilder.Build()
	seenCommit = suite.env.VoteBlock(config.VerifierNumber, 1, block)
	err = suite.chainState.SaveBftBlock(block, seenCommit)
	assert.NoError(t, err)

	verifiers = suite.chainState.CalVerifiers(block)
	assert.Len(t, verifiers, config.VerifierNumber)
}

func (suite *chainWriterSuite) TestChainState_getTopVerifiers(t *check.C) {
	config := suite.chainState.ChainConfig
	ver, _ := tests.ChangeVerifierAddress(nil)
	aliceAddr := common.HexToAddress("0x00005586B883Ec6dd4f8c26063E18eb4Bd228e59c3E9")

	// create topAddress
	var topAddress []common.Address
	for i := 0; i < config.VerifierNumber; i++ {
		topAddress = append(topAddress, ver[i].Address())
	}

	// create topPriority
	var topPriority []uint64
	for i := 0; i < config.VerifierNumber; i++ {
		topPriority = append(topPriority, uint64(2*config.VerifierNumber-i))
	}

	// add small priority
	resultAddress, _ := suite.chainState.getTopVerifiers(aliceAddr, 1, topAddress, topPriority)
	assert.False(t, aliceAddr.InSlice(resultAddress))

	// add same address
	resultAddress, _ = suite.chainState.getTopVerifiers(ver[1].Address(), uint64(2*config.VerifierNumber-1), topAddress, topPriority)
	assert.True(t, ver[1].Address().InSlice(resultAddress))

	// add alice address
	resultAddress, _ = suite.chainState.getTopVerifiers(aliceAddr, uint64(2*config.VerifierNumber-1), topAddress, topPriority)
	assert.Equal(t, aliceAddr, resultAddress[1])
}

func createRegisterTX(nonce uint64, amount *big.Int, account tests.Account) *model.Transaction {
	fs1 := model.NewSigner(big.NewInt(1))
	tx := model.NewRegisterTransaction(nonce, amount, g_testData.TestGasPrice, g_testData.TestGasLimit)
	signedTx, _ := tx.SignTx(account.Pk, fs1)
	return signedTx
}
