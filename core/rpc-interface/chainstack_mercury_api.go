// Copyright 2019, Keychain Foundation Ltd.
// This file is part of the dipperin-core library.
//
// The dipperin-core library is free software: you can redistribute
// it and/or modify it under the terms of the GNU Lesser General Public License
// as published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// The dipperin-core library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package rpc_interface

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/address-util"
	"github.com/dipperin/dipperin-core/common/config"
	"github.com/dipperin/dipperin-core/common/g-error"
	"github.com/dipperin/dipperin-core/common/hexutil"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/core/accounts"
	"github.com/dipperin/dipperin-core/core/chain-config"
	"github.com/dipperin/dipperin-core/core/contract"
	"github.com/dipperin/dipperin-core/core/dipperin/service"
	"github.com/dipperin/dipperin-core/core/economy-model"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/core/vm/common/utils"
	model2 "github.com/dipperin/dipperin-core/core/vm/model"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/dipperin/dipperin-core/third-party/rpc"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
)

type DipperinVenusApi struct {
	service *service.VenusFullChainService
}

// verify whether the chain is in sync
func (api *DipperinVenusApi) GetSyncStatus() bool {
	return api.service.GetSyncStatus()
}

// swagger:operation GET /url/CurrentBlock block information block
// ---
// summary: get the current block
// description: get the current block
// produces:
// - application/json
// responses:
//   "200":
//        "$ref": "#/responses/BlockResp"
func (api *DipperinVenusApi) CurrentBlock() (*BlockResp, error) {

	blockResp := &BlockResp{
		Header: model.Header{},
		Body: model.Body{
			Txs:  make([]*model.Transaction, 0),
			Vers: []model.AbstractVerification{},
		},
	}

	curBlock := api.service.CurrentBlock()

	blockResp.Header = *curBlock.Header().(*model.Header)
	blockResp.Body = *curBlock.Body().(*model.Body)

	//	log.Debug("the blockResp header is: ","header",blockResp.Header)
	log.Debug("the blockResp body is: ", "body", blockResp.Body)

	//log.Debug("the blockResp transactions is:","txs",*blockResp.Body.Txs[0])

	return blockResp, nil
}

// swagger:operation POST /url/GetBlockByNumber block information block
// ---
// summary: get the block by height
// description: get the block by height
// parameters:
// - name: number
//   in: body
//   description: block height
//   type: integer
//   required: true
// produces:
// - application/json
// responses:
//   "200":
//        "$ref": "#/responses/BlockResp"
func (api *DipperinVenusApi) GetBlockByNumber(number uint64) (*BlockResp, error) {

	blockResp := &BlockResp{
		Header: model.Header{},
		Body: model.Body{
			Txs:  make([]*model.Transaction, 0),
			Vers: []model.AbstractVerification{},
		},
	}

	curBlock, err := api.service.GetBlockByNumber(number)
	log.Info("DipperinVenusApi#GetBlockByNumber", "curBlock", curBlock)
	if err != nil || curBlock == nil {
		return nil, g_error.ErrBlockNotFound
	}

	//	log.Debug("the current block is: ","current block",*curBlock.(*model.Block))

	blockResp.Header = *curBlock.Header().(*model.Header)
	blockResp.Body = *curBlock.Body().(*model.Body)

	//	log.Debug("the blockResp header is: ","header",blockResp.Header)
	log.Debug("the blockResp body is: ", "body", blockResp.Body)

	return blockResp, nil
}

// swagger:operation POST /url/GetBlockByHash block information block
// ---
// summary: get the block by block hash
// description: get the block by block hash
// parameters:
// - name: hash
//   in: body
//   description: block hash
//   type: common.Hash
//   required: true
// produces:
// - application/json
// responses:
//   "200":
//        "$ref": "#/responses/BlockResp"
func (api *DipperinVenusApi) GetBlockByHash(hash common.Hash) (*BlockResp, error) {
	blockResp := &BlockResp{
		Header: model.Header{},
		Body: model.Body{
			Txs:  make([]*model.Transaction, 0),
			Vers: []model.AbstractVerification{},
		},
	}

	curBlock, err := api.service.GetBlockByHash(hash)
	if err != nil {
		return nil, err
	} else if curBlock == nil {
		return nil, errors.New(fmt.Sprintf("no block hash is %s", hash))
	}

	blockResp.Header = *curBlock.Header().(*model.Header)
	blockResp.Body = *curBlock.Body().(*model.Body)

	return blockResp, nil
}

// swagger:operation POST /url/GetBlockNumber block information block
// ---
// summary: get the height of the block by block hash
// description: get block height by block hash
// parameters:
// - name: hash
//   in: body
//   description: block hash
//   type: common.Hash
//   required: true
// produces:
// - application/json
// responses:
//   "200":
//        description: block height
func (api *DipperinVenusApi) GetBlockNumber(hash common.Hash) *uint64 {
	return api.service.GetBlockNumber(hash)
}

// get genesis block
// swagger:operation GET /url/GetGenesis block information block
// ---
// summary: get genesis block
// description: get genesis block
// produces:
// - application/json
// responses:
//   "200":
//        "$ref": "#/responses/BlockResp"
func (api *DipperinVenusApi) GetGenesis() (*BlockResp, error) {
	blockResp := &BlockResp{
		Header: model.Header{},
		Body: model.Body{
			Txs:  make([]*model.Transaction, 0),
			Vers: []model.AbstractVerification{},
		},
	}

	curBlock, err := api.service.GetGenesis()
	if err != nil {
		return nil, err
	}

	//	log.Debug("the current block is: ","current block",*curBlock.(*model.Block))

	blockResp.Header = *curBlock.Header().(*model.Header)
	blockResp.Body = *curBlock.Body().(*model.Body)

	//	log.Debug("the blockResp header is: ","header",blockResp.Header)
	log.Debug("the blockResp body is: ", "body", blockResp.Body)

	//	log.Debug("the blockResp transactions is:","txs",*blockResp.Body.Txs[0])

	return blockResp, nil
}

// consult block body by block hash
// swagger:operation POST /url/GetBlockBody block information block
// ---
// summary: consult block body by block hash
// description: consult block body by block hash
// parameters:
// - name: hash
//   in: body
//   description: block hash
//   type: common.Hash
//   required: true
// produces:
// - application/json
// responses:
//   "200":
//        "$ref": "#/responses/Body"
func (api *DipperinVenusApi) GetBlockBody(hash common.Hash) *model.Body {
	tmpBody := api.service.GetBlockBody(hash)

	return tmpBody.(*model.Body)
}

// get the current account balance
// swagger:operation POST /url/CurrentBalance account information CurBalanceResp
// ---
// summary: get all UTXO balances of current user
// description: get all UTXO balances of current user
// parameters:
// - name: addresses
//   in: body
//   description: address
//   type: []common.Address
//   required: true
// produces:
// - application/json
// responses:
//   "200":
//        "$ref": "#/responses/CurBalanceResp"
func (api *DipperinVenusApi) CurrentBalance(address common.Address) (resp *CurBalanceResp, err error) {
	balance := api.service.CurrentBalance(address)
	return &CurBalanceResp{
		Balance: (*hexutil.Big)(balance),
	}, nil
}

// fetch transaction data from TxID
// swagger:operation POST /url/Transaction transaction information TransactionReq
// ---
// summary: fetch transaction data, the height and block ID from TxID
// description: fetch transaction data, the height and block ID from TxID
// parameters:
// - name: hash
//   in: body
//   description: transaction hash
//   type: common.Hash
//   required: true
// produces:
// - application/json
// responses:
//   "200":
//        "$ref": "#/responses/TransactionResp"
func (api *DipperinVenusApi) Transaction(hash common.Hash) (resp *TransactionResp, err error) {

	tmpResp := TransactionResp{
		Transaction: &model.Transaction{},
		BlockHash:   common.Hash{},
		BlockNumber: 0,
		TxIndex:     0,
	}
	tmpResp.Transaction, tmpResp.BlockHash, tmpResp.BlockNumber, tmpResp.TxIndex, err = api.service.Transaction(hash)
	if err != nil {
		return nil, err
	}

	//log.Info("the resp.Transaction is: ","tx",tmpResp.Transaction)
	/*	log.Info("the resp.BlockHash is: ","blockHash",tmpResp.BlockHash)
		log.Info("the resp.BlockNumber is: ","blockNum",tmpResp.BlockNumber)
		log.Info("the resp.TxIndex is: ","txIndex",tmpResp.TxIndex)*/

	return &tmpResp, nil
}

// get the nonce needed for the transaction:
// swagger:operation POST /url/Transaction transaction information Transaction
// ---
// summary: get the nonce needed for the transaction
// description: get the nonce needed for the transaction
// parameters:
// - name: addr
//   in: body
//   description: address
//   type: common.Address
//   required: true
// produces:
// - application/json
// responses:
//   "200":
//        description: return nonce and the result
func (api *DipperinVenusApi) GetTransactionNonce(addr common.Address) (nonce uint64, err error) {
	return api.service.GetTransactionNonce(addr)
}

// create a new transaction Tx:
// swagger:operation POST /url/Transaction transaction information NewTransactionReq
// ---
// summary: create a new transaction Tx
// description: create a new transaction Tx
// parameters:
// - name: tra
//   in: body
//   description: address
//   type: model.Transaction
//   required: true
// produces:
// - application/json
// responses:
//   "200":
//        description:return TxHash and the operation result
func (api *DipperinVenusApi) NewTransaction(transactionRlpB []byte) (TxHash common.Hash, err error) {
	var transaction model.Transaction
	err = rlp.DecodeBytes(transactionRlpB, &transaction)
	if err != nil {
		log.TagError("decode client tx failed", "err", err)
		return common.Hash{}, err
	}

	log.Info("[NewTransaction] the tx is: ", "tx", transaction)
	if TxHash, err = api.service.NewTransaction(transaction); err != nil {
		log.TagError("add and broadcast client tx failed", "err", err)
	}

	log.Info("NewTransaction the txId is:", "txId", TxHash.Hex())
	return
}

// call contract
func (api *DipperinVenusApi) NewContract(transactionRlpB []byte, blockNum uint64) (resp string, err error) {
	var transaction model.Transaction
	err = rlp.DecodeBytes(transactionRlpB, &transaction)
	if err != nil {
		log.TagError("decode client tx failed", "err", err)
		return "", err
	}

	curBlock := api.service.CurrentBlock()
	if blockNum == 0 || curBlock.Number() < blockNum {
		blockNum = curBlock.Number()
	}
	return api.service.Call(&transaction, blockNum)
}

func (api *DipperinVenusApi) NewEstimateGas(transactionRlpB []byte) (resp hexutil.Uint64, err error) {
	var transaction model.Transaction
	err = rlp.DecodeBytes(transactionRlpB, &transaction)
	if err != nil {
		log.TagError("decode client tx failed", "err", err)
		return hexutil.Uint64(0), err
	}

	blockNum := api.service.CurrentBlock().Number()
	return api.service.EstimateGas(&transaction, blockNum)
}

//func (apiB *DipperinVenusApi) RetrieveSingleSC(req *req_params.SingleSCReq) *req_params.RetrieveSingleSCResp {
//	cslog.Info().Msg("fetch a single contract")
//	return &req_params.RetrieveSingleSCResp{
//		BaseResp: req_params.NewBaseRespWithErr(nil, "Sent successfully"),
//	}
//}
//
//func (apiB *DipperinVenusApi) NewSC(req *req_params.NewSCReq) *req_params.NewSCResp {
//	cslog.Info().Msg("create a new smart contract")
//	return &req_params.NewSCResp{
//		BaseResp: req_params.NewBaseRespWithErr(nil, "Sent successfully"),
//	}
//}
//
//func (apiB *DipperinVenusApi) GetSCConfig(req *req_params.SingleSCReq) *req_params.GetSCConfigResp {
//	cslog.Info().Msg("get contract settings")
//	return &req_params.GetSCConfigResp{
//		BaseResp:       req_params.NewBaseRespWithErr(nil, "Sent successfully"),
//		PagingRespBase: req_params.PagingRespBase{TotalCount: 103, TotalPages: 13},
//	}
//}

//func (apiB *DipperinVenusApi) GetBlockInfo(req *req_params.GetBlockInfoReq) *req_params.GetBlockInfoResp {
//	cslog.Info().Msg("get contract settings")
//	return &req_params.GetBlockInfoResp{
//		BaseResp:       req_params.NewBaseRespWithErr(nil, "Sent successfully"),
//		PagingRespBase: req_params.PagingRespBase{TotalCount: 104, TotalPages: 14},
//	}
//}

//func (apiB *DipperinVenusApi) GetContractInfo(eData *contract.ExtraDataForContract) (interface{}, error) {
//	return apiB.dipperin.GetContractInfo(eData)
//}
//

//func (apiB *DipperinVenusApi) GetContract(contractAddr common.Address) (interface{}, error) {
//	return apiB.dipperin.GetContract(contractAddr)
//}

// set mining address:
// swagger:operation POST /url/SetMineCoinBase mineOperation SetMineCoinBase
// ---
// summary: set mine CoinBase address
// description: set mine CoinBase address
// parameters:
// - name: addr
//   in: body
//   description: address
//   type: common.Address
//   required: true
// produces:
// - application/json
// responses:
//   "200":
//        description: return the operation result
func (api *DipperinVenusApi) SetMineCoinBase(addr common.Address) error {
	return api.service.SetMineCoinBase(addr)
}

func (api *DipperinVenusApi) GetMineCoinBase() common.Address {
	return api.service.MineMaster.CurrentCoinbaseAddress()
}

func (api *DipperinVenusApi) SetMineGasConfig(gasFloor, gasCeil uint64) error {
	return api.service.SetMineGasConfig(gasFloor, gasCeil)
}

// start mine:
// swagger:operation POST /url/StartMine mineOperation StartMine
// ---
// summary: start mine
// description: start mine
// produces:
// - application/json
// responses:
//   "200":
//        description: return the operation result
func (api *DipperinVenusApi) StartMine() error {
	return api.service.StartMine()
}

// stop mine:
// swagger:operation POST /url/StopMine mineOperation StopMine
// ---
// summary: stop mine
// description: stop mine
// produces:
// - application/json
// responses:
//   "200":
//        description: return the operation result
func (api *DipperinVenusApi) StopMine() error {
	return api.service.StopMine()
}

// establish wallet
// swagger:operation POST /url/EstablishWallet WalletOperation Wallet
// ---
// summary: establish wallet
// description: establish wallet
// parameters:
// - name: walletIdentifier
//   in: body
//   description: wallet identifier
//   type: accounts.WalletIdentifier
//   required: true
// - name: password
//   in: body
//   description: wallet password
//   type: string
//   required: true
// - name: passPhrase
//   in: body
//   description: wallet mnemonic passPhrase
//   type: string
//   required: true
// produces:
// - application/json
// responses:
//   "200":
//        description: return mnemonic and the operation result
func (api *DipperinVenusApi) EstablishWallet(password, passPhrase string, walletIdentifier accounts.WalletIdentifier) (mnemonic string, err error) {
	return api.service.EstablishWallet(walletIdentifier, password, passPhrase)
}

// open wallet
// swagger:operation POST /url/OpenWallet WalletOperation Wallet
// ---
// summary: open wallet
// description: open wallet
// parameters:
// - name: walletIdentifier
//   in: body
//   description: wallet identifier
//   type: accounts.WalletIdentifier
//   required: true
// - name: password
//   in: body
//   description: wallet password
//   type: string
//   required: true
// produces:
// - application/json
// responses:
//   "200":
//        description: return operation result
func (api *DipperinVenusApi) OpenWallet(password string, walletIdentifier accounts.WalletIdentifier) error {
	return api.service.OpenWallet(walletIdentifier, password)
}

// close wallet
// swagger:operation POST /url/CloseWallet WalletOperation Wallet
// ---
// summary: close wallet
// description: close wallet
// parameters:
// - name: walletIdentifier
//   in: body
//   description: wallet identifier
//   type: accounts.WalletIdentifier
//   required: true
// produces:
// - application/json
// responses:
//   "200":
//        description: return operation result
func (api *DipperinVenusApi) CloseWallet(walletIdentifier accounts.WalletIdentifier) error {
	return api.service.CloseWallet(walletIdentifier)
}

// restore wallet
// swagger:operation POST /url/RestoreWallet WalletOperation Wallet
// ---
// summary: restore wallet
// description: restore wallet
// parameters:
// - name: walletIdentifier
//   in: body
//   description: wallet identifier
//   type: accounts.WalletIdentifier
//   required: true
// - name: password
//   in: body
//   description: wallet password
//   type: string
//   required: true
// - name: passPhrase
//   in: body
//   description: wallet mnemonic passPhrase
//   type: string
//   required: true
// - name: mnemonic
//   in: body
//   description: wallet mnemonic
//   type: string
//   required: true
// produces:
// - application/json
// responses:
//   "200":
//        description: return operation result
func (api *DipperinVenusApi) RestoreWallet(password, mnemonic, passPhrase string, walletIdentifier accounts.WalletIdentifier) error {
	return api.service.RestoreWallet(walletIdentifier, password, passPhrase, mnemonic)
}

// list wallet
// swagger:operation POST /url/ListWallet WalletOperation Wallet
// ---
// summary: list wallet
// description: list wallet
// produces:
// - application/json
// responses:
//   "200":
//        description: return wallet identifier list and the operation result
func (api *DipperinVenusApi) ListWallet() ([]accounts.WalletIdentifier, error) {

	walletIdentifier, err := api.service.ListWallet()
	if err != nil {
		return []accounts.WalletIdentifier{}, err
	}

	return walletIdentifier, nil
}

//func (api *DipperinVenusApi) DipperinRpc(v *big.Int) (string, error) {
//    r := fmt.Sprintf("i am dipperin %v", v)
//    return r, nil
//    //return "why", errors.New(r)
//}

func BuildContractExtraData(op string, contractAdr common.Address, params string) []byte {
	erc20 := contract.ExtraDataForContract{
		ContractAddress: contractAdr,
		Action:          op,
		Params:          params,
	}
	erc20Str, _ := json.Marshal(erc20)
	return erc20Str
}

func (api *DipperinVenusApi) GetContractInfo(eData *contract.ExtraDataForContract) (interface{}, error) {
	return api.service.GetContractInfo(eData)
}

func (api *DipperinVenusApi) GetContract(contractAddr common.Address) (interface{}, error) {
	return api.service.GetContract(contractAddr)
}

func (api *DipperinVenusApi) ERC20TotalSupply(contractAddr common.Address) (interface{}, error) {
	extraData := contract.ExtraDataForContract{ContractAddress: contractAddr, Action: "TotalSupply", Params: "[]"}
	return api.service.GetContractInfo(&extraData)
}

func (api *DipperinVenusApi) ERC20Balance(contractAddr, owner common.Address) (interface{}, error) {
	adrStr := fmt.Sprintf("%v", owner)
	params := util.StringifyJson([]interface{}{adrStr})
	extraData := contract.ExtraDataForContract{ContractAddress: contractAddr, Action: "BalanceOf", Params: params}
	return api.service.GetContractInfo(&extraData)
}

func (api *DipperinVenusApi) ERC20Allowance(contractAddr, owner, spender common.Address) (interface{}, error) {
	ownerStr := fmt.Sprintf("%v", owner)
	spenderStr := fmt.Sprintf("%v", spender)
	params := util.StringifyJson([]interface{}{ownerStr, spenderStr})
	extraData := contract.ExtraDataForContract{ContractAddress: contractAddr, Action: "Allowance", Params: params}
	return api.service.GetContractInfo(&extraData)
}

func (api *DipperinVenusApi) ERC20Transfer(contractAddr, from, to common.Address, amount, gasPrice *big.Int, gasLimit uint64) (common.Hash, error) {

	destStr := fmt.Sprintf("%v", to)
	vStr := fmt.Sprintf("0x%x", amount)
	params := util.StringifyJson([]interface{}{destStr, vStr})
	extraData := BuildContractExtraData("Transfer", contractAddr, params)

	//send transaction
	return api.service.SendTransaction(from, contractAddr, big.NewInt(int64(0)), gasPrice, gasLimit, extraData, nil)
}

func (api *DipperinVenusApi) ERC20TransferFrom(contractAdr, owner, from, to common.Address, amount, gasPrice *big.Int, gasLimit uint64) (common.Hash, error) {

	srcStr := fmt.Sprintf("%v", owner)
	destStr := fmt.Sprintf("%v", to)
	vStr := fmt.Sprintf("0x%x", amount)
	params := util.StringifyJson([]interface{}{srcStr, destStr, vStr})
	extraData := BuildContractExtraData("TransferFrom", contractAdr, params)

	//send transaction
	return api.service.SendTransaction(from, contractAdr, big.NewInt(int64(0)), gasPrice, gasLimit, extraData, nil)
}

func (api *DipperinVenusApi) ERC20Approve(contractAdr, from, to common.Address, amount, gasPrice *big.Int, gasLimit uint64) (common.Hash, error) {

	adrStr := fmt.Sprintf("%v", to)
	vStr := fmt.Sprintf("0x%x", amount)
	params := util.StringifyJson([]interface{}{adrStr, vStr})
	extraData := BuildContractExtraData("Approve", contractAdr, params)

	//send transaction
	return api.service.SendTransaction(from, contractAdr, big.NewInt(int64(0)), gasPrice, gasLimit, extraData, nil)
}

func (api *DipperinVenusApi) CreateERC20(from common.Address, tokenName, tokenSymbol string, amount *big.Int, decimal int, gasPrice *big.Int, gasLimit uint64) (ERC20Resp, error) {
	erc20 := contract.BuiltInERC20Token{}
	erc20.Owner = from
	erc20.TokenDecimals = decimal
	erc20.TokenName = tokenName
	erc20.TokenSymbol = tokenSymbol
	erc20.TokenTotalSupply = amount

	es := util.StringifyJson(erc20)

	extra := contract.ExtraDataForContract{}
	extra.Action = "create"
	extra.Params = es
	contractAdr, _ := address_util.GenERC20Address()
	extra.ContractAddress = contractAdr

	txId, err := api.service.SendTransaction(from, contractAdr, big.NewInt(int64(0)), gasPrice, gasLimit, []byte(util.StringifyJson(extra)), nil)
	var resp ERC20Resp
	if err == nil {
		resp.TxId = txId
		resp.CtId = contractAdr
	}

	return resp, err
}

func (api *DipperinVenusApi) CheckBootNode() ([]string, error) {
	nodes := make([]string, len(chain_config.KBucketNodes))
	for i, kn := range chain_config.KBucketNodes {
		nodes[i] = fmt.Sprintf("%s", kn.String())
	}
	return nodes, nil
}

// list wallet account
// swagger:operation POST /url/ListWalletAccount WalletOperation Wallet
// ---
// summary: list wallet account
// description: list wallet account
// parameters:
// - name: walletIdentifier
//   in: body
//   description: wallet identifier
//   type: accounts.WalletIdentifier
//   required: true
// produces:
// - application/json
// responses:
//   "200":
//        description: return account list and the operation result
func (api *DipperinVenusApi) ListWalletAccount(walletIdentifier accounts.WalletIdentifier) ([]accounts.Account, error) {

	tmpAccounts, err := api.service.ListWalletAccount(walletIdentifier)

	/*	for _,account := range tmpAccounts{
		log.Info("the accounts is: ","accounts.Address",account.Address.Hex())
	}*/

	if err != nil {
		return []accounts.Account{}, err
	}

	return tmpAccounts, nil
}

func (api *DipperinVenusApi) StartRemainingService() {
	api.service.StartRemainingService()
}

// set pbft account address
// swagger:operation POST /url/ListWalletAccount PbftAddressOperation pbftAddress
// ---
// summary: set pbft account address
// description: set pbft account address
// parameters:
// - name: walletIdentifier
//   in: body
//   description: wallet identifier
//   type: accounts.WalletIdentifier
//   required: true
// produces:
// - application/json
// responses:
//   "200":
//        description: return account list and the operation result
func (api *DipperinVenusApi) SetBftSigner(address common.Address) error {
	log.Info("DipperinVenusApi SetBftSigner run")
	return api.service.SetBftSigner(address)
}

// wallet add account
// swagger:operation POST /url/AddAccount WalletOperation Wallet
// ---
// summary: wallet add account
// description: add account
// parameters:
// - name: walletIdentifier
//   in: body
//   description: wallet identifier
//   type: accounts.WalletIdentifier
//   required: true
// - name: derivationPath
//   in: body
//   description: the derived path
//   type: string
//   required: true
// produces:
// - application/json
// responses:
//   "200":
//        description: return the added account and the operation result
func (api *DipperinVenusApi) AddAccount(derivationPath string, walletIdentifier accounts.WalletIdentifier) (accounts.Account, error) {
	return api.service.AddAccount(walletIdentifier, derivationPath)
}

// send transaction
// swagger:operation POST /url/SendTransaction transactionOperation transaction
// ---
// summary: send transaction
// description: send transaction
// parameters:
// - name: from
//   in: body
//   description: the address that send coin
//   type: common.Address
//   required: true
// - name: to
//   in: body
//   description: the address that receive coin
//   type: common.Address
//   required: true
// - name: transactionFee
//   in: body
//   description: the transaction fee
//   type: *big.Int
//   required: true
// - name: data
//   in: body
//   description: the transaction extra data
//   type: []byte
//   required: true
// produces:
// - application/json
// responses:
//   "200":
//        description: return operation result
func (api *DipperinVenusApi) SendTransaction(from, to common.Address, value, gasPrice *big.Int, gasLimit uint64, data []byte, nonce *uint64) (common.Hash, error) {
	return api.service.SendTransaction(from, to, value, gasPrice, gasLimit, data, nonce)
}

func (api *DipperinVenusApi) SendTransactionContract(from, to common.Address, value, gasPrice *big.Int, gasLimit uint64, data []byte, nonce *uint64) (common.Hash, error) {
	return api.service.SendTransactionContract(from, to, value, gasPrice, gasLimit, data, nonce)
}

//send multiple-txs
func (api *DipperinVenusApi) SendTransactions(from common.Address, rpcTxs []model.RpcTransaction) (int, error) {
	return api.service.SendTransactions(from, rpcTxs)
}

//new send multiple-txs
func (api *DipperinVenusApi) NewSendTransactions(txs []model.Transaction) (int, error) {
	return api.service.NewSendTransactions(txs)
}

// get remote node height
func (api *DipperinVenusApi) RemoteHeight() uint64 {
	return api.service.RemoteHeight()

}

// send register transaction
// swagger:operation POST /url/SendRegisterTransaction transactionOperation transaction
// ---
// summary: send register transaction
// description: send register transaction
// parameters:
// - name: from
//   in: body
//   description: the address that send register transaction
//   type: common.Address
//   required: true
// - name: stake
//   in: body
//   description: the register pledge
//   type: *big.Int
//   required: true
// - name: fee
//   in: body
//   description: the transaction fee
//   type: *big.Int
//   required: true
// produces:
// - application/json
// responses:
//   "200":
//        description: return operation result
func (api *DipperinVenusApi) SendRegisterTransaction(from common.Address, stake, gasPrice *big.Int, gasLimit uint64, nonce *uint64) (common.Hash, error) {
	return api.service.SendRegisterTransaction(from, stake, gasPrice, gasLimit, nonce)
}

// send unstake transaction
// swagger:operation POST /url/SendCancelTransaction transactionOperation transaction
// ---
// summary: send cancel transaction
// description: send cancel transaction
// parameters:
// - name: from
//   in: body
//   description: the address that send register transaction
//   type: common.Address
//   required: true
// - name: fee
//   in: body
//   description: the transaction fee
//   type: *big.Int
//   required: true
// produces:
// - application/json
// responses:
//   "200":
//        description: return operation result
func (api *DipperinVenusApi) SendUnStakeTransaction(from common.Address, gasPrice *big.Int, gasLimit uint64, nonce *uint64) (common.Hash, error) {
	return api.service.SendUnStakeTransaction(from, gasPrice, gasLimit, nonce)
}

// send evidence transaction
// swagger:operation POST /url/SendEvidenceTransaction transactionOperation transaction
// ---
// summary: send evidence transaction
// description: send evidence transaction
// parameters:
// - name: from
//   in: body
//   description: the address that send register transaction
//   type: common.Address
//   required: true
// - name: target
//   in: body
//   description: the report target
//   type: common.Address
//   required: true
// - name: fee
//   in: body
//   description: the transaction fee
//   type: *big.Int
//   required: true
// - name: voteA
//   in: body
//   description: the report evidence voteA
//   type: *model.Verification
//   required: true
// - name: voteB
//   in: body
//   description: the report evidence voteB
//   type: *model.Verification
//   required: true
// produces:
// - application/json
// responses:
//   "200":
//        description: return operation result
func (api *DipperinVenusApi) SendEvidenceTransaction(from, target common.Address, gasPrice *big.Int, gasLimit uint64, voteA *model.VoteMsg, voteB *model.VoteMsg, nonce *uint64) (common.Hash, error) {
	return api.service.SendEvidenceTransaction(from, target, gasPrice, gasLimit, voteA, voteB, nonce)
}

// send cancel transaction
// swagger:operation POST /url/SendCancelTransaction transactionOperation transaction
// ---
// summary: send cancel transaction
// description: send cancel transaction
// parameters:
// - name: from
//   in: body
//   description: the address that send register transaction
//   type: common.Address
//   required: true
// - name: fee
//   in: body
//   description: the transaction fee
//   type: *big.Int
//   required: true
// produces:
// - application/json
// responses:
//   "200":
//        description: return operation result
func (api *DipperinVenusApi) SendCancelTransaction(from common.Address, gasPrice *big.Int, gasLimit uint64, nonce *uint64) (common.Hash, error) {
	return api.service.SendCancelTransaction(from, gasPrice, gasLimit, nonce)
}

// get verifiers info by round
// swagger:operation POST /url/GetVerifiersBySlot verifierInfo verifierInfo
// ---
// summary: get verifiers info by round
// description: get verifiers info by round
// parameters:
// - name: slotNum
//   in: body
//   description: the round (slotNum = blockNumber/slotSize)
//   type: uint64
//   required: true
// produces:
// - application/json
// responses:
//   "200":
//        description: return verifier address list and the operation result
func (api *DipperinVenusApi) GetVerifiersBySlot(slotNum uint64) ([]common.Address, error) {
	return api.service.GetVerifiers(slotNum), nil
}

func (api *DipperinVenusApi) GetSlotByNumber(blockNum uint64) (uint64, error) {
	block, _ := api.service.GetBlockByNumber(blockNum)
	if block == nil {
		return uint64(0), errors.New("invalid block number")
	}
	slot := api.service.GetSlot(block)
	return *slot, nil
}

// get current verifiers
// swagger:operation POST /url/GetCurVerifiers verifierInfo verifierInfo
// ---
// summary: get current verifiers
// description: get current verifiers
// produces:
// - application/json
// responses:
//   "200":
//        description: return verifier address list and the operation result
func (api *DipperinVenusApi) GetCurVerifiers() []common.Address {
	return api.service.GetCurVerifiers()
}

// get next verifiers
// swagger:operation POST /url/GetNextVerifiers verifierInfo verifierInfo
// ---
// summary: get next verifiers
// description: get next verifiers
// produces:
// - application/json
// responses:
//   "200":
//        description: return verifier address list and the operation result
func (api *DipperinVenusApi) GetNextVerifiers() []common.Address {
	return api.service.GetNextVerifiers()
}

// verifier status
// swagger:operation POST /url/VerifierStatus verifierInfo verifierInfo
// ---
// summary: verifier status
// description: verifier status
// produces:
// - application/json
// responses:
//   "200":
//        description: return verifier status and the operation result
func (api *DipperinVenusApi) VerifierStatus(address common.Address) (resp *VerifierStatus, err error) {
	state, stake, balance, reputation, isCurrentVerifier, err := api.service.VerifierStatus(address)
	return &VerifierStatus{
		Status:            state,
		Stake:             (*hexutil.Big)(stake),
		Balance:           (*hexutil.Big)(balance),
		Reputation:        reputation,
		IsCurrentVerifier: isCurrentVerifier,
	}, err
}

// get address stake
// swagger:operation POST /url/CurrentStake stakeInfo stakeInfo
// ---
// summary: get address stake
// description: get address stake
// produces:
// - application/json
// responses:
//   "200":
//        description: return address stake and the operation result
func (api *DipperinVenusApi) CurrentStake(address common.Address) (resp *CurStakeResp, err error) {
	stake := api.service.CurrentStake(address)
	return &CurStakeResp{
		Stake: (*hexutil.Big)(stake),
	}, nil
}

func (api *DipperinVenusApi) CurrentReputation(address common.Address) (uint64, error) {
	return api.service.CurrentReputation(address)
}

//get current practical verifiers
func (api *DipperinVenusApi) GetCurrentConnectPeers() ([]PeerInfoResp, error) {
	peersInfo := make([]PeerInfoResp, 0)
	tmpInfo := api.service.GetCurrentConnectPeers()

	for nodeId, address := range tmpInfo {
		verifier := PeerInfoResp{
			NodeId:  nodeId,
			Address: address,
		}
		peersInfo = append(peersInfo, verifier)
	}
	return peersInfo, nil
}

/*//SyncUsedAccounts
func (api *DipperinVenusApi)SyncUsedAccounts(walletIdentifier accounts.WalletIdentifier,MaxChangeValue ,MaxIndex uint32) error{
	return api.service.SyncUsedAccounts(walletIdentifier,MaxChangeValue,MaxIndex)
}*/

//get address nonce from wallet
func (api *DipperinVenusApi) GetAddressNonceFromWallet(address common.Address) (nonce uint64, err error) {
	return api.service.GetAddressNonceFromWallet(address)
}

func (api *DipperinVenusApi) GetChainConfig() (conf chain_config.ChainConfig, err error) {
	return api.service.GetChainConfig(), nil
}

func (api *DipperinVenusApi) GetBlockDiffVerifierInfo(blockNumber uint64) (map[economy_model.VerifierType][]common.Address, error) {
	return api.service.GetBlockDiffVerifierInfo(blockNumber)
}

func (api *DipperinVenusApi) GetVerifierDIPReward(blockNumber uint64) (map[economy_model.VerifierType]*hexutil.Big, error) {
	reward, err := api.service.GetVerifierDIPReward(blockNumber)
	if err != nil {
		return map[economy_model.VerifierType]*hexutil.Big{}, err
	}

	result := make(map[economy_model.VerifierType]*hexutil.Big, 0)
	for key, value := range reward {
		result[key] = (*hexutil.Big)(value)
	}
	return result, nil
}

func (api *DipperinVenusApi) GetMineMasterDIPReward(blockNumber uint64) (*hexutil.Big, error) {
	reward, err := api.service.GetMineMasterDIPReward(blockNumber)
	if err != nil {
		return nil, nil
	}

	return (*hexutil.Big)(reward), nil
}

func (api *DipperinVenusApi) GetBlockYear(blockNumber uint64) (uint64, error) {
	return api.service.GetBlockYear(blockNumber)
}

func (api *DipperinVenusApi) GetOneBlockTotalDIPReward(blockNumber uint64) (*hexutil.Big, error) {
	reward, err := api.service.GetOneBlockTotalDIPReward(blockNumber)
	if err != nil {
		return nil, err
	}

	return (*hexutil.Big)(reward), nil
}

func (api *DipperinVenusApi) GetInvestorInfo() map[string]*hexutil.Big {
	investorInfo := api.service.GetInvestorInfo()

	result := make(map[string]*hexutil.Big, 0)
	for key, value := range investorInfo {
		result[key.Hex()] = (*hexutil.Big)(value)
	}
	return result
}

func (api *DipperinVenusApi) GetDeveloperInfo() map[string]*hexutil.Big {
	developerInfo := api.service.GetDeveloperInfo()

	result := make(map[string]*hexutil.Big, 0)
	for key, value := range developerInfo {
		result[key.Hex()] = (*hexutil.Big)(value)
	}
	return result
}

func (api *DipperinVenusApi) GetAddressLockMoney(address common.Address) (*big.Int, error) {
	return api.service.GetAddressLockMoney(address)
}

func (api *DipperinVenusApi) GetInvestorLockDIP(address common.Address, blockNumber uint64) (*hexutil.Big, error) {
	lockValue, err := api.service.GetInvestorLockDIP(address, blockNumber)
	if err != nil {
		return nil, err
	}

	return (*hexutil.Big)(lockValue), nil
}

func (api *DipperinVenusApi) GetDeveloperLockDIP(address common.Address, blockNumber uint64) (*hexutil.Big, error) {
	lockValue, err := api.service.GetDeveloperLockDIP(address, blockNumber)
	if err != nil {
		return nil, err
	}

	return (*hexutil.Big)(lockValue), nil
}

func (api *DipperinVenusApi) GetFoundationInfo(usage economy_model.FoundationDIPUsage) map[string]*hexutil.Big {
	foundationInfo := api.service.GetFoundationInfo(usage)

	result := make(map[string]*hexutil.Big, 0)
	for key, value := range foundationInfo {
		result[key.Hex()] = (*hexutil.Big)(value)
	}
	return result
}

func (api *DipperinVenusApi) GetMaintenanceLockDIP(address common.Address, blockNumber uint64) (*hexutil.Big, error) {
	lockValue, err := api.service.GetMaintenanceLockDIP(address, blockNumber)
	if err != nil {
		return nil, err
	}

	return (*hexutil.Big)(lockValue), nil
}

func (api *DipperinVenusApi) GetReMainRewardLockDIP(address common.Address, blockNumber uint64) (*hexutil.Big, error) {
	lockValue, err := api.service.GetReMainRewardLockDIP(address, blockNumber)
	if err != nil {
		return nil, err
	}

	return (*hexutil.Big)(lockValue), nil
}

func (api *DipperinVenusApi) GetEarlyTokenLockDIP(address common.Address, blockNumber uint64) (*hexutil.Big, error) {
	lockValue, err := api.service.GetEarlyTokenLockDIP(address, blockNumber)
	if err != nil {
		return nil, err
	}

	return (*hexutil.Big)(lockValue), nil
}

func (api *DipperinVenusApi) GetMineMasterEDIPReward(blockNumber uint64, tokenDecimals int) (*hexutil.Big, error) {
	reward, err := api.service.GetMineMasterEDIPReward(blockNumber, tokenDecimals)
	if err != nil {
		return nil, nil
	}

	return (*hexutil.Big)(reward), nil
}

func (api *DipperinVenusApi) GetVerifierEDIPReward(blockNumber uint64, tokenDecimals int) (map[economy_model.VerifierType]*hexutil.Big, error) {
	reward, err := api.service.GetVerifierEDIPReward(blockNumber, tokenDecimals)
	if err != nil {
		return map[economy_model.VerifierType]*hexutil.Big{}, err
	}

	result := make(map[economy_model.VerifierType]*hexutil.Big, 0)
	for key, value := range reward {
		result[key] = (*hexutil.Big)(value)
	}
	return result, nil
}

func (api *DipperinVenusApi) NewBlock(ctx context.Context) (*rpc.Subscription, error) {
	return api.service.NewBlock(ctx)
}

func (api *DipperinVenusApi) SubscribeBlock(ctx context.Context) (*rpc.Subscription, error) {
	return api.service.SubscribeBlock(ctx)
}

func (api *DipperinVenusApi) StopDipperin() {
	api.service.StopDipperin()
}

func (api *DipperinVenusApi) GetABI(contractAddr common.Address) (*utils.WasmAbi, error) {
	return api.service.GetABI(contractAddr)
}

func (api *DipperinVenusApi) GetCode(contractAddr common.Address) ([]byte, error) {
	return api.service.GetCode(contractAddr)
}

func (api *DipperinVenusApi) SuggestGasPrice() (resp *CurBalanceResp, err error) {
	gasPrice, err := api.service.SuggestGasPrice()
	if err != nil {
		return nil, err
	}
	return &CurBalanceResp{
		Balance: (*hexutil.Big)(gasPrice),
	}, nil
}

func (api *DipperinVenusApi) GetContractAddressByTxHash(txHash common.Hash) (common.Address, error) {
	return api.service.GetContractAddressByTxHash(txHash)
}

func (api *DipperinVenusApi) GetLogs(blockHash common.Hash, fromBlock, toBlock uint64, Addresses []common.Address, Topics [][]common.Hash) ([]*model2.Log, error) {
	return api.service.GetLogs(blockHash, fromBlock, toBlock, Addresses, Topics)
}

//func (api *DipperinVenusApi) GetLogs(blockHash common.Hash, fromBlock, toBlock uint64, Addresses []common.Address, Topics [][]string) ([]*model2.Log, error) {
//	var tps [][]common.Hash
//	for _,ts  := range  Topics{
//		var tp []common.Hash
//		for _, t := range ts {
//			tp = append(tp, common.BytesToHash(crypto.Keccak256([]byte(t))))
//		}
//		tps = append(tps, tp)
//	}
//	return api.service.GetLogs(blockHash, fromBlock, toBlock, Addresses, tps)
//}

func (api *DipperinVenusApi) GetTxActualFee(txHash common.Hash) (resp *CurBalanceResp, err error) {
	fee, err := api.service.GetTxActualFee(txHash)
	if err != nil {
		return nil, err
	}
	return &CurBalanceResp{
		Balance: (*hexutil.Big)(fee),
	}, nil
}

func (api *DipperinVenusApi) GetReceiptByTxHash(txHash common.Hash) (*model2.Receipt, error) {
	return api.service.GetReceiptByTxHash(txHash)
}

func (api *DipperinVenusApi) GetReceiptsByBlockNum(num uint64) (model2.Receipts, error) {
	return api.service.GetReceiptsByBlockNum(num)
}

func (api *DipperinVenusApi) CallContract(from, to common.Address, data []byte, blockNum uint64) (string, error) {
	extraData, err := api.service.GetExtraData(to, data)
	if err != nil {
		return "", err
	}
	args := service.CallArgs{
		From: from,
		To:   &to,
		Data: extraData,
	}

	var gasLimit uint64
	curBlock := api.service.CurrentBlock()
	if blockNum == 0 || curBlock.Number() < blockNum {
		blockNum = curBlock.Number()
		gasLimit = curBlock.Header().GetGasLimit()
	} else {
		block, _ := api.service.GetBlockByNumber(blockNum)
		gasLimit = block.Header().GetGasLimit()
	}
	args.Gas = hexutil.Uint64(gasLimit)
	log.Info("API#CallContract start", "from", from, "to", to, "blockNum", blockNum)
	log.Info("API#CallContract start", "gasLimit", gasLimit)
	signedTx, err := api.service.MakeTmpSignedTx(args, blockNum)
	if err != nil {
		return "", err
	}
	return api.service.Call(signedTx, blockNum)
}

func (api *DipperinVenusApi) EstimateGas(from, to common.Address, value, gasPrice *big.Int, gasLimit uint64, data []byte, nonce *uint64) (hexutil.Uint64, error) {
	if value == nil {
		value = new(big.Int).SetUint64(0)
	}

	if gasPrice == nil {
		gasPrice = big.NewInt(0).SetInt64(config.DEFAULT_GAS_PRICE)
	}

	extraData, err := api.service.GetExtraData(to, data)
	if err != nil {
		return hexutil.Uint64(0), err
	}

	args := service.CallArgs{
		From:     from,
		To:       &to,
		Gas:      hexutil.Uint64(gasLimit),
		GasPrice: hexutil.Big(*gasPrice),
		Value:    hexutil.Big(*value),
		Data:     hexutil.Bytes(extraData),
	}

	blockNum := api.service.CurrentBlock().Number()
	log.Info("API#EstimateGas start", "from", from, "to", to, "blockNum", blockNum)
	log.Info("API#EstimateGas start", "value", value, "gasPrice", gasPrice, "gasLimit", gasLimit)
	signedTx, err := api.service.MakeTmpSignedTx(args, blockNum)
	if err != nil {
		return hexutil.Uint64(0), err
	}
	return api.service.EstimateGas(signedTx, blockNum)
}
