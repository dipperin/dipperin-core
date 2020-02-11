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

package components

import (
	"github.com/dipperin/dipperin-core/common"
	model2 "github.com/dipperin/dipperin-core/core/csbft/model"
	"github.com/dipperin/dipperin-core/core/model"
)

//go:generate mockgen -destination=./verification_mock_test.go -package=components github.com/dipperin/dipperin-core/core/model AbstractVerification
//go:generate mockgen -destination=./block_mock_test.go -package=components github.com/dipperin/dipperin-core/core/model AbstractBlock
//go:generate mockgen -destination=./node_mock_test.go -package=components github.com/dipperin/dipperin-core/core/csbft/statemachine ChainReader,MsgSigner,MsgSender,Validator,Fetcher
type Fetcher interface {
	FetchBlock(from common.Address, blockHash common.Hash) model.AbstractBlock
}

type FetcherConn interface {
	SendFetchBlockMsg(msgCode uint64, from common.Address, msg *model2.FetchBlockReqDecodeMsg) error
}
