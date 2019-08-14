package contract

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/tests/g-testData"
	"math/big"
)

const (
	createERC20JsonWithOwner = `{
		"owner":"0x1234",
		"token_name":"EOS",
		"token_decimals":18,
		"token_symbol":"EOS",
		"token_total_supply":"0x33b2e3c9fd0803ce8000000",
		"balances":{},
		"allowed":{}
	}`

	createERC20Json = `{
		"token_name":"EOS",
		"token_decimals":18,
		"token_symbol":"EOS",
		"token_total_supply":"0x33b2e3c9fd0803ce8000000",
		"balances":{},
		"allowed":{}
	}`
)

var (
	aliceAddr   = common.HexToAddress("0x00005586B883Ec6dd4f8c26063E18eb4Bd228e59c3E9")
	bobAddr     = common.HexToAddress("0x0000970e8128aB834E8EAC17aB8E3812f010678CF791")
	charlieAddr = common.HexToAddress("0x00007dbbf084F4a6CcC070568f7674d4c2CE8CD2709E")
	erc20Addr   = common.HexToAddress("0x00100000FA42f7315cD04D6774E58B54e92603e96d84")
)

func CreateSignedTx(to common.Address, extraData []byte) *model.Transaction {
	key1, _ := model.CreateKey()
	fs1 := model.NewSigner(big.NewInt(1))
	testTx1 := model.NewTransaction(0, to, big.NewInt(0), g_testData.TestGasPrice, g_testData.TestGasLimit, extraData)
	signedTx, _ := testTx1.SignTx(key1, fs1)
	return signedTx
}

func newTestToken() *BuiltInERC20Token {
	var token BuiltInERC20Token
	if err := util.ParseJson(createERC20Json, &token); err != nil {
		panic(err.Error())
	}
	token.Owner = aliceAddr
	return &token
}
