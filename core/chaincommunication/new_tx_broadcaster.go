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
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/log"
	"github.com/dipperin/dipperin-core/core/chainconfig"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third_party/p2p"
	"github.com/hashicorp/golang-lru"
	"go.uber.org/zap"
	"math"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

const (
	maxQueuedTxs = 128
	maxKnownTxs  = 32768

	// fixme: need test
	txSyncPackSize = 500 * 1024
)

func makeNewTxBroadcaster(config *NewTxBroadcasterConfig) *NewTxBroadcaster {
	service := &NewTxBroadcaster{
		NewTxBroadcasterConfig: config,
		handlers:               make(map[uint64]func(msg p2p.Msg, p PmAbstractPeer) error),
		txSyncC:                make(chan *txSync),
	}

	service.handlers[TxV1Msg] = service.onNewTx

	// start tx sync loop
	go service.txSyncLoop()

	return service
}

type NewTxBroadcasterConfig struct {
	P2PMsgDecoder P2PMsgDecoder
	TxPool        TxPool
	NodeConf      NodeConf
	Pm            PeerManager
}

type NewTxBroadcaster struct {
	*NewTxBroadcasterConfig

	handlers map[uint64]func(msg p2p.Msg, p PmAbstractPeer) error
	// key --> peer id, value --> txReceiver
	broadcast sync.Map

	// tx sync channel
	txSyncC chan *txSync
}

func (broadcaster *NewTxBroadcaster) MsgHandlers() map[uint64]func(msg p2p.Msg, p PmAbstractPeer) error {
	return broadcaster.handlers
}

func (broadcaster *NewTxBroadcaster) BroadcastTx(txs []model.AbstractTransaction) {
	// init txSet, key is peer id, value is tx array
	txSet := make(map[string][]model.AbstractTransaction)

	for i := range txs {
		peers := broadcaster.getPeersWithoutTx(txs[i].CalTxId())
		for _, peer := range peers {
			txSet[peer.ID()] = append(txSet[peer.ID()], txs[i])
		}
	}

	for pID, txs := range txSet {
		if cache, ok := broadcaster.broadcast.Load(pID); ok {
			receiver := cache.(*txReceiver)
			receiver.asyncSendTxs(txs)
		}
	}
}

func (broadcaster *NewTxBroadcaster) onNewTx(msg p2p.Msg, p PmAbstractPeer) error {
	// decode msg
	txs, err := broadcaster.P2PMsgDecoder.DecodeTxsMsg(msg)

	if err != nil {
		log.DLogger.Error("decode new tx msg failed", zap.Error(err))
		return err
	}

	// load txReceiver
	targetReceiver := broadcaster.getReceiver(p)

	// handle get txs
	for i := range txs {
		if txs[i] == nil {
			return errors.New("transaction is nil, tx index: " + strconv.Itoa(i))
		}

		// only for debug
		//txSender, err := txs[i].Sender(nil)
		//if err != nil {
		//	log.DLogger.Warn("get tx sender failed", zap.Error(err))
		//}
		//log.DLogger.Info("receive tx", "sender", txSender.Hex(), "tx id", txs[i].CalTxId())

		targetReceiver.markTx(txs[i].CalTxId())
	}

	// add to tx pool
	txPool := broadcaster.TxPool

	// todo for debug
	//startAt := time.Now()
	errs := txPool.AddRemotes(txs)
	//log.DLogger.Info("add tx pool AddRemotes use time", "t", time.Now().Sub(startAt), "node", p.NodeName())

	for i := range errs {
		if errs[i] != nil {
			log.DLogger.Debug("tx pool AddRemotes error", zap.Int("index", i), zap.Error(errs[i]))
			//You cannot return err here, otherwise the peer will be disconnected.
			return nil
		}
	}

	// broadcast txs to miner master
	go broadcaster.send2MinerMaster(txs)
	return nil
}

// get txs, broadcast txs to miner master
func (broadcaster *NewTxBroadcaster) send2MinerMaster(txs []model.AbstractTransaction) {
	// ensure node is the miner master, if it is not broadcast
	if broadcaster.NodeConf.GetNodeType() == chainconfig.NodeTypeOfMineMaster {
		return
	}

	// get miner master peer
	// init txSet, key is peer id, value is tx array
	txSet := make(map[string][]model.AbstractTransaction)

	for i := range txs {
		originPeers := broadcaster.getPeersWithoutTx(txs[i].CalTxId())

		// must broadcast to mine master
		for _, peer := range originPeers {
			if peer.NodeType() == chainconfig.NodeTypeOfMineMaster {
				txSet[peer.ID()] = append(txSet[peer.ID()], txs[i])
			}
		}

		// broadcast to near peers
		peers := originPeers[:int(math.Sqrt(float64(len(originPeers))))]
		for _, peer := range peers {
			txSet[peer.ID()] = append(txSet[peer.ID()], txs[i])
		}
	}

	// do broadcast
	for pID, txs := range txSet {
		if cache, ok := broadcaster.broadcast.Load(pID); ok {
			receiver := cache.(*txReceiver)
			receiver.asyncSendTxs(txs)
		}
	}
}

// get peer without tx
func (broadcaster *NewTxBroadcaster) getPeersWithoutTx(txHash common.Hash) []PmAbstractPeer {
	// get peers
	peers := broadcaster.Pm.GetPeers()

	var list []PmAbstractPeer

	for _, p := range peers {
		receiver := broadcaster.getReceiver(p)

		if !receiver.knownTxs.Contains(txHash) {
			list = append(list, p)
		}
	}

	return list
}

//TODO
func (broadcaster *NewTxBroadcaster) getReceiver(p PmAbstractPeer) *txReceiver {
	// load txReceiver
	var receiver *txReceiver

	if cache, ok := broadcaster.broadcast.Load(p.ID()); ok {
		receiver = cache.(*txReceiver)
	} else {
		receiver = broadcaster.newTxReceiver(p)
		broadcaster.broadcast.Store(p.ID(), receiver)
	}

	return receiver
}

// hack
func (broadcaster *NewTxBroadcaster) syncTxs(p PmAbstractPeer) {
	var txs []model.AbstractTransaction

	pending, _ := broadcaster.TxPool.Pending()

	for _, batch := range pending {
		txs = append(txs, batch...)
	}
	if len(txs) == 0 {
		return
	}

	select {
	case broadcaster.txSyncC <- &txSync{p, txs}:
	}
}

// loop, start on new NewTxBroadcaster
func (broadcaster *NewTxBroadcaster) txSyncLoop() {
	var (
		pending = make(map[string]*txSync)
		// ensure only one tx sync sending
		sending = false
		pack    = new(txSync)
		// result of the send
		done = make(chan error, 1)
	)

	// send func
	send := func(s *txSync) {
		size := common.StorageSize(0)
		pack.p = s.p
		pack.txs = pack.txs[:0]
		for i := 0; i < len(s.txs) && size < txSyncPackSize; i++ {
			pack.txs = append(pack.txs, s.txs[i])
			size += s.txs[i].Size()
		}

		// Remove the transactions that will be sent.
		s.txs = s.txs[:copy(s.txs, s.txs[len(pack.txs):])]
		if len(s.txs) == 0 {
			delete(pending, s.p.ID())
		}
		log.DLogger.Info("Sending batch of transactions", zap.Int("count", len(pack.txs)), zap.Any("bytes", size))
		sending = true

		getPeer := func() PmAbstractPeer {
			return broadcaster.Pm.GetPeer(s.p.ID())
		}

		go func() { done <- broadcaster.getReceiver(pack.p).sendTxs(pack.txs, getPeer) }()
	}

	// pick chooses the next pending sync
	pick := func() *txSync {
		fmt.Println(pending)
		if len(pending) == 0 {
			return nil
		}
		n := rand.Intn(len(pending)) + 1
		for _, s := range pending {
			if n--; n == 0 {
				return s
			}
		}
		return nil
	}

	for {
		select {
		case s := <-broadcaster.txSyncC:
			pending[s.p.ID()] = s
			fmt.Println(pending)
			if !sending {
				send(s)
			}
		case err := <-done:
			sending = false
			// Stop tracking peers that cause send failures.
			if err != nil {
				log.DLogger.Info("Transaction send failed", zap.Error(err))
				delete(pending, pack.p.ID())
			}
			// Schedule the next send.
			if s := pick(); s != nil {
				send(s)
			}
		}
	}

}

// tx sync
type txSync struct {
	p   PmAbstractPeer
	txs []model.AbstractTransaction
}

type txReceiver struct {
	//knownTxs  mapset.Set
	knownTxs  *lru.Cache
	queuedTxs chan []model.AbstractTransaction
	peerID    string
	peerName  string
}

func (broadcaster *NewTxBroadcaster) newTxReceiver(peer PmAbstractPeer) *txReceiver {
	knownTxs, _ := lru.New(1500)
	receiver := &txReceiver{
		peerID:   peer.ID(),
		peerName: peer.NodeName(),
		knownTxs: knownTxs,
		//knownTxs:  mapset.NewSet(),
		queuedTxs: make(chan []model.AbstractTransaction, maxQueuedTxs),
	}

	// start broadcast, send msg to txReceiver
	go func() {
		defer broadcaster.broadcast.Delete(peer.ID())

		getPeer := func() PmAbstractPeer {
			return broadcaster.Pm.GetPeer(peer.ID())
		}

		if err := receiver.broadcast(getPeer); err != nil {
			log.DLogger.Error("tx broadcast error", zap.Error(err), zap.String("peer name", peer.NodeName()))
			return
		}
	}()

	return receiver
}

// async send txs
func (r *txReceiver) asyncSendTxs(txs []model.AbstractTransaction) {
	select {
	case r.queuedTxs <- txs:
		//for _, tx := range txs {
		//	r.knownTxs.Add(tx.CalTxId())
		//}
		log.DLogger.Info("asyncSendTxs finished", zap.String("p", r.peerName))
	default:
		log.DLogger.Info("Dropping transaction propagation", zap.Int("count", len(txs)))
	}
}

// broadcast txs
func (r *txReceiver) broadcast(getPeer getPeerFunc) error {
	timer := time.NewTimer(5 * time.Minute)
	defer timer.Stop()

	for {
		select {
		case txs := <-r.queuedTxs:
			if err := r.sendTxs(txs, getPeer); err != nil {
				log.DLogger.Error("send txs err", zap.String("peer id", r.peerName), zap.Error(err))
				return err
			}
		case <-timer.C:
			return errors.New("txReceiver not broadcast more than 5 min, discard this Receiver")
		}

		timer.Reset(5 * time.Minute)
	}
}

// use p2p send msg--->txs
func (r *txReceiver) sendTxs(txs []model.AbstractTransaction, getPeer getPeerFunc) error {
	for _, tx := range txs {
		//r.knownTxs.Add(tx.CalTxId())
		r.knownTxs.Add(tx.CalTxId(), 1)
	}

	if peer := getPeer(); peer != nil {

		return peer.SendMsg(TxV1Msg, txs)
	}

	return errors.New("no found peer id " + r.peerName)
}

// mark a tx as known for the peer
func (r *txReceiver) markTx(txHash common.Hash) {
	//for r.knownTxs.Cardinality() >= maxKnownTxs {
	//	r.knownTxs.Pop()
	//}
	r.knownTxs.Add(txHash, 1)
}
