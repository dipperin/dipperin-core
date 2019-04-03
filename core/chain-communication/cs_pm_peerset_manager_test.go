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
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gopkg.in/check.v1"
	"net"
	"strconv"
	"testing"

	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/third-party/p2p"
)

// Hook up gocheck into the "go test" runner
func Test(t *testing.T) { check.TestingT(t) }

// this is the test case for peer set
// Test the complex logic of the peer set, add the node to the set, and remove the set from the node.
type PeerSetSuite struct {
	psManager *CsPmPeerSetManager
}

var _ = check.Suite(&PeerSetSuite{})

// initialize at the beginning of the suite
func (s *PeerSetSuite) SetUpSuite(c *check.C) {
	// set PbftMaxPeerCount

}

// do at the end of the suite
func (s *PeerSetSuite) TearDownSuite(c *check.C) {}

// initialization of each test case
func (s *PeerSetSuite) SetUpTest(c *check.C) {
	PbftMaxPeerCount = 2

	ps := &CsPmPeerSetManager{
		basePeers:            newPeerSet(),
		currentVerifierPeers: newPeerSet(),
		nextVerifierPeers:    newPeerSet(),
		verifierBootNode:     newPeerSet(),
	}

	s.psManager = ps
}

// do at the end of every test case
func (s *PeerSetSuite) TearDownTest(c *check.C) {
	s.psManager = nil
}

// node type = base, add peer to set
func (s *PeerSetSuite) Test_AddPeer_pmType_base(c *check.C) {
	// set pmType base
	s.psManager.pmType = base

	// set peer max number, PbftMaxPeerCount
	s.psManager.maxPeers = 4
	PbftMaxPeerCount = 1

	// add peer, base set full
	for i := 0; i < s.psManager.maxPeers; i++ {
		peer := &tPeer{id: strconv.Itoa(i)}
		c.Check(s.psManager.AddPeer(peer), check.IsNil)

		// assert: base set length should equals 1
		c.Check(s.psManager.basePeers.Len(), check.Equals, i+1)
	}

	// too many peers
	peer := &tPeer{id: "78979879879879"}
	c.Check(s.psManager.baseAddPeer(peer).Error(), check.Equals, "too many peers")
}

// node type = verifier, add peer to set
func (s *PeerSetSuite) Test_AddPeer_pmType_verifier(c *check.C) {
	// set pmType base
	s.psManager.pmType = verifier

	// set peer max number, PbftMaxPeerCount
	s.psManager.maxPeers = 8
	PbftMaxPeerCount = 4

	// add normal peer
	normalPeer := &tPeer{id: "111", nodeType: chain_config.NodeTypeOfNormal}
	c.Check(s.psManager.AddPeer(normalPeer), check.IsNil)
	c.Check(s.psManager.basePeers.Len(), check.Equals, 1)

	// add miner master peer
	minerMasterPeer := &tPeer{id: "222", nodeType: chain_config.NodeTypeOfMineMaster}
	c.Check(s.psManager.AddPeer(minerMasterPeer), check.IsNil)
	c.Check(s.psManager.basePeers.Len(), check.Equals, 2)

	// add verifier, but self isn't current verifier or next verifier
	verifierPeer := &tPeer{id: "333", nodeType: chain_config.NodeTypeOfVerifier}

	// set mock method for SelfIsCurrentVerifier & SelfIsNextVerifier
	mockSelfIsNotCurrentVerifierFunc := func() bool { return false }
	mockSelfIsNotNextVerifierFunc := func() bool { return false }

	s.psManager.selfIsCurrentVerifier = mockSelfIsNotCurrentVerifierFunc
	s.psManager.selfIsNextVerifier = mockSelfIsNotNextVerifierFunc

	c.Check(s.psManager.AddPeer(verifierPeer), check.IsNil)
	// Because the self node is not current or next, you need to add the remote node to the base.
	c.Check(s.psManager.basePeers.Len(), check.Equals, 3)

	// let self become current verifier
	mockSelfIsCurrentVerifierFunc := func() bool { return true }
	s.psManager.selfIsCurrentVerifier = mockSelfIsCurrentVerifierFunc

	// add verifier, is this the current or next round of verifier
	verifierPeer1 := &tPeer{id: "444", nodeType: chain_config.NodeTypeOfVerifier}

	// set mock method for IsCurrentVerifier & IsNextVerifier
	mockIsNotCurrentVerifierFunc := func(p PmAbstractPeer) bool { return false }
	mockIsNotNextVerifierFunc := func(p PmAbstractPeer) bool { return false }
	s.psManager.isCurrentVerifier = mockIsNotCurrentVerifierFunc
	s.psManager.isNextVerifier = mockIsNotNextVerifierFunc

	c.Check(s.psManager.AddPeer(verifierPeer1), check.IsNil)
	// self is the current round verifier, because the remote node is not current or next, so add the remote node to the base
	c.Check(s.psManager.basePeers.Len(), check.Equals, 4)

	// add verifier boot
	verifierBootPeer := &tPeer{id: "555", nodeType: chain_config.NodeTypeOfVerifierBoot}

	// set mock method for isVerifierBootNode
	mockIsVerifierBoot := func(p PmAbstractPeer) bool { return true }
	s.psManager.isVerifierBootNode = mockIsVerifierBoot

	c.Check(s.psManager.AddPeer(verifierBootPeer), check.IsNil)
	// The remote node is verifier boot, so add the remote node to the verifier boot peer set.
	c.Check(s.psManager.verifierBootNode.Len(), check.Equals, 1)

	// Self is the current round verifier, add verifier, remote peer is the current round verifier
	verifierPeer2 := &tPeer{id: "666", nodeType: chain_config.NodeTypeOfVerifier}

	// set mock method for IsCurrentVerifier
	mockIsCurrentVerifierFunc := func(p PmAbstractPeer) bool { return true }
	s.psManager.isCurrentVerifier = mockIsCurrentVerifierFunc

	c.Check(s.psManager.AddPeer(verifierPeer2), check.IsNil)
	// Self is the current round, remote is also the current round, so add remote to the current verifier peer set
	c.Check(s.psManager.currentVerifierPeers.Len(), check.Equals, 1)

	// Self is the current round verifier, add verifier, remote peer is the current round verifier, and the next round verifier
	verifierPeer3 := &tPeer{id: "777", nodeType: chain_config.NodeTypeOfVerifier}

	// set mock method for IsCurrentVerifier
	mockIsNextVerifierFunc := func(p PmAbstractPeer) bool { return true }
	s.psManager.isNextVerifier = mockIsNextVerifierFunc

	c.Check(s.psManager.AddPeer(verifierPeer3), check.IsNil)
	// Self is the current round, remote is also the current round, so add remote to the current verifier peer set
	c.Check(s.psManager.currentVerifierPeers.Len(), check.Equals, 2)
	c.Check(s.psManager.nextVerifierPeers.Len(), check.Equals, 0)

	// let self become current verifier
	mockSelfIsNextVerifierFunc := func() bool { return true }
	s.psManager.selfIsNextVerifier = mockSelfIsNextVerifierFunc

	// Now self is the current round and the next round of verifier, add verifier, remote is also the current round and the next round verifier
	verifierPeer4 := &tPeer{id: "888", nodeType: chain_config.NodeTypeOfVerifier}

	c.Check(s.psManager.AddPeer(verifierPeer4), check.IsNil)
	// Self is the current round, remote is also the current round, so add remote to the current verifier peer set
	c.Check(s.psManager.currentVerifierPeers.Len(), check.Equals, 3)
	// Self is the next round, remote is the next round, so add remote to next verifier peer set
	c.Check(s.psManager.nextVerifierPeers.Len(), check.Equals, 1)

	// At this time, the current verifier peer set and base are full.
	// too many peers: add base set
	normalPeer1 := &tPeer{id: "999", nodeType: chain_config.NodeTypeOfNormal}
	c.Check(s.psManager.AddPeer(normalPeer1).Error(), check.Equals, "too many peers")
	c.Check(s.psManager.basePeers.Len(), check.Equals, 4)

	// Now self is the current round and the next round of verifier, add verifier, remote is also the current round and the next round verifier
	verifierPeer5 := &tPeer{id: "1111", nodeType: chain_config.NodeTypeOfVerifier}

	c.Check(s.psManager.AddPeer(verifierPeer5).Error(), check.Equals, "too many peers")
	// Because it is full, it will throw an error, set length will not change
	c.Check(s.psManager.currentVerifierPeers.Len(), check.Equals, 3)
	//  Because the error is thrown in the add current verifier peer set, it won't go to next. The set length will not change.
	c.Check(s.psManager.nextVerifierPeers.Len(), check.Equals, 1)
}

// node type = verifier boot, add peer to set
func (s *PeerSetSuite) Test_AddPeer_pmType_verifier_boot(c *check.C) {
	// set pmType base
	s.psManager.pmType = boot

	// set peer max number, PbftMaxPeerCount
	s.psManager.maxPeers = 8
	PbftMaxPeerCount = 4

	// add normal peer
	normalPeer := &tPeer{id: "111", nodeType: chain_config.NodeTypeOfNormal}
	c.Check(s.psManager.AddPeer(normalPeer), check.IsNil)
	c.Check(s.psManager.basePeers.Len(), check.Equals, 1)

	// add miner master peer
	minerMasterPeer := &tPeer{id: "222", nodeType: chain_config.NodeTypeOfMineMaster}
	c.Check(s.psManager.AddPeer(minerMasterPeer), check.IsNil)
	c.Check(s.psManager.basePeers.Len(), check.Equals, 2)

	// add verifier, is this the current or next round of verifier
	verifierPeer1 := &tPeer{id: "444", nodeType: chain_config.NodeTypeOfVerifier}

	// set mock method for IsCurrentVerifier & IsNextVerifier
	mockIsNotCurrentVerifierFunc := func(p PmAbstractPeer) bool { return false }
	mockIsNotNextVerifierFunc := func(p PmAbstractPeer) bool { return false }
	s.psManager.isCurrentVerifier = mockIsNotCurrentVerifierFunc
	s.psManager.isNextVerifier = mockIsNotNextVerifierFunc

	c.Check(s.psManager.AddPeer(verifierPeer1), check.IsNil)
	// Self is the current round verifier, because the remote node is not current or next, so add the remote node to the base
	c.Check(s.psManager.basePeers.Len(), check.Equals, 3)

	// add verifier boot
	verifierBootPeer := &tPeer{id: "555", nodeType: chain_config.NodeTypeOfVerifierBoot}

	// set mock method for isVerifierBootNode
	mockIsVerifierBoot := func(p PmAbstractPeer) bool { return true }
	s.psManager.isVerifierBootNode = mockIsVerifierBoot

	c.Check(s.psManager.AddPeer(verifierBootPeer), check.IsNil)
	// The remote node is verifier boot, so add the remote node to the verifier boot peer set.
	c.Check(s.psManager.verifierBootNode.Len(), check.Equals, 1)

	// Self is the current round verifier, add verifier, remote peer is the current round verifier
	verifierPeer2 := &tPeer{id: "666", nodeType: chain_config.NodeTypeOfVerifier}

	// set mock method for IsCurrentVerifier
	mockIsCurrentVerifierFunc := func(p PmAbstractPeer) bool { return true }
	s.psManager.isCurrentVerifier = mockIsCurrentVerifierFunc

	c.Check(s.psManager.AddPeer(verifierPeer2), check.IsNil)
	//  Remote is also the current round, so add remote to current verifier peer set
	c.Check(s.psManager.currentVerifierPeers.Len(), check.Equals, 1)

	//Add verifier, remote peer is the current round verifier, and the next round verifier
	verifierPeer3 := &tPeer{id: "777", nodeType: chain_config.NodeTypeOfVerifier}

	// set mock method for IsCurrentVerifier
	mockIsNextVerifierFunc := func(p PmAbstractPeer) bool { return true }
	s.psManager.isNextVerifier = mockIsNextVerifierFunc

	c.Check(s.psManager.AddPeer(verifierPeer3), check.IsNil)
	// Self is the current round, remote is also the current round, so add remote to the current verifier peer set
	c.Check(s.psManager.currentVerifierPeers.Len(), check.Equals, 2)
	c.Check(s.psManager.nextVerifierPeers.Len(), check.Equals, 1)

	// Add verifier, remote is also the current round and the next round verifier
	verifierPeer4 := &tPeer{id: "888", nodeType: chain_config.NodeTypeOfVerifier}

	c.Check(s.psManager.AddPeer(verifierPeer4), check.IsNil)
	// Self is the current round, remote is also the current round, so add remote to the current verifier peer set
	c.Check(s.psManager.currentVerifierPeers.Len(), check.Equals, 3)
	// Self is the next round, remote is the next round, so add remote to next verifier peer set
	c.Check(s.psManager.nextVerifierPeers.Len(), check.Equals, 2)

	// Add verifier, remote is also the current round and the next round verifier
	verifierPeer5 := &tPeer{id: "1111", nodeType: chain_config.NodeTypeOfVerifier}

	c.Check(s.psManager.AddPeer(verifierPeer5), check.IsNil)
	// so add remote to current verifier peer set
	c.Check(s.psManager.currentVerifierPeers.Len(), check.Equals, 4)
	// so add remote to next verifier peer set
	c.Check(s.psManager.nextVerifierPeers.Len(), check.Equals, 3)
}

// test case: base peer set add peer
func (s *PeerSetSuite) Test_baseAddPeer(c *check.C) {
	// set peer max number
	s.psManager.maxPeers = 2

	// create mock peer
	peer := &tPeer{id: "12312313213"}

	c.Check(s.psManager.baseAddPeer(peer), check.IsNil)
	c.Check(s.psManager.basePeers.Len(), check.Equals, 1)
	c.Check(s.psManager.basePeers.Peer(peer.id), check.NotNil)
	c.Check(s.psManager.basePeers.Peer(peer.id).ID(), check.Equals, peer.id)

	peer2 := &tPeer{id: "456465465465"}
	c.Check(s.psManager.baseAddPeer(peer2), check.IsNil)
	c.Check(s.psManager.basePeers.Len(), check.Equals, 2)

	peer3 := &tPeer{id: "78979879879879"}
	c.Check(s.psManager.baseAddPeer(peer3).Error(), check.Equals, "too many peers")
}

// test case: Here we test the add node logic for the self node type verifier
func (s *PeerSetSuite) Test_verifierAddPeer(c *check.C) {
	// set peer max number
	s.psManager.maxPeers = 6

	// create mock peer node type normal
	peer := &tPeer{id: "12312313213", nodeType: chain_config.NodeTypeOfNormal}

	c.Check(s.psManager.verifierAddPeer(peer), check.IsNil)
	c.Check(s.psManager.basePeers.Len(), check.Equals, 1)

	// create mock peer node type miner master
	peer2 := &tPeer{id: "546456465456", nodeType: chain_config.NodeTypeOfMineMaster}
	c.Check(s.psManager.verifierAddPeer(peer2), check.IsNil)
	c.Check(s.psManager.basePeers.Len(), check.Equals, 2)

	// create mock peer node type verifier boot
	peer3 := &tPeer{id: "789798798789", nodeType: chain_config.NodeTypeOfVerifierBoot}

	// mock method
	mockIsVerifierBootNodeFunc := func(p PmAbstractPeer) bool {
		if p.ID() == peer3.id {
			return true
		}
		return false
	}

	s.psManager.isVerifierBootNode = mockIsVerifierBootNodeFunc
	c.Check(s.psManager.verifierAddPeer(peer3), check.IsNil)
	c.Check(s.psManager.verifierBootNode.Len(), check.Equals, 1)

	peer4 := &tPeer{id: "963963963963", nodeType: chain_config.NodeTypeOfVerifier}

	mockSelfIsCurrentVerifierFunc1 := func() bool {
		return false
	}

	mockSelfIsNextVerifierFunc1 := func() bool {
		return false
	}

	s.psManager.selfIsCurrentVerifier = mockSelfIsCurrentVerifierFunc1
	s.psManager.selfIsNextVerifier = mockSelfIsNextVerifierFunc1

	c.Check(s.psManager.verifierAddPeer(peer4), check.IsNil)
	c.Check(s.psManager.basePeers.Len(), check.Equals, 3)

}

// test case: Here we test the add node logic of the self node type for verifier boot
func (s *PeerSetSuite) Test_verifierBootAddPeer(c *check.C) {
	// set peer max number
	s.psManager.maxPeers = 6

	// create mock peer node type normal
	peer := &tPeer{id: "12312313213", nodeType: chain_config.NodeTypeOfNormal}

	c.Check(s.psManager.verifierBootAddPeer(peer), check.IsNil)
	c.Check(s.psManager.basePeers.Len(), check.Equals, 1)

	// create mock peer node type verifier boot
	peer3 := &tPeer{id: "789798798789", nodeType: chain_config.NodeTypeOfVerifierBoot}

	// mock method
	mockIsVerifierBootNodeFunc := func(p PmAbstractPeer) bool {
		if p.ID() == peer3.id {
			return true
		}
		return false
	}

	s.psManager.isVerifierBootNode = mockIsVerifierBootNodeFunc
	c.Check(s.psManager.verifierBootAddPeer(peer3), check.IsNil)
	c.Check(s.psManager.verifierBootNode.Len(), check.Equals, 1)

	// create mock peer node type verifier
	// test case: self node current verifier
	peer4 := &tPeer{id: "963963963963", nodeType: chain_config.NodeTypeOfVerifier}

	mockIsCurrentVerifierFunc1 := func(p PmAbstractPeer) bool {
		return true
	}

	mockIsNextVerifierFunc1 := func(p PmAbstractPeer) bool {
		return false
	}

	s.psManager.isCurrentVerifier = mockIsCurrentVerifierFunc1
	s.psManager.isNextVerifier = mockIsNextVerifierFunc1

	c.Check(s.psManager.verifierBootAddPeer(peer4), check.IsNil)
	c.Check(s.psManager.currentVerifierPeers.Len(), check.Equals, 1)
}

// test case: Here we test the add node logic for the self node type verifier
// i am not elected verifier
func (s *PeerSetSuite) Test_verifierAddVerifierSet_self_no_verifier(c *check.C) {
	// set peer max number
	s.psManager.maxPeers = 3

	// create mock peer node type verifier
	// test case: self node isn't current verifier or next verifier
	peer4 := &tPeer{id: "963963963963", nodeType: chain_config.NodeTypeOfVerifier}

	mockSelfIsCurrentVerifierFunc1 := func() bool {
		return false
	}

	mockSelfIsNextVerifierFunc1 := func() bool {
		return false
	}

	s.psManager.selfIsCurrentVerifier = mockSelfIsCurrentVerifierFunc1
	s.psManager.selfIsNextVerifier = mockSelfIsNextVerifierFunc1

	c.Check(s.psManager.verifierAddVerifierSet(peer4), check.IsNil)
	c.Check(s.psManager.basePeers.Len(), check.Equals, 1)
}

// test case: Here we test the added node logic with self node type verifier
// It is the current round verifier that is elected, and the remote peer is also the current round verifier.
func (s *PeerSetSuite) Test_verifierAddVerifierSet_self_current_verifier(c *check.C) {
	// set peer max number
	s.psManager.maxPeers = 3

	// create mock peer node type verifier
	// test case: self node current verifier
	peer4 := &tPeer{id: "963963963963", nodeType: chain_config.NodeTypeOfVerifier}

	mockSelfIsCurrentVerifierFunc1 := func() bool {
		return true
	}

	mockSelfIsNextVerifierFunc1 := func() bool {
		return false
	}

	mockIsCurrentVerifierFunc1 := func(p PmAbstractPeer) bool {
		return true
	}

	mockIsNextVerifierFunc1 := func(p PmAbstractPeer) bool {
		return false
	}

	s.psManager.selfIsCurrentVerifier = mockSelfIsCurrentVerifierFunc1
	s.psManager.selfIsNextVerifier = mockSelfIsNextVerifierFunc1
	s.psManager.isCurrentVerifier = mockIsCurrentVerifierFunc1
	s.psManager.isNextVerifier = mockIsNextVerifierFunc1

	c.Check(s.psManager.verifierAddVerifierSet(peer4), check.IsNil)
	c.Check(s.psManager.basePeers.Len(), check.Equals, 0)

	c.Check(s.psManager.currentVerifierPeers.Len(), check.Equals, 1)
}

// test case: Here we test the add node logic for the self node type verifier
// It is the current round verifier that is elected, and the remote peer is also the next round of verifier.
func (s *PeerSetSuite) Test_verifierAddVerifierSet_self_current_verifier_rn(c *check.C) {
	// set peer max number
	s.psManager.maxPeers = 3

	// create mock peer node type verifier
	// test case: self node current verifier
	peer4 := &tPeer{id: "963963963963", nodeType: chain_config.NodeTypeOfVerifier}

	mockSelfIsCurrentVerifierFunc1 := func() bool {
		return true
	}

	mockSelfIsNextVerifierFunc1 := func() bool {
		return false
	}

	mockIsCurrentVerifierFunc1 := func(p PmAbstractPeer) bool {
		return false
	}

	mockIsNextVerifierFunc1 := func(p PmAbstractPeer) bool {
		return true
	}

	s.psManager.selfIsCurrentVerifier = mockSelfIsCurrentVerifierFunc1
	s.psManager.selfIsNextVerifier = mockSelfIsNextVerifierFunc1
	s.psManager.isCurrentVerifier = mockIsCurrentVerifierFunc1
	s.psManager.isNextVerifier = mockIsNextVerifierFunc1

	c.Check(s.psManager.verifierAddVerifierSet(peer4), check.IsNil)
	c.Check(s.psManager.basePeers.Len(), check.Equals, 1)

	c.Check(s.psManager.currentVerifierPeers.Len(), check.Equals, 0)
}

// test case: Here we test the add node logic for the self node type verifier
// I am elected the current round verifier, remote peer is not elected verifier
func (s *PeerSetSuite) Test_verifierAddVerifierSet_self_current_verifier_rnv(c *check.C) {
	// set peer max number
	s.psManager.maxPeers = 3

	// create mock peer node type verifier
	// test case: self node current verifier
	peer4 := &tPeer{id: "963963963963", nodeType: chain_config.NodeTypeOfVerifier}

	mockSelfIsCurrentVerifierFunc1 := func() bool {
		return true
	}

	mockSelfIsNextVerifierFunc1 := func() bool {
		return false
	}

	mockIsCurrentVerifierFunc1 := func(p PmAbstractPeer) bool {
		return false
	}

	mockIsNextVerifierFunc1 := func(p PmAbstractPeer) bool {
		return false
	}

	s.psManager.selfIsCurrentVerifier = mockSelfIsCurrentVerifierFunc1
	s.psManager.selfIsNextVerifier = mockSelfIsNextVerifierFunc1
	s.psManager.isCurrentVerifier = mockIsCurrentVerifierFunc1
	s.psManager.isNextVerifier = mockIsNextVerifierFunc1

	c.Check(s.psManager.verifierAddVerifierSet(peer4), check.IsNil)
	c.Check(s.psManager.basePeers.Len(), check.Equals, 1)

	c.Check(s.psManager.currentVerifierPeers.Len(), check.Equals, 0)
}

// test case: Here we test the add node logic for the self node type verifier
// I am elected the next round of verifier, remote peer is not elected verifier
func (s *PeerSetSuite) Test_verifierAddVerifierSet_self_next_verifier_rnv(c *check.C) {
	// set peer max number
	s.psManager.maxPeers = 3

	// create mock peer node type verifier
	// test case: self node current verifier
	peer4 := &tPeer{id: "963963963963", nodeType: chain_config.NodeTypeOfVerifier}

	mockSelfIsCurrentVerifierFunc1 := func() bool {
		return false
	}

	mockSelfIsNextVerifierFunc1 := func() bool {
		return true
	}

	mockIsCurrentVerifierFunc1 := func(p PmAbstractPeer) bool {
		return false
	}

	mockIsNextVerifierFunc1 := func(p PmAbstractPeer) bool {
		return false
	}

	s.psManager.selfIsCurrentVerifier = mockSelfIsCurrentVerifierFunc1
	s.psManager.selfIsNextVerifier = mockSelfIsNextVerifierFunc1
	s.psManager.isCurrentVerifier = mockIsCurrentVerifierFunc1
	s.psManager.isNextVerifier = mockIsNextVerifierFunc1

	c.Check(s.psManager.verifierAddVerifierSet(peer4), check.IsNil)
	c.Check(s.psManager.basePeers.Len(), check.Equals, 1)

	c.Check(s.psManager.nextVerifierPeers.Len(), check.Equals, 0)
}

// test case: Here we test the add node logic for the self node type verifier
// I am elected the next round of verifier, remote peer elected next round verifier
func (s *PeerSetSuite) Test_verifierAddVerifierSet_self_next_verifier_nv(c *check.C) {
	// set peer max number
	s.psManager.maxPeers = 3

	// create mock peer node type verifier
	// test case: self node current verifier
	peer4 := &tPeer{id: "963963963963", nodeType: chain_config.NodeTypeOfVerifier}

	mockSelfIsCurrentVerifierFunc1 := func() bool {
		return false
	}

	mockSelfIsNextVerifierFunc1 := func() bool {
		return true
	}

	mockIsCurrentVerifierFunc1 := func(p PmAbstractPeer) bool {
		return false
	}

	mockIsNextVerifierFunc1 := func(p PmAbstractPeer) bool {
		return true
	}

	s.psManager.selfIsCurrentVerifier = mockSelfIsCurrentVerifierFunc1
	s.psManager.selfIsNextVerifier = mockSelfIsNextVerifierFunc1
	s.psManager.isCurrentVerifier = mockIsCurrentVerifierFunc1
	s.psManager.isNextVerifier = mockIsNextVerifierFunc1

	c.Check(s.psManager.verifierAddVerifierSet(peer4), check.IsNil)
	c.Check(s.psManager.basePeers.Len(), check.Equals, 0)

	c.Check(s.psManager.nextVerifierPeers.Len(), check.Equals, 1)
}

func (s *PeerSetSuite) Test_verifierAddBaseSet(c *check.C) {
	// set peer max number
	s.psManager.maxPeers = 2

	// set PbftMaxPeerCount
	PbftMaxPeerCount = 1

	// create mock peer
	peer := &tPeer{id: "12312313213"}
	c.Check(s.psManager.verifierAddBaseSet(peer), check.IsNil)

	peer3 := &tPeer{id: "78979879879879"}
	c.Check(s.psManager.verifierAddBaseSet(peer3).Error(), check.Equals, "too many peers")

}

// mock peer
type tPeer struct {
	nodeType uint64
	name     string
	id       string
	address  common.Address
}

func (p *tPeer) NodeName() string {
	return p.name
}

func (p *tPeer) NodeType() uint64 {
	return p.nodeType
}

func (p *tPeer) SendMsg(msgCode uint64, msg interface{}) error {
	panic("implement me")
}

func (p *tPeer) ID() string {
	return p.id
}

func (p *tPeer) ReadMsg() (p2p.Msg, error) {
	panic("implement me")
}

func (p *tPeer) GetHead() (common.Hash, uint64) {
	panic("implement me")
}

func (p *tPeer) SetHead(head common.Hash, height uint64) {
	panic("implement me")
}

func (p *tPeer) GetPeerRawUrl() string {
	panic("implement me")
}

func (p *tPeer) DisconnectPeer() {
	panic("implement me")
}

func (p *tPeer) RemoteVerifierAddress() (addr common.Address) {
	return p.address
}

func (p *tPeer) RemoteAddress() net.Addr {
	panic("implement me")
}

func (p *tPeer) SetRemoteVerifierAddress(addr common.Address) {
	panic("implement me")
}

func (p *tPeer) SetNodeName(name string) {
	panic("implement me")
}

func (p *tPeer) SetNodeType(nt uint64) {
	panic("implement me")
}

func (p *tPeer) SetPeerRawUrl(rawUrl string) {
	panic("implement me")
}

func (p *tPeer) SetNotRunning() {
	panic("implement me")
}

func (p *tPeer) IsRunning() bool {
	panic("implement me")
}

func (p *tPeer) GetCsPeerInfo() *p2p.CsPeerInfo {
	panic("implement me")
}

// node type = verifier boot, add peer to set
func (s *PeerSetSuite) Test_AddPeer_return_nil(c *check.C) {
	// set pmType base
	s.psManager.pmType = 6

	// set peer max number, PbftMaxPeerCount
	s.psManager.maxPeers = 8
	PbftMaxPeerCount = 4

	// add normal peer
	normalPeer := &tPeer{id: "111", nodeType: 6}
	c.Check(s.psManager.AddPeer(normalPeer), check.IsNil)
	c.Check(s.psManager.basePeers.Len(), check.Equals, 0)
}

func (s *PeerSetSuite) TestCsPmPeerSetManager_BestPeer(c *check.C) {
	s.psManager.pmType = base

	// set peer max number, PbftMaxPeerCount
	s.psManager.maxPeers = 8
	PbftMaxPeerCount = 4

	s.psManager.BestPeer()

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

	psm.RemovePeer("sss")
}

func TestCsPmPeerSetManager_removePeer(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// mock peer set
	mockPeerSetBasePeers := NewMockAbstractPeerSet(ctrl)

	psm := &CsPmPeerSetManager{
		basePeers: mockPeerSetBasePeers,
	}

	mockPeerSetBasePeers.EXPECT().Peer(gomock.Any()).Return(nil).AnyTimes()
	psm.removePeer("aaa")
}

func TestCsPmPeerSetManager_removePeer1(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// mock peer set
	mockPeerSetBasePeers := NewMockAbstractPeerSet(ctrl)

	psm := &CsPmPeerSetManager{
		basePeers: mockPeerSetBasePeers,
	}

	mockPeer := NewMockPmAbstractPeer(ctrl)
	mockPeer.EXPECT().NodeName().Return("test")
	mockPeer.EXPECT().DisconnectPeer().Return()
	mockPeerSetBasePeers.EXPECT().Peer(gomock.Any()).Return(mockPeer).AnyTimes()
	mockPeerSetBasePeers.EXPECT().RemovePeer(gomock.Any()).Return(errors.New("aaa")).AnyTimes()
	psm.removePeer("aaa")
}

func TestCsPmPeerSetManager_removeCurrentVerifierPeers(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// mock peer set
	mockPeerSetBasePeers := NewMockAbstractPeerSet(ctrl)

	psm := &CsPmPeerSetManager{
		currentVerifierPeers: mockPeerSetBasePeers,
	}

	mockPeer := NewMockPmAbstractPeer(ctrl)
	mockPeer.EXPECT().NodeName().Return("test")
	mockPeer.EXPECT().DisconnectPeer().Return()
	mockPeerSetBasePeers.EXPECT().Peer(gomock.Any()).Return(mockPeer).AnyTimes()
	mockPeerSetBasePeers.EXPECT().RemovePeer(gomock.Any()).Return(errors.New("aaa")).AnyTimes()
	psm.removeCurrentVerifierPeers("aaa")
}

func TestCsPmPeerSetManager_removeNextVerifierPeers(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// mock peer set
	mockPeerSetBasePeers := NewMockAbstractPeerSet(ctrl)

	psm := &CsPmPeerSetManager{
		nextVerifierPeers: mockPeerSetBasePeers,
	}

	mockPeer := NewMockPmAbstractPeer(ctrl)
	mockPeer.EXPECT().NodeName().Return("test")
	mockPeer.EXPECT().DisconnectPeer().Return()
	mockPeerSetBasePeers.EXPECT().Peer(gomock.Any()).Return(mockPeer).AnyTimes()
	mockPeerSetBasePeers.EXPECT().RemovePeer(gomock.Any()).Return(errors.New("aaa")).AnyTimes()
	psm.removeNextVerifierPeers("aaa")
}

func TestCsPmPeerSetManager_removeVerifierBootNodePeers(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// mock peer set
	mockPeerSetBasePeers := NewMockAbstractPeerSet(ctrl)

	psm := &CsPmPeerSetManager{
		verifierBootNode: mockPeerSetBasePeers,
	}

	mockPeer := NewMockPmAbstractPeer(ctrl)
	mockPeer.EXPECT().NodeName().Return("test")
	mockPeer.EXPECT().DisconnectPeer().Return()
	mockPeerSetBasePeers.EXPECT().Peer(gomock.Any()).Return(mockPeer).AnyTimes()
	mockPeerSetBasePeers.EXPECT().RemovePeer(gomock.Any()).Return(errors.New("aaa")).AnyTimes()
	psm.removeVerifierBootNodePeers("aaa")
}

func TestCsPmPeerSetManager_baseAddPeer(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// mock peer set
	mockPeerSetBasePeers := NewMockAbstractPeerSet(ctrl)

	psm := &CsPmPeerSetManager{
		basePeers: mockPeerSetBasePeers,
	}

	psm.maxPeers = 2
	mockPeerSetBasePeers.EXPECT().Len().Return(3)

	assert.NotNil(t, psm.baseAddPeer(nil))
}

func TestCsPmPeerSetManager_baseAddPeer1(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// mock peer set
	mockPeerSetBasePeers := NewMockAbstractPeerSet(ctrl)

	psm := &CsPmPeerSetManager{
		basePeers: mockPeerSetBasePeers,
	}

	psm.maxPeers = 4
	mockPeerSetBasePeers.EXPECT().Len().Return(3)
	mockPeerSetBasePeers.EXPECT().AddPeer(gomock.Any()).Return(nil)

	assert.NoError(t, psm.baseAddPeer(nil))
}

func TestCsPmPeerSetManager_baseAddPeer3(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// mock peer set
	mockPeerSetBasePeers := NewMockAbstractPeerSet(ctrl)

	psm := &CsPmPeerSetManager{
		basePeers: mockPeerSetBasePeers,
	}

	psm.maxPeers = 4
	mockPeerSetBasePeers.EXPECT().Len().Return(3)
	mockPeerSetBasePeers.EXPECT().AddPeer(gomock.Any()).Return(errors.New("aaa"))

	mockPeer := NewMockPmAbstractPeer(ctrl)
	mockPeer.EXPECT().ID().Return("aaaa")

	assert.EqualError(t, psm.baseAddPeer(mockPeer), "aaa")
}

func TestCsPmPeerSetManager_ChangeVerifiers(t *testing.T) {
	// todo
}

func TestCsPmPeerSetManager_collectAllPeers(t *testing.T) {
	// create mock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// mock peer set
	mPs1 := NewMockAbstractPeerSet(ctrl)
	mPs2 := NewMockAbstractPeerSet(ctrl)
	mPs3 := NewMockAbstractPeerSet(ctrl)

	psm := &CsPmPeerSetManager{
		basePeers: mPs1,
		currentVerifierPeers:mPs2,
		nextVerifierPeers:mPs3,
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
	mPeer := NewMockPmAbstractPeer(ctrl)

	vers := []common.Address{common.StringToAddress("aaa")}

	mPeer.EXPECT().RemoteVerifierAddress().Return(common.StringToAddress("aaa"))

	assert.Equal(t, true, peerInVers(mPeer, vers))

	vers2 := []common.Address{common.StringToAddress("aaa2")}
	mPeer.EXPECT().RemoteVerifierAddress().Return(common.StringToAddress("aaa"))
	assert.Equal(t, false, peerInVers(mPeer, vers2))

}

func Test_mergePeers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mPeer := NewMockPmAbstractPeer(ctrl)
	from := map[string]PmAbstractPeer{"a": mPeer}
	to := map[string]PmAbstractPeer{}

	mergePeers(to, from)
}
