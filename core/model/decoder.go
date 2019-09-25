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
	"github.com/ethereum/go-ethereum/rlp"
)

func MakeDefaultBlockDecoder() BlockDecoder {
	return &defaultBlockDecoder{}
}

type BlockDecoder interface {
	DecodeRlpHeaderFromBytes(data []byte) (AbstractHeader, error)
	DecodeRlpBodyFromBytes(data []byte) (AbstractBody, error)
	DecodeRlpBlockFromHeaderAndBodyBytes(headerB []byte, bodyB []byte) (AbstractBlock, error)
	DecodeRlpBlockFromBytes(data []byte) (AbstractBlock, error)
	DecodeRlpTransactionFromBytes(data []byte) (AbstractTransaction, error)
}

type defaultBlockDecoder struct{}

func (decoder *defaultBlockDecoder) DecodeRlpBlockFromHeaderAndBodyBytes(headerB []byte, bodyB []byte) (AbstractBlock, error) {
	var header Header
	if err := rlp.DecodeBytes(headerB, &header); err != nil {
		return nil, err
	}
	var body Body
	if err := rlp.DecodeBytes(bodyB, &body); err != nil {
		return nil, err
	}
	return &Block{header: &header, body: &body}, nil
}

func (decoder *defaultBlockDecoder) DecodeRlpHeaderFromBytes(data []byte) (AbstractHeader, error) {
	var header Header
	if err := rlp.DecodeBytes(data, &header); err != nil {
		return nil, err
	}
	return &header, nil
}

func (decoder *defaultBlockDecoder) DecodeRlpBodyFromBytes(data []byte) (AbstractBody, error) {
	var body Body
	if err := rlp.DecodeBytes(data, &body); err != nil {
		return nil, err
	}
	return &body, nil
}

func (decoder *defaultBlockDecoder) DecodeRlpBlockFromBytes(data []byte) (AbstractBlock, error) {
	var block Block
	if err := rlp.DecodeBytes(data, &block); err != nil {
		return nil, err
	}
	return &block, nil
}

func (decoder *defaultBlockDecoder) DecodeRlpTransactionFromBytes(data []byte) (AbstractTransaction, error) {
	var tx Transaction
	if err := rlp.DecodeBytes(data, &tx); err != nil {
		return nil, err
	}
	return &tx, nil
}
