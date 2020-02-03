package components

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/golang/mock/gomock"
	"github.com/dipperin/dipperin-core/common"
)

func TestNewFetcher(t *testing.T) {
	rsp := NewFetcher(nil)
	//c := &CsBftFetcher{
	//	fc:             nil,
	//	requests:       make(map[uint64]*FetchBlockReqMsg),
	//	fetchReqQueue:  make(chan *FetchBlockReqMsg, 1),
	//	fetchRespChan:  make(chan *FetchBlockRespMsg, 1),
	//	isFetchingChan: make(chan *IsFetchingMsg),
	//	rmReqChan:      make(chan uint64),
	//}
	//c.BaseService = *util.NewBaseService(log.DLogger, "cs_bft_fetcher", c)
	//assert.IsEqual(rsp, c)
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
	NewFetcher(nil).FetchBlock(ww, common.HexToHash(hashTmp))
}
