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

import (
	"fmt"
	"github.com/dipperin/dipperin-core/core/bloom"
	"github.com/dipperin/dipperin-core/core/mine/minemsg"
	"github.com/dipperin/dipperin-core/core/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDefaultExecutorBuilder(t *testing.T) {
	eb := NewDefaultExecutorBuilder()
	assert.NotNil(t, eb)
}

func Test_defaultExecutorBuilder_CreateExecutor(t *testing.T) {
	eb := NewDefaultExecutorBuilder()

	wm := localWorkMsg{}

	ws := fakeWorkSubmitter{}

	_, err := eb.CreateExecutor(&wm, 2, &ws)

	assert.Error(t, err)

	wm2 := localWorkMsg{
		code: minemsg.NewDefaultWorkMsg,
		work: &minemsg.DefaultWork{BlockHeader: model.Header{Bloom: iblt.NewBloom(iblt.NewBloomConfig(1<<12, 4))}},
	}

	res, err := eb.CreateExecutor(&wm2, 2, &ws)

	assert.NoError(t, err)

	fmt.Println(res)
}
