// Copyright 2018 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package mclock

import (
	"testing"
	"time"
)

func TestSimulated_Run(t *testing.T) {
	s := &Simulated{}
	s.insert(1, func() {})
	s.Now()
	s.After(1)
	go s.Sleep(1)
	s.WaitForTimers(1)
	s.Run(1)
	s.ActiveTimers()
	time.Sleep(100 * time.Millisecond)
}
