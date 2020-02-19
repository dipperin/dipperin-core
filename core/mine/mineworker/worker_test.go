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

package mineworker

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultWorker_SetCoinbaseAddress(t *testing.T) {
	worker := worker{}
	assert.Equal(t, common.Address{}, worker.CurrentCoinbaseAddress())
	tmpAddr := common.HexToAddress("0xfffffffabbbbaaa123123fff")
	worker.SetCoinbaseAddress(tmpAddr)
	assert.Equal(t, tmpAddr, worker.CurrentCoinbaseAddress())
}

func TestNewDefaultWorker(t *testing.T) {
	tmpAddr := common.HexToAddress("0xfffffffabbbbaaa123123fff")
	worker := newWorker(tmpAddr, 2, nil)
	assert.Equal(t, tmpAddr, worker.CurrentCoinbaseAddress())
}
