package statemachine

import (
	"testing"
	"github.com/golang/mock/gomock"
	"github.com/dipperin/dipperin-core/core/csbft/components"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/model"
)

func testState(t *testing.T) *StateHandler {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
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
	return NewStateHandler(config, DefaultConfig, bp)
}

func TestNewStateHandler(t *testing.T) {
	if state := testState(t); state == nil {
		t.Error("fail to NewStateHandler")
	}
}
