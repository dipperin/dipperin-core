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
	"errors"
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third-party/log"
	"time"
)

func NewRoundMsgWithSign(height uint64, round uint64, signer SignHashFunc, addr common.Address) *NewRoundMsg {
	msg := &NewRoundMsg{
		Height: height,
		Round:  round,
	}

	sign, err := signer(msg.Hash().Bytes())
	if err != nil {
		log.Warn("sign new round msg failed", "err", err)
		return nil
	}
	msg.Witness = &model.WitMsg{
		Address: addr,
		Sign:    sign,
	}

	return msg
}

type SignHashFunc func(hash []byte) ([]byte, error)

type NewRoundMsg struct {
	Height  uint64
	Round   uint64
	Witness *model.WitMsg
}

func (r NewRoundMsg) Hash() common.Hash {
	r.Witness = nil
	return common.RlpHashKeccak256(r)
}

func (r NewRoundMsg) Valid() error {
	if r.Witness == nil {
		return errors.New("new round msg witness can't be nil")
	}
	if err := r.Witness.Valid(r.Hash().Bytes()); err != nil {
		return err
	}
	return nil
}

func NewProposalWithSign(h, r uint64, blockID common.Hash, hashFunc SignHashFunc, addr common.Address) *Proposal {
	p := &Proposal{Height: h, Round: r, BlockID: blockID, Timestamp: time.Now()}
	sign, err := hashFunc(p.Hash().Bytes())
	if err != nil {
		log.Warn("sign new proposal msg failed", "err", err)
		return nil
	}
	p.Witness = &model.WitMsg{
		Sign:    sign,
		Address: addr,
	}
	return p
}

type Proposal struct {
	Height    uint64
	Round     uint64
	BlockID   common.Hash
	Timestamp time.Time
	Witness   *model.WitMsg
}

func (p Proposal) Hash() common.Hash {
	p.Witness = nil
	return common.RlpHashKeccak256(p)
}

func (p Proposal) ValidBlock(b model.AbstractBlock) error {
	if p.BlockID.IsEqual(b.Hash()) && p.Height == b.Number() {
		return nil
	}
	return errors.New(fmt.Sprintf("invalid proposal block, num: %v hash: %v, p h: %v p block id: %v", b.Number(), b.Hash(), p.Height, p.BlockID))
}
