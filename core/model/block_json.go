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
	"encoding/json"
	"errors"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/hexutil"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/core/bloom"
	"math/big"
)

func (h Header) MarshalJSON() ([]byte, error) {
	type Header struct {
		Version          hexutil.Uint64    `json:"version"  gencodec:"required"`
		Number           hexutil.Uint64    `json:"number"  gencodec:"required"`
		PreHash          common.Hash       `json:"pre_hash"  gencodec:"required"`
		Seed             common.Hash       `json:"seed"  gencodec:"required"`
		Diff             common.Difficulty `json:"diff"  gencodec:"required"`
		TimeStamp        *hexutil.Big      `json:"timestamp"  gencodec:"required"`
		CoinBase         common.Address    `json:"coinbase"  gencodec:"required"`
		Nonce            common.BlockNonce `json:"nonce"  gencodec:"required"`
		Bloom            iblt.BloomRLP     `json:"Bloom"        gencodec:"required"`
		TransactionRoot  common.Hash       `json:"txs_root"   gencodec:"required"`
		StateRoot        common.Hash       `json:"state_root" gencodec:"required"`
		HeaderRoot       common.Hash       `json:"header_root"  gencodec:"required"`
		VerificationRoot common.Hash       `json:"verification_root"  gencodec:"required"`
		InterLinksRoot   common.Hash       `json:"interlink_root"  gencodec:"required"`
	}
	var enc Header
	enc.Version = hexutil.Uint64(h.Version)
	enc.Number = hexutil.Uint64(h.Number)
	enc.PreHash = h.PreHash
	enc.Seed = h.Seed
	enc.Diff = h.Diff
	enc.TimeStamp = (*hexutil.Big)(h.TimeStamp)
	enc.CoinBase = h.CoinBase
	enc.Nonce = h.Nonce
	enc.Bloom = *h.Bloom.BloomRLP()
	enc.TransactionRoot = h.TransactionRoot
	enc.StateRoot = h.StateRoot
	enc.VerificationRoot = h.VerificationRoot
	enc.InterLinksRoot = h.InterlinkRoot
	return json.Marshal(&enc)
}

func (h *Header) UnmarshalJSON(input []byte) error {
	type Header struct {
		Version          *hexutil.Uint64    `json:"version"  gencodec:"required"`
		Number           *hexutil.Uint64    `json:"number"  gencodec:"required"`
		PreHash          *common.Hash       `json:"pre_hash"  gencodec:"required"`
		Seed             *common.Hash       `json:"seed"  gencodec:"required"`
		Diff             *common.Difficulty `json:"diff"  gencodec:"required"`
		TimeStamp        *hexutil.Big       `json:"timestamp"  gencodec:"required"`
		CoinBase         *common.Address    `json:"coinbase"  gencodec:"required"`
		Nonce            *common.BlockNonce `json:"nonce"  gencodec:"required"`
		Bloom            *iblt.BloomRLP     `json:"Bloom"        gencodec:"required"`
		TransactionRoot  *common.Hash       `json:"txs_root"   gencodec:"required"`
		StateRoot        *common.Hash       `json:"state_root" gencodec:"required"`
		HeaderRoot       *common.Hash       `json:"header_root"  gencodec:"required"`
		VerificationRoot *common.Hash       `json:"verification_root"  gencodec:"required"`
		InterLinksRoot   *common.Hash       `json:"interlink_root"  gencodec:"required"`
	}
	var dec Header
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if dec.Version == nil {
		return errors.New("missing required field 'version' for Header")
	}
	h.Version = uint64(*dec.Version)
	if dec.Number == nil {
		return errors.New("missing required field 'number' for Header")
	}
	h.Number = uint64(*dec.Number)
	if dec.PreHash == nil {
		return errors.New("missing required field 'prehash' for Header")
	}
	h.PreHash = *dec.PreHash
	if dec.Seed == nil {
		return errors.New("missing required field 'seed' for Header")
	}
	h.Seed = *dec.Seed
	if dec.Diff == nil {
		return errors.New("missing required field 'difficulty' for Header")
	}
	h.Diff = *dec.Diff
	if dec.TimeStamp == nil {
		return errors.New("missing required field 'timestamp' for Header")
	}
	h.TimeStamp = (*big.Int)(dec.TimeStamp)
	if dec.CoinBase == nil {
		return errors.New("missing required field 'coinbase' for Header")
	}
	h.CoinBase = *dec.CoinBase
	if dec.Nonce == nil {
		return errors.New("missing required field 'nonce' for Header")
	}
	h.Nonce = *dec.Nonce
	if dec.Bloom == nil {
		return errors.New("missing required field 'bloom' for Header")
	}
	h.Bloom = iblt.NewBloom(dec.Bloom.Config)
	dec.Bloom.CBloom(h.Bloom)
	if dec.TransactionRoot == nil {
		return errors.New("missing required field 'transactionRoot' for Header")
	}
	h.TransactionRoot = *dec.TransactionRoot
	if dec.StateRoot == nil {
		return errors.New("missing required field 'stateRoot ' for Header")
	}
	h.StateRoot = *dec.StateRoot
	if dec.VerificationRoot == nil {
		return errors.New("missing required field 'verificationRoot' for Header")
	}
	h.VerificationRoot = *dec.VerificationRoot
	if dec.InterLinksRoot == nil {
		return errors.New("missing required field 'interlinkRoot' for Header")
	}
	h.InterlinkRoot = *dec.InterLinksRoot
	return nil
}

func (b Block) MarshalJSON() ([]byte, error) {
	type block struct {
		Header *Header `json:"header"  gencodec:"required"`
		Body   *Body   `json:"body"  gencodec:"required"`
	}
	var enc block
	enc.Header = b.header
	enc.Body = b.body
	return json.Marshal(&enc)
}

func (b *Block) UnmarshalJSON(input []byte) error {
	type block struct {
		Header *Header `json:"header"  gencodec:"required"`
		Body   *Body   `json:"body"  gencodec:"required"`
	}
	var dec block
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if dec.Header == nil {
		return errors.New("missing header")
	}
	b.header = dec.Header
	if dec.Body == nil {
		return errors.New("missing body")
	}
	b.body = dec.Body
	return nil
}

//func (b Body) MarshalJSON() ([]byte, error) {
//
//}

func (b *Body) UnmarshalJSON(input []byte) error {
	return blockJsonHandler.DecodeBody(b, input)
}


var blockJsonHandler BlockJsonHandler

func init() {
	blockJsonHandler = &PBFTBlockJsonHandler{}
}

func SetBlockJsonHandler(h BlockJsonHandler) {
	blockJsonHandler = h
}

// only have interface use diff rlp handler
type BlockJsonHandler interface {
	DecodeBody(to *Body, input []byte) error
}

type PBFTBlockJsonHandler struct {}

func (h *PBFTBlockJsonHandler) DecodeBody(to *Body, input []byte) error {
	var from PBFTBody

	//log.Debug("PBFTBlockJsonHandler DecodeBody running~~~~~~~")

	if err := util.ParseJsonFromBytes(input, &from); err != nil {
		return err
	}

	to.Txs = from.Txs
	to.Vers = make([]AbstractVerification, len(from.Vers))
	util.InterfaceSliceCopy(to.Vers, from.Vers)
	to.Inters = from.Inters
	return nil
}


