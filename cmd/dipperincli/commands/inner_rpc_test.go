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
	"errors"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/core/rpcinterface"
	"testing"

	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/economymodel"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
	"os"
)

func Test_rpcCaller_GetBlockDiffVerifierInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := getRpcTestApp()
	app.Action = func(context *cli.Context) {
		c := &rpcCaller{}
		SyncStatus.Store(true)
		c.GetBlockDiffVerifierInfo(context)
	}
	assert.NoError(t, app.Run([]string{os.Args[0]}))

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}

		InnerRpcForbid = true
		caller.GetBlockDiffVerifierInfo(c)

		InnerRpcForbid = false
		SyncStatus.Store(false)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).Return(errors.New("test"))
		caller.GetBlockDiffVerifierInfo(c)

		SyncStatus.Store(true)
		caller.GetBlockDiffVerifierInfo(c)

		c.Set("p", "")
		caller.GetBlockDiffVerifierInfo(c)

		c.Set("p", "test")
		caller.GetBlockDiffVerifierInfo(c)

		c.Set("p", "1")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("test"))
		caller.GetBlockDiffVerifierInfo(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, method string, args ...interface{}) error {
			*result.(*map[economymodel.VerifierType][]common.Address) = map[economymodel.VerifierType][]common.Address{}
			return nil
		})
		caller.GetBlockDiffVerifierInfo(c)

	}
	assert.NoError(t, app.Run([]string{os.Args[0], "GetBlockDiffVerifierInfo"}))
	client = nil
}

func Test_printAddress(t *testing.T) {
	printAddress([]common.Address{common.HexToAddress("0x1234")})
}

func Test_rpcCaller_CheckVerifierType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := getRpcTestApp()
	app.Action = func(context *cli.Context) {
		c := &rpcCaller{}
		SyncStatus.Store(true)
		c.CheckVerifierType(context)
	}
	assert.NoError(t, app.Run([]string{os.Args[0]}))

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}

		InnerRpcForbid = true
		caller.CheckVerifierType(c)

		InnerRpcForbid = false
		SyncStatus.Store(false)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).Return(errors.New("test"))
		caller.CheckVerifierType(c)

		SyncStatus.Store(true)
		caller.CheckVerifierType(c)

		c.Set("p", "")
		caller.CheckVerifierType(c)

		c.Set("p", "test,test")
		caller.CheckVerifierType(c)

		c.Set("p", "1,test")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).Return(errors.New("test"))
		caller.CheckVerifierType(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), getDipperinRpcMethodByName("CurrentBlock")).DoAndReturn(func(result interface{}, method string, args ...interface{}) error {
			*result.(*rpcinterface.BlockResp) = rpcinterface.BlockResp{
				Header: model.Header{
					Number: 2,
				},
			}
			return nil
		}).Times(4)
		caller.CheckVerifierType(c)

		c.Set("p", "0,test")
		caller.CheckVerifierType(c)

		c.Set("p", "0,"+common.HexToAddress("0x1234").Hex())
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), getDipperinRpcMethodByName("GetBlockDiffVerifierInfo"), gomock.Any()).Return(errors.New("test"))
		caller.CheckVerifierType(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), getDipperinRpcMethodByName("GetBlockDiffVerifierInfo"), gomock.Any()).DoAndReturn(func(result interface{}, method string, args ...interface{}) error {
			tmp := map[economymodel.VerifierType][]common.Address{}
			tmp[economymodel.MasterVerifier] = []common.Address{
				common.HexToAddress("0x1234"),
			}
			*result.(*map[economymodel.VerifierType][]common.Address) = tmp
			return nil
		}).Times(2)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), getDipperinRpcMethodByName("GetBlockDiffVerifierInfo"), gomock.Any()).Return(nil).Times(1)
		caller.CheckVerifierType(c)
	}
	assert.NoError(t, app.Run([]string{os.Args[0], "CheckVerifierType"}))
	client = nil
}
