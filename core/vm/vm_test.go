package vm

/*
import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/core/vmcommon/common/utils"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"math/big"
	"testing"
	"time"
)

var vmcommon *VM

func TestDeployContract(t *testing.T) {
	InitVM()
	code, _ := ioutil.ReadFile("./map-string/map2.wasm")
	abiJson,_:=ioutil.ReadFile("./map-string/StringMap.cpp.abi.json")

	txCaller := &Caller{common.HexToAddress("0x999")}
	_, addr, _, err := vmcommon.Create(txCaller, code, abiJson, []byte{})

	assert.NoError(t, err)
	result := vmcommon.state.GetState(addr,[]byte("code"))
	assert.Equal(t,code,result)
}

func TestCallContract(t *testing.T){
	// Deploy contract before call it
	InitVM()
	code, _ := ioutil.ReadFile("./map-string/map2.wasm")
	abiJson,_:=ioutil.ReadFile("./map-string/StringMap.cpp.abi.json")

	txCaller := &Caller{common.HexToAddress("0x999")}
	_, addr, _, err := vmcommon.Create(txCaller, code, abiJson, []byte{})
	assert.NoError(t, err)

	// Call contract
	param := [][]byte{[]byte("alice"), utils.Int32ToBytes(100)} //key = "Alice" value = 100
	inputs := genInput(t,"setBalance",param)

	_, _, err = vmcommon.Call(txCaller, addr, inputs,0,big.NewInt(0))
	assert.NoError(t, err)

	//Verify result
	time.Sleep(time.Microsecond*100)
	vmcommon.state.GetState(addr,append([]byte{7}, []byte("balance")...))
}

func InitVM(){
	if vmcommon != nil {
		return
	}
	vmcommon = NewVM(Context{},&fakeStateDB{}, DEFAULT_VM_CONFIG)
}
*/
