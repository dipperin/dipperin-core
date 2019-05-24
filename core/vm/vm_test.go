package vm

import (
	"github.com/dipperin/dipperin-core/core/vm/common/utils"
	"github.com/dipperin/dipperin-core/common"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"math/big"
	"testing"
	"time"
)

var vm *VM

func TestDeployContract(t *testing.T) {
	InitVM()
	code, _ := ioutil.ReadFile("./map-string/map2.wasm")
	abiJson,_:=ioutil.ReadFile("./map-string/StringMap.cpp.abi.json")

	txCaller := &Caller{common.HexToAddress("0x999")}
	_, addr, _, err := vm.Create(txCaller, code, abiJson, big.NewInt(0))

	assert.NoError(t, err)
	result := vm.state.GetState(addr,[]byte("code"))
	assert.Equal(t,code,result)
}

func TestCallContract(t *testing.T){
	// Deploy contract before call it
	InitVM()
	code, _ := ioutil.ReadFile("./map-string/map2.wasm")
	abiJson,_:=ioutil.ReadFile("./map-string/StringMap.cpp.abi.json")

	txCaller := &Caller{common.HexToAddress("0x999")}
	_, addr, _, err := vm.Create(txCaller, code, abiJson, big.NewInt(0))
	assert.NoError(t, err)

	// Call contract
	param := [][]byte{[]byte("alice"), utils.Int32ToBytes(100)} //key = "Alice" value = 100
	inputs := genInput(t,"setBalance",param)

	_, _, err = vm.Call(txCaller, addr, inputs,big.NewInt(0))
	assert.NoError(t, err)

	//Verify result
	time.Sleep(time.Microsecond*100)
	vm.state.GetState(addr,append([]byte{7}, []byte("balance")...))
}

func InitVM(){
	if vm != nil {
		return
	}
	storage := NewStorage()
	vm = NewVM(Context{}, storage, DEFAULT_VM_CONFIG)
}
