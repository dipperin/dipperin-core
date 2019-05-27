package vm

import (
	"github.com/dipperin/dipperin-core/core/chain/state-processor"
	"math/big"
	"github.com/dipperin/dipperin-core/common"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/dipperin/dipperin-core/core/vm/model"
)

type storage struct {
	blockStateTrie state_processor.StateTrie
	contractLogs   map[common.Hash][]*model.Log
}

func NewStorage() *storage {
	db := ethdb.NewMemDatabase()
	tdb := state_processor.NewStateStorageWithCache(db)
	tr, err := tdb.OpenTrie(common.Hash{})
	if err != nil {
		panic(err)
	}
	logMap := make(map[common.Hash][]*model.Log)
	return &storage{
		blockStateTrie: tr,
		contractLogs:   logMap,
	}
}

func (s *storage) CreateAccount(a common.Address) {
	panic("implement me")
}

func (s *storage) SubBalance(a common.Address, b *big.Int) {
	panic("implement me")
}

func (s *storage) AddBalance(a common.Address, b *big.Int) {
	panic("implement me")
}

func (s *storage) GetBalance(common.Address) *big.Int {
	panic("implement me")
}

func (s *storage) GetNonce(common.Address) uint64 {
	panic("implement me")
}

func (s *storage) SetNonce(common.Address, uint64) {
	panic("implement me")
}

func (s *storage) GetCodeHash(common.Address) common.Hash {
	panic("implement me")
}

func (s *storage) GetCode(common.Address) []byte {
	panic("implement me")
}

func (s *storage) SetCode(common.Address, []byte) {
	panic("implement me")
}

func (s *storage) GetCodeSize(common.Address) int {
	panic("implement me")
}

func (s *storage) GetAbiHash(common.Address) common.Hash {
	panic("implement me")
}

func (s *storage) GetAbi(common.Address) []byte {
	panic("implement me")
}

func (s *storage) SetAbi(common.Address, []byte) {
	panic("implement me")
}

func (s *storage) AddRefund(uint64) {
	panic("implement me")
}

func (s *storage) SubRefund(uint64) {
	panic("implement me")
}

func (s *storage) GetRefund() uint64 {
	panic("implement me")
}

func (s *storage) GetCommittedState(common.Address, []byte) []byte {
	panic("implement me")
}

func (s *storage) GetState(addr common.Address, key []byte) []byte {
	log.Info("GetState Called", "contractAddr", addr.Bytes())

	key, err := rlp.EncodeToBytes(append(addr.Bytes(), key...))
	if err != nil {
		panic(err)
	}

	value, err := s.blockStateTrie.TryGet(key)
	if err != nil {
		panic(err)
	}

	log.Info("Get State", "key", string(key), "value", value)
	return value
}

func (s *storage) SetState(addr common.Address, key []byte, value []byte) {
	log.Info("SetState Called", "contractAddr", addr.Bytes())

	key, err := rlp.EncodeToBytes(append(addr.Bytes(), key...))
	if err != nil {
		panic(err)
	}

	err = s.blockStateTrie.TryUpdate(key, value)
	if err != nil {
		panic(err)
	}
	log.Info("State Saved", "key", string(key), "value", value)
}

func (s *storage) AddLog(addedLog *model.Log) {
	log.Info("AddLog Called")

	txHash := addedLog.TxHash
	contractLogs := s.GetLogs(txHash)
	addedLog.Index = uint(len(contractLogs) + 1)
	s.contractLogs[txHash] = append(contractLogs, addedLog)

	log.Info("Log Added", "txHash", txHash, "logs", s.contractLogs[txHash])
}

func (s *storage) GetLogs(txHash common.Hash) []*model.Log {
	return s.contractLogs[txHash]
}

func (s *storage) Suicide(common.Address) bool {
	panic("implement me")
}

func (s *storage) HasSuicided(common.Address) bool {
	panic("implement me")
}

func (s *storage) Exist(common.Address) bool {
	panic("implement me")
}

func (s *storage) Empty(common.Address) bool {
	panic("implement me")
}

func (s *storage) AddPreimage(common.Hash, []byte) {
	panic("implement me")
}

func (s *storage) ForEachStorage(common.Address, func(common.Hash, common.Hash) bool) {
	panic("implement me")
}

func (s *storage) TxHash() common.Hash {
	panic("implement me")
}

func (s *storage) TxIdx() uint32 {
	panic("implement me")
}

func (s *storage) RevertToSnapshot(int) {
	panic("implement me")
}

func (s *storage) Snapshot() int {
	panic("implement me")
}
