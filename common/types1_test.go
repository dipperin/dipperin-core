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

package common

import (
	"fmt"
	"github.com/dipperin/dipperin-core/common/hexutil"
	"github.com/dipperin/dipperin-core/common/log"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
	"math/big"
	"testing"
)

func TestBlockNonceFromInt(t *testing.T) {
	fmt.Println(BlockNonceFromInt(1))
}

func TestHash_ValidHashForDifficulty(t *testing.T) {
	hash := HexToHash("0xd0b3bc1c37915fb1d582dcfe98ea62328656141a099bbc46f6a3b9a38c01a9de")
	diff := HexToHash("0x0000566611000000000000000000000000000000000000000000000000000000")
	fmt.Println(hash.Cmp(diff))
}

//type tBig struct {
//	//X *big.Int `json:"x"`
//	X *CsBigInt `json:"x"`
//}

// test parsing data with big int
func TestParseJsonBigInt(t *testing.T) {
	//tStr := `{"x":"0x4140c78940f6a24fdffc78873d4490d2100000000000002"}`
	// used to check if the value is correct
	cStr := `[100000000000000000000000000000000000000000000000000000002]`
	var cI []*big.Int
	err := util.ParseJson(cStr, &cI)
	assert.NoError(t, err)

	//var tb tBig
	//err = util.ParseJson(tStr, &tb)
	//assert.NoError(t, err)
	//
	//assert.Equal(t, cI[0], tb.X.Int)
	//assert.Equal(t, tStr, util.StringifyJson(&tb))
}

func TestDifficulty_Big(t *testing.T) {
	fmt.Println((*hexutil.Big)(HexToDiff("0x201fffff").Big()).String())
	assert.Equal(t, big.NewInt(1), Big1)
	assert.Equal(t, big.NewInt(2), Big2)
	assert.Equal(t, big.NewInt(3), Big3)
	assert.Equal(t, big.NewInt(0), Big0)
	assert.Equal(t, big.NewInt(32), Big32)
	assert.Equal(t, big.NewInt(256), Big256)
	assert.Equal(t, big.NewInt(257), Big257)
}

func TestHexStringSameWithVM(t *testing.T) {
	log.InitLogger(log.LoggerConfig{
		Lvl:         zapcore.DebugLevel,
		FilePath:    "",
		Filename:    "",
		WithConsole: true,
		WithFile:    false,
	})
	HexStringSameWithVM("0x00005586B883Ec6dd4f8c26063E18eb4Bd228e59c3E9")
}
