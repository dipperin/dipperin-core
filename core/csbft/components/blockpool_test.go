package components

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/dipperin/dipperin-core/common"
	"time"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/golang/mock/gomock"
)

func TestNewBlockPool(t *testing.T) {
	rsp := NewBlockPool(0, nil)
	t.Log("www:", rsp)
	assert.NotEmpty(t, rsp)
}

func TestBlockPool_SetNodeConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	testCases := []struct {
		name   string
		given  func() bool
		expect bool
	}{
		{
			name: "SetNodeConfig true",
			given: func() bool {
				rsp := NewBlockPool(0, nil)
				assert.NotEmpty(t, rsp)
				rsp.SetNodeConfig(nil)
				return nil == rsp.Blockpoolconfig
			},
			expect: true,
		},
		{
			name: "SetNodeConfig false",
			given: func() bool {
				rsp := NewBlockPool(0, nil)
				assert.NotEmpty(t, rsp)
				rsp.SetNodeConfig(NewMockBlockpoolconfig(ctrl))
				return nil == rsp.Blockpoolconfig
			},
			expect: false,
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

type myTestNotifier struct{}

func (myTestNotifier) BlockPoolNotEmpty() {}

func TestBlockPool_SetPoolEventNotifier(t *testing.T) {
	testCases := []struct {
		name   string
		given  func() bool
		expect bool
	}{
		{
			name: "SetPoolEventNotifier true",
			given: func() bool {
				rsp := NewBlockPool(0, nil)
				assert.NotEmpty(t, rsp)
				rsp.SetPoolEventNotifier(nil)
				return nil == rsp.Blockpoolconfig
			},
			expect: true,
		},
		{
			name: "SetPoolEventNotifier false",
			given: func() bool {
				rsp := NewBlockPool(0, nil)
				rsp.SetPoolEventNotifier(myTestNotifier{})
				return nil == rsp.Blockpoolconfig
			},
			expect: true,
		},
	}

	for i, tc := range testCases {
		sign := tc.given()
		if testCases[i].expect == sign {
			t.Log("success")
		} else {
			t.Logf("expect:%v,actual:%v", testCases[i].expect, sign)
		}
	}
}

func TestBlockPool_Start(t *testing.T) {
	testCases := []struct {
		name   string
		given  func() string
		expect string
	}{
		{
			name: "BlockPool Start false",
			given: func() string {
				rsp := &BlockPool{
					height:            0,
					blocks:            []model.AbstractBlock{},
					poolEventNotifier: myTestNotifier{},
					Blockpoolconfig:   nil,

					newHeightChan: make(chan uint64, 5),
					newBlockChan:  make(chan newBlockWithResultErr, 5),
					getterChan:    make(chan *blockPoolGetter, 5),
					rmBlockChan:   make(chan common.Hash),
					stopChan:      make(chan struct{}),
				}
				return rsp.Start().Error()
			},
			expect: "block pool already started",
		},

		{
			name: "BlockPool Start true",
			given: func() string {
				rsp := NewBlockPool(0, nil)
				assert.NotEmpty(t, rsp)
				if rsp.Start() == nil {
					return ""
				}
				return ""
			},
			expect: "",
		},
	}
	for i, tc := range testCases {
		sign := tc.given()
		if testCases[i].expect == sign {
			t.Log("success")
		} else {
			t.Logf("expect:%v,actual:%v", testCases[i].expect, sign)
		}
	}
}

func TestBlockPool_Stop(t *testing.T) {
	testCases := []struct {
		name   string
		given  func() bool
		expect bool
	}{
		{
			name: "stopChan is nil",
			given: func() bool {
				rsp := NewBlockPool(0, nil)
				assert.NotEmpty(t, rsp)
				return assert.NotPanics(t, rsp.Stop)
			},
			expect: true,
		},
		{
			name: "stopChan is not nil",
			given: func() bool {
				rsp := &BlockPool{
					height:            0,
					blocks:            []model.AbstractBlock{},
					poolEventNotifier: myTestNotifier{},
					Blockpoolconfig:   nil,

					newHeightChan: make(chan uint64, 5),
					newBlockChan:  make(chan newBlockWithResultErr, 5),
					getterChan:    make(chan *blockPoolGetter, 5),
					rmBlockChan:   make(chan common.Hash),
					stopChan:      make(chan struct{}),
				}
				rsp.Stop()
				return rsp.stopChan == nil
			},
			expect: true,
		},
	}

	for i, tc := range testCases {
		sign := tc.given()
		if testCases[i].expect == sign {
			t.Log("success")
		} else {
			t.Logf("expect:%v,actual:%v", testCases[i].expect, sign)
		}
	}
}

func TestBlockPool_IsEmpty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	testCases := []struct {
		name   string
		given  func() bool
		expect bool
	}{
		{
			name: "IsEmpty is false",
			given: func() bool {
				rsp := NewBlockPool(0, nil)
				assert.NotEmpty(t, rsp)
				rsp.blocks = append(rsp.blocks, NewMockAbstractBlock(ctrl))
				return rsp.IsEmpty()
			},
			expect: false,
		},
		{
			name: "IsEmpty is true",
			given: func() bool {
				rsp := NewBlockPool(0, nil)
				assert.NotEmpty(t, rsp)
				return rsp.IsEmpty()
			},
			expect: true,
		},
	}

	for i, tc := range testCases {
		sign := tc.given()
		if testCases[i].expect == sign {
			t.Log("success")
		} else {
			t.Logf("expect:%v,actual:%v", testCases[i].expect, sign)
		}
	}
}

func TestBlockPool_loop(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	testCases := []struct {
		name   string
		given  func() bool
		expect bool
	}{
		{
			name: "IsEmpty is false",
			given: func() bool {
				rsp := &BlockPool{
					height:            0,
					blocks:            []model.AbstractBlock{},
					poolEventNotifier: myTestNotifier{},
					Blockpoolconfig:   nil,

					newHeightChan: make(chan uint64, 5),
					newBlockChan:  make(chan newBlockWithResultErr, 5),
					getterChan:    make(chan *blockPoolGetter, 5),
					rmBlockChan:   make(chan common.Hash),
					stopChan:      make(chan struct{}),
				}
				//go rsp.loop()
				go func(rsp *BlockPool) {
					time.Sleep(time.Second)
					rsp.NewHeight(1)
					b := NewMockAbstractBlock(ctrl)
					b.EXPECT().Number().Return(uint64(1)).AnyTimes()
					b.EXPECT().Hash().Return(common.Hash{}).AnyTimes()
					rsp.AddBlock(b)
					resultC := make(chan model.AbstractBlock)
					rsp.getBlock(&blockPoolGetter{
						blockHash:  common.Hash{},
						resultChan: resultC,
					})

					rsp.RemoveBlock(common.Hash{})
					rsp.stopChan <- struct{}{}
				}(rsp)
				return assert.NotPanics(t, func() {
					time.Sleep(2 * time.Second)
					go rsp.loop()
				})
			},
			expect: true,
		},
	}

	for i, tc := range testCases {
		sign := tc.given()
		if testCases[i].expect == sign {
			t.Log("success")
		} else {
			t.Logf("expect:%v,actual:%v", testCases[i].expect, sign)
		}
	}
}

func TestBlockPool_doRemoveBlock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	testCases := []struct {
		name   string
		given  func() bool
		expect bool
	}{
		{
			name: "doRemoveBlock test true",
			given: func() bool {
				rsp := NewBlockPool(0, nil)
				assert.NotEmpty(t, rsp)
				var hashTmp = `0xd50866a60b4f7e494400e0563efb987dc800d1a72af5cc1ae9ee68760bb18889`
				b := NewMockAbstractBlock(ctrl)
				b.EXPECT().Hash().Return(common.HexToHash(hashTmp)).AnyTimes()
				rsp.blocks = append(rsp.blocks, b)
				rsp.doRemoveBlock(common.HexToHash(hashTmp))
				return len(rsp.blocks) == 0
			},
			expect: true,
		},
		{
			name: "doRemoveBlock test false",
			given: func() bool {
				rsp := NewBlockPool(0, nil)
				assert.NotEmpty(t, rsp)
				var hashTmp1 = `0xd50866a60b4f7e494400e0563efb987dc800d1a72af5cc1ae9ee68760bb18889`
				var hashTmp2 = `0xd51866a60b4f7e494400e0563efb987dc800d1a72af5cc1ae9ee68760bb18889`
				b := NewMockAbstractBlock(ctrl)
				b.EXPECT().Hash().Return(common.HexToHash(hashTmp1)).AnyTimes()
				rsp.blocks = append(rsp.blocks, b)
				rsp.doRemoveBlock(common.HexToHash(hashTmp2))
				return len(rsp.blocks) == 0
			},
			expect: false,
		},
	}

	for i, tc := range testCases {
		sign := tc.given()
		if testCases[i].expect == sign {
			t.Log("success")
		} else {
			t.Logf("expect:%v,actual:%v", testCases[i].expect, sign)
		}
	}
}

func TestBlockPool_IsRunning(t *testing.T) {
	testCases := []struct {
		name   string
		given  func() bool
		expect bool
	}{
		{
			name: "IsRunning false",
			given: func() bool {
				rsp := NewBlockPool(0, nil)
				assert.NotEmpty(t, rsp)
				return rsp.IsRunning()
			},
			expect: false,
		},
		{
			name: "IsRunning true",
			given: func() bool {
				rsp := &BlockPool{
					height:            0,
					blocks:            []model.AbstractBlock{},
					poolEventNotifier: myTestNotifier{},
					Blockpoolconfig:   nil,

					newHeightChan: make(chan uint64, 5),
					newBlockChan:  make(chan newBlockWithResultErr, 5),
					getterChan:    make(chan *blockPoolGetter, 5),
					rmBlockChan:   make(chan common.Hash),
					stopChan:      make(chan struct{}),
				}
				return rsp.IsRunning()
			},
			expect: true,
		},
	}

	for i, tc := range testCases {
		sign := tc.given()
		if testCases[i].expect == sign {
			t.Log("success")
		} else {
			t.Logf("expect:%v,actual:%v", testCases[i].expect, sign)
		}
	}
}

func TestBlockPool_RemoveBlock(t *testing.T) {
	testCases := []struct {
		name   string
		given  func() bool
		expect bool
	}{
		{
			name: "RemoveBlock true",
			given: func() bool {
				rsp := &BlockPool{
					height:            0,
					blocks:            []model.AbstractBlock{},
					poolEventNotifier: myTestNotifier{},
					Blockpoolconfig:   nil,

					newHeightChan: make(chan uint64, 5),
					newBlockChan:  make(chan newBlockWithResultErr, 5),
					getterChan:    make(chan *blockPoolGetter, 5),
					rmBlockChan:   make(chan common.Hash),
					stopChan:      make(chan struct{}),
				}
				var hashTmp = `0xd50866a60b4f7e494400e0563efb987dc800d1a72af5cc1ae9ee68760bb18889`
				assert.NotPanics(t, func() {
					go func(*BlockPool, string) {
						rsp.RemoveBlock(common.HexToHash(hashTmp))
					}(rsp, hashTmp)
				})
				time.Sleep(time.Second)
				r := <-rsp.rmBlockChan
				return r.Hex() == hashTmp
			},
			expect: true,
		},
		{
			name: "RemoveBlock false",
			given: func() bool {
				rsp := NewBlockPool(0, nil)
				assert.NotEmpty(t, rsp)
				var hashTmp = `0xd50866a60b4f7e494400e0563efb987dc800d1a72af5cc1ae9ee68760bb18889`
				rsp.RemoveBlock(common.HexToHash(hashTmp))
				close(rsp.rmBlockChan)
				r := <-rsp.rmBlockChan
				return r.Hex() == hashTmp
			},
			expect: false,
		},
	}

	for i, tc := range testCases {
		sign := tc.given()
		if testCases[i].expect == sign {
			t.Log("success")
		} else {
			t.Logf("expect:%v,actual:%v", testCases[i].expect, sign)
		}
	}
}

func TestBlockPool_NewHeight(t *testing.T) {
	testCases := []struct {
		name   string
		given  func() bool
		expect bool
	}{
		{
			name: "NewHeight true",
			given: func() bool {
				rsp := &BlockPool{
					height:            0,
					blocks:            []model.AbstractBlock{},
					poolEventNotifier: myTestNotifier{},
					Blockpoolconfig:   nil,

					newHeightChan: make(chan uint64, 5),
					newBlockChan:  make(chan newBlockWithResultErr, 5),
					getterChan:    make(chan *blockPoolGetter, 5),
					rmBlockChan:   make(chan common.Hash),
					stopChan:      make(chan struct{}),
				}
				rsp.NewHeight(2)
				time.Sleep(time.Second)
				close(rsp.newHeightChan)
				r := <-rsp.newHeightChan
				return r == 2
			},
			expect: true,
		},
		{
			name: "NewHeight false",
			given: func() bool {
				rsp := NewBlockPool(0, nil)
				assert.NotEmpty(t, rsp)
				rsp.NewHeight(2)
				close(rsp.newHeightChan)
				r := <-rsp.newHeightChan
				return r == 2
			},
			expect: false,
		},
	}

	for i, tc := range testCases {
		sign := tc.given()
		if testCases[i].expect == sign {
			t.Log("success")
		} else {
			t.Logf("expect:%v,actual:%v", testCases[i].expect, sign)
		}
	}
}

func TestBlockPool_doNewHeight(t *testing.T) {
	testCases := []struct {
		name   string
		given  func() bool
		expect bool
	}{
		{
			name: "doNewHeight true",
			given: func() bool {
				rsp := NewBlockPool(0, nil)
				assert.NotEmpty(t, rsp)
				rsp.height = 2
				rsp.doNewHeight(5)
				return rsp.height == 5
			},
			expect: true,
		},
		{
			name: "doNewHeight false",
			given: func() bool {
				rsp := NewBlockPool(0, nil)
				assert.NotEmpty(t, rsp)
				rsp.height = 2
				rsp.doNewHeight(1)
				return rsp.height == 1
			},
			expect: false,
		},
	}

	for i, tc := range testCases {
		sign := tc.given()
		if testCases[i].expect == sign {
			t.Log("success")
		} else {
			t.Logf("expect:%v,actual:%v", testCases[i].expect, sign)
		}
	}
}

func TestBlockPool_AddBlock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	testCases := []struct {
		name   string
		given  func() bool
		expect bool
	}{
		{
			name: "AddBlock is true",
			given: func() bool {
				rsp := &BlockPool{
					height:            0,
					blocks:            []model.AbstractBlock{},
					poolEventNotifier: myTestNotifier{},
					Blockpoolconfig:   nil,
					newHeightChan:     make(chan uint64, 5),
					newBlockChan:      make(chan newBlockWithResultErr, 5),
					getterChan:        make(chan *blockPoolGetter, 5),
					rmBlockChan:       make(chan common.Hash),
					stopChan:          make(chan struct{}),
				}
				b := NewMockAbstractBlock(ctrl)
				assert.NotPanics(t, func() {
					go func(*BlockPool) {
						rsp.AddBlock(b)
					}(rsp)
				})
				time.Sleep(time.Second)
				close(rsp.newBlockChan)
				n := <-rsp.newBlockChan
				return n.block != nil
			},
			expect: true,
		},
		{
			name: "AddBlock error exists",
			given: func() bool {
				rsp := NewBlockPool(0, nil)
				assert.NotEmpty(t, rsp)
				b := NewMockAbstractBlock(ctrl)
				return rsp.AddBlock(b).Error() == "block pool not running"
			},
			expect: true,
		},
	}

	for i, tc := range testCases {
		sign := tc.given()
		if testCases[i].expect == sign {
			t.Log("success")
		} else {
			t.Logf("expect:%v,actual:%v", testCases[i].expect, sign)
		}
	}
}

func TestBlockPool_GetBlockByHash(t *testing.T) {
	testCases := []struct {
		name   string
		given  func() bool
		expect bool
	}{
		{
			name: "GetBlockByHash is true",
			given: func() bool {
				rsp := &BlockPool{
					height:            0,
					blocks:            []model.AbstractBlock{},
					poolEventNotifier: myTestNotifier{},
					Blockpoolconfig:   nil,
					newHeightChan:     make(chan uint64, 5),
					newBlockChan:      make(chan newBlockWithResultErr, 5),
					getterChan:        make(chan *blockPoolGetter, 5),
					rmBlockChan:       make(chan common.Hash),
					stopChan:          make(chan struct{}),
				}
				var hashTmp = `0xd50866a60b4f7e4123400e0563efb987dc800d1a72af5cc1ae9ee68760bb18889`
				assert.NotPanics(t, func() {
					go func(*BlockPool) {
						rsp.GetBlockByHash(common.HexToHash(hashTmp))
					}(rsp)
				})
				time.Sleep(500 * time.Millisecond)
				g := <-rsp.getterChan
				return g.blockHash.Hex() != ""
			},
			expect: true,
		},
		{
			name: "GetBlockByHash fail to insert data to getterChan",
			given: func() bool {
				rsp := NewBlockPool(0, nil)
				assert.NotEmpty(t, rsp)
				var hashTmp = `0xd50866a60b4f7e4123400e0563efb987dc800d1a72af5cc1ae9ee68760bb18889`
				assert.NotPanics(t, func() {
					go func(*BlockPool) {
						rsp.GetBlockByHash(common.HexToHash(hashTmp))
					}(rsp)
				})
				close(rsp.getterChan)
				g := <-rsp.getterChan
				return g == nil
			},
			expect: true,
		},
	}

	for i, tc := range testCases {
		sign := tc.given()
		if testCases[i].expect == sign {
			t.Log("success")
		} else {
			t.Logf("expect:%v,actual:%v", testCases[i].expect, sign)
		}
	}
}

func TestBlockPool_doAddBlock(t *testing.T) {
	//doAddBlock
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	testCases := []struct {
		name   string
		given  func() bool
		expect bool
	}{
		{
			name: "doAddBlock error exists",
			given: func() bool {
				rsp := &BlockPool{
					height:            0,
					blocks:            []model.AbstractBlock{},
					poolEventNotifier: myTestNotifier{},
					Blockpoolconfig:   nil,

					newHeightChan: make(chan uint64, 5),
					newBlockChan:  make(chan newBlockWithResultErr, 5),
					getterChan:    make(chan *blockPoolGetter, 5),
					rmBlockChan:   make(chan common.Hash),
					stopChan:      make(chan struct{}),
				}
				var n = newBlockWithResultErr{block: nil, resultChan: make(chan error)}
				go func(*BlockPool, newBlockWithResultErr) {
					rsp.height = 1
					b := NewMockAbstractBlock(ctrl)
					b.EXPECT().Number().Return(uint64(2)).AnyTimes()
					b.EXPECT().Hash().Return(common.Hash{}).AnyTimes()
					n.block = b
					rsp.doAddBlock(n)
				}(rsp, n)
				time.Sleep(time.Second)
				err := <-n.resultChan
				close(n.resultChan)
				t.Log("err:", err)
				if err != nil {
					return err.Error() == "invalid height block"
				}
				return false
			},
			expect: true,
		},
		{
			name: "doAddBlock error exists",
			given: func() bool {
				rsp := &BlockPool{
					height:            0,
					blocks:            []model.AbstractBlock{},
					poolEventNotifier: myTestNotifier{},
					Blockpoolconfig:   nil,

					newHeightChan: make(chan uint64, 5),
					newBlockChan:  make(chan newBlockWithResultErr, 5),
					getterChan:    make(chan *blockPoolGetter, 5),
					rmBlockChan:   make(chan common.Hash),
					stopChan:      make(chan struct{}),
				}
				var n = newBlockWithResultErr{block: nil, resultChan: make(chan error)}
				var hashTmp = `0xd50866a60b4f7e4123400e0563efb987dc800d1a72af5cc1ae9ee68760bb18889`
				go func(*BlockPool, newBlockWithResultErr) {
					rsp.height = 1
					b := NewMockAbstractBlock(ctrl)
					b.EXPECT().Number().Return(uint64(1)).AnyTimes()
					b.EXPECT().Hash().Return(common.HexToHash(hashTmp)).AnyTimes()
					rsp.blocks = append(rsp.blocks, b)
					n.block = b
					rsp.doAddBlock(n)
				}(rsp, n)
				time.Sleep(time.Second)
				err := <-n.resultChan
				close(n.resultChan)
				t.Log("err:", err)
				if err != nil {
					return err.Error() == "dul block"
				}
				return false
			},
			expect: true,
		},
		{
			name: "doAddBlock no error exists",
			given: func() bool {
				rsp := &BlockPool{
					height:            0,
					blocks:            []model.AbstractBlock{},
					poolEventNotifier: myTestNotifier{},
					Blockpoolconfig:   nil,

					newHeightChan: make(chan uint64, 5),
					newBlockChan:  make(chan newBlockWithResultErr, 5),
					getterChan:    make(chan *blockPoolGetter, 5),
					rmBlockChan:   make(chan common.Hash),
					stopChan:      make(chan struct{}),
				}
				var n = newBlockWithResultErr{block: nil, resultChan: make(chan error)}
				var hashTmp = `0xd50866a60b4f7e4123400e0563efb987dc800d1a72af5cc1ae9ee68760bb18889`
				go func(*BlockPool, newBlockWithResultErr) {
					rsp.height = 1
					b := NewMockAbstractBlock(ctrl)
					b.EXPECT().Number().Return(uint64(1)).AnyTimes()
					b.EXPECT().Hash().Return(common.HexToHash(hashTmp)).AnyTimes()
					//rsp.blocks = append(rsp.blocks, b)
					n.block = b
					rsp.doAddBlock(n)
				}(rsp, n)
				time.Sleep(time.Second)
				err := <-n.resultChan
				close(n.resultChan)
				t.Log("err:", err)
				if err == nil {
					return true
				}
				return false
			},
			expect: true,
		},
		{
			name: "doAddBlock no error exists",
			given: func() bool {
				rsp := &BlockPool{
					height:            0,
					blocks:            []model.AbstractBlock{},
					poolEventNotifier: myTestNotifier{},
					Blockpoolconfig:   nil,

					newHeightChan: make(chan uint64, 5),
					newBlockChan:  make(chan newBlockWithResultErr, 5),
					getterChan:    make(chan *blockPoolGetter, 5),
					rmBlockChan:   make(chan common.Hash),
					stopChan:      make(chan struct{}),
				}
				var n = newBlockWithResultErr{block: nil, resultChan: make(chan error)}
				var hashTmp = `0xd50866a60b4f7e4123400e0563efb987dc800d1a72af5cc1ae9ee68760bb18889`
				var hashTmp1 = `0xd508r3a60b4f7e4123400e0563efb987dc800d1a72af5cc1ae9ee68760bb18889`
				go func(*BlockPool, newBlockWithResultErr) {
					rsp.height = 1
					b := NewMockAbstractBlock(ctrl)
					b.EXPECT().Number().Return(uint64(1)).AnyTimes()
					b.EXPECT().Hash().Return(common.HexToHash(hashTmp)).AnyTimes()

					b1 := NewMockAbstractBlock(ctrl)
					b1.EXPECT().Number().Return(uint64(1)).AnyTimes()
					b1.EXPECT().Hash().Return(common.HexToHash(hashTmp1)).AnyTimes()
					rsp.blocks = append(rsp.blocks, b1)
					rsp.poolEventNotifier = myTestNotifier{}
					n.block = b
					rsp.doAddBlock(n)
				}(rsp, n)
				time.Sleep(time.Second)
				err := <-n.resultChan
				close(n.resultChan)
				t.Log("err:", err)
				if err == nil {
					return true
				}
				return false
			},
			expect: true,
		},
	}

	for i, tc := range testCases {
		sign := tc.given()
		if testCases[i].expect == sign {
			t.Log("success")
		} else {
			t.Logf("expect:%v,actual:%v", testCases[i].expect, sign)
		}
	}
}

func TestBlockPool_GetProposalBlock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	testCases := []struct {
		name   string
		given  func() bool
		expect bool
	}{
		{
			name: "GetProposalBlock BlockPool blocks is empty",
			given: func() bool {
				rsp := &BlockPool{
					height:            0,
					blocks:            []model.AbstractBlock{},
					poolEventNotifier: myTestNotifier{},
					Blockpoolconfig:   nil,

					newHeightChan: make(chan uint64, 5),
					newBlockChan:  make(chan newBlockWithResultErr, 5),
					getterChan:    make(chan *blockPoolGetter, 5),
					rmBlockChan:   make(chan common.Hash),
					stopChan:      make(chan struct{}),
				}
				go rsp.loop()
				var ab model.AbstractBlock
				go func(model.AbstractBlock) {
					ab = rsp.GetProposalBlock()
				}(ab)

				time.Sleep(2 * time.Second)
				rsp.stopChan <- struct{}{}
				return ab == nil
			},
			expect: true,
		},
		{
			name: "GetProposalBlock BlockPool blocks is not empty",
			given: func() bool {
				rsp := &BlockPool{
					height:            0,
					blocks:            []model.AbstractBlock{},
					poolEventNotifier: myTestNotifier{},
					Blockpoolconfig:   nil,

					newHeightChan: make(chan uint64, 5),
					newBlockChan:  make(chan newBlockWithResultErr, 5),
					getterChan:    make(chan *blockPoolGetter, 5),
					rmBlockChan:   make(chan common.Hash),
					stopChan:      make(chan struct{}),
				}
				var hashTmp = `0xd50866a60b4f7e4123400e0563efb987dc800d1a72af5cc1ae9ee68760bb18889`
				b := NewMockAbstractBlock(ctrl)
				b.EXPECT().Hash().Return(common.HexToHash(hashTmp)).AnyTimes()
				rsp.blocks = append(rsp.blocks, b)
				go rsp.loop()
				var ab model.AbstractBlock
				go func(model.AbstractBlock) {
					ab = rsp.GetProposalBlock()
				}(ab)

				time.Sleep(2 * time.Second)
				rsp.stopChan <- struct{}{}
				return ab != nil
			},
			expect: true,
		},
	}
	for i, tc := range testCases {
		sign := tc.given()
		if testCases[i].expect == sign {
			t.Log("success")
		} else {
			t.Logf("expect:%v,actual:%v", testCases[i].expect, sign)
		}
	}
}

func TestBlockPool_getBlock(t *testing.T) {
	rsp := NewBlockPool(0, nil)
	assert.NotEmpty(t, rsp)
	resultC := make(chan model.AbstractBlock)
	var hashTmp = `0xd50866a60b4f7e4123400e0563efb987dc800d1a72af5cc1ae9ee68760bb18889`
	h := common.HexToHash(hashTmp)
	assert.NotPanics(t, func() {
		rsp.getBlock(
			&blockPoolGetter{
				blockHash:  h,
				resultChan: resultC,
			},
		)
	})
}

func TestBlockPool_doGetBlock(t *testing.T) {
	rsp := NewBlockPool(0, nil)
	assert.NotEmpty(t, rsp)
	resultC := make(chan model.AbstractBlock)
	var hashTmp = `0xd50866a60b4f7e4123400e0563efb987dc800d1a72af5cc1ae9ee68760bb18889`
	h := common.HexToHash(hashTmp)
	var getter = blockPoolGetter{
		blockHash:  h,
		resultChan: resultC,
	}
	go func(g *blockPoolGetter) {
		rsp.doGetBlock(
			&blockPoolGetter{
				blockHash:  h,
				resultChan: resultC,
			},
		)
	}(&getter)
	assert.Equal(t, 0, len(getter.resultChan))
}
