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
	"fmt"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/common"
	"math/big"
	"github.com/dipperin/dipperin-core/common/hexutil"
)

var createERC20ConfigJsonStr = `{"token_name":"EOS","token_decimals":18,"token_symbol":"EOS","token_total_supply":"0x33b2e3c9fd0803ce8000000","balances":{},"allowed":{}}
`

var address = common.HexToAddress("0x0224FA42f7315cD04D6774E58B54e92603e96d84")
var address1 = common.HexToAddress("0x0224FA42f7315cD04D6774E58B54e92603e96d85")
var address2 = common.HexToAddress("0x0224FA42f7315cD04D6774E58B54e92603e96d86")

func newTestToken() (*BuiltInERC20Token) {
	var token BuiltInERC20Token
	if err := util.ParseJson(createERC20ConfigJsonStr, &token); err != nil {
		panic(err.Error())
	}
	token.Owner = address
	return &token
}

func TestBuiltInERC20Token_Decimals(t *testing.T) {
	token := newTestToken()
	assert.NotNil(t, token)
	assert.NotNil(t, token)

	assert.Equal(t, 18, token.Decimals())
}

func TestBuiltInERC20Token_Transfer(t *testing.T) {
	token := newTestToken()
	assert.NotNil(t, token)
	assert.NotNil(t, token)
	// TODO: fix me
	token.Balances[address.Hex()] = token.TokenTotalSupply
	token.CurSender = address
	assert.Nil(t, token.Transfer(address1, (*hexutil.Big)(big.NewInt(2))))

	assert.Equal(t, 0, big.NewInt(2).Cmp(token.Balances[address1.Hex()]))

	sb := token.TokenTotalSupply.Sub(token.TokenTotalSupply, big.NewInt(2))

	assert.Equal(t, 0, sb.Cmp(token.Balances[address.Hex()]))

	fmt.Println(util.StringifyJson(token))
}

func TestBuiltInERC20Token_Approve(t *testing.T) {
	token := newTestToken()
	assert.NotNil(t, token)
	assert.NotNil(t, token)
	token.CurSender = address
	assert.Equal(t, true, token.Approve(address1, (*hexutil.Big)(big.NewInt(2))))

	assert.Equal(t, 0, big.NewInt(2).Cmp(token.Allowed[address.Hex()][address1.Hex()]))
}

func TestBuiltInERC20Token_Allowance(t *testing.T) {
	token := newTestToken()
	assert.NotNil(t, token)
	token.CurSender = address
	assert.Equal(t, true, token.Approve(address1, (*hexutil.Big)(big.NewInt(2))))

	assert.Equal(t, 0, big.NewInt(2).Cmp(token.Allowance(address, address1)))

	assert.Equal(t, 0, big.NewInt(0).Cmp(token.Allowance(address, address2)))
}

func TestBuiltInERC20Token_TransferFrom(t *testing.T) {
	token := newTestToken()
	token.Balances[address.Hex()] = big.NewInt(2)
	assert.NotNil(t, token)
	token.CurSender = address
	assert.Equal(t, true, token.Approve(address1, (*hexutil.Big)(big.NewInt(2))))

	assert.Equal(t, 0, big.NewInt(2).Cmp(token.Allowed[address.Hex()][address1.Hex()]))
	token.CurSender = address1
	assert.Equal(t, true, token.TransferFrom(address, address2, (*hexutil.Big)(big.NewInt(1))))

	assert.Equal(t, 0, big.NewInt(1).Cmp(token.Balances[address2.Hex()]))

	assert.Equal(t, 0, big.NewInt(1).Cmp(token.Allowed[address.Hex()][address1.Hex()]))

	assert.Equal(t, false, token.TransferFrom(address, address2, (*hexutil.Big)(big.NewInt(3))))

}

func TestBuiltInERC20Token_IsValid(t *testing.T) {
	token := &BuiltInERC20Token{}
	assert.Error(t, token.IsValid(), ContractOwnerNilErr)

	token.Owner = common.HexToAddress("1234")
	assert.Error(t, token.IsValid(), ContractNameNilErr)

	token.TokenName = "56"
	token.TokenDecimals = -1
	assert.Error(t, token.IsValid(), ContractNumErr)

	token.TokenDecimals = 3
	assert.Error(t, token.IsValid(), ContractSupplyNilErr)

	token.TokenTotalSupply = big.NewInt(-1)
	assert.Error(t, token.IsValid(), ContractSupplyLess0Err)

	token.TokenTotalSupply = big.NewInt(11)
	assert.NoError(t, token.IsValid())
}

func TestBuiltInERC20Token_newToken(t *testing.T) {
	token := newToken(big.NewInt(18), "c", 3, "s", common.HexToAddress("1234"))
	assert.NotNil(t, token)
}

func TestBuiltInERC20Token_Encode(t *testing.T) {
	token := newTestToken()
	ret := token.Encode()
	assert.NotNil(t, ret)
}

func TestBuiltInERC20Token_require(t *testing.T) {
	token := newTestToken()

	assert.Equal(t, false, token.require(common.HexToAddress("1234"), big.NewInt(-1)))

	assert.Equal(t, false, token.require(common.HexToAddress("1234"), big.NewInt(1)))

	token.Balances[common.HexToAddress("1234").Hex()] = big.NewInt(12)
	assert.Equal(t, false, token.require(common.HexToAddress("1234"), big.NewInt(30)))

	assert.Equal(t, true, token.require(common.HexToAddress("1234"), big.NewInt(3)))
}

