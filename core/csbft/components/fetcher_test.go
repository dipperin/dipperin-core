package components

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/model"
	"time"
	csModel "github.com/dipperin/dipperin-core/core/csbft/model"
	"github.com/golang/mock/gomock"
)

type myFetcherConn struct {
}

func (myFetcherConn) SendFetchBlockMsg(msgCode uint64, from common.Address, msg *csModel.FetchBlockReqDecodeMsg) error {
	return nil
}

func TestNewFetcher(t *testing.T) {
	rsp := NewFetcher(myFetcherConn{})
	assert.NotEmpty(t, rsp)
}

func TestCsBftFetcher_FetchBlock(t *testing.T) {
	//ctrl := gomock.NewController(t)
	//defer ctrl.Finish()
	//var hashTmp = `0xd50866a60b4f7e494400e0563efb987dc800d1a72af5cc1ae9ee68760bb18889`
	//ww := common.Address{1, 2, 3}
	//fetcher := NewMockFetcher(ctrl)
	//var w model.AbstractBlock
	//gomock.InOrder(
	//	fetcher.EXPECT().FetchBlock(ww, common.HexToHash(hashTmp)).Return(w),
	//)
	//if f := fetcher.FetchBlock(ww, common.HexToHash(hashTmp)); f != w {
	//	t.Errorf("FetchBlock: got %v, want %v", f, w)
	//}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	testCases := []struct {
		name   string
		given  func() bool
		expect bool
	}{
		{
			name: "FetchBlock fail to start",
			given: func() bool {
				rsp := NewFetcher(myFetcherConn{})
				assert.NotEmpty(t, rsp)
				//rsp.BaseService.Start()
				//rsp.FetchBlock(common.Address{0x01},common.Hash{})
				return rsp.FetchBlock(common.Address{0x01}, common.Hash{}) == nil
			},
			expect: true,
		},
		{
			name: "FetchBlock true",
			given: func() bool {
				rsp := NewFetcher(myFetcherConn{})
				assert.NotEmpty(t, rsp)
				rsp.BaseService.Start()
				//rsp.FetchBlock(common.Address{0x01},common.Hash{})
				return rsp.FetchBlock(common.Address{0x01}, common.Hash{}) == nil
			},
			expect: true,
		},
	}

	for i, tc := range testCases {
		sign := tc.given()
		if testCases[i].expect == sign {
			t.Log("success")
		} else {
			t.Log("failure")
		}
	}
}

func TestCsBftFetcher_FetchBlockResp(t *testing.T) {
	rsp := NewFetcher(nil)
	assert.NotEmpty(t, rsp)
	assert.NotPanics(t, func() {
		rsp.FetchBlockResp(nil)
	})
}

func TestCsBftFetcher_rmReq(t *testing.T) {
	rsp := NewFetcher(nil)
	assert.NotEmpty(t, rsp)
	assert.NotPanics(t, func() {
		rsp.rmReq(1)
	})
}

func TestCsBftFetcher_OnStart(t *testing.T) {
	rsp := NewFetcher(nil)
	assert.NotEmpty(t, rsp)
	assert.NoError(t, rsp.OnStart())
}

func TestCsBftFetcher_OnStop(t *testing.T) {
	rsp := NewFetcher(nil)
	assert.NotEmpty(t, rsp)
	if assert.NoError(t, rsp.OnStart()) {
		time.Sleep(500 * time.Millisecond)
		assert.NotPanics(t, rsp.OnStop)
	}
}

func TestCsBftFetcher_OnReset(t *testing.T) {
	rsp := NewFetcher(nil)
	assert.NotEmpty(t, rsp)
	assert.NoError(t, rsp.OnReset())
}

func TestCsBftFetcher_loop(t *testing.T) {
	rsp := NewFetcher(nil)
	assert.NotEmpty(t, rsp)
	assert.NotPanics(t, func() {
		go func(rsp *CsBftFetcher) {
			rsp.loop()
		}(rsp)
	})
}

func TestCsBftFetcher_IsFetching(t *testing.T) {
	rsp := NewFetcher(nil)
	assert.NotEmpty(t, rsp)
	var hashTmp = `0xd50866a60b4f7e4123400e0563efb987dc800d1a72af5cc1ae9ee68760bb18889`
	assert.Equal(t, false, rsp.IsFetching(common.HexToHash(hashTmp)))
}

func TestCsBftFetcher_isFetching(t *testing.T) {
	rsp := NewFetcher(nil)
	assert.NotEmpty(t, rsp)
	var hashTmp = `0xd50866a60b4f7e4123400e0563efb987dc800d1a72af5cc1ae9ee68760bb18889`
	assert.Equal(t, false, rsp.isFetching(common.HexToHash(hashTmp)))
}

func TestCsBftFetcher_onFetchBlock(t *testing.T) {
	rsp := NewFetcher(nil)
	assert.NotEmpty(t, rsp)
	var hashTmp = `0xd50866a60b4f7e4123400e0563efb987dc800d1a72af5cc1ae9ee68760bb18889`
	var req = FetchBlockReqMsg{
		MsgId:      0,
		From:       common.HexToAddress(hashTmp),
		BlockHash:  common.HexToHash(hashTmp),
		ResultChan: make(chan model.AbstractBlock),
	}
	assert.Panics(t, func() {
		rsp.onFetchBlock(&req)
	})
}

func TestCsBftFetcher_onFetchResp(t *testing.T) {
	rsp := NewFetcher(nil)
	assert.NotEmpty(t, rsp)
	var req = FetchBlockRespMsg{
		MsgId: 1,
		Block: nil,
	}
	var hashTmp = `0xd50866a60b4f7e4123400e0563efb987dc800d1a72af5cc1ae9ee68760bb18889`
	r := FetchBlockReqMsg{
		MsgId:      0,
		From:       common.HexToAddress(hashTmp),
		BlockHash:  common.HexToHash(hashTmp),
		ResultChan: make(chan model.AbstractBlock),
	}
	rsp.requests[1] = &r
	assert.Panics(t, func() {
		rsp.onFetchResp(&req)
	})

	rsp1 := NewFetcher(nil)
	assert.NotEmpty(t, rsp1)
	assert.NotPanics(t, func() {
		rsp1.onFetchResp(&req)
	})
}

func TestCsBftFetcher_onResult(t *testing.T) {
	var hashTmp = `0xd50866a60b4f7e4123400e0563efb987dc800d1a72af5cc1ae9ee68760bb18889`
	r := FetchBlockReqMsg{
		MsgId:      0,
		From:       common.HexToAddress(hashTmp),
		BlockHash:  common.HexToHash(hashTmp),
		ResultChan: make(chan model.AbstractBlock),
	}
	assert.NotPanics(t, func() {
		r.onResult(nil)
	})
}
