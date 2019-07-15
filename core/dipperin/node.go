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

package dipperin

import (
	"github.com/dipperin/dipperin-core/third-party/log"
	"reflect"
	"sync"
	"sync/atomic"
	"time"
)

var (
	chokeTimeout = 20 * time.Second

)

type NodeService interface {
	Start() error
	Stop()
}

func NewCsNode(services []NodeService) *CsNode {
	return &CsNode{
		services: services,
	}
}

type CsNode struct {
	services []NodeService

	wg         sync.WaitGroup
	chokePoint uint32
}

func (n *CsNode) AddService(service NodeService) {
	n.services = append(n.services, service)
}

func (n *CsNode) Start() (err error) {
	startSuccess := true
	// If the startup is completed in 20 seconds, there is a service blocked (the virtual machine is too bad, it cannot be started in a few seconds)
	go func() {
		time.Sleep(chokeTimeout)
		i := atomic.LoadUint32(&n.chokePoint)
		if startSuccess && i < uint32(len(n.services)) {
			panic("node start choked by service:" + reflect.TypeOf(n.services[i]).String())
		}
	}()
	for _, s := range n.services {
		n.wg.Add(1)

		// if start err, stop all services and return the err
		if err = s.Start(); err != nil {
			log.Info("the err service is: ", "service", s)
			log.Error("start node service failed", "err", err)
			startSuccess = false
			n.Stop()
			break
		}
		atomic.AddUint32(&n.chokePoint, 1)
	}
	return
}

func (n *CsNode) Stop() {
	log.Info("call node stop")
	x := 0
	t := time.NewTimer(chokeTimeout)
	// check service stop choked
	go func() {
		select {
		case <-t.C:
			panic("node stop choked by service:" + reflect.TypeOf(n.services[x]).String())
		}
	}()
	i := int(atomic.LoadUint32(&n.chokePoint))
	for ; x < i; x++ {
		// stop shouldn't in go
		n.services[x].Stop()
		t.Reset(chokeTimeout)
		n.wg.Done()
	}
	t.Stop()
}

func (n *CsNode) Wait() {
	n.wg.Wait()
}
