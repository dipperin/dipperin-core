package chain_state

import (
	"context"
	"errors"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/cs-chain/chain-writer/middleware"
	"github.com/dipperin/dipperin-core/core/model"
	model2 "github.com/dipperin/dipperin-core/core/vm/model"
	"github.com/ethereum/go-ethereum/ethdb"
	"math/big"
)

// Filter can be used to retrieve and filter logs.
type Filter struct {
	ChainReader middleware.ChainInterface
	ChainIndex  ChainIndexer

	db        ethdb.Database
	addresses []common.Address
	topics    [][]common.Hash

	block      common.Hash // Block hash if filtering a single block
	begin, end int64       // Range interval if filtering multiple blocks

	matcher *Matcher
}

// NewRangeFilter creates a new filter which uses a bloom filter on blocks to
// figure out whether a particular block is interesting or not.
func NewRangeFilter(chainReader middleware.ChainInterface, chainIndex ChainIndexer, begin, end int64, addresses []common.Address, topics [][]common.Hash) *Filter {
	// Flatten the address and topic filter clauses into a single bloombits filter
	// system. Since the bloombits are not positional, nil topics are permitted,
	// which get flattened into a nil byte slice.
	var filters [][][]byte
	if len(addresses) > 0 {
		filter := make([][]byte, len(addresses))
		for i, address := range addresses {
			filter[i] = address.Bytes()
		}
		filters = append(filters, filter)
	}
	for _, topicList := range topics {
		filter := make([][]byte, len(topicList))
		for i, topic := range topicList {
			filter[i] = topic.Bytes()
		}
		filters = append(filters, filter)
	}
	size := BloomBitsBlocks

	// Create a generic filter and convert it into a range filter
	filter := newFilter(chainReader, chainIndex, addresses, topics)

	filter.matcher = NewMatcher(size, filters)
	filter.begin = begin
	filter.end = end

	return filter
}

// NewBlockFilter creates a new filter which directly inspects the contents of
// a block to figure out whether it is interesting or not.
func NewBlockFilter(chainIndex ChainIndexer, chainReader middleware.ChainInterface, block common.Hash, addresses []common.Address, topics [][]common.Hash) *Filter {
	// Create a generic filter and convert it into a block filter
	filter := newFilter(chainReader, chainIndex, addresses, topics)
	filter.block = block
	return filter
}

// newFilter creates a generic filter that can either filter based on a block hash,
// or based on range queries. The search criteria needs to be explicitly set.
func newFilter(chainReader middleware.ChainInterface, chainIndex ChainIndexer, addresses []common.Address, topics [][]common.Hash) *Filter {
	return &Filter{
		ChainReader: chainReader,
		addresses:   addresses,
		topics:      topics,
		//db:        backend.ChainDb(),
		ChainIndex: chainIndex,
	}
}

// Logs searches the blockchain for matching log entries, returning all from the
// first block that contains matches, updating the start of the filter accordingly.
func (f *Filter) Logs(ctx context.Context) ([]*model2.Log, error) {
	// If we're doing singleton block filtering, execute and return
	if f.block != (common.Hash{}) {
		// header, err := f.backend.HeaderByHash(ctx, f.block)
		header := f.ChainReader.GetHeaderByHash(f.block)
		if header == nil {
			return nil, errors.New("unknown block")
		}
		return f.blockLogs(ctx, &header)
	}
	// Figure out the limits of the filter range
	// header, _ := f.backend.HeaderByNumber(ctx, rpc.LatestBlockNumber)
	block := f.ChainReader.GetLatestNormalBlock()
	if block == nil {
		return nil, nil
	}
	head := block.Header().GetNumber()

	if f.begin == -1 {
		f.begin = int64(head)
	}
	end := uint64(f.end)
	if f.end == -1 {
		end = head
	}
	// Gather all indexed logs, and finish with non indexed ones
	var (
		logs []*model2.Log
		err  error
	)
	size := BloomBitsBlocks
	sections, _, _ := f.ChainIndex.Sections()
	if indexed := sections * size; indexed > uint64(f.begin) {
		if indexed > end {
			logs, err = f.indexedLogs(ctx, end)
		} else {
			logs, err = f.indexedLogs(ctx, indexed-1)
		}
		if err != nil {
			return logs, err
		}
	}
	rest, err := f.unindexedLogs(ctx, end)
	logs = append(logs, rest...)
	return logs, err
}

// indexedLogs returns the logs matching the filter criteria based on the bloom
// bits indexed available locally or via the network.
func (f *Filter) indexedLogs(ctx context.Context, end uint64) ([]*model2.Log, error) {
	// Create a matcher session and request servicing from the backend
	matches := make(chan uint64, 64)

	session, err := f.matcher.Start(ctx, uint64(f.begin), end, matches)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	f.ChainIndex.ServiceFilter(ctx, session)

	// Iterate over the matches until exhausted or context closed
	var logs []*model2.Log

	for {
		select {
		case number, ok := <-matches:
			// Abort if all matches have been fulfilled
			if !ok {
				err := session.Error()
				if err == nil {
					f.begin = int64(end) + 1
				}
				return logs, err
			}
			f.begin = int64(number) + 1

			// Retrieve the suggested block and pull any truly matching logs
			// header, err := f.backend.HeaderByNumber(ctx, rpc.BlockNumber(number))
			header := f.ChainReader.GetHeaderByNumber(number)
			if header == nil {
				return logs, err
			}
			found, err := f.checkMatches(ctx, &header)
			if err != nil {
				return logs, err
			}
			logs = append(logs, found...)

		case <-ctx.Done():
			return logs, ctx.Err()
		}
	}
}

// indexedLogs returns the logs matching the filter criteria based on raw block
// iteration and bloom matching.
func (f *Filter) unindexedLogs(ctx context.Context, end uint64) ([]*model2.Log, error) {
	var logs []*model2.Log

	for ; f.begin <= int64(end); f.begin++ {
		header := f.ChainReader.GetHeaderByNumber(uint64(f.begin))
		if header == nil {
			return logs, nil
		}
		found, err := f.blockLogs(ctx, &header)
		if err != nil {
			return logs, err
		}
		logs = append(logs, found...)
	}
	return logs, nil
}

// blockLogs returns the logs matching the filter criteria within a single block.
func (f *Filter) blockLogs(ctx context.Context, header *model.AbstractHeader) (logs []*model2.Log, err error) {
	if bloomFilter((*header).GetBloomLog(), f.addresses, f.topics) {
		found, err := f.checkMatches(ctx, header)
		if err != nil {
			return logs, err
		}
		logs = append(logs, found...)
	}
	return logs, nil
}

func (f *Filter) GetLogs(header *model.AbstractHeader) [][]*model2.Log {
	receipts := f.ChainReader.GetReceipts((*header).Hash(), (*header).GetNumber())
	if len(receipts) > 0 {
		logs := make([][]*model2.Log, len(receipts))
		for i, receipt := range receipts {
			logs[i] = receipt.Logs
		}
		return logs
	}
	return nil
}

// checkMatches checks if the receipts belonging to the given header contain any log events that
// match the filter criteria. This function is called when the bloom filter signals a potential match.
func (f *Filter) checkMatches(ctx context.Context, header *model.AbstractHeader) (logs []*model2.Log, err error) {
	// Get the logs of the block
	//logsList, err := f.backend.GetLogs(ctx, header)
	logsList := f.GetLogs(header)
	if logsList == nil {
		return nil, err
	}
	var unfiltered []*model2.Log
	for _, logs := range logsList {
		unfiltered = append(unfiltered, logs...)
	}
	logs = filterLogs(unfiltered, nil, nil, f.addresses, f.topics)
	if len(logs) > 0 {
		// We have matching logs, check if we need to resolve full logs via the light client
		if logs[0].TxHash == (common.Hash{}) {
			receipts := f.ChainReader.GetReceipts((*header).Hash(), (*header).GetNumber())
			if len(receipts) <= 0 {
				return nil, err
			}
			unfiltered = unfiltered[:0]
			for _, receipt := range receipts {
				unfiltered = append(unfiltered, receipt.Logs...)
			}
			logs = filterLogs(unfiltered, nil, nil, f.addresses, f.topics)
		}
		return logs, nil
	}
	return nil, nil
}

func includes(addresses []common.Address, a common.Address) bool {
	for _, addr := range addresses {
		if addr == a {
			return true
		}
	}

	return false
}

// filterLogs creates a slice of logs matching the given criteria.
func filterLogs(logs []*model2.Log, fromBlock, toBlock *big.Int, addresses []common.Address, topics [][]common.Hash) []*model2.Log {
	var ret []*model2.Log
Logs:
	for _, log := range logs {
		if fromBlock != nil && fromBlock.Int64() >= 0 && fromBlock.Uint64() > log.BlockNumber {
			continue
		}
		if toBlock != nil && toBlock.Int64() >= 0 && toBlock.Uint64() < log.BlockNumber {
			continue
		}

		if len(addresses) > 0 && !includes(addresses, log.Address) {
			continue
		}
		// If the to filtered topics is greater than the amount of topics in logs, skip.
		// todo to be understand
		if len(topics) > len(log.Topics) {
			continue Logs
		}
		for i, sub := range topics {
			match := len(sub) == 0 // empty rule set == wildcard
			for _, topic := range sub {
				if log.Topics[i] == topic {
					match = true
					break
				}
			}
			if !match {
				continue Logs
			}
		}
		ret = append(ret, log)
	}
	return ret
}

func bloomFilter(bloom model2.Bloom, addresses []common.Address, topics [][]common.Hash) bool {
	if len(addresses) > 0 {
		var included bool
		for _, addr := range addresses {
			if model2.BloomLookup(bloom, addr) {
				included = true
				break
			}
		}
		if !included {
			return false
		}
	}

	for _, sub := range topics {
		included := len(sub) == 0 // empty rule set == wildcard
		for _, topic := range sub {
			if model2.BloomLookup(bloom, topic) {
				included = true
				break
			}
		}
		if !included {
			return false
		}
	}
	return true
}
