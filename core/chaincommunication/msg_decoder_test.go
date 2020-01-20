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

package chaincommunication

import (
	"bytes"
	"github.com/dipperin/dipperin-core/common"
	iblt "github.com/dipperin/dipperin-core/core/bloom"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third_party/crypto"
	"github.com/dipperin/dipperin-core/third_party/p2p"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
	"math/big"
	"strings"
	"testing"
	"time"
)

func TestMakeDefaultMsgDecoder(t *testing.T) {
	assert.NotNil(t, MakeDefaultMsgDecoder())
}

func TestDefaultMsgDecoder_DecoderBlockMsg(t *testing.T) {
	decoder := MakeDefaultMsgDecoder()

	testCases := []struct {
		name           string
		givenAndExpect func() (p2p.Msg, struct {
			existErr bool
			block    *model.Block
		})
	}{
		{
			name: "invalid p2p message",
			givenAndExpect: func() (p2p.Msg, struct {
				existErr bool
				block    *model.Block
			}) {
				inputStr := "test msg"
				return p2p.Msg{
						Code:       0x1,
						Size:       uint32(len(inputStr)),
						Payload:    strings.NewReader(inputStr),
						ReceivedAt: time.Date(2012, 2, 2, 22, 33, 44, 0, time.Local),
					}, struct {
						existErr bool
						block    *model.Block
					}{existErr: true, block: nil}
			},
		},
		{
			name: "decode success",
			givenAndExpect: func() (p2p.Msg, struct {
				existErr bool
				block    *model.Block
			}) {
				header := &model.Header{Number: 1, PreHash: common.HexToHash("0x12312fa0929348"),
					Diff: common.HexToDiff("0x1effffff"), Bloom: iblt.NewBloom(model.DefaultBlockBloomConfig)}
				block := model.NewBlock(header, nil, nil)
				blockBts, _ := rlp.EncodeToBytes(block)
				return p2p.Msg{
						Code:       0x1,
						Size:       uint32(len(blockBts)),
						Payload:    bytes.NewReader(blockBts),
						ReceivedAt: time.Date(2012, 2, 2, 22, 33, 44, 0, time.Local),
					}, struct {
						existErr bool
						block    *model.Block
					}{existErr: false, block: block}
			},
		},
	}

	for _, tc := range testCases {
		given, expect := tc.givenAndExpect()

		block, err := decoder.DecoderBlockMsg(given)

		if expect.existErr {
			if !assert.Error(t, err) {
				t.Errorf("case:%s, existErr: true, got:nil", tc.name)
			}
			continue
		}
		if !assert.NoError(t, err) {
			t.Errorf("case:%s, existErr: false, got:%s", tc.name, err.Error())
			continue
		}
		if !assert.Equal(t, expect.block.Hash(), block.Hash()) {
			t.Errorf("case:%s, block:%+v, got:%+v", tc.name, expect.block, block)
		}
	}
}

func TestDefaultMsgDecoder_DecodeTxMsg(t *testing.T) {
	decoder := MakeDefaultMsgDecoder()

	testCases := []struct {
		name           string
		givenAndExpect func() (p2p.Msg, struct {
			existErr bool
			tx       *model.Transaction
		})
	}{
		{
			name: "invalid p2p message",
			givenAndExpect: func() (msg p2p.Msg, i struct {
				existErr bool
				tx       *model.Transaction
			}) {
				inputStr := "test msg"
				return p2p.Msg{
						Code:       0x1,
						Size:       uint32(len(inputStr)),
						Payload:    strings.NewReader(inputStr),
						ReceivedAt: time.Date(2012, 2, 2, 22, 33, 44, 0, time.Local),
					}, struct {
						existErr bool
						tx       *model.Transaction
					}{existErr: true, tx: nil}
			},
		},
		{
			name: "decode success",
			givenAndExpect: func() (msg p2p.Msg, i struct {
				existErr bool
				tx       *model.Transaction
			}) {
				tx := model.NewTransaction(10, common.HexToAddress("0121321432423534534534"),
					big.NewInt(10000), big.NewInt(1), 2*model.TxGas, []byte{})
				gasUsed, _ := model.IntrinsicGas(tx.ExtraData(), false, false)
				tx.PaddingActualTxFee(big.NewInt(0).Mul(big.NewInt(int64(gasUsed)), tx.GetGasPrice()))
				priKey, _ := crypto.HexToECDSA("289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232031")
				tx.SignTx(priKey, model.NewSigner(big.NewInt(1)))
				bts, _ := rlp.EncodeToBytes(tx)
				return p2p.Msg{
						Code:       0x1,
						Size:       uint32(len(bts)),
						Payload:    bytes.NewReader(bts),
						ReceivedAt: time.Date(2012, 2, 2, 22, 33, 44, 0, time.Local),
					}, struct {
						existErr bool
						tx       *model.Transaction
					}{existErr: false, tx: tx}
			},
		},
	}

	for _, tc := range testCases {
		given, expect := tc.givenAndExpect()

		tx, err := decoder.DecodeTxMsg(given)

		if expect.existErr {
			if !assert.Error(t, err) {
				t.Errorf("case:%s, existErr: true, got:nil", tc.name)
			}
			continue
		}
		if !assert.NoError(t, err) {
			t.Errorf("case:%s, existErr: false, got:%s", tc.name, err.Error())
			continue
		}
		if !assert.Equal(t, expect.tx.String(), tx.(*model.Transaction).String()) {
			t.Errorf("case:%s, transaction:%+v, got:%+v", tc.name, expect.tx, tx)
		}
	}
}

func TestDefaultMsgDecoder_DecodeTxsMsg(t *testing.T) {
	decoder := MakeDefaultMsgDecoder()

	testCases := []struct {
		name           string
		givenAndExpect func() (p2p.Msg, struct {
			existErr bool
			txs      []*model.Transaction
		})
	}{
		{
			name: "invalid p2p message",
			givenAndExpect: func() (msg p2p.Msg, i struct {
				existErr bool
				txs      []*model.Transaction
			}) {
				inputStr := "test msg"
				return p2p.Msg{
						Code:       0x1,
						Size:       uint32(len(inputStr)),
						Payload:    strings.NewReader(inputStr),
						ReceivedAt: time.Date(2012, 2, 2, 22, 33, 44, 0, time.Local),
					}, struct {
						existErr bool
						txs      []*model.Transaction
					}{existErr: true, txs: nil}
			},
		},
		{
			name: "decode success",
			givenAndExpect: func() (msg p2p.Msg, i struct {
				existErr bool
				txs      []*model.Transaction
			}) {
				tx := model.NewTransaction(10, common.HexToAddress("0121321432423534534534"),
					big.NewInt(10000), big.NewInt(1), 2*model.TxGas, []byte{})
				gasUsed, _ := model.IntrinsicGas(tx.ExtraData(), false, false)
				tx.PaddingActualTxFee(big.NewInt(0).Mul(big.NewInt(int64(gasUsed)), tx.GetGasPrice()))
				priKey, _ := crypto.HexToECDSA("289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232031")
				tx.SignTx(priKey, model.NewSigner(big.NewInt(1)))

				tx2 := model.NewTransaction(11, common.HexToAddress("0121321436623534534534"),
					big.NewInt(10000), big.NewInt(1), 2*model.TxGas, []byte{})
				gasUsed2, _ := model.IntrinsicGas(tx2.ExtraData(), false, false)
				tx2.PaddingActualTxFee(big.NewInt(0).Mul(big.NewInt(int64(gasUsed2)), tx2.GetGasPrice()))
				priKey2, _ := crypto.HexToECDSA("289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232032")
				tx2.SignTx(priKey2, model.NewSigner(big.NewInt(1)))

				txs := []*model.Transaction{tx, tx2}
				bts, _ := rlp.EncodeToBytes(txs)
				return p2p.Msg{
						Code:       0x1,
						Size:       uint32(len(bts)),
						Payload:    bytes.NewReader(bts),
						ReceivedAt: time.Date(2012, 2, 2, 22, 33, 44, 0, time.Local),
					}, struct {
						existErr bool
						txs      []*model.Transaction
					}{existErr: false, txs: txs}
			},
		},
	}

	for _, tc := range testCases {
		given, expect := tc.givenAndExpect()

		txs, err := decoder.DecodeTxsMsg(given)

		if expect.existErr {
			if !assert.Error(t, err) {
				t.Errorf("case:%s, existErr: true, got:nil", tc.name)
			}
			continue
		}
		if !assert.NoError(t, err) {
			t.Errorf("case:%s, existErr: false, got:%s", tc.name, err.Error())
			continue
		}
		for i := range expect.txs {
			if !assert.Equal(t, expect.txs[i].String(), txs[i].(*model.Transaction).String()) {
				t.Errorf("case:%s, transaction:%+v, got:%+v", tc.name, expect.txs[i], txs[i])
			}
		}
	}
}
