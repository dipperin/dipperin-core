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

package tests

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/log"
	"github.com/dipperin/dipperin-core/core/bloom"
	"github.com/dipperin/dipperin-core/core/chain"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/core/chain/chaindb"
	"github.com/dipperin/dipperin-core/core/chain/registerdb"
	"github.com/dipperin/dipperin-core/core/chain/state-processor"
	"github.com/dipperin/dipperin-core/core/economy-model"
	"github.com/dipperin/dipperin-core/core/model"
	model2 "github.com/dipperin/dipperin-core/core/vm/model"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"go.uber.org/zap"
	"math/big"
	"time"
)

type BftChainState interface {
	Chain
	SaveBftBlock(block model.AbstractBlock, seenCommits []model.AbstractVerification) error
}

type Chain interface {
	StateReader
	StateWriter
	VerifierHelper
	StateHelper
	ChainHelper
}

type StateWriter interface {
	SaveBlock(block model.AbstractBlock) error
}

type StateReader interface {
	Genesis() model.AbstractBlock
	CurrentBlock() model.AbstractBlock
	CurrentHeader() model.AbstractHeader
	GetBlock(hash common.Hash, number uint64) model.AbstractBlock
	GetBlockByHash(hash common.Hash) model.AbstractBlock
	GetBlockByNumber(number uint64) model.AbstractBlock
	HasBlock(hash common.Hash, number uint64) bool
	GetBody(hash common.Hash) model.AbstractBody
	GetBodyRLP(hash common.Hash) rlp.RawValue
	GetHeader(hash common.Hash, number uint64) model.AbstractHeader
	GetHeaderByHash(hash common.Hash) model.AbstractHeader
	GetHeaderByNumber(number uint64) model.AbstractHeader
	GetHeaderRLP(hash common.Hash) rlp.RawValue
	HasHeader(hash common.Hash, number uint64) bool
	GetBlockNumber(hash common.Hash) *uint64
	GetTransaction(txHash common.Hash) (model.AbstractTransaction, common.Hash, uint64, uint64)

	BlockProcessor(root common.Hash) (*chain.BlockProcessor, error)
	BlockProcessorByNumber(num uint64) (*chain.BlockProcessor, error)
}

type VerifierHelper interface {
	CurrentSeed() (common.Hash, uint64)
	IsChangePoint(block model.AbstractBlock, isProcessPackageBlock bool) bool
	GetLastChangePoint(block model.AbstractBlock) *uint64
	GetSlotByNum(num uint64) *uint64
	GetSlot(block model.AbstractBlock) *uint64
	GetCurrVerifiers() []common.Address
	GetVerifiers(round uint64) []common.Address
	GetNextVerifiers() []common.Address
	NumBeforeLastBySlot(slot uint64) *uint64
	BuildRegisterProcessor(preRoot common.Hash) (*registerdb.RegisterDB, error)
}

type StateHelper interface {
	GetStateStorage() state_processor.StateStorage
	CurrentState() (*state_processor.AccountStateDB, error)
	StateAtByBlockNumber(num uint64) (*state_processor.AccountStateDB, error)
	StateAtByStateRoot(root common.Hash) (*state_processor.AccountStateDB, error)
	BuildStateProcessor(preAccountStateRoot common.Hash) (*state_processor.AccountStateDB, error)
}

type ChainHelper interface {
	GetChainConfig() *chain_config.ChainConfig
	GetEconomyModel() economy_model.EconomyModel
	GetChainDB() chaindb.Database
}

type BlockBuilder struct {
	ChainState Chain
	PreBlock   model.AbstractBlock
	Txs        []*model.Transaction
	// commit list
	Vers          []model.AbstractVerification
	MinerPk       *ecdsa.PrivateKey
	InvalidTxList []model.AbstractTransaction
}

func (builder *BlockBuilder) SetVerifications(votes []model.AbstractVerification) {
	builder.Vers = votes
}

func (builder *BlockBuilder) SetPreBlock(block model.AbstractBlock) {
	builder.PreBlock = block
}

func (builder *BlockBuilder) SetMinerPk(pk *ecdsa.PrivateKey) {
	builder.MinerPk = pk
}

// build future block
func (builder *BlockBuilder) BuildFuture() model.AbstractBlock {
	coinbaseAddr := cs_crypto.GetNormalAddress(builder.MinerPk.PublicKey)
	if coinbaseAddr.IsEmpty() {
		panic("call NewBlockFromLastBlock, but coinbase address is empty")
	}
	curBlock := builder.PreBlock
	if curBlock == nil {
		panic("can't get current block when call NewBlockFromLastBlock")
	}

	curHeight := curBlock.Number()
	pubKey := &builder.MinerPk.PublicKey
	seed, proof := crypto.Evaluate(builder.MinerPk, builder.PreBlock.Seed().Bytes())

	header := &model.Header{
		Version:     curBlock.Version(),
		Number:      curHeight + 1,
		Seed:        seed,
		Proof:       proof,
		MinerPubKey: crypto.FromECDSAPub(pubKey),
		PreHash:     curBlock.Hash(),
		GasLimit:    chain_config.BlockGasLimit,

		// 一定要有，否则nonce和diff为空就会被判断成特殊块
		Diff:      builder.getDiff(),
		TimeStamp: big.NewInt(time.Now().Add(time.Second * 10).UnixNano()),
		CoinBase:  coinbaseAddr,
		Bloom:     iblt.NewBloom(model.DefaultBlockBloomConfig),
	}

	// set pre block verifications
	vers := builder.Vers
	pending := builder.getMappedTxs()

	// deal state
	processor, err := chain.NewBlockProcessor(builder.ChainState, curBlock.StateRoot(), builder.ChainState.GetStateStorage())
	if err != nil {
		panic("get state failed, err: " + err.Error())
	}

	txs := model.NewTransactionsByFeeAndNonce(nil, pending)
	txBuf, receipts := builder.commitTransactions(txs, processor, header, vers)

	var tmpTxs []*model.Transaction
	for _, tx := range txBuf {
		tmpTxs = append(tmpTxs, tx.(*model.Transaction))
	}

	if len(vers) == 0 && curHeight > 0 {
		panic(fmt.Sprintf("no verifications for height: %v", curHeight+1))
	}

	block := model.NewBlock(header, tmpTxs, vers)

	// calculate receipt hash
	receiptHash := model.DeriveSha(&receipts)
	block.SetReceiptHash(receiptHash)
	//block.SetBloomLog(model2.CreateBloom(receipts))
	//bloomLog := block.GetBloomLog()
	//log.DLogger.Info("BftBlockBuilder#BuildWaitPackBlock", "bloomLog", (&bloomLog).Hex(), "receipts", receipts, "bloomLogs2", fmt.Sprintf("%s", (&bloomLog).Hex()))

	linkList := model.NewInterLink(curBlock.GetInterlinks(), block)
	block.SetInterLinks(linkList)
	linkRoot := model.DeriveSha(linkList)
	block.SetInterLinkRoot(linkRoot)

	if err = processor.ProcessExceptTxs(block, builder.ChainState.GetEconomyModel(), true); err != nil {
		log.DLogger.Error("process state except txs failed", zap.Error(err))
		return nil
	}

	root, err := processor.Finalise()
	if err != nil {
		panic(err)
	}
	block.SetStateRoot(root)

	// deal register
	register, err := registerdb.NewRegisterDB(curBlock.GetRegisterRoot(), builder.ChainState.GetStateStorage(), builder.ChainState)
	if err = register.Process(block); err != nil {
		log.DLogger.Error("process register failed", zap.Error(err))
		return nil
	}
	registerRoot := register.Finalise()
	block.SetRegisterRoot(registerRoot)

	// calculate block nonce
	model.CalNonce(block)
	block.RefreshHashCache()
	log.DLogger.Info("calculate block nonce successful", zap.Uint64("num", block.Number()))
	return block
}

// build the wait-pack block
func (builder *BlockBuilder) Build() model.AbstractBlock {
	coinbaseAddr := cs_crypto.GetNormalAddress(builder.MinerPk.PublicKey)
	if coinbaseAddr.IsEmpty() {
		panic("call NewBlockFromLastBlock, but coinbase address is empty")
	}
	curBlock := builder.PreBlock
	if curBlock == nil {
		panic("can't get current block when call NewBlockFromLastBlock")
	}

	curHeight := curBlock.Number()
	pubKey := &builder.MinerPk.PublicKey
	seed, proof := crypto.Evaluate(builder.MinerPk, builder.PreBlock.Seed().Bytes())
	gasLimit := uint64(chain_config.BlockGasLimit)

	header := &model.Header{
		Version:     curBlock.Version(),
		Number:      curHeight + 1,
		Seed:        seed,
		Proof:       proof,
		MinerPubKey: crypto.FromECDSAPub(pubKey),
		PreHash:     curBlock.Hash(),

		//GasLimit:  &model.DefaultGasLimit,
		// 一定要有，否则nonce和diff为空就会被判断成特殊块
		Diff:      builder.getDiff(),
		TimeStamp: big.NewInt(time.Now().Add(time.Second * 3).UnixNano()),
		CoinBase:  coinbaseAddr,
		Bloom:     iblt.NewBloom(model.DefaultBlockBloomConfig),
		GasLimit:  gasLimit,
	}

	// set pre block verifications
	vers := builder.Vers
	pending := builder.getMappedTxs()

	// deal state
	processor, err := chain.NewBlockProcessor(builder.ChainState, curBlock.StateRoot(), builder.ChainState.GetStateStorage())
	if err != nil {
		panic("get state failed, err: " + err.Error())
	}

	txs := model.NewTransactionsByFeeAndNonce(nil, pending)
	txBuf, receipts := builder.commitTransactions(txs, processor, header, vers)

	var tmpTxs []*model.Transaction
	for _, tx := range txBuf {
		tmpTxs = append(tmpTxs, tx.(*model.Transaction))
	}

	if len(vers) == 0 && curHeight > 0 {
		panic(fmt.Sprintf("no verifications for height: %v", curHeight+1))
	}
	block := model.NewBlock(header, tmpTxs, vers)

	// calculate receipt hash
	receiptHash := model.DeriveSha(&receipts)
	block.SetReceiptHash(receiptHash)
	//block.SetBloomLog(model2.CreateBloom(receipts))
	//bloomLog := block.GetBloomLog()
	//log.DLogger.Info("BftBlockBuilder#BuildWaitPackBlock", "bloomLog", (&bloomLog).Hex(), "receipts", receipts, "bloomLogs2", fmt.Sprintf("%s", (&bloomLog).Hex()))

	linkList := model.NewInterLink(curBlock.GetInterlinks(), block)
	block.SetInterLinks(linkList)
	linkRoot := model.DeriveSha(linkList)
	block.SetInterLinkRoot(linkRoot)

	if err = processor.ProcessExceptTxs(block, builder.ChainState.GetEconomyModel(), true); err != nil {
		log.DLogger.Error("process state except txs failed", zap.Error(err))
		return nil
	}

	root, err := processor.Finalise()
	if err != nil {
		panic(err)
	}
	block.SetStateRoot(root)

	// deal register
	register, err := registerdb.NewRegisterDB(curBlock.GetRegisterRoot(), builder.ChainState.GetStateStorage(), builder.ChainState)
	if err = register.Process(block); err != nil {
		log.DLogger.Error("process register failed", zap.Error(err))
		return nil
	}
	registerRoot := register.Finalise()
	//log.DLogger.Info("set the register root is:","registerRoot",registerRoot.Hex(),"preRoot",curBlock.GetRegisterRoot().Hex())
	block.SetRegisterRoot(registerRoot)

	// calculate block nonce
	model.CalNonce(block)
	//refresh block hash
	block.RefreshHashCache()
	log.DLogger.Info("calculate block nonce successful", zap.Uint64("num", block.Number()))
	return block
}

// build special block
func (builder *BlockBuilder) BuildSpecialBlock() model.AbstractBlock {
	preBlock := builder.PreBlock
	pubKey := &builder.MinerPk.PublicKey
	seed, proof := crypto.Evaluate(builder.MinerPk, builder.PreBlock.Seed().Bytes())
	coinBaseAddr := cs_crypto.GetNormalAddress(builder.MinerPk.PublicKey)
	header := &model.Header{
		Version:     preBlock.Version(),
		Number:      preBlock.Number() + 1,
		Seed:        seed,
		Proof:       proof,
		MinerPubKey: crypto.FromECDSAPub(pubKey),
		PreHash:     preBlock.Hash(),
		Diff:        common.Difficulty{},
		TimeStamp:   big.NewInt(time.Now().Add(time.Second * 3).UnixNano()),
		CoinBase:    coinBaseAddr,
		Bloom:       iblt.NewBloom(model.DefaultBlockBloomConfig),
	}

	// set pre block verifications
	vers := builder.Vers

	// build block
	block := model.NewBlock(header, []*model.Transaction{}, vers)

	// calculate receipt hash
	var receipts model2.Receipts
	receiptHash := model.DeriveSha(&receipts)
	block.SetReceiptHash(receiptHash)
	//block.SetBloomLog(model2.CreateBloom(receipts))
	//bloomLog := block.GetBloomLog()
	//log.DLogger.Info("BftBlockBuilder#BuildWaitPackBlock", "bloomLog", (&bloomLog).Hex(), "receipts", receipts, "bloomLogs2", fmt.Sprintf("%s", (&bloomLog).Hex()))

	// set interlink root
	linkList := model.NewInterLink(preBlock.GetInterlinks(), block)
	block.SetInterLinks(linkList)
	linkRoot := model.DeriveSha(linkList)
	block.SetInterLinkRoot(linkRoot)

	// calculate state root
	processor, err := builder.ChainState.BlockProcessor(preBlock.StateRoot())
	if err = processor.ProcessExceptTxs(block, builder.ChainState.GetEconomyModel(), false); err != nil {
		log.DLogger.Error("process state failed", zap.Error(err))
	}

	root, err := processor.Finalise()
	if err != nil {
		log.DLogger.Error("finalise state failed", zap.Error(err))
	}
	block.SetStateRoot(root)

	// calculate register root
	registerPro, gErr := builder.ChainState.BuildRegisterProcessor(preBlock.GetRegisterRoot())
	if gErr != nil {
		log.DLogger.Error("get register processor failed", zap.Error(gErr))
	}

	if err = registerPro.Process(block); err != nil {
		log.DLogger.Error("process register failed", zap.Error(err))
	}
	registerRoot := registerPro.Finalise()
	block.SetRegisterRoot(registerRoot)
	block.RefreshHashCache()
	return block
}

func (builder *BlockBuilder) commitTransaction(conf *state_processor.TxProcessConfig, state *chain.BlockProcessor) error {
	snap := state.Snapshot()
	//err := state.ProcessTx(tx, height)
	/*	conf := state_processor.TxProcessConfig{
		Tx:      tx,
		TxIndex: txIndex,
		Header:  header,
		GetHash: state.GetBlockHashByNumber,
	}*/
	err := state.ProcessTxNew(conf)
	if err != nil {
		state.RevertToSnapshot(snap)
		return err
	}

	//updating tx fee
	conf.Tx.PaddingActualTxFee(conf.TxFee)
	return nil
}

/*func (builder *BlockBuilder) commitTransaction(tx model.AbstractTransaction, state *chain.BlockProcessor, txIndex int, header *model.Header) error {
	snap := state.Snapshot()
	//err := state.ProcessTx(tx, height)
	conf := state_processor.TxProcessConfig{
		Tx:      tx,
		TxIndex: txIndex,
		Header:  header,
		GetHash: state.GetBlockHashByNumber,
	}
	err := state.ProcessTxNew(&conf)
	if err != nil {
		state.RevertToSnapshot(snap)
		return err
	}
	return nil
}*/

func (builder *BlockBuilder) getDiff() common.Difficulty {
	if builder.PreBlock.Difficulty().Equal(common.Difficulty{}) {
		return common.HexToDiff("0x1fffffff")
	}
	return builder.PreBlock.Difficulty()
}

func (builder *BlockBuilder) commitTransactions(txs *model.TransactionsByFeeAndNonce, state *chain.BlockProcessor, header *model.Header, vers []model.AbstractVerification) (txBuf []model.AbstractTransaction, receipts model2.Receipts) {
	var invalidList []*model.Transaction
	gasLimit := header.GetGasLimit()
	gasUsed := header.GetGasUsed()
	for {
		// Retrieve the next transaction and abort if all done
		tx := txs.Peek()
		if tx == nil {
			break
		}
		//from, _ := tx.Sender(builder.nodeContext.TxSigner())
		conf := state_processor.TxProcessConfig{
			Tx:       tx,
			Header:   header,
			GetHash:  state.GetBlockHashByNumber,
			GasLimit: &gasLimit,
			GasUsed:  &gasUsed,
		}
		err := builder.commitTransaction(&conf, state)
		if err != nil {
			log.DLogger.Info("transaction is not processable because:", zap.Error(err), zap.Any("txID", tx.CalTxId()), zap.Uint64("nonce:", tx.Nonce()))
			txs.Pop()
			invalidList = append(invalidList, tx.(*model.Transaction))
		} else {
			receipt := tx.GetReceipt()
			if receipt == nil {
				log.DLogger.Info("empty receipt", zap.String("txId", tx.CalTxId().Hex()))
				txs.Pop()
				invalidList = append(invalidList, tx.(*model.Transaction))
			} else {
				txBuf = append(txBuf, tx)
				txs.Shift()
				receipts = append(receipts, receipt)
			}
		}
	}

	header.GasUsed = gasUsed
	// ProcessExceptTxs then finalise for fear that changing state root
	return
}

func (builder *BlockBuilder) getMappedTxs() map[common.Address][]model.AbstractTransaction {
	r := make(map[common.Address][]model.AbstractTransaction)
	for _, tx := range builder.Txs {
		if tx.Amount().Cmp(big.NewInt(0)) < 0 {
			builder.InvalidTxList = append(builder.InvalidTxList, tx)
			continue
		}
		sender, err := tx.Sender(nil)
		errPanic(err)
		r[sender] = append(r[sender], tx)
	}
	return r
}

func (builder *BlockBuilder) ClearInvalidTxList() {
	builder.InvalidTxList = []model.AbstractTransaction{}
}
