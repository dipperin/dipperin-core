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
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/ethereum/go-ethereum/rlp"
	"io"
)

var blockRlpHandler BlockRlpHandler

// maybe set handler outside
func SetBlockRlpHandler(h BlockRlpHandler) {
	blockRlpHandler = h
}

func init() {
	blockRlpHandler = &PBFTBlockRlpHandler{}
}

// only have interface use diff rlp handler
type BlockRlpHandler interface {
	DecodeBody(to *Body, s *rlp.Stream) error
}

type blockForRlp struct {
	Header *Header
	Body   *Body
}

func (b *Block) DecodeRLP(s *rlp.Stream) error {
	var eb blockForRlp
	_, size, _ := s.Kind()
	if err := s.Decode(&eb); err != nil {
		return err
	}
	b.header, b.body = eb.Header, eb.Body
	b.size.Store(common.StorageSize(rlp.ListSize(size)))
	return nil
}

// EncodeRLP serializes b into the Ethereum RLP block format.
func (b *Block) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, blockForRlp{
		Header: b.header,
		Body:   b.body,
	})
}


type PBFTBlockRlpHandler struct {}

func (h *PBFTBlockRlpHandler) DecodeBody(to *Body, s *rlp.Stream) error {
	var sBody PBFTBody
	if err := s.Decode(&sBody); err != nil {
		return err
	}
	to.Txs = sBody.Txs
	to.Vers = make([]AbstractVerification, len(sBody.Vers))
	util.InterfaceSliceCopy(to.Vers, sBody.Vers)
	to.Inters = sBody.Inters
	return nil
}

type PBFTBody struct {
	Txs  []*Transaction  `json:"transactions"`
	Vers []*VoteMsg `json:"commit_msg"`
	Inters InterLink	 `json:"interlinks"`
}

func (b *Body) DecodeRLP(s *rlp.Stream) error {
	return blockRlpHandler.DecodeBody(b, s)
}