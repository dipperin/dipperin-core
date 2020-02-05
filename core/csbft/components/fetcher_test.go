package components

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/golang/mock/gomock"
	"github.com/dipperin/dipperin-core/common"
)

func TestNewFetcher(t *testing.T) {
	rsp := NewFetcher(nil)
	assert.NotEmpty(t, rsp)
}

func TestCsBftFetcher_FetchBlock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	var hashTmp = `0xd50866a60b4f7e494400e0563efb987dc800d1a72af5cc1ae9ee68760bb18889`
	ww := [22]byte{123}
	fetcher := NewMockFetcher(ctrl)
	gomock.InOrder(
		fetcher.EXPECT().FetchBlock(ww, common.HexToHash(hashTmp)).Return(nil).AnyTimes(),
	)
	//NewFetcher(nil).FetchBlock(ww, common.HexToHash(hashTmp))
}

func TestCsBftFetcher_FetchBlockResp(t *testing.T) {
	rsp := NewFetcher(nil)
	rsp.FetchBlockResp(nil)
}

func TestCsBftFetcher_rmReq(t *testing.T)  {
	rsp := NewFetcher(nil)
	rsp.rmReq(1)
}

func TestCsBftFetcher_OnStart(t *testing.T) {
	rsp := NewFetcher(nil)
	assert.NoError(t, rsp.OnStart())
}

func TestCsBftFetcher_OnStop(t *testing.T) {
	rsp := NewFetcher(nil)
	rsp.OnStop()
}

func TestCsBftFetcher_OnReset(t *testing.T) {
	rsp := NewFetcher(nil)
	assert.NoError(t, rsp.OnReset())
}

func TestCsBftFetcher_IsFetching(t *testing.T) {
	rsp := NewFetcher(nil)
	var hashTmp = `0xd50866a60b4f7e4123400e0563efb987dc800d1a72af5cc1ae9ee68760bb18889`
	assert.Equal(t, false, rsp.isFetching(common.HexToHash(hashTmp)))
}
