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

package chain_communication

import (
	"bytes"
	"errors"
	"github.com/hashicorp/golang-lru"
	"testing"
	"time"

	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/chain-config"
	model2 "github.com/dipperin/dipperin-core/core/csbft/model"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/tests/factory"
	"github.com/dipperin/dipperin-core/third-party/p2p"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_makeNewBlockBroadcaster(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChain := NewMockChain(ctrl)
	mockPM := NewMockPeerManager(ctrl)
	mockPbftNode := NewMockPbftNode(ctrl)

	bbconfig := NewBlockBroadcasterConfig{
		Chain:    mockChain,
		Pm:       mockPM,
		PbftNode: mockPbftNode,
	}

	assert.NotNil(t, makeNewBlockBroadcaster(&bbconfig))
}

func TestNewBlockBroadcaster_MsgHandlers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChain := NewMockChain(ctrl)
	mockPM := NewMockPeerManager(ctrl)
	mockPbftNode := NewMockPbftNode(ctrl)

	bbconfig := NewBlockBroadcasterConfig{
		Chain:    mockChain,
		Pm:       mockPM,
		PbftNode: mockPbftNode,
	}

	bb := makeNewBlockBroadcaster(&bbconfig)

	assert.NotNil(t, bb.MsgHandlers()[NewBlockV1Msg])
}

func TestNewBlockBroadcaster_getReceiver(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChain := NewMockChain(ctrl)
	mockPM := NewMockPeerManager(ctrl)
	mockPbftNode := NewMockPbftNode(ctrl)

	bbconfig := NewBlockBroadcasterConfig{
		Chain:    mockChain,
		Pm:       mockPM,
		PbftNode: mockPbftNode,
	}

	bb := makeNewBlockBroadcaster(&bbconfig)

	mockPeer := NewMockPmAbstractPeer(ctrl)

	mockPeer.EXPECT().NodeName().Return("Test").AnyTimes()

	mockPeer.EXPECT().ID().Return("1").AnyTimes()

	receiver := bb.getReceiver(mockPeer)

	time.Sleep(100 * time.Millisecond)

	assert.NotNil(t, receiver)

	assert.Equal(t, bb.getReceiver(mockPeer), receiver)
}

func TestNewBlockBroadcaster_getPeersWithoutBlock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChain := NewMockChain(ctrl)
	mockPM := NewMockPeerManager(ctrl)
	mockPbftNode := NewMockPbftNode(ctrl)

	bbconfig := NewBlockBroadcasterConfig{
		Chain:    mockChain,
		Pm:       mockPM,
		PbftNode: mockPbftNode,
	}

	bb := makeNewBlockBroadcaster(&bbconfig)

	mockPeer := NewMockPmAbstractPeer(ctrl)

	mockPeer.EXPECT().ID().Return("1").AnyTimes()
	mockPeer.EXPECT().NodeName().Return("Test").AnyTimes()

	peers := make(map[string]PmAbstractPeer)

	peers[mockPeer.ID()] = mockPeer

	mockPM.EXPECT().GetPeers().Return(peers).Times(1)

	diff := common.HexToDiff("0x1effffff")
	fakeBlock := factory.CreateBlock2(diff, 1)

	assert.Equal(t, bb.getPeersWithoutBlock(fakeBlock), []PmAbstractPeer{mockPeer})
}

func TestNewBlockBroadcaster_broadcastBlock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChain := NewMockChain(ctrl)
	mockPM := NewMockPeerManager(ctrl)
	mockPbftNode := NewMockPbftNode(ctrl)

	bbconfig := NewBlockBroadcasterConfig{
		Chain:    mockChain,
		Pm:       mockPM,
		PbftNode: mockPbftNode,
	}

	bb := makeNewBlockBroadcaster(&bbconfig)

	mockPeer := NewMockPmAbstractPeer(ctrl)
	mockPeer.EXPECT().ID().Return("1").AnyTimes()
	mockPeer.EXPECT().NodeName().Return("Test").AnyTimes()
	mockPeer.EXPECT().NodeType().Return(uint64(chain_config.NodeTypeOfNormal)).AnyTimes()
	mockPeer.EXPECT().SendMsg(gomock.Any(), gomock.Any()).Return(nil).Times(1)

	mockPM.EXPECT().GetPeer("1").Return(mockPeer).Times(1)

	diff := common.HexToDiff("0x1effffff")
	fakeBlock := factory.CreateBlock2(diff, 1)

	bb.broadcastBlock(fakeBlock, []PmAbstractPeer{mockPeer})

	time.Sleep(100 * time.Millisecond)
}

func TestNewBlockBroadcaster_getTransferPeers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChain := NewMockChain(ctrl)
	mockPM := NewMockPeerManager(ctrl)
	mockPbftNode := NewMockPbftNode(ctrl)

	bbconfig := NewBlockBroadcasterConfig{
		Chain:    mockChain,
		Pm:       mockPM,
		PbftNode: mockPbftNode,
	}

	bb := makeNewBlockBroadcaster(&bbconfig)

	mockPeer := NewMockPmAbstractPeer(ctrl)

	assert.Equal(t, len(bb.getTransferPeers([]PmAbstractPeer{mockPeer})), 1)

	mockPeers := make([]PmAbstractPeer, 18)

	for i := 0; i <= 16; i++ {
		mockPeers = append(mockPeers, NewMockPmAbstractPeer(ctrl))
	}

	assert.Equal(t, len(bb.getTransferPeers(mockPeers)), 5)
}

func TestNewBlockBroadcaster_BroadcastBlock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChain := NewMockChain(ctrl)
	mockPM := NewMockPeerManager(ctrl)
	mockPbftNode := NewMockPbftNode(ctrl)

	bbconfig := NewBlockBroadcasterConfig{
		Chain:    mockChain,
		Pm:       mockPM,
		PbftNode: mockPbftNode,
	}

	bb := makeNewBlockBroadcaster(&bbconfig)

	mockPeer1 := NewMockPmAbstractPeer(ctrl)
	mockPeer1.EXPECT().ID().Return("1").AnyTimes()
	mockPeer1.EXPECT().NodeName().Return("Test").AnyTimes()
	mockPeer1.EXPECT().NodeType().Return(uint64(chain_config.NodeTypeOfNormal)).AnyTimes()
	mockPeer1.EXPECT().SendMsg(gomock.Any(), gomock.Any()).Return(nil).Times(1)

	mockPeer2 := NewMockPmAbstractPeer(ctrl)
	mockPeer2.EXPECT().ID().Return("2").AnyTimes()
	mockPeer2.EXPECT().NodeName().Return("Test2").AnyTimes()
	mockPeer2.EXPECT().NodeType().Return(uint64(chain_config.NodeTypeOfVerifier)).AnyTimes()
	mockPeer2.EXPECT().SendMsg(gomock.Any(), gomock.Any()).Return(nil).Times(1)

	peers := make(map[string]PmAbstractPeer)

	peers[mockPeer1.ID()] = mockPeer1
	peers[mockPeer2.ID()] = mockPeer2

	mockPM.EXPECT().GetPeers().Return(peers).Times(1)
	mockPM.EXPECT().GetPeer(mockPeer1.ID()).Return(mockPeer1).Times(1)
	mockPM.EXPECT().GetPeer(mockPeer2.ID()).Return(mockPeer2).Times(1)

	diff := common.HexToDiff("0x1effffff")
	fakeBlock := factory.CreateBlock2(diff, 1)

	bb.BroadcastBlock(fakeBlock)

	time.Sleep(100 * time.Millisecond)

}

func TestNewBlockBroadcaster_onNewBlock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChain := NewMockChain(ctrl)
	mockPM := NewMockPeerManager(ctrl)
	mockPbftNode := NewMockPbftNode(ctrl)

	bbconfig := NewBlockBroadcasterConfig{
		Chain:    mockChain,
		Pm:       mockPM,
		PbftNode: mockPbftNode,
	}

	bb := makeNewBlockBroadcaster(&bbconfig)

	msg := p2p.Msg{
		Payload: bytes.NewReader([]byte{}),
	}

	mockPeer := NewMockPmAbstractPeer(ctrl)
	mockPeer.EXPECT().ID().Return("1").AnyTimes()
	mockPeer.EXPECT().NodeName().Return("Test1").AnyTimes()

	err := bb.onNewBlock(msg, mockPeer)

	assert.Error(t, err)

	diff := common.HexToDiff("0x1effffff")
	fakeBlock := factory.CreateBlock2(diff, 1)

	payload, _ := rlp.EncodeToBytes(fakeBlock)

	msg = p2p.Msg{
		Payload: bytes.NewReader(payload),
	}

	mockPbftNode.EXPECT().OnNewWaitVerifyBlock(gomock.Any(), gomock.Any()).Times(1)

	err = bb.onNewBlock(msg, mockPeer)

	assert.NoError(t, err)
}

func TestNewBlockBroadcaster_newBlockReceiver(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChain := NewMockChain(ctrl)
	mockPM := NewMockPeerManager(ctrl)
	mockPbftNode := NewMockPbftNode(ctrl)

	bbconfig := NewBlockBroadcasterConfig{
		Chain:    mockChain,
		Pm:       mockPM,
		PbftNode: mockPbftNode,
	}

	bb := makeNewBlockBroadcaster(&bbconfig)

	mockPeer := NewMockPmAbstractPeer(ctrl)
	mockPeer.EXPECT().ID().Return("1").AnyTimes()
	mockPeer.EXPECT().NodeName().Return("Test1").AnyTimes()
	mockPeer.EXPECT().SendMsg(gomock.Any(), gomock.Any()).Return(errors.New("test")).Times(1)

	mockPM.EXPECT().GetPeer(gomock.Any()).Return(mockPeer).AnyTimes()

	receiver := bb.newBlockReceiver(mockPeer)

	diff := common.HexToDiff("0x1effffff")
	fakeBlock := factory.CreateBlock2(diff, 1)

	receiver.asyncSendBlock(fakeBlock)

	time.Sleep(100 * time.Millisecond)
}

func Test_blockReceiver_asyncSendBlock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPeer := NewMockPmAbstractPeer(ctrl)
	mockPeer.EXPECT().ID().Return("1").AnyTimes()
	mockPeer.EXPECT().NodeName().Return("Test1").AnyTimes()
	mockPeer.EXPECT().SendMsg(gomock.Any(), gomock.Any()).Return(errors.New("test")).Times(1)

	kb, _ := lru.New(500)

	getPeer := func() PmAbstractPeer {
		return mockPeer
	}

	receiver := &blockReceiver{
		peerID:          "1",
		peerName:        "Test",
		knownBlocks:     kb,
		queuedBlock:     make(chan model.AbstractBlock, maxQueuedBlock),
		queuedBlockHash: make(chan model.AbstractBlock, maxQueuedBlockHash),
	}

	diff := common.HexToDiff("0x1effffff")
	fakeBlock := factory.CreateBlock2(diff, 1)

	for i := 0; i < 5; i++ {
		receiver.asyncSendBlock(fakeBlock)
	}

	go func() {
		err := receiver.broadcast(getPeer)
		assert.Error(t, err)
	}()

	time.Sleep(100 * time.Millisecond)
}

func Test_blockReceiver_asyncSendBlockHash(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPeer := NewMockPmAbstractPeer(ctrl)
	mockPeer.EXPECT().ID().Return("1").AnyTimes()
	mockPeer.EXPECT().NodeName().Return("Test1").AnyTimes()
	mockPeer.EXPECT().SendMsg(gomock.Any(), gomock.Any()).Return(errors.New("test")).AnyTimes()

	kb, _ := lru.New(500)

	getPeer := func() PmAbstractPeer {
		return mockPeer
	}

	receiver := &blockReceiver{
		peerID:          "1",
		peerName:        "Test",
		knownBlocks:     kb,
		queuedBlock:     make(chan model.AbstractBlock, maxQueuedBlock),
		queuedBlockHash: make(chan model.AbstractBlock, maxQueuedBlockHash),
	}

	diff := common.HexToDiff("0x1effffff")
	fakeBlock := factory.CreateBlock2(diff, 1)

	for i := 0; i < 5; i++ {
		receiver.asyncSendBlockHash(fakeBlock)
	}

	go func() {
		err := receiver.broadcast(getPeer)
		assert.Error(t, err)
	}()

	time.Sleep(100 * time.Millisecond)
}

func Test_blockReceiver_asyncSendVerifyResult(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPeer := NewMockPmAbstractPeer(ctrl)
	mockPeer.EXPECT().ID().Return("1").AnyTimes()
	mockPeer.EXPECT().NodeName().Return("Test1").AnyTimes()
	mockPeer.EXPECT().SendMsg(gomock.Any(), gomock.Any()).Return(errors.New("test")).Times(1)

	kb, _ := lru.New(500)

	getPeer := func() PmAbstractPeer {
		return mockPeer
	}

	receiver := &blockReceiver{
		peerID:                 "1",
		peerName:               "Test",
		knownBlocks:            kb,
		queuedBlock:            make(chan model.AbstractBlock, maxQueuedBlock),
		queuedBlockHash:        make(chan model.AbstractBlock, maxQueuedBlockHash),
		queuedVerifyResult:     make(chan *model2.VerifyResult, maxQueuedBlock),
		queuedVerifyResultHash: make(chan *model2.VerifyResult, maxQueuedBlock),
	}

	diff := common.HexToDiff("0x1effffff")
	fakeBlock := factory.CreateBlock2(diff, 1)

	verifierBlock := &model2.VerifyResult{
		Block:       fakeBlock,
		SeenCommits: nil,
	}

	for i := 0; i < 5; i++ {
		receiver.asyncSendVerifyResult(verifierBlock)
	}

	go func() {
		err := receiver.broadcast(getPeer)
		assert.Error(t, err)
	}()

	time.Sleep(100 * time.Millisecond)
}

func Test_blockReceiver_asyncSendVerifyResultHash(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPeer := NewMockPmAbstractPeer(ctrl)
	mockPeer.EXPECT().ID().Return("1").AnyTimes()
	mockPeer.EXPECT().NodeName().Return("Test1").AnyTimes()
	mockPeer.EXPECT().SendMsg(gomock.Any(), gomock.Any()).Return(errors.New("test")).Times(1)

	kb, _ := lru.New(500)

	getPeer := func() PmAbstractPeer {
		return mockPeer
	}

	receiver := &blockReceiver{
		peerID:                 "1",
		peerName:               "Test",
		knownBlocks:            kb,
		queuedBlock:            make(chan model.AbstractBlock, maxQueuedBlock),
		queuedBlockHash:        make(chan model.AbstractBlock, maxQueuedBlockHash),
		queuedVerifyResult:     make(chan *model2.VerifyResult, maxQueuedBlock),
		queuedVerifyResultHash: make(chan *model2.VerifyResult, maxQueuedBlock),
	}

	diff := common.HexToDiff("0x1effffff")
	fakeBlock := factory.CreateBlock2(diff, 1)

	verifierBlock := &model2.VerifyResult{
		Block:       fakeBlock,
		SeenCommits: nil,
	}

	for i := 0; i < 5; i++ {
		receiver.asyncSendVerifyResultHash(verifierBlock)
	}

	go func() {
		err := receiver.broadcast(getPeer)
		assert.Error(t, err)
	}()

	time.Sleep(100 * time.Millisecond)
}

func Test_blockReceiver_sendBlock(t *testing.T) {

	kb, _ := lru.New(500)

	getPeer := func() PmAbstractPeer {
		return nil
	}

	diff := common.HexToDiff("0x1effffff")
	fakeBlock := factory.CreateBlock2(diff, 1)

	receiver := &blockReceiver{
		peerID:      "1",
		peerName:    "Test",
		knownBlocks: kb,
	}

	err := receiver.sendBlock(fakeBlock, getPeer)
	assert.Error(t, err)
}

func Test_blockReceiver_sendBlockHash(t *testing.T) {
	kb, _ := lru.New(500)

	getPeer := func() PmAbstractPeer {
		return nil
	}

	diff := common.HexToDiff("0x1effffff")
	fakeBlock := factory.CreateBlock2(diff, 1)

	receiver := &blockReceiver{
		peerID:      "1",
		peerName:    "Test",
		knownBlocks: kb,
	}

	err := receiver.sendBlockHash(fakeBlock, getPeer)
	assert.Error(t, err)
}

func Test_blockReceiver_sendVerifyResult(t *testing.T) {
	kb, _ := lru.New(500)

	getPeer := func() PmAbstractPeer {
		return nil
	}

	diff := common.HexToDiff("0x1effffff")
	fakeBlock := factory.CreateBlock2(diff, 1)

	verifierBlock := &model2.VerifyResult{
		Block:       fakeBlock,
		SeenCommits: nil,
	}

	receiver := &blockReceiver{
		peerID:      "1",
		peerName:    "Test",
		knownBlocks: kb,
	}

	err := receiver.sendVerifyResult(verifierBlock, getPeer)
	assert.Error(t, err)
}

func Test_blockReceiver_sendVerifyResultHash(t *testing.T) {
	kb, _ := lru.New(500)

	getPeer := func() PmAbstractPeer {
		return nil
	}

	diff := common.HexToDiff("0x1effffff")
	fakeBlock := factory.CreateBlock2(diff, 1)

	verifierBlock := &model2.VerifyResult{
		Block:       fakeBlock,
		SeenCommits: nil,
	}

	receiver := &blockReceiver{
		peerID:      "1",
		peerName:    "Test",
		knownBlocks: kb,
	}

	err := receiver.sendVerifyResultHash(verifierBlock, getPeer)
	assert.Error(t, err)
}
