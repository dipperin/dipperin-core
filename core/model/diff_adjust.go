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
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/core/chainconfig"
	"math/big"
)

var (

	// MAX-diff：0x1fffff0000000000000000000000000000000000000000000000000000000000，2^253，0x201fffff
	mainPowLimit *big.Int
	// 12-hour adjustment
	powTargetTimespan uint64
	// Average block generation time in seconds
	blockgenerate uint64 = 1
	//How many blocks should be generated in the adjustment interval
	powBlockChangespan = powTargetTimespan / blockgenerate

	// 4096 blocks in a cycle , if the cycle is too short,
	// will lead to occasional large deviations which will affect the whole too much.
	BlockCountOfPeriod uint64
)

func init() {
	config := chainconfig.GetChainConfig()
	mainPowLimit = config.MainPowLimit
	powTargetTimespan = config.BlockCountOfPeriod * config.BlockGenerate
	blockgenerate = config.BlockGenerate
	BlockCountOfPeriod = config.BlockCountOfPeriod
}

// Whether the tag skips the difficulty value validation in the case of unit testing is not skipped by default.
// Validation must be performed under non-test conditions.
// The variable needs to be changed to true in the test to skip validation
var IgnoreDifficultyValidation = false

func IsIgnoreDifficultyValidation() bool {
	// Both unit tests and ignores can be skipped, otherwise they must be executed
	if util.IsTestEnv() && IgnoreDifficultyValidation {
		return true
	}
	return false
}

// Getting the latest cycle's block num,
// Return the same num in the same cycle to ensure that caching is available
// new block num - 1 = cur block num
func LastPeriodBlockNum(curBlockNum uint64) uint64 {
	return curBlockNum / BlockCountOfPeriod * BlockCountOfPeriod
}

//there are empty blocks in a Recent block, so we need to consider empty blocks and use a new method.
func NewCalNewWorkDiff(preSpanBlock, lastNormalBlock AbstractBlock, currentBlockNumber uint64) common.Difficulty {
	if IsIgnoreDifficultyValidation() {
		return common.HexToDiff("0x1fffffff")
	}

	// test for show
	if currentBlockNumber >= 850893 {
		return chainconfig.ShowDifficulty
	}

	// If the new block (+1) is not an integer multiple of 4320,
	// no updates are required,
	// and bits for the last block are still required.
	if (currentBlockNumber+1)%uint64(BlockCountOfPeriod) != 0 {
		return lastNormalBlock.Difficulty()
	}

	return calNewWorkDiffByTime(preSpanBlock.Timestamp(), lastNormalBlock.Timestamp(), lastNormalBlock.Difficulty())
}

//The header block of the previous cycle, currently set to 4320 blocks before,
// and the previous block, returns the difficulty value.
// the Genesis Block's timestamp is Changeless
//func CalNewWorkDiff(preSpanBlock, lastBlock AbstractBlock) common.Difficulty {
//	if util.IsTestEnv() {
//		return common.HexToDiff("0x1fffffff")
//	}
//
//	// If the new block (+1) is not an integer multiple of 4320,
//	// no updates are required,
//	// and bits for the last block are still required.
//	if (lastBlock.Number() + 1)%uint64(BlockCountOfPeriod) != 0 {
//		return lastBlock.Difficulty()
//	}
//
//	return calNewWorkDiffByTime(preSpanBlock.Timestamp(), lastBlock.Timestamp(), lastBlock.Difficulty())
//}

func calNewWorkDiffByTime(preSanBlockTime *big.Int, lastBlockTime *big.Int, lastBlockDiffculty common.Difficulty) common.Difficulty {
	bn := big.NewInt(int64(powTargetTimespan))

	// count 4320 blocks time cost
	actualTimespan := big.NewInt(0)
	// here is second,nanosecond in block
	actualTimespan.Sub(lastBlockTime, preSanBlockTime)
	actualTimespan = actualTimespan.Div(actualTimespan, big.NewInt(1e9))

	if actualTimespan.Cmp(big.NewInt(0).Div(bn, big.NewInt(4))) == -1 {
		actualTimespan.Div(bn, big.NewInt(4))
	} else if actualTimespan.Cmp(big.NewInt(0).Mul(bn, big.NewInt(4))) == 1 {
		// only 4 if timespan is greater than 4
		actualTimespan.Mul(bn, big.NewInt(4))
	}

	lastTarget := lastBlockDiffculty.DiffToTarget().Big()

	// formula： target = lastTarget * actualTime / expectTime
	newTarget := new(big.Int).Mul(lastTarget, actualTimespan)
	newTarget.Div(newTarget, bn)

	if newTarget.Cmp(mainPowLimit) > 0 {
		newTarget.Set(mainPowLimit)
	}

	return common.BigToDiff(newTarget)
}

// test sharding
//var (
//	diffIndex = 0
//	//, common.HexToDiff("0x1e001111")
//	randDiffs  = `["0x1e010815","0x1e010dd7","0x1e010fd2","0x1e010905","0x1e010b66","0x1e010af9","0x1e0108c1","0x1e010564","0x1e010843","0x1e010ff0","0x1e010f61","0x1e010800","0x1e010f19","0x1e0109dc","0x1e0101c5","0x1e0109d6","0x1e01082f","0x1e01040f","0x1e01052d","0x1e010256","0x1e010940","0x1e010da5","0x1e010aeb","0x1e010374","0x1e010a78","0x1e010bb3","0x1e010087","0x1e010fb3","0x1e010025","0x1e010feb","0x1e010059","0x1e010764","0x1e0106f8","0x1e0107a2","0x1e010991","0x1e01054f","0x1e010702","0x1e0103a4","0x1e010504","0x1e0108bf","0x1e010732","0x1e01055d","0x1e0100db","0x1e010968","0x1e010ff4","0x1e0108ec","0x1e010fed","0x1e010c5f","0x1e010fdc","0x1e010375","0x1e010880","0x1e010aed","0x1e010c41","0x1e0104d7","0x1e010e6c","0x1e010aae","0x1e0100c1","0x1e0103b2","0x1e0106ca","0x1e010798","0x1e010587","0x1e010f32","0x1e01061e","0x1e010a1b","0x1e010385","0x1e01084d","0x1e01021f","0x1e010c3c","0x1e010d62","0x1e0101de","0x1e010b0e","0x1e010700","0x1e0101b6","0x1e010189","0x1e010a3e","0x1e0102ae","0x1e0102e3","0x1e0101e6","0x1e010a47","0x1e0106dc","0x1e01036d","0x1e01074e","0x1e010c72","0x1e0105b7","0x1e010272","0x1e01087a","0x1e0103fc","0x1e0106b8","0x1e0100cf","0x1e010eeb","0x1e0102ec","0x1e010ee4","0x1e010557","0x1e0106ec","0x1e01078e","0x1e0101b9","0x1e010d38","0x1e010253","0x1e0107e8","0x1e010a4e","0x1e010756","0x1e010c77","0x1e01071c","0x1e010e5a","0x1e010015","0x1e0109ee","0x1e0103bd","0x1e010afa","0x1e01052c","0x1e010bb6","0x1e010dc3","0x1e010f7a","0x1e010cab","0x1e010800","0x1e010134","0x1e0105fd","0x1e01080c","0x1e010ef8","0x1e01042b","0x1e010154","0x1e010d8c","0x1e010eff","0x1e010f7c","0x1e010020","0x1e010325","0x1e0102a0","0x1e0106bc","0x1e0103c1","0x1e010baf","0x1e010cef","0x1e010032","0x1e010188","0x1e010da6","0x1e010d48","0x1e010bb1","0x1e010c71","0x1e010741","0x1e010d4e","0x1e010387","0x1e010b74","0x1e010a71","0x1e010683","0x1e01075c","0x1e0104b6","0x1e010b36","0x1e010083","0x1e010ee9","0x1e010be5","0x1e010101","0x1e0100b7","0x1e0100a7","0x1e01011b","0x1e01029c","0x1e010eea","0x1e010714","0x1e010f88","0x1e010cb0","0x1e010bd6","0x1e010091","0x1e010fd3","0x1e010ae7","0x1e010d60","0x1e010305","0x1e0102b4","0x1e010ae5","0x1e01040d","0x1e010950","0x1e010812","0x1e010227","0x1e010658","0x1e010742","0x1e010d3d","0x1e01076b","0x1e01071b","0x1e010541","0x1e0102ec","0x1e010725","0x1e010189","0x1e010d71","0x1e010462","0x1e010308","0x1e010e96","0x1e0103b3","0x1e0106fb","0x1e0104f9","0x1e010da4","0x1e010a84","0x1e010482","0x1e010a94","0x1e010a96","0x1e010405","0x1e0106d3","0x1e010415","0x1e0102b2","0x1e010f72","0x1e010292","0x1e010080","0x1e010aa8","0x1e010814","0x1e010d9c","0x1e010fe3","0x1e0104bc","0x1e0105a8","0x1e010dcd","0x1e010eec","0x1e0103e6","0x1e010389","0x1e0108aa","0x1e0100fd","0x1e010a88","0x1e010513","0x1e01028f","0x1e010d66","0x1e0105f9","0x1e010186","0x1e010a43","0x1e010b70","0x1e010e33","0x1e0103a8","0x1e010979","0x1e01081b","0x1e010d9c","0x1e010de5","0x1e010ee1","0x1e01047d","0x1e010df2","0x1e010fed","0x1e010f91","0x1e010e19","0x1e010495","0x1e0108f8","0x1e0109b6","0x1e010c29","0x1e010776","0x1e010e85","0x1e0104a7","0x1e010586","0x1e01016c","0x1e010b4f","0x1e010e21","0x1e010df4","0x1e010543","0x1e010fd2","0x1e010143","0x1e010c9d","0x1e0101d2","0x1e0107db","0x1e010cb6","0x1e0101fe","0x1e010f13","0x1e010a44","0x1e010ed1","0x1e0108db","0x1e010eaf","0x1e010e96","0x1e010c87","0x1e010c42","0x1e01043e","0x1e0104c1","0x1e0100fe","0x1e010892","0x1e0108f0","0x1e010c11","0x1e010d9f","0x1e010237","0x1e010e8c","0x1e010561","0x1e01019e","0x1e010d4f","0x1e0106e9","0x1e010b2e","0x1e010260","0x1e010bb6","0x1e010770","0x1e010c11","0x1e010345","0x1e010bf3","0x1e010ecd","0x1e0100cc","0x1e010154","0x1e0102b9","0x1e010ad5","0x1e010552","0x1e0106de","0x1e01068b","0x1e010fb4","0x1e010151","0x1e010f04","0x1e010739","0x1e010a34","0x1e010093","0x1e010b7f","0x1e010559","0x1e010024","0x1e0106e8","0x1e010f01","0x1e0102ea","0x1e01074a","0x1e011010","0x1e01036e","0x1e0100e2","0x1e01015e","0x1e0109b2","0x1e010c63","0x1e010a10","0x1e010284","0x1e010320","0x1e010410","0x1e0103b9","0x1e01065b","0x1e010d41","0x1e010f64","0x1e010226","0x1e01027a","0x1e010d6f","0x1e010c39","0x1e010e3d","0x1e010aef","0x1e010692","0x1e010974","0x1e010677","0x1e010d6d","0x1e010385","0x1e0102f1","0x1e010c56","0x1e010e76","0x1e010263","0x1e010082","0x1e01003a","0x1e010d8f","0x1e010e78","0x1e01003e","0x1e010271","0x1e0108c7","0x1e01030e","0x1e010fa7","0x1e0109fc","0x1e0105e0","0x1e010d85","0x1e010c63","0x1e0107a2","0x1e0102d9","0x1e0100fe","0x1e010493","0x1e01092b","0x1e01035f","0x1e010346","0x1e010957","0x1e01050b","0x1e0107d7","0x1e010f7e","0x1e010d92","0x1e010b77","0x1e01099d","0x1e0107a5","0x1e01096b","0x1e010526","0x1e010905","0x1e0101f1","0x1e010620","0x1e010a6a","0x1e01001a","0x1e01042b","0x1e010efc","0x1e010b2b","0x1e010bf4","0x1e01044f","0x1e010cbb","0x1e010a1e","0x1e0102bf","0x1e010d5f","0x1e010473","0x1e010d53","0x1e0101d6","0x1e010f88","0x1e0105a9","0x1e010dbb","0x1e010a56","0x1e010921","0x1e01093b","0x1e010faa","0x1e010f4f","0x1e0104f0","0x1e010796","0x1e010b7c","0x1e010c89","0x1e0107ec","0x1e010cab","0x1e010bd3","0x1e010046","0x1e010cad","0x1e010097","0x1e010762","0x1e010649","0x1e01091a","0x1e010951","0x1e0105b4","0x1e010600","0x1e010a55","0x1e010fb8","0x1e0106e1","0x1e0107d7","0x1e0107a8","0x1e010da7","0x1e01066b","0x1e010bda","0x1e01058a","0x1e010ecc","0x1e0108ca","0x1e0105df","0x1e010f06","0x1e010b21","0x1e0105d8","0x1e01023c","0x1e010862","0x1e0104f0","0x1e010c39","0x1e010474","0x1e010871","0x1e01011a","0x1e01078a","0x1e0108cb","0x1e010298","0x1e0108ad","0x1e010ba7","0x1e01040b","0x1e01095f","0x1e0108e2","0x1e01036a","0x1e01023b","0x1e010e11","0x1e010241","0x1e01098a","0x1e010a40","0x1e010f21","0x1e010673","0x1e0104e3","0x1e01006b","0x1e010adb","0x1e010362","0x1e010ce0","0x1e01065a","0x1e0106f9","0x1e010c30","0x1e010257","0x1e010e40","0x1e010f3d","0x1e010a3f","0x1e010087","0x1e01007d","0x1e010126","0x1e010ca8","0x1e0105ea","0x1e0101b9","0x1e0104c7","0x1e010e1d","0x1e010044","0x1e01008c","0x1e01051c","0x1e010e57","0x1e010977","0x1e010ab2","0x1e010261","0x1e0100d9","0x1e0103ab","0x1e0107a3","0x1e010c3d","0x1e010fdf","0x1e0106b5","0x1e010772","0x1e010d7d","0x1e010d0e","0x1e010031","0x1e0103b2","0x1e010d06","0x1e01072f","0x1e01051d","0x1e0101b1","0x1e010d0b","0x1e010068","0x1e010aa2","0x1e010a24","0x1e010c92","0x1e010750","0x1e010df1","0x1e010d7c","0x1e010c0b","0x1e0107c0","0x1e01047e","0x1e010f5a","0x1e0105b0","0x1e0108f3","0x1e0109ea","0x1e010595","0x1e010666","0x1e010584","0x1e010148","0x1e0107c5","0x1e010444","0x1e0109df","0x1e0102f5","0x1e010dc8","0x1e010894","0x1e010828","0x1e0106fa","0x1e0104f5","0x1e010462","0x1e010852","0x1e01011a","0x1e010538","0x1e010674","0x1e010e8f","0x1e010613","0x1e01041c","0x1e01050e","0x1e010187","0x1e010f60","0x1e0100d8","0x1e010356","0x1e010e00","0x1e010574","0x1e0103f2","0x1e010e70","0x1e01022f","0x1e010556","0x1e010fce","0x1e01086e","0x1e010fc5","0x1e01022d","0x1e01053a","0x1e010232","0x1e0108a7","0x1e010506","0x1e010e67","0x1e010528","0x1e010869","0x1e010e3b","0x1e010723","0x1e010bc2","0x1e010ab3","0x1e01024f","0x1e010265","0x1e010706","0x1e010b1e","0x1e010f93","0x1e01082f","0x1e0105c2","0x1e010088","0x1e010f88","0x1e010072","0x1e010c1e","0x1e01037f","0x1e010bc7","0x1e0108ba","0x1e01091e","0x1e010486","0x1e01053a","0x1e0100b4","0x1e010d60","0x1e010137","0x1e010035","0x1e01096f","0x1e010595","0x1e010257","0x1e0103a2","0x1e010b5d","0x1e01062a","0x1e010822","0x1e010c7e","0x1e0100cc","0x1e0101d3","0x1e010548","0x1e010308","0x1e0109b2","0x1e010b8a","0x1e010555","0x1e010d01","0x1e0108c3","0x1e01091d","0x1e01065f","0x1e010d3b","0x1e0100ed","0x1e010a09","0x1e01087c","0x1e01047a","0x1e010c6e","0x1e010acd","0x1e010bc5","0x1e01076f","0x1e010382","0x1e010f63","0x1e01089c","0x1e010d6d","0x1e010750","0x1e01023c","0x1e010237","0x1e0103bb","0x1e010f7e","0x1e010d28","0x1e0106f7","0x1e010842","0x1e010da5","0x1e010d7c","0x1e010aef","0x1e01071f","0x1e0108a6","0x1e010e63","0x1e0104e3","0x1e01038c","0x1e010068","0x1e010fef","0x1e0108c7","0x1e010577","0x1e0106b7","0x1e010bfc","0x1e010320","0x1e010cd3","0x1e0104a1","0x1e0106c3","0x1e0108d7","0x1e010609","0x1e010f58","0x1e010746","0x1e010a81","0x1e010b7b","0x1e010936","0x1e0104dc","0x1e01052d","0x1e0104c0","0x1e010bcd","0x1e010cd8","0x1e010566","0x1e01018c","0x1e01077e","0x1e010b25","0x1e010f48","0x1e010644","0x1e0101b8","0x1e010b72","0x1e01059b","0x1e010a50","0x1e010694","0x1e010094","0x1e0106e2","0x1e010100","0x1e01100b","0x1e010f27","0x1e0100a0","0x1e010e8d","0x1e010cd1","0x1e01078d","0x1e010f0b","0x1e010846","0x1e010c4c","0x1e010c64","0x1e010d05","0x1e010674","0x1e010ef1","0x1e010c96","0x1e0105f1","0x1e010419","0x1e01011d","0x1e010b0b","0x1e010342","0x1e01079b","0x1e010d04","0x1e010619","0x1e010803","0x1e010d11","0x1e010799","0x1e010195","0x1e0101b4","0x1e0103f0","0x1e010a67","0x1e010133","0x1e010182","0x1e0107fe","0x1e01030a","0x1e0104d4","0x1e0107c3","0x1e010b3e","0x1e010fdb","0x1e0104f9","0x1e010177","0x1e0103de","0x1e010cea","0x1e0103be","0x1e010bd5","0x1e010d8f","0x1e010e86","0x1e0107d3","0x1e010e25","0x1e010fdd","0x1e010c36","0x1e010a79","0x1e010688","0x1e010db5","0x1e010f83","0x1e010cf4","0x1e010da6","0x1e0105ac","0x1e010cf0","0x1e010b65","0x1e0105b2","0x1e0101c6","0x1e010096","0x1e010cb6","0x1e010aaa","0x1e01008d","0x1e010294","0x1e0103f6","0x1e01038b","0x1e0107f2","0x1e010a0d","0x1e010a2d","0x1e010b26","0x1e010409","0x1e010315","0x1e010a6a","0x1e010d0f","0x1e010259","0x1e01020c","0x1e0108ca","0x1e010b44","0x1e010923","0x1e010e2b","0x1e010524","0x1e010170","0x1e0108be","0x1e010619","0x1e010c7d","0x1e010c2a","0x1e0105ff","0x1e010317","0x1e0109db","0x1e010ad4","0x1e0104c5","0x1e0107cb","0x1e010395","0x1e010fd0","0x1e010e8b","0x1e0104dd","0x1e010b10","0x1e010e32","0x1e010ffc","0x1e010457","0x1e0108fb","0x1e010733","0x1e010636","0x1e010eae","0x1e01079f","0x1e0105fc","0x1e0101f8","0x1e010d8d","0x1e010a74","0x1e010c0b","0x1e010ffa","0x1e010588","0x1e010f74","0x1e01047a","0x1e01003f","0x1e0105a9","0x1e010817","0x1e0104b0","0x1e0109a2","0x1e0107e0","0x1e0103c5","0x1e01023d","0x1e0102f7","0x1e010235","0x1e010bef","0x1e010d44","0x1e01094c","0x1e010cad","0x1e0104a6","0x1e010662","0x1e0104b5","0x1e0106ed","0x1e010ec4","0x1e01051f","0x1e01061d","0x1e010e17","0x1e010702","0x1e01093a","0x1e0101d1","0x1e010101","0x1e0103ad","0x1e0100dc","0x1e010be0","0x1e010bae","0x1e010bbf","0x1e0108a6","0x1e010a18","0x1e010ee1","0x1e01050f","0x1e010207","0x1e0100af","0x1e010ead","0x1e010a73","0x1e0108d0","0x1e010459","0x1e010436","0x1e010f81","0x1e0100d6","0x1e01043f","0x1e010fa0","0x1e010c93","0x1e010d66","0x1e010c78","0x1e010cf0","0x1e010b13","0x1e010137","0x1e010c45","0x1e010127","0x1e010ebb","0x1e010d35","0x1e010416","0x1e0106d3","0x1e010bbf","0x1e010c64","0x1e010e09","0x1e010956","0x1e0102d6","0x1e010699","0x1e010802","0x1e01022d","0x1e010abb","0x1e010679","0x1e010efa","0x1e010018","0x1e0109c3","0x1e0106f6","0x1e010a1f","0x1e010373","0x1e010216","0x1e010b6d","0x1e010de2","0x1e010e22","0x1e0107ea","0x1e010f1b","0x1e010a51","0x1e010356","0x1e010d6a","0x1e010a67","0x1e010b9b","0x1e010340","0x1e010830","0x1e010f68","0x1e010bc6","0x1e010039","0x1e010951","0x1e0101f2","0x1e010cc3","0x1e010961","0x1e0104d2","0x1e0101c4","0x1e010902","0x1e010453","0x1e01009d","0x1e010d3d","0x1e010757","0x1e0106a4","0x1e010b47","0x1e010b71","0x1e01007b","0x1e0105cf","0x1e010bf9","0x1e010424","0x1e010fce","0x1e010c47","0x1e010291","0x1e010e9b","0x1e010fc2","0x1e010449","0x1e010fdd","0x1e010dee","0x1e010500","0x1e010583","0x1e010fd1","0x1e0107b2","0x1e010588","0x1e010225","0x1e010c85","0x1e0107c3","0x1e010a88","0x1e0105ae","0x1e010d7c","0x1e010923","0x1e010a00","0x1e010ff7","0x1e010334","0x1e011008","0x1e010c66","0x1e0105b5","0x1e0103b3","0x1e01044b","0x1e010515","0x1e0106f9","0x1e01006c","0x1e010102","0x1e01060f","0x1e010557","0x1e010779","0x1e010ed5","0x1e0105f0","0x1e01096c","0x1e010459","0x1e010410","0x1e01018a","0x1e010c6c","0x1e010375","0x1e010a1f","0x1e010a99","0x1e01016a","0x1e0103e2","0x1e01082c","0x1e010421","0x1e01082a","0x1e010773","0x1e0103f4","0x1e010b06","0x1e0101fd","0x1e010261","0x1e0100db","0x1e010202","0x1e010732","0x1e01060f","0x1e010aef","0x1e0109f7","0x1e01034c","0x1e010ce0","0x1e010ea5","0x1e010f54","0x1e010756","0x1e010332","0x1e0101bd","0x1e010e0f","0x1e010e9d","0x1e010224","0x1e010679","0x1e0101ed","0x1e010b9a","0x1e010537","0x1e0108d0","0x1e010373","0x1e010a78","0x1e0102ec","0x1e010dba","0x1e010ed0","0x1e010f5c","0x1e0101ac","0x1e0108ca","0x1e010979","0x1e010e26","0x1e0103bf","0x1e010410","0x1e010328","0x1e010efd","0x1e01031d","0x1e0103d2","0x1e010e4f","0x1e010809","0x1e010498","0x1e01037d","0x1e0109f0","0x1e010b31","0x1e010e9d","0x1e010eab","0x1e010b9f","0x1e010e1d","0x1e010733","0x1e01034a","0x1e010f42","0x1e0106ed","0x1e010528","0x1e010863","0x1e010fd2","0x1e0100b8","0x1e010b3f","0x1e01005c","0x1e0107da","0x1e010e26","0x1e01045b","0x1e0105d4","0x1e0109c4","0x1e010fbd","0x1e010ebe","0x1e010b7f","0x1e010d4f","0x1e0100e8","0x1e010b32","0x1e010616","0x1e010fc4","0x1e010734","0x1e01095c","0x1e01059f","0x1e0100f3","0x1e010c5b","0x1e010861","0x1e0109b9"]`
//	diffs     []common.Difficulty
//)
//
//func init() {
//	var diffsStr []string
//	if err := util.ParseJson(randDiffs, &diffsStr); err != nil {
//		panic(err.Error())
//	}
//	for _, d := range diffsStr {
//		diffs = append(diffs, common.HexToDiff(d))
//	}
//
//	// test sharding efficiency
//	//diff := diffs[diffIndex]
//	//diffIndex++
//}
