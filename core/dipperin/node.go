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
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

var (
	chokeTimeout = 20 * time.Second
)

type ServiceType int

const (
	NeedWalletSignerService ServiceType = iota
	NotNeedWalletSignerService
)

type NodeService interface {
	Start() error
	Stop()
}

func NewCsNode(conf NodeConfig, components *BaseComponent) *CsNode {
	return &CsNode{
		nodeName:       conf.Name,
		ServiceManager: NewServiceManager(components),
	}
}

type CsNode struct {
	*ServiceManager
	nodeName string
}

type ServiceManager struct {
	components       *BaseComponent
	services         map[ServiceType][]NodeService

	chokePoints       sync.Map
	serviceNumber    uint32
	serviceStartFlag atomic.Value
	wg               sync.WaitGroup
}

func NewServiceManager(components *BaseComponent) *ServiceManager {
	s := make(map[ServiceType][]NodeService)
	s[NeedWalletSignerService] = make([]NodeService, 0)
	s[NotNeedWalletSignerService] = make([]NodeService, 0)
	manager := &ServiceManager{services: s, components: components}

	number := uint32(0)
	for _, service := range components.getNodeServices() {
		manager.AddService(service)
		number++
	}
	manager.serviceNumber = number
	return manager
}

func (m *ServiceManager) startService(t ServiceType) error{
	startSuccess := true
	// If the startup is completed in 20 seconds, there is a service blocked (the virtual machine is too bad, it cannot be started in a few seconds)
	go func() {
		time.Sleep(chokeTimeout)
		i,ok := m.chokePoints.Load(t)
		if !ok{
			panic("start service get chokePoints error")
		}
		if startSuccess && i.(uint32) < uint32(len(m.services[t])) {
			panic("node start choked by service:" + reflect.TypeOf(m.services[t][i.(uint32)]).String())
		}
	}()

	//print debug stack info
	//go printStackInfo(n.nodeName)
	var number uint32
	var err error
	for _, s := range m.services[t] {
		m.wg.Add(1)
		// if start err, stop all services and return the err
		log.Info("start service ","name",reflect.TypeOf(s).String())
		if err = s.Start(); err != nil {
			log.Info("the err service is: ", "service", s)
			log.Error("start node service failed", "err", err)
			startSuccess = false
			break
		}
		number++
	}
	m.chokePoints.Store(t,number)
	if err !=nil{
		m.Stop()
		return err
	}

	return nil
}

func (m *ServiceManager) stopService(service ServiceType){
	log.Info("call node stop")
	x := 0
	t := time.NewTimer(chokeTimeout)
	// check service stop choked
	go func() {
		select {
		case <-t.C:
			panic("node stop choked by service:" + reflect.TypeOf(m.services[service]).String())
		}
	}()

	i,ok := m.chokePoints.Load(service)
	if !ok{
		panic("stop service get chokePoints error")
	}
	for ; uint32(x) < i.(uint32); x++ {
		// stop shouldn't in go
		log.Info("stop service~~~","x",x,"service",reflect.TypeOf(m.services[service][x]).String())
		m.services[service][x].Stop()
		t.Reset(chokeTimeout)
		m.wg.Done()
	}
	t.Stop()
	return
}

func (m *ServiceManager) startRemainingServices() error{
	startFlag := m.serviceStartFlag.Load()
	if startFlag!=nil && startFlag.(bool){
		log.Info("the serviceStartFlag is true")
		return nil
	}
	err := m.components.setNodeSignerInfo()
	if err != nil {
		panic("serviceManager startRemainingServices err:" + err.Error())
	}
	err = m.startService(NeedWalletSignerService)
	if err !=nil{
		return err
	}
	m.serviceStartFlag.Store(true)
	return nil
}

func (m *ServiceManager) Start() error {
	//start Services not depend on the wallet signer
	err :=m.startService(NotNeedWalletSignerService)
	if err !=nil{
		return err
	}
	log.Info("the NoWalletStart is:","NoWalletStart",m.components.nodeConfig.NoWalletStart)
	if !m.components.nodeConfig.NoWalletStart {
		err =m.startRemainingServices()
	} else {
		// listen wallet manager. start other service when there is an soft wallet
		go func() {
			walletEvent := m.components.walletManager.SubScribeStartService()
			for {
				select {
				case <-walletEvent:
					//set wallet signer and start the other services
					m.startRemainingServices()
					return
				}
			}
		}()
	}

	return err
}

func (m *ServiceManager) Stop() {
	startFlag := m.serviceStartFlag.Load()
	if startFlag!=nil && startFlag.(bool){
		m.stopService(NeedWalletSignerService)
	}

	m.stopService(NotNeedWalletSignerService)
	m.serviceStartFlag.Store(false)
}

func (m *ServiceManager) Wait() {
	m.wg.Wait()
}

func (m *ServiceManager) AddService(service NodeService) {
	//m.services = append(m.services, service)
	serviceType := reflect.TypeOf(service).String()
	log.Info("the service type is:", "name", serviceType)

	switch serviceType {
	case "*service.VenusFullChainService", "*accounts.WalletManager",
		"*rpc_interface.Service", "*tx_pool.TxPool", "*g_metrics.PrometheusMetricsServer",
		"*vm_log_search.ChainIndexer":
		m.services[NotNeedWalletSignerService] = append(m.services[NotNeedWalletSignerService], service)
	default:
		m.services[NeedWalletSignerService] = append(m.services[NeedWalletSignerService], service)
	}
	return
}

/*func (n *CsNode) AddService(service NodeService) {
	n.services = append(n.services, service)
}*/
func logDebugStack() {
	buf := make([]byte, 5*1024*1024)
	log.Stack.Info("the runtime stack is:~~~~~~~~~~~~~~~~~~~~")
	buf = buf[:runtime.Stack(buf, true)]
	log.Stack.Info(string(buf))
	log.Stack.Info("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")

}

func printStackInfo(nodeName string) {
	tick := time.NewTicker(2 * time.Minute)
	for {
		select {
		case <-tick.C:
			logDebugStack()
		}
	}
}

/*func (n *CsNode) Start() (err error) {
	startSuccess := true
	// If the startup is completed in 20 seconds, there is a service blocked (the virtual machine is too bad, it cannot be started in a few seconds)
	go func() {
		time.Sleep(chokeTimeout)
		i := atomic.LoadUint32(&n.chokePoint)
		if startSuccess && i < uint32(len(n.services)) {
			panic("node start choked by service:" + reflect.TypeOf(n.services[i]).String())
		}
	}()

	//print debug stack info
	//go printStackInfo(n.nodeName)

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
*/
