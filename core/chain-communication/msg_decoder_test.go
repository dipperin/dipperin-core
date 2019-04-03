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

package chain_communication

import (
	"bytes"
	"github.com/dipperin/dipperin-core/core/model"
	"testing"

	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/tests/factory"
	"github.com/dipperin/dipperin-core/third-party/p2p"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
)

func TestMakeDefaultMsgDecoder(t *testing.T) {
	assert.NotNil(t, MakeDefaultMsgDecoder())
}

func Test_defaultMsgDecoder_DecoderBlockMsg(t *testing.T) {
	md := MakeDefaultMsgDecoder()
	msg := p2p.Msg{
		Payload: bytes.NewReader([]byte{}),
	}

	_, err := md.DecoderBlockMsg(msg)

	assert.Error(t, err)

	diff := common.HexToDiff("0x1effffff")
	fakeBlock := factory.CreateBlock2(diff, 1)

	payload, _ := rlp.EncodeToBytes(fakeBlock)

	msg = p2p.Msg{
		Payload: bytes.NewReader(payload),
	}

	block, err := md.DecoderBlockMsg(msg)

	assert.NoError(t, err)
	assert.NotNil(t, block)

}

func Test_defaultMsgDecoder_DecodeTxMsg(t *testing.T) {
	md := MakeDefaultMsgDecoder()
	msg := p2p.Msg{
		Payload: bytes.NewReader([]byte{}),
	}
	_, err := md.DecodeTxMsg(msg)

	assert.Error(t, err)

	tx, _ := factory.CreateTestTx()

	payload, _ := rlp.EncodeToBytes(tx)

	msg = p2p.Msg{
		Payload: bytes.NewReader(payload),
	}

	decodeTx, err := md.DecodeTxMsg(msg)

	assert.NoError(t, err)
	assert.NotNil(t, decodeTx)
}

func Test_defaultMsgDecoder_DecodeTxsMsg(t *testing.T) {
	md := MakeDefaultMsgDecoder()
	msg := p2p.Msg{
		Payload: bytes.NewReader([]byte{}),
	}
	_, err := md.DecodeTxsMsg(msg)

	assert.Error(t, err)

	tx, _ := factory.CreateTestTx()

	payload, _ := rlp.EncodeToBytes([]*model.Transaction{tx})

	msg = p2p.Msg{
		Payload: bytes.NewReader(payload),
	}

	decodeTxs, err := md.DecodeTxsMsg(msg)

	assert.NoError(t, err)
	assert.Equal(t, len(decodeTxs), 1)
}
