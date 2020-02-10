package csbftnode

import (
	"testing"
	"github.com/golang/mock/gomock"
	"github.com/dipperin/dipperin-core/core/csbft/statemachine"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/stretchr/testify/assert"
	"github.com/dipperin/dipperin-core/core/csbft/components"
	model2 "github.com/dipperin/dipperin-core/core/csbft/model"
	"github.com/dipperin/dipperin-core/third_party/p2p"
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

func TestCsBft_Stop(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	c := testReady(ctrl)
	if c.csbft == nil {
		t.Error("fail to NewCsBft")
		return
	}
	assert.Panics(t, c.csbft.Stop)
}

func TestCsBft_OnNewWaitVerifyBlock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	c := testReady(ctrl)
	if c.csbft == nil {
		t.Error("fail to NewCsBft")
		return
	}
	assert.NotPanics(t, func() {
		c.csbft.OnNewWaitVerifyBlock(c.block, "2222")
	})
}

func TestCsBft_broadcastFetchBlockMsg(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	c := testReady(ctrl)
	if c.csbft == nil {
		t.Error("fail to NewCsBft")
		return
	}
	assert.NotPanics(t, func() {
		c.sender.EXPECT().BroadcastMsg(uint64(model2.TypeOfSyncBlockMsg), common.Hash{})
		c.csbft.broadcastFetchBlockMsg(common.Hash{})
	})
}

func TestCsBft_OnNewMsg(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	c := testReady(ctrl)
	if c.csbft == nil {
		t.Error("fail to NewCsBft")
		return
	}
	assert.NoError(t, c.csbft.OnNewMsg(nil))
}

func TestCsBft_AddPeer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	c := testReady(ctrl)
	if c.csbft == nil {
		t.Error("fail to NewCsBft")
		return
	}
	assert.NoError(t, c.csbft.AddPeer(nil))
}

func TestCsBft_ChangePrimary(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	c := testReady(ctrl)
	if c.csbft == nil {
		t.Error("fail to NewCsBft")
		return
	}
	c.signer.EXPECT().GetAddress().Return(common.Address{1, 2, 3})
	assert.Panics(t, func() {
		c.csbft.ChangePrimary("222")
	})
}

func TestCsBft_canStart(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	c := testReady(ctrl)
	if c.csbft == nil {
		t.Error("fail to NewCsBft")
		return
	}
	c.signer.EXPECT().GetAddress().Return(common.Address{1, 2, 3})
	assert.Equal(t, false, c.csbft.canStart())
}

func testForIsCurrentVerifier(ctrl *gomock.Controller, a1, a2 common.Address) bool {
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

	var address1 []common.Address
	address1 = append(address1, a1)
	chain.EXPECT().GetCurrVerifiers().Return(address1).AnyTimes()
	var address []common.Address
	address = append(address, common.Address{1, 2, 3, 4})
	chain.EXPECT().GetNextVerifiers().Return(address).AnyTimes()
	signer.EXPECT().GetAddress().Return(a2).AnyTimes()
	return NewCsBft(config).isCurrentVerifier()
}

func TestCsBft_isCurrentVerifier(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	testCases := []struct {
		name   string
		given  func() bool
		expect bool
	}{
		{
			name: "IsCurrentVerifier true",
			given: func() bool {
				return testForIsCurrentVerifier(ctrl, common.Address{1, 2, 3}, common.Address{1, 2, 3})
			},
			expect: true,
		},

		{
			name: "IsCurrentVerifier false",
			given: func() bool {
				return testForIsCurrentVerifier(ctrl, common.Address{1, 2, 3}, common.Address{3, 2, 1})
			},
			expect: false,
		},
	}

	for i, tc := range testCases {
		sign := tc.given()
		if assert.Equal(t, testCases[i].expect, sign) {
			t.Log("success")
		} else {
			t.Log("failure")
		}
	}
}

func TestCsBft_OnNewP2PMsg(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	c := testReady(ctrl)
	if c.csbft == nil {
		t.Error("fail to NewCsBft")
		return
	}
	assert.NoError(t, c.csbft.OnNewP2PMsg(p2p.Msg{}, nil))
}

func TestCsBft_onSyncBlockMsg(t *testing.T)  {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	c := testReady(ctrl)
	if c.csbft == nil {
		t.Error("fail to NewCsBft")
		return
	}

	assert.NotPanics(t, func() {
		c.csbft.onSyncBlockMsg(common.Address{},common.Hash{})
	})
}
