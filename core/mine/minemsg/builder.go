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

package minemsg

import (
	"encoding/binary"
	"github.com/dipperin/dipperin-core/common/log"
	"github.com/dipperin/dipperin-core/core/model"
	"go.uber.org/zap"
)

func MakeDefaultWorkBuilder() *DefaultWorkBuilder {
	return &DefaultWorkBuilder{}
}

type DefaultWorkBuilder struct{}

func (builder *DefaultWorkBuilder) BuildWorks(newBlock model.AbstractBlock, workerLen int) (workMsgCode int, works []Work) {
	if newBlock == nil {
		log.DLogger.Warn("DefaultWorkBuilder build works, but got nil block")
		return
	}
	// this build only build model.Block's work
	block := newBlock
	header := block.Header().(*model.Header)

	for i := 0; i < workerLen; i++ {
		newHeader := *header
		binary.BigEndian.PutUint32(newHeader.Nonce[:4], uint32(i))

		log.DLogger.Info("BuildWorks", zap.String("verRoot", newHeader.VerificationRoot.Hex()), zap.Any("register root", newHeader.RegisterRoot))
		works = append(works, &DefaultWork{BlockHeader: newHeader})
		log.DLogger.Debug("DefaultWorkBuilder#BuildWorks", zap.Any("newHeader.Nonce", newHeader.Nonce))
	}
	workMsgCode = NewDefaultWorkMsg
	return
}
