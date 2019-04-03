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

package dipperin

import (
	"testing"
	"github.com/dipperin/dipperin-core/core/cs-chain/chain-writer"
	"github.com/dipperin/dipperin-core/core/cs-chain/chain-state"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/tests"
	"github.com/stretchr/testify/assert"
	"os"
	"github.com/dipperin/dipperin-core/common/util"
	"path/filepath"
	"github.com/dipperin/dipperin-core/tests/factory"
)

func createCsChain() *chain_state.ChainState {
	f := chain_writer.NewChainWriterFactory()
	chainState := chain_state.NewChainState(&chain_state.ChainStateConfig{
		DataDir:       "",
		WriterFactory: f,
		ChainConfig:   chain_config.GetChainConfig(),
	})
	f.SetChain(chainState)

	// Mainly the default initial verifier is in it, the outer test calls it to vote on the block
	tests.NewGenesisEnv(chainState.GetChainDB(), chainState.GetStateStorage(), nil)
	return chainState
}

func TestDefault(t *testing.T) {
	assert.NotNil(t, DefaultDataDir())
	assert.NotNil(t, DefaultMinerP2PConf())
	assert.NotNil(t, DefaultNodeConf())
	assert.NotNil(t, DefaultP2PConf())
}

func TestChainVerifiersReader(t *testing.T) {
	csChain := createCsChain()
	verifierNum := int(csChain.ChainConfig.VerifierNumber)
	reader := MakeVerifiersReader(csChain)

	verifiers := reader.CurrentVerifiers()
	assert.Len(t, verifiers, verifierNum)

	verifiers = reader.NextVerifiers()
	assert.Len(t, verifiers, verifierNum)

	verifier := reader.GetPBFTPrimaryNode()
	assert.Equal(t, verifiers[0], verifier)

	verifier = reader.PrimaryNode()
	assert.Equal(t, verifiers[0], verifier)

	assert.False(t, reader.ShouldChangeVerifier())
	assert.Equal(t, verifierNum, reader.VerifiersTotalCount())

	specialBlock := factory.CreateSpecialBlock(1)
	csChain.ChainDB.SaveBlock(specialBlock)
	csChain.ChainDB.SaveHeadBlockHash(specialBlock.Hash())

	verifier = reader.GetPBFTPrimaryNode()
	assert.Equal(t, verifiers[0], verifier)
}

func Test_getNodeList(t *testing.T) {
	path := filepath.Join("/tmp", "dipperin_Test_getNodeList")
	defer os.RemoveAll(path)

	file, err := os.Create(path)
	assert.NoError(t, err)

	_, err = file.Write(util.StringifyJsonToBytes([]string{url, "", "node"}))
	assert.NoError(t, err)

	NodeList := getNodeList(path)
	assert.NotNil(t, NodeList)

	NodeList = getNodeList(util.HomeDir())
	assert.Nil(t, NodeList)
}

func Test_loadNodeKeyFromFile(t *testing.T) {
	assert.NotNil(t, loadNodeKeyFromFile(""))
	defer os.RemoveAll("nodekey")
}