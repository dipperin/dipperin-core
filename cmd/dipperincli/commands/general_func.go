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


package commands

import (
	"github.com/dipperin/dipperin-core/third-party/log"
	"path/filepath"
	"math/big"
	"fmt"
	"strconv"
	"github.com/dipperin/dipperin-core/common/consts"

	"github.com/dipperin/dipperin-core/common"
	"errors"
	"github.com/dipperin/dipperin-core/common/hexutil"
)

//check address format
func CheckAndChangeHexToAddress(address string) (common.Address, error) {
	// Ignore 0x
	if len(address) - 2 != common.AddressLength * 2 {
		log.Error("the address is:","len",len(address),"addr",address)
		return common.Address{}, errors.New("address length is invalid")
	}

	if address[:2] != "0x" && address[:2] != "0X" {
		return common.Address{}, errors.New("address prefix should be 0x or 0X")
	}

	commonAddress := common.HexToAddress(address)

	addressType := commonAddress.GetAddressTypeStr()
	if addressType == "UnKnown" {
		return common.Address{}, errors.New("the address type error")
	}

	return commonAddress, nil
}

func ParseWalletPathAndName(inputPath string) (path, name string) {
	return inputPath, filepath.Base(inputPath)
}

func DecimalToInter(src string, unitBit int) (*big.Int, error) {
	length := len(src)
	if (length == 0) {
		return nil, errors.New("missing number")
	}
	scalingPos := 0
	if (src[0] < '0' && src[0] > '9') || (src[length-1] < '0' && src[length-1] > '9') {
		return nil, errors.New("the first and last character should be 0~9")
	}

	unit := 1
	for i:=0;i<unitBit;i++{
		unit *= 10
	}

	for i := 1; i < length-1; i++ {
		if src[i] < '0' || src[i] > '9' {
			if src[i] == '.' {
				scalingPos = i
			} else {
				errInfo := fmt.Sprintf("the character that index is:%v is invalid", i)
				return nil, errors.New(errInfo)
			}
		}
	}

	var interString string
	var decimalString string
	if scalingPos == 0 {
		interString = src
		decimalString = ""
	} else {
		interString = src[:scalingPos]
		decimalString = src[scalingPos+1:]
	}

	integerValue, err := strconv.ParseInt(interString, 10, 64)
	if err != nil {
		return nil, err
	}

	decimalLen := len(decimalString)
	if unitBit < decimalLen {
		return nil, errors.New("decimal length is invalid")
	}

	padding := make([]byte, unitBit-decimalLen)
	for index := range padding {
		padding[index] = '0'
	}

	tmpValue := append([]byte(decimalString), padding[:]...)

	decimalValue, err := strconv.ParseInt(string(tmpValue), 10, 64)
	if err != nil {
		return nil, err
	}

	totalValue := integerValue*int64(unit) + decimalValue
	return big.NewInt(totalValue), nil
}

func InterToDecimal(csCoinValue *hexutil.Big, unitBit int) (string, error) {
	if csCoinValue == nil {
		return "", errors.New("csCoinValue is nil")
	}

	if csCoinValue.ToInt().Cmp(big.NewInt(0)) == 0 {
		return "0", nil
	}

	coinValue := csCoinValue.ToInt().String()
	coinValueLen := len(coinValue)

	//remove 0 of the tail
	zeroNumber := 0
	for i := coinValueLen - 1; i > 0; i-- {
		if coinValue[i] == '0' {
			zeroNumber++
		} else {
			break
		}
		if zeroNumber == unitBit {
			break
		}
	}

	//log.Info("the coinValue is:","coinValue",coinValue,"coinValueLen",coinValueLen,"zeroNumber",zeroNumber)


	coinValue = coinValue[:coinValueLen-zeroNumber]
	if coinValueLen <= unitBit {
		padding := make([]byte, unitBit-coinValueLen)
		for index := range padding {
			padding[index] = '0'
		}
		tmpBytes := append(padding[:], coinValue[:]...)
		//log.Info("the tmpBytes is:","tmpBytes",tmpBytes,"string",string(tmpBytes[:]))
		return "0." + string(tmpBytes[:]), nil
	} else {
		scalingPos := coinValueLen - unitBit
		if zeroNumber >= unitBit {
			return coinValue[:scalingPos], nil
		} else {
			return coinValue[:scalingPos] + "." + coinValue[scalingPos:], nil
		}
	}
}

//check and change input money value
//input money unit is DIP
func MoneyValueToCSCoin(moneyValue string) (*big.Int, error) {
	return DecimalToInter(moneyValue, consts.DIPDecimalBits)
}

//CSCoin to money Value
func CSCoinToMoneyValue(csCoinValue *hexutil.Big) (string, error) {
	return InterToDecimal(csCoinValue, consts.DIPDecimalBits)
}
