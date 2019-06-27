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
	"github.com/dipperin/dipperin-core/common/address-util"
	"github.com/dipperin/dipperin-core/common/g-testData"
	"github.com/dipperin/dipperin-core/common/hexutil"
	"github.com/dipperin/dipperin-core/core/economy-model"
	"github.com/dipperin/dipperin-core/third-party/log"
	"testing"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/core/model"
	"fmt"
	"github.com/stretchr/testify/assert"
	"errors"
	"math/big"
	"github.com/golang/mock/gomock"
	"reflect"
)

func TestGetContractConf(t *testing.T) {
	//c := GetContractConf()
	//fmt.Println(util.StringifyJson(c))
}

func TestDecodeClientExtraData(t *testing.T) {
	//data := `0x7b22616374696f6e223a22637265617465222c22636f6e74726163745f61646472657373223a2230783030303164353031306132326231423742656465343237323130376231386136314335453065626537344638222c22706172616d73223a227b5c226f776e65725c223a5c22307830303030614164376632303362354331444230453537434442633438303866386431393538323130653138365c222c5c22746f6b656e5f6e616d655c223a5c2254657374315c222c5c22746f6b656e5f646563696d616c735c223a31382c5c22746f6b656e5f73796d626f6c5c223a5c2254315c222c5c22746f6b656e5f746f74616c5f737570706c795c223a5c22307866343234305c222c5c2262616c616e6365735c223a7b5c22307830303030614164376632303362354331444230453537434442633438303866386431393538323130653138365c223a5c22307866343234305c227d2c5c22616c6c6f7765645c223a7b7d7d227`
	//jb, err := hexutil.Decode(data)
	//fmt.Println("json str", string(jb))
	//assert.NoError(t, err)
	//var eData ExtraDataForContract
	//err = util.ParseJsonFromBytes(jb, &eData)
	//assert.NoError(t, err)
	//fmt.Println(eData)
}

func TestSaveContractId(t *testing.T) {
	SaveContractId("", "default_m0", common.HexToAddress("0x00000000000000000000000000000101231234124124"))
	addr, err := GetContractId("", "default_m0")
	//assert.NoError(t, err)
	fmt.Println(addr, err)
	//assert.Equal(t, "0x00000000000000000000000000000101231234124124", addr[len(addr)-1].Hex())

	SaveContractId(util.HomeDir(), "default_m0", common.HexToAddress("0x00000000000000000000000000000101231234124124"))
	addr, err = GetContractId(util.HomeDir(), "default_m0")
	assert.Error(t, err, errors.New("non-exsisted contract record"))
}

func Test_newInfoOfContract(t *testing.T) {
	info := newInfoOfContract(BuiltInERC20Token{})
	assert.NotNil(t, info)
}

func TestParseExtraDataForContract(t *testing.T) {
	dataStruct := &ExtraDataForContract{Action:"test"}
	data := util.StringifyJsonToBytes(dataStruct)
	result := ParseExtraDataForContract(data)
	assert.NotNil(t, result)
	assert.Equal(t, dataStruct.Action, result.Action)

	result = ParseExtraDataForContract([]byte("test"))
	assert.Nil(t, result)
}

func TestGetContractTempByType(t *testing.T) {
	ret, err := GetContractTempByType("test")
	fmt.Println(err)
	assert.Nil(t, ret)

	ret, err = GetContractTempByType("ERC20")
	assert.NoError(t, err)
}

func TestGetContractMethodArgs(t *testing.T) {
	ret, _ := GetContractMethodArgs("test", "")
	assert.Nil(t, ret)
	ret, _ = GetContractMethodArgs("ERC20", "nomtd")
	assert.Nil(t, ret)

	ret, err := GetContractMethodArgs("ERC20", "TransferFrom")
	fmt.Println(ret)
	assert.NoError(t, err)
}

func TestProcessor_Process(t *testing.T) {
	processor := &Processor{}
	tx := model.NewTransaction(0, common.HexToAddress("1234"), big.NewInt(10), g_testData.TestGasPrice,g_testData.TestGasLimit, []byte("test"))
	err := processor.Process(tx)
	assert.Error(t, err, CanNotParseContractErr)
}

func TestProcessor_DoCreate(t *testing.T) {
	mockCtl := gomock.NewController(t)
	mockConractDB := NewMockContractDB(mockCtl)
	mockAccountDB := NewMockAccountDB(mockCtl)
	processor := &Processor{contractDB:mockConractDB, accountDB:mockAccountDB}

	exData := &ExtraDataForContract{}
	_, err := processor.DoCreate(exData)
	assert.Error(t, err, ContractAdrEmptyErr)

	mockConractDB.EXPECT().ContractExist(gomock.Any()).Return(true)
	exData.ContractAddress = common.HexToAddress("1234")
	_, err = processor.DoCreate(exData)
	assert.NotNil(t, err)

	mockConractDB.EXPECT().ContractExist(gomock.Any()).Return(false)
	exData.ContractAddress = common.HexToAddress("1234")
	_, err = processor.DoCreate(exData)
	assert.NotNil(t, err)

	mockConractDB.EXPECT().ContractExist(gomock.Any()).Return(false)
	exData.ContractAddress = common.HexToAddress("0x00100000FA42f7315cD04D6774E58B54e92603e96d84")
	_, err = processor.DoCreate(exData)
	assert.NotNil(t, err)

	createERC20ConfigJsonStr := `{"token_name":"EOS","token_decimals":18,"token_symbol":"EOS","token_total_supply":"0x33b2e3c9fd0803ce8000000","balances":{},"allowed":{}}`
	mockConractDB.EXPECT().ContractExist(gomock.Any()).Return(false)
	exData.ContractAddress = common.HexToAddress("0x00100000FA42f7315cD04D6774E58B54e92603e96d84")
	exData.Params = createERC20ConfigJsonStr
	_, err = processor.DoCreate(exData)
	assert.NotNil(t, err)

	createERC20ConfigJsonStr = `{"owner":"0x1234","token_name":"EOS","token_decimals":18,"token_symbol":"EOS","token_total_supply":"0x33b2e3c9fd0803ce8000000","balances":{},"allowed":{}}`
	mockConractDB.EXPECT().ContractExist(gomock.Any()).Return(false)
	exData.ContractAddress = common.HexToAddress("0x00100000FA42f7315cD04D6774E58B54e92603e96d84")
	exData.Params = createERC20ConfigJsonStr
	_, err = processor.DoCreate(exData)
	assert.Nil(t, err)
}

func TestProcessor_Run(t *testing.T) {
	mockCtl := gomock.NewController(t)
	mockConractDB := NewMockContractDB(mockCtl)
	mockAccountDB := NewMockAccountDB(mockCtl)
	processor := &Processor{contractDB:mockConractDB, accountDB:mockAccountDB}

	exData := &ExtraDataForContract{}
	exData.ContractAddress = common.HexToAddress("1234")
	_, err := processor.Run(common.HexToAddress("5678"), exData)
	assert.NotNil(t, err)

	exData.ContractAddress = common.HexToAddress("0x00100000FA42f7315cD04D6774E58B54e92603e96d84")
	mockConractDB.EXPECT().GetContract(gomock.Any(), gomock.Any()).Return(reflect.Value{}, errors.New("test error"))
	_, err = processor.Run(common.HexToAddress("5678"), exData)
	assert.NotNil(t, err)

	tContract := BuiltInERC20Token{}
	mockConractDB.EXPECT().GetContract(gomock.Any(), gomock.Any()).Return(reflect.New(reflect.TypeOf(tContract)), nil)
	_, err = processor.Run(common.HexToAddress("5678"), exData)
	assert.NotNil(t, err)

	mockConractDB.EXPECT().GetContract(gomock.Any(), gomock.Any()).Return(reflect.New(reflect.TypeOf(tContract)), nil)
	exData.Action = "Decimals"
	_, err = processor.Run(common.HexToAddress("5678"), exData)
	assert.NotNil(t, err)

	mockConractDB.EXPECT().GetContract(gomock.Any(), gomock.Any()).Return(reflect.New(reflect.TypeOf(tContract)), nil)
	exData.Params = "[]"
	_, err = processor.Run(common.HexToAddress("5678"), exData)
	assert.NoError(t, err)
}

func TestProcessor_GetContractReadOnlyInfo(t *testing.T) {
	mockCtl := gomock.NewController(t)
	mockConractDB := NewMockContractDB(mockCtl)
	mockAccountDB := NewMockAccountDB(mockCtl)
	processor := &Processor{contractDB:mockConractDB, accountDB:mockAccountDB}

	exData := &ExtraDataForContract{}
	exData.ContractAddress = common.HexToAddress("1234")
	_, err := processor.GetContractReadOnlyInfo(exData)
	assert.NotNil(t, err)

	exData.ContractAddress = common.HexToAddress("0x00100000FA42f7315cD04D6774E58B54e92603e96d84")
	mockConractDB.EXPECT().GetContract(gomock.Any(), gomock.Any()).Return(reflect.Value{}, errors.New("test error"))
	_, err = processor.GetContractReadOnlyInfo(exData)
	assert.NotNil(t, err)

	tContract := BuiltInERC20Token{}
	mockConractDB.EXPECT().GetContract(gomock.Any(), gomock.Any()).Return(reflect.New(reflect.TypeOf(tContract)), nil)
	_, err = processor.GetContractReadOnlyInfo(exData)
	assert.NotNil(t, err)

	mockConractDB.EXPECT().GetContract(gomock.Any(), gomock.Any()).Return(reflect.New(reflect.TypeOf(tContract)), nil)
	exData.Action = "Decimals"
	_, err = processor.GetContractReadOnlyInfo(exData)
	assert.NotNil(t, err)

	mockConractDB.EXPECT().GetContract(gomock.Any(), gomock.Any()).Return(reflect.New(reflect.TypeOf(tContract)), nil)
	exData.Params = "[]"
	_, err = processor.GetContractReadOnlyInfo(exData)
	assert.NoError(t, err)
}

func TestGetERC20TxSize(t *testing.T){
	testContract := newTestToken()

	es := util.StringifyJson(testContract)
	extra := ExtraDataForContract{}
	extra.Action = "create"
	extra.Params = es
	contractAdr, _ := address_util.GenERC20Address()
	extra.ContractAddress = contractAdr

	log.Info("the es is:","es",es)
	log.Info("the contractAdr is:","contractAdr",contractAdr.Hex())

	extraData := []byte(util.StringifyJson(extra))
	log.Info("the extraData is:","extraData",hexutil.Encode(extraData))
	log.Info("the extraData Len is:","len",len(extraData))

	tx := model.NewTransaction(0,contractAdr,big.NewInt(0), g_testData.TestGasPrice,g_testData.TestGasLimit,extraData)
	key1, _ := model.CreateKey()
	fs := model.NewMercurySigner(big.NewInt(1))
	tx.SignTx(key1,fs)

	tx.RawSignatureValues()
	log.Info("the tx size is:","size",tx.Size(),"txFee",economy_model.GetMinimumTxFee(tx.Size()))
	log.Info("the tx hash is:","hash",tx.CalTxId().Hex())

	normalTx := model.NewTransaction(0,common.HexToAddress("0x00009865E43BEebad5fB771259F1660cD2aC4fD82557"),big.NewInt(10),g_testData.TestGasPrice,g_testData.TestGasLimit,[]byte{})
	normalTx.SignTx(key1,fs)


	log.Info("the tx is:","tx",normalTx)
	log.Info("the normal tx size is:","size",normalTx.Size(),"txFee",economy_model.GetMinimumTxFee(normalTx.Size()))
}