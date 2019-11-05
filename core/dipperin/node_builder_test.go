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
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/core/cs-chain"
	"github.com/dipperin/dipperin-core/tests/factory"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/third-party/rpc"
)

var (
	aliceAddr = common.HexToAddress("0x00005586B883Ec6dd4f8c26063E18eb4Bd228e59c3E9")
)

func createNodeConfig() *NodeConfig {
	return &NodeConfig{
		Name:                 "test",
		P2PListener:          ":30000",
		IPCPath:              "path",
		DataDir:              filepath.FromSlash(util.HomeDir() + "/dir"),
		HTTPHost:             "127.0.0.1",
		HTTPPort:             3335,
		WSHost:               "127.0.0.1",
		WSPort:               4335,
		NodeType:             0,
		SoftWalletPassword:   "123",
		SoftWalletPassPhrase: "pass_phrase",
		IsStartMine:          true,
		DefaultAccountKey:    "account_key",
		UploadURL:            "http://127.0.0.1:8080",
		NoDiscovery:          1,
		ExtraServiceFunc: ExtraServiceFunc(func(c ExtraServiceFuncConfig) (apis []rpc.API, services []NodeService) {
			return
		}),
	}
}

func TestNewBftNode(t *testing.T) {
	if runtime.GOOS != "linux" {
		return
	}

	cs_chain.GenesisSetUp = true
	nodeConfig := createNodeConfig()
	node := NewBftNode(*nodeConfig)
	assert.NotNil(t, node)
	os.RemoveAll(nodeConfig.DataDir)

	nodeConfig.NodeType = 1
	node = NewBftNode(*nodeConfig)
	assert.NotNil(t, node)
	os.RemoveAll(nodeConfig.DataDir)

	os.Setenv("boots_env", "mercury")
	nodeConfig.NodeType = 2
	node = NewBftNode(*nodeConfig)
	assert.NotNil(t, node)
	os.RemoveAll(nodeConfig.DataDir)
}

func TestMsgSender(t *testing.T) {
	if runtime.GOOS != "linux" {
		return
	}

	cs_chain.GenesisSetUp = true
	os.Setenv("boots_env", "test")

	nodeConfig := createNodeConfig()
	nodeConfig.Nat = "pmp:192.168.0.1"

	base := newBaseComponent(*nodeConfig)
	base.initFullChain()
	base.initP2PService()

	msgSender := MsgSender{
		csPm:              base.csPm,
		broadcastDelegate: base.broadcastDelegate,
	}
	block := factory.CreateBlock(1)
	msgSender.BroadcastEiBlock(block)
	msgSender.BroadcastMsg(0, block.Hash())
	msgSender.SendReqRoundMsg(0, []common.Address{aliceAddr}, block.Hash())
	os.RemoveAll(nodeConfig.DataDir)
}

type TestService struct {
}

func (t *TestService) Start() error {
	return nil
}

func (t *TestService) Stop() {
	return
}

func TestGetServiceName(t *testing.T) {
	var service NodeService

	service = &TestService{}

	serviceType := reflect.TypeOf(service)

	log.Info("the service name is:", "serviceType", serviceType.String())

	log.Info("the a type name is:", "name", reflect.TypeOf("hello").Name())
}
