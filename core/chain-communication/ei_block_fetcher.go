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
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/bloom"
	"github.com/dipperin/dipperin-core/core/model"
	"time"
)

var (
	fetchTimeout = 5 * time.Second
)

const (
	hashLimit     = 256
	maxQueueDist  = 32
	arriveTimeout = 500 * time.Millisecond
	gatherSlack   = 100 * time.Millisecond
	blockLimit    = 64
)

type inject struct {
	peerID  string
	catchup *catchup
}

type chainHeightFunc func() model.AbstractBlock

type getBlockByHashFunc func(common.Hash) model.AbstractBlock

type saveBlockFunc func(block model.AbstractBlock, seenCommits []model.AbstractVerification) error

type dlTask struct {
	peerID      string
	catchupList []*catchupRlp
}

// estimator req
type estimatorReqFunc func() error

// get tx pool map
type getTxPoolMapFunc func() map[common.Hash]model.AbstractTransaction

// broadcast
type broadcastFunc func(block model.AbstractBlock)

// get hash
type hashMsg struct {
	nodeName     string
	hash         common.Hash
	number       uint64
	peerID       string
	time         time.Time
	txBloom      *iblt.Bloom
	estimatorReq estimatorReqFunc
}

type ibltTask struct {
	peerID string
	data   *bloomBlockDataRLP
	time   time.Time
}

//
//
//// Make Ei Block Fetcher
//func NewEiBlockFetcher(chainHeight chainHeightFunc, getBlock getBlockByHashFunc, saveBlock saveBlockFunc, getTxPoolMap getTxPoolMapFunc, broadcast broadcastFunc) *EiBlockFetcher {
//	f := &EiBlockFetcher{
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
//		getBlock:     getBlock,
//		saveBlock:    saveBlock,
//		getTxPoolMap: getTxPoolMap,
//		task:         make(chan chan *ibltTask),
//		filter:       make(chan chan *dlTask),
//		done:         make(chan common.Hash),
//		quit:         nil,
//		broadcast:    broadcast,
//
//		filterDownloaderGauge: "ei_block_fetcher_filterDownloaderGauge",
//		notifyCountGauge: "ei_block_fetcher_notifyCountGauge",
//		notifiedGauge: "ei_block_fetcher_notifiedGauge",
//		fetchingGauge: "ei_block_fetcher_fetchingGauge",
//		fetchedGauge: "ei_block_fetcher_fetchedGauge",
//		finishedGauge: "ei_block_fetcher_finishedGauge",
//		queueGauge: "ei_block_fetcher_queueGauge",
//		queuesGauge: "ei_block_fetcher_queuesGauge",
//		queuedGauge: "ei_block_fetcher_queuedGauge",
//
//	}
//
//	g_metrics.CreateGauge(f.filterDownloaderGauge, "", nil)
//	g_metrics.CreateGauge(f.notifyCountGauge, "", nil)
//	g_metrics.CreateGauge(f.notifiedGauge, "", nil)
//	g_metrics.CreateGauge(f.fetchingGauge, "", nil)
//	g_metrics.CreateGauge(f.fetchedGauge, "", nil)
//	g_metrics.CreateGauge(f.finishedGauge, "", nil)
//	g_metrics.CreateGauge(f.queueGauge, "", nil)
//	g_metrics.CreateGauge(f.queuesGauge, "", nil)
//	g_metrics.CreateGauge(f.queuedGauge, "", nil)
//
//	return f
//}
//
//type EiBlockFetcher struct {
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
//	chainHeight  chainHeightFunc
//	getBlock     getBlockByHashFunc
//	saveBlock    saveBlockFunc
//	getTxPoolMap getTxPoolMapFunc
//	broadcast    broadcastFunc
//
//	task chan chan *ibltTask
//
//	filter chan chan *dlTask
//
//	done chan common.Hash
//
//	quit chan struct{}
//
//	// metrics
//	filterDownloaderGauge string
//	notifyCountGauge string
//	notifiedGauge string
//	fetchingGauge string
//	fetchedGauge string
//	finishedGauge string
//	queueGauge string
//	queuesGauge string
//	queuedGauge string
//}
//
//func (f *EiBlockFetcher) Notify(name, pID string, hash common.Hash, number uint64, time time.Time, estimatorReq estimatorReqFunc) error {
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
//		log.Info("block fetcher terminated")
//		return nil
//	}
//}
//
//// Call this function after getting VerifyResult
//func (f *EiBlockFetcher) DoTask(peerID string, data *bloomBlockDataRLP, time time.Time) {
//	taskC := make(chan *ibltTask)
//
//	select {
//	case f.task <- taskC:
//	case <-f.quit:
//		return
//	}
//	pbft_log.Debug("Dotask,1")
//	select {
//	case taskC <- &ibltTask{peerID: peerID, data: data, time: time}:
//	case <-f.quit:
//		return
//	}
//}
//
//func (f *EiBlockFetcher) DoFilter(peerID string, list []*catchupRlp) []*catchupRlp {
//	filterC := make(chan *dlTask)
//
//	select {
//	case f.filter <- filterC:
//	case <-f.quit:
//		return nil
//	}
//
//	select {
//	case filterC <- &dlTask{peerID: peerID, catchupList: list}:
//	case <-f.quit:
//		return nil
//	}
//
//	select {
//	case task := <-filterC:
//		return task.catchupList
//	case <-f.quit:
//		return nil
//	}
//}
//
//func (f *EiBlockFetcher) Start() error {
//	if f.quit != nil {
//		return errors.New("already started")
//	}
//	f.quit = make(chan struct{})
//	go f.loop()
//	return nil
//}
//
//func (f *EiBlockFetcher) Stop() {
//	if f.quit == nil {
//		return
//	}
//	close(f.quit)
//	f.quit = nil
//}
//
//func (f *EiBlockFetcher) collectMetrics() {
//	g_metrics.Set(f.notifyCountGauge, "", float64(len(f.notifyCount)))
//	var count float64 = 0
//	for _, arr := range f.notified {
//		count += float64(len(arr))
//	}
//	g_metrics.Set(f.notifiedGauge, "", count)
//	g_metrics.Set(f.fetchingGauge, "", float64(len(f.fetching)))
//	count = 0
//	for _, arr := range f.fetched {
//		count += float64(len(arr))
//	}
//	g_metrics.Set(f.fetchedGauge, "", count)
//	g_metrics.Set(f.finishedGauge, "", float64(len(f.finished)))
//	g_metrics.Set(f.queueGauge, "", float64(f.queue.Size()))
//	g_metrics.Set(f.queuesGauge, "", float64(len(f.queues)))
//	g_metrics.Set(f.queuedGauge, "", float64(len(f.queued)))
//}
//
//func (f *EiBlockFetcher) loop() {
//	// Iterate the vr fetching until a quit is requested
//	fetchTimer := time.NewTimer(0)
//	tickHandler := func() {
//		f.collectMetrics()
//	}
//	metricsCollectTicker := g_timer.SetPeriodAndRun(tickHandler, 20 * time.Second)
//	defer func() {
//		fetchTimer.Stop()
//		g_timer.StopWork(metricsCollectTicker)
//	}()
//
//	for {
//		// Clean up any expired vr fetches
//		for hash, msg := range f.fetching {
//			if time.Since(msg.time) > fetchTimeout {
//				f.forgetHash(hash)
//			}
//		}
//
//		f.handleInsert()
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
//		case filterTask := <-f.filter:
//			var task *dlTask
//			select {
//			case task = <-filterTask:
//			case <-f.quit:
//				return
//			}
//
//			// handle filter
//			f.filterDownloader(task)
//
//			select {
//			case filterTask <- task:
//			case <-f.quit:
//				return
//			}
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
//func (f *EiBlockFetcher) handleInsert() {
//	height := f.chainHeight().Number()
//	for !f.queue.Empty() {
//
//		op := f.queue.PopItem().(*inject)
//		hash := op.catchup.Block.Hash()
//		number := op.catchup.Block.Number()
//		pbft_log.Debug("add block to pool", "block number", number, "height", height)
//		if number > height+1 {
//			f.queue.Push(op, -int64(number))
//			break
//		}
//
//		// Otherwise if fresh and still unknown, try and import
//		if f.getBlock(hash) != nil {
//			f.forgetBlock(hash)
//			continue
//		}
//
//		f.insert(op.peerID, op.catchup)
//	}
//}
//
//func (f *EiBlockFetcher) handleNotify(notification *hashMsg, fetchTimer *time.Timer) (needBreak bool) {
//	// A vr was announced, make sure the peer isn't DOSing us
//	count := f.notifyCount[notification.peerID] + 1
//
//	if count > hashLimit {
//		log.Error("Peer exceeded outstanding announces", "peer", notification.peerID, "limit", hashLimit)
//		needBreak = true
//		return needBreak
//	}
//
//	// If we have a valid block number, check that it's potentially useful
//	if notification.number > 0 {
//		if dist := int64(notification.number) - int64(f.chainHeight().Number()); dist > maxQueueDist {
//			log.Debug("Peer discarded announcement", "peer", notification.peerID, "number", notification.number, "hash", notification.hash, "distance", dist)
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
//	//log.Info("----------get block new hash", "hash", notification.hash.Hex(), "height", notification.number,
//	//	"peer name", notification.nodeName)
//
//	if len(f.notified) == 1 {
//		f.rescheduleFetch(fetchTimer)
//	}
//
//	return needBreak
//}
//
//func (f *EiBlockFetcher) handleFetching(fetchTimer *time.Timer) {
//	// At least one vr's timer ran out, check for needing retrieval
//	request := make(map[string][]common.Hash)
//
//	for hash, msgList := range f.notified {
//		if time.Since(msgList[0].time) > arriveTimeout-gatherSlack {
//			msg := msgList[rand.Intn(len(msgList))]
//			f.forgetHash(hash)
//
//			// If the vr still didn't arrive, queue for fetching
//			if f.getBlock(hash) == nil {
//				request[msg.peerID] = append(request[msg.peerID], hash)
//				// add fetching
//				f.fetching[hash] = msg
//			}
//		}
//	}
//
//	// Send out all get vr requests
//	for peer, hashes := range request {
//		log.Debug("Fetching scheduled vrs", "peer", peer, "hashes", hashes)
//
//		req, hashes := f.fetching[hashes[0]].estimatorReq, hashes
//
//		go func() {
//			for _, hash := range hashes {
//				log.Debug("start fetching block vr", "hash", hash)
//				if err := req(); err != nil {
//					log.Error("ei fetcher handle fetcher estimator req failed", "err", err)
//				}
//			}
//		}()
//	}
//
//	// Schedule the next fetch if msg are still pending
//	f.rescheduleFetch(fetchTimer)
//}
//
//func (f *EiBlockFetcher) handleIBLTTask(task *ibltTask) {
//	pbft_log.Debug("Handle vr task")
//	data := task.data
//
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
//		if f.getBlock(data.Header.Hash()) == nil {
//			catchup := f.bloomData2catchUp(data)
//
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
//func (f *EiBlockFetcher) bloomData2catchUp(temData *bloomBlockDataRLP) *catchup {
//	preCommits := make([]model.AbstractVerification, len(temData.PreVerification))
//	util.InterfaceSliceCopy(preCommits, temData.PreVerification)
//
//	curCommits := make([]model.AbstractVerification, len(temData.CurVerification))
//	util.InterfaceSliceCopy(curCommits, temData.CurVerification)
//
//	data := &model.BloomBlockData{
//		Header:          temData.Header,
//		BloomRLP:        temData.BloomRLP,
//		PreVerification: preCommits,
//		CurVerification: curCommits,
//		Interlinks:      temData.Interlinks,
//	}
//
//	block, err := data.EiRecoverToBlock(f.getTxPoolMap())
//
//	if err != nil {
//		log.Error("ei fetcher ei recover to block failed", "err", err)
//		return nil
//	}
//
//	block.SetInterLinks(data.Interlinks)
//
//	return &catchup{Block: block, SeenCommit: curCommits}
//}
//
//func (f *EiBlockFetcher) filterDownloader(task *dlTask) {
//	g_metrics.Add(f.filterDownloaderGauge, "", 1)
//	defer g_metrics.Sub(f.filterDownloaderGauge, "", 1)
//
//	var pendingList []*catchup
//	for i := 0; i < len(task.catchupList); i++ {
//		// Match up a body to any possible completion request
//		matched := false
//		for hash := range f.notified {
//			if f.queued[hash] == nil {
//				if hash.IsEqual(task.catchupList[i].Block.Hash()) {
//					matched = true
//					if f.getBlock(hash) == nil {
//						tmpSeen := make([]model.AbstractVerification, len(task.catchupList[i].SeenCommit))
//						util.InterfaceSliceCopy(tmpSeen, task.catchupList[i].SeenCommit)
//
//						cup := &catchup{Block: task.catchupList[i].Block, SeenCommit: tmpSeen}
//
//						pendingList = append(pendingList, cup)
//
//					} else {
//						f.forgetHash(hash)
//					}
//				}
//			}
//		}
//
//		if matched {
//			task.catchupList = append(task.catchupList[:i], task.catchupList[i+1:]...)
//			i--
//			continue
//		}
//	}
//
//	for i := range pendingList {
//		f.enqueue(task.peerID, pendingList[i])
//	}
//}
//
//func (f *EiBlockFetcher) enqueue(peerID string, catchup *catchup) {
//	// ensure the peer isn't DOSing us
//	count := f.queues[peerID] + 1
//
//	if count > blockLimit {
//		log.Debug("Discarded propagated block, exceeded allowance", "peer", peerID, "number", catchup.Block.Number(), "hash", catchup.Block.Hash(), "limit", blockLimit)
//		f.forgetHash(catchup.Block.Hash())
//		return
//	}
//
//	// Discard any past or too distant blocks
//	if dist := int64(catchup.Block.Number()) - int64(f.chainHeight().Number()); dist > maxQueueDist {
//		log.Debug("Discarded propagated block, too far away", "peer", peerID, "number", catchup.Block.Number(), "hash", catchup.Block.Hash(), "distance", dist)
//		f.forgetHash(catchup.Block.Hash())
//		return
//	}
//
//	// Schedule the vr for future importing
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
//		log.Debug("Queued propagated block", "peer", peerID, "number", catchup.Block.Number(), "hash", catchup.Block.Hash(), "queued", f.queue.Size())
//	}
//
//}
//
//func (f *EiBlockFetcher) insert(peerID string, catchup *catchup) {
//	block := catchup.Block
//	pbft_log.Debug("Insert a block", "block", block.Number(), "height", f.chainHeight().Number())
//	go func() {
//		defer func() {
//			f.done <- block.Hash()
//		}()
//
//		parent := f.getBlock(block.PreHash())
//		if parent == nil {
//			pbft_log.Debug("Insert a block", "block", block.Number(), "height", f.chainHeight().Number(), "err", "Unknown parent of propagated block")
//			log.Error("Unknown parent of propagated block", "peer", peerID, "number", block.Number(), "hash", block.Hash(), "parent", block.PreHash())
//			return
//		}
//
//		if err := f.saveBlock(catchup.Block, catchup.SeenCommit); err != nil {
//			pbft_log.Debug("Save a block", "block", block.Number(), "height", f.chainHeight().Number(), "err", err)
//			log.Error("Propagated block import failed", "peer", peerID, "number", block.Number(), "hash", block.Hash(), "err", err)
//			return
//		}
//		pbft_log.Debug("Saved a block", "block", block.Number(), "height", f.chainHeight().Number())
//		log.Info("fetcher save block vr", "hash", block.Hash(), "number", block.Number())
//
//		go f.broadcast(catchup.Block)
//
//	}()
//}
//
//// rescheduleFetch resets the specified fetch timer to the next announce timeout.
//func (f *EiBlockFetcher) rescheduleFetch(fetch *time.Timer) {
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
//func (f *EiBlockFetcher) forgetHash(hash common.Hash) {
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
//func (f *EiBlockFetcher) forgetBlock(hash common.Hash) {
//	if insert := f.queued[hash]; insert != nil {
//		f.queues[insert.peerID]--
//		if f.queues[insert.peerID] == 0 {
//			delete(f.queues, insert.peerID)
//		}
//		delete(f.queued, hash)
//	}
//}
