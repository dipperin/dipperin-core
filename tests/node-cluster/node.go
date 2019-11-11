package node_cluster

import (
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/third-party/rpc"
	"strings"
	"path/filepath"
	"github.com/dipperin/dipperin-core/common/util"
)

type Node struct {
	Client  *rpc.Client
	Address common.Address
}

// 新建一个rpc client
func newRpcClient(host string, port string) *rpc.Client {
	if host == "" {
		host = "127.0.0.1"
	}
	client, err := rpc.Dial(fmt.Sprintf("http://%v:%v", host, port))
	if err != nil {
		panic(err.Error())
	}
	return client
}

// 新建一个ipc client
func newIpcClient(node string) *rpc.Client {
	path := filepath.Join(util.HomeDir(), "tmp/dipperin_apps/", node, "dipperin.ipc")
	client, err := rpc.Dial(path)
	if err != nil {
		panic("init rpc client failed: " + err.Error())
	}
	return client
}

func getRpcTXMethod(methodName string) string {
	return "dipperin_" + strings.ToLower(methodName[0:1]) + methodName[1:]
}
