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

package stateprocessor

import (
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	common2 "github.com/dipperin/dipperin-core/core/vm/base"
	"github.com/dipperin/dipperin-core/common/gerror"
	"github.com/dipperin/dipperin-core/common/log"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/common/util/jsonkv"
	"github.com/dipperin/dipperin-core/core/contract"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third_party/trie"
	"github.com/ethereum/go-ethereum/rlp"
	"go.uber.org/zap"
	"math/big"
	"reflect"
	"sort"
	"strings"
	"sync"
)

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
	smartContractData     map[common.Address]map[string][]byte
	logs                  map[common.Hash][]*model.Log

	stateChangeList *StateChangeList
	validRevisions  []revision
	nextRevisionId  int
	lock            sync.Mutex
}

func (state *AccountStateDB) PreStateRoot() common.Hash {
	return state.preStateRoot
}

func (state *AccountStateDB) getContractTrie(addr common.Address) (StateTrie, error) {

	//notice: can't get the trie if the contract root had been changed but not commit
	cRoot, err := state.blockStateTrie.TryGet(GetDataRootKey(addr))
	log.DLogger.Debug("get address contract root", zap.String("addr", addr.Hex()), zap.String("root", common.BytesToHash(cRoot).Hex()))
	if err != nil {
		log.DLogger.Info("no contract for addr", zap.String("addr", addr.Hex()))
		return nil, err
	}

	t, err := state.contractTrieCache.OpenTrie(common.BytesToHash(cRoot))
	if err != nil {
		//log.DLogger.Error("open contract trie failed", "")
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
		log.DLogger.Warn("invalid contract data", zap.Any("data", v))
		return gerror.ErrInvalidContractData
	}
	old := state.contractData[addr]
	state.contractData[addr] = v
	state.stateChangeList.append(contractChange{Account: &addr, Prev: old, Current: v, ChangeType: ContractChange})
	return nil
}

func (state *AccountStateDB) GetContract(addr common.Address, vType reflect.Type) (v reflect.Value, err error) {
	v = state.contractData[addr]
	if v.IsValid() && !v.IsNil() {
		return
	}

	//log.DLogger.Info("get contract", "addr", addr)
	kv, err := state.getContractKV(addr)
	if err != nil {
		return reflect.Value{}, err
	}

	nContract := reflect.New(vType)
	//change kv to value
	if err = jsonkv.KV2JsonObj(kv, nContract.Interface()); err != nil {
		log.DLogger.Debug("init contract error form db when call contract function")
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
		log.DLogger.Debug("get contract", zap.String(string(key), string(value)), zap.String("pre state", state.preStateRoot.Hex()))
		if addr.IsEqual(cAddr) {
			kv[string(key)] = string(value)
		} else {
			log.DLogger.Error("got invalid kv from contract mpt", zap.Uint8s("passKey", key), zap.String("contract addr", addr.Hex()))
		}
	}

	if len(kv) == 0 {
		return nil, gerror.ErrContractNotExist
	}
	return kv, nil
}

//func (state *AccountStateDB) ProcessHeader(Header model.AbstractHeader) (receipt model.AbstractReceipt, err error) {
//	log.DLogger.Debug("add reward to coinbase address", "addr", Header.CoinBaseAddress().Hex())
//	if state.IsEmptyAccount(Header.CoinBaseAddress()) {
//		if err = state.NewAccountState(Header.CoinBaseAddress()); err != nil {
//			return
//		}
//	}
//	return nil, state.AddBalance(Header.CoinBaseAddress(), big.NewInt(20*consts.DIP))
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
		smartContractData:     make(map[common.Address]map[string][]byte),
		finalisedContractRoot: map[common.Address]common.Hash{},
		stateChangeList:       newStateChangeList(),
		logs:                  map[common.Hash][]*model.Log{},
		alreadyFinalised:      false,
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
	enc, err1 := state.blockStateTrie.TryGet(GetNonceKey(addr))
	var res uint64
	if err1 != nil {
		return res, err1
	}
	if len(enc) == 0 {
		return res, gerror.ErrAccountNotExist
	}
	err2 := rlp.DecodeBytes(enc, &res)
	if err2 != nil {
		return res, err2
	}
	return res, nil
}

func (state *AccountStateDB) GetBalance(addr common.Address) (*big.Int, error) {
	empty := state.IsEmptyAccount(addr)
	if empty {
		return nil, gerror.ErrAccountNotExist
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
		return nil, gerror.ErrAccountNotExist
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
		return common.Hash{}, gerror.ErrAccountNotExist
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
		return uint64(0), gerror.ErrAccountNotExist
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
		return uint64(0), gerror.ErrAccountNotExist
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
		return performanceInitial, gerror.ErrAccountNotExist
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
		return 0, gerror.ErrAccountNotExist
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

func (state *AccountStateDB) GetDataRoot(addr common.Address) (common.Hash, error) {

	empty := state.IsEmptyAccount(addr)
	if empty {
		return common.Hash{}, gerror.ErrAccountNotExist
	}
	res := common.Hash{}
	enc, err1 := state.blockStateTrie.TryGet(GetDataRootKey(addr))
	if err1 != nil {

		return res, err1
	}
	res = common.BytesToHash(enc)
	return res, nil
}
func (state *AccountStateDB) GetStake(addr common.Address) (*big.Int, error) {
	empty := state.IsEmptyAccount(addr)
	if empty {
		return big.NewInt(0), gerror.ErrAccountNotExist
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

func (state *AccountStateDB) GetAbi(addr common.Address) ([]byte, error) {
	empty := state.IsEmptyAccount(addr)
	if empty {
		return []byte{}, gerror.ErrAccountNotExist
	}
	res := []byte{}
	enc, err1 := state.blockStateTrie.TryGet(GetAbiKey(addr))
	if err1 != nil {
		return res, err1
	}
	err2 := rlp.DecodeBytes(enc, &res)
	if err2 != nil {
		return res, err2
	}
	return res, nil
}

func (state *AccountStateDB) GetCode(addr common.Address) ([]byte, error) {
	empty := state.IsEmptyAccount(addr)
	if empty {
		return []byte{}, gerror.ErrAccountNotExist
	}
	res := []byte{}
	enc, err1 := state.blockStateTrie.TryGet(GetCodeKey(addr))
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
		return gerror.ErrAccountNotExist
	}

	if amount.Cmp(big.NewInt(0)) < 0 {
		log.DLogger.Debug("set address balance failed", zap.String("addr", addr.Hex()), zap.Any("amount", amount))
		return gerror.ErrBalanceNegative
	}

	log.DLogger.Debug("setBalance", zap.String("addr", addr.Hex()), zap.Any("v", amount), zap.String("pre state", state.preStateRoot.Hex()))
	newEnc, _ := rlp.EncodeToBytes(amount)
	balanceKey := GetBalanceKey(addr)
	//log.DLogger.Debug("SetBalance", "balanceKey", hexutil.Encode(balanceKey), "amount", amount.String())
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
		return gerror.ErrAccountNotExist
	}

	log.DLogger.Debug("setNonce", zap.String("addr", addr.Hex()), zap.Uint64("v", amount), zap.String("pre state", state.preStateRoot.Hex()))
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
		return gerror.ErrAccountNotExist
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
		return gerror.ErrAccountNotExist
	}

	log.DLogger.Debug("setTimeLock", zap.String("addr", addr.Hex()), zap.Any("v", timeLock), zap.String("pre state", state.preStateRoot.Hex()))
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
		return gerror.ErrAccountNotExist
	}

	log.DLogger.Debug("setHashLock", zap.String("addr", addr.Hex()), zap.String("v", hashLock.Hex()), zap.String("pre state", state.preStateRoot.Hex()))
	newEnc, _ := rlp.EncodeToBytes(hashLock)
	err := state.blockStateTrie.TryUpdate(GetHashLockKey(addr), newEnc)
	if err != nil {
		return err
	}
	return nil
}

func (state *AccountStateDB) SetAbi(addr common.Address, abi []byte) error {
	old, _ := state.GetAbi(addr)
	err := state.setAbi(addr, abi)
	if err != nil {
		return err
	}
	state.stateChangeList.append(abiChange{Account: &addr, Prev: old, Current: abi, ChangeType: AbiChange})
	return nil
}

func (state *AccountStateDB) setAbi(addr common.Address, abi []byte) error {
	empty := state.IsEmptyAccount(addr)
	if empty {
		return gerror.ErrAccountNotExist
	}

	log.DLogger.Debug("setAbi", zap.String("addr", addr.Hex()))
	newEnc, _ := rlp.EncodeToBytes(abi)
	err := state.blockStateTrie.TryUpdate(GetAbiKey(addr), newEnc)
	if err != nil {
		return err
	}
	return nil
}

func (state *AccountStateDB) SetCode(addr common.Address, code []byte) error {
	old, _ := state.GetCode(addr)
	err := state.setCode(addr, code)
	if err != nil {
		return err
	}
	state.stateChangeList.append(codeChange{Account: &addr, Prev: old, Current: code, ChangeType: CodeChange})
	return nil
}

func (state *AccountStateDB) setCode(addr common.Address, code []byte) error {
	empty := state.IsEmptyAccount(addr)
	if empty {
		return gerror.ErrAccountNotExist
	}
	log.DLogger.Debug("setCode", zap.String("addr", addr.Hex()))
	newEnc, _ := rlp.EncodeToBytes(code)
	err := state.blockStateTrie.TryUpdate(GetCodeKey(addr), newEnc)
	if err != nil {
		return err
	}
	return nil
}

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
		return gerror.ErrAccountNotExist
	}

	log.DLogger.Debug("setDataRoot", zap.String("addr", addr.Hex()), zap.String("v", dataRoot.Hex()), zap.String("pre state", state.preStateRoot.Hex()))
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
		return gerror.ErrAccountNotExist
	}

	log.DLogger.Debug("setStake", zap.String("addr", addr.Hex()), zap.Any("v", amount), zap.String("pre state", state.preStateRoot.Hex()))
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
	newValue := big.NewInt(0).Add(value, amount)
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
		return gerror.ErrAccountNotExist
	}

	log.DLogger.Debug("setCommitNum", zap.String("addr", addr.Hex()), zap.Uint64("v", amount), zap.String("pre state", state.preStateRoot.Hex()))
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
		return gerror.ErrAccountNotExist
	}

	log.DLogger.Debug("setPerformance", zap.String("addr", addr.Hex()), zap.Uint64("v", amount), zap.String("pre state", state.preStateRoot.Hex()))
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
		return gerror.ErrAccountNotExist
	}

	log.DLogger.Debug("setVerifyNum", zap.String("addr", addr.Hex()), zap.Uint64("v", amount), zap.String("pre state", state.preStateRoot.Hex()))
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
		return gerror.ErrAccountNotExist
	}

	log.DLogger.Debug("setLastElect", zap.String("addr", addr.Hex()), zap.Uint64("v", blockID), zap.String("pre state", state.preStateRoot.Hex()))
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

func (state *AccountStateDB) newContractAccount(addr common.Address) (acc *account, err error) {
	tempAccount := account{Nonce: 0, Balance: big.NewInt(0), TimeLock: big.NewInt(0), HashLock: common.Hash{}, DataRoot: common.Hash{}}
	err = state.blockStateTrie.TryUpdate(GetNonceKey(addr), tempAccount.NonceBytes())
	if err != nil {
		return nil, err
	}
	err = state.blockStateTrie.TryUpdate(GetBalanceKey(addr), tempAccount.BalanceBytes())
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
	err = state.blockStateTrie.TryUpdate(GetDataRootKey(addr), tempAccount.DataRootBytes())
	if err != nil {
		return nil, err
	}
	acc = &tempAccount
	return
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
		return nil, gerror.ErrAccountNotExist
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
		log.DLogger.Debug("commit contract", zap.String("addr", addr.Hex()), zap.String("root", root.Hex()), zap.String("pre state", state.preStateRoot.Hex()))
		//log.DLogger.Info("commit contract trie", "root", root.Hex())
		if err := state.contractTrieCache.TrieDB().Commit(root, false); err != nil {
			return err
		}
	}
	//log.DLogger.Info("commit contract trie end")
	return nil
}

// put contract data to trie
func (state *AccountStateDB) putContractDataToTrie(addr common.Address, data []byte) (StateTrie, error) {
	log.DLogger.Info("put contract", zap.String("addr", addr.String()))
	ct, err := state.getContractTrie(addr)
	//check err first and return err if not find trie, otherwise there isn't this trie if th ct is nil
	if err != nil && !strings.Contains(err.Error(), "missing trie node") {
		log.DLogger.Warn("can't get contract trie from db", zap.Error(err))
		return nil, err
	}
	// todo ct is nil?
	//else if ct == nil {
	//
	//}
	kv, err := jsonkv.JsonBytes2KV(data)
	if err != nil {
		return nil, err
	}
	for k, v := range kv {
		log.DLogger.Debug("putContractDataToTrie", zap.String("k", k), zap.String("v", v), zap.String("pre state", state.preStateRoot.Hex()))
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
		log.DLogger.Warn("commit contract data failed", zap.Error(err))
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
	state.clearChangeList() //TODO not tested
	return fStateRoot, err
	//}
}

// check if have been finalised
func (state *AccountStateDB) finalised() bool {
	return state.alreadyFinalised
}

func (state *AccountStateDB) Finalised() bool {
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
		log.DLogger.Info("finaliseContractData update contract root", zap.String("contract addr", addr.Hex()), zap.String("root", ch.Hex()))
		if err := state.blockStateTrie.TryUpdate(GetDataRootKey(addr), ch.Bytes()); err != nil {
			// change blockStateTrie to origin pre hash？If you want, clear the finalised contract root. But it is best to discard the AccountStateDB directly after the error is reported.
			//state.resetThisStateDB()
			log.DLogger.Error("Commit update contract root failed", zap.Error(err))
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
		log.DLogger.Debug("Finalise", zap.String("cur root", result.Hex()), zap.String("pre state", state.preStateRoot.Hex()))
		return
	}
	// finalise contracts
	if err = state.finaliseContractData(); err != nil {
		// change blockStateTrie to origin pre hash？
		// If you want, clear the finalised contract root. But it is best to discard the AccountStateDB directly after the error is reported.
		//state.resetThisStateDB()
		log.DLogger.Debug("Finalise failed", zap.Error(err), zap.String("pre state", state.preStateRoot.Hex()))
		result = common.Hash{}
		return
	}

	if err = state.finalSmartData(); err != nil {
		log.DLogger.Debug("Finalise smart data failed", zap.Error(err), zap.String("pre state", state.preStateRoot.Hex()))
		result = common.Hash{}
		return
	}

	state.alreadyFinalised = true
	result, err = state.blockStateTrie.Commit(nil)
	log.DLogger.Debug("Finalise", zap.Any("cur root", result), zap.String("pre state", state.preStateRoot.Hex()))
	return
}

func (state *AccountStateDB) IntermediateRoot() (result common.Hash, err error) {
	if state.finalised() {
		result = state.blockStateTrie.Hash()
		log.DLogger.Debug("Finalise", zap.String("cur root", result.Hex()), zap.String("pre state", state.preStateRoot.Hex()))
		return result, nil
	}
	// finalise contracts
	if err = state.finaliseContractData(); err != nil {
		// change blockStateTrie to origin pre hash？
		// If you want, clear the finalised contract root. But it is best to discard the AccountStateDB directly after the error is reported.
		// state.resetThisStateDB()
		log.DLogger.Debug("Finalise failed", zap.Error(err), zap.String("pre state", state.preStateRoot.Hex()))
		result = common.Hash{}
		return
	}

	// finalise smart contracts data
	if err = state.finalSmartData(); err != nil {
		log.DLogger.Debug("Finalise smart data failed", zap.Error(err), zap.String("pre state", state.preStateRoot.Hex()))
		result = common.Hash{}
		return
	}
	return state.blockStateTrie.Commit(nil)
}

//todo these processes are removed afterwards。
// todo Write a unit test for each transaction to cover all situations
/*func (state *AccountStateDB) ProcessTx(tx model.AbstractTransaction, height uint64) (err error) {
	// All transactions must be done with processBasicTx, and transactionBasicTx only deducts transaction fees. Amount is selectively handled in each type of transaction
	err = state.processBasicTx(tx)
	if err != nil {
		log.DLogger.Debug("processBasicTx failed", "err", err)
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
	default:
		err = g_error.ErrUnknownTxType
	}
	return
}*/

func (state *AccountStateDB) setTxReceiptPar(tx model.AbstractTransaction, par *model.ReceiptPara, blockGasUsed *uint64) error {
	if tx.GetType() == common.AddressTypeContractCreate || tx.GetType() == common.AddressTypeContractCall {
		return nil
	}

	root, err := state.IntermediateRoot()
	if err != nil {
		return err
	}

	gasUsed, err := model.IntrinsicGas(tx.ExtraData(), false, false)
	if err != nil {
		return err
	}

	//add block gasUsed
	*blockGasUsed += gasUsed

	par.CumulativeGasUsed = *blockGasUsed
	par.HandlerResult = false
	par.Root = root[:]
	par.Logs = []*model.Log{}
	return nil
}

type TxProcessConfig struct {
	Tx       model.AbstractTransaction
	Header   model.AbstractHeader
	GetHash  common2.GetHashFunc
	GasLimit *uint64
	GasUsed  *uint64
	TxFee    *big.Int
}

func (state *AccountStateDB) ProcessTxNew(conf *TxProcessConfig) (err error) {
	// All transactions must be done with processBasicTx, and transactionBasicTx only deducts transaction fees. Amount is selectively handled in each type of transaction
	if conf.Tx.GetType() != common.AddressTypeContractCall && conf.Tx.GetType() != common.AddressTypeContractCreate {
		err = state.processBasicTx(conf)
		if err != nil {
			log.DLogger.Debug("processBasicTx failed", zap.Error(err))
			return
		}
	}

	var par model.ReceiptPara
	switch conf.Tx.GetType() {
	case common.AddressTypeContractCall:
		par, err = state.ProcessContract(conf, false)
	case common.AddressTypeContractCreate:
		par, err = state.ProcessContract(conf, true)
	case common.AddressTypeNormal:
		err = state.processNormalTx(conf.Tx)
	case common.AddressTypeCross:
		err = state.processCrossTx(conf.Tx)
	case common.AddressTypeERC20:
		err = state.processERC20Tx(conf.Tx, conf.Header.GetNumber())
	case common.AddressTypeStake:
		err = state.processStakeTx(conf.Tx)
	case common.AddressTypeCancel:
		err = state.processCancelTx(conf.Tx, conf.Header.GetNumber())
	case common.AddressTypeUnStake:
		err = state.processUnStakeTx(conf.Tx)
	case common.AddressTypeEvidence:
		err = state.processEvidenceTx(conf.Tx)
	case common.AddressTypeEarlyReward:
		err = state.processEarlyTokenTx(conf.Tx, conf.Header.GetNumber())
	default:
		err = gerror.ErrUnknownTxType
	}
	if err != nil {
		return
	}

	err = state.setTxReceiptPar(conf.Tx, &par, conf.GasUsed)
	if err != nil {
		return
	}

	//updating tx fee
	conf.Tx.PaddingActualTxFee(conf.TxFee)
	conf.Tx.PaddingReceipt(par)
	return
}

func (state *AccountStateDB) processBasicTx(conf *TxProcessConfig) (err error) {
	sender, err := conf.Tx.Sender(nil)
	receiver := *(conf.Tx.To())
	if err != nil {
		log.DLogger.Debug("get Tx sender failed", zap.Error(err))
		return
	}
	if sender.IsEmpty() || receiver.IsEmpty() {
		log.DLogger.Warn("Tx ("+conf.Tx.CalTxId().Hex()+") but sender or receiver is empty", zap.Any("sender", sender), zap.Any("receiver", receiver))
		return gerror.ErrSenderOrReceiverIsEmpty
	}
	if empty := state.IsEmptyAccount(sender); empty {
		return gerror.ErrSenderNotExist
	}

	curNonce, _ := state.GetNonce(sender)
	if conf.Tx.Nonce() != curNonce {
		log.DLogger.Info("Tx nonce not match", zap.Uint64("Tx n", conf.Tx.Nonce()), zap.Uint64("cur account nonce", curNonce))
		return gerror.ErrTxNonceNotMatch
	}
	/*	if empty := state.IsEmptyAccount(receiver); empty {
		return g_error.ErrReceiverNotExist
	}*/
	//calculated gasUsed and sub the fee
	gasUsed, err := model.IntrinsicGas(conf.Tx.ExtraData(), false, false)
	if err != nil {
		return err
	}

	if gasUsed > conf.Tx.GetGasLimit() {
		return gerror.ErrTxGasUsedIsOverGasLimit
	}

	conf.TxFee = big.NewInt(0).Mul(big.NewInt(int64(gasUsed)), conf.Tx.GetGasPrice())
	err = state.SubBalance(sender, conf.TxFee)
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
	return gerror.ErrTxNotSupported
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
			return gerror.ErrProhibitFunctionCalled
		}
	}

	err = cProcessor.Process(tx)
	return
}

func (state *AccountStateDB) clearChangeList() {
	state.stateChangeList = newStateChangeList()
	state.validRevisions = state.validRevisions[:0]
	state.nextRevisionId = 0
}

//fixme
func (state *AccountStateDB) SetData(addr common.Address, key string, value []byte) (err error) {
	if state.IsEmptyAccount(addr) {
		return gerror.ErrAccountNotExist
	}

	log.DLogger.Debug("SetData", zap.Any("addr", addr), zap.String("key", key), zap.Uint8s("keybyte", []byte(key)), zap.Uint8s("value", value))
	var preValue []byte
	if state.smartContractData[addr] == nil {
		state.smartContractData[addr] = make(map[string][]byte)
	}
	preValue = state.smartContractData[addr][key]
	if value != nil {
		state.smartContractData[addr][key] = value
	} else {
		delete(state.smartContractData[addr], key)
		if len(state.smartContractData[addr]) == 0 {
			delete(state.smartContractData, addr)
		}
	}

	log.DLogger.Debug("SetData result", zap.String("key", key), zap.Uint8s("result", state.smartContractData[addr][key]))
	state.stateChangeList.append(dataChange{Account: &addr, Key: key, Prev: preValue, Current: value, ChangeType: DataChange})
	return
}

//fixme can not roll back
func (state *AccountStateDB) GetData(addr common.Address, key string) (data []byte) {
	if state.smartContractData[addr] != nil {
		log.DLogger.Debug("AccountStateDB#GetData", zap.Any("addr", addr))
		if state.smartContractData[addr][key] != nil {
			log.DLogger.Debug("GetData called", zap.String("key", key), zap.String("restring", string(state.smartContractData[addr][key])), zap.Uint8s("result", state.smartContractData[addr][key]))
			return state.smartContractData[addr][key]
		}
	}
	tier, err := state.getContractTrie(addr)
	if err != nil {
		return
	}
	data, err = tier.TryGet(GetContractFieldKey(addr, key))
	if err != nil {
		return
	}
	log.DLogger.Debug("GetData called", zap.String("key", key), zap.String("restring", string(data)), zap.Uint8s("result", data))
	return
}

func (state *AccountStateDB) AddLog(addedLog *model.Log) error {
	if addedLog == nil {
		return gerror.ErrAddedLogIsNil
	}

	log.DLogger.Info("AddLog Called")
	txHash := addedLog.TxHash
	old := state.GetLogs(txHash)
	current := append(old, addedLog)
	state.logs[txHash] = current
	//log.DLogger.Info("Log Added", "txHash", txHash, "logData", addedLog.Data)
	state.stateChangeList.append(logsChange{TxHash: &txHash, Prev: old, Current: current, ChangeType: LogsChange})
	return nil
}

func (state *AccountStateDB) GetLogs(txHash common.Hash) []*model.Log {
	return state.logs[txHash]
}

func (state *AccountStateDB) finalSmartData() error {
	for addr, data := range state.smartContractData {
		ct, err := state.putSmartDataToTrie(addr, data)
		if err != nil {
			return err
		}

		// You must commit trie to memory, and only use commit trie db in the commit.
		ch, err := ct.Commit(nil)
		if err != nil {
			return err
		}
		log.DLogger.Info("finaliseContractData update contract root", zap.String("contract addr", addr.Hex()), zap.String("root", ch.Hex()))

		if err := state.SetDataRoot(addr, ch); err != nil {
			// change blockStateTrie to origin pre hash？If you want, clear the finalised contract root. But it is best to discard the AccountStateDB directly after the error is reported.
			//state.resetThisStateDB()
			log.DLogger.Error("Commit update contract root failed", zap.Error(err))
			return err
		}
		state.finalisedContractRoot[addr] = ch
	}
	state.smartContractData = make(map[common.Address]map[string][]byte)
	return nil
}

// put contract data to trie
func (state *AccountStateDB) putSmartDataToTrie(addr common.Address, data map[string][]byte) (StateTrie, error) {
	log.DLogger.Info("put contract", zap.Any("addr", addr))
	ct, err := state.getContractTrie(addr)
	//check err first and return err if not find trie, otherwise there isn't this trie if th ct is nil
	if err != nil && !strings.Contains(err.Error(), "missing trie node") {
		log.DLogger.Warn("can't get contract trie from db", zap.Error(err))
		return nil, err
	}
	// todo ct is nil?
	for k, v := range data {
		log.DLogger.Debug("putContractDataToTrie", zap.String("k", k), zap.Uint8s("v", v), zap.String("pre state", state.preStateRoot.Hex()))
		if err := ct.TryUpdate(GetContractFieldKey(addr, k), v); err != nil {
			return nil, err
		}
	}
	return ct, nil
}
