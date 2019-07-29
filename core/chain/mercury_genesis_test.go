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
	"encoding/json"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/consts"
	"github.com/dipperin/dipperin-core/common/g-error"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/core/chain/registerdb"
	"github.com/dipperin/dipperin-core/core/chain/state-processor"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"testing"
)

func TestSetupGenesisBlock(t *testing.T) {
	err := os.Setenv("boots_env", "test")
	assert.NoError(t, err)

	gFPath := filepath.Join(util.HomeDir(), "softwares", "dipperin_deploy")
	assert.True(t, pathIsExist(gFPath))
	gFPath = filepath.Join(gFPath, "genesis.json")

	// delete genesis.json first
	os.Remove(gFPath)

	// get genesis
	defaultGenesis := createGenesis()

	// genesis is nil
	_, blockHash, err := SetupGenesisBlock(nil)
	assert.Error(t, err)
	assert.Equal(t, common.Hash{}, blockHash)

	// setup genesis successful
	defaultGenesis.Difficulty = common.Difficulty{}
	_, blockHash, err = SetupGenesisBlock(defaultGenesis)
	assert.NoError(t, err)
	assert.NotNil(t, blockHash)

	// genesis and stored genesis is match
	genesisBlock, err := defaultGenesis.Prepare()
	defaultGenesis.ChainDB.SaveBlockHash(genesisBlock.Hash(), genesisBlock.Number())
	assert.NoError(t, err)
	_, blockHash, err = SetupGenesisBlock(defaultGenesis)
	assert.NoError(t, err)
	assert.Equal(t, genesisBlock.Hash(), blockHash)

	// genesis and stored genesis not match
	block := createBlock(0)
	defaultGenesis.ChainDB.SaveBlockHash(block.Hash(), block.Number())
	assert.NoError(t, err)
	_, blockHash, err = SetupGenesisBlock(defaultGenesis)
	expect := &GenesisMismatchError{block.Hash(), genesisBlock.Hash()}
	log.Debug("GenesisMismatchError", "err", expect.Error())
	assert.Equal(t, expect, err)
	assert.Equal(t, genesisBlock.Hash(), blockHash)

	// config is nil
	defaultGenesis.Config = nil
	_, blockHash, err = SetupGenesisBlock(defaultGenesis)
	assert.Equal(t, errGenesisNoConfig, err)
	assert.Equal(t, common.Hash{}, blockHash)
}

func TestGenesis_Set(t *testing.T) {
	defaultGenesis := createGenesis()

	defaultGenesis.SetVerifiers(nil)
	assert.Nil(t, defaultGenesis.Verifiers)

	defaultGenesis.SetChainDB(nil)
	assert.Nil(t, defaultGenesis.ChainDB)

	defaultGenesis.SetAccountStateProcessor(nil)
	assert.Nil(t, defaultGenesis.AccountStateProcessor)
}

func TestDefaultGenesisBlock(t *testing.T) {
	gFPath := filepath.Join(util.HomeDir(), "softwares", "dipperin_deploy")
	assert.True(t, pathIsExist(gFPath))
	gFPath = filepath.Join(gFPath, "genesis.json")

	// delete genesis.json first
	os.Remove(gFPath)

	// write wrong info into file
	err := ioutil.WriteFile(gFPath, []byte{123}, 0666)
	assert.NoError(t, err)

	// can't get default genesis from file
	defaultGenesis := createGenesis()
	err = os.Remove(gFPath)
	assert.NoError(t, err)

	// write genesis into file
	accounts := make(map[string]int64)
	accounts[aliceAddr.String()] = 10
	cfg := genesisCfgFile{
		Nonce:     defaultGenesis.Nonce,
		Accounts:  accounts,
		Verifiers: []string{bobAddr.String()},
	}
	bytes, err := json.Marshal(cfg)
	assert.NoError(t, err)

	err = ioutil.WriteFile(gFPath, bytes, 0666)
	assert.NoError(t, err)

	// get info from genesis.json
	defaultGenesis = createGenesis()
	assert.Equal(t, cfg.Nonce, defaultGenesis.Nonce)
	assert.Equal(t, big.NewInt(0).Mul(big.NewInt(10), big.NewInt(consts.DIP)), defaultGenesis.Alloc[aliceAddr])
	assert.Equal(t, bobAddr, defaultGenesis.Verifiers[0])

	err = os.Remove(gFPath)
	assert.NoError(t, err)

	// negative alice balance
	accounts[aliceAddr.String()] = -1
	cfg.Accounts = accounts

	bytes, err = json.Marshal(cfg)
	assert.NoError(t, err)

	err = ioutil.WriteFile(gFPath, bytes, 0666)
	assert.NoError(t, err)
	assert.Panics(t, func() {
		createGenesis()
	})

	err = os.Remove(gFPath)
	assert.NoError(t, err)
}

func TestGenesis_Valid(t *testing.T) {
	defaultGenesis := createGenesis()
	assert.True(t, defaultGenesis.Valid())

	defaultGenesis.Alloc[aliceAddr] = big.NewInt(0).Add(consts.MaxAmount, big.NewInt(1000))
	assert.False(t, defaultGenesis.Valid())
}

func TestGenesis_Commit_Error(t *testing.T) {
	defaultGenesis := createGenesis()

	rProcessor, err := registerdb.MakeGenesisRegisterProcessor(fakeStateStorage{})
	assert.NoError(t, err)
	defaultGenesis.RegisterProcessor = rProcessor

	err = defaultGenesis.Commit(createBlock(0))
	assert.Equal(t, TrieError, err)

	sProcessor, err := state_processor.MakeGenesisAccountStateProcessor(fakeStateStorage{})
	assert.NoError(t, err)
	defaultGenesis.AccountStateProcessor = sProcessor

	err = defaultGenesis.Commit(createBlock(0))
	assert.Equal(t, TrieError, err)
}

func TestGenesis_Prepare_Error(t *testing.T) {
	defaultGenesis := createGenesis()

	// prepare registerDB error
	rProcessor, err := registerdb.MakeGenesisRegisterProcessor(fakeStateStorage{setErr: TrieError})
	assert.NoError(t, err)
	defaultGenesis.RegisterProcessor = rProcessor
	block, err := defaultGenesis.Prepare()
	assert.Equal(t, TrieError, err)

	// set balance error
	sProcessor, err := state_processor.MakeGenesisAccountStateProcessor(fakeStateStorage{errKey: "_nonce"})
	assert.NoError(t, err)
	defaultGenesis.AccountStateProcessor = sProcessor
	block, err = defaultGenesis.Prepare()
	assert.Error(t, err)
	assert.Nil(t, block)

	// the contract owner balance isn't enough
	sProcessor, err = state_processor.MakeGenesisAccountStateProcessor(fakeStateStorage{})
	assert.NoError(t, err)
	defaultGenesis.AccountStateProcessor = sProcessor
	block, err = defaultGenesis.Prepare()
	assert.Error(t, err)
	assert.Nil(t, block)

	// set balance error
	sProcessor, err = state_processor.MakeGenesisAccountStateProcessor(fakeStateStorage{setErr: TrieError})
	assert.NoError(t, err)
	defaultGenesis.AccountStateProcessor = sProcessor
	block, err = defaultGenesis.Prepare()
	assert.Equal(t, TrieError, err)
	assert.Nil(t, block)
}

func TestSetupGenesisBlock_Error(t *testing.T) {
	defaultGenesis := createGenesis()

	rProcessor, err := registerdb.MakeGenesisRegisterProcessor(fakeStateStorage{setErr: TrieError})
	assert.NoError(t, err)
	defaultGenesis.RegisterProcessor = rProcessor

	_, blockHash, err := SetupGenesisBlock(defaultGenesis)
	assert.Equal(t, TrieError, err)
	assert.Equal(t, common.Hash{}, blockHash)
}

func TestGenesis_SetEarlyTokenContract_Error(t *testing.T) {
	defaultGenesis := createGenesis()

	// get balance error
	sProcessor, err := state_processor.MakeGenesisAccountStateProcessor(fakeStateStorage{getErr: TrieError})
	assert.NoError(t, err)
	defaultGenesis.AccountStateProcessor = sProcessor

	err = defaultGenesis.SetEarlyTokenContract()
	assert.Equal(t, g_error.AccountNotExist, err)

	log.Info("")
	// set balance error
	sProcessor, err = state_processor.MakeGenesisAccountStateProcessor(fakeStateStorage{
		setErr:          TrieError,
		contractBalance: big.NewInt(0).Mul(big.NewInt(43693128000000000), big.NewInt(consts.DIP)),
	})
	assert.NoError(t, err)
	defaultGenesis.AccountStateProcessor = sProcessor

	err = defaultGenesis.SetEarlyTokenContract()
	assert.Equal(t, g_error.AccountNotExist, err)
}
