package dipperin

import (
	"github.com/dipperin/dipperin-core/core/chaincommunication"
	"github.com/dipperin/dipperin-core/core/chainconfig"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultNodeConf(t *testing.T) {
	conf := DefaultNodeConf()
	assert.Equal(t, chainconfig.AppName, conf.Name)
	assert.Equal(t, "/tmp/dipperin.ipc", conf.IPCPath)
	assert.Equal(t, "127.0.0.1", conf.HTTPHost)
	assert.Equal(t, 7777, conf.HTTPPort)
	assert.Equal(t, "127.0.0.1", conf.WSHost)
	assert.Equal(t, 8888, conf.WSPort)
	assert.Equal(t, 0, conf.IsUploadNodeData)
	assert.Equal(t, "", conf.UploadURL)
}

func TestDefaultP2PConf(t *testing.T) {
	conf := DefaultP2PConf()
	assert.Equal(t, false, conf.NoDiscovery)
	assert.Equal(t, chaincommunication.P2PMaxPeerCount, conf.MaxPeers)
	assert.Equal(t, ":60606", conf.ListenAddr)
}

func TestDefaultMinerP2PConf(t *testing.T) {
	conf := DefaultMinerP2PConf()
	assert.Equal(t, true, conf.NoDiscovery)
	assert.Equal(t, chaincommunication.P2PMaxPeerCount, conf.MaxPeers)
	assert.Equal(t, ":68080", conf.ListenAddr)
}
