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

package rpc_interface

import (
	"context"
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/g-error"
	"github.com/dipperin/dipperin-core/core/accounts"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/core/chain/state-processor"
	"github.com/dipperin/dipperin-core/core/contract"
	"github.com/dipperin/dipperin-core/core/dipperin/service"
	"github.com/dipperin/dipperin-core/core/economy-model"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/tests/g-mockFile"
	"github.com/dipperin/dipperin-core/tests/g-testData"
	"github.com/dipperin/dipperin-core/third-party/p2p/enode"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestDipperinMercuryApi(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	mp := NewMockPeerManager(controller)
	mc := g_mockFile.NewMockChainInterface(controller)
	mb := NewMockAbstractBlock(controller)
	mn := NewMockNodeConf(controller)
	mpeer := NewMockPmAbstractPeer(controller)
	mpm := NewMockAbstractPbftProtocolManager(controller)
	api := &DipperinVenusApi{service: &service.VenusFullChainService{
		DipperinConfig: &service.DipperinConfig{
			NormalPm:           mp,
			ChainReader:        mc,
			TxPool:             &fakeTxPool{},
			Broadcaster:        &fakeBroadcaster{},
			NodeConf:           mn,
			WalletManager:      &accounts.WalletManager{},
			MsgSigner:          &fakeMsgSigner{},
			PriorityCalculator: model.DefaultPriorityCalculator,
			PbftPm:             mpm,
		},
		TxValidator: &fakeTxV{},
	}}

	mp.EXPECT().IsSync().Return(true)
	assert.True(t, api.GetSyncStatus())

	// get block info
	mc.EXPECT().CurrentBlock().Return(mb).AnyTimes()
	mb.EXPECT().Header().Return(&model.Header{}).AnyTimes()
	mb.EXPECT().Body().Return(&model.Body{}).AnyTimes()
	_, err := api.CurrentBlock()
	assert.NoError(t, err)
	mc.EXPECT().GetBlockByNumber(gomock.Any()).Return(mb).AnyTimes()
	_, err = api.GetBlockByNumber(0)
	assert.NoError(t, err)
	mc.EXPECT().GetBlockByHash(gomock.Any()).Return(mb).AnyTimes()
	_, err = api.GetBlockByHash(common.Hash{})
	assert.NoError(t, err)
	xp := uint64(1)
	mc.EXPECT().GetBlockNumber(common.Hash{}).Return(&xp)
	x := api.GetBlockNumber(common.Hash{})
	assert.NotNil(t, x)

	// get genesis
	mc.EXPECT().Genesis().Return(mb)
	_, err = api.GetGenesis()
	assert.NoError(t, err)

	// get block body
	mc.EXPECT().GetBody(gomock.Any()).Return(&model.Body{})
	b := api.GetBlockBody(common.Hash{})
	assert.NotNil(t, b)
	adb, _ := NewEmptyAccountDB()
	mc.EXPECT().CurrentState().Return(adb, nil).AnyTimes()
	_, err = api.CurrentBalance(common.Address{})
	assert.NoError(t, err)

	// get transactions
	gasPrice := g_testData.TestGasPrice
	gasLimit := g_testData.TestGasLimit
	mc.EXPECT().GetTransaction(common.Hash{}).Return(&model.Transaction{}, common.Hash{}, uint64(1), uint64(0)).AnyTimes()
	_, err = api.Transaction(common.Hash{})
	assert.NoError(t, err)
	_, err = api.GetTransactionNonce(common.Address{})
	assert.Error(t, err)
	_, err = api.NewTransaction([]byte{})
	assert.Error(t, err)
	tx := model.NewTransaction(0, common.Address{}, big.NewInt(1), gasPrice, gasLimit, nil)
	tb, err := rlp.EncodeToBytes(tx)
	assert.NoError(t, err)
	_, err = api.NewTransaction(tb)
	assert.NoError(t, err)

	// node type
	mn.EXPECT().GetNodeType().Return(0).AnyTimes()
	assert.Error(t, api.SetMineCoinBase(common.Address{}))
	assert.Error(t, api.StartMine())
	assert.Error(t, api.StopMine())

	// soft wallet info
	mn.EXPECT().SoftWalletName().Return("test_wa").AnyTimes()
	mn.EXPECT().SoftWalletDir().Return("/tmp/test_wa").AnyTimes()
	_, err = api.EstablishWallet("", "", accounts.WalletIdentifier{})
	assert.Error(t, err)
	assert.Error(t, api.OpenWallet("", accounts.WalletIdentifier{}))
	assert.Error(t, api.CloseWallet(accounts.WalletIdentifier{}))
	assert.Error(t, api.RestoreWallet("", "", "", accounts.WalletIdentifier{}))
	_, err = api.ListWallet()
	assert.NoError(t, err)
	assert.NotEmpty(t, BuildContractExtraData("x", common.Address{}, ""))

	// get ERC20 info
	mc.EXPECT().CurrentHeader().Return(&model.Header{}).AnyTimes()
	_, err = api.GetContractInfo(&contract.ExtraDataForContract{})
	assert.Error(t, err)
	_, err = api.GetContract(common.Address{})
	assert.Error(t, err)
	_, err = api.ERC20TotalSupply(common.Address{})
	assert.Error(t, err)
	_, err = api.ERC20Balance(common.Address{}, common.Address{})
	assert.Error(t, err)
	_, err = api.ERC20Allowance(common.Address{}, common.Address{}, common.Address{})
	assert.Error(t, err)
	_, err = api.ERC20Transfer(common.Address{}, common.Address{}, common.Address{}, big.NewInt(1), g_testData.TestGasPrice, g_testData.TestGasLimit)
	assert.Error(t, err)
	_, err = api.ERC20TransferFrom(common.Address{}, common.Address{}, common.Address{}, common.Address{}, big.NewInt(1), g_testData.TestGasPrice, g_testData.TestGasLimit)
	assert.Error(t, err)
	_, err = api.ERC20Approve(common.Address{}, common.Address{}, common.Address{}, big.NewInt(1), g_testData.TestGasPrice, g_testData.TestGasLimit)
	assert.Error(t, err)
	_, err = api.CreateERC20(common.Address{}, "", "", big.NewInt(1), 2, g_testData.TestGasPrice, g_testData.TestGasLimit)
	assert.Error(t, err)

	// list wallet account
	n, _ := enode.ParseV4(fmt.Sprintf("enode://b832f4f2fe19dbc5604766bbb268a6d0f7ce9ce381b034b262a92f0ad8283a1b5fa058dea5269b66fbb2014a24fa7198c6dc2d8c9cbac7a348258fc20702561f@%v:%v", "127.0.0.1", 10003))
	chain_config.KBucketNodes = []*enode.Node{n}
	_, err = api.CheckBootNode()
	assert.NoError(t, err)
	_, err = api.ListWalletAccount(accounts.WalletIdentifier{})
	assert.Error(t, err)
	assert.Error(t, api.SetBftSigner(common.Address{}))
	_, err = api.AddAccount("", accounts.WalletIdentifier{})
	assert.Error(t, err)

	// send transaction
	nonce := uint64(1)
	value := big.NewInt(1)
	mp.EXPECT().BestPeer().Return(mpeer).AnyTimes()
	mpeer.EXPECT().GetHead().Return(common.Hash{}, uint64(1)).AnyTimes()
	api.RemoteHeight()
	_, err = api.SendTransaction(common.Address{}, common.Address{}, value, gasPrice, gasLimit, nil, &nonce)
	assert.Error(t, err)
	_, err = api.SendTransactions(common.Address{}, []model.RpcTransaction{})
	assert.Error(t, err)
	_, err = api.NewSendTransactions([]model.Transaction{})
	assert.NoError(t, err)
	_, err = api.SendRegisterTransaction(common.Address{}, value, gasPrice, gasLimit, &nonce)
	assert.Error(t, err)
	_, err = api.SendUnStakeTransaction(common.Address{}, gasPrice, gasLimit, &nonce)
	assert.Error(t, err)
	_, err = api.SendEvidenceTransaction(common.Address{}, common.Address{}, gasPrice, gasLimit, nil, nil, &nonce)
	assert.Error(t, err)
	_, err = api.SendCancelTransaction(common.Address{}, gasPrice, gasLimit, &nonce)
	assert.Error(t, err)
	_, err = api.SendTransactionContract(common.Address{}, common.Address{}, value, gasPrice, gasLimit, nil, &nonce)
	assert.Error(t, err)

	// get verifiers
	req := uint64(1)
	mc.EXPECT().GetVerifiers(gomock.Any()).Return([]common.Address{{}}).AnyTimes()
	mc.EXPECT().GetCurrVerifiers().Return([]common.Address{{}}).AnyTimes()
	mc.EXPECT().GetNextVerifiers().Return([]common.Address{{}}).AnyTimes()
	mc.EXPECT().GetSlot(mb).Return(&req).AnyTimes()
	mb.EXPECT().Number().Return(uint64(1)).AnyTimes()
	_, err = api.GetVerifiersBySlot(1)
	assert.NoError(t, err)
	slot, err := api.GetSlotByNumber(mb.Number())
	assert.Equal(t, req, slot)
	vs := api.GetCurVerifiers()
	assert.Len(t, vs, 1)
	vs = api.GetNextVerifiers()
	assert.Len(t, vs, 1)
	_, err = api.VerifierStatus(common.Address{})
	assert.NoError(t, err)
	_, err = api.CurrentStake(common.Address{})
	assert.NoError(t, err)
	_, err = api.CurrentReputation(common.Address{})
	assert.Error(t, err)

	// get connect peers
	mpm.EXPECT().GetCurrentConnectPeers().Return(map[string]common.Address{"x": {}})
	ps, err := api.GetCurrentConnectPeers()
	assert.NoError(t, err)
	assert.Len(t, ps, 1)
	_, err = api.GetAddressNonceFromWallet(common.Address{})
	assert.Error(t, err)
	_, err = api.GetChainConfig()
	assert.NoError(t, err)
	_, err = api.GetBlockDiffVerifierInfo(1)
	assert.Error(t, err)

	// reward verifier and miner
	mc.EXPECT().GetEconomyModel().Return(economy_model.MakeDipperinEconomyModel(nil, economy_model.DIPProportion)).AnyTimes()
	_, err = api.GetVerifierDIPReward(1)
	assert.NoError(t, err)
	_, err = api.GetMineMasterDIPReward(1)
	assert.NoError(t, err)
	_, err = api.GetBlockYear(1)
	assert.NoError(t, err)
	_, err = api.GetOneBlockTotalDIPReward(1)
	assert.NoError(t, err)
	api.GetInvestorInfo()
	api.GetDeveloperInfo()
	_, err = api.GetInvestorLockDIP(common.Address{}, 1)
	assert.Error(t, err)
	_, err = api.GetDeveloperLockDIP(common.Address{}, 1)
	assert.Error(t, err)
	api.GetFoundationInfo(0)
	_, err = api.GetMaintenanceLockDIP(common.Address{}, 1)
	assert.NoError(t, err)
	_, err = api.GetReMainRewardLockDIP(common.Address{}, 1)
	assert.NoError(t, err)
	_, err = api.GetEarlyTokenLockDIP(common.Address{}, 1)
	assert.NoError(t, err)
	_, err = api.GetMineMasterEDIPReward(1, 1)
	assert.NoError(t, err)
	_, err = api.GetVerifierEDIPReward(1, 1)
	assert.NoError(t, err)

	// subscribe block
	ctx := context.Background()
	_, err = api.NewBlock(ctx)
	assert.Error(t, err)
	_, err = api.SubscribeBlock(ctx)
	assert.Error(t, err)
	//api.StopDipperin()

	// get abi and logs
	WASMPath := g_testData.GetWASMPath("event", g_testData.CoreVmTestData)
	AbiPath := g_testData.GetAbiPath("event", g_testData.CoreVmTestData)
	code, abi := g_testData.GetCodeAbi(WASMPath, AbiPath)
	assert.NoError(t, adb.NewAccountState(common.Address{}))
	assert.NoError(t, adb.SetAbi(common.Address{}, abi))
	assert.NoError(t, adb.SetCode(common.Address{}, code))
	mb.EXPECT().StateRoot().Return(common.Hash{}).AnyTimes()
	mc.EXPECT().AccountStateDB(common.Hash{}).Return(adb, nil).AnyTimes()
	mc.EXPECT().GetLatestNormalBlock().Return(mb).AnyTimes()
	_, err = api.GetABI(common.Address{})
	assert.NoError(t, err)
	_, err = api.GetCode(common.Address{})
	assert.NoError(t, err)
	_, err = api.GetLogs(common.Hash{}, uint64(100), uint64(0), nil, nil)
	assert.Equal(t, g_error.BeginNumLargerError, err)

	// suggest gasPrice
	mb.EXPECT().Hash().Return(common.Hash{}).AnyTimes()
	_, err = api.SuggestGasPrice()
	assert.NoError(t, err)

	// get receipts
	mc.EXPECT().GetReceipts(common.Hash{}, uint64(1)).Return(nil).AnyTimes()
	_, err = api.GetContractAddressByTxHash(common.Hash{})
	assert.Equal(t, g_error.ErrReceiptIsNil, err)
	_, err = api.GetTxActualFee(common.Hash{})
	assert.Equal(t, g_error.ErrReceiptIsNil, err)
	_, err = api.GetReceiptByTxHash(common.Hash{})
	assert.Equal(t, g_error.ErrReceiptIsNil, err)
	_, err = api.GetReceiptsByBlockNum(uint64(1))
	assert.Equal(t, g_error.ErrReceiptIsNil, err)

	// new contract and estimate gas
	_, err = api.NewContract([]byte{}, uint64(0))
	assert.Error(t, err)
	_, err = api.NewContract(tb, uint64(0))
	assert.Equal(t, g_error.ErrInvalidContractType, err)
	_, err = api.NewEstimateGas([]byte{})
	assert.Error(t, err)
	_, err = api.NewEstimateGas(tb)
	assert.NoError(t, err)

	// call contract and estimate gas
	mc.EXPECT().StateAtByBlockNumber(uint64(1)).Return(adb, nil).AnyTimes()
	_, err = api.CallContract(common.Address{}, common.Address{}, nil, uint64(0))
	assert.Equal(t, g_error.ErrEmptyTxData, err)
	data, err := rlp.EncodeToBytes([]interface{}{code, abi})
	assert.NoError(t, err)
	_, err = api.CallContract(common.Address{}, common.Address{}, data, uint64(0))
	assert.Equal(t, accounts.ErrNotFindWallet, err)
	_, err = api.CallContract(common.Address{}, common.Address{}, data, uint64(1))
	assert.Equal(t, accounts.ErrNotFindWallet, err)
	_, err = api.EstimateGas(common.Address{}, common.Address{}, nil, nil, gasLimit, nil, &nonce)
	assert.Equal(t, g_error.ErrEmptyTxData, err)
	_, err = api.EstimateGas(common.Address{}, common.Address{}, nil, nil, gasLimit, data, &nonce)
	assert.Equal(t, accounts.ErrNotFindWallet, err)
}

func NewEmptyAccountDB() (*state_processor.AccountStateDB, state_processor.StateStorage) {
	storage := state_processor.NewStateStorageWithCache(ethdb.NewMemDatabase())
	db, err := state_processor.NewAccountStateDB(common.Hash{}, storage)
	if err != nil {
		panic(err)
	}
	return db, storage
}

type fakeTxV struct{}

func (f *fakeTxV) Valid(tx model.AbstractTransaction) error {
	return nil
}

type fakeTxPool struct{}

func (p *fakeTxPool) Stats() (int, int) {
	panic("implement me")
}

func (p *fakeTxPool) AddRemotes(txs []model.AbstractTransaction) []error {
	return nil
}

func (p *fakeTxPool) AddLocals(txs []model.AbstractTransaction) []error {
	return nil
}

func (p *fakeTxPool) AddRemote(tx model.AbstractTransaction) error {
	return nil
}

type fakeBroadcaster struct{}

func (f *fakeBroadcaster) BroadcastTx(txs []model.AbstractTransaction) {}

type fakeMsgSigner struct{}

func (f *fakeMsgSigner) SetBaseAddress(address common.Address) {}

func (f *fakeMsgSigner) GetAddress() common.Address {
	return common.Address{}
}
