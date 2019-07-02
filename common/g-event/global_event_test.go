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

package g_event

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAdd(t *testing.T) {
	finished := false
	go func() {
		time.Sleep(2 * time.Second)
		if !finished {
			panic("use too much time")
		}
	}()

	tEvent := "test"
	Add(tEvent)
	assert.Len(t, events, 1)

	sc1 := make(chan int)
	s := Subscribe(tEvent, sc1)
	go func() {
		x := <-sc1
		assert.Equal(t, 1, x)
		s.Unsubscribe()
	}()

	n := Send(tEvent, 1)
	assert.Equal(t, 1, n)

	time.Sleep(10 * time.Millisecond)
	n = Send(tEvent, 1)
	assert.Equal(t, 0, n)

	finished = true
}

func TestAddPanic(t *testing.T) {
	tEvent := "test"
	Add(tEvent)
	assert.Len(t, events, 1)

	//add again
	Add(tEvent)
}

func TestSubscribePanic(t *testing.T) {
	tEvent := "test1"
	assert.Panics(t, func() {
		Subscribe(tEvent, nil)
	})
}

func TestSendPanic(t *testing.T) {
	tEvent := "test2"
	assert.Panics(t, func() {
		Send(tEvent, 1)
	})
}
