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
	"encoding/json"
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/hexutil"
	"github.com/dipperin/dipperin-core/core/accounts"
	"github.com/dipperin/dipperin-core/core/bloom"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/core/rpc-interface"
	model2 "github.com/dipperin/dipperin-core/core/vm/model"
	"github.com/dipperin/dipperin-core/tests/factory"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
	"io/ioutil"
	"math/big"
	"os"
	"testing"
	"time"
)

func TestInitRpcClient(t *testing.T) {
	assert.Panics(t, func() {
		InitRpcClient(12345)
	})
}

func TestInitAccountInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client = NewMockRpcClient(ctrl)
	client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).Return(testErr).Times(2)
	InitAccountInfo(chain_config.NodeTypeOfVerifier, "", "", "")
	osExit = func(code int) {}

	client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
	client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).Return(testErr).Times(2)
	InitAccountInfo(chain_config.NodeTypeOfNormal, "", "", "")
}

func TestCheckDownloaderSyncStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client = NewMockRpcClient(ctrl)
	client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).Return(testErr)
	assert.Panics(t, func() {
		CheckDownloaderSyncStatus()
	})

	client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, args ...interface{}) error {
		result = false
		return nil
	})

	assert.NotPanics(t, func() {
		CheckDownloaderSyncStatus()
	})

	client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, args ...interface{}) error {
		*result.(*bool) = true
		return nil
	}).Times(1)

	client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, args ...interface{}) error {
		*result.(*bool) = true
		return testErr
	}).Times(1)

	client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, args ...interface{}) error {
		*result.(*bool) = false
		return nil
	}).Times(1)

	CheckSyncStatusDuration = 1 * time.Millisecond
	assert.NotPanics(t, func() {
		CheckDownloaderSyncStatus()
	})

	client = nil
}

func TestRpcCall(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := getRpcTestApp()
	app.Action = func(c *cli.Context) {
		assert.Panics(t, func() {
			RpcCall(c)
		})

		client = NewMockRpcClient(ctrl)
		RpcCall(c)
	}
	assert.NoError(t, app.Run([]string{os.Args[0]}))

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		RpcCall(c)
	}
	assert.NoError(t, app.Run([]string{os.Args[0], "unknown"}))

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
		RpcCall(c)
	}

	assert.NoError(t, app.Run([]string{os.Args[0], "GetDefaultAccountBalance"}))
	client = nil
}

func Test_getRpcParamFromString(t *testing.T) {
	assert.Equal(t, getRpcParamFromString(""), []string{})
	assert.Equal(t, getRpcParamFromString("test,test1"), []string{"test", "test1"})
}

func Test_getRpcMethodAndParam(t *testing.T) {
	app := getRpcTestApp()
	app.Action = func(c *cli.Context) {
		c.Set("p", "test")
		mName, cParams, err := getRpcMethodAndParam(c)

		assert.Equal(t, mName, "test")
		assert.Equal(t, cParams, []string{"test"})
		assert.NoError(t, err)
	}

	assert.NoError(t, app.Run([]string{os.Args[0], "test"}))
}

func Test_checkSync(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client = NewMockRpcClient(ctrl)
	client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).Return(testErr)
	SyncStatus.Store(false)
	assert.Equal(t, checkSync(), true)

	SyncStatus.Store(true)
	assert.Equal(t, checkSync(), false)
	client = nil
}

func TestRpcCaller_GetDefaultAccountBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := getRpcTestApp()
	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
		caller.GetDefaultAccountBalance(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		caller.GetDefaultAccountBalance(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, args ...interface{}) error {
			*result.(*rpc_interface.CurBalanceResp) = rpc_interface.CurBalanceResp{
				Balance: (*hexutil.Big)(big.NewInt(1)),
			}
			return nil
		})

		caller.GetDefaultAccountBalance(c)

	}

	assert.NoError(t, app.Run([]string{os.Args[0]}))
	client = nil
}

func TestRpcCaller_CurrentBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := getRpcTestApp()
	app.Action = func(c *cli.Context) {
		caller := &rpcCaller{}
		caller.CurrentBalance(c)
	}
	assert.NoError(t, app.Run([]string{os.Args[0]}))

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}

		c.Set("p", "test,test")
		caller.CurrentBalance(c)

		c.Set("p", "1234")
		caller.CurrentBalance(c)

		c.Set("p", from)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
		caller.CurrentBalance(c)

		c.Set("p", "")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), getDipperinRpcMethodByName("ListWallet")).DoAndReturn(func(result interface{}, method string, args ...interface{}) error {
			*result.(*[]accounts.WalletIdentifier) = []accounts.WalletIdentifier{
				{
					WalletType: 1,
					Path:       "",
					WalletName: "",
				},
			}
			return nil
		})
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), getDipperinRpcMethodByName("ListWalletAccount"), gomock.Any()).DoAndReturn(func(result interface{}, method string, args ...interface{}) error {
			*result.(*[]accounts.Account) = []accounts.Account{
				{
					Address: fromAddr,
				},
			}
			return nil
		})
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		caller.CurrentBalance(c)

		c.Set("p", from)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, method string, args ...interface{}) error {
			*result.(*rpc_interface.CurBalanceResp) = rpc_interface.CurBalanceResp{
				Balance: (*hexutil.Big)(big.NewInt(1)),
			}
			return nil
		})
		caller.CurrentBalance(c)
	}

	assert.NoError(t, app.Run([]string{os.Args[0], "CurrentBalance"}))
	client = nil
}

func Test_printBlockInfo(t *testing.T) {
	tx, _ := factory.CreateTestTx()
	m, _ := model.NewVoteMsgWithSign(uint64(1), uint64(1), common.HexToHash("0x1234"), model.PreVoteMessage, func(hash []byte) ([]byte, error) {
		return nil, nil
	}, common.HexToAddress("0x1234"))
	respBlock := rpc_interface.BlockResp{
		Body: model.Body{Txs: []*model.Transaction{tx}, Vers: []model.AbstractVerification{m}},
		Header: model.Header{
			Bloom: iblt.NewBloom(iblt.NewBloomConfig(8, 4)),
		},
	}

	printBlockInfo(respBlock)
}

func Test_printTransactionInfo(t *testing.T) {
	tx, _ := factory.CreateTestTx()
	respTx := rpc_interface.TransactionResp{
		Transaction: tx,
	}
	printTransactionInfo(respTx)
	printTransactionInfo(rpc_interface.TransactionResp{})
}

func TestRpcCaller_CurrentBlock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := getRpcTestApp()
	app.Action = func(c *cli.Context) {
		caller := &rpcCaller{}
		caller.CurrentBlock(c)
	}
	assert.NoError(t, app.Run([]string{os.Args[0]}))

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).Return(testErr)
		caller.CurrentBlock(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, method string, args ...interface{}) error {
			tx, _ := factory.CreateTestTx()
			m, _ := model.NewVoteMsgWithSign(uint64(1), uint64(1), common.HexToHash("0x1234"), model.PreVoteMessage, func(hash []byte) ([]byte, error) {
				return nil, nil
			}, fromAddr)
			*result.(*rpc_interface.BlockResp) = rpc_interface.BlockResp{
				Body: model.Body{Txs: []*model.Transaction{tx}, Vers: []model.AbstractVerification{m}},
				Header: model.Header{
					Bloom: iblt.NewBloom(iblt.NewBloomConfig(8, 4)),
				},
			}
			return nil
		})
		caller.CurrentBlock(c)
	}

	assert.NoError(t, app.Run([]string{os.Args[0], "CurrentBlock"}))
	client = nil
}

func TestRpcCaller_GetGenesis(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := getRpcTestApp()
	app.Action = func(c *cli.Context) {
		caller := &rpcCaller{}
		caller.GetGenesis(c)
	}
	assert.NoError(t, app.Run([]string{os.Args[0]}))

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).Return(testErr)
		caller.GetGenesis(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, method string, args ...interface{}) error {
			tx, _ := factory.CreateTestTx()
			m, _ := model.NewVoteMsgWithSign(uint64(1), uint64(1), common.HexToHash("0x1234"), model.PreVoteMessage, func(hash []byte) ([]byte, error) {
				return nil, nil
			}, fromAddr)
			*result.(*rpc_interface.BlockResp) = rpc_interface.BlockResp{
				Body: model.Body{Txs: []*model.Transaction{tx}, Vers: []model.AbstractVerification{m}},
				Header: model.Header{
					Bloom: iblt.NewBloom(iblt.NewBloomConfig(8, 4)),
				},
			}
			return nil
		})
		caller.GetGenesis(c)
	}

	assert.NoError(t, app.Run([]string{os.Args[0], "GetGenesis"}))
	client = nil
}

func TestRpcCaller_GetBlockByNumber(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := getRpcTestApp()
	app.Action = func(c *cli.Context) {
		caller := &rpcCaller{}
		caller.GetBlockByNumber(c)
	}
	assert.NoError(t, app.Run([]string{os.Args[0]}))

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}
		caller.GetBlockByNumber(c)

		c.Set("p", "")
		caller.GetBlockByNumber(c)

		c.Set("p", "s")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
		caller.GetBlockByNumber(c)

		c.Set("p", "1")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, args ...interface{}) error {
			tx, _ := factory.CreateTestTx()
			m, _ := model.NewVoteMsgWithSign(uint64(1), uint64(1), common.HexToHash("0x1234"), model.PreVoteMessage, func(hash []byte) ([]byte, error) {
				return nil, nil
			}, fromAddr)
			*result.(*rpc_interface.BlockResp) = rpc_interface.BlockResp{
				Body: model.Body{Txs: []*model.Transaction{tx}, Vers: []model.AbstractVerification{m}},
				Header: model.Header{
					Bloom: iblt.NewBloom(iblt.NewBloomConfig(8, 4)),
				},
			}
			return nil
		})
		caller.GetBlockByNumber(c)
	}

	assert.NoError(t, app.Run([]string{os.Args[0], "GetBlockByNumber"}))
	client = nil
}

func TestRpcCaller_GetBlockByHash(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := getRpcTestApp()
	app.Action = func(c *cli.Context) {
		caller := &rpcCaller{}
		caller.GetBlockByHash(c)
	}
	assert.NoError(t, app.Run([]string{os.Args[0]}))

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}
		caller.GetBlockByHash(c)

		c.Set("p", "")
		caller.GetBlockByHash(c)

		c.Set("p", "s")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
		caller.GetBlockByHash(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, args ...interface{}) error {
			tx, _ := factory.CreateTestTx()
			m, _ := model.NewVoteMsgWithSign(uint64(1), uint64(1), common.HexToHash("0x1234"), model.PreVoteMessage, func(hash []byte) ([]byte, error) {
				return nil, nil
			}, fromAddr)
			*result.(*rpc_interface.BlockResp) = rpc_interface.BlockResp{
				Body: model.Body{Txs: []*model.Transaction{tx}, Vers: []model.AbstractVerification{m}},
				Header: model.Header{
					Bloom: iblt.NewBloom(iblt.NewBloomConfig(8, 4)),
				},
			}
			return nil
		})
		caller.GetBlockByHash(c)
	}

	assert.NoError(t, app.Run([]string{os.Args[0], "GetBlockByHash"}))
	client = nil
}

func TestRpcCaller_StartMine(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := getRpcTestApp()
	app.Action = func(c *cli.Context) {
		SyncStatus.Store(true)
		caller := &rpcCaller{}
		caller.StartMine(c)
	}
	assert.NoError(t, app.Run([]string{os.Args[0]}))

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}
		SyncStatus.Store(false)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).Return(testErr)
		caller.StartMine(c)

		SyncStatus.Store(true)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).Return(testErr)
		caller.StartMine(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).Return(nil)
		caller.StartMine(c)
	}

	assert.NoError(t, app.Run([]string{os.Args[0], "StartMine"}))
	client = nil
}

func TestRpcCaller_StopMine(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := getRpcTestApp()
	app.Action = func(c *cli.Context) {
		caller := &rpcCaller{}
		caller.StopMine(c)
	}
	assert.NoError(t, app.Run([]string{os.Args[0]}))

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).Return(testErr)
		caller.StopMine(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).Return(nil)
		caller.StopMine(c)
	}

	assert.NoError(t, app.Run([]string{os.Args[0], "StopMine"}))
	client = nil
}

func TestRpcCaller_SetMineCoinBase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := getRpcTestApp()
	app.Action = func(c *cli.Context) {
		caller := &rpcCaller{}
		caller.SetMineCoinBase(c)
	}
	assert.NoError(t, app.Run([]string{os.Args[0]}))

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}

		c.Set("p", "")
		caller.SetMineCoinBase(c)

		c.Set("p", "coinBase")
		caller.SetMineCoinBase(c)

		c.Set("p", from)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
		caller.SetMineCoinBase(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		caller.SetMineCoinBase(c)
	}

	assert.NoError(t, app.Run([]string{os.Args[0], "SetMineCoinBase"}))
	client = nil
}

func TestRpcCaller_SetMineGasConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := getRpcTestApp()
	app.Action = func(c *cli.Context) {
		caller := &rpcCaller{}
		caller.SetMineGasConfig(c)
	}
	assert.NoError(t, app.Run([]string{os.Args[0]}))

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}

		c.Set("p", "")
		caller.SetMineGasConfig(c)

		c.Set("p", "gasFloor,gasCeil")
		caller.SetMineGasConfig(c)

		c.Set("p", fmt.Sprintf("%s,gasCeil", "10"))
		caller.SetMineGasConfig(c)

		c.Set("p", fmt.Sprintf("%s,%s", "10", "100"))
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
		caller.SetMineGasConfig(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		caller.SetMineGasConfig(c)
	}

	assert.NoError(t, app.Run([]string{os.Args[0], "SetMineGasConfig"}))
	client = nil
}

func TestRpcCaller_SendTx(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := getRpcTestApp()
	app.Action = func(c *cli.Context) {
		caller := &rpcCaller{}
		SyncStatus.Store(true)
		caller.SendTx(c)
	}
	assert.NoError(t, app.Run([]string{os.Args[0]}))

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}
		SyncStatus.Store(false)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).Return(nil)
		caller.SendTx(c)

		SyncStatus.Store(true)
		caller.SendTx(c)

		c.Set("p", "test")
		caller.SendTx(c)

		c.Set("p", "to,value,gasPrice,gasLimit,extraData")
		caller.SendTx(c)

		c.Set("p", fmt.Sprintf("%s,value,gasPrice,gasLimit,extraData", to))
		caller.SendTx(c)

		c.Set("p", fmt.Sprintf("%s,%s,gasPrice,gasLimit,extraData", to, "10dip"))
		caller.SendTx(c)

		c.Set("p", fmt.Sprintf("%s,%s,%s,gasLimit,extraData", to, "10dip", "10wu"))
		caller.SendTx(c)

		c.Set("p", fmt.Sprintf("%s,%s,%s,%s,extraData", to, "10dip", "10wu", "100"))
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
		caller.SendTx(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		caller.SendTx(c)
	}

	assert.NoError(t, app.Run([]string{os.Args[0], "SendTx"}))
	client = nil
}

func TestRpcCaller_SendTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := getRpcTestApp()
	app.Action = func(c *cli.Context) {
		caller := &rpcCaller{}
		SyncStatus.Store(true)
		caller.SendTransaction(c)
	}
	assert.NoError(t, app.Run([]string{os.Args[0]}))

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}
		SyncStatus.Store(false)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).Return(nil)
		caller.SendTransaction(c)

		SyncStatus.Store(true)
		caller.SendTransaction(c)

		c.Set("p", "test")
		caller.SendTransaction(c)

		c.Set("p", "from,to,value,gasPrice,gasLimit,extraData")
		caller.SendTransaction(c)

		c.Set("p", fmt.Sprintf("%s,to,value,gasPrice,gasLimit,extraData", from))
		caller.SendTransaction(c)

		c.Set("p", fmt.Sprintf("%s,%s,value,gasPrice,gasLimit,extraData", from, to))
		caller.SendTransaction(c)

		c.Set("p", fmt.Sprintf("%s,%s,%s,gasPrice,gasLimit,extraData", from, to, "10dip"))
		caller.SendTransaction(c)

		c.Set("p", fmt.Sprintf("%s,%s,%s,%s,gasLimit,extraData", from, to, "10dip", "10wu"))
		caller.SendTransaction(c)

		c.Set("p", fmt.Sprintf("%s,%s,%s,%s,%s,extraData", from, to, "10dip", "10wu", "100"))
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
		caller.SendTransaction(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		caller.SendTransaction(c)
	}

	assert.NoError(t, app.Run([]string{os.Args[0], "SendTransaction"}))
	client = nil
}

func Test_generateTxData(t *testing.T) {
	tx := model.CreateSignedTx(0, big.NewInt(100))
	transactionJson, err := json.Marshal(tx)
	assert.NoError(t, err)

	var expectTx *model.Transaction
	err = json.Unmarshal(transactionJson, &expectTx)
	assert.NoError(t, err)
	assert.Equal(t, expectTx.CalTxId(), tx.CalTxId())
}

func TestRpcCaller_Transaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := getRpcTestApp()
	app.Action = func(c *cli.Context) {
		caller := &rpcCaller{}
		SyncStatus.Store(true)
		caller.Transaction(c)
	}
	assert.NoError(t, app.Run([]string{os.Args[0]}))

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}
		SyncStatus.Store(false)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).Return(nil)
		caller.Transaction(c)

		SyncStatus.Store(true)
		caller.Transaction(c)

		c.Set("p", "")
		caller.Transaction(c)

		c.Set("p", "test")
		caller.Transaction(c)

		c.Set("p", "0x1234")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
		caller.Transaction(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, args ...interface{}) error {
			tx, _ := factory.CreateTestTx()
			*result.(*rpc_interface.TransactionResp) = rpc_interface.TransactionResp{
				Transaction: tx,
			}
			return nil
		})
		caller.Transaction(c)
	}

	assert.NoError(t, app.Run([]string{os.Args[0], "Transaction"}))
	client = nil
}

func TestRpcCaller_GetReceiptByTxHash(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := getRpcTestApp()
	app.Action = func(c *cli.Context) {
		caller := &rpcCaller{}
		caller.GetReceiptByTxHash(c)
	}
	assert.NoError(t, app.Run([]string{os.Args[0]}))

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}
		c.Set("p", "")
		caller.GetReceiptByTxHash(c)

		c.Set("p", "1234")
		caller.GetReceiptByTxHash(c)

		c.Set("p", txHash)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
		caller.GetReceiptByTxHash(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		caller.GetReceiptByTxHash(c)
	}

	assert.NoError(t, app.Run([]string{os.Args[0], "GetReceiptByTxHash"}))
	client = nil
}

func TestRpcCaller_GetReceiptsByBlockNum(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := getRpcTestApp()
	app.Action = func(c *cli.Context) {
		caller := &rpcCaller{}
		caller.GetReceiptsByBlockNum(c)
	}
	assert.NoError(t, app.Run([]string{os.Args[0]}))

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}
		c.Set("p", "")
		caller.GetReceiptsByBlockNum(c)

		c.Set("p", "aaa")
		caller.GetReceiptsByBlockNum(c)

		c.Set("p", "10")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
		caller.GetReceiptsByBlockNum(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		caller.GetReceiptsByBlockNum(c)
	}

	assert.NoError(t, app.Run([]string{os.Args[0], "GetReceiptsByBlockNum"}))
	client = nil
}

func TestRpcCaller_GetLogs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := getRpcTestApp()
	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}
		c.Set("p", "")
		caller.GetLogs(c)

		c.Set("p", "aaa")
		caller.GetLogs(c)

		jsonFile := `{
			"block_hash":"0x000023e18421a0abfceea172867b9b4a3bcf593edd0b504554bb7d1cf5f5e7b7",
			"addresses":["0x0014049F835be46352eD0Ec6B819272A2c8cF4feA10f"],
			"topics":[["0x0b5d2220daf8f0dfd95983d2ce625affbb7183c991271f49d818b4a64a268dbb"]]
		}`
		c.Set("p", jsonFile)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
		caller.GetLogs(c)

		jsonFile = `{
			"to_block":500,
			"addresses":["0x0014049F835be46352eD0Ec6B819272A2c8cF4feA10f"],
			"topics":[["0x0b5d2220daf8f0dfd95983d2ce625affbb7183c991271f49d818b4a64a268dbb"]]
		}`
		c.Set("p", jsonFile)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
		caller.GetLogs(c)

		log := model2.Log{}
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).SetArg(0, []model2.Log{log})
		caller.GetLogs(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		caller.GetLogs(c)
	}

	assert.NoError(t, app.Run([]string{os.Args[0], "GetLogs"}))
	client = nil
}

func TestRpcCaller_ListWallet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := getRpcTestApp()
	app.Action = func(c *cli.Context) {
		caller := &rpcCaller{}
		caller.ListWallet(c)
	}
	assert.NoError(t, app.Run([]string{os.Args[0]}))

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}
		c.Set("p", "")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).Return(testErr)
		caller.ListWallet(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, args ...interface{}) error {
			*result.(*[]accounts.WalletIdentifier) = []accounts.WalletIdentifier{
				{
					WalletType: 1,
					WalletName: "test",
					Path:       "",
				},
			}
			return nil
		})
		caller.ListWallet(c)
	}

	assert.NoError(t, app.Run([]string{os.Args[0], "ListWallet"}))
	client = nil
}

func TestRpcCaller_ListWalletAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := getRpcTestApp()
	app.Action = func(c *cli.Context) {
		caller := &rpcCaller{}
		caller.ListWalletAccount(c)
	}
	assert.NoError(t, app.Run([]string{os.Args[0]}))

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}
		c.Set("p", "test")
		caller.ListWalletAccount(c)

		c.Set("p", "")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
		caller.ListWalletAccount(c)

		c.Set("p", "Unknown, test")
		caller.ListWalletAccount(c)

		c.Set("p", "SoftWallet, test")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
		caller.ListWalletAccount(c)

		c.Set("p", "LedgerWallet, test")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
		caller.ListWalletAccount(c)

		c.Set("p", "TrezorWallet, test")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, args ...interface{}) error {
			*result.(*[]accounts.Account) = []accounts.Account{
				{
					Address: fromAddr,
				},
			}
			return nil
		})
		caller.ListWalletAccount(c)
	}

	assert.NoError(t, app.Run([]string{os.Args[0], "ListWalletAccount"}))
	client = nil
}

func TestRpcCaller_EstablishWallet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := getRpcTestApp()
	app.Action = func(c *cli.Context) {
		caller := &rpcCaller{}
		caller.EstablishWallet(c)
	}
	assert.NoError(t, app.Run([]string{os.Args[0]}))

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}
		c.Set("p", "")
		caller.EstablishWallet(c)

		c.Set("p", "test")

		caller.EstablishWallet(c)

		c.Set("p", "test,test,test")
		caller.EstablishWallet(c)

		c.Set("p", "SoftWallet,test,test")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
		caller.EstablishWallet(c)

		c.Set("p", "LedgerWallet,test,test")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
		caller.EstablishWallet(c)

		c.Set("p", "TrezorWallet,test,test")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, args ...interface{}) error {
			*result.(*string) = "test"
			return nil
		})
		caller.EstablishWallet(c)
	}

	assert.NoError(t, app.Run([]string{os.Args[0], "EstablishWallet"}))
	client = nil
}

func TestRpcCaller_RestoreWallet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := getRpcTestApp()
	app.Action = func(c *cli.Context) {
		caller := &rpcCaller{}
		caller.RestoreWallet(c)
	}
	assert.NoError(t, app.Run([]string{os.Args[0]}))

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}
		c.Set("p", "")
		caller.RestoreWallet(c)

		c.Set("p", "test")
		caller.RestoreWallet(c)

		c.Set("p", "test,test,test,test")
		caller.RestoreWallet(c)

		c.Set("p", "SoftWallet,test,test,test")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
		caller.RestoreWallet(c)

		c.Set("p", "LedgerWallet,test,test,test")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
		caller.RestoreWallet(c)

		c.Set("p", "TrezorWallet,test,test,test,test")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		caller.RestoreWallet(c)
	}

	assert.NoError(t, app.Run([]string{os.Args[0], "RestoreWallet"}))
	client = nil
}

func TestRpcCaller_OpenWallet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := getRpcTestApp()
	app.Action = func(c *cli.Context) {
		caller := &rpcCaller{}
		caller.OpenWallet(c)
	}
	assert.NoError(t, app.Run([]string{os.Args[0]}))

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}
		c.Set("p", "")
		caller.OpenWallet(c)

		c.Set("p", "test")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
		caller.OpenWallet(c)

		c.Set("p", "test,test,test")
		caller.OpenWallet(c)

		c.Set("p", "SoftWallet,test,test")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
		caller.OpenWallet(c)

		c.Set("p", "LedgerWallet,test,test")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
		caller.OpenWallet(c)

		c.Set("p", "TrezorWallet,test,test")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		caller.OpenWallet(c)
	}

	assert.NoError(t, app.Run([]string{os.Args[0], "OpenWallet"}))
	client = nil
}

func TestRpcCaller_CloseWallet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := getRpcTestApp()
	app.Action = func(c *cli.Context) {
		caller := &rpcCaller{}
		caller.CloseWallet(c)
	}
	assert.NoError(t, app.Run([]string{os.Args[0]}))

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}
		c.Set("p", "test")
		caller.CloseWallet(c)

		c.Set("p", "")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
		caller.CloseWallet(c)

		c.Set("p", "test,test")
		caller.CloseWallet(c)

		c.Set("p", "SoftWallet,test")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
		caller.CloseWallet(c)

		c.Set("p", "LedgerWallet,test")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
		caller.CloseWallet(c)

		c.Set("p", "TrezorWallet,test")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		caller.CloseWallet(c)
	}

	assert.NoError(t, app.Run([]string{os.Args[0], "CloseWallet"}))
	client = nil
}

func TestRpcCaller_AddAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := getRpcTestApp()
	app.Action = func(c *cli.Context) {
		caller := &rpcCaller{}
		caller.AddAccount(c)
	}
	assert.NoError(t, app.Run([]string{os.Args[0]}))

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}
		c.Set("p", "test")
		caller.AddAccount(c)

		c.Set("p", "")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
		caller.AddAccount(c)

		c.Set("p", "test,test")
		caller.AddAccount(c)

		c.Set("p", "SoftWallet,test")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
		caller.AddAccount(c)

		c.Set("p", "LedgerWallet,test")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
		caller.AddAccount(c)

		c.Set("p", "TrezorWallet,test")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, args ...interface{}) error {
			*result.(*accounts.Account) = accounts.Account{
				Address: fromAddr,
			}
			return nil
		})
		caller.AddAccount(c)
	}

	assert.NoError(t, app.Run([]string{os.Args[0], "AddAccount"}))
	client = nil
}

func TestRpcCaller_SendRegisterTx(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := getRpcTestApp()
	app.Action = func(c *cli.Context) {
		caller := &rpcCaller{}
		SyncStatus.Store(true)
		caller.SendRegisterTx(c)
	}
	assert.NoError(t, app.Run([]string{os.Args[0]}))

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}
		SyncStatus.Store(false)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).Return(testErr).Times(1)
		caller.SendRegisterTx(c)

		SyncStatus.Store(true)
		c.Set("p", "test")
		caller.SendRegisterTx(c)

		c.Set("p", "stake,gasPrice,gasLimit")
		caller.SendRegisterTx(c)

		c.Set("p", fmt.Sprintf("%s,gasPrice,gasLimit", "10dip"))
		caller.SendRegisterTx(c)

		c.Set("p", fmt.Sprintf("%s,%s,gasLimit", "10dip", "10wu"))
		caller.SendRegisterTx(c)

		c.Set("p", fmt.Sprintf("%s,%s,%s", "10dip", "10wu", "100"))
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
		caller.SendRegisterTx(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		caller.SendRegisterTx(c)
	}

	assert.NoError(t, app.Run([]string{os.Args[0], "SendRegisterTx"}))
	client = nil
}

func TestRpcCaller_SendRegisterTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := getRpcTestApp()
	app.Action = func(c *cli.Context) {
		caller := &rpcCaller{}
		SyncStatus.Store(true)
		caller.SendRegisterTransaction(c)
	}
	assert.NoError(t, app.Run([]string{os.Args[0]}))

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}
		SyncStatus.Store(false)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).Return(testErr).Times(1)
		caller.SendRegisterTransaction(c)

		SyncStatus.Store(true)
		c.Set("p", "test")
		caller.SendRegisterTransaction(c)

		c.Set("p", "from,stake,gasPrice,gasLimit")
		caller.SendRegisterTransaction(c)

		c.Set("p", fmt.Sprintf("%s,stake,gasPrice,gasLimit", from))
		caller.SendRegisterTransaction(c)

		c.Set("p", fmt.Sprintf("%s,%s,gasPrice,gasLimit", from, "10dip"))
		caller.SendRegisterTransaction(c)

		c.Set("p", fmt.Sprintf("%s,%s,%s,gasLimit", from, "10dip", "10wu"))
		caller.SendRegisterTransaction(c)

		c.Set("p", fmt.Sprintf("%s,%s,%s,%s", from, "10dip", "10wu", "100"))
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
		caller.SendRegisterTransaction(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		caller.SendRegisterTransaction(c)
	}

	assert.NoError(t, app.Run([]string{os.Args[0], "SendRegisterTransaction"}))
	client = nil
}

func TestRpcCaller_SendUnStakeTx(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := getRpcTestApp()
	app.Action = func(c *cli.Context) {
		caller := &rpcCaller{}
		SyncStatus.Store(true)
		caller.SendUnStakeTx(c)
	}
	assert.NoError(t, app.Run([]string{os.Args[0]}))

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}
		SyncStatus.Store(false)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).Return(testErr)
		caller.SendUnStakeTx(c)

		SyncStatus.Store(true)
		c.Set("p", "test")
		caller.SendUnStakeTx(c)

		c.Set("p", "gasPrice,gasLimit")
		caller.SendUnStakeTx(c)

		c.Set("p", fmt.Sprintf("%s,gasLimit", "10dip"))
		caller.SendUnStakeTx(c)

		c.Set("p", fmt.Sprintf("%s,%s", "10dip", "100"))
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
		caller.SendUnStakeTx(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		caller.SendUnStakeTx(c)
	}

	assert.NoError(t, app.Run([]string{os.Args[0], "SendUnStakeTx"}))
	client = nil
}

func TestRpcCaller_SendUnStakeTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := getRpcTestApp()
	app.Action = func(c *cli.Context) {
		caller := &rpcCaller{}
		SyncStatus.Store(true)
		caller.SendUnStakeTransaction(c)
	}
	assert.NoError(t, app.Run([]string{os.Args[0]}))

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}
		SyncStatus.Store(false)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).Return(testErr).Times(1)
		caller.SendUnStakeTransaction(c)

		SyncStatus.Store(true)
		c.Set("p", "test")
		caller.SendUnStakeTransaction(c)

		c.Set("p", "from,gasPrice,gasLimit")
		caller.SendUnStakeTransaction(c)

		c.Set("p", fmt.Sprintf("%s,gasPrice,gasLimit", from))
		caller.SendUnStakeTransaction(c)

		c.Set("p", fmt.Sprintf("%s,%s,gasLimit", from, "10dip"))
		caller.SendUnStakeTransaction(c)

		c.Set("p", fmt.Sprintf("%s,%s,%s", from, "10dip", "100"))
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
		caller.SendUnStakeTransaction(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		caller.SendUnStakeTransaction(c)
	}

	assert.NoError(t, app.Run([]string{os.Args[0], "SendUnStakeTransaction"}))
	client = nil
}

func TestRpcCaller_SendCancelTx(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := getRpcTestApp()
	app.Action = func(c *cli.Context) {
		caller := &rpcCaller{}
		SyncStatus.Store(true)
		caller.SendCancelTx(c)
	}
	assert.NoError(t, app.Run([]string{os.Args[0]}))

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}
		SyncStatus.Store(false)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).Return(testErr).Times(1)
		caller.SendCancelTx(c)

		SyncStatus.Store(true)
		c.Set("p", "test")
		caller.SendCancelTx(c)

		c.Set("p", "gasPrice,gasLimit")
		caller.SendCancelTx(c)

		c.Set("p", fmt.Sprintf("%s,gasLimit", "10dip"))
		caller.SendCancelTx(c)

		c.Set("p", fmt.Sprintf("%s,%s", "10dip", "100"))
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
		caller.SendCancelTx(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		caller.SendCancelTx(c)
	}

	assert.NoError(t, app.Run([]string{os.Args[0], "SendCancelTx"}))
	client = nil
}

func TestRpcCaller_SendCancelTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := getRpcTestApp()
	app.Action = func(c *cli.Context) {
		caller := &rpcCaller{}
		SyncStatus.Store(true)
		caller.SendCancelTransaction(c)
	}
	assert.NoError(t, app.Run([]string{os.Args[0]}))

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}
		SyncStatus.Store(false)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).Return(testErr).Times(1)
		caller.SendCancelTransaction(c)

		SyncStatus.Store(true)
		c.Set("p", "test")
		caller.SendCancelTransaction(c)

		c.Set("p", "from,gasPrice,gasLimit")
		caller.SendCancelTransaction(c)

		c.Set("p", fmt.Sprintf("%s,gasPrice,gasLimit", from))
		caller.SendCancelTransaction(c)

		c.Set("p", fmt.Sprintf("%s,%s,gasLimit", from, "10dip"))
		caller.SendCancelTransaction(c)

		c.Set("p", fmt.Sprintf("%s,%s,%s", from, "10dip", "100"))
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
		caller.SendCancelTransaction(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		caller.SendCancelTransaction(c)
	}

	assert.NoError(t, app.Run([]string{os.Args[0], "SendCancelTransaction"}))
	client = nil
}

func TestRpcCaller_GetVerifiersBySlot(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := getRpcTestApp()
	app.Action = func(c *cli.Context) {
		caller := &rpcCaller{}
		SyncStatus.Store(true)
		caller.GetVerifiersBySlot(c)
	}
	assert.NoError(t, app.Run([]string{os.Args[0]}))

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}
		SyncStatus.Store(false)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).Return(testErr).Times(1)
		caller.GetVerifiersBySlot(c)

		SyncStatus.Store(true)
		caller.GetVerifiersBySlot(c)

		c.Set("p", "")
		caller.GetVerifiersBySlot(c)

		c.Set("p", "test")
		caller.GetVerifiersBySlot(c)

		c.Set("p", "1")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
		caller.GetVerifiersBySlot(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, method string, args ...interface{}) error {
			*result.(*[]common.Address) = []common.Address{fromAddr}
			return nil
		})
		caller.GetVerifiersBySlot(c)
	}

	assert.NoError(t, app.Run([]string{os.Args[0], "GetVerifiersBySlot"}))
	client = nil
}

func TestRpcCaller_VerifierStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := getRpcTestApp()
	app.Action = func(c *cli.Context) {
		caller := &rpcCaller{}
		caller.VerifierStatus(c)
	}
	assert.NoError(t, app.Run([]string{os.Args[0]}))

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}
		c.Set("p", "test,test")
		caller.VerifierStatus(c)

		c.Set("p", "")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr).Times(2)
		caller.VerifierStatus(c)

		c.Set("p", "1234")
		caller.VerifierStatus(c)

		c.Set("p", from)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, method string, args ...interface{}) error {
			*result.(*rpc_interface.VerifierStatus) = rpc_interface.VerifierStatus{}
			return nil
		})
		caller.VerifierStatus(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, method string, args ...interface{}) error {
			*result.(*rpc_interface.VerifierStatus) = rpc_interface.VerifierStatus{
				Balance: (*hexutil.Big)(big.NewInt(1)),
				Status:  VerifierStatusNoRegistered,
			}
			return nil
		})
		caller.VerifierStatus(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, method string, args ...interface{}) error {
			*result.(*rpc_interface.VerifierStatus) = rpc_interface.VerifierStatus{
				Balance: (*hexutil.Big)(big.NewInt(1)),
				Status:  VerifiedStatusUnstaked,
			}
			return nil
		})
		caller.VerifierStatus(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, method string, args ...interface{}) error {
			*result.(*rpc_interface.VerifierStatus) = rpc_interface.VerifierStatus{
				Balance: (*hexutil.Big)(big.NewInt(1)),
				Status:  VerifierStatusRegistered,
			}
			return nil
		})
		caller.VerifierStatus(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, method string, args ...interface{}) error {
			*result.(*rpc_interface.VerifierStatus) = rpc_interface.VerifierStatus{
				Balance: (*hexutil.Big)(big.NewInt(1)),
				Stake:   (*hexutil.Big)(big.NewInt(1)),
				Status:  VerifierStatusRegistered,
			}
			return nil
		})
		caller.VerifierStatus(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, method string, args ...interface{}) error {
			*result.(*rpc_interface.VerifierStatus) = rpc_interface.VerifierStatus{
				Balance: (*hexutil.Big)(big.NewInt(1)),
				Stake:   (*hexutil.Big)(big.NewInt(1)),
				Status:  VerifiedStatusCanceled,
			}
			return nil
		})
		caller.VerifierStatus(c)
	}

	assert.NoError(t, app.Run([]string{os.Args[0], "VerifierStatus"}))
	client = nil
}

func TestRpcCaller_SetBftSigner(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := getRpcTestApp()
	app.Action = func(c *cli.Context) {
		caller := &rpcCaller{}
		caller.SetBftSigner(c)
	}
	assert.NoError(t, app.Run([]string{os.Args[0]}))

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}
		c.Set("p", "")
		caller.SetBftSigner(c)

		c.Set("p", "test")
		caller.SetBftSigner(c)

		c.Set("p", from)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		caller.SetBftSigner(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
		caller.SetBftSigner(c)
	}

	assert.NoError(t, app.Run([]string{os.Args[0], "SetBftSigner"}))
	client = nil
}

func TestRpcCaller_GetDefaultAccountStake(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := getRpcTestApp()
	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
		caller.GetDefaultAccountStake(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, args ...interface{}) error {
			*result.(*rpc_interface.CurBalanceResp) = rpc_interface.CurBalanceResp{
				Balance: (*hexutil.Big)(big.NewInt(1)),
			}
			return nil
		})
		caller.GetDefaultAccountStake(c)
	}

	assert.NoError(t, app.Run([]string{os.Args[0], "GetDefaultAccountStake"}))
	client = nil
}

func TestRpcCaller_CurrentStake(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := getRpcTestApp()
	app.Action = func(c *cli.Context) {
		caller := &rpcCaller{}
		caller.CurrentStake(c)
	}
	assert.NoError(t, app.Run([]string{os.Args[0]}))

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}
		c.Set("p", "test,test")
		caller.CurrentStake(c)

		c.Set("p", "")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr).Times(2)
		caller.CurrentStake(c)

		c.Set("p", "1234")
		caller.CurrentStake(c)

		c.Set("p", from)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, method string, args ...interface{}) error {
			*result.(*rpc_interface.CurBalanceResp) = rpc_interface.CurBalanceResp{}
			return nil
		})
		caller.CurrentStake(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, method string, args ...interface{}) error {
			*result.(*rpc_interface.CurBalanceResp) = rpc_interface.CurBalanceResp{
				Balance: (*hexutil.Big)(big.NewInt(1)),
			}
			return nil
		})
		caller.CurrentStake(c)
	}

	assert.NoError(t, app.Run([]string{os.Args[0], "CurrentStake"}))
	client = nil
}

func TestRpcCaller_CurrentReputation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := getRpcTestApp()
	app.Action = func(c *cli.Context) {
		caller := &rpcCaller{}
		caller.CurrentReputation(c)
	}
	assert.NoError(t, app.Run([]string{os.Args[0]}))

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}
		c.Set("p", "test,test")
		caller.CurrentReputation(c)

		c.Set("p", "")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr).Times(2)
		caller.CurrentReputation(c)

		c.Set("p", "1234")
		caller.CurrentReputation(c)

		c.Set("p", from)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		caller.CurrentReputation(c)
	}

	assert.NoError(t, app.Run([]string{os.Args[0], "CurrentReputation"}))
	client = nil
}

func TestRpcCaller_GetCurVerifiers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := getRpcTestApp()
	app.Action = func(c *cli.Context) {
		caller := &rpcCaller{}
		SyncStatus.Store(true)
		caller.GetCurVerifiers(c)
	}
	assert.NoError(t, app.Run([]string{os.Args[0]}))

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}
		SyncStatus.Store(false)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
		caller.GetCurVerifiers(c)

		SyncStatus.Store(true)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
		caller.GetCurVerifiers(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, method string, args ...interface{}) error {
			*result.(*[]common.Address) = []common.Address{fromAddr}
			return nil
		})
		caller.GetCurVerifiers(c)
	}

	assert.NoError(t, app.Run([]string{os.Args[0], "GetCurVerifiers"}))
	client = nil
}

func Test_inDefaultVs(t *testing.T) {
	address := fromAddr
	assert.Equal(t, inDefaultVs(address), false)
	address = chain_config.LocalVerifierAddress[0]
	assert.Equal(t, inDefaultVs(address), true)
}

func TestRpcCaller_GetNextVerifiers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := getRpcTestApp()
	app.Action = func(c *cli.Context) {
		caller := &rpcCaller{}
		SyncStatus.Store(true)
		caller.GetNextVerifiers(c)
	}
	assert.NoError(t, app.Run([]string{os.Args[0]}))

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}
		SyncStatus.Store(false)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
		caller.GetNextVerifiers(c)

		SyncStatus.Store(true)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
		caller.GetNextVerifiers(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, method string, args ...interface{}) error {
			*result.(*[]common.Address) = []common.Address{fromAddr}
			return nil
		})
		caller.GetNextVerifiers(c)
	}

	assert.NoError(t, app.Run([]string{os.Args[0], "GetNextVerifiers"}))
	client = nil
}

func Test_getNonceInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := getRpcTestApp()
	app.Action = func(c *cli.Context) {
		nonce, err := getNonceInfo(c)
		assert.Error(t, err)
		assert.Equal(t, nonce, uint64(0))
	}
	assert.NoError(t, app.Run([]string{os.Args[0]}))

	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		c.Set("p", "test,test")
		nonce, err := getNonceInfo(c)
		assert.Error(t, err)
		assert.Equal(t, nonce, uint64(0))

		c.Set("p", "")
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr).Times(2)
		nonce, err = getNonceInfo(c)
		assert.Error(t, err)
		assert.Equal(t, nonce, uint64(0))

		c.Set("p", "test")
		nonce, err = getNonceInfo(c)
		assert.Error(t, err)
		assert.Equal(t, nonce, uint64(0))

		c.Set("p", from)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		nonce, err = getNonceInfo(c)
		assert.NoError(t, err)
		assert.Equal(t, nonce, uint64(0))
	}

	assert.NoError(t, app.Run([]string{os.Args[0], "getNonceInfo"}))
	client = nil
}

func TestRpcCaller_GetTransactionNonce(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := getRpcTestApp()
	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}
		SyncStatus.Store(false)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
		caller.GetTransactionNonce(c)

		SyncStatus.Store(true)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr).Times(2)
		caller.GetTransactionNonce(c)

		c.Set("p", from)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
		caller.GetTransactionNonce(c)
	}

	assert.NoError(t, app.Run([]string{os.Args[0], "GetTransactionNonce"}))
	client = nil
}

func TestRpcCaller_GetAddressNonceFromWallet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := getRpcTestApp()
	app.Action = func(c *cli.Context) {
		client = NewMockRpcClient(ctrl)
		caller := &rpcCaller{}
		c.Set("p", from)
		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		caller.GetAddressNonceFromWallet(c)

		client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
		caller.GetAddressNonceFromWallet(c)
	}

	assert.NoError(t, app.Run([]string{os.Args[0], "GetAddressNonceFromWallet"}))
	client = nil
}

func testTempJSONFile(t *testing.T) (string, func()) {
	t.Helper()
	tf, err := ioutil.TempFile("", "*.json")
	if err != nil {
		t.Fatalf(":err: %s", err)
	}

	tf.Close()
	return tf.Name(), func() { os.Remove(tf.Name()) }
}

func Test_initWallet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	tf, tfClean := testTempJSONFile(t)
	defer tfClean()
	client = NewMockRpcClient(ctrl)

	client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
	err := initWallet(tf, "test", "test")
	assert.Error(t, err)

	client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	err = initWallet("test", "test", "test")

	client = nil
}

func Test_getDefaultAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client = NewMockRpcClient(ctrl)
	client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, method string, args ...interface{}) error {
		*result.(*[]accounts.WalletIdentifier) = []accounts.WalletIdentifier{{}}
		return nil
	})

	client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)
	address := getDefaultAccount()
	assert.Equal(t, address, common.Address{})
}

func Test_getDefaultWallet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client = NewMockRpcClient(ctrl)
	client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, method string, args ...interface{}) error {
		*result.(*[]accounts.WalletIdentifier) = []accounts.WalletIdentifier{{}}
		return nil
	})

	wallet := getDefaultWallet()
	assert.Equal(t, wallet, accounts.WalletIdentifier{})
}
