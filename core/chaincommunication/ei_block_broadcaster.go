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

/*
 	first of all to reduce useless ei transmission：
		knownBlocks mark block

		step 1 :
					Alice  ------- new block hash + tx bloom -----> Bob

		step 2 :
					there is a tangled place here:  direct code or a new design fetch

					Bob	   ------- Estimator -------> Alice

		step 3 :
					Alice  ------- invBloom -----> Bob

		step 4 :
					Bob decode invBloom get block

	test case :

		1.  best case（two node trading pools are similar) 　calculate the size of three transmitted data

		2.  worst case（two node trading pools are completely different) 　Calculate the size of three transmitted data

*/

//func makeEiBlockBroadcaster(config *EiBlockBroadcasterConfig) *EiBlockBroadcaster {
//	service := &EiBlockBroadcaster{
//		EiBlockBroadcasterConfig: config,
//
//		handlers: map[uint64]func(msg p2p.Msg, p PmAbstractPeer) error{},
//
//		receivedHashCount:      "v_block_hash_count",
//		receivedBlockCount:     "v_block_block_count",
//		failedInsertBlockCount: "failed_insert_block_count",
//	}
//
//	g_metrics.CreateCounter(service.receivedHashCount, "trace received v block hash", nil)
//	g_metrics.CreateCounter(service.receivedBlockCount, "trace received v block", nil)
//	g_metrics.CreateCounter(service.failedInsertBlockCount, "trace failed insert block count", nil)
//
//	service.handlers[EiNewBlockHashMsg] = service.onNewBlockHashMsg
//	service.handlers[EiEstimatorMsg] = service.onEstimatorMsg
//	service.handlers[EiNewBlockByBloomMsg] = service.onNewBloomBlock
//	return service
//}

//type EiBlockBroadcasterConfig struct {
//	Pm     PeerManager
//	Chain  Chain
//	TxPool TxPool
//}
//
//// Estimator & InvBloom
//type EiBlockBroadcaster struct {
//	*EiBlockBroadcasterConfig
//
//	handlers map[uint64]func(msg p2p.Msg, p PmAbstractPeer) error
//
//	// cache send new block msg
//	vResultBroadcast sync.Map
//
//	fetcher *EiBlockFetcher
//
//	receivedHashCount string
//	receivedBlockCount string
//	failedInsertBlockCount string
//}
//
//func (broadcaster *EiBlockBroadcaster) getTransport(p PmAbstractPeer) *eiBlockTransport {
//	// load transport
//	var transport *eiBlockTransport
//
//	if cache, ok := broadcaster.vResultBroadcast.Load(p.ID()); ok {
//		transport = cache.(*eiBlockTransport)
//	} else {
//		transport = broadcaster.newTransport(p)
//		broadcaster.vResultBroadcast.Store(p.ID(), transport)
//	}
//
//	return transport
//}
//
//// get peer without block
//func (broadcaster *EiBlockBroadcaster) getPeersWithoutBlock(block model.AbstractBlock) []PmAbstractPeer {
//	// get peers
//	peers := broadcaster.Pm.GetPeers()
//
//	var list []PmAbstractPeer
//
//	for _, p := range peers {
//		transport := broadcaster.getTransport(p)
//
//		if !transport.knownBlocks.Contains(block.Hash()) {
//
//			//log.DLogger.Info("ei block broadcast get peer without block ", "node name", p.NodeName())
//
//			list = append(list, p)
//		}
//	}
//
//	return list
//}
//
//func (broadcaster *EiBlockBroadcaster) newTransport(peer PmAbstractPeer) *eiBlockTransport {
//	transport := newEiBlockTransport(false, peer.ID(), peer.NodeName())
//
//	// start broadcast
//	go func() {
//		// if broadcast has err, break & remove peer from waitVerifyBroadcast
//		defer func() {
//			broadcaster.vResultBroadcast.Delete(peer.ID())
//			log.DLogger.Warn("delete ei broadcaster transport", "peer name", peer.NodeName())
//		}()
//
//		getPeer := func() PmAbstractPeer {
//			return broadcaster.Pm.GetPeer(peer.ID())
//		}
//
//		if err := transport.broadcast(getPeer); err != nil {
//			switch err {
//			case p2p.ErrShuttingDown:
//				log.DLogger.Warn("ei verified block broadcast err is shutting down", "peer name", peer.NodeName(), "is running", peer.IsRunning())
//				broadcaster.Pm.RemovePeer(peer.ID())
//			case BroadcastTimeoutErr:
//			default:
//				log.DLogger.Error("ei verified block broadcast failed", "err", err, "peer name", peer.NodeName(), "is running", peer.IsRunning())
//			}
//			return
//		}
//
//	}()
//
//	return transport
//}
//
//// msg handler
//func (broadcaster *EiBlockBroadcaster) MsgHandlers() map[uint64]func(msg p2p.Msg, p PmAbstractPeer) error {
//	return broadcaster.handlers
//}
//
//// Broadcast new block, send new block hash msg
//func (broadcaster *EiBlockBroadcaster) BroadcastBlock(block model.AbstractBlock) {
//	peers := broadcaster.getPeersWithoutBlock(block)
//
//	bHash := block.Hash()
//	txBloom := block.GetBlockTxsBloom()
//
//	msg := &eiBroadcastMsg{Height: block.Number(), BlockHash: bHash, TxBloom: txBloom}
//
//	for i := range peers {
//		//if strings.Compare(peers[i].NodeName(),"default_m0")==0 {
//		receiver := broadcaster.getTransport(peers[i])
//		// send block hash
//		receiver.asyncSendEiBroadcastMsg(msg)
//		//}
//	}
//}
//
//// remote peer handle new block hash msg
//func (broadcaster *EiBlockBroadcaster) onNewBlockHashMsg(msg p2p.Msg, p PmAbstractPeer) error {
//	//fmt.Println("===============EiBlockBroadcaster--onNewBlockHashMsg==============")
//	//log.DLogger.Info("receive verified block", "from", p.NodeName())
//	// decode msg
//	var data eiBroadcastMsg
//	if err := msg.Decode(&data); err != nil {
//		return err
//	}
//	// check data block hash
//	if data.BlockHash.IsEmpty() {
//		log.DLogger.Warn("new block hash msg data block hash is nil")
//		return nil
//	}
//	pbft_log.DLogger.Info("EiBlockBroadcaster#onNewBlockHashMsg receive verified block", "from", p.NodeName(), "Height", data.Height, "blockHash", data.BlockHash)
//
//	// Can't filter higher blocks (+x, x>1), otherwise you can't set head, downloader can't get new block and cause chain to get stuck
//
//	// check local chain has this block. if local chain has this block, no handle this msg
//	if broadcaster.hasBlock(data.BlockHash) {
//		return nil
//	}
//
//	// check data tx bloom
//	if data.TxBloom == nil {
//		log.DLogger.Warn("new block hash msg data tx bloom is nil")
//		return nil
//	}
//
//	// set remote block height
//	p.SetHead(data.BlockHash, data.Height)
//
//	g_metrics.Add(broadcaster.receivedHashCount, "", 1)
//
//	ts := broadcaster.getTransport(p)
//
//	ts.markHash(data.BlockHash)
//
//	go func() {
//		_ = broadcaster.fetcher.Notify(p.NodeName(), p.ID(), data.BlockHash, data.Height, time.Now(), func() error {
//			// peer tx pool constructs its Estimator locally
//			estimator := broadcaster.TxPool.GetTxsEstimator(data.TxBloom)
//			pbft_log.DLogger.Info("EiBlockBroadcaster#onNewBlockHashMsg send EiEstimatorMsg", "BlockHash", data.BlockHash)
//			return p.SendMsg(EiEstimatorMsg, &eiEstimatorReq{BlockHash: data.BlockHash, Estimator: estimator})
//		})
//	}()
//
//	return nil
//}
//
//// check peer local chain has block by block hash
//func (broadcaster *EiBlockBroadcaster) hasBlock(hash common.Hash) bool {
//	localBlock := broadcaster.Chain.GetBlockByHash(hash) // hack there is no need to query the block here
//
//	if localBlock != nil {
//		return true
//	}
//
//	return false
//}
//
//// receive request Estimator msg
//func (broadcaster *EiBlockBroadcaster) onEstimatorMsg(msg p2p.Msg, p PmAbstractPeer) error {
//	//fmt.Println("===============EiBlockBroadcaster--onEstimatorMsg==============")
//	start := time.Now()
//	log.DLogger.Info("start get block inv", "time", start.String())
//	var req eiEstimatorReq
//	if err := msg.Decode(&req); err != nil {
//		return err
//	}
//
//	// get Estimator
//	estimator := req.Estimator
//	if estimator == nil {
//		log.DLogger.Error("Estimator is nil", "block hash", req.BlockHash.Hex(), "peer id", p.ID())
//		return nil
//	}
//	log.DLogger.Info("finish get block inv 1", "time", time.Now().Sub(start).String())
//	// get block invBloom data
//	data := broadcaster.getBlockInvBloomData(req.BlockHash, req.Estimator)
//
//	log.DLogger.Info("finish get block inv 2", "time", time.Now().Sub(start).String())
//
//	if data == nil {
//		log.DLogger.Error("con't get block inv bloom data, data is nil")
//		return nil
//	}
//
//	return p.SendMsg(EiNewBlockByBloomMsg, data)
//}
//
//func (broadcaster *EiBlockBroadcaster) getBlockInvBloomData(bHash common.Hash, estimator *iblt.HybridEstimator) *model.BloomBlockData {
//	// get block
//	block := broadcaster.Chain.GetBlockByHash(bHash)
//	//fmt.Println("===============EiBlockBroadcaster--onEstimatorMsg==============data",block.GetVerifications())
//	if block == nil {
//		log.DLogger.Error("EiBlockBroadcaster#getBlockInvBloomData local chain con't get block", "bHash", bHash.Hex())
//		return nil
//	}
//	vers := broadcaster.Chain.GetSeenCommit(block.Number())
//
//	// get block inv bloom data
//	data := block.GetEiBloomBlockData(estimator)
//	if data == nil {
//		return nil
//	}
//	data.CurVerification = vers
//	return data
//}
//
//// receive target peer invBloom msg
//func (broadcaster *EiBlockBroadcaster) onNewBloomBlock(msg p2p.Msg, p PmAbstractPeer) error {
//	startTime := time.Now()
//	log.DLogger.Info("start onNewBloomBlock", "time", startTime.String())
//
//	//fmt.Println("===============EiBlockBroadcaster--RecoverToBlock==============msg")
//	// decode msg
//	var temData bloomBlockDataRLP
//	if err := msg.Decode(&temData); err != nil {
//		log.DLogger.Error("ei verify on new bloom block decode msg failed", "err", err)
//		return err
//	}
//
//	// check local chain
//	if broadcaster.hasBlock(temData.Header.Hash()) {
//		return nil
//	}
//	// There is a problem after placing OnVerifyResultBlock. If the other party is a future block, it will not be inserted successfully, and it cannot be set to the correct header.
//	p.SetHead(temData.Header.Hash(), temData.Header.Number)
//
//	g_metrics.Add(broadcaster.receivedBlockCount, "", 1)
//	// here will call the save block
//	broadcaster.fetcher.DoTask(p.ID(), &temData, time.Now())
//	return nil
//}
