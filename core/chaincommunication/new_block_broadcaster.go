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
	"github.com/dipperin/dipperin-core/common/gmetrics"
	"github.com/dipperin/dipperin-core/common/log"
	"github.com/dipperin/dipperin-core/core/chainconfig"
	model2 "github.com/dipperin/dipperin-core/core/csbft/model"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third_party/p2p"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/hashicorp/golang-lru"
	"go.uber.org/zap"
	"math"
	"reflect"
	"sync"
)

type NewBlockBroadcasterConfig struct {
	Chain    Chain
	Pm       PeerManager
	PbftNode PbftNode
}

func makeNewBlockBroadcaster(config *NewBlockBroadcasterConfig) *NewBlockBroadcaster {
	service := &NewBlockBroadcaster{
		NewBlockBroadcasterConfig: config,
		handlers:                  make(map[uint64]func(msg p2p.Msg, p PmAbstractPeer) error),
	}
	service.handlers[NewBlockV1Msg] = service.onNewBlock

	return service
}

type NewBlockBroadcaster struct {
	*NewBlockBroadcasterConfig
	handlers map[uint64]func(msg p2p.Msg, p PmAbstractPeer) error

	// key --> peer id, value --> blockReceiver
	// wait verify block use this
	waitVerifyBroadcast sync.Map
}

func (broadcaster *NewBlockBroadcaster) MsgHandlers() map[uint64]func(msg p2p.Msg, p PmAbstractPeer) error {
	return broadcaster.handlers
}

// get receiver
func (broadcaster *NewBlockBroadcaster) getReceiver(p PmAbstractPeer) *blockReceiver {
	// load txReceiver
	var receiver *blockReceiver

	if cache, ok := broadcaster.waitVerifyBroadcast.Load(p.ID()); ok {
		receiver = cache.(*blockReceiver)
	} else {
		receiver = broadcaster.newBlockReceiver(p)
		broadcaster.waitVerifyBroadcast.Store(p.ID(), receiver)
	}

	return receiver
}

// get peer without block
func (broadcaster *NewBlockBroadcaster) getPeersWithoutBlock(block model.AbstractBlock) []PmAbstractPeer {
	// get peers
	peers := broadcaster.Pm.GetPeers()

	var list []PmAbstractPeer

	for _, p := range peers {
		receiver := broadcaster.getReceiver(p)

		if !receiver.knownBlocks.Contains(block.Hash()) {
			list = append(list, p)
		}
	}

	return list
}

// broadcast new block
func (broadcaster *NewBlockBroadcaster) BroadcastBlock(block model.AbstractBlock) {
	log.DLogger.Info("new block broadcaster BroadcastBlock", zap.Uint64("num", block.Number()))
	log.DLogger.Debug("broadcast block", zap.Uint64("num", block.Number()), zap.Int("txs", block.TxCount()))
	peers := broadcaster.getPeersWithoutBlock(block)

	var vPeers []PmAbstractPeer
	var rPeers []PmAbstractPeer

	for i := range peers {
		if peers[i].NodeType() == chainconfig.NodeTypeOfVerifier {
			vPeers = append(vPeers, peers[i])
		} else {
			rPeers = append(rPeers, peers[i])
		}
	}

	transfer := broadcaster.getTransferPeers(rPeers)

	log.DLogger.Info("Miner broad cast block to", zap.Uint64("Height", block.Number()), zap.Int("v peer len", len(vPeers)), zap.Int("other peer len", len(transfer)))

	broadcaster.broadcastBlock(block, vPeers)
	broadcaster.broadcastBlock(block, transfer)
}

func (broadcaster *NewBlockBroadcaster) getTransferPeers(peers []PmAbstractPeer) []PmAbstractPeer {
	transferLen := int(math.Sqrt(float64(len(peers))))
	if transferLen < minBroadcastPeers {
		transferLen = minBroadcastPeers
	}
	if transferLen > len(peers) {
		transferLen = len(peers)
	}
	return peers[:transferLen]
}

func (broadcaster *NewBlockBroadcaster) broadcastBlock(block model.AbstractBlock, peers []PmAbstractPeer) {
	for i := range peers {
		receiver := broadcaster.getReceiver(peers[i])
		receiver.asyncSendBlock(block)
		log.DLogger.Debug("broadcast block", zap.String("to", peers[i].NodeName()), zap.Uint64("type", peers[i].NodeType()), zap.Uint64("num", block.Number()), zap.Int("txs", block.TxCount()))
	}
}

func (broadcaster *NewBlockBroadcaster) onNewBlock(msg p2p.Msg, p PmAbstractPeer) error {
	gmetrics.Add(gmetrics.ReceivedWaitVBlockCount, "", 1)

	log.DLogger.Debug("receive new block", zap.String("from", p.NodeName()))
	var block model.Block
	err := msg.Decode(&block)
	if err != nil {
		return err
	}
	//receiptHash := block.GetReceiptHash()
	//bloomLog := block.GetBloomLog()
	//log.DLogger.Info("NewBlockBroadcaster#onNewBlock", "bloomLog", (&bloomLog).Hex(), "receipts", receiptHash, "bloomLogs2", fmt.Sprintf("%s", (&bloomLog).Hex()))

	// load blockReceiver
	broadcaster.getReceiver(p).markBlock(&block)

	pbftNode := broadcaster.PbftNode
	log.DLogger.Info("Get new block", zap.String("from", p.NodeName()), zap.Bool("Is pbft", !reflect.ValueOf(pbftNode).IsNil()))
	if !reflect.ValueOf(pbftNode).IsNil() {
		pbftNode.OnNewWaitVerifyBlock(&block, p.ID())
	}
	return nil

}

func (broadcaster *NewBlockBroadcaster) newBlockReceiver(peer PmAbstractPeer) *blockReceiver {
	log.DLogger.Info("new block receiver", zap.String("p", peer.NodeName()))

	kb, _ := lru.New(500)
	receiver := &blockReceiver{
		peerID:          peer.ID(),
		peerName:        peer.NodeName(),
		knownBlocks:     kb,
		queuedBlock:     make(chan model.AbstractBlock, maxQueuedBlock),
		queuedBlockHash: make(chan model.AbstractBlock, maxQueuedBlockHash),
	}

	go func() {
		defer broadcaster.waitVerifyBroadcast.Delete(peer.ID())

		getPeer := func() PmAbstractPeer {
			return broadcaster.Pm.GetPeer(peer.ID())
		}

		if err := receiver.broadcast(getPeer); err != nil {
			log.DLogger.Error("wait verify block broadcast error", zap.Error(err), zap.String("peer name", peer.NodeName()))
			return
		}
	}()

	return receiver
}

// block receiver
type blockReceiver struct {
	knownBlocks     *lru.Cache
	queuedBlock     chan model.AbstractBlock
	queuedBlockHash chan model.AbstractBlock

	queuedVerifyResult     chan *model2.VerifyResult
	queuedVerifyResultHash chan *model2.VerifyResult

	peerID   string
	peerName string
}

// async send block
func (r *blockReceiver) asyncSendBlock(block model.AbstractBlock) {
	select {
	case r.queuedBlock <- block:
		//r.knownBlocks.Add(block.Hash())
	default:
		log.DLogger.Info("Dropping block propagation", zap.Uint64("number", block.Number()), zap.Any("hash", block.Hash()))
	}
}

// async send block hash
func (r *blockReceiver) asyncSendBlockHash(block model.AbstractBlock) {
	select {
	case r.queuedBlockHash <- block:
		//r.knownBlocks.Add(block.Hash())
	default:
		log.DLogger.Info("Dropping block propagation", zap.Uint64("number", block.Number()), zap.Any("hash", block.Hash()))
	}
}

// async send block
func (r *blockReceiver) asyncSendVerifyResult(result *model2.VerifyResult) {
	select {
	case r.queuedVerifyResult <- result:

		//r.knownBlocks.Add(result.Block.Hash())
	default:
		log.DLogger.Info("Dropping result propagation", zap.Uint64("number", result.Block.Number()), zap.Any("hash", result.Block.Hash()))
	}
}

// async send block hash
func (r *blockReceiver) asyncSendVerifyResultHash(result *model2.VerifyResult) {
	select {
	case r.queuedVerifyResultHash <- result:
		//r.knownBlocks.Add(result.Block.Hash())
	default:
		log.DLogger.Info("Dropping result propagation", zap.Uint64("number", result.Block.Number()), zap.Int("chan len", len(r.queuedVerifyResultHash)))
	}
}

func (r *blockReceiver) broadcast(getPeer getPeerFunc) error {
	for {
		select {
		case block := <-r.queuedBlock:
			if err := r.sendBlock(block, getPeer); err != nil {
				log.DLogger.Error("send block err", zap.Error(err))
				return err
			}

		case block := <-r.queuedBlockHash:
			//log.DLogger.Info("blockReceiver send wait verify block hash", "num", block.Number())
			if err := r.sendBlockHash(block, getPeer); err != nil {
				log.DLogger.Error("send block hash err", zap.Error(err))
				return err
			}

		case result := <-r.queuedVerifyResult:
			//log.DLogger.Info("blockReceiver send v result block", "num", result.Block.Number(), "commits len", len(result.SeenCommits))
			if err := r.sendVerifyResult(result, getPeer); err != nil {
				return err
			}

		case result := <-r.queuedVerifyResultHash:
			//log.DLogger.Info("blockReceiver send v result block hash", "node name", r.peerName)
			if err := r.sendVerifyResultHash(result, getPeer); err != nil {
				return err
			}
		}
	}
}

func (r *blockReceiver) sendBlock(block model.AbstractBlock, getPeer getPeerFunc) error {

	rlpValue, _ := rlp.EncodeToBytes(block)
	size := common.StorageSize(len(rlpValue))

	log.DLogger.Debug("send block size", zap.String("storage size", size.String()), zap.Int("block tx size", block.TxCount()))

	r.markBlock(block)

	if peer := getPeer(); peer != nil {
		log.DLogger.Debug("send block", zap.Uint64("block", block.Number()), zap.String("peer", peer.NodeName()))
		return peer.SendMsg(NewBlockV1Msg, block)
	}

	return errors.New("no found peer name :" + r.peerName)
}

func (r *blockReceiver) sendBlockHash(block model.AbstractBlock, getPeer getPeerFunc) error {
	r.markBlock(block)

	msg := &blockHashMsg{
		BlockHash:   block.Hash(),
		BlockNumber: block.Number(),
	}

	if peer := getPeer(); peer != nil {
		return peer.SendMsg(BlockHashesMsg, msg)
	}

	return errors.New("no found peer name :" + r.peerName)
}

func (r *blockReceiver) sendVerifyResult(result *model2.VerifyResult, getPeer getPeerFunc) error {
	r.markVerifyResult(result.Block.Hash())
	if peer := getPeer(); peer != nil {

		return peer.SendMsg(VerifyBlockResultMsg, result)
	}

	return errors.New("no found peer name :" + r.peerName)
}

func (r *blockReceiver) sendVerifyResultHash(result *model2.VerifyResult, getPeer getPeerFunc) error {
	r.markVerifyResult(result.Block.Hash())

	msg := &blockHashMsg{
		BlockHash:   result.Block.Hash(),
		BlockNumber: result.Block.Number(),
	}

	if peer := getPeer(); peer != nil {
		return peer.SendMsg(VerifyBlockHashResultMsg, msg)
	}

	return errors.New("no found peer name :" + r.peerName)
}

// mark a block as known for the peer
func (r *blockReceiver) markBlock(block model.AbstractBlock) {
	r.markVerifyResult(block.Hash())
}

func (r *blockReceiver) markVerifyResult(hash common.Hash) {
	r.knownBlocks.Add(hash, 1)
}

type blockHashMsg struct {
	BlockHash   common.Hash
	BlockNumber uint64
}
