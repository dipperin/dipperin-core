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
		l.Error("getRpcMethodAndParam error", "err", err)
		return
	}

	if !isParamValid(cParams, 7) {
		l.Error("parameters need：owner_address, token_name, token_symbol, token_total_supply, decimal, gasPrice,gasLimit")
		return
	}

	owner, err := CheckAndChangeHexToAddress(cParams[0])
	if err != nil {
		l.Error("the input address is invalid", "err", err)
		return
	}

	tokenName := cParams[1]
	tokenSymbol := cParams[2]

	decimal, err := strconv.Atoi(cParams[4])
	if err != nil {
		l.Error("the parameter decimal invalid", "err", err)
		return
	}

	if decimal > 18 {
		l.Error("the parameter decimal should be less than 17")
		return
	}

	tokenTotalSupply, err := DecimalToInter(cParams[3], decimal)
	if err != nil {
		l.Error("the parameter value invalid", "err", err)
		return
	}

	gasPrice, err := MoneyValueToCSCoin(cParams[5])
	if err != nil {
		l.Error("the parameter gasPrice invalid", "err", err)
		return
	}

	gasLimit, err := strconv.Atoi(cParams[6])
	if err != nil {
		l.Error("the parameter gaLimit invalid", "err", err)
		return
	}

	var resp rpc_interface.ERC20Resp
	//create contract, dest address must equal to contract address in transaction

	if err := client.Call(&resp, getDipperinRpcMethodByName("CreateERC20"), owner, tokenName, tokenSymbol, tokenTotalSupply, decimal, gasPrice, gasLimit); err != nil {
		l.Error("AnnounceERC20 failed", "err", err)
		return
	}
	l.Info("SendTransaction result", "txId", resp.TxId.Hex())
	l.Info("MUST record", "contract NO: ", resp.CtId.Hex())
}

func (caller *rpcCaller) ERC20TotalSupply(c *cli.Context) {
	_, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error", "err", err)
		return
	}

	if !isParamValid(cParams, 1) {
		l.Error("parameters need：contract address")
		return
	}

	contractAdr, err := CheckAndChangeHexToAddress(cParams[0])
	if err != nil {
		l.Error("the input address is invalid", "err", err)
		return
	}

	var resp *big.Int
	if err := client.Call(&resp, getDipperinRpcMethodByName("ERC20TotalSupply"), contractAdr); err != nil {
		l.Error("call GetContractInfo", "err", err)
		return
	}

	decimal := getERC20Decimal(contractAdr)
	ts, _ := InterToDecimal((*hexutil.Big)(resp), decimal)
	unit := getERC20Symbol(contractAdr)
	l.Info("contract info", "total supply", ts+unit)
}

func (caller *rpcCaller) ERC20Transfer(c *cli.Context) {
	_, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error", "err", err)
		return
	}

	if !isParamValid(cParams, 6) {
		l.Error("parameters need：contract address, owner, to_address, amount, gasPrice,gasLimit")
		return
	}

	contractAdr, err := CheckAndChangeHexToAddress(cParams[0])
	if err != nil {
		l.Error("the input address is invalid", "err", err)
		return
	}

	owner, err := CheckAndChangeHexToAddress(cParams[1])
	if err != nil {
		l.Error("the input address is invalid", "err", err)
		return
	}

	toAdr, err := CheckAndChangeHexToAddress(cParams[2])
	if err != nil {
		l.Error("the input address is invalid", "err", err)
		return
	}

	//get decimal
	decimal := getERC20Decimal(contractAdr)
	value, err := DecimalToInter(cParams[3], decimal)
	//value, err := MoneyValueToCSCoin(cParams[3])
	if err != nil {
		l.Error("the parameter value invalid", "err", err)
		return
	}

	gasPrice, err := MoneyValueToCSCoin(cParams[4])
	if err != nil {
		l.Error("the parameter gasPrice invalid", "err", err)
		return
	}

	gasLimit, err := strconv.Atoi(cParams[5])
	if err != nil {
		l.Error("the parameter gaLimit invalid", "err", err)
		return
	}

	//send transaction
	var resp common.Hash
	if err := client.Call(&resp, getDipperinRpcMethodByName("ERC20Transfer"), contractAdr, owner, toAdr, value, gasPrice, gasLimit); err != nil {
		l.Error("ERC20Transfer failed", "err", err)
		return
	}
	l.Info("ERC20Transfer result", "txId", resp.Hex())
}

func (caller *rpcCaller) ERC20TransferFrom(c *cli.Context) {
	_, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error", "err", err)
		return
	}

	if !isParamValid(cParams, 7) {
		l.Error("parameters need：contract address, owner, from_address, to_address, amount, gasPrice,gasLimit")
		return
	}

	contractAdr, err := CheckAndChangeHexToAddress(cParams[0])
	if err != nil {
		l.Error("the input address is invalid", "err", err)
		return
	}

	owner, err := CheckAndChangeHexToAddress(cParams[1])
	if err != nil {
		l.Error("the input address is invalid", "err", err)
		return
	}

	fromAdr, err := CheckAndChangeHexToAddress(cParams[2])
	if err != nil {
		l.Error("the input address is invalid", "err", err)
		return
	}

	toAdr, err := CheckAndChangeHexToAddress(cParams[3])
	if err != nil {
		l.Error("the input address is invalid", "err", err)
		return
	}

	//get decimal
	decimal := getERC20Decimal(contractAdr)
	value, err := DecimalToInter(cParams[4], decimal)
	//value, err := MoneyValueToCSCoin(cParams[4])
	if err != nil {
		l.Error("the parameter value invalid", "err", err)
		return
	}

	gasPrice, err := MoneyValueToCSCoin(cParams[5])
	if err != nil {
		l.Error("the parameter gasPrice invalid", "err", err)
		return
	}

	gasLimit, err := strconv.Atoi(cParams[6])
	if err != nil {
		l.Error("the parameter gaLimit invalid", "err", err)
		return
	}

	//send transaction
	var resp common.Hash
	if err := client.Call(&resp, getDipperinRpcMethodByName("ERC20TransferFrom"), contractAdr, owner, fromAdr, toAdr, value, gasPrice, gasLimit); err != nil {
		l.Error("ERC20TransferFrom failed", "err", err)
		return
	}
	l.Info("ERC20TransferFrom result", "txId", resp.Hex())
}

func (caller *rpcCaller) ERC20TokenName(c *cli.Context) {
	_, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error", "err", err)
		return
	}

	if !isParamValid(cParams, 1) {
		l.Error("parameters need：contract address")
		return
	}

	contractAdr, err := CheckAndChangeHexToAddress(cParams[0])
	if err != nil {
		l.Error("the input address is invalid", "err", err)
		return
	}

	extraData := contract.ExtraDataForContract{contractAdr, "Name", "[]"}

	//send transaction
	var resp string
	if err := client.Call(&resp, getDipperinRpcMethodByName("GetContractInfo"), &extraData); err != nil {
		l.Error("call GetContractInfo", "err", err)
		return
	}
	l.Info("contract info", "token name", resp)
}

func (caller *rpcCaller) ERC20TokenSymbol(c *cli.Context) {
	_, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error", "err", err)
		return
	}

	if !isParamValid(cParams, 1) {
		l.Error("parameters need：contract address")
		return
	}

	contractAdr, err := CheckAndChangeHexToAddress(cParams[0])
	if err != nil {
		l.Error("the input address is invalid", "err", err)
		return
	}

	extraData := contract.ExtraDataForContract{contractAdr, "Symbol", "[]"}

	//send transaction
	var resp string
	if err := client.Call(&resp, getDipperinRpcMethodByName("GetContractInfo"), &extraData); err != nil {
		l.Error("call GetContractInfo", "err", err)
		return
	}
	l.Info("contract info", "token symbol", resp)
}

func getERC20Symbol(contractAdr common.Address) string {
	extraData := contract.ExtraDataForContract{contractAdr, "Symbol", "[]"}

	//send transaction
	var resp string
	if err := client.Call(&resp, getDipperinRpcMethodByName("GetContractInfo"), &extraData); err != nil {
		l.Error("call GetContractInfo", "err", err)
		return ""
	}
	return resp
}

func getERC20Decimal(contractAdr common.Address) int {
	extraData := contract.ExtraDataForContract{contractAdr, "Decimals", "[]"}

	//send transaction
	var resp int
	if err := client.Call(&resp, getDipperinRpcMethodByName("GetContractInfo"), &extraData); err != nil {
		l.Error("call GetContractInfo", "err", err)
		return consts.DIPDecimalBits
	}

	return resp
}

func (caller *rpcCaller) ERC20TokenDecimals(c *cli.Context) {
	_, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error", "err", err)
		return
	}

	if !isParamValid(cParams, 1) {
		l.Error("parameters need：contract address")
		return
	}

	contractAdr, err := CheckAndChangeHexToAddress(cParams[0])
	if err != nil {
		l.Error("the input address is invalid", "err", err)
		return
	}

	extraData := contract.ExtraDataForContract{contractAdr, "Decimals", "[]"}

	//send transaction
	var resp int
	if err := client.Call(&resp, getDipperinRpcMethodByName("GetContractInfo"), &extraData); err != nil {
		l.Error("call GetContractInfo", "err", err)
		return
	}
	l.Info("contract info", "token decimals", resp)
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
		l.Error("getRpcMethodAndParam error", "err", err)
		return
	}

	if !isParamValid(cParams, 1) {
		l.Error("parameters need：contract address")
		return
	}

	contractAdr, err := CheckAndChangeHexToAddress(cParams[0])
	if err != nil {
		l.Error("the input address is invalid", "err", err)
		return
	}

	var resp interface{}

	if err := client.Call(&resp, getDipperinRpcMethodByName("GetContract"), contractAdr); err != nil {
		l.Error("call GetContract", "err", err)
		return
	}

	decimal := getERC20Decimal(contractAdr)
	unit := getERC20Symbol(contractAdr)

	if ct, ok := resp.(map[string]interface{}); ok {
		ts := common.FromHex(ct["token_total_supply"].(string))
		num := convert(ts)
		numBig := big.NewInt(int64(num))
		tokenNum, _ := InterToDecimal((*hexutil.Big)(numBig), decimal)
		l.Info("contract:", "owner", ct["owner"], "\nname", ct["token_name"], "\nsymbol", ct["token_symbol"], "\ndecimal", ct["token_decimals"], "\ntotal supply", tokenNum+unit)
		return
	}
	l.Error("call GetContract fail", "err", resp)

}

func (caller *rpcCaller) ERC20Allowance(c *cli.Context) {
	_, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error", "err", err)
		return
	}

	if !isParamValid(cParams, 3) {
		l.Error("parameters need：contract address, owner, spender")
		return
	}

	contractAdr, err := CheckAndChangeHexToAddress(cParams[0])
	if err != nil {
		l.Error("the input address is invalid", "err", err)
		return
	}

	owner, err := CheckAndChangeHexToAddress(cParams[1])
	if err != nil {
		l.Error("the input address is invalid", "err", err)
		return
	}

	spender, err := CheckAndChangeHexToAddress(cParams[2])
	if err != nil {
		l.Error("the input address is invalid", "err", err)
		return
	}

	var resp *big.Int
	if err := client.Call(&resp, getDipperinRpcMethodByName("ERC20Allowance"), contractAdr, owner, spender); err != nil {
		l.Error("call GetContractInfo", "err", err)
		return
	}

	decimal := getERC20Decimal(contractAdr)
	ts, _ := InterToDecimal((*hexutil.Big)(resp), decimal)
	unit := getERC20Symbol(contractAdr)
	l.Info("contract info", "address", spender, "token allowance is", ts+unit)
}

func (caller *rpcCaller) ERC20Approve(c *cli.Context) {
	_, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error", "err", err)
		return
	}

	if !isParamValid(cParams, 6) {
		l.Error("parameters need：contract address, owner, to_address, amount,gasPrice,gasLimit")
		return
	}

	contractAdr, err := CheckAndChangeHexToAddress(cParams[0])
	if err != nil {
		l.Error("the input address is invalid", "err", err)
		return
	}

	from, err := CheckAndChangeHexToAddress(cParams[1])
	if err != nil {
		l.Error("the input address is invalid", "err", err)
		return
	}

	to, err := CheckAndChangeHexToAddress(cParams[2])
	if err != nil {
		l.Error("the input address is invalid", "err", err)
		return
	}

	//get decimal
	decimal := getERC20Decimal(contractAdr)
	value, err := DecimalToInter(cParams[3], decimal)
	//value, err := MoneyValueToCSCoin(cParams[3])
	if err != nil {
		l.Error("the parameter value invalid", "err", err)
		return
	}

	//get totalsupply for validating approve
	var total *big.Int
	if err := client.Call(&total, getDipperinRpcMethodByName("ERC20TotalSupply"), contractAdr); err != nil {
		l.Error("call GetContractInfo", "err", err)
		return
	}

	if total.Cmp(value) == -1 {
		ts, _ := InterToDecimal((*hexutil.Big)(total), decimal)
		l.Error("approving credit exceeding", "total", ts)
		return
	}

	gasPrice, err := MoneyValueToCSCoin(cParams[4])
	if err != nil {
		l.Error("the parameter gasPrice invalid", "err", err)
		return
	}

	gasLimit, err := strconv.Atoi(cParams[5])
	if err != nil {
		l.Error("the parameter gaLimit invalid", "err", err)
		return
	}

	//send transaction
	var resp common.Hash
	if err := client.Call(&resp, getDipperinRpcMethodByName("ERC20Approve"), contractAdr, from, to, value, gasPrice,gasLimit); err != nil {
		l.Error("ERC20Approve failed", "err", err)
		return
	}
	l.Info("ERC20Approve result", "txId", resp.Hex())
}

func (caller *rpcCaller) ERC20Balance(c *cli.Context) {
	_, cParams, err := getRpcMethodAndParam(c)
	if err != nil {
		l.Error("getRpcMethodAndParam error", "err", err)
		return
	}

	if !isParamValid(cParams, 2) {
		l.Error("parameters need：contract address,owner address")
		return
	}

	contractAdr, err := CheckAndChangeHexToAddress(cParams[0])
	if err != nil {
		l.Error("the input address is invalid", "err", err)
		return
	}

	owner, err := CheckAndChangeHexToAddress(cParams[1])
	if err != nil {
		l.Error("the input address is invalid", "err", err)
		return
	}

	//send transaction
	var resp *hexutil.Big
	if err := client.Call(&resp, getDipperinRpcMethodByName("ERC20Balance"), contractAdr, owner); err != nil {
		l.Error("call GetContractInfo", "err", err)
		return
	}

	decimal := getERC20Decimal(contractAdr)
	ts, _ := InterToDecimal(resp, decimal)
	unit := getERC20Symbol(contractAdr)

	l.Info("contract info", "address", owner, "token balance", ts+unit)
}
