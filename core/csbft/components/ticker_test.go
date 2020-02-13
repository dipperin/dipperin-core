package components

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"time"
	"github.com/dipperin/dipperin-core/core/csbft/model"
)

func TestNewTimeoutTicker(t *testing.T) {
	assert.NotEmpty(t, NewTimeoutTicker())
}

func TestTimeoutTicker_OnStart(t *testing.T) {
	testCases := []struct {
		name   string
		given  func() bool
		expect bool
	}{
		{
			name: "TimeoutTicker OnStart true",
			given: func() bool {
				tt := &timeoutTicker{
					timer:    time.NewTimer(0),
					tickChan: make(chan TimeoutInfo, tickTockBufferSize),
					tockChan: make(chan TimeoutInfo, tickTockBufferSize),
				}
				return tt.OnStart() == nil
			},
			expect: true,
		},
		{
			name: "FetchBlock OnStart false",
			given: func() bool {
				tt := &timeoutTicker{
					timer:    time.NewTimer(0),
					tickChan: make(chan TimeoutInfo, tickTockBufferSize),
					tockChan: make(chan TimeoutInfo, tickTockBufferSize),
				}
				return tt.OnStart() != nil
			},
			expect: false,
		},
	}

	for i, tc := range testCases {
		sign := tc.given()
		if testCases[i].expect == sign {
			t.Log("success")
		} else {
			t.Logf("expect:%v,actual:%v", testCases[i].expect, sign)
		}
	}
}

func TestTimeoutTicker_OnStop(t *testing.T) {
	testCases := []struct {
		name   string
		given  func() bool
		expect bool
	}{
		{
			name: "TimeoutTicker OnStop true",
			given: func() bool {
				tt := &timeoutTicker{
					timer:    time.NewTimer(0),
					tickChan: make(chan TimeoutInfo, tickTockBufferSize),
					tockChan: make(chan TimeoutInfo, tickTockBufferSize),
				}
				return assert.NotPanics(t, tt.OnStop)
			},
			expect: true,
		},
	}

	for i, tc := range testCases {
		sign := tc.given()
		if testCases[i].expect == sign {
			t.Log("success")
		} else {
			t.Logf("expect:%v,actual:%v", testCases[i].expect, sign)
		}
	}
}

func TestTimeoutTicker_Chan(t *testing.T) {
	testCases := []struct {
		name   string
		given  func() bool
		expect bool
	}{
		{
			name: "TimeoutTicker Chan true",
			given: func() bool {
				tt := &timeoutTicker{
					timer:    time.NewTimer(0),
					tickChan: make(chan TimeoutInfo, tickTockBufferSize),
					tockChan: make(chan TimeoutInfo, tickTockBufferSize),
				}
				var tmp = TimeoutInfo{Duration: time.Second, Height: 1, Round: 1, Step: model.RoundStepNewHeight}
				tt.tockChan <- tmp
				c := <-tt.Chan()
				t.Log("c", c)
				return c == tmp
			},
			expect: true,
		},
	}

	for i, tc := range testCases {
		sign := tc.given()
		if testCases[i].expect == sign {
			t.Log("success")
		} else {
			t.Logf("expect:%v,actual:%v", testCases[i].expect, sign)
		}
	}
}

func TestTimeoutTicker_ScheduleTimeout(t *testing.T) {
	testCases := []struct {
		name   string
		given  func() bool
		expect bool
	}{
		{
			name: "TimeoutTicker ScheduleTimeout true",
			given: func() bool {
				tt := &timeoutTicker{
					timer:    time.NewTimer(0),
					tickChan: make(chan TimeoutInfo, tickTockBufferSize),
					tockChan: make(chan TimeoutInfo, tickTockBufferSize),
				}
				var tmp = TimeoutInfo{Duration: time.Second, Height: 1, Round: 1, Step: model.RoundStepNewHeight}
				tt.ScheduleTimeout(tmp)
				tc := <-tt.tickChan
				t.Log("tc:", tc)
				return tc == tmp
			},
			expect: true,
		},
	}

	for i, tc := range testCases {
		sign := tc.given()
		if testCases[i].expect == sign {
			t.Log("success")
		} else {
			t.Logf("expect:%v,actual:%v", testCases[i].expect, sign)
		}
	}
}

func TestTimeoutTicker_stopTimer(t *testing.T) {
	testCases := []struct {
		name   string
		given  func() bool
		expect bool
	}{
		{
			name: "TimeoutTicker stopTimer true",
			given: func() bool {
				tt := &timeoutTicker{
					timer:    time.NewTimer(0),
					tickChan: make(chan TimeoutInfo, tickTockBufferSize),
					tockChan: make(chan TimeoutInfo, tickTockBufferSize),
				}
				assert.NoError(t, tt.OnStart())
				return assert.NotPanics(t, tt.stopTimer)
			},
			expect: true,
		},
		{
			name: "TimeoutTicker stopTimer false",
			given: func() bool {
				tt := &timeoutTicker{
					timer:    time.NewTimer(0),
					tickChan: make(chan TimeoutInfo, tickTockBufferSize),
					tockChan: make(chan TimeoutInfo, tickTockBufferSize),
				}
				assert.NoError(t, tt.OnStart())
				assert.NotPanics(t, tt.OnStop)
				return assert.Panics(t, tt.stopTimer)
			},
			expect: true,
		},
	}

	for i, tc := range testCases {
		sign := tc.given()
		if testCases[i].expect == sign {
			t.Log("success")
		} else {
			t.Logf("expect:%v,actual:%v", testCases[i].expect, sign)
		}
	}
}

func TestTimeoutTicker_timeoutRoutine(t *testing.T) {
	testCases := []struct {
		name   string
		given  func() bool
		expect bool
	}{
		{
			name: "TimeoutTicker timeoutRoutine true",
			given: func() bool {
				tt := &timeoutTicker{
					timer:    time.NewTimer(0),
					tickChan: make(chan TimeoutInfo, tickTockBufferSize),
					tockChan: make(chan TimeoutInfo, tickTockBufferSize),
				}

				go tt.timeoutRoutine()
				return true
			},
			expect: true,
		},
	}

	for i, tc := range testCases {
		sign := tc.given()
		if testCases[i].expect == sign {
			t.Log("success")
		} else {
			t.Logf("expect:%v,actual:%v", testCases[i].expect, sign)
		}
	}
}
