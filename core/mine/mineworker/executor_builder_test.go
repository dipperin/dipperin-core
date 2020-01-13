package mineworker

import (
	iblt "github.com/dipperin/dipperin-core/core/bloom"
	"github.com/dipperin/dipperin-core/core/mine/minemsg"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewDefaultExecutorBuilder(t *testing.T) {
	executor := NewDefaultExecutorBuilder()
	assert.NotNil(t, executor)
}

func Test_defaultExecutorBuilder_CreateExecutor(t *testing.T) {
	// init
	executor := NewDefaultExecutorBuilder()

	// test case
	situations := []struct{
		name string
		given func() (msg workMsg, workCount int, submitter workSubmitter)
		expectResultLen int
		expectErr error
	}{
		{
			"unknown msg code",
			func() (msg workMsg, workCount int, submitter workSubmitter) {
				msg = &localWorkMsg{}
				workCount = 2
				submitter = &workManager{}
				return
			},
			0,
			UnknownMsgCodeErr,
		},
		{
				"normal situation",
				func() (msg workMsg, workCount int, submitter workSubmitter) {
					msg = &localWorkMsg{
						code: minemsg.NewDefaultWorkMsg,
						work: &minemsg.DefaultWork{BlockHeader: model.Header{Bloom: iblt.NewBloom(iblt.NewBloomConfig(1<<12, 4))}},
					}
					workCount = 2
					submitter = &workManager{}
					return
				},
				2,
				nil,
		},
	}
	// test
	for _, situation := range situations {
		msg, workCount, submitter := situation.given()
		workExecutor, workErr := executor.CreateExecutor(msg, workCount, submitter)
		// result
		assert.Equal(t, situation.expectResultLen, len(workExecutor))
		assert.Equal(t, situation.expectErr, workErr)
	}
}