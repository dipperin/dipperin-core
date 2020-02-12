package minemaster

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"math/big"
	"sync/atomic"
	"testing"
)

func Test_newMaster(t *testing.T) {
	mineConfig := MineConfig{}
	master := newMaster(mineConfig)
	assert.NotNil(t, master)
	assert.NotNil(t, master.MineConfig)
}

func Test_master_SetMsgSigner(t *testing.T) {
	// mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	signer := NewMockPbftSigner(ctrl)
	blockBuilder := NewMockBlockBuilder(ctrl)
	blockBuilder.EXPECT().SetMsgSigner(gomock.Any()).AnyTimes()
	blockBuilder.EXPECT().GetMsgSigner().Return(signer)
	// init
	mineConfig := MineConfig{
		BlockBuilder: blockBuilder,
	}
	getWorkersFunc := func() map[WorkerId]WorkerForMaster { return make(map[WorkerId]WorkerForMaster) }
	master := newMaster(mineConfig)
	master.workManager = newDefaultWorkManager(mineConfig)
	master.workDispatcher = newWorkDispatcher(mineConfig, getWorkersFunc)
	// test
	master.SetMsgSigner(signer)
	assert.NotNil(t, master.MineConfig.BlockBuilder.GetMsgSigner())
}

func Test_master_RetrieveReward(t *testing.T) {
	// mock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	blockBuilder := NewMockBlockBuilder(ctrl)
	blockBuilder.EXPECT().SetMsgSigner(gomock.Any()).AnyTimes()
	// init
	mineConfig := MineConfig{
		BlockBuilder: blockBuilder,
	}
	getWorkersFunc := func() map[WorkerId]WorkerForMaster { return make(map[WorkerId]WorkerForMaster) }
	master := newMaster(mineConfig)
	master.workManager = newDefaultWorkManager(mineConfig)
	master.workDispatcher = newWorkDispatcher(mineConfig, getWorkersFunc)
	// test
	master.RetrieveReward(common.HexToAddress("123"))
	assert.Equal(t, uint64(0), master.workManager.getPerformance(common.HexToAddress("123")))
	assert.Equal(t, big.NewInt(0), master.workManager.getReward(common.HexToAddress("123")))
}

func Test_master_GetReward(t *testing.T) {
	// init
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	addr := common.HexToAddress("123")
	// test case
	situations := []struct {
		name   string
		given  func() *master
		expect *big.Int
	}{
		{
			"default case",
			func() *master {
				mineConfig := MineConfig{}
				master := newMaster(mineConfig)
				master.workManager = newDefaultWorkManager(mineConfig)
				return master
			},
			big.NewInt(0),
		},
		// todo: can't set reward because of private interface and variable
	}
	// test
	for _, situation := range situations {
		master := situation.given()
		assert.Equal(t, situation.expect, master.GetReward(addr), situation.name)
	}
}

func Test_master_GetPerformance(t *testing.T) {
	// init
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	addr := common.HexToAddress("123")
	// test case
	situations := []struct {
		name   string
		given  func() *master
		expect uint64
	}{
		{
			"default case",
			func() *master {
				mineConfig := MineConfig{}
				master := newMaster(mineConfig)
				master.workManager = newDefaultWorkManager(mineConfig)
				return master
			},
			uint64(0),
		},
		// todo: can't set performance because of private interface and variable
	}
	// test
	for _, situation := range situations {
		master := situation.given()
		assert.Equal(t, situation.expect, master.GetPerformance(addr), situation.name)
	}
}

func Test_master_CurrentCoinbaseAddress(t *testing.T) {
	// init
	addr := common.HexToAddress("123")
	mineConfig := MineConfig{
		CoinbaseAddress: &atomic.Value{},
	}
	mineConfig.CoinbaseAddress.Store(addr)
	master := newMaster(mineConfig)
	// test
	assert.Equal(t, addr, master.CurrentCoinbaseAddress())
}

func Test_master_SetCoinbaseAddress(t *testing.T) {
	// init
	addr := common.HexToAddress("123")
	mineConfig := MineConfig{
		CoinbaseAddress:&atomic.Value{},
	}
	mineConfig.CoinbaseAddress.Store(addr)
	master := newMaster(mineConfig)
	assert.Equal(t, addr, master.MineConfig.CoinbaseAddress.Load())
	// test
	newAddr := common.HexToAddress("321")
	master.SetCoinbaseAddress(newAddr)
	assert.Equal(t, newAddr, master.MineConfig.CoinbaseAddress.Load())
}

func Test_master_SetMineGasConfig(t *testing.T) {
	// init
	defaultGasFloor := uint64(1)
	defaultGasCeil := uint64(999)
	mineConfig := MineConfig{
		GasFloor: &atomic.Value{},
		GasCeil:  &atomic.Value{},
	}
	mineConfig.GasFloor.Store(defaultGasFloor)
	mineConfig.GasCeil.Store(defaultGasCeil)
	master := newMaster(mineConfig)
	assert.Equal(t, defaultGasFloor, master.MineConfig.GasFloor.Load())
	assert.Equal(t, defaultGasCeil, master.MineConfig.GasCeil.Load())
	// test
	newGasFloor := uint64(5)
	newGasCeil := uint64(6660)
	master.SetMineGasConfig(newGasFloor, newGasCeil)
	assert.Equal(t, newGasFloor, master.MineConfig.GasFloor.Load())
	assert.Equal(t, newGasCeil, master.MineConfig.GasCeil.Load())
}
