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
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/accounts"
	"github.com/dipperin/dipperin-core/core/rpc-interface"
	"github.com/dipperin/dipperin-core/third-party/rpc"
	"github.com/dipperin/dipperin-core/third-party/log"
)

func main() {
	//c, err := rpc.Dial(fmt.Sprintf("ws://%v:%v", "localhost", 8005))
	//c, err := rpc.Dial(fmt.Sprintf("http://%v:%v", "10.200.0.139", 3035))
	c, err := rpc.Dial(fmt.Sprintf("http://%v:%v", "szly", 3035))
	if err != nil {
		panic(err)
	}

	curAddr := getDefaultAccount(c)

	var resp rpc_interface.CurBalanceResp
	if err := c.Call(&resp, "dipperin_currentBalance", curAddr); err != nil {
		log.Error("dipperin_call currentBalance failed", "err", err)
		return
	}
	log.Info("current balance", "balance", curAddr.Hex(), resp.Balance.ToInt().String())
}


func getDefaultAccount(c *rpc.Client) common.Address {
	var resp []accounts.WalletIdentifier

	if err := c.Call(&resp, "dipperin_listWallet"); err != nil {
		fmt.Println("Call ListWallet", "err", err)
		return common.Address{}
	}

	if len(resp) == 0 {
		return common.Address{}
	}

	var respA []accounts.Account

	if err := c.Call(&respA, "dipperin_listWalletAccount", resp[0]); err != nil {
		fmt.Println("Call ListWallet", "err", err)
		return common.Address{}
	}
	return respA[0].Address
}