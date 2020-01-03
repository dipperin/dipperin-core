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
	"github.com/dipperin/dipperin-core/common/g-error"
	"github.com/dipperin/dipperin-core/common/g-timer"
	"github.com/dipperin/dipperin-core/common/log"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third-party/p2p"
	"github.com/ethereum/go-ethereum/rlp"
	"go.uber.org/zap"
	"io"
	"strings"
	"sync/atomic"
	"time"
)

var (
	pollingInterval   = 10 * time.Second
	fetchBlockTimeout = 60 * time.Second
)

func MakeNewPbftDownloader(config *NewPbftDownloaderConfig) *NewPbftDownloader {
	service := &NewPbftDownloader{
		NewPbftDownloaderConfig: config,
		handlers:                map[uint64]func(msg p2p.Msg, p PmAbstractPeer) error{},
		blockC:                  make(chan *npbPack),
		quitCh:                  make(chan struct{}),
	}
	service.handlers[GetBlocksMsg] = service.onGetBlocks
	service.handlers[BlocksMsg] = service.onBlocks
	return service
}

type NewPbftDownloaderConfig struct {
	Chain    Chain
	Pm       PeerManager
	PbftNode PbftNode
	//fetcher  *EiBlockFetcher
	fetcher *BlockFetcher
}

type NewPbftDownloader struct {
	*NewPbftDownloaderConfig
	handlers      map[uint64]func(msg p2p.Msg, p PmAbstractPeer) error
	blockC        chan *npbPack
	quitCh        chan struct{}
	synchronising int32
}

type npbPack struct {
	peerID string
	blocks []*catchupRlp
}

type catchup struct {
	Block      model.AbstractBlock
	SeenCommit []model.AbstractVerification
}
type catchupRlp struct {
	Block      *model.Block
	SeenCommit []*model.VoteMsg
}

func (c *catchup) DecodeRLP(s *rlp.Stream) error {
	var from catchupRlp
	if err := s.Decode(&from); err != nil {
		return err
	}

	c.Block = from.Block
	c.SeenCommit = make([]model.AbstractVerification, len(from.SeenCommit))
	util.InterfaceSliceCopy(c.SeenCommit, from.SeenCommit)
	return nil
}

func (fd *NewPbftDownloader) MsgHandlers() map[uint64]func(msg p2p.Msg, p PmAbstractPeer) error {
	return fd.handlers
}

// TODO
func (fd *NewPbftDownloader) onGetBlocks(msg p2p.Msg, p PmAbstractPeer) error {
	log.DLogger.Info("receive get blocks msg1")
	var query getBlockHeaders
	if err := msg.Decode(&query); err != nil {
		return errors.New("decode error, invalid message")
	}
	log.DLogger.Info("receive get blocks msg2", zap.Uint64("OriginHeight", query.OriginHeight), zap.Uint64("amount", query.Amount), zap.String("remote node", p.NodeName()))

	var (
		blocks []*catchup
	)

	for len(blocks) < int(query.Amount) && len(blocks) < MaxBlockFetch {
		origin := fd.Chain.GetBlockByNumber(query.OriginHeight)
		cmt := fd.Chain.GetSeenCommit(query.OriginHeight)
		if origin == nil {
			log.DLogger.Info("can't get block for downloader", zap.Uint64("height", query.OriginHeight))
			break
		}

		blocks = append(blocks, &catchup{Block: origin, SeenCommit: cmt})
		query.OriginHeight += 1
	}

	// todo rm here
	//bb, _ := rlp.EncodeToBytes(blocks)
	//, "data size(Mb)", len(bb) * 8.0 / 1024 / 1024
	log.DLogger.Info("downloader send blocks to remote", zap.String("remote node", p.NodeName()), zap.Int("block len", len(blocks)))
	return p.SendMsg(BlocksMsg, blocks)
}

func (fd *NewPbftDownloader) onBlocks(msg p2p.Msg, p PmAbstractPeer) error {
	var blocks []*catchupRlp
	if err := msg.Decode(&blocks); err != nil {
		log.DLogger.Error("downloader decode blocks failed", zap.Error(err))
		return err
	}

	log.DLogger.Info("downloader receive blocks from", zap.String("node", p.NodeName()), zap.Int("block len", len(blocks)))
	if len(blocks) > 0 {
		// Filter out any explicitly requested block vr, deliver the rest to the downloader
		filter := fd.fetcher.DoFilter(p.ID(), blocks)
		log.DLogger.Info("fetcher filter catchup list", zap.Int("origin size", len(blocks)), zap.Int("filter", len(filter)))
		pack := &npbPack{
			peerID: p.ID(),
			blocks: filter,
		}
		select {
		case <-fd.quitCh:
			return quitErr
		case fd.blockC <- pack:
		}
	}

	return nil
}

func (fd *NewPbftDownloader) Start() error {
	log.DLogger.Debug("Start New PBFT Downloader")
	go fd.loop()
	return nil
}

func (fd *NewPbftDownloader) loop() {
	//tickHandler := func() {
	//	log.PBft.Debug("pbft downloader call download")
	//	if fd.getBestPeer() == nil {
	//		log.DLogger.Warn("downloader can't get best peer, do nothing")
	//		return
	//	}
	//
	//	log.DLogger.Info("downloader run sync")
	//	fd.runSync()
	//}
	forceSync := g_timer.SetPeriodAndRun(fd.runSync, pollingInterval)
	defer g_timer.StopWork(forceSync)

	<-fd.quitCh
}

func (fd *NewPbftDownloader) Stop() {
	close(fd.quitCh)
}

func (fd *NewPbftDownloader) getBestPeer() PmAbstractPeer {
	bestPeer := fd.Pm.BestPeer()

	// ensure bestPeer no nil
	if bestPeer == nil {
		log.DLogger.Info("====================bestPeer == nil====================")
		return nil
	}

	currentBlock := fd.Chain.CurrentBlock()
	if currentBlock == nil {
		log.DLogger.Error("=======================currentBlock is nil===================")
		return nil
	}

	// check local blockchain current block height < bestPeer height
	_, height := bestPeer.GetHead()
	if height <= currentBlock.Number() {
		log.DLogger.Info("local higher than best peer", zap.String("bestPeer", bestPeer.NodeName()), zap.Uint64("remote h", height), zap.Uint64("local h", currentBlock.Number()))
		return nil
	}

	log.DLogger.Info("downloader got best peer", zap.Uint64("p height", height), zap.String("p name", bestPeer.NodeName()))
	return bestPeer
}

//run synchronise
func (fd *NewPbftDownloader) runSync() {
	log.DLogger.Info("bft downloader run sync")
	if !atomic.CompareAndSwapInt32(&fd.synchronising, 0, 1) {
		log.DLogger.Info("downloader is busy")
		return
	}
	defer atomic.StoreInt32(&fd.synchronising, 0)

	// clear old blocks in chan
	for empty := false; !empty; {
		select {
		case <-fd.blockC:
		default:
			empty = true
		}
	}

	bestPeer := fd.getBestPeer()
	if bestPeer == nil {
		return
	}

	fd.fetchBlocks(bestPeer)
}

func (fd *NewPbftDownloader) fetchBlocks(bestPeer PmAbstractPeer) {
	_, height := bestPeer.GetHead()

	//current block may be reversed by the empty block
	rollBackNum := chain_config.GetChainConfig().RollBackNum
	curNumber := fd.Chain.CurrentBlock().Number()
	var nextNumber uint64
	if curNumber > rollBackNum {
		nextNumber = fd.Chain.GetBlockByNumber(curNumber - rollBackNum + 1).Number()
	}

	log.DLogger.Info("send get blocks msg", zap.String("remote peer name", bestPeer.NodeName()), zap.Uint64("remote peer height", height))

	if err := bestPeer.SendMsg(GetBlocksMsg, &getBlockHeaders{OriginHeight: nextNumber, Amount: MaxBlockFetch}); err != nil {
		log.DLogger.Error("first send get blocks msg failed", zap.Error(err))
		return
	}

	timeoutTimer := time.NewTimer(fetchBlockTimeout)
	defer timeoutTimer.Stop()

	quitCh := make(chan struct{})
	defer close(quitCh)
	for {
		select {
		case packet := <-fd.blockC:
			if packet.peerID != bestPeer.ID() {
				log.DLogger.Warn("Received skeleton from incorrect peer", zap.String("peer", packet.peerID))
				break
			}

			blocks := packet.blocks
			size := len(blocks)
			log.DLogger.Info("fetchBlocks1", zap.Int("blocks len", size))

			// no block return from remote
			if size <= 0 {
				return
			}

			// If the insertion is slow, it will cause a timeout, then the Peer is broken.
			if err := fd.importBlockResults(blocks); err != nil {
				log.DLogger.Error("downloader save block failed", zap.Error(err), zap.String("remote node", bestPeer.NodeName()))
				return
			}
			nextNumber += uint64(size)

			// stop when equal remote height
			if height == fd.Chain.CurrentBlock().Number() {
				return
			}

			// reset timeout timer
			timeoutTimer.Reset(fetchBlockTimeout)

			// run sync again
			go func() {
				if !bestPeer.IsRunning() {
					quitCh <- struct{}{}
					return
				}
				if err := bestPeer.SendMsg(GetBlocksMsg, &getBlockHeaders{OriginHeight: nextNumber, Amount: MaxBlockFetch}); err != nil {
					if err.Error() == io.EOF.Error() || strings.Contains(err.Error(), "use of closed network connection") {
						quitCh <- struct{}{}
						return
					}
					log.DLogger.Error("run sync send msg: GetBlocksMsg", zap.Error(err))
				}
			}()
		case <-timeoutTimer.C:
			log.DLogger.Warn("Waiting for fetchHeaders headers timed out", zap.String("node name", bestPeer.NodeName()))
			return

		case <-fd.quitCh:
			return

		case <-quitCh:
			log.DLogger.Warn("peer disconnect: GetBlocksMsg", zap.String("node name", bestPeer.NodeName()))
			return
		}
	}
}

func (fd *NewPbftDownloader) importBlockResults(list []*catchupRlp) error {
	log.DLogger.Info("insert blocks from downloader", zap.Int("len", len(list)))
	for _, b := range list {
		commits := make([]model.AbstractVerification, len(b.SeenCommit))
		util.InterfaceSliceCopy(commits, b.SeenCommit)

		if len(commits) > 0 {
			log.DLogger.Debug("pbft download call save block", zap.Uint64("block height", b.Block.Number()), zap.Int("commits", len(commits)), zap.String("commits", commits[0].GetBlockId().Hex()))
		} else {
			log.DLogger.Warn("commits is empty", zap.Uint64("height", b.Block.Number()))
		}

		log.DLogger.Info("importBlockResults save block number is:", zap.Uint64("blockNumber", b.Block.Number()))
		if err := fd.Chain.SaveBlock(b.Block, commits); err != nil {
			if err == g_error.ErrNormalBlockHeightTooLow {
				log.DLogger.Info("importBlockResults the block height is same as the current block ")
				continue
			} else {
				return err
			}
		}
	}

	return nil
}
