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
		timer:    time.NewTimer(time.Second),
		tickChan: make(chan TimeoutInfo, tickTockBufferSize),
		tockChan: make(chan TimeoutInfo, tickTockBufferSize),
	}
	assert.NoError(t, tt.OnStart())
}

func TestTimeoutTicker_OnStop(t *testing.T) {
	tt := &timeoutTicker{
		timer:    time.NewTimer(time.Second),
		tickChan: make(chan TimeoutInfo, tickTockBufferSize),
		tockChan: make(chan TimeoutInfo, tickTockBufferSize),
	}
	if assert.NoError(t, tt.OnStart()) {
		time.Sleep(500 * time.Millisecond)
		assert.NotPanics(t, tt.OnStop)
	}
}

func TestTimeoutTicker_Chan(t *testing.T) {
	tt := NewTimeoutTicker()
	assert.NotEmpty(t, tt)
	assert.NotPanics(t, func() {
		tt.Chan()
	})
}

func TestTimeoutTicker_ScheduleTimeout(t *testing.T) {
	tt := NewTimeoutTicker()
	assert.NotEmpty(t, tt)
	assert.NotPanics(t, func() {
		var w = TimeoutInfo{
			Duration: 500 * time.Millisecond,
			Height:   1,
			Round:    1,
			Step:     1,
		}
		tt.ScheduleTimeout(w)
	})
}

func TestTimeoutTicker_stopTimer(t *testing.T) {
	tt := &timeoutTicker{
		timer:    time.NewTimer(time.Second),
		tickChan: make(chan TimeoutInfo, tickTockBufferSize),
		tockChan: make(chan TimeoutInfo, tickTockBufferSize),
	}
	assert.NotPanics(t, tt.stopTimer)
}

func TestTimeoutTicker_timeoutRoutine(t *testing.T) {
	tt := &timeoutTicker{
		timer:    time.NewTimer(time.Second),
		tickChan: make(chan TimeoutInfo, tickTockBufferSize),
		tockChan: make(chan TimeoutInfo, tickTockBufferSize),
	}
	assert.NotPanics(t, func() {
		go tt.timeoutRoutine()
	})
}
