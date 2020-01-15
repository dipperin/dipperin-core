package dipperin

import (
	"github.com/dipperin/dipperin-core/common/gerror"
	"github.com/stretchr/testify/assert"
	"os"
	"runtime"
	"testing"
)

func TestNodeConfig_NodeConfigCheck(t *testing.T) {
	situations := []struct{
		name string
		given func() *NodeConfig
		expectErr error
	}{
		{
			"no wallet start and not empty soft wallet path",
			func() *NodeConfig {
				return &NodeConfig{
					NoWalletStart:true,
					SoftWalletPath:"test",
				}
			},
			gerror.NodeConfWalletError,
		},
		{
			"no wallet start and not empty soft wallet password",
			func() *NodeConfig {
				return &NodeConfig{
					NoWalletStart:true,
					SoftWalletPath:"",
					SoftWalletPassword:"test",
				}
			},
			gerror.NodeConfWalletError,
		},
		{
			"no wallet start and not empty soft wallet pass phrase",
			func() *NodeConfig {
				return &NodeConfig{
					NoWalletStart:true,
					SoftWalletPath:"",
					SoftWalletPassword:"",
					SoftWalletPassPhrase:"test",
				}
			},
			gerror.NodeConfWalletError,
		},
		{
			"no wallet start and empty soft wallet params",
			func() *NodeConfig {
				return &NodeConfig{
					NoWalletStart:true,
					SoftWalletPath:"",
					SoftWalletPassword:"",
					SoftWalletPassPhrase:"",
				}
			},
			nil,
		},
		{
			"wallet start and empty soft wallet password",
			func() *NodeConfig {
				return &NodeConfig{
					NoWalletStart:false,
					SoftWalletPassword:"",
				}
			},
			gerror.NodeConfWalletError,
		},
		{
			"wallet start and not empty soft wallet password",
			func() *NodeConfig {
				return &NodeConfig{
					NoWalletStart:false,
					SoftWalletPassword:"test",
				}
			},
			nil,
		},
	}
	// test
	for _, situation := range situations {
		config := situation.given()
		assert.Equal(t, situation.expectErr, config.NodeConfigCheck(), situation.name)
	}
}

func TestNodeConfig_GetIsStartMine(t *testing.T) {
	config := &NodeConfig{IsStartMine:true}
	assert.Equal(t, true, config.GetIsStartMine())
}

func TestNodeConfig_GetPMetricsPort(t *testing.T) {
	config := &NodeConfig{PMetricsPort:3333}
	assert.Equal(t, 3333, config.GetPMetricsPort())
}

func TestNodeConfig_GetAllowHosts(t *testing.T) {
	config := &NodeConfig{AllowHosts:[]string{"127.0.0.1"}}
	assert.Equal(t, []string{"127.0.0.1"}, config.GetAllowHosts())
}

func TestNodeConfig_GetIsUploadNodeData(t *testing.T) {
	config := &NodeConfig{IsUploadNodeData:123}
	assert.Equal(t, 123, config.GetIsUploadNodeData())
}

func TestNodeConfig_GetUploadURL(t *testing.T) {
	config := &NodeConfig{UploadURL:"test url"}
	assert.Equal(t, "test url", config.GetUploadURL())
}

func TestNodeConfig_GetNodeName(t *testing.T) {
	config := &NodeConfig{Name:"test name"}
	assert.Equal(t, "test name", config.GetNodeName())
}

func TestNodeConfig_GetNodeP2PPort(t *testing.T) {
	config := &NodeConfig{P2PListener:"3333"}
	assert.Equal(t, "3333", config.GetNodeP2PPort())
}

func TestNodeConfig_GetNodeHTTPPort(t *testing.T) {
	config := &NodeConfig{HTTPPort:3333}
	assert.Equal(t, "3333", config.GetNodeHTTPPort())
}

func TestNodeConfig_GetNodeType(t *testing.T) {
	config := &NodeConfig{NodeType:3333}
	assert.Equal(t, 3333, config.GetNodeType())
}

func TestNodeConfig_SoftWalletName(t *testing.T) {
	config := &NodeConfig{}
	assert.Equal(t, "CSWallet", config.SoftWalletName())
}

func TestNodeConfig_SoftWalletDir(t *testing.T) {
	config := &NodeConfig{DataDir:"test data dir"}
	assert.Equal(t, "test data dir", config.SoftWalletDir())
}
func TestNodeConfig_SoftWalletFile(t *testing.T) {
	situations := []struct{
		name string
		given func() *NodeConfig
		expect string
	}{
		{
			"empty soft wallet path",
			func() *NodeConfig {
				config := &NodeConfig{
					DataDir:"datadir",
					SoftWalletPath:"",
				}
				return config
			},
			"datadir/CSWallet",
		},
		{
				"not empty soft wallet path",
				func() *NodeConfig {
					config := &NodeConfig{SoftWalletPath:"test"}
					return config
				},
				"test",
		},
	}
	// test
	for _, situation := range situations {
		config := situation.given()
		assert.Equal(t, situation.expect, config.SoftWalletFile(), situation.name)
	}
}
func TestNodeConfig_FullChainDBDir(t *testing.T) {
	config := &NodeConfig{DataDir:"dir"}
	assert.Equal(t, "dir/full_chain_data", config.FullChainDBDir())
}
func TestNodeConfig_LightChainDBDir(t *testing.T) {
	config := &NodeConfig{DataDir:"dir"}
	assert.Equal(t, "dir/light_chain_data", config.LightChainDBDir())
}
func TestNodeConfig_IpcEndpoint(t *testing.T) {
	type situationStc struct {
		name string
		given func() *NodeConfig
		expect string
	}
	situations := []situationStc{
		{
			"empty IPCPath",
			func() *NodeConfig {
				config := &NodeConfig{IPCPath:""}
				return config
			},
			"",
		},
	}
	// append test case with different GOOS
	if runtime.GOOS == "windows" {
		cases := []situationStc{
			{
				"contain pipe prefix",
				func() *NodeConfig {
					config := &NodeConfig{IPCPath: `path\pipe\prefix`}
					return config
				},
				`\path\pipe\prefix`,
			},
			{
					"not contain pipe prefix",
					func() *NodeConfig {
						config := &NodeConfig{IPCPath: `path\any\prefix`}
						return config
					},
					`\\.\pipe\path\any\prefix`,
			},
		}
		situations = append(situations, cases...)
	} else {
		cases := []situationStc{
			{
				"base not equal",
				func() *NodeConfig {
					config := &NodeConfig{IPCPath:"test/abc"}
					return config
				},
				"test/abc",
			},
			{
					"base equal and empty data dir",
					func() *NodeConfig {
						config := &NodeConfig{
							IPCPath:"test",
							DataDir:"",
						}
						return config
					},
					os.TempDir() + "/" + "test",
			},
			{
						"base equal and not empty data dir",
						func() *NodeConfig {
							config := &NodeConfig{
								IPCPath:"test",
								DataDir: "datadir",
							}
							return config
						},
						"datadir/test",
			},
		}
		situations = append(situations, cases...)
	}
	// test
	for _, situation := range situations {
		config := situation.given()
		assert.Equal(t, situation.expect, config.IpcEndpoint(), situation.name)
	}
}

func TestNodeConfig_HttpEndpoint(t *testing.T) {
	config := &NodeConfig{HTTPHost:"", HTTPPort:3333}
	assert.Equal(t, "", config.HttpEndpoint())
	config.HTTPHost = "127.0.0.1"
	assert.Equal(t, "127.0.0.1:3333", config.HttpEndpoint())
}

func TestNodeConfig_WsEndpoint(t *testing.T) {
	// not empty host case
	config := &NodeConfig{WSHost:"", WSPort:3333}
	assert.Equal(t, "", config.WsEndpoint())
	// empty host case
	config.WSHost = "127.0.0.1"
	assert.Equal(t, "127.0.0.1:3333", config.WsEndpoint())
}