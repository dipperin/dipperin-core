package statemachine

import (
	"testing"
	"github.com/golang/mock/gomock"
	"github.com/dipperin/dipperin-core/core/csbft/components"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/stretchr/testify/assert"
	model2 "github.com/dipperin/dipperin-core/core/csbft/model"
)

type myStr struct {
	state        *StateHandler
	chain        *MockChainReader
	fetcher      *MockFetcher
	signer       *MockMsgSigner
	sender       *MockMsgSender
	validator    *MockValidator
	block        *MockAbstractBlock
	verification *MockAbstractVerification
}

func testState(ctrl *gomock.Controller) (my myStr) {
	bp := components.NewBlockPool(0, nil)
	chain := NewMockChainReader(ctrl)
	bp.SetNodeConfig(chain)
	fetcher := NewMockFetcher(ctrl)
	signer := NewMockMsgSigner(ctrl)
	sender := NewMockMsgSender(ctrl)
	validator := NewMockValidator(ctrl)
	verification := NewMockAbstractVerification(ctrl)
	var config = &BftConfig{
		ChainReader: chain,
		Fetcher:     fetcher,
		Signer:      signer,
		Sender:      sender,
		Validator:   validator,
	}
	var height uint64
	height = 1
	block := NewMockAbstractBlock(ctrl)
	block.EXPECT().Number().Return(height).AnyTimes()
	chain.EXPECT().CurrentBlock().Return(block).AnyTimes()
	chain.EXPECT().GetSeenCommit(height).Return([]model.AbstractVerification{verification}).AnyTimes()
	verification.EXPECT().GetRound().Return(uint64(2333)).AnyTimes()
	chain.EXPECT().IsChangePoint(gomock.Any(), gomock.Any()).Return(true).AnyTimes()
	chain.EXPECT().GetNextVerifiers().Return([]common.Address{}).AnyTimes()
	NewStateHandler(config, DefaultConfig, bp)
	my.state = NewStateHandler(config, DefaultConfig, bp)
	my.fetcher = fetcher
	my.signer = signer
	my.sender = sender
	my.validator = validator
	my.block = block
	my.verification = verification
	return
}

func TestNewStateHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	if s := testState(ctrl); s.state == nil {
		t.Error("fail to NewStateHandler")
	}
}

func TestStateHandler_OnStart(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s := testState(ctrl)
	if s.state == nil {
		t.Error("fail to NewStateHandler")
	}
	assert.NoError(t, s.state.OnStart())
}

func TestStateHandler_OnStop(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s := testState(ctrl)
	if s.state == nil {
		t.Error("fail to NewStateHandler")
	}
	assert.NotPanics(t, s.state.OnStop)
}

func TestStateHandler_OnReset(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s := testState(ctrl)
	if s.state == nil {
		t.Error("fail to NewStateHandler")
	}
	assert.NoError(t, s.state.OnReset())
}

func TestStateHandler_loop(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s := testState(ctrl)
	if s.state == nil {
		t.Error("fail to NewStateHandler")
	}
	assert.NotPanics(t, func() {
		go s.state.loop()
	})
}

func TestStateHandler_OnNewHeight(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s := testState(ctrl)
	if s.state == nil {
		t.Error("fail to NewStateHandler")
	}
	assert.NotPanics(t, func() {
		s.state.OnNewHeight(2)
	})
}

func TestStateHandler_OnNewRound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s := testState(ctrl)
	if s.state == nil {
		t.Error("fail to NewStateHandler")
	}
	assert.NotPanics(t, func() {
		s.state.OnNewRound(&model2.NewRoundMsg{
			Height: 2,
			Round:  2,
			Witness: &model.WitMsg{
				Address: common.Address{2, 3, 4},
				Sign:    []byte("2314"),
			},
		})
	})
}

func TestStateHandler_OnBlockPoolNotEmpty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s := testState(ctrl)
	if s.state == nil {
		t.Error("fail to NewStateHandler")
	}
	//var msg = &model2.NewRoundMsg{
	//	Height: 2,
	//	Round:  2,
	//	Witness: &model.WitMsg{
	//		Address: common.Address{2, 3, 4},
	//		Sign:    []byte("2222"),
	//	},
	//}
	s.signer.EXPECT().SignHash(gomock.Any()).Return([]byte("2222"), nil).AnyTimes()
	s.signer.EXPECT().GetAddress().Return(common.Address{2, 3, 4}).AnyTimes()
	s.sender.EXPECT().BroadcastMsg(uint64(model2.TypeOfNewRoundMsg), gomock.Any()).AnyTimes()
	assert.NotPanics(t, s.state.OnBlockPoolNotEmpty)
}

func TestStateHandler_OnNewProposal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s := testState(ctrl)
	if s.state == nil {
		t.Error("fail to NewStateHandler")
	}
	assert.NotPanics(t, func() {
		s.state.OnNewProposal(&model2.Proposal{})
	})
	//s.state.OnNewProposal(&model2.Proposal{})
}

func TestStateHandler_OnPreVote(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s := testState(ctrl)
	if s.state == nil {
		t.Error("fail to NewStateHandler")
	}
	assert.NotPanics(t, func() {
		s.state.OnPreVote(&model.VoteMsg{})
	})
}

func TestStateHandler_OnVote(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s := testState(ctrl)
	if s.state == nil {
		t.Error("fail to NewStateHandler")
	}
	assert.NotPanics(t, func() {
		s.state.bs.Height = 2
		s.state.bs.Votes.verifiers = append(s.state.bs.Votes.verifiers, common.Address{2, 3, 4})
		s.state.OnVote(&model.VoteMsg{
			Height: 2,
			Witness: &model.WitMsg{
				Address: common.Address{2, 3, 4},
				Sign:    make([]byte, 65),
			},
		})
	})
}

func TestStateHandler_OnTimeout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s := testState(ctrl)
	if s.state == nil {
		t.Error("fail to NewStateHandler")
	}
	assert.NotPanics(t, func() {
		s.state.OnTimeout(components.TimeoutInfo{})
	})
}

func TestStateHandler_broadcastNewRoundMsg(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s := testState(ctrl)
	if s.state == nil {
		t.Error("fail to NewStateHandler")
	}
	s.signer.EXPECT().SignHash(gomock.Any()).Return([]byte("2222"), nil).AnyTimes()
	s.signer.EXPECT().GetAddress().Return(common.Address{2, 3, 4}).AnyTimes()
	s.sender.EXPECT().BroadcastMsg(uint64(model2.TypeOfNewRoundMsg), gomock.Any()).AnyTimes()
	assert.NotPanics(t, s.state.broadcastNewRoundMsg)
}

func TestStateHandler_curProposer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s := testState(ctrl)
	if s.state == nil {
		t.Error("fail to NewStateHandler")
	}
	assert.NotNil(t, s.state.curProposer())
}

func TestStateHandler_broadcastReqRoundMsg(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s := testState(ctrl)
	if s.state == nil {
		t.Error("fail to NewStateHandler")
	}
	s.sender.EXPECT().SendReqRoundMsg(uint64(model2.TypeOfReqNewRoundMsg), []common.Address{}, gomock.Any()).AnyTimes()
	assert.NotPanics(t, func() {
		s.state.broadcastReqRoundMsg([]common.Address{})
	})
}

func TestStateHandler_fetchProposalBlock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s := testState(ctrl)
	if s.state == nil {
		t.Error("fail to NewStateHandler")
	}
	assert.NotPanics(t, func() {
		go s.state.fetchProposalBlock(common.Hash{}, common.Address{})
	})
}

func TestStateHandler_onEnterNewRound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s := testState(ctrl)
	if s.state == nil {
		t.Error("fail to NewStateHandler")
	}
	s.signer.EXPECT().SignHash(gomock.Any()).Return([]byte("2222"), nil).AnyTimes()
	s.signer.EXPECT().GetAddress().Return(common.Address{2, 3, 4}).AnyTimes()
	s.sender.EXPECT().BroadcastMsg(uint64(model2.TypeOfNewRoundMsg), gomock.Any()).AnyTimes()
	assert.NotPanics(t, func() {
		s.state.onEnterNewRound()
	})
}

func TestStateHandler_onEnterPropose(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s := testState(ctrl)
	if s.state == nil {
		t.Error("fail to NewStateHandler")
	}
	s.signer.EXPECT().GetAddress().Return(common.Address{2, 3, 4}).AnyTimes()
	assert.NotPanics(t, s.state.onEnterPropose)
}

//New functions
func TestStateHandler_finalBlock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s := testState(ctrl)
	s.signer.EXPECT().GetAddress().Return(common.Address{2, 3, 4}).AnyTimes()
	if s.state == nil {
		t.Error("fail to NewStateHandler")
	}
	s.signer.EXPECT().SignHash(gomock.Any()).Return([]byte("2222"), nil).AnyTimes()
	s.signer.EXPECT().GetAddress().Return(common.Address{2, 3, 4}).AnyTimes()
	s.sender.EXPECT().BroadcastMsg(gomock.Any(), gomock.Any()).AnyTimes()
	assert.NotPanics(t, func() {
		s.state.onEnterPrevote()
		s.state.onNewRoundTimeout()
		s.state.isCurProposer()
		s.state.onEnterPrecommit()
		s.state.onProposeTimeout()
		s.state.onPreVoteTimeout()
		s.state.onPreCommitTimeout()
		s.state.signAndPrevote(&model.VoteMsg{})
		s.state.bs.Height = 2
		s.state.bs.Votes.verifiers = append(s.state.bs.Votes.verifiers, common.Address{2, 3, 4})
		s.state.signAndVote(&model.VoteMsg{
			Height: 2,
			Witness: &model.WitMsg{
				Address: common.Address{2, 3, 4},
				Sign:    make([]byte, 65),
			},
		})
		s.state.addTimeoutCount("3232432")
	})
}
