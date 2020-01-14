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

package gevent

import (
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/ethereum/go-ethereum/event"
)

const (
	NewBlockInsertEvent = "on_new_block"
)

type Subscription interface {
	Err() <-chan error // returns the error channel
	Unsubscribe()      // cancels sending of events, closing the error channel
}

var events = map[string]*event.Feed{}

// Add event in global
func Add(name string) {
	if events[name] != nil {
		if util.IsTestEnv() {
			return
		}
		panic("add event already exist: " + name)
	}
	events[name] = &event.Feed{}
}

// Send event
func Send(to string, v interface{}) (sentN int) {
	if e := events[to]; e == nil {
		panic("send event called, but not add: " + to)
	} else {
		return e.Send(v)
	}
}

// Subscribe event
func Subscribe(name string, c interface{}) Subscription {
	if e := events[name]; e == nil {
		panic("subscribe event called, but not add: " + name)
	} else {
		return e.Subscribe(c)
	}
}
