package csbftnode

import (
	"testing"
	"github.com/golang/mock/gomock"
	"github.com/dipperin/dipperin-core/core/csbft/statemachine"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/stretchr/testify/assert"
	"github.com/dipperin/dipperin-core/core/csbft/components"
)

type myStr struct {
	csbft        *CsBft
	chain        *MockChainReader
	fetcher      *MockFetcher
	signer       *MockMsgSigner
	sender       *MockMsgSender
	validator    *MockValidator
	block        *MockAbstractBlock
	verification *MockAbstractVerification
}

func testReady(ctrl *gomock.Controller) (my myStr) {
	chain := NewMockChainReader(ctrl)
	fetcher := NewMockFetcher(ctrl)
	signer := NewMockMsgSigner(ctrl)
	sender := NewMockMsgSender(ctrl)
	validator := NewMockValidator(ctrl)
	block := NewMockAbstractBlock(ctrl)
	verification := NewMockAbstractVerification(ctrl)
	var config = &statemachine.BftConfig{
		ChainReader: chain,
		Fetcher:     fetcher,
		Signer:      signer,
		Sender:      sender,
		Validator:   validator,
	}
	var height uint64
	height = 1
	block.EXPECT().Number().Return(height).AnyTimes()
	chain.EXPECT().CurrentBlock().Return(block).AnyTimes()
	chain.EXPECT().GetSeenCommit(height).Return([]model.AbstractVerification{verification}).AnyTimes()
	verification.EXPECT().GetRound().Return(uint64(2333)).AnyTimes()
	chain.EXPECT().IsChangePoint(gomock.Any(), gomock.Any()).Return(true).AnyTimes()
	var address []common.Address
	address = append(address, common.Address{1, 2, 3, 4})
	chain.EXPECT().GetNextVerifiers().Return(address).AnyTimes()

	my.csbft = NewCsBft(config)
	my.fetcher = fetcher
	my.signer = signer
	my.sender = sender
	my.validator = validator
	my.block = block
	my.verification = verification
	return
}

func TestNewCsBft(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	if c := testReady(ctrl); c.csbft == nil {
		t.Error("fail to NewCsBft")
	}
}

func TestCsBft_OnEnterNewHeight(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := testReady(ctrl)
	if c.csbft == nil {
		t.Error("fail to NewCsBft")
		return
	}

	assert.NotPanics(t, func() {
		c.csbft.OnEnterNewHeight(1)
	})
}

func TestCsBft_SetFetcher(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := testReady(ctrl)
	if c.csbft == nil {
		t.Error("fail to NewCsBft")
		return
	}
	assert.NotPanics(t, func() {
		var s components.CsBftFetcher
		c.csbft.SetFetcher(&s)
	})
	assert.NotNil(t, c.csbft.stateHandler.Fetcher)
}

func TestCsBft_Start(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := testReady(ctrl)
	if c.csbft == nil {
		t.Error("fail to NewCsBft")
		return
	}
	c.signer.EXPECT().GetAddress().Return(common.Address{1, 2, 3})
	assert.NoError(t, c.csbft.Start())
}
