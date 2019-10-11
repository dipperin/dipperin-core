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
	"encoding/hex"
	"errors"
	"math/big"
	"net"
	"os"
	"testing"
	"time"

	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/tests"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"github.com/dipperin/dipperin-core/third-party/p2p"
	"github.com/dipperin/dipperin-core/third-party/p2p/enode"
	"github.com/dipperin/dipperin-core/third-party/p2p/enr"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestName(t *testing.T) {
	//t.Skip("aaa")
	var pyRecord, _ = hex.DecodeString("f884b8407098ad865b00a582051940cb9cf36836572411a47278783077011599ed5cd16b76f2635f4e234738f30813a89eb9137e3e3df5266e3a1f11df72ecf1145ccb9c01826964827634826970847f00000189736563703235366b31a103ca634cae0d49acb401d8a4c6b6fe8c55b70d115bf400769cc1400f3258cd31388375647082765f")
	var r enr.Record
	if err := rlp.DecodeBytes(pyRecord, &r); err != nil {
		t.Fatalf("can't decode: %v", err)
	}
	n, err := enode.New(enode.ValidSchemes, &r)
	if err != nil {
		t.Fatalf("can't verify record: %v", err)
	}

	assert.NotNil(t, n)

}

func TestCsProtocolManager_ShowPmInfo(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// node conf
	mockNodeConf := NewMockNodeConf(ctrl)
	mockNodeConf.EXPECT().GetNodeType().Return(chain_config.NodeTypeOfVerifier).AnyTimes()
	mockNodeConf.EXPECT().GetNodeName().Return("aaaa").AnyTimes()

	// chain
	mockChain := NewMockChain(ctrl)

	// p2p server
	mockP2pServer := NewMockP2PServer(ctrl)

	var pyRecord, _ = hex.DecodeString("f884b8407098ad865b00a582051940cb9cf36836572411a47278783077011599ed5cd16b76f2635f4e234738f30813a89eb9137e3e3df5266e3a1f11df72ecf1145ccb9c01826964827634826970847f00000189736563703235366b31a103ca634cae0d49acb401d8a4c6b6fe8c55b70d115bf400769cc1400f3258cd31388375647082765f")

	var r enr.Record
	if err := rlp.DecodeBytes(pyRecord, &r); err != nil {
		t.Fatalf("can't decode: %v", err)
	}
	n, err := enode.New(enode.ValidSchemes, &r)
	if err != nil {
		t.Fatalf("can't verify record: %v", err)
	}

	mockP2pServer.EXPECT().Self().Return(n)

	// verifier
	mockVerifiersReader := NewMockVerifiersReader(ctrl)

	// pbft node
	mockPbftNode := NewMockPbftNode(ctrl)

	// PbftSigner
	mockPbftSigner := NewMockPbftSigner(ctrl)

	cfg := &CsProtocolManagerConfig{
		ChainConfig:     *chain_config.GetChainConfig(),
		Chain:           mockChain,
		P2PServer:       mockP2pServer,
		NodeConf:        mockNodeConf,
		VerifiersReader: mockVerifiersReader,
		PbftNode:        mockPbftNode,
		MsgSigner:       mockPbftSigner,
	}

	// mock peer set
	mockPeerSetBasePeers := NewMockAbstractPeerSet(ctrl)
	mockPeerSetCurrentVerifierPeers := NewMockAbstractPeerSet(ctrl)
	mockPeerSetnNextVerifierPeers := NewMockAbstractPeerSet(ctrl)
	mockPeerSetVerifierBootNode := NewMockAbstractPeerSet(ctrl)

	// mock peer
	mockPeerSetBasePeers.EXPECT().GetPeersInfo().Return(nil).AnyTimes()
	mockPeerSetCurrentVerifierPeers.EXPECT().GetPeersInfo().Return(nil).AnyTimes()
	mockPeerSetnNextVerifierPeers.EXPECT().GetPeersInfo().Return(nil).AnyTimes()
	mockPeerSetVerifierBootNode.EXPECT().GetPeersInfo().Return(nil).AnyTimes()

	psm := &CsPmPeerSetManager{
		basePeers:            mockPeerSetBasePeers,
		currentVerifierPeers: mockPeerSetCurrentVerifierPeers,
		nextVerifierPeers:    mockPeerSetnNextVerifierPeers,
		verifierBootNode:     mockPeerSetVerifierBootNode,
	}

	pm := &CsProtocolManager{
		CsProtocolManagerConfig: cfg,
		peerSetManager:          psm,
	}

	pm.ShowPmInfo()
}

func TestCsProtocolManager_BestPeer(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// node conf
	mockNodeConf := NewMockNodeConf(ctrl)
	mockNodeConf.EXPECT().GetNodeType().Return(chain_config.NodeTypeOfVerifier).AnyTimes()
	mockNodeConf.EXPECT().GetNodeName().Return("aaaa").AnyTimes()

	// chain
	mockChain := NewMockChain(ctrl)

	// p2p server
	mockP2pServer := NewMockP2PServer(ctrl)

	// verifier
	mockVerifiersReader := NewMockVerifiersReader(ctrl)

	// pbft node
	mockPbftNode := NewMockPbftNode(ctrl)

	// PbftSigner
	mockPbftSigner := NewMockPbftSigner(ctrl)

	cfg := &CsProtocolManagerConfig{
		ChainConfig:     *chain_config.GetChainConfig(),
		Chain:           mockChain,
		P2PServer:       mockP2pServer,
		NodeConf:        mockNodeConf,
		VerifiersReader: mockVerifiersReader,
		PbftNode:        mockPbftNode,
		MsgSigner:       mockPbftSigner,
	}

	// mock peer set
	mockPeerSetBasePeers := NewMockAbstractPeerSet(ctrl)
	mockPeerSetCurrentVerifierPeers := NewMockAbstractPeerSet(ctrl)
	mockPeerSetnNextVerifierPeers := NewMockAbstractPeerSet(ctrl)
	mockPeerSetVerifierBootNode := NewMockAbstractPeerSet(ctrl)

	// mock peer
	mockPeer := NewMockPmAbstractPeer(ctrl)
	mockPeer.EXPECT().GetHead().Return(common.HexToHash("s"), uint64(11)).AnyTimes()
	mockPeerSetBasePeers.EXPECT().BestPeer().Return(mockPeer).AnyTimes()

	// mock peer
	mockPeer1 := NewMockPmAbstractPeer(ctrl)
	mockPeer1.EXPECT().GetHead().Return(common.HexToHash("s"), uint64(10)).AnyTimes()
	mockPeerSetCurrentVerifierPeers.EXPECT().BestPeer().Return(mockPeer1).AnyTimes()

	// mock peer
	mockPeer2 := NewMockPmAbstractPeer(ctrl)
	mockPeer2.EXPECT().GetHead().Return(common.HexToHash("s"), uint64(9)).AnyTimes()
	mockPeerSetnNextVerifierPeers.EXPECT().BestPeer().Return(mockPeer2).AnyTimes()

	// mock peer
	mockPeer3 := NewMockPmAbstractPeer(ctrl)
	mockPeer3.EXPECT().GetHead().Return(common.HexToHash("s"), uint64(8)).AnyTimes()
	mockPeerSetVerifierBootNode.EXPECT().BestPeer().Return(mockPeer3).AnyTimes()

	psm := &CsPmPeerSetManager{
		basePeers:            mockPeerSetBasePeers,
		currentVerifierPeers: mockPeerSetCurrentVerifierPeers,
		nextVerifierPeers:    mockPeerSetnNextVerifierPeers,
		verifierBootNode:     mockPeerSetVerifierBootNode,
	}

	pm := &CsProtocolManager{
		CsProtocolManagerConfig: cfg,
		peerSetManager:          psm,
	}

	bestPeer := pm.BestPeer()

	_, height := bestPeer.GetHead()

	assert.Equal(t, uint64(11), height)
}

func TestCsProtocolManager_RemovePeer(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// node conf
	mockNodeConf := NewMockNodeConf(ctrl)
	mockNodeConf.EXPECT().GetNodeType().Return(chain_config.NodeTypeOfVerifier).AnyTimes()
	mockNodeConf.EXPECT().GetNodeName().Return("aaaa").AnyTimes()

	// chain
	mockChain := NewMockChain(ctrl)

	// p2p server
	mockP2pServer := NewMockP2PServer(ctrl)

	// verifier
	mockVerifiersReader := NewMockVerifiersReader(ctrl)

	// pbft node
	mockPbftNode := NewMockPbftNode(ctrl)

	// PbftSigner
	mockPbftSigner := NewMockPbftSigner(ctrl)

	cfg := &CsProtocolManagerConfig{
		ChainConfig:     *chain_config.GetChainConfig(),
		Chain:           mockChain,
		P2PServer:       mockP2pServer,
		NodeConf:        mockNodeConf,
		VerifiersReader: mockVerifiersReader,
		PbftNode:        mockPbftNode,
		MsgSigner:       mockPbftSigner,
	}

	// mock peer set
	mockPeerSetBasePeers := NewMockAbstractPeerSet(ctrl)
	mockPeerSetCurrentVerifierPeers := NewMockAbstractPeerSet(ctrl)
	mockPeerSetnNextVerifierPeers := NewMockAbstractPeerSet(ctrl)
	mockPeerSetVerifierBootNode := NewMockAbstractPeerSet(ctrl)

	// mock peer
	mockPeerSetBasePeers.EXPECT().Peer(gomock.Any()).Return(nil).AnyTimes()
	mockPeerSetCurrentVerifierPeers.EXPECT().Peer(gomock.Any()).Return(nil).AnyTimes()
	mockPeerSetnNextVerifierPeers.EXPECT().Peer(gomock.Any()).Return(nil).AnyTimes()
	mockPeerSetVerifierBootNode.EXPECT().Peer(gomock.Any()).Return(nil).AnyTimes()

	psm := &CsPmPeerSetManager{
		basePeers:            mockPeerSetBasePeers,
		currentVerifierPeers: mockPeerSetCurrentVerifierPeers,
		nextVerifierPeers:    mockPeerSetnNextVerifierPeers,
		verifierBootNode:     mockPeerSetVerifierBootNode,
	}

	pm := &CsProtocolManager{
		CsProtocolManagerConfig: cfg,
		peerSetManager:          psm,
	}

	pm.RemovePeer("sss")
}

func TestCsProtocolManager_GetPeer(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// mock peer set
	mockPeerSetBasePeers := NewMockAbstractPeerSet(ctrl)
	mockPeerSetCurrentVerifierPeers := NewMockAbstractPeerSet(ctrl)
	mockPeerSetnNextVerifierPeers := NewMockAbstractPeerSet(ctrl)
	mockPeerSetVerifierBootNode := NewMockAbstractPeerSet(ctrl)

	// mock peer
	mockPeerSetBasePeers.EXPECT().Peer(gomock.Any()).Return(nil).AnyTimes()
	mockPeerSetCurrentVerifierPeers.EXPECT().Peer(gomock.Any()).Return(nil).AnyTimes()
	mockPeerSetnNextVerifierPeers.EXPECT().Peer(gomock.Any()).Return(nil).AnyTimes()
	mockPeerSetVerifierBootNode.EXPECT().Peer(gomock.Any()).Return(nil).AnyTimes()

	psm := &CsPmPeerSetManager{
		basePeers:            mockPeerSetBasePeers,
		currentVerifierPeers: mockPeerSetCurrentVerifierPeers,
		nextVerifierPeers:    mockPeerSetnNextVerifierPeers,
		verifierBootNode:     mockPeerSetVerifierBootNode,
	}

	pm := &CsProtocolManager{
		peerSetManager: psm,
	}

	assert.Nil(t, pm.GetPeer("aaa"))
}

func TestCsProtocolManager_GetPeer1(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// mock peer set
	mockPeerSetBasePeers := NewMockAbstractPeerSet(ctrl)

	mockPeer := NewMockPmAbstractPeer(ctrl)
	// mock peer
	mockPeerSetBasePeers.EXPECT().Peer(gomock.Any()).Return(mockPeer).AnyTimes()

	psm := &CsPmPeerSetManager{
		basePeers: mockPeerSetBasePeers,
	}

	pm := &CsProtocolManager{
		peerSetManager: psm,
	}

	assert.NotNil(t, pm.GetPeer("aaa"))
}

func TestCsProtocolManager_GetPeer2(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// mock peer set
	mockPeerSetBasePeers := NewMockAbstractPeerSet(ctrl)
	mockPeerSetCurrentVerifierPeers := NewMockAbstractPeerSet(ctrl)

	// mock peer
	mockPeer := NewMockPmAbstractPeer(ctrl)
	mockPeerSetBasePeers.EXPECT().Peer(gomock.Any()).Return(nil).AnyTimes()
	mockPeerSetCurrentVerifierPeers.EXPECT().Peer(gomock.Any()).Return(mockPeer).AnyTimes()

	psm := &CsPmPeerSetManager{
		basePeers:            mockPeerSetBasePeers,
		currentVerifierPeers: mockPeerSetCurrentVerifierPeers,
	}

	pm := &CsProtocolManager{
		peerSetManager: psm,
	}

	assert.NotNil(t, pm.GetPeer("aaa"))
}

func TestCsProtocolManager_GetPeer3(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// mock peer set
	mockPeerSetBasePeers := NewMockAbstractPeerSet(ctrl)
	mockPeerSetCurrentVerifierPeers := NewMockAbstractPeerSet(ctrl)
	mockPeerSetnNextVerifierPeers := NewMockAbstractPeerSet(ctrl)

	mockPeer := NewMockPmAbstractPeer(ctrl)
	// mock peer
	mockPeerSetBasePeers.EXPECT().Peer(gomock.Any()).Return(nil).AnyTimes()
	mockPeerSetCurrentVerifierPeers.EXPECT().Peer(gomock.Any()).Return(nil).AnyTimes()
	mockPeerSetnNextVerifierPeers.EXPECT().Peer(gomock.Any()).Return(mockPeer).AnyTimes()

	psm := &CsPmPeerSetManager{
		basePeers:            mockPeerSetBasePeers,
		currentVerifierPeers: mockPeerSetCurrentVerifierPeers,
		nextVerifierPeers:    mockPeerSetnNextVerifierPeers,
	}

	pm := &CsProtocolManager{
		peerSetManager: psm,
	}

	assert.NotNil(t, pm.GetPeer("aaa"))
}

func TestCsProtocolManager_GetPeer4(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPeer := NewMockPmAbstractPeer(ctrl)
	// mock peer set
	mockPeerSetBasePeers := NewMockAbstractPeerSet(ctrl)
	mockPeerSetCurrentVerifierPeers := NewMockAbstractPeerSet(ctrl)
	mockPeerSetnNextVerifierPeers := NewMockAbstractPeerSet(ctrl)
	mockPeerSetVerifierBootNode := NewMockAbstractPeerSet(ctrl)

	// mock peer
	mockPeerSetBasePeers.EXPECT().Peer(gomock.Any()).Return(nil).AnyTimes()
	mockPeerSetCurrentVerifierPeers.EXPECT().Peer(gomock.Any()).Return(nil).AnyTimes()
	mockPeerSetnNextVerifierPeers.EXPECT().Peer(gomock.Any()).Return(nil).AnyTimes()
	mockPeerSetVerifierBootNode.EXPECT().Peer(gomock.Any()).Return(mockPeer).AnyTimes()

	psm := &CsPmPeerSetManager{
		basePeers:            mockPeerSetBasePeers,
		currentVerifierPeers: mockPeerSetCurrentVerifierPeers,
		nextVerifierPeers:    mockPeerSetnNextVerifierPeers,
		verifierBootNode:     mockPeerSetVerifierBootNode,
	}

	pm := &CsProtocolManager{
		peerSetManager: psm,
	}

	assert.NotNil(t, pm.GetPeer("aaa"))
}

func TestCsProtocolManager_GetPeers(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// mock peer set
	mockPeerSetBasePeers := NewMockAbstractPeerSet(ctrl)
	mockPeerSetCurrentVerifierPeers := NewMockAbstractPeerSet(ctrl)
	mockPeerSetnNextVerifierPeers := NewMockAbstractPeerSet(ctrl)
	mockPeerSetVerifierBootNode := NewMockAbstractPeerSet(ctrl)

	// mock peer
	mockPeer := NewMockPmAbstractPeer(ctrl)
	mockPeerSetBasePeers.EXPECT().GetPeers().Return(map[string]PmAbstractPeer{"1": mockPeer}).AnyTimes()

	// mock peer
	mockPeer1 := NewMockPmAbstractPeer(ctrl)
	mockPeerSetCurrentVerifierPeers.EXPECT().GetPeers().Return(map[string]PmAbstractPeer{"2": mockPeer1}).AnyTimes()

	// mock peer
	mockPeer2 := NewMockPmAbstractPeer(ctrl)
	mockPeerSetnNextVerifierPeers.EXPECT().GetPeers().Return(map[string]PmAbstractPeer{"3": mockPeer2}).AnyTimes()

	// mock peer
	mockPeer3 := NewMockPmAbstractPeer(ctrl)
	mockPeerSetVerifierBootNode.EXPECT().GetPeers().Return(map[string]PmAbstractPeer{"4": mockPeer3}).AnyTimes()

	psm := &CsPmPeerSetManager{
		basePeers:            mockPeerSetBasePeers,
		currentVerifierPeers: mockPeerSetCurrentVerifierPeers,
		nextVerifierPeers:    mockPeerSetnNextVerifierPeers,
		verifierBootNode:     mockPeerSetVerifierBootNode,
	}

	pm := &CsProtocolManager{
		peerSetManager: psm,
	}

	peers := pm.GetPeers()

	assert.Equal(t, 4, len(peers))
}

func Test_mergePeersTo(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPeer1 := NewMockPmAbstractPeer(ctrl)

	fm := map[string]PmAbstractPeer{"1": mockPeer1}

	tm := make(map[string]PmAbstractPeer)

	mergePeersTo(fm, tm)

	assert.NotNil(t, tm["1"])
}

func TestCsProtocolManager_Protocols(t *testing.T) {
	pm := &CsProtocolManager{}
	assert.Equal(t, 1, len(pm.Protocols()))
	assert.Equal(t, 1, len(pm.Protocols()))
}

func TestCsProtocolManager_getCsProtocol(t *testing.T) {
	pm := &CsProtocolManager{}
	_ = os.Setenv("boots_env", "local")
	assert.Equal(t, chain_config.AppName+"_cs_local", pm.getCsProtocol().Name)
}

func TestCsProtocolManager_getCsProtocol1(t *testing.T) {
	pm := &CsProtocolManager{}
	_ = os.Setenv("boots_env", "mercury")
	assert.Equal(t, chain_config.AppName+"_cs", pm.getCsProtocol().Name)
}

func TestCsProtocolManager_getCsProtocol2(t *testing.T) {
	pm := &CsProtocolManager{}
	_ = os.Setenv("boots_env", "test")
	assert.Equal(t, chain_config.AppName+"_cs_test", pm.getCsProtocol().Name)
}

func TestCsProtocolManager_Start(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// node conf
	mockNodeConf := NewMockNodeConf(ctrl)
	mockNodeConf.EXPECT().GetNodeType().Return(chain_config.NodeTypeOfVerifier).AnyTimes()
	mockNodeConf.EXPECT().GetNodeName().Return("aaaa").AnyTimes()

	// chain
	mockChain := NewMockChain(ctrl)

	// p2p server
	mockP2pServer := NewMockP2PServer(ctrl)

	// verifier
	mockVerifiersReader := NewMockVerifiersReader(ctrl)

	// pbft node
	mockPbftNode := NewMockPbftNode(ctrl)

	// PbftSigner
	mockPbftSigner := NewMockPbftSigner(ctrl)

	cfg := &CsProtocolManagerConfig{
		ChainConfig:     *chain_config.GetChainConfig(),
		Chain:           mockChain,
		P2PServer:       mockP2pServer,
		NodeConf:        mockNodeConf,
		VerifiersReader: mockVerifiersReader,
		PbftNode:        mockPbftNode,
		MsgSigner:       mockPbftSigner,
	}

	// mock peer set
	mockPeerSetBasePeers := NewMockAbstractPeerSet(ctrl)
	mockPeerSetCurrentVerifierPeers := NewMockAbstractPeerSet(ctrl)
	mockPeerSetnNextVerifierPeers := NewMockAbstractPeerSet(ctrl)
	mockPeerSetVerifierBootNode := NewMockAbstractPeerSet(ctrl)

	psm := &CsPmPeerSetManager{
		basePeers:            mockPeerSetBasePeers,
		currentVerifierPeers: mockPeerSetCurrentVerifierPeers,
		nextVerifierPeers:    mockPeerSetnNextVerifierPeers,
		verifierBootNode:     mockPeerSetVerifierBootNode,
	}

	pm := &CsProtocolManager{
		BaseProtocolManager: BaseProtocolManager{
			msgHandlers: make(map[uint64]func(msg p2p.Msg, p PmAbstractPeer) error),
		},
		CsProtocolManagerConfig: cfg,
		maxPeers:                P2PMaxPeerCount,
		verifierBootNodes:       chain_config.VerifierBootNodes,
		stop:                    make(chan struct{}),
		peerSetManager:          psm,
	}

	_ = pm.Start()

	time.Sleep(30 * time.Millisecond)
}

//func TestCsProtocolManager_logCurPeersInfo(t *testing.T) {
//	pm := &CsProtocolManager{stop: make(chan struct{}),}
//
//	go pm.logCurPeersInfo()
//
//	pm.stop <- struct{}{}
//}

func TestCsProtocolManager_Stop(t *testing.T) {
	pm := &CsProtocolManager{
		BaseProtocolManager: BaseProtocolManager{
			msgHandlers: make(map[uint64]func(msg p2p.Msg, p PmAbstractPeer) error),
		},
	}

	pm.Stop()

	pm.stop = make(chan struct{})

	pm.Stop()
}

func TestCsProtocolManager_BroadcastMsg(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPeer := NewMockPmAbstractPeer(ctrl)
	mockPeer.EXPECT().SendMsg(gomock.Eq(uint64(1)), gomock.Any()).Return(errors.New("aaa"))
	mockPeer.EXPECT().SendMsg(gomock.Eq(uint64(2)), gomock.Any()).Return(nil)
	mockPeer.EXPECT().NodeName().Return("11212")

	// mock peer set
	mockPeerSetCurrentVerifierPeers := NewMockAbstractPeerSet(ctrl)

	mockPeerSetCurrentVerifierPeers.EXPECT().GetPeers().Return(map[string]PmAbstractPeer{"1": mockPeer}).Times(2)

	psm := &CsPmPeerSetManager{
		currentVerifierPeers: mockPeerSetCurrentVerifierPeers,
	}

	pm := &CsProtocolManager{
		peerSetManager: psm,
	}

	pm.BroadcastMsg(uint64(1), uint64(2333))
	pm.BroadcastMsg(uint64(2), uint64(2333))
}

func TestCsProtocolManager_BroadcastMsgToTargetVerifiers(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPeer := NewMockPmAbstractPeer(ctrl)
	mockPeer.EXPECT().SendMsg(gomock.Eq(uint64(1)), gomock.Any()).Return(errors.New("aaa")).AnyTimes()
	mockPeer.EXPECT().SendMsg(gomock.Eq(uint64(2)), gomock.Any()).Return(nil).AnyTimes()
	mockPeer.EXPECT().NodeName().Return("11212").AnyTimes()
	mockPeer.EXPECT().RemoteVerifierAddress().Return(common.HexToAddress("11")).AnyTimes()

	// mock peer set
	mockPeerSetCurrentVerifierPeers := NewMockAbstractPeerSet(ctrl)

	mockPeerSetCurrentVerifierPeers.EXPECT().GetPeers().Return(map[string]PmAbstractPeer{"1": mockPeer}).Times(2)
	mockPeerSetCurrentVerifierPeers.EXPECT().Len().Return(1).Times(2)

	psm := &CsPmPeerSetManager{
		currentVerifierPeers: mockPeerSetCurrentVerifierPeers,
	}

	pm := &CsProtocolManager{
		peerSetManager: psm,
	}

	pm.BroadcastMsgToTargetVerifiers(1, []common.Address{common.HexToAddress("11")}, uint64(33))
	pm.BroadcastMsgToTargetVerifiers(2, []common.Address{common.HexToAddress("11")}, uint64(33))

}

func TestCsProtocolManager_SendFetchBlockMsg(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPeer := NewMockPmAbstractPeer(ctrl)
	mockPeer.EXPECT().SendMsg(gomock.Eq(uint64(1)), gomock.Any()).Return(errors.New("aaa")).AnyTimes()
	mockPeer.EXPECT().SendMsg(gomock.Eq(uint64(2)), gomock.Any()).Return(nil).AnyTimes()
	mockPeer.EXPECT().NodeName().Return("11212").AnyTimes()
	mockPeer.EXPECT().RemoteVerifierAddress().Return(common.HexToAddress("11")).AnyTimes()

	// mock peer set
	mockPeerSetCurrentVerifierPeers := NewMockAbstractPeerSet(ctrl)

	mockPeerSetCurrentVerifierPeers.EXPECT().GetPeers().Return(map[string]PmAbstractPeer{"1": mockPeer}).Times(2)

	psm := &CsPmPeerSetManager{
		currentVerifierPeers: mockPeerSetCurrentVerifierPeers,
	}

	pm := &CsProtocolManager{
		peerSetManager: psm,
	}

	assert.EqualError(t, pm.SendFetchBlockMsg(1, common.HexToAddress("11"), nil), "aaa")

	assert.Nil(t, pm.SendFetchBlockMsg(2, common.HexToAddress("11"), nil))

	mockPeerSetCurrentVerifierPeers.EXPECT().GetPeers().Return(nil)

	assert.EqualError(t, pm.SendFetchBlockMsg(1, common.HexToAddress("11"), nil), "no verifier peer for fetcher")

}

func TestCsProtocolManager_ChangeVerifiers(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	address := common.HexToAddress("aaa")

	mVR := NewMockVerifiersReader(ctrl)
	mPbftS := NewMockPbftSigner(ctrl)
	mNodeConf := NewMockNodeConf(ctrl)
	mPn := NewMockPbftNode(ctrl)

	pm := &CsProtocolManager{CsProtocolManagerConfig: &CsProtocolManagerConfig{VerifiersReader: mVR, MsgSigner: mPbftS, NodeConf: mNodeConf, PbftNode: mPn}}

	mVR.EXPECT().CurrentVerifiers().Return([]common.Address{address})
	mVR.EXPECT().NextVerifiers().Return([]common.Address{address})
	mVR.EXPECT().ShouldChangeVerifier().Return(false)
	mPbftS.EXPECT().GetAddress().Return(address)
	mVR.EXPECT().NextVerifiers().Return([]common.Address{address})
	mNodeConf.EXPECT().GetNodeType().Return(verifier)
	mPn.EXPECT().ChangePrimary(gomock.Any())

	assert.Panics(t, func() {
		pm.ChangeVerifiers()
	})
}

func TestCsProtocolManager_GetCurrentConnectPeers(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockNodeConf := NewMockNodeConf(ctrl)
	mockNodeConf.EXPECT().GetNodeType().Return(chain_config.NodeTypeOfVerifier)

	mockPeer := NewMockPmAbstractPeer(ctrl)
	mockPeer.EXPECT().RemoteVerifierAddress().Return(common.HexToAddress("sss"))

	mockPeerSetCurrentVerifierPeers := NewMockAbstractPeerSet(ctrl)

	mockPeerSetCurrentVerifierPeers.EXPECT().GetPeers().Return(map[string]PmAbstractPeer{"1": mockPeer})

	psm := &CsPmPeerSetManager{
		currentVerifierPeers: mockPeerSetCurrentVerifierPeers,
	}

	pm := &CsProtocolManager{peerSetManager: psm, CsProtocolManagerConfig: &CsProtocolManagerConfig{NodeConf: mockNodeConf}}

	rM := pm.GetCurrentConnectPeers()

	assert.Equal(t, true, rM["1"].IsEqual(common.HexToAddress("sss")))
}

func TestCsProtocolManager_GetCurrentConnectPeers2(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockNodeConf := NewMockNodeConf(ctrl)
	mockNodeConf.EXPECT().GetNodeType().Return(chain_config.NodeTypeOfNormal)

	pm := &CsProtocolManager{CsProtocolManagerConfig: &CsProtocolManagerConfig{NodeConf: mockNodeConf}}

	pm.GetCurrentConnectPeers()
}

func TestCsProtocolManager_GetVerifierBootNode(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// mock peer set
	mockPeerSet := NewMockAbstractPeerSet(ctrl)

	mockPeerSet.EXPECT().GetPeers().Return(nil)

	psm := &CsPmPeerSetManager{
		verifierBootNode: mockPeerSet,
	}

	pm := &CsProtocolManager{peerSetManager: psm}

	pm.GetVerifierBootNode()
}

func TestCsProtocolManager_GetNextVerifierPeers(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// mock peer set
	mockPeerSet := NewMockAbstractPeerSet(ctrl)

	mockPeerSet.EXPECT().GetPeers().Return(nil)

	psm := &CsPmPeerSetManager{
		nextVerifierPeers: mockPeerSet,
	}

	pm := &CsProtocolManager{peerSetManager: psm}

	pm.GetNextVerifierPeers()
}

func TestCsProtocolManager_SelfIsBootNode(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockNodeConf := NewMockNodeConf(ctrl)
	mockNodeConf.EXPECT().GetNodeType().Return(chain_config.NodeTypeOfNormal)

	pm := &CsProtocolManager{CsProtocolManagerConfig: &CsProtocolManagerConfig{NodeConf: mockNodeConf}}

	assert.Equal(t, false, pm.SelfIsBootNode())
}

func TestCsProtocolManager_SelfIsBootNode1(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockNodeConf := NewMockNodeConf(ctrl)
	mockNodeConf.EXPECT().GetNodeType().Return(chain_config.NodeTypeOfVerifierBoot)

	mockP2PServer := NewMockP2PServer(ctrl)

	var pyRecord, _ = hex.DecodeString("f884b8407098ad865b00a582051940cb9cf36836572411a47278783077011599ed5cd16b76f2635f4e234738f30813a89eb9137e3e3df5266e3a1f11df72ecf1145ccb9c01826964827634826970847f00000189736563703235366b31a103ca634cae0d49acb401d8a4c6b6fe8c55b70d115bf400769cc1400f3258cd31388375647082765f")

	var r enr.Record
	if err := rlp.DecodeBytes(pyRecord, &r); err != nil {
		t.Fatalf("can't decode: %v", err)
	}
	n, err := enode.New(enode.ValidSchemes, &r)
	if err != nil {
		t.Fatalf("can't verify record: %v", err)
	}

	mockP2PServer.EXPECT().Self().Return(n)

	pm := &CsProtocolManager{CsProtocolManagerConfig: &CsProtocolManagerConfig{NodeConf: mockNodeConf, P2PServer: mockP2PServer}, verifierBootNodes: []*enode.Node{n}}

	assert.Equal(t, true, pm.SelfIsBootNode())
}

func TestCsProtocolManager_GetSelfNode(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockP2PServer := NewMockP2PServer(ctrl)

	var pyRecord, _ = hex.DecodeString("f884b8407098ad865b00a582051940cb9cf36836572411a47278783077011599ed5cd16b76f2635f4e234738f30813a89eb9137e3e3df5266e3a1f11df72ecf1145ccb9c01826964827634826970847f00000189736563703235366b31a103ca634cae0d49acb401d8a4c6b6fe8c55b70d115bf400769cc1400f3258cd31388375647082765f")

	var r enr.Record
	if err := rlp.DecodeBytes(pyRecord, &r); err != nil {
		t.Fatalf("can't decode: %v", err)
	}
	n, err := enode.New(enode.ValidSchemes, &r)
	if err != nil {
		t.Fatalf("can't verify record: %v", err)
	}

	mockP2PServer.EXPECT().Self().Return(n)

	pm := &CsProtocolManager{CsProtocolManagerConfig: &CsProtocolManagerConfig{P2PServer: mockP2PServer}, verifierBootNodes: []*enode.Node{n}}

	assert.Equal(t, n.ID().String(), pm.GetSelfNode().ID().String())
}

func TestCsProtocolManager_connectPeer(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockP2PServer := NewMockP2PServer(ctrl)

	var pyRecord, _ = hex.DecodeString("f884b8407098ad865b00a582051940cb9cf36836572411a47278783077011599ed5cd16b76f2635f4e234738f30813a89eb9137e3e3df5266e3a1f11df72ecf1145ccb9c01826964827634826970847f00000189736563703235366b31a103ca634cae0d49acb401d8a4c6b6fe8c55b70d115bf400769cc1400f3258cd31388375647082765f")

	var r enr.Record
	if err := rlp.DecodeBytes(pyRecord, &r); err != nil {
		t.Fatalf("can't decode: %v", err)
	}
	n, err := enode.New(enode.ValidSchemes, &r)
	if err != nil {
		t.Fatalf("can't verify record: %v", err)
	}

	mockP2PServer.EXPECT().AddPeer(gomock.Any()).Return()

	pm := &CsProtocolManager{CsProtocolManagerConfig: &CsProtocolManagerConfig{P2PServer: mockP2PServer}}

	pm.ConnectPeer(n)
}

func TestCsProtocolManager_MatchCurrentVerifiersToNext(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	address := common.HexToAddress("aaa")

	mVR := NewMockVerifiersReader(ctrl)
	mPbftS := NewMockPbftSigner(ctrl)

	pm := &CsProtocolManager{CsProtocolManagerConfig: &CsProtocolManagerConfig{VerifiersReader: mVR, MsgSigner: mPbftS}}

	mVR.EXPECT().NextVerifiers().Return([]common.Address{address})
	mVR.EXPECT().ShouldChangeVerifier().Return(true)
	mPbftS.EXPECT().GetAddress().Return(address)

	pm.MatchCurrentVerifiersToNext()
}

func TestCsProtocolManager_MatchCurrentVerifiersToNext1(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	address := common.HexToAddress("aaa")

	mVR := NewMockVerifiersReader(ctrl)
	mPbftS := NewMockPbftSigner(ctrl)

	mPs := NewMockAbstractPeerSet(ctrl)

	psm := &CsPmPeerSetManager{nextVerifierPeers: mPs}

	pm := &CsProtocolManager{CsProtocolManagerConfig: &CsProtocolManagerConfig{VerifiersReader: mVR, MsgSigner: mPbftS}, peerSetManager: psm}

	mVR.EXPECT().NextVerifiers().Return([]common.Address{address}).AnyTimes()
	mVR.EXPECT().ShouldChangeVerifier().Return(false)
	mPbftS.EXPECT().GetAddress().Return(address)
	mPs.EXPECT().Len().Return(2)
	totalVerifier = 3

	pm.MatchCurrentVerifiersToNext()
}

func TestCsProtocolManager_MatchCurrentVerifiersToNext2(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	address := common.HexToAddress("aaa")

	mVR := NewMockVerifiersReader(ctrl)
	mPbftS := NewMockPbftSigner(ctrl)

	mPeerSet1 := NewMockAbstractPeerSet(ctrl)
	mPeerSet2 := NewMockAbstractPeerSet(ctrl)
	mPeerSet3 := NewMockAbstractPeerSet(ctrl)

	psm := &CsPmPeerSetManager{
		basePeers:            mPeerSet1,
		currentVerifierPeers: mPeerSet2,
		nextVerifierPeers:    mPeerSet3,
	}

	pm := &CsProtocolManager{CsProtocolManagerConfig: &CsProtocolManagerConfig{VerifiersReader: mVR, MsgSigner: mPbftS}, peerSetManager: psm}

	mVR.EXPECT().NextVerifiers().Return([]common.Address{address}).AnyTimes()
	mVR.EXPECT().ShouldChangeVerifier().Return(false)
	mPbftS.EXPECT().GetAddress().Return(address)
	mPeerSet3.EXPECT().Len().Return(2)
	totalVerifier = 4

	mPeerSet1.EXPECT().GetPeers().Return(make(map[string]PmAbstractPeer))
	mPeerSet2.EXPECT().GetPeers().Return(make(map[string]PmAbstractPeer))

	pm.MatchCurrentVerifiersToNext()
}

func TestCsProtocolManager_pickNextVerifierFromPs(t *testing.T) {
	addr1 := common.StringToAddress("aaavvv")
	addr2 := common.StringToAddress("aaabbb")
	addr3 := common.StringToAddress("fdsf")
	nextVs := []common.Address{addr1, addr2, addr3}

	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mPeer1 := NewMockPmAbstractPeer(ctrl)
	mPeerSet1 := NewMockAbstractPeerSet(ctrl)
	mPeerSet2 := NewMockAbstractPeerSet(ctrl)

	psm := &CsPmPeerSetManager{
		basePeers:         mPeerSet1,
		nextVerifierPeers: mPeerSet2,
	}

	pm := &CsProtocolManager{peerSetManager: psm}

	mPeer1.EXPECT().RemoteVerifierAddress().Return(addr1).AnyTimes()
	mPeerSet1.EXPECT().GetPeers().Return(map[string]PmAbstractPeer{"1": mPeer1})
	mPeer1.EXPECT().ID().Return("1")
	mPeerSet2.EXPECT().Peer(gomock.Eq("1")).Return(mPeer1)
	mPeer1.EXPECT().ID().Return("1")
	mPeerSet1.EXPECT().RemovePeer(gomock.Eq("1")).Return(errors.New("fdsfds"))
	pm.pickNextVerifierFromPs(nextVs)

}

func TestCsProtocolManager_pickNextVerifierFromPs1(t *testing.T) {
	addr1 := common.StringToAddress("aaavvv")
	addr2 := common.StringToAddress("aaabbb")
	addr3 := common.StringToAddress("fdsf")
	nextVs := []common.Address{addr1, addr2, addr3}

	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mPeer1 := NewMockPmAbstractPeer(ctrl)
	mPeerSet1 := NewMockAbstractPeerSet(ctrl)
	mPeerSet2 := NewMockAbstractPeerSet(ctrl)

	psm := &CsPmPeerSetManager{
		basePeers:         mPeerSet1,
		nextVerifierPeers: mPeerSet2,
	}

	pm := &CsProtocolManager{peerSetManager: psm}

	mPeer1.EXPECT().RemoteVerifierAddress().Return(addr1).AnyTimes()
	mPeerSet1.EXPECT().GetPeers().Return(map[string]PmAbstractPeer{"1": mPeer1})
	mPeer1.EXPECT().ID().Return("1")
	mPeerSet2.EXPECT().Peer(gomock.Eq("1")).Return(nil)

	mPeerSet2.EXPECT().AddPeer(gomock.Any()).Return(errors.New("sdsd"))

	pm.pickNextVerifierFromPs(nextVs)

}

func TestCsProtocolManager_pickNextVerifierFromPs2(t *testing.T) {
	addr1 := common.StringToAddress("aaavvv")
	addr2 := common.StringToAddress("aaabbb")
	addr3 := common.StringToAddress("fdsf")
	nextVs := []common.Address{addr1, addr2, addr3}

	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mPeer1 := NewMockPmAbstractPeer(ctrl)
	mPeerSet1 := NewMockAbstractPeerSet(ctrl)
	mPeerSet2 := NewMockAbstractPeerSet(ctrl)

	psm := &CsPmPeerSetManager{
		basePeers:         mPeerSet1,
		nextVerifierPeers: mPeerSet2,
	}

	pm := &CsProtocolManager{peerSetManager: psm}

	mPeer1.EXPECT().RemoteVerifierAddress().Return(addr1).AnyTimes()
	mPeerSet1.EXPECT().GetPeers().Return(map[string]PmAbstractPeer{"1": mPeer1})
	mPeer1.EXPECT().ID().Return("1")
	mPeerSet2.EXPECT().Peer(gomock.Eq("1")).Return(nil)
	mPeerSet2.EXPECT().AddPeer(gomock.Any()).Return(nil)

	mPeer1.EXPECT().ID().Return("1")
	mPeerSet1.EXPECT().RemovePeer(gomock.Any()).Return(errors.New("sdsd"))

	pm.pickNextVerifierFromPs(nextVs)
}

func TestCsProtocolManager_pickNextVerifierFromCps(t *testing.T) {
	addr1 := common.StringToAddress("aaavvv")
	addr2 := common.StringToAddress("aaabbb")
	nextVs := []common.Address{addr1, addr2}

	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mPeer1 := NewMockPmAbstractPeer(ctrl)
	mPeer2 := NewMockPmAbstractPeer(ctrl)
	mPeer3 := NewMockPmAbstractPeer(ctrl)
	mPeerSet1 := NewMockAbstractPeerSet(ctrl)
	mPeerSet2 := NewMockAbstractPeerSet(ctrl)

	psm := &CsPmPeerSetManager{
		currentVerifierPeers: mPeerSet1,
		nextVerifierPeers:    mPeerSet2,
	}

	pm := &CsProtocolManager{peerSetManager: psm}

	mPeerSet1.EXPECT().GetPeers().Return(map[string]PmAbstractPeer{"1": mPeer1, "2": mPeer2, "3": mPeer3})
	mPeer1.EXPECT().RemoteVerifierAddress().Return(addr1).AnyTimes()
	mPeer2.EXPECT().RemoteVerifierAddress().Return(addr2).AnyTimes()
	mPeer3.EXPECT().RemoteVerifierAddress().Return(common.StringToAddress("fdsf")).AnyTimes()
	mPeer1.EXPECT().ID().Return("1")
	mPeer2.EXPECT().ID().Return("2")
	mPeerSet2.EXPECT().Peer(gomock.Eq("1")).Return(mPeer1)
	mPeerSet2.EXPECT().Peer(gomock.Eq("2")).Return(nil)
	mPeerSet2.EXPECT().AddPeer(gomock.Any()).Return(errors.New("dsadsa"))

	pm.pickNextVerifierFromCps(nextVs)
}

func TestCsProtocolManager_isVerifierBootNode(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var pyRecord, _ = hex.DecodeString("f884b8407098ad865b00a582051940cb9cf36836572411a47278783077011599ed5cd16b76f2635f4e234738f30813a89eb9137e3e3df5266e3a1f11df72ecf1145ccb9c01826964827634826970847f00000189736563703235366b31a103ca634cae0d49acb401d8a4c6b6fe8c55b70d115bf400769cc1400f3258cd31388375647082765f")

	var r enr.Record
	if err := rlp.DecodeBytes(pyRecord, &r); err != nil {
		t.Fatalf("can't decode: %v", err)
	}
	n, err := enode.New(enode.ValidSchemes, &r)
	if err != nil {
		t.Fatalf("can't verify record: %v", err)
	}

	mockPeer := NewMockPmAbstractPeer(ctrl)
	mockPeer.EXPECT().ID().Return(n.ID().String())

	pm := &CsProtocolManager{verifierBootNodes: []*enode.Node{n}}

	assert.Equal(t, true, pm.isVerifierBootNode(mockPeer))

	mockPeer.EXPECT().ID().Return("asdasddas")
	assert.Equal(t, false, pm.isVerifierBootNode(mockPeer))
}

func TestCsProtocolManager_isNextVerifierNode(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPeer := NewMockPmAbstractPeer(ctrl)

	mockVerifiersReader := NewMockVerifiersReader(ctrl)

	pm := &CsProtocolManager{CsProtocolManagerConfig: &CsProtocolManagerConfig{VerifiersReader: mockVerifiersReader}}

	mockVerifiersReader.EXPECT().ShouldChangeVerifier().Return(true)
	assert.Equal(t, false, pm.isNextVerifierNode(mockPeer))

	address := common.HexToAddress("aaa")
	mockVerifiersReader.EXPECT().ShouldChangeVerifier().Return(false)
	mockVerifiersReader.EXPECT().NextVerifiers().Return([]common.Address{address})
	mockPeer.EXPECT().RemoteVerifierAddress().Return(address)
	assert.Equal(t, true, pm.isNextVerifierNode(mockPeer))

	mockVerifiersReader.EXPECT().ShouldChangeVerifier().Return(false)
	mockVerifiersReader.EXPECT().NextVerifiers().Return([]common.Address{address})
	mockPeer.EXPECT().RemoteVerifierAddress().Return(common.HexToAddress("aab"))
	assert.Equal(t, false, pm.isNextVerifierNode(mockPeer))

}

func TestCsProtocolManager_isCurrentVerifierNode(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPeer := NewMockPmAbstractPeer(ctrl)

	mockVerifiersReader := NewMockVerifiersReader(ctrl)

	pm := &CsProtocolManager{CsProtocolManagerConfig: &CsProtocolManagerConfig{VerifiersReader: mockVerifiersReader}}

	address := common.HexToAddress("aaa")

	mockVerifiersReader.EXPECT().ShouldChangeVerifier().Return(false)
	mockVerifiersReader.EXPECT().CurrentVerifiers().Return([]common.Address{address})
	mockPeer.EXPECT().RemoteVerifierAddress().Return(common.HexToAddress("aab"))
	assert.Equal(t, false, pm.isCurrentVerifierNode(mockPeer))

	mockVerifiersReader.EXPECT().ShouldChangeVerifier().Return(false)
	mockVerifiersReader.EXPECT().CurrentVerifiers().Return([]common.Address{address})
	mockPeer.EXPECT().RemoteVerifierAddress().Return(address)
	assert.Equal(t, true, pm.isCurrentVerifierNode(mockPeer))

	mockVerifiersReader.EXPECT().ShouldChangeVerifier().Return(true)
	mockVerifiersReader.EXPECT().NextVerifiers().Return([]common.Address{address})
	mockPeer.EXPECT().RemoteVerifierAddress().Return(common.HexToAddress("aab"))
	assert.Equal(t, false, pm.isCurrentVerifierNode(mockPeer))

	mockVerifiersReader.EXPECT().ShouldChangeVerifier().Return(true)
	mockVerifiersReader.EXPECT().NextVerifiers().Return([]common.Address{address})
	mockPeer.EXPECT().RemoteVerifierAddress().Return(address)
	assert.Equal(t, true, pm.isCurrentVerifierNode(mockPeer))
}

func TestCsProtocolManager_selfIsCurrentVerifier(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	address := common.HexToAddress("aaa")

	mVR := NewMockVerifiersReader(ctrl)
	mPbftS := NewMockPbftSigner(ctrl)

	pm := &CsProtocolManager{CsProtocolManagerConfig: &CsProtocolManagerConfig{VerifiersReader: mVR, MsgSigner: mPbftS}}

	mVR.EXPECT().ShouldChangeVerifier().Return(true)
	mVR.EXPECT().CurrentVerifiers().Return([]common.Address{})
	mVR.EXPECT().NextVerifiers().Return([]common.Address{address})
	mPbftS.EXPECT().GetAddress().Return(address)
	assert.Equal(t, true, pm.SelfIsCurrentVerifier())

	mVR.EXPECT().ShouldChangeVerifier().Return(true)
	mVR.EXPECT().CurrentVerifiers().Return([]common.Address{})
	mVR.EXPECT().NextVerifiers().Return([]common.Address{address})
	mPbftS.EXPECT().GetAddress().Return(common.StringToAddress("dasds"))
	assert.Equal(t, false, pm.SelfIsCurrentVerifier())
}

func TestCsProtocolManager_selfIsVerifierBootNode(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	address := common.HexToAddress("aaa")

	mPbftS := NewMockPbftSigner(ctrl)
	pm := &CsProtocolManager{CsProtocolManagerConfig: &CsProtocolManagerConfig{MsgSigner: mPbftS}}
	chain_config.VerBootNodeAddress = append(chain_config.VerBootNodeAddress, address)

	mPbftS.EXPECT().GetAddress().Return(address)
	assert.Equal(t, true, pm.selfIsVerifierBootNode())

	mPbftS.EXPECT().GetAddress().Return(common.StringToAddress("dasds"))
	assert.Equal(t, false, pm.selfIsVerifierBootNode())
}

func TestCsProtocolManager_selfIsNextVerifier(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	address := common.HexToAddress("aaa")

	mVR := NewMockVerifiersReader(ctrl)
	mPbftS := NewMockPbftSigner(ctrl)

	pm := &CsProtocolManager{CsProtocolManagerConfig: &CsProtocolManagerConfig{VerifiersReader: mVR, MsgSigner: mPbftS}}

	mVR.EXPECT().NextVerifiers().Return([]common.Address{address})
	mVR.EXPECT().ShouldChangeVerifier().Return(false)
	mPbftS.EXPECT().GetAddress().Return(address)
	assert.Equal(t, true, pm.SelfIsNextVerifier())

	mVR.EXPECT().NextVerifiers().Return([]common.Address{address})
	mVR.EXPECT().ShouldChangeVerifier().Return(false)
	mPbftS.EXPECT().GetAddress().Return(address)
	assert.Equal(t, true, pm.SelfIsNextVerifier())

	mVR.EXPECT().NextVerifiers().Return([]common.Address{address})
	mVR.EXPECT().ShouldChangeVerifier().Return(true)
	mPbftS.EXPECT().GetAddress().Return(address)
	assert.Equal(t, false, pm.SelfIsNextVerifier())

	mVR.EXPECT().NextVerifiers().Return([]common.Address{address})
	mVR.EXPECT().ShouldChangeVerifier().Return(false)
	mPbftS.EXPECT().GetAddress().Return(address)
	assert.Equal(t, true, pm.SelfIsNextVerifier())

}

func TestCsProtocolManager_checkConnCount(t *testing.T) {
	pm1 := &CsProtocolManager{}
	pm1.pmType.Store(5)
	assert.Equal(t, false, pm1.checkConnCount())

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mps1 := NewMockAbstractPeerSet(ctrl)
	mps2 := NewMockAbstractPeerSet(ctrl)
	mps3 := NewMockAbstractPeerSet(ctrl)

	psm := &CsPmPeerSetManager{
		basePeers:            mps1,
		currentVerifierPeers: mps2,
		nextVerifierPeers:    mps3,
	}

	pm := &CsProtocolManager{peerSetManager: psm}

	pm.pmType.Store(base)
	pm.maxPeers = 2
	mps1.EXPECT().Len().Return(3)
	assert.Equal(t, true, pm.checkConnCount())

	pm.pmType.Store(verifier)
	pm.maxPeers = 4
	PbftMaxPeerCount = 2
	mps2.EXPECT().Len().Return(2)
	mps1.EXPECT().Len().Return(2)
	assert.Equal(t, true, pm.checkConnCount())

	pm.pmType.Store(boot)
	pm.maxPeers = 4
	PbftMaxPeerCount = 2
	mps2.EXPECT().Len().Return(2)
	mps1.EXPECT().Len().Return(2)
	assert.Equal(t, true, pm.checkConnCount())
}

func TestCsProtocolManager_handle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mPeer := NewMockPmAbstractPeer(ctrl)
	mps1 := NewMockAbstractPeerSet(ctrl)
	mps2 := NewMockAbstractPeerSet(ctrl)
	mps3 := NewMockAbstractPeerSet(ctrl)

	psm := &CsPmPeerSetManager{
		basePeers:            mps1,
		currentVerifierPeers: mps2,
		nextVerifierPeers:    mps3,
	}

	pm := &CsProtocolManager{peerSetManager: psm}

	pm.pmType.Store(verifier)
	pm.maxPeers = 4
	PbftMaxPeerCount = 2
	mps2.EXPECT().Len().Return(2)
	mps1.EXPECT().Len().Return(2)

	assert.EqualError(t, pm.handle(mPeer), "too many peers")
}

func TestCsProtocolManager_handle1(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mPeer := NewMockPmAbstractPeer(ctrl)
	mps1 := NewMockAbstractPeerSet(ctrl)
	mps2 := NewMockAbstractPeerSet(ctrl)
	mps3 := NewMockAbstractPeerSet(ctrl)

	psm := &CsPmPeerSetManager{
		basePeers:            mps1,
		currentVerifierPeers: mps2,
		nextVerifierPeers:    mps3,
	}

	mChain := NewMockChain(ctrl)
	mNodeConf := NewMockNodeConf(ctrl)
	mSigner := NewMockPbftSigner(ctrl)
	mP2PServer := NewMockP2PServer(ctrl)

	pm := &CsProtocolManager{
		CsProtocolManagerConfig: &CsProtocolManagerConfig{
			ChainConfig: *chain_config.GetChainConfig(),
			Chain:       mChain,
			NodeConf:    mNodeConf,
			MsgSigner:   mSigner,
			P2PServer:   mP2PServer,
		},
		peerSetManager: psm,
	}

	pm.pmType.Store(verifier)
	pm.maxPeers = 7
	PbftMaxPeerCount = 3
	mps2.EXPECT().Len().Return(2)
	//mps1.EXPECT().Len().Return(2)

	var pyRecord, _ = hex.DecodeString("f884b8407098ad865b00a582051940cb9cf36836572411a47278783077011599ed5cd16b76f2635f4e234738f30813a89eb9137e3e3df5266e3a1f11df72ecf1145ccb9c01826964827634826970847f00000189736563703235366b31a103ca634cae0d49acb401d8a4c6b6fe8c55b70d115bf400769cc1400f3258cd31388375647082765f")

	var r enr.Record
	if err := rlp.DecodeBytes(pyRecord, &r); err != nil {
		t.Fatalf("can't decode: %v", err)
	}
	n, err := enode.New(enode.ValidSchemes, &r)
	if err != nil {
		t.Fatalf("can't verify record: %v", err)
	}

	block := model.NewBlock(model.NewHeader(11, 101, common.HexToHash("ss"), common.HexToHash("fdfs"), common.StringToDiff("0x22"), big.NewInt(111), common.StringToAddress("fdsfds"), common.EncodeNonce(33)), nil, nil)

	// case 2
	mChain.EXPECT().GetBlockByNumber(gomock.Eq(uint64(0))).Return(block)

	//  send
	mChain.EXPECT().CurrentBlock().Return(block)
	mNodeConf.EXPECT().GetNodeType().Return(chain_config.NodeTypeOfNormal).Times(2)
	mNodeConf.EXPECT().GetNodeName().Return("dsadsad")
	mP2PServer.EXPECT().Self().Return(n)
	mPeer.EXPECT().SendMsg(gomock.Any(), gomock.Any()).Return(nil)

	hsData := HandShakeData{
		ProtocolVersion:    11,
		ChainID:            big.NewInt(2),
		NetworkId:          pm.ChainConfig.NetworkID,
		CurrentBlock:       common.HexToHash("aaa"),
		CurrentBlockHeight: 64,
		GenesisBlock:       block.Hash(),
		NodeType:           chain_config.NodeTypeOfVerifier,
		NodeName:           "test",
		RawUrl:             n.String(),
	}

	statusData := &StatusData{HandShakeData: hsData}

	hash := statusData.DataHash()

	assert.Equal(t, true, !hash.IsEmpty())

	account := tests.AccFactory.GenAccount()

	sign, err := account.SignHash(statusData.DataHash().Bytes())

	assert.NoError(t, err)

	assert.Equal(t, true, len(sign) > 0)

	statusData.Sign = sign

	statusData.PubKey = crypto.CompressPubkey(&account.Pk.PublicKey)

	size, r1, err := rlp.EncodeToReader(statusData)
	assert.NoError(t, err)

	msg := p2p.Msg{Code: StatusMsg, Size: uint32(size), Payload: r1}

	// read
	mPeer.EXPECT().ReadMsg().Return(msg, nil)

	nA := &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 1}
	mPeer.EXPECT().RemoteAddress().Return(nA)

	assert.EqualError(t, pm.handle(mPeer), "cs protocol version not match")

	time.Sleep(500 * time.Millisecond)
}

func TestCsProtocolManager_handle2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mPeer := NewMockPmAbstractPeer(ctrl)
	mps1 := NewMockAbstractPeerSet(ctrl)
	mps2 := NewMockAbstractPeerSet(ctrl)
	mps3 := NewMockAbstractPeerSet(ctrl)

	psm := &CsPmPeerSetManager{
		basePeers:            mps1,
		currentVerifierPeers: mps2,
		nextVerifierPeers:    mps3,
	}

	mChain := NewMockChain(ctrl)
	mNodeConf := NewMockNodeConf(ctrl)
	mSigner := NewMockPbftSigner(ctrl)
	mP2PServer := NewMockP2PServer(ctrl)

	mVReader := NewMockVerifiersReader(ctrl)

	pm := &CsProtocolManager{
		CsProtocolManagerConfig: &CsProtocolManagerConfig{
			ChainConfig:     *chain_config.GetChainConfig(),
			Chain:           mChain,
			NodeConf:        mNodeConf,
			MsgSigner:       mSigner,
			P2PServer:       mP2PServer,
			VerifiersReader: mVReader,
		},
		peerSetManager: psm,
	}

	pm.pmType.Store(verifier)
	pm.maxPeers = 7
	PbftMaxPeerCount = 3
	mps2.EXPECT().Len().Return(2)
	//mps1.EXPECT().Len().Return(2)

	var pyRecord, _ = hex.DecodeString("f884b8407098ad865b00a582051940cb9cf36836572411a47278783077011599ed5cd16b76f2635f4e234738f30813a89eb9137e3e3df5266e3a1f11df72ecf1145ccb9c01826964827634826970847f00000189736563703235366b31a103ca634cae0d49acb401d8a4c6b6fe8c55b70d115bf400769cc1400f3258cd31388375647082765f")

	var r enr.Record
	if err := rlp.DecodeBytes(pyRecord, &r); err != nil {
		t.Fatalf("can't decode: %v", err)
	}
	n, err := enode.New(enode.ValidSchemes, &r)
	if err != nil {
		t.Fatalf("can't verify record: %v", err)
	}

	block := model.NewBlock(model.NewHeader(11, 101, common.HexToHash("ss"), common.HexToHash("fdfs"), common.StringToDiff("0x22"), big.NewInt(111), common.StringToAddress("fdsfds"), common.EncodeNonce(33)), nil, nil)

	// case 2
	mChain.EXPECT().GetBlockByNumber(gomock.Eq(uint64(0))).Return(block)

	//  send
	mChain.EXPECT().CurrentBlock().Return(block)
	mNodeConf.EXPECT().GetNodeType().Return(chain_config.NodeTypeOfNormal).Times(2)
	mNodeConf.EXPECT().GetNodeName().Return("dsadsad")
	mP2PServer.EXPECT().Self().Return(n)
	mPeer.EXPECT().SendMsg(gomock.Any(), gomock.Any()).Return(nil)

	hsData := HandShakeData{
		ProtocolVersion:    chain_config.CsProtocolVersion,
		ChainID:            big.NewInt(2),
		NetworkId:          pm.ChainConfig.NetworkID,
		CurrentBlock:       common.HexToHash("aaa"),
		CurrentBlockHeight: 64,
		GenesisBlock:       block.Hash(),
		NodeType:           chain_config.NodeTypeOfVerifier,
		NodeName:           "test",
		RawUrl:             n.String(),
	}

	statusData := &StatusData{HandShakeData: hsData}

	hash := statusData.DataHash()

	assert.Equal(t, true, !hash.IsEmpty())

	account := tests.AccFactory.GenAccount()

	sign, err := account.SignHash(statusData.DataHash().Bytes())

	assert.NoError(t, err)

	assert.Equal(t, true, len(sign) > 0)

	statusData.Sign = sign

	statusData.PubKey = crypto.CompressPubkey(&account.Pk.PublicKey)

	size, r1, err := rlp.EncodeToReader(statusData)
	assert.NoError(t, err)

	msg := p2p.Msg{Code: StatusMsg, Size: uint32(size), Payload: r1}

	// read
	mPeer.EXPECT().ReadMsg().Return(msg, nil)

	mPeer.EXPECT().SetRemoteVerifierAddress(gomock.Any())
	mPeer.EXPECT().SetNodeType(gomock.Any())
	mPeer.EXPECT().SetNodeName(gomock.Any())
	mPeer.EXPECT().SetHead(gomock.Any(), gomock.Any())
	nA := &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 1}

	mPeer.EXPECT().RemoteAddress().Return(nA)
	mPeer.EXPECT().SetPeerRawUrl(gomock.Any())

	mVReader.EXPECT().ShouldChangeVerifier().Return(true)
	addr := common.StringToAddress("aaa")
	mVReader.EXPECT().NextVerifiers().Return([]common.Address{addr})
	mPeer.EXPECT().RemoteVerifierAddress().Return(addr)

	mps2.EXPECT().GetPeers().Return(map[string]PmAbstractPeer{"a": mPeer})
	mPeer.EXPECT().RemoteVerifierAddress().Return(addr)
	mPeer.EXPECT().RemoteVerifierAddress().Return(addr)

	assert.EqualError(t, pm.handle(mPeer), "current verifier address already in peer set")

	time.Sleep(500 * time.Millisecond)
}

func TestCsProtocolManager_handle3(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mPeer := NewMockPmAbstractPeer(ctrl)
	mps1 := NewMockAbstractPeerSet(ctrl)
	mps2 := NewMockAbstractPeerSet(ctrl)
	mps3 := NewMockAbstractPeerSet(ctrl)

	psm := &CsPmPeerSetManager{
		basePeers:            mps1,
		currentVerifierPeers: mps2,
		nextVerifierPeers:    mps3,
	}

	mChain := NewMockChain(ctrl)
	mNodeConf := NewMockNodeConf(ctrl)
	mSigner := NewMockPbftSigner(ctrl)
	mP2PServer := NewMockP2PServer(ctrl)

	mVReader := NewMockVerifiersReader(ctrl)

	pm := &CsProtocolManager{
		CsProtocolManagerConfig: &CsProtocolManagerConfig{
			ChainConfig:     *chain_config.GetChainConfig(),
			Chain:           mChain,
			NodeConf:        mNodeConf,
			MsgSigner:       mSigner,
			P2PServer:       mP2PServer,
			VerifiersReader: mVReader,
		},
		peerSetManager: psm,
	}

	pm.pmType.Store(verifier)
	pm.maxPeers = 7
	PbftMaxPeerCount = 3
	mps2.EXPECT().Len().Return(2)
	//mps1.EXPECT().Len().Return(2)

	var pyRecord, _ = hex.DecodeString("f884b8407098ad865b00a582051940cb9cf36836572411a47278783077011599ed5cd16b76f2635f4e234738f30813a89eb9137e3e3df5266e3a1f11df72ecf1145ccb9c01826964827634826970847f00000189736563703235366b31a103ca634cae0d49acb401d8a4c6b6fe8c55b70d115bf400769cc1400f3258cd31388375647082765f")

	var r enr.Record
	if err := rlp.DecodeBytes(pyRecord, &r); err != nil {
		t.Fatalf("can't decode: %v", err)
	}
	n, err := enode.New(enode.ValidSchemes, &r)
	if err != nil {
		t.Fatalf("can't verify record: %v", err)
	}

	block := model.NewBlock(model.NewHeader(11, 101, common.HexToHash("ss"), common.HexToHash("fdfs"), common.StringToDiff("0x22"), big.NewInt(111), common.StringToAddress("fdsfds"), common.EncodeNonce(33)), nil, nil)

	// case 2
	mChain.EXPECT().GetBlockByNumber(gomock.Eq(uint64(0))).Return(block)

	//  send
	mChain.EXPECT().CurrentBlock().Return(block)
	mNodeConf.EXPECT().GetNodeType().Return(chain_config.NodeTypeOfNormal).Times(2)
	mNodeConf.EXPECT().GetNodeName().Return("dsadsad")
	mP2PServer.EXPECT().Self().Return(n)
	mPeer.EXPECT().SendMsg(gomock.Any(), gomock.Any()).Return(nil)

	hsData := HandShakeData{
		ProtocolVersion:    chain_config.CsProtocolVersion,
		ChainID:            big.NewInt(2),
		NetworkId:          pm.ChainConfig.NetworkID,
		CurrentBlock:       common.HexToHash("aaa"),
		CurrentBlockHeight: 64,
		GenesisBlock:       block.Hash(),
		NodeType:           chain_config.NodeTypeOfVerifier,
		NodeName:           "test",
		RawUrl:             n.String(),
	}

	statusData := &StatusData{HandShakeData: hsData}

	hash := statusData.DataHash()

	assert.Equal(t, true, !hash.IsEmpty())

	account := tests.AccFactory.GenAccount()

	sign, err := account.SignHash(statusData.DataHash().Bytes())

	assert.NoError(t, err)

	assert.Equal(t, true, len(sign) > 0)

	statusData.Sign = sign

	statusData.PubKey = crypto.CompressPubkey(&account.Pk.PublicKey)

	size, r1, err := rlp.EncodeToReader(statusData)
	assert.NoError(t, err)

	msg := p2p.Msg{Code: StatusMsg, Size: uint32(size), Payload: r1}

	// read
	mPeer.EXPECT().ReadMsg().Return(msg, nil)

	mPeer.EXPECT().SetRemoteVerifierAddress(gomock.Any())
	mPeer.EXPECT().SetNodeType(gomock.Any())
	mPeer.EXPECT().SetNodeName(gomock.Any())
	mPeer.EXPECT().SetHead(gomock.Any(), gomock.Any())
	nA := &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 1}

	mPeer.EXPECT().RemoteAddress().Return(nA)
	mPeer.EXPECT().SetPeerRawUrl(gomock.Any())

	mVReader.EXPECT().ShouldChangeVerifier().Return(true)
	addr := common.StringToAddress("aaa")
	mVReader.EXPECT().NextVerifiers().Return([]common.Address{})

	mVReader.EXPECT().ShouldChangeVerifier().Return(false)
	mVReader.EXPECT().NextVerifiers().Return([]common.Address{addr})
	mPeer.EXPECT().RemoteVerifierAddress().Return(addr)

	mps3.EXPECT().GetPeers().Return(map[string]PmAbstractPeer{"a": mPeer})
	mPeer.EXPECT().RemoteVerifierAddress().Return(addr)
	mPeer.EXPECT().RemoteVerifierAddress().Return(addr)

	assert.EqualError(t, pm.handle(mPeer), "next verifier address already in peer set")

	time.Sleep(500 * time.Millisecond)
}

func TestCsProtocolManager_handleMsg(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mP := NewMockPmAbstractPeer(ctrl)

	pm := &CsProtocolManager{}

	mP.EXPECT().ReadMsg().Return(p2p.Msg{}, errors.New("fdsfds"))
	mP.EXPECT().NodeName().Return("ddd").AnyTimes()
	assert.EqualError(t, pm.handleMsg(mP), "fdsfds")
}

func TestCsProtocolManager_handleMsg1(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mP := NewMockPmAbstractPeer(ctrl)

	pm := &CsProtocolManager{}

	// blockHashMsg
	data := &blockHashMsg{BlockHash: common.HexToHash("vfd"), BlockNumber: 11}
	_, r, err := rlp.EncodeToReader(data)
	assert.NoError(t, err)
	msg := p2p.Msg{Code: 0x11, Size: uint32(ProtocolMaxMsgSize * 10), Payload: r}
	mP.EXPECT().ReadMsg().Return(msg, nil)

	assert.EqualError(t, pm.handleMsg(mP), msgTooLargeErr.Error())
}

func TestCsProtocolManager_handleMsg2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mP := NewMockPmAbstractPeer(ctrl)
	mPbftNode := NewMockPbftNode(ctrl)

	pm := &CsProtocolManager{
		CsProtocolManagerConfig: &CsProtocolManagerConfig{
			PbftNode: mPbftNode,
		},
	}

	// blockHashMsg
	data := struct {
		a string
	}{
		a: "aa",
	}

	pm.pmType.Store(verifier)
	size, r, err := rlp.EncodeToReader(data)
	assert.NoError(t, err)
	msg := p2p.Msg{Code: 0x101, Size: uint32(size), Payload: r}
	mP.EXPECT().ReadMsg().Return(msg, nil)
	mPbftNode.EXPECT().OnNewP2PMsg(gomock.Eq(msg), gomock.Any()).Return(nil)
	assert.Nil(t, pm.handleMsg(mP))
}

func TestCsProtocolManager_handleMsg3(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mP := NewMockPmAbstractPeer(ctrl)
	mPbftNode := NewMockPbftNode(ctrl)

	pm := &CsProtocolManager{
		CsProtocolManagerConfig: &CsProtocolManagerConfig{
			PbftNode: mPbftNode,
		},
	}

	// blockHashMsg
	data := struct {
		a string
	}{
		a: "aa",
	}

	pm.pmType.Store(verifier)
	size, r, err := rlp.EncodeToReader(data)
	assert.NoError(t, err)
	msg := p2p.Msg{Code: 0x101, Size: uint32(size), Payload: r}
	mP.EXPECT().ReadMsg().Return(msg, nil)
	mPbftNode.EXPECT().OnNewP2PMsg(gomock.Eq(msg), gomock.Any()).Return(errors.New("fdf"))
	assert.EqualError(t, pm.handleMsg(mP), "fdf")
}

func TestCsProtocolManager_handleMsg4(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mP := NewMockPmAbstractPeer(ctrl)
	mPbftNode := NewMockPbftNode(ctrl)
	mNodeConf := NewMockNodeConf(ctrl)

	pm := &CsProtocolManager{
		CsProtocolManagerConfig: &CsProtocolManagerConfig{
			PbftNode: mPbftNode,
			NodeConf: mNodeConf,
		},
	}

	// blockHashMsg
	data := struct {
		a string
	}{
		a: "aa",
	}

	pm.pmType.Store(base)
	size, r, err := rlp.EncodeToReader(data)
	assert.NoError(t, err)
	msg := p2p.Msg{Code: 0x101, Size: uint32(size), Payload: r}
	mP.EXPECT().ReadMsg().Return(msg, nil)
	mNodeConf.EXPECT().GetNodeType().Return(chain_config.NodeTypeOfNormal)
	assert.Nil(t, pm.handleMsg(mP))
}

func TestCsProtocolManager_handleMsg5(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mP := NewMockPmAbstractPeer(ctrl)
	mPbftNode := NewMockPbftNode(ctrl)
	mNodeConf := NewMockNodeConf(ctrl)

	pm := &CsProtocolManager{
		CsProtocolManagerConfig: &CsProtocolManagerConfig{
			PbftNode: mPbftNode,
			NodeConf: mNodeConf,
		},
	}

	// blockHashMsg
	data := struct {
		a string
	}{
		a: "aa",
	}

	pm.pmType.Store(base)
	size, r, err := rlp.EncodeToReader(data)
	assert.NoError(t, err)
	msg := p2p.Msg{Code: 0x44, Size: uint32(size), Payload: r}
	mP.EXPECT().ReadMsg().Return(msg, nil)
	assert.EqualError(t, pm.handleMsg(mP), msgHandleFuncNotFoundErr.Error())
}

func TestCsProtocolManager_handShake(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mChain := NewMockChain(ctrl)
	mNodeConf := NewMockNodeConf(ctrl)
	mSigner := NewMockPbftSigner(ctrl)
	mP2PServer := NewMockP2PServer(ctrl)

	mPeer := NewMockPmAbstractPeer(ctrl)

	pm := &CsProtocolManager{
		CsProtocolManagerConfig: &CsProtocolManagerConfig{
			ChainConfig: *chain_config.GetChainConfig(),
			Chain:       mChain,
			NodeConf:    mNodeConf,
			MsgSigner:   mSigner,
			P2PServer:   mP2PServer,
		},
	}

	var pyRecord, _ = hex.DecodeString("f884b8407098ad865b00a582051940cb9cf36836572411a47278783077011599ed5cd16b76f2635f4e234738f30813a89eb9137e3e3df5266e3a1f11df72ecf1145ccb9c01826964827634826970847f00000189736563703235366b31a103ca634cae0d49acb401d8a4c6b6fe8c55b70d115bf400769cc1400f3258cd31388375647082765f")

	var r enr.Record
	if err := rlp.DecodeBytes(pyRecord, &r); err != nil {
		t.Fatalf("can't decode: %v", err)
	}
	n, err := enode.New(enode.ValidSchemes, &r)
	if err != nil {
		t.Fatalf("can't verify record: %v", err)
	}

	block := model.NewBlock(model.NewHeader(11, 101, common.HexToHash("ss"), common.HexToHash("fdfs"), common.StringToDiff("0x22"), big.NewInt(111), common.StringToAddress("fdsfds"), common.EncodeNonce(33)), nil, nil)

	// case 1
	mChain.EXPECT().GetBlockByNumber(gomock.Eq(uint64(0))).Return(block)

	//  send
	mChain.EXPECT().CurrentBlock().Return(block)
	mNodeConf.EXPECT().GetNodeType().Return(chain_config.NodeTypeOfNormal).Times(2)
	mNodeConf.EXPECT().GetNodeName().Return("dsadsad")
	mP2PServer.EXPECT().Self().Return(n)
	mPeer.EXPECT().SendMsg(gomock.Any(), gomock.Any()).Return(nil)

	// read
	mPeer.EXPECT().ReadMsg().Return(p2p.Msg{}, errors.New("dfds"))
	assert.EqualError(t, pm.HandShake(mPeer), "can't read hand shake msg from remote")

	time.Sleep(600 * time.Millisecond)

	// case 2
	mChain.EXPECT().GetBlockByNumber(gomock.Eq(uint64(0))).Return(block)

	//  send
	mChain.EXPECT().CurrentBlock().Return(block)
	mNodeConf.EXPECT().GetNodeType().Return(chain_config.NodeTypeOfNormal).Times(2)
	mNodeConf.EXPECT().GetNodeName().Return("dsadsad")
	mP2PServer.EXPECT().Self().Return(n)
	mPeer.EXPECT().SendMsg(gomock.Any(), gomock.Any()).Return(nil)

	hsData := HandShakeData{
		ProtocolVersion:    1,
		ChainID:            big.NewInt(2),
		NetworkId:          1321,
		CurrentBlock:       common.HexToHash("aaa"),
		CurrentBlockHeight: 64,
		GenesisBlock:       common.HexToHash("sss"),
		NodeType:           2,
		NodeName:           "test",
		RawUrl:             n.String(),
	}

	statusData := &StatusData{HandShakeData: hsData}

	hash := statusData.DataHash()

	assert.Equal(t, true, !hash.IsEmpty())

	account := tests.AccFactory.GenAccount()

	sign, err := account.SignHash(statusData.DataHash().Bytes())

	assert.NoError(t, err)

	assert.Equal(t, true, len(sign) > 0)

	statusData.Sign = sign

	statusData.PubKey = crypto.CompressPubkey(&account.Pk.PublicKey)

	size, r1, err := rlp.EncodeToReader(statusData)
	assert.NoError(t, err)

	msg := p2p.Msg{Code: StatusMsg, Size: uint32(size), Payload: r1}

	// read
	mPeer.EXPECT().ReadMsg().Return(msg, nil)
	assert.EqualError(t, pm.HandShake(mPeer), "network id not match")

	time.Sleep(600 * time.Millisecond)
}

func TestCsProtocolManager_handShake1(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mChain := NewMockChain(ctrl)
	mNodeConf := NewMockNodeConf(ctrl)
	mSigner := NewMockPbftSigner(ctrl)
	mP2PServer := NewMockP2PServer(ctrl)

	mPeer := NewMockPmAbstractPeer(ctrl)

	pm := &CsProtocolManager{
		CsProtocolManagerConfig: &CsProtocolManagerConfig{
			ChainConfig: *chain_config.GetChainConfig(),
			Chain:       mChain,
			NodeConf:    mNodeConf,
			MsgSigner:   mSigner,
			P2PServer:   mP2PServer,
		},
	}

	var pyRecord, _ = hex.DecodeString("f884b8407098ad865b00a582051940cb9cf36836572411a47278783077011599ed5cd16b76f2635f4e234738f30813a89eb9137e3e3df5266e3a1f11df72ecf1145ccb9c01826964827634826970847f00000189736563703235366b31a103ca634cae0d49acb401d8a4c6b6fe8c55b70d115bf400769cc1400f3258cd31388375647082765f")

	var r enr.Record
	if err := rlp.DecodeBytes(pyRecord, &r); err != nil {
		t.Fatalf("can't decode: %v", err)
	}
	n, err := enode.New(enode.ValidSchemes, &r)
	if err != nil {
		t.Fatalf("can't verify record: %v", err)
	}

	block := model.NewBlock(model.NewHeader(11, 101, common.HexToHash("ss"), common.HexToHash("fdfs"), common.StringToDiff("0x22"), big.NewInt(111), common.StringToAddress("fdsfds"), common.EncodeNonce(33)), nil, nil)

	// case 2
	mChain.EXPECT().GetBlockByNumber(gomock.Eq(uint64(0))).Return(block)

	//  send
	mChain.EXPECT().CurrentBlock().Return(block)
	mNodeConf.EXPECT().GetNodeType().Return(chain_config.NodeTypeOfNormal).Times(2)
	mNodeConf.EXPECT().GetNodeName().Return("dsadsad")
	mP2PServer.EXPECT().Self().Return(n)
	mPeer.EXPECT().SendMsg(gomock.Any(), gomock.Any()).Return(nil)

	hsData := HandShakeData{
		ProtocolVersion:    1,
		ChainID:            big.NewInt(2),
		NetworkId:          pm.ChainConfig.NetworkID,
		CurrentBlock:       common.HexToHash("aaa"),
		CurrentBlockHeight: 64,
		GenesisBlock:       common.HexToHash("sss"),
		NodeType:           2,
		NodeName:           "test",
		RawUrl:             n.String(),
	}

	statusData := &StatusData{HandShakeData: hsData}

	hash := statusData.DataHash()

	assert.Equal(t, true, !hash.IsEmpty())

	account := tests.AccFactory.GenAccount()

	sign, err := account.SignHash(statusData.DataHash().Bytes())

	assert.NoError(t, err)

	assert.Equal(t, true, len(sign) > 0)

	statusData.Sign = sign

	statusData.PubKey = crypto.CompressPubkey(&account.Pk.PublicKey)

	size, r1, err := rlp.EncodeToReader(statusData)
	assert.NoError(t, err)

	msg := p2p.Msg{Code: StatusMsg, Size: uint32(size), Payload: r1}

	// read
	mPeer.EXPECT().ReadMsg().Return(msg, nil)
	assert.EqualError(t, pm.HandShake(mPeer), "genesis block not match")

	time.Sleep(600 * time.Millisecond)
}

func TestCsProtocolManager_handShake2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mChain := NewMockChain(ctrl)
	mNodeConf := NewMockNodeConf(ctrl)
	mSigner := NewMockPbftSigner(ctrl)
	mP2PServer := NewMockP2PServer(ctrl)

	mPeer := NewMockPmAbstractPeer(ctrl)

	pm := &CsProtocolManager{
		CsProtocolManagerConfig: &CsProtocolManagerConfig{
			ChainConfig: *chain_config.GetChainConfig(),
			Chain:       mChain,
			NodeConf:    mNodeConf,
			MsgSigner:   mSigner,
			P2PServer:   mP2PServer,
		},
	}

	var pyRecord, _ = hex.DecodeString("f884b8407098ad865b00a582051940cb9cf36836572411a47278783077011599ed5cd16b76f2635f4e234738f30813a89eb9137e3e3df5266e3a1f11df72ecf1145ccb9c01826964827634826970847f00000189736563703235366b31a103ca634cae0d49acb401d8a4c6b6fe8c55b70d115bf400769cc1400f3258cd31388375647082765f")

	var r enr.Record
	if err := rlp.DecodeBytes(pyRecord, &r); err != nil {
		t.Fatalf("can't decode: %v", err)
	}
	n, err := enode.New(enode.ValidSchemes, &r)
	if err != nil {
		t.Fatalf("can't verify record: %v", err)
	}

	block := model.NewBlock(model.NewHeader(11, 101, common.HexToHash("ss"), common.HexToHash("fdfs"), common.StringToDiff("0x22"), big.NewInt(111), common.StringToAddress("fdsfds"), common.EncodeNonce(33)), nil, nil)

	// case 2
	mChain.EXPECT().GetBlockByNumber(gomock.Eq(uint64(0))).Return(block)

	//  send
	mChain.EXPECT().CurrentBlock().Return(block)
	mNodeConf.EXPECT().GetNodeType().Return(chain_config.NodeTypeOfNormal).Times(2)
	mNodeConf.EXPECT().GetNodeName().Return("dsadsad")
	mP2PServer.EXPECT().Self().Return(n)
	mPeer.EXPECT().SendMsg(gomock.Any(), gomock.Any()).Return(nil)

	hsData := HandShakeData{
		ProtocolVersion:    0,
		ChainID:            big.NewInt(2),
		NetworkId:          pm.ChainConfig.NetworkID,
		CurrentBlock:       common.HexToHash("aaa"),
		CurrentBlockHeight: 64,
		GenesisBlock:       block.Hash(),
		NodeType:           chain_config.NodeTypeOfVerifier,
		NodeName:           "test",
		RawUrl:             n.String(),
	}

	statusData := &StatusData{HandShakeData: hsData}

	hash := statusData.DataHash()

	assert.Equal(t, true, !hash.IsEmpty())

	account := tests.AccFactory.GenAccount()

	sign, err := account.SignHash(statusData.DataHash().Bytes())

	assert.NoError(t, err)

	assert.Equal(t, true, len(sign) > 0)

	statusData.Sign = sign

	statusData.PubKey = crypto.CompressPubkey(&account.Pk.PublicKey)

	size, r1, err := rlp.EncodeToReader(statusData)
	assert.NoError(t, err)

	msg := p2p.Msg{Code: StatusMsg, Size: uint32(size), Payload: r1}

	// read
	mPeer.EXPECT().ReadMsg().Return(msg, nil)
	assert.EqualError(t, pm.HandShake(mPeer), "can't read hand shake msg")

	time.Sleep(600 * time.Millisecond)
}

func TestCsProtocolManager_handShake3(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mChain := NewMockChain(ctrl)
	mNodeConf := NewMockNodeConf(ctrl)
	mSigner := NewMockPbftSigner(ctrl)
	mP2PServer := NewMockP2PServer(ctrl)

	mPeer := NewMockPmAbstractPeer(ctrl)

	pm := &CsProtocolManager{
		CsProtocolManagerConfig: &CsProtocolManagerConfig{
			ChainConfig: *chain_config.GetChainConfig(),
			Chain:       mChain,
			NodeConf:    mNodeConf,
			MsgSigner:   mSigner,
			P2PServer:   mP2PServer,
		},
	}

	var pyRecord, _ = hex.DecodeString("f884b8407098ad865b00a582051940cb9cf36836572411a47278783077011599ed5cd16b76f2635f4e234738f30813a89eb9137e3e3df5266e3a1f11df72ecf1145ccb9c01826964827634826970847f00000189736563703235366b31a103ca634cae0d49acb401d8a4c6b6fe8c55b70d115bf400769cc1400f3258cd31388375647082765f")

	var r enr.Record
	if err := rlp.DecodeBytes(pyRecord, &r); err != nil {
		t.Fatalf("can't decode: %v", err)
	}
	n, err := enode.New(enode.ValidSchemes, &r)
	if err != nil {
		t.Fatalf("can't verify record: %v", err)
	}

	block := model.NewBlock(model.NewHeader(11, 101, common.HexToHash("ss"), common.HexToHash("fdfs"), common.StringToDiff("0x22"), big.NewInt(111), common.StringToAddress("fdsfds"), common.EncodeNonce(33)), nil, nil)

	// case 2
	mChain.EXPECT().GetBlockByNumber(gomock.Eq(uint64(0))).Return(block)

	//  send
	mChain.EXPECT().CurrentBlock().Return(block)
	mNodeConf.EXPECT().GetNodeType().Return(chain_config.NodeTypeOfNormal).Times(2)
	mNodeConf.EXPECT().GetNodeName().Return("dsadsad")
	mP2PServer.EXPECT().Self().Return(n)
	mPeer.EXPECT().SendMsg(gomock.Any(), gomock.Any()).Return(nil)

	hsData := HandShakeData{
		ProtocolVersion:    11,
		ChainID:            big.NewInt(2),
		NetworkId:          pm.ChainConfig.NetworkID,
		CurrentBlock:       common.HexToHash("aaa"),
		CurrentBlockHeight: 64,
		GenesisBlock:       block.Hash(),
		NodeType:           chain_config.NodeTypeOfVerifier,
		NodeName:           "test",
		RawUrl:             n.String(),
	}

	statusData := &StatusData{HandShakeData: hsData}

	hash := statusData.DataHash()

	assert.Equal(t, true, !hash.IsEmpty())

	account := tests.AccFactory.GenAccount()

	sign, err := account.SignHash(statusData.DataHash().Bytes())

	assert.NoError(t, err)

	assert.Equal(t, true, len(sign) > 0)

	statusData.Sign = sign

	statusData.PubKey = crypto.CompressPubkey(&account.Pk.PublicKey)

	size, r1, err := rlp.EncodeToReader(statusData)
	assert.NoError(t, err)

	msg := p2p.Msg{Code: StatusMsg, Size: uint32(size), Payload: r1}

	// read
	mPeer.EXPECT().ReadMsg().Return(msg, nil)
	assert.EqualError(t, pm.HandShake(mPeer), "cs protocol version not match")

	time.Sleep(600 * time.Millisecond)
}

func TestCsProtocolManager_handShake4(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mChain := NewMockChain(ctrl)
	mNodeConf := NewMockNodeConf(ctrl)
	mSigner := NewMockPbftSigner(ctrl)
	mP2PServer := NewMockP2PServer(ctrl)

	mPeer := NewMockPmAbstractPeer(ctrl)

	pm := &CsProtocolManager{
		CsProtocolManagerConfig: &CsProtocolManagerConfig{
			ChainConfig: *chain_config.GetChainConfig(),
			Chain:       mChain,
			NodeConf:    mNodeConf,
			MsgSigner:   mSigner,
			P2PServer:   mP2PServer,
		},
	}

	var pyRecord, _ = hex.DecodeString("f884b8407098ad865b00a582051940cb9cf36836572411a47278783077011599ed5cd16b76f2635f4e234738f30813a89eb9137e3e3df5266e3a1f11df72ecf1145ccb9c01826964827634826970847f00000189736563703235366b31a103ca634cae0d49acb401d8a4c6b6fe8c55b70d115bf400769cc1400f3258cd31388375647082765f")

	var r enr.Record
	if err := rlp.DecodeBytes(pyRecord, &r); err != nil {
		t.Fatalf("can't decode: %v", err)
	}
	n, err := enode.New(enode.ValidSchemes, &r)
	if err != nil {
		t.Fatalf("can't verify record: %v", err)
	}

	block := model.NewBlock(model.NewHeader(11, 101, common.HexToHash("ss"), common.HexToHash("fdfs"), common.StringToDiff("0x22"), big.NewInt(111), common.StringToAddress("fdsfds"), common.EncodeNonce(33)), nil, nil)

	// case 2
	mChain.EXPECT().GetBlockByNumber(gomock.Eq(uint64(0))).Return(block)

	//  send
	mChain.EXPECT().CurrentBlock().Return(block)
	mNodeConf.EXPECT().GetNodeType().Return(chain_config.NodeTypeOfNormal).Times(2)
	mNodeConf.EXPECT().GetNodeName().Return("dsadsad")
	mP2PServer.EXPECT().Self().Return(n)
	mPeer.EXPECT().SendMsg(gomock.Any(), gomock.Any()).Return(nil)

	hsData := HandShakeData{
		ProtocolVersion:    chain_config.CsProtocolVersion,
		ChainID:            big.NewInt(2),
		NetworkId:          pm.ChainConfig.NetworkID,
		CurrentBlock:       common.HexToHash("aaa"),
		CurrentBlockHeight: 64,
		GenesisBlock:       block.Hash(),
		NodeType:           chain_config.NodeTypeOfVerifier,
		NodeName:           "test",
		RawUrl:             n.String(),
	}

	statusData := &StatusData{HandShakeData: hsData}

	hash := statusData.DataHash()

	assert.Equal(t, true, !hash.IsEmpty())

	account := tests.AccFactory.GenAccount()

	sign, err := account.SignHash(statusData.DataHash().Bytes())

	assert.NoError(t, err)

	assert.Equal(t, true, len(sign) > 0)

	statusData.Sign = sign

	statusData.PubKey = crypto.CompressPubkey(&account.Pk.PublicKey)

	size, r1, err := rlp.EncodeToReader(statusData)
	assert.NoError(t, err)

	msg := p2p.Msg{Code: StatusMsg, Size: uint32(size), Payload: r1}

	// read
	mPeer.EXPECT().ReadMsg().Return(msg, nil)

	mPeer.EXPECT().SetRemoteVerifierAddress(gomock.Any())
	mPeer.EXPECT().SetNodeType(gomock.Any())
	mPeer.EXPECT().SetNodeName(gomock.Any())
	mPeer.EXPECT().SetHead(gomock.Any(), gomock.Any())
	nA := &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 1}

	mPeer.EXPECT().RemoteAddress().Return(nA)
	mPeer.EXPECT().SetPeerRawUrl(gomock.Any())

	assert.Nil(t, pm.HandShake(mPeer))

	time.Sleep(600 * time.Millisecond)
}

func TestCsProtocolManager_checkAndHandleVerBootNodes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pm := &CsProtocolManager{}
	pm.pmType.Store(base)
	pm.checkAndHandleVerBootNodes()
}

func TestCsProtocolManager_checkAndHandleVerBootNodes1(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mVReader := NewMockVerifiersReader(ctrl)
	mSigner := NewMockPbftSigner(ctrl)
	mPs := NewMockAbstractPeerSet(ctrl)

	psm := &CsPmPeerSetManager{
		verifierBootNode:     mPs,
		basePeers:            newPeerSet(),
		currentVerifierPeers: newPeerSet(),
		nextVerifierPeers:    newPeerSet(),
	}

	pm := &CsProtocolManager{
		CsProtocolManagerConfig: &CsProtocolManagerConfig{
			MsgSigner:       mSigner,
			VerifiersReader: mVReader,
		},
		peerSetManager: psm,
	}

	addr := common.StringToAddress("aaa")
	mVReader.EXPECT().ShouldChangeVerifier().Return(false).AnyTimes()
	mVReader.EXPECT().CurrentVerifiers().Return([]common.Address{addr}).AnyTimes()
	mVReader.EXPECT().NextVerifiers().Return([]common.Address{}).AnyTimes()
	mSigner.EXPECT().GetAddress().Return(addr).AnyTimes()

	mPs.EXPECT().BestPeer().Return(nil)
	mPs.EXPECT().Len().Return(3).AnyTimes()

	pm.pmType.Store(verifier)
	pm.checkAndHandleVerBootNodes()
}

func TestCsProtocolManager_checkAndHandleVerBootNodes2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mVReader := NewMockVerifiersReader(ctrl)
	mSigner := NewMockPbftSigner(ctrl)

	mP2pServer := NewMockP2PServer(ctrl)

	mPs := NewMockAbstractPeerSet(ctrl)

	psm := &CsPmPeerSetManager{
		verifierBootNode:     mPs,
		basePeers:            newPeerSet(),
		currentVerifierPeers: newPeerSet(),
		nextVerifierPeers:    newPeerSet(),
	}

	pm := &CsProtocolManager{
		CsProtocolManagerConfig: &CsProtocolManagerConfig{
			MsgSigner:       mSigner,
			VerifiersReader: mVReader,
			P2PServer:       mP2pServer,
		},
		peerSetManager: psm,
	}

	addr := common.StringToAddress("aaa")
	mVReader.EXPECT().ShouldChangeVerifier().Return(false).AnyTimes()
	mVReader.EXPECT().CurrentVerifiers().Return([]common.Address{addr}).AnyTimes()
	mVReader.EXPECT().NextVerifiers().Return([]common.Address{}).AnyTimes()
	mSigner.EXPECT().GetAddress().Return(addr).AnyTimes()
	mPs.EXPECT().BestPeer().Return(nil)

	mPs.EXPECT().Len().Return(2).AnyTimes()

	var pyRecord, _ = hex.DecodeString("f884b8407098ad865b00a582051940cb9cf36836572411a47278783077011599ed5cd16b76f2635f4e234738f30813a89eb9137e3e3df5266e3a1f11df72ecf1145ccb9c01826964827634826970847f00000189736563703235366b31a103ca634cae0d49acb401d8a4c6b6fe8c55b70d115bf400769cc1400f3258cd31388375647082765f")

	var r enr.Record
	if err := rlp.DecodeBytes(pyRecord, &r); err != nil {
		t.Fatalf("can't decode: %v", err)
	}
	n, err := enode.New(enode.ValidSchemes, &r)
	if err != nil {
		t.Fatalf("can't verify record: %v", err)
	}

	mP2pServer.EXPECT().Self().Return(n).AnyTimes()
	mPs.EXPECT().Peer(gomock.Any()).Return(nil).AnyTimes()

	pm.pmType.Store(verifier)
	pm.checkAndHandleVerBootNodes()
}

func TestCsProtocolManager_bootVerifierConnCheck(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockNodeConf := NewMockNodeConf(ctrl)

	pm := &CsProtocolManager{CsProtocolManagerConfig: &CsProtocolManagerConfig{NodeConf: mockNodeConf}, stop: make(chan struct{})}

	mockNodeConf.EXPECT().GetNodeType().Return(chain_config.NodeTypeOfNormal).Times(2)
	pm.bootVerifierConnCheck()

	mockNodeConf.EXPECT().GetNodeType().Return(chain_config.NodeTypeOfVerifier)

	go pm.bootVerifierConnCheck()
	time.Sleep(time.Millisecond)
	close(pm.stop)
}

func TestCsProtocolManager_RegisterCommunicationService(t *testing.T) {
	pm := &CsProtocolManager{BaseProtocolManager: BaseProtocolManager{
		msgHandlers: make(map[uint64]func(msg p2p.Msg, p PmAbstractPeer) error),
	}}

	pm.RegisterCommunicationService(&testCm{}, &testCe{})
}

func TestCsProtocolManager_GetCurrentVerifierPeers(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// mock peer set
	mockPeerSetCurrentVerifierPeers := NewMockAbstractPeerSet(ctrl)

	mockPeerSetCurrentVerifierPeers.EXPECT().GetPeers().Return(nil)

	psm := &CsPmPeerSetManager{
		currentVerifierPeers: mockPeerSetCurrentVerifierPeers,
	}

	pm := &CsProtocolManager{peerSetManager: psm}

	pm.GetCurrentVerifierPeers()
}

func TestCsProtocolManager_IsSync(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChain := NewMockChain(ctrl)

	pm := &CsProtocolManager{CsProtocolManagerConfig: &CsProtocolManagerConfig{Chain: mockChain}}

	mockChain.EXPECT().CurrentBlock().Return(nil)
	assert.Equal(t, true, pm.IsSync())

	// mock peer set
	mockPeerSetBasePeers := NewMockAbstractPeerSet(ctrl)
	mockPeerSetCurrentVerifierPeers := NewMockAbstractPeerSet(ctrl)
	mockPeerSetnNextVerifierPeers := NewMockAbstractPeerSet(ctrl)
	mockPeerSetVerifierBootNode := NewMockAbstractPeerSet(ctrl)

	psm := &CsPmPeerSetManager{
		basePeers:            mockPeerSetBasePeers,
		currentVerifierPeers: mockPeerSetCurrentVerifierPeers,
		nextVerifierPeers:    mockPeerSetnNextVerifierPeers,
		verifierBootNode:     mockPeerSetVerifierBootNode,
	}

	pm.peerSetManager = psm

	// case 2
	block := model.NewBlock(model.NewHeader(11, 101, common.HexToHash("ss"), common.HexToHash("fdfs"), common.StringToDiff("0x22"), big.NewInt(111), common.StringToAddress("fdsfds"), common.EncodeNonce(33)), nil, nil)
	mockChain.EXPECT().CurrentBlock().Return(block)

	mockPeerSetBasePeers.EXPECT().BestPeer().Return(nil)
	mockPeerSetCurrentVerifierPeers.EXPECT().BestPeer().Return(nil)
	mockPeerSetnNextVerifierPeers.EXPECT().BestPeer().Return(nil)
	mockPeerSetVerifierBootNode.EXPECT().BestPeer().Return(nil)

	assert.Equal(t, false, pm.IsSync())

	// case 3
	mockPeer := NewMockPmAbstractPeer(ctrl)
	mockPeer.EXPECT().GetHead().Return(common.HexToHash("sad"), uint64(102))
	mockChain.EXPECT().CurrentBlock().Return(block)
	mockPeerSetBasePeers.EXPECT().BestPeer().Return(mockPeer)
	mockPeerSetCurrentVerifierPeers.EXPECT().BestPeer().Return(nil)
	mockPeerSetnNextVerifierPeers.EXPECT().BestPeer().Return(nil)
	mockPeerSetVerifierBootNode.EXPECT().BestPeer().Return(nil)

	assert.Equal(t, false, pm.IsSync())

	// case 4
	mockPeer.EXPECT().GetHead().Return(common.HexToHash("sad"), uint64(500))
	mockChain.EXPECT().CurrentBlock().Return(block)
	mockPeerSetBasePeers.EXPECT().BestPeer().Return(mockPeer)
	mockPeerSetCurrentVerifierPeers.EXPECT().BestPeer().Return(nil)
	mockPeerSetnNextVerifierPeers.EXPECT().BestPeer().Return(nil)
	mockPeerSetVerifierBootNode.EXPECT().BestPeer().Return(nil)
	assert.Equal(t, true, pm.IsSync())

}

func TestCsProtocolManager_PrintPeerHealthCheck(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// mock peer set
	mockPeerSetBasePeers := NewMockAbstractPeerSet(ctrl)
	mockPeerSetCurrentVerifierPeers := NewMockAbstractPeerSet(ctrl)
	mockPeerSetnNextVerifierPeers := NewMockAbstractPeerSet(ctrl)
	mockPeerSetVerifierBootNode := NewMockAbstractPeerSet(ctrl)

	// mock peer
	mockPeer := NewMockPmAbstractPeer(ctrl)
	mockPeer.EXPECT().NodeName().Return("aaa")
	mockPeer.EXPECT().RemoteVerifierAddress().Return(common.StringToAddress("dfs"))
	mockPeer.EXPECT().IsRunning().Return(true)
	mockPeerSetBasePeers.EXPECT().GetPeers().Return(map[string]PmAbstractPeer{"1": mockPeer}).AnyTimes()

	// mock peer
	mockPeer1 := NewMockPmAbstractPeer(ctrl)
	mockPeer1.EXPECT().NodeName().Return("aaa")
	mockPeer1.EXPECT().RemoteVerifierAddress().Return(common.StringToAddress("dfs"))
	mockPeer1.EXPECT().IsRunning().Return(true)
	mockPeerSetCurrentVerifierPeers.EXPECT().GetPeers().Return(map[string]PmAbstractPeer{"2": mockPeer1}).AnyTimes()

	// mock peer
	mockPeer2 := NewMockPmAbstractPeer(ctrl)
	mockPeer2.EXPECT().NodeName().Return("aaa")
	mockPeer2.EXPECT().RemoteVerifierAddress().Return(common.StringToAddress("dfs"))
	mockPeer2.EXPECT().IsRunning().Return(true)
	mockPeerSetnNextVerifierPeers.EXPECT().GetPeers().Return(map[string]PmAbstractPeer{"3": mockPeer2}).AnyTimes()

	// mock peer
	mockPeer3 := NewMockPmAbstractPeer(ctrl)
	mockPeer3.EXPECT().NodeName().Return("aaa")
	mockPeer3.EXPECT().RemoteVerifierAddress().Return(common.StringToAddress("dfs"))
	mockPeer3.EXPECT().IsRunning().Return(true)
	mockPeerSetVerifierBootNode.EXPECT().GetPeers().Return(map[string]PmAbstractPeer{"4": mockPeer3}).AnyTimes()

	psm := &CsPmPeerSetManager{
		basePeers:            mockPeerSetBasePeers,
		currentVerifierPeers: mockPeerSetCurrentVerifierPeers,
		nextVerifierPeers:    mockPeerSetnNextVerifierPeers,
		verifierBootNode:     mockPeerSetVerifierBootNode,
	}

	pm := &CsProtocolManager{
		peerSetManager: psm,
	}

	pm.PrintPeerHealthCheck()
}

func Test_printPeerInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPeer := NewMockPmAbstractPeer(ctrl)
	mockPeer.EXPECT().NodeName().Return("aaa")
	mockPeer.EXPECT().RemoteVerifierAddress().Return(common.StringToAddress("dfs"))
	mockPeer.EXPECT().IsRunning().Return(true)

	pMap := map[string]PmAbstractPeer{"1": mockPeer}
	printPeerInfo("ss", pMap)
}

func Test_newCsProtocolManager(t *testing.T) {
	assert.Panics(t, func() {
		PbftMaxPeerCount = 101
		newCsProtocolManager(nil)
	})
}

func TestCsProtocolManager_HaveEnoughVerifiers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mPs1 := NewMockAbstractPeerSet(ctrl)
	mPs2 := NewMockAbstractPeerSet(ctrl)

	psm := &CsPmPeerSetManager{
		currentVerifierPeers: mPs1,
		nextVerifierPeers:    mPs2,
	}

	pm := &CsProtocolManager{
		CsProtocolManagerConfig: &CsProtocolManagerConfig{
			ChainConfig: *chain_config.GetChainConfig(),
		},
		peerSetManager: psm,
	}

	mPs1.EXPECT().Len().Return(2)
	mPs2.EXPECT().Len().Return(2)

	mc, mn := pm.HaveEnoughVerifiers(false)
	assert.Equal(t, uint(pm.ChainConfig.VerifierNumber-1-2), mc)
	assert.Equal(t, uint(pm.ChainConfig.VerifierNumber-1-2), mn)

	mPs1.EXPECT().Len().Return(pm.ChainConfig.VerifierNumber)
	mPs2.EXPECT().Len().Return(pm.ChainConfig.VerifierNumber)

	mc1, mn1 := pm.HaveEnoughVerifiers(false)
	assert.Equal(t, uint(0), mc1)
	assert.Equal(t, uint(0), mn1)
}

func TestCsProtocolManager_CurrentVerifierPeersSet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mPs1 := NewMockAbstractPeerSet(ctrl)

	psm := &CsPmPeerSetManager{
		currentVerifierPeers: mPs1,
	}

	pm := &CsProtocolManager{
		peerSetManager: psm,
	}

	pm.CurrentVerifierPeersSet()
}

func TestCsProtocolManager_NextVerifierPeersSet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mPs1 := NewMockAbstractPeerSet(ctrl)

	psm := &CsPmPeerSetManager{
		nextVerifierPeers: mPs1,
	}

	pm := &CsProtocolManager{
		peerSetManager: psm,
	}

	pm.NextVerifierPeersSet()
}

func TestCsProtocolManager_ConnectPeer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mP2pServer := NewMockP2PServer(ctrl)

	pm := &CsProtocolManager{
		CsProtocolManagerConfig: &CsProtocolManagerConfig{
			P2PServer: mP2pServer,
		},
	}

	mP2pServer.EXPECT().AddPeer(gomock.Any())

	pm.ConnectPeer(nil)
}

func TestCsProtocolManager_connectVBoots(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mPs1 := NewMockAbstractPeerSet(ctrl)

	psm := &CsPmPeerSetManager{
		verifierBootNode: mPs1,
	}

	pm := &CsProtocolManager{
		peerSetManager: psm,
	}

	mPs1.EXPECT().Len().Return(chain_config.GetChainConfig().VerifierBootNodeNumber - 1)
	pm.connectVBoots()
}

func TestCsProtocolManager_connectVBoots1(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mPeer := NewMockPmAbstractPeer(ctrl)
	mPs1 := NewMockAbstractPeerSet(ctrl)

	mP2pServer := NewMockP2PServer(ctrl)

	psm := &CsPmPeerSetManager{
		verifierBootNode: mPs1,
	}

	pm := &CsProtocolManager{
		peerSetManager: psm,
		CsProtocolManagerConfig: &CsProtocolManagerConfig{
			P2PServer: mP2pServer,
		},
	}

	mPs1.EXPECT().Len().Return(2)

	var pyRecord, _ = hex.DecodeString("f884b8407098ad865b00a582051940cb9cf36836572411a47278783077011599ed5cd16b76f2635f4e234738f30813a89eb9137e3e3df5266e3a1f11df72ecf1145ccb9c01826964827634826970847f00000189736563703235366b31a103ca634cae0d49acb401d8a4c6b6fe8c55b70d115bf400769cc1400f3258cd31388375647082765f")

	var r enr.Record
	if err := rlp.DecodeBytes(pyRecord, &r); err != nil {
		t.Fatalf("can't decode: %v", err)
	}
	n, err := enode.New(enode.ValidSchemes, &r)
	if err != nil {
		t.Fatalf("can't verify record: %v", err)
	}

	_, _ = tests.ChangeVerBootNodeAddress()

	mP2pServer.EXPECT().Self().Return(n)

	mPs1.EXPECT().Peer(gomock.Any()).Return(mPeer).Times(chain_config.GetChainConfig().VerifierBootNodeNumber)

	pm.connectVBoots()
}

func TestCsProtocolManager_connectVBoots2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	//mPeer := NewMockPmAbstractPeer(ctrl)
	mPs1 := NewMockAbstractPeerSet(ctrl)

	mP2pServer := NewMockP2PServer(ctrl)

	psm := &CsPmPeerSetManager{
		verifierBootNode: mPs1,
	}

	pm := &CsProtocolManager{
		peerSetManager: psm,
		CsProtocolManagerConfig: &CsProtocolManagerConfig{
			P2PServer: mP2pServer,
		},
	}

	mPs1.EXPECT().Len().Return(2)

	var pyRecord, _ = hex.DecodeString("f884b8407098ad865b00a582051940cb9cf36836572411a47278783077011599ed5cd16b76f2635f4e234738f30813a89eb9137e3e3df5266e3a1f11df72ecf1145ccb9c01826964827634826970847f00000189736563703235366b31a103ca634cae0d49acb401d8a4c6b6fe8c55b70d115bf400769cc1400f3258cd31388375647082765f")

	var r enr.Record
	if err := rlp.DecodeBytes(pyRecord, &r); err != nil {
		t.Fatalf("can't decode: %v", err)
	}
	n, err := enode.New(enode.ValidSchemes, &r)
	if err != nil {
		t.Fatalf("can't verify record: %v", err)
	}

	_, _ = tests.ChangeVerBootNodeAddress()

	chain_config.VerifierBootNodes = append(chain_config.VerifierBootNodes, n)

	mP2pServer.EXPECT().Self().Return(n)

	mPs1.EXPECT().Peer(gomock.Any()).Return(nil).AnyTimes()

	mP2pServer.EXPECT().AddPeer(gomock.Any()).AnyTimes()

	pm.connectVBoots()
}

func TestCsProtocolManager_chainHeightTooLow(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// node conf
	mockNodeConf := NewMockNodeConf(ctrl)
	mockNodeConf.EXPECT().GetNodeType().Return(chain_config.NodeTypeOfVerifier).AnyTimes()
	mockNodeConf.EXPECT().GetNodeName().Return("aaaa").AnyTimes()

	// chain
	mockChain := NewMockChain(ctrl)

	// p2p server
	mockP2pServer := NewMockP2PServer(ctrl)

	// verifier
	mockVerifiersReader := NewMockVerifiersReader(ctrl)

	// pbft node
	mockPbftNode := NewMockPbftNode(ctrl)

	// PbftSigner
	mockPbftSigner := NewMockPbftSigner(ctrl)

	cfg := &CsProtocolManagerConfig{
		ChainConfig:     *chain_config.GetChainConfig(),
		Chain:           mockChain,
		P2PServer:       mockP2pServer,
		NodeConf:        mockNodeConf,
		VerifiersReader: mockVerifiersReader,
		PbftNode:        mockPbftNode,
		MsgSigner:       mockPbftSigner,
	}

	// mock peer set
	mockPeerSetBasePeers := NewMockAbstractPeerSet(ctrl)
	mockPeerSetCurrentVerifierPeers := NewMockAbstractPeerSet(ctrl)
	mockPeerSetnNextVerifierPeers := NewMockAbstractPeerSet(ctrl)
	mockPeerSetVerifierBootNode := NewMockAbstractPeerSet(ctrl)

	// mock peer
	mockPeer := NewMockPmAbstractPeer(ctrl)
	mockPeer.EXPECT().GetHead().Return(common.HexToHash("s"), uint64(11)).AnyTimes()
	mockPeerSetBasePeers.EXPECT().BestPeer().Return(mockPeer).AnyTimes()

	// mock peer
	mockPeer1 := NewMockPmAbstractPeer(ctrl)
	mockPeer1.EXPECT().GetHead().Return(common.HexToHash("s"), uint64(10)).AnyTimes()
	mockPeerSetCurrentVerifierPeers.EXPECT().BestPeer().Return(mockPeer1).AnyTimes()

	// mock peer
	mockPeer2 := NewMockPmAbstractPeer(ctrl)
	mockPeer2.EXPECT().GetHead().Return(common.HexToHash("s"), uint64(9)).AnyTimes()
	mockPeerSetnNextVerifierPeers.EXPECT().BestPeer().Return(mockPeer2).AnyTimes()

	// mock peer
	mockPeer3 := NewMockPmAbstractPeer(ctrl)
	mockPeer3.EXPECT().GetHead().Return(common.HexToHash("s"), uint64(8)).AnyTimes()
	mockPeerSetVerifierBootNode.EXPECT().BestPeer().Return(mockPeer3).AnyTimes()

	psm := &CsPmPeerSetManager{
		basePeers:            mockPeerSetBasePeers,
		currentVerifierPeers: mockPeerSetCurrentVerifierPeers,
		nextVerifierPeers:    mockPeerSetnNextVerifierPeers,
		verifierBootNode:     mockPeerSetVerifierBootNode,
	}

	pm := &CsProtocolManager{
		CsProtocolManagerConfig: cfg,
		peerSetManager:          psm,
	}

	block := model.NewBlock(model.NewHeader(11, 3, common.HexToHash("ss"), common.HexToHash("fdfs"), common.StringToDiff("0x22"), big.NewInt(111), common.StringToAddress("fdsfds"), common.EncodeNonce(33)), nil, nil)

	mockChain.EXPECT().CurrentBlock().Return(block)

	assert.Equal(t, true, pm.chainHeightTooLow())
}

func TestCsProtocolManager_chainHeightTooLow1(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// node conf
	mockNodeConf := NewMockNodeConf(ctrl)
	mockNodeConf.EXPECT().GetNodeType().Return(chain_config.NodeTypeOfVerifier).AnyTimes()
	mockNodeConf.EXPECT().GetNodeName().Return("aaaa").AnyTimes()

	// chain
	mockChain := NewMockChain(ctrl)

	// p2p server
	mockP2pServer := NewMockP2PServer(ctrl)

	// verifier
	mockVerifiersReader := NewMockVerifiersReader(ctrl)

	// pbft node
	mockPbftNode := NewMockPbftNode(ctrl)

	// PbftSigner
	mockPbftSigner := NewMockPbftSigner(ctrl)

	cfg := &CsProtocolManagerConfig{
		ChainConfig:     *chain_config.GetChainConfig(),
		Chain:           mockChain,
		P2PServer:       mockP2pServer,
		NodeConf:        mockNodeConf,
		VerifiersReader: mockVerifiersReader,
		PbftNode:        mockPbftNode,
		MsgSigner:       mockPbftSigner,
	}

	// mock peer set
	mockPeerSetBasePeers := NewMockAbstractPeerSet(ctrl)
	mockPeerSetCurrentVerifierPeers := NewMockAbstractPeerSet(ctrl)
	mockPeerSetnNextVerifierPeers := NewMockAbstractPeerSet(ctrl)
	mockPeerSetVerifierBootNode := NewMockAbstractPeerSet(ctrl)

	// mock peer
	mockPeer := NewMockPmAbstractPeer(ctrl)
	mockPeer.EXPECT().GetHead().Return(common.HexToHash("s"), uint64(11)).AnyTimes()
	mockPeerSetBasePeers.EXPECT().BestPeer().Return(mockPeer).AnyTimes()

	// mock peer
	mockPeer1 := NewMockPmAbstractPeer(ctrl)
	mockPeer1.EXPECT().GetHead().Return(common.HexToHash("s"), uint64(10)).AnyTimes()
	mockPeerSetCurrentVerifierPeers.EXPECT().BestPeer().Return(mockPeer1).AnyTimes()

	// mock peer
	mockPeer2 := NewMockPmAbstractPeer(ctrl)
	mockPeer2.EXPECT().GetHead().Return(common.HexToHash("s"), uint64(9)).AnyTimes()
	mockPeerSetnNextVerifierPeers.EXPECT().BestPeer().Return(mockPeer2).AnyTimes()

	// mock peer
	mockPeer3 := NewMockPmAbstractPeer(ctrl)
	mockPeer3.EXPECT().GetHead().Return(common.HexToHash("s"), uint64(8)).AnyTimes()
	mockPeerSetVerifierBootNode.EXPECT().BestPeer().Return(mockPeer3).AnyTimes()

	psm := &CsPmPeerSetManager{
		basePeers:            mockPeerSetBasePeers,
		currentVerifierPeers: mockPeerSetCurrentVerifierPeers,
		nextVerifierPeers:    mockPeerSetnNextVerifierPeers,
		verifierBootNode:     mockPeerSetVerifierBootNode,
	}

	pm := &CsProtocolManager{
		CsProtocolManagerConfig: cfg,
		peerSetManager:          psm,
	}

	block := model.NewBlock(model.NewHeader(11, 13, common.HexToHash("ss"), common.HexToHash("fdfs"), common.StringToDiff("0x22"), big.NewInt(111), common.StringToAddress("fdsfds"), common.EncodeNonce(33)), nil, nil)

	mockChain.EXPECT().CurrentBlock().Return(block)

	assert.Equal(t, false, pm.chainHeightTooLow())
}
