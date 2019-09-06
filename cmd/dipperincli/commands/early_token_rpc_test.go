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
	"os"
	"testing"

	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/urfave/cli"
)

func TestRpcCaller_TransferEDIPToDIP(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := getRpcTestApp()
	app.Action = func(context *cli.Context) {
		c := &rpcCaller{}
		c.TransferEDIPToDIP(context)
	}
	assert.NoError(t, app.Run([]string{os.Args[0]}))

	app.Action = func(context *cli.Context) {
		c := &rpcCaller{}
		c.TransferEDIPToDIP(context)

		context.Set("p", "")
		c.TransferEDIPToDIP(context)

		context.Set("p", "from,eDIPValue,gasPrice,gasLimit")
		c.TransferEDIPToDIP(context)

		context.Set("p", fmt.Sprintf("%s,eDIPValue,gasPrice,gasLimit", from))
		c.TransferEDIPToDIP(context)

		context.Set("p", fmt.Sprintf("%s,%v,gasPrice,gasLimit", from, "10dip"))
		c.TransferEDIPToDIP(context)

		context.Set("p", fmt.Sprintf("%s,%v,gasPrice,gasLimit", from, "10"))
		c.TransferEDIPToDIP(context)

		context.Set("p", fmt.Sprintf("%s,%v,%v,gasLimit", from, "10", "1wu"))
		c.TransferEDIPToDIP(context)

		context.Set("p", fmt.Sprintf("%s,%v,%v,%v", from, "10", "1wu", "1000"))
		client = NewMockRpcClient(ctrl)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		c.TransferEDIPToDIP(context)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
		c.TransferEDIPToDIP(context)
	}
	assert.NoError(t, app.Run([]string{os.Args[0], "TransferEDIPToDIP"}))
	client = nil
}

func TestRpcCaller_SetExchangeRate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := getRpcTestApp()
	app.Action = func(context *cli.Context) {
		c := &rpcCaller{}
		c.SetExchangeRate(context)
	}
	assert.NoError(t, app.Run([]string{os.Args[0]}))

	app.Action = func(context *cli.Context) {
		c := &rpcCaller{}
		c.SetExchangeRate(context)

		context.Set("p", "")
		c.SetExchangeRate(context)

		context.Set("p", "from,exchangeRate,gasPrice,gasLimit")
		c.SetExchangeRate(context)

		context.Set("p", fmt.Sprintf("%s,exchangeRate,gasPrice,gasLimit", from))
		c.SetExchangeRate(context)

		context.Set("p", fmt.Sprintf("%s,%v,gasPrice,gasLimit", from, "10dip"))
		c.SetExchangeRate(context)

		context.Set("p", fmt.Sprintf("%s,%v,gasPrice,gasLimit", from, "10"))
		c.SetExchangeRate(context)

		context.Set("p", fmt.Sprintf("%s,%v,%v,gasLimit", from, "10", "1wu"))
		c.SetExchangeRate(context)

		context.Set("p", fmt.Sprintf("%s,%v,%v,%v", from, "10", "1wu", "1000"))
		client = NewMockRpcClient(ctrl)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		c.SetExchangeRate(context)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
		c.SetExchangeRate(context)
	}
	assert.NoError(t, app.Run([]string{os.Args[0], "SetExchangeRate"}))
	client = nil
}
