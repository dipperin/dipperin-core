package minemaster

import (
	"github.com/dipperin/dipperin-core/tests/mock/mine/minemaster-mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"sync/atomic"
	"testing"
)

func TestMakeMineMaster(t *testing.T) {
	// build mock
	ctrl := gomock.NewController(t)
	blockBuilderMock := minemaster_mock.NewMockBlockBuilder(ctrl)
	blockBroadcasterMock := minemaster_mock.NewMockBlockBroadcaster(ctrl)
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
