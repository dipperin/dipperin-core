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
	"github.com/dipperin/dipperin-core/core/model"
)

// todo do you need to broadcast again？？？
type wvBroadcastFunc func(block model.AbstractBlock)

//
//// Make Wv Ei Block Fetcher
//func NewWvEiBlockFetcher(config *WvEiBlockFetcherConfig, chainHeight chainHeightFunc, getTxPoolMap getTxPoolMapFunc, broadcast wvBroadcastFunc) *WvEiBlockFetcher {
//	return &WvEiBlockFetcher{
//		WvEiBlockFetcherConfig: config,
//
//		notify:       make(chan *hashMsg),
//		notifyCount:  make(map[string]int),
//		notified:     make(map[common.Hash][]*hashMsg),
//		fetching:     make(map[common.Hash]*hashMsg),
//		fetched:      make(map[common.Hash][]*hashMsg),
//		finished:     make(map[common.Hash]*hashMsg),
//		queue:        prque.New(nil),
//		queues:       make(map[string]int),
//		queued:       make(map[common.Hash]*inject),
//		chainHeight:  chainHeight,
//		getTxPoolMap: getTxPoolMap,
//		task:         make(chan chan *ibltTask),
//		done:         make(chan common.Hash),
//		quit:         nil,
//		broadcast:    broadcast,
//		blockPool:    newWaitVerifyBlockPool(),
//	}
//}
//
//type WvEiBlockFetcherConfig struct {
//	PbftNode PbftNode
//}
//
//type WvEiBlockFetcher struct {
//	*WvEiBlockFetcherConfig
//
//	notify chan *hashMsg
//
//	// key --> peer id, value notify count
//	notifyCount map[string]int
//	notified    map[common.Hash][]*hashMsg
//	fetching    map[common.Hash]*hashMsg
//	fetched     map[common.Hash][]*hashMsg
//	finished    map[common.Hash]*hashMsg
//
//	queue  *prque.Prque
//	queues map[string]int
//	queued map[common.Hash]*inject
//
//	chainHeight chainHeightFunc
//
//	task         chan chan *ibltTask
//	getTxPoolMap getTxPoolMapFunc
//	broadcast    wvBroadcastFunc
//
//	done chan common.Hash
//	quit chan struct{}
//
//	blockPool *waitVerifyBlockPool
//}
//
//// Receive the hash, determine the function to call this function after calling this block
//func (f *WvEiBlockFetcher) Notify(name, pID string, hash common.Hash, number uint64, time time.Time, estimatorReq estimatorReqFunc) error {
//	msg := &hashMsg{
//		nodeName:     name,
//		hash:         hash,
//		number:       number,
//		time:         time,
//		peerID:       pID,
//		estimatorReq: estimatorReq,
//	}
//
//	select {
//	case f.notify <- msg:
//		return nil
//	case <-f.quit:
//		log.DLogger.Info("wv ei block fetcher terminated")
//		return nil
//	}
//}
//
//func (f *WvEiBlockFetcher) DoTask(peerID string, data *bloomBlockDataRLP, time time.Time) {
//	taskC := make(chan *ibltTask)
//
//	select {
//	case f.task <- taskC:
//	case <-f.quit:
//		return
//	}
//	pbft_log.DLogger.Debug("Dotask,1")
//	select {
//	case taskC <- &ibltTask{peerID: peerID, data: data, time: time}:
//	case <-f.quit:
//		return
//	}
//}
//
//func (f *WvEiBlockFetcher) Start() error {
//	if f.quit != nil {
//		return errors.New("already started")
//	}
//	f.quit = make(chan struct{})
//	go f.loop()
//	return nil
//}
//
//func (f *WvEiBlockFetcher) Stop() {
//	if f.quit == nil {
//		return
//	}
//	close(f.quit)
//	f.quit = nil
//}
//
//func (f *WvEiBlockFetcher) loop() {
//	// Iterate the vr fetching until a quit is requested
//	fetchTimer := time.NewTimer(0)
//
//	for {
//		// Clean up any expired vr fetches
//		for hash, msg := range f.fetching {
//			if time.Since(msg.time) > fetchTimeout {
//				f.forgetHash(hash)
//			}
//		}
//
//		f.handlePush()
//
//		select {
//		case <-f.quit:
//			// Fetcher terminating, abort all operations
//			return
//
//		case notification := <-f.notify:
//			// handle notification
//			if needBreak := f.handleNotify(notification, fetchTimer); needBreak {
//				break
//			}
//
//		case <-fetchTimer.C:
//			// handle fetching
//			f.handleFetching(fetchTimer)
//
//		case cmTask := <-f.task:
//			var task *ibltTask
//			select {
//			case task = <-cmTask:
//			case <-f.quit:
//				return
//			}
//
//			// handle
//			f.handleIBLTTask(task)
//
//		case hash := <-f.done:
//			// A pending import finished, remove all traces of the notification
//			f.forgetHash(hash)
//			f.forgetBlock(hash)
//
//		}
//	}
//}
//
//func (f *WvEiBlockFetcher) handlePush() {
//	height := f.chainHeight().Number()
//	pbft_log.DLogger.Debug("fetch chain height number", "height", height)
//	for !f.queue.Empty() {
//
//		op := f.queue.PopItem().(*inject)
//		hash := op.catchup.Block.Hash()
//		number := op.catchup.Block.Number()
//		pbft_log.DLogger.Debug("to add block to block pool", "block number", number, "height", height)
//		if number > height+1 {
//			f.queue.Push(op, -int64(number))
//			break
//		}
//
//		if number < height {
//			f.forgetHash(hash)
//			f.forgetBlock(hash)
//			continue
//		}
//
//		// Otherwise if fresh and still unknown, try and import
//		if f.blockPool.getBlock(hash) != nil {
//			f.forgetHash(hash)
//			f.forgetBlock(hash)
//			continue
//		}
//
//		f.push(op.peerID, op.catchup)
//	}
//}
//
//func (f *WvEiBlockFetcher) handleNotify(notification *hashMsg, fetchTimer *time.Timer) (needBreak bool) {
//	// A vr was announced, make sure the peer isn't DOSing us
//	count := f.notifyCount[notification.peerID] + 1
//
//	if count > hashLimit {
//		log.DLogger.Error("wv ei fetcher Peer exceeded outstanding announces", "peer", notification.peerID, "limit", hashLimit)
//		needBreak = true
//		return needBreak
//	}
//
//	// If we have a valid block number, check that it's potentially useful
//	if notification.number > 0 {
//		chainHeight := f.chainHeight().Number()
//
//		if notification.number < chainHeight {
//			log.DLogger.Debug("notification number < chain height")
//			needBreak = true
//			return needBreak
//		}
//
//		if dist := int64(notification.number) - int64(f.chainHeight().Number()); dist > maxQueueDist {
//			log.DLogger.Debug("Peer discarded announcement", "peer", notification.peerID, "number", notification.number, "hash", notification.hash, "distance", dist)
//			needBreak = true
//			return needBreak
//		}
//	} else {
//		needBreak = true
//		return needBreak
//	}
//
//	// All is well, schedule the announce if vr's not yet downloading
//	if _, ok := f.fetching[notification.hash]; ok {
//		needBreak = true
//		return needBreak
//	}
//
//	if _, ok := f.finished[notification.hash]; ok {
//		needBreak = true
//		return needBreak
//	}
//
//	f.notifyCount[notification.peerID] = count
//	f.notified[notification.hash] = append(f.notified[notification.hash], notification)
//
//	if len(f.notified) == 1 {
//		f.rescheduleFetch(fetchTimer)
//	}
//
//	return needBreak
//}
//
//func (f *WvEiBlockFetcher) handleFetching(fetchTimer *time.Timer) {
//	// At least one vr's timer ran out, check for needing retrieval
//	request := make(map[string][]common.Hash)
//
//	for hash, msgList := range f.notified {
//		if time.Since(msgList[0].time) > arriveTimeout-gatherSlack {
//			msg := msgList[rand.Intn(len(msgList))]
//			f.forgetHash(hash)
//
//			// If the vr still didn't arrive, queue for fetching
//			if f.blockPool.getBlock(hash) == nil {
//				request[msg.peerID] = append(request[msg.peerID], hash)
//				// add fetching
//				f.fetching[hash] = msg
//			}
//		}
//	}
//
//	// Send out all get vr requests
//	for peer, hashes := range request {
//		log.DLogger.Debug("Fetching scheduled vrs", "peer", peer, "hashes", hashes)
//
//		req, hashes := f.fetching[hashes[0]].estimatorReq, hashes
//
//		go func() {
//			for _, hash := range hashes {
//				log.DLogger.Debug("start fetching wait verify block", "hash", hash)
//				if err := req(); err != nil {
//					log.DLogger.Error("ei fetcher handle fetcher estimator req failed", "err", err)
//				}
//			}
//		}()
//	}
//
//	// Schedule the next fetch if msg are still pending
//	f.rescheduleFetch(fetchTimer)
//}
//
//func (f *WvEiBlockFetcher) handleIBLTTask(task *ibltTask) {
//	log.DLogger.Debug("handle wv ei iblt task")
//
//	data := task.data
//	//pbft_log.DLogger.Debug("WvEiBlockFetcher#handleIBLTTask handle task.data","task.data", data)
//	// Filter fetcher-requested headers from other synchronisation algorithms
//	if msg := f.fetching[data.Header.Hash()]; msg != nil && msg.peerID == task.peerID && f.fetched[data.Header.Hash()] == nil && f.finished[data.Header.Hash()] == nil {
//		// If the delivered header does not match the promised number, drop the peer
//
//		if data.Header.Number != msg.number {
//			// todo peer drop
//			f.forgetHash(data.Header.Hash())
//		}
//
//		// Only keep if not imported by other means
//		if f.blockPool.getBlock(data.Header.Hash()) == nil {
//			catchup := f.bloomData2catchUp(data)
//			//pbft_log.DLogger.Debug("WvEiBlockFetcher#handleIBLTTask catchup","catchup", catchup)
//			if catchup != nil {
//				msg.time = task.time
//				f.finished[data.Header.Hash()] = msg
//
//				// put queue
//				f.enqueue(msg.peerID, catchup)
//			}
//
//		} else {
//			f.forgetHash(data.Header.Hash())
//		}
//	}
//
//}
//
//func (f *WvEiBlockFetcher) bloomData2catchUp(temData *bloomBlockDataRLP) *catchup {
//	preCommits := make([]model.AbstractVerification, len(temData.PreVerification))
//	util.InterfaceSliceCopy(preCommits, temData.PreVerification)
//
//	data := &model.BloomBlockData{
//		Header:          temData.Header,
//		BloomRLP:        temData.BloomRLP,
//		PreVerification: preCommits,
//		Interlinks:      temData.Interlinks,
//	}
//
//	block, err := data.EiRecoverToBlock(f.getTxPoolMap())
//
//	if err != nil {
//		log.DLogger.Error("ei fetcher ei recover to block failed", "err", err)
//		return nil
//	}
//
//	block.SetInterLinks(data.Interlinks)
//
//	return &catchup{Block: block, SeenCommit: nil}
//}
//
//func (f *WvEiBlockFetcher) enqueue(peerID string, catchup *catchup) {
//	// ensure the peer isn't DOSing us
//	count := f.queues[peerID] + 1
//
//	if count > blockLimit {
//		log.DLogger.Debug("Discarded propagated block, exceeded allowance", "peer", peerID, "number", catchup.Block.Number(), "hash", catchup.Block.Hash(), "limit", blockLimit)
//		f.forgetHash(catchup.Block.Hash())
//		return
//	}
//
//	if f.chainHeight().Number() > catchup.Block.Number() {
//		f.forgetHash(catchup.Block.Hash())
//		return
//	}
//
//	// Discard any past or too distant blocks
//	if dist := int64(catchup.Block.Number()) - int64(f.chainHeight().Number()); dist > maxQueueDist {
//		log.DLogger.Debug("Discarded propagated block, too far away", "peer", peerID, "number", catchup.Block.Number(), "hash", catchup.Block.Hash(), "distance", dist)
//		f.forgetHash(catchup.Block.Hash())
//		return
//	}
//	//pbft_log.DLogger.Debug("WvEiBlockFetcher#enqueue catchup is going to enqueue", "catchup", catchup)
//
//	// Schedule the wv for future importing
//	if _, ok := f.queued[catchup.Block.Hash()]; !ok {
//		op := &inject{
//			peerID:  peerID,
//			catchup: catchup,
//		}
//
//		f.queues[peerID] = count
//		f.queued[catchup.Block.Hash()] = op
//		f.queue.Push(op, -int64(catchup.Block.Number()))
//
//		log.DLogger.Debug("Queued propagated block", "peer", peerID, "number", catchup.Block.Number(), "hash", catchup.Block.Hash(), "queued", f.queue.Size())
//	}
//
//}
//
//func (f *WvEiBlockFetcher) push(peerID string, catchup *catchup) {
//	block := catchup.Block
//	pbft_log.DLogger.Debug("add block to block pool", "block", block.Number(), "height", f.chainHeight().Number())
//	go func() {
//		defer func() {
//			f.done <- block.Hash()
//		}()
//
//		f.blockPool.addBlock(block)
//
//		// determine if you are a verifier
//		pbftNode := f.PbftNode
//		if !reflect.ValueOf(pbftNode).IsNil() {
//			pbftNode.OnNewWaitVerifyBlock(block, peerID)
//			return
//		}
//
//		go f.broadcast(catchup.Block)
//	}()
//}
//
//// rescheduleFetch resets the specified fetch timer to the next announce timeout.
//func (f *WvEiBlockFetcher) rescheduleFetch(fetch *time.Timer) {
//	// Short circuit if no blocks are announced
//	if len(f.notified) == 0 {
//		return
//	}
//	// Otherwise find the earliest expiring announcement
//	earliest := time.Now()
//	for _, msgList := range f.notified {
//		if earliest.After(msgList[0].time) {
//			earliest = msgList[0].time
//		}
//	}
//	fetch.Reset(arriveTimeout - time.Since(earliest))
//}
//
//func (f *WvEiBlockFetcher) forgetHash(hash common.Hash) {
//	// Remove all pending announces and decrement DOS counters
//	for _, msg := range f.notified[hash] {
//		f.notifyCount[msg.peerID]--
//		if f.notifyCount[msg.peerID] == 0 {
//			delete(f.notifyCount, msg.peerID)
//		}
//	}
//
//	delete(f.notified, hash)
//
//	// Remove any pending fetches and decrement the DOS counters
//	if msg := f.fetching[hash]; msg != nil {
//		f.notifyCount[msg.peerID]--
//		if f.notifyCount[msg.peerID] == 0 {
//			delete(f.notifyCount, msg.peerID)
//		}
//		delete(f.fetching, hash)
//	}
//
//	// Remove any pending completion requests and decrement the DOS counters
//	for _, msg := range f.fetched[hash] {
//		f.notifyCount[msg.peerID]--
//		if f.notifyCount[msg.peerID] == 0 {
//			delete(f.notifyCount, msg.peerID)
//		}
//	}
//	delete(f.fetched, hash)
//
//	// Remove any pending completions and decrement the DOS counters
//	if msg := f.finished[hash]; msg != nil {
//		f.notifyCount[msg.peerID]--
//		if f.notifyCount[msg.peerID] == 0 {
//			delete(f.notifyCount, msg.peerID)
//		}
//		delete(f.finished, hash)
//	}
//}
//
//// forgetBlock removes all traces of a queued block from the fetcher's internal
//// state.
//func (f *WvEiBlockFetcher) forgetBlock(hash common.Hash) {
//	if insert := f.queued[hash]; insert != nil {
//		f.queues[insert.peerID]--
//		if f.queues[insert.peerID] == 0 {
//			delete(f.queues, insert.peerID)
//		}
//		delete(f.queued, hash)
//	}
//}
