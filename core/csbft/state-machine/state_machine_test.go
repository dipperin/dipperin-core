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

package state_machine

import (
	"github.com/dipperin/dipperin-core/common"
	model2 "github.com/dipperin/dipperin-core/core/csbft/model"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type H interface {
	Hash() common.Hash
}

func Sign(b H, signer *fakeSigner) *model.WitMsg {
	sign, err := signer.SignHash(b.Hash().Bytes())
	if err != nil {
		return nil
	}
	witness := &model.WitMsg{
		Address: signer.GetAddress(),
		Sign:    sign,
	}
	return witness
}

func MakeNewProposal(height uint64, round uint64, block model.AbstractBlock, key int) *model2.Proposal {
	sks, _ := CreateKey()
	signer := newFackSigner(sks[key])
	proposal := model2.Proposal{
		Height:    height,
		Round:     round,
		BlockID:   block.Hash(),
		Timestamp: time.Now(),
	}
	proposal.Witness = Sign(proposal, signer)
	return &proposal
}

func MakeNewProVote(height uint64, round uint64, block model.AbstractBlock, i int) *model.VoteMsg {
	sks, _ := CreateKey()
	signer := newFackSigner(sks[i])
	voteMsg := model.VoteMsg{
		Height:    height,
		Round:     round,
		BlockID:   block.Hash(),
		VoteType:  model.PreVoteMessage,
		Timestamp: time.Now(),
	}
	voteMsg.Witness = Sign(voteMsg, signer)
	return &voteMsg
}

func MakeNewVote(height uint64, round uint64, block model.AbstractBlock, i int) *model.VoteMsg {
	sks, _ := CreateKey()
	signer := newFackSigner(sks[i])
	voteMsg := model.VoteMsg{
		Height:    height,
		Round:     round,
		BlockID:   block.Hash(),
		VoteType:  model.VoteMessage,
		Timestamp: time.Now(),
	}
	voteMsg.Witness = Sign(voteMsg, signer)
	return &voteMsg
}

func MakeNewRound(height uint64, round uint64, i int) *model2.NewRoundMsg {
	sks, _ := CreateKey()
	signer := newFackSigner(sks[i])
	newRoundMsg := model2.NewRoundMsg{
		Height: height,
		Round:  round,
	}
	newRoundMsg.Witness = Sign(newRoundMsg, signer)
	return &newRoundMsg
}

func NewBftState(height uint64, round uint64, step model2.RoundStepType) BftState {
	_, v := CreateKey()
	bs := BftState{Height: height, Round: round, Step: step, BlockPoolNotEmpty: false, PreVotes: NewVoteSet(height, v), Votes: NewVoteSet(height, v), NewRound: NewNRoundSet(height, v), Proposal: NewProposalSet(), CurVerifiers: v, ProposalBlock: NewBlockSet()}
	return bs
}

func TestBftState_OnBlockPoolNotEmpty(t *testing.T) {
	height := uint64(5)
	round := uint64(4)
	step := model2.RoundStepNewHeight
	bs := NewBftState(height, round, step)
	bs.OnBlockPoolNotEmpty()
	assert.EqualValues(t, bs.BlockPoolNotEmpty, true)
	assert.EqualValues(t, bs.Step, step+1)
}

func TestBftState_OnNewProposal(t *testing.T) {
	block := model.CreateBlock(1, common.HexToHash("123"), 150)
	height := uint64(5)
	step := model2.RoundStepPropose

	round := uint64(3)
	bs := NewBftState(height, round, step)
	p := MakeNewProposal(height, round, block, 1)
	bs.OnNewProposal(p, block)
	assert.EqualValues(t, bs.Step, step)

	round = uint64(3)
	bs = NewBftState(height, round, step)
	p = MakeNewProposal(height, round, block, 0)
	bs.OnNewProposal(p, block)
	assert.EqualValues(t, bs.Step, step)

	round = uint64(4)
	bs = NewBftState(height, round, step)
	p = MakeNewProposal(height, round, block, 0)
	bs.OnNewProposal(p, block)
	assert.EqualValues(t, bs.Step, step+1)
}

func TestBftState_OnPreVote(t *testing.T) {
	block := model.CreateBlock(1, common.HexToHash("123"), 150)
	height := uint64(5)
	round := uint64(4)
	step := model2.RoundStepPropose
	bs := NewBftState(height, round, step)
	p := MakeNewProposal(height, round, block, 0)
	bs.OnNewProposal(p, block)
	assert.EqualValues(t, bs.Step, step+1)

	pv := MakeNewProVote(height, round, block, 0)
	bs.OnPreVote(pv)
	assert.EqualValues(t, bs.Step, step+1)

	pv = MakeNewProVote(height, round, block, 1)
	bs.OnPreVote(pv)
	assert.EqualValues(t, bs.Step, step+1)

	pv = MakeNewProVote(height, round, block, 2)
	bs.OnPreVote(pv)
	assert.EqualValues(t, bs.Step, step+2)

	pv = MakeNewProVote(height, round, block, 3)
	bs.OnPreVote(pv)
	assert.EqualValues(t, bs.Step, step+2)

	pv = MakeNewProVote(height, round, block, 4)
	bs.OnPreVote(pv)
	assert.EqualValues(t, bs.Step, step+2)

	pv = MakeNewProVote(height+1, round, block, 4)
	bs.OnPreVote(pv)
}

func TestBftState_OnNewRound(t *testing.T) {
	height := uint64(5)
	round := uint64(4)
	step := model2.RoundStepNewHeight
	bs := NewBftState(height, round, step)
	bs.OnBlockPoolNotEmpty()

	msg := MakeNewRound(height, round, 0)
	bs.OnNewRound(msg)
	assert.EqualValues(t, bs.Step, step+1)

	msg = MakeNewRound(height, round, 1)
	bs.OnNewRound(msg)
	assert.EqualValues(t, bs.Step, step+1)

	msg = MakeNewRound(height, round, 2)
	bs.OnNewRound(msg)
	assert.EqualValues(t, bs.Step, step+2)

	msg = MakeNewRound(height, round, 3)
	bs.OnNewRound(msg)
	assert.EqualValues(t, bs.Step, step+2)

	msg = MakeNewRound(height, round, 4)
	bs.OnNewRound(msg)
	assert.EqualValues(t, bs.Step, step+2)
}

func TestBftState_OnVote(t *testing.T) {
	block := model.CreateBlock(1, common.HexToHash("123"), 150)
	height := uint64(5)
	round := uint64(4)
	step := model2.RoundStepPreVote
	bs := NewBftState(height, round, step)

	v := MakeNewVote(height, round, block, 0)
	c, _ := bs.OnVote(v)
	assert.Equal(t, c, common.Hash{})

	v = MakeNewVote(height, round, block, 1)
	c, _ = bs.OnVote(v)
	assert.Equal(t, c, common.Hash{})

	v = MakeNewVote(height, round, block, 2)
	c, _ = bs.OnVote(v)
	assert.NotEqual(t, c, common.Hash{})

	v = MakeNewVote(height, round, block, 3)
	c, _ = bs.OnVote(v)
	assert.NotEqual(t, c, common.Hash{})

	v = MakeNewVote(height, round, block, 4)
	c, _ = bs.OnVote(v)
	assert.Equal(t, c, common.Hash{})
}

func TestBftState_OnPreCommitTimeout(t *testing.T) {
	height := uint64(5)
	round := uint64(4)
	step := model2.RoundStepPropose
	bs := NewBftState(height, round, step)
	bs.OnPreCommitTimeout()
	assert.EqualValues(t, bs.Step, model2.RoundStepNewRound)
}

func TestBftState_OnNewHeight(t *testing.T) {
	bs := NewBftState(uint64(5), uint64(4), model2.RoundStepPropose)
	bs.OnNewHeight(uint64(6), uint64(7), []common.Address{})

	assert.EqualValues(t, uint64(6), bs.Height)
	assert.EqualValues(t, uint64(7), bs.Round)
}

func TestBftState_enterNewRound(t *testing.T) {
	bs := NewBftState(uint64(5), uint64(4), model2.RoundStepPropose)
	bs.enterNewRound(uint64(7), uint64(7))

	assert.Equal(t, uint64(5), bs.Height)
}

func TestBftState_enterPrevote(t *testing.T) {
	bs := NewBftState(uint64(5), uint64(4), model2.RoundStepPropose)

	block := model.CreateBlock(5, common.HexToHash("123"), 150)
	p := MakeNewProposal(uint64(6), uint64(4), block, 1)
	bs.enterPreVote(p, block)
}

func TestBftState_enterPrecommit(t *testing.T) {
	bs := NewBftState(uint64(5), uint64(4), model2.RoundStepPropose)
	bs.enterPreCommit(uint64(4))
	assert.Equal(t, model2.RoundStepPropose, bs.Step)

	bs = NewBftState(uint64(5), uint64(4), model2.RoundStepPreVote)
	bs.enterPreCommit(uint64(6))
	assert.Equal(t, model2.RoundStepPreVote, bs.Step)
}

func TestBftState_validProposal(t *testing.T) {
	bs := NewBftState(uint64(5), uint64(4), model2.RoundStepPropose)
	block := model.CreateBlock(5, common.HexToHash("123"), 150)

	p := MakeNewProposal(uint64(6), uint64(4), block, 1)
	assert.Equal(t, false, bs.validProposal(p))

	p = MakeNewProposal(uint64(5), uint64(4), block, 3)
	assert.Equal(t, false, bs.validProposal(p))

	p = MakeNewProposal(uint64(5), uint64(4), block, 0)
	p.Witness.Sign = []byte{}
	assert.Equal(t, false, bs.validProposal(p))

	p = MakeNewProposal(uint64(5), uint64(4), block, 0)
	assert.Equal(t, true, bs.validProposal(p))

	bs.Proposal.Add(p)
	assert.Equal(t, false, bs.validProposal(p))
}
