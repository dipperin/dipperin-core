package sidechain

import (
	"fmt"
	"github.com/dipperin/dipperin-core/cmd/dipperin-across/node_cluster"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/log"
	"github.com/dipperin/dipperin-core/core/rpc-interface"
	"github.com/dipperin/dipperin-core/tests/g-testData"
	"go.uber.org/zap"
	"math/big"
	"time"
)

func CallSidechainContract(cluster *node_cluster.NodeCluster, nodeName string, contractAddr, oracleAddr common.Address, proof []byte) error {
	// transfer
	param := fmt.Sprintf("%s,%s", common.Bytes2Hex(proof), oracleAddr.Hex())
	data, err := g_testData.GetCallExtraData("transfer", param)
	if err != nil {
		return err
	}

	txHash, err := cluster.SendContract(nodeName, contractAddr, big.NewInt(0), data)
	log.DLogger.Info("Transfer sidechain contract processing", zap.String("txHash", txHash.String()))
	if err != nil {
		return err
	}

	for {
		time.Sleep(time.Second)
		resp, innerErr := cluster.Transaction(nodeName, txHash)
		if innerErr != nil {
			return innerErr
		}

		if resp.BlockNumber != 0 {
			break
		}
	}
	log.DLogger.Info("Transfer sidechain contract successful", zap.String("txHash", txHash.String()))
	return nil
}

func GetSPVProof(cluster *node_cluster.NodeCluster, from, to string, amount *big.Int) ([]byte, uint64, error) {
	txHash, err := cluster.SendTx(from, to, amount)
	if err != nil {
		return nil, 0, err
	}

	var resp *rpc_interface.TransactionResp
	for {
		time.Sleep(time.Second)
		resp, err = cluster.Transaction(from, txHash)
		if err != nil {
			return nil, 0, err
		}

		if resp.BlockNumber != 0 {
			break
		}
	}

	spvProof, err := cluster.GetSPVProof(from, txHash)
	if err != nil {
		return nil, 0, err
	}

	log.DLogger.Info("SendTx and GetSPVProof successful", zap.String("txHash", txHash.Hex()), zap.Uint64("height", resp.BlockNumber))
	return spvProof, resp.BlockNumber, nil
}

// 部署侧链智能合约，参数（交易人、chainID）
func InitSidechainContract(cluster *node_cluster.NodeCluster, from, to string, amount *big.Int, chainID uint64) (common.Address, error) {
	WASMPath := g_testData.GetWASMPath("sidechain", g_testData.CoreVmTestData)
	AbiPath := g_testData.GetAbiPath("sidechain", g_testData.CoreVmTestData)
	addr := cluster.NodeConfigure[to].Address
	param := fmt.Sprintf("%s,%v", addr, chainID)
	data, err := g_testData.GetCreateExtraData(WASMPath, AbiPath, param)
	if err != nil {
		return common.Address{}, err
	}

	create := common.HexToAddress(common.AddressContractCreate)
	txHash, err := cluster.SendContract(from, create, amount, data)
	if err != nil {
		return common.Address{}, err
	}

	for {
		time.Sleep(time.Second)
		resp, innerErr := cluster.Transaction(from, txHash)
		if innerErr != nil {
			return common.Address{}, innerErr
		}

		if resp.BlockNumber != 0 {
			break
		}
	}

	contractAddr, err := cluster.GetContractAddressByTxHash(from, txHash)
	if err != nil {
		return common.Address{}, err
	}

	log.DLogger.Info("Init sidechain contract successful", zap.String("txHash", txHash.Hex()), zap.String("contractAddr", contractAddr.Hex()))
	return contractAddr, nil
}
