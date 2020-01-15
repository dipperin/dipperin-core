package minemaster

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_newServer(t *testing.T) {
	server := newServer(nil, nil, nil)
	assert.NotNil(t, server)
}

