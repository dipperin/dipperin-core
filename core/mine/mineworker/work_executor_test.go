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
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/mine/minemsg"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/tests/factory"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type fakeWorkSubmitter struct {
}

func (submitter *fakeWorkSubmitter) SubmitWork(work minemsg.Work) {
}

func TestDefaultWorkExecutor_ChangeNonce(t *testing.T) {
	diff := common.HexToDiff("0x1effffff")
	block := factory.CreateBlock2(diff, 1)

	work := &minemsg.DefaultWork{
		BlockHeader: *(block.Header().(*model.Header)),
	}
	work.CalBlockRlpWithoutNonce()
	executor := NewDefaultWorkExecutor(work, &fakeWorkSubmitter{})
	result := make(chan minemsg.DefaultWork, 1)
	tick := time.Tick(20 * time.Second)

	found := false
out:
	for {
		select {
		case rw := <-result:
			rw.FillSealResult(block)
			log.Info("result work", "nonce", rw.BlockHeader.Nonce, "block header nonce", block.Nonce().Hex())
			assert.True(t, block.Hash().ValidHashForDifficulty(diff), block.Hash().Hex())
			break out
		case <-tick:
			t.Fatal("seal timeout")
			break out
		default:
			if !found && executor.ChangeNonce() {
				found = true
				log.Info("found nonce")
				result <- *executor.curWork
			}
		}
	}
}

func TestBytesCopy(t *testing.T) {
	tmpNonce := common.BlockNonce{0, 0, 0, 0, 1, 0, 0, 1, 0}
	fmt.Println(tmpNonce[:8])
}
