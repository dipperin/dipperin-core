package service

import (
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/core/accounts"
	"github.com/dipperin/dipperin-core/core/chain-communication"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/core/cs-chain/chain-state"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/core/vm/common/utils"
	"github.com/dipperin/dipperin-core/tests"
	"github.com/dipperin/dipperin-core/tests/g-testData"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
	"math/big"
	"os"
	"testing"
)

func TestMercuryFullChainService_Call(t *testing.T) {
	csChain := createCsChain(nil)
	config := DipperinConfig{ChainReader: csChain}
	service := MakeFullChainService(&config)

	// set abi
	WASMPath := g_testData.GetWASMPath("token-const", g_testData.CoreVmTestData)
	AbiPath := g_testData.GetAbiPath("token-const", g_testData.CoreVmTestData)
	tx := createContractTx(0, WASMPath, AbiPath, "DIPP,WU,10000", nil)
	block := createBlock(csChain, []*model.Transaction{tx}, nil)
	votes := createVerifiersVotes(block, csChain.ChainConfig.VerifierNumber, nil)
	err := csChain.SaveBftBlock(block, votes)
	assert.NoError(t, err)

	// call
	fmt.Println("-----------------------------------------------------------------------------")
	sender, err := tx.Sender(nil)
	contractAddr := cs_crypto.CreateContractAddress(sender, uint64(0))
	data, _ := rlp.EncodeToBytes([]interface{}{"getBalance", sender.String()})
	_, abi := g_testData.GetCodeAbi(WASMPath, AbiPath)
	extraData, err := utils.ParseCallContractData(abi, data)
	assert.NoError(t, err)
	tx = createSignedTx(1, contractAddr, g_testData.TestValue, extraData, nil)

	balance, err := service.Call(tx, 1)
	assert.NoError(t, err)
	assert.Equal(t, "10000", balance)
}

func TestMercuryFullChainService_ContractTransaction(t *testing.T) {
	csChain := createCsChain(nil)

	// create create contract
	WASMPath1 := g_testData.GetWASMPath("token", g_testData.CoreVmTestData)
	AbiPath1 := g_testData.GetAbiPath("token", g_testData.CoreVmTestData)
	tx1 := createContractTx(0, WASMPath1, AbiPath1, "DIPP,WU,10000", nil)
	WASMPath2 := g_testData.GetWASMPath("event", g_testData.CoreVmTestData)
	AbiPath2 := g_testData.GetAbiPath("event", g_testData.CoreVmTestData)
	tx2 := createContractTx(1, WASMPath2, AbiPath2, "", nil)
	block := createBlock(csChain, []*model.Transaction{tx1, tx2}, nil)
	votes := createVerifiersVotes(block, csChain.ChainConfig.VerifierNumber, nil)
	err := csChain.SaveBftBlock(block, votes)
	assert.NoError(t, err)

	chainIndex := chain_state.NewBloomIndexer(csChain, csChain.GetDB(), chain_state.BloomBitsBlocks, chain_state.BloomConfirms)
	config := &DipperinConfig{
		ChainReader: csChain,
		ChainIndex:  chainIndex,
	}
	service := VenusFullChainService{DipperinConfig: config}

	// get contract address
	sender, err := tx1.Sender(nil)
	contractAddr := cs_crypto.CreateContractAddress(sender, uint64(0))
	addr, err := service.GetContractAddressByTxHash(tx1.CalTxId())
	assert.NoError(t, err)
	assert.Equal(t, contractAddr, addr)

	// get receipt
	receipt1, err := service.GetReceiptByTxHash(tx1.CalTxId())
	assert.NoError(t, err)
	receipt2, err := service.GetReceiptByTxHash(tx2.CalTxId())
	assert.NoError(t, err)
	receipts, err := service.GetReceiptsByBlockNum(block.Number())
	assert.NoError(t, err)
	assert.Equal(t, receipt1, receipts[0])
	assert.Equal(t, receipt2, receipts[1])

	// get tx actual fee
	fee1, err := service.GetTxActualFee(tx1.CalTxId())
	assert.NoError(t, err)
	fee2, err := service.GetTxActualFee(tx2.CalTxId())
	assert.NoError(t, err)
	expectFee1 := receipt1.GasUsed * g_testData.TestGasPrice.Uint64()
	expectFee2 := receipt2.GasUsed * g_testData.TestGasPrice.Uint64()
	assert.Equal(t, expectFee1, fee1.Uint64())
	assert.Equal(t, expectFee2, fee2.Uint64())

	// get logs
	topic := []common.Hash{common.BytesToHash(crypto.Keccak256([]byte("Tranfer")))}
	logs, err := service.GetLogs(block.Hash(), uint64(0), uint64(100), []common.Address{addr}, [][]common.Hash{topic})
	assert.NoError(t, err)
	assert.Equal(t, receipt1.Logs, logs)

	err = service.Start()
	assert.NoError(t, err)
	logs, err = service.GetLogs(common.Hash{}, uint64(0), uint64(100), []common.Address{addr}, [][]common.Hash{topic})
	assert.NoError(t, err)
	assert.Equal(t, receipt1.Logs, logs)
}

func TestMercuryFullChainService_EstimateGas(t *testing.T) {
	csChain := createCsChain(nil)
	config := DipperinConfig{ChainReader: csChain}
	service := MakeFullChainService(&config)

	WASMPath := g_testData.GetWASMPath("token-const", g_testData.CoreVmTestData)
	AbiPath := g_testData.GetAbiPath("token-const", g_testData.CoreVmTestData)
	tx := createContractTx(0, WASMPath, AbiPath, "DIPP,WU,10000", nil)
	block := createBlock(csChain, []*model.Transaction{tx}, nil)
	votes := createVerifiersVotes(block, csChain.ChainConfig.VerifierNumber, nil)
	err := csChain.SaveBftBlock(block, votes)

	receipt, err := service.GetReceiptByTxHash(tx.CalTxId())
	assert.NoError(t, err)

	gas, err := service.EstimateGas(tx, 0)
	assert.NoError(t, err)
	assert.Equal(t, receipt.GasUsed, (uint64)(gas))
}

func TestMercuryFullChainService_MakeTmpSignedTx(t *testing.T) {
	manager := createWalletManager(t)
	defer os.Remove(util.HomeDir() + testPath)
	account, err := manager.Wallets[0].Accounts()
	assert.NoError(t, err)

	address := account[0].Address
	pk, err := manager.Wallets[0].GetSKFromAddress(address)
	testAccount := tests.NewAccount(pk, address)
	testAccounts := []tests.Account{*testAccount}

	serviceChain := createCsChainService(testAccounts)
	txPool := createTxPool(serviceChain.ChainState)
	serviceChain.TxPool = txPool

	broadcaster := chain_communication.NewBroadcastDelegate(txPool, fakeNodeConfig{}, fakePeerManager{}, serviceChain, fakePbftNode{})
	config := &DipperinConfig{
		NodeConf:      fakeNodeConfig{nodeType: chain_config.NodeTypeOfVerifier},
		WalletManager: manager,
		ChainReader:   serviceChain,
		TxPool:        txPool,
		ChainConfig:   *chain_config.GetChainConfig(),
		Broadcaster:   broadcaster,
	}

	service := VenusFullChainService{
		DipperinConfig: config,
		TxValidator:    fakeValidator{},
	}

	tx, err := service.MakeTmpSignedTx(CallArgs{}, 0)
	assert.NoError(t, err)
	sender, err := tx.Sender(nil)
	assert.NoError(t, err)
	assert.Equal(t, sender, address)
}

func TestMercuryFullChainService_SendTransactions(t *testing.T) {
	manager := createWalletManager(t)
	defer os.Remove(util.HomeDir() + testPath)
	account, err := manager.Wallets[0].Accounts()
	assert.NoError(t, err)

	address := account[0].Address
	pk, err := manager.Wallets[0].GetSKFromAddress(address)
	testAccount := tests.NewAccount(pk, address)
	testAccounts := []tests.Account{*testAccount}

	csChain := createCsChain(testAccounts)
	config := &DipperinConfig{
		WalletManager: manager,
		ChainReader:   csChain,
		TxPool:        createTxPool(csChain),
		ChainConfig:   *chain_config.GetChainConfig(),
	}

	service := VenusFullChainService{
		DipperinConfig: config,
		TxValidator:    fakeValidator{},
	}

	tx := model.RpcTransaction{
		To:       aliceAddr,
		Value:    big.NewInt(100),
		Nonce:    uint64(0),
		GasPrice: g_testData.TestGasPrice,
		GasLimit: g_testData.TestGasLimit,
	}

	// No error
	num, err := service.SendTransactions(address, []model.RpcTransaction{tx})
	assert.NoError(t, err)
	assert.Equal(t, 1, num)

	// tx pool AddLocals failed
	num, err = service.SendTransactions(address, []model.RpcTransaction{tx, tx})
	assert.Equal(t, "this transaction already in tx pool", err.Error())
	assert.Equal(t, 0, num)

	// Valid tx error
	service = VenusFullChainService{
		DipperinConfig: config,
		TxValidator:    fakeValidator{err: testErr},
	}
	num, err = service.SendTransactions(address, []model.RpcTransaction{tx})
	assert.Equal(t, testErr, err)
	assert.Equal(t, 0, num)

	// FindWalletFromAddress error
	num, err = service.SendTransactions(common.HexToAddress("123"), []model.RpcTransaction{tx})
	assert.Equal(t, accounts.ErrNotFindWallet, err)
	assert.Equal(t, 0, num)
}

func TestMercuryFullChainService_SendTransaction(t *testing.T) {
	manager := createWalletManager(t)
	defer os.Remove(util.HomeDir() + testPath)
	account, err := manager.Wallets[0].Accounts()
	assert.NoError(t, err)

	address := account[0].Address
	pk, err := manager.Wallets[0].GetSKFromAddress(address)
	testAccount := tests.NewAccount(pk, address)
	testAccounts := []tests.Account{*testAccount}

	serviceChain := createCsChainService(testAccounts)
	txPool := createTxPool(serviceChain.ChainState)
	serviceChain.TxPool = txPool

	// set abi
	WASMPath := g_testData.GetWASMPath("token", g_testData.CoreVmTestData)
	AbiPath := g_testData.GetAbiPath("token", g_testData.CoreVmTestData)
	tx := createContractTx(0, WASMPath, AbiPath, "DIPP,WU,10000", testAccounts)
	block := createBlock(serviceChain.ChainState, []*model.Transaction{tx}, nil)
	votes := createVerifiersVotes(block, serviceChain.ChainConfig.VerifierNumber, testAccounts)
	err = serviceChain.SaveBftBlock(block, votes)
	assert.NoError(t, err)

	broadcaster := chain_communication.NewBroadcastDelegate(txPool, fakeNodeConfig{}, fakePeerManager{}, serviceChain, fakePbftNode{})
	config := &DipperinConfig{
		NodeConf:      fakeNodeConfig{nodeType: chain_config.NodeTypeOfVerifier},
		WalletManager: manager,
		ChainReader:   serviceChain,
		TxPool:        txPool,
		ChainConfig:   *chain_config.GetChainConfig(),
		Broadcaster:   broadcaster,
	}

	service := VenusFullChainService{
		DipperinConfig: config,
		TxValidator:    fakeValidator{},
	}

	nonce := uint64(0)
	value := g_testData.TestValue
	hash, err := service.SendRegisterTransaction(address, value, g_testData.TestGasPrice, g_testData.TestGasLimit, &nonce)
	assert.NoError(t, err)
	assert.NotNil(t, hash)

	nonce = uint64(1)
	hash, err = service.SendCancelTransaction(address, g_testData.TestGasPrice, g_testData.TestGasLimit, &nonce)
	assert.NoError(t, err)
	assert.NotNil(t, hash)

	nonce = uint64(2)
	hash, err = service.SendUnStakeTransaction(address, g_testData.TestGasPrice, g_testData.TestGasLimit, &nonce)
	assert.NoError(t, err)
	assert.NotNil(t, hash)

	nonce = uint64(3)
	vote := &model.VoteMsg{}
	hash, err = service.SendEvidenceTransaction(address, aliceAddr, g_testData.TestGasPrice, g_testData.TestGasLimit, vote, vote, &nonce)
	assert.NoError(t, err)
	assert.NotNil(t, hash)

	nonce = uint64(4)
	hash, err = service.SendTransaction(address, aliceAddr, value, g_testData.TestGasPrice, g_testData.TestGasLimit, []byte{}, &nonce)
	assert.NoError(t, err)
	assert.NotNil(t, hash)

	nonce = uint64(5)
	hash, err = service.SendTransaction(common.Address{}, aliceAddr, value, g_testData.TestGasPrice, g_testData.TestGasLimit, []byte{}, &nonce)
	assert.Equal(t, "no default account in this node", err.Error())
	assert.Equal(t, common.Hash{}, hash)

	nonce = uint64(6)
	to := common.HexToAddress(common.AddressContractCreate)
	data, err := g_testData.GetCreateExtraData(WASMPath, AbiPath, "dipp,DIPP,10000")
	assert.NoError(t, err)
	gasLimit := g_testData.TestGasLimit * 100
	hash, err = service.SendTransactionContract(address, to, value, g_testData.TestGasPrice, gasLimit, data, &nonce)
	assert.NoError(t, err)

	nonce = uint64(7)
	to = cs_crypto.CreateContractAddress(address, uint64(0))
	data, err = g_testData.GetCallExtraData("getBalance", address.String())
	assert.NoError(t, err)
	hash, err = service.SendTransactionContract(address, to, value, g_testData.TestGasPrice, gasLimit, data, &nonce)
	assert.NoError(t, err)

	nonce = uint64(8)
	fs1 := model.NewSigner(big.NewInt(1))
	tx = model.NewTransaction(nonce, aliceAddr, value, g_testData.TestGasPrice, g_testData.TestGasLimit, []byte{})
	signedTx, _ := tx.SignTx(pk, fs1)
	hash, err = service.NewTransaction(*signedTx)
	assert.NoError(t, err)
	assert.NotNil(t, hash)

	hash, err = service.NewTransaction(*signedTx)
	assert.Equal(t, "this transaction already in tx pool", err.Error())
	assert.Equal(t, common.Hash{}, hash)
}

func TestMercuryFullChainService_SendTransaction_Error(t *testing.T) {
	manager := createWalletManager(t)
	defer os.Remove(util.HomeDir() + testPath)
	account, err := manager.Wallets[0].Accounts()
	assert.NoError(t, err)

	address := account[0].Address
	sk, err := manager.Wallets[0].GetSKFromAddress(address)
	testAccount := tests.NewAccount(sk, address)
	testAccounts := []tests.Account{*testAccount}
	csChain := createCsChainService(testAccounts)

	config := &DipperinConfig{
		NodeConf:      fakeNodeConfig{nodeType: chain_config.NodeTypeOfVerifier},
		WalletManager: manager,
		ChainReader:   csChain,
		ChainConfig:   *chain_config.GetChainConfig(),
		TxPool:        createTxPool(csChain.ChainState),
	}

	service := VenusFullChainService{
		DipperinConfig: config,
		TxValidator:    fakeValidator{err: testErr},
	}

	// signTxAndSend-valid error
	nonce := uint64(0)
	value := big.NewInt(100)
	hash, err := service.SendRegisterTransaction(address, value, g_testData.TestGasPrice, g_testData.TestGasLimit, &nonce)
	assert.Equal(t, testErr, err)
	assert.Equal(t, common.Hash{}, hash)

	hash, err = service.SendCancelTransaction(address, g_testData.TestGasPrice, g_testData.TestGasLimit, &nonce)
	assert.Equal(t, testErr, err)
	assert.Equal(t, common.Hash{}, hash)

	hash, err = service.SendUnStakeTransaction(address, g_testData.TestGasPrice, g_testData.TestGasLimit, &nonce)
	assert.Equal(t, testErr, err)
	assert.Equal(t, common.Hash{}, hash)

	vote := &model.VoteMsg{}
	hash, err = service.SendEvidenceTransaction(address, aliceAddr, g_testData.TestGasPrice, g_testData.TestGasLimit, vote, vote, &nonce)
	assert.Equal(t, testErr, err)
	assert.Equal(t, common.Hash{}, hash)

	hash, err = service.SendTransaction(address, aliceAddr, value, g_testData.TestGasPrice, g_testData.TestGasLimit, []byte{}, &nonce)
	assert.Equal(t, testErr, err)
	assert.Equal(t, common.Hash{}, hash)

	hash, err = service.NewTransaction(*createSignedTx(nonce, aliceAddr, value, []byte{}, testAccounts))
	assert.Equal(t, testErr, err)
	assert.Equal(t, common.Hash{}, hash)

	// getSendTxInfo error
	fakeAddr := common.HexToAddress("123")
	hash, err = service.SendRegisterTransaction(fakeAddr, value, g_testData.TestGasPrice, g_testData.TestGasLimit, &nonce)
	assert.Equal(t, accounts.ErrNotFindWallet, err)
	assert.Equal(t, common.Hash{}, hash)

	hash, err = service.SendCancelTransaction(fakeAddr, g_testData.TestGasPrice, g_testData.TestGasLimit, &nonce)
	assert.Equal(t, accounts.ErrNotFindWallet, err)
	assert.Equal(t, common.Hash{}, hash)

	hash, err = service.SendUnStakeTransaction(fakeAddr, g_testData.TestGasPrice, g_testData.TestGasLimit, &nonce)
	assert.Equal(t, accounts.ErrNotFindWallet, err)
	assert.Equal(t, common.Hash{}, hash)

	hash, err = service.SendEvidenceTransaction(fakeAddr, aliceAddr, g_testData.TestGasPrice, g_testData.TestGasLimit, vote, vote, &nonce)
	assert.Equal(t, accounts.ErrNotFindWallet, err)
	assert.Equal(t, common.Hash{}, hash)

	hash, err = service.SendTransaction(fakeAddr, aliceAddr, value, g_testData.TestGasPrice, g_testData.TestGasLimit, []byte{}, &nonce)
	assert.Equal(t, accounts.ErrNotFindWallet, err)
	assert.Equal(t, common.Hash{}, hash)

	// Type error
	config.NodeConf = fakeNodeConfig{nodeType: chain_config.NodeTypeOfMineMaster}
	service.DipperinConfig = config
	hash, err = service.SendRegisterTransaction(address, value, g_testData.TestGasPrice, g_testData.TestGasLimit, &nonce)
	assert.Equal(t, "the node isn't verifier", err.Error())
	assert.Equal(t, common.Hash{}, hash)

	hash, err = service.SendCancelTransaction(address, g_testData.TestGasPrice, g_testData.TestGasLimit, &nonce)
	assert.Equal(t, "the node isn't verifier", err.Error())
	assert.Equal(t, common.Hash{}, hash)

	hash, err = service.SendUnStakeTransaction(address, g_testData.TestGasPrice, g_testData.TestGasLimit, &nonce)
	assert.Equal(t, "the node isn't verifier", err.Error())
	assert.Equal(t, common.Hash{}, hash)

	hash, err = service.SendEvidenceTransaction(address, aliceAddr, g_testData.TestGasPrice, g_testData.TestGasLimit, vote, vote, &nonce)
	assert.Equal(t, "the node isn't verifier", err.Error())
	assert.Equal(t, common.Hash{}, hash)
}
