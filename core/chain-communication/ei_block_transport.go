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
	"github.com/dipperin/dipperin-core/core/bloom"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/hashicorp/golang-lru"
)

const (
	maxKnownBlocks     = 1024
	maxQueuedBlock     = 4
	maxQueuedBlockHash = 4
	minBroadcastPeers  = 4
)

var (
	BroadcastTimeoutErr = errors.New("eiBlockTransport not broadcast more than 5 min, discard this transport")
)

type getPeerFunc func() PmAbstractPeer

// ei step 1 broadcast data
type eiBroadcastMsg struct {
	// block height
	Height uint64
	// block hash
	BlockHash common.Hash
	//tx bloom in the block
	TxBloom *iblt.Bloom
}

// ei step 2 send Estimator
type eiEstimatorReq struct {
	BlockHash common.Hash
	Estimator *iblt.HybridEstimator
}

type bloomBlockDataRLP struct {
	Header          *model.Header
	BloomRLP        []byte
	PreVerification []*model.VoteMsg
	CurVerification []*model.VoteMsg
	Interlinks      []common.Hash
}

type eiBlockTransport struct {
	// here is to determine if it is for wait verify
	wvTs bool

	peerID                 string
	peerName               string
	//knownBlocks            mapset.Set
	knownBlocks            *lru.Cache
	queuedEiBroadcastMsg   chan *eiBroadcastMsg
	queuedWvEiBroadcastMsg chan *eiBroadcastMsg

	queuedEstimatorMsg   chan *eiEstimatorReq
	queuedWvEstimatorMsg chan *eiEstimatorReq

	queuedEiBlockBloomMsg   chan *model.BloomBlockData
	queuedWvEiBlockBloomMsg chan *model.BloomBlockData
}
//
//func newEiBlockTransport(wvTs bool, id, name string) *eiBlockTransport {
//	knownBlocks, _ := lru.New(800)
//	ts := &eiBlockTransport{
//		wvTs:        wvTs,
//		peerName:    name,
//		peerID:      id,
//		//knownBlocks: mapset.NewSet(),
//		knownBlocks: knownBlocks,
//	}
//
//	// Here is based on wvTs to determine what is needed to initialize chan
//	if ts.wvTs {
//		ts.queuedWvEiBroadcastMsg = make(chan *eiBroadcastMsg, maxQueuedBlockHash)
//		ts.queuedWvEstimatorMsg = make(chan *eiEstimatorReq, maxQueuedBlockHash)
//		ts.queuedWvEiBlockBloomMsg = make(chan *model.BloomBlockData, maxQueuedBlockHash)
//	} else {
//		ts.queuedEiBroadcastMsg = make(chan *eiBroadcastMsg, maxQueuedBlockHash)
//		ts.queuedEstimatorMsg = make(chan *eiEstimatorReq, maxQueuedBlockHash)
//		ts.queuedEiBlockBloomMsg = make(chan *model.BloomBlockData, maxQueuedBlockHash)
//	}
//
//	return ts
//}
//
//func (receiver *eiBlockTransport) broadcast(getPeer getPeerFunc) error {
//	timer := time.NewTimer(5 * time.Minute)
//	defer timer.Stop()
//
//	for {
//		select {
//		case msg := <-receiver.queuedEiBroadcastMsg:
//			if err := receiver.sendEiBroadcastMsg(msg, getPeer); err != nil {
//				log.Error("send ei failed", "err", err)
//				return err
//			}
//
//		case msg := <-receiver.queuedWvEiBroadcastMsg:
//			if err := receiver.sendWvEiBroadcastMsg(msg, getPeer); err != nil {
//				log.Error("send wv ei failed", "err", err)
//				return err
//			}
//
//		case msg := <-receiver.queuedEiBlockBloomMsg:
//			if err := receiver.sendEiBlockByBloomMsg(msg, getPeer); err != nil {
//				log.Error("send ei block bloom failed", "err", err)
//				return err
//			}
//
//		case msg := <-receiver.queuedWvEiBlockBloomMsg:
//			if err := receiver.sendWvEiBlockByBloomMsg(msg, getPeer); err != nil {
//				log.Error("send wv ei block bloom failed", "err", err)
//				return err
//			}
//
//		case msg := <-receiver.queuedWvEstimatorMsg:
//			if err := receiver.sendWvEstimatorMsg(msg, getPeer); err != nil {
//				log.Error("send wv estimator failed", "err", err)
//				return err
//			}
//
//		case msg := <-receiver.queuedEstimatorMsg:
//			if err := receiver.sendEstimatorMsg(msg, getPeer); err != nil {
//				log.Error("send estimator failed", "err", err)
//				return err
//			}
//		case <- timer.C:
//			return BroadcastTimeoutErr
//		}
//
//		timer.Reset(5 * time.Minute)
//	}
//}
//
//func (receiver *eiBlockTransport) asyncSendEiBroadcastMsg(msg *eiBroadcastMsg) {
//	//log.Info("eiBlockTransport#asyncSendEiBroadcastMsg", "receiver.wvTs", receiver.wvTs)
//	if receiver.wvTs {
//		select {
//		case receiver.queuedWvEiBroadcastMsg <- msg:
//		default:
//			log.Info("Dropping wv ei propagation", "hash", msg.BlockHash, "height", msg.Height)
//		}
//	} else {
//		select {
//		case receiver.queuedEiBroadcastMsg <- msg:
//		default:
//			log.Info("Dropping ei propagation", "hash", msg.BlockHash, "height", msg.Height)
//		}
//	}
//}
//
//func (receiver *eiBlockTransport) asyncEstimatorMsg(msg *eiEstimatorReq) {
//	if receiver.wvTs {
//		select {
//		case receiver.queuedWvEstimatorMsg <- msg:
//		default:
//			log.Info("Dropping wv estimator propagation", "hash", msg.BlockHash)
//		}
//	} else {
//		select {
//		case receiver.queuedEstimatorMsg <- msg:
//		default:
//			log.Info("Dropping ei propagation", "hash", msg.BlockHash)
//		}
//	}
//}
//
//func (receiver *eiBlockTransport) asyncEiBlockByBloomMsg(msg *model.BloomBlockData) {
//	if receiver.wvTs {
//		select {
//		case receiver.queuedWvEiBlockBloomMsg <- msg:
//		default:
//			log.Info("Dropping wv blockByBloom propagation", "hash", msg.Header.Hash().Hex(), "height", msg.Header.Number)
//		}
//	} else {
//		select {
//		case receiver.queuedEiBlockBloomMsg <- msg:
//		default:
//			log.Info("Dropping ei blockByBloom propagation", "hash", msg.Header.Hash().Hex(), "height", msg.Header.Number)
//		}
//	}
//}
//
//func (receiver *eiBlockTransport) sendEiBroadcastMsg(msg *eiBroadcastMsg, getPeer getPeerFunc) error {
//	// mark hash
//	receiver.knownBlocks.Add(msg.BlockHash, 1)
//	log.Info("eiBlockTransport#sendEiBroadcastMsg   send block to peer", "node", receiver.peerName)
//
//	if peer := getPeer(); peer != nil {
//		return peer.SendMsg(EiNewBlockHashMsg, msg)
//	}
//
//	return errors.New("send ei broadcast msg no found peer name :" + receiver.peerName)
//}
//
//func (receiver *eiBlockTransport) sendWvEiBroadcastMsg(msg *eiBroadcastMsg, getPeer getPeerFunc) error {
//	// mark hash
//	receiver.knownBlocks.Add(msg.BlockHash, 1)
//	//log.Info("send block to peer", "node", receiver.peerName)
//
//	if peer := getPeer(); peer != nil {
//		return peer.SendMsg(EiWaitVerifyBlockHashMsg, msg)
//	}
//
//	return errors.New("send wv ei broadcast msg no found peer name :" + receiver.peerName)
//}
//
//func (receiver *eiBlockTransport) sendEstimatorMsg(msg *eiEstimatorReq, getPeer getPeerFunc) error {
//	if peer := getPeer(); peer != nil {
//		return peer.SendMsg(EiEstimatorMsg, msg)
//	}
//
//	return errors.New("send estimator no found peer name :" + receiver.peerName)
//}
//
//func (receiver *eiBlockTransport) sendWvEstimatorMsg(msg *eiEstimatorReq, getPeer getPeerFunc) error {
//	if peer := getPeer(); peer != nil {
//		return peer.SendMsg(EiWaitVerifyEstimatorMsg, msg)
//	}
//
//	return errors.New("send wv estimator no found peer name :" + receiver.peerName)
//}
//
//func (receiver *eiBlockTransport) sendEiBlockByBloomMsg(msg *model.BloomBlockData, getPeer getPeerFunc) error {
//	if peer := getPeer(); peer != nil {
//		return peer.SendMsg(EiNewBlockByBloomMsg, msg)
//	}
//
//	return errors.New("send ei block bloom no found peer name :" + receiver.peerName)
//}
//
//func (receiver *eiBlockTransport) sendWvEiBlockByBloomMsg(msg *model.BloomBlockData, getPeer getPeerFunc) error {
//	if peer := getPeer(); peer != nil {
//		return peer.SendMsg(EiWaitVerifyBlockByBloomMsg, msg)
//	}
//
//	return errors.New("send wv ei block bloom no found peer name :" + receiver.peerName)
//}
//
//func (receiver *eiBlockTransport) markHash(hash common.Hash) {
//	//for receiver.knownBlocks.Cardinality() >= maxKnownBlocks {
//	//	receiver.knownBlocks.Pop()
//	//}
//	receiver.knownBlocks.Add(hash, 1)
//}
