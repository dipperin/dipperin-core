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
	"math/big"
	"testing"
	"time"

	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/consts"
	"github.com/dipperin/dipperin-core/common/hexutil"
	"github.com/dipperin/dipperin-core/core/accounts"
	"github.com/dipperin/dipperin-core/core/bloom"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/core/rpc-interface"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_loadDefaultAccountStake(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client = NewMockRpcClient(ctrl)

	client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("test"))

	loadDefaultAccountStake()

	client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, method string, args ...interface{}) error {
		*result.(*rpc_interface.CurBalanceResp) = rpc_interface.CurBalanceResp{}
		return nil
	})

	loadDefaultAccountStake()

	client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, method string, args ...interface{}) error {
		*result.(*rpc_interface.CurBalanceResp) = rpc_interface.CurBalanceResp{
			Balance: (*hexutil.Big)(big.NewInt(0)),
		}
		return nil
	})

	loadDefaultAccountStake()

	assert.Equal(t, defaultAccountStake, "0"+consts.CoinDIPName)

	client = nil
}

func TestPrintDefaultAccountStake(t *testing.T) {
	PrintDefaultAccountStake()
}

func TestAsyncLogElectionTx(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client = NewMockRpcClient(ctrl)

	logElectionTxTickerTime = 1 * time.Millisecond
	SyncStatus.Store(true)

	client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).Return(errors.New("test")).AnyTimes()
	timer := AsyncLogElectionTx()
	time.Sleep(2 * time.Millisecond)
	timer.Stop()
	time.Sleep(2 * time.Millisecond)
	client = nil
}

func Test_loadRegistedAccounts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client = NewMockRpcClient(ctrl)

	client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).Return(errors.New("test"))
	loadRegistedAccounts()

	client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, method string, args ...interface{}) error {
		*result.(*[]accounts.WalletIdentifier) = []accounts.WalletIdentifier{
			{},
		}
		return nil
	})
	client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("test"))
	loadRegistedAccounts()

	client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, method string, args ...interface{}) error {
		*result.(*[]accounts.WalletIdentifier) = []accounts.WalletIdentifier{
			{},
		}
		return nil
	})
	client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, method string, args ...interface{}) error {
		*result.(*[]accounts.Account) = []accounts.Account{
			{
				Address: common.HexToAddress("0x1234"),
			},
		}
		return nil
	})
	client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("test"))
	loadRegistedAccounts()

	client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, method string, args ...interface{}) error {
		*result.(*[]accounts.WalletIdentifier) = []accounts.WalletIdentifier{
			{},
		}
		return nil
	})
	client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, method string, args ...interface{}) error {
		*result.(*[]accounts.Account) = []accounts.Account{
			{
				Address: common.HexToAddress("0x1234"),
			},
		}
		return nil
	})
	client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, method string, args ...interface{}) error {
		*result.(*rpc_interface.VerifierStatus) = rpc_interface.VerifierStatus{
			Status: "Registered",
		}
		return nil
	})
	loadRegistedAccounts()

	client = nil
}

func Test_logElection(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client = NewMockRpcClient(ctrl)
	trackingAccounts = []accounts.Account{}
	logElection()

	trackingAccounts = append(trackingAccounts, accounts.Account{Address: common.HexToAddress("0x1234")})
	client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).Return(errors.New("test"))
	logElection()

	client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, method string, args ...interface{}) error {
		*result.(*rpc_interface.BlockResp) = rpc_interface.BlockResp{
			Header: model.Header{
				Bloom: iblt.NewBloom(iblt.NewBloomConfig(8, 4)),
			},
		}
		return nil
	})
	client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("test"))
	logElection()

	client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, method string, args ...interface{}) error {
		*result.(*rpc_interface.BlockResp) = rpc_interface.BlockResp{
			Header: model.Header{
				Bloom: iblt.NewBloom(iblt.NewBloomConfig(8, 4)),
			},
		}
		return nil
	})
	client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("test"))
	client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("test"))
	logElection()

	client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, method string, args ...interface{}) error {
		*result.(*rpc_interface.BlockResp) = rpc_interface.BlockResp{
			Header: model.Header{
				Bloom: iblt.NewBloom(iblt.NewBloomConfig(8, 4)),
			},
		}
		return nil
	})
	client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("test"))
	client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, method string, args ...interface{}) error {
		*result.(*rpc_interface.VerifierStatus) = rpc_interface.VerifierStatus{
			Status: VerifierStatusNoRegistered,
		}
		return nil
	})
	logElection()

	client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, method string, args ...interface{}) error {
		*result.(*rpc_interface.BlockResp) = rpc_interface.BlockResp{
			Header: model.Header{
				Bloom: iblt.NewBloom(iblt.NewBloomConfig(8, 4)),
			},
		}
		return nil
	})
	client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("test"))
	client.(*MockRpcClient).EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(result interface{}, method string, args ...interface{}) error {
		*result.(*rpc_interface.VerifierStatus) = rpc_interface.VerifierStatus{
			Status: VerifierStatusRegistered,
		}
		return nil
	})
	logElection()

	time.Sleep(1 * time.Millisecond)

	client = nil
}

func Test_isVerifier(t *testing.T) {
	address1 := common.HexToAddress("0x1234")
	address2 := common.HexToAddress("0x1233")

	verifierAddresses := []common.Address{address1}

	assert.Equal(t, isVerifier(address2, verifierAddresses), false)
	assert.Equal(t, isVerifier(address1, verifierAddresses), true)
}

func Test_addTrackingAccount(t *testing.T) {
	trackingAccounts = []accounts.Account{}
	addTrackingAccount(common.HexToAddress("0x1234"))
	addTrackingAccount(common.HexToAddress("0x1234"))

	assert.Equal(t, len(trackingAccounts), 1)
}

func Test_removeTrackingAccount(t *testing.T) {
	trackingAccounts = []accounts.Account{}
	addTrackingAccount(common.HexToAddress("0x1234"))
	addTrackingAccount(common.HexToAddress("0x1233"))
	removeTrackingAccount(common.HexToAddress("0x1234"))
	assert.Equal(t, len(trackingAccounts), 1)
}
