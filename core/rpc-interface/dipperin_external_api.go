package rpc_interface

import (
	"context"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/hexutil"
	"github.com/dipperin/dipperin-core/core/contract"
	"github.com/dipperin/dipperin-core/core/economy-model"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/core/vm/common/utils"
	model2 "github.com/dipperin/dipperin-core/core/vm/model"
	"github.com/dipperin/dipperin-core/third-party/rpc"
	"math/big"
)

type DipperExternalApi struct {
	allApis *DipperinVenusApi
}

// verify whether the chain is in sync
func (api *DipperExternalApi) GetSyncStatus() bool {
	return api.allApis.GetSyncStatus()
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
func (api *DipperExternalApi) CurrentBlock() (*BlockResp, error) {
	return api.allApis.CurrentBlock()
}


func (api *DipperExternalApi)CurrentBalance(address common.Address) (resp *CurBalanceResp, err error){
	return api.allApis.CurrentBalance(address)
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
func (api *DipperExternalApi) GetBlockByNumber(number uint64) (*BlockResp, error) {
	return api.allApis.GetBlockByNumber(number)
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
func (api *DipperExternalApi) GetBlockByHash(hash common.Hash) (*BlockResp, error) {
	return api.allApis.GetBlockByHash(hash)
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
func (api *DipperExternalApi) GetBlockNumber(hash common.Hash) *uint64 {
	return api.allApis.GetBlockNumber(hash)
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
func (api *DipperExternalApi) GetGenesis() (*BlockResp, error) {
	return api.allApis.GetGenesis()
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
func (api *DipperExternalApi) GetBlockBody(hash common.Hash) *model.Body {
	return api.allApis.GetBlockBody(hash)
}

func (api *DipperExternalApi) Transaction(hash common.Hash) (resp *TransactionResp, err error) {
	return api.allApis.Transaction(hash)
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
func (api *DipperExternalApi) GetTransactionNonce(addr common.Address) (nonce uint64, err error) {
	return api.allApis.GetTransactionNonce(addr)
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
func (api *DipperExternalApi) NewTransaction(transactionRlpB []byte) (TxHash common.Hash, err error) {
	return api.allApis.NewTransaction(transactionRlpB)
}

// call contract
func (api *DipperExternalApi) NewContract(transactionRlpB []byte, blockNum uint64) (resp string, err error) {
	return api.allApis.NewContract(transactionRlpB, blockNum)
}

func (api *DipperExternalApi) NewEstimateGas(transactionRlpB []byte) (resp hexutil.Uint64, err error) {
	return api.allApis.NewEstimateGas(transactionRlpB)
}

func (api *DipperExternalApi) GetContractInfo(eData *contract.ExtraDataForContract) (interface{}, error) {
	return api.allApis.GetContractInfo(eData)
}

func (api *DipperExternalApi) GetContract(contractAddr common.Address) (interface{}, error) {
	return api.allApis.GetContract(contractAddr)
}

func (api *DipperExternalApi) GetVerifiersBySlot(slotNum uint64) ([]common.Address, error) {
	return api.allApis.GetVerifiersBySlot(slotNum)
}

func (api *DipperExternalApi) GetSlotByNumber(blockNum uint64) (uint64, error) {
	return api.allApis.GetSlotByNumber(blockNum)
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
func (api *DipperExternalApi) GetCurVerifiers() []common.Address {
	return api.allApis.GetCurVerifiers()
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
func (api *DipperExternalApi) GetNextVerifiers() []common.Address {
	return api.allApis.GetNextVerifiers()
}

func (api *DipperExternalApi) VerifierStatus(address common.Address) (resp *VerifierStatus, err error) {
	return api.allApis.VerifierStatus(address)
}

func (api *DipperExternalApi) CurrentStake(address common.Address) (resp *CurStakeResp, err error) {
	return api.allApis.CurrentStake(address)
}

func (api *DipperExternalApi) CurrentReputation(address common.Address) (uint64, error) {
	return api.allApis.CurrentReputation(address)
}

func (api *DipperExternalApi) GetBlockDiffVerifierInfo(blockNumber uint64) (map[economy_model.VerifierType][]common.Address, error) {
	return api.allApis.GetBlockDiffVerifierInfo(blockNumber)
}

func (api *DipperExternalApi) GetVerifierDIPReward(blockNumber uint64) (map[economy_model.VerifierType]*hexutil.Big, error) {
	return api.allApis.GetVerifierDIPReward(blockNumber)
}

func (api *DipperExternalApi) GetMineMasterDIPReward(blockNumber uint64) (*hexutil.Big, error) {
	return api.allApis.GetMineMasterDIPReward(blockNumber)
}

func (api *DipperExternalApi) GetBlockYear(blockNumber uint64) (uint64, error) {
	return api.allApis.GetBlockYear(blockNumber)
}

func (api *DipperExternalApi) GetOneBlockTotalDIPReward(blockNumber uint64) (*hexutil.Big, error) {
	return api.allApis.GetOneBlockTotalDIPReward(blockNumber)
}

func (api *DipperExternalApi) GetInvestorInfo() map[string]*hexutil.Big {
	return api.allApis.GetInvestorInfo()
}

func (api *DipperExternalApi) GetDeveloperInfo() map[string]*hexutil.Big {
	return api.allApis.GetDeveloperInfo()
}

func (api *DipperExternalApi) GetAddressLockMoney(address common.Address) (*big.Int, error) {
	return api.allApis.GetAddressLockMoney(address)
}

func (api *DipperExternalApi) GetInvestorLockDIP(address common.Address, blockNumber uint64) (*hexutil.Big, error) {
	return api.allApis.GetInvestorLockDIP(address, blockNumber)
}

func (api *DipperExternalApi) GetDeveloperLockDIP(address common.Address, blockNumber uint64) (*hexutil.Big, error) {
	return api.allApis.GetDeveloperLockDIP(address, blockNumber)
}

func (api *DipperExternalApi) GetFoundationInfo(usage economy_model.FoundationDIPUsage) map[string]*hexutil.Big {
	return api.allApis.GetFoundationInfo(usage)
}

func (api *DipperExternalApi) GetMaintenanceLockDIP(address common.Address, blockNumber uint64) (*hexutil.Big, error) {
	return api.allApis.GetMaintenanceLockDIP(address, blockNumber)
}

func (api *DipperExternalApi) GetReMainRewardLockDIP(address common.Address, blockNumber uint64) (*hexutil.Big, error) {
	return api.allApis.GetReMainRewardLockDIP(address, blockNumber)
}

func (api *DipperExternalApi) GetEarlyTokenLockDIP(address common.Address, blockNumber uint64) (*hexutil.Big, error) {
	return api.allApis.GetEarlyTokenLockDIP(address, blockNumber)
}

func (api *DipperExternalApi) GetMineMasterEDIPReward(blockNumber uint64, tokenDecimals int) (*hexutil.Big, error) {
	return api.allApis.GetMineMasterEDIPReward(blockNumber, tokenDecimals)
}

func (api *DipperExternalApi) GetVerifierEDIPReward(blockNumber uint64, tokenDecimals int) (map[economy_model.VerifierType]*hexutil.Big, error) {
	return api.allApis.GetVerifierEDIPReward(blockNumber, tokenDecimals)
}

func (api *DipperExternalApi) GetABI(contractAddr common.Address) (*utils.WasmAbi, error) {
	return api.allApis.GetABI(contractAddr)
}

func (api *DipperExternalApi) GetCode(contractAddr common.Address) ([]byte, error) {
	return api.allApis.GetCode(contractAddr)
}

func (api *DipperExternalApi) SuggestGasPrice() (resp *CurBalanceResp, err error) {
	return api.allApis.SuggestGasPrice()
}

func (api *DipperExternalApi) GetContractAddressByTxHash(txHash common.Hash) (common.Address, error) {
	return api.allApis.GetContractAddressByTxHash(txHash)
}

func (api *DipperExternalApi) GetLogs(blockHash common.Hash, fromBlock, toBlock uint64, Addresses []common.Address, Topics [][]common.Hash) ([]*model2.Log, error) {
	return api.allApis.GetLogs(blockHash, fromBlock, toBlock, Addresses, Topics)
}

func (api *DipperExternalApi) GetTxActualFee(txHash common.Hash) (resp *CurBalanceResp, err error) {
	return api.allApis.GetTxActualFee(txHash)
}

func (api *DipperExternalApi) GetReceiptByTxHash(txHash common.Hash) (*model2.Receipt, error) {
	return api.allApis.GetReceiptByTxHash(txHash)
}

func (api *DipperExternalApi) GetReceiptsByBlockNum(num uint64) (model2.Receipts, error) {
	return api.allApis.GetReceiptsByBlockNum(num)
}

func (api *DipperExternalApi) NewBlock(ctx context.Context) (*rpc.Subscription, error) {
	return api.allApis.NewBlock(ctx)
}

func (api *DipperExternalApi) SubscribeBlock(ctx context.Context) (*rpc.Subscription, error) {
	return api.allApis.SubscribeBlock(ctx)
}
