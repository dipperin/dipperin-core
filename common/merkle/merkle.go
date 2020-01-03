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

package merkle

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
)

//	   Merkel tree vulnerability description, copy from BitCoin
/*     WARNING! If you're reading this because you're learning about cs_crypto
       and/or designing a new system that will use merkle trees, keep in mind
       that the following merkle tree algorithm has a serious flaw related to
       duplicate txids, resulting in a vulnerability (CVE-2012-2459).

       The reason is that if the number of hashes in the list at a given time
       is odd, the last one is duplicated before computing the next level (which
       is unusual in Merkle trees). This results in certain sequences of
       transactions leading to the same merkle root. For example, these two
       trees:

                    A               A
                  /  \            /   \
                B     C         B       C
               / \    |        / \     / \
              D   E   F       D   E   F   F
             / \ / \ / \     / \ / \ / \ / \
             1 2 3 4 5 6     1 2 3 4 5 6 5 6

       for transaction lists [1,2,3,4,5,6] and [1,2,3,4,5,6,5,6] (where 5 and
       6 are repeated) result in the same root hash A (because the hash of both
       of (F) and (F,F) is C).

       The vulnerability results from being able to send a block with such a
       transaction list, with the same merkle root, and the same block hash as
       the original without duplication, resulting in failed validation. If the
       receiving node proceeds to mark that block as permanently invalid
       however, it will fail to accept further unmodified (and thus potentially
       valid) versions of the same block. We defend against this by detecting
       the case where we would hash two identical hashes at the end of the list
       together, and treating that identically to the block having an invalid
       merkle root. Assuming no double-SHA256 collisions, this will detect all
       known ways of changing the transactions without affecting the merkle
       root.
*/

const (
	MaxRouteNumber = 32
)

// MerkleComputation Calculate RootHash and Hash path for specific leaf transactions based on the transaction hash list and leaf location.
func MerkleComputation(TdHash []common.Hash, BranchPos uint32, pMutated *bool) (pRoot common.Hash, pBranch []common.Hash) {

	// check the incoming parameters
	if pMutated == nil {
		//cslog.DLogger.Error().Msg("the pBranch or pMutated is nil")
		return
	}

	if len(TdHash) == 0 {
		//cslog.DLogger.Error().Msg("the tx number is 0")
		return
	} else if len(TdHash) > MaxRouteNumber {
		//cslog.DLogger.Error().Msg("the tx number exceed 2^32")
		return
	}

	if BranchPos >= uint32(len(TdHash)) {
		//cslog.DLogger.Error().Msg("the BranchPos exceed the transaction number")
		return
	}

	//First convert all leaf nodes into subtree hash values ​​and store them in the inner list, where the subscript is the layer number of the tree.
	inner := make([]common.Hash, MaxRouteNumber)
	var count uint32 = 0
	var mutated bool = false
	var MatchLevel int = -1

	// outer loop processing each leaf
	for count < uint32(len(TdHash)) {
		h := TdHash[count]
		Matchh := count == BranchPos

		//cslog.Printf("the Matchh is: %v",Matchh)
		count++

		var level int
		//cslog.Printf("the count is: %v",count)

		//When the leaf index +1 is odd, it is first buffered into the inner, but when it is even, the same layer's cache is calculated in the inner node.
		for level = 0; count&(uint32(0x01)<<uint32(level)) == 0x0; level++ {

			//record the accompanying path hash value
			if Matchh { //If the node to be recorded with the accompanying path is processed, when its index +1 is even, its accompanying node should be the left node, and this node is saved in the inner
				pBranch = append(pBranch, inner[level])
			} else if MatchLevel == level {
				//When Matchh is false, but MatchLevel == level, it means that the index of the node to record the accompanying path is +1, and the accompanying node is the current processing even node h.
				pBranch = append(pBranch, h)
				// When the record is completed, if the for loop is not finished, it means that the hash value of the next layer is still calculated. Therefore, Matchh needs to be set to true, so that the next accompanying path node hash can be recorded. At the same time, Matchh is set to true, so that in any case, the subsequent MatchLevel can record layer number of the next companion node to be saved.
				Matchh = true
			}

			//There are cases where the hash values ​​of the two nodes to be spliced ​​are the same. This case needs to be marked. a vulnerability belonging to the Merkel tree
			if inner[level] == h {
				mutated = true
			}
			// calculate the hash value of the next layer together with the node in the same layer cache in inner
			h = cs_crypto.Keccak256Hash(inner[level][:], h[:])
		}

		//Save the calculated hash value to inner when the leaf index +1 is odd
		inner[level] = h

		//Regardless of the index of the leaf node +1 is even or odd, the accompanying node recording process is processed on the even node.
		if Matchh {
			//When the index of the leaf node to be recorded +1 is an odd number, since it is processed at the next even node, Matchh is false when processing.
			//so use matchlevel to mark
			//when this loop exit, MatchLevel record the layer number of the next companion node to be saved
			MatchLevel = level
		}

	}
	/*
		cslog.DLogger.Debug().Msg("the inner data is: ")
		for TestNum,testdata := range inner{
			cslog.Printf("%v:%v",TestNum,testdata)
		}

		cslog.DLogger.Debug().Msg("the pBranch is: ")
		for TestNum1,testdata1 := range pBranch{
			cslog.Printf("%v:%v",TestNum1,testdata1)
		}
	*/
	//Get the final merkle root and path based on the obtained subtree hash list and branch
	var tmplevel int = 0
	//If the number of layer nodes is even, the uncalculated hash node value will not be retained in inner during the above calculation.
	for count&(0x01<<uint(tmplevel)) == 0x0 {
		tmplevel++
	}

	/*
		cslog.Printf("the count is: %v",count)
		cslog.Printf("the tmplevel is: %v",tmplevel)
	*/

	temph := inner[tmplevel]
	tempmatch := MatchLevel == tmplevel

	//When the sum of count and extended leaves is just 2 tmplevel power, the loop is jumped out.
	for count != (uint32(0x01) << uint(tmplevel)) {

		//When the processing layer number is the Branch node layer number to be recorded, the hash value of this node is the companion node to be recorded.
		if tempmatch {
			pBranch = append(pBranch, temph)
		}

		//For the node that needs to be copied, copy its own hash value to calculate the hash value of the next layer.
		temph = cs_crypto.Keccak256Hash(temph[:], temph[:])

		//If you need to add an extension node and do your own hash operation each layer, each time is equivalent to adding 2^level leaves.
		//Then the total number of leaves of expanded trees needs to increase accordingly.
		count += uint32(0x01) << uint(tmplevel)

		//The node copy and calculates the hash and processing the next layer of tree after finished.
		tmplevel++

		//If it is an even layer, you need to calculate the hash value saved in the inner to get the hash data of the next layer.
		for (count & (0x01 << uint(tmplevel))) == 0x0 {
			//save the remaining companion node values
			if tempmatch {
				pBranch = append(pBranch, inner[tmplevel])
			} else if MatchLevel == tmplevel {
				pBranch = append(pBranch, temph)
				tempmatch = true
			}

			temph = cs_crypto.Keccak256Hash(inner[tmplevel][:], temph[:])
			tmplevel++
		}
	}

	*pMutated = mutated
	pRoot = temph

	return pRoot, pBranch
}

// ComputeMerkleRoot calculate Merkelgen based on the transaction hash list
func ComputeMerkleRoot(TdHash []common.Hash, pMutated *bool) common.Hash {

	var pRoot common.Hash

	pRoot, _ = MerkleComputation(TdHash, 0, pMutated)

	return pRoot
}

// ComputeMerkleBranch get its merkle path based on the transaction hash list and its index
//func ComputeMerkleBranch(TdHash []common.Hash, index uint32) []common.Hash {
//
//	HashRoute := make([]common.Hash, MaxRouteNumber)
//
//	var pMutated bool
//	_, HashRoute = MerkleComputation(TdHash, index, &pMutated)
//
//	return HashRoute
//}

// ComputeMerkleRootFromBranch calculate its Merkel root value based on the transaction, its trading index in the block, and its merkle path
//func ComputeMerkleRootFromBranch(merkleroute []common.Hash, leaf common.Hash, index uint32) common.Hash {
//
//	TmpHash := leaf
//
//	for _, RouteHash := range merkleroute {
//
//		if index&0x01 == 1 { //when the index is odd the stitching order isRouteHash+TmpHash
//
//			TmpHash = cs_crypto.Keccak256Hash(RouteHash[:], TmpHash[:])
//
//		} else { //when the index is even the stitching order isTmpHash + RouteHash
//
//			TmpHash = cs_crypto.Keccak256Hash(TmpHash[:], RouteHash[:])
//		}
//
//		//the index value of the previous layer is index/2
//		index >>= 0x01
//	}
//
//	return TmpHash
//}
