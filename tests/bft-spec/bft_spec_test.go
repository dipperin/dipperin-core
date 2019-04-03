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


package bft_spec

import (
	"testing"
	"github.com/dipperin/dipperin-core/tests"
	"github.com/dipperin/dipperin-core/third-party/log/pbft_log"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/dipperin/dipperin-core/core/csbft/state-machine"
	"github.com/stretchr/testify/assert"
	"github.com/dipperin/dipperin-core/core/csbft/model"
	"github.com/dipperin/dipperin-core/common"
	model2 "github.com/dipperin/dipperin-core/core/model"
)

func TestNormalBft(t *testing.T) {
	vCount := 22
	pbft_log.InitPbftLogger(log.LvlDebug, "bft_spec", true)
	cluster := tests.NewBftCluster(tests.AccFactory.GenAccounts(vCount))

	nrs := cluster.NewRoundMsg(vCount, 1, 1)
	cluster.StatesIter(func(state *state_machine.BftState) {
		assert.Equal(t, model.RoundStepNewRound, state.Step)
		for _, nr := range nrs {
			state.OnNewRound(nr)
		}
		assert.Equal(t, model.RoundStepPropose, state.Step)
	})

	block := &tests.FakeBlockForBft{ Num: 1, PHash: common.HexToHash("0x123") }
	ps := cluster.NewProposal(vCount, 1, block)
	cluster.StatesIter(func(state *state_machine.BftState) {
		for _, p := range ps {
			state.OnNewProposal(p, block)
		}
		assert.Equal(t, model.RoundStepPreVote, state.Step)
	})

	vs := cluster.NewVote(vCount, 1, model2.PreVoteMessage, block)
	cluster.StatesIter(func(state *state_machine.BftState) {
		for _, v := range vs {
			state.OnPreVote(v)
		}
		assert.Equal(t, model.RoundStepPreCommit, state.Step)
	})

	vs2 := cluster.NewVote(vCount, 1, model2.VoteMessage, block)
	vs2Len := len(vs2)
	cluster.StatesIter(func(state *state_machine.BftState) {
		for i, v := range vs2 {
			rBlock, rVers := state.OnVote(v)
			if i == vs2Len - 1 {

				// Get the result of the vote
				assert.Equal(t, block.Hash(), rBlock)
				assert.Len(t, rVers, vCount)
			}
		}
	})
}
