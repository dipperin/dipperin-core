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

package cachedb

import (
	"errors"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third-party/log/pbft_log"
	"github.com/ethereum/go-ethereum/rlp"
)

// decode interface data
type CacheDataDecoder interface {
	DecodeSeenCommits(data []byte) ([]model.AbstractVerification, error)
}

type BFTCacheDataDecoder struct{}

func (d *BFTCacheDataDecoder) DecodeSeenCommits(data []byte) (result []model.AbstractVerification, err error) {

	if len(data) == 0 {
		return []model.AbstractVerification{}, errors.New("decode data is empty")
	}

	var from []*model.VoteMsg
	if err = rlp.DecodeBytes(data, &from); err != nil {
		pbft_log.Error("Decode seen commits error", "err", err)
		return []model.AbstractVerification{}, err
	}

	result = make([]model.AbstractVerification, len(from))
	util.InterfaceSliceCopy(result, from)
	pbft_log.Debug("Decode Seen commits", "length", len(result))
	return result, nil
}
