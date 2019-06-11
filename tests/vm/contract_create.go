package vm

import (
	"math/big"
	"github.com/dipperin/dipperin-core/core/economy-model"
	"github.com/dipperin/dipperin-core/third-party/rpc"
	"github.com/dipperin/dipperin-core/common"
	"strings"
	"github.com/dipperin/dipperin-core/third-party/log"
)

func LogTestPrint(function, msg string, ctx ...interface{}) {
	printMsg := "[~wjw~" + function + "]" + msg
	log.Info(printMsg, ctx...)
}

func getRpcTXMethod(methodName string) string {
	return "dipperin_" + strings.ToLower(methodName[0:1]) + methodName[1:]
}

func sendTransaction(client *rpc.Client, from, to common.Address, value int64, nonce uint64) (common.Hash, error) {

	var Tx common.Hash
	ExtraData := make([]byte, 0)
	if err := client.Call(&Tx, getRpcTXMethod("SendTransaction"), from, to, big.NewInt(value), economy_model.GetMinimumTxFee(500), ExtraData, nonce); err != nil {
		return Tx, err
	}
	return Tx, nil
}
