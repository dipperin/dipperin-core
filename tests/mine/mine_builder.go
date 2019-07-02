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

package mine_spec

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/chain-communication"
	"github.com/dipperin/dipperin-core/core/mine/minemaster"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third-party/p2p"
	"math/big"
)

func MasterBuilder() minemaster.Master {

	return &fakeMaster{}
}

type fakeMaster struct {
}

func (m *fakeMaster) SetMineGasConfig(gasFloor, gasCeil uint64) {
	panic("implement me")
}

func (m *fakeMaster) Start() {
	panic("implement me")
}

func (m *fakeMaster) Stop() {
	panic("implement me")
}

func (m *fakeMaster) CurrentCoinbaseAddress() common.Address {
	panic("implement me")
}

func (m *fakeMaster) SetCoinbaseAddress(addr common.Address) {
	panic("implement me")
}

func (m *fakeMaster) OnNewBlock(block model.AbstractBlock) {
	panic("implement me")
}

func (m *fakeMaster) Workers() map[minemaster.WorkerId]minemaster.WorkerForMaster {
	panic("implement me")
}

func (m *fakeMaster) GetReward(address common.Address) *big.Int {
	panic("implement me")
}

func (m *fakeMaster) GetPerformance(address common.Address) uint64 {
	panic("implement me")
}

func (m *fakeMaster) Mining() bool {
	panic("implement me")
}

func (m *fakeMaster) MineTxCount() int {
	panic("implement me")
}

func (m *fakeMaster) RetrieveReward(address common.Address) {
	panic("implement me")
}

func MasterServerBuilder() minemaster.MasterServer {
	return &FakeMasterServer{
		Workers: make(map[string]minemaster.WorkerForMaster),
	}
}

type FakeMasterServer struct {
	Workers map[string]minemaster.WorkerForMaster
}

func (ms *FakeMasterServer) RegisterWorker(worker minemaster.WorkerForMaster) {
	ms.Workers[string(worker.GetId())] = worker
}

func (ms *FakeMasterServer) UnRegisterWorker(workerId minemaster.WorkerId) {
	delete(ms.Workers, string(workerId))
}

func (ms *FakeMasterServer) ReceiveMsg(workerID minemaster.WorkerId, code uint64, msg interface{}) {

}

func (ms *FakeMasterServer) OnNewMsg(msg p2p.Msg, p chain_communication.PmAbstractPeer) error {
	panic("implement me")
}

func (ms *FakeMasterServer) SetMineMasterPeer(peer chain_communication.PmAbstractPeer) {
	panic("implement me")
}
