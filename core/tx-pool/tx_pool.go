// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package tx_pool

import (
	"errors"
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/log"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/core/bloom"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/core/chain/state-processor"
	"github.com/dipperin/dipperin-core/core/model"
	"go.uber.org/zap"
	"math/big"
	"sync"
	"time"

	"github.com/dipperin/dipperin-core/common/g-timer"
	"gopkg.in/karalabe/cookiejar.v2/collections/prque"
	"math"
	"sort"
)

//go:generate mockgen -destination=./block_chain_mock_test.go -package=tx_pool github.com/dipperin/dipperin-core/core/tx-pool BlockChain
type BlockChain interface {
	CurrentBlock() model.AbstractBlock
	GetBlockByNumber(number uint64) model.AbstractBlock
	StateAtByStateRoot(root common.Hash) (*state_processor.AccountStateDB, error)
}

// TODO：cannot package multiple election transaction in one round
// TxPoolConfig are the configuration parameters of the transaction pool.
type TxPoolConfig struct {
	NoLocals  bool          // Whether local transaction handling should be disabled
	Journal   string        // Journal of local transactions to survive node restarts
	Rejournal time.Duration // Time interval to regenerate the local transaction journal

	MinFee  *big.Int // Minimum fee to enforce for acceptance into the pool
	FeeBump uint64   // Minimum fee bump percentage to replace an already existing transaction (nonce)

	AccountSlots uint64 // Number of executable transaction slots guaranteed per account
	GlobalSlots  uint64 // Maximum number of executable transaction slots for all accounts
	AccountQueue uint64 // Maximum number of non-executable transaction slots permitted per account
	GlobalQueue  uint64 // Maximum number of non-executable transaction slots for all accounts

	Lifetime time.Duration // Maximum amount of time non-executable transaction are queued
}

// DefaultTxPoolConfig contains the default configurations for the transaction
// pool.

var DefaultTxPoolConfig = TxPoolConfig{
	//check tx validation when add local txs
	NoLocals:  true,
	Journal:   "transaction.rlp",
	Rejournal: time.Hour,

	FeeBump: 1,

	AccountSlots: 1024,
	GlobalSlots:  1024 * 20,
	AccountQueue: 1024 * 100,
	GlobalQueue:  1024 * 100,

	Lifetime: 3 * time.Hour,
}

type TxPool struct {
	//PoolConsensus consensus.TransactionValidator
	config      TxPoolConfig
	chainConfig chain_config.ChainConfig
	chain       BlockChain
	signer      model.Signer
	minFee      *big.Int

	mu sync.RWMutex

	currentState *state_processor.AccountStateDB
	pendingState *state_processor.ManagedState

	locals  *accountSet // Set of local transaction to exempt from eviction rules
	journal *txJournal  // Journal of local transaction to back up to disk

	pending map[common.Address]*txList   // All currently processable transactions
	queue   map[common.Address]*txList   // Queued but non-processable transactions
	beats   map[common.Address]time.Time // Last heartbeat from each known account
	all     *txLookup                    // All transactions to allow lookups
	feeList *txFeeList                   // All transactions sorted by price

	loopStopCtrl chan int
	wg           sync.WaitGroup  // for shutdown sync
	senderCacher *model.TxCacher //transaction sender/id concurrent calculator
	//waitPackTxs []model.AbstractTransaction
}

var (
	evictionInterval = time.Minute // Time interval to check for evictable transactions
)

func (pool *TxPool) Reset(oldHead, newHead *model.Header) {
	pool.mu.Lock()
	defer pool.mu.Unlock()
	pool.reset(oldHead, newHead)
}

// TODO
func (pool *TxPool) reset(oldHead, newHead *model.Header) {
	// If we're reorging an old state, reinject all dropped transactions
	var reinject []model.AbstractTransaction

	if oldHead != nil && oldHead.Hash() != newHead.PreHash {
		// If the reorg is too deep, avoid doing it (will happen during fast sync)
		oldNum := oldHead.Number
		newNum := newHead.Number

		if depth := uint64(math.Abs(float64(oldNum) - float64(newNum))); depth > 64 {
			log.DLogger.Debug("Skipping deep transaction reorg", zap.Uint64("depth", depth))
		} else {
			// Reorg seems shallow enough to pull in all transactions into memory
			var discarded, included []model.AbstractTransaction

			var (
				rem = pool.chain.GetBlockByNumber(oldHead.Number)
				add = pool.chain.GetBlockByNumber(newHead.Number)
			)
			for rem.Number() > add.Number() {
				discarded = append(discarded, rem.GetAbsTransactions()...)
				if rem = pool.chain.GetBlockByNumber(rem.Number() - 1); rem == nil {
					log.DLogger.Error("Unrooted old chain seen by tx pool", zap.Uint64("block", oldHead.Number), zap.Any("hash", oldHead.Hash()))
					return
				}
			}
			for add.Number() > rem.Number() {
				included = append(included, add.GetAbsTransactions()...)
				if add = pool.chain.GetBlockByNumber(add.Number() - 1); add == nil {
					log.DLogger.Error("Unrooted new chain seen by tx pool", zap.Uint64("block", newHead.Number), zap.Any("hash", newHead.Hash()))
					return
				}
			}
			for rem.Hash() != add.Hash() {
				discarded = append(discarded, rem.GetAbsTransactions()...)
				if rem = pool.chain.GetBlockByNumber(rem.Number() - 1); rem == nil {
					log.DLogger.Error("Unrooted old chain seen by tx pool", zap.Uint64("block", oldHead.Number), zap.Any("hash", oldHead.Hash()))
					return
				}
				included = append(included, add.GetAbsTransactions()...)
				if add = pool.chain.GetBlockByNumber(add.Number() - 1); add == nil {
					log.DLogger.Error("Unrooted new chain seen by tx pool", zap.Uint64("block", newHead.Number), zap.Any("hash", newHead.Hash()))
					return
				}
			}
			reinject = model.TxDifference(discarded, included)
		}
	}
	// Initialize the internal state to the current head
	if newHead == nil {
		newHead = pool.chain.CurrentBlock().Header().(*model.Header) // Special case during testing
	}
	statedb, err := pool.chain.StateAtByStateRoot(newHead.StateRoot)
	if err != nil {
		log.DLogger.Error("Failed to reset txpool state", zap.Error(err))
		return
	}
	log.DLogger.Info("TxPool reset stateDb")
	pool.currentState = statedb
	pool.pendingState = state_processor.ManageState(statedb)

	// Inject any transactions discarded due to reorgs
	log.DLogger.Debug("Reinjecting stale transactions", zap.Int("count", len(reinject)))
	//senderCacher.recover(pool.signer, reinject)
	// TODO  to understand this
	pool.senderCacher.TxRecover(reinject)
	pool.addTxsLocked(reinject, false)

	// validate the pool of pending transactions, this will remove
	// any transactions that have been included in the block or
	// have been invalidated because of another transaction (e.g.
	// higher gas price)
	pool.demoteUnexecutables()

	// Update all accounts to the latest known pending nonce
	for addr, list := range pool.pending {
		txs := list.Flatten() // Heavy but will be cached and is needed by the miner anyway
		pool.pendingState.SetNonce(addr, txs[len(txs)-1].Nonce()+1)
	}
	// Check the queue and move transactions over to the pending if possible
	// or remove those that have become invalid
	pool.promoteExecutables(nil)
}

func (pool *TxPool) TxsCaching(txs []*model.Transaction) {
	aTxs := make([]model.AbstractTransaction, len(txs))
	util.InterfaceSliceCopy(aTxs, txs)
	pool.senderCacher.TxRecover(aTxs)
}

func (pool *TxPool) AbstractTxsCaching(txs []model.AbstractTransaction) {
	pool.senderCacher.TxRecover(txs)
}

func (pool *TxPool) Start() error {
	log.DLogger.Debug("---------------------------------------transaction pool start")
	pool.wg.Add(1)
	go pool.loop()
	return nil
}

func (pool *TxPool) Stop() {
	log.DLogger.Debug("------------------------------------transaction pool stop")
	pool.loopStopCtrl <- 1
	//sync, wait loop stop
	pool.wg.Wait()
}

// todo use reset
func (pool *TxPool) loop() {
	defer pool.wg.Done()

	evictHandler := func() {
		// Handle inactive account transaction eviction
		pool.mu.Lock()
		for addr := range pool.queue {
			// Skip local transactions from the eviction mechanism
			if pool.locals.contains(addr) {
				continue
			}
			// Any non-locals old enough should be removed
			if time.Since(pool.beats[addr]) > pool.config.Lifetime {
				for _, tx := range pool.queue[addr].Flatten() {
					pool.removeTx(tx.CalTxId(), true)
				}
			}
		}
		pool.mu.Unlock()
	}
	evict := g_timer.SetPeriodAndRun(evictHandler, evictionInterval)
	defer g_timer.StopWork(evict)

	journalHandler := func() {
		// Handle local transaction journal rotation
		if pool.journal != nil {
			pool.mu.Lock()
			if err := pool.journal.rotate(pool.local()); err != nil {
				log.DLogger.Warn("Failed to rotate local tx journal", zap.Error(err))
			}
			pool.mu.Unlock()
		}
	}
	journal := g_timer.SetPeriodAndRun(journalHandler, pool.config.Rejournal)
	defer g_timer.StopWork(journal)

	// Keep waiting for and reacting to the various events
	for {
		select {
		// Be unsubscribed due to system stopped
		case <-pool.loopStopCtrl:
			return
		}
	}
}

// TxStatus is the current status of a transaction as seen by the pool.
type TxStatus uint

const (
	TxStatusUnknown TxStatus = iota
	TxStatusQueued
	TxStatusPending
)

// Status returns the status (unknown/pending/queued) of a batch of transactions
// identified by their hashes.
func (pool *TxPool) Status(hashes []common.Hash) []TxStatus {
	pool.mu.RLock()
	defer pool.mu.RUnlock()

	status := make([]TxStatus, len(hashes))
	for i, hash := range hashes {
		if tx := pool.all.Get(hash); tx != nil {
			from, _ := tx.Sender(pool.signer) // already validated
			if pool.pending[from] != nil && pool.pending[from].txs.items[tx.Nonce()] != nil {
				status[i] = TxStatusPending
			} else {
				status[i] = TxStatusQueued
			}
		}
	}
	return status
}

// Get returns a transaction if it is contained in the pool
// and nil otherwise.
func (pool *TxPool) Get(hash common.Hash) model.AbstractTransaction {
	return pool.all.Get(hash)
}

// journalTx adds the specified transaction to the local disk journal if it is
// deemed to have been sent from a local account.
func (pool *TxPool) journalTx(from common.Address, tx model.AbstractTransaction) {
	// Only journal if it's enabled and the transaction is local
	if pool.journal == nil || !pool.locals.contains(from) {
		return
	}
	if err := pool.journal.insert(tx); err != nil {
		log.DLogger.Warn("Failed to journal local transaction", zap.Error(err))
	}
}

// push a transaction into future queue, return true if success , and false if failure.
func (pool *TxPool) enqueueTx(hash common.Hash, tx model.AbstractTransaction) (bool, error) {
	// check if queue have transactions from this address. if not ,make a new txlist for this address.
	from, _ := tx.Sender(pool.signer)
	if pool.queue[from] == nil {
		pool.queue[from] = newTxList(false)
	}

	// if queue have txlist from this address, just insert this transaction .
	inserted, old := pool.queue[from].Add(tx, pool.config.FeeBump)

	// insertion fails, means the replace transaction fee does not exceed the FeeBump
	if !inserted {
		return false, errors.New("new fee is too low to replace old one")
	}
	// insertion success and replace an old transaction , and old transaction  should be removed
	if old != nil {
		pool.all.Remove(old.CalTxId())
		pool.feeList.Removed()
	}
	// insertion success and add the new transaction  to all and feelist.
	if pool.all.Get(hash) == nil {
		pool.all.Add(tx)
		pool.feeList.Put(tx)
	}

	return old != nil, nil
}

// validateTx checks whether a transaction is valid according to the consensus
// rules and adheres to some heuristic limits of the local node (price and size).
func (pool *TxPool) validateTx(tx model.AbstractTransaction, local bool) error {

	// Heuristic limit, reject transactions over 32KB to prevent DOS attacks
	/*	if err := middleware.ValidTxSize(tx); err != nil {
		return g_error.ErrTxOverSize
	}*/
	// Transactions can't be negative. This may never happen using RLP decoded
	// transactions but may occur if you create a transaction using the RPC.
	if tx.Amount().Sign() < 0 {
		return errors.New("tx value can not be negtive")
	}
	// Make sure the transaction is signed properly
	from, err := tx.Sender(pool.signer)
	if err != nil {
		log.DLogger.Error("txPool validateTx the err is:", zap.Error(err))
		return errors.New("invalid sender")
	}
	// Drop non-local transactions under our own minimal accepted gas price
	local = local || pool.locals.contains(from) // account may be local even if the transaction arrived from the network
	//log.DLogger.Info("[validateTx] the local is:", "local", local)
	//log.DLogger.Info("[validateTx] the pool.config.MinFee is: ", "mineFee", pool.config.MinFee)
	//log.DLogger.Info("[validateTx] the tx.fee is: ", "txFee", tx.Fee())

	gas, err := model.IntrinsicGas(tx.ExtraData(), tx.GetType() == common.AddressTypeContractCreate, true)
	if err != nil {
		return err
	}

	if gas > tx.GetGasLimit() {
		return fmt.Errorf("gas limit is to low, need:%v got:%v", gas, tx.GetGasLimit())
	}

	// Ensure the transaction adheres to nonce ordering
	curNonce, err := pool.currentState.GetNonce(from)
	//log.DLogger.Info("the curNonce is:", "curNonce", curNonce)
	//log.DLogger.Info("the tx nonce is:", "txNonce", tx.Nonce())
	//log.DLogger.Info("the pool.currentState.GetNonce result", zap.Error(err))

	if err != nil {
		//log.DLogger.Error("the pool.currentState.GetNonce result", zap.Error(err))
		return err
	}

	if curNonce > tx.Nonce() {
		//log.DLogger.Error("the curNonce is:", "curNonce", curNonce)
		//log.DLogger.Error("the tx nonce is:", "txNonce", tx.Nonce())
		return errors.New("tx nonce is invalid")
	}
	// Transactor should have enough funds to cover the costs
	// cost == V + GP * GL
	curBalance, err := pool.currentState.GetBalance(from)
	//fmt.Println("=======currentbalance======", curBalance, "tx cost", tx.Cost())
	if err != nil || curBalance.Cmp(tx.Cost()) < 0 {
		return errors.New(fmt.Sprintf("tx exceed balance limit, from:%v, cur balance:%v, cost:%v, err:%v", from.Hex(), curBalance.String(), tx.Cost().String(), err))
	}
	//TODO Add economy validator
	return nil
}

// add a single transaction to pool, local mark if this transaction is local or not, and return true if add success and false with err if failure.
func (pool *TxPool) add(tx model.AbstractTransaction, local bool) (bool, error) {
	hash := tx.CalTxId()
	// vaildate the transaction before add it to pool
	if pool.all.Get(hash) != nil {
		//from, _ := tx.Sender(pool.signer)
		//if pool.pending[from] == nil {
		//	log.DLogger.Info("TxPool  add "+from.Hex()+"to pending ", zap.Any("txId", tx.CalTxId()))
		//	pool.promoteExecutables([]common.Address{from})
		//	log.DLogger.Info("TxPool  add  after promoteExecutables", zap.Any("txId", tx.CalTxId()), zap.Bool("is in pending?", pool.pending[from] == nil))
		//} else {
			log.DLogger.Info("the the transaction is already in tx pool", zap.String("txId", hash.Hex()))
			return false, fmt.Errorf("this transaction already in tx pool")
		//}
	}

	if err := pool.validateTx(tx, local); err != nil {
		log.DLogger.Debug("Discarding invalid transaction", zap.Any("hash", hash), zap.Error(err))
		return false, err
	}

	// check if the pool is full or not.
	if uint64(pool.all.Count()) >= pool.config.GlobalSlots+pool.config.GlobalQueue {
		// If the new transaction is underpriced, don't accept it
		if !local && pool.feeList.UnderPriced(tx, pool.locals) {

			log.DLogger.Debug("Discarding underpriced transaction", zap.Any("hash", hash), zap.Any("gasPrice", tx.GetGasPrice()))

			return false, errors.New("transaction items too much")
		}
		// New transaction is better than worse ones, make room for it
		drop := pool.feeList.Discard(pool.all.Count()-int(pool.config.GlobalSlots+pool.config.GlobalQueue-1), pool.locals)
		for _, tx := range drop {
			log.DLogger.Debug("Discarding freshly underpriced transaction", zap.Any("hash", tx.CalTxId()), zap.Any("gasPrice", tx.GetGasPrice()))
			pool.removeTx(tx.CalTxId(), false)
		}
	}

	// If the transaction is replacing an already pending one, do directly
	from, _ := tx.Sender(pool.signer)
	// the pending list have the same address as this transaction, and the list contain a transaction with the same nonce.
	if list := pool.pending[from]; list != nil && list.Overlaps(tx) {
		// add this transaction into  pending list
		inserted, old := list.Add(tx, pool.config.FeeBump)

		// add failed, which means the fee is too low to replace the old one
		if !inserted {
			return false, errors.New("new fee is too low to replace old one")
		}

		// New transaction is replace an old transaction ,so need remove the old transaction.
		if old != nil {
			pool.all.Remove(old.CalTxId())
			pool.feeList.Removed()
		}
		//add this new transaction into all and feelist
		pool.all.Add(tx)
		pool.feeList.Put(tx)

		// add to journal if transaction is local
		pool.journalTx(from, tx)

		log.DLogger.Debug("Pooled new executable transaction", zap.Any("hash", hash), zap.Any("from", from), zap.Any("to", tx.To()))
		return old != nil, nil
	}

	//  If the transaction is not replacing an already pending one, add it to the future queue
	replace, err := pool.enqueueTx(hash, tx)
	if err != nil {
		return false, err
	}

	// Mark local addresses and journal local transactions
	if local {
		pool.locals.add(from)
	}
	pool.journalTx(from, tx)
	return replace, nil
}

// promoteTx moves a transaction to the pending  list of transactions
// and returns whether it was inserted or an older was better.
// Note, this method assumes the pool lock is held!
func (pool *TxPool) promoteTx(addr common.Address, hash common.Hash, tx model.AbstractTransaction) bool {
	// if the pending list do not have this address, make a new list for it.
	if pool.pending[addr] == nil {
		pool.pending[addr] = newTxList(true)
	}
	// otherwise pick the list of the same address
	list := pool.pending[addr]
	// add this transaction into  pending list
	inserted, old := list.Add(tx, pool.config.FeeBump)

	// add failed, which means the fee is too low to replace the old one, need to discard this transaction
	if !inserted {
		// this transaction fee is too low ,need to remove it from the pool and queue list.
		// here we don't remove it from the future queue ,need to delete it as well when this promote fails.
		pool.all.Remove(hash)
		pool.feeList.Removed()

		return false
	}
	// new transaction is better, and it replaced an old transaction in pending list. discard the old transaction
	if old != nil {
		pool.all.Remove(old.CalTxId())
		pool.feeList.Removed()

	}
	// add this transaction to all and feelist, normally it already in all and feelist. just add it to make it safe.
	// Failsafe to work around direct pending inserts (tests)
	if pool.all.Get(hash) == nil {
		pool.all.Add(tx)
		pool.feeList.Put(tx)
	}

	// reset the heartbeat and add the manage state nonce by 1
	pool.beats[addr] = time.Now()
	pool.pendingState.SetNonce(addr, tx.Nonce()+1)
	return true
}

// local retrieves all currently known local transactions, groupped by origin
// account and sorted by nonce. The returned transaction set is a copy and can be
// freely modified by calling code.
func (pool *TxPool) local() map[common.Address][]model.AbstractTransaction {
	txs := make(map[common.Address][]model.AbstractTransaction)
	for addr := range pool.locals.accounts {
		if pending := pool.pending[addr]; pending != nil {
			txs[addr] = append(txs[addr], pending.Flatten()...)
		}
		if queued := pool.queue[addr]; queued != nil {
			txs[addr] = append(txs[addr], queued.Flatten()...)
		}
	}
	return txs
}

// Pending retrieves all currently processable transactions, groupped by origin
// account and sorted by nonce. The returned transaction set is a copy and can be
// freely modified by calling code.
func (pool *TxPool) Pending() (map[common.Address][]model.AbstractTransaction, error) {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	pending := make(map[common.Address][]model.AbstractTransaction)
	for addr, list := range pool.pending {
		pending[addr] = list.Flatten()
	}
	return pending, nil
}

//get all transactions in future queue
func (pool *TxPool) Queueing() (map[common.Address][]model.AbstractTransaction, error) {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	queueing := make(map[common.Address][]model.AbstractTransaction)
	for addr, list := range pool.queue {
		queueing[addr] = list.Flatten()
	}
	return queueing, nil
}

// addressByHeartbeat is an account address tagged with its last activity timestamp.
type addressByHeartbeat struct {
	address   common.Address
	heartbeat time.Time
}

type addressesByHeartbeat []addressByHeartbeat

func (a addressesByHeartbeat) Len() int           { return len(a) }
func (a addressesByHeartbeat) Less(i, j int) bool { return a[i].heartbeat.Before(a[j].heartbeat) }
func (a addressesByHeartbeat) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func (pool *TxPool) promoteExecutables(accounts []common.Address) {
	// Track the promoted transactions to broadcast them at once
	//var promoted []model.AbstractTransaction

	// if accounts is nil then go through the whole pool and generate the accounts set.
	if accounts == nil {
		accounts = make([]common.Address, 0, len(pool.queue))
		for addr := range pool.queue {
			accounts = append(accounts, addr)
		}
	}
	// Iterate over all accounts and promote any executable transactions
	for _, addr := range accounts {
		list := pool.queue[addr]
		if list == nil {
			continue // Just in case someone calls with a non existing account
		}

		// Drop all transactions that are deemed too old (low nonce)
		oldNonce, _ := pool.currentState.GetNonce(addr)
		for _, tx := range list.FilterNonce(oldNonce) {
			hash := tx.CalTxId()
			log.DLogger.Debug("Removed old queued transaction", zap.Any("hash", hash))
			pool.all.Remove(hash)
			pool.feeList.Removed()
		}

		// Drop all transactions that are too costly (balance can not cover the cost)
		oldBalance, _ := pool.currentState.GetBalance(addr)
		drops, _ := list.Filter(oldBalance)
		for _, tx := range drops {
			hash := tx.CalTxId()
			log.DLogger.Debug("Removed unpayable queued transaction", zap.Any("hash", hash))
			pool.all.Remove(hash)
			pool.feeList.Removed()
		}

		// Gather all executable transactions with consecutive nonce and promote them to pending list
		newNonce := pool.pendingState.GetNonce(addr)
		pendingLen, _ := pool.stats()
		if uint64(pendingLen) < pool.config.GlobalSlots {
			for _, tx := range list.Ready(newNonce) {
				hash := tx.CalTxId()
				if pool.promoteTx(addr, hash, tx) {
					//log.DLogger.Debug("Promoting queued transaction", zap.Any("hash", hash))
					//promoted = append(promoted, tx)
				}
			}
		}
		// if the pending list for this address is oversize and this address is not in local whitelist, remove all the tail transactions.
		if !pool.locals.contains(addr) {
			for _, tx := range list.Cap(int(pool.config.AccountQueue)) {
				hash := tx.CalTxId()
				pool.all.Remove(hash)
				pool.feeList.Removed()
				log.DLogger.Debug("Removed cap-exceeding queued transaction", zap.Any("hash", hash))
			}
		}

		// Delete the entire queue entry if it became empty.
		if list.Empty() {
			delete(pool.queue, addr)
		}
	}

	// If the pending limit is overflown, start equalizing allowances
	pending := uint64(0)
	for _, list := range pool.pending {
		pending += uint64(list.Len())
	}

	// If the pending list exceed the max size, remove some transactions.
	// The rule for removing is first, sort all the nonlocal address with the length of the txlist,
	// add two address txlist with largest length and use the shorter one's length as threshold to cut the longer one.
	// and the add the third longest txlist and use its length as the new threshold and cut the first two txlist.
	// keep doing this until all address txlist have the same length.
	// if it still exceed, sub all txlist by one and keep doing this, stop if the total size is smaller then the max size or only local transactions left.
	if pending > pool.config.GlobalSlots {
		//pendingBeforeCap := pending
		// Assemble a spam order to penalize large transactors first
		spammers := prque.New()
		for addr, list := range pool.pending {
			// Only evict transactions from high rollers
			if !pool.locals.contains(addr) && uint64(list.Len()) > pool.config.AccountSlots {
				spammers.Push(addr, float32(list.Len()))
			}
		}
		// Gradually drop transactions from offenders
		offenders := []common.Address{}
		for pending > pool.config.GlobalSlots && !spammers.Empty() {
			// Retrieve the next offender if not local address
			offender, _ := spammers.Pop()
			offenders = append(offenders, offender.(common.Address))

			// Equalize balances until all the same or below threshold
			if len(offenders) > 1 {
				// Calculate the equalization threshold for all current offenders
				threshold := pool.pending[offender.(common.Address)].Len()

				// Iteratively reduce all offenders until below limit or threshold reached
				for pending > pool.config.GlobalSlots && pool.pending[offenders[len(offenders)-2]].Len() > threshold {
					for i := 0; i < len(offenders)-1; i++ {
						list := pool.pending[offenders[i]]
						for _, tx := range list.Cap(list.Len() - 1) {
							// Drop the transaction from the global pools too
							hash := tx.CalTxId()
							pool.all.Remove(hash)
							pool.feeList.Removed()

							//// Update the account nonce to the dropped transaction
							newNonce := pool.pendingState.GetNonce(offenders[i])
							if nonce := tx.Nonce(); newNonce > nonce {
								pool.pendingState.SetNonce(offenders[i], nonce)
							}
							log.DLogger.Debug("Removed fairness-exceeding pending transaction", zap.Any("hash", hash))
						}
						pending--
					}
				}
			}
		}
		// If still above threshold, reduce to limit or min allowance
		if pending > pool.config.GlobalSlots && len(offenders) > 0 {
			for pending > pool.config.GlobalSlots && uint64(pool.pending[offenders[len(offenders)-1]].Len()) > pool.config.AccountSlots {
				for _, addr := range offenders {
					list := pool.pending[addr]
					for _, tx := range list.Cap(list.Len() - 1) {
						// Drop the transaction from the global pools too
						hash := tx.CalTxId()
						pool.all.Remove(hash)
						pool.feeList.Removed()
						// Update the account nonce to the dropped transaction
						newNonce := pool.pendingState.GetNonce(addr)
						if nonce := tx.Nonce(); newNonce > nonce {
							pool.pendingState.SetNonce(addr, nonce)
						}
						log.DLogger.Debug("Removed fairness-exceeding pending transaction", zap.Any("hash", hash))
					}
					pending--
				}
			}
		}
	}
	// If we've queued more transactions than the hard limit, drop oldest ones
	queued := uint64(0)
	for _, list := range pool.queue {
		queued += uint64(list.Len())
	}
	if queued > pool.config.GlobalQueue {
		// Sort all accounts with queued transactions by heartbeat
		addresses := make(addressesByHeartbeat, 0, len(pool.queue))
		for addr := range pool.queue {
			if !pool.locals.contains(addr) { // don't drop locals
				addresses = append(addresses, addressByHeartbeat{addr, pool.beats[addr]})
			}
		}
		sort.Sort(addresses)

		// Drop transactions until the total is below the limit or only locals remain
		for drop := queued - pool.config.GlobalQueue; drop > 0 && len(addresses) > 0; {
			addr := addresses[len(addresses)-1]
			list := pool.queue[addr.address]

			addresses = addresses[:len(addresses)-1]

			// Drop all transactions if they are less than the overflow
			if size := uint64(list.Len()); size <= drop {
				for _, tx := range list.Flatten() {
					pool.removeTx(tx.CalTxId(), true)
				}
				drop -= size
				continue
			}
			// Otherwise drop only last few transactions
			txs := list.Flatten()
			for i := len(txs) - 1; i >= 0 && drop > 0; i-- {
				pool.removeTx(txs[i].CalTxId(), true)
				drop--
			}
		}
	}
}

// demoteUnexecutables removes invalid and processed transactions from the pools
// executable/pending queue and any subsequent transactions that become unexecutable
// are moved back into the future queue.
func (pool *TxPool) demoteUnexecutables() {
	// Iterate over all accounts and demote any non-executable transactions
	for addr, list := range pool.pending {
		nonce, _ := pool.currentState.GetNonce(addr)

		// Drop all transbuild dipperin-core failedactions that are deemed too old (low nonce)
		for _, tx := range list.FilterNonce(nonce) {
			hash := tx.CalTxId()
			log.DLogger.Debug("Removed old pending transaction", zap.Any("hash", hash))
			pool.all.Remove(hash)
			pool.feeList.Removed()
		}
		// Drop all transactions that are too costly (low balance or out of gas), and queue any invalids back for later

		balance, _ := pool.currentState.GetBalance(addr)
		drops, invalids := list.Filter(balance)
		for _, tx := range drops {
			hash := tx.CalTxId()
			log.DLogger.Debug("Removed unpayable pending transaction", zap.Any("hash", hash))
			pool.all.Remove(hash)
			pool.feeList.Removed()

		}
		for _, tx := range invalids {
			hash := tx.CalTxId()
			log.DLogger.Debug("Demoting pending transaction", zap.Any("hash", hash))
			pool.enqueueTx(hash, tx)
		}

		// If there's a gap in front, alert (should never happen) and postpone all transactions
		if list.Len() > 0 && list.txs.Get(nonce) == nil {
			for _, tx := range list.Cap(0) {
				hash := tx.CalTxId()
				log.DLogger.Error("Demoting invalidated transaction", zap.Any("hash", hash))
				pool.enqueueTx(hash, tx)
			}
		}

		// Delete the entire queue entry if it became empty.
		if list.Empty() {
			delete(pool.pending, addr)
			delete(pool.beats, addr)
		}
	}
}

func (pool *TxPool) RemoveTxs(newBlock model.AbstractBlock) {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	newBlock.TxIterator(func(i int, transaction model.AbstractTransaction) error {
		pool.removeTx(transaction.CalTxId(), true)
		return nil
	})
}

func (pool *TxPool) RemoveTxsBatch(txIds []common.Hash) {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	for _, txId := range txIds {
		pool.removeTx(txId, true)
	}
}

// removeTx removes a single transaction from the queue, moving all subsequent
// transactions back to the future queue.
// will clear the feelist index if outofbound is ture, which means totally removed
// if outofbound is false, the feelist wont be removed.
func (pool *TxPool) removeTx(hash common.Hash, outofbound bool) {
	// Fetch the transaction we wish to delete
	tx := pool.all.Get(hash)
	//if this transaction is not exist in pool directly return
	if tx == nil {
		return
	}
	//otherwise get the txlist for the same address as this transaction.
	addr, _ := tx.Sender(pool.signer)

	// Remove it from the pool
	pool.all.Remove(hash)
	if outofbound {
		pool.feeList.Removed()
	}

	// Remove the transaction from the pending lists and reset the account nonce
	if pending := pool.pending[addr]; pending != nil {
		// invalids are the transactions following by this transaction, should be remove from the pending list to future queue as well
		if removed, invalids := pending.Remove(tx); removed {
			// If no more pending transactions are left, remove the list
			if pending.Empty() {
				delete(pool.pending, addr)
				delete(pool.beats, addr)
			}

			// move the invalid transactions to future queue
			for _, tx := range invalids {
				pool.enqueueTx(tx.CalTxId(), tx)
			}

			// Update the account nonce if needed
			nonceGet := pool.pendingState.GetNonce(addr)
			if nonce := tx.Nonce(); nonceGet > nonce {
				pool.pendingState.SetNonce(addr, nonce)
			}
			return
		}
	}

	// Transaction is in the future queue
	if future := pool.queue[addr]; future != nil {
		future.Remove(tx)
		if future.Empty() {
			delete(pool.queue, addr)
		}
	}
}

// addTx enqueues a single transaction into the pool if it is valid.
func (pool *TxPool) addTx(tx model.AbstractTransaction, local bool) error {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	// Try to inject the transaction and update any state
	replace, err := pool.add(tx, local)
	if err != nil {
		//log.DLogger.Info("pool addTx err","err",err)
		return err
	}

	// If we added a new transaction, run promotion checks and return
	if !replace {
		from, _ := tx.Sender(pool.signer) // already validated
		pool.promoteExecutables([]common.Address{from})
	}
	log.DLogger.Debug("Add tx success", zap.String("txid", tx.CalTxId().Hex()))
	return nil
}

// addTxs attempts to queue a batch of transactions if they are valid.
func (pool *TxPool) addTxs(txs []model.AbstractTransaction, local bool) []error {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	return pool.addTxsLocked(txs, local)
}

// addTxsLocked attempts to queue a batch of transactions if they are valid,
// whilst assuming the transaction pool lock is already held.
func (pool *TxPool) addTxsLocked(txs []model.AbstractTransaction, local bool) []error {
	// Add the batch of transaction, tracking the accepted ones
	dirty := make(map[common.Address]struct{})
	errs := make([]error, len(txs))

	for i, tx := range txs {
		var replace bool
		if replace, errs[i] = pool.add(tx, local); errs[i] == nil && !replace {
			from, _ := tx.Sender(pool.signer) // already validated
			dirty[from] = struct{}{}
		} else {
			// only for debug
			//txSender, _ := tx.Sender(nil)
			//log.DLogger.Warn("add tx to pool failed", zap.Error(err)s[i], "tx sender", txSender.Hex())
		}
	}
	// Only reprocess the internal state if something was actually added
	if len(dirty) > 0 {
		addrs := make([]common.Address, 0, len(dirty))
		for addr := range dirty {
			addrs = append(addrs, addr)
		}
		pool.promoteExecutables(addrs)
	}
	return errs
}

// AddLocal enqueues a single transaction into the pool if it is valid, marking
// the sender as a local one in the mean time, ensuring it goes around the local
// pricing constraints.
func (pool *TxPool) AddLocal(tx model.AbstractTransaction) error {
	//log.DLogger.Info("add local tx is:","txId",tx.CalTxId().Hex())

	return pool.addTx(tx, !pool.config.NoLocals)
}

// AddRemote enqueues a single transaction into the pool if it is valid. If the
// sender is not among the locally tracked ones, full pricing constraints will
// apply.
func (pool *TxPool) AddRemote(tx model.AbstractTransaction) error {

	//log.DLogger.Info("add remote tx is:","txId",tx.CalTxId().Hex())

	return pool.addTx(tx, false)
}

// AddLocals enqueues a batch of transactions into the pool if they are valid,
// marking the senders as a local ones in the mean time, ensuring they go around
// the local pricing constraints.
func (pool *TxPool) AddLocals(txs []model.AbstractTransaction) []error {
	return pool.addTxs(txs, !pool.config.NoLocals)
}

// AddRemotes enqueues a batch of transactions into the pool if they are valid.
// If the senders are not among the locally tracked ones, full pricing constraints
// will apply.
func (pool *TxPool) AddRemotes(txs []model.AbstractTransaction) []error {
	return pool.addTxs(txs, false)
}

//// pick a packaged transaction
//func (pool *TxPool) RandTxsForPack() (result []model.AbstractTransaction) {
//	txsLen := len(pool.waitPackTxs)
//	if txsLen <= chain_config.MaxTxInBlock {
//		return pool.waitPackTxs
//	}
//	rand.Seed(time.Now().UnixNano())
//	x := rand.Intn(len(pool.waitPackTxs))
//	for i := 0; i < chain_config.MaxTxInBlock; i++ {
//		if x > txsLen-1 {
//			x = 0
//		}
//		result = append(result, pool.waitPackTxs[x])
//		x++
//	}
//	return
//}

func (pool *TxPool) ConvertPoolToMap() map[common.Hash]model.AbstractTransaction {
	m := make(map[common.Hash]model.AbstractTransaction)
	txs := pool.all.Flatten()
	for _, tx := range txs {
		m[tx.CalTxId()] = tx
	}
	return m
}

//add tx in TxPool
//func (pool *TxPool) NewTransaction(tx model.AbstractTransaction) error {
//	//todo validate transaction
//	log.DLogger.Debug("~~~~~~~~~~~~~~~~the txPool is: ", "txPool", *pool)
//	log.DLogger.Debug("~~~~~~~~~~~~~~~~the tx is: ", "tx", tx)
//
//	//add to transaction pool
//	for _,tmpTx := range pool.waitPackTxs{
//		if tmpTx.CalTxId() == tx.CalTxId() {
//			return nil
//		}
//	}
//
//	pool.waitPackTxs = append(pool.waitPackTxs, tx)
//	return nil
//}

//Remove the transactions that have been packaged and return the removed transactions
//func (pool *TxPool) RemoveTxs(txs []model.AbstractTransaction) []model.AbstractTransaction {
//	for _, tx := range txs {
//		for index, tmpTx := range pool.waitPackTxs {
//			if tmpTx.CalTxId() == tx.CalTxId() {
//				pool.waitPackTxs = append(pool.waitPackTxs[:index], pool.waitPackTxs[index+1:]...)
//				break
//			}
//		}
//	}
//	return pool.waitPackTxs
//}

// accountSet is simply a set of addresses to check for existence, and a signer
// capable of deriving addresses from transactions.
type accountSet struct {
	accounts map[common.Address]struct{}
	signer   model.Signer
}

// newAccountSet creates a new address set with an associated signer for sender
// derivations.
func newAccountSet(signer model.Signer) *accountSet {
	return &accountSet{
		accounts: make(map[common.Address]struct{}),
		signer:   signer,
	}
}

// contains checks if a given address is contained within the set.
func (as *accountSet) contains(addr common.Address) bool {
	_, exist := as.accounts[addr]
	return exist
}

// containsTx checks if the sender of a given tx is within the set. If the sender
// cannot be derived, this method returns false.
func (as *accountSet) containsTx(tx model.AbstractTransaction) bool {
	if addr, err := tx.Sender(as.signer); err == nil {
		return as.contains(addr)
	}
	return false
}

// add inserts a new address into the set to track.
func (as *accountSet) add(addr common.Address) {
	as.accounts[addr] = struct{}{}
}

// txLookup is used internally by TxPool to track transactions while allowing lookup without
// mutex contention.
//
// Note, although this type is properly protected against concurrent access, it
// is **not** a type that should ever be mutated or even exposed outside of the
// transaction pool, since its internal state is tightly coupled with the pools
// internal mechanisms. The sole purpose of the type is to permit out-of-bound
// peeking into the pool in TxPool.Get without having to acquire the widely scoped
// TxPool.mu mutex.
type txLookup struct {
	all  map[common.Hash]model.AbstractTransaction
	lock sync.RWMutex
}

// newTxLookup returns a new txLookup structure.
func newTxLookup() *txLookup {
	return &txLookup{
		all: make(map[common.Hash]model.AbstractTransaction),
	}
}

func (t *txLookup) Flatten() []model.AbstractTransaction {
	t.lock.RLock()
	defer t.lock.RUnlock()

	var txs []model.AbstractTransaction
	for _, v := range t.all {
		txs = append(txs, v)
	}
	return txs
}

// Range calls f on each key and value present in the map.
func (t *txLookup) Range(f func(hash common.Hash, tx model.AbstractTransaction) bool) {
	t.lock.RLock()
	defer t.lock.RUnlock()

	for key, value := range t.all {
		// stops if f returns false
		if !f(key, value) {
			break
		}
	}
}

// Get returns a transaction if it exists in the lookup, or nil if not found.
func (t *txLookup) Get(hash common.Hash) model.AbstractTransaction {
	t.lock.RLock()
	defer t.lock.RUnlock()

	return t.all[hash]
}

// Count returns the current number of items in the lookup.
func (t *txLookup) Count() int {
	t.lock.RLock()
	defer t.lock.RUnlock()

	return len(t.all)
}

// Add adds a transaction to the lookup.
func (t *txLookup) Add(tx model.AbstractTransaction) {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.all[tx.CalTxId()] = tx
}

// Remove removes a transaction from the lookup.
func (t *txLookup) Remove(hash common.Hash) {
	t.lock.Lock()
	defer t.lock.Unlock()

	delete(t.all, hash)
}

// first: pending, second: queued
func (pool *TxPool) Stats() (int, int) {
	pool.mu.RLock()
	defer pool.mu.RUnlock()

	return pool.stats()
}

// stats retrieves the current pool stats, namely the number of pending and the
// number of queued (non-executable) transactions.
func (pool *TxPool) stats() (int, int) {
	pending := 0
	for _, list := range pool.pending {
		pending += list.Len()
	}
	queued := 0
	for _, list := range pool.queue {
		queued += list.Len()
	}
	return pending, queued
}

// Create local tx pool hybrid estimator
func (pool *TxPool) GetTxsEstimator(broadcastBloom *iblt.Bloom) *iblt.HybridEstimator {
	// hack
	c := iblt.NewHybridEstimatorConfig()
	estimator := iblt.NewHybridEstimator(c)

	//startAt := time.Now()
	// get peer local tx pool all txs
	txs := pool.all.Flatten()
	//log.DLogger.Info("tx pool flatten", "use", time.Now().Sub(startAt))

	//startAt1 := time.Now()
	// get peer estimator
	for _, tx := range txs {
		//startAt2 := time.Now()
		b := tx.CalTxId().Bytes()
		//bloom_log.DLogger.Info("cal tx id", "use", time.Now().Sub(startAt2))
		if broadcastBloom.LookUp(b) {
			estimator.EncodeByte(b)
		}
	}
	//dis := time.Now().Sub(startAt1)
	//bloom_log.DLogger.Info("broadcastBloom.LookUp", "use", dis, "txs len", len(txs))
	//log.DLogger.Info("broadcastBloom.LookUp", "use", dis, "txs len", len(txs))

	return estimator
}
