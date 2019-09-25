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

package vm_log_search

import (
	"context"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/bitutil"
	"github.com/dipperin/dipperin-core/core/chain/chaindb"
	"github.com/dipperin/dipperin-core/core/cs-chain/chain-writer/middleware"
	"github.com/dipperin/dipperin-core/core/model"
	model2 "github.com/dipperin/dipperin-core/core/vm/model"
	"github.com/ethereum/go-ethereum/ethdb"
	"time"
)

const (
	// BloomServiceThreads is the number of goroutines used globally by an Ethereum
	// instance to service bloombits lookups for all running filters.
	//BloomServiceThreads = 16
	BloomServiceThreads = 8

	// bloomFilterThreads is the number of goroutines used locally per filter to
	// multiplex requests onto the global servicing goroutines.
	bloomFilterThreads = 3

	// bloomRetrievalBatch is the maximum number of bloom bit retrievals to service
	// in a single batch.
	bloomRetrievalBatch = 16

	// bloomRetrievalWait is the maximum time to wait for enough bloom bit requests
	// to accumulate request an entire batch (avoiding hysteresis).
	bloomRetrievalWait = time.Duration(0)

	// BloomBitsBlocks is the number of blocks a single bloom bit section vector
	// contains on the server side.
	BloomBitsBlocks uint64 = 4096

	// BloomConfirms is the number of confirmation blocks before a bloom section is
	// considered probably final and its rotated bits are calculated.
	BloomConfirms = 256
)

const (
	// bloomThrottling is the time to wait between processing two consecutive index
	// sections. It's useful during chain upgrades to prevent disk overload.
	bloomThrottling = 100 * time.Millisecond
)

// BloomIndexer implements a core.ChainIndexer, building up a rotated bloom bits index
// for the Ethereum header bloom filters, permitting blazing fast filtering.
type BloomIndexer struct {
	size uint64         // section size to generate bloombits for
	db   ethdb.Database // database instance to write index data and metadata into
	//chainReader middleware.ChainInterface
	gen     *Generator  // generator to rotate the bloom bits crating the bloom index
	section uint64      // Section is the section number being processed currently
	head    common.Hash // Head is the hash of the last header processed
}

// NewBloomIndexer returns a chain indexer that generates bloom bits data for the
// canonical chain for fast logs filtering.
func NewBloomIndexer(chainReader middleware.ChainInterface, db ethdb.Database, size, confirms uint64) *ChainIndexer {
	backend := &BloomIndexer{
		db:   db,
		size: size,
	}
	table := ethdb.NewTable(db, string(chaindb.BloomBitsIndexPrefix))

	return NewChainIndexer(chainReader, db, table, backend, size, confirms, bloomThrottling, "bloombits")
}

// Reset implements core.ChainIndexerBackend, starting a new bloombits index
// section.
func (b *BloomIndexer) Reset(ctx context.Context, section uint64, lastSectionHead common.Hash) error {
	gen, err := NewGenerator(uint(b.size))
	b.gen, b.section, b.head = gen, section, common.Hash{}
	return err
}

// Process implements core.ChainIndexerBackend, adding a new header's bloom into
// the index.
func (b *BloomIndexer) Process(ctx context.Context, header *model.AbstractHeader, bloom model2.Bloom) error {
	b.gen.AddBloom(uint((*header).GetNumber()-b.section*b.size), bloom)
	b.head = (*header).Hash()
	return nil
}

// Commit implements core.ChainIndexerBackend, finalizing the bloom section and
// writing it out into the database.
func (b *BloomIndexer) Commit() error {
	batch := b.db.NewBatch()
	for i := 0; i < model2.BloomBitLength; i++ {
		bits, err := b.gen.Bitset(uint(i))
		if err != nil {
			return err
		}
		chaindb.BatchSaveBloomBits(batch, b.head, uint(i), b.section, bitutil.CompressBytes(bits))
		//if err := batch.Put(bloomBitsKey(bit, section, head), bits); err != nil {
		//	log.Error("Failed to store bloom bits", "err", err)
		//}
	}
	return batch.Write()
}
