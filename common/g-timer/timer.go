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

package g_timer

import (
	"sync"
	"time"
)

var (
//errExisted = errors.New("timer existed")
//errNotExisted = errors.New("timer not existed")
)

type property struct {
	ticker *time.Ticker
	timer  *time.Timer
	stop   chan struct{}
}

var recordID int
var recordLock sync.RWMutex
var record = make(map[int]property)

/*
 * SetPeriodAndRun create period timer and run handle periodically , return timer id
 * handle : user function
 * interval: time period
 */
func SetPeriodAndRun(handle func(), interval time.Duration) int {
	recordLock.Lock()
	defer recordLock.Unlock()

	recordID++

	ticker := time.NewTicker(interval)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				handle()

			case <-quit:
				return
			}
		}
	}()
	record[recordID] = property{ticker, nil, quit}
	return recordID
}

/*
 * StopWork stop periodic timer
 */
func StopWork(id int) {
	recordLock.Lock()
	defer recordLock.Unlock()

	if tm, ok := record[id]; ok {
		close(tm.stop)
		delete(record, id)
	}
}
