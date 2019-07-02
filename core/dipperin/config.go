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
	"fmt"
	"github.com/dipperin/dipperin-core/core/dipperin/service"
	"github.com/dipperin/dipperin-core/third-party/rpc"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type ExtraServiceFunc func(c ExtraServiceFuncConfig) (apis []rpc.API, services []NodeService)
type ExtraServiceFuncConfig struct {
	service.DipperinConfig

	ChainService *service.MercuryFullChainService
}

// start a dipperin node config
// node config only in here, other module should use interface to use this config
type NodeConfig struct {
	// cur node name
	Name string

	P2PListener string

	IPCPath string
	// app data dir
	DataDir string
	// HTTPHost is the host interface on which to start the HTTP RPC server. If this
	// field is empty, no HTTP API endpoint will be started.
	HTTPHost string `toml:",omitempty"`

	// HTTPPort is the TCP port number on which to start the HTTP RPC server. The
	// default zero value is/ valid and will pick a port number randomly (useful
	// for ephemeral nodes).
	HTTPPort int `toml:",omitempty"`
	// WSHost is the host interface on which to start the websocket RPC server. If
	// this field is empty, no websocket API endpoint will be started.
	WSHost string `toml:",omitempty"`

	// WSPort is the TCP port number on which to start the websocket RPC server. The
	// default zero value is/ valid and will pick a port number randomly (useful for
	// ephemeral nodes).
	WSPort int `toml:",omitempty"`

	// 0 normal 1 mine master 2 verifier
	NodeType int

	// Set debug mode, 0 is single node not broadcast, 1 is multi-node with PBFT, 2 is multi-node with PBFT and election
	DebugMode int

	SoftWalletPassword   string
	SoftWalletPassPhrase string
	SoftWalletPath       string
	IsStartMine          bool

	//used to set the default account of pbft
	DefaultAccountKey string

	// Whether it includes browser function. 0 is not, 1 is
	IsScanner int

	IsUploadNodeData int

	//whether tps monitoring function 0 is not 1 is
	IsPerformance int
	// example http://127.0.0.1:8080
	UploadURL string

	NoDiscovery int
	Nat         string

	AllowHosts []string

	PMetricsPort int

	ExtraServiceFunc ExtraServiceFunc
}

func (conf NodeConfig) GetIsStartMine() bool {
	return conf.IsStartMine
}

func (conf NodeConfig) GetPMetricsPort() int {
	return conf.PMetricsPort
}

func (conf NodeConfig) GetAllowHosts() []string {
	return conf.AllowHosts
}

func (conf NodeConfig) GetIsUploadNodeData() int {
	return conf.IsUploadNodeData
}

func (conf NodeConfig) GetUploadURL() string {
	//switch os.Getenv(chain_config.BootEnvTagName) {
	//case "test":
	//	log.Agent("use test upload url for monitor")
	//	return fmt.Sprintf("http://%v:8887/api/Dipperin_nodes", chain_config.TestServer)
	// Mercury is configured directly through the startup parameters
	//case "mercury":
	//	log.Agent("use mercury upload url for monitor")
	//	return fmt.Sprintf("http://%v:8887/api/dipperin_nodes", chain_config.TestServer)
	//}
	return conf.UploadURL
}

func (conf NodeConfig) GetNodeName() string {
	return conf.Name
}

func (conf NodeConfig) GetNodeP2PPort() string {
	return conf.P2PListener
}

func (conf NodeConfig) GetNodeHTTPPort() string {
	return strconv.Itoa(conf.HTTPPort)
}

func (conf NodeConfig) GetNodeType() int {
	return conf.NodeType
}

func (conf NodeConfig) SoftWalletName() string {
	return "CSWallet"
}

func (conf NodeConfig) SoftWalletDir() string {
	return conf.DataDir
}

func (conf NodeConfig) SoftWalletFile() string {
	if conf.SoftWalletPath == "" {
		return filepath.Join(conf.SoftWalletDir(), conf.SoftWalletName())
	} else {
		return conf.SoftWalletPath
	}
}

// full and fast use same dir, fast to full should start sync history states
func (conf NodeConfig) FullChainDBDir() string {
	return filepath.Join(conf.DataDir, "full_chain_data")
}

func (conf NodeConfig) LightChainDBDir() string {
	return filepath.Join(conf.DataDir, "light_chain_data")
}

func (conf NodeConfig) IpcEndpoint() string {
	// Short circuit if IPC has not been enabled
	if conf.IPCPath == "" {
		return ""
	}
	// On windows we can only use plain top-level pipes
	if runtime.GOOS == "windows" {
		if strings.HasPrefix(conf.IPCPath, `\\.\pipe\`) {
			return conf.IPCPath
		}
		return `\\.\pipe\` + conf.IPCPath
	}
	// Resolve names into the data directory full paths otherwise
	if filepath.Base(conf.IPCPath) == conf.IPCPath {
		if conf.DataDir == "" {
			return filepath.Join(os.TempDir(), conf.IPCPath)
		}
		return filepath.Join(conf.DataDir, conf.IPCPath)
	}
	return conf.IPCPath
}
func (conf NodeConfig) HttpEndpoint() string {
	if conf.HTTPHost == "" {
		return ""
	}
	return fmt.Sprintf("%s:%d", conf.HTTPHost, conf.HTTPPort)
}
func (conf NodeConfig) WsEndpoint() string {
	if conf.WSHost == "" {
		return ""
	}
	return fmt.Sprintf("%s:%d", conf.WSHost, conf.WSPort)
}
