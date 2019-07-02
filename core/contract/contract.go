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
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/third-party/log"
	"io/ioutil"
	"path/filepath"
	"reflect"
)

var contracts = map[string]*InfoOfContract{
	consts.ERC20TypeName:      newInfoOfContract(BuiltInERC20Token{}),
	consts.EarlyTokenTypeName: newInfoOfContract(EarlyRewardContract{}),
}

// contract infomation
type InfoOfContract struct {
	TypeOfContract reflect.Type
	MethodArgs     map[string][]reflect.Type
}

// analyse contract information the parameter cannot be point
func newInfoOfContract(contract interface{}) *InfoOfContract {
	contractT := reflect.TypeOf(contract)
	// methods are New object's methods
	cvt := reflect.New(contractT).Type()
	info := &InfoOfContract{TypeOfContract: contractT, MethodArgs: map[string][]reflect.Type{}}
	nm := cvt.NumMethod()
	//log.Debug("initialize contract", "method num", nm, "contract", info.TypeOfContract.String())
	for i := 0; i < nm; i++ {
		tmpM := cvt.Method(i)
		info.MethodArgs[tmpM.Name] = []reflect.Type{}
		//cslog.Debug().Str("method name", tmpM.Name).Msg("init method")
		nIn := tmpM.Type.NumIn()
		for j := 0; j < nIn; j++ {
			if j != 0 {
				//cslog.Debug().Str("parameter", tmpM.Type.In(j).String()).Msg("parameter")
				info.MethodArgs[tmpM.Name] = append(info.MethodArgs[tmpM.Name], tmpM.Type.In(j))
			}
		}
	}
	return info
}

// extra data in transaction extra data
type ExtraDataForContract struct {
	// contract address
	ContractAddress common.Address `json:"contract_address"`
	//ContractType string `json:"contract_type"`
	// operation，create/contract function
	Action string `json:"action"`
	// parameters according to different methods, ["a", 123]
	Params string `json:"params"`
}

func SaveContractId(path, node string, id common.Address) {

	fPath := ""
	if path == "" {
		fPath = filepath.Join(util.HomeDir(), "tmp/dipperin_apps/", node, "contract")
	} else {
		fPath = filepath.Join(path, node, "contract")
	}

	//fd,_:=os.OpenFile(fPath,os.O_RDWR|os.O_CREATE|os.O_APPEND,0644)
	var ctId []string
	cont, e := ioutil.ReadFile(fPath)
	if e == nil {
		//file exist, ignore error while json parsing
		util.ParseJsonFromBytes(cont, &ctId)
	}
	ctId = append(ctId, id.Hex())
	content := util.StringifyJsonToBytes(ctId)
	ioutil.WriteFile(fPath, content, 0644)
}

func GetContractId(path, node string) (cAdr []common.Address, err error) {
	fPath := ""
	if path == "" {
		fPath = filepath.Join(util.HomeDir(), "tmp/dipperin_apps/", node, "contract")
	} else {
		fPath = filepath.Join(path, node, "contract")
	}
	var ctId []string
	cont, e := ioutil.ReadFile(fPath)
	if e == nil {
		//file exist, ignore error while json parsing
		util.ParseJsonFromBytes(cont, &ctId)
	}

	if len(ctId) == 0 {
		err = errors.New("non-exsisted contract record")
	}

	for i := 0; i < len(ctId); i++ {
		cAdr = append(cAdr, common.HexToAddress(ctId[i]))
	}
	return
}

// from transaction extra data, get contract infomation
func ParseExtraDataForContract(data []byte) *ExtraDataForContract {
	var eData ExtraDataForContract
	if err := util.ParseJsonFromBytes(data, &eData); err != nil {
		log.Info("contract information of extra data extracting error", "err", err)
		return nil
	} else {
		return &eData
	}
}

// lookup contract type
func GetContractTempByType(cType string) (reflect.Type, error) {
	if contracts[cType] == nil {
		return nil, errors.New("not found:" + cType)
	}
	return contracts[cType].TypeOfContract, nil
}

// get parameters from method
func GetContractMethodArgs(cType string, mName string) ([]reflect.Type, error) {
	if contracts[cType] == nil {
		return nil, errors.New("not found:" + cType)
	}
	if contracts[cType].MethodArgs[mName] == nil {
		return nil, errors.New("not found method:" + mName)
	}
	return contracts[cType].MethodArgs[mName], nil
}

// convert infomation to contract object
func ParseContractFromBytes(cTypeStr string, cb []byte) (interface{}, error) {
	ct, ctErr := GetContractTempByType(cTypeStr)
	if ctErr != nil {
		return nil, ctErr
	}
	// create contract
	nContract := reflect.New(ct)
	if err := util.ParseJsonFromBytes(cb, nContract.Interface()); err != nil {
		log.Debug("parse contract error", "err", err)
		return nil, err
	}
	return nContract.Interface(), nil
}

//type ERC20ExtraData struct {
//	// contract address
//	ContractAddress string `json:"contract_address"`
//
//	// contract method
//	ContractMethod string `json:"contract_method"`
//
//	CreateToken *erc20.CreateERC20Config `json:"create_token"`
//
//	//BalanceOf *erc20.BalanceOfParams `json:"balance_of"`
//
//	Transfer *erc20.TransferParams `json:"transfer"`
//
//	TransferFrom *erc20.TransferFromParams `json:"transfer_from"`
//
//	Approve *erc20.ApproveParams `json:"approve"`
//
//	//Allowance *erc20.AllowanceParams `json:"allowance"`
//}

//var contractInfoForClient atomic.Value
//// Deprecated get contract
//func GetContractConf() map[string]*ContractInfo {
//	if conf := contractInfoForClient.Load(); conf != nil {
//		return conf.(map[string]*ContractInfo)
//	}
//	conf := map[string]*ContractInfo{
//		consts.ERC20TypeName: erc20.GetContractConfig(),
//	}
//	contractInfoForClient.Store(conf)
//	return conf
//}
