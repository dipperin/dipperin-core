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

package mineworker

import (
	"github.com/dipperin/dipperin-core/common/g-metrics"
	"github.com/dipperin/dipperin-core/common/log"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"sync"
	"time"
)

const (
	//newWorkReceiveTimeOut = 1 * time.Second
	newWorkReceiveTimeOut = 500 * time.Millisecond
)

func NewMiner() *defaultMiner {
	m := &defaultMiner{
		newWorkChan: make(chan workExecutor),
	}
	return m
}

// the miner supports various mining algorithms and data, so only workExecutor is contained here. If algorithm or data is changed, then different workExecutor is inserted.
type defaultMiner struct {
	stopChan    chan struct{}
	newWorkChan chan workExecutor
	// the mission ongoing
	curWork workExecutor
	lock    sync.Mutex

	metricTimer *prometheus.Timer
	mineStartAt time.Time
}

func (miner *defaultMiner) receiveWork(work workExecutor) {
	t := time.NewTimer(newWorkReceiveTimeOut)
	select {
	case miner.newWorkChan <- work:
		t.Stop()
		// set receive timeout
	case <-t.C:
		log.DLogger.Info("receive new work time out, miner maybe stopped")
	}
}

func (miner *defaultMiner) startMine() {
	go miner.loop()
}

func (miner *defaultMiner) stopMine() {
	miner.lock.Lock()
	defer miner.lock.Unlock()

	if !util.StopChanClosed(miner.stopChan) {
		close(miner.stopChan)
	}

	miner.curWork = nil
}

func (miner *defaultMiner) loop() {
	// only one loop is authorized
	miner.lock.Lock()

	if !miner.stopped() {
		miner.lock.Unlock()

		log.DLogger.Info("call start miner, but miner already started")
		return
	}
	log.DLogger.Info("miner start mine loop")
	// reset stop chan
	miner.stopChan = make(chan struct{})

	miner.lock.Unlock()
out:
	for {
		// if two work both solve the puzzle and arrives consecutively,
		// then master will start 2 timers.
		select {
		case miner.curWork = <-miner.newWorkChan:
			miner.mineStartAt = time.Now()
			miner.metricTimer = g_metrics.NewTimer(g_metrics.FindNonceDuration)
			log.DLogger.Info("miner receive new work1")
		case <-miner.stopChan:
			log.DLogger.Info("stop mine")
			break out
		default:
			miner.doMine()
		}
	}
}

func (miner *defaultMiner) doMine() {
	//log.DLogger.Debug("miner do mine")
	if miner.curWork == nil {
		miner.waitNewWork()
	}
	// Submit if it is discovered, and wait for a new task
	if miner.curWork != nil && miner.curWork.ChangeNonce() {
		log.DLogger.Info("miner found nonce", zap.Duration("use time", time.Now().Sub(miner.mineStartAt)))
		miner.curWork.Submit()
		if miner.metricTimer != nil {
			miner.metricTimer.ObserveDuration()
			miner.metricTimer = nil
		}
		miner.waitNewWork()
	}
}

func (miner *defaultMiner) waitNewWork() {
	//log.DLogger.Debug("miner waitNewWork")
	select {
	case miner.curWork = <-miner.newWorkChan:
		miner.mineStartAt = time.Now()
		miner.metricTimer = g_metrics.NewTimer(g_metrics.FindNonceDuration)
		//log.DLogger.Info("miner receive new work2")
	case <-miner.stopChan:
	}
}

// check whether mining is stopped
func (miner *defaultMiner) stopped() bool {
	// if stop chan is closed, it means the mining is stopped
	return util.StopChanClosed(miner.stopChan)
}
