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
package stateprocessor

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/addressutil"
	"github.com/dipperin/dipperin-core/common/log"
	"github.com/dipperin/dipperin-core/core/accounts/accountsbase"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/core/vm/base/utils"
	"github.com/dipperin/dipperin-core/third_party/crypto"
	"github.com/dipperin/dipperin-core/third_party/crypto/cs-crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"math/big"
	"testing"
)

func TestAccountStateDB_ProcessContract_Error(t *testing.T) {
	to := common.HexToAddress(common.AddressContractCreate)

	type result struct {
		err error
	}

	testCases := []struct {
		name   string
		given  func() error
		expect result
	}{
		{
			name:"cannot encode negative *big.Int error",
			given: func() error {
				tx := model.NewTransaction(uint64(0), to, big.NewInt(100), testGasPrice, testGasLimit, nil)
				block := CreateBlock(1, common.Hash{}, []*model.Transaction{tx}, 5*testGasLimit)
				tmpGasLimit := block.GasLimit()
				gasUsed := block.GasUsed()
				config := &TxProcessConfig{
					Tx:       tx,
					Header:   block.Header(),
					GetHash:  getTestHashFunc(),
					GasLimit: &tmpGasLimit,
					GasUsed:  &gasUsed,
					TxFee:    big.NewInt(0),
				}

				db, root := CreateTestStateDB()
				tdb := NewStateStorageWithCache(db)
				processor, err := NewAccountStateDB(root, tdb)
				err = processor.ProcessTxNew(config)
				return err
			},
			expect:result{errors.New("rlp: cannot encode negative *big.Int")},
		},
		{
			name:"gas limit reached error",
			given: func() error {
				db, root := CreateTestStateDB()
				tdb := NewStateStorageWithCache(db)
				processor, _ := NewAccountStateDB(root, tdb)
				WASMPath := model.GetWASMPath("event", model.CoreVmTestData)
				abiPath := model.GetAbiPath("event", model.CoreVmTestData)
				tx := createContractTx(WASMPath, abiPath, 0, testGasLimit)
				block := CreateBlock(1, common.Hash{}, []*model.Transaction{tx}, uint64(100))
				tmpGasLimit := block.GasLimit()
				gasUsed := block.GasUsed()
				config := &TxProcessConfig{
					Tx:       tx,
					Header:   block.Header(),
					GetHash:  getTestHashFunc(),
					GasLimit: &tmpGasLimit,
					GasUsed:  &gasUsed,
					TxFee:    big.NewInt(0),
				}
				err := processor.ProcessTxNew(config)
				return err
			},
			expect:result{errors.New("gas limit reached")},
		},
	}

	for _,tc := range testCases{
		err:=tc.given()
		if err!=nil{
			assert.Equal(t,tc.expect.err,err)
		}
	}
}

func TestAccountStateDB_ProcessContractToken(t *testing.T) {
	singer := model.NewSigner(new(big.Int).SetInt64(int64(1)))

	ownSK, _ := crypto.GenerateKey()
	ownPk := ownSK.PublicKey
	ownAddress := cs_crypto.GetNormalAddress(ownPk)

	aliceSK, _ := crypto.GenerateKey()
	alicePk := aliceSK.PublicKey
	aliceAddress := cs_crypto.GetNormalAddress(alicePk)

	brotherSK, _ := crypto.GenerateKey()
	brotherPk := brotherSK.PublicKey
	brotherAddress := cs_crypto.GetNormalAddress(brotherPk)

	addressSlice := []common.Address{
		ownAddress,
		aliceAddress,
		brotherAddress,
	}

	WASMPath := model.GetWASMPath("token", model.CoreVmTestData)
	abiPath := model.GetAbiPath("token", model.CoreVmTestData)
	input := []string{"dipp", "DIPP", "1000000"}
	data, err := getCreateExtraData(WASMPath, abiPath, input)
	assert.NoError(t, err)

	addr := common.HexToAddress(common.AddressContractCreate)
	tx := model.NewTransaction(0, addr, big.NewInt(10), big.NewInt(1), 26427000, data)
	signCreateTx := getSignedTx(t, ownSK, tx, singer)

	gasLimit := testGasLimit * 10000000000
	block := CreateBlock(1, common.Hash{}, []*model.Transaction{signCreateTx}, gasLimit)
	processor, err := CreateProcessorAndInitAccount(t, addressSlice)

	tmpGasLimit := block.GasLimit()
	gasUsed := block.GasUsed()
	config := &TxProcessConfig{
		Tx:       tx,
		Header:   block.Header(),
		GetHash:  getTestHashFunc(),
		GasLimit: &tmpGasLimit,
		GasUsed:  &gasUsed,
		TxFee:    big.NewInt(0),
	}

	err = processor.ProcessTxNew(config)
	assert.NoError(t, err)

	contractAddr := cs_crypto.CreateContractAddress(ownAddress, uint64(0))
	_, err = processor.GetNonce(contractAddr)
	_, err = processor.GetCode(contractAddr)
	abi, err := processor.GetAbi(contractAddr)
	assert.NoError(t, err)
	processor.Commit()

	accountOwn := accountsbase.Account{ownAddress}
	//  合约调用getBalance方法  获取合约原始账户balance
	ownTransferNonce, err := processor.GetNonce(ownAddress)
	assert.NoError(t, err)
	err = processContractCall(t, contractAddr, abi, ownSK, processor, accountOwn, ownTransferNonce, "getBalance", ownAddress.Hex(), 2, singer)
	assert.NoError(t, err)

	gasUsed2 := uint64(0)
	//  合约调用  transfer方法 转账给alice
	ownTransferNonce++
	err = processContractCall(t, contractAddr, abi, ownSK, processor, accountOwn, ownTransferNonce, "transfer", aliceAddress.Hex()+",20", 3, singer)
	assert.NoError(t, err)

	//  合约调用getBalance方法  获取alice账户balance
	ownTransferNonce++
	err = processContractCall(t, contractAddr, abi, ownSK, processor, accountOwn, ownTransferNonce, "getBalance", aliceAddress.Hex(), 4, singer)
	assert.NoError(t, err)

	//  合约调用approve方法
	ownTransferNonce++
	callTxApprove, err := newContractCallTx(nil, &contractAddr, new(big.Int).SetUint64(1), uint64(1500000), "approve", brotherAddress.Hex()+",50", ownTransferNonce, abi)
	signCallTxApprove, err := callTxApprove.SignTx(ownSK, singer)

	assert.NoError(t, err)
	block5 := CreateBlock(5, common.Hash{}, []*model.Transaction{signCallTxApprove}, gasLimit)

	txConfig5 := &TxProcessConfig{
		Tx:       signCallTxApprove,
		Header:   block5.Header().(*model.Header),
		GetHash:  fakeGetBlockHash,
		GasLimit: &gasLimit,
		GasUsed:  &gasUsed2,
	}

	err = processor.ProcessTxNew(txConfig5)
	assert.NoError(t, err)
	processor.Commit()

	//  合约调用transferFrom方法
	callTxTransferFrom, err := newContractCallTx(nil, &contractAddr, new(big.Int).SetUint64(1), uint64(1500000), "transferFrom", ownAddress.Hex()+","+aliceAddress.Hex()+",5", 0, abi)
	assert.NoError(t, err)
	accountBrother := accountsbase.Account{Address: brotherAddress}
	assert.NoError(t, err)

	signCallTxTransferFrom, err := callTxTransferFrom.SignTx(brotherSK, singer)
	assert.NoError(t, err)
	block7 := CreateBlock(7, common.Hash{}, []*model.Transaction{signCallTxTransferFrom}, gasLimit)

	txConfig7 := &TxProcessConfig{
		Tx:       signCallTxTransferFrom,
		Header:   block7.Header().(*model.Header),
		GetHash:  fakeGetBlockHash,
		GasLimit: &gasLimit,
		GasUsed:  &gasUsed2,
	}

	err = processor.ProcessTxNew(txConfig7)
	assert.NoError(t, err)
	processor.Commit()

	//  合约调用getBalance方法  获取alice账户获得转账授权后的balance
	ownTransferNonce++
	err = processContractCall(t, contractAddr, abi, ownSK, processor, accountOwn, ownTransferNonce, "getBalance", aliceAddress.Hex(), 8, singer)
	assert.NoError(t, err)

	//  合约调用getBalance方法  获取own账户最终的balance
	ownTransferNonce++
	err = processContractCall(t, contractAddr, abi, ownSK, processor, accountOwn, ownTransferNonce, "getBalance", ownAddress.Hex(), 9, singer)
	assert.NoError(t, err)

	// 合约调用  transfer方法  转账给brother
	ownTransferNonce++
	err = processContractCall(t, contractAddr, abi, ownSK, processor, accountOwn, ownTransferNonce, "transfer", brotherAddress.Hex()+",28", 10, singer)
	assert.NoError(t, err)

	//  合约调用getBalance方法  获取own账户最终的balance
	ownTransferNonce++
	err = processContractCall(t, contractAddr, abi, ownSK, processor, accountOwn, ownTransferNonce, "getBalance", ownAddress.Hex(), 11, singer)
	assert.NoError(t, err)

	// 合约调用burn方法,将账户余额返还给own
	err = processContractCall(t, contractAddr, abi, brotherSK, processor, accountBrother, 1, "burn", "15", 12, singer)
	assert.NoError(t, err)

	// 合约调用getBalance方法,获取own的余额
	ownTransferNonce++
	err = processContractCall(t, contractAddr, abi, ownSK, processor, accountOwn, ownTransferNonce, "getBalance", ownAddress.Hex(), 13, singer)
	assert.NoError(t, err)

	// 合约调用setName方法，设置合约名
	ownTransferNonce++
	err = processContractCall(t, contractAddr, abi, ownSK, processor, accountOwn, ownTransferNonce, "setName", "wujinhai", 14, singer)
	assert.NoError(t, err)
}

func TestContractNewFeature(t *testing.T) {
	singer := model.NewSigner(new(big.Int).SetInt64(int64(1)))
	ownSK, _ := crypto.GenerateKey()
	ownPk := ownSK.PublicKey
	ownAddress := cs_crypto.GetNormalAddress(ownPk)

	aliceSK, _ := crypto.GenerateKey()
	alicePk := aliceSK.PublicKey
	aliceAddress := cs_crypto.GetNormalAddress(alicePk)

	addressSlice := []common.Address{
		ownAddress,
		aliceAddress,
	}

	WASMPath := model.GetWASMPath("demo", model.CoreVmTestData)
	abiPath := model.GetAbiPath("demo", model.CoreVmTestData)
	input := []string{"0x0000D36F282D8925B16Ed24CB637475e6a03B01E1056"}
	//0x0000D36F282D8925B16Ed24CB637475e6a03B01E1056
	data, err := getCreateExtraData(WASMPath, abiPath, input)
	assert.NoError(t, err)

	addr := common.HexToAddress(common.AddressContractCreate)
	tx := model.NewTransaction(0, addr, big.NewInt(0), big.NewInt(1), 26427000, data)
	signCreateTx := getSignedTx(t, ownSK, tx, singer)

	gasLimit := testGasLimit * 10000000000
	block := CreateBlock(1, common.Hash{}, []*model.Transaction{signCreateTx}, gasLimit)
	processor, err := CreateProcessorAndInitAccount(t, addressSlice)

	tmpGasLimit := block.GasLimit()
	gasUsed := block.GasUsed()
	config := &TxProcessConfig{
		Tx:       tx,
		Header:   block.Header(),
		GetHash:  getTestHashFunc(),
		GasLimit: &tmpGasLimit,
		GasUsed:  &gasUsed,
		TxFee:    big.NewInt(0),
	}

	err = processor.ProcessTxNew(config)
	assert.NoError(t, err)
}

func TestContractWithNewFeature(t *testing.T) {
	singer := model.NewSigner(new(big.Int).SetInt64(int64(1)))

	ownSK, _ := crypto.GenerateKey()
	ownPk := ownSK.PublicKey
	ownAddress := cs_crypto.GetNormalAddress(ownPk)

	aliceSK, _ := crypto.GenerateKey()
	alicePk := aliceSK.PublicKey
	aliceAddress := cs_crypto.GetNormalAddress(alicePk)

	brotherSK, _ := crypto.GenerateKey()
	brotherPk := brotherSK.PublicKey
	brotherAddress := cs_crypto.GetNormalAddress(brotherPk)

	addressSlice := []common.Address{
		ownAddress,
		aliceAddress,
		brotherAddress,
	}

	//WASMPath := g_testData.GetWASMPath("token", g_testData.CoreVmTestData)
	//abiPath := g_testData.GetAbiPath("token", g_testData.CoreVmTestData)
	WASMPath := model.GetWASMPath("token-param", model.CoreVmTestData)
	abiPath := model.GetAbiPath("token-param", model.CoreVmTestData)
	input := []string{"dipp", "DIPP", "1000000"}
	data, err := getCreateExtraData(WASMPath, abiPath, input)
	assert.NoError(t, err)

	addr := common.HexToAddress(common.AddressContractCreate)
	tx := model.NewTransaction(0, addr, big.NewInt(0), big.NewInt(1), 26427000, data)
	signCreateTx := getSignedTx(t, ownSK, tx, singer)

	gasLimit := testGasLimit * 10000000000
	block := CreateBlock(1, common.Hash{}, []*model.Transaction{signCreateTx}, gasLimit)
	processor, err := CreateProcessorAndInitAccount(t, addressSlice)

	tmpGasLimit := block.GasLimit()
	gasUsed := block.GasUsed()
	config := &TxProcessConfig{
		Tx:       tx,
		Header:   block.Header(),
		GetHash:  getTestHashFunc(),
		GasLimit: &tmpGasLimit,
		GasUsed:  &gasUsed,
		TxFee:    big.NewInt(0),
	}

	err = processor.ProcessTxNew(config)
	assert.NoError(t, err)

	contractAddr := cs_crypto.CreateContractAddress(ownAddress, uint64(0))
	_, err = processor.GetNonce(contractAddr)
	_, err = processor.GetCode(contractAddr)
	abi, err := processor.GetAbi(contractAddr)
	assert.NoError(t, err)
	//assert.Equal(t, code, tx.ExtraData())
	processor.Commit()

	accountOwn := accountsbase.Account{ownAddress}
	//  合约调用getBalance方法  获取合约原始账户balance
	ownTransferNonce, err := processor.GetNonce(ownAddress)
	assert.NoError(t, err)
	err = processContractCall(t, contractAddr, abi, ownSK, processor, accountOwn, ownTransferNonce, "getBalance", ownAddress.Hex(), 2, singer)
	assert.NoError(t, err)

	//gasUsed2 := uint64(0)
	//  合约调用  transfer方法 转账给alice
	ownTransferNonce++
	err = processContractCall(t, contractAddr, abi, ownSK, processor, accountOwn, ownTransferNonce, "transfer", aliceAddress.Hex()+",20", 3, singer)
	assert.NoError(t, err)

	//  合约调用getBalance方法  获取alice账户balance
	ownTransferNonce++
	err = processContractCall(t, contractAddr, abi, ownSK, processor, accountOwn, ownTransferNonce, "getBalance", aliceAddress.Hex(), 4, singer)
	assert.NoError(t, err)

}

func TestContractPaymentChannel(t *testing.T) {
	singer := model.NewSigner(new(big.Int).SetInt64(int64(1)))

	ownSK, _ := crypto.GenerateKey()
	ownPk := ownSK.PublicKey
	ownAddress := cs_crypto.GetNormalAddress(ownPk)

	aliceSK, _ := crypto.GenerateKey()
	alicePk := aliceSK.PublicKey
	aliceAddress := cs_crypto.GetNormalAddress(alicePk)
	//alicdAddr := address_util.PubKeyToAddress(alicePk, common.AddressTypeNormal)

	brotherSK, _ := crypto.GenerateKey()
	brotherPk := brotherSK.PublicKey
	brotherAddress := cs_crypto.GetNormalAddress(brotherPk)

	addressSlice := []common.Address{
		ownAddress,
		aliceAddress,
		brotherAddress,
	}

	WASMPath := model.GetWASMPath("PaymentChannel", model.CoreVmTestData)
	abiPath := model.GetAbiPath("PaymentChannel", model.CoreVmTestData)
	input := []string{aliceAddress.Hex(), "1573293024432297000", "10"}

	data, err := getCreateExtraData(WASMPath, abiPath, input)
	assert.NoError(t, err)

	addr := common.HexToAddress(common.AddressContractCreate)
	tx := model.NewTransaction(0, addr, big.NewInt(10), big.NewInt(1), 26427000, data)
	signCreateTx := getSignedTx(t, ownSK, tx, singer)

	gasLimit := testGasLimit * 10000000000
	block := CreateBlock(1, common.Hash{}, []*model.Transaction{signCreateTx}, gasLimit)
	processor, err := CreateProcessorAndInitAccount(t, addressSlice)

	tmpGasLimit := block.GasLimit()
	gasUsed := block.GasUsed()
	config := &TxProcessConfig{
		Tx:       tx,
		Header:   block.Header(),
		GetHash:  getTestHashFunc(),
		GasLimit: &tmpGasLimit,
		GasUsed:  &gasUsed,
		TxFee:    big.NewInt(0),
	}

	err = processor.ProcessTxNew(config)
	assert.NoError(t, err)

	contractAddr := cs_crypto.CreateContractAddress(ownAddress, uint64(0))
	_, err = processor.GetNonce(contractAddr)
	_, err = processor.GetCode(contractAddr)
	abi, err := processor.GetAbi(contractAddr)
	assert.NoError(t, err)
	processor.Commit()

	accountOwn := accountsbase.Account{ownAddress}
	//  合约调用extend方法，延长支付通道的最早可关闭时间
	ownTransferNonce, err := processor.GetNonce(ownAddress)
	assert.NoError(t, err)
	err = processContractCall(t, contractAddr, abi, ownSK, processor, accountOwn, ownTransferNonce, "extend", "1573293351343372000", 2, singer)
	assert.NoError(t, err)

	ownTransferNonce++

	// 合约调用  close 方法
	signMessage := contractAddr.Hex() + "1" + aliceAddress.Hex()
	signHash := crypto.Keccak256([]byte(signMessage))
	signature, err := crypto.Sign(signHash, ownSK)
	assert.NoError(t, err)

	err = processContractCall(t, contractAddr, abi, aliceSK, processor, accountOwn, 0,
		"close", "1,"+common.Bytes2Hex(signature), 3, singer)

	assert.NoError(t, err)

	//ownTransferNonce++
	err = processContractCall(t, contractAddr, abi, aliceSK, processor, accountOwn, 1, "close", "1,"+common.Bytes2Hex(signature), 4, singer)
	assert.Error(t, err)

}

func newContractCallTx(from *common.Address, to *common.Address, gasPrice *big.Int, gasLimit uint64, funcName string, input string, nonce uint64, code []byte) (tx *model.Transaction, err error) {
	// RLP([funcName][params])
	inputRlp, err := rlp.EncodeToBytes([]interface{}{
		funcName, input,
	})
	if err != nil {
		log.DLogger.Error("input rlp err")
		return
	}

	extraData, err := utils.ParseCallContractData(code, inputRlp)

	if err != nil {
		log.DLogger.Error("ParseCallContractData  inputRlp",  zap.Error(err))
		return
	}

	tx = model.NewTransaction(nonce, *to, nil, gasPrice, gasLimit, extraData)
	return tx, nil

}

//  合约调用getBalance方法
func processContractCall(t *testing.T, contractAddress common.Address, code []byte, priKey *ecdsa.PrivateKey, processor *AccountStateDB, accountOwn accountsbase.Account, nonce uint64, funcName string, params string, blockNum uint64, singer model.Signer) error {
	gasUsed2 := uint64(0)
	gasLimit := testGasLimit * 10000000000
	callTx, err := newContractCallTx(nil, &contractAddress, new(big.Int).SetUint64(1), uint64(1500000), funcName, params, nonce, code)
	assert.NoError(t, err)
	signCallTx, err := callTx.SignTx(priKey, singer)
	signPk, err := signCallTx.SenderPublicKey(singer)
	assert.NoError(t, err)
	addr := addressutil.PubKeyToAddress(*signPk, common.AddressTypeNormal)
	fmt.Println("processContractCall", "addr", addr)

	//sw.SignTx(accountOwn, callTx, nil)
	assert.NoError(t, err)
	block := CreateBlock(blockNum, common.Hash{}, []*model.Transaction{signCallTx}, gasLimit)
	log.DLogger.Info("callTx info", zap.Any("callTx",callTx))
	txConfig := &TxProcessConfig{
		Tx:       signCallTx,
		Header:   block.Header().(*model.Header),
		GetHash:  fakeGetBlockHash,
		GasLimit: &gasLimit,
		GasUsed:  &gasUsed2,
	}
	err = processor.ProcessTxNew(txConfig)
	//if funcName == "getBalance" {
	receipt := callTx.GetReceipt()
	fmt.Println("receipt  log", "receipt log", receipt)
	//}

	//assert.NoError(t, err)
	processor.Commit()
	return err
}

func CreateProcessorAndInitAccount(t *testing.T, addressSlice []common.Address) (*AccountStateDB, error) {
	db, root := CreateTestStateDB()
	processor, err := NewAccountStateDB(root, NewStateStorageWithCache(db))

	assert.NoError(t, err)
	processor.NewAccountState(addressSlice[0])
	err = processor.AddNonce(addressSlice[0], 0)
	processor.AddBalance(addressSlice[0], new(big.Int).SetInt64(int64(1000000000000000000)))
	for i := 1; i < len(addressSlice); i++ {
		fmt.Println("xxxxxxxxxxxxxxxxx", addressSlice[i])
		processor.NewAccountState(addressSlice[i])
		err = processor.AddNonce(addressSlice[i], 0)
		processor.AddBalance(addressSlice[i], new(big.Int).SetInt64(int64(10000000)))

	}
	return processor, err
}

func getSignedTx(t *testing.T, priKey *ecdsa.PrivateKey, tx *model.Transaction, singer model.Signer) *model.Transaction {
	signCreateTx, err := tx.SignTx(priKey, singer)
	assert.NoError(t, err)
	return signCreateTx
}

