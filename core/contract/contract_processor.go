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
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/core/model"
	"github.com/dipperin/dipperin-core/third-party/log"
	"reflect"
)

var (
	CanNotParseContractErr      = errors.New("cannot parse transaction extra data")
	ContractAdrEmptyErr         = errors.New("contract address can't be empty")
	ContractWithoutValidatorErr = errors.New("no found validator method，cannot create contract")
	ContractValidatorRetNilErr  = errors.New("contract validator method return nothing")
	ContractMethodRetNilErr     = errors.New("contract method return nothing")
	ContractMethodFailErr       = errors.New("contract method return false")
)

func NewProcessor(cDB ContractDB, blockHeight uint64) *Processor {
	return &Processor{contractDB: cDB, blockHeight: blockHeight}
}

type Processor struct {
	contractDB  ContractDB
	accountDB   AccountDB
	blockHeight uint64
}

func (p *Processor) SetAccountDB(db AccountDB) {
	p.accountDB = db
}

// when running the operation which can modify contract status，changeState decides whether change contract status
func (p *Processor) Process(tx model.AbstractTransaction) (err error) {
	eData := ParseExtraDataForContract(tx.ExtraData())
	if eData == nil {
		return CanNotParseContractErr
	}
	// must be to
	eData.ContractAddress = *tx.To()
	log.Debug("Processor Process", "eData.Action", eData.Action, "eData.ContractAddress", eData.ContractAddress)

	var result reflect.Value
	switch eData.Action {
	case "create":
		result, err = p.DoCreate(eData)
	default:
		sender, sErr := tx.Sender(nil)
		if sErr != nil {
			return sErr
		}
		result, err = p.Run(sender, eData)
	}
	// modify contract
	if err == nil {
		// TODO: check the type of contract address
		err = p.contractDB.PutContract(eData.ContractAddress, result)
	}

	return
}

// create contract
func (p *Processor) DoCreate(eData *ExtraDataForContract) (reflect.Value, error) {
	if eData.ContractAddress.IsEmpty() {
		return reflect.Value{}, ContractAdrEmptyErr
	}
	if p.contractDB.ContractExist(eData.ContractAddress) {
		return reflect.Value{}, errors.New(fmt.Sprintf("can't create contract, address already have contract data: %v", eData.ContractAddress))
	}

	contractType := eData.ContractAddress.GetAddressTypeStr()
	log.Debug("Processor doCreate")
	ct, ctErr := GetContractTempByType(contractType)
	if ctErr != nil {
		return reflect.Value{}, ctErr
	}
	nContract := reflect.New(ct)
	if err := util.ParseJson(eData.Params, nContract.Interface()); err != nil {
		return reflect.Value{}, err
	}

	// validate contract
	vMethod := nContract.MethodByName("IsValid")
	if vMethod.Kind() != reflect.Func {
		return reflect.Value{}, ContractWithoutValidatorErr
	}
	vResult := vMethod.Call([]reflect.Value{})
	if len(vResult) == 0 {
		return reflect.Value{}, ContractValidatorRetNilErr
	}
	if !vResult[0].IsNil() {
		return reflect.Value{}, vResult[0].Interface().(error)
	}

	// block height must be saved in state db, or meet hash collision
	nContract.Elem().FieldByName("CurBlockHeight").Set(reflect.ValueOf(p.blockHeight))
	//record balances
	owner := nContract.Elem().FieldByName("Owner").Interface().(common.Address).Hex()
	// todo if not ERC20 and EARLY TOKEN？
	amount := nContract.Elem().FieldByName("TokenTotalSupply")
	nContract.Elem().FieldByName("Balances").SetMapIndex(reflect.ValueOf(owner), amount)
	return nContract, nil
}

// run contract
func (p *Processor) Run(executorAddress common.Address, eData *ExtraDataForContract) (reflect.Value, error) {
	// get contract type
	contractType := eData.ContractAddress.GetAddressTypeStr()
	log.Debug("run contract method", "contract type", contractType, "method", eData.Action)
	// get contract from type
	ct, ctErr := GetContractTempByType(contractType)
	if ctErr != nil {
		return reflect.Value{}, ctErr
	}

	// get contract
	nContract, err := p.contractDB.GetContract(eData.ContractAddress, ct)
	if err != nil {
		return reflect.Value{}, err
	}

	// set caller address
	tmpF := nContract.Elem().FieldByName("CurSender")
	if tmpF.CanSet() {
		tmpF.Set(reflect.ValueOf(executorAddress))
	}
	// block height must be saved in state db, or meet hash collision
	nContract.Elem().FieldByName("CurBlockHeight").Set(reflect.ValueOf(p.blockHeight))

	tmpF = nContract.Elem().FieldByName("AccountDB")
	aDBV := reflect.ValueOf(p.accountDB)
	if tmpF.CanSet() && aDBV.IsValid() {
		tmpF.Set(reflect.ValueOf(p.accountDB))
	}

	method := nContract.MethodByName(eData.Action)
	if method.Kind() != reflect.Func {
		return reflect.Value{}, errors.New("not found method:" + eData.Action)
	}

	// parse user's input
	mArgs, ctmErr := GetContractMethodArgs(contractType, eData.Action)
	if ctmErr != nil {
		return reflect.Value{}, ctmErr
	}
	codec := NewParamsCodec(bufio.NewReader(bytes.NewBufferString(eData.Params)))
	rValue, pErr := codec.ParseRequestArguments(mArgs)
	if pErr != nil {
		log.Info("parse parameter error", "params", eData.Params, "err", pErr)
		for _, ma := range mArgs {
			log.Info("arg type", "arg t", ma.String())
		}
		return reflect.Value{}, pErr
	}
	result := method.Call(rValue)

	// check result
	if len(result) == 0 {
		log.Warn("contract method return nothing", "contract type", contractType, "contract address", eData.ContractAddress.Hex())
		return reflect.Value{}, ContractMethodRetNilErr
	}

	// method return bool
	if result[0].Kind() == reflect.Bool {
		if result[0].Bool() {
			return nContract, nil
		} else {
			return reflect.Value{}, ContractMethodFailErr
		}
	}
	// method return error
	if errR, ok := result[0].Interface().(error); ok {
		if !result[0].IsNil() {
			//return nil, result[0].Interface().(error)
			return reflect.Value{}, errR
		}
	}
	// empty means no err
	return nContract, nil
}

// get contract readonly infomation（not modify contract）
func (p *Processor) GetContractReadOnlyInfo(eData *ExtraDataForContract) (interface{}, error) {
	log.Info("GetContractReadOnlyInfo", "addr", eData.ContractAddress, "action", eData.Action)
	contractType := eData.ContractAddress.GetAddressTypeStr()
	// get contract by type
	ct, ctErr := GetContractTempByType(contractType)
	if ctErr != nil {
		return nil, ctErr
	}

	nContract, err := p.contractDB.GetContract(eData.ContractAddress, ct)
	if err != nil {
		return nil, err
	}

	method := nContract.MethodByName(eData.Action)
	if method.Kind() != reflect.Func {
		return nil, errors.New("not found method:" + eData.Action)
	}

	// convert user's input
	mArgs, ctmErr := GetContractMethodArgs(contractType, eData.Action)
	if ctmErr != nil {
		return nil, ctmErr
	}

	codec := NewParamsCodec(bufio.NewReader(bytes.NewBufferString(eData.Params)))
	rValue, pErr := codec.ParseRequestArguments(mArgs)
	if pErr != nil {
		log.Debug("parse parameter error", "err", pErr, "params", eData.Params, "m args", mArgs)
		return nil, pErr
	}
	result := method.Call(rValue)

	// check result
	if len(result) == 0 {
		log.Warn("contract method return nothing", "contract type", contractType, "contract address", eData.ContractAddress.Hex())
		return nil, ContractMethodRetNilErr
	}
	return result[0].Interface(), nil
}
