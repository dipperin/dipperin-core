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
	"github.com/dipperin/dipperin-core/common/g-error"
	"github.com/dipperin/dipperin-core/core/chain/state-processor"
	"github.com/dipperin/dipperin-core/core/economy-model"
	"github.com/dipperin/dipperin-core/core/model"
	model2 "github.com/dipperin/dipperin-core/core/vm/model"
	"github.com/dipperin/dipperin-core/tests/g-testData"
	"github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestValidateBlockTxs(t *testing.T) {
	block := &fakeBlock{}
	fakeChain := &fakeChainInterface{}
	blockContext := &BlockContext{Block: block, Chain: fakeChain}
	assert.Equal(t, g_error.ErrTxRootNotMatch, ValidateBlockTxs(blockContext)())

	block.txRoot = common.HexToHash("0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")
	assert.NoError(t, ValidateBlockTxs(blockContext)())

	block.isSpecial = true
	assert.NoError(t, ValidateBlockTxs(blockContext)())

	block.txs = []model.AbstractTransaction{&fakeTx{Receipt: &model2.Receipt{}}}
	block.txRoot = common.HexToHash("0xdcc044ba24b5184502aef321c170b8a52d570190d7b85a9ccfbd5d7b0754d2f8")
	assert.Equal(t, g_error.ErrTxInSpecialBlock, ValidateBlockTxs(blockContext)())

	block.isSpecial = false
	assert.Equal(t, "invalid sender", ValidateBlockTxs(blockContext)().Error())
}

func TestTxValidatorForRpcService_Valid(t *testing.T) {

	assert.Error(t, ValidTxSender(&fakeTx{
		sender: common.Address{0x11},
		fee:    big.NewInt(1),
	}, &fakeChainInterface{}, 0))

	assert.Error(t, ValidTxSender(&fakeTx{
		sender: common.Address{0x11},
		fee:    big.NewInt(100000),
	}, &fakeChainInterface{}, 0))

	adb, _ := NewEmptyAccountDB()
	assert.Error(t, ValidTxSender(&fakeTx{
		sender: common.Address{0x11},
		fee:    big.NewInt(100000),
	}, &fakeChainInterface{
		state: adb,
	}, 1))

	sender := common.Address{0x11}
	assert.NoError(t, adb.NewAccountState(sender))
	assert.NoError(t, adb.AddBalance(sender, big.NewInt(10000011)))
	assert.Error(t, ValidTxSender(&fakeTx{
		sender: sender,
		fee:    big.NewInt(100000),
	}, &fakeChainInterface{
		state: adb,
		block: &fakeBlock{},
		em:    &fakeEconomyModel{},
	}, 1))
	assert.Error(t, ValidTxSender(&fakeTx{
		sender: sender,
		fee:    big.NewInt(100000),
		amount: big.NewInt(10),
	}, &fakeChainInterface{
		state: adb,
		block: &fakeBlock{},
		em:    &fakeEconomyModel{lockM: big.NewInt(10000011)},
	}, 1))
	assert.NoError(t, ValidTxSender(&fakeTx{
		sender:   sender,
		GasLimit: g_testData.TestGasLimit,
		amount:   big.NewInt(10),
	}, &fakeChainInterface{
		state: adb,
		block: &fakeBlock{},
		em:    &fakeEconomyModel{lockM: big.NewInt(0)},
	}, 1))
}

func TestValidTxByType(t *testing.T) {
	_, _, passTx, passChain := getTxTestEnv(t)
	assert.NoError(t, ValidTxByType(passTx, passChain, 0)())

	passTx.txType = 0x9999
	assert.Equal(t, g_error.ErrInvalidTxType, ValidTxByType(passTx, passChain, 0)())

	passTx.txType = common.AddressTypeUnStake
	assert.Equal(t, state_processor.NotEnoughStakeErr, ValidTxByType(passTx, passChain, 0)())

	passTx.txType = common.AddressTypeContractCreate
	assert.NoError(t, ValidTxByType(passTx, passChain, 0)())

	passTx.txType = common.AddressTypeContractCall
	assert.NoError(t, ValidTxByType(passTx, passChain, 0)())
}

func Test_validTx(t *testing.T) {
	_, _, passTx, passChain := getTxTestEnv(t)
	assert.NoError(t, validTx(passTx, passChain, 0))
	passTx.txType = 0x9999
	assert.Equal(t, g_error.ErrInvalidTxType, validTx(passTx, passChain, 0))
	passTx.txType = common.AddressTypeUnStake
	assert.Equal(t, state_processor.NotEnoughStakeErr, validTx(passTx, passChain, 0))
}

func Test_validRegisterTx(t *testing.T) {
	tx := &fakeTx{amount: big.NewInt(100)}
	assert.Equal(t, g_error.ErrDelegatesNotEnough, validRegisterTx(tx, nil, 0))
	tx.amount = economy_model.MiniPledgeValue
	assert.NoError(t, validRegisterTx(tx, nil, 0))
}

func Test_validUnStakeTx(t *testing.T) {
	s, adb, passTx, passChain := getTxTestEnv(t)
	assert.Equal(t, state_processor.NotEnoughStakeErr, validUnStakeTx(passTx, passChain, 0))
	assert.NoError(t, adb.AddStake(s, big.NewInt(100)))
	assert.Equal(t, state_processor.SendCancelTxFirst, validUnStakeTx(passTx, passChain, 0))
	assert.NoError(t, adb.SetLastElect(s, uint64(1000)))
	assert.NoError(t, validUnStakeTx(passTx, passChain, 0))
}

func Test_validCancelTx(t *testing.T) {
	s, adb, passTx, passChain := getTxTestEnv(t)
	assert.Error(t, validCancelTx(passTx, passChain, 0))
	passTx.sender = common.Address{}
	assert.Error(t, validCancelTx(passTx, passChain, 0))
	passTx.sender = s
	passChain.state = nil
	assert.Error(t, validCancelTx(passTx, passChain, 0))
	passChain.state, _ = NewEmptyAccountDB()
	assert.Error(t, validCancelTx(passTx, passChain, 0))
	passChain.state = adb
	assert.NoError(t, adb.AddStake(s, big.NewInt(11)))
	assert.NoError(t, validCancelTx(passTx, passChain, 0))
	assert.NoError(t, adb.SetLastElect(s, 1))
	assert.Error(t, validCancelTx(passTx, passChain, 0))
}

func Test_validContractTx(t *testing.T) {
	assert.Error(t, validContractTx(&fakeTx{}, &fakeChainInterface{}, 0))
	s, _ := NewEmptyAccountDB()
	assert.Error(t, validContractTx(&fakeTx{}, &fakeChainInterface{state: s}, 0))
}

func Test_validEarlyTokenTx(t *testing.T) {
	assert.Nil(t, validEarlyTokenTx(nil, nil, 0))
}

func Test_validEvidenceTx(t *testing.T) {
	assert.Error(t, validEvidenceTx(&fakeTx{extraData: []byte{}}, &fakeChainInterface{}, 0))

	a, p := getPassConflictVote()
	pb, err := rlp.EncodeToBytes(p)
	assert.NoError(t, err)
	tmpAddr := a.Address()
	assert.Error(t, validEvidenceTx(&fakeTx{extraData: pb, to: &tmpAddr}, &fakeChainInterface{}, 0))

	_, adb, passTx, passChain := getTxTestEnv(t)
	passTx.extraData = pb
	passTx.to = &tmpAddr
	assert.NoError(t, adb.NewAccountState(tmpAddr))
	assert.Error(t, validEvidenceTx(passTx, passChain, 0))
	assert.NoError(t, adb.AddStake(tmpAddr, big.NewInt(100)))
	assert.NoError(t, validEvidenceTx(passTx, passChain, 0))
}

func Test_conflictVote(t *testing.T) {
	a, p := getPassConflictVote()
	pb, err := rlp.EncodeToBytes(p)
	assert.NoError(t, err)
	tmpAddr := common.Address{0x12}
	assert.Error(t, conflictVote(&fakeTx{extraData: pb, to: &tmpAddr}, &fakeChainInterface{}, 0))

	p.VoteB = a.getVoteMsg(0, 1, common.Hash{}, model.VoteMessage)
	pb, err = rlp.EncodeToBytes(p)
	assert.Error(t, conflictVote(&fakeTx{extraData: pb}, &fakeChainInterface{}, 0))

	p.VoteB.Height = 3
	pb, err = rlp.EncodeToBytes(p)
	assert.Error(t, conflictVote(&fakeTx{extraData: pb}, &fakeChainInterface{}, 0))

	p.VoteA.Height = 2
	pb, err = rlp.EncodeToBytes(p)
	assert.Error(t, conflictVote(&fakeTx{extraData: pb}, &fakeChainInterface{}, 0))

	assert.Error(t, conflictVote(&fakeTx{extraData: []byte{}}, &fakeChainInterface{}, 0))
}

func Test_validEvidenceTime(t *testing.T) {
	_, adb, passTx, passChain := getTxTestEnv(t)
	to := common.Address{0x12}
	passTx.to = &to
	assert.Error(t, validEvidenceTime(passTx, passChain, 0))

	norTo := cs_crypto.GetNormalAddressFromEvidence(to)
	assert.NoError(t, adb.NewAccountState(norTo))
	assert.NoError(t, adb.SetLastElect(norTo, 1))
	assert.Error(t, validEvidenceTime(passTx, passChain, 0))
}

func Test_validTargetStake(t *testing.T) {
	s, adb, passTx, passChain := getTxTestEnv(t)
	passTx.to = &s
	assert.Error(t, validTargetStake(passTx, passChain, 0))

	target := cs_crypto.GetNormalAddressFromEvidence(s)
	assert.NoError(t, adb.NewAccountState(target))
	assert.Error(t, validTargetStake(passTx, passChain, 0))

	assert.NoError(t, adb.AddStake(target, big.NewInt(10)))
	assert.NoError(t, validTargetStake(passTx, passChain, 0))
}

func Test_validUnStakeTime(t *testing.T) {
	assert.Error(t, validUnStakeTime(&fakeTx{}, &fakeChainInterface{}, 0))
	assert.Error(t, validUnStakeTime(&fakeTx{sender: common.Address{0x12}}, &fakeChainInterface{}, 0))
	adb, _ := NewEmptyAccountDB()
	assert.Error(t, validUnStakeTime(&fakeTx{sender: common.Address{0x12}}, &fakeChainInterface{state: adb}, 0))
	assert.NoError(t, adb.NewAccountState(common.Address{0x12}))
	assert.Error(t, validUnStakeTime(&fakeTx{sender: common.Address{0x12}}, &fakeChainInterface{state: adb}, 0))

	assert.NoError(t, adb.AddStake(common.Address{0x12}, big.NewInt(12)))
	assert.NoError(t, adb.SetLastElect(common.Address{0x12}, 12))
	assert.Error(t, validUnStakeTime(&fakeTx{sender: common.Address{0x12}}, &fakeChainInterface{
		state: adb,
		block: &fakeBlock{},
	}, 0))
}


