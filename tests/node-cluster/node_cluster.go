package node_cluster

import (
	"fmt"
	"path/filepath"
	"io/ioutil"
	"strconv"
	"github.com/dipperin/dipperin-core/third-party/rpc"
	"github.com/dipperin/dipperin-core/common/util"
	"errors"
	"github.com/dipperin/dipperin-core/common"
)

const (
	clusterConfFileName = "default_cluster.json"
)

type NodeCluster struct {
	ClusterConfigure []NodeConf
	NodeClient       map[string]*rpc.Client
	NodeConfigure    map[string]NodeConf
}

func (cluster NodeCluster) GetNodeMainAddress(name string) (common.Address, error) {
	address := cluster.NodeConfigure[name].Address
	if address == "" {
		return common.Address{}, errors.New(fmt.Sprintf("can't find %d main address", name))
	}
	return common.HexToAddress(address), nil
}

type NodeConf struct {
	NodeName    string `json:"node_name"`
	NodeURL     string `json:"node_url"`
	NodeID      string `json:"node_id"`
	P2PListener string `json:"p2p_listener"`
	HttpPort    int    `json:"http_port"`
	WsPort      int    `json:"ws_port"`
	Host        string `json:"host"`
	Address     string `json:"address"`
}

func CreateNodeCluster() (cluster *NodeCluster, err error) {
	configure, err := getClusterConfig()
	if err != nil {
		return nil, err
	}

	clientMap, NodeConfigure, err := newNodeClientAndConf(configure)
	if err != nil {
		return nil, err
	}

	return &NodeCluster{
		ClusterConfigure: configure,
		NodeClient:       clientMap,
		NodeConfigure:    NodeConfigure,
	}, nil
}

func getClusterConfig() (configure []NodeConf, err error) {
	confPath := filepath.Join(util.HomeDir(), clusterConfFileName)

	// 尝试从配置文件中加载配置
	if fb, err := ioutil.ReadFile(confPath); err == nil {
		if err = util.ParseJsonFromBytes(fb, &configure); err != nil {
			return []NodeConf{}, err
		}
	} else {
		return []NodeConf{}, err
	}
	return configure, nil
}

func newNodeClientAndConf(configure []NodeConf) (map[string]*rpc.Client, map[string]NodeConf, error) {
	clientMap := make(map[string]*rpc.Client)
	nodeConfigure := make(map[string]NodeConf)

	// verifiers
	for _, value := range configure {
		httpPort := strconv.Itoa(value.HttpPort)
		client := newRpcClient(value.Host, httpPort)
		clientMap[value.NodeName] = client
		nodeConfigure[value.NodeName] = value
	}
	return clientMap, nodeConfigure, nil
}
