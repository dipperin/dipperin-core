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
	"github.com/dipperin/dipperin-core/common/log"
	"github.com/dipperin/dipperin-core/core/chainconfig"
	"github.com/dipperin/dipperin-core/third_party/p2p"
	"go.uber.org/zap"
	"sync"
)

var maxPeers = 100

type p2pMsgHandler interface {
	OnNewMsg(msg p2p.Msg, p PmAbstractPeer) error
	// only for worker, master do nothing
	SetMineMasterPeer(peer PmAbstractPeer)
}

func NewMineProtocolManager(msgHandler p2pMsgHandler) *MineProtocolManager {
	return &MineProtocolManager{
		maxPeers:   maxPeers,
		peers:      newPeerSet(),
		msgHandler: msgHandler,
	}
}

type MineProtocolManager struct {
	maxPeers int
	peers    AbstractPeerSet

	msgHandler p2pMsgHandler

	wg sync.WaitGroup
}

func (pm *MineProtocolManager) handleMsg(p PmAbstractPeer) error {
	msg, err := p.ReadMsg()
	if err != nil {
		log.DLogger.Info("mine read msg from peer failed", zap.Error(err))
		return err
	}

	defer msg.Discard()
	if msg.Size > ProtocolMaxMsgSize {
		return msgTooLargeErr
	}

	// handle this msg
	if err = pm.msgHandler.OnNewMsg(msg, p); err != nil {
		p.SetNotRunning()
		return err
	}

	return nil
}

func (pm *MineProtocolManager) GetProtocol() p2p.Protocol {
	var version uint = chainconfig.MineProtocolVersion
	protocolName := chainconfig.AppName + "_mine"
	p := p2p.Protocol{Name: protocolName, Version: version, Length: 0x200}

	p.Run = func(peer *p2p.Peer, rw p2p.MsgReadWriter) error {
		log.DLogger.Info("new mine peer in", zap.String("protocol", protocolName), zap.Any("id", peer.ID()))
		// format with communication peer
		tmpPmPeer := newPeer(int(version), peer, rw)
		log.DLogger.Info("MineProtocolManager#GetProtocol", zap.Any("MineProtocolManager", pm))
		pm.wg.Add(1)
		defer pm.wg.Done()
		// read msg loop in here
		return pm.handle(tmpPmPeer)
	}
	return p
}

func (pm *MineProtocolManager) handle(p PmAbstractPeer) error {
	if pm.peers.Len() >= pm.maxPeers {
		return p2p.DiscTooManyPeers
	}

	//if err := pm.HandShake(p); err != nil {
	//	return err
	//}

	pm.msgHandler.SetMineMasterPeer(p)

	// add peer set
	if err := pm.peers.AddPeer(p); err != nil {
		log.DLogger.Error("peer set add peer failed", zap.Error(err), zap.String("p id", p.ID()))
		return err
	}

	defer pm.removePeer(p.ID())
	for {
		if err := pm.handleMsg(p); err != nil {
			log.DLogger.Info("handle mine peer msg failed, remove this peer", zap.Error(err))
			return err
		}
	}
}

func (pm *MineProtocolManager) removePeer(peerID string) {
	// Short circuit if the peer was already removed
	peer := pm.peers.Peer(peerID)
	if peer == nil {
		return
	}

	if err := pm.peers.RemovePeer(peerID); err != nil {
		log.DLogger.Error("mine peer removal failed", zap.String("peer", peerID), zap.Error(err))
	}

	// Hard disconnect at the networking layer
	peer.DisconnectPeer()
}

// receive all connects
//func (pm *MineProtocolManager) HandShake(p PmAbstractPeer) error {
//	return nil
//}
