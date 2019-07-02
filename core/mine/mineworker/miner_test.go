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
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type fakeWork struct {
	id             int
	count          int
	submitWorkChan chan int
}

func (fw *fakeWork) ChangeNonce() bool {
	//log.Debug("fake work change nonce", "cur count", fw.count)
	if fw.count < 10 {
		fw.count++
		return false
	}
	return true
}

func (fw *fakeWork) Submit() {
	//log.Debug("fakeWork submit", "id", fw.id, "fw", fw)
	fw.submitWorkChan <- fw.id
}

func (fw *fakeWork) waitForSubmit(t *testing.T) {
	select {
	case submitId := <-fw.submitWorkChan:
		assert.Equal(t, fw.id, submitId)
	case <-time.After(1 * time.Second):
		t.Fatal("mine time out", "id", fw.id)
	}
}

func TestNewDefaultMiner(t *testing.T) {
	m := NewMiner()
	fw := &fakeWork{id: 0, submitWorkChan: make(chan int)}
	// if the mining is not started, the sent of work to miners will not cause block
	m.receiveWork(fw)

	// a normal mining
	fw1 := &fakeWork{id: 1, submitWorkChan: make(chan int)}
	m.startMine()
	// wait start ok
	time.Sleep(100 * time.Millisecond)
	m.receiveWork(fw1)

	fw1.waitForSubmit(t)

	// test that it can start normally after stop
	m.stopMine()
	// wait for stop
	time.Sleep(100 * time.Millisecond)
	// confirm stop
	m.receiveWork(fw)

	m.startMine()
	time.Sleep(100 * time.Millisecond)

	for i := 20; i < 50; i++ {
		fwx := &fakeWork{id: i, submitWorkChan: make(chan int)}
		m.receiveWork(fwx)
		fwx.waitForSubmit(t)
	}
}

func TestReadFromAClosedChannel(t *testing.T) {
	tmp := make(chan struct{})
	close(tmp)
	if _, ok := <-tmp; !ok {
		log.Debug("channel closed")
	}
}

func TestCanNotStartTwice(t *testing.T) {
	m := NewMiner()
	m.startMine()
	m.startMine()
	time.Sleep(50 * time.Millisecond)
}
