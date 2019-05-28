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
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/accounts"
	"github.com/dipperin/dipperin-core/core/bloom"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/dipperin/dipperin-core/third-party/log/pbft_log"
	"time"
	"math/big"
	"fmt"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"github.com/dipperin/dipperin-core/core/chain"
	"github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
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

func (builder *BftBlockBuilder) commitTransaction(tx model.AbstractTransaction, state *chain.BlockProcessor, height uint64) error {
	snap := state.Snapshot()
	err := state.ProcessTx(tx, height)
	if err != nil {
		state.RevertToSnapshot(snap)
		return err
	}
	return nil
}

func (builder *BftBlockBuilder) commitTransactions(txs *model.TransactionsByFeeAndNonce, state *chain.BlockProcessor, header *model.Header, vers []model.AbstractVerification) (txBuf []model.AbstractTransaction) {
	var invalidList []*model.Transaction
	for {
		// Retrieve the next transaction and abort if all done
		tx := txs.Peek()
		if tx == nil {
			break
		}
		//from, _ := tx.Sender(builder.nodeContext.TxSigner())
		err := builder.commitTransaction(tx, state, header.Number)
		if err != nil {
			log.Info("transaction is not processable because:", "err", err, "txID", tx.CalTxId(), "nonce:", tx.Nonce())
			txs.Pop()
			invalidList = append(invalidList, tx.(*model.Transaction))
		} else {
			txBuf = append(txBuf, tx)
			txs.Shift()
		}
	}

	//if there are invalid txs remove them from tx pool
	if len(invalidList) != 0 {
		block := model.NewBlock(header, invalidList, vers)
		builder.TxPool.RemoveTxs(block)
		log.Info("remove invalid Txs from pool", "num of txs", len(invalidList))
	}

	// ProcessExceptTxs then finalise for fear that changing state root
	return
}

//build the wait-pack block
func (builder *BftBlockBuilder) BuildWaitPackBlock(coinbaseAddr common.Address) model.AbstractBlock {
	if coinbaseAddr.IsEmpty() {
		panic("call NewBlockFromLastBlock, but coinbase address is empty")
	}
	curBlock := builder.ChainReader.CurrentBlock()
	if curBlock == nil {
		panic("can't get current block when call NewBlockFromLastBlock")
	}
	curHeight := curBlock.Number()
	pbft_log.Debug("build wait pack block", "height", curHeight+1)

	pubKey := builder.MsgSigner.PublicKey()
	coinbaseAddr = cs_crypto.GetNormalAddress(*pubKey)

	account := accounts.Account{Address: coinbaseAddr}
	seed, proof, err := builder.MsgSigner.Evaluate(account, curBlock.Seed().Bytes())
	if err != nil {
		log.Error("VRF seed failed")
		return nil
	}

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
		Bloom: iblt.NewBloom(model.DefaultBlockBloomConfig),
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
	txBuf := builder.commitTransactions(txs, processor, header, vers)

	//log.Info("~~~~~~~~~~~~~~ the txBuf len is: ", "txBuf Len", len(txBuf))

	var tmpTxs []*model.Transaction
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
	pbft_log.Debug("build block", "preBlock root", curBlock.StateRoot().Hex(), "process result", root.Hex(), "this block", block.StateRoot())

	// deal register
	register, err := builder.ChainReader.BuildRegisterProcessor(curBlock.GetRegisterRoot())
	if err = register.Process(block); err != nil {
		log.Error("process register failed", "err", err)
		return nil
	}
	registerRoot := register.Finalise()
	block.SetRegisterRoot(registerRoot)
	pbft_log.Debug("build block", "block id", block.Hash().Hex(), "transaction", block.TxCount())
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
