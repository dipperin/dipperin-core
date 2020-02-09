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
	"github.com/dipperin/dipperin-core/common/gerror"
	"github.com/dipperin/dipperin-core/common/gmetrics"
	"github.com/dipperin/dipperin-core/common/log"
	"github.com/dipperin/dipperin-core/core/chaincommunication"
	"github.com/dipperin/dipperin-core/core/csbft/components"
	model2 "github.com/dipperin/dipperin-core/core/csbft/model"
	"github.com/dipperin/dipperin-core/core/csbft/statemachine"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third_party/p2p"
	"go.uber.org/zap"
	"time"
)

//go:generate mockgen -destination=./verification_mock_test.go -package=csbftnode github.com/dipperin/dipperin-core/core/model AbstractVerification
//go:generate mockgen -destination=./block_mock_test.go -package=csbftnode github.com/dipperin/dipperin-core/core/model AbstractBlock
//go:generate mockgen -destination=./node_mock_test.go -package=csbftnode github.com/dipperin/dipperin-core/core/csbft/statemachine ChainReader,MsgSigner,MsgSender,Validator,Fetcher
// new bft node
func NewCsBft(config *statemachine.BftConfig) *CsBft {
	bft := &CsBft{BftConfig: config}
	bp := components.NewBlockPool(0, nil)
	bp.SetNodeConfig(config.ChainReader)
	stateHandler := statemachine.NewStateHandler(config, statemachine.DefaultConfig, bp)
	bp.SetPoolEventNotifier(stateHandler)
	bft.blockPool = bp
	bft.stateHandler = stateHandler
	return bft
}

type CsBft struct {
	*statemachine.BftConfig

	blockPool    *components.BlockPool
	stateHandler *statemachine.StateHandler
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
	log.DLogger.Info("start CsBft", zap.Uint64("cur height", bft.ChainReader.CurrentBlock().Number()))
	if !bft.canStart() {
		log.DLogger.Info("isn't cur verifier, can't start CsBft", zap.Error(gerror.ErrIsNotCurVerifierCannotStartBft))
		return nil
	}

	if bft.stateHandler.IsRunning() && bft.blockPool.IsRunning() && bft.fetcher.IsRunning() {
		return nil
	}
	err := bft.stateHandler.Start()
	log.DLogger.Debug("start git", zap.Bool("is running", bft.stateHandler.IsRunning()), zap.Error(err))
	err = bft.blockPool.Start()
	log.DLogger.Debug("start pool", zap.Bool("is running", bft.blockPool.IsRunning()), zap.Error(err))
	err = bft.fetcher.Start()
	log.DLogger.Debug("start fetcher", zap.Bool("is running", bft.fetcher.IsRunning()), zap.Error(err))

	return nil
}

func (bft *CsBft) Stop() {
	log.DLogger.Info("stop CsBft", zap.Uint64("cur height", bft.ChainReader.CurrentBlock().Number()))

	bft.stateHandler.Stop()
	if err := bft.stateHandler.Reset(); err != nil {
		log.DLogger.Warn("reset state handler failed", zap.Error(err))
	}
	log.DLogger.Debug("Stop state handler", zap.Bool("state handler is running", bft.blockPool.IsRunning()))
	bft.blockPool.Stop()
	log.DLogger.Debug("Stop pool", zap.Bool("pool is running", bft.blockPool.IsRunning()))
	bft.fetcher.Stop()
	log.DLogger.Debug("Stop fetcher", zap.Bool("fetcher is running", bft.blockPool.IsRunning()))
	bft.fetcher.Reset()
}

func (bft *CsBft) OnNewWaitVerifyBlock(block model.AbstractBlock, id string) {
	//log.DLogger.Debug("cs onNewWatVerifyBlock")
	//check the node is or isn't current verifier node
	if !bft.stateHandler.IsRunning() || !bft.blockPool.IsRunning() {
		log.DLogger.Debug("cs onNewWatVerifyBlock, bft not running")
		return
	}
	log.DLogger.Info("cs bft OnNewWaitVerifyBlock", zap.Uint64("block num", block.Number()))

	// todo check block valid here?
	if err := bft.blockPool.AddBlock(block); err != nil {
		log.DLogger.Info("pool add block failed", zap.Error(err))
		return
	}
	// wait and sync block to other verifiers
	go bft.broadcastFetchBlockMsg(block.Hash())
}

func (bft *CsBft) broadcastFetchBlockMsg(blockHash common.Hash) {
	// maybe other node is receiving this block
	time.Sleep(500 * time.Millisecond)
	log.DLogger.Info("broadcast sync block msg", zap.String("hash", blockHash.Hex()))
	bft.Sender.BroadcastMsg(uint64(model2.TypeOfSyncBlockMsg), blockHash)
}

func (bft *CsBft) OnNewMsg(msg interface{}) error {
	return nil
}

func (bft *CsBft) AddPeer(p chaincommunication.PmAbstractPeer) error { return nil }

func (bft *CsBft) ChangePrimary(primary string) {
	log.DLogger.Debug("Change Primary Called")
	log.DLogger.Debug("Current num", zap.Uint64("num", bft.ChainReader.CurrentBlock().Number()))
	if bft.canStart() {
		log.DLogger.Debug("Start state handler")
		bft.Start()
		bft.stateHandler.NewHeight(bft.ChainReader.CurrentBlock().Number() + 1)
		return
	}
	log.DLogger.Debug("Stop state handler")
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
	log.DLogger.Info("CsBft isCurrentVerifier", zap.Any("cur vs", vs), zap.Any("cur account", curAccount))
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
func (bft *CsBft) OnNewP2PMsg(msg p2p.Msg, p chaincommunication.PmAbstractPeer) error {
	if !bft.stateHandler.IsRunning() {
		log.DLogger.Warn("[Node-OnNewMsg]receive bft msg, but state handler not started")
		return nil
	}

	switch model2.CsBftMsgType(msg.Code) {
	case model2.TypeOfNewRoundMsg:
		var m model2.NewRoundMsg
		if err := msg.Decode(&m); err != nil {
			return err
		}
		log.DLogger.Info("[Node-OnNewMsg]receive new round msg", zap.String("node", p.NodeName()), zap.Uint64("height", m.Height), zap.Uint64("round", m.Round))
		bft.stateHandler.NewRound(&m)
	case model2.TypeOfProposalMsg:
		var m model2.Proposal
		if err := msg.Decode(&m); err != nil {
			return err
		}
		log.DLogger.Info("[Node-OnNewMsg]receive proposal msg", zap.String("node", p.NodeName()), zap.Uint64("height", m.Height), zap.Uint64("round", m.Round), zap.String("block", m.BlockID.Hex()))
		bft.stateHandler.NewProposal(&m)
	case model2.TypeOfPreVoteMsg:
		var m model.VoteMsg
		if err := msg.Decode(&m); err != nil {
			return err
		}
		log.DLogger.Info("[Node-OnNewMsg]receive prevote msg", zap.String("node", p.NodeName()), zap.Uint64("height", m.Height), zap.Uint64("round", m.Round), zap.String("block", m.BlockID.Hex()))
		bft.stateHandler.PreVote(&m)

	case model2.TypeOfVoteMsg:
		var m model.VoteMsg
		if err := msg.Decode(&m); err != nil {
			return err
		}
		log.DLogger.Info("[Node-OnNewMsg]receive vote msg", zap.String("node", p.NodeName()), zap.Uint64("height", m.Height), zap.Uint64("round", m.Round), zap.String("block", m.BlockID.Hex()))
		bft.stateHandler.Vote(&m)

	case model2.TypeOfFetchBlockReqMsg:
		//fmt.Println("receive fetch block msg")
		log.DLogger.Info("[Node-OnNewMsg]receive fetch block msg", zap.String("from", p.NodeName()))
		var m model2.FetchBlockReqDecodeMsg
		if err := msg.Decode(&m); err != nil {
			return err
		}

		b := bft.blockPool.GetBlockByHash(m.BlockHash)
		if b == nil {
			b = bft.stateHandler.GetProposalBlock(m.BlockHash)
		}
		log.DLogger.Info("[Node-OnNewMsg] fetch result", zap.String("to", p.NodeName()), zap.Bool("block_is_nil", b == nil))

		// todo check will panic if b is nil?
		if b == nil {
			return nil
		}
		if err := p.SendMsg(uint64(model2.TypeOfFetchBlockRespMsg), &components.FetchBlockRespMsg{
			MsgId: m.MsgId,
			Block: b,
		}); err != nil {
			log.DLogger.Warn("[Node-OnNewMsg] send fetch block to client failed", zap.Error(err))
		}
		log.DLogger.Info("[Node-OnNewMsg] send fetch result 2")

	case model2.TypeOfFetchBlockRespMsg:
		var m model2.FetchBlockRespDecodeMsg
		if err := msg.Decode(&m); err != nil {
			log.DLogger.Debug("[Node-OnNewMsg] Decode Error, FetchBlockRespMsg", zap.Error(err))
			return err
		}
		log.DLogger.Info("[Node-OnNewMsg] receive fetch block resp", zap.String("node", p.NodeName()))
		bft.fetcher.FetchBlockResp(&components.FetchBlockRespMsg{
			MsgId: m.MsgId,
			Block: m.Block,
		})

	case model2.TypeOfSyncBlockMsg:
		log.DLogger.Info("[Node-OnNewMsg] receive sync block", zap.String("node", p.NodeName()))
		var m common.Hash
		if err := msg.Decode(&m); err != nil {
			return err
		}
		// coroutine is obliged
		go bft.onSyncBlockMsg(p.RemoteVerifierAddress(), m)
	case model2.TypeOfReqNewRoundMsg:
		var m model2.ReqRoundMsg
		if err := msg.Decode(&m); err != nil {
			log.DLogger.Error("decode req new round msg error", zap.Error(err))
			return err
		}
		log.DLogger.Info("[Node-OnNewMsg] receive req new round", zap.String("node", p.NodeName()), zap.Uint64("height", m.Height), zap.Uint64("round", m.Round))

		round := m.Round
		for {
			msg := bft.stateHandler.GetRoundMsg(m.Height, round)
			if msg != nil {
				log.DLogger.Info("[Node-OnNewMsg]  response round request", zap.String("to", p.NodeName()), zap.Uint64("height", m.Height), zap.Uint64("round", m.Round), zap.Any("msg", msg))
				if err := p.SendMsg(uint64(model2.TypeOfNewRoundMsg), msg); err != nil {
					log.DLogger.Error("response round request error", zap.Error(err))
				}
			} else {
				break
			}
			round++
		}

		//msg := bft.stateHandler.GetRoundMsg(m.Height, m.Round)
		//log.DLogger.Debug("[Node-OnNewMsg] response", "msg == nil", msg == nil)
		//if msg != nil {
		//	log.DLogger.Info("[Node-OnNewMsg]  response round request", "to", p.NodeName(), "height", m.Height, "round", m.Round)
		//	if err := p.SendMsg(uint64(model2.TypeOfNewRoundMsg), msg); err != nil {
		//		log.DLogger.Error("response round request error", zap.Error(err))
		//	}
		//}
	default:
		panic(fmt.Sprintf("unknown csbft msg, code: %v", msg.Code))
	}

	return nil
}

func (bft *CsBft) onSyncBlockMsg(from common.Address, h common.Hash) {
	gmetrics.Add(gmetrics.FetchBlockGoCount, "", 1)
	defer gmetrics.Sub(gmetrics.FetchBlockGoCount, "", 1)

	if from.IsEmpty() {
		log.DLogger.Warn("from is empty, do nothing for sync msg")
		return
	}
	if h.IsEmpty() {
		log.DLogger.Warn("block hash is empty, do nothing for sync msg")
		return
	}

	if !bft.blockPool.IsEmpty() {
		//log.DLogger.Warn("pool not empty, ignore sync block msg")
		return
	}

	// check have this block?
	b := bft.blockPool.GetBlockByHash(h)
	if b == nil {
		b = bft.stateHandler.GetProposalBlock(h)
	}
	if b != nil {
		log.DLogger.Info("onSyncBlockMsg already have this block")
		return
	}

	// synchronous acquisition of a
	b = bft.fetcher.FetchBlock(from, h)
	if b != nil {
		if err := bft.blockPool.AddBlock(b); err != nil {
			log.DLogger.Warn("fetcher add block failed", zap.Error(err))
		}
		return
	}
	log.DLogger.Info("fetch block failed")
}
