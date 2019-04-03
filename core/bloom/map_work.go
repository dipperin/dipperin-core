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
	"errors"
	"github.com/dipperin/dipperin-core/third-party/log"
	"runtime"
)

var (
	ErrNoOpFunc = errors.New("not set operated function")
	ErrNoData   = errors.New("not provide data")
)

type Operation interface {
	DoTask(i interface{}) interface{}
}

type MapWork struct {
	thdNum int
	op     Operation
	//result chan interface{}
	done   chan WorkResp
	inputs chan WorkReq
	stop   chan struct{}
	abort  chan struct{}
}

type WorkReq struct {
	index int
	req   interface{}
}

type WorkResp struct {
	index  int
	result interface{}
}

var (
	WorkMap = InitWorkMap(runtime.NumCPU())
)

func InitWorkMap(tNum int) *MapWork {
	work := MapWork{thdNum: tNum}

	work.done = make(chan WorkResp, tNum)
	work.inputs = make(chan WorkReq, tNum)
	//work.result = make(chan interface{}, tNum)
	work.stop = make(chan struct{}, tNum)
	work.abort = make(chan struct{})

	//start threads
	for i := 0; i < tNum; i++ {
		go work.workThread()
	}

	return &work
}

func (work *MapWork) SetOperation(task Operation) {
	work.op = task
}

func (work *MapWork) workThread() {
	defer func() {
		if r := recover(); r != nil {
			log.Debug(fmt.Sprintf("recover err：%s\n", r))
		}
	}()

	for {
		select {
		case in := <-work.inputs:
			ret := work.op.DoTask(in.req)
			resp := WorkResp{in.index, ret}
			work.done <- resp

		case <-work.stop:
			return
		}
	}
}

func (work *MapWork) StartWorks(items []interface{}) (chan interface{}, error) {
	var (
		in, out = 0, 0
		checked = make([]bool, len(items))
		inputs  = work.inputs
		cache   = make([]interface{}, len(items))
	)

	length := len(items)

	if length < 1 {
		return nil, ErrNoData
	}

	if work.op == nil {
		return nil, ErrNoOpFunc
	}
	result := make(chan interface{}, length)

	go func() {
		item := WorkReq{0, items[0]}
		for {
			select {
			case inputs <- item:
				if in++; in == length {
					//no more work to do
					inputs = nil
				} else {
					item.index = in
					item.req = items[in]
				}
			case resp := <-work.done:
				//fmt.Println(out, resp.index)
				cache[resp.index] = resp.result
				//fmt.Println(resp.index, out)
				for checked[resp.index] = true; checked[out]; out++ {
					//fmt.Println("@@", out, resp)
					result <- cache[out]
					if out == length-1 {
						//finish, end
						return
					}
				}
				//fmt.Println("$$$", out)

			case <-work.abort:
				return
			}
		}
	}()

	return result, nil
}

func (work *MapWork) AbortWorks() {
	work.abort <- struct{}{}
}

func (work *MapWork) StopWorks() {
	work.stop <- struct{}{}
}
