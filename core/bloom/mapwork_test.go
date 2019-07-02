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

package iblt

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

type addRet struct {
	v  int
	ok bool
}

func (ret addRet) DoTask(i interface{}) interface{} {
	result := addInt(i.(int))
	return result
}

func addInt(i int) addRet {
	time.Sleep(time.Microsecond * 1)
	return addRet{i, true}
}

func Test_Map_Work(t *testing.T) {
	LENGTH := 40000
	st := time.Now()
	mWorkMap := InitWorkMap(runtime.NumCPU())

	ret := addRet{0, false}
	mWorkMap.SetOperation(ret)

	raw := make([]interface{}, LENGTH)
	after := make([]addRet, LENGTH)

	for i := 0; i < LENGTH; i++ {
		raw[i] = i
	}

	retCh, err := mWorkMap.StartWorks(raw)

	j := 0

	if err == nil {
	pass:
		for {
			select {
			case r := <-retCh:
				after[j] = r.(addRet)
				j++
				//fmt.Println(r, j)
				if j == LENGTH {
					//fmt.Println("work processed")
					break pass
				}
			}
		}
	}

	fmt.Println(time.Now().Sub(st))

	st = time.Now()
	for i := 0; i < LENGTH; i++ {
		after[i] = addInt(i)
	}
	fmt.Println(time.Now().Sub(st))
}
