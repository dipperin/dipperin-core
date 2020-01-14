package minemsg

import (
	"encoding/binary"
	"errors"
	"github.com/dipperin/dipperin-core/common"
	iblt "github.com/dipperin/dipperin-core/core/bloom"
	"github.com/dipperin/dipperin-core/core/model"
	model_mock "github.com/dipperin/dipperin-core/tests/mock/model"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultWork_CalBlockRlpWithoutNonce(t *testing.T) {
	// init
	defaultWork := &DefaultWork{
		WorkerCoinbaseAddress: common.Address{},
		BlockHeader:           model.Header{Bloom: iblt.NewBloom(iblt.NewBloomConfig(1<<12, 4))},
		ResultNonce:           common.BlockNonce{},
		RlpPreCal:             nil,
	}

	// test
	defaultWork.CalBlockRlpWithoutNonce()
	assert.Equal(t, defaultWork.BlockHeader.RlpBlockWithoutNonce(), defaultWork.RlpPreCal)
}

func TestDefaultWork_CalHash(t *testing.T) {
	// test cases
	situations := []struct {
		name         string
		given        func() *DefaultWork
		expectResult common.Hash
		expectErr    error
	}{
		{
			"error of no rlp to be calculate",
			func() *DefaultWork {
				return &DefaultWork{
					BlockHeader: model.Header{Bloom: iblt.NewBloom(iblt.NewBloomConfig(1<<12, 4))},
				}
			},
			common.Hash{},
			errors.New("DefaultWork rlp be not calculated yet"),
		},
		{
			"normal case",
			func() *DefaultWork {
				defaultWork := &DefaultWork{
					BlockHeader: model.Header{Bloom: iblt.NewBloom(iblt.NewBloomConfig(1<<12, 4))},
				}
				defaultWork.RlpPreCal = defaultWork.BlockHeader.RlpBlockWithoutNonce()
				return defaultWork
			},
			common.HexToHash("c399c04115218674f164adc6a7ed0c7dff144587ca9ae7cf445bd8e03331cfca"),
			nil,
		},
	}
	// test
	for _, situation := range situations {
		dWork := situation.given()
		hash, err := dWork.CalHash()
		// expect result
		if situation.expectErr != nil {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
		assert.Equal(t, situation.expectResult.String(), hash.String())
	}
}

func TestDefaultWork_FillSealResult(t *testing.T) {
	// init
	defaultWork := &DefaultWork{
		WorkerCoinbaseAddress: common.Address{},
		BlockHeader:           model.Header{Number: 123},
		ResultNonce:           common.EncodeNonce(124),
		RlpPreCal:             nil,
	}

	// test case
	situations := []struct {
		name      string
		given     func() model.AbstractBlock
		expectErr error
	}{
		{
			"error of not in current block",
			func() model.AbstractBlock {
				// mock
				ctrl := gomock.NewController(t)
				mockBlock := model_mock.NewMockAbstractBlock(ctrl)
				// reject
				mockBlock.EXPECT().Number().Return(uint64(137)).AnyTimes()
				mockBlock.EXPECT().Nonce().Return(common.EncodeNonce(124)).AnyTimes()
				mockBlock.EXPECT().SetNonce(gomock.Any()).Return().AnyTimes()
				return mockBlock
			},
			errors.New("not nil"),
		},
		{
			"normal situation of fill result to current block",
			func() model.AbstractBlock {
				// mock
				ctrl := gomock.NewController(t)
				mockBlock := model_mock.NewMockAbstractBlock(ctrl)
				// reject
				mockBlock.EXPECT().Number().Return(uint64(123)).AnyTimes()
				mockBlock.EXPECT().Nonce().Return(common.EncodeNonce(124)).AnyTimes()
				mockBlock.EXPECT().SetNonce(gomock.Any()).Return().AnyTimes()
				return mockBlock
			},
			nil,
		},
	}
	// test
	for _, situation := range situations {
		absBlock := situation.given()
		err := defaultWork.FillSealResult(absBlock)
		// expect result
		if situation.expectErr != nil {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
		// check current block
		assert.Equal(t, defaultWork.ResultNonce, absBlock.Nonce())
	}
}

func TestDefaultWork_GetWorkerCoinbaseAddress(t *testing.T) {
	defaultWork := &DefaultWork{
		WorkerCoinbaseAddress: common.HexToAddress("test address"),
		BlockHeader:           model.Header{},
		ResultNonce:           common.BlockNonce{},
		RlpPreCal:             nil,
	}
	// test
	resAddress := defaultWork.GetWorkerCoinbaseAddress()
	assert.Equal(t, defaultWork.WorkerCoinbaseAddress.String(), resAddress.String())
}

func TestDefaultWork_SetWorkerCoinbaseAddress(t *testing.T) {
	defaultWork := &DefaultWork{
		WorkerCoinbaseAddress: common.Address{},
		BlockHeader:           model.Header{},
		ResultNonce:           common.BlockNonce{},
		RlpPreCal:             nil,
	}
	assert.Equal(t, "0x00000000000000000000000000000000000000000000", defaultWork.WorkerCoinbaseAddress.String())
	// test
	setAddress := common.HexToAddress("test address")
	defaultWork.SetWorkerCoinbaseAddress(setAddress)
	assert.Equal(t, setAddress.String(), defaultWork.WorkerCoinbaseAddress.String())
}

func TestDefaultWork_Split(t *testing.T) {
	w := &DefaultWork{BlockHeader: model.Header{}}
	works := w.Split(5)
	for i, w := range works {
		x := binary.BigEndian.Uint32(w.BlockHeader.Nonce[4:8])
		assert.Equal(t, uint32(i), x)
	}
}
