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

package model

import (
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/stretchr/testify/assert"
	"math/big"
	"runtime"
	"testing"
)

func TestNewTxCacher(t *testing.T) {
	cacher := NewTxCacher(runtime.NumCPU())
	assert.NotNil(t, cacher)
}

func TestTxCacher_TxRecover(t *testing.T) {
	log.InitLogger(log.LvlDebug)
	//go  func() {debug.PrintStack()}()
	cacher := NewTxCacher(runtime.NumCPU())
	assert.NotNil(t, cacher)
	txs := []AbstractTransaction{}
	for i := 0; i < 2 ; i++ {
		txs = append(txs, CreateSignedTx(uint64(i), big.NewInt(10)))
	}
	log.Debug("TestTxCacher_TxRecover", "txs", txs[0].CalTxId())


	cacher.TxRecover(txs)

	log.Debug("TestTxCacher_TxRecover after ", "txs", txs)
	//cacher.TxRecover([]AbstractTransaction{})
	cacher.StopTxCacher()
}

func TestTxCacher_StopTxCacher(t *testing.T) {
	cacher := NewTxCacher(runtime.NumCPU())
	assert.NotNil(t, cacher)
	cacher.StopTxCacher()
}
