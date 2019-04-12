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
	model2 "github.com/dipperin/dipperin-core/core/csbft/model"
	"github.com/stretchr/testify/assert"
	"testing"
	"github.com/dipperin/dipperin-core/common"
	"strconv"
	"github.com/dipperin/dipperin-core/core/model"
	"time"
	"reflect"
	"github.com/dipperin/dipperin-core/core/csbft/components"
)

func makeValidBlock(height uint64, round uint64) (fakeblock *FakeBlock, commits []model.AbstractVerification){
	fakeblock = &FakeBlock{height, common.HexToHash(strconv.Itoa(int(height))), nil}
	commits = append(commits,MakeNewVote(height, round, fakeblock, 0))
	commits = append(commits,MakeNewVote(height, round, fakeblock, 1))
	commits = append(commits,MakeNewVote(height, round, fakeblock, 2))
	return
}

func TestStateHandler_OnNewHeight(t *testing.T) {
	sh0 := NewFakeStateHandle(0)
	sh0.ChainReader.SaveBlock(makeValidBlock(uint64(1),uint64(1)))
	sh0.OnNewHeight(2)

	assert.Equal(t, uint64(2),sh0.bs.Height)
	assert.Equal(t, model2.RoundStepNewHeight, sh0.bs.Step)

	sh0.OnNewHeight(3)
	assert.Equal(t, uint64(2), sh0.bs.Height)

	sh0.ChainReader.SaveBlock(makeValidBlock(uint64(9),uint64(1)))
	sh0.OnNewHeight(10)
	assert.Equal(t, uint64(0), sh0.bs.Round)

	block,_ :=makeValidBlock(uint64(10),uint64(1))
	sh0.ChainReader.SaveBlock(block,nil)
	sh0.OnNewHeight(11)
}

func TestStateHandler_OnReset(t *testing.T) {
	sh0 := NewFakeStateHandle(0)
	assert.NotNil(t,sh0.Reset())
}

func TestStateHandler_OnBlockPoolNotEmpty(t *testing.T) {
	sh0 := NewFakeStateHandle(0)
	sh0.OnBlockPoolNotEmpty()

	assert.Equal(t,model2.RoundStepNewRound,sh0.bs.Step)
}

func TestStateHandler_GetProposalBlock(t *testing.T) {
	sh0 := NewFakeStateHandle(0)
	nothing := sh0.GetProposalBlock(common.HexToHash("1"))

	assert.Equal(t,nil,nothing)

	block := &FakeBlock{uint64(1), common.HexToHash("0x123"), nil}
	sh0.bs.OnNewProposal(MakeNewProposal(1, 1, block , 1),block)
	oneBlock := sh0.GetProposalBlock(common.HexToHash("0x123"))

	assert.NotNil(t,oneBlock)
	assert.Equal(t,uint64(1),oneBlock.Number())

	sh0.Stop()
	stoped := sh0.GetProposalBlock(common.HexToHash("1"))
	assert.Equal(t,nil,stoped)
}

func TestStateHandler_GetRoundMsg(t *testing.T) {
	sh0 := NewFakeStateHandle(0)
	sh0.bs.Round = 6

	low := sh0.GetRoundMsg(uint64(1),uint64(3))
	assert.NotNil(t,low)

	equal := sh0.GetRoundMsg(uint64(1),uint64(6))
	assert.NotNil(t,equal)

	high := sh0.GetRoundMsg(uint64(1),uint64(7))
	assert.Equal(t, (*model2.NewRoundMsg)(nil), high)

	none := sh0.GetRoundMsg(uint64(8),uint64(0))
	assert.Equal(t, (*model2.NewRoundMsg)(nil), none)
}

func TestStateHandler_EnterNewHeight(t *testing.T) {
	sh0 := NewFakeStateHandle(0)

	assert.Equal(t, sh0.bs.Height, uint64(1))
	assert.Equal(t, model2.RoundStepNewHeight, sh0.bs.Step)

	//Propose Stage
	fakeBlock := &FakeBlock{uint64(1), common.HexToHash("0x232"), nil}
	sh0.blockPool.AddBlock(fakeBlock)

	//New Round-> Propose -> Prevote
	sh0.NewRound(MakeNewRound(1, 1, 1))
	sh0.NewRound(MakeNewRound(1, 1, 2))
	sh0.NewRound(MakeNewRound(1, 1, 3))
	sh0.NewProposal(MakeNewProposal(1, 1, fakeBlock, 1))
	time.Sleep(100 * time.Millisecond)

	//Prevote stage
	sh0.PreVote(MakeNewProVote( 1, 1, fakeBlock, 1))
	sh0.PreVote(MakeNewProVote( 1, 1, fakeBlock, 2))
	sh0.PreVote(MakeNewProVote( 1, 1, fakeBlock, 3))

	//Precommit
	time.Sleep(100 * time.Millisecond)
	sh0.Vote(MakeNewVote(1, 1, fakeBlock, 1))
	sh0.Vote(MakeNewVote(1, 1, fakeBlock, 2))
	sh0.Vote(MakeNewVote(1, 1, fakeBlock, 3))

	time.Sleep(120 * time.Millisecond)
	assert.Equal(t, uint64(2), sh0.bs.Height)
}


func TestStateHandler_LockBlock(t *testing.T) {
	sh0 := NewFakeStateHandle(0)

	assert.Equal(t, sh0.bs.Height, uint64(1))
	assert.Equal(t, model2.RoundStepNewHeight, sh0.bs.Step)

	//Propose Stage
	fakeBlock := &FakeBlock{uint64(1), common.HexToHash("0x232"), nil}
	sh0.blockPool.AddBlock(fakeBlock)

	//New Round-> Propose -> Prevote
	sh0.NewRound( MakeNewRound(1,1,1))
	sh0.NewRound( MakeNewRound(1,1,2))
	sh0.NewRound( MakeNewRound(1,1,3))

	sh0.NewProposal(MakeNewProposal(1, 1, fakeBlock, 1))
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, model2.RoundStepPreVote, sh0.bs.Step)

	//Prevote stage
	sh0.PreVote(MakeNewProVote(1, 1, fakeBlock, 1))
	sh0.PreVote(MakeNewProVote(1, 1, fakeBlock, 2))
	sh0.PreVote(MakeNewProVote(1, 1, fakeBlock, 3))

	//Block is locked
	assert.Equal(t, model2.RoundStepPreCommit, sh0.bs.Step)
	assert.Equal(t, fakeBlock.Hash(),sh0.bs.LockedBlock.Hash())
	assert.Equal(t,uint64(1),sh0.bs.LockedRound)

	//Precommit time out
	time.Sleep(300 * time.Millisecond)
	assert.Equal(t, model2.RoundStepNewRound, sh0.bs.Step)
	assert.Equal(t, true, sh0.bs.LockedBlock.Hash().IsEqual(fakeBlock.Hash()))

	//New round
	sh0.NewRound( MakeNewRound(1,2,1))
	sh0.NewRound( MakeNewRound(1,2,2))
	sh0.NewRound( MakeNewRound(1,2,3))

	// Propose block 2
	fakeBlock2 := &FakeBlock{uint64(1), common.HexToHash("0x9999"), nil}
	sh0.blockPool.AddBlock(fakeBlock2)
	sh0.NewProposal(MakeNewProposal(1, 2, fakeBlock2, 2))
	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, model2.RoundStepPreVote, sh0.bs.Step)
	assert.Equal(t, fakeBlock.Hash(),sh0.bs.LockedBlock.Hash())
	assert.Equal(t,uint64(1),sh0.bs.LockedRound)

	//Prevote stage
	sh0.PreVote(MakeNewProVote(1, 2, fakeBlock2, 1))
	sh0.PreVote(MakeNewProVote(1, 2, fakeBlock2, 2))
	sh0.PreVote(MakeNewProVote(1, 2, fakeBlock2, 3))

	//Unlocked and lock on the new block
	time.Sleep(20 * time.Millisecond)
	assert.Equal(t, fakeBlock2.Hash(),sh0.bs.LockedBlock.Hash())
	assert.Equal(t, uint64(2), sh0.bs.LockedRound)
}

func TestStateHandler_ReceiveProposeThenRoundConfirm(t *testing.T) {
	sh0 := NewFakeStateHandle(0)
	sh0.ChainReader.SaveBlock(makeValidBlock(uint64(1),uint64(1)))
	sh0.OnNewHeight(2)

	assert.Equal(t, model2.RoundStepNewHeight, sh0.bs.Step)
	assert.Equal(t, uint64(2), sh0.bs.Height)
	assert.Equal(t, uint64(2), sh0.bs.Round)

	//Receive propose from node 2, before get 32 majority round confirmation
	fakeblock2 := &FakeBlock{uint64(2), common.HexToHash("0x232"), nil}
	sh0.blockPool.AddBlock(fakeblock2)
	sh0.NewProposal(MakeNewProposal(2, 2, fakeblock2, 2))
	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, model2.RoundStepNewRound, sh0.bs.Step)
	sh0.NewRound( MakeNewRound(2, 2, 1))
	sh0.NewRound( MakeNewRound(2, 2, 2))
	sh0.NewRound( MakeNewRound(2, 2, 3))

	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, model2.RoundStepPreVote, sh0.bs.Step)
}

func TestStateHandler_GetVotesThenGetProposal(t *testing.T) {
	sh0 := NewFakeStateHandle(0)
	sh0.ChainReader.SaveBlock(makeValidBlock(uint64(1),uint64(1)))
	sh0.OnNewHeight(2)

	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, model2.RoundStepNewHeight, sh0.bs.Step)
	assert.Equal(t, uint64(2), sh0.bs.Height)
	assert.Equal(t, uint64(2), sh0.bs.Round)

	//---------------------------------------------
	sh0.NewRound( MakeNewRound(2, 2, 1))
	sh0.NewRound( MakeNewRound(2, 2, 2))
	sh0.NewRound( MakeNewRound(2, 2, 3))
	//sh0.poolNotEmptyChan <- struct{}{}
	time.Sleep(100 * time.Millisecond)

	//assert.Equal(t, RoundStepNewRound, sh0.bs.Step)
	fakeblock2 := &FakeBlock{uint64(2), common.HexToHash("0x232"), nil}
	sh0.PreVote(MakeNewProVote(2, 2, fakeblock2, 1))
	sh0.PreVote(MakeNewProVote(2, 2, fakeblock2, 2))
	sh0.PreVote(MakeNewProVote(2, 2, fakeblock2, 3))
	time.Sleep(100 * time.Millisecond)

	//Receive propose from node 2, before get 32 majority round confirmation
	sh0.Fetcher = &FakeFetcher{fakeblock2}
	sh0.blockPool.AddBlock(fakeblock2)
	sh0.NewProposal(MakeNewProposal(2,2,fakeblock2, 2))
	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, model2.RoundStepPreCommit, sh0.bs.Step)
}

//When node is in propose stage, get a new height message
func TestInterupted_when_propose(t *testing.T) {
	sh0 := NewFakeStateHandle(0)
	sh0.ChainReader.SaveBlock(makeValidBlock(uint64(1),uint64(1)))
	sh0.OnNewHeight(2)

	time.Sleep(300 * time.Millisecond)
	assert.Equal(t, model2.RoundStepNewHeight, sh0.bs.Step)
	assert.Equal(t, uint64(2), sh0.bs.Height)
	assert.Equal(t, uint64(2), sh0.bs.Round)

	//---------------------------------------------
	sh0.blockPool.AddBlock(&FakeBlock{uint64(2), common.HexToHash("0x232"), nil})

	go func() {
		sh0.NewRound(MakeNewRound(2,1,1))
		sh0.NewRound(MakeNewRound(2,1,2))
		sh0.NewRound(MakeNewRound(2,1,3))
	}()
	go func() {
		var block3 model.AbstractBlock
		block3 = &FakeBlock{uint64(2), common.HexToHash("0xaa32"), nil}
		var v []model.AbstractVerification
		v = append(v, MakeNewVote(2,1,block3,1))
		sh0.ChainReader.SaveBlock(block3, v)
		sh0.NewHeight(3)
	}()
	time.Sleep(300 * time.Millisecond)
	assert.Equal(t, uint64(3), sh0.bs.Height)
	// assert.Equal(t, RoundStepNewHeight, sh0.bs.Step)
}

func TestInterupted_when_prevote(t *testing.T) {
	sh0 := NewFakeStateHandle(0)
	sh0.ChainReader.SaveBlock(makeValidBlock(uint64(1),uint64(1)))
	sh0.OnNewHeight(2)

	time.Sleep(600 * time.Millisecond)
	assert.Equal(t, model2.RoundStepNewHeight, sh0.bs.Step)
	assert.Equal(t, uint64(2), sh0.bs.Height)
	assert.Equal(t, uint64(2), sh0.bs.Round)

	//---------------------------------------------
	bk := &FakeBlock{uint64(2), common.HexToHash("0x232"), nil}
	sh0.blockPool.AddBlock(bk)
	
	sh0.NewRound(MakeNewRound(2,2,1))
	sh0.NewRound(MakeNewRound(2,2,2))
	sh0.NewRound(MakeNewRound(2,2,3))

	proposal := MakeNewProposal(2, 2, bk, 2)
	sh0.NewProposal(proposal)
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, model2.RoundStepPreVote, sh0.bs.Step)
	go func() {
		sh0.Vote(MakeNewProVote(2,2,bk,1))
		sh0.Vote(MakeNewProVote(2,2,bk,2))
		sh0.Vote(MakeNewProVote(2,2,bk,3))
	}()
	go func() {
		var block3 model.AbstractBlock
		block3 = &FakeBlock{uint64(2), common.HexToHash("0xaa32"), nil}
		var v []model.AbstractVerification
		v = append(v, MakeNewVote(2,2,block3,1))
		sh0.ChainReader.SaveBlock(block3, v)
		sh0.NewHeight(3)
	}()
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, model2.RoundStepNewHeight, sh0.bs.Step)
}

func TestInterupted_when_precommit(t *testing.T) {
	sh0 := NewFakeStateHandle(0)
	sh0.ChainReader.SaveBlock(makeValidBlock(uint64(1),uint64(1)))
	sh0.OnNewHeight(2)

	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, model2.RoundStepNewHeight, sh0.bs.Step)
	assert.Equal(t, uint64(2), sh0.bs.Height)
	assert.Equal(t, uint64(2), sh0.bs.Round)

	//---------------------------------------------
	bk := &FakeBlock{uint64(2), common.HexToHash("0x232"), nil}
	sh0.blockPool.AddBlock(bk)
	sh0.NewRound(MakeNewRound(2,2,1))
	sh0.NewRound(MakeNewRound(2,2,2))
	sh0.NewRound(MakeNewRound(2,2,3))

	proposal := MakeNewProposal(2, 2, bk, 2)
	sh0.NewProposal(proposal)
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, model2.RoundStepPreVote, sh0.bs.Step)
	

	sh0.PreVote(MakeNewProVote(2,2,bk,1))
	sh0.PreVote(MakeNewProVote(2,2,bk,2))
	sh0.PreVote(MakeNewProVote(2,2,bk,3))
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, model2.RoundStepPreCommit, sh0.bs.Step)

	go func() {
		sh0.Vote(MakeNewVote(2,2,bk,1))
		sh0.Vote(MakeNewVote(2,2,bk,2))
		sh0.Vote(MakeNewVote(2,2,bk,3))
	}()
	go func() {
		var block3 model.AbstractBlock
		block3 = &FakeBlock{uint64(2), common.HexToHash("0xaa32"), nil}
		var v []model.AbstractVerification
		v = append(v, MakeNewVote(2,2,block3,1))
		sh0.ChainReader.SaveBlock(block3, v)
		sh0.NewHeight(3)
	}()
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, model2.RoundStepNewHeight, sh0.bs.Step)
}

func TestStateHandler_SetFetcher(t *testing.T) {
	sh0 := NewFakeStateHandle(0)
	fetcher := FakeFetcher{}
	sh0.SetFetcher(&fetcher)
	assert.Equal(t,reflect.ValueOf(&fetcher).Pointer(),reflect.ValueOf(sh0.Fetcher).Pointer())
}

func TestStateHandler_OnNewProposal(t *testing.T) {
	sh0 := NewFakeStateHandle(0)
	bk := &FakeBlock{uint64(1), common.HexToHash("0x232"), nil}
	proposal := MakeNewProposal(1, 1, bk, 1)
	sh0.OnNewProposal(proposal)

	sh0.blockPool.AddBlock(bk)
	sh0.Validator = &FakeValidtor{NotValid:true}
	sh0.OnNewProposal(proposal)

	assert.Equal(t,0,len(sh0.bs.Proposal.Proposals))
}

func TestStateHandler_OnTimeout(t *testing.T) {
	sh0 := NewFakeStateHandle(0)
	assert.Panics(t, func() {
		sh0.OnTimeout(components.TimeoutInfo{Height:uint64(1),Round:uint64(1),Step:model2.RoundStepNewHeight})
	})

	sh0.OnTimeout(components.TimeoutInfo{Height:uint64(1),Round:uint64(1),Step:model2.RoundStepPropose})
	sh0.OnTimeout(components.TimeoutInfo{Height:uint64(1),Round:uint64(1),Step:model2.RoundStepPreVote})
	sh0.OnTimeout(components.TimeoutInfo{Height:uint64(1),Round:uint64(1),Step:model2.RoundStepPreCommit})
}

func TestStateHandler_OnNewRound(t *testing.T) {
	sh0 := NewFakeStateHandle(0)
	bk := &FakeBlock{uint64(1), common.HexToHash("0x232"), nil}
	sh0.blockPool.AddBlock(bk)

	sh0.NewRound(MakeNewRound(1,2,1))
	sh0.NewRound(MakeNewRound(1,2,2))
	sh0.NewRound(MakeNewRound(1,2,3))

	time.Sleep(10 * time.Millisecond)
	assert.Equal(t,uint64(2),sh0.bs.Round)
}

