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
	"errors"
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/address-util"
	"github.com/dipperin/dipperin-core/common/hexutil"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/core/economy-model"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"math/big"
	"reflect"
	"testing"
)

func TestSaveContractId(t *testing.T) {

	SaveContractId("", "default_m0", erc20Addr)
	addr, err := GetContractId("", "default_m0")
	fmt.Println(addr, err)

	SaveContractId(util.HomeDir(), "default_m0", erc20Addr)
	addr, err = GetContractId(util.HomeDir(), "default_m0")
	assert.Error(t, err, errors.New("non-exsisted contract record"))
}

func Test_newInfoOfContract(t *testing.T) {
	info := newInfoOfContract(BuiltInERC20Token{})
	assert.NotNil(t, info)
}

func TestParseExtraDataForContract(t *testing.T) {
	dataStruct := &ExtraDataForContract{Action: "test"}
	data := util.StringifyJsonToBytes(dataStruct)
	result := ParseExtraDataForContract(data)
	assert.NotNil(t, result)
	assert.Equal(t, dataStruct.Action, result.Action)

	result = ParseExtraDataForContract([]byte("test"))
	assert.Nil(t, result)
}

func TestGetContractTempByType(t *testing.T) {
	ret, err := GetContractTempByType("test")
	assert.Error(t, err)
	assert.Nil(t, ret)

	ret, err = GetContractTempByType("ERC20")
	assert.NoError(t, err)
}

func TestGetContractMethodArgs(t *testing.T) {
	ret, err := GetContractMethodArgs("test", "")
	assert.Error(t, err)
	assert.Nil(t, ret)

	ret, err = GetContractMethodArgs("ERC20", "nomtd")
	assert.Error(t, err)
	assert.Nil(t, ret)

	ret, err = GetContractMethodArgs("ERC20", "TransferFrom")
	assert.NoError(t, err)
	assert.NotNil(t, ret)
}

func TestProcessor_Process(t *testing.T) {
	mockCtl := gomock.NewController(t)
	mockContactDB := NewMockContractDB(mockCtl)
	mockAccountDB := NewMockAccountDB(mockCtl)
	returnValue := reflect.New(reflect.TypeOf(BuiltInERC20Token{}))
	mockContactDB.EXPECT().ContractExist(gomock.Any()).Return(false).AnyTimes()
	mockContactDB.EXPECT().GetContract(gomock.Any(), gomock.Any()).Return(returnValue, nil).AnyTimes()
	mockContactDB.EXPECT().PutContract(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	processor := NewProcessor(mockContactDB, uint64(1))
	processor.SetAccountDB(mockAccountDB)

	tx := CreateSignedTx(erc20Addr, []byte("test"))
	err := processor.Process(tx)
	assert.Error(t, err, CanNotParseContractErr)

	// create contract
	m := ExtraDataForContract{
		ContractAddress: erc20Addr,
		Action:          "create",
		Params:          createERC20JsonWithOwner,
	}
	extraData := util.StringifyJsonToBytes(m)
	tx = CreateSignedTx(erc20Addr, extraData)
	err = processor.Process(tx)
	assert.NoError(t, err)

	// call contract
	m.Action = "Decimals"
	m.Params = "[]"
	extraData = util.StringifyJsonToBytes(m)
	tx = model.NewTransaction(0, common.Address{}, big.NewInt(0), big.NewInt(0), 0, extraData)
	err = processor.Process(tx)
	assert.Error(t, err)

	tx = CreateSignedTx(erc20Addr, extraData)
	err = processor.Process(tx)
	assert.NoError(t, err)
}

func TestProcessor_DoCreate(t *testing.T) {
	mockCtl := gomock.NewController(t)
	mockContactDB := NewMockContractDB(mockCtl)
	mockAccountDB := NewMockAccountDB(mockCtl)
	processor := &Processor{contractDB: mockContactDB, accountDB: mockAccountDB}

	exData := &ExtraDataForContract{}
	_, err := processor.DoCreate(exData)
	assert.Error(t, err, ContractAdrEmptyErr)

	mockContactDB.EXPECT().ContractExist(gomock.Any()).Return(true)
	exData.ContractAddress = common.HexToAddress("1234")
	_, err = processor.DoCreate(exData)
	assert.Error(t, err)

	mockContactDB.EXPECT().ContractExist(gomock.Any()).Return(false).AnyTimes()
	_, err = processor.DoCreate(exData)
	assert.Error(t, err)

	exData.ContractAddress = erc20Addr
	_, err = processor.DoCreate(exData)
	assert.Error(t, err)

	exData.Params = createERC20Json
	_, err = processor.DoCreate(exData)
	assert.Error(t, err)

	exData.Params = createERC20JsonWithOwner
	_, err = processor.DoCreate(exData)
	assert.NoError(t, err)
}

func TestProcessor_Run(t *testing.T) {
	mockCtl := gomock.NewController(t)
	mockContactDB := NewMockContractDB(mockCtl)
	mockAccountDB := NewMockAccountDB(mockCtl)
	processor := &Processor{contractDB: mockContactDB, accountDB: mockAccountDB}

	exData := &ExtraDataForContract{}
	exData.ContractAddress = common.HexToAddress("1234")
	_, err := processor.Run(aliceAddr, exData)
	assert.Error(t, err)

	exData.ContractAddress = erc20Addr
	mockContactDB.EXPECT().GetContract(gomock.Any(), gomock.Any()).Return(reflect.Value{}, errors.New("test error"))
	_, err = processor.Run(aliceAddr, exData)
	assert.Error(t, err)

	returnValue := reflect.New(reflect.TypeOf(BuiltInERC20Token{}))
	mockContactDB.EXPECT().GetContract(gomock.Any(), gomock.Any()).Return(returnValue, nil).AnyTimes()
	_, err = processor.Run(aliceAddr, exData)
	assert.Error(t, err)

	exData.Action = "Decimals"
	_, err = processor.Run(aliceAddr, exData)
	assert.Error(t, err)

	exData.Params = "[]"
	_, err = processor.Run(aliceAddr, exData)
	assert.NoError(t, err)
}

func TestProcessor_GetContractReadOnlyInfo(t *testing.T) {
	mockCtl := gomock.NewController(t)
	mockContactDB := NewMockContractDB(mockCtl)
	mockAccountDB := NewMockAccountDB(mockCtl)
	processor := &Processor{contractDB: mockContactDB, accountDB: mockAccountDB}

	exData := &ExtraDataForContract{}
	exData.ContractAddress = common.HexToAddress("1234")
	_, err := processor.GetContractReadOnlyInfo(exData)
	assert.Error(t, err)

	exData.ContractAddress = erc20Addr
	mockContactDB.EXPECT().GetContract(gomock.Any(), gomock.Any()).Return(reflect.Value{}, errors.New("test error"))
	_, err = processor.GetContractReadOnlyInfo(exData)
	assert.Error(t, err)

	returnValue := reflect.New(reflect.TypeOf(BuiltInERC20Token{}))
	mockContactDB.EXPECT().GetContract(gomock.Any(), gomock.Any()).Return(returnValue, nil).AnyTimes()
	_, err = processor.GetContractReadOnlyInfo(exData)
	assert.Error(t, err)

	exData.Action = "Decimals"
	_, err = processor.GetContractReadOnlyInfo(exData)
	assert.Error(t, err)

	exData.Params = "[]"
	_, err = processor.GetContractReadOnlyInfo(exData)
	assert.NoError(t, err)
}

func TestGetERC20TxSize(t *testing.T) {
	testContract := newTestToken()

	es := util.StringifyJson(testContract)
	extra := ExtraDataForContract{}
	extra.Action = "create"
	extra.Params = es
	contractAddr, _ := address_util.GenERC20Address()
	extra.ContractAddress = contractAddr

	log.Info("the es is:", "es", es)
	log.Info("the contractAddr is:", "contractAddr", contractAddr.Hex())

	extraData := []byte(util.StringifyJson(extra))
	log.Info("the extraData is:", "extraData", hexutil.Encode(extraData))
	log.Info("the extraData Len is:", "len", len(extraData))

	tx := CreateSignedTx(contractAddr, extraData)
	log.Info("the tx size is:", "size", tx.Size(), "txFee", economy_model.GetMinimumTxFee(tx.Size()))
	log.Info("the tx hash is:", "hash", tx.CalTxId().Hex())

	normalTx := CreateSignedTx(bobAddr, []byte{})
	log.Info("the normal tx size is:", "size", normalTx.Size(), "txFee", economy_model.GetMinimumTxFee(normalTx.Size()))
}
