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

package mineworker

import "github.com/dipperin/dipperin-core/core/mine/minemsg"

// implement executorBuilder

func NewDefaultExecutorBuilder() *defaultExecutorBuilder {
	return &defaultExecutorBuilder{}
}

type defaultExecutorBuilder struct{}

func (builder *defaultExecutorBuilder) CreateExecutor(msg workMsg, workCount int, submitter workSubmitter) (result []workExecutor, err error) {
	switch msg.MsgCode() {
	case minemsg.NewDefaultWorkMsg:
		return builder.createDefaultWorkExecutor(msg, workCount, submitter)
	}
	return nil, UnknownMsgCodeErr
}

// create default work executor
func (builder *defaultExecutorBuilder) createDefaultWorkExecutor(msg workMsg, workCount int, submitter workSubmitter) (result []workExecutor, err error) {
	tmpWork := &minemsg.DefaultWork{}
	if err = msg.Decode(tmpWork); err != nil {
		return
	}
	//pre-calculate block rlp
	tmpWork.CalBlockRlpWithoutNonce()
	works := tmpWork.Split(workCount)
	for _, w := range works {
		result = append(result, NewDefaultWorkExecutor(w, submitter))
	}
	return
}
