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
	"github.com/dipperin/dipperin-core/common/util"
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
	return
	controller := gomock.NewController(t)
	defer controller.Finish()
	mc := g_mockFile.NewMockChainInterface(controller)
	//mb := NewMockAbstractBlock(controller)
	//mb.EXPECT().IsSpecial().Return(false)
	//mb.EXPECT().Version().Return(uint64(0))
	header1323008 := model.NewHeader(0, 1323008, common.HexToHash("0x00001730080caf977f01521360ccf597d2b9728a0c8e27a6d71b26f78814e1f0"), common.HexToHash("0xce091940ad1cbe6bf3556d1487799cee509ea098df562138a4369ba1b57c4706"), common.HexToDiff("0x1e17f011"), big.NewInt(1576820079415692271), common.HexToAddress("0x00004f011d62285f527ce47D189EBA821601dAf8A16B"), common.BlockNonceFromHex("0x000000000000000000000000000000000000000000000000000000000004ec6a"))

	seed1327040 := common.HexToHash("0xce091940ad1cbe6bf3556d1487799cee509ea098df562138a4369ba1b57c4706")

	header1327040 := model.NewHeader(0, 1327040, common.HexToHash("0x00001730080caf977f01521360ccf597d2b9728a0c8e27a6d71b26f78814e1f0"), seed1327040, common.HexToDiff("0x1e17f011"), big.NewInt(1576820079415692271), common.HexToAddress("0x00004f011d62285f527ce47D189EBA821601dAf8A16B"), common.BlockNonceFromHex("0x000000000000000000000000000000000000000000000000000000000004ec6a"))
	header1327040.TransactionRoot = common.HexToHash("0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")
	header1327040.StateRoot = common.HexToHash("0x1b33e1134caf4ce3a54e1706f8cec97f450aba432aa9d5a8179e4e4267d354cd")
	header1327040.InterlinkRoot = common.HexToHash("0x6874d5ede0aac4c07eff23f41add89d4e7926aa265027396aa750e2283988172")
	header1327040.RegisterRoot = common.HexToHash("0x652719fecee37fd80fb1900ba99e0aef818fe6215bca79aec0fad5482e725ef3")
	header1327040.ReceiptHash = common.HexToHash("0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")
	header1327040.MinerPubKey = []byte{
		4, 101, 24, 121, 44, 169, 195, 129, 28, 111, 196, 132, 76, 233, 192, 210, 55, 200, 121, 8, 92, 8, 69, 53, 254, 132, 252, 66, 137, 146, 186, 178, 236, 252, 248, 8, 116, 185, 206, 107, 218, 41, 153, 97, 94, 213, 172, 47, 0, 18, 122, 163, 2, 220, 149, 38, 185, 33, 63, 8, 29, 40, 205, 133, 152,
	}
	header1327040.Proof = []byte{
		42, 182, 51, 209, 226, 169, 67, 203, 252, 70, 87, 101, 49, 188, 159, 38, 241, 106, 158, 16, 35, 230, 212, 150, 140, 116, 206, 54, 137, 49, 132, 208, 129, 159, 149, 207, 182, 183, 10, 112, 236, 23, 66, 44, 163, 28, 158, 86, 158, 4, 187, 31, 47, 16, 81, 98, 150, 252, 150, 32, 114, 109, 87, 157, 4, 95, 158, 205, 136, 150, 201, 238, 239, 230, 178, 129, 223, 153, 222, 95, 42, 209, 113, 79, 197, 51, 226, 234, 151, 94, 76, 59, 88, 134, 151, 12, 18, 179, 56, 203, 208, 221, 23, 10, 96, 64, 233, 170, 64, 85, 92, 234, 14, 162, 47, 120, 19, 23, 228, 200, 151, 137, 6, 171, 154, 129, 214, 61, 163,
	}
	str1327040 := getVoteMsgJsonArr1327040()
	voteMsg1327040 := strsToVoteMsg(str1327040)

	header1327040.GasLimit = uint64(3360000000)
	mc.EXPECT().GetChainConfig().Return(chain_config.GetChainConfig()).AnyTimes()
	mc.EXPECT().GetBlockByNumber(uint64(1323008)).Return(model.NewBlock(header1323008, nil, nil)).Times(1)
	mc.EXPECT().GetBlockByNumber(uint64(1327039)).Return().Times(1)
	block1327040 := model.NewBlock(header1327040, nil, voteMsg1327040)

	t.Log("block1327040 VerificationRoot", block1327040.VerificationRoot())
	header1327040.SetVerificationRoot(common.HexToHash("0x316e71157ad841bd0a48eb297db710b4ca692c5e3ba47b77d495a22d00610599"))
	t.Log("block1327040 VerificationRoot after", block1327040.VerificationRoot())

	mc.EXPECT().GetLatestNormalBlock().Return(model.NewBlock(header1323008, nil, nil)).AnyTimes()

	assert.Error(t, NewBftChainWriter(&middleware.BftBlockContext{
		BlockContext: middleware.BlockContext{Block: block1327040, Chain: mc},
	}, mc).SaveBlock())
}

func getVoteMsgJsonArr1327040() []string {
	return []string{
		`{"Height":1327039,"Round":78,"BlockID":"0x00001730080caf977f01521360ccf597d2b9728a0c8e27a6d71b26f78814e1f0","VoteType":1,"Timestamp":"0001-01-01T00:00:00Z","Witness":{"address":"0x00005966ffecc91f26bf9b6bb2b7e09a0d00465e2940","sign":"pepeJXp/Sjy/qTz5xRUel/rcG4Z/q8wtwcWa6QD92mBULV0td7Lt1HJSoB2rL0ro1kLmZH8wJNoBCIfD1969nQA="}}`,
		`{"Height":1327039,"Round":78,"BlockID":"0x00001730080caf977f01521360ccf597d2b9728a0c8e27a6d71b26f78814e1f0","VoteType":1,"Timestamp":"0001-01-01T00:00:00Z","Witness":{"address":"0x0000ea4b978ad52d5cbefad6f4e3db1688ed6b659137","sign":"IYFm8enmKjH05nT3kOLoRp3MFND4nKtDt8sryY88DnBispYujyX5nZd2hmavyfGlLryMU0Ygm9uhueEwplolSAE="}}
`, `{"Height":1327039,"Round":78,"BlockID":"0x00001730080caf977f01521360ccf597d2b9728a0c8e27a6d71b26f78814e1f0","VoteType":1,"Timestamp":"0001-01-01T00:00:00Z","Witness":{"address":"0x00000c6b87dd03d80d229031dadc2cd10fd8a8c24133","sign":"U4WFF0T6a92N/Ce36ys0HB9+UApmRSwJPBfrwk6X1xsJjHLf/mhtq6V7zb97uJs3BwAwlFB2+FzFntE6X5WYSwE="}}
`, `{"Height":1327039,"Round":78,"BlockID":"0x00001730080caf977f01521360ccf597d2b9728a0c8e27a6d71b26f78814e1f0","VoteType":1,"Timestamp":"0001-01-01T00:00:00Z","Witness":{"address":"0x00004179d57e45cb3b54d6faef69e746bf240e287978","sign":"19aGgnREBgXFtk0eqXVTwE+HPSEucQ9aj6fSMCOQCzAvX69zqzdMqOh8JO8nyld3uR0nnx+V+WJE3uL2TJT9OQA="}}
`, `{"Height":1327039,"Round":78,"BlockID":"0x00001730080caf977f01521360ccf597d2b9728a0c8e27a6d71b26f78814e1f0","VoteType":1,"Timestamp":"0001-01-01T00:00:00Z","Witness":{"address":"0x0000daa28ec52c284ca84aacbd69039269d7a36247c2","sign":"VRIFdqCSvTNiG8s+OxdMsQxiezGjIa6gaVPwlYWJTnMQejJbDWyb7axdnuPEg4W9M1kC8XqohrnEmJnFurNyVwA="}}
`, `{"Height":1327039,"Round":78,"BlockID":"0x00001730080caf977f01521360ccf597d2b9728a0c8e27a6d71b26f78814e1f0","VoteType":1,"Timestamp":"0001-01-01T00:00:00Z","Witness":{"address":"0x0000d8fba7dd94654e3dfbfa78b76bb40db605b4e6ca","sign":"jryz06La4C/8X5RlbT04t8jDV4/dV6tnI/a+a7RGqThKF3CeiYozI12nUSq3P5rX09ywCb0tuBjIH9825dJtZgA="}}
`, `{"Height":1327039,"Round":78,"BlockID":"0x00001730080caf977f01521360ccf597d2b9728a0c8e27a6d71b26f78814e1f0","VoteType":1,"Timestamp":"0001-01-01T00:00:00Z","Witness":{"address":"0x00002b9af8390c3ca1da57054dc8947c9f92a1ca0c99","sign":"LLwVL+4QjHRT/DIuCD+8uPMyNnv2x4JUvlmJopdpRwMmaBRETP3Pq3jI2hWxpqqt790TTGYHHMsSyCTZoU4HTQE="}}
`, `{"Height":1327039,"Round":78,"BlockID":"0x00001730080caf977f01521360ccf597d2b9728a0c8e27a6d71b26f78814e1f0","VoteType":1,"Timestamp":"0001-01-01T00:00:00Z","Witness":{"address":"0x00000add04ac4d527de866cde4c93ae1662214617b12","sign":"ykn1xz0Oif/K/SmDhGCOFRSKIxcxIrSINt0Akxl0fd5z5Qszmy7g37X+2HC0c0gdAEiC/8s5IEacoCCx8HmVDAA="}}
`, `{"Height":1327039,"Round":78,"BlockID":"0x00001730080caf977f01521360ccf597d2b9728a0c8e27a6d71b26f78814e1f0","VoteType":1,"Timestamp":"0001-01-01T00:00:00Z","Witness":{"address":"0x0000f984742b330ec987c3df79c71ce1e729498cc613","sign":"AwnR9PQvAWmhgOL/toheht4FlGG9rUPknckCIweGCZ48R6ZhdSrUi9yZ1a3wUEc48K9G2XSB6mg9kUN6Dw+9ZQE="}}
`, `{"Height":1327039,"Round":78,"BlockID":"0x00001730080caf977f01521360ccf597d2b9728a0c8e27a6d71b26f78814e1f0","VoteType":1,"Timestamp":"0001-01-01T00:00:00Z","Witness":{"address":"0x0000fcd90394e9e902220bfd7a2d04431b3cfba2e2dd","sign":"OJpauCgdxO1Z8TxABkVeejBBbUFJk2sO9noAV+N91VgDRjsj2Jvke5jwi59FFQq0WTp2fCozWPxcTfT96hmhFgE="}}
`, `{"Height":1327039,"Round":78,"BlockID":"0x00001730080caf977f01521360ccf597d2b9728a0c8e27a6d71b26f78814e1f0","VoteType":1,"Timestamp":"0001-01-01T00:00:00Z","Witness":{"address":"0x0000d552c7ec7735668571a8ed92b57122f21fcea599","sign":"/CnCJjUHopSjR0IQX1OmBNKI2ym3WdXV9TWcphbWN0QzVbqoxbNZ2ZQKoDL7TJ2tNpzslBNk0WUpitErFY+OBAA="}}
`, `{"Height":1327039,"Round":78,"BlockID":"0x00001730080caf977f01521360ccf597d2b9728a0c8e27a6d71b26f78814e1f0","VoteType":1,"Timestamp":"0001-01-01T00:00:00Z","Witness":{"address":"0x0000b57898eb80649b2f9993d8a3941ed195961368e9","sign":"R6W3K2mAVttFpTZZX69QbdQBzc5nymFxgpQVzww74oRZASSQwX7elin3g9HAUv5+LtWI6NtmflAo/5debDvRZwA="}}
`, `{"Height":1327039,"Round":78,"BlockID":"0x00001730080caf977f01521360ccf597d2b9728a0c8e27a6d71b26f78814e1f0","VoteType":1,"Timestamp":"0001-01-01T00:00:00Z","Witness":{"address":"0x0000918c773880b462929ace4f975ccfed9be2d8efc9","sign":"X+SkVJcT+HYhabH7vIf5nlxYgtpA9eRgxzb6cYU4/YgXJea8Z8se0dnwDAQRcHes0fMm10sa/ctERbFQbG66QwA="}}
`, `{"Height":1327039,"Round":78,"BlockID":"0x00001730080caf977f01521360ccf597d2b9728a0c8e27a6d71b26f78814e1f0","VoteType":1,"Timestamp":"0001-01-01T00:00:00Z","Witness":{"address":"0x0000f55be671e8ff2184b0c7181ae0e4cd92429c034c","sign":"8tcb5uc1JJV8V9MtxnoHUEVmHc7TZb2jVC4fdTRdSqcPr8fUpQqtBK1znR5x5W+pOGA2U8cosrVEW1mL2t8UhwA="}}
`, `{"Height":1327039,"Round":78,"BlockID":"0x00001730080caf977f01521360ccf597d2b9728a0c8e27a6d71b26f78814e1f0","VoteType":1,"Timestamp":"0001-01-01T00:00:00Z","Witness":{"address":"0x00005eccf0aaa6e8f451078448a182970e80cbdd253b","sign":"WXgu7x0avREoGHP9Uh88WQeViO7zdFfD51/7hfjn/clGyox82mvo/ILumSAmTCSK/IGPoB9RZZ4lebRB6fzY1wE="}}`,
	}
}

func strsToVoteMsg(strs []string) []model.AbstractVerification {
	var voteMsgs []model.AbstractVerification
	for _, s := range strs {
		var ver model.VoteMsg
		util.ParseJson(s, &ver)
		voteMsgs = append(voteMsgs, ver)
	}
	return voteMsgs
}
