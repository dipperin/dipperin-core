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

package model

import (
	"sync"
)

const (
	txWaterLevel = 10
)

//singleton
var txConcurrentCacher *TxCacher = nil

type TxCacherRequest struct {
	txs []AbstractTransaction
	inc int
}

// txCacher is a helper structure to concurrently recover transaction
// id from digital signatures on background threads.
type TxCacher struct {
	threads int
	tasks   chan *TxCacherRequest
	wg      sync.WaitGroup
}

// NewTxCacher creates a new transaction id background cacher and starts
// as many processing goroutines as allowed by the GOMAXPROCS on construction.
func NewTxCacher(threads int) *TxCacher {
	if txConcurrentCacher != nil {
		return txConcurrentCacher
	}

	cacher := &TxCacher{
		tasks:   make(chan *TxCacherRequest, threads),
		threads: threads,
	}
	for i := 0; i < threads; i++ {
		go cacher.cache(i)
	}

	txConcurrentCacher = cacher

	return cacher
}

// cache is an infinite loop, caching transaction id from various forms of
// data structures.
func (cacher *TxCacher) cache(no int) {
	for task := range cacher.tasks {
		length := len(task.txs)
		for i := 0; i < length; i += task.inc {
			if task.txs[i] == nil {
				break
			}
			task.txs[i].CalTxId()
		}
		cacher.wg.Done()
	}
}

// TxRecover recovers the id from a batch of transactions and caches them
// back into the same data structures. There is no validation being done, nor
// any reaction to invalid signatures. That is up to calling code later.
func (cacher *TxCacher) TxRecover(txs []AbstractTransaction) {
	if len(txs) == 0 {
		return
	}
	// Ensure we have meaningful task sizes and schedule the recoveries
	tasks := cacher.threads
	if len(txs) < tasks*txWaterLevel {
		tasks = 1
	}
	for i := 0; i < tasks; i++ {
		cacher.wg.Add(1)
		cacher.tasks <- &TxCacherRequest{
			txs: txs[i:],
			inc: tasks,
		}
	}

	cacher.wg.Wait()
}

func (cacher *TxCacher) StopTxCacher() {
	close(cacher.tasks)
	txConcurrentCacher = nil
}
