package minemaster

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"sync/atomic"
	"testing"
)

func TestMakeMineMaster(t *testing.T) {
	// build mock
	ctrl := gomock.NewController(t)
	blockBuilderMock := NewMockBlockBuilder(ctrl)
	blockBroadcasterMock := NewMockBlockBroadcaster(ctrl)
	// build config
	mineConfig := MineConfig{
		GasFloor:         &atomic.Value{},
		GasCeil:          &atomic.Value{},
		CoinbaseAddress:  &atomic.Value{},
		BlockBuilder:     blockBuilderMock,
		BlockBroadcaster: blockBroadcasterMock,
	}
	// assert
	master, masterSvr := MakeMineMaster(mineConfig)
	assert.NotNil(t, master)
	assert.NotNil(t, masterSvr)
}
