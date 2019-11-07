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
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/chain-communication"
	"github.com/dipperin/dipperin-core/core/mine/minemsg"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third-party/p2p"
	"math/big"
	"sync/atomic"
)

type WorkerId string

type Master interface {
	Start()
	Stop()
	CurrentCoinbaseAddress() common.Address
	SetCoinbaseAddress(addr common.Address)
	SetMineGasConfig(gasFloor, gasCeil uint64)
	OnNewBlock(block model.AbstractBlock)
	Workers() map[WorkerId]WorkerForMaster

	GetReward(address common.Address) *big.Int
	GetPerformance(address common.Address) uint64

	// whether the mining is ongoing
	Mining() bool
	// cur mine block tx count
	MineTxCount() int

	SetMsgSigner(MsgSigner chain_communication.PbftSigner)

	GetMsgSigner() chain_communication.PbftSigner
	SpendableMaster
	// Done: 1. add get worker's work,
	// Done: 2. worker's coin count method,
	// Done: 3. add withdrawal coins method
	// :‑Þ
}

type mineMaster interface {
	registerWorker(worker WorkerForMaster)
	unRegisterWorker(workerId WorkerId)
	startWaitTimer()
	getWorker(id WorkerId) WorkerForMaster
}

type SpendableMaster interface {
	RetrieveReward(address common.Address)
}

type WorkerForMaster interface {
	Start()
	Stop()
	GetId() WorkerId

	SendNewWork(msgCode int, work minemsg.Work)

	SetCoinbase(coinbase common.Address)
	CurrentCoinbaseAddress() common.Address
	//Nickname() string
}

type MasterServer interface {
	RegisterWorker(worker WorkerForMaster)
	UnRegisterWorker(workerId WorkerId)
	ReceiveMsg(workerID WorkerId, code uint64, msg interface{})
	// for p2p msg
	OnNewMsg(msg p2p.Msg, p chain_communication.PmAbstractPeer) error
	// only for worker, do nothing
	SetMineMasterPeer(peer chain_communication.PmAbstractPeer)
}

// there is only one workManager to manage all worker's works
type workManager interface {
	submitBlock(workerAddress common.Address, block model.AbstractBlock)
	getPerformance(address common.Address) uint64
	getReward(address common.Address) *big.Int
	onNewBlock(block model.AbstractBlock)
	SetMsgSigner(MsgSigner chain_communication.PbftSigner)
	spendableWorkManager
	partialSpendableWorkManager
}

type dispatcher interface {
	onNewBlock(block model.AbstractBlock) error
	dispatchNewWork() error
	curWorkBlock() model.AbstractBlock
	SetMsgSigner(MsgSigner chain_communication.PbftSigner)
}

type spendableWorkManager interface {
	clearPerformance(address common.Address)
	clearReward(address common.Address)
}

type partialSpendableWorkManager interface {
	subtractPerformance(address common.Address, performance uint64)
	subtractReward(address common.Address, reward *big.Int)
}

// workerPerformance keeps records of all worker's work
// different worker could have different types of workerPerformance.
type workerPerformance interface {
	getPerformance() uint64
	updatePerformance()
	setPerformance(uint64)
}

// rewardDistributor distributes the rewards based on the performance
type rewardDistributor interface {
	divideReward(coinbase *big.Int) map[common.Address]*big.Int
}

// calculableBlock could calculate coinbase and transaction fees from
// the callee block
//type calculableBlock interface {
//	GetCoinbase() *big.Int
//	GetTransactionFees() *big.Int
//}

//type NodeContext interface {
//	CoinbaseAddress() common.Address
//	BuildWorks(newBlock model.AbstractBlock, workerLen int) (workMsgCode int, works []minemsg.Work)
//	BuildWaitPackBlock(coinbaseAddr common.Address) model.AbstractBlock
//	BroadcastWaitVerifyBlock(block model.AbstractBlock)
//
//	ChainReader() state_processor.ChainReader
//}

type BlockBuilder interface {
	SetMsgSigner(MsgSigner chain_communication.PbftSigner)
	GetMsgSigner() chain_communication.PbftSigner
	BuildWaitPackBlock(coinbaseAddr common.Address, gasFloor, gasCeil uint64) model.AbstractBlock
}

type BlockBroadcaster interface {
	BroadcastMinedBlock(block model.AbstractBlock)
}

type MineConfig struct {
	GasFloor         *atomic.Value // Target gas floor for mined blocks.
	GasCeil          *atomic.Value // Target gas ceiling for mined blocks.
	CoinbaseAddress  *atomic.Value
	BlockBuilder     BlockBuilder
	BlockBroadcaster BlockBroadcaster
}

func (conf *MineConfig) GetMsgSigner() chain_communication.PbftSigner {
	return conf.BlockBuilder.GetMsgSigner()
}

func (conf *MineConfig) SetMsgSigner(MsgSigner chain_communication.PbftSigner) {
	conf.BlockBuilder.SetMsgSigner(MsgSigner)
}

func (conf *MineConfig) GetGasFloor() (result uint64) {
	if v := conf.GasFloor.Load(); v != nil {
		return v.(uint64)
	}
	return
}

func (conf *MineConfig) GetGasCeil() (result uint64) {
	if v := conf.GasCeil.Load(); v != nil {
		return v.(uint64)
	}
	return
}

func (conf *MineConfig) GetCoinbaseAddr() (result common.Address) {
	if v := conf.CoinbaseAddress.Load(); v != nil {
		return v.(common.Address)
	}
	return
}
