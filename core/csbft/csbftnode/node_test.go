package csbftnode

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestNewCsBft(t *testing.T) {
	assert.Panics(t, func() {
		NewCsBft(nil)
	})
}
