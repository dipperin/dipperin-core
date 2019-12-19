package state_processor

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/address-util"
	"github.com/dipperin/dipperin-core/common/math"
	"github.com/dipperin/dipperin-core/core/accounts"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/core/vm/common/utils"
	"github.com/dipperin/dipperin-core/tests/g-testData"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestAccountStateDB_ProcessContract(t *testing.T) {
	WASMPath := g_testData.GetWASMPath("event", g_testData.CoreVmTestData)
	abiPath := g_testData.GetAbiPath("event", g_testData.CoreVmTestData)
	tx1 := createContractTx(WASMPath, abiPath, 0, testGasLimit)
	contractAddr := cs_crypto.CreateContractAddress(aliceAddr, 0)
	name := []byte("ProcessContract")
	param := [][]byte{name}
	tx2 := callContractTx(&contractAddr, "returnString", param, 1)

	db, root := CreateTestStateDB()
	tdb := NewStateStorageWithCache(db)
	processor, err := NewAccountStateDB(root, tdb)
	assert.NoError(t, err)

	block := CreateBlock(1, common.Hash{}, []*model.Transaction{tx1, tx2}, 5*testGasLimit)
	tmpGasLimit := block.GasLimit()
	gasUsed := block.GasUsed()
	config := &TxProcessConfig{
		Tx:       tx1,
		Header:   block.Header(),
		GetHash:  getTestHashFunc(),
		GasLimit: &tmpGasLimit,
		GasUsed:  &gasUsed,
		TxFee:    big.NewInt(0),
	}
	err = processor.ProcessTxNew(config)
	assert.NoError(t, err)

	receipt1 := tx1.GetReceipt()
	assert.Len(t, receipt1.Logs, 0)

	fmt.Println("---------------------------")

	config.Tx = tx2
	err = processor.ProcessTxNew(config)
	assert.NoError(t, err)

	receipt2 := tx2.GetReceipt()
	assert.Len(t, receipt2.Logs, 1)

	log1 := receipt2.Logs[0]
	assert.Equal(t, tx2.CalTxId(), log1.TxHash)
	assert.Equal(t, contractAddr, log1.Address)
}

func TestAccountStateDB_ProcessContract_Error(t *testing.T) {
	to := common.HexToAddress(common.AddressContractCreate)
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
	assert.NoError(t, err)
	err = processor.ProcessTxNew(config)
	assert.Equal(t, "rlp: cannot encode negative *big.Int", err.Error())

	WASMPath := g_testData.GetWASMPath("event", g_testData.CoreVmTestData)
	abiPath := g_testData.GetAbiPath("event", g_testData.CoreVmTestData)
	tx = createContractTx(WASMPath, abiPath, 0, testGasLimit)
	block = CreateBlock(1, common.Hash{}, []*model.Transaction{tx}, uint64(100))
	tmpGasLimit = block.GasLimit()
	gasUsed = block.GasUsed()
	config = &TxProcessConfig{
		Tx:       tx,
		Header:   block.Header(),
		GetHash:  getTestHashFunc(),
		GasLimit: &tmpGasLimit,
		GasUsed:  &gasUsed,
		TxFee:    big.NewInt(0),
	}
	err = processor.ProcessTxNew(config)
	assert.Equal(t, "gas limit reached", err.Error())
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

	WASMPath := g_testData.GetWASMPath("token", g_testData.CoreVmTestData)
	abiPath := g_testData.GetAbiPath("token", g_testData.CoreVmTestData)
	//WASMPath := g_testData.GetWASMPath("token-param", g_testData.CoreVmTestData)
	//abiPath := g_testData.GetAbiPath("token-param", g_testData.CoreVmTestData)
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
	contractNonce, err := processor.GetNonce(contractAddr)
	log.Info("TestAccountStateDB_ProcessContract", "contractNonce", contractNonce)
	_, err = processor.GetCode(contractAddr)
	abi, err := processor.GetAbi(contractAddr)
	//log.Info("TestAccountStateDB_ProcessContract", "code  get from state", code)
	assert.NoError(t, err)
	//assert.Equal(t, code, tx.ExtraData())
	processor.Commit()

	accountOwn := accounts.Account{ownAddress}
	//  合约调用getBalance方法  获取合约原始账户balance
	ownTransferNonce, err := processor.GetNonce(ownAddress)
	assert.NoError(t, err)
	err = processContractCall(t, contractAddr, big.NewInt(0), abi, ownSK, processor, accountOwn, ownTransferNonce, "getBalance", ownAddress.Hex(), 2, singer)
	assert.NoError(t, err)

	gasUsed2 := uint64(0)
	//  合约调用  transfer方法 转账给alice
	ownTransferNonce++
	err = processContractCall(t, contractAddr, big.NewInt(0), abi, ownSK, processor, accountOwn, ownTransferNonce, "transfer", aliceAddress.Hex()+",20", 3, singer)
	assert.NoError(t, err)

	//  合约调用getBalance方法  获取alice账户balance
	ownTransferNonce++
	err = processContractCall(t, contractAddr, big.NewInt(0), abi, ownSK, processor, accountOwn, ownTransferNonce, "getBalance", aliceAddress.Hex(), 4, singer)
	assert.NoError(t, err)

	//  合约调用approve方法
	log.Info("==========================================")
	ownTransferNonce++
	callTxApprove, err := newContractCallTx(nil, &contractAddr, big.NewInt(0), new(big.Int).SetUint64(1), uint64(1500000), "approve", brotherAddress.Hex()+",50", ownTransferNonce, abi)
	//accountAlice := accounts.Account{aliceAddress}
	signCallTxApprove, err := callTxApprove.SignTx(ownSK, singer)

	assert.NoError(t, err)
	block5 := CreateBlock(5, common.Hash{}, []*model.Transaction{signCallTxApprove}, gasLimit)
	log.Info("signCallTxApprove info", "signCallTxApprove", signCallTxApprove)

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

	//  合约调用getApproveBalance方法  获取own授权给brother账户balance
	/*err = processContractCall(t, contractAddr, abi, ownSK,  processor, accountOwn, 5, "getApproveBalance", "0x000062be10f46b5d01Ecd9b502c4bA3d6131f6fc2e41,0x00004179D57e45Cb3b54D6FAEF69e746bf240E287978", 6)
	assert.NoError(t, err)*/

	//  合约调用transferFrom方法
	log.Info("==========================================")
	callTxTransferFrom, err := newContractCallTx(nil, &contractAddr, big.NewInt(0), new(big.Int).SetUint64(1), uint64(1500000), "transferFrom", ownAddress.Hex()+","+aliceAddress.Hex()+",5", 0, abi)
	assert.NoError(t, err)
	accountBrother := accounts.Account{Address: brotherAddress}
	assert.NoError(t, err)

	signCallTxTransferFrom, err := callTxTransferFrom.SignTx(brotherSK, singer)
	assert.NoError(t, err)
	block7 := CreateBlock(7, common.Hash{}, []*model.Transaction{signCallTxTransferFrom}, gasLimit)
	log.Info("signCallTxTransferFrom info", "signCallTxTransferFrom", signCallTxTransferFrom)

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
	err = processContractCall(t, contractAddr, big.NewInt(0), abi, ownSK, processor, accountOwn, ownTransferNonce, "getBalance", aliceAddress.Hex(), 8, singer)
	assert.NoError(t, err)

	//  合约调用getBalance方法  获取own账户最终的balance
	ownTransferNonce++
	err = processContractCall(t, contractAddr, big.NewInt(0), abi, ownSK, processor, accountOwn, ownTransferNonce, "getBalance", ownAddress.Hex(), 9, singer)
	assert.NoError(t, err)

	// 合约调用  transfer方法  转账给brother
	ownTransferNonce++
	err = processContractCall(t, contractAddr, big.NewInt(0), abi, ownSK, processor, accountOwn, ownTransferNonce, "transfer", brotherAddress.Hex()+",28", 10, singer)
	assert.NoError(t, err)

	//  合约调用getBalance方法  获取own账户最终的balance
	ownTransferNonce++
	err = processContractCall(t, contractAddr, big.NewInt(0), abi, ownSK, processor, accountOwn, ownTransferNonce, "getBalance", ownAddress.Hex(), 11, singer)
	assert.NoError(t, err)

	// 合约调用burn方法,将账户余额返还给own
	err = processContractCall(t, contractAddr, big.NewInt(0), abi, brotherSK, processor, accountBrother, 1, "burn", "15", 12, singer)
	assert.NoError(t, err)

	// 合约调用getBalance方法,获取own的余额
	ownTransferNonce++
	err = processContractCall(t, contractAddr, big.NewInt(0), abi, ownSK, processor, accountOwn, ownTransferNonce, "getBalance", ownAddress.Hex(), 13, singer)
	assert.NoError(t, err)

	// 合约调用setName方法，设置合约名
	ownTransferNonce++
	err = processContractCall(t, contractAddr, big.NewInt(0), abi, ownSK, processor, accountOwn, ownTransferNonce, "setName", "wujinhai", 14, singer)
	assert.NoError(t, err)

	log.Info("TestAccountStateDB_ProcessContract++", "callRecipt", "", "err", err)
}

func TestContractNewFeature(t *testing.T) {
	log.InitLogger(log.LvlDebug)
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

	//WASMPath := g_testData.GetWASMPath("token", g_testData.CoreVmTestData)
	//abiPath := g_testData.GetAbiPath("token", g_testData.CoreVmTestData)
	WASMPath := g_testData.GetWASMPath("demo", g_testData.CoreVmTestData)
	abiPath := g_testData.GetAbiPath("demo", g_testData.CoreVmTestData)
	input := []string{"0x0000D36F282D8925B16Ed24CB637475e6a03B01E1056"}
	//0x0000d36F282D8925B16Ed24cb637475e6A03b01E1056
	//0x0000d36F282D8925B16Ed24cb637475e6A03b01E1056
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
	WASMPath := g_testData.GetWASMPath("token-param", g_testData.CoreVmTestData)
	abiPath := g_testData.GetAbiPath("token-param", g_testData.CoreVmTestData)
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
	contractNonce, err := processor.GetNonce(contractAddr)
	log.Info("TestAccountStateDB_ProcessContract", "contractNonce", contractNonce)
	code, err := processor.GetCode(contractAddr)
	abi, err := processor.GetAbi(contractAddr)
	log.Info("TestAccountStateDB_ProcessContract", "code  get from state", code)
	assert.NoError(t, err)
	//assert.Equal(t, code, tx.ExtraData())
	processor.Commit()

	accountOwn := accounts.Account{ownAddress}
	//  合约调用getBalance方法  获取合约原始账户balance
	ownTransferNonce, err := processor.GetNonce(ownAddress)
	assert.NoError(t, err)
	err = processContractCall(t, contractAddr, big.NewInt(0), abi, ownSK, processor, accountOwn, ownTransferNonce, "getBalance", ownAddress.Hex(), 2, singer)
	assert.NoError(t, err)

	//gasUsed2 := uint64(0)
	//  合约调用  transfer方法 转账给alice
	ownTransferNonce++
	err = processContractCall(t, contractAddr, big.NewInt(0), abi, ownSK, processor, accountOwn, ownTransferNonce, "transfer", aliceAddress.Hex()+",20", 3, singer)
	assert.NoError(t, err)

	//  合约调用getBalance方法  获取alice账户balance
	ownTransferNonce++
	err = processContractCall(t, contractAddr, big.NewInt(0), abi, ownSK, processor, accountOwn, ownTransferNonce, "getBalance", aliceAddress.Hex(), 4, singer)
	assert.NoError(t, err)

}

//func TestContractpayableCallNotPayableMulti(t *testing.T) {
//	singer := model.NewSigner(new(big.Int).SetInt64(int64(1)))
//
//	ownSK, _ := crypto.GenerateKey()
//	ownPk := ownSK.PublicKey
//	ownAddress := cs_crypto.GetNormalAddress(ownPk)
//
//	aliceSK, _ := crypto.GenerateKey()
//	alicePk := aliceSK.PublicKey
//	aliceAddress := cs_crypto.GetNormalAddress(alicePk)
//
//	brotherSK, _ := crypto.GenerateKey()
//	brotherPk := brotherSK.PublicKey
//	brotherAddress := cs_crypto.GetNormalAddress(brotherPk)
//
//	addressSlice := []common.Address{
//		ownAddress,
//		aliceAddress,
//		brotherAddress,
//	}
//
//	//WASMPath := g_testData.GetWASMPath("token", g_testData.CoreVmTestData)
//	//abiPath := g_testData.GetAbiPath("token", g_testData.CoreVmTestData)
//	WASMPath := g_testData.GetWASMPath("payableCallNotPayableMulti", g_testData.CoreVmTestData)
//	abiPath := g_testData.GetAbiPath("payableCallNotPayableMulti", g_testData.CoreVmTestData)
//	input := []string{"dipp", "DIPP", "10000000"}
//	data, err := getCreateExtraData(WASMPath, abiPath, input)
//	assert.NoError(t, err)
//
//	addr := common.HexToAddress(common.AddressContractCreate)
//	tx := model.NewTransaction(0, addr, big.NewInt(0), big.NewInt(1), 26427000, data)
//	signCreateTx := getSignedTx(t, ownSK, tx, singer)
//
//	gasLimit := testGasLimit * 10000000000
//	block := CreateBlock(1, common.Hash{}, []*model.Transaction{signCreateTx}, gasLimit)
//	processor, err := CreateProcessorAndInitAccount(t, addressSlice)
//
//	tmpGasLimit := block.GasLimit()
//	gasUsed := block.GasUsed()
//	config := &TxProcessConfig{
//		Tx:       tx,
//		Header:   block.Header(),
//		GetHash:  getTestHashFunc(),
//		GasLimit: &tmpGasLimit,
//		GasUsed:  &gasUsed,
//		TxFee:    big.NewInt(0),
//	}
//
//	err = processor.ProcessTxNew(config)
//	assert.NoError(t, err)
//
//	contractAddr := cs_crypto.CreateContractAddress(ownAddress, uint64(0))
//	contractNonce, err := processor.GetNonce(contractAddr)
//	log.Info("TestAccountStateDB_ProcessContract", "contractNonce", contractNonce)
//	code, err := processor.GetCode(contractAddr)
//	abi, err := processor.GetAbi(contractAddr)
//	log.Info("TestAccountStateDB_ProcessContract", "code  get from state", code)
//	assert.NoError(t, err)
//	//assert.Equal(t, code, tx.ExtraData())
//	processor.Commit()
//
//	accountOwn := accounts.Account{ownAddress}
//	//  合约调用getBalance方法  获取合约原始账户balance
//	ownTransferNonce, err := processor.GetNonce(ownAddress)
//	assert.NoError(t, err)
//	err = processContractCall(t, contractAddr, abi, ownSK, processor, accountOwn, ownTransferNonce, "getBalance", ownAddress.Hex(), 2, singer)
//	assert.NoError(t, err)
//
//	//gasUsed2 := uint64(0)
//	//  合约调用  transfer方法 转账给alice
//	ownTransferNonce++
//	err = processContractCall(t, contractAddr, abi, ownSK, processor, accountOwn, ownTransferNonce, "transfer", aliceAddress.Hex()+",20", 3, singer)
//	assert.NoError(t, err)
//
//	//  合约调用getBalance方法  获取alice账户balance
//	ownTransferNonce++
//	err = processContractCall(t, contractAddr, abi, ownSK, processor, accountOwn, ownTransferNonce, "getBalance", aliceAddress.Hex(), 4, singer)
//	assert.NoError(t, err)
//
//}

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

	//WASMPath := g_testData.GetWASMPath("token", g_testData.CoreVmTestData)
	//abiPath := g_testData.GetAbiPath("token", g_testData.CoreVmTestData)
	WASMPath := g_testData.GetWASMPath("PaymentChannel", g_testData.CoreVmTestData)
	abiPath := g_testData.GetAbiPath("PaymentChannel", g_testData.CoreVmTestData)
	fmt.Println("aliceAddr hex", aliceAddress.Hex())
	fmt.Println("ownAddr hex", ownAddress.Hex())
	//input := []string{"123456789012345678901234","1573293024432297000","10"}
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
	log.Info("TestContractPaymentChannel contractAddr", "contractAddr", contractAddr)
	contractNonce, err := processor.GetNonce(contractAddr)
	log.Info("TestAccountStateDB_ProcessContract", "contractNonce", contractNonce)
	code, err := processor.GetCode(contractAddr)
	abi, err := processor.GetAbi(contractAddr)
	log.Info("TestAccountStateDB_ProcessContract", "code  get from state", code)
	assert.NoError(t, err)
	//assert.Equal(t, code, tx.ExtraData())
	processor.Commit()

	accountOwn := accounts.Account{ownAddress}
	//  合约调用extend方法，延长支付通道的最早可关闭时间
	ownTransferNonce, err := processor.GetNonce(ownAddress)
	assert.NoError(t, err)
	err = processContractCall(t, contractAddr, big.NewInt(0), abi, ownSK, processor, accountOwn, ownTransferNonce, "extend", "1573293351343372000", 2, singer)
	assert.NoError(t, err)

	//gasUsed2 := uint64(0)
	//合约调用 错误调用 extend
	ownTransferNonce++
	//err = processContractCall(t, contractAddr, abi, ownSK, processor, accountOwn, ownTransferNonce, "extend", "1573293024432297000", 3, singer)
	//assert.Error(t, err)

	// 合约调用  close 方法
	signMessage := contractAddr.Hex() + "1" + aliceAddress.Hex()
	log.Info("TestContractPaymentChannel#signature", "signMessage", signMessage)
	signHash := crypto.Keccak256([]byte(signMessage))
	signature, err := crypto.Sign(signHash, ownSK)
	log.Info("TestContractPaymentChannel#signature", "signature", signature, "signHash", signHash, "sign byte", common.Bytes2Hex(signature))
	assert.NoError(t, err)

	err = processContractCall(t, contractAddr, big.NewInt(0), abi, aliceSK, processor, accountOwn, 0,
		"close", "1,"+common.Bytes2Hex(signature), 3, singer)

	assert.NoError(t, err)

	//  合约再次调用close方法，报错
	//ownTransferNonce++
	err = processContractCall(t, contractAddr, big.NewInt(0), abi, aliceSK, processor, accountOwn, 1, "close", "1,"+common.Bytes2Hex(signature), 4, singer)
	assert.Error(t, err)

}

func TestContractPortfolioManage(t *testing.T) {
	log.InitLogger(log.LvlDebug)
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

	//WASMPath := g_testData.GetWASMPath("token", g_testData.CoreVmTestData)
	//abiPath := g_testData.GetAbiPath("token", g_testData.CoreVmTestData)
	WASMPath := g_testData.GetWASMPath("PortfolioManage", g_testData.CoreVmTestData)
	abiPath := g_testData.GetAbiPath("PortfolioManage", g_testData.CoreVmTestData)
	fmt.Println("aliceAddr hex", aliceAddress.Hex())
	fmt.Println("ownAddr hex", ownAddress.Hex())
	//input := []string{"123456789012345678901234","1573293024432297000","10"}
	input := []string{}

	data, err := getCreateExtraData(WASMPath, abiPath, input)
	assert.NoError(t, err)

	addr := common.HexToAddress(common.AddressContractCreate)
	fmt.Println("bigint", new(big.Int).Mul(new(big.Int).SetInt64(10), math.BigPow(10, 18)))
	eachValue := new(big.Int).Mul(new(big.Int).SetInt64(10), math.BigPow(10, 18))
	tx := model.NewTransaction(0, addr, eachValue, new(big.Int).SetInt64(1), 26427000, data)
	signCreateTx := getSignedTx(t, ownSK, tx, singer)

	gasLimit := uint64(3360000000)
	blockNum := uint64(1)
	//gasLimit := testGasLimit * 10000000000
	block := CreateBlock(blockNum, common.Hash{}, []*model.Transaction{signCreateTx}, gasLimit)
	blockNum++
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

	receipt := tx.GetReceipt()
	fmt.Println("receipt  log", "receipt log", receipt)
	contractAddr := cs_crypto.CreateContractAddress(ownAddress, uint64(0))
	log.Info("TestContractPortfolioManage contractAddr", "contractAddr", contractAddr)
	contractNonce, err := processor.GetNonce(contractAddr)
	log.Info("TestContractPortfolioManage  ", "contractNonce", contractNonce)
	code, err := processor.GetCode(contractAddr)
	abi, err := processor.GetAbi(contractAddr)
	log.Info("TestContractPortfolioManage  ", "code  get from state", code)
	assert.NoError(t, err)
	//assert.Equal(t, code, tx.ExtraData())
	processor.Commit()

	accountOwn := accounts.Account{ownAddress}
	//  合约调用createPortfolio方法，创建投资组合
	ownTransferNonce, err := processor.GetNonce(ownAddress)
	assert.NoError(t, err)
	aliceNonce := uint64(0)
	err = processContractCall(t, contractAddr, eachValue, abi, aliceSK, processor, accountOwn, aliceNonce, "createPortfolio", "winner,winner", blockNum, singer)
	blockNum++
	aliceNonce++
	assert.NoError(t, err)

	//gasUsed2 := uint64(0)
	//合约调用 错误调用 extend

	//err = processContractCall(t, contractAddr, abi, ownSK, processor, accountOwn, ownTransferNonce, "extend", "1573293024432297000", 3, singer)
	//assert.Error(t, err)

	// 合约调用  createOrder 方法

	err = processContractCall(t, contractAddr, eachValue, abi, aliceSK, processor, accountOwn, aliceNonce,
		"createOrder", "winner,1234,1,000002,100,100,15000000000", blockNum, singer)
	blockNum++
	aliceNonce++
	assert.NoError(t, err)

	// 合约 第二次  调用  createOrder 方法
	err = processContractCall(t, contractAddr, big.NewInt(0), abi, aliceSK, processor, accountOwn, aliceNonce,
		"createOrder", "winner,12345,1,000002,100,100,15000000000", blockNum, singer)
	blockNum++
	aliceNonce++
	assert.NoError(t, err)

	// 合约调用  dealOrder 方法
	//ownTransferNonce++
	err = processContractCall(t, contractAddr, big.NewInt(0), abi, ownSK, processor, accountOwn, ownTransferNonce,
		"dealOrder", "winner,1234", blockNum, singer)
	blockNum++
	ownTransferNonce++

	assert.NoError(t, err)

	// 合约调用  revocationOrder 方法
	//ownTransferNonce++
	err = processContractCall(t, contractAddr, big.NewInt(0), abi, aliceSK, processor, accountOwn, aliceNonce,
		"revocationOrder", "winner,12345", 6, singer)

	assert.NoError(t, err)

}

func newContractCallTx(from *common.Address, to *common.Address, amount *big.Int, gasPrice *big.Int, gasLimit uint64, funcName string, input string, nonce uint64, code []byte) (tx *model.Transaction, err error) {
	// RLP([funcName][params])
	inputRlp, err := rlp.EncodeToBytes([]interface{}{
		funcName, input,
	})
	if err != nil {
		log.Error("input rlp err")
		return
	}

	extraData, err := utils.ParseCallContractData(code, inputRlp)

	if err != nil {
		log.Error("ParseCallContractData  inputRlp", "err", err)
		return
	}

	tx = model.NewTransaction(nonce, *to, amount, gasPrice, gasLimit, extraData)
	return tx, nil

}

type  contractCallConf struct {
	contractAddress common.Address
	amount *big.Int
	code []byte
	priKey *ecdsa.PrivateKey
	processor *AccountStateDB
	accountOwn *accounts.Account
	nonce uint64
	funcName string
	params string
	blockNum uint64
	singer *model.Signer
}

//  合约调用getBalance方法
func processContractCall(t *testing.T, contractAddress common.Address, amount *big.Int, code []byte, priKey *ecdsa.PrivateKey, processor *AccountStateDB, accountOwn accounts.Account, nonce uint64, funcName string, params string, blockNum uint64, singer model.Signer) error {
	gasUsed2 := uint64(0)
	gasLimit := testGasLimit * 10000000000
	log.Info("processContractCall=================================================")
	callTx, err := newContractCallTx(nil, &contractAddress, amount, new(big.Int).SetUint64(1), uint64(150000000000), funcName, params, nonce, code)
	assert.NoError(t, err)
	signCallTx, err := callTx.SignTx(priKey, singer)
	signPk, err := signCallTx.SenderPublicKey(singer)
	assert.NoError(t, err)
	addr := address_util.PubKeyToAddress(*signPk, common.AddressTypeNormal)
	fmt.Println("processContractCall", "addr", addr)

	//sw.SignTx(accountOwn, callTx, nil)
	assert.NoError(t, err)
	block := CreateBlock(blockNum, common.Hash{}, []*model.Transaction{signCallTx}, gasLimit)
	log.Info("callTx info", "callTx", callTx)
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
	processor.AddBalance(addressSlice[0], new(big.Int).Mul(new(big.Int).SetInt64(1000000), math.BigPow(10, 18)))
	for i := 1; i < len(addressSlice); i++ {
		fmt.Println("xxxxxxxxxxxxxxxxx", addressSlice[i])
		processor.NewAccountState(addressSlice[i])
		err = processor.AddNonce(addressSlice[i], 0)
		processor.AddBalance(addressSlice[i], new(big.Int).Mul(new(big.Int).SetInt64(1000000), math.BigPow(10, 18)))

	}
	return processor, err
}

func getSignedTx(t *testing.T, priKey *ecdsa.PrivateKey, tx *model.Transaction, singer model.Signer) *model.Transaction {
	signCreateTx, err := tx.SignTx(priKey, singer)
	assert.NoError(t, err)
	return signCreateTx
}
