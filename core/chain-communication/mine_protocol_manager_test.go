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
	"testing"
	"time"

	"github.com/dipperin/dipperin-core/third-party/p2p"
	"github.com/dipperin/dipperin-core/third-party/p2p/enode"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type mockP2pMsgHandler struct {
	withError bool
}

func (pmh *mockP2pMsgHandler) OnNewMsg(msg p2p.Msg, p PmAbstractPeer) error {
	if pmh.withError {
		return errors.New("test")
	}
	return nil
}

func (pmh *mockP2pMsgHandler) SetMineMasterPeer(peer PmAbstractPeer) {

}

type mockMsgWriter struct {
}

func (mockMsgWriter) WriteMsg(p2p.Msg) error {
	return nil
}

type mockMsgReader struct {
}

func (mockMsgReader) ReadMsg() (p2p.Msg, error) {
	return p2p.Msg{
		Code:    uint64(1),
		Payload: bytes.NewReader([]byte{}),
	}, nil
}

type mockMsgReadWriter struct {
	mockMsgReader
	mockMsgWriter
}

func TestNewMineProtocolManager(t *testing.T) {
	assert.NotNil(t, NewMineProtocolManager(&mockP2pMsgHandler{}))
}

func TestMineProtocolManager_GetProtocol(t *testing.T) {
	mpm := NewMineProtocolManager(&mockP2pMsgHandler{})
	protocol := mpm.GetProtocol()

	assert.Equal(t, protocol.Length, uint64(512))
	assert.Equal(t, protocol.Name, "dipperin_mine")
	assert.Equal(t, protocol.Version, uint(1))
}

func TestMineProtocolManager_GetProtocol_Run(t *testing.T) {
	mpm := NewMineProtocolManager(&mockP2pMsgHandler{})
	protocol := mpm.GetProtocol()
	p2pPeer := p2p.NewPeer(enode.ID{0}, "test", []p2p.Cap{})
	go protocol.Run(p2pPeer, &mockMsgReadWriter{
		mockMsgReader{},
		mockMsgWriter{},
	})
	time.Sleep(50 * time.Millisecond)
}

func TestMineProtocolManager_handleMsg(t *testing.T) {
	mpm := NewMineProtocolManager(&mockP2pMsgHandler{})

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// peer
	mockPeer := NewMockPmAbstractPeer(ctrl)
	mockPeer.EXPECT().ReadMsg().Return(p2p.Msg{}, errors.New("test")).Times(1)

	assert.Error(t, mpm.handleMsg(mockPeer))

	mockPeer.EXPECT().ReadMsg().Return(p2p.Msg{
		Code:    uint64(1),
		Payload: bytes.NewReader([]byte{}),
		Size:    ProtocolMaxMsgSize + 1,
	}, nil).Times(1)

	assert.Error(t, mpm.handleMsg(mockPeer))

	mpm2 := NewMineProtocolManager(&mockP2pMsgHandler{
		withError: true,
	})

	mockPeer.EXPECT().ReadMsg().Return(p2p.Msg{
		Code:    uint64(1),
		Payload: bytes.NewReader([]byte{}),
	}, nil).Times(1)

	mockPeer.EXPECT().SetNotRunning().Times(1)

	assert.Error(t, mpm2.handleMsg(mockPeer))

}

func TestMineProtocolManager_handle(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// peer
	mockPeer := NewMockPmAbstractPeer(ctrl)
	mockPeerSet := NewMockAbstractPeerSet(ctrl)

	mpm := &MineProtocolManager{
		msgHandler: &mockP2pMsgHandler{},
		maxPeers:   1,
		peers:      mockPeerSet,
	}

	mockPeer.EXPECT().ID().Return("test").AnyTimes()

	mockPeerSet.EXPECT().Len().Return(2).AnyTimes()

	assert.Error(t, mpm.handle(mockPeer))

	mpm2 := &MineProtocolManager{
		msgHandler: &mockP2pMsgHandler{},
		maxPeers:   100,
		peers:      mockPeerSet,
	}

	mockPeerSet.EXPECT().AddPeer(mockPeer).Return(errors.New("test")).Times(1)

	assert.Error(t, mpm2.handle(mockPeer))

	mockPeerSet.EXPECT().AddPeer(mockPeer).Return(nil).Times(1)

	mockPeerSet.EXPECT().Peer(mockPeer.ID()).Return(nil).Times(1)

	mockPeer.EXPECT().ReadMsg().Return(p2p.Msg{}, errors.New("test")).Times(1)

	assert.Error(t, mpm2.handle(mockPeer))

}

func TestMineProtocolManager_removePeer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// peer
	mockPeer := NewMockPmAbstractPeer(ctrl)
	mockPeerSet := NewMockAbstractPeerSet(ctrl)

	mpm := &MineProtocolManager{
		msgHandler: &mockP2pMsgHandler{},
		maxPeers:   1,
		peers:      mockPeerSet,
	}

	mockPeerSet.EXPECT().Peer("").Return(mockPeer).Times(1)

	mockPeerSet.EXPECT().RemovePeer("").Return(errors.New("test")).Times(1)

	mockPeer.EXPECT().DisconnectPeer().Times(1)

	mpm.removePeer("")
}
