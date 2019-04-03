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

	"github.com/dipperin/dipperin-core/common"
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

func Test_rpcCaller_AnnounceERC20(t *testing.T) {
	app := getRpcTestApp()
	app.Action = func(context *cli.Context) {
		c := &rpcCaller{}
		c.AnnounceERC20(context)

		wrapRpcArgs(context, "1", "")
		c.AnnounceERC20(context)

		wrapRpcArgs(context, "1", "x,y,z,z,z,z")
		c.AnnounceERC20(context)

		wrapRpcArgs(context, "1", "0x00005033874289F4F823A896700D94274683535cF0E1,y,z,z,z,z")
		c.AnnounceERC20(context)

		wrapRpcArgs(context, "1", "0x00005033874289F4F823A896700D94274683535cF0E1,y,z,10,z1,z2")
		c.AnnounceERC20(context)

		wrapRpcArgs(context, "1", "0x00005033874289F4F823A896700D94274683535cF0E1,y,z,10,1,z2")
		c.AnnounceERC20(context)

		wrapRpcArgs(context, "1", "0x00005033874289F4F823A896700D94274683535cF0E1,y,z,10,-1,z2")
		c.AnnounceERC20(context)

		assert.Panics(t, func() {
			wrapRpcArgs(context, "1", "0x00005033874289F4F823A896700D94274683535cF0E1,y,z,10,1,2")
			c.AnnounceERC20(context)
		})
	}
	assert.NoError(t, app.Run([]string{ os.Args[0] }))
}

func Test_rpcCaller_ERC20TotalSupply(t *testing.T) {
	app := getRpcTestApp()
	app.Action = func(context *cli.Context) {
		c := &rpcCaller{}
		c.ERC20TotalSupply(context)

		wrapRpcArgs(context, "1", "")
		c.ERC20TotalSupply(context)
	}
	assert.NoError(t, app.Run([]string{ os.Args[0] }))
}

func Test_rpcCaller_ERC20Transfer(t *testing.T) {
	app := getRpcTestApp()
	app.Action = func(context *cli.Context) {
		c := &rpcCaller{}
		c.ERC20Transfer(context)

		wrapRpcArgs(context, "ERC20Transfer", "")
		c.ERC20Transfer(context)

		wrapRpcArgs(context, "ERC20Transfer", "w")
		c.ERC20Transfer(context)

		wrapRpcArgs(context, "ERC20Transfer", "w,x,y,z,h")
		c.ERC20Transfer(context)

		wrapRpcArgs(context, "ERC20Transfer", "0x00100f35adf022a8aaAbef59abB97665788CDdbA30e3,x,y,z,h")
		c.ERC20Transfer(context)

		wrapRpcArgs(context, "ERC20Transfer", "0x00100f35adf022a8aaAbef59abB97665788CDdbA30e3,0x0000D07252C7A396Cc444DC0196A8b43c1A4B6c53532,y,z,h")
		c.ERC20Transfer(context)

		assert.Panics(t, func() {
			wrapRpcArgs(context, "ERC20Transfer", "0x00100f35adf022a8aaAbef59abB97665788CDdbA30e3,0x0000D07252C7A396Cc444DC0196A8b43c1A4B6c53532,0x0000B04985A7ccc00ab023d9bC40E241F9DF0379d8c4,z,h")
			c.ERC20Transfer(context)
		})
		//wrapRpcArgs(context, "1", "0x00100f35adf022a8aaAbef59abB97665788CDdbA30e3,0x0000D07252C7A396Cc444DC0196A8b43c1A4B6c53532,0x0000B04985A7ccc00ab023d9bC40E241F9DF0379d8c4,4,0.00001")
		//c.ERC20Transfer(context)
	}
	assert.NoError(t, app.Run([]string{ os.Args[0] }))
}

func Test_rpcCaller_ERC20TransferFrom(t *testing.T) {
	app := getRpcTestApp()
	app.Action = func(context *cli.Context) {
		c := &rpcCaller{}
		c.ERC20TransferFrom(context)

		wrapRpcArgs(context, "ERC20TransferFrom", "")
		c.ERC20TransferFrom(context)

		wrapRpcArgs(context, "ERC20TransferFrom", "x")
		c.ERC20TransferFrom(context)

		wrapRpcArgs(context, "ERC20TransferFrom", "0x00100f35adf022a8aaAbef59abB97665788CDdbA30e3,v,w,x,y,z")
		c.ERC20TransferFrom(context)

		wrapRpcArgs(context, "ERC20TransferFrom", "0x00100f35adf022a8aaAbef59abB97665788CDdbA30e3,0x0000D07252C7A396Cc444DC0196A8b43c1A4B6c53532,w,x,y,z")
		c.ERC20TransferFrom(context)

		wrapRpcArgs(context, "ERC20TransferFrom", "0x00100f35adf022a8aaAbef59abB97665788CDdbA30e3,0x0000D07252C7A396Cc444DC0196A8b43c1A4B6c53532,0x0000D07252C7A396Cc444DC0196A8b43c1A4B6c53532,x,y,z")
		c.ERC20TransferFrom(context)

		assert.Panics(t, func() {
			wrapRpcArgs(context, "ERC20TransferFrom", "0x00100f35adf022a8aaAbef59abB97665788CDdbA30e3,0x0000D07252C7A396Cc444DC0196A8b43c1A4B6c53532,0x0000D07252C7A396Cc444DC0196A8b43c1A4B6c53532,0x0000D07252C7A396Cc444DC0196A8b43c1A4B6c53532,y,z")
			c.ERC20TransferFrom(context)
		})
	}
	assert.NoError(t, app.Run([]string{ os.Args[0] }))
}

func Test_rpcCaller_ERC20TokenName(t *testing.T) {
	app := getRpcTestApp()
	app.Action = func(context *cli.Context) {
		c := &rpcCaller{}
		c.ERC20TokenName(context)

		wrapRpcArgs(context, "ERC20TokenName", "")
		c.ERC20TokenName(context)

		wrapRpcArgs(context, "ERC20TokenName", "x")
		c.ERC20TokenName(context)

		assert.Panics(t, func() {
			wrapRpcArgs(context, "ERC20TokenName", "0x0010a0dC0928eA10F28ceB47eb2de950195789eb9E96")
			c.ERC20TokenName(context)
		})
	}
	assert.NoError(t, app.Run([]string{ os.Args[0] }))
}

func Test_rpcCaller_ERC20TokenSymbol(t *testing.T) {
	app := getRpcTestApp()
	app.Action = func(context *cli.Context) {
		c := &rpcCaller{}
		c.ERC20TokenSymbol(context)

		wrapRpcArgs(context, "ERC20TokenSymbol", "")
		c.ERC20TokenSymbol(context)

		wrapRpcArgs(context, "ERC20TokenSymbol", "x")
		c.ERC20TokenSymbol(context)

		assert.Panics(t, func() {
			wrapRpcArgs(context, "ERC20TokenSymbol", "0x0010a0dC0928eA10F28ceB47eb2de950195789eb9E96")
			c.ERC20TokenSymbol(context)
		})
	}
	assert.NoError(t, app.Run([]string{ os.Args[0] }))
}

func Test_getERC20Symbol(t *testing.T) {
	assert.Panics(t, func() {
		getERC20Symbol(common.HexToAddress("0x123"))
	})
}

func Test_getERC20Decimal(t *testing.T) {
	//type args struct {
	//	contractAdr common.Address
	//}
	//tests := []struct {
	//	name string
	//	args args
	//	want int
	//}{
	//	// TODO: Add test cases.
	//}
	//for _, tt := range tests {
	//	t.Run(tt.name, func(t *testing.T) {
	//		if got := getERC20Decimal(tt.args.contractAdr); got != tt.want {
	//			t.Errorf("getERC20Decimal() = %v, want %v", got, tt.want)
	//		}
	//	})
	//}

	assert.Panics(t, func() {
		getERC20Decimal(common.HexToAddress("0x123"))
	})
}

func Test_rpcCaller_ERC20TokenDecimals(t *testing.T) {
	app := getRpcTestApp()
	app.Action = func(context *cli.Context) {
		c := &rpcCaller{}
		c.ERC20TokenDecimals(context)

		wrapRpcArgs(context, "ERC20TokenDecimals", "")
		c.ERC20TokenDecimals(context)

		wrapRpcArgs(context, "ERC20TokenDecimals", "x")
		c.ERC20TokenDecimals(context)

		assert.Panics(t, func() {
			wrapRpcArgs(context, "ERC20TokenDecimals", "0x0010a0dC0928eA10F28ceB47eb2de950195789eb9E96")
			c.ERC20TokenDecimals(context)
		})
	}
	assert.NoError(t, app.Run([]string{ os.Args[0] }))
}

func Test_convert(t *testing.T) {
	ret := convert(common.FromHex("123456"))
	assert.Equal(t, ret, 0x123456)
}

func Test_rpcCaller_ERC20GetInfo(t *testing.T) {
	app := getRpcTestApp()
	app.Action = func(context *cli.Context) {
		c := &rpcCaller{}
		c.ERC20GetInfo(context)

		wrapRpcArgs(context, "ERC20GetInfo", "")
		c.ERC20GetInfo(context)

		wrapRpcArgs(context, "ERC20GetInfo", "x")
		c.ERC20GetInfo(context)

		assert.Panics(t, func() {
			wrapRpcArgs(context, "ERC20GetInfo", "0x0010a0dC0928eA10F28ceB47eb2de950195789eb9E96")
			c.ERC20GetInfo(context)
		})
	}
	assert.NoError(t, app.Run([]string{ os.Args[0] }))
}

func Test_rpcCaller_ERC20Allowance(t *testing.T) {
	app := getRpcTestApp()
	app.Action = func(context *cli.Context) {
		c := &rpcCaller{}
		c.ERC20Allowance(context)

		wrapRpcArgs(context, "ERC20Allowance", "")
		c.ERC20Allowance(context)

		wrapRpcArgs(context, "ERC20Allowance", "x,y,z")
		c.ERC20Allowance(context)

		wrapRpcArgs(context, "ERC20Allowance", "0x0010a0dC0928eA10F28ceB47eb2de950195789eb9E96,y,z")
		c.ERC20Allowance(context)

		wrapRpcArgs(context, "ERC20Allowance", "0x0010a0dC0928eA10F28ceB47eb2de950195789eb9E96,0x0010a0dC0928eA10F28ceB47eb2de950195789eb9E96,z")
		c.ERC20Allowance(context)

		assert.Panics(t, func() {
			wrapRpcArgs(context, "ERC20Allowance", "0x0010a0dC0928eA10F28ceB47eb2de950195789eb9E96,0x0010a0dC0928eA10F28ceB47eb2de950195789eb9E96,0x0010a0dC0928eA10F28ceB47eb2de950195789eb9E96")
			c.ERC20Allowance(context)
		})
	}
	assert.NoError(t, app.Run([]string{ os.Args[0] }))
}

func Test_rpcCaller_ERC20Approve(t *testing.T) {
	app := getRpcTestApp()
	app.Action = func(context *cli.Context) {
		c := &rpcCaller{}
		c.ERC20Approve(context)

		wrapRpcArgs(context, "ERC20Approve", "")
		c.ERC20Approve(context)

		wrapRpcArgs(context, "ERC20Approve", "v,w,x,y,z")
		c.ERC20Approve(context)

		wrapRpcArgs(context, "ERC20Approve", "0x0010a0dC0928eA10F28ceB47eb2de950195789eb9E96,w,x,y,z")
		c.ERC20Approve(context)

		wrapRpcArgs(context, "ERC20Approve", "0x0010a0dC0928eA10F28ceB47eb2de950195789eb9E96,0x0010a0dC0928eA10F28ceB47eb2de950195789eb9E96,x,y,z")
		c.ERC20Approve(context)

		assert.Panics(t, func() {
			wrapRpcArgs(context, "ERC20Approve", "0x0010a0dC0928eA10F28ceB47eb2de950195789eb9E96,0x0010a0dC0928eA10F28ceB47eb2de950195789eb9E96,0x0010a0dC0928eA10F28ceB47eb2de950195789eb9E96,y,z")
			c.ERC20Approve(context)
		})
	}
	assert.NoError(t, app.Run([]string{ os.Args[0] }))
}

func Test_rpcCaller_ERC20Balance(t *testing.T) {
	app := getRpcTestApp()
	app.Action = func(context *cli.Context) {
		c := &rpcCaller{}
		c.ERC20Balance(context)

		wrapRpcArgs(context, "ERC20Balance", "")
		c.ERC20Balance(context)

		wrapRpcArgs(context, "ERC20Balance", "x,y")
		c.ERC20Balance(context)

		wrapRpcArgs(context, "ERC20Balance", "0x0010a0dC0928eA10F28ceB47eb2de950195789eb9E96,y")
		c.ERC20Balance(context)

		assert.Panics(t, func() {
			wrapRpcArgs(context, "ERC20Balance", "0x0010a0dC0928eA10F28ceB47eb2de950195789eb9E96,0x0010a0dC0928eA10F28ceB47eb2de950195789eb9E96")
			c.ERC20Balance(context)
		})
	}
	assert.NoError(t, app.Run([]string{ os.Args[0] }))
}


func Test_ERC20Size(t *testing.T) {

}
