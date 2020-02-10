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

package chainstate

import (
	"github.com/dipperin/dipperin-core/core/chainconfig"
	"github.com/dipperin/dipperin-core/core/cschain/chainwriter"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	
	"github.com/dipperin/dipperin-core/common"
)

func TestNewChainState(t *testing.T) {
	writerF := chainwriter.NewChainWriterFactory()
	cs := NewChainState(&ChainStateConfig{
		DataDir:       "",
		WriterFactory: writerF,
		ChainConfig:   chainconfig.GetChainConfig(),
	})
	
	_, err := cs.BlockProcessorByNumber(22)
	assert.Error(t, err)
	
	_, err = cs.CurrentState()
	assert.Error(t, err)
	
	_, err = cs.StateAtByBlockNumber(11)
	assert.Error(t, err)
	assert.NotNil(t, cs)
	
	_, err = cs.AccountStateDB(common.Hash{})
	assert.NoError(t, err)
}

func TestChainState_initConfigAndDB(t *testing.T) {
	cs := &ChainState{ChainStateConfig: &ChainStateConfig{}}
	
	assert.Nil(t, cs.GetDB())
	assert.Nil(t, cs.ChainConfig)
	assert.Nil(t, cs.ChainDB)
	assert.Nil(t, cs.StateStorage)
	assert.Nil(t, cs.EconomyModel)
	
	cs.initConfigAndDB("test")
	
	assert.NotNil(t, cs.ChainConfig)
	assert.NotNil(t, cs.ChainDB)
	assert.NotNil(t, cs.StateStorage)
	assert.NotNil(t, cs.EconomyModel)
}

func Test_initEthDB(t *testing.T) {
	testMemDB := initEthDB("mem")
	assert.NotNil(t, testMemDB)
	
	testMemDB2 := initEthDB("test")
	assert.NotNil(t, testMemDB2)
	
	testMemDB3 := initEthDB("")
	assert.NotNil(t, testMemDB3)
	
	dbPath := "/tmp/chainState/test_init_db"
	defer os.RemoveAll(dbPath)
	assert.NotNil(t, initEthDB(dbPath))
	assert.Panics(t, func() {
		initEthDB(dbPath)
	})
}
