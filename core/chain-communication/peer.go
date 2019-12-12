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
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/dipperin/dipperin-core/third-party/p2p"
	"github.com/dipperin/dipperin-core/third-party/p2p/enode"
	"net"
	"sync"
)

var (
	errClosed            = errors.New("peer set is closed")
	errAlreadyRegistered = errors.New("peer is already registered")
	errNotRegistered     = errors.New("peer is not registered")
	emptyPeerIDErr       = errors.New("peer not valid, empty peer id")
)

func newPeer(version int, p2pPeer P2PPeer, rw p2p.MsgReadWriter) PmAbstractPeer {
	return &peer{
		p2pPeer: p2pPeer,
		id:      p2pPeer.ID().String(),
		rw:      rw,
		version: version,
	}
}

// interface are awesome
type P2PPeer interface {
	Disconnect(reason p2p.DiscReason)
	ID() enode.ID
	RemoteAddr() net.Addr
}

type peer struct {
	p2pPeer P2PPeer
	version int    // Protocol version negotiated
	id      string // p2p peer id string
	//*p2p.Peer
	rw p2p.MsgReadWriter

	head   common.Hash
	height uint64

	// The default address used by the other party for communication to determine if it is a verifier
	verifierAddress common.Address
	// remote node type
	nodeType uint64
	// remote node name
	nodeName string
	// remote raw url
	rawUrl string

	lock sync.RWMutex

	// mark if the connection is already unavailable
	notRunning bool
}

func (p *peer) GetCsPeerInfo() *p2p.CsPeerInfo {
	return &p2p.CsPeerInfo{
		ID:              p.id,
		NodeName:        p.nodeName,
		NodeType:        p.nodeType,
		RawUrl:          p.rawUrl,
		HeadHash:        p.head.Hex(),
		HeadHeight:      p.height,
		VerifierAddress: p.verifierAddress.Hex(),
	}
}

func (p *peer) SetNotRunning() {
	log.Info("set peer not running", "node name", p.nodeName)
	p.notRunning = true
}

func (p *peer) IsRunning() bool {
	return !p.notRunning
}

func (p *peer) NodeName() string {
	return p.nodeName
}

func (p *peer) SetNodeType(nt uint64) {
	p.nodeType = nt
}

func (p *peer) SetNodeName(name string) {
	p.nodeName = name
}

func (p *peer) RemoteAddress() net.Addr {
	return p.p2pPeer.RemoteAddr()
}

func (p *peer) RemoteVerifierAddress() (addr common.Address) {
	return p.verifierAddress
}

func (p *peer) SetRemoteVerifierAddress(addr common.Address) {
	p.verifierAddress = addr
}

func (p *peer) GetPeerRawUrl() string {
	return p.rawUrl
}

func (p *peer) SetPeerRawUrl(rawUrl string) {
	p.rawUrl = rawUrl
}

func (p *peer) SetHead(hash common.Hash, height uint64) {
	p.lock.Lock()
	defer p.lock.Unlock()

	copy(p.head[:], hash[:])
	p.height = height
}

func (p *peer) DisconnectPeer() {
	log.Info("call peer disconnect", "p", p.nodeName)
	p.p2pPeer.Disconnect(p2p.DiscQuitting)
	p.SetNotRunning()
}

func (p *peer) ReadMsg() (p2p.Msg, error) {
	return p.rw.ReadMsg()
}

func (p *peer) ID() string {
	return p.id
}

func (p *peer) NodeType() uint64 {
	return p.nodeType
}

func (p *peer) SendMsg(msgCode uint64, msg interface{}) error {
	return p2p.Send(p.rw, uint64(msgCode), msg)
}

func (p *peer) GetHead() (common.Hash, uint64) {
	p.lock.Lock()
	defer p.lock.Unlock()

	return p.head, p.height
}

type peerSet struct {
	peers  map[string]PmAbstractPeer
	lock   sync.RWMutex
	closed bool
}

func (ps *peerSet) BestPeer() PmAbstractPeer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	var (
		bestPeer   PmAbstractPeer
		bestHeight uint64
	)

	log.Debug("get BestPeer the peers number is:", "number", len(ps.peers))
	for _, p := range ps.peers {
		_, height := p.GetHead()
		log.Debug("get best peer", "nodeName", p.NodeName(), "p height", height)
		if bestPeer == nil || height > bestHeight {
			bestPeer, bestHeight = p, height
		}
	}

	return bestPeer
}

func (ps *peerSet) GetPeers() map[string]PmAbstractPeer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	peers := make(map[string]PmAbstractPeer)
	for key, value := range ps.peers {
		peers[key] = value
	}

	return peers
}

// newPeerSet creates a new peer set to track the active participants.
func newPeerSet() *peerSet {
	return &peerSet{
		peers: make(map[string]PmAbstractPeer),
	}
}

// AddPeerSet injects a new peer into the working set, or returns an error if the
// peer is already known. If a new peer it registered, its broadcast loop is also
// started.
func (ps *peerSet) AddPeer(p PmAbstractPeer) error {
	if p.ID() == "" {
		return emptyPeerIDErr
	}

	ps.lock.Lock()
	defer ps.lock.Unlock()

	if ps.closed {
		return errClosed
	}

	if _, ok := ps.peers[p.ID()]; ok {
		log.Warn("duplicate peer replace old", "p name", p.NodeName())
		//You must return error to ensure that the handle exits. If it is a replacement, may be have a problem: AddPeer grabs the lock first than RemovePeer, after Add replaces the old unlock, Remove continues to remove the newly added peer.
		return errors.New("duplicate peer error")
	}
	ps.peers[p.ID()] = p

	return nil
}

// Unregister removes a remote peer from the active set, disabling any further
// actions to/from that particular entity.
func (ps *peerSet) RemovePeer(id string) error {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	if _, ok := ps.peers[id]; !ok {
		return errNotRegistered
	}

	delete(ps.peers, id)
	return nil
}

func (ps *peerSet) ReplacePeers(newPeers map[string]PmAbstractPeer) {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	ps.peers = newPeers
}

// Peer retrieves the registered peer with the given id.
func (ps *peerSet) Peer(id string) PmAbstractPeer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	return ps.peers[id]
}

// Len returns if the current number of GetPeers in the set.
func (ps *peerSet) Len() int {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	return len(ps.peers)
}

// Close disconnects all GetPeers.
// No new GetPeers can be registered after Close has returned.
func (ps *peerSet) Close() {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	for _, p := range ps.peers {
		p.DisconnectPeer()
	}
}

func (ps *peerSet) GetPeersInfo() []*p2p.CsPeerInfo {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	var info []*p2p.CsPeerInfo
	for _, p := range ps.peers {
		info = append(info, p.GetCsPeerInfo())
	}

	return info
}
