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
	"github.com/dipperin/dipperin-core/tests/g-testData"
	"testing"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/bloom"
	"time"
	"github.com/dipperin/dipperin-core/third-party/log"
	"math/big"
	"github.com/dipperin/dipperin-core/third-party/crypto"
)

type testPool struct {
	txs []model.AbstractTransaction
}

func (p *testPool) AddLocal(tx model.AbstractTransaction) error {
	p.txs = append(p.txs, tx)
	return nil
}

func (testPool) AddRemote(tx model.AbstractTransaction) error {
	panic("implement me")
}

func (testPool) AddLocals(txs []model.AbstractTransaction) []error {
	panic("implement me")
}

func (testPool) AddRemotes(txs []model.AbstractTransaction) []error {
	panic("implement me")
}

func (testPool) ConvertPoolToMap() map[common.Hash]model.AbstractTransaction {
	panic("implement me")
}

func (testPool) Stats() (int, int) {
	panic("implement me")
}

func (p *testPool) GetTxsEstimator(broadcastBloom *iblt.Bloom) *iblt.HybridEstimator {
	// hack
	c := iblt.NewHybridEstimatorConfig()
	estimator := iblt.NewHybridEstimator(c)
	// get peer local tx pool all txs

	startAt1 := time.Now()
	// get peer estimator
	for _, tx := range p.txs {
		b := tx.CalTxId().Bytes()
		if broadcastBloom.LookUp(b) {
			estimator.EncodeByte(b)
		}
	}
	log.Info("broadcastBloom.LookUp", "use", time.Now().Sub(startAt1))

	return estimator
}

func (testPool) Pending() (map[common.Address][]model.AbstractTransaction, error) {
	panic("implement me")
}

func (testPool) Queueing() (map[common.Address][]model.AbstractTransaction, error) {
	panic("implement me")
}

func TestTxPool_GetTxsEstimator(t *testing.T) {
	// 100 transactions in block
	bloom := getBloom(100)
	p := addTxPool(100)

	for i := 0; i < 10; i ++ {
		p.GetTxsEstimator(bloom)
	}
}

func addTxPool(n int) *testPool {
	p := &testPool{}

	for i := 0; i < 100; i++ {
		p.AddLocal(getTx())
	}

	return p
}

func getTx() model.AbstractTransaction {
	pk, _ := crypto.GenerateKey()
	return transaction(1, common.HexToAddress("0x123"), big.NewInt(123), big.NewInt(321),g_testData.TestGasLimit, pk)
}

func getBloom(n int) *iblt.Bloom {
	pk, _ := crypto.GenerateKey()

	bloom := iblt.NewBloom(iblt.DeriveBloomConfig(n))

	for i := uint64(0); i < uint64(n); i++ {
		tx := transaction(i, common.HexToAddress("0x123"), big.NewInt(123), big.NewInt(321),g_testData.TestGasLimit, pk)
		bloom.Digest(tx.CalTxId().Bytes())
	}

	return bloom
}
