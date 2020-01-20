package dipperin

import (
	"errors"
	"github.com/dipperin/dipperin-core/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewMinerNode(t *testing.T) {
	situations := []struct {
		name      string
		given     func() (string, string, int, string)
		expectNil bool
		expectErr error
	}{
		{
			"empty coinbase",
			func() (master string, coinbase string, minerCount int, p2pListenAddr string) {
				master = "test master"
				coinbase = ""
				minerCount = 15
				p2pListenAddr = "test addr"
				return
			},
			true,
			errors.New("coinbase count not right"),
		},
		{
			"zero miner count",
			func() (master string, coinbase string, minerCount int, p2pListenAddr string) {
				master = "test master"
				coinbase = "test coin base"
				minerCount = 0
				p2pListenAddr = "test addr"
				return
			},
			true,
			errors.New("miner count not right"),
		},
		{
			"parse master node failed",
			func() (master string, coinbase string, minerCount int, p2pListenAddr string) {
				master = "test master"
				coinbase = "test coin base"
				minerCount = 15
				p2pListenAddr = "test addr"
				return
			},
			true,
			errors.New("parse master node faield"),
		},
		{
			"normal case",
			func() (master string, coinbase string, minerCount int, p2pListenAddr string) {
				master = "enode://b832f4f2fe19dbc5604766bbb268a6d0f7ce9ce381b034b262a92f0ad8283a1b5fa058dea5269b66fbb2014a24fa7198c6dc2d8c9cbac7a348258fc20702561f@127.0.0.1:10003"
				coinbase = "coinbase"
				minerCount = 1
				p2pListenAddr = common.HexToAddress("0x00005586B883Ec6dd4f8c26063E18eb4Bd228e59c3E9").String()
				return
			},
			false,
			nil,
		},
	}
	// test
	for _, situation := range situations {
		master, coinbase, minerCount, p2pListenAddr := situation.given()
		node, err := NewMinerNode(master, coinbase, minerCount, p2pListenAddr)
		// check result
		if situation.expectNil == true {
			assert.Nil(t, node)
		} else {
			assert.NotNil(t, node)
		}
		if situation.expectErr != nil {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}
