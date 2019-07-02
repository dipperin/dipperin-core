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

package contract

import (
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/stretchr/testify/assert"
	"testing"
	//"github.com/dipperin/dipperin-core/core/chain"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/core/chain/state-processor"
	"github.com/dipperin/dipperin-core/core/contract"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/dipperin/dipperin-core/third-party/log/mpt_log"
	"reflect"
)

func TestERC20Check(t *testing.T) {
	return
	mpt_log.InitMptLogger(log.LvlDebug, "TestERC20Check", true)
	dataDir := "/home/qydev/tmp/dipperin_apps/node/full_chain_data"
	kvDB, err := ethdb.NewLDBDatabase(dataDir, 0, 0)
	assert.NoError(t, err)

	c := newFkc(kvDB)
	state, err := c.fullChain.CurrentState()
	assert.NoError(t, err)

	erc20V, err := state.GetContract(common.HexToAddress("0x0010BCdbc9822289edba939e235E6435A39B1fcE2785"), reflect.TypeOf(contract.BuiltInERC20Token{}))
	assert.NoError(t, err)

	erc20 := erc20V.Interface().(*contract.BuiltInERC20Token)
	b := erc20.BalanceOf(common.HexToAddress("0x00005a55a149b9935F4Dde63631EF2Beb8A70dAcd62D"))
	fmt.Println(b.ToInt().String())
}

func newFkc(db ethdb.Database) *fakeContext {
	model.SetBlockRlpHandler(&model.PBFTBlockRlpHandler{})
	model.SetBlockJsonHandler(&model.PBFTBlockJsonHandler{})

	//c := &fakeContext{
	//	db: db,
	//	blockDecoder: model.MakeDefaultBlockDecoder(),
	//	chainConfig: *chain_config.GetChainConfig(),
	//	stateStorage: state_processor.NewStateStorageWithCache(db),
	//}
	//c.fullChain = chain.MakeFullChain(c)

	return nil
}

type Chain interface {
	CurrentState() (*state_processor.AccountStateDB, error)
}

type fakeContext struct {
	db           ethdb.Database
	blockDecoder model.BlockDecoder
	fullChain    Chain
	chainConfig  chain_config.ChainConfig
	stateStorage state_processor.StateStorage
}
