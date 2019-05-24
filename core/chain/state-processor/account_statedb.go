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
	"github.com/dipperin/dipperin-core/common/util/json-kv"
	"github.com/dipperin/dipperin-core/core/contract"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/dipperin/dipperin-core/third-party/log/mpt_log"
	"github.com/dipperin/dipperin-core/third-party/trie"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
	"errors"
	"sync"
	"sort"
	"fmt"
	"strings"
	"reflect"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/common/g-error"
)

type account struct {
	Nonce   uint64
	Balance *big.Int
	// Verifier stake, commit number, produce num, verifier num
	Stake       *big.Int
	CommitNum   uint64
	VerifyNum   uint64
	LastElect   uint64
	Performance uint64

	HashLock common.Hash `rlp:"nil"`
	TimeLock *big.Int
	// not need, merkle root of the contract storage trie
	//ContractRoot common.Hash `rlp:"nil"`
	// merkle root of the triple-layered smart contraction data storage trie
	DataRoot common.Hash `rlp:"nil"`
}

var (
	SenderOrReceiverIsEmptyErr = errors.New("sender or receiver is empty")
)

const (
	nonceKeySuffix     = "_nonce"
	balanceKeySuffix   = "_balance"
	hashLockKeySuffix  = "_hashLock"
	timeLockKeySuffix  = "_timeLock"
	contractRootSuffix = "_contract_root"
	dataRootSuffix     = "_data_root"
	stakeKeySuffix     = "_stake"
	commitNumKeySuffix = "_commit_num"
	verifyNumKeySuffix = "_verify_num"
	lastElectKeySuffix = "_last_elect"
	performanceSuffix  = "_performance"
)

func GetContractFieldKey(address common.Address, key string) []byte {
	return append(address[:], []byte(key)...)
}

// get the real key without hash and address
func GetContractAddrAndKey(key []byte) (common.Address, []byte) {
	//the key is larger than addr because there is one character at least
	if len(key) > common.AddressLength {
		return common.BytesToAddress(key[:common.AddressLength]), key[common.AddressLength:]
	}
	return common.Address{}, nil
}

func GetContractRootKey(address common.Address) []byte {
	return append(address[:], []byte(contractRootSuffix)...)
}

func GetNonceKey(address common.Address) []byte {
	return append(address[:], []byte(nonceKeySuffix)...)
}

func GetBalanceKey(address common.Address) []byte {
	return append(address[:], []byte(balanceKeySuffix)...)
}

func GetHashLockKey(address common.Address) []byte {
	return append(address[:], []byte(hashLockKeySuffix)...)
}

func GetTimeLockKey(address common.Address) []byte {
	return append(address[:], []byte(timeLockKeySuffix)...)
}

func GetDataRootKey(address common.Address) []byte {
	return append(address[:], []byte(dataRootSuffix)...)
}

func GetStakeKey(address common.Address) []byte {
	return append(address[:], []byte(stakeKeySuffix)...)
}

func GetCommitNumKey(address common.Address) []byte {
	return append(address[:], []byte(commitNumKeySuffix)...)
}

func GetVerifyNumKey(address common.Address) []byte {
	return append(address[:], []byte(verifyNumKeySuffix)...)
}

func GetLastElectKey(address common.Address) []byte {
	return append(address[:], []byte(lastElectKeySuffix)...)
}

func GetPerformanceKey(address common.Address) []byte {
	return append(address[:], []byte(performanceSuffix)...)
}

func (a *account) getNonce() uint64 {
	return a.Nonce
}

func (a *account) setNonce(n uint64) {
	a.Nonce = n
}

//todo later use  rlp.EncodeToBytes(bytes.TrimLeft(value[:], "\x00")) method to save storage, decode method need use split
func (a *account) NonceBytes() []byte {
	v, _ := rlp.EncodeToBytes(a.Nonce)
	return v
}

func (a *account) BalanceBytes() []byte {
	v, _ := rlp.EncodeToBytes(a.Balance)
	return v
}

func (a *account) CommitNumBytes() []byte {
	v, _ := rlp.EncodeToBytes(a.CommitNum)
	return v
}

func (a *account) VerifyNumBytes() []byte {
	v, _ := rlp.EncodeToBytes(a.VerifyNum)
	return v
}

func (a *account) PerformanceBytes() []byte {
	v, _ := rlp.EncodeToBytes(a.Performance)
	return v
}

func (a *account) StakeBytes() []byte {
	v, _ := rlp.EncodeToBytes(a.Stake)
	return v
}

func (a *account) LastElectBytes() []byte {
	v, _ := rlp.EncodeToBytes(a.LastElect)
	return v
}

func (a *account) HashLockBytes() []byte {
	v, _ := rlp.EncodeToBytes(a.HashLock)
	return v
}

func (a *account) TimeLockBytes() []byte {
	v, _ := rlp.EncodeToBytes(a.TimeLock)
	return v
}

/*func (a *account) ContractRootBytes() []byte {
	v, _ := rlp.EncodeToBytes(a.ContractRoot)
	return v
}*/

func (a *account) DataRootBytes() []byte {
	v, _ := rlp.EncodeToBytes(a.DataRoot)
	return v
}

type revision struct {
	id          int
	changeIndex int
}

//call the Process function of this struct to get the state root
type AccountStateDB struct {
	preStateRoot common.Hash

	blockStateTrie StateTrie
	storage        StateStorage

	//each AccountStateDB own individual contract storage. new it when used
	contractTrieCache     StateStorage
	contractData          map[common.Address]reflect.Value
	finalisedContractRoot map[common.Address]common.Hash
	alreadyFinalised      bool

	stateChangeList *StateChangeList
	validRevisions  []revision
	nextRevisionId  int

	lock sync.Mutex
}

func (state *AccountStateDB) PreStateRoot() common.Hash {
	return state.preStateRoot
}

func (state *AccountStateDB) getContractTrie(addr common.Address) (StateTrie, error) {

	//notice: can't get the trie if the contract root had been changed but not commit
	cRoot, err := state.blockStateTrie.TryGet(GetContractRootKey(addr))
	mpt_log.Debug("get address contract root", "addr", addr.Hex(), "root", common.BytesToHash(cRoot).Hex())
	if err != nil {
		log.Info("no contract for addr", "addr", addr.Hex())
		return nil, err
	}

	t, err := state.contractTrieCache.OpenTrie(common.BytesToHash(cRoot))
	if err != nil {
		//log.Error("open contract trie failed", "")
		return nil, err
	}

	return t, err
}

func (state *AccountStateDB) ContractExist(addr common.Address) bool {
	//not find if the err isn't nil
	if _, err := state.getContractKV(addr); err != nil {
		return false
	} else {
		return true
	}
}

//not save the return data if there is an error. only save it in DB when commit in the end
func (state *AccountStateDB) PutContract(addr common.Address, v reflect.Value) error {
	if !v.IsValid() || v.IsNil() {
		log.Warn("invalid contract data", "data", v)
		return errors.New("invalid contract data")
	}

	state.contractData[addr] = v
	return nil
}

func (state *AccountStateDB) GetContract(addr common.Address, vType reflect.Type) (v reflect.Value, err error) {
	v = state.contractData[addr]
	if v.IsValid() && !v.IsNil() {
		return
	}

	//log.Info("get contract", "addr", addr)
	kv, err := state.getContractKV(addr)
	if err != nil {
		return reflect.Value{}, err
	}

	nContract := reflect.New(vType)
	//change kv to value
	if err = json_kv.KV2JsonObj(kv, nContract.Interface()); err != nil {
		log.Debug("init contract error form db when call contract function")
		return reflect.Value{}, err
	}
	state.contractData[addr] = nContract

	return nContract, err
}

//get contract data
func (state *AccountStateDB) getContractKV(addr common.Address) (kv map[string]string, err error) {
	t, err := state.getContractTrie(addr)
	if err != nil {
		return nil, err
	}

	kv = map[string]string{}
	it := trie.NewIterator(t.NodeIterator(nil))

	for it.Next() {
		cAddr, key := GetContractAddrAndKey(t.GetKey(it.Key))
		value := it.Value
		mpt_log.Debug("get contract", string(key), string(value), "pre state", state.preStateRoot.Hex())
		if addr.IsEqual(cAddr) {
			kv[string(key)] = string(value)
		} else {
			log.Error("got invalid kv from contract mpt", "passKey", key, "contract addr", addr.Hex())
		}
	}

	if len(kv) == 0 {
		return nil, errors.New(fmt.Sprintf("contract %v not exist", addr))
	}
	return kv, nil
}

//func (state *AccountStateDB) ProcessHeader(header model.AbstractHeader) (receipt model.AbstractReceipt, err error) {
//	log.Debug("add reward to coinbase address", "addr", header.CoinBaseAddress().Hex())
//	if state.IsEmptyAccount(header.CoinBaseAddress()) {
//		if err = state.NewAccountState(header.CoinBaseAddress()); err != nil {
//			return
//		}
//	}
//	return nil, state.AddBalance(header.CoinBaseAddress(), big.NewInt(20*consts.DIP))
//}

// add chain reader out side
func NewAccountStateDB(preStateRoot common.Hash, db StateStorage) (*AccountStateDB, error) {
	tr, err := db.OpenTrie(preStateRoot)
	if err != nil {
		return nil, err
	}
	stateDB := &AccountStateDB{
		preStateRoot:   preStateRoot,
		blockStateTrie: tr,
		storage:        db,

		contractTrieCache:     NewStateStorageWithCache(db.DiskDB()),
		contractData:          map[common.Address]reflect.Value{},
		finalisedContractRoot: map[common.Address]common.Hash{},
		stateChangeList:       newStateChangeList(),
	}
	return stateDB, nil
}

// Copy creates a deep, independent copy of the state.
// Snapshots of the copied state cannot be applied to the copy.
func (state *AccountStateDB) Copy() *AccountStateDB {
	state.lock.Lock()
	defer state.lock.Unlock()
	//todo copy dont maintain the changelist, implement later
	// Copy all the basic fields, initialize the memory ones
	statedb := &AccountStateDB{
		preStateRoot:    state.preStateRoot,
		blockStateTrie:  state.storage.CopyTrie(state.blockStateTrie),
		storage:         state.storage,
		stateChangeList: newStateChangeList(),

		contractTrieCache:     NewStateStorageWithCache(state.storage.DiskDB()),
		contractData:          map[common.Address]reflect.Value{},
		finalisedContractRoot: map[common.Address]common.Hash{},
		//todo: if there is a question because not copy early contract in here
	}
	return statedb
}

// Snapshot returns an identifier for the Current revision of the state.
func (state *AccountStateDB) Snapshot() int {
	id := state.nextRevisionId
	state.nextRevisionId++
	state.validRevisions = append(state.validRevisions, revision{id, state.stateChangeList.length()})
	return id
}

// RevertToSnapshot reverts all state changes made since the given revision.
func (state *AccountStateDB) RevertToSnapshot(revid int) {
	idx := sort.Search(len(state.validRevisions), func(i int) bool {
		return state.validRevisions[i].id >= revid
	})
	if idx == len(state.validRevisions) || state.validRevisions[idx].id != revid {
		panic(fmt.Errorf("revision id %v cannot be reverted", revid))
	}
	snapshot := state.validRevisions[idx].changeIndex
	state.stateChangeList.revert(state, snapshot)
	state.validRevisions = state.validRevisions[:idx]
}

func (state *AccountStateDB) IsEmptyAccount(addr common.Address) bool {
	_, err := state.GetNonce(addr)
	if err != nil {
		return true
	}
	return false
}

func (state *AccountStateDB) GetNonce(addr common.Address) (uint64, error) {
	//log.Info("AccountStateDB GetNonce the addr is: ","addr",addr.Hex())
	enc, err1 := state.blockStateTrie.TryGet(GetNonceKey(addr))
	var res uint64
	if err1 != nil {
		return res, err1
	}
	if len(enc) == 0 {
		return res, g_error.AccountNotExist
	}
	err2 := rlp.DecodeBytes(enc, &res)
	if err2 != nil {
		return res, err2
	}
	return res, nil
}

func (state *AccountStateDB) GetBalance(addr common.Address) (*big.Int, error) {
	//log.Info("AccountStateDB GetBalance the addr is: ","addr",addr.Hex())
	empty := state.IsEmptyAccount(addr)
	if empty {
		return nil, g_error.AccountNotExist
	}
	enc, err1 := state.blockStateTrie.TryGet(GetBalanceKey(addr))
	if err1 != nil {
		return nil, err1
	}
	res := new(big.Int)
	err := rlp.DecodeBytes(enc, &res)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (state *AccountStateDB) GetTimeLock(addr common.Address) (*big.Int, error) {
	empty := state.IsEmptyAccount(addr)
	if empty {
		return nil, g_error.AccountNotExist
	}
	res := new(big.Int)
	enc, err1 := state.blockStateTrie.TryGet(GetTimeLockKey(addr))
	if err1 != nil {
		return nil, err1
	}
	err2 := rlp.DecodeBytes(enc, &res)
	if err2 != nil {
		return res, err2
	}
	return res, nil
}

func (state *AccountStateDB) GetHashLock(addr common.Address) (common.Hash, error) {
	empty := state.IsEmptyAccount(addr)
	if empty {
		return common.Hash{}, g_error.AccountNotExist
	}
	res := common.Hash{}
	enc, err1 := state.blockStateTrie.TryGet(GetHashLockKey(addr))
	if err1 != nil {
		return res, err1
	}
	err2 := rlp.DecodeBytes(enc, &res)
	if err2 != nil {
		return res, err2
	}
	return res, nil
}

func (state *AccountStateDB) GetCommitNum(addr common.Address) (uint64, error) {
	empty := state.IsEmptyAccount(addr)
	if empty {
		return uint64(0), g_error.AccountNotExist
	}
	res := uint64(0)
	enc, err1 := state.blockStateTrie.TryGet(GetCommitNumKey(addr))
	if err1 != nil {
		return uint64(0), err1
	}
	err2 := rlp.DecodeBytes(enc, &res)
	if err2 != nil {
		return res, err2
	}
	return res, nil
}

func (state *AccountStateDB) GetVerifyNum(addr common.Address) (uint64, error) {
	empty := state.IsEmptyAccount(addr)
	if empty {
		return uint64(0), g_error.AccountNotExist
	}
	res := uint64(0)
	enc, err1 := state.blockStateTrie.TryGet(GetVerifyNumKey(addr))
	if err1 != nil {
		return uint64(0), err1
	}
	err2 := rlp.DecodeBytes(enc, &res)
	if err2 != nil {
		return res, err2
	}
	return res, nil
}

func (state *AccountStateDB) GetPerformance(addr common.Address) (uint64, error) {
	empty := state.IsEmptyAccount(addr)
	if empty {
		return performanceInitial, g_error.AccountNotExist
	}
	res := performanceInitial
	enc, err1 := state.blockStateTrie.TryGet(GetPerformanceKey(addr))
	if err1 != nil {
		return res, err1
	}
	err2 := rlp.DecodeBytes(enc, &res)
	if err2 != nil {
		return res, err2
	}
	return res, nil
}

func (state *AccountStateDB) GetLastElect(addr common.Address) (uint64, error) {
	empty := state.IsEmptyAccount(addr)
	if empty {
		return 0, g_error.AccountNotExist
	}
	enc, err1 := state.blockStateTrie.TryGet(GetLastElectKey(addr))
	var res uint64
	if err1 != nil {
		return res, err1
	}
	err2 := rlp.DecodeBytes(enc, &res)
	if err2 != nil {
		return res, err2
	}
	return res, nil
}

//func (state *AccountStateDB) GetContractRoot(addr common.Address) (common.Hash, error) {
//    empty := state.IsEmptyAccount(addr)
//    if empty {
//        return common.Hash{}, g_error.AccountNotExist
//    }
//    res := common.Hash{}
//    enc, err1 := state.blockStateTrie.TryGet(GetContractRootKey(addr))
//    if err1 != nil {
//        return res, err1
//    }
//    err2 := rlp.DecodeBytes(enc, &res)
//    if err2 != nil {
//        return res, err2
//    }
//    return res, nil
//}

func (state *AccountStateDB) GetDataRoot(addr common.Address) (common.Hash, error) {
	empty := state.IsEmptyAccount(addr)
	if empty {
		return common.Hash{}, g_error.AccountNotExist
	}
	res := common.Hash{}
	enc, err1 := state.blockStateTrie.TryGet(GetDataRootKey(addr))
	if err1 != nil {
		return res, err1
	}
	err2 := rlp.DecodeBytes(enc, &res)
	if err2 != nil {
		return res, err2
	}
	return res, nil
}
func (state *AccountStateDB) GetStake(addr common.Address) (*big.Int, error) {
	empty := state.IsEmptyAccount(addr)
	if empty {
		return big.NewInt(0), g_error.AccountNotExist
	}
	res := big.NewInt(0)
	enc, err1 := state.blockStateTrie.TryGet(GetStakeKey(addr))
	if err1 != nil {
		return res, err1
	}
	err2 := rlp.DecodeBytes(enc, &res)
	if err2 != nil {
		return res, err2
	}
	return res, nil
}

func (state *AccountStateDB) SetBalance(addr common.Address, amount *big.Int) error {
	old, _ := state.GetBalance(addr)
	err := state.setBalance(addr, amount)
	if err != nil {
		return err
	}
	state.stateChangeList.append(balanceChange{Account: &addr, Prev: old, Current: amount, ChangeType: BalanceChange})
	return nil
}

//setBalance do not change the changelist, usually called by the revert operation.
func (state *AccountStateDB) setBalance(addr common.Address, amount *big.Int) error {
	empty := state.IsEmptyAccount(addr)
	if empty {
		return g_error.AccountNotExist
	}

	if amount.Cmp(big.NewInt(0)) < 0 {
		log.Debug("set address balance failed", "addr", addr.Hex(), "amount", amount)
		return g_error.BalanceNegErr
	}

	mpt_log.Debug("setBalance", "addr", addr.Hex(), "v", amount, "pre state", state.preStateRoot.Hex())
	newEnc, _ := rlp.EncodeToBytes(amount)
	balanceKey := GetBalanceKey(addr)
	//log.Debug("SetBalance", "balanceKey", hexutil.Encode(balanceKey), "amount", amount.String())
	err := state.blockStateTrie.TryUpdate(balanceKey, newEnc)
	if err != nil {
		return err
	}
	return nil
}

func (state *AccountStateDB) AddBalance(addr common.Address, amount *big.Int) error {
	value, err := state.GetBalance(addr)
	if err != nil {
		return err
	}
	newValue := big.NewInt(0).Add(value, amount)
	err = state.SetBalance(addr, newValue)
	if err != nil {
		return err
	}
	return nil
}

func (state *AccountStateDB) SubBalance(addr common.Address, amount *big.Int) error {
	value, err := state.GetBalance(addr)
	if err != nil {
		return err
	}

	newValue := big.NewInt(0).Sub(value, amount)

	err = state.SetBalance(addr, newValue)
	if err != nil {
		return err
	}
	return nil
}

func (state *AccountStateDB) SetNonce(addr common.Address, amount uint64) error {
	old, _ := state.GetNonce(addr)
	err := state.setNonce(addr, amount)
	if err != nil {
		return err
	}
	state.stateChangeList.append(nonceChange{Account: &addr, Prev: old, Current: amount, ChangeType: NonceChange})
	return nil
}

//setNonce do not change the changelist, usually called by the revert operation.
func (state *AccountStateDB) setNonce(addr common.Address, amount uint64) error {
	empty := state.IsEmptyAccount(addr)
	if empty {
		return g_error.AccountNotExist
	}

	mpt_log.Debug("setNonce", "addr", addr.Hex(), "v", amount, "pre state", state.preStateRoot.Hex())
	newEnc, _ := rlp.EncodeToBytes(amount)
	err := state.blockStateTrie.TryUpdate(GetNonceKey(addr), newEnc)
	if err != nil {
		return err
	}
	return nil
}

func (state *AccountStateDB) AddNonce(addr common.Address, amount uint64) error {
	empty := state.IsEmptyAccount(addr)
	if empty {
		return g_error.AccountNotExist
	}
	var nonce uint64
	enc, _ := state.blockStateTrie.TryGet(GetNonceKey(addr))
	rlp.DecodeBytes(enc, &nonce)
	newNonce := nonce + amount
	state.SetNonce(addr, newNonce)
	return nil
}

func (state *AccountStateDB) SetTimeLock(addr common.Address, timeLock *big.Int) error {
	old, _ := state.GetTimeLock(addr)
	err := state.setTimeLock(addr, timeLock)
	if err != nil {
		return err
	}
	state.stateChangeList.append(timeLockChange{Account: &addr, Prev: old, Current: timeLock, ChangeType: TimeLockChange})
	return nil
}

//setTimeLock do not change the changelist, usually called by the revert operation.
func (state *AccountStateDB) setTimeLock(addr common.Address, timeLock *big.Int) error {
	empty := state.IsEmptyAccount(addr)
	if empty {
		return g_error.AccountNotExist
	}

	mpt_log.Debug("setTimeLock", "addr", addr.Hex(), "v", timeLock, "pre state", state.preStateRoot.Hex())
	newEnc, _ := rlp.EncodeToBytes(timeLock)
	err := state.blockStateTrie.TryUpdate(GetTimeLockKey(addr), newEnc)
	if err != nil {
		return err
	}
	return nil
}

func (state *AccountStateDB) SetHashLock(addr common.Address, hashLock common.Hash) error {
	old, _ := state.GetHashLock(addr)
	err := state.setHashLock(addr, hashLock)
	if err != nil {
		return err
	}
	state.stateChangeList.append(hashLockChange{Account: &addr, Prev: old, Current: hashLock, ChangeType: HashLockChange})
	return nil
}

//setHashLock do not change the changelist, usually called by the revert operation.
func (state *AccountStateDB) setHashLock(addr common.Address, hashLock common.Hash) error {
	empty := state.IsEmptyAccount(addr)
	if empty {
		return g_error.AccountNotExist
	}

	mpt_log.Debug("setHashLock", "addr", addr.Hex(), "v", hashLock.Hex(), "pre state", state.preStateRoot.Hex())
	newEnc, _ := rlp.EncodeToBytes(hashLock)
	err := state.blockStateTrie.TryUpdate(GetHashLockKey(addr), newEnc)
	if err != nil {
		return err
	}
	return nil
}

//func (state *AccountStateDB) SetContractRoot(addr common.Address, contractRoot common.Hash) error {
//    old, _ := state.GetContractRoot(addr)
//    err := state.setContractRoot(addr, contractRoot)
//    if err != nil {
//        return err
//    }
//    state.stateChangeList.append(contractRootChange{Account: &addr, Prev: old, Current: contractRoot, ChangeType: ContractRootChange})
//    return nil
//}

//setContractRoot do not change the changelist, usually called by the revert operation.
//func (state *AccountStateDB) setContractRoot(addr common.Address, contractRoot common.Hash) error {
//    empty := state.IsEmptyAccount(addr)
//    if empty {
//        return g_error.AccountNotExist
//    }
//    newEnc, _ := rlp.EncodeToBytes(contractRoot)
//    err := state.blockStateTrie.TryUpdate(GetContractRootKey(addr), newEnc)
//    if err != nil {
//        return err
//    }
//    return nil
//}

func (state *AccountStateDB) SetDataRoot(addr common.Address, dataRoot common.Hash) error {
	old, _ := state.GetDataRoot(addr)
	err := state.setDataRoot(addr, dataRoot)
	if err != nil {
		return err
	}
	state.stateChangeList.append(dataRootChange{Account: &addr, Prev: old, Current: dataRoot, ChangeType: DataRootChange})
	return nil
}

//setDataRoot do not change the changelist, usually called by the revert operation.
func (state *AccountStateDB) setDataRoot(addr common.Address, dataRoot common.Hash) error {
	empty := state.IsEmptyAccount(addr)
	if empty {
		return g_error.AccountNotExist
	}

	mpt_log.Debug("setDataRoot", "addr", addr.Hex(), "v", dataRoot.Hex(), "pre state", state.preStateRoot.Hex())
	newEnc, _ := rlp.EncodeToBytes(dataRoot)
	err := state.blockStateTrie.TryUpdate(GetDataRootKey(addr), newEnc)
	if err != nil {
		return err
	}
	return nil
}

func (state *AccountStateDB) SetStake(addr common.Address, amount *big.Int) error {
	old, _ := state.GetStake(addr)
	err := state.setStake(addr, amount)
	if err != nil {
		return err
	}
	state.stateChangeList.append(stakeChange{Account: &addr, Prev: old, Current: amount, ChangeType: StakeChange})
	return nil
}

//setStake do not change the changelist, usually called by the revert operation.
func (state *AccountStateDB) setStake(addr common.Address, amount *big.Int) error {
	empty := state.IsEmptyAccount(addr)
	if empty {
		return g_error.AccountNotExist
	}

	mpt_log.Debug("setStake", "addr", addr.Hex(), "v", amount, "pre state", state.preStateRoot.Hex())
	newEnc, _ := rlp.EncodeToBytes(amount)
	err := state.blockStateTrie.TryUpdate(GetStakeKey(addr), newEnc)
	if err != nil {
		return err
	}
	return nil
}

func (state *AccountStateDB) AddStake(addr common.Address, amount *big.Int) error {
	value, err := state.GetStake(addr)
	if err != nil {
		return err
	}
	log.Info("the old stake is:", "value", value)
	newValue := big.NewInt(0).Add(value, amount)
	log.Info("the new stake is:", "newValue", newValue)
	return state.SetStake(addr, newValue)
}

func (state *AccountStateDB) SubStake(addr common.Address, amount *big.Int) error {
	value, err := state.GetStake(addr)
	if err != nil {
		return err
	}
	newValue := big.NewInt(0).Sub(value, amount)
	return state.SetStake(addr, newValue)
}

func (state *AccountStateDB) SetCommitNum(addr common.Address, amount uint64) error {
	old, _ := state.GetCommitNum(addr)
	err := state.setCommitNum(addr, amount)
	if err != nil {
		return err
	}
	state.stateChangeList.append(commitNumChange{Account: &addr, Prev: old, Current: amount, ChangeType: CommitNumChange})
	return nil
}

//setCommitNum do not change the changelist, usually called by the revert operation.
func (state *AccountStateDB) setCommitNum(addr common.Address, amount uint64) error {
	empty := state.IsEmptyAccount(addr)
	if empty {
		return g_error.AccountNotExist
	}

	mpt_log.Debug("setCommitNum", "addr", addr.Hex(), "v", amount, "pre state", state.preStateRoot.Hex())
	newEnc, _ := rlp.EncodeToBytes(amount)
	err := state.blockStateTrie.TryUpdate(GetCommitNumKey(addr), newEnc)
	if err != nil {
		return err
	}
	return nil
}

func (state *AccountStateDB) SetPerformance(addr common.Address, amount uint64) error {
	old, _ := state.GetPerformance(addr)
	err := state.setPerformance(addr, amount)
	if err != nil {
		return err
	}
	state.stateChangeList.append(performanceChange{&addr, old, amount, PerformanceChange})
	return nil
}

func (state *AccountStateDB) setPerformance(addr common.Address, amount uint64) error {
	empty := state.IsEmptyAccount(addr)
	if empty {
		return g_error.AccountNotExist
	}

	mpt_log.Debug("setPerformance", "addr", addr.Hex(), "v", amount, "pre state", state.preStateRoot.Hex())
	newEnc, _ := rlp.EncodeToBytes(amount)
	err := state.blockStateTrie.TryUpdate(GetPerformanceKey(addr), newEnc)
	if err != nil {
		return err
	}
	return nil
}

func (state *AccountStateDB) SetVerifyNum(addr common.Address, amount uint64) error {
	old, _ := state.GetVerifyNum(addr)
	err := state.setVerifyNum(addr, amount)
	if err != nil {
		return err
	}

	state.stateChangeList.append(verifyNumChange{Account: &addr, Prev: old, Current: amount, ChangeType: VerifyNumChange})
	return nil
}

//setVerifyNum do not change the changelist, usually called by the revert operation.
func (state *AccountStateDB) setVerifyNum(addr common.Address, amount uint64) error {
	empty := state.IsEmptyAccount(addr)
	if empty {
		return g_error.AccountNotExist
	}

	mpt_log.Debug("setVerifyNum", "addr", addr.Hex(), "v", amount, "pre state", state.preStateRoot.Hex())
	newEnc, _ := rlp.EncodeToBytes(amount)
	err := state.blockStateTrie.TryUpdate(GetVerifyNumKey(addr), newEnc)
	if err != nil {
		return err
	}
	return nil
}

func (state *AccountStateDB) SetLastElect(addr common.Address, blockID uint64) error {
	old, _ := state.GetLastElect(addr)
	err := state.setLastElect(addr, blockID)
	if err != nil {
		return err
	}
	state.stateChangeList.append(lastElectChange{Account: &addr, Prev: old, Current: blockID, ChangeType: LastElectChange})
	return nil
}

//setLastElect do not change the changelist, usually called by the revert operation.
func (state *AccountStateDB) setLastElect(addr common.Address, blockID uint64) error {
	empty := state.IsEmptyAccount(addr)
	if empty {
		return g_error.AccountNotExist
	}

	mpt_log.Debug("setLastElect", "addr", addr.Hex(), "v", blockID, "pre state", state.preStateRoot.Hex())
	encBlockId, _ := rlp.EncodeToBytes(blockID)
	err := state.blockStateTrie.TryUpdate(GetLastElectKey(addr), encBlockId)
	if err != nil {
		return err
	}
	return nil
}

func (state *AccountStateDB) NewAccountState(addr common.Address) error {
	_, err := state.newAccountState(addr)
	if err != nil {
		return err
	}
	state.stateChangeList.append(newAccountChange{Account: &addr, ChangeType: NewAccountChange})
	return nil
}

func (state *AccountStateDB) newAccountState(addr common.Address) (acc *account, err error) {
	tempAccount := account{Nonce: 0, Balance: big.NewInt(0), TimeLock: big.NewInt(0), Stake: big.NewInt(0), CommitNum: uint64(0), VerifyNum: uint64(0), Performance: performanceInitial, LastElect: uint64(0), HashLock: common.Hash{}, DataRoot: common.Hash{}}
	err = state.blockStateTrie.TryUpdate(GetNonceKey(addr), tempAccount.NonceBytes())
	if err != nil {
		return nil, err
	}
	err = state.blockStateTrie.TryUpdate(GetBalanceKey(addr), tempAccount.BalanceBytes())
	if err != nil {
		return nil, err
	}
	err = state.blockStateTrie.TryUpdate(GetStakeKey(addr), tempAccount.StakeBytes())
	if err != nil {
		return nil, err
	}
	err = state.blockStateTrie.TryUpdate(GetCommitNumKey(addr), tempAccount.CommitNumBytes())
	if err != nil {
		return nil, err
	}
	err = state.blockStateTrie.TryUpdate(GetVerifyNumKey(addr), tempAccount.VerifyNumBytes())
	if err != nil {
		return nil, err
	}
	err = state.blockStateTrie.TryUpdate(GetPerformanceKey(addr), tempAccount.PerformanceBytes())
	if err != nil {
		return nil, err
	}
	err = state.blockStateTrie.TryUpdate(GetHashLockKey(addr), tempAccount.HashLockBytes())
	if err != nil {
		return nil, err
	}
	err = state.blockStateTrie.TryUpdate(GetTimeLockKey(addr), tempAccount.TimeLockBytes())
	if err != nil {
		return nil, err
	}
	//err = state.blockStateTrie.TryUpdate(GetContractRootKey(addr), tempAccount.ContractRootBytes())
	//if err != nil {
	//    return nil, err
	//}
	err = state.blockStateTrie.TryUpdate(GetDataRootKey(addr), tempAccount.DataRootBytes())
	if err != nil {
		return nil, err
	}
	err = state.blockStateTrie.TryUpdate(GetLastElectKey(addr), tempAccount.LastElectBytes())
	if err != nil {
		return nil, err
	}
	acc = &tempAccount
	return

}

func (state *AccountStateDB) DeleteAccountState(addr common.Address) error {
	err := state.deleteAccountState(addr)
	if err != nil {
		return err
	}
	state.stateChangeList.append(deleteAccountChange{Account: &addr, ChangeType: DeleteAccountChange})
	return nil
}

// deleteStateObject removes the given object from the state trie.
func (state *AccountStateDB) deleteAccountState(addr common.Address) (err error) {
	err = state.blockStateTrie.TryDelete(GetNonceKey(addr))
	if err != nil {
		return err
	}
	err = state.blockStateTrie.TryDelete(GetBalanceKey(addr))
	if err != nil {
		return err
	}
	err = state.blockStateTrie.TryDelete(GetTimeLockKey(addr))
	if err != nil {
		return err
	}
	err = state.blockStateTrie.TryDelete(GetHashLockKey(addr))
	if err != nil {
		return err
	}
	err = state.blockStateTrie.TryDelete(GetContractRootKey(addr))
	if err != nil {
		return err
	}
	err = state.blockStateTrie.TryDelete(GetDataRootKey(addr))
	if err != nil {
		return err
	}
	err = state.blockStateTrie.TryDelete(GetStakeKey(addr))
	if err != nil {
		return err
	}
	err = state.blockStateTrie.TryDelete(GetCommitNumKey(addr))
	if err != nil {
		return err
	}
	err = state.blockStateTrie.TryDelete(GetVerifyNumKey(addr))
	if err != nil {
		return err
	}
	err = state.blockStateTrie.TryDelete(GetPerformanceKey(addr))
	if err != nil {
		return err
	}
	err = state.blockStateTrie.TryDelete(GetLastElectKey(addr))
	if err != nil {
		return err
	}
	return nil
}

func (state *AccountStateDB) GetAccountState(addr common.Address) (*account, error) {
	account := new(account)
	if state.IsEmptyAccount(addr) {
		return nil, g_error.AccountNotExist
	}
	nonce, err := state.GetNonce(addr)
	if err != nil {
		return nil, err
	}
	account.Nonce = nonce
	balance, err := state.GetBalance(addr)
	if err != nil {
		return nil, err
	}
	account.Balance = balance
	hashLock, err := state.GetHashLock(addr)
	if err != nil {
		return nil, err
	}
	account.HashLock = hashLock
	timeLock, err := state.GetTimeLock(addr)
	if err != nil {
		return nil, err
	}
	account.TimeLock = timeLock
	stake, err := state.GetStake(addr)
	if err != nil {
		return nil, err
	}
	account.Stake = stake
	commitNum, err := state.GetCommitNum(addr)
	if err != nil {
		return nil, err
	}
	account.CommitNum = commitNum
	verifyNum, err := state.GetVerifyNum(addr)
	if err != nil {
		return nil, err
	}
	account.VerifyNum = verifyNum
	lastElect, err := state.GetLastElect(addr)
	if err != nil {
		return nil, err
	}
	account.LastElect = lastElect
	performance, err := state.GetPerformance(addr)
	if err != nil {
		return nil, err
	}
	account.Performance = performance
	//contractRoot, err := state.GetContractRoot(addr)
	//if err != nil {
	//    return nil, err
	//}
	//account.ContractRoot = contractRoot
	dataRoot, err := state.GetDataRoot(addr)
	if err != nil {
		return nil, err
	}
	account.DataRoot = dataRoot
	// Insert into the live set.
	return account, nil
}

// commit contract data
func (state *AccountStateDB) commitContractData() error {
	for addr, root := range state.finalisedContractRoot {
		mpt_log.Debug("commit contract", "addr", addr.Hex(), "root", root.Hex(), "pre state", state.preStateRoot.Hex())
		//log.Info("commit contract trie", "root", root.Hex())
		if err := state.contractTrieCache.TrieDB().Commit(root, false); err != nil {
			return err
		}
	}
	//log.Info("commit contract trie end")
	return nil
}

// put contract data to trie
func (state *AccountStateDB) putContractDataToTrie(addr common.Address, data []byte) (StateTrie, error) {
	mpt_log.Info("put contract", "addr", addr)
	ct, err := state.getContractTrie(addr)
	//check err first and return err if not find trie, otherwise there isn't this trie if th ct is nil
	if err != nil && !strings.Contains(err.Error(), "missing trie node") {
		log.Warn("can't get contract trie from db", "err", err)
		return nil, err
	}
	// todo ct is nil?
	//else if ct == nil {
	//
	//}
	kv, err := json_kv.JsonBytes2KV(data)
	if err != nil {
		return nil, err
	}
	for k, v := range kv {
		mpt_log.Debug("putContractDataToTrie", "k", k, "v", v, "pre state", state.preStateRoot.Hex())
		if err := ct.TryUpdate(GetContractFieldKey(addr, k), []byte(v)); err != nil {
			return nil, err
		}
	}
	return ct, nil
}

func (state *AccountStateDB) Commit() (common.Hash, error) {

	//must finalise ,otherwise the state root of contract will be incorrect
	fStateRoot, err := state.Finalise()
	if err != nil {
		return common.Hash{}, err
	}

	// commit contracts
	if err := state.commitContractData(); err != nil {
		//need clear finalised contract root if yes. At the same time, you'd better throw away this accountStateDB
		//state.resetThisStateDB()
		log.Warn("commit contract data failed", "err", err)
		return common.Hash{}, err
	}

	//it's difficult to do reference in here,because we don't know if the data of the leaf callback is contract data,maybe
	//balance or other data
	//if root, err := state.blockStateTrie.Commit(nil); err != nil {
	//    //state.resetThisStateDB()
	//    return common.Hash{}, err
	//} else {
	//    if !fStateRoot.IsEqual(root) {
	//        // maybe panic here
	//        return common.Hash{}, errors.New("finalised state root not match commit state root")
	//    }
	//have committed in the finalise
	err = state.storage.TrieDB().Commit(fStateRoot, false)
	return fStateRoot, err
	//}
}

// check if have been finalised
func (state *AccountStateDB) finalised() bool {
	return state.alreadyFinalised
}

func (state *AccountStateDB) finaliseContractData() error {
	for addr, data := range state.contractData {
		ct, err := state.putContractDataToTrie(addr, util.StringifyJsonToBytes(data.Interface()))
		if err != nil {
			return err
		}

		// You must commit trie to memory, and only use commit trie db in the commit.
		ch, err := ct.Commit(nil)
		if err != nil {
			return err
		}
		mpt_log.Info("finaliseContractData update contract root", "contract addr", addr.Hex(), "root", ch.Hex())
		if err := state.blockStateTrie.TryUpdate(GetContractRootKey(addr), ch.Bytes()); err != nil {
			// change blockStateTrie to origin pre hash？If you want, clear the finalised contract root. But it is best to discard the AccountStateDB directly after the error is reported.
			//state.resetThisStateDB()
			log.Error("Commit update contract root failed", "err", err)
			return err
		}
		state.finalisedContractRoot[addr] = ch
	}
	return nil
}

// deleteEmptyAccount bool true.
// Doing a trie commit logic here is more complicated, so don't consider committing for the time being.
// If finalised, don't change any state outside, otherwise there will be problems.
func (state *AccountStateDB) Finalise() (result common.Hash, err error) {

	if state.finalised() {
		result = state.blockStateTrie.Hash()
		mpt_log.Debug("Finalise", "cur root", result.Hex(), "pre state", state.preStateRoot.Hex())
		return result, nil
	}
	// finalise contracts
	if err := state.finaliseContractData(); err != nil {
		// change blockStateTrie to origin pre hash？
		// If you want, clear the finalised contract root. But it is best to discard the AccountStateDB directly after the error is reported.
		//state.resetThisStateDB()
		mpt_log.Debug("Finalise failed", "err", err, "pre state", state.preStateRoot.Hex())
		result = common.Hash{}
		return result, err
	}

	state.alreadyFinalised = true

	result, err = state.blockStateTrie.Commit(nil)
	mpt_log.Debug("Finalise", "cur root", result.Hex(), "pre state", state.preStateRoot.Hex())
	return
}

//todo these processes are removed afterwards。
// todo Write a unit test for each transaction to cover all situations
func (state *AccountStateDB) ProcessTx(tx model.AbstractTransaction, height uint64) (err error) {
	// All transactions must be done with processBasicTx, and transactionBasicTx only deducts transaction fees. Amount is selectively handled in each type of transaction
	err = state.processBasicTx(tx)
	if err != nil {
		log.Debug("processBasicTx failed", "err", err)
		return
	}
	switch tx.GetType() {
	case common.AddressTypeNormal:
		err = state.processNormalTx(tx)
	case common.AddressTypeCross:
		err = state.processCrossTx(tx)
	case common.AddressTypeERC20:
		err = state.processERC20Tx(tx, height)
		// Verifier relate transaction processor
	case common.AddressTypeStake:
		err = state.processStakeTx(tx)
	case common.AddressTypeCancel:
		err = state.processCancelTx(tx, height)
	case common.AddressTypeUnStake:
		err = state.processUnStakeTx(tx)
	case common.AddressTypeEvidence:
		err = state.processEvidenceTx(tx)
	case common.AddressTypeEarlyReward:
		err = state.processEarlyTokenTx(tx, height)
	case common.AddressTypeSmartContract:
		err = state.ProcessSmartContract(tx, height)
	default:
		err = g_error.UnknownTxTypeErr
	}
	return
}

func (state *AccountStateDB) processBasicTx(tx model.AbstractTransaction) (err error) {
	sender, err := tx.Sender(nil)
	receiver := *(tx.To())
	if err != nil {
		log.Debug("get tx sender failed", "err", err)
		return
	}
	if sender.IsEmpty() || receiver.IsEmpty() {
		log.Warn("tx ("+tx.CalTxId().Hex()+") but sender or receiver is empty", "sender", sender, "receiver", receiver)
		return SenderOrReceiverIsEmptyErr
	}
	if empty := state.IsEmptyAccount(sender); empty {
		return SenderNotExistErr
	}

	curNonce, _ := state.GetNonce(sender)
	if tx.Nonce() != curNonce {
		log.Info("tx nonce not match", "tx n", tx.Nonce(), "cur account nonce", curNonce)
		return g_error.ErrTxNonceNotMatch
	}
	/*	if empty := state.IsEmptyAccount(receiver); empty {
			return ReceiverNotExistErr
		}*/
	err = state.SubBalance(sender, tx.Fee())
	if err != nil {
		return
	}
	err = state.AddNonce(sender, uint64(1))
	if err != nil {
		return
	}
	return
}

func (state *AccountStateDB) processNormalTx(tx model.AbstractTransaction) (err error) {

	sender, _ := tx.Sender(nil)
	receiver := *(tx.To())
	if empty := state.IsEmptyAccount(receiver); empty {
		state.NewAccountState(receiver)
	}
	err = state.SubBalance(sender, tx.Amount())
	if err != nil {
		return
	}
	err = state.AddBalance(receiver, tx.Amount())
	if err != nil {
		return
	}
	return
}

func (state *AccountStateDB) processCrossTx(tx model.AbstractTransaction) (err error) {
	// TODO:
	return errors.New("not support now")
}

func (state *AccountStateDB) processERC20Tx(tx model.AbstractTransaction, blockHeight uint64) (err error) {
	cProcessor := contract.NewProcessor(state, blockHeight)
	err = cProcessor.Process(tx)
	if err != nil {
		return
	}
	return
}

func (state *AccountStateDB) processEarlyTokenTx(tx model.AbstractTransaction, blockHeight uint64) (err error) {

	cProcessor := contract.NewProcessor(state, blockHeight)
	cProcessor.SetAccountDB(state)
	eData := contract.ParseExtraDataForContract(tx.ExtraData())
	if eData == nil {
		return contract.CanNotParseContractErr
	}

	for _, prohibitFunc := range contract.ProhibitFunction {
		if eData.Action == prohibitFunc {
			return errors.New("can't use this contract function")
		}
	}

	err = cProcessor.Process(tx)
	return
}

