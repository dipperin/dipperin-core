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


package main

import (
	"context"
	"github.com/dipperin/dipperin-core/third-party/rpc"
	"time"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/dipperin/dipperin-core/core/rpc-interface"
	"github.com/dipperin/dipperin-core/core/chain-config"
)

type HaHa struct {
	ID uint
}

type TestApi struct {}

func (api *TestApi) SubscribeBlock(ctx context.Context) (*rpc.Subscription, error) {
	notifier, supported := rpc.NotifierFromContext(ctx)
	if !supported {
		return &rpc.Subscription{}, rpc.ErrNotificationsUnsupported
	}

	rpcSub := notifier.CreateSubscription()

	go func() {
		//blockCh := make(chan model.Block)
		//blockSub := service.nodeContext.ChainReader().SubscribeBlockEvent(blockCh)

		for {
			select {
			case <-time.NewTicker(2 * time.Second).C:
				log.Info("notify client")
				notifier.Notify(rpcSub.ID, HaHa{ ID: 123 })
				//if rV, err := rlp.EncodeToBytes(b); err != nil {
				//	log.Error("block can't encode to bytes", "err", err)
				//} else {
				//	if err := notifier.Notify(rpcSub.ID, rV); err != nil {
				//		log.Error("can't notify wallet", "err", err)
				//	}
				//}

			case <-rpcSub.Err():
				return
			case <-notifier.Closed():
				return
			}
		}

	}()

	return rpcSub, nil
}

func main() {
	s := rpc_interface.MakeRpcService(&fakeNodeConf{}, []rpc.API{
		{
			Namespace: "dipperin",
			Version:   chain_config.Version,
			Service:   &TestApi{},
			Public:    true,
		},
	}, []string{"*"})

	if err := s.Start(); err != nil {
		panic(err)
	}

	select {}
}

type fakeNodeConf struct {

}

func (c *fakeNodeConf)IpcEndpoint() string {
	return "/tmp/rpc_sub_test"
}

func (c *fakeNodeConf)HttpEndpoint() string {
	return ":10001"
}

func (c *fakeNodeConf)WsEndpoint() string {
	return ":10002"
}