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
	"errors"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/consts"
	"github.com/dipperin/dipperin-core/common/gerror"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/core/chain/registerdb"
	"github.com/dipperin/dipperin-core/core/chain/stateprocessor"
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

	type result struct {
		err error
	}
	testCases := []struct {
		name   string
		given  func() error
		expect result
	}{
		{
			name:"genesis is nil",
			given: func() error {
				_, _, err := SetupGenesisBlock(nil)
				return err
			},
			expect:result{errors.New("genesis can't be nil")},
		},
		{
			name:"setup genesis successful",
			given: func() error {
				defaultGenesis.Difficulty = common.Difficulty{}
				_, _, err := SetupGenesisBlock(defaultGenesis)
				return err
			},
			expect:result{nil},
		},
		{
			name:"genesis and stored genesis is match",
			given: func() error {
				genesisBlock, _ := defaultGenesis.Prepare()
				defaultGenesis.ChainDB.SaveBlockHash(genesisBlock.Hash(), genesisBlock.Number())
				_, _, err = SetupGenesisBlock(defaultGenesis)
				return err
			},
			expect:result{nil},
		},
		{
			name:"config is nil",
			given: func() error {
				defaultGenesis.Config = nil
				_, _, err = SetupGenesisBlock(defaultGenesis)
				return err
			},
			expect:result{errGenesisNoConfig},
		},
	}

	for _,tc:=range testCases{
		err:=tc.given()
		if err!=nil{
			assert.Equal(t,tc.expect.err,err)
		}else {
			assert.NoError(t,err)
		}
	}

	genesisBlock, err := defaultGenesis.Prepare()
	block := createBlock(0)
	defaultGenesis.ChainDB.SaveBlockHash(block.Hash(), block.Number())
	_, blockHash, err := SetupGenesisBlock(defaultGenesis)
	assert.NotEqualf(t, genesisBlock.Hash(), blockHash,"not equal")
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

	type result struct {
		err error
	}
	testCases := []struct {
		name   string
		given  func() error
		expect result
	}{
		{
			name:"Register Processor",
			given: func() error {
				rProcessor, err := registerdb.MakeGenesisRegisterProcessor(fakeStateStorage{})
				defaultGenesis.RegisterProcessor = rProcessor
				err = defaultGenesis.Commit(createBlock(0))
				return err
			},
			expect:result{TrieError},
		},
		{
			name:"AccountState Processor",
			given: func() error {
				sProcessor, err := stateprocessor.MakeGenesisAccountStateProcessor(fakeStateStorage{})
				defaultGenesis.AccountStateProcessor = sProcessor
				err = defaultGenesis.Commit(createBlock(0))
				return err
			},
			expect:result{TrieError},
		},
	}

	for _,tc:=range testCases{
		err:=tc.given()
		if err!=nil{
			assert.Equal(t,tc.expect.err,err)
		}else {
			assert.NoError(t,err)
		}
	}
}

func TestGenesis_Prepare_Error(t *testing.T) {
	defaultGenesis := createGenesis()

	type result struct {
		err error
	}
	testCases := []struct {
		name   string
		given  func() error
		expect result
	}{
		{
			name:"prepare registerDB error",
			given: func() error {
				rProcessor, err := registerdb.MakeGenesisRegisterProcessor(fakeStateStorage{setErr: TrieError})
				defaultGenesis.RegisterProcessor = rProcessor
				_, err = defaultGenesis.Prepare()
				return err
			},
			expect:result{TrieError},
		},
		{
			name:"set balance error",
			given: func() error {
				sProcessor, err := stateprocessor.MakeGenesisAccountStateProcessor(fakeStateStorage{errKey: "_nonce"})
				defaultGenesis.AccountStateProcessor = sProcessor
				_, err = defaultGenesis.Prepare()
				return err
			},
			expect:result{errors.New("account does not exist")},
		},
		{
			name:"the contract owner balance isn't enough",
			given: func() error {
				sProcessor, err := stateprocessor.MakeGenesisAccountStateProcessor(fakeStateStorage{})
				defaultGenesis.AccountStateProcessor = sProcessor
				_, err = defaultGenesis.Prepare()
				return err
			},
			expect:result{errors.New("the contract owner balance isn't enough")},
		},
		{
			name:"Set TrieError",
			given: func() error {
				sProcessor, err := stateprocessor.MakeGenesisAccountStateProcessor(fakeStateStorage{setErr: TrieError})
				assert.NoError(t, err)
				defaultGenesis.AccountStateProcessor = sProcessor
				_, err = defaultGenesis.Prepare()
				return err
			},
			expect:result{TrieError},
		},
	}

	for _,tc:=range testCases{
		err:=tc.given()
		if err!=nil{
			assert.Equal(t,tc.expect.err,err)
		}else {
			assert.NoError(t,err)
		}
	}
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

	type result struct {
		err error
	}
	testCases := []struct {
		name   string
		given  func() error
		expect result
	}{
		{
			name:"get balance error",
			given: func() error {
				sProcessor, err := stateprocessor.MakeGenesisAccountStateProcessor(fakeStateStorage{getErr: TrieError})
				defaultGenesis.AccountStateProcessor = sProcessor
				err = defaultGenesis.SetEarlyTokenContract()
				return err
			},
			expect:result{gerror.ErrAccountNotExist},
		},
		{
			name:"set balance error",
			given: func() error {
				sProcessor, err := stateprocessor.MakeGenesisAccountStateProcessor(fakeStateStorage{
					setErr:          TrieError,
					contractBalance: big.NewInt(0).Mul(big.NewInt(43693128000000000), big.NewInt(consts.DIP)),
				})
				defaultGenesis.AccountStateProcessor = sProcessor
				err = defaultGenesis.SetEarlyTokenContract()
				return err
			},
			expect:result{gerror.ErrAccountNotExist},
		},
	}

	for _,tc:=range testCases{
		err:=tc.given()
		if err!=nil{
			assert.Equal(t,tc.expect.err,err)
		}else {
			assert.NoError(t,err)
		}
	}
}