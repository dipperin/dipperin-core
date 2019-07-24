package transaction

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/consts"
	"github.com/dipperin/dipperin-core/common/hexutil"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/core/chain/state-processor"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/core/vm/common/utils"
	"github.com/dipperin/dipperin-core/tests/factory"
	"github.com/dipperin/dipperin-core/tests/g-testData"
	"github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
	"time"
)

func TestDebugTxRlp(t *testing.T) {
	txData, err := hexutil.Decode("0xf869e280960000970e8128ab834e8eac17ab8e3812f010678cf79180808203e80182a41080f844a04a6a9d72b370b8f3f7ff04a8224ee11aea5c993eb08e2851094c1aaf48d68527a07e1505cbdd7c7d25481065620a1288a2f6dc76e0aa156da970617caa148cc2c93880")
	assert.NoError(t, err)
	var transaction model.Transaction

	err = rlp.DecodeBytes(txData, &transaction)
	assert.NoError(t, err)

	log.Info("the tx is:", "transaction", transaction)

	log.Info("the tx extraData is:", "extraData", hexutil.Encode(transaction.ExtraData()))
}

func TestTxSize(t *testing.T) {
	keyAlice, _ := model.CreateKey()
	ms := model.NewMercurySigner(big.NewInt(1))
	tempTx := model.NewTransaction(uint64(0), factory.BobAddrV, big.NewInt(1000), g_testData.TestGasPrice, g_testData.TestGasLimit, []byte{})
	tempTx.SignTx(keyAlice, ms)
	log.Info("the tx size is:", "size", tempTx.Size())

	bytes, err := tempTx.EncodeRlpToBytes()
	assert.NoError(t, err)

	log.Info("the tx rlpBytes len is:", "len", len(bytes))
}

func TestCalculateMiniTxFee(t *testing.T) {
	//normal tx fee
	extraData := make([]byte, 0)
	for i := 0; i < 50*1024; i++ {
		extraData = append(extraData, byte(i%2))
	}

	log.Info("the extra data is:", "extraData", hexutil.Encode(extraData))
	tempTx := model.NewTransaction(uint64(0), factory.BobAddrV, big.NewInt(1000), g_testData.TestGasPrice, g_testData.TestGasLimit, extraData)
	keyAlice, _ := model.CreateKey()
	ms := model.NewMercurySigner(big.NewInt(1))
	tempTx.SignTx(keyAlice, ms)

	txData, err := tempTx.EncodeRlpToBytes()
	assert.NoError(t, err)

	log.Info("the txSize is:", "txSize", tempTx.Size(), "txRlpLen", len(txData))

	gasUsed, err := model.IntrinsicGas(extraData, false, false)
	assert.NoError(t, err)
	log.Info("the gasUsed is:", "gasUsed", gasUsed)
}

func createTestStateDB(addrInfo map[common.Address]*big.Int) (ethdb.Database, common.Hash) {
	db := ethdb.NewMemDatabase()

	//todo The new method does not take the tree from the underlying database
	tdb := state_processor.NewStateStorageWithCache(db)
	processor, _ := state_processor.NewAccountStateDB(common.Hash{}, tdb)

	for addr, balance := range addrInfo {
		processor.NewAccountState(addr)
		processor.AddBalance(addr, balance)
		processor.AddNonce(addr, 0)
	}
	root, _ := processor.Commit()
	tdb.TrieDB().Commit(root, false)
	return db, root
}

func createBlock(num uint64, preHash common.Hash, txList []*model.Transaction, limit uint64) *model.Block {
	header := model.NewHeader(1, num, preHash, common.HexToHash("123456"), common.HexToDiff("1fffffff"), big.NewInt(time.Now().UnixNano()), factory.BobAddrV, common.BlockNonce{})

	// vote
	var voteList []model.AbstractVerification
	header.GasLimit = limit
	block := model.NewBlock(header, txList, voteList)

	// calculate block nonce
	model.CalNonce(block)
	block.RefreshHashCache()
	return block
}

func TestWASMContactMiniTxFee(t *testing.T) {
	params := "dipp,DIPP,1000000"

	WASMTokenPath := g_testData.GetWASMPath("token-const",g_testData.CoreVmTestData)
	AbiTokenPath := g_testData.GetAbiPath("token-const",g_testData.CoreVmTestData)
	extraData, err := g_testData.GetCreateExtraData(WASMTokenPath, AbiTokenPath, params)
	extraData, err = utils.ParseCreateContractData(extraData)
	assert.NoError(t, err)

	to := common.HexToAddress(common.AddressContractCreate)
	value := big.NewInt(0)
	gasPrice := big.NewInt(1)
	gasLimit := big.NewInt(2 * consts.DIP)
	tempTx := model.NewTransactionSc(0, &to, value, gasPrice, gasLimit.Uint64(), extraData)

	keyAlice, _ := model.CreateKey()
	ms := model.NewMercurySigner(big.NewInt(1))
	tempTx.SignTx(keyAlice, ms)

	log.Info("the tx extra data size is:", "extraData size", len(tempTx.ExtraData()))

	//creat test stateDB
	sender := cs_crypto.GetNormalAddress(keyAlice.PublicKey)
	db, root := createTestStateDB(map[common.Address]*big.Int{sender: big.NewInt(0).Mul(big.NewInt(100), big.NewInt(consts.DIP))})
	processor, err := state_processor.NewAccountStateDB(root, state_processor.NewStateStorageWithCache(db))
	assert.NoError(t, err)

	//creat process config
	block := createBlock(1, common.Hash{}, []*model.Transaction{tempTx}, chain_config.MaxGasLimit)
	gasUsed := uint64(0)
	confGasLimit := gasLimit.Uint64()
	txConfigCreate := &state_processor.TxProcessConfig{
		Tx: tempTx,
		GetHash: func(number uint64) common.Hash {
			return common.Hash{}
		},
		Header:   block.Header(),
		GasLimit: &confGasLimit,
		GasUsed:  &gasUsed,
	}

	err = processor.ProcessTxNew(txConfigCreate)
	assert.NoError(t, err)

	receipt, err := txConfigCreate.Tx.GetReceipt()
	assert.NoError(t, err)

	log.Info("the contract tx gasUsed is:", "gasUsed", receipt.GasUsed)
	log.Info("the contract tx used TxFee is:", "txFee", txConfigCreate.Tx.(*model.Transaction).GetActualTxFee())

	log.Info("the contract tx size is: ", "size", txConfigCreate.Tx.Size())
}
