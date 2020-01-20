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
	"errors"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/chainconfig"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_newCsPmPeerSetManager(t *testing.T) {
	assert.NotNil(t, newCsPmPeerSetManager(1, 1,
		func() bool {
			return false
		}, func() bool {
			return false
		}, func(p PmAbstractPeer) bool {
			return false
		}, func(p PmAbstractPeer) bool {
			return false
		}, func(p PmAbstractPeer) bool {
			return false
		}))
}

func TestCsPmPeerSetManager_AddPeer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCases := []struct {
		name   string
		given  func() (*CsPmPeerSetManager, PmAbstractPeer)
		expect bool // exist error
	}{
		{
			name: "base - too many peers",
			given: func() (*CsPmPeerSetManager, PmAbstractPeer) {
				ps := NewMockAbstractPeerSet(ctrl)
				ps.EXPECT().Len().Return(maxPeers + 1)
				return &CsPmPeerSetManager{
					pmType:    base,
					maxPeers:  maxPeers,
					basePeers: ps,
				}, nil
			},
			expect: true,
		},
		{
			name: "base - add peer error",
			given: func() (*CsPmPeerSetManager, PmAbstractPeer) {
				p := NewMockPmAbstractPeer(ctrl)
				p.EXPECT().ID().Return("test")

				ps := NewMockAbstractPeerSet(ctrl)
				ps.EXPECT().Len().Return(maxPeers - 1)
				ps.EXPECT().AddPeer(p).Return(errors.New("mock error"))

				return &CsPmPeerSetManager{
					pmType:    base,
					maxPeers:  maxPeers,
					basePeers: ps,
				}, p
			},
			expect: true,
		},
		{
			name: "verifier - NodeTypeOfNormal - too many peers",
			given: func() (manager *CsPmPeerSetManager, abstractPeer PmAbstractPeer) {
				p := NewMockPmAbstractPeer(ctrl)
				p.EXPECT().NodeType().Return(uint64(chainconfig.NodeTypeOfNormal))

				ps := NewMockAbstractPeerSet(ctrl)
				ps.EXPECT().Len().Return(maxPeers)

				return &CsPmPeerSetManager{
					pmType:    verifier,
					maxPeers:  maxPeers,
					basePeers: ps,
				}, p
			},
			expect: true,
		},
		{
			name: "verifier - NodeTypeOfNormal - add peer error",
			given: func() (manager *CsPmPeerSetManager, abstractPeer PmAbstractPeer) {
				p := NewMockPmAbstractPeer(ctrl)
				p.EXPECT().NodeType().Return(uint64(chainconfig.NodeTypeOfNormal))

				ps := NewMockAbstractPeerSet(ctrl)
				ps.EXPECT().Len().Return(maxPeers - PbftMaxPeerCount - 1)
				ps.EXPECT().AddPeer(p).Return(errors.New("mock error"))

				return &CsPmPeerSetManager{
					pmType:    verifier,
					maxPeers:  maxPeers,
					basePeers: ps,
				}, p
			},
			expect: true,
		},
		{
			name: "verifier - NodeTypeOfMineMaster - too many peers",
			given: func() (manager *CsPmPeerSetManager, abstractPeer PmAbstractPeer) {
				p := NewMockPmAbstractPeer(ctrl)
				p.EXPECT().NodeType().Return(uint64(chainconfig.NodeTypeOfMineMaster))

				ps := NewMockAbstractPeerSet(ctrl)
				ps.EXPECT().Len().Return(maxPeers)

				return &CsPmPeerSetManager{
					pmType:    verifier,
					maxPeers:  maxPeers,
					basePeers: ps,
				}, p
			},
			expect: true,
		},
		{
			name: "verifier - NodeTypeOfMineMaster - add peer error",
			given: func() (manager *CsPmPeerSetManager, abstractPeer PmAbstractPeer) {
				p := NewMockPmAbstractPeer(ctrl)
				p.EXPECT().NodeType().Return(uint64(chainconfig.NodeTypeOfMineMaster))

				ps := NewMockAbstractPeerSet(ctrl)
				ps.EXPECT().Len().Return(maxPeers - PbftMaxPeerCount - 1)
				ps.EXPECT().AddPeer(p).Return(errors.New("mock error"))

				return &CsPmPeerSetManager{
					pmType:    verifier,
					maxPeers:  maxPeers,
					basePeers: ps,
				}, p
			},
			expect: true,
		},
		{
			name: "verifier - NodeTypeOfVerifier - add peer to base set",
			given: func() (manager *CsPmPeerSetManager, abstractPeer PmAbstractPeer) {
				p := NewMockPmAbstractPeer(ctrl)
				p.EXPECT().NodeType().Return(uint64(chainconfig.NodeTypeOfVerifier))

				ps := NewMockAbstractPeerSet(ctrl)
				ps.EXPECT().Len().Return(maxPeers - PbftMaxPeerCount - 1)
				ps.EXPECT().AddPeer(p).Return(nil)

				return &CsPmPeerSetManager{
					pmType:    verifier,
					maxPeers:  maxPeers,
					basePeers: ps,
					selfIsCurrentVerifier: func() bool {
						return false
					},
					selfIsNextVerifier: func() bool {
						return false
					},
				}, p
			},
			expect: false,
		},
		{
			name: "verifier - NodeTypeOfVerifier - add peer to currentVerifier set",
			given: func() (manager *CsPmPeerSetManager, abstractPeer PmAbstractPeer) {
				p := NewMockPmAbstractPeer(ctrl)
				p.EXPECT().NodeType().Return(uint64(chainconfig.NodeTypeOfVerifier))
				p.EXPECT().NodeName().Return("test")

				ps := NewMockAbstractPeerSet(ctrl)
				ps.EXPECT().Len().Return(PbftMaxPeerCount - 2)
				ps.EXPECT().AddPeer(p).Return(nil)

				return &CsPmPeerSetManager{
					pmType:               verifier,
					maxPeers:             maxPeers,
					currentVerifierPeers: ps,
					selfIsCurrentVerifier: func() bool {
						return true
					},
					selfIsNextVerifier: func() bool {
						return true
					},
					isCurrentVerifier: func(p PmAbstractPeer) bool {
						return true
					},
					isNextVerifier: func(p PmAbstractPeer) bool {
						return false
					},
				}, p
			},
			expect: false,
		},
		{
			name: "verifier - NodeTypeOfVerifier - add peer to nextVerifier set",
			given: func() (manager *CsPmPeerSetManager, abstractPeer PmAbstractPeer) {
				p := NewMockPmAbstractPeer(ctrl)
				p.EXPECT().NodeType().Return(uint64(chainconfig.NodeTypeOfVerifier))
				p.EXPECT().NodeName().Return("test")

				ps := NewMockAbstractPeerSet(ctrl)
				ps.EXPECT().Len().Return(PbftMaxPeerCount - 2)
				ps.EXPECT().AddPeer(p).Return(nil)

				return &CsPmPeerSetManager{
					pmType:            verifier,
					maxPeers:          maxPeers,
					nextVerifierPeers: ps,
					selfIsCurrentVerifier: func() bool {
						return true
					},
					selfIsNextVerifier: func() bool {
						return true
					},
					isCurrentVerifier: func(p PmAbstractPeer) bool {
						return false
					},
					isNextVerifier: func(p PmAbstractPeer) bool {
						return true
					},
				}, p
			},
			expect: false,
		},
		{
			name: "verifier - NodeTypeOfVerifierBoot - the remote peer isn't verifier boot node",
			given: func() (manager *CsPmPeerSetManager, abstractPeer PmAbstractPeer) {
				p := NewMockPmAbstractPeer(ctrl)
				p.EXPECT().NodeType().Return(uint64(chainconfig.NodeTypeOfVerifierBoot))

				return &CsPmPeerSetManager{
					pmType:   verifier,
					maxPeers: maxPeers,
					isVerifierBootNode: func(p PmAbstractPeer) bool {
						return false
					},
				}, p
			},
			expect: false,
		},
		{
			name: "verifier - NodeTypeOfVerifierBoot - the remote peer is verifier boot node",
			given: func() (manager *CsPmPeerSetManager, abstractPeer PmAbstractPeer) {
				p := NewMockPmAbstractPeer(ctrl)
				p.EXPECT().NodeType().Return(uint64(chainconfig.NodeTypeOfVerifierBoot))
				p.EXPECT().NodeName().Return("test")

				ps := NewMockAbstractPeerSet(ctrl)
				ps.EXPECT().AddPeer(p).Return(nil)

				return &CsPmPeerSetManager{
					pmType:           verifier,
					maxPeers:         maxPeers,
					verifierBootNode: ps,
					isVerifierBootNode: func(p PmAbstractPeer) bool {
						return true
					},
				}, p
			},
			expect: false,
		},
		{
			name: "boot - NodeTypeOfNormal - add peer to base set",
			given: func() (manager *CsPmPeerSetManager, abstractPeer PmAbstractPeer) {
				p := NewMockPmAbstractPeer(ctrl)
				p.EXPECT().NodeType().Return(uint64(chainconfig.NodeTypeOfNormal))

				ps := NewMockAbstractPeerSet(ctrl)
				ps.EXPECT().Len().Return(maxPeers - PbftMaxPeerCount - 1)
				ps.EXPECT().AddPeer(p).Return(nil)

				return &CsPmPeerSetManager{
					pmType:    boot,
					maxPeers:  maxPeers,
					basePeers: ps,
					isVerifierBootNode: func(p PmAbstractPeer) bool {
						return true
					},
				}, p
			},
			expect: false,
		},
		{
			name: "boot - NodeTypeOfVerifier - add peer to currentVerifier set",
			given: func() (manager *CsPmPeerSetManager, abstractPeer PmAbstractPeer) {
				p := NewMockPmAbstractPeer(ctrl)
				p.EXPECT().NodeType().Return(uint64(chainconfig.NodeTypeOfVerifier))
				p.EXPECT().NodeName().Return("test")

				ps := NewMockAbstractPeerSet(ctrl)
				ps.EXPECT().Len().Return(PbftMaxPeerCount - 1)
				ps.EXPECT().AddPeer(p).Return(nil)

				return &CsPmPeerSetManager{
					pmType:               boot,
					maxPeers:             maxPeers,
					currentVerifierPeers: ps,
					isCurrentVerifier: func(p PmAbstractPeer) bool {
						return true
					},
					isNextVerifier: func(p PmAbstractPeer) bool {
						return false
					},
				}, p
			},
			expect: false,
		},
		{
			name: "boot - NodeTypeOfVerifier - add peer to nextVerifier set",
			given: func() (manager *CsPmPeerSetManager, abstractPeer PmAbstractPeer) {
				p := NewMockPmAbstractPeer(ctrl)
				p.EXPECT().NodeType().Return(uint64(chainconfig.NodeTypeOfVerifier))
				p.EXPECT().NodeName().Return("test")

				ps := NewMockAbstractPeerSet(ctrl)
				ps.EXPECT().Len().Return(PbftMaxPeerCount - 1)
				ps.EXPECT().AddPeer(p).Return(nil)

				return &CsPmPeerSetManager{
					pmType:            boot,
					maxPeers:          maxPeers,
					nextVerifierPeers: ps,
					isCurrentVerifier: func(p PmAbstractPeer) bool {
						return false
					},
					isNextVerifier: func(p PmAbstractPeer) bool {
						return true
					},
				}, p
			},
			expect: false,
		},
		{
			name: "boot - NodeTypeOfVerifierBoot - add peer",
			given: func() (manager *CsPmPeerSetManager, abstractPeer PmAbstractPeer) {
				p := NewMockPmAbstractPeer(ctrl)
				p.EXPECT().NodeType().Return(uint64(chainconfig.NodeTypeOfVerifierBoot))
				p.EXPECT().NodeName().Return("test")

				ps := NewMockAbstractPeerSet(ctrl)
				ps.EXPECT().AddPeer(p).Return(nil)

				return &CsPmPeerSetManager{
					pmType:           boot,
					maxPeers:         maxPeers,
					verifierBootNode: ps,
					isVerifierBootNode: func(p PmAbstractPeer) bool {
						return true
					},
				}, p
			},
			expect: false,
		},
	}

	for _, tc := range testCases {
		psm, p := tc.given()
		err := psm.AddPeer(p)
		if !assert.Equal(t, tc.expect, err != nil) {
			t.Errorf("case:%s, expect: exist error(%v), got:%v", tc.name, tc.expect, err)
		}
	}
}

func TestCsPmPeerSetManager_BestPeer(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

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

	bestPeer := psm.BestPeer()

	hash, height := bestPeer.GetHead()

	assert.Equal(t, true, hash.IsEqual(common.HexToHash("s")))
	assert.Equal(t, uint64(11), height)
}

func TestCsPmPeerSetManager_RemovePeer(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// mock peer set
	mockPeerSetBasePeers := NewMockAbstractPeerSet(ctrl)
	mockPeerSetCurrentVerifierPeers := NewMockAbstractPeerSet(ctrl)
	mockPeerSetnNextVerifierPeers := NewMockAbstractPeerSet(ctrl)
	mockPeerSetVerifierBootNode := NewMockAbstractPeerSet(ctrl)

	// mock peer
	mockPeerSetBasePeers.EXPECT().Peer(gomock.Any()).Return(nil)
	mockPeerSetCurrentVerifierPeers.EXPECT().Peer(gomock.Any()).Return(nil)
	mockPeerSetnNextVerifierPeers.EXPECT().Peer(gomock.Any()).Return(nil)
	mockPeerSetVerifierBootNode.EXPECT().Peer(gomock.Any()).Return(nil)

	psm := &CsPmPeerSetManager{
		basePeers:            mockPeerSetBasePeers,
		currentVerifierPeers: mockPeerSetCurrentVerifierPeers,
		nextVerifierPeers:    mockPeerSetnNextVerifierPeers,
		verifierBootNode:     mockPeerSetVerifierBootNode,
	}

	psm.RemovePeer("test")
}

func TestCsProtocolManager_removePeer(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	peerId := "test peer"

	// mock peer
	p := NewMockPmAbstractPeer(ctrl)
	p.EXPECT().NodeName().Return(peerId)
	p.EXPECT().DisconnectPeer()

	// mock peer set
	ps := NewMockAbstractPeerSet(ctrl)
	ps.EXPECT().Peer(peerId).Return(p)
	ps.EXPECT().RemovePeer(peerId).Return(nil)

	psm := &CsPmPeerSetManager{
		basePeers: ps,
	}

	psm.removePeer(peerId)
}

func TestCsProtocolManager_removeCurrentVerifierPeers(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	peerId := "test peer"

	// mock peer
	p := NewMockPmAbstractPeer(ctrl)
	p.EXPECT().NodeName().Return(peerId)
	p.EXPECT().DisconnectPeer()

	// mock peer set
	ps := NewMockAbstractPeerSet(ctrl)
	ps.EXPECT().Peer(peerId).Return(p)
	ps.EXPECT().RemovePeer(peerId).Return(nil)

	psm := &CsPmPeerSetManager{
		currentVerifierPeers: ps,
	}

	psm.removeCurrentVerifierPeers(peerId)
}

func TestCsProtocolManager_removeNextVerifierPeers(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	peerId := "test peer"

	// mock peer
	p := NewMockPmAbstractPeer(ctrl)
	p.EXPECT().NodeName().Return(peerId)
	p.EXPECT().DisconnectPeer()

	// mock peer set
	ps := NewMockAbstractPeerSet(ctrl)
	ps.EXPECT().Peer(peerId).Return(p)
	ps.EXPECT().RemovePeer(peerId).Return(nil)

	psm := &CsPmPeerSetManager{
		nextVerifierPeers: ps,
	}

	psm.removeNextVerifierPeers(peerId)
}

func TestCsProtocolManager_removeVerifierBootNodePeers(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	peerId := "test peer"

	// mock peer
	p := NewMockPmAbstractPeer(ctrl)
	p.EXPECT().NodeName().Return(peerId)
	p.EXPECT().DisconnectPeer()

	// mock peer set
	ps := NewMockAbstractPeerSet(ctrl)
	ps.EXPECT().Peer(peerId).Return(p)
	ps.EXPECT().RemovePeer(peerId).Return(nil)

	psm := &CsPmPeerSetManager{
		verifierBootNode: ps,
	}

	psm.removeVerifierBootNodePeers(peerId)
}

func TestCsPmPeerSetManager_collectAllPeers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mPs1 := NewMockAbstractPeerSet(ctrl)
	mPs2 := NewMockAbstractPeerSet(ctrl)
	mPs3 := NewMockAbstractPeerSet(ctrl)

	psm := &CsPmPeerSetManager{
		basePeers:            mPs1,
		currentVerifierPeers: mPs2,
		nextVerifierPeers:    mPs3,
	}

	mPs1.EXPECT().GetPeers().Return(map[string]PmAbstractPeer{})
	mPs2.EXPECT().GetPeers().Return(map[string]PmAbstractPeer{})
	mPs3.EXPECT().GetPeers().Return(map[string]PmAbstractPeer{})

	assert.Len(t, psm.collectAllPeers(), 0)
}

func Test_filterPeers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mPeer1 := NewMockPmAbstractPeer(ctrl)
	mPeer2 := NewMockPmAbstractPeer(ctrl)
	mPeer3 := NewMockPmAbstractPeer(ctrl)

	addr1 := common.StringToAddress("aaa")
	addr2 := common.StringToAddress("bbb")
	addr3 := common.StringToAddress("ccc")

	from := map[string]PmAbstractPeer{
		"a": mPeer1,
		"b": mPeer2,
		"c": mPeer3,
	}

	curs := []common.Address{addr1}
	nexts := []common.Address{addr2}

	mPeer1.EXPECT().RemoteVerifierAddress().Return(addr1).Times(2)
	mPeer2.EXPECT().RemoteVerifierAddress().Return(addr2).Times(2)
	mPeer3.EXPECT().RemoteVerifierAddress().Return(addr3).Times(2)

	baseM, curM, nextM := filterPeers(from, curs, nexts)

	assert.Len(t, baseM, 1)
	assert.Len(t, curM, 1)
	assert.Len(t, nextM, 1)
}

func Test_peerInVers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCases := []struct {
		name   string
		given  func() (PmAbstractPeer, []common.Address)
		expect bool //peer in RemoteVerifier
	}{
		{
			name: "peer not in RemoteVerifier",
			given: func() (abstractPeer PmAbstractPeer, addresses []common.Address) {
				p := NewMockPmAbstractPeer(ctrl)
				p.EXPECT().RemoteVerifierAddress().Return(common.StringToAddress("test"))
				abstractPeer = p
				addresses = []common.Address{common.StringToAddress("testNN")}
				return
			},
			expect: false,
		},
		{
			name: "peer in RemoteVerifier",
			given: func() (abstractPeer PmAbstractPeer, addresses []common.Address) {
				addrStr := "test"
				p := NewMockPmAbstractPeer(ctrl)
				p.EXPECT().RemoteVerifierAddress().Return(common.StringToAddress(addrStr))
				abstractPeer = p
				addresses = []common.Address{common.StringToAddress(addrStr)}
				return
			},
			expect: true,
		},
	}

	for _, tc := range testCases {
		p, vers := tc.given()
		if !assert.Equal(t, tc.expect, peerInVers(p, vers)) {
			t.Errorf("case:%s, expect:%v", tc.name, tc.expect)
		}
	}
}

func Test_mergePeers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mPeer := NewMockPmAbstractPeer(ctrl)

	from := map[string]PmAbstractPeer{"a": mPeer}
	to := map[string]PmAbstractPeer{}

	mergePeers(to, from)

	assert.Equal(t, mPeer, to["a"])
}
