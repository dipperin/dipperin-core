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

package statemachine

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/model"
)

//go:generate mockgen -destination=./verification_mock_test.go -package=statemachine github.com/dipperin/dipperin-core/core/model AbstractVerification
//go:generate mockgen -destination=./block_mock_test.go -package=statemachine github.com/dipperin/dipperin-core/core/model AbstractBlock
//go:generate mockgen -destination=./state_mock_test.go -package=statemachine github.com/dipperin/dipperin-core/core/csbft/statemachine ChainReader,MsgSigner,MsgSender,Validator,Fetcher
type ChainReader interface {
	GetSeenCommit(height uint64) []model.AbstractVerification
	SaveBlock(block model.AbstractBlock, seenCommits []model.AbstractVerification) error
	CurrentBlock() model.AbstractBlock
	IsChangePoint(block model.AbstractBlock, isProcessPackageBlock bool) bool
	GetNextVerifiers() []common.Address
	GetCurrVerifiers() []common.Address
}

type MsgSigner interface {
	SignHash(hash []byte) ([]byte, error)
	GetAddress() common.Address
}

type MsgSender interface {
	BroadcastMsg(msgCode uint64, msg interface{})
	SendReqRoundMsg(msgCode uint64, from []common.Address, msg interface{})
	//iblt
	BroadcastEiBlock(block model.AbstractBlock)
}
type Validator interface {
	FullValid(block model.AbstractBlock) error
}

type Fetcher interface {
	FetchBlock(from common.Address, blockHash common.Hash) model.AbstractBlock
}
