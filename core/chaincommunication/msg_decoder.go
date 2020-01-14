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
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third_party/p2p"
)

func MakeDefaultMsgDecoder() *defaultMsgDecoder {
	return &defaultMsgDecoder{}
}

type defaultMsgDecoder struct{}

func (decoder *defaultMsgDecoder) DecoderBlockMsg(msg p2p.Msg) (model.AbstractBlock, error) {
	var block model.Block
	if err := msg.Decode(&block); err != nil {
		return nil, err
	}
	return &block, nil
}

func (decoder *defaultMsgDecoder) DecodeTxMsg(msg p2p.Msg) (model.AbstractTransaction, error) {
	var tx model.Transaction
	if err := msg.Decode(&tx); err != nil {
		return nil, err
	}
	return &tx, nil
}

func (decoder *defaultMsgDecoder) DecodeTxsMsg(msg p2p.Msg) (result []model.AbstractTransaction, err error) {

	var txs []*model.Transaction
	if err = msg.Decode(&txs); err != nil {
		return
	}
	result = make([]model.AbstractTransaction, len(txs))
	util.InterfaceSliceCopy(result, txs)
	return
}
