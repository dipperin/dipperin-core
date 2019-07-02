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
	"fmt"
	"runtime/debug"
	"sync/atomic"
)

type OnlyOneGo struct {
	entered   int32
	lastStack []byte
}

func (oog *OnlyOneGo) Enter() {
	// If there is no swaped, it is already 1
	if !atomic.CompareAndSwapInt32(&oog.entered, 0, 1) {
		fmt.Println("last stack =======", string(oog.lastStack))
		fmt.Println("cur stack =======", string(debug.Stack()))
		panic("duplicate enter for one go check")
	}
	oog.lastStack = debug.Stack()
}

func (oog *OnlyOneGo) Quit() {
	if !atomic.CompareAndSwapInt32(&oog.entered, 1, 0) {
		return
	}
	oog.lastStack = []byte{}
}
