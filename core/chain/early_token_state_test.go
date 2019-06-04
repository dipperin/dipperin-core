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


package chain

import (
	"fmt"
	"github.com/dipperin/dipperin-core/core/chain/state-processor"
	"github.com/dipperin/dipperin-core/core/contract"
	"github.com/dipperin/dipperin-core/core/economy-model"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/dipperin/dipperin-core/third-party/log/mpt_log"
	"github.com/stretchr/testify/assert"
	"math/big"
	"reflect"
	"testing"
)


/*

test foundation contract execution is correct

*/

func TestProcessEarlyContract(t *testing.T) {
	mpt_log.InitMptLogger(log.LvlDebug, "TestProcessEarlyContract", true)

	eModel := economy_model.MakeDipperinEconomyModel(&earlyContractFakeChainService{}, economy_model.DIPProportion)
	cReader := &fakeAccountDBChain{}
	kvDB, originStateRoot := createTestStateDB(t)

	// Test processing block can correctly modify the contract value
	trDB := state_processor.NewStateStorageWithCache(kvDB)
	aStateDB, err := NewBlockProcessor(cReader, originStateRoot, trDB)
	assert.NoError(t, err)

	tmpB := createBlock(20)
	err = aStateDB.Process(tmpB, eModel)
	assert.NoError(t, err)

	hash, err:=aStateDB.Finalise()
	assert.NoError(t,err)
	fmt.Println(hash.Hex())


	//the root of finalise and commit
	fHash, err := aStateDB.Finalise()
	assert.NoError(t, err)

	// start with kv db completely
	trDB = state_processor.NewStateStorageWithCache(kvDB)
	aStateDB, err = NewBlockProcessor(cReader, originStateRoot, trDB)
	fmt.Println("1111111111111",aStateDB.Finalised())
	assert.NoError(t, err)
	err = aStateDB.Process(tmpB, eModel)
	assert.NoError(t, err)

	// Do a commit, then take the contract data from kvDB
	fmt.Println("22222222222222",aStateDB.Finalised())
	cHash, err := aStateDB.Commit()
	fmt.Println("333333333333333",aStateDB.Finalised())
	assert.NoError(t, err)
	assert.Equal(t, fHash, cHash)

	// Take out contract data from KVDB
	trDB = state_processor.NewStateStorageWithCache(kvDB)
	aStateDB, err = NewBlockProcessor(cReader, cHash, trDB)
	fmt.Println("44444444444444",aStateDB.Finalised())
	assert.NoError(t, err)


	earlyTCV, err := aStateDB.GetContract(contract.EarlyContractAddress, reflect.TypeOf(contract.EarlyRewardContract{}))
	assert.NoError(t, err)

	earlyTC := earlyTCV.Interface().(*contract.EarlyRewardContract)
	assert.NoError(t, err)
	assert.NotNil(t, earlyTC.Balances[aliceAddr.Hex()])
	assert.Equal(t, 1, earlyTC.Balances[aliceAddr.Hex()].Cmp(big.NewInt(0)))
}