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

package chain_config

import (
	"fmt"
	"net"
	"os"
	"testing"
	"github.com/dipperin/dipperin-core/common/hexutil"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/core/chain-config/env-conf"
	"github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
	"github.com/dipperin/dipperin-core/third-party/p2p/enode"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"io/ioutil"
)

func TestLoadDefaultBootID(t *testing.T) {
	if GetChainConfig().VerifierNumber != 22 {
		panic("verifier num is not 22")
	}

	pkByte, err := hexutil.Decode(env_conf.MercuryBootNodePublicKey)
	assert.NoError(t, err)

	pk := cs_crypto.ToECDSAPub(pkByte)

	n := enode.NewV4(pk, net.ParseIP("14.17.65.122"), 30301, 30301)
	fmt.Println(n.String())
}

func TestGetChainConfig(t *testing.T) {
	chainConfig := GetChainConfig()
	assert.Equal(t, uint64(110), chainConfig.SlotSize)
	assert.Equal(t, 22, chainConfig.VerifierNumber)

	err := os.Setenv("boots_env", "test")
	assert.NoError(t, err)

	chainConfig = defaultChainConfig()
	assert.Equal(t, uint64(1), chainConfig.NetworkID)

	err = os.Setenv("boots_env", "mercury")
	assert.NoError(t, err)

	chainConfig = defaultChainConfig()
	assert.Equal(t, uint64(99), chainConfig.NetworkID)
}

func TestGetCurBootsEnv(t *testing.T) {
	err := os.Setenv("boots_env", "mercury")
	assert.NoError(t, err)

	env := GetCurBootsEnv()
	assert.Equal(t, "mercury", env)
}

func TestDefaultDataDir(t *testing.T) {
	dir := DefaultDataDir()
	assert.NotEqual(t, "", dir)
}

func TestInitBootNodes(t *testing.T) {
	if !util.IsTestEnv() {
		return
	}

	err := os.Setenv("boots_env", "")
	assert.NoError(t, err)
	InitBootNodes("")
	assert.Equal(t, 1, len(VerifierBootNodes))
	assert.Equal(t, 1, len(KBucketNodes))

	err = os.Setenv("boots_env", "test")
	assert.NoError(t, err)
	resetNodes()
	InitBootNodes("")
	assert.Equal(t, 4, len(VerifierBootNodes))
	assert.Equal(t, 1, len(KBucketNodes))

	err = os.Setenv("boots_env", "mercury")
	assert.NoError(t, err)
	VerifierBootNodes = []*enode.Node{}
	KBucketNodes = []*enode.Node{}
	InitBootNodes("")
	assert.Equal(t, 4, len(VerifierBootNodes))
	assert.Equal(t, 1, len(KBucketNodes))
}

func TestLoadNodesFromFile(t *testing.T) {
	gFPath := filepath.Join(util.HomeDir(), "test.txt")

	// No error
	url := []string{"enode://8b610c5400bfdb355c7a204beb65cb261fe5e89cc2c1837dc3cf752d16df65cfd95c6ee79be3720ef9dc6ba0b6876c63530a6352cb18298afb2b282b111ec7cf@192.168.122.102:40006"}
	json := util.StringifyJsonToBytes(url)
	err := ioutil.WriteFile(gFPath, json, 0666)
	assert.NoError(t, err)

	nodes := LoadNodesFromFile(gFPath)
	assert.Equal(t, url[0], nodes[0].String())

	err = os.Remove(gFPath)
	assert.NoError(t, err)

	// ParseV4 error
	url = []string{"enode://@192.168.122.102:40006"}
	json = util.StringifyJsonToBytes(url)
	err = ioutil.WriteFile(gFPath, json, 0666)
	assert.NoError(t, err)

	nodes = LoadNodesFromFile(gFPath)
	assert.Len(t, nodes, 0)

	err = os.Remove(gFPath)
	assert.NoError(t, err)

	// ParseJsonFromBytes error
	json = []byte{123}
	err = ioutil.WriteFile(gFPath, json, 0666)
	assert.NoError(t, err)

	nodes = LoadNodesFromFile(gFPath)
	assert.Len(t, nodes, 0)

	err = os.Remove(gFPath)
	assert.NoError(t, err)
}

func resetNodes() {
	VerifierBootNodes = []*enode.Node{}
	KBucketNodes = []*enode.Node{}
}

func Test_initMercuryBoots(t *testing.T) {
	initMercuryBoots("")
}
