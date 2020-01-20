package mineworker

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_newWorker(t *testing.T) {
	worker := newWorker(common.HexToAddress("123"), 2, nil)
	assert.NotNil(t, worker)
	assert.Equal(t, 2, len(worker.miners))
}

func Test_worker_CurrentCoinbaseAddress(t *testing.T) {
	addr := common.HexToAddress("123")
	worker := newWorker(addr, 2, nil)
	coinbaseAddr := worker.CurrentCoinbaseAddress()
	assert.Equal(t, coinbaseAddr, worker.coinbaseAddress.Load())
}

func Test_worker_Miners(t *testing.T) {
	worker := newWorker(common.HexToAddress("123"), 2, nil)
	assert.NotNil(t, worker)
	assert.Equal(t, 2, len(worker.Miners()))
}

func Test_worker_SetCoinbaseAddress(t *testing.T) {
	addr := common.HexToAddress("123")
	worker := newWorker(addr, 2, nil)
	worker.SetCoinbaseAddress(addr)
	assert.Equal(t, addr, worker.coinbaseAddress.Load())
}

func Test_worker_Start(t *testing.T) {

}

func Test_worker_Stop(t *testing.T) {

}

func Test_worker_register(t *testing.T) {

}

func Test_worker_unRegister(t *testing.T) {

}
