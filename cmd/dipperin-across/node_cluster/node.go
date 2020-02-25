package node_cluster

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/third-party/rpc"
	"path/filepath"
	"strings"
)

type Node struct {
	Client  *rpc.Client
	Address common.Address
}

// 新建一个ipc client
func newIpcClient(node, env string) *rpc.Client {
	path := filepath.Join(util.HomeDir(), "tmp/dipperin_apps/", env, node, "dipperin.ipc")
	client, err := rpc.Dial(path)
	if err != nil {
		panic("init rpc client failed: " + err.Error())
	}
	return client
}

func getRpcTXMethod(methodName string) string {
	return "dipperin_" + strings.ToLower(methodName[0:1]) + methodName[1:]
}
