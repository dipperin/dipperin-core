package minemaster

import (
	"errors"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/mine/minemsg"
	"github.com/dipperin/dipperin-core/core/model"
	model_mock "github.com/dipperin/dipperin-core/tests/mock/model"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"sync/atomic"
	"testing"
)

func Test_newWorkDispatcher(t *testing.T) {
	mineConfig := MineConfig{
		GasFloor:         nil,
		GasCeil:          nil,
		CoinbaseAddress:  nil,
		BlockBuilder:     nil,
		BlockBroadcaster: nil,
	}
	getWorkersFunc := func() map[WorkerId]WorkerForMaster { return make(map[WorkerId]WorkerForMaster) }
	dispatcher := newWorkDispatcher(mineConfig, getWorkersFunc)
	assert.NotNil(t, dispatcher)
}

func Test_workDispatcher_curWorkBlock(t *testing.T) {
	// init
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mineConfig := MineConfig{
		GasFloor:         nil,
		GasCeil:          nil,
		CoinbaseAddress:  nil,
		BlockBuilder:     nil,
		BlockBroadcaster: nil,
	}
	getWorkersFunc := func() map[WorkerId]WorkerForMaster { return make(map[WorkerId]WorkerForMaster) }
	dispatcher := newWorkDispatcher(mineConfig, getWorkersFunc)

	// test case
	situations := []struct {
		name        string
		given       func() model.AbstractBlock
		expectIsNil bool
	}{
		{
			"nil case",
			func() model.AbstractBlock {
				return nil
			},
			true,
		},
		{
			"normal case which has block",
			func() model.AbstractBlock {
				absBlock := model_mock.NewMockAbstractBlock(ctrl)
				return absBlock
			},
			false,
		},
	}
	// test
	for _, situation := range situations {
		absBlock := situation.given()
		dispatcher.curBlock = absBlock
		// test result
		if situation.expectIsNil {
			assert.Nil(t, dispatcher.curWorkBlock())
		} else {
			assert.NotNil(t, dispatcher.curWorkBlock())
		}
	}
}

func Test_workDispatcher_dispatchNewWork(t *testing.T) {
	// init
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockBlockBuilder := NewMockBlockBuilder(ctrl)
	mockBlockBuilder.EXPECT().BuildWaitPackBlock(gomock.Any(), gomock.Any(), gomock.Any()).Return(func() model.AbstractBlock {
		absBlock := model_mock.NewMockAbstractBlock(ctrl)
		absBlock.EXPECT().Header().Return(&model.Header{})
		return absBlock
	}()).AnyTimes()
	mineConfig := MineConfig{
		GasFloor:         &atomic.Value{},
		GasCeil:          &atomic.Value{},
		CoinbaseAddress:  &atomic.Value{},
		BlockBuilder:     mockBlockBuilder,
		BlockBroadcaster: nil,
	}
	dispatcher := newWorkDispatcher(mineConfig, nil)

	// test case
	situations := []struct {
		name      string
		given     func() getWorkersFunc
		expectErr error
	}{
		{
			"empty workers map",
			func() getWorkersFunc {
				return func() map[WorkerId]WorkerForMaster {
					return make(map[WorkerId]WorkerForMaster)
				}
			},
			errors.New("no worker to dispatch work"),
		},
		{
			"normal case",
			func() getWorkersFunc {
				return func() map[WorkerId]WorkerForMaster {
					// mock
					mockWorkerForMaster := NewMockWorkerForMaster(ctrl)
					mockWorkerForMaster.EXPECT().SendNewWork(gomock.Any(), gomock.Any()).Return().AnyTimes()
					// build
					workerMap := make(map[WorkerId]WorkerForMaster)
					workerMap["test worker id"] = mockWorkerForMaster
					return workerMap
				}
			},
			nil,
		},
	}
	// test
	for _, situation := range situations {
		dispatcher.getWorkersFunc = situation.given()
		res := dispatcher.dispatchNewWork()
		// test result
		if situation.expectErr != nil {
			assert.Error(t, res)
		} else {
			assert.NoError(t, res)
		}
	}
}

func Test_workDispatcher_makeNewWorks(t *testing.T) {

}

func Test_workDispatcher_onNewBlock(t *testing.T) {
	// init
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// test case
	situations := []struct{
		name string
		given func() (*workDispatcher, int)
		expectCode int
		expectWorkLen int
	}{
		{
			"nil abstract of dispatcher",
			func() (*workDispatcher, int) {
				// mock
				mockBlockBuilder := NewMockBlockBuilder(ctrl)
				mockBlockBuilder.EXPECT().BuildWaitPackBlock(gomock.Any(),gomock.Any(),gomock.Any()).Return(func() model.AbstractBlock {
					return nil
				}()).AnyTimes()
				// build dispatcher
				mineConfig := MineConfig{
					GasFloor:         &atomic.Value{},
					GasCeil:          &atomic.Value{},
					CoinbaseAddress:  &atomic.Value{},
					BlockBuilder:     mockBlockBuilder,
					BlockBroadcaster: nil,
				}
				return newWorkDispatcher(mineConfig, nil), 5

			},
			0,
			0,
		},
		{
			"normal build",
			func() (*workDispatcher, int) {
				// mock
				mockBlockBuilder := NewMockBlockBuilder(ctrl)
				mockBlockBuilder.EXPECT().BuildWaitPackBlock(gomock.Any(),gomock.Any(),gomock.Any()).Return(func() model.AbstractBlock {
					// mock block
					mockBlock := model_mock.NewMockAbstractBlock(ctrl)
					mockBlock.EXPECT().Header().Return(func() *model.Header {
						nonce := common.EncodeNonce(137)
						return &model.Header{Nonce: nonce}
					}())
					// return block
					return mockBlock
				}()).AnyTimes()
				// build dispatcher
				mineConfig := MineConfig{
					GasFloor:         &atomic.Value{},
					GasCeil:          &atomic.Value{},
					CoinbaseAddress:  &atomic.Value{},
					BlockBuilder:     mockBlockBuilder,
					BlockBroadcaster: nil,
				}
				dispatcher := newWorkDispatcher(mineConfig, nil)
				// return dispatcher
				return dispatcher, 3
			},
			minemsg.NewDefaultWorkMsg,
			3,
		},
	}
	// test
	for _, situation := range situations {
		dispatcher, workerLen := situation.given()
		resCode, resWork := dispatcher.makeNewWorks(workerLen)
		// result
		assert.Equal(t, situation.expectCode, resCode)
		assert.Equal(t, situation.expectWorkLen, len(resWork))
	}
}
