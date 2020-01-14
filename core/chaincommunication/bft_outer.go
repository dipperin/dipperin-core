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
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/g-metrics"
	"github.com/dipperin/dipperin-core/common/log"
	"github.com/dipperin/dipperin-core/common/util"
	model2 "github.com/dipperin/dipperin-core/core/csbft/model"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third-party/p2p"
	"github.com/hashicorp/golang-lru"
	"go.uber.org/zap"
	"sync"
	"time"
)

type NewBftOuterConfig struct {
	Chain Chain
	Pm    PeerManager
}

func NewBftOuter(config *NewBftOuterConfig) *BftOuter {
	b := &BftOuter{
		NewBftOuterConfig: config,
		handlers:          make(map[uint64]func(msg p2p.Msg, p PmAbstractPeer) error),
	}
	b.handlers[GetVerifyResultMsg] = b.onGetVerifiedResult
	b.handlers[VerifyBlockHashResultMsg] = b.onVerifiedResultBlockHash
	b.handlers[VerifyBlockResultMsg] = b.onVerifiedResultBlock

	return b
}

// csbft outer msg(broadcast verify result)
type BftOuter struct {
	*NewBftOuterConfig
	blockFetcher *BlockFetcher

	handlers map[uint64]func(msg p2p.Msg, p PmAbstractPeer) error

	// key --> peer id, value --> blockReceiver
	vResultBroadcast sync.Map
}

func (broadcaster *BftOuter) SetBlockFetcher(blockFetcher *BlockFetcher) {
	broadcaster.blockFetcher = blockFetcher
}

func (broadcaster *BftOuter) MsgHandlers() map[uint64]func(msg p2p.Msg, p PmAbstractPeer) error {
	return broadcaster.handlers
}

// Broadcast this verification result so everyone can receive the result
func (broadcaster *BftOuter) BroadcastVerifiedBlock(vr *model2.VerifyResult) {
	peers := broadcaster.getPeersWithoutBlock(vr.Block.Hash())
	log.DLogger.Info("broadcast verified block", zap.Int("p len", len(peers)), zap.String("hash", vr.Block.Hash().Hex()))

	for i := range peers {
		receiver := broadcaster.getReceiver(peers[i])
		// send block hash
		receiver.asyncSendVerifyResultHash(vr)
	}
}

// get verify result block hash
func (broadcaster *BftOuter) onVerifiedResultBlockHash(msg p2p.Msg, p PmAbstractPeer) error {
	g_metrics.Add(g_metrics.ReceivedHashCount, "", 1)

	var resultHash blockHashMsg
	if err := msg.Decode(&resultHash); err != nil {
		log.DLogger.Error("con't decode block hash msg", zap.Uint64("msg code", msg.Code), zap.Error(err))
		return err
	}

	if broadcaster.Chain.CurrentBlock().Number() >= resultHash.BlockNumber {
		return nil
	}

	log.DLogger.Info("receive new verified block hash msg", zap.String("node", p.NodeName()), zap.Uint64("height", resultHash.BlockNumber))
	// refresh peer's height
	p.SetHead(resultHash.BlockHash, resultHash.BlockNumber)

	// load receiver
	receiver := broadcaster.getReceiver(p)
	receiver.markVerifyResult(resultHash.BlockHash)

	if broadcaster.Chain.GetBlockByNumber(resultHash.BlockNumber) != nil {
		log.DLogger.Info("already have result block, no fetch block", zap.Uint64("height", resultHash.BlockNumber))
		return nil
	}

	// add fetch result block task
	go func() {

		log.DLogger.Debug("block fetcher notify", zap.Uint64("number", resultHash.BlockNumber))

		broadcaster.blockFetcher.Notify(p.ID(), resultHash.BlockHash, resultHash.BlockNumber, time.Now(), func(hash common.Hash) error {
			return p.SendMsg(GetVerifyResultMsg, resultHash.BlockNumber)
		})
	}()

	//go broadcaster.broadcastHash(&resultHash)

	return nil
}

func (broadcaster *BftOuter) onGetVerifiedResult(msg p2p.Msg, p PmAbstractPeer) error {
	var height uint64
	if err := msg.Decode(&height); err != nil {
		return err
	}

	block := broadcaster.Chain.GetBlockByNumber(height)
	seen := broadcaster.Chain.GetSeenCommit(height)
	if block == nil || len(seen) == 0 {
		log.DLogger.Error("can't load block for remote", zap.Int("seen len", len(seen)), zap.Bool("block_is_nil?", block == nil))
		return nil
	}

	receiver := broadcaster.getReceiver(p)
	receiver.asyncSendVerifyResult(&model2.VerifyResult{
		Block:       block,
		SeenCommits: seen,
	})

	return nil
}

// real verify result
func (broadcaster *BftOuter) onVerifiedResultBlock(msg p2p.Msg, p PmAbstractPeer) error {
	g_metrics.Add(g_metrics.ReceivedBlockCount, "", 1)

	var result model2.VerifyResultRlp
	if err := msg.Decode(&result); err != nil {
		log.DLogger.Warn("decode v result failed", zap.Error(err))
		return err
	}

	commits := make([]model.AbstractVerification, len(result.SeenCommits))
	util.InterfaceSliceCopy(commits, result.SeenCommits)

	// here will call the save block
	broadcaster.blockFetcher.DoTask(p.ID(), &model2.VerifyResult{
		Block:       &result.Block,
		SeenCommits: commits,
	}, time.Now())

	return nil
}

// get receiver
func (broadcaster *BftOuter) getReceiver(p PmAbstractPeer) *blockReceiver {
	// load txReceiver
	var receiver *blockReceiver

	if cache, ok := broadcaster.vResultBroadcast.Load(p.ID()); ok {
		receiver = cache.(*blockReceiver)
	} else {
		receiver = broadcaster.newBlockReceiver(p)
		broadcaster.vResultBroadcast.Store(p.ID(), receiver)
	}

	return receiver
}

func (broadcaster *BftOuter) newBlockReceiver(peer PmAbstractPeer) *blockReceiver {

	kb, _ := lru.New(500)

	receiver := &blockReceiver{
		peerID:                 peer.ID(),
		peerName:               peer.NodeName(),
		knownBlocks:            kb,
		queuedVerifyResult:     make(chan *model2.VerifyResult, maxQueuedBlock),
		queuedVerifyResultHash: make(chan *model2.VerifyResult, maxQueuedBlockHash),
	}

	go func() {
		defer func() {
			broadcaster.vResultBroadcast.Delete(peer.ID())
		}()

		getPeer := func() PmAbstractPeer {
			return broadcaster.Pm.GetPeer(peer.ID())
		}

		if err := receiver.broadcast(getPeer); err != nil {
			log.DLogger.Error("broadcast verified block result hash error", zap.Error(err), zap.String("peer name", peer.NodeName()))
			return
		}

	}()

	return receiver
}

// get peer without block
func (broadcaster *BftOuter) getPeersWithoutBlock(hash common.Hash) []PmAbstractPeer {
	// get peers
	peers := broadcaster.Pm.GetPeers()

	var list []PmAbstractPeer

	for _, p := range peers {
		receiver := broadcaster.getReceiver(p)

		if !receiver.knownBlocks.Contains(hash) {
			list = append(list, p)
		}
	}

	return list
}
