package statemachine

import (
	"testing"
	"github.com/dipperin/dipperin-core/common"
	"github.com/stretchr/testify/assert"
	"github.com/dipperin/dipperin-core/core/csbft/model"
	"github.com/golang/mock/gomock"
	model2 "github.com/dipperin/dipperin-core/core/model"
)

func TestNewNRoundSet(t *testing.T) {
	var addresses []common.Address
	addresses = append(addresses, common.Address{1, 2, 3})
	r := NewNRoundSet(5, addresses)
	assert.NotEmpty(t, r)
	assert.NotEmpty(t, r.MissingAtRound(1))
	assert.Equal(t, false, r.EnoughAtRound(1))
	assert.Equal(t, false, r.enoughAtRound(1))
	assert.NotNil(t, r.shouldCatchUpTo())
	assert.Error(t, r.Add(&model.NewRoundMsg{}), "")
	assert.Equal(t, false, r.isCurrentVerifier(common.Address{}))
	assert.Equal(t, false, r.hasMaj32(1))
	assert.Equal(t, false, r.hasHalfUp(1))
}

func TestNewProposalSet(t *testing.T) {
	p := NewProposalSet()
	assert.NotEmpty(t, p)
	assert.Equal(t, false, p.Have(1))
	p.Proposals[1] = &model.Proposal{}
	assert.Equal(t, true, p.Have(1))
	assert.NotPanics(t, func() {
		p.Add(&model.Proposal{})
	})
	assert.Empty(t, p.GetProposal(1))
}

func TestNewBlockSet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ab := NewMockAbstractBlock(ctrl)
	b := NewBlockSet()
	assert.NotEmpty(t, b)
	assert.Nil(t, b.GetBlock(1))
	assert.NotPanics(t, func() {
		b.AddBlock(ab, 3)
	})
	ab.EXPECT().Hash().Return(common.Hash{}).AnyTimes()
	assert.NotEmpty(t, b.GetBlockByHash(common.Hash{}))
}

func TestNewVoteSet(t *testing.T) {
	var addresses []common.Address
	addresses = append(addresses, common.Address{1, 2, 3})
	v := NewVoteSet(5, addresses)
	assert.NotEmpty(t, v)
	assert.Empty(t, v.FinalVerifications(1))
	assert.NotEmpty(t, v.VotesEnough(1))
	assert.Empty(t, v.roundVotes(1))
	assert.Empty(t, v.roundBlockVotes(1))
	assert.EqualError(t, v.validVote(&model2.VoteMsg{}), "invalid vote, witness is nil")
	assert.Equal(t, false, v.isCurrentVerifier(common.Address{}))
	assert.EqualError(t, v.AddVote(&model2.VoteMsg{}), "vote height not match")
}
