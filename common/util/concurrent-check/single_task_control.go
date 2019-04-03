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

package concurrent_check

import (
	"github.com/dipperin/dipperin-core/third-party/log"
	"sync"
	"sync/atomic"
)

type SingleTaskControl struct {
	isDoing int32
	count int
	lock sync.Mutex
}

// Return whether swapped ， true represent can do
func (stc *SingleTaskControl) TryDo(count int) bool {
	if count == 0 {
		return false
	}
	if atomic.CompareAndSwapInt32(&stc.isDoing, 0, 1) {
		stc.count = count
		return true
	}
	log.Warn("try start task, but last task haven't done", "remain count", stc.count)
	return false
}

func (stc *SingleTaskControl) Done() {
	stc.lock.Lock()
	defer stc.lock.Unlock()

	stc.count--
	if stc.count <=0 {
		if !atomic.CompareAndSwapInt32(&stc.isDoing, 1, 0) {
			log.Warn("try call task done, but not started")
		}
	}
}