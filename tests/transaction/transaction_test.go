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
	"github.com/dipperin/dipperin-core/tests/vm"
	"github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
	"time"
)

func TestDebugTxRlp(t *testing.T){
	txData ,err:= hexutil.Decode("0xf902def90295019600120000000000000000000000000000000000000000808080806483989680b90271f9026eb8eb0061736d01000000010d0360017f0060027f7f00600000021d0203656e76067072696e7473000003656e76087072696e74735f6c00010304030202000405017001010105030100020615037f01419088040b7f00419088040b7f004186080b073405066d656d6f727902000b5f5f686561705f6261736503010a5f5f646174615f656e64030204696e697400030568656c6c6f00040a450302000b02000b3d01017f230041106b220124004180081000200141203a000f2001410f6a41011001200010002001410a3a000e2001410e6a41011001200141106a24000b0b0d01004180080b0668656c6c6f00b9017d5b0a202020207b0a2020202020202020226e616d65223a2022696e6974222c0a202020202020202022696e70757473223a205b5d2c0a2020202020202020226f757470757473223a205b5d2c0a202020202020202022636f6e7374616e74223a202266616c7365222c0a20202020202020202274797065223a202266756e6374696f6e220a202020207d2c0a202020207b0a2020202020202020226e616d65223a202268656c6c6f222c0a202020202020202022696e70757473223a205b0a2020202020202020202020207b0a20202020202020202020202020202020226e616d65223a20226e616d65222c0a202020202020202020202020202020202274797065223a2022737472696e67220a2020202020202020202020207d0a20202020202020205d2c0a2020202020202020226f757470757473223a205b5d2c0a202020202020202022636f6e7374616e74223a202274727565222c0a20202020202020202274797065223a202266756e6374696f6e220a202020207d0a5d0a80f844a094fdf6afa4600fcd86ceee9cb86c7edf3a70c6de61ccd100e7728d2c7f3a00d0a013062e19e9b265ad6ea25923f008d931b8c3df2f41d2d852a90cec30483d230b3980")
	assert.NoError(t,err)

	var transaction model.Transaction

	err = rlp.DecodeBytes(txData, &transaction)
	assert.NoError(t,err)

	log.Info("the tx is:","transaction",transaction)

	log.Info("the tx extraData is:","extraData",hexutil.Encode(transaction.ExtraData()))
}

func TestTxSize(t *testing.T){
	keyAlice, _ := model.CreateKey()
	ms := model.NewMercurySigner(big.NewInt(1))
	tempTx := model.NewTransaction(uint64(0), factory.BobAddrV, big.NewInt(1000), big.NewInt(10000), []byte{})
	tempTx.SignTx(keyAlice, ms)
	log.Info("the tx size is:","size",tempTx.Size())

	bytes,err := tempTx.EncodeRlpToBytes()
	assert.NoError(t,err)

	log.Info("the tx rlpBytes len is:","len",len(bytes))
}

func TestCalculateMiniTxFee(t *testing.T){
	//normal tx fee
	extraData := make([]byte,0)
	for i:=0;i<50*1024;i++{
		extraData = append(extraData,byte(i%2))
	}

	log.Info("the extra data is:","extraData",hexutil.Encode(extraData))
	tempTx := model.NewTransaction(uint64(0), factory.BobAddrV, big.NewInt(1000), big.NewInt(10000), extraData)
	keyAlice, _ := model.CreateKey()
	ms := model.NewMercurySigner(big.NewInt(1))
	tempTx.SignTx(keyAlice, ms)

	txData,err := tempTx.EncodeRlpToBytes()
	assert.NoError(t,err)

	log.Info("the txSize is:","txSize",tempTx.Size(),"txRlpLen",len(txData))

	gasUsed,err := model.IntrinsicGas(extraData,false,false)
	assert.NoError(t,err)
	log.Info("the gasUsed is:","gasUsed",gasUsed)
}

func createTestStateDB(addrInfo map[common.Address]*big.Int) (ethdb.Database, common.Hash){
	db := ethdb.NewMemDatabase()

	//todo The new method does not take the tree from the underlying database
	tdb := state_processor.NewStateStorageWithCache(db)
	processor, _ := state_processor.NewAccountStateDB(common.Hash{}, tdb)

	for addr,balance := range addrInfo{
		processor.NewAccountState(addr)
		processor.AddBalance(addr,balance)
		processor.AddNonce(addr,0)
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

func TestWASMContactMiniTxFee(t *testing.T){
	params := "dipp,DIPP,1000000"
	extraData := vm.GetCreateExtraData(t, vm.WASMTokenPath, vm.AbiTokenPath, params)

	extraData, err := utils.ParseCreateContractData(extraData)
	assert.NoError(t,err)

	to := common.HexToAddress(common.AddressContractCreate)
	value := big.NewInt(0)
	gasPrice := big.NewInt(1)
	gasLimit := big.NewInt(2 * consts.DIP)
	tempTx := model.NewTransactionSc(0, &to, value, gasPrice, gasLimit.Uint64(), extraData)

	keyAlice, _ := model.CreateKey()
	ms := model.NewMercurySigner(big.NewInt(1))
	tempTx.SignTx(keyAlice, ms)

	log.Info("the tx extra data size is:","extraData size",len(tempTx.ExtraData()))

	//creat test stateDB
	sender :=  cs_crypto.GetNormalAddress(keyAlice.PublicKey)
	db, root := createTestStateDB(map[common.Address]*big.Int{sender:big.NewInt(100*consts.DIP)})
	processor, err := state_processor.NewAccountStateDB(root, state_processor.NewStateStorageWithCache(db))
	assert.NoError(t, err)

	//creat process config
	block := createBlock(1, common.Hash{}, []*model.Transaction{tempTx},chain_config.MaxGasLimit)
	tempTx.PaddingTxIndex(0)
	gasUsed := uint64(0)
	confGasLimit := gasLimit.Uint64()
	txConfigCreate := &state_processor.TxProcessConfig{
		Tx:       tempTx,
		GetHash:  func (number uint64) common.Hash{
			return common.Hash{}
		},
		Header:block.Header(),
		GasLimit: &confGasLimit,
		GasUsed:  &gasUsed,
	}

	err = processor.ProcessTxNew(txConfigCreate)
	assert.NoError(t,err)

	receipt ,err:= txConfigCreate.Tx.GetReceipt()
	assert.NoError(t,err)

	log.Info("the contract tx gasUsed is:","gasUsed",receipt.GasUsed)
	log.Info("the contract tx used TxFee is:","txFee",txConfigCreate.Tx.(*model.Transaction).GetActualTxFee())

	log.Info("the contract tx size is: ","size",txConfigCreate.Tx.Size())
}









