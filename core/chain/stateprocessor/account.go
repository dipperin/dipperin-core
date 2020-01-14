package stateprocessor

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
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
	Code     []byte
	Abi      []byte
}

const (
	nonceKeySuffix    = "_nonce"
	balanceKeySuffix  = "_balance"
	hashLockKeySuffix = "_hashLock"
	timeLockKeySuffix = "_timeLock"
	//contractRootSuffix = "_contract_root"
	dataRootSuffix     = "_data_root"
	stakeKeySuffix     = "_stake"
	commitNumKeySuffix = "_commit_num"
	verifyNumKeySuffix = "_verify_num"
	lastElectKeySuffix = "_last_elect"
	performanceSuffix  = "_performance"
	abiSuffix          = "_abi"
	codeSuffix         = "_code"
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

func GetAbiKey(address common.Address) []byte {
	return append(address[:], []byte(abiSuffix)...)
}

func GetCodeKey(address common.Address) []byte {
	return append(address[:], []byte(codeSuffix)...)
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

func (a *account) DataRootBytes() []byte {
	v, _ := rlp.EncodeToBytes(a.DataRoot)
	return v
}
