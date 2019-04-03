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

package contract

import (
	"errors"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/consts"
	"github.com/dipperin/dipperin-core/common/hexutil"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/core/economy-model"
	"github.com/dipperin/dipperin-core/third-party/log"
	"math/big"
	"sync"
)

const (
	DecimalUnits = 3
	tokenName = "EarlyToken"
	tokenSymbol = "EarlyReward"
)
var EarlyContractAddress = common.HexToAddress("0x00110000000000000000000000000000000000000000")

var ProhibitFunction = []string{"create", "RewardMineMaster", "RewardVerifier"}

type EarlyRewardContract struct {
	BuiltInERC20Token
	Early economy_model.Foundation `json:"-"`

	//remaining token equl to DIP
	NeedDIP *big.Int `json:"need_coin"`
	//token cashed for DIP
	ChangeToDIPToken *big.Int `json:"change_to_coin_token"`

	ExchangeRate []int64 `json:"exchange_rate"`

	Lock sync.Mutex
}

type EarlyRewardContractForMarshaling struct {
	Erc20 BuiltInERC20Token `json:"erc_20"`

	//remaining token equl to DIP
	NeedDIP *hexutil.Big `json:"need_DIP"`
	//token exchanged for DIP
	ChangeToDIPToken *hexutil.Big `json:"change_to_DIP_token"`
	//exchange ratio
	ExchangeRate []int64 `json:"exchange_rate"`
}

var EarlyRewardContractStr string

func init(){
	foundation := economy_model.MakeDipperinFoundation(economy_model.DIPProportion)
	owner := economy_model.EarlyTokenAddresses[0]
	decimalBase := big.NewInt(0).Exp(big.NewInt(10), big.NewInt(int64(DecimalUnits)), nil)
	initAmount := big.NewInt(0).Mul(economy_model.EarlyTokenAmount, decimalBase)

	contract,err := MakeEarlyRewardContract(foundation, initAmount, economy_model.InitExchangeRate, tokenName, DecimalUnits, tokenSymbol, owner)
	if err !=nil{
		panic("early_token init panic")
	}

	EarlyRewardContractStr = util.StringifyJson(contract)
}

//DIP = eDIP*const.DIP*exChangeRate/EarlyTokenExchangeBase/decimalBase
func calcNeedDIP(eDIP *big.Int, decimalUnits int, exChangeRate int64) *big.Int {
	decimalBase := big.NewInt(0).Exp(big.NewInt(10), big.NewInt(int64(decimalUnits)), nil)
	actualNeedDIP := big.NewInt(0)
	actualNeedDIP.Mul(eDIP, big.NewInt(exChangeRate))
	actualNeedDIP.Mul(actualNeedDIP, big.NewInt(consts.DIP))
	actualNeedDIP.Div(actualNeedDIP, big.NewInt(economy_model.EarlyTokenExchangeBase))
	actualNeedDIP.Div(actualNeedDIP, decimalBase)

	return actualNeedDIP
}

func MakeEarlyRewardContract(foundation economy_model.Foundation, initAmount *big.Int, initExchangeRate int64, tokenName string, decimalUnits int, tokenSymbol string, owner common.Address) (*EarlyRewardContract, error) {
	actualNeedDIP := calcNeedDIP(initAmount, decimalUnits, initExchangeRate)
	if actualNeedDIP.Cmp(economy_model.EarlyTokenDIP) == 1 {
		return nil, errors.New("the DIP isn't enough")
	}
	return &EarlyRewardContract{
		BuiltInERC20Token: *newToken(initAmount, tokenName, decimalUnits, tokenSymbol, owner),
		Early:             foundation,
		NeedDIP:           actualNeedDIP,
		ChangeToDIPToken:  big.NewInt(0),
		ExchangeRate:      []int64{initExchangeRate},
	}, nil
}

func (earlyToken EarlyRewardContract) MarshalJSON() ([]byte, error) {
	marshalData := &EarlyRewardContractForMarshaling{
		Erc20:            earlyToken.BuiltInERC20Token,
		NeedDIP:          (*hexutil.Big)(earlyToken.NeedDIP),
		ChangeToDIPToken: (*hexutil.Big)(earlyToken.ChangeToDIPToken),
		ExchangeRate:     earlyToken.ExchangeRate,
	}
	return util.StringifyJsonToBytesWithErr(marshalData)
}

func (earlyToken *EarlyRewardContract) UnmarshalJSON(input []byte) error {

	var unmarshalData EarlyRewardContractForMarshaling
	if err := util.ParseJsonFromBytes(input, &unmarshalData); err != nil {
		return err
	}

	earlyToken.BuiltInERC20Token = unmarshalData.Erc20
	earlyToken.NeedDIP = unmarshalData.NeedDIP.ToInt()
	earlyToken.ChangeToDIPToken = unmarshalData.ChangeToDIPToken.ToInt()
	earlyToken.ExchangeRate = unmarshalData.ExchangeRate

	return nil
}

func (earlyToken *EarlyRewardContract) IsValid() error {
	return nil
}

func (earlyToken *EarlyRewardContract) getDipperinFoundation() economy_model.Foundation {
	if earlyToken.Early != nil {
		return earlyToken.Early
	}

	earlyToken.Early = economy_model.MakeDipperinFoundation(economy_model.DIPProportion)
	return earlyToken.Early
}

//get current exchangeRate
func (earlyToken *EarlyRewardContract) GetExchangeRate() (exchangeRate int64) {
	return earlyToken.ExchangeRate[len(earlyToken.ExchangeRate)-1]
}

//set exchange rate，only maintenance autherized
func (earlyToken *EarlyRewardContract) SetExchangeRate(from common.Address, exchangeRate int64) error {
	earlyToken.Lock.Lock()
	defer earlyToken.Lock.Unlock()
	//validate address
	addressType := earlyToken.getDipperinFoundation().GetAddressType(from)
	if addressType != economy_model.MaintenanceAddress {
		return errors.New("the address isn't foundation maintenance address")
	}
	if exchangeRate == earlyToken.ExchangeRate[len(earlyToken.ExchangeRate)-1] {
		return nil
	}

	//calculate DIP needed
	decimal := earlyToken.Decimals()
	decimalBase := big.NewInt(0).Exp(big.NewInt(10), big.NewInt(int64(decimal)), nil)
	tokenAmount := big.NewInt(0).Mul(economy_model.EarlyTokenAmount, decimalBase)
	notExchangeEDIP := big.NewInt(0).Sub(tokenAmount, earlyToken.ChangeToDIPToken)

	needDIP := calcNeedDIP(notExchangeEDIP, decimal, exchangeRate)
	//refund
	cmpResult := needDIP.Cmp(earlyToken.NeedDIP)

	if cmpResult == -1 {
		//after exchange rate modified，the real DIP less than Token DIP, refund creator DIP
		earlyToken.AccountDB.AddBalance(earlyToken.Owner, big.NewInt(0).Sub(earlyToken.NeedDIP, needDIP))
		earlyToken.NeedDIP = needDIP
		earlyToken.ExchangeRate = append(earlyToken.ExchangeRate, exchangeRate)
	} else if cmpResult == 1 {
		//after exchange rate modified，the real DIP more than Token DIP, minus DIP
		earlyToken.AccountDB.SubBalance(from, big.NewInt(0).Sub(needDIP, earlyToken.NeedDIP))
		earlyToken.NeedDIP = needDIP
		earlyToken.ExchangeRate = append(earlyToken.ExchangeRate, exchangeRate)
	}
	return nil
}

func (earlyToken *EarlyRewardContract) TransferEDIPToDIP(from common.Address, eDIPValue *hexutil.Big) error {
	earlyToken.Lock.Lock()
	defer earlyToken.Lock.Unlock()

	//check whether address is normal
	addressType := earlyToken.getDipperinFoundation().GetAddressType(from)
	if addressType != economy_model.NotFoundationAddress {
		log.Info("the addressType is:", "addressType", addressType)
		return errors.New("the address isn't NotFoundationAddress")
	}

	if earlyToken.Balances[from.Hex()].Cmp(eDIPValue.ToInt()) == -1 {
		return errors.New("the token isn't enough")
	}

	log.Info("the eDIP value is:", "eDIP", eDIPValue)

	//calculate DIP needed
	decimal := earlyToken.Decimals()
	decimalBase := big.NewInt(0).Exp(big.NewInt(10), big.NewInt(int64(decimal)), nil)
	DIP := big.NewInt(0)
	currentExchangeRate := earlyToken.ExchangeRate[len(earlyToken.ExchangeRate)-1]

	log.Info("the currentExchangeRate is:", "currentExchangeRate", currentExchangeRate)

	DIP.Mul(eDIPValue.ToInt(), big.NewInt(currentExchangeRate))
	DIP.Mul(DIP, big.NewInt(consts.DIP))
	DIP.Div(DIP, big.NewInt(economy_model.EarlyTokenExchangeBase))
	DIP.Div(DIP, decimalBase)
	earlyToken.NeedDIP.Sub(earlyToken.NeedDIP, DIP)

	earlyToken.Balances[from.Hex()].Sub(earlyToken.Balances[from.Hex()], eDIPValue.ToInt())
	earlyToken.AccountDB.AddBalance(from, DIP)

	return nil
}

func (earlyToken *EarlyRewardContract) Destroy(from common.Address) error {
	return errors.New("can't destroy the contract")
}

func (earlyToken *EarlyRewardContract) Transfer(toAddress common.Address, hValue *hexutil.Big) error {
	addressType := earlyToken.getDipperinFoundation().GetAddressType(earlyToken.CurSender)
	if addressType != economy_model.NotFoundationAddress {
		return errors.New("the address should be normalAddress")
	} else {
		return earlyToken.BuiltInERC20Token.Transfer(toAddress, hValue)
	}
}

func (earlyToken *EarlyRewardContract) TransferFrom(fromAddress, toAddress common.Address, value *hexutil.Big) bool {
	addressType := earlyToken.getDipperinFoundation().GetAddressType(fromAddress)
	if addressType != economy_model.NotFoundationAddress {
		return false
	}

	addressType = earlyToken.getDipperinFoundation().GetAddressType(earlyToken.CurSender)
	if addressType != economy_model.NotFoundationAddress {
		return false
	}

	return earlyToken.BuiltInERC20Token.TransferFrom(fromAddress, toAddress, value)
}

func (earlyToken *EarlyRewardContract) Approve(spenderAddress common.Address, value *hexutil.Big) bool {
	addressType := earlyToken.getDipperinFoundation().GetAddressType(earlyToken.CurSender)
	if addressType != economy_model.NotFoundationAddress {
		return false
	} else {
		return earlyToken.BuiltInERC20Token.Approve(spenderAddress, value)
	}
}

func (earlyToken *EarlyRewardContract) Allowance(ownerAddress, spenderAddress common.Address) *big.Int {
	addressType := earlyToken.getDipperinFoundation().GetAddressType(ownerAddress)
	if addressType != economy_model.NotFoundationAddress {
		return big.NewInt(0)
	}

	addressType = earlyToken.getDipperinFoundation().GetAddressType(spenderAddress)
	if addressType != economy_model.NotFoundationAddress {
		return big.NewInt(0)
	}

	addressType = earlyToken.getDipperinFoundation().GetAddressType(earlyToken.CurSender)
	if addressType != economy_model.NotFoundationAddress {
		return big.NewInt(0)
	}

	return earlyToken.BuiltInERC20Token.Allowance(ownerAddress, spenderAddress)
}

//give mineMaster extra bonus every block
func (earlyToken *EarlyRewardContract) RewardMineMaster(DIPReward *big.Int, blockNumber uint64, rewardAddress common.Address) error {
	rewardEDIP, err := earlyToken.getDipperinFoundation().GetMineMasterEDIPReward(DIPReward, blockNumber, earlyToken.Decimals())
	if err != nil {
		return err
	}

	log.Info("the token owner value is:", "value", earlyToken.Balances[earlyToken.Owner.Hex()])
	log.Info("the rewardEDIP value is:", "rewardEDIP", rewardEDIP)
	if rewardEDIP.Cmp(big.NewInt(0)) == 0 {
		return nil
	}

	if earlyToken.Balances[earlyToken.Owner.Hex()].Cmp(rewardEDIP) == -1 {
		return errors.New("RewardMineMaster the Early token contract token isn't enough")
	}

	if _, ok := earlyToken.Balances[rewardAddress.Hex()]; !ok {
		earlyToken.Balances[rewardAddress.Hex()] = big.NewInt(0)
	}

	earlyToken.Balances[rewardAddress.Hex()].Add(earlyToken.Balances[rewardAddress.Hex()], rewardEDIP)
	earlyToken.Balances[earlyToken.Owner.Hex()].Sub(earlyToken.Balances[earlyToken.Owner.Hex()], rewardEDIP)

	return nil
}

//giev Verifier extra bonus every block
func (earlyToken *EarlyRewardContract) RewardVerifier(DIPReward map[economy_model.VerifierType]*big.Int, blockNumber uint64, verifierAddress map[economy_model.VerifierType][]common.Address) error {
	rewardEDIP, err := earlyToken.getDipperinFoundation().GetVerifierEDIPReward(DIPReward, blockNumber, earlyToken.Decimals())
	if err != nil {
		return err
	}

	for verifierType, rewardVale := range rewardEDIP {
		if rewardVale.Cmp(big.NewInt(0)) == 0 {
			return nil
		}

		if earlyToken.Balances[earlyToken.Owner.Hex()].Cmp(rewardVale) == -1 {
			return errors.New("RewardVerifier the Early token contract token isn't enough")
		}

		if verifierType == economy_model.MasterVerifier {
			masterAddress := verifierAddress[economy_model.MasterVerifier][0].Hex()
			if _, ok := earlyToken.Balances[masterAddress]; !ok {
				earlyToken.Balances[masterAddress] = big.NewInt(0)
			}
			earlyToken.Balances[masterAddress].Add(earlyToken.Balances[masterAddress], rewardVale)
			earlyToken.Balances[earlyToken.Owner.Hex()].Sub(earlyToken.Balances[earlyToken.Owner.Hex()], rewardVale)
		} else {
			for _, address := range verifierAddress[verifierType] {
				if _, ok := earlyToken.Balances[address.Hex()]; !ok {
					earlyToken.Balances[address.Hex()] = big.NewInt(0)
				}
				earlyToken.Balances[address.Hex()] = big.NewInt(0).Add(earlyToken.Balances[address.Hex()], rewardVale)
				earlyToken.Balances[earlyToken.Owner.Hex()] = big.NewInt(0).Sub(earlyToken.Balances[earlyToken.Owner.Hex()], rewardVale)
			}
		}
	}

	return nil
}
