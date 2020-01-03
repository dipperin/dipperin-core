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

package components

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/log"
	"github.com/dipperin/dipperin-core/common/util"
	model2 "github.com/dipperin/dipperin-core/core/csbft/model"
	"github.com/dipperin/dipperin-core/core/model"
	"go.uber.org/zap"
	"time"
)

const (
	fetchTimeout = 12 * time.Second
)

func NewFetcher(fc FetcherConn) *CsBftFetcher {
	fetcher := &CsBftFetcher{
		fc:             fc,
		requests:       make(map[uint64]*FetchBlockReqMsg),
		fetchReqQueue:  make(chan *FetchBlockReqMsg, 1),
		fetchRespChan:  make(chan *FetchBlockRespMsg, 1),
		isFetchingChan: make(chan *IsFetchingMsg),
		rmReqChan:      make(chan uint64),
	}
	fetcher.BaseService = *util.NewBaseService(log.DLogger, "cs_bft_fetcher", fetcher)
	log.DLogger.Info("NewFetcher", zap.Any("fetcher.BaseService", fetcher.BaseService))
	return fetcher
}

type CsBftFetcher struct {
	util.BaseService

	fc FetcherConn

	requests map[uint64]*FetchBlockReqMsg

	fetchReqQueue  chan *FetchBlockReqMsg
	fetchRespChan  chan *FetchBlockRespMsg
	isFetchingChan chan *IsFetchingMsg
	rmReqChan      chan uint64
}

func (f *CsBftFetcher) FetchBlock(from common.Address, blockHash common.Hash) model.AbstractBlock {
	if !f.IsRunning() {
		log.DLogger.Error("call fetch block, but fetcher not started")
		return nil
	}
	log.DLogger.Info("CsBftFetcher#FetchBlock  call fetch block", zap.String("block hash", blockHash.Hex()), zap.String("from", from.Hex()))
	req := &FetchBlockReqMsg{
		MsgId:      uint64(time.Now().UnixNano()),
		From:       from,
		BlockHash:  blockHash,
		ResultChan: make(chan model.AbstractBlock, 1),
	}

	f.fetchReqQueue <- req

	//log.DLogger.Info("CsBftFetcher#FetchBlock: ", "req", req,"from", from)
	select {
	case result := <-req.ResultChan:
		f.rmReq(req.MsgId)
		//log.DLogger.Info("CsBftFetcher#FetchBlock fetch block success", "block hash", result.Hash().Hex(), "from", from.Hex() )
		return result

	case <-time.After(fetchTimeout):
		log.DLogger.Warn("fetch block timeout", zap.String("block hash", blockHash.Hex()), zap.String("from", from.Hex()))
		// rm req
		f.rmReq(req.MsgId)
	}

	return nil
}

func (f *CsBftFetcher) FetchBlockResp(resp *FetchBlockRespMsg) {
	if f.IsRunning() {
		f.fetchRespChan <- resp
	} else {
		log.DLogger.Warn("receive fetch block resp, but fetcher not started")
		log.DLogger.Warn("receive fetch block resp, but fetcher not started")
	}
}

func (f *CsBftFetcher) rmReq(id uint64) {
	if f.IsRunning() {
		f.rmReqChan <- id
	}
}

func (f *CsBftFetcher) OnStart() error {
	go f.loop()
	return nil
}

func (f *CsBftFetcher) OnStop() {}

func (f *CsBftFetcher) OnReset() error { return nil }

func (f *CsBftFetcher) loop() {
	// clear old requests
	f.requests = make(map[uint64]*FetchBlockReqMsg)
	for {
		select {
		case req := <-f.fetchReqQueue:
			f.onFetchBlock(req)

		case resp := <-f.fetchRespChan:
			f.onFetchResp(resp)

		case msg := <-f.isFetchingChan:
			msg.Result <- f.isFetching(msg.BlockHash)

		case rId := <-f.rmReqChan:
			delete(f.requests, rId)

		case <-f.Quit():
			log.DLogger.Info("bft fetcher stopped")
			return
		}
	}
}

func (f *CsBftFetcher) IsFetching(hash common.Hash) bool {
	if !f.IsRunning() {
		return false
	}
	rc := make(chan bool)
	f.isFetchingChan <- &IsFetchingMsg{BlockHash: hash, Result: rc}
	return <-rc
}

// only initiated by loop, otherwise concurrence problem will occur
func (f *CsBftFetcher) isFetching(h common.Hash) bool {
	for _, r := range f.requests {
		if r.BlockHash.IsEqual(h) {
			return true
		}
	}
	return false
}

// action when request message of FetchBlock is received
func (f *CsBftFetcher) onFetchBlock(req *FetchBlockReqMsg) {
	if f.isFetching(req.BlockHash) {
		req.onResult(nil)
		log.DLogger.Info("is fetching block", zap.Any("hash", req.BlockHash))
		return
	}
	if len(f.requests) > 5 {
		req.onResult(nil)
		log.DLogger.Warn("too many fetches", zap.Int("req len", len(f.requests)))
		return
	}
	f.requests[req.MsgId] = req
	if err := f.fc.SendFetchBlockMsg(uint64(model2.TypeOfFetchBlockReqMsg), req.From, &model2.FetchBlockReqDecodeMsg{
		MsgId:     uint64(req.MsgId),
		BlockHash: req.BlockHash,
	}); err != nil {
		req.onResult(nil)
		log.DLogger.Warn("send fetch req failed", zap.Error(err))
		return
	}
}

// action when response message of FetchBlock is received
func (f *CsBftFetcher) onFetchResp(resp *FetchBlockRespMsg) {
	req := f.requests[resp.MsgId]
	if req == nil {
		log.DLogger.Info("receive fetch block resp, but req has been removed")
		return
	}

	log.DLogger.Info("onFetchResp1", zap.Uint64("block height", resp.Block.Number()))
	req.onResult(resp.Block)
}

type IsFetchingMsg struct {
	BlockHash common.Hash
	Result    chan bool
}

type FetchBlockReqMsg struct {
	MsgId     uint64
	From      common.Address
	BlockHash common.Hash

	ResultChan chan model.AbstractBlock `json:"-" rlp:"-"`
}

func (req *FetchBlockReqMsg) onResult(block model.AbstractBlock) {
	select {
	case req.ResultChan <- block:
	case <-time.After(100 * time.Millisecond):
		log.DLogger.Warn("can't send fetch resp to ResultChan, maybe already timeout")
	}
}

type FetchBlockRespMsg struct {
	MsgId uint64
	Block model.AbstractBlock
}
