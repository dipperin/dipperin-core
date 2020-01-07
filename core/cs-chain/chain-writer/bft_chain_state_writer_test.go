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

package chain_writer

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/core/cs-chain/chain-writer/middleware"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/tests/g-mockFile"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestNewBftChainWriter(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	mc := g_mockFile.NewMockChainInterface(controller)
	mb := NewMockAbstractBlock(controller)
	mb.EXPECT().IsSpecial().Return(true)
	mb.EXPECT().Version().Return(uint64(100))
	mc.EXPECT().GetChainConfig().Return(chain_config.GetChainConfig()).AnyTimes()

	assert.Error(t, NewBftChainWriter(&middleware.BftBlockContext{
		BlockContext: middleware.BlockContext{Block: mb, Chain: mc},
	}, mc).SaveBlock())

}

func TestBftChainWriter_SaveBlock(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	mc := g_mockFile.NewMockChainInterface(controller)
	mb := NewMockAbstractBlock(controller)
	mb.EXPECT().IsSpecial().Return(true)
	mb.EXPECT().Version().Return(uint64(100))
	mc.EXPECT().GetChainConfig().Return(chain_config.GetChainConfig()).AnyTimes()

	assert.Error(t, NewBftChainWriterWithoutVotes(&middleware.BftBlockContextWithoutVotes{
		BlockContext: middleware.BlockContext{Block: mb, Chain: mc},
	}, mc).SaveBlock())
}

func TestBftChainWriter_SaveBlock_For1327041(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	mc := g_mockFile.NewMockChainInterface(controller)
	//mb := NewMockAbstractBlock(controller)
	//mb.EXPECT().IsSpecial().Return(false)
	//mb.EXPECT().Version().Return(uint64(0))
	header1323008 := model.NewHeader(0, 1323008, common.HexToHash("0x00001730080caf977f01521360ccf597d2b9728a0c8e27a6d71b26f78814e1f0"), common.HexToHash("0xce091940ad1cbe6bf3556d1487799cee509ea098df562138a4369ba1b57c4706"), common.HexToDiff("0x1e17f011"), big.NewInt(1576820079415692271), common.HexToAddress("0x00004f011d62285f527ce47D189EBA821601dAf8A16B"), common.BlockNonceFromHex("0x000000000000000000000000000000000000000000000000000000000004ec6a"))

	header1327040 := model.NewHeader(0, 1327040, common.HexToHash("0x00001730080caf977f01521360ccf597d2b9728a0c8e27a6d71b26f78814e1f0"), common.HexToHash("0xce091940ad1cbe6bf3556d1487799cee509ea098df562138a4369ba1b57c4706"), common.HexToDiff("0x1e17f011"), big.NewInt(1576820079415692271), common.HexToAddress("0x00004f011d62285f527ce47D189EBA821601dAf8A16B"), common.BlockNonceFromHex("0x000000000000000000000000000000000000000000000000000000000004ec6a"))
	header1327040.TransactionRoot = common.HexToHash("0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")
	header1327040.StateRoot = common.HexToHash("0x1b33e1134caf4ce3a54e1706f8cec97f450aba432aa9d5a8179e4e4267d354cd")
	header1327040.InterlinkRoot = common.HexToHash("0x6874d5ede0aac4c07eff23f41add89d4e7926aa265027396aa750e2283988172")
	header1327040.RegisterRoot = common.HexToHash("0x652719fecee37fd80fb1900ba99e0aef818fe6215bca79aec0fad5482e725ef3")
	header1327040.ReceiptHash = common.HexToHash("0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")

	header1327040.GasLimit = uint64(3360000000)
	mc.EXPECT().GetChainConfig().Return(chain_config.GetChainConfig()).AnyTimes()
	mc.EXPECT().GetBlockByNumber(uint64(1323008)).Return(model.NewBlock(header1323008, nil, nil)).AnyTimes()
	block1327040 := model.NewBlock(header1327040, nil, nil)
	t.Log("block1327040 VerificationRoot", block1327040.VerificationRoot())
	header1327040.SetVerificationRoot(common.HexToHash("0x316e71157ad841bd0a48eb297db710b4ca692c5e3ba47b77d495a22d00610599"))
	t.Log("block1327040 VerificationRoot after", block1327040.VerificationRoot())

	mc.EXPECT().GetLatestNormalBlock().Return(model.NewBlock(header1323008, nil, nil)).AnyTimes()

	assert.Error(t, NewBftChainWriter(&middleware.BftBlockContext{
		BlockContext: middleware.BlockContext{Block: block1327040, Chain: mc},
	}, mc).SaveBlock())
}
