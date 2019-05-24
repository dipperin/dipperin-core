package vm

import (
	"bytes"
	platOnLifeExec "github.com/dipperin/dipperin-core/third-party/life/exec"
	"fmt"
	"github.com/perlin-network/life/exec"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"testing"
)

func TestPlatOnLifeMap(t *testing.T){
	fileCode, err := ioutil.ReadFile("./sample/test.wasm")
	assert.NoError(t, err)

	vm, err := platOnLifeExec.NewVirtualMachine(fileCode, DEFAULT_VM_CONFIG, nil, nil)
	assert.NoError(t,err)

	entryID, ok := vm.GetFunctionExport("main")
	assert.Equal(t,true,ok)

	_, err = vm.Run(entryID)
	assert.NoError(t,err)
}


func TestOriginLifeMap(t *testing.T){
	fileCode, err := ioutil.ReadFile("./sample/test.wasm")
	assert.NoError(t, err)

	vm, err := exec.NewVirtualMachine(fileCode, exec.VMConfig{
		EnableJIT:          false,
		DefaultMemoryPages: 128,
		DefaultTableSize:   65536,
	}, nil, nil)
	assert.NoError(t,err)

	entryID, ok := vm.GetFunctionExport("_Z10getmapdatav")
	assert.Equal(t,true,ok)

	data, err := vm.Run(entryID)
	assert.NoError(t,err)

	fmt.Printf("the data is: %d",data)
}

func TestIOByteReader(t *testing.T){
	testData := []byte{
		0x00,0x01,0x02,0x03,0x04,0x05,
	}

	r:=bytes.NewReader(testData)
	byte0,err :=r.ReadByte()
	assert.NoError(t,err)
	assert.Equal(t,byte(0x00),byte0)

	readByte := make([]byte,len(testData)-1)
	n,err :=io.ReadFull(r,readByte)
	assert.NoError(t,err)
	assert.Equal(t,len(testData)-1,n)
	assert.Equal(t,[]byte{0x01,0x02,0x03,0x04,0x05,},readByte)
}

