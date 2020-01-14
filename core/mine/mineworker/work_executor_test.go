package mineworker

import (
	"github.com/dipperin/dipperin-core/core/mine/minemsg"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewDefaultWorkExecutor(t *testing.T) {
	worker := &minemsg.DefaultWork{}
	workManager := &workManager{}
	executor := NewDefaultWorkExecutor(worker, workManager)
	assert.NotNil(t, executor)
}

//func Test_defaultWorkExecutor_ChangeNonce(t *testing.T) {
//	// test case
//	situations := []struct{
//		name string
//		given func() *defaultWorkExecutor
//		expectRes bool
//	}{
//		{
//			"meet the requirement",
//			func() *defaultWorkExecutor {
//				// init
//				blockHeader := model.Header{
//					Number:   100,
//					PreHash:  common.HexToHash("0x12312fa0929348"),
//					Diff:     common.HexToDiff("0x1effffff"),
//					CoinBase: common.HexToAddress("123"),
//					GasLimit: 1000000,
//					GasUsed:  1000,
//					Nonce:    common.EncodeNonce(100),
//					Bloom:    iblt.NewBloom(model.DefaultBlockBloomConfig),
//				}
//				work := minemsg.DefaultWork{
//					WorkerCoinbaseAddress: common.HexToAddress("123"),
//					BlockHeader:           blockHeader,
//					ResultNonce:           common.BlockNonce{},
//					RlpPreCal:             nil,
//				}
//				workExecutor := NewDefaultWorkExecutor(&work, nil)
//				// return
//				return workExecutor
//			},
//			true,
//		},
//		//{
//		//		"not meet the requirement",
//		//		func() defaultWorkExecutor {
//		//
//		//		},
//		//		false,
//		//},
//		//{
//		//			"cal hash failed",
//		//			func() defaultWorkExecutor {
//		//
//		//			},
//		//			false,
//		//},
//	}
//	// test
//	for _, situation := range situations {
//		executor := situation.given()
//		res := executor.ChangeNonce()
//		assert.Equal(t, situation.expectRes, res)
//	}
//
//}
