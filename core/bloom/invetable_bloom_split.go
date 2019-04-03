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


package iblt

import (
	"bytes"
	"encoding/binary"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/ethereum/go-ethereum/rlp"
	"reflect"
	"sort"
)

// InsertRLP inserts any k, v to IBLT using RLP coding
// it prints error if input is invalid or encode throws error
func (b *InvBloom) InsertRLP(k, v interface{}) {
	if k == nil {
		log.Error("insert key is nil", "key", k)
		return
	}

	if v == nil {
		log.Error("insert value is nil", "value", v)
		return
	}

	kBytes, kError := rlp.EncodeToBytes(k)
	vBytes, vError := rlp.EncodeToBytes(v)

	if kError != nil {
		log.Error("key RLP encode error", "error", kError, "key", k)
		return
	}

	if vError != nil {
		log.Error("value RLP encode error", "error", vError, "value", v)
		return
	}

	b.insertRLP(kBytes, vBytes)
}

// insertRLP does the actual insert to IBLT
// then digests the classical bloom filter
func (b *InvBloom) insertRLP(k, v []byte) {
	m := b.constructKeyValuePairs(k, v)
	b.insertMap(m)
}

// key is hashed with salt added to form a new key
// to new key is taken from the first BucketKeyHashLength bytes of salted and hashed key
// last BucketSerNumLength bytes are used to mark the index of byte slices.
// BucketValueLength defines how long the sub-array is.
// keyIdx ranges from [0, ceil(len(v)/BucketValueLength)]
// Data = [key, serNum, val]
// keys are usually the hash of val in most cases
// serNum ranges from 0x00... up to 0xFF..., it is set by the config
// values are the actual encoding parts
// This Data is then hashed to form `DataHash`
func (b *InvBloom) constructKeyValuePairs(k, v []byte) map[common.Hash]Data {
	res := make(map[common.Hash]Data)

	// append salt
	hash := hash(append(k, b.salt))

	// TODO, add overflow safety check
	for byteIdx, keyIdx := 0, uint16(0); ; keyIdx++ {
		// do-while byteIdx+BucketValueLength exceeds len(v)
		key := b.NewData()

		copy(key[0:b.config.KeyLen], hash)
		binary.BigEndian.PutUint16(key[b.config.KeyLen:b.config.KeyLen+b.config.SerNumLen], keyIdx)

		if byteIdx+int(b.config.ValLen) < len(v) {
			// normal loops, inserts fix length bytes
			copy(key[b.config.KeyLen+b.config.SerNumLen:], v[byteIdx:byteIdx+int(b.config.ValLen)])
			res[key.Hash()] = b.NewData()
			copy(res[key.Hash()], key)
		} else {
			// last loop, inserts all the remaining bytes
			copy(key[b.config.KeyLen+b.config.SerNumLen:], v[byteIdx:])
			res[key.Hash()] = b.NewData()
			copy(res[key.Hash()], key)
			break
		}

		byteIdx += int(b.config.ValLen)
	}
	return res
}

// mapToSlice converts the keys and values of a map to two key value slices
func mapToSlice(m map[common.Hash]Data) (key []Data) {
	for _, v := range m {
		key = append(key, v)
	}

	return
}

// b is already the subtraction of two IBLTs,
// if the decode succeeds,
// reconstruct the two original sets
func (b *InvBloom) ListRLP() (Alice, Bob [][]byte, err error) {
	alice := make(map[common.Hash]Data)
	bob := make(map[common.Hash]Data)

	// This method gets the unique elements of alice and bob
	success := b.Decode(alice, bob)

	if !success {
		return nil, nil, ErrDecodeFailed
	}

	Alice = b.reconstructByteSlice(alice)

	Bob = b.reconstructByteSlice(bob)
	return Alice, Bob, nil
}

// reconstructByteSlice returns valid slices by the ordered hash key in m.
func (b *InvBloom) reconstructByteSlice(m map[common.Hash]Data) (byteSlice [][]byte) {
	keySlice := mapToSlice(m)

	// sorts keySlice by its key hash, thus the same
	// groups of values are sorted together by key hash
	sort.Sort(SortData(keySlice))

	byteIdx := uint16(0)
	byteIdx--

	keyIdx := uint16(0)
	oldHash := make([]byte, b.config.BktConfig.HashLen)

	for _, key := range keySlice {

		// if the last hash is not equal to current one,
		// means this is a new group of slice.
		if !keyHashIsEqual(oldHash, key[:]) {
			// update and keep track of current hash
			copy(oldHash, key[:])

			// initialize a new byte slice
			byteSlice = append(byteSlice, []byte{})

			// reset the keyIdx and increment byteIdx
			keyIdx = 0
			byteIdx++
		}

		// to ensure keyIdx correctness
		if binary.BigEndian.Uint16(key[b.config.KeyLen:b.config.KeyLen+b.config.SerNumLen]) == keyIdx {
			byteSlice[byteIdx] = append(byteSlice[byteIdx], key[b.config.KeyLen+b.config.SerNumLen:]...)
			keyIdx++
		} else {
			log.Error("wrong bloom slice index", "keyIdx", keyIdx, "key", key)
		}
	}

	for i := range byteSlice {
		// must remove the trailing zeros
		byteSlice[i] = bytes.TrimRight(byteSlice[i], "\x00")
	}
	return
}

// keyHashIsEqual compares byte slice a and b, return true if
// they are equal element-wise
func keyHashIsEqual(a, b []byte) bool {
	for idx, ele := range a {
		if ele != b[idx] {
			return false
		}
	}
	return true
}

// This is an accessory function to convert the input into a map
func interfaceToInterfaceMap(i interface{}) map[interface{}]interface{} {
	if reflect.TypeOf(i).Kind() != reflect.Map {
		panic("input is not a map")
	}

	m := reflect.ValueOf(i)
	res := make(map[interface{}]interface{})

	for _, key := range m.MapKeys() {
		res[key.Interface()] = m.MapIndex(key).Interface()
	}

	return res
}
