package components

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"time"
)

func TestNewTimeoutTicker(t *testing.T) {
	assert.NotEmpty(t, NewTimeoutTicker())
}

func TestTimeoutTicker_OnStart(t *testing.T) {
	tt := &timeoutTicker{
		timer:    time.NewTimer(0),
		tickChan: make(chan TimeoutInfo, tickTockBufferSize),
		tockChan: make(chan TimeoutInfo, tickTockBufferSize),
	}
	assert.NoError(t, tt.OnStart())
}
