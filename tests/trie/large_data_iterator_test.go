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

package trie

import (
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
	trie2 "github.com/dipperin/dipperin-core/third-party/trie"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"math/big"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"testing"
)

type ContractBase struct {
	CurSender common.Address `json:"-"`
}

type BuiltInERC20Token struct {
	ContractBase

	Owner            common.Address                 `json:"owner"`
	TokenName        string                         `json:"token_name"`
	TokenDecimals    uint                           `json:"token_decimals"`
	TokenSymbol      string                         `json:"token_symbol"`
	TokenTotalSupply *big.Int                       `json:"token_total_supply"`
	Balances         map[string]*big.Int            `json:"balances"`
	Allowed          map[string]map[string]*big.Int `json:"allowed"`
}

func newERC20(owner common.Address) *BuiltInERC20Token {
	c := &BuiltInERC20Token{
		Owner:            owner,
		TokenName:        owner.Hex() + "_x",
		TokenDecimals:    18,
		TokenSymbol:      "ccc",
		TokenTotalSupply: big.NewInt(200000000000),
		Balances: map[string]*big.Int{
			"0x123213": big.NewInt(89498),
			"0x123211": big.NewInt(89498),
		},
		Allowed: map[string]map[string]*big.Int{
			"0x1231232": {"0x321312": big.NewInt(2313)},
			"0x1231231": {"0x321315": big.NewInt(231334)},
			"0x1231211": {"0x1321315": big.NewInt(23134)},
		},
	}
	return c
}

func TestCToKV(t *testing.T) {
	c := newERC20(common.HexToAddress("0x123123213"))
	kv := obj2KV2(c)

	// kv to json str
	cJsonStr, err := kv2JsonStr(kv)
	assert.NoError(t, err)

	// parse from kv
	var tmpC BuiltInERC20Token
	err = util.ParseJson(cJsonStr, &tmpC)
	assert.NoError(t, err)
	assert.Equal(t, c.Owner, tmpC.Owner)
	assert.Equal(t, c.TokenTotalSupply, tmpC.TokenTotalSupply)
	assert.Equal(t, c.Allowed["0x1231232"]["0x321312"], tmpC.Allowed["0x1231232"]["0x321312"])
}

/*Test 1

1. Create a contract structure and a random number of accounts,
2. Flatten the contract under each account and deposit it in MPT.
3. Load contract data of some accounts from MPT and test efficiency

Test 2

Putting the root of the contract in a certain value of the account to see if it's more efficient than test 1
*/
func TestLargeDataIterator(t *testing.T) {
	ldbPath := filepath.Join(util.HomeDir(), "tmp", "test_mpt_data")
	db, err := ethdb.NewLDBDatabase(ldbPath, 0, 0)
	if err != nil {
		panic(err)
	}
	defer func() {
		db.Close()
		os.RemoveAll(ldbPath)
	}()

	trie, err := trie2.New(common.Hash{}, trie2.NewDatabase(db))
	assert.NoError(t, err)

	addresses := genAddresses(10000)

	// save contract
	for _, addr := range addresses {
		c := newERC20(addr)
		cKV := obj2KV(c)
		for k, v := range cKV {
			trie.TryUpdate(getContractKey(addr, k), getContractValue(v))
		}
	}
}

func getContractKey(addr common.Address, ck string) []byte {
	return append(addr.Bytes(), []byte("_ERC20_"+ck)...)
}

func getContractValue(v interface{}) []byte {
	fmt.Println(reflect.TypeOf(v), v)
	return []byte("321")
}

func genAddresses(count int) (result []common.Address) {
	for i := 0; i < count; i++ {
		pk, err := crypto.GenerateKey()
		if err != nil {
			panic(err)
		}
		result = append(result, cs_crypto.GetNormalAddress(pk.PublicKey))
	}
	return
}

func obj2KV(obj interface{}) map[string]interface{} {
	value := gjson.Parse(util.StringifyJson(obj))

	result := map[string]interface{}{}
	json2KV("", value, result)

	return result
}

func json2KV(key string, json gjson.Result, result map[string]interface{}) {
	if json.IsObject() {

		json.ForEach(func(key1, value gjson.Result) bool {

			if key == "" {
				json2KV(key1.Str, value, result)
			} else {
				json2KV(key+"."+key1.Str, value, result)
			}

			return true
		})

	} else if json.IsArray() {
		index := 0
		json.ForEach(func(key1, value gjson.Result) bool {

			if key == "" {
				json2KV(strconv.Itoa(index), value, result)
			} else {
				json2KV(key+"."+strconv.Itoa(index), value, result)
			}

			index++
			return true
		})

	} else {
		result[key] = json.Value()
	}
}
