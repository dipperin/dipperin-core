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
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/consts"
	"github.com/dipperin/dipperin-core/common/hexutil"
	"github.com/dipperin/dipperin-core/core/contract"
	"github.com/dipperin/dipperin-core/core/rpc-interface"
	"github.com/urfave/cli"
	"go.uber.org/zap"
	"math/big"
	"strconv"
)

func buildERC20Token(owner common.Address, tokenName string, tokenSymbol string, supply *big.Int, decimal int) *contract.BuiltInERC20Token {

	//erc20 token
	erc20Token := contract.BuiltInERC20Token{}
	erc20Token.Owner = owner
	erc20Token.TokenName = tokenName
	erc20Token.TokenDecimals = decimal
	erc20Token.TokenSymbol = tokenSymbol
	erc20Token.TokenTotalSupply = supply

	return &erc20Token
}

/*func BuildContractExtraData(op string, contractAdr common.Address, params string) []byte {
	erc20 := contract.ExtraDataForContract{
		ContractAddress: contractAdr,
		Action: op,
		Params: params,
	}
	erc20Str, _ := json.Marshal(erc20)
	return erc20Str
}*/

func isParamValid(params []string, num int) bool {
	if len(params) != num {
		return false
	}
	for _, p := range params {
		if p == "" {
			return false
		}
	}
	return true
}

func (caller *rpcCaller) AnnounceERC20(c *cli.Context) {
	_, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error", zap.Error(err))
		return
	}

	if !isParamValid(cParams, 7) {
		l.Error("parameters need：owner_address,token_name,token_symbol,token_total_supply,decimal,gasPrice,gasLimit")
		return
	}

	owner, err := CheckAndChangeHexToAddress(cParams[0])
	if err != nil {
		l.Error("the input address is invalid", zap.Error(err))
		return
	}

	tokenName := cParams[1]
	tokenSymbol := cParams[2]

	decimal, err := strconv.Atoi(cParams[4])
	if err != nil {
		l.Error("the parameter decimal invalid", zap.Error(err))
		return
	}

	if decimal > 18 {
		l.Error("the parameter decimal should be less than 17")
		return
	}

	tokenTotalSupply, err := DecimalToInter(cParams[3], decimal)
	if err != nil {
		l.Error("the parameter value invalid", zap.Error(err))
		return
	}

	gasPrice, err := MoneyValueToCSCoin(cParams[5])
	if err != nil {
		l.Error("the parameter gasPrice invalid", zap.Error(err))
		return
	}

	gasLimit, err := strconv.ParseUint(cParams[6], 10, 64)
	if err != nil {
		l.Error("the parameter gasLimit invalid", zap.Error(err))
		return
	}

	var resp rpc_interface.ERC20Resp
	if err = client.Call(&resp, getDipperinRpcMethodByName("CreateERC20"), owner, tokenName, tokenSymbol, tokenTotalSupply, decimal, gasPrice, gasLimit); err != nil {
		l.Error("AnnounceERC20 failed", zap.Error(err))
		return
	}
	l.Info("SendTransaction result", zap.String("txId", resp.TxId.Hex()))
	l.Info("MUST record", zap.String("contract NO: ", resp.CtId.Hex()))
}

func (caller *rpcCaller) ERC20TotalSupply(c *cli.Context) {
	_, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error", zap.Error(err))
		return
	}

	if !isParamValid(cParams, 1) {
		l.Error("parameters need：contract address")
		return
	}

	contractAdr, err := CheckAndChangeHexToAddress(cParams[0])
	if err != nil {
		l.Error("the input address is invalid", zap.Error(err))
		return
	}

	var resp *big.Int
	if err = client.Call(&resp, getDipperinRpcMethodByName("ERC20TotalSupply"), contractAdr); err != nil {
		l.Error("call ERC20TotalSupply", zap.Error(err))
		return
	}

	decimal := getERC20Decimal(contractAdr)
	ts, _ := InterToDecimal((*hexutil.Big)(resp), decimal)
	unit := getERC20Symbol(contractAdr)
	l.Info("contract info", zap.String("total supply", ts+unit))
}

func (caller *rpcCaller) ERC20Transfer(c *cli.Context) {
	_, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error", zap.Error(err))
		return
	}

	if !isParamValid(cParams, 6) {
		l.Error("parameters need：contract_address,owner,to_address,amount,gasPrice,gasLimit")
		return
	}

	contractAdr, err := CheckAndChangeHexToAddress(cParams[0])
	if err != nil {
		l.Error("the input address is invalid", zap.Error(err))
		return
	}

	owner, err := CheckAndChangeHexToAddress(cParams[1])
	if err != nil {
		l.Error("the input address is invalid", zap.Error(err))
		return
	}

	toAdr, err := CheckAndChangeHexToAddress(cParams[2])
	if err != nil {
		l.Error("the input address is invalid", zap.Error(err))
		return
	}

	//get decimal
	decimal := getERC20Decimal(contractAdr)
	value, err := DecimalToInter(cParams[3], decimal)
	if err != nil {
		l.Error("the parameter value invalid", zap.Error(err))
		return
	}

	gasPrice, err := MoneyValueToCSCoin(cParams[4])
	if err != nil {
		l.Error("the parameter gasPrice invalid", zap.Error(err))
		return
	}

	gasLimit, err := strconv.ParseUint(cParams[5], 10, 64)
	if err != nil {
		l.Error("the parameter gasLimit invalid", zap.Error(err))
		return
	}

	//send transaction
	var resp common.Hash
	if err = client.Call(&resp, getDipperinRpcMethodByName("ERC20Transfer"), contractAdr, owner, toAdr, value, gasPrice, gasLimit); err != nil {
		l.Error("ERC20Transfer failed", zap.Error(err))
		return
	}
	l.Info("ERC20Transfer result", zap.String("txId", resp.Hex()))
}

func (caller *rpcCaller) ERC20TransferFrom(c *cli.Context) {
	_, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error", zap.Error(err))
		return
	}

	if !isParamValid(cParams, 7) {
		l.Error("parameters need：contract_address,owner,from_address,to_address,amount,gasPrice,gasLimit")
		return
	}

	contractAdr, err := CheckAndChangeHexToAddress(cParams[0])
	if err != nil {
		l.Error("the input address is invalid", zap.Error(err))
		return
	}

	owner, err := CheckAndChangeHexToAddress(cParams[1])
	if err != nil {
		l.Error("the input address is invalid", zap.Error(err))
		return
	}

	fromAdr, err := CheckAndChangeHexToAddress(cParams[2])
	if err != nil {
		l.Error("the input address is invalid", zap.Error(err))
		return
	}

	toAdr, err := CheckAndChangeHexToAddress(cParams[3])
	if err != nil {
		l.Error("the input address is invalid", zap.Error(err))
		return
	}

	//get decimal
	decimal := getERC20Decimal(contractAdr)
	value, err := DecimalToInter(cParams[4], decimal)
	if err != nil {
		l.Error("the parameter value invalid", zap.Error(err))
		return
	}

	gasPrice, err := MoneyValueToCSCoin(cParams[5])
	if err != nil {
		l.Error("the parameter gasPrice invalid", zap.Error(err))
		return
	}

	gasLimit, err := strconv.ParseUint(cParams[6], 10, 64)
	if err != nil {
		l.Error("the parameter gasLimit invalid", zap.Error(err))
		return
	}

	//send transaction
	var resp common.Hash
	if err = client.Call(&resp, getDipperinRpcMethodByName("ERC20TransferFrom"), contractAdr, owner, fromAdr, toAdr, value, gasPrice, gasLimit); err != nil {
		l.Error("ERC20TransferFrom failed", zap.Error(err))
		return
	}
	l.Info("ERC20TransferFrom result", zap.String("txId", resp.Hex()))
}

func (caller *rpcCaller) ERC20TokenName(c *cli.Context) {
	_, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error", zap.Error(err))
		return
	}

	if !isParamValid(cParams, 1) {
		l.Error("parameters need：contract address")
		return
	}

	contractAdr, err := CheckAndChangeHexToAddress(cParams[0])
	if err != nil {
		l.Error("the input address is invalid", zap.Error(err))
		return
	}

	extraData := contract.ExtraDataForContract{contractAdr, "Name", "[]"}

	//send transaction
	var resp string
	if err = client.Call(&resp, getDipperinRpcMethodByName("GetContractInfo"), &extraData); err != nil {
		l.Error("call GetContractInfo", zap.Error(err))
		return
	}
	l.Info("contract info", zap.String("token name", resp))
}

func (caller *rpcCaller) ERC20TokenSymbol(c *cli.Context) {
	_, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error", zap.Error(err))
		return
	}

	if !isParamValid(cParams, 1) {
		l.Error("parameters need：contract address")
		return
	}

	contractAdr, err := CheckAndChangeHexToAddress(cParams[0])
	if err != nil {
		l.Error("the input address is invalid", zap.Error(err))
		return
	}

	extraData := contract.ExtraDataForContract{contractAdr, "Symbol", "[]"}

	//send transaction
	var resp string
	if err = client.Call(&resp, getDipperinRpcMethodByName("GetContractInfo"), &extraData); err != nil {
		l.Error("call GetContractInfo", zap.Error(err))
		return
	}
	l.Info("contract info", zap.String("token symbol", resp))
}

func getERC20Symbol(contractAdr common.Address) string {
	extraData := contract.ExtraDataForContract{contractAdr, "Symbol", "[]"}

	//send transaction
	var resp string
	if err := client.Call(&resp, getDipperinRpcMethodByName("GetContractInfo"), &extraData); err != nil {
		l.Error("call GetContractInfo", zap.Error(err))
		return ""
	}
	return resp
}

func getERC20Decimal(contractAdr common.Address) int {
	extraData := contract.ExtraDataForContract{contractAdr, "Decimals", "[]"}

	//send transaction
	var resp int
	if err := client.Call(&resp, getDipperinRpcMethodByName("GetContractInfo"), &extraData); err != nil {
		l.Error("call GetContractInfo", zap.Error(err))
		return consts.DIPDecimalBits
	}

	return resp
}

func (caller *rpcCaller) ERC20TokenDecimals(c *cli.Context) {
	_, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error", zap.Error(err))
		return
	}

	if !isParamValid(cParams, 1) {
		l.Error("parameters need：contract address")
		return
	}

	contractAdr, err := CheckAndChangeHexToAddress(cParams[0])
	if err != nil {
		l.Error("the input address is invalid", zap.Error(err))
		return
	}

	extraData := contract.ExtraDataForContract{contractAdr, "Decimals", "[]"}

	//send transaction
	var resp int
	if err = client.Call(&resp, getDipperinRpcMethodByName("GetContractInfo"), &extraData); err != nil {
		l.Error("call GetContractInfo", zap.Error(err))
		return
	}
	l.Info("contract info", zap.Int("token decimals", resp))
}

func convert(raw []byte) int {
	l := len(raw)
	sum := 0
	for i := 0; i < l; i++ {
		sum = sum << 8
		sum += int(uint8(raw[i]))
	}
	return sum
}

func (caller *rpcCaller) ERC20GetInfo(c *cli.Context) {
	_, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error", zap.Error(err))
		return
	}

	if !isParamValid(cParams, 1) {
		l.Error("parameters need：contract_address")
		return
	}

	contractAdr, err := CheckAndChangeHexToAddress(cParams[0])
	if err != nil {
		l.Error("the input address is invalid", zap.Error(err))
		return
	}

	var resp interface{}

	if err = client.Call(&resp, getDipperinRpcMethodByName("GetContract"), contractAdr); err != nil {
		l.Error("call GetContract", zap.Error(err))
		return
	}

	decimal := getERC20Decimal(contractAdr)
	unit := getERC20Symbol(contractAdr)
	if ct, ok := resp.(map[string]interface{}); ok {
		ts := common.FromHex(ct["token_total_supply"].(string))
		num := convert(ts)
		numBig := big.NewInt(int64(num))
		tokenNum, _ := InterToDecimal((*hexutil.Big)(numBig), decimal)
		l.Info("contract:", zap.Any("ct", ct), zap.String("total supply", tokenNum+unit))
		return
	}
	l.Error("call ERC20GetInfo fail", zap.Any("err", resp))
}

func (caller *rpcCaller) ERC20Allowance(c *cli.Context) {
	_, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error", zap.Error(err))
		return
	}

	if !isParamValid(cParams, 3) {
		l.Error("parameters need：contract_address,owner,spender")
		return
	}

	contractAdr, err := CheckAndChangeHexToAddress(cParams[0])
	if err != nil {
		l.Error("the input address is invalid", zap.Error(err))
		return
	}

	owner, err := CheckAndChangeHexToAddress(cParams[1])
	if err != nil {
		l.Error("the input address is invalid", zap.Error(err))
		return
	}

	spender, err := CheckAndChangeHexToAddress(cParams[2])
	if err != nil {
		l.Error("the input address is invalid", zap.Error(err))
		return
	}

	var resp *big.Int
	if err = client.Call(&resp, getDipperinRpcMethodByName("ERC20Allowance"), contractAdr, owner, spender); err != nil {
		l.Error("call ERC20Allowance", zap.Error(err))
		return
	}

	decimal := getERC20Decimal(contractAdr)
	ts, _ := InterToDecimal((*hexutil.Big)(resp), decimal)
	unit := getERC20Symbol(contractAdr)
	l.Info("contract info", zap.String("token_allowance", ts+unit))
}

func (caller *rpcCaller) ERC20Approve(c *cli.Context) {
	_, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error", zap.Error(err))
		return
	}

	if !isParamValid(cParams, 6) {
		l.Error("parameters need：contract_address,owner,to_address,amount,gasPrice,gasLimit")
		return
	}

	contractAdr, err := CheckAndChangeHexToAddress(cParams[0])
	if err != nil {
		l.Error("the input address is invalid", zap.Error(err))
		return
	}

	from, err := CheckAndChangeHexToAddress(cParams[1])
	if err != nil {
		l.Error("the input address is invalid", zap.Error(err))
		return
	}

	to, err := CheckAndChangeHexToAddress(cParams[2])
	if err != nil {
		l.Error("the input address is invalid", zap.Error(err))
		return
	}

	//get decimal
	decimal := getERC20Decimal(contractAdr)
	value, err := DecimalToInter(cParams[3], decimal)
	if err != nil {
		l.Error("the parameter value invalid", zap.Error(err))
		return
	}

	//get totalsupply for validating approve
	var total *big.Int
	if err := client.Call(&total, getDipperinRpcMethodByName("ERC20TotalSupply"), contractAdr); err != nil {
		l.Error("call ERC20TotalSupply", zap.Error(err))
		return
	}

	if total.Cmp(value) == -1 {
		ts, _ := InterToDecimal((*hexutil.Big)(total), decimal)
		l.Error("approving credit exceeding", zap.String("total", ts))
		return
	}

	gasPrice, err := MoneyValueToCSCoin(cParams[4])
	if err != nil {
		l.Error("the parameter gasPrice invalid", zap.Error(err))
		return
	}

	gasLimit, err := strconv.ParseUint(cParams[5], 10, 64)
	if err != nil {
		l.Error("the parameter gasLimit invalid", zap.Error(err))
		return
	}

	//send transaction
	var resp common.Hash
	if err := client.Call(&resp, getDipperinRpcMethodByName("ERC20Approve"), contractAdr, from, to, value, gasPrice, gasLimit); err != nil {
		l.Error("ERC20Approve failed", zap.Error(err))
		return
	}
	l.Info("ERC20Approve result", zap.String("txId", resp.Hex()))
}

func (caller *rpcCaller) ERC20Balance(c *cli.Context) {
	_, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error", zap.Error(err))
		return
	}

	if !isParamValid(cParams, 2) {
		l.Error("parameters need：contract_address,owner_address")
		return
	}

	contractAdr, err := CheckAndChangeHexToAddress(cParams[0])
	if err != nil {
		l.Error("the input address is invalid", zap.Error(err))
		return
	}

	owner, err := CheckAndChangeHexToAddress(cParams[1])
	if err != nil {
		l.Error("the input address is invalid", zap.Error(err))
		return
	}

	//send transaction
	var resp *hexutil.Big
	if err = client.Call(&resp, getDipperinRpcMethodByName("ERC20Balance"), contractAdr, owner); err != nil {
		l.Error("call ERC20Balance", zap.Error(err))
		return
	}

	decimal := getERC20Decimal(contractAdr)
	ts, _ := InterToDecimal(resp, decimal)
	unit := getERC20Symbol(contractAdr)

	l.Info("contract info", zap.Any("address", owner), zap.String("token balance", ts+unit))
}
