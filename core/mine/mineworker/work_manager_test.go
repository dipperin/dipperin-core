package mineworker

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_newWorkManager(t *testing.T) {
	manager := newWorkManager(nil, nil, nil)
	assert.NotNil(t, manager)
}