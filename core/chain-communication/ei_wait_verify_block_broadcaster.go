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
//
//func makeEiWaitVerifyBlockBroadcaster(config *EiWaitVerifyBlockBroadcasterConfig) *EiWaitVerifyBlockBroadcaster {
//	service := &EiWaitVerifyBlockBroadcaster{
//		EiWaitVerifyBlockBroadcasterConfig: config,
//		handlers:                           map[uint64]func(msg p2p.Msg, p PmAbstractPeer) error{},
//
//		receivedHashCount:    "wait_v_block_hash_count",
//		receivedBlockCount:   "wait_v_block_block_count",
//		newTransportCount:    "new_transport_count",
//		deleteTransportCount: "delete_transport_count",
//	}
//	g_metrics.CreateCounter(service.receivedHashCount, "trace received wait v block hash", nil)
//	g_metrics.CreateCounter(service.receivedBlockCount, "trace received wait v block", nil)
//	g_metrics.CreateCounter(service.newTransportCount, "trace new transport count", nil)
//	g_metrics.CreateCounter(service.deleteTransportCount, "trace delete transport count", nil)
//
//	service.handlers[EiWaitVerifyBlockHashMsg] = service.onNewBlockHashMsg
//	service.handlers[EiWaitVerifyEstimatorMsg] = service.onEstimatorMsg
//	service.handlers[EiWaitVerifyBlockByBloomMsg] = service.onNewBloomBlock
//	return service
//}

//type EiWaitVerifyBlockBroadcasterConfig struct {
//	Pm       PeerManager
//	NodeConf NodeConf
//	TxPool   TxPool
//	Chain    Chain
//}
//
//// Estimator & InvBloom
//type EiWaitVerifyBlockBroadcaster struct {
//	*EiWaitVerifyBlockBroadcasterConfig
//
//	handlers map[uint64]func(msg p2p.Msg, p PmAbstractPeer) error
//
//	fetcher *WvEiBlockFetcher
//
//	// cache send new block msg
//	waitVerifyBroadcast sync.Map
//
//	receivedHashCount string
//	receivedBlockCount string
//
//	// new transport count
//	newTransportCount    string
//	deleteTransportCount string
//}
//
//func (broadcaster *EiWaitVerifyBlockBroadcaster) getTransport(p PmAbstractPeer) *eiBlockTransport {
//	// load transport
//	var transport *eiBlockTransport
//
//	if cache, ok := broadcaster.waitVerifyBroadcast.Load(p.ID()); ok {
//		transport = cache.(*eiBlockTransport)
//	} else {
//		transport = broadcaster.newTransport(p)
//		broadcaster.waitVerifyBroadcast.Store(p.ID(), transport)
//	}
//
//	return transport
//}
//
//// get peer without block
//func (broadcaster *EiWaitVerifyBlockBroadcaster) getPeersWithoutBlock(block model.AbstractBlock) []PmAbstractPeer {
//	// get peers
//	peers := broadcaster.Pm.GetPeers()
//	pbft_log.DLogger.Debug("EiWaitVerifyBlockBroadcaster_getPeersWithoutBlock", "block", block, "peers", peers)
//
//	var list []PmAbstractPeer
//
//	for _, p := range peers {
//		transport := broadcaster.getTransport(p)
//
//		if !transport.knownBlocks.Contains(block.Hash()) {
//			list = append(list, p)
//		}
//	}
//
//	return list
//}
//
//func (broadcaster *EiWaitVerifyBlockBroadcaster) newTransport(peer PmAbstractPeer) *eiBlockTransport {
//	transport := newEiBlockTransport(true, peer.ID(), peer.NodeName())
//	g_metrics.Add(broadcaster.newTransportCount, "", 1)
//
//	// start broadcast
//	go func() {
//		// if broadcast has err, break & remove peer from waitVerifyBroadcast
//		defer func() {
//			broadcaster.waitVerifyBroadcast.Delete(peer.ID())
//			g_metrics.Add(broadcaster.deleteTransportCount, "", 1)
//		}()
//
//		getPeer := func() PmAbstractPeer {
//			return broadcaster.Pm.GetPeer(peer.ID())
//		}
//
//		if err := transport.broadcast(getPeer); err != nil {
//			switch err {
//			case p2p.ErrShuttingDown:
//				log.DLogger.Warn("broadcast err is shutting down", "peer name", peer.NodeName(), "is running", peer.IsRunning())
//				broadcaster.Pm.RemovePeer(peer.ID())
//			case BroadcastTimeoutErr:
//			default:
//				log.DLogger.Error("ei wait verify block broadcast failed", "err", err, "peer name", peer.NodeName(), "is running", peer.IsRunning())
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
//func (broadcaster *EiWaitVerifyBlockBroadcaster) MsgHandlers() map[uint64]func(msg p2p.Msg, p PmAbstractPeer) error {
//	return broadcaster.handlers
//}
//
//// Broadcast new block, send new block hash msg
//func (broadcaster *EiWaitVerifyBlockBroadcaster) BroadcastBlock(block model.AbstractBlock) {
//	if broadcaster.NodeConf.GetNodeType() == chain_config.NodeTypeOfMineMaster {
//		broadcaster.fetcher.blockPool.addBlock(block)
//	}
//	//log.DLogger.Info("EiWaitVerifyBlockBroadcaster#BroadcastBlock", "broadcaster", broadcaster)
//	peers := broadcaster.getPeersWithoutBlock(block)
//	pbft_log.DLogger.Debug("EiWaitVerifyBlockBroadcaster_BroadcastBlock", "block", block, "peers", peers)
//
//	bHash := block.Hash()
//	txBloom := block.GetBlockTxsBloom()
//
//	msg := &eiBroadcastMsg{Height: block.Number(), BlockHash: bHash, TxBloom: txBloom}
//
//	var vPeers []PmAbstractPeer
//	var rPeers []PmAbstractPeer
//
//	for i := range peers {
//		if peers[i].NodeType() == chain_config.NodeTypeOfVerifier { //&& strings.Compare(peers[i].NodeName(),"default_v0")==0
//			vPeers = append(vPeers, peers[i])
//		} else {
//			if peers[i].NodeType() == chain_config.NodeTypeOfNormal {
//				rPeers = append(rPeers, peers[i])
//			}
//		}
//	}
//
//	// Send the msg to a subset of our peers
//	transferLen := int(math.Sqrt(float64(len(rPeers))))
//	if transferLen < minBroadcastPeers {
//		transferLen = minBroadcastPeers
//	}
//	if transferLen > len(rPeers) {
//		transferLen = len(rPeers)
//	}
//	transfer := rPeers[:transferLen]
//
//	//log.DLogger.Info("Miner broad cast block to", "Height", block.Number(), "v peer len", len(vPeers), "other peer len", len(transfer))
//	for i := range vPeers {
//		receiver := broadcaster.getTransport(vPeers[i])
//		//log.DLogger.Info("Miner broad cast block to","Peer",receiver.peerName,"Height",block.Number(), "block.hash", block.Hash().Hex())
//		receiver.asyncSendEiBroadcastMsg(msg)
//	}
//
//	for i := range transfer {
//		receiver := broadcaster.getTransport(transfer[i])
//		//log.DLogger.Info("Miner broad cast block to","Peer",receiver.peer.NodeName(),"Height",block.Number())
//		receiver.asyncSendEiBroadcastMsg(msg)
//	}
//
//}
//
//// remote peer handle new block hash msg
//func (broadcaster *EiWaitVerifyBlockBroadcaster) onNewBlockHashMsg(msg p2p.Msg, p PmAbstractPeer) error {
//	// decode msg
//	pbft_log.DLogger.Debug("EiWaitVerifyBlockBroadcaster#onNewBlockHashMsg  Receive on new block hash msg", "from", p.NodeName())
//	//pbft_log.DLogger.Debug("EiWaitVerifyBlockBroadcaster#onNewBlockHashMsg  Receive on new block hash msg", "from", p.NodeName())
//
//	var data eiBroadcastMsg
//	if err := msg.Decode(&data); err != nil {
//		return err
//	}
//
//	// check data block hash
//	if data.BlockHash.IsEmpty() {
//		log.DLogger.Warn("ei wait verify broadcast msg hash is nil", "p name", p.NodeName())
//		return nil
//	}
//
//	// check local chain has this block. if local chain has this block, no handle this msg
//	if broadcaster.hasBlock(data.BlockHash) {
//		return nil
//	}
//
//	// check block pool has this block
//	if broadcaster.fetcher.blockPool.getBlock(data.BlockHash) != nil {
//		return nil
//	}
//
//	// Cannot filter higher blocks (+x, x>1)
//
//	// check data tx bloom
//	if data.TxBloom == nil {
//		log.DLogger.Warn("ei wait verify broadcast msg tx bloom is nil")
//		return nil
//	}
//
//	g_metrics.Add(broadcaster.receivedHashCount, "", 1)
//
//	ts := broadcaster.getTransport(p)
//	ts.markHash(data.BlockHash)
//
//	pbft_log.DLogger.Debug("send EiWaitVerifyEstimatorMsg", "to", p.NodeName())
//	//pbft_log.DLogger.Debug("send EiWaitVerifyEstimatorMsg", "to", p.NodeName())
//
//	go func() {
//		_ = broadcaster.fetcher.Notify(p.NodeName(), p.ID(), data.BlockHash, data.Height, time.Now(), func() error {
//
//			txPool := broadcaster.TxPool
//			pending, queued := txPool.Stats()
//			log.DLogger.Info("cur tx pool", "pending", pending, "queued", queued)
//
//			//startAt := time.Now()
//
//			// peer tx pool constructs its Estimator locally
//			estimator := broadcaster.TxPool.GetTxsEstimator(data.TxBloom)
//
//			//log.DLogger.Info("EiWaitVerifyBlockBroadcaster GetTxsEstimator use time", "t", time.Now().Sub(startAt), "node name", p.NodeName())
//
//			return p.SendMsg(EiWaitVerifyEstimatorMsg, &eiEstimatorReq{BlockHash: data.BlockHash, Estimator: estimator})
//		})
//	}()
//
//	return nil
//}
//
//// check peer local chain has block by block hash
//func (broadcaster *EiWaitVerifyBlockBroadcaster) hasBlock(hash common.Hash) bool {
//	localBlock := broadcaster.Chain.GetBlockByHash(hash)
//
//	if localBlock != nil {
//		return true
//	}
//
//	return false
//}
//
//// receive request Estimator msg
//func (broadcaster *EiWaitVerifyBlockBroadcaster) onEstimatorMsg(msg p2p.Msg, p PmAbstractPeer) error {
//	//fmt.Println("===============EiWaitVerifyBlockBroadcaster--onEstimatorMsg==============")
//	pbft_log.DLogger.Debug("recieve onEstimatiorMsg", "from", p.NodeName())
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
//
//	// get block invBloom data
//	//startAt := time.Now()
//	data := broadcaster.getBlockInvBloomData(req.BlockHash, req.Estimator)
//	//log.DLogger.Info("EiWaitVerifyBlockBroadcaster getBlockInvBloomData use time", "t", time.Now().Sub(startAt), "node name", p.NodeName())
//	if data == nil {
//		log.DLogger.Error("can't get block inv bloom data, data is nil")
//		return nil
//	}
//
//	ts := broadcaster.getTransport(p)
//
//	ts.asyncEiBlockByBloomMsg(data)
//
//	return nil
//}
//
//func (broadcaster *EiWaitVerifyBlockBroadcaster) getBlockInvBloomData(bHash common.Hash, estimator *iblt.HybridEstimator) *model.BloomBlockData {
//	// get block
//	block := broadcaster.fetcher.blockPool.getBlock(bHash)
//	if block == nil {
//		log.DLogger.Error("local chain can't get block", "bHash", bHash.Hex())
//		return nil
//	}
//
//	// get block inv bloom data
//	data := block.GetEiBloomBlockData(estimator)
//
//	if data == nil {
//		return nil
//	}
//
//	return data
//}
//
//// receive target peer invBloom msg
//func (broadcaster *EiWaitVerifyBlockBroadcaster) onNewBloomBlock(msg p2p.Msg, p PmAbstractPeer) error {
//	//log.DLogger.Info("receive wait verify block", "remote", p.NodeName())
//	// decode msg
//	var temData bloomBlockDataRLP
//	if err := msg.Decode(&temData); err != nil {
//		log.DLogger.Error("ei wait verify on new bloom block decode msg failed", "err", err)
//		return err
//	}
//
//	// check local chain
//	if broadcaster.hasBlock(temData.Header.Hash()) {
//		return nil
//	}
//
//	if broadcaster.fetcher.blockPool.getBlock(temData.Header.Hash()) != nil {
//		return nil
//	}
//
//	g_metrics.Add(broadcaster.receivedBlockCount, "", 1)
//
//	broadcaster.fetcher.DoTask(p.ID(), &temData, time.Now())
//
//	return nil
//}
//
//func newWaitVerifyBlockPool() *waitVerifyBlockPool {
//	c, err := lru.New(30)
//	if err != nil {
//		panic(err)
//	}
//	return &waitVerifyBlockPool{blocks: c}
//}
//
//type waitVerifyBlockPool struct {
//	blocks *lru.Cache
//}
//
//func (pool *waitVerifyBlockPool) addBlock(block model.AbstractBlock) {
//	pool.blocks.Add(block.Hash(), block)
//}
//
//func (pool *waitVerifyBlockPool) getBlock(hash common.Hash) model.AbstractBlock {
//	if b, ok := pool.blocks.Get(hash); ok {
//		return b.(model.AbstractBlock)
//	}
//	return nil
//}
