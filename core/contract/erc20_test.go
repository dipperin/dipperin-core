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
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/hexutil"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestGetContractConfig(t *testing.T) {
	config := GetContractConfig()
	fmt.Println(config)

	token := newTestToken()
	assert.NotNil(t, token)
	assert.Equal(t, 18, token.Decimals())
	assert.Equal(t, "EOS", token.Name())
	assert.Equal(t, "EOS", token.Symbol())
	assert.Equal(t, "1000000000000000000000000000", token.TotalSupply().String())
}

func TestBuiltInERC20Token_Transfer(t *testing.T) {
	token := newTestToken()
	assert.NotNil(t, token)

	token.Balances[aliceAddr.Hex()] = token.TokenTotalSupply
	token.CurSender = aliceAddr
	assert.NoError(t, token.Transfer(bobAddr, (*hexutil.Big)(big.NewInt(2))))
	assert.Equal(t, big.NewInt(2), token.BalanceOf(bobAddr).ToInt())

	sb := token.TokenTotalSupply.Sub(token.TokenTotalSupply, big.NewInt(2))
	assert.Equal(t, sb, token.BalanceOf(aliceAddr).ToInt())

	data := util.StringifyJsonToBytes(token)
	var resp *BuiltInERC20Token
	err := util.ParseJsonFromBytes(data, &resp)
	assert.NoError(t, err)
	assert.Equal(t, token.Allowed, resp.Allowed)
}

func TestBuiltInERC20Token_Approve(t *testing.T) {
	token := newTestToken()
	assert.NotNil(t, token)

	token.CurSender = aliceAddr
	assert.Equal(t, true, token.Approve(bobAddr, (*hexutil.Big)(big.NewInt(2))))
	assert.Equal(t, 0, big.NewInt(2).Cmp(token.Allowed[aliceAddr.Hex()][bobAddr.Hex()]))
}

func TestBuiltInERC20Token_Allowance(t *testing.T) {
	token := newTestToken()
	assert.NotNil(t, token)

	token.CurSender = aliceAddr
	assert.Equal(t, true, token.Approve(bobAddr, (*hexutil.Big)(big.NewInt(2))))
	assert.Equal(t, 0, big.NewInt(2).Cmp(token.Allowance(aliceAddr, bobAddr)))
	assert.Equal(t, 0, big.NewInt(0).Cmp(token.Allowance(aliceAddr, charlieAddr)))

	data := util.StringifyJsonToBytes(token)
	var resp *BuiltInERC20Token
	err := util.ParseJsonFromBytes(data, &resp)
	assert.NoError(t, err)
	assert.Equal(t, token.Allowed, resp.Allowed)
}

func TestBuiltInERC20Token_TransferFrom(t *testing.T) {
	token := newTestToken()
	token.Balances[aliceAddr.Hex()] = big.NewInt(2)
	token.CurSender = aliceAddr

	assert.Equal(t, true, token.Approve(bobAddr, (*hexutil.Big)(big.NewInt(2))))
	assert.Equal(t, 0, big.NewInt(2).Cmp(token.Allowed[aliceAddr.Hex()][bobAddr.Hex()]))

	token.CurSender = bobAddr
	assert.Equal(t, true, token.TransferFrom(aliceAddr, charlieAddr, (*hexutil.Big)(big.NewInt(1))))
	assert.Equal(t, 0, big.NewInt(1).Cmp(token.Balances[charlieAddr.Hex()]))
	assert.Equal(t, 0, big.NewInt(1).Cmp(token.Allowed[aliceAddr.Hex()][bobAddr.Hex()]))
	assert.Equal(t, false, token.TransferFrom(aliceAddr, charlieAddr, (*hexutil.Big)(big.NewInt(3))))

	data := util.StringifyJsonToBytes(token)
	var resp *BuiltInERC20Token
	err := util.ParseJsonFromBytes(data, &resp)
	assert.NoError(t, err)
	assert.Equal(t, token.Allowed, resp.Allowed)
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
