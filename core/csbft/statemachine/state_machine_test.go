package statemachine

import (
	"testing"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/csbft/model"
	model2 "github.com/dipperin/dipperin-core/core/model"
)

// Events Process
func TestNewBftState(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ab := NewMockAbstractBlock(ctrl)
	var b = &BftState{
		LockedBlock: ab,
	}
	assert.NotEmpty(t, b)
	assert.NotPanics(t, func() {
		b.OnNewHeight(1, 1, []common.Address{})
	})
	assert.NotPanics(t, b.OnBlockPoolNotEmpty)
	assert.NotPanics(t, func() {
		b.OnNewRound(&model.NewRoundMsg{})
	})
	assert.NotPanics(t, func() {
		b.OnNewProposal(&model.Proposal{}, nil)
	})
	assert.NotPanics(t, func() {
		b.OnPreVote(&model2.VoteMsg{
			Height: 2,
			Witness: &model2.WitMsg{
				Address: common.Address{2, 3, 4},
				Sign:    make([]byte, 65),
			},
		})
	})
	assert.NotPanics(t, func() {
		b.Height = 2
		h, av := b.OnVote(&model2.VoteMsg{
			Height: 2,
			Witness: &model2.WitMsg{
				Address: common.Address{2, 3, 4},
				Sign:    make([]byte, 65),
			},
		})
		assert.NotNil(t, h)
		assert.Nil(t, av)
	})
}

// State Change
func TestStateChange(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ab := NewMockAbstractBlock(ctrl)
	var b = &BftState{
		LockedBlock: ab,
	}
	assert.NotEmpty(t, b)
	assert.NotPanics(t, func() {
		b.enterNewHeight(1, 1, []common.Address{})
	})
	assert.NotPanics(t, func() {
		b.enterNewRound(1, 1)
	})
	assert.NotPanics(t, func() {
		b.enterPropose(1)
	})
	assert.NotPanics(t, func() {
		b.enterPreVote(&model.Proposal{}, nil)
	})
	assert.NotPanics(t, func() {
		b.enterPreCommit(1)
	})

}

//Timeout actions
func TestTimeoutActions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ab := NewMockAbstractBlock(ctrl)
	var b = &BftState{
		LockedBlock: ab,
	}
	assert.NotEmpty(t, b)
	assert.NotPanics(t, b.OnProposeTimeout)
	assert.NotPanics(t, b.OnPreVoteTimeout)
	assert.NotPanics(t, b.OnPreCommitTimeout)
}

//New Functions
func TestNewFunctions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ab := NewMockAbstractBlock(ctrl)

	var addresses []common.Address
	addresses = append(addresses, common.Address{1, 2, 3})
	var b = &BftState{
		LockedBlock:  ab,
		CurVerifiers: addresses,
	}
	assert.NotEmpty(t, b)
	assert.Nil(t, b.getNewRoundReqList())
	b.Proposal = &ProposalSet{Proposals: make(map[uint64]*model.Proposal)}
	b.Proposal.Proposals[2] = &model.Proposal{Height: 2, Round: 2}
	b.Height = 2
	assert.Equal(t, false, b.validProposal(&model.Proposal{Height: 2, Round: 2}))
	assert.NotPanics(t, b.tryEnterPropose)
	ab.EXPECT().Hash().Return(common.Hash{}).AnyTimes()
	assert.NotNil(t, b.makePrevote())
	assert.NotNil(t, b.makeVote())
	assert.NotNil(t, b.curProposer())
	assert.NotNil(t, b.proposerAtRound(1))
}
