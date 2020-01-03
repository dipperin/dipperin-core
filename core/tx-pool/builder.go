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

package tx_pool

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/log"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/core/model"
	"go.uber.org/zap"
	"runtime"
	"time"
)

func NewTxPool(config TxPoolConfig, chainConfig chain_config.ChainConfig, chain BlockChain) *TxPool {
	//todo need rethink this new method.
	pool := &TxPool{
		config:      config,
		chainConfig: chainConfig,
		chain:       chain,
		//PoolConsensus: consensus,
		signer: model.MakeSigner(&chainConfig, chain.CurrentBlock().Number()),
		minFee: config.MinFee,

		pending: make(map[common.Address]*txList),
		queue:   make(map[common.Address]*txList),
		beats:   make(map[common.Address]time.Time),

		all: newTxLookup(),
	}
	pool.locals = newAccountSet(pool.signer)
	pool.feeList = newTxFeeList(pool.all)
	pool.reset(nil, chain.CurrentBlock().Header().(*model.Header))

	// If local transactions and journaling is enabled, load from disk
	if !config.NoLocals && config.Journal != "" {
		pool.journal = newTxJournal(config.Journal)

		if err := pool.journal.load(pool.AddLocals); err != nil {
			log.DLogger.Warn("Failed to load transaction journal", zap.Error(err))
		}
		if err := pool.journal.rotate(pool.local()); err != nil {
			log.DLogger.Warn("Failed to rotate transaction journal", zap.Error(err))
		}
	}

	pool.loopStopCtrl = make(chan int)

	//transaction cacher
	pool.senderCacher = model.NewTxCacher(runtime.NumCPU())
	//pool.wg.Add(1)
	//go pool.loop()

	return pool

}
