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


package common

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/dipperin/dipperin-core/common/consts"
	"github.com/dipperin/dipperin-core/common/hexutil"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
	"strconv"
)

const (
	HashLength = 32
	// The first two bytes are used to identify the type of the address, and 0x0001 is the normal address.
	AddressLength = 22
	NonceLength   = 32
	DiffLength    = 4
)

type TxType int
// address type
const (
	AddressTypeNormal   = 0x0000
	AddressTypeCross    = 0x0001
	AddressTypeStake    = 0x0002
	AddressTypeCancel   = 0x0003
	AddressTypeUnStake  = 0x0004
	AddressTypeEvidence = 0x0005
	AddressTypeERC20    = 0x0010
	AddressTypeEarlyReward    = 0x0011
	AddressTypeSmartContract = 0x0016
)

func (txType TxType) String() string {
	switch txType {
	case AddressTypeNormal:
		return "normal transaction"
	case AddressTypeCross:
		return "cross chain transaction"
	case AddressTypeStake:
		return "stake transaction"
	case AddressTypeCancel:
		return "cancel transaction"
	case AddressTypeUnStake:
		return "unstake transaction"
	case AddressTypeEvidence:
		return "evidence transaction"
	case AddressTypeERC20:
		return "erc20 transaction"
	default:
		return fmt.Sprintf("unkonw tx:%v", int(txType))
	}
}

const (
	AddressStake   = "0x00020000000000000000000000000000000000000000"
	AddressCancel  = "0x00030000000000000000000000000000000000000000"
	AddressUnStake = "0x00040000000000000000000000000000000000000000"
)

// Dipperin hash
type Hash [HashLength]byte

func (h Hash) String() string {
	return h.Hex()
}

// Returns an exact copy of the provided bytes
func CopyHash(h *Hash) (*Hash) {
	copied:=Hash{}
	if len(h)==0{
		return &copied
	}
	copy(copied[:], h[:])
	return &copied
}

func (h Hash) IsEqual(oh Hash) bool {
	return bytes.Equal(h[:], oh[:])
}
func (h Hash) IsEmpty() bool {
	emptyHash := Hash{}
	return bytes.Equal(h[:], emptyHash[:])
}

func (h Hash) MarshalJSON() ([]byte, error) {
	hb := (hexutil.Bytes)(h[:])
	return util.StringifyJsonToBytesWithErr(hb)
}

func (h *Hash) UnmarshalJSON(input []byte) error {
	var hb hexutil.Bytes
	if err := hb.UnmarshalJSON(input); err != nil {
		return err
	}
	copy((*h)[:], hb)
	return nil
}

// compare hash h with hash oh, if h bigger than oh return 1, if h equal to oh return 0,if h smaller than oh return -1
func (h Hash) Cmp(oh Hash) (res int) {
	for i := 0; i < HashLength-1; i++ {
		if h[i] != oh[i] {
			if h[i] > oh[i] {
				res = 1
			} else {
				res = -1
			}
			break
		}
	}
	return
}

func (h *Hash) Clear() {
	*h = Hash{}
}

// do rlp hash on an object
func RlpHashKeccak256(v interface{}) (h Hash) {
	if b, err := rlp.EncodeToBytes(v); err != nil {
		return
	} else {
		return BytesToHash(crypto.Keccak256(b))
	}
}

func BytesToHash(b []byte) (result Hash) {
	result.SetBytes(b)
	return
}

// Sets the hash to the value of b. If b is larger than len(h), 'b' will be cropped (from the left).
func (h *Hash) SetBytes(b []byte) {
	if len(b) > len(h) {
		b = b[len(b)-HashLength:]
	}

	copy(h[HashLength-len(b):], b)
}

func (h Hash) Hex() string          { return hexutil.Encode(h[:]) }
func (h Hash) HexWithout0x() string { return hexutil.EncodeWithout0x(h[:]) }
func (h Hash) Str() string          { return string(h[:]) }
func (h Hash) Bytes() []byte        { return h[:] }
func (h Hash) Big() *big.Int        { return new(big.Int).SetBytes(h[:]) }

func HexToHash(s string) Hash    { return BytesToHash(FromHex(s)) }
func StringToHash(s string) Hash { return BytesToHash([]byte(s)) }
func BigToHash(b *big.Int) Hash  { return BytesToHash(b.Bytes()) }

//func (h Hash) TerminalString() string {
//	return fmt.Sprintf("%x…%x", h[:3], h[29:])
//	return ""
//}

//implement stringer interface
//func (h Hash) String() string {
//	return h.Hex()
//}

//func (h Hash) Format(s fmt.State, c rune) {
//	fmt.Fprintf(s, "%"+string(c), h[:])
//}

// ValidHashForDifficulty check is it a hash that meets the corresponding difficulty?
func (h Hash) ValidHashForDifficulty(difficulty Difficulty) bool {
	//log.Debug("h ValidHashForDifficulty", "hash", h.Hex(), "diff", difficulty.Hex())
	result := h.Cmp(difficulty.DiffToTarget())
	if result <= 0 {
		return true
	} else {
		//log.Info("===============received-hash==============",h.Hex())
		//log.Info("===============difftohash==============",difficulty.DiffToTarget().Hex())
		return false
	}
	//prefix := strings.Repeat("0", difficulty)
	//return strings.HasPrefix(h.HexWithout0x(), prefix)
}

type UnprefixedAddress Address

// Dipperin Address
type Address [AddressLength]byte

func (addr Address) String() string {
	return addr.Hex()
}

func (addr Address) MarshalJSON() ([]byte, error) {
	hb := (hexutil.Bytes)(addr[:])
	return util.StringifyJsonToBytesWithErr(hb)
}

func (addr *Address) UnmarshalJSON(input []byte) error {
	var hb hexutil.Bytes
	if err := hb.UnmarshalJSON(input); err != nil {
		return err
	}
	copy((*addr)[:], hb)
	return nil
}

// GetAddressType get the type of address based on the first two digits of the address
func (addr Address) GetAddressType() TxType {
	//cslog.Debug().Str("addr", addr.Hex()).Msg("GetAddressType")
	//if len(addr) < 2 {
	//	//cslog.Warn().Msg("GetAddressType but the length of the address is less than 2")
	//	return AddressTypeNormal
	//}
	return TxType(binary.BigEndian.Uint16(addr[:2]))
}

// GetAddressTypeStr get a text description of the address type
func (addr Address) GetAddressTypeStr() string {
	switch addr.GetAddressType() {
	case AddressTypeNormal:
		return "Normal"
	case AddressTypeERC20:
		return consts.ERC20TypeName
	case AddressTypeCross:
		return "CrossChain"
	case AddressTypeStake:
		return "Stake"
	case AddressTypeCancel:
		return "Cancel"
	case AddressTypeUnStake:
		return "UnStake"
	case AddressTypeEvidence:
		return "Evidence"
	case AddressTypeEarlyReward:
		return consts.EarlyTokenTypeName
	}
	return "UnKnown"
}

func (addr Address) IsEmpty() bool {
	emptyAddress := Address{}
	if bytes.Equal(addr[:], emptyAddress[:]) {
		return true
	}
	return false
}

func (addr Address) IsEqual(oaddr Address) bool {
	return bytes.Equal(addr[:], oaddr[:])
}

func (addr Address) IsEqualWithoutType(oaddr Address) bool {
	return bytes.Equal(addr[2:], oaddr[2:])
}

func (addr *Address) Clear() {
	*addr = Address{}
}

func BytesToAddress(b []byte) (result Address) {
	result.SetBytes(b)
	return
}

// Sets the address to the value of b. If b is larger than len(a) it will panic
func (addr *Address) SetBytes(b []byte) {
	if len(b) > len(addr) {
		b = b[len(b)-AddressLength:]
	}
	copy(addr[AddressLength-len(b):], b)
}

// Get the string representation of the underlying address
func (addr Address) Str() string   { return string(addr[:]) }
func (addr Address) Bytes() []byte { return addr[:] }
func (addr Address) Big() *big.Int { return new(big.Int).SetBytes(addr[:]) }
func (addr Address) Hash() Hash    { return BytesToHash(addr[:]) }
func (addr Address) InSlice(addresses []Address) bool {
	for _, a := range addresses {
		if addr.IsEqual(a) {
			return true
		}
	}
	return false
}

func StringToAddress(s string) Address { return BytesToAddress([]byte(s)) }
func BigToAddress(b *big.Int) Address  { return BytesToAddress(b.Bytes()) }
func HexToAddress(s string) Address    { return BytesToAddress(FromHex(s)) }

func (addr Address) Hex() string {
	unchecksummed := hex.EncodeToString(addr[:])
	hash := crypto.Keccak256([]byte(unchecksummed))

	result := []byte(unchecksummed)
	for i := 0; i < len(result); i++ {
		hashByte := hash[i/2]
		if i%2 == 0 {
			hashByte = hashByte >> 4
		} else {
			hashByte &= 0xf
		}
		if result[i] > '9' && hashByte > 7 {
			result[i] -= 32
		}
	}
	return "0x" + string(result)
}

//Dipperin BlockHeader Difficulty
type Difficulty [DiffLength]byte

func (d Difficulty) MarshalJSON() ([]byte, error) {
	hb := (hexutil.Bytes)(d[:])
	return util.StringifyJsonToBytesWithErr(hb)
}
func (d *Difficulty) UnmarshalJSON(input []byte) error {
	var hb hexutil.Bytes
	if err := hb.UnmarshalJSON(input); err != nil {
		return err
	}
	copy((*d)[:], hb)
	return nil
}

func BytesToDiff(b []byte) (result Difficulty) {
	result.SetBytes(b)
	return
}

func (d Difficulty) Equal(d2 Difficulty) bool {
	return bytes.Equal(d[:], d2[:])
}

func (d *Difficulty) SetBytes(b []byte) {
	if len(b) > len(d) {
		b = b[len(b)-DiffLength:]
	}
	copy(d[DiffLength-len(b):], b)
}

func (d Difficulty) HexWithout0x() string { return hexutil.EncodeWithout0x(d[:]) }
func (d Difficulty) Str() string          { return string(d[:]) }
func (d Difficulty) Bytes() []byte        { return d[:] }

func HexToDiff(s string) Difficulty    { return BytesToDiff(FromHex(s)) }
func StringToDiff(s string) Difficulty { return BytesToDiff([]byte(s)) }

func (d Difficulty) Hex() string {
	return hexutil.Encode(d[:])
}

func (d Difficulty) DiffToTarget() (target Hash) {
	a := HashLength - d[0]
	//cslog.Debug().Interface("a", a).Interface("d", d).Msg("DiffToTarget")
	if a + 2 > HashLength - 1 || HashLength < d[0] {
		log.Error("DiffToTarget failed", "diff", d.Hex())
		panic("The first digit of diff cannot be less than 3 and cannot be greater than 0x20")
	}
	target[a] = d[1]
	target[a+1] = d[2]
	target[a+2] = d[3]
	return
}

func (d Difficulty) Big() *big.Int { return d.DiffToTarget().Big() }

func BigToDiff(b *big.Int) Difficulty {
	tempint := bigToCompact(b)
	temphex := strconv.FormatInt(int64(tempint), 16)
	return HexToDiff(temphex)
}

func bigToCompact(n *big.Int) uint32 {
	// No need to do any work if it's zero.
	if n.Sign() == 0 {
		return 0
	}

	// Since the base for the exponent is 256, the exponent can be treated
	// as the number of bytes.  So, shift the number right or left
	// accordingly.  This is equivalent to:
	// mantissa = mantissa / 256^(exponent-3)
	var mantissa uint32
	exponent := uint(len(n.Bytes()))
	if exponent <= 3 {
		mantissa = uint32(n.Bits()[0])
		mantissa <<= 8 * (3 - exponent)
	} else {
		// Use a copy to avoid modifying the caller's original number.
		tn := new(big.Int).Set(n)
		mantissa = uint32(tn.Rsh(tn, 8*(exponent-3)).Bits()[0])
	}

	// When the mantissa already has the sign bit set, the number is too
	// large to fit into the available 23-bits, so divide the number by 256
	// and increment the exponent accordingly.
	if mantissa&0x00800000 != 0 {
		mantissa >>= 8
		exponent++
	}

	// Pack the exponent, sign bit, and mantissa into an unsigned 32-bit
	// int and return it.
	compact := uint32(exponent<<24) | mantissa
	if n.Sign() < 0 {
		compact |= 0x00800000
	}
	return compact
}

// Dipperin Block Nonce
type BlockNonce [NonceLength]byte

func (bn BlockNonce) MarshalJSON() ([]byte, error) {
	hb := (hexutil.Bytes)(bn[:])
	return util.StringifyJsonToBytesWithErr(hb)
}
func (bn *BlockNonce) UnmarshalJSON(input []byte) error {
	var hb hexutil.Bytes
	if err := hb.UnmarshalJSON(input); err != nil {
		return err
	}
	copy((*bn)[:], hb)
	return nil
}

//Compare block nonce
func (bn BlockNonce) IsEqual(obn BlockNonce) bool{
	return bytes.Equal(bn[:], obn[:])
}

// calculate hex
func (bn BlockNonce) Hex() string {
	return hexutil.Encode(bn[:])
}

func BlockNonceFromHex(hexStr string) (result BlockNonce) {
	rb, _ := hexutil.Decode(hexStr)
	copy(result[:], rb)
	return
}

func BlockNonceFromInt(i uint32) (result BlockNonce) {
	binary.BigEndian.PutUint32(result[:], i)
	return
}

// EncodeNonce converts the given integer to a block nonce.
func EncodeNonce(i uint64) BlockNonce {
	var n BlockNonce
	binary.BigEndian.PutUint32(n[:], uint32(i))
	return n
}

// CsBigInt support for parsing large numbers of strings, RLP encoding does not support it
//type CsBigInt struct {
//	*big.Int
//}
//func (csBitInt CsBigInt) MarshalJSON() ([]byte, error) {
//	tmpS := hexutil.EncodeBig(csBitInt.Int)
//	resultB := []byte(`"` + tmpS + `"`)
//	return resultB, nil
//}
//
//func (csBitInt *CsBigInt) UnmarshalJSON(input []byte) error {
//	tmpStr := strings.Replace(string(input), `"`, "", -1)
//	x, err := hexutil.DecodeBig(tmpStr)
//	if err != nil {
//		return err
//	}
//	csBitInt.Int = x
//	return nil
//}
//func NewCsBigInt(x int64) *CsBigInt {
//	return &CsBigInt{ Int: big.NewInt(x) }
//}