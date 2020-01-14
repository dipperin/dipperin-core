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
	"github.com/dipperin/dipperin-core/common/log"
	"github.com/dipperin/dipperin-core/core/chainconfig"
	"github.com/dipperin/dipperin-core/third_party/p2p"
	"go.uber.org/zap"
	"reflect"
	"sort"
	"sync"
)

func newCsPmPeerSetManager(pmType, maxPeers int, selfIsNextVerifier selfIsNextVerifier, selfIsCurrentVerifier selfIsCurrentVerifier, isCurrentVerifier isCurrentVerifier, isNextVerifier isNextVerifier, isVerifierBootNode isVerifierBootNode) *CsPmPeerSetManager {
	ps := &CsPmPeerSetManager{
		pmType:                pmType,
		maxPeers:              maxPeers,
		selfIsNextVerifier:    selfIsNextVerifier,
		selfIsCurrentVerifier: selfIsCurrentVerifier,
		isNextVerifier:        isNextVerifier,
		isCurrentVerifier:     isCurrentVerifier,
		isVerifierBootNode:    isVerifierBootNode,

		basePeers:            newPeerSet(),
		currentVerifierPeers: newPeerSet(),
		nextVerifierPeers:    newPeerSet(),
		verifierBootNode:     newPeerSet(),
	}

	return ps
}

// Determine if self is next verifier
type selfIsNextVerifier func() bool

// Determine if self is current verifier
type selfIsCurrentVerifier func() bool

// Determine if remote peer is current verifier
type isCurrentVerifier func(p PmAbstractPeer) bool

// Determine if remote peer is next verifier
type isNextVerifier func(p PmAbstractPeer) bool

// Determine if remote peer is verifier boot
type isVerifierBootNode func(p PmAbstractPeer) bool

type CsPmPeerSetManager struct {
	// pm type
	pmType int
	// Maximum number of peer connections
	maxPeers int
	// base peer set
	basePeers AbstractPeerSet
	// current verifier peer set
	currentVerifierPeers AbstractPeerSet
	// next round verifier peer set
	nextVerifierPeers AbstractPeerSet
	// verifier boot nodes peer set
	verifierBootNode AbstractPeerSet

	selfIsNextVerifier    selfIsNextVerifier
	selfIsCurrentVerifier selfIsCurrentVerifier

	isCurrentVerifier  isCurrentVerifier
	isNextVerifier     isNextVerifier
	isVerifierBootNode isVerifierBootNode

	changeVerifiersLock sync.Mutex
}

func (ps *CsPmPeerSetManager) AddPeer(p PmAbstractPeer) error {
	ps.changeVerifiersLock.Lock()
	defer ps.changeVerifiersLock.Unlock()

	switch ps.pmType {
	case base:
		return ps.baseAddPeer(p)
	case verifier:
		return ps.verifierAddPeer(p)
	case boot:
		return ps.verifierBootAddPeer(p)
	}

	return nil
}

func (ps *CsPmPeerSetManager) BestPeer() PmAbstractPeer {
	norP := ps.basePeers.BestPeer()
	curP := ps.currentVerifierPeers.BestPeer()
	nextP := ps.nextVerifierPeers.BestPeer()
	vbP := ps.verifierBootNode.BestPeer()
	psArr := []PmAbstractPeer{norP, curP, nextP, vbP}

	sort.Slice(psArr, func(i, j int) bool {
		iv := reflect.ValueOf(psArr[i])
		jv := reflect.ValueOf(psArr[j])
		if !iv.IsValid() || iv.IsNil() {
			return false
		} else if !jv.IsValid() || jv.IsNil() {
			return true
		}
		_, iNum := psArr[i].GetHead()
		_, jNum := psArr[j].GetHead()
		if iNum > jNum {
			return true
		}
		return false
	})

	rv0 := reflect.ValueOf(psArr[0])

	if !rv0.IsValid() || rv0.IsNil() {
		return nil
	}

	return psArr[0]
}

// remove a peer from all set
func (ps *CsPmPeerSetManager) RemovePeer(pid string) {
	ps.changeVerifiersLock.Lock()
	defer ps.changeVerifiersLock.Unlock()

	ps.removePeer(pid)
	ps.removeCurrentVerifierPeers(pid)
	ps.removeNextVerifierPeers(pid)
	ps.removeVerifierBootNodePeers(pid)
}

func (ps *CsPmPeerSetManager) removePeer(peerID string) {
	// Short circuit if the peer was already removed
	peer := ps.basePeers.Peer(peerID)
	if peer == nil {
		return
	}

	log.DLogger.Info("Removing Dipperin peer", zap.String("peer", peerID), zap.String("peerName", peer.NodeName()))

	if err := ps.basePeers.RemovePeer(peerID); err != nil {
		log.DLogger.Error("Peer removal failed", zap.String("peer", peerID), zap.Error(err))
	}

	// Hard disconnect at the networking layer
	peer.DisconnectPeer()
}

func (ps *CsPmPeerSetManager) removeCurrentVerifierPeers(peerID string) {
	// Short circuit if the peer was already removed
	peer := ps.currentVerifierPeers.Peer(peerID)
	if peer == nil {
		return
	}

	log.DLogger.Info("Removing Dipperin nextVerifierPeers", zap.String("peer", peerID), zap.String("peerName", peer.NodeName()))

	if err := ps.currentVerifierPeers.RemovePeer(peerID); err != nil {
		log.DLogger.Error("Peer removal failed", zap.String("peer", peerID), zap.Error(err))
	}

	// Hard disconnect at the networking layer
	peer.DisconnectPeer()
}

func (ps *CsPmPeerSetManager) removeNextVerifierPeers(peerID string) {
	// Short circuit if the peer was already removed
	peer := ps.nextVerifierPeers.Peer(peerID)
	if peer == nil {
		return
	}

	log.DLogger.Info("Removing Dipperin nextVerifierPeers", zap.String("peer", peerID), zap.String("peerName", peer.NodeName()))

	if err := ps.nextVerifierPeers.RemovePeer(peerID); err != nil {
		log.DLogger.Error("Peer removal failed", zap.String("peer", peerID), zap.Error(err))
	}

	// Hard disconnect at the networking layer
	peer.DisconnectPeer()
}

func (ps *CsPmPeerSetManager) removeVerifierBootNodePeers(peerID string) {
	// Short circuit if the peer was already removed
	peer := ps.verifierBootNode.Peer(peerID)
	if peer == nil {
		return
	}

	log.DLogger.Info("Removing Dipperin verifierBootNode", zap.String("peer", peerID), zap.String("peerName", peer.NodeName()))

	if err := ps.verifierBootNode.RemovePeer(peerID); err != nil {
		log.DLogger.Error("Peer removal failed", zap.String("peer", peerID), zap.Error(err))
	}

	// Hard disconnect at the networking layer
	peer.DisconnectPeer()
}

// base add peer
func (ps *CsPmPeerSetManager) baseAddPeer(p PmAbstractPeer) error {
	// check peer connection count
	if ps.basePeers.Len() >= ps.maxPeers {
		return p2p.DiscTooManyPeers
	}

	if err := ps.basePeers.AddPeer(p); err != nil {
		log.DLogger.Error("peer set add peer failed", zap.Error(err), zap.String("p id", p.ID()))
		return err
	}

	return nil
}

// verifier add peer
func (ps *CsPmPeerSetManager) verifierAddPeer(p PmAbstractPeer) error {
	// check remote peer node type
	switch p.NodeType() {
	case chainconfig.NodeTypeOfNormal, chainconfig.NodeTypeOfMineMaster:
		// If the remote node type is normal or mine master,
		// we need to put this node in the base peer set
		return ps.verifierAddBaseSet(p)

	case chainconfig.NodeTypeOfVerifier:
		// If the remote node type is verifier
		return ps.verifierAddVerifierSet(p)

	case chainconfig.NodeTypeOfVerifierBoot:
		// the remote peer is verifier boot node ?
		if ps.isVerifierBootNode(p) {
			log.DLogger.Info("verifier add verifier boot node", zap.String("nodeName", p.NodeName()))
			return ps.verifierBootNode.AddPeer(p)
		}

	default:
		return errors.New("remote peer node type illegal")

	}

	return nil
}

// verifier boot add peer
func (ps *CsPmPeerSetManager) verifierBootAddPeer(p PmAbstractPeer) error {
	switch p.NodeType() {
	case chainconfig.NodeTypeOfNormal, chainconfig.NodeTypeOfMineMaster:
		// If the remote node type is normal or mine master,
		// we need to put this node in the base peer set
		return ps.verifierAddBaseSet(p)

	case chainconfig.NodeTypeOfVerifier:
		var needAddBase = true

		if ps.isCurrentVerifier(p) {
			// check peer connection count
			if ps.currentVerifierPeers.Len() >= PbftMaxPeerCount {
				return p2p.DiscTooManyPeers
			}

			log.DLogger.Info("verifier boot node add current verifier peer", zap.String("peerName", p.NodeName()))
			if err := ps.currentVerifierPeers.AddPeer(p); err != nil {
				return err
			}

			needAddBase = false
		}

		if ps.isNextVerifier(p) {
			// check peer connection count
			if ps.nextVerifierPeers.Len() >= PbftMaxPeerCount {
				return p2p.DiscTooManyPeers
			}

			log.DLogger.Info("verifier boot node add next verifier peer", zap.String("peerName", p.NodeName()))
			if err := ps.nextVerifierPeers.AddPeer(p); err != nil {
				return err
			}

			needAddBase = false
		}

		if needAddBase {
			return ps.verifierAddBaseSet(p)
		}

		return nil

	case chainconfig.NodeTypeOfVerifierBoot:
		// the remote peer is verifier boot node ?
		if ps.isVerifierBootNode(p) {
			log.DLogger.Info("verifier boot node add verifier boot node peer", zap.String("peerName", p.NodeName()))
			return ps.verifierBootNode.AddPeer(p)
		}

	default:
		return errors.New("remote peer node type illegal")

	}

	return nil
}

// self node type is verifier, remote peer node type is base, add peer to base peer set
func (ps *CsPmPeerSetManager) verifierAddBaseSet(p PmAbstractPeer) error {
	// check peer connection count
	// If the self node type is verifier,
	// the maximum len of base peer set for the node is maxPeers - PbftMaxPeerCount
	if ps.basePeers.Len() >= (ps.maxPeers - PbftMaxPeerCount) {
		return p2p.DiscTooManyPeers
	}

	// add remote peer to base peer set
	return ps.basePeers.AddPeer(p)
}

// self node type is verifier, remote peer node type is verifier, add peer
func (ps *CsPmPeerSetManager) verifierAddVerifierSet(p PmAbstractPeer) error {
	// if self node isn't current verifier or next verifier, add remote peer to base set
	if !(ps.selfIsCurrentVerifier() || ps.selfIsNextVerifier()) {
		return ps.verifierAddBaseSet(p)
	}

	var needAddBase = true

	// Determine if the remote node is current verifier
	if ps.isCurrentVerifier(p) && ps.selfIsCurrentVerifier() {
		// check peer connection count
		if ps.currentVerifierPeers.Len() >= PbftMaxPeerCount-1 {
			return p2p.DiscTooManyPeers
		}

		log.DLogger.Info("verifier add current verifier peer", zap.String("peerName", p.NodeName()))
		if err := ps.currentVerifierPeers.AddPeer(p); err != nil {
			return err
		}

		needAddBase = false
	}

	// Determine if the remote node is next verifier
	if ps.isNextVerifier(p) && ps.selfIsNextVerifier() {
		// check peer connection count
		if ps.nextVerifierPeers.Len() >= PbftMaxPeerCount-1 {
			return p2p.DiscTooManyPeers
		}

		log.DLogger.Info("verifier add next verifier peer", zap.String("peerName", p.NodeName()))
		if err := ps.nextVerifierPeers.AddPeer(p); err != nil {
			return err
		}

		needAddBase = false
	}

	if needAddBase {
		// add peer to base set
		return ps.verifierAddBaseSet(p)
	}

	return nil
}

// You can't have a task peer in change verifiers, and you can't remove it.
// Take out all the peers and re-classify them into each set so that you don't have to move them.
func (ps *CsPmPeerSetManager) OrganizeVerifiersSet(curs []common.Address, nexts []common.Address) {
	ps.changeVerifiersLock.Lock()
	defer ps.changeVerifiersLock.Unlock()
	log.DLogger.Info("do OrganizeVerifiersSet", zap.Any("curs", curs), zap.Any("nexts", nexts))

	// collect all peers
	peers := ps.collectAllPeers()
	// classify
	newBase, newCur, newNext := filterPeers(peers, curs, nexts)
	ps.basePeers.ReplacePeers(newBase)
	ps.currentVerifierPeers.ReplacePeers(newCur)
	ps.nextVerifierPeers.ReplacePeers(newNext)

	// todo check and reduce peer count

	// logs
	log.DLogger.Info("after change verifiers", zap.Int("total len", len(peers)), zap.Int("new base", len(newBase)), zap.Int("new cur", len(newCur)), zap.Int("new next", len(newNext)))
	// only for debug
	//for _, c := range newCur {
	//	log.DLogger.Info("cur verifier", "name", c.NodeName(), "addr", c.RemoteVerifierAddress())
	//}
	//for _, c := range newNext {
	//	log.DLogger.Info("next verifier", "name", c.NodeName(), "addr", c.RemoteVerifierAddress())
	//}
}

func (ps *CsPmPeerSetManager) collectAllPeers() map[string]PmAbstractPeer {
	peers := make(map[string]PmAbstractPeer)
	mergePeers(peers, ps.basePeers.GetPeers())
	mergePeers(peers, ps.currentVerifierPeers.GetPeers())
	mergePeers(peers, ps.nextVerifierPeers.GetPeers())
	return peers
}

func filterPeers(from map[string]PmAbstractPeer, curs []common.Address, nexts []common.Address) (base map[string]PmAbstractPeer, cur map[string]PmAbstractPeer, next map[string]PmAbstractPeer) {
	base = make(map[string]PmAbstractPeer)
	cur = make(map[string]PmAbstractPeer)
	next = make(map[string]PmAbstractPeer)

	for id, p := range from {
		isCurrent := peerInVers(p, curs)
		if isCurrent {
			cur[id] = p
		}
		isNext := peerInVers(p, nexts)
		if isNext {
			next[id] = p
		}
		if !isCurrent && !isNext {
			if len(base) < NormalMaxPeerCount {
				base[id] = p
			} else {
				log.DLogger.Info("too many base peers, do disconnect", zap.String("p name", p.NodeName()))
				p.DisconnectPeer()
			}
		}
	}
	return
}

func peerInVers(p PmAbstractPeer, vers []common.Address) bool {
	for _, v := range vers {
		if p.RemoteVerifierAddress().IsEqual(v) {
			return true
		}
	}
	return false
}

func mergePeers(to map[string]PmAbstractPeer, from map[string]PmAbstractPeer) {
	for k, v := range from {
		to[k] = v
	}
}
