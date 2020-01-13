package minemaster

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/tests/mock/mine/minemaster-mock"
	"github.com/dipperin/dipperin-core/tests/mock/model-mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func Test_defaultWorkManager_subtractPerformance(t *testing.T) {
	// init
	address := common.HexToAddress("123")
	blocksMined := uint64(10)
	// build
	mineConfig := MineConfig{
		GasFloor:         nil,
		GasCeil:          nil,
		CoinbaseAddress:  nil,
		BlockBuilder:     nil,
		BlockBroadcaster: nil,
	}
	manager := newDefaultWorkManager(mineConfig)
	mPerformance := make(map[common.Address]workerPerformance)
	mPerformance[address] = &defaultPerformance{blocksMined}
	manager.performance = mPerformance
	// test
	manager.subtractPerformance(address, blocksMined - 1)
	assert.Equal(t, uint64(1), manager.performance[address].getPerformance())
}

func Test_defaultWorkManager_subtractReward(t *testing.T) {
	// init
	address := common.HexToAddress("123")
	reward := big.NewInt(10)
	// build
	mineConfig := MineConfig{
		GasFloor:         nil,
		GasCeil:          nil,
		CoinbaseAddress:  nil,
		BlockBuilder:     nil,
		BlockBroadcaster: nil,
	}
	manager := newDefaultWorkManager(mineConfig)
	mReward := make(map[common.Address]*big.Int)
	mReward[address] = reward
	manager.reward = mReward
	// test
	paramInt := big.NewInt(9)
	manager.subtractReward(address, paramInt)
	assert.Equal(t, big.NewInt(1), manager.reward[address])
}

func Test_defaultWorkManager_clearPerformance(t *testing.T) {
	// init
	address1 := common.HexToAddress("123")
	address2 := common.HexToAddress("456")
	blocksMined := uint64(10)
	// build
	mineConfig := MineConfig{
		GasFloor:         nil,
		GasCeil:          nil,
		CoinbaseAddress:  nil,
		BlockBuilder:     nil,
		BlockBroadcaster: nil,
	}
	manager := newDefaultWorkManager(mineConfig)
	mPerformance := make(map[common.Address]workerPerformance)
	mPerformance[address1] = &defaultPerformance{blocksMined}
	mPerformance[address2] = &defaultPerformance{blocksMined}
	manager.performance = mPerformance
	// test
	assert.Equal(t, 2, len(manager.performance))
	manager.clearPerformance(address1)
	assert.Equal(t, 1, len(manager.performance))
}

func Test_defaultWorkManager_clearReward(t *testing.T) {
	// init
	address1 := common.HexToAddress("123")
	address2 := common.HexToAddress("456")
	reward := big.NewInt(10)
	// build
	mineConfig := MineConfig{
		GasFloor:         nil,
		GasCeil:          nil,
		CoinbaseAddress:  nil,
		BlockBuilder:     nil,
		BlockBroadcaster: nil,
	}
	manager := newDefaultWorkManager(mineConfig)
	mReward := make(map[common.Address]*big.Int)
	mReward[address1] = reward
	mReward[address2] = reward
	manager.reward = mReward
	// test
	assert.Equal(t, 2, len(manager.reward))
	manager.clearReward(address1)
	assert.Equal(t, 1, len(manager.reward))
}

func Test_defaultWorkManager_getReward(t *testing.T) {
	// init
	address1 := common.HexToAddress("123")
	reward := big.NewInt(10)
	// build
	mineConfig := MineConfig{
		GasFloor:         nil,
		GasCeil:          nil,
		CoinbaseAddress:  nil,
		BlockBuilder:     nil,
		BlockBroadcaster: nil,
	}
	manager := newDefaultWorkManager(mineConfig)
	mReward := make(map[common.Address]*big.Int)
	mReward[address1] = reward
	manager.reward = mReward
	// test
	assert.Equal(t, reward, manager.getReward(address1))
}

func Test_defaultWorkManager_getPerformance(t *testing.T) {
	// init
	address1 := common.HexToAddress("123")
	blocksMined := uint64(10)
	// build
	mineConfig := MineConfig{
		GasFloor:         nil,
		GasCeil:          nil,
		CoinbaseAddress:  nil,
		BlockBuilder:     nil,
		BlockBroadcaster: nil,
	}
	manager := newDefaultWorkManager(mineConfig)
	mPerformance := make(map[common.Address]workerPerformance)
	mPerformance[address1] = &defaultPerformance{blocksMined}
	manager.performance = mPerformance
	// test
	assert.Equal(t, blocksMined, manager.getPerformance(address1))
}

func Test_defaultWorkManager_divideReward(t *testing.T) {
	// mock and init
	address := common.HexToAddress("123")
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// build
	mineConfig := MineConfig{}
	manager := newDefaultWorkManager(mineConfig)
	// reward
	mReward := make(map[common.Address]*big.Int)
	mReward[address] = big.NewInt(10)
	manager.reward = mReward
	manager.totalReward = big.NewInt(10)
	// performance
	mPerformance := make(map[common.Address]workerPerformance)
	mPerformance[address] = &defaultPerformance{5}
	manager.performance = mPerformance
	// test
	resReward := manager.divideReward(big.NewInt(11))
	assert.Equal(t, big.NewInt(21), resReward[address])
	assert.Equal(t, big.NewInt(21), manager.reward[address])
}

func Test_defaultPerformance_getPerformance(t *testing.T) {
	// init
	performance := &defaultPerformance{blocksMined:10}
	assert.Equal(t, uint64(10), performance.getPerformance())
}

func Test_defaultPerformance_setPerformance(t *testing.T) {
	// init
	performance := &defaultPerformance{blocksMined:10}
	performance.setPerformance(15)
	assert.Equal(t, uint64(15), performance.blocksMined)
}

func Test_defaultPerformance_updatePerformance(t *testing.T) {
	// init
	performance := &defaultPerformance{blocksMined:10}
	performance.updatePerformance()
	assert.Equal(t, uint64(11), performance.blocksMined)
}

func Test_newDefaultPerformance(t *testing.T) {
	performance := newDefaultPerformance()
	assert.NotNil(t, performance)
}

func Test_newDefaultWorkManager(t *testing.T) {
	// init
	mineConfig := MineConfig{
		GasFloor:         nil,
		GasCeil:          nil,
		CoinbaseAddress:  nil,
		BlockBuilder:     nil,
		BlockBroadcaster: nil,
	}
	manager := newDefaultWorkManager(mineConfig)
	assert.NotNil(t, manager)
}

func Test_defaultWorkManager_submitBlock(t *testing.T) {
	// mock and init
	address := common.HexToAddress("123")
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// mock block broadcaster
	mockBlockBroadcaster := minemaster_mock.NewMockBlockBroadcaster(ctrl)
	mockBlockBroadcaster.EXPECT().BroadcastMinedBlock(gomock.Any()).Return().AnyTimes()
	// mock block
	mockBlock := model_mock.NewMockAbstractBlock(ctrl)
	// build
	mineConfig := MineConfig{
		GasFloor:         nil,
		GasCeil:          nil,
		CoinbaseAddress:  nil,
		BlockBuilder:     nil,
		BlockBroadcaster: mockBlockBroadcaster,
	}
	manager := newDefaultWorkManager(mineConfig)

	// test case
	situations := []struct{
		name string
		given func () map[common.Address]workerPerformance
		expectPerformance uint64
	}{
		{
			"empty performance map",
			func() map[common.Address]workerPerformance {
				wp := make(map[common.Address]workerPerformance)
				return wp
			},
			1,
		},
		{
				"already has performance map",
				func() map[common.Address]workerPerformance {
					wp := make(map[common.Address]workerPerformance)
					wp[address] = &defaultPerformance{blocksMined:10}
					return wp
				},
				11,
		},
	}
	// test
	for _, situation := range situations {
		manager.performance = situation.given()
		manager.submitBlock(address, mockBlock)
		assert.Equal(t, situation.expectPerformance, manager.getPerformance(address))
	}
}

func Test_defaultWorkManager_onNewBlock(t *testing.T) {
	// mock and init
	address := common.HexToAddress("123")
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// build
	mineConfig := MineConfig{
		GasFloor:         nil,
		GasCeil:          nil,
		CoinbaseAddress:  nil,
		BlockBuilder:     nil,
		BlockBroadcaster: nil,
	}
	manager := newDefaultWorkManager(mineConfig)
	// reward
	mReward := make(map[common.Address]*big.Int)
	mReward[address] = big.NewInt(10)
	manager.reward = mReward
	manager.totalReward = big.NewInt(10)
	// performance
	mPerformance := make(map[common.Address]workerPerformance)
	mPerformance[address] = &defaultPerformance{5}
	manager.performance = mPerformance

	// test case
	situations := []struct{
		name string
		given func () model.AbstractBlock
		expectReward *big.Int
	}{
		{
			"empty performance map",
			func() model.AbstractBlock {
				// mock block
				mockBlock := model_mock.NewMockAbstractBlock(ctrl)
				mockBlock.EXPECT().CoinBase().Return(big.NewInt(10)).AnyTimes()
				mockBlock.EXPECT().GetTransactionFees().Return(big.NewInt(1)).AnyTimes()
				return mockBlock
			},
			big.NewInt(21),
		},
	}
	// test
	for _, situation := range situations {
		block := situation.given()
		manager.onNewBlock(block)
		assert.Equal(t, situation.expectReward, manager.getReward(address))
	}
}