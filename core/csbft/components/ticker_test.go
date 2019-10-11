// Copyright 2019, Keychain Foundation Ltd.
// This file is part of the dipperin-core library.
//
// The dipperin-core library is free software: you can redistribute
// it and/or modify it under the terms of the GNU Lesser General Public License
// as published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// The dipperin-core library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package components

import (
	"fmt"
	cmn "github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/core/csbft/model"
	"github.com/dipperin/dipperin-core/third-party/log"
	"testing"
	"time"
)

func TestNewTimeoutTicker(t *testing.T) {
	log.InitLogger(log.LvlDebug)
	tt := NewTimeoutTicker()
	tt.Start()
	go func() {
		select {
		//time.Sleep(time.Second * 25)
		case info := <-tt.Chan():
			fmt.Println("info", info, time.Now())
		}

	}()

	go func() {

		tt.ScheduleTimeout(TimeoutInfo{Duration: time.Duration(time.Millisecond * 10), Height: 1, Round: 2, Step: model.RoundStepNewRound})
		fmt.Println("one", time.Now())
		time.Sleep(time.Millisecond * 20)
		tt.ScheduleTimeout(TimeoutInfo{Duration: time.Duration(time.Millisecond * 10), Height: 1, Round: 2, Step: model.RoundStepNewRound})
	}()

	tt.Chan()
	tt.Stop()

}

func Test_timeoutRoutine(t *testing.T) {
	tt := &timeoutTicker{
		timer:    time.NewTimer(0),
		tickChan: make(chan TimeoutInfo, tickTockBufferSize),
		tockChan: make(chan TimeoutInfo, tickTockBufferSize),
	}
	tt.BaseService = *cmn.NewBaseService(nil, "TimeoutTicker", tt)
	tt.Start()

	tt.tickChan <- TimeoutInfo{Height: 0, Round: 10}
	go tt.timeoutRoutine()
	time.Sleep(time.Millisecond * 5)
	tt.Stop()

}

func Test_timeoutRoutine2(t *testing.T) {
	tt := &timeoutTicker{
		timer:    time.NewTimer(0),
		tickChan: make(chan TimeoutInfo, tickTockBufferSize),
		tockChan: make(chan TimeoutInfo, tickTockBufferSize),
	}
	tt.BaseService = *cmn.NewBaseService(nil, "TimeoutTicker", tt)
	tt.Start()
	//tt.stopTimer() // don't want to fire until the first scheduled timeout
	//ticker := NewTimeoutTicker()
	tt.tickChan <- TimeoutInfo{Height: 100, Round: 15}
	time.Sleep(time.Millisecond * 2)
	tt.tickChan <- TimeoutInfo{Height: 50, Round: 15}
	go tt.timeoutRoutine()
	time.Sleep(time.Millisecond * 5)
	tt.Stop()

}

func TestTimeoutTicker_ScheduleTimeout(t *testing.T) {
	ticker := NewTimeoutTicker()
	ticker.Start()
	ticker.ScheduleTimeout(TimeoutInfo{Duration: 3 * time.Millisecond, Height: uint64(1), Round: uint64(1), Step: 1})
	ticker.ScheduleTimeout(TimeoutInfo{Duration: 3 * time.Millisecond, Height: uint64(1), Round: uint64(1), Step: 2})
	time.Sleep(time.Millisecond * 6)
	ticker.Stop()
}
