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
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/model"
)

func TestValidateVotes(t *testing.T) {
	_, _, _, passChain := getTxTestEnv(t)
	passChain.slot = 1
	assert.Error(t, ValidateVotes(&BftBlockContext{
		BlockContext: BlockContext{Chain: passChain},
	})())

	a := NewAccount()
	b := NewAccount()
	va := a.getVoteMsg(1, 1, common.Hash{}, model.VoteMessage)
	vb := b.getVoteMsg(1, 1, common.Hash{}, model.VoteMessage)
	passChain.verifiers = []common.Address{a.Address(), b.Address()}
	//passChain.VerificationRoot =
	fb := fakeBlock{vs: []model.AbstractVerification{va, vb}}
	fb.SetVerifications(fb.vs)
	assert.NoError(t, ValidateVotes(&BftBlockContext{
		BlockContext: BlockContext{Chain: passChain, Block: &fb},
		Votes:        []model.AbstractVerification{va, vb},
	})())

	assert.Error(t, ValidateVotes(&BftBlockContext{
		BlockContext: BlockContext{Chain: passChain, Block: &fakeBlock{vs: []model.AbstractVerification{va}, hash: common.Hash{0x12}}},
		Votes:        []model.AbstractVerification{va},
	})())
}

func TestValidateVotesForBFT(t *testing.T) {
	assert.Error(t, ValidateVotesForBFT(&BlockContext{
		Chain: &fakeChainInterface{},
		Block: &fakeBlock{},
	})())
}

func Test_validateVotes(t *testing.T) {
	_, _, _, passChain := getTxTestEnv(t)
	passChain.slot = 1
	assert.Error(t, validateVotes(&fakeBlock{num: 1}, passChain))
	assert.Error(t, validateVotes(&fakeBlock{vRoot: common.HexToHash("0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")}, passChain))

	a := &Account{Pk: crypto.HexToECDSAErrPanic("fe10ee89565549d616d43c4e71b61d46a963fdb69489093a57cacf06836ecd91")}
	v := a.getVoteMsg(1, 1, common.Hash{}, model.VoteMessage)
	passChain.verifiers = []common.Address{a.Address()}
	passChain.block.hash = common.Hash{0x12}
	assert.Error(t, validateVotes(&fakeBlock{
		vRoot: common.HexToHash("0x231153ba3f22bb7ab9586adf88d46cf1cbeba8ae9966bcac9d6ea163144a9536"),
		vs:    []model.AbstractVerification{v},
	}, passChain))
}

func Test_validVotesForBlock(t *testing.T) {
	_, _, _, passChain := getTxTestEnv(t)
	passChain.slot = 1
	a := &Account{Pk: crypto.HexToECDSAErrPanic("fe10ee89565549d616d43c4e71b61d46a963fdb69489093a57cacf06836ecd91")}
	v := a.getVoteMsg(1, 1, common.Hash{}, model.VoteMessage)
	assert.Error(t, validVotesForBlock([]model.AbstractVerification{v, v}, &fakeBlock{}, []common.Address{a.Address()}))

	assert.Error(t, validVotesForBlock([]model.AbstractVerification{v}, &fakeBlock{isSpecial: true}, []common.Address{a.Address()}))
	v.VoteType = model.VerBootNodeVoteMessage
	assert.Error(t, validVotesForBlock([]model.AbstractVerification{v}, &fakeBlock{isSpecial: true}, []common.Address{a.Address()}))

	v2 := a.getVoteMsg(1, 1, common.Hash{}, model.AliveVerifierVoteMessage)
	assert.Error(t, validVotesForBlock([]model.AbstractVerification{v, v2}, &fakeBlock{isSpecial: true}, []common.Address{a.Address()}))

	v2 = a.getVoteMsg(1, 1, common.Hash{}, model.VoteMessage)
	assert.Error(t, validVotesForBlock([]model.AbstractVerification{v, v2}, &fakeBlock{isSpecial: true}, []common.Address{a.Address()}))

	assert.Error(t, validVotesForBlock([]model.AbstractVerification{v}, &fakeBlock{}, []common.Address{a.Address()}))

	v = a.getVoteMsg(1, 1, common.Hash{}, model.VoteMessage)
	assert.Error(t, validVotesForBlock([]model.AbstractVerification{v}, &fakeBlock{}, []common.Address{a.Address(), a.Address(), a.Address(), a.Address()}))

	assert.Error(t, validVotesForBlock([]model.AbstractVerification{v}, &fakeBlock{}, []common.Address{{}}))
}
