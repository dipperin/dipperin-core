package minemsg

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/model"
	model_mock "github.com/dipperin/dipperin-core/tests/mock/model-mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMakeDefaultWorkBuilder(t *testing.T) {
	builder := MakeDefaultWorkBuilder()
	assert.NotNil(t, builder)
}

func TestDefaultWorkBuilder_BuildWorks(t *testing.T) {
	// init
	builder := MakeDefaultWorkBuilder()
	assert.NotNil(t, builder)

	// test situations
	situations := []struct{
		name string
		given func() (block model.AbstractBlock, workLen int)
		expectCode int
		expectWorkLens int
	}{
		{
			"normal new block",
			func() (model.AbstractBlock, int) {
				// mock
				ctrl := gomock.NewController(t)
				mockBlock := model_mock.NewMockAbstractBlock(ctrl)
				// reject
				mockBlock.EXPECT().Header().Return(func() *model.Header {
					nonce := common.EncodeNonce(137)
					return &model.Header{Nonce: nonce}
				}())
				return mockBlock, 10
			},
			16,
			10,
		},
		{
				"invalid nil block",
				func() (model.AbstractBlock, int) {
					return nil, 10
				},
				0,
				0,
		},
	}
	// test
	for _, situation := range situations {
		block, workLen := situation.given()
		code, works := builder.BuildWorks(block, workLen)
		assert.Equal(t, situation.expectCode, code)
		assert.Equal(t, situation.expectWorkLens, len(works))
	}
}