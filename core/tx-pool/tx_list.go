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
	"container/heap"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third-party/log"
	"math/big"
	"sort"
)

// nonceHeap is a heap.Interface implementation over 64bit unsigned integers for
// retrieving sorted transactions from the possibly gapped future queue.
type nonceHeap []uint64

func (h nonceHeap) Len() int           { return len(h) }
func (h nonceHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h nonceHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *nonceHeap) Push(x interface{}) {
	*h = append(*h, x.(uint64))
}

func (h *nonceHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// txSortedMap is a nonce->transaction hash map with a heap based index to allow
// iterating over the contents in a nonce-incrementing way.
type txSortedMap struct {
	items map[uint64]model.AbstractTransaction // Hash map storing the transaction data
	index *nonceHeap                           // Heap of nonces of all the stored transactions (non-strict mode)
	cache []model.AbstractTransaction
}

// newTxSortedMap creates a new nonce-sorted transaction map.
func newTxSortedMap() *txSortedMap {
	return &txSortedMap{
		items: make(map[uint64]model.AbstractTransaction),
		index: new(nonceHeap),
	}
}

// Get retrieves the current transactions associated with the given nonce.
func (m *txSortedMap) Get(nonce uint64) model.AbstractTransaction {
	return m.items[nonce]
}

// Put inserts a new transaction into the map, also updating the map's nonce
// index. If a transaction already exists with the same nonce, it's overwritten.
func (m *txSortedMap) Put(tx model.AbstractTransaction) {
	nonce := tx.Nonce()
	if m.items[nonce] == nil {
		heap.Push(m.index, nonce)
	}
	m.items[nonce], m.cache = tx, nil
}

// FilterNonce removes all transactions from the map with a nonce lower than the
// provided threshold. Every removed transaction is returned for any post-removal
// maintenance.
func (m *txSortedMap) FilterNonce(threshold uint64) []model.AbstractTransaction {
	var removed []model.AbstractTransaction

	// Pop off heap items until the threshold is reached
	for m.index.Len() > 0 && (*m.index)[0] < threshold {
		nonce := heap.Pop(m.index).(uint64)
		removed = append(removed, m.items[nonce])
		delete(m.items, nonce)
	}
	// If we had a cached order, shift the front
	if m.cache != nil {
		m.cache = m.cache[len(removed):]
	}
	return removed
}

// Sort tries to Sort all transactions from the map and writes them to cache
func (m *txSortedMap) Sort() {
	// If the sorting was not cached yet, create and cache it
	if m.cache == nil {
		m.cache = make([]model.AbstractTransaction, 0, len(m.items))
		for _, tx := range m.items {
			m.cache = append(m.cache, tx)
		}

		sort.Sort(TxByNonce(m.cache))
	}
}

// Filter iterates over the list of transactions and removes all of them for which
// the specified function evaluates to true.
func (m *txSortedMap) Filter(filter func(model.AbstractTransaction) bool) []model.AbstractTransaction {
	var removed []model.AbstractTransaction

	m.Sort()

	// slices that share the same memory backing
	*m.index = (*m.index)[:0]
	newCache := m.cache[:0]

	// Collect all the transactions to filter out
	for _, tx := range m.cache {
		if filter(tx) {
			removed = append(removed, tx)
			delete(m.items, tx.Nonce())
		} else {
			*m.index = append(*m.index, tx.Nonce())
			newCache = append(newCache, tx)
		}
	}

	heap.Init(m.index)
	m.cache = newCache

	return removed
}

// Cap places a hard limit on the number of items, returning all transactions
// exceeding that limit.
func (m *txSortedMap) Cap(threshold int) []model.AbstractTransaction {
	// Short circuit if the number of items is under the limit
	if len(m.items) <= threshold {
		return nil
	}
	// Otherwise gather and drop the highest nonce'd transactions
	var drops []model.AbstractTransaction

	sort.Sort(*m.index)
	// removing transactions with larger nonce
	for size := len(m.items); size > threshold; size-- {
		drops = append(drops, m.items[(*m.index)[size-1]])
		delete(m.items, (*m.index)[size-1])
	}
	*m.index = (*m.index)[:threshold]
	heap.Init(m.index)

	// If we had a cache, shift the back
	if m.cache != nil {
		m.cache = m.cache[:len(m.cache)-len(drops)]
	}
	return drops
}

// Remove deletes a transaction from the maintained map, returning whether the
// transaction was found.
func (m *txSortedMap) Remove(nonce uint64) bool {
	// Short circuit if no transaction is present
	_, ok := m.items[nonce]
	if !ok {
		return false
	}
	// Otherwise delete the transaction and fix the heap index
	for i := 0; i < m.index.Len(); i++ {
		if (*m.index)[i] == nonce {
			heap.Remove(m.index, i)
			break
		}
	}
	delete(m.items, nonce)
	m.cache = nil

	return true
}

// Ready retrieves a sequentially increasing list of transactions starting at the
// provided nonce that is ready for processing. The returned transactions will be
// removed from the list.
//
// Note, all transactions with nonces lower than start will also be returned to
// prevent getting into and invalid state. This is not something that should ever
// happen but better to be self correcting than failing!
func (m *txSortedMap) Ready(start uint64) []model.AbstractTransaction {
	// Short circuit if no transactions are available
	if m.index.Len() == 0 || (*m.index)[0] > start {
		return nil
	}
	// Otherwise start accumulating incremental transactions
	var ready []model.AbstractTransaction
	for next := (*m.index)[0]; m.index.Len() > 0 && (*m.index)[0] == next; next++ {
		ready = append(ready, m.items[next])
		delete(m.items, next)
		heap.Pop(m.index)
	}
	m.cache = nil

	return ready
}

// Len returns the length of the transaction map.
func (m *txSortedMap) Len() int {
	return len(m.items)
}

// Flatten creates a nonce-sorted slice of transactions based on the loosely
// sorted internal representation. The result of the sorting is cached in case
// it's requested again before any modifications are made to the contents.
func (m *txSortedMap) Flatten() []model.AbstractTransaction {
	m.Sort()

	// Copy the cache to prevent accidental modifications
	txs := make([]model.AbstractTransaction, len(m.cache))
	copy(txs, m.cache)
	return txs
}

type TxByNonce []model.AbstractTransaction

func (t TxByNonce) Len() int           { return len(t) }
func (t TxByNonce) Less(i, j int) bool { return t[i].Nonce() < t[j].Nonce() }
func (t TxByNonce) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }

// txList is a "list" of transactions belonging to an account, sorted by account
// nonce. The same type can be used both for storing contiguous transactions for
// the executable/pending queue; and for storing gapped transactions for the non-
// executable/future queue, with minor behavioral changes.
type txList struct {
	strict bool         // Whether nonces are strictly continuous or not
	txs    *txSortedMap // Heap indexed sorted hash map of the transactions

	costCap *big.Int // Amount of the highest costing transaction (reset only if exceeds balance)
}

// newTxList create a new transaction list for maintaining nonce-indexable fast,
// gapped, sortable transaction lists.
func newTxList(strict bool) *txList {
	return &txList{
		strict:  strict,
		txs:     newTxSortedMap(),
		costCap: new(big.Int),
	}
}

// Overlaps returns whether the transaction specified has the same nonce as one
// already contained within the list.
func (l *txList) Overlaps(tx model.AbstractTransaction) bool {
	return l.txs.Get(tx.Nonce()) != nil
}

// Add tries to insert a new transaction into the list, returning whether the
// transaction was accepted, and if yes, any previous transaction it replaced.
//
// If the new transaction is accepted into the list, the lists' cost thresholds
// is also potentially updated.
func (l *txList) Add(tx model.AbstractTransaction, feeBump uint64) (bool, model.AbstractTransaction) {
	// If there's an older better transaction, abort
	old := l.txs.Get(tx.Nonce())
	if old != nil {
		threshold := new(big.Int).Div(new(big.Int).Mul(old.Fee(), big.NewInt(100+int64(feeBump))), big.NewInt(100))
		// threshold = old.Fee() * (1 + feeBump/100)
		// Have to ensure that the new fee is higher than the fee as well as
		// checking the percentage threshold to ensure that this is accurate
		// for low fee replacements
		if old.Fee().Cmp(tx.Fee()) >= 0 || threshold.Cmp(tx.Fee()) > 0 {
			// old transaction fee is higher than current one
			// OR
			// the input price bump is less than price bump threshold
			return false, nil
		}
	}

	// Otherwise overwrite the old transaction with the current one
	l.txs.Put(tx)

	// if the fee is higher than the current costCap, update costCap
	if cost := tx.Cost(); l.costCap.Cmp(cost) < 0 {
		l.costCap = cost
	}

	return true, old
}

// FilterNonce removes all transactions from the list with a nonce lower than the
// provided threshold. Every removed transaction is returned for any post-removal
// maintenance.
func (l *txList) FilterNonce(threshold uint64) []model.AbstractTransaction {
	return l.txs.FilterNonce(threshold)
}

// Filter removes all transactions from the list with a cost higher
// than the provided thresholds. Every removed transaction is returned for any
// post-removal maintenance. Strict-mode invalidated transactions are also
// returned.
//
// This method uses the cached costCap to quickly decide if there's even
// a point in calculating all the fees or if the balance covers all. If the threshold
// is lower than the cap, the cap will be reset to a new high after removing
// the newly invalidated transactions.
func (l *txList) Filter(threshold *big.Int) ([]model.AbstractTransaction, []model.AbstractTransaction) {
	// If all transactions are below the threshold, short circuit
	if l.costCap.Cmp(threshold) <= 0 {
		return nil, nil
	}
	l.costCap = new(big.Int).Set(threshold) // Lower the caps to the thresholds

	// Filter out all the transactions above the account's funds
	removed := l.txs.Filter(func(tx model.AbstractTransaction) bool { return tx.Cost().Cmp(threshold) > 0 })

	// If the list was strict, filter anything above the lowest nonce
	var invalids []model.AbstractTransaction

	if l.strict && len(removed) > 0 {
		lowest := removed[0].Nonce()

		invalids = l.txs.Filter(func(tx model.AbstractTransaction) bool { return tx.Nonce() > lowest })
	}
	return removed, invalids
}

// Cap places a hard limit on the number of items, returning all transactions
// exceeding that limit.
func (l *txList) Cap(threshold int) []model.AbstractTransaction {
	return l.txs.Cap(threshold)
}

// Remove deletes a transaction from the maintained list, returning whether the
// transaction was found, and also returning any transaction invalidated due to
// the deletion (strict mode only).
func (l *txList) Remove(tx model.AbstractTransaction) (bool, []model.AbstractTransaction) {
	// Remove the transaction from the set
	nonce := tx.Nonce()
	if removed := l.txs.Remove(nonce); !removed {
		return false, nil
	}

	// In strict mode, filter out non-executable transactions
	if l.strict {
		return true, l.txs.Filter(func(tx model.AbstractTransaction) bool { return tx.Nonce() > nonce })
	}
	return true, nil
}

// Ready retrieves a sequentially increasing list of transactions starting at the
// provided nonce that is ready for processing. The returned transactions will be
// removed from the list.
//
// Note, all transactions with nonces lower than start will also be returned to
// prevent getting into and invalid state. This is not something that should ever
// happen but better to be self correcting than failing!
func (l *txList) Ready(start uint64) []model.AbstractTransaction {
	return l.txs.Ready(start)
}

// Len returns the length of the transaction list.
func (l *txList) Len() int {
	return l.txs.Len()
}

// Empty returns whether the list of transactions is empty or not.
func (l *txList) Empty() bool {
	return l.Len() == 0
}

// Flatten creates a nonce-sorted slice of transactions based on the loosely
// sorted internal representation. The result of the sorting is cached in case
// it's requested again before any modifications are made to the contents.
func (l *txList) Flatten() []model.AbstractTransaction {
	return l.txs.Flatten()
}

// priceHeap is a heap.Interface implementation over transactions for retrieving
// price-sorted transactions to discard when the pool fills up.
type priceHeap []model.AbstractTransaction

func (h priceHeap) Len() int      { return len(h) }
func (h priceHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h priceHeap) Less(i, j int) bool {
	// Sort primarily by price, returning the cheaper one
	switch h[i].Fee().Cmp(h[j].Fee()) {
	case -1:
		return true
	case 1:
		return false
	}
	// If the prices match, stabilize via nonces (high nonce is worse)
	return h[i].Nonce() > h[j].Nonce()
}

func (h *priceHeap) Push(x interface{}) {
	*h = append(*h, x.(model.AbstractTransaction))
}

func (h *priceHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// txFeeList is a price-sorted heap to allow operating on transactions pool
// contents in a price-incrementing way.
type txFeeList struct {
	all    *txLookup  // Pointer to the map of all transactions
	items  *priceHeap // Heap of prices of all the stored transactions
	stales int        // Number of stale price points to (re-heap trigger)
}

// newTxFeeList creates a new price-sorted transaction heap.
func newTxFeeList(all *txLookup) *txFeeList {
	return &txFeeList{
		all:   all,
		items: new(priceHeap),
	}
}

// Put inserts a new transaction into the heap.
func (l *txFeeList) Put(tx model.AbstractTransaction) {
	heap.Push(l.items, tx)
}

// Removed notifies the prices transaction list that an old transaction dropped
// from the pool. The list will just keep a counter of stale objects and update
// the heap if a large enough ratio of transactions go stale.
func (l *txFeeList) Removed() {
	// Bump the stale counter, but exit if still too low (< 25%)
	l.stales++
	if l.stales <= len(*l.items)/4 {
		return
	}
	// Seems we've reached a critical number of stale transactions, newHeap
	newHeap := make(priceHeap, 0, l.all.Count())

	l.stales, l.items = 0, &newHeap
	l.all.Range(func(hash common.Hash, tx model.AbstractTransaction) bool {
		*l.items = append(*l.items, tx)
		return true
	})
	heap.Init(l.items)
}

// Cap finds all the transactions below the given fee threshold, drops them
// from the fee list and returns them for further removal from the entire pool.
func (l *txFeeList) Cap(threshold *big.Int, local *accountSet) []model.AbstractTransaction {
	drop := make([]model.AbstractTransaction, 0, 128) // Remote under priced transactions to drop
	save := make([]model.AbstractTransaction, 0, 64)  // Local under priced transactions to keep

	for len(*l.items) > 0 {
		// Discard stale transactions if found during cleanup
		tx := heap.Pop(l.items).(model.AbstractTransaction)
		if l.all.Get(tx.CalTxId()) == nil {
			l.stales--
			continue
		}
		// Stop the discards if we've reached the threshold
		if tx.Fee().Cmp(threshold) >= 0 {
			save = append(save, tx)
			break
		}
		// Non stale transaction found, discard unless local
		if local.containsTx(tx) {
			save = append(save, tx)
		} else {
			drop = append(drop, tx)
		}
	}
	for _, tx := range save {
		heap.Push(l.items, tx)
	}
	return drop
}

// UnderPriced checks whether a transaction is cheaper than (or as cheap as) the
// lowest feeList transaction currently being tracked.
func (l *txFeeList) UnderPriced(tx model.AbstractTransaction, local *accountSet) bool {
	// Local transactions cannot be under feeList
	if local.containsTx(tx) {
		return false
	}
	// Discard stale price points if found at the heap start
	for len(*l.items) > 0 {
		head := []model.AbstractTransaction(*l.items)[0]
		if l.all.Get(head.CalTxId()) == nil {
			l.stales--
			heap.Pop(l.items)
			continue
		}
		break
	}
	// Check if the transaction is under feeList or not
	if len(*l.items) == 0 {
		log.Error("Pricing query for empty pool") // This cannot happen, print to catch programming errors
		return false
	}
	cheapest := []model.AbstractTransaction(*l.items)[0]
	return cheapest.Fee().Cmp(tx.Fee()) >= 0
}

// Discard finds a number of most under feeList transactions, removes them from the
// feeList list and returns them for further removal from the entire pool.
func (l *txFeeList) Discard(count int, local *accountSet) []model.AbstractTransaction {
	drop := make([]model.AbstractTransaction, 0, count) // Remote under feeList transactions to drop
	save := make([]model.AbstractTransaction, 0, 64)    // Local under feeList transactions to keep

	for len(*l.items) > 0 && count > 0 {
		// Discard stale transactions if found during cleanup
		tx := heap.Pop(l.items).(model.AbstractTransaction)
		if l.all.Get(tx.CalTxId()) == nil {
			l.stales--
			continue
		}
		// Non stale transaction found, discard unless local
		if local.containsTx(tx) {
			save = append(save, tx)
		} else {
			drop = append(drop, tx)
			count--
		}
	}
	for _, tx := range save {
		heap.Push(l.items, tx)
	}
	return drop
}
