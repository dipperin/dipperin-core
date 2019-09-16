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

package csbftnode

import (
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/g-error"
	"github.com/dipperin/dipperin-core/common/g-metrics"
	"github.com/dipperin/dipperin-core/core/chain-communication"
	"github.com/dipperin/dipperin-core/core/csbft/components"
	model2 "github.com/dipperin/dipperin-core/core/csbft/model"
	"github.com/dipperin/dipperin-core/core/csbft/state-machine"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/dipperin/dipperin-core/third-party/p2p"
	"time"
)

// new bft node
func NewCsBft(config *state_machine.BftConfig) *CsBft {
	bft := &CsBft{BftConfig: config}
	bp := components.NewBlockPool(0, nil)
	bp.SetNodeConfig(config.ChainReader)
	stateHandler := state_machine.NewStateHandler(config, state_machine.DefaultConfig, bp)
	bp.SetPoolEventNotifier(stateHandler)
	bft.blockPool = bp
	bft.stateHandler = stateHandler
	return bft
}

type CsBft struct {
	*state_machine.BftConfig

	blockPool    *components.BlockPool
	stateHandler *state_machine.StateHandler
	fetcher      *components.CsBftFetcher
}

// when new block insert to chain, call here notify state change
func (bft *CsBft) OnEnterNewHeight(h uint64) {
	bft.stateHandler.NewHeight(h)
}

func (bft *CsBft) SetFetcher(fetcher *components.CsBftFetcher) {
	bft.fetcher = fetcher
	bft.stateHandler.SetFetcher(fetcher)
}

/*func (bft *CsBft) SendFetchBlockMsg(msgCode uint64, from common.Address, msg *model2.FetchBlockReqDecodeMsg) error {
    //return bft.nodeContext.FetcherConnAdaptCsBft().SendFetchBlockMsg(msgCode, from, msg)
    return bft.FetcherConnAdaptCsBft.SendFetchBlockMsg(msgCode, from, msg)
}*/

func (bft *CsBft) Start() error {
	log.PBft.Info("start CsBft", "cur height", bft.ChainReader.CurrentBlock().Number())
	if !bft.canStart() {
		log.PBft.Info("isn't cur verifier, can't start CsBft", "err", g_error.ErrIsNotCurVerifierCannotStartBft)
		return nil
	}

	if bft.stateHandler.IsRunning() && bft.blockPool.IsRunning() && bft.fetcher.IsRunning() {
		return nil
	}
	err := bft.stateHandler.Start()
	log.PBft.Debug("start git", "is running", bft.stateHandler.IsRunning(), "err", err)
	err = bft.blockPool.Start()
	log.PBft.Debug("start pool", "is running", bft.blockPool.IsRunning(), "err", err)
	err = bft.fetcher.Start()
	log.PBft.Debug("start fetcher", "is running", bft.fetcher.IsRunning(), "err", err)

	return nil
}

func (bft *CsBft) Stop() {
	log.PBft.Info("stop CsBft", "cur height", bft.ChainReader.CurrentBlock().Number())

	bft.stateHandler.Stop()
	if err := bft.stateHandler.Reset(); err != nil {
		log.Warn("reset state handler failed", "err", err)
	}
	log.PBft.Debug("Stop state handler", "state handler is running", bft.blockPool.IsRunning())
	bft.blockPool.Stop()
	log.PBft.Debug("Stop pool", "pool is running", bft.blockPool.IsRunning())
	bft.fetcher.Stop()
	log.PBft.Debug("Stop fetcher", "fetcher is running", bft.blockPool.IsRunning())
	bft.fetcher.Reset()
}

func (bft *CsBft) OnNewWaitVerifyBlock(block model.AbstractBlock, id string) {
	//log.PBft.Debug("cs onNewWatVerifyBlock")
	//check the node is or isn't current verifier node
	if !bft.stateHandler.IsRunning() || !bft.blockPool.IsRunning() {
		log.PBft.Debug("cs onNewWatVerifyBlock, bft not running")
		return
	}
	log.PBft.Info("cs bft OnNewWaitVerifyBlock", "block num", block.Number())

	// todo check block valid here?
	if err := bft.blockPool.AddBlock(block); err != nil {
		log.PBft.Info("pool add block failed", "err", err)
		return
	}
	// wait and sync block to other verifiers
	go bft.broadcastFetchBlockMsg(block.Hash())
}

func (bft *CsBft) broadcastFetchBlockMsg(blockHash common.Hash) {
	// maybe other node is receiving this block
	time.Sleep(500 * time.Millisecond)
	log.PBft.Info("broadcast sync block msg", "hash", blockHash.Hex())
	bft.Sender.BroadcastMsg(uint64(model2.TypeOfSyncBlockMsg), blockHash)
}

func (bft *CsBft) OnNewMsg(msg interface{}) error {
	return nil
}

func (bft *CsBft) AddPeer(p chain_communication.PmAbstractPeer) error { return nil }

func (bft *CsBft) ChangePrimary(primary string) {
	log.PBft.Debug("Change Primary Called")
	log.PBft.Debug("Current num", "num", bft.ChainReader.CurrentBlock().Number())
	if bft.canStart() {
		log.PBft.Debug("Start state handler")
		bft.Start()
		bft.stateHandler.NewHeight(bft.ChainReader.CurrentBlock().Number() + 1)
		return
	}
	log.PBft.Debug("Stop state handler")
	bft.Stop()
}

// determine whether it should start
func (bft *CsBft) canStart() bool {
	curB := bft.ChainReader.CurrentBlock()
	// The second parameter is true only if it is packaged. If it is a switch point, it should take next
	if bft.ChainReader.IsChangePoint(curB, false) {
		return bft.isNextVerifier()
	}
	return bft.isCurrentVerifier()
}

func (bft *CsBft) isCurrentVerifier() bool {
	vs := bft.ChainReader.GetCurrVerifiers()
	curAccount := bft.Signer.GetAddress()
	log.PBft.Info("CsBft isCurrentVerifier", "cur vs", vs, "cur account", curAccount)
	for _, v := range vs {
		if v.IsEqual(curAccount) {
			return true
		}
	}
	return false
}

func (bft *CsBft) isNextVerifier() bool {
	vs := bft.ChainReader.GetNextVerifiers()
	curAccount := bft.Signer.GetAddress()
	for _, v := range vs {
		if v.IsEqual(curAccount) {
			return true
		}
	}
	return false
}

// The processing here can't be blocked, it must be quickly put into a coroutine and returned after processing, otherwise msg read will be blocked.
func (bft *CsBft) OnNewP2PMsg(msg p2p.Msg, p chain_communication.PmAbstractPeer) error {
	if !bft.stateHandler.IsRunning() {
		log.PBft.Warn("[Node-OnNewMsg]receive bft msg, but state handler not started")
		return nil
	}

	switch model2.CsBftMsgType(msg.Code) {
	case model2.TypeOfNewRoundMsg:
		var m model2.NewRoundMsg
		if err := msg.Decode(&m); err != nil {
			return err
		}
		log.PBft.Info("[Node-OnNewMsg]receive new round msg", "node", p.NodeName(), "height", m.Height, "round", m.Round)
		bft.stateHandler.NewRound(&m)
	case model2.TypeOfProposalMsg:
		var m model2.Proposal
		if err := msg.Decode(&m); err != nil {
			return err
		}
		log.PBft.Info("[Node-OnNewMsg]receive proposal msg", "node", p.NodeName(), "height", m.Height, "round", m.Round, "block", m.BlockID.Hex())
		bft.stateHandler.NewProposal(&m)
	case model2.TypeOfPreVoteMsg:
		var m model.VoteMsg
		if err := msg.Decode(&m); err != nil {
			return err
		}
		log.PBft.Info("[Node-OnNewMsg]receive prevote msg", "node", p.NodeName(), "height", m.Height, "round", m.Round, "block", m.BlockID.Hex())
		bft.stateHandler.PreVote(&m)

	case model2.TypeOfVoteMsg:
		var m model.VoteMsg
		if err := msg.Decode(&m); err != nil {
			return err
		}
		log.PBft.Info("[Node-OnNewMsg]receive vote msg", "node", p.NodeName(), "height", m.Height, "round", m.Round, "block", m.BlockID.Hex())
		bft.stateHandler.Vote(&m)

	case model2.TypeOfFetchBlockReqMsg:
		//fmt.Println("receive fetch block msg")
		log.PBft.Info("[Node-OnNewMsg]receive fetch block msg", "from", p.NodeName())
		var m model2.FetchBlockReqDecodeMsg
		if err := msg.Decode(&m); err != nil {
			return err
		}

		b := bft.blockPool.GetBlockByHash(m.BlockHash)
		if b == nil {
			b = bft.stateHandler.GetProposalBlock(m.BlockHash)
		}
		log.PBft.Info("[Node-OnNewMsg] fetch result", "to", p.NodeName(), "block_is_nil", b == nil)

		// todo check will panic if b is nil?
		if b == nil {
			return nil
		}
		if err := p.SendMsg(uint64(model2.TypeOfFetchBlockRespMsg), &components.FetchBlockRespMsg{
			MsgId: m.MsgId,
			Block: b,
		}); err != nil {
			log.PBft.Warn("[Node-OnNewMsg] send fetch block to client failed", "err", err)
		}
		log.PBft.Info("[Node-OnNewMsg] send fetch result 2")

	case model2.TypeOfFetchBlockRespMsg:
		var m model2.FetchBlockRespDecodeMsg
		if err := msg.Decode(&m); err != nil {
			log.PBft.Debug("[Node-OnNewMsg] Decode Error, FetchBlockRespMsg", "err", err)
			return err
		}
		log.PBft.Info("[Node-OnNewMsg] receive fetch block resp", "node", p.NodeName())
		bft.fetcher.FetchBlockResp(&components.FetchBlockRespMsg{
			MsgId: m.MsgId,
			Block: m.Block,
		})

	case model2.TypeOfSyncBlockMsg:
		log.PBft.Info("[Node-OnNewMsg] receive sync block", "node", p.NodeName())
		var m common.Hash
		if err := msg.Decode(&m); err != nil {
			return err
		}
		// coroutine is obliged
		go bft.onSyncBlockMsg(p.RemoteVerifierAddress(), m)
	case model2.TypeOfReqNewRoundMsg:
		var m model2.ReqRoundMsg
		if err := msg.Decode(&m); err != nil {
			log.PBft.Error("decode req new round msg error", "err", err)
			return err
		}
		log.PBft.Info("[Node-OnNewMsg] receive req new round", "node", p.NodeName(), "height", m.Height, "round", m.Round)

		round := m.Round
		for msg := bft.stateHandler.GetRoundMsg(m.Height, round);msg != nil ; round++  {
			log.PBft.Info("[Node-OnNewMsg]  response round request", "to", p.NodeName(), "height", m.Height, "round", round)
			if err := p.SendMsg(uint64(model2.TypeOfNewRoundMsg), msg); err != nil {
				log.PBft.Error("response round request error", "err", err)
			}
		}

		//msg := bft.stateHandler.GetRoundMsg(m.Height, m.Round)
		//log.PBft.Debug("[Node-OnNewMsg] response", "msg == nil", msg == nil)
		//if msg != nil {
		//	log.PBft.Info("[Node-OnNewMsg]  response round request", "to", p.NodeName(), "height", m.Height, "round", m.Round)
		//	if err := p.SendMsg(uint64(model2.TypeOfNewRoundMsg), msg); err != nil {
		//		log.PBft.Error("response round request error", "err", err)
		//	}
		//}
	default:
		panic(fmt.Sprintf("unknown csbft msg, code: %v", msg.Code))
	}

	return nil
}

func (bft *CsBft) onSyncBlockMsg(from common.Address, h common.Hash) {
	g_metrics.Add(g_metrics.FetchBlockGoCount, "", 1)
	defer g_metrics.Sub(g_metrics.FetchBlockGoCount, "", 1)

	if from.IsEmpty() {
		log.PBft.Warn("from is empty, do nothing for sync msg")
		return
	}
	if h.IsEmpty() {
		log.PBft.Warn("block hash is empty, do nothing for sync msg")
		return
	}

	if !bft.blockPool.IsEmpty() {
		//log.PBft.Warn("pool not empty, ignore sync block msg")
		return
	}

	// check have this block?
	b := bft.blockPool.GetBlockByHash(h)
	if b == nil {
		b = bft.stateHandler.GetProposalBlock(h)
	}
	if b != nil {
		log.PBft.Info("onSyncBlockMsg already have this block")
		return
	}

	// synchronous acquisition of a
	b = bft.fetcher.FetchBlock(from, h)
	if b != nil {
		if err := bft.blockPool.AddBlock(b); err != nil {
			log.PBft.Warn("fetcher add block failed", "err", err)
		}
		return
	}
	log.PBft.Info("fetch block failed")
}
