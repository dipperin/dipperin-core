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
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/tests/factory"
	"github.com/dipperin/dipperin-core/tests/peer"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

var workers map[WorkerId]WorkerForMaster

var block *model.Block

func fakeGetWorkersFunc() map[WorkerId]WorkerForMaster {
	return workers
}

func Test_newWorkDispatcher(t *testing.T) {
	nwd := newWorkDispatcher(testMineConfig, fakeGetWorkersFunc)

	assert.NotNil(t, nwd)
}

type fakeBlockBuilder struct {
}

func (b *fakeBlockBuilder) GetMsgSigner() chain_communication.PbftSigner {
	return nil
}

func (b *fakeBlockBuilder) SetMsgSigner(MsgSigner chain_communication.PbftSigner) {

}

func (b *fakeBlockBuilder) BuildWaitPackBlock(coinbaseAddr common.Address, gasFloor, gasCeil uint64) model.AbstractBlock {
	return block
}

func Test_workDispatcher_onNewBlock(t *testing.T) {
	testMineConfig = MineConfig{
		CoinbaseAddress: &atomic.Value{},
		BlockBuilder:    &fakeBlockBuilder{},
		GasFloor:        &atomic.Value{},
		GasCeil:         &atomic.Value{},
	}

	gasFloor := chain_config.BlockGasLimit
	gasCeil := chain_config.BlockGasLimit
	testMineConfig.GasFloor.Store(uint64(gasFloor))
	testMineConfig.GasCeil.Store(uint64(gasCeil))

	nwd := newWorkDispatcher(testMineConfig, fakeGetWorkersFunc)

	diff := common.HexToDiff("0x1effffff")
	fakeBlock := factory.CreateBlock2(diff, 2)

	block = fakeBlock

	err := nwd.onNewBlock(fakeBlock)
	assert.Error(t, err)

	p := peer_spec.PeerBuilder()
	workers = map[WorkerId]WorkerForMaster{
		"123": newRemoteWorker(p, common.HexToAddress("0x1234"), "123"),
	}

	err = nwd.onNewBlock(fakeBlock)
	assert.NoError(t, err)

	fakeBlock2 := factory.CreateBlock2(diff, 1)
	err = nwd.onNewBlock(fakeBlock2)
	assert.Error(t, err)

}
