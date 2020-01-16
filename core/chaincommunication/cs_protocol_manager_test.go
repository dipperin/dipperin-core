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
	"encoding/hex"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/chainconfig"
	"github.com/dipperin/dipperin-core/third_party/p2p"
	"github.com/dipperin/dipperin-core/third_party/p2p/enode"
	"github.com/dipperin/dipperin-core/third_party/p2p/enr"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestCsProtocolManager_ShowPmInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	p2pServer := NewMockP2PServer(ctrl)
	p2pServer.EXPECT().Self().Return(&enode.Node{}).Times(1)

	nodeConf := NewMockNodeConf(ctrl)
	nodeConf.EXPECT().GetNodeName().Return("test").Times(1)
	nodeConf.EXPECT().GetNodeType().Return(1).Times(1)

	ps := NewMockAbstractPeerSet(ctrl)
	ps.EXPECT().GetPeersInfo().Return([]*p2p.CsPeerInfo{}).Times(4)

	pm := &CsProtocolManager{
		CsProtocolManagerConfig: &CsProtocolManagerConfig{
			P2PServer: p2pServer,
			NodeConf:  nodeConf,
		},
		peerSetManager: &CsPmPeerSetManager{
			basePeers:            ps,
			currentVerifierPeers: ps,
			nextVerifierPeers:    ps,
			verifierBootNode:     ps,
		},
	}

	assert.NotNil(t, pm.ShowPmInfo())
}

func TestCsProtocolManager_BestPeer(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// node conf
	mockNodeConf := NewMockNodeConf(ctrl)
	mockNodeConf.EXPECT().GetNodeType().Return(chainconfig.NodeTypeOfVerifier).AnyTimes()
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
		ChainConfig:     *chainconfig.GetChainConfig(),
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

	// mock peer - expect
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
	mockNodeConf.EXPECT().GetNodeType().Return(chainconfig.NodeTypeOfVerifier).AnyTimes()
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
		ChainConfig:     *chainconfig.GetChainConfig(),
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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCases := []struct {
		name           string
		givenAndExpect func() (manager *CsProtocolManager, peerId string, abstractPeer PmAbstractPeer)
	}{
		{
			name: "get peer by base",
			givenAndExpect: func() (manager *CsProtocolManager, peerId string, abstractPeer PmAbstractPeer) {
				peerId = "test peer"

				abstractPeer = NewMockPmAbstractPeer(ctrl)

				ps := NewMockAbstractPeerSet(ctrl)
				ps.EXPECT().Peer(peerId).Return(abstractPeer).Times(1)

				psm := &CsPmPeerSetManager{
					basePeers: ps,
				}
				manager = &CsProtocolManager{peerSetManager: psm}
				return
			},
		},
		{
			name: "get peer by currentVerifier",
			givenAndExpect: func() (manager *CsProtocolManager, peerId string, abstractPeer PmAbstractPeer) {
				peerId = "test peer"

				abstractPeer = NewMockPmAbstractPeer(ctrl)

				basePs := NewMockAbstractPeerSet(ctrl)
				basePs.EXPECT().Peer(peerId).Return(nil).Times(1)

				currPs := NewMockAbstractPeerSet(ctrl)
				currPs.EXPECT().Peer(peerId).Return(abstractPeer).Times(1)

				psm := &CsPmPeerSetManager{
					basePeers:            basePs,
					currentVerifierPeers: currPs,
				}
				manager = &CsProtocolManager{peerSetManager: psm}
				return
			},
		},
		{
			name: "get peer by nextVerifier",
			givenAndExpect: func() (manager *CsProtocolManager, peerId string, abstractPeer PmAbstractPeer) {
				peerId = "test peer"

				abstractPeer = NewMockPmAbstractPeer(ctrl)

				basePs := NewMockAbstractPeerSet(ctrl)
				basePs.EXPECT().Peer(peerId).Return(nil).Times(1)

				currPs := NewMockAbstractPeerSet(ctrl)
				currPs.EXPECT().Peer(peerId).Return(nil).Times(1)

				nextPs := NewMockAbstractPeerSet(ctrl)
				nextPs.EXPECT().Peer(peerId).Return(abstractPeer).Times(1)

				psm := &CsPmPeerSetManager{
					basePeers:            basePs,
					currentVerifierPeers: currPs,
					nextVerifierPeers:    nextPs,
				}
				manager = &CsProtocolManager{peerSetManager: psm}
				return
			},
		},
		{
			name: "get peer by verifierBootNode",
			givenAndExpect: func() (manager *CsProtocolManager, peerId string, abstractPeer PmAbstractPeer) {
				peerId = "test peer"

				abstractPeer = NewMockPmAbstractPeer(ctrl)

				basePs := NewMockAbstractPeerSet(ctrl)
				basePs.EXPECT().Peer(peerId).Return(nil).Times(1)

				currPs := NewMockAbstractPeerSet(ctrl)
				currPs.EXPECT().Peer(peerId).Return(nil).Times(1)

				nextPs := NewMockAbstractPeerSet(ctrl)
				nextPs.EXPECT().Peer(peerId).Return(nil).Times(1)

				bootPs := NewMockAbstractPeerSet(ctrl)
				bootPs.EXPECT().Peer(peerId).Return(abstractPeer).Times(1)

				psm := &CsPmPeerSetManager{
					basePeers:            basePs,
					currentVerifierPeers: currPs,
					nextVerifierPeers:    nextPs,
					verifierBootNode:     bootPs,
				}
				manager = &CsProtocolManager{peerSetManager: psm}
				return
			},
		},
	}

	for _, tc := range testCases {
		g1, g2, expect := tc.givenAndExpect()
		p := g1.GetPeer(g2)
		if !assert.Equal(t, expect, p) {
			t.Errorf("case:%s, expec:%+v, got:%+v", tc.name, expect, p)
		}
	}
}

func TestCsProtocolManager_GetPeers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	baseP := NewMockPmAbstractPeer(ctrl)
	basePs := NewMockAbstractPeerSet(ctrl)
	basePs.EXPECT().GetPeers().Return(map[string]PmAbstractPeer{"base": baseP}).Times(1)

	currP := NewMockPmAbstractPeer(ctrl)
	currPs := NewMockAbstractPeerSet(ctrl)
	currPs.EXPECT().GetPeers().Return(map[string]PmAbstractPeer{"curr": currP}).Times(1)

	nextP := NewMockPmAbstractPeer(ctrl)
	nextPs := NewMockAbstractPeerSet(ctrl)
	nextPs.EXPECT().GetPeers().Return(map[string]PmAbstractPeer{"next": nextP}).Times(1)

	bootP := NewMockPmAbstractPeer(ctrl)
	bootPs := NewMockAbstractPeerSet(ctrl)
	bootPs.EXPECT().GetPeers().Return(map[string]PmAbstractPeer{"boot": bootP}).Times(1)

	pm := &CsProtocolManager{
		peerSetManager: &CsPmPeerSetManager{
			basePeers:            basePs,
			currentVerifierPeers: currPs,
			nextVerifierPeers:    nextPs,
			verifierBootNode:     bootPs,
		},
	}

	result := pm.GetPeers()

	assert.Equal(t, baseP, result["base"])
	assert.Equal(t, currP, result["curr"])
	assert.Equal(t, nextP, result["next"])
	assert.Equal(t, bootP, result["boot"])
}

func Test_mergePeersTo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPeer1 := NewMockPmAbstractPeer(ctrl)

	fm := map[string]PmAbstractPeer{"1": mockPeer1}

	tm := make(map[string]PmAbstractPeer)

	mergePeersTo(fm, tm)

	assert.Equal(t, mockPeer1, tm["1"])
}

func TestCsProtocolManager_Protocols(t *testing.T) {
	testCases := []struct {
		name   string
		given  func() *CsProtocolManager
		expect int // protocol count
	}{
		{
			name: "CsProtocolManager.protocols isn't empty",
			given: func() *CsProtocolManager {
				return &CsProtocolManager{
					BaseProtocolManager: BaseProtocolManager{
						protocols: []p2p.Protocol{{Name: "p1"}, {Name: "p2"}},
					},
				}
			},
			expect: 2,
		},
		{
			name: "CsProtocolManager.protocols is empty",
			given: func() *CsProtocolManager {
				return &CsProtocolManager{
					BaseProtocolManager: BaseProtocolManager{
						protocols: []p2p.Protocol{},
					},
				}
			},
			expect: 1,
		},
	}

	for _, tc := range testCases {
		psLen := len(tc.given().Protocols())
		if !assert.Equal(t, tc.expect, psLen) {
			t.Errorf("case:%s, expect:%d, got:%d", tc.name, tc.expect, psLen)
		}
	}
}

func TestCsProtocolManager_getCsProtocol(t *testing.T) {
	testCases := []struct {
		name   string
		given  func() // set boot env
		expect string // Protocol.Name
	}{
		{
			name: "boot env is BootEnvMercury",
			given: func() {
				os.Setenv(chainconfig.BootEnvTagName, chainconfig.BootEnvMercury)
			},
			expect: chainconfig.AppName + "_cs",
		},
		{
			name: "boot env is BootEnvTest",
			given: func() {
				os.Setenv(chainconfig.BootEnvTagName, chainconfig.BootEnvTest)
			},
			expect: chainconfig.AppName + "_cs_test",
		},
		{
			name: "boot env is BootEnvVenus",
			given: func() {
				os.Setenv(chainconfig.BootEnvTagName, chainconfig.BootEnvVenus)
			},
			expect: chainconfig.AppName + "_vs",
		},
		{
			name: "boot env is local(default)",
			given: func() {
				os.Setenv(chainconfig.BootEnvTagName, chainconfig.BootEnvLocal)
			},
			expect: chainconfig.AppName + "_cs_local",
		},
	}

	pm := &CsProtocolManager{}

	for _, tc := range testCases {
		tc.given()
		p := pm.getCsProtocol()
		if !assert.Equal(t, tc.expect, p.Name) {
			t.Errorf("case:%s expect:%s got:%s", tc.name, tc.expect, p.Name)
		}
	}
}

func Test_newCsProtocolManager(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	nodeConf := NewMockNodeConf(ctrl)
	nodeConf.EXPECT().GetNodeType().Return(chainconfig.NodeTypeOfNormal).Times(1)

	assert.NotNil(t, newCsProtocolManager(&CsProtocolManagerConfig{
		NodeConf: nodeConf,
	}))
}

// TODO: need check ?
//todo Start
//todo Stop
//todo BroadcastMsg
//todo BroadcastMsgToTargetVerifiers

func TestCsProtocolManager_SendFetchBlockMsg(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCases := []struct {
		name   string
		given  func() (*CsProtocolManager, common.Address)
		expect bool // exist error
	}{
		{
			name: "no verifier peer for fetcher",
			given: func() (*CsProtocolManager, common.Address) {
				ps := NewMockAbstractPeerSet(ctrl)
				ps.EXPECT().GetPeers().Return(map[string]PmAbstractPeer{}).Times(1)
				psm := &CsPmPeerSetManager{
					currentVerifierPeers: ps,
				}
				return &CsProtocolManager{
					peerSetManager: psm,
				}, common.Address{}
			},
			expect: true,
		},
		{
			name: "send success",
			given: func() (*CsProtocolManager, common.Address) {
				addr := common.Address{1, 2, 3}
				p := NewMockPmAbstractPeer(ctrl)
				p.EXPECT().RemoteVerifierAddress().Return(addr).Times(1)
				p.EXPECT().SendMsg(gomock.Any(), gomock.Any()).Return(nil).Times(1)
				p.EXPECT().NodeName().Return("test peer").Times(1)

				ps := NewMockAbstractPeerSet(ctrl)
				ps.EXPECT().GetPeers().Return(map[string]PmAbstractPeer{
					"test peer": p,
				}).Times(1)
				psm := &CsPmPeerSetManager{
					currentVerifierPeers: ps,
				}
				return &CsProtocolManager{
					peerSetManager: psm,
				}, addr
			},
			expect: false,
		},
	}

	for _, tc := range testCases {
		pm, addr := tc.given()
		err := pm.SendFetchBlockMsg(0, addr, nil)
		if !assert.Equal(t, tc.expect, err != nil) {
			t.Errorf("case:%s, expect:%v, err:%v", tc.name, tc.expect, err)
		}
	}
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

func TestCsProtocolManager_HaveEnoughVerifiers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	verNum := 10
	currLen := 1
	nextLen := 2

	currPs := NewMockAbstractPeerSet(ctrl)
	currPs.EXPECT().Len().Return(currLen).Times(1)

	nextPs := NewMockAbstractPeerSet(ctrl)
	nextPs.EXPECT().Len().Return(nextLen).Times(1)

	psm := &CsPmPeerSetManager{
		currentVerifierPeers: currPs,
		nextVerifierPeers:    nextPs,
	}

	pm := &CsProtocolManager{
		CsProtocolManagerConfig: &CsProtocolManagerConfig{
			ChainConfig: chainconfig.ChainConfig{
				VerifierNumber: verNum,
			},
		},
		peerSetManager: psm,
	}

	mc, mn := pm.HaveEnoughVerifiers(false)

	assert.Equal(t, uint(verNum-1-currLen), mc)
	assert.Equal(t, uint(verNum-1-nextLen), mn)
}

func TestCsProtocolManager_GetCurrentConnectPeers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	addr := common.HexToAddress("sss")
	addrKey := "test peer"

	mockNodeConf := NewMockNodeConf(ctrl)
	mockNodeConf.EXPECT().GetNodeType().Return(chainconfig.NodeTypeOfVerifier)

	mockPeer := NewMockPmAbstractPeer(ctrl)
	mockPeer.EXPECT().RemoteVerifierAddress().Return(addr)

	mockPeerSetCurrentVerifierPeers := NewMockAbstractPeerSet(ctrl)

	mockPeerSetCurrentVerifierPeers.EXPECT().GetPeers().Return(map[string]PmAbstractPeer{addrKey: mockPeer})

	psm := &CsPmPeerSetManager{
		currentVerifierPeers: mockPeerSetCurrentVerifierPeers,
	}

	pm := &CsProtocolManager{peerSetManager: psm, CsProtocolManagerConfig: &CsProtocolManagerConfig{NodeConf: mockNodeConf}}

	rM := pm.GetCurrentConnectPeers()

	assert.True(t, rM[addrKey].IsEqual(addr))
}

func TestCsProtocolManager_GetVerifierBootNode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPeerSet := NewMockAbstractPeerSet(ctrl)
	mockPeerSet.EXPECT().GetPeers().Return(nil)

	psm := &CsPmPeerSetManager{
		verifierBootNode: mockPeerSet,
	}

	pm := &CsProtocolManager{peerSetManager: psm}

	assert.Nil(t, pm.GetVerifierBootNode())
}

func TestCsProtocolManager_GetNextVerifierPeers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPeerSet := NewMockAbstractPeerSet(ctrl)
	mockPeerSet.EXPECT().GetPeers().Return(nil)

	psm := &CsPmPeerSetManager{
		nextVerifierPeers: mockPeerSet,
	}

	pm := &CsProtocolManager{peerSetManager: psm}

	assert.Nil(t, pm.GetNextVerifierPeers())
}

func TestCsProtocolManager_CurrentVerifierPeersSet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ps := NewMockAbstractPeerSet(ctrl)

	psm := &CsPmPeerSetManager{
		currentVerifierPeers: ps,
	}

	pm := &CsProtocolManager{
		peerSetManager: psm,
	}

	assert.Equal(t, ps, pm.CurrentVerifierPeersSet())
}

func TestCsProtocolManager_NextVerifierPeersSet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ps := NewMockAbstractPeerSet(ctrl)

	psm := &CsPmPeerSetManager{
		nextVerifierPeers: ps,
	}

	pm := &CsProtocolManager{
		peerSetManager: psm,
	}

	assert.Equal(t, ps, pm.NextVerifierPeersSet())
}

func TestCsProtocolManager_SelfIsBootNode(t *testing.T) {
	testCases := []struct {
		name   string
		given  func() *CsProtocolManager
		expect bool
	}{
		{
			name: "self is boot node",
			given: func() *CsProtocolManager {
				pm := &CsProtocolManager{}
				pm.pmType.Store(boot)
				return pm
			},
			expect: true,
		},
		{
			name: "self isn't boot node",
			given: func() *CsProtocolManager {
				pm := &CsProtocolManager{}
				pm.pmType.Store(base)
				return pm
			},
			expect: false,
		},
	}

	for _, tc := range testCases {
		if !assert.Equal(t, tc.expect, tc.given().SelfIsBootNode()) {
			t.Errorf("case:%s expect:%v", tc.name, tc.expect)
		}
	}
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

	selfNode := pm.GetSelfNode()

	assert.Equal(t, n.ID().String(), selfNode.ID().String())
}
