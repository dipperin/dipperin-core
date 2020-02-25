package oracle

import (
	"fmt"
	"github.com/dipperin/dipperin-core/cmd/dipperin-across/node_cluster"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/log"
	"github.com/dipperin/dipperin-core/core/dipperin/service"
	"github.com/dipperin/dipperin-core/core/rpc-interface"
	"github.com/dipperin/dipperin-core/core/spv"
	"github.com/dipperin/dipperin-core/tests/g-testData"
	"github.com/ethereum/go-ethereum/rlp"
	"go.uber.org/zap"
	"math/big"
	"time"
)

func SubscribeChainBlock(clusterA, clusterB *node_cluster.NodeCluster, nodeName string, contractAddr common.Address) {
	block := make(chan service.SubBlockResp)
	subscription, err := clusterB.SubscribeChainBlockEvent(nodeName, block)
	if err != nil {
		log.DLogger.Error("SubscribeChainBlockEvent failed")
		panic(err)
	}

	defer subscription.Unsubscribe()
	for {
		select {
		case blockInfo := <-block:
			innerErr := CallOracleContract(clusterA, clusterB, nodeName, contractAddr, blockInfo.Number)
			if innerErr != nil {
				log.DLogger.Error("CallOracleContract failed")
				panic(innerErr)
			}
		}
	}
}

// 部署预研机智能合约
func InitOracleContract(cluster *node_cluster.NodeCluster, nodeName string) (common.Address, error) {
	WASMPath := g_testData.GetWASMPath("oracle", g_testData.CoreVmTestData)
	AbiPath := g_testData.GetAbiPath("oracle", g_testData.CoreVmTestData)
	data, err := g_testData.GetCreateExtraData(WASMPath, AbiPath, "")
	if err != nil {
		return common.Address{}, err
	}

	to := common.HexToAddress(common.AddressContractCreate)
	txHash, err := cluster.SendContract(nodeName, to, big.NewInt(0), data)
	if err != nil {
		return common.Address{}, err
	}

	for {
		time.Sleep(time.Second)
		resp, innerErr := cluster.Transaction(nodeName, txHash)
		if innerErr != nil {
			return common.Address{}, innerErr
		}

		if resp.BlockNumber != 0 {
			break
		}
	}

	contractAddr, err := cluster.GetContractAddressByTxHash(nodeName, txHash)
	if err != nil {
		return common.Address{}, err
	}

	log.DLogger.Info("Init oracle contract successful", zap.String("txHash", txHash.Hex()), zap.String("contractAddr", contractAddr.Hex()))
	return contractAddr, nil
}

// 下载B链区块，存入A链
func CallOracleContract(clusterA, clusterB *node_cluster.NodeCluster, nodeName string, contractAddr common.Address, height uint64) error {
	var (
		block rpc_interface.BlockResp
		err   error
	)
	for {
		time.Sleep(time.Second)
		block, err = clusterB.GetBlockByNumber(nodeName, height)
		if err == nil {
			break
		}
	}

	spvHeader := spv.SPVHeader{
		ChainID: block.Header.ChainID,
		Hash:    block.Header.Hash(),
		Height:  block.Header.Number,
		TxRoot:  block.Header.TransactionRoot,
	}
	key := fmt.Sprintf("%v", block.Header.Number)
	value, err := rlp.EncodeToBytes(spvHeader)
	if err != nil {
		return err
	}

	// set block
	param := fmt.Sprintf("%s,%s", key, common.Bytes2Hex(value))
	data, err := g_testData.GetCallExtraData("setHeader", param)
	if err != nil {
		return err
	}

	txHash, err := clusterA.SendContract(nodeName, contractAddr, big.NewInt(0), data)
	if err != nil {
		return err
	}

	for {
		time.Sleep(time.Second)
		resp, innerErr := clusterA.Transaction(nodeName, txHash)
		if innerErr != nil {
			return innerErr
		}

		if resp.BlockNumber != 0 {
			break
		}
	}
	log.DLogger.Info("Set spv header successful", zap.String("txHash", txHash.String()), zap.Uint64("chainID", spvHeader.ChainID), zap.Uint64("height", spvHeader.Height))
	return nil
}
