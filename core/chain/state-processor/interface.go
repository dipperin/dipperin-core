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

package state_processor

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/third-party/trie"
	"github.com/ethereum/go-ethereum/ethdb"
	"math/big"
	"reflect"
)

type StateTrie interface {
	TryGet(key []byte) ([]byte, error)
	TryUpdate(key, value []byte) error
	TryDelete(key []byte) error
	Commit(onleaf trie.LeafCallback) (common.Hash, error)
	Hash() common.Hash
	NodeIterator(startKey []byte) trie.NodeIterator
	GetKey([]byte) []byte // TODO(fjl): remove this when SecureTrie is removed
	Prove(key []byte, fromLevel uint, proofDb ethdb.Putter) error
}

type StateStorage interface {
	// OpenTrie opens the main account trie.
	OpenTrie(root common.Hash) (StateTrie, error)

	// OpenStorageTrie opens the storage trie of an account.
	OpenStorageTrie(addrHash, root common.Hash) (StateTrie, error)

	// CopyTrie returns an independent copy of the given trie.
	CopyTrie(StateTrie) StateTrie

	// TrieDB retrieves the low level trie database used for data storage.
	TrieDB() *trie.Database

	DiskDB() ethdb.Database
}

type AccountStateReader interface {
	GetNonce(addr common.Address) (uint64, error)
	GetBalance(addr common.Address) (*big.Int, error)
	GetContractRoot(addr common.Address) (common.Hash, error)
	GetLastElect(addr common.Address) (uint64, error)
	GetStake(addr common.Address) (*big.Int, error)
	GetCommitNum(addr common.Address) (uint64, error)
	GetVerifyNum(addr common.Address) (uint64, error)
	GetProduceNum(addr common.Address) (uint64, error)
}

//type AccountStateWriter interface {
//
//}

// only genesis use
type AccountStateProcessor interface {
	NewAccountState(addr common.Address) (err error)
	GetBalance(addr common.Address) (*big.Int, error)
	SetBalance(addr common.Address, amount *big.Int) error
	Finalise() (root common.Hash, err error)
	Commit() (root common.Hash, err error)

	PutContract(addr common.Address, v reflect.Value) error
	GetContract(addr common.Address, vType reflect.Type) (v reflect.Value, err error)
	ContractExist(addr common.Address) bool
}

// state_transaction use
// Message represents a message sent to a contract.
type Message interface {
	From() common.Address
	To() *common.Address
	GasPrice() *big.Int
	Gas() uint64
	SetGas(gas uint64)
	Value() *big.Int
	Nonce() uint64
	CheckNonce() bool
	Data() []byte
}
