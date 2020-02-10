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
	"github.com/dipperin/dipperin-core/core/economymodel"
	"github.com/dipperin/dipperin-core/core/model"
	cs_crypto "github.com/dipperin/dipperin-core/third_party/crypto/cs-crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestValidateBlockTxs(t *testing.T) {
	block := &fakeBlock{}
	fakeChain := &fakeChainInterface{block: &fakeBlock{}}
	blockContext := &BlockContext{Block: block, Chain: fakeChain}
	assert.Equal(t, gerror.ErrTxRootNotMatch, ValidateBlockTxs(blockContext)())
	
	block.txRoot = common.HexToHash("0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")
	assert.NoError(t, ValidateBlockTxs(blockContext)())
	
	block.isSpecial = true
	assert.NoError(t, ValidateBlockTxs(blockContext)())
	
	block.txs = []model.AbstractTransaction{&fakeTx{Receipt: &model.Receipt{}}}
	block.txRoot = common.HexToHash("0xdcc044ba24b5184502aef321c170b8a52d570190d7b85a9ccfbd5d7b0754d2f8")
	assert.Equal(t, gerror.ErrTxInSpecialBlock, ValidateBlockTxs(blockContext)())
	
	block.isSpecial = false
	assert.Equal(t, "invalid sender", ValidateBlockTxs(blockContext)().Error())
}

func TestTxValidatorForRpcService_Valid(t *testing.T) {
	conf := newValidTxSenderNeedConfig(&fakeChainInterface{block: &fakeBlock{}}, 0)
	
	assert.Error(t, ValidTxSender(&fakeTx{
		sender: common.Address{0x11},
		fee:    big.NewInt(1),
	}, conf))
	
	assert.Error(t, ValidTxSender(&fakeTx{
		sender: common.Address{0x11},
		fee:    big.NewInt(100000),
	}, conf))
	
	adb, _ := NewEmptyAccountDB()
	fakeChain := &fakeChainInterface{
		state: adb,
		block: &fakeBlock{},
	}
	conf = newValidTxSenderNeedConfig(fakeChain, 1)
	assert.Error(t, ValidTxSender(&fakeTx{
		sender: common.Address{0x11},
		fee:    big.NewInt(100000),
	}, conf))
	
	sender := common.Address{0x11}
	assert.NoError(t, adb.NewAccountState(sender))
	assert.NoError(t, adb.AddBalance(sender, big.NewInt(10000011)))
	fakeChain = &fakeChainInterface{
		state: adb,
		block: &fakeBlock{},
		em:    &fakeEconomyModel{},
	}
	conf = newValidTxSenderNeedConfig(fakeChain, 1)
	assert.Error(t, ValidTxSender(&fakeTx{
		sender: sender,
		fee:    big.NewInt(100000),
	}, conf))
	
	fakeChain = &fakeChainInterface{
		state: adb,
		block: &fakeBlock{},
		em:    &fakeEconomyModel{lockM: big.NewInt(10000011)},
	}
	conf = newValidTxSenderNeedConfig(fakeChain, 1)
	assert.Error(t, ValidTxSender(&fakeTx{
		sender: sender,
		fee:    big.NewInt(100000),
		amount: big.NewInt(10),
	}, conf))
	
	fakeChain = &fakeChainInterface{
		state: adb,
		block: &fakeBlock{},
		em:    &fakeEconomyModel{lockM: big.NewInt(0)},
	}
	conf = newValidTxSenderNeedConfig(fakeChain, 1)
	assert.NoError(t, ValidTxSender(&fakeTx{
		sender:   sender,
		GasLimit: 2 * model.TxGas,
		amount:   big.NewInt(10),
	}, conf))
}

func TestValidTxByType(t *testing.T) {
	_, _, passTx, passChain := getTxTestEnv(t)
	assert.NoError(t, ValidTxByType(passTx, passChain, 0)())
	
	passTx.txType = 0x9999
	assert.Equal(t, gerror.ErrInvalidTxType, ValidTxByType(passTx, passChain, 0)())
	
	passTx.txType = common.AddressTypeUnStake
	assert.Equal(t, gerror.ErrTxSenderStakeNotEnough, ValidTxByType(passTx, passChain, 0)())
	
	passTx.txType = common.AddressTypeContractCreate
	assert.NoError(t, ValidTxByType(passTx, passChain, 0)())
	
	passTx.txType = common.AddressTypeContractCall
	assert.NoError(t, ValidTxByType(passTx, passChain, 0)())
}

func Test_validTx(t *testing.T) {
	_, _, passTx, passChain := getTxTestEnv(t)
	conf := newValidTxSenderNeedConfig(passChain, 0)
	assert.NoError(t, validTx(passTx, conf))
	passTx.txType = 0x9999
	assert.Equal(t, gerror.ErrInvalidTxType, validTx(passTx, conf))
	passTx.txType = common.AddressTypeUnStake
	assert.Equal(t, gerror.ErrTxSenderStakeNotEnough, validTx(passTx, conf))
}

func Test_validRegisterTx(t *testing.T) {
	tx := &fakeTx{amount: big.NewInt(100)}
	conf := newValidTxSenderNeedConfig(&fakeChainInterface{block: &fakeBlock{}}, 0)
	assert.Equal(t, gerror.ErrTxDelegatesNotEnough, validRegisterTx(tx, conf))
	tx.amount = economymodel.MiniPledgeValue
	assert.NoError(t, validRegisterTx(tx, conf))
}

func Test_validUnStakeTx(t *testing.T) {
	s, adb, passTx, passChain := getTxTestEnv(t)
	conf := newValidTxSenderNeedConfig(passChain, 0)
	assert.Equal(t, gerror.ErrTxSenderStakeNotEnough, validUnStakeTx(passTx, conf))
	assert.NoError(t, adb.AddStake(s, big.NewInt(100)))
	assert.Equal(t, gerror.ValidateSendCancelTxFirst, validUnStakeTx(passTx, conf))
	assert.NoError(t, adb.SetLastElect(s, uint64(1000)))
	assert.NoError(t, validUnStakeTx(passTx, conf))
}

func Test_validCancelTx(t *testing.T) {
	s, adb, passTx, passChain := getTxTestEnv(t)
	conf := newValidTxSenderNeedConfig(passChain, 0)
	assert.Error(t, validCancelTx(passTx, conf))
	passTx.sender = common.Address{}
	assert.Error(t, validCancelTx(passTx, conf))
	passTx.sender = s
	passChain.state = nil
	assert.Error(t, validCancelTx(passTx, conf))
	passChain.state, _ = NewEmptyAccountDB()
	assert.Error(t, validCancelTx(passTx, conf))
	passChain.state = adb
	assert.NoError(t, adb.AddStake(s, big.NewInt(11)))
	assert.NoError(t, validCancelTx(passTx, conf))
	assert.NoError(t, adb.SetLastElect(s, 1))
	assert.Error(t, validCancelTx(passTx, conf))
}

func Test_validContractTx(t *testing.T) {
	conf := newValidTxSenderNeedConfig(&fakeChainInterface{block: &fakeBlock{}}, 0)
	assert.Error(t, validContractTx(&fakeTx{}, conf))
	s, _ := NewEmptyAccountDB()
	conf = newValidTxSenderNeedConfig(&fakeChainInterface{state: s, block: &fakeBlock{}}, 0)
	assert.Error(t, validContractTx(&fakeTx{}, conf))
}

func Test_validEarlyTokenTx(t *testing.T) {
	assert.Nil(t, validEarlyTokenTx(nil, nil))
}

func Test_validEvidenceTx(t *testing.T) {
	conf := newValidTxSenderNeedConfig(&fakeChainInterface{block: &fakeBlock{}}, 0)
	assert.Error(t, validEvidenceTx(&fakeTx{extraData: []byte{}}, conf))
	
	a, p := getPassConflictVote()
	pb, err := rlp.EncodeToBytes(p)
	assert.NoError(t, err)
	tmpAddr := a.Address()
	assert.Panics(t, func() { validEvidenceTx(&fakeTx{extraData: pb, to: &tmpAddr}, conf) })
	
	_, adb, passTx, passChain := getTxTestEnv(t)
	passTx.extraData = pb
	passTx.to = &tmpAddr
	assert.NoError(t, adb.NewAccountState(tmpAddr))
	conf = newValidTxSenderNeedConfig(passChain, 0)
	assert.Error(t, validEvidenceTx(passTx, conf))
	assert.NoError(t, adb.AddStake(tmpAddr, big.NewInt(100)))
	assert.NoError(t, validEvidenceTx(passTx, conf))
}

func Test_conflictVote(t *testing.T) {
	a, p := getPassConflictVote()
	pb, err := rlp.EncodeToBytes(p)
	assert.NoError(t, err)
	tmpAddr := common.Address{0x12}
	conf := newValidTxSenderNeedConfig(&fakeChainInterface{block: &fakeBlock{}}, 0)
	assert.Error(t, conflictVote(&fakeTx{extraData: pb, to: &tmpAddr}, conf))
	
	p.VoteB = a.getVoteMsg(0, 1, common.Hash{}, model.VoteMessage)
	pb, err = rlp.EncodeToBytes(p)
	assert.Error(t, conflictVote(&fakeTx{extraData: pb}, conf))
	
	p.VoteB.Height = 3
	pb, err = rlp.EncodeToBytes(p)
	assert.Error(t, conflictVote(&fakeTx{extraData: pb}, conf))
	
	p.VoteA.Height = 2
	pb, err = rlp.EncodeToBytes(p)
	assert.Error(t, conflictVote(&fakeTx{extraData: pb}, conf))
	
	assert.Error(t, conflictVote(&fakeTx{extraData: []byte{}}, conf))
}

func Test_validEvidenceTime(t *testing.T) {
	_, adb, passTx, passChain := getTxTestEnv(t)
	to := common.Address{0x12}
	passTx.to = &to
	conf := newValidTxSenderNeedConfig(passChain, 0)
	assert.Error(t, validEvidenceTime(passTx, conf))
	
	norTo := cs_crypto.GetNormalAddressFromEvidence(to)
	assert.NoError(t, adb.NewAccountState(norTo))
	assert.NoError(t, adb.SetLastElect(norTo, testBlockNum+1))
	assert.Error(t, validEvidenceTime(passTx, conf))
}

func Test_validTargetStake(t *testing.T) {
	s, adb, passTx, passChain := getTxTestEnv(t)
	passTx.to = &s
	conf := newValidTxSenderNeedConfig(passChain, 0)
	assert.Error(t, validTargetStake(passTx, conf))
	
	target := cs_crypto.GetNormalAddressFromEvidence(s)
	assert.NoError(t, adb.NewAccountState(target))
	assert.Error(t, validTargetStake(passTx, conf))
	
	assert.NoError(t, adb.AddStake(target, big.NewInt(10)))
	assert.NoError(t, validTargetStake(passTx, conf))
}

func Test_validUnStakeTime(t *testing.T) {
	conf := newValidTxSenderNeedConfig(&fakeChainInterface{block: &fakeBlock{}}, 0)
	assert.Error(t, validUnStakeTime(&fakeTx{}, conf))
	assert.Panics(t, func() { validUnStakeTime(&fakeTx{sender: common.Address{0x12}}, conf) })
	
	adb, _ := NewEmptyAccountDB()
	conf = newValidTxSenderNeedConfig(&fakeChainInterface{state: adb, block: &fakeBlock{}}, 0)
	assert.Error(t, validUnStakeTime(&fakeTx{sender: common.Address{0x12}}, conf))
	assert.NoError(t, adb.NewAccountState(common.Address{0x12}))
	
	conf = newValidTxSenderNeedConfig(&fakeChainInterface{state: adb, block: &fakeBlock{}}, 0)
	assert.Error(t, validUnStakeTime(&fakeTx{sender: common.Address{0x12}}, conf))
	
	assert.NoError(t, adb.AddStake(common.Address{0x12}, big.NewInt(12)))
	assert.NoError(t, adb.SetLastElect(common.Address{0x12}, 12))
	conf = newValidTxSenderNeedConfig(&fakeChainInterface{
		state: adb,
		block: &fakeBlock{},
	}, 0)
	assert.Error(t, validUnStakeTime(&fakeTx{sender: common.Address{0x12}}, conf))
}
