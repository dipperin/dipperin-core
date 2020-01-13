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

package model

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/third-party/crypto"
	"github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"io/ioutil"
	"log"
	"math/big"
	"os/exec"
	"path/filepath"
	"regexp"
	"time"
)

var (
	alicePriv = "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232031"
	bobPriv   = "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232032"
	aliceAddr = common.HexToAddress("0x00005586B883Ec6dd4f8c26063E18eb4Bd228e59c3E9")
	bobAddr   = common.HexToAddress("0x0000970e8128aB834E8EAC17aB8E3812f010678CF791")
)

func CreateKey() (*ecdsa.PrivateKey, *ecdsa.PrivateKey) {
	key1, _ := crypto.HexToECDSA(alicePriv)
	key2, _ := crypto.HexToECDSA(bobPriv)
	return key1, key2
}

func CreateSignedVote(height, round uint64, blockId common.Hash, voteType VoteMsgType) *VoteMsg {
	voteA := NewVoteMsg(height, round, blockId, voteType)
	key, _ := CreateKey()
	sign, _ := crypto.Sign(voteA.Hash().Bytes(), key)
	voteA.Witness.Address = aliceAddr
	voteA.Witness.Sign = sign
	return voteA
}

func CreateSignedTx(nonce uint64, amount *big.Int) *Transaction {
	key1, _ := CreateKey()
	fs1 := NewSigner(big.NewInt(1))
	testTx1 := NewTransaction(nonce, bobAddr, amount, TestGasPrice, TestGasLimit, []byte{})
	signedTx, _ := testTx1.SignTx(key1, fs1)
	return signedTx
}

func CreateSignedTxList(n int) []*Transaction {
	keyAlice, _ := CreateKey()
	ms := NewSigner(big.NewInt(1))

	var res []*Transaction
	for i := 0; i < n; i++ {
		tempTx := NewTransaction(uint64(i), bobAddr, big.NewInt(1000), TestGasPrice, TestGasLimit, []byte{})
		gasUsed, _ := IntrinsicGas(tempTx.ExtraData(), false, false)
		tempTx.PaddingActualTxFee(big.NewInt(0).Mul(big.NewInt(int64(gasUsed)), TestGasPrice))
		tempTx.SignTx(keyAlice, ms)
		res = append(res, tempTx)
	}
	return res
}

func CreateBlock(num uint64, preHash common.Hash, txsNum int) *Block {
	header := NewHeader(0, num, preHash, common.HexToHash("123456"), common.HexToDiff("1fffffff"), big.NewInt(time.Now().UnixNano()), aliceAddr, common.BlockNonce{})

	// tx list
	txList := CreateSignedTxList(txsNum)

	// vote
	var voteList []AbstractVerification
	block := NewBlock(header, txList, voteList)

	// calculate block nonce
	CalNonce(block)
	block.RefreshHashCache()
	return block
}

func createTestTx() (*Transaction, *Transaction) {
	_, key2 := CreateKey()
	fs2 := NewSigner(big.NewInt(3))
	hashLock := cs_crypto.Keccak256Hash([]byte("123"))
	tx1 := CreateSignedTx(10, big.NewInt(100))
	tx2 := CreateRawLockTx(1, hashLock, big.NewInt(34564), big.NewInt(10000), TestGasPrice, TestGasLimit, aliceAddr, bobAddr)
	tx2.SignTx(key2, fs2)
	return tx1, tx2
}

// calculate block nonce
func CalNonce(block *Block) {

	preRlp := block.header.RlpBlockWithoutNonce()
	work := &defaultWork{
		header: block.header,
		preRlp: preRlp,
	}
	work.changeNonce()
}

type defaultWork struct {
	header *Header
	preRlp []byte
}

func (work *defaultWork) changeNonce() {
	nonce := work.header.Nonce
	for {
		for index := len(nonce) - 1; index >= 0; {
			if nonce[index] < 255 {
				nonce[index]++
				break
			} else {
				nonce[index] = 0
				index--
			}
		}
		copy(work.header.Nonce[:], nonce[:])
		bHash, err := work.calHash()
		if err == nil {
			if bHash.ValidHashForDifficulty(work.header.GetDifficulty()) {
				break
			}
		}
	}
}

func (work *defaultWork) calHash() (common.Hash, error) {
	if len(work.preRlp) == 0 {
		return common.Hash{}, errors.New("DefaultWork rlp be not calculated yet")
	}

	//extract nonce
	nonce := work.header.Nonce
	raw := append(work.preRlp, nonce[:]...)

	//calculate hash
	return cs_crypto.Keccak256Hash(raw), nil
}

var (
	TestGasPrice = big.NewInt(1)
	TestGasLimit = 2 * TxGas
	TestValue    = big.NewInt(100)
)

type ContractTestParameter struct {
	NodeName      string
	WASMPath      string
	AbiPath       string
	InitInputPara string
	CallFuncName  string
	CallInputPara string
}

var ContractTestPar = ContractTestParameter{
	NodeName:      "default_v0",
	WASMPath:      GetWASMPath("dipclib_test", DIPCTestContract),
	AbiPath:       GetAbiPath("dipclib_test", DIPCTestContract),
	InitInputPara: "",
	CallFuncName:  "libTest",
	CallInputPara: "",
}

var testCorePath = "src/github.com/dipperin/dipperin-core/core/vm/test-data"
var testDIPCPath = "c++/src/dipc/testcontract/"

type ContractPathType uint8

const (
	CoreVmTestData ContractPathType = iota
	DIPCTestContract
)

var contractPath = map[ContractPathType]string{
	CoreVmTestData:   testCorePath,
	DIPCTestContract: testDIPCPath,
}

func getGoPath() string {
	cmd := exec.Command("go", "env")
	str, _ := cmd.Output()

	pattern := `GOPATH="(.*?)"{1}?`
	regular1, _ := regexp.Compile(pattern)
	out := regular1.FindStringSubmatch(string(str))
	return string(out[1])
}

func GetWASMPath(fileName string, pathType ContractPathType) string {
	var tmpPath string
	if pathType == CoreVmTestData {
		tmpPath = getGoPath()
		log.Println("get WASM path", tmpPath)
	} else {
		tmpPath = util.HomeDir()
	}
	log.Println("the tmpPath is:", tmpPath)
	path := filepath.Join(tmpPath, contractPath[pathType])
	return filepath.Join(path, fmt.Sprintf("%s/%s.wasm", fileName, fileName))
}

func GetAbiPath(fileName string, pathType ContractPathType) string {
	var tmpPath string
	if pathType == CoreVmTestData {
		tmpPath = getGoPath()
	} else {
		tmpPath = util.HomeDir()
	}
	path := filepath.Join(tmpPath, contractPath[pathType])
	return filepath.Join(path, fmt.Sprintf("%s/%s.cpp.abi.json", fileName, fileName))
}

func GetCodeAbi(code, abi string) ([]byte, []byte) {
	fileCode, err := ioutil.ReadFile(code)
	if err != nil {
		panic(fmt.Sprintf("Read code failed, err=%s", err.Error()))
	}
	fileABI, err := ioutil.ReadFile(abi)
	if err != nil {
		panic(fmt.Sprintf("Read abi failed, err=%s", err.Error()))
	}
	return fileCode, fileABI
}

func GetCallExtraData(funcName, param string) ([]byte, error) {
	input := []interface{}{
		funcName,
		param,
	}

	result, err := rlp.EncodeToBytes(input)
	return result, err
}

func GetCreateExtraData(wasmPath, abiPath string, init string) ([]byte, error) {
	// GetContractExtraData
	WASMBytes, err := ioutil.ReadFile(wasmPath)
	if err != nil {
		return WASMBytes, err
	}

	abiBytes, err := ioutil.ReadFile(abiPath)
	if err != nil {
		return abiBytes, err
	}

	var rlpParams []interface{}
	if init == "" {
		rlpParams = []interface{}{
			WASMBytes, abiBytes,
		}
	} else {
		rlpParams = []interface{}{
			WASMBytes, abiBytes, init,
		}
	}

	data, err := rlp.EncodeToBytes(rlpParams)
	return data, err
}
