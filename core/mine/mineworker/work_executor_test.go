package mineworker

import (
	"github.com/dipperin/dipperin-core/common"
	iblt "github.com/dipperin/dipperin-core/core/bloom"
	"github.com/dipperin/dipperin-core/core/mine/minemsg"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewDefaultWorkExecutor(t *testing.T) {
	worker := &minemsg.DefaultWork{}
	workManager := &workManager{}
	executor := NewDefaultWorkExecutor(worker, workManager)
	assert.NotNil(t, executor)
}

func Test_defaultWorkExecutor_ChangeNonce(t *testing.T) {
	// test case
	situations := []struct{
		name string
		given func() *defaultWorkExecutor
		expectRes bool
	}{
		{
			"meet the requirement",
			func() *defaultWorkExecutor {
				// init
				blockHeader := model.Header{
					Number:   100,
					PreHash:  common.HexToHash("0x000004e9937f5754db7e3953189befb1938fa11f57afbbc995a0c44beb08adfe"),
					Diff:     common.HexToDiff("0x1e17f011"),
					CoinBase: common.HexToAddress("0x00004f011d62285f527ce47D189EBA821601dAf8A16B"),
					GasLimit: 1000000,
					GasUsed:  1000,
					Nonce:    common.EncodeNonce(563001),
					Bloom:    iblt.NewBloom(iblt.NewBloomConfig(1<<12, 4)),
				}
				work := minemsg.DefaultWork{
					WorkerCoinbaseAddress: common.HexToAddress("0x00004f011d62285f527ce47D189EBA821601dAf8A16B"),
					BlockHeader:           blockHeader,
					ResultNonce:           common.BlockNonce{},
					RlpPreCal:             nil,
				}
				work.RlpPreCal = work.BlockHeader.RlpBlockWithoutNonce()
				workExecutor := NewDefaultWorkExecutor(&work, nil)
				// return
				return workExecutor
			},
			false,		// todo: set the condition which meet the requirement
		},
		{
			"not meet the requirement",
			func() *defaultWorkExecutor {
				// init
				blockHeader := model.Header{
					Number:   100,
					PreHash:  common.HexToHash("0x12312fa0929348"),
					Diff:     common.HexToDiff("0x1effffff"),
					CoinBase: common.HexToAddress("123"),
					GasLimit: 1000000,
					GasUsed:  1000,
					Nonce:    common.EncodeNonce(100),
					Bloom:    iblt.NewBloom(iblt.NewBloomConfig(1<<12, 4)),
				}
				work := minemsg.DefaultWork{
					WorkerCoinbaseAddress: common.HexToAddress("123"),
					BlockHeader:           blockHeader,
					ResultNonce:           common.BlockNonce{},
					RlpPreCal:             nil,
				}
				work.RlpPreCal = work.BlockHeader.RlpBlockWithoutNonce()
				workExecutor := NewDefaultWorkExecutor(&work, nil)
				// return
				return workExecutor
			},
			false,
		},
		{
				"cal hash failed",
				func() *defaultWorkExecutor {
					// init
					blockHeader := model.Header{
						Number:   100,
						PreHash:  common.HexToHash("0x12312fa0929348"),
						Diff:     common.HexToDiff("0x1effffff"),
						CoinBase: common.HexToAddress("123"),
						GasLimit: 1000000,
						GasUsed:  1000,
						Nonce:    common.EncodeNonce(100),
						Bloom:    iblt.NewBloom(model.DefaultBlockBloomConfig),
					}
					work := minemsg.DefaultWork{
						WorkerCoinbaseAddress: common.HexToAddress("123"),
						BlockHeader:           blockHeader,
						ResultNonce:           common.BlockNonce{},
						RlpPreCal:             nil,
					}
					workExecutor := NewDefaultWorkExecutor(&work, nil)
					// return
					return workExecutor
				},
				false,
		},
	}
	// test
	for _, situation := range situations {
		executor := situation.given()
		res := executor.ChangeNonce()
		assert.Equal(t, situation.expectRes, res, situation.name)
	}

}
