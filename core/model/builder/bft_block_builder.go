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

package builder

import (
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/g-metrics"
	"github.com/dipperin/dipperin-core/core/accounts"
	"github.com/dipperin/dipperin-core/core/bloom"
	"github.com/dipperin/dipperin-core/core/chain"
	"github.com/dipperin/dipperin-core/core/chain-communication"
	"github.com/dipperin/dipperin-core/core/chain/state-processor"

	"github.com/dipperin/dipperin-core/core/model"
	model2 "github.com/dipperin/dipperin-core/core/vm/model"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
	"github.com/dipperin/dipperin-core/third-party/log"
	"math/big"
	"time"
)

// context must have chainReader state_processor.ChainReader, stateProcessorBuilder stateProcessorBuilder, accountStorage state_processor.StateStorage, txPool txPool
func MakeBftBlockBuilder(config ModelConfig) *BftBlockBuilder {
	return &BftBlockBuilder{
		ModelConfig: config,
	}
}

type BftBlockBuilder struct {
	//nodeContext NodeContext
	ModelConfig
}

func (builder *BftBlockBuilder) GetMsgSigner() chain_communication.PbftSigner {
	return builder.MsgSigner
}

func (builder *BftBlockBuilder) SetMsgSigner(MsgSigner chain_communication.PbftSigner) {
	builder.MsgSigner = MsgSigner
}

func (builder *BftBlockBuilder) commitTransaction(conf *state_processor.TxProcessConfig, state *chain.BlockProcessor) error {
	snap := state.Snapshot()
	err := state.ProcessTxNew(conf)
	if err != nil {
		state.RevertToSnapshot(snap)
		return err
	}

	return nil
}

func (builder *BftBlockBuilder) commitTransactions(txs *model.TransactionsByFeeAndNonce, state *chain.BlockProcessor, header *model.Header, vers []model.AbstractVerification) (txBuf []model.AbstractTransaction, receipts model2.Receipts) {
	var invalidList []*model.Transaction
	log.Info("BftBlockBuilder#commitTransactions  start ~~~~~++")
	gasUsed := uint64(0)
	gasLimit := header.GasLimit
	for {
		// Retrieve the next transaction and abort if all done
		tx := txs.Peek()
		if tx == nil {
			break
		}
		log.Info("BftBlockBuilder#commitTransactions ", "tx hash", tx.CalTxId())
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
			log.Info("transaction is not processable because:", "err", err, "txID", tx.CalTxId(), "nonce:", tx.Nonce())
			txs.Pop()
			invalidList = append(invalidList, tx.(*model.Transaction))
		} else {
			receipt := tx.GetReceipt()
			if receipt == nil {
				log.Info("cant get tx receipt", "txId", tx.CalTxId().Hex())
				txs.Pop()
				invalidList = append(invalidList, tx.(*model.Transaction))
			} else {
				txBuf = append(txBuf, tx)
				txs.Shift()
				receipts = append(receipts, receipt)
			}
		}
	}

	//if there are invalid txs remove them from tx pool
	if len(invalidList) != 0 {
		block := model.NewBlock(header, invalidList, vers)
		builder.TxPool.RemoveTxs(block)
		log.Info("remove invalid Txs from pool", "num of txs", len(invalidList))
	}

	//update gasUsed in header
	header.GasUsed = gasUsed
	// ProcessExceptTxs then finalise for fear that changing state root
	return
}

//build the wait-pack block
func (builder *BftBlockBuilder) BuildWaitPackBlock(coinbaseAddr common.Address, gasFloor, gasCeil uint64) model.AbstractBlock {
	//trace pack block duration
	timer:=g_metrics.NewTimer(g_metrics.PackageBlockDuration)
	defer timer.ObserveDuration()

	if coinbaseAddr.IsEmpty() {
		panic("call NewBlockFromLastBlock, but coinbase address is empty")
	}
	curBlock := builder.ChainReader.CurrentBlock()
	if curBlock == nil {
		panic("can't get current block when call NewBlockFromLastBlock")
	}
	curHeight := curBlock.Number()
	log.PBft.Debug("build wait pack block", "height", curHeight+1)

	pubKey := builder.MsgSigner.PublicKey()
	coinbaseAddr = cs_crypto.GetNormalAddress(*pubKey)

	account := accounts.Account{Address: coinbaseAddr}
	seed, proof, err := builder.MsgSigner.Evaluate(account, curBlock.Seed().Bytes())
	if err != nil {
		log.Error("VRF seed failed")
		return nil
	}

	lastNormalBlock := builder.ChainReader.GetLatestNormalBlock()
	tmpValue := CalcGasLimit(lastNormalBlock.(*model.Block), gasFloor, gasCeil)
	log.Info("build block", "gasFloor", gasFloor, "gasCeil", gasCeil, "newGasLimit", tmpValue, "coinBase", coinbaseAddr)

	header := &model.Header{
		Version:     curBlock.Version(),
		Number:      curHeight + 1,
		Seed:        seed,
		Proof:       proof,
		MinerPubKey: crypto.FromECDSAPub(pubKey),
		PreHash:     curBlock.Hash(),
		Diff:        builder.GetDifficulty(),
		TimeStamp:   big.NewInt(time.Now().Add(time.Second * 3).UnixNano()),
		CoinBase:    coinbaseAddr,
		// TODO:
		Bloom:    iblt.NewBloom(model.DefaultBlockBloomConfig),
		GasLimit: tmpValue,
	}

	// set pre block verifications
	vers := builder.ChainReader.GetSeenCommit(curHeight)
	for index, v := range vers {
		if v.GetBlockHash() != curBlock.Hash().Hex() {
			panic(fmt.Sprintf("build block, but vote hash not match, index: %v, b hash: %v, vb hash: %v", index, curBlock.Hash().Hex(), v.GetBlockHash()))
		}
	}

	pending, err := builder.TxPool.Pending()
	if err != nil {
		log.Error("Failed to fetch pending transactions", "err", err)
		return nil
	}

	// deal state
	processor, err := builder.ChainReader.BlockProcessor(curBlock.StateRoot())
	//processor, err := builder.BuildStateProcessor.BuildStateProcessor(curBlock.StateRoot(), builder.ChainReader, builder.StateStorage)

	log.Info("~~~~~~~~~~~~~~~the pending len is:", "number", len(pending))
	txs := model.NewTransactionsByFeeAndNonce(builder.TxSigner, pending)
	txBuf, receipts := builder.commitTransactions(txs, processor, header, vers)

	//log.Info("~~~~~~~~~~~~~~ the txBuf len is: ", "txBuf Len", len(txBuf))

	var tmpTxs []*model.Transaction
	//var logs []*model2.Log
	for _, tx := range txBuf {
		//log.Info("the packaged tx is:", "txId", tx.CalTxId().Hex())
		tmpTxs = append(tmpTxs, tx.(*model.Transaction))

	}

	if len(vers) == 0 && curHeight > 0 {
		log.Warn("can't load pre verifications for height", "new height", curHeight+1)
		return nil
	}
	log.Info("build bft block1", "vers", len(vers), "height", curHeight)

	block := model.NewBlock(header, tmpTxs, vers)
	if block.Number() == 1 && !block.VerificationRoot().IsEqual(model.EmptyVerfRoot) {
		panic(fmt.Sprintf("invalid v root: %v", block.VerificationRoot()))
	}
	log.Info("build bft block2", "vers", len(block.GetVerifications()))

	//calculate receipt hash
	receiptHash := model.DeriveSha(&receipts)
	block.SetReceiptHash(receiptHash)
	//block.SetBloomLog(model2.CreateBloom(receipts))
	//bloomLog := block.GetBloomLog()
	//log.Info("BftBlockBuilder#BuildWaitPackBlock", "bloomLog", (&bloomLog).Hex(), "receipts", receipts, "bloomLogs2", fmt.Sprintf("%s", (&bloomLog).Hex()))

	//TODO:calculate block interlinks
	linkList := model.NewInterLink(curBlock.GetInterlinks(), block)
	block.SetInterLinks(linkList)
	//TODO:calculate RLP
	linkRoot := model.DeriveSha(linkList)
	block.SetInterLinkRoot(linkRoot)

	//Txs processed in commitTransactions(),no use process()
	log.Info("the processor is:", "processor", processor)
	if err = processor.ProcessExceptTxs(block, builder.ChainReader.GetEconomyModel(), true); err != nil {
		log.Error("process state except txs failed", "err", err)
		return nil
	}

	root, err := processor.Finalise()
	if err != nil {
		panic(err)
	}

	log.Info("the block build calculated block stateRoot is:", "blockNumber", block.Number(), "stateRoot", root.Hex())
	block.SetStateRoot(root)
	log.PBft.Debug("build block", "preBlock root", curBlock.StateRoot().Hex(), "process result", root.Hex(), "this block", block.StateRoot())

	// deal register
	register, err := builder.ChainReader.BuildRegisterProcessor(curBlock.GetRegisterRoot())
	if err = register.Process(block); err != nil {
		log.Error("process register failed", "err", err)
		return nil
	}
	registerRoot := register.Finalise()
	block.SetRegisterRoot(registerRoot)
	log.PBft.Debug("build block", "block id", block.Hash().Hex(), "transaction", block.TxCount())
	return block
}

//func (builder *DefaultBlockBuilder) NewBlockFromLastBlock(coinbaseAddr common.Address) model.AbstractBlock {
//	if coinbaseAddr.IsEmpty() {
//		panic("call NewBlockFromLastBlock, but coinbase address is empty")
//	}
//	txs := builder.txPool.RandTxsForPack()
//	// this builder only build default tx
//	var tmpTxs []*model.Transaction
//	for _, tx := range txs {
//		tmpTxs = append(tmpTxs, tx.(*model.Transaction))
//	}
//	curBlock := builder.chainDB.CurrentBlock()
//	if curBlock == nil {
//		panic("can't get current block when call NewBlockFromLastBlock")
//	}
//
//	header := &model.Header{
//		Version: curBlock.Version(),
//		Number: curBlock.Number() + 1,
//		PreHash: curBlock.Hash(),
//		//FIXME add a right seed.
//		Seed: curBlock.Hash(),
//		Diff: builder.GetDifficulty(),
//		CoinBase: coinbaseAddr,
//		//StateRoot: ,
//		// TODO:
//		//TimeStamp: ,
//		//Bloom: ,
//		//HeaderRoot: ,
//		//VerificationRoot: ,
//	}
//	// TODO: add v msg
//	block := model.NewBlock(header, tmpTxs[:], nil)
//
//	// process state root
//	processor, err := builder.stateProcessorBuilder(curBlock.StateRoot(), builder.chainDB, builder.accountStorage)
//	if err != nil {
//		log.Error("can't create state processor", "err", err)
//		return nil
//	}
//	if _, err := processor.Process(block); err != nil {
//		log.Warn("pack block builder process block err", "err", err)
//		return nil
//	}
//	roots := processor.Finalise()
//	root := roots[state_processor.AccountStateRootKey]
//	if root.IsEmpty() {
//		log.Warn("pack block builder process block got empty state root")
//		return nil
//	}
//
//	block.SetStateRoot(root)
//	return block
//}

func (builder *BftBlockBuilder) GetDifficulty() common.Difficulty {
	chainReader := builder.ChainReader

	curBlock := chainReader.CurrentBlock()
	lastPNum := model.LastPeriodBlockNum(curBlock.Number())
	if lastPNum == 0 {
		lastPNum = 1
	}

	log.Debug("call GetDifficulty for mine, get lastPeriodBlock", "num", lastPNum)
	lastPeriodBlock := chainReader.GetBlockByNumber(lastPNum)
	if curBlock == nil {
		panic("mine master get difficulty error,block is nil")
	}

	//find the neighbor normal block
	findBlock := chainReader.GetLatestNormalBlock()

	diff := model.NewCalNewWorkDiff(lastPeriodBlock, findBlock, curBlock.Number())

	log.Debug("mine master difficulty", "diff", diff.Hex())
	return diff
}

// CalcGasLimit computes the gas limit of the next block after parent. It aims
// to keep the baseline gas above the provided floor, and increase it towards the
// ceil if the blocks are full. If the ceil is exceeded, it will always decrease
// the gas allowance.
//　以parentGasUsed > parentGasLimit * (2/3)为判断，若大于则contrib-decay正，limit增加，
//  若小于则取决于差值有多大，若非常接近2/3则会增加1. 与此同时，gas limit的调整量在上一个块limit的１/1024之间
//  gas limit要大于系统定义最小limit:5000. 与此同时其要介于传入的gasFloor 和gasCeil之间．增量和减量都是以decay来修改的
//  gasFloor和gasCeil由矿工进行设置，因此调节会趋近矿工设置范围区块
//  共识处需要对gas limit进行检查，查看其修改量是否超出上一个limit的1/1024.并且其介于系统设置最大最小值之间．
func CalcGasLimit(parent *model.Block, gasFloor, gasCeil uint64) uint64 {
	// contrib = (parentGasUsed * 3 / 2) / 1024
	contrib := (parent.GasUsed() + parent.GasUsed()/2) / model2.GasLimitBoundDivisor

	// decay = parentGasLimit / 1024 -1
	decay := parent.GasLimit()/model2.GasLimitBoundDivisor - 1

	log.Info("the contrib and decay is:", "contrib", contrib, "decay", decay)
	/*
		strategy: gasLimit of block-to-mine is set based on parent's
		gasUsed value.  if parentGasUsed > parentGasLimit * (2/3) then we
		increase it, otherwise lower it (or leave it unchanged if it's right
		at that usage) the amount increased/decreased depends on how far away
		from parentGasLimit * (2/3) parentGasUsed is.
	*/
	limit := parent.GasLimit() - decay + contrib
	if limit < model2.MinGasLimit {
		limit = model2.MinGasLimit
	}

	log.Info("the limit after change is:", "limit", limit)
	// If we're outside our allowed gas range, we try to hone towards them
	if limit < gasFloor {
		limit = parent.GasLimit() + decay
		if limit > gasFloor {
			limit = gasFloor
		}
	} else if limit > gasCeil {
		limit = parent.GasLimit() - decay
		if limit < gasCeil {
			limit = gasCeil
		}
	}
	return limit
}
