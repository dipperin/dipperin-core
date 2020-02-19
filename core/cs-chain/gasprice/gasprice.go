package gasprice

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/config"
	"github.com/dipperin/dipperin-core/common/consts"
	"github.com/dipperin/dipperin-core/common/g-error"
	"github.com/dipperin/dipperin-core/core/model"
	"math/big"
	"sort"
	"sync"
)

var maxPrice = big.NewInt(500 * consts.GDIPUNIT)

type GasPriceConfig struct {
	Blocks     int
	Percentile int
	Default    *big.Int `toml:",omitempty"`
}

var DefaultGasPriceConfig = GasPriceConfig{
	Blocks:     20,
	Percentile: 60,
	Default:    big.NewInt(config.DEFAULT_GAS_PRICE),
}

// Oracle recommends gas prices based on the content of recent
// blocks. Suitable for both light and full clients.
type Oracle struct {
	chainReader Chain
	lastBlock   common.Hash
	lastPrice   *big.Int
	cacheLock   sync.RWMutex
	fetchLock   sync.Mutex

	checkBlocks, maxEmpty, maxBlocks int
	percentile                       int
}

// NewOracle returns a new oracle.
func NewOracle(chainReader Chain, params GasPriceConfig) *Oracle {
	blocks := params.Blocks
	if blocks < 1 {
		blocks = 1
	}
	percent := params.Percentile
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	return &Oracle{
		chainReader: chainReader,
		lastPrice:   params.Default,
		checkBlocks: blocks,
		maxEmpty:    blocks / 2,
		maxBlocks:   blocks * 5,
		percentile:  percent,
	}
}

// SuggestPrice returns the recommended gas price.
func (gpo *Oracle) SuggestPrice() (*big.Int, error) {
	gpo.cacheLock.RLock()
	lastBlock := gpo.lastBlock
	lastPrice := gpo.lastPrice
	gpo.cacheLock.RUnlock()

	block := gpo.chainReader.CurrentBlock()
	blockHash := block.Hash()
	if block.Hash() == lastBlock {
		return lastPrice, nil
	}

	gpo.fetchLock.Lock()
	defer gpo.fetchLock.Unlock()

	// try checking the cache again, maybe the last fetch fetched what we need
	gpo.cacheLock.RLock()
	lastBlock = gpo.lastBlock
	lastPrice = gpo.lastPrice
	gpo.cacheLock.RUnlock()
	if blockHash == lastBlock {
		return lastPrice, nil
	}

	blockNum := block.Number()
	ch := make(chan getBlockPricesResult, gpo.checkBlocks)
	sent := 0
	exp := 0
	var blockPrices []*big.Int
	for sent < gpo.checkBlocks && blockNum > 0 {
		go gpo.getBlockPrices(blockNum, ch)
		sent++
		exp++
		blockNum--
	}
	maxEmpty := gpo.maxEmpty
	for exp > 0 {
		res := <-ch
		if res.err != nil {
			return lastPrice, res.err
		}
		exp--
		if res.price != nil {
			blockPrices = append(blockPrices, res.price)
			continue
		}
		if maxEmpty > 0 {
			maxEmpty--
			continue
		}
		if blockNum > 0 && sent < gpo.maxBlocks {
			go gpo.getBlockPrices(blockNum, ch)
			sent++
			exp++
			blockNum--
		}
	}
	price := lastPrice
	if len(blockPrices) > 0 {
		sort.Sort(bigIntArray(blockPrices))
		price = blockPrices[(len(blockPrices)-1)*gpo.percentile/100]
	}
	if price.Cmp(maxPrice) > 0 {
		price = new(big.Int).Set(maxPrice)
	}

	gpo.cacheLock.Lock()
	gpo.lastBlock = blockHash
	gpo.lastPrice = price
	gpo.cacheLock.Unlock()
	return price, nil
}

type getBlockPricesResult struct {
	price *big.Int
	err   error
}

type transactionsByGasPrice []*model.Transaction

func (t transactionsByGasPrice) Len() int      { return len(t) }
func (t transactionsByGasPrice) Swap(i, j int) { t[i], t[j] = t[j], t[i] }
func (t transactionsByGasPrice) Less(i, j int) bool {
	return t[i].GetGasPrice().Cmp(t[j].GetGasPrice()) < 0
}

// getBlockPrices calculates the lowest transaction gas price in a given block
// and sends it to the result channel. If the block is empty, price is nil.
func (gpo *Oracle) getBlockPrices(blockNum uint64, ch chan getBlockPricesResult) {
	block := gpo.chainReader.GetBlockByNumber(blockNum)
	if block == nil {
		ch <- getBlockPricesResult{nil, g_error.BlockIsNilError}
		return
	}

	blockTxs := block.GetTransactions()
	txs := make([]*model.Transaction, len(blockTxs))
	copy(txs, blockTxs)
	sort.Sort(transactionsByGasPrice(txs))
	for _, tx := range txs {
		sender, err := tx.Sender(nil)
		if err == nil && sender != block.CoinBaseAddress() {
			ch <- getBlockPricesResult{tx.GetGasPrice(), nil}
			return
		}
	}
	ch <- getBlockPricesResult{nil, nil}
}

type bigIntArray []*big.Int

func (s bigIntArray) Len() int           { return len(s) }
func (s bigIntArray) Less(i, j int) bool { return s[i].Cmp(s[j]) < 0 }
func (s bigIntArray) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
