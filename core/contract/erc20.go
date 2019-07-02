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

package contract

import (
	"errors"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/hexutil"
	"github.com/dipperin/dipperin-core/common/number"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/third-party/log"
	"math/big"
)

type BaseERC20 struct {
	ContractBase

	Owner            common.Address                 `json:"owner"`
	TokenName        string                         `json:"token_name"`
	TokenDecimals    int                            `json:"token_decimals"`
	TokenSymbol      string                         `json:"token_symbol"`
	TokenTotalSupply *big.Int                       `json:"token_total_supply"`
	Balances         map[string]*big.Int            `json:"balances"`
	Allowed          map[string]map[string]*big.Int `json:"allowed"`
}

/*
	dipperin built in erc20 contract
*/

// implement erc20 interface
type BuiltInERC20Token struct {
	BaseERC20
}
type builtInERC20TokenForMarshaling struct {
	Owner            common.Address          `json:"owner"`
	TokenName        string                  `json:"token_name"`
	TokenDecimals    int                     `json:"token_decimals"`
	TokenSymbol      string                  `json:"token_symbol"`
	TokenTotalSupply *hexutil.Big            `json:"token_total_supply"`
	Balances         map[string]*hexutil.Big `json:"balances"`
	// TODO:Allowed need json serialization?
	Allowed map[string]map[string]*hexutil.Big `json:"allowed"`
}

var (
	ContractOwnerNilErr    = errors.New("contract owner empty")
	ContractNameNilErr     = errors.New("contract name empty")
	ContractNumErr         = errors.New("contract decimal minus")
	ContractSupplyNilErr   = errors.New("contract TokenTotalSupply empty")
	ContractSupplyLess0Err = errors.New("contract TokenTotalSupply must more than 0")
)

func (token BuiltInERC20Token) MarshalJSON() ([]byte, error) {
	bm := &builtInERC20TokenForMarshaling{
		Owner:            token.Owner,
		TokenName:        token.TokenName,
		TokenDecimals:    token.TokenDecimals,
		TokenSymbol:      token.TokenSymbol,
		TokenTotalSupply: (*hexutil.Big)(token.TokenTotalSupply),
		Balances:         map[string]*hexutil.Big{},
		Allowed:          map[string]map[string]*hexutil.Big{},
	}
	for k, b := range token.Balances {
		bm.Balances[k] = (*hexutil.Big)(b)
	}
	for k1, v := range token.Allowed {
		if bm.Allowed[k1] == nil {
			bm.Allowed[k1] = map[string]*hexutil.Big{}
		}
		for k2, a := range v {
			bm.Allowed[k1][k2] = (*hexutil.Big)(a)
		}
	}
	return util.StringifyJsonToBytesWithErr(bm)
}

func (token *BuiltInERC20Token) UnmarshalJSON(input []byte) error {
	var bm builtInERC20TokenForMarshaling
	if err := util.ParseJsonFromBytes(input, &bm); err != nil {
		return err
	}
	token.Owner = bm.Owner
	token.TokenName = bm.TokenName
	token.TokenDecimals = bm.TokenDecimals
	token.TokenSymbol = bm.TokenSymbol
	token.TokenTotalSupply = (*big.Int)(bm.TokenTotalSupply)
	token.Balances = map[string]*big.Int{}
	token.Allowed = map[string]map[string]*big.Int{}

	for k, b := range bm.Balances {
		token.Balances[k] = (*big.Int)(b)
	}
	for k1, v := range bm.Allowed {
		if token.Allowed[k1] == nil {
			token.Allowed[k1] = map[string]*big.Int{}
		}
		for k2, a := range v {
			token.Allowed[k1][k2] = (*big.Int)(a)
		}
	}
	return nil
}

func (token *BuiltInERC20Token) IsValid() error {
	switch {
	case token.Owner.IsEmpty():
		return ContractOwnerNilErr
	case token.TokenName == "":
		return ContractNameNilErr
	case token.TokenDecimals < 0:
		return ContractNumErr
	case token.TokenTotalSupply == nil:
		return ContractSupplyNilErr
	case token.TokenTotalSupply.Cmp(big.NewInt(0)) != 1:
		return ContractSupplyLess0Err
	}
	return nil
}

// create ERC20 Token
//func NewERC20Token(config CreateERC20Config, owner common.Address) *BuiltInERC20Token {
//	return newToken(
//		config.InitAmount,
//		config.TokenName,
//		config.DecimalUnits,
//		config.TokenSymbol,
//		owner,
//	)
//
//}

// new token
func newToken(initAmount *big.Int, tokenName string, decimalUnits int, tokenSymbol string, owner common.Address) *BuiltInERC20Token {
	token := &BuiltInERC20Token{
		BaseERC20{
			Owner:            owner,
			TokenName:        tokenName,
			TokenTotalSupply: initAmount,
			TokenSymbol:      tokenSymbol,
			TokenDecimals:    decimalUnits,
			Balances:         make(map[string]*big.Int),
			Allowed:          make(map[string]map[string]*big.Int),
		},
	}

	// init balances
	token.Balances[owner.Hex()] = initAmount

	return token
}

// for save db ---> rlp no support map
func (token *BuiltInERC20Token) Encode() []byte {
	return util.StringifyJsonToBytes(token)
}

// validate Value
func (token *BuiltInERC20Token) require(address common.Address, value *big.Int) bool {
	// accout token
	aBalance := token.getBalanceForAddress(address)

	log.Debug("call address require", "aBalance", aBalance.String(), "value", value.String())

	// value more than 0
	if value.Cmp(big.NewInt(0)) <= 0 {
		return false
	}

	// step 1 check address balance > 0
	// account balance more than 0
	if aBalance.Cmp(big.NewInt(0)) <= 0 {
		return false
	}

	// step 2 check address balance >= value
	// account balance more than value
	if aBalance.Cmp(value) < 0 {
		return false
	}

	return true
}

// token config info
func GetContractConfig() *ContractInfo {
	info := &ContractInfo{
		// description
		Description: "for ICO",
		// initial parameters
		InitArgs: []*ContractArg{
			{Name: "owner", Description: "owner", ArgType: "common.Address"},
			{Name: "token_name", Description: "name", ArgType: "string"},
			{Name: "token_decimals", Description: "decimal", ArgType: "[]byte"},
			{Name: "token_symbol", Description: "symbol", ArgType: "string"},
			{Name: "token_total_supply", Description: "total supply", ArgType: "[]byte"},
			//{Name: "balances", Description: "balance", ArgType: "map[string][]byte"},
			//{Name: "allowed", Description: "", ArgType: "map[string]map[string]*big.Int"},
		},
		// the contract method
		Methods: []*ContractMethod{
			// configuration infomation
			{Name: "Name", Description: "get contract name", Args: []*ContractArg{}, Return: &ContractArg{Name: "name", Description: "contract name", ArgType: "string"}},
			{Name: "Symbol", Description: "get contract symbol", Args: []*ContractArg{}, Return: &ContractArg{Name: "symbol", Description: "contract symbol", ArgType: "string"}},
			{Name: "Decimals", Description: "get decimal", Args: []*ContractArg{}, Return: &ContractArg{Name: "decimals", Description: "decimal", ArgType: "*big.Int"}},
			{Name: "TotalSupply", Description: "get total supply", Args: []*ContractArg{}, Return: &ContractArg{Name: "totalSupply", Description: "total supply", ArgType: "*big.Int"}},
			// readonly infomation
			{Name: "BalanceOf", Description: "check account token balance", Args: []*ContractArg{
				{Name: "address", Description: "account address", ArgType: "common.Address"},
			}, Return: &ContractArg{Name: "balance", Description: "balance", ArgType: "*big.Int"}},
			{Name: "Allowance", Description: "token allowance", Args: []*ContractArg{
				{Name: "ownerAddress", Description: "token owner address", ArgType: "common.Address"},
				{Name: "spenderAddress", Description: "third party address", ArgType: "common.Address"},
			}, Return: &ContractArg{Name: "limit", Description: "allowance", ArgType: "*big.Int"}},
			// need transaction
			{Name: "Transfer", Description: "token transfer", Args: []*ContractArg{
				{Name: "toAddress", Description: "receiver address", ArgType: "common.Address"},
				{Name: "value", Description: "amount", ArgType: "*big.Int"},
			}, Return: &ContractArg{Name: "err", Description: "result", ArgType: "error"}, TxMethod: true},
			{Name: "TransferFrom", Description: "3rd party transfer token", Args: []*ContractArg{
				{Name: "fromAddress", Description: "token owner address", ArgType: "common.Address"},
				{Name: "toAddress", Description: "receiver address", ArgType: "common.Address"},
				{Name: "value", Description: "amount", ArgType: "*big.Int"},
			}, Return: &ContractArg{Name: "err", Description: "result", ArgType: "error"}, TxMethod: true},
			{Name: "Approve", Description: "approve token", Args: []*ContractArg{
				{Name: "spenderAddress", Description: "supplier address", ArgType: "common.Address"},
				{Name: "value", Description: "allowance amount", ArgType: "*big.Int"},
			}, Return: &ContractArg{Name: "err", Description: "result", ArgType: "error"}, TxMethod: true},
		},
	}

	return info
}

// token name
func (token *BuiltInERC20Token) Name() string {
	return token.TokenName
}

// token symbol
func (token *BuiltInERC20Token) Symbol() string {
	return token.TokenSymbol
}

// token decimals
func (token *BuiltInERC20Token) Decimals() int {
	return token.TokenDecimals
}

// Token total supply
func (token *BuiltInERC20Token) TotalSupply() *big.Int {
	return token.TokenTotalSupply
}

// check token Balance
func (token *BuiltInERC20Token) BalanceOf(address common.Address) *hexutil.Big {
	tmpBalance := token.getBalanceForAddress(address)
	//cslog.Debug().Str("addr", address.Hex()).Interface("balance byte", token.Balances[address.Hex()]).Interface("balance", tmpBalance).Msg("balance of erc20 address")
	log.Info("balance here", "adr", address)
	return (*hexutil.Big)(big.NewInt(0).Set(tmpBalance))
}

// transfer token
//func (token *BuiltInERC20Token) Transfer(toAddress common.Address, value *big.Int) error {
func (token *BuiltInERC20Token) Transfer(toAddress common.Address, hValue *hexutil.Big) error {
	log.Debug("call ERC20 Transfer")
	value := (*big.Int)(hValue)
	senderAddress := token.CurSender
	// check value > sender balance // or sender balance == 0
	if !token.require(senderAddress, value) {
		return errors.New("remainder not enough，addr:" + senderAddress.Hex())
	}

	sBalance := token.getBalanceForAddress(senderAddress)
	tBalance := token.getBalanceForAddress(toAddress)

	token.Balances[senderAddress.Hex()] = sBalance.Sub(sBalance, value)

	token.Balances[toAddress.Hex()] = tBalance.Add(tBalance, value)

	log.Debug("ERC20 transfer", "from address", senderAddress.Hex(), "to address", toAddress.Hex())
	// TODO record operation

	return nil

}

// acquire the balance of an address
func (token *BuiltInERC20Token) getBalanceForAddress(addr common.Address) *big.Int {
	balance := token.Balances[addr.Hex()]
	if balance == nil {
		balance = big.NewInt(0)
	}
	return balance
}

//  transfer token from third party
//func (token *BuiltInERC20Token) TransferFrom(fromAddress, toAddress common.Address, value *big.Int) bool {
func (token *BuiltInERC20Token) TransferFrom(fromAddress, toAddress common.Address, hValue *hexutil.Big) bool {
	senderAddress := token.CurSender
	allowance := token.Allowance(fromAddress, senderAddress)
	value := (*big.Int)(hValue)
	// check value
	if !(allowance.Cmp(value) >= 0) {
		return false
	}

	if !token.require(fromAddress, value) {
		return false
	}

	fBalance := token.getBalanceForAddress(fromAddress)

	tBalance := token.getBalanceForAddress(toAddress)

	token.Balances[fromAddress.Hex()] = fBalance.Sub(fBalance, value)

	token.Balances[toAddress.Hex()] = tBalance.Add(tBalance, value)

	if allowance.Cmp(big.NewInt(0).SetBytes(number.MaxUint256.Bytes())) < 0 {
		if token.Allowed[fromAddress.Hex()] != nil && token.Allowed[fromAddress.Hex()][senderAddress.Hex()] != nil {
			token.Allowed[fromAddress.Hex()][senderAddress.Hex()] = allowance.Sub(allowance, value)
		}
	}

	// TODO record operation

	return true

}

// approve Token
func (token *BuiltInERC20Token) Approve(spenderAddress common.Address, hValue *hexutil.Big) bool {
	senderAddress := token.CurSender
	// step 1 check map is nil
	if token.Allowed[senderAddress.Hex()] == nil {
		am := make(map[string]*big.Int)
		token.Allowed[senderAddress.Hex()] = am
	}

	value := (*big.Int)(hValue)
	token.Allowed[senderAddress.Hex()][spenderAddress.Hex()] = value

	// TODO record operation

	return true

}

// check token allowance
func (token *BuiltInERC20Token) Allowance(ownerAddress, spenderAddress common.Address) *big.Int {
	// step 1 check map is nil
	if token.Allowed[ownerAddress.Hex()] == nil {
		return big.NewInt(0)
	}

	// step 2 check map map is nil
	if token.Allowed[ownerAddress.Hex()][spenderAddress.Hex()] == nil {
		return big.NewInt(0)
	}

	return token.Allowed[ownerAddress.Hex()][spenderAddress.Hex()]
}
