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

package gtimer

import (
	"fmt"
	"testing"
	"time"
)

type si struct {
	no int
}

type inf interface {
	show()
}

func (s *si) show() {
	fmt.Println(s.no)
}

func change(in inf, v int) {
	m, ok := in.(*si)
	if ok {
		n := m
		n.no = v
	} else {
		fmt.Println("not si")
	}
}

func TestSetPeriodAndRun(t *testing.T) {
	exp := si{3}

	id := SetPeriodAndRun(exp.show, time.Millisecond)
	fmt.Println("timer id:", id)

	//Analog stop timer logic, no need end the timer
	stop := make(chan bool)
	go func() {
		for i := 0; i < 4; i++ {
			time.Sleep(time.Millisecond)
			change(&exp, i)
		}
		stop <- true
	}()

	//wait to exit
	for {
		select {
		case <-stop:
			StopWork(id)
			fmt.Println("exit")
			return
		}
	}
}
