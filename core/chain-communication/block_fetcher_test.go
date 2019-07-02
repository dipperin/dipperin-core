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
	"math/big"
	"testing"
	"time"

	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/prque"
	model2 "github.com/dipperin/dipperin-core/core/csbft/model"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/tests/factory"
	"github.com/stretchr/testify/assert"
)

func getVr(hash common.Hash) error {
	// mock get vr p2p send request
	return nil
}

func chainHeight() model.AbstractBlock {
	// mock get current block

	block := factory.CreateBlock(2)

	return block
}

func getBlockByHashReturnNil(hash common.Hash) model.AbstractBlock {
	return nil
}

func saveBlock(block model.AbstractBlock, seenCommits []model.AbstractVerification) error {
	return nil
}

func broadcast(b *model2.VerifyResult) {
}

func getFetcher() *BlockFetcher {
	return NewBlockFetcher(chainHeight, getBlockByHashReturnNil, saveBlock, broadcast)
}

func TestNewBlockFetcher(t *testing.T) {
	fetcher := getFetcher()

	assert.NotNil(t, fetcher)
}

func TestBlockFetcher_Start(t *testing.T) {
	fetchTimeout = time.Millisecond
	defer func() { fetchTimeout = 5 * time.Second }()

	fetcher := getFetcher()
	assert.NotNil(t, fetcher)

	_ = fetcher.Start()

	time.Sleep(2 * time.Millisecond)
}

func TestBlockFetcher_Stop(t *testing.T) {
	fetchTimeout = time.Millisecond
	defer func() { fetchTimeout = 5 * time.Second }()

	fetcher := getFetcher()
	assert.NotNil(t, fetcher)

	go func() {
		_ = fetcher.Start()
	}()

	time.Sleep(2 * time.Millisecond)

	fetcher.Stop()
}

func TestBlockFetcher_Notify_CheckMapInit(t *testing.T) {
	// init test env, start fetcher loop
	fetcher := getFetcher()
	assert.NotNil(t, fetcher)

	go func() {
		_ = fetcher.Start()
	}()

	time.Sleep(200 * time.Millisecond)

	// check map init
	assert.Equal(t, 0, len(fetcher.notifyCount))
	assert.Equal(t, 0, len(fetcher.notified))
	assert.Equal(t, 0, len(fetcher.fetching))
	assert.Equal(t, 0, len(fetcher.fetched))
	assert.Equal(t, 0, len(fetcher.finished))
}

func TestBlockFetcher_Notify_AddNotify(t *testing.T) {
	// init test env, start fetcher loop
	fetcher := getFetcher()
	assert.NotNil(t, fetcher)

	go func() {
		_ = fetcher.Start()
	}()

	time.Sleep(200 * time.Millisecond)

	// use Notify
	_ = fetcher.Notify("001", common.HexToHash("aaa"), 3, time.Now(), getVr)

	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, 1, len(fetcher.notifyCount))
	assert.Equal(t, 1, len(fetcher.notified))
	assert.Equal(t, 0, len(fetcher.fetching))
	assert.Equal(t, 0, len(fetcher.fetched))
	assert.Equal(t, 0, len(fetcher.finished))
}

func TestBlockFetcher_Notify_handleFetching(t *testing.T) {
	var getVrReqCount int

	getVrFn := func(hash common.Hash) error {
		getVrReqCount++
		return nil
	}

	fetcher := NewBlockFetcher(chainHeight, getBlockByHashReturnNil, saveBlock, broadcast)
	assert.NotNil(t, fetcher)

	go func() {
		_ = fetcher.Start()
	}()

	time.Sleep(200 * time.Millisecond)

	// use Notify
	_ = fetcher.Notify("001", common.HexToHash("aaa"), 3, time.Now(), getVrFn)

	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, 1, len(fetcher.notifyCount))
	assert.Equal(t, 1, len(fetcher.notified))

	time.Sleep(arriveTimeout)

	// go handle fetcher

	time.Sleep(arriveTimeout)

	assert.Equal(t, 0, len(fetcher.notifyCount))
	assert.Equal(t, 0, len(fetcher.notified))
	assert.Equal(t, 1, len(fetcher.fetching))

	assert.Equal(t, 1, getVrReqCount)
}

func TestBlockFetcher_DoTask(t *testing.T) {
	var getVrReqCount int

	getVrFn := func(hash common.Hash) error {
		getVrReqCount++
		return nil
	}

	block := factory.CreateBlock(3)

	getBlock := func(hash common.Hash) model.AbstractBlock {

		if hash.IsEqual(block.PreHash()) {
			return factory.CreateBlock(2)
		}

		return nil
	}

	var saveBlockCount int

	saveBlock := func(block model.AbstractBlock, seenCommits []model.AbstractVerification) error {
		saveBlockCount++
		return nil
	}

	fetcher := NewBlockFetcher(chainHeight, getBlock, saveBlock, broadcast)

	go func() {
		_ = fetcher.Start()
	}()

	time.Sleep(200 * time.Millisecond)

	// use Notify
	_ = fetcher.Notify("001", block.Hash(), 3, time.Now(), getVrFn)

	time.Sleep(arriveTimeout)

	// go handle fetcher

	time.Sleep(arriveTimeout)

	assert.Equal(t, 1, getVrReqCount)

	vr := &model2.VerifyResult{Block: block, SeenCommits: nil}

	fetcher.DoTask("001", vr, time.Now())

	time.Sleep(1 * time.Second)

	assert.Equal(t, 1, saveBlockCount)

}

func TestArr(t *testing.T) {
	var list []int

	for i := 0; i < 10; i++ {
		list = append(list, i)
	}

	for i := 0; i < len(list); i++ {
		matched := false

		if list[i] == 0 {
			matched = true
		}

		if matched {
			list = append(list[:i], list[i+1:]...)
			i--
			continue
		}
	}

	for i := range list {
		println(list[i])
	}
}

func TestBlockFetcher_DoFilter(t *testing.T) {
	// init test env, start fetcher loop
	fetcher := getFetcher()
	assert.NotNil(t, fetcher)

	go func() {
		_ = fetcher.Start()
	}()

	time.Sleep(200 * time.Millisecond)

	fetcher.DoFilter("123", []*catchupRlp{})
}

func TestBlockFetcher_loop(t *testing.T) {
	type fields struct {
		notify           chan *vrMsg
		notifyCount      map[string]int
		notified         map[common.Hash][]*vrMsg
		fetching         map[common.Hash]*vrMsg
		fetched          map[common.Hash][]*vrMsg
		finished         map[common.Hash]*vrMsg
		queue            *prque.Prque
		queues           map[string]int
		queued           map[common.Hash]*inject
		chainHeight      chainHeightFunc
		getBlock         getBlockByHashFunc
		saveBlock        saveBlockFunc
		blockBroadcaster blockBroadcasterFunc
		task             chan chan *vrTask
		filter           chan chan *dlTask
		done             chan common.Hash
		quit             chan struct{}
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &BlockFetcher{
				notify:           tt.fields.notify,
				notifyCount:      tt.fields.notifyCount,
				notified:         tt.fields.notified,
				fetching:         tt.fields.fetching,
				fetched:          tt.fields.fetched,
				finished:         tt.fields.finished,
				queue:            tt.fields.queue,
				queues:           tt.fields.queues,
				queued:           tt.fields.queued,
				chainHeight:      tt.fields.chainHeight,
				getBlock:         tt.fields.getBlock,
				saveBlock:        tt.fields.saveBlock,
				blockBroadcaster: tt.fields.blockBroadcaster,
				task:             tt.fields.task,
				filter:           tt.fields.filter,
				done:             tt.fields.done,
				quit:             tt.fields.quit,
			}
			f.loop()
		})
	}
}

func TestBlockFetcher_handleInsert(t *testing.T) {
	type fields struct {
		notify           chan *vrMsg
		notifyCount      map[string]int
		notified         map[common.Hash][]*vrMsg
		fetching         map[common.Hash]*vrMsg
		fetched          map[common.Hash][]*vrMsg
		finished         map[common.Hash]*vrMsg
		queue            *prque.Prque
		queues           map[string]int
		queued           map[common.Hash]*inject
		chainHeight      chainHeightFunc
		getBlock         getBlockByHashFunc
		saveBlock        saveBlockFunc
		blockBroadcaster blockBroadcasterFunc
		task             chan chan *vrTask
		filter           chan chan *dlTask
		done             chan common.Hash
		quit             chan struct{}
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &BlockFetcher{
				notify:           tt.fields.notify,
				notifyCount:      tt.fields.notifyCount,
				notified:         tt.fields.notified,
				fetching:         tt.fields.fetching,
				fetched:          tt.fields.fetched,
				finished:         tt.fields.finished,
				queue:            tt.fields.queue,
				queues:           tt.fields.queues,
				queued:           tt.fields.queued,
				chainHeight:      tt.fields.chainHeight,
				getBlock:         tt.fields.getBlock,
				saveBlock:        tt.fields.saveBlock,
				blockBroadcaster: tt.fields.blockBroadcaster,
				task:             tt.fields.task,
				filter:           tt.fields.filter,
				done:             tt.fields.done,
				quit:             tt.fields.quit,
			}
			f.handleInsert()
		})
	}
}

func TestBlockFetcher_handleNotify(t *testing.T) {
	// init test env, start fetcher loop
	fetcher := getFetcher()
	assert.NotNil(t, fetcher)

	notification := &vrMsg{hash: common.HexToHash("aaa"), number: 64, peerID: "test1"}

	for i := 0; i < hashLimit; i++ {
		fetcher.notifyCount[notification.peerID] = fetcher.notifyCount[notification.peerID] + 1
	}

	fetchTimer := time.NewTimer(0)
	assert.Equal(t, true, fetcher.handleNotify(notification, fetchTimer))
}

func TestBlockFetcher_handleNotify1(t *testing.T) {
	// init test env, start fetcher loop
	fetcher := getFetcher()
	assert.NotNil(t, fetcher)

	notification := &vrMsg{hash: common.HexToHash("aaa"), number: 64, peerID: "test1"}

	fetcher.chainHeight = func() model.AbstractBlock {
		return model.NewBlock(model.NewHeader(11, 10, common.HexToHash("ss"), common.HexToHash("fdfs"), common.StringToDiff("0x22"), big.NewInt(111), common.StringToAddress("fdsfds"), common.EncodeNonce(33)), nil, nil)
	}

	fetchTimer := time.NewTimer(0)
	assert.Equal(t, true, fetcher.handleNotify(notification, fetchTimer))
}

func TestBlockFetcher_handleNotify2(t *testing.T) {
	// init test env, start fetcher loop
	fetcher := getFetcher()
	assert.NotNil(t, fetcher)

	notification := &vrMsg{hash: common.HexToHash("aaa"), number: 64, peerID: "test1"}

	fetcher.chainHeight = func() model.AbstractBlock {
		return model.NewBlock(model.NewHeader(11, 62, common.HexToHash("ss"), common.HexToHash("fdfs"), common.StringToDiff("0x22"), big.NewInt(111), common.StringToAddress("fdsfds"), common.EncodeNonce(33)), nil, nil)
	}

	fetcher.fetching[notification.hash] = &vrMsg{}

	fetchTimer := time.NewTimer(0)
	assert.Equal(t, true, fetcher.handleNotify(notification, fetchTimer))
}

func TestBlockFetcher_handleNotify3(t *testing.T) {
	// init test env, start fetcher loop
	fetcher := getFetcher()
	assert.NotNil(t, fetcher)

	notification := &vrMsg{hash: common.HexToHash("aaa"), number: 64, peerID: "test1"}

	fetcher.chainHeight = func() model.AbstractBlock {
		return model.NewBlock(model.NewHeader(11, 62, common.HexToHash("ss"), common.HexToHash("fdfs"), common.StringToDiff("0x22"), big.NewInt(111), common.StringToAddress("fdsfds"), common.EncodeNonce(33)), nil, nil)
	}

	fetcher.finished[notification.hash] = &vrMsg{}

	fetchTimer := time.NewTimer(0)
	assert.Equal(t, true, fetcher.handleNotify(notification, fetchTimer))
}

func TestBlockFetcher_handleFetching(t *testing.T) {
	type fields struct {
		notify           chan *vrMsg
		notifyCount      map[string]int
		notified         map[common.Hash][]*vrMsg
		fetching         map[common.Hash]*vrMsg
		fetched          map[common.Hash][]*vrMsg
		finished         map[common.Hash]*vrMsg
		queue            *prque.Prque
		queues           map[string]int
		queued           map[common.Hash]*inject
		chainHeight      chainHeightFunc
		getBlock         getBlockByHashFunc
		saveBlock        saveBlockFunc
		blockBroadcaster blockBroadcasterFunc
		task             chan chan *vrTask
		filter           chan chan *dlTask
		done             chan common.Hash
		quit             chan struct{}
	}
	type args struct {
		fetchTimer *time.Timer
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &BlockFetcher{
				notify:           tt.fields.notify,
				notifyCount:      tt.fields.notifyCount,
				notified:         tt.fields.notified,
				fetching:         tt.fields.fetching,
				fetched:          tt.fields.fetched,
				finished:         tt.fields.finished,
				queue:            tt.fields.queue,
				queues:           tt.fields.queues,
				queued:           tt.fields.queued,
				chainHeight:      tt.fields.chainHeight,
				getBlock:         tt.fields.getBlock,
				saveBlock:        tt.fields.saveBlock,
				blockBroadcaster: tt.fields.blockBroadcaster,
				task:             tt.fields.task,
				filter:           tt.fields.filter,
				done:             tt.fields.done,
				quit:             tt.fields.quit,
			}
			f.handleFetching(tt.args.fetchTimer)
		})
	}
}

func TestBlockFetcher_handleVrTask(t *testing.T) {
	type fields struct {
		notify           chan *vrMsg
		notifyCount      map[string]int
		notified         map[common.Hash][]*vrMsg
		fetching         map[common.Hash]*vrMsg
		fetched          map[common.Hash][]*vrMsg
		finished         map[common.Hash]*vrMsg
		queue            *prque.Prque
		queues           map[string]int
		queued           map[common.Hash]*inject
		chainHeight      chainHeightFunc
		getBlock         getBlockByHashFunc
		saveBlock        saveBlockFunc
		blockBroadcaster blockBroadcasterFunc
		task             chan chan *vrTask
		filter           chan chan *dlTask
		done             chan common.Hash
		quit             chan struct{}
	}
	type args struct {
		task *vrTask
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &BlockFetcher{
				notify:           tt.fields.notify,
				notifyCount:      tt.fields.notifyCount,
				notified:         tt.fields.notified,
				fetching:         tt.fields.fetching,
				fetched:          tt.fields.fetched,
				finished:         tt.fields.finished,
				queue:            tt.fields.queue,
				queues:           tt.fields.queues,
				queued:           tt.fields.queued,
				chainHeight:      tt.fields.chainHeight,
				getBlock:         tt.fields.getBlock,
				saveBlock:        tt.fields.saveBlock,
				blockBroadcaster: tt.fields.blockBroadcaster,
				task:             tt.fields.task,
				filter:           tt.fields.filter,
				done:             tt.fields.done,
				quit:             tt.fields.quit,
			}
			f.handleVrTask(tt.args.task)
		})
	}
}

func TestBlockFetcher_filterDownloader(t *testing.T) {
	// init test env, start fetcher loop
	fetcher := getFetcher()
	assert.NotNil(t, fetcher)

	block := model.NewBlock(model.NewHeader(11, 10, common.HexToHash("ss"), common.HexToHash("fdfs"), common.StringToDiff("0x22"), big.NewInt(111), common.StringToAddress("fdsfds"), common.EncodeNonce(33)), nil, nil)

	task := &dlTask{
		peerID:      "aaa",
		catchupList: []*catchupRlp{{Block: block, SeenCommit: []*model.VoteMsg{}}},
	}

	fetcher.notified[block.Hash()] = []*vrMsg{}

	fetcher.filterDownloader(task)

}

func TestBlockFetcher_enqueue(t *testing.T) {
	// init test env, start fetcher loop
	fetcher := getFetcher()
	assert.NotNil(t, fetcher)

	for i := 0; i < blockLimit; i++ {
		fetcher.queues["test"] = fetcher.queues["test"] + 1
	}

	block := model.NewBlock(model.NewHeader(11, 10, common.HexToHash("ss"), common.HexToHash("fdfs"), common.StringToDiff("0x22"), big.NewInt(111), common.StringToAddress("fdsfds"), common.EncodeNonce(33)), nil, nil)

	catchup := &catchup{Block: block, SeenCommit: nil}

	fetcher.enqueue("test", catchup)
}

func TestBlockFetcher_enqueue1(t *testing.T) {
	// init test env, start fetcher loop
	fetcher := getFetcher()
	assert.NotNil(t, fetcher)

	block := model.NewBlock(model.NewHeader(11, 500, common.HexToHash("ss"), common.HexToHash("fdfs"), common.StringToDiff("0x22"), big.NewInt(111), common.StringToAddress("fdsfds"), common.EncodeNonce(33)), nil, nil)

	catchup := &catchup{Block: block, SeenCommit: nil}

	fetcher.enqueue("test", catchup)
}

func TestBlockFetcher_insert(t *testing.T) {
	// init test env, start fetcher loop
	fetcher := getFetcher()
	assert.NotNil(t, fetcher)

	block := model.NewBlock(model.NewHeader(11, 500, common.HexToHash("ss"), common.HexToHash("fdfs"), common.StringToDiff("0x22"), big.NewInt(111), common.StringToAddress("fdsfds"), common.EncodeNonce(33)), nil, nil)

	catchup := &catchup{Block: block, SeenCommit: nil}

	fetcher.getBlock = func(hashes common.Hash) model.AbstractBlock {
		return nil
	}

	fetcher.insert("test", catchup)
	time.Sleep(500 * time.Millisecond)
}

func TestBlockFetcher_insert1(t *testing.T) {
	// init test env, start fetcher loop
	fetcher := getFetcher()
	assert.NotNil(t, fetcher)

	block := model.NewBlock(model.NewHeader(11, 500, common.HexToHash("ss"), common.HexToHash("fdfs"), common.StringToDiff("0x22"), big.NewInt(111), common.StringToAddress("fdsfds"), common.EncodeNonce(33)), nil, nil)

	catchup := &catchup{Block: block, SeenCommit: nil}

	fetcher.getBlock = func(hashes common.Hash) model.AbstractBlock {
		return model.NewBlock(model.NewHeader(11, 499, common.HexToHash("ssdfs"), common.HexToHash("fdfsf"), common.StringToDiff("0x22"), big.NewInt(111), common.StringToAddress("fdsfds"), common.EncodeNonce(33)), nil, nil)
	}

	fetcher.saveBlock = func(block model.AbstractBlock, seenCommits []model.AbstractVerification) error {
		return errors.New("sss")
	}

	fetcher.insert("test", catchup)

	time.Sleep(500 * time.Millisecond)
}

func TestBlockFetcher_rescheduleFetch(t *testing.T) {
	type fields struct {
		notify           chan *vrMsg
		notifyCount      map[string]int
		notified         map[common.Hash][]*vrMsg
		fetching         map[common.Hash]*vrMsg
		fetched          map[common.Hash][]*vrMsg
		finished         map[common.Hash]*vrMsg
		queue            *prque.Prque
		queues           map[string]int
		queued           map[common.Hash]*inject
		chainHeight      chainHeightFunc
		getBlock         getBlockByHashFunc
		saveBlock        saveBlockFunc
		blockBroadcaster blockBroadcasterFunc
		task             chan chan *vrTask
		filter           chan chan *dlTask
		done             chan common.Hash
		quit             chan struct{}
	}
	type args struct {
		fetch *time.Timer
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &BlockFetcher{
				notify:           tt.fields.notify,
				notifyCount:      tt.fields.notifyCount,
				notified:         tt.fields.notified,
				fetching:         tt.fields.fetching,
				fetched:          tt.fields.fetched,
				finished:         tt.fields.finished,
				queue:            tt.fields.queue,
				queues:           tt.fields.queues,
				queued:           tt.fields.queued,
				chainHeight:      tt.fields.chainHeight,
				getBlock:         tt.fields.getBlock,
				saveBlock:        tt.fields.saveBlock,
				blockBroadcaster: tt.fields.blockBroadcaster,
				task:             tt.fields.task,
				filter:           tt.fields.filter,
				done:             tt.fields.done,
				quit:             tt.fields.quit,
			}
			f.rescheduleFetch(tt.args.fetch)
		})
	}
}

func TestBlockFetcher_forgetHash(t *testing.T) {
	type fields struct {
		notify           chan *vrMsg
		notifyCount      map[string]int
		notified         map[common.Hash][]*vrMsg
		fetching         map[common.Hash]*vrMsg
		fetched          map[common.Hash][]*vrMsg
		finished         map[common.Hash]*vrMsg
		queue            *prque.Prque
		queues           map[string]int
		queued           map[common.Hash]*inject
		chainHeight      chainHeightFunc
		getBlock         getBlockByHashFunc
		saveBlock        saveBlockFunc
		blockBroadcaster blockBroadcasterFunc
		task             chan chan *vrTask
		filter           chan chan *dlTask
		done             chan common.Hash
		quit             chan struct{}
	}
	type args struct {
		hash common.Hash
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &BlockFetcher{
				notify:           tt.fields.notify,
				notifyCount:      tt.fields.notifyCount,
				notified:         tt.fields.notified,
				fetching:         tt.fields.fetching,
				fetched:          tt.fields.fetched,
				finished:         tt.fields.finished,
				queue:            tt.fields.queue,
				queues:           tt.fields.queues,
				queued:           tt.fields.queued,
				chainHeight:      tt.fields.chainHeight,
				getBlock:         tt.fields.getBlock,
				saveBlock:        tt.fields.saveBlock,
				blockBroadcaster: tt.fields.blockBroadcaster,
				task:             tt.fields.task,
				filter:           tt.fields.filter,
				done:             tt.fields.done,
				quit:             tt.fields.quit,
			}
			f.forgetHash(tt.args.hash)
		})
	}
}

func TestBlockFetcher_forgetBlock(t *testing.T) {
	type fields struct {
		notify           chan *vrMsg
		notifyCount      map[string]int
		notified         map[common.Hash][]*vrMsg
		fetching         map[common.Hash]*vrMsg
		fetched          map[common.Hash][]*vrMsg
		finished         map[common.Hash]*vrMsg
		queue            *prque.Prque
		queues           map[string]int
		queued           map[common.Hash]*inject
		chainHeight      chainHeightFunc
		getBlock         getBlockByHashFunc
		saveBlock        saveBlockFunc
		blockBroadcaster blockBroadcasterFunc
		task             chan chan *vrTask
		filter           chan chan *dlTask
		done             chan common.Hash
		quit             chan struct{}
	}
	type args struct {
		hash common.Hash
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &BlockFetcher{
				notify:           tt.fields.notify,
				notifyCount:      tt.fields.notifyCount,
				notified:         tt.fields.notified,
				fetching:         tt.fields.fetching,
				fetched:          tt.fields.fetched,
				finished:         tt.fields.finished,
				queue:            tt.fields.queue,
				queues:           tt.fields.queues,
				queued:           tt.fields.queued,
				chainHeight:      tt.fields.chainHeight,
				getBlock:         tt.fields.getBlock,
				saveBlock:        tt.fields.saveBlock,
				blockBroadcaster: tt.fields.blockBroadcaster,
				task:             tt.fields.task,
				filter:           tt.fields.filter,
				done:             tt.fields.done,
				quit:             tt.fields.quit,
			}
			f.forgetBlock(tt.args.hash)
		})
	}
}
