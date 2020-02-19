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
	"bytes"
	"github.com/dipperin/dipperin-core/tests/g-testData"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestTransaction_MarshalJSON(t *testing.T) {
	txsent1, txsent2 := createTestTx()
	_, err1 := txsent1.MarshalJSON()
	_, err2 := txsent2.MarshalJSON()
	assert.NoError(t, err1)
	assert.NoError(t, err2)
}

func TestTransaction_UnmarshalJSON(t *testing.T) {
	txsent1, txsent2 := createTestTx()
	json1, err1 := txsent1.MarshalJSON()
	json2, err2 := txsent2.MarshalJSON()
	assert.NoError(t, err1)
	assert.NoError(t, err2)
	var txget1 Transaction
	var txget2 Transaction
	err3 := txget1.UnmarshalJSON(json1)
	assert.NoError(t, err3)
	err4 := txget2.UnmarshalJSON(json2)
	assert.NoError(t, err4)

	assert.Equal(t, txsent1.CalTxId(), txget1.CalTxId())
	assert.Equal(t, txsent1.wit.R.Cmp(txget1.wit.R), 0)
	assert.Equal(t, txsent1.wit.S.Cmp(txget1.wit.S), 0)
	assert.Equal(t, txsent1.wit.V.Cmp(txget1.wit.V), 0)
	assert.True(t, bytes.Equal(txsent1.wit.HashKey, txget1.wit.HashKey))

	assert.Equal(t, txsent2.CalTxId(), txget2.CalTxId())
	assert.Equal(t, txsent2.wit.R.Cmp(txget2.wit.R), 0)
	assert.Equal(t, txsent2.wit.S.Cmp(txget2.wit.S), 0)
	assert.Equal(t, txsent2.wit.V.Cmp(txget2.wit.V), 0)
	assert.True(t, bytes.Equal(txsent2.wit.HashKey, txget2.wit.HashKey))
}

func TestTxData_MarshalJSON(t *testing.T) {
	d := txData{Amount: big.NewInt(10), Price: g_testData.TestGasPrice, GasLimit: g_testData.TestGasLimit}
	result, err := d.MarshalJSON()
	assert.NotNil(t, result)
	assert.NoError(t, err)
}

func TestTxData_UnmarshalJSON(t *testing.T) {
	d := txData{Amount: big.NewInt(10), Price: g_testData.TestGasPrice, GasLimit: g_testData.TestGasLimit}
	result, err := d.MarshalJSON()
	assert.NotNil(t, result)
	assert.NoError(t, err)
	err = d.UnmarshalJSON(result)
	assert.NoError(t, err)
}

func TestWitness_MarshalJSON(t *testing.T) {
	w := witness{R: big.NewInt(10), S: big.NewInt(20), V: big.NewInt(100)}
	result, err := w.MarshalJSON()
	assert.NotNil(t, result)
	assert.NoError(t, err)
}

func TestWitness_UnmarshalJSON(t *testing.T) {
	w := witness{R: big.NewInt(10), S: big.NewInt(20), V: big.NewInt(100)}
	result, err := w.MarshalJSON()
	assert.NotNil(t, result)
	assert.NoError(t, err)
	err = w.UnmarshalJSON(result)
	assert.NoError(t, err)
}
