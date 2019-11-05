package call_node_rpc

import (
	"fmt"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/core/accounts"
	"github.com/dipperin/dipperin-core/tests/vm"
	"github.com/dipperin/dipperin-core/third-party/rpc"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
	"time"
)

func Test_MineMasterServiceStart(t *testing.T){
	var err error
	var client *rpc.Client
	if client, err = rpc.Dial(fmt.Sprintf("http://%v:%d", "127.0.0.1", 20016)); err != nil {
		panic("init rpc client failed: " + err.Error())
	}
	assert.NoError(t,err)

	//恢复钱包
	mnemonic := "chicken coconut winner february brown topple pond bird endless salt filter journey mass ramp milk tuition card seat worth school length rain slice ozone"
	password := "123"
	identifier := accounts.WalletIdentifier{
		WalletType:accounts.SoftWallet,
		Path: filepath.Join(util.HomeDir(),"tmp/dipperin_apps/default_m1/CSWallet"),
		WalletName:"CSWallet",
	}
	var resp interface{}
	err = client.Call(resp, vm.GetRpcTXMethod("RestoreWallet"), password, mnemonic, "", identifier)
	assert.NoError(t,err)

	time.Sleep(100*time.Millisecond)
	err = client.Call(resp,vm.GetRpcTXMethod("StartRemainingService"))
	assert.NoError(t,err)
}

func Test_MineMasterStartMine(t *testing.T){
	var err error
	var client *rpc.Client
	if client, err = rpc.Dial(fmt.Sprintf("http://%v:%d", "127.0.0.1", 20016)); err != nil {
		panic("init rpc client failed: " + err.Error())
	}
	assert.NoError(t,err)

	var resp interface{}
	err = client.Call(resp, vm.GetRpcTXMethod("StartMine"))
	assert.NoError(t,err)
}
