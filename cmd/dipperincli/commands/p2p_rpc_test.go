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
	"github.com/dipperin/dipperin-core/third-party/p2p"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/urfave/cli"
	"os"
)

func Test_getP2pRpcMethodByName(t *testing.T) {
	assert.Equal(t, getP2pRpcMethodByName("Test"), "p2p_test")
}

func Test_rpcCaller_AddPeer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := getRpcTestApp()
	app.Action = func(context *cli.Context) {
		c := &rpcCaller{}
		c.AddPeer(context)
	}
	assert.NoError(t, app.Run([]string{os.Args[0]}))

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}
		caller.AddPeer(c)

		c.Set("p", "")
		caller.AddPeer(c)

		c.Set("p", "test")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
		caller.AddPeer(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, method string, args ...interface{}) error {
			*result.(*error) = testErr
			return nil
		})
		caller.AddPeer(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		caller.AddPeer(c)
	}
	assert.NoError(t, app.Run([]string{os.Args[0], "AddPeer"}))
	client = nil
}

func Test_rpcCaller_Peers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := getRpcTestApp()
	app.Action = func(context *cli.Context) {
		c := &rpcCaller{}
		c.Peers(context)
	}
	assert.NoError(t, app.Run([]string{os.Args[0]}))

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}
		//caller.Peers(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).Return(errors.New("test"))
		caller.Peers(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, method string, args ...interface{}) error {
			*result.(*[]*p2p.PeerInfo) = []*p2p.PeerInfo{{}}
			return nil
		})
		caller.Peers(c)
	}
	assert.NoError(t, app.Run([]string{os.Args[0], "Peers"}))
	client = nil
}

//func Test_rpcCaller_Debug(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	app := getRpcTestApp()
//	app.Action = func(c *cli.Context) {
//		client = NewMockRpcClient(ctrl)
//		caller := &rpcCaller{}
//
//		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("test"))
//		caller.Debug(c)
//
//		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, method string, args ...interface{}) error {
//			tx, _ := factory.CreateTestTx()
//			m, _ := model.NewVoteMsgWithSign(uint64(1), uint64(1), common.HexToHash("0x1234"), model.PreVoteMessage, func(hash []byte) ([]byte, error) {
//				return nil, nil
//			}, common.HexToAddress("0x1234"))
//			*result.(*rpc_interface.BlockResp) = rpc_interface.BlockResp{
//				Body: model.Body{Txs: []*model.Transaction{tx}, Vers: []model.AbstractVerification{m}},
//				Header: model.Header{
//					Bloom: iblt.NewBloom(iblt.NewBloomConfig(8, 4)),
//				},
//			}
//			return nil
//		})
//		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).Return(errors.New("test"))
//		caller.Debug(c)
//
//		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, method string, args ...interface{}) error {
//			tx, _ := factory.CreateTestTx()
//			m, _ := model.NewVoteMsgWithSign(uint64(1), uint64(1), common.HexToHash("0x1234"), model.PreVoteMessage, func(hash []byte) ([]byte, error) {
//				return nil, nil
//			}, common.HexToAddress("0x1234"))
//			*result.(*rpc_interface.BlockResp) = rpc_interface.BlockResp{
//				Body: model.Body{Txs: []*model.Transaction{tx}, Vers: []model.AbstractVerification{m}},
//				Header: model.Header{
//					Bloom: iblt.NewBloom(iblt.NewBloomConfig(8, 4)),
//				},
//			}
//			return nil
//		})
//		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).Return(nil)
//		caller.Debug(c)
//	}
//	assert.NoError(t, app.Run([]string{os.Args[0]}))
//	client = nil
//}
