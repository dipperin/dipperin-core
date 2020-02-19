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
	"github.com/stretchr/testify/assert"
	"math/big"
	"os"
	"testing"

	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/consts"
	"github.com/dipperin/dipperin-core/common/hexutil"
	"github.com/golang/mock/gomock"
	"github.com/urfave/cli"
)

func Test_buildERC20Token(t *testing.T) {
	c := buildERC20Token(common.HexToAddress("0x123"), "xxx", "EOS", big.NewInt(123), 9)
	assert.NotNil(t, c)
}

func Test_isParamValid(t *testing.T) {
	assert.True(t, isParamValid([]string{"1"}, 1))
	assert.False(t, isParamValid([]string{"1", "2"}, 1))
	assert.False(t, isParamValid([]string{"", "2"}, 2))
}

func TestRpcCaller_AnnounceERC20(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := getRpcTestApp()
	app.Action = func(context *cli.Context) {
		c := &rpcCaller{}
		c.AnnounceERC20(context)
	}
	assert.NoError(t, app.Run([]string{os.Args[0]}))

	app.Action = func(context *cli.Context) {
		c := &rpcCaller{}
		c.AnnounceERC20(context)

		context.Set("p", "")
		c.AnnounceERC20(context)

		context.Set("p", "owner_address,name,symbol,supply,decimal,gasPrice,gasLimit")
		c.AnnounceERC20(context)

		context.Set("p", fmt.Sprintf("%s,name,symbol,supply,decimal,gasPrice,gasLimit", from))
		c.AnnounceERC20(context)

		context.Set("p", fmt.Sprintf("%s,name,symbol,supply,%s,gasPrice,gasLimit", from, "20"))
		c.AnnounceERC20(context)

		context.Set("p", fmt.Sprintf("%s,name,symbol,supply,%s,gasPrice,gasLimit", from, "18"))
		c.AnnounceERC20(context)

		context.Set("p", fmt.Sprintf("%s,name,symbol,%s,%s,gasPrice,gasLimit", from, "1000", "18"))
		c.AnnounceERC20(context)

		context.Set("p", fmt.Sprintf("%s,name,symbol,%s,%s,%v,gasLimit", from, "1000", "18", "10wu"))
		c.AnnounceERC20(context)

		context.Set("p", fmt.Sprintf("%s,name,symbol,%s,%s,%v,%v", from, "1000", "18", "10wu", "100"))
		client = NewMockRpcClient(ctrl)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
		c.AnnounceERC20(context)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		c.AnnounceERC20(context)
	}
	assert.NoError(t, app.Run([]string{os.Args[0], "AnnounceERC20"}))
	client = nil
}

func TestRpcCaller_ERC20Transfer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := getRpcTestApp()
	app.Action = func(context *cli.Context) {
		c := &rpcCaller{}
		c.ERC20Transfer(context)
	}
	assert.NoError(t, app.Run([]string{os.Args[0]}))

	app.Action = func(context *cli.Context) {
		c := &rpcCaller{}
		context.Set("p", "")
		c.ERC20Transfer(context)

		context.Set("p", "contractAddr,owner,to,amount,gasPrice,gasLimit")
		c.ERC20Transfer(context)

		context.Set("p", fmt.Sprintf("%s,owner,to,amount,gasPrice,gasLimit", contractAddr))
		c.ERC20Transfer(context)

		context.Set("p", fmt.Sprintf("%s,%s,to,amount,gasPrice,gasLimit", contractAddr, from))
		c.ERC20Transfer(context)

		client = NewMockRpcClient(ctrl)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).SetArg(0, consts.DIPDecimalBits)
		context.Set("p", fmt.Sprintf("%s,%s,%s,amount,gasPrice,gasLimit", contractAddr, from, to))
		c.ERC20Transfer(context)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).SetArg(0, consts.DIPDecimalBits)
		context.Set("p", fmt.Sprintf("%s,%s,%s,%v,gasPrice,gasLimit", contractAddr, from, to, "10"))
		c.ERC20Transfer(context)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).SetArg(0, consts.DIPDecimalBits)
		context.Set("p", fmt.Sprintf("%s,%s,%s,%v,%v,gasLimit", contractAddr, from, to, "10", "10wu"))
		c.ERC20Transfer(context)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).SetArg(0, consts.DIPDecimalBits)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		context.Set("p", fmt.Sprintf("%s,%s,%s,%v,%v,%v", contractAddr, from, to, "10", "10wu", "100"))
		c.ERC20Transfer(context)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).SetArg(0, consts.DIPDecimalBits)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
		c.ERC20Transfer(context)
	}
	assert.NoError(t, app.Run([]string{os.Args[0], "ERC20Transfer"}))
	client = nil
}

func TestRpcCaller_ERC20TransferFrom(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := getRpcTestApp()
	app.Action = func(context *cli.Context) {
		c := &rpcCaller{}
		c.ERC20TransferFrom(context)
	}
	assert.NoError(t, app.Run([]string{os.Args[0]}))

	app.Action = func(context *cli.Context) {
		c := &rpcCaller{}
		context.Set("p", "")
		c.ERC20TransferFrom(context)

		context.Set("p", "contractAddr,owner,from,to,amount,gasPrice,gasLimit")
		c.ERC20TransferFrom(context)

		context.Set("p", fmt.Sprintf("%s,owner,from,to,amount,gasPrice,gasLimit", contractAddr))
		c.ERC20TransferFrom(context)

		context.Set("p", fmt.Sprintf("%s,%s,from,to,amount,gasPrice,gasLimit", contractAddr, from))
		c.ERC20TransferFrom(context)

		context.Set("p", fmt.Sprintf("%s,%s,%s,to,amount,gasPrice,gasLimit", contractAddr, from, from))
		c.ERC20TransferFrom(context)

		client = NewMockRpcClient(ctrl)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).SetArg(0, consts.DIPDecimalBits)
		context.Set("p", fmt.Sprintf("%s,%s,%s,%s,amount,gasPrice,gasLimit", contractAddr, from, from, to))
		c.ERC20TransferFrom(context)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).SetArg(0, consts.DIPDecimalBits)
		context.Set("p", fmt.Sprintf("%s,%s,%s,%s,%s,gasPrice,gasLimit", contractAddr, from, from, to, "10"))
		c.ERC20TransferFrom(context)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).SetArg(0, consts.DIPDecimalBits)
		context.Set("p", fmt.Sprintf("%s,%s,%s,%s,%s,%s,gasLimit", contractAddr, from, from, to, "10", "10wu"))
		c.ERC20TransferFrom(context)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).SetArg(0, consts.DIPDecimalBits)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
		context.Set("p", fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s", contractAddr, from, from, to, "10", "10wu", "100"))
		c.ERC20TransferFrom(context)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).SetArg(0, consts.DIPDecimalBits)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		c.ERC20TransferFrom(context)
	}
	assert.NoError(t, app.Run([]string{os.Args[0], "ERC20TransferFrom"}))
	client = nil
}

func TestRpcCaller_ERC20Allowance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := getRpcTestApp()
	app.Action = func(context *cli.Context) {
		c := &rpcCaller{}
		c.ERC20Allowance(context)
	}
	assert.NoError(t, app.Run([]string{os.Args[0]}))

	app.Action = func(context *cli.Context) {
		c := &rpcCaller{}
		context.Set("p", "")
		c.ERC20Allowance(context)

		context.Set("p", "contractAddr,owner,spender")
		c.ERC20Allowance(context)

		context.Set("p", fmt.Sprintf("%s,owner,spender", contractAddr))
		c.ERC20Allowance(context)

		context.Set("p", fmt.Sprintf("%s,%s,spender", contractAddr, from))
		c.ERC20Allowance(context)

		context.Set("p", fmt.Sprintf("%s,%s,%s", contractAddr, from, to))
		client = NewMockRpcClient(ctrl)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
		c.ERC20Allowance(context)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).SetArg(0, big.NewInt(1e18))
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).SetArg(0, consts.DIPDecimalBits)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).SetArg(0, "dip")
		c.ERC20Allowance(context)
	}
	assert.NoError(t, app.Run([]string{os.Args[0], "ERC20Allowance"}))
	client = nil
}

func TestRpcCaller_ERC20Approve(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := getRpcTestApp()
	app.Action = func(context *cli.Context) {
		c := &rpcCaller{}
		c.ERC20Approve(context)
	}
	assert.NoError(t, app.Run([]string{os.Args[0]}))

	app.Action = func(context *cli.Context) {
		c := &rpcCaller{}
		context.Set("p", "")
		c.ERC20Approve(context)

		context.Set("p", "contractAddr,owner,to,amount,gasPrice,gasLimit")
		c.ERC20Approve(context)

		context.Set("p", fmt.Sprintf("%s,owner,to,amount,gasPrice,gasLimit", contractAddr))
		c.ERC20Approve(context)

		context.Set("p", fmt.Sprintf("%s,%s,to,amount,gasPrice,gasLimit", contractAddr, from))
		c.ERC20Approve(context)

		client = NewMockRpcClient(ctrl)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).SetArg(0, consts.KUnitDecimalBits)
		context.Set("p", fmt.Sprintf("%s,%s,%s,amount,gasPrice,gasLimit", contractAddr, from, to))
		c.ERC20Approve(context)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).SetArg(0, consts.KUnitDecimalBits)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
		context.Set("p", fmt.Sprintf("%s,%s,%s,%s,gasPrice,gasLimit", contractAddr, from, to, "10"))
		c.ERC20Approve(context)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).SetArg(0, consts.KUnitDecimalBits)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).SetArg(0, big.NewInt(10))
		context.Set("p", fmt.Sprintf("%s,%s,%s,%s,gasPrice,gasLimit", contractAddr, from, to, "10"))
		c.ERC20Approve(context)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).SetArg(0, consts.KUnitDecimalBits)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).SetArg(0, big.NewInt(1e5))
		context.Set("p", fmt.Sprintf("%s,%s,%s,%s,gasPrice,gasLimit", contractAddr, from, to, "10"))
		c.ERC20Approve(context)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).SetArg(0, consts.KUnitDecimalBits)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).SetArg(0, big.NewInt(1e5))
		context.Set("p", fmt.Sprintf("%s,%s,%s,%s,%s,gasLimit", contractAddr, from, to, "10", "1wu"))
		c.ERC20Approve(context)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).SetArg(0, consts.KUnitDecimalBits)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).SetArg(0, big.NewInt(1e5))
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		context.Set("p", fmt.Sprintf("%s,%s,%s,%s,%s,%s", contractAddr, from, to, "10", "1wu", "100"))
		c.ERC20Approve(context)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).SetArg(0, consts.KUnitDecimalBits)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).SetArg(0, big.NewInt(1e5))
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
		c.ERC20Approve(context)
	}
	assert.NoError(t, app.Run([]string{os.Args[0], "ERC20Approve"}))
	client = nil
}

func TestRpcCaller_ERC20TotalSupply(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := getRpcTestApp()
	app.Action = func(context *cli.Context) {
		c := &rpcCaller{}
		c.ERC20TotalSupply(context)
	}
	assert.NoError(t, app.Run([]string{os.Args[0]}))

	app.Action = func(context *cli.Context) {
		c := &rpcCaller{}
		context.Set("p", "")
		c.ERC20TotalSupply(context)

		context.Set("p", "contractAddr")
		c.ERC20TotalSupply(context)

		context.Set("p", contractAddr)
		client = NewMockRpcClient(ctrl)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
		c.ERC20TotalSupply(context)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).SetArg(0, big.NewInt(1e18))
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).SetArg(0, consts.DIPDecimalBits)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).SetArg(0, "dip")
		c.ERC20TotalSupply(context)
	}
	assert.NoError(t, app.Run([]string{os.Args[0], "ERC20TotalSupply"}))
	client = nil
}

func TestRpcCaller_ERC20Balance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := getRpcTestApp()
	app.Action = func(context *cli.Context) {
		c := &rpcCaller{}
		c.ERC20Balance(context)
	}
	assert.NoError(t, app.Run([]string{os.Args[0]}))

	app.Action = func(context *cli.Context) {
		c := &rpcCaller{}
		context.Set("p", "")
		c.ERC20Balance(context)

		context.Set("p", "contractAddr,owner")
		c.ERC20Balance(context)

		context.Set("p", fmt.Sprintf("%s,owner", contractAddr))
		c.ERC20Balance(context)

		context.Set("p", fmt.Sprintf("%s,%s", contractAddr, from))
		client = NewMockRpcClient(ctrl)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
		c.ERC20Balance(context)

		balance := big.NewInt(1e18)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).SetArg(0, (*hexutil.Big)(balance))
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).SetArg(0, consts.DIPDecimalBits)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).SetArg(0, "dip")
		c.ERC20Balance(context)
	}
	assert.NoError(t, app.Run([]string{os.Args[0], "ERC20Balance"}))
	client = nil
}

func TestRpcCaller_ERC20TokenName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := getRpcTestApp()
	app.Action = func(context *cli.Context) {
		c := &rpcCaller{}
		c.ERC20TokenName(context)
	}
	assert.NoError(t, app.Run([]string{os.Args[0]}))

	app.Action = func(context *cli.Context) {
		c := &rpcCaller{}
		context.Set("p", "")
		c.ERC20TokenName(context)

		context.Set("p", "contractAddr")
		c.ERC20TokenName(context)

		context.Set("p", contractAddr)
		client = NewMockRpcClient(ctrl)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
		c.ERC20TokenName(context)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).SetArg(0, "name")
		c.ERC20TokenName(context)
	}
	assert.NoError(t, app.Run([]string{os.Args[0], "ERC20TokenName"}))
	client = nil
}

func TestRpcCaller_ERC20TokenSymbol(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := getRpcTestApp()
	app.Action = func(context *cli.Context) {
		c := &rpcCaller{}
		c.ERC20TokenSymbol(context)
	}
	assert.NoError(t, app.Run([]string{os.Args[0]}))

	app.Action = func(context *cli.Context) {
		c := &rpcCaller{}
		context.Set("p", "")
		c.ERC20TokenSymbol(context)

		context.Set("p", "contractAddr")
		c.ERC20TokenSymbol(context)

		context.Set("p", contractAddr)
		client = NewMockRpcClient(ctrl)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
		c.ERC20TokenSymbol(context)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).SetArg(0, "symbol")
		c.ERC20TokenSymbol(context)
	}
	assert.NoError(t, app.Run([]string{os.Args[0], "ERC20TokenSymbol"}))
	client = nil
}

func TestRpcCaller_ERC20TokenDecimals(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := getRpcTestApp()
	app.Action = func(context *cli.Context) {
		c := &rpcCaller{}
		c.ERC20TokenDecimals(context)
	}
	assert.NoError(t, app.Run([]string{os.Args[0]}))

	app.Action = func(context *cli.Context) {
		c := &rpcCaller{}
		context.Set("p", "")
		c.ERC20TokenDecimals(context)

		context.Set("p", "contractAddr")
		c.ERC20TokenDecimals(context)

		context.Set("p", contractAddr)
		client = NewMockRpcClient(ctrl)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
		c.ERC20TokenDecimals(context)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).SetArg(0, 18)
		c.ERC20TokenDecimals(context)
	}
	assert.NoError(t, app.Run([]string{os.Args[0], "ERC20TokenDecimals"}))
	client = nil
}

func TestRpcCaller_ERC20GetInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := getRpcTestApp()
	app.Action = func(context *cli.Context) {
		c := &rpcCaller{}
		c.ERC20GetInfo(context)
	}
	assert.NoError(t, app.Run([]string{os.Args[0]}))

	app.Action = func(context *cli.Context) {
		c := &rpcCaller{}
		context.Set("p", "")
		c.ERC20GetInfo(context)

		context.Set("p", "contractAddr")
		c.ERC20GetInfo(context)

		context.Set("p", contractAddr)
		client = NewMockRpcClient(ctrl)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).MaxTimes(3)
		c.ERC20GetInfo(context)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
		c.ERC20GetInfo(context)

		mapInterface := make(map[string]interface{})
		mapInterface["token_total_supply"] = "1000000000000000000"
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).SetArg(0, mapInterface)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).SetArg(0, consts.DIPDecimalBits)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).SetArg(0, "dip")
		c.ERC20GetInfo(context)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).SetArg(0, mapInterface)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr).MaxTimes(2)
		c.ERC20GetInfo(context)
	}
	assert.NoError(t, app.Run([]string{os.Args[0], "ERC20GetInfo"}))
	client = nil
}
