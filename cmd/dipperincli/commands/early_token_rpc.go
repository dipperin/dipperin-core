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


package commands

import (
	"github.com/urfave/cli"
	"fmt"
	"github.com/dipperin/dipperin-core/core/contract"
	"github.com/dipperin/dipperin-core/common"
	"strconv"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/core/rpc-interface"
)

func (caller *rpcCaller) TransferEDIPToDIP(c *cli.Context) {
	_, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error", "err", err)
		return
	}

	if len(cParams) != 3 {
		l.Error("EarlyTokenTransferEDIPToDIP needs at least：from, eDIPValue,transactionFee")
		return
	}

	from, err := CheckAndChangeHexToAddress(cParams[0])
	if err != nil {
		l.Error("the from is invalid", "err", err)
		return
	}

	eDIPValue,err := DecimalToInter(cParams[1],contract.DecimalUnits)
	if err !=nil{
		l.Error("the eDIPValue is invalid", "err", err)
		return
	}

	txFee, err := MoneyValueToCSCoin(cParams[2])
	if err != nil {
		l.Error("the parameter transactionFee invalid", "err", err)
		return
	}

	vStr := fmt.Sprintf("0x%x", eDIPValue)
	params := util.StringifyJson([]interface{}{ from.Hex(), vStr })
	contractAdr:= contract.EarlyContractAddress
	extraData := rpc_interface.BuildContractExtraData("TransferEDIPToDIP",contractAdr,params)

	//send transaction
	var resp common.Hash
	if err := client.Call(&resp, getDipperinRpcMethodByName("SendTransaction"),from, contractAdr,0, txFee, extraData, nil); err != nil {
		l.Error("Call a send transaction", "err", err)
		return
	}
	l.Info("SendTransaction result", "txId", resp.Hex())
}

func (caller *rpcCaller) SetExchangeRate(c *cli.Context) {
	_, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error", "err", err)
		return
	}

	if len(cParams) != 3 {
		l.Error("SetExchangeRate needs at least：from, exchangeRate,transactionFee")
		return
	}

	from, err := CheckAndChangeHexToAddress( cParams[0])
	if err != nil {
		l.Error("the from is invalid", "err", err)
		return
	}

	exChangeRate := cParams[1]

	txFee, err := MoneyValueToCSCoin(cParams[2])
	if err != nil {
		l.Error("the parameter transactionFee invalid", "err", err)
		return
	}

	value ,err:= strconv.Atoi(exChangeRate)
	if err !=nil{
		l.Error("the parameter exChangeRate invalid","err",err)
	}

	//svStr := fmt.Sprintf("0x%x", value)
	params := util.StringifyJson([]interface{}{ from.Hex(), int64(value)})
	contractAdr:= contract.EarlyContractAddress
	extraData := rpc_interface.BuildContractExtraData("SetExchangeRate",contractAdr,params)

	//send transaction
	var resp common.Hash
	if err := client.Call(&resp, getDipperinRpcMethodByName("SendTransaction"),from, contractAdr,0, txFee, extraData, nil); err != nil {
		l.Error("call sending transaction", "err", err)
		return
	}
	l.Info("SendTransaction result", "txId", resp.Hex())
}


