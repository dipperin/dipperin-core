package node_cluster

import (
	"context"
	"errors"
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/g-event"
	"github.com/dipperin/dipperin-core/common/log"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/core/rpc-interface"
	"github.com/dipperin/dipperin-core/third-party/rpc"
	"go.uber.org/zap"
	"io/ioutil"
	"math/big"
	"path/filepath"
)

var (
	gasPrice = big.NewInt(7)
	gasLimit = uint64(5000000)
)

type NodeCluster struct {
	ClusterConfigure []NodeConf
	NodeClient       map[string]*rpc.Client
	NodeConfigure    map[string]NodeConf
}

func (cluster NodeCluster) GetNodeMainAddress(nodeName string) (common.Address, error) {
	address := cluster.NodeConfigure[nodeName].Address
	if address == "" {
		return common.Address{}, errors.New(fmt.Sprintf("can't find %s main address", nodeName))
	}
	return common.HexToAddress(address), nil
}

func (cluster NodeCluster) GetAddressBalance(nodeName string) (*big.Int, error) {
	client := cluster.NodeClient[nodeName]
	addr := common.HexToAddress(cluster.NodeConfigure[nodeName].Address)
	var resp rpc_interface.CurBalanceResp
	if err := client.Call(&resp, getRpcTXMethod("CurrentBalance"), addr); err != nil {
		return big.NewInt(0), err
	}
	log.DLogger.Info("GetAddressBalance", zap.String("addr", addr.Hex()), zap.Any("balance", resp.Balance.ToInt()))
	return resp.Balance.ToInt(), nil
}

func (cluster NodeCluster) SubscribeChainBlockEvent(nodeName string, channel interface{}) (g_event.Subscription, error) {
	client := cluster.NodeClient[nodeName]
	return client.Subscribe(context.Background(), "dipperin", channel, "subscribeBlock")
}

func (cluster NodeCluster) GetSPVProof(nodeName string, txHash common.Hash) ([]byte, error) {
	client := cluster.NodeClient[nodeName]
	var resp []byte
	err := client.Call(&resp, getRpcTXMethod("GetSPVProof"), txHash)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (cluster NodeCluster) SendTx(nodeName, to string, amount *big.Int) (common.Hash, error) {
	client := cluster.NodeClient[nodeName]
	from := common.HexToAddress(cluster.NodeConfigure[nodeName].Address)
	toAddr := common.HexToAddress(cluster.NodeConfigure[to].Address)

	var resp common.Hash
	err := client.Call(&resp, getRpcTXMethod("SendTransaction"), from, toAddr, amount, gasPrice, gasLimit, []byte{}, nil)
	if err != nil {
		return common.Hash{}, err
	}
	return resp, nil
}

func (cluster NodeCluster) SendContract(nodeName string, to common.Address, amount *big.Int, data []byte) (common.Hash, error) {
	client := cluster.NodeClient[nodeName]
	from := common.HexToAddress(cluster.NodeConfigure[nodeName].Address)

	var resp common.Hash
	err := client.Call(&resp, getRpcTXMethod("SendTransactionContract"), from, to, amount, gasPrice, gasLimit, data, nil)
	if err != nil {
		return common.Hash{}, err
	}
	return resp, nil
}

func (cluster NodeCluster) CallContract(nodeName string, to common.Address, data []byte) (string, error) {
	client := cluster.NodeClient[nodeName]
	from := common.HexToAddress(cluster.NodeConfigure[nodeName].Address)

	var resp string
	err := client.Call(&resp, getRpcTXMethod("CallContract"), from, to, data, uint64(0))
	if err != nil {
		return "", err
	}
	return resp, nil
}

func (cluster NodeCluster) Transaction(nodeName string, txHash common.Hash) (*rpc_interface.TransactionResp, error) {
	client := cluster.NodeClient[nodeName]

	var resp *rpc_interface.TransactionResp
	err := client.Call(&resp, getRpcTXMethod("Transaction"), txHash)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (cluster NodeCluster) GetContractAddressByTxHash(nodeName string, txHash common.Hash) (common.Address, error) {
	client := cluster.NodeClient[nodeName]

	var resp common.Address
	err := client.Call(&resp, getRpcTXMethod("GetContractAddressByTxHash"), txHash)
	if err != nil {
		return common.Address{}, err
	}
	return resp, nil
}

func (cluster NodeCluster) GetBlockByNumber(nodeName string, num uint64) (rpc_interface.BlockResp, error) {
	client := cluster.NodeClient[nodeName]
	var respBlock rpc_interface.BlockResp
	err := client.Call(&respBlock, getRpcTXMethod("GetBlockByNumber"), num)
	if err != nil {
		return rpc_interface.BlockResp{}, err
	}
	return respBlock, nil
}

func (cluster NodeCluster) CurrentBlock(nodeName string) (rpc_interface.BlockResp, error) {
	client := cluster.NodeClient[nodeName]
	var respBlock rpc_interface.BlockResp
	err := client.Call(&respBlock, getRpcTXMethod("CurrentBlock"))
	if err != nil {
		return rpc_interface.BlockResp{}, err
	}
	return respBlock, nil
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

func CreateIpcNodeCluster(env string) (*NodeCluster, error) {
	configure, err := getClusterConfig(env)
	if err != nil {
		return nil, err
	}

	clientMap, NodeConfigure, err := newIpcClientAndConf(configure, env)
	if err != nil {
		return nil, err
	}

	return &NodeCluster{
		ClusterConfigure: configure,
		NodeClient:       clientMap,
		NodeConfigure:    NodeConfigure,
	}, nil
}

func getClusterConfig(env string) (configure []NodeConf, err error) {
	var path string
	switch env {
	case "local":
		path = filepath.Join(util.HomeDir(), "default_cluster_local.json")
	case "tps":
		path = filepath.Join(util.HomeDir(), "default_cluster_tps.json")
	}
	// 尝试从配置文件中加载配置
	var fb []byte
	if fb, err = ioutil.ReadFile(path); err != nil {
		return []NodeConf{}, err
	}

	if err = util.ParseJsonFromBytes(fb, &configure); err != nil {
		return []NodeConf{}, err
	}

	return configure, nil
}

func newIpcClientAndConf(configure []NodeConf, env string) (map[string]*rpc.Client, map[string]NodeConf, error) {
	clientMap := make(map[string]*rpc.Client)
	nodeConfigure := make(map[string]NodeConf)

	// verifiers
	for _, value := range configure {
		clientMap[value.NodeName] = newIpcClient(value.NodeName, env)
		nodeConfigure[value.NodeName] = value
	}
	return clientMap, nodeConfigure, nil
}
