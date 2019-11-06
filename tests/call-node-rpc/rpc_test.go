package call_node_rpc

import (
	"fmt"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/core/accounts"
	"github.com/dipperin/dipperin-core/core/rpc-interface"
	"github.com/dipperin/dipperin-core/tests/vm"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/dipperin/dipperin-core/third-party/rpc"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func Test_IPC(t *testing.T){
	var err error
	var client *rpc.Client
	if client, err = rpc.Dial("/home/qydev/tmp/dipperin_apps/default_v0/dipperin.ipc"); err != nil {
		panic("init rpc client failed: " + err.Error())
	}
	assert.NoError(t, err)

	//恢复钱包
	mnemonic := "chicken coconut winner february brown topple pond bird endless salt filter journey mass ramp milk tuition card seat worth school length rain slice ozone"
	password := "123"
	identifier := accounts.WalletIdentifier{
		WalletType: accounts.SoftWallet,
		Path:       filepath.Join(util.HomeDir(), "tmp/dipperin_apps/default_v0/CSWallet2"),
		WalletName: "CSWallet2",
	}
	var resp interface{}
	err = client.Call(resp, vm.GetRpcTXMethod("RestoreWallet"), password, mnemonic, "", identifier)
	assert.NoError(t, err)
}

func Test_Http(t *testing.T){
	var err error
	var client *rpc.Client
	if client, err = rpc.Dial(fmt.Sprintf("http://%v:%d", "127.0.0.1", 10019)); err != nil {
		panic("init rpc client failed: " + err.Error())
	}
	assert.NoError(t, err)

	//恢复钱包
	mnemonic := "chicken coconut winner february brown topple pond bird endless salt filter journey mass ramp milk tuition card seat worth school length rain slice ozone"
	password := "123"
	identifier := accounts.WalletIdentifier{
		WalletType: accounts.SoftWallet,
		Path:       filepath.Join(util.HomeDir(), "tmp/dipperin_apps/default_v0/CSWallet3"),
		WalletName: "CSWallet3",
	}
	var resp interface{}
	err = client.Call(resp, vm.GetRpcTXMethod("RestoreWallet"), password, mnemonic, "", identifier)
	log.Info("the error is: ","err",err)
	assert.Error(t, err)

	var resp1 rpc_interface.BlockResp
	err = client.Call(&resp1, vm.GetRpcTXMethod("CurrentBlock"))
	assert.NoError(t,err)
	log.Info("the current block is:","block",resp1)
}

func Test_WebSocket(t *testing.T){
	var err error
	var client *rpc.Client
	if client, err = rpc.Dial(fmt.Sprintf("ws://%v:%d", "127.0.0.1", 10020)); err != nil {
		panic("init rpc client failed: " + err.Error())
	}
	assert.NoError(t, err)

	//恢复钱包
	mnemonic := "chicken coconut winner february brown topple pond bird endless salt filter journey mass ramp milk tuition card seat worth school length rain slice ozone"
	password := "123"
	identifier := accounts.WalletIdentifier{
		WalletType: accounts.SoftWallet,
		Path:       filepath.Join(util.HomeDir(), "tmp/dipperin_apps/default_v0/CSWallet4"),
		WalletName: "CSWallet4",
	}
	var resp interface{}
	err = client.Call(resp, vm.GetRpcTXMethod("RestoreWallet"), password, mnemonic, "", identifier)
	log.Info("the error is: ","err",err)
	assert.Error(t, err)

	var resp1 rpc_interface.BlockResp
	err = client.Call(&resp1, vm.GetRpcTXMethod("CurrentBlock"))
	assert.NoError(t,err)
	log.Info("the current block is:","block",resp1)
}