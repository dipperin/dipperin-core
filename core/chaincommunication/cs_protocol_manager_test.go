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
	p2pServer.EXPECT().Self().Return(&enode.Node{})

	nodeConf := NewMockNodeConf(ctrl)
	nodeConf.EXPECT().GetNodeName().Return("test")
	nodeConf.EXPECT().GetNodeType().Return(1)

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
				ps.EXPECT().Peer(peerId).Return(abstractPeer)

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
				basePs.EXPECT().Peer(peerId).Return(nil)

				currPs := NewMockAbstractPeerSet(ctrl)
				currPs.EXPECT().Peer(peerId).Return(abstractPeer)

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
				basePs.EXPECT().Peer(peerId).Return(nil)

				currPs := NewMockAbstractPeerSet(ctrl)
				currPs.EXPECT().Peer(peerId).Return(nil)

				nextPs := NewMockAbstractPeerSet(ctrl)
				nextPs.EXPECT().Peer(peerId).Return(abstractPeer)

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
				basePs.EXPECT().Peer(peerId).Return(nil)

				currPs := NewMockAbstractPeerSet(ctrl)
				currPs.EXPECT().Peer(peerId).Return(nil)

				nextPs := NewMockAbstractPeerSet(ctrl)
				nextPs.EXPECT().Peer(peerId).Return(nil)

				bootPs := NewMockAbstractPeerSet(ctrl)
				bootPs.EXPECT().Peer(peerId).Return(abstractPeer)

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
	basePs.EXPECT().GetPeers().Return(map[string]PmAbstractPeer{"base": baseP})

	currP := NewMockPmAbstractPeer(ctrl)
	currPs := NewMockAbstractPeerSet(ctrl)
	currPs.EXPECT().GetPeers().Return(map[string]PmAbstractPeer{"curr": currP})

	nextP := NewMockPmAbstractPeer(ctrl)
	nextPs := NewMockAbstractPeerSet(ctrl)
	nextPs.EXPECT().GetPeers().Return(map[string]PmAbstractPeer{"next": nextP})

	bootP := NewMockPmAbstractPeer(ctrl)
	bootPs := NewMockAbstractPeerSet(ctrl)
	bootPs.EXPECT().GetPeers().Return(map[string]PmAbstractPeer{"boot": bootP})

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
	nodeConf.EXPECT().GetNodeType().Return(chainconfig.NodeTypeOfNormal)

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
				ps.EXPECT().GetPeers().Return(map[string]PmAbstractPeer{})
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
				p.EXPECT().RemoteVerifierAddress().Return(addr)
				p.EXPECT().SendMsg(gomock.Any(), gomock.Any()).Return(nil)
				p.EXPECT().NodeName().Return("test peer")

				ps := NewMockAbstractPeerSet(ctrl)
				ps.EXPECT().GetPeers().Return(map[string]PmAbstractPeer{
					"test peer": p,
				})
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
	mVR.EXPECT().CurrentVerifiers().Return([]common.Address{address})
	mVR.EXPECT().NextVerifiers().Return([]common.Address{address})
	mVR.EXPECT().ShouldChangeVerifier().Return(false)

	mPbftS := NewMockPbftSigner(ctrl)
	mPbftS.EXPECT().GetAddress().Return(address)

	mNodeConf := NewMockNodeConf(ctrl)
	mNodeConf.EXPECT().GetNodeType().Return(verifier)

	mPn := NewMockPbftNode(ctrl)
	mPn.EXPECT().ChangePrimary(gomock.Any())

	pm := &CsProtocolManager{CsProtocolManagerConfig: &CsProtocolManagerConfig{VerifiersReader: mVR, MsgSigner: mPbftS, NodeConf: mNodeConf, PbftNode: mPn}}

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
	currPs.EXPECT().Len().Return(currLen)

	nextPs := NewMockAbstractPeerSet(ctrl)
	nextPs.EXPECT().Len().Return(nextLen)

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

func TestCsProtocolManager_ConnectPeer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	node := &enode.Node{}

	p2pServer := NewMockP2PServer(ctrl)
	p2pServer.EXPECT().AddPeer(node)

	pm := &CsProtocolManager{CsProtocolManagerConfig: &CsProtocolManagerConfig{
		P2PServer: p2pServer,
	}}

	pm.ConnectPeer(node)
}

func TestCsProtocolManager_MatchCurrentVerifiersToNext(t *testing.T) {
	// TODO
}

func TestCsProtocolManager_pickNextVerifierFromPs(t *testing.T) {
	// TODO
}

func TestCsProtocolManager_pickNextVerifierFromCps(t *testing.T) {
	// TODO
}

func TestCsProtocolManager_selfPmType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCases := []struct {
		name   string
		given  func() *CsProtocolManager
		expect int
	}{
		{
			name: "get type by cache",
			given: func() *CsProtocolManager {
				pm := &CsProtocolManager{}
				pm.pmType.Store(15)
				return pm
			},
			expect: 15,
		},
		{
			name: "get base",
			given: func() *CsProtocolManager {
				nodeConf := NewMockNodeConf(ctrl)
				nodeConf.EXPECT().GetNodeType().Return(chainconfig.NodeTypeOfNormal)

				return &CsProtocolManager{
					CsProtocolManagerConfig: &CsProtocolManagerConfig{
						NodeConf: nodeConf,
					},
				}
			},
			expect: base,
		},
		{
			name: "get verifier",
			given: func() *CsProtocolManager {
				nodeConf := NewMockNodeConf(ctrl)
				nodeConf.EXPECT().GetNodeType().Return(chainconfig.NodeTypeOfVerifier)

				return &CsProtocolManager{
					CsProtocolManagerConfig: &CsProtocolManagerConfig{
						NodeConf: nodeConf,
					},
				}
			},
			expect: verifier,
		},
		{
			name: "get boot",
			given: func() *CsProtocolManager {
				node := &enode.Node{}

				p2pServer := NewMockP2PServer(ctrl)
				p2pServer.EXPECT().Self().Return(node)

				nodeConf := NewMockNodeConf(ctrl)
				nodeConf.EXPECT().GetNodeType().Return(chainconfig.NodeTypeOfVerifierBoot)

				return &CsProtocolManager{
					CsProtocolManagerConfig: &CsProtocolManagerConfig{
						NodeConf:  nodeConf,
						P2PServer: p2pServer,
					},
					verifierBootNodes: []*enode.Node{node},
				}
			},
			expect: boot,
		},
	}

	for _, tc := range testCases {
		spt := tc.given().selfPmType()
		if !assert.Equal(t, tc.expect, spt) {
			t.Errorf("case: %s, expect:%d, got:%d", tc.name, tc.expect, spt)
		}
	}
}

func TestCsProtocolManager_isVerifierBootNode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCases := []struct {
		name   string
		given  func() (*CsProtocolManager, PmAbstractPeer)
		expect bool
	}{
		{
			name: "peer is a verifier boot node",
			given: func() (*CsProtocolManager, PmAbstractPeer) {
				node := &enode.Node{}

				p := NewMockPmAbstractPeer(ctrl)
				p.EXPECT().ID().Return(node.ID().String())

				return &CsProtocolManager{
					verifierBootNodes: []*enode.Node{node},
				}, p
			},
			expect: true,
		},
		{
			name: "peer isn't a verifier boot node",
			given: func() (*CsProtocolManager, PmAbstractPeer) {
				node := &enode.Node{}

				p := NewMockPmAbstractPeer(ctrl)
				p.EXPECT().ID().Return("test id")

				return &CsProtocolManager{
					verifierBootNodes: []*enode.Node{node},
				}, p
			},
			expect: false,
		},
	}

	for _, tc := range testCases {
		pm, p := tc.given()
		if !assert.Equal(t, tc.expect, pm.isVerifierBootNode(p)) {
			t.Errorf("case:%s, expect:%v", tc.name, tc.expect)
		}
	}
}

func TestCsProtocolManager_isCurrentVerifierNode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCases := []struct {
		name   string
		given  func() (*CsProtocolManager, PmAbstractPeer)
		expect bool
	}{
		{
			name: "peer is the current round verifier",
			given: func() (*CsProtocolManager, PmAbstractPeer) {
				addr := common.Address{1, 2, 3}

				vr := NewMockVerifiersReader(ctrl)
				vr.EXPECT().ShouldChangeVerifier().Return(false)
				vr.EXPECT().CurrentVerifiers().Return([]common.Address{addr})

				p := NewMockPmAbstractPeer(ctrl)
				p.EXPECT().RemoteVerifierAddress().Return(addr)

				return &CsProtocolManager{
					CsProtocolManagerConfig: &CsProtocolManagerConfig{
						VerifiersReader: vr,
					},
				}, p
			},
			expect: true,
		},
		{
			name: "peer isn't the current round verifier",
			given: func() (*CsProtocolManager, PmAbstractPeer) {
				addr := common.Address{1, 2, 3}

				vr := NewMockVerifiersReader(ctrl)
				vr.EXPECT().ShouldChangeVerifier().Return(true)
				vr.EXPECT().NextVerifiers().Return([]common.Address{addr})

				p := NewMockPmAbstractPeer(ctrl)
				p.EXPECT().RemoteVerifierAddress().Return(common.Address{})

				return &CsProtocolManager{
					CsProtocolManagerConfig: &CsProtocolManagerConfig{
						VerifiersReader: vr,
					},
				}, p
			},
			expect: false,
		},
	}

	for _, tc := range testCases {
		pm, p := tc.given()
		if !assert.Equal(t, tc.expect, pm.isCurrentVerifierNode(p)) {
			t.Errorf("case:%s, expect:%v", tc.name, tc.expect)
		}
	}
}

func TestCsProtocolManager_SelfIsCurrentVerifier(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCases := []struct {
		name   string
		given  func() *CsProtocolManager
		expect bool
	}{
		{
			name: "self is current verifier",
			given: func() *CsProtocolManager {
				addr := common.Address{1, 2, 3}

				vr := NewMockVerifiersReader(ctrl)
				vr.EXPECT().ShouldChangeVerifier().Return(false)
				vr.EXPECT().CurrentVerifiers().Return([]common.Address{addr})

				s := NewMockPbftSigner(ctrl)
				s.EXPECT().GetAddress().Return(addr)

				return &CsProtocolManager{
					CsProtocolManagerConfig: &CsProtocolManagerConfig{
						VerifiersReader: vr,
						MsgSigner:       s,
					},
				}
			},
			expect: true,
		},
		{
			name: "self isn't current verifier",
			given: func() *CsProtocolManager {

				vr := NewMockVerifiersReader(ctrl)
				vr.EXPECT().ShouldChangeVerifier().Return(true)
				vr.EXPECT().NextVerifiers().Return([]common.Address{{1, 2, 3}})

				s := NewMockPbftSigner(ctrl)
				s.EXPECT().GetAddress().Return(common.Address{2, 9, 3})

				return &CsProtocolManager{
					CsProtocolManagerConfig: &CsProtocolManagerConfig{
						VerifiersReader: vr,
						MsgSigner:       s,
					},
				}
			},
			expect: false,
		},
	}

	for _, tc := range testCases {
		if !assert.Equal(t, tc.expect, tc.given().SelfIsCurrentVerifier()) {
			t.Errorf("case:%s, expect:%v", tc.name, tc.expect)
		}
	}
}

func TestCsProtocolManager_selfIsVerifierBootNode(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCases := []struct {
		name   string
		given  func() *CsProtocolManager
		expect bool
	}{
		{
			name: "self is verifier boot node",
			given: func() *CsProtocolManager {
				addr := common.Address{1, 2, 3}

				s := NewMockPbftSigner(ctrl)
				s.EXPECT().GetAddress().Return(addr)

				chainconfig.VerBootNodeAddress = []common.Address{addr}
				return &CsProtocolManager{
					CsProtocolManagerConfig: &CsProtocolManagerConfig{
						MsgSigner: s,
					},
				}
			},
			expect: true,
		},
		{
			name: "self isn't verifier boot node",
			given: func() *CsProtocolManager {
				addr := common.Address{1, 2, 3}

				s := NewMockPbftSigner(ctrl)
				s.EXPECT().GetAddress().Return(addr)

				chainconfig.VerBootNodeAddress = []common.Address{{6, 9, 3}}
				return &CsProtocolManager{
					CsProtocolManagerConfig: &CsProtocolManagerConfig{
						MsgSigner: s,
					},
				}
			},
			expect: false,
		},
	}

	for _, tc := range testCases {
		if !assert.Equal(t, tc.expect, tc.given().selfIsVerifierBootNode()) {
			t.Errorf("case:%s, expect:%v", tc.name, tc.expect)
		}
	}
}

func TestCsProtocolManager_SelfIsNextVerifier(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCases := []struct {
		name   string
		given  func() *CsProtocolManager
		expect bool
	}{
		{
			name: "self is next verifier",
			given: func() *CsProtocolManager {
				addr := common.Address{1, 2, 3}

				vr := NewMockVerifiersReader(ctrl)
				vr.EXPECT().ShouldChangeVerifier().Return(false)
				vr.EXPECT().NextVerifiers().Return([]common.Address{addr})

				s := NewMockPbftSigner(ctrl)
				s.EXPECT().GetAddress().Return(addr)

				return &CsProtocolManager{
					CsProtocolManagerConfig: &CsProtocolManagerConfig{
						VerifiersReader: vr,
						MsgSigner:       s,
					},
				}
			},
			expect: true,
		},
		{
			name: "self isn't next verifier",
			given: func() *CsProtocolManager {

				vr := NewMockVerifiersReader(ctrl)
				vr.EXPECT().ShouldChangeVerifier().Return(true)
				vr.EXPECT().NextVerifiers().Return([]common.Address{{1, 2, 3}})

				s := NewMockPbftSigner(ctrl)
				s.EXPECT().GetAddress().Return(common.Address{1, 2, 3})

				return &CsProtocolManager{
					CsProtocolManagerConfig: &CsProtocolManagerConfig{
						VerifiersReader: vr,
						MsgSigner:       s,
					},
				}
			},
			expect: false,
		},
	}

	for _, tc := range testCases {
		if !assert.Equal(t, tc.expect, tc.given().SelfIsNextVerifier()) {
			t.Errorf("case:%s, expect:%v", tc.name, tc.expect)
		}
	}
}

func TestCsProtocolManager_checkConnCount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCases := []struct {
		name   string
		given  func() *CsProtocolManager
		expect bool // is max
	}{
		{
			name: "check the number of connections(default)",
			given: func() *CsProtocolManager {
				pm := &CsProtocolManager{}
				pm.pmType.Store(15)
				return pm
			},
			expect: false,
		},
		{
			name: "check base",
			given: func() *CsProtocolManager {
				nodeConf := NewMockNodeConf(ctrl)
				nodeConf.EXPECT().GetNodeType().Return(chainconfig.NodeTypeOfNormal)

				ps := NewMockAbstractPeerSet(ctrl)
				ps.EXPECT().Len().Return(maxPeers)

				psm := &CsPmPeerSetManager{
					basePeers: ps,
					maxPeers:  maxPeers,
				}

				return &CsProtocolManager{
					CsProtocolManagerConfig: &CsProtocolManagerConfig{
						NodeConf: nodeConf,
					},
					peerSetManager: psm,
				}
			},
			expect: true,
		},
		{
			name: "check verifier",
			given: func() *CsProtocolManager {
				nodeConf := NewMockNodeConf(ctrl)
				nodeConf.EXPECT().GetNodeType().Return(chainconfig.NodeTypeOfVerifier)

				basePs := NewMockAbstractPeerSet(ctrl)
				basePs.EXPECT().Len().Return(maxPeers)

				currPs := NewMockAbstractPeerSet(ctrl)
				currPs.EXPECT().Len().Return(PbftMaxPeerCount)

				psm := &CsPmPeerSetManager{
					basePeers:            basePs,
					currentVerifierPeers: currPs,
					maxPeers:             maxPeers,
				}

				return &CsProtocolManager{
					CsProtocolManagerConfig: &CsProtocolManagerConfig{
						NodeConf: nodeConf,
					},
					peerSetManager: psm,
				}
			},
			expect: true,
		},
		{
			name: "check boot",
			given: func() *CsProtocolManager {
				node := &enode.Node{}

				p2pServer := NewMockP2PServer(ctrl)
				p2pServer.EXPECT().Self().Return(node)

				nodeConf := NewMockNodeConf(ctrl)
				nodeConf.EXPECT().GetNodeType().Return(chainconfig.NodeTypeOfVerifierBoot)

				basePs := NewMockAbstractPeerSet(ctrl)
				basePs.EXPECT().Len().Return(maxPeers)

				currPs := NewMockAbstractPeerSet(ctrl)
				currPs.EXPECT().Len().Return(PbftMaxPeerCount)

				psm := &CsPmPeerSetManager{
					basePeers:            basePs,
					currentVerifierPeers: currPs,
					maxPeers:             maxPeers,
				}

				return &CsProtocolManager{
					CsProtocolManagerConfig: &CsProtocolManagerConfig{
						NodeConf:  nodeConf,
						P2PServer: p2pServer,
					},
					verifierBootNodes: []*enode.Node{node},
					peerSetManager:    psm,
				}
			},
			expect: true,
		},
	}

	for _, tc := range testCases {
		if !assert.Equal(t, tc.expect, tc.given().checkConnCount()) {
			t.Errorf("case: %s, expect:%v", tc.name, tc.expect)
		}
	}
}

func TestCsProtocolManager_handle(t *testing.T) {
	// TODO
}

func TestCsProtocolManager_handleMsg(t *testing.T) {
	// TODO
}

func TestCsProtocolManager_HandShake(t *testing.T) {
	// TODO
}

func TestCsProtocolManager_checkAndHandleVerBootNodes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCases := []struct {
		name  string
		given func() *CsProtocolManager
		// expect
	}{
		{
			name: "protocol type is base",
			given: func() *CsProtocolManager {
				nodeConf := NewMockNodeConf(ctrl)
				nodeConf.EXPECT().GetNodeType().Return(chainconfig.NodeTypeOfMineMaster)

				return &CsProtocolManager{
					CsProtocolManagerConfig: &CsProtocolManagerConfig{
						NodeConf: nodeConf,
					},
				}
			},
		},
		{
			name: "chain height too low",
			given: func() *CsProtocolManager {
				nodeConf := NewMockNodeConf(ctrl)
				nodeConf.EXPECT().GetNodeType().Return(chainconfig.NodeTypeOfVerifier)

				basePs := NewMockAbstractPeerSet(ctrl)
				basePs.EXPECT().BestPeer().Return(nil)
				currPs := NewMockAbstractPeerSet(ctrl)
				currPs.EXPECT().BestPeer().Return(nil)
				nextPs := NewMockAbstractPeerSet(ctrl)
				nextPs.EXPECT().BestPeer().Return(nil)
				bootPs := NewMockAbstractPeerSet(ctrl)
				bootPs.EXPECT().BestPeer().Return(nil)

				psm := &CsPmPeerSetManager{
					basePeers:            basePs,
					currentVerifierPeers: currPs,
					nextVerifierPeers:    nextPs,
					verifierBootNode:     bootPs,
				}

				return &CsProtocolManager{
					CsProtocolManagerConfig: &CsProtocolManagerConfig{
						NodeConf: nodeConf,
					},
					peerSetManager: psm,
				}
			},
		},
		{
			name: "self is current verifier",
			given: func() *CsProtocolManager {
				nodeConf := NewMockNodeConf(ctrl)
				nodeConf.EXPECT().GetNodeType().Return(chainconfig.NodeTypeOfVerifier)

				pNum := uint64(15)

				p := NewMockPmAbstractPeer(ctrl)
				p.EXPECT().GetHead().Return(common.Hash{}, pNum)

				basePs := NewMockAbstractPeerSet(ctrl)
				basePs.EXPECT().BestPeer().Return(nil)
				currPs := NewMockAbstractPeerSet(ctrl)
				currPs.EXPECT().BestPeer().Return(nil)
				nextPs := NewMockAbstractPeerSet(ctrl)
				nextPs.EXPECT().BestPeer().Return(nil)
				bootPs := NewMockAbstractPeerSet(ctrl)
				bootPs.EXPECT().BestPeer().Return(p)
				bootPs.EXPECT().Len().Return(chainconfig.GetChainConfig().VerifierBootNodeNumber - 1)

				psm := &CsPmPeerSetManager{
					basePeers:            basePs,
					currentVerifierPeers: currPs,
					nextVerifierPeers:    nextPs,
					verifierBootNode:     bootPs,
				}

				addr := common.Address{1, 2, 3}

				vr := NewMockVerifiersReader(ctrl)
				vr.EXPECT().ShouldChangeVerifier().Return(false)
				vr.EXPECT().CurrentVerifiers().Return([]common.Address{addr})

				s := NewMockPbftSigner(ctrl)
				s.EXPECT().GetAddress().Return(addr)

				block := NewMockAbstractBlock(ctrl)
				block.EXPECT().Number().Return(pNum)

				c := NewMockChain(ctrl)
				c.EXPECT().CurrentBlock().Return(block)

				return &CsProtocolManager{
					CsProtocolManagerConfig: &CsProtocolManagerConfig{
						NodeConf:        nodeConf,
						VerifiersReader: vr,
						MsgSigner:       s,
						Chain:           c,
					},
					peerSetManager: psm,
				}
			},
		},
	}

	for _, tc := range testCases {
		if !assert.NotPanics(t, tc.given().checkAndHandleVerBootNodes) {
			t.Errorf("case: %s", tc.name)
		}
	}
}

func TestCsProtocolManager_bootVerifierConnCheck(t *testing.T) {
	// TODO
}

func TestCsProtocolManager_RegisterCommunicationService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cs := NewMockCommunicationService(ctrl)
	cs.EXPECT().MsgHandlers().Return(map[uint64]func(msg p2p.Msg, p PmAbstractPeer) error{
		33: func(msg p2p.Msg, p PmAbstractPeer) error {
			return nil
		},
	})

	ce := NewMockCommunicationExecutable(ctrl)

	pm := &CsProtocolManager{BaseProtocolManager: BaseProtocolManager{
		msgHandlers: make(map[uint64]func(msg p2p.Msg, p PmAbstractPeer) error),
	}}

	pm.RegisterCommunicationService(cs, ce)
}

func TestCsProtocolManager_GetCurrentVerifierPeers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	key := "test peer"
	p := NewMockPmAbstractPeer(ctrl)

	ps := NewMockAbstractPeerSet(ctrl)
	ps.EXPECT().GetPeers().Return(map[string]PmAbstractPeer{key: p})

	pm := &CsProtocolManager{peerSetManager: &CsPmPeerSetManager{
		currentVerifierPeers: ps,
	}}

	assert.Equal(t, p, pm.GetCurrentVerifierPeers()[key])
}

func TestCsProtocolManager_IsSync(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCases := []struct {
		name   string
		given  func() *CsProtocolManager
		expect bool
	}{
		{
			name: "current block is nil",
			given: func() *CsProtocolManager {
				c := NewMockChain(ctrl)
				c.EXPECT().CurrentBlock().Return(nil)

				return &CsProtocolManager{CsProtocolManagerConfig: &CsProtocolManagerConfig{
					Chain: c,
				}}
			},
			expect: true,
		},
		{
			name: "best peer is nil",
			given: func() *CsProtocolManager {
				block := NewMockAbstractBlock(ctrl)

				c := NewMockChain(ctrl)
				c.EXPECT().CurrentBlock().Return(block)

				basePs := NewMockAbstractPeerSet(ctrl)
				basePs.EXPECT().BestPeer().Return(nil)

				currPs := NewMockAbstractPeerSet(ctrl)
				currPs.EXPECT().BestPeer().Return(nil)

				nextPs := NewMockAbstractPeerSet(ctrl)
				nextPs.EXPECT().BestPeer().Return(nil)

				bootPs := NewMockAbstractPeerSet(ctrl)
				bootPs.EXPECT().BestPeer().Return(nil)

				return &CsProtocolManager{CsProtocolManagerConfig: &CsProtocolManagerConfig{
					Chain: c,
				},
					peerSetManager: &CsPmPeerSetManager{
						basePeers:            basePs,
						currentVerifierPeers: currPs,
						nextVerifierPeers:    nextPs,
						verifierBootNode:     bootPs,
					}}
			},
			expect: false,
		},
		{
			name: "peer current block number + 10 > best peer height",
			given: func() *CsProtocolManager {
				height := uint64(33)

				block := NewMockAbstractBlock(ctrl)
				block.EXPECT().Number().Return(height)

				c := NewMockChain(ctrl)
				c.EXPECT().CurrentBlock().Return(block)

				p := NewMockPmAbstractPeer(ctrl)
				p.EXPECT().GetHead().Return(common.Hash{}, height)

				basePs := NewMockAbstractPeerSet(ctrl)
				basePs.EXPECT().BestPeer().Return(nil)

				currPs := NewMockAbstractPeerSet(ctrl)
				currPs.EXPECT().BestPeer().Return(nil)

				nextPs := NewMockAbstractPeerSet(ctrl)
				nextPs.EXPECT().BestPeer().Return(nil)

				bootPs := NewMockAbstractPeerSet(ctrl)
				bootPs.EXPECT().BestPeer().Return(p)

				return &CsProtocolManager{CsProtocolManagerConfig: &CsProtocolManagerConfig{
					Chain: c,
				},
					peerSetManager: &CsPmPeerSetManager{
						basePeers:            basePs,
						currentVerifierPeers: currPs,
						nextVerifierPeers:    nextPs,
						verifierBootNode:     bootPs,
					}}
			},
			expect: false,
		},
		{
			name: "peer current block number + 10 < best peer height",
			given: func() *CsProtocolManager {
				height := uint64(33)

				block := NewMockAbstractBlock(ctrl)
				block.EXPECT().Number().Return(height)

				c := NewMockChain(ctrl)
				c.EXPECT().CurrentBlock().Return(block)

				p := NewMockPmAbstractPeer(ctrl)
				p.EXPECT().GetHead().Return(common.Hash{}, height+11)

				basePs := NewMockAbstractPeerSet(ctrl)
				basePs.EXPECT().BestPeer().Return(nil)

				currPs := NewMockAbstractPeerSet(ctrl)
				currPs.EXPECT().BestPeer().Return(nil)

				nextPs := NewMockAbstractPeerSet(ctrl)
				nextPs.EXPECT().BestPeer().Return(nil)

				bootPs := NewMockAbstractPeerSet(ctrl)
				bootPs.EXPECT().BestPeer().Return(p)

				return &CsProtocolManager{CsProtocolManagerConfig: &CsProtocolManagerConfig{
					Chain: c,
				},
					peerSetManager: &CsPmPeerSetManager{
						basePeers:            basePs,
						currentVerifierPeers: currPs,
						nextVerifierPeers:    nextPs,
						verifierBootNode:     bootPs,
					}}
			},
			expect: true,
		},
	}

	for _, tc := range testCases {
		if !assert.Equal(t, tc.expect, tc.given().IsSync()) {
			t.Errorf("case:%s expect:%v", tc.name, tc.expect)
		}
	}
}
