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

func Test_defaultWorkExecutor_ChangeNonce(t *testing.T) {
	// todo
}
