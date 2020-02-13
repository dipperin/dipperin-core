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

package chaincommunication

import (
	"net"
	"testing"

	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/chainconfig"
	"github.com/dipperin/dipperin-core/third_party/p2p"
	"github.com/dipperin/dipperin-core/third_party/p2p/enode"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_newPeer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockP2PPeer := NewMockP2PPeer(ctrl)
	mockReadWriter := NewMockMsgReadWriter(ctrl)

	mockP2PPeer.EXPECT().ID().Return(enode.HexID("0x00000000000000806ad9b61fa5ae014307ebdc964253adcd9f2c0a392aa11abc"))

	assert.NotNil(t, newPeer(1, mockP2PPeer, mockReadWriter))

}

func Test_peer_GetCsPeerInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockP2PPeer := NewMockP2PPeer(ctrl)
	mockReadWriter := NewMockMsgReadWriter(ctrl)

	mockP2PPeer.EXPECT().ID().Return(enode.HexID("0x00000000000000806ad9b61fa5ae014307ebdc964253adcd9f2c0a392aa11abc"))

	peer := newPeer(1, mockP2PPeer, mockReadWriter)

	assert.NotNil(t, peer.GetCsPeerInfo())
}

func Test_peer_SetNotRunning(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockP2PPeer := NewMockP2PPeer(ctrl)
	mockReadWriter := NewMockMsgReadWriter(ctrl)

	mockP2PPeer.EXPECT().ID().Return(enode.HexID("0x00000000000000806ad9b61fa5ae014307ebdc964253adcd9f2c0a392aa11abc"))

	peer := newPeer(1, mockP2PPeer, mockReadWriter)

	assert.Equal(t, peer.IsRunning(), true)

	peer.SetNotRunning()

	assert.Equal(t, peer.IsRunning(), false)
}

func Test_peer_SetNodeName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockP2PPeer := NewMockP2PPeer(ctrl)
	mockReadWriter := NewMockMsgReadWriter(ctrl)

	mockP2PPeer.EXPECT().ID().Return(enode.HexID("0x00000000000000806ad9b61fa5ae014307ebdc964253adcd9f2c0a392aa11abc"))

	peer := newPeer(1, mockP2PPeer, mockReadWriter)

	assert.Equal(t, peer.NodeName(), "")

	peer.SetNodeName("test")

	assert.Equal(t, peer.NodeName(), "test")
}

func Test_peer_SetNodeType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockP2PPeer := NewMockP2PPeer(ctrl)
	mockReadWriter := NewMockMsgReadWriter(ctrl)

	mockP2PPeer.EXPECT().ID().Return(enode.HexID("0x00000000000000806ad9b61fa5ae014307ebdc964253adcd9f2c0a392aa11abc"))

	peer := newPeer(1, mockP2PPeer, mockReadWriter)

	assert.Equal(t, peer.NodeType(), uint64(0))

	peer.SetNodeType(chainconfig.NodeTypeOfNormal)

	assert.Equal(t, peer.NodeType(), uint64(chainconfig.NodeTypeOfNormal))
}

func Test_peer_RemoteAddress(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockP2PPeer := NewMockP2PPeer(ctrl)
	mockReadWriter := NewMockMsgReadWriter(ctrl)

	mockP2PPeer.EXPECT().ID().Return(enode.HexID("0x00000000000000806ad9b61fa5ae014307ebdc964253adcd9f2c0a392aa11abc"))

	addr := &net.IPAddr{}

	mockP2PPeer.EXPECT().RemoteAddr().Return(addr)

	peer := newPeer(1, mockP2PPeer, mockReadWriter)

	assert.Equal(t, peer.RemoteAddress(), addr)
}

func Test_peer_SetRemoteVerifierAddress(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockP2PPeer := NewMockP2PPeer(ctrl)
	mockReadWriter := NewMockMsgReadWriter(ctrl)

	mockP2PPeer.EXPECT().ID().Return(enode.HexID("0x00000000000000806ad9b61fa5ae014307ebdc964253adcd9f2c0a392aa11abc"))

	peer := newPeer(1, mockP2PPeer, mockReadWriter)

	assert.Equal(t, peer.RemoteVerifierAddress(), common.HexToAddress("0x0"))

	peer.SetRemoteVerifierAddress(common.HexToAddress("0x123"))

	assert.Equal(t, peer.RemoteVerifierAddress(), common.HexToAddress("0x123"))

}

func Test_peer_SetPeerRawUrl(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockP2PPeer := NewMockP2PPeer(ctrl)
	mockReadWriter := NewMockMsgReadWriter(ctrl)

	mockP2PPeer.EXPECT().ID().Return(enode.HexID("0x00000000000000806ad9b61fa5ae014307ebdc964253adcd9f2c0a392aa11abc"))

	peer := newPeer(1, mockP2PPeer, mockReadWriter)

	assert.Equal(t, peer.GetPeerRawUrl(), "")

	peer.SetPeerRawUrl("test")

	assert.Equal(t, peer.GetPeerRawUrl(), "test")
}

func Test_peer_SetHead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockP2PPeer := NewMockP2PPeer(ctrl)
	mockReadWriter := NewMockMsgReadWriter(ctrl)

	mockP2PPeer.EXPECT().ID().Return(enode.HexID("0x00000000000000806ad9b61fa5ae014307ebdc964253adcd9f2c0a392aa11abc"))

	peer := newPeer(1, mockP2PPeer, mockReadWriter)

	head, height := peer.GetHead()

	assert.Equal(t, head, common.HexToHash("0x0"))
	assert.Equal(t, height, uint64(0))

	peer.SetHead(common.HexToHash("0x1"), uint64(1))

	head, height = peer.GetHead()

	assert.Equal(t, head, common.HexToHash("0x1"))
	assert.Equal(t, height, uint64(1))
}

func Test_peer_DisconnectPeer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockP2PPeer := NewMockP2PPeer(ctrl)
	mockReadWriter := NewMockMsgReadWriter(ctrl)

	mockP2PPeer.EXPECT().ID().Return(enode.HexID("0x00000000000000806ad9b61fa5ae014307ebdc964253adcd9f2c0a392aa11abc"))
	mockP2PPeer.EXPECT().Disconnect(gomock.Any())

	peer := newPeer(1, mockP2PPeer, mockReadWriter)

	peer.DisconnectPeer()

	assert.Equal(t, peer.IsRunning(), false)
}

func Test_peer_ReadMsg(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockP2PPeer := NewMockP2PPeer(ctrl)
	mockReadWriter := NewMockMsgReadWriter(ctrl)

	mockP2PPeer.EXPECT().ID().Return(enode.HexID("0x00000000000000806ad9b61fa5ae014307ebdc964253adcd9f2c0a392aa11abc"))

	peer := newPeer(1, mockP2PPeer, mockReadWriter)

	msg := p2p.Msg{}

	mockReadWriter.EXPECT().ReadMsg().Return(msg, nil)

	resMsg, err := peer.ReadMsg()

	assert.Equal(t, resMsg, msg)
	assert.NoError(t, err)
}

func Test_peer_ID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockP2PPeer := NewMockP2PPeer(ctrl)
	mockReadWriter := NewMockMsgReadWriter(ctrl)

	mockP2PPeer.EXPECT().ID().Return(enode.HexID("0x00000000000000806ad9b61fa5ae014307ebdc964253adcd9f2c0a392aa11abc"))

	peer := newPeer(1, mockP2PPeer, mockReadWriter)

	assert.Equal(t, peer.ID(), "00000000000000806ad9b61fa5ae014307ebdc964253adcd9f2c0a392aa11abc")
}

func Test_peer_SendMsg(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockP2PPeer := NewMockP2PPeer(ctrl)
	mockReadWriter := NewMockMsgReadWriter(ctrl)

	mockP2PPeer.EXPECT().ID().Return(enode.HexID("0x00000000000000806ad9b61fa5ae014307ebdc964253adcd9f2c0a392aa11abc"))

	peer := newPeer(1, mockP2PPeer, mockReadWriter)

	mockReadWriter.EXPECT().WriteMsg(gomock.Any()).Return(nil)

	err := peer.SendMsg(uint64(1), "test")

	assert.NoError(t, err)
}

func Test_newPeerSet(t *testing.T) {
	assert.NotNil(t, newPeerSet())
}

func Test_peerSet_AddPeer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	peerSet := newPeerSet()

	mockPeer := NewMockPmAbstractPeer(ctrl)
	mockPeer.EXPECT().ID().Return("")

	err := peerSet.AddPeer(mockPeer)

	assert.Error(t, err)

	mockPeer.EXPECT().ID().Return("1").AnyTimes()
	mockPeer.EXPECT().NodeName().Return("test").AnyTimes()

	peerSet.closed = true

	err = peerSet.AddPeer(mockPeer)

	assert.Error(t, err)

	peerSet.closed = false

	err = peerSet.AddPeer(mockPeer)

	assert.NoError(t, err)

	err = peerSet.AddPeer(mockPeer)

	assert.Error(t, err)

}

func Test_peerSet_RemovePeer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	peerSet := newPeerSet()

	mockPeer := NewMockPmAbstractPeer(ctrl)
	mockPeer.EXPECT().ID().Return("1").AnyTimes()
	mockPeer.EXPECT().NodeName().Return("test").AnyTimes()

	err := peerSet.RemovePeer(mockPeer.ID())

	assert.Error(t, err)

	peerSet.AddPeer(mockPeer)

	assert.Equal(t, len(peerSet.GetPeers()), 1)

	err = peerSet.RemovePeer(mockPeer.ID())

	assert.NoError(t, err)

	assert.Equal(t, len(peerSet.GetPeers()), 0)
}

func Test_peerSet_ReplacePeers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	peerSet := newPeerSet()

	mockPeer := NewMockPmAbstractPeer(ctrl)
	mockPeer.EXPECT().ID().Return("1").AnyTimes()
	mockPeer.EXPECT().NodeName().Return("test").AnyTimes()

	peers := make(map[string]PmAbstractPeer)
	peers[mockPeer.ID()] = mockPeer

	peerSet.ReplacePeers(peers)

	assert.Equal(t, peerSet.GetPeers(), peers)
}

func Test_peerSet_BestPeer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	peerSet := newPeerSet()

	mockPeer := NewMockPmAbstractPeer(ctrl)
	mockPeer.EXPECT().ID().Return("1").AnyTimes()
	mockPeer.EXPECT().NodeName().Return("test").AnyTimes()
	mockPeer.EXPECT().GetHead().Return(common.HexToHash("0x1"), uint64(2))

	peerSet.AddPeer(mockPeer)

	assert.Equal(t, peerSet.BestPeer(), mockPeer)
}

func Test_peerSet_Peer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	peerSet := newPeerSet()

	mockPeer := NewMockPmAbstractPeer(ctrl)
	mockPeer.EXPECT().ID().Return("1").AnyTimes()
	mockPeer.EXPECT().NodeName().Return("test").AnyTimes()

	peerSet.AddPeer(mockPeer)

	assert.Equal(t, peerSet.Peer(mockPeer.ID()), mockPeer)
}

func Test_peerSet_Len(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	peerSet := newPeerSet()

	mockPeer := NewMockPmAbstractPeer(ctrl)
	mockPeer.EXPECT().ID().Return("1").AnyTimes()
	mockPeer.EXPECT().NodeName().Return("test").AnyTimes()

	peerSet.AddPeer(mockPeer)

	assert.Equal(t, peerSet.Len(), 1)
}

func Test_peerSet_Close(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	peerSet := newPeerSet()

	mockPeer := NewMockPmAbstractPeer(ctrl)
	mockPeer.EXPECT().ID().Return("1").AnyTimes()
	mockPeer.EXPECT().NodeName().Return("test").AnyTimes()
	mockPeer.EXPECT().DisconnectPeer()

	peerSet.AddPeer(mockPeer)

	peerSet.Close()
}

func Test_peerSet_GetPeersInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	peerSet := newPeerSet()

	mockPeer := NewMockPmAbstractPeer(ctrl)
	mockPeer.EXPECT().ID().Return("1").AnyTimes()
	mockPeer.EXPECT().NodeName().Return("test").AnyTimes()
	mockPeer.EXPECT().GetCsPeerInfo().Return(&p2p.CsPeerInfo{})

	peerSet.AddPeer(mockPeer)

	assert.Equal(t, len(peerSet.GetPeersInfo()), 1)
}
