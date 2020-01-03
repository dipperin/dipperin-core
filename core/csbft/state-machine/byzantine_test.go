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
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	model2 "github.com/dipperin/dipperin-core/core/csbft/model"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestByzantine(t *testing.T) {
	fakeblock := &FakeBlock{uint64(1), common.HexToHash("0x232"), nil}
	fakeblock2 := &FakeBlock{uint64(1), common.HexToHash("0x2d32"), nil}
	sh0 := NewFakeStateHandle(0)
	sh1 := NewFakeStateHandle(1)
	sh2 := NewFakeStateHandle(2)
	sh3 := NewFakeStateHandle(3)

	sh0.blockPool.AddBlock(fakeblock)
	sh1.blockPool.AddBlock(fakeblock)
	sh2.blockPool.AddBlock(fakeblock)
	sh3.blockPool.AddBlock(fakeblock)

	//Confirm round number
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, model2.RoundStepNewRound, sh0.bs.Step)
	assert.Equal(t, model2.RoundStepNewRound, sh1.bs.Step)
	assert.Equal(t, model2.RoundStepNewRound, sh2.bs.Step)
	assert.Equal(t, model2.RoundStepNewRound, sh3.bs.Step)

	sh0.NewRound(MakeNewRound(1, 1, 1))
	sh0.NewRound(MakeNewRound(1, 1, 2))
	sh0.NewRound(MakeNewRound(1, 1, 3))
	sh1.NewRound(MakeNewRound(1, 1, 0))
	sh1.NewRound(MakeNewRound(1, 1, 2))
	sh1.NewRound(MakeNewRound(1, 1, 3))
	sh2.NewRound(MakeNewRound(1, 1, 0))
	sh2.NewRound(MakeNewRound(1, 1, 1))
	sh2.NewRound(MakeNewRound(1, 1, 3))
	sh3.NewRound(MakeNewRound(1, 1, 0))
	sh3.NewRound(MakeNewRound(1, 1, 1))
	sh3.NewRound(MakeNewRound(1, 1, 2))

	//Sh1 is proposed
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, model2.RoundStepPropose, sh0.bs.Step)
	assert.Equal(t, model2.RoundStepPreVote, sh1.bs.Step)
	assert.Equal(t, model2.RoundStepPropose, sh2.bs.Step)
	assert.Equal(t, model2.RoundStepPropose, sh3.bs.Step)

	fmt.Println("**********", "fb", fakeblock.Hash())

	//Feed the proposal
	msg := MakeNewProposal(1, 1, fakeblock, 1)
	sh2.NewProposal(msg)
	sh3.NewProposal(msg)
	time.Sleep(10 * time.Millisecond)
	assert.Equal(t, model2.RoundStepPreVote, sh3.bs.Step)

	//Sh0 is the byzantine node.
	//Sh1,2 received vote from 1,2
	//Sh3 received vote from 0,2,3

	sh1.PreVote(MakeNewProVote(1, 1, fakeblock, 2))
	sh2.PreVote(MakeNewProVote(1, 1, fakeblock, 1))
	sh3.PreVote(MakeNewProVote(1, 1, fakeblock, 0))
	sh3.PreVote(MakeNewProVote(1, 1, fakeblock, 2))

	//Sh3 locked at fackblock
	//Sh1,Sh2 have no lock
	time.Sleep(20 * time.Millisecond)
	assert.Equal(t, model2.RoundStepPreVote, sh1.bs.Step)
	assert.Equal(t, model2.RoundStepPreVote, sh2.bs.Step)
	assert.Equal(t, model2.RoundStepPreCommit, sh3.bs.Step)
	assert.Equal(t, sh1.bs.LockedBlock, nil)
	assert.Equal(t, sh2.bs.LockedBlock, nil)
	assert.Equal(t, sh3.bs.LockedBlock.Hash().Hex(), fakeblock.Hash().Hex())

	sh1.blockPool.RemoveBlock(fakeblock.Hash())
	sh2.blockPool.RemoveBlock(fakeblock.Hash())
	sh3.blockPool.RemoveBlock(fakeblock.Hash())
	sh1.blockPool.AddBlock(fakeblock2)
	sh2.blockPool.AddBlock(fakeblock2)
	sh3.blockPool.AddBlock(fakeblock2)

	time.Sleep(300 * time.Millisecond)
	assert.Equal(t, model2.RoundStepNewRound, sh1.bs.Step)
	assert.Equal(t, model2.RoundStepNewRound, sh2.bs.Step)
	assert.Equal(t, model2.RoundStepNewRound, sh3.bs.Step)
	//enter a new round
	sh1.NewRound(MakeNewRound(1, 2, 0))
	sh1.NewRound(MakeNewRound(1, 2, 2))
	sh1.NewRound(MakeNewRound(1, 2, 3))
	sh2.NewRound(MakeNewRound(1, 2, 0))
	sh2.NewRound(MakeNewRound(1, 2, 1))
	sh2.NewRound(MakeNewRound(1, 2, 3))
	sh3.NewRound(MakeNewRound(1, 2, 0))
	sh3.NewRound(MakeNewRound(1, 2, 1))
	sh3.NewRound(MakeNewRound(1, 2, 2))

	//Sh2 is proposed
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, uint64(2), sh1.bs.Round)
	assert.Equal(t, uint64(2), sh2.bs.Round)
	assert.Equal(t, uint64(2), sh3.bs.Round)
	assert.Equal(t, model2.RoundStepPropose, sh1.bs.Step)
	assert.Equal(t, model2.RoundStepPreVote, sh2.bs.Step)
	assert.Equal(t, model2.RoundStepPropose, sh3.bs.Step)

	proposal2 := MakeNewProposal(1, 2, fakeblock2, 2)
	sh1.NewProposal(proposal2)
	sh3.NewProposal(proposal2)

	time.Sleep(60 * time.Millisecond)
	assert.Equal(t, model2.RoundStepPreVote, sh1.bs.Step)
	assert.Equal(t, model2.RoundStepPreVote, sh2.bs.Step)
	assert.Equal(t, model2.RoundStepPreVote, sh3.bs.Step)

	//Sh0 make a conflict vote, vote on fakeBlock2
	//Sh3 locked on fakeblock, would not do vote
	sh1.PreVote(MakeNewProVote(1, 2, fakeblock2, 0))
	sh1.PreVote(MakeNewProVote(1, 2, fakeblock2, 2))
	sh2.PreVote(MakeNewProVote(1, 2, fakeblock2, 0))
	sh2.PreVote(MakeNewProVote(1, 2, fakeblock2, 1))
	sh3.PreVote(MakeNewProVote(1, 2, fakeblock2, 0))
	sh3.PreVote(MakeNewProVote(1, 2, fakeblock2, 1))
	sh3.PreVote(MakeNewProVote(1, 2, fakeblock2, 2))

	//Sh3 clear lock, and then locked on fakeBlock2
	time.Sleep(20 * time.Millisecond)
	assert.Equal(t, uint64(2), sh1.bs.Round)
	assert.Equal(t, uint64(2), sh2.bs.Round)
	assert.Equal(t, uint64(2), sh3.bs.Round)
	assert.Equal(t, fakeblock2, sh1.bs.LockedBlock)
	assert.Equal(t, fakeblock2, sh2.bs.LockedBlock)
	assert.Equal(t, fakeblock2, sh3.bs.LockedBlock)
	assert.Equal(t, model2.RoundStepPreCommit, sh1.bs.Step)
	assert.Equal(t, model2.RoundStepPreCommit, sh2.bs.Step)
	assert.Equal(t, model2.RoundStepPreCommit, sh3.bs.Step)
}
