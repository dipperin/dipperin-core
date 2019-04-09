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
	"github.com/stretchr/testify/assert"
	"runtime"
	"testing"
)

func TestNodeConfig_FullChainDBDir(t *testing.T) {
	nodeConfig := NodeConfig{}
	result := nodeConfig.FullChainDBDir()
	assert.Equal(t, "full_chain_data", result)
}

func TestNodeConfig_GetAllowHosts(t *testing.T) {
	nodeConfig := NodeConfig{AllowHosts: []string{}}
	result := nodeConfig.GetAllowHosts()
	assert.Equal(t, []string{}, result)
}

func TestNodeConfig_GetIsStartMine(t *testing.T) {
	nodeConfig := NodeConfig{}
	result := nodeConfig.GetIsStartMine()
	assert.Equal(t, false, result)
}

func TestNodeConfig_GetIsUploadNodeData(t *testing.T) {
	nodeConfig := NodeConfig{}
	result := nodeConfig.GetIsUploadNodeData()
	assert.Equal(t, 0, result)
}

func TestNodeConfig_GetNodeHTTPPort(t *testing.T) {
	nodeConfig := NodeConfig{}
	result := nodeConfig.GetNodeHTTPPort()
	assert.Equal(t, "0", result)
}

func TestNodeConfig_GetNodeName(t *testing.T) {
	nodeConfig := NodeConfig{}
	result := nodeConfig.GetNodeName()
	assert.Equal(t, "", result)
}

func TestNodeConfig_GetNodeP2PPort(t *testing.T) {
	nodeConfig := NodeConfig{}
	result := nodeConfig.GetNodeP2PPort()
	assert.Equal(t, "", result)
}

func TestNodeConfig_GetNodeType(t *testing.T) {
	nodeConfig := NodeConfig{}
	result := nodeConfig.GetNodeType()
	assert.Equal(t, 0, result)
}

func TestNodeConfig_GetPMetricsPort(t *testing.T) {
	nodeConfig := NodeConfig{}
	result := nodeConfig.GetPMetricsPort()
	assert.Equal(t, 0, result)
}

func TestNodeConfig_GetUploadURL(t *testing.T) {
	nodeConfig := NodeConfig{}
	result := nodeConfig.GetUploadURL()
	assert.Equal(t, "", result)
}

func TestNodeConfig_HttpEndpoint(t *testing.T) {
	nodeConfig := NodeConfig{}
	result := nodeConfig.HttpEndpoint()
	assert.Equal(t, "", result)

	nodeConfig = NodeConfig{HTTPHost: "host"}
	result = nodeConfig.HttpEndpoint()
	assert.Equal(t, "host:0", result)
}

func TestNodeConfig_IpcEndpoint(t *testing.T) {
	nodeConfig := NodeConfig{}
	result := nodeConfig.IpcEndpoint()
	assert.Equal(t, "", result)

	nodeConfig = NodeConfig{
		IPCPath: "path",
	}
	result = nodeConfig.IpcEndpoint()
	if runtime.GOOS == "windows" {
		assert.Equal(t, `\\.\pipe\`+nodeConfig.IPCPath, result)
	} else {
		assert.Equal(t, "/tmp/path", result)
	}

	nodeConfig = NodeConfig{
		IPCPath: "path",
		DataDir: "dir",
	}
	result = nodeConfig.IpcEndpoint()
	if runtime.GOOS == "windows" {
		assert.Equal(t, `\\.\pipe\`+nodeConfig.IPCPath, result)
	} else {
		assert.Equal(t, "dir/path", result)
	}

}

func TestNodeConfig_LightChainDBDir(t *testing.T) {
	nodeConfig := NodeConfig{}
	result := nodeConfig.LightChainDBDir()
	assert.Equal(t, "light_chain_data", result)
}

func TestNodeConfig_SoftWalletDir(t *testing.T) {
	nodeConfig := NodeConfig{}
	result := nodeConfig.SoftWalletDir()
	assert.Equal(t, "", result)
}

func TestNodeConfig_SoftWalletFile(t *testing.T) {
	nodeConfig := NodeConfig{}
	result := nodeConfig.SoftWalletFile()
	assert.Equal(t, "CSWallet", result)

	nodeConfig = NodeConfig{SoftWalletPath: "path"}
	result = nodeConfig.SoftWalletFile()
	assert.Equal(t, "path", result)
}

func TestNodeConfig_SoftWalletName(t *testing.T) {
	nodeConfig := NodeConfig{}
	result := nodeConfig.SoftWalletName()
	assert.Equal(t, "CSWallet", result)
}

func TestNodeConfig_WsEndpoint(t *testing.T) {
	nodeConfig := NodeConfig{}
	result := nodeConfig.WsEndpoint()
	assert.Equal(t, "", result)

	nodeConfig = NodeConfig{WSHost: "host"}
	result = nodeConfig.WsEndpoint()
	assert.Equal(t, "host:0", result)
}
