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
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/core/contract"
	"github.com/dipperin/dipperin-core/core/rpc-interface"
	"github.com/urfave/cli"
	"go.uber.org/zap"
	"strconv"
)

func (caller *rpcCaller) TransferEDIPToDIP(c *cli.Context) {
	_, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error", zap.Error(err))
		return
	}

	if len(cParams) != 4 {
		l.Error("EarlyTokenTransferEDIPToDIP needs at least：from,eDIPValue,gasPrice,gasLimit")
		return
	}

	from, err := CheckAndChangeHexToAddress(cParams[0])
	if err != nil {
		l.Error("the from is invalid", zap.Error(err))
		return
	}

	eDIPValue, err := DecimalToInter(cParams[1], contract.DecimalUnits)
	if err != nil {
		l.Error("the eDIPValue is invalid", zap.Error(err))
		return
	}

	gasPrice, err := MoneyValueToCSCoin(cParams[2])
	if err != nil {
		l.Error("the parameter gasPrice invalid", zap.Error(err))
		return
	}

	gasLimit, err := strconv.ParseUint(cParams[3], 10, 64)
	if err != nil {
		l.Error("the parameter gasLimit invalid", zap.Error(err))
		return
	}

	vStr := fmt.Sprintf("0x%x", eDIPValue)
	params := util.StringifyJson([]interface{}{from.Hex(), vStr})
	contractAdr := contract.EarlyContractAddress
	extraData := rpc_interface.BuildContractExtraData("TransferEDIPToDIP", contractAdr, params)

	//send transaction
	var resp common.Hash
	if err := client.Call(&resp, getDipperinRpcMethodByName("SendTransaction"), from, contractAdr, 0, gasPrice, gasLimit, extraData, nil); err != nil {
		l.Error("Call a send transaction", zap.Error(err))
		return
	}
	l.Info("SendTransaction result", zap.String("txId", resp.Hex()))
}

func (caller *rpcCaller) SetExchangeRate(c *cli.Context) {
	_, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error", zap.Error(err))
		return
	}

	if len(cParams) != 4 {
		l.Error("SetExchangeRate needs at least：from,exchangeRate,gasPrice,gasLimit")
		return
	}

	from, err := CheckAndChangeHexToAddress(cParams[0])
	if err != nil {
		l.Error("the from is invalid", zap.Error(err))
		return
	}

	value, err := strconv.Atoi(cParams[1])
	if err != nil {
		l.Error("the parameter exChangeRate invalid", zap.Error(err))
	}

	gasPrice, err := MoneyValueToCSCoin(cParams[2])
	if err != nil {
		l.Error("the parameter gasPrice invalid", zap.Error(err))
		return
	}

	gasLimit, err := strconv.ParseUint(cParams[3], 10, 64)
	if err != nil {
		l.Error("the parameter gasLimit invalid", zap.Error(err))
		return
	}

	//svStr := fmt.Sprintf("0x%x", value)
	params := util.StringifyJson([]interface{}{from.Hex(), int64(value)})
	contractAdr := contract.EarlyContractAddress
	extraData := rpc_interface.BuildContractExtraData("SetExchangeRate", contractAdr, params)

	//send transaction
	var resp common.Hash
	if err := client.Call(&resp, getDipperinRpcMethodByName("SendTransaction"), from, contractAdr, 0, gasPrice, gasLimit, extraData, nil); err != nil {
		l.Error("call sending transaction", zap.Error(err))
		return
	}
	l.Info("SendTransaction result", zap.String("txId", resp.Hex()))
}
