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
	"github.com/dipperin/dipperin-core/common"
	"testing"
	"github.com/stretchr/testify/assert"
	"time"
	"math/big"
	"fmt"
)

func TestGetTarget(t *testing.T) {

	// CalNewWorkDiff if testing env
	firstTime, _ := time.Parse("2006-01-02 15:04:05", "2018-09-20 10:32:45")
	lastTime, _ := time.Parse("2006-01-02 15:04:05", "2018-09-20 11:11:31")
	prevB := &Block{
		header: &Header{},
		body:   &Body{},
	}
	header1 := &Header{Number: 4319, Diff: common.HexToDiff("0x1effffff"), TimeStamp: big.NewInt(lastTime.Unix())}
	prevB.header = header1

	prevSpanB := &Block{header: &Header{},
		body: &Body{},
	}
	header2 := &Header{Number: 1, Diff: common.HexToDiff("0x1effffff"), TimeStamp: big.NewInt(firstTime.Unix())}
	prevSpanB.header = header2

	result := calNewWorkDiffByTime(prevSpanB.Timestamp(), prevB.Timestamp(), prevB.Difficulty())
	fmt.Println(result.Hex())

	//firstTime2, _ := time.Parse("2006-01-02 15:04:05", "2018-09-20 10:32:45")
	//lastTime2, _ := time.Parse("2006-01-02 15:04:05", "2018-09-20 10:33:31")
	//prevB2 := &Block{
	//	header: &Header{},
	//	body:   &Body{},
	//}
	//header11 := &Header{Number: 500, Diff: common.HexToDiff("0x1fffffff"), TimeStamp: big.NewInt(lastTime2.Unix())}
	//prevB2.header = header11
	//
	//prevSpanB2 := &Block{header: &Header{},
	//	body: &Body{},
	//}
	//header22 := &Header{Number: 1, Diff: common.HexToDiff("0x1fffffff"), TimeStamp: big.NewInt(firstTime2.Unix())}
	//prevSpanB2.header = header22
	//
	//result = calNewWorkDiffByTime(prevSpanB2.Timestamp(), prevB2.Timestamp(), prevB2.Difficulty())
	//fmt.Println(result.Hex())
}

func TestLastPeriodBlockNum(t *testing.T) {
	assert.Equal(t, 0, int(LastPeriodBlockNum(0)))
	// The new block num is an integer multiple of BlockCountOfPeriod, so it should return the value of the last cycle. Since the current block is passed here, it will be added by 1, so here it should be added by -1
	assert.Equal(t, 0, int(LastPeriodBlockNum(BlockCountOfPeriod-1)))
	assert.Equal(t, int(BlockCountOfPeriod), int(LastPeriodBlockNum(BlockCountOfPeriod)))
	assert.Equal(t, int(BlockCountOfPeriod), int(LastPeriodBlockNum(BlockCountOfPeriod+1)))
	assert.Equal(t, int(BlockCountOfPeriod), int(LastPeriodBlockNum(BlockCountOfPeriod*2-1)))
	assert.Equal(t, int(BlockCountOfPeriod*2), int(LastPeriodBlockNum(BlockCountOfPeriod*2)))
	assert.Equal(t, int(BlockCountOfPeriod*2), int(LastPeriodBlockNum(BlockCountOfPeriod*2+1)))
}

func TestNewCalNewWorkDiff(t *testing.T) {
	block1 := CreateBlock(5, common.HexToHash("123"), 2)
	block2 := CreateBlock(6, common.HexToHash("123"), 2)

	result := NewCalNewWorkDiff(block1, block2, 12)
	assert.Equal(t, block2.Difficulty(), result)

	result = NewCalNewWorkDiff(block1, block2, BlockCountOfPeriod-1)
	assert.NotEqual(t, block2.Difficulty(), result)

	IgnoreDifficultyValidation = true
	result = NewCalNewWorkDiff(block1, block2, 12)
	assert.Equal(t, common.HexToDiff("0x1fffffff"), result)
}
