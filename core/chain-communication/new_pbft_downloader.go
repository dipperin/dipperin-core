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
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/dipperin/dipperin-core/third-party/p2p"
	"github.com/ethereum/go-ethereum/rlp"
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
	log.Info("receive get blocks msg1")
	var query getBlockHeaders
	if err := msg.Decode(&query); err != nil {
		return errors.New("decode error, invalid message")
	}
	log.Info("receive get blocks msg2", "OriginHeight", query.OriginHeight, "amount", query.Amount, "remote node", p.NodeName())

	var (
		blocks []*catchup
	)

	for len(blocks) < int(query.Amount) && len(blocks) < MaxBlockFetch {
		origin := fd.Chain.GetBlockByNumber(query.OriginHeight)
		cmt := fd.Chain.GetSeenCommit(query.OriginHeight)
		if origin == nil {
			log.Info("can't get block for downloader", "height", query.OriginHeight)
			break
		}

		blocks = append(blocks, &catchup{Block: origin, SeenCommit: cmt})
		query.OriginHeight += 1
	}

	// todo rm here
	//bb, _ := rlp.EncodeToBytes(blocks)
	//, "data size(Mb)", len(bb) * 8.0 / 1024 / 1024
	log.Info("downloader send blocks to remote", "remote node", p.NodeName(), "block len", len(blocks))
	return p.SendMsg(BlocksMsg, blocks)
}

func (fd *NewPbftDownloader) onBlocks(msg p2p.Msg, p PmAbstractPeer) error {
	var blocks []*catchupRlp
	if err := msg.Decode(&blocks); err != nil {
		log.Error("downloader decode blocks failed", "err", err)
		return err
	}

	log.Info("downloader receive blocks from", "node", p.NodeName(), "block len", len(blocks))
	if len(blocks) > 0 {
		// Filter out any explicitly requested block vr, deliver the rest to the downloader
		filter := fd.fetcher.DoFilter(p.ID(), blocks)
		log.Info("fetcher filter catchup list", "origin size", len(blocks), "filter", len(filter))
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
	log.PBft.Debug("Start New PBFT Downloader")
	go fd.loop()
	return nil
}

func (fd *NewPbftDownloader) loop() {
	tickHandler := func() {
		log.PBft.Debug("pbft downloader call download")
		if fd.getBestPeer() == nil {
			log.Warn("downloader can't get best peer, do nothing")
			return
		}

		log.Info("downloader run sync")
		fd.runSync()
	}
	forceSync := g_timer.SetPeriodAndRun(tickHandler, pollingInterval)
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
		log.Info("====================bestPeer == nil====================")
		return nil
	}

	currentBlock := fd.Chain.CurrentBlock()
	if currentBlock == nil {
		log.Error("=======================currentBlock is nil===================")
		return nil
	}

	// check local blockchain current block height < bestPeer height
	_, height := bestPeer.GetHead()
	if height <= currentBlock.Number() {
		log.Info("local higher than best peer", "bestPeer", bestPeer.NodeName(), "remote h", height, "local h", currentBlock.Number())
		return nil
	}

	log.Info("downloader got best peer", "p height", height, "p name", bestPeer.NodeName())
	return bestPeer
}

//run synchronise
func (fd *NewPbftDownloader) runSync() {
	log.Info("bft downloader run sync")
	log.PBft.Debug("bft downloader run sync")
	if !atomic.CompareAndSwapInt32(&fd.synchronising, 0, 1) {
		log.Info("downloader is busy")
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

	log.Info("send get blocks msg", "remote peer name", bestPeer.NodeName(), "remote peer height", height)
	go func() {
		if err := bestPeer.SendMsg(GetBlocksMsg, &getBlockHeaders{OriginHeight: nextNumber, Amount: MaxBlockFetch}); err != nil {
			log.Warn("send get blocks msg failed", "err", err)
		}
	}()

	timeoutTimer := time.NewTimer(fetchBlockTimeout)
	defer timeoutTimer.Stop()
	for {
		select {
		case packet := <-fd.blockC:
			if packet.peerID != bestPeer.ID() {
				log.Warn("Received skeleton from incorrect peer", "peer", packet.peerID)
				break
			}

			blocks := packet.blocks
			size := len(blocks)
			log.Info("fetchBlocks1", "blocks len", size)

			// no block return from remote
			if size == 0 {
				return
			}

			// If the insertion is slow, it will cause a timeout, then the Peer is broken.
			if size > 0 {
				if err := fd.importBlockResults(blocks); err != nil {
					log.Error("downloader save block failed", "err", err, "remote node", bestPeer.NodeName())
					return
				}
				nextNumber += uint64(len(blocks))
			}

			// stop when equal remote height
			if height == fd.Chain.CurrentBlock().Number() {
				return
			}

			// reset timeout timer
			timeoutTimer.Reset(fetchBlockTimeout)
			// run sync again
			go bestPeer.SendMsg(GetBlocksMsg, &getBlockHeaders{OriginHeight: nextNumber, Amount: MaxBlockFetch})
		case <-timeoutTimer.C:
			log.Warn("Waiting for fetchHeaders headers timed out", "node name", bestPeer.NodeName())
			return

		case <-fd.quitCh:
			return

		}
	}
}

func (fd *NewPbftDownloader) importBlockResults(list []*catchupRlp) error {
	log.Info("insert blocks from downloader", "len", len(list))
	for _, b := range list {
		commits := make([]model.AbstractVerification, len(b.SeenCommit))
		util.InterfaceSliceCopy(commits, b.SeenCommit)

		if len(commits) > 0 {
			log.PBft.Debug("pbft download call save block", "block height", b.Block.Number(), "commits", len(commits), "commits", commits[0].GetBlockId().Hex())
		} else {
			log.PBft.Warn("commits is empty", "height", b.Block.Number())
		}

		log.Info("importBlockResults save block number is:", "blockNumber", b.Block.Number())
		if err := fd.Chain.SaveBlock(b.Block, commits); err != nil {
			if err == g_error.ErrNormalBlockHeightTooLow {
				log.Info("importBlockResults the block height is same as the current block ")
				continue
			} else {
				return err
			}
		}
	}

	return nil
}
