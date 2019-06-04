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
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
	"sort"
	"io"
)

// journalEntry is a modification entry in the state change journal that can be
// reverted on demand.
type StateChange interface {
	// revert undoes the changes introduced by this journal entry.
	revert(*AccountStateDB)

	// dirtied returns the Ethereum address modified by this journal entry.
	dirtied() *common.Address

	recover(*AccountStateDB)

	getType() int

	digest(sc StateChange) StateChange
}

// journal contains the list of state modifications applied since the last state
// commit. These are tracked to be able to be reverted in case of an execution
// exception or revertal request.
type StateChangeList struct {
	changes []StateChange          // Current changes tracked by the journal
	dirties map[common.Address]int // Dirty accounts and the number of changes
}

type StateChangeRLP struct {
	StateType   uint64
	StateChange []byte
}

func (scl *StateChangeList) DecodeRLP(s *rlp.Stream) (err error) {
	buf, err := s.Raw()
	if err != nil {
		return err
	}

	var states []StateChangeRLP
	err = rlp.DecodeBytes(buf, &states)
	if err != nil {
		return err
	}

	scl.dirties = make(map[common.Address]int)
	for _, state := range states {
		switch state.StateType {
		case NewAccountChange:
			var change newAccountChange
			rlp.DecodeBytes(state.StateChange, &change)
			scl.append(change)
		case BalanceChange:
			var change balanceChange
			rlp.DecodeBytes(state.StateChange, &change)
			scl.append(change)
		case NonceChange:
			var change nonceChange
			rlp.DecodeBytes(state.StateChange, &change)
			scl.append(change)
		case HashLockChange:
			var change hashLockChange
			rlp.DecodeBytes(state.StateChange, &change)
			scl.append(change)
		case TimeLockChange:
			var change timeLockChange
			rlp.DecodeBytes(state.StateChange, &change)
			scl.append(change)
		case DataRootChange:
			var change dataRootChange
			rlp.DecodeBytes(state.StateChange, &change)
			scl.append(change)
		case StakeChange:
			var change stakeChange
			rlp.DecodeBytes(state.StateChange, &change)
			scl.append(change)
		case CommitNumChange:
			var change commitNumChange
			rlp.DecodeBytes(state.StateChange, &change)
			scl.append(change)
		case VerifyNumChange:
			var change verifyNumChange
			rlp.DecodeBytes(state.StateChange, &change)
			scl.append(change)
		case PerformanceChange:
			var change performanceChange
			rlp.DecodeBytes(state.StateChange, &change)
			scl.append(change)
		case LastElectChange:
			var change lastElectChange
			rlp.DecodeBytes(state.StateChange, &change)
			scl.append(change)
		case DeleteAccountChange:
			var change deleteAccountChange
			rlp.DecodeBytes(state.StateChange, &change)
			scl.append(change)
		case AbiChange:
			var change abiChange
			rlp.DecodeBytes(state.StateChange, &change)
			scl.append(change)
		case CodeChange:
			var change codeChange
			rlp.DecodeBytes(state.StateChange,&change)
			scl.append(change)
		case DataChange:
			var change dataChange
			rlp.DecodeBytes(state.StateChange,&change)
			scl.append(change)
		default:
			panic("no type")
		}
	}
	return nil
}

func (scl *StateChangeList) EncodeRLP(w io.Writer) error {

	var sList []StateChangeRLP
	//var errr error
	for _, change := range scl.changes {
		s := StateChangeRLP{}
		s.StateType = uint64(change.getType())
		s.StateChange, _ = rlp.EncodeToBytes(change)
		sList = append(sList, s)
	}


	err := rlp.Encode(w, sList)
	if err != nil {
		return err
	}
	return nil
}

func (scl *StateChangeList) Len() int {
	return len(scl.changes)
}

func (scl *StateChangeList) Less(i, j int) bool {
	return scl.changes[i].getType() < scl.changes[j].getType()
}

func (scl *StateChangeList) Swap(i, j int) {
	scl.changes[i], scl.changes[j] = scl.changes[j], scl.changes[i]
}

// newStateChangeList create a new statechangelist.
func newStateChangeList() *StateChangeList {
	return &StateChangeList{
		dirties: make(map[common.Address]int),
	}
}

// append inserts a new modification entry to the end of the change journal.
func (scl *StateChangeList) append(change StateChange) {
	scl.changes = append(scl.changes, change)
	if addr := change.dirtied(); addr != nil {
		scl.dirties[*addr]++
	}
}

// revert undoes a batch of state change along with any reverted
// dirty handling too.
func (scl *StateChangeList) revert(statedb *AccountStateDB, snapshot int) {
	for i := len(scl.changes) - 1; i >= snapshot; i-- {
		// Undo the changes made by the operation
		scl.changes[i].revert(statedb)
		// Drop any dirty tracking induced by the change
		if addr := scl.changes[i].dirtied(); addr != nil {
			if scl.dirties[*addr]--; scl.dirties[*addr] == 0 {
				delete(scl.dirties, *addr)
			}
		}
	}

	scl.changes = scl.changes[:snapshot]
}

//digest combine all the same type change from same address.
func (scl *StateChangeList) digest() *StateChangeList {
	totalChange := make(map[common.Address][]StateChange)
	for _, change := range scl.changes {
		addr := *change.dirtied()
		if totalChange[addr] == nil {
			totalChange[addr] = []StateChange{}
		}
		totalChange[addr] = append(totalChange[addr], change)
	}

	newscl := newStateChangeList()
	for _, stateChangeSlice := range totalChange {
		// state changes with same address
		// map of `ChangeType` to `last stateChange`
		changes := make(map[int]StateChange)

		for _, change := range stateChangeSlice {
			// for every change in same address
			ChangeType := change.getType()

			if changes[ChangeType] == nil {
				// if nil, initialize with the first change state
				changes[ChangeType] = change
			} else {
				// if not, digest and update the previous state change
				changes[ChangeType] = change.digest(changes[ChangeType])
			}
		}
		for _, change := range changes {
			newscl.append(change)
		}
	}

	return newscl
}

//recover should used after digest,if not digest, deleteAccountChange may produce some err.
func (scl *StateChangeList) recover(statedb *AccountStateDB) {
	sort.Sort(scl)
	for _, change := range scl.changes {
		change.recover(statedb)
	}
}

// dirty explicitly sets an address to dirty, even if the change entries would
// otherwise suggest it as clean. use with caution.
func (scl *StateChangeList) dirty(addr common.Address) {
	scl.dirties[addr]++
}

// length returns the Current number of state change.
func (scl *StateChangeList) length() int {
	return len(scl.changes)
}

const (
	NewAccountChange   = iota
	BalanceChange
	NonceChange
	HashLockChange
	TimeLockChange
	// not use this
	//ContractRootChange
	DataRootChange
	StakeChange
	CommitNumChange
	VerifyNumChange
	PerformanceChange
	LastElectChange
	AbiChange
	CodeChange
	DataChange

	DeleteAccountChange
)

type (
	newAccountChange struct {
		Account    *common.Address
		ChangeType uint64
	}

	deleteAccountChange struct {
		Account    *common.Address
		ChangeType uint64
	}
	balanceChange struct {
		Account    *common.Address
		Prev       *big.Int
		Current    *big.Int
		ChangeType uint64
	}
	nonceChange struct {
		Account    *common.Address
		Prev       uint64
		Current    uint64
		ChangeType uint64
	}
	hashLockChange struct {
		Account    *common.Address
		Prev       common.Hash
		Current    common.Hash
		ChangeType uint64
	}
	timeLockChange struct {
		Account    *common.Address
		Prev       *big.Int
		Current    *big.Int
		ChangeType uint64
	}
	//contractRootChange struct {
	//	Account    *common.Address
	//	Prev       common.Hash
	//	Current    common.Hash
	//	ChangeType uint64
	//}
	dataRootChange struct {
		Account    *common.Address
		Prev       common.Hash
		Current    common.Hash
		ChangeType uint64
	}
	stakeChange struct {
		Account    *common.Address
		Prev       *big.Int
		Current    *big.Int
		ChangeType uint64
	}
	commitNumChange struct {
		Account    *common.Address
		Prev       uint64
		Current    uint64
		ChangeType uint64
	}
	verifyNumChange struct {
		Account    *common.Address
		Prev       uint64
		Current    uint64
		ChangeType uint64
	}
	performanceChange struct {
		Account    *common.Address
		Prev       uint64
		Current    uint64
		ChangeType uint64
	}
	lastElectChange struct {
		Account    *common.Address
		Prev       uint64
		Current    uint64
		ChangeType uint64
	}
	abiChange struct{
		Account *common.Address
		Prev []byte
		Current []byte
		ChangeType uint64
	}
	codeChange struct{
		Account *common.Address
		Prev []byte
		Current []byte
		ChangeType uint64
	}
	dataChange struct{
		Account *common.Address
		Key string
		Prev []byte
		Current []byte
		ChangeType uint64
	}
)

func (sc deleteAccountChange) revert(s *AccountStateDB) {
	s.newAccountState(*sc.Account)
}

func (sc deleteAccountChange) dirtied() *common.Address {
	return sc.Account
}

func (sc deleteAccountChange) recover(s *AccountStateDB) {
	s.deleteAccountState(*sc.Account)
}

func (sc deleteAccountChange) getType() int {
	return int(sc.ChangeType)
}

func (sc deleteAccountChange) digest(change StateChange) StateChange {
	panic("delete same account twice ")
}

func (sc newAccountChange) revert(s *AccountStateDB) {
	s.deleteAccountState(*sc.Account)
}

func (sc newAccountChange) recover(s *AccountStateDB) {
	s.newAccountState(*sc.Account)
}

func (sc newAccountChange) dirtied() *common.Address {
	return sc.Account
}
func (sc newAccountChange) getType() int {
	return int(sc.ChangeType)
}

func (sc newAccountChange) digest(change StateChange) StateChange {
	panic("creat same account twice ")
}





func (sc lastElectChange) revert(s *AccountStateDB) {
	s.setLastElect(*sc.Account, sc.Prev)
}

func (sc lastElectChange) recover(s *AccountStateDB) {
	s.setLastElect(*sc.Account, sc.Current)
}

func (sc lastElectChange) dirtied() *common.Address {
	return sc.Account
}

func (sc lastElectChange) getType() int {
	return int(sc.ChangeType)
}

func (sc lastElectChange) digest(change StateChange) StateChange {
	if change.getType() == LastElectChange {
		c := change.(lastElectChange)
		return lastElectChange{Account: sc.Account, Prev: c.Prev, Current: sc.Current, ChangeType: LastElectChange}
	}
	return nil
}






func (sc verifyNumChange) revert(s *AccountStateDB) {
	s.setVerifyNum(*sc.Account, sc.Prev)
}

func (sc verifyNumChange) recover(s *AccountStateDB) {
	s.setVerifyNum(*sc.Account, sc.Current)
}

func (sc verifyNumChange) dirtied() *common.Address {
	return sc.Account
}

func (sc verifyNumChange) getType() int {
	return int(sc.ChangeType)
}
func (sc verifyNumChange) digest(change StateChange) StateChange {
	if change.getType() == VerifyNumChange {
		c := change.(verifyNumChange)
		return verifyNumChange{Account: sc.Account, Prev: c.Prev, Current: sc.Current, ChangeType: VerifyNumChange}
	}
	return nil
}





func (sc commitNumChange) revert(s *AccountStateDB) {
	s.setCommitNum(*sc.Account, sc.Prev)
}
func (sc commitNumChange) recover(s *AccountStateDB) {
	s.setCommitNum(*sc.Account, sc.Current)
}

func (sc commitNumChange) dirtied() *common.Address {
	return sc.Account
}

func (sc commitNumChange) getType() int {
	return int(sc.ChangeType)
}
func (sc commitNumChange) digest(change StateChange) StateChange {
	if change.getType() == CommitNumChange {
		c := change.(commitNumChange)
		return commitNumChange{Account: sc.Account, Prev: c.Prev, Current: sc.Current, ChangeType: CommitNumChange}
	}
	return nil
}




func (sc performanceChange) revert(s *AccountStateDB) {
	s.setPerformance(*sc.Account, sc.Prev)
}

func (sc performanceChange) recover(s *AccountStateDB) {
	s.setPerformance(*sc.Account, sc.Current)
}

func (sc performanceChange) dirtied() *common.Address {
	return sc.Account
}

func (sc performanceChange) getType() int {
	return int(sc.ChangeType)
}
func (sc performanceChange) digest(change StateChange) StateChange {
	if change.getType() == PerformanceChange {
		c := change.(performanceChange)
		return performanceChange{Account: sc.Account, Prev: c.Prev, Current: sc.Current, ChangeType: PerformanceChange}
	}
	return nil
}




func (sc stakeChange) revert(s *AccountStateDB) {
	s.setStake(*sc.Account, sc.Prev)
}

func (sc stakeChange) recover(s *AccountStateDB) {
	s.setStake(*sc.Account, sc.Current)
}

func (sc stakeChange) dirtied() *common.Address {
	return sc.Account
}

func (sc stakeChange) getType() int {
	return int(sc.ChangeType)
}
func (sc stakeChange) digest(change StateChange) StateChange {
	if change.getType() == StakeChange {
		c := change.(stakeChange)
		return stakeChange{Account: sc.Account, Prev: c.Prev, Current: sc.Current, ChangeType: StakeChange}
	}
	return nil
}





func (sc dataRootChange) revert(s *AccountStateDB) {
	s.setDataRoot(*sc.Account, sc.Prev)
}

func (sc dataRootChange) recover(s *AccountStateDB) {
	s.setDataRoot(*sc.Account, sc.Current)
}

func (sc dataRootChange) dirtied() *common.Address {
	return sc.Account
}

func (sc dataRootChange) getType() int {
	return int(sc.ChangeType)
}
func (sc dataRootChange) digest(change StateChange) StateChange {
	if change.getType() == DataRootChange {
		c := change.(dataRootChange)
		return dataRootChange{Account: sc.Account, Prev: c.Prev, Current: sc.Current, ChangeType: DataRootChange}
	}
	return nil
}


func (sc timeLockChange) revert(s *AccountStateDB) {
	s.setTimeLock(*sc.Account, sc.Prev)
}

func (sc timeLockChange) recover(s *AccountStateDB) {
	s.setTimeLock(*sc.Account, sc.Current)
}

func (sc timeLockChange) dirtied() *common.Address {
	return sc.Account
}

func (sc timeLockChange) getType() int {
	return int(sc.ChangeType)
}
func (sc timeLockChange) digest(change StateChange) StateChange {
	if change.getType() == TimeLockChange {
		c := change.(timeLockChange)
		return timeLockChange{Account: sc.Account, Prev: c.Prev, Current: sc.Current, ChangeType: TimeLockChange}
	}
	return nil
}





func (sc hashLockChange) revert(s *AccountStateDB) {
	s.setHashLock(*sc.Account, sc.Prev)
}

func (sc hashLockChange) recover(s *AccountStateDB) {
	s.setHashLock(*sc.Account, sc.Current)
}

func (sc hashLockChange) dirtied() *common.Address {
	return sc.Account
}

func (sc hashLockChange) getType() int {
	return int(sc.ChangeType)
}
func (sc hashLockChange) digest(change StateChange) StateChange {
	if change.getType() == HashLockChange {
		c := change.(hashLockChange)
		return hashLockChange{Account: sc.Account, Prev: c.Prev, Current: sc.Current, ChangeType: HashLockChange}
	}
	return nil
}





func (sc balanceChange) revert(s *AccountStateDB) {
	s.setBalance(*sc.Account, sc.Prev)
}

func (sc balanceChange) recover(s *AccountStateDB) {
	s.setBalance(*sc.Account, sc.Current)
}

func (sc balanceChange) dirtied() *common.Address {
	return sc.Account
}

func (sc balanceChange) getType() int {
	return int(sc.ChangeType)
}
func (sc balanceChange) digest(change StateChange) StateChange {
	if change.getType() == BalanceChange {
		c := change.(balanceChange)
		return balanceChange{Account: sc.Account, Prev: c.Prev, Current: sc.Current, ChangeType: BalanceChange}
	}
	return nil
}





func (sc nonceChange) revert(s *AccountStateDB) {
	s.setNonce(*sc.Account, sc.Prev)
}

func (sc nonceChange) recover(s *AccountStateDB) {
	s.setNonce(*sc.Account, sc.Current)
}

func (sc nonceChange) dirtied() *common.Address {
	return sc.Account
}

func (sc nonceChange) getType() int {
	return int(sc.ChangeType)
}
func (sc nonceChange) digest(change StateChange) StateChange {
	if change.getType() == NonceChange {
		c := change.(nonceChange)
		return nonceChange{Account: sc.Account, Prev: c.Prev, Current: sc.Current, ChangeType: NonceChange}
	}
	return nil
}


func (sc abiChange) revert(s *AccountStateDB) {
	s.setAbi(*sc.Account, sc.Prev)
}

func (sc abiChange) recover(s *AccountStateDB) {
	s.setAbi(*sc.Account, sc.Current)
}

func (sc abiChange) dirtied() *common.Address {
	return sc.Account
}

func (sc abiChange) getType() int {
	return int(sc.ChangeType)
}
func (sc abiChange) digest(change StateChange) StateChange {
	if change.getType() == AbiChange {
		c := change.(abiChange)
		return abiChange{Account: sc.Account, Prev: c.Prev, Current: sc.Current, ChangeType: AbiChange}
	}
	return nil
}


func (sc codeChange) revert(s *AccountStateDB) {
	s.setCode(*sc.Account, sc.Prev)
}

func (sc codeChange) recover(s *AccountStateDB) {
	s.setCode(*sc.Account, sc.Current)
}

func (sc codeChange) dirtied() *common.Address {
	return sc.Account
}

func (sc codeChange) getType() int {
	return int(sc.ChangeType)
}
func (sc codeChange) digest(change StateChange) StateChange {
	if change.getType() == CodeChange {
		c := change.(codeChange)
		return codeChange{Account: sc.Account, Prev: c.Prev, Current: sc.Current, ChangeType: CodeChange}
	}
	return nil
}


func (sc dataChange) revert(s *AccountStateDB) {
	s.SetData(*sc.Account,sc.Key,sc.Prev)
}

func (sc dataChange) recover(s *AccountStateDB) {
	s.SetData(*sc.Account,sc.Key,sc.Current)
}

func (sc dataChange) dirtied() *common.Address {
	return sc.Account
}

func (sc dataChange) getType() int {
	return int(sc.ChangeType)
}
func (sc dataChange) digest(change StateChange) StateChange {
	if change.getType() == DataChange {
		c := change.(dataChange)
		return dataChange{Account: sc.Account,Key:c.Key, Prev: c.Prev, Current: sc.Current, ChangeType: DataChange}
	}
	return nil
}