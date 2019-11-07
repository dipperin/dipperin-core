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

package minemaster

import (
	"errors"
	"github.com/dipperin/dipperin-core/core/chain-communication"
	"math/big"
	"sync/atomic"
	"testing"
	"time"

	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/g-event"
	"github.com/dipperin/dipperin-core/core/mine/minemsg"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/tests/factory"
	"github.com/stretchr/testify/assert"
)

var testMineConfig = MineConfig{
	CoinbaseAddress: &atomic.Value{},
}

type mockDispatch struct {
	curBlock        model.AbstractBlock
	dispatchError   error
	onNewBlockError error
}

func (dispatcher *mockDispatch)SetMsgSigner(MsgSigner chain_communication.PbftSigner){

}

func (dispatcher *mockDispatch) onNewBlock(block model.AbstractBlock) error {
	return dispatcher.onNewBlockError
}

func (dispatcher *mockDispatch) dispatchNewWork() error {
	return dispatcher.dispatchError
}

func (dispatcher *mockDispatch) curWorkBlock() model.AbstractBlock {
	return dispatcher.curBlock
}

type mockWorker struct {
	coinbase common.Address
	workerId WorkerId
}

func (worker *mockWorker) Start() {

}

func (worker *mockWorker) Stop() {

}

func (worker *mockWorker) GetId() WorkerId {
	return worker.workerId
}

func (worker *mockWorker) SendNewWork(msgCode int, work minemsg.Work) {

}

func (worker *mockWorker) SetCoinbase(coinbase common.Address) {
	worker.coinbase = coinbase
}

func (worker *mockWorker) CurrentCoinbaseAddress() common.Address {
	return worker.coinbase
}

func Test_newMaster(t *testing.T) {
	assert.NotNil(t, newMaster(MineConfig{
		CoinbaseAddress: &atomic.Value{},
	}))
}

func Test_master_Start(t *testing.T) {
	diff := common.HexToDiff("0x1effffff")
	block := factory.CreateBlock2(diff, 1)
	dispatch := mockDispatch{curBlock: block}

	nM := testMasterBuilder(testMineConfig)

	nM.setWorkDispatcher(&dispatch)

	nM.Start()

	assert.Equal(t, nM.Mining(), true)

	nM.Start()

	nM.Stop()

	dispatch2 := mockDispatch{curBlock: block, dispatchError: errors.New("test")}

	nM.setWorkDispatcher(&dispatch2)

	assert.Panics(t, func() {
		nM.Start()
	})
}

func Test_master_startWaitTimer(t *testing.T) {
	diff := common.HexToDiff("0x1effffff")
	block := factory.CreateBlock2(diff, 1)
	dispatch := mockDispatch{curBlock: block}

	nM := testMasterBuilder(testMineConfig)

	nM.setWorkDispatcher(&dispatch)

	nM.Start()

	waitTimeout = 1 * time.Millisecond

	nM.startWaitTimer()

	time.Sleep(2 * time.Millisecond)

	nM.Stop()

	nM.startWaitTimer()

	nM.Start()

	dispatch2 := mockDispatch{curBlock: block, dispatchError: errors.New("test")}

	nM.setWorkDispatcher(&dispatch2)

	nM.startWaitTimer()

	time.Sleep(2 * time.Millisecond)

}

func Test_master_Stop(t *testing.T) {
	diff := common.HexToDiff("0x1effffff")
	block := factory.CreateBlock2(diff, 1)
	dispatch := mockDispatch{curBlock: block}

	nM := testMasterBuilder(testMineConfig)

	nM.setWorkDispatcher(&dispatch)

	nM.Start()

	nM.Stop()

	nM.Stop()
}

func Test_master_stopped(t *testing.T) {
	diff := common.HexToDiff("0x1effffff")
	block := factory.CreateBlock2(diff, 1)
	dispatch := mockDispatch{curBlock: block}

	nM := testMasterBuilder(testMineConfig)

	nM.setWorkDispatcher(&dispatch)

	assert.Equal(t, nM.stopped(), true)

	nM.Start()

	assert.Equal(t, nM.stopped(), false)
}

func Test_master_registerWorker(t *testing.T) {

	diff := common.HexToDiff("0x1effffff")
	block := factory.CreateBlock2(diff, 1)
	dispatch := mockDispatch{curBlock: block}

	worker := mockWorker{
		workerId: "1",
	}

	nM := testMasterBuilder(testMineConfig)
	nM.setWorkDispatcher(&dispatch)

	nM.Start()

	nM.registerWorker(&worker)

	time.Sleep(100 * time.Millisecond)

	nM.Stop()

	nM.Start()

	assert.Equal(t, len(nM.Workers()), 1)

	assert.Equal(t, nM.getWorker("1"), &worker)

	nM.unRegisterWorker("1")

	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, len(nM.Workers()), 0)

}

func Test_master_MineTxCount(t *testing.T) {
	dispatch := mockDispatch{curBlock: nil}

	nM := testMasterBuilder(testMineConfig)
	nM.setWorkDispatcher(&dispatch)

	assert.Equal(t, nM.MineTxCount(), 0)

	diff := common.HexToDiff("0x1effffff")
	block := factory.CreateBlock2(diff, 1)
	dispatch2 := mockDispatch{curBlock: block}

	nM.setWorkDispatcher(&dispatch2)

	assert.Equal(t, nM.MineTxCount(), 2)
}

func Test_master_SetCoinbaseAddress(t *testing.T) {
	dispatch := mockDispatch{curBlock: nil}

	nM := testMasterBuilder(testMineConfig)
	nM.setWorkDispatcher(&dispatch)

	nM.SetCoinbaseAddress(common.HexToAddress("0x00000000000000000000000000000000000000001234"))

	assert.Equal(t, nM.CurrentCoinbaseAddress().Hex(), "0x00000000000000000000000000000000000000001234")
}

func Test_master_OnNewBlock(t *testing.T) {
	dispatch := mockDispatch{curBlock: nil}

	nM := testMasterBuilder(testMineConfig)

	nM.setWorkDispatcher(&dispatch)

	diff := common.HexToDiff("0x1effffff")
	block := factory.CreateBlock2(diff, 1)

	nM.OnNewBlock(block)

	nM.Start()

	nM.OnNewBlock(block)
}

func Test_master_GetPerformance(t *testing.T) {
	dispatch := mockDispatch{curBlock: nil}

	nM := testMasterBuilder(testMineConfig)

	nM.setWorkDispatcher(&dispatch)

	assert.Equal(t, nM.GetPerformance(common.HexToAddress("0x1234")), uint64(0x0))
}

func Test_master_doOnNewBlock(t *testing.T) {
	dispatch := mockDispatch{curBlock: nil}

	nM := testMasterBuilder(testMineConfig)
	nM.setWorkDispatcher(&dispatch)

	diff := common.HexToDiff("0x1effffff")
	block := factory.CreateBlock2(diff, 2)

	nM.Start()

	g_event.Send(g_event.NewBlockInsertEvent, *block)

	block2 := factory.CreateBlock2(diff, 1)

	g_event.Send(g_event.NewBlockInsertEvent, *block2)

	dispatch2 := mockDispatch{curBlock: nil, onNewBlockError: errors.New("test")}

	nM.setWorkDispatcher(&dispatch2)

	g_event.Send(g_event.NewBlockInsertEvent, *block)

	nM.Stop()
}

func Test_master_GetReward(t *testing.T) {
	dispatch := mockDispatch{curBlock: nil}

	nM := testMasterBuilder(testMineConfig)

	nM.setWorkDispatcher(&dispatch)

	assert.Equal(t, nM.GetReward(common.HexToAddress("0x1234")), big.NewInt(0))
}

func Test_master_RetrieveReward(t *testing.T) {
	dispatch := mockDispatch{curBlock: nil}

	nM := testMasterBuilder(testMineConfig)

	nM.setWorkDispatcher(&dispatch)

	nM.RetrieveReward(common.HexToAddress("0x1234"))
}

func testMasterBuilder(config MineConfig) *master {

	manager := newDefaultWorkManager(config)

	return &master{
		MineConfig: config,

		workManager:          manager,
		workers:              map[WorkerId]WorkerForMaster{},
		registerWorkerChan:   make(chan WorkerForMaster),
		unRegisterWorkerChan: make(chan WorkerId),
		onNewBlockChan:       make(chan model.AbstractBlock),
		startTimerChan:       make(chan struct{}),
		timeoutChan:          make(chan struct{}),
	}
}

func init() {
	g_event.Add(g_event.NewBlockInsertEvent)
}
