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
	"go.uber.org/zap"
	"time"

	"github.com/dipperin/dipperin-core/common/log"
	cmn "github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/core/csbft/model"
)

var (
	tickTockBufferSize = 10
)

type TimeoutInfo struct {
	Duration time.Duration       `json:"duration"`
	Height   uint64              `json:"height"`
	Round    uint64              `json:"round"`
	Step     model.RoundStepType `json:"step"`
}

// TimeoutTicker is a timer that schedules timeouts
// conditional on the height/round/step in the TimeoutInfo.
// The TimeoutInfo.Duration may be non-positive.
type TimeoutTicker interface {
	Start() error
	Stop()
	Chan() <-chan TimeoutInfo       // on which to receive a timeout
	ScheduleTimeout(ti TimeoutInfo) // reset the timer

	//SetLogger(logger log.Logger)
}

// timeoutTicker wraps time.Timer,
// scheduling timeouts only for greater height/round/step
// than what it's already seen.
// Timeouts are scheduled along the tickChan,
// and fired on the tockChan.
type timeoutTicker struct {
	cmn.BaseService

	timer    *time.Timer
	tickChan chan TimeoutInfo // for scheduling timeouts
	tockChan chan TimeoutInfo // for notifying about them
}

// NewTimeoutTicker returns a new TimeoutTicker.
func NewTimeoutTicker() TimeoutTicker {
	tt := &timeoutTicker{
		timer:    time.NewTimer(0),
		tickChan: make(chan TimeoutInfo, tickTockBufferSize),
		tockChan: make(chan TimeoutInfo, tickTockBufferSize),
	}
	tt.BaseService = *cmn.NewBaseService(nil, "TimeoutTicker", tt)
	tt.stopTimer() // don't want to fire until the first scheduled timeout
	return tt
}

// OnStart implements cmn.Service. It starts the timeout routine.
func (t *timeoutTicker) OnStart() error {
	//t.BaseService.OnStart()
	go t.timeoutRoutine()

	return nil
}

// OnStop implements cmn.Service. It stops the timeout routine.
func (t *timeoutTicker) OnStop() {
	//t.BaseService.OnStop()
	t.stopTimer()
}

// Chan returns a channel on which timeouts are sent.
func (t *timeoutTicker) Chan() <-chan TimeoutInfo {
	return t.tockChan
}

// ScheduleTimeout schedules a new timeout by sending on the internal tickChan.
// The timeoutRoutine is always available to read from tickChan, so this won't block.
// The scheduling may fail if the timeoutRoutine has already scheduled a timeout for a later height/round/step.
func (t *timeoutTicker) ScheduleTimeout(ti TimeoutInfo) {
	log.DLogger.Debug("Schedule Timeout", zap.Any("timeout info", ti), zap.Bool("is running", t.BaseService.IsRunning()))
	t.tickChan <- ti
}

//-------------------------------------------------------------
// stop the timer and drain if necessary
func (t *timeoutTicker) stopTimer() {
	// Stop() returns false if it was already fired or was stopped
	if !t.timer.Stop() {
		select {
		case <-t.timer.C:
		default:
			t.Logger.Debug("Timer already stopped")
		}
	}
}

// send on tickChan to start a new timer.
// timers are interupted and replaced by new ticks from later steps
// timeouts of 0 on the tickChan will be immediately relayed to the tockChan
func (t *timeoutTicker) timeoutRoutine() {
	var ti TimeoutInfo
	for {
		select {
		case newti := <-t.tickChan:
			log.DLogger.Debug("Received tick", zap.Any("old_ti", ti), zap.Any("new_ti", newti))
			// ignore tickers for old height/round/step
			if newti.Height < ti.Height {
				log.DLogger.Warn("height too low, ignore ticker")
				continue
			} else if newti.Height == ti.Height {
				if newti.Round < ti.Round {
					log.DLogger.Warn("round too low, ignore ticker")
					continue
				} else if newti.Round == ti.Round {
					if ti.Step > 0 && newti.Step < ti.Step {
						log.DLogger.Warn("step must higher than latest, ignore ticker", zap.Any("new step", newti.Step), zap.Any("cur step", ti.Step))
						continue
					}
				}
			}
			// stop the last timer
			t.stopTimer()
			// update TimeoutInfo and reset timer
			// NOTE time.Timer allows duration to be non-positive
			ti = newti
			t.timer.Reset(ti.Duration)
			log.DLogger.Info("ticker reset", zap.Duration("dur", ti.Duration), zap.Uint64("height", ti.Height), zap.Uint64("round", ti.Round), zap.Any("step", ti.Step))
		case <-t.timer.C:
			log.DLogger.Info("ticker timed out", zap.Duration("dur", ti.Duration), zap.Uint64("height", ti.Height), zap.Uint64("round", ti.Round), zap.Any("step", ti.Step))
			// go routine here guarantees timeoutRoutine doesn't block.
			// Determinism comes from playback in the receiveRoutine.
			// We can eliminate it by merging the timeoutRoutine into receiveRoutine
			//  and managing the timeouts ourselves with a millisecond ticker
			go func(toi TimeoutInfo) { t.tockChan <- toi }(ti)
		case <-t.Quit():
			return
		}
	}
}
