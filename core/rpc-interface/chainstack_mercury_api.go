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
    "github.com/dipperin/dipperin-core/common"
    "github.com/dipperin/dipperin-core/common/hexutil"
    "github.com/dipperin/dipperin-core/core/accounts"
    "github.com/dipperin/dipperin-core/core/chain-config"
    "github.com/dipperin/dipperin-core/core/contract"
    "github.com/dipperin/dipperin-core/core/economy-model"
    "github.com/dipperin/dipperin-core/core/model"
    "github.com/dipperin/dipperin-core/third-party/log"
    "github.com/dipperin/dipperin-core/common/util"
    "context"
    "fmt"
    "github.com/ethereum/go-ethereum/rlp"
    "github.com/dipperin/dipperin-core/third-party/rpc"
    "math/big"
    "github.com/dipperin/dipperin-core/common/address-util"
    "encoding/json"
    "github.com/dipperin/dipperin-core/core/dipperin/service"
)

type DipperinMercuryApi struct {
    service *service.MercuryFullChainService
}

// verify whether the chain is in sync
func (api *DipperinMercuryApi) GetSyncStatus() bool {
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
func (api *DipperinMercuryApi) CurrentBlock() (*BlockResp, error) {

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
func (api *DipperinMercuryApi) GetBlockByNumber(number uint64) (*BlockResp, error) {

    blockResp := &BlockResp{
        Header: model.Header{},
        Body: model.Body{
            Txs:  make([]*model.Transaction, 0),
            Vers: []model.AbstractVerification{},
        },
    }

    curBlock, err := api.service.GetBlockByNumber(number)
    if err != nil {
        return nil, err
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
func (api *DipperinMercuryApi) GetBlockByHash(hash common.Hash) (*BlockResp, error) {
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
    }

    //	log.Debug("the current block is: ","current block",*curBlock.(*model.Block))

    blockResp.Header = *curBlock.Header().(*model.Header)
    blockResp.Body = *curBlock.Body().(*model.Body)

    //	log.Debug("the blockResp header is: ","header",blockResp.Header)
    log.Debug("the blockResp body is: ", "body", blockResp.Body)

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
func (api *DipperinMercuryApi) GetBlockNumber(hash common.Hash) *uint64 {
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
func (api *DipperinMercuryApi) GetGenesis() (*BlockResp, error) {
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
func (api *DipperinMercuryApi) GetBlockBody(hash common.Hash) *model.Body {
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
func (api *DipperinMercuryApi) CurrentBalance(address common.Address) (resp *CurBalanceResp, err error) {
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
func (api *DipperinMercuryApi) Transaction(hash common.Hash) (resp *TransactionResp, err error) {

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
func (api *DipperinMercuryApi) GetTransactionNonce(addr common.Address) (nonce uint64, err error) {
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
func (api *DipperinMercuryApi) NewTransaction(transactionRlpB []byte) (TxHash common.Hash, err error) {
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

//func (apiB *DipperinMercuryApi) RetrieveSingleSC(req *req_params.SingleSCReq) *req_params.RetrieveSingleSCResp {
//	cslog.Info().Msg("fetch a single contract")
//	return &req_params.RetrieveSingleSCResp{
//		BaseResp: req_params.NewBaseRespWithErr(nil, "Sent successfully"),
//	}
//}
//
//func (apiB *DipperinMercuryApi) NewSC(req *req_params.NewSCReq) *req_params.NewSCResp {
//	cslog.Info().Msg("create a new smart contract")
//	return &req_params.NewSCResp{
//		BaseResp: req_params.NewBaseRespWithErr(nil, "Sent successfully"),
//	}
//}
//
//func (apiB *DipperinMercuryApi) GetSCConfig(req *req_params.SingleSCReq) *req_params.GetSCConfigResp {
//	cslog.Info().Msg("get contract settings")
//	return &req_params.GetSCConfigResp{
//		BaseResp:       req_params.NewBaseRespWithErr(nil, "Sent successfully"),
//		PagingRespBase: req_params.PagingRespBase{TotalCount: 103, TotalPages: 13},
//	}
//}

//func (apiB *DipperinMercuryApi) GetBlockInfo(req *req_params.GetBlockInfoReq) *req_params.GetBlockInfoResp {
//	cslog.Info().Msg("get contract settings")
//	return &req_params.GetBlockInfoResp{
//		BaseResp:       req_params.NewBaseRespWithErr(nil, "Sent successfully"),
//		PagingRespBase: req_params.PagingRespBase{TotalCount: 104, TotalPages: 14},
//	}
//}


//func (apiB *DipperinMercuryApi) GetContractInfo(eData *contract.ExtraDataForContract) (interface{}, error) {
//	return apiB.dipperin.GetContractInfo(eData)
//}
//

//func (apiB *DipperinMercuryApi) GetContract(contractAddr common.Address) (interface{}, error) {
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
func (api *DipperinMercuryApi) SetMineCoinBase(addr common.Address) error {
    return api.service.SetMineCoinBase(addr)
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
func (api *DipperinMercuryApi) StartMine() error {
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
func (api *DipperinMercuryApi) StopMine() error {
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
func (api *DipperinMercuryApi) EstablishWallet(password, passPhrase string, walletIdentifier accounts.WalletIdentifier) (mnemonic string, err error) {
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
func (api *DipperinMercuryApi) OpenWallet(password string, walletIdentifier accounts.WalletIdentifier) error {
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
func (api *DipperinMercuryApi) CloseWallet(walletIdentifier accounts.WalletIdentifier) error {
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
func (api *DipperinMercuryApi) RestoreWallet(password, mnemonic, passPhrase string, walletIdentifier accounts.WalletIdentifier) error {
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
func (api *DipperinMercuryApi) ListWallet() ([]accounts.WalletIdentifier, error) {

    walletIdentifier, err := api.service.ListWallet()
    if err != nil {
        return []accounts.WalletIdentifier{}, err
    }

    return walletIdentifier, nil
}

//func (api *DipperinMercuryApi) DipperinRpc(v *big.Int) (string, error) {
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

func (api *DipperinMercuryApi) GetContractInfo(eData *contract.ExtraDataForContract) (interface{}, error) {
    return api.service.GetContractInfo(eData)
}

func (api *DipperinMercuryApi) GetContract(contractAddr common.Address) (interface{}, error) {
    return api.service.GetContract(contractAddr)
}

func (api *DipperinMercuryApi) ERC20TotalSupply(contractAddr common.Address) (interface{}, error) {
    extraData := contract.ExtraDataForContract{ContractAddress: contractAddr, Action: "TotalSupply", Params: "[]"}
    return api.service.GetContractInfo(&extraData)
}

func (api *DipperinMercuryApi) ERC20Balance(contractAddr, owner common.Address) (interface{}, error) {
    adrStr := fmt.Sprintf("%v", owner)
    params := util.StringifyJson([]interface{}{adrStr})
    extraData := contract.ExtraDataForContract{ContractAddress: contractAddr, Action: "BalanceOf", Params: params}
    return api.service.GetContractInfo(&extraData)
}

func (api *DipperinMercuryApi) ERC20Allowance(contractAddr, owner, spender common.Address) (interface{}, error) {
    ownerStr := fmt.Sprintf("%v", owner)
    spenderStr := fmt.Sprintf("%v", spender)
    params := util.StringifyJson([]interface{}{ownerStr, spenderStr})
    extraData := contract.ExtraDataForContract{ContractAddress: contractAddr, Action: "Allowance", Params: params}
    return api.service.GetContractInfo(&extraData)
}

func (api *DipperinMercuryApi) ERC20Transfer(contractAddr, from, to common.Address, amount, txFee *big.Int) (common.Hash, error) {

    destStr := fmt.Sprintf("%v", to)
    vStr := fmt.Sprintf("0x%x", amount)
    params := util.StringifyJson([]interface{}{destStr, vStr})
    extraData := BuildContractExtraData("Transfer", contractAddr, params)

    //send transaction
    return api.service.SendTransaction(from, contractAddr, big.NewInt(int64(0)), txFee, extraData, nil)
}

func (api *DipperinMercuryApi) ERC20TransferFrom(contractAdr, owner, from, to common.Address, amount, txFee *big.Int) (common.Hash, error) {

    srcStr := fmt.Sprintf("%v", owner)
    destStr := fmt.Sprintf("%v", to)
    vStr := fmt.Sprintf("0x%x", amount)
    params := util.StringifyJson([]interface{}{srcStr, destStr, vStr})
    extraData := BuildContractExtraData("TransferFrom", contractAdr, params)

    //send transaction
    return api.service.SendTransaction(from, contractAdr, big.NewInt(int64(0)), txFee, extraData, nil)
}

func (api *DipperinMercuryApi) ERC20Approve(contractAdr, from, to common.Address, amount, txFee *big.Int) (common.Hash, error) {

    adrStr := fmt.Sprintf("%v", to)
    vStr := fmt.Sprintf("0x%x", amount)
    params := util.StringifyJson([]interface{}{adrStr, vStr})
    extraData := BuildContractExtraData("Approve", contractAdr, params)

    //send transaction
    return api.service.SendTransaction(from, contractAdr, big.NewInt(int64(0)), txFee, extraData, nil)
}

func (api *DipperinMercuryApi) CreateERC20(from common.Address, tokenName, tokenSymbol string, amount *big.Int, decimal int, fee *big.Int) (ERC20Resp, error) {
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

    txId, err := api.service.SendTransaction(from, contractAdr, big.NewInt(int64(0)), fee, []byte(util.StringifyJson(extra)), nil)
    var resp ERC20Resp
    if err == nil {
        resp.TxId = txId
        resp.CtId = contractAdr
    }

    return resp, err
}

func (api *DipperinMercuryApi) CheckBootNode() ([]string, error) {
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
func (api *DipperinMercuryApi) ListWalletAccount(walletIdentifier accounts.WalletIdentifier) ([]accounts.Account, error) {

    tmpAccounts, err := api.service.ListWalletAccount(walletIdentifier)

    /*	for _,account := range tmpAccounts{
            log.Info("the accounts is: ","accounts.Address",account.Address.Hex())
        }*/

    if err != nil {
        return []accounts.Account{}, err
    }

    return tmpAccounts, nil
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
func (api *DipperinMercuryApi) SetBftSigner(address common.Address) error {
    log.Info("DipperinMercuryApi SetBftSigner run")
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
func (api *DipperinMercuryApi) AddAccount(derivationPath string, walletIdentifier accounts.WalletIdentifier) (accounts.Account, error) {
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
func (api *DipperinMercuryApi) SendTransaction(from, to common.Address, value, transactionFee *big.Int, data []byte, nonce *uint64) (common.Hash, error) {
    return api.service.SendTransaction(from, to, value, transactionFee, data, nonce)
}

func (api *DipperinMercuryApi) SendTransactionContractCreate(from, to common.Address,value,gasLimit, gasPrice *big.Int, data []byte, nonce *uint64 ) (common.Hash, error) {
    return api.service.SendTransactionContractCreate(from,to, value, gasLimit, gasPrice, data, nonce)
}

//send multiple-txs
func (api *DipperinMercuryApi) SendTransactions(from common.Address, rpcTxs []model.RpcTransaction) (int, error) {
    return api.service.SendTransactions(from, rpcTxs)
}

//new send multiple-txs
func (api *DipperinMercuryApi) NewSendTransactions(txs []model.Transaction) (int, error) {
    return api.service.NewSendTransactions(txs)
}

// get remote node height
func (api *DipperinMercuryApi) RemoteHeight() uint64 {
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
func (api *DipperinMercuryApi) SendRegisterTransaction(from common.Address, stake, fee *big.Int, nonce *uint64) (common.Hash, error) {
    return api.service.SendRegisterTransaction(from, stake, fee, nonce)
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
func (api *DipperinMercuryApi) SendUnStakeTransaction(from common.Address, fee *big.Int, nonce *uint64) (common.Hash, error) {
    return api.service.SendUnStakeTransaction(from, fee, nonce)
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
func (api *DipperinMercuryApi) SendEvidenceTransaction(from, target common.Address, fee *big.Int, voteA *model.VoteMsg, voteB *model.VoteMsg, nonce *uint64) (common.Hash, error) {
    return api.service.SendEvidenceTransaction(from, target, fee, voteA, voteB, nonce)
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
func (api *DipperinMercuryApi) SendCancelTransaction(from common.Address, fee *big.Int, nonce *uint64) (common.Hash, error) {
    return api.service.SendCancelTransaction(from, fee, nonce)
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
func (api *DipperinMercuryApi) GetVerifiersBySlot(slotNum uint64) ([]common.Address, error) {
    return api.service.GetVerifiers(slotNum), nil
}

func (api *DipperinMercuryApi) GetSlot(block model.AbstractBlock) *uint64 {
    return api.service.GetSlot(block)
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
func (api *DipperinMercuryApi) GetCurVerifiers() ([]common.Address) {
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
func (api *DipperinMercuryApi) GetNextVerifiers() ([]common.Address) {
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
func (api *DipperinMercuryApi) VerifierStatus(address common.Address) (resp *VerifierStatus, err error) {
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
func (api *DipperinMercuryApi) CurrentStake(address common.Address) (resp *CurStakeResp, err error) {
    stake := api.service.CurrentStake(address)
    return &CurStakeResp{
        Stake: (*hexutil.Big)(stake),
    }, nil
}

func (api *DipperinMercuryApi) CurrentReputation(address common.Address) (uint64, error) {
    return api.service.CurrentReputation(address)
}

//get current practical verifiers
func (api *DipperinMercuryApi) GetCurrentConnectPeers() ([]PeerInfoResp, error) {
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
func (api *DipperinMercuryApi)SyncUsedAccounts(walletIdentifier accounts.WalletIdentifier,MaxChangeValue ,MaxIndex uint32) error{
	return api.service.SyncUsedAccounts(walletIdentifier,MaxChangeValue,MaxIndex)
}*/

//get address nonce from wallet
func (api *DipperinMercuryApi) GetAddressNonceFromWallet(address common.Address) (nonce uint64, err error) {
    return api.service.GetAddressNonceFromWallet(address)
}

func (api *DipperinMercuryApi) GetChainConfig() (conf chain_config.ChainConfig, err error) {
    return api.service.GetChainConfig(), nil
}

func (api *DipperinMercuryApi) GetBlockDiffVerifierInfo(blockNumber uint64) (map[economy_model.VerifierType][]common.Address, error) {
    return api.service.GetBlockDiffVerifierInfo(blockNumber)
}

func (api *DipperinMercuryApi) GetVerifierDIPReward(blockNumber uint64) (map[economy_model.VerifierType]*hexutil.Big, error) {
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

func (api *DipperinMercuryApi) GetMineMasterDIPReward(blockNumber uint64) (*hexutil.Big, error) {
    reward, err := api.service.GetMineMasterDIPReward(blockNumber)
    if err != nil {
        return nil, nil
    }

    return (*hexutil.Big)(reward), nil
}

func (api *DipperinMercuryApi) GetBlockYear(blockNumber uint64) (uint64, error) {
    return api.service.GetBlockYear(blockNumber)
}

func (api *DipperinMercuryApi) GetOneBlockTotalDIPReward(blockNumber uint64) (*hexutil.Big, error) {
    reward, err := api.service.GetOneBlockTotalDIPReward(blockNumber)
    if err != nil {
        return nil, err
    }

    return (*hexutil.Big)(reward), nil
}

func (api *DipperinMercuryApi) GetInvestorInfo() map[string]*hexutil.Big {
    investorInfo := api.service.GetInvestorInfo()

    result := make(map[string]*hexutil.Big, 0)
    for key, value := range investorInfo {
        result[key.Hex()] = (*hexutil.Big)(value)
    }
    return result
}

func (api *DipperinMercuryApi) GetDeveloperInfo() map[string]*hexutil.Big {
    developerInfo := api.service.GetDeveloperInfo()

    result := make(map[string]*hexutil.Big, 0)
    for key, value := range developerInfo {
        result[key.Hex()] = (*hexutil.Big)(value)
    }
    return result
}

func (api *DipperinMercuryApi) GetAddressLockMoney(address common.Address) (*big.Int,error){
    return api.service.GetAddressLockMoney(address)
}

func (api *DipperinMercuryApi) GetInvestorLockDIP(address common.Address, blockNumber uint64) (*hexutil.Big, error) {
    lockValue, err := api.service.GetInvestorLockDIP(address, blockNumber)
    if err != nil {
        return nil, err
    }

    return (*hexutil.Big)(lockValue), nil
}

func (api *DipperinMercuryApi) GetDeveloperLockDIP(address common.Address, blockNumber uint64) (*hexutil.Big, error) {
    lockValue, err := api.service.GetDeveloperLockDIP(address, blockNumber)
    if err != nil {
        return nil, err
    }

    return (*hexutil.Big)(lockValue), nil
}

func (api *DipperinMercuryApi) GetFoundationInfo(usage economy_model.FoundationDIPUsage) map[string]*hexutil.Big {
    foundationInfo := api.service.GetFoundationInfo(usage)

    result := make(map[string]*hexutil.Big, 0)
    for key, value := range foundationInfo {
        result[key.Hex()] = (*hexutil.Big)(value)
    }
    return result
}

func (api *DipperinMercuryApi) GetMaintenanceLockDIP(address common.Address, blockNumber uint64) (*hexutil.Big, error) {
    lockValue, err := api.service.GetMaintenanceLockDIP(address, blockNumber)
    if err != nil {
        return nil, err
    }

    return (*hexutil.Big)(lockValue), nil
}

func (api *DipperinMercuryApi) GetReMainRewardLockDIP(address common.Address, blockNumber uint64) (*hexutil.Big, error) {
    lockValue, err := api.service.GetReMainRewardLockDIP(address, blockNumber)
    if err != nil {
        return nil, err
    }

    return (*hexutil.Big)(lockValue), nil
}

func (api *DipperinMercuryApi) GetEarlyTokenLockDIP(address common.Address, blockNumber uint64) (*hexutil.Big, error) {
    lockValue, err := api.service.GetEarlyTokenLockDIP(address, blockNumber)
    if err != nil {
        return nil, err
    }

    return (*hexutil.Big)(lockValue), nil
}

func (api *DipperinMercuryApi) GetMineMasterEDIPReward(blockNumber uint64, tokenDecimals int) (*hexutil.Big, error) {
    reward, err := api.service.GetMineMasterEDIPReward(blockNumber, tokenDecimals)
    if err != nil {
        return nil, nil
    }

    return (*hexutil.Big)(reward), nil
}

func (api *DipperinMercuryApi) GetVerifierEDIPReward(blockNumber uint64, tokenDecimals int) (map[economy_model.VerifierType]*hexutil.Big, error) {
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

func (api *DipperinMercuryApi) NewBlock(ctx context.Context) (*rpc.Subscription, error) {
    return api.service.NewBlock(ctx)
}

func (api *DipperinMercuryApi) SubscribeBlock(ctx context.Context) (*rpc.Subscription, error) {
    return api.service.SubscribeBlock(ctx)
}

func (api *DipperinMercuryApi) StopDipperin() {
    api.service.StopDipperin()
}
