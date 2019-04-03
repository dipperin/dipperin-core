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

package minemsg

import (
	"encoding/binary"
	"github.com/dipperin/dipperin-core/core/bloom"
	"testing"

	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/tests/factory"
	"github.com/stretchr/testify/assert"
)

func TestDefaultWork_FillSealResult(t *testing.T) {
	diff := common.HexToDiff("0x1effffff")
	block := factory.CreateBlock2(diff, 1)
	w := &DefaultWork{BlockHeader: model.Header{Number: 1, Nonce: common.BlockNonce{0, 1}}, ResultNonce: common.BlockNonce{0, 2}}
	err := w.FillSealResult(block)
	assert.NoError(t, err)
	assert.Equal(t, common.BlockNonce{0, 2}, block.Nonce())

	block2 := factory.CreateBlock2(diff, 2)
	w2 := &DefaultWork{BlockHeader: model.Header{Number: 1, Nonce: common.BlockNonce{0, 1}}, ResultNonce: common.BlockNonce{0, 2}}
	err = w2.FillSealResult(block2)
	assert.Error(t, err)
}

func TestDefaultWork_SetWorkerCoinbaseAddress(t *testing.T) {
	w := &DefaultWork{BlockHeader: model.Header{}}
	w.SetWorkerCoinbaseAddress(common.Address{11, 22})
	assert.Equal(t, common.Address{11, 22}, w.GetWorkerCoinbaseAddress())
}

func TestDefaultWork_Split(t *testing.T) {
	w := &DefaultWork{BlockHeader: model.Header{}}
	works := w.Split(5)
	for i, w := range works {
		x := binary.BigEndian.Uint32(w.BlockHeader.Nonce[4:8])
		assert.Equal(t, uint32(i), x)
	}
}

func TestDefaultWork_CalHash(t *testing.T) {
	w := &DefaultWork{BlockHeader: model.Header{Bloom: iblt.NewBloom(iblt.NewBloomConfig(1<<12, 4))}}
	_, err := w.CalHash()
	assert.Error(t, err)

	w.CalBlockRlpWithoutNonce()
	_,err = w.CalHash()
	assert.NoError(t, err)
}
