package cs_crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/binary"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/third_party/crypto"
	"github.com/dipperin/dipperin-core/third_party/crypto/secp256k1"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"
	"math/big"
)

// ValidateSignatureValues verifies whether the signature values are valid with
// the given chain rules. The v value is assumed to be either 0 or 1.
func ValidateSignatureValues(v byte, r, s *big.Int, homestead bool) bool {
	if r.Cmp(common.Big1) < 0 || s.Cmp(common.Big1) < 0 {
		return false
	}
	// reject upper range of s values (ECDSA malleability)
	// see discussion in secp256k1/libsecp256k1/include/secp256k1.h
	if homestead && s.Cmp(crypto.Secp256k1halfN) > 0 {
		return false
	}
	// Frontier: allow s to be in full N range
	return r.Cmp(crypto.Secp256k1N) < 0 && s.Cmp(crypto.Secp256k1N) < 0 && (v == 0 || v == 1)
}

// don't use this method, use GetNormalAddress
func PubkeyToAddress(p ecdsa.PublicKey) common.Address {
	pubBytes := crypto.FromECDSAPub(&p)
	return common.BytesToAddress(crypto.Keccak256(pubBytes[1:])[12:])
}

// CreateAddress creates an ethereum address given the bytes and the nonce
func CreateAddress(b common.Address, nonce uint64) common.Address {
	data, _ := rlp.EncodeToBytes([]interface{}{b, nonce})
	return common.BytesToAddress(crypto.Keccak256(data)[12:])
}

// CreateAddress2 creates an ethereum address given the address bytes, initial
// contract code hash and a salt.
func CreateAddress2(b common.Address, salt [32]byte, inithash []byte) common.Address {
	return common.BytesToAddress(crypto.Keccak256([]byte{0xff}, b.Bytes(), salt[:], inithash)[12:])
}

// Keccak256Hash calculates and returns the Keccak256 hash of the input data,
// converting it to an internal Hash data structure.
func Keccak256Hash(data ...[]byte) (h common.Hash) {
	d := sha3.NewLegacyKeccak256()
	for _, b := range data {
		d.Write(b)
	}
	d.Sum(h[:0])
	return h
}

func ToECDSAPub(pub []byte) *ecdsa.PublicKey {
	if len(pub) == 0 {
		return nil
	}
	x, y := elliptic.Unmarshal(crypto.S256(), pub)
	return &ecdsa.PublicKey{Curve: crypto.S256(), X: x, Y: y}
}

//检查签名的r，s，v是否越界
//todo 加测试
func ValidSigValue(r, s, v *big.Int) bool {
	if r.Cmp(common.Big1) < 0 || s.Cmp(common.Big1) < 0 {
		return false
	}
	// reject upper range of s values (ECDSA malleability)
	// see discussion in secp256k1/libsecp256k1/include/secp256k1.h
	if s.Cmp(crypto.Secp256k1halfN) > 0 {
		return false
	}
	return r.Cmp(crypto.Secp256k1N) < 0 && s.Cmp(crypto.Secp256k1N) < 0 && (v.Cmp(common.Big1) == 0 || v.Cmp(common.Big0) == 0)
}

func GetNormalAddress(p ecdsa.PublicKey) common.Address {
	pubBytes := crypto.FromECDSAPub(&p)
	var tmpTypeB [2]byte
	binary.BigEndian.PutUint16(tmpTypeB[:], uint16(common.AddressTypeNormal))
	tmpAddr := crypto.Keccak256(pubBytes[1:])[12:]
	return common.BytesToAddress(append(tmpTypeB[:], tmpAddr...))
}

func GetLockAddress(alice, bob common.Address) common.Address {
	res, err := rlp.EncodeToBytes([]interface{}{alice, bob})
	if err != nil {
		return common.Address{}
	}
	var tmpTypeB [2]byte
	binary.BigEndian.PutUint16(tmpTypeB[:], uint16(common.AddressTypeCross))
	tmpAddr := crypto.Keccak256(res[:])[12:]
	return common.BytesToAddress(append(tmpTypeB[:], tmpAddr...))
}

func GetEvidenceAddress(target common.Address) common.Address {
	var tmpType [2]byte
	binary.BigEndian.PutUint16(tmpType[:], uint16(common.AddressTypeEvidence))
	var trueAdd = target[2:]
	var evAdd []byte
	evAdd = append(evAdd, tmpType[0], tmpType[1])
	evAdd = append(evAdd, trueAdd...)
	return common.BytesToAddress(evAdd)
}

func GetNormalAddressFromEvidence(target common.Address) common.Address {
	var tmpType [2]byte
	binary.BigEndian.PutUint16(tmpType[:], uint16(common.AddressTypeNormal))
	var trueAdd = target[2:]
	var evAdd []byte
	evAdd = append(evAdd, tmpType[0], tmpType[1])
	evAdd = append(evAdd, trueAdd...)
	return common.BytesToAddress(evAdd)
}

func GetContractAddress(address common.Address) common.Address {
	var tmpType [2]byte
	binary.BigEndian.PutUint16(tmpType[:], uint16(common.AddressTypeERC20))

	var trueAdd = address[2:]
	var cAdd []byte
	cAdd = append(cAdd, tmpType[0], tmpType[1])
	cAdd = append(cAdd, trueAdd...)
	return common.BytesToAddress(cAdd)
}

func CreateContractAddress(b common.Address, nonce uint64) common.Address {
	var tmpTypeB [2]byte
	binary.BigEndian.PutUint16(tmpTypeB[:], uint16(common.AddressTypeContractCall))
	data, _ := rlp.EncodeToBytes([]interface{}{b, nonce})
	tempAddr := crypto.Keccak256(data)[12:]
	addr := append(tmpTypeB[:], tempAddr...)
	return common.BytesToAddress(addr)
}

//func zeroBytes(bytes []byte) {
//	for i := range bytes {
//		bytes[i] = 0
//	}
//}

func RecoverAddressFromSig(hash common.Hash, sig []byte) (common.Address, error) {
	recoverPk, err := secp256k1.RecoverPubkey(hash[:], sig)
	if err != nil {
		return common.Address{}, err
	}
	return GetNormalAddress(*ToECDSAPub(recoverPk)), nil
}
