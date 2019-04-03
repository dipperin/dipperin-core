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


package call_node_rpc

import (
	"context"
	"github.com/dipperin/dipperin-core/third-party/rpc"
	"testing"
	"fmt"
	"time"
	"github.com/dipperin/dipperin-core/third-party/log"
)

func TestWs(t *testing.T) {
	return
	log.InitLogger(log.LvlDebug)
	_, err := rpc.Dial(fmt.Sprintf("ws://%v:%v", "localhost", 7002))
	if err != nil {
		panic(err.Error())
	}
	time.Sleep(1 * time.Second)
}

type HaHa struct {
	ID uint
}

// CsSubscribe registers a subscripion under the "cs" namespace.
func csSubscribe(c *rpc.Client, ctx context.Context, channel interface{}, args ...interface{}) (*rpc.ClientSubscription, error) {
	return c.Subscribe(ctx, "dipperin", channel, args...)
}

func TestSubscribe(t *testing.T){
	return
	client, err := rpc.Dial("ws://localhost:10002")
	defer client.Close()
	//client, err := rpc.Dial("ws://10.200.0.139:10002")
	log.Info("the err is:","err",err)
	if err != nil {
		panic(err)
	}

	xx := make(chan HaHa)
	sub, err := csSubscribe(client, context.Background(), xx, "subscribeBlock")
	log.Info("the err is:","err",err)
	if err != nil {
		panic(err)
	}
	for {
		select {
		case err := <-sub.Err():
			log.Info("sub result", "err", err)
		case h := <-xx:
			// todo makes another prompt, otherwise issuing commands is difficult to manipulate
			log.Info("sdfnoiwef", "h", h)
		}
	}
}
